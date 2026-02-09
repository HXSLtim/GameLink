package playercertification

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service/external"
)

var (
	// ErrNotFound 记录不存在
	ErrNotFound = repository.ErrNotFound
	// ErrValidation 输入校验失败
	ErrValidation = errors.New("validation failed")
	// ErrPlayerNotFound 陪玩师不存在
	ErrPlayerNotFound = errors.New("player not found")
	// ErrAlreadyCertified 已完成实名认证
	ErrAlreadyCertified = errors.New("already certified")
	// ErrPendingCertification 有待审核的认证申请
	ErrPendingCertification = errors.New("pending certification exists")
	// ErrInvalidStatus 无效的状态变更
	ErrInvalidStatus = errors.New("invalid status transition")
)

// PlayerCertificationService 陪玩师实名认证服务
//
// 第三方接口依赖:
// - IdentityVerifier: 身份证实名验证（二要素验证）
//
// TODO: 生产环境部署前需要配置真实的身份证验证服务提供商
// 参考: api/internal/service/external/thirdparty.go
type PlayerCertificationService struct {
	certs            repository.PlayerCertificationRepository
	players          repository.PlayerRepository
	identityVerifier external.IdentityVerifier // 身份证验证器
}

// NewPlayerCertificationService 创建陪玩师实名认证服务
func NewPlayerCertificationService(
	certs repository.PlayerCertificationRepository,
	players repository.PlayerRepository,
) *PlayerCertificationService {
	return &PlayerCertificationService{
		certs:            certs,
		players:          players,
		identityVerifier: external.NewMockIdentityVerifier(), // 默认使用 Mock 验证器
	}
}

// SetIdentityVerifier 设置身份证验证器
// 生产环境应调用此方法设置真实的验证器
func (s *PlayerCertificationService) SetIdentityVerifier(verifier external.IdentityVerifier) {
	s.identityVerifier = verifier
}

// ApplyInput 申请实名认证输入
type ApplyInput struct {
	PlayerID       uint64
	RealName       string
	IDCardNo       string
	IDCardFrontURL string
	IDCardBackURL  string
	PhotoURL       string // 可选
	VoiceURL       string // 可选
}

// VerifyInput 审核实名认证输入
type VerifyInput struct {
	CertID       uint64
	Status       model.CertificationStatus // verified/rejected
	VerifiedBy   uint64
	RejectReason string
}

// Apply 申请实名认证
//
// 流程:
// 1. 验证陪玩师是否存在
// 2. 验证输入参数
// 3. 调用第三方身份证验证接口（验证姓名和身份证号是否匹配）
// 4. 检查是否已有认证记录
// 5. 创建或更新认证记录
//
// TODO: 生产环境需要配置真实的身份证验证服务
func (s *PlayerCertificationService) Apply(ctx context.Context, input ApplyInput) (*model.PlayerCertification, error) {
	// 验证陪玩师是否存在
	if _, err := s.players.Get(ctx, input.PlayerID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrPlayerNotFound
		}
		return nil, fmt.Errorf("get player: %w", err)
	}

	// 验证输入
	if input.RealName == "" {
		return nil, fmt.Errorf("%w: real name is required", ErrValidation)
	}
	if input.IDCardNo == "" {
		return nil, fmt.Errorf("%w: id card number is required", ErrValidation)
	}
	if input.IDCardFrontURL == "" {
		return nil, fmt.Errorf("%w: id card front image is required", ErrValidation)
	}
	if input.IDCardBackURL == "" {
		return nil, fmt.Errorf("%w: id card back image is required", ErrValidation)
	}

	// =========================================================================
	// 第三方接口调用：身份证实名验证
	// TODO: 生产环境需要配置真实的身份证验证服务提供商（阿里云/腾讯云）
	// 当前使用 Mock 验证器，所有格式正确的身份证都会通过
	// =========================================================================
	if s.identityVerifier != nil {
		result, err := s.identityVerifier.VerifyIdentity(ctx, input.RealName, input.IDCardNo)
		if err != nil {
			return nil, fmt.Errorf("identity verification failed: %w", err)
		}
		if !result.Verified {
			return nil, fmt.Errorf("%w: %s", ErrValidation, result.Message)
		}
	}

	// 检查是否已有认证记录
	existing, err := s.certs.GetByPlayerID(ctx, input.PlayerID)
	if err == nil {
		switch existing.Status {
		case model.CertificationStatusVerified:
			return nil, ErrAlreadyCertified
		case model.CertificationStatusPending:
			return nil, ErrPendingCertification
		case model.CertificationStatusRejected:
			// 被拒绝后可以重新申请，更新现有记录
			existing.RealName = input.RealName
			existing.IDCardNo = input.IDCardNo
			existing.IDCardFrontURL = input.IDCardFrontURL
			existing.IDCardBackURL = input.IDCardBackURL
			existing.PhotoURL = input.PhotoURL
			existing.VoiceURL = input.VoiceURL
			existing.Status = model.CertificationStatusPending
			existing.RejectReason = ""
			existing.VerifiedAt = nil
			existing.VerifiedBy = nil

			if err := s.certs.Update(ctx, existing); err != nil {
				return nil, fmt.Errorf("update certification: %w", err)
			}
			return existing, nil
		}
	}

	// 创建新的认证记录
	cert := &model.PlayerCertification{
		PlayerID:       input.PlayerID,
		RealName:       input.RealName,
		IDCardNo:       input.IDCardNo,
		IDCardFrontURL: input.IDCardFrontURL,
		IDCardBackURL:  input.IDCardBackURL,
		PhotoURL:       input.PhotoURL,
		VoiceURL:       input.VoiceURL,
		Status:         model.CertificationStatusPending,
	}

	if err := s.certs.Create(ctx, cert); err != nil {
		return nil, fmt.Errorf("create certification: %w", err)
	}

	return cert, nil
}

