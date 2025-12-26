package withdraw

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/settlementcompany"
	withdrawrepo "gamelink/internal/repository/withdraw"
	"gamelink/pkg/apierr"
)

var (
	// ErrNotFound 提现记录不存在
	ErrNotFound = apierr.NotFound("withdrawal not found")
	// ErrNoSettlementCompany 陪玩师未分配结算公司
	ErrNoSettlementCompany = apierr.BadRequest("player has no settlement company assigned")
	// ErrCompanyInactive 结算公司已禁用
	ErrCompanyInactive = apierr.BadRequest("settlement company is inactive")
	// ErrInvalidAmount 无效金额
	ErrInvalidAmount = apierr.BadRequest("invalid withdrawal amount")
)

// WithdrawRoutingService 提现分流服务
// Requirements: 13.1, 13.3, 13.4, 13.5
type WithdrawRoutingService struct {
	withdrawRepo   withdrawrepo.WithdrawRepository
	settlementRepo settlementcompany.SettlementCompanyRepository
}

// NewWithdrawRoutingService 创建提现分流服务
func NewWithdrawRoutingService(
	withdrawRepo withdrawrepo.WithdrawRepository,
	settlementRepo settlementcompany.SettlementCompanyRepository,
) *WithdrawRoutingService {
	return &WithdrawRoutingService{
		withdrawRepo:   withdrawRepo,
		settlementRepo: settlementRepo,
	}
}

// RouteWithdrawal 根据陪玩师所属公司确定提现处理主体
// Requirements: 13.1
// Property 13: 提现分流一致性 - 系统分配的结算公司必须与陪玩师当前生效的结算公司分配一致
func (s *WithdrawRoutingService) RouteWithdrawal(ctx context.Context, withdraw *model.Withdraw) (*model.SettlementCompany, error) {
	// 获取陪玩师当前的结算公司分配
	assignment, err := s.settlementRepo.GetCurrentAssignment(ctx, withdraw.PlayerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNoSettlementCompany.WithDetails(fmt.Sprintf("player %d has no settlement company assigned", withdraw.PlayerID))
		}
		return nil, apierr.InternalError("failed to get player assignment").WithDetails(err.Error())
	}

	// 获取结算公司信息
	company, err := s.settlementRepo.Get(ctx, assignment.SettlementCompanyID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apierr.InternalError("settlement company not found").WithDetails(err.Error())
		}
		return nil, apierr.InternalError("failed to get settlement company").WithDetails(err.Error())
	}

	// 验证结算公司状态
	if !company.IsActive() {
		return nil, ErrCompanyInactive.WithDetails(fmt.Sprintf("settlement company %s is inactive", company.Name))
	}

	// 设置提现分流信息
	withdraw.SetRoutingInfo(company)

	return company, nil
}

// ProcessWithdrawalRouting 处理提现分流并生成工资发放记录
// Requirements: 13.3, 13.4, 13.5
func (s *WithdrawRoutingService) ProcessWithdrawalRouting(ctx context.Context, withdrawID uint64, taxDeductedCents int64) (*model.SalaryPaymentRecord, error) {
	// 获取提现记录
	withdraw, err := s.withdrawRepo.Get(ctx, withdrawID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("failed to get withdrawal").WithDetails(err.Error())
	}

	// 验证提现状态
	if withdraw.Status != model.WithdrawStatusApproved {
		return nil, apierr.BadRequest("withdrawal must be approved before processing")
	}

	// 如果还没有分配结算公司，进行分流
	if withdraw.SettlementCompanyID == nil {
		_, err := s.RouteWithdrawal(ctx, withdraw)
		if err != nil {
			return nil, err
		}
	}

	// 设置代扣个税和实际到账金额
	withdraw.TaxDeductedCents = taxDeductedCents
	withdraw.ActualAmountCents = withdraw.CalculateActualAmount()

	// 更新提现记录
	if err := s.withdrawRepo.Update(ctx, withdraw); err != nil {
		return nil, apierr.InternalError("failed to update withdrawal").WithDetails(err.Error())
	}

	// 生成工资发放记录
	record := &model.SalaryPaymentRecord{
		WithdrawID:            withdraw.ID,
		PlayerID:              withdraw.PlayerID,
		SettlementCompanyID:   *withdraw.SettlementCompanyID,
		SettlementCompanyName: withdraw.SettlementCompanyName,
		AmountCents:           withdraw.AmountCents,
		TaxDeductedCents:      withdraw.TaxDeductedCents,
		ActualAmountCents:     withdraw.ActualAmountCents,
		BankAccount:           withdraw.PaymentBankAccount,
		Status:                "pending",
	}

	return record, nil
}

