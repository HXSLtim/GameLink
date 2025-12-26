package admin

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service/ordertimeout"
	"gamelink/pkg/apierr"
)

// OrderTimeoutHandler 订单超时管理接口
type OrderTimeoutHandler struct {
	svc *ordertimeout.OrderTimeoutService
}

// NewOrderTimeoutHandler 创建Handler
func NewOrderTimeoutHandler(svc *ordertimeout.OrderTimeoutService) *OrderTimeoutHandler {
	return &OrderTimeoutHandler{svc: svc}
}

// ============================================================================
// 配置管理
// ============================================================================

// ListConfigs
// @Summary      获取订单超时配置列表
// @Tags         Admin/OrderTimeout
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  model.APIResponse[[]model.OrderTimeoutConfig]
// @Router       /admin/order-timeout/configs [get]
func (h *OrderTimeoutHandler) ListConfigs(c *gin.Context) {
	configs, err := h.svc.ListConfigs(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, configs)
}

// GetConfig
// @Summary      获取订单超时配置
// @Tags         Admin/OrderTimeout
// @Security     BearerAuth
// @Param        key   path  string  true  "配置键"
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.OrderTimeoutConfig]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/order-timeout/configs/{key} [get]
func (h *OrderTimeoutHandler) GetConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		respondBadRequest(c, "config key is required")
		return
	}

	config, err := h.svc.GetConfig(c.Request.Context(), key)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, config)
}

// SaveConfigRequest 保存配置请求
type SaveConfigRequest struct {
	Key         string `json:"key" binding:"required"`
	Value       string `json:"value" binding:"required"`
	Description string `json:"description"`
}

// SaveConfig
// @Summary      保存订单超时配置
// @Tags         Admin/OrderTimeout
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  SaveConfigRequest  true  "配置信息"
// @Success      200  {object}  model.SuccessResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/order-timeout/configs [post]
func (h *OrderTimeoutHandler) SaveConfig(c *gin.Context) {
	var req SaveConfigRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	if err := h.svc.SaveConfig(c.Request.Context(), req.Key, req.Value, req.Description); err != nil {
		respondError(c, err)
		return
	}
	respondMsg(c, "config saved")
}

// DeleteConfig
// @Summary      删除订单超时配置
// @Tags         Admin/OrderTimeout
// @Security     BearerAuth
// @Param        key   path  string  true  "配置键"
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/order-timeout/configs/{key} [delete]
func (h *OrderTimeoutHandler) DeleteConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		respondBadRequest(c, "config key is required")
		return
	}

	if err := h.svc.SaveConfig(c.Request.Context(), key, "", ""); err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
}

// ============================================================================
// 超时日志
// ============================================================================

// ListTimeoutLogs
// @Summary      获取订单超时日志列表
// @Tags         Admin/OrderTimeout
// @Security     BearerAuth
// @Param        page         query  int     false  "页码"
// @Param        pageSize     query  int     false  "每页数量"
// @Param        orderId      query  int     false  "订单ID筛选"
// @Param        timeoutType  query  string  false  "超时类型" Enums(payment_timeout,accept_timeout)
// @Param        action       query  string  false  "处理动作" Enums(canceled,refunded,notified)
// @Produce      json
// @Success      200  {object}  model.APIResponse[[]model.OrderTimeoutLog]
// @Router       /admin/order-timeout/logs [get]
func (h *OrderTimeoutHandler) ListTimeoutLogs(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	orderID, ok := QueryUint64PtrAndRespond(c, "orderId", apierr.ErrInvalidOrderID)
	if !ok {
		return
	}

	var timeoutType *model.OrderTimeoutType
	if v := c.Query("timeoutType"); v != "" {
		t := model.OrderTimeoutType(v)
		timeoutType = &t
	}

	var action *model.OrderTimeoutAction
	if v := c.Query("action"); v != "" {
		a := model.OrderTimeoutAction(v)
		action = &a
	}

	opts := repository.OrderTimeoutLogListOptions{
		Page:        page,
		PageSize:    pageSize,
		OrderID:     orderID,
		TimeoutType: timeoutType,
		Action:      action,
	}

	logs, pagination, err := h.svc.ListLogsPaged(c.Request.Context(), opts)
	if err != nil {
		respondError(c, err)
		return
	}
	respondList(c, logs, pagination)
}

