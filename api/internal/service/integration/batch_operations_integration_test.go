package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository/activity"
	"gamelink/internal/repository/commission"
	"gamelink/internal/repository/team"
	activityservice "gamelink/internal/service/activity"
	commissionservice "gamelink/internal/service/commission"
	teamservice "gamelink/internal/service/team"
)

// ============================================================================
// Team Batch Operations Tests
// ============================================================================

func TestTeamService_BatchDeleteTeams(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := team.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create test teams
	var teamIDs []uint64
	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, "batch_delete_team")
		player := CreateTestPlayer(t, db, user)

		teamObj := &model.Team{
			Name:            "Batch Delete Team",
			LeaderID:        player.ID,
			Status:          model.TeamStatusActive,
			MaxMembers:      5,
			MemberCount:     1,
			IncomeShareType: "equal",
		}
		require.NoError(t, repo.Create(ctx, teamObj))
		teamIDs = append(teamIDs, teamObj.ID)
	}

	// Add a team with active order (should fail)
	user := CreateUniqueTestUser(t, db, "batch_delete_team_order")
	player := CreateTestPlayer(t, db, user)
	order := CreateTestOrder(t, db, user, player, model.OrderStatusInProgress)

	busyTeam := &model.Team{
		Name:            "Busy Team",
		LeaderID:        player.ID,
		Status:          model.TeamStatusActive,
		MaxMembers:      5,
		MemberCount:     1,
		IncomeShareType: "equal",
		CurrentOrderID:  &order.ID,
	}
	require.NoError(t, repo.Create(ctx, busyTeam))
	teamIDs = append(teamIDs, busyTeam.ID)

	// Batch delete
	result, err := svc.BatchDeleteTeams(ctx, teamIDs)
	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.FailedIDs, 1)
	assert.Len(t, result.Errors, 1)

	// Verify deleted teams
	for _, id := range teamIDs[:3] {
		_, err := repo.GetByID(ctx, id)
		assert.Error(t, err)
	}

	// Verify busy team still exists
	_, err = repo.GetByID(ctx, busyTeam.ID)
	assert.NoError(t, err)
}

func TestTeamService_BatchUpdateTeamStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := team.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create test teams
	var teamIDs []uint64
	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, "batch_status_team")
		player := CreateTestPlayer(t, db, user)

		teamObj := &model.Team{
			Name:            "Batch Status Team",
			LeaderID:        player.ID,
			Status:          model.TeamStatusActive,
			MaxMembers:      5,
			MemberCount:     1,
			IncomeShareType: "equal",
		}
		require.NoError(t, repo.Create(ctx, teamObj))
		teamIDs = append(teamIDs, teamObj.ID)
	}

	// Add a team with active order (should fail)
	user := CreateUniqueTestUser(t, db, "batch_status_team_order")
	player := CreateTestPlayer(t, db, user)
	order := CreateTestOrder(t, db, user, player, model.OrderStatusInProgress)

	busyTeam := &model.Team{
		Name:            "Busy Team Status",
		LeaderID:        player.ID,
		Status:          model.TeamStatusActive,
		MaxMembers:      5,
		MemberCount:     1,
		IncomeShareType: "equal",
		CurrentOrderID:  &order.ID,
	}
	require.NoError(t, repo.Create(ctx, busyTeam))
	teamIDs = append(teamIDs, busyTeam.ID)

	// Batch update status to inactive
	result, err := svc.BatchUpdateTeamStatus(ctx, teamIDs, model.TeamStatusInactive)
	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)

	// Verify status updates
	for _, id := range teamIDs[:3] {
		team, err := repo.GetByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.TeamStatusInactive, team.Status)
	}

	// Verify busy team still active
	busyTeamResult, err := repo.GetByID(ctx, busyTeam.ID)
	require.NoError(t, err)
	assert.Equal(t, model.TeamStatusActive, busyTeamResult.Status)
}

