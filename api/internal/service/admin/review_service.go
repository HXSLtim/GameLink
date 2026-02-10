package admin

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	
	"log/slog"
	
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/common"
	feedservice "gamelink/internal/service/feed"
	"gamelink/pkg/apierr"
)

// --- Review management ---

// ListReviews 列出评价。
func (s *AdminService) ListReviews(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, *model.Pagination, error) {
	if s.tx == nil {
		return nil, nil, apierr.InternalError("事务管理器未配置")
	}
	page := repository.NormalizePage(opts.Page)
	size := repository.NormalizePageSize(opts.PageSize)
	items, total, err := s.repos().Reviews.List(ctx, repository.ReviewListOptions{
		Page: page, PageSize: size, OrderID: opts.OrderID, UserID: opts.UserID, PlayerID: opts.PlayerID, DateFrom: opts.DateFrom, DateTo: opts.DateTo,
	})
	if err != nil {
		return nil, nil, err
	}
	p := buildPagination(page, size, total)
	return items, &p, nil
}

// GetReview 返回评价详情。
func (s *AdminService) GetReview(ctx context.Context, id uint64) (*model.Review, error) {
	if s.tx == nil {
		return nil, apierr.InternalError("事务管理器未配置")
	}
	item, err := s.repos().Reviews.Get(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get review")
	}
	return item, nil
}

// CreateReview 新建评价。
func (s *AdminService) CreateReview(ctx context.Context, r model.Review) (*model.Review, error) {
	if !r.Score.Valid() || r.OrderID == 0 || r.UserID == 0 || r.PlayerID == 0 {
		return nil, ErrValidation
	}
	if s.tx == nil {
		return nil, apierr.InternalError("事务管理器未配置")
	}
	err := s.tx.WithTx(ctx, func(txr *common.Repos) error { return txr.Reviews.Create(ctx, &r) })
	if err != nil {
		return nil, WrapError(err, "create review")
	}
	s.appendLogAsync(ctx, string(model.OpEntityReview), r.ID, string(model.OpActionCreate), map[string]any{"order_id": r.OrderID, "player_id": r.PlayerID})
	return &r, nil
}

// UpdateReview 修改评价分数/内容。
func (s *AdminService) UpdateReview(ctx context.Context, id uint64, score model.Rating, content string) (*model.Review, error) {
	if !score.Valid() {
		return nil, ErrValidation
	}
	if s.tx == nil {
		return nil, apierr.InternalError("事务管理器未配置")
	}
	var item *model.Review
	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		obj, err := r.Reviews.Get(ctx, id)
		if err != nil {
			return WrapError(err, "get review")
		}
		obj.Score = score
		obj.Content = strings.TrimSpace(content)
		if err := r.Reviews.Update(ctx, obj); err != nil {
			return WrapError(err, "update review")
		}
		item = obj
		return nil
	})
	if err != nil {
		return nil, WrapError(err, "update review transaction")
	}
	s.appendLogAsync(ctx, string(model.OpEntityReview), id, string(model.OpActionUpdate), nil)
	return item, nil
}

// DeleteReview 删除评价。
func (s *AdminService) DeleteReview(ctx context.Context, id uint64) error {
	if s.tx == nil {
		return errors.New("transaction manager not configured")
	}
	err := s.tx.WithTx(ctx, func(r *common.Repos) error { return r.Reviews.Delete(ctx, id) })
	return WrapError(err, "delete review")
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

// ReportReview 举报评价
func (s *AdminService) ReportReview(ctx context.Context, reviewID, reporterID uint64, reason, evidence string) (uint64, error) {
	if s.tx == nil {
		return 0, apierr.InternalError("事务管理器未配置")
	}

	var reportID uint64
	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		// 验证评价是否存在
		_, err := r.Reviews.Get(ctx, reviewID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}
			return err
		}

		// 创建举报记录
		report := &model.ReviewReport{
			ReviewID:   reviewID,
			ReporterID: reporterID,
			Reason:     reason,
			Evidence:   evidence,
			Status:     model.ReviewReportStatusPending,
		}

		if err := r.ReviewReports.Create(ctx, report); err != nil {
			return err
		}

		reportID = report.ID

		// 标记评价为已举报
		review, err := r.Reviews.Get(ctx, reviewID)
		if err != nil {
			return err
		}
		review.IsReported = true
		return r.Reviews.Update(ctx, review)
	})

	if err != nil {
		return 0, WrapError(err, "report review")
	}

	return reportID, nil
}

