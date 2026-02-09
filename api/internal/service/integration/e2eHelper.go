// Package integration provides end-to-end integration test utilities.
package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// E2EData holds all test data for E2E tests.
type E2EData struct {
	Customer     *TestUserWithCreds
	Player       *TestUserWithCreds
	PlayerUser   *TestUserWithCreds // The user account behind the player
	Admin        *TestUserWithCreds
	CS           *TestUserWithCreds
	Game         *GameData
	ServiceItem  *ServiceItemData
	VIPLevel     *VIPLevelData
	Coupon       *CouponData
	ReferralCode *ReferralData
}

// TestUserWithCreds combines a user with their authentication credentials.
type TestUserWithCreds struct {
	UserID   uint64
	Phone    string
	Password string
	Token    string
	Name     string
	Email    string
	Role     string
}

// GameData holds game information.
type GameData struct {
	GameID   uint64
	GameKey  string
	GameName string
	Category string
}

// ServiceItemData holds service item information.
type ServiceItemData struct {
	ItemID         uint64
	ItemCode       string
	Name           string
	Category       string
	BasePriceCents int64
	CommissionRate float64
}

// VIPLevelData holds VIP level information.
type VIPLevelData struct {
	LevelID     uint64
	Slug        string
	Title       string
	ExpRequired int64
	Discount    float64
}

// CouponData holds coupon information.
type CouponData struct {
	TemplateID   uint64
	CouponID     uint64
	Name         string
	DeductCents  int64
	MinAmount    int64
	ValidityDays int
}

// ReferralData holds referral information.
type ReferralData struct {
	Code        string
	ReferrerID  uint64
	ReferrerAmt int64
	RefereeAmt  int64
}

// HTTPTestClient wraps an httptest.Server for making HTTP requests.
type HTTPTestClient struct {
	Server    *httptest.Server
	BaseURL   string
	Token     string
	Transport *http.Transport
}

// NewHTTPTestClient creates a new HTTP test client.
func NewHTTPTestClient(server *httptest.Server) *HTTPTestClient {
	return &HTTPTestClient{
		Server:  server,
		BaseURL: server.URL,
	}
}

// SetToken sets the authentication token for subsequent requests.
func (c *HTTPTestClient) SetToken(token string) {
	c.Token = token
}

// Post makes a POST request with JSON body.
func (c *HTTPTestClient) Post(t *testing.T, path string, body interface{}, expectStatus int) *http.Response {
	t.Helper()

	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		require.NoError(t, err)
		bodyReader = bytes.NewReader(jsonBytes)
	}

	url := c.BaseURL + path
	req, err := http.NewRequest("POST", url, bodyReader)
	require.NoError(t, err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Equal(t, expectStatus, resp.StatusCode, "unexpected status code for POST %s", path)
	return resp
}

// Get makes a GET request.
func (c *HTTPTestClient) Get(t *testing.T, path string, expectStatus int) *http.Response {
	t.Helper()

	url := c.BaseURL + path
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Equal(t, expectStatus, resp.StatusCode, "unexpected status code for GET %s", path)
	return resp
}

// Put makes a PUT request with JSON body.
func (c *HTTPTestClient) Put(t *testing.T, path string, body interface{}, expectStatus int) *http.Response {
	t.Helper()

	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		require.NoError(t, err)
		bodyReader = bytes.NewReader(jsonBytes)
	}

	url := c.BaseURL + path
	req, err := http.NewRequest("PUT", url, bodyReader)
	require.NoError(t, err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Equal(t, expectStatus, resp.StatusCode, "unexpected status code for PUT %s", path)
	return resp
}

// Delete makes a DELETE request.
func (c *HTTPTestClient) Delete(t *testing.T, path string, expectStatus int) *http.Response {
	t.Helper()

	url := c.BaseURL + path
	req, err := http.NewRequest("DELETE", url, nil)
	require.NoError(t, err)

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Equal(t, expectStatus, resp.StatusCode, "unexpected status code for DELETE %s", path)
	return resp
}

