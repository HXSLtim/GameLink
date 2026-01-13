package admin

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	reviewservice "gamelink/internal/service/review"
	"gamelink/pkg/apierr"
)

// ReviewSettingsHandler 评价展示设置处理器
type ReviewSettingsHandler struct {
	svc *reviewservice.SettingsService
}

// NewReviewSettingsHandler 创建评价展示设置处理器实例
func NewReviewSettingsHandler(svc *reviewservice.SettingsService) *ReviewSettingsHandler {
	return &ReviewSettingsHandler{svc: svc}
}

// GetReviewSettings
// @Summary      获取评价展示设置
// @Description  获取当前的评价展示规则配置
// @Tags         Admin/ReviewSettings
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  model.ReviewDisplaySettings
// @Failure      500  {object}  model.ErrorResponse
// @Router       /admin/review-settings [get]
func (h *ReviewSettingsHandler) GetReviewSettings(c *gin.Context) {
	settings, err := h.svc.GetSettings(c.Request.Context())
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}
	respondSuccess(c, settings)
}

// UpdateReviewSettingsPayload 更新评价展示设置请求体
type UpdateReviewSettingsPayload struct {
	// 排序方式：time/score/likes
	// @Enum time, score, likes
	SortBy *string `json:"sortBy"`
	// 最低评分阈值（1-5）
	MinScore *int `json:"minScore"`
	// 是否显示匿名评价
	ShowAnonymous *bool `json:"showAnonymous"`
	// 每页显示数量（1-100）
	PageSize *int `json:"pageSize"`
	// 是否自动批准评价
	AutoApprove *bool `json:"autoApprove"`
	// 自动批准最低评分（1-5）
	AutoApproveMinRating *int `json:"autoApproveMinRating"`
}

// UpdateReviewSettings
// @Summary      更新评价展示设置
// @Description  更新评价展示规则配置，支持部分更新
// @Tags         Admin/ReviewSettings
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  UpdateReviewSettingsPayload  true  "设置信息"
// @Success      200  {object}  model.ReviewDisplaySettings
// @Failure      400  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /admin/review-settings [put]
func (h *ReviewSettingsHandler) UpdateReviewSettings(c *gin.Context) {
	var p UpdateReviewSettingsPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		respondBadRequest(c, apierr.ErrInvalidJSONPayload)
		return
	}

	// 构建更新输入
	input := reviewservice.UpdateSettingsInput{}

	if p.SortBy != nil {
		sortBy := model.ReviewSortBy(*p.SortBy)
		if !sortBy.Valid() {
			respondBadRequest(c, "invalid sortBy value, must be one of: time, score, likes")
			return
		}
		input.SortBy = &sortBy
	}

	if p.MinScore != nil {
		input.MinScore = p.MinScore
	}

	if p.ShowAnonymous != nil {
		input.ShowAnonymous = p.ShowAnonymous
	}

	if p.PageSize != nil {
		input.PageSize = p.PageSize
	}

	if p.AutoApprove != nil {
		input.AutoApprove = p.AutoApprove
	}

	if p.AutoApproveMinRating != nil {
		input.AutoApproveMinRating = p.AutoApproveMinRating
	}

	settings, err := h.svc.UpdateSettingsPartial(c.Request.Context(), input)
	if err != nil {
		if _, ok := err.(*model.ValidationError); ok {
			respondBadRequest(c, err.Error())
			return
		}
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondUpdated(c, settings)
}
