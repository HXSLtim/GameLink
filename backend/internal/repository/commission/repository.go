package commission

import (
	"context"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"gorm.io/gorm"
)

// CommissionRepository 抽成记录仓储接口
type CommissionRepository interface {
	// 抽成规则
	CreateRule(ctx context.Context, rule *model.CommissionRule) error
	GetRule(ctx context.Context, id uint64) (*model.CommissionRule, error)
	GetDefaultRule(ctx context.Context) (*model.CommissionRule, error)
	GetRuleForOrder(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error)
	ListRules(ctx context.Context, opts CommissionRuleListOptions) ([]model.CommissionRule, int64, error)
	UpdateRule(ctx context.Context, rule *model.CommissionRule) error
	DeleteRule(ctx context.Context, id uint64) error

	// 抽成记录
	CreateRecord(ctx context.Context, record *model.CommissionRecord) error
	GetRecord(ctx context.Context, id uint64) (*model.CommissionRecord, error)
	GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error)
	ListRecords(ctx context.Context, opts CommissionRecordListOptions) ([]model.CommissionRecord, int64, error)
	UpdateRecord(ctx context.Context, record *model.CommissionRecord) error

	// 月度结算
	CreateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error
	GetSettlement(ctx context.Context, id uint64) (*model.MonthlySettlement, error)
	GetSettlementByPlayerMonth(ctx context.Context, playerID uint64, month string) (*model.MonthlySettlement, error)
	ListSettlements(ctx context.Context, opts SettlementListOptions) ([]model.MonthlySettlement, int64, error)
	UpdateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error

	// 统计查询
	GetMonthlyStats(ctx context.Context, month string) (*MonthlyStats, error)
	GetPlayerMonthlyIncome(ctx context.Context, playerID uint64, month string) (int64, error)
}

// CommissionRuleListOptions 抽成规则查询选项
type CommissionRuleListOptions struct {
	Type     *string
	GameID   *uint64
	PlayerID *uint64
	IsActive *bool
	Page     int
	PageSize int
}

// CommissionRecordListOptions 抽成记录查询选项
type CommissionRecordListOptions struct {
	OrderID          *uint64
	PlayerID         *uint64
	SettlementStatus *string
	SettlementMonth  *string
	DateFrom         *time.Time
	DateTo           *time.Time
	Page             int
	PageSize         int
}

// SettlementListOptions 月度结算查询选项
type SettlementListOptions struct {
	PlayerID        *uint64
	SettlementMonth *string
	Status          *string
	Page            int
	PageSize        int
}

// MonthlyStats 月度统计数据
type MonthlyStats struct {
	TotalOrders       int64
	TotalIncome       int64
	TotalCommission   int64
	TotalPlayerIncome int64
}

type commissionRepository struct {
	db *gorm.DB
}

// NewCommissionRepository 创建抽成记录仓储
func NewCommissionRepository(db *gorm.DB) CommissionRepository {
	return &commissionRepository{db: db}
}

// CreateRule 创建抽成规则
func (r *commissionRepository) CreateRule(ctx context.Context, rule *model.CommissionRule) error {
	return r.db.WithContext(ctx).Create(rule).Error
}

