package activity

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	activityrepo "gamelink/internal/repository/activity"
)

// ActivityRepository defines the interface for activity repository operations
type ActivityRepository interface {
	ListActivities(ctx context.Context, opts activityrepo.ActivityListOptions) ([]model.Activity, int64, error)
	GetActiveActivities(ctx context.Context) ([]model.Activity, error)
	GetVisibleActivities(ctx context.Context) ([]model.Activity, error)
	GetActivityByID(ctx context.Context, id uint64) (*model.Activity, error)
	CreateActivity(ctx context.Context, activity *model.Activity) error
	UpdateActivity(ctx context.Context, activity *model.Activity) error
	DeleteActivity(ctx context.Context, id uint64) error
	UpdateActivityStatus(ctx context.Context, id uint64, status model.ActivityStatus) error
	GetRewardByID(ctx context.Context, id uint64) (*model.ActivityReward, error)
	GetRewardsByActivityID(ctx context.Context, activityID uint64) ([]model.ActivityReward, error)
	CreateReward(ctx context.Context, reward *model.ActivityReward) error
	UpdateReward(ctx context.Context, reward *model.ActivityReward) error
	DeleteReward(ctx context.Context, id uint64) error
	CountUserParticipations(ctx context.Context, userID, activityID uint64) (int64, error)
	CountTodayParticipations(ctx context.Context, activityID uint64) (int64, error)
	CreateParticipation(ctx context.Context, participation *model.ActivityParticipation) error
	GetUserParticipations(ctx context.Context, userID uint64, limit int) ([]model.ActivityParticipation, error)
	ListParticipations(ctx context.Context, opts activityrepo.ParticipationListOptions) ([]model.ActivityParticipation, int64, error)
	IncrementParticipants(ctx context.Context, activityID uint64) error
	DecrementRewardStock(ctx context.Context, rewardID uint64) error
	IncrementDailyStats(ctx context.Context, activityID uint64) error
	GetActivityStats(ctx context.Context, activityID uint64) (map[string]any, error)
	GetAllActivityStats(ctx context.Context) (map[string]any, error)
	ResetTodayParticipants(ctx context.Context) error
}

// CouponService defines the interface for coupon service operations
type CouponService interface {
	IssueCoupon(ctx context.Context, userID, templateID uint64, source model.CouponSource) (*model.Coupon, error)
}

// Service 活动业务逻辑层
type Service struct {
	repo      ActivityRepository
	couponSvc CouponService
}

// NewActivityService 创建活动服务
func NewActivityService(repo ActivityRepository, couponSvc CouponService) *Service {
	return &Service{
		repo:      repo,
		couponSvc: couponSvc,
	}
}

// ============================================================================
// 活动管理
// ============================================================================

// ListActivities 获取活动列表
func (s *Service) ListActivities(ctx context.Context, opts activityrepo.ActivityListOptions) ([]model.Activity, int64, error) {
	activities, total, err := s.repo.ListActivities(ctx, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("list activities: %w", err)
	}
	return activities, total, nil
}

// GetActiveActivities 获取进行中的活动列表（用户端）
func (s *Service) GetActiveActivities(ctx context.Context) ([]model.Activity, error) {
	activities, err := s.repo.GetActiveActivities(ctx)
	if err != nil {
		return nil, fmt.Errorf("get active activities: %w", err)
	}
	return activities, nil
}

// GetVisibleActivities 获取可见的活动列表（包含预热期）
func (s *Service) GetVisibleActivities(ctx context.Context) ([]model.Activity, error) {
	activities, err := s.repo.GetVisibleActivities(ctx)
	if err != nil {
		return nil, fmt.Errorf("get visible activities: %w", err)
	}
	return activities, nil
}

// GetActivity 获取活动详情
func (s *Service) GetActivity(ctx context.Context, id uint64) (*model.Activity, error) {
	activity, err := s.repo.GetActivityByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get activity: %w", err)
	}
	return activity, nil
}

