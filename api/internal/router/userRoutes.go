package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	userhandler "gamelink/internal/handler/user"
	adminrepo "gamelink/internal/repository/admin"
	favoriterepo "gamelink/internal/repository/favorite"
	playerrepo "gamelink/internal/repository/player"
	authservice "gamelink/internal/service/auth"
)

// registerUserRoutes 注册用户端路由
func registerUserRoutes(api *gin.RouterGroup, authMiddleware gin.HandlerFunc, services *appServices) {
	userGroup := api.Group("/user")
	userGroup.Use(authMiddleware)
	{
		userhandler.RegisterOrderRoutes(userGroup, services.orderSvc, authMiddleware)
		userhandler.RegisterDisputeRoutes(userGroup, services.disputeSvc, authMiddleware)
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
		userhandler.RegisterSettingsRoutes(userGroup, services.userSettingsSvc, services.notificationSettingsSvc, authMiddleware)
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
		userhandler.RegisterDisputeRoutes(userGroup, services.disputeSvc, authMiddleware)
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
		userhandler.RegisterSettingsRoutes(userGroup, services.userSettingsSvc, services.notificationSettingsSvc, authMiddleware)
		userhandler.RegisterCustomerServiceRoutes(
			userGroup,
			services.chatSvc,
			services.userRepo,
			adminrepo.NewRoleRepository(orm),
			authMiddleware,
		)
		// 主订单路由（订单拆分与转单）
		orderGroupHandler := userhandler.NewOrderGroupHandler(services.orderSvc, services.orderGroupRepo)
		orderGroupHandler.RegisterRoutes(userGroup)

		// 角色切换路由（小程序用户/陪玩师切换）
		roleSvc := authservice.NewRoleService(services.userRepo, services.playerRepo)
		roleHandler := userhandler.NewRoleHandler(roleSvc)
		roleHandler.RegisterRoutes(userGroup)

		// 收藏路由
		favoriteRepo := favoriterepo.NewRepository(orm)
		playerRepo := playerrepo.NewPlayerRepository(orm) // favoriteHandler 需要 playerrepo 子包类型
		favoriteHandler := userhandler.NewFavoriteHandler(favoriteRepo, playerRepo)
		userhandler.RegisterFavoriteRoutes(userGroup, favoriteHandler, authMiddleware)

		// 用户资料路由
		userhandler.RegisterProfileRoutes(userGroup, services.userRepo, authMiddleware)

		// 订单统计路由
		userhandler.RegisterOrderStatsRoutes(userGroup, services.orderRepo, authMiddleware)

		// VIP 用户信息路由
		userhandler.RegisterVipInfoRoutes(userGroup, services.vipSvc, services.userRepo, authMiddleware)

		// 钱包交易记录路由
		userhandler.RegisterWalletTransactionsRoutes(userGroup, services.paymentRepo, authMiddleware)

		// 在线状态路由
		userhandler.RegisterPresenceRoutes(userGroup, services.presenceSvc, authMiddleware)

		// 游戏房间路由
		userhandler.RegisterGameRoomRoutes(userGroup, services.gameRoomSvc, authMiddleware)

		// 快速匹配路由
		userhandler.RegisterLFGRoutes(userGroup, services.lfgSvc, authMiddleware)

		// 语音通话路由
		userhandler.RegisterVoiceRoutes(userGroup, services.trtcSvc, authMiddleware)
	}

	// 兼容旧版前端路径：/users/*
	legacyUsersGroup := api.Group("/users")
	legacyUsersGroup.Use(authMiddleware)
	{
		userhandler.RegisterChatRoutes(legacyUsersGroup, services.chatSvc, authMiddleware)
		userhandler.RegisterProfileRoutes(legacyUsersGroup, services.userRepo, authMiddleware)
	}

	// 上传相关路由（注册在 /api/v1 下）
	if services.uploadSvc != nil {
		userhandler.RegisterUploadRoutes(api, authMiddleware, services.uploadSvc)
	}

	// 修改密码路由（注册在 /api/v1 下，不是 /user 下）
	userhandler.RegisterChangePasswordRoutes(api, services.userRepo, authMiddleware)
}
