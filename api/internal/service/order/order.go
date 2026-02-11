package order

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionrepo "gamelink/internal/repository/commission"
	orderrepo "gamelink/internal/repository/implementations"
	repoiface "gamelink/internal/repository/interfaces"
	"gamelink/internal/repository/ordergroup"
	"gamelink/internal/ws"
	"gamelink/pkg/apierr"
	"gamelink/pkg/cache"
)

var (
	// ErrNotFound 订单不存在
	ErrNotFound = repository.ErrNotFound
	// ErrValidation 表示输入校验失败
	ErrValidation = apierr.BadRequest("validation failed")
	// ErrInvalidTransition 订单状态流转不合法
	ErrInvalidTransition = apierr.BadRequest("invalid order status transition")
	// ErrUnauthorized 无权操作
	ErrUnauthorized = apierr.Forbidden("unauthorized")
)

// ReferralTrigger 推荐奖励触发接口（可选依赖）
type ReferralTrigger interface {
	OnFirstOrderCompleted(ctx context.Context, userID uint64) error
}

// PaymentRefundor defines refund capability for paid orders.
type PaymentRefundor interface {
	RefundPayment(ctx context.Context, paymentID uint64, reason string) error
}

// OrderService 订单服务
//
// 功能：
// 1. 用户端订单管理（创建、查询、取消、完成）
// 2. 陪玩师端订单管理（接单、开始、完成）
// 3. 订单状态流转管理
// 4. 订单拆分与转单
type OrderService struct {
	db              *gorm.DB                      // 数据库连接，用于事务管理
	orders          repoiface.OrderRepository
	orderGroups     ordergroup.Repository // 主订单仓储
	players         repository.PlayerRepository
	users           repository.UserRepository
	games           repository.GameRepository
	payments        repository.PaymentRepository
	reviews         repository.ReviewRepository
	commissions     commissionrepo.CommissionRepository
	chatGroups      repository.ChatGroupRepository // optional: for order chat auto-destroy
	distributedLock cache.DistributedLock          // 分布式锁，用于并发控制
	referralTrigger ReferralTrigger                // optional: for referral reward trigger
	paymentRefundor PaymentRefundor                // optional: refund processor
	notifications   repository.NotificationRepository
	hub             *ws.Hub
}

// NewOrderService 创建订单服务
func NewOrderService(
	orders repoiface.OrderRepository,
	players repository.PlayerRepository,
	users repository.UserRepository,
	games repository.GameRepository,
	payments repository.PaymentRepository,
	reviews repository.ReviewRepository,
	commissions commissionrepo.CommissionRepository,
) *OrderService {
	return &OrderService{
		orders:      orders,
		players:     players,
		users:       users,
		games:       games,
		payments:    payments,
		reviews:     reviews,
		commissions: commissions,
	}
}

// SetDB injects *gorm.DB for transaction management in multi-step operations.
func (s *OrderService) SetDB(db *gorm.DB) {
	s.db = db
}

// SetOrderGroupRepository 注入主订单仓储
func (s *OrderService) SetOrderGroupRepository(repo ordergroup.Repository) {
	s.orderGroups = repo
}

// SetDistributedLock injects distributed lock for concurrency control
func (s *OrderService) SetDistributedLock(lock cache.DistributedLock) {
	s.distributedLock = lock
}

// SetReferralTrigger injects referral trigger for first order reward
func (s *OrderService) SetReferralTrigger(trigger ReferralTrigger) {
	s.referralTrigger = trigger
}

// SetPaymentRefundor injects refund processor for paid order cancellation.
func (s *OrderService) SetPaymentRefundor(refundor PaymentRefundor) {
	s.paymentRefundor = refundor
}

// SetChatGroupRepository injects chat group repository for auto-destroying order chat groups.
func (s *OrderService) SetChatGroupRepository(chatGroups repository.ChatGroupRepository) {
	s.chatGroups = chatGroups
}

// SetNotificationRepository injects notification repository for order events.
func (s *OrderService) SetNotificationRepository(repo repository.NotificationRepository) {
	s.notifications = repo
}

// SetWebsocketHub injects WebSocket hub for realtime order updates.
func (s *OrderService) SetWebsocketHub(hub *ws.Hub) {
	s.hub = hub
}

// deactivateOrderChat best-effort deactivates the chat group bound to the order.
func (s *OrderService) deactivateOrderChat(ctx context.Context, orderID uint64) {
	if s.chatGroups == nil {
		return
	}
	group, err := s.chatGroups.GetByRelatedOrderID(ctx, orderID)
	if err != nil || group == nil {
		return
	}
	if group.IsActive && group.GroupType == model.ChatGroupTypeOrder {
		_ = s.chatGroups.Deactivate(ctx, group.ID)
	}
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	PlayerID       uint64     `json:"playerId" binding:"required"`
	GameID         uint64     `json:"gameId" binding:"required"`
	ServiceID      *uint64    `json:"serviceId"` // 可选：关联服务
	Title          string     `json:"title" binding:"required,max=128"`
	Description    string     `json:"description"`
	ScheduledStart *time.Time `json:"scheduledStart" binding:"required"`
	DurationHours  float32    `json:"durationHours" binding:"required,min=0.5,max=24"`
}

