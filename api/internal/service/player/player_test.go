package player

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	repoiface "gamelink/internal/repository/interfaces"
)

// ============================================================================
// Mock Implementations
// ============================================================================

// MockPlayerRepository is a mock implementation of PlayerRepository
type MockPlayerRepository struct {
	mock.Mock
}

func (m *MockPlayerRepository) Create(ctx context.Context, player *model.Player) error {
	args := m.Called(ctx, player)
	if args.Error(0) == nil && player.ID == 0 {
		player.ID = 1
	}
	return args.Error(0)
}

func (m *MockPlayerRepository) Get(ctx context.Context, id uint64) (*model.Player, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Player), args.Error(1)
}

func (m *MockPlayerRepository) Update(ctx context.Context, player *model.Player) error {
	args := m.Called(ctx, player)
	return args.Error(0)
}

func (m *MockPlayerRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPlayerRepository) List(ctx context.Context) ([]model.Player, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Player), args.Error(1)
}

func (m *MockPlayerRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.Player), args.Get(1).(int64), args.Error(2)
}

func (m *MockPlayerRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, status *model.VerificationStatus) ([]model.Player, int64, error) {
	args := m.Called(ctx, page, pageSize, keyword, status)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.Player), args.Get(1).(int64), args.Error(2)
}

func (m *MockPlayerRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Player), args.Error(1)
}

func (m *MockPlayerRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.Player, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Player), args.Error(1)
}

func (m *MockPlayerRepository) BatchUpdateRank(ctx context.Context, ids []uint64, rank string) (int64, error) {
	args := m.Called(ctx, ids, rank)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPlayerRepository) BatchUpdateHourlyRate(ctx context.Context, ids []uint64, rateCents int64) (int64, error) {
	args := m.Called(ctx, ids, rateCents)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPlayerRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.VerificationStatus) (int64, error) {
	args := m.Called(ctx, ids, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPlayerRepository) BatchDelete(ctx context.Context, ids []uint64) (int64, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(int64), args.Error(1)
}

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context) ([]model.User, error)                   { return nil, nil }
func (m *MockUserRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *MockUserRepository) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *MockUserRepository) Count(ctx context.Context, opts repository.UserListOptions) (int, error) { return 0, nil }
func (m *MockUserRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.User, error)      { return nil, nil }
func (m *MockUserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error)     { return nil, nil }
func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error)     { return nil, nil }
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error)    { return nil, nil }
func (m *MockUserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error)    { return nil, nil }
func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error                   { return nil }
func (m *MockUserRepository) Delete(ctx context.Context, id uint64) error                          { return nil }
func (m *MockUserRepository) UpdatePassword(ctx context.Context, userID uint64, newPassword string) error { return nil }
func (m *MockUserRepository) GetByWeChatOpenID(ctx context.Context, openID string) (*model.User, error) { return nil, nil }
func (m *MockUserRepository) GetByWeChatUnionID(ctx context.Context, unionID string) (*model.User, error) { return nil, nil }

// MockGameRepository is a mock implementation of GameRepository
type MockGameRepository struct {
	mock.Mock
}

func (m *MockGameRepository) Get(ctx context.Context, id uint64) (*model.Game, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Game), args.Error(1)
}

func (m *MockGameRepository) Create(ctx context.Context, game *model.Game) error { return nil }
func (m *MockGameRepository) Update(ctx context.Context, game *model.Game) error { return nil }
func (m *MockGameRepository) Delete(ctx context.Context, id uint64) error        { return nil }
func (m *MockGameRepository) List(ctx context.Context) ([]model.Game, error)     { return nil, nil }
func (m *MockGameRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Game, int64, error) {
	return nil, 0, nil
}
func (m *MockGameRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string) ([]model.Game, int64, error) {
	return nil, 0, nil
}
func (m *MockGameRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.Game, error) { return nil, nil }
func (m *MockGameRepository) BatchDelete(ctx context.Context, ids []uint64) (int64, error)     { return 0, nil }
func (m *MockGameRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error) { return 0, nil }
func (m *MockGameRepository) BatchUpdateSortOrder(ctx context.Context, updates map[uint64]int) (int64, error) { return 0, nil }
func (m *MockGameRepository) BatchUpdateCategory(ctx context.Context, ids []uint64, category string) (int64, error) { return 0, nil }

// MockOrderQuery is a mock implementation of OrderQuery
type MockOrderQuery struct {
	mock.Mock
}

