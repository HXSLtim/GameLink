package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository/team"
)

// ============================================================================
// Team CRUD Tests
// ============================================================================

func TestTeamRepository_CreateTeam(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := team.NewTeamRepository(db)
	ctx := context.Background()

	// Create leader player
	user := CreateUniqueTestUser(t, db, "team_leader")
	player := CreateTestPlayer(t, db, user)

	teamObj := &model.Team{
		Name:            "Test Team",
		Description:     "Test team description",
		LeaderID:        player.ID,
		Status:          model.TeamStatusActive,
		MaxMembers:      5,
		MemberCount:     1,
		IncomeShareType: "equal",
	}

	err := repo.Create(ctx, teamObj)
	require.NoError(t, err)
	assert.NotZero(t, teamObj.ID)

	// Verify
	got, err := repo.GetByID(ctx, teamObj.ID)
	require.NoError(t, err)
	assert.Equal(t, "Test Team", got.Name)
	assert.Equal(t, model.TeamStatusActive, got.Status)
}

func TestTeamRepository_UpdateTeam(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := team.NewTeamRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "team_update")
	player := CreateTestPlayer(t, db, user)

	teamObj := &model.Team{
		Name:            "Update Test",
		LeaderID:        player.ID,
		Status:          model.TeamStatusActive,
		MaxMembers:      5,
		MemberCount:     1,
		IncomeShareType: "equal",
	}
	require.NoError(t, repo.Create(ctx, teamObj))

	// Update
	teamObj.Name = "Updated Name"
	teamObj.Description = "Updated description"
	teamObj.MaxMembers = 10
	err := repo.Update(ctx, teamObj)
	require.NoError(t, err)

	// Verify
	got, err := repo.GetByID(ctx, teamObj.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", got.Name)
	assert.Equal(t, 10, got.MaxMembers)
}

func TestTeamRepository_DeleteTeam(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := team.NewTeamRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "team_delete")
	player := CreateTestPlayer(t, db, user)

	teamObj := &model.Team{
		Name:            "Delete Test",
		LeaderID:        player.ID,
		Status:          model.TeamStatusActive,
		MaxMembers:      5,
		MemberCount:     1,
		IncomeShareType: "equal",
	}
	require.NoError(t, repo.Create(ctx, teamObj))

	err := repo.Delete(ctx, teamObj.ID)
	require.NoError(t, err)

	// Verify deleted (soft delete)
	_, err = repo.GetByID(ctx, teamObj.ID)
	assert.Error(t, err)
}

func TestTeamRepository_ListTeams(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := team.NewTeamRepository(db)
	ctx := context.Background()

	// Create multiple teams
	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, "team_list")
		player := CreateTestPlayer(t, db, user)
		teamObj := &model.Team{
			Name:            "List Test Team",
			LeaderID:        player.ID,
			Status:          model.TeamStatusActive,
			MaxMembers:      5,
			MemberCount:     1,
			IncomeShareType: "equal",
		}
		require.NoError(t, repo.Create(ctx, teamObj))
	}

	teams, total, err := repo.List(ctx, team.TeamListOptions{
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, teams, 3)
}

func TestTeamRepository_UpdateStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := team.NewTeamRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "team_status")
	player := CreateTestPlayer(t, db, user)

	teamObj := &model.Team{
		Name:            "Status Test",
		LeaderID:        player.ID,
		Status:          model.TeamStatusActive,
		MaxMembers:      5,
		MemberCount:     1,
		IncomeShareType: "equal",
	}
	require.NoError(t, repo.Create(ctx, teamObj))

	// Update to busy
	err := repo.UpdateStatus(ctx, teamObj.ID, model.TeamStatusBusy)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, teamObj.ID)
	require.NoError(t, err)
	assert.Equal(t, model.TeamStatusBusy, got.Status)
}

// ============================================================================
// Team Member Tests
// ============================================================================