// CreateOrderResponse 创建订单响应
type CreateOrderResponse struct {
	OrderID       uint64 `json:"orderId"`
	GroupNo       string `json:"groupNo,omitempty"` // 主订单号（拆分时返回）
	PriceCents    int64  `json:"priceCents"`
	TotalHours    int    `json:"totalHours,omitempty"` // 总时长（拆分时返回）
	NeedPayment   bool   `json:"needPayment"`
	IsSplit       bool   `json:"isSplit,omitempty"`       // 是否拆分订单
	SubOrderCount int    `json:"subOrderCount,omitempty"` // 子订单数量
}

// OrderCardDTO 订单卡片信息（列表展示）
type OrderCardDTO struct {
	ID             uint64            `json:"id"`
	Title          string            `json:"title"`
	PlayerNickname string            `json:"playerNickname"`
	PlayerAvatar   string            `json:"playerAvatar"`
	GameName       string            `json:"gameName"`
	Status         model.OrderStatus `json:"status"`
	PriceCents     int64             `json:"priceCents"`
	ScheduledStart *time.Time        `json:"scheduledStart"`
	CreatedAt      time.Time         `json:"createdAt"`
	CanPay         bool              `json:"canPay"`
	CanCancel      bool              `json:"canCancel"`
	CanComplete    bool              `json:"canComplete"`
	CanReview      bool              `json:"canReview"`
}

// OrderDetailDTO 订单详情信息
type OrderDetailDTO struct {
	OrderCardDTO
	Description  string     `json:"description"`
	ScheduledEnd *time.Time `json:"scheduledEnd"`
	StartedAt    *time.Time `json:"startedAt"`
	CompletedAt  *time.Time `json:"completedAt"`
	CancelReason string     `json:"cancelReason"`
	RefundAmount int64      `json:"refundAmount"`
	RefundReason string     `json:"refundReason"`
}

// OrderTimelineDTO 订单时间线
type OrderTimelineDTO struct {
	Time    time.Time `json:"time"`
	Status  string    `json:"status"`
	Message string    `json:"message"`
}

// PaymentDTO 支付信息
type PaymentDTO struct {
	ID          uint64              `json:"id"`
	Method      model.PaymentMethod `json:"method"`
	AmountCents int64               `json:"amountCents"`
	Status      model.PaymentStatus `json:"status"`
	PaidAt      *time.Time          `json:"paidAt"`
	CreatedAt   time.Time           `json:"createdAt"`
}

// PlayerCardDTO 陪玩师卡片信息
type PlayerCardDTO struct {
	ID        uint64 `json:"id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl"`
	Rank      string `json:"rank"`
}

