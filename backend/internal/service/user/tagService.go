/**
 * @file user tag service
 * @description 用户标签业务逻辑层
 */

package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/cache"
)

// UserTagService 用户标签服务
type UserTagService struct {
	tagRepo  repository.UserTagRepository
	userRepo repository.UserRepository
	cache    cache.Cache
}

// NewUserTagService 创建用户标签服务
func NewUserTagService(
	tagRepo repository.UserTagRepository,
	userRepo repository.UserRepository,
	cache cache.Cache,
) *UserTagService {
	return &UserTagService{
		tagRepo:  tagRepo,
		userRepo: userRepo,
		cache:    cache,
	}
}

// CreateTag 创建标签
// Parameters:
//   - name: 标签名称
//   - color: 颜色代码（如：#FF6B6B）
//   - description: 标签描述
//
// Returns:
//   - 创建的标签
//   - 错误信息
func (s *UserTagService) CreateTag(ctx context.Context, name, color, description string) (*model.UserTag, error) {
	// 参数验证
	if len(name) == 0 || len(name) > 64 {
		return nil, errors.New("标签名称长度必须在1-64个字符之间")
	}

	// 验证颜色格式
	if !isValidColor(color) {
		return nil, errors.New("颜色格式不正确，必须是#开头的7位十六进制颜色代码")
	}

	// 检查标签是否已存在
	existingTags, err := s.tagRepo.ListTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取标签列表失败: %w", err)
	}

	for _, tag := range existingTags {
		if tag.Name == name {
			return nil, errors.New("标签名称已存在")
		}
	}

	// 创建标签
	tag := &model.UserTag{
		Name:        name,
		Color:       color,
		Description: description,
	}

	if err := s.tagRepo.CreateTag(ctx, tag); err != nil {
		return nil, fmt.Errorf("创建标签失败: %w", err)
	}

	// 清除缓存
	if s.cache != nil {
		s.cache.Delete(ctx, "user_tags:list")
	}

	return tag, nil
}

// GetTag 获取标签详情
func (s *UserTagService) GetTag(ctx context.Context, id uint64) (*model.UserTag, error) {
	return s.tagRepo.GetTag(ctx, id)
}

// ListTags 获取所有标签列表
// @return 标签列表（按创建时间倒序）
func (s *UserTagService) ListTags(ctx context.Context) ([]model.UserTag, error) {
	// 尝试从缓存获取
	cacheKey := "user_tags:list"
	if s.cache != nil {
		if cached, ok, err := s.cache.Get(ctx, cacheKey); err == nil && ok && cached != "" {
			var tags []model.UserTag
			if err := json.Unmarshal([]byte(cached), &tags); err == nil {
				return tags, nil
			}
		}
	}

	// 缓存未命中，从数据库获取
	tags, err := s.tagRepo.ListTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取标签列表失败: %w", err)
	}

	// 写入缓存（1小时过期）
	if s.cache != nil && len(tags) > 0 {
		if data, err := json.Marshal(tags); err == nil {
			s.cache.Set(ctx, cacheKey, string(data), 3600)
		}
	}

	return tags, nil
}

// UpdateTag 更新标签
func (s *UserTagService) UpdateTag(ctx context.Context, id uint64, name, color, description string) error {
	// 验证参数
	if len(name) > 64 {
		return errors.New("标签名称长度不能超过64个字符")
	}

	if color != "" && !isValidColor(color) {
		return errors.New("颜色格式不正确")
	}

	// 检查标签是否存在
	tag, err := s.tagRepo.GetTag(ctx, id)
	if err != nil {
		return fmt.Errorf("获取标签失败: %w", err)
	}

	// 更新字段
	if name != "" {
		tag.Name = name
	}
	if color != "" {
		tag.Color = color
	}
	tag.Description = description

	// 执行更新
	if err := s.tagRepo.UpdateTag(ctx, tag); err != nil {
		return fmt.Errorf("更新标签失败: %w", err)
	}

	// 清除缓存
	if s.cache != nil {
		s.cache.Delete(ctx, "user_tags:list")
		s.cache.Delete(ctx, fmt.Sprintf("user_tag:%d", id))
	}

	return nil
}

