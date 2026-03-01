package admin

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionrepo "gamelink/internal/repository/commission"
	repoiface "gamelink/internal/repository/interfaces"
	withdrawrepo "gamelink/internal/repository/withdraw"
	"gamelink/pkg/apierr"
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
	TotalUsers       int64 `json:"totalUsers"`
	TotalPlayers     int64 `json:"totalPlayers"`
	TotalOrders      int64 `json:"totalOrders"`
	TodayOrders      int64 `json:"todayOrders"`
	TodayRevenue     int64 `json:"todayRevenue"`
	MonthRevenue     int64 `json:"monthRevenue"`
	PendingWithdraws int64 `json:"pendingWithdraws"`
	ActiveServices   int64 `json:"activeServices"`
}

// MonthlyRevenueData 月度收入数据
type MonthlyRevenueData struct {
	Month           string `json:"month"`
	TotalRevenue    int64  `json:"totalRevenue"`
	TotalCommission int64  `json:"totalCommission"`
	TotalOrders     int64  `json:"totalOrders"`
}

// RegisterDashboardRoutes 注册管理端Dashboard路由
func RegisterDashboardRoutes(
	router gin.IRouter,
	userRepo repository.UserRepository,
	playerRepo repository.PlayerRepository,
	orderRepo repoiface.OrderQuery,
	withdrawRepo withdrawrepo.WithdrawRepository,
	serviceItemRepo repository.ServiceItemRepository,
	commissionRepo commissionrepo.CommissionRepository,
) {
	group := router.Group("/dashboard")
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
// @Success      200            {object}  model.DashboardOverviewStats
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/dashboard/overview [get]
func getDashboardOverviewHandler(
	c *gin.Context,
	userRepo repository.UserRepository,
	playerRepo repository.PlayerRepository,
	orderRepo repoiface.OrderQuery,
	withdrawRepo withdrawrepo.WithdrawRepository,
	serviceItemRepo repository.ServiceItemRepository,
) {
	ctx := c.Request.Context()
	stats := &DashboardOverviewStats{}

	users, err := userRepo.List(ctx)
	if err != nil {
		respondError(c, apierr.InternalError("获取用户统计失败").WithDetails(err.Error()))
		return
	}
	stats.TotalUsers = int64(len(users))

	_, totalPlayers, err := playerRepo.ListPaged(ctx, 1, 1)
	if err != nil {
		respondError(c, apierr.InternalError("获取陪玩统计失败").WithDetails(err.Error()))
		return
	}
	stats.TotalPlayers = totalPlayers

	orders, total, err := orderRepo.List(ctx, repoiface.OrderListOptions{Page: 1, PageSize: 1})
	if err != nil {
		respondError(c, apierr.InternalError("获取订单统计失败").WithDetails(err.Error()))
		return
	}
	_ = orders
	stats.TotalOrders = total

	todayStart := time.Now().Truncate(24 * time.Hour)
	todayOrders, todayTotal, err := orderRepo.List(ctx, repoiface.OrderListOptions{
		DateFrom: &todayStart,
		Page:     1,
		PageSize: 10000,
	})
	if err != nil {
		respondError(c, apierr.InternalError("获取今日订单统计失败").WithDetails(err.Error()))
		return
	}
	stats.TodayOrders = todayTotal

	var todayRevenue int64
	for _, order := range todayOrders {
		if order.Status == model.OrderStatusCompleted {
			todayRevenue += order.TotalPriceCents
		}
	}
	stats.TodayRevenue = todayRevenue

	monthStart := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	monthOrders, _, err := orderRepo.List(ctx, repoiface.OrderListOptions{
		DateFrom: &monthStart,
		Statuses: []model.OrderStatus{model.OrderStatusCompleted},
		Page:     1,
		PageSize: 10000,
	})
	if err != nil {
		respondError(c, apierr.InternalError("获取本月营收失败").WithDetails(err.Error()))
		return
	}
	var monthRevenue int64
	for _, order := range monthOrders {
		monthRevenue += order.TotalPriceCents
	}
	stats.MonthRevenue = monthRevenue

	pendingStatus := model.WithdrawStatusPending
	_, pendingTotal, err := withdrawRepo.List(ctx, withdrawrepo.WithdrawListOptions{
		Status:   &pendingStatus,
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		respondError(c, apierr.InternalError("获取提现统计失败").WithDetails(err.Error()))
		return
	}
	stats.PendingWithdraws = pendingTotal

	isActive := true
	_, activeTotal, err := serviceItemRepo.List(ctx, repository.ServiceItemListOptions{
		IsActive: &isActive,
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		respondError(c, apierr.InternalError("获取服务统计失败").WithDetails(err.Error()))
		return
	}
	stats.ActiveServices = activeTotal

	respondSuccess(c, *stats)
}

// getRecentOrdersHandler 获取最近订单
// @Summary      获取最近订单
// @Description  获取最近的订单列表
// @Tags         Admin - Dashboard
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        limit  query     int  false  "数量限制" default(10)
// @Success      200    {object}  model.SuccessResponse
// @Failure      400    {object}  model.ErrorResponse
// @Failure      401    {object}  model.ErrorResponse
// @Failure      500    {object}  model.ErrorResponse
// @Router       /admin/dashboard/recent-orders [get]
func getRecentOrdersHandler(c *gin.Context, orderRepo repoiface.OrderQuery) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	orders, _, err := orderRepo.List(c.Request.Context(), repoiface.OrderListOptions{
		Page:     1,
		PageSize: limit,
	})
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondSuccess(c, gin.H{"orders": orders})
}

// getRecentWithdrawsHandler 获取最近提现
// @Summary      获取最近提现
// @Description  获取最近的提现申请列表
// @Tags         Admin - Dashboard
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        limit  query     int  false  "数量限制" default(10)
// @Success      200    {object}  model.SuccessResponse
// @Failure      400    {object}  model.ErrorResponse
// @Failure      401    {object}  model.ErrorResponse
// @Failure      500    {object}  model.ErrorResponse
// @Router       /admin/dashboard/recent-withdraws [get]
func getRecentWithdrawsHandler(c *gin.Context, withdrawRepo withdrawrepo.WithdrawRepository) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	withdraws, _, err := withdrawRepo.List(c.Request.Context(), withdrawrepo.WithdrawListOptions{
		Page:     1,
		PageSize: limit,
	})
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondSuccess(c, gin.H{"withdraws": withdraws})
}

// getMonthlyRevenueHandler 获取月度收入趋势
// @Summary      获取月度收入趋势
// @Description  获取指定月数的月度收入趋势数据
// @Tags         Admin - Dashboard
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        months  query     int  false  "月数范围" default(12)
// @Success      200     {object}  model.SuccessResponse
// @Failure      400     {object}  model.ErrorResponse
// @Failure      401     {object}  model.ErrorResponse
// @Failure      500     {object}  model.ErrorResponse
// @Router       /admin/dashboard/monthly-revenue [get]
func getMonthlyRevenueHandler(c *gin.Context, commissionRepo commissionrepo.CommissionRepository) {
	months := 12
	if monthsStr := c.Query("months"); monthsStr != "" {
		if m, err := strconv.Atoi(monthsStr); err == nil && m > 0 && m <= 24 {
			months = m
		}
	}

	ctx := c.Request.Context()
	revenueData := make([]MonthlyRevenueData, 0, months)

	now := time.Now()
	for i := months - 1; i >= 0; i-- {
		month := now.AddDate(0, -i, 0).Format("2006-01")

		stats, err := commissionRepo.GetMonthlyStats(ctx, month)
		if err == nil && stats != nil {
			revenueData = append(revenueData, MonthlyRevenueData{
				Month:           month,
				TotalRevenue:    stats.TotalIncome,
				TotalCommission: stats.TotalCommission,
				TotalOrders:     stats.TotalOrders,
			})
		} else {
			revenueData = append(revenueData, MonthlyRevenueData{
				Month:           month,
				TotalRevenue:    0,
				TotalCommission: 0,
				TotalOrders:     0,
			})
		}
	}

	respondSuccess(c, gin.H{"revenue": revenueData})
}
