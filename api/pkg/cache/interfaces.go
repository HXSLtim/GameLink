package cache

// CacheConfigProvider defines the cache configuration interface
type CacheConfigProvider interface {
	GetCacheType() string
	GetRedisConfig() RedisConfig
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}
