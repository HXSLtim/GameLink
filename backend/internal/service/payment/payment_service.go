package payment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

var (
	// ErrNotFound 支付记录不存在
	ErrNotFound = repository.ErrNotFound
	// ErrValidation 表示输入校验失败
	ErrValidation = errors.New("validation failed")
	// ErrOrderAlreadyPaid 订单已支付
	ErrOrderAlreadyPaid = errors.New("order already paid")
	// ErrInvalidOrderStatus 订单状态不正确
	ErrInvalidOrderStatus = errors.New("invalid order status")
)

// PaymentService 支付服务
//
// 功能：
// 1. 创建支付
// 2. 查询支付状态
// 3. 取消支付
// 4. 处理支付回调（Mock版本）
type PaymentService struct {
	payments repository.PaymentRepository
	orders   repository.OrderRepository
}

// NewPaymentService 创建支付服务
func NewPaymentService(
	payments repository.PaymentRepository,
	orders repository.OrderRepository,
) *PaymentService {
	return &PaymentService{
		payments: payments,
		orders:   orders,
	}
}

// CreatePaymentRequest 创建支付请求
type CreatePaymentRequest struct {
	OrderID uint64              `json:"orderId" binding:"required"`
	Method  model.PaymentMethod `json:"method" binding:"required,oneof=wechat alipay"`
}

// CreatePaymentResponse 创建支付响应
type CreatePaymentResponse struct {
	PaymentID uint64                 `json:"paymentId"`
	PayInfo   map[string]interface{} `json:"payInfo"` // 支付参数（对接支付SDK）
}

// PaymentStatusResponse 支付状态响应
type PaymentStatusResponse struct {
	PaymentID uint64              `json:"paymentId"`
	OrderID   uint64              `json:"orderId"`
	Status    model.PaymentStatus `json:"status"`
	PaidAt    *time.Time          `json:"paidAt"`
}

// CreatePayment 创建支付
func (s *PaymentService) CreatePayment(ctx context.Context, userID uint64, req CreatePaymentRequest) (*CreatePaymentResponse, error) {
	// 验证订单
	order, err := s.orders.Get(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}

	// 权限检查
	if order.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// 状态检查：只有 pending 状态可以支付
	if order.Status != model.OrderStatusPending {
		return nil, ErrInvalidOrderStatus
	}

	// 检查是否已有支付记录
	orderIDPtr := &req.OrderID
	existingPayments, _, err := s.payments.List(ctx, repository.PaymentListOptions{
		OrderID:  orderIDPtr,
		Page:     1,
		PageSize: 1,
	})
	if err == nil && len(existingPayments) > 0 {
		// 检查是否已支付
		if existingPayments[0].Status == model.PaymentStatusPaid {
			return nil, ErrOrderAlreadyPaid
		}
	}

	// 创建支付记录
	payment := &model.Payment{
		OrderID:     req.OrderID,
		UserID:      userID,
		Method:      req.Method,
		AmountCents: order.PriceCents,
		Currency:    order.Currency,
		Status:      model.PaymentStatusPending,
	}

	if err := s.payments.Create(ctx, payment); err != nil {
		return nil, err
	}

	// Mock: 生成支付参数
	payInfo := s.generateMockPayInfo(payment.ID, req.Method, order.PriceCents)

	// Mock: 自动标记为已支付（仅用于测试）
	// 在生产环境中，这里应该等待支付回调
	if err := s.mockPaymentSuccess(ctx, payment.ID, order); err != nil {
		return nil, err
	}

	return &CreatePaymentResponse{
		PaymentID: payment.ID,
		PayInfo:   payInfo,
	}, nil
}

// GetPaymentStatus 查询支付状态
func (s *PaymentService) GetPaymentStatus(ctx context.Context, paymentID uint64) (*PaymentStatusResponse, error) {
	payment, err := s.payments.Get(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	return &PaymentStatusResponse{
		PaymentID: payment.ID,
		OrderID:   payment.OrderID,
		Status:    payment.Status,
		PaidAt:    payment.PaidAt,
	}, nil
}

// CancelPayment 取消支付
func (s *PaymentService) CancelPayment(ctx context.Context, userID uint64, paymentID uint64) error {
	payment, err := s.payments.Get(ctx, paymentID)
	if err != nil {
		return err
	}

	// 权限检查
	if payment.UserID != userID {
		return errors.New("unauthorized")
	}

	// 状态检查：只有 pending 状态可以取消
	if payment.Status != model.PaymentStatusPending {
		return errors.New("cannot cancel payment")
	}

	// 更新支付状态
	payment.Status = model.PaymentStatusFailed

	return s.payments.Update(ctx, payment)
}

// generateMockPayInfo 生成 Mock 支付参数
func (s *PaymentService) generateMockPayInfo(paymentID uint64, method model.PaymentMethod, amountCents int64) map[string]interface{} {
	payInfo := map[string]interface{}{
		"paymentId":   paymentID,
		"method":      string(method),
		"amountCents": amountCents,
		"currency":    "CNY",
		"timestamp":   time.Now().Unix(),
	}

	// 根据支付方式生成不同的参数
	switch method {
	case model.PaymentMethodWeChat:
		payInfo["prepay_id"] = fmt.Sprintf("wx_mock_%d", paymentID)
		payInfo["code_url"] = fmt.Sprintf("weixin://wxpay/bizpayurl?pr=mock%d", paymentID)
	case model.PaymentMethodAlipay:
		payInfo["trade_no"] = fmt.Sprintf("ali_mock_%d", paymentID)
		payInfo["qr_code"] = fmt.Sprintf("https://qr.alipay.com/mock%d", paymentID)
	}

	return payInfo
}

// mockPaymentSuccess Mock 支付成功（仅用于测试）
func (s *PaymentService) mockPaymentSuccess(ctx context.Context, paymentID uint64, order *model.Order) error {
	// 更新支付状态
	payment, err := s.payments.Get(ctx, paymentID)
	if err != nil {
		return err
	}

	now := time.Now()
	payment.Status = model.PaymentStatusPaid
	payment.PaidAt = &now
	payment.ProviderTradeNo = fmt.Sprintf("mock_trade_%d", paymentID)

	if err := s.payments.Update(ctx, payment); err != nil {
		return err
	}

	// 更新订单状态
	order.Status = model.OrderStatusConfirmed
	if err := s.orders.Update(ctx, order); err != nil {
		return err
	}

	return nil
}

// HandlePaymentCallback 处理支付回调（预留接口）
func (s *PaymentService) HandlePaymentCallback(ctx context.Context, provider string, data map[string]interface{}) error {
	// TODO: 实现真实的支付回调处理逻辑
	// 1. 验证签名
	// 2. 验证金额
	// 3. 更新支付状态
	// 4. 更新订单状态
	return nil
}

// RefundPayment 退款（预留接口）
func (s *PaymentService) RefundPayment(ctx context.Context, paymentID uint64, reason string) error {
	// TODO: 实现真实的退款逻辑
	// 1. 调用支付提供商退款接口
	// 2. 更新支付状态
	// 3. 更新订单状态
	payment, err := s.payments.Get(ctx, paymentID)
	if err != nil {
		return err
	}

	// Mock: 直接标记为已退款
	now := time.Now()
	payment.Status = model.PaymentStatusRefunded
	payment.RefundedAt = &now

	return s.payments.Update(ctx, payment)
}
