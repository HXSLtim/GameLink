package admin

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// DashboardService Dashboard统计服务接口
type DashboardService interface {
	GetOverviewStats(ctx context.Context) (*DashboardOverviewStats, error)
	GetRecentOrders(ctx context.Context, limit int) ([]model.Order, error)
	GetRecentWithdraws(ctx context.Context, limit int) ([]model.Withdraw, error)
	GetMonthlyRevenue(ctx context.Context, months int) ([]MonthlyRevenueData, error)
}

// DashboardOverviewStats 总览统计
type DashboardOverviewStats struct {
	TotalUsers       int64   `json:"totalUsers"`
	TotalPlayers     int64   `json:"totalPlayers"`
	TotalOrders      int64   `json:"totalOrders"`
	TodayOrders      int64   `json:"todayOrders"`
	TodayRevenue     int64   `json:"todayRevenue"`
	MonthRevenue     int64   `json:"monthRevenue"`
	PendingWithdraws int64   `json:"pendingWithdraws"`
	ActiveServices   int64   `json:"activeServices"`
}

// MonthlyRevenueData 月度收入数据
type MonthlyRevenueData struct {
	Month            string  `json:"month"`
	TotalRevenue     int64   `json:"totalRevenue"`
	TotalCommission  int64   `json:"totalCommission"`
	TotalOrders      int64   `json:"totalOrders"`
}

// RegisterDashboardRoutes 注册管理端Dashboard路由
func RegisterDashboardRoutes(
	router gin.IRouter,
	userRepo repository.UserRepository,
	playerRepo repository.PlayerRepository,
	orderRepo repository.OrderRepository,
	withdrawRepo repository.WithdrawRepository,
	serviceItemRepo repository.ServiceItemRepository,
	commissionRepo repository.CommissionRepository,
) {
	group := router.Group("/admin/dashboard")
	{
		group.GET("/overview", func(c *gin.Context) {
			getDashboardOverviewHandler(c, userRepo, playerRepo, orderRepo, withdrawRepo, serviceItemRepo)
		})
		group.GET("/recent-orders", func(c *gin.Context) {
			getRecentOrdersHandler(c, orderRepo)
		})
		group.GET("/recent-withdraws", func(c *gin.Context) {
			getRecentWithdrawsHandler(c, withdrawRepo)
		})
		group.GET("/monthly-revenue", func(c *gin.Context) {
			getMonthlyRevenueHandler(c, commissionRepo)
		})
	}
}

// getDashboardOverviewHandler 获取Dashboard总览
// @Summary      获取Dashboard总览
// @Description  管理员Dashboard总览数据
// @Tags         Admin - Dashboard
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.APIResponse[DashboardOverviewStats]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/dashboard/overview [get]
func getDashboardOverviewHandler(
	c *gin.Context,
	userRepo repository.UserRepository,
	playerRepo repository.PlayerRepository,
	orderRepo repository.OrderRepository,
	withdrawRepo repository.WithdrawRepository,
	serviceItemRepo repository.ServiceItemRepository,
) {
	ctx := c.Request.Context()
	stats := &DashboardOverviewStats{}

	// 总用户数
	users, _ := userRepo.List(ctx)
	stats.TotalUsers = int64(len(users))

	// 总陪玩师�?	players, total, _ := playerRepo.ListPaged(ctx, 1, 1)
	_ = players
	stats.TotalPlayers = total

	// 总订单数
	orders, total, _ := orderRepo.List(ctx, repository.OrderListOptions{Page: 1, PageSize: 1})
	_ = orders
	stats.TotalOrders = total

	// 今日订单�?	todayStart := time.Now().Truncate(24 * time.Hour)
	todayOrders, todayTotal, _ := orderRepo.List(ctx, repository.OrderListOptions{
		DateFrom: &todayStart,
		Page:     1,
		PageSize: 10000,
	})
	stats.TodayOrders = todayTotal

	// 今日收入
	var todayRevenue int64
	for _, order := range todayOrders {
		if order.Status == model.OrderStatusCompleted {
			todayRevenue += order.TotalPriceCents
		}
	}
	stats.TodayRevenue = todayRevenue

	// 本月收入
	monthStart := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	monthOrders, _, _ := orderRepo.List(ctx, repository.OrderListOptions{
		DateFrom: &monthStart,
		Statuses: []model.OrderStatus{model.OrderStatusCompleted},
		Page:     1,
		PageSize: 10000,
	})
	var monthRevenue int64
	for _, order := range monthOrders {
		monthRevenue += order.TotalPriceCents
	}
	stats.MonthRevenue = monthRevenue

	// 待审批提�?	pendingStatus := model.WithdrawStatusPending
	_, pendingTotal, _ := withdrawRepo.List(ctx, repository.WithdrawListOptions{
		Status:   &pendingStatus,
		Page:     1,
		PageSize: 1,
	})
	stats.PendingWithdraws = pendingTotal

	// 活跃服务�?	isActive := true
	_, activeTotal, _ := serviceItemRepo.List(ctx, repository.ServiceItemListOptions{
		IsActive: &isActive,
		Page:     1,
		PageSize: 1,
	})
	stats.ActiveServices = activeTotal

	respondJSON(c, http.StatusOK, model.APIResponse[DashboardOverviewStats]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *stats,
	})
}

