package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service/contentcategory"
)

// ContentCategoryHandler 内容分类管理处理器
type ContentCategoryHandler struct {
	svc *contentcategory.ContentCategoryService
}

// NewContentCategoryHandler 创建内容分类管理处理器
func NewContentCategoryHandler(svc *contentcategory.ContentCategoryService) *ContentCategoryHandler {
	return &ContentCategoryHandler{svc: svc}
}

// List 列出内容分类
// @Summary      列出内容分类
// @Tags         Admin - Content Category
// @Security     BearerAuth
// @Param        page     query     int     false  "页码"
// @Param        pageSize query     int     false  "每页数量"
// @Param        keyword  query     string  false  "关键词"
// @Param        status   query     string  false  "状态"
// @Success      200  {object}  model.APIResponse[contentcategory.ListResponse]
// @Router       /admin/content/categories [get]
func (h *ContentCategoryHandler) List(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	var status *model.ContentCategoryStatus
	if s := c.Query("status"); s != "" {
		st := model.ContentCategoryStatus(s)
		status = &st
	}

	resp, err := h.svc.List(c.Request.Context(), contentcategory.ListRequest{
		Page:     page,
		PageSize: pageSize,
		Keyword:  c.Query("keyword"),
		Status:   status,
	})
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[*contentcategory.ListResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    resp,
	})
}

// Get 获取分类详情
// @Summary      获取分类详情
// @Tags         Admin - Content Category
// @Security     BearerAuth
// @Param        id   path      int  true  "分类ID"
// @Success      200  {object}  model.APIResponse[contentcategory.CategoryDTO]
// @Router       /admin/content/categories/{id} [get]
func (h *ContentCategoryHandler) Get(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}

	category, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		writeJSONError(c, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[*contentcategory.CategoryDTO]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    category,
	})
}

// Create 创建分类
// @Summary      创建分类
// @Tags         Admin - Content Category
// @Security     BearerAuth
// @Param        body body      contentcategory.CreateRequest  true  "创建请求"
// @Success      201  {object}  model.APIResponse[contentcategory.CategoryDTO]
// @Router       /admin/content/categories [post]
func (h *ContentCategoryHandler) Create(c *gin.Context) {
	var req contentcategory.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		if err == contentcategory.ErrDuplicate {
			writeJSONError(c, http.StatusConflict, "分类名称已存在")
			return
		}
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusCreated, model.APIResponse[*contentcategory.CategoryDTO]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "创建成功",
		Data:    category,
	})
}

// Update 更新分类
// @Summary      更新分类
// @Tags         Admin - Content Category
// @Security     BearerAuth
// @Param        id   path      int                            true  "分类ID"
// @Param        body body      contentcategory.UpdateRequest  true  "更新请求"
// @Success      200  {object}  model.APIResponse[any]
// @Router       /admin/content/categories/{id} [put]
func (h *ContentCategoryHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req contentcategory.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.svc.Update(c.Request.Context(), id, req); err != nil {
		if err == contentcategory.ErrNotFound {
			writeJSONError(c, http.StatusNotFound, "分类不存在")
			return
		}
		if err == contentcategory.ErrDuplicate {
			writeJSONError(c, http.StatusConflict, "分类名称已存在")
			return
		}
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "更新成功",
	})
}

// Delete 删除分类
// @Summary      删除分类
// @Tags         Admin - Content Category
// @Security     BearerAuth
// @Param        id                  path      int  true   "分类ID"
// @Param        migrateToCategoryId query     int  false  "迁移目标分类ID"
// @Success      200  {object}  model.APIResponse[any]
// @Router       /admin/content/categories/{id} [delete]
func (h *ContentCategoryHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}

	migrateTo, _ := queryUint64Ptr(c, "migrateToCategoryId")

	if err := h.svc.Delete(c.Request.Context(), id, migrateTo); err != nil {
		if err == contentcategory.ErrNotFound {
			writeJSONError(c, http.StatusNotFound, "分类不存在")
			return
		}
		if err == contentcategory.ErrHasFeeds {
			writeJSONError(c, http.StatusBadRequest, "分类下有动态，请指定迁移目标分类")
			return
		}
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "删除成功",
	})
}
