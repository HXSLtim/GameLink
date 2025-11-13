package player

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// Mock repositories for testing
type mockPlayerRepository struct{}

func (m *mockPlayerRepository) List(ctx context.Context) ([]model.Player, error) {
	return []model.Player{}, nil
}

func (m *mockPlayerRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	return []model.Player{
		{
			Base:               model.Base{ID: 1},
			UserID:             1,
			Nickname:           "TestPlayer",
			Bio:                "Test bio",
			Rank:               "Diamond",
			RatingAverage:      4.5,
			RatingCount:        10,
			HourlyRateCents:    10000,
			MainGameID:         1,
			VerificationStatus: model.VerificationVerified,
		},
	}, 1, nil
}

func (m *mockPlayerRepository) Get(ctx context.Context, id uint64) (*model.Player, error) {
	return &model.Player{
		Base:               model.Base{ID: id},
		UserID:             1,
		Nickname:           "TestPlayer",
		Bio:                "Test bio",
		Rank:               "Diamond",
		RatingAverage:      4.5,
		RatingCount:        10,
		HourlyRateCents:    10000,
		MainGameID:         1,
		VerificationStatus: model.VerificationVerified,
	}, nil
}

func (m *mockPlayerRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	return &model.Player{
		Base:               model.Base{ID: userID},
		UserID:             userID,
		Nickname:           "TestPlayer",
		Bio:                "Test bio",
		Rank:               "Diamond",
		RatingAverage:      4.5,
		RatingCount:        10,
		HourlyRateCents:    10000,
		MainGameID:         1,
		VerificationStatus: model.VerificationVerified,
	}, nil
}

func (m *mockPlayerRepository) Create(ctx context.Context, player *model.Player) error {
	player.ID = 1
	return nil
}

func (m *mockPlayerRepository) Update(ctx context.Context, player *model.Player) error {
	return nil
}

