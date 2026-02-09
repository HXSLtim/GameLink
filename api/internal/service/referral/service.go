package referral

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gamelink/internal/model"
	referralrepo "gamelink/internal/repository/referral"
)

// Repository 定义推荐仓库接口
type Repository interface {
	GetAllConfigs(ctx context.Context) ([]model.ReferralConfig, error)
	GetConfig(ctx context.Context, key string) (*model.ReferralConfig, error)
	SetConfig(ctx context.Context, key, value, description string) error
	ListCodes(ctx context.Context, opts referralrepo.CodeListOptions) ([]model.ReferralCode, int64, error)
	GetCodeByID(ctx context.Context, id uint64) (*model.ReferralCode, error)
	GetCodeByCode(ctx context.Context, code string) (*model.ReferralCode, error)
	GetUserCode(ctx context.Context, userID uint64, refType model.ReferralType) (*model.ReferralCode, error)
	CreateCode(ctx context.Context, item *model.ReferralCode) error
	UpdateCode(ctx context.Context, item *model.ReferralCode) error
	DeleteCode(ctx context.Context, id uint64) error
	IncrementCodeUseCount(ctx context.Context, codeID uint64) error
	ListReferrals(ctx context.Context, opts referralrepo.ReferralListOptions) ([]model.Referral, int64, error)
	GetReferralByID(ctx context.Context, id uint64) (*model.Referral, error)
	GetReferralByReferee(ctx context.Context, refereeID uint64) (*model.Referral, error)
	CreateReferral(ctx context.Context, item *model.Referral) error
	UpdateReferralStatus(ctx context.Context, id uint64, status model.ReferralStatus) error
	DeleteReferral(ctx context.Context, id uint64) error
	GetUserReferrals(ctx context.Context, userID uint64, limit int) ([]model.Referral, error)
	ListRewards(ctx context.Context, opts referralrepo.RewardListOptions) ([]model.ReferralReward, int64, error)
	GetRewardByID(ctx context.Context, id uint64) (*model.ReferralReward, error)
	CreateReward(ctx context.Context, item *model.ReferralReward) error
	UpdateRewardStatus(ctx context.Context, id uint64, status model.ReferralRewardStatus, failureReason string) error
	GetUserRewards(ctx context.Context, userID uint64, limit int) ([]model.ReferralReward, error)
	GetReferralStats(ctx context.Context) (map[string]any, error)
	GetUserReferralStats(ctx context.Context, userID uint64) (map[string]any, error)
}

// Service 推荐服务
type Service struct {
	repo Repository
}

// NewReferralService 创建推荐服务
func NewReferralService(repo Repository) *Service {
	return &Service{repo: repo}
}

// ============================================================================
// 配置管理
// ============================================================================

// GetAllConfigs 获取所有配置
func (s *Service) GetAllConfigs(ctx context.Context) ([]model.ReferralConfig, error) {
	return s.repo.GetAllConfigs(ctx)
}

// GetConfig 获取配置
func (s *Service) GetConfig(ctx context.Context, key string) (string, error) {
	config, err := s.repo.GetConfig(ctx, key)
	if err != nil {
		return "", err
	}
	return config.ConfigValue, nil
}

// SetConfig 设置配置
func (s *Service) SetConfig(ctx context.Context, key, value, description string) error {
	return s.repo.SetConfig(ctx, key, value, description)
}

// IsEnabled 检查推荐系统是否启用
func (s *Service) IsEnabled(ctx context.Context) bool {
	value, err := s.GetConfig(ctx, model.ReferralConfigEnabled)
	if err != nil {
		return false
	}
	return value == "true" || value == "1"
}

// GetExpireDays 获取邀请码过期天数
func (s *Service) GetExpireDays(ctx context.Context) int {
	value, err := s.GetConfig(ctx, model.ReferralConfigExpireDays)
	if err != nil {
		return 30 // 默认30天
	}
	days, _ := strconv.Atoi(value)
	if days <= 0 {
		return 30
	}
	return days
}

func (s *Service) resolveRewardConfig(ctx context.Context, refType model.ReferralType) (model.RewardType, int64) {
	typeKey := model.ReferralConfigUserRewardType
	amountKey := model.ReferralConfigUserRewardAmount
	switch refType {
	case model.ReferralTypePlayerToPlayer, model.ReferralTypeUserToPlayer:
		typeKey = model.ReferralConfigPlayerRewardType
		amountKey = model.ReferralConfigPlayerRewardAmount
	}

	rewardType := model.RewardTypeCash
	if value, err := s.GetConfig(ctx, typeKey); err == nil && value != "" {
		switch value {
		case string(model.RewardTypeCash):
			rewardType = model.RewardTypeCash
		case string(model.RewardTypeCoupon):
			rewardType = model.RewardTypeCoupon
		case string(model.RewardTypePoints):
			rewardType = model.RewardTypePoints
		}
	}

	var amount int64
	if value, err := s.GetConfig(ctx, amountKey); err == nil && value != "" {
		if parsed, parseErr := strconv.ParseInt(value, 10, 64); parseErr == nil && parsed > 0 {
			amount = parsed
		}
	}

	return rewardType, amount
}