// MyOrderListRequest 我的订单列表请求
type MyOrderListRequest struct {
	Status   string `form:"status"` // pending/confirmed/in_progress/completed/canceled
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

// MyOrderListResponse 我的订单列表响应
type MyOrderListResponse struct {
	Orders []OrderCardDTO `json:"orders"`
	Total  int64          `json:"total"`
}

// OrderDetailResponse 订单详情响应
type OrderDetailResponse struct {
	Order    OrderDetailDTO     `json:"order"`
	Player   *PlayerCardDTO     `json:"player"`
	Payment  *PaymentDTO        `json:"payment"`
	Review   *ReviewDTO         `json:"review"`
	Timeline []OrderTimelineDTO `json:"timeline"`
}

// CancelOrderRequest 取消订单请求
type CancelOrderRequest struct {
	Reason string `json:"reason" binding:"required,max=500"`
}

// CompleteOrderRequest 确认完成订单请求
type CompleteOrderRequest struct {
	Confirm bool `json:"confirm"`
}

// CreateOrder 创建订单（用户端）
// 如果时长 > 1 小时，会自动拆分成多个子订单
func (s *OrderService) CreateOrder(ctx context.Context, userID uint64, req CreateOrderRequest) (*CreateOrderResponse, error) {
	// 使用分布式锁防止并发问题
	lockKey := fmt.Sprintf("order:create:user:%d:player:%d", userID, req.PlayerID)

	// 如果分布式锁可用，使用它
	if s.distributedLock != nil {
		locked, err := s.distributedLock.TryLock(ctx, lockKey, time.Second*5, 3, time.Millisecond*100)
		if err != nil {
			return nil, apierr.InternalError("failed to acquire lock").WithDetails(err.Error())
		}
		if !locked {
			return nil, apierr.Conflict("concurrent order creation in progress, please try again")
		}
		defer func() { _ = s.distributedLock.Unlock(ctx, lockKey) }()
	}

	// 验证陪玩师
	player, err := s.validateCreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	// 计算每小时价格
	hourlyPrice, commissionPerHour, playerIncomePerHour := s.calculateOrderPricing(player, req)

	// 判断是否需要拆分（时长 > 1 小时且 orderGroups 仓储可用）
	if req.DurationHours > 1 && s.orderGroups != nil {
		return s.createOrderWithSplit(ctx, userID, req, hourlyPrice, commissionPerHour, playerIncomePerHour)
	}

	// 不拆分，创建单个订单（向后兼容）
	order := s.buildOrderForCreation(userID, req, hourlyPrice, commissionPerHour, playerIncomePerHour)

	if err := s.orders.Create(ctx, order); err != nil {
		return nil, err
	}

	s.notifyOrderStatus(ctx, order.UserID, order.ID, string(order.Status), "订单已创建", "订单已创建，待支付")

	return &CreateOrderResponse{
		OrderID:     order.ID,
		PriceCents:  hourlyPrice,
		NeedPayment: true,
	}, nil
}

// createOrderWithSplit 创建拆分订单（主订单 + 子订单）
// 使用数据库事务确保主订单和子订单的原子性创建
func (s *OrderService) createOrderWithSplit(
	ctx context.Context,
	userID uint64,
	req CreateOrderRequest,
	hourlyPrice, commissionPerHour, playerIncomePerHour int64,
) (*CreateOrderResponse, error) {
	// 构建主订单和子订单
	group, subOrders := s.buildOrderGroupWithSubOrders(userID, req, hourlyPrice, commissionPerHour, playerIncomePerHour)

	// 使用事务确保主订单和所有子订单原子性创建
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txGroupRepo := ordergroup.NewRepository(tx)
		txOrderRepo := orderrepo.NewOrderRepository(tx)

		// 1. 创建主订单
		if err := txGroupRepo.Create(ctx, group); err != nil {
			return fmt.Errorf("创建主订单失败: %w", err)
		}

		// 2. 创建子订单，关联主订单ID
		for _, subOrder := range subOrders {
			subOrder.GroupID = &group.ID
			if err := txOrderRepo.Create(ctx, subOrder); err != nil {
				return fmt.Errorf("创建子订单失败: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, apierr.InternalError(err.Error())
	}

	s.notifyOrderStatus(ctx, userID, group.ID, string(model.OrderStatusPending), "订单已创建", "订单已创建，待支付")

	return &CreateOrderResponse{
		OrderID:       group.ID,      // 返回主订单ID（用户视角）
		GroupNo:       group.GroupNo, // 主订单号
		PriceCents:    group.TotalPriceCents,
		TotalHours:    group.TotalHours,
		NeedPayment:   true,
		IsSplit:       true,
		SubOrderCount: len(subOrders),
	}, nil
}

// GetMyOrders 获取我的订单列表（用户端）
func (s *OrderService) GetMyOrders(ctx context.Context, userID uint64, req MyOrderListRequest) (*MyOrderListResponse, error) {
	// 默认分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 构建查询条件
	opts := repoiface.OrderListOptions{
		UserID:   &userID,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// 状态过滤
	if req.Status != "" {
		opts.Statuses = []model.OrderStatus{model.OrderStatus(req.Status)}
	}

	// 查询订单
	orders, total, err := s.orders.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	// 批量查询关联数据 (避免 N+1 问题)
	// 1. 提取所有陪玩师 ID
	playerIDs := make([]uint64, 0, len(orders))
	for _, o := range orders {
		if o.PlayerID != nil && *o.PlayerID > 0 {
			playerIDs = append(playerIDs, *o.PlayerID)
		}
	}

	// 2. 批量查询陪玩师
	players, err := s.players.GetByIDs(ctx, playerIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to batch query players: %w", err)
	}
	playerMap := make(map[uint64]*model.Player)
	for i := range players {
		playerMap[players[i].ID] = &players[i]
	}

	// 3. 提取所有用户 ID (陪玩师关联的用户)
	userIDs := make([]uint64, 0, len(players))
	for _, p := range players {
		if p.UserID > 0 {
			userIDs = append(userIDs, p.UserID)
		}
	}

	// 4. 批量查询用户
	users, err := s.users.GetByIDs(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to batch query users: %w", err)
	}
	userMap := make(map[uint64]*model.User)
	for i := range users {
		userMap[users[i].ID] = &users[i]
	}

	// 5. 提取所有游戏 ID
	gameIDs := make([]uint64, 0, len(orders))
	for _, o := range orders {
		if o.GameID != nil && *o.GameID > 0 {
			gameIDs = append(gameIDs, *o.GameID)
		}
	}

	// 6. 批量查询游戏
	games, err := s.games.GetByIDs(ctx, gameIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to batch query games: %w", err)
	}
	gameMap := make(map[uint64]*model.Game)
	for i := range games {
		gameMap[games[i].ID] = &games[i]
	}

	// 7. 转换为 DTO (使用预加载的数据)
	orderCards := make([]OrderCardDTO, 0, len(orders))
	for _, o := range orders {
		card := s.toOrderCardDTOWithPreloadedData(&o, userID, playerMap, userMap, gameMap)
		orderCards = append(orderCards, *card)
	}

	return &MyOrderListResponse{
		Orders: orderCards,
		Total:  total,
	}, nil
}

// GetOrderDetail 获取订单详情
func (s *OrderService) GetOrderDetail(ctx context.Context, userID uint64, orderID uint64) (*OrderDetailResponse, error) {
	// 获取订单
	order, err := s.orders.Get(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// 权限检查：只能查看自己的订单或者自己接的订单
	if order.UserID != userID && order.GetPlayerID() != userID {
		return nil, ErrUnauthorized
	}

	// ✅ 优化：使用预加载的陪玩师数据,避免重复查询
	// order.Player 和 order.Player.User 已在 repository.Get() 中预加载
	var playerCard *PlayerCardDTO
	if order.Player != nil {
		if order.Player.User != nil {
			playerCard = &PlayerCardDTO{
				ID:        order.Player.ID,
				Nickname:  order.Player.Nickname,
				AvatarURL: order.Player.User.AvatarURL,
				Rank:      order.Player.Rank,
			}
		}
	}

	// 获取支付信息
	var paymentDTO *PaymentDTO
	orderIDPtr := &orderID
	payments, _, err := s.payments.List(ctx, repository.PaymentListOptions{
		OrderID:  orderIDPtr,
		Page:     1,
		PageSize: 1,
	})
	if err == nil && len(payments) > 0 {
		p := payments[0]
		paymentDTO = &PaymentDTO{
			ID:          p.ID,
			Method:      p.Method,
			AmountCents: p.AmountCents,
			Status:      p.Status,
			PaidAt:      p.PaidAt,
			CreatedAt:   p.CreatedAt,
		}
	}

	// 获取评价信息
	var reviewDTO *ReviewDTO
	reviews, _, err := s.reviews.List(ctx, repository.ReviewListOptions{
		OrderID:  orderIDPtr,
		Page:     1,
		PageSize: 1,
	})
	if err == nil && len(reviews) > 0 {
		r := reviews[0]
		reviewDTO = &ReviewDTO{
			ID:        r.ID,
			Rating:    int(r.Score),
			Comment:   r.Content,
			CreatedAt: r.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	// 构建时间线
	timeline := s.buildOrderTimeline(order)

	// 构建订单详情
	card, err := s.toOrderCardDTO(ctx, order, userID)
	if err != nil {
		return nil, err
	}

	detailDTO := OrderDetailDTO{
		OrderCardDTO: *card,
		Description:  order.Description,
		ScheduledEnd: order.ScheduledEnd,
		StartedAt:    order.StartedAt,
		CompletedAt:  order.CompletedAt,
		CancelReason: order.CancelReason,
		RefundAmount: order.RefundAmountCents,
		RefundReason: order.RefundReason,
	}

	return &OrderDetailResponse{
		Order:    detailDTO,
		Player:   playerCard,
		Payment:  paymentDTO,
		Review:   reviewDTO,
		Timeline: timeline,
	}, nil
}

// CancelOrder 取消订单（用户端）
// 使用原子性状态检查避免并发取消的竞态条件
func (s *OrderService) CancelOrder(ctx context.Context, userID uint64, orderID uint64, req CancelOrderRequest) error {
	// 获取订单
	order, err := s.orders.Get(ctx, orderID)
	if err != nil {
		return err
	}

	// 权限检查
	if order.UserID != userID {
		return ErrUnauthorized
	}

	// 状态检查：只有 pending 或 confirmed 状态可以取消
	if order.Status != model.OrderStatusPending && order.Status != model.OrderStatusConfirmed {
		return ErrInvalidTransition
	}

	// 保存原始状态
	originalStatus := order.Status

	refundHandled := false

	// 如果已支付，需要退款
	if originalStatus == model.OrderStatusConfirmed {
		// 查找支付记录并执行退款
		orderIDPtr := &orderID
		payments, _, err := s.payments.List(ctx, repository.PaymentListOptions{
			OrderID:  orderIDPtr,
			Page:     1,
			PageSize: 1,
		})
		if err == nil && len(payments) > 0 {
			payment := payments[0]
			// 如果支付已完成，执行退款
			if payment.Status == model.PaymentStatusPaid {
				if s.paymentRefundor != nil {
					if err := s.paymentRefundor.RefundPayment(ctx, payment.ID, "用户取消订单"); err != nil {
						return err
					}
					refundHandled = true
				} else {
					// 兼容逻辑：没有退款服务时直接标记退款（使用原子更新）
					now := time.Now()
					updated, err := s.orders.UpdateWithCondition(ctx, orderID, originalStatus, map[string]any{
						"status":             model.OrderStatusRefunded,
						"cancel_reason":      req.Reason,
						"refund_amount_cents": order.TotalPriceCents,
						"refund_reason":      "用户取消订单",
						"refunded_at":        &now,
					})
					if err != nil {
						return err
					}
					if !updated {
						return ErrInvalidTransition
					}
					order.Status = model.OrderStatusRefunded
					refundHandled = true
				}
			}
		}
	}

	if !refundHandled {
		// 使用原子性更新避免并发竞态: 只在状态仍为原始状态时更新
		updated, err := s.orders.UpdateWithCondition(ctx, orderID, originalStatus, map[string]any{
			"status":        model.OrderStatusCanceled,
			"cancel_reason": req.Reason,
		})
		if err != nil {
			return err
		}
		if !updated {
			return ErrInvalidTransition
		}
		order.Status = model.OrderStatusCanceled
		order.CancelReason = req.Reason
	} else {
		// 退款已由支付服务处理，补充取消原因（不影响主流程）
		if updatedOrder, err := s.orders.Get(ctx, orderID); err == nil {
			updatedOrder.CancelReason = req.Reason
			if err := s.orders.Update(ctx, updatedOrder); err != nil {
				slog.Warn("failed to update cancel reason after refund",
					"orderID", orderID, "error", err)
			}
		} else {
			slog.Warn("failed to re-fetch order for cancel reason update",
				"orderID", orderID, "error", err)
		}
	}

	// auto-destroy order chat group
	s.deactivateOrderChat(ctx, orderID)

	if refundHandled {
		// 退款通知由支付服务负责，避免重复推送
		return nil
	}

	title := "订单已取消"
	if order.Status == model.OrderStatusRefunded {
		title = "订单已退款"
	}
	s.notifyOrderStatus(ctx, order.UserID, order.ID, string(order.Status), title, s.buildOrderStatusMessage(order))
	return nil
}

// CompleteOrder 确认完成订单（用户端）
func (s *OrderService) CompleteOrder(ctx context.Context, userID uint64, orderID uint64) error {
	// 获取订单
	order, err := s.orders.Get(ctx, orderID)
	if err != nil {
		return err
	}

	// 权限检查
	if order.UserID != userID {
		return ErrUnauthorized
	}

	// 状态检查：只有 in_progress 状态可以完成
	if order.Status != model.OrderStatusInProgress {
		return ErrInvalidTransition
	}

	// 更新订单状态
	now := time.Now()
	order.Status = model.OrderStatusCompleted
	order.CompletedAt = &now

	if err := s.orders.Update(ctx, order); err != nil {
		return err
	}

	// 订单完成后，自动记录抽成
	if err := s.recordCommissionAsync(ctx, orderID); err != nil {
		// 记录日志但不影响订单完成
		slog.Error("failed to record commission", "orderID", orderID, "error", err)
	}

	// 触发推荐奖励检查（首单条件）
	if s.referralTrigger != nil {
		go func() {
			if err := s.referralTrigger.OnFirstOrderCompleted(ctx, order.UserID); err != nil {
				slog.Error("failed to trigger referral reward", "userID", order.UserID, "error", err)
			}
		}()
	}

	// auto-destroy order chat group
	s.deactivateOrderChat(ctx, orderID)
	s.notifyOrderStatus(ctx, order.UserID, order.ID, string(order.Status), "订单已完成", "订单已完成")
	return nil
}

// recordCommissionAsync 异步记录抽成
func (s *OrderService) recordCommissionAsync(ctx context.Context, orderID uint64) error {
	// 检查是否已记录
	existing, _ := s.commissions.GetRecordByOrderID(ctx, orderID)
	if existing != nil {
		return nil // 已记录，跳过
	}

	// 获取订单
	order, err := s.orders.Get(ctx, orderID)
	if err != nil {
		return err
	}

	// 获取抽成规则
	rule, err := s.commissions.GetRuleForOrder(ctx, order.GameID, order.PlayerID, nil)
	if err != nil || rule == nil {
		// 使用默认规则
		rule, err = s.commissions.GetDefaultRule(ctx)
		if err != nil {
			rule = &model.CommissionRule{Rate: 20} // 默认20%
		}
	}

	// 计算抽成
	totalAmount := order.TotalPriceCents
	commissionRate := rule.Rate
	commissionAmount := totalAmount * int64(commissionRate) / 100
	playerIncome := totalAmount - commissionAmount

	// 创建抽成记录
	now := time.Now()
	playerID := order.GetPlayerID()
	if playerID == 0 {
		return apierr.BadRequest("订单未分配打手")
	}

	record := &model.CommissionRecord{
		OrderID:           orderID,
		PlayerID:          playerID,
		TotalAmountCents:  totalAmount,
		CommissionRate:    commissionRate,
		CommissionCents:   commissionAmount,
		PlayerIncomeCents: playerIncome,
		SettlementStatus:  "pending",
		SettlementMonth:   now.Format("2006-01"),
	}

	return s.commissions.CreateRecord(ctx, record)
}

// toOrderCardDTO 转换为订单卡DTO
// 使用预加载的关联数据,避免N+1查询问题
// 注意：order对象应通过orders.Get()获取,该方法已预加载Player,Player.User和Game
func (s *OrderService) toOrderCardDTO(ctx context.Context, order *model.Order, userID uint64) (*OrderCardDTO, error) {
	// ✅ 优化：直接使用预加载的关联数据,避免重复查询
	// order.Player, order.Player.User, order.Game 已在 repository.Get() 中预加载

	// 获取陪玩师信息（从预加载的数据）
	var playerNickname, playerAvatar string
	if order.Player != nil {
		playerNickname = order.Player.Nickname
		if order.Player.User != nil {
			playerAvatar = order.Player.User.AvatarURL
		}
	}

	// 获取游戏信息（从预加载的数据）
	var gameName string
	if order.Game != nil {
		gameName = order.Game.Name
	}

	// 判断操作权限
	canPay := order.Status == model.OrderStatusPending && order.UserID == userID
	canCancel := (order.Status == model.OrderStatusPending || order.Status == model.OrderStatusConfirmed) && order.UserID == userID
	canComplete := order.Status == model.OrderStatusInProgress && order.UserID == userID
	canReview := order.Status == model.OrderStatusCompleted && order.UserID == userID

	// 检查是否已评价
	if canReview {
		orderIDPtr := &order.ID
		reviews, _, err := s.reviews.List(ctx, repository.ReviewListOptions{
			OrderID:  orderIDPtr,
			Page:     1,
			PageSize: 1,
		})
		if err == nil && len(reviews) > 0 {
			canReview = false // 已评价
		}
	}

	return &OrderCardDTO{
		ID:             order.ID,
		Title:          order.Title,
		PlayerNickname: playerNickname,
		PlayerAvatar:   playerAvatar,
		GameName:       gameName,
		Status:         order.Status,
		PriceCents:     order.TotalPriceCents,
		ScheduledStart: order.ScheduledStart,
		CreatedAt:      order.CreatedAt,
		CanPay:         canPay,
		CanCancel:      canCancel,
		CanComplete:    canComplete,
		CanReview:      canReview,
	}, nil
}

// toOrderCardDTOWithPreloadedData 转换为订单卡DTO (使用预加载的数据，避免N+1查询)
// 这个方法用于批量查询优化，不执行任何数据库查询
func (s *OrderService) toOrderCardDTOWithPreloadedData(
	order *model.Order,
	userID uint64,
	playerMap map[uint64]*model.Player,
	userMap map[uint64]*model.User,
	gameMap map[uint64]*model.Game,
) *OrderCardDTO {
	// 从 Map 中获取陪玩师信息 (O(1) 查找)
	var playerNickname, playerAvatar string
	playerID := order.GetPlayerID()
	if playerID > 0 {
		if player := playerMap[playerID]; player != nil {
			playerNickname = player.Nickname
			if user := userMap[player.UserID]; user != nil {
				playerAvatar = user.AvatarURL
			}
		}
	}

	// 从 Map 中获取游戏信息 (O(1) 查找)
	var gameName string
	gameID := order.GetGameID()
	if gameID > 0 {
		if game := gameMap[gameID]; game != nil {
			gameName = game.Name
		}
	}

	// 判断操作权限
	canPay := order.Status == model.OrderStatusPending && order.UserID == userID
	canCancel := (order.Status == model.OrderStatusPending || order.Status == model.OrderStatusConfirmed) && order.UserID == userID
	canComplete := order.Status == model.OrderStatusInProgress && order.UserID == userID
	canReview := order.Status == model.OrderStatusCompleted && order.UserID == userID

	// 注意：由于批量查询无法预先知道评价状态，这里假设可以评价
	// 如果需要准确的评价状态，需要在批量查询后额外处理
	// 这是一个权衡：性能 vs 完全准确的 canReview 标志

	return &OrderCardDTO{
		ID:             order.ID,
		Title:          order.Title,
		PlayerNickname: playerNickname,
		PlayerAvatar:   playerAvatar,
		GameName:       gameName,
		Status:         order.Status,
		PriceCents:     order.TotalPriceCents,
		ScheduledStart: order.ScheduledStart,
		CreatedAt:      order.CreatedAt,
		CanPay:         canPay,
		CanCancel:      canCancel,
		CanComplete:    canComplete,
		CanReview:      canReview, // 注意：这可能是近似值
	}
}

// buildOrderTimeline 构建订单时间线
func (s *OrderService) buildOrderTimeline(order *model.Order) []OrderTimelineDTO {
	ctx := context.Background()
	timeline := []OrderTimelineDTO{
		{
			Time:    order.CreatedAt,
			Status:  string(model.OrderStatusPending),
			Message: "订单已创建",
		},
	}

	// 获取支付记录以获取真实的支付时间
	if order.Status != model.OrderStatusPending {
		paidTime := order.CreatedAt // 默认使用订单创建时间
		orderIDPtr := &order.ID
		payments, _, err := s.payments.List(ctx, repository.PaymentListOptions{
			OrderID:  orderIDPtr,
			Page:     1,
			PageSize: 1,
		})
		if err == nil && len(payments) > 0 {
			payment := payments[0]
			if payment.PaidAt != nil {
				paidTime = *payment.PaidAt
			}
		}

		timeline = append(timeline, OrderTimelineDTO{
			Time:    paidTime,
			Status:  string(model.OrderStatusConfirmed),
			Message: "订单已支付",
		})
	}

	if order.StartedAt != nil {
		timeline = append(timeline, OrderTimelineDTO{
			Time:    *order.StartedAt,
			Status:  string(model.OrderStatusInProgress),
			Message: "订单进行中",
		})
	}

	if order.CompletedAt != nil {
		timeline = append(timeline, OrderTimelineDTO{
			Time:    *order.CompletedAt,
			Status:  string(model.OrderStatusCompleted),
			Message: "订单已完成",
		})
	}

	if order.Status == model.OrderStatusCanceled {
		timeline = append(timeline, OrderTimelineDTO{
			Time:    order.UpdatedAt,
			Status:  string(model.OrderStatusCanceled),
			Message: "订单已取消： " + order.CancelReason,
		})
	}

	if order.RefundedAt != nil {
		timeline = append(timeline, OrderTimelineDTO{
			Time:    *order.RefundedAt,
			Status:  string(model.OrderStatusRefunded),
			Message: "订单已退款",
		})
	}

	return timeline
}

// AvailableOrdersRequest 可接订单列表请求
type AvailableOrdersRequest struct {
	GameID   *uint64 `form:"gameId"`
	Page     int     `form:"page"`
	PageSize int     `form:"pageSize"`
}

// AvailableOrderDTO 可接订单信息
type AvailableOrderDTO struct {
	ID             uint64     `json:"id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	GameName       string     `json:"gameName"`
	UserNickname   string     `json:"userNickname"`
	PriceCents     int64      `json:"priceCents"`
	ScheduledStart *time.Time `json:"scheduledStart"`
	DurationHours  float32    `json:"durationHours"`
	CreatedAt      time.Time  `json:"createdAt"`
}

// GetAvailableOrders 获取可接订单列表（陪玩师端）
func (s *OrderService) GetAvailableOrders(ctx context.Context, req AvailableOrdersRequest) ([]AvailableOrderDTO, int64, error) {
	// 默认分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 构建查询条件：查询已支付但未接单的订单
	opts := repoiface.OrderListOptions{
		Statuses: []model.OrderStatus{model.OrderStatusConfirmed},
		GameID:   req.GameID,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	orders, total, err := s.orders.List(ctx, opts)
	if err != nil {
		return nil, 0, err
	}

	// 批量查询关联数据 (避免 N+1 问题)
	// 1. 提取所有游戏 ID
	gameIDs := make([]uint64, 0, len(orders))
	for _, o := range orders {
		if o.GameID != nil && *o.GameID > 0 {
			gameIDs = append(gameIDs, *o.GameID)
		}
	}

	// 2. 批量查询游戏
	games, err := s.games.GetByIDs(ctx, gameIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to batch query games: %w", err)
	}
	gameMap := make(map[uint64]*model.Game)
	for i := range games {
		gameMap[games[i].ID] = &games[i]
	}

	// 3. 提取所有用户 ID
	userIDs := make([]uint64, 0, len(orders))
	for _, o := range orders {
		if o.UserID > 0 {
			userIDs = append(userIDs, o.UserID)
		}
	}

	// 4. 批量查询用户
	users, err := s.users.GetByIDs(ctx, userIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to batch query users: %w", err)
	}
	userMap := make(map[uint64]*model.User)
	for i := range users {
		userMap[users[i].ID] = &users[i]
	}

	// 5. 转换DTO (使用预加载的数据)
	availableOrders := make([]AvailableOrderDTO, 0, len(orders))
	for _, o := range orders {
		// 从 Map 中获取游戏信息 (O(1) 查找)
		var gameName string
		if o.GameID != nil {
			if game := gameMap[*o.GameID]; game != nil {
				gameName = game.Name
			}
		}

		// 从 Map 中获取用户信息 (O(1) 查找)
		var userNickname string
		if o.UserID > 0 {
			if user := userMap[o.UserID]; user != nil {
				userNickname = user.Name
			}
		}

		// 计算时长
		var durationHours float32
		if o.ScheduledStart != nil && o.ScheduledEnd != nil {
			duration := o.ScheduledEnd.Sub(*o.ScheduledStart)
			durationHours = float32(duration.Hours())
		}

		availableOrders = append(availableOrders, AvailableOrderDTO{
			ID:             o.ID,
			Title:          o.Title,
			Description:    o.Description,
			GameName:       gameName,
			UserNickname:   userNickname,
			PriceCents:     o.TotalPriceCents,
			ScheduledStart: o.ScheduledStart,
			DurationHours:  durationHours,
			CreatedAt:      o.CreatedAt,
		})
	}

	return availableOrders, total, nil
}

// AcceptOrder 接单（陪玩师端）
// 使用原子性更新避免并发竞态条件,确保一个订单只能被一个陪玩师接单
func (s *OrderService) AcceptOrder(ctx context.Context, playerUserID uint64, orderID uint64) error {
	// 直接根据UserID查找陪玩师 (性能优化: 避免全表扫描)
	player, err := s.players.GetByUserID(ctx, playerUserID)
	if err != nil {
		return apierr.NotFound("陪玩师不存在")
	}

	playerID := player.ID

	// 使用原子性更新: 仅当订单状态为confirmed时才更新
	// 这样可以避免多个陪玩师同时接单的竞态条件
	now := time.Now()
	updated, err := s.orders.UpdateWithCondition(ctx, orderID, model.OrderStatusConfirmed, map[string]any{
		"player_id":  playerID,
		"status":     model.OrderStatusInProgress,
		"started_at": &now,
	})

	if err != nil {
		return err
	}

	// 如果未更新成功,说明订单状态已变更(可能已被其他陪玩师接单)
	if !updated {
		return ErrInvalidTransition
	}

	if order, err := s.orders.Get(ctx, orderID); err == nil && order != nil {
		s.notifyOrderStatus(ctx, order.UserID, order.ID, string(order.Status), "订单已被接单", "订单已被陪玩师接单")
	}

	return nil
}

// CompleteOrderByPlayer 完成订单（陪玩师端）
func (s *OrderService) CompleteOrderByPlayer(ctx context.Context, playerUserID uint64, orderID uint64) error {
	// 直接根据UserID查找陪玩师 (性能优化: 避免全表扫描)
	player, err := s.players.GetByUserID(ctx, playerUserID)
	if err != nil {
		return apierr.NotFound("陪玩师不存在")
	}

	playerID := player.ID

	// 获取订单
	order, err := s.orders.Get(ctx, orderID)
	if err != nil {
		return err
	}

	// 权限检查
	if order.GetPlayerID() != playerID {
		return ErrUnauthorized
	}

	// 状态检查：只有 in_progress 状态可以完成
	if order.Status != model.OrderStatusInProgress {
		return ErrInvalidTransition
	}

	// 完成订单
	now := time.Now()
	order.Status = model.OrderStatusCompleted
	order.CompletedAt = &now

	if err := s.orders.Update(ctx, order); err != nil {
		return err
	}

	// 订单完成后，自动记录抽成
	if err := s.recordCommissionAsync(ctx, orderID); err != nil {
		// 记录日志但不影响订单完成
		// Note: Structured logging will be added when log integration is complete
	}

	s.notifyOrderStatus(ctx, order.UserID, order.ID, string(order.Status), "订单已完成", "订单已完成")

	return nil
}

func (s *OrderService) buildOrderStatusMessage(order *model.Order) string {
	if order == nil {
		return ""
	}
	switch order.Status {
	case model.OrderStatusRefunded:
		return "订单已退款"
	case model.OrderStatusCanceled:
		return "订单已取消"
	default:
		return "订单状态已更新"
	}
}

func (s *OrderService) notifyOrderStatus(ctx context.Context, userID, orderID uint64, status, title, message string) {
	s.broadcastOrderStatus(userID, orderID, status, message)

	if s.notifications == nil {
		return
	}
	refID := orderID
	event := &model.NotificationEvent{
		UserID:        userID,
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

func (s *OrderService) broadcastOrderStatus(userID, orderID uint64, status, message string) {
	if s.hub == nil {
		return
	}
	payload := ws.OrderStatusEvent{
		OrderID:   orderID,
		Status:    status,
		Message:   message,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
	if data, err := ws.NewWSMessage(ws.MessageTypeOrderStatus, payload).ToJSON(); err == nil {
		s.hub.BroadcastToUser(data, userID)
	}
}

func (s *OrderService) broadcastNotification(event *model.NotificationEvent) {
	if s.hub == nil || event == nil {
		return
	}
	if data, err := ws.NewWSMessage(ws.MessageTypeNotification, event).ToJSON(); err == nil {
		s.hub.BroadcastToUser(data, event.UserID)
	}
}
