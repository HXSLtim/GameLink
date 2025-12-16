package review

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	orderimpl "gamelink/internal/repository/implementations"
	"gamelink/internal/repository/notification"
	"gamelink/internal/repository/operationlog"
	reviewrepo "gamelink/internal/repository/review"
	"gamelink/internal/repository/reviewreply"
	"gamelink/internal/repository/reviewreport"
	"gamelink/internal/repository/user"
	"gamelink/pkg/testutil"
)

func setupReviewTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.User{},
		&model.Player{},
		&model.Game{},
		&model.Order{},
		&model.Review{},
		&model.ReviewReply{},
		&model.ReviewReport{},
		&model.NotificationEvent{},
		&model.OperationLog{},
	)
	return db
}

func createReviewTestData(t *testing.T, db *gorm.DB) (customer *model.User, playerUser *model.User, player *model.Player, order *model.Order) {
	t.Helper()

	// 创建用户
	customer = &model.User{
		Phone:        "13800001001",
		Email:        "customer_review@test.com",
		Name:         "Review Customer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	// 创建陪玩师用户
	playerUser = &model.User{
		Phone:        "13800001002",
		Email:        "player_review@test.com",
		Name:         "Review Player",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(playerUser).Error)

	// 创建游戏
	game := &model.Game{Key: "lol", Name: "英雄联盟", Category: "moba"}
	require.NoError(t, db.Create(game).Error)

	// 创建陪玩师
	player = &model.Player{
		UserID:             playerUser.ID,
		Nickname:           "Pro Review Player",
		MainGameID:         game.ID,
		HourlyRateCents:    5000,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(player).Error)

	// 创建已完成订单
	scheduledStart := time.Now().Add(-24 * time.Hour)
	scheduledEnd := scheduledStart.Add(2 * time.Hour)
	completedAt := time.Now().Add(-1 * time.Hour)
	order = &model.Order{
		UserID:          customer.ID,
		ItemID:          game.ID,
		Title:           "测试订单",
		Status:          model.OrderStatusCompleted,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
		ScheduledEnd:    &scheduledEnd,
		CompletedAt:     &completedAt,
	}
	order.SetPlayerID(player.ID)
	order.SetGameID(game.ID)
	require.NoError(t, db.Create(order).Error)

	return
}

func createReviewService(db *gorm.DB) *ReviewService {
	reviewRepo := reviewrepo.NewReviewRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	userRepo := user.NewUserRepository(db)
	replyRepo := reviewreply.NewReviewReplyRepository(db)
	reportRepo := reviewreport.NewReviewReportRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	opLogRepo := operationlog.NewOperationLogRepository(db)

	return NewReviewService(reviewRepo, orderRepo, playerRepo, userRepo, replyRepo, reportRepo, notificationRepo, opLogRepo)
}

func TestReviewService_CreateReview(t *testing.T) {
	db := setupReviewTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, order := createReviewTestData(t, db)
	svc := createReviewService(db)

	t.Run("创建评价成功", func(t *testing.T) {
		resp, err := svc.CreateReview(context.Background(), customer.ID, CreateReviewRequest{
			OrderID: order.ID,
			Rating:  5,
			Comment: "服务非常好！",
		})
		require.NoError(t, err)
		assert.NotZero(t, resp.ReviewID)

		// 验证陪玩师评分已更新
		var updatedPlayer model.Player
		require.NoError(t, db.First(&updatedPlayer, player.ID).Error)
		assert.Equal(t, float32(5), updatedPlayer.RatingAverage)
		assert.Equal(t, uint32(1), updatedPlayer.RatingCount)
	})

	t.Run("重复评价应失败", func(t *testing.T) {
		// 创建新订单
		scheduledStart := time.Now().Add(-24 * time.Hour)
		completedAt := time.Now().Add(-1 * time.Hour)
		newOrder := &model.Order{
			UserID:          customer.ID,
			ItemID:          1,
			Title:           "重复评价测试订单",
			Status:          model.OrderStatusCompleted,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
			CompletedAt:     &completedAt,
		}
		newOrder.SetPlayerID(player.ID)
		require.NoError(t, db.Create(newOrder).Error)

		// 第一次评价
		_, err := svc.CreateReview(context.Background(), customer.ID, CreateReviewRequest{
			OrderID: newOrder.ID,
			Rating:  4,
			Comment: "第一次评价",
		})
		require.NoError(t, err)

		// 第二次评价应失败
		_, err = svc.CreateReview(context.Background(), customer.ID, CreateReviewRequest{
			OrderID: newOrder.ID,
			Rating:  5,
			Comment: "第二次评价",
		})
		assert.ErrorIs(t, err, ErrAlreadyReviewed)
	})

	t.Run("评价未完成订单应失败", func(t *testing.T) {
		// 创建未完成订单
		scheduledStart := time.Now().Add(24 * time.Hour)
		pendingOrder := &model.Order{
			UserID:          customer.ID,
			ItemID:          1,
			Title:           "未完成订单",
			Status:          model.OrderStatusPending,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		pendingOrder.SetPlayerID(player.ID)
		require.NoError(t, db.Create(pendingOrder).Error)

		_, err := svc.CreateReview(context.Background(), customer.ID, CreateReviewRequest{
			OrderID: pendingOrder.ID,
			Rating:  5,
			Comment: "评价未完成订单",
		})
		assert.ErrorIs(t, err, ErrOrderNotCompleted)
	})

	t.Run("评价他人订单应失败", func(t *testing.T) {
		// 创建另一个用户
		otherUser := &model.User{
			Phone:        "13800001003",
			Email:        "other_review@test.com",
			Name:         "Other User",
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(otherUser).Error)

		// 创建另一个用户的订单
		scheduledStart := time.Now().Add(-24 * time.Hour)
		completedAt := time.Now().Add(-1 * time.Hour)
		otherOrder := &model.Order{
			UserID:          otherUser.ID,
			ItemID:          1,
			Title:           "他人订单",
			Status:          model.OrderStatusCompleted,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
			CompletedAt:     &completedAt,
		}
		otherOrder.SetPlayerID(player.ID)
		require.NoError(t, db.Create(otherOrder).Error)

		_, err := svc.CreateReview(context.Background(), customer.ID, CreateReviewRequest{
			OrderID: otherOrder.ID,
			Rating:  5,
			Comment: "评价他人订单",
		})
		assert.ErrorIs(t, err, ErrUnauthorized)
	})
}

func TestReviewService_GetMyReviews(t *testing.T) {
	db := setupReviewTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, playerUser, player, _ := createReviewTestData(t, db)
	svc := createReviewService(db)

	// 创建多个评价
	for i := 0; i < 3; i++ {
		scheduledStart := time.Now().Add(-24 * time.Hour)
		completedAt := time.Now().Add(-1 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          1,
			Title:           "测试订单",
			Status:          model.OrderStatusCompleted,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
			CompletedAt:     &completedAt,
		}
		order.SetPlayerID(player.ID)
		require.NoError(t, db.Create(order).Error)

		review := &model.Review{
			OrderID:  order.ID,
			UserID:   customer.ID,
			PlayerID: player.ID,
			Score:    model.Rating(4 + i%2),
			Content:  "测试评价内容",
			Status:   model.ReviewStatusApproved,
		}
		require.NoError(t, db.Create(review).Error)
	}

	t.Run("获取我的评价列表", func(t *testing.T) {
		resp, err := svc.GetMyReviews(context.Background(), customer.ID, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(3), resp.Total)
		assert.Len(t, resp.Reviews, 3)
	})

	t.Run("分页获取评价", func(t *testing.T) {
		resp, err := svc.GetMyReviews(context.Background(), customer.ID, 1, 2)
		require.NoError(t, err)
		assert.Equal(t, int64(3), resp.Total)
		assert.Len(t, resp.Reviews, 2)
	})

	t.Run("其他用户无评价", func(t *testing.T) {
		resp, err := svc.GetMyReviews(context.Background(), playerUser.ID, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(0), resp.Total)
		assert.Empty(t, resp.Reviews)
	})
}

func TestReviewService_GetPlayerReviews(t *testing.T) {
	db := setupReviewTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, _ := createReviewTestData(t, db)
	svc := createReviewService(db)

	// 创建多个评价
	for i := 0; i < 5; i++ {
		scheduledStart := time.Now().Add(-24 * time.Hour)
		completedAt := time.Now().Add(-1 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          1,
			Title:           "测试订单",
			Status:          model.OrderStatusCompleted,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
			CompletedAt:     &completedAt,
		}
		order.SetPlayerID(player.ID)
		require.NoError(t, db.Create(order).Error)

		review := &model.Review{
			OrderID:  order.ID,
			UserID:   customer.ID,
			PlayerID: player.ID,
			Score:    model.Rating(3 + i%3),
			Content:  "陪玩师评价内容",
			Status:   model.ReviewStatusApproved,
		}
		require.NoError(t, db.Create(review).Error)
	}

	t.Run("获取陪玩师评价列表", func(t *testing.T) {
		reviews, total, err := svc.GetPlayerReviews(context.Background(), player.ID, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, reviews, 5)
	})

	t.Run("分页获取陪玩师评价", func(t *testing.T) {
		reviews, total, err := svc.GetPlayerReviews(context.Background(), player.ID, 1, 3)
		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, reviews, 3)
	})
}

func TestReviewService_ReplyReview(t *testing.T) {
	db := setupReviewTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, playerUser, player, order := createReviewTestData(t, db)
	svc := createReviewService(db)

	// 创建评价
	review := &model.Review{
		OrderID:  order.ID,
		UserID:   customer.ID,
		PlayerID: player.ID,
		Score:    5,
		Content:  "很好的服务",
		Status:   model.ReviewStatusApproved,
	}
	require.NoError(t, db.Create(review).Error)

	t.Run("陪玩师回复评价成功", func(t *testing.T) {
		resp, err := svc.ReplyReview(context.Background(), playerUser.ID, review.ID, ReplyReviewRequest{
			Content: "感谢您的好评！",
		})
		require.NoError(t, err)
		assert.NotZero(t, resp.ReplyID)
	})

	t.Run("非陪玩师回复应失败", func(t *testing.T) {
		_, err := svc.ReplyReview(context.Background(), customer.ID, review.ID, ReplyReviewRequest{
			Content: "用户尝试回复",
		})
		assert.Error(t, err)
	})

	t.Run("回复内容过长应失败", func(t *testing.T) {
		longContent := make([]byte, 600)
		for i := range longContent {
			longContent[i] = 'a'
		}
		_, err := svc.ReplyReview(context.Background(), playerUser.ID, review.ID, ReplyReviewRequest{
			Content: string(longContent),
		})
		assert.Error(t, err)
	})
}

func TestReviewService_ApproveReview(t *testing.T) {
	db := setupReviewTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, order := createReviewTestData(t, db)
	svc := createReviewService(db)

	t.Run("批准待审核评价", func(t *testing.T) {
		review := &model.Review{
			OrderID:  order.ID,
			UserID:   customer.ID,
			PlayerID: player.ID,
			Score:    5,
			Content:  "待审核评价",
			Status:   model.ReviewStatusPending,
		}
		require.NoError(t, db.Create(review).Error)

		adminID := uint64(1)
		err := svc.ApproveReview(context.Background(), review.ID, &adminID)
		require.NoError(t, err)

		// 验证状态
		var updated model.Review
		require.NoError(t, db.First(&updated, review.ID).Error)
		assert.Equal(t, model.ReviewStatusApproved, updated.Status)
	})

	t.Run("批准已通过评价应失败", func(t *testing.T) {
		review := &model.Review{
			OrderID:  order.ID,
			UserID:   customer.ID,
			PlayerID: player.ID,
			Score:    4,
			Content:  "已通过评价",
			Status:   model.ReviewStatusApproved,
		}
		require.NoError(t, db.Create(review).Error)

		adminID := uint64(1)
		err := svc.ApproveReview(context.Background(), review.ID, &adminID)
		assert.Error(t, err)
	})
}

func TestReviewService_RejectReview(t *testing.T) {
	db := setupReviewTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, order := createReviewTestData(t, db)
	svc := createReviewService(db)

	t.Run("拒绝待审核评价", func(t *testing.T) {
		review := &model.Review{
			OrderID:  order.ID,
			UserID:   customer.ID,
			PlayerID: player.ID,
			Score:    1,
			Content:  "包含敏感内容的评价",
			Status:   model.ReviewStatusPending,
		}
		require.NoError(t, db.Create(review).Error)

		adminID := uint64(1)
		err := svc.RejectReview(context.Background(), review.ID, "包含敏感词", &adminID)
		require.NoError(t, err)

		// 验证状态
		var updated model.Review
		require.NoError(t, db.First(&updated, review.ID).Error)
		assert.Equal(t, model.ReviewStatusRejected, updated.Status)
		assert.Equal(t, "包含敏感词", updated.RejectionReason)
	})

	t.Run("拒绝原因为空应失败", func(t *testing.T) {
		review := &model.Review{
			OrderID:  order.ID,
			UserID:   customer.ID,
			PlayerID: player.ID,
			Score:    2,
			Content:  "另一个待审核评价",
			Status:   model.ReviewStatusPending,
		}
		require.NoError(t, db.Create(review).Error)

		adminID := uint64(1)
		err := svc.RejectReview(context.Background(), review.ID, "", &adminID)
		assert.Error(t, err)
	})
}

func TestReviewService_BatchApprove(t *testing.T) {
	db := setupReviewTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, order := createReviewTestData(t, db)
	svc := createReviewService(db)

	t.Run("批量批准评价", func(t *testing.T) {
		var reviewIDs []uint64
		for i := 0; i < 3; i++ {
			review := &model.Review{
				OrderID:  order.ID,
				UserID:   customer.ID,
				PlayerID: player.ID,
				Score:    model.Rating(3 + i),
				Content:  "批量审核评价",
				Status:   model.ReviewStatusPending,
			}
			require.NoError(t, db.Create(review).Error)
			reviewIDs = append(reviewIDs, review.ID)
		}

		adminID := uint64(1)
		err := svc.BatchApprove(context.Background(), reviewIDs, &adminID)
		require.NoError(t, err)

		// 验证所有评价状态
		for _, id := range reviewIDs {
			var review model.Review
			require.NoError(t, db.First(&review, id).Error)
			assert.Equal(t, model.ReviewStatusApproved, review.Status)
		}
	})

	t.Run("空列表应失败", func(t *testing.T) {
		adminID := uint64(1)
		err := svc.BatchApprove(context.Background(), []uint64{}, &adminID)
		assert.Error(t, err)
	})
}

func TestReviewService_BatchReject(t *testing.T) {
	db := setupReviewTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, order := createReviewTestData(t, db)
	svc := createReviewService(db)

	t.Run("批量拒绝评价", func(t *testing.T) {
		var reviewIDs []uint64
		for i := 0; i < 3; i++ {
			review := &model.Review{
				OrderID:  order.ID,
				UserID:   customer.ID,
				PlayerID: player.ID,
				Score:    model.Rating(1 + i),
				Content:  "批量拒绝评价",
				Status:   model.ReviewStatusPending,
			}
			require.NoError(t, db.Create(review).Error)
			reviewIDs = append(reviewIDs, review.ID)
		}

		adminID := uint64(1)
		err := svc.BatchReject(context.Background(), reviewIDs, "批量拒绝原因", &adminID)
		require.NoError(t, err)

		// 验证所有评价状态
		for _, id := range reviewIDs {
			var review model.Review
			require.NoError(t, db.First(&review, id).Error)
			assert.Equal(t, model.ReviewStatusRejected, review.Status)
		}
	})

	t.Run("空原因应失败", func(t *testing.T) {
		review := &model.Review{
			OrderID:  order.ID,
			UserID:   customer.ID,
			PlayerID: player.ID,
			Score:    3,
			Content:  "测试评价",
			Status:   model.ReviewStatusPending,
		}
		require.NoError(t, db.Create(review).Error)

		adminID := uint64(1)
		err := svc.BatchReject(context.Background(), []uint64{review.ID}, "", &adminID)
		assert.Error(t, err)
	})
}

func TestReviewService_ListPendingReviews(t *testing.T) {
	db := setupReviewTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, order := createReviewTestData(t, db)
	svc := createReviewService(db)

	// 创建待审核评价
	for i := 0; i < 5; i++ {
		review := &model.Review{
			OrderID:  order.ID,
			UserID:   customer.ID,
			PlayerID: player.ID,
			Score:    model.Rating(3 + i%3),
			Content:  "待审核评价",
			Status:   model.ReviewStatusPending,
		}
		require.NoError(t, db.Create(review).Error)
	}

	// 创建已通过评价
	for i := 0; i < 3; i++ {
		review := &model.Review{
			OrderID:  order.ID,
			UserID:   customer.ID,
			PlayerID: player.ID,
			Score:    5,
			Content:  "已通过评价",
			Status:   model.ReviewStatusApproved,
		}
		require.NoError(t, db.Create(review).Error)
	}

	t.Run("获取待审核评价列表", func(t *testing.T) {
		reviews, total, err := svc.ListPendingReviews(context.Background(), 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, reviews, 5)
		for _, r := range reviews {
			assert.Equal(t, model.ReviewStatusPending, r.Status)
		}
	})

	t.Run("分页获取待审核评价", func(t *testing.T) {
		reviews, total, err := svc.ListPendingReviews(context.Background(), 1, 3)
		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, reviews, 3)
	})
}

