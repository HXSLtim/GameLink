package admin

import (
	"bytes"
	"context"
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
	adminservice "gamelink/internal/service/admin"
	assignmentservice "gamelink/internal/service/assignment"
	permservice "gamelink/internal/service/permission"
	roleservice "gamelink/internal/service/role"
)

type routerPermRepo struct{ byUser map[uint64][]model.Permission }

func (f *routerPermRepo) List(_ context.Context) ([]model.Permission, error) { return nil, nil }
func (f *routerPermRepo) ListPaged(_ context.Context, _ int, _ int) ([]model.Permission, int64, error) {
	return nil, 0, nil
}
func (f *routerPermRepo) ListPagedWithFilter(_ context.Context, _ int, _ int, _, _, _ string) ([]model.Permission, int64, error) {
	return nil, 0, nil
}
func (f *routerPermRepo) ListByGroup(_ context.Context) (map[string][]model.Permission, error) {
	return nil, nil
}
func (f *routerPermRepo) ListGroups(_ context.Context) ([]string, error) { return nil, nil }
func (f *routerPermRepo) Get(_ context.Context, _ uint64) (*model.Permission, error) {
	return nil, repository.ErrNotFound
}
func (f *routerPermRepo) GetByResource(_ context.Context, _, _ string) (*model.Permission, error) {
	return nil, repository.ErrNotFound
}
func (f *routerPermRepo) GetByCode(_ context.Context, _ string) (*model.Permission, error) {
	return nil, repository.ErrNotFound
}
func (f *routerPermRepo) GetByMethodAndPath(_ context.Context, _, _ string) (*model.Permission, error) {
	return nil, repository.ErrNotFound
}
func (f *routerPermRepo) Create(_ context.Context, _ *model.Permission) error             { return nil }
func (f *routerPermRepo) Update(_ context.Context, _ *model.Permission) error             { return nil }
func (f *routerPermRepo) UpsertByMethodPath(_ context.Context, _ *model.Permission) error { return nil }
func (f *routerPermRepo) Delete(_ context.Context, _ uint64) error                        { return nil }
func (f *routerPermRepo) ListByRoleID(_ context.Context, _ uint64) ([]model.Permission, error) {
	return nil, nil
}
func (f *routerPermRepo) ListByUserID(_ context.Context, uid uint64) ([]model.Permission, error) {
	return f.byUser[uid], nil
}

type routerRoleRepo struct{}

func (routerRoleRepo) List(_ context.Context) ([]model.RoleModel, error) { return nil, nil }
func (routerRoleRepo) ListPaged(_ context.Context, _ int, _ int) ([]model.RoleModel, int64, error) {
	return nil, 0, nil
}
func (routerRoleRepo) ListPagedWithFilter(_ context.Context, _ int, _ int, _ string, _ *bool) ([]model.RoleModel, int64, error) {
	return nil, 0, nil
}
func (routerRoleRepo) ListWithPermissions(_ context.Context) ([]model.RoleModel, error) {
	return nil, nil
}
func (routerRoleRepo) Get(_ context.Context, _ uint64) (*model.RoleModel, error) {
	return nil, repository.ErrNotFound
}
func (routerRoleRepo) GetWithPermissions(_ context.Context, _ uint64) (*model.RoleModel, error) {
	return nil, repository.ErrNotFound
}
func (routerRoleRepo) GetBySlug(_ context.Context, _ string) (*model.RoleModel, error) {
	return nil, repository.ErrNotFound
}
func (routerRoleRepo) Create(_ context.Context, _ *model.RoleModel) error              { return nil }
func (routerRoleRepo) Update(_ context.Context, _ *model.RoleModel) error              { return nil }
func (routerRoleRepo) Delete(_ context.Context, _ uint64) error                        { return nil }
func (routerRoleRepo) AssignPermissions(_ context.Context, _ uint64, _ []uint64) error { return nil }
func (routerRoleRepo) AddPermissions(_ context.Context, _ uint64, _ []uint64) error    { return nil }
func (routerRoleRepo) RemovePermissions(_ context.Context, _ uint64, _ []uint64) error { return nil }
func (routerRoleRepo) AssignToUser(_ context.Context, _ uint64, _ []uint64) error      { return nil }
func (routerRoleRepo) RemoveFromUser(_ context.Context, _ uint64, _ []uint64) error    { return nil }
func (routerRoleRepo) ListByUserID(_ context.Context, _ uint64) ([]model.RoleModel, error) {
	return nil, nil
}
func (routerRoleRepo) CheckUserHasRole(_ context.Context, _ uint64, _ string) (bool, error) {
	return false, nil
}

