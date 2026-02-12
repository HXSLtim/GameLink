package player

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/handler/testutil"
	"gamelink/internal/model"
	gamerepo "gamelink/internal/repository/game"
	gamerankrepo "gamelink/internal/repository/gamerank"
	orderrepo "gamelink/internal/repository/implementations"
	playerrepo "gamelink/internal/repository/player"
	playerschedulerepo "gamelink/internal/repository/playerschedule"
	playerservicerepo "gamelink/internal/repository/playerservice"
	reviewrepo "gamelink/internal/repository/review"
	userrepo "gamelink/internal/repository/user"
	playerservice "gamelink/internal/service/player"
	"gamelink/pkg/cache"
)

type playerExtraTestContext struct {
	Router      *gin.Engine
	DB          *gorm.DB
	User        *model.User
	Player      *model.Player
	Game        *model.Game
	Rank        *model.GameRank
	ServiceMgmt *playerservice.ServiceManagement
	ScheduleSvc *playerservice.ScheduleService
	PlayerSvc   *playerservice.PlayerService
}

func setupPlayerExtraTest(t *testing.T) *playerExtraTestContext {
	t.Helper()

	db := testutil.SetupTestDB(t)
	router := testutil.SetupGinTest(t)

	user := testutil.CreateAdminUser(t, db, model.RolePlayer)
	player := testutil.CreateTestPlayer(t, db, user.ID)
	game := testutil.CreateTestGame(t, db)

	rank := &model.GameRank{
		GameID:      game.ID,
		Name:        "Gold",
		Level:       5,
		PriceCents:  5000,
		IsActive:    true,
		SortOrder:   1,
		Description: "Test rank",
	}
	require.NoError(t, db.Create(rank).Error)

	playerRepo := playerrepo.NewPlayerRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	gameRepo := gamerepo.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	orderRepo := orderrepo.NewOrderRepository(db)
	reviewRepo := reviewrepo.NewReviewRepository(db)
	playerTagRepo := userrepo.NewPlayerTagRepository(db)
	memCache := cache.NewMemory()

	playerSvc := playerservice.NewPlayerService(playerRepo, userRepo, gameRepo, orderRepo, reviewRepo, playerTagRepo, memCache)

	serviceRepo := playerservicerepo.NewPlayerServiceRepository(db)
	serviceMgmt := playerservice.NewServiceManagement(serviceRepo, playerRepo, gameRepo, rankRepo)

	scheduleRepo := playerschedulerepo.NewPlayerScheduleRepository(db)
	scheduleSvc := playerservice.NewScheduleService(scheduleRepo, playerRepo)

	return &playerExtraTestContext{
		Router:      router,
		DB:          db,
		User:        user,
		Player:      player,
		Game:        game,
		Rank:        rank,
		ServiceMgmt: serviceMgmt,
		ScheduleSvc: scheduleSvc,
		PlayerSvc:   playerSvc,
	}
}

func (ctx *playerExtraTestContext) registerRoutes() {
	authMiddleware := func(c *gin.Context) {
		c.Set("user_id", ctx.User.ID)
		c.Next()
	}
	playerGroup := ctx.Router.Group("/player")
	RegisterServiceRoutes(playerGroup, ctx.ServiceMgmt, authMiddleware)
	RegisterScheduleRoutes(playerGroup, ctx.ScheduleSvc, authMiddleware)
	RegisterStatsRoutes(playerGroup, ctx.PlayerSvc, authMiddleware)
}

