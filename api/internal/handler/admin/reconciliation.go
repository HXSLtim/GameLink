package admin

import (
	"errors"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	reconservice "gamelink/internal/service/reconciliation"
	"gamelink/pkg/apierr"
)

// ReconciliationHandler handles admin reconciliation APIs.
type ReconciliationHandler struct {
	svc *reconservice.Service
}

// NewReconciliationHandler creates a reconciliation handler.
func NewReconciliationHandler(svc *reconservice.Service) *ReconciliationHandler {
	return &ReconciliationHandler{svc: svc}
}

// ReconciliationCreateRequest is the create payload.
type ReconciliationCreateRequest struct {
	ReconciliationNo   string                            `json:"reconciliationNo"`
	ReconciliationDate time.Time                         `json:"reconciliationDate" binding:"required"`
	Type               model.ReconciliationType          `json:"type" binding:"required"`
	PeriodStart        time.Time                         `json:"periodStart" binding:"required"`
	PeriodEnd          time.Time                         `json:"periodEnd" binding:"required"`
	Abstract           string                            `json:"abstract"`
	Details            []ReconciliationCreateDetailInput `json:"details"`
}

// ReconciliationCreateDetailInput is one create detail line.
type ReconciliationCreateDetailInput struct {
	ExternalType   string    `json:"externalType" binding:"required"`
	ExternalNo     string    `json:"externalNo" binding:"required"`
	ExternalAmount int64     `json:"externalAmount"`
	ExternalDate   time.Time `json:"externalDate" binding:"required"`
	InternalType   string    `json:"internalType" binding:"required"`
	InternalNo     string    `json:"internalNo" binding:"required"`
	InternalAmount int64     `json:"internalAmount"`
	InternalDate   time.Time `json:"internalDate" binding:"required"`
	Remark         string    `json:"remark"`
}

// ReconciliationExecuteRequest is the execute payload.
type ReconciliationExecuteRequest struct {
	Status *model.ReconciliationStatus `json:"status"`
}

// List lists reconciliations.
func (h *ReconciliationHandler) List(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	opts := reconservice.ListInput{
		Page:     page,
		PageSize: pageSize,
	}

	if v := strings.TrimSpace(c.Query("type")); v != "" {
		tp := model.ReconciliationType(v)
		opts.Type = &tp
	}
	if v := strings.TrimSpace(c.Query("status")); v != "" {
		status := model.ReconciliationStatus(v)
		opts.Status = &status
	}

	dateFrom, err := queryTimePtr(c, "date_from")
	if err != nil {
		respondBadRequest(c, apierr.ErrInvalidDateFrom)
		return
	}
	dateTo, err := queryTimePtr(c, "date_to")
	if err != nil {
		respondBadRequest(c, apierr.ErrInvalidDateTo)
		return
	}
	opts.DateFrom = dateFrom
	opts.DateTo = dateTo

	items, pagination, listErr := h.svc.List(c.Request.Context(), opts)
	if listErr != nil {
		respondError(c, listErr)
		return
	}

	respondList(c, items, pagination)
}

// Get gets reconciliation detail with lines.
func (h *ReconciliationHandler) Get(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	rec, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, rec)
}

// Create creates a reconciliation record.
func (h *ReconciliationHandler) Create(c *gin.Context) {
	var req ReconciliationCreateRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	details := make([]reconservice.CreateDetailInput, 0, len(req.Details))
	for _, d := range req.Details {
		details = append(details, reconservice.CreateDetailInput{
			ExternalType:   d.ExternalType,
			ExternalNo:     d.ExternalNo,
			ExternalAmount: d.ExternalAmount,
			ExternalDate:   d.ExternalDate,
			InternalType:   d.InternalType,
			InternalNo:     d.InternalNo,
			InternalAmount: d.InternalAmount,
			InternalDate:   d.InternalDate,
			Remark:         d.Remark,
		})
	}

	created, err := h.svc.Create(c.Request.Context(), reconservice.CreateInput{
		ReconciliationNo:   req.ReconciliationNo,
		ReconciliationDate: req.ReconciliationDate,
		Type:               req.Type,
		PeriodStart:        req.PeriodStart,
		PeriodEnd:          req.PeriodEnd,
		Abstract:           req.Abstract,
		Details:            details,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, created)
}

// Execute executes reconciliation and transitions status.
func (h *ReconciliationHandler) Execute(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	adminUserID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req ReconciliationExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		respondError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return
	}

	out, err := h.svc.Execute(c.Request.Context(), id, adminUserID, reconservice.ExecuteInput{
		TargetStatus: req.Status,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, out)
}

// RegisterReconciliationRoutes registers reconciliation admin routes.
func RegisterReconciliationRoutes(rg *gin.RouterGroup, svc *reconservice.Service, pm *middleware.PermissionMiddleware) {
	h := NewReconciliationHandler(svc)

	group := rg.Group("/reconciliations")
	group.Use(pm.RequireAuth())
	{
		group.GET("", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/reconciliations"), h.List)
		group.GET("/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/reconciliations/:id"), h.Get)
		group.POST("", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/reconciliations"), h.Create)
		group.POST("/:id/execute", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/reconciliations/:id/execute"), h.Execute)
	}
}
