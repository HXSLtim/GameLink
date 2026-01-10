// Package integration provides CRUD integration tests for SettlementCompany module.
package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository/settlementcompany"
)

// ============================================================================
// SettlementCompany CRUD Integration Tests
// ============================================================================

// TestSettlementCompanyRepository_Create tests creating a new settlement company
func TestSettlementCompanyRepository_Create(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	admin := CreateUniqueTestUser(t, db, "admin_create")

	company := &model.SettlementCompany{
		Name:              "Test Settlement Company",
		CreditCode:        fmt.Sprintf("91110000%010d", time.Now().UnixNano()%10000000000),
		TaxRegistrationNo: "123456789012345678",
		BankName:          "Test Bank",
		BankAccount:       "1234567890",
		BankBranch:        "Test Branch",
		ContactName:       "Test Contact",
		ContactPhone:      "13800138000",
		Address:           "Test Address",
		Status:            model.CompanyStatusActive,
		CreatedBy:         admin.ID,
	}

	err := repo.Create(ctx, company)
	require.NoError(t, err)
	assert.NotZero(t, company.ID)
	assert.NotZero(t, company.CreatedAt)
}

// TestSettlementCompanyRepository_Get tests retrieving a company by ID
func TestSettlementCompanyRepository_Get(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	company := CreateTestSettlementCompany(t, db, "Get Test Company")

	got, err := repo.Get(ctx, company.ID)
	require.NoError(t, err)
	assert.Equal(t, company.Name, got.Name)
	assert.Equal(t, company.CreditCode, got.CreditCode)
	assert.Equal(t, model.CompanyStatusActive, got.Status)
}

// TestSettlementCompanyRepository_Get_NonExistent tests getting non-existent company
func TestSettlementCompanyRepository_Get_NonExistent(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	_, err := repo.Get(ctx, 99999)
	assert.Error(t, err)
}

// TestSettlementCompanyRepository_GetByCreditCode tests retrieving by credit code
func TestSettlementCompanyRepository_GetByCreditCode(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	company := CreateTestSettlementCompany(t, db, "CreditCode Test")

	got, err := repo.GetByCreditCode(ctx, company.CreditCode)
	require.NoError(t, err)
	assert.Equal(t, company.ID, got.ID)
	assert.Equal(t, company.Name, got.Name)
}

// TestSettlementCompanyRepository_GetByCreditCode_NotFound tests with non-existent credit code
func TestSettlementCompanyRepository_GetByCreditCode_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	_, err := repo.GetByCreditCode(ctx, "NONEXISTENT123456")
	assert.Error(t, err)
}

