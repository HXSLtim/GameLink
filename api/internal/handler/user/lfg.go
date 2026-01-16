package user

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	lfgservice "gamelink/internal/service/lfg"
	"gamelink/pkg/apierr"
)

// LFGHandler 快速匹配处理器
type LFGHandler struct {
	svc *lfgservice.Service
}

// NewLFGHandler 创建快速匹配处理器
func NewLFGHandler(svc *lfgservice.Service) *LFGHandler {
	return &LFGHandler{svc: svc}
}

// CreateRequest 创建匹配请求
// @Summary 创建匹配请求
// @Tags 用户-快速匹配
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body lfgservice.CreateLFGRequest true "创建请求"
// @Success 200 {object} model.LFGRequest
// @Router /user/lfg [post]
func (h *LFGHandler) CreateRequest(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	var req lfgservice.CreateLFGRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("参数错误").WithDetails(err.Error()))
		return
	}

	// 验证请求类型
	if req.RequestType != model.LFGFindPlayer && req.RequestType != model.LFGFindTeam {
		resp.Error(c, apierr.BadRequest("无效的请求类型"))
		return
	}

	request, err := h.svc.CreateRequest(c.Request.Context(), userID, &req)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.Created(c, request)
}

// GetRequest 获取匹配请求详情
// @Summary 获取匹配请求详情
// @Tags 用户-快速匹配
// @Produce json
// @Param id path int true "请求ID"
// @Success 200 {object} model.LFGRequest
// @Router /user/lfg/{id} [get]
func (h *LFGHandler) GetRequest(c *gin.Context) {
	requestID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	request, err := h.svc.GetRequest(c.Request.Context(), requestID)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, request)
}

// GetMyActiveRequest 获取我的活跃请求
// @Summary 获取我的活跃匹配请求
// @Tags 用户-快速匹配
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.LFGRequest
// @Router /user/lfg/active [get]
func (h *LFGHandler) GetMyActiveRequest(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	request, err := h.svc.GetActiveRequest(c.Request.Context(), userID)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, request)
}

// CancelRequest 取消匹配请求
// @Summary 取消匹配请求
// @Tags 用户-快速匹配
// @Produce json
// @Security BearerAuth
// @Param id path int true "请求ID"
// @Success 200 {object} model.SuccessResponse
// @Router /user/lfg/{id} [delete]
func (h *LFGHandler) CancelRequest(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	requestID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	if err := h.svc.CancelRequest(c.Request.Context(), requestID, userID); err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, gin.H{"message": "请求已取消"})
}

// ListRequests 列出匹配请求
// @Summary 列出匹配请求
// @Tags 用户-快速匹配
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param gameId query int false "游戏ID"
// @Param requestType query string false "请求类型"
// @Success 200 {object} resp.PagedResponse
// @Router /user/lfg [get]
func (h *LFGHandler) ListRequests(c *gin.Context) {
	page := resp.ParseQueryInt(c, "page", 1)
	pageSize := resp.ParseQueryInt(c, "pageSize", 20)
	gameID := resp.ParseQueryUint64Ptr(c, "gameId")

	opts := repository.LFGRequestListOptions{
		Page:     page,
		PageSize: pageSize,
		GameID:   gameID,
	}

	// 解析请求类型
	if requestType := c.Query("requestType"); requestType != "" {
		rt := model.LFGRequestType(requestType)
		opts.RequestType = &rt
	}

	// 默认只显示等待中的请求
	pendingStatus := model.LFGPending
	opts.Status = &pendingStatus

	requests, total, err := h.svc.ListRequests(c.Request.Context(), opts)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.List(c, requests, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// ListPendingRequests 列出等待中的请求
// @Summary 列出等待中的匹配请求
// @Tags 用户-快速匹配
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param gameId query int false "游戏ID"
// @Success 200 {object} resp.PagedResponse
// @Router /user/lfg/pending [get]
func (h *LFGHandler) ListPendingRequests(c *gin.Context) {
	page := resp.ParseQueryInt(c, "page", 1)
	pageSize := resp.ParseQueryInt(c, "pageSize", 20)
	gameID := resp.ParseQueryUint64Ptr(c, "gameId")

	requests, total, err := h.svc.ListPendingRequests(c.Request.Context(), gameID, page, pageSize)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.List(c, requests, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// GetMyRequests 获取我的请求列表
// @Summary 获取我的匹配请求列表
// @Tags 用户-快速匹配
// @Produce json
// @Security BearerAuth
// @Param status query string false "状态筛选"
// @Success 200 {array} model.LFGRequest
// @Router /user/lfg/my [get]
func (h *LFGHandler) GetMyRequests(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	var status *model.LFGRequestStatus
	if statusStr := c.Query("status"); statusStr != "" {
		s := model.LFGRequestStatus(statusStr)
		status = &s
	}

	requests, err := h.svc.ListUserRequests(c.Request.Context(), userID, status)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, requests)
}

// AcceptRequest 接受匹配请求
// @Summary 接受匹配请求（陪玩师接单）
// @Tags 用户-快速匹配
// @Produce json
// @Security BearerAuth
// @Param id path int true "请求ID"
// @Success 200 {object} model.ChatGroup
// @Router /user/lfg/{id}/accept [post]
func (h *LFGHandler) AcceptRequest(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	requestID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	room, err := h.svc.AcceptRequest(c.Request.Context(), requestID, userID)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, room)
}

// FindMatches 查找匹配
// @Summary 查找可匹配的请求
// @Tags 用户-快速匹配
// @Produce json
// @Param id path int true "请求ID"
// @Param limit query int false "返回数量限制"
// @Success 200 {array} model.LFGRequest
// @Router /user/lfg/{id}/matches [get]
func (h *LFGHandler) FindMatches(c *gin.Context) {
	requestID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	limit := resp.ParseQueryInt(c, "limit", 10)

	matches, err := h.svc.FindMatches(c.Request.Context(), requestID, limit)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, matches)
}

// GetPendingCount 获取等待中的请求数量
// @Summary 获取等待中的请求数量
// @Tags 用户-快速匹配
// @Produce json
// @Param gameId query int false "游戏ID"
// @Success 200 {object} model.SuccessResponse
// @Router /user/lfg/count [get]
func (h *LFGHandler) GetPendingCount(c *gin.Context) {
	gameID := resp.ParseQueryUint64Ptr(c, "gameId")

	count, err := h.svc.CountPending(c.Request.Context(), gameID)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, gin.H{"count": count})
}

// ============================================================================
// Route Registration
// ============================================================================

// RegisterLFGRoutes 注册快速匹配路由
func RegisterLFGRoutes(rg *gin.RouterGroup, svc *lfgservice.Service, authMiddleware gin.HandlerFunc) {
	h := NewLFGHandler(svc)

	lfgGroup := rg.Group("/lfg")
	{
		// 公开路由
		lfgGroup.GET("", h.ListRequests)
		lfgGroup.GET("/pending", h.ListPendingRequests)
		lfgGroup.GET("/count", h.GetPendingCount)
		lfgGroup.GET("/:id", h.GetRequest)
		lfgGroup.GET("/:id/matches", h.FindMatches)

		// 需要认证的路由
		authGroup := lfgGroup.Group("")
		authGroup.Use(authMiddleware)
		{
			authGroup.POST("", h.CreateRequest)
			authGroup.DELETE("/:id", h.CancelRequest)
			authGroup.GET("/active", h.GetMyActiveRequest)
			authGroup.GET("/my", h.GetMyRequests)
			authGroup.POST("/:id/accept", h.AcceptRequest)
		}
	}
}
