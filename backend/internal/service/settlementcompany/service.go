package settlementcompany

import (
	"context"
	"errors"
	"fmt"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/settlementcompany"
	"gamelink/pkg/apierr"
	"reflect"
	"time"
)

var (
	// ErrNotFound 结算公司不存在
	ErrNotFound = apierr.NotFound("settlement company not found")
	// ErrValidation 表示输入校验失败
	ErrValidation = apierr.BadRequest("validation failed")
	// ErrCreditCodeExists 统一社会信用代码已存在
	ErrCreditCodeExists = apierr.Conflict("credit code already exists")
	// ErrInvalidCreditCode 统一社会信用代码格式无效
	ErrInvalidCreditCode = apierr.BadRequest("invalid credit code format")
	// ErrCompanyInactive 结算公司已禁用
	ErrCompanyInactive = apierr.BadRequest("settlement company is inactive")
	// ErrPlayerNotFound 陪玩师不存在
	ErrPlayerNotFound = apierr.NotFound("player not found")
	// ErrAssignmentNotFound 分配记录不存在
	ErrAssignmentNotFound = apierr.NotFound("assignment not found")
)

// SettlementCompanyService 结算公司服务
// Requirements: 11.1, 11.2, 11.3, 11.4, 11.5
type SettlementCompanyService struct {
	repo       settlementcompany.SettlementCompanyRepository
	playerRepo repository.PlayerRepository
}

// NewSettlementCompanyService 创建结算公司服务
func NewSettlementCompanyService(
	repo settlementcompany.SettlementCompanyRepository,
	playerRepo repository.PlayerRepository,
) *SettlementCompanyService {
	return &SettlementCompanyService{
		repo:       repo,
		playerRepo: playerRepo,
	}
}

// CreateCompany 创建结算公司
// Requirements: 11.1, 11.2
func (s *SettlementCompanyService) CreateCompany(ctx context.Context, req *model.CreateSettlementCompanyRequest, createdBy uint64) (*model.SettlementCompany, error) {
	// 验证统一社会信用代码格式
	// Property 11: 统一社会信用代码格式验证
	if !model.ValidateCreditCode(req.CreditCode) {
		return nil, ErrInvalidCreditCode.WithDetails("credit code must be 18 characters and match the standard format")
	}

	// 检查统一社会信用代码唯一性
	existing, err := s.repo.GetByCreditCode(ctx, req.CreditCode)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, apierr.InternalError("failed to check credit code uniqueness").WithDetails(err.Error())
	}
	if existing != nil {
		return nil, ErrCreditCodeExists.WithDetails(fmt.Sprintf("credit code %s already exists", req.CreditCode))
	}

	// 创建结算公司
	company := &model.SettlementCompany{
		Name:              req.Name,
		CreditCode:        req.CreditCode,
		TaxRegistrationNo: req.TaxRegistrationNo,
		BankName:          req.BankName,
		BankAccount:       req.BankAccount,
		BankBranch:        req.BankBranch,
		ContactName:       req.ContactName,
		ContactPhone:      req.ContactPhone,
		Address:           req.Address,
		Status:            model.CompanyStatusActive,
		PlayerCount:       0,
		CreatedBy:         createdBy,
	}

	if err := s.repo.Create(ctx, company); err != nil {
		return nil, apierr.InternalError("failed to create settlement company").WithDetails(err.Error())
	}

	return company, nil
}

// GetCompany 获取结算公司
// Requirements: 11.3
func (s *SettlementCompanyService) GetCompany(ctx context.Context, id uint64) (*model.SettlementCompany, error) {
	company, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("failed to get settlement company").WithDetails(err.Error())
	}
	return company, nil
}

