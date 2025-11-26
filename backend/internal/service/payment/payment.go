package payment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/apierr"
	"gamelink/internal/cache"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	repoiface "gamelink/internal/repository/interfaces"
)

var (
	// ErrNotFound 支付记录不存在
	ErrNotFound = apierr.NotFound("payment not found")
	// ErrValidation 表示输入校验失败
	ErrValidation = apierr.BadRequest("validation failed")
	// ErrOrderAlreadyPaid 订单已支付
	ErrOrderAlreadyPaid = apierr.BadRequest("order already paid")
	// ErrInvalidOrderStatus 订单状态不正确
	ErrInvalidOrderStatus = apierr.BadRequest("invalid order status")
)

// PaymentService 支付服务
//
// 功能：
// 1. 创建支付
// 2. 查询支付状态
// 3. 取消支付
// 4. 处理支付回调（Mock版本）
type PaymentService struct {
	payments        repository.PaymentRepository
	orders          repoiface.OrderReadWriter
	providers       map[model.PaymentMethod]ProviderClient
	distributedLock cache.DistributedLock
	wallets         repository.WalletRepository // 分布式锁，用于并发控制
}

// NewPaymentService 创建支付服务
func NewPaymentService(
	payments repository.PaymentRepository,
	orders repoiface.OrderReadWriter,
) *PaymentService {
	return &PaymentService{
		payments: payments,
		orders:   orders,
		providers: map[model.PaymentMethod]ProviderClient{
			model.PaymentMethodWeChat: wechatProvider{},
			model.PaymentMethodAlipay: alipayProvider{},
		},
	}
}

// SetDistributedLock injects distributed lock for concurrency control
func (s *PaymentService) SetDistributedLock(lock cache.DistributedLock) {
	s.distributedLock = lock
}

// SetWalletRepository injects wallet repository for refund credit.
func (s *PaymentService) SetWalletRepository(repo repository.WalletRepository) {
	s.wallets = repo
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
	// 使用分布式锁确保幂等性
	lockKey := fmt.Sprintf("payment:create:order:%d:user:%d", req.OrderID, userID)

	// 如果分布式锁可用，使用它
	if s.distributedLock != nil {
		locked, err := s.distributedLock.TryLock(ctx, lockKey, time.Second*10, 3, time.Millisecond*50)
		if err != nil {
			return nil, apierr.InternalError("failed to acquire lock").WithDetails(err.Error())
		}
		if !locked {
			return nil, apierr.Conflict("payment creation in progress, please try again")
		}
		defer s.distributedLock.Unlock(ctx, lockKey)
	}

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

	// 检查是否已有支付记录（幂等性检查）
	orderIDPtr := &req.OrderID
	existingPayments, _, err := s.payments.List(ctx, repository.PaymentListOptions{
		OrderID:  orderIDPtr,
		Page:     1,
		PageSize: 10, // 检查更多记录以确保找到pending状态的支付
	})
	if err == nil && len(existingPayments) > 0 {
		// 检查是否已支付
		for _, existingPayment := range existingPayments {
			if existingPayment.Status == model.PaymentStatusPaid {
				return nil, ErrOrderAlreadyPaid
			}
			// 如果存在pending状态的支付，返回已存在
			if existingPayment.Status == model.PaymentStatusPending {
				return nil, apierr.Conflict("payment already exists for this order").WithDetails("Please check your payment status")
			}
		}
	}

	// 创建支付记录
	payment := &model.Payment{
		OrderID:     req.OrderID,
		UserID:      userID,
		Method:      req.Method,
		AmountCents: order.TotalPriceCents,
		Currency:    order.Currency,
		Status:      model.PaymentStatusPending,
	}

	if err := s.payments.Create(ctx, payment); err != nil {
		return nil, err
	}

	// Mock: 生成支付参数
	payInfo := s.generateMockPayInfo(payment.ID, req.Method, order.TotalPriceCents)

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
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
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
		if errors.Is(err, repository.ErrNotFound) {
			return repository.ErrNotFound
		}
		// 对于取消支付操作，其他获取失败视为内部错误（500）
		return apierr.InternalError("failed to get payment").WithDetails(err.Error())
	}

	// 权限检查
	if payment.UserID != userID {
		return apierr.Forbidden("unauthorized")
	}

	// 状态检查：只有 pending 状态可以取消
	if payment.Status != model.PaymentStatusPending {
		return apierr.BadRequest("cannot cancel payment")
	}

	// 更新支付状态
	payment.Status = model.PaymentStatusFailed

	if err := s.payments.Update(ctx, payment); err != nil {
		return WrapError(err, "update payment")
	}

	return nil
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

