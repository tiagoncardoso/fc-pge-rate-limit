package cache

import "time"

type CacheInterface interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl time.Duration) error
	Update(key string, value string, ttl time.Duration) error
	Delete(key string) error
}
