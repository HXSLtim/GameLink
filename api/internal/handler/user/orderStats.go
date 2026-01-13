package user

import (
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/model"
	repoiface "gamelink/internal/repository/interfaces"
	"gamelink/pkg/apierr"
)

// OrderStatsHandler 订单统计处理器
type OrderStatsHandler struct {
	orderRepo repoiface.OrderRepository
}

// NewOrderStatsHandler 创建订单统计处理器
func NewOrderStatsHandler(orderRepo repoiface.OrderRepository) *OrderStatsHandler {
	return &OrderStatsHandler{orderRepo: orderRepo}
}

// OrderStatsResponse 订单统计响应
type OrderStatsResponse struct {
	TotalCount          int   `json:"totalCount"`
	MonthlyCount        int   `json:"monthlyCount"`
	MonthlyChange       int   `json:"monthlyChange"`
	PendingCount        int   `json:"pendingCount"`
	InProgressCount     int   `json:"inProgressCount"`
	CompletedCount      int   `json:"completedCount"`
	CanceledCount       int   `json:"canceledCount"`
	TotalSpentCents     int64 `json:"totalSpentCents"`
	AvgOrderAmountCents int64 `json:"avgOrderAmountCents"`
}

// GetOrderStats 获取订单统计
// @Summary      获取订单统计
// @Description  获取当前用户的订单统计数据
// @Tags         User - Orders
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  OrderStatsResponse
// @Failure      401  {object}  apierr.APIError
// @Failure      500  {object}  apierr.APIError
// @Router       /user/orders/stats [get]
func (h *OrderStatsHandler) GetOrderStats(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	ctx := c.Request.Context()

	// 获取各状态订单数量
	stats := OrderStatsResponse{}

	// 查询所有订单
	allOrders, total, err := h.orderRepo.List(ctx, repoiface.OrderListOptions{
		UserID:   &userID,
		Page:     1,
		PageSize: 10000, // 获取所有订单用于统计
	})
	if err != nil {
		resp.Error(c, apierr.InternalError("获取订单统计失败").WithDetails(err.Error()))
		return
	}

	stats.TotalCount = int(total)

	// 统计各状态数量和总消费
	now := time.Now()
	thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastMonth := thisMonth.AddDate(0, -1, 0)

	var monthlyCount, lastMonthCount int

	for _, order := range allOrders {
		switch order.Status {
		case model.OrderStatusPending:
			stats.PendingCount++
		case model.OrderStatusInProgress:
			stats.InProgressCount++
		case model.OrderStatusCompleted:
			stats.CompletedCount++
			stats.TotalSpentCents += order.TotalPriceCents
		case model.OrderStatusCanceled:
			stats.CanceledCount++
		}

		// 本月订单
		if order.CreatedAt.After(thisMonth) || order.CreatedAt.Equal(thisMonth) {
			monthlyCount++
		}
		// 上月订单
		if order.CreatedAt.After(lastMonth) && order.CreatedAt.Before(thisMonth) {
			lastMonthCount++
		}
	}

	stats.MonthlyCount = monthlyCount
	stats.MonthlyChange = monthlyCount - lastMonthCount

	// 计算平均订单金额
	if stats.CompletedCount > 0 {
		stats.AvgOrderAmountCents = stats.TotalSpentCents / int64(stats.CompletedCount)
	}

	resp.OK(c, stats)
}

// RegisterOrderStatsRoutes 注册订单统计路由
func RegisterOrderStatsRoutes(rg *gin.RouterGroup, orderRepo repoiface.OrderRepository, _ gin.HandlerFunc) {
	h := NewOrderStatsHandler(orderRepo)
	rg.GET("/orders/stats", h.GetOrderStats)
}
