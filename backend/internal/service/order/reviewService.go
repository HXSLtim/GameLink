package order

import (
	"context"
	"strings"

	"gamelink/pkg/apierr"
	"gamelink/internal/model"
	"gamelink/pkg/safety"
	"gamelink/internal/repository"
	repoiface "gamelink/internal/repository/interfaces"
	feedservice "gamelink/internal/service/feed"
)

var (
	// ErrReviewNotFound 评价不存在
	ErrReviewNotFound = repository.ErrNotFound
	// ErrReviewValidation 表示输入校验失败
	ErrReviewValidation = apierr.BadRequest("验证失败")
	// ErrAlreadyReviewed 已评价
	ErrAlreadyReviewed = apierr.BadRequest("订单已评价")
	// ErrOrderNotCompleted 订单未完成
	ErrOrderNotCompleted = apierr.BadRequest("订单未完成")
	// ErrReviewUnauthorized 无权操作
	ErrReviewUnauthorized = apierr.Unauthorized("无权操作")
)

// ReviewService 评价服务
//
// 功能：
// 1. 创建评价
// 2. 查询评价列表
// 3. 更新陪玩师评分
// 4. 评价回复管理
type ReviewService struct {
	reviews       repository.ReviewRepository
	orders        repoiface.OrderReader
	players       repository.PlayerRepository
	users         repository.UserRepository
	replies       repository.ReviewReplyRepository
	notifications repository.NotificationRepository
}

// NewReviewService 创建评价服务
func NewReviewService(
	reviews repository.ReviewRepository,
	orders repoiface.OrderReader,
	players repository.PlayerRepository,
	users repository.UserRepository,
	replies repository.ReviewReplyRepository,
	notifications repository.NotificationRepository,
) *ReviewService {
	return &ReviewService{
		reviews:       reviews,
		orders:        orders,
		players:       players,
		users:         users,
		replies:       replies,
		notifications: notifications,
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
		return nil, ErrReviewUnauthorized
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
		return nil, apierr.BadRequest("验证失败: " + err.Error())
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
		return nil, ErrReviewUnauthorized
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

	status := "pending"
	note := result.Reason
	switch result.Decision {
	case feedservice.ModerationDecisionApprove:
		status = "approved"
	case feedservice.ModerationDecisionReject:
		status = "rejected"
	case feedservice.ModerationDecisionManual:
		status = "pending"
	}

	reply.Status = status
	if status != "pending" || note != "" {
		if err := s.replies.UpdateStatus(ctx, reply.ID, status, note); err != nil {
			return nil, err
		}
		reply.ModerationNote = note
	}

	return &ReplyReviewResponse{ReplyID: reply.ID, Status: reply.Status}, nil
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
		return nil, apierr.BadRequest("验证失败: " + err.Error())
	}

	// 获取回复
	reply, err := s.replies.Get(ctx, replyID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, apierr.NotFound("回复不存在")
		}
		return nil, err
	}

	// 权限检查：只能更新自己的回复
	if reply.AuthorID != userID {
		return nil, ErrReviewUnauthorized
	}

	// 更新回复内容
	reply.Content = strings.TrimSpace(req.Content)

	// 重新进行内容审核
	engine := feedservice.NewDefaultModerationEngine()
	result, err := engine.Evaluate(ctx, feedservice.ModerationInput{Content: reply.Content})
	if err != nil {
		return nil, err
	}

	status := "pending"
	note := result.Reason
	switch result.Decision {
	case feedservice.ModerationDecisionApprove:
		status = "approved"
	case feedservice.ModerationDecisionReject:
		status = "rejected"
	case feedservice.ModerationDecisionManual:
		status = "pending"
	}

	reply.Status = status
	reply.ModerationNote = note

	// 更新回复
	if err := s.replies.Update(ctx, reply); err != nil {
		return nil, err
	}

	// 发送通知给评价者
	review, err := s.reviews.Get(ctx, reply.ReviewID)
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
		_ = s.notifications.Create(ctx, notification)
	}

	return &UpdateReplyResponse{ReplyID: reply.ID, Status: reply.Status}, nil
}

// DeleteReply 删除评价回复
func (s *ReviewService) DeleteReply(ctx context.Context, userID, replyID uint64) error {
	// 获取回复
	reply, err := s.replies.Get(ctx, replyID)
	if err != nil {
		if err == repository.ErrNotFound {
			return apierr.NotFound("回复不存在")
		}
		return err
	}

	// 权限检查：只能删除自己的回复
	if reply.AuthorID != userID {
		return ErrReviewUnauthorized
	}

	// 删除回复
	if err := s.replies.Delete(ctx, replyID); err != nil {
		return err
	}

	// 发送通知给评价者
	review, err := s.reviews.Get(ctx, reply.ReviewID)
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
		_ = s.notifications.Create(ctx, notification)
	}

	return nil
}
