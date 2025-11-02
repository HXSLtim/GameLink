package permission

import (
    "context"
    "encoding/json"
    "reflect"
    "testing"
    "time"

    "gamelink/internal/cache"
    "gamelink/internal/model"
    "gamelink/internal/repository"
)

// mockPermissionRepository 是一个简单的mock实现。
type mockPermissionRepository struct {
	repository.PermissionRepository
}

// mockCache 是一个简单的mock实现。
type mockCache struct {
	cache.Cache
}

// TestNewPermissionService 测试构造函数。
func TestNewPermissionService(t *testing.T) {
	permRepo := &mockPermissionRepository{}
	cache := &mockCache{}

	svc := NewPermissionService(permRepo, cache)

	if svc == nil {
		t.Fatal("NewPermissionService returned nil")
	}

	if svc.permissions != permRepo {
		t.Error("permissions repository not set correctly")
	}

	if svc.cache != cache {
		t.Error("cache not set correctly")
	}
}

// ---- Fakes for tests ----
type fakeCache struct{
    cache.Cache
    values map[string]string
    deleted []string
}

func newFakeCache() *fakeCache { return &fakeCache{values: make(map[string]string)} }
func (c *fakeCache) Get(ctx context.Context, key string) (string, bool, error) {
    v, ok := c.values[key]
    return v, ok, nil
}
func (c *fakeCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
    c.values[key] = value
    return nil
}
func (c *fakeCache) Delete(ctx context.Context, key string) error {
    c.deleted = append(c.deleted, key)
    delete(c.values, key)
    return nil
}
func (c *fakeCache) Close(context.Context) error { return nil }

type fakePermRepo struct{
    repository.PermissionRepository
    byID map[uint64]*model.Permission
    byMethodPath map[string]*model.Permission
    listAll []model.Permission
    groups map[string][]model.Permission
    byRole map[uint64][]model.Permission
    byUser map[uint64][]model.Permission
}

func newFakePermRepo() *fakePermRepo { return &fakePermRepo{byID: make(map[uint64]*model.Permission), byMethodPath: make(map[string]*model.Permission), groups: make(map[string][]model.Permission), byRole: make(map[uint64][]model.Permission), byUser: make(map[uint64][]model.Permission)} }

func (r *fakePermRepo) List(ctx context.Context) ([]model.Permission, error) { return r.listAll, nil }
func (r *fakePermRepo) ListPaged(ctx context.Context, page, pageSize int) ([]model.Permission, int64, error) { return r.listAll, int64(len(r.listAll)), nil }
func (r *fakePermRepo) ListByGroup(ctx context.Context) (map[string][]model.Permission, error) { return r.groups, nil }
func (r *fakePermRepo) ListGroups(ctx context.Context) ([]string, error) { keys := make([]string,0,len(r.groups)); for k := range r.groups { keys = append(keys,k) }; return keys, nil }
func (r *fakePermRepo) Get(ctx context.Context, id uint64) (*model.Permission, error) { if p, ok := r.byID[id]; ok { return p, nil }; return nil, repository.ErrNotFound }
func (r *fakePermRepo) GetByResource(ctx context.Context, resource, action string) (*model.Permission, error) { return nil, repository.ErrNotFound }
func (r *fakePermRepo) GetByMethodAndPath(ctx context.Context, method, path string) (*model.Permission, error) { key := method+" "+path; if p, ok := r.byMethodPath[key]; ok { return p, nil }; return nil, repository.ErrNotFound }
func (r *fakePermRepo) Create(ctx context.Context, perm *model.Permission) error { if perm.ID==0 { perm.ID = uint64(len(r.byID)+1) }; r.byID[perm.ID]=perm; r.byMethodPath[string(perm.Method)+" "+perm.Path]=perm; r.listAll = append(r.listAll,*perm); if perm.Group!="" { r.groups[perm.Group] = append(r.groups[perm.Group], *perm) }; return nil }
func (r *fakePermRepo) Update(ctx context.Context, perm *model.Permission) error { r.byID[perm.ID]=perm; r.byMethodPath[string(perm.Method)+" "+perm.Path]=perm; return nil }
func (r *fakePermRepo) UpsertByMethodPath(ctx context.Context, perm *model.Permission) error { key := string(perm.Method)+" "+perm.Path; if existing, ok := r.byMethodPath[key]; ok { perm.ID = existing.ID; return r.Update(ctx, perm) }; return r.Create(ctx, perm) }
func (r *fakePermRepo) Delete(ctx context.Context, id uint64) error { delete(r.byID,id); return nil }
func (r *fakePermRepo) ListByRoleID(ctx context.Context, roleID uint64) ([]model.Permission, error) { return r.byRole[roleID], nil }
func (r *fakePermRepo) ListByUserID(ctx context.Context, userID uint64) ([]model.Permission, error) { return r.byUser[userID], nil }