func (m *MockOrderQuery) List(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderQuery) Get(ctx context.Context, id uint64) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderQuery) GetByOrderNo(ctx context.Context, orderNo string) (*model.Order, error) {
	args := m.Called(ctx, orderNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// MockReviewRepository is a mock implementation of ReviewRepository
type MockReviewRepository struct {
	mock.Mock
}

func (m *MockReviewRepository) List(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.Review), args.Get(1).(int64), args.Error(2)
}

func (m *MockReviewRepository) Create(ctx context.Context, review *model.Review) error { return nil }
func (m *MockReviewRepository) Get(ctx context.Context, id uint64) (*model.Review, error) { return nil, nil }
func (m *MockReviewRepository) Update(ctx context.Context, review *model.Review) error { return nil }
func (m *MockReviewRepository) Delete(ctx context.Context, id uint64) error { return nil }
func (m *MockReviewRepository) GetByOrderID(ctx context.Context, orderID uint64) (*model.Review, error) { return nil, nil }
func (m *MockReviewRepository) ListPending(ctx context.Context, page, pageSize int) ([]model.Review, int64, error) { return nil, 0, nil }
func (m *MockReviewRepository) UpdateStatus(ctx context.Context, id uint64, status model.ReviewStatus, rejectionReason string) error { return nil }
func (m *MockReviewRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.ReviewStatus, rejectionReason string) error { return nil }
func (m *MockReviewRepository) GetStats(ctx context.Context) (repository.ReviewStats, error) { return repository.ReviewStats{}, nil }
func (m *MockReviewRepository) GetTrend(ctx context.Context, days int) ([]repository.DateValue, error) { return nil, nil }
func (m *MockReviewRepository) GetTopPlayersByReviewCount(ctx context.Context, limit int) ([]repository.PlayerReviewStats, error) { return nil, nil }
func (m *MockReviewRepository) GetTopPlayersByRating(ctx context.Context, limit int) ([]repository.PlayerReviewStats, error) { return nil, nil }
func (m *MockReviewRepository) GetGameStats(ctx context.Context) ([]repository.GameReviewStats, error) { return nil, nil }

// MockPlayerTagRepository is a mock implementation of PlayerTagRepository
type MockPlayerTagRepository struct {
	mock.Mock
}

func (m *MockPlayerTagRepository) GetTags(ctx context.Context, playerID uint64) ([]string, error) {
	args := m.Called(ctx, playerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockPlayerTagRepository) ReplaceTags(ctx context.Context, playerID uint64, tags []string) error {
	args := m.Called(ctx, playerID, tags)
	return args.Error(0)
}

// MockCache is a mock implementation of cache.Cache
type MockCache struct {
	mock.Mock
	data map[string]string
}

func NewMockCache() *MockCache {
	return &MockCache{data: make(map[string]string)}
}

func (m *MockCache) Get(ctx context.Context, key string) (string, bool, error) {
	if val, ok := m.data[key]; ok {
		return val, true, nil
	}
	return "", false, nil
}

func (m *MockCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	m.data[key] = value
	return nil
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	delete(m.data, key)
	return nil
}

func (m *MockCache) Close(ctx context.Context) error {
	return nil
}

func (m *MockCache) GetRedisClient() interface{} {
	return nil
}

// ============================================================================
// Tests
// ============================================================================

func TestNewPlayerService(t *testing.T) {
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	orders := &MockOrderQuery{}
	reviews := &MockReviewRepository{}
	playerTags := &MockPlayerTagRepository{}
	cache := NewMockCache()

	svc := NewPlayerService(players, users, games, orders, reviews, playerTags, cache)

	assert.NotNil(t, svc)
}

func TestPlayerService_ListPlayers(t *testing.T) {
	tests := []struct {
		name        string
		req         PlayerListRequest
		setupMock   func(*MockPlayerRepository, *MockUserRepository, *MockGameRepository, *MockOrderQuery)
		expectError bool
		expectCount int
	}{
		{
			name: "successful list players",
			req: PlayerListRequest{
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(players *MockPlayerRepository, users *MockUserRepository, games *MockGameRepository, orders *MockOrderQuery) {
				player := model.Player{
					UserID:             1,
					Nickname:           "TestPlayer",
					VerificationStatus: model.VerificationVerified,
					MainGameID:         1,
					HourlyRateCents:    5000,
				}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 10).Return([]model.Player{player}, int64(1), nil)

				user := &model.User{Name: "TestUser", AvatarURL: "https://example.com/avatar.jpg"}
				user.ID = 1
				users.On("Get", mock.Anything, uint64(1)).Return(user, nil)

				game := &model.Game{Name: "王者荣耀"}
				game.ID = 1
				games.On("Get", mock.Anything, uint64(1)).Return(game, nil)

				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return(
					[]model.Order{}, int64(5), nil,
				)
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name: "empty list",
			req: PlayerListRequest{
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(players *MockPlayerRepository, users *MockUserRepository, games *MockGameRepository, orders *MockOrderQuery) {
				players.On("ListPaged", mock.Anything, 1, 10).Return([]model.Player{}, int64(0), nil)
			},
			expectError: false,
			expectCount: 0,
		},
		{
			name: "filter by game ID",
			req: PlayerListRequest{
				GameID:   ptrUint64(1),
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(players *MockPlayerRepository, users *MockUserRepository, games *MockGameRepository, orders *MockOrderQuery) {
				player1 := model.Player{
					UserID:             1,
					Nickname:           "Player1",
					VerificationStatus: model.VerificationVerified,
					MainGameID:         1,
				}
				player1.ID = 1
				player2 := model.Player{
					UserID:             2,
					Nickname:           "Player2",
					VerificationStatus: model.VerificationVerified,
					MainGameID:         2, // Different game
				}
				player2.ID = 2
				players.On("ListPaged", mock.Anything, 1, 10).Return([]model.Player{player1, player2}, int64(2), nil)

				user := &model.User{Name: "TestUser"}
				user.ID = 1
				users.On("Get", mock.Anything, uint64(1)).Return(user, nil)

				game := &model.Game{Name: "王者荣耀"}
				game.ID = 1
				games.On("Get", mock.Anything, uint64(1)).Return(game, nil)

				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return(
					[]model.Order{}, int64(0), nil,
				)
			},
			expectError: false,
			expectCount: 1, // Only player1 matches
		},
		{
			name: "default pagination",
			req:  PlayerListRequest{}, // No page/pageSize specified
			setupMock: func(players *MockPlayerRepository, users *MockUserRepository, games *MockGameRepository, orders *MockOrderQuery) {
				players.On("ListPaged", mock.Anything, 1, 20).Return([]model.Player{}, int64(0), nil)
			},
			expectError: false,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			users := &MockUserRepository{}
			games := &MockGameRepository{}
			orders := &MockOrderQuery{}
			reviews := &MockReviewRepository{}
			playerTags := &MockPlayerTagRepository{}
			cache := NewMockCache()

			tt.setupMock(players, users, games, orders)

			svc := NewPlayerService(players, users, games, orders, reviews, playerTags, cache)
			resp, err := svc.ListPlayers(context.Background(), tt.req)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Len(t, resp.Players, tt.expectCount)
			}
		})
	}
}

func TestPlayerService_GetPlayerDetail(t *testing.T) {
	tests := []struct {
		name        string
		playerID    uint64
		setupMock   func(*MockPlayerRepository, *MockUserRepository, *MockGameRepository, *MockOrderQuery, *MockReviewRepository, *MockPlayerTagRepository)
		expectError bool
	}{
		{
			name:     "successful get player detail",
			playerID: 1,
			setupMock: func(players *MockPlayerRepository, users *MockUserRepository, games *MockGameRepository, orders *MockOrderQuery, reviews *MockReviewRepository, playerTags *MockPlayerTagRepository) {
				player := &model.Player{
					UserID:             1,
					Nickname:           "TestPlayer",
					VerificationStatus: model.VerificationVerified,
					MainGameID:         1,
				}
				player.ID = 1
				players.On("Get", mock.Anything, uint64(1)).Return(player, nil)

				user := &model.User{Name: "TestUser", AvatarURL: "https://example.com/avatar.jpg"}
				user.ID = 1
				users.On("Get", mock.Anything, uint64(1)).Return(user, nil)

				game := &model.Game{Name: "王者荣耀"}
				game.ID = 1
				games.On("Get", mock.Anything, uint64(1)).Return(game, nil)

				playerTags.On("GetTags", mock.Anything, uint64(1)).Return([]string{"高端局", "上分快"}, nil)

				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return(
					[]model.Order{}, int64(10), nil,
				)

				reviews.On("List", mock.Anything, mock.AnythingOfType("repository.ReviewListOptions")).Return(
					[]model.Review{}, int64(0), nil,
				)
			},
			expectError: false,
		},
		{
			name:     "player not found",
			playerID: 999,
			setupMock: func(players *MockPlayerRepository, users *MockUserRepository, games *MockGameRepository, orders *MockOrderQuery, reviews *MockReviewRepository, playerTags *MockPlayerTagRepository) {
				players.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
		{
			name:     "player not verified",
			playerID: 2,
			setupMock: func(players *MockPlayerRepository, users *MockUserRepository, games *MockGameRepository, orders *MockOrderQuery, reviews *MockReviewRepository, playerTags *MockPlayerTagRepository) {
				player := &model.Player{
					UserID:             2,
					Nickname:           "PendingPlayer",
					VerificationStatus: model.VerificationPending,
				}
				player.ID = 2
				players.On("Get", mock.Anything, uint64(2)).Return(player, nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			users := &MockUserRepository{}
			games := &MockGameRepository{}
			orders := &MockOrderQuery{}
			reviews := &MockReviewRepository{}
			playerTags := &MockPlayerTagRepository{}
			cache := NewMockCache()

			tt.setupMock(players, users, games, orders, reviews, playerTags)

			svc := NewPlayerService(players, users, games, orders, reviews, playerTags, cache)
			resp, err := svc.GetPlayerDetail(context.Background(), tt.playerID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.playerID, resp.Player.ID)
			}
		})
	}
}

// Helper function
func ptrUint64(v uint64) *uint64 {
	return &v
}

func TestPlayerService_ApplyAsPlayer(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint64
		req         ApplyPlayerRequest
		setupMock   func(*MockPlayerRepository, *MockUserRepository, *MockGameRepository, *MockPlayerTagRepository)
		expectError bool
	}{
		{
			name:   "successful apply as player",
			userID: 1,
			req: ApplyPlayerRequest{
				Nickname:        "NewPlayer",
				Bio:             "I'm a great player",
				MainGameID:      1,
				Rank:            "Diamond",
				HourlyRateCents: 5000,
				Tags:            []string{"高端局", "上分快"},
			},
			setupMock: func(players *MockPlayerRepository, users *MockUserRepository, games *MockGameRepository, playerTags *MockPlayerTagRepository) {
				user := &model.User{Name: "TestUser"}
				user.ID = 1
				users.On("Get", mock.Anything, uint64(1)).Return(user, nil)
				users.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

				// No existing player
				players.On("ListPaged", mock.Anything, 1, 1).Return([]model.Player{}, int64(0), nil)
				players.On("Create", mock.Anything, mock.AnythingOfType("*model.Player")).Return(nil)

				game := &model.Game{Name: "王者荣耀"}
				game.ID = 1
				games.On("Get", mock.Anything, uint64(1)).Return(game, nil)

				playerTags.On("ReplaceTags", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name:   "user not found",
			userID: 999,
			req: ApplyPlayerRequest{
				Nickname:        "NewPlayer",
				MainGameID:      1,
				Rank:            "Diamond",
				HourlyRateCents: 5000,
			},
			setupMock: func(players *MockPlayerRepository, users *MockUserRepository, games *MockGameRepository, playerTags *MockPlayerTagRepository) {
				users.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
		{
			name:   "already a player",
			userID: 1,
			req: ApplyPlayerRequest{
				Nickname:        "NewPlayer",
				MainGameID:      1,
				Rank:            "Diamond",
				HourlyRateCents: 5000,
			},
			setupMock: func(players *MockPlayerRepository, users *MockUserRepository, games *MockGameRepository, playerTags *MockPlayerTagRepository) {
				user := &model.User{Name: "TestUser"}
				user.ID = 1
				users.On("Get", mock.Anything, uint64(1)).Return(user, nil)

				// Already a player
				existingPlayer := model.Player{UserID: 1}
				existingPlayer.ID = 1
				players.On("ListPaged", mock.Anything, 1, 1).Return([]model.Player{existingPlayer}, int64(1), nil)
			},
			expectError: true,
		},
		{
			name:   "invalid game ID",
			userID: 1,
			req: ApplyPlayerRequest{
				Nickname:        "NewPlayer",
				MainGameID:      999,
				Rank:            "Diamond",
				HourlyRateCents: 5000,
			},
			setupMock: func(players *MockPlayerRepository, users *MockUserRepository, games *MockGameRepository, playerTags *MockPlayerTagRepository) {
				user := &model.User{Name: "TestUser"}
				user.ID = 1
				users.On("Get", mock.Anything, uint64(1)).Return(user, nil)

				players.On("ListPaged", mock.Anything, 1, 1).Return([]model.Player{}, int64(0), nil)

				games.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			users := &MockUserRepository{}
			games := &MockGameRepository{}
			orders := &MockOrderQuery{}
			reviews := &MockReviewRepository{}
			playerTags := &MockPlayerTagRepository{}
			cache := NewMockCache()

			tt.setupMock(players, users, games, playerTags)

			svc := NewPlayerService(players, users, games, orders, reviews, playerTags, cache)
			resp, err := svc.ApplyAsPlayer(context.Background(), tt.userID, tt.req)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, model.VerificationPending, resp.VerificationStatus)
			}
		})
	}
}

func TestPlayerService_SetPlayerOnlineStatus(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint64
		status      string
		setupMock   func(*MockPlayerRepository)
		expectError bool
	}{
		{
			name:   "set online status",
			userID: 1,
			status: "online",
			setupMock: func(players *MockPlayerRepository) {
				player := model.Player{UserID: 1}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)
			},
			expectError: false,
		},
		{
			name:   "set offline status",
			userID: 1,
			status: "offline",
			setupMock: func(players *MockPlayerRepository) {
				player := model.Player{UserID: 1}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)
			},
			expectError: false,
		},
		{
			name:   "player not found",
			userID: 999,
			status: "online",
			setupMock: func(players *MockPlayerRepository) {
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{}, int64(0), nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			users := &MockUserRepository{}
			games := &MockGameRepository{}
			orders := &MockOrderQuery{}
			reviews := &MockReviewRepository{}
			playerTags := &MockPlayerTagRepository{}
			cache := NewMockCache()

			tt.setupMock(players)

			svc := NewPlayerService(players, users, games, orders, reviews, playerTags, cache)
			err := svc.SetPlayerOnlineStatus(context.Background(), tt.userID, tt.status)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPlayerService_GetPlayerProfile(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint64
		setupMock   func(*MockPlayerRepository, *MockUserRepository, *MockGameRepository, *MockOrderQuery, *MockReviewRepository, *MockPlayerTagRepository)
		expectError bool
	}{
		{
			name:   "successful get profile",
			userID: 1,
			setupMock: func(players *MockPlayerRepository, users *MockUserRepository, games *MockGameRepository, orders *MockOrderQuery, reviews *MockReviewRepository, playerTags *MockPlayerTagRepository) {
				player := model.Player{
					UserID:             1,
					Nickname:           "TestPlayer",
					VerificationStatus: model.VerificationVerified,
					MainGameID:         1,
				}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)
				players.On("Get", mock.Anything, uint64(1)).Return(&player, nil)

				user := &model.User{Name: "TestUser", AvatarURL: "https://example.com/avatar.jpg"}
				user.ID = 1
				users.On("Get", mock.Anything, uint64(1)).Return(user, nil)

				game := &model.Game{Name: "王者荣耀"}
				game.ID = 1
				games.On("Get", mock.Anything, uint64(1)).Return(game, nil)

				playerTags.On("GetTags", mock.Anything, uint64(1)).Return([]string{"高端局"}, nil)

				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return(
					[]model.Order{}, int64(0), nil,
				)

				reviews.On("List", mock.Anything, mock.AnythingOfType("repository.ReviewListOptions")).Return(
					[]model.Review{}, int64(0), nil,
				)
			},
			expectError: false,
		},
		{
			name:   "player not found",
			userID: 999,
			setupMock: func(players *MockPlayerRepository, users *MockUserRepository, games *MockGameRepository, orders *MockOrderQuery, reviews *MockReviewRepository, playerTags *MockPlayerTagRepository) {
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{}, int64(0), nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			users := &MockUserRepository{}
			games := &MockGameRepository{}
			orders := &MockOrderQuery{}
			reviews := &MockReviewRepository{}
			playerTags := &MockPlayerTagRepository{}
			cache := NewMockCache()

			tt.setupMock(players, users, games, orders, reviews, playerTags)

			svc := NewPlayerService(players, users, games, orders, reviews, playerTags, cache)
			resp, err := svc.GetPlayerProfile(context.Background(), tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestParseUint64(t *testing.T) {
	tests := []struct {
		input    string
		expected uint64
		hasError bool
	}{
		{"123", 123, false},
		{"0", 0, false},
		{"999999", 999999, false},
		{"abc", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := parseUint64(tt.input)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestPlayerCardDTO_Structure(t *testing.T) {
	dto := PlayerCardDTO{
		ID:              1,
		UserID:          100,
		Nickname:        "TestPlayer",
		AvatarURL:       "https://example.com/avatar.jpg",
		Bio:             "Test bio",
		Rank:            "Diamond",
		RatingAverage:   4.5,
		RatingCount:     100,
		HourlyRateCents: 5000,
		MainGame:        "王者荣耀",
		IsOnline:        true,
		OrderCount:      50,
	}

	assert.Equal(t, uint64(1), dto.ID)
	assert.Equal(t, uint64(100), dto.UserID)
	assert.Equal(t, "TestPlayer", dto.Nickname)
	assert.Equal(t, "Diamond", dto.Rank)
	assert.Equal(t, float32(4.5), dto.RatingAverage)
	assert.Equal(t, int64(5000), dto.HourlyRateCents)
	assert.True(t, dto.IsOnline)
}

func TestPlayerDetailDTO_Structure(t *testing.T) {
	dto := PlayerDetailDTO{
		PlayerCardDTO: PlayerCardDTO{
			ID:       1,
			Nickname: "TestPlayer",
		},
		Tags:           []string{"高端局", "上分快"},
		GoodRatio:      0.95,
		AvgResponseMin: 15,
	}

	assert.Equal(t, uint64(1), dto.ID)
	assert.Len(t, dto.Tags, 2)
	assert.Equal(t, float32(0.95), dto.GoodRatio)
	assert.Equal(t, 15, dto.AvgResponseMin)
}

func TestPlayerErrors(t *testing.T) {
	assert.NotNil(t, ErrNotFound)
	assert.NotNil(t, ErrValidation)
	assert.NotNil(t, ErrPlayerNotVerified)
	assert.NotNil(t, ErrAlreadyPlayer)
}

func TestPlayerService_UpdatePlayerProfile(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint64
		req         UpdatePlayerProfileRequest
		setupMock   func(*MockPlayerRepository, *MockPlayerTagRepository)
		expectError bool
	}{
		{
			name:   "successful update profile",
			userID: 1,
			req: UpdatePlayerProfileRequest{
				Nickname:        "UpdatedNickname",
				Bio:             "Updated bio",
				Rank:            "Master",
				HourlyRateCents: 6000,
				Tags:            []string{"新标签"},
			},
			setupMock: func(players *MockPlayerRepository, playerTags *MockPlayerTagRepository) {
				player := model.Player{
					UserID:             1,
					Nickname:           "OldNickname",
					Bio:                "Old bio",
					Rank:               "Diamond",
					HourlyRateCents:    5000,
					VerificationStatus: model.VerificationVerified,
				}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)
				players.On("Update", mock.Anything, mock.AnythingOfType("*model.Player")).Return(nil)
				playerTags.On("ReplaceTags", mock.Anything, uint64(1), []string{"新标签"}).Return(nil)
			},
			expectError: false,
		},
		{
			name:   "update profile without rank and rate",
			userID: 1,
			req: UpdatePlayerProfileRequest{
				Nickname: "UpdatedNickname",
				Bio:      "Updated bio",
			},
			setupMock: func(players *MockPlayerRepository, playerTags *MockPlayerTagRepository) {
				player := model.Player{
					UserID:             1,
					Nickname:           "OldNickname",
					Rank:               "Diamond",
					HourlyRateCents:    5000,
					VerificationStatus: model.VerificationVerified,
				}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)
				players.On("Update", mock.Anything, mock.AnythingOfType("*model.Player")).Return(nil)
			},
			expectError: false,
		},
		{
			name:   "player not found",
			userID: 999,
			req: UpdatePlayerProfileRequest{
				Nickname: "UpdatedNickname",
			},
			setupMock: func(players *MockPlayerRepository, playerTags *MockPlayerTagRepository) {
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{}, int64(0), nil)
			},
			expectError: true,
		},
		{
			name:   "update error",
			userID: 1,
			req: UpdatePlayerProfileRequest{
				Nickname: "UpdatedNickname",
			},
			setupMock: func(players *MockPlayerRepository, playerTags *MockPlayerTagRepository) {
				player := model.Player{UserID: 1}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)
				players.On("Update", mock.Anything, mock.AnythingOfType("*model.Player")).Return(repository.ErrNotFound)
			},
			expectError: true,
		},
		{
			name:   "list error",
			userID: 1,
			req: UpdatePlayerProfileRequest{
				Nickname: "UpdatedNickname",
			},
			setupMock: func(players *MockPlayerRepository, playerTags *MockPlayerTagRepository) {
				players.On("ListPaged", mock.Anything, 1, 100).Return(nil, int64(0), repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			users := &MockUserRepository{}
			games := &MockGameRepository{}
			orders := &MockOrderQuery{}
			reviews := &MockReviewRepository{}
			playerTags := &MockPlayerTagRepository{}
			cache := NewMockCache()

			tt.setupMock(players, playerTags)

			svc := NewPlayerService(players, users, games, orders, reviews, playerTags, cache)
			err := svc.UpdatePlayerProfile(context.Background(), tt.userID, tt.req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPlayerService_GetPlayerOnlineStatusByUserID(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint64
		setupMock    func(*MockPlayerRepository)
		setupCache   func(*MockCache)
		expectError  bool
		expectStatus string
	}{
		{
			name:   "get online status - online",
			userID: 1,
			setupMock: func(players *MockPlayerRepository) {
				player := model.Player{UserID: 1}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)
			},
			setupCache: func(cache *MockCache) {
				cache.data["player:online:1"] = "online"
			},
			expectError:  false,
			expectStatus: "online",
		},
		{
			name:   "get online status - offline",
			userID: 1,
			setupMock: func(players *MockPlayerRepository) {
				player := model.Player{UserID: 1}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)
			},
			setupCache:   func(cache *MockCache) {},
			expectError:  false,
			expectStatus: "offline",
		},
		{
			name:   "player not found",
			userID: 999,
			setupMock: func(players *MockPlayerRepository) {
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{}, int64(0), nil)
			},
			setupCache:  func(cache *MockCache) {},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			users := &MockUserRepository{}
			games := &MockGameRepository{}
			orders := &MockOrderQuery{}
			reviews := &MockReviewRepository{}
			playerTags := &MockPlayerTagRepository{}
			cache := NewMockCache()

			tt.setupMock(players)
			tt.setupCache(cache)

			svc := NewPlayerService(players, users, games, orders, reviews, playerTags, cache)
			status, err := svc.GetPlayerOnlineStatusByUserID(context.Background(), tt.userID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectStatus, status)
			}
		})
	}
}

func TestPlayerService_getPlayerIDByUserID_CacheHit(t *testing.T) {
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	orders := &MockOrderQuery{}
	reviews := &MockReviewRepository{}
	playerTags := &MockPlayerTagRepository{}
	cache := NewMockCache()

	// Set cache value
	cache.data["player:user_id:1"] = "100"

	svc := NewPlayerService(players, users, games, orders, reviews, playerTags, cache)
	playerID, err := svc.getPlayerIDByUserID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, uint64(100), playerID)
	// Verify no database call was made
	players.AssertNotCalled(t, "ListPaged")
}

func TestPlayerService_getPlayerIDByUserID_CacheMiss(t *testing.T) {
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	orders := &MockOrderQuery{}
	reviews := &MockReviewRepository{}
	playerTags := &MockPlayerTagRepository{}
	cache := NewMockCache()

	player := model.Player{UserID: 1}
	player.ID = 100
	players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)

	svc := NewPlayerService(players, users, games, orders, reviews, playerTags, cache)
	playerID, err := svc.getPlayerIDByUserID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, uint64(100), playerID)
	// Verify cache was set
	assert.Equal(t, "100", cache.data["player:user_id:1"])
}

func TestPlayerService_calculateGoodRatio(t *testing.T) {
	tests := []struct {
		name        string
		playerID    uint64
		setupMock   func(*MockReviewRepository)
		expectRatio float32
	}{
		{
			name:     "all good reviews",
			playerID: 1,
			setupMock: func(reviews *MockReviewRepository) {
				reviewList := []model.Review{
					{Score: 5},
					{Score: 4},
					{Score: 5},
				}
				reviews.On("List", mock.Anything, mock.AnythingOfType("repository.ReviewListOptions")).Return(reviewList, int64(3), nil)
			},
			expectRatio: 1.0,
		},
		{
			name:     "mixed reviews",
			playerID: 1,
			setupMock: func(reviews *MockReviewRepository) {
				reviewList := []model.Review{
					{Score: 5},
					{Score: 4},
					{Score: 3},
					{Score: 2},
				}
				reviews.On("List", mock.Anything, mock.AnythingOfType("repository.ReviewListOptions")).Return(reviewList, int64(4), nil)
			},
			expectRatio: 0.5, // 2 good (4,5) out of 4
		},
		{
			name:     "no reviews",
			playerID: 1,
			setupMock: func(reviews *MockReviewRepository) {
				reviews.On("List", mock.Anything, mock.AnythingOfType("repository.ReviewListOptions")).Return([]model.Review{}, int64(0), nil)
			},
			expectRatio: 0.0,
		},
		{
			name:     "review list error",
			playerID: 1,
			setupMock: func(reviews *MockReviewRepository) {
				reviews.On("List", mock.Anything, mock.AnythingOfType("repository.ReviewListOptions")).Return(nil, int64(0), repository.ErrNotFound)
			},
			expectRatio: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			users := &MockUserRepository{}
			games := &MockGameRepository{}
			orders := &MockOrderQuery{}
			reviews := &MockReviewRepository{}
			playerTags := &MockPlayerTagRepository{}
			cache := NewMockCache()

			tt.setupMock(reviews)

			svc := NewPlayerService(players, users, games, orders, reviews, playerTags, cache)
			ratio := svc.calculateGoodRatio(context.Background(), tt.playerID)

			assert.Equal(t, tt.expectRatio, ratio)
		})
	}
}

func TestPlayerService_calculateAvgResponseTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		playerID   uint64
		setupMock  func(*MockOrderQuery)
		expectTime int
	}{
		{
			name:     "orders with response time",
			playerID: 1,
			setupMock: func(orders *MockOrderQuery) {
				startedAt1 := now.Add(10 * time.Minute)
				startedAt2 := now.Add(20 * time.Minute)
				orderList := []model.Order{
					{StartedAt: &startedAt1},
					{StartedAt: &startedAt2},
				}
				orderList[0].CreatedAt = now
				orderList[1].CreatedAt = now
				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return(orderList, int64(2), nil)
			},
			expectTime: 15, // (10 + 20) / 2
		},
		{
			name:     "no orders",
			playerID: 1,
			setupMock: func(orders *MockOrderQuery) {
				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return([]model.Order{}, int64(0), nil)
			},
			expectTime: 30, // default
		},
		{
			name:     "orders without started time",
			playerID: 1,
			setupMock: func(orders *MockOrderQuery) {
				orderList := []model.Order{
					{StartedAt: nil},
					{StartedAt: nil},
				}
				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return(orderList, int64(2), nil)
			},
			expectTime: 30, // default when no valid orders
		},
		{
			name:     "order list error",
			playerID: 1,
			setupMock: func(orders *MockOrderQuery) {
				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return(nil, int64(0), repository.ErrNotFound)
			},
			expectTime: 30, // default on error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			users := &MockUserRepository{}
			games := &MockGameRepository{}
			orders := &MockOrderQuery{}
			reviews := &MockReviewRepository{}
			playerTags := &MockPlayerTagRepository{}
			cache := NewMockCache()

			tt.setupMock(orders)

			svc := NewPlayerService(players, users, games, orders, reviews, playerTags, cache)
			avgTime := svc.calculateAvgResponseTime(context.Background(), tt.playerID)

			assert.Equal(t, tt.expectTime, avgTime)
		})
	}
}

func TestPlayerService_calculateRepeatRate(t *testing.T) {
	tests := []struct {
		name       string
		playerID   uint64
		setupMock  func(*MockOrderQuery)
		expectRate float32
	}{
		{
			name:     "with repeat users",
			playerID: 1,
			setupMock: func(orders *MockOrderQuery) {
				orderList := []model.Order{
					{UserID: 1},
					{UserID: 1}, // repeat
					{UserID: 2},
					{UserID: 3},
				}
				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return(orderList, int64(4), nil)
			},
			expectRate: float32(1) / float32(3), // 1 repeat user out of 3 unique users
		},
		{
			name:     "no repeat users",
			playerID: 1,
			setupMock: func(orders *MockOrderQuery) {
				orderList := []model.Order{
					{UserID: 1},
					{UserID: 2},
					{UserID: 3},
				}
				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return(orderList, int64(3), nil)
			},
			expectRate: 0.0,
		},
		{
			name:     "all repeat users",
			playerID: 1,
			setupMock: func(orders *MockOrderQuery) {
				orderList := []model.Order{
					{UserID: 1},
					{UserID: 1},
					{UserID: 2},
					{UserID: 2},
				}
				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return(orderList, int64(4), nil)
			},
			expectRate: 1.0, // 2 repeat users out of 2 unique users
		},
		{
			name:     "no orders",
			playerID: 1,
			setupMock: func(orders *MockOrderQuery) {
				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return([]model.Order{}, int64(0), nil)
			},
			expectRate: 0.0,
		},
		{
			name:     "order list error",
			playerID: 1,
			setupMock: func(orders *MockOrderQuery) {
				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return(nil, int64(0), repository.ErrNotFound)
			},
			expectRate: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			users := &MockUserRepository{}
			games := &MockGameRepository{}
			orders := &MockOrderQuery{}
			reviews := &MockReviewRepository{}
			playerTags := &MockPlayerTagRepository{}
			cache := NewMockCache()

			tt.setupMock(orders)

			svc := NewPlayerService(players, users, games, orders, reviews, playerTags, cache)
			rate := svc.calculateRepeatRate(context.Background(), tt.playerID)

			assert.Equal(t, tt.expectRate, rate)
		})
	}
}

