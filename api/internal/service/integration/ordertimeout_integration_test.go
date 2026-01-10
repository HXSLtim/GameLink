// Package integration provides integration tests for services.
package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/ordertimeout"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// OrderTimeoutConfig Tests
// ============================================================================

func TestOrderTimeoutRepository_GetConfig(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)
	config := CreateTestOrderTimeoutConfig(t, db, "test_config_get", "30", "Test config")

	result, err := repo.GetConfig(ctx, config.ConfigKey)
	require.NoError(t, err)
	assert.Equal(t, "30", result.ConfigValue)
}

func TestOrderTimeoutRepository_GetConfig_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	_, err := repo.GetConfig(ctx, "nonexistent_key")
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestOrderTimeoutRepository_ListConfigs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	CreateTestOrderTimeoutConfig(t, db, "config_list_1", "10", "Config 1")
	CreateTestOrderTimeoutConfig(t, db, "config_list_2", "20", "Config 2")

	configs, err := repo.ListConfigs(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(configs), 2)
}

func TestOrderTimeoutRepository_SaveConfig_Create(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	config := &model.OrderTimeoutConfig{
		ConfigKey:   "new_config_key",
		ConfigValue: "60",
		Description: "New config",
	}

	err := repo.SaveConfig(ctx, config)
	require.NoError(t, err)

	// Verify
	result, err := repo.GetConfig(ctx, "new_config_key")
	require.NoError(t, err)
	assert.Equal(t, "60", result.ConfigValue)
}

func TestOrderTimeoutRepository_SaveConfig_Update(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)
	CreateTestOrderTimeoutConfig(t, db, "update_config", "30", "Original")

	// Update
	config := &model.OrderTimeoutConfig{
		ConfigKey:   "update_config",
		ConfigValue: "45",
		Description: "Updated",
	}
	err := repo.SaveConfig(ctx, config)
	require.NoError(t, err)

	// Verify
	result, err := repo.GetConfig(ctx, "update_config")
	require.NoError(t, err)
	assert.Equal(t, "45", result.ConfigValue)
	assert.Equal(t, "Updated", result.Description)
}

func TestOrderTimeoutRepository_DeleteConfig(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)
	CreateTestOrderTimeoutConfig(t, db, "delete_config", "30", "To delete")

	err := repo.DeleteConfig(ctx, "delete_config")
	require.NoError(t, err)

	_, err = repo.GetConfig(ctx, "delete_config")
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

// ============================================================================
// OrderTimeoutLog Tests
// ============================================================================

func TestOrderTimeoutRepository_CreateLog(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	user := CreateUniqueTestUser(t, db, "ot_log_create")
	playerUser := CreateUniqueTestUser(t, db, "ot_log_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "ot_log_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusPending, 10000)

	log := &model.OrderTimeoutLog{
		OrderID:     order.ID,
		TimeoutType: model.OrderTimeoutTypePayment,
		TimeoutAt:   time.Now(),
		Action:      model.OrderTimeoutActionCanceled,
		Remark:      "Payment timeout",
	}

	err := repo.CreateLog(ctx, log)
	require.NoError(t, err)
	assert.NotZero(t, log.ID)
}

func TestOrderTimeoutRepository_GetLog(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	user := CreateUniqueTestUser(t, db, "ot_log_get")
	playerUser := CreateUniqueTestUser(t, db, "ot_log_get_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "ot_log_get_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusPending, 10000)
	log := CreateTestOrderTimeoutLog(t, db, order, model.OrderTimeoutTypePayment, model.OrderTimeoutActionCanceled)

	result, err := repo.GetLog(ctx, log.ID)
	require.NoError(t, err)
	assert.Equal(t, log.ID, result.ID)
	assert.Equal(t, order.ID, result.OrderID)
}

func TestOrderTimeoutRepository_GetLog_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	_, err := repo.GetLog(ctx, 99999)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestOrderTimeoutRepository_GetLogWithOrder(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	user := CreateUniqueTestUser(t, db, "ot_log_with_order")
	playerUser := CreateUniqueTestUser(t, db, "ot_log_with_order_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "ot_log_with_order_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusPending, 10000)
	log := CreateTestOrderTimeoutLog(t, db, order, model.OrderTimeoutTypeAccept, model.OrderTimeoutActionNotified)

	result, err := repo.GetLogWithOrder(ctx, log.ID)
	require.NoError(t, err)
	assert.NotNil(t, result.Order)
	assert.Equal(t, order.ID, result.Order.ID)
}

