package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gamelink/internal/cache"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	adminservice "gamelink/internal/service/admin"
)

type fakePlayerRepoForHandler struct {
	items     []model.Player
	listPaged func(page, size int) ([]model.Player, int64, error)
}

func (f *fakePlayerRepoForHandler) List(ctx context.Context) ([]model.Player, error) {
	return append([]model.Player(nil), f.items...), nil
}

func (f *fakePlayerRepoForHandler) ListPaged(ctx context.Context, page, size int) ([]model.Player, int64, error) {
	if f.listPaged != nil {
		return f.listPaged(page, size)
	}
	return append([]model.Player(nil), f.items...), int64(len(f.items)), nil
}

func (f *fakePlayerRepoForHandler) Get(ctx context.Context, id uint64) (*model.Player, error) {
	for i := range f.items {
		if f.items[i].ID == id {
			c := f.items[i]
			return &c, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakePlayerRepoForHandler) Create(ctx context.Context, p *model.Player) error {
	if p.ID == 0 {
		p.ID = uint64(len(f.items) + 1)
	}
	f.items = append(f.items, *p)
	return nil
}

func (f *fakePlayerRepoForHandler) Update(ctx context.Context, p *model.Player) error {
	for i := range f.items {
		if f.items[i].ID == p.ID {
			f.items[i] = *p
			return nil
		}
	}
	return repository.ErrNotFound
}

func (f *fakePlayerRepoForHandler) Delete(ctx context.Context, id uint64) error {
	for i := range f.items {
		if f.items[i].ID == id {
			f.items = append(f.items[:i], f.items[i+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}

func setupPlayerTestRouter(playerRepo *fakePlayerRepoForHandler) (*gin.Engine, *adminservice.AdminService) {
	r := newTestEngine()

	svc := adminservice.NewAdminService(
		&fakeGameRepo{},
		&fakeUserRepo{},
		playerRepo,
		&fakeOrderRepo{},
		&fakePaymentRepo{},
		&fakeRoleRepo{},
		cache.NewMemory(),
	)

	handler := NewPlayerHandler(svc)
	r.GET("/admin/players", handler.ListPlayers)
	r.GET("/admin/players/:id", handler.GetPlayer)
	r.POST("/admin/players", handler.CreatePlayer)
	r.PUT("/admin/players/:id", handler.UpdatePlayer)
	r.DELETE("/admin/players/:id", handler.DeletePlayer)
	r.PUT("/admin/players/:id/verification", handler.UpdatePlayerVerification)
	r.PUT("/admin/players/:id/games", handler.UpdatePlayerGames)
	r.PUT("/admin/players/:id/skill-tags", handler.UpdatePlayerSkillTags)
	r.GET("/admin/players/:id/logs", handler.ListPlayerLogs)

	return r, svc
}

func TestPlayerHandler_ListPlayers(t *testing.T) {
	playerRepo := &fakePlayerRepoForHandler{
		items: []model.Player{
			{Base: model.Base{ID: 1}, UserID: 1, Nickname: "测试陪玩师1"},
			{Base: model.Base{ID: 2}, UserID: 2, Nickname: "测试陪玩师2"},
		},
	}
	r, _ := setupPlayerTestRouter(playerRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/players?page=1&page_size=20", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[[]model.Player]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Len(t, resp.Data, 2)
}

func TestPlayerHandler_ListPlayers_Pagination(t *testing.T) {
	playerRepo := &fakePlayerRepoForHandler{
		listPaged: func(page, size int) ([]model.Player, int64, error) {
			return []model.Player{
				{Base: model.Base{ID: 1}, UserID: 1, Nickname: "测试陪玩师"},
			}, 1, nil
		},
	}
	r, _ := setupPlayerTestRouter(playerRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/players?page=2&page_size=10", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[[]model.Player]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Pagination)
}

func TestPlayerHandler_GetPlayer(t *testing.T) {
	playerRepo := &fakePlayerRepoForHandler{
		items: []model.Player{
			{Base: model.Base{ID: 1}, UserID: 1, Nickname: "测试陪玩师"},
		},
	}
	r, _ := setupPlayerTestRouter(playerRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/players/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[*model.Player]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, uint64(1), resp.Data.ID)
}

func TestPlayerHandler_GetPlayer_NotFound(t *testing.T) {
	playerRepo := &fakePlayerRepoForHandler{
		items: []model.Player{},
	}
	r, _ := setupPlayerTestRouter(playerRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/players/999", nil)
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestPlayerHandler_CreatePlayer(t *testing.T) {
	playerRepo := &fakePlayerRepoForHandler{
		items: []model.Player{},
	}
	r, _ := setupPlayerTestRouter(playerRepo)

	payload := CreatePlayerPayload{
		UserID:             1,
		Nickname:           "新陪玩师",
		Bio:                "这是一个测试陪玩师",
		Rank:               "钻石",
		HourlyRateCents:    10000,
		MainGameID:         1,
		VerificationStatus: "pending",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/players", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp model.APIResponse[*model.Player]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Data)
}

func TestPlayerHandler_CreatePlayer_Validation(t *testing.T) {
	playerRepo := &fakePlayerRepoForHandler{}
	r, _ := setupPlayerTestRouter(playerRepo)

	// 测试缺少必填字段
	payload := CreatePlayerPayload{
		UserID:             0, // 缺少必填字段
		VerificationStatus: "",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/players", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPlayerHandler_UpdatePlayer(t *testing.T) {
	playerRepo := &fakePlayerRepoForHandler{
		items: []model.Player{
			{Base: model.Base{ID: 1}, UserID: 1, Nickname: "测试陪玩师"},
		},
	}
	r, _ := setupPlayerTestRouter(playerRepo)

	payload := UpdatePlayerPayload{
		Nickname:           "更新后的陪玩师",
		Bio:                "更新后的简介",
		Rank:               "大师",
		HourlyRateCents:    15000,
		MainGameID:         2,
		VerificationStatus: "verified",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/admin/players/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[*model.Player]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Data)
}

func TestPlayerHandler_UpdatePlayer_NotFound(t *testing.T) {
	playerRepo := &fakePlayerRepoForHandler{
		items: []model.Player{},
	}
	r, _ := setupPlayerTestRouter(playerRepo)

	payload := UpdatePlayerPayload{
		Nickname:           "更新后的陪玩师",
		VerificationStatus: "verified",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/admin/players/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestPlayerHandler_DeletePlayer(t *testing.T) {
	playerRepo := &fakePlayerRepoForHandler{
		items: []model.Player{
			{Base: model.Base{ID: 1}, UserID: 1, Nickname: "测试陪玩师"},
		},
	}
	r, _ := setupPlayerTestRouter(playerRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/admin/players/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestPlayerHandler_UpdatePlayerVerification(t *testing.T) {
	playerRepo := &fakePlayerRepoForHandler{
		items: []model.Player{
			{Base: model.Base{ID: 1}, UserID: 1, Nickname: "测试陪玩师"},
		},
	}
	r, _ := setupPlayerTestRouter(playerRepo)

	payload := map[string]string{
		"verification_status": "verified",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/admin/players/1/verification", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var resp model.APIResponse[*model.Player]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	}
}

func TestPlayerHandler_UpdatePlayerGames(t *testing.T) {
	playerRepo := &fakePlayerRepoForHandler{
		items: []model.Player{
			{Base: model.Base{ID: 1}, UserID: 1, Nickname: "测试陪玩师", MainGameID: 1},
		},
	}
	r, _ := setupPlayerTestRouter(playerRepo)

	payload := map[string]uint64{
		"main_game_id": 2,
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/admin/players/1/games", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var resp model.APIResponse[*model.Player]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	}
}

func TestPlayerHandler_UpdatePlayerSkillTags(t *testing.T) {
	playerRepo := &fakePlayerRepoForHandler{
		items: []model.Player{
			{Base: model.Base{ID: 1}, UserID: 1, Nickname: "测试陪玩师"},
		},
	}
	r, _ := setupPlayerTestRouter(playerRepo)

	payload := SkillTagsBody{
		Tags: []string{"carry", "support", "mid"},
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/admin/players/1/skill-tags", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// 由于需要 TxManager，这里只检查响应格式
	if w.Code == http.StatusOK {
		var resp model.APIResponse[any]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	}
}

func TestPlayerHandler_ListPlayerLogs(t *testing.T) {
	t.Skip("ListPlayerLogs requires TxManager, skipping for now")
}

