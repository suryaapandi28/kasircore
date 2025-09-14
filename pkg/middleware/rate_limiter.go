package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

var ctx = context.Background()

// NewRateLimiter bikin middleware limit request
func NewRateLimiter(rdb *redis.Client, limit int, window time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			key := "rate_limit:" + ip

			// Tambah 1 hit ke redis
			count, err := rdb.Incr(ctx, key).Result()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "server error"})
			}

			if count == 1 {
				// set expiry pertama kali
				rdb.Expire(ctx, key, window)
			}

			if count > int64(limit) {
				// terlalu banyak request
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Terlalu banyak request, coba lagi nanti",
				})
			}

			return next(c)
		}
	}
}
