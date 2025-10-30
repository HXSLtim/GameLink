package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/auth"
	"gamelink/internal/cache"
	mw "gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service"
)

// testRouterConfig 用于配置测试路由器的选项
type testRouterConfig struct {
	permRepo *fakePermissionRepo
	roleRepo *fakeRoleRepo
}

// buildTestRouter constructs a Gin engine with admin routes and a provided service.
func buildTestRouter(svc *service.AdminService) *gin.Engine {
	return buildTestRouterWithConfig(svc, nil)
}

// buildTestRouterWithConfig constructs a Gin engine with custom repo configuration.
func buildTestRouterWithConfig(svc *service.AdminService, cfg *testRouterConfig) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1")

	// Create mock permission middleware
	jwtMgr := auth.NewJWTManager("test-secret", 24*time.Hour)

	var permRepo *fakePermissionRepo
	var roleRepo *fakeRoleRepo

	if cfg != nil && cfg.permRepo != nil {
		permRepo = cfg.permRepo
	} else {
		permRepo = newFakePermissionRepo()
	}

	if cfg != nil && cfg.roleRepo != nil {
		roleRepo = cfg.roleRepo
	} else {
		roleRepo = newFakeRoleRepo()
	}

	// 使用内存缓存避免 nil pointer
	cache := cache.NewMemory()
	permService := service.NewPermissionService(permRepo, cache)
	roleService := service.NewRoleService(roleRepo, cache)
	permMiddleware := mw.NewPermissionMiddleware(jwtMgr, permService, roleService)

	RegisterRoutes(api, svc, permMiddleware)
	return r
}

// generateTestJWT 生成测试用的 JWT 令牌
func generateTestJWT(userID uint64, role string) string {
	jwtMgr := auth.NewJWTManager("test-secret", 24*time.Hour)
	token, _ := jwtMgr.GenerateToken(userID, role)
	return token
}

