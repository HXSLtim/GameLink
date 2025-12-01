package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	menusvc "gamelink/internal/service/menu"
)

type MenuHandler struct {
	svc *menusvc.Service
}

func NewMenuHandler(svc *menusvc.Service) *MenuHandler {
	return &MenuHandler{svc: svc}
}

// RegisterMenuRoutes 管理端菜单 CRUD
func RegisterMenuRoutes(router gin.IRouter, svc *menusvc.Service) {
	h := NewMenuHandler(svc)
	group := router.Group("/menus")
	group.GET("", h.List)
	group.POST("", h.Create)
	group.GET("/:id", h.Get)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

// List 菜单列表（可选 parentId）
func (h *MenuHandler) List(c *gin.Context) {
	parentID, _ := queryUint64Ptr(c, "parentId")

	// Check if pagination is requested
	pageStr := c.Query("page")
	pageSizeStr := c.Query("pageSize")
	if pageSizeStr == "" {
		pageSizeStr = c.Query("page_size")
	}

	if pageStr == "" && pageSizeStr == "" {
		// No pagination params, return all (backward compatibility)
		menus, err := h.svc.List(c.Request.Context(), parentID)
		if err != nil {
			writeJSONError(c, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(c, http.StatusOK, model.APIResponse[[]model.Menu]{
			Success: true,
			Code:    http.StatusOK,
			Message: "OK",
			Data:    ensureSlice(menus),
		})
		return
	}

	// Pagination requested
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	menus, total, err := h.svc.ListPaged(c.Request.Context(), page, pageSize, parentID)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := 0
	if pageSize > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	pagination := &model.Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	writeJSON(c, http.StatusOK, model.APIResponse[[]model.Menu]{
		Success:    true,
		Code:       http.StatusOK,
		Message:    "OK",
		Data:       ensureSlice(menus),
		Pagination: pagination,
	})
}

// Create 创建菜单
func (h *MenuHandler) Create(c *gin.Context) {
	var req model.Menu
	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusCreated, model.APIResponse[model.Menu]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "Created",
		Data:    req,
	})
}

// Get 获取菜单详情
func (h *MenuHandler) Get(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	menu, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			writeJSONError(c, http.StatusNotFound, "menu not found")
			return
		}
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[model.Menu]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *menu,
	})
}

// Update 更新菜单
func (h *MenuHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req model.Menu
	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	req.ID = id
	if err := h.svc.Update(c.Request.Context(), &req); err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "Updated",
		Data:    req,
	})
}

// Delete 删除菜单
func (h *MenuHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "Deleted",
	})
}