// getRecentOrdersHandler 获取最近订�?// @Summary      获取最近订�?// @Description  管理员查看最近订�?// @Tags         Admin - Dashboard
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        limit          query     int     false  "数量限制"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/dashboard/recent-orders [get]
func getRecentOrdersHandler(c *gin.Context, orderRepo repository.OrderRepository) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	orders, _, err := orderRepo.List(c.Request.Context(), repository.OrderListOptions{
		Page:     1,
		PageSize: limit,
	})
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: map[string]interface{}{
			"orders": orders,
		},
	})
}

// getRecentWithdrawsHandler 获取最近提�?// @Summary      获取最近提�?// @Description  管理员查看最近提现申�?// @Tags         Admin - Dashboard
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        limit          query     int     false  "数量限制"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/dashboard/recent-withdraws [get]
func getRecentWithdrawsHandler(c *gin.Context, withdrawRepo repository.WithdrawRepository) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	withdraws, _, err := withdrawRepo.List(c.Request.Context(), repository.WithdrawListOptions{
		Page:     1,
		PageSize: limit,
	})
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: map[string]interface{}{
			"withdraws": withdraws,
		},
	})
}

// getMonthlyRevenueHandler 获取月度收入趋势
// @Summary      获取月度收入趋势
// @Description  管理员查看最近几个月的收入趋�?// @Tags         Admin - Dashboard
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        months         query     int     false  "月数（默�?2�?
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/dashboard/monthly-revenue [get]
func getMonthlyRevenueHandler(c *gin.Context, commissionRepo repository.CommissionRepository) {
	months := 12
	if monthsStr := c.Query("months"); monthsStr != "" {
		if m, err := strconv.Atoi(monthsStr); err == nil && m > 0 && m <= 24 {
			months = m
		}
	}

	ctx := c.Request.Context()
	revenueData := make([]MonthlyRevenueData, 0, months)

	// 计算每个月的数据
	now := time.Now()
	for i := months - 1; i >= 0; i-- {
		month := now.AddDate(0, -i, 0).Format("2006-01")

		// 获取月度统计
		stats, err := commissionRepo.GetMonthlyStats(ctx, month)
		if err == nil {
			revenueData = append(revenueData, MonthlyRevenueData{
				Month:           month,
				TotalRevenue:    stats.TotalIncome,
				TotalCommission: stats.TotalCommission,
				TotalOrders:     stats.TotalOrders,
			})
		} else {
			// 如果没有数据，填�?
			revenueData = append(revenueData, MonthlyRevenueData{
				Month:           month,
				TotalRevenue:    0,
				TotalCommission: 0,
				TotalOrders:     0,
			})
		}
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: map[string]interface{}{
			"revenue": revenueData,
		},
	})
}

