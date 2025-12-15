package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	adminrepo "gamelink/internal/repository/admin"
	"gamelink/internal/repository/game"
	"gamelink/internal/repository/menu"
	"gamelink/internal/repository/permission"
	"gamelink/internal/repository/serviceitem"
	"gamelink/internal/repository/stats"
	"gamelink/internal/repository/user"
	"gamelink/internal/repository/wallet"
	adminservice "gamelink/internal/service/admin"
	"gamelink/pkg/cache"
	"gamelink/pkg/testutil"
)

// createAdminServiceForPlayerTest 创建用于陪玩师测试的AdminService
func createAdminServiceForPlayerTest(db *gorm.DB, memCache cache.Cache) *adminservice.AdminService {
	userRepo := user.NewUserRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	gameRepo := game.NewGameRepository(db)
	roleRepo := adminrepo.NewRoleRepository(db)
	serviceItemRepo := serviceitem.NewServiceItemRepository(db)
	permRepo := permission.NewPermissionRepository(db)
	menuRepo := menu.NewMenuRepository(db)
	statsRepo := stats.NewStatsRepository(db)
	walletRepo := wallet.NewWalletRepository(db)
	return adminservice.NewAdminService(gameRepo, userRepo, playerRepo, nil, nil, roleRepo, serviceItemRepo, permRepo, menuRepo, statsRepo, walletRepo, memCache)
}

// setupAdminPlayerTestDB 设置管理端陪玩师测试数据库
func setupAdminPlayerTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.User{},
		&model.Game{},
		&model.Player{},
		&model.RoleModel{},
		&model.Permission{},
		&model.OperationLog{},
	)
	return db
}

// TestAdminPlayerVerification 测试管理端陪玩师审核流程
func TestAdminPlayerVerification(t *testing.T) {
	db := setupAdminPlayerTestDB(t)
	defer testutil.CleanDB(t, db)

	// 创建管理员用户
	admin := &model.User{
		Phone:        "13900000001",
		Email:        "admin@test.com",
		Name:         "Admin",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(admin).Error)

	// 创建陪玩师用户
	playerUser := &model.User{
		Phone:        "13900000002",
		Email:        "player@test.com",
		Name:         "Player",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(playerUser).Error)

	// 创建游戏
	game := &model.Game{Key: "lol", Name: "英雄联盟", Category: "moba"}
	require.NoError(t, db.Create(game).Error)

	// 创建待审核的陪玩师
	player := &model.Player{
		UserID:             playerUser.ID,
		Nickname:           "测试陪玩",
		Bio:                "专业陪玩",
		MainGameID:         game.ID,
		HourlyRateCents:    5000,
		VerificationStatus: model.VerificationPending,
	}
	require.NoError(t, db.Create(player).Error)

	// 创建AdminService
	memCache := cache.NewMemory()
	svc := createAdminServiceForPlayerTest(db, memCache)

	t.Run("审核通过陪玩师", func(t *testing.T) {
		result, err := svc.UpdatePlayerVerification(context.Background(), player.ID, adminservice.UpdateVerificationInput{
			Nickname:           player.Nickname,
			Bio:                player.Bio,
			HourlyRateCents:    player.HourlyRateCents,
			MainGameID:         player.MainGameID,
			VerificationStatus: model.VerificationVerified,
			VerifiedBy:         admin.ID,
			Remark:             "资料完整，审核通过",
		})
		require.NoError(t, err)

		// 验证状态更新
		assert.Equal(t, model.VerificationVerified, result.VerificationStatus)
		assert.NotNil(t, result.VerifiedAt)
		assert.NotNil(t, result.VerifiedBy)
		assert.Equal(t, admin.ID, *result.VerifiedBy)
		assert.Equal(t, "资料完整，审核通过", result.VerifyRemark)
		assert.Empty(t, result.RejectReason) // 通过时不应有拒绝原因
	})

	t.Run("审核拒绝陪玩师", func(t *testing.T) {
		// 创建另一个待审核陪玩师
		playerUser2 := &model.User{
			Phone:        "13900000003",
			Email:        "player2@test.com",
			Name:         "Player2",
			Role:         model.RolePlayer,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(playerUser2).Error)

		player2 := &model.Player{
			UserID:             playerUser2.ID,
			Nickname:           "待审核陪玩",
			MainGameID:         game.ID,
			HourlyRateCents:    3000,
			VerificationStatus: model.VerificationPending,
		}
		require.NoError(t, db.Create(player2).Error)

		result, err := svc.UpdatePlayerVerification(context.Background(), player2.ID, adminservice.UpdateVerificationInput{
			Nickname:           player2.Nickname,
			Bio:                player2.Bio,
			HourlyRateCents:    player2.HourlyRateCents,
			MainGameID:         player2.MainGameID,
			VerificationStatus: model.VerificationRejected,
			VerifiedBy:         admin.ID,
			Remark:             "资料不完整，请补充游戏截图",
		})
		require.NoError(t, err)

		// 验证状态更新
		assert.Equal(t, model.VerificationRejected, result.VerificationStatus)
		assert.NotNil(t, result.VerifiedAt)
		assert.Equal(t, admin.ID, *result.VerifiedBy)
		assert.Equal(t, "资料不完整，请补充游戏截图", result.RejectReason)
		assert.Equal(t, "资料不完整，请补充游戏截图", result.VerifyRemark)
	})

	t.Run("重新审核已拒绝的陪玩师", func(t *testing.T) {
		// 创建已拒绝的陪玩师
		playerUser3 := &model.User{
			Phone:        "13900000004",
			Email:        "player3@test.com",
			Name:         "Player3",
			Role:         model.RolePlayer,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(playerUser3).Error)

		player3 := &model.Player{
			UserID:             playerUser3.ID,
			Nickname:           "已拒绝陪玩",
			MainGameID:         game.ID,
			HourlyRateCents:    4000,
			VerificationStatus: model.VerificationRejected,
			RejectReason:       "之前的拒绝原因",
		}
		require.NoError(t, db.Create(player3).Error)

		// 重新审核通过
		result, err := svc.UpdatePlayerVerification(context.Background(), player3.ID, adminservice.UpdateVerificationInput{
			Nickname:           player3.Nickname,
			Bio:                "补充了资料",
			HourlyRateCents:    player3.HourlyRateCents,
			MainGameID:         player3.MainGameID,
			VerificationStatus: model.VerificationVerified,
			VerifiedBy:         admin.ID,
			Remark:             "资料已补充完整，重新审核通过",
		})
		require.NoError(t, err)

		assert.Equal(t, model.VerificationVerified, result.VerificationStatus)
		assert.Empty(t, result.RejectReason) // 通过后应清空拒绝原因
		assert.Equal(t, "资料已补充完整，重新审核通过", result.VerifyRemark)
	})
}

