package router

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"gamelink/internal/model"
	activityrepo "gamelink/internal/repository/activity"
	adminrepo "gamelink/internal/repository/admin"
	alertrepo "gamelink/internal/repository/alert"
	chatrepo "gamelink/internal/repository/chat"
	collectionentityrepo "gamelink/internal/repository/collectionentity"
	"gamelink/internal/repository/common"
	commissionrepo "gamelink/internal/repository/commission"
	contentrepo "gamelink/internal/repository/content"
	contentcategoryrepo "gamelink/internal/repository/contentcategory"
	couponrepo "gamelink/internal/repository/coupon"
	gamerepo "gamelink/internal/repository/game"
	gamerankrepo "gamelink/internal/repository/gamerank"
	orderrepo "gamelink/internal/repository/implementations"
	lfgrepo "gamelink/internal/repository/lfg"
	notificationsettingsrepo "gamelink/internal/repository/notification"
	ordermodelsrepo "gamelink/internal/repository/order"
	ordergrouprepo "gamelink/internal/repository/ordergroup"
	ordertimeoutrepo "gamelink/internal/repository/ordertimeout"
	playercertificationrepo "gamelink/internal/repository/playercertification"
	playerrankrepo "gamelink/internal/repository/playerrank"
	playerschedulerepo "gamelink/internal/repository/playerschedule"
	playerservicerepo "gamelink/internal/repository/playerservice"
	presencerepo "gamelink/internal/repository/presence"
	rechargerepo "gamelink/internal/repository/recharge"
	referralrepo "gamelink/internal/repository/referral"
	reviewdisplaysettingsrepo "gamelink/internal/repository/reviewdisplaysettings"
	routingrulerepo "gamelink/internal/repository/routingrule"
	sensitivewordrepo "gamelink/internal/repository/sensitiveword"
	serviceitemrepo "gamelink/internal/repository/serviceitem"
	settlementcompanyrepo "gamelink/internal/repository/settlementcompany"
	teamrepo "gamelink/internal/repository/team"
	uploadrepo "gamelink/internal/repository/upload"
	userrepo "gamelink/internal/repository/user"
	userblockrepo "gamelink/internal/repository/userblock"
	viprepo "gamelink/internal/repository/vip"
	withdrawrepo "gamelink/internal/repository/withdraw"
	activityservice "gamelink/internal/service/activity"
	analyticsservice "gamelink/internal/service/analytics"
	chatservice "gamelink/internal/service/chat"
	commissionservice "gamelink/internal/service/commission"
	contentservice "gamelink/internal/service/content"
	contentcategoryservice "gamelink/internal/service/contentcategory"
	couponservice "gamelink/internal/service/coupon"
	"gamelink/internal/service/external"
	gamerankservice "gamelink/internal/service/gamerank"
	gameroomservice "gamelink/internal/service/gameroom"
	giftservice "gamelink/internal/service/gift"
	itemservice "gamelink/internal/service/item"
	kpiservice "gamelink/internal/service/kpi"
	lfgservice "gamelink/internal/service/lfg"
	monitorservice "gamelink/internal/service/monitor"
	notificationservice "gamelink/internal/service/notification"
	orderservice "gamelink/internal/service/order"
	ordertimeoutservice "gamelink/internal/service/ordertimeout"
	paymentservice "gamelink/internal/service/payment"
	serviceplayer "gamelink/internal/service/player"
	playercertificationservice "gamelink/internal/service/playercertification"
	playerrankservice "gamelink/internal/service/playerrank"
	presenceservice "gamelink/internal/service/presence"
	rechargeservice "gamelink/internal/service/recharge"
	referralservice "gamelink/internal/service/referral"
	reviewservice "gamelink/internal/service/review"
	routingruleservice "gamelink/internal/service/routingrule"
	sensitivewordservice "gamelink/internal/service/sensitiveword"
	statisticsservice "gamelink/internal/service/statistics"
	teamservice "gamelink/internal/service/team"
	trtcservice "gamelink/internal/service/trtc"
	uploadservice "gamelink/internal/service/upload"
	userservice "gamelink/internal/service/user"
	userblockservice "gamelink/internal/service/userblock"
	vipservice "gamelink/internal/service/vip"
	walletservice "gamelink/internal/service/wallet"
	withdrawservice "gamelink/internal/service/withdraw"
	"gamelink/internal/ws"
	"gamelink/pkg/cache"
	"gamelink/pkg/config"
	"gamelink/pkg/scheduler"
)