func TestOrderTimeoutRepository_ListLogsByOrderID(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	user := CreateUniqueTestUser(t, db, "ot_logs_by_order")
	playerUser := CreateUniqueTestUser(t, db, "ot_logs_by_order_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "ot_logs_by_order_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusPending, 10000)

	CreateTestOrderTimeoutLog(t, db, order, model.OrderTimeoutTypePayment, model.OrderTimeoutActionNotified)
	CreateTestOrderTimeoutLog(t, db, order, model.OrderTimeoutTypePayment, model.OrderTimeoutActionCanceled)

	logs, err := repo.ListLogsByOrderID(ctx, order.ID)
	require.NoError(t, err)
	assert.Len(t, logs, 2)
}

func TestOrderTimeoutRepository_ListLogsPaged(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	// Create multiple logs
	for i := 0; i < 5; i++ {
		user := CreateUniqueTestUser(t, db, fmt.Sprintf("ot_logs_paged_%d", i))
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("ot_logs_paged_player_%d", i))
		player := CreateTestPlayer(t, db, playerUser)
		game := CreateTestGame(t, db, fmt.Sprintf("ot_logs_paged_game_%d_%d", i, time.Now().UnixNano()))
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusPending, 10000)
		CreateTestOrderTimeoutLog(t, db, order, model.OrderTimeoutTypePayment, model.OrderTimeoutActionCanceled)
	}

	logs, total, err := repo.ListLogsPaged(ctx, repository.OrderTimeoutLogListOptions{
		Page:     1,
		PageSize: 3,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, logs, 3)
}

func TestOrderTimeoutRepository_ListLogsPaged_FilterByType(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	user := CreateUniqueTestUser(t, db, "ot_logs_type")
	playerUser := CreateUniqueTestUser(t, db, "ot_logs_type_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "ot_logs_type_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusPending, 10000)

	CreateTestOrderTimeoutLog(t, db, order, model.OrderTimeoutTypePayment, model.OrderTimeoutActionCanceled)
	CreateTestOrderTimeoutLog(t, db, order, model.OrderTimeoutTypeAccept, model.OrderTimeoutActionNotified)

	timeoutType := model.OrderTimeoutTypePayment
	logs, total, err := repo.ListLogsPaged(ctx, repository.OrderTimeoutLogListOptions{
		TimeoutType: &timeoutType,
		Page:        1,
		PageSize:    10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	for _, l := range logs {
		assert.Equal(t, model.OrderTimeoutTypePayment, l.TimeoutType)
	}
}

func TestOrderTimeoutRepository_GetLogStats(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	// Create logs with different types
	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, fmt.Sprintf("ot_stats_payment_%d", i))
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("ot_stats_payment_player_%d", i))
		player := CreateTestPlayer(t, db, playerUser)
		game := CreateTestGame(t, db, fmt.Sprintf("ot_stats_payment_game_%d_%d", i, time.Now().UnixNano()))
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusPending, 10000)
		CreateTestOrderTimeoutLog(t, db, order, model.OrderTimeoutTypePayment, model.OrderTimeoutActionCanceled)
	}

	for i := 0; i < 2; i++ {
		user := CreateUniqueTestUser(t, db, fmt.Sprintf("ot_stats_accept_%d", i))
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("ot_stats_accept_player_%d", i))
		player := CreateTestPlayer(t, db, playerUser)
		game := CreateTestGame(t, db, fmt.Sprintf("ot_stats_accept_game_%d_%d", i, time.Now().UnixNano()))
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusConfirmed, 10000)
		CreateTestOrderTimeoutLog(t, db, order, model.OrderTimeoutTypeAccept, model.OrderTimeoutActionNotified)
	}

	stats, err := repo.GetLogStats(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, stats[model.OrderTimeoutTypePayment], int64(3))
	assert.GreaterOrEqual(t, stats[model.OrderTimeoutTypeAccept], int64(2))
}

// ============================================================================
// OrderServiceAssignment Tests
// ============================================================================

func TestOrderTimeoutRepository_CreateAssignment(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	user := CreateUniqueTestUser(t, db, "ot_assign_create")
	playerUser := CreateUniqueTestUser(t, db, "ot_assign_create_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "ot_assign_create_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusInProgress, 10000)
	serviceUser := CreateUniqueTestUser(t, db, "ot_assign_service")

	assignment := &model.OrderServiceAssignment{
		OrderID:       order.ID,
		ServiceUserID: serviceUser.ID,
		Status:        model.ServiceAssignmentStatusAssigned,
		AssignedAt:    time.Now(),
		AssignType:    "auto",
	}

	err := repo.CreateAssignment(ctx, assignment)
	require.NoError(t, err)
	assert.NotZero(t, assignment.ID)
}