func (m *mockPlayerRepository) Delete(ctx context.Context, id uint64) error {
	return nil
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
		Role:      model.RolePlayer,
		Status:    model.UserStatusActive,
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
	user.ID = 1
	return nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *model.User) error {
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

type mockGameRepository struct{}

func (m *mockGameRepository) List(ctx context.Context) ([]model.Game, error) {
	return []model.Game{}, nil
}

func (m *mockGameRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Game, int64, error) {
	return []model.Game{}, 0, nil
}

func (m *mockGameRepository) Get(ctx context.Context, id uint64) (*model.Game, error) {
	return &model.Game{
		Base: model.Base{ID: id},
		Name: "TestGame",
	}, nil
}

func (m *mockGameRepository) Create(ctx context.Context, game *model.Game) error {
	return nil
}

func (m *mockGameRepository) Update(ctx context.Context, game *model.Game) error {
	return nil
}

func (m *mockGameRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

type mockOrderRepository struct{}

func (m *mockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	return nil
}

func (m *mockOrderRepository) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	return []model.Order{}, 10, nil
}

func (m *mockOrderRepository) Get(ctx context.Context, id uint64) (*model.Order, error) {
	return &model.Order{}, nil
}

func (m *mockOrderRepository) Update(ctx context.Context, order *model.Order) error {
	return nil
}

func (m *mockOrderRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

type mockReviewRepository struct{}

func (m *mockReviewRepository) List(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
	return []model.Review{}, 0, nil
}

func (m *mockReviewRepository) Get(ctx context.Context, id uint64) (*model.Review, error) {
	return &model.Review{}, nil
}

func (m *mockReviewRepository) Create(ctx context.Context, review *model.Review) error {
	return nil
}

func (m *mockReviewRepository) Update(ctx context.Context, review *model.Review) error {
	return nil
}

func (m *mockReviewRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

type mockPlayerTagRepository struct{}

func (m *mockPlayerTagRepository) GetTags(ctx context.Context, playerID uint64) ([]string, error) {
	return []string{"friendly", "skilled"}, nil
}

func (m *mockPlayerTagRepository) ReplaceTags(ctx context.Context, playerID uint64, tags []string) error {
	return nil
}

type mockCache struct{}

func (m *mockCache) Get(ctx context.Context, key string) (string, bool, error) {
	return "", false, nil
}

func (m *mockCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return nil
}

func (m *mockCache) Delete(ctx context.Context, key string) error {
	return nil
}

func (m *mockCache) Clear(ctx context.Context) error {
	return nil
}

func (m *mockCache) Close(ctx context.Context) error {
	return nil
}

func TestListPlayers(t *testing.T) {
	svc := NewPlayerService(
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockOrderRepository{},
		&mockReviewRepository{},
		&mockPlayerTagRepository{},
		&mockCache{},
	)

	resp, err := svc.ListPlayers(context.Background(), PlayerListRequest{
		Page:     1,
		PageSize: 20,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}

	if len(resp.Players) != 1 {
		t.Errorf("expected 1 player, got %d", len(resp.Players))
	}
}

func TestGetPlayerDetail(t *testing.T) {
	svc := NewPlayerService(
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockOrderRepository{},
		&mockReviewRepository{},
		&mockPlayerTagRepository{},
		&mockCache{},
	)

	resp, err := svc.GetPlayerDetail(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}

	if resp.Player.Nickname != "TestPlayer" {
		t.Errorf("expected TestPlayer, got %s", resp.Player.Nickname)
	}
}

func TestApplyAsPlayer(t *testing.T) {
	svc := NewPlayerService(
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockOrderRepository{},
		&mockReviewRepository{},
		&mockPlayerTagRepository{},
		&mockCache{},
	)

	// 使用 UserID=2，因为 mockPlayerRepository 返回的是 UserID=1 的陪玩师
	resp, err := svc.ApplyAsPlayer(context.Background(), 2, ApplyPlayerRequest{
		Nickname:        "NewPlayer",
		Bio:             "I'm a new player",
		MainGameID:      1,
		Rank:            "Diamond",
		HourlyRateCents: 10000,
		Tags:            []string{"friendly"},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}

	if resp.VerificationStatus != model.VerificationPending {
		t.Errorf("expected pending status, got %s", resp.VerificationStatus)
	}
}

func TestUpdatePlayerProfile(t *testing.T) {
	svc := NewPlayerService(
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockOrderRepository{},
		&mockReviewRepository{},
		&mockPlayerTagRepository{},
		&mockCache{},
	)

	// 测试更新资料
	req := UpdatePlayerProfileRequest{
		Nickname:        "UpdatedPlayer",
		Bio:             "Updated bio",
		HourlyRateCents: 15000,
		Tags:            []string{"skilled", "friendly"},
	}

	err := svc.UpdatePlayerProfile(context.Background(), 1, req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSetPlayerOnlineStatus(t *testing.T) {
	svc := NewPlayerService(
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockOrderRepository{},
		&mockReviewRepository{},
		&mockPlayerTagRepository{},
		&mockCache{},
	)

	// 测试设置在线状态
	err := svc.SetPlayerOnlineStatus(context.Background(), 1, true)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 测试设置离线状态
	err = svc.SetPlayerOnlineStatus(context.Background(), 1, false)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestListPlayersWithFilters(t *testing.T) {
	svc := NewPlayerService(
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockOrderRepository{},
		&mockReviewRepository{},
		&mockPlayerTagRepository{},
		&mockCache{},
	)

	// 测试带过滤条件的列表
	gameID := uint64(1)
	minRating := float32(4.0)
	resp, err := svc.ListPlayers(context.Background(), PlayerListRequest{
		Page:      1,
		PageSize:  20,
		GameID:    &gameID,
		MinRating: &minRating,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}
}

func TestGetPlayerProfile(t *testing.T) {
	svc := NewPlayerService(
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockOrderRepository{},
		&mockReviewRepository{},
		&mockPlayerTagRepository{},
		&mockCache{},
	)

	// 测试获取陪玩师个人资料
	resp, err := svc.GetPlayerProfile(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}

	if resp.Player.Nickname != "TestPlayer" {
		t.Errorf("expected nickname 'TestPlayer', got '%s'", resp.Player.Nickname)
	}
}

// TestGetPlayerDetail_WithStats 测试获取陪玩师详情（包含统计数据）
func TestGetPlayerDetail_WithStats(t *testing.T) {
	// 创建mock repositories with data
	playerID := uint64(1)
	orderRepo := &mockOrderRepositoryWithData{
		orders: []model.Order{
			{
				Base:     model.Base{ID: 1, CreatedAt: time.Now()},
				UserID:   1,
				PlayerID: &playerID,
				Status:   model.OrderStatusCompleted,
			},
			{
				Base:     model.Base{ID: 2, CreatedAt: time.Now()},
				UserID:   1,
				PlayerID: &playerID,
				Status:   model.OrderStatusCompleted,
			},
		},
	}

	reviewRepo := &mockReviewRepositoryWithData{
		reviews: []model.Review{
			{
				Base:     model.Base{ID: 1, CreatedAt: time.Now()},
				UserID:   1,
				PlayerID: playerID,
				Score:    5,
			},
			{
				Base:     model.Base{ID: 2, CreatedAt: time.Now()},
				UserID:   1,
				PlayerID: playerID,
				Score:    4,
			},
			{
				Base:     model.Base{ID: 3, CreatedAt: time.Now()},
				UserID:   1,
				PlayerID: playerID,
				Score:    3,
			},
		},
	}

	svc := NewPlayerService(
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		orderRepo,
		reviewRepo,
		&mockPlayerTagRepository{},
		&mockCache{},
	)

	resp, err := svc.GetPlayerDetail(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}

	// 验证统计数据
	if resp.Stats.TotalOrders == 0 {
		t.Error("expected total orders > 0")
	}

	// 验证好评率（3个评价中2个是4分以上，应该是2/3 ≈ 0.67）
	if resp.Player.GoodRatio < 0.6 || resp.Player.GoodRatio > 0.7 {
		t.Errorf("expected good ratio around 0.67, got %f", resp.Player.GoodRatio)
	}

	// 验证评价列表
	if len(resp.Reviews) == 0 {
		t.Error("expected reviews")
	}
}

// mockOrderRepositoryWithData 提供数据的mock订单仓库
type mockOrderRepositoryWithData struct {
	orders []model.Order
}

func (m *mockOrderRepositoryWithData) Create(ctx context.Context, order *model.Order) error {
	return nil
}

func (m *mockOrderRepositoryWithData) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	var filtered []model.Order
	for _, o := range m.orders {
		if opts.PlayerID != nil && o.PlayerID != nil && *o.PlayerID != *opts.PlayerID {
			continue
		}
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
		filtered = append(filtered, o)
	}
	return filtered, int64(len(filtered)), nil
}

func (m *mockOrderRepositoryWithData) Get(ctx context.Context, id uint64) (*model.Order, error) {
	return &model.Order{}, nil
}

func (m *mockOrderRepositoryWithData) Update(ctx context.Context, order *model.Order) error {
	return nil
}

func (m *mockOrderRepositoryWithData) Delete(ctx context.Context, id uint64) error {
	return nil
}

// mockReviewRepositoryWithData 提供数据的mock评价仓库
type mockReviewRepositoryWithData struct {
	reviews []model.Review
}

func (m *mockReviewRepositoryWithData) List(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
	var filtered []model.Review
	for _, r := range m.reviews {
		if opts.PlayerID != nil && r.PlayerID != *opts.PlayerID {
			continue
		}
		filtered = append(filtered, r)
	}
	// 限制返回数量
	if opts.PageSize > 0 && len(filtered) > opts.PageSize {
		filtered = filtered[:opts.PageSize]
	}
	return filtered, int64(len(filtered)), nil
}

func (m *mockReviewRepositoryWithData) Get(ctx context.Context, id uint64) (*model.Review, error) {
	return &model.Review{}, nil
}

func (m *mockReviewRepositoryWithData) Create(ctx context.Context, review *model.Review) error {
	return nil
}

func (m *mockReviewRepositoryWithData) Update(ctx context.Context, review *model.Review) error {
	return nil
}

func (m *mockReviewRepositoryWithData) Delete(ctx context.Context, id uint64) error {
	return nil
}

// TestGetPlayerDetail_CalculateRepeatRate 测试复购率计算
func TestGetPlayerDetail_CalculateRepeatRate(t *testing.T) {
	now := time.Now()
	playerID := uint64(1)

	// 创建有复购的订单数据（用户1有2个订单，用户2有1个订单）
	orderRepo := &mockOrderRepositoryWithData{
		orders: []model.Order{
			{
				Base:     model.Base{ID: 1, CreatedAt: now},
				UserID:   1,
				PlayerID: &playerID,
				Status:   model.OrderStatusCompleted,
			},
			{
				Base:     model.Base{ID: 2, CreatedAt: now},
				UserID:   1,
				PlayerID: &playerID,
				Status:   model.OrderStatusCompleted,
			},
			{
				Base:     model.Base{ID: 3, CreatedAt: now},
				UserID:   2,
				PlayerID: &playerID,
				Status:   model.OrderStatusCompleted,
			},
		},
	}

	svc := NewPlayerService(
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		orderRepo,
		&mockReviewRepository{},
		&mockPlayerTagRepository{},
		&mockCache{},
	)

	resp, err := svc.GetPlayerDetail(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 2个用户，1个有复购，复购率应该是0.5
	if resp.Stats.RepeatRate < 0.4 || resp.Stats.RepeatRate > 0.6 {
		t.Errorf("expected repeat rate around 0.5, got %f", resp.Stats.RepeatRate)
	}
}

// TestGetPlayerDetail_CalculateAvgResponseTime 测试平均响应时间计算
func TestGetPlayerDetail_CalculateAvgResponseTime(t *testing.T) {
	now := time.Now()
	playerID := uint64(1)
	startedAt1 := now.Add(30 * time.Minute)
	startedAt2 := now.Add(60 * time.Minute)

	orderRepo := &mockOrderRepositoryWithData{
		orders: []model.Order{
			{
				Base:      model.Base{ID: 1, CreatedAt: now},
				UserID:    1,
				PlayerID:  &playerID,
				Status:    model.OrderStatusCompleted,
				StartedAt: &startedAt1,
			},
			{
				Base:      model.Base{ID: 2, CreatedAt: now},
				UserID:    1,
				PlayerID:  &playerID,
				Status:    model.OrderStatusCompleted,
				StartedAt: &startedAt2,
			},
		},
	}

	svc := NewPlayerService(
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		orderRepo,
		&mockReviewRepository{},
		&mockPlayerTagRepository{},
		&mockCache{},
	)

	resp, err := svc.GetPlayerDetail(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 平均响应时间应该是 (30+60)/2 = 45分钟
	if resp.Player.AvgResponseMin < 40 || resp.Player.AvgResponseMin > 50 {
		t.Errorf("expected avg response time around 45 minutes, got %d", resp.Player.AvgResponseMin)
	}
}
