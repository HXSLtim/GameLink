package collectionentity

import (
	"context"
	"errors"
	"fmt"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/collectionentity"
	"gamelink/pkg/apierr"
)

var (
	// ErrNotFound 收款主体不存在
	ErrNotFound = apierr.NotFound("collection entity not found")
	// ErrValidation 表示输入校验失败
	ErrValidation = apierr.BadRequest("validation failed")
	// ErrCreditCodeExists 统一社会信用代码已存在
	ErrCreditCodeExists = apierr.Conflict("credit code already exists")
	// ErrInvalidCreditCode 统一社会信用代码格式无效
	ErrInvalidCreditCode = apierr.BadRequest("invalid credit code format")
	// ErrEntityInactive 收款主体已禁用
	ErrEntityInactive = apierr.BadRequest("collection entity is inactive")
	// ErrChannelNotFound 支付渠道配置不存在
	ErrChannelNotFound = apierr.NotFound("payment channel config not found")
	// ErrChannelExists 支付渠道配置已存在
	ErrChannelExists = apierr.Conflict("payment channel config already exists")
	// ErrDefaultEntityRequired 必须有默认收款主体
	ErrDefaultEntityRequired = apierr.BadRequest("at least one default collection entity is required")
)

// CollectionEntityService 收款主体服务
// Requirements: 15.1, 15.2, 15.3, 15.4, 15.5
type CollectionEntityService struct {
	repo collectionentity.CollectionEntityRepository
}

// NewCollectionEntityService 创建收款主体服务
func NewCollectionEntityService(
	repo collectionentity.CollectionEntityRepository,
) *CollectionEntityService {
	return &CollectionEntityService{
		repo: repo,
	}
}

// CreateEntity 创建收款主体
// Requirements: 15.1, 15.2
func (s *CollectionEntityService) CreateEntity(ctx context.Context, req *model.CreateCollectionEntityRequest, createdBy uint64) (*model.CollectionEntity, error) {
	// 验证统一社会信用代码格式
	// Property 11: 统一社会信用代码格式验证
	if !model.ValidateCreditCode(req.CreditCode) {
		return nil, ErrInvalidCreditCode.WithDetails("credit code must be 18 characters and match the standard format")
	}

	// 检查统一社会信用代码唯一性
	existing, err := s.repo.GetByCreditCode(ctx, req.CreditCode)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, apierr.InternalError("failed to check credit code uniqueness").WithDetails(err.Error())
	}
	if existing != nil {
		return nil, ErrCreditCodeExists.WithDetails(fmt.Sprintf("credit code %s already exists", req.CreditCode))
	}

	// 创建收款主体
	entity := &model.CollectionEntity{
		Name:              req.Name,
		CreditCode:        req.CreditCode,
		TaxRegistrationNo: req.TaxRegistrationNo,
		Status:            model.EntityStatusActive,
		IsDefault:         req.IsDefault,
		CreatedBy:         createdBy,
	}

	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, apierr.InternalError("failed to create collection entity").WithDetails(err.Error())
	}

	// 如果设置为默认，需要取消其他默认主体
	if req.IsDefault {
		if err := s.repo.SetDefault(ctx, entity.ID); err != nil {
			// 记录错误但不影响主流程
		}
	}

	return entity, nil
}

// GetEntity 获取收款主体
// Requirements: 15.3
func (s *CollectionEntityService) GetEntity(ctx context.Context, id uint64) (*model.CollectionEntity, error) {
	entity, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("failed to get collection entity").WithDetails(err.Error())
	}
	return entity, nil
}

