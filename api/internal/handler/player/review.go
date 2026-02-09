package player

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	orderservice "gamelink/internal/service/order"
	"gamelink/pkg/apierr"
)

// RegisterReviewRoutes 注册陪玩师评价回复路由。
func RegisterReviewRoutes(router gin.IRouter, svc *orderservice.ReviewService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/reviews")
	group.Use(authMiddleware)
	group.POST(":id/reply", func(c *gin.Context) { replyReviewHandler(c, svc) })
}

// replyReviewHandler 回复评价
// @Summary      回复评价
// @Description  陪玩师回复用户对自己的评价
// @Tags         Player - Review
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                             true  "Bearer {token}"
// @Param        id             path      int                                true  "Review ID"
// @Param        request        body      orderservice.ReplyReviewRequest   true  "Reply content"
// @Success      200            {object}  orderservice.ReplyReviewResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      403            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /player/reviews/{id}/reply [post]
func replyReviewHandler(c *gin.Context, svc *orderservice.ReviewService) {
	userID := getUserIDFromContext(c)
	reviewID, err := parseUintParam(c, "id")
	if err != nil {
		respondError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	var req orderservice.ReplyReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := svc.ReplyReview(c.Request.Context(), userID, reviewID, req)
	if err != nil {
		if err != nil {
			respondError(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, orderservice.ErrUnauthorized) {
			respondError(c, http.StatusForbidden, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON[*orderservice.ReplyReviewResponse](c, http.StatusOK, model.APIResponse[*orderservice.ReplyReviewResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "回复已提交",
		Data:    resp,
	})
}
