package handlers

import (
	"net/http"

	"github.com/fastlearn/backend/api/middleware"
	"github.com/fastlearn/backend/services"
	"github.com/fastlearn/backend/utils"
	"github.com/labstack/echo/v4"
)

type SubscriptionHandler struct {
	svc *services.SubscriptionService
}

func NewSubscriptionHandler(svc *services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc}
}

func (h *SubscriptionHandler) Plans(c echo.Context) error {
	return c.JSON(http.StatusOK, []map[string]any{
		{
			"id": "free", "name": "Free", "price": 0,
			"videos_month": 3, "max_duration_min": 30,
			"features": []string{"HTML export", "7-day history"},
		},
		{
			"id": "pro", "name": "Pro", "price": 1500,
			"videos_month": 30, "max_duration_min": 180,
			"features": []string{"PDF + Markdown", "Q&A generation", "30-day history"},
		},
		{
			"id": "premium", "name": "Premium", "price": 3900,
			"videos_month": -1, "max_duration_min": -1,
			"features": []string{"Unlimited", "API access", "JSON export", "Lifetime history"},
		},
	})
}

func (h *SubscriptionHandler) Mine(c echo.Context) error {
	userID := middleware.GetUserID(c)

	sub, err := h.svc.GetMine(c.Request().Context(), userID)
	if err != nil {
		return utils.InternalError(c)
	}
	return c.JSON(http.StatusOK, sub)
}

func (h *SubscriptionHandler) Usage(c echo.Context) error {
	userID := middleware.GetUserID(c)

	usage, err := h.svc.GetUsage(c.Request().Context(), userID)
	if err != nil {
		return utils.InternalError(c)
	}
	return c.JSON(http.StatusOK, usage)
}

func (h *SubscriptionHandler) Upgrade(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"contact": map[string]string{
			"type":  "whatsapp",
			"phone": "+7 747 925 50 91",
			"url":   "https://wa.me/77479255091",
		},
		"plans": []map[string]string{
			{"name": "Pro", "price": "1500 ₸/month", "videos": "30/month"},
			{"name": "Premium", "price": "3900 ₸/month", "videos": "Unlimited"},
		},
	})
}

func (h *SubscriptionHandler) StripeWebhook(c echo.Context) error {
	if err := h.svc.HandleStripeWebhook(c.Request()); err != nil {
		return utils.BadRequest(c, "invalid webhook payload")
	}
	return c.JSON(http.StatusOK, map[string]string{"received": "true"})
}
