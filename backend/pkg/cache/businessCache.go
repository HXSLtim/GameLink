package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// BusinessCache 业务数据缓存服务
// 封装常用业务数据的缓存操作
type BusinessCache struct {
	hotspot *HotspotCache
}

// 缓存 key 前缀
const (
	PrefixUser       = "user"
	PrefixPlayer     = "player"
	PrefixOrder      = "order"
	PrefixGame       = "game"
	PrefixService    = "service"
	PrefixConfig     = "config"
	PrefixStats      = "stats"
	PrefixMenu       = "menu"
	PrefixPermission = "perm"
	PrefixRole       = "role"
)

// 缓存 TTL 配置
const (
	TTLUserProfile   = 5 * time.Minute  // 用户资料
	TTLPlayerProfile = 5 * time.Minute  // 陪玩师资料
	TTLPlayerList    = 2 * time.Minute  // 陪玩师列表
	TTLGameList      = 10 * time.Minute // 游戏列表
	TTLServiceList   = 10 * time.Minute // 服务项列表
	TTLConfig        = 30 * time.Minute // 系统配置
	TTLStats         = 1 * time.Minute  // 统计数据
	TTLMenu          = 15 * time.Minute // 菜单数据
	TTLPermission    = 10 * time.Minute // 权限数据
	TTLRole          = 10 * time.Minute // 角色数据
	TTLOrderDetail   = 3 * time.Minute  // 订单详情
	TTLHotPlayer     = 30 * time.Second // 热门陪玩师（本地缓存）
)

// NewBusinessCache 创建业务缓存服务
func NewBusinessCache(remote Cache) *BusinessCache {
	config := DefaultHotspotConfig()
	config.LocalTTL = 30 * time.Second
	config.RemoteTTL = 5 * time.Minute
	config.HotThreshold = 50 // 50 次访问视为热点

	return &BusinessCache{
		hotspot: NewHotspotCache(remote, config),
	}
}

// NewBusinessCacheWithConfig 使用自定义配置创建业务缓存
func NewBusinessCacheWithConfig(remote Cache, config HotspotConfig) *BusinessCache {
	return &BusinessCache{
		hotspot: NewHotspotCache(remote, config),
	}
}

// ==================== 用户相关缓存 ====================

// GetUser 获取用户缓存
func (bc *BusinessCache) GetUser(ctx context.Context, userID uint64, dest interface{}) (bool, error) {
	key := CacheKey(PrefixUser, userID)
	return bc.hotspot.GetJSON(ctx, key, dest)
}

// SetUser 设置用户缓存
func (bc *BusinessCache) SetUser(ctx context.Context, userID uint64, user interface{}) error {
	key := CacheKey(PrefixUser, userID)
	return bc.hotspot.SetJSON(ctx, key, user, TTLUserProfile)
}

// DeleteUser 删除用户缓存
func (bc *BusinessCache) DeleteUser(ctx context.Context, userID uint64) error {
	key := CacheKey(PrefixUser, userID)
	return bc.hotspot.Delete(ctx, key)
}

// ==================== 陪玩师相关缓存 ====================

// GetPlayer 获取陪玩师缓存
func (bc *BusinessCache) GetPlayer(ctx context.Context, playerID uint64, dest interface{}) (bool, error) {
	key := CacheKey(PrefixPlayer, playerID)
	return bc.hotspot.GetJSON(ctx, key, dest)
}

// SetPlayer 设置陪玩师缓存
func (bc *BusinessCache) SetPlayer(ctx context.Context, playerID uint64, player interface{}) error {
	key := CacheKey(PrefixPlayer, playerID)
	return bc.hotspot.SetJSON(ctx, key, player, TTLPlayerProfile)
}

// DeletePlayer 删除陪玩师缓存
func (bc *BusinessCache) DeletePlayer(ctx context.Context, playerID uint64) error {
	key := CacheKey(PrefixPlayer, playerID)
	return bc.hotspot.Delete(ctx, key)
}

// GetPlayerList 获取陪玩师列表缓存
func (bc *BusinessCache) GetPlayerList(ctx context.Context, gameID uint64, page, pageSize int, dest interface{}) (bool, error) {
	key := CacheKey(PrefixPlayer, "list", gameID, page, pageSize)
	return bc.hotspot.GetJSON(ctx, key, dest)
}

// SetPlayerList 设置陪玩师列表缓存
func (bc *BusinessCache) SetPlayerList(ctx context.Context, gameID uint64, page, pageSize int, list interface{}) error {
	key := CacheKey(PrefixPlayer, "list", gameID, page, pageSize)
	return bc.hotspot.SetJSON(ctx, key, list, TTLPlayerList)
}

