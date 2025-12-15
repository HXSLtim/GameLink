package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	playerhandler "gamelink/internal/handler/player"
	"gamelink/internal/model"
	gamerepo "gamelink/internal/repository/game"
	orderrepo "gamelink/internal/repository/implementations"
	playertagrepo "gamelink/internal/repository/playertag"
	reviewrepo "gamelink/internal/repository/review"
	playerrepo "gamelink/internal/repository/user"
	playerservice "gamelink/internal/service/player"
	"gamelink/pkg/cache"
	"gamelink/pkg/testutil"
)

// setupPlayerTestDB 设置陪玩师测试数据库
func setupPlayerTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.User{},
		&model.Game{},
		&model.Player{},
		&model.Order{},
		&model.Review{},
		&model.ServiceItem{},
	)
	return db
}

// setupPlayerService 创建陪玩师服务
func setupPlayerService(db *gorm.DB) *playerservice.PlayerService {
	memCache := cache.NewMemory()
	playerRepo := playerrepo.NewPlayerRepository(db)
	userRepo := playerrepo.NewUserRepository(db)
	gameRepo := gamerepo.NewGameRepository(db)
	orderRepo := orderrepo.NewOrderRepository(db)
	reviewRepo := reviewrepo.NewReviewRepository(db)
	playerTagRepo := playertagrepo.NewPlayerTagRepository(db)

	return playerservice.NewPlayerService(playerRepo, userRepo, gameRepo, orderRepo, reviewRepo, playerTagRepo, memCache)
}