// HandlePaymentCallback 处理支付回调
//
// 注意：这是一个简化版本，生产环境需要：
// 1. 验证支付提供商签名
// 2. 防止重复回调
// 3. 使用事务确保数据一致性
func (s *PaymentService) HandlePaymentCallback(ctx context.Context, provider string, data map[string]interface{}) error {
	// 获取支付ID
	paymentID, ok := data["payment_id"].(uint64)
	if !ok {
		// 尝试从其他字段获取
		if id, ok := data["payment_id"].(float64); ok {
			paymentID = uint64(id)
		} else {
			return errors.New("missing payment_id in callback data")
		}
	}

	// 获取支付记录
	payment, err := s.payments.Get(ctx, paymentID)
	if err != nil {
		return err
	}

	// 验证支付状态：只处理待支付的订单
	if payment.Status != model.PaymentStatusPending {
		// 已经处理过，返回成功避免重复处理
		return nil
	}

	// 验证支付提供商
	expectedProvider := string(payment.Method)
	if provider != expectedProvider {
		return fmt.Errorf("provider mismatch: expected %s, got %s", expectedProvider, provider)
	}

	// 验证金额
	if amountCents, ok := data["amount_cents"].(int64); ok {
		if amountCents != payment.AmountCents {
			return fmt.Errorf("amount mismatch: expected %d, got %d", payment.AmountCents, amountCents)
		}
	}

	// 获取订单
	order, err := s.orders.Get(ctx, payment.OrderID)
	if err != nil {
		return err
	}

	// 更新支付状态
	now := time.Now()
	payment.Status = model.PaymentStatusPaid
	payment.PaidAt = &now

	// 设置第三方交易号
	if tradeNo, ok := data["trade_no"].(string); ok {
		payment.ProviderTradeNo = tradeNo
	} else {
		payment.ProviderTradeNo = fmt.Sprintf("%s_%d_%d", provider, paymentID, now.Unix())
	}

	if err := s.payments.Update(ctx, payment); err != nil {
		return err
	}

	// 更新订单状态为已确认
	order.Status = model.OrderStatusConfirmed
	if err := s.orders.Update(ctx, order); err != nil {
		return err
	}

	return nil
}

// RefundPayment 退款
//
// 注意：这是一个简化版本，生产环境需要：
// 1. 调用真实的支付提供商退款API
// 2. 处理部分退款
// 3. 处理退款失败重试
func (s *PaymentService) RefundPayment(ctx context.Context, paymentID uint64, reason string) error {
	payment, err := s.payments.Get(ctx, paymentID)
	if err != nil {
		return err
	}

	// 验证支付状态：只有已支付的订单可以退款
	if payment.Status != model.PaymentStatusPaid {
		return fmt.Errorf("payment status must be paid, current: %s", payment.Status)
	}

	client, ok := s.providers[payment.Method]
	if !ok {
		client = genericProvider{}
	}
	tradeNo, raw, refundedAt, err := client.Refund(ctx, payment, reason)
	if err != nil {
		return err
	}

	// 更新支付状态
	payment.Status = model.PaymentStatusRefunded
	payment.RefundedAt = &refundedAt
	payment.ProviderTradeNo = tradeNo
	payment.ProviderRaw = raw

	if err := s.payments.Update(ctx, payment); err != nil {
		return err
	}

	// 更新订单状态
	order, err := s.orders.Get(ctx, payment.OrderID)
	if err != nil {
		return err
	}

	order.Status = model.OrderStatusRefunded
	order.RefundAmountCents = payment.AmountCents
	order.RefundReason = reason
	order.RefundedAt = &refundedAt

	if err := s.orders.Update(ctx, order); err != nil {
		return err
	}

	// 钱包回滚：如配置了钱包仓储，将退款金额计入余额
	if s.wallets != nil {
		if err := s.creditWallet(ctx, payment.UserID, payment.AmountCents); err != nil {
			return err
		}
	}

	return nil
}

// List 返回支付列表
func (s *PaymentService) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize <= 0 || opts.PageSize > 100 {
		opts.PageSize = 20
	}
	return s.payments.List(ctx, opts)
}

func (s *PaymentService) creditWallet(ctx context.Context, userID uint64, amount int64) error {
	w, err := s.wallets.GetByUserID(ctx, userID)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFound) {
			return err
		}
		w = &model.Wallet{UserID: userID}
	}
	w.BalanceCents += amount
	return s.wallets.Save(ctx, w)
}
