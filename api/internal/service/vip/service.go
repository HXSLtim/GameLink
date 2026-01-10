package vip

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/cache"
)

// Cache key patterns and TTL
const (
	cacheKeyVipLevels       = "vip:levels:all"
	cacheKeyVipActiveLevels = "vip:levels:active"
	cacheKeyVipDefaultLevel = "vip:levels:default"
	cacheKeyVipConfigs      = "vip:configs:all"
	cacheKeyVipConfig       = "vip:config:%s"
	cacheTTLVipLevels       = 1 * time.Hour
	cacheTTLVipConfigs      = 1 * time.Hour
)

// Service VIP业务逻辑层
type Service struct {
	repo  repository.VipRepository
	cache cache.Cache
}

// NewVipService 创建VIP服务
func NewVipService(repo repository.VipRepository) *Service {
	return &Service{repo: repo}
}

// NewVipServiceWithCache 创建带缓存的VIP服务
func NewVipServiceWithCache(repo repository.VipRepository, c cache.Cache) *Service {
	return &Service{repo: repo, cache: c}
}

// SetCache 设置缓存实例
func (s *Service) SetCache(c cache.Cache) {
	s.cache = c
}

// invalidateLevelCaches 清除VIP等级相关缓存
func (s *Service) invalidateLevelCaches(ctx context.Context) {
	if s.cache == nil {
		return
	}
	_ = s.cache.Delete(ctx, cacheKeyVipLevels)
	_ = s.cache.Delete(ctx, cacheKeyVipActiveLevels)
	_ = s.cache.Delete(ctx, cacheKeyVipDefaultLevel)
}

// invalidateConfigCaches 清除VIP配置相关缓存
func (s *Service) invalidateConfigCaches(ctx context.Context, key string) {
	if s.cache == nil {
		return
	}
	_ = s.cache.Delete(ctx, cacheKeyVipConfigs)
	if key != "" {
		_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyVipConfig, key))
	}
}

// ============================================================================
// VIP等级管理
// ============================================================================

// CreateLevel 创建VIP等级
func (s *Service) CreateLevel(ctx context.Context, level *model.VipLevel) error {
	// 检查Slug是否已存在
	existing, err := s.repo.GetLevelBySlug(ctx, level.Slug)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("check slug: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("slug %s already exists", level.Slug)
	}

	if err := s.repo.CreateLevel(ctx, level); err != nil {
		return fmt.Errorf("create level: %w", err)
	}

	// 清除缓存
	s.invalidateLevelCaches(ctx)

	return nil
}

// GetLevel 获取VIP等级
func (s *Service) GetLevel(ctx context.Context, id uint64) (*model.VipLevel, error) {
	level, err := s.repo.GetLevel(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get level: %w", err)
	}
	return level, nil
}

// GetLevelBySlug 根据Slug获取VIP等级
func (s *Service) GetLevelBySlug(ctx context.Context, slug string) (*model.VipLevel, error) {
	level, err := s.repo.GetLevelBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("get level by slug: %w", err)
	}
	return level, nil
}

// GetDefaultLevel 获取默认VIP等级
func (s *Service) GetDefaultLevel(ctx context.Context) (*model.VipLevel, error) {
	level, err := s.repo.GetDefaultLevel(ctx)
	if err != nil {
		return nil, fmt.Errorf("get default level: %w", err)
	}
	return level, nil
}

// ListLevels 获取所有VIP等级
func (s *Service) ListLevels(ctx context.Context) ([]model.VipLevel, error) {
	// 尝试从缓存获取
	if s.cache != nil {
		if val, ok, _ := s.cache.Get(ctx, cacheKeyVipLevels); ok {
			var levels []model.VipLevel
			if err := json.Unmarshal([]byte(val), &levels); err == nil {
				return levels, nil
			}
		}
	}

	levels, err := s.repo.ListLevels(ctx)
	if err != nil {
		return nil, fmt.Errorf("list levels: %w", err)
	}

	// 写入缓存
	if s.cache != nil {
		if data, err := json.Marshal(levels); err == nil {
			_ = s.cache.Set(ctx, cacheKeyVipLevels, string(data), cacheTTLVipLevels)
		}
	}

	return levels, nil
}

