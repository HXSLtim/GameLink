package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	adminhandler "gamelink/internal/handler/admin"
	"gamelink/internal/model"
	"gamelink/internal/repository/common"
	"gamelink/internal/repository/game"
	"gamelink/internal/repository/implementations"
	"gamelink/internal/repository/menu"
	"gamelink/internal/repository/order"
	"gamelink/internal/repository/permission"
	"gamelink/internal/repository/review"
	"gamelink/internal/repository/role"
	"gamelink/internal/repository/sensitiveword"
	"gamelink/internal/repository/serviceitem"
	"gamelink/internal/repository/stats"
	"gamelink/internal/repository/user"
	"gamelink/internal/repository/wallet"
	adminservice "gamelink/internal/service/admin"
	"gamelink/pkg/cache"
	"gamelink/pkg/testutil"
)

// TestReviewModerationFlow 测试完整的评价审核流程
// 需求: 2.1, 2.2, 2.3
func TestReviewModerationFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateReviewModerationModels(t, db)

	seed := seedReviewModerationData(t, db)
	router, _ := setupReviewModerationRouter(t, db, seed.adminUserID)

	ctx := context.Background()
	reviewRepo := review.NewReviewRepository(db)

	// 1. 创建待审核评价
	pendingReview := &model.Review{
		OrderID:  seed.orderID,
		UserID:   seed.userID,
		PlayerID: seed.playerID,
		Score:    model.Rating(5),
		Content:  "Great service!",
		Status:   model.ReviewStatusPending,
	}
	err := reviewRepo.Create(ctx, pendingReview)
	require.NoError(t, err)

	// 2. 获取待审核评价列表
	listResp := doJSON(router, http.MethodGet, "/api/v1/admin/reviews/pending?page=1&pageSize=20", nil, "")
	require.Equal(t, http.StatusOK, listResp.Code, "list pending reviews should succeed")

	var listParsed apiResp[[]model.Review]
	err = json.Unmarshal(listResp.Body.Bytes(), &listParsed)
	require.NoError(t, err)
	assert.True(t, listParsed.Success)
	assert.GreaterOrEqual(t, len(listParsed.Data), 1, "should have at least one pending review")

	// 验证评价在待审核列表中
	found := false
	for _, r := range listParsed.Data {
		if r.ID == pendingReview.ID {
			found = true
			assert.Equal(t, model.ReviewStatusPending, r.Status)
			break
		}
	}
	assert.True(t, found, "pending review should be in the list")

	// 3. 批准评价
	approveResp := doJSON(router, http.MethodPut, "/api/v1/admin/reviews/"+uintToStr(pendingReview.ID)+"/approve", nil, "")
	require.Equal(t, http.StatusOK, approveResp.Code, "approve review should succeed")

	var approveParsed apiResp[any]
	err = json.Unmarshal(approveResp.Body.Bytes(), &approveParsed)
	require.NoError(t, err)
	assert.True(t, approveParsed.Success)
	assert.Equal(t, "review approved", approveParsed.Message)

	// 4. 验证评价状态已更新为已批准
	approvedReview, err := reviewRepo.Get(ctx, pendingReview.ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewStatusApproved, approvedReview.Status)

	// 5. 创建另一个待审核评价用于拒绝测试
	rejectReview := &model.Review{
		OrderID:  seed.orderID,
		UserID:   seed.userID,
		PlayerID: seed.playerID,
		Score:    model.Rating(1),
		Content:  "Bad service",
		Status:   model.ReviewStatusPending,
	}
	err = reviewRepo.Create(ctx, rejectReview)
	require.NoError(t, err)

	// 6. 拒绝评价
	rejectPayload := map[string]interface{}{
		"reason": "内容不符合规范",
	}
	rejectResp := doJSON(router, http.MethodPut, "/api/v1/admin/reviews/"+uintToStr(rejectReview.ID)+"/reject", rejectPayload, "")
	require.Equal(t, http.StatusOK, rejectResp.Code, "reject review should succeed")

	var rejectParsed apiResp[any]
	err = json.Unmarshal(rejectResp.Body.Bytes(), &rejectParsed)
	require.NoError(t, err)
	assert.True(t, rejectParsed.Success)
	assert.Equal(t, "review rejected", rejectParsed.Message)

	// 7. 验证评价状态已更新为已拒绝
	rejectedReview, err := reviewRepo.Get(ctx, rejectReview.ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewStatusRejected, rejectedReview.Status)
	assert.Equal(t, "内容不符合规范", rejectedReview.RejectionReason)

	// 8. 验证不能重复审核已批准的评价
	reApproveResp := doJSON(router, http.MethodPut, "/api/v1/admin/reviews/"+uintToStr(pendingReview.ID)+"/approve", nil, "")
	assert.NotEqual(t, http.StatusOK, reApproveResp.Code, "should not be able to re-approve")

	// 9. 验证不能重复审核已拒绝的评价
	reRejectResp := doJSON(router, http.MethodPut, "/api/v1/admin/reviews/"+uintToStr(rejectReview.ID)+"/reject", rejectPayload, "")
	assert.NotEqual(t, http.StatusOK, reRejectResp.Code, "should not be able to re-reject")
}