// UpdateEntity 更新收款主体
// Requirements: 15.1
func (s *CollectionEntityService) UpdateEntity(ctx context.Context, id uint64, req *model.UpdateCollectionEntityRequest, updatedBy uint64) (*model.CollectionEntity, error) {
	entity, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("failed to get collection entity").WithDetails(err.Error())
	}

	// 记录修改历史
	changes := s.detectChanges(entity, req)

	// 应用更新
	if req.Name != nil {
		entity.Name = *req.Name
	}
	if req.TaxRegistrationNo != nil {
		entity.TaxRegistrationNo = *req.TaxRegistrationNo
	}
	entity.UpdatedBy = &updatedBy

	if err := s.repo.Update(ctx, entity); err != nil {
		return nil, apierr.InternalError("failed to update collection entity").WithDetails(err.Error())
	}

	// 保存修改历史
	for _, change := range changes {
		change.CollectionEntityID = id
		change.ChangedBy = updatedBy
		if err := s.repo.CreateHistory(ctx, &change); err != nil {
			// 记录历史失败不影响主流程
			continue
		}
	}

	return entity, nil
}

// detectChanges 检测字段变更
func (s *CollectionEntityService) detectChanges(entity *model.CollectionEntity, req *model.UpdateCollectionEntityRequest) []model.CollectionEntityHistory {
	var changes []model.CollectionEntityHistory

	if req.Name != nil && *req.Name != entity.Name {
		changes = append(changes, model.CollectionEntityHistory{
			FieldName: "name",
			OldValue:  entity.Name,
			NewValue:  *req.Name,
		})
	}
	if req.TaxRegistrationNo != nil && *req.TaxRegistrationNo != entity.TaxRegistrationNo {
		changes = append(changes, model.CollectionEntityHistory{
			FieldName: "tax_registration_no",
			OldValue:  entity.TaxRegistrationNo,
			NewValue:  *req.TaxRegistrationNo,
		})
	}

	return changes
}

// ListEntities 查询收款主体列表
// Requirements: 15.3
func (s *CollectionEntityService) ListEntities(ctx context.Context, req *model.ListCollectionEntitiesRequest) (*model.ListCollectionEntitiesResponse, error) {
	opts := collectionentity.ListOptions{
		Keyword:   req.Keyword,
		Page:      req.Page,
		PageSize:  req.PageSize,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	if req.Status != "" {
		opts.Status = &req.Status
	}

	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 || opts.PageSize > 100 {
		opts.PageSize = 20
	}

	entities, total, err := s.repo.List(ctx, opts)
	if err != nil {
		return nil, apierr.InternalError("failed to list collection entities").WithDetails(err.Error())
	}

	return &model.ListCollectionEntitiesResponse{
		Total:    total,
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Entities: entities,
	}, nil
}

// ToggleEntityStatus 切换收款主体状态
// Requirements: 15.5
func (s *CollectionEntityService) ToggleEntityStatus(ctx context.Context, id uint64, enabled bool, updatedBy uint64) error {
	entity, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return apierr.InternalError("failed to get collection entity").WithDetails(err.Error())
	}

	var newStatus model.EntityStatus
	if enabled {
		newStatus = model.EntityStatusActive
	} else {
		newStatus = model.EntityStatusInactive
	}

	// 记录状态变更历史
	if entity.Status != newStatus {
		history := &model.CollectionEntityHistory{
			CollectionEntityID: id,
			FieldName:          "status",
			OldValue:           string(entity.Status),
			NewValue:           string(newStatus),
			ChangedBy:          updatedBy,
		}
		if err := s.repo.CreateHistory(ctx, history); err != nil {
			// 记录历史失败不影响主流程
		}
	}

	if err := s.repo.ToggleStatus(ctx, id, newStatus); err != nil {
		return apierr.InternalError("failed to toggle collection entity status").WithDetails(err.Error())
	}

	return nil
}

// GetEntityHistory 获取收款主体修改历史
func (s *CollectionEntityService) GetEntityHistory(ctx context.Context, entityID uint64) ([]model.CollectionEntityHistory, error) {
	// 先验证主体存在
	_, err := s.repo.Get(ctx, entityID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("failed to get collection entity").WithDetails(err.Error())
	}

	histories, err := s.repo.GetHistory(ctx, entityID)
	if err != nil {
		return nil, apierr.InternalError("failed to get collection entity history").WithDetails(err.Error())
	}

	return histories, nil
}

