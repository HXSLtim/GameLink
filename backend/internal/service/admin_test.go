package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gamelink/internal/cache"
	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// ---- Fakes ----

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
func (f *fakeGameRepo) Update(ctx context.Context, g *model.Game) error { return nil }
func (f *fakeGameRepo) Delete(ctx context.Context, id uint64) error     { return nil }

type fakeUserRepo struct{ last *model.User }

func (f *fakeUserRepo) List(ctx context.Context) ([]model.User, error) { return nil, nil }
func (f *fakeUserRepo) ListPaged(ctx context.Context, page, size int) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (f *fakeUserRepo) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (f *fakeUserRepo) Get(ctx context.Context, id uint64) (*model.User, error) {
	if f.last != nil && f.last.ID == id {
		return f.last, nil
	}
	return nil, repository.ErrNotFound
}
func (f *fakeUserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return f.last, nil
}
func (f *fakeUserRepo) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	return f.last, nil
}
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

// ---- Tests ----

func TestService_ListGames_UsesCacheAndInvalidatesOnWrite(t *testing.T) {
	gRepo := &fakeGameRepo{items: []model.Game{{Base: model.Base{ID: 1}, Key: "lol", Name: "League"}}}
	s := NewAdminService(gRepo, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, &fakePaymentRepo{}, &fakeRoleRepo{}, cache.NewMemory())

	ctx := context.Background()

	// First call should hit repo.
	got1, err := s.ListGames(ctx)
	if err != nil {
		t.Fatalf("ListGames err: %v", err)
	}
	if gRepo.listCalls != 1 {
		t.Fatalf("expected 1 repo call, got %d", gRepo.listCalls)
	}
	if len(got1) != 1 || got1[0].Key != "lol" {
		t.Fatalf("unexpected games: %#v", got1)
	}

	// Second call should be served from cache.
	got2, err := s.ListGames(ctx)
	if err != nil {
		t.Fatalf("ListGames err: %v", err)
	}
	if gRepo.listCalls != 1 {
		t.Fatalf("expected cached result, repo calls=%d", gRepo.listCalls)
	}
	if len(got2) != 1 {
		t.Fatalf("unexpected cached games: %#v", got2)
	}

	// Write invalidates cache, next list hits repo again.
	_, err = s.CreateGame(ctx, CreateGameInput{Key: "dota2", Name: "DOTA2"})
	if err != nil {
		t.Fatalf("CreateGame err: %v", err)
	}
	_, _ = s.ListGames(ctx)
	if gRepo.listCalls != 2 {
		t.Fatalf("expected cache invalidation; repo calls=%d", gRepo.listCalls)
	}
}

func TestService_CreateGame_Validation(t *testing.T) {
	s := NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, &fakePaymentRepo{}, &fakeRoleRepo{}, cache.NewMemory())
	if _, err := s.CreateGame(context.Background(), CreateGameInput{Key: "", Name: ""}); err == nil {
		t.Fatalf("expected validation error for empty key/name")
	}
}

func TestService_UpdateOrder_Validation(t *testing.T) {
	now := time.Now()
	order := &model.Order{Base: model.Base{ID: 1}, Status: model.OrderStatusPending}
	oRepo := &fakeOrderRepo{obj: order}
	s := NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, oRepo, &fakePaymentRepo{}, &fakeRoleRepo{}, cache.NewMemory())

	// invalid status
	_, err := s.UpdateOrder(context.Background(), 1, UpdateOrderInput{Status: "bad", PriceCents: 1, Currency: model.CurrencyCNY})
	if err == nil {
		t.Fatalf("expected validation error for bad status")
	}

	// invalid currency
	_, err = s.UpdateOrder(context.Background(), 1, UpdateOrderInput{Status: model.OrderStatusConfirmed, PriceCents: 1, Currency: "XYZ"})
	if err == nil {
		t.Fatalf("expected validation error for bad currency")
	}

	// negative price
	_, err = s.UpdateOrder(context.Background(), 1, UpdateOrderInput{Status: model.OrderStatusConfirmed, PriceCents: -1, Currency: model.CurrencyCNY})
	if err == nil {
		t.Fatalf("expected validation error for negative price")
	}

	// valid update
	in := UpdateOrderInput{Status: model.OrderStatusCanceled, PriceCents: 100, Currency: model.CurrencyUSD, ScheduledStart: &now}
	out, err := s.UpdateOrder(context.Background(), 1, in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Status != model.OrderStatusCanceled || out.PriceCents != 100 || out.Currency != model.CurrencyUSD {
		t.Fatalf("unexpected order after update: %#v", out)
	}
}

