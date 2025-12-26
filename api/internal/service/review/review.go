package review

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	repoiface "gamelink/internal/repository/interfaces"
	feedservice "gamelink/internal/service/feed"
	"gamelink/pkg/safety"
)

var (
	// ErrNotFound 评价不存在
	ErrNotFound = repository.ErrNotFound
	// ErrValidation 表示输入校验失败
	ErrValidation = errors.New("validation failed")
	// ErrAlreadyReviewed 已评价
	ErrAlreadyReviewed = errors.New("order already reviewed")
	// ErrOrderNotCompleted 订单未完成
	ErrOrderNotCompleted = errors.New("order not completed")
	// ErrUnauthorized 无权操作
	ErrUnauthorized = errors.New("unauthorized")
)

// ReviewService 评价服务
//
// 功能：
// 1. 创建评价
// 2. 查询评价列表
// 3. 更新陪玩师评分
// 4. 举报管理
// 5. 评价回复管理
type ReviewService struct {
	reviews       repository.ReviewRepository
	orders        repoiface.OrderReader
	players       repository.PlayerRepository
	users         repository.UserRepository
	replies       repository.ReviewReplyRepository
	reports       repository.ReviewReportRepository
	notifications repository.NotificationRepository
	opLogs        repository.OperationLogRepository
}

// NewReviewService 创建评价服务
func NewReviewService(
	reviews repository.ReviewRepository,
	orders repoiface.OrderReader,
	players repository.PlayerRepository,
	users repository.UserRepository,
	replies repository.ReviewReplyRepository,
	reports repository.ReviewReportRepository,
	notifications repository.NotificationRepository,
	opLogs repository.OperationLogRepository,
) *ReviewService {
	return &ReviewService{
		reviews:       reviews,
		orders:        orders,
		players:       players,
		users:         users,
		replies:       replies,
		reports:       reports,
		notifications: notifications,
		opLogs:        opLogs,
	}
}

// CreateReviewRequest 创建评价请求
type CreateReviewRequest struct {
	OrderID   uint64   `json:"orderId" binding:"required"`
	Rating    int      `json:"rating" binding:"required,min=1,max=5"`
	Comment   string   `json:"comment" binding:"max=500"`
	Tags      []string `json:"tags"`      // 评价标签
	Anonymous bool     `json:"anonymous"` // 是否匿名
}

// CreateReviewResponse 创建评价响应
type CreateReviewResponse struct {
	ReviewID uint64 `json:"reviewId"`
}

// MyReviewDTO 我的评价信息
type MyReviewDTO struct {
	ReviewDTO
	OrderTitle     string `json:"orderTitle"`
	PlayerNickname string `json:"playerNickname"`
}

// ReviewDTO 评价信息
type ReviewDTO struct {
	ID            uint64 `json:"id"`
	OrderID       uint64 `json:"orderId"`
	Rating        int    `json:"rating"`
	Comment       string `json:"comment"`
	UserNickname  string `json:"userNickname"`
	UserAvatarURL string `json:"userAvatarUrl"`
	CreatedAt     string `json:"createdAt"`
}

// MyReviewListResponse 我的评价列表响应
type MyReviewListResponse struct {
	Reviews []MyReviewDTO `json:"reviews"`
	Total   int64         `json:"total"`
}

// ReplyReviewRequest 陪玩师回复评价请求
type ReplyReviewRequest struct {
	Content string `json:"content"`
}

// ReplyReviewResponse 陪玩师回复评价响应
type ReplyReviewResponse struct {
	ReplyID uint64 `json:"replyId"`
	Status  string `json:"status"`
}

