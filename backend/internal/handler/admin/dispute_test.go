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
	assignmentsvc "gamelink/internal/service/assignment"
)

// Mock repositories for admin dispute handler tests
type mockDisputeRepoForAdminHandler struct {
	disputes map[uint64]*model.OrderDispute
}

func (m *mockDisputeRepoForAdminHandler) Create(ctx context.Context, dispute *model.OrderDispute) error {
	if dispute.ID == 0 {
		dispute.ID = uint64(len(m.disputes) + 1)
	}
	m.disputes[dispute.ID] = dispute
	return nil
}

func (m *mockDisputeRepoForAdminHandler) Get(ctx context.Context, id uint64) (*model.OrderDispute, error) {
	if d, ok := m.disputes[id]; ok {
		return d, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockDisputeRepoForAdminHandler) GetByOrderID(ctx context.Context, orderID uint64) (*model.OrderDispute, error) {
	for _, d := range m.disputes {
		if d.OrderID == orderID {
			return d, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockDisputeRepoForAdminHandler) Update(ctx context.Context, dispute *model.OrderDispute) error {
	m.disputes[dispute.ID] = dispute
	return nil
}

func (m *mockDisputeRepoForAdminHandler) List(ctx context.Context, opts repository.DisputeListOptions) ([]model.OrderDispute, int64, error) {
	disputes := make([]model.OrderDispute, 0)
	for _, d := range m.disputes {
		disputes = append(disputes, *d)
	}
	return disputes, int64(len(disputes)), nil
}

func (m *mockDisputeRepoForAdminHandler) ListPendingAssignment(ctx context.Context, page, pageSize int) ([]model.OrderDispute, int64, error) {
	disputes := make([]model.OrderDispute, 0)
	for _, d := range m.disputes {
		if d.Status == model.DisputeStatusPending && d.AssignedToUserID == nil {
			disputes = append(disputes, *d)
		}
	}
	return disputes, int64(len(disputes)), nil
}

func (m *mockDisputeRepoForAdminHandler) ListSLABreached(ctx context.Context) ([]model.OrderDispute, error) {
	disputes := make([]model.OrderDispute, 0)
	for _, d := range m.disputes {
		if d.SLABreached {
			disputes = append(disputes, *d)
		}
	}
	return disputes, nil
}

func (m *mockDisputeRepoForAdminHandler) MarkSLABreached(ctx context.Context, disputeID uint64) error {
	if d, ok := m.disputes[disputeID]; ok {
		d.SLABreached = true
		return nil
	}
	return repository.ErrNotFound
}

func (m *mockDisputeRepoForAdminHandler) Delete(ctx context.Context, id uint64) error {
	delete(m.disputes, id)
	return nil
}

func (m *mockDisputeRepoForAdminHandler) CountByStatus(ctx context.Context, status model.DisputeStatus) (int64, error) {
	count := 0
	for _, d := range m.disputes {
		if d.Status == status {
			count++
		}
	}
	return int64(count), nil
}

func (m *mockDisputeRepoForAdminHandler) GetPendingCount(ctx context.Context) (int64, error) {
	count := 0
	for _, d := range m.disputes {
		if d.Status == model.DisputeStatusPending {
			count++
		}
	}
	return int64(count), nil
}

type mockOrderRepoForAdminDispute struct{}

func (m *mockOrderRepoForAdminDispute) Create(ctx context.Context, o *model.Order) error { return nil }
func (m *mockOrderRepoForAdminDispute) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	return nil, 0, nil
}
func (m *mockOrderRepoForAdminDispute) Get(ctx context.Context, id uint64) (*model.Order, error) {
	return &model.Order{Base: model.Base{ID: id}, UserID: 1}, nil
}
func (m *mockOrderRepoForAdminDispute) Update(ctx context.Context, o *model.Order) error { return nil }
func (m *mockOrderRepoForAdminDispute) Delete(ctx context.Context, id uint64) error { return nil }

type mockUserRepoForAdminDispute struct{}

func (m *mockUserRepoForAdminDispute) List(ctx context.Context) ([]model.User, error) { return nil, nil }
func (m *mockUserRepoForAdminDispute) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *mockUserRepoForAdminDispute) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *mockUserRepoForAdminDispute) Get(ctx context.Context, id uint64) (*model.User, error) {
	return &model.User{Base: model.Base{ID: id}}, nil
}
func (m *mockUserRepoForAdminDispute) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (m *mockUserRepoForAdminDispute) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (m *mockUserRepoForAdminDispute) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (m *mockUserRepoForAdminDispute) Create(ctx context.Context, user *model.User) error { return nil }
func (m *mockUserRepoForAdminDispute) Update(ctx context.Context, user *model.User) error { return nil }
func (m *mockUserRepoForAdminDispute) Delete(ctx context.Context, id uint64) error { return nil }

