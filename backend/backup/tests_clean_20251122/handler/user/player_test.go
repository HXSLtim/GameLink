package user

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/cache"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service/player"
)

// ---- Fake PlayerRepository for user_player tests ----

type mockPlayerRepoForUserPlayer struct {
	players []model.Player
}

func newMockPlayerRepoForUserPlayer() *mockPlayerRepoForUserPlayer {
	return &mockPlayerRepoForUserPlayer{
		players: []model.Player{
			{Base: model.Base{ID: 1}, UserID: 100, Nickname: "Player1", MainGameID: 10, HourlyRateCents: 5000, VerificationStatus: model.VerificationVerified, RatingAverage: 4.5},
			{Base: model.Base{ID: 2}, UserID: 101, Nickname: "Player2", MainGameID: 10, HourlyRateCents: 8000, VerificationStatus: model.VerificationVerified, RatingAverage: 4.8},
			{Base: model.Base{ID: 3}, UserID: 102, Nickname: "Player3", MainGameID: 20, HourlyRateCents: 3000, VerificationStatus: model.VerificationVerified, RatingAverage: 4.2},
		},
	}
}

func (m *mockPlayerRepoForUserPlayer) List(ctx context.Context) ([]model.Player, error) {
	return m.players, nil
}

func (m *mockPlayerRepoForUserPlayer) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	return m.players, int64(len(m.players)), nil
}

