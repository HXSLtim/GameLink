package content

import (
	"context"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// ReportStatus 举报状态
type ReportStatus string

const (
	ReportStatusPending   ReportStatus = "pending"   // 待处理
	ReportStatusProcessed ReportStatus = "processed" // 已处理
	ReportStatusDismissed ReportStatus = "dismissed" // 已驳回
)

// ReportAction 举报处理动作
type ReportAction string

const (
	ReportActionDeleteContent ReportAction = "delete_content" // 删除内容
	ReportActionWarnUser      ReportAction = "warn_user"      // 警告用户
	ReportActionDismiss       ReportAction = "dismiss"        // 驳回举报
)

// FeedReportService 动态举报处理服务
type FeedReportService struct {
	feedRepo  repository.FeedRepository
	opLogRepo repository.OperationLogRepository
}

// NewFeedReportService 创建动态举报处理服务
func NewFeedReportService(
	feedRepo repository.FeedRepository,
	opLogRepo repository.OperationLogRepository,
) *FeedReportService {
	return &FeedReportService{
		feedRepo:  feedRepo,
		opLogRepo: opLogRepo,
	}
}

// FeedReportDTO 动态举报DTO
type FeedReportDTO struct {
	ID           uint64       `json:"id"`
	FeedID       uint64       `json:"feedId"`
	Feed         *FeedSummary `json:"feed,omitempty"`
	ReporterID   uint64       `json:"reporterId"`
	ReporterName string       `json:"reporterName,omitempty"`
	Reason       string       `json:"reason"`
	Status       string       `json:"status"`
	Result       string       `json:"result,omitempty"`
	HandledBy    *uint64      `json:"handledBy,omitempty"`
	HandlerName  string       `json:"handlerName,omitempty"`
	HandledAt    string       `json:"handledAt,omitempty"`
	CreatedAt    string       `json:"createdAt"`
}

// FeedSummary 动态摘要（用于举报详情）
type FeedSummary struct {
	ID               uint64   `json:"id"`
	AuthorID         uint64   `json:"authorId"`
	AuthorName       string   `json:"authorName,omitempty"`
	AuthorAvatar     string   `json:"authorAvatar,omitempty"`
	Content          string   `json:"content"`
	Images           []string `json:"images,omitempty"`
	ModerationStatus string   `json:"moderationStatus"`
	CreatedAt        string   `json:"createdAt"`
}

// ListFeedReportsRequest 列出动态举报请求
type ListFeedReportsRequest struct {
	Page       int        `form:"page"`
	PageSize   int        `form:"pageSize"`
	FeedID     *uint64    `form:"feedId"`
	ReporterID *uint64    `form:"reporterId"`
	Status     *string    `form:"status"`
	DateFrom   *time.Time `form:"dateFrom"`
	DateTo     *time.Time `form:"dateTo"`
}

// ListFeedReportsResponse 列出动态举报响应
type ListFeedReportsResponse struct {
	Items []FeedReportDTO `json:"items"`
	Total int64           `json:"total"`
}

// ProcessReportRequest 处理举报请求
type ProcessReportRequest struct {
	Action ReportAction `json:"action" binding:"required"`
	Result string       `json:"result"`
}

// ListFeedReports 列出动态举报
func (s *FeedReportService) ListFeedReports(ctx context.Context, req ListFeedReportsRequest) (*ListFeedReportsResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	reports, total, err := s.feedRepo.ListReports(ctx, repository.FeedReportListOptions{
		Page:       req.Page,
		PageSize:   req.PageSize,
		FeedID:     req.FeedID,
		ReporterID: req.ReporterID,
		Status:     req.Status,
		DateFrom:   req.DateFrom,
		DateTo:     req.DateTo,
	})
	if err != nil {
		return nil, err
	}

	dtos := make([]FeedReportDTO, 0, len(reports))
	for _, r := range reports {
		// 包含动态内容，方便管理员审核
		dtos = append(dtos, *s.toReportDTOWithFeed(ctx, &r))
	}

	return &ListFeedReportsResponse{
		Items: dtos,
		Total: total,
	}, nil
}

// GetFeedReport 获取举报详情
func (s *FeedReportService) GetFeedReport(ctx context.Context, id uint64) (*FeedReportDTO, error) {
	report, err := s.feedRepo.GetReport(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toReportDTO(report), nil
}

// ProcessFeedReport 处理动态举报
func (s *FeedReportService) ProcessFeedReport(ctx context.Context, reportID uint64, req ProcessReportRequest, handlerID uint64) error {
	report, err := s.feedRepo.GetReport(ctx, reportID)
	if err != nil {
		return err
	}

	now := time.Now()
	report.HandledBy = &handlerID
	report.HandledAt = &now
	report.Result = req.Result

	switch req.Action {
	case ReportActionDeleteContent:
		// 删除被举报的动态
		if err := s.feedRepo.UpdateModeration(ctx, report.FeedID, model.FeedModerationRemoved, "举报处理：删除内容", &handlerID); err != nil {
			return err
		}
		report.Status = string(ReportStatusProcessed)
		s.logOperation(ctx, report.FeedID, handlerID, model.OpActionDelete, "举报处理")

	case ReportActionWarnUser:
		// 警告用户（这里只更新举报状态，实际警告逻辑可扩展）
		report.Status = string(ReportStatusProcessed)
		s.logOperation(ctx, report.FeedID, handlerID, model.OpActionWarnUser, req.Result)

	case ReportActionDismiss:
		// 驳回举报
		report.Status = string(ReportStatusDismissed)
		s.logOperation(ctx, reportID, handlerID, model.OpActionDismissReport, req.Result)

	default:
		return ErrAdminValidation
	}

	return s.feedRepo.UpdateReport(ctx, report)
}

func (s *FeedReportService) logOperation(ctx context.Context, entityID uint64, actorID uint64, action model.OperationAction, reason string) {
	if s.opLogRepo == nil {
		return
	}
	_ = s.opLogRepo.Append(ctx, &model.OperationLog{
		EntityType:  string(model.OpEntityFeedReport),
		EntityID:    entityID,
		ActorUserID: &actorID,
		Action:      string(action),
		Reason:      reason,
	})
}

func (s *FeedReportService) toReportDTO(report *model.FeedReport) *FeedReportDTO {
	dto := &FeedReportDTO{
		ID:         report.ID,
		FeedID:     report.FeedID,
		ReporterID: report.Reporter,
		Reason:     report.Reason,
		Status:     report.Status,
		Result:     report.Result,
		HandledBy:  report.HandledBy,
		CreatedAt:  report.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if report.HandledAt != nil {
		dto.HandledAt = report.HandledAt.Format("2006-01-02 15:04:05")
	}
	return dto
}

// toReportDTOWithFeed 转换为DTO并包含动态内容
func (s *FeedReportService) toReportDTOWithFeed(ctx context.Context, report *model.FeedReport) *FeedReportDTO {
	dto := s.toReportDTO(report)

	// 获取动态内容
	feed, err := s.feedRepo.Get(ctx, report.FeedID)
	if err == nil && feed != nil {
		images := make([]string, 0)
		for _, img := range feed.Images {
			images = append(images, img.URL)
		}
		dto.Feed = &FeedSummary{
			ID:               feed.ID,
			AuthorID:         feed.AuthorID,
			Content:          feed.Content,
			Images:           images,
			ModerationStatus: string(feed.ModerationStatus),
			CreatedAt:        feed.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return dto
}