func TestTeamService_BatchAddMembers(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := team.NewTeamRepository(db)
	svc := teamservice.NewTeamService(repo)
	ctx := context.Background()

	// Create team
	leaderUser := CreateUniqueTestUser(t, db, "batch_add_leader")
	leader := CreateTestPlayer(t, db, leaderUser)

	teamObj := &model.Team{
		Name:            "Batch Add Members Team",
		LeaderID:        leader.ID,
		Status:          model.TeamStatusActive,
		MaxMembers:      10,
		MemberCount:     1,
		IncomeShareType: "equal",
	}
	require.NoError(t, repo.Create(ctx, teamObj))

	// Create players to add
	var playerIDs []uint64
	for i := 0; i < 5; i++ {
		user := CreateUniqueTestUser(t, db, "batch_add_member")
		player := CreateTestPlayer(t, db, user)
		playerIDs = append(playerIDs, player.ID)
	}

	// Create a player already in another team (should fail)
	anotherUser := CreateUniqueTestUser(t, db, "batch_add_another")
	anotherPlayer := CreateTestPlayer(t, db, anotherUser)

	anotherTeam := &model.Team{
		Name:            "Another Team",
		LeaderID:        anotherPlayer.ID,
		Status:          model.TeamStatusActive,
		MaxMembers:      5,
		MemberCount:     1,
		IncomeShareType: "equal",
	}
	require.NoError(t, repo.Create(ctx, anotherTeam))
	playerIDs = append(playerIDs, anotherPlayer.ID)

	// Batch add members
	result, err := svc.BatchAddMembers(ctx, teamObj.ID, playerIDs)
	require.NoError(t, err)
	assert.Equal(t, 5, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.FailedPlayerIDs, 1)

	// Verify team member count
	team, err := repo.GetByID(ctx, teamObj.ID)
	require.NoError(t, err)
	assert.Equal(t, 6, team.MemberCount) // 1 leader + 5 new members
}

// ============================================================================
// Activity Batch Operations Tests
// ============================================================================

func TestActivityService_BatchDeleteActivities(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := activity.NewActivityRepository(db)
	svc := activityservice.NewActivityService(repo, nil)
	ctx := context.Background()

	// Create test activities
	var activityIDs []uint64
	for i := 0; i < 3; i++ {
		activity := &model.Activity{
			Name:         "Batch Delete Activity",
			Type:         model.ActivityTypeCoupon,
			Status:       model.ActivityStatusDraft,
			StartAt:      time.Now().Add(24 * time.Hour),
			EndAt:        time.Now().Add(48 * time.Hour),
			TotalLimit:   1000,
			DailyLimit:   100,
			PerUserLimit: 1,
		}
		require.NoError(t, repo.CreateActivity(ctx, activity))
		activityIDs = append(activityIDs, activity.ID)
	}

	// Add an active activity (should fail)
	activeActivity := &model.Activity{
		Name:         "Active Activity",
		Type:         model.ActivityTypeCoupon,
		Status:       model.ActivityStatusActive,
		StartAt:      time.Now().Add(-24 * time.Hour),
		EndAt:        time.Now().Add(24 * time.Hour),
		TotalLimit:   1000,
		DailyLimit:   100,
		PerUserLimit: 1,
	}
	require.NoError(t, repo.CreateActivity(ctx, activeActivity))
	activityIDs = append(activityIDs, activeActivity.ID)

	// Batch delete
	result, err := svc.BatchDeleteActivities(ctx, activityIDs)
	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)

	// Verify deleted activities
	for _, id := range activityIDs[:3] {
		_, err := repo.GetActivityByID(ctx, id)
		assert.Error(t, err)
	}

	// Verify active activity still exists
	_, err = repo.GetActivityByID(ctx, activeActivity.ID)
	assert.NoError(t, err)
}

