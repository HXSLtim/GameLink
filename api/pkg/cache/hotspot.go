package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// HotspotCache 热点数据缓存
// 支持多级缓存：本地内存 -> Redis -> 数据库
type HotspotCache struct {
	local   *memoryCache // 本地缓存（L1）
	remote  Cache        // 远程缓存（L2，如 Redis）
	stats   *CacheStats  // 缓存统计
	config  HotspotConfig
	mu      sync.RWMutex
	hotKeys map[string]*hotKeyInfo // 热点 key 追踪
}

// HotspotConfig 热点缓存配置
type HotspotConfig struct {
	LocalTTL     time.Duration // 本地缓存 TTL
	RemoteTTL    time.Duration // 远程缓存 TTL
	LocalMaxSize int           // 本地缓存最大条目数
	HotThreshold int64         // 热点阈值（访问次数）
	HotWindow    time.Duration // 热点统计窗口
	EnableStats  bool          // 是否启用统计
	AutoPromote  bool          // 自动提升热点到本地缓存
}

// CacheStats 缓存统计
type CacheStats struct {
	LocalHits    int64 // 本地缓存命中
	LocalMisses  int64 // 本地缓存未命中
	RemoteHits   int64 // 远程缓存命中
	RemoteMisses int64 // 远程缓存未命中
	TotalGets    int64 // 总获取次数
	TotalSets    int64 // 总设置次数
	TotalDeletes int64 // 总删除次数
}

// hotKeyInfo 热点 key 信息
type hotKeyInfo struct {
	accessCount int64
	lastAccess  time.Time
	isHot       bool
}

// DefaultHotspotConfig 默认热点缓存配置
func DefaultHotspotConfig() HotspotConfig {
	return HotspotConfig{
		LocalTTL:     30 * time.Second, // 本地缓存 30 秒
		RemoteTTL:    5 * time.Minute,  // 远程缓存 5 分钟
		LocalMaxSize: 10000,            // 本地最多 10000 条
		HotThreshold: 100,              // 100 次访问视为热点
		HotWindow:    time.Minute,      // 1 分钟统计窗口
		EnableStats:  true,
		AutoPromote:  true,
	}
}

// NewHotspotCache 创建热点缓存
func NewHotspotCache(remote Cache, config HotspotConfig) *HotspotCache {
	hc := &HotspotCache{
		local:   NewMemory().(*memoryCache),
		remote:  remote,
		config:  config,
		stats:   &CacheStats{},
		hotKeys: make(map[string]*hotKeyInfo),
	}

	// 启动热点清理协程
	go hc.cleanupHotKeys()

	return hc
}

// Get 获取缓存值（多级查找）
func (hc *HotspotCache) Get(ctx context.Context, key string) (string, bool, error) {
	if hc.config.EnableStats {
		atomic.AddInt64(&hc.stats.TotalGets, 1)
	}

	// 记录访问
	hc.recordAccess(key)

	// L1: 本地缓存
	if value, ok, err := hc.local.Get(ctx, key); err == nil && ok {
		if hc.config.EnableStats {
			atomic.AddInt64(&hc.stats.LocalHits, 1)
		}
		return value, true, nil
	}

	if hc.config.EnableStats {
		atomic.AddInt64(&hc.stats.LocalMisses, 1)
	}

	// L2: 远程缓存
	if hc.remote != nil {
		if value, ok, err := hc.remote.Get(ctx, key); err == nil && ok {
			if hc.config.EnableStats {
				atomic.AddInt64(&hc.stats.RemoteHits, 1)
			}

			// 如果是热点 key，提升到本地缓存
			if hc.config.AutoPromote && hc.isHotKey(key) {
				_ = hc.local.Set(ctx, key, value, hc.config.LocalTTL)
			}

			return value, true, nil
		}

		if hc.config.EnableStats {
			atomic.AddInt64(&hc.stats.RemoteMisses, 1)
		}
	}

	return "", false, nil
}

// Set 设置缓存值
func (hc *HotspotCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	if hc.config.EnableStats {
		atomic.AddInt64(&hc.stats.TotalSets, 1)
	}

	// 设置远程缓存
	if hc.remote != nil {
		remoteTTL := ttl
		if remoteTTL == 0 {
			remoteTTL = hc.config.RemoteTTL
		}
		if err := hc.remote.Set(ctx, key, value, remoteTTL); err != nil {
			return err
		}
	}

	// 如果是热点 key 或显式设置，也写入本地缓存
	if hc.isHotKey(key) {
		localTTL := hc.config.LocalTTL
		if ttl > 0 && ttl < localTTL {
			localTTL = ttl
		}
		_ = hc.local.Set(ctx, key, value, localTTL)
	}

	return nil
}

// SetLocal 仅设置本地缓存（用于热点数据预热）
func (hc *HotspotCache) SetLocal(ctx context.Context, key, value string, ttl time.Duration) error {
	if ttl == 0 {
		ttl = hc.config.LocalTTL
	}
	return hc.local.Set(ctx, key, value, ttl)
}

// SetBoth 同时设置本地和远程缓存
func (hc *HotspotCache) SetBoth(ctx context.Context, key, value string, localTTL, remoteTTL time.Duration) error {
	if hc.config.EnableStats {
		atomic.AddInt64(&hc.stats.TotalSets, 1)
	}

	// 设置本地缓存
	if localTTL == 0 {
		localTTL = hc.config.LocalTTL
	}
	if err := hc.local.Set(ctx, key, value, localTTL); err != nil {
		return err
	}

	// 设置远程缓存
	if hc.remote != nil {
		if remoteTTL == 0 {
			remoteTTL = hc.config.RemoteTTL
		}
		return hc.remote.Set(ctx, key, value, remoteTTL)
	}

	return nil
}

