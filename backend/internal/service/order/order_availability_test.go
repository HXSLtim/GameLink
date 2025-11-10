package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

func TestOrderService_GetAvailableOrders_DefaultsAndMapping(t *testing.T) {
	t.Helper()

	start := time.Now().Add(time.Hour)
	end := start.Add(2 * time.Hour)
	gameID := uint64(99)
	userID := uint64(42)

	orderRepo := &spyAvailableOrderRepository{
		orders: []model.Order{
			{
				Base: model.Base{
					ID:        1,
					CreatedAt: start.Add(-15 * time.Minute),
				},
				Title:           "Need boost",
				Description:     "Carry me please",
				UserID:          userID,
				TotalPriceCents: 18800,
				GameID:          &gameID,
				ScheduledStart:  &start,
				ScheduledEnd:    &end,
			},
		},
		total: 1,
	}

	userRepo := &stubUserRepository{
		mockUserRepository: &mockUserRepository{},
		data: map[uint64]*model.User{
			userID: {
				Base: model.Base{ID: userID},
				Name: "Alice",
			},
		},
	}

	gameRepo := &stubGameRepository{
		mockGameRepository: &mockGameRepository{},
		data: map[uint64]*model.Game{
			gameID: {
				Base: model.Base{ID: gameID},
				Name: "Valorant",
			},
		},
	}

	svc := &OrderService{
		orders:      orderRepo,
		players:     &mockPlayerRepository{},
		users:       userRepo,
		games:       gameRepo,
		payments:    &mockPaymentRepository{},
		reviews:     &mockReviewRepository{},
		commissions: &mockCommissionRepository{},
	}

	result, total, err := svc.GetAvailableOrders(context.Background(), AvailableOrdersRequest{})
	if err != nil {
		t.Fatalf("GetAvailableOrders returned error: %v", err)
	}

	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}
	if orderRepo.lastOpts.Page != 1 {
		t.Fatalf("expected default Page=1, got %d", orderRepo.lastOpts.Page)
	}
	if orderRepo.lastOpts.PageSize != 20 {
		t.Fatalf("expected default PageSize=20, got %d", orderRepo.lastOpts.PageSize)
	}
	if len(orderRepo.lastOpts.Statuses) != 1 || orderRepo.lastOpts.Statuses[0] != model.OrderStatusConfirmed {
		t.Fatalf("expected statuses filter to be confirmed, got %+v", orderRepo.lastOpts.Statuses)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 available order, got %d", len(result))
	}

	dto := result[0]
	if dto.GameName != "Valorant" {
		t.Fatalf("expected game name Valorant, got %s", dto.GameName)
	}
	if dto.UserNickname != "Alice" {
		t.Fatalf("expected user nickname Alice, got %s", dto.UserNickname)
	}
	if dto.DurationHours != 2 {
		t.Fatalf("expected duration 2 hours, got %f", dto.DurationHours)
	}
	if dto.PriceCents != 18800 {
		t.Fatalf("expected price 18800, got %d", dto.PriceCents)
	}
}