// DeleteTag 删除标签
func (s *UserTagService) DeleteTag(ctx context.Context, id uint64) error {
	// 检查标签是否存在
	if _, err := s.tagRepo.GetTag(ctx, id); err != nil {
		return fmt.Errorf("标签不存在: %w", err)
	}

	// 删除标签
	if err := s.tagRepo.DeleteTag(ctx, id); err != nil {
		return fmt.Errorf("删除标签失败: %w", err)
	}

	// 清除缓存
	if s.cache != nil {
		s.cache.Delete(ctx, "user_tags:list")
		s.cache.Delete(ctx, fmt.Sprintf("user_tag:%d", id))
	}

	return nil
}

// AddTagToUser 为用户添加标签
func (s *UserTagService) AddTagToUser(ctx context.Context, userID uint64, tagID uint64) error {
	// 检查用户是否存在
	if _, err := s.userRepo.Get(ctx, userID); err != nil {
		return fmt.Errorf("用户不存在: %w", err)
	}

	// 检查标签是否存在
	if _, err := s.tagRepo.GetTag(ctx, tagID); err != nil {
		return fmt.Errorf("标签不存在: %w", err)
	}

	// 检查用户是否已有此标签
	userTags, err := s.tagRepo.GetUserTags(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户标签失败: %w", err)
	}

	for _, tag := range userTags {
		if tag.ID == tagID {
			return errors.New("用户已拥有此标签")
		}
	}

	// 添加标签
	if err := s.tagRepo.AddTagToUser(ctx, userID, tagID); err != nil {
		return fmt.Errorf("添加标签失败: %w", err)
	}

	return nil
}

// RemoveTagFromUser 从用户移除标签
func (s *UserTagService) RemoveTagFromUser(ctx context.Context, userID uint64, tagID uint64) error {
	if err := s.tagRepo.RemoveTagFromUser(ctx, userID, tagID); err != nil {
		return fmt.Errorf("移除标签失败: %w", err)
	}
	return nil
}

// GetUserTags 获取用户的所有标签
func (s *UserTagService) GetUserTags(ctx context.Context, userID uint64) ([]model.UserTag, error) {
	// 检查用户是否存在
	if _, err := s.userRepo.Get(ctx, userID); err != nil {
		return nil, fmt.Errorf("用户不存在: %w", err)
	}

	tags, err := s.tagRepo.GetUserTags(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户标签失败: %w", err)
	}

	return tags, nil
}

// BatchSetUserTags 批量设置用户标签（覆盖式）
func (s *UserTagService) BatchSetUserTags(ctx context.Context, userID uint64, tagIDs []uint64) error {
	// 检查用户是否存在
	if _, err := s.userRepo.Get(ctx, userID); err != nil {
		return fmt.Errorf("用户不存在: %w", err)
	}

	// 验证所有标签是否存在
	for _, tagID := range tagIDs {
		if _, err := s.tagRepo.GetTag(ctx, tagID); err != nil {
			return fmt.Errorf("标签%d不存在: %w", tagID, err)
		}
	}

	// 使用事务批量设置
	if err := s.tagRepo.BatchSetUserTags(ctx, userID, tagIDs); err != nil {
		return fmt.Errorf("批量设置标签失败: %w", err)
	}

	return nil
}

// GetUsersByTag 获取拥有指定标签的所有用户（分页）
func (s *UserTagService) GetUsersByTag(ctx context.Context, tagID uint64, page, pageSize int) ([]model.User, int64, error) {
	// 检查标签是否存在
	if _, err := s.tagRepo.GetTag(ctx, tagID); err != nil {
		return nil, 0, fmt.Errorf("标签不存在: %w", err)
	}

	users, total, err := s.tagRepo.GetUsersByTag(ctx, tagID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户列表失败: %w", err)
	}

	return users, total, nil
}

// isValidColor 验证颜色格式（#FF6B6B）
func isValidColor(color string) bool {
	if color == "" {
		return true // 允许空值（不更新颜色）
	}

	// 正则表达式匹配 #RRGGBB 格式
	pattern := `^#[0-9A-Fa-f]{6}$`
	matched, _ := regexp.MatchString(pattern, color)
	return matched
}
