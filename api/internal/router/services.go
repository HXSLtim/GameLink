package router

import (
	"gamelink/internal/model"
	adminrepo "gamelink/internal/repository/admin"
	alertrepo "gamelink/internal/repository/alert"
	chatrepo "gamelink/internal/repository/chat"
	collectionentityrepo "gamelink/internal/repository/collectionentity"
	commissionrepo "gamelink/internal/repository/commission"
	contentrepo "gamelink/internal/repository/content"
	contentcategoryrepo "gamelink/internal/repository/contentcategory"
	gamerepo "gamelink/internal/repository/game"
	gamerankrepo "gamelink/internal/repository/gamerank"
	orderrepo "gamelink/internal/repository/implementations"
	ordermodelsrepo "gamelink/internal/repository/order"
	ordertimeoutrepo "gamelink/internal/repository/ordertimeout"
	playercertificationrepo "gamelink/internal/repository/playercertification"
	playerrankrepo "gamelink/internal/repository/playerrank"
	reviewdisplaysettingsrepo "gamelink/internal/repository/reviewdisplaysettings"
	routingrulerepo "gamelink/internal/repository/routingrule"
	sensitivewordrepo "gamelink/internal/repository/sensitiveword"
	serviceitemrepo "gamelink/internal/repository/serviceitem"
	userrepo "gamelink/internal/repository/user"
	userblockrepo "gamelink/internal/repository/userblock"
	withdrawrepo "gamelink/internal/repository/withdraw"
	analyticsservice "gamelink/internal/service/analytics"
	chatservice "gamelink/internal/service/chat"
	commissionservice "gamelink/internal/service/commission"
	contentservice "gamelink/internal/service/content"
	contentcategoryservice "gamelink/internal/service/contentcategory"
	gamerankservice "gamelink/internal/service/gamerank"
	giftservice "gamelink/internal/service/gift"
	itemservice "gamelink/internal/service/item"
	kpiservice "gamelink/internal/service/kpi"
	monitorservice "gamelink/internal/service/monitor"
	orderservice "gamelink/internal/service/order"
	ordertimeoutservice "gamelink/internal/service/ordertimeout"
	paymentservice "gamelink/internal/service/payment"
	serviceplayer "gamelink/internal/service/player"
	playercertificationservice "gamelink/internal/service/playercertification"
	playerrankservice "gamelink/internal/service/playerrank"
	reviewservice "gamelink/internal/service/review"
	routingruleservice "gamelink/internal/service/routingrule"
	sensitivewordservice "gamelink/internal/service/sensitiveword"
	statisticsservice "gamelink/internal/service/statistics"
	userservice "gamelink/internal/service/user"
	userblockservice "gamelink/internal/service/userblock"
	walletservice "gamelink/internal/service/wallet"
	"gamelink/internal/ws"
	"gamelink/pkg/cache"
	"gamelink/pkg/scheduler"

	"gorm.io/gorm"
)

// appServices 包含所有领域服务实例和调度器句柄，供路由注册使用。
type appServices struct {
	commissionSvc       *commissionservice.CommissionService
	serviceItemSvc      *itemservice.ServiceItemService
	giftSvc             *giftservice.GiftService
	orderSvc            *orderservice.OrderService
	paymentSvc          *paymentservice.PaymentService
	playerSvc           *serviceplayer.PlayerService
	reviewSvc           *orderservice.ReviewService
	disputeSvc          *orderservice.DisputeService
	earningsSvc         *userservice.EarningsService
	chatSvc             *chatservice.ChatService
	feedSvc             *contentservice.FeedService
	notificationSvc     *contentservice.NotificationService
	walletSvc           *walletservice.WalletService
	settlementScheduler *scheduler.SettlementScheduler
	chatRetention       *scheduler.ChatRetentionScheduler
	// Monitor services
	wsHub       *ws.Hub
	realtimeSvc *monitorservice.RealtimeService
	alertRepo   model.AlertRepository
	// Analytics service
	analyticsSvc *analyticsservice.AnalyticsService
	// KPI service
	kpiSvc *kpiservice.KPIService
	// User management services
	tagSvc   *userservice.UserTagService
	batchSvc *userservice.BatchOperationService
	// Review stats service
	reviewStatsSvc *reviewservice.ReviewStatsService
	// Content management services
	adminFeedSvc       *contentservice.AdminFeedService
	chatModerationSvc  *contentservice.ChatModerationService
	feedReportSvc      *contentservice.FeedReportService
	contentStatsSvc    *contentservice.ContentStatsService
	contentCategorySvc *contentcategoryservice.ContentCategoryService
	// Sensitive word service
	sensitiveWordSvc *sensitivewordservice.SensitiveWordService
	// Review settings service
	reviewSettingsSvc *reviewservice.SettingsService
	// Routing rule service
	routingRuleSvc *routingruleservice.RoutingRuleService
	// Statistics service (统计指标)
	statisticsSvc       *statisticsservice.Service
	statisticsEvaluator *statisticsservice.TagEvaluator
	statisticsHooks     *statisticsservice.EventHooks
	// Player rank services (陪玩师等级/认证)
	gameRankSvc            *gamerankservice.GameRankService
	playerRankSvc          *playerrankservice.PlayerRankService
	playerCertificationSvc *playercertificationservice.PlayerCertificationService
	// Order timeout service (订单超时处理)
	orderTimeoutSvc *ordertimeoutservice.OrderTimeoutService
	// User block service (用户拉黑)
	userBlockSvc *userblockservice.UserBlockService
}

