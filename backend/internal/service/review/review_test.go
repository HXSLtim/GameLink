package review

import (
	"context"
	"testing"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// Mock repositories
type mockReviewRepository struct {
	reviews map[uint64]*model.Review
}

type mockReviewReplyRepository struct{}

func (m *mockReviewReplyRepository) Create(ctx context.Context, reply *model.ReviewReply) error {
	return nil
}

func (m *mockReviewReplyRepository) ListByReview(ctx context.Context, reviewID uint64) ([]model.ReviewReply, error) {
	return []model.ReviewReply{}, nil
}

func (m *mockReviewReplyRepository) UpdateStatus(ctx context.Context, replyID uint64, status string, note string) error {
	return nil
}

func newMockReviewRepository() *mockReviewRepository {
	return &mockReviewRepository{
		reviews: make(map[uint64]*model.Review),
	}
}

func (m *mockReviewRepository) List(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
	var result []model.Review
	for _, r := range m.reviews {
		if opts.OrderID != nil && r.OrderID != *opts.OrderID {
			continue
		}
		if opts.UserID != nil && r.UserID != *opts.UserID {
			continue
		}
		if opts.PlayerID != nil && r.PlayerID != *opts.PlayerID {
			continue
		}
		result = append(result, *r)
	}
	return result, int64(len(result)), nil
}

func (m *mockReviewRepository) Get(ctx context.Context, id uint64) (*model.Review, error) {
	if review, ok := m.reviews[id]; ok {
		return review, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockReviewRepository) Create(ctx context.Context, review *model.Review) error {
	review.ID = uint64(len(m.reviews) + 1)
	m.reviews[review.ID] = review
	return nil
}

func (m *mockReviewRepository) Update(ctx context.Context, review *model.Review) error {
	if _, ok := m.reviews[review.ID]; !ok {
		return repository.ErrNotFound
	}
	m.reviews[review.ID] = review
	return nil
}

func (m *mockReviewRepository) Delete(ctx context.Context, id uint64) error {
	delete(m.reviews, id)
	return nil
}

type mockOrderRepository struct {
	orders map[uint64]*model.Order
}

func newMockOrderRepository() *mockOrderRepository {
	return &mockOrderRepository{
		orders: make(map[uint64]*model.Order),
	}
}

func (m *mockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	order.ID = uint64(len(m.orders) + 1)
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepository) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	return []model.Order{}, 0, nil
}

func (m *mockOrderRepository) Get(ctx context.Context, id uint64) (*model.Order, error) {
	if order, ok := m.orders[id]; ok {
		return order, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockOrderRepository) Update(ctx context.Context, order *model.Order) error {
	if _, ok := m.orders[order.ID]; !ok {
		return repository.ErrNotFound
	}
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepository) Delete(ctx context.Context, id uint64) error {
	delete(m.orders, id)
	return nil
}

type mockPlayerRepository struct{}

func (m *mockPlayerRepository) List(ctx context.Context) ([]model.Player, error) {
	return []model.Player{}, nil
}

func (m *mockPlayerRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	return []model.Player{}, 0, nil
}

func (m *mockPlayerRepository) Get(ctx context.Context, id uint64) (*model.Player, error) {
	return &model.Player{
		Base:          model.Base{ID: id},
		Nickname:      "TestPlayer",
		RatingAverage: 4.5,
		RatingCount:   10,
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
		Base:          model.Base{ID: 1},
		UserID:        userID,
		Nickname:      "TestPlayer",
		RatingAverage: 4.5,
		RatingCount:   10,
	}, nil
}

type mockUserRepository struct{}

func (m *mockUserRepository) List(ctx context.Context) ([]model.User, error) {
	return []model.User{}, nil
}

func (m *mockUserRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	return []model.User{}, 0, nil
}

func (m *mockUserRepository) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	return []model.User{}, 0, nil
}

func (m *mockUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	return &model.User{
		Base:      model.Base{ID: id},
		Name:      "TestUser",
		AvatarURL: "http://example.com/avatar.jpg",
	}, nil
}

func (m *mockUserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, nil
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}

func (m *mockUserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, nil
}

func (m *mockUserRepository) Create(ctx context.Context, user *model.User) error {
	return nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *model.User) error {
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

func TestCreateReview(t *testing.T) {
	reviewRepo := newMockReviewRepository()
	orderRepo := newMockOrderRepository()
	svc := NewReviewService(reviewRepo, orderRepo, &mockPlayerRepository{}, &mockUserRepository{}, &mockReviewReplyRepository{})

	// 创建已完成的订单
	playerID := uint64(1)
	order := &model.Order{
		UserID:   1,
		PlayerID: &playerID,
		Status:   model.OrderStatusCompleted,
		ItemID:   1,
	}
	_ = orderRepo.Create(context.Background(), order)

	// 测试创建评价
	resp, err := svc.CreateReview(context.Background(), 1, CreateReviewRequest{
		OrderID: order.ID,
		Rating:  5,
		Comment: "Very good player!",
		Tags:    []string{"friendly", "skilled"},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}

	if resp.ReviewID == 0 {
		t.Error("expected review ID, got 0")
	}

	// 验证评价已创建
	review, err := reviewRepo.Get(context.Background(), resp.ReviewID)
	if err != nil {
		t.Fatalf("expected to find review, got error: %v", err)
	}

	if review.Score != 5 {
		t.Errorf("expected rating 5, got %d", review.Score)
	}

	if review.Content != "Very good player!" {
		t.Errorf("expected comment 'Very good player!', got %s", review.Content)
	}
}

func TestCreateReviewOrderNotCompleted(t *testing.T) {
	reviewRepo := newMockReviewRepository()
	orderRepo := newMockOrderRepository()
	svc := NewReviewService(reviewRepo, orderRepo, &mockPlayerRepository{}, &mockUserRepository{}, &mockReviewReplyRepository{})

	// 创建进行中的订单
	playerID := uint64(1)
	order := &model.Order{
		UserID:   1,
		PlayerID: &playerID,
		Status:   model.OrderStatusInProgress,
		ItemID:   1,
	}
	_ = orderRepo.Create(context.Background(), order)

	// 测试为未完成订单创建评价（应该失败）
	_, err := svc.CreateReview(context.Background(), 1, CreateReviewRequest{
		OrderID: order.ID,
		Rating:  5,
		Comment: "Good",
	})

	if err != ErrOrderNotCompleted {
		t.Errorf("expected ErrOrderNotCompleted, got %v", err)
	}
}

func TestCreateReviewUnauthorized(t *testing.T) {
	reviewRepo := newMockReviewRepository()
	orderRepo := newMockOrderRepository()
	svc := NewReviewService(reviewRepo, orderRepo, &mockPlayerRepository{}, &mockUserRepository{}, &mockReviewReplyRepository{})

	// 创建其他用户的订单
	playerID := uint64(1)
	order := &model.Order{
		UserID:   2,
		PlayerID: &playerID,
		Status:   model.OrderStatusCompleted,
		ItemID:   1,
	}
	_ = orderRepo.Create(context.Background(), order)

	// 测试用户1评价用户2的订单（应该失败）
	_, err := svc.CreateReview(context.Background(), 1, CreateReviewRequest{
		OrderID: order.ID,
		Rating:  5,
		Comment: "Good",
	})

	if err != ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestCreateReviewAlreadyReviewed(t *testing.T) {
	reviewRepo := newMockReviewRepository()
	orderRepo := newMockOrderRepository()
	svc := NewReviewService(reviewRepo, orderRepo, &mockPlayerRepository{}, &mockUserRepository{}, &mockReviewReplyRepository{})

	// 创建已完成的订单
	playerID := uint64(1)
	order := &model.Order{
		UserID:   1,
		PlayerID: &playerID,
		Status:   model.OrderStatusCompleted,
		ItemID:   1,
	}
	_ = orderRepo.Create(context.Background(), order)

	// 第一次评价
	_, err := svc.CreateReview(context.Background(), 1, CreateReviewRequest{
		OrderID: order.ID,
		Rating:  5,
		Comment: "Good",
	})

	if err != nil {
		t.Fatalf("expected no error on first review, got %v", err)
	}

	// 第二次评价（应该失败）
	_, err = svc.CreateReview(context.Background(), 1, CreateReviewRequest{
		OrderID: order.ID,
		Rating:  4,
		Comment: "Second review",
	})

	if err != ErrAlreadyReviewed {
		t.Errorf("expected ErrAlreadyReviewed, got %v", err)
	}
}

func TestGetMyReviews(t *testing.T) {
	reviewRepo := newMockReviewRepository()
	orderRepo := newMockOrderRepository()
	svc := NewReviewService(reviewRepo, orderRepo, &mockPlayerRepository{}, &mockUserRepository{}, &mockReviewReplyRepository{})

	// 创建订单和评价
	playerID := uint64(1)
	order := &model.Order{
		UserID:   1,
		PlayerID: &playerID,
		Title:    "Test Order",
		Status:   model.OrderStatusCompleted,
		ItemID:   1,
	}
	_ = orderRepo.Create(context.Background(), order)

	review := &model.Review{
		OrderID:  order.ID,
		UserID:   1,
		PlayerID: playerID,
		Score:    5,
		Content:  "Great!",
	}
	_ = reviewRepo.Create(context.Background(), review)

	// 测试获取我的评价列表
	resp, err := svc.GetMyReviews(context.Background(), 1, 1, 20)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}

	if len(resp.Reviews) != 1 {
		t.Errorf("expected 1 review, got %d", len(resp.Reviews))
	}

	if resp.Reviews[0].Rating != 5 {
		t.Errorf("expected rating 5, got %d", resp.Reviews[0].Rating)
	}

	if resp.Reviews[0].OrderTitle != "Test Order" {
		t.Errorf("expected order title 'Test Order', got %s", resp.Reviews[0].OrderTitle)
	}
}

func TestGetPlayerReviews(t *testing.T) {
	reviewRepo := newMockReviewRepository()
	orderRepo := newMockOrderRepository()
	svc := NewReviewService(reviewRepo, orderRepo, &mockPlayerRepository{}, &mockUserRepository{}, &mockReviewReplyRepository{})

	// 创建多个评价
	playerID := uint64(1)
	for i := 1; i <= 3; i++ {
		review := &model.Review{
			OrderID:  uint64(i),
			UserID:   uint64(i),
			PlayerID: playerID,
			Score:    model.Rating(i + 2), // 3, 4, 5
			Content:  "Review " + string(rune(i)),
		}
		_ = reviewRepo.Create(context.Background(), review)
	}

	// 测试获取陪玩师的评价列表
	reviews, total, err := svc.GetPlayerReviews(context.Background(), 1, 1, 20)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if total != 3 {
		t.Errorf("expected 3 reviews, got %d", total)
	}

	if len(reviews) != 3 {
		t.Errorf("expected 3 reviews in list, got %d", len(reviews))
	}
}
