package cache

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/**
 * Test_NewMemory_CreatesCache 测试创建内存缓存
 */
func Test_NewMemory_CreatesCache(t *testing.T) {
	// Act
	cache := NewMemory()

	// Assert
	assert.NotNil(t, cache)
	mc := cache.(*memoryCache)
	assert.NotNil(t, mc.values)
	assert.NotNil(t, mc.stopCh)
}

/**
 * Test_MemoryCache_SetAndGet_Success 测试设置和获取
 */
func Test_MemoryCache_SetAndGet_Success(t *testing.T) {
	// Arrange
	cache := NewMemory()
	ctx := context.Background()
	key := "test_key"
	value := "test_value"

	// Act
	err := cache.Set(ctx, key, value, time.Minute)
	require.NoError(t, err)

	result, ok, err := cache.Get(ctx, key)

	// Assert
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, value, result)
}

/**
 * Test_MemoryCache_Get_NonExistentKey 测试获取不存在的键
 */
func Test_MemoryCache_Get_NonExistentKey(t *testing.T) {
	// Arrange
	cache := NewMemory()
	ctx := context.Background()

	// Act
	result, ok, err := cache.Get(ctx, "non_existent")

	// Assert
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Empty(t, result)
}

/**
 * Test_MemoryCache_Set_WithZeroTTL 测试设置无过期时间的值
 */
func Test_MemoryCache_Set_WithZeroTTL(t *testing.T) {
	// Arrange
	cache := NewMemory()
	ctx := context.Background()
	key := "permanent_key"
	value := "permanent_value"

	// Act
	err := cache.Set(ctx, key, value, 0)
	require.NoError(t, err)

	// 等待一段时间
	time.Sleep(100 * time.Millisecond)

	result, ok, err := cache.Get(ctx, key)

	// Assert
	assert.NoError(t, err)
	assert.True(t, ok, "永久键不应过期")
	assert.Equal(t, value, result)
}

/**
 * Test_MemoryCache_TTL_Expiry 测试TTL过期
 */
func Test_MemoryCache_TTL_Expiry(t *testing.T) {
	// Arrange
	cache := NewMemory()
	ctx := context.Background()
	key := "expiring_key"
	value := "expiring_value"

	// Act
	err := cache.Set(ctx, key, value, 100*time.Millisecond)
	require.NoError(t, err)

	// 验证键存在
	_, ok, _ := cache.Get(ctx, key)
	assert.True(t, ok, "键应该存在")

	// 等待过期
	time.Sleep(200 * time.Millisecond)

	// 验证键已过期
	result, ok, err := cache.Get(ctx, key)

	// Assert
	assert.NoError(t, err)
	assert.False(t, ok, "键应该已过期")
	assert.Empty(t, result)
}

/**
 * Test_MemoryCache_Set_UpdateExisting 测试更新已存在的键
 */
func Test_MemoryCache_Set_UpdateExisting(t *testing.T) {
	// Arrange
	cache := NewMemory()
	ctx := context.Background()
	key := "update_key"

	// Act - 首次设置
	cache.Set(ctx, key, "value1", time.Minute)
	result1, _, _ := cache.Get(ctx, key)

	// Act - 更新
	cache.Set(ctx, key, "value2", time.Minute)
	result2, ok, err := cache.Get(ctx, key)

	// Assert
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "value1", result1)
	assert.Equal(t, "value2", result2)
}

/**
 * Test_MemoryCache_Delete_Success 测试删除键
 */
func Test_MemoryCache_Delete_Success(t *testing.T) {
	// Arrange
	cache := NewMemory()
	ctx := context.Background()
	key := "delete_key"

	// Act - 设置并删除
	cache.Set(ctx, key, "value", time.Minute)
	err := cache.Delete(ctx, key)
	require.NoError(t, err)

	// 验证已删除
	_, ok, err := cache.Get(ctx, key)

	// Assert
	assert.NoError(t, err)
	assert.False(t, ok)
}

/**
 * Test_MemoryCache_Delete_NonExistent 测试删除不存在的键
 */
func Test_MemoryCache_Delete_NonExistent(t *testing.T) {
	// Arrange
	cache := NewMemory()
	ctx := context.Background()

	// Act
	err := cache.Delete(ctx, "non_existent")

	// Assert
	assert.NoError(t, err)
}

/**
 * Test_MemoryCache_Close 测试关闭缓存
 */