// appServices 包含所有领域服务实例和调度器句柄，供路由注册使用。
type appServices struct {
	commissionSvc           *commissionservice.CommissionService
	serviceItemSvc          *itemservice.ServiceItemService
	giftSvc                 *giftservice.GiftService
	orderSvc                *orderservice.OrderService
	orderGroupRepo          ordergrouprepo.Repository // 主订单仓储
	paymentSvc              *paymentservice.PaymentService
	playerSvc               *serviceplayer.PlayerService
	reviewSvc               *orderservice.ReviewService
	disputeSvc              *orderservice.DisputeService
	earningsSvc             *userservice.EarningsService
	chatSvc                 *chatservice.ChatService
	feedSvc                 *contentservice.FeedService
	notificationSvc         *contentservice.NotificationService
	uploadSvc               *uploadservice.Service
	walletSvc               *walletservice.WalletService
	userSettingsSvc         *userservice.SettingsService
	notificationSettingsSvc *notificationservice.SettingsService
	playerServiceMgmtSvc    *serviceplayer.ServiceManagement
	playerScheduleSvc       *serviceplayer.ScheduleService
	settlementScheduler     *scheduler.SettlementScheduler
	chatRetention           *scheduler.ChatRetentionScheduler
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
	// VIP service (VIP会员)
	vipSvc *vipservice.Service
	// Coupon service (优惠券)
	couponSvc *couponservice.Service
	// Recharge service (充值)
	rechargeSvc *rechargeservice.Service
	// Activity service (活动)
	activitySvc *activityservice.Service
	// Team service (团队)
	teamSvc *teamservice.TeamService
	// Referral service (推荐)
	referralSvc *referralservice.Service
	// Withdraw routing service (提现分流)
	withdrawRoutingSvc *withdrawservice.WithdrawRoutingService
	// Presence service (在线状态)
	presenceSvc *presenceservice.Service
	// GameRoom service (游戏房间)
	gameRoomSvc *gameroomservice.Service
	// LFG service (快速匹配)
	lfgSvc *lfgservice.Service
	// TRTC service (语音通话)
	trtcSvc *trtcservice.Service
	// Business scheduler (综合业务调度器)
	businessScheduler *scheduler.BusinessScheduler
	// Referral trigger service (推荐奖励触发)
	referralTriggerSvc *referralservice.TriggerService
}

