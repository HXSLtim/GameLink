package admin

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	apierr "gamelink/internal/handler"
	"gamelink/internal/logging"
	"gamelink/internal/model"
	assignmentservice "gamelink/internal/service/assignment"
)

// AssignmentHandler 管理客服指派与争议调解接口。
type AssignmentHandler struct {
	svc *assignmentservice.Service
}

// NewAssignmentHandler 创建实例。
func NewAssignmentHandler(svc *assignmentservice.Service) *AssignmentHandler {
	return &AssignmentHandler{svc: svc}
}

// ListPendingAssignments 列出待指派订单。
// @Summary      待指派订单列表
// @Tags         Admin/Assignments
// @Security     BearerAuth
// @Param        page       query  int false "页码"
// @Param        page_size  query  int false "每页数量"
// @Produce      json
// @Success      200 {object} model.APIResponse[[]assignmentservice.PendingAssignment]
// @Router       /admin/orders/pending-assign [get]
func (h *AssignmentHandler) ListPendingAssignments(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}
	items, total, err := h.svc.ListPendingAssignments(c.Request.Context(), page, pageSize)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	pagination := &model.Pagination{Page: page, PageSize: pageSize, Total: int(total)}
	if pageSize > 0 {
		pagination.TotalPages = (pagination.Total + pageSize - 1) / pageSize
		pagination.HasNext = page < pagination.TotalPages
		pagination.HasPrev = page > 1
	}
	traceID := assignmentservice.TraceIDFromContext(c.Request.Context())
	writeJSON(c, http.StatusOK, model.APIResponse[[]assignmentservice.PendingAssignment]{
		Success:    true,
		Code:       http.StatusOK,
		Message:    "OK",
		Data:       items,
		Pagination: pagination,
		TraceID:    traceID,
	})
}

// ListCandidates 获取指派候选人。
// @Summary      获取订单候选陪玩师
// @Tags         Admin/Assignments
// @Security     BearerAuth
// @Param        id     path   int  true  "订单ID"
// @Param        limit  query  int  false "返回数量"
// @Produce      json
// @Success      200 {object} model.APIResponse[[]assignmentservice.Candidate]
// @Router       /admin/orders/{id}/candidates [get]
func (h *AssignmentHandler) ListCandidates(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	limit, err := queryIntDefault(c, "limit", 10)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidPageSize)
		return
	}
	candidates, err := h.svc.ListCandidates(c.Request.Context(), id, limit)
	if err != nil {
		switch {
		case errors.Is(err, assignmentservice.ErrNotFound):
			_ = c.Error(assignmentservice.ErrNotFound)
		case errors.Is(err, assignmentservice.ErrValidation):
			writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidParameter)
		default:
			writeJSONError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	traceID := assignmentservice.TraceIDFromContext(c.Request.Context())
	writeJSON(c, http.StatusOK, model.APIResponse[[]assignmentservice.Candidate]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    candidates,
		TraceID: traceID,
	})
}

// AssignRequest 指派请求体。
type AssignRequest struct {
	PlayerID uint64 `json:"player_id" binding:"required"`
	Source   string `json:"source"`
}

// Assign 指派订单。
// @Summary      指派陪玩师
// @Tags         Admin/Assignments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int            true  "订单ID"
// @Param        request  body  AssignRequest  true  "指派信息"
// @Success      200 {object} model.APIResponse[*model.Order]
// @Router       /admin/orders/{id}/assign [post]
func (h *AssignmentHandler) Assign(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	var req AssignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}
	traceID := assignmentservice.TraceIDFromContext(c.Request.Context())
	actorID, ok := logging.ActorUserIDFromContext(c.Request.Context())
	var actorPtr *uint64
	if ok {
		actorPtr = &actorID
	}
	input := assignmentservice.AssignInput{
		PlayerID:    req.PlayerID,
		Source:      model.OrderAssignmentSource(strings.ToLower(strings.TrimSpace(req.Source))),
		ActorUserID: actorPtr,
		TraceID:     traceID,
	}
	order, err := h.svc.Assign(c.Request.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, assignmentservice.ErrValidation):
			writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidParameter)
			return
		case errors.Is(err, assignmentservice.ErrNotFound):
			_ = c.Error(assignmentservice.ErrNotFound)
			return
		default:
			writeJSONError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	writeJSON(c, http.StatusOK, model.APIResponse[*model.Order]{
		Success: true,
		Code:    http.StatusOK,
		Message: "updated",
		Data:    order,
		TraceID: traceID,
	})
}

