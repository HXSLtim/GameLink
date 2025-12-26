package router

import (
	"github.com/gin-gonic/gin"

	userhandler "gamelink/internal/handler/user"
)

// registerUserRoutes 注册用户端路由
func registerUserRoutes(api *gin.RouterGroup, authMiddleware gin.HandlerFunc, services *appServices) {
	userGroup := api.Group("/user")
	userGroup.Use(authMiddleware)
	{
		userhandler.RegisterOrderRoutes(userGroup, services.orderSvc, authMiddleware)
		userhandler.RegisterPaymentRoutes(userGroup, services.paymentSvc, authMiddleware)
		userhandler.RegisterWalletRoutes(userGroup, services.walletSvc, authMiddleware)
		userhandler.RegisterPlayerRoutes(userGroup, services.playerSvc, authMiddleware)
		userhandler.RegisterReviewRoutes(userGroup, services.reviewSvc, authMiddleware)
		userhandler.RegisterGiftRoutes(userGroup, services.giftSvc, services.serviceItemSvc, authMiddleware)
		userhandler.RegisterChatRoutes(userGroup, services.chatSvc, authMiddleware)
		userhandler.RegisterFeedRoutes(userGroup, services.feedSvc, authMiddleware)
		userhandler.RegisterBlockRoutes(userGroup, services.userBlockSvc, authMiddleware)
	}
}