// Verify 审核实名认证
func (s *PlayerCertificationService) Verify(ctx context.Context, input VerifyInput) (*model.PlayerCertification, error) {
	cert, err := s.certs.Get(ctx, input.CertID)
	if err != nil {
		return nil, err
	}

	// 验证状态变更是否合法
	if cert.Status != model.CertificationStatusPending {
		return nil, fmt.Errorf("%w: can only verify pending certifications", ErrInvalidStatus)
	}

	if input.Status != model.CertificationStatusVerified && input.Status != model.CertificationStatusRejected {
		return nil, fmt.Errorf("%w: status must be verified or rejected", ErrValidation)
	}

	// 更新状态
	now := time.Now()
	cert.Status = input.Status
	cert.VerifiedBy = &input.VerifiedBy
	cert.RejectReason = input.RejectReason

	if input.Status == model.CertificationStatusVerified {
		cert.VerifiedAt = &now
	}

	if err := s.certs.Update(ctx, cert); err != nil {
		return nil, fmt.Errorf("update certification: %w", err)
	}

	return cert, nil
}

// Get 获取实名认证记录
func (s *PlayerCertificationService) Get(ctx context.Context, id uint64) (*model.PlayerCertification, error) {
	return s.certs.GetWithPlayer(ctx, id)
}

// GetByPlayerID 根据陪玩师ID获取实名认证记录
func (s *PlayerCertificationService) GetByPlayerID(ctx context.Context, playerID uint64) (*model.PlayerCertification, error) {
	return s.certs.GetByPlayerID(ctx, playerID)
}

// ListPaged 分页获取实名认证记录
func (s *PlayerCertificationService) ListPaged(ctx context.Context, opts repository.PlayerCertificationListOptions) ([]model.PlayerCertification, *model.Pagination, error) {
	certs, total, err := s.certs.ListPaged(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	pagination := &model.Pagination{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Total:    int(total),
	}

	return certs, pagination, nil
}

// ListPending 获取待审核的实名认证记录
func (s *PlayerCertificationService) ListPending(ctx context.Context, page, pageSize int) ([]model.PlayerCertification, *model.Pagination, error) {
	certs, total, err := s.certs.ListPending(ctx, page, pageSize)
	if err != nil {
		return nil, nil, err
	}

	pagination := &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	}

	return certs, pagination, nil
}

// Delete 删除实名认证记录
func (s *PlayerCertificationService) Delete(ctx context.Context, id uint64) error {
	return s.certs.Delete(ctx, id)
}

// GetStats 获取实名认证统计
func (s *PlayerCertificationService) GetStats(ctx context.Context) (map[model.CertificationStatus]int64, error) {
	return s.certs.CountByStatus(ctx)
}

// GetPendingCount 获取待审核数量
func (s *PlayerCertificationService) GetPendingCount(ctx context.Context) (int64, error) {
	return s.certs.GetPendingCount(ctx)
}