// CreateActivity 创建活动
func (s *Service) CreateActivity(ctx context.Context, activity *model.Activity) error {
	if err := s.validateActivity(activity); err != nil {
		return err
	}

	// 设置默认状态
	if activity.Status == "" {
		activity.Status = model.ActivityStatusDraft
	}

	if err := s.repo.CreateActivity(ctx, activity); err != nil {
		return fmt.Errorf("create activity: %w", err)
	}
	return nil
}

// UpdateActivity 更新活动
func (s *Service) UpdateActivity(ctx context.Context, activity *model.Activity) error {
	existing, err := s.repo.GetActivityByID(ctx, activity.ID)
	if err != nil {
		return fmt.Errorf("get activity: %w", err)
	}

	// 已结束或已取消的活动不能修改
	if existing.Status == model.ActivityStatusEnded || existing.Status == model.ActivityStatusCanceled {
		return errors.New("已结束或已取消的活动不能修改")
	}

	if err := s.validateActivity(activity); err != nil {
		return err
	}

	if err := s.repo.UpdateActivity(ctx, activity); err != nil {
		return fmt.Errorf("update activity: %w", err)
	}
	return nil
}

// DeleteActivity 删除活动
func (s *Service) DeleteActivity(ctx context.Context, id uint64) error {
	activity, err := s.repo.GetActivityByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get activity: %w", err)
	}

	// 进行中的活动不能删除
	if activity.Status == model.ActivityStatusActive {
		return errors.New("进行中的活动不能删除")
	}

	if err := s.repo.DeleteActivity(ctx, id); err != nil {
		return fmt.Errorf("delete activity: %w", err)
	}
	return nil
}

// UpdateActivityStatus 更新活动状态
func (s *Service) UpdateActivityStatus(ctx context.Context, id uint64, status model.ActivityStatus) error {
	activity, err := s.repo.GetActivityByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get activity: %w", err)
	}

	// 验证状态流转
	if err := s.validateStatusTransition(activity.Status, status); err != nil {
		return err
	}

	if err := s.repo.UpdateActivityStatus(ctx, id, status); err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	return nil
}

// validateActivity 验证活动配置
func (s *Service) validateActivity(activity *model.Activity) error {
	if activity.Name == "" {
		return errors.New("活动名称不能为空")
	}
	if activity.StartAt.IsZero() {
		return errors.New("活动开始时间不能为空")
	}
	if activity.EndAt.IsZero() {
		return errors.New("活动结束时间不能为空")
	}
	if !activity.EndAt.After(activity.StartAt) {
		return errors.New("活动结束时间必须晚于开始时间")
	}
	if activity.PreheatAt != nil && !activity.StartAt.After(*activity.PreheatAt) {
		return errors.New("活动开始时间必须晚于预热时间")
	}
	return nil
}

// validateStatusTransition 验证状态流转
func (s *Service) validateStatusTransition(from, to model.ActivityStatus) error {
	validTransitions := map[model.ActivityStatus][]model.ActivityStatus{
		model.ActivityStatusDraft:    {model.ActivityStatusPreheat, model.ActivityStatusActive, model.ActivityStatusCanceled},
		model.ActivityStatusPreheat:  {model.ActivityStatusActive, model.ActivityStatusPaused, model.ActivityStatusCanceled},
		model.ActivityStatusActive:   {model.ActivityStatusPaused, model.ActivityStatusEnded, model.ActivityStatusCanceled},
		model.ActivityStatusPaused:   {model.ActivityStatusActive, model.ActivityStatusEnded, model.ActivityStatusCanceled},
		model.ActivityStatusEnded:    {},
		model.ActivityStatusCanceled: {},
	}

	allowed, ok := validTransitions[from]
	if !ok {
		return fmt.Errorf("无效的当前状态: %s", from)
	}

	for _, s := range allowed {
		if s == to {
			return nil
		}
	}
	return fmt.Errorf("不允许从 %s 状态转换到 %s 状态", from, to)
}

// ============================================================================
// 活动奖励管理
// ============================================================================

