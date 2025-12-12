package admin

import (
	"strings"

	"github.com/gin-gonic/gin"

	mw "gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	svc "gamelink/internal/service/collectionentity"
)

// CollectionEntityHandler 处理收款主体管理接口
// Requirements: 15.1-15.5
type CollectionEntityHandler struct {
	svc *svc.CollectionEntityService
}

// NewCollectionEntityHandler 创建Handler
func NewCollectionEntityHandler(svc *svc.CollectionEntityService) *CollectionEntityHandler {
	return &CollectionEntityHandler{svc: svc}
}

// ListCollectionEntities
// @Summary      列出收款主体
// @Tags         Admin/CollectionEntities
// @Security     BearerAuth
// @Param        page       query  int     false  "页码"
// @Param        pageSize   query  int     false  "每页数量"
// @Param        status     query  string  false  "状态过滤" Enums(active, inactive)
// @Param        keyword    query  string  false  "关键词搜索"
// @Param        sortBy     query  string  false  "排序字段" Enums(name, created_at, total_collection_cents)
// @Param        sortOrder  query  string  false  "排序方向" Enums(asc, desc)
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.ListCollectionEntitiesResponse]
// @Router       /admin/collection-entities [get]
func (h *CollectionEntityHandler) ListCollectionEntities(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	req := &model.ListCollectionEntitiesRequest{
		Page:      page,
		PageSize:  pageSize,
		Status:    model.EntityStatus(strings.TrimSpace(c.Query("status"))),
		Keyword:   strings.TrimSpace(c.Query("keyword")),
		SortBy:    strings.TrimSpace(c.Query("sortBy")),
		SortOrder: strings.TrimSpace(c.Query("sortOrder")),
	}

	resp, err := h.svc.ListEntities(c.Request.Context(), req)
	if err != nil {
		respondError(c, err)
		return
	}

	totalPages := int((resp.Total + int64(resp.PageSize) - 1) / int64(resp.PageSize))
	pagination := &model.Pagination{
		Page:       resp.Page,
		PageSize:   resp.PageSize,
		Total:      int(resp.Total),
		TotalPages: totalPages,
		HasNext:    resp.Page < totalPages,
		HasPrev:    resp.Page > 1,
	}
	respondList(c, resp.Entities, pagination)
}

// GetCollectionEntity
// @Summary      获取收款主体详情
// @Tags         Admin/CollectionEntities
// @Security     BearerAuth
// @Param        id   path  int  true  "收款主体ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.CollectionEntity]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/collection-entities/{id} [get]
func (h *CollectionEntityHandler) GetCollectionEntity(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	entity, err := h.svc.GetEntity(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, entity)
}

// CreateCollectionEntity
// @Summary      创建收款主体
// @Tags         Admin/CollectionEntities
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  model.CreateCollectionEntityRequest  true  "收款主体信息"
// @Success      201  {object}  model.APIResponse[model.CollectionEntity]
// @Failure      400  {object}  model.ErrorResponse
// @Failure      409  {object}  model.ErrorResponse
// @Router       /admin/collection-entities [post]
func (h *CollectionEntityHandler) CreateCollectionEntity(c *gin.Context) {
	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req model.CreateCollectionEntityRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	entity, err := h.svc.CreateEntity(c.Request.Context(), &req, adminID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, entity)
}

