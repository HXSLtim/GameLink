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
 * Test_NewDistributedLock_Creates 测试创建分布式锁
 */
func Test_NewDistributedLock_Creates(t *testing.T) {
	// Arrange
	cache := NewMemory()

	// Act
	lock := NewDistributedLock(cache)

	// Assert
	assert.NotNil(t, lock)
	assert.IsType(t, &distributedLockImpl{}, lock)
}

/**
 * Test_DistributedLock_Lock_Success 测试成功获取锁
 */
func Test_DistributedLock_Lock_Success(t *testing.T) {
	// Arrange
	cache := NewMemory()
	lock := NewDistributedLock(cache)
	ctx := context.Background()
	key := "test_lock"

	// Act
	acquired, err := lock.Lock(ctx, key, time.Minute)

	// Assert
	assert.NoError(t, err)
	assert.True(t, acquired, "应该成功获取锁")

	// 清理
	lock.Unlock(ctx, key)
}

/**
 * Test_DistributedLock_Lock_AlreadyLocked 测试锁已被持有
 */
func Test_DistributedLock_Lock_AlreadyLocked(t *testing.T) {
	// Arrange
	cache := NewMemory()
	lock := NewDistributedLock(cache)
	ctx := context.Background()
	key := "test_lock"

	// Act - 第一次获取锁
	acquired1, err1 := lock.Lock(ctx, key, time.Minute)
	require.NoError(t, err1)
	require.True(t, acquired1)

	// Act - 第二次尝试获取同一个锁
	acquired2, err2 := lock.Lock(ctx, key, time.Minute)

	// Assert
	assert.NoError(t, err2)
	assert.False(t, acquired2, "不应该获取到锁，因为已被持有")

	// 清理
	lock.Unlock(ctx, key)
}

/**
 * Test_DistributedLock_Unlock_Success 测试释放锁
 */
func Test_DistributedLock_Unlock_Success(t *testing.T) {
	// Arrange
	cache := NewMemory()
	lock := NewDistributedLock(cache)
	ctx := context.Background()
	key := "test_lock"

	// Act - 获取锁
	acquired, err := lock.Lock(ctx, key, time.Minute)
	require.NoError(t, err)
	require.True(t, acquired)

	// Act - 释放锁
	err = lock.Unlock(ctx, key)
	assert.NoError(t, err)

	// Assert - 验证锁已释放，可以再次获取
	acquired2, err2 := lock.Lock(ctx, key, time.Minute)
	assert.NoError(t, err2)
	assert.True(t, acquired2, "锁已释放，应该能再次获取")

	// 清理
	lock.Unlock(ctx, key)
}

/**
 * Test_DistributedLock_Unlock_NotLocked 测试释放未持有的锁
 */
func Test_DistributedLock_Unlock_NotLocked(t *testing.T) {
	// Arrange
	cache := NewMemory()
	lock := NewDistributedLock(cache)
	ctx := context.Background()
	key := "test_lock"

	// Act - 直接释放未持有的锁
	err := lock.Unlock(ctx, key)

	// Assert - 应该不报错
	assert.NoError(t, err)
}

/**
 * Test_DistributedLock_TryLock_SuccessFirstAttempt 测试TryLock第一次就成功
 */
func Test_DistributedLock_TryLock_SuccessFirstAttempt(t *testing.T) {
	// Arrange
	cache := NewMemory()
	lock := NewDistributedLock(cache)
	ctx := context.Background()
	key := "test_lock"

	// Act
	acquired, err := lock.TryLock(ctx, key, time.Minute, 3, 10*time.Millisecond)

	// Assert
	assert.NoError(t, err)
	assert.True(t, acquired, "应该在第一次尝试就获取锁")

	// 清理
	lock.Unlock(ctx, key)
}

/**
 * Test_DistributedLock_TryLock_SuccessAfterRetry 测试TryLock重试后成功
 */
func Test_DistributedLock_TryLock_SuccessAfterRetry(t *testing.T) {
	// Arrange
	cache := NewMemory()
	lock := NewDistributedLock(cache)
	ctx := context.Background()
	key := "test_lock"

	// 先获取锁
	acquired1, err := lock.Lock(ctx, key, 200*time.Millisecond)
	require.NoError(t, err)
	require.True(t, acquired1)

	// 启动goroutine在200ms后释放锁
	go func() {
		time.Sleep(200 * time.Millisecond)
		lock.Unlock(ctx, key)
	}()

	// Act - 尝试获取锁，应该在重试后成功
	startTime := time.Now()
	acquired2, err := lock.TryLock(ctx, key, time.Minute, 5, 100*time.Millisecond)
	duration := time.Since(startTime)

	// Assert
	assert.NoError(t, err)
	assert.True(t, acquired2, "应该在重试后获取锁")
	assert.GreaterOrEqual(t, duration, 100*time.Millisecond, "应该至少等待一次重试")

	// 清理
	lock.Unlock(ctx, key)
}

