package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// testCacheConfigProvider is a test implementation of CacheConfigProvider
type testCacheConfigProvider struct {
	cacheType string
	redisCfg  RedisConfig
}

func (t *testCacheConfigProvider) GetCacheType() string {
	return t.cacheType
}

func (t *testCacheConfigProvider) GetRedisConfig() RedisConfig {
	return t.redisCfg
}

/**
 * Test_New_WithMemoryType 测试创建内存缓存
 */
func Test_New_WithMemoryType(t *testing.T) {
	// Arrange
	cfg := &testCacheConfigProvider{cacheType: "memory"}

	// Act
	cache, err := New(cfg)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, cache)
	assert.IsType(t, &memoryCache{}, cache)
}

/**
 * Test_New_WithEmptyType 测试空类型默认为内存缓存
 */
func Test_New_WithEmptyType(t *testing.T) {
	// Arrange
	cfg := &testCacheConfigProvider{cacheType: ""}

	// Act
	cache, err := New(cfg)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, cache)
	assert.IsType(t, &memoryCache{}, cache)
}

/**
 * Test_New_WithUnknownType 测试未知类型默认为内存缓存
 */
func Test_New_WithUnknownType(t *testing.T) {
	// Arrange
	cfg := &testCacheConfigProvider{cacheType: "unknown"}

	// Act
	cache, err := New(cfg)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, cache)
	assert.IsType(t, &memoryCache{}, cache)
}

/**
 * Test_New_WithRedisType_InvalidConfig 测试Redis配置错误
 */
func Test_New_WithRedisType_InvalidConfig(t *testing.T) {
	// Arrange
	cfg := &testCacheConfigProvider{
		cacheType: "redis",
		redisCfg: RedisConfig{
			Addr: "invalid:99999", // 无效地址
		},
	}

	// Act
	cache, err := New(cfg)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, cache)
}
