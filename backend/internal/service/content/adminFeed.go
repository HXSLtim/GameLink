package content

import (
	"context"
	"errors"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service/sensitiveword"
)

var (
	// ErrAdminNotFound 动态不存在
	ErrAdminNotFound = repository.ErrNotFound
	// ErrAdminValidation 输入校验失败
	ErrAdminValidation = errors.New("validation failed")
)

// AdminFeedService 动态管理服务（管理员）
type AdminFeedService struct {
	feedRepo      repository.FeedRepository
	sensitiveWord *sensitiveword.SensitiveWordService
	opLogRepo     repository.OperationLogRepository
}

// NewAdminFeedService 创建动态管理服务
func NewAdminFeedService(
	feedRepo repository.FeedRepository,
	sensitiveWord *sensitiveword.SensitiveWordService,
	opLogRepo repository.OperationLogRepository,
) *AdminFeedService {
	return &AdminFeedService{
		feedRepo:      feedRepo,
		sensitiveWord: sensitiveWord,
		opLogRepo:     opLogRepo,
	}
}

// AdminFeedDTO 动态DTO（管理员视图）
type AdminFeedDTO struct {
	ID                 uint64                     `json:"id"`
	AuthorID           uint64                     `json:"authorId"`
	Content            string                     `json:"content"`
	CategoryID         *uint64                    `json:"categoryId,omitempty"`
	CategoryName       string                     `json:"categoryName,omitempty"`
	Visibility         model.FeedVisibility       `json:"visibility"`
	ModerationStatus   model.FeedModerationStatus `json:"moderationStatus"`
	ModerationNote     string                     `json:"moderationNote,omitempty"`
	Images             []AdminFeedImageDTO        `json:"images"`
	Metrics            model.FeedMetricFields     `json:"metrics"`
	CreatedAt          string                     `json:"createdAt"`
	UpdatedAt          string                     `json:"updatedAt"`
	HasSensitiveWords  bool                       `json:"hasSensitiveWords,omitempty"`
	HighlightedContent string                     `json:"highlightedContent,omitempty"`
}