// CancelAssignRequest 回退指派请求体。
type CancelAssignRequest struct {
	Reason string `json:"reason"`
}

// CancelAssignment 回退订单指派。
// @Summary      回退指派
// @Tags         Admin/Assignments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                  true  "订单ID"
// @Param        request  body  CancelAssignRequest  true  "回退原因"
// @Success      200 {object} model.APIResponse[*model.Order]
// @Router       /admin/orders/{id}/assign/cancel [post]
func (h *AssignmentHandler) CancelAssignment(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	var req CancelAssignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}
	traceID := assignmentservice.TraceIDFromContext(c.Request.Context())
	actorID, ok := logging.ActorUserIDFromContext(c.Request.Context())
	var actorPtr *uint64
	if ok {
		actorPtr = &actorID
	}
	order, err := h.svc.CancelAssignment(c.Request.Context(), id, assignmentservice.CancelAssignInput{
		Reason:      strings.TrimSpace(req.Reason),
		ActorUserID: actorPtr,
		TraceID:     traceID,
	})
	if err != nil {
		switch {
		case errors.Is(err, assignmentservice.ErrNotFound):
			_ = c.Error(assignmentservice.ErrNotFound)
			return
		default:
			writeJSONError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	writeJSON(c, http.StatusOK, model.APIResponse[*model.Order]{
		Success: true,
		Code:    http.StatusOK,
		Message: "updated",
		Data:    order,
		TraceID: traceID,
	})
}

// MediateRequest 调解请求体。
type MediateRequest struct {
	Resolution        string  `json:"resolution" binding:"required"`
	Note              string  `json:"note"`
	RefundAmountCents *int64  `json:"refund_amount_cents"`
	ReassignPlayerID  *uint64 `json:"reassign_player_id"`
}

// MediateDispute 调解争议。
// @Summary      调解订单争议
// @Tags         Admin/Assignments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int             true  "订单ID"
// @Param        request  body  MediateRequest  true  "调解信息"
// @Success      200 {object} model.APIResponse[*model.OrderDispute]
// @Router       /admin/orders/{id}/mediate [post]
func (h *AssignmentHandler) MediateDispute(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	var req MediateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}
	traceID := assignmentservice.TraceIDFromContext(c.Request.Context())
	actorID, ok := logging.ActorUserIDFromContext(c.Request.Context())
	var actorPtr *uint64
	if ok {
		actorPtr = &actorID
	}
	input := assignmentservice.MediateInput{
		Resolution:        model.OrderDisputeResolution(strings.ToLower(strings.TrimSpace(req.Resolution))),
		Note:              strings.TrimSpace(req.Note),
		RefundAmountCents: req.RefundAmountCents,
		ReassignPlayerID:  req.ReassignPlayerID,
		ActorUserID:       actorPtr,
		TraceID:           traceID,
	}
	dispute, err := h.svc.MediateDispute(c.Request.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, assignmentservice.ErrValidation):
			writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidParameter)
			return
		case errors.Is(err, assignmentservice.ErrNotFound):
			_ = c.Error(assignmentservice.ErrNotFound)
			return
		default:
			writeJSONError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	writeJSON(c, http.StatusOK, model.APIResponse[*model.OrderDispute]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    dispute,
		TraceID: traceID,
	})
}

// ListDisputes 获取订单争议列表。
// @Summary      查询订单争议
// @Tags         Admin/Assignments
// @Security     BearerAuth
// @Param        id  path int true "订单ID"
// @Produce      json
// @Success      200 {object} model.APIResponse[[]model.OrderDispute]
// @Router       /admin/orders/{id}/disputes [get]
func (h *AssignmentHandler) ListDisputes(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	disputes, err := h.svc.ListDisputes(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, assignmentservice.ErrNotFound):
			_ = c.Error(assignmentservice.ErrNotFound)
			return
		default:
			writeJSONError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	traceID := assignmentservice.TraceIDFromContext(c.Request.Context())
	writeJSON(c, http.StatusOK, model.APIResponse[[]model.OrderDispute]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    ensureSlice(disputes),
		TraceID: traceID,
	})
}
