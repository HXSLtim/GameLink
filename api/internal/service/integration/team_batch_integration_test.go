// Package integration provides integration tests for team batch operations.
package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	teamrepo "gamelink/internal/repository/team"
	teamservice "gamelink/internal/service/team"
)

// ============================================================================
// Test Helpers for Team Batch Operations
// ============================================================================

// createTestTeamForBatch creates a test team with a leader for batch tests.
func createTestTeamForBatch(t *testing.T, db *gorm.DB, name string, leaderID uint64) *model.Team {
	t.Helper()
	teamObj := &model.Team{
		Base:            model.Base{ExtJSON: "{}"},
		Name:            name,
		LeaderID:        leaderID,
		Status:          model.TeamStatusActive,
		MaxMembers:      5,
		MemberCount:     1,
		IncomeShareType: "equal",
	}
	if err := db.Create(teamObj).Error; err != nil {
		t.Fatalf("Failed to create test team: %v", err)
	}

	// Create leader member record
	member := &model.TeamMember{
		Base:      model.Base{ExtJSON: "{}"},
		TeamID:    teamObj.ID,
		PlayerID:  leaderID,
		Role:      model.TeamMemberRoleLeader,
		Status:    model.TeamMemberStatusActive,
		JoinedAt:  time.Now(),
		SortOrder: 0,
	}
	if err := db.Create(member).Error; err != nil {
		t.Fatalf("Failed to create team leader member: %v", err)
	}

	return teamObj
}

// createTestTeamWithMembersForBatch creates a test team with multiple members for batch tests.
func createTestTeamWithMembersForBatch(t *testing.T, db *gorm.DB, name string, leaderID uint64, memberCount int) (*model.Team, []uint64) {
	t.Helper()
	teamObj := createTestTeamForBatch(t, db, name, leaderID)

	var memberIDs []uint64
	for i := 0; i < memberCount; i++ {
		user := CreateUniqueTestUser(t, db, name+"_member"+string(rune('1'+i)))
		player := CreateTestPlayer(t, db, user)

		member := &model.TeamMember{
			Base:      model.Base{ExtJSON: "{}"},
			TeamID:    teamObj.ID,
			PlayerID:  player.ID,
			Role:      model.TeamMemberRoleMember,
			Status:    model.TeamMemberStatusActive,
			JoinedAt:  time.Now(),
			SortOrder: i + 1,
		}
		if err := db.Create(member).Error; err != nil {
			t.Fatalf("Failed to create team member: %v", err)
		}
		memberIDs = append(memberIDs, player.ID)
	}

	// Update member count
	teamObj.MemberCount = memberCount + 1
	db.Save(teamObj)

	return teamObj, memberIDs
}

// ============================================================================
// BatchDeleteTeams Tests
// ============================================================================

func TestTeamService_BatchDeleteTeams_AllSuccess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create test teams
	var teamIDs []uint64
	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, "batch_delete_all_success")
		player := CreateTestPlayer(t, db, user)
		teamObj := createTestTeamForBatch(t, db, "Test Team "+string(rune('1'+i)), player.ID)
		teamIDs = append(teamIDs, teamObj.ID)
	}

	// Batch delete all teams
	result, err := svc.BatchDeleteTeams(ctx, teamIDs)
	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Empty(t, result.FailedIDs)
	assert.Empty(t, result.Errors)

	// Verify all teams are deleted
	for _, id := range teamIDs {
		_, err := repo.GetByID(ctx, id)
		assert.Error(t, err)
	}
}

func TestTeamService_BatchDeleteTeams_WithActiveOrder(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create teams without active orders
	var validTeamIDs []uint64
	for i := 0; i < 2; i++ {
		user := CreateUniqueTestUser(t, db, "batch_delete_valid")
		player := CreateTestPlayer(t, db, user)
		teamObj := createTestTeamForBatch(t, db, "Valid Team", player.ID)
		validTeamIDs = append(validTeamIDs, teamObj.ID)
	}

	// Create team with active order (should fail)
	user := CreateUniqueTestUser(t, db, "batch_delete_with_order")
	player := CreateTestPlayer(t, db, user)
	order := CreateTestOrder(t, db, user, player, model.OrderStatusInProgress)

	busyTeam := createTestTeamForBatch(t, db, "Busy Team", player.ID)
	orderID := order.ID
	busyTeam.CurrentOrderID = &orderID
	db.Save(busyTeam)

	// Mix with valid team IDs
	allTeamIDs := append(validTeamIDs, busyTeam.ID)

	// Batch delete
	result, err := svc.BatchDeleteTeams(ctx, allTeamIDs)
	require.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Contains(t, result.FailedIDs, busyTeam.ID)
	assert.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0], "团队有进行中的订单")

	// Verify valid teams are deleted
	for _, id := range validTeamIDs {
		_, err := repo.GetByID(ctx, id)
		assert.Error(t, err)
	}

	// Verify busy team still exists
	team, err := repo.GetByID(ctx, busyTeam.ID)
	require.NoError(t, err)
	assert.NotNil(t, team.CurrentOrderID)
}