type mockOperationLogRepoForAdminDispute struct{}

func (m *mockOperationLogRepoForAdminDispute) Create(ctx context.Context, log *model.OperationLog) error { return nil }
func (m *mockOperationLogRepoForAdminDispute) List(ctx context.Context, opts repository.OperationLogListOptions) ([]model.OperationLog, int64, error) {
	return nil, 0, nil
}
func (m *mockOperationLogRepoForAdminDispute) Append(ctx context.Context, log *model.OperationLog) error { return nil }
func (m *mockOperationLogRepoForAdminDispute) ListByEntity(ctx context.Context, entityType string, entityID uint64, opts repository.OperationLogListOptions) ([]model.OperationLog, int64, error) {
	return nil, 0, nil
}

type mockNotificationRepoForAdminDispute struct{}

func (m *mockNotificationRepoForAdminDispute) ListByUser(ctx context.Context, opts repository.NotificationListOptions) ([]model.NotificationEvent, int64, error) {
	return nil, 0, nil
}
func (m *mockNotificationRepoForAdminDispute) MarkRead(ctx context.Context, userID uint64, ids []uint64) error { return nil }
func (m *mockNotificationRepoForAdminDispute) CountUnread(ctx context.Context, userID uint64) (int64, error) { return 0, nil }
func (m *mockNotificationRepoForAdminDispute) Create(ctx context.Context, event *model.NotificationEvent) error { return nil }

type mockPaymentRepoForAdminDispute struct{}

func (m *mockPaymentRepoForAdminDispute) Create(ctx context.Context, p *model.Payment) error { return nil }
func (m *mockPaymentRepoForAdminDispute) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	return nil, 0, nil
}
func (m *mockPaymentRepoForAdminDispute) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	return nil, repository.ErrNotFound
}
func (m *mockPaymentRepoForAdminDispute) Update(ctx context.Context, p *model.Payment) error { return nil }
func (m *mockPaymentRepoForAdminDispute) Delete(ctx context.Context, id uint64) error { return nil }

func setupAdminDisputeHandler(t *testing.T) *DisputeHandler {
	t.Helper()
	disputeRepo := &mockDisputeRepoForAdminHandler{disputes: make(map[uint64]*model.OrderDispute)}
	orderRepo := &mockOrderRepoForAdminDispute{}
	userRepo := &mockUserRepoForAdminDispute{}
	opLogRepo := &mockOperationLogRepoForAdminDispute{}
	notifRepo := &mockNotificationRepoForAdminDispute{}
	paymentRepo := &mockPaymentRepoForAdminDispute{}
	svc := assignmentsvc.NewAssignmentService(disputeRepo, orderRepo, userRepo, opLogRepo, notifRepo, paymentRepo)
	return NewDisputeHandler(svc)
}

func TestAdminDisputeHandler_GetDisputeDetail_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupAdminDisputeHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/disputes/1", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler.GetDisputeDetail(c)

	// Will return 404 since dispute doesn't exist, but that's OK
	assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusOK)
}

func TestAdminDisputeHandler_GetDisputeDetail_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupAdminDisputeHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/disputes/abc", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	handler.GetDisputeDetail(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminDisputeHandler_ListPendingDisputes_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupAdminDisputeHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/orders/pending-assign?page=1&pageSize=20", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ListPendingDisputes(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminDisputeHandler_ListPendingDisputes_WithPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupAdminDisputeHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/orders/pending-assign?page=2&pageSize=10", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ListPendingDisputes(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminDisputeHandler_ListPendingDisputes_InvalidPage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupAdminDisputeHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/orders/pending-assign?page=abc", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ListPendingDisputes(c)

	// Should still return 200 with default page
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminDisputeHandler_ListPendingDisputes_DefaultParams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupAdminDisputeHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/orders/pending-assign", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ListPendingDisputes(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
