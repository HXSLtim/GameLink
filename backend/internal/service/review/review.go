package review

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gamelink/internal/model"
	"gamelink/internal/pkg/safety"
	"gamelink/internal/repository"
	feedservice "gamelink/internal/service/feed"
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
type ReviewService struct {
	reviews repository.ReviewRepository
	orders  repository.OrderRepository
	players repository.PlayerRepository
	users   repository.UserRepository
	replies repository.ReviewReplyRepository
}

// NewReviewService 创建评价服务
func NewReviewService(
	reviews repository.ReviewRepository,
	orders repository.OrderRepository,
	players repository.PlayerRepository,
	users repository.UserRepository,
	replies repository.ReviewReplyRepository,
) *ReviewService {
	return &ReviewService{
		reviews: reviews,
		orders:  orders,
		players: players,
		users:   users,
		replies: replies,
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
