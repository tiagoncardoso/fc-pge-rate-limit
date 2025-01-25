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

func (r *RedisConfig) Set(key string, value string, ttl time.Duration) error {
	return r.redisClient.Set(r.ctx, key, value, ttl).Err()
}

func (r *RedisConfig) Get(key string) (string, error) {
	return r.redisClient.Get(r.ctx, key).Result()
}

func (r *RedisConfig) Update(key string, value string, ttl time.Duration) error {
	if ttl != -1 {
		return r.redisClient.SetXX(r.ctx, key, value, ttl).Err()
	}
	return r.redisClient.SetXX(r.ctx, key, value, -1).Err()
}

func (r *RedisConfig) Delete(key string) error {
	return r.redisClient.Del(r.ctx, key).Err()
}
