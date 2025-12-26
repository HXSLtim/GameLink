package team

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	teamrepo "gamelink/internal/repository/team"
)

type MockTeamRepository struct {
	mock.Mock
}

func (m *MockTeamRepository) Create(ctx context.Context, team *model.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *MockTeamRepository) GetByID(ctx context.Context, id uint64) (*model.Team, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Team), args.Error(1)
}

func (m *MockTeamRepository) Update(ctx context.Context, team *model.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *MockTeamRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTeamRepository) List(ctx context.Context, opts teamrepo.TeamListOptions) ([]model.Team, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.Team), args.Get(1).(int64), args.Error(2)
}

func (m *MockTeamRepository) UpdateStatus(ctx context.Context, id uint64, status model.TeamStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockTeamRepository) UpdateMemberCount(ctx context.Context, id uint64, delta int) error {
	args := m.Called(ctx, id, delta)
	return args.Error(0)
}

func (m *MockTeamRepository) UpdateLeader(ctx context.Context, id uint64, leaderID uint64) error {
	args := m.Called(ctx, id, leaderID)
	return args.Error(0)
}

func (m *MockTeamRepository) GetTeamStats(ctx context.Context) (*teamrepo.TeamStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*teamrepo.TeamStats), args.Error(1)
}

func (m *MockTeamRepository) CreateMember(ctx context.Context, member *model.TeamMember) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockTeamRepository) GetMemberByTeamAndPlayer(ctx context.Context, teamID, playerID uint64) (*model.TeamMember, error) {
	args := m.Called(ctx, teamID, playerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TeamMember), args.Error(1)
}

func (m *MockTeamRepository) GetActiveMemberByPlayer(ctx context.Context, playerID uint64) (*model.TeamMember, error) {
	args := m.Called(ctx, playerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TeamMember), args.Error(1)
}

func (m *MockTeamRepository) UpdateMember(ctx context.Context, member *model.TeamMember) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockTeamRepository) GetActiveMembers(ctx context.Context, teamID uint64) ([]model.TeamMember, error) {
	args := m.Called(ctx, teamID)
	return args.Get(0).([]model.TeamMember), args.Error(1)
}

func (m *MockTeamRepository) ListMembers(ctx context.Context, opts teamrepo.MemberListOptions) ([]model.TeamMember, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.TeamMember), args.Get(1).(int64), args.Error(2)
}

func (m *MockTeamRepository) GetNextLeader(ctx context.Context, teamID uint64, excludePlayerID uint64) (*model.TeamMember, error) {
	args := m.Called(ctx, teamID, excludePlayerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TeamMember), args.Error(1)
}

func (m *MockTeamRepository) CreateInvite(ctx context.Context, invite *model.TeamInvite) error {
	args := m.Called(ctx, invite)
	return args.Error(0)
}

func (m *MockTeamRepository) GetInviteByID(ctx context.Context, id uint64) (*model.TeamInvite, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TeamInvite), args.Error(1)
}

func (m *MockTeamRepository) GetPendingInvite(ctx context.Context, teamID, playerID uint64) (*model.TeamInvite, error) {
	args := m.Called(ctx, teamID, playerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TeamInvite), args.Error(1)
}

func (m *MockTeamRepository) UpdateInvite(ctx context.Context, invite *model.TeamInvite) error {
	args := m.Called(ctx, invite)
	return args.Error(0)
}

func (m *MockTeamRepository) ListInvites(ctx context.Context, opts teamrepo.InviteListOptions) ([]model.TeamInvite, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.TeamInvite), args.Get(1).(int64), args.Error(2)
}

func TestNewTeamService(t *testing.T) {
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)
	require.NotNil(t, svc)
}

func TestTeamService_CreateTeam(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		repo.On("GetActiveMemberByPlayer", ctx, uint64(1)).Return(nil, repository.ErrNotFound)
		repo.On("Create", ctx, mock.AnythingOfType("*model.Team")).Return(nil)
		repo.On("CreateMember", ctx, mock.AnythingOfType("*model.TeamMember")).Return(nil)
		team := &model.Team{Name: "Test Team"}
		err := svc.CreateTeam(ctx, team, 1)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), team.LeaderID)
	})

	t.Run("player already in team", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		existingMember := &model.TeamMember{PlayerID: 1}
		repo.On("GetActiveMemberByPlayer", ctx, uint64(1)).Return(existingMember, nil)
		team := &model.Team{Name: "Test Team"}
		err := svc.CreateTeam(ctx, team, 1)
		require.Error(t, err)
	})

	t.Run("check error", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		repo.On("GetActiveMemberByPlayer", ctx, uint64(1)).Return(nil, errors.New("db error"))
		team := &model.Team{Name: "Test Team"}
		err := svc.CreateTeam(ctx, team, 1)
		require.Error(t, err)
	})
}