// GetHotPlayers 获取热门陪玩师（本地缓存优先）
func (bc *BusinessCache) GetHotPlayers(ctx context.Context, limit int, dest interface{}) (bool, error) {
	key := CacheKey(PrefixPlayer, "hot", limit)
	return bc.hotspot.GetJSON(ctx, key, dest)
}

// SetHotPlayers 设置热门陪玩师（同时写入本地和远程）
func (bc *BusinessCache) SetHotPlayers(ctx context.Context, limit int, players interface{}) error {
	key := CacheKey(PrefixPlayer, "hot", limit)
	data, err := json.Marshal(players)
	if err != nil {
		return err
	}
	// 热门数据同时写入本地和远程缓存
	return bc.hotspot.SetBoth(ctx, key, string(data), TTLHotPlayer, TTLPlayerList)
}

// ==================== 游戏相关缓存 ====================

// GetGameList 获取游戏列表缓存
func (bc *BusinessCache) GetGameList(ctx context.Context, dest interface{}) (bool, error) {
	key := CacheKey(PrefixGame, "list")
	return bc.hotspot.GetJSON(ctx, key, dest)
}

// SetGameList 设置游戏列表缓存
func (bc *BusinessCache) SetGameList(ctx context.Context, games interface{}) error {
	key := CacheKey(PrefixGame, "list")
	return bc.hotspot.SetJSON(ctx, key, games, TTLGameList)
}

// GetGame 获取单个游戏缓存
func (bc *BusinessCache) GetGame(ctx context.Context, gameID uint64, dest interface{}) (bool, error) {
	key := CacheKey(PrefixGame, gameID)
	return bc.hotspot.GetJSON(ctx, key, dest)
}

// SetGame 设置单个游戏缓存
func (bc *BusinessCache) SetGame(ctx context.Context, gameID uint64, game interface{}) error {
	key := CacheKey(PrefixGame, gameID)
	return bc.hotspot.SetJSON(ctx, key, game, TTLGameList)
}

// InvalidateGameCache 使游戏缓存失效
func (bc *BusinessCache) InvalidateGameCache(ctx context.Context) error {
	// 删除游戏列表缓存
	return bc.hotspot.Delete(ctx, CacheKey(PrefixGame, "list"))
}

// ==================== 服务项相关缓存 ====================

// GetServiceList 获取服务项列表缓存
func (bc *BusinessCache) GetServiceList(ctx context.Context, serviceType string, dest interface{}) (bool, error) {
	key := CacheKey(PrefixService, "list", serviceType)
	return bc.hotspot.GetJSON(ctx, key, dest)
}

// SetServiceList 设置服务项列表缓存
func (bc *BusinessCache) SetServiceList(ctx context.Context, serviceType string, services interface{}) error {
	key := CacheKey(PrefixService, "list", serviceType)
	return bc.hotspot.SetJSON(ctx, key, services, TTLServiceList)
}

// ==================== 订单相关缓存 ====================

// GetOrder 获取订单缓存
func (bc *BusinessCache) GetOrder(ctx context.Context, orderID uint64, dest interface{}) (bool, error) {
	key := CacheKey(PrefixOrder, orderID)
	return bc.hotspot.GetJSON(ctx, key, dest)
}

// SetOrder 设置订单缓存
func (bc *BusinessCache) SetOrder(ctx context.Context, orderID uint64, order interface{}) error {
	key := CacheKey(PrefixOrder, orderID)
	return bc.hotspot.SetJSON(ctx, key, order, TTLOrderDetail)
}

// DeleteOrder 删除订单缓存
func (bc *BusinessCache) DeleteOrder(ctx context.Context, orderID uint64) error {
	key := CacheKey(PrefixOrder, orderID)
	return bc.hotspot.Delete(ctx, key)
}

// ==================== 配置相关缓存 ====================

// GetConfig 获取系统配置缓存
func (bc *BusinessCache) GetConfig(ctx context.Context, configKey string, dest interface{}) (bool, error) {
	key := CacheKey(PrefixConfig, configKey)
	return bc.hotspot.GetJSON(ctx, key, dest)
}

// SetConfig 设置系统配置缓存
func (bc *BusinessCache) SetConfig(ctx context.Context, configKey string, config interface{}) error {
	key := CacheKey(PrefixConfig, configKey)
	return bc.hotspot.SetJSON(ctx, key, config, TTLConfig)
}

