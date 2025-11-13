package assignment

import (
        "context"
        "encoding/json"
        "errors"
        "time"

        "gamelink/internal/logging"
        "gamelink/internal/model"
        "gamelink/internal/repository"
        "gamelink/internal/repository/common"
)

var (
        // ErrValidation 输入校验失败。
        ErrValidation = errors.New("validation failed")
        // ErrNotFound 标识资源不存在。
        ErrNotFound = repository.ErrNotFound
)

// TxManager 复用通用事务接口。
type TxManager interface {
        WithTx(ctx context.Context, fn func(r *common.Repos) error) error
}

// RecommendationClient 通过外部风控/推荐系统获取候选陪玩师。
type RecommendationClient interface {
        Recommend(ctx context.Context, order model.Order) ([]CandidateRecommendation, error)
}

// CandidateRecommendation 表示推荐接口返回的数据。
type CandidateRecommendation struct {
        PlayerID uint64
        Score    float64
        Reason   string
}

// Service 聚合客服指派与争议处理逻辑。
type Service struct {
        orders        repository.OrderRepository
        players       repository.PlayerRepository
        disputes      repository.OrderDisputeRepository
        opLogs        repository.OperationLogRepository
        notifications repository.NotificationRepository
        tx            TxManager
        recommender   RecommendationClient
        slaDuration   time.Duration
}

// NewService 构造 Service。
func NewService(
        orders repository.OrderRepository,
        players repository.PlayerRepository,
        disputes repository.OrderDisputeRepository,
        opLogs repository.OperationLogRepository,
        notifications repository.NotificationRepository,
) *Service {
        return &Service{
                orders:        orders,
                players:       players,
                disputes:      disputes,
                opLogs:        opLogs,
                notifications: notifications,
                slaDuration:   30 * time.Minute,
        }
}

// SetTxManager 注入事务管理器。
func (s *Service) SetTxManager(tx TxManager) { s.tx = tx }

// SetRecommendationClient 注入推荐服务。
func (s *Service) SetRecommendationClient(c RecommendationClient) { s.recommender = c }

// PendingAssignment 表示待处理指派单。
type PendingAssignment struct {
        OrderID        uint64                 `json:"orderId"`
        UserID         uint64                 `json:"userId"`
        Status         model.OrderStatus      `json:"status"`
        AssignmentSource model.OrderAssignmentSource `json:"assignmentSource"`
        CreatedAt      time.Time              `json:"createdAt"`
        SLADeadline    time.Time              `json:"slaDeadline"`
        SLARemaining   int64                  `json:"slaRemainingSeconds"`
        IsOverdue      bool                   `json:"isOverdue"`
}

// Candidate 用于前端展示候选陪玩师。
type Candidate struct {
        PlayerID   uint64  `json:"playerId"`
        Nickname   string  `json:"nickname"`
        HourlyRate int64   `json:"hourlyRateCents"`
        Score      float64 `json:"score"`
        Source     string  `json:"source"`
        Reason     string  `json:"reason,omitempty"`
}

// AssignInput 指派请求。
type AssignInput struct {
        PlayerID    uint64
        Source      model.OrderAssignmentSource
        ActorUserID *uint64
        TraceID     string
}

// CancelAssignInput 指派回退请求。
type CancelAssignInput struct {
        Reason      string
        ActorUserID *uint64
        TraceID     string
}

// DisputeRequest 用户/陪玩师发起争议。
type DisputeRequest struct {
        RaisedBy       model.OrderDisputeRaisedBy
        RaisedByUserID *uint64
        Reason         string
        EvidenceURLs   []string
        TraceID        string
}

// MediateInput 客服调解请求。
type MediateInput struct {
        Resolution        model.OrderDisputeResolution
        Note              string
        RefundAmountCents *int64
        ReassignPlayerID  *uint64
        ActorUserID       *uint64
        TraceID           string
}

