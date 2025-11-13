package feed

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/pkg/safety"
	"gamelink/internal/repository"
	"gamelink/internal/service"
)

const (
	maxFeedImages       = 9
	maxFeedImageSize    = 10 * 1024 * 1024
	maxFeedContentRunes = 1000
	maxReportRunes      = 500
)

// Service orchestrates feed publishing, listing and moderation.
type Service struct {
	repo       repository.FeedRepository
	moderation ModerationEngine
}

// NewService builds a feed service instance.
func NewService(repo repository.FeedRepository, moderation ModerationEngine) *Service {
	if moderation == nil {
		moderation = NewDefaultModerationEngine()
	}
	return &Service{repo: repo, moderation: moderation}
}

// CreateFeedRequest describes payload for creating feeds.
type CreateFeedRequest struct {
	Content    string               `json:"content"`
	Visibility model.FeedVisibility `json:"visibility"`
	Images     []FeedImageInput     `json:"images"`
}

// FeedImageInput describes uploaded image metadata.
type FeedImageInput struct {
	URL       string `json:"url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	SizeBytes int64  `json:"sizeBytes"`
}

// FeedView is a DTO for returning feed information.
type FeedView struct {
	ID               uint64               `json:"id"`
	AuthorID         uint64               `json:"authorId"`
	Content          string               `json:"content"`
	Visibility       model.FeedVisibility `json:"visibility"`
	ModerationStatus string               `json:"moderationStatus"`
	ModerationNote   string               `json:"moderationNote,omitempty"`
	CreatedAt        time.Time            `json:"createdAt"`
	Images           []FeedImageView      `json:"images"`
}

// FeedImageView is serialized feed image.
type FeedImageView struct {
	URL       string `json:"url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	SizeBytes int64  `json:"sizeBytes"`
	Order     int    `json:"order"`
}

// ListFeedsRequest contains filters for timeline fetch.
type ListFeedsRequest struct {
	Cursor string
	Limit  int
}

// ListFeedsResponse returns feed slice with cursor.
type ListFeedsResponse struct {
	Items      []FeedView `json:"items"`
	NextCursor string     `json:"nextCursor,omitempty"`
}

// CreateFeed handles publishing with validation and moderation.
func (s *Service) CreateFeed(ctx context.Context, authorID uint64, req CreateFeedRequest) (*FeedView, error) {
	if err := validateVisibility(req.Visibility); err != nil {
		return nil, fmt.Errorf("%w: %v", service.ErrValidation, err)
	}
	if err := safety.ValidateText(req.Content, maxFeedContentRunes); err != nil {
		return nil, fmt.Errorf("%w: %v", service.ErrValidation, err)
	}
	if len(req.Images) > maxFeedImages {
		return nil, fmt.Errorf("%w: 图片数量超过限制", service.ErrValidation)
	}

	images := make([]model.FeedImage, 0, len(req.Images))
	imageURLs := make([]string, 0, len(req.Images))
	for idx, img := range req.Images {
		url := strings.TrimSpace(img.URL)
		if url == "" {
			return nil, fmt.Errorf("%w: 第 %d 张图片URL为空", service.ErrValidation, idx+1)
		}
		if img.SizeBytes > maxFeedImageSize {
			return nil, fmt.Errorf("%w: 第 %d 张图片超过 10MB", service.ErrValidation, idx+1)
		}
		images = append(images, model.FeedImage{
			URL:       url,
			Width:     img.Width,
			Height:    img.Height,
			SizeBytes: img.SizeBytes,
			Order:     idx,
		})
		imageURLs = append(imageURLs, url)
	}

	feed := &model.Feed{
		AuthorID:   authorID,
		Content:    strings.TrimSpace(req.Content),
		Visibility: req.Visibility,
		Images:     images,
	}
	if err := s.repo.Create(ctx, feed); err != nil {
		return nil, err
	}

	result, err := s.moderation.Evaluate(ctx, ModerationInput{Content: feed.Content, ImageURLs: imageURLs})
	if err != nil {
		return nil, err
	}
	switch result.Decision {
	case ModerationDecisionApprove:
		if err := s.repo.UpdateModeration(ctx, feed.ID, model.FeedModerationApproved, result.Reason, false); err != nil {
			return nil, err
		}
		feed.ModerationStatus = model.FeedModerationApproved
		feed.ModerationNote = result.Reason
	case ModerationDecisionReject:
		if err := s.repo.UpdateModeration(ctx, feed.ID, model.FeedModerationRejected, result.Reason, false); err != nil {
			return nil, err
		}
		feed.ModerationStatus = model.FeedModerationRejected
		feed.ModerationNote = result.Reason
	case ModerationDecisionManual:
		if result.Reason != "" {
			if err := s.repo.UpdateModeration(ctx, feed.ID, model.FeedModerationPending, result.Reason, false); err != nil {
				return nil, err
			}
			feed.ModerationNote = result.Reason
		}
	default:
		// no-op, keep pending
	}

	return toFeedView(feed), nil
}

// ListFeeds returns timeline for user.
func (s *Service) ListFeeds(ctx context.Context, userID uint64, req ListFeedsRequest) (*ListFeedsResponse, error) {
	var cursorValue *uint64
	if req.Cursor != "" {
		parsed, err := strconv.ParseUint(req.Cursor, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: cursor 无效", service.ErrValidation)
		}
		cursorValue = &parsed
	}

	feeds, err := s.repo.List(ctx, repository.FeedListOptions{
		Limit:        req.Limit,
		CursorBefore: cursorValue,
		OnlyApproved: true,
	})
	if err != nil {
		return nil, err
	}

	resp := &ListFeedsResponse{Items: make([]FeedView, 0, len(feeds))}
	for _, f := range feeds {
		feedCopy := f
		resp.Items = append(resp.Items, *toFeedView(&feedCopy))
	}
	if len(feeds) > 0 {
		last := feeds[len(feeds)-1]
		resp.NextCursor = strconv.FormatUint(last.ID, 10)
	}
	return resp, nil
}

// ReportFeed allows users to flag content.
func (s *Service) ReportFeed(ctx context.Context, reporterID, feedID uint64, reason string) error {
	if err := safety.ValidateText(reason, maxReportRunes); err != nil {
		return fmt.Errorf("%w: %v", service.ErrValidation, err)
	}
	if _, err := s.repo.Get(ctx, feedID); err != nil {
		return err
	}
	report := &model.FeedReport{
		FeedID:   feedID,
		Reporter: reporterID,
		Reason:   strings.TrimSpace(reason),
	}
	return s.repo.CreateReport(ctx, report)
}

func validateVisibility(visibility model.FeedVisibility) error {
	switch visibility {
	case "":
		return nil
	case model.FeedVisibilityPublic, model.FeedVisibilityFollowers, model.FeedVisibilityPrivate:
		return nil
	default:
		return fmt.Errorf("visibility 不支持: %s", visibility)
	}
}

func toFeedView(feed *model.Feed) *FeedView {
	images := make([]FeedImageView, 0, len(feed.Images))
	for _, img := range feed.Images {
		images = append(images, FeedImageView{
			URL:       img.URL,
			Width:     img.Width,
			Height:    img.Height,
			SizeBytes: img.SizeBytes,
			Order:     img.Order,
		})
	}
	return &FeedView{
		ID:               feed.ID,
		AuthorID:         feed.AuthorID,
		Content:          feed.Content,
		Visibility:       feed.Visibility,
		ModerationStatus: string(feed.ModerationStatus),
		ModerationNote:   feed.ModerationNote,
		CreatedAt:        feed.CreatedAt,
		Images:           images,
	}
}