// CompleteWithdrawalPayment 完成提现付款
// Requirements: 13.5
func (s *WithdrawRoutingService) CompleteWithdrawalPayment(ctx context.Context, withdrawID uint64, bankTransactionNo string) error {
	withdraw, err := s.withdrawRepo.Get(ctx, withdrawID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return apierr.InternalError("failed to get withdrawal").WithDetails(err.Error())
	}

	now := time.Now()
	withdraw.BankTransactionNo = bankTransactionNo
	withdraw.PaidAt = &now
	withdraw.Status = model.WithdrawStatusCompleted
	withdraw.CompletedAt = &now

	if err := s.withdrawRepo.Update(ctx, withdraw); err != nil {
		return apierr.InternalError("failed to update withdrawal").WithDetails(err.Error())
	}

	return nil
}

// GetWithdrawal 获取提现记录
func (s *WithdrawRoutingService) GetWithdrawal(ctx context.Context, id uint64) (*model.Withdraw, error) {
	withdraw, err := s.withdrawRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("failed to get withdrawal").WithDetails(err.Error())
	}
	return withdraw, nil
}

// GetPlayerCurrentCompany 获取陪玩师当前的结算公司
func (s *WithdrawRoutingService) GetPlayerCurrentCompany(ctx context.Context, playerID uint64) (*model.SettlementCompany, error) {
	assignment, err := s.settlementRepo.GetCurrentAssignment(ctx, playerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNoSettlementCompany
		}
		return nil, apierr.InternalError("failed to get player assignment").WithDetails(err.Error())
	}

	company, err := s.settlementRepo.Get(ctx, assignment.SettlementCompanyID)
	if err != nil {
		return nil, apierr.InternalError("failed to get settlement company").WithDetails(err.Error())
	}

	return company, nil
}

// WithdrawRoutingStatsService 提现分流统计服务
// Requirements: 14.1, 14.2, 14.3, 14.4, 14.5
type WithdrawRoutingStatsService struct {
	withdrawRepo   withdrawrepo.WithdrawRepository
	settlementRepo settlementcompany.SettlementCompanyRepository
}

// NewWithdrawRoutingStatsService 创建提现分流统计服务
func NewWithdrawRoutingStatsService(
	withdrawRepo withdrawrepo.WithdrawRepository,
	settlementRepo settlementcompany.SettlementCompanyRepository,
) *WithdrawRoutingStatsService {
	return &WithdrawRoutingStatsService{
		withdrawRepo:   withdrawRepo,
		settlementRepo: settlementRepo,
	}
}

// GetRoutingStats 获取提现分流统计
// Requirements: 14.1, 14.2
// Property 14: 提现分流统计准确性
func (s *WithdrawRoutingStatsService) GetRoutingStats(ctx context.Context, req *model.WithdrawRoutingStatsRequest) (*model.WithdrawRoutingStatsResponse, error) {
	stats, err := s.withdrawRepo.GetRoutingStats(ctx, req.DateFrom, req.DateTo)
	if err != nil {
		return nil, apierr.InternalError("failed to get routing stats").WithDetails(err.Error())
	}
	return stats, nil
}