func TestService_UpdatePayment_Validation(t *testing.T) {
	p := &model.Payment{Base: model.Base{ID: 1}, Status: model.PaymentStatusPending}
	pRepo := &fakePaymentRepo{obj: p}
	s := NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, pRepo, &fakeRoleRepo{}, cache.NewMemory())

	// invalid status
	_, err := s.UpdatePayment(context.Background(), 1, UpdatePaymentInput{Status: "oops"})
	if err == nil {
		t.Fatalf("expected validation error for bad status")
	}

	// valid update
	raw := json.RawMessage(`{"from":"test"}`)
	out, err := s.UpdatePayment(context.Background(), 1, UpdatePaymentInput{Status: model.PaymentStatusPaid, ProviderRaw: raw})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Status != model.PaymentStatusPaid {
		t.Fatalf("unexpected status: %v", out.Status)
	}
	if string(out.ProviderRaw) != string(raw) {
		t.Fatalf("unexpected raw: %s", string(out.ProviderRaw))
	}
}

func TestService_CreateUser_HashesPassword(t *testing.T) {
	uRepo := &fakeUserRepo{}
	s := NewAdminService(&fakeGameRepo{}, uRepo, &fakePlayerRepo{}, &fakeOrderRepo{}, &fakePaymentRepo{}, &fakeRoleRepo{}, cache.NewMemory())
	out, err := s.CreateUser(context.Background(), CreateUserInput{
		Phone: "1", Email: "a@b", Password: "secret123", Name: "alice", Role: model.RoleUser, Status: model.UserStatusActive,
	})
	if err != nil {
		t.Fatalf("CreateUser err: %v", err)
	}
	if out.PasswordHash == "" || out.PasswordHash == "secret123" {
		t.Fatalf("password not hashed properly: %q", out.PasswordHash)
	}
	if bcrypt.CompareHashAndPassword([]byte(out.PasswordHash), []byte("secret123")) != nil {
		t.Fatalf("stored hash does not match password")
	}
}

func TestService_CreateUser_PasswordTooShort(t *testing.T) {
	s := NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, &fakePaymentRepo{}, &fakeRoleRepo{}, cache.NewMemory())
	_, err := s.CreateUser(context.Background(), CreateUserInput{
		Name: "alice", Password: "123", Role: model.RoleUser, Status: model.UserStatusActive,
	})
	if err == nil {
		t.Fatalf("expected validation error for short password")
	}
}

func TestService_UpdateUser_PasswordOptional(t *testing.T) {
	uRepo := &fakeUserRepo{}
	s := NewAdminService(&fakeGameRepo{}, uRepo, &fakePlayerRepo{}, &fakeOrderRepo{}, &fakePaymentRepo{}, &fakeRoleRepo{}, cache.NewMemory())

	// create user first
	u, err := s.CreateUser(context.Background(), CreateUserInput{
		Name: "bob", Password: "password1", Role: model.RoleUser, Status: model.UserStatusActive,
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	original := u.PasswordHash

	// update with nil password
	u2, err := s.UpdateUser(context.Background(), u.ID, UpdateUserInput{Name: "bob", Role: model.RoleUser, Status: model.UserStatusActive})
	if err != nil {
		t.Fatalf("update user: %v", err)
	}
	if u2.PasswordHash != original {
		t.Fatalf("password should not change when nil")
	}

	// update with blank password pointer (should be ignored)
	empty := ""
	u3, err := s.UpdateUser(context.Background(), u.ID, UpdateUserInput{Name: "bob", Role: model.RoleUser, Status: model.UserStatusActive, Password: &empty})
	if err != nil {
		t.Fatalf("update user blank: %v", err)
	}
	if u3.PasswordHash != original {
		t.Fatalf("password should not change when blank")
	}
}

func TestService_OrderStateMachine(t *testing.T) {
	o := &model.Order{Base: model.Base{ID: 1}, Status: model.OrderStatusPending}
	oRepo := &fakeOrderRepo{obj: o}
	s := NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, oRepo, &fakePaymentRepo{}, &fakeRoleRepo{}, cache.NewMemory())

	// pending -> confirmed ok
	if _, err := s.UpdateOrder(context.Background(), 1, UpdateOrderInput{Status: model.OrderStatusConfirmed, PriceCents: 1, Currency: model.CurrencyCNY}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// confirmed -> pending not allowed
	if _, err := s.UpdateOrder(context.Background(), 1, UpdateOrderInput{Status: model.OrderStatusPending, PriceCents: 1, Currency: model.CurrencyCNY}); err == nil {
		t.Fatalf("expected invalid transition error")
	}
}

func TestService_PaymentStateMachine(t *testing.T) {
	p := &model.Payment{Base: model.Base{ID: 1}, Status: model.PaymentStatusPending}
	pRepo := &fakePaymentRepo{obj: p}
	s := NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, pRepo, &fakeRoleRepo{}, cache.NewMemory())

	// pending -> paid ok
	if _, err := s.UpdatePayment(context.Background(), 1, UpdatePaymentInput{Status: model.PaymentStatusPaid}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// paid -> failed not allowed
	if _, err := s.UpdatePayment(context.Background(), 1, UpdatePaymentInput{Status: model.PaymentStatusFailed}); err == nil {
		t.Fatalf("expected invalid transition error")
	}
}
