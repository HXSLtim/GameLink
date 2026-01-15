// Package testutil provides comprehensive testing utilities for handler layer testing.
// This package includes helpers for making HTTP requests, asserting responses,
// creating test data, and setting up test environments for Gin handlers.
package testutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/service/integration"
)

// ============================================================================
// Test Context & Setup
// ============================================================================

// TestContext holds all testing components for handler tests.
type TestContext struct {
	T          testing.TB
	Router     *gin.Engine
	DB         *gorm.DB
	Service    interface{} // Can be AdminService, OrderService, etc.
	AdminUser  *model.User
	AdminToken string
	BaseURL    string
}

// SetupGinTest initializes Gin test mode and returns a test router.
func SetupGinTest(t testing.TB) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// SetupTestDB initializes the test database with proper cleanup.
func SetupTestDB(t testing.TB) *gorm.DB {
	t.Helper()
	// Type assert to *testing.T
	if tt, ok := t.(*testing.T); ok {
		return integration.SetupTestDB(tt)
	}
	// Fallback for testing.B
	return integration.SetupTestDB(&testing.T{})
}

// ============================================================================
// HTTP Request Helpers
// ============================================================================

// RequestOption configures HTTP requests.
type RequestOption func(*http.Request)

// WithHeader adds a header to the request.
func WithHeader(key, value string) RequestOption {
	return func(r *http.Request) {
		r.Header.Set(key, value)
	}
}

