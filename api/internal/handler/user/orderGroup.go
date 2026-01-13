package user

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/repository/ordergroup"
	"gamelink/internal/service/order"
)

// OrderGroupHandler 主订单处理器
type OrderGroupHandler struct {
	orderSvc  *order.OrderService
	groupRepo ordergroup.Repository
}

// NewOrderGroupHandler 创建主订单处理器
func NewOrderGroupHandler(orderSvc *order.OrderService, groupRepo ordergroup.Repository) *OrderGroupHandler {
	return &OrderGroupHandler{
		orderSvc:  orderSvc,
		groupRepo: groupRepo,
	}
}

// RegisterRoutes 注册路由
func (h *OrderGroupHandler) RegisterRoutes(router gin.IRouter) {
	group := router.Group("/order-groups")
	{
		group.GET("", h.listMyOrderGroups)
		group.GET("/:id", h.getOrderGroup)
		group.GET("/:id/sub-orders", h.getSubOrders)
		group.POST("/:subOrderId/transfer", h.transferSubOrder)
		group.POST("/batch-transfer", h.batchTransferSubOrders)
	}
}

// listMyOrderGroups 获取我的主订单列表
// @Summary      获取我的主订单列表
// @Description  用户端获取主订单列表（聚合视图）
// @Tags         User - OrderGroup
// @Security     BearerAuth
// @Produce      json
// @Param        status    query   string  false  "状态筛选"
// @Param        page      query   int     false  "页码"
// @Param        pageSize  query   int     false  "每页数量"
// @Success      200  {object}  OrderGroupListResponse
// @Router       /user/order-groups [get]
func (h *OrderGroupHandler) listMyOrderGroups(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		respondError(c, 401, "未登录")
		return
	}

	var opts ordergroup.ListOptions
	if page, _ := strconv.Atoi(c.Query("page")); page > 0 {
		opts.Page = page
	}
	if pageSize, _ := strconv.Atoi(c.Query("pageSize")); pageSize > 0 {
		opts.PageSize = pageSize
	}

	groups, total, err := h.groupRepo.ListByUser(c.Request.Context(), userID, opts)
	if err != nil {
		respondError(c, 500, "查询失败")
		return
	}

	respondSuccess(c, "获取成功", gin.H{
		"items": groups,
		"total": total,
		"page":  opts.Page,
	})
}

// getOrderGroup 获取主订单详情
// @Summary      获取主订单详情
// @Description  获取主订单及其子订单详情
// @Tags         User - OrderGroup
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  int  true  "主订单ID"
// @Success      200  {object}  OrderGroupDTO
// @Router       /user/order-groups/{id} [get]
func (h *OrderGroupHandler) getOrderGroup(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		respondError(c, 401, "未登录")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, 400, "无效的ID")
		return
	}

	group, err := h.groupRepo.GetWithSubOrders(c.Request.Context(), id)
	if err != nil {
		respondError(c, 404, "订单不存在")
		return
	}

	// 验证权限
	if group.UserID != userID {
		respondError(c, 403, "无权访问")
		return
	}

	respondSuccess(c, "获取成功", group)
}

// getSubOrders 获取子订单列表
// @Summary      获取子订单列表
// @Description  获取主订单下的所有子订单
// @Tags         User - OrderGroup
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  int  true  "主订单ID"
// @Success      200  {array}   OrderGroupDTO
// @Router       /user/order-groups/{id}/sub-orders [get]
func (h *OrderGroupHandler) getSubOrders(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		respondError(c, 401, "未登录")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, 400, "无效的ID")
		return
	}

	group, err := h.groupRepo.GetWithSubOrders(c.Request.Context(), id)
	if err != nil {
		respondError(c, 404, "订单不存在")
		return
	}

	if group.UserID != userID {
		respondError(c, 403, "无权访问")
		return
	}

	respondSuccess(c, "获取成功", group.SubOrders)
}

// transferSubOrder 转单
// @Summary      转单
// @Description  将子订单转给另一个陪玩师
// @Tags         User - OrderGroup
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        subOrderId  path  int  true  "子订单ID"
// @Param        request     body  order.TransferSubOrderRequest  true  "转单请求"
// @Success      200  {object}  order.TransferSubOrderResponse
// @Router       /user/order-groups/{subOrderId}/transfer [post]
func (h *OrderGroupHandler) transferSubOrder(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		respondError(c, 401, "未登录")
		return
	}

	subOrderID, err := strconv.ParseUint(c.Param("subOrderId"), 10, 64)
	if err != nil {
		respondError(c, 400, "无效的子订单ID")
		return
	}

	var req order.TransferSubOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, 400, "参数错误")
		return
	}
	req.SubOrderID = subOrderID

	result, err := h.orderSvc.TransferSubOrder(c.Request.Context(), userID, req)
	if err != nil {
		respondAPIError(c, err)
		return
	}

	respondSuccess(c, "转单成功", result)
}

// batchTransferSubOrders 批量转单
// @Summary      批量转单
// @Description  将多个子订单转给另一个陪玩师
// @Tags         User - OrderGroup
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  order.BatchTransferRequest  true  "批量转单请求"
// @Success      200  {object}  order.BatchTransferResponse
// @Router       /user/order-groups/batch-transfer [post]
func (h *OrderGroupHandler) batchTransferSubOrders(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		respondError(c, 401, "未登录")
		return
	}

	var req order.BatchTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, 400, "参数错误")
		return
	}

	result, err := h.orderSvc.BatchTransferSubOrders(c.Request.Context(), userID, req)
	if err != nil {
		respondAPIError(c, err)
		return
	}

	respondSuccess(c, "批量转单完成", result)
}

// OrderGroupListResponse 主订单列表响应
type OrderGroupListResponse struct {
	Items []OrderGroupDTO `json:"items"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
}

// OrderGroupDTO 主订单DTO
type OrderGroupDTO struct {
	ID              uint64   `json:"id"`
	GroupNo         string   `json:"groupNo"`
	GameName        string   `json:"gameName"`
	TotalPriceCents int64    `json:"totalPriceCents"`
	TotalHours      int      `json:"totalHours"`
	CompletedHours  int      `json:"completedHours"`
	Status          string   `json:"status"`
	Title           string   `json:"title"`
	CreatedAt       string   `json:"createdAt"`
	Players         []string `json:"players"` // 服务过的陪玩师列表
}
