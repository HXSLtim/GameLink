package user

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/apierr"
	"gamelink/internal/model"
	assignmentservice "gamelink/internal/service/assignment"
	"gamelink/internal/service/order"
)

// RegisterOrderRoutes 注册用户端订单路由
func RegisterOrderRoutes(router gin.IRouter, svc *order.OrderService, assignSvc *assignmentservice.Service, authMiddleware gin.HandlerFunc) {
	group := router.Group("/user/orders")
	group.Use(authMiddleware) // 需要认证
	group.POST("", func(c *gin.Context) { createOrderHandler(c, svc) })
	group.GET("", func(c *gin.Context) { getMyOrdersHandler(c, svc) })
	group.GET("/:id", func(c *gin.Context) { getOrderDetailHandler(c, svc) })
	group.PUT("/:id/cancel", func(c *gin.Context) { cancelOrderHandler(c, svc) })
	group.PUT("/:id/complete", func(c *gin.Context) { completeOrderHandler(c, svc) })
	group.POST("/:id/dispute", func(c *gin.Context) { createDisputeHandler(c, assignSvc) })
}

// createOrderHandler 创建订单
// @Summary      创建订单
// @Description  用户创建陪玩订单
// @Tags         User - Orders
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                       true  "Bearer {token}"
// @Param        request        body      order.CreateOrderRequest     true  "创建订单请求"
// @Success      200            {object}  model.APIResponse[order.CreateOrderResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /user/orders [post]
func createOrderHandler(c *gin.Context, svc *order.OrderService) {
	userID := getUserIDFromContext(c)

	var req order.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := svc.CreateOrder(c.Request.Context(), userID, req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[order.CreateOrderResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "订单创建成功",
		Data:    *resp,
	})
}

// getMyOrdersHandler 获取我的订单列表
// @Summary      获取我的订单列表
// @Description  获取当前用户的订单列�?// @Tags         User - Orders
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        status         query     string  false  "订单状�?
// @Param        page           query     int     false  "页码"
// @Param        pageSize       query     int     false  "每页数量"
// @Success      200            {object}  model.APIResponse[order.MyOrderListResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /user/orders [get]
func getMyOrdersHandler(c *gin.Context, svc *order.OrderService) {
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

// getOrderDetailHandler 获取订单详情
// @Summary      获取订单详情
// @Description  获取订单详细信息
// @Tags         User - Orders
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      int     true  "订单ID"
// @Success      200            {object}  model.APIResponse[order.OrderDetailResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Failure      404            {object}  model.APIResponse[any]
// @Router       /user/orders/{id} [get]
func getOrderDetailHandler(c *gin.Context, svc *order.OrderService) {
	userID := getUserIDFromContext(c)

	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}

	resp, err := svc.GetOrderDetail(c.Request.Context(), userID, orderID)
	if err != nil {
		if err == order.ErrNotFound {
			respondError(c, http.StatusNotFound, err.Error())
			return
		}
		if err == order.ErrUnauthorized {
			respondError(c, http.StatusForbidden, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[order.OrderDetailResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *resp,
	})
}

// cancelOrderHandler 取消订单
// @Summary      取消订单
// @Description  用户取消订单
// @Tags         User - Orders
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                     true  "Bearer {token}"
// @Param        id             path      int                        true  "订单ID"
// @Param        request        body      order.CancelOrderRequest   true  "取消原因"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /user/orders/{id}/cancel [put]
func cancelOrderHandler(c *gin.Context, svc *order.OrderService) {
	userID := getUserIDFromContext(c)

	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}

	var req order.CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := svc.CancelOrder(c.Request.Context(), userID, orderID, req); err != nil {
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
		Message: "订单已取消",
	})
}

// completeOrderHandler 完成订单
// @Summary      完成订单
// @Description  用户确认订单完成
// @Tags         User - Orders
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      int     true  "订单ID"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /user/orders/{id}/complete [put]
func completeOrderHandler(c *gin.Context, svc *order.OrderService) {
	userID := getUserIDFromContext(c)

	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}

	if err := svc.CompleteOrder(c.Request.Context(), userID, orderID); err != nil {
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

type disputeRequest struct {
	Reason   string   `json:"reason" binding:"required"`
	Evidence []string `json:"evidence"`
}

// createDisputeHandler 用户发起争议
// @Summary      发起订单争议
// @Tags         User - Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int             true  "订单ID"
// @Param        request  body  disputeRequest  true  "争议信息"
// @Success      201 {object} model.APIResponse[*model.OrderDispute]
// @Failure      400 {object} model.APIResponse[any]
// @Router       /user/orders/{id}/dispute [post]
func createDisputeHandler(c *gin.Context, svc *assignmentservice.Service) {
	if svc == nil {
		respondError(c, http.StatusServiceUnavailable, "dispute service unavailable")
		return
	}
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	var req disputeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}
	traceID := assignmentservice.TraceIDFromContext(c.Request.Context())
	actor := userID
	dispute, err := svc.CreateDispute(c.Request.Context(), orderID, assignmentservice.DisputeRequest{
		RaisedBy:       model.OrderDisputeRaisedByUser,
		RaisedByUserID: &actor,
		Reason:         strings.TrimSpace(req.Reason),
		EvidenceURLs:   req.Evidence,
		TraceID:        traceID,
	})
	if err != nil {
		switch {
		case errors.Is(err, assignmentservice.ErrValidation):
			respondError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, assignmentservice.ErrNotFound):
			respondError(c, http.StatusNotFound, err.Error())
		default:
			respondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondJSON(c, http.StatusCreated, model.APIResponse[*model.OrderDispute]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "争议已创建",
		Data:    dispute,
		TraceID: traceID,
	})
}

// getUserIDFromContext 从上下文获取用户ID
func getUserIDFromContext(c *gin.Context) uint64 {
	// 从 JWT 中间件设置的上下文中获取用户ID
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
