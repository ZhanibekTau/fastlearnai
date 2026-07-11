package api

import (
	"net/http"

	"github.com/fastlearn/backend/api/handlers"
	mw "github.com/fastlearn/backend/api/middleware"
	"github.com/fastlearn/backend/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(
	authH *handlers.AuthHandler,
	videoH *handlers.VideoHandler,
	subH *handlers.SubscriptionHandler,
	adminH *handlers.AdminHandler,
	cfg *config.Config,
) *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(mw.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	e.Use(mw.NewRateLimiter())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	v1 := e.Group("/fastlearn/v1")

	// ── Auth (public) ────────────────────────────────────────────
	auth := v1.Group("/auth")
	auth.POST("/register", authH.Register)              // → returns access_token + refresh_token
	auth.POST("/login", authH.Login)                    // → returns access_token + refresh_token
	auth.POST("/refresh", authH.Refresh)                // → new token pair via refresh_token
	auth.POST("/forgot-password", authH.ForgotPassword) // → generates + emails a new password

	// ── Auth (protected) ─────────────────────────────────────────
	authP := v1.Group("/auth", mw.JWTAuth(cfg))
	authP.GET("/me", authH.Me) // current user profile
	authP.POST("/logout", authH.Logout)

	// ── Subscriptions (public) ───────────────────────────────────
	subs := v1.Group("/subscriptions")
	subs.GET("/plans", subH.Plans)

	// ── Billing (public webhook) ─────────────────────────────────
	v1.POST("/billing/webhook", subH.StripeWebhook)

	// ── Videos (public) ─────────────────────────────────────────
	v1.GET("/videos/:id/public", videoH.GetPublic)

	// ── Protected ────────────────────────────────────────────────
	p := v1.Group("", mw.JWTAuth(cfg))

	// Videos
	p.POST("/videos", videoH.Upload)
	p.GET("/videos", videoH.List)
	p.GET("/videos/:id", videoH.Get)
	p.DELETE("/videos/:id", videoH.Delete)
	p.POST("/videos/:id/regenerate", videoH.Regenerate)
	p.GET("/videos/:id/download", videoH.Download)

	// Subscriptions
	p.GET("/subscriptions/me", subH.Mine)
	p.GET("/subscriptions/usage", subH.Usage)
	p.GET("/subscriptions/upgrade", subH.Upgrade)

	// ── Admin (manual plan grants, gated by ADMIN_EMAILS) ─────────
	admin := v1.Group("/admin", mw.JWTAuth(cfg), mw.AdminOnly(cfg))
	admin.GET("/stats", adminH.Stats)
	admin.GET("/users", adminH.SearchUsers)
	admin.POST("/users/:id/subscription", adminH.GrantSubscription)

	return e
}