func TestTeamService_BatchDeleteTeams_NonExistentTeams(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create valid team
	user := CreateUniqueTestUser(t, db, "batch_delete_nonexistent")
	player := CreateTestPlayer(t, db, user)
	teamObj := CreateTestTeam(t, db, "Valid Team", player.ID)

	// Mix with non-existent team IDs
	teamIDs := []uint64{teamObj.ID, 99999, 88888}

	// Batch delete
	result, err := svc.BatchDeleteTeams(ctx, teamIDs)
	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)
	assert.Len(t, result.FailedIDs, 2)

	// Verify valid team is deleted
	_, err = repo.GetByID(ctx, teamObj.ID)
	assert.Error(t, err)
}

func TestTeamService_BatchDeleteTeams_WithMembers(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create team with members (should still be deletable)
	leaderUser := CreateUniqueTestUser(t, db, "batch_delete_leader")
	leader := CreateTestPlayer(t, db, leaderUser)
	teamObj, _ := createTestTeamWithMembersForBatch(t, db, "Team With Members", leader.ID, 3)

	// Batch delete
	result, err := svc.BatchDeleteTeams(ctx, []uint64{teamObj.ID})
	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)

	// Verify team is deleted
	_, err = repo.GetByID(ctx, teamObj.ID)
	assert.Error(t, err)
}

func TestTeamService_BatchDeleteTeams_EmptyList(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Empty list should return error
	result, err := svc.BatchDeleteTeams(ctx, []uint64{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "团队ID列表不能为空")
}

// ============================================================================
// BatchUpdateTeamStatus Tests
// ============================================================================

func TestTeamService_BatchUpdateTeamStatus_ToActive(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create inactive teams
	var teamIDs []uint64
	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, "batch_activate")
		player := CreateTestPlayer(t, db, user)
		teamObj := createTestTeamForBatch(t, db, "Inactive Team", player.ID)
		teamObj.Status = model.TeamStatusInactive
		db.Save(teamObj)
		teamIDs = append(teamIDs, teamObj.ID)
	}

	// Batch activate
	result, err := svc.BatchUpdateTeamStatus(ctx, teamIDs, model.TeamStatusActive)
	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)

	// Verify status updates
	for _, id := range teamIDs {
		team, err := repo.GetByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.TeamStatusActive, team.Status)
	}
}

func TestTeamService_BatchUpdateTeamStatus_ToInactive(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create active teams
	var validTeamIDs []uint64
	for i := 0; i < 2; i++ {
		user := CreateUniqueTestUser(t, db, "batch_deactivate_valid")
		player := CreateTestPlayer(t, db, user)
		teamObj := createTestTeamForBatch(t, db, "Active Team", player.ID)
		validTeamIDs = append(validTeamIDs, teamObj.ID)
	}

	// Create team with active order (cannot deactivate)
	user := CreateUniqueTestUser(t, db, "batch_deactivate_busy")
	player := CreateTestPlayer(t, db, user)
	order := CreateTestOrder(t, db, user, player, model.OrderStatusInProgress)

	busyTeam := createTestTeamForBatch(t, db, "Busy Team", player.ID)
	orderID := order.ID
	busyTeam.CurrentOrderID = &orderID
	db.Save(busyTeam)

	// Mix team IDs
	allTeamIDs := append(validTeamIDs, busyTeam.ID)

	// Batch deactivate
	result, err := svc.BatchUpdateTeamStatus(ctx, allTeamIDs, model.TeamStatusInactive)
	require.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Contains(t, result.FailedIDs, busyTeam.ID)
	assert.Contains(t, result.Errors[0], "团队有进行中的订单")

	// Verify valid teams are deactivated
	for _, id := range validTeamIDs {
		team, err := repo.GetByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.TeamStatusInactive, team.Status)
	}

	// Verify busy team is still active
	busyTeamResult, err := repo.GetByID(ctx, busyTeam.ID)
	require.NoError(t, err)
	assert.Equal(t, model.TeamStatusActive, busyTeamResult.Status)
}