func TestTeamService_GetTeam(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		expected := &model.Team{Name: "Test Team"}
		expected.ID = 1
		repo.On("GetByID", ctx, uint64(1)).Return(expected, nil)
		result, err := svc.GetTeam(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, expected.Name, result.Name)
	})

	t.Run("not found", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		repo.On("GetByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound)
		_, err := svc.GetTeam(ctx, 999)
		require.Error(t, err)
	})
}

func TestTeamService_ListTeams(t *testing.T) {
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		teams := []model.Team{{Name: "Team 1"}, {Name: "Team 2"}}
		opts := teamrepo.TeamListOptions{Page: 1, PageSize: 10}
		repo.On("List", ctx, opts).Return(teams, int64(2), nil)
		result, total, err := svc.ListTeams(ctx, opts)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(2), total)
	})
}

func TestTeamService_GetTeamStats(t *testing.T) {
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		stats := &teamrepo.TeamStats{TotalTeams: 10, ActiveTeams: 8}
		repo.On("GetTeamStats", ctx).Return(stats, nil)
		result, err := svc.GetTeamStats(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(10), result.TotalTeams)
	})
}

func TestTeamService_GetTeamMembers(t *testing.T) {
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		members := []model.TeamMember{{PlayerID: 1}, {PlayerID: 2}}
		repo.On("GetActiveMembers", ctx, uint64(1)).Return(members, nil)
		result, err := svc.GetTeamMembers(ctx, 1)
		require.NoError(t, err)
		assert.Len(t, result, 2)
	})
}

func TestTeamService_ListMembers(t *testing.T) {
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		members := []model.TeamMember{{PlayerID: 1}}
		opts := teamrepo.MemberListOptions{Page: 1, PageSize: 10}
		repo.On("ListMembers", ctx, opts).Return(members, int64(1), nil)
		result, total, err := svc.ListMembers(ctx, opts)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), total)
	})
}

func TestTeamService_ListInvites(t *testing.T) {
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		invites := []model.TeamInvite{{PlayerID: 1}}
		opts := teamrepo.InviteListOptions{Page: 1, PageSize: 10}
		repo.On("ListInvites", ctx, opts).Return(invites, int64(1), nil)
		result, total, err := svc.ListInvites(ctx, opts)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), total)
	})
}

func TestTeamService_GetInvite(t *testing.T) {
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		invite := &model.TeamInvite{PlayerID: 1}
		invite.ID = 1
		repo.On("GetInviteByID", ctx, uint64(1)).Return(invite, nil)
		result, err := svc.GetInvite(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), result.PlayerID)
	})

	t.Run("not found", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		repo.On("GetInviteByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound)
		_, err := svc.GetInvite(ctx, 999)
		require.Error(t, err)
	})
}

func TestTeamService_UpdateTeam(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		existing := &model.Team{MemberCount: 3, MaxMembers: 5}
		existing.ID = 1
		updated := &model.Team{Name: "New Name", MaxMembers: 5}
		updated.ID = 1
		repo.On("GetByID", ctx, uint64(1)).Return(existing, nil)
		repo.On("Update", ctx, updated).Return(nil)
		err := svc.UpdateTeam(ctx, updated)
		require.NoError(t, err)
	})

	t.Run("cannot reduce max members", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		existing := &model.Team{MemberCount: 4, MaxMembers: 5}
		existing.ID = 1
		updated := &model.Team{MaxMembers: 3}
		updated.ID = 1
		repo.On("GetByID", ctx, uint64(1)).Return(existing, nil)
		err := svc.UpdateTeam(ctx, updated)
		require.Error(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		updated := &model.Team{MaxMembers: 5}
		updated.ID = 999
		repo.On("GetByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound)
		err := svc.UpdateTeam(ctx, updated)
		require.Error(t, err)
	})
}

func TestTeamService_DeleteTeam(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		team := &model.Team{CurrentOrderID: nil}
		team.ID = 1
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		repo.On("Delete", ctx, uint64(1)).Return(nil)
		err := svc.DeleteTeam(ctx, 1)
		require.NoError(t, err)
	})

	t.Run("has active order", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		orderID := uint64(100)
		team := &model.Team{CurrentOrderID: &orderID}
		team.ID = 1
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		err := svc.DeleteTeam(ctx, 1)
		require.Error(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		repo.On("GetByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound)
		err := svc.DeleteTeam(ctx, 999)
		require.Error(t, err)
	})
}