// SetDefaultEntity 设置默认收款主体
// Requirements: 16.3
func (s *CollectionEntityService) SetDefaultEntity(ctx context.Context, id uint64, updatedBy uint64) error {
	entity, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return apierr.InternalError("failed to get collection entity").WithDetails(err.Error())
	}

	// 验证主体是活跃状态
	if !entity.IsActive() {
		return ErrEntityInactive.WithDetails("cannot set inactive entity as default")
	}

	if err := s.repo.SetDefault(ctx, id); err != nil {
		return apierr.InternalError("failed to set default collection entity").WithDetails(err.Error())
	}

	// 记录历史
	history := &model.CollectionEntityHistory{
		CollectionEntityID: id,
		FieldName:          "is_default",
		OldValue:           fmt.Sprintf("%v", entity.IsDefault),
		NewValue:           "true",
		ChangedBy:          updatedBy,
	}
	if err := s.repo.CreateHistory(ctx, history); err != nil {
		// 记录历史失败不影响主流程
	}

	return nil
}

// GetDefaultEntity 获取默认收款主体
// Requirements: 16.3, 17.4
func (s *CollectionEntityService) GetDefaultEntity(ctx context.Context) (*model.CollectionEntity, error) {
	entity, err := s.repo.GetDefault(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound.WithDetails("no default collection entity found")
		}
		return nil, apierr.InternalError("failed to get default collection entity").WithDetails(err.Error())
	}
	return entity, nil
}

// ConfigurePaymentChannel 配置支付渠道
// Requirements: 15.4
func (s *CollectionEntityService) ConfigurePaymentChannel(ctx context.Context, entityID uint64, req *model.ConfigurePaymentChannelRequest, createdBy uint64) (*model.PaymentChannelConfig, error) {
	// 验证收款主体存在
	entity, err := s.repo.Get(ctx, entityID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("failed to get collection entity").WithDetails(err.Error())
	}

	// 验证主体是活跃状态
	if !entity.IsActive() {
		return nil, ErrEntityInactive.WithDetails("cannot configure payment channel for inactive entity")
	}

	// 检查是否已存在相同渠道配置
	existing, err := s.repo.GetChannelByEntityAndMethod(ctx, entityID, req.Channel)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, apierr.InternalError("failed to check existing channel").WithDetails(err.Error())
	}

	if existing != nil {
		// 更新现有配置
		existing.MerchantNo = req.MerchantNo
		existing.MerchantKey = req.MerchantKey
		existing.CallbackURL = req.CallbackURL
		existing.Enabled = req.Enabled
		existing.Priority = req.Priority
		existing.Remark = req.Remark

		if err := s.repo.UpdateChannel(ctx, existing); err != nil {
			return nil, apierr.InternalError("failed to update payment channel").WithDetails(err.Error())
		}
		return existing, nil
	}

	// 创建新配置
	channel := &model.PaymentChannelConfig{
		CollectionEntityID: entityID,
		Channel:            req.Channel,
		MerchantNo:         req.MerchantNo,
		MerchantKey:        req.MerchantKey,
		CallbackURL:        req.CallbackURL,
		Enabled:            req.Enabled,
		Priority:           req.Priority,
		Remark:             req.Remark,
	}

	if err := s.repo.CreateChannel(ctx, channel); err != nil {
		return nil, apierr.InternalError("failed to create payment channel").WithDetails(err.Error())
	}

	return channel, nil
}

// GetPaymentChannel 获取支付渠道配置
func (s *CollectionEntityService) GetPaymentChannel(ctx context.Context, channelID uint64) (*model.PaymentChannelConfig, error) {
	channel, err := s.repo.GetChannel(ctx, channelID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrChannelNotFound
		}
		return nil, apierr.InternalError("failed to get payment channel").WithDetails(err.Error())
	}
	return channel, nil
}

