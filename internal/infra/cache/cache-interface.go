package cache

type CacheInterface interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl int) error
	Delete(key string) error
}
