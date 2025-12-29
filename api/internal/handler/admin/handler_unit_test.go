// Package admin provides unit tests for admin handlers.
package admin

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	adminservice "gamelink/internal/service/admin"
	"gamelink/pkg/apierr"
)

// MockAdminService is a mock implementation of AdminService for testing.
type MockAdminService struct {
	// User operations
	GetUserFunc                func(ctx any, id uint64) (*model.User, error)
	ListUsersFunc              func(ctx any, opts any) ([]model.User, *model.Pagination, error)
	CreateUserFunc             func(ctx any, input adminservice.CreateUserInput) (*model.User, error)
	UpdateUserFunc             func(ctx any, id uint64, input adminservice.UpdateUserInput) (*model.User, error)
	DeleteUserFunc             func(ctx any, id uint64) error
	UpdateUserStatusFunc       func(ctx any, id uint64, status model.UserStatus) (*model.User, error)
	UpdateUserRoleFunc         func(ctx any, id uint64, role model.Role) (*model.User, error)
	GetUserStatsFunc           func(ctx any) (*adminservice.UserStatsResponse, error)
	RegisterUserAndPlayerFunc  func(ctx any, u adminservice.CreateUserInput, p adminservice.CreatePlayerInput) (*model.User, *model.Player, error)
	ListOperationLogsFunc      func(ctx any, entityType string, entityID uint64, opts any) ([]model.OperationLog, *model.Pagination, error)

	// Player operations
	GetPlayerFunc                  func(ctx any, id uint64) (*model.Player, error)
	ListPlayersFunc                func(ctx any, page, pageSize int) ([]model.Player, *model.Pagination, error)
	ListPlayersPagedFunc           func(ctx any, page, pageSize int) ([]model.Player, *model.Pagination, error)
	ListPlayersPagedWithFilterFunc func(ctx any, page, pageSize int, keyword string, status *model.VerificationStatus) ([]model.Player, *model.Pagination, error)
	CreatePlayerFunc               func(ctx any, input adminservice.CreatePlayerInput) (*model.Player, error)
	UpdatePlayerFunc               func(ctx any, id uint64, input adminservice.UpdatePlayerInput) (*model.Player, error)
	DeletePlayerFunc               func(ctx any, id uint64) error
	UpdatePlayerVerificationFunc   func(ctx any, id uint64, status model.VerificationStatus) (*model.Player, error)

	// Order operations
	GetOrderFunc         func(ctx any, id uint64) (*model.Order, error)
	ListOrdersFunc       func(ctx any, opts any) ([]model.Order, *model.Pagination, error)
	CreateOrderFunc      func(ctx any, input adminservice.CreateOrderInput) (*model.Order, error)
	UpdateOrderFunc      func(ctx any, id uint64, input adminservice.UpdateOrderInput) (*model.Order, error)
	DeleteOrderFunc      func(ctx any, id uint64) error
	AssignOrderFunc      func(ctx any, orderID, playerID uint64) (*model.Order, error)
	ConfirmOrderFunc     func(ctx any, id uint64, note string) (*model.Order, error)
	StartOrderFunc       func(ctx any, id uint64, note string) (*model.Order, error)
	CompleteOrderFunc    func(ctx any, id uint64, note string) (*model.Order, error)
	CancelOrderFunc      func(ctx any, id uint64, reason string) (*model.Order, error)
	RefundOrderFunc      func(ctx any, id uint64, input any) (*model.Order, error)
	ReviewOrderFunc      func(ctx any, id uint64, input any) (*model.Review, error)

	// Payment operations
	GetPaymentFunc   func(ctx any, id uint64) (*model.Payment, error)
	ListPaymentsFunc func(ctx any, opts any) ([]model.Payment, *model.Pagination, error)
	CreatePaymentFunc func(ctx any, input any) (*model.Payment, error)
	UpdatePaymentFunc func(ctx any, id uint64, input any) (*model.Payment, error)
	DeletePaymentFunc func(ctx any, id uint64) error
	RefundPaymentFunc func(ctx any, id uint64, input any) (*model.Payment, error)
	CapturePaymentFunc func(ctx any, id uint64, input any) (*model.Payment, error)
}

func (m *MockAdminService) GetUser(ctx any, id uint64) (*model.User, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(ctx, id)
	}
	return &model.User{Base: model.Base{ID: id}}, nil
}

func (m *MockAdminService) ListUsersWithOptions(ctx any, opts any) ([]model.User, *model.Pagination, error) {
	if m.ListUsersFunc != nil {
		return m.ListUsersFunc(ctx, opts)
	}
	return []model.User{}, &model.Pagination{Page: 1, PageSize: 20, Total: 0}, nil
}

func (m *MockAdminService) CreateUser(ctx any, input adminservice.CreateUserInput) (*model.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, input)
	}
	return &model.User{Base: model.Base{ID: 1}}, nil
}

func (m *MockAdminService) UpdateUser(ctx any, id uint64, input adminservice.UpdateUserInput) (*model.User, error) {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(ctx, id, input)
	}
	return &model.User{Base: model.Base{ID: id}}, nil
}