// GetReward 获取奖励详情
func (s *Service) GetReward(ctx context.Context, id uint64) (*model.ActivityReward, error) {
	reward, err := s.repo.GetRewardByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get reward: %w", err)
	}
	return reward, nil
}

// GetRewardsByActivityID 获取活动的奖励列表
func (s *Service) GetRewardsByActivityID(ctx context.Context, activityID uint64) ([]model.ActivityReward, error) {
	rewards, err := s.repo.GetRewardsByActivityID(ctx, activityID)
	if err != nil {
		return nil, fmt.Errorf("get rewards: %w", err)
	}
	return rewards, nil
}

// CreateReward 创建奖励
func (s *Service) CreateReward(ctx context.Context, reward *model.ActivityReward) error {
	// 验证活动存在
	if _, err := s.repo.GetActivityByID(ctx, reward.ActivityID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return errors.New("活动不存在")
		}
		return fmt.Errorf("get activity: %w", err)
	}

	if err := s.validateReward(reward); err != nil {
		return err
	}

	// 设置剩余库存
	if reward.TotalStock > 0 {
		reward.RemainingStock = reward.TotalStock
	}

	if err := s.repo.CreateReward(ctx, reward); err != nil {
		return fmt.Errorf("create reward: %w", err)
	}
	return nil
}

// UpdateReward 更新奖励
func (s *Service) UpdateReward(ctx context.Context, reward *model.ActivityReward) error {
	existing, err := s.repo.GetRewardByID(ctx, reward.ID)
	if err != nil {
		return fmt.Errorf("get reward: %w", err)
	}

	if err := s.validateReward(reward); err != nil {
		return err
	}

	// 如果修改了总库存，调整剩余库存
	if reward.TotalStock != existing.TotalStock {
		diff := reward.TotalStock - existing.TotalStock
		reward.RemainingStock = existing.RemainingStock + diff
		if reward.RemainingStock < 0 {
			reward.RemainingStock = 0
		}
	}

	if err := s.repo.UpdateReward(ctx, reward); err != nil {
		return fmt.Errorf("update reward: %w", err)
	}
	return nil
}

// DeleteReward 删除奖励
func (s *Service) DeleteReward(ctx context.Context, id uint64) error {
	if err := s.repo.DeleteReward(ctx, id); err != nil {
		return fmt.Errorf("delete reward: %w", err)
	}
	return nil
}

// validateReward 验证奖励配置
func (s *Service) validateReward(reward *model.ActivityReward) error {
	if reward.CouponTemplateID == 0 {
		return errors.New("优惠券模板ID不能为空")
	}
	if reward.CouponCount <= 0 {
		return errors.New("发放数量必须大于0")
	}
	if reward.Probability < 1 || reward.Probability > 100 {
		return errors.New("发放概率必须在1-100之间")
	}
	return nil
}

// ============================================================================
// 用户参与活动
// ============================================================================

