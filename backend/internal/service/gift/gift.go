package gift

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionrepo "gamelink/internal/repository/commission"
	serviceitemrepo "gamelink/internal/repository/serviceitem"
)

var (
	// ErrNotFound 记录不存�?
	ErrNotFound = repository.ErrNotFound
	// ErrValidation 表示输入校验失败
	ErrValidation = errors.New("validation failed")
	// ErrInvalidGiftItem 无效的礼物项�?
	ErrInvalidGiftItem = errors.New("invalid gift item")
)

// GiftService 礼物服务（基于统一订单系统�?
type GiftService struct {
	items       serviceitemrepo.ServiceItemRepository
	orders      repository.OrderRepository
	players     repository.PlayerRepository
	commissions commissionrepo.CommissionRepository
}

// NewGiftService 创建礼物服务
func NewGiftService(
	items serviceitemrepo.ServiceItemRepository,
	orders repository.OrderRepository,
	players repository.PlayerRepository,
	commissions commissionrepo.CommissionRepository,
) *GiftService {
	return &GiftService{
		items:       items,
		orders:      orders,
		players:     players,
		commissions: commissions,
	}
}

// SendGiftRequest 赠送礼物请�?
type SendGiftRequest struct {
	PlayerID    uint64  `json:"playerId" binding:"required"`              // 接收礼物的陪玩师
	GiftItemID  uint64  `json:"giftItemId" binding:"required"`            // 礼物项目ID
	Quantity    int     `json:"quantity" binding:"required,min=1,max=99"` // 数量
	Message     string  `json:"message" binding:"max=200"`                // 留言
	IsAnonymous bool    `json:"isAnonymous"`                              // 是否匿名
	OrderID     *uint64 `json:"orderId"`                                  // 关联的护航订单（可选）
}

// SendGift 赠送礼�?
func (s *GiftService) SendGift(ctx context.Context, userID uint64, req SendGiftRequest) (*GiftOrderResponse, error) {
	// 1. 验证礼物项目
	giftItem, err := s.items.Get(ctx, req.GiftItemID)
	if err != nil {
		return nil, err
	}

	// 确保是礼物类�?
	if !giftItem.IsGift() {
		return nil, ErrInvalidGiftItem
	}

	// 确保礼物是激活状�?
	if !giftItem.IsActive {
		return nil, errors.New("gift item is not active")
	}

	// 2. 验证陪玩�?
	player, err := s.players.Get(ctx, req.PlayerID)
	if err != nil {
		return nil, fmt.Errorf("player not found: %w", err)
	}

	// 3. 计算价格和抽�?
	platformCommission, playerIncome := giftItem.CalculateCommission(req.Quantity)
	totalPrice := giftItem.BasePriceCents * int64(req.Quantity)

	// 4. 生成订单�?
	orderNo := generateOrderNo("GIFT")

	// 5. 创建礼物订单
	order := &model.Order{
		OrderNo:           orderNo,
		UserID:            userID,
		ItemID:            req.GiftItemID,
		PlayerID:          &req.PlayerID, // 礼物订单的PlayerID就是接收�?
		RecipientPlayerID: &req.PlayerID, // 明确标识接收�?
		Quantity:          req.Quantity,
		UnitPriceCents:    giftItem.BasePriceCents,
		TotalPriceCents:   totalPrice,
		CommissionCents:   platformCommission,
		PlayerIncomeCents: playerIncome,
		Status:            model.OrderStatusPending,
		Title:             fmt.Sprintf("赠送礼物：%s", giftItem.Name),
		Description:       fmt.Sprintf("�?%s 赠�?%s x%d", player.Nickname, giftItem.Name, req.Quantity),
		GiftMessage:       req.Message,
		IsAnonymous:       req.IsAnonymous,
	}

	// 创建订单
	if err := s.orders.Create(ctx, order); err != nil {
		return nil, err
	}

	// TODO: 自动完成支付（这里简化，实际需要支付流程）
	// 礼物订单支付后立即送达
	if err := s.deliverGift(ctx, order); err != nil {
		return nil, err
	}

	return &GiftOrderResponse{
		OrderID:     order.ID,
		OrderNo:     order.OrderNo,
		PlayerID:    req.PlayerID,
		GiftName:    giftItem.Name,
		Quantity:    req.Quantity,
		TotalPrice:  totalPrice,
		Status:      string(order.Status),
		DeliveredAt: order.DeliveredAt,
	}, nil
}

// deliverGift 送达礼物
func (s *GiftService) deliverGift(ctx context.Context, order *model.Order) error {
	// 更新订单状态为已完�?
	now := time.Now()
	order.Status = model.OrderStatusCompleted
	order.DeliveredAt = &now
	order.CompletedAt = &now

	if err := s.orders.Update(ctx, order); err != nil {
		return err
	}

	// 记录抽成
	record := &model.CommissionRecord{
		OrderID:           order.ID,
		PlayerID:          *order.RecipientPlayerID,
		TotalAmountCents:  order.TotalPriceCents,
		CommissionRate:    int(order.CommissionCents * 100 / order.TotalPriceCents),
		CommissionCents:   order.CommissionCents,
		PlayerIncomeCents: order.PlayerIncomeCents,
		SettlementStatus:  "pending",
		SettlementMonth:   now.Format("2006-01"),
	}

	if err := s.commissions.CreateRecord(ctx, record); err != nil {
		// 记录抽成失败不影响礼物送达
		// TODO: 记录日志
	}

	// TODO: 发送通知给陪玩师

	return nil
}

