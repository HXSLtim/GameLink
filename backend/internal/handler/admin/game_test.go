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

type fakeGameRepo struct {
	items     []model.Game
	listCalls int
	listPaged func(page, size int) ([]model.Game, int64, error)
}

func (f *fakeGameRepo) List(ctx context.Context) ([]model.Game, error) {
	f.listCalls++
	return append([]model.Game(nil), f.items...), nil
}

func (f *fakeGameRepo) ListPaged(ctx context.Context, page, size int) ([]model.Game, int64, error) {
	if f.listPaged != nil {
		return f.listPaged(page, size)
	}
	return append([]model.Game(nil), f.items...), int64(len(f.items)), nil
}

func (f *fakeGameRepo) Get(ctx context.Context, id uint64) (*model.Game, error) {
	for i := range f.items {
		if f.items[i].ID == id {
			c := f.items[i]
			return &c, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeGameRepo) Create(ctx context.Context, g *model.Game) error {
	if g.ID == 0 {
		g.ID = uint64(len(f.items) + 1)
	}
	f.items = append(f.items, *g)
	return nil
}

func (f *fakeGameRepo) Update(ctx context.Context, g *model.Game) error {
	for i := range f.items {
		if f.items[i].ID == g.ID {
			f.items[i] = *g
			return nil
		}
	}
	return repository.ErrNotFound
}

func (f *fakeGameRepo) Delete(ctx context.Context, id uint64) error {
	for i := range f.items {
		if f.items[i].ID == id {
			f.items = append(f.items[:i], f.items[i+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}

type fakeUserRepo struct{ last *model.User }

func (f *fakeUserRepo) List(ctx context.Context) ([]model.User, error) {
	if f.last != nil {
		return []model.User{*f.last}, nil
	}
	return []model.User{}, nil
}

func (f *fakeUserRepo) ListPaged(ctx context.Context, page, size int) ([]model.User, int64, error) {
	if f.last != nil {
		return []model.User{*f.last}, 1, nil
	}
	return []model.User{}, 0, nil
}

func (f *fakeUserRepo) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	if f.last != nil {
		return []model.User{*f.last}, 1, nil
	}
	return []model.User{}, 0, nil
}

func (f *fakeUserRepo) Get(ctx context.Context, id uint64) (*model.User, error) {
	if f.last != nil && f.last.ID == id {
		return f.last, nil
	}
	return nil, repository.ErrNotFound
}

func (f *fakeUserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) { return f.last, nil }
func (f *fakeUserRepo) FindByPhone(ctx context.Context, phone string) (*model.User, error) { return f.last, nil }
func (f *fakeUserRepo) GetByPhone(ctx context.Context, phone string) (*model.User, error) { return f.last, nil }
func (f *fakeUserRepo) Create(ctx context.Context, u *model.User) error {
	if u.ID == 0 {
		u.ID = 1
	}
	f.last = u
	return nil
}
func (f *fakeUserRepo) Update(ctx context.Context, u *model.User) error { f.last = u; return nil }
func (f *fakeUserRepo) Delete(ctx context.Context, id uint64) error     { return nil }

type fakePlayerRepo struct{}

func (f *fakePlayerRepo) List(ctx context.Context) ([]model.Player, error) { return []model.Player{}, nil }
func (f *fakePlayerRepo) ListPaged(ctx context.Context, page, size int) ([]model.Player, int64, error) {
	return nil, 0, nil
}
func (f *fakePlayerRepo) Get(ctx context.Context, id uint64) (*model.Player, error) {
	return nil, repository.ErrNotFound
}
func (f *fakePlayerRepo) Create(ctx context.Context, p *model.Player) error { return nil }
func (f *fakePlayerRepo) Update(ctx context.Context, p *model.Player) error { return nil }
func (f *fakePlayerRepo) Delete(ctx context.Context, id uint64) error       { return nil }
func (f *fakePlayerRepo) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) { return nil, repository.ErrNotFound }

type fakeOrderRepo struct{ obj *model.Order }

func (f *fakeOrderRepo) List(ctx context.Context, _ repository.OrderListOptions) ([]model.Order, int64, error) {
	return nil, 0, nil
}
func (f *fakeOrderRepo) Create(ctx context.Context, o *model.Order) error {
	if o.ID == 0 {
		o.ID = 1
	}
	f.obj = o
	return nil
}
func (f *fakeOrderRepo) Get(ctx context.Context, id uint64) (*model.Order, error) {
	if f.obj == nil {
		return nil, repository.ErrNotFound
	}
	return f.obj, nil
}
func (f *fakeOrderRepo) Update(ctx context.Context, o *model.Order) error { f.obj = o; return nil }
func (f *fakeOrderRepo) Delete(ctx context.Context, id uint64) error      { return nil }

type fakePaymentRepo struct{ obj *model.Payment }

func (f *fakePaymentRepo) List(ctx context.Context, _ repository.PaymentListOptions) ([]model.Payment, int64, error) {
	return nil, 0, nil
}
func (f *fakePaymentRepo) Create(ctx context.Context, p *model.Payment) error {
	if p.ID == 0 {
		p.ID = 1
	}
	f.obj = p
	return nil
}
func (f *fakePaymentRepo) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	if f.obj == nil {
		return nil, repository.ErrNotFound
	}
	return f.obj, nil
}
func (f *fakePaymentRepo) Update(ctx context.Context, p *model.Payment) error { f.obj = p; return nil }
func (f *fakePaymentRepo) Delete(ctx context.Context, id uint64) error        { return nil }

type fakeRoleRepo struct{}

func (f *fakeRoleRepo) List(ctx context.Context) ([]model.RoleModel, error) { return nil, nil }
func (f *fakeRoleRepo) ListPaged(ctx context.Context, page, pageSize int) ([]model.RoleModel, int64, error) {
	return nil, 0, nil
}
func (f *fakeRoleRepo) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, isSystem *bool) ([]model.RoleModel, int64, error) {
	return nil, 0, nil
}
func (f *fakeRoleRepo) ListWithPermissions(ctx context.Context) ([]model.RoleModel, error) { return nil, nil }
func (f *fakeRoleRepo) Get(ctx context.Context, id uint64) (*model.RoleModel, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeRoleRepo) GetWithPermissions(ctx context.Context, id uint64) (*model.RoleModel, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeRoleRepo) GetBySlug(ctx context.Context, slug string) (*model.RoleModel, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeRoleRepo) Create(ctx context.Context, role *model.RoleModel) error { return nil }
func (f *fakeRoleRepo) Update(ctx context.Context, role *model.RoleModel) error { return nil }
func (f *fakeRoleRepo) Delete(ctx context.Context, id uint64) error             { return nil }
func (f *fakeRoleRepo) AssignPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	return nil
}
func (f *fakeRoleRepo) AddPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	return nil
}
func (f *fakeRoleRepo) RemovePermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	return nil
}
func (f *fakeRoleRepo) ListByUserID(ctx context.Context, userID uint64) ([]model.RoleModel, error) {
	return nil, nil
}
func (f *fakeRoleRepo) AssignToUser(ctx context.Context, userID uint64, roleIDs []uint64) error {
	return nil
}
func (f *fakeRoleRepo) RemoveFromUser(ctx context.Context, userID uint64, roleIDs []uint64) error {
	return nil
}
func (f *fakeRoleRepo) CheckUserHasRole(ctx context.Context, userID uint64, roleSlug string) (bool, error) {
	return false, nil
}

