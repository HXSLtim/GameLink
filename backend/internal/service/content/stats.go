package content

import (
	"context"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// ContentStatsService 内容统计服务
type ContentStatsService struct {
	feedRepo    repository.FeedRepository
	messageRepo repository.ChatMessageRepository
}

// NewContentStatsService 创建内容统计服务
func NewContentStatsService(
	feedRepo repository.FeedRepository,
	messageRepo repository.ChatMessageRepository,
) *ContentStatsService {
	return &ContentStatsService{
		feedRepo:    feedRepo,
		messageRepo: messageRepo,
	}
}

// ContentStatsDTO 内容统计DTO
type ContentStatsDTO struct {
	TotalFeeds    int64                                `json:"totalFeeds"`
	PendingFeeds  int64                                `json:"pendingFeeds"`
	ApprovedFeeds int64                                `json:"approvedFeeds"`
	RejectedFeeds int64                                `json:"rejectedFeeds"`
	TotalMessages int64                                `json:"totalMessages"`
	FeedsByStatus map[model.FeedModerationStatus]int64 `json:"feedsByStatus"`
	FeedTrend     []repository.DateValue               `json:"feedTrend"`
}

// GetStats 获取内容统计
func (s *ContentStatsService) GetStats(ctx context.Context, trendDays int) (*ContentStatsDTO, error) {
	if trendDays < 1 {
		trendDays = 30
	}

	// 获取动态状态统计
	feedsByStatus, err := s.feedRepo.CountByStatus(ctx)
	if err != nil {
		return nil, err
	}

	// 计算总数
	var totalFeeds int64
	for _, count := range feedsByStatus {
		totalFeeds += count
	}

	// 获取动态趋势
	feedTrend, err := s.feedRepo.GetTrend(ctx, trendDays)
	if err != nil {
		feedTrend = []repository.DateValue{}
	}

	// 获取消息统计（简化版，可扩展）
	messages, totalMessages, _ := s.messageRepo.ListForModeration(ctx, repository.ChatMessageModerationListOptions{
		Page:     1,
		PageSize: 1,
	})
	_ = messages // 只需要总数

	return &ContentStatsDTO{
		TotalFeeds:    totalFeeds,
		PendingFeeds:  feedsByStatus[model.FeedModerationPending],
		ApprovedFeeds: feedsByStatus[model.FeedModerationApproved],
		RejectedFeeds: feedsByStatus[model.FeedModerationRejected],
		TotalMessages: totalMessages,
		FeedsByStatus: feedsByStatus,
		FeedTrend:     feedTrend,
	}, nil
}
