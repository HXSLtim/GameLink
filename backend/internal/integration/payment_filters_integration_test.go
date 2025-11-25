package integration

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	userhandler "gamelink/internal/handler/user"
	"gamelink/internal/model"
	"gamelink/internal/repository/commission"
	"gamelink/internal/repository/game"
	orderimpl "gamelink/internal/repository/implementations"
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/review"
	"gamelink/internal/repository/user"
	ordersvc "gamelink/internal/service/order"
	paymentsvc "gamelink/internal/service/payment"
	"gamelink/internal/testutil"
)

// 支付列表过滤：状态+日期
func TestPaymentListFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateOrderModels(t, db)
	seed := seedOrderData(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	commissionRepo := commission.NewCommissionRepository(db)

	orderService := ordersvc.NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	paymentService := paymentsvc.NewPaymentService(paymentRepo, orderRepo)

	router := gin.New()
	api := router.Group("/api/v1")
	userAuth := fakeAuthMiddleware(seed.userID)
	userhandler.RegisterOrderRoutes(api, orderService, userAuth)
	userhandler.RegisterPaymentRoutes(api, paymentService, userAuth)

	// create two paid payments via handler
	makePayment := func(title string) {
		scheduledStart := time.Now().Add(30 * time.Minute).UTC()
		createOrderPayload := map[string]interface{}{
			"playerId":       seed.playerID,
			"gameId":         seed.gameID,
			"title":          title,
			"scheduledStart": scheduledStart.Format(time.RFC3339),
			"durationHours":  1.0,
		}
		w := doJSON(router, http.MethodPost, "/api/v1/user/orders", createOrderPayload, "")
		if w.Code != http.StatusOK {
			t.Fatalf("create order status=%d body=%s", w.Code, w.Body.String())
		}
		var orderResp apiResp[ordersvc.CreateOrderResponse]
		_ = json.Unmarshal(w.Body.Bytes(), &orderResp)
		payPayload := map[string]interface{}{"orderId": orderResp.Data.OrderID, "method": "alipay"}
		payResp := doJSON(router, http.MethodPost, "/api/v1/user/payments", payPayload, "")
		if payResp.Code != http.StatusOK {
			t.Fatalf("pay status=%d body=%s", payResp.Code, payResp.Body.String())
		}
	}
	makePayment("P1")
	makePayment("P2")

	// status filter paid should return >=2
	listResp := doJSON(router, http.MethodGet, "/api/v1/user/payments?status=paid", nil, "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list status=%d body=%s", listResp.Code, listResp.Body.String())
	}
	var listParsed apiResp[struct {
		Items []model.Payment `json:"items"`
		Total int64           `json:"total"`
	}]
	_ = json.Unmarshal(listResp.Body.Bytes(), &listParsed)
	if listParsed.Data.Total < 2 {
		t.Fatalf("expected at least 2 paid payments, got %+v", listParsed.Data)
	}

	// dateFrom in future should return empty
	future := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	listFuture := doJSON(router, http.MethodGet, "/api/v1/user/payments?dateFrom="+url.QueryEscape(future), nil, "")
	var futureParsed apiResp[struct {
		Items []model.Payment `json:"items"`
		Total int64           `json:"total"`
	}]
	_ = json.Unmarshal(listFuture.Body.Bytes(), &futureParsed)
	if futureParsed.Data.Total != 0 {
		t.Fatalf("expected 0 payments for future filter, got %+v", futureParsed.Data)
	}
}
