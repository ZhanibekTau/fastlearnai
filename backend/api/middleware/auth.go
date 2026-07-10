package middleware

import (
	"strings"

	"github.com/fastlearn/backend/config"
	"github.com/fastlearn/backend/utils"
	"github.com/labstack/echo/v4"
)

const UserIDKey = "user_id"
const UserEmailKey = "user_email"

func JWTAuth(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				return utils.Unauthorized(c)
			}
			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims, err := utils.ParseToken(tokenStr, cfg.JWTSecret)
			if err != nil {
				return utils.Unauthorized(c)
			}
			c.Set(UserIDKey, claims.UserID)
			c.Set(UserEmailKey, claims.Email)
			return next(c)
		}
	}
}

func GetUserID(c echo.Context) string {
	id, _ := c.Get(UserIDKey).(string)
	return id
}

func GetUserEmail(c echo.Context) string {
	email, _ := c.Get(UserEmailKey).(string)
	return email
}

// AdminOnly must run after JWTAuth. It rejects any authenticated user whose
// email is not in cfg.AdminEmails.
func AdminOnly(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !cfg.IsAdminEmail(GetUserEmail(c)) {
				return utils.Forbidden(c, "admin access required")
			}
			return next(c)
		}
	}
}
