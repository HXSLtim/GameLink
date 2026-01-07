//go:build integration

package order

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gamelink/internal/model"
	commissionrepo "gamelink/internal/repository/commission"
	gamerepo "gamelink/internal/repository/game"
	orderrepo "gamelink/internal/repository/implementations"
	ordermodelsrepo "gamelink/internal/repository/order"
	ordergrouprepo "gamelink/internal/repository/ordergroup"
	userrepo "gamelink/internal/repository/user"
)

// testDB 测试数据库连接
var testDB *gorm.DB

// TestMain 设置测试环境
func TestMain(m *testing.M) {
	// 从环境变量获取测试数据库配置
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		host := getEnvOrDefault("TEST_DB_HOST", "localhost")
		port := getEnvOrDefault("TEST_DB_PORT", "5433")
		user := getEnvOrDefault("TEST_DB_USER", "gamelink")
		password := getEnvOrDefault("TEST_DB_PASSWORD", "gamelink123")
		dbname := getEnvOrDefault("TEST_DB_NAME", "gamelink_test")
		dsn = "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	}

	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}

	// 自动迁移测试表
	err = testDB.AutoMigrate(
		&model.User{},
		&model.Player{},
		&model.Game{},
		&model.ServiceItem{},
		&model.Order{},
		&model.OrderGroup{},
		&model.Payment{},
		&model.Review{},
		&model.CommissionRule{},
		&model.CommissionRecord{},
	)
	if err != nil {
		panic("failed to migrate test database: " + err.Error())
	}

	// 运行测试
	code := m.Run()

	// 清理
	sqlDB, _ := testDB.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}

	os.Exit(code)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// setupTestData 设置测试数据
func setupTestData(t *testing.T) (userID, playerID, gameID, serviceItemID uint64, cleanup func()) {
	ctx := context.Background()

	// 使用纳秒级时间戳确保唯一性
	uniqueSuffix := time.Now().UnixNano()

	// 创建测试用户
	user := &model.User{
		Name:  "Test User",
		Email: fmt.Sprintf("test%d@example.com", uniqueSuffix),
		Phone: fmt.Sprintf("138%08d", uniqueSuffix%100000000),
	}
	require.NoError(t, testDB.WithContext(ctx).Create(user).Error)

	// 创建测试游戏
	game := &model.Game{
		Name:     "Test Game",
		IsActive: true,
	}
	require.NoError(t, testDB.WithContext(ctx).Create(game).Error)

	// 创建测试服务项目（解决 item_id 外键约束）
	serviceItem := &model.ServiceItem{
		ItemCode:       fmt.Sprintf("TEST_%d", uniqueSuffix),
		Name:           "Test Escort Service",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		GameID:         &game.ID,
		BasePriceCents: 5000,
		ServiceHours:   1,
		CommissionRate: 0.20,
		IsActive:       true,
	}
	require.NoError(t, testDB.WithContext(ctx).Create(serviceItem).Error)

	// 创建测试陪玩师
	player := &model.Player{
		UserID:             user.ID,
		Nickname:           "Test Player",
		HourlyRateCents:    5000,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, testDB.WithContext(ctx).Create(player).Error)

	cleanup = func() {
		testDB.WithContext(ctx).Unscoped().Delete(&model.Order{}, "user_id = ?", user.ID)
		testDB.WithContext(ctx).Unscoped().Delete(&model.OrderGroup{}, "user_id = ?", user.ID)
		testDB.WithContext(ctx).Unscoped().Delete(&model.Player{}, "id = ?", player.ID)
		testDB.WithContext(ctx).Unscoped().Delete(&model.ServiceItem{}, "id = ?", serviceItem.ID)
		testDB.WithContext(ctx).Unscoped().Delete(&model.Game{}, "id = ?", game.ID)
		testDB.WithContext(ctx).Unscoped().Delete(&model.User{}, "id = ?", user.ID)
	}

	return user.ID, player.ID, game.ID, serviceItem.ID, cleanup
}

// createTestService 创建测试服务
func createTestService(t *testing.T) (*OrderService, ordergrouprepo.Repository) {
	orderRepo := orderrepo.NewOrderRepository(testDB)
	playerRepo := userrepo.NewPlayerRepository(testDB)
	userRepo := userrepo.NewUserRepository(testDB)
	gameRepo := gamerepo.NewGameRepository(testDB)
	paymentRepo := ordermodelsrepo.NewPaymentRepository(testDB)
	reviewRepo := ordermodelsrepo.NewReviewRepository(testDB)
	commissionRepo := commissionrepo.NewCommissionRepository(testDB)
	orderGroupRepo := ordergrouprepo.NewRepository(testDB)

	svc := NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	svc.SetOrderGroupRepository(orderGroupRepo)

	return svc, orderGroupRepo
}

