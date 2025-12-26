package referral

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gamelink/internal/model"
	referralrepo "gamelink/internal/repository/referral"
)

// Service 推荐服务
type Service struct {
	repo *referralrepo.Repository
}

// NewReferralService 创建推荐服务
func NewReferralService(repo *referralrepo.Repository) *Service {
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

	referral := &model.Referral{
		ReferrerID:       req.ReferrerID,
		RefereeID:        req.RefereeID,
		CodeID:           req.CodeID,
		Type:             req.Type,
		Level:            level,
		Status:           model.ReferralStatusPending,
		RefereeCondition: req.RefereeCondition,
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

	// TODO: 根据奖励类型执行发放逻辑
	// - cash: 调用钱包服务增加余额
	// - coupon: 调用优惠券服务发放优惠券
	// - points: 调用积分服务增加积分

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
