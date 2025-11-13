package user

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionrepo "gamelink/internal/repository/commission"
	ordersvc "gamelink/internal/service/order"
)

type filterOrderRepo struct{ items []model.Order }

func (r *filterOrderRepo) Create(ctx context.Context, o *model.Order) error {
	r.items = append(r.items, *o)
	return nil
}
func (r *filterOrderRepo) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	out := make([]model.Order, 0, len(r.items))
	for _, o := range r.items {
		if opts.UserID != nil && o.UserID != *opts.UserID {
			continue
		}
		if len(opts.Statuses) > 0 {
			match := false
			for _, st := range opts.Statuses {
				if o.Status == st {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		out = append(out, o)
	}
	return out, int64(len(out)), nil
}
func (r *filterOrderRepo) Get(ctx context.Context, id uint64) (*model.Order, error) {
	for i := range r.items {
		if r.items[i].ID == id {
			o := r.items[i]
			return &o, nil
		}
	}
	return nil, repository.ErrNotFound
}
func (r *filterOrderRepo) Update(ctx context.Context, o *model.Order) error {
	for i := range r.items {
		if r.items[i].ID == o.ID {
			r.items[i] = *o
			break
		}
	}
	return nil
}
func (r *filterOrderRepo) Delete(ctx context.Context, id uint64) error { return nil }

type dummyPlayers struct{}

func (dummyPlayers) List(context.Context) ([]model.Player, error) { return nil, nil }
func (dummyPlayers) ListPaged(context.Context, int, int) ([]model.Player, int64, error) {
	return []model.Player{{Base: model.Base{ID: 1}, UserID: 1, Nickname: "p", HourlyRateCents: 1000}}, 1, nil
}
func (dummyPlayers) Get(context.Context, uint64) (*model.Player, error) {
	return &model.Player{Base: model.Base{ID: 1}, UserID: 1, Nickname: "p", HourlyRateCents: 1000}, nil
}
func (dummyPlayers) GetByUserID(context.Context, uint64) (*model.Player, error) {
	return &model.Player{Base: model.Base{ID: 1}, UserID: 1, Nickname: "p", HourlyRateCents: 1000}, nil
}
func (dummyPlayers) Create(context.Context, *model.Player) error { return nil }
func (dummyPlayers) Update(context.Context, *model.Player) error { return nil }
func (dummyPlayers) Delete(context.Context, uint64) error        { return nil }

type dummyUsers struct{}

func (dummyUsers) List(context.Context) ([]model.User, error) { return nil, nil }
func (dummyUsers) ListPaged(context.Context, int, int) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (dummyUsers) ListWithFilters(context.Context, repository.UserListOptions) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (dummyUsers) Get(context.Context, uint64) (*model.User, error) {
	return &model.User{Base: model.Base{ID: 1}, AvatarURL: "a"}, nil
}
func (dummyUsers) GetByPhone(context.Context, string) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (dummyUsers) FindByEmail(context.Context, string) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (dummyUsers) FindByPhone(context.Context, string) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (dummyUsers) Create(context.Context, *model.User) error { return nil }
func (dummyUsers) Update(context.Context, *model.User) error { return nil }
func (dummyUsers) Delete(context.Context, uint64) error      { return nil }

type dummyGames struct{}

func (dummyGames) List(context.Context) ([]model.Game, error) { return nil, nil }
func (dummyGames) ListPaged(context.Context, int, int) ([]model.Game, int64, error) {
	return nil, 0, nil
}
func (dummyGames) Get(context.Context, uint64) (*model.Game, error) {
	return &model.Game{Base: model.Base{ID: 1}, Name: "g"}, nil
}
func (dummyGames) Create(context.Context, *model.Game) error { return nil }
func (dummyGames) Update(context.Context, *model.Game) error { return nil }
func (dummyGames) Delete(context.Context, uint64) error      { return nil }

type dummyPayments struct{}

func (dummyPayments) Create(context.Context, *model.Payment) error { return nil }
func (dummyPayments) List(context.Context, repository.PaymentListOptions) ([]model.Payment, int64, error) {
	return nil, 0, nil
}
func (dummyPayments) Get(context.Context, uint64) (*model.Payment, error) {
	return nil, repository.ErrNotFound
}
func (dummyPayments) Update(context.Context, *model.Payment) error { return nil }
func (dummyPayments) Delete(context.Context, uint64) error         { return nil }

type dummyReviews struct{}

func (dummyReviews) List(context.Context, repository.ReviewListOptions) ([]model.Review, int64, error) {
	return nil, 0, nil
}
func (dummyReviews) Get(context.Context, uint64) (*model.Review, error) {
	return nil, repository.ErrNotFound
}
func (dummyReviews) Create(context.Context, *model.Review) error { return nil }
func (dummyReviews) Update(context.Context, *model.Review) error { return nil }
func (dummyReviews) Delete(context.Context, uint64) error        { return nil }

type dummyCommissions struct{}

func (dummyCommissions) CreateRule(context.Context, *model.CommissionRule) error { return nil }
func (dummyCommissions) GetRule(context.Context, uint64) (*model.CommissionRule, error) {
	return nil, repository.ErrNotFound
}
func (dummyCommissions) GetDefaultRule(context.Context) (*model.CommissionRule, error) {
	return &model.CommissionRule{Rate: 20}, nil
}
func (dummyCommissions) GetRuleForOrder(context.Context, *uint64, *uint64, *string) (*model.CommissionRule, error) {
	return nil, repository.ErrNotFound
}
func (dummyCommissions) ListRules(context.Context, commissionrepo.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) {
	return nil, 0, nil
}
func (dummyCommissions) UpdateRule(context.Context, *model.CommissionRule) error     { return nil }
func (dummyCommissions) DeleteRule(context.Context, uint64) error                    { return nil }
func (dummyCommissions) CreateRecord(context.Context, *model.CommissionRecord) error { return nil }
func (dummyCommissions) GetRecord(context.Context, uint64) (*model.CommissionRecord, error) {
	return nil, repository.ErrNotFound
}
func (dummyCommissions) GetRecordByOrderID(context.Context, uint64) (*model.CommissionRecord, error) {
	return nil, repository.ErrNotFound
}
func (dummyCommissions) ListRecords(context.Context, commissionrepo.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) {
	return nil, 0, nil
}
func (dummyCommissions) UpdateRecord(context.Context, *model.CommissionRecord) error      { return nil }
func (dummyCommissions) CreateSettlement(context.Context, *model.MonthlySettlement) error { return nil }
func (dummyCommissions) GetSettlement(context.Context, uint64) (*model.MonthlySettlement, error) {
	return nil, repository.ErrNotFound
}
func (dummyCommissions) GetSettlementByPlayerMonth(context.Context, uint64, string) (*model.MonthlySettlement, error) {
	return nil, repository.ErrNotFound
}
func (dummyCommissions) ListSettlements(context.Context, commissionrepo.SettlementListOptions) ([]model.MonthlySettlement, int64, error) {
	return nil, 0, nil
}
func (dummyCommissions) UpdateSettlement(context.Context, *model.MonthlySettlement) error { return nil }
func (dummyCommissions) GetMonthlyStats(context.Context, string) (*commissionrepo.MonthlyStats, error) {
	return &commissionrepo.MonthlyStats{}, nil
}
func (dummyCommissions) GetPlayerMonthlyIncome(context.Context, uint64, string) (int64, error) {
	return 0, nil
}

func setupOrderFiltersRouter(repo *filterOrderRepo) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_id", uint64(1)); c.Next() })
	svc := ordersvc.NewOrderService(repo, dummyPlayers{}, dummyUsers{}, dummyGames{}, dummyPayments{}, dummyReviews{}, dummyCommissions{})
	assignSvc := newAssignmentServiceStub(repo, dummyPlayers{})
	RegisterOrderRoutes(r, svc, assignSvc, func(c *gin.Context) { c.Next() })
	return r
}

func TestUserOrder_ListFilters(t *testing.T) {
	repo := &filterOrderRepo{items: []model.Order{
		{Base: model.Base{ID: 1}, UserID: 1, Title: "A", Status: model.OrderStatusPending},
		{Base: model.Base{ID: 2}, UserID: 1, Title: "B", Status: model.OrderStatusConfirmed},
		{Base: model.Base{ID: 3}, UserID: 2, Title: "C", Status: model.OrderStatusPending},
	}}
	r := setupOrderFiltersRouter(repo)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/user/orders?status=pending&page=1&pageSize=10", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("%d", w.Code)
	}
}
