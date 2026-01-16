package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	userhandler "gamelink/internal/handler/user"
	favoriterepo "gamelink/internal/repository/favorite"
	orderrepo "gamelink/internal/repository/implementations"
	paymentrepo "gamelink/internal/repository/payment"
	playerrepo "gamelink/internal/repository/player"
	userrepo "gamelink/internal/repository/user"
	authservice "gamelink/internal/service/auth"
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
		userhandler.RegisterVipRoutes(userGroup, services.vipSvc, authMiddleware)
		userhandler.RegisterCouponRoutes(userGroup, services.couponSvc, authMiddleware)
		userhandler.RegisterRechargeRoutes(userGroup, services.rechargeSvc, authMiddleware)
		userhandler.RegisterActivityRoutes(userGroup, services.activitySvc, authMiddleware)
		userhandler.RegisterReferralRoutes(userGroup, services.referralSvc, authMiddleware)
		// 主订单路由（订单拆分与转单）
		orderGroupHandler := userhandler.NewOrderGroupHandler(services.orderSvc, services.orderGroupRepo)
		orderGroupHandler.RegisterRoutes(userGroup)
	}
}

// registerUserRoutesWithRoleSwitch 注册用户端路由（包含角色切换，需要 orm）
func registerUserRoutesWithRoleSwitch(api *gin.RouterGroup, authMiddleware gin.HandlerFunc, services *appServices, orm *gorm.DB) {
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
		userhandler.RegisterVipRoutes(userGroup, services.vipSvc, authMiddleware)
		userhandler.RegisterCouponRoutes(userGroup, services.couponSvc, authMiddleware)
		userhandler.RegisterRechargeRoutes(userGroup, services.rechargeSvc, authMiddleware)
		userhandler.RegisterActivityRoutes(userGroup, services.activitySvc, authMiddleware)
		userhandler.RegisterReferralRoutes(userGroup, services.referralSvc, authMiddleware)
		// 主订单路由（订单拆分与转单）
		orderGroupHandler := userhandler.NewOrderGroupHandler(services.orderSvc, services.orderGroupRepo)
		orderGroupHandler.RegisterRoutes(userGroup)

		// 角色切换路由（小程序用户/陪玩师切换）
		userRepo := userrepo.NewUserRepository(orm)
		playerRepo := playerrepo.NewPlayerRepository(orm)
		roleSvc := authservice.NewRoleService(userRepo, playerRepo)
		roleHandler := userhandler.NewRoleHandler(roleSvc)
		roleHandler.RegisterRoutes(userGroup)

		// 收藏路由
		favoriteRepo := favoriterepo.NewRepository(orm)
		favoriteHandler := userhandler.NewFavoriteHandler(favoriteRepo, playerRepo)
		userhandler.RegisterFavoriteRoutes(userGroup, favoriteHandler, authMiddleware)

		// 用户资料路由
		userhandler.RegisterProfileRoutes(userGroup, userRepo, authMiddleware)

		// 订单统计路由
		orderRepo := orderrepo.NewOrderRepository(orm)
		userhandler.RegisterOrderStatsRoutes(userGroup, orderRepo, authMiddleware)

		// VIP 用户信息路由
		userhandler.RegisterVipInfoRoutes(userGroup, services.vipSvc, userRepo, authMiddleware)

		// 钱包交易记录路由
		paymentRepo := paymentrepo.NewPaymentRepository(orm)
		userhandler.RegisterWalletTransactionsRoutes(userGroup, paymentRepo, authMiddleware)

		// 在线状态路由
		userhandler.RegisterPresenceRoutes(userGroup, services.presenceSvc, authMiddleware)

		// 游戏房间路由
		userhandler.RegisterGameRoomRoutes(userGroup, services.gameRoomSvc, authMiddleware)

		// 快速匹配路由
		userhandler.RegisterLFGRoutes(userGroup, services.lfgSvc, authMiddleware)

		// 语音通话路由
		userhandler.RegisterVoiceRoutes(userGroup, services.trtcSvc, authMiddleware)
	}

	// 修改密码路由（注册在 /api/v1 下，不是 /user 下）
	userRepo := userrepo.NewUserRepository(orm)
	userhandler.RegisterChangePasswordRoutes(api, userRepo, authMiddleware)
}

