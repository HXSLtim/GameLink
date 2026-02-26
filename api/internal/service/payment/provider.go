package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gamelink/internal/model"
)

type ProviderClient interface {
	Refund(ctx context.Context, p *model.Payment, reason string) (providerTradeNo string, providerRaw json.RawMessage, refundedAt time.Time, err error)
}

// wechatOrderCreator defines provider capability for WeChat-like create-order APIs.
type wechatOrderCreator interface {
	CreateOrder(ctx context.Context, orderID, description string, amountCents int64, clientIP string) (map[string]interface{}, error)
}

// alipayOrderCreator defines provider capability for Alipay-like create-order APIs.
type alipayOrderCreator interface {
	CreateOrder(ctx context.Context, orderID, subject string, amountCents int64) (map[string]interface{}, error)
}

type wechatProvider struct{}

func (wechatProvider) Refund(ctx context.Context, p *model.Payment, reason string) (string, json.RawMessage, time.Time, error) {
	now := time.Now()
	tradeNo := fmt.Sprintf("wx_refund_%d_%d", p.ID, now.Unix())
	raw := map[string]interface{}{
		"channel":       "wechat",
		"payment_id":    p.ID,
		"refund_reason": reason,
		"refunded_at":   now.Unix(),
	}
	b, _ := json.Marshal(raw)
	return tradeNo, json.RawMessage(b), now, nil
}

type alipayProvider struct{}

func (alipayProvider) Refund(ctx context.Context, p *model.Payment, reason string) (string, json.RawMessage, time.Time, error) {
	now := time.Now()
	tradeNo := fmt.Sprintf("ali_refund_%d_%d", p.ID, now.Unix())
	raw := map[string]interface{}{
		"channel":       "alipay",
		"payment_id":    p.ID,
		"refund_reason": reason,
		"refunded_at":   now.Unix(),
	}
	b, _ := json.Marshal(raw)
	return tradeNo, json.RawMessage(b), now, nil
}

type genericProvider struct{}

func (genericProvider) Refund(ctx context.Context, p *model.Payment, reason string) (string, json.RawMessage, time.Time, error) {
	now := time.Now()
	tradeNo := fmt.Sprintf("refund_%d_%d", p.ID, now.Unix())
	raw := map[string]interface{}{
		"channel":       "generic",
		"payment_id":    p.ID,
		"refund_reason": reason,
		"refunded_at":   now.Unix(),
	}
	b, _ := json.Marshal(raw)
	return tradeNo, json.RawMessage(b), now, nil
}

// failClosedProvider returns explicit errors for unsafe production fallbacks.
type failClosedProvider struct {
	method model.PaymentMethod
	cause  error
}

func (p failClosedProvider) Refund(ctx context.Context, payment *model.Payment, reason string) (string, json.RawMessage, time.Time, error) {
	_ = ctx
	_ = payment
	_ = reason
	return "", nil, time.Time{}, p.unavailableError()
}

func (p failClosedProvider) unavailableError() error {
	if p.cause != nil {
		return fmt.Errorf("payment provider %s unavailable in production: %w", p.method, p.cause)
	}
	return fmt.Errorf("payment provider %s unavailable in production", p.method)
}
