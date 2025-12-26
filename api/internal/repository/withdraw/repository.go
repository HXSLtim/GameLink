package withdraw

import (
	"context"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"gorm.io/gorm"
)

// WithdrawRepository 提现记录仓储接口
type WithdrawRepository interface {
	// Create 创建提现记录
	Create(ctx context.Context, withdraw *model.Withdraw) error
	// Get 获取提现记录
	Get(ctx context.Context, id uint64) (*model.Withdraw, error)
	// Update 更新提现记录
	Update(ctx context.Context, withdraw *model.Withdraw) error
	// List 查询提现记录列表
	List(ctx context.Context, opts WithdrawListOptions) ([]model.Withdraw, int64, error)
	// GetPlayerBalance 获取陪玩师余额信息
	GetPlayerBalance(ctx context.Context, playerID uint64) (*PlayerBalance, error)

	// 提现分流相关方法 - Requirements: 14.1-14.5
	// ListByCompany 按结算公司查询提现列表
	ListByCompany(ctx context.Context, opts WithdrawByCompanyOptions) ([]model.Withdraw, int64, error)
	// GetRoutingStats 获取提现分流统计
	GetRoutingStats(ctx context.Context, dateFrom, dateTo *time.Time) (*model.WithdrawRoutingStatsResponse, error)
	// GetRoutingStatsByCompany 获取单个结算公司的提现统计
	GetRoutingStatsByCompany(ctx context.Context, companyID uint64, dateFrom, dateTo *time.Time) (*model.WithdrawRoutingStats, error)
	// CreateSalaryPaymentRecord 创建工资发放记录
	CreateSalaryPaymentRecord(ctx context.Context, record *model.SalaryPaymentRecord) error
	// GetSalaryPaymentRecordByWithdrawID 根据提现ID获取工资发放记录
	GetSalaryPaymentRecordByWithdrawID(ctx context.Context, withdrawID uint64) (*model.SalaryPaymentRecord, error)
	// UpdateSalaryPaymentRecord 更新工资发放记录
	UpdateSalaryPaymentRecord(ctx context.Context, record *model.SalaryPaymentRecord) error
}

// WithdrawListOptions 提现记录查询选项
type WithdrawListOptions struct {
	PlayerID *uint64
	UserID   *uint64
	Status   *model.WithdrawStatus
	DateFrom *time.Time
	DateTo   *time.Time
	Page     int
	PageSize int
}

// WithdrawByCompanyOptions 按结算公司查询提现选项
type WithdrawByCompanyOptions struct {
	SettlementCompanyID *uint64
	Status              *model.WithdrawStatus
	DateFrom            *time.Time
	DateTo              *time.Time
	Page                int
	PageSize            int
}

// PlayerBalance 陪玩师余额信息
type PlayerBalance struct {
	TotalEarnings    int64 // 累计收益
	WithdrawTotal    int64 // 累计提现
	PendingWithdraw  int64 // 待处理提现
	AvailableBalance int64 // 可提现余额
	PendingBalance   int64 // 待结算余额
}

type withdrawRepository struct {
	db *gorm.DB
}

// NewWithdrawRepository 创建提现记录仓储
func NewWithdrawRepository(db *gorm.DB) WithdrawRepository {
	return &withdrawRepository{db: db}
}

// Create 创建提现记录
func (r *withdrawRepository) Create(ctx context.Context, withdraw *model.Withdraw) error {
	return r.db.WithContext(ctx).Create(withdraw).Error
}

// Get 获取提现记录
func (r *withdrawRepository) Get(ctx context.Context, id uint64) (*model.Withdraw, error) {
	var withdraw model.Withdraw
	err := r.db.WithContext(ctx).First(&withdraw, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &withdraw, nil
}

// Update 更新提现记录
func (r *withdrawRepository) Update(ctx context.Context, withdraw *model.Withdraw) error {
	return r.db.WithContext(ctx).Save(withdraw).Error
}

// List 查询提现记录列表
func (r *withdrawRepository) List(ctx context.Context, opts WithdrawListOptions) ([]model.Withdraw, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Withdraw{})

	// 过滤条件
	if opts.PlayerID != nil {
		query = query.Where("player_id = ?", *opts.PlayerID)
	}
	if opts.UserID != nil {
		query = query.Where("user_id = ?", *opts.UserID)
	}
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
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
	var withdraws []model.Withdraw
	err := query.Order("created_at DESC").Offset(offset).Limit(opts.PageSize).Find(&withdraws).Error
	if err != nil {
		return nil, 0, err
	}

	return withdraws, total, nil
}

