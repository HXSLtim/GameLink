// Package integration provides integration tests for services.
package integration

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository/settlementcompany"
	withdrawrepo "gamelink/internal/repository/withdraw"
	withdrawservice "gamelink/internal/service/withdraw"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithdrawRoutingService_RouteWithdrawal(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)

	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	playerUser := CreateUniqueTestUser(t, db, "route_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create settlement company
	company := CreateTestSettlementCompany(t, db, "Test Company")

	// Create player company assignment
	assignment := &model.PlayerCompanyAssignment{
		PlayerID:            testPlayer.ID,
		SettlementCompanyID: company.ID,
		EffectiveDate:       time.Now().Add(-24 * time.Hour),
		IsCurrent:           true,
		AssignedBy:          playerUser.ID,
	}
	err := db.Create(assignment).Error
	require.NoError(t, err)

	// Create withdraw record
	withdraw := &model.Withdraw{
		PlayerID:    testPlayer.ID,
		AmountCents: 10000,
		Method:      model.WithdrawMethodAlipay,
		Status:      model.WithdrawStatusPending,
	}
	err = withdrawRepo.Create(ctx, withdraw)
	require.NoError(t, err)

	// Route withdrawal
	routedCompany, err := svc.RouteWithdrawal(ctx, withdraw)
	require.NoError(t, err)
	assert.Equal(t, company.ID, routedCompany.ID)
	assert.Equal(t, company.Name, routedCompany.Name)

	// Verify withdraw has routing info
	assert.NotNil(t, withdraw.SettlementCompanyID)
	assert.Equal(t, company.ID, *withdraw.SettlementCompanyID)
}

func TestWithdrawRoutingService_RouteWithdrawal_NoAssignment(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)

	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data without settlement assignment
	playerUser := CreateUniqueTestUser(t, db, "no_assign_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create withdraw record
	withdraw := &model.Withdraw{
		PlayerID:    testPlayer.ID,
		AmountCents: 10000,
		Method:      model.WithdrawMethodAlipay,
		Status:      model.WithdrawStatusPending,
	}
	err := withdrawRepo.Create(ctx, withdraw)
	require.NoError(t, err)

	// Try to route withdrawal
	_, err = svc.RouteWithdrawal(ctx, withdraw)
	assert.Error(t, err)
	assert.Equal(t, withdrawservice.ErrNoSettlementCompany, err)
}

func TestWithdrawRoutingService_RouteWithdrawal_InactiveCompany(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)

	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	playerUser := CreateUniqueTestUser(t, db, "inactive_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create inactive settlement company
	company := CreateTestSettlementCompany(t, db, "Inactive Company")
	company.Status = model.CompanyStatusInactive
	db.Save(company)

	// Create player company assignment
	assignment := &model.PlayerCompanyAssignment{
		PlayerID:            testPlayer.ID,
		SettlementCompanyID: company.ID,
		EffectiveDate:       time.Now().Add(-24 * time.Hour),
		IsCurrent:           true,
		AssignedBy:          playerUser.ID,
	}
	err := db.Create(assignment).Error
	require.NoError(t, err)

	// Create withdraw record
	withdraw := &model.Withdraw{
		PlayerID:    testPlayer.ID,
		AmountCents: 10000,
		Method:      model.WithdrawMethodAlipay,
		Status:      model.WithdrawStatusPending,
	}
	err = withdrawRepo.Create(ctx, withdraw)
	require.NoError(t, err)

	// Try to route withdrawal
	_, err = svc.RouteWithdrawal(ctx, withdraw)
	assert.Error(t, err)
	assert.Equal(t, withdrawservice.ErrCompanyInactive, err)
}

func TestWithdrawRoutingService_GetWithdrawal(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)

	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	playerUser := CreateUniqueTestUser(t, db, "get_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create withdraw record
	withdraw := &model.Withdraw{
		PlayerID:    testPlayer.ID,
		AmountCents: 15000,
		Method:      model.WithdrawMethodWeChat,
		Status:      model.WithdrawStatusPending,
	}
	err := withdrawRepo.Create(ctx, withdraw)
	require.NoError(t, err)

	// Get withdrawal
	result, err := svc.GetWithdrawal(ctx, withdraw.ID)
	require.NoError(t, err)
	assert.Equal(t, withdraw.ID, result.ID)
	assert.Equal(t, testPlayer.ID, result.PlayerID)
	assert.Equal(t, int64(15000), result.AmountCents)
	assert.Equal(t, model.WithdrawMethodWeChat, result.Method)
}

func TestWithdrawRoutingService_GetWithdrawal_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)

	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Try to get non-existent withdrawal
	_, err := svc.GetWithdrawal(ctx, 99999)
	assert.Error(t, err)
	assert.Equal(t, withdrawservice.ErrNotFound, err)
}