// GetRule 获取抽成规则
func (r *commissionRepository) GetRule(ctx context.Context, id uint64) (*model.CommissionRule, error) {
	var rule model.CommissionRule
	err := r.db.WithContext(ctx).First(&rule, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &rule, nil
}

// GetDefaultRule 获取默认抽成规则
func (r *commissionRepository) GetDefaultRule(ctx context.Context) (*model.CommissionRule, error) {
	var rule model.CommissionRule
	err := r.db.WithContext(ctx).
		Where("type = ? AND is_active = ?", "default", true).
		Where("game_id IS NULL AND player_id IS NULL AND service_type IS NULL").
		First(&rule).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &rule, nil
}

// GetRuleForOrder 获取订单适用的抽成规则（优先级：特定玩家 > 特定游戏 > 特定服务类型 > 默认�?
func (r *commissionRepository) GetRuleForOrder(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
	var rule model.CommissionRule

	// 1. 尝试获取玩家专属规则
	if playerID != nil {
		err := r.db.WithContext(ctx).
			Where("player_id = ? AND is_active = ?", *playerID, true).
			First(&rule).Error
		if err == nil {
			return &rule, nil
		}
	}

	// 2. 尝试获取游戏专属规则
	if gameID != nil {
		err := r.db.WithContext(ctx).
			Where("game_id = ? AND is_active = ?", *gameID, true).
			Where("player_id IS NULL").
			First(&rule).Error
		if err == nil {
			return &rule, nil
		}
	}

	// 3. 尝试获取服务类型规则
	if serviceType != nil {
		err := r.db.WithContext(ctx).
			Where("service_type = ? AND is_active = ?", *serviceType, true).
			Where("player_id IS NULL AND game_id IS NULL").
			First(&rule).Error
		if err == nil {
			return &rule, nil
		}
	}

	// 4. 返回默认规则
	return r.GetDefaultRule(ctx)
}

// ListRules 查询抽成规则列表
func (r *commissionRepository) ListRules(ctx context.Context, opts CommissionRuleListOptions) ([]model.CommissionRule, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.CommissionRule{})

	// 过滤条件
	if opts.Type != nil {
		query = query.Where("type = ?", *opts.Type)
	}
	if opts.GameID != nil {
		query = query.Where("game_id = ?", *opts.GameID)
	}
	if opts.PlayerID != nil {
		query = query.Where("player_id = ?", *opts.PlayerID)
	}
	if opts.IsActive != nil {
		query = query.Where("is_active = ?", *opts.IsActive)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
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
	var rules []model.CommissionRule
	err := query.Order("created_at DESC").Offset(offset).Limit(opts.PageSize).Find(&rules).Error
	if err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

// UpdateRule 更新抽成规则
func (r *commissionRepository) UpdateRule(ctx context.Context, rule *model.CommissionRule) error {
	return r.db.WithContext(ctx).Save(rule).Error
}

// DeleteRule 删除抽成规则
func (r *commissionRepository) DeleteRule(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.CommissionRule{}, id).Error
}

// CreateRecord 创建抽成记录
func (r *commissionRepository) CreateRecord(ctx context.Context, record *model.CommissionRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

// GetRecord 获取抽成记录
func (r *commissionRepository) GetRecord(ctx context.Context, id uint64) (*model.CommissionRecord, error) {
	var record model.CommissionRecord
	err := r.db.WithContext(ctx).First(&record, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &record, nil
}

// GetRecordByOrderID 根据订单ID获取抽成记录
func (r *commissionRepository) GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
	var record model.CommissionRecord
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &record, nil
}

// ListRecords 查询抽成记录列表
func (r *commissionRepository) ListRecords(ctx context.Context, opts CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.CommissionRecord{})

	// 过滤条件
	if opts.OrderID != nil {
		query = query.Where("order_id = ?", *opts.OrderID)
	}
	if opts.PlayerID != nil {
		query = query.Where("player_id = ?", *opts.PlayerID)
	}
	if opts.SettlementStatus != nil {
		query = query.Where("settlement_status = ?", *opts.SettlementStatus)
	}
	if opts.SettlementMonth != nil {
		query = query.Where("settlement_month = ?", *opts.SettlementMonth)
	}
	if opts.DateFrom != nil {
		query = query.Where("created_at >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		query = query.Where("created_at < ?", *opts.DateTo)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
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
	var records []model.CommissionRecord
	err := query.Order("created_at DESC").Offset(offset).Limit(opts.PageSize).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// UpdateRecord 更新抽成记录
func (r *commissionRepository) UpdateRecord(ctx context.Context, record *model.CommissionRecord) error {
	return r.db.WithContext(ctx).Save(record).Error
}

// CreateSettlement 创建月度结算
func (r *commissionRepository) CreateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	return r.db.WithContext(ctx).Create(settlement).Error
}

// GetSettlement 获取月度结算
func (r *commissionRepository) GetSettlement(ctx context.Context, id uint64) (*model.MonthlySettlement, error) {
	var settlement model.MonthlySettlement
	err := r.db.WithContext(ctx).First(&settlement, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &settlement, nil
}

// GetSettlementByPlayerMonth 根据玩家和月份获取结�?
func (r *commissionRepository) GetSettlementByPlayerMonth(ctx context.Context, playerID uint64, month string) (*model.MonthlySettlement, error) {
	var settlement model.MonthlySettlement
	err := r.db.WithContext(ctx).
		Where("player_id = ? AND settlement_month = ?", playerID, month).
		First(&settlement).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &settlement, nil
}

// ListSettlements 查询月度结算列表
func (r *commissionRepository) ListSettlements(ctx context.Context, opts SettlementListOptions) ([]model.MonthlySettlement, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.MonthlySettlement{})

	// 过滤条件
	if opts.PlayerID != nil {
		query = query.Where("player_id = ?", *opts.PlayerID)
	}
	if opts.SettlementMonth != nil {
		query = query.Where("settlement_month = ?", *opts.SettlementMonth)
	}
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
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
	var settlements []model.MonthlySettlement
	err := query.Order("settlement_month DESC, created_at DESC").
		Offset(offset).Limit(opts.PageSize).Find(&settlements).Error
	if err != nil {
		return nil, 0, err
	}

	return settlements, total, nil
}

// UpdateSettlement 更新月度结算
func (r *commissionRepository) UpdateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	return r.db.WithContext(ctx).Save(settlement).Error
}

// GetMonthlyStats 获取月度统计数据
func (r *commissionRepository) GetMonthlyStats(ctx context.Context, month string) (*MonthlyStats, error) {
	stats := &MonthlyStats{}

	// 统计该月所有已结算记录
	var result struct {
		TotalOrders       int64
		TotalIncome       int64
		TotalCommission   int64
		TotalPlayerIncome int64
	}

	err := r.db.WithContext(ctx).
		Model(&model.CommissionRecord{}).
		Where("settlement_month = ? AND settlement_status = ?", month, "settled").
		Select(`
			COUNT(*) as total_orders,
			SUM(total_amount_cents) as total_income,
			SUM(commission_cents) as total_commission,
			SUM(player_income_cents) as total_player_income
		`).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	stats.TotalOrders = result.TotalOrders
	stats.TotalIncome = result.TotalIncome
	stats.TotalCommission = result.TotalCommission
	stats.TotalPlayerIncome = result.TotalPlayerIncome

	return stats, nil
}

// GetPlayerMonthlyIncome 获取玩家月度收入
func (r *commissionRepository) GetPlayerMonthlyIncome(ctx context.Context, playerID uint64, month string) (int64, error) {
	var income int64
	err := r.db.WithContext(ctx).
		Model(&model.CommissionRecord{}).
		Where("player_id = ? AND settlement_month = ?", playerID, month).
		Select("COALESCE(SUM(player_income_cents), 0)").
		Scan(&income).Error
	if err != nil {
		return 0, err
	}
	return income, nil
}

