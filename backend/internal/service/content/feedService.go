package content

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gamelink/pkg/apierr"
	"gamelink/internal/model"
	"gamelink/pkg/safety"
	"gamelink/internal/repository"
)

const (
	maxFeedImages       = 9
	maxFeedImageSize    = 10 * 1024 * 1024
	maxFeedContentRunes = 1000
	maxReportRunes      = 500
)

// FeedService orchestrates feed publishing, listing and moderation.
type FeedService struct {
	repo       repository.FeedRepository
	moderation FeedModerationEngine
}

// NewFeedService builds a feed service instance.
func NewFeedService(repo repository.FeedRepository, moderation FeedModerationEngine) *FeedService {
	if moderation == nil {
		moderation = NewDefaultFeedModerationEngine()
	}
	return &FeedService{repo: repo, moderation: moderation}
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
func (s *FeedService) CreateFeed(ctx context.Context, authorID uint64, req CreateFeedRequest) (*FeedView, error) {
	if err := validateFeedVisibility(req.Visibility); err != nil {
		return nil, apierr.BadRequest("动态可见性验证失败: " + err.Error())
	}
	if err := safety.ValidateText(req.Content, maxFeedContentRunes); err != nil {
		return nil, apierr.BadRequest("动态内容验证失败: " + err.Error())
	}
	if len(req.Images) > maxFeedImages {
		return nil, apierr.BadRequest("图片数量超过限制")
	}

	images := make([]model.FeedImage, 0, len(req.Images))
	imageURLs := make([]string, 0, len(req.Images))
	for idx, img := range req.Images {
		url := strings.TrimSpace(img.URL)
		if url == "" {
			return nil, apierr.BadRequest(fmt.Sprintf("第%d张图片URL为空", idx+1))
		}
		if img.SizeBytes > maxFeedImageSize {
			return nil, apierr.BadRequest(fmt.Sprintf("第%d张图片超过10MB", idx+1))
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

	result, err := s.moderation.Evaluate(ctx, FeedModerationInput{Content: feed.Content, ImageURLs: imageURLs})
	if err != nil {
		return nil, err
	}
	switch result.Decision {
	case FeedModerationDecisionApprove:
		if err := s.repo.UpdateModeration(ctx, feed.ID, model.FeedModerationApproved, result.Reason, false); err != nil {
			return nil, err
		}
		feed.ModerationStatus = model.FeedModerationApproved
		feed.ModerationNote = result.Reason
	case FeedModerationDecisionReject:
		if err := s.repo.UpdateModeration(ctx, feed.ID, model.FeedModerationRejected, result.Reason, false); err != nil {
			return nil, err
		}
		feed.ModerationStatus = model.FeedModerationRejected
		feed.ModerationNote = result.Reason
	case FeedModerationDecisionManual:
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
func (s *FeedService) ListFeeds(ctx context.Context, userID uint64, req ListFeedsRequest) (*ListFeedsResponse, error) {
	var cursorValue *uint64
	if req.Cursor != "" {
		parsed, err := strconv.ParseUint(req.Cursor, 10, 64)
		if err != nil {
			return nil, apierr.BadRequest("cursor 无效: " + err.Error())
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
func (s *FeedService) ReportFeed(ctx context.Context, reporterID, feedID uint64, reason string) error {
	if err := safety.ValidateText(reason, maxReportRunes); err != nil {
		return apierr.BadRequest("验证失败: " + err.Error())
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

func validateFeedVisibility(visibility model.FeedVisibility) error {
	switch visibility {
	case "":
		return nil
	case model.FeedVisibilityPublic, model.FeedVisibilityFollowers, model.FeedVisibilityPrivate:
		return nil
	default:
		return apierr.BadRequest("visibility 不支持: " + string(visibility))
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