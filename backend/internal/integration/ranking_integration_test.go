package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	orderimpl "gamelink/internal/repository/implementations"
	"gamelink/internal/repository/ranking"
	rankingsvc "gamelink/internal/service/ranking"
	"gamelink/pkg/testutil"
)

func setupRankingTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.User{},
		&model.Player{},
		&model.Game{},
		&model.Order{},
		&model.PlayerRanking{},
		&model.RankingReward{},
	)
	return db
}

func TestRankingCalculation(t *testing.T) {
	db := setupRankingTestDB(t)
	defer testutil.CleanDB(t, db)

	// 创建游戏
	gameModel := &model.Game{Key: "lol", Name: "英雄联盟", Category: "moba"}
	require.NoError(t, db.Create(gameModel).Error)

	// 创建多个陪玩师
	players := make([]*model.Player, 3)
	for i := 0; i < 3; i++ {
		playerUser := &model.User{
			Phone:        "1380000010" + string(rune('0'+i)),
			Email:        "player" + string(rune('0'+i)) + "@test.com",
			Name:         "Player " + string(rune('0'+i)),
			Role:         model.RolePlayer,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(playerUser).Error)

		players[i] = &model.Player{
			UserID:             playerUser.ID,
			Nickname:           "Pro " + string(rune('0'+i)),
			MainGameID:         gameModel.ID,
			HourlyRateCents:    int64(5000 + i*1000),
			VerificationStatus: model.VerificationVerified,
		}
		require.NoError(t, db.Create(players[i]).Error)
	}

	// 创建用户
	customer := &model.User{
		Phone:        "13900000001",
		Email:        "customer@test.com",
		Name:         "Customer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	// 创建已完成订单（不同数量给不同陪玩师）
	now := time.Now()
	month := now.Format("2006-01")

	// Player 0: 5 orders, 50000 cents
	for i := 0; i < 5; i++ {
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "Order",
			Status:          model.OrderStatusCompleted,
			UnitPriceCents:  10000,
			TotalPriceCents: 10000,
			CompletedAt:     &now,
		}
		order.SetPlayerID(players[0].ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)
	}

	// Player 1: 3 orders, 45000 cents
	for i := 0; i < 3; i++ {
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "Order",
			Status:          model.OrderStatusCompleted,
			UnitPriceCents:  15000,
			TotalPriceCents: 15000,
			CompletedAt:     &now,
		}
		order.SetPlayerID(players[1].ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)
	}

	// Player 2: 1 order, 20000 cents
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "Order",
		Status:          model.OrderStatusCompleted,
		UnitPriceCents:  20000,
		TotalPriceCents: 20000,
		CompletedAt:     &now,
	}
	order.SetPlayerID(players[2].ID)
	order.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(order).Error)

	// 创建服务
	rankingRepo := ranking.NewRankingRepository(db)
	commissionRepo := ranking.NewRankingCommissionRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	svc := rankingsvc.NewRankingService(rankingRepo, commissionRepo, orderRepo)

	t.Run("计算月度排名", func(t *testing.T) {
		err := svc.CalculateMonthlyRankings(context.Background(), month)
		require.NoError(t, err)

		// 验证排名已创建
		var rankings []model.PlayerRanking
		require.NoError(t, db.Where("period_value = ?", month).Find(&rankings).Error)
		assert.GreaterOrEqual(t, len(rankings), 3) // 至少3个排名（单量+金额）
	})

	t.Run("获取陪玩师排名信息", func(t *testing.T) {
		info, err := svc.GetPlayerRankingInfo(context.Background(), players[0].ID, month)
		require.NoError(t, err)
		assert.Equal(t, players[0].ID, info.PlayerID)
		assert.Equal(t, month, info.Month)
	})
}