func TestActivityService_BatchUpdateActivityStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := activity.NewActivityRepository(db)
	svc := activityservice.NewActivityService(repo, nil)
	ctx := context.Background()

	// Create test activities
	var activityIDs []uint64
	for i := 0; i < 3; i++ {
		activity := &model.Activity{
			Name:         "Batch Status Activity",
			Type:         model.ActivityTypeCoupon,
			Status:       model.ActivityStatusDraft,
			StartAt:      time.Now().Add(24 * time.Hour),
			EndAt:        time.Now().Add(48 * time.Hour),
			TotalLimit:   1000,
			DailyLimit:   100,
			PerUserLimit: 1,
		}
		require.NoError(t, repo.CreateActivity(ctx, activity))
		activityIDs = append(activityIDs, activity.ID)
	}

	// Add an ended activity (invalid status transition)
	endedActivity := &model.Activity{
		Name:         "Ended Activity",
		Type:         model.ActivityTypeCoupon,
		Status:       model.ActivityStatusEnded,
		StartAt:      time.Now().Add(-48 * time.Hour),
		EndAt:        time.Now().Add(-24 * time.Hour),
		TotalLimit:   1000,
		DailyLimit:   100,
		PerUserLimit: 1,
	}
	require.NoError(t, repo.CreateActivity(ctx, endedActivity))
	activityIDs = append(activityIDs, endedActivity.ID)

	// Batch update status
	result, err := svc.BatchUpdateActivityStatus(ctx, activityIDs, model.ActivityStatusActive)
	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)

	// Verify status updates
	for _, id := range activityIDs[:3] {
		activity, err := repo.GetActivityByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.ActivityStatusActive, activity.Status)
	}

	// Verify ended activity status unchanged
	ended, err := repo.GetActivityByID(ctx, endedActivity.ID)
	require.NoError(t, err)
	assert.Equal(t, model.ActivityStatusEnded, ended.Status)
}

func TestActivityService_BatchPublishActivities(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := activity.NewActivityRepository(db)
	svc := activityservice.NewActivityService(repo, nil)
	ctx := context.Background()

	// Create test activities
	var activityIDs []uint64
	for i := 0; i < 3; i++ {
		activity := &model.Activity{
			Name:         "Batch Publish Activity",
			Type:         model.ActivityTypeCoupon,
			Status:       model.ActivityStatusDraft,
			IsVisible:    false,
			StartAt:      time.Now().Add(24 * time.Hour),
			EndAt:        time.Now().Add(48 * time.Hour),
			TotalLimit:   1000,
			DailyLimit:   100,
			PerUserLimit: 1,
		}
		require.NoError(t, repo.CreateActivity(ctx, activity))
		activityIDs = append(activityIDs, activity.ID)
	}

	// Batch publish (set isVisible = true)
	result, err := svc.BatchPublishActivities(ctx, activityIDs, true)
	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)

	// Verify visibility
	for _, id := range activityIDs {
		activity, err := repo.GetActivityByID(ctx, id)
		require.NoError(t, err)
		assert.True(t, activity.IsVisible)
	}

	// Batch unpublish (set isVisible = false)
	result, err = svc.BatchPublishActivities(ctx, activityIDs, false)
	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)

	// Verify visibility
	for _, id := range activityIDs {
		activity, err := repo.GetActivityByID(ctx, id)
		require.NoError(t, err)
		assert.False(t, activity.IsVisible)
	}
}

// ============================================================================
// Commission Batch Operations Tests
// ============================================================================

func TestCommissionService_BatchDeleteCommissionRules(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := commission.NewCommissionRepository(db)
	svc := commissionservice.NewCommissionService(repo, nil, nil)
	ctx := context.Background()

	// Create test rules
	var ruleIDs []uint64
	for i := 0; i < 3; i++ {
		rule := &model.CommissionRule{
			Name:     "Batch Delete Rule",
			Type:     model.CommissionRuleTypeDefault,
			Rate:     20,
			IsActive: true,
		}
		require.NoError(t, repo.CreateRule(ctx, rule))
		ruleIDs = append(ruleIDs, rule.ID)
	}

	// Add non-existent rule (should fail)
	ruleIDs = append(ruleIDs, 99999)

	// Batch delete
	result, err := svc.BatchDeleteCommissionRules(ctx, ruleIDs)
	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)

	// Verify deleted rules
	for _, id := range ruleIDs[:3] {
		_, err := repo.GetRule(ctx, id)
		assert.Error(t, err)
	}
}

