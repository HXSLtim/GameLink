package earnings

import (
	"context"
	"errors"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	withdrawrepo "gamelink/internal/repository/withdraw"
)

var (
	// ErrNotFound 记录不存在
	ErrNotFound = repository.ErrNotFound
	// ErrValidation 表示输入校验失败
	ErrValidation = errors.New("validation failed")
	// ErrInsufficientBalance 余额不足
	ErrInsufficientBalance = errors.New("insufficient balance")
	// ErrUnauthorized 无权操作
	ErrUnauthorized = errors.New("unauthorized")
)

// WithdrawStatus 提现状态
type WithdrawStatus string

const (
	// WithdrawPending 待处理
	WithdrawPending WithdrawStatus = "pending"
	// WithdrawApproved 已批准
	WithdrawApproved WithdrawStatus = "approved"
	// WithdrawRejected 已拒绝
	WithdrawRejected WithdrawStatus = "rejected"
	// WithdrawCompleted 已完成
	WithdrawCompleted WithdrawStatus = "completed"
)

// EarningsService 收益服务
//
// 功能：
// 1. 收益概览
// 2. 收益趋势
// 3. 提现管理
type EarningsService struct {
	players   repository.PlayerRepository
	orders    repository.OrderRepository
	withdraws withdrawrepo.WithdrawRepository
}

// NewEarningsService 创建收益服务
func NewEarningsService(
	players repository.PlayerRepository,
	orders repository.OrderRepository,
	withdraws withdrawrepo.WithdrawRepository,
) *EarningsService {
	return &EarningsService{
		players:   players,
		orders:    orders,
		withdraws: withdraws,
	}
}

// EarningsSummaryResponse 收益概览响应
type EarningsSummaryResponse struct {
	TodayEarnings    int64 `json:"todayEarnings"`    // 今日收益（分）
	MonthEarnings    int64 `json:"monthEarnings"`    // 本月收益
	TotalEarnings    int64 `json:"totalEarnings"`    // 累计收益
	AvailableBalance int64 `json:"availableBalance"` // 可提现余额
	PendingBalance   int64 `json:"pendingBalance"`   // 待结算余额
	WithdrawTotal    int64 `json:"withdrawTotal"`    // 累计提现
}

// DailyEarningDTO 每日收益
type DailyEarningDTO struct {
	Date       string `json:"date"`       // YYYY-MM-DD
	Earnings   int64  `json:"earnings"`   // 当日收益
	OrderCount int    `json:"orderCount"` // 订单数
}

// EarningsTrendResponse 收益趋势响应
type EarningsTrendResponse struct {
	Trend []DailyEarningDTO `json:"trend"`
}

// WithdrawRequest 提现请求
type WithdrawRequest struct {
	AmountCents int64  `json:"amountCents" binding:"required,min=10000"` // 最低100元
	Method      string `json:"method" binding:"required,oneof=alipay wechat bank"`
	AccountInfo string `json:"accountInfo" binding:"required"` // 账号信息
}

// WithdrawResponse 提现响应
type WithdrawResponse struct {
	WithdrawID uint64 `json:"withdrawId"`
	Status     string `json:"status"`
}