// GetTimeoutLog
// @Summary      获取订单超时日志详情
// @Tags         Admin/OrderTimeout
// @Security     BearerAuth
// @Param        id   path  int  true  "日志ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.OrderTimeoutLog]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/order-timeout/logs/{id} [get]
func (h *OrderTimeoutHandler) GetTimeoutLog(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	log, err := h.svc.GetLog(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, log)
}

// GetTimeoutLogsByOrder
// @Summary      获取订单的超时日志
// @Tags         Admin/OrderTimeout
// @Security     BearerAuth
// @Param        orderId   path  int  true  "订单ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[[]model.OrderTimeoutLog]
// @Router       /admin/orders/{orderId}/timeout-logs [get]
func (h *OrderTimeoutHandler) GetTimeoutLogsByOrder(c *gin.Context) {
	orderID, ok := ParseIDAndRespond(c, "orderId")
	if !ok {
		return
	}

	logs, err := h.svc.ListLogsByOrderID(c.Request.Context(), orderID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, logs)
}

// GetTimeoutLogStats
// @Summary      获取订单超时日志统计
// @Tags         Admin/OrderTimeout
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  model.APIResponse[map[string]int64]
// @Router       /admin/order-timeout/logs/stats [get]
func (h *OrderTimeoutHandler) GetTimeoutLogStats(c *gin.Context) {
	stats, err := h.svc.GetLogStats(c.Request.Context())
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
// 客服分配
// ============================================================================

// ListServiceAssignments
// @Summary      获取客服分配列表
// @Tags         Admin/OrderTimeout
// @Security     BearerAuth
// @Param        page           query  int     false  "页码"
// @Param        pageSize       query  int     false  "每页数量"
// @Param        orderId        query  int     false  "订单ID筛选"
// @Param        serviceUserId  query  int     false  "客服ID筛选"
// @Param        status         query  string  false  "状态筛选" Enums(assigned,joined,left,completed)
// @Param        assignType     query  string  false  "分配方式" Enums(auto,manual)
// @Produce      json
// @Success      200  {object}  model.APIResponse[[]model.OrderServiceAssignment]
// @Router       /admin/order-timeout/assignments [get]
func (h *OrderTimeoutHandler) ListServiceAssignments(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	orderID, ok := QueryUint64PtrAndRespond(c, "orderId", apierr.ErrInvalidOrderID)
	if !ok {
		return
	}

	serviceUserID, ok := QueryUint64PtrAndRespond(c, "serviceUserId", apierr.ErrInvalidUserID)
	if !ok {
		return
	}

	var status *model.ServiceAssignmentStatus
	if v := c.Query("status"); v != "" {
		s := model.ServiceAssignmentStatus(v)
		status = &s
	}

	opts := repository.ServiceAssignmentListOptions{
		Page:          page,
		PageSize:      pageSize,
		OrderID:       orderID,
		ServiceUserID: serviceUserID,
		Status:        status,
		AssignType:    c.Query("assignType"),
	}

	assignments, pagination, err := h.svc.ListAssignmentsPaged(c.Request.Context(), opts)
	if err != nil {
		respondError(c, err)
		return
	}
	respondList(c, assignments, pagination)
}

// GetServiceAssignment
// @Summary      获取客服分配详情
// @Tags         Admin/OrderTimeout
// @Security     BearerAuth
// @Param        id   path  int  true  "分配记录ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.OrderServiceAssignment]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/order-timeout/assignments/{id} [get]
func (h *OrderTimeoutHandler) GetServiceAssignment(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	assignment, err := h.svc.GetAssignment(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, assignment)
}

// AssignServiceRequest 分配客服请求
type AssignServiceRequest struct {
	OrderID       uint64  `json:"orderId" binding:"required"`
	ServiceUserID uint64  `json:"serviceUserId" binding:"required"`
	ChatGroupID   *uint64 `json:"chatGroupId"`
	Remark        string  `json:"remark"`
}

// AssignService
// @Summary      手动分配客服
// @Tags         Admin/OrderTimeout
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  AssignServiceRequest  true  "分配信息"
// @Success      201  {object}  model.APIResponse[model.OrderServiceAssignment]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/order-timeout/assignments [post]
func (h *OrderTimeoutHandler) AssignService(c *gin.Context) {
	var req AssignServiceRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	assignment, err := h.svc.AssignService(c.Request.Context(), ordertimeout.AssignServiceInput{
		OrderID:       req.OrderID,
		ServiceUserID: req.ServiceUserID,
		ChatGroupID:   req.ChatGroupID,
		AssignType:    "manual",
		Remark:        req.Remark,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, assignment)
}

// UpdateAssignmentStatusRequest 更新分配状态请求
type UpdateAssignmentStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=assigned joined left completed"`
}

// UpdateAssignmentStatus
// @Summary      更新客服分配状态
// @Tags         Admin/OrderTimeout
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                            true  "分配记录ID"
// @Param        request  body  UpdateAssignmentStatusRequest  true  "状态信息"
// @Success      200  {object}  model.SuccessResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/order-timeout/assignments/{id}/status [put]
func (h *OrderTimeoutHandler) UpdateAssignmentStatus(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req UpdateAssignmentStatusRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	if err := h.svc.UpdateAssignmentStatus(c.Request.Context(), id, model.ServiceAssignmentStatus(req.Status)); err != nil {
		respondError(c, err)
		return
	}
	respondMsg(c, "status updated")
}

// DeleteServiceAssignment
// @Summary      删除客服分配记录
// @Tags         Admin/OrderTimeout
// @Security     BearerAuth
// @Param        id   path  int  true  "分配记录ID"
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/order-timeout/assignments/{id} [delete]
func (h *OrderTimeoutHandler) DeleteServiceAssignment(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	if err := h.svc.DeleteAssignment(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
}

// GetAssignmentByOrder
// @Summary      获取订单的客服分配
// @Tags         Admin/OrderTimeout
// @Security     BearerAuth
// @Param        orderId   path  int  true  "订单ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.OrderServiceAssignment]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/orders/{orderId}/service-assignment [get]
func (h *OrderTimeoutHandler) GetAssignmentByOrder(c *gin.Context) {
	orderID, ok := ParseIDAndRespond(c, "orderId")
	if !ok {
		return
	}

	assignment, err := h.svc.GetAssignmentByOrderID(c.Request.Context(), orderID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, assignment)
}

// GetAssignmentStats
// @Summary      获取客服分配统计
// @Tags         Admin/OrderTimeout
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  model.APIResponse[map[string]int64]
// @Router       /admin/order-timeout/assignments/stats [get]
func (h *OrderTimeoutHandler) GetAssignmentStats(c *gin.Context) {
	stats, err := h.svc.GetAssignmentStats(c.Request.Context())
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

// RegisterOrderTimeoutRoutes 注册订单超时管理路由
func RegisterOrderTimeoutRoutes(router gin.IRouter, svc *ordertimeout.OrderTimeoutService, pm *middleware.PermissionMiddleware) {
	h := NewOrderTimeoutHandler(svc)

	group := router.Group("/order-timeout")
	group.Use(pm.RequireAuth())
	{
		// 配置管理
		group.GET("/configs", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/order-timeout/configs"), h.ListConfigs)
		group.GET("/configs/:key", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/order-timeout/configs/:key"), h.GetConfig)
		group.POST("/configs", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/order-timeout/configs"), h.SaveConfig)
		group.DELETE("/configs/:key", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/order-timeout/configs/:key"), h.DeleteConfig)

		// 超时日志
		group.GET("/logs", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/order-timeout/logs"), h.ListTimeoutLogs)
		group.GET("/logs/stats", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/order-timeout/logs/stats"), h.GetTimeoutLogStats)
		group.GET("/logs/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/order-timeout/logs/:id"), h.GetTimeoutLog)

		// 客服分配
		group.GET("/assignments", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/order-timeout/assignments"), h.ListServiceAssignments)
		group.GET("/assignments/stats", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/order-timeout/assignments/stats"), h.GetAssignmentStats)
		group.GET("/assignments/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/order-timeout/assignments/:id"), h.GetServiceAssignment)
		group.POST("/assignments", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/order-timeout/assignments"), h.AssignService)
		group.PUT("/assignments/:id/status", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/order-timeout/assignments/:id/status"), h.UpdateAssignmentStatus)
		group.DELETE("/assignments/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/order-timeout/assignments/:id"), h.DeleteServiceAssignment)
	}

	// 订单相关路由（挂在 /orders 下）
	ordersGroup := router.Group("/orders")
	ordersGroup.Use(pm.RequireAuth())
	{
		ordersGroup.GET("/:orderId/timeout-logs", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/orders/:orderId/timeout-logs"), h.GetTimeoutLogsByOrder)
		ordersGroup.GET("/:orderId/service-assignment", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/orders/:orderId/service-assignment"), h.GetAssignmentByOrder)
	}
}