func TestTeamService_BatchUpdateTeamStatus_ToBusy(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create active teams
	var teamIDs []uint64
	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, "batch_busy")
		player := CreateTestPlayer(t, db, user)
		teamObj := createTestTeamForBatch(t, db, "Active Team", player.ID)
		teamIDs = append(teamIDs, teamObj.ID)
	}

	// Batch set to busy
	result, err := svc.BatchUpdateTeamStatus(ctx, teamIDs, model.TeamStatusBusy)
	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)

	// Verify status updates
	for _, id := range teamIDs {
		team, err := repo.GetByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.TeamStatusBusy, team.Status)
	}
}

func TestTeamService_BatchUpdateTeamStatus_NonExistentTeams(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create valid team
	user := CreateUniqueTestUser(t, db, "batch_status_nonexistent")
	player := CreateTestPlayer(t, db, user)
	teamObj := CreateTestTeam(t, db, "Valid Team", player.ID)

	// Mix with non-existent team IDs
	teamIDs := []uint64{teamObj.ID, 99999, 88888}

	// Batch update status
	result, err := svc.BatchUpdateTeamStatus(ctx, teamIDs, model.TeamStatusInactive)
	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)
	assert.Len(t, result.FailedIDs, 2)

	// Verify valid team status updated
	team, err := repo.GetByID(ctx, teamObj.ID)
	require.NoError(t, err)
	assert.Equal(t, model.TeamStatusInactive, team.Status)
}

func TestTeamService_BatchUpdateTeamStatus_InvalidStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create team
	user := CreateUniqueTestUser(t, db, "batch_status_invalid")
	player := CreateTestPlayer(t, db, user)
	teamObj := CreateTestTeam(t, db, "Test Team", player.ID)

	// Try to update with invalid status
	result, err := svc.BatchUpdateTeamStatus(ctx, []uint64{teamObj.ID}, "invalid_status")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "无效的团队状态")

	// Verify team status unchanged
	team, err := repo.GetByID(ctx, teamObj.ID)
	require.NoError(t, err)
	assert.Equal(t, model.TeamStatusActive, team.Status)
}

func TestTeamService_BatchUpdateTeamStatus_EmptyList(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Empty list should return error
	result, err := svc.BatchUpdateTeamStatus(ctx, []uint64{}, model.TeamStatusActive)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "团队ID列表不能为空")
}

// ============================================================================
// BatchAddMembers Tests
// ============================================================================

func TestTeamService_BatchAddMembers_AllSuccess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create team
	leaderUser := CreateUniqueTestUser(t, db, "batch_add_all_success")
	leader := CreateTestPlayer(t, db, leaderUser)
	teamObj := createTestTeamForBatch(t, db, "Test Team", leader.ID)

	// Create players to add
	var playerIDs []uint64
	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, "batch_add_member")
		player := CreateTestPlayer(t, db, user)
		playerIDs = append(playerIDs, player.ID)
	}

	// Batch add members
	result, err := svc.BatchAddMembers(ctx, teamObj.ID, playerIDs)
	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Empty(t, result.FailedPlayerIDs)
	assert.Empty(t, result.Errors)

	// Verify team member count
	team, err := repo.GetByID(ctx, teamObj.ID)
	require.NoError(t, err)
	assert.Equal(t, 4, team.MemberCount) // 1 leader + 3 new members

	// Verify all members are active
	members, err := repo.GetActiveMembers(ctx, teamObj.ID)
	require.NoError(t, err)
	assert.Len(t, members, 4)
}

