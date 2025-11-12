package user

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"

    "gamelink/internal/model"
)

func TestUserPayment_Create_BadJSON(t *testing.T) {
    payRepo := newFakePaymentRepo()
    ordRepo := &fakeOrderRepoPay{items: map[uint64]*model.Order{1: {Base: model.Base{ID:1}, UserID:1, TotalPriceCents:1000, Currency:model.CurrencyCNY, Status:model.OrderStatusPending}}}
    r := setupUserPaymentRouter(payRepo, ordRepo)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/user/payments", bytes.NewReader([]byte("{")))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

func TestUserPayment_Cancel_InvalidID(t *testing.T) {
    payRepo := newFakePaymentRepo()
    ordRepo := &fakeOrderRepoPay{items: map[uint64]*model.Order{}}
    r := setupUserPaymentRouter(payRepo, ordRepo)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/user/payments/abc/cancel", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

func TestUserPayment_Cancel_Unauthorized(t *testing.T) {
    payRepo := newFakePaymentRepo()
    ordRepo := &fakeOrderRepoPay{items: map[uint64]*model.Order{1: {Base: model.Base{ID:1}, UserID:2, TotalPriceCents:1000, Currency:model.CurrencyCNY, Status:model.OrderStatusPending}}}
    r := setupUserPaymentRouter(payRepo, ordRepo)
    payRepo.items[3] = &model.Payment{Base: model.Base{ID:3}, UserID:2, OrderID:1, Method:model.PaymentMethodWeChat, AmountCents:1000, Currency:model.CurrencyCNY, Status:model.PaymentStatusPending}
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/user/payments/3/cancel", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusInternalServerError { t.Fatalf("%d", w.Code) }
}
