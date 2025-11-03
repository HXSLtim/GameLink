package service

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionservice "gamelink/internal/service/commission"
	giftservice "gamelink/internal/service/gift"
	itemservice "gamelink/internal/service/item"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupIntegrationTestDB 设置集成测试数据库
func setupIntegrationTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// 迁移所有表
	err = db.AutoMigrate(
		&model.User{},
		&model.Player{},
		&model.Game{},
		&model.ServiceItem{},
		&model.Order{},
		&model.CommissionRule{},
		&model.CommissionRecord{},
		&model.MonthlySettlement{},
	)
	require.NoError(t, err)

	return db
}

// TestGiftOrderFlow 测试完整的礼物赠送流程
func TestGiftOrderFlow(t *testing.T) {
	db := setupIntegrationTestDB(t)
	ctx := context.Background()

	// 初始化repositories
	gameRepo := &mockGameRepo{db: db}
	playerRepo := &mockPlayerRepo{db: db}
	serviceItemRepo := repository.NewServiceItemRepository(db)
	orderRepo := &mockOrderRepo{db: db}
	commissionRepo := repository.NewCommissionRepository(db)

	// 初始化services
	itemSvc := itemservice.NewServiceItemService(serviceItemRepo, gameRepo, playerRepo)
	giftSvc := giftservice.NewGiftService(serviceItemRepo, orderRepo, playerRepo, commissionRepo)

	// Step 1: 创建测试数据
	t.Run("准备测试数据", func(t *testing.T) {
		// 创建游戏
		game := &model.Game{
			Key:  "lol",
			Name: "英雄联盟",
		}
		db.Create(game)

		// 创建陪玩师
		user := &model.User{
			Name:  "测试用户",
			Email: "test@test.com",
			Phone: "13800138000",
		}
		db.Create(user)

		player := &model.Player{
			UserID:          user.ID,
			Nickname:        "测试陪玩师",
			HourlyRateCents: 50000,
		}
		db.Create(player)

		// 创建默认抽成规则
		defaultRule := &model.CommissionRule{
			Name:     "默认抽成",
			Type:     "default",
			Rate:     20,
			IsActive: true,
		}
		db.Create(defaultRule)
	})

	// Step 2: 管理员创建礼物
	var giftID uint64
	t.Run("创建礼物服务项目", func(t *testing.T) {
		req := itemservice.CreateServiceItemRequest{
			ItemCode:       "ESCORT_GIFT_ROSE_PREMIUM",
			Name:           "高端玫瑰",
			Description:    "送给陪玩师表达感谢",
			SubCategory:    model.SubCategoryGift,
			BasePriceCents: 10000,
			ServiceHours:   0, // 礼物为0
			CommissionRate: 0.20,
			MinUsers:       1,
			MaxPlayers:     1,
		}

		item, err := itemSvc.CreateServiceItem(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.NotZero(t, item.ID)
		assert.True(t, item.IsGift())
		assert.Equal(t, 0, item.ServiceHours)

		giftID = item.ID
	})

	// Step 3: 用户赠送礼物
	var orderID uint64
	t.Run("用户赠送礼物", func(t *testing.T) {
		var player model.Player
		db.First(&player)

		req := giftservice.SendGiftRequest{
			PlayerID:    player.ID,
			GiftItemID:  giftID,
			Quantity:    3,
			Message:     "感谢你的陪伴！",
			IsAnonymous: false,
		}

		resp, err := giftSvc.SendGift(ctx, 100, req)

		// 验证礼物订单创建
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, player.ID, resp.PlayerID)
		assert.Equal(t, "高端玫瑰", resp.GiftName)
		assert.Equal(t, 3, resp.Quantity)
		assert.Equal(t, int64(30000), resp.TotalPrice) // 10000 * 3
		assert.NotNil(t, resp.DeliveredAt) // 立即送达

		orderID = resp.OrderID
	})

	// Step 4: 验证订单创建正确
	t.Run("验证礼物订单", func(t *testing.T) {
		var order model.Order
		err := db.First(&order, orderID).Error

		assert.NoError(t, err)
		assert.Equal(t, giftID, order.ItemID)
		assert.Equal(t, 3, order.Quantity)
		assert.Equal(t, int64(10000), order.UnitPriceCents)
		assert.Equal(t, int64(30000), order.TotalPriceCents)
		assert.Equal(t, int64(6000), order.CommissionCents)      // 20%
		assert.Equal(t, int64(24000), order.PlayerIncomeCents)   // 80%
		assert.Equal(t, "感谢你的陪伴！", order.GiftMessage)
		assert.False(t, order.IsAnonymous)
		assert.NotNil(t, order.RecipientPlayerID)
		assert.NotNil(t, order.DeliveredAt)
		assert.Equal(t, model.OrderStatusCompleted, order.Status)
		assert.True(t, order.IsGiftOrder())
	})

	// Step 5: 验证抽成记录
	t.Run("验证抽成记录自动创建", func(t *testing.T) {
		var commissionRecord model.CommissionRecord
		err := db.Where("order_id = ?", orderID).First(&commissionRecord).Error

		assert.NoError(t, err)
		assert.Equal(t, orderID, commissionRecord.OrderID)
		assert.Equal(t, int64(30000), commissionRecord.TotalAmountCents)
		assert.Equal(t, 20, commissionRecord.CommissionRate)
		assert.Equal(t, int64(6000), commissionRecord.CommissionCents)
		assert.Equal(t, int64(24000), commissionRecord.PlayerIncomeCents)
		assert.Equal(t, "pending", commissionRecord.SettlementStatus)
	})

	// Step 6: 陪玩师查看收到的礼物
	t.Run("陪玩师查看收到的礼物", func(t *testing.T) {
		var player model.Player
		db.First(&player)

		resp, err := giftSvc.GetPlayerReceivedGifts(ctx, player.ID, 1, 10)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Gifts, 1)
		assert.Equal(t, "高端玫瑰", resp.Gifts[0].GiftName)
		assert.Equal(t, 3, resp.Gifts[0].Quantity)
		assert.Equal(t, int64(30000), resp.Gifts[0].TotalPrice)
		assert.Equal(t, int64(24000), resp.Gifts[0].Income)
	})

	// Step 7: 月度结算
	t.Run("月度自动结算", func(t *testing.T) {
		var player model.Player
		db.First(&player)

		commissionSvc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

		month := time.Now().Format("2006-01")
		err := commissionSvc.SettleMonth(ctx, month)

		assert.NoError(t, err)

		// 验证月度结算记录
		var settlement model.MonthlySettlement
		err = db.Where("player_id = ? AND settlement_month = ?", player.ID, month).
			First(&settlement).Error

		assert.NoError(t, err)
		assert.Equal(t, int64(1), settlement.TotalOrderCount)
		assert.Equal(t, int64(30000), settlement.TotalAmountCents)
		assert.Equal(t, int64(6000), settlement.TotalCommissionCents)
		assert.Equal(t, int64(24000), settlement.TotalIncomeCents)
	})
}