func TestRankingRewardCreation(t *testing.T) {
	db := setupRankingTestDB(t)
	defer testutil.CleanDB(t, db)

	rankingRepo := ranking.NewRankingRepository(db)
	commissionRepo := ranking.NewRankingCommissionRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	svc := rankingsvc.NewRankingService(rankingRepo, commissionRepo, orderRepo)

	t.Run("创建排名奖励规则", func(t *testing.T) {
		reward, err := svc.CreateRankingReward(context.Background(), rankingsvc.CreateRankingRewardRequest{
			RankingType: model.RankingTypeOrderCount,
			Period:      "monthly",
			RankStart:   1,
			RankEnd:     3,
			RewardType:  "commission",
			RewardValue: 500, // 5元奖励
			Description: "月度单量前三名奖励",
		})
		require.NoError(t, err)
		assert.NotZero(t, reward.ID)
		assert.Equal(t, model.RankingTypeOrderCount, reward.RankingType)
		assert.Equal(t, int64(500), reward.RewardValue)
	})
}

func TestGiftOrderExcludedFromRanking(t *testing.T) {
	db := setupRankingTestDB(t)
	defer testutil.CleanDB(t, db)

	// 创建游戏
	gameModel := &model.Game{Key: "lol", Name: "英雄联盟", Category: "moba"}
	require.NoError(t, db.Create(gameModel).Error)

	// 创建陪玩师
	playerUser := &model.User{
		Phone:        "13800000001",
		Email:        "player@test.com",
		Name:         "Player",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(playerUser).Error)

	player := &model.Player{
		UserID:             playerUser.ID,
		Nickname:           "Pro Player",
		MainGameID:         gameModel.ID,
		HourlyRateCents:    5000,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(player).Error)

	// 创建用户
	customer := &model.User{
		Phone:        "13900000001",
		Email:        "customer@test.com",
		Name:         "Customer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	now := time.Now()
	month := now.Format("2006-01")

	// 创建普通订单
	normalOrder := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "Normal Order",
		Status:          model.OrderStatusCompleted,
		UnitPriceCents:  10000,
		TotalPriceCents: 10000,
		CompletedAt:     &now,
	}
	normalOrder.SetPlayerID(player.ID)
	normalOrder.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(normalOrder).Error)

	// 创建礼物订单（应被排除）
	recipientID := player.ID
	giftOrder := &model.Order{
		UserID:            customer.ID,
		ItemID:            gameModel.ID,
		Title:             "Gift Order",
		Status:            model.OrderStatusCompleted,
		UnitPriceCents:    5000,
		TotalPriceCents:   5000,
		CompletedAt:       &now,
		RecipientPlayerID: &recipientID, // 设置接收者表示礼物订单
	}
	giftOrder.SetPlayerID(player.ID)
	giftOrder.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(giftOrder).Error)

	// 创建服务并计算排名
	rankingRepo := ranking.NewRankingRepository(db)
	commissionRepo := ranking.NewRankingCommissionRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	svc := rankingsvc.NewRankingService(rankingRepo, commissionRepo, orderRepo)

	err := svc.CalculateMonthlyRankings(context.Background(), month)
	require.NoError(t, err)

	// 验证排名只计算了普通订单
	var rankings []model.PlayerRanking
	require.NoError(t, db.Where("player_id = ? AND period_value = ? AND ranking_type = ?",
		player.ID, month, model.RankingTypeOrderCount).Find(&rankings).Error)

	if len(rankings) > 0 {
		// 订单数应该是1（只有普通订单）
		assert.Equal(t, int64(1), rankings[0].OrderCount)
	}
}