// initServices 初始化领域服务和调度任务（但不启动调度器）。
func initServices(orm *gorm.DB, cacheClient cache.Cache, cfg config.AppConfig) *appServices {
	// 仓库实例（仅在此函数内部复用）
	userRepo := userrepo.NewUserRepository(orm)
	playerRepo := userrepo.NewPlayerRepository(orm)
	gameRepo := gamerepo.NewGameRepository(orm)
	gameRankRepo := gamerankrepo.NewGameRankRepository(orm)
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
	playerServiceRepo := playerservicerepo.NewPlayerServiceRepository(orm)
	playerScheduleRepo := playerschedulerepo.NewPlayerScheduleRepository(orm)
	feedRepo := contentrepo.NewFeedRepository(orm)
	notificationRepo := contentrepo.NewNotificationRepository(orm)
	walletRepo := userrepo.NewWalletRepository(orm)
	uploadRepo := uploadrepo.NewUploadRepository(orm)
	userSettingsRepo := userrepo.NewUserSettingsRepository(orm)
	notificationSettingsRepo := notificationsettingsrepo.NewNotificationSettingRepository(orm)

	externalCfg := external.NewConfig(cfg.ExternalAPI)

	// 领域服务
	commissionSvc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)
	serviceItemSvc := itemservice.NewServiceItemService(serviceItemRepo, gameRepo, playerRepo)
	giftSvc := giftservice.NewGiftService(serviceItemRepo, orderRepo, playerRepo, userRepo, commissionRepo)
	uow := common.NewUnitOfWork(orm)
	orderSvc := orderservice.NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	orderSvc.SetTxManager(uow) // 注入事务管理器
	// 注入聊天群仓库用于订单聊天自动销毁
	orderSvc.SetChatGroupRepository(chatGroupRepo)
	// 注入主订单仓库用于订单拆分
	orderGroupRepo := ordergrouprepo.NewRepository(orm)
	orderSvc.SetOrderGroupRepository(orderGroupRepo)
	paymentSvc := paymentservice.NewPaymentService(paymentRepo, orderRepo)
	paymentSvc.SetTxManager(uow) // 注入事务管理器
	paymentSvc.SetExternalConfig(externalCfg)
	paymentSvc.SetWalletRepository(walletRepo)
	orderSvc.SetPaymentRefundor(paymentSvc)
	playerSvc := serviceplayer.NewPlayerService(playerRepo, userRepo, gameRepo, orderRepo, reviewRepo, playerTagRepo, cacheClient)
	playerServiceMgmtSvc := serviceplayer.NewServiceManagement(playerServiceRepo, playerRepo, gameRepo, gameRankRepo)
	playerScheduleSvc := serviceplayer.NewScheduleService(playerScheduleRepo, playerRepo)
	reviewSvc := orderservice.NewReviewService(reviewRepo, orderRepo, playerRepo, userRepo, reviewReplyRepo, notificationRepo)
	disputeRepo := ordermodelsrepo.NewDisputeRepository(orm)
	operationLogRepo := adminrepo.NewOperationLogRepository(orm)
	disputeSvc := orderservice.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)
	earningsSvc := userservice.NewEarningsService(playerRepo, orderRepo, withdrawRepo)
	chatSvc := chatservice.NewChatService(chatGroupRepo, chatMemberRepo, chatMessageRepo, chatReportRepo, userRepo, cacheClient)
	feedSvc := contentservice.NewFeedService(feedRepo, nil)
	notificationSvc := contentservice.NewNotificationService(notificationRepo)
	uploadSvc := uploadservice.NewService(uploadRepo, externalCfg)
	walletSvc := walletservice.NewWalletService(walletRepo, paymentRepo, orderRepo)
	userSettingsSvc := userservice.NewSettingsService(userSettingsRepo)
	notificationSettingsSvc := notificationservice.NewSettingsService(notificationSettingsRepo)

	// Create distributed lock for schedulers
	distributedLock := cache.NewDistributedLock(cacheClient)

	// 调度器（先构造，调用方负责 Start/Stop）
	settlementScheduler := scheduler.NewSettlementScheduler(commissionSvc, distributedLock)
	chatRetention := scheduler.NewChatRetentionScheduler(chatGroupRepo, chatMessageRepo, distributedLock, 30)
	businessScheduler := scheduler.NewBusinessScheduler(orm, distributedLock)

	// 推荐奖励触发服务
	referralTriggerSvc := referralservice.NewTriggerService(orm)
	// 注入推荐触发服务到订单服务（用于首单奖励）
	orderSvc.SetReferralTrigger(referralTriggerSvc)

	// Monitor services - WebSocket Hub with Redis support for horizontal scaling
	// If cacheClient uses Redis, WebSocket will enable multi-instance support via Pub/Sub
	// If cacheClient uses memory, WebSocket will run in single-instance mode (backward compatible)
	var redisClient *redis.Client
	if client := cacheClient.GetRedisClient(); client != nil {
		if rc, ok := client.(*redis.Client); ok {
			redisClient = rc
		}
	}
	wsHub := ws.NewHubWithRedis(redisClient)
	alertRepo := alertrepo.NewAlertRepository(orm)
	realtimeSvc := monitorservice.NewRealtimeService(wsHub, orm)

	// Realtime event wiring
	orderSvc.SetNotificationRepository(notificationRepo)
	orderSvc.SetWebsocketHub(wsHub)
	chatSvc.SetWebsocketHub(wsHub)
	paymentSvc.SetNotificationRepository(notificationRepo)
	paymentSvc.SetWebsocketHub(wsHub)

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

	// VIP service (VIP会员)
	vipRepo := viprepo.NewVipRepository(orm)
	vipSvc := vipservice.NewVipService(vipRepo)

	// Coupon service (优惠券)
	couponRepo := couponrepo.NewCouponRepository(orm)
	couponSvc := couponservice.NewCouponService(couponRepo)

	// Recharge service (充值)
	rechargeRepo := rechargerepo.NewRechargeRepository(orm)
	rechargeSvc := rechargeservice.NewRechargeService(rechargeRepo, walletRepo, couponSvc)
	// 注入推荐触发服务到充值服务（用于首充奖励）
	rechargeSvc.SetReferralTrigger(referralTriggerSvc)

	// Activity service (活动)
	activityRepo := activityrepo.NewActivityRepository(orm)
	activitySvc := activityservice.NewActivityService(activityRepo, couponSvc)

	// Team service (团队)
	teamRepo := teamrepo.NewTeamRepository(orm)
	teamSvc := teamservice.NewTeamService(teamRepo)

	// Referral service (推荐)
	referralRepo := referralrepo.NewReferralRepository(orm)
	referralSvc := referralservice.NewReferralService(referralRepo)

	// Withdraw routing service (提现分流)
	settlementCompanyRepo := settlementcompanyrepo.NewSettlementCompanyRepository(orm)
	withdrawRoutingSvc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementCompanyRepo)
	// Inject wallet repository for refund operations
	withdrawRoutingSvc.SetWalletRepository(walletRepo)

	// Presence service (在线状态)
	presenceRepo := presencerepo.NewPresenceRepository(orm)
	presenceSvc := presenceservice.NewPresenceServiceWithCache(presenceRepo, cacheClient)

	// GameRoom service (游戏房间)
	gameRoomSvc := gameroomservice.NewService(chatGroupRepo, chatMemberRepo, cacheClient)

	// LFG service (快速匹配)
	lfgRepo := lfgrepo.NewLFGRepository(orm)
	lfgSvc := lfgservice.NewService(lfgRepo, gameRoomSvc, cacheClient)

	// TRTC service (语音通话) - 可选，需要配置
	var trtcSvc *trtcservice.Service
	// TRTC 配置从环境变量读取，如果未配置则跳过
	// trtcConfig := &trtcservice.Config{
	//     SDKAppID:  uint64(os.Getenv("TRTC_SDK_APP_ID")),
	//     SecretKey: os.Getenv("TRTC_SECRET_KEY"),
	//     ExpireSec: 86400 * 7,
	// }
	// if trtcConfig.SDKAppID > 0 && trtcConfig.SecretKey != "" {
	//     trtcSvc = trtcservice.NewService(trtcConfig, chatGroupRepo, chatMemberRepo, cacheClient)
	// }

	return &appServices{
		commissionSvc:           commissionSvc,
		serviceItemSvc:          serviceItemSvc,
		giftSvc:                 giftSvc,
		orderSvc:                orderSvc,
		orderGroupRepo:          orderGroupRepo,
		paymentSvc:              paymentSvc,
		playerSvc:               playerSvc,
		playerServiceMgmtSvc:    playerServiceMgmtSvc,
		playerScheduleSvc:       playerScheduleSvc,
		reviewSvc:               reviewSvc,
		disputeSvc:              disputeSvc,
		earningsSvc:             earningsSvc,
		chatSvc:                 chatSvc,
		feedSvc:                 feedSvc,
		notificationSvc:         notificationSvc,
		uploadSvc:               uploadSvc,
		walletSvc:               walletSvc,
		userSettingsSvc:         userSettingsSvc,
		notificationSettingsSvc: notificationSettingsSvc,
		settlementScheduler:     settlementScheduler,
		chatRetention:           chatRetention,
		wsHub:                   wsHub,
		realtimeSvc:             realtimeSvc,
		alertRepo:               alertRepo,
		analyticsSvc:            analyticsSvc,
		kpiSvc:                  kpiSvc,
		tagSvc:                  tagSvc,
		batchSvc:                batchSvc,
		reviewStatsSvc:          reviewStatsSvc,
		adminFeedSvc:            adminFeedSvc,
		chatModerationSvc:       chatModerationSvc,
		feedReportSvc:           feedReportSvc,
		contentStatsSvc:         contentStatsSvc,
		contentCategorySvc:      contentCategorySvc,
		sensitiveWordSvc:        sensitiveWordSvc,
		reviewSettingsSvc:       reviewSettingsSvc,
		routingRuleSvc:          routingRuleSvc,
		statisticsSvc:           statisticsSvc,
		statisticsEvaluator:     statisticsEvaluator,
		statisticsHooks:         statisticsHooks,
		// Player rank services
		gameRankSvc:            gameRankSvc,
		playerRankSvc:          playerRankSvc,
		playerCertificationSvc: playerCertificationSvc,
		// Order timeout service
		orderTimeoutSvc: orderTimeoutSvc,
		// User block service
		userBlockSvc: userBlockSvc,
		// VIP service
		vipSvc: vipSvc,
		// Coupon service
		couponSvc: couponSvc,
		// Recharge service
		rechargeSvc: rechargeSvc,
		// Activity service
		activitySvc: activitySvc,
		// Team service
		teamSvc: teamSvc,
		// Referral service
		referralSvc: referralSvc,
		// Withdraw routing service
		withdrawRoutingSvc: withdrawRoutingSvc,
		// Presence service
		presenceSvc: presenceSvc,
		// GameRoom service
		gameRoomSvc: gameRoomSvc,
		// LFG service
		lfgSvc: lfgSvc,
		// TRTC service
		trtcSvc: trtcSvc,
		// Business scheduler
		businessScheduler: businessScheduler,
		// Referral trigger service
		referralTriggerSvc: referralTriggerSvc,
	}
}