func TestCommissionService_BatchUpdateCommissionRuleStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := commission.NewCommissionRepository(db)
	svc := commissionservice.NewCommissionService(repo, nil, nil)
	ctx := context.Background()

	// Create test rules
	var ruleIDs []uint64
	for i := 0; i < 3; i++ {
		rule := &model.CommissionRule{
			Name:     "Batch Status Rule",
			Type:     model.CommissionRuleTypeDefault,
			Rate:     20,
			IsActive: true,
		}
		require.NoError(t, repo.CreateRule(ctx, rule))
		ruleIDs = append(ruleIDs, rule.ID)
	}

	// Add non-existent rule (should fail)
	ruleIDs = append(ruleIDs, 99999)

	// Batch disable (set isActive = false)
	result, err := svc.BatchUpdateCommissionRuleStatus(ctx, ruleIDs, false)
	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)

	// Verify status updates
	for _, id := range ruleIDs[:3] {
		rule, err := repo.GetRule(ctx, id)
		require.NoError(t, err)
		assert.False(t, rule.IsActive)
	}

	// Batch enable (set isActive = true)
	result, err = svc.BatchUpdateCommissionRuleStatus(ctx, ruleIDs[:3], true)
	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)

	// Verify status updates
	for _, id := range ruleIDs[:3] {
		rule, err := repo.GetRule(ctx, id)
		require.NoError(t, err)
		assert.True(t, rule.IsActive)
	}
}

// ============================================================================
// Edge Cases Tests
// ============================================================================

func TestBatchOperations_EmptyIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Team service - empty IDs
	teamRepo := team.NewTeamRepository(db)
	teamSvc := teamservice.NewTeamService(teamRepo)
	_, err := teamSvc.BatchDeleteTeams(ctx, []uint64{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "团队ID列表不能为空")

	_, err = teamSvc.BatchUpdateTeamStatus(ctx, []uint64{}, model.TeamStatusActive)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "团队ID列表不能为空")

	_, err = teamSvc.BatchAddMembers(ctx, 1, []uint64{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "陪玩师ID列表不能为空")

	// Activity service - empty IDs
	activityRepo := activity.NewActivityRepository(db)
	activitySvc := activityservice.NewActivityService(activityRepo, nil)
	_, err = activitySvc.BatchDeleteActivities(ctx, []uint64{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "活动ID列表不能为空")

	_, err = activitySvc.BatchUpdateActivityStatus(ctx, []uint64{}, model.ActivityStatusActive)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "活动ID列表不能为空")

	_, err = activitySvc.BatchPublishActivities(ctx, []uint64{}, true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "活动ID列表不能为空")

	// Commission service - empty IDs
	commissionRepo := commission.NewCommissionRepository(db)
	commissionSvc := commissionservice.NewCommissionService(commissionRepo, nil, nil)
	_, err = commissionSvc.BatchDeleteCommissionRules(ctx, []uint64{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "规则ID列表不能为空")

	_, err = commissionSvc.BatchUpdateCommissionRuleStatus(ctx, []uint64{}, true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "规则ID列表不能为空")
}

func TestBatchOperations_InvalidStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Team service - invalid status
	teamRepo := team.NewTeamRepository(db)
	teamSvc := teamservice.NewTeamService(teamRepo)
	_, err := teamSvc.BatchUpdateTeamStatus(ctx, []uint64{1, 2}, "invalid_status")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无效的团队状态")

	// Activity service - invalid status
	activityRepo := activity.NewActivityRepository(db)
	activitySvc := activityservice.NewActivityService(activityRepo, nil)
	_, err = activitySvc.BatchUpdateActivityStatus(ctx, []uint64{1, 2}, "invalid_status")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无效的活动状态")
}