// initServices 初始化领域服务和调度任务（但不启动调度器）。
func initServices(orm *gorm.DB, cacheClient cache.Cache) *appServices {
	// 仓库实例（仅在此函数内部复用）
	userRepo := userrepo.NewUserRepository(orm)
	playerRepo := userrepo.NewPlayerRepository(orm)
	gameRepo := gamerepo.NewGameRepository(orm)
	orderRepo := orderrepo.NewOrderRepository(orm)
	chatGroupRepo := chatrepo.NewChatGroupRepository(orm)
	chatMemberRepo := chatrepo.NewChatMemberRepository(orm)
	chatMessageRepo := chatrepo.NewChatMessageRepository(orm)
	chatReportRepo := chatrepo.NewChatReportRepository(orm)
	paymentRepo := ordermodelsrepo.NewPaymentRepository(orm)
	reviewRepo := ordermodelsrepo.NewReviewRepository(orm)
	reviewReplyRepo := ordermodelsrepo.NewReviewReplyRepository(orm)
	playerTagRepo := userrepo.NewPlayerTagRepository(orm)
	withdrawRepo := withdrawrepo.NewWithdrawRepository(orm)
	commissionRepo := commissionrepo.NewCommissionRepository(orm)
	serviceItemRepo := serviceitemrepo.NewServiceItemRepository(orm)
	feedRepo := contentrepo.NewFeedRepository(orm)
	notificationRepo := contentrepo.NewNotificationRepository(orm)
	walletRepo := userrepo.NewWalletRepository(orm)

	// 领域服务
	commissionSvc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)
	serviceItemSvc := itemservice.NewServiceItemService(serviceItemRepo, gameRepo, playerRepo)
	giftSvc := giftservice.NewGiftService(serviceItemRepo, orderRepo, playerRepo, commissionRepo)
	orderSvc := orderservice.NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	// 注入聊天群仓库用于订单聊天自动销毁
	orderSvc.SetChatGroupRepository(chatGroupRepo)
	paymentSvc := paymentservice.NewPaymentService(paymentRepo, orderRepo)
	playerSvc := serviceplayer.NewPlayerService(playerRepo, userRepo, gameRepo, orderRepo, reviewRepo, playerTagRepo, cacheClient)
	reviewSvc := orderservice.NewReviewService(reviewRepo, orderRepo, playerRepo, userRepo, reviewReplyRepo, notificationRepo)
	disputeRepo := ordermodelsrepo.NewDisputeRepository(orm)
	operationLogRepo := adminrepo.NewOperationLogRepository(orm)
	disputeSvc := orderservice.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)
	earningsSvc := userservice.NewEarningsService(playerRepo, orderRepo, withdrawRepo)
	chatSvc := chatservice.NewChatService(chatGroupRepo, chatMemberRepo, chatMessageRepo, chatReportRepo, cacheClient)
	feedSvc := contentservice.NewFeedService(feedRepo, nil)
	notificationSvc := contentservice.NewNotificationService(notificationRepo)
	walletSvc := walletservice.NewWalletService(walletRepo, paymentRepo, orderRepo)

	// 调度器（先构造，调用方负责 Start/Stop）
	settlementScheduler := scheduler.NewSettlementScheduler(commissionSvc)
	chatRetention := scheduler.NewChatRetentionScheduler(chatGroupRepo, chatMessageRepo, 30)

	// Monitor services
	wsHub := ws.NewHub()
	alertRepo := alertrepo.NewAlertRepository(orm)
	realtimeSvc := monitorservice.NewRealtimeService(wsHub, orm)

	// Analytics service
	analyticsSvc := analyticsservice.NewAnalyticsService(orm)

	// KPI service
	kpiSvc := kpiservice.NewKPIService(orm)

	// User management services
	tagRepo := userrepo.NewUserTagRepository(orm)
	notifRepo := contentrepo.NewNotificationRepository(orm)
	tagSvc := userservice.NewUserTagService(tagRepo, userRepo, cacheClient)
	batchSvc := userservice.NewBatchOperationService(orm, userRepo, tagRepo, notifRepo)

	// Review stats service
	reviewStatsSvc := reviewservice.NewReviewStatsService(reviewRepo)

	// Review settings service
	reviewDisplaySettingsRepo := reviewdisplaysettingsrepo.New(orm)
	reviewSettingsSvc := reviewservice.NewSettingsService(reviewDisplaySettingsRepo)

	// Content management services
	contentCategoryRepo := contentcategoryrepo.NewContentCategoryRepository(orm)
	sensitiveWordRepo := sensitivewordrepo.NewSensitiveWordRepository(orm)
	sensitiveWordSvc := sensitivewordservice.NewSensitiveWordService(sensitiveWordRepo)
	adminFeedSvc := contentservice.NewAdminFeedService(feedRepo, sensitiveWordSvc, operationLogRepo)
	chatModerationSvc := contentservice.NewChatModerationService(chatMessageRepo, chatMemberRepo, sensitiveWordSvc, operationLogRepo)
	feedReportSvc := contentservice.NewFeedReportService(feedRepo, operationLogRepo)
	contentStatsSvc := contentservice.NewContentStatsService(feedRepo, chatMessageRepo)
	contentCategorySvc := contentcategoryservice.NewContentCategoryService(contentCategoryRepo)

	// Routing rule service
	routingRuleRepo := routingrulerepo.NewRoutingRuleRepository(orm)
	collectionEntityRepo := collectionentityrepo.NewCollectionEntityRepository(orm)
	routingRuleSvc := routingruleservice.NewRoutingRuleService(routingRuleRepo, collectionEntityRepo)

	// Statistics service (统计指标)
	statisticsSvc := statisticsservice.NewService(orm)
	statisticsEvaluator := statisticsservice.NewTagEvaluator(orm)
	statisticsHooks := statisticsservice.NewEventHooks(statisticsSvc)

	// Player rank services (陪玩师等级/认证)
	gameRankRepo := gamerankrepo.NewGameRankRepository(orm)
	playerRankRepo := playerrankrepo.NewPlayerRankRepository(orm)
	playerCertificationRepo := playercertificationrepo.NewPlayerCertificationRepository(orm)
	gameRankSvc := gamerankservice.NewGameRankService(gameRankRepo, gameRepo)
	playerRankSvc := playerrankservice.NewPlayerRankService(playerRankRepo, gameRankRepo, playerRepo, gameRepo)
	playerCertificationSvc := playercertificationservice.NewPlayerCertificationService(playerCertificationRepo, playerRepo)

	// Order timeout service (订单超时处理)
	orderTimeoutRepo := ordertimeoutrepo.NewOrderTimeoutRepository(orm)
	orderTimeoutSvc := ordertimeoutservice.NewOrderTimeoutService(orderTimeoutRepo, orderRepo, userRepo)

	// User block service (用户拉黑)
	userBlockRepo := userblockrepo.NewUserBlockRepository(orm)
	userBlockSvc := userblockservice.NewUserBlockService(userBlockRepo, userRepo)

	return &appServices{
		commissionSvc:       commissionSvc,
		serviceItemSvc:      serviceItemSvc,
		giftSvc:             giftSvc,
		orderSvc:            orderSvc,
		paymentSvc:          paymentSvc,
		playerSvc:           playerSvc,
		reviewSvc:           reviewSvc,
		disputeSvc:          disputeSvc,
		earningsSvc:         earningsSvc,
		chatSvc:             chatSvc,
		feedSvc:             feedSvc,
		notificationSvc:     notificationSvc,
		walletSvc:           walletSvc,
		settlementScheduler: settlementScheduler,
		chatRetention:       chatRetention,
		wsHub:               wsHub,
		realtimeSvc:         realtimeSvc,
		alertRepo:           alertRepo,
		analyticsSvc:        analyticsSvc,
		kpiSvc:              kpiSvc,
		tagSvc:              tagSvc,
		batchSvc:            batchSvc,
		reviewStatsSvc:      reviewStatsSvc,
		adminFeedSvc:        adminFeedSvc,
		chatModerationSvc:   chatModerationSvc,
		feedReportSvc:       feedReportSvc,
		contentStatsSvc:     contentStatsSvc,
		contentCategorySvc:  contentCategorySvc,
		sensitiveWordSvc:    sensitiveWordSvc,
		reviewSettingsSvc:   reviewSettingsSvc,
		routingRuleSvc:      routingRuleSvc,
		statisticsSvc:       statisticsSvc,
		statisticsEvaluator: statisticsEvaluator,
		statisticsHooks:     statisticsHooks,
		// Player rank services
		gameRankSvc:            gameRankSvc,
		playerRankSvc:          playerRankSvc,
		playerCertificationSvc: playerCertificationSvc,
		// Order timeout service
		orderTimeoutSvc: orderTimeoutSvc,
		// User block service
		userBlockSvc: userBlockSvc,
	}
}