// ==================== 统计相关缓存 ====================

// GetStats 获取统计数据缓存
func (bc *BusinessCache) GetStats(ctx context.Context, statsKey string, dest interface{}) (bool, error) {
	key := CacheKey(PrefixStats, statsKey)
	return bc.hotspot.GetJSON(ctx, key, dest)
}

// SetStats 设置统计数据缓存
func (bc *BusinessCache) SetStats(ctx context.Context, statsKey string, stats interface{}) error {
	key := CacheKey(PrefixStats, statsKey)
	return bc.hotspot.SetJSON(ctx, key, stats, TTLStats)
}

// GetDashboardStats 获取仪表盘统计缓存
func (bc *BusinessCache) GetDashboardStats(ctx context.Context, dest interface{}) (bool, error) {
	return bc.GetStats(ctx, "dashboard", dest)
}

// SetDashboardStats 设置仪表盘统计缓存
func (bc *BusinessCache) SetDashboardStats(ctx context.Context, stats interface{}) error {
	return bc.SetStats(ctx, "dashboard", stats)
}

// ==================== 菜单和权限缓存 ====================

// GetUserMenus 获取用户菜单缓存
func (bc *BusinessCache) GetUserMenus(ctx context.Context, userID uint64, dest interface{}) (bool, error) {
	key := CacheKey(PrefixMenu, "user", userID)
	return bc.hotspot.GetJSON(ctx, key, dest)
}

// SetUserMenus 设置用户菜单缓存
func (bc *BusinessCache) SetUserMenus(ctx context.Context, userID uint64, menus interface{}) error {
	key := CacheKey(PrefixMenu, "user", userID)
	return bc.hotspot.SetJSON(ctx, key, menus, TTLMenu)
}

// DeleteUserMenus 删除用户菜单缓存
func (bc *BusinessCache) DeleteUserMenus(ctx context.Context, userID uint64) error {
	key := CacheKey(PrefixMenu, "user", userID)
	return bc.hotspot.Delete(ctx, key)
}

// ==================== 通用方法 ====================

// GetOrLoad 获取缓存，不存在则加载
func (bc *BusinessCache) GetOrLoad(ctx context.Context, key string, dest interface{}, loader func() (interface{}, error), ttl time.Duration) error {
	// 尝试从缓存获取
	if ok, err := bc.hotspot.GetJSON(ctx, key, dest); err == nil && ok {
		return nil
	}

	// 加载数据
	data, err := loader()
	if err != nil {
		return err
	}

	// 设置缓存
	if err := bc.hotspot.SetJSON(ctx, key, data, ttl); err != nil {
		// 缓存设置失败不影响返回
		fmt.Printf("cache set error: %v\n", err)
	}

	// 将数据复制到 dest
	jsonData, _ := json.Marshal(data)
	return json.Unmarshal(jsonData, dest)
}

// Warmup 预热缓存
func (bc *BusinessCache) Warmup(ctx context.Context, keys []string, loader func(key string) (string, error)) error {
	return bc.hotspot.Warmup(ctx, keys, loader)
}

// Stats 获取缓存统计
func (bc *BusinessCache) Stats() CacheStats {
	return bc.hotspot.Stats()
}

// HitRate 获取缓存命中率
func (bc *BusinessCache) HitRate() float64 {
	return bc.hotspot.HitRate()
}

// LocalHitRate 获取本地缓存命中率
func (bc *BusinessCache) LocalHitRate() float64 {
	return bc.hotspot.LocalHitRate()
}

// HotKeys 获取热点 key 列表
func (bc *BusinessCache) HotKeys() []string {
	return bc.hotspot.HotKeys()
}

// Close 关闭缓存
func (bc *BusinessCache) Close(ctx context.Context) error {
	return bc.hotspot.Close(ctx)
}

// Delete 删除指定 key
func (bc *BusinessCache) Delete(ctx context.Context, key string) error {
	return bc.hotspot.Delete(ctx, key)
}

// InvalidatePattern 使匹配模式的缓存失效（简化实现）
func (bc *BusinessCache) InvalidatePattern(ctx context.Context, prefix string) error {
	// 注意：这是简化实现，实际生产环境可能需要使用 Redis SCAN 命令
	// 这里只删除已知的热点 key
	hotKeys := bc.hotspot.HotKeys()
	for _, key := range hotKeys {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			_ = bc.hotspot.Delete(ctx, key)
		}
	}
	return nil
}
