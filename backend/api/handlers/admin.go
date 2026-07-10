package handlers

import (
	"net/http"

	"github.com/fastlearn/backend/models"
	"github.com/fastlearn/backend/services"
	"github.com/fastlearn/backend/utils"
	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	authSvc *services.AuthService
	subSvc  *services.SubscriptionService
}

func NewAdminHandler(authSvc *services.AuthService, subSvc *services.SubscriptionService) *AdminHandler {
	return &AdminHandler{authSvc: authSvc, subSvc: subSvc}
}

type adminUserResult struct {
	*models.User
	Subscription *models.Subscription `json:"subscription"`
}

// SearchUsers handles GET /admin/users?email=...
func (h *AdminHandler) SearchUsers(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
		return utils.BadRequest(c, "email query param required")
	}

	ctx := c.Request().Context()
	users, err := h.authSvc.SearchUsers(ctx, email)
	if err != nil {
		return utils.InternalError(c)
	}

	results := make([]adminUserResult, 0, len(users))
	for i := range users {
		sub, err := h.subSvc.GetMine(ctx, users[i].ID)
		if err != nil {
			return utils.InternalError(c)
		}
		results = append(results, adminUserResult{User: &users[i], Subscription: sub})
	}
	return c.JSON(http.StatusOK, results)
}

type grantPlanRequest struct {
	Plan   string `json:"plan" validate:"required"`
	Months int    `json:"months"`
}

// GrantSubscription handles POST /admin/users/:id/subscription
func (h *AdminHandler) GrantSubscription(c echo.Context) error {
	userID := c.Param("id")

	var req grantPlanRequest
	if err := c.Bind(&req); err != nil || req.Plan == "" {
		return utils.BadRequest(c, "plan is required")
	}

	sub, err := h.subSvc.AdminGrantPlan(c.Request().Context(), userID, models.Plan(req.Plan), req.Months)
	if err != nil {
		if err == services.ErrInvalidPlan {
			return utils.BadRequest(c, "invalid plan")
		}
		return utils.InternalError(c)
	}
	return c.JSON(http.StatusOK, sub)
}