// ListReviewReports 列出举报
func (s *AdminService) ListReviewReports(ctx context.Context, page, pageSize int, reviewID, reporterID *uint64, status *model.ReviewReportStatus, dateFrom, dateTo *time.Time) ([]ReviewReportDTO, *model.Pagination, error) {
	if s.tx == nil {
		return nil, nil, apierr.InternalError("事务管理器未配置")
	}
	reports, total, err := s.repos().ReviewReports.List(ctx, repository.ReviewReportListOptions{
		Page:       page,
		PageSize:   pageSize,
		ReviewID:   reviewID,
		ReporterID: reporterID,
		Status:     status,
		DateFrom:   dateFrom,
		DateTo:     dateTo,
	})
	if err != nil {
		return nil, nil, WrapError(err, "list review reports")
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

	p := &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	}

	return reportDTOs, p, nil
}

// GetReviewReport 获取举报详情
func (s *AdminService) GetReviewReport(ctx context.Context, id uint64) (*ReviewReportDTO, error) {
	if s.tx == nil {
		return nil, apierr.InternalError("事务管理器未配置")
	}
	report, err := s.repos().ReviewReports.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, WrapError(err, "get review report")
	}

	dto := &ReviewReportDTO{
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

	return dto, nil
}

// HandleReportResponse 处理举报响应
type HandleReportResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// HandleReviewReport 处理举报
func (s *AdminService) HandleReviewReport(ctx context.Context, reportID, handlerID uint64, action, note string) (*HandleReportResponse, error) {
	if s.tx == nil {
		return nil, apierr.InternalError("事务管理器未配置")
	}

	var message string
	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		// 获取举报记录
		report, err := r.ReviewReports.Get(ctx, reportID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}
			return err
		}

		// 检查举报是否已处理
		if report.Status != model.ReviewReportStatusPending {
			return apierr.BadRequest("report already handled")
		}

		// 获取被举报的评价
		review, err := r.Reviews.Get(ctx, report.ReviewID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}
			return err
		}

		now := time.Now()

		switch action {
		case "delete":
			// 删除评价
			oldStatus := review.Status
			review.Status = model.ReviewStatusDeleted
			if err := r.Reviews.Update(ctx, review); err != nil {
				return err
			}

			// 记录删除评价的操作日志
			if r.OpLogs != nil {
				metadata := fmt.Sprintf(`{"old_status":"%s","new_status":"%s","report_id":%d,"action":"delete"}`, oldStatus, model.ReviewStatusDeleted, reportID)
				log := &model.OperationLog{
					EntityType:   string(model.OpEntityReview),
					EntityID:     review.ID,
					ActorUserID:  &handlerID,
					Action:       string(model.OpActionDelete),
					Reason:       fmt.Sprintf("处理举报：%s", note),
					MetadataJSON: []byte(metadata),
				}
				_ = r.OpLogs.Append(ctx, log)
			}

			// 更新举报状态为已通过
			report.Status = model.ReviewReportStatusApproved
			report.HandledBy = &handlerID
			report.HandledAt = &now
			report.HandlingNote = note
			if report.HandlingNote == "" {
				report.HandlingNote = "评价已删除"
			}
			message = "评价已删除"

		case "warn":
			// 警告评价者（保留评价，但标记为已处理）
			report.Status = model.ReviewReportStatusApproved
			report.HandledBy = &handlerID
			report.HandledAt = &now
			report.HandlingNote = note
			if report.HandlingNote == "" {
				report.HandlingNote = "已警告评价者"
			}
			message = "已警告评价者"

		case "reject":
			// 驳回举报
			report.Status = model.ReviewReportStatusRejected
			report.HandledBy = &handlerID
			report.HandledAt = &now
			report.HandlingNote = note
			if report.HandlingNote == "" {
				report.HandlingNote = "举报不成立"
			}
			message = "举报已驳回"

		default:
			return apierr.BadRequest("invalid action")
		}

		// 更新举报记录
		if err := r.ReviewReports.Update(ctx, report); err != nil {
			return err
		}

		// 检查是否还有其他待处理的举报，如果没有则取消评价的举报标记
		pendingStatus := model.ReviewReportStatusPending
		reviewIDPtr := &report.ReviewID
		pendingReports, _, err := r.ReviewReports.List(ctx, repository.ReviewReportListOptions{
			ReviewID: reviewIDPtr,
			Status:   &pendingStatus,
			Page:     1,
			PageSize: 1,
		})
		if err == nil && len(pendingReports) == 0 && review.IsReported {
			review.IsReported = false
			if err := r.Reviews.Update(ctx, review); err != nil {
				// 记录错误但不影响举报处理
				slog.Warn("failed to update review reported status", slog.Any("error", err))
			}
		}

		// 记录处理举报的操作日志
		if r.OpLogs != nil {
			metadata := fmt.Sprintf(`{"report_id":%d,"action":"%s","note":"%s"}`, reportID, action, note)
			log := &model.OperationLog{
				EntityType:   string(model.OpEntityReview),
				EntityID:     review.ID,
				ActorUserID:  &handlerID,
				Action:       string(model.OpActionHandleReport),
				Reason:       fmt.Sprintf("处理举报：%s", message),
				MetadataJSON: []byte(metadata),
			}
			_ = r.OpLogs.Append(ctx, log)
		}

		return nil
	})

	if err != nil {
		return nil, WrapError(err, "handle review report")
	}

	return &HandleReportResponse{
		Status:  "success",
		Message: message,
	}, nil
}

