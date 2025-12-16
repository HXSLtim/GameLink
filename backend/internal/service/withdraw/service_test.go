package withdraw

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository/settlementcompany"
	withdrawrepo "gamelink/internal/repository/withdraw"
	"gamelink/pkg/testutil"
)

func setupWithdrawTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.User{},
		&model.Player{},
		&model.Withdraw{},
		&model.SettlementCompany{},
		&model.PlayerCompanyAssignment{},
	)
	return db
}

func createWithdrawTestData(t *testing.T, db *gorm.DB) (*model.Player, *model.SettlementCompany) {
	t.Helper()

	// 创建用户
	user := &model.User{
		Phone:        "13800000001",
		Email:        "player@test.com",
		Name:         "Test Player",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(user).Error)

	// 创建陪玩师
	player := &model.Player{
		UserID:             user.ID,
		Nickname:           "Pro Player",
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(player).Error)

	// 创建结算公司
	company := &model.SettlementCompany{
		Name:       "测试结算公司",
		CreditCode: "91110000MA00ABCD12",
		Status:     model.CompanyStatusActive,
		BankName:   "测试银行",
		BankBranch: "测试支行",
		CreatedBy:  user.ID,
	}
	require.NoError(t, db.Create(company).Error)

	// 创建分配关系
	now := time.Now()
	assignment := &model.PlayerCompanyAssignment{
		PlayerID:            player.ID,
		SettlementCompanyID: company.ID,
		EffectiveDate:       now,
		AssignedBy:          user.ID,
		IsCurrent:           true,
	}
	require.NoError(t, db.Create(assignment).Error)

	return player, company
}

func createWithdrawService(db *gorm.DB) *WithdrawRoutingService {
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	return NewWithdrawRoutingService(withdrawRepo, settlementRepo)
}

func createWithdrawStatsService(db *gorm.DB) *WithdrawRoutingStatsService {
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	return NewWithdrawRoutingStatsService(withdrawRepo, settlementRepo)
}

// uint64Ptr returns a pointer to the given uint64 value
func uint64Ptr(v uint64) *uint64 {
	return &v
}

func TestWithdrawRoutingService_RouteWithdrawal(t *testing.T) {
	db := setupWithdrawTestDB(t)
	defer testutil.CleanDB(t, db)

	player, company := createWithdrawTestData(t, db)
	svc := createWithdrawService(db)
	ctx := context.Background()

	t.Run("分流成功", func(t *testing.T) {
		withdraw := &model.Withdraw{
			PlayerID:    player.ID,
			AmountCents: 100000,
			Status:      model.WithdrawStatusPending,
		}
		require.NoError(t, db.Create(withdraw).Error)

		routedCompany, err := svc.RouteWithdrawal(ctx, withdraw)
		require.NoError(t, err)
		assert.Equal(t, company.ID, routedCompany.ID)
		assert.Equal(t, company.Name, routedCompany.Name)
		assert.NotNil(t, withdraw.SettlementCompanyID)
		assert.Equal(t, company.Name, withdraw.SettlementCompanyName)
	})

	t.Run("陪玩师无结算公司分配", func(t *testing.T) {
		// 创建没有分配的陪玩师
		user2 := &model.User{
			Phone:        "13800000002",
			Email:        "player2@test.com",
			Name:         "Test Player 2",
			Role:         model.RolePlayer,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(user2).Error)

		player2 := &model.Player{
			UserID:             user2.ID,
			Nickname:           "Player 2",
			VerificationStatus: model.VerificationVerified,
		}
		require.NoError(t, db.Create(player2).Error)

		withdraw := &model.Withdraw{
			PlayerID:    player2.ID,
			AmountCents: 50000,
			Status:      model.WithdrawStatusPending,
		}
		require.NoError(t, db.Create(withdraw).Error)

		_, err := svc.RouteWithdrawal(ctx, withdraw)
		assert.Error(t, err)
	})

	t.Run("结算公司已禁用", func(t *testing.T) {
		// 创建禁用的结算公司
		inactiveCompany := &model.SettlementCompany{
			Name:       "禁用公司",
			CreditCode: "91110000MA00EFGH34",
			Status:     model.CompanyStatusInactive,
			CreatedBy:  1,
		}
		require.NoError(t, db.Create(inactiveCompany).Error)

		// 创建新陪玩师并分配到禁用公司
		user3 := &model.User{
			Phone:        "13800000003",
			Email:        "player3@test.com",
			Name:         "Test Player 3",
			Role:         model.RolePlayer,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(user3).Error)

		player3 := &model.Player{
			UserID:             user3.ID,
			Nickname:           "Player 3",
			VerificationStatus: model.VerificationVerified,
		}
		require.NoError(t, db.Create(player3).Error)

		now := time.Now()
		assignment := &model.PlayerCompanyAssignment{
			PlayerID:            player3.ID,
			SettlementCompanyID: inactiveCompany.ID,
			EffectiveDate:       now,
			AssignedBy:          user3.ID,
			IsCurrent:           true,
		}
		require.NoError(t, db.Create(assignment).Error)

		withdraw := &model.Withdraw{
			PlayerID:    player3.ID,
			AmountCents: 50000,
			Status:      model.WithdrawStatusPending,
		}
		require.NoError(t, db.Create(withdraw).Error)

		_, err := svc.RouteWithdrawal(ctx, withdraw)
		assert.Error(t, err)
	})
}

