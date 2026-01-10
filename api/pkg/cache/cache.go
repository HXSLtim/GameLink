package cache

import (
	"context"
	"time"
)

// Cache 提供简单的 KV 缓存接口。
type Cache interface {
	Get(ctx context.Context, key string) (value string, ok bool, err error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Close(ctx context.Context) error

	// GetRedisClient returns the underlying Redis client if using Redis cache,
	// or nil if using memory cache. This is used for WebSocket Pub/Sub.
	GetRedisClient() interface{}
}

// New 根据配置创建缓存实例。
func New(cfg CacheConfigProvider) (Cache, error) {
	switch cfg.GetCacheType() {
	case "", "memory":
		return NewMemory(), nil
	case "redis":
		return NewRedis(cfg.GetRedisConfig())
	default:
		return NewMemory(), nil
	}
}