// ParticipateActivity 用户参与活动
func (s *Service) ParticipateActivity(ctx context.Context, userID, activityID, rewardID uint64, clientIP string) (*model.ActivityParticipation, error) {
	// 获取活动
	activity, err := s.repo.GetActivityByID(ctx, activityID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errors.New("活动不存在")
		}
		return nil, fmt.Errorf("get activity: %w", err)
	}

	// 检查活动状态
	if !activity.IsActive() {
		return nil, errors.New("活动未开始或已结束")
	}

	// 检查每人限制
	if activity.PerUserLimit > 0 {
		count, err := s.repo.CountUserParticipations(ctx, userID, activityID)
		if err != nil {
			return nil, fmt.Errorf("count participations: %w", err)
		}
		if int(count) >= activity.PerUserLimit {
			return nil, errors.New("已达到参与上限")
		}
	}

	// 检查每日限制
	if activity.DailyLimit > 0 {
		count, err := s.repo.CountTodayParticipations(ctx, activityID)
		if err != nil {
			return nil, fmt.Errorf("count today participations: %w", err)
		}
		if int(count) >= activity.DailyLimit {
			return nil, errors.New("今日参与名额已满")
		}
	}

	// 检查总量限制
	if activity.TotalLimit > 0 && activity.TotalClaimed >= activity.TotalLimit {
		return nil, errors.New("活动名额已满")
	}

	// 获取奖励配置
	reward, err := s.repo.GetRewardByID(ctx, rewardID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errors.New("奖励配置不存在")
		}
		return nil, fmt.Errorf("get reward: %w", err)
	}

	// 检查奖励是否属于该活动
	if reward.ActivityID != activityID {
		return nil, errors.New("奖励配置不属于该活动")
	}

	// 检查库存
	if reward.TotalStock > 0 && reward.RemainingStock <= 0 {
		return nil, errors.New("奖励库存不足")
	}

	// 发放优惠券
	var couponIDs []uint64
	if s.couponSvc != nil {
		for i := 0; i < reward.CouponCount; i++ {
			coupon, err := s.couponSvc.IssueCoupon(ctx, userID, reward.CouponTemplateID, model.CouponSourceActivity)
			if err != nil {
				fmt.Printf("issue coupon error: %v\n", err)
				continue
			}
			couponIDs = append(couponIDs, coupon.ID)
		}
	}

	// 创建参与记录
	participation := &model.ActivityParticipation{
		ActivityID: activityID,
		UserID:     userID,
		RewardID:   rewardID,
		CouponIDs:  fmt.Sprintf("%v", couponIDs),
		ClaimedAt:  time.Now(),
		ClientIP:   clientIP,
	}

	if err := s.repo.CreateParticipation(ctx, participation); err != nil {
		return nil, fmt.Errorf("create participation: %w", err)
	}

	// 更新统计
	_ = s.repo.IncrementParticipants(ctx, activityID)
	_ = s.repo.DecrementRewardStock(ctx, rewardID)
	_ = s.repo.IncrementDailyStats(ctx, activityID)

	return participation, nil
}

// GetUserParticipations 获取用户参与记录
func (s *Service) GetUserParticipations(ctx context.Context, userID uint64, limit int) ([]model.ActivityParticipation, error) {
	if limit <= 0 {
		limit = 20
	}
	participations, err := s.repo.GetUserParticipations(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get user participations: %w", err)
	}
	return participations, nil
}

// ListParticipations 获取参与记录列表（管理端）
func (s *Service) ListParticipations(ctx context.Context, opts activityrepo.ParticipationListOptions) ([]model.ActivityParticipation, int64, error) {
	participations, total, err := s.repo.ListParticipations(ctx, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("list participations: %w", err)
	}
	return participations, total, nil
}

// ============================================================================
// 统计
// ============================================================================

// GetActivityStats 获取活动统计
func (s *Service) GetActivityStats(ctx context.Context, activityID uint64) (map[string]any, error) {
	stats, err := s.repo.GetActivityStats(ctx, activityID)
	if err != nil {
		return nil, fmt.Errorf("get activity stats: %w", err)
	}
	return stats, nil
}

// GetAllActivityStats 获取所有活动统计概览
func (s *Service) GetAllActivityStats(ctx context.Context) (map[string]any, error) {
	stats, err := s.repo.GetAllActivityStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all activity stats: %w", err)
	}
	return stats, nil
}

// ============================================================================
// 定时任务
// ============================================================================

// ResetTodayParticipants 重置今日参与人数（定时任务）
func (s *Service) ResetTodayParticipants(ctx context.Context) error {
	if err := s.repo.ResetTodayParticipants(ctx); err != nil {
		return fmt.Errorf("reset today participants: %w", err)
	}
	return nil
}