func TestWithdrawRoutingService_ProcessWithdrawalRouting(t *testing.T) {
	db := setupWithdrawTestDB(t)
	defer testutil.CleanDB(t, db)

	player, company := createWithdrawTestData(t, db)
	svc := createWithdrawService(db)
	ctx := context.Background()

	t.Run("处理提现分流成功", func(t *testing.T) {
		companyID := company.ID
		withdraw := &model.Withdraw{
			PlayerID:              player.ID,
			AmountCents:           100000,
			Status:                model.WithdrawStatusApproved,
			SettlementCompanyID:   &companyID,
			SettlementCompanyName: company.Name,
		}
		require.NoError(t, db.Create(withdraw).Error)

		record, err := svc.ProcessWithdrawalRouting(ctx, withdraw.ID, 3000)
		require.NoError(t, err)
		assert.Equal(t, withdraw.ID, record.WithdrawID)
		assert.Equal(t, player.ID, record.PlayerID)
		assert.Equal(t, company.ID, record.SettlementCompanyID)
		assert.Equal(t, int64(100000), record.AmountCents)
		assert.Equal(t, int64(3000), record.TaxDeductedCents)
		assert.Equal(t, int64(97000), record.ActualAmountCents)
	})

	t.Run("未审批的提现无法处理", func(t *testing.T) {
		withdraw := &model.Withdraw{
			PlayerID:    player.ID,
			AmountCents: 50000,
			Status:      model.WithdrawStatusPending,
		}
		require.NoError(t, db.Create(withdraw).Error)

		_, err := svc.ProcessWithdrawalRouting(ctx, withdraw.ID, 1500)
		assert.Error(t, err)
	})

	t.Run("提现不存在", func(t *testing.T) {
		_, err := svc.ProcessWithdrawalRouting(ctx, 99999, 1000)
		assert.Error(t, err)
	})
}

func TestWithdrawRoutingService_CompleteWithdrawalPayment(t *testing.T) {
	db := setupWithdrawTestDB(t)
	defer testutil.CleanDB(t, db)

	player, company := createWithdrawTestData(t, db)
	svc := createWithdrawService(db)
	ctx := context.Background()

	t.Run("完成提现付款成功", func(t *testing.T) {
		companyID := company.ID
		withdraw := &model.Withdraw{
			PlayerID:              player.ID,
			AmountCents:           100000,
			Status:                model.WithdrawStatusApproved,
			SettlementCompanyID:   &companyID,
			SettlementCompanyName: company.Name,
		}
		require.NoError(t, db.Create(withdraw).Error)

		err := svc.CompleteWithdrawalPayment(ctx, withdraw.ID, "BANK_TX_123456")
		require.NoError(t, err)

		// 验证状态更新
		var updated model.Withdraw
		require.NoError(t, db.First(&updated, withdraw.ID).Error)
		assert.Equal(t, model.WithdrawStatusCompleted, updated.Status)
		assert.Equal(t, "BANK_TX_123456", updated.BankTransactionNo)
		assert.NotNil(t, updated.PaidAt)
		assert.NotNil(t, updated.CompletedAt)
	})

	t.Run("提现不存在", func(t *testing.T) {
		err := svc.CompleteWithdrawalPayment(ctx, 99999, "TX123")
		assert.Error(t, err)
	})
}

func TestWithdrawRoutingService_GetWithdrawal(t *testing.T) {
	db := setupWithdrawTestDB(t)
	defer testutil.CleanDB(t, db)

	player, _ := createWithdrawTestData(t, db)
	svc := createWithdrawService(db)
	ctx := context.Background()

	t.Run("获取提现记录成功", func(t *testing.T) {
		withdraw := &model.Withdraw{
			PlayerID:    player.ID,
			AmountCents: 100000,
			Status:      model.WithdrawStatusPending,
		}
		require.NoError(t, db.Create(withdraw).Error)

		found, err := svc.GetWithdrawal(ctx, withdraw.ID)
		require.NoError(t, err)
		assert.Equal(t, withdraw.ID, found.ID)
		assert.Equal(t, player.ID, found.PlayerID)
		assert.Equal(t, int64(100000), found.AmountCents)
	})

	t.Run("提现不存在", func(t *testing.T) {
		_, err := svc.GetWithdrawal(ctx, 99999)
		assert.Error(t, err)
	})
}

