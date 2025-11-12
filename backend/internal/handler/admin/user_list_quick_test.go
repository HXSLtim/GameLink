package admin

import (
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

type fakeUsersRepo5 struct{ items []model.User }
func (r *fakeUsersRepo5) List(ctx context.Context) ([]model.User, error) { return r.items, nil }
func (r *fakeUsersRepo5) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) { return r.items, int64(len(r.items)), nil }
func (r *fakeUsersRepo5) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) { return r.items, int64(len(r.items)), nil }
func (r *fakeUsersRepo5) Get(ctx context.Context, id uint64) (*model.User, error) { for i:=range r.items{ if r.items[i].ID==id { u:=r.items[i]; return &u, nil } } ; return nil, repository.ErrNotFound }
func (r *fakeUsersRepo5) GetByPhone(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (r *fakeUsersRepo5) FindByEmail(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (r *fakeUsersRepo5) FindByPhone(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (r *fakeUsersRepo5) Create(context.Context, *model.User) error { return nil }
func (r *fakeUsersRepo5) Update(context.Context, *model.User) error { return nil }
func (r *fakeUsersRepo5) Delete(context.Context, uint64) error { return nil }

type fakeOrdersRepo5 struct{ items []model.Order }
func (r *fakeOrdersRepo5) Create(context.Context, *model.Order) error { return nil }
func (r *fakeOrdersRepo5) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) { return r.items, int64(len(r.items)), nil }
func (r *fakeOrdersRepo5) Get(ctx context.Context, id uint64) (*model.Order, error) { for i:=range r.items{ if r.items[i].ID==id { o:=r.items[i]; return &o, nil } } ; return nil, repository.ErrNotFound }
func (r *fakeOrdersRepo5) Update(context.Context, *model.Order) error { return nil }
func (r *fakeOrdersRepo5) Delete(context.Context, uint64) error { return nil }

func setupUserHandler5() *UserHandler {
    u := []model.User{{Base: model.Base{ID:1}, Name:"U", Role:model.RoleUser, Status:model.UserStatusActive}}
    o := []model.Order{{Base: model.Base{ID:1}, UserID:1, TotalPriceCents:1000, Currency:model.CurrencyCNY, Status:model.OrderStatusPending}}
    svc := adminservice.NewAdminService(nil, &fakeUsersRepo5{items:u}, nil, &fakeOrdersRepo5{items:o}, nil, nil, cache.NewMemory())
    return NewUserHandler(svc)
}

func TestAdminUser_ListUsers_FilterAndDates(t *testing.T) {
    gin.SetMode(gin.TestMode)
    h := setupUserHandler5()
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest(http.MethodGet, "/admin/users?page=1&page_size=10&role=admin,user&status=active&date_from=2025-01-01&date_to=2025-01-02&keyword=u", nil)
    h.ListUsers(c)
    if w.Code != http.StatusOK { t.Fatalf("%d", w.Code) }
}

func TestAdminUser_ListUsers_InvalidDate(t *testing.T) {
    h := setupUserHandler5()
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest(http.MethodGet, "/admin/users?date_from=bad", nil)
    h.ListUsers(c)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

func TestAdminUser_ListUserOrders(t *testing.T) {
    h := setupUserHandler5()
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key:"id", Value:"1"}}
    c.Request = httptest.NewRequest(http.MethodGet, "/admin/users/1/orders?page=1&page_size=10&status=pending", nil)
    h.ListUserOrders(c)
    if w.Code != http.StatusOK { t.Fatalf("%d", w.Code) }
}