// AutoUpdateActivityStatus 自动更新活动状态（定时任务）
func (s *Service) AutoUpdateActivityStatus(ctx context.Context) error {
	now := time.Now()

	// 获取所有活动
	activities, _, err := s.repo.ListActivities(ctx, activityrepo.ActivityListOptions{
		Page:     1,
		PageSize: 1000,
	})
	if err != nil {
		return fmt.Errorf("list activities: %w", err)
	}

	for _, activity := range activities {
		var newStatus model.ActivityStatus

		switch activity.Status {
		case model.ActivityStatusDraft:
			// 草稿状态：检查是否到预热时间
			if activity.PreheatAt != nil && now.After(*activity.PreheatAt) && now.Before(activity.StartAt) {
				newStatus = model.ActivityStatusPreheat
			} else if now.After(activity.StartAt) && now.Before(activity.EndAt) {
				newStatus = model.ActivityStatusActive
			}
		case model.ActivityStatusPreheat:
			// 预热状态：检查是否到开始时间
			if now.After(activity.StartAt) && now.Before(activity.EndAt) {
				newStatus = model.ActivityStatusActive
			}
		case model.ActivityStatusActive:
			// 进行中：检查是否到结束时间
			if now.After(activity.EndAt) {
				newStatus = model.ActivityStatusEnded
			}
		}

		if newStatus != "" && newStatus != activity.Status {
			_ = s.repo.UpdateActivityStatus(ctx, activity.ID, newStatus)
		}
	}

	return nil
}

// ============================================================================
// 批量操作
// ============================================================================

// BatchOperationResult 批量操作结果
type BatchOperationResult struct {
	SuccessCount int      `json:"successCount"`
	FailedCount  int      `json:"failedCount"`
	FailedIDs    []uint64 `json:"failedIds"`
	Errors       []string `json:"errors"`
}

// BatchDeleteActivities 批量删除活动
func (s *Service) BatchDeleteActivities(ctx context.Context, ids []uint64) (*BatchOperationResult, error) {
	if len(ids) == 0 {
		return nil, errors.New("活动ID列表不能为空")
	}

	result := &BatchOperationResult{
		FailedIDs: make([]uint64, 0),
		Errors:    make([]string, 0),
	}

	for _, id := range ids {
		err := s.DeleteActivity(ctx, id)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, id)
			result.Errors = append(result.Errors, fmt.Sprintf("活动%d: %s", id, err.Error()))
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

// BatchUpdateActivityStatus 批量更新活动状态
func (s *Service) BatchUpdateActivityStatus(ctx context.Context, ids []uint64, status model.ActivityStatus) (*BatchOperationResult, error) {
	if len(ids) == 0 {
		return nil, errors.New("活动ID列表不能为空")
	}

	// 验证状态值
	validStatuses := map[model.ActivityStatus]bool{
		model.ActivityStatusDraft:    true,
		model.ActivityStatusPreheat:  true,
		model.ActivityStatusActive:   true,
		model.ActivityStatusPaused:   true,
		model.ActivityStatusEnded:    true,
		model.ActivityStatusCanceled: true,
	}
	if !validStatuses[status] {
		return nil, errors.New("无效的活动状态")
	}

	result := &BatchOperationResult{
		FailedIDs: make([]uint64, 0),
		Errors:    make([]string, 0),
	}

	for _, id := range ids {
		err := s.UpdateActivityStatus(ctx, id, status)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, id)
			result.Errors = append(result.Errors, fmt.Sprintf("活动%d: %s", id, err.Error()))
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

// BatchPublishActivities 批量发布/取消发布活动
func (s *Service) BatchPublishActivities(ctx context.Context, ids []uint64, isVisible bool) (*BatchOperationResult, error) {
	if len(ids) == 0 {
		return nil, errors.New("活动ID列表不能为空")
	}

	result := &BatchOperationResult{
		FailedIDs: make([]uint64, 0),
		Errors:    make([]string, 0),
	}

	for _, id := range ids {
		activity, err := s.repo.GetActivityByID(ctx, id)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, id)
			result.Errors = append(result.Errors, fmt.Sprintf("活动%d: %s", id, err.Error()))
			continue
		}

		// 更新可见性
		activity.IsVisible = isVisible
		if err := s.repo.UpdateActivity(ctx, activity); err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, id)
			result.Errors = append(result.Errors, fmt.Sprintf("活动%d: %s", id, err.Error()))
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}