func (m *MockAdminService) DeleteUser(ctx any, id uint64) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(ctx, id)
	}
	return nil
}

func (m *MockAdminService) UpdateUserStatus(ctx any, id uint64, status model.UserStatus) (*model.User, error) {
	if m.UpdateUserStatusFunc != nil {
		return m.UpdateUserStatusFunc(ctx, id, status)
	}
	return &model.User{Base: model.Base{ID: id}}, nil
}

func (m *MockAdminService) UpdateUserRole(ctx any, id uint64, role model.Role) (*model.User, error) {
	if m.UpdateUserRoleFunc != nil {
		return m.UpdateUserRoleFunc(ctx, id, role)
	}
	return &model.User{Base: model.Base{ID: id}}, nil
}

func (m *MockAdminService) GetUserStats(ctx any) (*adminservice.UserStatsResponse, error) {
	if m.GetUserStatsFunc != nil {
		return m.GetUserStatsFunc(ctx)
	}
	return &adminservice.UserStatsResponse{TotalUsers: 100}, nil
}

func (m *MockAdminService) RegisterUserAndPlayer(ctx any, u adminservice.CreateUserInput, p adminservice.CreatePlayerInput) (*model.User, *model.Player, error) {
	if m.RegisterUserAndPlayerFunc != nil {
		return m.RegisterUserAndPlayerFunc(ctx, u, p)
	}
	return &model.User{Base: model.Base{ID: 1}}, &model.Player{Base: model.Base{ID: 1}}, nil
}

func (m *MockAdminService) ListOperationLogs(ctx any, entityType string, entityID uint64, opts any) ([]model.OperationLog, *model.Pagination, error) {
	if m.ListOperationLogsFunc != nil {
		return m.ListOperationLogsFunc(ctx, entityType, entityID, opts)
	}
	return []model.OperationLog{}, &model.Pagination{Page: 1, PageSize: 20}, nil
}

func (m *MockAdminService) GetPlayer(ctx any, id uint64) (*model.Player, error) {
	if m.GetPlayerFunc != nil {
		return m.GetPlayerFunc(ctx, id)
	}
	return &model.Player{Base: model.Base{ID: id}}, nil
}

func (m *MockAdminService) ListPlayersPaged(ctx any, page, pageSize int) ([]model.Player, *model.Pagination, error) {
	if m.ListPlayersPagedFunc != nil {
		return m.ListPlayersPagedFunc(ctx, page, pageSize)
	}
	return []model.Player{}, &model.Pagination{Page: 1, PageSize: pageSize}, nil
}

func (m *MockAdminService) ListPlayersPagedWithFilter(ctx any, page, pageSize int, keyword string, status *model.VerificationStatus) ([]model.Player, *model.Pagination, error) {
	if m.ListPlayersPagedWithFilterFunc != nil {
		return m.ListPlayersPagedWithFilterFunc(ctx, page, pageSize, keyword, status)
	}
	return []model.Player{}, &model.Pagination{Page: 1, PageSize: pageSize}, nil
}

func (m *MockAdminService) CreatePlayer(ctx any, input adminservice.CreatePlayerInput) (*model.Player, error) {
	if m.CreatePlayerFunc != nil {
		return m.CreatePlayerFunc(ctx, input)
	}
	return &model.Player{Base: model.Base{ID: 1}}, nil
}

func (m *MockAdminService) UpdatePlayer(ctx any, id uint64, input adminservice.UpdatePlayerInput) (*model.Player, error) {
	if m.UpdatePlayerFunc != nil {
		return m.UpdatePlayerFunc(ctx, id, input)
	}
	return &model.Player{Base: model.Base{ID: id}}, nil
}

func (m *MockAdminService) DeletePlayer(ctx any, id uint64) error {
	if m.DeletePlayerFunc != nil {
		return m.DeletePlayerFunc(ctx, id)
	}
	return nil
}

func (m *MockAdminService) UpdatePlayerVerification(ctx any, id uint64, status model.VerificationStatus) (*model.Player, error) {
	if m.UpdatePlayerVerificationFunc != nil {
		return m.UpdatePlayerVerificationFunc(ctx, id, status)
	}
	return &model.Player{Base: model.Base{ID: id}}, nil
}

func (m *MockAdminService) GetOrder(ctx any, id uint64) (*model.Order, error) {
	if m.GetOrderFunc != nil {
		return m.GetOrderFunc(ctx, id)
	}
	return &model.Order{Base: model.Base{ID: id}}, nil
}

func (m *MockAdminService) ListOrders(ctx any, opts any) ([]model.Order, *model.Pagination, error) {
	if m.ListOrdersFunc != nil {
		return m.ListOrdersFunc(ctx, opts)
	}
	return []model.Order{}, &model.Pagination{Page: 1, PageSize: 20}, nil
}

func (m *MockAdminService) CreateOrder(ctx any, input adminservice.CreateOrderInput) (*model.Order, error) {
	if m.CreateOrderFunc != nil {
		return m.CreateOrderFunc(ctx, input)
	}
	return &model.Order{Base: model.Base{ID: 1}}, nil
}

