package admin

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	menusvc "gamelink/internal/service/admin"
	"gamelink/pkg/apierr"
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
			respondError(c, err)
			return
		}
		respondSuccess(c, ensureSlice(menus))
		return
	}

	// Pagination requested
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	menus, total, err := h.svc.ListPaged(c.Request.Context(), page, pageSize, parentID)
	if err != nil {
		respondError(c, err)
		return
	}

	totalPages := 0
	if pageSize > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	respondList(c, ensureSlice(menus), &model.Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	})
}

// Create 创建菜单
func (h *MenuHandler) Create(c *gin.Context) {
	var req model.Menu
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, req)
}

// Get 获取菜单详情
func (h *MenuHandler) Get(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	menu, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			respondError(c, apierr.NotFound("menu not found"))
			return
		}
		respondError(c, err)
		return
	}
	respondSuccess(c, *menu)
}

// Update 更新菜单
func (h *MenuHandler) Update(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	var req model.Menu
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}
	req.ID = id
	if err := h.svc.Update(c.Request.Context(), &req); err != nil {
		respondError(c, err)
		return
	}
	respondUpdated(c, req)
}

// Delete 删除菜单
func (h *MenuHandler) Delete(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
}

// ListMyMenus 获取当前用户可访问的菜单
// @Summary      获取当前用户菜单
// @Description  根据当前用户权限返回可访问的菜单列表，超级管理员返回所有菜单
// @Tags         Admin - Menus
// @Security     BearerAuth
// @Success      200  {array}   model.Menu
// @Router       /admin/menus/me [get]
// @Router       /admin/me/menus [get]
func (h *MenuHandler) ListMyMenus(c *gin.Context) {
	userIDVal, ok := c.Get(middleware.UserIDKey)
	if !ok {
		respondError(c, apierr.Unauthorized("未登录"))
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
			respondError(c, permErr)
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
		respondError(c, err)
		return
	}

	// 过滤隐藏菜单（hidden=true 的菜单不应该出现在用户可访问的菜单列表中）
	menus = filterHiddenMenus(menus)

	// 构建菜单树
	menuTree := buildMenuTree(menus)

	respondSuccess(c, menuTree)
}

// filterHiddenMenus 过滤掉隐藏的菜单
func filterHiddenMenus(menus []model.Menu) []model.Menu {
	result := make([]model.Menu, 0, len(menus))
	for _, menu := range menus {
		if !menu.Hidden {
			result = append(result, menu)
		}
	}
	return result
}

// buildMenuTree 将扁平菜单列表构建为树形结构
func buildMenuTree(menus []model.Menu) []model.Menu {
	if len(menus) == 0 {
		return []model.Menu{}
	}

	// 创建菜单副本并建立 ID -> Menu 指针映射
	menuMap := make(map[uint64]*model.Menu)
	menuCopies := make([]model.Menu, len(menus))
	copy(menuCopies, menus)

	for i := range menuCopies {
		menuCopies[i].Children = []model.Menu{}
		menuMap[menuCopies[i].ID] = &menuCopies[i]
	}

	// 收集根节点 ID
	var rootIDs []uint64
	for i := range menuCopies {
		menu := &menuCopies[i]
		if menu.ParentID == nil || *menu.ParentID == 0 {
			// 根节点
			rootIDs = append(rootIDs, menu.ID)
		} else if parent, ok := menuMap[*menu.ParentID]; ok {
			// 添加到父节点的 children（使用指针操作，直接修改 map 中的数据）
			parent.Children = append(parent.Children, *menu)
		} else {
			// 父菜单不在权限范围内，作为根节点
			rootIDs = append(rootIDs, menu.ID)
		}
	}

	// 递归更新所有节点的 children
	var updateChildren func(menu *model.Menu)
	updateChildren = func(menu *model.Menu) {
		for i := range menu.Children {
			childID := menu.Children[i].ID
			if childPtr, ok := menuMap[childID]; ok {
				menu.Children[i].Children = childPtr.Children
				updateChildren(&menu.Children[i])
			}
		}
	}

	// 构建最终的根节点列表
	roots := make([]model.Menu, 0, len(rootIDs))
	for _, id := range rootIDs {
		if menuPtr, ok := menuMap[id]; ok {
			updateChildren(menuPtr)
			roots = append(roots, *menuPtr)
		}
	}

	return roots
}

