package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type MemoryCache struct {
	cache map[string]cacheItem
	mu    sync.RWMutex
	ctx   context.Context
}

type cacheItem struct {
	value      string
	expiration int64
}

var (
	instance *MemoryCache
	once     sync.Once
)

func NewMemoryCache(ctx context.Context) *MemoryCache {
	once.Do(func() {
		instance = &MemoryCache{
			cache: make(map[string]cacheItem),
			ctx:   ctx,
		}
	})
	return instance
}

func (m *MemoryCache) Set(key string, value string, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	expiration := time.Now().Add(ttl).UnixNano()

	m.cache[key] = cacheItem{
		value:      value,
		expiration: expiration,
	}

	return nil
}

func (m *MemoryCache) Get(key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, found := m.cache[key]
	if !found || time.Now().UnixNano() > item.expiration {
		return "", fmt.Errorf("key not found or expired")
	}

	return item.value, nil
}

func (m *MemoryCache) Update(key string, value string, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, found := m.cache[key]; !found {
		return fmt.Errorf("key not found")
	}

	if ttl != -1 {
		expiration := time.Now().Add(ttl).UnixNano()
		m.cache[key] = cacheItem{
			value:      value,
			expiration: expiration,
		}
	}

	m.cache[key] = cacheItem{
		value:      value,
		expiration: m.cache[key].expiration,
	}
	return nil
}

func (m *MemoryCache) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, found := m.cache[key]; !found {
		return fmt.Errorf("key not found")
	}
	delete(m.cache, key)
	return nil
}
