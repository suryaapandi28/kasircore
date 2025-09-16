package utils

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// ResetPenalty akan menghapus penalty untuk IP tertentu di Redis
func ResetPenalty(rdb *redis.Client, ip string) error {
	_, err := rdb.Del(ctx, "rate_penalty:"+ip, "rate_penalty_count:"+ip).Result()
	return err
}