// ListPendingAssignments 返回未指派或待处理的订单。
func (s *Service) ListPendingAssignments(ctx context.Context, page, pageSize int) ([]PendingAssignment, int64, error) {
        opts := repository.OrderListOptions{
                Page:     page,
                PageSize: pageSize,
                Statuses: []model.OrderStatus{model.OrderStatusPending, model.OrderStatusConfirmed},
        }
        orders, _, err := s.orders.List(ctx, opts)
        if err != nil {
                return nil, 0, err
        }
        now := time.Now()
        result := make([]PendingAssignment, 0, len(orders))
        var total int64
        for _, order := range orders {
                if order.PlayerID != nil && *order.PlayerID != 0 {
                        continue
                }
                deadline := order.CreatedAt.Add(s.slaDuration)
                remaining := int64(deadline.Sub(now).Seconds())
                if remaining < 0 {
                        remaining = 0
                }
                result = append(result, PendingAssignment{
                        OrderID:         order.ID,
                        UserID:          order.UserID,
                        Status:          order.Status,
                        AssignmentSource: order.AssignmentSource,
                        CreatedAt:       order.CreatedAt,
                        SLADeadline:     deadline,
                        SLARemaining:    remaining,
                        IsOverdue:       now.After(deadline),
                })
                total++
        }
        return result, total, nil
}

// ListCandidates 获取指派候选人。
func (s *Service) ListCandidates(ctx context.Context, orderID uint64, limit int) ([]Candidate, error) {
        if limit <= 0 {
                        limit = 10
        }
        order, err := s.orders.Get(ctx, orderID)
        if err != nil {
                return nil, err
        }
        pageSize := limit
        players, _, err := s.players.ListPaged(ctx, 1, pageSize)
        if err != nil {
                return nil, err
        }
        candidates := make([]Candidate, 0, len(players))
        for _, player := range players {
                candidates = append(candidates, Candidate{
                        PlayerID:   player.ID,
                        Nickname:   player.Nickname,
                        HourlyRate: player.HourlyRateCents,
                        Score:      0,
                        Source:     "roster",
                })
        }
        if s.recommender != nil {
                recs, err := s.recommender.Recommend(ctx, *order)
                if err == nil {
                        for _, rec := range recs {
                                candidates = append(candidates, Candidate{
                                        PlayerID: rec.PlayerID,
                                        Score:    rec.Score,
                                        Source:   "recommendation",
                                        Reason:   rec.Reason,
                                })
                        }
                }
        }
        return candidates, nil
}

// Assign 指派陪玩师。
func (s *Service) Assign(ctx context.Context, orderID uint64, in AssignInput) (*model.Order, error) {
        if s.tx == nil {
                return nil, errors.New("transaction manager not configured")
        }
        if in.PlayerID == 0 {
                return nil, ErrValidation
        }
        if in.Source == "" {
                in.Source = model.OrderAssignmentSourceManual
        }
        var updated *model.Order
        err := s.tx.WithTx(ctx, func(r *common.Repos) error {
                order, err := r.Orders.Get(ctx, orderID)
                if err != nil {
                        return err
                }
                switch order.Status {
                case model.OrderStatusCanceled, model.OrderStatusCompleted, model.OrderStatusRefunded:
                        return ErrValidation
                }
                if _, err := r.Players.Get(ctx, in.PlayerID); err != nil {
                        return err
                }
                order.SetPlayerID(in.PlayerID)
                order.AssignmentSource = in.Source
                if order.DisputeStatus == model.OrderDisputeStatusPending {
                        order.DisputeStatus = model.OrderDisputeStatusInMediation
                }
                if err := r.Orders.Update(ctx, order); err != nil {
                        return err
                }
                updated = order
                metadata := map[string]any{
                        "player_id": in.PlayerID,
                        "source":    in.Source,
                }
                raw, _ := json.Marshal(metadata)
                log := &model.OperationLog{
                        EntityType:   string(model.OpEntityOrder),
                        EntityID:     order.ID,
                        ActorUserID:  in.ActorUserID,
                        Action:       string(model.OpActionAssignPlayer),
                        MetadataJSON: raw,
                        TraceID:      in.TraceID,
                }
                if err := r.OpLogs.Append(ctx, log); err != nil {
                        return err
                }
                return nil
        })
        if err != nil {
                return nil, err
        }
        return updated, nil
}

