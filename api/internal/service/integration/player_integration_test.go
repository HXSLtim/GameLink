// Package integration provides integration tests for services.
package integration

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository/game"
	repoiface "gamelink/internal/repository/interfaces"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/review"
	"gamelink/internal/repository/user"
	playerservice "gamelink/internal/service/player"
	"gamelink/pkg/cache"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockOrderQuery implements repoiface.OrderQuery for testing
type mockOrderQuery struct {
	orders []model.Order
}

func (m *mockOrderQuery) List(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
	var result []model.Order
	for _, o := range m.orders {
		if opts.PlayerID != nil && o.PlayerID != nil && *o.PlayerID == *opts.PlayerID {
			result = append(result, o)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockOrderQuery) Get(ctx context.Context, id uint64) (*model.Order, error) {
	for _, o := range m.orders {
		if o.ID == id {
			return &o, nil
		}
	}
	return nil, nil
}

// mockPlayerTagRepo implements repository.PlayerTagRepository for testing
type mockPlayerTagRepo struct {
	tags map[uint64][]string
}

func (m *mockPlayerTagRepo) GetTags(ctx context.Context, playerID uint64) ([]string, error) {
	if tags, ok := m.tags[playerID]; ok {
		return tags, nil
	}
	return []string{}, nil
}

func (m *mockPlayerTagRepo) ReplaceTags(ctx context.Context, playerID uint64, tags []string) error {
	m.tags[playerID] = tags
	return nil
}

// mockCache implements cache.Cache for testing
type mockCache struct {
	data map[string]string
}

func newMockCache() *mockCache {
	return &mockCache{data: make(map[string]string)}
}

func (m *mockCache) Get(ctx context.Context, key string) (string, bool, error) {
	v, ok := m.data[key]
	return v, ok, nil
}

func (m *mockCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	m.data[key] = value
	return nil
}

func (m *mockCache) Delete(ctx context.Context, key string) error {
	delete(m.data, key)
	return nil
}

func (m *mockCache) Exists(ctx context.Context, key string) (bool, error) {
	_, ok := m.data[key]
	return ok, nil
}

func (m *mockCache) Incr(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (m *mockCache) Expire(ctx context.Context, key string, ttl interface{}) error {
	return nil
}

func (m *mockCache) TTL(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (m *mockCache) Close(ctx context.Context) error {
	return nil
}

func (m *mockCache) GetRedisClient() interface{} {
	return nil
}

var _ cache.Cache = (*mockCache)(nil)

func TestPlayerService_ListPlayers(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	playerRepo := player.NewPlayerRepository(db)
	userRepo := user.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	mockOrders := &mockOrderQuery{}
	mockTags := &mockPlayerTagRepo{tags: make(map[uint64][]string)}
	mockCacheInst := newMockCache()

	svc := playerservice.NewPlayerService(playerRepo, userRepo, gameRepo, mockOrders, reviewRepo, mockTags, mockCacheInst)

	// Create test data
	testGame := CreateTestGame(t, db, "test_game_list")

	// Create multiple players
	for i := 0; i < 5; i++ {
		playerUser := CreateUniqueTestUser(t, db, "list_player")
		p := CreateTestPlayer(t, db, playerUser)
		p.MainGameID = testGame.ID
		p.HourlyRateCents = int64(5000 + i*1000)
		p.RatingAverage = float32(4.0 + float32(i)*0.2)
		db.Save(p)
	}

	// List players
	resp, err := svc.ListPlayers(ctx, playerservice.PlayerListRequest{
		Page:     1,
		PageSize: 20,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(resp.Players), 5)
}

func TestPlayerService_ListPlayers_FilterByGame(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerRepo := player.NewPlayerRepository(db)
	userRepo := user.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	mockOrders := &mockOrderQuery{}
	mockTags := &mockPlayerTagRepo{tags: make(map[uint64][]string)}
	mockCacheInst := newMockCache()

	svc := playerservice.NewPlayerService(playerRepo, userRepo, gameRepo, mockOrders, reviewRepo, mockTags, mockCacheInst)

	// Create test games
	game1 := CreateTestGame(t, db, "filter_game_1")
	game2 := CreateTestGame(t, db, "filter_game_2")

	// Create players for game1
	for i := 0; i < 3; i++ {
		playerUser := CreateUniqueTestUser(t, db, "game1_player")
		p := CreateTestPlayer(t, db, playerUser)
		p.MainGameID = game1.ID
		db.Save(p)
	}

	// Create players for game2
	for i := 0; i < 2; i++ {
		playerUser := CreateUniqueTestUser(t, db, "game2_player")
		p := CreateTestPlayer(t, db, playerUser)
		p.MainGameID = game2.ID
		db.Save(p)
	}

	// Filter by game1
	resp, err := svc.ListPlayers(ctx, playerservice.PlayerListRequest{
		GameID:   &game1.ID,
		Page:     1,
		PageSize: 20,
	})
	require.NoError(t, err)
	assert.Equal(t, 3, len(resp.Players))
}

func TestPlayerService_GetPlayerDetail(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerRepo := player.NewPlayerRepository(db)
	userRepo := user.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	mockOrders := &mockOrderQuery{}
	mockTags := &mockPlayerTagRepo{tags: make(map[uint64][]string)}
	mockCacheInst := newMockCache()

	svc := playerservice.NewPlayerService(playerRepo, userRepo, gameRepo, mockOrders, reviewRepo, mockTags, mockCacheInst)

	// Create test data
	testGame := CreateTestGame(t, db, "detail_game")
	playerUser := CreateUniqueTestUser(t, db, "detail_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testPlayer.MainGameID = testGame.ID
	testPlayer.Bio = "Test bio"
	testPlayer.HourlyRateCents = 5000
	db.Save(testPlayer)

	// Add tags
	mockTags.tags[testPlayer.ID] = []string{"friendly", "skilled"}

	// Get detail
	resp, err := svc.GetPlayerDetail(ctx, testPlayer.ID)
	require.NoError(t, err)
	assert.Equal(t, testPlayer.ID, resp.Player.ID)
	assert.Equal(t, testPlayer.Bio, resp.Player.Bio)
	assert.Equal(t, testGame.Name, resp.Player.MainGame)
	assert.Equal(t, []string{"friendly", "skilled"}, resp.Player.Tags)
}

func TestPlayerService_GetPlayerDetail_NotVerified(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerRepo := player.NewPlayerRepository(db)
	userRepo := user.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	mockOrders := &mockOrderQuery{}
	mockTags := &mockPlayerTagRepo{tags: make(map[uint64][]string)}
	mockCacheInst := newMockCache()

	svc := playerservice.NewPlayerService(playerRepo, userRepo, gameRepo, mockOrders, reviewRepo, mockTags, mockCacheInst)

	// Create unverified player
	playerUser := CreateUniqueTestUser(t, db, "unverified_player")
	testPlayer := &model.Player{
		UserID:             playerUser.ID,
		Nickname:           "Unverified",
		VerificationStatus: model.VerificationPending,
	}
	db.Create(testPlayer)

	// Try to get detail
	_, err := svc.GetPlayerDetail(ctx, testPlayer.ID)
	assert.Error(t, err)
	assert.Equal(t, playerservice.ErrPlayerNotVerified, err)
}

func TestPlayerService_ApplyAsPlayer(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerRepo := player.NewPlayerRepository(db)
	userRepo := user.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	mockOrders := &mockOrderQuery{}
	mockTags := &mockPlayerTagRepo{tags: make(map[uint64][]string)}
	mockCacheInst := newMockCache()

	svc := playerservice.NewPlayerService(playerRepo, userRepo, gameRepo, mockOrders, reviewRepo, mockTags, mockCacheInst)

	// Create test data
	testGame := CreateTestGame(t, db, "apply_game")
	testUser := CreateUniqueTestUser(t, db, "apply_user")

	// Apply as player
	req := playerservice.ApplyPlayerRequest{
		Nickname:        "New Player",
		Bio:             "I am a skilled player",
		MainGameID:      testGame.ID,
		Rank:            "Diamond",
		HourlyRateCents: 5000,
		Tags:            []string{"friendly", "patient"},
	}

	resp, err := svc.ApplyAsPlayer(ctx, testUser.ID, req)
	require.NoError(t, err)
	assert.NotZero(t, resp.PlayerID)
	assert.Equal(t, model.VerificationPending, resp.VerificationStatus)

	// Verify in database
	var savedPlayer model.Player
	err = db.First(&savedPlayer, resp.PlayerID).Error
	require.NoError(t, err)
	assert.Equal(t, req.Nickname, savedPlayer.Nickname)
	assert.Equal(t, req.Bio, savedPlayer.Bio)
	assert.Equal(t, req.Rank, savedPlayer.Rank)
}

func TestPlayerService_ApplyAsPlayer_AlreadyPlayer(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerRepo := player.NewPlayerRepository(db)
	userRepo := user.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	mockOrders := &mockOrderQuery{}
	mockTags := &mockPlayerTagRepo{tags: make(map[uint64][]string)}
	mockCacheInst := newMockCache()

	svc := playerservice.NewPlayerService(playerRepo, userRepo, gameRepo, mockOrders, reviewRepo, mockTags, mockCacheInst)

	// Create existing player
	testGame := CreateTestGame(t, db, "already_game")
	playerUser := CreateUniqueTestUser(t, db, "already_player")
	CreateTestPlayer(t, db, playerUser)

	// Try to apply again
	req := playerservice.ApplyPlayerRequest{
		Nickname:        "Another Player",
		MainGameID:      testGame.ID,
		Rank:            "Gold",
		HourlyRateCents: 3000,
	}

	_, err := svc.ApplyAsPlayer(ctx, playerUser.ID, req)
	assert.Error(t, err)
	assert.Equal(t, playerservice.ErrAlreadyPlayer, err)
}

func TestPlayerService_UpdatePlayerProfile(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerRepo := player.NewPlayerRepository(db)
	userRepo := user.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	mockOrders := &mockOrderQuery{}
	mockTags := &mockPlayerTagRepo{tags: make(map[uint64][]string)}
	mockCacheInst := newMockCache()

	svc := playerservice.NewPlayerService(playerRepo, userRepo, gameRepo, mockOrders, reviewRepo, mockTags, mockCacheInst)

	// Create test player
	playerUser := CreateUniqueTestUser(t, db, "update_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Update profile
	req := playerservice.UpdatePlayerProfileRequest{
		Nickname:        "Updated Nickname",
		Bio:             "Updated bio",
		Rank:            "Master",
		HourlyRateCents: 8000,
		Tags:            []string{"expert", "fast"},
	}

	err := svc.UpdatePlayerProfile(ctx, playerUser.ID, req)
	require.NoError(t, err)

	// Verify update
	var updatedPlayer model.Player
	err = db.First(&updatedPlayer, testPlayer.ID).Error
	require.NoError(t, err)
	assert.Equal(t, req.Nickname, updatedPlayer.Nickname)
	assert.Equal(t, req.Bio, updatedPlayer.Bio)
	assert.Equal(t, req.Rank, updatedPlayer.Rank)
	assert.Equal(t, req.HourlyRateCents, updatedPlayer.HourlyRateCents)
}

func TestPlayerService_SetPlayerOnlineStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerRepo := player.NewPlayerRepository(db)
	userRepo := user.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	mockOrders := &mockOrderQuery{}
	mockTags := &mockPlayerTagRepo{tags: make(map[uint64][]string)}
	mockCacheInst := newMockCache()

	svc := playerservice.NewPlayerService(playerRepo, userRepo, gameRepo, mockOrders, reviewRepo, mockTags, mockCacheInst)

	// Create test player
	playerUser := CreateUniqueTestUser(t, db, "online_player")
	CreateTestPlayer(t, db, playerUser)

	// Set online
	err := svc.SetPlayerOnlineStatus(ctx, playerUser.ID, true)
	require.NoError(t, err)

	// Set offline
	err = svc.SetPlayerOnlineStatus(ctx, playerUser.ID, false)
	require.NoError(t, err)
}

func TestPlayerService_GetPlayerProfile(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerRepo := player.NewPlayerRepository(db)
	userRepo := user.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	mockOrders := &mockOrderQuery{}
	mockTags := &mockPlayerTagRepo{tags: make(map[uint64][]string)}
	mockCacheInst := newMockCache()

	svc := playerservice.NewPlayerService(playerRepo, userRepo, gameRepo, mockOrders, reviewRepo, mockTags, mockCacheInst)

	// Create test player
	testGame := CreateTestGame(t, db, "profile_game")
	playerUser := CreateUniqueTestUser(t, db, "profile_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testPlayer.MainGameID = testGame.ID
	db.Save(testPlayer)

	// Get own profile
	resp, err := svc.GetPlayerProfile(ctx, playerUser.ID)
	require.NoError(t, err)
	assert.Equal(t, testPlayer.ID, resp.Player.ID)
}