func (m *MockAdminService) AssignOrder(ctx any, orderID, playerID uint64) (*model.Order, error) {
	if m.AssignOrderFunc != nil {
		return m.AssignOrderFunc(ctx, orderID, playerID)
	}
	return &model.Order{Base: model.Base{ID: orderID}}, nil
}

func (m *MockAdminService) ConfirmOrder(ctx any, id uint64, note string) (*model.Order, error) {
	if m.ConfirmOrderFunc != nil {
		return m.ConfirmOrderFunc(ctx, id, note)
	}
	return &model.Order{Base: model.Base{ID: id}}, nil
}

func (m *MockAdminService) StartOrder(ctx any, id uint64, note string) (*model.Order, error) {
	if m.StartOrderFunc != nil {
		return m.StartOrderFunc(ctx, id, note)
	}
	return &model.Order{Base: model.Base{ID: id}}, nil
}

func (m *MockAdminService) CompleteOrder(ctx any, id uint64, note string) (*model.Order, error) {
	if m.CompleteOrderFunc != nil {
		return m.CompleteOrderFunc(ctx, id, note)
	}
	return &model.Order{Base: model.Base{ID: id}}, nil
}

func (m *MockAdminService) CancelOrder(ctx any, id uint64, reason string) (*model.Order, error) {
	if m.CancelOrderFunc != nil {
		return m.CancelOrderFunc(ctx, id, reason)
	}
	return &model.Order{Base: model.Base{ID: id}}, nil
}

func (m *MockAdminService) RefundOrder(ctx any, id uint64, input any) (*model.Order, error) {
	if m.RefundOrderFunc != nil {
		return m.RefundOrderFunc(ctx, id, input)
	}
	return &model.Order{Base: model.Base{ID: id}}, nil
}

func (m *MockAdminService) ReviewOrder(ctx any, id uint64, input any) (*model.Review, error) {
	if m.ReviewOrderFunc != nil {
		return m.ReviewOrderFunc(ctx, id, input)
	}
	return &model.Review{Base: model.Base{ID: 1}}, nil
}

func (m *MockAdminService) GetPayment(ctx any, id uint64) (*model.Payment, error) {
	if m.GetPaymentFunc != nil {
		return m.GetPaymentFunc(ctx, id)
	}
	return &model.Payment{Base: model.Base{ID: id}}, nil
}

func (m *MockAdminService) ListPayments(ctx any, opts any) ([]model.Payment, *model.Pagination, error) {
	if m.ListPaymentsFunc != nil {
		return m.ListPaymentsFunc(ctx, opts)
	}
	return []model.Payment{}, &model.Pagination{Page: 1, PageSize: 20}, nil
}

func (m *MockAdminService) CreatePayment(ctx any, input any) (*model.Payment, error) {
	if m.CreatePaymentFunc != nil {
		return m.CreatePaymentFunc(ctx, input)
	}
	return &model.Payment{Base: model.Base{ID: 1}}, nil
}

func (m *MockAdminService) UpdatePayment(ctx any, id uint64, input any) (*model.Payment, error) {
	if m.UpdatePaymentFunc != nil {
		return m.UpdatePaymentFunc(ctx, id, input)
	}
	return &model.Payment{Base: model.Base{ID: id}}, nil
}

func (m *MockAdminService) DeletePayment(ctx any, id uint64) error {
	if m.DeletePaymentFunc != nil {
		return m.DeletePaymentFunc(ctx, id)
	}
	return nil
}

func (m *MockAdminService) RefundPayment(ctx any, id uint64, input any) (*model.Payment, error) {
	if m.RefundPaymentFunc != nil {
		return m.RefundPaymentFunc(ctx, id, input)
	}
	return &model.Payment{Base: model.Base{ID: id}}, nil
}

func (m *MockAdminService) CapturePayment(ctx any, id uint64, input any) (*model.Payment, error) {
	if m.CapturePaymentFunc != nil {
		return m.CapturePaymentFunc(ctx, id, input)
	}
	return &model.Payment{Base: model.Base{ID: id}}, nil
}

// Helper function to set up test context
func setupTestContext() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// Helper to create test request
func makeTestRequest(router *gin.Engine, method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// Helper to parse response
func parseResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	return response
}

// ============================================================================
// UserHandler Tests
// ============================================================================

func TestUserHandler_Unit_GetUser_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		GetUserFunc: func(ctx any, id uint64) (*model.User, error) {
			return &model.User{
				Base:   model.Base{ID: id},
				Name:   "Test User",
				Phone:  "13800138000",
				Role:   model.RoleUser,
				Status: model.UserStatusActive,
			}, nil
		},
	}
	handler := NewUserHandler(mockSvc)
	router.GET("/users/:id", handler.GetUser)

	w := makeTestRequest(router, "GET", "/users/123", nil, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(123), data["id"])
	assert.Equal(t, "Test User", data["name"])
}