// Mock implementations for integration tests
type mockGameRepo struct {
	db *gorm.DB
}

func (m *mockGameRepo) Get(ctx context.Context, id uint64) (*model.Game, error) {
	var game model.Game
	err := m.db.WithContext(ctx).First(&game, id).Error
	return &game, err
}

func (m *mockGameRepo) Create(ctx context.Context, game *model.Game) error { return nil }
func (m *mockGameRepo) Update(ctx context.Context, game *model.Game) error { return nil }
func (m *mockGameRepo) Delete(ctx context.Context, id uint64) error { return nil }
func (m *mockGameRepo) List(ctx context.Context) ([]model.Game, error) { return nil, nil }
func (m *mockGameRepo) GetByKey(ctx context.Context, key string) (*model.Game, error) { return nil, nil }

type mockPlayerRepo struct {
	db *gorm.DB
}

func (m *mockPlayerRepo) Get(ctx context.Context, id uint64) (*model.Player, error) {
	var player model.Player
	err := m.db.WithContext(ctx).First(&player, id).Error
	return &player, err
}

func (m *mockPlayerRepo) Create(ctx context.Context, player *model.Player) error { return nil }
func (m *mockPlayerRepo) Update(ctx context.Context, player *model.Player) error { return nil }
func (m *mockPlayerRepo) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	return nil, 0, nil
}

type mockOrderRepo struct {
	db *gorm.DB
}

func (m *mockOrderRepo) Create(ctx context.Context, order *model.Order) error {
	return m.db.WithContext(ctx).Create(order).Error
}

func (m *mockOrderRepo) Get(ctx context.Context, id uint64) (*model.Order, error) {
	var order model.Order
	err := m.db.WithContext(ctx).First(&order, id).Error
	return &order, err
}

func (m *mockOrderRepo) Update(ctx context.Context, order *model.Order) error {
	return m.db.WithContext(ctx).Save(order).Error
}

func (m *mockOrderRepo) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	query := m.db.WithContext(ctx).Model(&model.Order{})

	if opts.PlayerID != nil {
		query = query.Where("player_id = ? OR recipient_player_id = ?", *opts.PlayerID, *opts.PlayerID)
	}

	if len(opts.Statuses) > 0 {
		query = query.Where("status IN ?", opts.Statuses)
	}

	var total int64
	query.Count(&total)

	var orders []model.Order
	query.Offset((opts.Page - 1) * opts.PageSize).Limit(opts.PageSize).Find(&orders)

	return orders, total, nil
}


