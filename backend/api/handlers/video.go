package handlers

import (
	"net/http"
	"strconv"

	"github.com/fastlearn/backend/api/middleware"
	"github.com/fastlearn/backend/services"
	"github.com/fastlearn/backend/utils"
	"github.com/labstack/echo/v4"
)

type VideoHandler struct {
	svc *services.VideoService
}

func NewVideoHandler(svc *services.VideoService) *VideoHandler {
	return &VideoHandler{svc: svc}
}

func (h *VideoHandler) Upload(c echo.Context) error {
	userID := middleware.GetUserID(c)

	file, header, _ := c.Request().FormFile("file")
	youtubeURL := c.FormValue("youtube_url")
	language := c.FormValue("language")
	if language == "" {
		language = "auto"
	}

	if file == nil && youtubeURL == "" {
		return utils.BadRequest(c, "provide a file or youtube_url")
	}

	video, err := h.svc.Upload(c.Request().Context(), services.UploadInput{
		UserID:     userID,
		File:       file,
		FileHeader: header,
		YouTubeURL: youtubeURL,
		Language:   language,
	})
	if err != nil {
		if err == services.ErrLimitReached {
			return c.JSON(http.StatusPaymentRequired, map[string]any{
				"message": "Monthly limit reached",
				"upgrade": map[string]string{
					"type": "whatsapp",
					"url":  "https://wa.me/77479255091",
				},
			})
		}
		return utils.InternalError(c)
	}
	return c.JSON(http.StatusCreated, video)
}

func (h *VideoHandler) List(c echo.Context) error {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.QueryParam("page"))

	result, err := h.svc.List(c.Request().Context(), userID, page, 20)
	if err != nil {
		return utils.InternalError(c)
	}
	return c.JSON(http.StatusOK, result)
}

func (h *VideoHandler) Get(c echo.Context) error {
	userID := middleware.GetUserID(c)

	video, err := h.svc.Get(c.Request().Context(), c.Param("id"), userID)
	if err != nil {
		if err == services.ErrNotFound {
			return utils.NotFound(c, "video not found")
		}
		return utils.InternalError(c)
	}
	return c.JSON(http.StatusOK, video)
}

func (h *VideoHandler) Delete(c echo.Context) error {
	userID := middleware.GetUserID(c)

	if err := h.svc.Delete(c.Request().Context(), c.Param("id"), userID); err != nil {
		if err == services.ErrNotFound {
			return utils.NotFound(c, "video not found")
		}
		return utils.InternalError(c)
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "deleted"})
}

func (h *VideoHandler) Regenerate(c echo.Context) error {
	userID := middleware.GetUserID(c)

	if err := h.svc.Regenerate(c.Request().Context(), c.Param("id"), userID); err != nil {
		if err == services.ErrNotFound {
			return utils.NotFound(c, "video not found")
		}
		return utils.InternalError(c)
	}
	return c.JSON(http.StatusAccepted, map[string]string{"message": "regeneration queued"})
}

func (h *VideoHandler) Download(c echo.Context) error {
	userID := middleware.GetUserID(c)
	format := c.QueryParam("format")
	if format == "" {
		format = "pdf"
	}

	data, filename, err := h.svc.Download(c.Request().Context(), c.Param("id"), userID, format)
	if err != nil {
		switch err {
		case services.ErrNotFound:
			return utils.NotFound(c, "video not found")
		case services.ErrUnsupportedFormat:
			return utils.BadRequest(c, "unsupported export format")
		case services.ErrUpgradeRequired:
			return c.JSON(http.StatusPaymentRequired, map[string]any{
				"message": "PDF export is available on Pro and Premium plans",
				"upgrade": map[string]string{
					"type": "whatsapp",
					"url":  "https://wa.me/77479255091",
				},
			})
		default:
			return utils.InternalError(c)
		}
	}

	c.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="`+filename+`"`)
	return c.Blob(http.StatusOK, "application/pdf", data)
}

func (h *VideoHandler) GetPublic(c echo.Context) error {
	video, err := h.svc.GetPublic(c.Request().Context(), c.Param("id"))
	if err != nil {
		return utils.NotFound(c, "not found")
	}
	return c.JSON(http.StatusOK, video)
}
