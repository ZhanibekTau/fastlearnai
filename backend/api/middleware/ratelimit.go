package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type rateLimiter struct {
	mu       sync.Mutex
	counters map[string]*counter
}

type counter struct {
	count    int
	resetAt  time.Time
}

func NewRateLimiter() echo.MiddlewareFunc {
	rl := &rateLimiter{counters: make(map[string]*counter)}

	go func() {
		for range time.Tick(time.Minute) {
			rl.mu.Lock()
			now := time.Now()
			for k, c := range rl.counters {
				if now.After(c.resetAt) {
					delete(rl.counters, k)
				}
			}
			rl.mu.Unlock()
		}
	}()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := c.RealIP()
			rl.mu.Lock()
			cnt, ok := rl.counters[key]
			if !ok || time.Now().After(cnt.resetAt) {
				rl.counters[key] = &counter{count: 1, resetAt: time.Now().Add(time.Minute)}
				rl.mu.Unlock()
				return next(c)
			}
			cnt.count++
			limit := cnt.count
			rl.mu.Unlock()

			if limit > 100 {
				return c.JSON(http.StatusTooManyRequests, map[string]string{"message": "rate limit exceeded"})
			}
			return next(c)
		}
	}
}