func TestTeamService_UpdateTeamStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		team := &model.Team{Status: model.TeamStatusActive, CurrentOrderID: nil}
		team.ID = 1
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		repo.On("UpdateStatus", ctx, uint64(1), model.TeamStatusInactive).Return(nil)
		err := svc.UpdateTeamStatus(ctx, 1, model.TeamStatusInactive)
		require.NoError(t, err)
	})

	t.Run("cannot set inactive with order", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		orderID := uint64(100)
		team := &model.Team{CurrentOrderID: &orderID}
		team.ID = 1
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		err := svc.UpdateTeamStatus(ctx, 1, model.TeamStatusInactive)
		require.Error(t, err)
	})
}

func TestTeamService_AddMember(t *testing.T) {
	ctx := context.Background()

	t.Run("success new member", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		team := &model.Team{MemberCount: 2, MaxMembers: 5}
		team.ID = 1
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		repo.On("GetActiveMemberByPlayer", ctx, uint64(2)).Return(nil, repository.ErrNotFound)
		repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(2)).Return(nil, repository.ErrNotFound)
		repo.On("CreateMember", ctx, mock.AnythingOfType("*model.TeamMember")).Return(nil)
		repo.On("UpdateMemberCount", ctx, uint64(1), 1).Return(nil)
		err := svc.AddMember(ctx, 1, 2)
		require.NoError(t, err)
	})

	t.Run("team full", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		team := &model.Team{MemberCount: 5, MaxMembers: 5}
		team.ID = 1
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		err := svc.AddMember(ctx, 1, 2)
		require.Error(t, err)
	})

	t.Run("player already in team", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		team := &model.Team{MemberCount: 2, MaxMembers: 5}
		team.ID = 1
		existingMember := &model.TeamMember{PlayerID: 2}
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		repo.On("GetActiveMemberByPlayer", ctx, uint64(2)).Return(existingMember, nil)
		err := svc.AddMember(ctx, 1, 2)
		require.Error(t, err)
	})

	t.Run("reactivate former member", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		team := &model.Team{MemberCount: 2, MaxMembers: 5}
		team.ID = 1
		formerMember := &model.TeamMember{TeamID: 1, PlayerID: 2, Status: model.TeamMemberStatusLeft}
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		repo.On("GetActiveMemberByPlayer", ctx, uint64(2)).Return(nil, repository.ErrNotFound)
		repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(2)).Return(formerMember, nil)
		repo.On("UpdateMember", ctx, mock.AnythingOfType("*model.TeamMember")).Return(nil)
		repo.On("UpdateMemberCount", ctx, uint64(1), 1).Return(nil)
		err := svc.AddMember(ctx, 1, 2)
		require.NoError(t, err)
	})
}

func TestTeamService_RemoveMember(t *testing.T) {
	ctx := context.Background()

	t.Run("success regular member", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		team := &model.Team{CurrentOrderID: nil}
		team.ID = 1
		member := &model.TeamMember{TeamID: 1, PlayerID: 2, Role: model.TeamMemberRoleMember, Status: model.TeamMemberStatusActive}
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(2)).Return(member, nil)
		repo.On("UpdateMember", ctx, mock.AnythingOfType("*model.TeamMember")).Return(nil)
		repo.On("UpdateMemberCount", ctx, uint64(1), -1).Return(nil)
		err := svc.RemoveMember(ctx, 1, 2, false)
		require.NoError(t, err)
	})

	t.Run("has active order", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		orderID := uint64(100)
		team := &model.Team{CurrentOrderID: &orderID}
		team.ID = 1
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		err := svc.RemoveMember(ctx, 1, 2, false)
		require.Error(t, err)
	})

	t.Run("member not active", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		team := &model.Team{CurrentOrderID: nil}
		team.ID = 1
		member := &model.TeamMember{Status: model.TeamMemberStatusLeft}
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(2)).Return(member, nil)
		err := svc.RemoveMember(ctx, 1, 2, false)
		require.Error(t, err)
	})

	t.Run("leader leaves transfer to next", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		team := &model.Team{CurrentOrderID: nil}
		team.ID = 1
		leader := &model.TeamMember{TeamID: 1, PlayerID: 1, Role: model.TeamMemberRoleLeader, Status: model.TeamMemberStatusActive}
		nextLeader := &model.TeamMember{TeamID: 1, PlayerID: 2, Role: model.TeamMemberRoleMember, Status: model.TeamMemberStatusActive}
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(leader, nil)
		repo.On("GetNextLeader", ctx, uint64(1), uint64(1)).Return(nextLeader, nil)
		repo.On("UpdateMember", ctx, mock.AnythingOfType("*model.TeamMember")).Return(nil)
		repo.On("UpdateLeader", ctx, uint64(1), uint64(2)).Return(nil)
		repo.On("UpdateMemberCount", ctx, uint64(1), -1).Return(nil)
		err := svc.RemoveMember(ctx, 1, 1, false)
		require.NoError(t, err)
	})

	t.Run("last member leaves delete team", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		team := &model.Team{CurrentOrderID: nil}
		team.ID = 1
		leader := &model.TeamMember{TeamID: 1, PlayerID: 1, Role: model.TeamMemberRoleLeader, Status: model.TeamMemberStatusActive}
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(leader, nil)
		repo.On("GetNextLeader", ctx, uint64(1), uint64(1)).Return(nil, repository.ErrNotFound)
		repo.On("Delete", ctx, uint64(1)).Return(nil)
		err := svc.RemoveMember(ctx, 1, 1, false)
		require.NoError(t, err)
	})
}