// WithAuth adds Bearer token authentication.
func WithAuth(token string) RequestOption {
	return func(r *http.Request) {
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
}

// WithQuery adds query parameters to the request.
func WithQuery(params map[string]string) RequestOption {
	return func(r *http.Request) {
		q := r.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		r.URL.RawQuery = q.Encode()
	}
}

// MakeRequest creates and executes an HTTP request in the test context.
func MakeRequest(t testing.TB, router *gin.Engine, method, path string, body interface{}, opts ...RequestOption) *httptest.ResponseRecorder {
	t.Helper()

	var bodyReader *bytes.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		require.NoError(t, err, "Failed to marshal request body")
		bodyReader = bytes.NewReader(jsonData)
	} else {
		bodyReader = bytes.NewReader([]byte{})
	}

	// Create request
	req, err := http.NewRequest(method, path, bodyReader)
	require.NoError(t, err, "Failed to create request")

	// Set default headers
	req.Header.Set("Content-Type", "application/json")

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	// Execute request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

// MakeAuthenticatedRequest makes a request with Bearer token authentication.
func MakeAuthenticatedRequest(t testing.TB, router *gin.Engine, method, path, token string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	return MakeRequest(t, router, method, path, body, WithAuth(token))
}

// ============================================================================
// Response Assertions
// ============================================================================

// AssertSuccess asserts that the response indicates success (HTTP 2xx).
func AssertSuccess(t testing.TB, w *httptest.ResponseRecorder, expectedStatus ...int) {
	t.Helper()

	status := http.StatusOK
	if len(expectedStatus) > 0 {
		status = expectedStatus[0]
	}

	assert.Equal(t, status, w.Code, "Status code should match")

	// Verify it's JSON response
	contentType := w.Header().Get("Content-Type")
	assert.True(t, strings.Contains(contentType, "application/json"),
		"Response should be JSON, got: %s", contentType)
}

// AssertError asserts that the response indicates an error (HTTP 4xx/5xx).
func AssertError(t testing.TB, w *httptest.ResponseRecorder, expectedStatus int) {
	t.Helper()

	assert.Equal(t, expectedStatus, w.Code, "Status code should match expected error status")
}

// AssertJSONBody asserts that the response body contains specific JSON fields.
func AssertJSONBody(t testing.TB, w *httptest.ResponseRecorder, expectedFields map[string]interface{}) {
	t.Helper()

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Response should be valid JSON")

	for key, expectedValue := range expectedFields {
		actualValue, exists := response[key]
		assert.True(t, exists, "Response should contain key: %s", key)
		assert.Equal(t, expectedValue, actualValue, "Field %s should match", key)
	}
}

// AssertErrorMessage asserts that the error response contains the expected message.
func AssertErrorMessage(t testing.TB, w *httptest.ResponseRecorder, expectedMessage string) {
	t.Helper()

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Response should be valid JSON")

	message, exists := response["message"]
	assert.True(t, exists, "Error response should contain 'message' field")
	assert.True(t, strings.Contains(fmt.Sprintf("%v", message), expectedMessage),
		"Error message should contain '%s', got: %v", expectedMessage, message)
}

// AssertPagination asserts that the response contains valid pagination data.
func AssertPagination(t testing.TB, w *httptest.ResponseRecorder, expectedPage, expectedPageSize int) {
	t.Helper()

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Response should be valid JSON")

	data, exists := response["data"]
	require.True(t, exists, "Response should contain 'data' field")

	dataMap, ok := data.(map[string]interface{})
	require.True(t, ok, "Data should be an object")

	// Check pagination exists
	pagination, exists := dataMap["pagination"]
	require.True(t, exists, "Data should contain 'pagination' field")

	paginationMap, ok := pagination.(map[string]interface{})
	require.True(t, ok, "Pagination should be an object")

	// Check pagination fields
	assert.Equal(t, float64(expectedPage), paginationMap["page"], "Page number should match")
	assert.Equal(t, float64(expectedPageSize), paginationMap["page_size"], "Page size should match")
}

// AssertListResponse asserts that the response contains a list with pagination.
func AssertListResponse(t testing.TB, w *httptest.ResponseRecorder) (items []interface{}, pagination map[string]interface{}) {
	t.Helper()

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Response should be valid JSON")

	assert.True(t, response["success"].(bool), "Response should indicate success")

	data, exists := response["data"]
	require.True(t, exists, "Response should contain 'data' field")

	dataMap, ok := data.(map[string]interface{})
	require.True(t, ok, "Data should be an object")

	itemsVal, exists := dataMap["items"]
	require.True(t, exists, "Data should contain 'items' field")

	itemsList, ok := itemsVal.([]interface{})
	require.True(t, ok, "Items should be a list")

	paginationVal, exists := dataMap["pagination"]
	require.True(t, exists, "Data should contain 'pagination' field")

	paginationMap, ok := paginationVal.(map[string]interface{})
	require.True(t, ok, "Pagination should be an object")

	return itemsList, paginationMap
}

// AssertDeleted asserts that the response indicates successful deletion.
func AssertDeleted(t testing.TB, w *httptest.ResponseRecorder) {
	t.Helper()

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Response should be valid JSON")

	assert.True(t, response["success"].(bool), "Response should indicate success")
	assert.Equal(t, float64(http.StatusOK), response["code"], "Should return 200 status")
}

// ============================================================================
// Test Data Creation Helpers
// ============================================================================

// CreateAdminUser creates a test admin user in the database.
func CreateAdminUser(t testing.TB, db *gorm.DB, role model.Role) *model.User {
	t.Helper()

	ts := time.Now().UnixNano()
	openID := fmt.Sprintf("openid_%d", ts)
	admin := &model.User{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Phone:         fmt.Sprintf("138%011d", ts%100000000000),
		Name:          fmt.Sprintf("Test Admin %d", ts),
		Email:         fmt.Sprintf("admin_%d@test.com", ts),
		Role:          role,
		Status:        model.UserStatusActive,
		WeChatOpenID:  &openID,                       // 添加唯一的 OpenID（指针类型）
		WeChatUnionID: fmt.Sprintf("unionid_%d", ts), // 添加唯一的 UnionID
	}

	require.NoError(t, db.Create(admin).Error, "Failed to create admin user")
	return admin
}

// CreateSuperAdmin creates a super admin user for testing.
func CreateSuperAdmin(t testing.TB, db *gorm.DB) *model.User {
	t.Helper()

	admin := CreateAdminUser(t, db, model.RoleAdmin)
	require.NoError(t, db.Model(admin).Update("role", string(model.RoleSlugSuperAdmin)).Error,
		"Failed to update admin to super admin role")
	return admin
}

// GenerateTestToken generates a test JWT token for the given user.
func GenerateTestToken(userID uint64) string {
	return fmt.Sprintf("test-token-%d", userID)
}

// CreateTestOrder creates a test order in the database.
// Note: This function automatically creates a ServiceItem for the order since ItemID is required.
func CreateTestOrder(t testing.TB, db *gorm.DB, userID, playerID, gameID uint64, status model.OrderStatus) *model.Order {
	t.Helper()

	var playerIDPtr *uint64
	if playerID != 0 {
		playerIDPtr = &playerID
	}

	var gameIDPtr *uint64
	if gameID != 0 {
		gameIDPtr = &gameID
	}

	// Create a ServiceItem first since ItemID is required (NOT NULL constraint)
	item := CreateTestServiceItem(t, db, gameID, 10000, true)

	order := &model.Order{
		Base: model.Base{
			ExtJSON: "{}",
		},
		UserID:          userID,
		ItemID:          item.ID,
		PlayerID:        playerIDPtr,
		GameID:          gameIDPtr,
		Title:           "Test Order",
		Description:     "Test order description",
		UnitPriceCents:  10000, // $100.00
		TotalPriceCents: 10000, // $100.00
		Currency:        model.CurrencyUSD,
		Status:          status,
		ScheduledStart:  timePtr(time.Now().Add(1 * time.Hour)),
		ScheduledEnd:    timePtr(time.Now().Add(2 * time.Hour)),
	}

	require.NoError(t, db.Create(order).Error, "Failed to create test order")
	return order
}

// CreateTestOrderWithItem creates a test order with a service item in the database.
func CreateTestOrderWithItem(t testing.TB, db *gorm.DB, userID, playerID, itemID uint64, status model.OrderStatus) *model.Order {
	t.Helper()

	var playerIDPtr *uint64
	if playerID != 0 {
		playerIDPtr = &playerID
	}

	order := &model.Order{
		Base: model.Base{
			ExtJSON: "{}",
		},
		UserID:          userID,
		ItemID:          itemID,
		PlayerID:        playerIDPtr,
		Title:           "Test Order",
		Description:     "Test order description",
		TotalPriceCents: 10000, // $100.00
		Currency:        model.CurrencyUSD,
		Status:          status,
		ScheduledStart:  timePtr(time.Now().Add(1 * time.Hour)),
		ScheduledEnd:    timePtr(time.Now().Add(2 * time.Hour)),
	}

	require.NoError(t, db.Create(order).Error, "Failed to create test order")
	return order
}

// CreateTestPayment creates a test payment in the database.
func CreateTestPayment(t testing.TB, db *gorm.DB, orderID, userID uint64, status model.PaymentStatus) *model.Payment {
	t.Helper()

	payment := &model.Payment{
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderID:     orderID,
		UserID:      userID,
		Method:      model.PaymentMethodAlipay,
		AmountCents: 10000,
		Currency:    model.CurrencyUSD,
		Status:      status,
	}

	require.NoError(t, db.Create(payment).Error, "Failed to create test payment")
	return payment
}

// CreateTestGame creates a test game in the database.
func CreateTestGame(t testing.TB, db *gorm.DB) *model.Game {
	t.Helper()

	ts := time.Now().UnixNano()
	game := &model.Game{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Key:     fmt.Sprintf("test-game-%d", ts),
		Name:    fmt.Sprintf("Test Game %d", ts),
		IconURL: "https://example.com/icon.png",
	}

	require.NoError(t, db.Create(game).Error, "Failed to create test game")
	return game
}

// CreateTestPlayer creates a test player in the database.
func CreateTestPlayer(t testing.TB, db *gorm.DB, userID uint64) *model.Player {
	t.Helper()

	ts := time.Now().UnixNano()
	player := &model.Player{
		Base: model.Base{
			ExtJSON: "{}",
		},
		UserID:             userID,
		Nickname:           fmt.Sprintf("Player %d %d", userID, ts),
		Bio:                "Test player bio",
		HourlyRateCents:    5000, // $50.00
		VerificationStatus: model.VerificationVerified,
	}

	require.NoError(t, db.Create(player).Error, "Failed to create test player")
	return player
}

// CreateTestServiceItem creates a test service item in the database.
func CreateTestServiceItem(t testing.TB, db *gorm.DB, gameID uint64, priceCents int64, isActive bool) *model.ServiceItem {
	t.Helper()

	ts := time.Now().UnixNano()
	item := &model.ServiceItem{
		ItemCode:       fmt.Sprintf("ITEM-%d", ts),
		Name:           fmt.Sprintf("Service Item %d", ts),
		BasePriceCents: priceCents,
		ServiceHours:   1,
		IsActive:       isActive,
		GameID:         &gameID,
		Description:    "Test service item",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		CommissionRate: 0.20,
	}

	require.NoError(t, db.Create(item).Error, "Failed to create test service item")
	return item
}

// ============================================================================
// Path & Query Building Helpers
// ============================================================================

// BuildPath builds a URL path with parameters.
func BuildPath(base string, params map[string]string) string {
	path := base
	for key, value := range params {
		path = strings.Replace(path, ":"+key, value, 1)
	}
	return path
}

// BuildURL builds a complete URL with query parameters.
func BuildURL(base string, pathParams, queryParams map[string]string) string {
	path := BuildPath(base, pathParams)

	if len(queryParams) > 0 {
		values := url.Values{}
		for k, v := range queryParams {
			values.Set(k, v)
		}
		path += "?" + values.Encode()
	}

	return path
}

// Uint64ToStr converts uint64 to string for path parameters.
func Uint64ToStr(v uint64) string {
	return strconv.FormatUint(v, 10)
}

// UintToString converts uint64 to string (alias for Uint64ToStr).
func UintToString(v uint64) string {
	return strconv.FormatUint(v, 10)
}

// ============================================================================
// Response Parsing Helpers
// ============================================================================

// ParseJSONResponse parses the response body into the provided interface.
func ParseJSONResponse(t testing.TB, w *httptest.ResponseRecorder, v interface{}) {
	t.Helper()

	err := json.Unmarshal(w.Body.Bytes(), v)
	require.NoError(t, err, "Failed to parse JSON response")
}

// GetResponseData extracts the data field from a standard API response.
func GetResponseData(t testing.TB, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()

	var response struct {
		Success bool                   `json:"success"`
		Code    int                    `json:"code"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to parse response")

	return response.Data
}

// GetResponseList extracts items and pagination from a list response.
func GetResponseList(t testing.TB, w *httptest.ResponseRecorder) ([]interface{}, map[string]interface{}) {
	t.Helper()

	data := GetResponseData(t, w)

	itemsVal, exists := data["items"]
	require.True(t, exists, "Response should contain items")

	itemsList, ok := itemsVal.([]interface{})
	require.True(t, ok, "Items should be a list")

	paginationVal, exists := data["pagination"]
	require.True(t, exists, "Response should contain pagination")

	paginationMap, ok := paginationVal.(map[string]interface{})
	require.True(t, ok, "Pagination should be an object")

	return itemsList, paginationMap
}

// ============================================================================
// Time Helpers
// ============================================================================

// TimePtr returns a pointer to the provided time.
func timePtr(t time.Time) *time.Time {
	return &t
}

// ParseRFC3339Time parses an RFC3339 formatted time string.
func ParseRFC3339Time(t testing.TB, timeStr string) time.Time {
	t.Helper()

	parsed, err := time.Parse(time.RFC3339, timeStr)
	require.NoError(t, err, "Failed to parse time string")
	return parsed
}

// ============================================================================
// Custom Assertions
// ============================================================================

// AssertOrderStatus asserts that an order has the expected status.
func AssertOrderStatus(t testing.TB, db *gorm.DB, orderID uint64, expectedStatus model.OrderStatus) {
	t.Helper()

	var order model.Order
	err := db.First(&order, orderID).Error
	require.NoError(t, err, "Failed to fetch order")
	assert.Equal(t, expectedStatus, order.Status, "Order status should match")
}

// AssertPaymentStatus asserts that a payment has the expected status.
func AssertPaymentStatus(t testing.TB, db *gorm.DB, paymentID uint64, expectedStatus model.PaymentStatus) {
	t.Helper()

	var payment model.Payment
	err := db.First(&payment, paymentID).Error
	require.NoError(t, err, "Failed to fetch payment")
	assert.Equal(t, expectedStatus, payment.Status, "Payment status should match")
}

// AssertUserExists asserts that a user exists in the database.
func AssertUserExists(t testing.TB, db *gorm.DB, userID uint64) {
	t.Helper()

	var user model.User
	err := db.First(&user, userID).Error
	assert.NoError(t, err, "User should exist")
}

// AssertUserDeleted asserts that a user has been deleted from the database.
func AssertUserDeleted(t testing.TB, db *gorm.DB, userID uint64) {
	t.Helper()

	var user model.User
	err := db.First(&user, userID).Error
	assert.Error(t, err, "User should be deleted")
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "Error should be ErrRecordNotFound")
}

// AssertRecordCount asserts the number of records in a table.
func AssertRecordCount(t testing.TB, db *gorm.DB, model interface{}, expectedCount int) {
	t.Helper()

	var count int64
	err := db.Model(model).Count(&count).Error
	require.NoError(t, err, "Failed to count records")
	assert.Equal(t, int64(expectedCount), count, "Record count should match")
}
