package admin

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service/userblock"
	"gamelink/pkg/apierr"
)

// UserBlockHandler 用户拉黑管理接口
type UserBlockHandler struct {
	svc *userblock.UserBlockService
}

// NewUserBlockHandler 创建Handler
func NewUserBlockHandler(svc *userblock.UserBlockService) *UserBlockHandler {
	return &UserBlockHandler{svc: svc}
}

// ============================================================================
// 拉黑记录管理
// ============================================================================

// ListUserBlocks
// @Summary      获取用户拉黑列表
// @Tags         Admin/UserBlock
// @Security     BearerAuth
// @Param        page         query  int     false  "页码"
// @Param        pageSize     query  int     false  "每页数量"
// @Param        blockerId    query  int     false  "拉黑发起人ID"
// @Param        blockedId    query  int     false  "被拉黑人ID"
// @Param        blockerType  query  string  false  "发起人类型" Enums(user,player)
// @Param        blockedType  query  string  false  "被拉黑人类型" Enums(user,player)
// @Param        status       query  string  false  "状态" Enums(active,canceled,admin_canceled)
// @Produce      json
// @Success      200  {array}   model.UserBlock
// @Router       /admin/user-blocks [get]
func (h *UserBlockHandler) ListUserBlocks(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	blockerID, ok := QueryUint64PtrAndRespond(c, "blockerId", apierr.ErrInvalidUserID)
	if !ok {
		return
	}

	blockedID, ok := QueryUint64PtrAndRespond(c, "blockedId", apierr.ErrInvalidUserID)
	if !ok {
		return
	}

	var blockerType *model.BlockUserType
	if v := c.Query("blockerType"); v != "" {
		t := model.BlockUserType(v)
		blockerType = &t
	}

	var blockedType *model.BlockUserType
	if v := c.Query("blockedType"); v != "" {
		t := model.BlockUserType(v)
		blockedType = &t
	}

	var status *model.BlockStatus
	if v := c.Query("status"); v != "" {
		s := model.BlockStatus(v)
		status = &s
	}

	opts := repository.UserBlockListOptions{
		Page:        page,
		PageSize:    pageSize,
		BlockerID:   blockerID,
		BlockedID:   blockedID,
		BlockerType: blockerType,
		BlockedType: blockedType,
		Status:      status,
	}

	blocks, pagination, err := h.svc.ListPaged(c.Request.Context(), opts)
	if err != nil {
		respondError(c, err)
		return
	}
	respondList(c, blocks, pagination)
}

// GetUserBlock
// @Summary      获取用户拉黑详情
// @Tags         Admin/UserBlock
// @Security     BearerAuth
// @Param        id   path  int  true  "拉黑记录ID"
// @Produce      json
// @Success      200  {object}  model.UserBlock
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/user-blocks/{id} [get]
func (h *UserBlockHandler) GetUserBlock(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	block, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, block)
}

// CheckBlockStatus
// @Summary      检查两个用户之间的拉黑状态
// @Tags         Admin/UserBlock
// @Security     BearerAuth
// @Param        userId1  query  int  true  "用户1 ID"
// @Param        userId2  query  int  true  "用户2 ID"
// @Produce      json
// @Success      200  {object}  map[string]bool
// @Router       /admin/user-blocks/check [get]
func (h *UserBlockHandler) CheckBlockStatus(c *gin.Context) {
	userID1Ptr, ok := QueryUint64PtrAndRespond(c, "userId1", apierr.ErrInvalidUserID)
	if !ok {
		return
	}
	if userID1Ptr == nil {
		respondBadRequest(c, "userId1 is required")
		return
	}
	userID1 := *userID1Ptr

	userID2Ptr, ok := QueryUint64PtrAndRespond(c, "userId2", apierr.ErrInvalidUserID)
	if !ok {
		return
	}
	if userID2Ptr == nil {
		respondBadRequest(c, "userId2 is required")
		return
	}
	userID2 := *userID2Ptr

	// 检查双向拉黑状态
	isBlocked, err := h.svc.IsBlocked(c.Request.Context(), userID1, userID2)
	if err != nil {
		respondError(c, err)
		return
	}

	// 检查单向拉黑状态
	user1BlockedUser2, err := h.svc.IsBlockedBy(c.Request.Context(), userID1, userID2)
	if err != nil {
		respondError(c, err)
		return
	}

	user2BlockedUser1, err := h.svc.IsBlockedBy(c.Request.Context(), userID2, userID1)
	if err != nil {
		respondError(c, err)
		return
	}

	result := map[string]bool{
		"isBlocked":         isBlocked,
		"user1BlockedUser2": user1BlockedUser2,
		"user2BlockedUser1": user2BlockedUser1,
	}
	respondSuccess(c, result)
}

