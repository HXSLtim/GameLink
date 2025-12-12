package settlementcompany_test

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository/settlementcompany"
	svc "gamelink/internal/service/settlementcompany"
	"gamelink/pkg/testutil"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/require"
)

// **Feature: payment-finance-module, Property 12: 陪玩师结算公司分配唯一性**
// **Validates: Requirements 12.1, 12.4**
//
// Property 12: Player Settlement Company Assignment Uniqueness
// *For any* point in time, each player can only have one active settlement company assignment.
// When a new assignment takes effect, the old assignment automatically ends.

// TestProperty12_PlayerAssignmentUniqueness tests that each player has at most one current assignment
func TestProperty12_PlayerAssignmentUniqueness(t *testing.T) {
	// Setup test database
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.SettlementCompany{},
		&model.PlayerCompanyAssignment{},
		&model.SettlementCompanyHistory{},
		&model.User{},
		&model.Player{},
	)
	defer testutil.CleanDB(t, db)

	repo := settlementcompany.NewSettlementCompanyRepository(db)
	service := svc.NewSettlementCompanyService(repo, nil)
	ctx := context.Background()

	// Create test user for createdBy field
	testUser := &model.User{
		Name:         "test_admin",
		Email:        "admin@test.com",
		PasswordHash: "hashed_password",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
	}
	require.NoError(t, db.Create(testUser).Error)

	// Create test player
	testPlayer := &model.Player{
		UserID:             testUser.ID,
		Nickname:           "TestPlayer",
		VerificationStatus: model.VerificationPending,
	}
	require.NoError(t, db.Create(testPlayer).Error)

	// Create multiple settlement companies for testing
	companies := make([]*model.SettlementCompany, 3)
	for i := 0; i < 3; i++ {
		company, err := service.CreateCompany(ctx, &model.CreateSettlementCompanyRequest{
			Name:       genCompanyName(i),
			CreditCode: genValidCreditCodeForTest(i),
		}, testUser.ID)
		require.NoError(t, err)
		companies[i] = company
	}

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property 12.1: After any number of assignments, player should have at most one current assignment
	properties.Property("player should have at most one current assignment after multiple assignments", prop.ForAll(
		func(assignmentCount int) bool {
			// Perform multiple assignments to different companies
			for i := 0; i < assignmentCount; i++ {
				companyIndex := i % len(companies)
				_, err := service.AssignPlayerToCompany(ctx, &model.AssignPlayerToCompanyRequest{
					PlayerID:            testPlayer.ID,
					SettlementCompanyID: companies[companyIndex].ID,
					EffectiveDate:       time.Now(),
					Reason:              "Property test assignment",
				}, testUser.ID)
				if err != nil {
					return false
				}
			}

			// Count current assignments for this player
			var currentCount int64
			err := db.Model(&model.PlayerCompanyAssignment{}).
				Where("player_id = ? AND is_current = ?", testPlayer.ID, true).
				Count(&currentCount).Error
			if err != nil {
				return false
			}

			// Property: At most one current assignment
			return currentCount <= 1
		},
		gen.IntRange(1, 10),
	))

	// Property 12.2: New assignment should end previous assignment
	properties.Property("new assignment should end previous assignment", prop.ForAll(
		func(companyIndex int) bool {
			// Get current assignment before new assignment
			oldAssignment, _ := service.GetCurrentAssignment(ctx, testPlayer.ID)

			// Make new assignment to a different company
			newCompanyIndex := (companyIndex + 1) % len(companies)
			_, err := service.AssignPlayerToCompany(ctx, &model.AssignPlayerToCompanyRequest{
				PlayerID:            testPlayer.ID,
				SettlementCompanyID: companies[newCompanyIndex].ID,
				EffectiveDate:       time.Now(),
				Reason:              "Property test - new assignment",
			}, testUser.ID)
			if err != nil {
				return false
			}

			// If there was an old assignment, verify it's no longer current
			if oldAssignment != nil {
				var oldRecord model.PlayerCompanyAssignment
				err := db.First(&oldRecord, oldAssignment.ID).Error
				if err != nil {
					return false
				}
				// Old assignment should no longer be current
				if oldRecord.IsCurrent {
					return false
				}
				// Old assignment should have an end date
				if oldRecord.EndDate == nil {
					return false
				}
			}

			return true
		},
		gen.IntRange(0, 2),
	))

	// Property 12.3: Current assignment should match the last assigned company
	properties.Property("current assignment should match last assigned company", prop.ForAll(
		func(companyIndex int) bool {
			targetCompany := companies[companyIndex%len(companies)]

			// Assign to target company
			_, err := service.AssignPlayerToCompany(ctx, &model.AssignPlayerToCompanyRequest{
				PlayerID:            testPlayer.ID,
				SettlementCompanyID: targetCompany.ID,
				EffectiveDate:       time.Now(),
				Reason:              "Property test - verify current",
			}, testUser.ID)
			if err != nil {
				return false
			}

			// Get current assignment
			current, err := service.GetCurrentAssignment(ctx, testPlayer.ID)
			if err != nil {
				return false
			}

			// Current assignment should be for the target company
			return current.SettlementCompanyID == targetCompany.ID
		},
		gen.IntRange(0, 2),
	))

	properties.TestingRun(t)
}

