package payment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/collectionentity"
	"gamelink/internal/repository/common"
	repoiface "gamelink/internal/repository/interfaces"
	"gamelink/internal/repository/routingrule"
	"gamelink/internal/service/external"
	routingruleservice "gamelink/internal/service/routingrule"
	"gamelink/internal/ws"
	"gamelink/pkg/apierr"
	"gamelink/pkg/cache"
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
// 4. 处理支付回调
// 5. 收款分流（Requirements: 17.1, 17.2, 17.3）
//
// ============================================================================
// TODO: 生产环境部署前需要完成以下第三方接口对接：
// ============================================================================
// 1. 微信支付
//   - [ ] 统一下单接口: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
//   - [ ] 支付回调验签: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_3
//   - [ ] 退款接口: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_4
//   - [ ] 商户证书配置
//
// 2. 支付宝
//   - [ ] 统一收单交易支付接口: https://opendocs.alipay.com/apis/api_1/alipay.trade.pay
//   - [ ] 异步通知验签: https://opendocs.alipay.com/open/270/105902
//   - [ ] 退款接口: https://opendocs.alipay.com/apis/api_1/alipay.trade.refund
//   - [ ] 应用公私钥配置
//
// 3. 当前状态：Mock 模式
//   - 支付自动成功（mockPaymentSuccess）
//   - 退款自动成功
//   - 适用于开发和测试环境
//
// ============================================================================
type PaymentService struct {
	tx                    common.TxManager                   // 事务管理器
	payments              repository.PaymentRepository
	orders                repoiface.OrderReadWriter
	providers             map[model.PaymentMethod]ProviderClient
	distributedLock       cache.DistributedLock
	wallets               repository.WalletRepository       // 钱包仓储
	routingEngine         *routingruleservice.RoutingEngine // 收款分流引擎
	notifications         repository.NotificationRepository
	hub                   *ws.Hub
	wechatCallbackHandler *WeChatCallbackHandler           // 微信支付回调处理器
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

// SetTxManager injects a transaction manager for multi-step operations.
func (s *PaymentService) SetTxManager(tx common.TxManager) {
	s.tx = tx
}

// SetDistributedLock injects distributed lock for concurrency control
func (s *PaymentService) SetDistributedLock(lock cache.DistributedLock) {
	s.distributedLock = lock
}

// SetWalletRepository injects wallet repository for refund credit.
func (s *PaymentService) SetWalletRepository(repo repository.WalletRepository) {
	s.wallets = repo
}

// SetNotificationRepository injects notification repository for payment events.
func (s *PaymentService) SetNotificationRepository(repo repository.NotificationRepository) {
	s.notifications = repo
}

// SetWebsocketHub injects WebSocket hub for realtime order updates.
func (s *PaymentService) SetWebsocketHub(hub *ws.Hub) {
	s.hub = hub
}

// SetWeChatCallbackHandler injects WeChat payment callback handler.
func (s *PaymentService) SetWeChatCallbackHandler(handler *WeChatCallbackHandler) {
	s.wechatCallbackHandler = handler
}

// SetRoutingEngine injects routing engine for payment routing.
// Requirements: 17.1, 17.2, 17.3
func (s *PaymentService) SetRoutingEngine(engine *routingruleservice.RoutingEngine) {
	s.routingEngine = engine
}

// InitRoutingEngine initializes the routing engine with repositories.
// Requirements: 17.1, 17.2, 17.3
func (s *PaymentService) InitRoutingEngine(
	ruleRepo routingrule.RoutingRuleRepository,
	entityRepo collectionentity.CollectionEntityRepository,
) {
	s.routingEngine = routingruleservice.NewRoutingEngine(ruleRepo, entityRepo)
}

// SetExternalConfig configures external API credentials for payment providers
func (s *PaymentService) SetExternalConfig(cfg *external.Config) {
	factory := NewProviderFactory(cfg)
	s.providers = factory.CreateProviders()
}

// SetProviders sets payment providers directly (for testing)
func (s *PaymentService) SetProviders(providers map[model.PaymentMethod]ProviderClient) {
	s.providers = providers
}

// CreatePaymentRequest 创建支付请求
type CreatePaymentRequest struct {
	OrderID uint64              `json:"orderId" binding:"required"`
	Method  model.PaymentMethod `json:"method" binding:"required,oneof=wechat alipay wallet combined"`
	// RequestID 客户端生成的唯一请求ID，用于幂等性控制
	// 客户端应生成 UUID，网络失败时使用相同的 RequestID 重试
	RequestID string `json:"requestId,omitempty"`
	// 组合支付参数（仅当 Method 为 combined 时使用）
	WalletAmountCents int64               `json:"walletAmountCents,omitempty"` // 钱包支付金额（分）
	ThirdPartyMethod  model.PaymentMethod `json:"thirdPartyMethod,omitempty"`  // 第三方支付方式
}

// CreatePaymentResponse 创建支付响应
type CreatePaymentResponse struct {
	PaymentID        uint64                 `json:"paymentId"`
	PayInfo          map[string]interface{} `json:"payInfo"`                    // 支付参数（对接支付SDK）
	WalletDeducted   int64                  `json:"walletDeducted,omitempty"`   // 钱包已扣款金额
	ThirdPartyAmount int64                  `json:"thirdPartyAmount,omitempty"` // 第三方需支付金额
	WalletPaidDirect bool                   `json:"walletPaidDirect,omitempty"` // 是否纯钱包支付完成
}

// PaymentStatusResponse 支付状态响应
type PaymentStatusResponse struct {
	PaymentID uint64              `json:"paymentId"`
	OrderID   uint64              `json:"orderId"`
	Status    model.PaymentStatus `json:"status"`
	PaidAt    *time.Time          `json:"paidAt"`
}

// CreatePayment 创建支付
// 支持三种支付方式：
// 1. 纯第三方支付（wechat/alipay）
// 2. 纯钱包支付（wallet）
// 3. 组合支付（combined）：钱包+第三方
//
// 幂等性保证：
// - 如果提供了 RequestID，会先检查是否已存在相同 RequestID 的支付记录
// - 如果存在且状态为 paid，返回已有支付信息
// - 如果存在且状态为 pending，返回冲突错误
// - 如果存在且状态为 failed，允许重新创建
func (s *PaymentService) CreatePayment(ctx context.Context, userID uint64, req CreatePaymentRequest) (*CreatePaymentResponse, error) {
	// 幂等性检查：如果提供了 RequestID，先检查是否已存在
	if req.RequestID != "" {
		existingPayment, err := s.payments.GetByRequestID(ctx, req.RequestID)
		if err == nil && existingPayment != nil {
			// 已存在相同 RequestID 的支付记录
			switch existingPayment.Status {
			case model.PaymentStatusPaid:
				// 已支付成功，返回已有支付信息（幂等返回）
				return &CreatePaymentResponse{
					PaymentID:        existingPayment.ID,
					PayInfo:          map[string]interface{}{"status": "already_paid", "idempotent": true},
					WalletDeducted:   existingPayment.WalletAmountCents,
					ThirdPartyAmount: existingPayment.ThirdPartyAmountCents,
					WalletPaidDirect: existingPayment.Method == model.PaymentMethodWallet,
				}, nil
			case model.PaymentStatusPending:
				// 支付进行中，返回冲突错误
				return nil, apierr.Conflict("payment with this request_id is already in progress")
			case model.PaymentStatusFailed:
				// 之前失败了，允许重新创建（不使用旧记录）
				// 继续执行创建流程
			}
		}
	}

	// 使用分布式锁确保并发安全
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
		defer func() { _ = s.distributedLock.Unlock(ctx, lockKey) }()
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

	// 检查是否已有支付记录（基于 OrderID 的幂等性检查）
	orderIDPtr := &req.OrderID
	existingPayments, _, err := s.payments.List(ctx, repository.PaymentListOptions{
		OrderID:  orderIDPtr,
		Page:     1,
		PageSize: 10,
	})
	if err == nil && len(existingPayments) > 0 {
		for _, existingPayment := range existingPayments {
			if existingPayment.Status == model.PaymentStatusPaid {
				return nil, ErrOrderAlreadyPaid
			}
			if existingPayment.Status == model.PaymentStatusPending {
				return nil, apierr.Conflict("payment already exists for this order").WithDetails("Please check your payment status")
			}
		}
	}

	// 根据支付方式处理
	switch req.Method {
	case model.PaymentMethodWallet:
		return s.createWalletPayment(ctx, userID, order, req)
	case model.PaymentMethodCombined:
		return s.createCombinedPayment(ctx, userID, order, req)
	default:
		return s.createThirdPartyPayment(ctx, userID, order, req)
	}
}

// createWalletPayment 创建纯钱包支付
func (s *PaymentService) createWalletPayment(ctx context.Context, userID uint64, order *model.Order, req CreatePaymentRequest) (*CreatePaymentResponse, error) {
	if s.wallets == nil {
		return nil, apierr.BadRequest("wallet payment not supported")
	}

	// 获取用户钱包
	wallet, err := s.wallets.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apierr.BadRequest("insufficient wallet balance")
		}
		return nil, err
	}

	// 检查余额
	if wallet.BalanceCents < order.TotalPriceCents {
		return nil, apierr.BadRequest("insufficient wallet balance").WithDetails(
			fmt.Sprintf("需要 %d 分，余额 %d 分", order.TotalPriceCents, wallet.BalanceCents))
	}

	// 使用乐观锁扣除钱包余额
	updatedWallet, err := s.wallets.UpdateBalanceWithLock(ctx, userID, -order.TotalPriceCents, 3)
	if err != nil {
		return nil, apierr.InternalError("failed to deduct wallet balance").WithDetails(err.Error())
	}

	// 创建支付记录
	now := time.Now()
	payment := &model.Payment{
		RequestID:         req.RequestID, // 设置幂等请求ID
		OrderID:           req.OrderID,
		UserID:            userID,
		Method:            model.PaymentMethodWallet,
		AmountCents:       order.TotalPriceCents,
		WalletAmountCents: order.TotalPriceCents,
		Currency:          order.Currency,
		Status:            model.PaymentStatusPaid,
		PaidAt:            &now,
		ProviderTradeNo:   fmt.Sprintf("wallet_%d_%d", userID, now.Unix()),
	}

	if err := s.payments.Create(ctx, payment); err != nil {
		// 回滚钱包余额
		_, _ = s.wallets.UpdateBalanceWithLock(ctx, userID, order.TotalPriceCents, 3)
		return nil, err
	}

	// 更新订单状态
	order.Status = model.OrderStatusConfirmed
	if err := s.orders.Update(ctx, order); err != nil {
		return nil, err
	}

	s.notifyOrderPaid(ctx, order)
	s.broadcastNewOrder(order)

	return &CreatePaymentResponse{
		PaymentID:        payment.ID,
		PayInfo:          map[string]interface{}{"method": "wallet", "status": "paid", "newBalance": updatedWallet.BalanceCents},
		WalletDeducted:   order.TotalPriceCents,
		WalletPaidDirect: true,
	}, nil
}