// ============================================================================
// 邀请码管理
// ============================================================================

// ListCodes 获取邀请码列表
func (s *Service) ListCodes(ctx context.Context, opts referralrepo.CodeListOptions) ([]model.ReferralCode, int64, error) {
	return s.repo.ListCodes(ctx, opts)
}

// GetCodeByID 根据ID获取邀请码
func (s *Service) GetCodeByID(ctx context.Context, id uint64) (*model.ReferralCode, error) {
	return s.repo.GetCodeByID(ctx, id)
}

// GetCodeByCode 根据邀请码获取
func (s *Service) GetCodeByCode(ctx context.Context, code string) (*model.ReferralCode, error) {
	return s.repo.GetCodeByCode(ctx, code)
}

// CreateCodeRequest 创建邀请码请求
type CreateCodeRequest struct {
	UserID   uint64
	Type     model.ReferralType
	MaxUse   int
	ExpireAt *time.Time
}

// CreateCode 创建邀请码
func (s *Service) CreateCode(ctx context.Context, req CreateCodeRequest) (*model.ReferralCode, error) {
	code := &model.ReferralCode{
		UserID:   req.UserID,
		Type:     req.Type,
		IsActive: true,
		MaxUse:   req.MaxUse,
		ExpireAt: req.ExpireAt,
	}

	if err := s.repo.CreateCode(ctx, code); err != nil {
		return nil, fmt.Errorf("create code: %w", err)
	}

	return code, nil
}

// GetOrCreateUserCode 获取或创建用户邀请码
func (s *Service) GetOrCreateUserCode(ctx context.Context, userID uint64, refType model.ReferralType) (*model.ReferralCode, error) {
	// 先尝试获取现有邀请码
	code, err := s.repo.GetUserCode(ctx, userID, refType)
	if err == nil {
		return code, nil
	}

	// 不存在则创建新的
	expireDays := s.GetExpireDays(ctx)
	expireAt := time.Now().AddDate(0, 0, expireDays)

	return s.CreateCode(ctx, CreateCodeRequest{
		UserID:   userID,
		Type:     refType,
		ExpireAt: &expireAt,
	})
}

// UpdateCodeRequest 更新邀请码请求
type UpdateCodeRequest struct {
	ID       uint64
	IsActive *bool
	MaxUse   *int
	ExpireAt *time.Time
}

// UpdateCode 更新邀请码
func (s *Service) UpdateCode(ctx context.Context, req UpdateCodeRequest) (*model.ReferralCode, error) {
	code, err := s.repo.GetCodeByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if req.IsActive != nil {
		code.IsActive = *req.IsActive
	}
	if req.MaxUse != nil {
		code.MaxUse = *req.MaxUse
	}
	if req.ExpireAt != nil {
		code.ExpireAt = req.ExpireAt
	}

	if err := s.repo.UpdateCode(ctx, code); err != nil {
		return nil, fmt.Errorf("update code: %w", err)
	}

	return code, nil
}

// DeleteCode 删除邀请码
func (s *Service) DeleteCode(ctx context.Context, id uint64) error {
	return s.repo.DeleteCode(ctx, id)
}

// ValidateCode 验证邀请码
func (s *Service) ValidateCode(ctx context.Context, codeStr string) (*model.ReferralCode, error) {
	code, err := s.repo.GetCodeByCode(ctx, codeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid code: %w", err)
	}

	if !code.IsValid() {
		return nil, fmt.Errorf("code is expired or reached max use")
	}

	return code, nil
}

// ============================================================================
// 推荐记录管理
// ============================================================================

// ListReferrals 获取推荐记录列表
func (s *Service) ListReferrals(ctx context.Context, opts referralrepo.ReferralListOptions) ([]model.Referral, int64, error) {
	return s.repo.ListReferrals(ctx, opts)
}

// GetReferralByID 根据ID获取推荐记录
func (s *Service) GetReferralByID(ctx context.Context, id uint64) (*model.Referral, error) {
	return s.repo.GetReferralByID(ctx, id)
}

