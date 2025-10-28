package admin

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	apierr "gamelink/internal/handler"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service"
)

// ReviewHandler 管理评价接口。
type ReviewHandler struct{ svc *service.AdminService }

func NewReviewHandler(s *service.AdminService) *ReviewHandler { return &ReviewHandler{svc: s} }

// ListReviews
// @Summary      评价列表
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Produce      json
// @Param        page       query  int  false  "页码"
// @Param        page_size  query  int  false  "每页数量"
// @Param        order_id   query  int  false  "订单ID"
// @Param        user_id    query  int  false  "用户ID"
// @Param        player_id  query  int  false  "陪玩师ID"
// @Param        date_from  query  string false "开始时间"
// @Param        date_to    query  string false "结束时间"
// @Success      200  {object}  map[string]any
// @Router       /admin/reviews [get]
func (h *ReviewHandler) ListReviews(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}
	var orderID, userID, playerID *uint64
	if v, err := queryUint64Ptr(c, "order_id"); err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidOrderID)
		return
	} else {
		orderID = v
	}
	if v, err := queryUint64Ptr(c, "user_id"); err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidUserID)
		return
	} else {
		userID = v
	}
	if v, err := queryUint64Ptr(c, "player_id"); err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidPlayerID)
		return
	} else {
		playerID = v
	}
	dateFrom, err := queryTimePtr(c, "date_from")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidDateFrom)
		return
	}
	dateTo, err := queryTimePtr(c, "date_to")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidDateTo)
		return
	}
	items, p, err := h.svc.ListReviews(c.Request.Context(), repository.ReviewListOptions{Page: page, PageSize: pageSize, OrderID: orderID, UserID: userID, PlayerID: playerID, DateFrom: dateFrom, DateTo: dateTo})
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}
	items = ensureSlice(items)
	writeJSON(c, 200, model.APIResponse[[]model.Review]{Success: true, Code: 200, Message: "OK", Data: items, Pagination: p})
}

// GetReview
// @Summary      获取评价
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "评价ID"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/reviews/{id} [get]
func (h *ReviewHandler) GetReview(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	}
	item, err := h.svc.GetReview(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}
	writeJSON(c, 200, model.APIResponse[*model.Review]{Success: true, Code: 200, Message: "OK", Data: item})
}

// CreateReview
// @Summary      创建评价
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  CreateReviewPayload  true  "评价"
// @Success      201  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Router       /admin/reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var p CreateReviewPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidJSONPayload)
		return
	}
	r := model.Review{OrderID: p.OrderID, UserID: p.UserID, PlayerID: p.PlayerID, Score: model.Rating(p.Score), Content: strings.TrimSpace(p.Content)}
	out, err := h.svc.CreateReview(c.Request.Context(), r)
	if errors.Is(err, service.ErrValidation) {
		_ = c.Error(service.ErrValidation)
		return
	}
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}
	writeJSON(c, 201, model.APIResponse[*model.Review]{Success: true, Code: 201, Message: "created", Data: out})
}

// UpdateReview
// @Summary      更新评价
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                    true  "评价ID"
// @Param        request  body  UpdateReviewPayload    true  "评价"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/reviews/{id} [put]
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	}
	var p UpdateReviewPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidJSONPayload)
		return
	}
	out, err := h.svc.UpdateReview(c.Request.Context(), id, model.Rating(p.Score), p.Content)
	if errors.Is(err, service.ErrValidation) {
		_ = c.Error(service.ErrValidation)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}
	writeJSON(c, 200, model.APIResponse[*model.Review]{Success: true, Code: 200, Message: "updated", Data: out})
}

// DeleteReview
// @Summary      删除评价
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "评价ID"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/reviews/{id} [delete]
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	}
	if err := h.svc.DeleteReview(c.Request.Context(), id); errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	} else if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}
	writeJSON(c, 200, model.APIResponse[any]{Success: true, Code: 200, Message: "deleted"})
}

// ListReviewLogs
// @Summary      获取评价操作日志
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Produce      json
// @Param        id           path   int  true  "评价ID"
// @Param        page         query  int  false "页码"
// @Param        page_size    query  int  false "每页数量"
// @Param        action       query  string false "动作过滤" Enums(create,update,delete)
// @Param        actor_user_id query int   false "操作者用户ID"
// @Param        date_from    query  string false "开始时间"
// @Param        date_to      query  string false "结束时间"
// @Param        export       query  string false "导出格式" Enums(csv)
// @Param        fields       query  string false "导出列（逗号分隔）"
// @Param        header_lang  query  string false "列头语言" Enums(en,zh)
// @Success      200  {object}  map[string]any
// @Router       /admin/reviews/{id}/logs [get]
func (h *ReviewHandler) ListReviewLogs(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	}
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}
	var actorID *uint64
	if v, err := queryUint64Ptr(c, "actor_user_id"); err == nil {
		actorID = v
	}
	var dateFrom, dateTo *time.Time
	if v, err := queryTimePtr(c, "date_from"); err == nil {
		dateFrom = v
	} else if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidDateFrom)
		return
	}
	if v, err := queryTimePtr(c, "date_to"); err == nil {
		dateTo = v
	} else if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidDateTo)
		return
	}
	opts := repository.OperationLogListOptions{Page: page, PageSize: pageSize, Action: strings.TrimSpace(c.Query("action")), ActorUserID: actorID, DateFrom: dateFrom, DateTo: dateTo}
	items, p, err := h.svc.ListOperationLogs(c.Request.Context(), "review", id, opts)
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}
	if strings.EqualFold(strings.TrimSpace(c.Query("export")), "csv") {
		exportOperationLogsCSV(c, "review", id, items)
		return
	}
	items = ensureSlice(items)
	writeJSON(c, 200, model.APIResponse[[]model.OperationLog]{Success: true, Code: 200, Message: "OK", Data: items, Pagination: p})
}

// ListPlayerReviews
// @Summary      获取陪玩师的评价
// @Tags         Admin/Players
// @Security     BearerAuth
// @Produce      json
// @Param        id         path   int  true  "陪玩师ID"
// @Param        page       query  int  false  "页码"
// @Param        page_size  query  int  false  "每页数量"
// @Success      200  {object}  map[string]any
// @Router       /admin/players/{id}/reviews [get]
func (h *ReviewHandler) ListPlayerReviews(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	}
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}
	pid := id
	items, p, err := h.svc.ListReviews(c.Request.Context(), repository.ReviewListOptions{Page: page, PageSize: pageSize, PlayerID: &pid})
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}
	items = ensureSlice(items)
	writeJSON(c, 200, model.APIResponse[[]model.Review]{Success: true, Code: 200, Message: "OK", Data: items, Pagination: p})
}

type CreateReviewPayload struct {
	OrderID  uint64 `json:"order_id" binding:"required"`
	UserID   uint64 `json:"user_id" binding:"required"`
	PlayerID uint64 `json:"player_id" binding:"required"`
	Score    uint8  `json:"score" binding:"required"`
	Content  string `json:"content"`
}

type UpdateReviewPayload struct {
	Score   uint8  `json:"score" binding:"required"`
	Content string `json:"content"`
}