// ============================================================================
// 批量操作相关请求和响应结构
// ============================================================================

// BatchDeleteMenusRequest 批量删除菜单请求
type BatchDeleteMenusRequest struct {
	IDs []uint64 `json:"ids" binding:"required,min=1"`
}

// BatchUpdateMenuStatusRequest 批量更新菜单状态请求
type BatchUpdateMenuStatusRequest struct {
	Updates []struct {
		ID     uint64 `json:"id" binding:"required"`
		Hidden *bool  `json:"hidden"`
	} `json:"updates" binding:"required,min=1"`
}

// BatchUpdateMenuSortRequest 批量更新菜单排序请求
type BatchUpdateMenuSortRequest struct {
	Updates []struct {
		ID        uint64 `json:"id" binding:"required"`
		SortOrder int    `json:"sortOrder"`
	} `json:"updates" binding:"required,min=1"`
}

// ============================================================================
// 批量操作处理器方法
// ============================================================================

// BatchDelete 批量删除菜单
// @Summary      批量删除菜单
// @Description  批量删除多个菜单，返回成功和失败的数量
// @Tags         Admin - Menus
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                      true  "Bearer {token}"
// @Param        request        body      BatchDeleteMenusRequest      true  "批量删除请求"
// @Success      200            {object}  menusvc.BatchDeleteResult
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/menus/batch [delete]
func (h *MenuHandler) BatchDelete(c *gin.Context) {
	var req BatchDeleteMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "参数验证失败: "+err.Error())
		return
	}

	result, err := h.svc.BatchDelete(c.Request.Context(), req.IDs)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccessWithMsg(c, "批量删除完成", result)
}

// BatchUpdateStatus 批量更新菜单状态
// @Summary      批量更新菜单状态
// @Description  批量更新多个菜单的显示/隐藏状态，返回成功和失败的数量
// @Tags         Admin - Menus
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                          true  "Bearer {token}"
// @Param        request        body      BatchUpdateMenuStatusRequest    true  "批量更新状态请求"
// @Success      200            {object}  menusvc.BatchUpdateStatusResult
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/menus/batch/status [put]
func (h *MenuHandler) BatchUpdateStatus(c *gin.Context) {
	var req BatchUpdateMenuStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "参数验证失败: "+err.Error())
		return
	}

	// 转换请求格式
	updates := make([]menusvc.MenuStatusUpdate, len(req.Updates))
	for i, u := range req.Updates {
		updates[i] = menusvc.MenuStatusUpdate{
			ID:     u.ID,
			Hidden: u.Hidden,
		}
	}

	result, err := h.svc.BatchUpdateStatus(c.Request.Context(), updates)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccessWithMsg(c, "批量状态更新完成", result)
}

// BatchUpdateSort 批量更新菜单排序
// @Summary      批量更新菜单排序
// @Description  批量更新多个菜单的排序值，返回成功和失败的数量
// @Tags         Admin - Menus
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                        true  "Bearer {token}"
// @Param        request        body      BatchUpdateMenuSortRequest    true  "批量更新排序请求"
// @Success      200            {object}  menusvc.BatchUpdateSortResult
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/menus/batch/sort [put]
func (h *MenuHandler) BatchUpdateSort(c *gin.Context) {
	var req BatchUpdateMenuSortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "参数验证失败: "+err.Error())
		return
	}

	// 转换请求格式
	updates := make([]menusvc.MenuSortUpdate, len(req.Updates))
	for i, u := range req.Updates {
		updates[i] = menusvc.MenuSortUpdate{
			ID:        u.ID,
			SortOrder: u.SortOrder,
		}
	}

	result, err := h.svc.BatchUpdateSort(c.Request.Context(), updates)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccessWithMsg(c, "批量排序更新完成", result)
}
