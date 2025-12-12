package settlementcompany

import (
	"context"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"gorm.io/gorm"
)

// SettlementCompanyRepository 结算公司仓储接口
// Requirements: 11.3, 12.5
type SettlementCompanyRepository interface {
	// Company CRUD operations
	Create(ctx context.Context, company *model.SettlementCompany) error
	Get(ctx context.Context, id uint64) (*model.SettlementCompany, error)
	GetByCreditCode(ctx context.Context, creditCode string) (*model.SettlementCompany, error)
	Update(ctx context.Context, company *model.SettlementCompany) error
	List(ctx context.Context, opts ListOptions) ([]model.SettlementCompany, int64, error)
	ToggleStatus(ctx context.Context, id uint64, status model.CompanyStatus) error

	// Player assignment operations
	AssignPlayer(ctx context.Context, assignment *model.PlayerCompanyAssignment) error
	GetCurrentAssignment(ctx context.Context, playerID uint64) (*model.PlayerCompanyAssignment, error)
	GetAssignmentHistory(ctx context.Context, playerID uint64) ([]model.PlayerCompanyAssignment, error)
	EndCurrentAssignment(ctx context.Context, playerID uint64, endDate time.Time) error
	BatchAssignPlayers(ctx context.Context, assignments []model.PlayerCompanyAssignment) error

	// History operations
	CreateHistory(ctx context.Context, history *model.SettlementCompanyHistory) error
	GetHistory(ctx context.Context, companyID uint64) ([]model.SettlementCompanyHistory, error)

	// Statistics
	GetPlayerCount(ctx context.Context, companyID uint64) (int, error)
	UpdatePlayerCount(ctx context.Context, companyID uint64) error
}

// ListOptions 结算公司查询选项
type ListOptions struct {
	Status    *model.CompanyStatus
	Keyword   string
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
}

type settlementCompanyRepository struct {
	db *gorm.DB
}

// NewSettlementCompanyRepository 创建结算公司仓储
func NewSettlementCompanyRepository(db *gorm.DB) SettlementCompanyRepository {
	return &settlementCompanyRepository{db: db}
}

// Create 创建结算公司
func (r *settlementCompanyRepository) Create(ctx context.Context, company *model.SettlementCompany) error {
	return r.db.WithContext(ctx).Create(company).Error
}

// Get 获取结算公司
func (r *settlementCompanyRepository) Get(ctx context.Context, id uint64) (*model.SettlementCompany, error) {
	var company model.SettlementCompany
	err := r.db.WithContext(ctx).First(&company, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &company, nil
}

// GetByCreditCode 根据统一社会信用代码获取结算公司
func (r *settlementCompanyRepository) GetByCreditCode(ctx context.Context, creditCode string) (*model.SettlementCompany, error) {
	var company model.SettlementCompany
	err := r.db.WithContext(ctx).Where("credit_code = ?", creditCode).First(&company).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &company, nil
}

// Update 更新结算公司
func (r *settlementCompanyRepository) Update(ctx context.Context, company *model.SettlementCompany) error {
	return r.db.WithContext(ctx).Save(company).Error
}

// List 查询结算公司列表
// Requirements: 11.3
func (r *settlementCompanyRepository) List(ctx context.Context, opts ListOptions) ([]model.SettlementCompany, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.SettlementCompany{})

	// 过滤条件
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}
	if opts.Keyword != "" {
		keyword := "%" + opts.Keyword + "%"
		query = query.Where("name LIKE ? OR credit_code LIKE ?", keyword, keyword)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	orderClause := "created_at DESC"
	if opts.SortBy != "" {
		order := "ASC"
		if opts.SortOrder == "desc" {
			order = "DESC"
		}
		switch opts.SortBy {
		case "name":
			orderClause = "name " + order
		case "created_at":
			orderClause = "created_at " + order
		case "player_count":
			orderClause = "player_count " + order
		}
	}

	// 分页
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 {
		opts.PageSize = 20
	}
	offset := (opts.Page - 1) * opts.PageSize

	// 查询数据
	var companies []model.SettlementCompany
	err := query.Order(orderClause).Offset(offset).Limit(opts.PageSize).Find(&companies).Error
	if err != nil {
		return nil, 0, err
	}

	return companies, total, nil
}