func TestPlayerService_CRUD(t *testing.T) {
	ctx := setupPlayerExtraTest(t)
	ctx.registerRoutes()

	createPayload := map[string]interface{}{
		"gameId":      ctx.Game.ID,
		"rankId":      ctx.Rank.ID,
		"description": "Test service",
	}
	createResp := testutil.MakeRequest(t, ctx.Router, http.MethodPost, "/player/services", createPayload)
	testutil.AssertSuccess(t, createResp)

	var createBody model.APIResponse[model.PlayerService]
	require.NoError(t, json.Unmarshal(createResp.Body.Bytes(), &createBody))
	serviceID := createBody.Data.ID
	assert.NotZero(t, serviceID)

	listResp := testutil.MakeRequest(t, ctx.Router, http.MethodGet, "/player/services", nil)
	testutil.AssertSuccess(t, listResp)
	var listBody model.APIResponse[[]model.PlayerService]
	require.NoError(t, json.Unmarshal(listResp.Body.Bytes(), &listBody))
	assert.Len(t, listBody.Data, 1)

	updatePayload := map[string]interface{}{
		"description": "Updated service",
	}
	updatePath := testutil.BuildPath("/player/services/:id", map[string]string{
		"id": testutil.Uint64ToStr(serviceID),
	})
	updateResp := testutil.MakeRequest(t, ctx.Router, http.MethodPut, updatePath, updatePayload)
	testutil.AssertSuccess(t, updateResp)

	statusPayload := map[string]interface{}{
		"isActive": false,
	}
	statusPath := testutil.BuildPath("/player/services/:id/status", map[string]string{
		"id": testutil.Uint64ToStr(serviceID),
	})
	statusResp := testutil.MakeRequest(t, ctx.Router, http.MethodPut, statusPath, statusPayload)
	testutil.AssertSuccess(t, statusResp)

	deletePath := testutil.BuildPath("/player/services/:id", map[string]string{
		"id": testutil.Uint64ToStr(serviceID),
	})
	deleteResp := testutil.MakeRequest(t, ctx.Router, http.MethodDelete, deletePath, nil)
	testutil.AssertSuccess(t, deleteResp)
}

func TestPlayerSchedule_GetUpdate(t *testing.T) {
	ctx := setupPlayerExtraTest(t)
	ctx.registerRoutes()

	getResp := testutil.MakeRequest(t, ctx.Router, http.MethodGet, "/player/schedule", nil)
	testutil.AssertSuccess(t, getResp)

	payload := map[string]interface{}{
		"weeklySchedule": map[string]interface{}{
			"monday": map[string]interface{}{
				"enabled": true,
				"start":   "09:00",
				"end":     "18:00",
			},
		},
		"autoOffline":     true,
		"maxOrdersPerDay": 5,
	}
	updateResp := testutil.MakeRequest(t, ctx.Router, http.MethodPut, "/player/schedule", payload)
	testutil.AssertSuccess(t, updateResp)
}

func TestPlayerStats_Endpoints(t *testing.T) {
	ctx := setupPlayerExtraTest(t)
	ctx.registerRoutes()

	orderUser := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	order1 := testutil.CreateTestOrder(t, ctx.DB, orderUser.ID, ctx.Player.ID, ctx.Game.ID, model.OrderStatusCompleted)
	order2 := testutil.CreateTestOrder(t, ctx.DB, orderUser.ID, ctx.Player.ID, ctx.Game.ID, model.OrderStatusCompleted)
	_ = testutil.CreateTestOrder(t, ctx.DB, orderUser.ID, ctx.Player.ID, ctx.Game.ID, model.OrderStatusPending)

	require.NoError(t, ctx.DB.Model(order1).Update("player_income_cents", int64(1200)).Error)
	require.NoError(t, ctx.DB.Model(order2).Update("player_income_cents", int64(800)).Error)

	todayResp := testutil.MakeRequest(t, ctx.Router, http.MethodGet, "/player/stats/today", nil)
	testutil.AssertSuccess(t, todayResp)

	var todayBody model.APIResponse[playerservice.PlayerStatsToday]
	require.NoError(t, json.Unmarshal(todayResp.Body.Bytes(), &todayBody))
	assert.Equal(t, int64(3), todayBody.Data.OrderCount)
	assert.Equal(t, int64(2000), todayBody.Data.EarningsCents)

	overviewResp := testutil.MakeRequest(t, ctx.Router, http.MethodGet, "/player/stats/overview", nil)
	testutil.AssertSuccess(t, overviewResp)

	var overviewBody model.APIResponse[playerservice.PlayerStatsOverview]
	require.NoError(t, json.Unmarshal(overviewResp.Body.Bytes(), &overviewBody))
	assert.Equal(t, int64(3), overviewBody.Data.TotalOrders)
	assert.Equal(t, int64(2000), overviewBody.Data.TotalEarningsCents)
	assert.Equal(t, ctx.Player.RatingAverage, overviewBody.Data.RatingAverage)
}
