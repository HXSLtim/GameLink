package user

import (
        "context"

        "gamelink/internal/model"
        "gamelink/internal/repository"
        "gamelink/internal/repository/common"
        assignmentservice "gamelink/internal/service/assignment"
)

type stubDisputeRepo struct{}

func (stubDisputeRepo) Create(context.Context, *model.OrderDispute) error { return nil }
func (stubDisputeRepo) Update(context.Context, *model.OrderDispute) error { return nil }
func (stubDisputeRepo) ListByOrder(context.Context, uint64) ([]model.OrderDispute, error) { return nil, nil }
func (stubDisputeRepo) GetLatestByOrder(context.Context, uint64) (*model.OrderDispute, error) {
        return nil, repository.ErrNotFound
}

type stubOpLogRepo struct{}

func (stubOpLogRepo) Append(context.Context, *model.OperationLog) error { return nil }
func (stubOpLogRepo) ListByEntity(context.Context, string, uint64, repository.OperationLogListOptions) ([]model.OperationLog, int64, error) {
        return nil, 0, nil
}

type stubNotificationRepo struct{}

func (stubNotificationRepo) ListByUser(context.Context, repository.NotificationListOptions) ([]model.NotificationEvent, int64, error) {
        return nil, 0, nil
}
func (stubNotificationRepo) MarkRead(context.Context, uint64, []uint64) error { return nil }
func (stubNotificationRepo) CountUnread(context.Context, uint64) (int64, error) { return 0, nil }
func (stubNotificationRepo) Create(context.Context, *model.NotificationEvent) error { return nil }

type stubTxManager struct {
        orders   repository.OrderRepository
        players  repository.PlayerRepository
        disputes repository.OrderDisputeRepository
        opLogs   repository.OperationLogRepository
}

func (s stubTxManager) WithTx(ctx context.Context, fn func(r *common.Repos) error) error {
        repos := &common.Repos{
                Orders:   s.orders,
                Players:  s.players,
                Disputes: s.disputes,
                OpLogs:   s.opLogs,
        }
        return fn(repos)
}

func newAssignmentServiceStub(orders repository.OrderRepository, players repository.PlayerRepository) *assignmentservice.Service {
        disputes := stubDisputeRepo{}
        opLogs := stubOpLogRepo{}
        notifications := stubNotificationRepo{}
        svc := assignmentservice.NewService(orders, players, disputes, opLogs, notifications)
        svc.SetTxManager(stubTxManager{orders: orders, players: players, disputes: disputes, opLogs: opLogs})
        return svc
}