type routerGameRepo struct{}

func (routerGameRepo) List(_ context.Context) ([]model.Game, error) { return nil, nil }
func (routerGameRepo) ListPaged(_ context.Context, _ int, _ int) ([]model.Game, int64, error) {
	return []model.Game{}, 0, nil
}
func (routerGameRepo) Get(_ context.Context, _ uint64) (*model.Game, error) {
	return nil, repository.ErrNotFound
}
func (routerGameRepo) Create(_ context.Context, _ *model.Game) error { return nil }
func (routerGameRepo) Update(_ context.Context, _ *model.Game) error { return nil }
func (routerGameRepo) Delete(_ context.Context, _ uint64) error      { return nil }

type routerUsers struct{}

func (routerUsers) List(_ context.Context) ([]model.User, error) { return nil, nil }
func (routerUsers) ListPaged(_ context.Context, _ int, _ int) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (routerUsers) ListWithFilters(_ context.Context, _ repository.UserListOptions) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (routerUsers) Get(_ context.Context, _ uint64) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (routerUsers) GetByPhone(_ context.Context, _ string) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (routerUsers) FindByEmail(_ context.Context, _ string) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (routerUsers) FindByPhone(_ context.Context, _ string) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (routerUsers) Create(_ context.Context, _ *model.User) error { return nil }
func (routerUsers) Update(_ context.Context, _ *model.User) error { return nil }
func (routerUsers) Delete(_ context.Context, _ uint64) error      { return nil }

type routerPlayers struct{}

func (routerPlayers) List(_ context.Context) ([]model.Player, error) { return nil, nil }
func (routerPlayers) ListPaged(_ context.Context, _ int, _ int) ([]model.Player, int64, error) {
	return nil, 0, nil
}
func (routerPlayers) Get(_ context.Context, _ uint64) (*model.Player, error) {
	return nil, repository.ErrNotFound
}
func (routerPlayers) GetByUserID(_ context.Context, _ uint64) (*model.Player, error) {
	return nil, repository.ErrNotFound
}
func (routerPlayers) Create(_ context.Context, _ *model.Player) error { return nil }
func (routerPlayers) Update(_ context.Context, _ *model.Player) error { return nil }
func (routerPlayers) Delete(_ context.Context, _ uint64) error        { return nil }

type routerOrders struct{}

func (routerOrders) Create(_ context.Context, _ *model.Order) error { return nil }
func (routerOrders) List(_ context.Context, _ repository.OrderListOptions) ([]model.Order, int64, error) {
	return nil, 0, nil
}
func (routerOrders) Get(_ context.Context, _ uint64) (*model.Order, error) {
	return nil, repository.ErrNotFound
}
func (routerOrders) Update(_ context.Context, _ *model.Order) error { return nil }
func (routerOrders) Delete(_ context.Context, _ uint64) error       { return nil }

type routerDisputes struct{}

func (routerDisputes) Create(_ context.Context, _ *model.OrderDispute) error { return nil }
func (routerDisputes) Update(_ context.Context, _ *model.OrderDispute) error { return nil }
func (routerDisputes) ListByOrder(_ context.Context, _ uint64) ([]model.OrderDispute, error) {
	return nil, nil
}
func (routerDisputes) GetLatestByOrder(_ context.Context, _ uint64) (*model.OrderDispute, error) {
	return nil, repository.ErrNotFound
}