func TestPlayerService_getPlayerReviews(t *testing.T) {
	tests := []struct {
		name        string
		playerID    uint64
		limit       int
		setupMock   func(*MockReviewRepository, *MockUserRepository)
		expectError bool
		expectCount int
	}{
		{
			name:     "successful get reviews",
			playerID: 1,
			limit:    5,
			setupMock: func(reviews *MockReviewRepository, users *MockUserRepository) {
				review := model.Review{
					UserID:  1,
					Score:   5,
					Content: "Great player!",
				}
				review.ID = 1
				review.CreatedAt = time.Now()
				reviews.On("List", mock.Anything, mock.AnythingOfType("repository.ReviewListOptions")).Return([]model.Review{review}, int64(1), nil)

				user := &model.User{Name: "TestUser", AvatarURL: "https://example.com/avatar.jpg"}
				user.ID = 1
				users.On("Get", mock.Anything, uint64(1)).Return(user, nil)
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name:     "review list error",
			playerID: 1,
			limit:    5,
			setupMock: func(reviews *MockReviewRepository, users *MockUserRepository) {
				reviews.On("List", mock.Anything, mock.AnythingOfType("repository.ReviewListOptions")).Return(nil, int64(0), repository.ErrNotFound)
			},
			expectError: true,
			expectCount: 0,
		},
		{
			name:     "user not found - skip review",
			playerID: 1,
			limit:    5,
			setupMock: func(reviews *MockReviewRepository, users *MockUserRepository) {
				review := model.Review{
					UserID:  999,
					Score:   5,
					Content: "Great player!",
				}
				review.ID = 1
				review.CreatedAt = time.Now()
				reviews.On("List", mock.Anything, mock.AnythingOfType("repository.ReviewListOptions")).Return([]model.Review{review}, int64(1), nil)
				users.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: false,
			expectCount: 0, // Review skipped because user not found
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			users := &MockUserRepository{}
			games := &MockGameRepository{}
			orders := &MockOrderQuery{}
			reviews := &MockReviewRepository{}
			playerTags := &MockPlayerTagRepository{}
			cache := NewMockCache()

			tt.setupMock(reviews, users)

			svc := NewPlayerService(players, users, games, orders, reviews, playerTags, cache)
			result, err := svc.getPlayerReviews(context.Background(), tt.playerID, tt.limit)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectCount)
			}
		})
	}
}

