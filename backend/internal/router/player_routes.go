package router

import (
	"github.com/gin-gonic/gin"

	playerhandler "gamelink/internal/handler/player"
)

// registerPlayerRoutes 注册陪玩端路由
func registerPlayerRoutes(api *gin.RouterGroup, authMiddleware gin.HandlerFunc, services *appServices) {
	playerGroup := api.Group("/player")
	playerGroup.Use(authMiddleware)
	{
		playerhandler.RegisterProfileRoutes(playerGroup, services.playerSvc, authMiddleware)
		playerhandler.RegisterOrderRoutes(playerGroup, services.orderSvc, authMiddleware)
		playerhandler.RegisterEarningsRoutes(playerGroup, services.earningsSvc, authMiddleware)
		playerhandler.RegisterCommissionRoutes(playerGroup, services.commissionSvc, authMiddleware)
		playerhandler.RegisterGiftRoutes(playerGroup, services.giftSvc, authMiddleware)
		playerhandler.RegisterReviewRoutes(playerGroup, services.reviewSvc, authMiddleware)
	}
}