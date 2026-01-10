package admin

import (
	"strings"

	"github.com/gin-gonic/gin"

	mw "gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	svc "gamelink/internal/service/routingrule"
)

// RoutingRuleHandler 处理分流规则管理接口
// Requirements: 16.1-16.5
type RoutingRuleHandler struct {
	svc *svc.RoutingRuleService
}

// NewRoutingRuleHandler 创建Handler
func NewRoutingRuleHandler(svc *svc.RoutingRuleService) *RoutingRuleHandler {
	return &RoutingRuleHandler{svc: svc}
}

// ListRoutingRules
// @Summary      列出分流规则
// @Tags         Admin/RoutingRules
// @Security     BearerAuth
// @Param        page           query  int     false  "页码"
// @Param        pageSize       query  int     false  "每页数量"
// @Param        status         query  string  false  "状态过滤" Enums(active, inactive)
// @Param        targetEntityId query  int     false  "目标收款主体ID"
// @Param        keyword        query  string  false  "关键词搜索"
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.ListRoutingRulesResponse]
// @Router       /admin/routing-rules [get]
func (h *RoutingRuleHandler) ListRoutingRules(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	var targetEntityID *uint64
	if idStr := c.Query("targetEntityId"); idStr != "" {
		id, err := parseUint64(idStr)
		if err != nil {
			respondBadRequest(c, "invalid targetEntityId")
			return
		}
		targetEntityID = &id
	}

	req := &model.ListRoutingRulesRequest{
		Page:           page,
		PageSize:       pageSize,
		Status:         model.RuleStatus(strings.TrimSpace(c.Query("status"))),
		TargetEntityID: targetEntityID,
		Keyword:        strings.TrimSpace(c.Query("keyword")),
	}

	resp, err := h.svc.ListRules(c.Request.Context(), req)
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
	respondList(c, resp.Rules, pagination)
}

// GetRoutingRule
// @Summary      获取分流规则详情
// @Tags         Admin/RoutingRules
// @Security     BearerAuth
// @Param        id   path  int  true  "分流规则ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.RoutingRule]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/routing-rules/{id} [get]
func (h *RoutingRuleHandler) GetRoutingRule(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	rule, err := h.svc.GetRule(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, rule)
}

// CreateRoutingRule
// @Summary      创建分流规则
// @Tags         Admin/RoutingRules
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  model.CreateRoutingRuleRequest  true  "分流规则信息"
// @Success      201  {object}  model.APIResponse[model.RoutingRule]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/routing-rules [post]
func (h *RoutingRuleHandler) CreateRoutingRule(c *gin.Context) {
	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req model.CreateRoutingRuleRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	rule, err := h.svc.CreateRule(c.Request.Context(), &req, adminID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, rule)
}