// TestSensitiveWordAutoMarking 测试敏感词自动标记
// 需求: 2.4
func TestSensitiveWordAutoMarking(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateReviewModerationModels(t, db)

	seed := seedReviewModerationData(t, db)
	router, _ := setupReviewModerationRouter(t, db, seed.adminUserID)

	ctx := context.Background()
	reviewRepo := review.NewReviewRepository(db)
	sensitiveWordRepo := sensitiveword.NewSensitiveWordRepository(db)

	// 1. 添加敏感词
	sensitiveWords := []model.SensitiveWord{
		{
			Word:     "垃圾",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityHigh,
		},
		{
			Word:     "骗子",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityHigh,
		},
		{
			Word:     "差评",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityMedium,
		},
	}

	for _, sw := range sensitiveWords {
		err := sensitiveWordRepo.Create(ctx, &sw)
		require.NoError(t, err)
	}

	// 2. 创建包含敏感词的评价
	reviewWithSensitiveWord := &model.Review{
		OrderID:  seed.orderID,
		UserID:   seed.userID,
		PlayerID: seed.playerID,
		Score:    model.Rating(1),
		Content:  "这个陪玩师是垃圾，简直是骗子！",
		Status:   model.ReviewStatusPending,
	}
	err := reviewRepo.Create(ctx, reviewWithSensitiveWord)
	require.NoError(t, err)

	// 3. 获取待审核评价列表，验证包含敏感词的评价被标记
	listResp := doJSON(router, http.MethodGet, "/api/v1/admin/reviews/pending?page=1&pageSize=20", nil, "")
	require.Equal(t, http.StatusOK, listResp.Code)

	var listParsed apiResp[[]model.Review]
	err = json.Unmarshal(listResp.Body.Bytes(), &listParsed)
	require.NoError(t, err)

	// 验证评价在列表中
	found := false
	for _, r := range listParsed.Data {
		if r.ID == reviewWithSensitiveWord.ID {
			found = true
			assert.Equal(t, model.ReviewStatusPending, r.Status)
			// 评价应该包含敏感词
			assert.Contains(t, r.Content, "垃圾")
			assert.Contains(t, r.Content, "骗子")
			break
		}
	}
	assert.True(t, found, "review with sensitive words should be in pending list")

	// 4. 创建不包含敏感词的评价
	reviewWithoutSensitiveWord := &model.Review{
		OrderID:  seed.orderID,
		UserID:   seed.userID,
		PlayerID: seed.playerID,
		Score:    model.Rating(5),
		Content:  "服务很好，非常满意！",
		Status:   model.ReviewStatusPending,
	}
	err = reviewRepo.Create(ctx, reviewWithoutSensitiveWord)
	require.NoError(t, err)

	// 5. 验证不包含敏感词的评价也在待审核列表中
	listResp2 := doJSON(router, http.MethodGet, "/api/v1/admin/reviews/pending?page=1&pageSize=20", nil, "")
	require.Equal(t, http.StatusOK, listResp2.Code)

	var listParsed2 apiResp[[]model.Review]
	err = json.Unmarshal(listResp2.Body.Bytes(), &listParsed2)
	require.NoError(t, err)

	found = false
	for _, r := range listParsed2.Data {
		if r.ID == reviewWithoutSensitiveWord.ID {
			found = true
			assert.Equal(t, model.ReviewStatusPending, r.Status)
			break
		}
	}
	assert.True(t, found, "review without sensitive words should also be in pending list")

	// 6. 管理员可以批准不包含敏感词的评价
	approveResp := doJSON(router, http.MethodPut, "/api/v1/admin/reviews/"+uintToStr(reviewWithoutSensitiveWord.ID)+"/approve", nil, "")
	require.Equal(t, http.StatusOK, approveResp.Code)

	// 7. 管理员应该拒绝包含敏感词的评价
	rejectPayload := map[string]interface{}{
		"reason": "包含敏感词汇",
	}
	rejectResp := doJSON(router, http.MethodPut, "/api/v1/admin/reviews/"+uintToStr(reviewWithSensitiveWord.ID)+"/reject", rejectPayload, "")
	require.Equal(t, http.StatusOK, rejectResp.Code)

	// 8. 验证最终状态
	approvedReview, err := reviewRepo.Get(ctx, reviewWithoutSensitiveWord.ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewStatusApproved, approvedReview.Status)

	rejectedReview, err := reviewRepo.Get(ctx, reviewWithSensitiveWord.ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewStatusRejected, rejectedReview.Status)
	assert.Equal(t, "包含敏感词汇", rejectedReview.RejectionReason)
}

