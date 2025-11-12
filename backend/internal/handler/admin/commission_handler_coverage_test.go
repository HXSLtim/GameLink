package admin

import (
    "bytes"
    "context"
    "errors"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionrepo "gamelink/internal/repository/commission"
	"gamelink/internal/service/commission"
)

// fakeCommissionRepo 实现CommissionRepository接口
type fakeCommissionRepo struct {
    rules       []model.CommissionRule
    records     []model.CommissionRecord
    settlements []model.MonthlySettlement
    createError error
    getError    error
    statsErr    error
}

func (f *fakeCommissionRepo) CreateRule(ctx context.Context, rule *model.CommissionRule) error {
	if f.createError != nil {
		return f.createError
	}
	if rule.ID == 0 {
		rule.ID = uint64(len(f.rules) + 1)
	}
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()
	f.rules = append(f.rules, *rule)
	return nil
}

func (f *fakeCommissionRepo) GetRule(ctx context.Context, id uint64) (*model.CommissionRule, error) {
	if f.getError != nil {
		return nil, f.getError
	}
	for i := range f.rules {
		if f.rules[i].ID == id {
			return &f.rules[i], nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeCommissionRepo) GetDefaultRule(ctx context.Context) (*model.CommissionRule, error) {
	for i := range f.rules {
		if f.rules[i].Type == "default" {
			return &f.rules[i], nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeCommissionRepo) GetRuleForOrder(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
	return f.GetDefaultRule(ctx)
}

func (f *fakeCommissionRepo) ListRules(ctx context.Context, opts commissionrepo.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) {
	return f.rules, int64(len(f.rules)), nil
}

func (f *fakeCommissionRepo) UpdateRule(ctx context.Context, rule *model.CommissionRule) error {
	for i := range f.rules {
		if f.rules[i].ID == rule.ID {
			f.rules[i] = *rule
			return nil
		}
	}
	return repository.ErrNotFound
}

func (f *fakeCommissionRepo) DeleteRule(ctx context.Context, id uint64) error {
	for i := range f.rules {
		if f.rules[i].ID == id {
			f.rules = append(f.rules[:i], f.rules[i+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}

func (f *fakeCommissionRepo) CreateRecord(ctx context.Context, record *model.CommissionRecord) error {
	if record.ID == 0 {
		record.ID = uint64(len(f.records) + 1)
	}
	f.records = append(f.records, *record)
	return nil
}

func (f *fakeCommissionRepo) GetRecord(ctx context.Context, id uint64) (*model.CommissionRecord, error) {
	for i := range f.records {
		if f.records[i].ID == id {
			return &f.records[i], nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeCommissionRepo) GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
	for i := range f.records {
		if f.records[i].OrderID == orderID {
			return &f.records[i], nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeCommissionRepo) ListRecords(ctx context.Context, opts commissionrepo.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) {
	return f.records, int64(len(f.records)), nil
}

func (f *fakeCommissionRepo) UpdateRecord(ctx context.Context, record *model.CommissionRecord) error {
	for i := range f.records {
		if f.records[i].ID == record.ID {
			f.records[i] = *record
			return nil
		}
	}
	return repository.ErrNotFound
}

func (f *fakeCommissionRepo) CreateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	if settlement.ID == 0 {
		settlement.ID = uint64(len(f.settlements) + 1)
	}
	f.settlements = append(f.settlements, *settlement)
	return nil
}

func (f *fakeCommissionRepo) GetSettlement(ctx context.Context, id uint64) (*model.MonthlySettlement, error) {
	for i := range f.settlements {
		if f.settlements[i].ID == id {
			return &f.settlements[i], nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeCommissionRepo) GetSettlementByPlayerMonth(ctx context.Context, playerID uint64, month string) (*model.MonthlySettlement, error) {
	for i := range f.settlements {
		if f.settlements[i].PlayerID == playerID && f.settlements[i].SettlementMonth == month {
			return &f.settlements[i], nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeCommissionRepo) ListSettlements(ctx context.Context, opts commissionrepo.SettlementListOptions) ([]model.MonthlySettlement, int64, error) {
	return f.settlements, int64(len(f.settlements)), nil
}

func (f *fakeCommissionRepo) UpdateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	for i := range f.settlements {
		if f.settlements[i].ID == settlement.ID {
			f.settlements[i] = *settlement
			return nil
		}
	}
	return repository.ErrNotFound
}

func (f *fakeCommissionRepo) GetMonthlyStats(ctx context.Context, month string) (*commissionrepo.MonthlyStats, error) {
    if f.statsErr != nil {
        return nil, f.statsErr
    }
    var totalOrders, totalCommission, totalPlayerIncome int64
    recordCount := 0
	
	for _, record := range f.records {
		if record.SettlementMonth == month {
			totalOrders += record.TotalAmountCents
			totalCommission += record.CommissionCents
			totalPlayerIncome += record.PlayerIncomeCents
			recordCount++
		}
	}
	
	return &commissionrepo.MonthlyStats{
		TotalOrders:       int64(recordCount),
		TotalIncome:       totalOrders,
		TotalCommission:   totalCommission,
		TotalPlayerIncome: totalPlayerIncome,
	}, nil
}

func (f *fakeCommissionRepo) GetPlayerMonthlyIncome(ctx context.Context, playerID uint64, month string) (int64, error) {
	var total int64
	for _, record := range f.records {
		if record.PlayerID == playerID && record.SettlementMonth == month {
			total += record.PlayerIncomeCents
		}
	}
	return total, nil
}

// fakeCommissionOrderRepo 简化的订单仓储
type fakeCommissionOrderRepo struct {
	orders []model.Order
}

func (f *fakeCommissionOrderRepo) Get(ctx context.Context, id uint64) (*model.Order, error) {
	for i := range f.orders {
		if f.orders[i].ID == id {
			return &f.orders[i], nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeCommissionOrderRepo) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	return f.orders, int64(len(f.orders)), nil
}

func (f *fakeCommissionOrderRepo) Create(ctx context.Context, order *model.Order) error {
	if order.ID == 0 {
		order.ID = uint64(len(f.orders) + 1)
	}
	f.orders = append(f.orders, *order)
	return nil
}

func (f *fakeCommissionOrderRepo) Update(ctx context.Context, order *model.Order) error {
	for i := range f.orders {
		if f.orders[i].ID == order.ID {
			f.orders[i] = *order
			return nil
		}
	}
	return repository.ErrNotFound
}

func (f *fakeCommissionOrderRepo) Delete(ctx context.Context, id uint64) error {
	return nil
}

// fakeCommissionPlayerRepo 简化的玩家仓储
type fakeCommissionPlayerRepo struct {
	players []model.Player
}

func (f *fakeCommissionPlayerRepo) Get(ctx context.Context, id uint64) (*model.Player, error) {
	for i := range f.players {
		if f.players[i].ID == id {
			return &f.players[i], nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeCommissionPlayerRepo) List(ctx context.Context) ([]model.Player, error) {
	return f.players, nil
}

func (f *fakeCommissionPlayerRepo) ListPaged(ctx context.Context, page, size int) ([]model.Player, int64, error) {
	return f.players, int64(len(f.players)), nil
}

func (f *fakeCommissionPlayerRepo) Create(ctx context.Context, player *model.Player) error {
	if player.ID == 0 {
		player.ID = uint64(len(f.players) + 1)
	}
	f.players = append(f.players, *player)
	return nil
}

func (f *fakeCommissionPlayerRepo) Update(ctx context.Context, player *model.Player) error {
	for i := range f.players {
		if f.players[i].ID == player.ID {
			f.players[i] = *player
			return nil
		}
	}
	return repository.ErrNotFound
}

func (f *fakeCommissionPlayerRepo) Delete(ctx context.Context, id uint64) error {
	return nil
}

// fakeScheduler 模拟调度器
type fakeScheduler struct {
    triggerError error
    lastMonth    string
}

func (f *fakeScheduler) TriggerSettlement(month string) error {
	f.lastMonth = month
	return f.triggerError
}

func TestTriggerSettlement_DefaultAndError(t *testing.T) {
    gin.SetMode(gin.TestMode)
    _ = commission.NewCommissionService(&fakeCommissionRepo{}, &fakeCommissionOrderRepo{}, &fakeCommissionPlayerRepo{})

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    sch := &fakeScheduler{}
    c.Request = httptest.NewRequest(http.MethodPost, "/admin/commission/settlements/trigger", nil)
    triggerSettlementHandler(c, sch)
    assert.Equal(t, http.StatusOK, w.Code)
    if sch.lastMonth == "" { t.Fatal("expected non-empty month") }

    w2 := httptest.NewRecorder()
    c2, _ := gin.CreateTestContext(w2)
    sch2 := &fakeScheduler{triggerError: errors.New("boom")}
    c2.Request = httptest.NewRequest(http.MethodPost, "/admin/commission/settlements/trigger?month=2025-01", nil)
    triggerSettlementHandler(c2, sch2)
    assert.Equal(t, http.StatusInternalServerError, w2.Code)
}

func TestGetPlatformStats_ErrorPath(t *testing.T) {
    gin.SetMode(gin.TestMode)
    commRepo := &fakeCommissionRepo{statsErr: errors.New("db")}
    orderRepo := &fakeCommissionOrderRepo{}
    playerRepo := &fakeCommissionPlayerRepo{}
    svc := commission.NewCommissionService(commRepo, orderRepo, playerRepo)

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest(http.MethodGet, "/admin/commission/stats?month=2025-01", nil)
    getPlatformStatsHandler(c, svc)
    assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestCreateCommissionRuleHandler_WithActualCall 测试实际调用createCommissionRuleHandler
func TestCreateCommissionRuleHandler_WithActualCall(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("成功创建默认规则", func(t *testing.T) {
		commRepo := &fakeCommissionRepo{}
		orderRepo := &fakeCommissionOrderRepo{}
		playerRepo := &fakeCommissionPlayerRepo{}
		svc := commission.NewCommissionService(commRepo, orderRepo, playerRepo)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := commission.CreateCommissionRuleRequest{
			Name:        "默认抽成规则",
			Description: "平台默认20%抽成",
			Type:        "default",
			Rate:        20,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPost, "/admin/commission/rules", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		createCommissionRuleHandler(c, svc)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.APIResponse[model.CommissionRule]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "默认抽成规则", response.Data.Name)
		assert.Equal(t, 20, response.Data.Rate)
		assert.Equal(t, uint64(1), response.Data.ID)
	})

	t.Run("无效JSON返回400", func(t *testing.T) {
		commRepo := &fakeCommissionRepo{}
		orderRepo := &fakeCommissionOrderRepo{}
		playerRepo := &fakeCommissionPlayerRepo{}
		svc := commission.NewCommissionService(commRepo, orderRepo, playerRepo)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/admin/commission/rules", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		createCommissionRuleHandler(c, svc)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("rate超出范围返回400", func(t *testing.T) {
		commRepo := &fakeCommissionRepo{}
		orderRepo := &fakeCommissionOrderRepo{}
		playerRepo := &fakeCommissionPlayerRepo{}
		svc := commission.NewCommissionService(commRepo, orderRepo, playerRepo)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := commission.CreateCommissionRuleRequest{
			Name: "无效规则",
			Type: "default",
			Rate: 150, // 超过100
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPost, "/admin/commission/rules", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		createCommissionRuleHandler(c, svc)

		// rate验证失败会在binding阶段被拦截，返回400
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestUpdateCommissionRuleHandler_WithActualCall 测试实际调用updateCommissionRuleHandler
func TestUpdateCommissionRuleHandler_WithActualCall(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("成功更新规则", func(t *testing.T) {
		commRepo := &fakeCommissionRepo{
			rules: []model.CommissionRule{
				{
					ID:   1,
					Name: "旧规则",
					Type: "default",
					Rate: 20,
				},
			},
		}
		orderRepo := &fakeCommissionOrderRepo{}
		playerRepo := &fakeCommissionPlayerRepo{}
		svc := commission.NewCommissionService(commRepo, orderRepo, playerRepo)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		newName := "更新后的规则"
		newRate := 25
		reqBody := commission.UpdateCommissionRuleRequest{
			Name: &newName,
			Rate: &newRate,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPut, "/admin/commission/rules/1", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		updateCommissionRuleHandler(c, svc)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.APIResponse[any]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("无效ID返回400", func(t *testing.T) {
		commRepo := &fakeCommissionRepo{}
		orderRepo := &fakeCommissionOrderRepo{}
		playerRepo := &fakeCommissionPlayerRepo{}
		svc := commission.NewCommissionService(commRepo, orderRepo, playerRepo)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		name := "测试"
		reqBody := commission.UpdateCommissionRuleRequest{Name: &name}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPut, "/admin/commission/rules/invalid", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		updateCommissionRuleHandler(c, svc)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("规则不存在返回500", func(t *testing.T) {
		commRepo := &fakeCommissionRepo{rules: []model.CommissionRule{}}
		orderRepo := &fakeCommissionOrderRepo{}
		playerRepo := &fakeCommissionPlayerRepo{}
		svc := commission.NewCommissionService(commRepo, orderRepo, playerRepo)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		name := "测试"
		reqBody := commission.UpdateCommissionRuleRequest{Name: &name}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPut, "/admin/commission/rules/999", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "999"}}

		updateCommissionRuleHandler(c, svc)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

// TestTriggerSettlementHandler_WithActualCall 测试实际调用triggerSettlementHandler
func TestTriggerSettlementHandler_WithActualCall(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("成功触发结算-指定月份", func(t *testing.T) {
		scheduler := &fakeScheduler{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/admin/commission/settlements/trigger?month=2024-11", nil)

		triggerSettlementHandler(c, scheduler)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "2024-11", scheduler.lastMonth)
		
		var response model.APIResponse[any]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Contains(t, response.Message, "2024-11")
	})

	t.Run("成功触发结算-默认上月", func(t *testing.T) {
		scheduler := &fakeScheduler{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/admin/commission/settlements/trigger", nil)

		triggerSettlementHandler(c, scheduler)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, scheduler.lastMonth)
	})

	t.Run("调度器错误返回500", func(t *testing.T) {
		scheduler := &fakeScheduler{
			triggerError: assert.AnError,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/admin/commission/settlements/trigger?month=2024-11", nil)

		triggerSettlementHandler(c, scheduler)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

// TestGetPlatformStatsHandler_WithActualCall 测试实际调用getPlatformStatsHandler
func TestGetPlatformStatsHandler_WithActualCall(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("成功获取统计-指定月份", func(t *testing.T) {
		commRepo := &fakeCommissionRepo{
			records: []model.CommissionRecord{
				{
					ID:                 1,
					OrderID:            1,
					PlayerID:           1,
					TotalAmountCents:   10000,
					CommissionCents:    2000,
					PlayerIncomeCents:  8000,
					SettlementMonth:    "2024-11",
				},
				{
					ID:                 2,
					OrderID:            2,
					PlayerID:           1,
					TotalAmountCents:   20000,
					CommissionCents:    4000,
					PlayerIncomeCents:  16000,
					SettlementMonth:    "2024-11",
				},
			},
		}
		orderRepo := &fakeCommissionOrderRepo{}
		playerRepo := &fakeCommissionPlayerRepo{}
		svc := commission.NewCommissionService(commRepo, orderRepo, playerRepo)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/admin/commission/stats?month=2024-11", nil)

		getPlatformStatsHandler(c, svc)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response model.APIResponse[commission.PlatformStatsResponse]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "2024-11", response.Data.Month)
		assert.Equal(t, int64(30000), response.Data.TotalIncome)
		assert.Equal(t, int64(6000), response.Data.TotalCommission)
		assert.Equal(t, int64(24000), response.Data.TotalPlayerIncome)
	})

	t.Run("成功获取统计-默认当前月", func(t *testing.T) {
		commRepo := &fakeCommissionRepo{records: []model.CommissionRecord{}}
		orderRepo := &fakeCommissionOrderRepo{}
		playerRepo := &fakeCommissionPlayerRepo{}
		svc := commission.NewCommissionService(commRepo, orderRepo, playerRepo)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/admin/commission/stats", nil)

		getPlatformStatsHandler(c, svc)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response model.APIResponse[commission.PlatformStatsResponse]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotEmpty(t, response.Data.Month)
	})
}

// TestRegisterCommissionRoutes_Coverage 测试路由注册覆盖率
func TestRegisterCommissionRoutes_Coverage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	commRepo := &fakeCommissionRepo{}
	orderRepo := &fakeCommissionOrderRepo{}
	playerRepo := &fakeCommissionPlayerRepo{}
	svc := commission.NewCommissionService(commRepo, orderRepo, playerRepo)
	scheduler := &fakeScheduler{}

	RegisterCommissionRoutes(router, svc, scheduler)

	routes := router.Routes()
	assert.NotEmpty(t, routes)

	routeMap := make(map[string]bool)
	for _, route := range routes {
		routeMap[route.Method+":"+route.Path] = true
	}

	assert.True(t, routeMap["POST:/admin/commission/rules"])
	assert.True(t, routeMap["PUT:/admin/commission/rules/:id"])
	assert.True(t, routeMap["POST:/admin/commission/settlements/trigger"])
	assert.True(t, routeMap["GET:/admin/commission/stats"])
}