func TestTeamService_KickMember(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		team := &model.Team{CurrentOrderID: nil}
		team.ID = 1
		leader := &model.TeamMember{Role: model.TeamMemberRoleLeader}
		target := &model.TeamMember{TeamID: 1, PlayerID: 3, Role: model.TeamMemberRoleMember, Status: model.TeamMemberStatusActive}
		repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(2)).Return(leader, nil)
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(3)).Return(target, nil)
		repo.On("UpdateMember", ctx, mock.AnythingOfType("*model.TeamMember")).Return(nil)
		repo.On("UpdateMemberCount", ctx, uint64(1), -1).Return(nil)
		err := svc.KickMember(ctx, 1, 2, 3)
		require.NoError(t, err)
	})

	t.Run("not leader", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		member := &model.TeamMember{Role: model.TeamMemberRoleMember}
		repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(2)).Return(member, nil)
		err := svc.KickMember(ctx, 1, 2, 3)
		require.Error(t, err)
	})

	t.Run("cannot kick self", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		leader := &model.TeamMember{Role: model.TeamMemberRoleLeader}
		repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(leader, nil)
		err := svc.KickMember(ctx, 1, 1, 1)
		require.Error(t, err)
	})
}

func TestTeamService_TransferLeader(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		team := &model.Team{LeaderID: 1}
		team.ID = 1
		currentLeader := &model.TeamMember{TeamID: 1, PlayerID: 1, Role: model.TeamMemberRoleLeader}
		newLeader := &model.TeamMember{TeamID: 1, PlayerID: 2, Role: model.TeamMemberRoleMember, Status: model.TeamMemberStatusActive}
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(currentLeader, nil)
		repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(2)).Return(newLeader, nil)
		repo.On("UpdateMember", ctx, mock.AnythingOfType("*model.TeamMember")).Return(nil)
		repo.On("UpdateLeader", ctx, uint64(1), uint64(2)).Return(nil)
		err := svc.TransferLeader(ctx, 1, 1, 2)
		require.NoError(t, err)
	})

	t.Run("not current leader", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		team := &model.Team{LeaderID: 1}
		team.ID = 1
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		err := svc.TransferLeader(ctx, 1, 2, 3)
		require.Error(t, err)
	})

	t.Run("target not active", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		team := &model.Team{LeaderID: 1}
		team.ID = 1
		currentLeader := &model.TeamMember{TeamID: 1, PlayerID: 1, Role: model.TeamMemberRoleLeader}
		newLeader := &model.TeamMember{TeamID: 1, PlayerID: 2, Status: model.TeamMemberStatusLeft}
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(currentLeader, nil)
		repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(2)).Return(newLeader, nil)
		err := svc.TransferLeader(ctx, 1, 1, 2)
		require.Error(t, err)
	})
}

func TestTeamService_GetPlayerTeam(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		member := &model.TeamMember{TeamID: 1, PlayerID: 1}
		team := &model.Team{Name: "Test Team"}
		team.ID = 1
		repo.On("GetActiveMemberByPlayer", ctx, uint64(1)).Return(member, nil)
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		result, err := svc.GetPlayerTeam(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, "Test Team", result.Name)
	})

	t.Run("not in team", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		repo.On("GetActiveMemberByPlayer", ctx, uint64(1)).Return(nil, repository.ErrNotFound)
		_, err := svc.GetPlayerTeam(ctx, 1)
		require.Error(t, err)
	})
}