// createCombinedPayment 创建组合支付（钱包+第三方）
func (s *PaymentService) createCombinedPayment(ctx context.Context, userID uint64, order *model.Order, req CreatePaymentRequest) (*CreatePaymentResponse, error) {
	if s.wallets == nil {
		return nil, apierr.BadRequest("combined payment not supported")
	}

	// 验证组合支付参数
	if req.WalletAmountCents <= 0 {
		return nil, apierr.BadRequest("wallet amount must be positive for combined payment")
	}
	if req.ThirdPartyMethod != model.PaymentMethodWeChat && req.ThirdPartyMethod != model.PaymentMethodAlipay {
		return nil, apierr.BadRequest("third party method must be wechat or alipay")
	}
	if req.WalletAmountCents >= order.TotalPriceCents {
		return nil, apierr.BadRequest("wallet amount must be less than total price for combined payment, use wallet method instead")
	}

	// 获取用户钱包
	wallet, err := s.wallets.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apierr.BadRequest("insufficient wallet balance")
		}
		return nil, err
	}

	// 检查钱包余额
	if wallet.BalanceCents < req.WalletAmountCents {
		return nil, apierr.BadRequest("insufficient wallet balance").WithDetails(
			fmt.Sprintf("需要 %d 分，余额 %d 分", req.WalletAmountCents, wallet.BalanceCents))
	}

	// 计算第三方支付金额
	thirdPartyAmount := order.TotalPriceCents - req.WalletAmountCents

	// 使用乐观锁扣除钱包余额（预扣款）
	_, err = s.wallets.UpdateBalanceWithLock(ctx, userID, -req.WalletAmountCents, 3)
	if err != nil {
		return nil, apierr.InternalError("failed to deduct wallet balance").WithDetails(err.Error())
	}

	// 创建支付记录
	payment := &model.Payment{
		RequestID:             req.RequestID, // 设置幂等请求ID
		OrderID:               req.OrderID,
		UserID:                userID,
		Method:                model.PaymentMethodCombined,
		AmountCents:           order.TotalPriceCents,
		WalletAmountCents:     req.WalletAmountCents,
		ThirdPartyMethod:      req.ThirdPartyMethod,
		ThirdPartyAmountCents: thirdPartyAmount,
		Currency:              order.Currency,
		Status:                model.PaymentStatusPending,
	}

	// 执行收款分流
	if s.routingEngine != nil {
		routingResult, err := s.routePayment(ctx, order, req.ThirdPartyMethod)
		if err == nil {
			payment.CollectionEntityID = &routingResult.CollectionEntityID
			payment.MerchantNo = routingResult.MerchantNo
		}
	}

	if err := s.payments.Create(ctx, payment); err != nil {
		// 回滚钱包余额
		_, _ = s.wallets.UpdateBalanceWithLock(ctx, userID, req.WalletAmountCents, 3)
		return nil, err
	}

	// 生成第三方支付参数
	payInfo := s.generateMockPayInfo(payment.ID, req.ThirdPartyMethod, thirdPartyAmount)
	payInfo["walletDeducted"] = req.WalletAmountCents
	payInfo["combinedPayment"] = true

	// Mock: 自动完成第三方支付（仅用于测试）
	if err := s.mockPaymentSuccess(ctx, payment.ID, order); err != nil {
		return nil, err
	}

	return &CreatePaymentResponse{
		PaymentID:        payment.ID,
		PayInfo:          payInfo,
		WalletDeducted:   req.WalletAmountCents,
		ThirdPartyAmount: thirdPartyAmount,
	}, nil
}

