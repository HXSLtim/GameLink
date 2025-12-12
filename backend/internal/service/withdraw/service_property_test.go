package withdraw_test

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository/settlementcompany"
	withdrawrepo "gamelink/internal/repository/withdraw"
	svc "gamelink/internal/service/withdraw"
	"gamelink/pkg/testutil"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/require"
)

// **Feature: payment-finance-module, Property 13: 提现分流一致性**
// **Validates: Requirements 13.1**
//
// Property 13: Withdrawal Routing Consistency
// *For any* player withdrawal request, the settlement company assigned by the system
// must match the player's currently effective settlement company assignment.

// TestProperty13_WithdrawRoutingConsistency tests that withdrawal routing matches player's current assignment
func TestProperty13_WithdrawRoutingConsistency(t *testing.T) {
	// Setup test database
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.SettlementCompany{},
		&model.PlayerCompanyAssignment{},
		&model.SettlementCompanyHistory{},
		&model.User{},
		&model.Player{},
		&model.Withdraw{},
		&model.SalaryPaymentRecord{},
		&model.Order{},
	)
	defer testutil.CleanDB(t, db)

	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	service := svc.NewWithdrawRoutingService(withdrawRepo, settlementRepo)
	ctx := context.Background()

	// Create test admin user
	testAdmin := &model.User{
		Name:         "test_admin",
		Email:        "admin@test.com",
		PasswordHash: "hashed_password",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
	}
	require.NoError(t, db.Create(testAdmin).Error)

	// Create multiple settlement companies
	companies := make([]*model.SettlementCompany, 3)
	for i := 0; i < 3; i++ {
		company := &model.SettlementCompany{
			Name:        genCompanyName(i),
			CreditCode:  genValidCreditCode(i),
			BankName:    "Test Bank",
			BankAccount: genBankAccount(i),
			Status:      model.CompanyStatusActive,
			CreatedBy:   testAdmin.ID,
		}
		require.NoError(t, db.Create(company).Error)
		companies[i] = company
	}

	// Create test players with their users
	players := make([]*model.Player, 5)
	for i := 0; i < 5; i++ {
		playerUser := &model.User{
			Name:         genPlayerName(i),
			Email:        genPlayerName(i) + "@test.com",
			Phone:        "1380000000" + string(rune('0'+i)),
			PasswordHash: "hashed_password",
			Role:         model.RolePlayer,
			Status:       model.UserStatusActive,
		}
		require.NoError(t, db.Create(playerUser).Error)

		player := &model.Player{
			UserID:             playerUser.ID,
			Nickname:           genPlayerName(i),
			VerificationStatus: model.VerificationVerified,
		}
		require.NoError(t, db.Create(player).Error)
		players[i] = player

		// Assign each player to a company
		companyIndex := i % len(companies)
		assignment := &model.PlayerCompanyAssignment{
			PlayerID:            player.ID,
			SettlementCompanyID: companies[companyIndex].ID,
			EffectiveDate:       time.Now().Add(-24 * time.Hour),
			Reason:              "Initial assignment",
			AssignedBy:          testAdmin.ID,
			IsCurrent:           true,
		}
		require.NoError(t, db.Create(assignment).Error)
	}

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property 13.1: Routed company must match player's current assignment
	properties.Property("withdrawal routing must match player's current settlement company", prop.ForAll(
		func(playerIndex int, amountCents int64) bool {
			player := players[playerIndex%len(players)]

			// Get player's current assignment
			var currentAssignment model.PlayerCompanyAssignment
			err := db.Where("player_id = ? AND is_current = ?", player.ID, true).
				First(&currentAssignment).Error
			if err != nil {
				return false
			}

			// Create a withdrawal request
			withdraw := &model.Withdraw{
				PlayerID:    player.ID,
				UserID:      player.UserID,
				AmountCents: amountCents,
				Method:      model.WithdrawMethodBank,
				AccountInfo: "test_account",
				Status:      model.WithdrawStatusPending,
			}

			// Route the withdrawal
			routedCompany, err := service.RouteWithdrawal(ctx, withdraw)
			if err != nil {
				return false
			}

			// Property: Routed company must match current assignment
			return routedCompany.ID == currentAssignment.SettlementCompanyID
		},
		gen.IntRange(0, 4),
		gen.Int64Range(1000, 1000000), // 10 to 10000 yuan in cents
	))

	// Property 13.2: Withdrawal should have correct routing info after routing
	properties.Property("withdrawal should have correct routing info after routing", prop.ForAll(
		func(playerIndex int) bool {
			player := players[playerIndex%len(players)]

			// Get player's current assignment
			var currentAssignment model.PlayerCompanyAssignment
			err := db.Where("player_id = ? AND is_current = ?", player.ID, true).
				First(&currentAssignment).Error
			if err != nil {
				return false
			}

			// Get the expected company
			var expectedCompany model.SettlementCompany
			err = db.First(&expectedCompany, currentAssignment.SettlementCompanyID).Error
			if err != nil {
				return false
			}

			// Create and route a withdrawal
			withdraw := &model.Withdraw{
				PlayerID:    player.ID,
				UserID:      player.UserID,
				AmountCents: 50000, // 500 yuan
				Method:      model.WithdrawMethodBank,
				AccountInfo: "test_account",
				Status:      model.WithdrawStatusPending,
			}

			_, err = service.RouteWithdrawal(ctx, withdraw)
			if err != nil {
				return false
			}

			// Verify routing info is set correctly
			if withdraw.SettlementCompanyID == nil {
				return false
			}
			if *withdraw.SettlementCompanyID != expectedCompany.ID {
				return false
			}
			if withdraw.SettlementCompanyName != expectedCompany.Name {
				return false
			}
			if withdraw.PaymentBankAccount != expectedCompany.BankAccount {
				return false
			}

			return true
		},
		gen.IntRange(0, 4),
	))

	// Property 13.3: Routing should fail for players without assignment
	// Create a single unassigned player for this test (outside the property function)
	unassignedUser := &model.User{
		Name:         "unassigned_player",
		Email:        "unassigned@test.com",
		Phone:        "13900000000",
		PasswordHash: "hashed_password",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
	}
	require.NoError(t, db.Create(unassignedUser).Error)

	unassignedPlayer := &model.Player{
		UserID:             unassignedUser.ID,
		Nickname:           "UnassignedPlayer",
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(unassignedPlayer).Error)

	properties.Property("routing should fail for players without settlement company", prop.ForAll(
		func(amountCents int64) bool {
			// Try to route withdrawal for unassigned player
			withdraw := &model.Withdraw{
				PlayerID:    unassignedPlayer.ID,
				UserID:      unassignedPlayer.UserID,
				AmountCents: amountCents,
				Method:      model.WithdrawMethodBank,
				AccountInfo: "test_account",
				Status:      model.WithdrawStatusPending,
			}

			_, err := service.RouteWithdrawal(ctx, withdraw)

			// Property: Should return error for unassigned player
			return err != nil
		},
		gen.Int64Range(1000, 100000),
	))

	properties.TestingRun(t)
}