// TestPlayerApply 测试申请成为陪玩师
func TestPlayerApply(t *testing.T) {
	db := setupPlayerTestDB(t)
	defer testutil.CleanDB(t, db)

	// 创建测试用户
	user := &model.User{
		Phone:        "13800000001",
		Email:        "test@example.com",
		Name:         "Test User",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(user).Error)

	// 创建测试游戏
	game := &model.Game{
		Key:      "lol",
		Name:     "英雄联盟",
		Category: "moba",
	}
	require.NoError(t, db.Create(game).Error)

	svc := setupPlayerService(db)

	// 设置路由
	gin.SetMode(gin.TestMode)
	router := gin.New()
	group := router.Group("/player")
	playerhandler.RegisterProfileRoutes(group, svc, fakeAuthMiddleware(user.ID))

	t.Run("成功申请成为陪玩师", func(t *testing.T) {
		req := playerservice.ApplyPlayerRequest{
			Nickname:        "峡谷守护者",
			Bio:             "专业陪玩，技术过硬",
			MainGameID:      game.ID,
			Rank:            "钻石",
			HourlyRateCents: 5000,
			Tags:            []string{"技术流", "耐心"},
		}

		w := doJSON(router, "POST", "/player/apply", req, "")
		assert.Equal(t, http.StatusOK, w.Code)

		var resp apiResp[playerservice.ApplyPlayerResponse]
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.True(t, resp.Success)
		assert.Equal(t, model.VerificationPending, resp.Data.VerificationStatus)
		assert.Greater(t, resp.Data.PlayerID, uint64(0))
	})

	t.Run("重复申请应失败", func(t *testing.T) {
		req := playerservice.ApplyPlayerRequest{
			Nickname:        "另一个昵称",
			Bio:             "再次申请",
			MainGameID:      game.ID,
			Rank:            "王者",
			HourlyRateCents: 8000,
		}

		w := doJSON(router, "POST", "/player/apply", req, "")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("无效游戏ID应失败", func(t *testing.T) {
		// 创建新用户
		user2 := &model.User{
			Phone:        "13800000002",
			Email:        "test2@example.com",
			Name:         "Test User 2",
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(user2).Error)

		router2 := gin.New()
		group2 := router2.Group("/player")
		playerhandler.RegisterProfileRoutes(group2, svc, fakeAuthMiddleware(user2.ID))

		req := playerservice.ApplyPlayerRequest{
			Nickname:        "测试昵称",
			Bio:             "测试简介",
			MainGameID:      99999, // 不存在的游戏ID
			Rank:            "钻石",
			HourlyRateCents: 5000,
		}

		w := doJSON(router2, "POST", "/player/apply", req, "")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

// TestPlayerProfile 测试陪玩师资料管理
func TestPlayerProfile(t *testing.T) {
	db := setupPlayerTestDB(t)
	defer testutil.CleanDB(t, db)

	// 创建测试用户
	user := &model.User{
		Phone:        "13800000010",
		Email:        "player@example.com",
		Name:         "Player User",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(user).Error)

	// 创建测试游戏
	game := &model.Game{
		Key:      "valorant",
		Name:     "无畏契约",
		Category: "fps",
	}
	require.NoError(t, db.Create(game).Error)

	// 创建陪玩师资料
	player := &model.Player{
		UserID:             user.ID,
		Nickname:           "王牌射手",
		Bio:                "FPS专业选手",
		Rank:               "不朽",
		RatingAverage:      4.8,
		RatingCount:        100,
		HourlyRateCents:    8000,
		MainGameID:         game.ID,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(player).Error)

	svc := setupPlayerService(db)

	// 设置路由
	gin.SetMode(gin.TestMode)
	router := gin.New()
	group := router.Group("/player")
	playerhandler.RegisterProfileRoutes(group, svc, fakeAuthMiddleware(user.ID))

	t.Run("获取陪玩师资料", func(t *testing.T) {
		w := doJSON(router, "GET", "/player/profile", nil, "")
		assert.Equal(t, http.StatusOK, w.Code)

		var resp apiResp[playerservice.PlayerDetailResponse]
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.True(t, resp.Success)
		assert.Equal(t, "王牌射手", resp.Data.Player.Nickname)
		assert.Equal(t, "无畏契约", resp.Data.Player.MainGame)
	})

	t.Run("更新陪玩师资料", func(t *testing.T) {
		req := playerservice.UpdatePlayerProfileRequest{
			Nickname:        "超级射手",
			Bio:             "更新后的简介",
			Rank:            "辐射",
			HourlyRateCents: 10000,
			Tags:            []string{"枪法精准", "战术大师"},
		}

		w := doJSON(router, "PUT", "/player/profile", req, "")
		assert.Equal(t, http.StatusOK, w.Code)

		// 验证更新
		var updatedPlayer model.Player
		require.NoError(t, db.First(&updatedPlayer, player.ID).Error)
		assert.Equal(t, "超级射手", updatedPlayer.Nickname)
		assert.Equal(t, int64(10000), updatedPlayer.HourlyRateCents)
	})

	t.Run("设置在线状态", func(t *testing.T) {
		req := playerservice.SetPlayerStatusRequest{
			Online: true,
		}

		w := doJSON(router, "PUT", "/player/status", req, "")
		assert.Equal(t, http.StatusOK, w.Code)

		// 设置离线
		req.Online = false
		w = doJSON(router, "PUT", "/player/status", req, "")
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestPlayerList 测试陪玩师列表查询
func TestPlayerList(t *testing.T) {
	db := setupPlayerTestDB(t)
	defer testutil.CleanDB(t, db)

	// 创建测试游戏
	game1 := &model.Game{Key: "lol", Name: "英雄联盟", Category: "moba"}
	game2 := &model.Game{Key: "csgo", Name: "CS:GO", Category: "fps"}
	require.NoError(t, db.Create(game1).Error)
	require.NoError(t, db.Create(game2).Error)

	// 创建多个陪玩师
	players := []struct {
		phone    string
		email    string
		nickname string
		gameID   uint64
		rate     int64
		rating   float32
		status   model.VerificationStatus
	}{
		{"13800000101", "p1@test.com", "陪玩师A", game1.ID, 5000, 4.9, model.VerificationVerified},
		{"13800000102", "p2@test.com", "陪玩师B", game1.ID, 8000, 4.5, model.VerificationVerified},
		{"13800000103", "p3@test.com", "陪玩师C", game2.ID, 6000, 4.7, model.VerificationVerified},
		{"13800000104", "p4@test.com", "待审核", game1.ID, 5000, 0, model.VerificationPending}, // 未审核
	}

	for _, p := range players {
		user := &model.User{
			Phone:        p.phone,
			Email:        p.email,
			Name:         p.nickname,
			Role:         model.RolePlayer,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(user).Error)

		player := &model.Player{
			UserID:             user.ID,
			Nickname:           p.nickname,
			MainGameID:         p.gameID,
			HourlyRateCents:    p.rate,
			RatingAverage:      p.rating,
			VerificationStatus: p.status,
		}
		require.NoError(t, db.Create(player).Error)
	}

	svc := setupPlayerService(db)

	t.Run("获取所有已审核陪玩师", func(t *testing.T) {
		resp, err := svc.ListPlayers(context.Background(), playerservice.PlayerListRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		// 只返回已审核的陪玩师（3个）
		assert.Equal(t, 3, len(resp.Players))
	})

	t.Run("按游戏筛选", func(t *testing.T) {
		resp, err := svc.ListPlayers(context.Background(), playerservice.PlayerListRequest{
			GameID:   &game1.ID,
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		// LOL游戏只有2个已审核陪玩师
		assert.Equal(t, 2, len(resp.Players))
	})

	t.Run("按价格筛选", func(t *testing.T) {
		minPrice := int64(6000)
		resp, err := svc.ListPlayers(context.Background(), playerservice.PlayerListRequest{
			MinPrice: &minPrice,
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		// 价格>=6000的有2个
		assert.Equal(t, 2, len(resp.Players))
	})

	t.Run("按评分筛选", func(t *testing.T) {
		minRating := float32(4.6)
		resp, err := svc.ListPlayers(context.Background(), playerservice.PlayerListRequest{
			MinRating: &minRating,
			Page:      1,
			PageSize:  10,
		})
		require.NoError(t, err)
		// 评分>=4.6的有2个
		assert.Equal(t, 2, len(resp.Players))
	})
}

// TestPlayerDetail 测试陪玩师详情
func TestPlayerDetail(t *testing.T) {
	db := setupPlayerTestDB(t)
	defer testutil.CleanDB(t, db)

	// 创建测试数据
	user := &model.User{
		Phone:        "13800000200",
		Email:        "detail@test.com",
		Name:         "Detail Test",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
		AvatarURL:    "https://example.com/avatar.jpg",
	}
	require.NoError(t, db.Create(user).Error)

	game := &model.Game{Key: "dota2", Name: "DOTA 2", Category: "moba"}
	require.NoError(t, db.Create(game).Error)

	player := &model.Player{
		UserID:             user.ID,
		Nickname:           "DOTA大神",
		Bio:                "专业DOTA2陪玩",
		Rank:               "冠绝",
		RatingAverage:      4.9,
		RatingCount:        200,
		HourlyRateCents:    12000,
		MainGameID:         game.ID,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(player).Error)

	svc := setupPlayerService(db)

	t.Run("获取已审核陪玩师详情", func(t *testing.T) {
		resp, err := svc.GetPlayerDetail(context.Background(), player.ID)
		require.NoError(t, err)
		assert.Equal(t, "DOTA大神", resp.Player.Nickname)
		assert.Equal(t, "DOTA 2", resp.Player.MainGame)
		assert.Equal(t, float32(4.9), resp.Player.RatingAverage)
		assert.Equal(t, "https://example.com/avatar.jpg", resp.Player.AvatarURL)
	})

	t.Run("获取未审核陪玩师详情应失败", func(t *testing.T) {
		// 创建未审核陪玩师
		user2 := &model.User{
			Phone:        "13800000201",
			Email:        "pending@test.com",
			Name:         "Pending",
			Role:         model.RolePlayer,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(user2).Error)

		pendingPlayer := &model.Player{
			UserID:             user2.ID,
			Nickname:           "待审核",
			MainGameID:         game.ID,
			VerificationStatus: model.VerificationPending,
		}
		require.NoError(t, db.Create(pendingPlayer).Error)

		_, err := svc.GetPlayerDetail(context.Background(), pendingPlayer.ID)
		assert.ErrorIs(t, err, playerservice.ErrPlayerNotVerified)
	})
}

// TestPlayerOrderAccept 测试陪玩师接单流程
func TestPlayerOrderAccept(t *testing.T) {
	db := setupPlayerTestDB(t)
	defer testutil.CleanDB(t, db)

	// 创建用户
	customer := &model.User{
		Phone:        "13800000300",
		Email:        "customer@test.com",
		Name:         "Customer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	playerUser := &model.User{
		Phone:        "13800000301",
		Email:        "player@test.com",
		Name:         "Player",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(playerUser).Error)

	// 创建游戏和服务项目
	game := &model.Game{Key: "lol", Name: "英雄联盟", Category: "moba"}
	require.NoError(t, db.Create(game).Error)

	serviceItem := &model.ServiceItem{
		ItemCode:       "escort-lol",
		Name:           "LOL陪玩",
		Category:       "escort",
		BasePriceCents: 5000,
		IsActive:       true,
	}
	require.NoError(t, db.Create(serviceItem).Error)

	// 创建陪玩师
	player := &model.Player{
		UserID:             playerUser.ID,
		Nickname:           "测试陪玩",
		MainGameID:         game.ID,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(player).Error)

	// 创建待接单订单
	order := &model.Order{
		OrderNo:         "TEST001",
		UserID:          customer.ID,
		ItemID:          serviceItem.ID,
		GameID:          &game.ID,
		Status:          model.OrderStatusPending,
		TotalPriceCents: 5000,
		Currency:        model.CurrencyCNY,
		OrderConfig:     "{}",
	}
	require.NoError(t, db.Create(order).Error)

	t.Run("订单状态为pending", func(t *testing.T) {
		var o model.Order
		require.NoError(t, db.First(&o, order.ID).Error)
		assert.Equal(t, model.OrderStatusPending, o.Status)
	})

	t.Run("陪玩师资料已创建", func(t *testing.T) {
		var p model.Player
		require.NoError(t, db.First(&p, player.ID).Error)
		assert.Equal(t, "测试陪玩", p.Nickname)
		assert.Equal(t, model.VerificationVerified, p.VerificationStatus)
	})
}
