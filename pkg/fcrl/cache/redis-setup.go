package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisConfig struct {
	redisClient *redis.Client
	ctx         context.Context
}

func NewRedisConfig(host string, port int, pass string, ctx context.Context) *RedisConfig {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: pass,
		DB:       0,
	})

	return &RedisConfig{
		redisClient: rdb,
		ctx:         ctx,
	}
}

func (r *RedisConfig) Set(key string, value string, ttl int) error {
	return r.redisClient.Set(r.ctx, key, value, time.Duration(ttl)*time.Second).Err()
}

func (r *RedisConfig) Get(key string) (string, error) {
	return r.redisClient.Get(r.ctx, key).Result()
}

func (r *RedisConfig) Delete(key string) error {
	return r.redisClient.Del(r.ctx, key).Err()
}