func TestTeamRepository_CreateMember(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := team.NewTeamRepository(db)
	ctx := context.Background()

	// Create team with leader
	leaderUser := CreateUniqueTestUser(t, db, "team_member_leader")
	leader := CreateTestPlayer(t, db, leaderUser)

	teamObj := &model.Team{
		Name:            "Member Test",
		LeaderID:        leader.ID,
		Status:          model.TeamStatusActive,
		MaxMembers:      5,
		MemberCount:     1,
		IncomeShareType: "equal",
	}
	require.NoError(t, repo.Create(ctx, teamObj))

	// Add leader as member
	leaderMember := &model.TeamMember{
		TeamID:   teamObj.ID,
		PlayerID: leader.ID,
		Role:     model.TeamMemberRoleLeader,
		Status:   model.TeamMemberStatusActive,
		JoinedAt: time.Now(),
	}
	err := repo.CreateMember(ctx, leaderMember)
	require.NoError(t, err)
	assert.NotZero(t, leaderMember.ID)

	// Add another member
	memberUser := CreateUniqueTestUser(t, db, "team_member")
	member := CreateTestPlayer(t, db, memberUser)

	teamMember := &model.TeamMember{
		TeamID:   teamObj.ID,
		PlayerID: member.ID,
		Role:     model.TeamMemberRoleMember,
		Status:   model.TeamMemberStatusActive,
		JoinedAt: time.Now(),
	}
	err = repo.CreateMember(ctx, teamMember)
	require.NoError(t, err)
	assert.NotZero(t, teamMember.ID)
}

func TestTeamRepository_GetActiveMembers(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := team.NewTeamRepository(db)
	ctx := context.Background()

	leaderUser := CreateUniqueTestUser(t, db, "active_members_leader")
	leader := CreateTestPlayer(t, db, leaderUser)

	teamObj := &model.Team{
		Name:            "Active Members Test",
		LeaderID:        leader.ID,
		Status:          model.TeamStatusActive,
		MaxMembers:      5,
		MemberCount:     3,
		IncomeShareType: "equal",
	}
	require.NoError(t, repo.Create(ctx, teamObj))

	// Add leader
	require.NoError(t, repo.CreateMember(ctx, &model.TeamMember{
		TeamID:   teamObj.ID,
		PlayerID: leader.ID,
		Role:     model.TeamMemberRoleLeader,
		Status:   model.TeamMemberStatusActive,
		JoinedAt: time.Now(),
	}))

	// Add 2 more active members
	for i := 0; i < 2; i++ {
		memberUser := CreateUniqueTestUser(t, db, "active_member")
		member := CreateTestPlayer(t, db, memberUser)
		require.NoError(t, repo.CreateMember(ctx, &model.TeamMember{
			TeamID:   teamObj.ID,
			PlayerID: member.ID,
			Role:     model.TeamMemberRoleMember,
			Status:   model.TeamMemberStatusActive,
			JoinedAt: time.Now(),
		}))
	}

	members, err := repo.GetActiveMembers(ctx, teamObj.ID)
	require.NoError(t, err)
	assert.Len(t, members, 3)
}

func TestTeamRepository_UpdateMember(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := team.NewTeamRepository(db)
	ctx := context.Background()

	leaderUser := CreateUniqueTestUser(t, db, "update_member_leader")
	leader := CreateTestPlayer(t, db, leaderUser)

	teamObj := &model.Team{
		Name:            "Update Member Test",
		LeaderID:        leader.ID,
		Status:          model.TeamStatusActive,
		MaxMembers:      5,
		MemberCount:     1,
		IncomeShareType: "equal",
	}
	require.NoError(t, repo.Create(ctx, teamObj))

	member := &model.TeamMember{
		TeamID:   teamObj.ID,
		PlayerID: leader.ID,
		Role:     model.TeamMemberRoleMember,
		Status:   model.TeamMemberStatusActive,
		JoinedAt: time.Now(),
	}
	require.NoError(t, repo.CreateMember(ctx, member))

	// Update to leader
	member.Role = model.TeamMemberRoleLeader
	err := repo.UpdateMember(ctx, member)
	require.NoError(t, err)

	got, err := repo.GetMemberByID(ctx, member.ID)
	require.NoError(t, err)
	assert.Equal(t, model.TeamMemberRoleLeader, got.Role)
}