// ToggleStatus 切换结算公司状态
// Requirements: 11.4
func (r *settlementCompanyRepository) ToggleStatus(ctx context.Context, id uint64, status model.CompanyStatus) error {
	return r.db.WithContext(ctx).Model(&model.SettlementCompany{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// AssignPlayer 分配陪玩师到结算公司
// Requirements: 12.1
func (r *settlementCompanyRepository) AssignPlayer(ctx context.Context, assignment *model.PlayerCompanyAssignment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 结束当前分配
		if err := tx.Model(&model.PlayerCompanyAssignment{}).
			Where("player_id = ? AND is_current = ?", assignment.PlayerID, true).
			Updates(map[string]interface{}{
				"is_current": false,
				"end_date":   assignment.EffectiveDate,
			}).Error; err != nil {
			return err
		}

		// 创建新分配
		assignment.IsCurrent = true
		if err := tx.Create(assignment).Error; err != nil {
			return err
		}

		// 更新结算公司的陪玩师数量
		return r.updatePlayerCountInTx(tx, assignment.SettlementCompanyID)
	})
}

// GetCurrentAssignment 获取陪玩师当前的结算公司分配
// Requirements: 12.5
func (r *settlementCompanyRepository) GetCurrentAssignment(ctx context.Context, playerID uint64) (*model.PlayerCompanyAssignment, error) {
	var assignment model.PlayerCompanyAssignment
	err := r.db.WithContext(ctx).
		Where("player_id = ? AND is_current = ?", playerID, true).
		Preload("SettlementCompany").
		First(&assignment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &assignment, nil
}

// GetAssignmentHistory 获取陪玩师的结算公司分配历史
// Requirements: 12.5
func (r *settlementCompanyRepository) GetAssignmentHistory(ctx context.Context, playerID uint64) ([]model.PlayerCompanyAssignment, error) {
	var assignments []model.PlayerCompanyAssignment
	err := r.db.WithContext(ctx).
		Where("player_id = ?", playerID).
		Preload("SettlementCompany").
		Order("effective_date DESC").
		Find(&assignments).Error
	if err != nil {
		return nil, err
	}
	return assignments, nil
}

// EndCurrentAssignment 结束当前分配
func (r *settlementCompanyRepository) EndCurrentAssignment(ctx context.Context, playerID uint64, endDate time.Time) error {
	return r.db.WithContext(ctx).Model(&model.PlayerCompanyAssignment{}).
		Where("player_id = ? AND is_current = ?", playerID, true).
		Updates(map[string]interface{}{
			"is_current": false,
			"end_date":   endDate,
		}).Error
}

// BatchAssignPlayers 批量分配陪玩师到结算公司
// Requirements: 12.3
func (r *settlementCompanyRepository) BatchAssignPlayers(ctx context.Context, assignments []model.PlayerCompanyAssignment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		companyIDs := make(map[uint64]bool)

		for i := range assignments {
			// 结束当前分配
			if err := tx.Model(&model.PlayerCompanyAssignment{}).
				Where("player_id = ? AND is_current = ?", assignments[i].PlayerID, true).
				Updates(map[string]interface{}{
					"is_current": false,
					"end_date":   assignments[i].EffectiveDate,
				}).Error; err != nil {
				return err
			}

			// 创建新分配
			assignments[i].IsCurrent = true
			if err := tx.Create(&assignments[i]).Error; err != nil {
				return err
			}

			companyIDs[assignments[i].SettlementCompanyID] = true
		}

		// 更新所有涉及的结算公司的陪玩师数量
		for companyID := range companyIDs {
			if err := r.updatePlayerCountInTx(tx, companyID); err != nil {
				return err
			}
		}

		return nil
	})
}

// CreateHistory 创建结算公司修改历史
// Requirements: 11.5
func (r *settlementCompanyRepository) CreateHistory(ctx context.Context, history *model.SettlementCompanyHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// GetHistory 获取结算公司修改历史
func (r *settlementCompanyRepository) GetHistory(ctx context.Context, companyID uint64) ([]model.SettlementCompanyHistory, error) {
	var histories []model.SettlementCompanyHistory
	err := r.db.WithContext(ctx).
		Where("settlement_company_id = ?", companyID).
		Order("created_at DESC").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

// GetPlayerCount 获取结算公司的陪玩师数量
func (r *settlementCompanyRepository) GetPlayerCount(ctx context.Context, companyID uint64) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.PlayerCompanyAssignment{}).
		Where("settlement_company_id = ? AND is_current = ?", companyID, true).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// UpdatePlayerCount 更新结算公司的陪玩师数量
func (r *settlementCompanyRepository) UpdatePlayerCount(ctx context.Context, companyID uint64) error {
	return r.updatePlayerCountInTx(r.db.WithContext(ctx), companyID)
}

// updatePlayerCountInTx 在事务中更新陪玩师数量
func (r *settlementCompanyRepository) updatePlayerCountInTx(tx *gorm.DB, companyID uint64) error {
	var count int64
	if err := tx.Model(&model.PlayerCompanyAssignment{}).
		Where("settlement_company_id = ? AND is_current = ?", companyID, true).
		Count(&count).Error; err != nil {
		return err
	}

	return tx.Model(&model.SettlementCompany{}).
		Where("id = ?", companyID).
		Update("player_count", count).Error
}