// TestAdminPlayerList 测试管理端陪玩师列表
func TestAdminPlayerList(t *testing.T) {
	db := setupAdminPlayerTestDB(t)
	defer testutil.CleanDB(t, db)

	// 创建游戏
	game := &model.Game{Key: "lol", Name: "英雄联盟", Category: "moba"}
	require.NoError(t, db.Create(game).Error)

	// 创建多个陪玩师
	statuses := []model.VerificationStatus{
		model.VerificationPending,
		model.VerificationVerified,
		model.VerificationVerified,
		model.VerificationRejected,
	}

	for i, status := range statuses {
		user := &model.User{
			Phone:        "1390000010" + string(rune('0'+i)),
			Email:        "player" + string(rune('0'+i)) + "@test.com",
			Name:         "Player" + string(rune('0'+i)),
			Role:         model.RolePlayer,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(user).Error)

		player := &model.Player{
			UserID:             user.ID,
			Nickname:           "陪玩师" + string(rune('0'+i)),
			MainGameID:         game.ID,
			HourlyRateCents:    int64(5000 + i*1000),
			VerificationStatus: status,
		}
		require.NoError(t, db.Create(player).Error)
	}

	memCache := cache.NewMemory()
	svc := createAdminServiceForPlayerTest(db, memCache)

	t.Run("获取所有陪玩师列表", func(t *testing.T) {
		players, err := svc.ListPlayers(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 4, len(players))
	})

	t.Run("获取单个陪玩师详情", func(t *testing.T) {
		// 获取第一个陪玩师
		var firstPlayer model.Player
		require.NoError(t, db.First(&firstPlayer).Error)

		player, err := svc.GetPlayer(context.Background(), firstPlayer.ID)
		require.NoError(t, err)
		assert.Equal(t, firstPlayer.ID, player.ID)
		assert.Equal(t, firstPlayer.Nickname, player.Nickname)
	})
}