// GetPlayerBalance 获取陪玩师余额信息
func (r *withdrawRepository) GetPlayerBalance(ctx context.Context, playerID uint64) (*PlayerBalance, error) {
	balance := &PlayerBalance{}

	// 计算累计收益（从已完成订单）
	var totalEarnings int64
	err := r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("player_id = ? AND status = ?", playerID, model.OrderStatusCompleted).
		Select("COALESCE(SUM(total_price_cents), 0)").
		Scan(&totalEarnings).Error
	if err != nil {
		return nil, err
	}
	balance.TotalEarnings = totalEarnings

	// 计算累计提现（已完成的提现）
	var withdrawTotal int64
	err = r.db.WithContext(ctx).
		Model(&model.Withdraw{}).
		Where("player_id = ? AND status = ?", playerID, model.WithdrawStatusCompleted).
		Select("COALESCE(SUM(amount_cents), 0)").
		Scan(&withdrawTotal).Error
	if err != nil {
		return nil, err
	}
	balance.WithdrawTotal = withdrawTotal

	// 计算待处理提现（pending approved 状态）
	var pendingWithdraw int64
	err = r.db.WithContext(ctx).
		Model(&model.Withdraw{}).
		Where("player_id = ? AND status IN ?", playerID, []model.WithdrawStatus{
			model.WithdrawStatusPending,
			model.WithdrawStatusApproved,
		}).
		Select("COALESCE(SUM(amount_cents), 0)").
		Scan(&pendingWithdraw).Error
	if err != nil {
		return nil, err
	}
	balance.PendingWithdraw = pendingWithdraw

	// 计算待结算余额（进行中的订单）
	var pendingBalance int64
	err = r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("player_id = ? AND status = ?", playerID, model.OrderStatusInProgress).
		Select("COALESCE(SUM(total_price_cents), 0)").
		Scan(&pendingBalance).Error
	if err != nil {
		return nil, err
	}
	balance.PendingBalance = pendingBalance

	// 计算可提现余额 = 累计收益 - 累计提现 - 待处理提现 - 待结算余额
	balance.AvailableBalance = totalEarnings - withdrawTotal - pendingWithdraw - pendingBalance
	if balance.AvailableBalance < 0 {
		balance.AvailableBalance = 0
	}

	return balance, nil
}

// ListByCompany 按结算公司查询提现列表
// Requirements: 14.5
func (r *withdrawRepository) ListByCompany(ctx context.Context, opts WithdrawByCompanyOptions) ([]model.Withdraw, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Withdraw{})

	// 过滤条件
	if opts.SettlementCompanyID != nil {
		query = query.Where("settlement_company_id = ?", *opts.SettlementCompanyID)
	}
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
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
	var withdraws []model.Withdraw
	err := query.Order("created_at DESC").Offset(offset).Limit(opts.PageSize).Find(&withdraws).Error
	if err != nil {
		return nil, 0, err
	}

	return withdraws, total, nil
}

