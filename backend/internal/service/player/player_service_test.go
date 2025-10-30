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
