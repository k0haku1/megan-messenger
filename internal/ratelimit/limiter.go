package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *Limiter {
	return &Limiter{rdb: rdb}
}

func (l *Limiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	if l.rdb == nil {
		return true, nil
	}

	redisKey := fmt.Sprintf("ratelimit:%s", key)
	count, err := l.rdb.Incr(ctx, redisKey).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		if err := l.rdb.Expire(ctx, redisKey, window).Err(); err != nil {
			return false, err
		}
	}

	return count <= int64(limit), nil
}
