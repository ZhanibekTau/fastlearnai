package middleware

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	ansiReset   = "\033[0m"
	ansiBold    = "\033[1m"
	ansiDim     = "\033[2m"
	ansiRed     = "\033[31m"
	ansiGreen   = "\033[32m"
	ansiYellow  = "\033[33m"
	ansiBlue    = "\033[34m"
	ansiMagenta = "\033[35m"
	ansiCyan    = "\033[36m"
	ansiGray    = "\033[90m"
)

func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)

			path := c.Request().URL.Path
			if path == "/health" {
				return err
			}

			status := c.Response().Status
			method := c.Request().Method
			elapsed := time.Since(start).Round(time.Millisecond)

			fmt.Printf("%s%s%s  %s%-7s%s  %-44s  %s%s%3d%s  %s%v%s\n",
				ansiGray, start.Format("15:04:05"), ansiReset,
				methodColor(method), method, ansiReset,
				path,
				statusColor(status), ansiBold, status, ansiReset,
				ansiGray, elapsed, ansiReset,
			)

			return err
		}
	}
}

func methodColor(m string) string {
	switch m {
	case "GET":
		return ansiBlue
	case "POST":
		return ansiGreen
	case "PUT", "PATCH":
		return ansiYellow
	case "DELETE":
		return ansiRed
	}
	return ansiDim
}

func statusColor(s int) string {
	switch {
	case s < 300:
		return ansiGreen
	case s < 400:
		return ansiCyan
	case s < 500:
		return ansiYellow
	default:
		return ansiRed
	}
}