func TestUserHandler_Unit_GetUser_NotFound(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		GetUserFunc: func(ctx any, id uint64) (*model.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	handler := NewUserHandler(mockSvc)
	router.GET("/users/:id", handler.GetUser)

	w := makeTestRequest(router, "GET", "/users/999", nil, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_Unit_GetUser_InvalidID(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{}
	handler := NewUserHandler(mockSvc)
	router.GET("/users/:id", handler.GetUser)

	w := makeTestRequest(router, "GET", "/users/invalid", nil, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Unit_ListUsers_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		ListUsersFunc: func(ctx any, opts any) ([]model.User, *model.Pagination, error) {
			return []model.User{
					{Base: model.Base{ID: 1}, Name: "User 1"},
					{Base: model.Base{ID: 2}, Name: "User 2"},
				}, &model.Pagination{
					Page:       1,
					PageSize:   20,
					Total:      2,
					TotalPages: 1,
				}, nil
		},
	}
	handler := NewUserHandler(mockSvc)
	router.GET("/users", handler.ListUsers)

	w := makeTestRequest(router, "GET", "/users?page=1&pageSize=20", nil, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].([]interface{})
	assert.Len(t, data, 2)
}

func TestUserHandler_Unit_ListUsers_WithFilters(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		ListUsersFunc: func(ctx any, opts any) ([]model.User, *model.Pagination, error) {
			return []model.User{}, &model.Pagination{Page: 1, PageSize: 20, Total: 0}, nil
		},
	}
	handler := NewUserHandler(mockSvc)
	router.GET("/users", handler.ListUsers)

	w := makeTestRequest(router, "GET", "/users?role=user&status=active&keyword=test", nil, nil)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_Unit_CreateUser_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		CreateUserFunc: func(ctx any, input adminservice.CreateUserInput) (*model.User, error) {
			return &model.User{
				Base:   model.Base{ID: 1},
				Name:   input.Name,
				Phone:  input.Phone,
				Role:   input.Role,
				Status: input.Status,
			}, nil
		},
	}
	handler := NewUserHandler(mockSvc)
	router.POST("/users", handler.CreateUser)

	payload := map[string]interface{}{
		"phone":    "13800138000",
		"email":    "test@example.com",
		"password": "password123",
		"name":     "Test User",
		"role":     "user",
		"status":   "active",
	}

	w := makeTestRequest(router, "POST", "/users", payload, nil)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "Test User", data["name"])
}

func TestUserHandler_Unit_CreateUser_InvalidEmail(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{}
	handler := NewUserHandler(mockSvc)
	router.POST("/users", handler.CreateUser)

	payload := map[string]interface{}{
		"email":    "invalid-email",
		"password": "password123",
		"name":     "Test User",
		"role":     "user",
		"status":   "active",
	}

	w := makeTestRequest(router, "POST", "/users", payload, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Unit_CreateUser_InvalidPhone(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{}
	handler := NewUserHandler(mockSvc)
	router.POST("/users", handler.CreateUser)

	payload := map[string]interface{}{
		"phone":    "123",
		"email":    "test@example.com",
		"password": "password123",
		"name":     "Test User",
		"role":     "user",
		"status":   "active",
	}

	w := makeTestRequest(router, "POST", "/users", payload, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Unit_UpdateUser_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		UpdateUserFunc: func(ctx any, id uint64, input adminservice.UpdateUserInput) (*model.User, error) {
			return &model.User{
				Base:   model.Base{ID: id},
				Name:   *input.Name,
				Phone:  *input.Phone,
				Status: *input.Status,
			}, nil
		},
	}
	handler := NewUserHandler(mockSvc)
	router.PUT("/users/:id", handler.UpdateUser)

	payload := map[string]interface{}{
		"name":   "Updated Name",
		"phone":  "13900139000",
		"role":   "user",
		"status": "active",
	}

	w := makeTestRequest(router, "PUT", "/users/123", payload, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "Updated Name", data["name"])
}

func TestUserHandler_Unit_UpdateUser_NotFound(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		UpdateUserFunc: func(ctx any, id uint64, input adminservice.UpdateUserInput) (*model.User, error) {
			return nil, apierr.NotFound("user not found")
		},
	}
	handler := NewUserHandler(mockSvc)
	router.PUT("/users/:id", handler.UpdateUser)

	payload := map[string]interface{}{
		"name":   "Updated Name",
		"role":   "user",
		"status": "active",
	}

	w := makeTestRequest(router, "PUT", "/users/999", payload, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_Unit_DeleteUser_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		DeleteUserFunc: func(ctx any, id uint64) error {
			return nil
		},
	}
	handler := NewUserHandler(mockSvc)
	router.DELETE("/users/:id", handler.DeleteUser)

	w := makeTestRequest(router, "DELETE", "/users/123", nil, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
}

func TestUserHandler_Unit_DeleteUser_NotFound(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		DeleteUserFunc: func(ctx any, id uint64) error {
			return apierr.NotFound("user not found")
		},
	}
	handler := NewUserHandler(mockSvc)
	router.DELETE("/users/:id", handler.DeleteUser)

	w := makeTestRequest(router, "DELETE", "/users/999", nil, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_Unit_UpdateUserStatus_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		UpdateUserStatusFunc: func(ctx any, id uint64, status model.UserStatus) (*model.User, error) {
			return &model.User{
				Base:   model.Base{ID: id},
				Status: status,
			}, nil
		},
	}
	handler := NewUserHandler(mockSvc)
	router.PUT("/users/:id/status", handler.UpdateUserStatus)

	payload := map[string]string{
		"status": "inactive",
	}

	w := makeTestRequest(router, "PUT", "/users/123/status", payload, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "inactive", data["status"])
}

func TestUserHandler_Unit_UpdateUserRole_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		UpdateUserRoleFunc: func(ctx any, id uint64, role model.Role) (*model.User, error) {
			return &model.User{
				Base: model.Base{ID: id},
				Role: role,
			}, nil
		},
	}
	handler := NewUserHandler(mockSvc)
	router.PUT("/users/:id/role", handler.UpdateUserRole)

	payload := map[string]string{
		"role": "admin",
	}

	w := makeTestRequest(router, "PUT", "/users/123/role", payload, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "admin", data["role"])
}