// TestProperty13_WithdrawRoutingAfterReassignment tests routing after company reassignment
func TestProperty13_WithdrawRoutingAfterReassignment(t *testing.T) {
	// Setup test database
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.SettlementCompany{},
		&model.PlayerCompanyAssignment{},
		&model.SettlementCompanyHistory{},
		&model.User{},
		&model.Player{},
		&model.Withdraw{},
		&model.SalaryPaymentRecord{},
		&model.Order{},
	)
	defer testutil.CleanDB(t, db)

	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	service := svc.NewWithdrawRoutingService(withdrawRepo, settlementRepo)
	ctx := context.Background()

	// Create test admin
	testAdmin := &model.User{
		Name:         "reassign_admin",
		Email:        "reassign_admin@test.com",
		PasswordHash: "hashed_password",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
	}
	require.NoError(t, db.Create(testAdmin).Error)

	// Create two companies
	companyA := &model.SettlementCompany{
		Name:        "Company A",
		CreditCode:  "91110000300000000A",
		BankName:    "Bank A",
		BankAccount: "1111111111",
		Status:      model.CompanyStatusActive,
		CreatedBy:   testAdmin.ID,
	}
	require.NoError(t, db.Create(companyA).Error)

	companyB := &model.SettlementCompany{
		Name:        "Company B",
		CreditCode:  "91110000300000000B",
		BankName:    "Bank B",
		BankAccount: "2222222222",
		Status:      model.CompanyStatusActive,
		CreatedBy:   testAdmin.ID,
	}
	require.NoError(t, db.Create(companyB).Error)

	// Create test player
	playerUser := &model.User{
		Name:         "reassign_player",
		Email:        "reassign_player@test.com",
		Phone:        "13800000099",
		PasswordHash: "hashed_password",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
	}
	require.NoError(t, db.Create(playerUser).Error)

	player := &model.Player{
		UserID:             playerUser.ID,
		Nickname:           "ReassignPlayer",
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(player).Error)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	// Property 13.4: After reassignment, routing should use new company
	properties.Property("routing should use new company after reassignment", prop.ForAll(
		func(iterations int) bool {
			companies := []*model.SettlementCompany{companyA, companyB}

			for i := 0; i < iterations; i++ {
				targetCompany := companies[i%2]

				// End current assignment if exists
				db.Model(&model.PlayerCompanyAssignment{}).
					Where("player_id = ? AND is_current = ?", player.ID, true).
					Updates(map[string]interface{}{
						"is_current": false,
						"end_date":   time.Now(),
					})

				// Create new assignment
				assignment := &model.PlayerCompanyAssignment{
					PlayerID:            player.ID,
					SettlementCompanyID: targetCompany.ID,
					EffectiveDate:       time.Now(),
					Reason:              "Reassignment test",
					AssignedBy:          testAdmin.ID,
					IsCurrent:           true,
				}
				if err := db.Create(assignment).Error; err != nil {
					return false
				}

				// Route withdrawal
				withdraw := &model.Withdraw{
					PlayerID:    player.ID,
					UserID:      player.UserID,
					AmountCents: 10000,
					Method:      model.WithdrawMethodBank,
					AccountInfo: "test_account",
					Status:      model.WithdrawStatusPending,
				}

				routedCompany, err := service.RouteWithdrawal(ctx, withdraw)
				if err != nil {
					return false
				}

				// Verify routing matches new assignment
				if routedCompany.ID != targetCompany.ID {
					return false
				}
			}

			return true
		},
		gen.IntRange(2, 10),
	))

	properties.TestingRun(t)
}