// TestIntegration_CreateOrderWithSplit 测试订单拆分创建完整流程
func TestIntegration_CreateOrderWithSplit(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not available")
	}

	userID, playerID, gameID, serviceItemID, cleanup := setupTestData(t)
	defer cleanup()

	svc, groupRepo := createTestService(t)
	ctx := context.Background()

	// 创建 3 小时订单（应该拆分成 3 个子订单）
	scheduledStart := time.Now().Add(time.Hour)
	req := CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		ServiceID:      &serviceItemID,
		Title:          "Integration Test Order",
		Description:    "Testing order split",
		ScheduledStart: &scheduledStart,
		DurationHours:  3,
	}

	resp, err := svc.CreateOrder(ctx, userID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.IsSplit)
	assert.Equal(t, 3, resp.SubOrderCount)
	assert.Equal(t, 3, resp.TotalHours)
	assert.NotEmpty(t, resp.GroupNo)

	// 验证主订单
	group, err := groupRepo.GetWithSubOrders(ctx, resp.OrderID)
	require.NoError(t, err)
	assert.Equal(t, userID, group.UserID)
	assert.Equal(t, gameID, group.GameID)
	assert.Equal(t, 3, group.TotalHours)
	assert.Equal(t, model.OrderGroupStatusPending, group.Status)

	// 验证子订单
	assert.Len(t, group.SubOrders, 3)
	for i, sub := range group.SubOrders {
		assert.Equal(t, i+1, sub.HourIndex)
		assert.True(t, sub.IsSubOrder)
		assert.True(t, sub.CanTransfer)
		assert.Equal(t, &group.ID, sub.GroupID)
	}
}

// TestIntegration_CreateOrderNoSplit 测试短时长订单不拆分
func TestIntegration_CreateOrderNoSplit(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not available")
	}

	userID, playerID, gameID, serviceItemID, cleanup := setupTestData(t)
	defer cleanup()

	svc, _ := createTestService(t)
	ctx := context.Background()

	// 创建 1 小时订单（不应该拆分）
	scheduledStart := time.Now().Add(time.Hour)
	req := CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		ServiceID:      &serviceItemID,
		Title:          "Short Order",
		ScheduledStart: &scheduledStart,
		DurationHours:  1,
	}

	resp, err := svc.CreateOrder(ctx, userID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.IsSplit)
	assert.Equal(t, 0, resp.SubOrderCount)
}

// TestIntegration_TransferSubOrder 测试转单完整流程
func TestIntegration_TransferSubOrder(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not available")
	}

	userID, playerID, gameID, serviceItemID, cleanup := setupTestData(t)
	defer cleanup()

	ctx := context.Background()

	// 创建第二个陪玩师用于转单
	newPlayer := &model.Player{
		UserID:             userID + 1000, // 使用不同的 UserID
		Nickname:           "New Player",
		HourlyRateCents:    6000,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, testDB.WithContext(ctx).Create(newPlayer).Error)
	defer testDB.WithContext(ctx).Delete(newPlayer)

	svc, groupRepo := createTestService(t)

	// 创建 3 小时订单
	scheduledStart := time.Now().Add(time.Hour)
	createResp, err := svc.CreateOrder(ctx, userID, CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		ServiceID:      &serviceItemID,
		Title:          "Transfer Test Order",
		ScheduledStart: &scheduledStart,
		DurationHours:  3,
	})
	require.NoError(t, err)

	// 获取第二个子订单进行转单
	group, err := groupRepo.GetWithSubOrders(ctx, createResp.OrderID)
	require.NoError(t, err)
	require.Len(t, group.SubOrders, 3)

	subOrderID := group.SubOrders[1].ID // 第二个小时的订单

	// 执行转单
	transferResp, err := svc.TransferSubOrder(ctx, userID, TransferSubOrderRequest{
		SubOrderID:   subOrderID,
		NewPlayerID:  newPlayer.ID,
		TransferNote: "Integration test transfer",
	})

	require.NoError(t, err)
	assert.True(t, transferResp.Success)
	assert.NotZero(t, transferResp.NewSubOrderID)

	// 验证原订单状态
	var originalOrder model.Order
	require.NoError(t, testDB.WithContext(ctx).First(&originalOrder, subOrderID).Error)
	assert.Equal(t, model.OrderStatusCanceled, originalOrder.Status)
	assert.False(t, originalOrder.CanTransfer)
	assert.NotNil(t, originalOrder.TransferTo)

	// 验证新订单
	var newOrder model.Order
	require.NoError(t, testDB.WithContext(ctx).First(&newOrder, transferResp.NewSubOrderID).Error)
	assert.Equal(t, newPlayer.ID, *newOrder.PlayerID)
	assert.Equal(t, originalOrder.HourIndex, newOrder.HourIndex)
	assert.True(t, newOrder.IsSubOrder)
	assert.Equal(t, &subOrderID, newOrder.TransferFrom)
}

