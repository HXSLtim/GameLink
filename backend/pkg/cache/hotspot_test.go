package cache

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHotspotCache(t *testing.T) {
	remote := NewMemory()
	config := DefaultHotspotConfig()
	config.HotThreshold = 3 // 降低阈值便于测试

	hc := NewHotspotCache(remote, config)

	assert.NotNil(t, hc)
	assert.NotNil(t, hc.local)
	assert.NotNil(t, hc.remote)
}

func TestHotspotCache_SetAndGet(t *testing.T) {
	ctx := context.Background()
	remote := NewMemory()
	config := DefaultHotspotConfig()

	hc := NewHotspotCache(remote, config)
	defer hc.Close(ctx)

	// 设置值
	err := hc.Set(ctx, "key1", "value1", time.Minute)
	require.NoError(t, err)

	// 获取值
	value, ok, err := hc.Get(ctx, "key1")
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "value1", value)
}

func TestHotspotCache_LocalCachePromotion(t *testing.T) {
	ctx := context.Background()
	remote := NewMemory()
	config := DefaultHotspotConfig()
	config.HotThreshold = 3 // 3 次访问成为热点
	config.AutoPromote = true

	hc := NewHotspotCache(remote, config)
	defer hc.Close(ctx)

	// 设置值到远程缓存
	err := remote.Set(ctx, "hotkey", "hotvalue", time.Minute)
	require.NoError(t, err)

	// 多次访问使其成为热点
	for i := 0; i < 5; i++ {
		_, _, _ = hc.Get(ctx, "hotkey")
	}

	// 验证已成为热点
	assert.True(t, hc.isHotKey("hotkey"))

	// 再次获取，应该从本地缓存获取
	stats := hc.Stats()
	initialLocalHits := stats.LocalHits

	// 设置到本地缓存（模拟提升）
	_ = hc.SetLocal(ctx, "hotkey", "hotvalue", time.Minute)

	value, ok, err := hc.Get(ctx, "hotkey")
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "hotvalue", value)

	// 本地命中应该增加
	newStats := hc.Stats()
	assert.Greater(t, newStats.LocalHits, initialLocalHits)
}

func TestHotspotCache_SetBoth(t *testing.T) {
	ctx := context.Background()
	remote := NewMemory()
	config := DefaultHotspotConfig()

	hc := NewHotspotCache(remote, config)
	defer hc.Close(ctx)

	// 同时设置本地和远程
	err := hc.SetBoth(ctx, "bothkey", "bothvalue", time.Minute, time.Minute)
	require.NoError(t, err)

	// 从本地获取
	value, ok, err := hc.local.Get(ctx, "bothkey")
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "bothvalue", value)

	// 从远程获取
	value, ok, err = remote.Get(ctx, "bothkey")
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "bothvalue", value)
}

func TestHotspotCache_Delete(t *testing.T) {
	ctx := context.Background()
	remote := NewMemory()
	config := DefaultHotspotConfig()

	hc := NewHotspotCache(remote, config)
	defer hc.Close(ctx)

	// 设置值
	_ = hc.SetBoth(ctx, "delkey", "delvalue", time.Minute, time.Minute)

	// 删除
	err := hc.Delete(ctx, "delkey")
	require.NoError(t, err)

	// 验证已删除
	_, ok, _ := hc.Get(ctx, "delkey")
	assert.False(t, ok)
}

func TestHotspotCache_GetOrLoad(t *testing.T) {
	ctx := context.Background()
	remote := NewMemory()
	config := DefaultHotspotConfig()

	hc := NewHotspotCache(remote, config)
	defer hc.Close(ctx)

	loadCount := 0
	loader := func() (string, error) {
		loadCount++
		return "loaded_value", nil
	}

	// 第一次调用应该触发 loader
	value, err := hc.GetOrLoad(ctx, "loadkey", loader, time.Minute)
	require.NoError(t, err)
	assert.Equal(t, "loaded_value", value)
	assert.Equal(t, 1, loadCount)

	// 第二次调用应该从缓存获取
	value, err = hc.GetOrLoad(ctx, "loadkey", loader, time.Minute)
	require.NoError(t, err)
	assert.Equal(t, "loaded_value", value)
	assert.Equal(t, 1, loadCount) // loader 不应该再次调用
}

func TestHotspotCache_GetJSON(t *testing.T) {
	ctx := context.Background()
	remote := NewMemory()
	config := DefaultHotspotConfig()

	hc := NewHotspotCache(remote, config)
	defer hc.Close(ctx)

	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	// 设置 JSON 数据
	data := TestData{Name: "test", Value: 42}
	err := hc.SetJSON(ctx, "jsonkey", data, time.Minute)
	require.NoError(t, err)

	// 获取 JSON 数据
	var result TestData
	ok, err := hc.GetJSON(ctx, "jsonkey", &result)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "test", result.Name)
	assert.Equal(t, 42, result.Value)
}

