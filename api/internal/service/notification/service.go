package notification

import (
	"context"
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