func TestAdminRoutes_UnauthorizedWhenTokenConfigured(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("ADMIN_TOKEN", "secret")

	svc := service.NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, &fakePaymentRepo{}, &fakeRoleRepo{}, nil)
	r := buildTestRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/games", nil)
	req.Header.Set("Authorization", "Bearer "+generateTestJWT(1, "admin"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAdminRoutes_ListGames_Envelope(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("ADMIN_AUTH_MODE", "jwt")

	svc := service.NewAdminService(&fakeGameRepo{items: []model.Game{{Name: "A"}}}, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, &fakePaymentRepo{}, &fakeRoleRepo{}, nil)
	r := buildTestRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/games?page=1&page_size=10", nil)
	req.Header.Set("Authorization", "Bearer "+generateTestJWT(1, "admin"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Success bool         `json:"success"`
		Code    int          `json:"code"`
		Data    []model.Game `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if !body.Success || body.Code != 200 || len(body.Data) == 0 {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

func TestAdminRoutes_UpdateOrder_AcceptsCancelledSpelling(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("ADMIN_AUTH_MODE", "jwt")

	o := &model.Order{Base: model.Base{ID: 1}, Status: model.OrderStatusPending, Currency: model.CurrencyCNY}
	oRepo := &fakeOrderRepo{obj: o}
	svc := service.NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, oRepo, &fakePaymentRepo{}, &fakeRoleRepo{}, nil)
	r := buildTestRouter(svc)

	payload := map[string]any{
		"status":      "cancelled", //nolint:misspell // intentionally testing legacy spelling
		"price_cents": 100,
		"currency":    "USD",
	}
	buf, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/orders/1", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+generateTestJWT(1, "admin"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body=%s", w.Code, w.Body.String())
	}
	if oRepo.obj.Status != model.OrderStatusCanceled {
		t.Fatalf("expected status normalized to 'canceled', got %q", oRepo.obj.Status)
	}
}

// ---- minimal fakes for integration tests ----

type fakeGameRepo struct{ items []model.Game }

func (f *fakeGameRepo) List(ctx context.Context) ([]model.Game, error) {
	return append([]model.Game(nil), f.items...), nil
}
func (f *fakeGameRepo) ListPaged(ctx context.Context, page, size int) ([]model.Game, int64, error) {
	return append([]model.Game(nil), f.items...), int64(len(f.items)), nil
}
func (f *fakeGameRepo) Get(ctx context.Context, id uint64) (*model.Game, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeGameRepo) Create(ctx context.Context, g *model.Game) error { return nil }
func (f *fakeGameRepo) Update(ctx context.Context, g *model.Game) error { return nil }
func (f *fakeGameRepo) Delete(ctx context.Context, id uint64) error     { return nil }

type fakeUserRepo struct{}

func (f *fakeUserRepo) List(ctx context.Context) ([]model.User, error) { return nil, nil }
func (f *fakeUserRepo) ListPaged(ctx context.Context, page, size int) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (f *fakeUserRepo) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (f *fakeUserRepo) Get(ctx context.Context, id uint64) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeUserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeUserRepo) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeUserRepo) Create(ctx context.Context, u *model.User) error { return nil }
func (f *fakeUserRepo) Update(ctx context.Context, u *model.User) error { return nil }
func (f *fakeUserRepo) Delete(ctx context.Context, id uint64) error     { return nil }

type fakePlayerRepo struct{}

func (f *fakePlayerRepo) List(ctx context.Context) ([]model.Player, error) { return nil, nil }
func (f *fakePlayerRepo) ListPaged(ctx context.Context, page, size int) ([]model.Player, int64, error) {
	return nil, 0, nil
}
func (f *fakePlayerRepo) Get(ctx context.Context, id uint64) (*model.Player, error) {
	return nil, repository.ErrNotFound
}
func (f *fakePlayerRepo) Create(ctx context.Context, p *model.Player) error { return nil }
func (f *fakePlayerRepo) Update(ctx context.Context, p *model.Player) error { return nil }
func (f *fakePlayerRepo) Delete(ctx context.Context, id uint64) error       { return nil }

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

type fakePaymentRepo struct{}

func (f *fakePaymentRepo) List(ctx context.Context, _ repository.PaymentListOptions) ([]model.Payment, int64, error) {
	return nil, 0, nil
}
func (f *fakePaymentRepo) Create(ctx context.Context, p *model.Payment) error { return nil }
func (f *fakePaymentRepo) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	return nil, repository.ErrNotFound
}
func (f *fakePaymentRepo) Update(ctx context.Context, p *model.Payment) error { return nil }
func (f *fakePaymentRepo) Delete(ctx context.Context, id uint64) error        { return nil }

func TestPaymentHandler_InvalidTime_Returns400(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("ADMIN_AUTH_MODE", "jwt")

	pay := &model.Payment{Base: model.Base{ID: 1}, Status: model.PaymentStatusPending}
	pRepo := &fakePaymentRepoWithObj{obj: pay}
	svc := service.NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, pRepo, &fakeRoleRepo{}, nil)
	r := buildTestRouter(svc)

	payload := map[string]any{
		"status":  "paid",
		"paid_at": "not-a-time",
	}
	buf, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/payments/1", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+generateTestJWT(1, "admin"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid paid_at, got %d", w.Code)
	}
}

func TestAdminUserValidation_InvalidEmailAndPhone(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("ADMIN_AUTH_MODE", "jwt")

	svc := service.NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, &fakePaymentRepo{}, &fakeRoleRepo{}, nil)
	r := buildTestRouter(svc)

	// invalid email
	badEmail := map[string]any{
		"email":    "not-an-email",
		"password": "Abcdef1!",
		"name":     "user",
		"role":     "user",
		"status":   "active",
	}
	buf, _ := json.Marshal(badEmail)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+generateTestJWT(1, "admin"))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid email, got %d", w.Code)
	}

	// invalid phone
	badPhone := map[string]any{
		"phone":    "12345",
		"password": "Abcdef1!",
		"name":     "user",
		"role":     "user",
		"status":   "active",
	}
	buf2, _ := json.Marshal(badPhone)
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users", bytes.NewReader(buf2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+generateTestJWT(1, "admin"))
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid phone, got %d", w2.Code)
	}
}

func TestAdmin_CreateUserWithPlayer_InvalidEmail(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("ADMIN_AUTH_MODE", "jwt")

	svc := service.NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, &fakePaymentRepo{}, &fakeRoleRepo{}, nil)
	r := buildTestRouter(svc)

	body := map[string]any{
		"email":    "bad",
		"password": "Abcdef1!",
		"name":     "n",
		"role":     "user",
		"status":   "active",
		"player":   map[string]any{"verification_status": "pending"},
	}
	buf, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users/with-player", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+generateTestJWT(1, "admin"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid email, got %d", w.Code)
	}
}

type fakePaymentRepoWithObj struct{ obj *model.Payment }

func (f *fakePaymentRepoWithObj) List(ctx context.Context, _ repository.PaymentListOptions) ([]model.Payment, int64, error) {
	return nil, 0, nil
}
func (f *fakePaymentRepoWithObj) Create(ctx context.Context, p *model.Payment) error {
	f.obj = p
	return nil
}
func (f *fakePaymentRepoWithObj) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	if f.obj == nil {
		return nil, repository.ErrNotFound
	}
	return f.obj, nil
}
func (f *fakePaymentRepoWithObj) Update(ctx context.Context, p *model.Payment) error {
	f.obj = p
	return nil
}
func (f *fakePaymentRepoWithObj) Delete(ctx context.Context, id uint64) error { return nil }

// fakeRoleRepo implements repository.RoleRepository for testing
// 支持自定义用户角色映射，便于测试不同权限场景
type fakeRoleRepo struct {
	userRoles map[uint64][]model.RoleModel // userID -> roles
}

func newFakeRoleRepo() *fakeRoleRepo {
	return &fakeRoleRepo{
		userRoles: make(map[uint64][]model.RoleModel),
	}
}

// setUserRoles 设置用户拥有的角色（用于测试）
func (f *fakeRoleRepo) setUserRoles(userID uint64, roles []model.RoleModel) {
	f.userRoles[userID] = roles
}

func (f *fakeRoleRepo) List(ctx context.Context) ([]model.RoleModel, error) { return nil, nil }
func (f *fakeRoleRepo) ListPaged(ctx context.Context, page, pageSize int) ([]model.RoleModel, int64, error) {
	return nil, 0, nil
}
func (f *fakeRoleRepo) ListWithPermissions(ctx context.Context) ([]model.RoleModel, error) {
	return nil, nil
}
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
	// 查找用户的角色配置
	if roles, ok := f.userRoles[userID]; ok {
		return roles, nil
	}
	// 默认返回 super_admin（向后兼容）
	return []model.RoleModel{
		{
			Base:        model.Base{ID: 1},
			Slug:        string(model.RoleSlugSuperAdmin),
			Name:        "超级管理员",
			Description: "测试用超级管理员",
			IsSystem:    true,
		},
	}, nil
}
func (f *fakeRoleRepo) AssignToUser(ctx context.Context, userID uint64, roleIDs []uint64) error {
	return nil
}
func (f *fakeRoleRepo) RemoveFromUser(ctx context.Context, userID uint64, roleIDs []uint64) error {
	return nil
}
func (f *fakeRoleRepo) CheckUserHasRole(ctx context.Context, userID uint64, roleSlug string) (bool, error) {
	roles, err := f.ListByUserID(ctx, userID)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		if role.Slug == roleSlug {
			return true, nil
		}
	}
	return false, nil
}

// fakePermissionRepo implements repository.PermissionRepository for testing
// 支持自定义用户权限映射，便于测试不同权限场景
type fakePermissionRepo struct {
	userPermissions map[uint64][]model.Permission // userID -> permissions
}

func newFakePermissionRepo() *fakePermissionRepo {
	return &fakePermissionRepo{
		userPermissions: make(map[uint64][]model.Permission),
	}
}

// setUserPermissions 设置用户拥有的权限（用于测试）
func (f *fakePermissionRepo) setUserPermissions(userID uint64, permissions []model.Permission) {
	f.userPermissions[userID] = permissions
}

func (f *fakePermissionRepo) List(ctx context.Context) ([]model.Permission, error) { return nil, nil }
func (f *fakePermissionRepo) ListPaged(ctx context.Context, page, pageSize int) ([]model.Permission, int64, error) {
	return nil, 0, nil
}
func (f *fakePermissionRepo) ListByGroup(ctx context.Context, group string) ([]model.Permission, error) {
	return nil, nil
}
func (f *fakePermissionRepo) ListGroups(ctx context.Context) ([]string, error) { return nil, nil }
func (f *fakePermissionRepo) Get(ctx context.Context, id uint64) (*model.Permission, error) {
	return nil, repository.ErrNotFound
}
func (f *fakePermissionRepo) GetByMethodAndPath(ctx context.Context, method model.HTTPMethod, path string) (*model.Permission, error) {
	return nil, repository.ErrNotFound
}
func (f *fakePermissionRepo) GetByCode(ctx context.Context, code string) (*model.Permission, error) {
	return nil, repository.ErrNotFound
}
func (f *fakePermissionRepo) Create(ctx context.Context, permission *model.Permission) error {
	return nil
}
func (f *fakePermissionRepo) CreateBatch(ctx context.Context, permissions []model.Permission) error {
	return nil
}
func (f *fakePermissionRepo) Update(ctx context.Context, permission *model.Permission) error {
	return nil
}
func (f *fakePermissionRepo) Delete(ctx context.Context, id uint64) error { return nil }
func (f *fakePermissionRepo) UpsertByMethodPath(ctx context.Context, permission *model.Permission) error {
	return nil
}
func (f *fakePermissionRepo) ListByRoleID(ctx context.Context, roleID uint64) ([]model.Permission, error) {
	return nil, nil
}
func (f *fakePermissionRepo) ListByUserID(ctx context.Context, userID uint64) ([]model.Permission, error) {
	// 查找用户的权限配置
	if perms, ok := f.userPermissions[userID]; ok {
		return perms, nil
	}
	// 默认返回空（由 super_admin 快速通道处理）
	return nil, nil
}

// ========== RBAC 自定义角色权限测试 ==========

// TestCustomRole_WithSpecificPermission 测试自定义角色只拥有特定权限
func TestCustomRole_WithSpecificPermission(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("ADMIN_AUTH_MODE", "jwt")

	// 创建自定义角色：game_viewer（只能查看游戏列表）
	roleRepo := newFakeRoleRepo()
	permRepo := newFakePermissionRepo()

	// 设置用户角色：userID=2 拥有 game_viewer 角色
	roleRepo.setUserRoles(2, []model.RoleModel{
		{
			Base:        model.Base{ID: 10},
			Slug:        "game_viewer",
			Name:        "游戏查看员",
			Description: "只能查看游戏列表",
			IsSystem:    false,
		},
	})

	// 设置用户权限：userID=2 只拥有 GET /api/v1/admin/games 权限
	permRepo.setUserPermissions(2, []model.Permission{
		{
			Base:   model.Base{ID: 1},
			Method: model.HTTPMethodGET,
			Path:   "/api/v1/admin/games",
			Group:  "/admin/games",
			Code:   "admin.games.list",
		},
	})

	svc := service.NewAdminService(&fakeGameRepo{items: []model.Game{{Name: "Test Game"}}}, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, &fakePaymentRepo{}, &fakeRoleRepo{}, nil)
	r := buildTestRouterWithConfig(svc, &testRouterConfig{
		permRepo: permRepo,
		roleRepo: roleRepo,
	})

	// 测试：可以访问有权限的接口
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/games?page=1&page_size=10", nil)
	req.Header.Set("Authorization", "Bearer "+generateTestJWT(2, "game_viewer"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for allowed endpoint, got %d; body=%s", w.Code, w.Body.String())
	}

	// 验证返回的数据
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if success, ok := body["success"].(bool); !ok || !success {
		t.Fatalf("expected success=true, got: %v", body)
	}
}

// TestCustomRole_WithoutPermission 测试自定义角色访问无权限接口被拒绝
func TestCustomRole_WithoutPermission(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("ADMIN_AUTH_MODE", "jwt")

	// 创建自定义角色：game_viewer（只能查看游戏列表）
	roleRepo := newFakeRoleRepo()
	permRepo := newFakePermissionRepo()

	// 设置用户角色：userID=2 拥有 game_viewer 角色
	roleRepo.setUserRoles(2, []model.RoleModel{
		{
			Base:        model.Base{ID: 10},
			Slug:        "game_viewer",
			Name:        "游戏查看员",
			Description: "只能查看游戏列表",
			IsSystem:    false,
		},
	})

	// 设置用户权限：userID=2 只拥有 GET /api/v1/admin/games 权限（没有 POST 权限）
	permRepo.setUserPermissions(2, []model.Permission{
		{
			Base:   model.Base{ID: 1},
			Method: model.HTTPMethodGET,
			Path:   "/api/v1/admin/games",
			Group:  "/admin/games",
			Code:   "admin.games.list",
		},
	})

	svc := service.NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, &fakePaymentRepo{}, &fakeRoleRepo{}, nil)
	r := buildTestRouterWithConfig(svc, &testRouterConfig{
		permRepo: permRepo,
		roleRepo: roleRepo,
	})

	// 测试：无权限访问 POST /games（创建游戏）应被拒绝
	body := map[string]any{
		"key":  "test-game",
		"name": "Test Game",
	}
	buf, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/games", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+generateTestJWT(2, "game_viewer"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for forbidden endpoint, got %d; body=%s", w.Code, w.Body.String())
	}

	// 验证错误消息包含权限不足提示
	if !bytes.Contains(w.Body.Bytes(), []byte("权限不足")) {
		t.Fatalf("expected '权限不足' in error message, got: %s", w.Body.String())
	}
}

// TestSuperAdmin_HasAllPermissions 测试超级管理员可以访问所有接口
func TestSuperAdmin_HasAllPermissions(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("ADMIN_AUTH_MODE", "jwt")

	roleRepo := newFakeRoleRepo()
	permRepo := newFakePermissionRepo()

	// 设置用户角色：userID=1 拥有 super_admin 角色
	roleRepo.setUserRoles(1, []model.RoleModel{
		{
			Base:        model.Base{ID: 1},
			Slug:        string(model.RoleSlugSuperAdmin),
			Name:        "超级管理员",
			Description: "拥有所有权限",
			IsSystem:    true,
		},
	})

	// super_admin 不需要设置具体权限，会通过快速通道

	svc := service.NewAdminService(&fakeGameRepo{items: []model.Game{{Name: "Test"}}}, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, &fakePaymentRepo{}, &fakeRoleRepo{}, nil)
	r := buildTestRouterWithConfig(svc, &testRouterConfig{
		permRepo: permRepo,
		roleRepo: roleRepo,
	})

	// 测试：super_admin 可以访问任意接口（这里测试 GET /games 和 POST /games）
	tests := []struct {
		method string
		path   string
		body   map[string]any
	}{
		{http.MethodGet, "/api/v1/admin/games?page=1&page_size=10", nil},
		{http.MethodPost, "/api/v1/admin/games", map[string]any{"key": "test", "name": "Test"}},
	}

	for _, tt := range tests {
		var req *http.Request
		if tt.body != nil {
			buf, _ := json.Marshal(tt.body)
			req = httptest.NewRequest(tt.method, tt.path, bytes.NewReader(buf))
			req.Header.Set("Content-Type", "application/json")
		} else {
			req = httptest.NewRequest(tt.method, tt.path, nil)
		}
		req.Header.Set("Authorization", "Bearer "+generateTestJWT(1, "admin"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// super_admin 应该能访问所有接口（200 或其他非 403）
		if w.Code == http.StatusForbidden {
			t.Fatalf("super_admin should not get 403 for %s %s, got: %s", tt.method, tt.path, w.Body.String())
		}
	}
}

// TestCustomRole_MultiplePermissions 测试自定义角色拥有多个权限
func TestCustomRole_MultiplePermissions(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("ADMIN_AUTH_MODE", "jwt")

	// 创建自定义角色：game_manager（可以查看和创建游戏）
	roleRepo := newFakeRoleRepo()
	permRepo := newFakePermissionRepo()

	roleRepo.setUserRoles(3, []model.RoleModel{
		{
			Base:        model.Base{ID: 11},
			Slug:        "game_manager",
			Name:        "游戏管理员",
			Description: "可以查看和创建游戏",
			IsSystem:    false,
		},
	})

	// 设置多个权限：GET 和 POST /games
	permRepo.setUserPermissions(3, []model.Permission{
		{
			Base:   model.Base{ID: 1},
			Method: model.HTTPMethodGET,
			Path:   "/api/v1/admin/games",
			Group:  "/admin/games",
			Code:   "admin.games.list",
		},
		{
			Base:   model.Base{ID: 2},
			Method: model.HTTPMethodPOST,
			Path:   "/api/v1/admin/games",
			Group:  "/admin/games",
			Code:   "admin.games.create",
		},
	})

	svc := service.NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, &fakePaymentRepo{}, &fakeRoleRepo{}, nil)
	r := buildTestRouterWithConfig(svc, &testRouterConfig{
		permRepo: permRepo,
		roleRepo: roleRepo,
	})

	// 测试 GET：应该成功
	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/admin/games?page=1&page_size=10", nil)
	req1.Header.Set("Authorization", "Bearer "+generateTestJWT(3, "game_manager"))
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Fatalf("expected 200 for GET /games, got %d; body=%s", w1.Code, w1.Body.String())
	}

	// 测试 POST：应该成功
	body := map[string]any{"key": "test", "name": "Test"}
	buf, _ := json.Marshal(body)
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/admin/games", bytes.NewReader(buf))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+generateTestJWT(3, "game_manager"))
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK && w2.Code != http.StatusCreated {
		t.Fatalf("expected 200/201 for POST /games, got %d; body=%s", w2.Code, w2.Body.String())
	}

	// 测试 DELETE：应该被拒绝（没有删除权限）
	req3 := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/games/1", nil)
	req3.Header.Set("Authorization", "Bearer "+generateTestJWT(3, "game_manager"))
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	if w3.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for DELETE /games/1, got %d; body=%s", w3.Code, w3.Body.String())
	}
}
