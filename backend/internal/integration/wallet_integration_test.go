package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	userhandler "gamelink/internal/handler/user"
	"gamelink/internal/model"
	orderimpl "gamelink/internal/repository/implementations"
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/user"
	"gamelink/internal/repository/wallet"
	paymentservice "gamelink/internal/service/payment"
	walletservice "gamelink/internal/service/wallet"
	"gamelink/internal/testutil"
)

// 钱包集成：充值 -> 查询余额 -> 查询支付列表
func TestWalletRechargeAndBalance(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateWalletModels(t, db)

	// seed user
	userRepo := user.NewUserRepository(db)
	u := &model.User{Name: "WalletUser", Email: "wallet@example.com", Phone: "18000000000", PasswordHash: "x", Role: model.RoleUser, Status: model.UserStatusActive}
	if err := userRepo.Create(context.Background(), u); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	orderRepo := orderimpl.NewOrderRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	walletRepo := wallet.NewWalletRepository(db)
	paymentSvc := paymentservice.NewPaymentService(paymentRepo, orderRepo)
	walletSvc := walletservice.NewService(walletRepo, paymentRepo, orderRepo)

	router := gin.New()
	api := router.Group("/api/v1")
	auth := fakeAuthMiddleware(u.ID)
	userGroup := api.Group("/user")
	userhandler.RegisterWalletRoutes(userGroup, walletSvc, auth)
	userhandler.RegisterPaymentRoutes(userGroup, paymentSvc, auth)

	// 充值
	rechargePayload := map[string]interface{}{
		"amountCents": 5000,
		"method":      "alipay",
	}
	rechargeResp := doJSON(router, http.MethodPost, "/api/v1/user/wallet/recharge", rechargePayload, "")
	if rechargeResp.Code != http.StatusOK {
		t.Fatalf("recharge status=%d body=%s", rechargeResp.Code, rechargeResp.Body.String())
	}
	var rechargeParsed apiResp[walletservice.RechargeResponse]
	if err := json.Unmarshal(rechargeResp.Body.Bytes(), &rechargeParsed); err != nil {
		t.Fatalf("parse recharge resp: %v", err)
	}
	if rechargeParsed.Data.Balance != 5000 {
		t.Fatalf("expected balance 5000, got %d", rechargeParsed.Data.Balance)
	}

	// 查询余额
	balanceResp := doJSON(router, http.MethodGet, "/api/v1/user/wallet/balance", nil, "")
	if balanceResp.Code != http.StatusOK {
		t.Fatalf("balance status=%d body=%s", balanceResp.Code, balanceResp.Body.String())
	}
	var balanceParsed apiResp[model.Wallet]
	_ = json.Unmarshal(balanceResp.Body.Bytes(), &balanceParsed)
	if balanceParsed.Data.BalanceCents != 5000 {
		t.Fatalf("balance mismatch: %d", balanceParsed.Data.BalanceCents)
	}

	// 支付列表
	listResp := doJSON(router, http.MethodGet, "/api/v1/user/payments?status=paid", nil, "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("payments status=%d body=%s", listResp.Code, listResp.Body.String())
	}
	var listParsed apiResp[struct {
		Items []model.Payment `json:"items"`
		Total int64           `json:"total"`
	}]
	_ = json.Unmarshal(listResp.Body.Bytes(), &listParsed)
	if listParsed.Data.Total == 0 || len(listParsed.Data.Items) == 0 {
		t.Fatalf("expected payments, got %+v", listParsed.Data)
	}
}

// 充值非法金额应返回400
func TestWalletRechargeInvalidAmount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateWalletModels(t, db)

	userRepo := user.NewUserRepository(db)
	u := &model.User{Name: "WalletUser", Email: "wallet2@example.com", Phone: "18000000002", PasswordHash: "x", Role: model.RoleUser, Status: model.UserStatusActive}
	_ = userRepo.Create(context.Background(), u)

	orderRepo := orderimpl.NewOrderRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	walletRepo := wallet.NewWalletRepository(db)
	paymentSvc := paymentservice.NewPaymentService(paymentRepo, orderRepo)
	walletSvc := walletservice.NewService(walletRepo, paymentRepo, orderRepo)

	router := gin.New()
	api := router.Group("/api/v1")
	auth := fakeAuthMiddleware(u.ID)
	userGroup := api.Group("/user")
	userhandler.RegisterWalletRoutes(userGroup, walletSvc, auth)
	userhandler.RegisterPaymentRoutes(userGroup, paymentSvc, auth)

	rechargePayload := map[string]interface{}{
		"amountCents": 0, // invalid
		"method":      "alipay",
	}
	resp := doJSON(router, http.MethodPost, "/api/v1/user/wallet/recharge", rechargePayload, "")
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid amount, got %d body=%s", resp.Code, resp.Body.String())
	}
}