// ListActiveLevels 获取所有启用的VIP等级
func (s *Service) ListActiveLevels(ctx context.Context) ([]model.VipLevel, error) {
	// 尝试从缓存获取
	if s.cache != nil {
		if val, ok, _ := s.cache.Get(ctx, cacheKeyVipActiveLevels); ok {
			var levels []model.VipLevel
			if err := json.Unmarshal([]byte(val), &levels); err == nil {
				return levels, nil
			}
		}
	}

	levels, err := s.repo.ListActiveLevels(ctx)
	if err != nil {
		return nil, fmt.Errorf("list active levels: %w", err)
	}

	// 写入缓存
	if s.cache != nil {
		if data, err := json.Marshal(levels); err == nil {
			_ = s.cache.Set(ctx, cacheKeyVipActiveLevels, string(data), cacheTTLVipLevels)
		}
	}

	return levels, nil
}

// ListLevelsPaged 分页获取VIP等级
func (s *Service) ListLevelsPaged(ctx context.Context, opts repository.VipLevelListOptions) ([]model.VipLevel, int64, error) {
	levels, total, err := s.repo.ListLevelsPaged(ctx, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("list levels paged: %w", err)
	}
	return levels, total, nil
}

// UpdateLevel 更新VIP等级
func (s *Service) UpdateLevel(ctx context.Context, level *model.VipLevel) error {
	// 检查是否存在
	existing, err := s.repo.GetLevel(ctx, level.ID)
	if err != nil {
		return fmt.Errorf("get level: %w", err)
	}

	// 如果Slug变更，检查新Slug是否已存在
	if level.Slug != existing.Slug {
		other, err := s.repo.GetLevelBySlug(ctx, level.Slug)
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("check slug: %w", err)
		}
		if other != nil {
			return fmt.Errorf("slug %s already exists", level.Slug)
		}
	}

	if err := s.repo.UpdateLevel(ctx, level); err != nil {
		return fmt.Errorf("update level: %w", err)
	}

	// 清除缓存
	s.invalidateLevelCaches(ctx)

	return nil
}

// DeleteLevel 删除VIP等级
func (s *Service) DeleteLevel(ctx context.Context, id uint64) error {
	if err := s.repo.DeleteLevel(ctx, id); err != nil {
		return fmt.Errorf("delete level: %w", err)
	}

	// 清除缓存
	s.invalidateLevelCaches(ctx)

	return nil
}

// SetDefaultLevel 设置默认VIP等级
func (s *Service) SetDefaultLevel(ctx context.Context, id uint64) error {
	if err := s.repo.SetDefaultLevel(ctx, id); err != nil {
		return fmt.Errorf("set default level: %w", err)
	}

	// 清除缓存
	s.invalidateLevelCaches(ctx)

	return nil
}

// GetLevelByExp 根据经验值获取对应的VIP等级
func (s *Service) GetLevelByExp(ctx context.Context, exp int64) (*model.VipLevel, error) {
	level, err := s.repo.GetLevelByExp(ctx, exp)
	if err != nil {
		return nil, fmt.Errorf("get level by exp: %w", err)
	}
	return level, nil
}

// BatchUpdateLevelStatus 批量更新VIP等级状态
func (s *Service) BatchUpdateLevelStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error) {
	affected, err := s.repo.BatchUpdateLevelStatus(ctx, ids, isActive)
	if err != nil {
		return 0, fmt.Errorf("batch update level status: %w", err)
	}

	// 清除缓存
	s.invalidateLevelCaches(ctx)

	return affected, nil
}

// BatchDeleteLevels 批量删除VIP等级
func (s *Service) BatchDeleteLevels(ctx context.Context, ids []uint64) (int64, error) {
	affected, err := s.repo.BatchDeleteLevels(ctx, ids)
	if err != nil {
		return 0, fmt.Errorf("batch delete levels: %w", err)
	}

	// 清除缓存
	s.invalidateLevelCaches(ctx)

	return affected, nil
}

// ============================================================================
// VIP配置管理
// ============================================================================

// GetConfig 获取VIP配置
func (s *Service) GetConfig(ctx context.Context, key string) (*model.VipConfig, error) {
	config, err := s.repo.GetConfig(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("get config: %w", err)
	}
	return config, nil
}

