package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

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

func reportFeedHandler(c *gin.Context, svc *feedservice.Service) {
	userID := getUserIDFromContext(c)
	idParam := c.Param("id")
	feedID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "feedId 无效")
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