// CreateReview 创建评价
func (s *ReviewService) CreateReview(ctx context.Context, userID uint64, req CreateReviewRequest) (*CreateReviewResponse, error) {
	// 验证订单
	order, err := s.orders.Get(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}

	// 权限检查：只能评价自己的订单
	if order.UserID != userID {
		return nil, ErrUnauthorized
	}

	// 状态检查：只有已完成的订单可以评价
	if order.Status != model.OrderStatusCompleted {
		return nil, ErrOrderNotCompleted
	}

	// 检查是否已评价
	orderIDPtr := &req.OrderID
	existingReviews, _, err := s.reviews.List(ctx, repository.ReviewListOptions{
		OrderID:  orderIDPtr,
		Page:     1,
		PageSize: 1,
	})
	if err == nil && len(existingReviews) > 0 {
		return nil, ErrAlreadyReviewed
	}

	// 创建评价
	playerID := order.GetPlayerID()
	review := &model.Review{
		OrderID:  req.OrderID,
		UserID:   userID,
		PlayerID: playerID,
		Score:    model.Rating(req.Rating),
		Content:  req.Comment,
	}

	if err := s.reviews.Create(ctx, review); err != nil {
		return nil, err
	}

	// 更新陪玩师评分
	if playerID > 0 {
		if err := s.updatePlayerRating(ctx, playerID); err != nil {
			// 更新评分失败不影响评价创建
		}
	}

	return &CreateReviewResponse{
		ReviewID: review.ID,
	}, nil
}