func TestTeamService_BatchAddMembers_SomeAlreadyMembers(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create team with existing members
	leaderUser := CreateUniqueTestUser(t, db, "batch_add_existing")
	leader := CreateTestPlayer(t, db, leaderUser)
	teamObj, existingMemberIDs := createTestTeamWithMembersForBatch(t, db, "Test Team", leader.ID, 2)

	// Try to add existing members again
	var playerIDs []uint64
	// Add one existing member
	playerIDs = append(playerIDs, existingMemberIDs[0])
	// Add one new member
	newUser := CreateUniqueTestUser(t, db, "batch_add_new")
	newPlayer := CreateTestPlayer(t, db, newUser)
	playerIDs = append(playerIDs, newPlayer.ID)

	// Batch add members
	result, err := svc.BatchAddMembers(ctx, teamObj.ID, playerIDs)
	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount) // Only new member added
	assert.Equal(t, 1, result.FailedCount)  // Existing member failed
	assert.Len(t, result.FailedPlayerIDs, 1)
	assert.Contains(t, result.Errors[0], "陪玩师已在其他团队中")

	// Verify team member count (unchanged for existing member, +1 for new)
	team, err := repo.GetByID(ctx, teamObj.ID)
	require.NoError(t, err)
	assert.Equal(t, 4, team.MemberCount) // 1 leader + 2 existing + 1 new
}

func TestTeamService_BatchAddMembers_NonExistentPlayers(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create team
	leaderUser := CreateUniqueTestUser(t, db, "batch_add_nonexistent_player")
	leader := CreateTestPlayer(t, db, leaderUser)
	teamObj := createTestTeamForBatch(t, db, "Test Team", leader.ID)

	// Create one valid player
	user := CreateUniqueTestUser(t, db, "batch_add_valid")
	player := CreateTestPlayer(t, db, user)

	// Mix with non-existent player IDs
	playerIDs := []uint64{player.ID, 99999, 88888}

	// Batch add members
	result, err := svc.BatchAddMembers(ctx, teamObj.ID, playerIDs)
	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)
	assert.Len(t, result.FailedPlayerIDs, 2)

	// Verify team member count
	team, err := repo.GetByID(ctx, teamObj.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, team.MemberCount) // 1 leader + 1 new member
}

func TestTeamService_BatchAddMembers_NonExistentTeam(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create players
	user := CreateUniqueTestUser(t, db, "batch_add_nonexistent_team")
	player := CreateTestPlayer(t, db, user)

	// Try to add to non-existent team
	result, err := svc.BatchAddMembers(ctx, 99999, []uint64{player.ID})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "团队不存在")
}

func TestTeamService_BatchAddMembers_TeamFull(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create team with max members
	leaderUser := CreateUniqueTestUser(t, db, "batch_add_full")
	leader := CreateTestPlayer(t, db, leaderUser)
	teamObj, _ := createTestTeamWithMembersForBatch(t, db, "Full Team", leader.ID, 4) // 1 leader + 4 members = 5 max

	// Try to add more players
	var playerIDs []uint64
	for i := 0; i < 2; i++ {
		user := CreateUniqueTestUser(t, db, "batch_add_overflow")
		player := CreateTestPlayer(t, db, user)
		playerIDs = append(playerIDs, player.ID)
	}

	// Batch add members (should fail as team is full)
	result, err := svc.BatchAddMembers(ctx, teamObj.ID, playerIDs)
	require.NoError(t, err)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)
	assert.Len(t, result.FailedPlayerIDs, 2)
	// Both should fail with "团队已满"
	for _, e := range result.Errors {
		assert.Contains(t, e, "团队已满")
	}

	// Verify team member count unchanged
	team, err := repo.GetByID(ctx, teamObj.ID)
	require.NoError(t, err)
	assert.Equal(t, 5, team.MemberCount) // Still at max
}

func TestTeamService_BatchAddMembers_EmptyList(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create team
	leaderUser := CreateUniqueTestUser(t, db, "batch_add_empty")
	leader := CreateTestPlayer(t, db, leaderUser)
	teamObj := createTestTeamForBatch(t, db, "Test Team", leader.ID)

	// Empty list should return error
	result, err := svc.BatchAddMembers(ctx, teamObj.ID, []uint64{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "陪玩师ID列表不能为空")
}

func TestTeamService_BatchAddMembers_MixedScenarios(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create team
	leaderUser := CreateUniqueTestUser(t, db, "batch_add_mixed")
	leader := CreateTestPlayer(t, db, leaderUser)
	teamObj, existingMemberIDs := createTestTeamWithMembersForBatch(t, db, "Test Team", leader.ID, 2)

	var playerIDs []uint64

	// Add 2 valid new players
	for i := 0; i < 2; i++ {
		user := CreateUniqueTestUser(t, db, "batch_add_mixed_new")
		player := CreateTestPlayer(t, db, user)
		playerIDs = append(playerIDs, player.ID)
	}

	// Add 1 existing member (should fail)
	playerIDs = append(playerIDs, existingMemberIDs[0])

	// Add 1 player from another team (should fail)
	anotherUser := CreateUniqueTestUser(t, db, "batch_add_mixed_another")
	anotherPlayer := CreateTestPlayer(t, db, anotherUser)
	createTestTeamForBatch(t, db, "Another Team", anotherPlayer.ID)
	playerIDs = append(playerIDs, anotherPlayer.ID)

	// Add 1 non-existent player (should fail)
	playerIDs = append(playerIDs, 99999)

	// Batch add members
	result, err := svc.BatchAddMembers(ctx, teamObj.ID, playerIDs)
	require.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 3, result.FailedCount)
	assert.Len(t, result.FailedPlayerIDs, 3)
	assert.Len(t, result.Errors, 3)

	// Verify team member count
	team, err := repo.GetByID(ctx, teamObj.ID)
	require.NoError(t, err)
	assert.Equal(t, 5, team.MemberCount) // 1 leader + 2 existing + 2 new
}