// CancelAssignment 回退指派。
func (s *Service) CancelAssignment(ctx context.Context, orderID uint64, in CancelAssignInput) (*model.Order, error) {
        if s.tx == nil {
                return nil, errors.New("transaction manager not configured")
        }
        var updated *model.Order
        err := s.tx.WithTx(ctx, func(r *common.Repos) error {
                order, err := r.Orders.Get(ctx, orderID)
                if err != nil {
                        return err
                }
                order.PlayerID = nil
                order.AssignmentSource = model.OrderAssignmentSourceRollback
                if err := r.Orders.Update(ctx, order); err != nil {
                        return err
                }
                updated = order
                metadata := map[string]any{"reason": in.Reason}
                raw, _ := json.Marshal(metadata)
                log := &model.OperationLog{
                        EntityType:   string(model.OpEntityOrder),
                        EntityID:     order.ID,
                        ActorUserID:  in.ActorUserID,
                        Action:       string(model.OpActionAssignRollback),
                        Reason:       in.Reason,
                        MetadataJSON: raw,
                        TraceID:      in.TraceID,
                }
                if err := r.OpLogs.Append(ctx, log); err != nil {
                        return err
                }
                return nil
        })
        if err != nil {
                return nil, err
        }
        return updated, nil
}

// CreateDispute 发起争议。
func (s *Service) CreateDispute(ctx context.Context, orderID uint64, req DisputeRequest) (*model.OrderDispute, error) {
        if s.tx == nil {
                return nil, errors.New("transaction manager not configured")
        }
        var created *model.OrderDispute
        err := s.tx.WithTx(ctx, func(r *common.Repos) error {
                order, err := r.Orders.Get(ctx, orderID)
                if err != nil {
                        return err
                }
                withinWindow := time.Since(order.CreatedAt) <= 24*time.Hour
                if order.CompletedAt != nil {
                        withinWindow = withinWindow || time.Since(*order.CompletedAt) <= 24*time.Hour
                }
                if !withinWindow {
                        return ErrValidation
                }
                dispute := &model.OrderDispute{
                        OrderID:          orderID,
                        RaisedBy:         req.RaisedBy,
                        RaisedByUserID:   req.RaisedByUserID,
                        Reason:           req.Reason,
                        EvidenceURLs:     req.EvidenceURLs,
                        Status:           model.OrderDisputeStatusPending,
                        ResponseDeadline: time.Now().Add(s.slaDuration),
                        TraceID:          req.TraceID,
                }
                if err := r.Disputes.Create(ctx, dispute); err != nil {
                        return err
                }
                order.DisputeStatus = model.OrderDisputeStatusPending
                if err := r.Orders.Update(ctx, order); err != nil {
                        return err
                }
                metadata := map[string]any{
                        "raised_by": req.RaisedBy,
                }
                raw, _ := json.Marshal(metadata)
                logEntry := &model.OperationLog{
                        EntityType:   string(model.OpEntityOrder),
                        EntityID:     order.ID,
                        ActorUserID:  req.RaisedByUserID,
                        Action:       string(model.OpActionCreateDispute),
                        MetadataJSON: raw,
                        TraceID:      req.TraceID,
                }
                if err := r.OpLogs.Append(ctx, logEntry); err != nil {
                        return err
                }
                created = dispute
                return nil
        })
        if err != nil {
                return nil, err
        }
        return created, nil
}

// ListDisputes 查询争议列表。
func (s *Service) ListDisputes(ctx context.Context, orderID uint64) ([]model.OrderDispute, error) {
        return s.disputes.ListByOrder(ctx, orderID)
}

