package review

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/safety"
)

var (
	// ErrReportNotFound 举报不存在
	ErrReportNotFound = errors.New("report not found")
	// ErrReviewNotFound 评价不存在
	ErrReviewNotFound = errors.New("review not found")
	// ErrInvalidReportAction 无效的举报处理操作
	ErrInvalidReportAction = errors.New("invalid report action")
	// ErrReportAlreadyHandled 举报已处理
	ErrReportAlreadyHandled = errors.New("report already handled")
)

// ReportReviewRequest 举报评价请求
type ReportReviewRequest struct {
	Reason   string `json:"reason" binding:"required,max=500"`
	Evidence string `json:"evidence" binding:"max=1000"`
}

// ReportReviewResponse 举报评价响应
type ReportReviewResponse struct {
	ReportID uint64 `json:"reportId"`
}

// ListReportsRequest 列出举报请求
type ListReportsRequest struct {
	Page       int                       `json:"page"`
	PageSize   int                       `json:"pageSize"`
	ReviewID   *uint64                   `json:"reviewId"`
	ReporterID *uint64                   `json:"reporterId"`
	Status     *model.ReviewReportStatus `json:"status"`
	DateFrom   *time.Time                `json:"dateFrom"`
	DateTo     *time.Time                `json:"dateTo"`
}

// ListReportsResponse 列出举报响应
type ListReportsResponse struct {
	Reports []ReviewReportDTO `json:"reports"`
	Total   int64             `json:"total"`
}

// ReviewReportDTO 举报信息DTO
type ReviewReportDTO struct {
	ID           uint64                   `json:"id"`
	ReviewID     uint64                   `json:"reviewId"`
	ReporterID   uint64                   `json:"reporterId"`
	ReporterName string                   `json:"reporterName"`
	Reason       string                   `json:"reason"`
	Evidence     string                   `json:"evidence,omitempty"`
	Status       model.ReviewReportStatus `json:"status"`
	HandledBy    *uint64                  `json:"handledBy,omitempty"`
	HandlerName  string                   `json:"handlerName,omitempty"`
	HandledAt    *time.Time               `json:"handledAt,omitempty"`
	HandlingNote string                   `json:"handlingNote,omitempty"`
	CreatedAt    time.Time                `json:"createdAt"`
}

// HandleReportRequest 处理举报请求
type HandleReportRequest struct {
	Action string `json:"action" binding:"required,oneof=delete warn reject"`
	Note   string `json:"note" binding:"max=500"`
}

// HandleReportResponse 处理举报响应
type HandleReportResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// ReportReview 举报评价
func (s *ReviewService) ReportReview(ctx context.Context, reviewID, reporterID uint64, req ReportReviewRequest) (*ReportReviewResponse, error) {
	// 验证输入
	if err := safety.ValidateText(req.Reason, 500); err != nil {
		return nil, fmt.Errorf("%w: reason: %v", ErrValidation, err)
	}
	if req.Evidence != "" {
		if err := safety.ValidateText(req.Evidence, 1000); err != nil {
			return nil, fmt.Errorf("%w: evidence: %v", ErrValidation, err)
		}
	}

	// 验证评价是否存在
	review, err := s.reviews.Get(ctx, reviewID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrReviewNotFound
		}
		return nil, err
	}

	// 创建举报记录
	report := &model.ReviewReport{
		ReviewID:   reviewID,
		ReporterID: reporterID,
		Reason:     req.Reason,
		Evidence:   req.Evidence,
		Status:     model.ReviewReportStatusPending,
	}

	if err := s.reports.Create(ctx, report); err != nil {
		return nil, err
	}

	// 标记评价为已举报
	review.IsReported = true
	if err := s.reviews.Update(ctx, review); err != nil {
		// 记录错误但不影响举报创建
		// TODO: 添加日志记录
	}

	return &ReportReviewResponse{
		ReportID: report.ID,
	}, nil
}

// ListReports 列出举报
func (s *ReviewService) ListReports(ctx context.Context, req ListReportsRequest) (*ListReportsResponse, error) {
	// 默认分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 查询举报列表
	reports, total, err := s.reports.List(ctx, repository.ReviewReportListOptions{
		Page:       req.Page,
		PageSize:   req.PageSize,
		ReviewID:   req.ReviewID,
		ReporterID: req.ReporterID,
		Status:     req.Status,
		DateFrom:   req.DateFrom,
		DateTo:     req.DateTo,
	})
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	reportDTOs := make([]ReviewReportDTO, 0, len(reports))
	for _, report := range reports {
		dto := ReviewReportDTO{
			ID:           report.ID,
			ReviewID:     report.ReviewID,
			ReporterID:   report.ReporterID,
			Reason:       report.Reason,
			Evidence:     report.Evidence,
			Status:       report.Status,
			HandledBy:    report.HandledBy,
			HandledAt:    report.HandledAt,
			HandlingNote: report.HandlingNote,
			CreatedAt:    report.CreatedAt,
		}

		// 获取举报人信息
		if reporter, err := s.users.Get(ctx, report.ReporterID); err == nil {
			dto.ReporterName = reporter.Name
		}

		// 获取处理人信息
		if report.HandledBy != nil {
			if handler, err := s.users.Get(ctx, *report.HandledBy); err == nil {
				dto.HandlerName = handler.Name
			}
		}

		reportDTOs = append(reportDTOs, dto)
	}

	return &ListReportsResponse{
		Reports: reportDTOs,
		Total:   total,
	}, nil
}

