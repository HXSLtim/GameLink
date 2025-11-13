package assignment_test

import (
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/logging"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/common"
	notificationrepo "gamelink/internal/repository/notification"
	operationlogrepo "gamelink/internal/repository/operation_log"
	orderrepo "gamelink/internal/repository/order"
	orderdisputerepo "gamelink/internal/repository/orderdispute"
	playerrepo "gamelink/internal/repository/player"
	assignmentservice "gamelink/internal/service/assignment"
)

func setupAssignmentTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&model.User{},
		&model.Player{},
		&model.Order{},
		&model.OrderDispute{},
		&model.OperationLog{},
		&model.NotificationEvent{},
	))
	return db
}

func TestAssignmentService_DisputeLifecycle(t *testing.T) {
	db := setupAssignmentTestDB(t)
	ctx := logging.WithRequestID(context.Background(), "trace-dispute-flow")

	user := &model.User{Name: "customer", Email: "c@example.com", Phone: "13800000001"}
	require.NoError(t, db.Create(user).Error)
	agent := &model.User{Name: "agent", Email: "agent@example.com", Phone: "13800000002"}
	require.NoError(t, db.Create(agent).Error)
	playerUser := &model.User{Name: "player", Email: "p@example.com", Phone: "13800000003"}
	require.NoError(t, db.Create(playerUser).Error)
	player := &model.Player{UserID: playerUser.ID, Nickname: "pro", HourlyRateCents: 1000}
	require.NoError(t, db.Create(player).Error)

	gameID := uint64(1)
	order := &model.Order{
		UserID:           user.ID,
		ItemID:           1,
		GameID:           &gameID,
		Title:            "Boost",
		Status:           model.OrderStatusConfirmed,
		TotalPriceCents:  2000,
		AssignmentSource: model.OrderAssignmentSourceUnknown,
	}
	require.NoError(t, db.Create(order).Error)

	orderRepo := orderrepo.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	disputeRepo := orderdisputerepo.NewRepository(db)
	opLogRepo := operationlogrepo.NewOperationLogRepository(db)
	notificationRepo := notificationrepo.NewNotificationRepository(db)

	svc := assignmentservice.NewService(orderRepo, playerRepo, disputeRepo, opLogRepo, notificationRepo)
	svc.SetTxManager(common.NewUnitOfWork(db))

	raisedBy := user.ID
	dispute, err := svc.CreateDispute(ctx, order.ID, assignmentservice.DisputeRequest{
		RaisedBy:       model.OrderDisputeRaisedByUser,
		RaisedByUserID: &raisedBy,
		Reason:         "服务不达标",
		EvidenceURLs:   []string{"https://example.com/proof.png"},
		TraceID:        "trace-dispute-flow",
	})
	require.NoError(t, err)
	require.Equal(t, model.OrderDisputeStatusPending, dispute.Status)

	_, err = svc.Assign(ctx, order.ID, assignmentservice.AssignInput{
		PlayerID:    player.ID,
		Source:      model.OrderAssignmentSourceManual,
		ActorUserID: &agent.ID,
		TraceID:     "trace-dispute-flow",
	})
	require.NoError(t, err)

	refund := int64(1500)
	resolved, err := svc.MediateDispute(ctx, order.ID, assignmentservice.MediateInput{
		Resolution:        model.OrderDisputeResolutionRefund,
		RefundAmountCents: &refund,
		ActorUserID:       &agent.ID,
		TraceID:           "trace-dispute-flow",
	})
	require.NoError(t, err)
	require.Equal(t, model.OrderDisputeStatusResolved, resolved.Status)
	require.Equal(t, refund, resolved.RefundAmountCents)
	require.WithinDuration(t, time.Now(), *resolved.HandledAt, time.Second)

	updatedOrder, err := orderRepo.Get(ctx, order.ID)
	require.NoError(t, err)
	require.Equal(t, model.OrderStatusRefunded, updatedOrder.Status)
	require.Equal(t, refund, updatedOrder.RefundAmountCents)

	logs, _, err := opLogRepo.ListByEntity(ctx, string(model.OpEntityOrder), order.ID, repository.OperationLogListOptions{})
	require.NoError(t, err)
	foundTrace := false
	for _, l := range logs {
		if l.TraceID == "trace-dispute-flow" && l.Action == string(model.OpActionMediateDispute) {
			foundTrace = true
			break
		}
	}
	require.True(t, foundTrace)

	events, _, err := notificationRepo.ListByUser(ctx, repository.NotificationListOptions{UserID: user.ID})
	require.NoError(t, err)
	require.NotEmpty(t, events)
}
