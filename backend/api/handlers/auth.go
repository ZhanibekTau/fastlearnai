package handlers

import (
	"net/http"

	"github.com/fastlearn/backend/api/middleware"
	"github.com/fastlearn/backend/config"
	"github.com/fastlearn/backend/models"
	"github.com/fastlearn/backend/services"
	"github.com/fastlearn/backend/utils"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	svc *services.AuthService
	cfg *config.Config
}

func NewAuthHandler(svc *services.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{svc: svc, cfg: cfg}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}

	resp, err := h.svc.Register(c.Request().Context(), req)
	if err != nil {
		if err == services.ErrEmailExists {
			return utils.BadRequest(c, "email already registered")
		}
		return utils.InternalError(c)
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}

	resp, err := h.svc.Login(c.Request().Context(), req)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			return utils.Unauthorized(c)
		}
		return utils.InternalError(c)
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.Bind(&body); err != nil || body.RefreshToken == "" {
		return utils.BadRequest(c, "refresh_token required")
	}

	resp, err := h.svc.Refresh(c.Request().Context(), body.RefreshToken)
	if err != nil {
		return utils.Unauthorized(c)
	}
	return c.JSON(http.StatusOK, resp)
}

// Me returns the currently authenticated user. Requires JWT.
func (h *AuthHandler) Me(c echo.Context) error {
	userID := middleware.GetUserID(c)
	user, err := h.svc.Me(c.Request().Context(), userID)
	if err != nil {
		return utils.NotFound(c, "user not found")
	}
	return c.JSON(http.StatusOK, map[string]any{
		"id":             user.ID,
		"email":          user.Email,
		"name":           user.Name,
		"oauth_provider": user.OAuthProvider,
		"created_at":     user.CreatedAt,
		"updated_at":     user.UpdatedAt,
		"is_admin":       h.cfg.IsAdminEmail(user.Email),
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	// Stateless JWT: client drops the token. Server-side blacklist is Phase 2.
	return c.JSON(http.StatusOK, map[string]string{"message": "logged out"})
}