// Helper functions
func genCompanyName(index int) string {
	return "Routing Test Company " + string(rune('A'+index))
}

func genValidCreditCode(index int) string {
	codes := []string{
		"91110000400000000A",
		"91110000400000000B",
		"91110000400000000C",
		"91110000400000000D",
		"91110000400000000E",
	}
	return codes[index%len(codes)]
}

func genBankAccount(index int) string {
	return "622000000000000" + string(rune('0'+index))
}

func genPlayerName(index int) string {
	return "RoutingTestPlayer" + string(rune('0'+index))
}

// **Feature: payment-finance-module, Property 14: 提现分流统计准确性**
// **Validates: Requirements 14.1**
//
// Property 14: Withdrawal Routing Statistics Accuracy
// *For any* settlement company's withdrawal statistics, the total withdrawal amount
// must equal the sum of all completed withdrawal records for that company.

// TestProperty14_WithdrawRoutingStatsAccuracy tests that routing statistics are accurate
func TestProperty14_WithdrawRoutingStatsAccuracy(t *testing.T) {
	// Setup test database
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.SettlementCompany{},
		&model.PlayerCompanyAssignment{},
		&model.SettlementCompanyHistory{},
		&model.User{},
		&model.Player{},
		&model.Withdraw{},
		&model.SalaryPaymentRecord{},
		&model.Order{},
	)
	defer testutil.CleanDB(t, db)

	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	ctx := context.Background()

	// Create test admin user
	testAdmin := &model.User{
		Name:         "stats_admin",
		Email:        "stats_admin@test.com",
		PasswordHash: "hashed_password",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
	}
	require.NoError(t, db.Create(testAdmin).Error)

	// Create settlement companies
	companies := make([]*model.SettlementCompany, 3)
	for i := 0; i < 3; i++ {
		company := &model.SettlementCompany{
			Name:        "Stats Test Company " + string(rune('A'+i)),
			CreditCode:  genStatsTestCreditCode(i),
			BankName:    "Test Bank",
			BankAccount: genStatsBankAccount(i),
			Status:      model.CompanyStatusActive,
			CreatedBy:   testAdmin.ID,
		}
		require.NoError(t, db.Create(company).Error)
		companies[i] = company
	}

	// Create test players
	players := make([]*model.Player, 5)
	for i := 0; i < 5; i++ {
		playerUser := &model.User{
			Name:         "StatsPlayer" + string(rune('0'+i)),
			Email:        "statsplayer" + string(rune('0'+i)) + "@test.com",
			Phone:        "1390000010" + string(rune('0'+i)),
			PasswordHash: "hashed_password",
			Role:         model.RolePlayer,
			Status:       model.UserStatusActive,
		}
		require.NoError(t, db.Create(playerUser).Error)

		player := &model.Player{
			UserID:             playerUser.ID,
			Nickname:           "StatsPlayer" + string(rune('0'+i)),
			VerificationStatus: model.VerificationVerified,
		}
		require.NoError(t, db.Create(player).Error)
		players[i] = player
	}

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	// Property 14.1: Total stats should equal sum of individual withdrawals
	properties.Property("total withdrawal stats should equal sum of completed withdrawals", prop.ForAll(
		func(withdrawalCounts []int, amounts []int64) bool {
			// Clean up previous test data
			db.Exec("DELETE FROM withdraws")

			// Create withdrawals for each company
			expectedTotals := make(map[uint64]int64)
			expectedCounts := make(map[uint64]int64)

			for i, company := range companies {
				count := withdrawalCounts[i%len(withdrawalCounts)] % 10 // Limit to 0-9 withdrawals
				if count < 0 {
					count = 0
				}

				for j := 0; j < count; j++ {
					amount := amounts[(i*10+j)%len(amounts)]
					if amount < 100 {
						amount = 100 // Minimum 1 yuan
					}
					if amount > 10000000 {
						amount = 10000000 // Maximum 100000 yuan
					}

					player := players[(i+j)%len(players)]
					now := time.Now()

					withdraw := &model.Withdraw{
						PlayerID:              player.ID,
						UserID:                player.UserID,
						AmountCents:           amount,
						Method:                model.WithdrawMethodBank,
						AccountInfo:           "test_account",
						Status:                model.WithdrawStatusCompleted,
						SettlementCompanyID:   &company.ID,
						SettlementCompanyName: company.Name,
						PaymentBankAccount:    company.BankAccount,
						TaxDeductedCents:      amount / 10, // 10% tax
						ActualAmountCents:     amount - amount/10,
						CompletedAt:           &now,
					}

					if err := db.Create(withdraw).Error; err != nil {
						return false
					}

					expectedTotals[company.ID] += amount
					expectedCounts[company.ID]++
				}
			}

			// Get routing stats
			stats, err := withdrawRepo.GetRoutingStats(ctx, nil, nil)
			if err != nil {
				return false
			}

			// Verify totals match
			var totalExpectedAmount int64
			var totalExpectedCount int64
			for _, amount := range expectedTotals {
				totalExpectedAmount += amount
			}
			for _, count := range expectedCounts {
				totalExpectedCount += count
			}

			if stats.TotalAmountCents != totalExpectedAmount {
				t.Logf("Total amount mismatch: expected %d, got %d", totalExpectedAmount, stats.TotalAmountCents)
				return false
			}

			if stats.TotalWithdrawals != totalExpectedCount {
				t.Logf("Total count mismatch: expected %d, got %d", totalExpectedCount, stats.TotalWithdrawals)
				return false
			}

			// Verify per-company stats
			for _, companyStat := range stats.ByCompany {
				expectedAmount := expectedTotals[companyStat.SettlementCompanyID]
				expectedCount := expectedCounts[companyStat.SettlementCompanyID]

				if companyStat.TotalAmountCents != expectedAmount {
					t.Logf("Company %d amount mismatch: expected %d, got %d",
						companyStat.SettlementCompanyID, expectedAmount, companyStat.TotalAmountCents)
					return false
				}

				if companyStat.TotalWithdrawals != expectedCount {
					t.Logf("Company %d count mismatch: expected %d, got %d",
						companyStat.SettlementCompanyID, expectedCount, companyStat.TotalWithdrawals)
					return false
				}
			}

			return true
		},
		gen.SliceOfN(3, gen.IntRange(0, 10)),
		gen.SliceOfN(30, gen.Int64Range(100, 10000000)),
	))

	// Property 14.2: Stats should only include completed withdrawals
	properties.Property("stats should only include completed withdrawals", prop.ForAll(
		func(completedCount, pendingCount int) bool {
			// Clean up previous test data
			db.Exec("DELETE FROM withdraws")

			company := companies[0]
			player := players[0]
			now := time.Now()

			// Create completed withdrawals
			var expectedTotal int64
			for i := 0; i < completedCount; i++ {
				amount := int64(10000 + i*1000)
				withdraw := &model.Withdraw{
					PlayerID:              player.ID,
					UserID:                player.UserID,
					AmountCents:           amount,
					Method:                model.WithdrawMethodBank,
					AccountInfo:           "test_account",
					Status:                model.WithdrawStatusCompleted,
					SettlementCompanyID:   &company.ID,
					SettlementCompanyName: company.Name,
					CompletedAt:           &now,
				}
				if err := db.Create(withdraw).Error; err != nil {
					return false
				}
				expectedTotal += amount
			}

			// Create pending withdrawals (should not be counted)
			for i := 0; i < pendingCount; i++ {
				withdraw := &model.Withdraw{
					PlayerID:              player.ID,
					UserID:                player.UserID,
					AmountCents:           int64(50000 + i*1000),
					Method:                model.WithdrawMethodBank,
					AccountInfo:           "test_account",
					Status:                model.WithdrawStatusPending,
					SettlementCompanyID:   &company.ID,
					SettlementCompanyName: company.Name,
				}
				if err := db.Create(withdraw).Error; err != nil {
					return false
				}
			}

			// Get stats
			stats, err := withdrawRepo.GetRoutingStats(ctx, nil, nil)
			if err != nil {
				return false
			}

			// Verify only completed withdrawals are counted
			if stats.TotalWithdrawals != int64(completedCount) {
				t.Logf("Count mismatch: expected %d completed, got %d", completedCount, stats.TotalWithdrawals)
				return false
			}

			if stats.TotalAmountCents != expectedTotal {
				t.Logf("Amount mismatch: expected %d, got %d", expectedTotal, stats.TotalAmountCents)
				return false
			}

			return true
		},
		gen.IntRange(0, 10),
		gen.IntRange(0, 10),
	))

	// Property 14.3: Percentage should sum to 100% (or 0% if no withdrawals)
	properties.Property("company percentages should sum to approximately 100%", prop.ForAll(
		func(amounts []int64) bool {
			// Clean up previous test data
			db.Exec("DELETE FROM withdraws")

			now := time.Now()
			hasWithdrawals := false

			// Create withdrawals for each company
			for i, company := range companies {
				amount := amounts[i%len(amounts)]
				if amount < 100 {
					continue // Skip if amount too small
				}
				if amount > 10000000 {
					amount = 10000000
				}

				player := players[i%len(players)]
				withdraw := &model.Withdraw{
					PlayerID:              player.ID,
					UserID:                player.UserID,
					AmountCents:           amount,
					Method:                model.WithdrawMethodBank,
					AccountInfo:           "test_account",
					Status:                model.WithdrawStatusCompleted,
					SettlementCompanyID:   &company.ID,
					SettlementCompanyName: company.Name,
					CompletedAt:           &now,
				}
				if err := db.Create(withdraw).Error; err != nil {
					return false
				}
				hasWithdrawals = true
			}

			if !hasWithdrawals {
				return true // No withdrawals, skip percentage check
			}

			// Get stats
			stats, err := withdrawRepo.GetRoutingStats(ctx, nil, nil)
			if err != nil {
				return false
			}

			// Sum percentages
			var totalPercentage float64
			for _, companyStat := range stats.ByCompany {
				totalPercentage += companyStat.Percentage
			}

			// Allow small floating point error
			if totalPercentage < 99.9 || totalPercentage > 100.1 {
				t.Logf("Percentage sum: %.2f%% (expected ~100%%)", totalPercentage)
				return false
			}

			return true
		},
		gen.SliceOfN(3, gen.Int64Range(100, 10000000)),
	))

	properties.TestingRun(t)
}

func genStatsTestCreditCode(index int) string {
	codes := []string{
		"91110000500000000A",
		"91110000500000000B",
		"91110000500000000C",
	}
	return codes[index%len(codes)]
}

func genStatsBankAccount(index int) string {
	return "622000000000001" + string(rune('0'+index))
}
