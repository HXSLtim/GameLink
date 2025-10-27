package cache

import (
	"context"
	"time"

	"gamelink/internal/config"
)

// Cache 提供简单的 KV 缓存接口。
type Cache interface {
	Get(ctx context.Context, key string) (value string, ok bool, err error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Close(ctx context.Context) error
}

// New 根据配置创建缓存实例。
func New(cfg config.CacheConfig) (Cache, error) {
	switch cfg.Type {
	case "", "memory":
		return NewMemory(), nil
	case "redis":
		return NewRedis(cfg.Redis)
	default:
		return NewMemory(), nil
	}
}