// ============================================================================
// BatchOperationResult Format Tests
// ============================================================================

func TestTeamService_BatchOperationResult_Format(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create teams
	var teamIDs []uint64
	for i := 0; i < 5; i++ {
		user := CreateUniqueTestUser(t, db, "batch_format")
		player := CreateTestPlayer(t, db, user)
		teamObj := createTestTeamForBatch(t, db, "Format Test Team", player.ID)
		teamIDs = append(teamIDs, teamObj.ID)
	}

	// Add non-existent IDs
	teamIDs = append(teamIDs, 99999, 88888)

	// Batch delete
	result, err := svc.BatchDeleteTeams(ctx, teamIDs)
	require.NoError(t, err)

	// Verify BatchDeleteTeamsResult format
	assert.Equal(t, 5, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)
	assert.Len(t, result.FailedIDs, 2)
	assert.Len(t, result.Errors, 2)

	// Verify error format
	for i, err := range result.Errors {
		assert.Contains(t, err, "团队")
		assert.Contains(t, err, ":")
		// Error format: "团队<ID>: <error message>"
		if i < len(result.FailedIDs) {
			assert.Contains(t, err, string(rune(result.FailedIDs[i])))
		}
	}
}

func TestTeamService_BatchUpdateStatusResult_Format(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create teams
	var teamIDs []uint64
	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, "batch_status_format")
		player := CreateTestPlayer(t, db, user)
		teamObj := createTestTeamForBatch(t, db, "Status Format Team", player.ID)
		teamIDs = append(teamIDs, teamObj.ID)
	}

	// Add non-existent IDs
	teamIDs = append(teamIDs, 99999)

	// Batch update status
	result, err := svc.BatchUpdateTeamStatus(ctx, teamIDs, model.TeamStatusInactive)
	require.NoError(t, err)

	// Verify BatchUpdateTeamStatusResult format
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.FailedIDs, 1)
	assert.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0], "团队")
}

func TestTeamService_BatchAddMembersResult_Format(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create team
	leaderUser := CreateUniqueTestUser(t, db, "batch_add_format")
	leader := CreateTestPlayer(t, db, leaderUser)
	teamObj := CreateTestTeam(t, db, "Format Test Team", leader.ID)

	// Create players and add non-existent IDs
	var playerIDs []uint64
	for i := 0; i < 2; i++ {
		user := CreateUniqueTestUser(t, db, "batch_add_format_member")
		player := CreateTestPlayer(t, db, user)
		playerIDs = append(playerIDs, player.ID)
	}
	playerIDs = append(playerIDs, 99999)

	// Batch add members
	result, err := svc.BatchAddMembers(ctx, teamObj.ID, playerIDs)
	require.NoError(t, err)

	// Verify BatchAddMembersResult format
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.FailedPlayerIDs, 1)
	assert.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0], "陪玩师")
}

// ============================================================================
// Database State Verification Tests
// ============================================================================

