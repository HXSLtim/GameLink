package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service/order"
)

// RegisterPlayerOrderRoutes 注册陪玩师端订单管理路由
func RegisterPlayerOrderRoutes(router gin.IRouter, svc *order.OrderService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/player/orders")
	group.Use(authMiddleware) // 需要认证
	{
		group.GET("/available", func(c *gin.Context) { getAvailableOrdersHandler(c, svc) })
		group.POST("/:id/accept", func(c *gin.Context) { acceptOrderHandler(c, svc) })
		group.GET("/my", func(c *gin.Context) { getMyAcceptedOrdersHandler(c, svc) })
		group.PUT("/:id/complete", func(c *gin.Context) { completeOrderByPlayerHandler(c, svc) })
	}
}

// getAvailableOrdersHandler 获取可接订单列表
// @Summary      获取可接订单列表
// @Description  获取订单大厅的可接订单列表
// @Tags         Player - Orders
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        gameId         query     int     false  "游戏ID"
// @Param        page           query     int     false  "页码"
// @Param        pageSize       query     int     false  "每页数量"
// @Success      200            {object}  model.APIResponse[map[string]any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /player/orders/available [get]
func getAvailableOrdersHandler(c *gin.Context, svc *order.OrderService) {
	var req order.AvailableOrdersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	orders, total, err := svc.GetAvailableOrders(c.Request.Context(), req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[map[string]any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: map[string]any{
			"orders": orders,
			"total":  total,
		},
	})
}

// acceptOrderHandler 接单
// @Summary      接单
// @Description  陪玩师接单
// @Tags         Player - Orders
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      int     true  "订单ID"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /player/orders/{id}/accept [post]
func acceptOrderHandler(c *gin.Context, svc *order.OrderService) {
	userID := getUserIDFromContext(c)

	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, ErrInvalidID)
		return
	}

	if err := svc.AcceptOrder(c.Request.Context(), userID, orderID); err != nil {
		if err == order.ErrInvalidTransition {
			respondError(c, http.StatusBadRequest, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "接单成功",
	})
}

// getMyAcceptedOrdersHandler 获取我接的订单
// @Summary      获取我接的订单
// @Description  获取陪玩师接的订单列表
// @Tags         Player - Orders
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        status         query     string  false  "订单状态"
// @Param        page           query     int     false  "页码"
// @Param        pageSize       query     int     false  "每页数量"
// @Success      200            {object}  model.APIResponse[order.MyOrderListResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /player/orders/my [get]
func getMyAcceptedOrdersHandler(c *gin.Context, svc *order.OrderService) {
	userID := getUserIDFromContext(c)

	var req order.MyOrderListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := svc.GetMyOrders(c.Request.Context(), userID, req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[order.MyOrderListResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *resp,
	})
}

// completeOrderByPlayerHandler 完成订单（陪玩师端）
// @Summary      完成订单
// @Description  陪玩师完成订单
// @Tags         Player - Orders
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      int     true  "订单ID"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /player/orders/{id}/complete [put]
func completeOrderByPlayerHandler(c *gin.Context, svc *order.OrderService) {
	userID := getUserIDFromContext(c)

	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, ErrInvalidID)
		return
	}

	if err := svc.CompleteOrderByPlayer(c.Request.Context(), userID, orderID); err != nil {
		if err == order.ErrUnauthorized {
			respondError(c, http.StatusForbidden, err.Error())
			return
		}
		if err == order.ErrInvalidTransition {
			respondError(c, http.StatusBadRequest, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "订单已完成",
	})
}
