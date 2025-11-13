//go:build team_service
// +build team_service

package team

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/model"
)

// Service provides team order assignment and payout operations.
type Service struct {
	teams             teamRepository
	members           teamMemberRepository
	assignments       assignmentRepository
	assignmentMembers assignmentMemberRepository
	payoutPlans       payoutPlanRepository
	orders            orderRepository
	opLogs            opLogRepository

	// now is injectable time provider for tests.
	now func() time.Time
}

// --- repositories (minimal interfaces required by tests) ---

type teamRepository interface{}

type teamMemberRepository interface {
	GetByTeamAndUser(ctx context.Context, teamID, userID uint64) (*model.TeamMember, error)
}

type assignmentRepository interface {
	Create(ctx context.Context, a *model.TeamOrderAssignment) error
	GetByOrderID(ctx context.Context, orderID uint64) (*model.TeamOrderAssignment, error)
	Update(ctx context.Context, a *model.TeamOrderAssignment) error
	ListExpired(ctx context.Context, now time.Time) ([]model.TeamOrderAssignment, error)
}

type assignmentMemberRepository interface {
	Replace(ctx context.Context, assignmentID uint64, members []*model.TeamAssignmentMember) error
}

type payoutPlanRepository interface {
	Upsert(ctx context.Context, plan *model.TeamPayoutPlan) error
	GetByAssignment(ctx context.Context, assignmentID uint64) (*model.TeamPayoutPlan, error)
}

type orderRepository interface {
	AcquireForSnatch(ctx context.Context, orderID uint64) (*model.Order, error)
	UpdateColumns(ctx context.Context, orderID uint64, values map[string]any) error
}

type opLogRepository interface {
	Append(ctx context.Context, log *model.OperationLog) error
}

// --- request/response types used by tests ---

type SnatchRequest struct {
	OrderID uint64
	TeamID  uint64
}

type SnatchResponse struct {
	AssignmentID uint64
}

type PayoutShare struct {
	MemberID uint64
	Percent  int
}

type PayoutPlanRequest struct {
	OrderID    uint64
	TeamID     uint64
	ProfitMode model.TeamProfitMode
	Shares     []PayoutShare
}

type PayoutPlanResponse struct {
	ProfitMode model.TeamProfitMode
	Shares     []PayoutShare
}

func (s *Service) nowTime() time.Time {
	if s.now != nil {
		return s.now()
	}
	return time.Now().UTC()
}

// SnatchOrder allows a team leader to claim an order for team dispatching.
func (s *Service) SnatchOrder(ctx context.Context, leaderID uint64, req SnatchRequest) (*SnatchResponse, error) {
	// validate leader role
	member, err := s.members.GetByTeamAndUser(ctx, req.TeamID, leaderID)
	if err != nil {
		return nil, err
	}
	if member.Role != model.TeamRoleLeader || member.Status != model.TeamMemberStatusActive {
		return nil, errors.New("permission denied: not active team leader")
	}

	// acquire order for snatch
	order, err := s.orders.AcquireForSnatch(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}
	// basic validation according to test data
	if order.Status != model.OrderStatusPending {
		return nil, fmt.Errorf("invalid order status: %s", order.Status)
	}
	if order.QueueType != model.OrderQueueTypeTeam {
		return nil, errors.New("order is not in team queue")
	}

	// create assignment (first wins, others will fail inside repo)
	now := s.nowTime()
	assignment := &model.TeamOrderAssignment{
		OrderID:          req.OrderID,
		TeamID:           req.TeamID,
		Status:           model.TeamAssignmentStatusDispatching,
		DispatchDeadline: now.Add(10 * time.Minute),
		LockedAt:         &now,
	}
	if err := s.assignments.Create(ctx, assignment); err != nil {
		return nil, err
	}

	// update order assignment fields
	if err := s.orders.UpdateColumns(ctx, req.OrderID, map[string]any{
		"assigned_team_id":  req.TeamID,
		"assignment_source": model.OrderAssignmentSourceTeam,
	}); err != nil {
		return nil, err
	}

	// append simple op log (best-effort)
	_ = s.opLogs.Append(ctx, &model.OperationLog{
		EntityType: string(model.OpEntityOrder),
		EntityID:   req.OrderID,
		Action:     string(model.OpActionAssignPlayer),
	})

	return &SnatchResponse{AssignmentID: assignment.ID}, nil
}

// ReleaseExpiredAssignments releases timed-out team assignments and clears order assignment.
func (s *Service) ReleaseExpiredAssignments(ctx context.Context) (int, error) {
	now := s.nowTime()
	expired, err := s.assignments.ListExpired(ctx, now)
	if err != nil {
		return 0, err
	}
	count := 0
	for i := range expired {
		a := expired[i]
		a.Status = model.TeamAssignmentStatusReleased
		a.ReleasedAt = &now
		if err := s.assignments.Update(ctx, &a); err != nil {
			return count, err
		}
		// clear order assigned team
		if err := s.orders.UpdateColumns(ctx, a.OrderID, map[string]any{
			"assigned_team_id": nil,
		}); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

// UpsertPayoutPlan creates or updates payout plan for an assignment. Leader only.
func (s *Service) UpsertPayoutPlan(ctx context.Context, leaderID uint64, req PayoutPlanRequest) (*PayoutPlanResponse, error) {
	// validate leader role
	member, err := s.members.GetByTeamAndUser(ctx, req.TeamID, leaderID)
	if err != nil {
		return nil, err
	}
	if member.Role != model.TeamRoleLeader || member.Status != model.TeamMemberStatusActive {
		return nil, errors.New("permission denied: not active team leader")
	}

	// fetch assignment by order
	assignment, err := s.assignments.GetByOrderID(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if assignment.TeamID != req.TeamID {
		return nil, errors.New("assignment team mismatch")
	}

	// validate custom shares sum to 100
	if req.ProfitMode == model.TeamProfitModeCustom {
		sum := 0
		for _, s := range req.Shares {
			sum += s.Percent
		}
		if sum != 100 {
			return nil, errors.New("shares percent must sum to 100")
		}
	}

	// persist plan
	plan := &model.TeamPayoutPlan{
		AssignmentID: assignment.ID,
		ProfitMode:   req.ProfitMode,
	}
	if err := s.payoutPlans.Upsert(ctx, plan); err != nil {
		return nil, err
	}

	return &PayoutPlanResponse{
		ProfitMode: req.ProfitMode,
		Shares:     req.Shares,
	}, nil
}