// UpdateCompany 更新结算公司
// Requirements: 11.5
func (s *SettlementCompanyService) UpdateCompany(ctx context.Context, id uint64, req *model.UpdateSettlementCompanyRequest, updatedBy uint64) (*model.SettlementCompany, error) {
	company, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("failed to get settlement company").WithDetails(err.Error())
	}

	// 记录修改历史
	changes := s.detectChanges(company, req)

	// 应用更新
	if req.Name != nil {
		company.Name = *req.Name
	}
	if req.TaxRegistrationNo != nil {
		company.TaxRegistrationNo = *req.TaxRegistrationNo
	}
	if req.BankName != nil {
		company.BankName = *req.BankName
	}
	if req.BankAccount != nil {
		company.BankAccount = *req.BankAccount
	}
	if req.BankBranch != nil {
		company.BankBranch = *req.BankBranch
	}
	if req.ContactName != nil {
		company.ContactName = *req.ContactName
	}
	if req.ContactPhone != nil {
		company.ContactPhone = *req.ContactPhone
	}
	if req.Address != nil {
		company.Address = *req.Address
	}
	company.UpdatedBy = &updatedBy

	if err := s.repo.Update(ctx, company); err != nil {
		return nil, apierr.InternalError("failed to update settlement company").WithDetails(err.Error())
	}

	// 保存修改历史
	for _, change := range changes {
		change.SettlementCompanyID = id
		change.ChangedBy = updatedBy
		if err := s.repo.CreateHistory(ctx, &change); err != nil {
			// 记录历史失败不影响主流程，只记录日志
			continue
		}
	}

	return company, nil
}

// detectChanges 检测字段变更
func (s *SettlementCompanyService) detectChanges(company *model.SettlementCompany, req *model.UpdateSettlementCompanyRequest) []model.SettlementCompanyHistory {
	var changes []model.SettlementCompanyHistory

	if req.Name != nil && *req.Name != company.Name {
		changes = append(changes, model.SettlementCompanyHistory{
			FieldName: "name",
			OldValue:  company.Name,
			NewValue:  *req.Name,
		})
	}
	if req.TaxRegistrationNo != nil && *req.TaxRegistrationNo != company.TaxRegistrationNo {
		changes = append(changes, model.SettlementCompanyHistory{
			FieldName: "tax_registration_no",
			OldValue:  company.TaxRegistrationNo,
			NewValue:  *req.TaxRegistrationNo,
		})
	}
	if req.BankName != nil && *req.BankName != company.BankName {
		changes = append(changes, model.SettlementCompanyHistory{
			FieldName: "bank_name",
			OldValue:  company.BankName,
			NewValue:  *req.BankName,
		})
	}
	if req.BankAccount != nil && *req.BankAccount != company.BankAccount {
		changes = append(changes, model.SettlementCompanyHistory{
			FieldName: "bank_account",
			OldValue:  company.BankAccount,
			NewValue:  *req.BankAccount,
		})
	}
	if req.BankBranch != nil && *req.BankBranch != company.BankBranch {
		changes = append(changes, model.SettlementCompanyHistory{
			FieldName: "bank_branch",
			OldValue:  company.BankBranch,
			NewValue:  *req.BankBranch,
		})
	}
	if req.ContactName != nil && *req.ContactName != company.ContactName {
		changes = append(changes, model.SettlementCompanyHistory{
			FieldName: "contact_name",
			OldValue:  company.ContactName,
			NewValue:  *req.ContactName,
		})
	}
	if req.ContactPhone != nil && *req.ContactPhone != company.ContactPhone {
		changes = append(changes, model.SettlementCompanyHistory{
			FieldName: "contact_phone",
			OldValue:  company.ContactPhone,
			NewValue:  *req.ContactPhone,
		})
	}
	if req.Address != nil && *req.Address != company.Address {
		changes = append(changes, model.SettlementCompanyHistory{
			FieldName: "address",
			OldValue:  company.Address,
			NewValue:  *req.Address,
		})
	}

	return changes
}

