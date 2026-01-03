// package middleware

// import (
// 	"net/http"
// 	"time"

// 	"github.com/labstack/echo/v4"
// 	"github.com/redis/go-redis/v9"
// 	"github.com/suryaapandi28/kasircore/pkg/response"
// 	"golang.org/x/net/context"
// )

// var ctx = context.Background()

// // NewRateLimiter bikin middleware limit request
// func NewRateLimiter(rdb *redis.Client, limit int, window time.Duration) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			ip := c.RealIP()
// 			key := "rate_limit:" + ip

// 			// Tambah 1 hit ke redis
// 			count, err := rdb.Incr(ctx, key).Result()
// 			if err != nil {
// 				// return c.JSON(http.StatusInternalServerError, map[string]string{"error": "server error"})
// 				return c.JSON(http.StatusInternalServerError, response.ErrorResponse(429, "server error"))
// 			}

// 			if count == 1 {
// 				// set expiry pertama kali
// 				rdb.Expire(ctx, key, window)
// 			}

// 			if count > int64(limit) {
// 				// terlalu banyak request
// 				// return c.JSON(http.StatusTooManyRequests, map[string]string{
// 				// 	"error": "Terlalu banyak request, coba lagi nanti",

// 				// })
// 				return c.JSON(http.StatusInternalServerError, response.ErrorResponse(429, "Terlalu banyak request, coba lagi nanti"))
// 			}

// 			return next(c)
// 		}
// 	}
// }

package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/suryaapandi28/kasircore/pkg/response"
	"golang.org/x/net/context"
)

var ctx = context.Background()

// NewRateLimiter bikin middleware limit request
func NewRateLimiter(rdb *redis.Client, limit int, window time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			key := "rate_limit:" + ip
			penaltyKey := "rate_penalty:" + ip
			penaltyCountKey := "rate_penalty_count:" + ip

			// cek apakah ip masih dalam penalty
			exists, _ := rdb.Exists(ctx, penaltyKey).Result()
			if exists > 0 {
				return c.JSON(http.StatusTooManyRequests,
					response.ErrorResponse(http.StatusTooManyRequests,
						"Anda diblokir sementara karena terlalu banyak request"))
			}

			// hitung request
			count, err := rdb.Incr(ctx, key).Result()
			if err != nil {
				return c.JSON(http.StatusInternalServerError,
					response.ErrorResponse(http.StatusInternalServerError, "server error"))
			}

			if count == 1 {
				// set expiry window (misal 10 detik)
				rdb.Expire(ctx, key, window)
			}

			if count > int64(limit) {
				// ambil jumlah pelanggaran sebelumnya
				violations, _ := rdb.Get(ctx, penaltyCountKey).Int()
				violations++

				// tentukan durasi penalty
				var penaltyDuration time.Duration
				switch violations {
				case 1:
					penaltyDuration = 1 * time.Minute
				case 2:
					penaltyDuration = 5 * time.Minute
				case 3:
					penaltyDuration = 15 * time.Minute
				case 4:
					penaltyDuration = 1 * time.Hour
				default:
					penaltyDuration = 0 // banned permanen
				}

				// set penalty
				if penaltyDuration > 0 {
					// banned sementara
					rdb.Set(ctx, penaltyKey, 1, penaltyDuration)
				} else {
					// banned permanen (tanpa expire)
					rdb.Set(ctx, penaltyKey, 1, 0)
				}

				// update jumlah pelanggaran (reset setiap 24 jam kalau mau)
				rdb.Set(ctx, penaltyCountKey, violations, 24*time.Hour)

				// respon error
				if penaltyDuration > 0 {
					return c.JSON(http.StatusTooManyRequests,
						response.ErrorResponse(http.StatusTooManyRequests,
							"Terlalu banyak request, coba lagi nanti"))
				} else {
					return c.JSON(http.StatusTooManyRequests,
						response.ErrorResponse(http.StatusTooManyRequests,
							"Anda diblokir permanen karena terlalu banyak request"))
				}
			}

			return next(c)
		}
	}
}