// GetMyReviews 获取我的评价列表
func (s *ReviewService) GetMyReviews(ctx context.Context, userID uint64, page, pageSize int) (*MyReviewListResponse, error) {
	// 默认分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 查询评价
	userIDPtr := &userID
	reviews, total, err := s.reviews.List(ctx, repository.ReviewListOptions{
		UserID:   userIDPtr,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	myReviews := make([]MyReviewDTO, 0, len(reviews))
	for _, r := range reviews {
		// 获取订单信息
		order, err := s.orders.Get(ctx, r.OrderID)
		if err != nil {
			continue
		}

		// 获取陪玩师信息
		player, err := s.players.Get(ctx, r.PlayerID)
		if err != nil {
			continue
		}

		// 获取用户信息
		user, err := s.users.Get(ctx, r.UserID)
		if err != nil {
			continue
		}

		myReviews = append(myReviews, MyReviewDTO{
			ReviewDTO: ReviewDTO{
				ID:            r.ID,
				OrderID:       r.OrderID,
				Rating:        int(r.Score),
				Comment:       r.Content,
				UserNickname:  user.Name,
				UserAvatarURL: user.AvatarURL,
				CreatedAt:     r.CreatedAt.Format("2006-01-02 15:04:05"),
			},
			OrderTitle:     order.Title,
			PlayerNickname: player.Nickname,
		})
	}

	return &MyReviewListResponse{
		Reviews: myReviews,
		Total:   total,
	}, nil
}

// GetPlayerReviews 获取陪玩师的评价列表
func (s *ReviewService) GetPlayerReviews(ctx context.Context, playerID uint64, page, pageSize int) ([]ReviewDTO, int64, error) {
	// 默认分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 查询评价
	playerIDPtr := &playerID
	reviews, total, err := s.reviews.List(ctx, repository.ReviewListOptions{
		PlayerID: playerIDPtr,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, 0, err
	}

	// 转换为 DTO
	reviewDTOs := make([]ReviewDTO, 0, len(reviews))
	for _, r := range reviews {
		user, err := s.users.Get(ctx, r.UserID)
		if err != nil {
			continue
		}

		reviewDTOs = append(reviewDTOs, ReviewDTO{
			ID:            r.ID,
			OrderID:       r.OrderID,
			Rating:        int(r.Score),
			Comment:       r.Content,
			UserNickname:  user.Name,
			UserAvatarURL: user.AvatarURL,
			CreatedAt:     r.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return reviewDTOs, total, nil
}

// ReplyReview 陪玩师回复评价
func (s *ReviewService) ReplyReview(ctx context.Context, userID, reviewID uint64, req ReplyReviewRequest) (*ReplyReviewResponse, error) {
	if err := safety.ValidateText(req.Content, 500); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}

	player, err := s.players.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	authorID := userID
	playerID := player.ID

	review, err := s.reviews.Get(ctx, reviewID)
	if err != nil {
		return nil, err
	}
	if review.PlayerID != playerID {
		return nil, ErrUnauthorized
	}

	reply := &model.ReviewReply{
		ReviewID: reviewID,
		AuthorID: authorID,
		Content:  strings.TrimSpace(req.Content),
	}
	if err := s.replies.Create(ctx, reply); err != nil {
		return nil, err
	}

	engine := feedservice.NewDefaultModerationEngine()
	result, err := engine.Evaluate(ctx, feedservice.ModerationInput{Content: reply.Content})
	if err != nil {
		return nil, err
	}

	var status model.ReviewReplyStatus
	note := result.Reason
	switch result.Decision {
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
	if status != model.ReviewReplyStatusPending || note != "" {
		if err := s.replies.UpdateStatus(ctx, reply.ID, string(status), note); err != nil {
			return nil, err
		}
		reply.ModerationNote = note
	}

	// 记录操作日志
	if s.opLogs != nil {
		metadata := fmt.Sprintf(`{"reply_id":%d,"status":"%s"}`, reply.ID, status)
		log := &model.OperationLog{
			EntityType:   string(model.OpEntityReview),
			EntityID:     reviewID,
			ActorUserID:  &userID,
			Action:       string(model.OpActionReply),
			Reason:       "回复评价",
			MetadataJSON: []byte(metadata),
		}
		_ = s.opLogs.Append(ctx, log)
	}

	return &ReplyReviewResponse{ReplyID: reply.ID, Status: string(reply.Status)}, nil
}

// updatePlayerRating 更新陪玩师评分
func (s *ReviewService) updatePlayerRating(ctx context.Context, playerID uint64) error {
	// 获取陪玩师
	player, err := s.players.Get(ctx, playerID)
	if err != nil {
		return err
	}

	// 获取所有评价
	playerIDPtr := &playerID
	reviews, _, err := s.reviews.List(ctx, repository.ReviewListOptions{
		PlayerID: playerIDPtr,
		Page:     1,
		PageSize: 10000, // 获取所有评价
	})
	if err != nil {
		return err
	}

	if len(reviews) == 0 {
		return nil
	}

	// 计算平均评分
	var totalScore int
	for _, r := range reviews {
		totalScore += int(r.Score)
	}

	player.RatingAverage = float32(totalScore) / float32(len(reviews))
	player.RatingCount = uint32(len(reviews))

	return s.players.Update(ctx, player)
}

// ApproveReviewRequest 批准评价请求
type ApproveReviewRequest struct {
	ReviewID uint64 `json:"reviewId" binding:"required"`
}

// RejectReviewRequest 拒绝评价请求
type RejectReviewRequest struct {
	ReviewID uint64 `json:"reviewId" binding:"required"`
	Reason   string `json:"reason" binding:"required,max=500"`
}

// BatchApproveRequest 批量批准评价请求
type BatchApproveRequest struct {
	ReviewIDs []uint64 `json:"reviewIds" binding:"required,min=1"`
}

// BatchRejectRequest 批量拒绝评价请求
type BatchRejectRequest struct {
	ReviewIDs []uint64 `json:"reviewIds" binding:"required,min=1"`
	Reason    string   `json:"reason" binding:"required,max=500"`
}

// ListPendingReviews 获取待审核评价列表
func (s *ReviewService) ListPendingReviews(ctx context.Context, page, pageSize int) ([]model.Review, int64, error) {
	// 默认分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.reviews.ListPending(ctx, page, pageSize)
}

// ApproveReview 批准评价
func (s *ReviewService) ApproveReview(ctx context.Context, reviewID uint64, actorUserID *uint64) error {
	// 获取评价
	review, err := s.reviews.Get(ctx, reviewID)
	if err != nil {
		return err
	}

	// 检查状态：只有待审核的评价可以批准
	if review.Status != model.ReviewStatusPending {
		return errors.New("只能批准待审核的评价")
	}

	oldStatus := review.Status

	// 更新状态为已通过
	if err := s.reviews.UpdateStatus(ctx, reviewID, model.ReviewStatusApproved, ""); err != nil {
		return err
	}

	// 记录操作日志
	if s.opLogs != nil {
		metadata := fmt.Sprintf(`{"old_status":"%s","new_status":"%s"}`, oldStatus, model.ReviewStatusApproved)
		log := &model.OperationLog{
			EntityType:   string(model.OpEntityReview),
			EntityID:     reviewID,
			ActorUserID:  actorUserID,
			Action:       string(model.OpActionApprove),
			Reason:       "批准评价",
			MetadataJSON: []byte(metadata),
		}
		_ = s.opLogs.Append(ctx, log)
	}

	return nil
}

// RejectReview 拒绝评价
func (s *ReviewService) RejectReview(ctx context.Context, reviewID uint64, reason string, actorUserID *uint64) error {
	// 验证拒绝原因
	if reason == "" {
		return errors.New("拒绝原因不能为空")
	}

	// 获取评价
	review, err := s.reviews.Get(ctx, reviewID)
	if err != nil {
		return err
	}

	// 检查状态：只有待审核的评价可以拒绝
	if review.Status != model.ReviewStatusPending {
		return errors.New("只能拒绝待审核的评价")
	}

	oldStatus := review.Status

	// 更新状态为已拒绝
	if err := s.reviews.UpdateStatus(ctx, reviewID, model.ReviewStatusRejected, reason); err != nil {
		return err
	}

	// 记录操作日志
	if s.opLogs != nil {
		metadata := fmt.Sprintf(`{"old_status":"%s","new_status":"%s","rejection_reason":"%s"}`, oldStatus, model.ReviewStatusRejected, reason)
		log := &model.OperationLog{
			EntityType:   string(model.OpEntityReview),
			EntityID:     reviewID,
			ActorUserID:  actorUserID,
			Action:       string(model.OpActionReject),
			Reason:       reason,
			MetadataJSON: []byte(metadata),
		}
		_ = s.opLogs.Append(ctx, log)
	}

	return nil
}

// BatchApprove 批量批准评价
func (s *ReviewService) BatchApprove(ctx context.Context, reviewIDs []uint64, actorUserID *uint64) error {
	if len(reviewIDs) == 0 {
		return errors.New("评价ID列表不能为空")
	}

	// 验证所有评价都是待审核状态
	for _, id := range reviewIDs {
		review, err := s.reviews.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("获取评价 %d 失败: %w", id, err)
		}
		if review.Status != model.ReviewStatusPending {
			return fmt.Errorf("评价 %d 不是待审核状态", id)
		}
	}

	// 批量更新状态
	if err := s.reviews.BatchUpdateStatus(ctx, reviewIDs, model.ReviewStatusApproved, ""); err != nil {
		return err
	}

	// 记录操作日志
	if s.opLogs != nil {
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
			_ = s.opLogs.Append(ctx, log)
		}
	}

	return nil
}

// BatchReject 批量拒绝评价
func (s *ReviewService) BatchReject(ctx context.Context, reviewIDs []uint64, reason string, actorUserID *uint64) error {
	if len(reviewIDs) == 0 {
		return errors.New("评价ID列表不能为空")
	}

	if reason == "" {
		return errors.New("拒绝原因不能为空")
	}

	// 验证所有评价都是待审核状态
	for _, id := range reviewIDs {
		review, err := s.reviews.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("获取评价 %d 失败: %w", id, err)
		}
		if review.Status != model.ReviewStatusPending {
			return fmt.Errorf("评价 %d 不是待审核状态", id)
		}
	}

	// 批量更新状态
	if err := s.reviews.BatchUpdateStatus(ctx, reviewIDs, model.ReviewStatusRejected, reason); err != nil {
		return err
	}

	// 记录操作日志
	if s.opLogs != nil {
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
			_ = s.opLogs.Append(ctx, log)
		}
	}

	return nil
}

// UpdateReplyRequest 更新回复请求
type UpdateReplyRequest struct {
	Content string `json:"content" binding:"required,max=500"`
}

// UpdateReplyResponse 更新回复响应
type UpdateReplyResponse struct {
	ReplyID uint64 `json:"replyId"`
	Status  string `json:"status"`
}

// UpdateReply 更新评价回复
func (s *ReviewService) UpdateReply(ctx context.Context, userID, replyID uint64, req UpdateReplyRequest) (*UpdateReplyResponse, error) {
	// 验证内容
	if err := safety.ValidateText(req.Content, 500); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}

	// 获取回复
	reply, err := s.replies.Get(ctx, replyID)
	if err != nil {
		return nil, err
	}

	// 权限检查：只能更新自己的回复
	if reply.AuthorID != userID {
		return nil, ErrUnauthorized
	}

	oldContent := reply.Content

	// 更新回复内容
	reply.Content = strings.TrimSpace(req.Content)

	// 重新进行内容审核
	engine := feedservice.NewDefaultModerationEngine()
	result, err := engine.Evaluate(ctx, feedservice.ModerationInput{Content: reply.Content})
	if err != nil {
		return nil, err
	}

	var status model.ReviewReplyStatus
	note := result.Reason
	switch result.Decision {
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
	if err := s.replies.Update(ctx, reply); err != nil {
		return nil, err
	}

	// 记录操作日志
	if s.opLogs != nil {
		metadata := fmt.Sprintf(`{"reply_id":%d,"old_content":"%s","new_content":"%s","status":"%s"}`, replyID, oldContent, reply.Content, status)
		log := &model.OperationLog{
			EntityType:   string(model.OpEntityReview),
			EntityID:     reply.ReviewID,
			ActorUserID:  &userID,
			Action:       string(model.OpActionUpdateReply),
			Reason:       "更新回复",
			MetadataJSON: []byte(metadata),
		}
		_ = s.opLogs.Append(ctx, log)
	}

	// 发送通知给评价者
	review, err := s.reviews.Get(ctx, reply.ReviewID)
	if err == nil && review.UserID != userID {
		notification := &model.NotificationEvent{
			UserID:        review.UserID,
			Title:         "评价回复已更新",
			Message:       fmt.Sprintf("陪玩师更新了对您评价的回复"),
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "review_reply",
			ReferenceID:   &replyID,
		}
		_ = s.notifications.Create(ctx, notification)
	}

	return &UpdateReplyResponse{ReplyID: reply.ID, Status: string(reply.Status)}, nil
}

// DeleteReply 删除评价回复
func (s *ReviewService) DeleteReply(ctx context.Context, userID, replyID uint64) error {
	// 获取回复
	reply, err := s.replies.Get(ctx, replyID)
	if err != nil {
		return err
	}

	// 权限检查：只能删除自己的回复
	if reply.AuthorID != userID {
		return ErrUnauthorized
	}

	reviewID := reply.ReviewID

	// 删除回复
	if err := s.replies.Delete(ctx, replyID); err != nil {
		return err
	}

	// 记录操作日志
	if s.opLogs != nil {
		metadata := fmt.Sprintf(`{"reply_id":%d,"content":"%s"}`, replyID, reply.Content)
		log := &model.OperationLog{
			EntityType:   string(model.OpEntityReview),
			EntityID:     reviewID,
			ActorUserID:  &userID,
			Action:       string(model.OpActionDeleteReply),
			Reason:       "删除回复",
			MetadataJSON: []byte(metadata),
		}
		_ = s.opLogs.Append(ctx, log)
	}

	// 发送通知给评价者
	review, err := s.reviews.Get(ctx, reviewID)
	if err == nil && review.UserID != userID {
		notification := &model.NotificationEvent{
			UserID:        review.UserID,
			Title:         "评价回复已删除",
			Message:       fmt.Sprintf("陪玩师删除了对您评价的回复"),
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "review_reply",
			ReferenceID:   &replyID,
		}
		_ = s.notifications.Create(ctx, notification)
	}

	return nil
}

// GetUsersByIDs 批量获取用户信息
func (s *ReviewService) GetUsersByIDs(ctx context.Context, ids []uint64) ([]model.User, error) {
	if len(ids) == 0 {
		return []model.User{}, nil
	}
	return s.users.GetByIDs(ctx, ids)
}
