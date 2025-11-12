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

type wechatProvider struct{}

func (wechatProvider) Refund(ctx context.Context, p *model.Payment, reason string) (string, json.RawMessage, time.Time, error) {
    now := time.Now()
    tradeNo := fmt.Sprintf("wx_refund_%d_%d", p.ID, now.Unix())
    raw := map[string]interface{}{
        "channel": "wechat",
        "payment_id": p.ID,
        "refund_reason": reason,
        "refunded_at": now.Unix(),
    }
    b, _ := json.Marshal(raw)
    return tradeNo, json.RawMessage(b), now, nil
}

type alipayProvider struct{}

func (alipayProvider) Refund(ctx context.Context, p *model.Payment, reason string) (string, json.RawMessage, time.Time, error) {
    now := time.Now()
    tradeNo := fmt.Sprintf("ali_refund_%d_%d", p.ID, now.Unix())
    raw := map[string]interface{}{
        "channel": "alipay",
        "payment_id": p.ID,
        "refund_reason": reason,
        "refunded_at": now.Unix(),
    }
    b, _ := json.Marshal(raw)
    return tradeNo, json.RawMessage(b), now, nil
}

type genericProvider struct{}

func (genericProvider) Refund(ctx context.Context, p *model.Payment, reason string) (string, json.RawMessage, time.Time, error) {
    now := time.Now()
    tradeNo := fmt.Sprintf("refund_%d_%d", p.ID, now.Unix())
    raw := map[string]interface{}{
        "channel": "generic",
        "payment_id": p.ID,
        "refund_reason": reason,
        "refunded_at": now.Unix(),
    }
    b, _ := json.Marshal(raw)
    return tradeNo, json.RawMessage(b), now, nil
}
