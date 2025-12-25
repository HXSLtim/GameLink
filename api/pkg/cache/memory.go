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
	stopCh chan struct{}
}

// NewMemory 创建内存缓存。
func NewMemory() Cache {
	c := &memoryCache{
		values: make(map[string]memoryEntry),
		stopCh: make(chan struct{}),
	}
	go c.janitor(1 * time.Minute)
	return c
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
	close(c.stopCh)
	c.mu.Lock()
	c.values = make(map[string]memoryEntry)
	c.mu.Unlock()
	return nil
}

func (c *memoryCache) janitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			c.mu.Lock()
			for k, v := range c.values {
				if !v.expiry.IsZero() && now.After(v.expiry) {
					delete(c.values, k)
				}
			}
			c.mu.Unlock()
		case <-c.stopCh:
			return
		}
	}
}