// TestIntegration_BatchTransferSubOrders 测试批量转单
func TestIntegration_BatchTransferSubOrders(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not available")
	}

	userID, playerID, gameID, serviceItemID, cleanup := setupTestData(t)
	defer cleanup()

	ctx := context.Background()

	// 创建第二个陪玩师
	newPlayer := &model.Player{
		UserID:             userID + 2000,
		Nickname:           "Batch Transfer Player",
		HourlyRateCents:    6000,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, testDB.WithContext(ctx).Create(newPlayer).Error)
	defer testDB.WithContext(ctx).Delete(newPlayer)

	svc, groupRepo := createTestService(t)

	// 创建 4 小时订单
	scheduledStart := time.Now().Add(time.Hour)
	createResp, err := svc.CreateOrder(ctx, userID, CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		ServiceID:      &serviceItemID,
		Title:          "Batch Transfer Test",
		ScheduledStart: &scheduledStart,
		DurationHours:  4,
	})
	require.NoError(t, err)

	// 获取子订单
	group, err := groupRepo.GetWithSubOrders(ctx, createResp.OrderID)
	require.NoError(t, err)
	require.Len(t, group.SubOrders, 4)

	// 批量转单第 2、3、4 小时
	subOrderIDs := []uint64{
		group.SubOrders[1].ID,
		group.SubOrders[2].ID,
		group.SubOrders[3].ID,
	}

	batchResp, err := svc.BatchTransferSubOrders(ctx, userID, BatchTransferRequest{
		SubOrderIDs:  subOrderIDs,
		NewPlayerID:  newPlayer.ID,
		TransferNote: "Batch transfer test",
	})

	require.NoError(t, err)
	assert.Equal(t, 3, batchResp.SuccessCount)
	assert.Equal(t, 0, batchResp.FailedCount)
	assert.Len(t, batchResp.NewOrderIDs, 3)
}

// TestIntegration_GetTransferableSubOrders 测试获取可转单子订单
func TestIntegration_GetTransferableSubOrders(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not available")
	}

	userID, playerID, gameID, serviceItemID, cleanup := setupTestData(t)
	defer cleanup()

	ctx := context.Background()
	svc, groupRepo := createTestService(t)

	// 创建 3 小时订单
	scheduledStart := time.Now().Add(time.Hour)
	createResp, err := svc.CreateOrder(ctx, userID, CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		ServiceID:      &serviceItemID,
		Title:          "Transferable Test",
		ScheduledStart: &scheduledStart,
		DurationHours:  3,
	})
	require.NoError(t, err)

	// 将第一个子订单标记为已完成（不可转）
	group, err := groupRepo.GetWithSubOrders(ctx, createResp.OrderID)
	require.NoError(t, err)

	firstSubOrder := group.SubOrders[0]
	firstSubOrder.Status = model.OrderStatusCompleted
	firstSubOrder.CanTransfer = false
	require.NoError(t, testDB.WithContext(ctx).Save(&firstSubOrder).Error)

	// 获取可转单子订单
	transferable, err := svc.GetTransferableSubOrders(ctx, createResp.OrderID)

	require.NoError(t, err)
	assert.Len(t, transferable, 2) // 只有第 2、3 小时可转
}

// TestIntegration_OrderGroupStatusUpdate 测试主订单状态更新
func TestIntegration_OrderGroupStatusUpdate(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not available")
	}

	userID, playerID, gameID, serviceItemID, cleanup := setupTestData(t)
	defer cleanup()

	ctx := context.Background()
	svc, groupRepo := createTestService(t)

	// 创建 3 小时订单
	scheduledStart := time.Now().Add(time.Hour)
	createResp, err := svc.CreateOrder(ctx, userID, CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		ServiceID:      &serviceItemID,
		Title:          "Status Update Test",
		ScheduledStart: &scheduledStart,
		DurationHours:  3,
	})
	require.NoError(t, err)

	// 获取主订单
	group, err := groupRepo.GetWithSubOrders(ctx, createResp.OrderID)
	require.NoError(t, err)

	// 将所有子订单标记为已完成
	for i := range group.SubOrders {
		group.SubOrders[i].Status = model.OrderStatusCompleted
		require.NoError(t, testDB.WithContext(ctx).Save(&group.SubOrders[i]).Error)
	}

	// 更新主订单状态
	group.UpdateStatusFromSubOrders(group.SubOrders)
	require.NoError(t, groupRepo.Update(ctx, group))

	// 验证状态
	updatedGroup, err := groupRepo.Get(ctx, createResp.OrderID)
	require.NoError(t, err)
	assert.Equal(t, model.OrderGroupStatusCompleted, updatedGroup.Status)
	assert.Equal(t, 3, updatedGroup.CompletedHours)
}
