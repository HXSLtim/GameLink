package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service/review"
)

// RegisterUserReviewRoutes 注册用户端评价路由
func RegisterUserReviewRoutes(router gin.IRouter, svc *review.ReviewService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/user/reviews")
	group.Use(authMiddleware) // 需要认证
	{
		group.POST("", func(c *gin.Context) { createReviewHandler(c, svc) })
		group.GET("/my", func(c *gin.Context) { getMyReviewsHandler(c, svc) })
	}
}

// createReviewHandler 创建评价
// @Summary      创建评价
// @Description  为已完成订单创建评价
// @Tags         User - Reviews
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                        true  "Bearer {token}"
// @Param        request        body      review.CreateReviewRequest    true  "创建评价请求"
// @Success      200            {object}  model.APIResponse[review.CreateReviewResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /user/reviews [post]
func createReviewHandler(c *gin.Context, svc *review.ReviewService) {
	userID := getUserIDFromContext(c)

	var req review.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := svc.CreateReview(c.Request.Context(), userID, req)
	if err != nil {
		if err == review.ErrAlreadyReviewed {
			respondError(c, http.StatusBadRequest, err.Error())
			return
		}
		if err == review.ErrOrderNotCompleted {
			respondError(c, http.StatusBadRequest, err.Error())
			return
		}
		if err == review.ErrUnauthorized {
			respondError(c, http.StatusForbidden, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[review.CreateReviewResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "评价创建成功",
		Data:    *resp,
	})
}

// getMyReviewsHandler 获取我的评价列表
// @Summary      获取我的评价列表
// @Description  获取当前用户的评价列表
// @Tags         User - Reviews
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        page           query     int     false  "页码"
// @Param        pageSize       query     int     false  "每页数量"
// @Success      200            {object}  model.APIResponse[review.MyReviewListResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /user/reviews/my [get]
func getMyReviewsHandler(c *gin.Context, svc *review.ReviewService) {
	userID := getUserIDFromContext(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	resp, err := svc.GetMyReviews(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[review.MyReviewListResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *resp,
	})
}
