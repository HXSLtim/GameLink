package admin

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/cache"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/common"
)

type foRepo struct{ m map[uint64]*model.Order }

func newFo() *foRepo { return &foRepo{m: map[uint64]*model.Order{}} }
func (r *foRepo) Create(ctx context.Context, o *model.Order) error {
	if o.ID == 0 {
		o.ID = uint64(len(r.m) + 1)
	}
	r.m[o.ID] = o
	return nil
}
func (r *foRepo) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	out := []model.Order{}
	for _, v := range r.m {
		out = append(out, *v)
	}
	return out, int64(len(out)), nil
}
func (r *foRepo) Get(ctx context.Context, id uint64) (*model.Order, error) {
	v := r.m[id]
	if v == nil {
		return nil, repository.ErrNotFound
	}
	return v, nil
}
func (r *foRepo) Update(ctx context.Context, o *model.Order) error { r.m[o.ID] = o; return nil }
func (r *foRepo) Delete(ctx context.Context, id uint64) error      { delete(r.m, id); return nil }

type fpRepo struct{ m map[uint64]*model.Payment }

func newFp() *fpRepo { return &fpRepo{m: map[uint64]*model.Payment{}} }
func (r *fpRepo) Create(ctx context.Context, p *model.Payment) error {
	if p.ID == 0 {
		p.ID = uint64(len(r.m) + 1)
	}
	r.m[p.ID] = p
	return nil
}
func (r *fpRepo) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	out := []model.Payment{}
	for _, v := range r.m {
		if opts.OrderID != nil && v.OrderID == *opts.OrderID {
			out = append(out, *v)
		}
	}
	return out, int64(len(out)), nil
}
func (r *fpRepo) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	v := r.m[id]
	if v == nil {
		return nil, repository.ErrNotFound
	}
	return v, nil
}
func (r *fpRepo) Update(ctx context.Context, p *model.Payment) error { r.m[p.ID] = p; return nil }
func (r *fpRepo) Delete(ctx context.Context, id uint64) error        { delete(r.m, id); return nil }

type fuRepo struct{}

func (fuRepo) List(context.Context) ([]model.User, error)                       { return nil, nil }
func (fuRepo) ListPaged(context.Context, int, int) ([]model.User, int64, error) { return nil, 0, nil }
func (fuRepo) ListWithFilters(context.Context, repository.UserListOptions) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (fuRepo) Get(context.Context, uint64) (*model.User, error) {
	return &model.User{Base: model.Base{ID: 1}, Name: "u", Role: model.RoleUser, Status: model.UserStatusActive}, nil
}
func (fuRepo) GetByPhone(context.Context, string) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (fuRepo) FindByEmail(context.Context, string) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (fuRepo) FindByPhone(context.Context, string) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (fuRepo) Create(context.Context, *model.User) error { return nil }
func (fuRepo) Update(context.Context, *model.User) error { return nil }
func (fuRepo) Delete(context.Context, uint64) error      { return nil }

type fplRepo struct{}

func (fplRepo) List(context.Context) ([]model.Player, error) { return nil, nil }
func (fplRepo) ListPaged(context.Context, int, int) ([]model.Player, int64, error) {
	return []model.Player{{Base: model.Base{ID: 1}, UserID: 1, Nickname: "p"}}, 1, nil
}
func (fplRepo) Get(context.Context, uint64) (*model.Player, error) {
	return &model.Player{Base: model.Base{ID: 1}, UserID: 1, Nickname: "p"}, nil
}
func (fplRepo) GetByUserID(context.Context, uint64) (*model.Player, error) {
	return &model.Player{Base: model.Base{ID: 1}, UserID: 1, Nickname: "p"}, nil
}
func (fplRepo) Create(context.Context, *model.Player) error { return nil }
func (fplRepo) Update(context.Context, *model.Player) error { return nil }
func (fplRepo) Delete(context.Context, uint64) error        { return nil }

type frRole struct{}

func (frRole) List(context.Context) ([]model.RoleModel, error) { return nil, nil }
func (frRole) ListPaged(context.Context, int, int) ([]model.RoleModel, int64, error) {
	return nil, 0, nil
}
func (frRole) ListPagedWithFilter(context.Context, int, int, string, *bool) ([]model.RoleModel, int64, error) {
	return nil, 0, nil
}
func (frRole) ListWithPermissions(context.Context) ([]model.RoleModel, error) { return nil, nil }
func (frRole) Get(context.Context, uint64) (*model.RoleModel, error) {
	return nil, repository.ErrNotFound
}
func (frRole) GetWithPermissions(context.Context, uint64) (*model.RoleModel, error) {
	return nil, repository.ErrNotFound
}
func (frRole) GetBySlug(context.Context, string) (*model.RoleModel, error) {
	return &model.RoleModel{Base: model.Base{ID: 2}}, nil
}
func (frRole) Create(context.Context, *model.RoleModel) error                  { return nil }
func (frRole) Update(context.Context, *model.RoleModel) error                  { return nil }
func (frRole) Delete(context.Context, uint64) error                            { return nil }
func (frRole) AssignPermissions(context.Context, uint64, []uint64) error       { return nil }
func (frRole) AddPermissions(context.Context, uint64, []uint64) error          { return nil }
func (frRole) RemovePermissions(context.Context, uint64, []uint64) error       { return nil }
func (frRole) AssignToUser(context.Context, uint64, []uint64) error            { return nil }
func (frRole) RemoveFromUser(context.Context, uint64, []uint64) error          { return nil }
func (frRole) ListByUserID(context.Context, uint64) ([]model.RoleModel, error) { return nil, nil }
func (frRole) CheckUserHasRole(context.Context, uint64, string) (bool, error)  { return true, nil }