// ListPendingReviews 获取待审核评价列表
func (s *AdminService) ListPendingReviews(ctx context.Context, page, pageSize int) ([]model.Review, int64, error) {
	if s.tx == nil {
		return nil, 0, apierr.InternalError("事务管理器未配置")
	}
	reviews, total, err := s.repos().Reviews.ListPending(ctx, page, pageSize)
	if err != nil {
		return nil, 0, WrapError(err, "list pending reviews")
	}
	return reviews, total, nil
}

// ApproveReview 批准评价
func (s *AdminService) ApproveReview(ctx context.Context, reviewID uint64, reason string, actorUserID *uint64) error {
	if s.tx == nil {
		return apierr.InternalError("事务管理器未配置")
	}

	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		// 获取评价
		review, err := r.Reviews.Get(ctx, reviewID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}
			return err
		}

		// 检查状态：只有待审核的评价可以批准
		if review.Status != model.ReviewStatusPending {
			return apierr.BadRequest("只能批准待审核的评价")
		}

		oldStatus := review.Status

		// 更新状态为已通过
		if err := r.Reviews.UpdateStatus(ctx, reviewID, model.ReviewStatusApproved, ""); err != nil {
			return err
		}

		// 记录操作日志
		if r.OpLogs != nil {
			metadata := fmt.Sprintf(`{"old_status":"%s","new_status":"%s"}`, oldStatus, model.ReviewStatusApproved)
			log := &model.OperationLog{
				EntityType:   string(model.OpEntityReview),
				EntityID:     reviewID,
				ActorUserID:  actorUserID,
				Action:       string(model.OpActionApprove),
				Reason:       reason,
				MetadataJSON: []byte(metadata),
			}
			_ = r.OpLogs.Append(ctx, log)
		}

		return nil
	})

	if err != nil {
		return WrapError(err, "approve review")
	}

	return nil
}

// RejectReview 拒绝评价
func (s *AdminService) RejectReview(ctx context.Context, reviewID uint64, reason string, actorUserID *uint64) error {
	if s.tx == nil {
		return apierr.InternalError("事务管理器未配置")
	}

	// 验证拒绝原因
	if reason == "" {
		return apierr.BadRequest("拒绝原因不能为空")
	}

	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		// 获取评价
		review, err := r.Reviews.Get(ctx, reviewID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}
			return err
		}

		// 检查状态：只有待审核的评价可以拒绝
		if review.Status != model.ReviewStatusPending {
			return apierr.BadRequest("只能拒绝待审核的评价")
		}

		oldStatus := review.Status

		// 更新状态为已拒绝
		if err := r.Reviews.UpdateStatus(ctx, reviewID, model.ReviewStatusRejected, reason); err != nil {
			return err
		}

		// 记录操作日志
		if r.OpLogs != nil {
			metadata := fmt.Sprintf(`{"old_status":"%s","new_status":"%s","rejection_reason":"%s"}`, oldStatus, model.ReviewStatusRejected, reason)
			log := &model.OperationLog{
				EntityType:   string(model.OpEntityReview),
				EntityID:     reviewID,
				ActorUserID:  actorUserID,
				Action:       string(model.OpActionReject),
				Reason:       reason,
				MetadataJSON: []byte(metadata),
			}
			_ = r.OpLogs.Append(ctx, log)
		}

		return nil
	})

	if err != nil {
		return WrapError(err, "reject review")
	}

	return nil
}

