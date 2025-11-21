package earnings

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	withdrawrepo "gamelink/internal/repository/withdraw"
)

// Mock repositories
type mockPlayerRepository struct{}

func (m *mockPlayerRepository) List(ctx context.Context) ([]model.Player, error) {
	return []model.Player{}, nil
}

func (m *mockPlayerRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	return []model.Player{
		{
			Base:   model.Base{ID: 1},
			UserID: 1,
		},
	}, 1, nil
}

func (m *mockPlayerRepository) Get(ctx context.Context, id uint64) (*model.Player, error) {
	return &model.Player{
		Base:   model.Base{ID: id},
		UserID: 1,
	}, nil
}

func (m *mockPlayerRepository) Create(ctx context.Context, player *model.Player) error {
	return nil
}

func (m *mockPlayerRepository) Update(ctx context.Context, player *model.Player) error {
	return nil
}

func (m *mockPlayerRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

func (m *mockPlayerRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	return &model.Player{
		Base:   model.Base{ID: 1},
		UserID: userID,
	}, nil
}

type mockOrderRepository struct {
	orders []model.Order
}

func newMockOrderRepository() *mockOrderRepository {
	return &mockOrderRepository{
		orders: []model.Order{},
	}
}

func (m *mockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	order.ID = uint64(len(m.orders) + 1)
	m.orders = append(m.orders, *order)
	return nil
}

func (m *mockOrderRepository) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	var result []model.Order
	for _, o := range m.orders {
		// Filter by PlayerID
		if opts.PlayerID != nil && o.PlayerID != nil && *o.PlayerID != *opts.PlayerID {
			continue
		}
		// Filter by Statuses
		if len(opts.Statuses) > 0 {
			match := false
			for _, s := range opts.Statuses {
				if o.Status == s {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		// Filter by DateFrom
		if opts.DateFrom != nil && o.CreatedAt.Before(*opts.DateFrom) {
			continue
		}
		// Filter by DateTo
		if opts.DateTo != nil && o.CreatedAt.After(*opts.DateTo) {
			continue
		}
		result = append(result, o)
	}
	return result, int64(len(result)), nil
}

func (m *mockOrderRepository) Get(ctx context.Context, id uint64) (*model.Order, error) {
	for _, o := range m.orders {
		if o.ID == id {
			return &o, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockOrderRepository) Update(ctx context.Context, order *model.Order) error {
	for i, o := range m.orders {
		if o.ID == order.ID {
			m.orders[i] = *order
			return nil
		}
	}
	return repository.ErrNotFound
}

func (m *mockOrderRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

// mockWithdrawRepository mock 提现仓库
type mockWithdrawRepository struct {
	data   map[uint64]*model.Withdraw
	nextID uint64
}

func newMockWithdrawRepository() *mockWithdrawRepository {
	return &mockWithdrawRepository{
		data:   make(map[uint64]*model.Withdraw),
		nextID: 0,
	}
}

func (m *mockWithdrawRepository) Create(ctx context.Context, withdraw *model.Withdraw) error {
	m.nextID++
	withdraw.ID = m.nextID
	m.data[withdraw.ID] = withdraw
	return nil
}

func (m *mockWithdrawRepository) Get(ctx context.Context, id uint64) (*model.Withdraw, error) {
	if w, ok := m.data[id]; ok {
		return w, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockWithdrawRepository) Update(ctx context.Context, withdraw *model.Withdraw) error {
	if _, ok := m.data[withdraw.ID]; ok {
		m.data[withdraw.ID] = withdraw
		return nil
	}
	return repository.ErrNotFound
}

func (m *mockWithdrawRepository) List(ctx context.Context, opts withdrawrepo.WithdrawListOptions) ([]model.Withdraw, int64, error) {
	var result []model.Withdraw
	for _, w := range m.data {
		// Filter by PlayerID if specified
		if opts.PlayerID != nil && w.PlayerID != *opts.PlayerID {
			continue
		}
		result = append(result, *w)
	}
	return result, int64(len(result)), nil
}

func (m *mockWithdrawRepository) GetPlayerBalance(ctx context.Context, playerID uint64) (*withdrawrepo.PlayerBalance, error) {
	// 返回错误，让服务使用备用计算方法（基于订单）
	return nil, repository.ErrNotFound
}

func TestGetEarningsSummary(t *testing.T) {
	orderRepo := newMockOrderRepository()
	withdrawRepo := newMockWithdrawRepository()
	playerID := uint64(1)
	svc := NewEarningsService(&mockPlayerRepository{}, orderRepo, withdrawRepo)

	// 创建一些已完成的订单
	now := time.Now()
	for i := 0; i < 5; i++ {
		createdAt := now.Add(time.Duration(-i*24) * time.Hour)
		order := &model.Order{
			Base: model.Base{
				CreatedAt: createdAt,
			},
			PlayerID:        &playerID,
			Status:          model.OrderStatusCompleted,
			TotalPriceCents: 10000,
		}
		_ = orderRepo.Create(context.Background(), order)
	}

	// 测试获取收益概览
	resp, err := svc.GetEarningsSummary(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}

	// 验证累计收益
	if resp.TotalEarnings != 50000 { // 5 * 10000
		t.Errorf("expected total earnings 50000, got %d", resp.TotalEarnings)
	}

	// 验证今日收益（应该有一个订单）
	if resp.TodayEarnings != 10000 {
		t.Errorf("expected today earnings 10000, got %d", resp.TodayEarnings)
	}

	// 验证可提现余额（80%）
	expectedAvailable := resp.TotalEarnings * 8 / 10
	if resp.AvailableBalance != expectedAvailable {
		t.Errorf("expected available balance %d, got %d", expectedAvailable, resp.AvailableBalance)
	}

	// 验证待结算余额（20%）
	expectedPending := resp.TotalEarnings * 2 / 10
	if resp.PendingBalance != expectedPending {
		t.Errorf("expected pending balance %d, got %d", expectedPending, resp.PendingBalance)
	}
}

func TestGetEarningsTrend(t *testing.T) {
	orderRepo := newMockOrderRepository()
	withdrawRepo := newMockWithdrawRepository()
	playerID := uint64(1)
	svc := NewEarningsService(&mockPlayerRepository{}, orderRepo, withdrawRepo)

	// 创建过去7天的订单
	now := time.Now()
	for i := 0; i < 7; i++ {
		createdAt := now.Add(time.Duration(-i*24) * time.Hour)
		order := &model.Order{
			Base: model.Base{
				CreatedAt: createdAt,
			},
			PlayerID:        &playerID,
			Status:          model.OrderStatusCompleted,
			TotalPriceCents: int64((i + 1) * 1000),
		}
		_ = orderRepo.Create(context.Background(), order)
	}

	// 测试获取7天收益趋势
	resp, err := svc.GetEarningsTrend(context.Background(), 1, 7)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}

	if len(resp.Trend) != 7 {
		t.Errorf("expected 7 days trend, got %d", len(resp.Trend))
	}

	// 验证趋势数据格式
	for _, day := range resp.Trend {
		if day.Date == "" {
			t.Error("expected date to be set")
		}
		if day.Earnings < 0 {
			t.Error("earnings should not be negative")
		}
	}
}

func TestGetEarningsTrendInvalidDays(t *testing.T) {
	orderRepo := newMockOrderRepository()
	withdrawRepo := newMockWithdrawRepository()
	svc := NewEarningsService(&mockPlayerRepository{}, orderRepo, withdrawRepo)

	// 测试小于7天（应该自动调整为7）
	resp, err := svc.GetEarningsTrend(context.Background(), 1, 5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.Trend) != 7 {
		t.Errorf("expected 7 days trend (auto-adjusted), got %d", len(resp.Trend))
	}

	// 测试大于90天（应该自动调整为90）
	resp, err = svc.GetEarningsTrend(context.Background(), 1, 100)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.Trend) != 90 {
		t.Errorf("expected 90 days trend (auto-adjusted), got %d", len(resp.Trend))
	}
}

func TestRequestWithdraw(t *testing.T) {
	orderRepo := newMockOrderRepository()
	withdrawRepo := newMockWithdrawRepository()
	playerID := uint64(1)
	svc := NewEarningsService(&mockPlayerRepository{}, orderRepo, withdrawRepo)

	// 创建足够的收益
	for i := 0; i < 10; i++ {
		order := &model.Order{
			Base: model.Base{
				CreatedAt: time.Now(),
			},
			PlayerID:        &playerID,
			Status:          model.OrderStatusCompleted,
			TotalPriceCents: 20000,
		}
		_ = orderRepo.Create(context.Background(), order)
	}

	// 测试申请提现
	resp, err := svc.RequestWithdraw(context.Background(), 1, WithdrawRequest{
		AmountCents: 10000, // 100元
		Method:      "alipay",
		AccountInfo: "test@example.com",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}

	if resp.WithdrawID == 0 {
		t.Error("expected withdraw ID, got 0")
	}

	if resp.Status != string(WithdrawPending) {
		t.Errorf("expected status pending, got %s", resp.Status)
	}
}

func TestRequestWithdrawInsufficientBalance(t *testing.T) {
	orderRepo := newMockOrderRepository()
	withdrawRepo := newMockWithdrawRepository()
	playerID := uint64(1)
	svc := NewEarningsService(&mockPlayerRepository{}, orderRepo, withdrawRepo)

	// 只创建少量收益
	order := &model.Order{
		Base: model.Base{
			CreatedAt: time.Now(),
		},
		PlayerID:        &playerID,
		Status:          model.OrderStatusCompleted,
		TotalPriceCents: 5000, // 50元
	}
	_ = orderRepo.Create(context.Background(), order)

	// 测试提现超过余额（应该失败）
	_, err := svc.RequestWithdraw(context.Background(), 1, WithdrawRequest{
		AmountCents: 10000, // 100元
		Method:      "alipay",
		AccountInfo: "test@example.com",
	})

	if err != ErrInsufficientBalance {
		t.Errorf("expected ErrInsufficientBalance, got %v", err)
	}
}

func TestGetWithdrawHistory(t *testing.T) {
	orderRepo := newMockOrderRepository()
	withdrawRepo := newMockWithdrawRepository()
	svc := NewEarningsService(&mockPlayerRepository{}, orderRepo, withdrawRepo)

	// 测试获取提现记录（当前返回空列表）
	resp, err := svc.GetWithdrawHistory(context.Background(), 1, 1, 20)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}

	if resp.Total != 0 {
		t.Errorf("expected 0 records, got %d", resp.Total)
	}

	if len(resp.Records) != 0 {
		t.Errorf("expected empty records, got %d", len(resp.Records))
	}
}

func TestFindPlayerByUserID(t *testing.T) {
	orderRepo := newMockOrderRepository()
	withdrawRepo := newMockWithdrawRepository()
	svc := NewEarningsService(&mockPlayerRepository{}, orderRepo, withdrawRepo)

	// 测试查找陪玩师
	player, err := svc.findPlayerByUserID(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if player == nil {
		t.Fatal("expected player, got nil")
	}

	if player.UserID != 1 {
		t.Errorf("expected user ID 1, got %d", player.UserID)
	}
}

func TestFindPlayerByUserIDNotFound(t *testing.T) {
	orderRepo := newMockOrderRepository()
	withdrawRepo := newMockWithdrawRepository()
	svc := NewEarningsService(&mockPlayerRepository{}, orderRepo, withdrawRepo)

	// 测试查找不存在的陪玩师
	_, err := svc.findPlayerByUserID(context.Background(), 999)

	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