// CreateReferralRequest 创建推荐记录请求
type CreateReferralRequest struct {
	ReferrerID       uint64
	RefereeID        uint64
	CodeID           *uint64
	Type             model.ReferralType
	Level            int
	RefereeCondition string
}

// CreateReferral 创建推荐记录
func (s *Service) CreateReferral(ctx context.Context, req CreateReferralRequest) (*model.Referral, error) {
	// 检查被推荐人是否已有推荐记录
	_, err := s.repo.GetReferralByReferee(ctx, req.RefereeID)
	if err == nil {
		return nil, fmt.Errorf("referee already has referral record")
	}

	level := req.Level
	if level <= 0 {
		level = 1
	}

	rewardType, rewardAmount := s.resolveRewardConfig(ctx, req.Type)

	referral := &model.Referral{
		ReferrerID:       req.ReferrerID,
		RefereeID:        req.RefereeID,
		CodeID:           req.CodeID,
		Type:             req.Type,
		Level:            level,
		Status:           model.ReferralStatusPending,
		RefereeCondition: req.RefereeCondition,
		RewardType:       rewardType,
		RewardAmountCents: rewardAmount,
	}

	if err := s.repo.CreateReferral(ctx, referral); err != nil {
		return nil, fmt.Errorf("create referral: %w", err)
	}

	// 增加邀请码使用次数
	if req.CodeID != nil {
		_ = s.repo.IncrementCodeUseCount(ctx, *req.CodeID)
	}

	return referral, nil
}

// UseCodeRequest 使用邀请码请求
type UseCodeRequest struct {
	Code      string
	RefereeID uint64
}

// UseCode 使用邀请码（被邀请人注册时调用）
func (s *Service) UseCode(ctx context.Context, req UseCodeRequest) (*model.Referral, error) {
	// 验证邀请码
	code, err := s.ValidateCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	// 不能邀请自己
	if code.UserID == req.RefereeID {
		return nil, fmt.Errorf("cannot refer yourself")
	}

	// 创建推荐记录
	return s.CreateReferral(ctx, CreateReferralRequest{
		ReferrerID:       code.UserID,
		RefereeID:        req.RefereeID,
		CodeID:           &code.ID,
		Type:             code.Type,
		Level:            1,
		RefereeCondition: "registered",
	})
}

// CompleteReferral 完成推荐（被推荐人满足条件时调用）
func (s *Service) CompleteReferral(ctx context.Context, refereeID uint64, condition string) error {
	referral, err := s.repo.GetReferralByReferee(ctx, refereeID)
	if err != nil {
		return fmt.Errorf("get referral: %w", err)
	}

	// 检查条件是否匹配
	if referral.RefereeCondition != "" && referral.RefereeCondition != condition {
		return nil // 条件不匹配，不处理
	}

	// 已完成则跳过
	if referral.Status != model.ReferralStatusPending {
		return nil
	}

	// 更新状态为已完成
	if err := s.repo.UpdateReferralStatus(ctx, referral.ID, model.ReferralStatusCompleted); err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	return nil
}

// UpdateReferralStatus 更新推荐状态（管理员）
func (s *Service) UpdateReferralStatus(ctx context.Context, id uint64, status model.ReferralStatus) error {
	return s.repo.UpdateReferralStatus(ctx, id, status)
}

// GetUserReferrals 获取用户的推荐记录
func (s *Service) GetUserReferrals(ctx context.Context, userID uint64, limit int) ([]model.Referral, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.GetUserReferrals(ctx, userID, limit)
}

// ============================================================================
// 奖励管理
// ============================================================================

// ListRewards 获取奖励记录列表
func (s *Service) ListRewards(ctx context.Context, opts referralrepo.RewardListOptions) ([]model.ReferralReward, int64, error) {
	return s.repo.ListRewards(ctx, opts)
}

// GetRewardByID 根据ID获取奖励记录
func (s *Service) GetRewardByID(ctx context.Context, id uint64) (*model.ReferralReward, error) {
	return s.repo.GetRewardByID(ctx, id)
}

// CreateRewardRequest 创建奖励请求
type CreateRewardRequest struct {
	ReferralID  uint64
	UserID      uint64
	Type        model.RewardType
	AmountCents int64
	CouponID    *uint64
}

// CreateReward 创建奖励记录
func (s *Service) CreateReward(ctx context.Context, req CreateRewardRequest) (*model.ReferralReward, error) {
	reward := &model.ReferralReward{
		ReferralID:  req.ReferralID,
		UserID:      req.UserID,
		Type:        req.Type,
		AmountCents: req.AmountCents,
		CouponID:    req.CouponID,
		Status:      model.ReferralRewardStatusPending,
	}

	if err := s.repo.CreateReward(ctx, reward); err != nil {
		return nil, fmt.Errorf("create reward: %w", err)
	}

	return reward, nil
}