// WithdrawRecordDTO 提现记录
type WithdrawRecordDTO struct {
	ID          uint64     `json:"id"`
	AmountCents int64      `json:"amountCents"`
	Method      string     `json:"method"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	ProcessedAt *time.Time `json:"processedAt"`
}

// WithdrawHistoryResponse 提现记录响应
type WithdrawHistoryResponse struct {
	Records []WithdrawRecordDTO `json:"records"`
	Total   int64               `json:"total"`
}

// GetEarningsSummary 获取收益概览
func (s *EarningsService) GetEarningsSummary(ctx context.Context, userID uint64) (*EarningsSummaryResponse, error) {
	// 查找陪玩师
	player, err := s.findPlayerByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	playerIDPtr := &player.ID

	// 今日收益
	todayStart := time.Now().Truncate(24 * time.Hour)
	todayEnd := todayStart.Add(24 * time.Hour)
	todayEarnings, err := s.calculateEarnings(ctx, playerIDPtr, &todayStart, &todayEnd)
	if err != nil {
		todayEarnings = 0
	}

	// 本月收益
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)
	monthEarnings, err := s.calculateEarnings(ctx, playerIDPtr, &monthStart, &monthEnd)
	if err != nil {
		monthEarnings = 0
	}

	// 累计收益
	totalEarnings, err := s.calculateEarnings(ctx, playerIDPtr, nil, nil)
	if err != nil {
		totalEarnings = 0
	}

	// 从数据库获取余额信息
	var balance *withdrawrepo.PlayerBalance
	if s.withdraws != nil {
		balance, err = s.withdraws.GetPlayerBalance(ctx, player.ID)
	}
	if s.withdraws == nil || err != nil {
		// 如果获取失败或withdraws为nil，使用计算值
		// 80%可提现，20%待结算
		availableBalance := totalEarnings * 8 / 10
		pendingBalance := totalEarnings - availableBalance
		balance = &withdrawrepo.PlayerBalance{
			TotalEarnings:    totalEarnings,
			WithdrawTotal:    0,
			PendingWithdraw:  0,
			AvailableBalance: availableBalance,
			PendingBalance:   pendingBalance,
		}
	}

	return &EarningsSummaryResponse{
		TodayEarnings:    todayEarnings,
		MonthEarnings:    monthEarnings,
		TotalEarnings:    balance.TotalEarnings,
		AvailableBalance: balance.AvailableBalance,
		PendingBalance:   balance.PendingBalance,
		WithdrawTotal:    balance.WithdrawTotal,
	}, nil
}

// GetEarningsTrend 获取收益趋势
func (s *EarningsService) GetEarningsTrend(ctx context.Context, userID uint64, days int) (*EarningsTrendResponse, error) {
	// 查找陪玩师
	player, err := s.findPlayerByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if days < 7 {
		days = 7
	}
	if days > 90 {
		days = 90
	}

	playerIDPtr := &player.ID
	trend := make([]DailyEarningDTO, 0, days)

	// 计算每日收益
	now := time.Now()
	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i).Truncate(24 * time.Hour)
		dateEnd := date.Add(24 * time.Hour)

		earnings, err := s.calculateEarnings(ctx, playerIDPtr, &date, &dateEnd)
		if err != nil {
			earnings = 0
		}

		orderCount, err := s.countOrders(ctx, playerIDPtr, &date, &dateEnd)
		if err != nil {
			orderCount = 0
		}

		trend = append(trend, DailyEarningDTO{
			Date:       date.Format("2006-01-02"),
			Earnings:   earnings,
			OrderCount: int(orderCount),
		})
	}

	return &EarningsTrendResponse{
		Trend: trend,
	}, nil
}

// RequestWithdraw 申请提现
func (s *EarningsService) RequestWithdraw(ctx context.Context, userID uint64, req WithdrawRequest) (*WithdrawResponse, error) {
	// 查找陪玩师
	player, err := s.findPlayerByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 获取可提现余额
	summary, err := s.GetEarningsSummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 检查余额
	if summary.AvailableBalance < req.AmountCents {
		return nil, ErrInsufficientBalance
	}

	// 创建提现记录
	withdraw := &model.Withdraw{
		PlayerID:    player.ID,
		UserID:      userID,
		AmountCents: req.AmountCents,
		Method:      model.WithdrawMethod(req.Method),
		AccountInfo: req.AccountInfo, // TODO: 需要加密存储敏感信息
		Status:      model.WithdrawStatusPending,
	}

	if err := s.withdraws.Create(ctx, withdraw); err != nil {
		return nil, err
	}

	return &WithdrawResponse{
		WithdrawID: withdraw.ID,
		Status:     string(withdraw.Status),
	}, nil
}

// GetWithdrawHistory 获取提现记录
func (s *EarningsService) GetWithdrawHistory(ctx context.Context, userID uint64, page, pageSize int) (*WithdrawHistoryResponse, error) {
	// 查找陪玩师
	player, err := s.findPlayerByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 默认分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 从数据库获取提现记录
	playerID := player.ID
	withdraws, total, err := s.withdraws.List(ctx, withdrawrepo.WithdrawListOptions{
		PlayerID: &playerID,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	records := make([]WithdrawRecordDTO, 0, len(withdraws))
	for _, w := range withdraws {
		records = append(records, WithdrawRecordDTO{
			ID:          w.ID,
			AmountCents: w.AmountCents,
			Method:      string(w.Method),
			Status:      string(w.Status),
			CreatedAt:   w.CreatedAt,
			ProcessedAt: w.ProcessedAt,
		})
	}

	return &WithdrawHistoryResponse{
		Records: records,
		Total:   total,
	}, nil
}

// calculateEarnings 计算收益
func (s *EarningsService) calculateEarnings(ctx context.Context, playerID *uint64, dateFrom, dateTo *time.Time) (int64, error) {
	// 查询已完成的订单
	orders, _, err := s.orders.List(ctx, repository.OrderListOptions{
		PlayerID: playerID,
		Statuses: []model.OrderStatus{model.OrderStatusCompleted},
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Page:     1,
		PageSize: 10000, // 获取所有订单
	})
	if err != nil {
		return 0, err
	}

	var total int64
	for _, o := range orders {
		total += o.TotalPriceCents
	}

	return total, nil
}

// countOrders 统计订单数
func (s *EarningsService) countOrders(ctx context.Context, playerID *uint64, dateFrom, dateTo *time.Time) (int64, error) {
	_, total, err := s.orders.List(ctx, repository.OrderListOptions{
		PlayerID: playerID,
		Statuses: []model.OrderStatus{model.OrderStatusCompleted},
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		return 0, err
	}

	return total, nil
}

// findPlayerByUserID 根据用户ID查找陪玩师
func (s *EarningsService) findPlayerByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	players, _, err := s.players.ListPaged(ctx, 1, 100)
	if err != nil {
		return nil, err
	}

	for _, p := range players {
		if p.UserID == userID {
			return &p, nil
		}
	}

	return nil, ErrNotFound
}