// UpdateRoutingRule
// @Summary      更新分流规则
// @Tags         Admin/RoutingRules
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                              true  "分流规则ID"
// @Param        request  body  model.UpdateRoutingRuleRequest  true  "分流规则信息"
// @Success      200  {object}  model.APIResponse[model.RoutingRule]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/routing-rules/{id} [put]
func (h *RoutingRuleHandler) UpdateRoutingRule(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req model.UpdateRoutingRuleRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	rule, err := h.svc.UpdateRule(c.Request.Context(), id, &req, adminID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondUpdated(c, rule)
}

// DeleteRoutingRule
// @Summary      删除分流规则
// @Tags         Admin/RoutingRules
// @Security     BearerAuth
// @Param        id   path  int  true  "分流规则ID"
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/routing-rules/{id} [delete]
func (h *RoutingRuleHandler) DeleteRoutingRule(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	err := h.svc.DeleteRule(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccessWithMsg[any](c, "routing rule deleted", nil)
}

// ToggleRoutingRuleStatus
// @Summary      启用/禁用分流规则
// @Tags         Admin/RoutingRules
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                true  "分流规则ID"
// @Param        request  body  ToggleStatusPayload  true  "状态信息"
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/routing-rules/{id}/toggle [post]
func (h *RoutingRuleHandler) ToggleRoutingRuleStatus(c *gin.Context) {
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

	err := h.svc.ToggleRuleStatus(c.Request.Context(), id, payload.Enabled, adminID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccessWithMsg[any](c, "status updated", nil)
}

// GetRoutingRuleHistory
// @Summary      获取分流规则修改历史
// @Tags         Admin/RoutingRules
// @Security     BearerAuth
// @Param        id   path  int  true  "分流规则ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[[]model.RoutingRuleHistory]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/routing-rules/{id}/history [get]
func (h *RoutingRuleHandler) GetRoutingRuleHistory(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	histories, err := h.svc.GetRuleHistory(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, histories)
}

// SetDefaultEntityPayload 设置默认收款主体请求
type SetDefaultEntityPayload struct {
	EntityID uint64 `json:"entityId" binding:"required"`
}

// SetDefaultEntity
// @Summary      设置默认收款主体
// @Tags         Admin/RoutingRules
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  SetDefaultEntityPayload  true  "默认主体信息"
// @Success      200  {object}  model.SuccessResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/routing-rules/set-default [post]
func (h *RoutingRuleHandler) SetDefaultEntity(c *gin.Context) {
	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var payload SetDefaultEntityPayload
	if !ValidateAndRespond(c, &payload) {
		return
	}

	err := h.svc.SetDefaultEntity(c.Request.Context(), payload.EntityID, adminID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccessWithMsg[any](c, "default entity updated", nil)
}

// GetDefaultEntity
// @Summary      获取默认收款主体
// @Tags         Admin/RoutingRules
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.CollectionEntity]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/routing-rules/default-entity [get]
func (h *RoutingRuleHandler) GetDefaultEntity(c *gin.Context) {
	entity, err := h.svc.GetDefaultEntity(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, entity)
}

// TestRouting
// @Summary      测试分流规则
// @Tags         Admin/RoutingRules
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  model.RoutingTestRequest  true  "测试参数"
// @Success      200  {object}  model.APIResponse[model.RoutingTestResponse]
// @Router       /admin/routing-rules/test [post]
func (h *RoutingRuleHandler) TestRouting(c *gin.Context) {
	var req model.RoutingTestRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	result, err := h.svc.TestRouting(c.Request.Context(), &req)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, result)
}

// ReorderPrioritiesPayload 重新排序优先级请求
type ReorderPrioritiesPayload struct {
	RuleIDs []uint64 `json:"ruleIds" binding:"required,min=1"`
}

// ReorderPriorities
// @Summary      重新排序分流规则优先级
// @Tags         Admin/RoutingRules
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  ReorderPrioritiesPayload  true  "规则ID列表（按新优先级顺序）"
// @Success      200  {object}  model.SuccessResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/routing-rules/reorder [post]
func (h *RoutingRuleHandler) ReorderPriorities(c *gin.Context) {
	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var payload ReorderPrioritiesPayload
	if !ValidateAndRespond(c, &payload) {
		return
	}

	err := h.svc.ReorderPriorities(c.Request.Context(), payload.RuleIDs, adminID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccessWithMsg[any](c, "priorities reordered", nil)
}

// RegisterRoutingRuleRoutes 注册分流规则管理路由
// Requirements: 16.1-16.5
func RegisterRoutingRuleRoutes(router gin.IRouter, svc *svc.RoutingRuleService, pm *mw.PermissionMiddleware) {
	handler := NewRoutingRuleHandler(svc)

	// 分流规则管理
	router.GET("/routing-rules", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/routing-rules"), handler.ListRoutingRules)
	router.GET("/routing-rules/default-entity", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/routing-rules/default-entity"), handler.GetDefaultEntity)
	router.POST("/routing-rules/set-default", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/routing-rules/set-default"), handler.SetDefaultEntity)
	router.POST("/routing-rules/test", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/routing-rules/test"), handler.TestRouting)
	router.POST("/routing-rules/reorder", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/routing-rules/reorder"), handler.ReorderPriorities)
	router.GET("/routing-rules/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/routing-rules/:id"), handler.GetRoutingRule)
	router.POST("/routing-rules", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/routing-rules"), handler.CreateRoutingRule)
	router.PUT("/routing-rules/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/routing-rules/:id"), handler.UpdateRoutingRule)
	router.DELETE("/routing-rules/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/routing-rules/:id"), handler.DeleteRoutingRule)
	router.POST("/routing-rules/:id/toggle", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/routing-rules/:id/toggle"), handler.ToggleRoutingRuleStatus)
	router.GET("/routing-rules/:id/history", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/routing-rules/:id/history"), handler.GetRoutingRuleHistory)

	// 批量操作
	router.POST("/routing-rules/batch/status", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/routing-rules/batch/status"), handler.BatchUpdateRoutingRuleStatus)
	router.POST("/routing-rules/batch/delete", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/routing-rules/batch/delete"), handler.BatchDeleteRoutingRules)
}