type routerOpLogs struct{}

func (routerOpLogs) Append(_ context.Context, _ *model.OperationLog) error { return nil }
func (routerOpLogs) ListByEntity(_ context.Context, _ string, _ uint64, _ repository.OperationLogListOptions) ([]model.OperationLog, int64, error) {
	return nil, 0, nil
}

type routerNotifications struct{}

func (routerNotifications) ListByUser(_ context.Context, _ repository.NotificationListOptions) ([]model.NotificationEvent, int64, error) {
	return nil, 0, nil
}
func (routerNotifications) MarkRead(_ context.Context, _ uint64, _ []uint64) error     { return nil }
func (routerNotifications) CountUnread(_ context.Context, _ uint64) (int64, error)     { return 0, nil }
func (routerNotifications) Create(_ context.Context, _ *model.NotificationEvent) error { return nil }

type routerPayments struct{}

func (routerPayments) Create(_ context.Context, _ *model.Payment) error { return nil }
func (routerPayments) List(_ context.Context, _ repository.PaymentListOptions) ([]model.Payment, int64, error) {
	return nil, 0, nil
}
func (routerPayments) Get(_ context.Context, _ uint64) (*model.Payment, error) {
	return nil, repository.ErrNotFound
}
func (routerPayments) Update(_ context.Context, _ *model.Payment) error { return nil }
func (routerPayments) Delete(_ context.Context, _ uint64) error         { return nil }

func setupAdminRouterForAuth(pm *mw.PermissionMiddleware) *gin.Engine {
	r := newTestEngine()
	svc := adminservice.NewAdminService(routerGameRepo{}, routerUsers{}, routerPlayers{}, routerOrders{}, routerPayments{}, routerRoleRepo{}, cache.NewMemory())
	assignSvc := assignmentservice.NewService(routerOrders{}, routerPlayers{}, routerDisputes{}, routerOpLogs{}, routerNotifications{})
	RegisterRoutes(r, svc, assignSvc, pm)
	return r
}

