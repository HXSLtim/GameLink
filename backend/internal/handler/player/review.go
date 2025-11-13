package player

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service"
	reviewservice "gamelink/internal/service/review"
)

// RegisterReviewRoutes 注册陪玩师评价回复路由。
func RegisterReviewRoutes(router gin.IRouter, svc *reviewservice.ReviewService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/reviews")
	group.Use(authMiddleware)
	group.POST(":id/reply", func(c *gin.Context) { replyReviewHandler(c, svc) })
}

func replyReviewHandler(c *gin.Context, svc *reviewservice.ReviewService) {
	userID := getUserIDFromContext(c)
	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "reviewId 无效")
		return
	}
	var req reviewservice.ReplyReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := svc.ReplyReview(c.Request.Context(), userID, reviewID, req)
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			respondError(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, reviewservice.ErrUnauthorized) {
			respondError(c, http.StatusForbidden, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(c, http.StatusOK, model.APIResponse[*reviewservice.ReplyReviewResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "回复已提交",
		Data:    resp,
	})
}
