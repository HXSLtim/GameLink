package admin

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "context"

    "github.com/gin-gonic/gin"

    "gamelink/internal/cache"
    "gamelink/internal/model"
    "gamelink/internal/repository"
    adminservice "gamelink/internal/service/admin"
)

type fakeUserRepo2 struct{ store map[uint64]*model.User; next uint64 }
func newFakeUserRepo() *fakeUserRepo2 { return &fakeUserRepo2{store: map[uint64]*model.User{}, next:1} }
func (f *fakeUserRepo2) List(ctx context.Context) ([]model.User, error) { out := make([]model.User,0,len(f.store)); for _, v := range f.store { out = append(out, *v) } ; return out, nil }
func (f *fakeUserRepo2) ListPaged(ctx context.Context, page int, pageSize int) ([]model.User, int64, error) { _=page; _=pageSize; out := make([]model.User,0,len(f.store)); for _, v := range f.store { out = append(out, *v) } ; return out, int64(len(out)), nil }
func (f *fakeUserRepo2) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) { _=opts; out := make([]model.User,0,len(f.store)); for _, v := range f.store { out = append(out, *v) } ; return out, int64(len(out)), nil }
func (f *fakeUserRepo2) Get(ctx context.Context, id uint64) (*model.User, error) { v := f.store[id]; if v==nil { return nil, repository.ErrNotFound } ; return v, nil }
func (f *fakeUserRepo2) GetByPhone(ctx context.Context, phone string) (*model.User, error) { _=phone; return nil, repository.ErrNotFound }
func (f *fakeUserRepo2) FindByEmail(ctx context.Context, email string) (*model.User, error) { _=email; return nil, repository.ErrNotFound }
func (f *fakeUserRepo2) FindByPhone(ctx context.Context, phone string) (*model.User, error) { _=phone; return nil, repository.ErrNotFound }
func (f *fakeUserRepo2) Create(ctx context.Context, u *model.User) error { _=ctx; u.ID = f.next; f.next++; f.store[u.ID] = u; return nil }
func (f *fakeUserRepo2) Update(ctx context.Context, u *model.User) error { _=ctx; f.store[u.ID] = u; return nil }
func (f *fakeUserRepo2) Delete(ctx context.Context, id uint64) error { _=ctx; delete(f.store, id); return nil }

type fakeRoleRepo2 struct{}
func (fakeRoleRepo2) List(ctx context.Context) ([]model.RoleModel, error) { return nil, nil }
func (fakeRoleRepo2) ListPaged(ctx context.Context, page int, pageSize int) ([]model.RoleModel, int64, error) { return nil, 0, nil }
func (fakeRoleRepo2) ListPagedWithFilter(ctx context.Context, page int, pageSize int, keyword string, isSystem *bool) ([]model.RoleModel, int64, error) { return nil, 0, nil }
func (fakeRoleRepo2) ListWithPermissions(ctx context.Context) ([]model.RoleModel, error) { return nil, nil }
func (fakeRoleRepo2) Get(ctx context.Context, id uint64) (*model.RoleModel, error) { return nil, repository.ErrNotFound }
func (fakeRoleRepo2) GetWithPermissions(ctx context.Context, id uint64) (*model.RoleModel, error) { return nil, repository.ErrNotFound }
func (fakeRoleRepo2) GetBySlug(ctx context.Context, slug string) (*model.RoleModel, error) { return nil, repository.ErrNotFound }
func (fakeRoleRepo2) Create(ctx context.Context, role *model.RoleModel) error { return nil }
func (fakeRoleRepo2) Update(ctx context.Context, role *model.RoleModel) error { return nil }
func (fakeRoleRepo2) Delete(ctx context.Context, id uint64) error { return nil }
func (fakeRoleRepo2) AssignPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error { return nil }
func (fakeRoleRepo2) AddPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error { return nil }
func (fakeRoleRepo2) RemovePermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error { return nil }
func (fakeRoleRepo2) AssignToUser(ctx context.Context, userID uint64, roleIDs []uint64) error { return nil }
func (fakeRoleRepo2) RemoveFromUser(ctx context.Context, userID uint64, roleIDs []uint64) error { return nil }
func (fakeRoleRepo2) ListByUserID(ctx context.Context, userID uint64) ([]model.RoleModel, error) { return nil, nil }
func (fakeRoleRepo2) CheckUserHasRole(ctx context.Context, userID uint64, roleSlug string) (bool, error) { return true, nil }

type dummyRepo2 struct{}
func (dummyRepo2) List(ctx context.Context) ([]model.Game, error) { return nil, nil }
func (dummyRepo2) ListPaged(ctx context.Context, page int, pageSize int) ([]model.Game, int64, error) { return nil, 0, nil }
func (dummyRepo2) Get(ctx context.Context, id uint64) (*model.Game, error) { return &model.Game{Name:"g"}, nil }
func (dummyRepo2) Create(ctx context.Context, game *model.Game) error { return nil }
func (dummyRepo2) Update(ctx context.Context, game *model.Game) error { return nil }
func (dummyRepo2) Delete(ctx context.Context, id uint64) error { return nil }

