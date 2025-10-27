package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"gamelink/internal/config"
)

// redisCache 是生产环境下的 Redis 缓存实现。
type redisCache struct {
	client *redis.Client
}

// NewRedis 创建 Redis 缓存实例。
func NewRedis(cfg config.RedisConfig) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}

	return &redisCache{client: client}, nil
}

func (c *redisCache) Get(ctx context.Context, key string) (string, bool, error) {
	result, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return result, true, nil
}

func (c *redisCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *redisCache) Close(context.Context) error {
	return c.client.Close()
}