func setupGameTestRouter(gameRepo *fakeGameRepo) (*gin.Engine, *adminservice.AdminService) {
	r := newTestEngine()

	svc := adminservice.NewAdminService(
		gameRepo,
		&fakeUserRepo{},
		&fakePlayerRepo{},
		&fakeOrderRepo{},
		&fakePaymentRepo{},
		&fakeRoleRepo{},
		cache.NewMemory(),
	)

	handler := NewGameHandler(svc)
	r.GET("/admin/games", handler.ListGames)
	r.GET("/admin/games/:id", handler.GetGame)
	r.POST("/admin/games", handler.CreateGame)
	r.PUT("/admin/games/:id", handler.UpdateGame)
	r.DELETE("/admin/games/:id", handler.DeleteGame)
	r.GET("/admin/games/:id/logs", handler.ListGameLogs)

	return r, svc
}

func TestGameHandler_ListGames(t *testing.T) {
	gameRepo := &fakeGameRepo{
		items: []model.Game{
			{Base: model.Base{ID: 1}, Key: "lol", Name: "英雄联盟"},
			{Base: model.Base{ID: 2}, Key: "dota2", Name: "DOTA 2"},
		},
	}
	r, _ := setupGameTestRouter(gameRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/games?page=1&page_size=20", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[[]model.Game]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Len(t, resp.Data, 2)
}

func TestGameHandler_ListGames_Pagination(t *testing.T) {
	gameRepo := &fakeGameRepo{
		listPaged: func(page, size int) ([]model.Game, int64, error) {
			return []model.Game{
				{Base: model.Base{ID: 1}, Key: "lol", Name: "英雄联盟"},
			}, 1, nil
		},
	}
	r, _ := setupGameTestRouter(gameRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/games?page=2&page_size=10", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[[]model.Game]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Pagination)
}

func TestGameHandler_GetGame(t *testing.T) {
	gameRepo := &fakeGameRepo{
		items: []model.Game{
			{Base: model.Base{ID: 1}, Key: "lol", Name: "英雄联盟"},
		},
	}
	r, _ := setupGameTestRouter(gameRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/games/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[*model.Game]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, uint64(1), resp.Data.ID)
}

func TestGameHandler_GetGame_NotFound(t *testing.T) {
	gameRepo := &fakeGameRepo{
		items: []model.Game{},
	}
	r, _ := setupGameTestRouter(gameRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/games/999", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGameHandler_GetGame_InvalidID(t *testing.T) {
	gameRepo := &fakeGameRepo{}
	r, _ := setupGameTestRouter(gameRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/games/invalid", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGameHandler_CreateGame(t *testing.T) {
	gameRepo := &fakeGameRepo{
		items: []model.Game{},
	}
	r, _ := setupGameTestRouter(gameRepo)

	payload := GamePayload{
		Key:         "test-game",
		Name:        "测试游戏",
		Category:    "MOBA",
		IconURL:     "https://example.com/icon.png",
		Description: "这是一个测试游戏",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/games", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp model.APIResponse[*model.Game]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, "test-game", resp.Data.Key)
}

func TestGameHandler_CreateGame_Validation(t *testing.T) {
	gameRepo := &fakeGameRepo{}
	r, _ := setupGameTestRouter(gameRepo)

	payload := GamePayload{
		Key:  "",
		Name: "测试游戏",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/games", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/admin/games", bytes.NewReader([]byte("{invalid json")))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestGameHandler_UpdateGame(t *testing.T) {
	gameRepo := &fakeGameRepo{
		items: []model.Game{
			{Base: model.Base{ID: 1}, Key: "lol", Name: "英雄联盟"},
		},
	}
	r, _ := setupGameTestRouter(gameRepo)

	payload := GamePayload{
		Key:         "lol",
		Name:        "英雄联盟（更新）",
		Category:    "MOBA",
		IconURL:     "https://example.com/icon.png",
		Description: "更新后的描述",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/admin/games/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[*model.Game]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Data)
}

func TestGameHandler_UpdateGame_NotFound(t *testing.T) {
	gameRepo := &fakeGameRepo{
		items: []model.Game{},
	}
	r, _ := setupGameTestRouter(gameRepo)

	payload := GamePayload{
		Key:  "lol",
		Name: "英雄联盟",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/admin/games/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestGameHandler_DeleteGame(t *testing.T) {
	gameRepo := &fakeGameRepo{
		items: []model.Game{
			{Base: model.Base{ID: 1}, Key: "lol", Name: "英雄联盟"},
		},
	}
	r, _ := setupGameTestRouter(gameRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/admin/games/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestGameHandler_DeleteGame_NotFound(t *testing.T) {
	gameRepo := &fakeGameRepo{
		items: []model.Game{},
	}
	r, _ := setupGameTestRouter(gameRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/admin/games/999", nil)
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestGameHandler_ListGameLogs(t *testing.T) {
	t.Skip("ListGameLogs requires TxManager, skipping for now")
}