// GetUserBlocksByUser
// @Summary      获取用户的拉黑列表
// @Tags         Admin/UserBlock
// @Security     BearerAuth
// @Param        id   path  int  true  "用户ID"
// @Produce      json
// @Success      200  {array}   model.UserBlock
// @Router       /admin/users/{id}/blocks [get]
func (h *UserBlockHandler) GetUserBlocksByUser(c *gin.Context) {
	userID, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	blocks, err := h.svc.ListByBlocker(c.Request.Context(), userID, true)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, blocks)
}

// GetUserBlockedByList
// @Summary      获取拉黑该用户的列表
// @Tags         Admin/UserBlock
// @Security     BearerAuth
// @Param        id   path  int  true  "用户ID"
// @Produce      json
// @Success      200  {array}   model.UserBlock
// @Router       /admin/users/{id}/blocked-by [get]
func (h *UserBlockHandler) GetUserBlockedByList(c *gin.Context) {
	userID, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	blocks, err := h.svc.ListByBlocked(c.Request.Context(), userID, true)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, blocks)
}

// GetUserBlockStats
// @Summary      获取用户拉黑统计
// @Tags         Admin/UserBlock
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]int64
// @Router       /admin/user-blocks/stats [get]
func (h *UserBlockHandler) GetUserBlockStats(c *gin.Context) {
	stats, err := h.svc.GetStats(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}

	// 转换为 string key 的 map
	result := make(map[string]int64)
	for k, v := range stats {
		result[string(k)] = v
	}
	respondSuccess(c, result)
}

// ============================================================================
// 管理操作
// ============================================================================

// AdminUnblockRequest 管理员取消拉黑请求
type AdminUnblockRequest struct {
	Remark string `json:"remark"`
}

// AdminUnblock
// @Summary      管理员强制取消拉黑
// @Tags         Admin/UserBlock
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                   true  "拉黑记录ID"
// @Param        request  body  AdminUnblockRequest   false "备注信息"
// @Success      200  {object}  model.SuccessResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/user-blocks/{id}/unblock [post]
func (h *UserBlockHandler) AdminUnblock(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req AdminUnblockRequest
	_ = c.ShouldBindJSON(&req) // 可选参数

	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	if err := h.svc.AdminUnblock(c.Request.Context(), id, adminID, req.Remark); err != nil {
		respondError(c, err)
		return
	}
	respondMsg(c, "unblocked")
}

// BatchUnblockRequest 批量取消拉黑请求
type BatchUnblockRequest struct {
	BlockIDs []uint64 `json:"block_ids" binding:"required,min=1,max=100"`
	Remark   string   `json:"remark"`
}

// BatchUnblock
// @Summary      批量取消拉黑
// @Description  批量取消多个拉黑记录（管理员强制解除）
// @Tags         Admin/UserBlock
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchUnblockRequest  true  "批量取消请求"
// @Success      200  {object}  BatchOperationResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/user-blocks/batch/unblock [post]
func (h *UserBlockHandler) BatchUnblock(c *gin.Context) {
	var req BatchUnblockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.BlockIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("block_ids is required"))
		return
	}
	if len(req.BlockIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 blocks per batch"))
		return
	}

	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	result, err := h.svc.BatchUnblock(c.Request.Context(), req.BlockIDs, adminID, req.Remark)
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch unblock failed").WithDetails(err.Error()))
		return
	}

	response := convertBatchOperationResult(result)
	respondSuccess(c, response)
}

// DeleteUserBlock
// @Summary      删除拉黑记录
// @Tags         Admin/UserBlock
// @Security     BearerAuth
// @Param        id   path  int  true  "拉黑记录ID"
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/user-blocks/{id} [delete]
func (h *UserBlockHandler) DeleteUserBlock(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
}

// BatchDeleteRequest 批量删除请求
type BatchDeleteRequest struct {
	BlockIDs []uint64 `json:"block_ids" binding:"required,min=1,max=100"`
}