func TestOrderService_GetAvailableOrders_GameFilterAndFallback(t *testing.T) {
	t.Helper()

	gameID := uint64(7)
	orderRepo := &spyAvailableOrderRepository{
		orders: []model.Order{
			{
				Base:   model.Base{ID: 10},
				Title:  "Any game ok",
				UserID: 100,
				GameID: &gameID,
				Status: model.OrderStatusConfirmed,
			},
		},
		total: 1,
	}

	svc := &OrderService{
		orders:  orderRepo,
		players: &mockPlayerRepository{},
		users: &stubUserRepository{
			mockUserRepository: &mockUserRepository{},
			data:               map[uint64]*model.User{},
		},
		games: &stubGameRepository{
			mockGameRepository: &mockGameRepository{},
			data:               map[uint64]*model.Game{},
		},
		payments:    &mockPaymentRepository{},
		reviews:     &mockReviewRepository{},
		commissions: &mockCommissionRepository{},
	}

	reqGameID := uint64(7)
	result, total, err := svc.GetAvailableOrders(context.Background(), AvailableOrdersRequest{
		GameID:   &reqGameID,
		Page:     2,
		PageSize: 30,
	})
	if err != nil {
		t.Fatalf("GetAvailableOrders returned error: %v", err)
	}

	if orderRepo.lastOpts.GameID == nil || *orderRepo.lastOpts.GameID != reqGameID {
		t.Fatalf("expected GameID filter %d, got %+v", reqGameID, orderRepo.lastOpts.GameID)
	}
	if orderRepo.lastOpts.Page != 2 || orderRepo.lastOpts.PageSize != 30 {
		t.Fatalf("expected paging options to be preserved, got %+v", orderRepo.lastOpts)
	}

	if total != 1 || len(result) != 1 {
		t.Fatalf("expected single result, got total=%d len=%d", total, len(result))
	}

	dto := result[0]
	if dto.GameName != "" {
		t.Fatalf("expected empty game name when lookup fails, got %s", dto.GameName)
	}
	if dto.UserNickname != "" {
		t.Fatalf("expected empty user nickname when lookup fails, got %s", dto.UserNickname)
	}
	if dto.DurationHours != 0 {
		t.Fatalf("expected duration 0 without schedule, got %f", dto.DurationHours)
	}
}

func TestOrderService_GetAvailableOrders_PropagatesErrors(t *testing.T) {
	t.Helper()

	orderRepo := &spyAvailableOrderRepository{
		err: errors.New("db down"),
	}

	svc := &OrderService{
		orders:      orderRepo,
		players:     &mockPlayerRepository{},
		users:       &mockUserRepository{},
		games:       &mockGameRepository{},
		payments:    &mockPaymentRepository{},
		reviews:     &mockReviewRepository{},
		commissions: &mockCommissionRepository{},
	}

	if _, _, err := svc.GetAvailableOrders(context.Background(), AvailableOrdersRequest{}); err == nil {
		t.Fatal("expected error to be propagated from repository")
	}
}

func TestOrderService_BuildOrderTimeline_FullFlow(t *testing.T) {
	t.Helper()

	paidAt := time.Now().Add(5 * time.Minute)
	start := time.Now()
	started := start.Add(30 * time.Minute)
	completed := started.Add(time.Hour)
	refunded := completed.Add(30 * time.Minute)

	svc := &OrderService{
		payments: &stubPaymentRepository{
			mockPaymentRepository: &mockPaymentRepository{},
			listFn: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
				return []model.Payment{
					{
						Base:    model.Base{ID: 1},
						OrderID: 1,
						PaidAt:  &paidAt,
					},
				}, 1, nil
			},
		},
	}

	order := &model.Order{
		Base: model.Base{
			ID:        1,
			CreatedAt: start,
			UpdatedAt: completed.Add(15 * time.Minute),
		},
		Status:         model.OrderStatusCanceled,
		ScheduledStart: &start,
		StartedAt:      &started,
		CompletedAt:    &completed,
		CancelReason:   "user canceled",
		RefundedAt:     &refunded,
	}

	timeline := svc.buildOrderTimeline(order)
	if len(timeline) != 6 {
		t.Fatalf("expected 6 timeline entries, got %d", len(timeline))
	}

	expectedStatuses := []string{
		string(model.OrderStatusPending),
		string(model.OrderStatusConfirmed),
		string(model.OrderStatusInProgress),
		string(model.OrderStatusCompleted),
		string(model.OrderStatusCanceled),
		string(model.OrderStatusRefunded),
	}

	for i, status := range expectedStatuses {
		if timeline[i].Status != status {
			t.Fatalf("expected status %s at index %d, got %s", status, i, timeline[i].Status)
		}
	}

	if !timeline[1].Time.Equal(paidAt) {
		t.Fatalf("expected payment time to use PaidAt, got %v", timeline[1].Time)
	}
	if !timeline[2].Time.Equal(*order.StartedAt) {
		t.Fatalf("expected in-progress time to equal StartedAt, got %v", timeline[2].Time)
	}
	if !timeline[4].Time.Equal(order.UpdatedAt) {
		t.Fatalf("expected cancel time to equal UpdatedAt, got %v", timeline[4].Time)
	}
	if !timeline[5].Time.Equal(*order.RefundedAt) {
		t.Fatalf("expected refund time to equal RefundedAt, got %v", timeline[5].Time)
	}
}