// genCompanyName generates a unique company name for testing
func genCompanyName(index int) string {
	return "Test Company " + string(rune('A'+index))
}

// genValidCreditCodeForTest generates a valid credit code for testing
func genValidCreditCodeForTest(index int) string {
	// Generate unique valid credit codes
	codes := []string{
		"91110000100000000A",
		"91110000100000000B",
		"91110000100000000C",
		"91110000100000000D",
		"91110000100000000E",
	}
	return codes[index%len(codes)]
}

// TestProperty12_BatchAssignmentUniqueness tests batch assignment maintains uniqueness
func TestProperty12_BatchAssignmentUniqueness(t *testing.T) {
	// Setup test database
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.SettlementCompany{},
		&model.PlayerCompanyAssignment{},
		&model.SettlementCompanyHistory{},
		&model.User{},
		&model.Player{},
	)
	defer testutil.CleanDB(t, db)

	repo := settlementcompany.NewSettlementCompanyRepository(db)
	service := svc.NewSettlementCompanyService(repo, nil)
	ctx := context.Background()

	// Create test user
	testUser := &model.User{
		Name:         "batch_test_admin",
		Email:        "batch_admin@test.com",
		PasswordHash: "hashed_password",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
	}
	require.NoError(t, db.Create(testUser).Error)

	// Create multiple test players (each needs a unique user)
	players := make([]*model.Player, 5)
	for i := 0; i < 5; i++ {
		// Create a unique user for each player
		playerUser := &model.User{
			Name:         genPlayerNickname(i),
			Email:        genPlayerNickname(i) + "@test.com",
			Phone:        "1380000000" + string(rune('0'+i)),
			PasswordHash: "hashed_password",
			Role:         model.RolePlayer,
			Status:       model.UserStatusActive,
		}
		require.NoError(t, db.Create(playerUser).Error)

		player := &model.Player{
			UserID:             playerUser.ID,
			Nickname:           genPlayerNickname(i),
			VerificationStatus: model.VerificationPending,
		}
		require.NoError(t, db.Create(player).Error)
		players[i] = player
	}

	// Create settlement companies
	companies := make([]*model.SettlementCompany, 2)
	for i := 0; i < 2; i++ {
		company, err := service.CreateCompany(ctx, &model.CreateSettlementCompanyRequest{
			Name:       "Batch Test Company " + string(rune('A'+i)),
			CreditCode: genValidCreditCodeForBatchTest(i),
		}, testUser.ID)
		require.NoError(t, err)
		companies[i] = company
	}

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	// Property 12.4: After batch assignment, each player should have exactly one current assignment
	properties.Property("batch assignment should maintain uniqueness for all players", prop.ForAll(
		func(companyIndex int) bool {
			targetCompany := companies[companyIndex%len(companies)]
			playerIDs := make([]uint64, len(players))
			for i, p := range players {
				playerIDs[i] = p.ID
			}

			// Batch assign all players
			_, err := service.BatchAssignPlayers(ctx, &model.BatchAssignPlayersRequest{
				PlayerIDs:           playerIDs,
				SettlementCompanyID: targetCompany.ID,
				EffectiveDate:       time.Now(),
				Reason:              "Batch property test",
			}, testUser.ID)
			if err != nil {
				return false
			}

			// Verify each player has exactly one current assignment
			for _, player := range players {
				var currentCount int64
				err := db.Model(&model.PlayerCompanyAssignment{}).
					Where("player_id = ? AND is_current = ?", player.ID, true).
					Count(&currentCount).Error
				if err != nil || currentCount != 1 {
					return false
				}
			}

			return true
		},
		gen.IntRange(0, 1),
	))

	properties.TestingRun(t)
}

func genPlayerNickname(index int) string {
	return "BatchTestPlayer" + string(rune('0'+index))
}

func genValidCreditCodeForBatchTest(index int) string {
	codes := []string{
		"91110000200000000A",
		"91110000200000000B",
	}
	return codes[index%len(codes)]
}
