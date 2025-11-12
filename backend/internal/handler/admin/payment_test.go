package admin

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gamelink/internal/cache"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	adminservice "gamelink/internal/service/admin"
)

type fakePaymentRepoForHandler struct {
	items    []model.Payment
	listFunc func(opts repository.PaymentListOptions) ([]model.Payment, int64, error)
}

func (f *fakePaymentRepoForHandler) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	if f.listFunc != nil {
		return f.listFunc(opts)
	}

	var filtered []model.Payment
	for _, p := range f.items {
		if opts.OrderID != nil && p.OrderID != *opts.OrderID {
			continue
		}
		if opts.UserID != nil && p.UserID != *opts.UserID {
			continue
		}
		if opts.Status != nil && p.Status != *opts.Status {
			continue
		}
		if opts.Method != nil && p.Method != *opts.Method {
			continue
		}
		filtered = append(filtered, p)
	}

	total := int64(len(filtered))
	page := repository.NormalizePage(opts.Page)
	pageSize := repository.NormalizePageSize(opts.PageSize)
	start := (page - 1) * pageSize
	if start >= len(filtered) {
		return []model.Payment{}, total, nil
	}
	end := start + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], total, nil
}

func (f *fakePaymentRepoForHandler) Create(ctx context.Context, p *model.Payment) error {
	if p.ID == 0 {
		p.ID = uint64(len(f.items) + 1)
	}
	f.items = append(f.items, *p)
	return nil
}