func TestOrderService_BuildOrderTimeline_Fallbacks(t *testing.T) {
	t.Helper()

	svc := &OrderService{
		payments: &stubPaymentRepository{
			mockPaymentRepository: &mockPaymentRepository{},
			listFn: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
				return nil, 0, errors.New("cache miss")
			},
		},
	}

	created := time.Now()
	order := &model.Order{
		Base: model.Base{
			ID:        2,
			CreatedAt: created,
		},
		Status: model.OrderStatusConfirmed,
	}

	timeline := svc.buildOrderTimeline(order)
	if len(timeline) != 2 {
		t.Fatalf("expected 2 timeline entries, got %d", len(timeline))
	}
	if !timeline[1].Time.Equal(created) {
		t.Fatalf("expected confirmed entry to fallback to order creation time, got %v", timeline[1].Time)
	}

	pending := &model.Order{
		Base: model.Base{
			ID:        3,
			CreatedAt: created,
		},
		Status: model.OrderStatusPending,
	}
	pendingTimeline := svc.buildOrderTimeline(pending)
	if len(pendingTimeline) != 1 {
		t.Fatalf("expected only creation entry for pending order, got %d", len(pendingTimeline))
	}
}

func TestOrderService_RecordCommissionAsync(t *testing.T) {
	t.Helper()

	orderRepo := newMockOrderRepository()
	playerID := uint64(501)
	order := &model.Order{
		Base:            model.Base{ID: 1},
		TotalPriceCents: 10000,
	}
	order.SetPlayerID(playerID)
	orderRepo.orders[order.ID] = order

	commissionRepo := &recordingCommissionRepository{
		mockCommissionRepository: &mockCommissionRepository{},
		ruleErr:                  errors.New("no specific rule"),
		defaultRule: &model.CommissionRule{
			Rate: 25,
		},
	}

	svc := &OrderService{
		orders:      orderRepo,
		players:     &mockPlayerRepository{},
		users:       &mockUserRepository{},
		games:       &mockGameRepository{},
		payments:    &mockPaymentRepository{},
		reviews:     &mockReviewRepository{},
		commissions: commissionRepo,
	}

	if err := svc.recordCommissionAsync(context.Background(), order.ID); err != nil {
		t.Fatalf("recordCommissionAsync returned error: %v", err)
	}

	if len(commissionRepo.createdRecords) != 1 {
		t.Fatalf("expected commission record to be created, got %d", len(commissionRepo.createdRecords))
	}

	record := commissionRepo.createdRecords[0]
	if record.PlayerID != playerID {
		t.Fatalf("expected record player %d, got %d", playerID, record.PlayerID)
	}
	if record.CommissionRate != 25 {
		t.Fatalf("expected commission rate 25, got %d", record.CommissionRate)
	}
	if record.CommissionCents != 2500 {
		t.Fatalf("expected commission cents 2500, got %d", record.CommissionCents)
	}
	if record.PlayerIncomeCents != 7500 {
		t.Fatalf("expected player income 7500, got %d", record.PlayerIncomeCents)
	}
}