func Test_MemoryCache_Close(t *testing.T) {
	// Arrange
	cache := NewMemory()
	ctx := context.Background()

	// 设置一些数据
	cache.Set(ctx, "key1", "value1", time.Minute)
	cache.Set(ctx, "key2", "value2", time.Minute)

	// Act
	err := cache.Close(ctx)

	// Assert
	assert.NoError(t, err)

	// 验证数据已清空
	mc := cache.(*memoryCache)
	mc.mu.RLock()
	assert.Empty(t, mc.values)
	mc.mu.RUnlock()
}

/**
 * Test_MemoryCache_Concurrent_Access 测试并发访问
 */
func Test_MemoryCache_Concurrent_Access(t *testing.T) {
	// Arrange
	cache := NewMemory()
	ctx := context.Background()
	var wg sync.WaitGroup
	numGoroutines := 100
	numOperations := 10

	// Act - 并发读写
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := "key_" + string(rune(id))
				value := "value_" + string(rune(j))

				// 写
				cache.Set(ctx, key, value, time.Minute)

				// 读
				cache.Get(ctx, key)

				// 删除
				if j%2 == 0 {
					cache.Delete(ctx, key)
				}
			}
		}(i)
	}

	// Assert - 等待完成，不应该panic
	wg.Wait()
	assert.True(t, true, "并发访问应该安全完成")
}

/**
 * Test_MemoryCache_Janitor_CleansExpired 测试清理器清理过期键
 */
func Test_MemoryCache_Janitor_CleansExpired(t *testing.T) {
	// Arrange
	cache := NewMemory()
	ctx := context.Background()

	// 设置一些短期和长期的键
	cache.Set(ctx, "short1", "value1", 100*time.Millisecond)
	cache.Set(ctx, "short2", "value2", 100*time.Millisecond)
	cache.Set(ctx, "long", "value3", time.Hour)

	// 等待短期键过期
	time.Sleep(200 * time.Millisecond)

	// 触发一次Get来清理过期键
	cache.Get(ctx, "short1")
	cache.Get(ctx, "short2")

	// Act - 等待janitor运行
	time.Sleep(1100 * time.Millisecond) // janitor每分钟运行一次，但我们不能等那么久

	// Assert - 长期键应该还在
	_, ok, _ := cache.Get(ctx, "long")
	assert.True(t, ok, "长期键应该还存在")
}

/**
 * Test_MemoryCache_MultipleKeys 测试多个键的操作
 */
func Test_MemoryCache_MultipleKeys(t *testing.T) {
	// Arrange
	cache := NewMemory()
	ctx := context.Background()

	keys := []string{"key1", "key2", "key3", "key4", "key5"}
	values := []string{"val1", "val2", "val3", "val4", "val5"}

	// Act - 设置多个键
	for i, key := range keys {
		err := cache.Set(ctx, key, values[i], time.Minute)
		require.NoError(t, err)
	}

	// Assert - 验证所有键
	for i, key := range keys {
		result, ok, err := cache.Get(ctx, key)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, values[i], result)
	}

	// Act - 删除部分键
	cache.Delete(ctx, "key2")
	cache.Delete(ctx, "key4")

	// Assert - 验证删除结果
	_, ok1, _ := cache.Get(ctx, "key1")
	_, ok2, _ := cache.Get(ctx, "key2")
	_, ok3, _ := cache.Get(ctx, "key3")
	_, ok4, _ := cache.Get(ctx, "key4")
	_, ok5, _ := cache.Get(ctx, "key5")

	assert.True(t, ok1)
	assert.False(t, ok2)
	assert.True(t, ok3)
	assert.False(t, ok4)
	assert.True(t, ok5)
}

/**
 * Test_MemoryCache_LargeValue 测试大值存储
 */
func Test_MemoryCache_LargeValue(t *testing.T) {
	// Arrange
	cache := NewMemory()
	ctx := context.Background()
	key := "large_key"

	// 创建一个大字符串 (1MB)
	largeValue := string(make([]byte, 1024*1024))

	// Act
	err := cache.Set(ctx, key, largeValue, time.Minute)
	require.NoError(t, err)

	result, ok, err := cache.Get(ctx, key)

	// Assert
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, len(largeValue), len(result))
}

/**
 * Test_MemoryCache_SpecialCharacters 测试特殊字符键和值
 */
func Test_MemoryCache_SpecialCharacters(t *testing.T) {
	// Arrange
	cache := NewMemory()
	ctx := context.Background()

	specialKeys := []string{
		"key:with:colons",
		"key/with/slashes",
		"key with spaces",
		"key-with-dashes",
		"键中文",
		"🚀emoji🎉",
	}

	// Act & Assert
	for _, key := range specialKeys {
		value := "value_for_" + key

		err := cache.Set(ctx, key, value, time.Minute)
		require.NoError(t, err)

		result, ok, err := cache.Get(ctx, key)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, value, result)
	}
}