// IssueReward 发放奖励
func (s *Service) IssueReward(ctx context.Context, rewardID uint64) error {
	reward, err := s.repo.GetRewardByID(ctx, rewardID)
	if err != nil {
		return err
	}

	if reward.Status != model.ReferralRewardStatusPending {
		return fmt.Errorf("reward already processed")
	}

	// Issue reward based on reward type (future implementation)
	// - cash: Call wallet service to increase balance
	// - coupon: Call coupon service to issue coupon
	// - points: Call points service to add points
	// For now, just mark as issued without actual distribution

	// 更新状态为已发放
	if err := s.repo.UpdateRewardStatus(ctx, rewardID, model.ReferralRewardStatusIssued, ""); err != nil {
		return fmt.Errorf("update reward status: %w", err)
	}

	// 更新推荐记录状态
	if err := s.repo.UpdateReferralStatus(ctx, reward.ReferralID, model.ReferralStatusRewarded); err != nil {
		return fmt.Errorf("update referral status: %w", err)
	}

	return nil
}

// FailReward 标记奖励发放失败
func (s *Service) FailReward(ctx context.Context, rewardID uint64, reason string) error {
	return s.repo.UpdateRewardStatus(ctx, rewardID, model.ReferralRewardStatusFailed, reason)
}

// GetUserRewards 获取用户的奖励记录
func (s *Service) GetUserRewards(ctx context.Context, userID uint64, limit int) ([]model.ReferralReward, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.GetUserRewards(ctx, userID, limit)
}

// ============================================================================
// 统计
// ============================================================================

// GetReferralStats 获取推荐统计
func (s *Service) GetReferralStats(ctx context.Context) (map[string]any, error) {
	return s.repo.GetReferralStats(ctx)
}

// GetUserReferralStats 获取用户推荐统计
func (s *Service) GetUserReferralStats(ctx context.Context, userID uint64) (map[string]any, error) {
	return s.repo.GetUserReferralStats(ctx, userID)
}

// ============================================================================
// 批量操作
// ============================================================================

// BatchOperationResult 批量操作结果
type BatchOperationResult struct {
	SuccessCount int      `json:"successCount"`
	FailedCount  int      `json:"failedCount"`
	FailedIDs    []uint64 `json:"failedIds,omitempty"`
	TotalCount   int      `json:"totalCount"`
}

// BatchUpdateCodesStatusRequest 批量更新邀请码状态请求
type BatchUpdateCodesStatusRequest struct {
	IDs      []uint64
	IsActive bool
}

// BatchUpdateCodesStatus 批量更新邀请码状态
func (s *Service) BatchUpdateCodesStatus(ctx context.Context, ids []uint64, isActive bool) (*BatchOperationResult, error) {
	result := &BatchOperationResult{
		TotalCount: len(ids),
		FailedIDs:  make([]uint64, 0),
	}

	for _, id := range ids {
		req := UpdateCodeRequest{
			ID:       id,
			IsActive: &isActive,
		}
		if _, err := s.UpdateCode(ctx, req); err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, id)
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

// BatchDeleteCodes 批量删除邀请码
func (s *Service) BatchDeleteCodes(ctx context.Context, ids []uint64) (*BatchOperationResult, error) {
	result := &BatchOperationResult{
		TotalCount: len(ids),
		FailedIDs:  make([]uint64, 0),
	}

	for _, id := range ids {
		if err := s.DeleteCode(ctx, id); err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, id)
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

// BatchUpdateReferralsStatus 批量更新推荐状态
func (s *Service) BatchUpdateReferralsStatus(ctx context.Context, ids []uint64, status model.ReferralStatus) (*BatchOperationResult, error) {
	result := &BatchOperationResult{
		TotalCount: len(ids),
		FailedIDs:  make([]uint64, 0),
	}

	for _, id := range ids {
		if err := s.UpdateReferralStatus(ctx, id, status); err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, id)
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

// BatchDeleteReferrals 批量删除推荐记录
func (s *Service) BatchDeleteReferrals(ctx context.Context, ids []uint64) (*BatchOperationResult, error) {
	result := &BatchOperationResult{
		TotalCount: len(ids),
		FailedIDs:  make([]uint64, 0),
	}

	for _, id := range ids {
		if err := s.repo.DeleteReferral(ctx, id); err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, id)
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}
