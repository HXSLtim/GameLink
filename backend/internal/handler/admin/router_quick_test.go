package admin

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
    "os"

    "github.com/gin-gonic/gin"

    "gamelink/internal/auth"
    "gamelink/internal/cache"
    "gamelink/internal/model"
    mw "gamelink/internal/handler/middleware"
    "gamelink/internal/repository"
    permservice "gamelink/internal/service/permission"
    roleservice "gamelink/internal/service/role"
    stats "gamelink/internal/service/stats"
)

type permRepoStub struct{}
func (permRepoStub) List(context.Context) ([]model.Permission, error) { return nil, nil }
func (permRepoStub) ListPaged(context.Context, int, int) ([]model.Permission, int64, error) { return nil, 0, nil }
func (permRepoStub) ListPagedWithFilter(context.Context, int, int, string, string, string) ([]model.Permission, int64, error) { return nil, 0, nil }
func (permRepoStub) ListByGroup(context.Context) (map[string][]model.Permission, error) { return nil, nil }
func (permRepoStub) ListGroups(context.Context) ([]string, error) { return nil, nil }
func (permRepoStub) Get(context.Context, uint64) (*model.Permission, error) { return nil, repository.ErrNotFound }
func (permRepoStub) GetByResource(context.Context, string, string) (*model.Permission, error) { return nil, repository.ErrNotFound }
func (permRepoStub) GetByCode(context.Context, string) (*model.Permission, error) { return nil, repository.ErrNotFound }
func (permRepoStub) GetByMethodAndPath(context.Context, string, string) (*model.Permission, error) { return nil, repository.ErrNotFound }
func (permRepoStub) Create(context.Context, *model.Permission) error { return nil }
func (permRepoStub) Update(context.Context, *model.Permission) error { return nil }
func (permRepoStub) UpsertByMethodPath(context.Context, *model.Permission) error { return nil }
func (permRepoStub) Delete(context.Context, uint64) error { return nil }
func (permRepoStub) ListByRoleID(context.Context, uint64) ([]model.Permission, error) { return nil, nil }
func (permRepoStub) ListByUserID(context.Context, uint64) ([]model.Permission, error) {
    return []model.Permission{
        {Method: model.HTTPMethodGET, Path: "/api/v1/admin/stats/dashboard"},
        {Method: model.HTTPMethodGET, Path: "/api/v1/admin/stats/revenue-trend"},
        {Method: model.HTTPMethodGET, Path: "/api/v1/admin/stats/user-growth"},
        {Method: model.HTTPMethodGET, Path: "/api/v1/admin/stats/orders"},
        {Method: model.HTTPMethodGET, Path: "/api/v1/admin/stats/top-players"},
        {Method: model.HTTPMethodGET, Path: "/api/v1/admin/stats/audit/overview"},
        {Method: model.HTTPMethodGET, Path: "/api/v1/admin/stats/audit/trend"},
    }, nil
}

type roleRepoStub struct{}
func (roleRepoStub) List(context.Context) ([]model.RoleModel, error) { return nil, nil }
func (roleRepoStub) ListPaged(context.Context, int, int) ([]model.RoleModel, int64, error) { return nil, 0, nil }
func (roleRepoStub) ListPagedWithFilter(context.Context, int, int, string, *bool) ([]model.RoleModel, int64, error) { return nil, 0, nil }
func (roleRepoStub) ListWithPermissions(context.Context) ([]model.RoleModel, error) { return nil, nil }
func (roleRepoStub) Get(context.Context, uint64) (*model.RoleModel, error) { return nil, repository.ErrNotFound }
func (roleRepoStub) GetWithPermissions(context.Context, uint64) (*model.RoleModel, error) { return nil, repository.ErrNotFound }
func (roleRepoStub) GetBySlug(context.Context, string) (*model.RoleModel, error) { return nil, repository.ErrNotFound }
func (roleRepoStub) Create(context.Context, *model.RoleModel) error { return nil }
func (roleRepoStub) Update(context.Context, *model.RoleModel) error { return nil }
func (roleRepoStub) Delete(context.Context, uint64) error { return nil }
func (roleRepoStub) AssignPermissions(context.Context, uint64, []uint64) error { return nil }
func (roleRepoStub) AddPermissions(context.Context, uint64, []uint64) error { return nil }
func (roleRepoStub) RemovePermissions(context.Context, uint64, []uint64) error { return nil }
func (roleRepoStub) AssignToUser(context.Context, uint64, []uint64) error { return nil }
func (roleRepoStub) RemoveFromUser(context.Context, uint64, []uint64) error { return nil }
func (roleRepoStub) ListByUserID(context.Context, uint64) ([]model.RoleModel, error) { return nil, nil }
func (roleRepoStub) CheckUserHasRole(context.Context, uint64, string) (bool, error) { return true, nil }