// HandleReport 处理举报
func (s *ReviewService) HandleReport(ctx context.Context, reportID, handlerID uint64, req HandleReportRequest) (*HandleReportResponse, error) {
	// 验证输入
	if req.Note != "" {
		if err := safety.ValidateText(req.Note, 500); err != nil {
			return nil, fmt.Errorf("%w: note: %v", ErrValidation, err)
		}
	}

	// 获取举报记录
	report, err := s.reports.Get(ctx, reportID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrReportNotFound
		}
		return nil, err
	}

	// 检查举报是否已处理
	if report.Status != model.ReviewReportStatusPending {
		return nil, ErrReportAlreadyHandled
	}

	// 获取被举报的评价
	review, err := s.reviews.Get(ctx, report.ReviewID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrReviewNotFound
		}
		return nil, err
	}

	var message string
	now := time.Now()

	switch req.Action {
	case "delete":
		// 删除评价
		oldStatus := review.Status
		review.Status = model.ReviewStatusDeleted
		if err := s.reviews.Update(ctx, review); err != nil {
			return nil, err
		}

		// 记录删除评价的操作日志
		if s.opLogs != nil {
			metadata := fmt.Sprintf(`{"old_status":"%s","new_status":"%s","report_id":%d,"action":"delete"}`, oldStatus, model.ReviewStatusDeleted, reportID)
			log := &model.OperationLog{
				EntityType:   string(model.OpEntityReview),
				EntityID:     review.ID,
				ActorUserID:  &handlerID,
				Action:       string(model.OpActionDelete),
				Reason:       fmt.Sprintf("处理举报：%s", req.Note),
				MetadataJSON: []byte(metadata),
			}
			_ = s.opLogs.Append(ctx, log)
		}

		// 更新举报状态为已通过
		report.Status = model.ReviewReportStatusApproved
		report.HandledBy = &handlerID
		report.HandledAt = &now
		report.HandlingNote = req.Note
		if report.HandlingNote == "" {
			report.HandlingNote = "评价已删除"
		}
		message = "评价已删除"

	case "warn":
		// 警告评价者（保留评价，但标记为已处理）
		// TODO: 发送警告通知给评价者
		report.Status = model.ReviewReportStatusApproved
		report.HandledBy = &handlerID
		report.HandledAt = &now
		report.HandlingNote = req.Note
		if report.HandlingNote == "" {
			report.HandlingNote = "已警告评价者"
		}
		message = "已警告评价者"

	case "reject":
		// 驳回举报
		report.Status = model.ReviewReportStatusRejected
		report.HandledBy = &handlerID
		report.HandledAt = &now
		report.HandlingNote = req.Note
		if report.HandlingNote == "" {
			report.HandlingNote = "举报不成立"
		}

		// 如果没有其他待处理的举报，取消评价的举报标记
		pendingStatus := model.ReviewReportStatusPending
		reviewIDPtr := &report.ReviewID
		pendingReports, _, err := s.reports.List(ctx, repository.ReviewReportListOptions{
			ReviewID: reviewIDPtr,
			Status:   &pendingStatus,
			Page:     1,
			PageSize: 1,
		})
		if err == nil && len(pendingReports) == 0 {
			review.IsReported = false
			if err := s.reviews.Update(ctx, review); err != nil {
				// 记录错误但不影响举报处理
			}
		}
		message = "举报已驳回"

	default:
		return nil, ErrInvalidReportAction
	}

	// 更新举报记录
	if err := s.reports.Update(ctx, report); err != nil {
		return nil, err
	}

	// 记录处理举报的操作日志
	if s.opLogs != nil {
		metadata := fmt.Sprintf(`{"report_id":%d,"action":"%s","note":"%s"}`, reportID, req.Action, req.Note)
		log := &model.OperationLog{
			EntityType:   string(model.OpEntityReview),
			EntityID:     review.ID,
			ActorUserID:  &handlerID,
			Action:       string(model.OpActionHandleReport),
			Reason:       fmt.Sprintf("处理举报：%s", message),
			MetadataJSON: []byte(metadata),
		}
		_ = s.opLogs.Append(ctx, log)
	}

	return &HandleReportResponse{
		Status:  "success",
		Message: message,
	}, nil
}
