package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository/ordergroup"
	"gamelink/internal/service/order"
	"gamelink/pkg/apierr"
)

// RegisterOrderGroupRoutes 注册管理端主订单路由
func RegisterOrderGroupRoutes(router gin.IRouter, orderSvc *order.OrderService, groupRepo ordergroup.Repository) {
	group := router.Group("/order-groups")
	{
		group.GET("", listOrderGroupsHandler(groupRepo))
		group.GET("/:id", getOrderGroupHandler(groupRepo))
		group.GET("/:id/sub-orders", getOrderGroupSubOrdersHandler(groupRepo))
		group.POST("/:subOrderId/transfer", transferSubOrderHandler(orderSvc))
		group.POST("/batch-transfer", batchTransferHandler(orderSvc))
	}
}

// listOrderGroupsHandler 获取主订单列表
// @Summary      获取主订单列表
// @Description  管理端获取所有主订单列表
// @Tags         Admin - OrderGroup
// @Security     BearerAuth
// @Produce      json
// @Param        userId    query   int     false  "用户ID筛选"
// @Param        status    query   string  false  "状态筛选"
// @Param        page      query   int     false  "页码"
// @Param        pageSize  query   int     false  "每页数量"
// @Success      200  {object}  OrderGroupListResponse
// @Router       /admin/order-groups [get]
func listOrderGroupsHandler(repo ordergroup.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var opts ordergroup.ListOptions

		if userID, _ := strconv.ParseUint(c.Query("userId"), 10, 64); userID > 0 {
			opts.UserID = &userID
		}
		if status := c.Query("status"); status != "" {
			s := model.OrderGroupStatus(status)
			opts.Status = &s
		}
		if page, _ := strconv.Atoi(c.Query("page")); page > 0 {
			opts.Page = page
		}
		if pageSize, _ := strconv.Atoi(c.Query("pageSize")); pageSize > 0 {
			opts.PageSize = pageSize
		}

		groups, total, err := repo.List(c.Request.Context(), opts)
		if err != nil {
			respondError(c, apierr.InternalError("查询失败"))
			return
		}

		respondSuccess(c, OrderGroupListResponse{
			Items: groups,
			Total: total,
			Page:  opts.Page,
		})
	}
}

// getOrderGroupHandler 获取主订单详情
// @Summary      获取主订单详情
// @Description  获取主订单及其子订单详情
// @Tags         Admin - OrderGroup
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  int  true  "主订单ID"
// @Success      200  {object}  model.OrderGroup
// @Router       /admin/order-groups/{id} [get]
func getOrderGroupHandler(repo ordergroup.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			respondBadRequest(c, "无效的ID")
			return
		}

		group, err := repo.GetWithSubOrders(c.Request.Context(), id)
		if err != nil {
			respondError(c, apierr.NotFound("订单不存在"))
			return
		}

		respondSuccess(c, group)
	}
}

// getOrderGroupSubOrdersHandler 获取子订单列表
// @Summary      获取子订单列表
// @Description  获取主订单下的所有子订单
// @Tags         Admin - OrderGroup
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  int  true  "主订单ID"
// @Success      200  {array}   model.Order
// @Router       /admin/order-groups/{id}/sub-orders [get]
func getOrderGroupSubOrdersHandler(repo ordergroup.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			respondBadRequest(c, "无效的ID")
			return
		}

		group, err := repo.GetWithSubOrders(c.Request.Context(), id)
		if err != nil {
			respondError(c, apierr.NotFound("订单不存在"))
			return
		}

		respondSuccess(c, group.SubOrders)
	}
}

// transferSubOrderHandler 管理员转单
// @Summary      管理员转单
// @Description  管理员将子订单转给另一个陪玩师
// @Tags         Admin - OrderGroup
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        subOrderId  path  int  true  "子订单ID"
// @Param        request     body  order.TransferSubOrderRequest  true  "转单请求"
// @Success      200  {object}  order.TransferSubOrderResponse
// @Router       /admin/order-groups/{subOrderId}/transfer [post]
func transferSubOrderHandler(svc *order.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := getAdminID(c)

		subOrderID, err := strconv.ParseUint(c.Param("subOrderId"), 10, 64)
		if err != nil {
			respondBadRequest(c, "无效的子订单ID")
			return
		}

		var req order.TransferSubOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			respondBadRequest(c, err.Error())
			return
		}
		req.SubOrderID = subOrderID

		result, err := svc.TransferSubOrder(c.Request.Context(), adminID, req)
		if err != nil {
			respondError(c, err)
			return
		}

		respondSuccess(c, result)
	}
}

// batchTransferHandler 管理员批量转单
// @Summary      管理员批量转单
// @Description  管理员将多个子订单转给另一个陪玩师
// @Tags         Admin - OrderGroup
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  order.BatchTransferRequest  true  "批量转单请求"
// @Success      200  {object}  order.BatchTransferResponse
// @Router       /admin/order-groups/batch-transfer [post]
func batchTransferHandler(svc *order.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := getAdminID(c)

		var req order.BatchTransferRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			respondBadRequest(c, err.Error())
			return
		}

		result, err := svc.BatchTransferSubOrders(c.Request.Context(), adminID, req)
		if err != nil {
			respondError(c, err)
			return
		}

		respondSuccess(c, result)
	}
}

// getAdminID 从上下文获取管理员ID
func getAdminID(c *gin.Context) uint64 {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	userID, ok := userIDVal.(uint64)
	if !ok {
		return 0
	}
	return userID
}

// OrderGroupListResponse 主订单列表响应
type OrderGroupListResponse struct {
	Items []model.OrderGroup `json:"items"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
}