// GetConfigValue 获取VIP配置值
func (s *Service) GetConfigValue(ctx context.Context, key string) (string, error) {
	config, err := s.repo.GetConfig(ctx, key)
	if err != nil {
		return "", fmt.Errorf("get config: %w", err)
	}
	return config.ConfigValue, nil
}

// GetConfigInt64 获取VIP配置值（int64）
func (s *Service) GetConfigInt64(ctx context.Context, key string) (int64, error) {
	value, err := s.GetConfigValue(ctx, key)
	if err != nil {
		return 0, err
	}
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse config value: %w", err)
	}
	return result, nil
}

// ListConfigs 获取所有VIP配置
func (s *Service) ListConfigs(ctx context.Context) ([]model.VipConfig, error) {
	// 尝试从缓存获取
	if s.cache != nil {
		if val, ok, _ := s.cache.Get(ctx, cacheKeyVipConfigs); ok {
			var configs []model.VipConfig
			if err := json.Unmarshal([]byte(val), &configs); err == nil {
				return configs, nil
			}
		}
	}

	configs, err := s.repo.ListConfigs(ctx)
	if err != nil {
		return nil, fmt.Errorf("list configs: %w", err)
	}

	// 写入缓存
	if s.cache != nil {
		if data, err := json.Marshal(configs); err == nil {
			_ = s.cache.Set(ctx, cacheKeyVipConfigs, string(data), cacheTTLVipConfigs)
		}
	}

	return configs, nil
}

// SaveConfig 保存VIP配置
func (s *Service) SaveConfig(ctx context.Context, config *model.VipConfig) error {
	if err := s.repo.SaveConfig(ctx, config); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	// 清除缓存
	s.invalidateConfigCaches(ctx, config.ConfigKey)

	return nil
}

// SetConfig 设置VIP配置（简化方法）
func (s *Service) SetConfig(ctx context.Context, key, value, description string) error {
	// 尝试获取现有配置
	config, err := s.repo.GetConfig(ctx, key)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("get config: %w", err)
	}

	if config == nil {
		config = &model.VipConfig{
			ConfigKey:   key,
			ConfigValue: value,
			Description: description,
		}
	} else {
		config.ConfigValue = value
		if description != "" {
			config.Description = description
		}
	}

	return s.SaveConfig(ctx, config)
}

// DeleteConfig 删除VIP配置
func (s *Service) DeleteConfig(ctx context.Context, key string) error {
	if err := s.repo.DeleteConfig(ctx, key); err != nil {
		return fmt.Errorf("delete config: %w", err)
	}

	// 清除缓存
	s.invalidateConfigCaches(ctx, key)

	return nil
}

// ============================================================================
// VIP业务逻辑
// ============================================================================

// GetUnlockThreshold 获取VIP解锁门槛
func (s *Service) GetUnlockThreshold(ctx context.Context) (consumeThreshold, rechargeThreshold int64, err error) {
	consumeThreshold, _ = s.GetConfigInt64(ctx, model.VipConfigUnlockByConsume)
	rechargeThreshold, _ = s.GetConfigInt64(ctx, model.VipConfigUnlockByRecharge)
	return consumeThreshold, rechargeThreshold, nil
}

// CheckVipUnlock 检查用户是否满足VIP解锁条件
func (s *Service) CheckVipUnlock(ctx context.Context, totalConsume, totalRecharge int64) (bool, error) {
	consumeThreshold, rechargeThreshold, err := s.GetUnlockThreshold(ctx)
	if err != nil {
		return false, err
	}

	// 任一条件满足即可解锁
	if consumeThreshold > 0 && totalConsume >= consumeThreshold {
		return true, nil
	}
	if rechargeThreshold > 0 && totalRecharge >= rechargeThreshold {
		return true, nil
	}
	return false, nil
}

// CalculateVipLevel 根据经验值计算VIP等级
func (s *Service) CalculateVipLevel(ctx context.Context, exp int64) (*model.VipLevel, error) {
	return s.GetLevelByExp(ctx, exp)
}

// GetVipExpireDays 获取VIP过期天数（0=永久）
func (s *Service) GetVipExpireDays(ctx context.Context) (int64, error) {
	days, err := s.GetConfigInt64(ctx, model.VipConfigExpireDays)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return 0, nil // 默认永久
		}
		return 0, err
	}
	return days, nil
}