func TestOrderTimeoutRepository_GetAssignment(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	user := CreateUniqueTestUser(t, db, "ot_assign_get")
	playerUser := CreateUniqueTestUser(t, db, "ot_assign_get_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "ot_assign_get_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusInProgress, 10000)
	serviceUser := CreateUniqueTestUser(t, db, "ot_assign_get_service")
	assignment := CreateTestOrderServiceAssignment(t, db, order, serviceUser, model.ServiceAssignmentStatusAssigned)

	result, err := repo.GetAssignment(ctx, assignment.ID)
	require.NoError(t, err)
	assert.Equal(t, assignment.ID, result.ID)
	assert.Equal(t, order.ID, result.OrderID)
}

func TestOrderTimeoutRepository_GetAssignment_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	_, err := repo.GetAssignment(ctx, 99999)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestOrderTimeoutRepository_GetAssignmentByOrderID(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	user := CreateUniqueTestUser(t, db, "ot_assign_by_order")
	playerUser := CreateUniqueTestUser(t, db, "ot_assign_by_order_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "ot_assign_by_order_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusInProgress, 10000)
	serviceUser := CreateUniqueTestUser(t, db, "ot_assign_by_order_service")
	CreateTestOrderServiceAssignment(t, db, order, serviceUser, model.ServiceAssignmentStatusAssigned)

	result, err := repo.GetAssignmentByOrderID(ctx, order.ID)
	require.NoError(t, err)
	assert.Equal(t, order.ID, result.OrderID)
}

func TestOrderTimeoutRepository_ListAssignmentsByServiceUser(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	serviceUser := CreateUniqueTestUser(t, db, "ot_assign_list_service")

	// Create multiple assignments for the same service user
	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, fmt.Sprintf("ot_assign_list_%d", i))
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("ot_assign_list_player_%d", i))
		player := CreateTestPlayer(t, db, playerUser)
		game := CreateTestGame(t, db, fmt.Sprintf("ot_assign_list_game_%d_%d", i, time.Now().UnixNano()))
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusInProgress, 10000)
		CreateTestOrderServiceAssignment(t, db, order, serviceUser, model.ServiceAssignmentStatusAssigned)
	}

	assignments, err := repo.ListAssignmentsByServiceUser(ctx, serviceUser.ID, nil)
	require.NoError(t, err)
	assert.Len(t, assignments, 3)
}

func TestOrderTimeoutRepository_ListAssignmentsPaged(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	// Create multiple assignments
	for i := 0; i < 5; i++ {
		user := CreateUniqueTestUser(t, db, fmt.Sprintf("ot_assign_paged_%d", i))
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("ot_assign_paged_player_%d", i))
		player := CreateTestPlayer(t, db, playerUser)
		game := CreateTestGame(t, db, fmt.Sprintf("ot_assign_paged_game_%d_%d", i, time.Now().UnixNano()))
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusInProgress, 10000)
		serviceUser := CreateUniqueTestUser(t, db, fmt.Sprintf("ot_assign_paged_service_%d", i))
		CreateTestOrderServiceAssignment(t, db, order, serviceUser, model.ServiceAssignmentStatusAssigned)
	}

	assignments, total, err := repo.ListAssignmentsPaged(ctx, repository.ServiceAssignmentListOptions{
		Page:     1,
		PageSize: 3,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, assignments, 3)
}

func TestOrderTimeoutRepository_UpdateAssignmentStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	user := CreateUniqueTestUser(t, db, "ot_assign_status")
	playerUser := CreateUniqueTestUser(t, db, "ot_assign_status_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "ot_assign_status_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusInProgress, 10000)
	serviceUser := CreateUniqueTestUser(t, db, "ot_assign_status_service")
	assignment := CreateTestOrderServiceAssignment(t, db, order, serviceUser, model.ServiceAssignmentStatusAssigned)

	// Update to joined
	err := repo.UpdateAssignmentStatus(ctx, assignment.ID, model.ServiceAssignmentStatusJoined)
	require.NoError(t, err)

	// Verify
	updated, err := repo.GetAssignment(ctx, assignment.ID)
	require.NoError(t, err)
	assert.Equal(t, model.ServiceAssignmentStatusJoined, updated.Status)
	assert.NotNil(t, updated.JoinedAt)
}

