package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionrepo "gamelink/internal/repository/commission"
	serviceitemrepo "gamelink/internal/repository/serviceitem"
	withdrawrepo "gamelink/internal/repository/withdraw"
)

func TestGetDashboardOverviewHandler(t *testing.T) {
	t.Parallel()
	userRepo := &dashboardUserRepo{users: []model.User{
		{Base: model.Base{ID: 1}, Name: "u1"},
		{Base: model.Base{ID: 2}, Name: "u2"},
		{Base: model.Base{ID: 3}, Name: "u3"},
	}}
	playerRepo := &fakePlayerRepoForHandler{
		listPaged: func(page, size int) ([]model.Player, int64, error) {
			return nil, 5, nil
		},
	}

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	lastMonth := now.AddDate(0, -1, 0)
	orders := []model.Order{
		{Base: model.Base{ID: 1, CreatedAt: now.Add(-2 * time.Hour)}, Status: model.OrderStatusCompleted, TotalPriceCents: 1_000},
		{Base: model.Base{ID: 2, CreatedAt: now.Add(-1 * time.Hour)}, Status: model.OrderStatusPending, TotalPriceCents: 500},
		{Base: model.Base{ID: 3, CreatedAt: yesterday}, Status: model.OrderStatusCompleted, TotalPriceCents: 2_000},
		{Base: model.Base{ID: 4, CreatedAt: lastMonth}, Status: model.OrderStatusCompleted, TotalPriceCents: 3_000},
	}
	orderRepo := &fakeOrderRepoForHandler{
		items: orders,
	}
	orderRepo.listFunc = func(opts repository.OrderListOptions) ([]model.Order, int64, error) {
		switch {
		case opts.DateFrom == nil:
			return append([]model.Order(nil), orders...), int64(len(orders)), nil
		default:
			filtered := make([]model.Order, 0, len(orders))
			for _, o := range orders {
				if opts.DateFrom != nil && o.CreatedAt.Before(*opts.DateFrom) {
					continue
				}
				if len(opts.Statuses) > 0 {
					match := false
					for _, status := range opts.Statuses {
						if o.Status == status {
							match = true
							break
						}
					}
					if !match {
						continue
					}
				}
				filtered = append(filtered, o)
			}
			return filtered, int64(len(filtered)), nil
		}
	}

	pendingWithdraw := model.Withdraw{ID: 1, Status: model.WithdrawStatusPending}
	approvedWithdraw := model.Withdraw{ID: 2, Status: model.WithdrawStatusApproved}
	withdrawRepo := &dashboardWithdrawRepo{
		items: []model.Withdraw{pendingWithdraw, approvedWithdraw},
	}

	serviceItemRepo := &dashboardServiceItemRepo{
		items: []model.ServiceItem{
			{ID: 1, IsActive: true},
			{ID: 2, IsActive: true},
			{ID: 3, IsActive: false},
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/admin/dashboard/overview", nil)
	c.Request = req

	getDashboardOverviewHandler(c, userRepo, playerRepo, orderRepo, withdrawRepo, serviceItemRepo)

	require.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[DashboardOverviewStats]
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.True(t, resp.Success)
	require.Equal(t, int64(3), resp.Data.TotalUsers)
	require.Equal(t, int64(5), resp.Data.TotalPlayers)
	require.Equal(t, int64(len(orders)), resp.Data.TotalOrders)
	require.Equal(t, int64(2), resp.Data.TodayOrders)
	require.Equal(t, int64(1_000), resp.Data.TodayRevenue)
	require.Equal(t, int64(1_000+2_000), resp.Data.MonthRevenue)
	require.Equal(t, int64(1), resp.Data.PendingWithdraws)
	require.Equal(t, int64(2), resp.Data.ActiveServices)
}

func TestGetRecentOrdersHandler_LimitValidation(t *testing.T) {
	t.Parallel()
	var receivedLimit int
	orderRepo := &fakeOrderRepoForHandler{
		listFunc: func(opts repository.OrderListOptions) ([]model.Order, int64, error) {
			receivedLimit = opts.PageSize
			return []model.Order{{Base: model.Base{ID: 1}}}, 1, nil
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/admin/dashboard/recent-orders?limit=200", nil)
	c.Request = req

	getRecentOrdersHandler(c, orderRepo)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, 10, receivedLimit, "invalid limit should fall back to default")

	// valid limit
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	req2 := httptest.NewRequest(http.MethodGet, "/admin/dashboard/recent-orders?limit=7", nil)
	c2.Request = req2
	getRecentOrdersHandler(c2, orderRepo)
	require.Equal(t, http.StatusOK, w2.Code)
	require.Equal(t, 7, receivedLimit)
}

func TestGetRecentWithdrawsHandler_Error(t *testing.T) {
	t.Parallel()
	withdrawRepo := &dashboardWithdrawRepo{
		err: assertErr("boom"),
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/admin/dashboard/recent-withdraws", nil)

	getRecentWithdrawsHandler(c, withdrawRepo)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetMonthlyRevenueHandler_CustomMonths(t *testing.T) {
	t.Parallel()
	fakeRepo := &dashboardCommissionRepo{
		stats: map[string]*commissionrepo.MonthlyStats{
			time.Now().Format("2006-01"): {
				TotalOrders:     10,
				TotalIncome:     5000,
				TotalCommission: 1000,
			},
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/admin/dashboard/monthly-revenue?months=3", nil)
	c.Request = req
	getMonthlyRevenueHandler(c, fakeRepo)
	require.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[map[string]interface{}]
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.True(t, resp.Success)

	revenue, ok := resp.Data["revenue"].([]interface{})
	require.True(t, ok)
	require.Len(t, revenue, 3)
}

// --- Fakes -------------------------------------------------------------

type dashboardUserRepo struct {
	users []model.User
}

func (r *dashboardUserRepo) List(ctx context.Context) ([]model.User, error) {
	return append([]model.User(nil), r.users...), nil
}

func (r *dashboardUserRepo) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	return append([]model.User(nil), r.users...), int64(len(r.users)), nil
}

func (r *dashboardUserRepo) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	return append([]model.User(nil), r.users...), int64(len(r.users)), nil
}

func (r *dashboardUserRepo) Get(ctx context.Context, id uint64) (*model.User, error) {
	for i := range r.users {
		if r.users[i].ID == id {
			user := r.users[i]
			return &user, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *dashboardUserRepo) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	for i := range r.users {
		if r.users[i].Phone == phone {
			user := r.users[i]
			return &user, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *dashboardUserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return r.GetByPhone(ctx, email)
}

func (r *dashboardUserRepo) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	return r.GetByPhone(ctx, phone)
}

func (r *dashboardUserRepo) Create(ctx context.Context, user *model.User) error {
	if user.ID == 0 {
		user.ID = uint64(len(r.users) + 1)
	}
	r.users = append(r.users, *user)
	return nil
}

func (r *dashboardUserRepo) Update(ctx context.Context, user *model.User) error {
	for i := range r.users {
		if r.users[i].ID == user.ID {
			r.users[i] = *user
			return nil
		}
	}
	return repository.ErrNotFound
}

func (r *dashboardUserRepo) Delete(ctx context.Context, id uint64) error {
	for i := range r.users {
		if r.users[i].ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}

type dashboardWithdrawRepo struct {
	items []model.Withdraw
	err   error
}

func (r *dashboardWithdrawRepo) Create(ctx context.Context, withdraw *model.Withdraw) error {
	return nil
}
func (r *dashboardWithdrawRepo) Get(ctx context.Context, id uint64) (*model.Withdraw, error) {
	return nil, repository.ErrNotFound
}
func (r *dashboardWithdrawRepo) Update(ctx context.Context, withdraw *model.Withdraw) error {
	return nil
}

func (r *dashboardWithdrawRepo) List(ctx context.Context, opts withdrawrepo.WithdrawListOptions) ([]model.Withdraw, int64, error) {
	if r.err != nil {
		return nil, 0, r.err
	}
	var filtered []model.Withdraw
	for _, w := range r.items {
		if opts.Status != nil && w.Status != *opts.Status {
			continue
		}
		filtered = append(filtered, w)
	}
	limit := opts.PageSize
	if limit <= 0 || limit > len(filtered) {
		limit = len(filtered)
	}
	return append([]model.Withdraw(nil), filtered[:limit]...), int64(len(filtered)), nil
}

func (r *dashboardWithdrawRepo) GetPlayerBalance(ctx context.Context, playerID uint64) (*withdrawrepo.PlayerBalance, error) {
	return &withdrawrepo.PlayerBalance{}, nil
}

type dashboardServiceItemRepo struct {
	items []model.ServiceItem
}

func (r *dashboardServiceItemRepo) Create(ctx context.Context, item *model.ServiceItem) error {
	return nil
}
func (r *dashboardServiceItemRepo) Get(ctx context.Context, id uint64) (*model.ServiceItem, error) {
	return nil, repository.ErrNotFound
}
func (r *dashboardServiceItemRepo) GetByCode(ctx context.Context, itemCode string) (*model.ServiceItem, error) {
	return nil, repository.ErrNotFound
}
func (r *dashboardServiceItemRepo) Update(ctx context.Context, item *model.ServiceItem) error {
	return nil
}
func (r *dashboardServiceItemRepo) Delete(ctx context.Context, id uint64) error { return nil }
func (r *dashboardServiceItemRepo) BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) error {
	return nil
}
func (r *dashboardServiceItemRepo) BatchUpdatePrice(ctx context.Context, ids []uint64, basePriceCents int64) error {
	return nil
}
func (r *dashboardServiceItemRepo) GetGifts(ctx context.Context, page, pageSize int) ([]model.ServiceItem, int64, error) {
	return nil, 0, nil
}
func (r *dashboardServiceItemRepo) GetGameServices(ctx context.Context, gameID uint64, subCategory *model.ServiceItemSubCategory) ([]model.ServiceItem, error) {
	return nil, nil
}

func (r *dashboardServiceItemRepo) List(ctx context.Context, opts serviceitemrepo.ServiceItemListOptions) ([]model.ServiceItem, int64, error) {
	var filtered []model.ServiceItem
	for _, item := range r.items {
		if opts.IsActive != nil && item.IsActive != *opts.IsActive {
			continue
		}
		filtered = append(filtered, item)
	}
	return filtered, int64(len(filtered)), nil
}

type assertErr string

func (e assertErr) Error() string { return string(e) }

type dashboardCommissionRepo struct {
	stats map[string]*commissionrepo.MonthlyStats
}

func (r *dashboardCommissionRepo) CreateRule(ctx context.Context, rule *model.CommissionRule) error {
	return nil
}
func (r *dashboardCommissionRepo) GetRule(ctx context.Context, id uint64) (*model.CommissionRule, error) {
	return nil, repository.ErrNotFound
}
func (r *dashboardCommissionRepo) GetDefaultRule(ctx context.Context) (*model.CommissionRule, error) {
	return nil, repository.ErrNotFound
}
func (r *dashboardCommissionRepo) GetRuleForOrder(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
	return nil, repository.ErrNotFound
}
func (r *dashboardCommissionRepo) ListRules(ctx context.Context, opts commissionrepo.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) {
	return nil, 0, nil
}
func (r *dashboardCommissionRepo) UpdateRule(ctx context.Context, rule *model.CommissionRule) error {
	return nil
}
func (r *dashboardCommissionRepo) DeleteRule(ctx context.Context, id uint64) error { return nil }
func (r *dashboardCommissionRepo) CreateRecord(ctx context.Context, record *model.CommissionRecord) error {
	return nil
}
func (r *dashboardCommissionRepo) GetRecord(ctx context.Context, id uint64) (*model.CommissionRecord, error) {
	return nil, repository.ErrNotFound
}
func (r *dashboardCommissionRepo) GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
	return nil, repository.ErrNotFound
}
func (r *dashboardCommissionRepo) ListRecords(ctx context.Context, opts commissionrepo.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) {
	return nil, 0, nil
}
func (r *dashboardCommissionRepo) UpdateRecord(ctx context.Context, record *model.CommissionRecord) error {
	return nil
}
func (r *dashboardCommissionRepo) CreateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	return nil
}
func (r *dashboardCommissionRepo) GetSettlement(ctx context.Context, id uint64) (*model.MonthlySettlement, error) {
	return nil, repository.ErrNotFound
}
func (r *dashboardCommissionRepo) GetSettlementByPlayerMonth(ctx context.Context, playerID uint64, month string) (*model.MonthlySettlement, error) {
	return nil, repository.ErrNotFound
}
func (r *dashboardCommissionRepo) ListSettlements(ctx context.Context, opts commissionrepo.SettlementListOptions) ([]model.MonthlySettlement, int64, error) {
	return nil, 0, nil
}
func (r *dashboardCommissionRepo) UpdateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	return nil
}
func (r *dashboardCommissionRepo) GetMonthlyStats(ctx context.Context, month string) (*commissionrepo.MonthlyStats, error) {
	if stats, ok := r.stats[month]; ok {
		return stats, nil
	}
	return nil, repository.ErrNotFound
}
func (r *dashboardCommissionRepo) GetPlayerMonthlyIncome(ctx context.Context, playerID uint64, month string) (int64, error) {
	return 0, nil
}

// helper to quickly create recorder/context pair with query parameters.
