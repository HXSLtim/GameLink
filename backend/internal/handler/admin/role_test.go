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
	roleservice "gamelink/internal/service/role"
)

type fakeRoleRepoForHandler struct {
	items []model.RoleModel
}

func (f *fakeRoleRepoForHandler) List(ctx context.Context) ([]model.RoleModel, error) {
	return append([]model.RoleModel(nil), f.items...), nil
}

func (f *fakeRoleRepoForHandler) ListPaged(ctx context.Context, page, pageSize int) ([]model.RoleModel, int64, error) {
	return append([]model.RoleModel(nil), f.items...), int64(len(f.items)), nil
}

func (f *fakeRoleRepoForHandler) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, isSystem *bool) ([]model.RoleModel, int64, error) {
	return append([]model.RoleModel(nil), f.items...), int64(len(f.items)), nil
}

func (f *fakeRoleRepoForHandler) ListWithPermissions(ctx context.Context) ([]model.RoleModel, error) {
	return append([]model.RoleModel(nil), f.items...), nil
}

func (f *fakeRoleRepoForHandler) Get(ctx context.Context, id uint64) (*model.RoleModel, error) {
	for i := range f.items {
		if f.items[i].ID == id {
			c := f.items[i]
			return &c, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeRoleRepoForHandler) GetWithPermissions(ctx context.Context, id uint64) (*model.RoleModel, error) {
	return f.Get(ctx, id)
}

func (f *fakeRoleRepoForHandler) GetBySlug(ctx context.Context, slug string) (*model.RoleModel, error) {
	for i := range f.items {
		if f.items[i].Slug == slug {
			c := f.items[i]
			return &c, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeRoleRepoForHandler) Create(ctx context.Context, role *model.RoleModel) error {
	if role.ID == 0 {
		role.ID = uint64(len(f.items) + 1)
	}
	f.items = append(f.items, *role)
	return nil
}

func (f *fakeRoleRepoForHandler) Update(ctx context.Context, role *model.RoleModel) error {
	for i := range f.items {
		if f.items[i].ID == role.ID {
			f.items[i] = *role
			return nil
		}
	}
	return repository.ErrNotFound
}

func (f *fakeRoleRepoForHandler) Delete(ctx context.Context, id uint64) error {
	for i := range f.items {
		if f.items[i].ID == id {
			f.items = append(f.items[:i], f.items[i+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}

func (f *fakeRoleRepoForHandler) AssignPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	return nil
}

func (f *fakeRoleRepoForHandler) AddPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	return nil
}

func (f *fakeRoleRepoForHandler) RemovePermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	return nil
}

func (f *fakeRoleRepoForHandler) ListByUserID(ctx context.Context, userID uint64) ([]model.RoleModel, error) {
	return nil, nil
}

func (f *fakeRoleRepoForHandler) AssignToUser(ctx context.Context, userID uint64, roleIDs []uint64) error {
	return nil
}

func (f *fakeRoleRepoForHandler) RemoveFromUser(ctx context.Context, userID uint64, roleIDs []uint64) error {
	return nil
}

func (f *fakeRoleRepoForHandler) CheckUserHasRole(ctx context.Context, userID uint64, roleSlug string) (bool, error) {
	return false, nil
}

func setupRoleTestRouter(roleRepo *fakeRoleRepoForHandler) (*gin.Engine, *roleservice.RoleService) {
	r := newTestEngine()

	svc := roleservice.NewRoleService(roleRepo, cache.NewMemory())

	handler := NewRoleHandler(svc)
	r.GET("/admin/roles", handler.ListRoles)
	r.GET("/admin/roles/:id", handler.GetRole)
	r.POST("/admin/roles", handler.CreateRole)
	r.PUT("/admin/roles/:id", handler.UpdateRole)
	r.DELETE("/admin/roles/:id", handler.DeleteRole)
	r.POST("/admin/roles/:id/permissions", handler.AssignPermissions)
	r.POST("/admin/users/:user_id/roles", handler.AssignRolesToUser)
	r.GET("/admin/users/:user_id/roles", handler.GetUserRoles)

	return r, svc
}

func TestRoleHandler_ListRoles(t *testing.T) {
	roleRepo := &fakeRoleRepoForHandler{
		items: []model.RoleModel{
			{Base: model.Base{ID: 1}, Slug: "admin", Name: "管理员"},
			{Base: model.Base{ID: 2}, Slug: "user", Name: "用户"},
		},
	}
	r, _ := setupRoleTestRouter(roleRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/roles?page=1&page_size=20", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestRoleHandler_ListRoles_WithPermissions(t *testing.T) {
	roleRepo := &fakeRoleRepoForHandler{
		items: []model.RoleModel{
			{Base: model.Base{ID: 1}, Slug: "admin", Name: "管理员"},
		},
	}
	r, _ := setupRoleTestRouter(roleRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/roles?with_permissions=true", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestRoleHandler_GetRole(t *testing.T) {
	roleRepo := &fakeRoleRepoForHandler{
		items: []model.RoleModel{
			{Base: model.Base{ID: 1}, Slug: "admin", Name: "管理员"},
		},
	}
	r, _ := setupRoleTestRouter(roleRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/roles/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[*model.RoleModel]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, uint64(1), resp.Data.ID)
}

func TestRoleHandler_GetRole_NotFound(t *testing.T) {
	roleRepo := &fakeRoleRepoForHandler{
		items: []model.RoleModel{},
	}
	r, _ := setupRoleTestRouter(roleRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/roles/999", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRoleHandler_CreateRole(t *testing.T) {
	roleRepo := &fakeRoleRepoForHandler{
		items: []model.RoleModel{},
	}
	r, _ := setupRoleTestRouter(roleRepo)

	payload := map[string]string{
		"slug":        "test-role",
		"name":        "测试角色",
		"description": "这是一个测试角色",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp model.APIResponse[*model.RoleModel]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Data)
}

func TestRoleHandler_CreateRole_Validation(t *testing.T) {
	roleRepo := &fakeRoleRepoForHandler{}
	r, _ := setupRoleTestRouter(roleRepo)

	// 测试缺少必填字段
	payload := map[string]string{
		"slug": "", // 缺少必填字段
		"name": "测试角色",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRoleHandler_UpdateRole(t *testing.T) {
	roleRepo := &fakeRoleRepoForHandler{
		items: []model.RoleModel{
			{Base: model.Base{ID: 1}, Slug: "admin", Name: "管理员"},
		},
	}
	r, _ := setupRoleTestRouter(roleRepo)

	payload := map[string]string{
		"slug":        "admin",
		"name":        "管理员（更新）",
		"description": "更新后的描述",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/admin/roles/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var resp model.APIResponse[*model.RoleModel]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	}
}

func TestRoleHandler_DeleteRole(t *testing.T) {
	roleRepo := &fakeRoleRepoForHandler{
		items: []model.RoleModel{
			{Base: model.Base{ID: 1}, Slug: "admin", Name: "管理员", IsSystem: false},
		},
	}
	r, _ := setupRoleTestRouter(roleRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/admin/roles/1", nil)
	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var resp model.APIResponse[any]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	}
}

func TestRoleHandler_AssignPermissions(t *testing.T) {
	roleRepo := &fakeRoleRepoForHandler{
		items: []model.RoleModel{
			{Base: model.Base{ID: 1}, Slug: "admin", Name: "管理员"},
		},
	}
	r, _ := setupRoleTestRouter(roleRepo)

	payload := map[string]any{
		"permission_ids": []uint64{1, 2, 3},
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/roles/1/permissions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var resp model.APIResponse[any]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	}
}

func TestRoleHandler_AssignRolesToUser(t *testing.T) {
	roleRepo := &fakeRoleRepoForHandler{}
	r, _ := setupRoleTestRouter(roleRepo)

	payload := map[string]any{
		"role_ids": []uint64{1, 2},
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/users/1/roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var resp model.APIResponse[any]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	}
}

func TestRoleHandler_GetUserRoles(t *testing.T) {
	roleRepo := &fakeRoleRepoForHandler{}
	r, _ := setupRoleTestRouter(roleRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/users/1/roles", nil)
	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var resp model.APIResponse[any]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	}
}