func TestTeamService_CreateInvite(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		team := &model.Team{MemberCount: 2, MaxMembers: 5}
		team.ID = 1
		inviter := &model.TeamMember{TeamID: 1, PlayerID: 1, Status: model.TeamMemberStatusActive}
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(inviter, nil)
		repo.On("GetActiveMemberByPlayer", ctx, uint64(2)).Return(nil, repository.ErrNotFound)
		repo.On("GetPendingInvite", ctx, uint64(1), uint64(2)).Return(nil, repository.ErrNotFound)
		repo.On("CreateInvite", ctx, mock.AnythingOfType("*model.TeamInvite")).Return(nil)
		err := svc.CreateInvite(ctx, 1, 1, 2, "Join us!")
		require.NoError(t, err)
	})

	t.Run("team full", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		team := &model.Team{MemberCount: 5, MaxMembers: 5}
		team.ID = 1
		inviter := &model.TeamMember{TeamID: 1, PlayerID: 1, Status: model.TeamMemberStatusActive}
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(inviter, nil)
		err := svc.CreateInvite(ctx, 1, 1, 2, "Join us!")
		require.Error(t, err)
	})

	t.Run("target already in team", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		team := &model.Team{MemberCount: 2, MaxMembers: 5}
		team.ID = 1
		inviter := &model.TeamMember{TeamID: 1, PlayerID: 1, Status: model.TeamMemberStatusActive}
		targetMember := &model.TeamMember{PlayerID: 2}
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(inviter, nil)
		repo.On("GetActiveMemberByPlayer", ctx, uint64(2)).Return(targetMember, nil)
		err := svc.CreateInvite(ctx, 1, 1, 2, "Join us!")
		require.Error(t, err)
	})

	t.Run("inviter not active", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		team := &model.Team{MemberCount: 2, MaxMembers: 5}
		team.ID = 1
		inviter := &model.TeamMember{TeamID: 1, PlayerID: 1, Status: model.TeamMemberStatusLeft}
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(inviter, nil)
		err := svc.CreateInvite(ctx, 1, 1, 2, "Join us!")
		require.Error(t, err)
	})

	t.Run("pending invite exists", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		team := &model.Team{MemberCount: 2, MaxMembers: 5}
		team.ID = 1
		inviter := &model.TeamMember{TeamID: 1, PlayerID: 1, Status: model.TeamMemberStatusActive}
		existingInvite := &model.TeamInvite{TeamID: 1, PlayerID: 2}
		repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
		repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(inviter, nil)
		repo.On("GetActiveMemberByPlayer", ctx, uint64(2)).Return(nil, repository.ErrNotFound)
		repo.On("GetPendingInvite", ctx, uint64(1), uint64(2)).Return(existingInvite, nil)
		err := svc.CreateInvite(ctx, 1, 1, 2, "Join us!")
		require.Error(t, err)
	})
}

func TestTeamService_AcceptInvite(t *testing.T) {
	ctx := context.Background()

	t.Run("not the invitee", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		invite := &model.TeamInvite{PlayerID: 2, Status: model.TeamInviteStatusPending}
		invite.ID = 1
		repo.On("GetInviteByID", ctx, uint64(1)).Return(invite, nil)
		err := svc.AcceptInvite(ctx, 1, 3)
		require.Error(t, err)
	})

	t.Run("invite already processed", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		invite := &model.TeamInvite{PlayerID: 2, Status: model.TeamInviteStatusAccepted}
		invite.ID = 1
		repo.On("GetInviteByID", ctx, uint64(1)).Return(invite, nil)
		err := svc.AcceptInvite(ctx, 1, 2)
		require.Error(t, err)
	})

	t.Run("invite not found", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		repo.On("GetInviteByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound)
		err := svc.AcceptInvite(ctx, 999, 2)
		require.Error(t, err)
	})
}

func TestTeamService_RejectInvite(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		invite := &model.TeamInvite{PlayerID: 2, Status: model.TeamInviteStatusPending}
		invite.ID = 1
		repo.On("GetInviteByID", ctx, uint64(1)).Return(invite, nil)
		repo.On("UpdateInvite", ctx, mock.AnythingOfType("*model.TeamInvite")).Return(nil)
		err := svc.RejectInvite(ctx, 1, 2)
		require.NoError(t, err)
	})

	t.Run("not the invitee", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		invite := &model.TeamInvite{PlayerID: 2, Status: model.TeamInviteStatusPending}
		invite.ID = 1
		repo.On("GetInviteByID", ctx, uint64(1)).Return(invite, nil)
		err := svc.RejectInvite(ctx, 1, 3)
		require.Error(t, err)
	})

	t.Run("invite already processed", func(t *testing.T) {
		repo := &MockTeamRepository{}
		svc := NewTeamService(repo)
		invite := &model.TeamInvite{PlayerID: 2, Status: model.TeamInviteStatusAccepted}
		invite.ID = 1
		repo.On("GetInviteByID", ctx, uint64(1)).Return(invite, nil)
		err := svc.RejectInvite(ctx, 1, 2)
		require.Error(t, err)
	})
}