// TestSettlementCompanyRepository_Update tests updating a company
func TestSettlementCompanyRepository_Update(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	company := CreateTestSettlementCompany(t, db, "Update Test Company")

	// Update company details
	company.Name = "Updated Company Name"
	company.ContactPhone = "13900139000"
	updater := uint64(999)
	company.UpdatedBy = &updater

	err := repo.Update(ctx, company)
	require.NoError(t, err)

	got, err := repo.Get(ctx, company.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Company Name", got.Name)
	assert.Equal(t, "13900139000", got.ContactPhone)
}

// TestSettlementCompanyRepository_Delete tests deleting a company
func TestSettlementCompanyRepository_Delete(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	company := CreateTestSettlementCompany(t, db, "Delete Test Company")

	err := repo.Delete(ctx, company.ID)
	require.NoError(t, err)

	// Verify company is deleted (soft delete)
	_, err = repo.Get(ctx, company.ID)
	assert.Error(t, err)
}

// TestSettlementCompanyRepository_List tests listing companies
func TestSettlementCompanyRepository_List(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	// Create multiple companies
	for i := 0; i < 3; i++ {
		CreateTestSettlementCompany(t, db, fmt.Sprintf("List Company %d", i))
	}

	companies, total, err := repo.List(ctx, settlementcompany.ListOptions{
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(3))
	assert.GreaterOrEqual(t, len(companies), 3)
}

// TestSettlementCompanyRepository_List_ByStatus tests filtering by status
func TestSettlementCompanyRepository_List_ByStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	// Create active companies
	for i := 0; i < 2; i++ {
		company := CreateTestSettlementCompany(t, db, fmt.Sprintf("Active Company %d", i))
		company.Status = model.CompanyStatusActive
		db.Save(company)
	}

	// Create inactive company
	inactiveCompany := CreateTestSettlementCompany(t, db, "Inactive Company")
	inactiveCompany.Status = model.CompanyStatusInactive
	db.Save(inactiveCompany)

	// List only active companies
	activeStatus := model.CompanyStatusActive
	companies, total, err := repo.List(ctx, settlementcompany.ListOptions{
		Status:   &activeStatus,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(2))
	for _, c := range companies {
		assert.Equal(t, model.CompanyStatusActive, c.Status)
	}
}

// TestSettlementCompanyRepository_List_ByKeyword tests searching by keyword
func TestSettlementCompanyRepository_List_ByKeyword(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	uniqueKeyword := fmt.Sprintf("UniqueKeyword_%d", time.Now().UnixNano())
	company := CreateTestSettlementCompany(t, db, uniqueKeyword+" Company")

	companies, total, err := repo.List(ctx, settlementcompany.ListOptions{
		Keyword:  uniqueKeyword,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	found := false
	for _, c := range companies {
		if c.ID == company.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "Company should be found in search results")
}

// TestSettlementCompanyRepository_List_Sorting tests sorting options
func TestSettlementCompanyRepository_List_Sorting(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	// Create companies with different player counts
	for i := 0; i < 3; i++ {
		company := CreateTestSettlementCompany(t, db, fmt.Sprintf("Sort Company %d", i))
		company.PlayerCount = i
		db.Save(company)
	}

	// Sort by player count descending
	companies, _, err := repo.List(ctx, settlementcompany.ListOptions{
		SortBy:    "player_count",
		SortOrder: "desc",
		Page:      1,
		PageSize:  10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(companies), 3)
	// First company should have highest player count
	if len(companies) >= 3 {
		assert.GreaterOrEqual(t, companies[0].PlayerCount, companies[1].PlayerCount)
		assert.GreaterOrEqual(t, companies[1].PlayerCount, companies[2].PlayerCount)
	}
}

// TestSettlementCompanyRepository_ToggleStatus tests toggling company status
func TestSettlementCompanyRepository_ToggleStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	company := CreateTestSettlementCompany(t, db, "Toggle Status Company")
	assert.Equal(t, model.CompanyStatusActive, company.Status)

	// Deactivate
	err := repo.ToggleStatus(ctx, company.ID, model.CompanyStatusInactive)
	require.NoError(t, err)

	got, err := repo.Get(ctx, company.ID)
	require.NoError(t, err)
	assert.Equal(t, model.CompanyStatusInactive, got.Status)

	// Reactivate
	err = repo.ToggleStatus(ctx, company.ID, model.CompanyStatusActive)
	require.NoError(t, err)

	got, err = repo.Get(ctx, company.ID)
	require.NoError(t, err)
	assert.Equal(t, model.CompanyStatusActive, got.Status)
}

// TestSettlementCompanyRepository_AssignPlayer tests assigning a player to a company
func TestSettlementCompanyRepository_AssignPlayer(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	company := CreateTestSettlementCompany(t, db, "Assign Player Company")
	user := CreateUniqueTestUser(t, db, "player_to_assign")
	player := CreateTestPlayer(t, db, user)

	now := time.Now()
	assignment := &model.PlayerCompanyAssignment{
		PlayerID:            player.ID,
		SettlementCompanyID: company.ID,
		EffectiveDate:       now,
		Reason:              "Initial assignment",
		AssignedBy:          1,
	}

	err := repo.AssignPlayer(ctx, assignment)
	require.NoError(t, err)

	// Verify assignment
	current, err := repo.GetCurrentAssignment(ctx, player.ID)
	require.NoError(t, err)
	assert.Equal(t, company.ID, current.SettlementCompanyID)
	assert.True(t, current.IsCurrent)
}

// TestSettlementCompanyRepository_AssignPlayer_Reassignment tests reassigning a player
func TestSettlementCompanyRepository_AssignPlayer_Reassignment(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	company1 := CreateTestSettlementCompany(t, db, "First Company")
	company2 := CreateTestSettlementCompany(t, db, "Second Company")
	user := CreateUniqueTestUser(t, db, "player_to_reassign")
	player := CreateTestPlayer(t, db, user)

	// Initial assignment
	assignment1 := &model.PlayerCompanyAssignment{
		PlayerID:            player.ID,
		SettlementCompanyID: company1.ID,
		EffectiveDate:       time.Now(),
		Reason:              "Initial assignment",
		AssignedBy:          1,
	}
	require.NoError(t, repo.AssignPlayer(ctx, assignment1))

	// Reassign to company2
	time.Sleep(time.Millisecond) // Ensure different timestamp
	assignment2 := &model.PlayerCompanyAssignment{
		PlayerID:            player.ID,
		SettlementCompanyID: company2.ID,
		EffectiveDate:       time.Now(),
		Reason:              "Reassignment",
		AssignedBy:          1,
	}
	err := repo.AssignPlayer(ctx, assignment2)
	require.NoError(t, err)

	// Verify new assignment
	current, err := repo.GetCurrentAssignment(ctx, player.ID)
	require.NoError(t, err)
	assert.Equal(t, company2.ID, current.SettlementCompanyID)

	// Verify history
	history, err := repo.GetAssignmentHistory(ctx, player.ID)
	require.NoError(t, err)
	assert.Len(t, history, 2)
}

// TestSettlementCompanyRepository_GetCurrentAssignment tests getting current assignment
func TestSettlementCompanyRepository_GetCurrentAssignment(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	company := CreateTestSettlementCompany(t, db, "Current Assign Company")
	user := CreateUniqueTestUser(t, db, "player_current_assign")
	player := CreateTestPlayer(t, db, user)

	// No assignment initially
	_, err := repo.GetCurrentAssignment(ctx, player.ID)
	assert.Error(t, err)

	// Create assignment
	assignment := &model.PlayerCompanyAssignment{
		PlayerID:            player.ID,
		SettlementCompanyID: company.ID,
		EffectiveDate:       time.Now(),
		Reason:              "Test assignment",
		AssignedBy:          1,
	}
	require.NoError(t, repo.AssignPlayer(ctx, assignment))

	// Now should have current assignment
	current, err := repo.GetCurrentAssignment(ctx, player.ID)
	require.NoError(t, err)
	assert.Equal(t, company.ID, current.SettlementCompanyID)
	assert.NotNil(t, current.SettlementCompany)
}

// TestSettlementCompanyRepository_GetAssignmentHistory tests getting assignment history
func TestSettlementCompanyRepository_GetAssignmentHistory(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	company1 := CreateTestSettlementCompany(t, db, "History Company 1")
	company2 := CreateTestSettlementCompany(t, db, "History Company 2")
	user := CreateUniqueTestUser(t, db, "player_history")
	player := CreateTestPlayer(t, db, user)

	// First assignment
	assignment1 := &model.PlayerCompanyAssignment{
		PlayerID:            player.ID,
		SettlementCompanyID: company1.ID,
		EffectiveDate:       time.Now().Add(-30 * 24 * time.Hour),
		Reason:              "First assignment",
		AssignedBy:          1,
	}
	require.NoError(t, repo.AssignPlayer(ctx, assignment1))

	// Reassign
	time.Sleep(time.Millisecond)
	assignment2 := &model.PlayerCompanyAssignment{
		PlayerID:            player.ID,
		SettlementCompanyID: company2.ID,
		EffectiveDate:       time.Now(),
		Reason:              "Reassignment",
		AssignedBy:          1,
	}
	require.NoError(t, repo.AssignPlayer(ctx, assignment2))

	// Get history
	history, err := repo.GetAssignmentHistory(ctx, player.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(history), 2)
	// Should be ordered by effective_date DESC
	assert.Equal(t, company2.ID, history[0].SettlementCompanyID)
}

// TestSettlementCompanyRepository_EndCurrentAssignment tests ending current assignment
func TestSettlementCompanyRepository_EndCurrentAssignment(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	company := CreateTestSettlementCompany(t, db, "End Assign Company")
	user := CreateUniqueTestUser(t, db, "player_end_assign")
	player := CreateTestPlayer(t, db, user)

	// Create assignment
	assignment := &model.PlayerCompanyAssignment{
		PlayerID:            player.ID,
		SettlementCompanyID: company.ID,
		EffectiveDate:       time.Now(),
		Reason:              "Test assignment",
		AssignedBy:          1,
	}
	require.NoError(t, repo.AssignPlayer(ctx, assignment))

	// End assignment
	endDate := time.Now()
	err := repo.EndCurrentAssignment(ctx, player.ID, endDate)
	require.NoError(t, err)

	// Verify no current assignment
	_, err = repo.GetCurrentAssignment(ctx, player.ID)
	assert.Error(t, err)

	// Verify history shows ended assignment
	history, err := repo.GetAssignmentHistory(ctx, player.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(history), 1)
	assert.NotNil(t, history[0].EndDate)
}

// TestSettlementCompanyRepository_BatchAssignPlayers tests batch player assignment
func TestSettlementCompanyRepository_BatchAssignPlayers(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	company := CreateTestSettlementCompany(t, db, "Batch Assign Company")

	// Create multiple players
	var assignments []model.PlayerCompanyAssignment
	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, fmt.Sprintf("batch_player_%d", i))
		player := CreateTestPlayer(t, db, user)
		assignments = append(assignments, model.PlayerCompanyAssignment{
			PlayerID:            player.ID,
			SettlementCompanyID: company.ID,
			EffectiveDate:       time.Now(),
			Reason:              "Batch assignment",
			AssignedBy:          1,
		})
	}

	err := repo.BatchAssignPlayers(ctx, assignments)
	require.NoError(t, err)

	// Verify all assignments
	for _, assignment := range assignments {
		current, err := repo.GetCurrentAssignment(ctx, assignment.PlayerID)
		require.NoError(t, err)
		assert.Equal(t, company.ID, current.SettlementCompanyID)
	}

	// Verify company player count updated
	updatedCompany, err := repo.Get(ctx, company.ID)
	require.NoError(t, err)
	assert.Equal(t, 3, updatedCompany.PlayerCount)
}

// TestSettlementCompanyRepository_GetPlayerCount tests getting player count
func TestSettlementCompanyRepository_GetPlayerCount(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	company := CreateTestSettlementCompany(t, db, "Player Count Company")

	// Initially no players
	count, err := repo.GetPlayerCount(ctx, company.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	// Assign players
	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, fmt.Sprintf("count_player_%d", i))
		player := CreateTestPlayer(t, db, user)
		assignment := &model.PlayerCompanyAssignment{
			PlayerID:            player.ID,
			SettlementCompanyID: company.ID,
			EffectiveDate:       time.Now(),
			Reason:              "Count test",
			AssignedBy:          1,
		}
		require.NoError(t, repo.AssignPlayer(ctx, assignment))
	}

	// Check count
	count, err = repo.GetPlayerCount(ctx, company.ID)
	require.NoError(t, err)
	assert.Equal(t, 3, count)
}

// TestSettlementCompanyRepository_CreateHistory tests creating history record
func TestSettlementCompanyRepository_CreateHistory(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	company := CreateTestSettlementCompany(t, db, "History Test Company")

	history := &model.SettlementCompanyHistory{
		SettlementCompanyID: company.ID,
		FieldName:           "status",
		OldValue:            "active",
		NewValue:            "inactive",
		ChangedBy:           1,
	}

	err := repo.CreateHistory(ctx, history)
	require.NoError(t, err)
	assert.NotZero(t, history.ID)

	// Verify history
	histories, err := repo.GetHistory(ctx, company.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(histories), 1)
}

// TestSettlementCompanyRepository_GetHistory tests getting company history
func TestSettlementCompanyRepository_GetHistory(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	company := CreateTestSettlementCompany(t, db, "Get History Company")

	// Create multiple history records
	for i := 0; i < 3; i++ {
		history := &model.SettlementCompanyHistory{
			SettlementCompanyID: company.ID,
			FieldName:           fmt.Sprintf("field_%d", i),
			OldValue:            fmt.Sprintf("old_%d", i),
			NewValue:            fmt.Sprintf("new_%d", i),
			ChangedBy:           1,
		}
		require.NoError(t, repo.CreateHistory(ctx, history))
	}

	// Get history
	histories, err := repo.GetHistory(ctx, company.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(histories), 3)
}

// TestSettlementCompanyRepository_BatchUpdateStatus tests batch status update
func TestSettlementCompanyRepository_BatchUpdateStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	// Create companies
	var companyIDs []uint64
	for i := 0; i < 3; i++ {
		company := CreateTestSettlementCompany(t, db, fmt.Sprintf("Batch Status %d", i))
		companyIDs = append(companyIDs, company.ID)
	}

	// Batch update status
	err := repo.BatchUpdateStatus(ctx, companyIDs, model.CompanyStatusInactive)
	require.NoError(t, err)

	// Verify all updated
	for _, id := range companyIDs {
		company, err := repo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.CompanyStatusInactive, company.Status)
	}
}

// TestSettlementCompanyRepository_BatchDelete tests batch delete
func TestSettlementCompanyRepository_BatchDelete(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	// Create companies
	var companyIDs []uint64
	for i := 0; i < 3; i++ {
		company := CreateTestSettlementCompany(t, db, fmt.Sprintf("Batch Delete %d", i))
		companyIDs = append(companyIDs, company.ID)
	}

	// Batch delete
	err := repo.BatchDelete(ctx, companyIDs)
	require.NoError(t, err)

	// Verify all deleted
	for _, id := range companyIDs {
		_, err := repo.Get(ctx, id)
		assert.Error(t, err)
	}
}

// TestSettlementCompanyRepository_GetByIDsWithPlayerCount tests batch getting with player counts
func TestSettlementCompanyRepository_GetByIDsWithPlayerCount(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	ctx := context.Background()

	// Create companies
	var companyIDs []uint64
	for i := 0; i < 3; i++ {
		company := CreateTestSettlementCompany(t, db, fmt.Sprintf("Batch Get %d", i))
		companyIDs = append(companyIDs, company.ID)
	}

	// Batch get
	companies, err := repo.GetByIDsWithPlayerCount(ctx, companyIDs)
	require.NoError(t, err)
	assert.Len(t, companies, 3)
}
