package reconciliation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
)

// Service provides reconciliation business logic.
type Service struct {
	repo repository.ReconciliationRepository
	now  func() time.Time
}

// NewService creates a reconciliation service.
func NewService(repo repository.ReconciliationRepository) *Service {
	return &Service{
		repo: repo,
		now:  time.Now,
	}
}

// ListInput describes list query parameters.
type ListInput struct {
	Page     int
	PageSize int
	Type     *model.ReconciliationType
	Status   *model.ReconciliationStatus
	DateFrom *time.Time
	DateTo   *time.Time
}

// CreateDetailInput describes one reconciliation line.
type CreateDetailInput struct {
	ExternalType   string
	ExternalNo     string
	ExternalAmount int64
	ExternalDate   time.Time

	InternalType   string
	InternalNo     string
	InternalAmount int64
	InternalDate   time.Time

	Remark string
}

// CreateInput describes reconciliation creation parameters.
type CreateInput struct {
	ReconciliationNo   string
	ReconciliationDate time.Time
	Type               model.ReconciliationType
	PeriodStart        time.Time
	PeriodEnd          time.Time
	Abstract           string
	Details            []CreateDetailInput
}

// ExecuteInput describes execute/transition input.
type ExecuteInput struct {
	TargetStatus *model.ReconciliationStatus
}

// List returns reconciliations with pagination.
func (s *Service) List(ctx context.Context, input ListInput) ([]model.Reconciliation, *model.Pagination, error) {
	opts := repository.ReconciliationListOptions{
		Page:     repository.NormalizePage(input.Page),
		PageSize: repository.NormalizePageSize(input.PageSize),
		Type:     input.Type,
		Status:   input.Status,
		DateFrom: input.DateFrom,
		DateTo:   input.DateTo,
	}
	items, total, err := s.repo.List(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	page, pageSize, totalInt, totalPages, hasNext, hasPrev := repository.BuildPagination(opts.Page, opts.PageSize, total)
	return items, &model.Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      totalInt,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}, nil
}

// Get returns one reconciliation with lines.
func (s *Service) Get(ctx context.Context, id uint64) (*model.Reconciliation, error) {
	return s.repo.Get(ctx, id, true)
}

// Create creates a reconciliation.
func (s *Service) Create(ctx context.Context, input CreateInput) (*model.Reconciliation, error) {
	if input.ReconciliationDate.IsZero() {
		return nil, apierr.BadRequest("reconciliationDate is required")
	}
	if input.Type == "" {
		return nil, apierr.BadRequest("type is required")
	}
	if input.PeriodStart.IsZero() || input.PeriodEnd.IsZero() {
		return nil, apierr.BadRequest("periodStart and periodEnd are required")
	}
	if input.PeriodStart.After(input.PeriodEnd) {
		return nil, apierr.BadRequest("periodStart cannot be after periodEnd")
	}

	recNo := input.ReconciliationNo
	if recNo == "" {
		recNo = s.generateReconciliationNo()
	}

	details := make([]model.ReconciliationDetail, 0, len(input.Details))
	totalRecords := len(input.Details)
	matchedRecords := 0
	var differenceAmount int64
	for i, d := range input.Details {
		diff := d.ExternalAmount - d.InternalAmount
		status := "mismatch"
		if diff == 0 {
			status = "matched"
			matchedRecords++
		}
		if diff < 0 {
			differenceAmount -= diff
		} else {
			differenceAmount += diff
		}

		details = append(details, model.ReconciliationDetail{
			LineNo:           i + 1,
			ExternalType:     d.ExternalType,
			ExternalNo:       d.ExternalNo,
			ExternalAmount:   d.ExternalAmount,
			ExternalDate:     d.ExternalDate,
			InternalType:     d.InternalType,
			InternalNo:       d.InternalNo,
			InternalAmount:   d.InternalAmount,
			InternalDate:     d.InternalDate,
			Status:           status,
			DifferenceAmount: diff,
			Remark:           d.Remark,
		})
	}

	rec := &model.Reconciliation{
		ReconciliationNo:   recNo,
		ReconciliationDate: input.ReconciliationDate,
		Type:               input.Type,
		Status:             model.ReconciliationStatusPending,
		PeriodStart:        input.PeriodStart,
		PeriodEnd:          input.PeriodEnd,
		TotalRecords:       totalRecords,
		MatchedRecords:     matchedRecords,
		DifferenceAmount:   differenceAmount,
		Abstract:           input.Abstract,
		Details:            details,
	}

	if err := s.repo.Create(ctx, rec); err != nil {
		return nil, err
	}

	return rec, nil
}

// Execute runs reconciliation and performs status transition.
func (s *Service) Execute(ctx context.Context, id uint64, adminUserID uint64, input ExecuteInput) (*model.Reconciliation, error) {
	if adminUserID == 0 {
		return nil, apierr.Unauthorized("missing admin user")
	}

	if input.TargetStatus != nil {
		switch *input.TargetStatus {
		case model.ReconciliationStatusSuccess, model.ReconciliationStatusFailed, model.ReconciliationStatusException:
		default:
			return nil, apierr.BadRequest("status must be success, failed, or exception")
		}
	}

	out, err := s.repo.Execute(ctx, id, repository.ReconciliationExecuteOptions{
		ProcessedBy: adminUserID,
		ForceStatus: input.TargetStatus,
	})
	if err != nil {
		if errors.Is(err, repository.ErrInvalidStatusTransition) {
			return nil, apierr.BadRequest("reconciliation status does not allow execute")
		}
		return nil, err
	}
	return out, nil
}

func (s *Service) generateReconciliationNo() string {
	now := s.now().UTC()
	return fmt.Sprintf("RCN-%s-%04d", now.Format("20060102150405"), now.UnixNano()%10000)
}
