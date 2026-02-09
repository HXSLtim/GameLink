package public

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
)

// PublicReviewHandler handles public review endpoints.
type PublicReviewHandler struct {
	reviews repository.ReviewRepository
	users   repository.UserRepository
}

// NewPublicReviewHandler creates a public review handler.
func NewPublicReviewHandler(reviews repository.ReviewRepository, users repository.UserRepository) *PublicReviewHandler {
	return &PublicReviewHandler{reviews: reviews, users: users}
}

// RegisterRoutes registers public review routes.
func (h *PublicReviewHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/players/:id/reviews", h.ListPlayerReviews)
}

type publicReviewDTO struct {
	ID            uint64 `json:"id"`
	Rating        int    `json:"rating"`
	Comment       string `json:"comment"`
	UserNickname  string `json:"userNickname"`
	UserAvatarURL string `json:"userAvatarUrl"`
	CreatedAt     string `json:"createdAt"`
}

// ListPlayerReviews 获取陪玩师评价列表（公开）
// @Summary 获取陪玩师评价列表
// @Description 获取陪玩师的公开评价列表
// @Tags 公共-评价
// @Accept json
// @Produce json
// @Param id path int true "陪玩师ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param rating query int false "评分筛选"
// @Success 200 {object} resp.PagedResponse
// @Failure 500 {object} apierr.APIError
// @Router /public/players/{id}/reviews [get]
func (h *PublicReviewHandler) ListPlayerReviews(c *gin.Context) {
	playerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resp.Error(c, apierr.BadRequest("无效的陪玩师ID"))
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	var rating *int
	if ratingStr := c.Query("rating"); ratingStr != "" {
		if v, parseErr := strconv.Atoi(ratingStr); parseErr == nil {
			rating = &v
		}
	}

	status := model.ReviewStatusApproved
	isPublic := true
	playerIDPtr := &playerID
	reviews, total, err := h.reviews.List(c.Request.Context(), repository.ReviewListOptions{
		PlayerID: playerIDPtr,
		Page:     page,
		PageSize: pageSize,
		Status:   &status,
		IsPublic: &isPublic,
		Rating:   rating,
	})
	if err != nil {
		resp.Error(c, apierr.InternalError("获取评价列表失败"))
		return
	}

	items := make([]publicReviewDTO, 0, len(reviews))
	for _, r := range reviews {
		nickname := "匿名用户"
		avatar := ""
		if !r.IsAnonymous {
			if user, err := h.users.Get(c.Request.Context(), r.UserID); err == nil {
				if user.Nickname != "" {
					nickname = user.Nickname
				} else if user.Name != "" {
					nickname = user.Name
				}
				avatar = user.AvatarURL
			}
		}
		items = append(items, publicReviewDTO{
			ID:            r.ID,
			Rating:        int(r.Score),
			Comment:       r.Content,
			UserNickname:  nickname,
			UserAvatarURL: avatar,
			CreatedAt:     r.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	resp.List(c, items, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}