// ============================================================================
// Additional Tests for Coverage Improvement
// ============================================================================

func TestTeamService_CreateTeam_CreateError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	repo.On("GetActiveMemberByPlayer", ctx, uint64(1)).Return(nil, repository.ErrNotFound)
	repo.On("Create", ctx, mock.AnythingOfType("*model.Team")).Return(errors.New("db error"))

	team := &model.Team{Name: "Test Team"}
	err := svc.CreateTeam(ctx, team, 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "创建团队失败")
}

func TestTeamService_CreateTeam_CreateMemberError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	repo.On("GetActiveMemberByPlayer", ctx, uint64(1)).Return(nil, repository.ErrNotFound)
	repo.On("Create", ctx, mock.AnythingOfType("*model.Team")).Return(nil)
	repo.On("CreateMember", ctx, mock.AnythingOfType("*model.TeamMember")).Return(errors.New("db error"))

	team := &model.Team{Name: "Test Team"}
	err := svc.CreateTeam(ctx, team, 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "创建队长成员记录失败")
}

func TestTeamService_UpdateTeamStatus_RepoError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	team := &model.Team{Status: model.TeamStatusActive, CurrentOrderID: nil}
	team.ID = 1
	repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
	repo.On("UpdateStatus", ctx, uint64(1), model.TeamStatusBusy).Return(errors.New("db error"))

	err := svc.UpdateTeamStatus(ctx, 1, model.TeamStatusBusy)
	require.Error(t, err)
}

func TestTeamService_AddMember_GetByIDError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	repo.On("GetByID", ctx, uint64(1)).Return(nil, errors.New("db error"))

	err := svc.AddMember(ctx, 1, 2)
	require.Error(t, err)
}

func TestTeamService_AddMember_CheckPlayerError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	team := &model.Team{MemberCount: 2, MaxMembers: 5}
	team.ID = 1
	repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
	repo.On("GetActiveMemberByPlayer", ctx, uint64(2)).Return(nil, errors.New("db error"))

	err := svc.AddMember(ctx, 1, 2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "检查陪玩师团队状态失败")
}

func TestTeamService_AddMember_GetMemberError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	team := &model.Team{MemberCount: 2, MaxMembers: 5}
	team.ID = 1
	repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
	repo.On("GetActiveMemberByPlayer", ctx, uint64(2)).Return(nil, repository.ErrNotFound)
	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(2)).Return(nil, errors.New("db error"))

	err := svc.AddMember(ctx, 1, 2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "检查成员记录失败")
}

func TestTeamService_AddMember_ReactivateError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	team := &model.Team{MemberCount: 2, MaxMembers: 5}
	team.ID = 1
	formerMember := &model.TeamMember{TeamID: 1, PlayerID: 2, Status: model.TeamMemberStatusLeft}
	repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
	repo.On("GetActiveMemberByPlayer", ctx, uint64(2)).Return(nil, repository.ErrNotFound)
	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(2)).Return(formerMember, nil)
	repo.On("UpdateMember", ctx, mock.AnythingOfType("*model.TeamMember")).Return(errors.New("db error"))

	err := svc.AddMember(ctx, 1, 2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "重新激活成员失败")
}

func TestTeamService_AddMember_CreateMemberError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	team := &model.Team{MemberCount: 2, MaxMembers: 5}
	team.ID = 1
	repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
	repo.On("GetActiveMemberByPlayer", ctx, uint64(2)).Return(nil, repository.ErrNotFound)
	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(2)).Return(nil, repository.ErrNotFound)
	repo.On("CreateMember", ctx, mock.AnythingOfType("*model.TeamMember")).Return(errors.New("db error"))

	err := svc.AddMember(ctx, 1, 2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "创建成员失败")
}

func TestTeamService_RemoveMember_GetMemberError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	team := &model.Team{CurrentOrderID: nil}
	team.ID = 1
	repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(2)).Return(nil, errors.New("db error"))

	err := svc.RemoveMember(ctx, 1, 2, false)
	require.Error(t, err)
}

func TestTeamService_RemoveMember_LeaderGetNextError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	team := &model.Team{CurrentOrderID: nil}
	team.ID = 1
	leader := &model.TeamMember{TeamID: 1, PlayerID: 1, Role: model.TeamMemberRoleLeader, Status: model.TeamMemberStatusActive}
	repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(leader, nil)
	repo.On("GetNextLeader", ctx, uint64(1), uint64(1)).Return(nil, errors.New("db error"))

	err := svc.RemoveMember(ctx, 1, 1, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "查找下一个队长失败")
}

