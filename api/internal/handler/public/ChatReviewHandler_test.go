package public

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/handler/testutil"
	"gamelink/internal/model"
	chatrepo "gamelink/internal/repository/chat"
	orderrepo "gamelink/internal/repository/order"
	userrepo "gamelink/internal/repository/user"
)

type publicChatReviewContext struct {
	Router     *gin.Engine
	DB         *gorm.DB
	Reviewer   *model.User
	PlayerUser *model.User
	Player     *model.Player
}

func setupPublicChatReviewTest(t *testing.T) *publicChatReviewContext {
	t.Helper()

	db := testutil.SetupTestDB(t)
	router := testutil.SetupGinTest(t)

	reviewer := testutil.CreateAdminUser(t, db, model.RoleUser)
	playerUser := testutil.CreateAdminUser(t, db, model.RolePlayer)
	player := testutil.CreateTestPlayer(t, db, playerUser.ID)

	return &publicChatReviewContext{
		Router:     router,
		DB:         db,
		Reviewer:   reviewer,
		PlayerUser: playerUser,
		Player:     player,
	}
}

func (ctx *publicChatReviewContext) registerRoutes() {
	publicGroup := ctx.Router.Group("/public")
	chatGroupRepo := chatrepo.NewChatGroupRepository(ctx.DB)
	reviewRepo := orderrepo.NewReviewRepository(ctx.DB)
	userRepo := userrepo.NewUserRepository(ctx.DB)

	chatHandler := NewPublicChatHandler(chatGroupRepo)
	chatHandler.RegisterRoutes(publicGroup)

	reviewHandler := NewPublicReviewHandler(reviewRepo, userRepo)
	reviewHandler.RegisterRoutes(publicGroup)
}

func TestPublicChat_ListPublicChannels(t *testing.T) {
	ctx := setupPublicChatReviewTest(t)
	ctx.registerRoutes()

	activeGroup := &model.ChatGroup{
		GroupName: "Public Lobby",
		GroupType: model.ChatGroupTypePublic,
		CreatedBy: ctx.PlayerUser.ID,
		IsActive:  true,
	}
	inactiveGroup := &model.ChatGroup{
		GroupName: "Inactive",
		GroupType: model.ChatGroupTypePublic,
		CreatedBy: ctx.PlayerUser.ID,
		IsActive:  false,
	}
	privateGroup := &model.ChatGroup{
		GroupName: "Private",
		GroupType: model.ChatGroupTypePrivate,
		CreatedBy: ctx.PlayerUser.ID,
		IsActive:  true,
	}
	require.NoError(t, ctx.DB.Create(activeGroup).Error)
	require.NoError(t, ctx.DB.Create(inactiveGroup).Error)
	require.NoError(t, ctx.DB.Create(privateGroup).Error)

	w := testutil.MakeRequest(t, ctx.Router, http.MethodGet, "/public/chat/public-channels", nil)
	testutil.AssertSuccess(t, w)

	var resp model.APIResponse[[]model.ChatGroup]
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, activeGroup.ID, resp.Data[0].ID)
}

func TestPublicReview_ListPlayerReviews(t *testing.T) {
	ctx := setupPublicChatReviewTest(t)
	ctx.registerRoutes()

	review := &model.Review{
		OrderID:     1,
		UserID:      ctx.Reviewer.ID,
		PlayerID:    ctx.Player.ID,
		Score:       model.Rating5,
		Content:     "Great service",
		Status:      model.ReviewStatusApproved,
		IsPublic:    true,
		IsAnonymous: false,
		ExpireAt:    time.Now().Add(7 * 24 * time.Hour),
	}
	require.NoError(t, ctx.DB.Create(review).Error)

	path := testutil.BuildURL("/public/players/:id/reviews", map[string]string{
		"id": testutil.Uint64ToStr(ctx.Player.ID),
	}, map[string]string{"rating": "5"})

	w := testutil.MakeRequest(t, ctx.Router, http.MethodGet, path, nil)
	testutil.AssertSuccess(t, w)

	var resp model.APIResponse[[]map[string]interface{}]
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp.Data, 1)
}