// AdminFeedImageDTO 动态图片DTO
type AdminFeedImageDTO struct {
	ID     uint64 `json:"id"`
	URL    string `json:"url"`
	Order  int    `json:"order"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// AdminListFeedsRequest 列出动态请求
type AdminListFeedsRequest struct {
	Page             int                         `form:"page"`
	PageSize         int                         `form:"pageSize"`
	AuthorID         *uint64                     `form:"authorId"`
	CategoryID       *uint64                     `form:"categoryId"`
	Keyword          string                      `form:"keyword"`
	ModerationStatus *model.FeedModerationStatus `form:"moderationStatus"`
	DateFrom         *time.Time                  `form:"dateFrom"`
	DateTo           *time.Time                  `form:"dateTo"`
}

// AdminListFeedsResponse 列出动态响应
type AdminListFeedsResponse struct {
	Feeds []AdminFeedDTO `json:"feeds"`
	Total int64          `json:"total"`
}

// AdminModerationRequest 审核请求
type AdminModerationRequest struct {
	Note string `json:"note"`
}

// AdminBatchModerationRequest 批量审核请求
type AdminBatchModerationRequest struct {
	FeedIDs []uint64 `json:"feedIds" binding:"required,min=1"`
	Note    string   `json:"note"`
}

// ListFeeds 列出动态（管理员）
func (s *AdminFeedService) ListFeeds(ctx context.Context, req AdminListFeedsRequest) (*AdminListFeedsResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	feeds, total, err := s.feedRepo.ListPaged(ctx, repository.FeedPagedListOptions{
		Page:             req.Page,
		PageSize:         req.PageSize,
		AuthorID:         req.AuthorID,
		CategoryID:       req.CategoryID,
		Keyword:          req.Keyword,
		ModerationStatus: req.ModerationStatus,
		DateFrom:         req.DateFrom,
		DateTo:           req.DateTo,
	})
	if err != nil {
		return nil, err
	}

	dtos := make([]AdminFeedDTO, 0, len(feeds))
	for _, f := range feeds {
		dto := s.toDTO(&f)
		if s.sensitiveWord != nil {
			result, _ := s.sensitiveWord.DetectSensitiveWords(ctx, sensitiveword.DetectSensitiveWordsRequest{
				Content: f.Content,
			})
			if result != nil {
				dto.HasSensitiveWords = result.HasSensitiveWords
				dto.HighlightedContent = result.HighlightedContent
			}
		}
		dtos = append(dtos, *dto)
	}

	return &AdminListFeedsResponse{Feeds: dtos, Total: total}, nil
}

// GetFeed 获取动态详情
func (s *AdminFeedService) GetFeed(ctx context.Context, id uint64) (*AdminFeedDTO, error) {
	feed, err := s.feedRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	dto := s.toDTO(feed)
	if s.sensitiveWord != nil {
		result, _ := s.sensitiveWord.DetectSensitiveWords(ctx, sensitiveword.DetectSensitiveWordsRequest{
			Content: feed.Content,
		})
		if result != nil {
			dto.HasSensitiveWords = result.HasSensitiveWords
			dto.HighlightedContent = result.HighlightedContent
		}
	}
	return dto, nil
}

// ApproveFeed 批准动态
func (s *AdminFeedService) ApproveFeed(ctx context.Context, id uint64, moderatorID uint64, note string) error {
	if err := s.feedRepo.UpdateModeration(ctx, id, model.FeedModerationApproved, note, &moderatorID); err != nil {
		return err
	}
	s.logOperation(ctx, id, moderatorID, model.OpActionApprove, note)
	return nil
}

// RejectFeed 拒绝动态
func (s *AdminFeedService) RejectFeed(ctx context.Context, id uint64, moderatorID uint64, note string) error {
	if note == "" {
		return ErrAdminValidation
	}
	if err := s.feedRepo.UpdateModeration(ctx, id, model.FeedModerationRejected, note, &moderatorID); err != nil {
		return err
	}
	s.logOperation(ctx, id, moderatorID, model.OpActionReject, note)
	return nil
}

// DeleteFeed 删除动态
func (s *AdminFeedService) DeleteFeed(ctx context.Context, id uint64, moderatorID uint64, reason string) error {
	if err := s.feedRepo.UpdateModeration(ctx, id, model.FeedModerationRemoved, reason, &moderatorID); err != nil {
		return err
	}
	s.logOperation(ctx, id, moderatorID, model.OpActionDelete, reason)
	return nil
}

// BatchApproveFeed 批量批准动态
func (s *AdminFeedService) BatchApproveFeed(ctx context.Context, feedIDs []uint64, moderatorID uint64, note string) error {
	if len(feedIDs) == 0 {
		return ErrAdminValidation
	}
	if err := s.feedRepo.BatchUpdateModeration(ctx, feedIDs, model.FeedModerationApproved, note, &moderatorID); err != nil {
		return err
	}
	for _, id := range feedIDs {
		s.logOperation(ctx, id, moderatorID, model.OpActionBatchApprove, note)
	}
	return nil
}

// BatchRejectFeed 批量拒绝动态
func (s *AdminFeedService) BatchRejectFeed(ctx context.Context, feedIDs []uint64, moderatorID uint64, note string) error {
	if len(feedIDs) == 0 || note == "" {
		return ErrAdminValidation
	}
	if err := s.feedRepo.BatchUpdateModeration(ctx, feedIDs, model.FeedModerationRejected, note, &moderatorID); err != nil {
		return err
	}
	for _, id := range feedIDs {
		s.logOperation(ctx, id, moderatorID, model.OpActionBatchReject, note)
	}
	return nil
}

func (s *AdminFeedService) logOperation(ctx context.Context, feedID uint64, actorID uint64, action model.OperationAction, reason string) {
	if s.opLogRepo == nil {
		return
	}
	_ = s.opLogRepo.Append(ctx, &model.OperationLog{
		EntityType:  string(model.OpEntityFeed),
		EntityID:    feedID,
		ActorUserID: &actorID,
		Action:      string(action),
		Reason:      reason,
	})
}

func (s *AdminFeedService) toDTO(feed *model.Feed) *AdminFeedDTO {
	dto := &AdminFeedDTO{
		ID:               feed.ID,
		AuthorID:         feed.AuthorID,
		Content:          feed.Content,
		CategoryID:       feed.CategoryID,
		Visibility:       feed.Visibility,
		ModerationStatus: feed.ModerationStatus,
		ModerationNote:   feed.ModerationNote,
		Metrics:          feed.Metrics,
		Images:           make([]AdminFeedImageDTO, 0, len(feed.Images)),
		CreatedAt:        feed.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        feed.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	if feed.Category != nil {
		dto.CategoryName = feed.Category.Name
	}
	for _, img := range feed.Images {
		dto.Images = append(dto.Images, AdminFeedImageDTO{
			ID: img.ID, URL: img.URL, Order: img.Order, Width: img.Width, Height: img.Height,
		})
	}
	return dto
}