// ListCompanies 查询结算公司列表
// Requirements: 11.3
func (s *SettlementCompanyService) ListCompanies(ctx context.Context, req *model.ListSettlementCompaniesRequest) (*model.ListSettlementCompaniesResponse, error) {
	opts := settlementcompany.ListOptions{
		Keyword:   req.Keyword,
		Page:      req.Page,
		PageSize:  req.PageSize,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	if req.Status != "" {
		opts.Status = &req.Status
	}

	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 || opts.PageSize > 100 {
		opts.PageSize = 20
	}

	companies, total, err := s.repo.List(ctx, opts)
	if err != nil {
		return nil, apierr.InternalError("failed to list settlement companies").WithDetails(err.Error())
	}

	return &model.ListSettlementCompaniesResponse{
		Total:     total,
		Page:      opts.Page,
		PageSize:  opts.PageSize,
		Companies: companies,
	}, nil
}

// ToggleCompanyStatus 切换结算公司状态
// Requirements: 11.4
func (s *SettlementCompanyService) ToggleCompanyStatus(ctx context.Context, id uint64, enabled bool, updatedBy uint64) error {
	company, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return apierr.InternalError("failed to get settlement company").WithDetails(err.Error())
	}

	var newStatus model.CompanyStatus
	if enabled {
		newStatus = model.CompanyStatusActive
	} else {
		newStatus = model.CompanyStatusInactive
	}

	// 记录状态变更历史
	if company.Status != newStatus {
		history := &model.SettlementCompanyHistory{
			SettlementCompanyID: id,
			FieldName:           "status",
			OldValue:            string(company.Status),
			NewValue:            string(newStatus),
			ChangedBy:           updatedBy,
		}
		if err := s.repo.CreateHistory(ctx, history); err != nil {
			// 记录历史失败不影响主流程
		}
	}

	if err := s.repo.ToggleStatus(ctx, id, newStatus); err != nil {
		return apierr.InternalError("failed to toggle settlement company status").WithDetails(err.Error())
	}

	return nil
}

// GetCompanyHistory 获取结算公司修改历史
// Requirements: 11.5
func (s *SettlementCompanyService) GetCompanyHistory(ctx context.Context, companyID uint64) ([]model.SettlementCompanyHistory, error) {
	// 先验证公司存在
	_, err := s.repo.Get(ctx, companyID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("failed to get settlement company").WithDetails(err.Error())
	}

	histories, err := s.repo.GetHistory(ctx, companyID)
	if err != nil {
		return nil, apierr.InternalError("failed to get settlement company history").WithDetails(err.Error())
	}

	return histories, nil
}

// Ensure reflect is used (for future use in generic change detection)
var _ = reflect.TypeOf

// AssignPlayerToCompany 分配陪玩师到结算公司
// Requirements: 12.1, 12.4
// Property 12: 陪玩师结算公司分配唯一性
func (s *SettlementCompanyService) AssignPlayerToCompany(ctx context.Context, req *model.AssignPlayerToCompanyRequest, assignedBy uint64) (*model.PlayerCompanyAssignment, error) {
	// 验证结算公司存在且处于活跃状态
	company, err := s.repo.Get(ctx, req.SettlementCompanyID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound.WithDetails("settlement company not found")
		}
		return nil, apierr.InternalError("failed to get settlement company").WithDetails(err.Error())
	}

	if !company.IsActive() {
		return nil, ErrCompanyInactive.WithDetails("cannot assign player to inactive company")
	}

	// 验证陪玩师存在
	if s.playerRepo != nil {
		_, err := s.playerRepo.Get(ctx, req.PlayerID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, ErrPlayerNotFound.WithDetails(fmt.Sprintf("player %d not found", req.PlayerID))
			}
			return nil, apierr.InternalError("failed to get player").WithDetails(err.Error())
		}
	}

	// 创建分配记录
	// Property 12: 陪玩师结算公司分配唯一性 - 仓储层会自动结束当前分配
	assignment := &model.PlayerCompanyAssignment{
		PlayerID:            req.PlayerID,
		SettlementCompanyID: req.SettlementCompanyID,
		EffectiveDate:       req.EffectiveDate,
		Reason:              req.Reason,
		AssignedBy:          assignedBy,
		IsCurrent:           true,
	}

	if err := s.repo.AssignPlayer(ctx, assignment); err != nil {
		return nil, apierr.InternalError("failed to assign player to company").WithDetails(err.Error())
	}

	return assignment, nil
}

