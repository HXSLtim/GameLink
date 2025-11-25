package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/apierr"
	"gamelink/internal/model"
	"gamelink/internal/service"
	feedservice "gamelink/internal/service/feed"
)

// RegisterFeedRoutes 注册社区动态路由。
func RegisterFeedRoutes(router gin.IRouter, svc *feedservice.Service, authMiddleware gin.HandlerFunc) {
	group := router.Group("/feeds")
	group.Use(authMiddleware)
	group.POST("", func(c *gin.Context) { createFeedHandler(c, svc) })
	group.GET("", func(c *gin.Context) { listFeedsHandler(c, svc) })
	group.POST(":id/report", func(c *gin.Context) { reportFeedHandler(c, svc) })
}

// createFeedHandler 发布动态
// @Summary      发布动态
// @Description  用户发布社区动态，包括文字、图片等内容
// @Tags         User - Feed
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                          true  "Bearer {token}"
// @Param        request        body      feedservice.CreateFeedRequest   true  "Feed content"
// @Success      200            {object}  feedservice.FeedView
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /user/feeds [post]
func createFeedHandler(c *gin.Context, svc *feedservice.Service) {
	userID := getUserIDFromContext(c)
	var req feedservice.CreateFeedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	feed, err := svc.CreateFeed(c.Request.Context(), userID, req)
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			respondError(c, http.StatusBadRequest, err.Error())
		} else {
			respondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    feed,
	})
}

// listFeedsHandler 获取动态列表
// @Summary      获取动态列表
// @Description  获取社区动态列表，支持游标分页
// @Tags         User - Feed
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        limit          query     int     false  "Limit (default 20)"
// @Param        cursor         query     string  false  "Cursor for pagination"
// @Success      200            {object}  feedservice.ListFeedsResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /user/feeds [get]
func listFeedsHandler(c *gin.Context, svc *feedservice.Service) {
	userID := getUserIDFromContext(c)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cursor := c.Query("cursor")
	resp, err := svc.ListFeeds(c.Request.Context(), userID, feedservice.ListFeedsRequest{
		Cursor: cursor,
		Limit:  limit,
	})
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			respondError(c, http.StatusBadRequest, err.Error())
		} else {
			respondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    resp,
	})
}

// reportFeedHandler 举报动态
// @Summary      举报动态
// @Description  举报不当社区动态内容
// @Tags         User - Feed
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string              true  "Bearer {token}"
// @Param        id             path      int                 true  "Feed ID"
// @Param        request        body      object{reason=string}  true  "Report reason"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /user/feeds/{id}/report [post]
func reportFeedHandler(c *gin.Context, svc *feedservice.Service) {
	userID := getUserIDFromContext(c)
	feedID, err := parseUintParam(c, "id")
	if err != nil {
		respondError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	var body struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := svc.ReportFeed(c.Request.Context(), userID, feedID, body.Reason); err != nil {
		if errors.Is(err, service.ErrValidation) {
			respondError(c, http.StatusBadRequest, err.Error())
		} else {
			respondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "举报成功",
	})
}