func TestOrderService_RecordCommissionAsync_EdgeCases(t *testing.T) {
	t.Helper()

	t.Run("skips when record exists", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		order := &model.Order{
			Base:            model.Base{ID: 2},
			TotalPriceCents: 5000,
		}
		order.SetPlayerID(1)
		orderRepo.orders[order.ID] = order

		commissionRepo := &recordingCommissionRepository{
			mockCommissionRepository: &mockCommissionRepository{},
			existingRecord:           &model.CommissionRecord{ID: 99},
		}

		svc := &OrderService{
			orders:      orderRepo,
			players:     &mockPlayerRepository{},
			users:       &mockUserRepository{},
			games:       &mockGameRepository{},
			payments:    &mockPaymentRepository{},
			reviews:     &mockReviewRepository{},
			commissions: commissionRepo,
		}

		if err := svc.recordCommissionAsync(context.Background(), order.ID); err != nil {
			t.Fatalf("expected nil error when record already exists, got %v", err)
		}
		if len(commissionRepo.createdRecords) != 0 {
			t.Fatalf("expected no new record when one already exists")
		}
	})

	t.Run("errors when player missing", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		orderWithoutPlayer := &model.Order{
			Base:            model.Base{ID: 3},
			TotalPriceCents: 8000,
		}
		orderRepo.orders[orderWithoutPlayer.ID] = orderWithoutPlayer

		svc := &OrderService{
			orders:      orderRepo,
			players:     &mockPlayerRepository{},
			users:       &mockUserRepository{},
			games:       &mockGameRepository{},
			payments:    &mockPaymentRepository{},
			reviews:     &mockReviewRepository{},
			commissions: &recordingCommissionRepository{mockCommissionRepository: &mockCommissionRepository{}},
		}

		if err := svc.recordCommissionAsync(context.Background(), orderWithoutPlayer.ID); err == nil {
			t.Fatal("expected error when order has no player assigned")
		}
	})
}

type spyAvailableOrderRepository struct {
	orders   []model.Order
	total    int64
	lastOpts repository.OrderListOptions
	err      error
}

func (s *spyAvailableOrderRepository) Create(ctx context.Context, order *model.Order) error {
	return nil
}

func (s *spyAvailableOrderRepository) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	s.lastOpts = opts
	if s.err != nil {
		return nil, 0, s.err
	}
	return s.orders, s.total, nil
}

func (s *spyAvailableOrderRepository) Get(ctx context.Context, id uint64) (*model.Order, error) {
	for i := range s.orders {
		if s.orders[i].ID == id {
			return &s.orders[i], nil
		}
	}
	return nil, repository.ErrNotFound
}

func (s *spyAvailableOrderRepository) Update(ctx context.Context, order *model.Order) error {
	return nil
}

func (s *spyAvailableOrderRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

type stubUserRepository struct {
	*mockUserRepository
	data map[uint64]*model.User
}

func (s *stubUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	if user, ok := s.data[id]; ok {
		copy := *user
		return &copy, nil
	}
	return nil, repository.ErrNotFound
}

type stubGameRepository struct {
	*mockGameRepository
	data map[uint64]*model.Game
}

func (s *stubGameRepository) Get(ctx context.Context, id uint64) (*model.Game, error) {
	if game, ok := s.data[id]; ok {
		copy := *game
		return &copy, nil
	}
	return nil, repository.ErrNotFound
}

type stubPaymentRepository struct {
	*mockPaymentRepository
	listFn func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error)
}

func (s *stubPaymentRepository) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	if s.listFn != nil {
		return s.listFn(ctx, opts)
	}
	return s.mockPaymentRepository.List(ctx, opts)
}

type recordingCommissionRepository struct {
	*mockCommissionRepository
	existingRecord *model.CommissionRecord
	rule           *model.CommissionRule
	ruleErr        error
	defaultRule    *model.CommissionRule
	defaultErr     error
	createdRecords []*model.CommissionRecord
}

func (r *recordingCommissionRepository) GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
	return r.existingRecord, nil
}

func (r *recordingCommissionRepository) GetRuleForOrder(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
	if r.rule != nil {
		return r.rule, nil
	}
	if r.ruleErr != nil {
		return nil, r.ruleErr
	}
	return nil, nil
}

func (r *recordingCommissionRepository) GetDefaultRule(ctx context.Context) (*model.CommissionRule, error) {
	if r.defaultErr != nil {
		return nil, r.defaultErr
	}
	if r.defaultRule != nil {
		return r.defaultRule, nil
	}
	return &model.CommissionRule{Rate: 20}, nil
}

func (r *recordingCommissionRepository) CreateRecord(ctx context.Context, record *model.CommissionRecord) error {
	copy := *record
	r.createdRecords = append(r.createdRecords, &copy)
	return nil
}