func TestUserHandler_Unit_GetUserStats_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		GetUserStatsFunc: func(ctx any) (*adminservice.UserStatsResponse, error) {
			return &adminservice.UserStatsResponse{
				TotalUsers:    100,
				ActiveUsers:   80,
				InactiveUsers: 20,
			}, nil
		},
	}
	handler := NewUserHandler(mockSvc)
	router.GET("/users/stats", handler.GetUserStats)

	w := makeTestRequest(router, "GET", "/users/stats", nil, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(100), data["total_users"])
	assert.Equal(t, float64(80), data["active_users"])
}

func TestUserHandler_Unit_CreateUserWithPlayer_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		RegisterUserAndPlayerFunc: func(ctx any, u adminservice.CreateUserInput, p adminservice.CreatePlayerInput) (*model.User, *model.Player, error) {
			return &model.User{
					Base: model.Base{ID: 1},
					Name: u.Name,
				}, &model.Player{
					Base:               model.Base{ID: 1},
					Nickname:           p.Nickname,
					VerificationStatus: p.VerificationStatus,
				}, nil
		},
	}
	handler := NewUserHandler(mockSvc)
	router.POST("/users/with-player", handler.CreateUserWithPlayer)

	payload := map[string]interface{}{
		"phone":    "13800138000",
		"email":    "test@example.com",
		"password": "password123",
		"name":     "Test User",
		"role":     "player",
		"status":   "active",
		"player": map[string]interface{}{
			"nickname":               "TestPlayer",
			"bio":                    "Test bio",
			"verification_status":    "pending",
			"hourly_rate_cents":      5000,
			"main_game_id":           1,
		},
	}

	w := makeTestRequest(router, "POST", "/users/with-player", payload, nil)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.NotNil(t, data["user"])
	assert.NotNil(t, data["player"])
}

func TestUserHandler_Unit_BatchDeleteUsers_Success(t *testing.T) {
	router := setupTestContext()
	deleteCount := 0
	mockSvc := &MockAdminService{
		DeleteUserFunc: func(ctx any, id uint64) error {
			deleteCount++
			return nil
		},
	}
	handler := NewUserHandler(mockSvc)
	router.POST("/users/batch-delete", handler.BatchDeleteUsers)

	payload := map[string]interface{}{
		"ids": []uint64{1, 2, 3},
	}

	w := makeTestRequest(router, "POST", "/users/batch-delete", payload, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(3), data["deleted"])
}

func TestUserHandler_Unit_BatchDeleteUsers_PartialFailure(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		DeleteUserFunc: func(ctx any, id uint64) error {
			if id == 2 {
				return errors.New("delete failed")
			}
			return nil
		},
	}
	handler := NewUserHandler(mockSvc)
	router.POST("/users/batch-delete", handler.BatchDeleteUsers)

	payload := map[string]interface{}{
		"ids": []uint64{1, 2, 3},
	}

	w := makeTestRequest(router, "POST", "/users/batch-delete", payload, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["deleted"])
	assert.Equal(t, float64(1), data["failed"])
}

// ============================================================================
// PlayerHandler Tests
// ============================================================================

func TestPlayerHandler_Unit_GetPlayer_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		GetPlayerFunc: func(ctx any, id uint64) (*model.Player, error) {
			return &model.Player{
				Base:     model.Base{ID: id},
				Nickname: "Test Player",
			}, nil
		},
	}
	handler := NewPlayerHandler(mockSvc)
	router.GET("/players/:id", handler.GetPlayer)

	w := makeTestRequest(router, "GET", "/players/123", nil, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(123), data["id"])
	assert.Equal(t, "Test Player", data["nickname"])
}

func TestPlayerHandler_Unit_GetPlayer_NotFound(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		GetPlayerFunc: func(ctx any, id uint64) (*model.Player, error) {
			return nil, adminservice.ErrNotFound
		},
	}
	handler := NewPlayerHandler(mockSvc)
	router.GET("/players/:id", handler.GetPlayer)

	w := makeTestRequest(router, "GET", "/players/999", nil, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPlayerHandler_Unit_ListPlayers_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		ListPlayersPagedFunc: func(ctx any, page, pageSize int) ([]model.Player, *model.Pagination, error) {
			return []model.Player{
					{Base: model.Base{ID: 1}, Nickname: "Player 1"},
					{Base: model.Base{ID: 2}, Nickname: "Player 2"},
				}, &model.Pagination{
					Page:       1,
					PageSize:   20,
					Total:      2,
					TotalPages: 1,
				}, nil
		},
	}
	handler := NewPlayerHandler(mockSvc)
	router.GET("/players", handler.ListPlayers)

	w := makeTestRequest(router, "GET", "/players?page=1&pageSize=20", nil, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].([]interface{})
	assert.Len(t, data, 2)
}

