package cache

import (
	"context"
	"sync"
	"time"
)

type memoryEntry struct {
	value  string
	expiry time.Time
}

// memoryCache 是开发环境下的本地缓存实现。
type memoryCache struct {
	mu     sync.RWMutex
	values map[string]memoryEntry
}

// NewMemory 创建内存缓存。
func NewMemory() Cache {
	return &memoryCache{
		values: make(map[string]memoryEntry),
	}
}

func (c *memoryCache) Get(_ context.Context, key string) (string, bool, error) {
	c.mu.RLock()
	entry, ok := c.values[key]
	c.mu.RUnlock()
	if !ok {
		return "", false, nil
	}
	if !entry.expiry.IsZero() && time.Now().After(entry.expiry) {
		c.mu.Lock()
		delete(c.values, key)
		c.mu.Unlock()
		return "", false, nil
	}
	return entry.value, true, nil
}

func (c *memoryCache) Set(_ context.Context, key, value string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := memoryEntry{value: value}
	if ttl > 0 {
		entry.expiry = time.Now().Add(ttl)
	}
	c.values[key] = entry
	return nil
}

func (c *memoryCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	delete(c.values, key)
	c.mu.Unlock()
	return nil
}

func (c *memoryCache) Close(context.Context) error {
	c.mu.Lock()
	c.values = make(map[string]memoryEntry)
	c.mu.Unlock()
	return nil
}