func TestTeamService_RemoveMember_TransferLeaderError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	team := &model.Team{CurrentOrderID: nil}
	team.ID = 1
	leader := &model.TeamMember{TeamID: 1, PlayerID: 1, Role: model.TeamMemberRoleLeader, Status: model.TeamMemberStatusActive}
	nextLeader := &model.TeamMember{TeamID: 1, PlayerID: 2, Role: model.TeamMemberRoleMember, Status: model.TeamMemberStatusActive}
	repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(leader, nil)
	repo.On("GetNextLeader", ctx, uint64(1), uint64(1)).Return(nextLeader, nil)
	repo.On("UpdateMember", ctx, mock.AnythingOfType("*model.TeamMember")).Return(errors.New("db error")).Once()

	err := svc.RemoveMember(ctx, 1, 1, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "转让队长失败")
}

func TestTeamService_RemoveMember_UpdateLeaderError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	team := &model.Team{CurrentOrderID: nil}
	team.ID = 1
	leader := &model.TeamMember{TeamID: 1, PlayerID: 1, Role: model.TeamMemberRoleLeader, Status: model.TeamMemberStatusActive}
	nextLeader := &model.TeamMember{TeamID: 1, PlayerID: 2, Role: model.TeamMemberRoleMember, Status: model.TeamMemberStatusActive}
	repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(leader, nil)
	repo.On("GetNextLeader", ctx, uint64(1), uint64(1)).Return(nextLeader, nil)
	repo.On("UpdateMember", ctx, mock.AnythingOfType("*model.TeamMember")).Return(nil).Once()
	repo.On("UpdateLeader", ctx, uint64(1), uint64(2)).Return(errors.New("db error"))

	err := svc.RemoveMember(ctx, 1, 1, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "更新团队队长失败")
}

func TestTeamService_RemoveMember_UpdateMemberError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	team := &model.Team{CurrentOrderID: nil}
	team.ID = 1
	member := &model.TeamMember{TeamID: 1, PlayerID: 2, Role: model.TeamMemberRoleMember, Status: model.TeamMemberStatusActive}
	repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(2)).Return(member, nil)
	repo.On("UpdateMember", ctx, mock.AnythingOfType("*model.TeamMember")).Return(errors.New("db error"))

	err := svc.RemoveMember(ctx, 1, 2, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "更新成员状态失败")
}

func TestTeamService_KickMember_GetLeaderError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(2)).Return(nil, errors.New("db error"))

	err := svc.KickMember(ctx, 1, 2, 3)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "获取队长信息失败")
}

func TestTeamService_TransferLeader_GetCurrentLeaderError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	team := &model.Team{LeaderID: 1}
	team.ID = 1
	repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(nil, errors.New("db error"))

	err := svc.TransferLeader(ctx, 1, 1, 2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "获取当前队长信息失败")
}

func TestTeamService_TransferLeader_GetNewLeaderError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	team := &model.Team{LeaderID: 1}
	team.ID = 1
	currentLeader := &model.TeamMember{TeamID: 1, PlayerID: 1, Role: model.TeamMemberRoleLeader}
	repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(currentLeader, nil)
	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(2)).Return(nil, errors.New("db error"))

	err := svc.TransferLeader(ctx, 1, 1, 2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "获取新队长信息失败")
}

func TestTeamService_TransferLeader_UpdateCurrentLeaderError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	team := &model.Team{LeaderID: 1}
	team.ID = 1
	currentLeader := &model.TeamMember{TeamID: 1, PlayerID: 1, Role: model.TeamMemberRoleLeader}
	newLeader := &model.TeamMember{TeamID: 1, PlayerID: 2, Role: model.TeamMemberRoleMember, Status: model.TeamMemberStatusActive}
	repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(currentLeader, nil)
	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(2)).Return(newLeader, nil)
	repo.On("UpdateMember", ctx, mock.AnythingOfType("*model.TeamMember")).Return(errors.New("db error")).Once()

	err := svc.TransferLeader(ctx, 1, 1, 2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "更新原队长角色失败")
}

func TestTeamService_TransferLeader_UpdateNewLeaderError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	team := &model.Team{LeaderID: 1}
	team.ID = 1
	currentLeader := &model.TeamMember{TeamID: 1, PlayerID: 1, Role: model.TeamMemberRoleLeader}
	newLeader := &model.TeamMember{TeamID: 1, PlayerID: 2, Role: model.TeamMemberRoleMember, Status: model.TeamMemberStatusActive}
	repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(currentLeader, nil)
	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(2)).Return(newLeader, nil)
	repo.On("UpdateMember", ctx, mock.AnythingOfType("*model.TeamMember")).Return(nil).Once()
	repo.On("UpdateMember", ctx, mock.AnythingOfType("*model.TeamMember")).Return(errors.New("db error")).Once()

	err := svc.TransferLeader(ctx, 1, 1, 2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "更新新队长角色失败")
}

