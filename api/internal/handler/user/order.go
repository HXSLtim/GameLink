package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service/order"
	"gamelink/pkg/apierr"
)

// CreateOrderResponse 创建订单响应（类型别名）
type CreateOrderResponse = order.CreateOrderResponse

// MyOrderListResponse 我的订单列表响应（类型别名）
type MyOrderListResponse = order.MyOrderListResponse

// OrderDetailResponse 订单详情响应（类型别名）
type OrderDetailResponse = order.OrderDetailResponse

// RegisterOrderRoutes 注册用户端订单路
func RegisterOrderRoutes(router gin.IRouter, svc *order.OrderService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/orders")
	group.Use(authMiddleware) // 需要认证
	group.POST("", func(c *gin.Context) { createOrderHandler(c, svc) })
	group.GET("", func(c *gin.Context) { getMyOrdersHandler(c, svc) })
	group.GET("/:id", func(c *gin.Context) { getOrderDetailHandler(c, svc) })
	group.PUT("/:id/cancel", func(c *gin.Context) { cancelOrderHandler(c, svc) })
	group.PUT("/:id/complete", func(c *gin.Context) { completeOrderHandler(c, svc) })
}

// createOrderHandler 创建订单
// @Summary      创建订单
// @Description  用户创建陪玩订单
// @Tags         User - Orders
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                       true  "Bearer {token}"
// @Param        request        body      order.CreateOrderRequest     true  "创建订单请求"
// @Success      200            {object}  CreateOrderResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /user/orders [post]
func createOrderHandler(c *gin.Context, svc *order.OrderService) {
	userID := getUserIDFromContext(c)

	var req order.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}

	resp, err := svc.CreateOrder(c.Request.Context(), userID, req)
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err.(*apierr.APIError))
			return
		}
		// Check for repository.ErrNotFound (player or game not found)
		if errors.Is(err, repository.ErrNotFound) || apierr.IsNotFound(err) {
			respondAPIError(c, apierr.NotFound("陪玩师或游戏不存在"))
			return
		}
		respondAPIError(c, apierr.InternalError("创建订单失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "订单创建成功", *resp)
}

func getOrderMessagesHandler(c *gin.Context, svc *order.OrderService) {
	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
	})
}

// getMyOrdersHandler 获取我的订单列表
// @Summary      获取我的订单列表
// @Description  获取当前用户的订单列表，支持状态过滤和分页
// @Tags         User - Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        status    query     string     false  "Status filter" Enums(pending,confirmed,in_progress,completed,canceled,refunded)
// @Param        page      query     int        false  "Page number" default(1)
// @Param        pageSize  query     int        false  "Page size" default(20)
// @Success      200       {object}  MyOrderListResponse
// @Failure      400       {object}  model.ErrorResponse
// @Failure      401       {object}  model.ErrorResponse
// @Router       /user/orders [get]
func getMyOrdersHandler(c *gin.Context, svc *order.OrderService) {
	userID := getUserIDFromContext(c)

	var req order.MyOrderListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的查询参 "+err.Error())
		return
	}

	resp, err := svc.GetMyOrders(c.Request.Context(), userID, req)
	if err != nil {
		respondAPIError(c, apierr.InternalError("获取订单列表失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "OK", *resp)
}

// getOrderDetailHandler 获取订单详情
// @Summary      获取订单详情
// @Description  获取指定订单的详细信
// @Tags         User - Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      uint64  true  "订单ID"
// @Success      200   {object}  OrderDetailResponse
// @Failure      400   {object}  model.ErrorResponse
// @Failure      401   {object}  model.ErrorResponse
// @Failure      404   {object}  model.ErrorResponse
// @Router       /user/orders/{id} [get]
func getOrderDetailHandler(c *gin.Context, svc *order.OrderService) {
	userID := getUserIDFromContext(c)

	orderID, err := parseUintParam(c, "id")
	if err != nil {
		respondError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}

	resp, err := svc.GetOrderDetail(c.Request.Context(), userID, orderID)
	if err != nil {
		respondAPIError(c, err)
		return
	}

	respondSuccess(c, "OK", *resp)
}

// cancelOrderHandler 取消订单
// @Summary      取消订单
// @Description  用户取消指定订单
// @Tags         User - Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      uint64                     true  "订单ID"
// @Param        request  body      order.CancelOrderRequest   true  "取消原因"
// @Success      200      {object}  model.SuccessResponse
// @Failure      400      {object}  apierr.APIError
// @Failure      401      {object}  apierr.APIError
// @Failure      403      {object}  apierr.APIError
// @Failure      404      {object}  apierr.APIError
// @Failure      500      {object}  apierr.APIError
// @Router       /user/orders/{id}/cancel [put]
func cancelOrderHandler(c *gin.Context, svc *order.OrderService) {
	userID := getUserIDFromContext(c)

	orderID, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidID))
		return
	}

	var req order.CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return
	}

	if err := svc.CancelOrder(c.Request.Context(), userID, orderID, req); err != nil {
		// Check for not found error first
		if errors.Is(err, repository.ErrNotFound) || errors.Is(err, order.ErrNotFound) {
			respondAPIError(c, apierr.NotFound("订单不存在"))
			return
		}
		if err == order.ErrUnauthorized {
			respondAPIError(c, apierr.Forbidden("无权限取消此订单"))
			return
		}
		if err == order.ErrInvalidTransition {
			respondAPIError(c, apierr.BadRequest(err.Error()))
			return
		}
		respondAPIError(c, apierr.InternalError("取消订单失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "订单已取消", struct{}{})
}

// completeOrderHandler 完成订单
// @Summary      完成订单
// @Description  用户确认订单完成
// @Tags         User - Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      uint64  true  "订单ID"
// @Success      200   {object}  model.SuccessResponse
// @Failure      400   {object}  apierr.APIError
// @Failure      401   {object}  apierr.APIError
// @Failure      403   {object}  apierr.APIError
// @Failure      404   {object}  apierr.APIError
// @Failure      500   {object}  apierr.APIError
// @Router       /user/orders/{id}/complete [put]
func completeOrderHandler(c *gin.Context, svc *order.OrderService) {
	userID := getUserIDFromContext(c)

	orderID, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidID))
		return
	}

	if err := svc.CompleteOrder(c.Request.Context(), userID, orderID); err != nil {
		// Check for not found error first
		if errors.Is(err, repository.ErrNotFound) || errors.Is(err, order.ErrNotFound) {
			respondAPIError(c, apierr.NotFound("订单不存在"))
			return
		}
		if err == order.ErrUnauthorized {
			respondAPIError(c, apierr.Forbidden("无权限完成此订单"))
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

// getUserIDFromContext 从上下文获取用户ID
func getUserIDFromContext(c *gin.Context) uint64 {
	// JWT 中间件设置的上下文中获取用户ID
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
