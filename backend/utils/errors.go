package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrorResponse(c echo.Context, status int, message string) error {
	return c.JSON(status, APIError{Code: status, Message: message})
}

func BadRequest(c echo.Context, message string) error {
	return ErrorResponse(c, http.StatusBadRequest, message)
}

func Unauthorized(c echo.Context) error {
	return ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
}

func Forbidden(c echo.Context, message string) error {
	return ErrorResponse(c, http.StatusForbidden, message)
}

func NotFound(c echo.Context, message string) error {
	return ErrorResponse(c, http.StatusNotFound, message)
}

func InternalError(c echo.Context) error {
	return ErrorResponse(c, http.StatusInternalServerError, "internal server error")
}