func TestPlayerHandler_Unit_ListPlayers_WithFilter(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		ListPlayersPagedWithFilterFunc: func(ctx any, page, pageSize int, keyword string, status *model.VerificationStatus) ([]model.Player, *model.Pagination, error) {
			return []model.Player{}, &model.Pagination{Page: 1, PageSize: 20, Total: 0}, nil
		},
	}
	handler := NewPlayerHandler(mockSvc)
	router.GET("/players", handler.ListPlayers)

	w := makeTestRequest(router, "GET", "/players?keyword=test&status=pending", nil, nil)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPlayerHandler_Unit_CreatePlayer_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		CreatePlayerFunc: func(ctx any, input adminservice.CreatePlayerInput) (*model.Player, error) {
			return &model.Player{
				Base:               model.Base{ID: 1},
				UserID:             input.UserID,
				Nickname:           input.Nickname,
				VerificationStatus: input.VerificationStatus,
			}, nil
		},
	}
	handler := NewPlayerHandler(mockSvc)
	router.POST("/players", handler.CreatePlayer)

	payload := map[string]interface{}{
		"user_id":             123,
		"nickname":            "New Player",
		"bio":                 "Player bio",
		"hourly_rate_cents":   5000,
		"main_game_id":        1,
		"verification_status": "pending",
	}

	w := makeTestRequest(router, "POST", "/players", payload, nil)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "New Player", data["nickname"])
}

func TestPlayerHandler_Unit_UpdatePlayer_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		UpdatePlayerFunc: func(ctx any, id uint64, input adminservice.UpdatePlayerInput) (*model.Player, error) {
			return &model.Player{
				Base:     model.Base{ID: id},
				Nickname: *input.Nickname,
				Bio:      *input.Bio,
			}, nil
		},
	}
	handler := NewPlayerHandler(mockSvc)
	router.PUT("/players/:id", handler.UpdatePlayer)

	payload := map[string]interface{}{
		"nickname":            "Updated Player",
		"bio":                 "Updated bio",
		"verification_status": "verified",
	}

	w := makeTestRequest(router, "PUT", "/players/123", payload, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "Updated Player", data["nickname"])
}

func TestPlayerHandler_Unit_DeletePlayer_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		DeletePlayerFunc: func(ctx any, id uint64) error {
			return nil
		},
	}
	handler := NewPlayerHandler(mockSvc)
	router.DELETE("/players/:id", handler.DeletePlayer)

	w := makeTestRequest(router, "DELETE", "/players/123", nil, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
}

func TestPlayerHandler_Unit_DeletePlayer_NotFound(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		DeletePlayerFunc: func(ctx any, id uint64) error {
			return adminservice.ErrNotFound
		},
	}
	handler := NewPlayerHandler(mockSvc)
	router.DELETE("/players/:id", handler.DeletePlayer)

	w := makeTestRequest(router, "DELETE", "/players/999", nil, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ============================================================================
// OrderHandler Tests
// ============================================================================

func TestOrderHandler_Unit_GetOrder_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		GetOrderFunc: func(ctx any, id uint64) (*model.Order, error) {
			return &model.Order{
				Base:        model.Base{ID: id},
				Title:       "Test Order",
				TotalPriceCents: 10000,
			}, nil
		},
	}
	handler := NewOrderHandler(mockSvc)
	router.GET("/orders/:id", handler.GetOrder)

	w := makeTestRequest(router, "GET", "/orders/123", nil, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(123), data["id"])
	assert.Equal(t, "Test Order", data["title"])
}

func TestOrderHandler_Unit_GetOrder_NotFound(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		GetOrderFunc: func(ctx any, id uint64) (*model.Order, error) {
			return nil, apierr.NotFound("order not found")
		},
	}
	handler := NewOrderHandler(mockSvc)
	router.GET("/orders/:id", handler.GetOrder)

	w := makeTestRequest(router, "GET", "/orders/999", nil, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestOrderHandler_Unit_ListOrders_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		ListOrdersFunc: func(ctx any, opts any) ([]model.Order, *model.Pagination, error) {
			return []model.Order{
					{Base: model.Base{ID: 1}, Title: "Order 1"},
					{Base: model.Base{ID: 2}, Title: "Order 2"},
				}, &model.Pagination{
					Page:       1,
					PageSize:   20,
					Total:      2,
					TotalPages: 1,
				}, nil
		},
	}
	handler := NewOrderHandler(mockSvc)
	router.GET("/orders", handler.ListOrders)

	w := makeTestRequest(router, "GET", "/orders?page=1&pageSize=20", nil, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].([]interface{})
	assert.Len(t, data, 2)
}