func TestWithdrawRoutingService_GetPlayerCurrentCompany(t *testing.T) {
	db := setupWithdrawTestDB(t)
	defer testutil.CleanDB(t, db)

	player, company := createWithdrawTestData(t, db)
	svc := createWithdrawService(db)
	ctx := context.Background()

	t.Run("获取陪玩师当前结算公司成功", func(t *testing.T) {
		found, err := svc.GetPlayerCurrentCompany(ctx, player.ID)
		require.NoError(t, err)
		assert.Equal(t, company.ID, found.ID)
		assert.Equal(t, company.Name, found.Name)
	})

	t.Run("陪玩师无结算公司", func(t *testing.T) {
		_, err := svc.GetPlayerCurrentCompany(ctx, 99999)
		assert.Error(t, err)
	})
}

func TestWithdrawRoutingStatsService_GetRoutingStats(t *testing.T) {
	db := setupWithdrawTestDB(t)
	defer testutil.CleanDB(t, db)

	player, company := createWithdrawTestData(t, db)
	svc := createWithdrawStatsService(db)
	ctx := context.Background()

	// 创建测试提现记录
	now := time.Now()
	for i := 0; i < 5; i++ {
		withdraw := &model.Withdraw{
			PlayerID:              player.ID,
			AmountCents:           int64((i + 1) * 10000),
			TaxDeductedCents:      int64((i + 1) * 300),
			ActualAmountCents:     int64((i+1)*10000 - (i+1)*300),
			Status:                model.WithdrawStatusCompleted,
			SettlementCompanyID:   uint64Ptr(company.ID),
			SettlementCompanyName: company.Name,
			CompletedAt:           &now,
		}
		require.NoError(t, db.Create(withdraw).Error)
	}

	t.Run("获取分流统计成功", func(t *testing.T) {
		dateFrom := now.AddDate(0, 0, -1)
		dateTo := now.AddDate(0, 0, 1)
		stats, err := svc.GetRoutingStats(ctx, &model.WithdrawRoutingStatsRequest{
			DateFrom: &dateFrom,
			DateTo:   &dateTo,
		})
		require.NoError(t, err)
		assert.NotNil(t, stats)
	})
}

