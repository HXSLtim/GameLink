package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// RegisterStatsAnalysisRoutes 注册管理端统计分析路�?func RegisterStatsAnalysisRoutes(
	router gin.IRouter,
	orderRepo repository.OrderRepository,
	commissionRepo repository.CommissionRepository,
	serviceItemRepo repository.ServiceItemRepository,
) {
	group := router.Group("/admin/stats")
	{
		group.GET("/service-items", func(c *gin.Context) {
			getServiceItemStatsHandler(c, orderRepo, serviceItemRepo)
		})
		group.GET("/top-players", func(c *gin.Context) {
			getTopPlayersHandler(c, commissionRepo)
		})
		group.GET("/gift-stats", func(c *gin.Context) {
			getAdminGiftStatsHandler(c, orderRepo, serviceItemRepo)
		})
		group.GET("/revenue-by-game", func(c *gin.Context) {
			getRevenueByGameHandler(c, orderRepo)
		})
	}
}

// getServiceItemStatsHandler 服务项目统计
// @Summary      服务项目统计
// @Description  统计各服务项目的销售情�?// @Tags         Admin - Stats
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/stats/service-items [get]
func getServiceItemStatsHandler(
	c *gin.Context,
	orderRepo repository.OrderRepository,
	serviceItemRepo repository.ServiceItemRepository,
) {
	ctx := c.Request.Context()

	// 获取所有服务项�?	items, _, err := serviceItemRepo.List(ctx, repository.ServiceItemListOptions{
		Page:     1,
		PageSize: 1000,
	})
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 统计每个服务项目的订单数和收�?	type ItemStats struct {
		ItemID      uint64 `json:"itemId"`
		ItemCode    string `json:"itemCode"`
		ItemName    string `json:"itemName"`
		SubCategory string `json:"subCategory"`
		OrderCount  int64  `json:"orderCount"`
		TotalRevenue int64 `json:"totalRevenue"`
	}

	stats := make([]ItemStats, 0, len(items))
	for _, item := range items {
		// 查询该服务项目的所有订�?		orders, _, _ := orderRepo.List(ctx, repository.OrderListOptions{
			// TODO: 需要添�?ItemID 过滤
			Statuses: []model.OrderStatus{model.OrderStatusCompleted},
			Page:     1,
			PageSize: 10000,
		})

		var orderCount int64
		var totalRevenue int64
		for _, order := range orders {
			if order.ItemID == item.ID {
				orderCount++
				totalRevenue += order.TotalPriceCents
			}
		}

		stats = append(stats, ItemStats{
			ItemID:       item.ID,
			ItemCode:     item.ItemCode,
			ItemName:     item.Name,
			SubCategory:  string(item.SubCategory),
			OrderCount:   orderCount,
			TotalRevenue: totalRevenue,
		})
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: map[string]interface{}{
			"items": stats,
		},
	})
}

// getTopPlayersHandler 获取Top陪玩�?// @Summary      获取Top陪玩�?// @Description  按收入排名获取Top陪玩�?// @Tags         Admin - Stats
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        month          query     string  false  "月份(YYYY-MM)"
// @Param        limit          query     int     false  "数量"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/stats/top-players [get]
func getTopPlayersHandler(c *gin.Context, commissionRepo repository.CommissionRepository) {
	ctx := c.Request.Context()
	month := c.DefaultQuery("month", time.Now().Format("2006-01"))
	limit := 10

	// TODO: 实现Top陪玩师查�?	// 暂时返回空数�?	_ = ctx
	_ = month

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: map[string]interface{}{
			"players": []interface{}{},
			"month":   month,
			"limit":   limit,
		},
	})
}

// getAdminGiftStatsHandler 获取礼物统计
// @Summary      获取礼物统计
// @Description  统计礼物的销售情�?// @Tags         Admin - Stats
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/stats/gift-stats [get]
func getAdminGiftStatsHandler(
	c *gin.Context,
	orderRepo repository.OrderRepository,
	serviceItemRepo repository.ServiceItemRepository,
) {
	ctx := c.Request.Context()

	// 获取所有礼�?	gifts, _, _ := serviceItemRepo.GetGifts(ctx, 1, 1000)

	type GiftStat struct {
		GiftID      uint64 `json:"giftId"`
		GiftName    string `json:"giftName"`
		TotalSent   int64  `json:"totalSent"`
		TotalRevenue int64 `json:"totalRevenue"`
	}

	giftStats := make([]GiftStat, 0, len(gifts))
	for _, gift := range gifts {
		// 统计该礼物的销�?		orders, _, _ := orderRepo.List(ctx, repository.OrderListOptions{
			Statuses: []model.OrderStatus{model.OrderStatusCompleted},
			Page:     1,
			PageSize: 10000,
		})

		var totalSent int64
		var totalRevenue int64
		for _, order := range orders {
			if order.ItemID == gift.ID && order.IsGiftOrder() {
				totalSent += int64(order.Quantity)
				totalRevenue += order.TotalPriceCents
			}
		}

		giftStats = append(giftStats, GiftStat{
			GiftID:       gift.ID,
			GiftName:     gift.Name,
			TotalSent:    totalSent,
			TotalRevenue: totalRevenue,
		})
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: map[string]interface{}{
			"gifts": giftStats,
		},
	})
}

// getRevenueByGameHandler 按游戏统计收�?// @Summary      按游戏统计收�?// @Description  统计各游戏的订单和收入情�?// @Tags         Admin - Stats
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/stats/revenue-by-game [get]
func getRevenueByGameHandler(c *gin.Context, orderRepo repository.OrderRepository) {
	ctx := c.Request.Context()

	// 获取所有已完成订单
	orders, _, err := orderRepo.List(ctx, repository.OrderListOptions{
		Statuses: []model.OrderStatus{model.OrderStatusCompleted},
		Page:     1,
		PageSize: 10000,
	})
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 按游戏ID分组统计
	gameStats := make(map[uint64]struct {
		OrderCount int64
		Revenue    int64
	})

	for _, order := range orders {
		gameID := order.GetGameID()
		if gameID > 0 {
			stats := gameStats[gameID]
			stats.OrderCount++
			stats.Revenue += order.TotalPriceCents
			gameStats[gameID] = stats
		}
	}

	// 转换为数�?	type GameRevenue struct {
		GameID      uint64 `json:"gameId"`
		OrderCount  int64  `json:"orderCount"`
		TotalRevenue int64 `json:"totalRevenue"`
	}

	result := make([]GameRevenue, 0, len(gameStats))
	for gameID, stats := range gameStats {
		result = append(result, GameRevenue{
			GameID:       gameID,
			OrderCount:   stats.OrderCount,
			TotalRevenue: stats.Revenue,
		})
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: map[string]interface{}{
			"games": result,
		},
	})
}