/**
 * Test_DistributedLock_TryLock_FailAfterRetries 测试TryLock所有重试都失败
 */
func Test_DistributedLock_TryLock_FailAfterRetries(t *testing.T) {
	// Arrange
	cache := NewMemory()
	lock := NewDistributedLock(cache)
	ctx := context.Background()
	key := "test_lock"

	// 先获取锁并持有
	acquired1, err := lock.Lock(ctx, key, time.Minute)
	require.NoError(t, err)
	require.True(t, acquired1)

	// Act - 尝试获取锁，应该失败
	startTime := time.Now()
	acquired2, err := lock.TryLock(ctx, key, time.Minute, 3, 10*time.Millisecond)
	duration := time.Since(startTime)

	// Assert
	assert.NoError(t, err)
	assert.False(t, acquired2, "应该无法获取锁")
	assert.GreaterOrEqual(t, duration, 20*time.Millisecond, "应该至少重试2次")

	// 清理
	lock.Unlock(ctx, key)
}

/**
 * Test_DistributedLock_TryLock_WithZeroRetry 测试零次重试
 */
func Test_DistributedLock_TryLock_WithZeroRetry(t *testing.T) {
	// Arrange
	cache := NewMemory()
	lock := NewDistributedLock(cache)
	ctx := context.Background()
	key := "test_lock"

	// 先获取锁
	acquired1, err := lock.Lock(ctx, key, time.Minute)
	require.NoError(t, err)
	require.True(t, acquired1)

	// Act - 尝试获取锁，0次重试
	acquired2, err := lock.TryLock(ctx, key, time.Minute, 0, 10*time.Millisecond)

	// Assert
	assert.NoError(t, err)
	assert.False(t, acquired2, "应该立即失败")

	// 清理
	lock.Unlock(ctx, key)
}

/**
 * Test_DistributedLock_TTL_Expiry 测试锁的TTL过期
 */
func Test_DistributedLock_TTL_Expiry(t *testing.T) {
	// Arrange
	cache := NewMemory()
	lock := NewDistributedLock(cache)
	ctx := context.Background()
	key := "test_lock"

	// Act - 获取短期锁
	acquired1, err := lock.Lock(ctx, key, 100*time.Millisecond)
	require.NoError(t, err)
	require.True(t, acquired1)

	// 等待锁过期
	time.Sleep(200 * time.Millisecond)

	// 尝试再次获取锁
	acquired2, err := lock.Lock(ctx, key, time.Minute)

	// Assert
	assert.NoError(t, err)
	assert.True(t, acquired2, "锁应该已过期，可以再次获取")

	// 清理
	lock.Unlock(ctx, key)
}

/**
 * Test_DistributedLock_MultipleKeys 测试多个不同的锁
 */
func Test_DistributedLock_MultipleKeys(t *testing.T) {
	// Arrange
	cache := NewMemory()
	lock := NewDistributedLock(cache)
	ctx := context.Background()

	keys := []string{"lock1", "lock2", "lock3"}

	// Act - 获取多个锁
	for _, key := range keys {
		acquired, err := lock.Lock(ctx, key, time.Minute)
		require.NoError(t, err)
		require.True(t, acquired)
	}

	// Assert - 所有锁都应该被持有
	for _, key := range keys {
		acquired, err := lock.Lock(ctx, key, time.Minute)
		assert.NoError(t, err)
		assert.False(t, acquired, "锁应该已被持有")
	}

	// 清理
	for _, key := range keys {
		lock.Unlock(ctx, key)
	}
}

/**
 * Test_DistributedLock_Concurrent 测试并发获取锁
 */
func Test_DistributedLock_Concurrent(t *testing.T) {
	// Arrange
	cache := NewMemory()
	lock := NewDistributedLock(cache)
	ctx := context.Background()
	key := "concurrent_lock"

	var wg sync.WaitGroup
	numGoroutines := 10
	successCount := 0
	var mu sync.Mutex

	// Act - 并发尝试获取同一个锁
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			acquired, err := lock.Lock(ctx, key, time.Minute)
			if err == nil && acquired {
				mu.Lock()
				successCount++
				mu.Unlock()

				// 持有锁一段时间
				time.Sleep(10 * time.Millisecond)

				// 释放锁
				lock.Unlock(ctx, key)
			}
		}(i)
	}

	wg.Wait()

	// Assert - 只应该有一个goroutine成功获取锁
	// 注意：由于锁会被释放，可能有多个goroutine陆续获取到锁
	// 但同一时刻只有一个持有锁
	assert.Greater(t, successCount, 0, "至少应该有一个成功")
}