// BatchAssignPlayers 批量分配陪玩师到结算公司
// Requirements: 12.3
func (s *SettlementCompanyService) BatchAssignPlayers(ctx context.Context, req *model.BatchAssignPlayersRequest, assignedBy uint64) (int, error) {
	// 验证结算公司存在且处于活跃状态
	company, err := s.repo.Get(ctx, req.SettlementCompanyID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return 0, ErrNotFound.WithDetails("settlement company not found")
		}
		return 0, apierr.InternalError("failed to get settlement company").WithDetails(err.Error())
	}

	if !company.IsActive() {
		return 0, ErrCompanyInactive.WithDetails("cannot assign players to inactive company")
	}

	// 验证所有陪玩师存在
	if s.playerRepo != nil {
		for _, playerID := range req.PlayerIDs {
			_, err := s.playerRepo.Get(ctx, playerID)
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					return 0, ErrPlayerNotFound.WithDetails(fmt.Sprintf("player %d not found", playerID))
				}
				return 0, apierr.InternalError("failed to get player").WithDetails(err.Error())
			}
		}
	}

	// 创建分配记录
	assignments := make([]model.PlayerCompanyAssignment, len(req.PlayerIDs))
	for i, playerID := range req.PlayerIDs {
		assignments[i] = model.PlayerCompanyAssignment{
			PlayerID:            playerID,
			SettlementCompanyID: req.SettlementCompanyID,
			EffectiveDate:       req.EffectiveDate,
			Reason:              req.Reason,
			AssignedBy:          assignedBy,
			IsCurrent:           true,
		}
	}

	if err := s.repo.BatchAssignPlayers(ctx, assignments); err != nil {
		return 0, apierr.InternalError("failed to batch assign players").WithDetails(err.Error())
	}

	return len(assignments), nil
}

// GetCurrentAssignment 获取陪玩师当前的结算公司分配
// Requirements: 12.5
func (s *SettlementCompanyService) GetCurrentAssignment(ctx context.Context, playerID uint64) (*model.PlayerCompanyAssignment, error) {
	assignment, err := s.repo.GetCurrentAssignment(ctx, playerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrAssignmentNotFound.WithDetails(fmt.Sprintf("no current assignment for player %d", playerID))
		}
		return nil, apierr.InternalError("failed to get current assignment").WithDetails(err.Error())
	}
	return assignment, nil
}

// GetAssignmentHistory 获取陪玩师的结算公司分配历史
// Requirements: 12.5
func (s *SettlementCompanyService) GetAssignmentHistory(ctx context.Context, playerID uint64) (*model.PlayerAssignmentHistoryResponse, error) {
	// 验证陪玩师存在
	if s.playerRepo != nil {
		_, err := s.playerRepo.Get(ctx, playerID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, ErrPlayerNotFound.WithDetails(fmt.Sprintf("player %d not found", playerID))
			}
			return nil, apierr.InternalError("failed to get player").WithDetails(err.Error())
		}
	}

	assignments, err := s.repo.GetAssignmentHistory(ctx, playerID)
	if err != nil {
		return nil, apierr.InternalError("failed to get assignment history").WithDetails(err.Error())
	}

	return &model.PlayerAssignmentHistoryResponse{
		Total:       int64(len(assignments)),
		Assignments: assignments,
	}, nil
}

// EndCurrentAssignment 结束陪玩师当前的结算公司分配
// Requirements: 12.4
func (s *SettlementCompanyService) EndCurrentAssignment(ctx context.Context, playerID uint64, endDate time.Time) error {
	// 检查是否有当前分配
	_, err := s.repo.GetCurrentAssignment(ctx, playerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrAssignmentNotFound.WithDetails(fmt.Sprintf("no current assignment for player %d", playerID))
		}
		return apierr.InternalError("failed to get current assignment").WithDetails(err.Error())
	}

	if err := s.repo.EndCurrentAssignment(ctx, playerID, endDate); err != nil {
		return apierr.InternalError("failed to end current assignment").WithDetails(err.Error())
	}

	return nil
}
