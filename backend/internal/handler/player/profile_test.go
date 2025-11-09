package player

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service/player"
)

// Fake repositories for profile tests (reuse from other test files where available)
type fakePlayerTagRepository struct{}

func (m *fakePlayerTagRepository) GetTags(ctx context.Context, playerID uint64) ([]string, error) {
	return []string{}, nil
}
func (m *fakePlayerTagRepository) ReplaceTags(ctx context.Context, playerID uint64, tags []string) error {
	return nil
}

type fakeCache struct{}

func (f *fakeCache) Get(ctx context.Context, key string) (value string, ok bool, err error) {
	return "", false, nil
}
func (f *fakeCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return nil
}
func (f *fakeCache) Delete(ctx context.Context, key string) error { return nil }
func (f *fakeCache) Close(ctx context.Context) error              { return nil }

type mockPlayerRepoForProfile struct {
	players map[uint64]*model.Player
}

func newMockPlayerRepoForUserPlayer() *mockPlayerRepoForProfile {
	return &mockPlayerRepoForProfile{
		players: map[uint64]*model.Player{
			1: {Base: model.Base{ID: 1}, UserID: 100, Nickname: "ExistingPlayer", VerificationStatus: model.VerificationVerified},
		},
	}
}

func (m *mockPlayerRepoForProfile) List(ctx context.Context) ([]model.Player, error) {
	var result []model.Player
	for _, p := range m.players {
		result = append(result, *p)
	}
	return result, nil
}
func (m *mockPlayerRepoForProfile) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	var result []model.Player
	for _, p := range m.players {
		result = append(result, *p)
	}
	return result, int64(len(result)), nil
}
func (m *mockPlayerRepoForProfile) Get(ctx context.Context, id uint64) (*model.Player, error) {
	if p, ok := m.players[id]; ok {
		return p, nil
	}
	return nil, repository.ErrNotFound
}
func (m *mockPlayerRepoForProfile) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	for _, p := range m.players {
		if p.UserID == userID {
			return p, nil
		}
	}
	return nil, repository.ErrNotFound
}
func (m *mockPlayerRepoForProfile) Create(ctx context.Context, player *model.Player) error {
	player.ID = uint64(len(m.players) + 1)
	m.players[player.ID] = player
	return nil
}
func (m *mockPlayerRepoForProfile) Update(ctx context.Context, player *model.Player) error {
	m.players[player.ID] = player
	return nil
}
func (m *mockPlayerRepoForProfile) Delete(ctx context.Context, id uint64) error {
	delete(m.players, id)
	return nil
}
func (m *mockPlayerRepoForProfile) ListByGameID(ctx context.Context, gameID uint64) ([]model.Player, error) {
	return []model.Player{}, nil
}

// ---- Tests for player_profile.go ----

func TestApplyAsPlayerHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	playerRepo := newMockPlayerRepoForUserPlayer()
	playerSvc := player.NewPlayerService(playerRepo, &fakeUserRepository{}, &fakeGameRepository{}, newFakeOrderRepositoryForEarnings(), &fakeReviewRepository{}, &fakePlayerTagRepository{}, &fakeCache{})

	router := gin.New()
	router.POST("/player/apply", func(c *gin.Context) {
		// Simulate authenticated user
		c.Set("user_id", uint64(500))
		applyAsPlayerHandler(c, playerSvc)
	})

	reqBody := player.ApplyPlayerRequest{
		Nickname:        "NewPlayer",
		Bio:             "I'm a new player",
		MainGameID:      10,
		Rank:            "Diamond",
		HourlyRateCents: 6000,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/player/apply", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp model.APIResponse[player.ApplyPlayerResponse]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("Expected success=true, got %v", resp.Success)
	}
}

func TestApplyAsPlayerHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	playerRepo := newMockPlayerRepoForUserPlayer()
	playerSvc := player.NewPlayerService(playerRepo, &fakeUserRepository{}, &fakeGameRepository{}, newFakeOrderRepositoryForEarnings(), &fakeReviewRepository{}, &fakePlayerTagRepository{}, &fakeCache{})

	router := gin.New()
	router.POST("/player/apply", func(c *gin.Context) {
		c.Set("user_id", uint64(500))
		applyAsPlayerHandler(c, playerSvc)
	})

	req := httptest.NewRequest(http.MethodPost, "/player/apply", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

func TestApplyAsPlayerHandler_AlreadyPlayer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	playerRepo := newMockPlayerRepoForUserPlayer()
	playerSvc := player.NewPlayerService(playerRepo, &fakeUserRepository{}, &fakeGameRepository{}, newFakeOrderRepositoryForEarnings(), &fakeReviewRepository{}, &fakePlayerTagRepository{}, &fakeCache{})

	router := gin.New()
	router.POST("/player/apply", func(c *gin.Context) {
		// User 100 already has a player record
		c.Set("user_id", uint64(100))
		applyAsPlayerHandler(c, playerSvc)
	})

	reqBody := player.ApplyPlayerRequest{
		Nickname:        "DuplicatePlayer",
		MainGameID:      10,
		Rank:            "Gold",
		HourlyRateCents: 5000,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/player/apply", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

func TestGetPlayerProfileHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	playerRepo := newMockPlayerRepoForUserPlayer()
	playerSvc := player.NewPlayerService(playerRepo, &fakeUserRepository{}, &fakeGameRepository{}, newFakeOrderRepositoryForEarnings(), &fakeReviewRepository{}, &fakePlayerTagRepository{}, &fakeCache{})

	router := gin.New()
	router.GET("/player/profile", func(c *gin.Context) {
		// User 100 has player ID 1
		c.Set("user_id", uint64(100))
		getPlayerProfileHandler(c, playerSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/player/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp model.APIResponse[player.PlayerDetailResponse]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("Expected success=true, got %v", resp.Success)
	}
}

func TestGetPlayerProfileHandler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	playerRepo := newMockPlayerRepoForUserPlayer()
	playerSvc := player.NewPlayerService(playerRepo, &fakeUserRepository{}, &fakeGameRepository{}, newFakeOrderRepositoryForEarnings(), &fakeReviewRepository{}, &fakePlayerTagRepository{}, &fakeCache{})

	router := gin.New()
	router.GET("/player/profile", func(c *gin.Context) {
		// User 9999 does not have a player record
		c.Set("user_id", uint64(9999))
		getPlayerProfileHandler(c, playerSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/player/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected status 404, got %d", w.Code)
	}
}

func TestUpdatePlayerProfileHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	playerRepo := newMockPlayerRepoForUserPlayer()
	playerSvc := player.NewPlayerService(playerRepo, &fakeUserRepository{}, &fakeGameRepository{}, newFakeOrderRepositoryForEarnings(), &fakeReviewRepository{}, &fakePlayerTagRepository{}, &fakeCache{})

	router := gin.New()
	router.PUT("/player/profile", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		updatePlayerProfileHandler(c, playerSvc)
	})

	reqBody := player.UpdatePlayerProfileRequest{
		Nickname:        "UpdatedNickname",
		Bio:             "Updated introduction",
		Rank:            "Platinum",
		HourlyRateCents: 7000,
		Tags:            []string{"pro", "friendly"},
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/player/profile", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdatePlayerProfileHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	playerRepo := newMockPlayerRepoForUserPlayer()
	playerSvc := player.NewPlayerService(playerRepo, &fakeUserRepository{}, &fakeGameRepository{}, newFakeOrderRepositoryForEarnings(), &fakeReviewRepository{}, &fakePlayerTagRepository{}, &fakeCache{})

	router := gin.New()
	router.PUT("/player/profile", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		updatePlayerProfileHandler(c, playerSvc)
	})

	req := httptest.NewRequest(http.MethodPut, "/player/profile", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

func TestSetPlayerStatusHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	playerRepo := newMockPlayerRepoForUserPlayer()
	playerSvc := player.NewPlayerService(playerRepo, &fakeUserRepository{}, &fakeGameRepository{}, newFakeOrderRepositoryForEarnings(), &fakeReviewRepository{}, &fakePlayerTagRepository{}, &fakeCache{})

	router := gin.New()
	router.PUT("/player/status", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		setPlayerStatusHandler(c, playerSvc)
	})

	reqBody := player.SetPlayerStatusRequest{
		Online: true,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/player/status", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSetPlayerStatusHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	playerRepo := newMockPlayerRepoForUserPlayer()
	playerSvc := player.NewPlayerService(playerRepo, &fakeUserRepository{}, &fakeGameRepository{}, newFakeOrderRepositoryForEarnings(), &fakeReviewRepository{}, &fakePlayerTagRepository{}, &fakeCache{})

	router := gin.New()
	router.PUT("/player/status", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		setPlayerStatusHandler(c, playerSvc)
	})

	req := httptest.NewRequest(http.MethodPut, "/player/status", bytes.NewBuffer([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}
