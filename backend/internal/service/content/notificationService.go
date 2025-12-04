package content

import (
	"context"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// NotificationService manages notification center workflows.
type NotificationService struct {
	repo repository.NotificationRepository
}

// NewNotificationService constructs a notification service.
func NewNotificationService(repo repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

// NotificationListRequest wraps pagination filters.
type NotificationListRequest struct {
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
	IsRead        bool                       `json:"isRead"`
	CreatedAt     time.Time                  `json:"createdAt"`
}

// NotificationListResponse holds notifications with pagination.
type NotificationListResponse struct {
	Items       []NotificationView `json:"items"`
	Page        int                `json:"page"`
	PageSize    int                `json:"pageSize"`
	Total       int64              `json:"total"`
	UnreadCount int64              `json:"unreadCount"`
}

// ListNotifications fetches notifications by user.
func (s *NotificationService) ListNotifications(ctx context.Context, userID uint64, req NotificationListRequest) (*NotificationListResponse, error) {
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
	resp := &NotificationListResponse{
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
			view.IsRead = true
		}
		resp.Items = append(resp.Items, view)
	}
	return resp, nil
}

// MarkNotificationsRead marks notifications as read.
func (s *NotificationService) MarkNotificationsRead(ctx context.Context, userID uint64, ids []uint64) error {
	return s.repo.MarkRead(ctx, userID, ids)
}

// MarkAllNotificationsRead marks all notifications as read.
func (s *NotificationService) MarkAllNotificationsRead(ctx context.Context, userID uint64) error {
	return s.repo.MarkAllRead(ctx, userID)
}

// GetUnreadNotificationCount returns unread notifications count.
func (s *NotificationService) GetUnreadNotificationCount(ctx context.Context, userID uint64) (int64, error) {
	return s.repo.CountUnread(ctx, userID)
}

// DeleteNotification deletes a notification by ID.
func (s *NotificationService) DeleteNotification(ctx context.Context, userID uint64, id uint64) error {
	return s.repo.Delete(ctx, userID, id)
}