func TestTeamService_BatchDeleteTeams_DatabaseState(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create teams with members
	var teamIDs []uint64
	for i := 0; i < 2; i++ {
		user := CreateUniqueTestUser(t, db, "batch_delete_db_state")
		player := CreateTestPlayer(t, db, user)
		teamObj, _ := createTestTeamWithMembersForBatch(t, db, "DB State Team", player.ID, 2)
		teamIDs = append(teamIDs, teamObj.ID)
	}

	// Get initial member count
	var initialMemberCount int64
	db.Model(&model.TeamMember{}).Where("team_id IN ?", teamIDs).Count(&initialMemberCount)
	assert.Equal(t, int64(6), initialMemberCount) // 2 teams * (1 leader + 2 members)

	// Batch delete teams
	result, err := svc.BatchDeleteTeams(ctx, teamIDs)
	require.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount)

	// Verify teams are deleted (soft delete)
	var teamCount int64
	db.Model(&model.Team{}).Unscoped().Where("id IN ? AND deleted_at IS NOT NULL", teamIDs).Count(&teamCount)
	assert.Equal(t, int64(2), teamCount)

	// Verify team members are also soft deleted (GORM cascades)
	var memberCount int64
	db.Model(&model.TeamMember{}).Where("team_id IN ?", teamIDs).Count(&memberCount)
	assert.Equal(t, int64(0), memberCount) // All members soft deleted

	// Verify players still exist
	var playerCount int64
	db.Model(&model.Player{}).Where("id IN (SELECT player_id FROM team_members WHERE team_id IN ?)", teamIDs).Count(&playerCount)
	// Players should still exist, only team and member records are deleted
	assert.Greater(t, playerCount, int64(0))
}

func TestTeamService_BatchUpdateStatus_DatabaseState(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create teams with various initial statuses
	user1 := CreateUniqueTestUser(t, db, "batch_update_db_active")
	player1 := CreateTestPlayer(t, db, user1)
	team1 := CreateTestTeam(t, db, "Active Team", player1.ID)
	team1.Status = model.TeamStatusActive
	db.Save(team1)

	user2 := CreateUniqueTestUser(t, db, "batch_update_db_inactive")
	player2 := CreateTestPlayer(t, db, user2)
	team2 := CreateTestTeam(t, db, "Inactive Team", player2.ID)
	team2.Status = model.TeamStatusInactive
	db.Save(team2)

	teamIDs := []uint64{team1.ID, team2.ID}

	// Batch update to busy
	result, err := svc.BatchUpdateTeamStatus(ctx, teamIDs, model.TeamStatusBusy)
	require.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount)

	// Verify database state
	var teams []model.Team
	db.Where("id IN ?", teamIDs).Find(&teams)
	for _, team := range teams {
		assert.Equal(t, model.TeamStatusBusy, team.Status)
	}
}

func TestTeamService_BatchAddMembers_DatabaseState(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := teamrepo.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create team
	leaderUser := CreateUniqueTestUser(t, db, "batch_add_db_state")
	leader := CreateTestPlayer(t, db, leaderUser)
	teamObj := CreateTestTeam(t, db, "DB State Team", leader.ID)

	// Get initial state
	var initialMemberCount int64
	db.Model(&model.TeamMember{}).Where("team_id = ? AND status = ?", teamObj.ID, model.TeamMemberStatusActive).Count(&initialMemberCount)
	assert.Equal(t, int64(1), initialMemberCount) // Only leader

	// Create and add members
	var playerIDs []uint64
	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, "batch_add_db_member")
		player := CreateTestPlayer(t, db, user)
		playerIDs = append(playerIDs, player.ID)
	}

	// Batch add members
	result, err := svc.BatchAddMembers(ctx, teamObj.ID, playerIDs)
	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)

	// Verify team member count updated
	var finalMemberCount int64
	db.Model(&model.TeamMember{}).Where("team_id = ? AND status = ?", teamObj.ID, model.TeamMemberStatusActive).Count(&finalMemberCount)
	assert.Equal(t, int64(4), finalMemberCount) // 1 leader + 3 members

	// Verify team.MemberCount field updated
	team, err := repo.GetByID(ctx, teamObj.ID)
	require.NoError(t, err)
	assert.Equal(t, 4, team.MemberCount)

	// Verify all members have correct role and status
	var members []model.TeamMember
	db.Where("team_id = ? AND status = ?", teamObj.ID, model.TeamMemberStatusActive).Find(&members)
	assert.Len(t, members, 4)

	leaderCount := 0
	memberCount := 0
	for _, m := range members {
		if m.Role == model.TeamMemberRoleLeader {
			leaderCount++
		} else {
			memberCount++
		}
	}
	assert.Equal(t, 1, leaderCount)
	assert.Equal(t, 3, memberCount)
}