// TestBatchModeration 测试批量审核
// 需求: 2.5
func TestBatchModeration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateReviewModerationModels(t, db)

	seed := seedReviewModerationData(t, db)
	router, _ := setupReviewModerationRouter(t, db, seed.adminUserID)

	ctx := context.Background()
	reviewRepo := review.NewReviewRepository(db)

	// 1. 创建多个待审核评价
	reviewIDs := make([]uint64, 0, 5)
	for i := 0; i < 5; i++ {
		pendingReview := &model.Review{
			OrderID:  seed.orderID,
			UserID:   seed.userID,
			PlayerID: seed.playerID,
			Score:    model.Rating(4 + i%2), // 4 or 5
			Content:  "Test review " + string(rune('A'+i)),
			Status:   model.ReviewStatusPending,
		}
		err := reviewRepo.Create(ctx, pendingReview)
		require.NoError(t, err)
		reviewIDs = append(reviewIDs, pendingReview.ID)
	}

	// 2. 批量批准前3个评价
	batchApprovePayload := map[string]interface{}{
		"reviewIds": reviewIDs[:3],
	}
	approveResp := doJSON(router, http.MethodPut, "/api/v1/admin/reviews/batch-approve", batchApprovePayload, "")
	require.Equal(t, http.StatusOK, approveResp.Code, "batch approve should succeed")

	var approveParsed apiResp[any]
	err := json.Unmarshal(approveResp.Body.Bytes(), &approveParsed)
	require.NoError(t, err)
	assert.True(t, approveParsed.Success)
	assert.Equal(t, "reviews approved", approveParsed.Message)

	// 3. 验证前3个评价状态已更新为已批准
	for i := 0; i < 3; i++ {
		approvedReview, err := reviewRepo.Get(ctx, reviewIDs[i])
		require.NoError(t, err)
		assert.Equal(t, model.ReviewStatusApproved, approvedReview.Status, "review %d should be approved", i)
	}

	// 4. 验证后2个评价仍然是待审核状态
	for i := 3; i < 5; i++ {
		pendingReview, err := reviewRepo.Get(ctx, reviewIDs[i])
		require.NoError(t, err)
		assert.Equal(t, model.ReviewStatusPending, pendingReview.Status, "review %d should still be pending", i)
	}

	// 5. 批量拒绝后2个评价
	batchRejectPayload := map[string]interface{}{
		"reviewIds": reviewIDs[3:],
		"reason":    "批量拒绝测试",
	}
	rejectResp := doJSON(router, http.MethodPut, "/api/v1/admin/reviews/batch-reject", batchRejectPayload, "")
	require.Equal(t, http.StatusOK, rejectResp.Code, "batch reject should succeed")

	var rejectParsed apiResp[any]
	err = json.Unmarshal(rejectResp.Body.Bytes(), &rejectParsed)
	require.NoError(t, err)
	assert.True(t, rejectParsed.Success)
	assert.Equal(t, "reviews rejected", rejectParsed.Message)

	// 6. 验证后2个评价状态已更新为已拒绝
	for i := 3; i < 5; i++ {
		rejectedReview, err := reviewRepo.Get(ctx, reviewIDs[i])
		require.NoError(t, err)
		assert.Equal(t, model.ReviewStatusRejected, rejectedReview.Status, "review %d should be rejected", i)
		assert.Equal(t, "批量拒绝测试", rejectedReview.RejectionReason)
	}

	// 7. 测试批量操作的原子性：尝试批量批准包含已批准评价的列表
	mixedIDs := []uint64{reviewIDs[0], reviewIDs[3]} // 一个已批准，一个已拒绝
	mixedPayload := map[string]interface{}{
		"reviewIds": mixedIDs,
	}
	mixedResp := doJSON(router, http.MethodPut, "/api/v1/admin/reviews/batch-approve", mixedPayload, "")
	assert.NotEqual(t, http.StatusOK, mixedResp.Code, "batch approve with non-pending reviews should fail")

	// 8. 验证原子性：所有评价状态应该保持不变
	review0, err := reviewRepo.Get(ctx, reviewIDs[0])
	require.NoError(t, err)
	assert.Equal(t, model.ReviewStatusApproved, review0.Status, "review 0 should still be approved")

	review3, err := reviewRepo.Get(ctx, reviewIDs[3])
	require.NoError(t, err)
	assert.Equal(t, model.ReviewStatusRejected, review3.Status, "review 3 should still be rejected")
}