// BatchApproveReviews 批量批准评价
func (s *AdminService) BatchApproveReviews(ctx context.Context, reviewIDs []uint64, actorUserID *uint64) error {
	if s.tx == nil {
		return apierr.InternalError("事务管理器未配置")
	}

	if len(reviewIDs) == 0 {
		return apierr.BadRequest("评价ID列表不能为空")
	}

	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		// 验证所有评价都是待审核状态
		for _, id := range reviewIDs {
			review, err := r.Reviews.Get(ctx, id)
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					return apierr.NotFound("评价不存在: " + string(rune(id)))
				}
				return err
			}
			if review.Status != model.ReviewStatusPending {
				return apierr.BadRequest("评价不是待审核状态: " + string(rune(id)))
			}
		}

		// 批量更新状态
		if err := r.Reviews.BatchUpdateStatus(ctx, reviewIDs, model.ReviewStatusApproved, ""); err != nil {
			return err
		}

		// 记录操作日志
		if r.OpLogs != nil {
			for _, id := range reviewIDs {
				metadata := fmt.Sprintf(`{"old_status":"%s","new_status":"%s","batch":true}`, model.ReviewStatusPending, model.ReviewStatusApproved)
				log := &model.OperationLog{
					EntityType:   string(model.OpEntityReview),
					EntityID:     id,
					ActorUserID:  actorUserID,
					Action:       string(model.OpActionApprove),
					Reason:       "批量批准评价",
					MetadataJSON: []byte(metadata),
				}
				_ = r.OpLogs.Append(ctx, log)
			}
		}

		return nil
	})

	if err != nil {
		return WrapError(err, "batch approve reviews")
	}

	return nil
}

// BatchRejectReviews 批量拒绝评价
func (s *AdminService) BatchRejectReviews(ctx context.Context, reviewIDs []uint64, reason string, actorUserID *uint64) error {
	if s.tx == nil {
		return apierr.InternalError("事务管理器未配置")
	}

	if len(reviewIDs) == 0 {
		return apierr.BadRequest("评价ID列表不能为空")
	}

	if reason == "" {
		return apierr.BadRequest("拒绝原因不能为空")
	}

	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		// 验证所有评价都是待审核状态
		for _, id := range reviewIDs {
			review, err := r.Reviews.Get(ctx, id)
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					return apierr.NotFound("评价不存在: " + string(rune(id)))
				}
				return err
			}
			if review.Status != model.ReviewStatusPending {
				return apierr.BadRequest("评价不是待审核状态: " + string(rune(id)))
			}
		}

		// 批量更新状态
		if err := r.Reviews.BatchUpdateStatus(ctx, reviewIDs, model.ReviewStatusRejected, reason); err != nil {
			return err
		}

		// 记录操作日志
		if r.OpLogs != nil {
			for _, id := range reviewIDs {
				metadata := fmt.Sprintf(`{"old_status":"%s","new_status":"%s","rejection_reason":"%s","batch":true}`, model.ReviewStatusPending, model.ReviewStatusRejected, reason)
				log := &model.OperationLog{
					EntityType:   string(model.OpEntityReview),
					EntityID:     id,
					ActorUserID:  actorUserID,
					Action:       string(model.OpActionReject),
					Reason:       reason,
					MetadataJSON: []byte(metadata),
				}
				_ = r.OpLogs.Append(ctx, log)
			}
		}

		return nil
	})

	if err != nil {
		return WrapError(err, "batch reject reviews")
	}

	return nil
}

