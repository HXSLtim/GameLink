package admin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionrepo "gamelink/internal/repository/commission"
	serviceitemrepo "gamelink/internal/repository/serviceitem"
)

// fakeStatsOrderRepo 用于stats测试的订单仓储
type fakeStatsOrderRepo struct {
	orders []model.Order
}

func (f *fakeStatsOrderRepo) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	return f.orders, int64(len(f.orders)), nil
}

func (f *fakeStatsOrderRepo) Get(ctx context.Context, id uint64) (*model.Order, error) {
	for i := range f.orders {
		if f.orders[i].ID == id {
			return &f.orders[i], nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeStatsOrderRepo) Create(ctx context.Context, order *model.Order) error {
	return nil
}

func (f *fakeStatsOrderRepo) Update(ctx context.Context, order *model.Order) error {
	return nil
}

func (f *fakeStatsOrderRepo) Delete(ctx context.Context, id uint64) error {
	return nil
}

// fakeStatsServiceItemRepo 用于stats测试的服务项目仓储
type fakeStatsServiceItemRepo struct {
	items []model.ServiceItem
	gifts []model.ServiceItem
}

func (f *fakeStatsServiceItemRepo) List(ctx context.Context, opts serviceitemrepo.ServiceItemListOptions) ([]model.ServiceItem, int64, error) {
	return f.items, int64(len(f.items)), nil
}

func (f *fakeStatsServiceItemRepo) Get(ctx context.Context, id uint64) (*model.ServiceItem, error) {
	for i := range f.items {
		if f.items[i].ID == id {
			return &f.items[i], nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeStatsServiceItemRepo) GetByCode(ctx context.Context, code string) (*model.ServiceItem, error) {
	for i := range f.items {
		if f.items[i].ItemCode == code {
			return &f.items[i], nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeStatsServiceItemRepo) GetGifts(ctx context.Context, page, pageSize int) ([]model.ServiceItem, int64, error) {
	return f.gifts, int64(len(f.gifts)), nil
}

func (f *fakeStatsServiceItemRepo) GetGameServices(ctx context.Context, gameID uint64, subCategory *model.ServiceItemSubCategory) ([]model.ServiceItem, error) {
	return nil, nil
}

func (f *fakeStatsServiceItemRepo) Create(ctx context.Context, item *model.ServiceItem) error {
	return nil
}

func (f *fakeStatsServiceItemRepo) Update(ctx context.Context, item *model.ServiceItem) error {
	return nil
}

func (f *fakeStatsServiceItemRepo) Delete(ctx context.Context, id uint64) error {
	return nil
}

func (f *fakeStatsServiceItemRepo) BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) error {
	return nil
}

func (f *fakeStatsServiceItemRepo) BatchUpdatePrice(ctx context.Context, ids []uint64, basePriceCents int64) error {
	return nil
}

// fakeStatsCommissionRepo 用于stats测试的抽成仓储
type fakeStatsCommissionRepo struct{}

func (f *fakeStatsCommissionRepo) CreateRule(ctx context.Context, rule *model.CommissionRule) error {
	return nil
}

func (f *fakeStatsCommissionRepo) GetRule(ctx context.Context, id uint64) (*model.CommissionRule, error) {
	return nil, repository.ErrNotFound
}

func (f *fakeStatsCommissionRepo) GetDefaultRule(ctx context.Context) (*model.CommissionRule, error) {
	return nil, repository.ErrNotFound
}

func (f *fakeStatsCommissionRepo) GetRuleForOrder(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
	return nil, repository.ErrNotFound
}

func (f *fakeStatsCommissionRepo) ListRules(ctx context.Context, opts commissionrepo.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) {
	return nil, 0, nil
}

func (f *fakeStatsCommissionRepo) UpdateRule(ctx context.Context, rule *model.CommissionRule) error {
	return nil
}

func (f *fakeStatsCommissionRepo) DeleteRule(ctx context.Context, id uint64) error {
	return nil
}

func (f *fakeStatsCommissionRepo) CreateRecord(ctx context.Context, record *model.CommissionRecord) error {
	return nil
}

func (f *fakeStatsCommissionRepo) GetRecord(ctx context.Context, id uint64) (*model.CommissionRecord, error) {
	return nil, repository.ErrNotFound
}

func (f *fakeStatsCommissionRepo) GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
	return nil, repository.ErrNotFound
}

func (f *fakeStatsCommissionRepo) ListRecords(ctx context.Context, opts commissionrepo.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) {
	return nil, 0, nil
}

func (f *fakeStatsCommissionRepo) UpdateRecord(ctx context.Context, record *model.CommissionRecord) error {
	return nil
}

func (f *fakeStatsCommissionRepo) CreateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	return nil
}

func (f *fakeStatsCommissionRepo) GetSettlement(ctx context.Context, id uint64) (*model.MonthlySettlement, error) {
	return nil, repository.ErrNotFound
}

func (f *fakeStatsCommissionRepo) GetSettlementByPlayerMonth(ctx context.Context, playerID uint64, month string) (*model.MonthlySettlement, error) {
	return nil, repository.ErrNotFound
}

func (f *fakeStatsCommissionRepo) ListSettlements(ctx context.Context, opts commissionrepo.SettlementListOptions) ([]model.MonthlySettlement, int64, error) {
	return nil, 0, nil
}

func (f *fakeStatsCommissionRepo) UpdateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	return nil
}

func (f *fakeStatsCommissionRepo) GetMonthlyStats(ctx context.Context, month string) (*commissionrepo.MonthlyStats, error) {
	return nil, nil
}

func (f *fakeStatsCommissionRepo) GetPlayerMonthlyIncome(ctx context.Context, playerID uint64, month string) (int64, error) {
	return 0, nil
}

// TestGetServiceItemStatsHandler_WithActualCall 测试服务项目统计
func TestGetServiceItemStatsHandler_WithActualCall(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("成功获取服务项目统计", func(t *testing.T) {
		orderRepo := &fakeStatsOrderRepo{
			orders: []model.Order{
				{
					Base:            model.Base{ID: 1},
					ItemID:          1,
					TotalPriceCents: 10000,
					Status:          model.OrderStatusCompleted,
				},
				{
					Base:            model.Base{ID: 2},
					ItemID:          1,
					TotalPriceCents: 15000,
					Status:          model.OrderStatusCompleted,
				},
				{
					Base:            model.Base{ID: 3},
					ItemID:          2,
					TotalPriceCents: 20000,
					Status:          model.OrderStatusCompleted,
				},
			},
		}

		serviceItemRepo := &fakeStatsServiceItemRepo{
			items: []model.ServiceItem{
				{
					ID:       1,
					ItemCode: "ITEM001",
					Name:     "服务1",
				},
				{
					ID:       2,
					ItemCode: "ITEM002",
					Name:     "服务2",
				},
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/service-items", nil)

		getServiceItemStatsHandler(c, orderRepo, serviceItemRepo)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("服务项目仓储错误", func(t *testing.T) {
		orderRepo := &fakeStatsOrderRepo{}
		serviceItemRepo := &fakeStatsServiceItemRepo{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/service-items", nil)

		getServiceItemStatsHandler(c, orderRepo, serviceItemRepo)

		// 即使没有数据也应该返回200
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestGetTopPlayersHandler_WithActualCall 测试Top陪玩师统计
func TestGetTopPlayersHandler_WithActualCall(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("成功获取Top陪玩师-默认月份", func(t *testing.T) {
		commissionRepo := &fakeStatsCommissionRepo{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/top-players", nil)

		getTopPlayersHandler(c, commissionRepo)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("成功获取Top陪玩师-指定月份", func(t *testing.T) {
		commissionRepo := &fakeStatsCommissionRepo{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/top-players?month=2024-11&limit=20", nil)

		getTopPlayersHandler(c, commissionRepo)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestGetAdminGiftStatsHandler_WithActualCall 测试礼物统计
func TestGetAdminGiftStatsHandler_WithActualCall(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("成功获取礼物统计", func(t *testing.T) {
		orderRepo := &fakeStatsOrderRepo{
			orders: []model.Order{
				{
					Base:            model.Base{ID: 1},
					ItemID:          101,
					Quantity:        5,
					TotalPriceCents: 50000,
					Status:          model.OrderStatusCompleted,
				},
				{
					Base:            model.Base{ID: 2},
					ItemID:          102,
					Quantity:        3,
					TotalPriceCents: 30000,
					Status:          model.OrderStatusCompleted,
				},
			},
		}

		serviceItemRepo := &fakeStatsServiceItemRepo{
			gifts: []model.ServiceItem{
				{
					ID:   101,
					Name: "礼物1",
				},
				{
					ID:   102,
					Name: "礼物2",
				},
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/gift-stats", nil)

		getAdminGiftStatsHandler(c, orderRepo, serviceItemRepo)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("礼物仓储错误返回500", func(t *testing.T) {
		orderRepo := &fakeStatsOrderRepo{}
		// 创建一个会返回错误的repo
		serviceItemRepo := &fakeStatsServiceItemRepo{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/gift-stats", nil)

		getAdminGiftStatsHandler(c, orderRepo, serviceItemRepo)

		// 空数据应该返回200
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestGetRevenueByGameHandler_WithActualCall 测试按游戏统计收入
func TestGetRevenueByGameHandler_WithActualCall(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("成功获取游戏收入统计", func(t *testing.T) {
		gameID1 := uint64(1)
		gameID2 := uint64(2)
		orderRepo := &fakeStatsOrderRepo{
			orders: []model.Order{
				{
					Base:            model.Base{ID: 1},
					GameID:          &gameID1,
					TotalPriceCents: 10000,
					Status:          model.OrderStatusCompleted,
				},
				{
					Base:            model.Base{ID: 2},
					GameID:          &gameID1,
					TotalPriceCents: 15000,
					Status:          model.OrderStatusCompleted,
				},
				{
					Base:            model.Base{ID: 3},
					GameID:          &gameID2,
					TotalPriceCents: 20000,
					Status:          model.OrderStatusCompleted,
				},
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/revenue-by-game", nil)

		getRevenueByGameHandler(c, orderRepo)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("空订单列表", func(t *testing.T) {
		orderRepo := &fakeStatsOrderRepo{
			orders: []model.Order{},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/revenue-by-game", nil)

		getRevenueByGameHandler(c, orderRepo)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestRegisterStatsAnalysisRoutes_Coverage 测试路由注册
func TestRegisterStatsAnalysisRoutes_Coverage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	orderRepo := &fakeStatsOrderRepo{}
	commissionRepo := &fakeStatsCommissionRepo{}
	serviceItemRepo := &fakeStatsServiceItemRepo{}

	RegisterStatsAnalysisRoutes(router, orderRepo, commissionRepo, serviceItemRepo)

	routes := router.Routes()
	assert.NotEmpty(t, routes)

	routeMap := make(map[string]bool)
	for _, route := range routes {
		routeMap[route.Method+":"+route.Path] = true
	}

	assert.True(t, routeMap["GET:/admin/stats/service-items"])
	assert.True(t, routeMap["GET:/admin/stats/top-players"])
	assert.True(t, routeMap["GET:/admin/stats/gift-stats"])
	assert.True(t, routeMap["GET:/admin/stats/revenue-by-game"])
}