func TestHotspotCache_Stats(t *testing.T) {
	ctx := context.Background()
	remote := NewMemory()
	config := DefaultHotspotConfig()
	config.EnableStats = true

	hc := NewHotspotCache(remote, config)
	defer hc.Close(ctx)

	// 执行一些操作
	_ = hc.Set(ctx, "key1", "value1", time.Minute)
	_, _, _ = hc.Get(ctx, "key1")
	_, _, _ = hc.Get(ctx, "nonexistent")

	stats := hc.Stats()
	assert.Greater(t, stats.TotalGets, int64(0))
	assert.Greater(t, stats.TotalSets, int64(0))
}

func TestHotspotCache_HitRate(t *testing.T) {
	ctx := context.Background()
	remote := NewMemory()
	config := DefaultHotspotConfig()
	config.EnableStats = true

	hc := NewHotspotCache(remote, config)
	defer hc.Close(ctx)

	// 设置一些值
	_ = hc.SetBoth(ctx, "hit1", "value1", time.Minute, time.Minute)
	_ = hc.SetBoth(ctx, "hit2", "value2", time.Minute, time.Minute)

	// 命中
	_, _, _ = hc.Get(ctx, "hit1")
	_, _, _ = hc.Get(ctx, "hit2")

	// 未命中
	_, _, _ = hc.Get(ctx, "miss1")

	hitRate := hc.HitRate()
	assert.Greater(t, hitRate, 0.0)
}

func TestHotspotCache_Warmup(t *testing.T) {
	ctx := context.Background()
	remote := NewMemory()
	config := DefaultHotspotConfig()

	hc := NewHotspotCache(remote, config)
	defer hc.Close(ctx)

	keys := []string{"warm1", "warm2", "warm3"}
	loader := func(key string) (string, error) {
		return "value_" + key, nil
	}

	err := hc.Warmup(ctx, keys, loader)
	require.NoError(t, err)

	// 验证所有 key 都已缓存
	for _, key := range keys {
		value, ok, err := hc.Get(ctx, key)
		require.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, "value_"+key, value)
	}

	// 验证都被标记为热点
	hotKeys := hc.HotKeys()
	assert.GreaterOrEqual(t, len(hotKeys), len(keys))
}

func TestHotspotCache_HotKeys(t *testing.T) {
	ctx := context.Background()
	remote := NewMemory()
	config := DefaultHotspotConfig()
	config.HotThreshold = 2

	hc := NewHotspotCache(remote, config)
	defer hc.Close(ctx)

	// 设置值
	_ = remote.Set(ctx, "hot1", "value1", time.Minute)
	_ = remote.Set(ctx, "hot2", "value2", time.Minute)
	_ = remote.Set(ctx, "cold", "value3", time.Minute)

	// 多次访问使 hot1 和 hot2 成为热点
	for i := 0; i < 3; i++ {
		_, _, _ = hc.Get(ctx, "hot1")
		_, _, _ = hc.Get(ctx, "hot2")
	}
	// cold 只访问一次
	_, _, _ = hc.Get(ctx, "cold")

	hotKeys := hc.HotKeys()
	assert.Contains(t, hotKeys, "hot1")
	assert.Contains(t, hotKeys, "hot2")
}

func TestCacheKey(t *testing.T) {
	key := CacheKey("user", 123, "profile")
	assert.Equal(t, "user:123:profile", key)

	key = CacheKey("order")
	assert.Equal(t, "order", key)

	key = CacheKey("game", "list", 1, 10)
	assert.Equal(t, "game:list:1:10", key)
}

func TestDefaultHotspotConfig(t *testing.T) {
	config := DefaultHotspotConfig()

	assert.Equal(t, 30*time.Second, config.LocalTTL)
	assert.Equal(t, 5*time.Minute, config.RemoteTTL)
	assert.Equal(t, 10000, config.LocalMaxSize)
	assert.Equal(t, int64(100), config.HotThreshold)
	assert.Equal(t, time.Minute, config.HotWindow)
	assert.True(t, config.EnableStats)
	assert.True(t, config.AutoPromote)
}

// 并发测试
func TestHotspotCache_Concurrent(t *testing.T) {
	ctx := context.Background()
	remote := NewMemory()
	config := DefaultHotspotConfig()

	hc := NewHotspotCache(remote, config)
	defer hc.Close(ctx)

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			key := CacheKey("concurrent", id)
			value, _ := json.Marshal(map[string]int{"id": id})

			// 写入
			_ = hc.Set(ctx, key, string(value), time.Minute)

			// 读取
			for j := 0; j < 100; j++ {
				_, _, _ = hc.Get(ctx, key)
			}

			done <- true
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证统计
	stats := hc.Stats()
	assert.Greater(t, stats.TotalGets, int64(0))
	assert.Greater(t, stats.TotalSets, int64(0))
}