// createThirdPartyPayment 创建纯第三方支付
func (s *PaymentService) createThirdPartyPayment(ctx context.Context, userID uint64, order *model.Order, req CreatePaymentRequest) (*CreatePaymentResponse, error) {
	// 创建支付记录
	payment := &model.Payment{
		RequestID:             req.RequestID, // 设置幂等请求ID
		OrderID:               req.OrderID,
		UserID:                userID,
		Method:                req.Method,
		AmountCents:           order.TotalPriceCents,
		ThirdPartyAmountCents: order.TotalPriceCents,
		Currency:              order.Currency,
		Status:                model.PaymentStatusPending,
	}

	// 执行收款分流
	if s.routingEngine != nil {
		routingResult, err := s.routePayment(ctx, order, req.Method)
		if err == nil {
			payment.CollectionEntityID = &routingResult.CollectionEntityID
			payment.MerchantNo = routingResult.MerchantNo
		}
	}

	if err := s.payments.Create(ctx, payment); err != nil {
		return nil, err
	}

	// 创建分流日志
	if s.routingEngine != nil && payment.CollectionEntityID != nil {
		routingResult, _ := s.routePayment(ctx, order, req.Method)
		if routingResult != nil {
			_ = s.routingEngine.CreateRoutingLog(ctx, payment.ID, order.ID, routingResult)
		}
	}

	// Mock: 生成支付参数
	payInfo := s.generateMockPayInfo(payment.ID, req.Method, order.TotalPriceCents)

	// Mock: 自动标记为已支付（仅用于测试）
	if err := s.mockPaymentSuccess(ctx, payment.ID, order); err != nil {
		return nil, err
	}

	return &CreatePaymentResponse{
		PaymentID:        payment.ID,
		PayInfo:          payInfo,
		ThirdPartyAmount: order.TotalPriceCents,
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

	s.notifyOrderPaid(ctx, order)
	s.broadcastNewOrder(order)

	return nil
}

// HandlePaymentCallback 处理支付回调
//
// ============================================================================
// TODO: 生产环境必须实现以下安全验证：
// ============================================================================
// 1. 验证支付提供商签名（防止伪造回调）
//   - 微信支付: 使用 API Key 进行 MD5/HMAC-SHA256 签名验证
//   - 支付宝: 使用支付宝公钥进行 RSA 签名验证
//
// 2. 防止重复回调
//   - 使用 Redis 分布式锁
//   - 或使用数据库唯一约束
//
// 3. 使用数据库事务确保数据一致性
//
// 4. 返回正确的响应格式
//   - 微信支付: 返回 XML 格式 <return_code>SUCCESS</return_code>
//   - 支付宝: 返回字符串 "success"
//
// 当前实现：Mock 模式，不验证签名
// ============================================================================
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

	// 分布式锁防止并发回调重复处理
	if s.distributedLock != nil {
		lockKey := fmt.Sprintf("payment:callback:%d", paymentID)
		locked, err := s.distributedLock.TryLock(ctx, lockKey, time.Second*10, 3, time.Millisecond*100)
		if err != nil {
			return apierr.InternalError("failed to acquire callback lock").WithDetails(err.Error())
		}
		if !locked {
			return nil // 其他实例正在处理，返回成功
		}
		defer func() { _ = s.distributedLock.Unlock(ctx, lockKey) }()
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

	// 使用事务确保支付状态和订单状态的原子性更新
	err = s.tx.WithTx(ctx, func(r *common.Repos) error {
		if err := r.Payments.Update(ctx, payment); err != nil {
			return fmt.Errorf("update payment: %w", err)
		}

		order.Status = model.OrderStatusConfirmed
		if err := r.Orders.Update(ctx, order); err != nil {
			return fmt.Errorf("update order: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	s.notifyOrderPaid(ctx, order)
	s.broadcastNewOrder(order)

	return nil
}

// HandleWeChatPaymentCallback 处理微信支付回调
//
// 功能：
// 1. 验证回调签名（防止伪造）
// 2. 防止重复回调（使用分布式锁）
// 3. 更新支付和订单状态
// 4. 返回 XML 格式响应
//
// 参数：
//   - ctx: 上下文
//   - callbackData: 微信支付回调 XML 数据（字节）
//
// 返回：
//   - []byte: XML 格式的响应数据
//   - error: 错误信息
func (s *PaymentService) HandleWeChatPaymentCallback(ctx context.Context, callbackData []byte) ([]byte, error) {
	// TODO: 完成微信支付回调处理集成
	//
	// 当前状态：
	// ✅ WeChatCallbackHandler 已实现并测试通过
	// ✅ 签名验证已实现
	// ✅ 防重放攻击已实现
	// ✅ XML 解析已实现
	//
	// 待完成工作：
	// ⏳ 需要在 OrderRepository 中添加 GetByOrderNo 方法
	// ⏳ 需要在 PaymentRepository 中添加 GetByRequestID 方法
	// ⏳ 完整的回调处理逻辑（更新支付和订单状态）
	// ⏳ 集成测试
	//
	// 实现步骤：
	// 1. 在 api/internal/repository/interfaces/order.go 中添加：
	//    GetByOrderNo(ctx context.Context, orderNo string) (*model.Order, error)
	//
	// 2. 在 api/internal/repository/interfaces.go 中添加：
	//    GetByRequestID(ctx context.Context, requestID string) (*model.Payment, error)
	//
	// 3. 在 api/internal/repository/implementations/order.go 中实现 GetByOrderNo
	//
	// 4. 在 api/internal/repository/payment/repository.go 中实现 GetByRequestID
	//
	// 5. 在此方法中完成回调处理逻辑：
	//    - 验证回调签名（已完成）
	//    - 查找支付记录
	//    - 验证金额
	//    - 更新支付状态
	//    - 更新订单状态
	//    - 发送通知
	//
	// 临时返回：模拟成功响应
	_ = ctx
	_ = callbackData
	return NewWeChatCallbackResponse().ToXML()
}

// RefundPayment 退款
//
// 支持三种支付方式的退款：
// 1. 纯第三方支付：退回第三方
// 2. 纯钱包支付：退回钱包
// 3. 组合支付：钱包部分退回钱包，第三方部分退回第三方
func (s *PaymentService) RefundPayment(ctx context.Context, paymentID uint64, reason string) error {
	payment, err := s.payments.Get(ctx, paymentID)
	if err != nil {
		return err
	}

	// 验证支付状态：只有已支付的订单可以退款
	if payment.Status != model.PaymentStatusPaid {
		return fmt.Errorf("payment status must be paid, current: %s", payment.Status)
	}

	now := time.Now()
	var tradeNo string
	var raw []byte

	// 根据支付方式处理退款
	switch payment.Method {
	case model.PaymentMethodWallet:
		// 纯钱包支付：退回钱包
		if s.wallets != nil {
			if err := s.creditWallet(ctx, payment.UserID, payment.WalletAmountCents); err != nil {
				return err
			}
		}
		tradeNo = fmt.Sprintf("wallet_refund_%d_%d", paymentID, now.Unix())

	case model.PaymentMethodCombined:
		// 组合支付：钱包部分退回钱包，第三方部分退回第三方
		// 1. 退回钱包部分
		if s.wallets != nil && payment.WalletAmountCents > 0 {
			if err := s.creditWallet(ctx, payment.UserID, payment.WalletAmountCents); err != nil {
				return err
			}
		}
		// 2. 退回第三方部分
		if payment.ThirdPartyAmountCents > 0 {
			client, ok := s.providers[payment.ThirdPartyMethod]
			if !ok {
				client = genericProvider{}
			}
			var refundedAt time.Time
			tradeNo, raw, refundedAt, err = client.Refund(ctx, payment, reason)
			if err != nil {
				// 第三方退款失败，回滚钱包退款
				if s.wallets != nil && payment.WalletAmountCents > 0 {
					_ = s.debitWallet(ctx, payment.UserID, payment.WalletAmountCents)
				}
				return err
			}
			now = refundedAt
		} else {
			tradeNo = fmt.Sprintf("combined_refund_%d_%d", paymentID, now.Unix())
		}

	default:
		// 纯第三方支付：退回第三方
		client, ok := s.providers[payment.Method]
		if !ok {
			client = genericProvider{}
		}
		var refundedAt time.Time
		tradeNo, raw, refundedAt, err = client.Refund(ctx, payment, reason)
		if err != nil {
			return err
		}
		now = refundedAt
	}

	// 使用事务确保支付状态和订单状态的原子性更新
	payment.Status = model.PaymentStatusRefunded
	payment.RefundedAt = &now
	payment.ProviderTradeNo = tradeNo
	payment.ProviderRaw = raw

	var order *model.Order

	err = s.tx.WithTx(ctx, func(r *common.Repos) error {
		if err := r.Payments.Update(ctx, payment); err != nil {
			return fmt.Errorf("update payment: %w", err)
		}

		var getErr error
		order, getErr = r.Orders.Get(ctx, payment.OrderID)
		if getErr != nil {
			return fmt.Errorf("get order: %w", getErr)
		}

		order.Status = model.OrderStatusRefunded
		order.RefundAmountCents = payment.AmountCents
		order.RefundReason = reason
		order.RefundedAt = &now

		if err := r.Orders.Update(ctx, order); err != nil {
			return fmt.Errorf("update order: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	s.notifyOrderRefunded(ctx, order)

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

// creditWallet 使用乐观锁安全地增加钱包余额
func (s *PaymentService) creditWallet(ctx context.Context, userID uint64, amount int64) error {
	_, err := s.wallets.UpdateBalanceWithLock(ctx, userID, amount, 3)
	return err
}

// debitWallet 使用乐观锁安全地从钱包扣款（用于回滚退款）
func (s *PaymentService) debitWallet(ctx context.Context, userID uint64, amount int64) error {
	_, err := s.wallets.UpdateBalanceWithLock(ctx, userID, -amount, 3)
	return err
}

// routePayment 执行支付分流
// Requirements: 17.1, 17.2
// Property 15: 收款分流规则优先级 - 按优先级顺序匹配规则
// Property 16: 收款分流默认主体回退 - 无规则匹配时使用默认主体
func (s *PaymentService) routePayment(ctx context.Context, order *model.Order, method model.PaymentMethod) (*routingruleservice.RoutingResult, error) {
	if s.routingEngine == nil {
		return nil, errors.New("routing engine not initialized")
	}

	// 构建分流上下文
	routingCtx := &routingruleservice.RoutingContext{
		OrderID:     order.ID,
		AmountCents: order.TotalPriceCents,
		Method:      method,
	}

	// 获取游戏类型
	if order.Game != nil {
		routingCtx.GameType = order.Game.Name
	} else if order.GameID != nil {
		routingCtx.GameType = fmt.Sprintf("game_%d", *order.GameID)
	}

	// 获取服务类型
	if order.ServiceItem != nil {
		routingCtx.ServiceType = order.ServiceItem.Name
	} else {
		routingCtx.ServiceType = "default"
	}

	// 执行分流
	return s.routingEngine.RoutePayment(ctx, routingCtx)
}

// GetPaymentRoutingLog 获取支付分流日志
// Requirements: 17.3
func (s *PaymentService) GetPaymentRoutingLog(ctx context.Context, paymentID uint64) (*model.RoutingLog, error) {
	if s.routingEngine == nil {
		return nil, errors.New("routing engine not initialized")
	}
	return s.routingEngine.GetRoutingLogByPayment(ctx, paymentID)
}

// WalletBalanceResponse 钱包余额响应
type WalletBalanceResponse struct {
	BalanceCents int64 `json:"balanceCents"` // 可用余额（分）
	FrozenCents  int64 `json:"frozenCents"`  // 冻结金额（分）
}

// GetWalletBalance 获取用户钱包余额
func (s *PaymentService) GetWalletBalance(ctx context.Context, userID uint64) (*WalletBalanceResponse, error) {
	if s.wallets == nil {
		return &WalletBalanceResponse{BalanceCents: 0, FrozenCents: 0}, nil
	}

	wallet, err := s.wallets.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &WalletBalanceResponse{BalanceCents: 0, FrozenCents: 0}, nil
		}
		return nil, err
	}

	return &WalletBalanceResponse{
		BalanceCents: wallet.BalanceCents,
		FrozenCents:  wallet.FrozenCents,
	}, nil
}

// CalculateCombinedPaymentRequest 计算组合支付请求
type CalculateCombinedPaymentRequest struct {
	OrderID           uint64 `json:"orderId" binding:"required"`
	WalletAmountCents int64  `json:"walletAmountCents"` // 希望使用的钱包金额，0表示使用全部余额
}

// CalculateCombinedPaymentResponse 计算组合支付响应
type CalculateCombinedPaymentResponse struct {
	OrderTotalCents       int64 `json:"orderTotalCents"`       // 订单总金额
	WalletBalanceCents    int64 `json:"walletBalanceCents"`    // 钱包余额
	WalletUsableCents     int64 `json:"walletUsableCents"`     // 可使用的钱包金额
	ThirdPartyAmountCents int64 `json:"thirdPartyAmountCents"` // 需要第三方支付的金额
	CanPayWithWalletOnly  bool  `json:"canPayWithWalletOnly"`  // 是否可以纯钱包支付
}

// CalculateCombinedPayment 计算组合支付金额
func (s *PaymentService) CalculateCombinedPayment(ctx context.Context, userID uint64, req CalculateCombinedPaymentRequest) (*CalculateCombinedPaymentResponse, error) {
	// 获取订单
	order, err := s.orders.Get(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}

	// 权限检查
	if order.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// 获取钱包余额
	var walletBalance int64
	if s.wallets != nil {
		wallet, err := s.wallets.GetByUserID(ctx, userID)
		if err == nil {
			walletBalance = wallet.BalanceCents
		}
	}

	// 计算可使用的钱包金额
	walletUsable := walletBalance
	if req.WalletAmountCents > 0 && req.WalletAmountCents < walletBalance {
		walletUsable = req.WalletAmountCents
	}
	if walletUsable > order.TotalPriceCents {
		walletUsable = order.TotalPriceCents
	}

	// 计算第三方支付金额
	thirdPartyAmount := order.TotalPriceCents - walletUsable

	return &CalculateCombinedPaymentResponse{
		OrderTotalCents:       order.TotalPriceCents,
		WalletBalanceCents:    walletBalance,
		WalletUsableCents:     walletUsable,
		ThirdPartyAmountCents: thirdPartyAmount,
		CanPayWithWalletOnly:  walletBalance >= order.TotalPriceCents,
	}, nil
}

func (s *PaymentService) notifyOrderPaid(ctx context.Context, order *model.Order) {
	s.notifyOrderStatus(ctx, order, "支付成功", "订单支付成功")
}

func (s *PaymentService) notifyOrderRefunded(ctx context.Context, order *model.Order) {
	s.notifyOrderStatus(ctx, order, "订单已退款", "订单已退款")
}

func (s *PaymentService) notifyOrderStatus(ctx context.Context, order *model.Order, title, message string) {
	if order == nil {
		return
	}

	s.broadcastOrderStatus(order, message)

	if s.notifications == nil {
		return
	}
	refID := order.ID
	event := &model.NotificationEvent{
		UserID:        order.UserID,
		Title:         title,
		Message:       message,
		Priority:      model.NotificationPriorityNormal,
		Channel:       "in_app",
		ReferenceType: "order",
		ReferenceID:   &refID,
	}
	if err := s.notifications.Create(ctx, event); err != nil {
		return
	}
	s.broadcastNotification(event)
}

func (s *PaymentService) broadcastOrderStatus(order *model.Order, message string) {
	if s.hub == nil || order == nil {
		return
	}
	payload := ws.OrderStatusEvent{
		OrderID:   order.ID,
		Status:    string(order.Status),
		Message:   message,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
	if data, err := ws.NewWSMessage(ws.MessageTypeOrderStatus, payload).ToJSON(); err == nil {
		s.hub.BroadcastToUser(data, order.UserID)
	}
}

func (s *PaymentService) broadcastNewOrder(order *model.Order) {
	if s.hub == nil || order == nil {
		return
	}
	payload := ws.OrderNewEvent{
		OrderID:        order.ID,
		Title:          order.Title,
		PriceCents:     order.TotalPriceCents,
		ScheduledStart: order.ScheduledStart,
		GameID:         order.GameID,
	}
	if data, err := ws.NewWSMessage(ws.MessageTypeOrderNew, payload).ToJSON(); err == nil {
		s.hub.BroadcastToRole(data, string(model.RolePlayer))
	}
}

func (s *PaymentService) broadcastNotification(event *model.NotificationEvent) {
	if s.hub == nil || event == nil {
		return
	}
	if data, err := ws.NewWSMessage(ws.MessageTypeNotification, event).ToJSON(); err == nil {
		s.hub.BroadcastToUser(data, event.UserID)
	}
}