// DecodeJSON decodes JSON response body into target.
func DecodeJSON[T any](t *testing.T, resp *http.Response) T {
	t.Helper()

	var result T
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	return result
}

// LoginRequest represents a login request.
type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// LoginResponse represents a login response.
type LoginResponse struct {
	Code int `json:"code"`
	Data struct {
		Token string `json:"token"`
		User  struct {
			ID    uint64 `json:"id"`
			Name  string `json:"name"`
			Phone string `json:"phone"`
			Role  string `json:"role"`
		} `json:"user"`
	} `json:"data"`
	Message string `json:"message"`
}

// Login performs login and returns token.
func (c *HTTPTestClient) Login(t *testing.T, phone, password string) (string, uint64) {
	t.Helper()

	req := LoginRequest{
		Phone:    phone,
		Password: password,
	}

	resp := c.Post(t, "/api/v1/auth/login", req, http.StatusOK)
	result := DecodeJSON[LoginResponse](t, resp)

	require.NotEmpty(t, result.Data.Token, "login should return token")
	return result.Data.Token, result.Data.User.ID
}

// SetupE2EData creates comprehensive test data for E2E tests.
func SetupE2EData(t *testing.T, db *gorm.DB) *E2EData {
	t.Helper()

	data := &E2EData{}

	// Create users
	customerUser := CreateTestUserWithPassword(t, db, "customer", "CustomerPass123!")
	data.Customer = &TestUserWithCreds{
		UserID:   customerUser.ID,
		Phone:    customerUser.Phone,
		Password: "CustomerPass123!",
		Name:     customerUser.Name,
		Email:    customerUser.Email,
		Role:     string(customerUser.Role),
	}

	playerUser := CreateTestUserWithPassword(t, db, "player_user", "PlayerPass123!")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	data.PlayerUser = &TestUserWithCreds{
		UserID:   playerUser.ID,
		Phone:    playerUser.Phone,
		Password: "PlayerPass123!",
		Name:     playerUser.Name,
		Email:    playerUser.Email,
		Role:     string(playerUser.Role),
	}
	data.Player = &TestUserWithCreds{
		UserID:   testPlayer.UserID,
		Phone:    playerUser.Phone,
		Password: "PlayerPass123!",
		Name:     playerUser.Name,
		Email:    playerUser.Email,
		Role:     "player",
	}

	adminUser := CreateTestUserWithPassword(t, db, "admin", "AdminPass123!")
	adminUser.Role = "admin"
	db.Save(adminUser)
	data.Admin = &TestUserWithCreds{
		UserID:   adminUser.ID,
		Phone:    adminUser.Phone,
		Password: "AdminPass123!",
		Name:     adminUser.Name,
		Email:    adminUser.Email,
		Role:     "admin",
	}

	csUser := CreateTestUserWithPassword(t, db, "cs", "CSPass123!")
	data.CS = &TestUserWithCreds{
		UserID:   csUser.ID,
		Phone:    csUser.Phone,
		Password: "CSPass123!",
		Name:     csUser.Name,
		Email:    csUser.Email,
		Role:     "admin",
	}

	// Create game
	testGame := CreateTestGame(t, db, "TestGame")
	data.Game = &GameData{
		GameID:   testGame.ID,
		GameKey:  testGame.Key,
		GameName: testGame.Name,
		Category: testGame.Category,
	}

	// Create service item
	serviceItem := CreateTestServiceItem(t, db, testGame, "陪玩服务", 5000)
	data.ServiceItem = &ServiceItemData{
		ItemID:         serviceItem.ID,
		ItemCode:       serviceItem.ItemCode,
		Name:           serviceItem.Name,
		Category:       serviceItem.Category,
		BasePriceCents: serviceItem.BasePriceCents,
		CommissionRate: serviceItem.CommissionRate,
	}

	// Create VIP level
	vipLevel := CreateTestVipLevel(t, db, "gold", 1000)
	data.VIPLevel = &VIPLevelData{
		LevelID:     vipLevel.ID,
		Slug:        vipLevel.Slug,
		Title:       vipLevel.Title,
		ExpRequired: vipLevel.ExpRequired,
		Discount:    vipLevel.OrderDiscount,
	}

	// Create coupon
	couponTemplate := CreateTestCouponTemplate(t, db, "SAVE10", 1000)
	data.Coupon = &CouponData{
		TemplateID:   couponTemplate.ID,
		Name:         couponTemplate.Name,
		DeductCents:  couponTemplate.DeductAmountCents,
		MinAmount:    couponTemplate.MinAmountCents,
		ValidityDays: couponTemplate.ValidityDays,
	}

	return data
}