// 测试排名奖励应用
func TestRankingRewardApplication(t *testing.T) {
	db := setupRankingTestDB(t)
	defer testutil.CleanDB(t, db)

	rankingRepo := ranking.NewRankingRepository(db)
	commissionRepo := ranking.NewRankingCommissionRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	svc := rankingsvc.NewRankingService(rankingRepo, commissionRepo, orderRepo)

	// 创建奖励规则：前3名获得奖励
	_, err := svc.CreateRankingReward(context.Background(), rankingsvc.CreateRankingRewardRequest{
		RankingType: model.RankingTypeOrderCount,
		Period:      "monthly",
		RankStart:   1,
		RankEnd:     1,
		RewardType:  "commission",
		RewardValue: 1000, // 第1名奖励10元
		Description: "月度单量冠军奖励",
	})
	require.NoError(t, err)

	_, err = svc.CreateRankingReward(context.Background(), rankingsvc.CreateRankingRewardRequest{
		RankingType: model.RankingTypeOrderCount,
		Period:      "monthly",
		RankStart:   2,
		RankEnd:     3,
		RewardType:  "commission",
		RewardValue: 500, // 第2-3名奖励5元
		Description: "月度单量亚军季军奖励",
	})
	require.NoError(t, err)

	// 创建游戏和陪玩师
	gameModel := &model.Game{Key: "apex", Name: "Apex Legends", Category: "fps"}
	require.NoError(t, db.Create(gameModel).Error)

	players := make([]*model.Player, 3)
	for i := range 3 {
		playerUser := &model.User{
			Phone:        "1390000100" + string(rune('0'+i)),
			Email:        "reward_player" + string(rune('0'+i)) + "@test.com",
			Name:         "RewardPlayer " + string(rune('0'+i)),
			Role:         model.RolePlayer,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(playerUser).Error)

		players[i] = &model.Player{
			UserID:             playerUser.ID,
			Nickname:           "RewardPro " + string(rune('0'+i)),
			MainGameID:         gameModel.ID,
			HourlyRateCents:    5000,
			VerificationStatus: model.VerificationVerified,
		}
		require.NoError(t, db.Create(players[i]).Error)
	}

	customer := &model.User{
		Phone:        "13900000099",
		Email:        "reward_customer@test.com",
		Name:         "RewardCustomer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	now := time.Now()
	month := now.Format("2006-01")

	// 创建订单：Player0=10单, Player1=5单, Player2=3单
	orderCounts := []int{10, 5, 3}
	for i, count := range orderCounts {
		for j := 0; j < count; j++ {
			order := &model.Order{
				UserID:          customer.ID,
				ItemID:          gameModel.ID,
				Title:           "Reward Order",
				Status:          model.OrderStatusCompleted,
				UnitPriceCents:  5000,
				TotalPriceCents: 5000,
				CompletedAt:     &now,
			}
			order.SetPlayerID(players[i].ID)
			order.SetGameID(gameModel.ID)
			require.NoError(t, db.Create(order).Error)
		}
	}

	// 计算排名
	err = svc.CalculateMonthlyRankings(context.Background(), month)
	require.NoError(t, err)

	// 验证排名和奖励
	var rankings []model.PlayerRanking
	require.NoError(t, db.Where("period_value = ? AND ranking_type = ?",
		month, model.RankingTypeOrderCount).Order("rank ASC").Find(&rankings).Error)

	require.GreaterOrEqual(t, len(rankings), 3)

	// 第1名应该是Player0，奖励1000
	assert.Equal(t, players[0].ID, rankings[0].PlayerID)
	assert.Equal(t, 1, rankings[0].Rank)
	assert.Equal(t, int64(1000), rankings[0].BonusCents)

	// 第2名应该是Player1，奖励500
	assert.Equal(t, players[1].ID, rankings[1].PlayerID)
	assert.Equal(t, 2, rankings[1].Rank)
	assert.Equal(t, int64(500), rankings[1].BonusCents)

	// 第3名应该是Player2，奖励500
	assert.Equal(t, players[2].ID, rankings[2].PlayerID)
	assert.Equal(t, 3, rankings[2].Rank)
	assert.Equal(t, int64(500), rankings[2].BonusCents)
}

// 测试收入排名计算
func TestIncomeRankingCalculation(t *testing.T) {
	db := setupRankingTestDB(t)
	defer testutil.CleanDB(t, db)

	gameModel := &model.Game{Key: "dota2", Name: "Dota 2", Category: "moba"}
	require.NoError(t, db.Create(gameModel).Error)

	// 创建3个陪玩师
	players := make([]*model.Player, 3)
	for i := range 3 {
		playerUser := &model.User{
			Phone:        "1380000200" + string(rune('0'+i)),
			Email:        "income_player" + string(rune('0'+i)) + "@test.com",
			Name:         "IncomePlayer " + string(rune('0'+i)),
			Role:         model.RolePlayer,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(playerUser).Error)

		players[i] = &model.Player{
			UserID:             playerUser.ID,
			Nickname:           "IncomePro " + string(rune('0'+i)),
			MainGameID:         gameModel.ID,
			HourlyRateCents:    int64(5000 + i*1000),
			VerificationStatus: model.VerificationVerified,
		}
		require.NoError(t, db.Create(players[i]).Error)
	}

	customer := &model.User{
		Phone:        "13900000199",
		Email:        "income_customer@test.com",
		Name:         "IncomeCustomer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	now := time.Now()
	month := now.Format("2006-01")

	// 创建订单：Player0=50000分, Player1=80000分, Player2=30000分
	// 收入排名应该是 Player1 > Player0 > Player2
	incomes := []int64{50000, 80000, 30000}
	for i, income := range incomes {
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "Income Order",
			Status:          model.OrderStatusCompleted,
			UnitPriceCents:  income,
			TotalPriceCents: income,
			CompletedAt:     &now,
		}
		order.SetPlayerID(players[i].ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)
	}

	rankingRepo := ranking.NewRankingRepository(db)
	commissionRepo := ranking.NewRankingCommissionRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	svc := rankingsvc.NewRankingService(rankingRepo, commissionRepo, orderRepo)

	err := svc.CalculateMonthlyRankings(context.Background(), month)
	require.NoError(t, err)

	// 验证收入排名
	var rankings []model.PlayerRanking
	require.NoError(t, db.Where("period_value = ? AND ranking_type = ?",
		month, model.RankingTypeIncome).Order("rank ASC").Find(&rankings).Error)

	require.GreaterOrEqual(t, len(rankings), 3)

	// 第1名应该是Player1（80000分）
	assert.Equal(t, players[1].ID, rankings[0].PlayerID)
	assert.Equal(t, 1, rankings[0].Rank)
	assert.Equal(t, int64(80000), rankings[0].IncomeCents)

	// 第2名应该是Player0（50000分）
	assert.Equal(t, players[0].ID, rankings[1].PlayerID)
	assert.Equal(t, 2, rankings[1].Rank)

	// 第3名应该是Player2（30000分）
	assert.Equal(t, players[2].ID, rankings[2].PlayerID)
	assert.Equal(t, 3, rankings[2].Rank)
}

// 测试空订单月份的排名计算
func TestRankingEmptyMonth(t *testing.T) {
	db := setupRankingTestDB(t)
	defer testutil.CleanDB(t, db)

	rankingRepo := ranking.NewRankingRepository(db)
	commissionRepo := ranking.NewRankingCommissionRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	svc := rankingsvc.NewRankingService(rankingRepo, commissionRepo, orderRepo)

	// 计算一个没有订单的月份
	err := svc.CalculateMonthlyRankings(context.Background(), "2020-01")
	require.NoError(t, err)

	// 验证没有排名记录
	var rankings []model.PlayerRanking
	require.NoError(t, db.Where("period_value = ?", "2020-01").Find(&rankings).Error)
	assert.Empty(t, rankings)
}

// 测试获取不存在的陪玩师排名信息
func TestGetNonExistentPlayerRanking(t *testing.T) {
	db := setupRankingTestDB(t)
	defer testutil.CleanDB(t, db)

	rankingRepo := ranking.NewRankingRepository(db)
	commissionRepo := ranking.NewRankingCommissionRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	svc := rankingsvc.NewRankingService(rankingRepo, commissionRepo, orderRepo)

	month := time.Now().Format("2006-01")

	// 获取不存在的陪玩师排名
	info, err := svc.GetPlayerRankingInfo(context.Background(), 99999, month)
	require.NoError(t, err)
	assert.Equal(t, uint64(99999), info.PlayerID)
	assert.Equal(t, month, info.Month)
	assert.Equal(t, 0, info.BestRank) // 没有排名
}