func TestOrderTimeoutRepository_UpdateAssignmentStatus_Completed(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	user := CreateUniqueTestUser(t, db, "ot_assign_completed")
	playerUser := CreateUniqueTestUser(t, db, "ot_assign_completed_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "ot_assign_completed_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusInProgress, 10000)
	serviceUser := CreateUniqueTestUser(t, db, "ot_assign_completed_service")
	assignment := CreateTestOrderServiceAssignment(t, db, order, serviceUser, model.ServiceAssignmentStatusJoined)

	// Update to completed
	err := repo.UpdateAssignmentStatus(ctx, assignment.ID, model.ServiceAssignmentStatusCompleted)
	require.NoError(t, err)

	// Verify
	updated, err := repo.GetAssignment(ctx, assignment.ID)
	require.NoError(t, err)
	assert.Equal(t, model.ServiceAssignmentStatusCompleted, updated.Status)
	assert.NotNil(t, updated.LeftAt)
}

func TestOrderTimeoutRepository_DeleteAssignment(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	user := CreateUniqueTestUser(t, db, "ot_assign_delete")
	playerUser := CreateUniqueTestUser(t, db, "ot_assign_delete_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "ot_assign_delete_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusInProgress, 10000)
	serviceUser := CreateUniqueTestUser(t, db, "ot_assign_delete_service")
	assignment := CreateTestOrderServiceAssignment(t, db, order, serviceUser, model.ServiceAssignmentStatusAssigned)

	err := repo.DeleteAssignment(ctx, assignment.ID)
	require.NoError(t, err)

	_, err = repo.GetAssignment(ctx, assignment.ID)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestOrderTimeoutRepository_GetAssignmentStats(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	// Create assignments with different statuses
	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, fmt.Sprintf("ot_stats_assigned_%d", i))
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("ot_stats_assigned_player_%d", i))
		player := CreateTestPlayer(t, db, playerUser)
		game := CreateTestGame(t, db, fmt.Sprintf("ot_stats_assigned_game_%d_%d", i, time.Now().UnixNano()))
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusInProgress, 10000)
		serviceUser := CreateUniqueTestUser(t, db, fmt.Sprintf("ot_stats_assigned_service_%d", i))
		CreateTestOrderServiceAssignment(t, db, order, serviceUser, model.ServiceAssignmentStatusAssigned)
	}

	for i := 0; i < 2; i++ {
		user := CreateUniqueTestUser(t, db, fmt.Sprintf("ot_stats_joined_%d", i))
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("ot_stats_joined_player_%d", i))
		player := CreateTestPlayer(t, db, playerUser)
		game := CreateTestGame(t, db, fmt.Sprintf("ot_stats_joined_game_%d_%d", i, time.Now().UnixNano()))
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusInProgress, 10000)
		serviceUser := CreateUniqueTestUser(t, db, fmt.Sprintf("ot_stats_joined_service_%d", i))
		CreateTestOrderServiceAssignment(t, db, order, serviceUser, model.ServiceAssignmentStatusJoined)
	}

	stats, err := repo.GetAssignmentStats(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, stats[model.ServiceAssignmentStatusAssigned], int64(3))
	assert.GreaterOrEqual(t, stats[model.ServiceAssignmentStatusJoined], int64(2))
}

func TestOrderTimeoutRepository_GetActiveAssignmentCount(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := ordertimeout.NewOrderTimeoutRepository(db)

	serviceUser := CreateUniqueTestUser(t, db, "ot_active_count_service")

	// Create active assignments (assigned and joined)
	for i := 0; i < 2; i++ {
		user := CreateUniqueTestUser(t, db, fmt.Sprintf("ot_active_assigned_%d", i))
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("ot_active_assigned_player_%d", i))
		player := CreateTestPlayer(t, db, playerUser)
		game := CreateTestGame(t, db, fmt.Sprintf("ot_active_assigned_game_%d_%d", i, time.Now().UnixNano()))
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusInProgress, 10000)
		CreateTestOrderServiceAssignment(t, db, order, serviceUser, model.ServiceAssignmentStatusAssigned)
	}

	// Create completed assignment (should not count)
	user := CreateUniqueTestUser(t, db, "ot_active_completed")
	playerUser := CreateUniqueTestUser(t, db, "ot_active_completed_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "ot_active_completed_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 10000)
	CreateTestOrderServiceAssignment(t, db, order, serviceUser, model.ServiceAssignmentStatusCompleted)

	count, err := repo.GetActiveAssignmentCount(ctx, serviceUser.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}