func TestOrderHandler_Unit_CreateOrder_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		CreateOrderFunc: func(ctx any, input adminservice.CreateOrderInput) (*model.Order, error) {
			return &model.Order{
				Base:            model.Base{ID: 1},
				UserID:          input.UserID,
				GameID:          input.GameID,
				Title:           input.Title,
				TotalPriceCents: input.TotalPriceCents,
			}, nil
		},
	}
	handler := NewOrderHandler(mockSvc)
	router.POST("/orders", handler.CreateOrder)

	payload := map[string]interface{}{
		"user_id":           123,
		"game_id":           1,
		"title":             "Test Order",
		"total_price_cents": 10000,
		"currency":          "CNY",
	}

	w := makeTestRequest(router, "POST", "/orders", payload, nil)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "Test Order", data["title"])
}

func TestOrderHandler_Unit_AssignOrder_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		AssignOrderFunc: func(ctx any, orderID, playerID uint64) (*model.Order, error) {
			return &model.Order{
				Base:     model.Base{ID: orderID},
				PlayerID: &playerID,
			}, nil
		},
	}
	handler := NewOrderHandler(mockSvc)
	router.POST("/orders/:id/assign", handler.AssignOrder)

	payload := map[string]interface{}{
		"player_id": 456,
	}

	w := makeTestRequest(router, "POST", "/orders/123/assign", payload, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.NotNil(t, data["player_id"])
}

func TestOrderHandler_Unit_ConfirmOrder_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		ConfirmOrderFunc: func(ctx any, id uint64, note string) (*model.Order, error) {
			return &model.Order{
				Base:   model.Base{ID: id},
				Status: model.OrderStatusConfirmed,
			}, nil
		},
	}
	handler := NewOrderHandler(mockSvc)
	router.POST("/orders/:id/confirm", handler.ConfirmOrder)

	w := makeTestRequest(router, "POST", "/orders/123/confirm", map[string]string{"note": "test"}, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "confirmed", data["status"])
}

func TestOrderHandler_Unit_StartOrder_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		StartOrderFunc: func(ctx any, id uint64, note string) (*model.Order, error) {
			return &model.Order{
				Base:   model.Base{ID: id},
				Status: model.OrderStatusInProgress,
			}, nil
		},
	}
	handler := NewOrderHandler(mockSvc)
	router.POST("/orders/:id/start", handler.StartOrder)

	w := makeTestRequest(router, "POST", "/orders/123/start", map[string]string{"note": "starting"}, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "in_progress", data["status"])
}

func TestOrderHandler_Unit_CompleteOrder_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		CompleteOrderFunc: func(ctx any, id uint64, note string) (*model.Order, error) {
			return &model.Order{
				Base:   model.Base{ID: id},
				Status: model.OrderStatusCompleted,
			}, nil
		},
	}
	handler := NewOrderHandler(mockSvc)
	router.POST("/orders/:id/complete", handler.CompleteOrder)

	w := makeTestRequest(router, "POST", "/orders/123/complete", map[string]string{"note": "done"}, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "completed", data["status"])
}

// ============================================================================
// PaymentHandler Tests (simplified mock-based tests)
// ============================================================================

// Payment handler tests would follow similar patterns
func TestPaymentHandler_Unit_GetPayment_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		GetPaymentFunc: func(ctx any, id uint64) (*model.Payment, error) {
			return &model.Payment{
				Base:            model.Base{ID: id},
				AmountCents:     10000,
				Status:          model.PaymentStatusSucceeded,
			}, nil
		},
	}
	handler := NewPaymentHandler(mockSvc)
	router.GET("/payments/:id", handler.GetPayment)

	w := makeTestRequest(router, "GET", "/payments/123", nil, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(123), data["id"])
}

func TestPaymentHandler_Unit_GetPayment_NotFound(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		GetPaymentFunc: func(ctx any, id uint64) (*model.Payment, error) {
			return nil, apierr.NotFound("payment not found")
		},
	}
	handler := NewPaymentHandler(mockSvc)
	router.GET("/payments/:id", handler.GetPayment)

	w := makeTestRequest(router, "GET", "/payments/999", nil, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPaymentHandler_Unit_ListPayments_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		ListPaymentsFunc: func(ctx any, opts any) ([]model.Payment, *model.Pagination, error) {
			return []model.Payment{
					{Base: model.Base{ID: 1}, AmountCents: 5000},
					{Base: model.Base{ID: 2}, AmountCents: 10000},
				}, &model.Pagination{
					Page:       1,
					PageSize:   20,
					Total:      2,
					TotalPages: 1,
				}, nil
		},
	}
	handler := NewPaymentHandler(mockSvc)
	router.GET("/payments", handler.ListPayments)

	w := makeTestRequest(router, "GET", "/payments?page=1&pageSize=20", nil, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].([]interface{})
	assert.Len(t, data, 2)
}

func TestPaymentHandler_Unit_CreatePayment_Success(t *testing.T) {
	router := setupTestContext()
	mockSvc := &MockAdminService{
		CreatePaymentFunc: func(ctx any, input any) (*model.Payment, error) {
			return &model.Payment{
				Base:        model.Base{ID: 1},
				AmountCents: 10000,
				Status:      model.PaymentStatusPending,
			}, nil
		},
	}
	handler := NewPaymentHandler(mockSvc)
	router.POST("/payments", handler.CreatePayment)

	payload := map[string]interface{}{
		"order_id":      123,
		"amount_cents":  10000,
		"currency":      "CNY",
		"method":        "alipay",
	}

	w := makeTestRequest(router, "POST", "/payments", payload, nil)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp["success"].(bool))
}