// TestBatchModerationEmptyList 测试批量审核空列表
func TestBatchModerationEmptyList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateReviewModerationModels(t, db)

	seed := seedReviewModerationData(t, db)
	router, _ := setupReviewModerationRouter(t, db, seed.adminUserID)

	// 测试空列表批量批准
	emptyApprovePayload := map[string]interface{}{
		"reviewIds": []uint64{},
	}
	approveResp := doJSON(router, http.MethodPut, "/api/v1/admin/reviews/batch-approve", emptyApprovePayload, "")
	assert.NotEqual(t, http.StatusOK, approveResp.Code, "batch approve with empty list should fail")

	// 测试空列表批量拒绝
	emptyRejectPayload := map[string]interface{}{
		"reviewIds": []uint64{},
		"reason":    "test",
	}
	rejectResp := doJSON(router, http.MethodPut, "/api/v1/admin/reviews/batch-reject", emptyRejectPayload, "")
	assert.NotEqual(t, http.StatusOK, rejectResp.Code, "batch reject with empty list should fail")
}

// TestRejectReviewWithoutReason 测试拒绝评价时必须提供原因
func TestRejectReviewWithoutReason(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateReviewModerationModels(t, db)

	seed := seedReviewModerationData(t, db)
	router, _ := setupReviewModerationRouter(t, db, seed.adminUserID)

	ctx := context.Background()
	reviewRepo := review.NewReviewRepository(db)

	// 创建待审核评价
	pendingReview := &model.Review{
		OrderID:  seed.orderID,
		UserID:   seed.userID,
		PlayerID: seed.playerID,
		Score:    model.Rating(3),
		Content:  "Test review",
		Status:   model.ReviewStatusPending,
	}
	err := reviewRepo.Create(ctx, pendingReview)
	require.NoError(t, err)

	// 尝试拒绝评价但不提供原因
	emptyReasonPayload := map[string]interface{}{
		"reason": "",
	}
	rejectResp := doJSON(router, http.MethodPut, "/api/v1/admin/reviews/"+uintToStr(pendingReview.ID)+"/reject", emptyReasonPayload, "")
	assert.NotEqual(t, http.StatusOK, rejectResp.Code, "reject without reason should fail")

	// 验证评价状态未改变
	review, err := reviewRepo.Get(ctx, pendingReview.ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewStatusPending, review.Status, "review should still be pending")
}

// Helper functions

type reviewModerationSeed struct {
	adminUserID uint64
	userID      uint64
	playerID    uint64
	orderID     uint64
}

func migrateReviewModerationModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.User{},
		&model.Player{},
		&model.Game{},
		&model.ServiceItem{},
		&model.Order{},
		&model.Payment{},
		&model.Review{},
		&model.ReviewReply{},
		&model.ReviewReport{},
		&model.SensitiveWord{},
		&model.OperationLog{},
	); err != nil {
		t.Fatalf("migrate review moderation models: %v", err)
	}
}