type dummyPlayerRepo2 struct{}
func (dummyPlayerRepo2) List(ctx context.Context) ([]model.Player, error) { return nil, nil }
func (dummyPlayerRepo2) ListPaged(ctx context.Context, page int, pageSize int) ([]model.Player, int64, error) { return nil, 0, nil }
func (dummyPlayerRepo2) Get(ctx context.Context, id uint64) (*model.Player, error) { return &model.Player{Nickname:"p"}, nil }
func (dummyPlayerRepo2) Create(ctx context.Context, player *model.Player) error { return nil }
func (dummyPlayerRepo2) Update(ctx context.Context, player *model.Player) error { return nil }
func (dummyPlayerRepo2) Delete(ctx context.Context, id uint64) error { return nil }

type dummyOrderRepo2 struct{}
func (dummyOrderRepo2) Create(ctx context.Context, order *model.Order) error { return nil }
func (dummyOrderRepo2) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) { return nil, 0, nil }
func (dummyOrderRepo2) Get(ctx context.Context, id uint64) (*model.Order, error) { return nil, repository.ErrNotFound }
func (dummyOrderRepo2) Update(ctx context.Context, order *model.Order) error { return nil }
func (dummyOrderRepo2) Delete(ctx context.Context, id uint64) error { return nil }

type dummyPaymentRepo2 struct{}
func (dummyPaymentRepo2) Create(ctx context.Context, payment *model.Payment) error { return nil }
func (dummyPaymentRepo2) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) { return nil, 0, nil }
func (dummyPaymentRepo2) Get(ctx context.Context, id uint64) (*model.Payment, error) { return nil, repository.ErrNotFound }
func (dummyPaymentRepo2) Update(ctx context.Context, payment *model.Payment) error { return nil }
func (dummyPaymentRepo2) Delete(ctx context.Context, id uint64) error { return nil }

func setupUserHandler() *UserHandler {
    svc := adminservice.NewAdminService(dummyRepo2{}, newFakeUserRepo(), dummyPlayerRepo2{}, dummyOrderRepo2{}, dummyPaymentRepo2{}, fakeRoleRepo2{}, cache.NewMemory())
    return NewUserHandler(svc)
}

func TestAdminUser_CreateUser_SuccessAndValidation(t *testing.T) {
    h := setupUserHandler()
    gin.SetMode(gin.TestMode)

    // success
    payload := CreateUserPayload{Phone:"13800138000", Email:"u@example.com", Password:"Abc123", Name:"U", Role:"user", Status:"active"}
    b, _ := json.Marshal(payload)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest(http.MethodPost, "/admin/users", bytes.NewReader(b))
    c.Request.Header.Set("Content-Type", "application/json")
    h.CreateUser(c)
    if w.Code != http.StatusCreated { t.Fatalf("%d", w.Code) }

    // invalid email
    bad := CreateUserPayload{Phone:"13800138000", Email:"bad", Password:"Abc123", Name:"U", Role:"user", Status:"active"}
    bb, _ := json.Marshal(bad)
    w2 := httptest.NewRecorder()
    c2, _ := gin.CreateTestContext(w2)
    c2.Request = httptest.NewRequest(http.MethodPost, "/admin/users", bytes.NewReader(bb))
    c2.Request.Header.Set("Content-Type", "application/json")
    h.CreateUser(c2)
    if w2.Code != http.StatusBadRequest { t.Fatalf("%d", w2.Code) }
}

func TestAdminUser_UpdateStatusAndRole(t *testing.T) {
    h := setupUserHandler()
    // pre-create user
    _, _ = h.svc.CreateUser(context.Background(), adminservice.CreateUserInput{Phone:"13800138000", Email:"u@example.com", Password:"Abc123", Name:"U", Role:model.RoleUser, Status:model.UserStatusActive})

    // update status
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key:"id", Value:"1"}}
    c.Request = httptest.NewRequest(http.MethodPut, "/admin/users/1/status", bytes.NewReader([]byte(`{"status":"inactive"}`)))
    c.Request.Header.Set("Content-Type", "application/json")
    h.UpdateUserStatus(c)
    if w.Code != http.StatusOK { t.Fatalf("%d", w.Code) }

    // update role
    w2 := httptest.NewRecorder()
    c2, _ := gin.CreateTestContext(w2)
    c2.Params = gin.Params{{Key:"id", Value:"1"}}
    c2.Request = httptest.NewRequest(http.MethodPut, "/admin/users/1/role", bytes.NewReader([]byte(`{"role":"admin"}`)))
    c2.Request.Header.Set("Content-Type", "application/json")
    h.UpdateUserRole(c2)
    if w2.Code != http.StatusOK { t.Fatalf("%d", w2.Code) }
}

func TestAdminUser_GetUser_InvalidID(t *testing.T) {
    h := setupUserHandler()
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key:"id", Value:"abc"}}
    h.GetUser(c)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}