// ============================================================================
// Table-Driven Tests for Common Patterns
// ============================================================================

func TestUserHandler_Unit_IDParsing_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{"Valid ID", "/users/123", http.StatusOK},
		{"Invalid ID - String", "/users/abc", http.StatusBadRequest},
		{"Invalid ID - Negative", "/users/-1", http.StatusBadRequest},
		{"Invalid ID - Zero", "/users/0", http.StatusOK}, // 0 is valid uint64
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestContext()
			mockSvc := &MockAdminService{
				GetUserFunc: func(ctx any, id uint64) (*model.User, error) {
					return &model.User{Base: model.Base{ID: id}}, nil
				},
			}
			handler := NewUserHandler(mockSvc)
			router.GET("/users/:id", handler.GetUser)

			w := makeTestRequest(router, "GET", tt.url, nil, nil)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestUserHandler_Unit_ValidationErrors_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		payload        map[string]interface{}
		wantStatus     int
		shouldBeSuccess bool
	}{
		{
			name: "Valid User Creation",
			payload: map[string]interface{}{
				"phone":    "13800138000",
				"email":    "test@example.com",
				"password": "password123",
				"name":     "Test User",
				"role":     "user",
				"status":   "active",
			},
			wantStatus:     http.StatusCreated,
			shouldBeSuccess: true,
		},
		{
			name: "Missing Required Fields",
			payload: map[string]interface{}{
				"phone": "13800138000",
			},
			wantStatus:     http.StatusBadRequest,
			shouldBeSuccess: false,
		},
		{
			name: "Invalid Email",
			payload: map[string]interface{}{
				"email":    "not-an-email",
				"password": "password123",
				"name":     "Test User",
				"role":     "user",
				"status":   "active",
			},
			wantStatus:     http.StatusBadRequest,
			shouldBeSuccess: false,
		},
		{
			name: "Invalid Phone",
			payload: map[string]interface{}{
				"phone":    "123",
				"password": "password123",
				"name":     "Test User",
				"role":     "user",
				"status":   "active",
			},
			wantStatus:     http.StatusBadRequest,
			shouldBeSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestContext()
			mockSvc := &MockAdminService{
				CreateUserFunc: func(ctx any, input adminservice.CreateUserInput) (*model.User, error) {
					return &model.User{Base: model.Base{ID: 1}}, nil
				},
			}
			handler := NewUserHandler(mockSvc)
			router.POST("/users", handler.CreateUser)

			w := makeTestRequest(router, "POST", "/users", tt.payload, nil)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.shouldBeSuccess {
				resp := parseResponse(t, w)
				assert.True(t, resp["success"].(bool))
			}
		})
	}
}

// ============================================================================
// Error Response Tests
// ============================================================================

func TestErrorHandling_StandardizedResponses(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func() *MockAdminService
		path       string
		method     string
		wantStatus int
		wantCode   string
	}{
		{
			name: "Not Found Error",
			setupMock: func() *MockAdminService {
				return &MockAdminService{
					GetUserFunc: func(ctx any, id uint64) (*model.User, error) {
						return nil, apierr.NotFound("user not found")
					},
				}
			},
			path:       "/users/999",
			method:     "GET",
			wantStatus: http.StatusNotFound,
			wantCode:   "not_found",
		},
		{
			name: "Validation Error",
			setupMock: func() *MockAdminService {
				return &MockAdminService{
					CreateUserFunc: func(ctx any, input adminservice.CreateUserInput) (*model.User, error) {
						return nil, apierr.ErrValidation
					},
				}
			},
			path:       "/users",
			method:     "POST",
			wantStatus: http.StatusBadRequest,
			wantCode:   "bad_request",
		},
		{
			name: "Internal Server Error",
			setupMock: func() *MockAdminService {
				return &MockAdminService{
					GetUserFunc: func(ctx any, id uint64) (*model.User, error) {
						return nil, errors.New("database connection failed")
					},
				}
			},
			path:       "/users/1",
			method:     "GET",
			wantStatus: http.StatusInternalServerError,
			wantCode:   "internal_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestContext()
			mockSvc := tt.setupMock()
			handler := NewUserHandler(mockSvc)

			if tt.method == "GET" {
				router.GET("/users/:id", handler.GetUser)
				w := makeTestRequest(router, "GET", tt.path, nil, nil)
				assert.Equal(t, tt.wantStatus, w.Code)
			} else if tt.method == "POST" {
				router.POST("/users", handler.CreateUser)
				payload := map[string]interface{}{
					"password": "password123",
					"name":     "Test User",
					"role":     "user",
					"status":   "active",
				}
				w := makeTestRequest(router, "POST", tt.path, payload, nil)
				assert.Equal(t, tt.wantStatus, w.Code)
			}
		})
	}
}