func TestPlayerService_getPlayerStats(t *testing.T) {
	tests := []struct {
		name        string
		playerID    uint64
		setupMock   func(*MockOrderQuery)
		expectError bool
	}{
		{
			name:     "successful get stats",
			playerID: 1,
			setupMock: func(orders *MockOrderQuery) {
				// First call for total orders
				orders.On("List", mock.Anything, mock.MatchedBy(func(opts repoiface.OrderListOptions) bool {
					return opts.Statuses == nil
				})).Return([]model.Order{}, int64(10), nil).Once()

				// Second call for completed orders
				orders.On("List", mock.Anything, mock.MatchedBy(func(opts repoiface.OrderListOptions) bool {
					return len(opts.Statuses) > 0 && opts.Statuses[0] == model.OrderStatusCompleted
				})).Return([]model.Order{}, int64(8), nil).Once()

				// Third call for repeat rate calculation
				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return([]model.Order{}, int64(0), nil)
			},
			expectError: false,
		},
		{
			name:     "total orders error",
			playerID: 1,
			setupMock: func(orders *MockOrderQuery) {
				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return(nil, int64(0), repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			users := &MockUserRepository{}
			games := &MockGameRepository{}
			orders := &MockOrderQuery{}
			reviews := &MockReviewRepository{}
			playerTags := &MockPlayerTagRepository{}
			cache := NewMockCache()

			tt.setupMock(orders)

			svc := NewPlayerService(players, users, games, orders, reviews, playerTags, cache)
			stats, err := svc.getPlayerStats(context.Background(), tt.playerID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, stats)
			}
		})
	}
}