// GetRoutingStats 获取提现分流统计
// Requirements: 14.1
// Property 14: 提现分流统计准确性
func (r *withdrawRepository) GetRoutingStats(ctx context.Context, dateFrom, dateTo *time.Time) (*model.WithdrawRoutingStatsResponse, error) {
	query := r.db.WithContext(ctx).Model(&model.Withdraw{}).
		Where("status = ?", model.WithdrawStatusCompleted)

	if dateFrom != nil {
		query = query.Where("completed_at >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("completed_at < ?", *dateTo)
	}

	// 获取总体统计
	var totalStats struct {
		TotalWithdrawals       int64
		TotalAmountCents       int64
		TotalTaxDeductedCents  int64
		TotalActualAmountCents int64
	}
	err := query.Select(`
		COUNT(*) as total_withdrawals,
		COALESCE(SUM(amount_cents), 0) as total_amount_cents,
		COALESCE(SUM(tax_deducted_cents), 0) as total_tax_deducted_cents,
		COALESCE(SUM(actual_amount_cents), 0) as total_actual_amount_cents
	`).Scan(&totalStats).Error
	if err != nil {
		return nil, err
	}

	// 按结算公司分组统计
	type companyStats struct {
		SettlementCompanyID    uint64
		SettlementCompanyName  string
		TotalWithdrawals       int64
		TotalAmountCents       int64
		TotalTaxDeductedCents  int64
		TotalActualAmountCents int64
	}

	var byCompanyStats []companyStats
	groupQuery := r.db.WithContext(ctx).Model(&model.Withdraw{}).
		Where("status = ? AND settlement_company_id IS NOT NULL", model.WithdrawStatusCompleted)

	if dateFrom != nil {
		groupQuery = groupQuery.Where("completed_at >= ?", *dateFrom)
	}
	if dateTo != nil {
		groupQuery = groupQuery.Where("completed_at < ?", *dateTo)
	}

	err = groupQuery.Select(`
		settlement_company_id,
		settlement_company_name,
		COUNT(*) as total_withdrawals,
		COALESCE(SUM(amount_cents), 0) as total_amount_cents,
		COALESCE(SUM(tax_deducted_cents), 0) as total_tax_deducted_cents,
		COALESCE(SUM(actual_amount_cents), 0) as total_actual_amount_cents
	`).Group("settlement_company_id, settlement_company_name").Scan(&byCompanyStats).Error
	if err != nil {
		return nil, err
	}

	// 构建响应
	response := &model.WithdrawRoutingStatsResponse{
		TotalWithdrawals:       totalStats.TotalWithdrawals,
		TotalAmountCents:       totalStats.TotalAmountCents,
		TotalTaxDeductedCents:  totalStats.TotalTaxDeductedCents,
		TotalActualAmountCents: totalStats.TotalActualAmountCents,
		ByCompany:              make([]model.WithdrawRoutingStats, len(byCompanyStats)),
	}

	for i, cs := range byCompanyStats {
		var percentage float64
		if totalStats.TotalAmountCents > 0 {
			percentage = float64(cs.TotalAmountCents) / float64(totalStats.TotalAmountCents) * 100
		}
		var avgAmount int64
		if cs.TotalWithdrawals > 0 {
			avgAmount = cs.TotalAmountCents / cs.TotalWithdrawals
		}

		response.ByCompany[i] = model.WithdrawRoutingStats{
			SettlementCompanyID:    cs.SettlementCompanyID,
			SettlementCompanyName:  cs.SettlementCompanyName,
			TotalWithdrawals:       cs.TotalWithdrawals,
			TotalAmountCents:       cs.TotalAmountCents,
			TotalTaxDeductedCents:  cs.TotalTaxDeductedCents,
			TotalActualAmountCents: cs.TotalActualAmountCents,
			AverageAmountCents:     avgAmount,
			Percentage:             percentage,
		}
	}

	return response, nil
}

// GetRoutingStatsByCompany 获取单个结算公司的提现统计
// Requirements: 14.5
func (r *withdrawRepository) GetRoutingStatsByCompany(ctx context.Context, companyID uint64, dateFrom, dateTo *time.Time) (*model.WithdrawRoutingStats, error) {
	query := r.db.WithContext(ctx).Model(&model.Withdraw{}).
		Where("settlement_company_id = ? AND status = ?", companyID, model.WithdrawStatusCompleted)

	if dateFrom != nil {
		query = query.Where("completed_at >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("completed_at < ?", *dateTo)
	}

	var stats struct {
		SettlementCompanyName  string
		TotalWithdrawals       int64
		TotalAmountCents       int64
		TotalTaxDeductedCents  int64
		TotalActualAmountCents int64
	}

	err := query.Select(`
		MAX(settlement_company_name) as settlement_company_name,
		COUNT(*) as total_withdrawals,
		COALESCE(SUM(amount_cents), 0) as total_amount_cents,
		COALESCE(SUM(tax_deducted_cents), 0) as total_tax_deducted_cents,
		COALESCE(SUM(actual_amount_cents), 0) as total_actual_amount_cents
	`).Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	var avgAmount int64
	if stats.TotalWithdrawals > 0 {
		avgAmount = stats.TotalAmountCents / stats.TotalWithdrawals
	}

	return &model.WithdrawRoutingStats{
		SettlementCompanyID:    companyID,
		SettlementCompanyName:  stats.SettlementCompanyName,
		TotalWithdrawals:       stats.TotalWithdrawals,
		TotalAmountCents:       stats.TotalAmountCents,
		TotalTaxDeductedCents:  stats.TotalTaxDeductedCents,
		TotalActualAmountCents: stats.TotalActualAmountCents,
		AverageAmountCents:     avgAmount,
	}, nil
}

// CreateSalaryPaymentRecord 创建工资发放记录
// Requirements: 13.4
func (r *withdrawRepository) CreateSalaryPaymentRecord(ctx context.Context, record *model.SalaryPaymentRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

// GetSalaryPaymentRecordByWithdrawID 根据提现ID获取工资发放记录
func (r *withdrawRepository) GetSalaryPaymentRecordByWithdrawID(ctx context.Context, withdrawID uint64) (*model.SalaryPaymentRecord, error) {
	var record model.SalaryPaymentRecord
	err := r.db.WithContext(ctx).Where("withdraw_id = ?", withdrawID).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &record, nil
}

// UpdateSalaryPaymentRecord 更新工资发放记录
func (r *withdrawRepository) UpdateSalaryPaymentRecord(ctx context.Context, record *model.SalaryPaymentRecord) error {
	return r.db.WithContext(ctx).Save(record).Error
}