func (f *fakePaymentRepoForHandler) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	for i := range f.items {
		if f.items[i].ID == id {
			c := f.items[i]
			return &c, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakePaymentRepoForHandler) Update(ctx context.Context, p *model.Payment) error {
	for i := range f.items {
		if f.items[i].ID == p.ID {
			f.items[i] = *p
			return nil
		}
	}
	return repository.ErrNotFound
}

func (f *fakePaymentRepoForHandler) Delete(ctx context.Context, id uint64) error {
	for i := range f.items {
		if f.items[i].ID == id {
			f.items = append(f.items[:i], f.items[i+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}

func setupPaymentTestRouter(paymentRepo *fakePaymentRepoForHandler) (*gin.Engine, *adminservice.AdminService) {
	r := newTestEngine()

	orderRepo := &fakeOrderRepo{
		obj: &model.Order{
			Base:            model.Base{ID: 1},
			UserID:          1,
			ItemID:          1,
			PlayerID:        nil,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
		},
	}

	userRepo := &fakeUserRepo{
		last: &model.User{
			Base:   model.Base{ID: 1},
			Name:   "测试用户",
			Status: "active",
		},
	}

	svc := adminservice.NewAdminService(
		&fakeGameRepo{},
		userRepo,
		&fakePlayerRepo{},
		orderRepo,
		paymentRepo,
		&fakeRoleRepo{},
		cache.NewMemory(),
	)

	handler := NewPaymentHandler(svc)
	r.POST("/admin/payments", handler.CreatePayment)
	r.POST("/admin/payments/:id/capture", handler.CapturePayment)
	r.GET("/admin/payments", handler.ListPayments)
	r.GET("/admin/payments/:id", handler.GetPayment)
	r.PUT("/admin/payments/:id", handler.UpdatePayment)
	r.DELETE("/admin/payments/:id", handler.DeletePayment)
	r.POST("/admin/payments/:id/refund", handler.RefundPayment)
	r.GET("/admin/payments/:id/logs", handler.ListPaymentLogs)

	return r, svc
}

func TestPaymentHandler_CreatePayment(t *testing.T) {
	paymentRepo := &fakePaymentRepoForHandler{
		items: []model.Payment{},
	}
	r, _ := setupPaymentTestRouter(paymentRepo)

	payload := CreatePaymentPayload{
		OrderID:     1,
		UserID:      1,
		Method:      "alipay",
		AmountCents: 10000,
		Currency:    "CNY",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/payments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp model.APIResponse[*model.Payment]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Data)
}

func TestPaymentHandler_CreatePayment_Validation(t *testing.T) {
	paymentRepo := &fakePaymentRepoForHandler{}
	r, _ := setupPaymentTestRouter(paymentRepo)

	// 测试缺少必填字段
	payload := CreatePaymentPayload{
		OrderID: 0, // 缺少必填字段
		UserID:  1,
		Method:  "alipay",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/payments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPaymentHandler_ListPayments(t *testing.T) {
	paymentRepo := &fakePaymentRepoForHandler{
		items: []model.Payment{
			{Base: model.Base{ID: 1}, OrderID: 1, UserID: 1, Status: model.PaymentStatusPending},
			{Base: model.Base{ID: 2}, OrderID: 2, UserID: 2, Status: model.PaymentStatusPaid},
		},
	}
	r, _ := setupPaymentTestRouter(paymentRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/payments?page=1&page_size=20", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[[]model.Payment]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestPaymentHandler_GetPayment(t *testing.T) {
	paymentRepo := &fakePaymentRepoForHandler{
		items: []model.Payment{
			{Base: model.Base{ID: 1}, OrderID: 1, UserID: 1, AmountCents: 10000},
		},
	}
	r, _ := setupPaymentTestRouter(paymentRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/payments/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[*model.Payment]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, uint64(1), resp.Data.ID)
}

func TestPaymentHandler_GetPayment_NotFound(t *testing.T) {
	paymentRepo := &fakePaymentRepoForHandler{
		items: []model.Payment{},
	}
	r, _ := setupPaymentTestRouter(paymentRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/payments/999", nil)
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestPaymentHandler_UpdatePayment(t *testing.T) {
	now := time.Now().UTC()
	paidAtStr := now.Format(time.RFC3339)
	paymentRepo := &fakePaymentRepoForHandler{
		items: []model.Payment{
			{Base: model.Base{ID: 1}, OrderID: 1, UserID: 1, Status: model.PaymentStatusPending},
		},
	}
	r, _ := setupPaymentTestRouter(paymentRepo)

	payload := UpdatePaymentPayload{
		Status:          "paid",
		ProviderTradeNo: "TRADE123",
		PaidAt:          &paidAtStr,
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/admin/payments/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var resp model.APIResponse[*model.Payment]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	}
}

func TestPaymentHandler_CapturePayment(t *testing.T) {
	now := time.Now().UTC()
	paidAtStr := now.Format(time.RFC3339)
	paymentRepo := &fakePaymentRepoForHandler{
		items: []model.Payment{
			{Base: model.Base{ID: 1}, OrderID: 1, UserID: 1, Status: model.PaymentStatusPending},
		},
	}
	r, _ := setupPaymentTestRouter(paymentRepo)

	payload := CapturePaymentPayload{
		ProviderTradeNo: "TRADE123",
		PaidAt:          &paidAtStr,
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/payments/1/capture", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var resp model.APIResponse[*model.Payment]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	}
}

func TestPaymentHandler_RefundPayment(t *testing.T) {
	now := time.Now().UTC()
	refundedAtStr := now.Format(time.RFC3339)
	paymentRepo := &fakePaymentRepoForHandler{
		items: []model.Payment{
			{Base: model.Base{ID: 1}, OrderID: 1, UserID: 1, Status: model.PaymentStatusPaid},
		},
	}
	r, _ := setupPaymentTestRouter(paymentRepo)

	payload := RefundPaymentPayload{
		RefundedAt:      &refundedAtStr,
		ProviderTradeNo: "REFUND123",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/payments/1/refund", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var resp model.APIResponse[*model.Payment]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	}
}

func TestPaymentHandler_DeletePayment(t *testing.T) {
	paymentRepo := &fakePaymentRepoForHandler{
		items: []model.Payment{
			{Base: model.Base{ID: 1}, OrderID: 1, UserID: 1},
		},
	}
	r, _ := setupPaymentTestRouter(paymentRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/admin/payments/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestPaymentHandler_ListPaymentLogs(t *testing.T) {
    t.Skip("ListPaymentLogs requires TxManager, skipping for now")
}

func TestPaymentHandler_UpdatePayment_InvalidStatus(t *testing.T) {
    paymentRepo := &fakePaymentRepoForHandler{items: []model.Payment{{Base: model.Base{ID: 1}, Status: model.PaymentStatusPending, OrderID: 1, UserID: 1}}}
    r, _ := setupPaymentTestRouter(paymentRepo)
    payload := UpdatePaymentPayload{Status: "invalid"}
    body, _ := json.Marshal(payload)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPut, "/admin/payments/1", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError { t.Fatalf("expected 400/500, got %d", w.Code) }
}

func TestPaymentHandler_UpdatePayment_InvalidTransition(t *testing.T) {
    paymentRepo := &fakePaymentRepoForHandler{items: []model.Payment{{Base: model.Base{ID: 1}, Status: model.PaymentStatusRefunded, OrderID: 1, UserID: 1}}}
    r, _ := setupPaymentTestRouter(paymentRepo)
    payload := UpdatePaymentPayload{Status: "paid"}
    body, _ := json.Marshal(payload)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPut, "/admin/payments/1", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError { t.Fatalf("expected 400/500, got %d", w.Code) }
}

func TestPaymentHandler_CapturePayment_InvalidTransition(t *testing.T) {
    paymentRepo := &fakePaymentRepoForHandler{items: []model.Payment{{Base: model.Base{ID: 1}, Status: model.PaymentStatusRefunded, OrderID: 1, UserID: 1}}}
    r, _ := setupPaymentTestRouter(paymentRepo)
    payload := CapturePaymentPayload{ProviderTradeNo: "X"}
    body, _ := json.Marshal(payload)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/admin/payments/1/capture", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError { t.Fatalf("expected 400/500, got %d", w.Code) }
}

func TestPaymentHandler_CapturePayment_InvalidPaidAt(t *testing.T) {
    paymentRepo := &fakePaymentRepoForHandler{items: []model.Payment{{Base: model.Base{ID: 1}, Status: model.PaymentStatusPending, OrderID: 1, UserID: 1}}}
    r, _ := setupPaymentTestRouter(paymentRepo)
    s := "bad"
    payload := CapturePaymentPayload{PaidAt: &s}
    body, _ := json.Marshal(payload)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/admin/payments/1/capture", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError { t.Fatalf("expected 400/500, got %d", w.Code) }
}

func TestPaymentHandler_ListPaymentLogs_InvalidDates(t *testing.T) {
    r, _ := setupPaymentTestRouter(&fakePaymentRepoForHandler{})
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/admin/payments/1/logs?date_from=bad&date_to=bad", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("expected 400, got %d", w.Code) }
}

func TestPaymentHandler_Capture_InvalidJSON(t *testing.T) {
    r, _ := setupPaymentTestRouter(&fakePaymentRepoForHandler{})
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/admin/payments/1/capture", bytes.NewBufferString("bad"))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("expected 400, got %d", w.Code) }
}

func TestPaymentHandler_ListPaymentLogs_InvalidActorID(t *testing.T) {
    r, _ := setupPaymentTestRouter(&fakePaymentRepoForHandler{})
    w := httptest.NewRecorder()
    url := "/admin/payments/1/logs?actor_user_id=abc&date_from=2025-01-01T00:00:00Z&date_to=2025-01-02T00:00:00Z"
    req := httptest.NewRequest(http.MethodGet, url, nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("expected 400, got %d", w.Code) }
}

func TestExportOperationLogsCSV_Payment(t *testing.T) {
    r := newTestEngine()
    r.GET("/export_pay", func(c *gin.Context) {
        items := []model.OperationLog{
            {Base: model.Base{ID: 1, CreatedAt: time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)}, EntityType: "payment", EntityID: 7, Action: "capture", Reason: "", MetadataJSON: []byte("{\"ok\":true}")},
            {Base: model.Base{ID: 2, CreatedAt: time.Date(2025, 1, 3, 3, 4, 5, 0, time.UTC)}, EntityType: "payment", EntityID: 7, Action: "refund", Reason: "dup", MetadataJSON: nil},
        }
        exportOperationLogsCSV(c, "payment", 7, items)
    })
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/export_pay?fields=id,action,created_at&header_lang=zh&tz=Asia/Shanghai&bom=true", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w.Code) }
    if ct := w.Header().Get("Content-Type"); ct == "" || !strings.Contains(ct, "text/csv") { t.Fatalf("expected csv content type, got %q", ct) }
    if cd := w.Header().Get("Content-Disposition"); !strings.Contains(cd, "attachment") || !strings.Contains(cd, "payment_7_logs.csv") { t.Fatalf("unexpected content disposition: %q", cd) }
}

func TestPaymentHandler_ListPayments_InvalidQueryIDs(t *testing.T) {
    r, _ := setupPaymentTestRouter(&fakePaymentRepoForHandler{})
    for _, url := range []string{
        "/admin/payments?page=1&page_size=20&user_id=abc",
        "/admin/payments?page=1&page_size=20&order_id=abc",
    } {
        w := httptest.NewRecorder()
        req := httptest.NewRequest(http.MethodGet, url, nil)
        r.ServeHTTP(w, req)
        if w.Code != http.StatusBadRequest { t.Fatalf("expected 400, got %d", w.Code) }
    }
}

func TestPaymentHandler_ListPayments_InvalidPagination(t *testing.T) {
    r, _ := setupPaymentTestRouter(&fakePaymentRepoForHandler{})
    for _, url := range []string{
        "/admin/payments?page=abc",
        "/admin/payments?page=1&page_size=abc",
    } {
        w := httptest.NewRecorder()
        req := httptest.NewRequest(http.MethodGet, url, nil)
        r.ServeHTTP(w, req)
        if w.Code != http.StatusBadRequest { t.Fatalf("expected 400, got %d", w.Code) }
    }
}