func TestPlayerService_getOnlineStatusKey(t *testing.T) {
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	orders := &MockOrderQuery{}
	reviews := &MockReviewRepository{}
	playerTags := &MockPlayerTagRepository{}
	cache := NewMockCache()

	svc := NewPlayerService(players, users, games, orders, reviews, playerTags, cache)

	key := svc.getOnlineStatusKey(123)
	assert.Equal(t, "player:online:123", key)
}

func TestPlayerService_getPlayerOnlineStatus(t *testing.T) {
	tests := []struct {
		name       string
		playerID   uint64
		setupCache func(*MockCache)
		expectBool bool
	}{
		{
			name:     "online",
			playerID: 1,
			setupCache: func(cache *MockCache) {
				cache.data["player:online:1"] = "online"
			},
			expectBool: true,
		},
		{
			name:     "busy",
			playerID: 1,
			setupCache: func(cache *MockCache) {
				cache.data["player:online:1"] = "busy"
			},
			expectBool: true,
		},
		{
			name:       "offline - no cache",
			playerID:   1,
			setupCache: func(cache *MockCache) {},
			expectBool: false,
		},
		{
			name:     "offline - explicit",
			playerID: 1,
			setupCache: func(cache *MockCache) {
				cache.data["player:online:1"] = "offline"
			},
			expectBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			users := &MockUserRepository{}
			games := &MockGameRepository{}
			orders := &MockOrderQuery{}
			reviews := &MockReviewRepository{}
			playerTags := &MockPlayerTagRepository{}
			cache := NewMockCache()

			tt.setupCache(cache)

			svc := NewPlayerService(players, users, games, orders, reviews, playerTags, cache)
			result := svc.getPlayerOnlineStatus(context.Background(), tt.playerID)

			assert.Equal(t, tt.expectBool, result)
		})
	}
}

