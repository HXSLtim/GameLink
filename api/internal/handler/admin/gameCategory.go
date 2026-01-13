package admin

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	adminservice "gamelink/internal/service/admin"
)

// GameCategory 游戏分类模型（类型别名）
type GameCategory = model.GameCategory

// GameCategoryHandler 处理游戏分类管理接口
type GameCategoryHandler struct {
	svc *adminservice.AdminService
}

// NewGameCategoryHandler 创建Handler
func NewGameCategoryHandler(svc *adminservice.AdminService) *GameCategoryHandler {
	return &GameCategoryHandler{svc: svc}
}

// CreateCategory
// @Summary      创建游戏分类
// @Tags         Admin/GameCategories
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  CreateCategoryRequest  true  "分类信息"
// @Success      201  {object}  GameCategory
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/game-categories [post]
func (h *GameCategoryHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	category, err := h.svc.CreateGameCategory(c.Request.Context(), adminservice.CreateGameCategoryInput{
		Name:        req.Name,
		Description: req.Description,
		IconURL:     req.IconURL,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, category)
}

// GetCategory
// @Summary      获取游戏分类
// @Tags         Admin/GameCategories
// @Security     BearerAuth
// @Param        id   path  int  true  "分类ID"
// @Produce      json
// @Success      200  {object}  GameCategory
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/game-categories/{id} [get]
func (h *GameCategoryHandler) GetCategory(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	category, err := h.svc.GetGameCategory(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, category)
}

// ListCategories
// @Summary      列出游戏分类
// @Tags         Admin/GameCategories
// @Security     BearerAuth
// @Param        page       query  int     false  "页码"
// @Param        pageSize   query  int     false  "每页数量"
// @Param        keyword    query  string  false  "关键词搜索"
// @Param        isActive   query  bool    false  "是否启用"
// @Produce      json
// @Success      200  {array}   GameCategory
// @Router       /admin/game-categories [get]
func (h *GameCategoryHandler) ListCategories(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	keyword := c.Query("keyword")
	var isActive *bool
	if isActiveStr := c.Query("isActive"); isActiveStr != "" {
		if isActiveStr == "true" {
			active := true
			isActive = &active
		} else if isActiveStr == "false" {
			active := false
			isActive = &active
		}
	}

	categories, pagination, err := h.svc.ListGameCategoriesPaged(c.Request.Context(), page, pageSize, keyword, isActive)
	if err != nil {
		respondError(c, err)
		return
	}
	respondList(c, categories, pagination)
}

// UpdateCategory
// @Summary      更新游戏分类
// @Tags         Admin/GameCategories
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                      true  "分类ID"
// @Param        request  body  UpdateCategoryRequest   true  "分类信息"
// @Success      200  {object}  GameCategory
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/game-categories/{id} [put]
func (h *GameCategoryHandler) UpdateCategory(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req UpdateCategoryRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	category, err := h.svc.UpdateGameCategory(c.Request.Context(), id, adminservice.UpdateGameCategoryInput{
		Name:        req.Name,
		Description: req.Description,
		IconURL:     req.IconURL,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondUpdated(c, category)
}

// DeleteCategory
// @Summary      删除游戏分类
// @Tags         Admin/GameCategories
// @Security     BearerAuth
// @Param        id   path  int  true  "分类ID"
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/game-categories/{id} [delete]
func (h *GameCategoryHandler) DeleteCategory(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	err := h.svc.DeleteGameCategory(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
}

// BatchUpdateStatus
// @Summary      批量更新分类状态
// @Description  批量启用/禁用游戏分类
// @Tags         Admin/GameCategories
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchUpdateCategoryStatusRequest  true  "分类ID列表和状态"
// @Success      200  {object}  BatchOperationResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/game-categories/batch/status [post]
func (h *GameCategoryHandler) BatchUpdateStatus(c *gin.Context) {
	var req BatchUpdateCategoryStatusRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	if len(req.CategoryIDs) == 0 {
		respondAPIError(c, adminservice.ErrValidation.WithDetails("category_ids is required"))
		return
	}
	if len(req.CategoryIDs) > 100 {
		respondAPIError(c, adminservice.ErrValidation.WithDetails("maximum 100 categories per batch"))
		return
	}

	result, err := h.svc.BatchUpdateGameCategoriesStatus(c.Request.Context(), req.CategoryIDs, req.IsActive)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, result)
}

// BatchDeleteCategories
// @Summary      批量删除游戏分类
// @Description  批量删除游戏分类，最多100条
// @Tags         Admin/GameCategories
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchDeleteCategoriesRequest  true  "分类ID列表"
// @Success      200  {object}  BatchOperationResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/game-categories/batch/delete [post]
func (h *GameCategoryHandler) BatchDeleteCategories(c *gin.Context) {
	var req BatchDeleteCategoriesRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	if len(req.CategoryIDs) == 0 {
		respondAPIError(c, adminservice.ErrValidation.WithDetails("category_ids is required"))
		return
	}
	if len(req.CategoryIDs) > 100 {
		respondAPIError(c, adminservice.ErrValidation.WithDetails("maximum 100 categories per batch"))
		return
	}

	result, err := h.svc.BatchDeleteGameCategories(c.Request.Context(), req.CategoryIDs)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, result)
}

// ============================================================================
// Request/Response DTOs
// ============================================================================

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"max=500"`
	IconURL     string `json:"iconUrl" binding:"omitempty,max=255,url"`
	SortOrder   int    `json:"sortOrder"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=50"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	IconURL     *string `json:"iconUrl" binding:"omitempty,max=255,url"`
	SortOrder   *int    `json:"sortOrder"`
	IsActive    *bool   `json:"isActive"`
}

// BatchUpdateCategoryStatusRequest 批量更新分类状态请求
type BatchUpdateCategoryStatusRequest struct {
	CategoryIDs []uint64 `json:"categoryIds" binding:"required,min=1,max=100"`
	IsActive    bool     `json:"isActive"`
}

// BatchDeleteCategoriesRequest 批量删除分类请求
type BatchDeleteCategoriesRequest struct {
	CategoryIDs []uint64 `json:"categoryIds" binding:"required,min=1,max=100"`
}