// Delete 删除缓存
func (hc *HotspotCache) Delete(ctx context.Context, key string) error {
	if hc.config.EnableStats {
		atomic.AddInt64(&hc.stats.TotalDeletes, 1)
	}

	// 删除本地缓存
	_ = hc.local.Delete(ctx, key)

	// 删除远程缓存
	if hc.remote != nil {
		return hc.remote.Delete(ctx, key)
	}

	return nil
}

// Close 关闭缓存
func (hc *HotspotCache) Close(ctx context.Context) error {
	_ = hc.local.Close(ctx)
	if hc.remote != nil {
		return hc.remote.Close(ctx)
	}
	return nil
}

// GetOrLoad 获取缓存，如果不存在则加载
func (hc *HotspotCache) GetOrLoad(ctx context.Context, key string, loader func() (string, error), ttl time.Duration) (string, error) {
	// 先尝试获取
	if value, ok, err := hc.Get(ctx, key); err == nil && ok {
		return value, nil
	}

	// 加载数据
	value, err := loader()
	if err != nil {
		return "", err
	}

	// 设置缓存
	if err := hc.Set(ctx, key, value, ttl); err != nil {
		// 设置失败不影响返回
		return value, nil
	}

	return value, nil
}

// GetJSON 获取 JSON 缓存并反序列化
func (hc *HotspotCache) GetJSON(ctx context.Context, key string, dest interface{}) (bool, error) {
	value, ok, err := hc.Get(ctx, key)
	if err != nil || !ok {
		return false, err
	}

	if err := json.Unmarshal([]byte(value), dest); err != nil {
		return false, err
	}

	return true, nil
}

// SetJSON 序列化并设置 JSON 缓存
func (hc *HotspotCache) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return hc.Set(ctx, key, string(data), ttl)
}

// Warmup 预热缓存
func (hc *HotspotCache) Warmup(ctx context.Context, keys []string, loader func(key string) (string, error)) error {
	for _, key := range keys {
		value, err := loader(key)
		if err != nil {
			continue // 跳过加载失败的 key
		}

		// 预热数据直接写入本地和远程缓存
		_ = hc.SetBoth(ctx, key, value, hc.config.LocalTTL, hc.config.RemoteTTL)

		// 标记为热点
		hc.markHot(key)
	}
	return nil
}

// Stats 获取缓存统计
func (hc *HotspotCache) Stats() CacheStats {
	return CacheStats{
		LocalHits:    atomic.LoadInt64(&hc.stats.LocalHits),
		LocalMisses:  atomic.LoadInt64(&hc.stats.LocalMisses),
		RemoteHits:   atomic.LoadInt64(&hc.stats.RemoteHits),
		RemoteMisses: atomic.LoadInt64(&hc.stats.RemoteMisses),
		TotalGets:    atomic.LoadInt64(&hc.stats.TotalGets),
		TotalSets:    atomic.LoadInt64(&hc.stats.TotalSets),
		TotalDeletes: atomic.LoadInt64(&hc.stats.TotalDeletes),
	}
}

// HitRate 计算缓存命中率
func (hc *HotspotCache) HitRate() float64 {
	stats := hc.Stats()
	total := stats.LocalHits + stats.LocalMisses
	if total == 0 {
		return 0
	}
	return float64(stats.LocalHits+stats.RemoteHits) / float64(total)
}

// LocalHitRate 计算本地缓存命中率
func (hc *HotspotCache) LocalHitRate() float64 {
	stats := hc.Stats()
	total := stats.LocalHits + stats.LocalMisses
	if total == 0 {
		return 0
	}
	return float64(stats.LocalHits) / float64(total)
}

// recordAccess 记录访问
func (hc *HotspotCache) recordAccess(key string) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	info, ok := hc.hotKeys[key]
	if !ok {
		info = &hotKeyInfo{}
		hc.hotKeys[key] = info
	}

	info.accessCount++
	info.lastAccess = time.Now()

	// 检查是否达到热点阈值
	if !info.isHot && info.accessCount >= hc.config.HotThreshold {
		info.isHot = true
	}
}

// isHotKey 检查是否为热点 key
func (hc *HotspotCache) isHotKey(key string) bool {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	if info, ok := hc.hotKeys[key]; ok {
		return info.isHot
	}
	return false
}

// markHot 标记为热点
func (hc *HotspotCache) markHot(key string) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	info, ok := hc.hotKeys[key]
	if !ok {
		info = &hotKeyInfo{}
		hc.hotKeys[key] = info
	}
	info.isHot = true
	info.lastAccess = time.Now()
}

// cleanupHotKeys 清理过期的热点 key 信息
func (hc *HotspotCache) cleanupHotKeys() {
	ticker := time.NewTicker(hc.config.HotWindow)
	defer ticker.Stop()

	for range ticker.C {
		hc.mu.Lock()
		now := time.Now()
		for key, info := range hc.hotKeys {
			// 超过统计窗口未访问，重置计数
			if now.Sub(info.lastAccess) > hc.config.HotWindow {
				info.accessCount = 0
				info.isHot = false
			}
			// 超过 2 倍窗口未访问，删除记录
			if now.Sub(info.lastAccess) > 2*hc.config.HotWindow {
				delete(hc.hotKeys, key)
			}
		}
		hc.mu.Unlock()
	}
}

// HotKeys 获取当前热点 key 列表
func (hc *HotspotCache) HotKeys() []string {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	var keys []string
	for key, info := range hc.hotKeys {
		if info.isHot {
			keys = append(keys, key)
		}
	}
	return keys
}

// CacheKey 生成缓存 key 的辅助函数
func CacheKey(prefix string, parts ...interface{}) string {
	key := prefix
	for _, part := range parts {
		key += fmt.Sprintf(":%v", part)
	}
	return key
}