func TestPlayerStatsDTO_Structure(t *testing.T) {
	stats := PlayerStatsDTO{
		TotalOrders:     100,
		CompletedOrders: 90,
		RepeatRate:      0.35,
	}

	assert.Equal(t, int64(100), stats.TotalOrders)
	assert.Equal(t, int64(90), stats.CompletedOrders)
	assert.Equal(t, float32(0.35), stats.RepeatRate)
}

func TestReviewDTO_Structure(t *testing.T) {
	dto := ReviewDTO{
		ID:            1,
		UserNickname:  "TestUser",
		UserAvatarURL: "https://example.com/avatar.jpg",
		Rating:        5,
		Comment:       "Great service!",
		CreatedAt:     "2024-01-15 10:30:00",
	}

	assert.Equal(t, uint64(1), dto.ID)
	assert.Equal(t, "TestUser", dto.UserNickname)
	assert.Equal(t, 5, dto.Rating)
	assert.Equal(t, "Great service!", dto.Comment)
}

func TestPlayerListRequest_Structure(t *testing.T) {
	gameID := uint64(1)
	minPrice := int64(1000)
	maxPrice := int64(5000)
	minRating := float32(4.0)

	req := PlayerListRequest{
		GameID:     &gameID,
		MinPrice:   &minPrice,
		MaxPrice:   &maxPrice,
		MinRating:  &minRating,
		OnlineOnly: true,
		SortBy:     "rating",
		Page:       1,
		PageSize:   20,
	}

	assert.Equal(t, uint64(1), *req.GameID)
	assert.Equal(t, int64(1000), *req.MinPrice)
	assert.Equal(t, int64(5000), *req.MaxPrice)
	assert.Equal(t, float32(4.0), *req.MinRating)
	assert.True(t, req.OnlineOnly)
	assert.Equal(t, "rating", req.SortBy)
}

