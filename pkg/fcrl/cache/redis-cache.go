package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	redisClient *redis.Client
	ctx         context.Context
}

func NewRedisCache(host string, port int, pass string, ctx context.Context) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: pass,
		DB:       0,
	})

	return &RedisCache{
		redisClient: rdb,
		ctx:         ctx,
	}
}

func NewRedisClientCache(rdb *redis.Client, ctx context.Context) *RedisCache {
	return &RedisCache{
		redisClient: rdb,
		ctx:         ctx,
	}
}

func (r *RedisCache) Set(key string, value string, ttl time.Duration) error {
	return r.redisClient.Set(r.ctx, key, value, ttl).Err()
}

func (r *RedisCache) Get(key string) (string, error) {
	return r.redisClient.Get(r.ctx, key).Result()
}

func (r *RedisCache) Update(key string, value string, ttl time.Duration) error {
	if ttl != -1 {
		return r.redisClient.SetXX(r.ctx, key, value, ttl).Err()
	}
	return r.redisClient.SetXX(r.ctx, key, value, -1).Err()
}

func (r *RedisCache) Delete(key string) error {
	return r.redisClient.Del(r.ctx, key).Err()
}