// ---- Tests ----
func TestPermission_Create_ValidationAndDuplicate(t *testing.T) {
    repo := newFakePermRepo()
    c := newFakeCache()
    svc := NewPermissionService(repo, c)
    ctx := context.Background()

    // missing required
    if err := svc.CreatePermission(ctx, &model.Permission{Path: "/x"}); err == nil { t.Fatal("expected validation error: missing method") }
    if err := svc.CreatePermission(ctx, &model.Permission{Method: model.HTTPMethodGET}); err == nil { t.Fatal("expected validation error: missing path") }

    // duplicate method+path
    p := &model.Permission{Method: model.HTTPMethodGET, Path: "/a"}
    _ = repo.Create(ctx, p)
    if err := svc.CreatePermission(ctx, &model.Permission{Method: model.HTTPMethodGET, Path: "/a"}); err == nil { t.Fatal("expected duplicate error") }
}

func TestPermission_Create_Update_Delete_InvalidateCache(t *testing.T) {
    repo := newFakePermRepo(); c := newFakeCache(); svc := NewPermissionService(repo, c); ctx := context.Background()
    p := &model.Permission{Method: model.HTTPMethodPOST, Path: "/p", Group: "g"}
    if err := svc.CreatePermission(ctx, p); err != nil { t.Fatalf("create err: %v", err) }
    if len(c.deleted)==0 { t.Fatalf("expected cache invalidation on create") }

    p.Description = "updated"
    if err := svc.UpdatePermission(ctx, p); err != nil { t.Fatalf("update err: %v", err) }
    if len(c.deleted) < 2 { t.Fatalf("expected cache invalidation on update") }

    if err := svc.DeletePermission(ctx, p.ID); err != nil { t.Fatalf("delete err: %v", err) }
    if len(c.deleted) < 3 { t.Fatalf("expected cache invalidation on delete") }
}

func TestPermission_Update_RequireID(t *testing.T) {
    svc := NewPermissionService(newFakePermRepo(), newFakeCache())
    if err := svc.UpdatePermission(context.Background(), &model.Permission{}); err == nil { t.Fatal("expected validation error for missing ID") }
}

func TestPermission_Upsert_InvalidateCache(t *testing.T) {
    repo := newFakePermRepo(); c := newFakeCache(); svc := NewPermissionService(repo, c)
    p := &model.Permission{Method: model.HTTPMethodGET, Path: "/u"}
    if err := svc.UpsertPermission(context.Background(), p); err != nil { t.Fatalf("upsert err: %v", err) }
    if len(c.deleted)==0 { t.Fatal("expected cache invalidation on upsert") }
}