func TestTeamRepository_MemberLeave(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := team.NewTeamRepository(db)
	ctx := context.Background()

	leaderUser := CreateUniqueTestUser(t, db, "leave_leader")
	leader := CreateTestPlayer(t, db, leaderUser)

	teamObj := &model.Team{
		Name:            "Leave Test",
		LeaderID:        leader.ID,
		Status:          model.TeamStatusActive,
		MaxMembers:      5,
		MemberCount:     2,
		IncomeShareType: "equal",
	}
	require.NoError(t, repo.Create(ctx, teamObj))

	memberUser := CreateUniqueTestUser(t, db, "leave_member")
	memberPlayer := CreateTestPlayer(t, db, memberUser)

	member := &model.TeamMember{
		TeamID:   teamObj.ID,
		PlayerID: memberPlayer.ID,
		Role:     model.TeamMemberRoleMember,
		Status:   model.TeamMemberStatusActive,
		JoinedAt: time.Now(),
	}
	require.NoError(t, repo.CreateMember(ctx, member))

	// Member leaves
	now := time.Now()
	member.Status = model.TeamMemberStatusLeft
	member.LeftAt = &now
	err := repo.UpdateMember(ctx, member)
	require.NoError(t, err)

	got, err := repo.GetMemberByID(ctx, member.ID)
	require.NoError(t, err)
	assert.Equal(t, model.TeamMemberStatusLeft, got.Status)
	assert.NotNil(t, got.LeftAt)
}

// ============================================================================
// Team Invite Tests
// ============================================================================

func TestTeamRepository_CreateInvite(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := team.NewTeamRepository(db)
	ctx := context.Background()

	leaderUser := CreateUniqueTestUser(t, db, "invite_leader")
	leader := CreateTestPlayer(t, db, leaderUser)

	teamObj := &model.Team{
		Name:            "Invite Test",
		LeaderID:        leader.ID,
		Status:          model.TeamStatusActive,
		MaxMembers:      5,
		MemberCount:     1,
		IncomeShareType: "equal",
	}
	require.NoError(t, repo.Create(ctx, teamObj))

	inviteeUser := CreateUniqueTestUser(t, db, "invitee")
	invitee := CreateTestPlayer(t, db, inviteeUser)

	invite := &model.TeamInvite{
		TeamID:    teamObj.ID,
		PlayerID:  invitee.ID,
		InviterID: leader.ID,
		Status:    model.TeamInviteStatusPending,
		ExpireAt:  time.Now().Add(24 * time.Hour),
		Message:   "Join our team!",
	}

	err := repo.CreateInvite(ctx, invite)
	require.NoError(t, err)
	assert.NotZero(t, invite.ID)
}

func TestTeamRepository_AcceptInvite(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := team.NewTeamRepository(db)
	ctx := context.Background()

	leaderUser := CreateUniqueTestUser(t, db, "accept_leader")
	leader := CreateTestPlayer(t, db, leaderUser)

	teamObj := &model.Team{
		Name:            "Accept Invite Test",
		LeaderID:        leader.ID,
		Status:          model.TeamStatusActive,
		MaxMembers:      5,
		MemberCount:     1,
		IncomeShareType: "equal",
	}
	require.NoError(t, repo.Create(ctx, teamObj))

	inviteeUser := CreateUniqueTestUser(t, db, "accept_invitee")
	invitee := CreateTestPlayer(t, db, inviteeUser)

	invite := &model.TeamInvite{
		TeamID:    teamObj.ID,
		PlayerID:  invitee.ID,
		InviterID: leader.ID,
		Status:    model.TeamInviteStatusPending,
		ExpireAt:  time.Now().Add(24 * time.Hour),
	}
	require.NoError(t, repo.CreateInvite(ctx, invite))

	// Accept invite
	invite.Status = model.TeamInviteStatusAccepted
	err := repo.UpdateInvite(ctx, invite)
	require.NoError(t, err)

	got, err := repo.GetInviteByID(ctx, invite.ID)
	require.NoError(t, err)
	assert.Equal(t, model.TeamInviteStatusAccepted, got.Status)
}

