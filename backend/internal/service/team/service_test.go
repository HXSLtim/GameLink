package team

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	orderrepo "gamelink/internal/repository/order"
)

func TestServiceSnatchOrderConcurrent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	env := newTestEnv()

	const (
		leaderID = uint64(11)
		teamID   = uint64(77)
		orderID  = uint64(900)
	)

	env.members.addMember(teamID, leaderID, model.TeamRoleLeader, model.TeamMemberStatusActive)
	env.orders.saveOrder(&model.Order{Base: model.Base{ID: orderID}, QueueType: model.OrderQueueTypeTeam, RequiredMembers: 3, Status: model.OrderStatusPending})

	var success int32
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := env.svc.SnatchOrder(ctx, leaderID, SnatchRequest{OrderID: orderID, TeamID: teamID}); err == nil {
				atomic.AddInt32(&success, 1)
			}
		}()
	}
	wg.Wait()

	require.Equal(t, int32(1), success, "only one snatch should succeed")
	assignment, err := env.assignments.GetByOrderID(ctx, orderID)
	require.NoError(t, err)
	require.Equal(t, teamID, assignment.TeamID)
	require.Equal(t, model.TeamAssignmentStatusDispatching, assignment.Status)
	require.NotNil(t, env.orders.orders[orderID].AssignedTeamID)
	require.Equal(t, teamID, *env.orders.orders[orderID].AssignedTeamID)
}

func TestServiceReleaseExpiredAssignments(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	env := newTestEnv()

	const (
		teamID  = uint64(50)
		orderID = uint64(200)
	)

	assignedTeamID := teamID
	env.orders.saveOrder(&model.Order{Base: model.Base{ID: orderID}, AssignedTeamID: &assignedTeamID})

	env.nowValue = env.nowValue.Add(-10 * time.Minute)
	err := env.assignments.Create(ctx, &model.TeamOrderAssignment{
		OrderID:          orderID,
		TeamID:           teamID,
		Status:           model.TeamAssignmentStatusDispatching,
		DispatchDeadline: env.nowValue.Add(-time.Minute),
		LockedAt:         env.nowValue.Add(-5 * time.Minute),
	})
	require.NoError(t, err)

	env.nowValue = env.nowValue.Add(15 * time.Minute)
	released, err := env.svc.ReleaseExpiredAssignments(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, released)

	updated, err := env.assignments.GetByOrderID(ctx, orderID)
	require.NoError(t, err)
	require.Equal(t, model.TeamAssignmentStatusReleased, updated.Status)
	require.NotNil(t, updated.ReleasedAt)
	require.Nil(t, env.orders.orders[orderID].AssignedTeamID)
}

func TestServiceUpsertPayoutPlan(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	env := newTestEnv()

	const (
		leaderID = uint64(4)
		teamID   = uint64(5)
		orderID  = uint64(6)
	)

	env.members.addMember(teamID, leaderID, model.TeamRoleLeader, model.TeamMemberStatusActive)
	require.NoError(t, env.assignments.Create(ctx, &model.TeamOrderAssignment{OrderID: orderID, TeamID: teamID}))

	resp, err := env.svc.UpsertPayoutPlan(ctx, leaderID, PayoutPlanRequest{
		OrderID:    orderID,
		TeamID:     teamID,
		ProfitMode: model.TeamProfitModeCustom,
		Shares: []PayoutShare{
			{MemberID: 1, Percent: 60},
			{MemberID: 2, Percent: 40},
		},
	})
	require.NoError(t, err)
	require.Equal(t, model.TeamProfitModeCustom, resp.ProfitMode)
	require.Len(t, resp.Shares, 2)

	plan, err := env.payoutPlans.GetByAssignment(ctx, env.assignments.mustGet(orderID).ID)
	require.NoError(t, err)
	require.Equal(t, model.TeamProfitModeCustom, plan.ProfitMode)
}

func TestServiceUpsertPayoutPlanValidation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	env := newTestEnv()

	const (
		leaderID = uint64(8)
		teamID   = uint64(9)
		orderID  = uint64(10)
	)

	env.members.addMember(teamID, leaderID, model.TeamRoleLeader, model.TeamMemberStatusActive)
	require.NoError(t, env.assignments.Create(ctx, &model.TeamOrderAssignment{OrderID: orderID, TeamID: teamID}))

	_, err := env.svc.UpsertPayoutPlan(ctx, leaderID, PayoutPlanRequest{
		OrderID: orderID,
		TeamID:  teamID,
		Shares:  []PayoutShare{{MemberID: 1, Percent: 90}},
	})
	require.Error(t, err)
}

// --- test helpers ---

type testEnv struct {
	svc          *Service
	teams        *fakeTeamRepo
	members      *fakeTeamMemberRepo
	assignments  *fakeAssignmentRepo
	assignMember *fakeAssignmentMemberRepo
	payoutPlans  *fakePayoutPlanRepo
	orders       *fakeOrderRepo
	opLogs       *fakeOpLogRepo
	nowValue     time.Time
}

