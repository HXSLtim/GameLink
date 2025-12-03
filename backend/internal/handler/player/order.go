package player

import (
	"github.com/gin-gonic/gin"

	"gamelink/pkg/apierr"
	"gamelink/internal/model"
	"gamelink/internal/service/order"
)

// SuccessResponseSwagger 通用成功响应（Swagger 用，避免泛型）
type SuccessResponseSwagger struct {
	Success    bool              `json:"success"`
	Code       int               `json:"code"`
	Message    string            `json:"message"`
	Data       interface{}       `json:"data,omitempty"`
	Pagination *model.Pagination `json:"pagination,omitempty"`
	TraceID    string            `json:"traceId,omitempty"`
}

// RegisterOrderRoutes 注册陪玩师端订单管理路由
func RegisterOrderRoutes(router gin.IRouter, svc *order.OrderService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/orders")
	group.Use(authMiddleware)
	group.GET("/available", func(c *gin.Context) { getAvailableOrdersHandler(c, svc) })
	group.POST(":id/accept", func(c *gin.Context) { acceptOrderHandler(c, svc) })
	group.GET("/my", func(c *gin.Context) { getMyAcceptedOrdersHandler(c, svc) })
	group.PUT(":id/complete", func(c *gin.Context) { completeOrderByPlayerHandler(c, svc) })
}

// getAvailableOrdersHandler 获取可接订单列表
// @Summary      获取可接订单列表
// @Description  获取订单大厅的可接订单列表
// @Tags         Player - Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        gameId   query     uint64  false  "游戏ID"
// @Param        page     query     int     false  "页码" default(1)
// @Param        pageSize query     int     false  "每页数量" default(20)
// @Success      200      {object}  SuccessResponseSwagger
// @Failure      400      {object}  apierr.APIError
// @Failure      401      {object}  apierr.APIError
// @Failure      500      {object}  apierr.APIError
// @Router       /player/orders/available [get]
func getAvailableOrdersHandler(c *gin.Context, svc *order.OrderService) {
	var req order.AvailableOrdersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("无效的查询参数").WithDetails(err.Error()))
		return
	}

	orders, total, err := svc.GetAvailableOrders(c.Request.Context(), req)
	if err != nil {
		respondAPIError(c, apierr.InternalError("获取可接订单列表失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "OK", map[string]interface{}{
		"orders": orders,
		"total":  total,
	})
}

// acceptOrderHandler 接单
// @Summary      接单
// @Description  陪玩师接单
// @Tags         Player - Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id  path  uint64  true  "订单ID"
// @Success      200  {object}  SuccessResponseSwagger
// @Failure      400  {object}  apierr.APIError
// @Failure      401  {object}  apierr.APIError
// @Failure      403  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Failure      409  {object}  apierr.APIError
// @Failure      500  {object}  apierr.APIError
// @Router       /player/orders/{id}/accept [post]
func acceptOrderHandler(c *gin.Context, svc *order.OrderService) {
	userID := getUserIDFromContext(c)

	orderID, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest("无效的ID"))
		return
	}

	if err := svc.AcceptOrder(c.Request.Context(), userID, orderID); err != nil {
		if err == order.ErrInvalidTransition {
			respondAPIError(c, apierr.BadRequest(err.Error()))
			return
		}
		if apierr.IsNotFound(err) {
			respondAPIError(c, apierr.NotFound(err.Error()))
			return
		}
		if err == order.ErrUnauthorized {
			respondAPIError(c, apierr.Forbidden("无权限操作此订单"))
			return
		}
		respondAPIError(c, apierr.InternalError("接单失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "接单成功", struct{}{})
}

// getMyAcceptedOrdersHandler 获取我的已接订单
// @Summary      获取我的已接订单
// @Description  获取陪玩师已接的订单列表
// @Tags         Player - Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        status    query     string  false  "状态过滤" Enums(pending,confirmed,in_progress,completed,canceled)
// @Param        page      query     int     false  "页码" default(1)
// @Param        pageSize  query     int     false  "每页数量" default(20)
// @Success      200       {object}  SuccessResponseSwagger
// @Failure      400       {object}  apierr.APIError
// @Failure      401       {object}  apierr.APIError
// @Failure      500       {object}  apierr.APIError
// @Router       /player/orders/my [get]
func getMyAcceptedOrdersHandler(c *gin.Context, svc *order.OrderService) {
	userID := getUserIDFromContext(c)

	var req order.MyOrderListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("无效的查询参数").WithDetails(err.Error()))
		return
	}

	resp, err := svc.GetMyOrders(c.Request.Context(), userID, req)
	if err != nil {
		respondAPIError(c, apierr.InternalError("获取订单列表失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "OK", *resp)
}

// completeOrderByPlayerHandler 完成订单（陪玩师端）
// @Summary      完成订单
// @Description  陪玩师完成订单
// @Tags         Player - Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id  path  uint64  true  "订单ID"
// @Success      200  {object}  SuccessResponseSwagger
// @Failure      400  {object}  apierr.APIError
// @Failure      401  {object}  apierr.APIError
// @Failure      403  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Failure      409  {object}  apierr.APIError
// @Failure      500  {object}  apierr.APIError
// @Router       /player/orders/{id}/complete [put]
func completeOrderByPlayerHandler(c *gin.Context, svc *order.OrderService) {
	userID := getUserIDFromContext(c)

	orderID, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest("无效的ID"))
		return
	}

	if err := svc.CompleteOrderByPlayer(c.Request.Context(), userID, orderID); err != nil {
		if err == order.ErrUnauthorized {
			respondAPIError(c, apierr.Forbidden("无权限操作此订单"))
			return
		}
		if err == order.ErrInvalidTransition {
			respondAPIError(c, apierr.BadRequest(err.Error()))
			return
		}
		respondAPIError(c, apierr.InternalError("完成订单失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "订单已完成", struct{}{})
}
