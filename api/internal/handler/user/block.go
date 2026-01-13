package user

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/model"
	"gamelink/internal/service/userblock"
	"gamelink/pkg/apierr"
)

// BlockHandler 用户拉黑接口（用户端）
type BlockHandler struct {
	svc *userblock.UserBlockService
}

// NewBlockHandler 创建Handler
func NewBlockHandler(svc *userblock.UserBlockService) *BlockHandler {
	return &BlockHandler{svc: svc}
}

// BlockRequest 拉黑请求
type BlockRequest struct {
	BlockedID   uint64              `json:"blockedId" binding:"required"`
	BlockedType model.BlockUserType `json:"blockedType"` // 可选，默认 player
	Reason      string              `json:"reason"`
}

// Block
// @Summary      拉黑用户
// @Tags         User/Block
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BlockRequest  true  "拉黑信息"
// @Success      201  {object}  model.UserBlock
// @Failure      400  {object}  model.ErrorResponse
// @Router       /user/blocks [post]
func (h *BlockHandler) Block(c *gin.Context) {
	var req BlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("invalid request").WithDetails(err.Error()))
		return
	}

	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	// 默认拉黑陪玩师
	blockedType := req.BlockedType
	if blockedType == "" {
		blockedType = model.BlockUserTypePlayer
	}

	block, err := h.svc.Block(c.Request.Context(), userblock.BlockInput{
		BlockerID:   userID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   req.BlockedID,
		BlockedType: blockedType,
		Reason:      req.Reason,
	})
	if err != nil {
		resp.Error(c, err)
		return
	}
	resp.Created(c, block)
}

// Unblock
// @Summary      取消拉黑
// @Tags         User/Block
// @Security     BearerAuth
// @Param        blockedId   path  int  true  "被拉黑用户ID"
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /user/blocks/{blockedId} [delete]
func (h *BlockHandler) Unblock(c *gin.Context) {
	blockedID, ok := resp.ParseIDOrFail(c, "blockedId")
	if !ok {
		return
	}

	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	if err := h.svc.Unblock(c.Request.Context(), userID, blockedID); err != nil {
		resp.Error(c, err)
		return
	}
	resp.OK(c, gin.H{"message": "unblocked"})
}

// ListMyBlocks
// @Summary      获取我的拉黑列表
// @Tags         User/Block
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   model.UserBlock
// @Router       /user/blocks [get]
func (h *BlockHandler) ListMyBlocks(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	blocks, err := h.svc.ListByBlocker(c.Request.Context(), userID, true)
	if err != nil {
		resp.Error(c, err)
		return
	}
	resp.OK(c, blocks)
}

// CheckBlocked
// @Summary      检查是否拉黑了某用户
// @Tags         User/Block
// @Security     BearerAuth
// @Param        targetId   query  int  true  "目标用户ID"
// @Produce      json
// @Success      200  {object}  map[string]bool
// @Router       /user/blocks/check [get]
func (h *BlockHandler) CheckBlocked(c *gin.Context) {
	targetID, ok := resp.ParseQueryIDOrFail(c, "targetId")
	if !ok {
		return
	}

	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	// 检查我是否拉黑了对方
	iBlockedTarget, err := h.svc.IsBlockedBy(c.Request.Context(), userID, targetID)
	if err != nil {
		resp.Error(c, err)
		return
	}

	// 检查对方是否拉黑了我
	targetBlockedMe, err := h.svc.IsBlockedBy(c.Request.Context(), targetID, userID)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, gin.H{
		"iBlockedTarget":  iBlockedTarget,
		"targetBlockedMe": targetBlockedMe,
		"isBlocked":       iBlockedTarget || targetBlockedMe,
	})
}

// RegisterBlockRoutes 注册用户拉黑路由（用户端）
func RegisterBlockRoutes(router gin.IRouter, svc *userblock.UserBlockService, authMiddleware gin.HandlerFunc) {
	h := NewBlockHandler(svc)

	group := router.Group("/blocks")
	group.Use(authMiddleware)
	{
		group.GET("", h.ListMyBlocks)
		group.GET("/check", h.CheckBlocked)
		group.POST("", h.Block)
		group.DELETE("/:blockedId", h.Unblock)
	}
}