func TestReviewService_UpdateReply(t *testing.T) {
	db := setupReviewTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, playerUser, player, order := createReviewTestData(t, db)
	svc := createReviewService(db)

	// 创建评价
	review := &model.Review{
		OrderID:  order.ID,
		UserID:   customer.ID,
		PlayerID: player.ID,
		Score:    5,
		Content:  "很好的服务",
		Status:   model.ReviewStatusApproved,
	}
	require.NoError(t, db.Create(review).Error)

	// 创建回复
	reply := &model.ReviewReply{
		ReviewID: review.ID,
		AuthorID: playerUser.ID,
		Content:  "原始回复内容",
		Status:   "approved",
	}
	require.NoError(t, db.Create(reply).Error)

	t.Run("更新回复成功", func(t *testing.T) {
		resp, err := svc.UpdateReply(context.Background(), playerUser.ID, reply.ID, UpdateReplyRequest{
			Content: "更新后的回复内容",
		})
		require.NoError(t, err)
		assert.Equal(t, reply.ID, resp.ReplyID)
	})

	t.Run("非作者更新应失败", func(t *testing.T) {
		_, err := svc.UpdateReply(context.Background(), customer.ID, reply.ID, UpdateReplyRequest{
			Content: "用户尝试更新",
		})
		assert.ErrorIs(t, err, ErrUnauthorized)
	})
}