type txStub struct{ r *common.Repos }

func (t txStub) WithTx(ctx context.Context, fn func(r *common.Repos) error) error { return fn(t.r) }

type opLogsStub struct{}

func (opLogsStub) Append(ctx context.Context, log *model.OperationLog) error { return nil }
func (opLogsStub) ListByEntity(ctx context.Context, entityType string, entityID uint64, opts repository.OperationLogListOptions) ([]model.OperationLog, int64, error) {
	return nil, 0, nil
}

type tagsStub struct{}

func (tagsStub) GetTags(ctx context.Context, playerID uint64) ([]string, error) {
	return []string{"a"}, nil
}
func (tagsStub) ReplaceTags(ctx context.Context, playerID uint64, tags []string) error { return nil }

func TestAdminService_UpdateOrder_InvalidTransition(t *testing.T) {
	orders := newFo()
	orders.m[1] = &model.Order{Base: model.Base{ID: 1}, Status: model.OrderStatusCanceled, TotalPriceCents: 100, Currency: model.CurrencyCNY}
	svc := NewAdminService(nil, fuRepo{}, fplRepo{}, orders, newFp(), frRole{}, cache.NewMemory())
	_, err := svc.UpdateOrder(context.Background(), 1, UpdateOrderInput{Status: model.OrderStatusConfirmed, TotalPriceCents: 100, Currency: model.CurrencyCNY})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestAdminService_RefundOrder_Validation(t *testing.T) {
	orders := newFo()
	orders.m[2] = &model.Order{Base: model.Base{ID: 2}, Status: model.OrderStatusConfirmed, TotalPriceCents: 100, Currency: model.CurrencyCNY}
	svc := NewAdminService(nil, fuRepo{}, fplRepo{}, orders, newFp(), frRole{}, cache.NewMemory())
	_, err := svc.RefundOrder(context.Background(), 2, RefundOrderInput{Reason: ""})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestAdminService_CreatePayment_Invalid(t *testing.T) {
	orders := newFo()
	orders.m[3] = &model.Order{Base: model.Base{ID: 3}, UserID: 1, TotalPriceCents: 100, Currency: model.CurrencyCNY}
	svc := NewAdminService(nil, fuRepo{}, fplRepo{}, orders, newFp(), frRole{}, cache.NewMemory())
	_, err := svc.CreatePayment(context.Background(), CreatePaymentInput{OrderID: 3, UserID: 1, Method: "", AmountCents: 100, Currency: model.CurrencyCNY})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestAdminService_CapturePayment_FromPending(t *testing.T) {
	orders := newFo()
	orders.m[4] = &model.Order{Base: model.Base{ID: 4}, UserID: 1, TotalPriceCents: 100, Currency: model.CurrencyCNY}
	pays := newFp()
	_ = pays.Create(context.Background(), &model.Payment{OrderID: 4, UserID: 1, Method: model.PaymentMethodWeChat, AmountCents: 100, Currency: model.CurrencyCNY, Status: model.PaymentStatusPending})
	svc := NewAdminService(nil, fuRepo{}, fplRepo{}, orders, pays, frRole{}, cache.NewMemory())
	_, err := svc.CapturePayment(context.Background(), 1, CapturePaymentInput{ProviderTradeNo: "x"})
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestAdminService_GetOrderRefunds_Summary(t *testing.T) {
	orders := newFo()
	refundedAt := time.Now().UTC()
	orders.m[5] = &model.Order{Base: model.Base{ID: 5}, UserID: 1, TotalPriceCents: 100, Currency: model.CurrencyCNY, Status: model.OrderStatusRefunded, RefundAmountCents: 80, RefundReason: "r", RefundedAt: &refundedAt}
	pays := newFp()
	svc := NewAdminService(nil, fuRepo{}, fplRepo{}, orders, pays, frRole{}, cache.NewMemory())
	items, err := svc.GetOrderRefunds(context.Background(), 5)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if len(items) == 0 {
		t.Fatalf("no items")
	}
}

func TestAdminService_UpdatePlayerSkillTags_WithTx(t *testing.T) {
	svc := NewAdminService(nil, fuRepo{}, fplRepo{}, nil, nil, frRole{}, cache.NewMemory())
	svc.SetTxManager(txStub{r: &common.Repos{Players: fplRepo{}, Tags: tagsStub{}, OpLogs: opLogsStub{}}})
	if err := svc.UpdatePlayerSkillTags(context.Background(), 1, []string{"a", "b"}); err != nil {
		t.Fatalf("%v", err)
	}
}