func seedReviewModerationData(t *testing.T, db *gorm.DB) reviewModerationSeed {
	t.Helper()
	ctx := context.Background()

	userRepo := user.NewUserRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	gameRepo := game.NewGameRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	serviceItemRepo := serviceitem.NewServiceItemRepository(db)

	// 创建管理员用户
	adminUser := &model.User{
		Name:         "AdminUser",
		Email:        "admin@example.com",
		Phone:        "10000000001",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleAdmin,
	}
	if err := userRepo.Create(ctx, adminUser); err != nil {
		t.Fatalf("seed admin user: %v", err)
	}

	// 创建普通用户
	normalUser := &model.User{
		Name:         "NormalUser",
		Email:        "user@example.com",
		Phone:        "10000000002",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleUser,
	}
	if err := userRepo.Create(ctx, normalUser); err != nil {
		t.Fatalf("seed normal user: %v", err)
	}

	// 创建陪玩师用户
	playerUser := &model.User{
		Name:         "PlayerUser",
		Email:        "player@example.com",
		Phone:        "10000000003",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RolePlayer,
	}
	if err := userRepo.Create(ctx, playerUser); err != nil {
		t.Fatalf("seed player user: %v", err)
	}

	// 创建陪玩师
	player := &model.Player{
		UserID:          playerUser.ID,
		Nickname:        "TestPlayer",
		HourlyRateCents: 5000,
	}
	if err := playerRepo.Create(ctx, player); err != nil {
		t.Fatalf("seed player: %v", err)
	}

	// 创建游戏
	gameModel := &model.Game{
		Key:  "lol",
		Name: "League of Legends",
	}
	if err := gameRepo.Create(ctx, gameModel); err != nil {
		t.Fatalf("seed game: %v", err)
	}

	// 创建服务项目
	playerIDPtr := player.ID
	gameIDPtr := gameModel.ID
	serviceItem := &model.ServiceItem{
		ItemCode:       "TEST001",
		Name:           "Test Service",
		Description:    "Test service description",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		PlayerID:       &playerIDPtr,
		GameID:         &gameIDPtr,
		BasePriceCents: 5000,
		ServiceHours:   1,
		IsActive:       true,
	}
	if err := serviceItemRepo.Create(ctx, serviceItem); err != nil {
		t.Fatalf("seed service item: %v", err)
	}

	// 创建订单
	scheduledStart := time.Now().Add(1 * time.Hour)
	orderModel := &model.Order{
		UserID:          normalUser.ID,
		ItemID:          serviceItem.ID,
		PlayerID:        &playerIDPtr,
		GameID:          &gameIDPtr,
		Title:           "Test Order",
		Status:          model.OrderStatusCompleted,
		Quantity:        1,
		UnitPriceCents:  5000,
		TotalPriceCents: 5000,
		ScheduledStart:  &scheduledStart,
	}
	if err := orderRepo.Create(ctx, orderModel); err != nil {
		t.Fatalf("seed order: %v", err)
	}

	return reviewModerationSeed{
		adminUserID: adminUser.ID,
		userID:      normalUser.ID,
		playerID:    player.ID,
		orderID:     orderModel.ID,
	}
}

func setupReviewModerationRouter(t *testing.T, db *gorm.DB, adminUserID uint64) (*gin.Engine, *adminservice.AdminService) {
	t.Helper()

	// 创建所有需要的仓储
	userRepo := user.NewUserRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	gameRepo := game.NewGameRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	paymentRepo := order.NewPaymentRepository(db)
	roleRepo := role.NewRoleRepository(db)
	serviceItemRepo := serviceitem.NewServiceItemRepository(db)
	permRepo := permission.NewPermissionRepository(db)
	menuRepo := menu.NewMenuRepository(db)
	statsRepo := stats.NewStatsRepository(db)
	walletRepo := wallet.NewWalletRepository(db)
	memCache := cache.NewMemory()

	// 创建管理服务
	adminService := adminservice.NewAdminService(
		gameRepo,
		userRepo,
		playerRepo,
		orderRepo,
		paymentRepo,
		roleRepo,
		serviceItemRepo,
		permRepo,
		menuRepo,
		statsRepo,
		walletRepo,
		memCache,
	)

	// 设置事务管理器
	adminService.SetTxManager(common.NewUnitOfWork(db))

	// 创建路由
	router := gin.New()
	api := router.Group("/api/v1")
	adminAuth := fakeAuthMiddleware(adminUserID)

	// 注册管理员路由
	adminGroup := api.Group("/admin")
	adminGroup.Use(adminAuth)

	reviewHandler := adminhandler.NewReviewHandler(adminService)
	adminGroup.GET("/reviews/pending", reviewHandler.ListPendingReviews)
	adminGroup.PUT("/reviews/:id/approve", reviewHandler.ApproveReview)
	adminGroup.PUT("/reviews/:id/reject", reviewHandler.RejectReview)
	adminGroup.PUT("/reviews/batch-approve", reviewHandler.BatchApproveReviews)
	adminGroup.PUT("/reviews/batch-reject", reviewHandler.BatchRejectReviews)

	// Review report endpoints
	adminGroup.POST("/reviews/:id/reports", reviewHandler.CreateReviewReport)
	adminGroup.GET("/review-reports", reviewHandler.ListReviewReports)
	adminGroup.GET("/review-reports/:id", reviewHandler.GetReviewReport)
	adminGroup.PUT("/review-reports/:id/handle", reviewHandler.HandleReviewReport)

	return router, adminService
}