// CreateOrderRequest represents an order creation request.
type CreateOrderRequest struct {
	PlayerID       uint64 `json:"playerId"`
	GameID         uint64 `json:"gameId"`
	ItemID         uint64 `json:"itemId"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	DurationHours  int    `json:"durationHours"`
	ScheduledStart string `json:"scheduledStart,omitempty"`
}

// OrderResponse represents an order response.
type OrderResponse struct {
	Code int `json:"code"`
	Data struct {
		OrderID     uint64 `json:"orderId"`
		OrderNo     string `json:"orderNo"`
		PriceCents  int64  `json:"priceCents"`
		NeedPayment bool   `json:"needPayment"`
		Status      string `json:"status"`
	} `json:"data"`
	Message string `json:"message"`
}

// CreateOrder creates an order via HTTP API.
func (c *HTTPTestClient) CreateOrder(t *testing.T, req CreateOrderRequest) OrderResponse {
	t.Helper()

	resp := c.Post(t, "/api/v1/orders", req, http.StatusOK)
	return DecodeJSON[OrderResponse](t, resp)
}

// GetOrderResponse represents a get order response.
type GetOrderResponse struct {
	Code int `json:"code"`
	Data struct {
		ID              uint64  `json:"id"`
		OrderNo         string  `json:"orderNo"`
		Status          string  `json:"status"`
		TotalPriceCents int64   `json:"totalPriceCents"`
		PlayerID        *uint64 `json:"playerId"`
		UserID          uint64  `json:"userId"`
	} `json:"data"`
	Message string `json:"message"`
}

// GetOrder retrieves an order by ID.
func (c *HTTPTestClient) GetOrder(t *testing.T, orderID uint64) GetOrderResponse {
	t.Helper()

	resp := c.Get(t, fmt.Sprintf("/api/v1/orders/%d", orderID), http.StatusOK)
	return DecodeJSON[GetOrderResponse](t, resp)
}

// AcceptOrder accepts an order (player endpoint).
func (c *HTTPTestClient) AcceptOrder(t *testing.T, orderID uint64) {
	t.Helper()

	c.Post(t, fmt.Sprintf("/api/v1/player/orders/%d/accept", orderID), nil, http.StatusOK)
}

// StartOrder starts an order (player endpoint).
func (c *HTTPTestClient) StartOrder(t *testing.T, orderID uint64) {
	t.Helper()

	c.Post(t, fmt.Sprintf("/api/v1/player/orders/%d/start", orderID), nil, http.StatusOK)
}

// CompleteOrderRequest represents complete order request.
type CompleteOrderRequest struct {
	Remark string `json:"remark,omitempty"`
}

// CompleteOrder completes an order (player endpoint).
func (c *HTTPTestClient) CompleteOrder(t *testing.T, orderID uint64, req CompleteOrderRequest) {
	t.Helper()

	c.Post(t, fmt.Sprintf("/api/v1/player/orders/%d/complete", orderID), req, http.StatusOK)
}

// CreatePaymentRequest represents payment creation request.
type CreatePaymentRequest struct {
	OrderID           uint64 `json:"orderId"`
	Method            string `json:"method"`
	WalletAmountCents int64  `json:"walletAmountCents,omitempty"`
	ThirdPartyMethod  string `json:"thirdPartyMethod,omitempty"`
}

// PaymentResponse represents payment response.
type PaymentResponse struct {
	Code int `json:"code"`
	Data struct {
		PaymentID        uint64      `json:"paymentId"`
		WalletPaidDirect bool        `json:"walletPaidDirect"`
		WalletDeducted   int64       `json:"walletDeducted"`
		ThirdPartyAmount int64       `json:"thirdPartyAmount"`
		PayInfo          interface{} `json:"payInfo"`
	} `json:"data"`
	Message string `json:"message"`
}

// CreatePayment creates a payment.
func (c *HTTPTestClient) CreatePayment(t *testing.T, req CreatePaymentRequest) PaymentResponse {
	t.Helper()

	resp := c.Post(t, "/api/v1/payments", req, http.StatusOK)
	return DecodeJSON[PaymentResponse](t, resp)
}

// InitiateDisputeRequest represents dispute initiation request.
type InitiateDisputeRequest struct {
	OrderID      uint64   `json:"orderId"`
	Type         string   `json:"type"`
	Reason       string   `json:"reason"`
	EvidenceText string   `json:"evidenceText,omitempty"`
	EvidenceURLs []string `json:"evidenceUrls,omitempty"`
}

// DisputeResponse represents dispute response.
type DisputeResponse struct {
	Code int `json:"code"`
	Data struct {
		DisputeID   uint64 `json:"disputeId"`
		TraceID     string `json:"traceId"`
		SLADeadline string `json:"slaDeadline"`
	} `json:"data"`
	Message string `json:"message"`
}

// InitiateDispute initiates a dispute.
func (c *HTTPTestClient) InitiateDispute(t *testing.T, req InitiateDisputeRequest) DisputeResponse {
	t.Helper()

	resp := c.Post(t, "/api/v1/disputes", req, http.StatusOK)
	return DecodeJSON[DisputeResponse](t, resp)
}

// ResolveDisputeRequest represents dispute resolution request.
type ResolveDisputeRequest struct {
	Resolution    string `json:"resolution"`
	ResolveRemark string `json:"resolveRemark"`
	RefundAmount  int64  `json:"refundAmount,omitempty"`
}

// ResolveDispute resolves a dispute (admin endpoint).
func (c *HTTPTestClient) ResolveDispute(t *testing.T, disputeID uint64, req ResolveDisputeRequest) {
	t.Helper()

	c.Post(t, fmt.Sprintf("/api/v1/admin/disputes/%d/resolve", disputeID), req, http.StatusOK)
}

// RechargeWalletRequest represents wallet recharge request.
type RechargeWalletRequest struct {
	AmountCents int64  `json:"amountCents"`
	Method      string `json:"method"`
}

// RechargeWallet recharges user wallet.
func (c *HTTPTestClient) RechargeWallet(t *testing.T, req RechargeWalletRequest) {
	t.Helper()

	c.Post(t, "/api/v1/wallet/recharge", req, http.StatusOK)
}

// WalletBalanceResponse represents wallet balance response.
type WalletBalanceResponse struct {
	Code int `json:"code"`
	Data struct {
		BalanceCents int64 `json:"balanceCents"`
		FrozenCents  int64 `json:"frozenCents"`
	} `json:"data"`
	Message string `json:"message"`
}

// GetWalletBalance retrieves wallet balance.
func (c *HTTPTestClient) GetWalletBalance(t *testing.T) WalletBalanceResponse {
	t.Helper()

	resp := c.Get(t, "/api/v1/wallet/balance", http.StatusOK)
	return DecodeJSON[WalletBalanceResponse](t, resp)
}

// WithdrawRequest represents withdraw request.
type WithdrawRequest struct {
	AmountCents   int64  `json:"amountCents"`
	Method        string `json:"method"`
	AccountName   string `json:"accountName"`
	AccountNumber string `json:"accountNumber"`
	BankName      string `json:"bankName,omitempty"`
}

// WithdrawResponse represents withdraw response.
type WithdrawResponse struct {
	Code int `json:"code"`
	Data struct {
		WithdrawID uint64 `json:"withdrawId"`
		Status     string `json:"status"`
	} `json:"data"`
	Message string `json:"message"`
}

// CreateWithdraw creates a withdraw request.
func (c *HTTPTestClient) CreateWithdraw(t *testing.T, req WithdrawRequest) WithdrawResponse {
	t.Helper()

	resp := c.Post(t, "/api/v1/player/withdraws", req, http.StatusOK)
	return DecodeJSON[WithdrawResponse](t, resp)
}

// ApproveWithdraw approves a withdraw request (admin endpoint).
func (c *HTTPTestClient) ApproveWithdraw(t *testing.T, withdrawID uint64) {
	t.Helper()

	c.Post(t, fmt.Sprintf("/api/v1/admin/withdraws/%d/approve", withdrawID), nil, http.StatusOK)
}

// WaitForCondition waits for a condition to be true or timeout.
func WaitForCondition(t *testing.T, condition func() bool, timeout time.Duration, msg string) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		if condition() {
			return
		}

		if time.Now().After(deadline) {
			t.Fatalf("timeout waiting for condition: %s", msg)
		}

		<-ticker.C
	}
}

// AssertOrderStatus asserts order status in database.
func AssertOrderStatus(t *testing.T, db *gorm.DB, orderID uint64, expectedStatus string) {
	t.Helper()

	var order struct {
		Status string
	}
	err := db.Table("orders").Select("status").Where("id = ?", orderID).First(&order).Error
	require.NoError(t, err)
	require.Equal(t, expectedStatus, order.Status)
}

// AssertPaymentStatus asserts payment status in database.
func AssertPaymentStatus(t *testing.T, db *gorm.DB, paymentID uint64, expectedStatus string) {
	t.Helper()

	var payment struct {
		Status string
	}
	err := db.Table("payments").Select("status").Where("id = ?", paymentID).First(&payment).Error
	require.NoError(t, err)
	require.Equal(t, expectedStatus, payment.Status)
}

// AssertWalletBalance asserts wallet balance in database.
func AssertWalletBalance(t *testing.T, db *gorm.DB, userID uint64, expectedBalance, expectedFrozen int64) {
	t.Helper()

	var wallet struct {
		BalanceCents int64
		FrozenCents  int64
	}
	err := db.Table("wallets").Select("balance_cents, frozen_cents").Where("user_id = ?", userID).First(&wallet).Error
	require.NoError(t, err)
	require.Equal(t, expectedBalance, wallet.BalanceCents, "balance mismatch")
	require.Equal(t, expectedFrozen, wallet.FrozenCents, "frozen mismatch")
}

// AssertCommissionRecord asserts commission record exists.
func AssertCommissionRecord(t *testing.T, db *gorm.DB, orderID uint64) {
	t.Helper()

	var count int64
	err := db.Table("commission_records").Where("order_id = ?", orderID).Count(&count).Error
	require.NoError(t, err)
	require.Greater(t, count, int64(0), "commission record should exist")
}

// AssertDisputeStatus asserts dispute status in database.
func AssertDisputeStatus(t *testing.T, db *gorm.DB, disputeID uint64, expectedStatus string) {
	t.Helper()

	var dispute struct {
		Status string
	}
	err := db.Table("order_disputes").Select("status").Where("id = ?", disputeID).First(&dispute).Error
	require.NoError(t, err)
	require.Equal(t, expectedStatus, dispute.Status)
}

// TimeTravel modifies timestamps for testing time-based operations.
func TimeTravel(t *testing.T, db *gorm.DB, orderID uint64, daysAgo int) {
	t.Helper()

	completedAt := time.Now().Add(-time.Duration(daysAgo) * 24 * time.Hour)
	err := db.Table("orders").Where("id = ?", orderID).Update("completed_at", completedAt).Error
	require.NoError(t, err)
}
