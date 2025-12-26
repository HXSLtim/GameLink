package playercertification

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// gormPlayerCertificationRepository 陪玩师实名认证仓储实现
type gormPlayerCertificationRepository struct {
	db *gorm.DB
}

// NewPlayerCertificationRepository 创建陪玩师实名认证仓储
func NewPlayerCertificationRepository(db *gorm.DB) repository.PlayerCertificationRepository {
	return &gormPlayerCertificationRepository{db: db}
}

// Create 创建实名认证记录
func (r *gormPlayerCertificationRepository) Create(ctx context.Context, cert *model.PlayerCertification) error {
	return r.db.WithContext(ctx).Create(cert).Error
}

// Get 根据ID获取实名认证记录
func (r *gormPlayerCertificationRepository) Get(ctx context.Context, id uint64) (*model.PlayerCertification, error) {
	var cert model.PlayerCertification
	if err := r.db.WithContext(ctx).First(&cert, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &cert, nil
}

// GetWithPlayer 获取实名认证记录及关联的陪玩师信息
func (r *gormPlayerCertificationRepository) GetWithPlayer(ctx context.Context, id uint64) (*model.PlayerCertification, error) {
	var cert model.PlayerCertification
	if err := r.db.WithContext(ctx).
		Preload("Player").
		Preload("Verifier").
		First(&cert, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &cert, nil
}

// GetByPlayerID 根据陪玩师ID获取实名认证记录
func (r *gormPlayerCertificationRepository) GetByPlayerID(ctx context.Context, playerID uint64) (*model.PlayerCertification, error) {
	var cert model.PlayerCertification
	if err := r.db.WithContext(ctx).
		Where("player_id = ?", playerID).
		First(&cert).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &cert, nil
}

// ListPaged 分页获取实名认证记录
func (r *gormPlayerCertificationRepository) ListPaged(ctx context.Context, opts repository.PlayerCertificationListOptions) ([]model.PlayerCertification, int64, error) {
	page := repository.NormalizePage(opts.Page)
	pageSize := repository.NormalizePageSize(opts.PageSize)
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.PlayerCertification{})

	// 陪玩师ID筛选
	if opts.PlayerID != nil {
		query = query.Where("player_id = ?", *opts.PlayerID)
	}

	// 状态筛选
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}

	// 多状态筛选
	if len(opts.Statuses) > 0 {
		query = query.Where("status IN ?", opts.Statuses)
	}

	// 关键词搜索（真实姓名）
	if opts.Keyword != "" {
		likePattern := "%" + opts.Keyword + "%"
		query = query.Where("real_name ILIKE ?", likePattern)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var certs []model.PlayerCertification
	if err := query.
		Preload("Player").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&certs).Error; err != nil {
		return nil, 0, err
	}

	return certs, total, nil
}

// ListPending 获取待审核的实名认证记录
func (r *gormPlayerCertificationRepository) ListPending(ctx context.Context, page, pageSize int) ([]model.PlayerCertification, int64, error) {
	return r.ListPaged(ctx, repository.PlayerCertificationListOptions{
		Page:     page,
		PageSize: pageSize,
		Status:   ptrCertificationStatus(model.CertificationStatusPending),
	})
}

// Update 更新实名认证记录
func (r *gormPlayerCertificationRepository) Update(ctx context.Context, cert *model.PlayerCertification) error {
	tx := r.db.WithContext(ctx).Model(cert).Updates(map[string]any{
		"real_name":         cert.RealName,
		"id_card_no":        cert.IDCardNo,
		"id_card_front_url": cert.IDCardFrontURL,
		"id_card_back_url":  cert.IDCardBackURL,
		"status":            cert.Status,
		"verified_at":       cert.VerifiedAt,
		"verified_by":       cert.VerifiedBy,
		"reject_reason":     cert.RejectReason,
		"photo_url":         cert.PhotoURL,
		"voice_url":         cert.VoiceURL,
		"ext_json":          cert.ExtJSON,
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// UpdateStatus 更新认证状态
func (r *gormPlayerCertificationRepository) UpdateStatus(ctx context.Context, id uint64, status model.CertificationStatus, verifiedBy *uint64, rejectReason string) error {
	updates := map[string]any{
		"status":        status,
		"reject_reason": rejectReason,
	}
	if verifiedBy != nil {
		updates["verified_by"] = *verifiedBy
	}
	if status == model.CertificationStatusVerified {
		updates["verified_at"] = gorm.Expr("NOW()")
	}

	tx := r.db.WithContext(ctx).Model(&model.PlayerCertification{}).Where("id = ?", id).Updates(updates)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// Delete 删除实名认证记录
func (r *gormPlayerCertificationRepository) Delete(ctx context.Context, id uint64) error {
	tx := r.db.WithContext(ctx).Delete(&model.PlayerCertification{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// CountByStatus 统计各状态的认证数量
func (r *gormPlayerCertificationRepository) CountByStatus(ctx context.Context) (map[model.CertificationStatus]int64, error) {
	type result struct {
		Status model.CertificationStatus
		Count  int64
	}
	var results []result
	if err := r.db.WithContext(ctx).
		Model(&model.PlayerCertification{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[model.CertificationStatus]int64)
	for _, r := range results {
		counts[r.Status] = r.Count
	}
	return counts, nil
}

// GetPendingCount 获取待审核数量
func (r *gormPlayerCertificationRepository) GetPendingCount(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.PlayerCertification{}).
		Where("status = ?", model.CertificationStatusPending).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ptrCertificationStatus 返回状态指针
func ptrCertificationStatus(s model.CertificationStatus) *model.CertificationStatus {
	return &s
}