// UpdateCollectionEntity
// @Summary      更新收款主体
// @Tags         Admin/CollectionEntities
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                                   true  "收款主体ID"
// @Param        request  body  model.UpdateCollectionEntityRequest  true  "收款主体信息"
// @Success      200  {object}  model.APIResponse[model.CollectionEntity]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/collection-entities/{id} [put]
func (h *CollectionEntityHandler) UpdateCollectionEntity(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req model.UpdateCollectionEntityRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	entity, err := h.svc.UpdateEntity(c.Request.Context(), id, &req, adminID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondUpdated(c, entity)
}

// ToggleCollectionEntityStatus
// @Summary      启用/禁用收款主体
// @Tags         Admin/CollectionEntities
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                true  "收款主体ID"
// @Param        request  body  ToggleStatusPayload  true  "状态信息"
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/collection-entities/{id}/toggle [post]
func (h *CollectionEntityHandler) ToggleCollectionEntityStatus(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var payload ToggleStatusPayload
	if !ValidateAndRespond(c, &payload) {
		return
	}

	err := h.svc.ToggleEntityStatus(c.Request.Context(), id, payload.Enabled, adminID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccessWithMsg[any](c, "status updated", nil)
}

// GetCollectionEntityHistory
// @Summary      获取收款主体修改历史
// @Tags         Admin/CollectionEntities
// @Security     BearerAuth
// @Param        id   path  int  true  "收款主体ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[[]model.CollectionEntityHistory]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/collection-entities/{id}/history [get]
func (h *CollectionEntityHandler) GetCollectionEntityHistory(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	histories, err := h.svc.GetEntityHistory(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, histories)
}

// SetDefaultCollectionEntity
// @Summary      设置默认收款主体
// @Tags         Admin/CollectionEntities
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "收款主体ID"
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/collection-entities/{id}/set-default [post]
func (h *CollectionEntityHandler) SetDefaultCollectionEntity(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	err := h.svc.SetDefaultEntity(c.Request.Context(), id, adminID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccessWithMsg[any](c, "default entity updated", nil)
}

// ConfigurePaymentChannel
// @Summary      配置支付渠道
// @Tags         Admin/CollectionEntities
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                                true  "收款主体ID"
// @Param        request  body  model.ConfigurePaymentChannelRequest  true  "支付渠道配置"
// @Success      200  {object}  model.APIResponse[model.PaymentChannelConfig]
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/collection-entities/{id}/channels [post]
func (h *CollectionEntityHandler) ConfigurePaymentChannel(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req model.ConfigurePaymentChannelRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	channel, err := h.svc.ConfigurePaymentChannel(c.Request.Context(), id, &req, adminID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, channel)
}

// ListPaymentChannels
// @Summary      获取收款主体的支付渠道列表
// @Tags         Admin/CollectionEntities
// @Security     BearerAuth
// @Param        id   path  int  true  "收款主体ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[[]model.PaymentChannelConfig]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/collection-entities/{id}/channels [get]
func (h *CollectionEntityHandler) ListPaymentChannels(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	channels, err := h.svc.ListPaymentChannels(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, channels)
}

// UpdatePaymentChannel
// @Summary      更新支付渠道配置
// @Tags         Admin/CollectionEntities
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id         path  int                                true  "收款主体ID"
// @Param        channelId  path  int                                true  "支付渠道ID"
// @Param        request    body  model.ConfigurePaymentChannelRequest  true  "支付渠道配置"
// @Success      200  {object}  model.APIResponse[model.PaymentChannelConfig]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/collection-entities/{id}/channels/{channelId} [put]
func (h *CollectionEntityHandler) UpdatePaymentChannel(c *gin.Context) {
	_, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	channelID, ok := ParseIDAndRespond(c, "channelId")
	if !ok {
		return
	}

	var req model.ConfigurePaymentChannelRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	channel, err := h.svc.UpdatePaymentChannel(c.Request.Context(), channelID, &req)
	if err != nil {
		respondError(c, err)
		return
	}
	respondUpdated(c, channel)
}

// DeletePaymentChannel
// @Summary      删除支付渠道配置
// @Tags         Admin/CollectionEntities
// @Security     BearerAuth
// @Param        id         path  int  true  "收款主体ID"
// @Param        channelId  path  int  true  "支付渠道ID"
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/collection-entities/{id}/channels/{channelId} [delete]
func (h *CollectionEntityHandler) DeletePaymentChannel(c *gin.Context) {
	_, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	channelID, ok := ParseIDAndRespond(c, "channelId")
	if !ok {
		return
	}

	err := h.svc.DeletePaymentChannel(c.Request.Context(), channelID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccessWithMsg[any](c, "payment channel deleted", nil)
}

// TogglePaymentChannelStatus
// @Summary      启用/禁用支付渠道
// @Tags         Admin/CollectionEntities
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id         path  int                true  "收款主体ID"
// @Param        channelId  path  int                true  "支付渠道ID"
// @Param        request    body  ToggleStatusPayload  true  "状态信息"
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/collection-entities/{id}/channels/{channelId}/toggle [post]
func (h *CollectionEntityHandler) TogglePaymentChannelStatus(c *gin.Context) {
	_, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	channelID, ok := ParseIDAndRespond(c, "channelId")
	if !ok {
		return
	}

	var payload ToggleStatusPayload
	if !ValidateAndRespond(c, &payload) {
		return
	}

	err := h.svc.TogglePaymentChannelStatus(c.Request.Context(), channelID, payload.Enabled)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccessWithMsg[any](c, "channel status updated", nil)
}

// RegisterCollectionEntityRoutes 注册收款主体管理路由
// Requirements: 15.1-15.5
func RegisterCollectionEntityRoutes(router gin.IRouter, svc *svc.CollectionEntityService, pm *mw.PermissionMiddleware) {
	handler := NewCollectionEntityHandler(svc)

	// 收款主体管理
	router.GET("/collection-entities", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/collection-entities"), handler.ListCollectionEntities)
	router.GET("/collection-entities/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/collection-entities/:id"), handler.GetCollectionEntity)
	router.POST("/collection-entities", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/collection-entities"), handler.CreateCollectionEntity)
	router.PUT("/collection-entities/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/collection-entities/:id"), handler.UpdateCollectionEntity)
	router.POST("/collection-entities/:id/toggle", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/collection-entities/:id/toggle"), handler.ToggleCollectionEntityStatus)
	router.GET("/collection-entities/:id/history", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/collection-entities/:id/history"), handler.GetCollectionEntityHistory)
	router.POST("/collection-entities/:id/set-default", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/collection-entities/:id/set-default"), handler.SetDefaultCollectionEntity)

	// 支付渠道配置
	router.GET("/collection-entities/:id/channels", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/collection-entities/:id/channels"), handler.ListPaymentChannels)
	router.POST("/collection-entities/:id/channels", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/collection-entities/:id/channels"), handler.ConfigurePaymentChannel)
	router.PUT("/collection-entities/:id/channels/:channelId", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/collection-entities/:id/channels/:channelId"), handler.UpdatePaymentChannel)
	router.DELETE("/collection-entities/:id/channels/:channelId", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/collection-entities/:id/channels/:channelId"), handler.DeletePaymentChannel)
	router.POST("/collection-entities/:id/channels/:channelId/toggle", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/collection-entities/:id/channels/:channelId/toggle"), handler.TogglePaymentChannelStatus)
}