// ListPaymentChannels 获取收款主体的所有支付渠道配置
// Requirements: 15.4
func (s *CollectionEntityService) ListPaymentChannels(ctx context.Context, entityID uint64) ([]model.PaymentChannelConfig, error) {
	// 验证收款主体存在
	_, err := s.repo.Get(ctx, entityID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("failed to get collection entity").WithDetails(err.Error())
	}

	channels, err := s.repo.ListChannelsByEntity(ctx, entityID)
	if err != nil {
		return nil, apierr.InternalError("failed to list payment channels").WithDetails(err.Error())
	}

	return channels, nil
}

// UpdatePaymentChannel 更新支付渠道配置
// Requirements: 15.4
func (s *CollectionEntityService) UpdatePaymentChannel(ctx context.Context, channelID uint64, req *model.ConfigurePaymentChannelRequest) (*model.PaymentChannelConfig, error) {
	channel, err := s.repo.GetChannel(ctx, channelID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrChannelNotFound
		}
		return nil, apierr.InternalError("failed to get payment channel").WithDetails(err.Error())
	}

	// 更新配置
	channel.MerchantNo = req.MerchantNo
	if req.MerchantKey != "" {
		channel.MerchantKey = req.MerchantKey
	}
	channel.CallbackURL = req.CallbackURL
	channel.Enabled = req.Enabled
	channel.Priority = req.Priority
	channel.Remark = req.Remark

	if err := s.repo.UpdateChannel(ctx, channel); err != nil {
		return nil, apierr.InternalError("failed to update payment channel").WithDetails(err.Error())
	}

	return channel, nil
}

// DeletePaymentChannel 删除支付渠道配置
func (s *CollectionEntityService) DeletePaymentChannel(ctx context.Context, channelID uint64) error {
	_, err := s.repo.GetChannel(ctx, channelID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrChannelNotFound
		}
		return apierr.InternalError("failed to get payment channel").WithDetails(err.Error())
	}

	if err := s.repo.DeleteChannel(ctx, channelID); err != nil {
		return apierr.InternalError("failed to delete payment channel").WithDetails(err.Error())
	}

	return nil
}

// TogglePaymentChannelStatus 切换支付渠道状态
// Requirements: 15.4
func (s *CollectionEntityService) TogglePaymentChannelStatus(ctx context.Context, channelID uint64, enabled bool) error {
	channel, err := s.repo.GetChannel(ctx, channelID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrChannelNotFound
		}
		return apierr.InternalError("failed to get payment channel").WithDetails(err.Error())
	}

	channel.Enabled = enabled
	if err := s.repo.UpdateChannel(ctx, channel); err != nil {
		return apierr.InternalError("failed to toggle payment channel status").WithDetails(err.Error())
	}

	return nil
}

// GetChannelByEntityAndMethod 根据收款主体和支付方式获取渠道配置
// Requirements: 17.2
func (s *CollectionEntityService) GetChannelByEntityAndMethod(ctx context.Context, entityID uint64, method model.PaymentMethod) (*model.PaymentChannelConfig, error) {
	channel, err := s.repo.GetChannelByEntityAndMethod(ctx, entityID, method)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrChannelNotFound.WithDetails(fmt.Sprintf("no enabled channel for method %s", method))
		}
		return nil, apierr.InternalError("failed to get payment channel").WithDetails(err.Error())
	}
	return channel, nil
}

// UpdateCollectionStats 更新收款统计
// Requirements: 18.1
func (s *CollectionEntityService) UpdateCollectionStats(ctx context.Context, entityID uint64, amountCents int64) error {
	if err := s.repo.UpdateCollectionStats(ctx, entityID, amountCents); err != nil {
		return apierr.InternalError("failed to update collection stats").WithDetails(err.Error())
	}
	return nil
}

// GetCollectionStats 获取收款统计
// Requirements: 18.1
func (s *CollectionEntityService) GetCollectionStats(ctx context.Context, entityID uint64) (totalCents int64, count int64, err error) {
	totalCents, count, err = s.repo.GetCollectionStats(ctx, entityID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return 0, 0, ErrNotFound
		}
		return 0, 0, apierr.InternalError("failed to get collection stats").WithDetails(err.Error())
	}
	return totalCents, count, nil
}