func newTestEnv() *testEnv {
	teams := &fakeTeamRepo{}
	members := newFakeTeamMemberRepo()
	assignments := newFakeAssignmentRepo()
	assignMembers := newFakeAssignmentMemberRepo()
	payoutPlans := newFakePayoutPlanRepo()
	orders := newFakeOrderRepo()
	opLogs := newFakeOpLogRepo()
	env := &testEnv{
		teams:        teams,
		members:      members,
		assignments:  assignments,
		assignMember: assignMembers,
		payoutPlans:  payoutPlans,
		orders:       orders,
		opLogs:       opLogs,
		nowValue:     time.Now().UTC(),
	}
	svc := &Service{
		teams:             teams,
		members:           members,
		assignments:       assignments,
		assignmentMembers: assignMembers,
		payoutPlans:       payoutPlans,
		orders:            orders,
		opLogs:            opLogs,
	}
	svc.now = func() time.Time {
		return env.nowValue
	}
	env.svc = svc
	return env
}

type fakeTeamRepo struct{}

func (r *fakeTeamRepo) Create(ctx context.Context, team *model.Team) error { return nil }
func (r *fakeTeamRepo) Get(ctx context.Context, id uint64) (*model.Team, error) {
	return nil, repository.ErrNotFound
}
func (r *fakeTeamRepo) ListByLeader(ctx context.Context, leaderID uint64) ([]model.Team, error) {
	return nil, nil
}
func (r *fakeTeamRepo) ListByUser(ctx context.Context, userID uint64) ([]model.Team, error) {
	return nil, nil
}

type fakeTeamMemberRepo struct {
	mu     sync.Mutex
	byTeam map[uint64]map[uint64]*model.TeamMember
	byID   map[uint64]*model.TeamMember
	nextID uint64
}

func newFakeTeamMemberRepo() *fakeTeamMemberRepo {
	return &fakeTeamMemberRepo{byTeam: make(map[uint64]map[uint64]*model.TeamMember), byID: make(map[uint64]*model.TeamMember)}
}

func (r *fakeTeamMemberRepo) addMember(teamID, userID uint64, role model.TeamRole, status model.TeamMemberStatus) uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	member := &model.TeamMember{TeamID: teamID, UserID: userID, Role: role, Status: status}
	member.ID = r.nextID
	if _, ok := r.byTeam[teamID]; !ok {
		r.byTeam[teamID] = make(map[uint64]*model.TeamMember)
	}
	r.byTeam[teamID][userID] = member
	r.byID[member.ID] = member
	return member.ID
}

func (r *fakeTeamMemberRepo) Create(ctx context.Context, member *model.TeamMember) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	member.ID = r.nextID
	if _, ok := r.byTeam[member.TeamID]; !ok {
		r.byTeam[member.TeamID] = make(map[uint64]*model.TeamMember)
	}
	cp := *member
	r.byTeam[member.TeamID][member.UserID] = &cp
	r.byID[cp.ID] = &cp
	return nil
}

func (r *fakeTeamMemberRepo) ListByTeam(ctx context.Context, teamID uint64) ([]model.TeamMember, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]model.TeamMember, 0)
	for _, member := range r.byTeam[teamID] {
		result = append(result, *member)
	}
	return result, nil
}

func (r *fakeTeamMemberRepo) GetByTeamAndUser(ctx context.Context, teamID, userID uint64) (*model.TeamMember, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if teamMembers, ok := r.byTeam[teamID]; ok {
		if member, ok := teamMembers[userID]; ok {
			cp := *member
			return &cp, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *fakeTeamMemberRepo) Update(ctx context.Context, member *model.TeamMember) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if existing, ok := r.byID[member.ID]; ok {
		existing.Role = member.Role
		existing.Status = member.Status
		existing.ProfitShareDefault = member.ProfitShareDefault
		return nil
	}
	return repository.ErrNotFound
}

type fakeAssignmentRepo struct {
	mu      sync.Mutex
	byOrder map[uint64]*model.TeamOrderAssignment
	nextID  uint64
}

func newFakeAssignmentRepo() *fakeAssignmentRepo {
	return &fakeAssignmentRepo{byOrder: make(map[uint64]*model.TeamOrderAssignment)}
}

func (r *fakeAssignmentRepo) Create(ctx context.Context, assignment *model.TeamOrderAssignment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if existing, ok := r.byOrder[assignment.OrderID]; ok && existing.Status != model.TeamAssignmentStatusReleased {
		return fmt.Errorf("assignment exists")
	}
	r.nextID++
	if assignment.ID == 0 {
		assignment.ID = r.nextID
	}
	cp := *assignment
	r.byOrder[assignment.OrderID] = &cp
	return nil
}

func (r *fakeAssignmentRepo) GetByOrderID(ctx context.Context, orderID uint64) (*model.TeamOrderAssignment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if assignment, ok := r.byOrder[orderID]; ok {
		cp := *assignment
		return &cp, nil
	}
	return nil, repository.ErrNotFound
}

func (r *fakeAssignmentRepo) Update(ctx context.Context, assignment *model.TeamOrderAssignment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byOrder[assignment.OrderID]; !ok {
		return repository.ErrNotFound
	}
	cp := *assignment
	r.byOrder[assignment.OrderID] = &cp
	return nil
}

func (r *fakeAssignmentRepo) ListExpired(ctx context.Context, now time.Time) ([]model.TeamOrderAssignment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []model.TeamOrderAssignment
	for _, assignment := range r.byOrder {
		if assignment.Status == model.TeamAssignmentStatusDispatching && assignment.DispatchDeadline.Before(now) {
			result = append(result, *assignment)
		}
	}
	return result, nil
}

func (r *fakeAssignmentRepo) mustGet(orderID uint64) *model.TeamOrderAssignment {
	r.mu.Lock()
	defer r.mu.Unlock()
	if assignment, ok := r.byOrder[orderID]; ok {
		return assignment
	}
	return nil
}

type fakeAssignmentMemberRepo struct {
	mu      sync.Mutex
	records map[uint64][]model.TeamAssignmentMember
}

func newFakeAssignmentMemberRepo() *fakeAssignmentMemberRepo {
	return &fakeAssignmentMemberRepo{records: make(map[uint64][]model.TeamAssignmentMember)}
}

func (r *fakeAssignmentMemberRepo) Replace(ctx context.Context, assignmentID uint64, members []*model.TeamAssignmentMember) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cloned := make([]model.TeamAssignmentMember, len(members))
	for i, m := range members {
		cloned[i] = *m
	}
	r.records[assignmentID] = cloned
	return nil
}