// 新用户未充值时余额为0，支付列表为空
func TestWalletBalanceZeroForNewUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateWalletModels(t, db)

	userRepo := user.NewUserRepository(db)
	u := &model.User{Name: "WalletUser4", Email: "wallet4@example.com", Phone: "18000000004", PasswordHash: "x", Role: model.RoleUser, Status: model.UserStatusActive}
	_ = userRepo.Create(context.Background(), u)

	orderRepo := orderimpl.NewOrderRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	walletRepo := wallet.NewWalletRepository(db)
	paymentSvc := paymentservice.NewPaymentService(paymentRepo, orderRepo)
	walletSvc := walletservice.NewService(walletRepo, paymentRepo, orderRepo)

	router := gin.New()
	api := router.Group("/api/v1")
	auth := fakeAuthMiddleware(u.ID)
	userGroup := api.Group("/user")
	userhandler.RegisterWalletRoutes(userGroup, walletSvc, auth)
	userhandler.RegisterPaymentRoutes(userGroup, paymentSvc, auth)

	// 余额为0
	balanceResp := doJSON(router, http.MethodGet, "/api/v1/user/wallet/balance", nil, "")
	var balanceParsed apiResp[model.Wallet]
	_ = json.Unmarshal(balanceResp.Body.Bytes(), &balanceParsed)
	if balanceParsed.Data.BalanceCents != 0 {
		t.Fatalf("expected zero balance, got %d", balanceParsed.Data.BalanceCents)
	}

	// 支付列表为空
	listResp := doJSON(router, http.MethodGet, "/api/v1/user/payments?status=paid", nil, "")
	var listParsed apiResp[struct {
		Items []model.Payment `json:"items"`
		Total int64           `json:"total"`
	}]
	_ = json.Unmarshal(listResp.Body.Bytes(), &listParsed)
	if listParsed.Data.Total != 0 || len(listParsed.Data.Items) != 0 {
		t.Fatalf("expected empty payments, got %+v", listParsed.Data)
	}
}

// 多次充值累加、支付记录应按时间倒序返回
func TestWalletRechargeAccumulationAndPaymentOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateWalletModels(t, db)

	userRepo := user.NewUserRepository(db)
	u := &model.User{Name: "WalletUser3", Email: "wallet3@example.com", Phone: "18000000003", PasswordHash: "x", Role: model.RoleUser, Status: model.UserStatusActive}
	_ = userRepo.Create(context.Background(), u)

	orderRepo := orderimpl.NewOrderRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	walletRepo := wallet.NewWalletRepository(db)
	paymentSvc := paymentservice.NewPaymentService(paymentRepo, orderRepo)
	walletSvc := walletservice.NewService(walletRepo, paymentRepo, orderRepo)

	router := gin.New()
	api := router.Group("/api/v1")
	auth := fakeAuthMiddleware(u.ID)
	userGroup := api.Group("/user")
	userhandler.RegisterWalletRoutes(userGroup, walletSvc, auth)
	userhandler.RegisterPaymentRoutes(userGroup, paymentSvc, auth)

	// 充值两次
	doRecharge := func(amount int64) {
		payload := map[string]interface{}{"amountCents": amount, "method": "wechat"}
		resp := doJSON(router, http.MethodPost, "/api/v1/user/wallet/recharge", payload, "")
		if resp.Code != http.StatusOK {
			t.Fatalf("recharge %d status=%d body=%s", amount, resp.Code, resp.Body.String())
		}
	}
	doRecharge(3000)
	doRecharge(2000)

	// 余额应为 5000
	balanceResp := doJSON(router, http.MethodGet, "/api/v1/user/wallet/balance", nil, "")
	var balanceParsed apiResp[model.Wallet]
	_ = json.Unmarshal(balanceResp.Body.Bytes(), &balanceParsed)
	if balanceParsed.Data.BalanceCents != 5000 {
		t.Fatalf("expected balance 5000, got %d", balanceParsed.Data.BalanceCents)
	}

	// 支付列表应按创建时间倒序，且总数为2
	listResp := doJSON(router, http.MethodGet, "/api/v1/user/payments?status=paid&pageSize=10", nil, "")
	var listParsed apiResp[struct {
		Items []model.Payment `json:"items"`
		Total int64           `json:"total"`
	}]
	_ = json.Unmarshal(listResp.Body.Bytes(), &listParsed)
	if listParsed.Data.Total != 2 || len(listParsed.Data.Items) != 2 {
		t.Fatalf("expected 2 payments, got %+v", listParsed.Data)
	}
	if !listParsed.Data.Items[0].CreatedAt.After(listParsed.Data.Items[1].CreatedAt) && !listParsed.Data.Items[0].CreatedAt.Equal(listParsed.Data.Items[1].CreatedAt) {
		t.Fatalf("payments not ordered desc by created_at: %+v", listParsed.Data.Items)
	}
}

func migrateWalletModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.User{},
		&model.Order{},
		&model.Payment{},
		&model.Wallet{},
	); err != nil {
		t.Fatalf("migrate wallet models: %v", err)
	}
}