func TestReviewService_DeleteReply(t *testing.T) {
	db := setupReviewTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, playerUser, player, order := createReviewTestData(t, db)
	svc := createReviewService(db)

	// 创建评价
	review := &model.Review{
		OrderID:  order.ID,
		UserID:   customer.ID,
		PlayerID: player.ID,
		Score:    5,
		Content:  "很好的服务",
		Status:   model.ReviewStatusApproved,
	}
	require.NoError(t, db.Create(review).Error)

	t.Run("删除回复成功", func(t *testing.T) {
		// 创建回复
		reply := &model.ReviewReply{
			ReviewID: review.ID,
			AuthorID: playerUser.ID,
			Content:  "待删除回复",
			Status:   "approved",
		}
		require.NoError(t, db.Create(reply).Error)

		err := svc.DeleteReply(context.Background(), playerUser.ID, reply.ID)
		require.NoError(t, err)

		// 验证已删除
		var deleted model.ReviewReply
		err = db.First(&deleted, reply.ID).Error
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("非作者删除应失败", func(t *testing.T) {
		reply := &model.ReviewReply{
			ReviewID: review.ID,
			AuthorID: playerUser.ID,
			Content:  "另一个回复",
			Status:   "approved",
		}
		require.NoError(t, db.Create(reply).Error)

		err := svc.DeleteReply(context.Background(), customer.ID, reply.ID)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})
}

func TestReviewService_GetUsersByIDs(t *testing.T) {
	db := setupReviewTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, playerUser, _, _ := createReviewTestData(t, db)
	svc := createReviewService(db)

	t.Run("批量获取用户", func(t *testing.T) {
		users, err := svc.GetUsersByIDs(context.Background(), []uint64{customer.ID, playerUser.ID})
		require.NoError(t, err)
		assert.Len(t, users, 2)
	})

	t.Run("空列表返回空结果", func(t *testing.T) {
		users, err := svc.GetUsersByIDs(context.Background(), []uint64{})
		require.NoError(t, err)
		assert.Empty(t, users)
	})
}

func TestReviewService_UpdatePlayerRating(t *testing.T) {
	db := setupReviewTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, _ := createReviewTestData(t, db)
	svc := createReviewService(db)

	// 创建多个评价来测试评分计算
	scores := []int{5, 4, 5, 3, 4}
	for i, score := range scores {
		scheduledStart := time.Now().Add(-24 * time.Hour)
		completedAt := time.Now().Add(-1 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          1,
			Title:           "评分测试订单",
			Status:          model.OrderStatusCompleted,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
			CompletedAt:     &completedAt,
		}
		order.SetPlayerID(player.ID)
		require.NoError(t, db.Create(order).Error)

		_, err := svc.CreateReview(context.Background(), customer.ID, CreateReviewRequest{
			OrderID: order.ID,
			Rating:  score,
			Comment: "评分测试" + string(rune('0'+i)),
		})
		require.NoError(t, err)
	}

	// 验证平均评分
	var updatedPlayer model.Player
	require.NoError(t, db.First(&updatedPlayer, player.ID).Error)

	// 计算预期平均分: (5+4+5+3+4)/5 = 4.2
	expectedAvg := float32(21) / float32(5)
	assert.InDelta(t, expectedAvg, updatedPlayer.RatingAverage, 0.01)
	assert.Equal(t, uint32(5), updatedPlayer.RatingCount)
}
