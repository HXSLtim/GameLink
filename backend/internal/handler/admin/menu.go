package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	menusvc "gamelink/internal/service/admin"
)

type MenuHandler struct {
	svc     *menusvc.MenuService
	permSvc *menusvc.PermissionService
	roleSvc *menusvc.RoleService
}

func NewMenuHandler(svc *menusvc.MenuService, permSvc *menusvc.PermissionService) *MenuHandler {
	return &MenuHandler{svc: svc, permSvc: permSvc}
}

// NewMenuHandlerWithRoleService 创建带角色服务的菜单处理器实例
func NewMenuHandlerWithRoleService(svc *menusvc.MenuService, permSvc *menusvc.PermissionService, roleSvc *menusvc.RoleService) *MenuHandler {
	return &MenuHandler{svc: svc, permSvc: permSvc, roleSvc: roleSvc}
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

// ListMyMenus 获取当前用户可访问的菜单
// @Summary      获取当前用户菜单
// @Description  根据当前用户权限返回可访问的菜单列表，超级管理员返回所有菜单
// @Tags         Admin - Menus
// @Security     BearerAuth
// @Success      200  {object}  model.APIResponse[[]model.Menu]
// @Router       /admin/menus/me [get]
// @Router       /admin/me/menus [get]
func (h *MenuHandler) ListMyMenus(c *gin.Context) {
	userIDVal, ok := c.Get(middleware.UserIDKey)
	if !ok {
		writeJSONError(c, http.StatusUnauthorized, "未登录")
		return
	}
	userID, _ := userIDVal.(uint64)

	var menus []model.Menu
	var err error

	// 检查是否为超级管理员（Requirements 8.1）
	isSuperAdmin := false
	if h.roleSvc != nil {
		isSuperAdmin, _ = h.roleSvc.CheckUserIsSuperAdmin(c.Request.Context(), userID)
	}

	if isSuperAdmin {
		// 超级管理员返回所有菜单
		menus, err = h.svc.List(c.Request.Context(), nil)
	} else {
		// 获取用户权限列表
		perms, permErr := h.permSvc.ListPermissionsByUserID(c.Request.Context(), userID)
		if permErr != nil {
			writeJSONError(c, http.StatusInternalServerError, permErr.Error())
			return
		}

		// 提取权限码
		var codes []string
		for _, p := range perms {
			if p.Code == "*" {
				// 如果有 * 权限码，也视为超级管理员
				menus, err = h.svc.List(c.Request.Context(), nil)
				goto buildTree
			}
			if p.Code != "" {
				codes = append(codes, p.Code)
			}
		}

		// 根据权限过滤菜单
		menus, err = h.svc.ListAccessible(c.Request.Context(), codes)
	}

buildTree:
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 构建菜单树
	menuTree := buildMenuTree(menus)

	writeJSON(c, http.StatusOK, model.APIResponse[[]model.Menu]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    menuTree,
	})
}

// buildMenuTree 将扁平菜单列表构建为树形结构
func buildMenuTree(menus []model.Menu) []model.Menu {
	menuMap := make(map[uint64]*model.Menu)
	var roots []model.Menu

	// 先建立 ID -> Menu 映射（使用指针）
	for i := range menus {
		menus[i].Children = []model.Menu{}
		menuMap[menus[i].ID] = &menus[i]
	}

	// 构建树
	for i := range menus {
		menu := &menus[i]
		if menu.ParentID == nil || *menu.ParentID == 0 {
			// 根节点
			roots = append(roots, *menu)
		} else if parent, ok := menuMap[*menu.ParentID]; ok {
			// 添加到父节点的 children
			parent.Children = append(parent.Children, *menu)
		} else {
			// 父菜单不在权限范围内，作为根节点
			roots = append(roots, *menu)
		}
	}

	// 更新根节点的 children（因为之前是值拷贝）
	for i := range roots {
		if menuPtr, ok := menuMap[roots[i].ID]; ok {
			roots[i].Children = menuPtr.Children
		}
	}

	return roots
}
