package user

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	gameroomservice "gamelink/internal/service/gameroom"
	"gamelink/pkg/apierr"
)

// GameRoomHandler 游戏房间处理器
type GameRoomHandler struct {
	svc *gameroomservice.Service
}

// NewGameRoomHandler 创建游戏房间处理器
func NewGameRoomHandler(svc *gameroomservice.Service) *GameRoomHandler {
	return &GameRoomHandler{svc: svc}
}

// CreateRoom 创建房间
// @Summary 创建游戏房间
// @Tags 用户-游戏房间
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body gameroomservice.CreateRoomRequest true "创建房间请求"
// @Success 200 {object} model.ChatGroup
// @Router /user/rooms [post]
func (h *GameRoomHandler) CreateRoom(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	var req gameroomservice.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("参数错误").WithDetails(err.Error()))
		return
	}

	room, err := h.svc.CreateRoom(c.Request.Context(), userID, &req)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.Created(c, room)
}

// GetRoom 获取房间详情
// @Summary 获取房间详情
// @Tags 用户-游戏房间
// @Produce json
// @Param id path int true "房间ID"
// @Success 200 {object} model.ChatGroup
// @Router /user/rooms/{id} [get]
func (h *GameRoomHandler) GetRoom(c *gin.Context) {
	roomID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	room, err := h.svc.GetRoom(c.Request.Context(), roomID)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, room)
}

// ListRooms 列出房间
// @Summary 列出游戏房间
// @Tags 用户-游戏房间
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param gameId query int false "游戏ID"
// @Param groupType query string false "房间类型"
// @Param isPrivate query bool false "是否私密"
// @Success 200 {object} resp.PagedResponse
// @Router /user/rooms [get]
func (h *GameRoomHandler) ListRooms(c *gin.Context) {
	page := resp.ParseQueryInt(c, "page", 1)
	pageSize := resp.ParseQueryInt(c, "pageSize", 20)
	gameID := resp.ParseQueryUint64Ptr(c, "gameId")

	opts := repository.GameRoomListOptions{
		Page:     page,
		PageSize: pageSize,
		GameID:   gameID,
	}

	// 解析房间类型
	if groupType := c.Query("groupType"); groupType != "" {
		gt := model.ChatGroupType(groupType)
		opts.GroupType = &gt
	}

	// 解析是否私密
	if isPrivateStr := c.Query("isPrivate"); isPrivateStr != "" {
		isPrivate := isPrivateStr == "true"
		opts.IsPrivate = &isPrivate
	}

	rooms, total, err := h.svc.ListRooms(c.Request.Context(), opts)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.List(c, rooms, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// ListPublicRooms 列出公开房间
// @Summary 列出公开房间
// @Tags 用户-游戏房间
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param gameId query int false "游戏ID"
// @Success 200 {object} resp.PagedResponse
// @Router /user/rooms/public [get]
func (h *GameRoomHandler) ListPublicRooms(c *gin.Context) {
	page := resp.ParseQueryInt(c, "page", 1)
	pageSize := resp.ParseQueryInt(c, "pageSize", 20)
	gameID := resp.ParseQueryUint64Ptr(c, "gameId")

	rooms, total, err := h.svc.ListPublicRooms(c.Request.Context(), gameID, page, pageSize)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.List(c, rooms, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// UpdateRoom 更新房间
// @Summary 更新房间信息
// @Tags 用户-游戏房间
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "房间ID"
// @Param body body gameroomservice.UpdateRoomRequest true "更新请求"
// @Success 200 {object} model.ChatGroup
// @Router /user/rooms/{id} [put]
func (h *GameRoomHandler) UpdateRoom(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	roomID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	var req gameroomservice.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("参数错误").WithDetails(err.Error()))
		return
	}

	room, err := h.svc.UpdateRoom(c.Request.Context(), roomID, userID, &req)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, room)
}

// CloseRoom 关闭房间
// @Summary 关闭房间
// @Tags 用户-游戏房间
// @Produce json
// @Security BearerAuth
// @Param id path int true "房间ID"
// @Success 200 {object} model.SuccessResponse
// @Router /user/rooms/{id} [delete]
func (h *GameRoomHandler) CloseRoom(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	roomID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	if err := h.svc.CloseRoom(c.Request.Context(), roomID, userID); err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, gin.H{"message": "房间已关闭"})
}

// JoinRoom 加入房间
// @Summary 加入房间
// @Tags 用户-游戏房间
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "房间ID"
// @Param body body JoinRoomRequest false "加入请求"
// @Success 200 {object} model.SuccessResponse
// @Router /user/rooms/{id}/join [post]
func (h *GameRoomHandler) JoinRoom(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	roomID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	var req JoinRoomRequest
	_ = c.ShouldBindJSON(&req) // 密码可选

	if err := h.svc.JoinRoom(c.Request.Context(), roomID, userID, req.Password); err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, gin.H{"message": "已加入房间"})
}

