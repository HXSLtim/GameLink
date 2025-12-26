package statistics

import (
	"context"
	"fmt"
	"reflect"

	"gamelink/internal/model"

	"gorm.io/gorm"
)

// TagEvaluator 标签评估器
type TagEvaluator struct {
	db *gorm.DB
}

// NewTagEvaluator 创建标签评估器
func NewTagEvaluator(db *gorm.DB) *TagEvaluator {
	return &TagEvaluator{db: db}
}

// EvaluateUserTags 评估用户标签
func (e *TagEvaluator) EvaluateUserTags(ctx context.Context, userID uint64) ([]uint64, error) {
	// 获取用户统计
	var stats model.UserStatistics
	if err := e.db.WithContext(ctx).Where("user_id = ?", userID).First(&stats).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("get user statistics: %w", err)
	}

	// 获取用户类型的标签阈值配置
	var thresholds []model.TagThreshold
	if err := e.db.WithContext(ctx).
		Where("target_type = ? AND is_active = true", model.TagTargetUser).
		Order("priority DESC").
		Find(&thresholds).Error; err != nil {
		return nil, fmt.Errorf("get tag thresholds: %w", err)
	}

	// 评估每个阈值
	matchedTagIDs := make([]uint64, 0)
	for _, threshold := range thresholds {
		value, err := e.getFieldValue(&stats, threshold.MetricField)
		if err != nil {
			continue // 字段不存在，跳过
		}

		if threshold.Evaluate(value) {
			matchedTagIDs = append(matchedTagIDs, threshold.TagID)
		}
	}

	return matchedTagIDs, nil
}

// EvaluatePlayerTags 评估陪玩师标签
func (e *TagEvaluator) EvaluatePlayerTags(ctx context.Context, playerID uint64) ([]uint64, error) {
	var stats model.PlayerStatistics
	if err := e.db.WithContext(ctx).Where("player_id = ?", playerID).First(&stats).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("get player statistics: %w", err)
	}

	var thresholds []model.TagThreshold
	if err := e.db.WithContext(ctx).
		Where("target_type = ? AND is_active = true", model.TagTargetPlayer).
		Order("priority DESC").
		Find(&thresholds).Error; err != nil {
		return nil, fmt.Errorf("get tag thresholds: %w", err)
	}

	matchedTagIDs := make([]uint64, 0)
	for _, threshold := range thresholds {
		value, err := e.getFieldValue(&stats, threshold.MetricField)
		if err != nil {
			continue
		}

		if threshold.Evaluate(value) {
			matchedTagIDs = append(matchedTagIDs, threshold.TagID)
		}
	}

	return matchedTagIDs, nil
}

// getFieldValue 通过反射获取结构体字段值
func (e *TagEvaluator) getFieldValue(obj interface{}, fieldName string) (float64, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return 0, fmt.Errorf("field %s not found", fieldName)
	}

	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(field.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(field.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return field.Float(), nil
	default:
		return 0, fmt.Errorf("unsupported field type: %s", field.Kind())
	}
}

// SyncUserTags 同步用户标签（根据统计自动打标签）
func (e *TagEvaluator) SyncUserTags(ctx context.Context, userID uint64) error {
	tagIDs, err := e.EvaluateUserTags(ctx, userID)
	if err != nil {
		return err
	}

	// 删除旧的自动标签
	if err := e.db.WithContext(ctx).
		Where("user_id = ? AND is_auto = true", userID).
		Delete(&model.UserTagRelation{}).Error; err != nil {
		return fmt.Errorf("delete old auto tags: %w", err)
	}

	// 添加新的自动标签
	for _, tagID := range tagIDs {
		relation := &model.UserTagRelation{
			UserID: userID,
			TagID:  tagID,
			IsAuto: true,
		}
		if err := e.db.WithContext(ctx).Create(relation).Error; err != nil {
			// 忽略重复错误
			continue
		}
	}

	return nil
}