func TestApplyPlayerRequest_Structure(t *testing.T) {
	req := ApplyPlayerRequest{
		Nickname:        "NewPlayer",
		Bio:             "I'm a great player",
		MainGameID:      1,
		Rank:            "Diamond",
		HourlyRateCents: 5000,
		Tags:            []string{"高端局", "上分快"},
		ProofImages:     []string{"https://example.com/proof1.jpg"},
	}

	assert.Equal(t, "NewPlayer", req.Nickname)
	assert.Equal(t, uint64(1), req.MainGameID)
	assert.Equal(t, int64(5000), req.HourlyRateCents)
	assert.Len(t, req.Tags, 2)
	assert.Len(t, req.ProofImages, 1)
}

func TestUpdatePlayerProfileRequest_Structure(t *testing.T) {
	req := UpdatePlayerProfileRequest{
		Nickname:        "UpdatedNickname",
		Bio:             "Updated bio",
		Rank:            "Master",
		HourlyRateCents: 6000,
		Tags:            []string{"新标签"},
	}

	assert.Equal(t, "UpdatedNickname", req.Nickname)
	assert.Equal(t, "Updated bio", req.Bio)
	assert.Equal(t, "Master", req.Rank)
	assert.Equal(t, int64(6000), req.HourlyRateCents)
	assert.Len(t, req.Tags, 1)
}

func TestSetPlayerStatusRequest_Structure(t *testing.T) {
	req := SetPlayerStatusRequest{
		Online: true,
	}

	assert.True(t, req.Online)
}