func TestTeamService_CreateInvite_GetTeamError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	repo.On("GetByID", ctx, uint64(1)).Return(nil, errors.New("db error"))

	err := svc.CreateInvite(ctx, 1, 1, 2, "Join us!")
	require.Error(t, err)
}

func TestTeamService_CreateInvite_GetInviterError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	team := &model.Team{MemberCount: 2, MaxMembers: 5}
	team.ID = 1
	repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(nil, errors.New("db error"))

	err := svc.CreateInvite(ctx, 1, 1, 2, "Join us!")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "邀请者不是团队成员")
}

func TestTeamService_CreateInvite_CheckTargetError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	team := &model.Team{MemberCount: 2, MaxMembers: 5}
	team.ID = 1
	inviter := &model.TeamMember{TeamID: 1, PlayerID: 1, Status: model.TeamMemberStatusActive}
	repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(inviter, nil)
	repo.On("GetActiveMemberByPlayer", ctx, uint64(2)).Return(nil, errors.New("db error"))

	err := svc.CreateInvite(ctx, 1, 1, 2, "Join us!")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "检查目标陪玩师团队状态失败")
}

func TestTeamService_CreateInvite_CheckPendingError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	team := &model.Team{MemberCount: 2, MaxMembers: 5}
	team.ID = 1
	inviter := &model.TeamMember{TeamID: 1, PlayerID: 1, Status: model.TeamMemberStatusActive}
	repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(1)).Return(inviter, nil)
	repo.On("GetActiveMemberByPlayer", ctx, uint64(2)).Return(nil, repository.ErrNotFound)
	repo.On("GetPendingInvite", ctx, uint64(1), uint64(2)).Return(nil, errors.New("db error"))

	err := svc.CreateInvite(ctx, 1, 1, 2, "Join us!")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "检查邀请状态失败")
}

func TestTeamService_AcceptInvite_Success(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	invite := &model.TeamInvite{
		TeamID:   1,
		PlayerID: 2,
		Status:   model.TeamInviteStatusPending,
		ExpireAt: time.Now().Add(24 * time.Hour),
	}
	invite.ID = 1
	team := &model.Team{MemberCount: 2, MaxMembers: 5}
	team.ID = 1

	repo.On("GetInviteByID", ctx, uint64(1)).Return(invite, nil)
	repo.On("GetByID", ctx, uint64(1)).Return(team, nil)
	repo.On("GetActiveMemberByPlayer", ctx, uint64(2)).Return(nil, repository.ErrNotFound)
	repo.On("GetMemberByTeamAndPlayer", ctx, uint64(1), uint64(2)).Return(nil, repository.ErrNotFound)
	repo.On("CreateMember", ctx, mock.AnythingOfType("*model.TeamMember")).Return(nil)
	repo.On("UpdateMemberCount", ctx, uint64(1), 1).Return(nil)
	repo.On("UpdateInvite", ctx, mock.AnythingOfType("*model.TeamInvite")).Return(nil)

	err := svc.AcceptInvite(ctx, 1, 2)
	require.NoError(t, err)
}

func TestTeamService_AcceptInvite_Expired(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	invite := &model.TeamInvite{
		TeamID:   1,
		PlayerID: 2,
		Status:   model.TeamInviteStatusPending,
		ExpireAt: time.Now().Add(-24 * time.Hour), // Expired
	}
	invite.ID = 1

	repo.On("GetInviteByID", ctx, uint64(1)).Return(invite, nil)
	repo.On("UpdateInvite", ctx, mock.AnythingOfType("*model.TeamInvite")).Return(nil)

	err := svc.AcceptInvite(ctx, 1, 2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "已过期")
}

func TestTeamService_AcceptInvite_AddMemberError(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	invite := &model.TeamInvite{
		TeamID:   1,
		PlayerID: 2,
		Status:   model.TeamInviteStatusPending,
		ExpireAt: time.Now().Add(24 * time.Hour),
	}
	invite.ID = 1
	team := &model.Team{MemberCount: 5, MaxMembers: 5} // Full team
	team.ID = 1

	repo.On("GetInviteByID", ctx, uint64(1)).Return(invite, nil)
	repo.On("GetByID", ctx, uint64(1)).Return(team, nil)

	err := svc.AcceptInvite(ctx, 1, 2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "团队已满")
}

func TestTeamService_RejectInvite_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := &MockTeamRepository{}
	svc := NewTeamService(repo)

	repo.On("GetInviteByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	err := svc.RejectInvite(ctx, 999, 2)
	require.Error(t, err)
}
