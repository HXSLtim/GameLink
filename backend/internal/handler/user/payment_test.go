package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service/payment"
)

// ---- Fake OrderRepository for payment tests ----

type fakeOrderRepositoryForPayment struct {
	orders map[uint64]*model.Order
}

func newFakeOrderRepositoryForPayment() *fakeOrderRepositoryForPayment {
	gameID := uint64(1)
	return &fakeOrderRepositoryForPayment{
		orders: map[uint64]*model.Order{
			10: {Base: model.Base{ID: 10}, UserID: 100, GameID: &gameID, ItemID: 1, OrderNo: "PAYMENT-TEST-01", Status: model.OrderStatusPending, TotalPriceCents: 5000, UnitPriceCents: 5000},
			11: {Base: model.Base{ID: 11}, UserID: 101, GameID: &gameID, ItemID: 1, OrderNo: "PAYMENT-TEST-02", Status: model.OrderStatusPending, TotalPriceCents: 8000, UnitPriceCents: 8000},
		},
	}
}

func (m *fakeOrderRepositoryForPayment) Create(ctx context.Context, o *model.Order) error {
	o.ID = uint64(len(m.orders) + 1)
	m.orders[o.ID] = o
	return nil
}

func (m *fakeOrderRepositoryForPayment) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	var result []model.Order
	for _, o := range m.orders {
		result = append(result, *o)
	}
	return result, int64(len(result)), nil
}

func (m *fakeOrderRepositoryForPayment) Get(ctx context.Context, id uint64) (*model.Order, error) {
	if o, ok := m.orders[id]; ok {
		return o, nil
	}
	return nil, repository.ErrNotFound
}

func (m *fakeOrderRepositoryForPayment) Update(ctx context.Context, o *model.Order) error {
	m.orders[o.ID] = o
	return nil
}

func (m *fakeOrderRepositoryForPayment) Delete(ctx context.Context, id uint64) error {
	delete(m.orders, id)
	return nil
}

// ---- Fake PaymentRepository for user_payment tests ----

type mockPaymentRepoForUserPayment struct {
	payments map[uint64]*model.Payment
}

func newMockPaymentRepoForUserPayment() *mockPaymentRepoForUserPayment {
	return &mockPaymentRepoForUserPayment{
		payments: map[uint64]*model.Payment{
			1: {Base: model.Base{ID: 1}, OrderID: 10, UserID: 100, AmountCents: 5000, Status: model.PaymentStatusPending},
			2: {Base: model.Base{ID: 2}, OrderID: 11, UserID: 101, AmountCents: 8000, Status: model.PaymentStatusPaid},
		},
	}
}

func (m *mockPaymentRepoForUserPayment) Create(ctx context.Context, payment *model.Payment) error {
	payment.ID = uint64(len(m.payments) + 1)
	m.payments[payment.ID] = payment
	return nil
}

func (m *mockPaymentRepoForUserPayment) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	var filtered []model.Payment
	for _, p := range m.payments {
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
		filtered = append(filtered, *p)
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

func (m *mockPaymentRepoForUserPayment) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	if p, ok := m.payments[id]; ok {
		return p, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockPaymentRepoForUserPayment) Update(ctx context.Context, payment *model.Payment) error {
	m.payments[payment.ID] = payment
	return nil
}

func (m *mockPaymentRepoForUserPayment) Delete(ctx context.Context, id uint64) error {
	delete(m.payments, id)
	return nil
}

// ---- Tests for user_payment.go ----

func TestCreatePaymentHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	paymentRepo := newMockPaymentRepoForUserPayment()
	paymentSvc := payment.NewPaymentService(paymentRepo, newFakeOrderRepositoryForPayment())

	router := gin.New()
	router.POST("/user/payments", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		createPaymentHandler(c, paymentSvc)
	})

	reqBody := payment.CreatePaymentRequest{
		OrderID: 10,
		Method:  "alipay",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/user/payments", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp model.APIResponse[payment.CreatePaymentResponse]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("Expected success=true, got %v", resp.Success)
	}
}

func TestCreatePaymentHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	paymentRepo := newMockPaymentRepoForUserPayment()
	paymentSvc := payment.NewPaymentService(paymentRepo, newFakeOrderRepositoryForPayment())

	router := gin.New()
	router.POST("/user/payments", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		createPaymentHandler(c, paymentSvc)
	})

	req := httptest.NewRequest(http.MethodPost, "/user/payments", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

func TestGetPaymentStatusHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	paymentRepo := newMockPaymentRepoForUserPayment()
	paymentSvc := payment.NewPaymentService(paymentRepo, newFakeOrderRepositoryForPayment())

	router := gin.New()
	router.GET("/user/payments/:id", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		getPaymentStatusHandler(c, paymentSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/user/payments/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp model.APIResponse[payment.PaymentStatusResponse]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("Expected success=true, got %v", resp.Success)
	}
}

func TestGetPaymentStatusHandler_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	paymentRepo := newMockPaymentRepoForUserPayment()
	paymentSvc := payment.NewPaymentService(paymentRepo, newFakeOrderRepositoryForPayment())

	router := gin.New()
	router.GET("/user/payments/:id", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		getPaymentStatusHandler(c, paymentSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/user/payments/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

func TestGetPaymentStatusHandler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	paymentRepo := newMockPaymentRepoForUserPayment()
	paymentSvc := payment.NewPaymentService(paymentRepo, newFakeOrderRepositoryForPayment())

	router := gin.New()
	router.GET("/user/payments/:id", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		getPaymentStatusHandler(c, paymentSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/user/payments/9999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected status 404, got %d", w.Code)
	}
}

func TestCancelPaymentHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	paymentRepo := newMockPaymentRepoForUserPayment()
	paymentSvc := payment.NewPaymentService(paymentRepo, newFakeOrderRepositoryForPayment())

	router := gin.New()
	router.POST("/user/payments/:id/cancel", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		cancelPaymentHandler(c, paymentSvc)
	})

	req := httptest.NewRequest(http.MethodPost, "/user/payments/1/cancel", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCancelPaymentHandler_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	paymentRepo := newMockPaymentRepoForUserPayment()
	paymentSvc := payment.NewPaymentService(paymentRepo, newFakeOrderRepositoryForPayment())

	router := gin.New()
	router.POST("/user/payments/:id/cancel", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		cancelPaymentHandler(c, paymentSvc)
	})

	req := httptest.NewRequest(http.MethodPost, "/user/payments/invalid/cancel", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}
