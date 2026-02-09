package notification

import (
	"context"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// Service manages notification center workflows.
type Service struct {
	repo repository.NotificationRepository
}

// NewService constructs a notification service.
func NewService(repo repository.NotificationRepository) *Service {
	return &Service{repo: repo}
}

// ListRequest wraps pagination filters.
type ListRequest struct {
	Page       int
	PageSize   int
	UnreadOnly bool
	Priorities []model.NotificationPriority
}

// NotificationView renders response.
type NotificationView struct {
	ID            uint64                     `json:"id"`
	Title         string                     `json:"title"`
	Message       string                     `json:"message"`
	Priority      model.NotificationPriority `json:"priority"`
	Channel       string                     `json:"channel"`
	ReferenceType string                     `json:"referenceType"`
	ReferenceID   *uint64                    `json:"referenceId,omitempty"`
	ReadAt        *time.Time                 `json:"readAt,omitempty"`
	CreatedAt     time.Time                  `json:"createdAt"`
}

// ListResponse holds notifications with pagination.
type ListResponse struct {
	Items       []NotificationView `json:"items"`
	Page        int                `json:"page"`
	PageSize    int                `json:"pageSize"`
	Total       int64              `json:"total"`
	UnreadCount int64              `json:"unreadCount"`
}

// List fetches notifications by user.
func (s *Service) List(ctx context.Context, userID uint64, req ListRequest) (*ListResponse, error) {
	var unreadFilter *bool
	if req.UnreadOnly {
		val := true
		unreadFilter = &val
	}
	items, total, err := s.repo.ListByUser(ctx, repository.NotificationListOptions{
		Page:     req.Page,
		PageSize: req.PageSize,
		UserID:   userID,
		Unread:   unreadFilter,
		Priority: req.Priorities,
	})
	if err != nil {
		return nil, err
	}
	unreadCount, err := s.repo.CountUnread(ctx, userID)
	if err != nil {
		return nil, err
	}
	resp := &ListResponse{
		Items:       make([]NotificationView, 0, len(items)),
		Page:        req.Page,
		PageSize:    req.PageSize,
		Total:       total,
		UnreadCount: unreadCount,
	}
	for _, item := range items {
		view := NotificationView{
			ID:            item.ID,
			Title:         item.Title,
			Message:       item.Message,
			Priority:      item.Priority,
			Channel:       item.Channel,
			ReferenceType: item.ReferenceType,
			ReferenceID:   item.ReferenceID,
			CreatedAt:     item.CreatedAt,
		}
		if item.ReadAt != nil {
			view.ReadAt = item.ReadAt
		}
		resp.Items = append(resp.Items, view)
	}
	return resp, nil
}

// MarkRead marks notifications as read.
func (s *Service) MarkRead(ctx context.Context, userID uint64, ids []uint64) error {
	return s.repo.MarkRead(ctx, userID, ids)
}

// GetUnreadCount returns unread notifications count.
func (s *Service) GetUnreadCount(ctx context.Context, userID uint64) (int64, error) {
	return s.repo.CountUnread(ctx, userID)
}

// SendNotificationRequest 发送通知请求
type SendNotificationRequest struct {
	UserID        uint64                     // 接收者用户ID
	Title         string                     // 通知标题
	Message       string                     // 通知内容
	Priority      model.NotificationPriority // 优先级
	Channel       string                     // 渠道（默认 "in_app"）
	ReferenceType string                     // 关联类型（如 order, gift, coupon）
	ReferenceID   *uint64                    // 关联ID
}

// Send 发送通知
func (s *Service) Send(ctx context.Context, req SendNotificationRequest) error {
	if req.Channel == "" {
		req.Channel = "in_app"
	}
	if req.Priority == "" {
		req.Priority = model.NotificationPriorityNormal
	}

	notification := &model.NotificationEvent{
		UserID:        req.UserID,
		Title:         req.Title,
		Message:       req.Message,
		Priority:      req.Priority,
		Channel:       req.Channel,
		ReferenceType: req.ReferenceType,
		ReferenceID:   req.ReferenceID,
	}

	return s.repo.Create(ctx, notification)
}

// SendGiftNotification 发送礼物通知
func (s *Service) SendGiftNotification(ctx context.Context, playerUserID uint64, senderName string, giftName string, quantity int, orderID uint64) error {
	title := "收到新礼物"
	message := senderName + " 赠送了您 " + giftName
	if quantity > 1 {
		message += " x" + formatInt(quantity)
	}

	return s.Send(ctx, SendNotificationRequest{
		UserID:        playerUserID,
		Title:         title,
		Message:       message,
		Priority:      model.NotificationPriorityNormal,
		ReferenceType: "gift",
		ReferenceID:   &orderID,
	})
}

// SendOrderNotification 发送订单通知
func (s *Service) SendOrderNotification(ctx context.Context, userID uint64, title, message string, orderID uint64) error {
	return s.Send(ctx, SendNotificationRequest{
		UserID:        userID,
		Title:         title,
		Message:       message,
		Priority:      model.NotificationPriorityNormal,
		ReferenceType: "order",
		ReferenceID:   &orderID,
	})
}

// SendVipNotification 发送 VIP 相关通知
func (s *Service) SendVipNotification(ctx context.Context, userID uint64, title, message string) error {
	return s.Send(ctx, SendNotificationRequest{
		UserID:        userID,
		Title:         title,
		Message:       message,
		Priority:      model.NotificationPriorityNormal,
		ReferenceType: "vip",
	})
}

// SendCouponNotification 发送优惠券通知
func (s *Service) SendCouponNotification(ctx context.Context, userID uint64, title, message string, couponID *uint64) error {
	return s.Send(ctx, SendNotificationRequest{
		UserID:        userID,
		Title:         title,
		Message:       message,
		Priority:      model.NotificationPriorityNormal,
		ReferenceType: "coupon",
		ReferenceID:   couponID,
	})
}

// SendReferralRewardNotification 发送推荐奖励通知
func (s *Service) SendReferralRewardNotification(ctx context.Context, userID uint64, rewardType string, amount int64) error {
	title := "推荐奖励已到账"
	var message string
	if rewardType == "cash" {
		message = "恭喜您获得推荐奖励 ¥" + formatCents(amount) + "，已存入您的钱包"
	} else {
		message = "恭喜您获得推荐奖励，优惠券已发放到您的账户"
	}

	return s.Send(ctx, SendNotificationRequest{
		UserID:        userID,
		Title:         title,
		Message:       message,
		Priority:      model.NotificationPriorityNormal,
		ReferenceType: "referral",
	})
}

// formatInt 格式化整数
func formatInt(n int) string {
	return fmt.Sprintf("%d", n)
}

// formatCents 格式化分为元
func formatCents(cents int64) string {
	yuan := cents / 100
	fen := cents % 100
	if fen == 0 {
		return fmt.Sprintf("%d", yuan)
	}
	return fmt.Sprintf("%d.%02d", yuan, fen)
}