func TestWithdrawRoutingStatsService_ListWithdrawalsByCompany(t *testing.T) {
	db := setupWithdrawTestDB(t)
	defer testutil.CleanDB(t, db)

	player, company := createWithdrawTestData(t, db)
	svc := createWithdrawStatsService(db)
	ctx := context.Background()

	// 创建测试提现记录
	for i := 0; i < 3; i++ {
		withdraw := &model.Withdraw{
			PlayerID:              player.ID,
			AmountCents:           int64((i + 1) * 10000),
			Status:                model.WithdrawStatusCompleted,
			SettlementCompanyID:   uint64Ptr(company.ID),
			SettlementCompanyName: company.Name,
		}
		require.NoError(t, db.Create(withdraw).Error)
	}

	t.Run("按公司查询提现列表", func(t *testing.T) {
		resp, err := svc.ListWithdrawalsByCompany(ctx, &model.ListWithdrawsByCompanyRequest{
			SettlementCompanyID: uint64Ptr(company.ID),
			Page:                1,
			PageSize:            10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(3), resp.Total)
		assert.Len(t, resp.Withdraws, 3)
	})

	t.Run("分页查询", func(t *testing.T) {
		resp, err := svc.ListWithdrawalsByCompany(ctx, &model.ListWithdrawsByCompanyRequest{
			SettlementCompanyID: uint64Ptr(company.ID),
			Page:                1,
			PageSize:            2,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(3), resp.Total)
		assert.Len(t, resp.Withdraws, 2)
	})

	t.Run("默认分页参数", func(t *testing.T) {
		resp, err := svc.ListWithdrawalsByCompany(ctx, &model.ListWithdrawsByCompanyRequest{
			SettlementCompanyID: uint64Ptr(company.ID),
			Page:                0,
			PageSize:            0,
		})
		require.NoError(t, err)
		assert.Equal(t, 1, resp.Page)
		assert.Equal(t, 20, resp.PageSize)
	})
}

func TestWithdrawRoutingStatsService_GenerateRoutingReport(t *testing.T) {
	db := setupWithdrawTestDB(t)
	defer testutil.CleanDB(t, db)

	player, company := createWithdrawTestData(t, db)
	svc := createWithdrawStatsService(db)
	ctx := context.Background()

	// 创建测试提现记录
	now := time.Now()
	for i := 0; i < 3; i++ {
		withdraw := &model.Withdraw{
			PlayerID:              player.ID,
			AmountCents:           int64((i + 1) * 10000),
			TaxDeductedCents:      int64((i + 1) * 300),
			ActualAmountCents:     int64((i+1)*10000 - (i+1)*300),
			Status:                model.WithdrawStatusCompleted,
			SettlementCompanyID:   uint64Ptr(company.ID),
			SettlementCompanyName: company.Name,
			CompletedAt:           &now,
		}
		require.NoError(t, db.Create(withdraw).Error)
	}

	t.Run("生成月度报表", func(t *testing.T) {
		report, err := svc.GenerateRoutingReport(ctx, &model.WithdrawRoutingReportRequest{
			ReportType: "monthly",
			Year:       now.Year(),
			Month:      int(now.Month()),
		})
		require.NoError(t, err)
		assert.Equal(t, "monthly", report.ReportType)
		assert.Equal(t, now.Year(), report.Year)
		assert.NotZero(t, report.GeneratedAt)
	})

	t.Run("生成季度报表", func(t *testing.T) {
		quarter := (int(now.Month())-1)/3 + 1
		report, err := svc.GenerateRoutingReport(ctx, &model.WithdrawRoutingReportRequest{
			ReportType: "quarterly",
			Year:       now.Year(),
			Quarter:    quarter,
		})
		require.NoError(t, err)
		assert.Equal(t, "quarterly", report.ReportType)
		assert.Equal(t, quarter, report.Quarter)
	})

	t.Run("生成年度报表", func(t *testing.T) {
		report, err := svc.GenerateRoutingReport(ctx, &model.WithdrawRoutingReportRequest{
			ReportType: "yearly",
			Year:       now.Year(),
		})
		require.NoError(t, err)
		assert.Equal(t, "yearly", report.ReportType)
	})

	t.Run("无效的报表类型", func(t *testing.T) {
		_, err := svc.GenerateRoutingReport(ctx, &model.WithdrawRoutingReportRequest{
			ReportType: "invalid",
			Year:       now.Year(),
		})
		assert.Error(t, err)
	})

	t.Run("无效的月份", func(t *testing.T) {
		_, err := svc.GenerateRoutingReport(ctx, &model.WithdrawRoutingReportRequest{
			ReportType: "monthly",
			Year:       now.Year(),
			Month:      13,
		})
		assert.Error(t, err)
	})

	t.Run("无效的季度", func(t *testing.T) {
		_, err := svc.GenerateRoutingReport(ctx, &model.WithdrawRoutingReportRequest{
			ReportType: "quarterly",
			Year:       now.Year(),
			Quarter:    5,
		})
		assert.Error(t, err)
	})
}

func TestWithdrawRoutingStatsService_GetCompanyWithdrawalStats(t *testing.T) {
	db := setupWithdrawTestDB(t)
	defer testutil.CleanDB(t, db)

	player, company := createWithdrawTestData(t, db)
	svc := createWithdrawStatsService(db)
	ctx := context.Background()

	// 创建测试提现记录
	now := time.Now()
	for i := 0; i < 3; i++ {
		withdraw := &model.Withdraw{
			PlayerID:              player.ID,
			AmountCents:           int64((i + 1) * 10000),
			TaxDeductedCents:      int64((i + 1) * 300),
			ActualAmountCents:     int64((i+1)*10000 - (i+1)*300),
			Status:                model.WithdrawStatusCompleted,
			SettlementCompanyID:   uint64Ptr(company.ID),
			SettlementCompanyName: company.Name,
			CompletedAt:           &now,
		}
		require.NoError(t, db.Create(withdraw).Error)
	}

	t.Run("获取公司提现统计成功", func(t *testing.T) {
		dateFrom := now.AddDate(0, 0, -1)
		dateTo := now.AddDate(0, 0, 1)
		stats, err := svc.GetCompanyWithdrawalStats(ctx, company.ID, &dateFrom, &dateTo)
		require.NoError(t, err)
		assert.NotNil(t, stats)
	})

	t.Run("公司不存在", func(t *testing.T) {
		_, err := svc.GetCompanyWithdrawalStats(ctx, 99999, nil, nil)
		assert.Error(t, err)
	})
}

// 测试错误变量
func TestWithdrawErrors(t *testing.T) {
	assert.NotNil(t, ErrNotFound)
	assert.NotNil(t, ErrNoSettlementCompany)
	assert.NotNil(t, ErrCompanyInactive)
	assert.NotNil(t, ErrInvalidAmount)
}