func TestWithdrawRoutingService_CompleteWithdrawalPayment(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)

	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	playerUser := CreateUniqueTestUser(t, db, "complete_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create approved withdraw record
	withdraw := &model.Withdraw{
		PlayerID:          testPlayer.ID,
		AmountCents:       20000,
		Method:            model.WithdrawMethodBank,
		Status:            model.WithdrawStatusApproved,
		TaxDeductedCents:  2000,
		ActualAmountCents: 18000,
	}
	err := withdrawRepo.Create(ctx, withdraw)
	require.NoError(t, err)

	// Complete withdrawal payment
	bankTransactionNo := "BANK_TXN_123456"
	err = svc.CompleteWithdrawalPayment(ctx, withdraw.ID, bankTransactionNo)
	require.NoError(t, err)

	// Verify withdrawal status
	result, err := svc.GetWithdrawal(ctx, withdraw.ID)
	require.NoError(t, err)
	assert.Equal(t, model.WithdrawStatusCompleted, result.Status)
	assert.Equal(t, bankTransactionNo, result.BankTransactionNo)
	assert.NotNil(t, result.PaidAt)
	assert.NotNil(t, result.CompletedAt)
}

func TestWithdrawRoutingService_GetPlayerCurrentCompany(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)

	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	playerUser := CreateUniqueTestUser(t, db, "current_company_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create settlement company
	company := CreateTestSettlementCompany(t, db, "Current Company")

	// Create player company assignment
	assignment := &model.PlayerCompanyAssignment{
		PlayerID:            testPlayer.ID,
		SettlementCompanyID: company.ID,
		EffectiveDate:       time.Now().Add(-24 * time.Hour),
		IsCurrent:           true,
		AssignedBy:          playerUser.ID,
	}
	err := db.Create(assignment).Error
	require.NoError(t, err)

	// Get player's current company
	result, err := svc.GetPlayerCurrentCompany(ctx, testPlayer.ID)
	require.NoError(t, err)
	assert.Equal(t, company.ID, result.ID)
	assert.Equal(t, company.Name, result.Name)
}

func TestWithdrawRoutingStatsService_GetRoutingStats(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)

	svc := withdrawservice.NewWithdrawRoutingStatsService(withdrawRepo, settlementRepo)

	// Create test data
	playerUser := CreateUniqueTestUser(t, db, "stats_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create settlement company
	company := CreateTestSettlementCompany(t, db, "Stats Company")

	// Create completed withdrawals
	now := time.Now()
	for i := 0; i < 3; i++ {
		withdraw := &model.Withdraw{
			PlayerID:              testPlayer.ID,
			AmountCents:           int64(10000 + i*5000),
			Method:                model.WithdrawMethodAlipay,
			Status:                model.WithdrawStatusCompleted,
			SettlementCompanyID:   &company.ID,
			SettlementCompanyName: company.Name,
			TaxDeductedCents:      int64(1000 + i*500),
			ActualAmountCents:     int64(9000 + i*4500),
			CompletedAt:           &now,
		}
		err := withdrawRepo.Create(ctx, withdraw)
		require.NoError(t, err)
	}

	// Get routing stats
	dateFrom := now.Add(-24 * time.Hour)
	dateTo := now.Add(24 * time.Hour)
	req := &model.WithdrawRoutingStatsRequest{
		DateFrom: &dateFrom,
		DateTo:   &dateTo,
	}

	stats, err := svc.GetRoutingStats(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, int64(3), stats.TotalWithdrawals)
	assert.Greater(t, stats.TotalAmountCents, int64(0))
	assert.Greater(t, stats.TotalTaxDeductedCents, int64(0))
	assert.Greater(t, stats.TotalActualAmountCents, int64(0))
}

func TestWithdrawRoutingStatsService_ListWithdrawalsByCompany(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)

	svc := withdrawservice.NewWithdrawRoutingStatsService(withdrawRepo, settlementRepo)

	// Create test data
	playerUser := CreateUniqueTestUser(t, db, "list_by_company_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create settlement company
	company := CreateTestSettlementCompany(t, db, "List Company")

	// Create withdrawals for this company
	for i := 0; i < 5; i++ {
		withdraw := &model.Withdraw{
			PlayerID:              testPlayer.ID,
			AmountCents:           int64(10000 + i*1000),
			Method:                model.WithdrawMethodAlipay,
			Status:                model.WithdrawStatusCompleted,
			SettlementCompanyID:   &company.ID,
			SettlementCompanyName: company.Name,
		}
		err := withdrawRepo.Create(ctx, withdraw)
		require.NoError(t, err)
	}

	// List withdrawals by company
	req := &model.ListWithdrawsByCompanyRequest{
		SettlementCompanyID: &company.ID,
		Page:                1,
		PageSize:            10,
	}

	resp, err := svc.ListWithdrawalsByCompany(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, int64(5), resp.Total)
	assert.Len(t, resp.Withdraws, 5)
}

