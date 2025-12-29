package admin

import (
	"github.com/gin-gonic/gin"

	"gamelink/pkg/apierr"
)

// BatchUpdateGameStatusRequest 批量更新游戏状态请求
type BatchUpdateGameStatusRequest struct {
	GameIDs []uint64 `json:"gameIds" binding:"required,min=1,max=100"`
	IsActive bool    `json:"isActive" binding:"required"`
}

// BatchUpdateGameCategoryRequest 批量更新游戏分类请求
type BatchUpdateGameCategoryRequest struct {
	GameIDs  []uint64 `json:"gameIds" binding:"required,min=1,max=100"`
	Category string   `json:"category" binding:"required"`
}

// BatchUpdateGamesStatus 批量更新游戏状态
// @Summary      批量更新游戏状态
// @Description  批量启用/禁用游戏
// @Tags         Admin/Games
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchUpdateGameStatusRequest  true  "游戏ID列表和状态"
// @Success      200  {object}  adminservice.BatchOperationResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/games/batch/status [post]
func (h *GameHandler) BatchUpdateGamesStatus(c *gin.Context) {
	var req BatchUpdateGameStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.GameIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("game_ids is required"))
		return
	}
	if len(req.GameIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 games per batch"))
		return
	}

	result, err := h.svc.BatchUpdateGamesStatusWithResponse(c.Request.Context(), req.GameIDs, req.IsActive)
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch update games status failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}

// BatchUpdateGamesCategory 批量更新游戏分类
// @Summary      批量更新游戏分类
// @Description  批量更改游戏分类
// @Tags         Admin/Games
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchUpdateGameCategoryRequest  true  "游戏ID列表和分类"
// @Success      200  {object}  adminservice.BatchOperationResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/games/batch/category [post]
func (h *GameHandler) BatchUpdateGamesCategory(c *gin.Context) {
	var req BatchUpdateGameCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.GameIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("game_ids is required"))
		return
	}
	if len(req.GameIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 games per batch"))
		return
	}
	if req.Category == "" {
		respondAPIError(c, apierr.BadRequest("category is required"))
		return
	}

	result, err := h.svc.BatchUpdateGamesCategory(c.Request.Context(), req.GameIDs, req.Category)
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch update games category failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}