// MediateDispute 处理争议并执行退款/重派/驳回。
func (s *Service) MediateDispute(ctx context.Context, orderID uint64, in MediateInput) (*model.OrderDispute, error) {
        if s.tx == nil {
                return nil, errors.New("transaction manager not configured")
        }
        var (
                resolved     *model.OrderDispute
                orderSnapshot model.Order
        )
        err := s.tx.WithTx(ctx, func(r *common.Repos) error {
                order, err := r.Orders.Get(ctx, orderID)
                if err != nil {
                        return err
                }
                dispute, err := r.Disputes.GetLatestByOrder(ctx, orderID)
                if err != nil {
                        return err
                }
                if dispute.Status == model.OrderDisputeStatusResolved {
                        return ErrValidation
                }
                now := time.Now()
                dispute.Status = model.OrderDisputeStatusResolved
                dispute.Resolution = in.Resolution
                dispute.ResolutionNote = in.Note
                dispute.HandledByID = in.ActorUserID
                dispute.HandledAt = &now
                dispute.TraceID = in.TraceID
                dispute.RespondedAt = &now
                if in.RefundAmountCents != nil {
                        dispute.RefundAmountCents = *in.RefundAmountCents
                }
                if err := r.Disputes.Update(ctx, dispute); err != nil {
                        return err
                }
                switch in.Resolution {
                case model.OrderDisputeResolutionRefund:
                        if in.RefundAmountCents != nil {
                                order.RefundAmountCents = *in.RefundAmountCents
                        }
                        order.Status = model.OrderStatusRefunded
                        order.DisputeStatus = model.OrderDisputeStatusResolved
                case model.OrderDisputeResolutionReassign:
                        if in.ReassignPlayerID != nil && *in.ReassignPlayerID != 0 {
                                if _, err := r.Players.Get(ctx, *in.ReassignPlayerID); err != nil {
                                        return err
                                }
                                order.SetPlayerID(*in.ReassignPlayerID)
                                order.AssignmentSource = model.OrderAssignmentSourceManual
                                order.Status = model.OrderStatusConfirmed
                        } else {
                                order.PlayerID = nil
                                order.AssignmentSource = model.OrderAssignmentSourceRollback
                                order.Status = model.OrderStatusPending
                        }
                        order.DisputeStatus = model.OrderDisputeStatusResolved
                case model.OrderDisputeResolutionReject:
                        order.DisputeStatus = model.OrderDisputeStatusResolved
                }
                if err := r.Orders.Update(ctx, order); err != nil {
                        return err
                }
                metadata := map[string]any{
                        "resolution": dispute.Resolution,
                        "note":       dispute.ResolutionNote,
                        "refund":     dispute.RefundAmountCents,
                }
                raw, _ := json.Marshal(metadata)
                logEntry := &model.OperationLog{
                        EntityType:   string(model.OpEntityOrder),
                        EntityID:     order.ID,
                        ActorUserID:  in.ActorUserID,
                        Action:       string(model.OpActionMediateDispute),
                        MetadataJSON: raw,
                        TraceID:      in.TraceID,
                }
                if err := r.OpLogs.Append(ctx, logEntry); err != nil {
                        return err
                }
                resolved = dispute
                orderSnapshot = *order
                return nil
        })
        if err != nil {
                return nil, err
        }
        if resolved != nil && s.notifications != nil {
                go s.sendDisputeNotification(context.Background(), orderSnapshot, *resolved)
        }
        return resolved, nil
}

func (s *Service) sendDisputeNotification(ctx context.Context, order model.Order, dispute model.OrderDispute) {
        if s.notifications == nil {
                return
        }
        title := "订单争议已处理"
        message := "您的订单争议已完成调解，结果：" + string(dispute.Resolution)
        evt := &model.NotificationEvent{
                UserID:        order.UserID,
                Title:         title,
                Message:       message,
                Priority:      model.NotificationPriorityHigh,
                ReferenceType: "order_dispute",
                ReferenceID:   &dispute.OrderID,
        }
        _ = s.notifications.Create(ctx, evt)
        if order.PlayerID != nil {
                if player, err := s.players.Get(ctx, order.GetPlayerID()); err == nil {
                        evt := &model.NotificationEvent{
                                UserID:        player.UserID,
                                Title:         title,
                                Message:       message,
                                Priority:      model.NotificationPriorityNormal,
                                ReferenceType: "order_dispute",
                                ReferenceID:   &dispute.OrderID,
                        }
                        _ = s.notifications.Create(ctx, evt)
                }
        }
}

// TraceIDFromContext 帮助函数，默认读取 logging 包上下文。
func TraceIDFromContext(ctx context.Context) string {
        if traceID, ok := logging.RequestIDFromContext(ctx); ok {
                return traceID
        }
        return ""
}