// BatchDelete
// @Summary      批量删除拉黑记录
// @Description  批量删除多个拉黑记录
// @Tags         Admin/UserBlock
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchDeleteRequest  true  "批量删除请求"
// @Success      200  {object}  BatchOperationResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/user-blocks/batch/delete [post]
func (h *UserBlockHandler) BatchDelete(c *gin.Context) {
	var req BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.BlockIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("block_ids is required"))
		return
	}
	if len(req.BlockIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 blocks per batch"))
		return
	}

	result, err := h.svc.BatchDelete(c.Request.Context(), req.BlockIDs)
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch delete failed").WithDetails(err.Error()))
		return
	}

	response := convertBatchOperationResult(result)
	respondSuccess(c, response)
}

// batchBlockRequest 批量拉黑请求
type batchBlockRequest struct {
	Blocks []blockInputItem `json:"blocks" binding:"required,min=1"`
}

// blockInputItem 单个拉黑项
type blockInputItem struct {
	BlockerID   uint64              `json:"blockerId" binding:"required"`
	BlockerType model.BlockUserType `json:"blockerType" binding:"required,oneof=user player"`
	BlockedID   uint64              `json:"blockedId" binding:"required"`
	BlockedType model.BlockUserType `json:"blockedType" binding:"required,oneof=user player"`
	Reason      string              `json:"reason"`
}

// BatchBlock
// @Summary      批量拉黑用户
// @Description  批量创建拉黑关系
// @Tags         Admin/UserBlock
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  batchBlockRequest  true  "批量拉黑请求"
// @Success      200  {object}  BatchOperationResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/user-blocks/batch [post]
func (h *UserBlockHandler) BatchBlock(c *gin.Context) {
	var req batchBlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.Blocks) == 0 {
		respondAPIError(c, apierr.BadRequest("blocks is required"))
		return
	}
	if len(req.Blocks) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 blocks per batch"))
		return
	}

	// Convert to service input
	items := make([]userblock.BlockInputItemForBatch, len(req.Blocks))
	for i, b := range req.Blocks {
		items[i] = userblock.BlockInputItemForBatch{
			BlockerID:   b.BlockerID,
			BlockerType: b.BlockerType,
			BlockedID:   b.BlockedID,
			BlockedType: b.BlockedType,
			Reason:      b.Reason,
		}
		_ = i // Avoid unused variable warning
	}

	result, err := h.svc.BatchBlock(c.Request.Context(), items)
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch block failed").WithDetails(err.Error()))
		return
	}

	// Convert service result to handler response format
	response := convertBatchOperationResult(result)
	respondSuccess(c, response)
}

// convertBatchOperationResult 转换 service 层的批量操作结果为 handler 层的响应格式
func convertBatchOperationResult(result *userblock.BatchOperationResult) BatchOperationResponse {
	response := BatchOperationResponse{
		SuccessCount: result.SuccessCount,
		FailedCount:  result.FailedCount,
		TotalCount:   result.TotalCount,
		FailedItems:  make([]BatchOperationError, 0, len(result.FailedIDs)),
	}

	for _, id := range result.FailedIDs {
		response.FailedItems = append(response.FailedItems, BatchOperationError{
			ID:      id,
			Message: "operation failed",
		})
	}

	return response
}

// RegisterUserBlockRoutes 注册用户拉黑管理路由
func RegisterUserBlockRoutes(router gin.IRouter, svc *userblock.UserBlockService, pm *middleware.PermissionMiddleware) {
	h := NewUserBlockHandler(svc)

	group := router.Group("/user-blocks")
	group.Use(pm.RequireAuth())
	{
		// 列表和查询
		group.GET("", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/user-blocks"), h.ListUserBlocks)
		group.GET("/stats", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/user-blocks/stats"), h.GetUserBlockStats)
		group.GET("/check", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/user-blocks/check"), h.CheckBlockStatus)
		group.GET("/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/user-blocks/:id"), h.GetUserBlock)

		// 管理操作
		group.POST("/:id/unblock", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/user-blocks/:id/unblock"), h.AdminUnblock)
		group.POST("/batch/unblock", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/user-blocks/batch/unblock"), h.BatchUnblock)
		group.DELETE("/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/user-blocks/:id"), h.DeleteUserBlock)
		group.POST("/batch/delete", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/user-blocks/batch/delete"), h.BatchDelete)
		group.POST("/batch", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/user-blocks/batch"), h.BatchBlock)
	}

	// 用户相关路由
	usersGroup := router.Group("/users")
	usersGroup.Use(pm.RequireAuth())
	{
		usersGroup.GET("/:id/blocks", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/:id/blocks"), h.GetUserBlocksByUser)
		usersGroup.GET("/:id/blocked-by", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/:id/blocked-by"), h.GetUserBlockedByList)
	}
}