type statsRepoStub struct{}
func (statsRepoStub) Dashboard(ctx context.Context) (repository.Dashboard, error) { return repository.Dashboard{}, nil }
func (statsRepoStub) RevenueTrend(ctx context.Context, days int) ([]repository.DateValue, error) { return []repository.DateValue{{Date:"2025-01-01", Value:int64(days)}}, nil }
func (statsRepoStub) UserGrowth(ctx context.Context, days int) ([]repository.DateValue, error) { return []repository.DateValue{{Date:"2025-01-01", Value:int64(days)}}, nil }
func (statsRepoStub) OrdersByStatus(ctx context.Context) (map[string]int64, error) { return map[string]int64{"pending":1}, nil }
func (statsRepoStub) TopPlayers(ctx context.Context, limit int) ([]repository.PlayerTop, error) { return []repository.PlayerTop{{PlayerID:1, Nickname:"p"}}, nil }
func (statsRepoStub) AuditOverview(ctx context.Context, from *time.Time, to *time.Time) (map[string]int64, map[string]int64, error) { return map[string]int64{"order":1}, map[string]int64{"create":1}, nil }
func (statsRepoStub) AuditTrend(ctx context.Context, from *time.Time, to *time.Time, entity string, action string) ([]repository.DateValue, error) { return []repository.DateValue{{Date:"2025-01-01", Value:1}}, nil }

func TestAdminRegisterStatsRoutes_WithJWT(t *testing.T) {
    os.Setenv("ADMIN_AUTH_MODE", "jwt")
    secret := "test-secret"
    jwt := auth.NewJWTManager(secret, 2*time.Hour)
    token, _ := jwt.GenerateToken(1, "super_admin")
    permSvc := permservice.NewPermissionService(permRepoStub{}, cache.NewMemory())
    roleSvc := roleservice.NewRoleService(roleRepoStub{}, cache.NewMemory())
    pm := mw.NewPermissionMiddleware(jwt, permSvc, roleSvc)
    gin.SetMode(gin.TestMode)
    r := gin.New()
    RegisterStatsRoutes(r, stats.NewStatsService(statsRepoStub{}), pm)

    authHeader := "Bearer " + token
    for _, path := range []string{
        "/admin/stats/dashboard",
        "/admin/stats/revenue-trend?days=7",
        "/admin/stats/user-growth?days=7",
        "/admin/stats/orders",
        "/admin/stats/top-players?limit=5",
        "/admin/stats/audit/overview?from=2025-01-01T00:00:00Z&to=2025-01-02T00:00:00Z",
        "/admin/stats/audit/trend?from=2025-01-01T00:00:00Z&to=2025-01-02T00:00:00Z&entity=order&action=create",
    } {
        w := httptest.NewRecorder()
        req := httptest.NewRequest(http.MethodGet, path, nil)
        req.Header.Set("Authorization", authHeader)
        r.ServeHTTP(w, req)
        if w.Code != http.StatusOK {
            t.Fatalf("%s -> %d", path, w.Code)
        }
    }
}
