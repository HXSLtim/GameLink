package admin

import (
    "context"
    "encoding/json"
    "testing"
    "time"

    "gamelink/internal/cache"
    "gamelink/internal/model"
    "gamelink/internal/repository"
    "gamelink/internal/repository/common"
)

type timelineOrders struct{ o model.Order }
func (t timelineOrders) Create(context.Context, *model.Order) error { return nil }
func (t timelineOrders) List(context.Context, repository.OrderListOptions) ([]model.Order, int64, error) { return nil, 0, nil }
func (t timelineOrders) Get(context.Context, uint64) (*model.Order, error) { return &t.o, nil }
func (t timelineOrders) Update(context.Context, *model.Order) error { return nil }
func (t timelineOrders) Delete(context.Context, uint64) error { return nil }

type timelinePayments struct{ items []model.Payment }
func (p timelinePayments) Create(context.Context, *model.Payment) error { return nil }
func (p timelinePayments) List(context.Context, repository.PaymentListOptions) ([]model.Payment, int64, error) { return p.items, int64(len(p.items)), nil }
func (p timelinePayments) Get(context.Context, uint64) (*model.Payment, error) { return nil, repository.ErrNotFound }
func (p timelinePayments) Update(context.Context, *model.Payment) error { return nil }
func (p timelinePayments) Delete(context.Context, uint64) error { return nil }

type timelineOpLogs struct{ logs []model.OperationLog }
func (l timelineOpLogs) Append(context.Context, *model.OperationLog) error { return nil }
func (l timelineOpLogs) ListByEntity(context.Context, string, uint64, repository.OperationLogListOptions) ([]model.OperationLog, int64, error) { return l.logs, int64(len(l.logs)), nil }

type txTL struct{ repos common.Repos }
func (t *txTL) WithTx(ctx context.Context, fn func(r *common.Repos) error) error { return fn(&t.repos) }

func TestGetOrderTimeline_BuildsFromLogsAndPayments(t *testing.T) {
    paidAt := time.Date(2025,1,2,3,4,5,0,time.UTC)
    refundedAt := time.Date(2025,1,3,3,4,5,0,time.UTC)
    pay1 := model.Payment{Base: model.Base{ID: 7}, Status: model.PaymentStatusPaid, Method: model.PaymentMethodWeChat, AmountCents: 1000, PaidAt: &paidAt}
    pay2 := model.Payment{Base: model.Base{ID: 8}, Status: model.PaymentStatusRefunded, Method: model.PaymentMethodAlipay, AmountCents: 500, RefundedAt: &refundedAt}
    o := timelineOrders{o: model.Order{Base: model.Base{ID:1}, Status: model.OrderStatusConfirmed}}
    meta := map[string]any{"note":"服务确认","status": string(model.OrderStatusConfirmed), "from_status": string(model.OrderStatusPending)}
    mb, _ := json.Marshal(meta)
    logs := []model.OperationLog{{Base: model.Base{ID: 1, CreatedAt: time.Date(2025,1,2,2,2,2,0,time.UTC)}, Action: string(model.OpActionConfirm), MetadataJSON: mb}}
    svc := NewAdminService(nil, nil, nil, o, timelinePayments{items: []model.Payment{pay1, pay2}}, nil, cache.NewMemory())
    svc.SetTxManager(&txTL{repos: common.Repos{OpLogs: timelineOpLogs{logs: logs}}})
    items, err := svc.GetOrderTimeline(context.Background(), 1)
    if err != nil { t.Fatalf("%v", err) }
    if len(items) < 3 { t.Fatalf("expected >=3 items, got %d", len(items)) }
}

type txReject struct{}
func (txReject) WithTx(context.Context, func(r *common.Repos) error) error { return repository.ErrNotFound }

func TestGetOrderTimeline_TxRollback(t *testing.T) {
    svc := NewAdminService(nil, nil, nil, timelineOrders{o: model.Order{Base: model.Base{ID:5}, Status: model.OrderStatusConfirmed}}, timelinePayments{items: []model.Payment{}}, nil, cache.NewMemory())
    svc.SetTxManager(txReject{})
    _, err := svc.GetOrderTimeline(context.Background(), 5)
    if err == nil { t.Fatal("expected tx error") }
}

func TestGetOrderRefunds_SummaryAddedWhenNoPaymentMatches(t *testing.T) {
    updated := time.Date(2025,1,4,1,1,1,0,time.UTC)
    o := timelineOrders{o: model.Order{Base: model.Base{ID:2, UpdatedAt: updated}, Status: model.OrderStatusCompleted, RefundAmountCents: 300, RefundReason: "时长补差"}}
    svc := NewAdminService(nil, nil, nil, o, timelinePayments{items: []model.Payment{}}, nil, cache.NewMemory())
    items, err := svc.GetOrderRefunds(context.Background(), 2)
    if err != nil { t.Fatalf("%v", err) }
    hasSummary := false
    for _, it := range items { if it.AmountCents == 300 { hasSummary = true } }
    if !hasSummary { t.Fatal("expected summary refund item") }
}

func TestGetOrderRefunds_NoSummaryWhenPaymentMatches(t *testing.T) {
    refundedAt := time.Date(2025,2,1,0,0,0,0,time.UTC)
    o := timelineOrders{o: model.Order{Base: model.Base{ID:3}, Status: model.OrderStatusCompleted, RefundAmountCents: 500, RefundReason: "补差", RefundedAt: &refundedAt}}
    pays := []model.Payment{{Base: model.Base{ID: 9}, Status: model.PaymentStatusRefunded, AmountCents: 500, RefundedAt: &refundedAt}}
    svc := NewAdminService(nil, nil, nil, o, timelinePayments{items: pays}, nil, cache.NewMemory())
    items, err := svc.GetOrderRefunds(context.Background(), 3)
    if err != nil { t.Fatalf("%v", err) }
    count := 0
    for _, it := range items { if it.AmountCents == 500 { count++ } }
    if count != 1 { t.Fatalf("expected only payment-based refund item, got %d", count) }
}