func (r *fakeAssignmentMemberRepo) ListByAssignment(ctx context.Context, assignmentID uint64) ([]model.TeamAssignmentMember, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if members, ok := r.records[assignmentID]; ok {
		copied := make([]model.TeamAssignmentMember, len(members))
		copy(copied, members)
		return copied, nil
	}
	return nil, nil
}

func (r *fakeAssignmentMemberRepo) UpdateState(ctx context.Context, assignmentID, memberID uint64, state model.TeamAssignmentMemberState, confirmedAt *time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	members := r.records[assignmentID]
	for i := range members {
		if members[i].MemberID == memberID {
			members[i].State = state
			members[i].ConfirmedAt = confirmedAt
			r.records[assignmentID][i] = members[i]
			return nil
		}
	}
	return repository.ErrNotFound
}

type fakePayoutPlanRepo struct {
	mu    sync.Mutex
	plans map[uint64]*model.TeamPayoutPlan
}

func newFakePayoutPlanRepo() *fakePayoutPlanRepo {
	return &fakePayoutPlanRepo{plans: make(map[uint64]*model.TeamPayoutPlan)}
}

func (r *fakePayoutPlanRepo) Upsert(ctx context.Context, plan *model.TeamPayoutPlan) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *plan
	r.plans[plan.AssignmentID] = &cp
	return nil
}

func (r *fakePayoutPlanRepo) GetByAssignment(ctx context.Context, assignmentID uint64) (*model.TeamPayoutPlan, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if plan, ok := r.plans[assignmentID]; ok {
		cp := *plan
		return &cp, nil
	}
	return nil, repository.ErrNotFound
}

type fakeOrderRepo struct {
	mu     sync.Mutex
	orders map[uint64]*model.Order
}

func newFakeOrderRepo() *fakeOrderRepo {
	return &fakeOrderRepo{orders: make(map[uint64]*model.Order)}
}

func (r *fakeOrderRepo) ListPool(ctx context.Context, opts orderrepo.PoolListOptions) ([]model.Order, *uint64, error) {
	return nil, nil, nil
}

func (r *fakeOrderRepo) AcquireForSnatch(ctx context.Context, orderID uint64) (*model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if order, ok := r.orders[orderID]; ok {
		cp := *order
		return &cp, nil
	}
	return nil, repository.ErrNotFound
}

func (r *fakeOrderRepo) UpdateColumns(ctx context.Context, orderID uint64, values map[string]any) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, ok := r.orders[orderID]
	if !ok {
		return repository.ErrNotFound
	}
	for key, val := range values {
		switch key {
		case "assigned_team_id":
			if val == nil {
				order.AssignedTeamID = nil
			} else {
				switch v := val.(type) {
				case uint64:
					order.AssignedTeamID = &v
				case *uint64:
					order.AssignedTeamID = v
				}
			}
		case "assignment_source":
			switch v := val.(type) {
			case model.OrderAssignmentSource:
				order.AssignmentSource = v
			case string:
				order.AssignmentSource = model.OrderAssignmentSource(v)
			}
		}
	}
	return nil
}

func (r *fakeOrderRepo) saveOrder(order *model.Order) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
}

type fakeOpLogRepo struct {
	mu   sync.Mutex
	logs []model.OperationLog
}

func newFakeOpLogRepo() *fakeOpLogRepo { return &fakeOpLogRepo{} }

func (r *fakeOpLogRepo) Append(ctx context.Context, log *model.OperationLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logs = append(r.logs, *log)
	return nil
}

func (r *fakeOpLogRepo) ListByEntity(ctx context.Context, entityType string, entityID uint64, opts repository.OperationLogListOptions) ([]model.OperationLog, int64, error) {
	return nil, 0, nil
}