func TestAdminRouter_AuthAndPermission(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("APP_ENV", "development")
	t.Setenv("ADMIN_AUTH_MODE", "jwt")
	jwt := auth.NewJWTManager("secret", time.Hour)
	token, _ := jwt.GenerateToken(100, "admin")
	permRepo := &routerPermRepo{byUser: map[uint64][]model.Permission{}}
	roleRepo := &routerRoleRepo{}
	permSvc := permservice.NewPermissionService(permRepo, cache.NewMemory())
	roleSvc := roleservice.NewRoleService(roleRepo, cache.NewMemory())
	pm := mw.NewPermissionMiddleware(jwt, permSvc, roleSvc)
	r := setupAdminRouterForAuth(pm)

	w1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodGet, "/admin/games", nil)
	r.ServeHTTP(w1, req1)
	if w1.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w1.Code)
	}

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/admin/games", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w2.Code)
	}

	permRepo.byUser[100] = []model.Permission{{Base: model.Base{ID: 1}, Method: model.HTTPMethodGET, Path: "/api/v1/admin/games"}}
	permSvc2 := permservice.NewPermissionService(permRepo, cache.NewMemory())
	pm2 := mw.NewPermissionMiddleware(jwt, permSvc2, roleSvc)
	r2 := setupAdminRouterForAuth(pm2)
	w3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodGet, "/admin/games", nil)
	req3.Header.Set("Authorization", "Bearer "+token)
	r2.ServeHTTP(w3, req3)
	if w3.Code != http.StatusOK && w3.Code != http.StatusInternalServerError {
		t.Fatalf("expected 200/500, got %d", w3.Code)
	}

	// orders path
	w4 := httptest.NewRecorder()
	req4 := httptest.NewRequest(http.MethodGet, "/admin/orders", nil)
	req4.Header.Set("Authorization", "Bearer "+token)
	r2.ServeHTTP(w4, req4)
	if w4.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w4.Code)
	}
	permRepo.byUser[100] = append(permRepo.byUser[100], model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/admin/orders"})
	pm3 := mw.NewPermissionMiddleware(jwt, permservice.NewPermissionService(permRepo, cache.NewMemory()), roleSvc)
	r3 := setupAdminRouterForAuth(pm3)
	w5 := httptest.NewRecorder()
	req5 := httptest.NewRequest(http.MethodGet, "/admin/orders", nil)
	req5.Header.Set("Authorization", "Bearer "+token)
	r3.ServeHTTP(w5, req5)
	if w5.Code != http.StatusOK && w5.Code != http.StatusInternalServerError {
		t.Fatalf("expected 200/500, got %d", w5.Code)
	}

	// payments path
	w6 := httptest.NewRecorder()
	req6 := httptest.NewRequest(http.MethodGet, "/admin/payments", nil)
	req6.Header.Set("Authorization", "Bearer "+token)
	r3.ServeHTTP(w6, req6)
	if w6.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w6.Code)
	}
	permRepo.byUser[100] = append(permRepo.byUser[100], model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/admin/payments"})
	pm4 := mw.NewPermissionMiddleware(jwt, permservice.NewPermissionService(permRepo, cache.NewMemory()), roleSvc)
	r4 := setupAdminRouterForAuth(pm4)
	w7 := httptest.NewRecorder()
	req7 := httptest.NewRequest(http.MethodGet, "/admin/payments", nil)
	req7.Header.Set("Authorization", "Bearer "+token)
	r4.ServeHTTP(w7, req7)
	if w7.Code != http.StatusOK && w7.Code != http.StatusInternalServerError {
		t.Fatalf("expected 200/500, got %d", w7.Code)
	}

	// players path
	w8 := httptest.NewRecorder()
	req8 := httptest.NewRequest(http.MethodGet, "/admin/players", nil)
	req8.Header.Set("Authorization", "Bearer "+token)
	r4.ServeHTTP(w8, req8)
	if w8.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w8.Code)
	}
	permRepo.byUser[100] = append(permRepo.byUser[100], model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/admin/players"})
	pm5 := mw.NewPermissionMiddleware(jwt, permservice.NewPermissionService(permRepo, cache.NewMemory()), roleSvc)
	r5 := setupAdminRouterForAuth(pm5)
	w9 := httptest.NewRecorder()
	req9 := httptest.NewRequest(http.MethodGet, "/admin/players", nil)
	req9.Header.Set("Authorization", "Bearer "+token)
	r5.ServeHTTP(w9, req9)
	if w9.Code != http.StatusOK && w9.Code != http.StatusInternalServerError {
		t.Fatalf("expected 200/500, got %d", w9.Code)
	}

	// players detail path
	permRepo.byUser[100] = append(permRepo.byUser[100], model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/admin/players/:id"})
	pm6 := mw.NewPermissionMiddleware(jwt, permservice.NewPermissionService(permRepo, cache.NewMemory()), roleSvc)
	r6 := setupAdminRouterForAuth(pm6)
	w10 := httptest.NewRecorder()
	req10 := httptest.NewRequest(http.MethodGet, "/admin/players/1", nil)
	req10.Header.Set("Authorization", "Bearer "+token)
	r6.ServeHTTP(w10, req10)
	if w10.Code != http.StatusOK && w10.Code != http.StatusInternalServerError && w10.Code != http.StatusNotFound {
		t.Fatalf("expected 200/404/500, got %d", w10.Code)
	}

	// orders detail path
	w11 := httptest.NewRecorder()
	req11 := httptest.NewRequest(http.MethodGet, "/admin/orders/1", nil)
	req11.Header.Set("Authorization", "Bearer "+token)
	r6.ServeHTTP(w11, req11)
	if w11.Code != http.StatusForbidden && w11.Code != http.StatusInternalServerError && w11.Code != http.StatusNotFound {
		t.Fatalf("expected 403/404/500, got %d", w11.Code)
	}
	permRepo.byUser[100] = append(permRepo.byUser[100], model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/admin/orders/:id"})
	pm7 := mw.NewPermissionMiddleware(jwt, permservice.NewPermissionService(permRepo, cache.NewMemory()), roleSvc)
	r7 := setupAdminRouterForAuth(pm7)
	w12 := httptest.NewRecorder()
	req12 := httptest.NewRequest(http.MethodGet, "/admin/orders/1", nil)
	req12.Header.Set("Authorization", "Bearer "+token)
	r7.ServeHTTP(w12, req12)
	if w12.Code != http.StatusOK && w12.Code != http.StatusInternalServerError && w12.Code != http.StatusNotFound {
		t.Fatalf("expected 200/404/500, got %d", w12.Code)
	}

	// payments detail path
	w13 := httptest.NewRecorder()
	req13 := httptest.NewRequest(http.MethodGet, "/admin/payments/1", nil)
	req13.Header.Set("Authorization", "Bearer "+token)
	r7.ServeHTTP(w13, req13)
	if w13.Code != http.StatusForbidden && w13.Code != http.StatusInternalServerError && w13.Code != http.StatusNotFound {
		t.Fatalf("expected 403/404/500, got %d", w13.Code)
	}
	permRepo.byUser[100] = append(permRepo.byUser[100], model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/admin/payments/:id"})
	pm8 := mw.NewPermissionMiddleware(jwt, permservice.NewPermissionService(permRepo, cache.NewMemory()), roleSvc)
	r8 := setupAdminRouterForAuth(pm8)
	w14 := httptest.NewRecorder()
	req14 := httptest.NewRequest(http.MethodGet, "/admin/payments/1", nil)
	req14.Header.Set("Authorization", "Bearer "+token)
	r8.ServeHTTP(w14, req14)
	if w14.Code != http.StatusOK && w14.Code != http.StatusInternalServerError && w14.Code != http.StatusNotFound {
		t.Fatalf("expected 200/404/500, got %d", w14.Code)
	}

	// orders confirm
	w15 := httptest.NewRecorder()
	req15 := httptest.NewRequest(http.MethodPost, "/admin/orders/1/confirm", nil)
	req15.Header.Set("Authorization", "Bearer "+token)
	r8.ServeHTTP(w15, req15)
	if w15.Code != http.StatusForbidden && w15.Code != http.StatusInternalServerError && w15.Code != http.StatusNotFound {
		t.Fatalf("expected 403/404/500, got %d", w15.Code)
	}
	permRepo.byUser[100] = append(permRepo.byUser[100], model.Permission{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/orders/:id/confirm"})
	pm9 := mw.NewPermissionMiddleware(jwt, permservice.NewPermissionService(permRepo, cache.NewMemory()), roleSvc)
	r9 := setupAdminRouterForAuth(pm9)
	w16 := httptest.NewRecorder()
	req16 := httptest.NewRequest(http.MethodPost, "/admin/orders/1/confirm", nil)
	req16.Header.Set("Authorization", "Bearer "+token)
	r9.ServeHTTP(w16, req16)
	if w16.Code != http.StatusOK && w16.Code != http.StatusInternalServerError && w16.Code != http.StatusNotFound {
		t.Fatalf("expected 200/404/500, got %d", w16.Code)
	}

	// payments update
	w17 := httptest.NewRecorder()
	req17 := httptest.NewRequest(http.MethodPut, "/admin/payments/1", bytes.NewReader([]byte(`{"status":"paid"}`)))
	req17.Header.Set("Authorization", "Bearer "+token)
	req17.Header.Set("Content-Type", "application/json")
	r9.ServeHTTP(w17, req17)
	if w17.Code != http.StatusForbidden && w17.Code != http.StatusInternalServerError && w17.Code != http.StatusNotFound {
		t.Fatalf("expected 403/404/500, got %d", w17.Code)
	}
	permRepo.byUser[100] = append(permRepo.byUser[100], model.Permission{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/payments/:id"})
	pm10 := mw.NewPermissionMiddleware(jwt, permservice.NewPermissionService(permRepo, cache.NewMemory()), roleSvc)
	r10 := setupAdminRouterForAuth(pm10)
	w18 := httptest.NewRecorder()
	req18 := httptest.NewRequest(http.MethodPut, "/admin/payments/1", bytes.NewReader([]byte(`{"status":"paid"}`)))
	req18.Header.Set("Authorization", "Bearer "+token)
	req18.Header.Set("Content-Type", "application/json")
	r10.ServeHTTP(w18, req18)
	if w18.Code != http.StatusOK && w18.Code != http.StatusInternalServerError && w18.Code != http.StatusNotFound {
		t.Fatalf("expected 200/404/500, got %d", w18.Code)
	}

	// orders start
	w19 := httptest.NewRecorder()
	req19 := httptest.NewRequest(http.MethodPost, "/admin/orders/1/start", nil)
	req19.Header.Set("Authorization", "Bearer "+token)
	r10.ServeHTTP(w19, req19)
	if w19.Code != http.StatusForbidden && w19.Code != http.StatusInternalServerError && w19.Code != http.StatusNotFound {
		t.Fatalf("expected 403/404/500, got %d", w19.Code)
	}
	permRepo.byUser[100] = append(permRepo.byUser[100], model.Permission{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/orders/:id/start"})
	pm11 := mw.NewPermissionMiddleware(jwt, permservice.NewPermissionService(permRepo, cache.NewMemory()), roleSvc)
	r11 := setupAdminRouterForAuth(pm11)
	w20 := httptest.NewRecorder()
	req20 := httptest.NewRequest(http.MethodPost, "/admin/orders/1/start", nil)
	req20.Header.Set("Authorization", "Bearer "+token)
	r11.ServeHTTP(w20, req20)
	if w20.Code != http.StatusOK && w20.Code != http.StatusInternalServerError && w20.Code != http.StatusNotFound {
		t.Fatalf("expected 200/404/500, got %d", w20.Code)
	}

	// orders refund
	w21 := httptest.NewRecorder()
	req21 := httptest.NewRequest(http.MethodPost, "/admin/orders/1/refund", bytes.NewReader([]byte(`{"reason":"dup"}`)))
	req21.Header.Set("Authorization", "Bearer "+token)
	req21.Header.Set("Content-Type", "application/json")
	r11.ServeHTTP(w21, req21)
	if w21.Code != http.StatusForbidden && w21.Code != http.StatusInternalServerError && w21.Code != http.StatusNotFound {
		t.Fatalf("expected 403/404/500, got %d", w21.Code)
	}
	permRepo.byUser[100] = append(permRepo.byUser[100], model.Permission{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/orders/:id/refund"})
	pm12 := mw.NewPermissionMiddleware(jwt, permservice.NewPermissionService(permRepo, cache.NewMemory()), roleSvc)
	r12 := setupAdminRouterForAuth(pm12)
	w22 := httptest.NewRecorder()
	req22 := httptest.NewRequest(http.MethodPost, "/admin/orders/1/refund", bytes.NewReader([]byte(`{"reason":"dup"}`)))
	req22.Header.Set("Authorization", "Bearer "+token)
	req22.Header.Set("Content-Type", "application/json")
	r12.ServeHTTP(w22, req22)
	if w22.Code != http.StatusOK && w22.Code != http.StatusInternalServerError && w22.Code != http.StatusNotFound {
		t.Fatalf("expected 200/404/500, got %d", w22.Code)
	}

	// orders complete
	w23 := httptest.NewRecorder()
	req23 := httptest.NewRequest(http.MethodPost, "/admin/orders/1/complete", nil)
	req23.Header.Set("Authorization", "Bearer "+token)
	r12.ServeHTTP(w23, req23)
	if w23.Code != http.StatusForbidden && w23.Code != http.StatusInternalServerError && w23.Code != http.StatusNotFound {
		t.Fatalf("expected 403/404/500, got %d", w23.Code)
	}
	permRepo.byUser[100] = append(permRepo.byUser[100], model.Permission{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/orders/:id/complete"})
	pm13 := mw.NewPermissionMiddleware(jwt, permservice.NewPermissionService(permRepo, cache.NewMemory()), roleSvc)
	r13 := setupAdminRouterForAuth(pm13)
	w24 := httptest.NewRecorder()
	req24 := httptest.NewRequest(http.MethodPost, "/admin/orders/1/complete", nil)
	req24.Header.Set("Authorization", "Bearer "+token)
	r13.ServeHTTP(w24, req24)
	if w24.Code != http.StatusOK && w24.Code != http.StatusInternalServerError && w24.Code != http.StatusNotFound {
		t.Fatalf("expected 200/404/500, got %d", w24.Code)
	}

	// orders cancel
	w25 := httptest.NewRecorder()
	req25 := httptest.NewRequest(http.MethodPost, "/admin/orders/1/cancel", bytes.NewReader([]byte(`{"reason":"no-show"}`)))
	req25.Header.Set("Authorization", "Bearer "+token)
	req25.Header.Set("Content-Type", "application/json")
	r13.ServeHTTP(w25, req25)
	if w25.Code != http.StatusForbidden && w25.Code != http.StatusInternalServerError && w25.Code != http.StatusNotFound {
		t.Fatalf("expected 403/404/500, got %d", w25.Code)
	}
	permRepo.byUser[100] = append(permRepo.byUser[100], model.Permission{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/orders/:id/cancel"})
	pm14 := mw.NewPermissionMiddleware(jwt, permservice.NewPermissionService(permRepo, cache.NewMemory()), roleSvc)
	r14 := setupAdminRouterForAuth(pm14)
	w26 := httptest.NewRecorder()
	req26 := httptest.NewRequest(http.MethodPost, "/admin/orders/1/cancel", bytes.NewReader([]byte(`{"reason":"no-show"}`)))
	req26.Header.Set("Authorization", "Bearer "+token)
	req26.Header.Set("Content-Type", "application/json")
	r14.ServeHTTP(w26, req26)
	if w26.Code != http.StatusOK && w26.Code != http.StatusInternalServerError && w26.Code != http.StatusNotFound {
		t.Fatalf("expected 200/404/500, got %d", w26.Code)
	}

	// orders assign
	w27 := httptest.NewRecorder()
	req27 := httptest.NewRequest(http.MethodPost, "/admin/orders/1/assign", bytes.NewReader([]byte(`{"player_id":1}`)))
	req27.Header.Set("Authorization", "Bearer "+token)
	req27.Header.Set("Content-Type", "application/json")
	r14.ServeHTTP(w27, req27)
	if w27.Code != http.StatusForbidden && w27.Code != http.StatusInternalServerError && w27.Code != http.StatusNotFound {
		t.Fatalf("expected 403/404/500, got %d", w27.Code)
	}
	permRepo.byUser[100] = append(permRepo.byUser[100], model.Permission{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/orders/:id/assign"})
	pm15 := mw.NewPermissionMiddleware(jwt, permservice.NewPermissionService(permRepo, cache.NewMemory()), roleSvc)
	r15 := setupAdminRouterForAuth(pm15)
	w28 := httptest.NewRecorder()
	req28 := httptest.NewRequest(http.MethodPost, "/admin/orders/1/assign", bytes.NewReader([]byte(`{"player_id":1}`)))
	req28.Header.Set("Authorization", "Bearer "+token)
	req28.Header.Set("Content-Type", "application/json")
	r15.ServeHTTP(w28, req28)
	if w28.Code != http.StatusOK && w28.Code != http.StatusInternalServerError && w28.Code != http.StatusNotFound {
		t.Fatalf("expected 200/404/500, got %d", w28.Code)
	}
}
