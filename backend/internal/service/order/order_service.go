package order

import (
	"context"
	"errors"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

var (
	// ErrNotFound 订单不存在
	ErrNotFound = repository.ErrNotFound
	// ErrValidation 表示输入校验失败
	ErrValidation = errors.New("validation failed")
	// ErrInvalidTransition 订单状态流转不合法
	ErrInvalidTransition = errors.New("invalid order status transition")
	// ErrUnauthorized 无权操作
	ErrUnauthorized = errors.New("unauthorized")
)

// OrderService 订单服务
//
// 功能：
// 1. 用户端订单管理（创建、查询、取消、完成）
// 2. 陪玩师端订单管理（接单、开始、完成）
// 3. 订单状态流转管理
type OrderService struct {
	orders   repository.OrderRepository
	players  repository.PlayerRepository
	users    repository.UserRepository
	games    repository.GameRepository
	payments repository.PaymentRepository
	reviews  repository.ReviewRepository
}

// NewOrderService 创建订单服务
func NewOrderService(
	orders repository.OrderRepository,
	players repository.PlayerRepository,
	users repository.UserRepository,
	games repository.GameRepository,
	payments repository.PaymentRepository,
	reviews repository.ReviewRepository,
) *OrderService {
	return &OrderService{
		orders:   orders,
		players:  players,
		users:    users,
		games:    games,
		payments: payments,
		reviews:  reviews,
	}
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	PlayerID       uint64     `json:"playerId" binding:"required"`
	GameID         uint64     `json:"gameId" binding:"required"`
	Title          string     `json:"title" binding:"required,max=128"`
	Description    string     `json:"description"`
	ScheduledStart *time.Time `json:"scheduledStart" binding:"required"`
	DurationHours  float32    `json:"durationHours" binding:"required,min=1,max=24"`
}

// CreateOrderResponse 创建订单响应
type CreateOrderResponse struct {
	OrderID     uint64 `json:"orderId"`
	PriceCents  int64  `json:"priceCents"`
	NeedPayment bool   `json:"needPayment"`
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

// OrderTimelineDTO 订单时间轴
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

// ReviewDTO 评价信息
type ReviewDTO struct {
	ID        uint64 `json:"id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"createdAt"`
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
func (s *OrderService) CreateOrder(ctx context.Context, userID uint64, req CreateOrderRequest) (*CreateOrderResponse, error) {
	// 验证陪玩师
	player, err := s.players.Get(ctx, req.PlayerID)
	if err != nil {
		return nil, err
	}

	// 验证游戏
	_, err = s.games.Get(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	// 计算价格
	priceCents := int64(float32(player.HourlyRateCents) * req.DurationHours)

	// 计算结束时间
	scheduledEnd := req.ScheduledStart.Add(time.Duration(req.DurationHours * float32(time.Hour)))

	// 创建订单
	order := &model.Order{
		UserID:         userID,
		PlayerID:       req.PlayerID,
		GameID:         req.GameID,
		Title:          req.Title,
		Description:    req.Description,
		Status:         model.OrderStatusPending, // 待支付
		PriceCents:     priceCents,
		Currency:       model.CurrencyCNY,
		ScheduledStart: req.ScheduledStart,
		ScheduledEnd:   &scheduledEnd,
	}

	if err := s.orders.Create(ctx, order); err != nil {
		return nil, err
	}

	return &CreateOrderResponse{
		OrderID:     order.ID,
		PriceCents:  priceCents,
		NeedPayment: true,
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
	opts := repository.OrderListOptions{
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

	// 转换为 DTO
	orderCards := make([]OrderCardDTO, 0, len(orders))
	for _, o := range orders {
		card, err := s.toOrderCardDTO(ctx, &o, userID)
		if err != nil {
			continue
		}
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
		return nil, err
	}

	// 权限检查：只能查看自己的订单或者自己接的订单
	if order.UserID != userID && order.PlayerID != userID {
		return nil, ErrUnauthorized
	}

	// 获取陪玩师信息
	var playerCard *PlayerCardDTO
	if order.PlayerID > 0 {
		player, err := s.players.Get(ctx, order.PlayerID)
		if err == nil {
			user, err := s.users.Get(ctx, player.UserID)
			if err == nil {
				playerCard = &PlayerCardDTO{
					ID:        player.ID,
					Nickname:  player.Nickname,
					AvatarURL: user.AvatarURL,
					Rank:      player.Rank,
				}
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

	// 构建时间轴
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

	// 状态检查：只有 pending 和 confirmed 状态可以取消
	if order.Status != model.OrderStatusPending && order.Status != model.OrderStatusConfirmed {
		return ErrInvalidTransition
	}

	// 更新订单状态
	order.Status = model.OrderStatusCanceled
	order.CancelReason = req.Reason

	// 如果已支付，需要退款
	if order.Status == model.OrderStatusConfirmed {
		// TODO: 处理退款逻辑
		now := time.Now()
		order.RefundAmountCents = order.PriceCents
		order.RefundReason = "用户取消订单"
		order.RefundedAt = &now
		order.Status = model.OrderStatusRefunded
	}

	return s.orders.Update(ctx, order)
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

	return s.orders.Update(ctx, order)
}

// toOrderCardDTO 转换为订单卡片 DTO
func (s *OrderService) toOrderCardDTO(ctx context.Context, order *model.Order, userID uint64) (*OrderCardDTO, error) {
	// 获取陪玩师信息
	var playerNickname, playerAvatar string
	if order.PlayerID > 0 {
		player, err := s.players.Get(ctx, order.PlayerID)
		if err == nil {
			playerNickname = player.Nickname
			user, err := s.users.Get(ctx, player.UserID)
			if err == nil {
				playerAvatar = user.AvatarURL
			}
		}
	}

	// 获取游戏信息
	var gameName string
	if order.GameID > 0 {
		game, err := s.games.Get(ctx, order.GameID)
		if err == nil {
			gameName = game.Name
		}
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
		PriceCents:     order.PriceCents,
		ScheduledStart: order.ScheduledStart,
		CreatedAt:      order.CreatedAt,
		CanPay:         canPay,
		CanCancel:      canCancel,
		CanComplete:    canComplete,
		CanReview:      canReview,
	}, nil
}

// buildOrderTimeline 构建订单时间轴
func (s *OrderService) buildOrderTimeline(order *model.Order) []OrderTimelineDTO {
	timeline := []OrderTimelineDTO{
		{
			Time:    order.CreatedAt,
			Status:  string(model.OrderStatusPending),
			Message: "订单已创建",
		},
	}

	if order.Status != model.OrderStatusPending {
		// TODO: 根据支付时间添加确认时间
		timeline = append(timeline, OrderTimelineDTO{
			Time:    order.CreatedAt,
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
			Message: "订单已取消: " + order.CancelReason,
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
	opts := repository.OrderListOptions{
		Statuses: []model.OrderStatus{model.OrderStatusConfirmed},
		GameID:   req.GameID,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	orders, total, err := s.orders.List(ctx, opts)
	if err != nil {
		return nil, 0, err
	}

	// 转换为 DTO
	availableOrders := make([]AvailableOrderDTO, 0, len(orders))
	for _, o := range orders {
		// 获取游戏信息
		var gameName string
		if o.GameID > 0 {
			game, err := s.games.Get(ctx, o.GameID)
			if err == nil {
				gameName = game.Name
			}
		}

		// 获取用户信息
		var userNickname string
		if o.UserID > 0 {
			user, err := s.users.Get(ctx, o.UserID)
			if err == nil {
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
			PriceCents:     o.PriceCents,
			ScheduledStart: o.ScheduledStart,
			DurationHours:  durationHours,
			CreatedAt:      o.CreatedAt,
		})
	}

	return availableOrders, total, nil
}

// AcceptOrder 接单（陪玩师端）
func (s *OrderService) AcceptOrder(ctx context.Context, playerUserID uint64, orderID uint64) error {
	// 查找陪玩师
	players, _, err := s.players.ListPaged(ctx, 1, 100)
	if err != nil {
		return err
	}

	var playerID uint64
	for _, p := range players {
		if p.UserID == playerUserID {
			playerID = p.ID
			break
		}
	}

	if playerID == 0 {
		return errors.New("player not found")
	}

	// 获取订单
	order, err := s.orders.Get(ctx, orderID)
	if err != nil {
		return err
	}

	// 状态检查：只有 confirmed 状态可以接单
	if order.Status != model.OrderStatusConfirmed {
		return ErrInvalidTransition
	}

	// 接单
	order.PlayerID = playerID
	order.Status = model.OrderStatusInProgress
	now := time.Now()
	order.StartedAt = &now

	return s.orders.Update(ctx, order)
}

// CompleteOrderByPlayer 完成订单（陪玩师端）
func (s *OrderService) CompleteOrderByPlayer(ctx context.Context, playerUserID uint64, orderID uint64) error {
	// 查找陪玩师
	players, _, err := s.players.ListPaged(ctx, 1, 100)
	if err != nil {
		return err
	}

	var playerID uint64
	for _, p := range players {
		if p.UserID == playerUserID {
			playerID = p.ID
			break
		}
	}

	if playerID == 0 {
		return errors.New("player not found")
	}

	// 获取订单
	order, err := s.orders.Get(ctx, orderID)
	if err != nil {
		return err
	}

	// 权限检查
	if order.PlayerID != playerID {
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

	return s.orders.Update(ctx, order)
}