func TestWithdrawRoutingStatsService_GenerateRoutingReport_Monthly(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)

	svc := withdrawservice.NewWithdrawRoutingStatsService(withdrawRepo, settlementRepo)

	// Create test data
	playerUser := CreateUniqueTestUser(t, db, "report_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create settlement company
	company := CreateTestSettlementCompany(t, db, "Report Company")

	// Create completed withdrawals
	now := time.Now()
	for i := 0; i < 3; i++ {
		withdraw := &model.Withdraw{
			PlayerID:              testPlayer.ID,
			AmountCents:           int64(10000),
			Method:                model.WithdrawMethodAlipay,
			Status:                model.WithdrawStatusCompleted,
			SettlementCompanyID:   &company.ID,
			SettlementCompanyName: company.Name,
			TaxDeductedCents:      1000,
			ActualAmountCents:     9000,
			CompletedAt:           &now,
		}
		err := withdrawRepo.Create(ctx, withdraw)
		require.NoError(t, err)
	}

	// Generate monthly report
	req := &model.WithdrawRoutingReportRequest{
		ReportType: "monthly",
		Year:       now.Year(),
		Month:      int(now.Month()),
	}

	report, err := svc.GenerateRoutingReport(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, "monthly", report.ReportType)
	assert.Equal(t, now.Year(), report.Year)
	assert.Equal(t, int(now.Month()), report.Month)
	assert.Equal(t, int64(3), report.TotalWithdrawals)
	assert.Equal(t, int64(30000), report.TotalAmountCents)
}

func TestWithdrawRoutingStatsService_GetCompanyWithdrawalStats(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)

	svc := withdrawservice.NewWithdrawRoutingStatsService(withdrawRepo, settlementRepo)

	// Create test data
	playerUser := CreateUniqueTestUser(t, db, "company_stats_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create settlement company
	company := CreateTestSettlementCompany(t, db, "Company Stats")

	// Create completed withdrawals
	now := time.Now()
	for i := 0; i < 4; i++ {
		withdraw := &model.Withdraw{
			PlayerID:              testPlayer.ID,
			AmountCents:           int64(5000),
			Method:                model.WithdrawMethodAlipay,
			Status:                model.WithdrawStatusCompleted,
			SettlementCompanyID:   &company.ID,
			SettlementCompanyName: company.Name,
			TaxDeductedCents:      500,
			ActualAmountCents:     4500,
			CompletedAt:           &now,
		}
		err := withdrawRepo.Create(ctx, withdraw)
		require.NoError(t, err)
	}

	// Get company withdrawal stats
	dateFrom := now.Add(-24 * time.Hour)
	dateTo := now.Add(24 * time.Hour)

	stats, err := svc.GetCompanyWithdrawalStats(ctx, company.ID, &dateFrom, &dateTo)
	require.NoError(t, err)
	assert.Equal(t, company.ID, stats.SettlementCompanyID)
	assert.Equal(t, int64(4), stats.TotalWithdrawals)
	assert.Equal(t, int64(20000), stats.TotalAmountCents)
	assert.Equal(t, int64(2000), stats.TotalTaxDeductedCents)
	assert.Equal(t, int64(18000), stats.TotalActualAmountCents)
}

func TestWithdrawRepository_GetPlayerBalance(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repository
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)

	// Create test data
	playerUser := CreateUniqueTestUser(t, db, "balance_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create completed orders for earnings
	testGame := CreateTestGame(t, db, "balance_game")
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, CreateUniqueTestUser(t, db, "order_user"), testPlayer, testGame, model.OrderStatusCompleted, 10000)
		_ = order
	}

	// Create completed withdrawals
	for i := 0; i < 2; i++ {
		withdraw := &model.Withdraw{
			PlayerID:    testPlayer.ID,
			AmountCents: 5000,
			Method:      model.WithdrawMethodAlipay,
			Status:      model.WithdrawStatusCompleted,
		}
		err := withdrawRepo.Create(ctx, withdraw)
		require.NoError(t, err)
	}

	// Get player balance
	balance, err := withdrawRepo.GetPlayerBalance(ctx, testPlayer.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(30000), balance.TotalEarnings)    // 3 orders * 10000
	assert.Equal(t, int64(10000), balance.WithdrawTotal)    // 2 withdrawals * 5000
	assert.Equal(t, int64(0), balance.PendingWithdraw)      // No pending withdrawals
	assert.Equal(t, int64(20000), balance.AvailableBalance) // 30000 - 10000
}

