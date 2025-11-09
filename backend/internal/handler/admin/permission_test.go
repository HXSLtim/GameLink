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
	permissionservice "gamelink/internal/service/permission"
)

type fakePermissionRepoForHandler struct {
	items []model.Permission
}

func (f *fakePermissionRepoForHandler) List(ctx context.Context) ([]model.Permission, error) {
	return append([]model.Permission(nil), f.items...), nil
}

func (f *fakePermissionRepoForHandler) ListPaged(ctx context.Context, page, pageSize int) ([]model.Permission, int64, error) {
	return append([]model.Permission(nil), f.items...), int64(len(f.items)), nil
}

func (f *fakePermissionRepoForHandler) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword, method, group string) ([]model.Permission, int64, error) {
	return append([]model.Permission(nil), f.items...), int64(len(f.items)), nil
}

func (f *fakePermissionRepoForHandler) Get(ctx context.Context, id uint64) (*model.Permission, error) {
	for i := range f.items {
		if f.items[i].ID == id {
			c := f.items[i]
			return &c, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakePermissionRepoForHandler) GetByCode(ctx context.Context, code string) (*model.Permission, error) {
	for i := range f.items {
		if f.items[i].Code == code {
			p := f.items[i]
			return &p, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakePermissionRepoForHandler) GetBySlug(ctx context.Context, slug string) (*model.Permission, error) {
	return nil, repository.ErrNotFound
}

func (f *fakePermissionRepoForHandler) GetByResource(ctx context.Context, resource, action string) (*model.Permission, error) {
	return nil, repository.ErrNotFound
}

func (f *fakePermissionRepoForHandler) GetByMethodAndPath(ctx context.Context, method, path string) (*model.Permission, error) {
	for i := range f.items {
		if string(f.items[i].Method) == method && f.items[i].Path == path {
			c := f.items[i]
			return &c, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakePermissionRepoForHandler) Create(ctx context.Context, p *model.Permission) error {
	if p.ID == 0 {
		p.ID = uint64(len(f.items) + 1)
	}
	f.items = append(f.items, *p)
	return nil
}

func (f *fakePermissionRepoForHandler) Update(ctx context.Context, p *model.Permission) error {
	for i := range f.items {
		if f.items[i].ID == p.ID {
			f.items[i] = *p
			return nil
		}
	}
	return repository.ErrNotFound
}

func (f *fakePermissionRepoForHandler) Delete(ctx context.Context, id uint64) error {
	for i := range f.items {
		if f.items[i].ID == id {
			f.items = append(f.items[:i], f.items[i+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}

func (f *fakePermissionRepoForHandler) UpsertByMethodPath(ctx context.Context, p *model.Permission) error {
	existing, _ := f.GetByMethodAndPath(ctx, string(p.Method), p.Path)
	if existing != nil {
		return f.Update(ctx, p)
	}
	return f.Create(ctx, p)
}

func (f *fakePermissionRepoForHandler) ListByGroup(ctx context.Context) (map[string][]model.Permission, error) {
	groups := make(map[string][]model.Permission)
	for _, p := range f.items {
		groups[p.Group] = append(groups[p.Group], p)
	}
	return groups, nil
}

func (f *fakePermissionRepoForHandler) ListGroups(ctx context.Context) ([]string, error) {
	groups := make(map[string]bool)
	for _, p := range f.items {
		if p.Group != "" {
			groups[p.Group] = true
		}
	}
	result := make([]string, 0, len(groups))
	for g := range groups {
		result = append(result, g)
	}
	return result, nil
}

func (f *fakePermissionRepoForHandler) ListByRoleID(ctx context.Context, roleID uint64) ([]model.Permission, error) {
	return nil, nil
}

func (f *fakePermissionRepoForHandler) ListByUserID(ctx context.Context, userID uint64) ([]model.Permission, error) {
	return nil, nil
}

func setupPermissionTestRouter(permissionRepo *fakePermissionRepoForHandler) (*gin.Engine, *permissionservice.PermissionService) {
	r := newTestEngine()

	svc := permissionservice.NewPermissionService(permissionRepo, cache.NewMemory())

	handler := NewPermissionHandler(svc)
	r.GET("/admin/permissions", handler.ListPermissions)
	r.GET("/admin/permissions/:id", handler.GetPermission)
	r.POST("/admin/permissions", handler.CreatePermission)
	r.PUT("/admin/permissions/:id", handler.UpdatePermission)
	r.DELETE("/admin/permissions/:id", handler.DeletePermission)
	r.GET("/admin/roles/:role_id/permissions", handler.GetRolePermissions)
	r.GET("/admin/users/:user_id/permissions", handler.GetUserPermissions)
	r.GET("/admin/permissions/groups", handler.GetPermissionGroups)

	return r, svc
}

func TestPermissionHandler_ListPermissions(t *testing.T) {
	permissionRepo := &fakePermissionRepoForHandler{
		items: []model.Permission{
			{Base: model.Base{ID: 1}, Method: model.HTTPMethodGET, Path: "/api/users"},
			{Base: model.Base{ID: 2}, Method: model.HTTPMethodPOST, Path: "/api/users"},
		},
	}
	r, _ := setupPermissionTestRouter(permissionRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/permissions?page=1&page_size=20", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestPermissionHandler_GetPermission(t *testing.T) {
	permissionRepo := &fakePermissionRepoForHandler{
		items: []model.Permission{
			{Base: model.Base{ID: 1}, Method: model.HTTPMethodGET, Path: "/api/users"},
		},
	}
	r, _ := setupPermissionTestRouter(permissionRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/permissions/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[*model.Permission]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, uint64(1), resp.Data.ID)
}

func TestPermissionHandler_GetPermission_NotFound(t *testing.T) {
	permissionRepo := &fakePermissionRepoForHandler{
		items: []model.Permission{},
	}
	r, _ := setupPermissionTestRouter(permissionRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/permissions/999", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPermissionHandler_CreatePermission(t *testing.T) {
	permissionRepo := &fakePermissionRepoForHandler{
		items: []model.Permission{},
	}
	r, _ := setupPermissionTestRouter(permissionRepo)

	payload := map[string]any{
		"method":      "GET",
		"path":        "/api/test",
		"code":        "test:read",
		"group":       "test",
		"description": "测试权限",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/permissions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp model.APIResponse[*model.Permission]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Data)
}

func TestPermissionHandler_CreatePermission_Validation(t *testing.T) {
	permissionRepo := &fakePermissionRepoForHandler{}
	r, _ := setupPermissionTestRouter(permissionRepo)

	// 测试缺少必填字段
	payload := map[string]any{
		"method": "", // 缺少必填字段
		"path":   "/api/test",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/permissions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPermissionHandler_UpdatePermission(t *testing.T) {
	permissionRepo := &fakePermissionRepoForHandler{
		items: []model.Permission{
			{Base: model.Base{ID: 1}, Method: model.HTTPMethodGET, Path: "/api/users"},
		},
	}
	r, _ := setupPermissionTestRouter(permissionRepo)

	payload := map[string]any{
		"code":        "users:read:updated",
		"group":       "users",
		"description": "更新后的描述",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/admin/permissions/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var resp model.APIResponse[*model.Permission]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	}
}

func TestPermissionHandler_DeletePermission(t *testing.T) {
	permissionRepo := &fakePermissionRepoForHandler{
		items: []model.Permission{
			{Base: model.Base{ID: 1}, Method: model.HTTPMethodGET, Path: "/api/users"},
		},
	}
	r, _ := setupPermissionTestRouter(permissionRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/admin/permissions/1", nil)
	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var resp model.APIResponse[any]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	}
}

func TestPermissionHandler_GetRolePermissions(t *testing.T) {
	permissionRepo := &fakePermissionRepoForHandler{}
	r, _ := setupPermissionTestRouter(permissionRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/roles/1/permissions", nil)
	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var resp model.APIResponse[any]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	}
}

func TestPermissionHandler_GetUserPermissions(t *testing.T) {
	permissionRepo := &fakePermissionRepoForHandler{}
	r, _ := setupPermissionTestRouter(permissionRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/users/1/permissions", nil)
	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var resp model.APIResponse[any]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	}
}

func TestPermissionHandler_GetPermissionGroups(t *testing.T) {
	permissionRepo := &fakePermissionRepoForHandler{
		items: []model.Permission{
			{Base: model.Base{ID: 1}, Method: model.HTTPMethodGET, Path: "/api/users", Group: "users"},
			{Base: model.Base{ID: 2}, Method: model.HTTPMethodPOST, Path: "/api/orders", Group: "orders"},
		},
	}
	r, _ := setupPermissionTestRouter(permissionRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/permissions/groups", nil)
	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var resp model.APIResponse[any]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	}
}

