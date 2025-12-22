package admin

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service/sensitiveword"
	"gamelink/pkg/apierr"
)

// SensitiveWordHandler 敏感词管理接口
type SensitiveWordHandler struct {
	svc *sensitiveword.SensitiveWordService
}

// NewSensitiveWordHandler 创建敏感词处理器
func NewSensitiveWordHandler(svc *sensitiveword.SensitiveWordService) *SensitiveWordHandler {
	return &SensitiveWordHandler{svc: svc}
}

// ListSensitiveWords
// @Summary      列出敏感词
// @Tags         Admin/SensitiveWords
// @Security     BearerAuth
// @Produce      json
// @Param        page       query  int     false  "页码"
// @Param        pageSize   query  int     false  "每页数量"
// @Param        keyword    query  string  false  "关键词搜索"
// @Param        category   query  string  false  "分类" Enums(political,pornographic,violent,advertising,other)
// @Param        severity   query  string  false  "严重程度" Enums(low,medium,high)
// @Success      200  {object}  model.APIResponse[sensitiveword.ListSensitiveWordsResponse]
// @Router       /admin/sensitive-words [get]
func (h *SensitiveWordHandler) ListSensitiveWords(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	var category *model.SensitiveWordCategory
	if categoryStr := strings.TrimSpace(c.Query("category")); categoryStr != "" {
		cat := model.SensitiveWordCategory(categoryStr)
		if !cat.Valid() {
			respondBadRequest(c, "invalid category")
			return
		}
		category = &cat
	}

	var severity *model.SensitiveWordSeverity
	if severityStr := strings.TrimSpace(c.Query("severity")); severityStr != "" {
		sev := model.SensitiveWordSeverity(severityStr)
		if !sev.Valid() {
			respondBadRequest(c, "invalid severity")
			return
		}
		severity = &sev
	}

	req := sensitiveword.ListSensitiveWordsRequest{
		Page:     page,
		PageSize: pageSize,
		Keyword:  strings.TrimSpace(c.Query("keyword")),
		Category: category,
		Severity: severity,
	}

	resp, err := h.svc.ListSensitiveWords(c.Request.Context(), req)
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondSuccess(c, resp)
}

// AddSensitiveWord
// @Summary      添加敏感词
// @Tags         Admin/SensitiveWords
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  sensitiveword.AddSensitiveWordRequest  true  "敏感词信息"
// @Success      201  {object}  model.APIResponse[sensitiveword.SensitiveWordDTO]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/sensitive-words [post]
func (h *SensitiveWordHandler) AddSensitiveWord(c *gin.Context) {
	var req sensitiveword.AddSensitiveWordRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	dto, err := h.svc.AddSensitiveWord(c.Request.Context(), req)
	if errors.Is(err, sensitiveword.ErrValidation) {
		respondBadRequest(c, "validation failed")
		return
	}
	if errors.Is(err, sensitiveword.ErrDuplicate) {
		respondBadRequest(c, "sensitive word already exists")
		return
	}
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondCreated(c, dto)
}

// UpdateSensitiveWord
// @Summary      更新敏感词
// @Tags         Admin/SensitiveWords
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                                      true  "敏感词ID"
// @Param        request  body  sensitiveword.UpdateSensitiveWordRequest  true  "敏感词信息"
// @Success      200  {object}  model.APIResponse[any]
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/sensitive-words/{id} [put]
func (h *SensitiveWordHandler) UpdateSensitiveWord(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req sensitiveword.UpdateSensitiveWordRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	err := h.svc.UpdateSensitiveWord(c.Request.Context(), id, req)
	if errors.Is(err, sensitiveword.ErrValidation) {
		respondBadRequest(c, "validation failed")
		return
	}
	if errors.Is(err, sensitiveword.ErrDuplicate) {
		respondBadRequest(c, "sensitive word already exists")
		return
	}
	if errors.Is(err, sensitiveword.ErrNotFound) {
		respondError(c, apierr.NotFound("sensitive word not found"))
		return
	}
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondMsg(c, "updated")
}

// DeleteSensitiveWord
// @Summary      删除敏感词
// @Tags         Admin/SensitiveWords
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "敏感词ID"
// @Success      200  {object}  model.APIResponse[any]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/sensitive-words/{id} [delete]
func (h *SensitiveWordHandler) DeleteSensitiveWord(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	err := h.svc.DeleteSensitiveWord(c.Request.Context(), id)
	if errors.Is(err, sensitiveword.ErrNotFound) {
		respondError(c, apierr.NotFound("sensitive word not found"))
		return
	}
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondDeleted(c)
}

// DetectSensitiveWords
// @Summary      检测敏感词
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  sensitiveword.DetectSensitiveWordsRequest  true  "检测内容"
// @Success      200  {object}  model.APIResponse[sensitiveword.DetectSensitiveWordsResponse]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/reviews/detect-sensitive [post]
func (h *SensitiveWordHandler) DetectSensitiveWords(c *gin.Context) {
	var req sensitiveword.DetectSensitiveWordsRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	resp, err := h.svc.DetectSensitiveWords(c.Request.Context(), req)
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondSuccess(c, resp)
}