func TestTeamRepository_RejectInvite(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := team.NewTeamRepository(db)
	ctx := context.Background()

	leaderUser := CreateUniqueTestUser(t, db, "reject_leader")
	leader := CreateTestPlayer(t, db, leaderUser)

	teamObj := &model.Team{
		Name:            "Reject Invite Test",
		LeaderID:        leader.ID,
		Status:          model.TeamStatusActive,
		MaxMembers:      5,
		MemberCount:     1,
		IncomeShareType: "equal",
	}
	require.NoError(t, repo.Create(ctx, teamObj))

	inviteeUser := CreateUniqueTestUser(t, db, "reject_invitee")
	invitee := CreateTestPlayer(t, db, inviteeUser)

	invite := &model.TeamInvite{
		TeamID:    teamObj.ID,
		PlayerID:  invitee.ID,
		InviterID: leader.ID,
		Status:    model.TeamInviteStatusPending,
		ExpireAt:  time.Now().Add(24 * time.Hour),
	}
	require.NoError(t, repo.CreateInvite(ctx, invite))

	// Reject invite
	invite.Status = model.TeamInviteStatusRejected
	err := repo.UpdateInvite(ctx, invite)
	require.NoError(t, err)

	got, err := repo.GetInviteByID(ctx, invite.ID)
	require.NoError(t, err)
	assert.Equal(t, model.TeamInviteStatusRejected, got.Status)
}

func TestTeamRepository_ExpireInvites(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := team.NewTeamRepository(db)
	ctx := context.Background()

	leaderUser := CreateUniqueTestUser(t, db, "expire_leader")
	leader := CreateTestPlayer(t, db, leaderUser)

	teamObj := &model.Team{
		Name:            "Expire Invite Test",
		LeaderID:        leader.ID,
		Status:          model.TeamStatusActive,
		MaxMembers:      5,
		MemberCount:     1,
		IncomeShareType: "equal",
	}
	require.NoError(t, repo.Create(ctx, teamObj))

	inviteeUser := CreateUniqueTestUser(t, db, "expire_invitee")
	invitee := CreateTestPlayer(t, db, inviteeUser)

	// Create expired invite
	invite := &model.TeamInvite{
		TeamID:    teamObj.ID,
		PlayerID:  invitee.ID,
		InviterID: leader.ID,
		Status:    model.TeamInviteStatusPending,
		ExpireAt:  time.Now().Add(-time.Hour), // Already expired
	}
	require.NoError(t, repo.CreateInvite(ctx, invite))

	// Run expire job
	affected, err := repo.ExpireInvites(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, affected, int64(1))

	got, err := repo.GetInviteByID(ctx, invite.ID)
	require.NoError(t, err)
	assert.Equal(t, model.TeamInviteStatusExpired, got.Status)
}

// ============================================================================
// Team Stats Tests
// ============================================================================

func TestTeamRepository_GetTeamStats(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := team.NewTeamRepository(db)
	ctx := context.Background()

	// Create some teams
	for i := 0; i < 2; i++ {
		user := CreateUniqueTestUser(t, db, "stats_team")
		player := CreateTestPlayer(t, db, user)
		teamObj := &model.Team{
			Name:            "Stats Test Team",
			LeaderID:        player.ID,
			Status:          model.TeamStatusActive,
			MaxMembers:      5,
			MemberCount:     1,
			IncomeShareType: "equal",
		}
		require.NoError(t, repo.Create(ctx, teamObj))
	}

	stats, err := repo.GetTeamStats(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, stats.TotalTeams, int64(2))
	assert.GreaterOrEqual(t, stats.ActiveTeams, int64(2))
}

func TestTeamRepository_IncrementStats(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := team.NewTeamRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "increment_stats")
	player := CreateTestPlayer(t, db, user)

	teamObj := &model.Team{
		Name:            "Increment Stats Test",
		LeaderID:        player.ID,
		Status:          model.TeamStatusActive,
		MaxMembers:      5,
		MemberCount:     1,
		IncomeShareType: "equal",
	}
	require.NoError(t, repo.Create(ctx, teamObj))

	// Increment stats
	err := repo.IncrementStats(ctx, teamObj.ID, 1, 10000)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, teamObj.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, got.TotalOrderCount)
	assert.Equal(t, int64(10000), got.TotalIncomeCents)
}