func TestPermission_ListByRoleID_CacheHitAndMiss(t *testing.T) {
    repo := newFakePermRepo(); c := newFakeCache(); svc := NewPermissionService(repo, c); ctx := context.Background()
    roleID := uint64(1)
    perms := []model.Permission{{Method: model.HTTPMethodGET, Path: "/a"}}
    // cache miss -> set
    repo.byRole[roleID] = perms
    got, err := svc.ListPermissionsByRoleID(ctx, roleID)
    if err != nil { t.Fatalf("err: %v", err) }
    if !reflect.DeepEqual(got, perms) { t.Fatalf("unexpected perms on miss") }
    // now cache hit
    key := "admin:permissions:role:1"
    data, _ := json.Marshal(perms)
    _ = c.Set(ctx, key, string(data), time.Minute)
    got2, err := svc.ListPermissionsByRoleID(ctx, roleID)
    if err != nil { t.Fatalf("err: %v", err) }
    if !reflect.DeepEqual(got2, perms) { t.Fatalf("unexpected perms on hit") }
}

func TestPermission_ListByUserID_CacheHitAndMiss(t *testing.T) {
    repo := newFakePermRepo(); c := newFakeCache(); svc := NewPermissionService(repo, c); ctx := context.Background()
    userID := uint64(9)
    perms := []model.Permission{{Method: model.HTTPMethodPOST, Path: "/b"}}
    // miss
    repo.byUser[userID] = perms
    got, err := svc.ListPermissionsByUserID(ctx, userID)
    if err != nil || !reflect.DeepEqual(got, perms) { t.Fatalf("unexpected result on miss: %v", err) }
    // hit
    key := "admin:permissions:user:9"
    data, _ := json.Marshal(perms)
    _ = c.Set(ctx, key, string(data), time.Minute)
    got2, _ := svc.ListPermissionsByUserID(ctx, userID)
    if !reflect.DeepEqual(got2, perms) { t.Fatalf("unexpected perms on hit") }
}

func TestPermission_CheckUserHasPermission(t *testing.T) {
    repo := newFakePermRepo(); c := newFakeCache(); svc := NewPermissionService(repo, c)
    repo.byUser[5] = []model.Permission{{Method: model.HTTPMethodGET, Path: "/ok"}}
    ok, err := svc.CheckUserHasPermission(context.Background(), 5, model.HTTPMethodGET, "/ok")
    if err != nil || !ok { t.Fatalf("expected true, got %v err=%v", ok, err) }
    ok, _ = svc.CheckUserHasPermission(context.Background(), 5, model.HTTPMethodGET, "/no")
    if ok { t.Fatalf("expected false") }
}

func TestPermission_ListAndGroups(t *testing.T) {
    repo := newFakePermRepo(); c := newFakeCache(); svc := NewPermissionService(repo, c)
    p := &model.Permission{Method: model.HTTPMethodGET, Path: "/g", Group: "grp"}
    _ = repo.Create(context.Background(), p)
    list, err := svc.ListPermissions(context.Background())
    if err != nil || len(list) != 1 { t.Fatalf("list err=%v len=%d", err, len(list)) }
    grouped, err := svc.ListPermissionsByGroup(context.Background(), "grp")
    if err != nil || len(grouped) != 1 { t.Fatalf("group err=%v len=%d", err, len(grouped)) }
    groups, _ := svc.ListPermissionGroups(context.Background())
    if len(groups) != 1 || groups[0] != "grp" { t.Fatalf("groups unexpected: %v", groups) }
}

func TestPermission_Paged_Normalization(t *testing.T) {
    repo := newFakePermRepo(); repo.listAll = []model.Permission{{},{},{}}
    svc := NewPermissionService(repo, newFakeCache())
    // page <1 -> 1 ; pageSize out of range -> 20
    list, total, err := svc.ListPermissionsPaged(context.Background(), 0, 1000)
    if err != nil || len(list) != 3 || total != 3 { t.Fatalf("unexpected paged result: err=%v total=%d", err, total) }
}

func TestPermission_GetPermission(t *testing.T) {
    repo := newFakePermRepo(); svc := NewPermissionService(repo, newFakeCache())
    p := &model.Permission{Base: model.Base{ID: 42}, Method: model.HTTPMethodGET, Path: "/x"}
    repo.byID[p.ID] = p
    got, err := svc.GetPermission(context.Background(), 42)
    if err != nil || got == nil || got.ID != 42 { t.Fatalf("get err=%v", err) }
}
