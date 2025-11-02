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
	// GetPlayerBalance 获取陪玩师余额信�?
	GetPlayerBalance(ctx context.Context, playerID uint64) (*PlayerBalance, error)
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

// PlayerBalance 陪玩师余额信�?
type PlayerBalance struct {
	TotalEarnings    int64 // 累计收益
	WithdrawTotal    int64 // 累计提现
	PendingWithdraw  int64 // 待处理提�?
	AvailableBalance int64 // 可提现余�?
	PendingBalance   int64 // 待结算余�?
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

// GetPlayerBalance 获取陪玩师余额信�?
func (r *withdrawRepository) GetPlayerBalance(ctx context.Context, playerID uint64) (*PlayerBalance, error) {
	balance := &PlayerBalance{}

	// 计算累计收益（从已完成订单）
	var totalEarnings int64
	err := r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("player_id = ? AND status = ?", playerID, model.OrderStatusCompleted).
		Select("COALESCE(SUM(price_cents), 0)").
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

	// 计算待处理提现（pending �?approved 状态）
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

	// 计算待结算余额（进行中的订单�?
	var pendingBalance int64
	err = r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("player_id = ? AND status = ?", playerID, model.OrderStatusInProgress).
		Select("COALESCE(SUM(price_cents), 0)").
		Scan(&pendingBalance).Error
	if err != nil {
		return nil, err
	}
	balance.PendingBalance = pendingBalance

	// 计算可提现余�?= 累计收益 - 累计提现 - 待处理提�?- 待结算余�?
	balance.AvailableBalance = totalEarnings - withdrawTotal - pendingWithdraw - pendingBalance
	if balance.AvailableBalance < 0 {
		balance.AvailableBalance = 0
	}

	return balance, nil
}

