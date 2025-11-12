package order

import (
    "context"
    "testing"
    "time"

    "gamelink/internal/model"
    "gamelink/internal/repository"
)

type mockChatGroupRepo struct{
	lastOrderID uint64
	lastDeactivatedID uint64
	group *model.ChatGroup
}

func (m *mockChatGroupRepo) Create(ctx context.Context, group *model.ChatGroup) error { return nil }
func (m *mockChatGroupRepo) Get(ctx context.Context, id uint64) (*model.ChatGroup, error) { return m.group, nil }
func (m *mockChatGroupRepo) GetByRelatedOrderID(ctx context.Context, orderID uint64) (*model.ChatGroup, error) {
	m.lastOrderID = orderID
	return m.group, nil
}
func (m *mockChatGroupRepo) ListByUser(ctx context.Context, userID uint64, opts repository.ChatGroupListOptions) ([]model.ChatGroup, int64, error) { return nil, 0, nil }
func (m *mockChatGroupRepo) ListMembers(ctx context.Context, groupID uint64, opts repository.ChatGroupMemberListOptions) ([]model.ChatGroupMember, int64, error) { return nil, 0, nil }
func (m *mockChatGroupRepo) Update(ctx context.Context, group *model.ChatGroup) error { return nil }
func (m *mockChatGroupRepo) Deactivate(ctx context.Context, id uint64) error { m.lastDeactivatedID = id; return nil }
func (m *mockChatGroupRepo) ListDeactivatedBefore(ctx context.Context, cutoffTime time.Time, limit int) ([]model.ChatGroup, error) { return nil, nil }
func (m *mockChatGroupRepo) DeleteByIDs(ctx context.Context, ids []uint64) error { return nil }

// Test that CancelOrder triggers auto-deactivation of order chat group
func TestCancelOrder_AutoDeactivateOrderChat(t *testing.T) {
	orderRepo := newMockOrderRepository()
	// seed order: pending & owned by user 100
	order := &model.Order{ Base: model.Base{ID: 1}, UserID: 100, Status: model.OrderStatusPending }
	orderRepo.orders[1] = order

	svc := NewOrderService(orderRepo, &mockPlayerRepository{}, &mockUserRepository{}, &mockGameRepository{}, &mockPaymentRepository{}, &mockReviewRepository{}, &mockCommissionRepository{})
	chatRepo := &mockChatGroupRepo{ group: &model.ChatGroup{ Base: model.Base{ID: 55}, GroupType: model.ChatGroupTypeOrder, IsActive: true } }
	svc.SetChatGroupRepository(chatRepo)

	if err := svc.CancelOrder(context.Background(), 100, 1, CancelOrderRequest{Reason: "x"}); err != nil {
		t.Fatalf("cancel: %v", err)
	}
	if chatRepo.lastDeactivatedID != 55 {
		t.Fatalf("expected deactivated group=55, got %d", chatRepo.lastDeactivatedID)
	}
}

// Test that CompleteOrder triggers auto-deactivation
func TestCompleteOrder_AutoDeactivateOrderChat(t *testing.T) {
	orderRepo := newMockOrderRepository()
	// seed order: in_progress & owned by user 100
	now := time.Now().Add(-time.Hour)
	order := &model.Order{ Base: model.Base{ID: 2}, UserID: 100, Status: model.OrderStatusInProgress, StartedAt: &now }
	orderRepo.orders[2] = order

	svc := NewOrderService(orderRepo, &mockPlayerRepository{}, &mockUserRepository{}, &mockGameRepository{}, &mockPaymentRepository{}, &mockReviewRepository{}, &mockCommissionRepository{})
	chatRepo := &mockChatGroupRepo{ group: &model.ChatGroup{ Base: model.Base{ID: 77}, GroupType: model.ChatGroupTypeOrder, IsActive: true } }
	svc.SetChatGroupRepository(chatRepo)

	if err := svc.CompleteOrder(context.Background(), 100, 2); err != nil {
		t.Fatalf("complete: %v", err)
	}
	if chatRepo.lastDeactivatedID != 77 {
		t.Fatalf("expected deactivated group=77, got %d", chatRepo.lastDeactivatedID)
	}
}

// Test that CompleteOrderByPlayer triggers auto-deactivation
func TestCompleteOrderByPlayer_NoAutoDeactivateOrderChat(t *testing.T) {
	orderRepo := newMockOrderRepository()
	// seed order: in_progress & assigned to player 1
	pid := uint64(1)
	order := &model.Order{ Base: model.Base{ID: 3}, UserID: 200, Status: model.OrderStatusInProgress }
	order.PlayerID = &pid
	orderRepo.orders[3] = order

	svc := NewOrderService(orderRepo, &mockPlayerRepository{}, &mockUserRepository{}, &mockGameRepository{}, &mockPaymentRepository{}, &mockReviewRepository{}, &mockCommissionRepository{})
	chatRepo := &mockChatGroupRepo{ group: &model.ChatGroup{ Base: model.Base{ID: 88}, GroupType: model.ChatGroupTypeOrder, IsActive: true } }
	svc.SetChatGroupRepository(chatRepo)

    if err := svc.CompleteOrderByPlayer(context.Background(), 1, 3); err != nil {
        t.Fatalf("complete by player: %v", err)
    }
    if chatRepo.lastDeactivatedID != 0 {
        t.Fatalf("expected no auto-deactivate, got %d", chatRepo.lastDeactivatedID)
    }
}