// LeaveRoom 离开房间
// @Summary 离开房间
// @Tags 用户-游戏房间
// @Produce json
// @Security BearerAuth
// @Param id path int true "房间ID"
// @Success 200 {object} model.SuccessResponse
// @Router /user/rooms/{id}/leave [post]
func (h *GameRoomHandler) LeaveRoom(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	roomID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	if err := h.svc.LeaveRoom(c.Request.Context(), roomID, userID); err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, gin.H{"message": "已离开房间"})
}

// StartGame 开始游戏
// @Summary 开始游戏
// @Tags 用户-游戏房间
// @Produce json
// @Security BearerAuth
// @Param id path int true "房间ID"
// @Success 200 {object} model.SuccessResponse
// @Router /user/rooms/{id}/start [post]
func (h *GameRoomHandler) StartGame(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	roomID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	if err := h.svc.StartGame(c.Request.Context(), roomID, userID); err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, gin.H{"message": "游戏已开始"})
}

// FinishGame 结束游戏
// @Summary 结束游戏
// @Tags 用户-游戏房间
// @Produce json
// @Security BearerAuth
// @Param id path int true "房间ID"
// @Success 200 {object} model.SuccessResponse
// @Router /user/rooms/{id}/finish [post]
func (h *GameRoomHandler) FinishGame(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	roomID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	if err := h.svc.FinishGame(c.Request.Context(), roomID, userID); err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, gin.H{"message": "游戏已结束"})
}

// KickMember 踢出成员
// @Summary 踢出成员
// @Tags 用户-游戏房间
// @Produce json
// @Security BearerAuth
// @Param id path int true "房间ID"
// @Param userId path int true "用户ID"
// @Success 200 {object} model.SuccessResponse
// @Router /user/rooms/{id}/kick/{userId} [post]
func (h *GameRoomHandler) KickMember(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	roomID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	targetUserID, ok := resp.ParseIDOrFail(c, "userId")
	if !ok {
		return
	}

	if err := h.svc.KickMember(c.Request.Context(), roomID, userID, targetUserID); err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, gin.H{"message": "成员已被踢出"})
}

// GetRoomMembers 获取房间成员
// @Summary 获取房间成员列表
// @Tags 用户-游戏房间
// @Produce json
// @Param id path int true "房间ID"
// @Success 200 {array} model.ChatGroupMember
// @Router /user/rooms/{id}/members [get]
func (h *GameRoomHandler) GetRoomMembers(c *gin.Context) {
	roomID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	members, err := h.svc.GetRoomMembers(c.Request.Context(), roomID)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, members)
}

// GetMyRooms 获取我的房间
// @Summary 获取我的房间列表
// @Tags 用户-游戏房间
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} resp.PagedResponse
// @Router /user/rooms/my [get]
func (h *GameRoomHandler) GetMyRooms(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	page := resp.ParseQueryInt(c, "page", 1)
	pageSize := resp.ParseQueryInt(c, "pageSize", 20)

	rooms, total, err := h.svc.GetUserRooms(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.List(c, rooms, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// ToggleReady 切换准备状态
// @Summary 切换准备状态
// @Tags 用户-游戏房间
// @Produce json
// @Security BearerAuth
// @Param id path int true "房间ID"
// @Success 200 {object} map[string]interface{}
// @Router /user/rooms/{id}/ready [post]
func (h *GameRoomHandler) ToggleReady(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	roomID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	isReady, err := h.svc.ToggleReady(c.Request.Context(), roomID, userID)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, gin.H{"isReady": isReady})
}

// ============================================================================
// Request DTOs
// ============================================================================

// JoinRoomRequest 加入房间请求
type JoinRoomRequest struct {
	Password string `json:"password"`
}

// ============================================================================
// Route Registration
// ============================================================================

// RegisterGameRoomRoutes 注册游戏房间路由
func RegisterGameRoomRoutes(rg *gin.RouterGroup, svc *gameroomservice.Service, authMiddleware gin.HandlerFunc) {
	h := NewGameRoomHandler(svc)

	roomsGroup := rg.Group("/rooms")
	{
		// 公开路由
		roomsGroup.GET("", h.ListRooms)
		roomsGroup.GET("/public", h.ListPublicRooms)
		roomsGroup.GET("/:id", h.GetRoom)
		roomsGroup.GET("/:id/members", h.GetRoomMembers)

		// 需要认证的路由
		authGroup := roomsGroup.Group("")
		authGroup.Use(authMiddleware)
		{
			authGroup.POST("", h.CreateRoom)
			authGroup.PUT("/:id", h.UpdateRoom)
			authGroup.DELETE("/:id", h.CloseRoom)
			authGroup.POST("/:id/join", h.JoinRoom)
			authGroup.POST("/:id/leave", h.LeaveRoom)
			authGroup.POST("/:id/start", h.StartGame)
			authGroup.POST("/:id/finish", h.FinishGame)
			authGroup.POST("/:id/ready", h.ToggleReady)
			authGroup.POST("/:id/kick/:userId", h.KickMember)
			authGroup.GET("/my", h.GetMyRooms)
		}
	}
}