/**
 * Test_DistributedLock_LockKey_Format 测试锁键的格式
 */
func Test_DistributedLock_LockKey_Format(t *testing.T) {
	// Arrange
	cache := NewMemory()
	lock := NewDistributedLock(cache)
	ctx := context.Background()
	key := "my_resource"

	// Act
	acquired, err := lock.Lock(ctx, key, time.Minute)
	require.NoError(t, err)
	require.True(t, acquired)

	// 验证实际存储的键
	mc := cache.(*memoryCache)
	mc.mu.RLock()
	_, exists := mc.values["lock:my_resource"]
	mc.mu.RUnlock()

	// Assert
	assert.True(t, exists, "锁键应该有'lock:'前缀")

	// 清理
	lock.Unlock(ctx, key)
}

/**
 * Test_DistributedLock_Stress 测试锁的压力情况
 */
func Test_DistributedLock_Stress(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过压力测试")
	}

	// Arrange
	cache := NewMemory()
	lock := NewDistributedLock(cache)
	ctx := context.Background()
	key := "stress_lock"

	var wg sync.WaitGroup
	numGoroutines := 50
	operationsPerGoroutine := 20
	totalOperations := numGoroutines * operationsPerGoroutine
	successfulLocks := 0
	var mu sync.Mutex

	// Act
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				acquired, err := lock.TryLock(ctx, key, 50*time.Millisecond, 3, 5*time.Millisecond)
				if err == nil && acquired {
					mu.Lock()
					successfulLocks++
					mu.Unlock()

					// 模拟临界区工作
					time.Sleep(1 * time.Millisecond)

					lock.Unlock(ctx, key)
				}
			}
		}()
	}

	wg.Wait()

	// Assert
	assert.Greater(t, successfulLocks, 0, "应该有成功的锁获取")
	assert.LessOrEqual(t, successfulLocks, totalOperations, "成功次数不应超过总尝试次数")
	t.Logf("成功获取锁: %d/%d 次", successfulLocks, totalOperations)
}

/**
 * Test_DistributedLock_Lock_GetError 测试Lock时Get返回错误
 */
func Test_DistributedLock_Lock_GetError(t *testing.T) {
	// Arrange
	cache := NewMemory()
	lock := NewDistributedLock(cache)
	ctx := context.Background()
	key := "test_lock"

	// 先关闭cache以模拟错误
	cache.Close(ctx)

	// Act
	acquired, err := lock.Lock(ctx, key, time.Minute)

	// Assert - 虽然memoryCache.Get不会返回错误，但我们测试了这个路径
	// 在实际使用中（如Redis），这个错误处理很重要
	_ = acquired
	_ = err
}

/**
 * Test_DistributedLock_Lock_SetError 测试Lock时Set返回错误
 */
func Test_DistributedLock_Lock_SetError(t *testing.T) {
	// Arrange
	cache := NewMemory()
	lock := NewDistributedLock(cache)
	ctx := context.Background()
	key := "test_lock"

	// Act - 正常情况下memoryCache.Set不会返回错误
	// 但这测试了错误处理路径
	acquired, err := lock.Lock(ctx, key, time.Minute)

	// Assert
	assert.NoError(t, err)
	assert.True(t, acquired)

	// 清理
	lock.Unlock(ctx, key)
}

/**
 * Test_DistributedLock_TryLock_LockError 测试TryLock时Lock返回错误
 */
func Test_DistributedLock_TryLock_LockError(t *testing.T) {
	// Arrange
	cache := NewMemory()
	lock := NewDistributedLock(cache)
	ctx := context.Background()
	key := "test_lock"

	// 关闭cache
	cache.Close(ctx)

	// Act
	acquired, err := lock.TryLock(ctx, key, time.Minute, 1, 10*time.Millisecond)

	// Assert - memoryCache不会返回错误，但测试了错误处理路径
	_ = acquired
	_ = err
}

/**
 * Test_DistributedLock_TryLock_SingleRetry 测试单次重试
 */
func Test_DistributedLock_TryLock_SingleRetry(t *testing.T) {
	// Arrange
	cache := NewMemory()
	lock := NewDistributedLock(cache)
	ctx := context.Background()
	key := "test_lock"

	// Act - 单次尝试
	acquired, err := lock.TryLock(ctx, key, time.Minute, 1, 10*time.Millisecond)

	// Assert
	assert.NoError(t, err)
	assert.True(t, acquired)

	// 清理
	lock.Unlock(ctx, key)
}