// GiftOrderResponse 礼物订单响应
type GiftOrderResponse struct {
	OrderID     uint64     `json:"orderId"`
	OrderNo     string     `json:"orderNo"`
	PlayerID    uint64     `json:"playerId"`
	GiftName    string     `json:"giftName"`
	Quantity    int        `json:"quantity"`
	TotalPrice  int64      `json:"totalPrice"`
	Status      string     `json:"status"`
	DeliveredAt *time.Time `json:"deliveredAt"`
}

// GetPlayerReceivedGifts 获取陪玩师收到的礼物
func (s *GiftService) GetPlayerReceivedGifts(ctx context.Context, playerID uint64, page, pageSize int) (*ReceivedGiftsResponse, error) {
	// 查询该陪玩师收到的礼物订�?
	orders, _, err := s.orders.List(ctx, repository.OrderListOptions{
		PlayerID: &playerID,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, err
	}

	// 只返回礼物订�?
	gifts := make([]ReceivedGiftDTO, 0)
	for _, order := range orders {
		if order.IsGiftOrder() {
			// 获取礼物项目信息
			item, err := s.items.Get(ctx, order.ItemID)
			if err != nil {
				continue
			}

			gifts = append(gifts, ReceivedGiftDTO{
				OrderID:     order.ID,
				OrderNo:     order.OrderNo,
				GiftName:    item.Name,
				GiftIcon:    item.IconURL,
				Quantity:    order.Quantity,
				TotalPrice:  order.TotalPriceCents,
				Income:      order.PlayerIncomeCents,
				Message:     order.GiftMessage,
				IsAnonymous: order.IsAnonymous,
				SenderName:  "", // TODO: 如果不匿名，获取发送者信�?
				DeliveredAt: order.DeliveredAt,
				CreatedAt:   order.CreatedAt,
			})
		}
	}

	return &ReceivedGiftsResponse{
		Gifts: gifts,
		Total: int64(len(gifts)),
	}, nil
}

// ReceivedGiftDTO 收到的礼物DTO
type ReceivedGiftDTO struct {
	OrderID     uint64     `json:"orderId"`
	OrderNo     string     `json:"orderNo"`
	GiftName    string     `json:"giftName"`
	GiftIcon    string     `json:"giftIcon"`
	Quantity    int        `json:"quantity"`
	TotalPrice  int64      `json:"totalPrice"`
	Income      int64      `json:"income"`
	Message     string     `json:"message"`
	IsAnonymous bool       `json:"isAnonymous"`
	SenderName  string     `json:"senderName,omitempty"`
	DeliveredAt *time.Time `json:"deliveredAt"`
	CreatedAt   time.Time  `json:"createdAt"`
}

// ReceivedGiftsResponse 收到的礼物列表响�?
type ReceivedGiftsResponse struct {
	Gifts []ReceivedGiftDTO `json:"gifts"`
	Total int64             `json:"total"`
}

// GetGiftStats 获取陪玩师礼物统�?
func (s *GiftService) GetGiftStats(ctx context.Context, playerID uint64) (*GiftStatsResponse, error) {
	// 查询所有礼物订�?
	playerIDPtr := &playerID
	orders, _, err := s.orders.List(ctx, repository.OrderListOptions{
		PlayerID: playerIDPtr,
		Statuses: []model.OrderStatus{model.OrderStatusCompleted},
		Page:     1,
		PageSize: 10000,
	})
	if err != nil {
		return nil, err
	}

	var totalGifts int64
	var totalIncome int64
	var giftCount int64

	for _, order := range orders {
		if order.IsGiftOrder() {
			totalGifts += int64(order.Quantity)
			totalIncome += order.PlayerIncomeCents
			giftCount++
		}
	}

	return &GiftStatsResponse{
		TotalGiftsReceived: totalGifts,
		TotalGiftIncome:    totalIncome,
		TotalGiftOrders:    giftCount,
	}, nil
}

// GiftStatsResponse 礼物统计响应
type GiftStatsResponse struct {
	TotalGiftsReceived int64 `json:"totalGiftsReceived"` // 收到的礼物总数
	TotalGiftIncome    int64 `json:"totalGiftIncome"`    // 礼物总收入（分）
	TotalGiftOrders    int64 `json:"totalGiftOrders"`    // 礼物订单�?
}

// generateOrderNo 生成订单�?
func generateOrderNo(prefix string) string {
	return fmt.Sprintf("%s%d", prefix, time.Now().UnixNano())
}