// UpdateReviewReply 更新评价回复
func (s *AdminService) UpdateReviewReply(ctx context.Context, userID, replyID uint64, content string) (map[string]interface{}, error) {
	if s.tx == nil {
		return nil, apierr.InternalError("事务管理器未配置")
	}

	var result map[string]interface{}
	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		// 获取回复
		reply, err := r.ReviewReplies.Get(ctx, replyID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}
			return err
		}

		// 权限检查：只能更新自己的回复
		if reply.AuthorID != userID {
			return ErrUnauthorized
		}

		oldContent := reply.Content

		// 更新回复内容
		reply.Content = strings.TrimSpace(content)

		// 重新进行内容审核
		engine := feedservice.NewDefaultModerationEngine()
		moderationResult, err := engine.Evaluate(ctx, feedservice.ModerationInput{Content: reply.Content})
		if err != nil {
			return err
		}

		var status model.ReviewReplyStatus
		note := moderationResult.Reason
		switch moderationResult.Decision {
		case feedservice.ModerationDecisionApprove:
			status = model.ReviewReplyStatusApproved
		case feedservice.ModerationDecisionReject:
			status = model.ReviewReplyStatusRejected
		case feedservice.ModerationDecisionManual:
			status = model.ReviewReplyStatusPending
		default:
			status = model.ReviewReplyStatusPending
		}

		reply.Status = status
		reply.ModerationNote = note

		// 更新回复
		if err := r.ReviewReplies.Update(ctx, reply); err != nil {
			return err
		}

		// 记录操作日志
		if r.OpLogs != nil {
			metadata := fmt.Sprintf(`{"reply_id":%d,"old_content":"%s","new_content":"%s","status":"%s"}`, replyID, oldContent, reply.Content, status)
			log := &model.OperationLog{
				EntityType:   string(model.OpEntityReview),
				EntityID:     reply.ReviewID,
				ActorUserID:  &userID,
				Action:       string(model.OpActionUpdateReply),
				Reason:       "更新回复",
				MetadataJSON: []byte(metadata),
			}
			_ = r.OpLogs.Append(ctx, log)
		}

		// 发送通知给评价者
		review, err := r.Reviews.Get(ctx, reply.ReviewID)
		if err == nil && review.UserID != userID {
			notification := &model.NotificationEvent{
				UserID:        review.UserID,
				Title:         "评价回复已更新",
				Message:       "陪玩师更新了对您评价的回复",
				Channel:       "web",
				Priority:      model.NotificationPriorityNormal,
				ReferenceType: "review_reply",
				ReferenceID:   &replyID,
			}
			_ = r.Notifications.Create(ctx, notification)
		}

		result = map[string]interface{}{
			"replyId": reply.ID,
			"status":  string(reply.Status),
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteReviewReply 删除评价回复
func (s *AdminService) DeleteReviewReply(ctx context.Context, userID, replyID uint64) error {
	if s.tx == nil {
		return apierr.InternalError("事务管理器未配置")
	}

	return s.tx.WithTx(ctx, func(r *common.Repos) error {
		// 获取回复
		reply, err := r.ReviewReplies.Get(ctx, replyID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}
			return err
		}

		// 权限检查：只能删除自己的回复
		if reply.AuthorID != userID {
			return ErrUnauthorized
		}

		reviewID := reply.ReviewID

		// 删除回复
		if err := r.ReviewReplies.Delete(ctx, replyID); err != nil {
			return err
		}

		// 记录操作日志
		if r.OpLogs != nil {
			metadata := fmt.Sprintf(`{"reply_id":%d,"content":"%s"}`, replyID, reply.Content)
			log := &model.OperationLog{
				EntityType:   string(model.OpEntityReview),
				EntityID:     reviewID,
				ActorUserID:  &userID,
				Action:       string(model.OpActionDeleteReply),
				Reason:       "删除回复",
				MetadataJSON: []byte(metadata),
			}
			_ = r.OpLogs.Append(ctx, log)
		}

		// 发送通知给评价者
		review, err := r.Reviews.Get(ctx, reviewID)
		if err == nil && review.UserID != userID {
			notification := &model.NotificationEvent{
				UserID:        review.UserID,
				Title:         "评价回复已删除",
				Message:       "陪玩师删除了对您评价的回复",
				Channel:       "web",
				Priority:      model.NotificationPriorityNormal,
				ReferenceType: "review_reply",
				ReferenceID:   &replyID,
			}
			_ = r.Notifications.Create(ctx, notification)
		}

		return nil
	})
}