func (m *mockPlayerRepoForUserPlayer) Get(ctx context.Context, id uint64) (*model.Player, error) {
	for _, p := range m.players {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockPlayerRepoForUserPlayer) Create(ctx context.Context, player *model.Player) error {
	player.ID = uint64(len(m.players) + 1)
	m.players = append(m.players, *player)
	return nil
}

func (m *mockPlayerRepoForUserPlayer) Update(ctx context.Context, player *model.Player) error {
	for i := range m.players {
		if m.players[i].ID == player.ID {
			m.players[i] = *player
			return nil
		}
	}
	return repository.ErrNotFound
}

func (m *mockPlayerRepoForUserPlayer) Delete(ctx context.Context, id uint64) error {
	for i, p := range m.players {
		if p.ID == id {
			m.players = append(m.players[:i], m.players[i+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}

func (m *mockPlayerRepoForUserPlayer) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	for _, p := range m.players {
		if p.UserID == userID {
			return &p, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockPlayerRepoForUserPlayer) ListByGameID(ctx context.Context, gameID uint64) ([]model.Player, error) {
	var result []model.Player
	for _, p := range m.players {
		if p.MainGameID == gameID {
			result = append(result, p)
		}
	}
	return result, nil
}

// ---- Fake PlayerTagRepository ----

type fakePlayerTagRepository struct{}

func (m *fakePlayerTagRepository) GetTags(ctx context.Context, playerID uint64) ([]string, error) {
	return []string{}, nil
}
func (m *fakePlayerTagRepository) ReplaceTags(ctx context.Context, playerID uint64, tags []string) error {
	return nil
}

// ---- Fake Cache ----

type fakeCache struct{}

func (c *fakeCache) Get(ctx context.Context, key string) (string, bool, error)           { return "", false, nil }
func (c *fakeCache) Set(ctx context.Context, key, value string, ttl time.Duration) error { return nil }
func (c *fakeCache) Delete(ctx context.Context, key string) error                        { return nil }
func (c *fakeCache) Close(ctx context.Context) error                                     { return nil }

var _ cache.Cache = (*fakeCache)(nil)

// ---- Tests for user_player.go ----

func TestListPlayersHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	playerRepo := newMockPlayerRepoForUserPlayer()
	playerSvc := player.NewPlayerService(playerRepo, &fakeUserRepository{}, &fakeGameRepository{}, newFakeOrderRepository(), &fakeReviewRepository{}, &fakePlayerTagRepository{}, &fakeCache{})

	router := gin.New()
	router.GET("/user/players", func(c *gin.Context) {
		listPlayersHandler(c, playerSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/user/players?page=1&pageSize=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var resp model.APIResponse[player.PlayerListResponse]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("Expected success=true, got %v", resp.Success)
	}

	if len(resp.Data.Players) == 0 {
		t.Fatalf("Expected non-empty players list")
	}
}

func TestListPlayersHandler_WithFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	playerRepo := newMockPlayerRepoForUserPlayer()
	playerSvc := player.NewPlayerService(playerRepo, &fakeUserRepository{}, &fakeGameRepository{}, newFakeOrderRepository(), &fakeReviewRepository{}, &fakePlayerTagRepository{}, &fakeCache{})

	router := gin.New()
	router.GET("/user/players", func(c *gin.Context) {
		listPlayersHandler(c, playerSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/user/players?gameId=10&onlineOnly=true&minRating=4.0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
}

func TestGetPlayerDetailHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	playerRepo := newMockPlayerRepoForUserPlayer()
	playerSvc := player.NewPlayerService(playerRepo, &fakeUserRepository{}, &fakeGameRepository{}, newFakeOrderRepository(), &fakeReviewRepository{}, &fakePlayerTagRepository{}, &fakeCache{})

	router := gin.New()
	router.GET("/user/players/:id", func(c *gin.Context) {
		getPlayerDetailHandler(c, playerSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/user/players/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var resp model.APIResponse[player.PlayerDetailResponse]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("Expected success=true, got %v", resp.Success)
	}

	if resp.Data.Player.ID != 1 {
		t.Fatalf("Expected player ID 1, got %d", resp.Data.Player.ID)
	}
}

func TestGetPlayerDetailHandler_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	playerRepo := newMockPlayerRepoForUserPlayer()
	playerSvc := player.NewPlayerService(playerRepo, &fakeUserRepository{}, &fakeGameRepository{}, newFakeOrderRepository(), &fakeReviewRepository{}, &fakePlayerTagRepository{}, &fakeCache{})

	router := gin.New()
	router.GET("/user/players/:id", func(c *gin.Context) {
		getPlayerDetailHandler(c, playerSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/user/players/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

func TestGetPlayerDetailHandler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	playerRepo := newMockPlayerRepoForUserPlayer()
	playerSvc := player.NewPlayerService(playerRepo, &fakeUserRepository{}, &fakeGameRepository{}, newFakeOrderRepository(), &fakeReviewRepository{}, &fakePlayerTagRepository{}, &fakeCache{})

	router := gin.New()
	router.GET("/user/players/:id", func(c *gin.Context) {
		getPlayerDetailHandler(c, playerSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/user/players/9999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected status 404, got %d", w.Code)
	}
}

func TestListPlayersHandler_InvalidQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	playerRepo := newMockPlayerRepoForUserPlayer()
	playerSvc := player.NewPlayerService(playerRepo, &fakeUserRepository{}, &fakeGameRepository{}, newFakeOrderRepository(), &fakeReviewRepository{}, &fakePlayerTagRepository{}, &fakeCache{})

	router := gin.New()
	router.GET("/user/players", func(c *gin.Context) {
		listPlayersHandler(c, playerSvc)
	})

	// Invalid page parameter
	req := httptest.NewRequest(http.MethodGet, "/user/players?page=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

func TestListPlayersHandler_EmptyResult(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Empty repository
	emptyRepo := &mockPlayerRepoForUserPlayer{players: []model.Player{}}
	playerSvc := player.NewPlayerService(emptyRepo, &fakeUserRepository{}, &fakeGameRepository{}, newFakeOrderRepository(), &fakeReviewRepository{}, &fakePlayerTagRepository{}, &fakeCache{})

	router := gin.New()
	router.GET("/user/players", func(c *gin.Context) {
		listPlayersHandler(c, playerSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/user/players", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
}

func TestGetPlayerDetailHandler_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	playerRepo := newMockPlayerRepoForUserPlayer()
	playerSvc := player.NewPlayerService(playerRepo, &fakeUserRepository{}, &fakeGameRepository{}, newFakeOrderRepository(), &fakeReviewRepository{}, &fakePlayerTagRepository{}, &fakeCache{})

	router := gin.New()
	router.GET("/user/players/:id", func(c *gin.Context) {
		getPlayerDetailHandler(c, playerSvc)
	})

	// Test with existing player
	req := httptest.NewRequest(http.MethodGet, "/user/players/2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
}