// ListWithdrawalsByCompany 按结算公司查询提现列表
// Requirements: 14.5
func (s *WithdrawRoutingStatsService) ListWithdrawalsByCompany(ctx context.Context, req *model.ListWithdrawsByCompanyRequest) (*model.ListWithdrawsByCompanyResponse, error) {
	opts := withdrawrepo.WithdrawByCompanyOptions{
		SettlementCompanyID: req.SettlementCompanyID,
		Status:              req.Status,
		DateFrom:            req.DateFrom,
		DateTo:              req.DateTo,
		Page:                req.Page,
		PageSize:            req.PageSize,
	}

	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 || opts.PageSize > 100 {
		opts.PageSize = 20
	}

	withdraws, total, err := s.withdrawRepo.ListByCompany(ctx, opts)
	if err != nil {
		return nil, apierr.InternalError("failed to list withdrawals by company").WithDetails(err.Error())
	}

	return &model.ListWithdrawsByCompanyResponse{
		Total:     total,
		Page:      opts.Page,
		PageSize:  opts.PageSize,
		Withdraws: withdraws,
	}, nil
}

// GenerateRoutingReport 生成提现分流报表
// Requirements: 14.3, 14.4
func (s *WithdrawRoutingStatsService) GenerateRoutingReport(ctx context.Context, req *model.WithdrawRoutingReportRequest) (*model.WithdrawRoutingReport, error) {
	var dateFrom, dateTo time.Time

	switch req.ReportType {
	case "monthly":
		if req.Month < 1 || req.Month > 12 {
			return nil, apierr.BadRequest("invalid month for monthly report")
		}
		dateFrom = time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
		dateTo = dateFrom.AddDate(0, 1, 0)
	case "quarterly":
		if req.Quarter < 1 || req.Quarter > 4 {
			return nil, apierr.BadRequest("invalid quarter for quarterly report")
		}
		startMonth := (req.Quarter-1)*3 + 1
		dateFrom = time.Date(req.Year, time.Month(startMonth), 1, 0, 0, 0, 0, time.UTC)
		dateTo = dateFrom.AddDate(0, 3, 0)
	case "yearly":
		dateFrom = time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
		dateTo = dateFrom.AddDate(1, 0, 0)
	default:
		return nil, apierr.BadRequest("invalid report type")
	}

	stats, err := s.withdrawRepo.GetRoutingStats(ctx, &dateFrom, &dateTo)
	if err != nil {
		return nil, apierr.InternalError("failed to get routing stats for report").WithDetails(err.Error())
	}

	report := &model.WithdrawRoutingReport{
		ReportType:             req.ReportType,
		Year:                   req.Year,
		Month:                  req.Month,
		Quarter:                req.Quarter,
		TotalWithdrawals:       stats.TotalWithdrawals,
		TotalAmountCents:       stats.TotalAmountCents,
		TotalTaxDeductedCents:  stats.TotalTaxDeductedCents,
		TotalActualAmountCents: stats.TotalActualAmountCents,
		ByCompany:              stats.ByCompany,
		GeneratedAt:            time.Now(),
	}

	return report, nil
}

// GetCompanyWithdrawalStats 获取单个结算公司的提现统计
// Requirements: 14.5
func (s *WithdrawRoutingStatsService) GetCompanyWithdrawalStats(ctx context.Context, companyID uint64, dateFrom, dateTo *time.Time) (*model.WithdrawRoutingStats, error) {
	// 验证结算公司存在
	_, err := s.settlementRepo.Get(ctx, companyID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apierr.NotFound("settlement company not found")
		}
		return nil, apierr.InternalError("failed to get settlement company").WithDetails(err.Error())
	}

	stats, err := s.withdrawRepo.GetRoutingStatsByCompany(ctx, companyID, dateFrom, dateTo)
	if err != nil {
		return nil, apierr.InternalError("failed to get company withdrawal stats").WithDetails(err.Error())
	}

	return stats, nil
}
