package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	permissionservice "gamelink/internal/service/admin"
)

// Permission 权限模型（类型别名）
type Permission = model.Permission

// PermissionHandler 权限管理处理器
type PermissionHandler struct {
	permissionSvc *permissionservice.PermissionService
	roleSvc       *permissionservice.RoleService
}

// NewPermissionHandler 创建权限处理器实例
func NewPermissionHandler(permissionSvc *permissionservice.PermissionService) *PermissionHandler {
	return &PermissionHandler{permissionSvc: permissionSvc}
}

// NewPermissionHandlerWithRoleService 创建带角色服务的权限处理器实例
func NewPermissionHandlerWithRoleService(permissionSvc *permissionservice.PermissionService, roleSvc *permissionservice.RoleService) *PermissionHandler {
	return &PermissionHandler{permissionSvc: permissionSvc, roleSvc: roleSvc}
}

// ListPermissions 获取权限列表
// @Summary      获取权限列表
// @Description  管理员获取系统权限列表，支持分页和过滤
// @Tags         Admin - Permissions
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        page           query     int     false  "页码" default(1)
// @Param        pageSize       query     int     false  "每页数量" default(20)
// @Param        keyword        query     string  false  "关键词搜索（匹配code, path, description）"
// @Param        method         query     string  false  "HTTP方法过滤（GET, POST, PUT, PATCH, DELETE）"
// @Param        group          query     string  false  "权限分组过滤"
// @Success      200            {object}  model.APIResponse[[]Permission]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/permissions [get]
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	// 构建查询选项
	opts := permissionservice.PermissionListOptions{
		Page:     page,
		PageSize: pageSize,
		Keyword:  c.Query("keyword"),
		Method:   c.Query("method"),
		Group:    c.Query("group"),
	}

	// 解析 is_system 参数
	if isSystemStr := c.Query("is_system"); isSystemStr != "" {
		isSystem := isSystemStr == "true"
		opts.IsSystem = &isSystem
	}

	permissions, total, err := h.permissionSvc.ListPermissionsPagedWithFilter(c.Request.Context(), opts)
	if err != nil {
		respondError(c, err)
		return
	}

	// 计算总页数
	totalInt := int(total)
	totalPages := (totalInt + pageSize - 1) / pageSize

	respondList(c, permissions, &model.Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      totalInt,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	})
}

// GetPermission 获取权限详情
// @Summary      获取权限详情
// @Description  管理员根据ID获取指定权限的详细信息
// @Tags         Admin - Permissions
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      uint    true  "权限ID"
// @Success      200            {object}  model.APIResponse[Permission]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      404            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/permissions/{id} [get]
func (h *PermissionHandler) GetPermission(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	permission, err := h.permissionSvc.GetPermission(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, permission)
}

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Method      model.HTTPMethod `json:"method" binding:"required"`
	Path        string           `json:"path" binding:"required,max=255"`
	Code        string           `json:"code" binding:"max=128"`
	Group       string           `json:"group" binding:"max=64"`
	Description string           `json:"description" binding:"max=255"`
	ParentID    *uint64          `json:"parentId"`
	SortOrder   int              `json:"sortOrder"`
	IsSystem    bool             `json:"isSystem"`
}

// UpdatePermissionRequest 全量更新权限请求
type UpdatePermissionRequest struct {
	Group       string  `json:"group" binding:"max=64"`
	Description string  `json:"description" binding:"max=255"`
	ParentID    *uint64 `json:"parentId"`
	SortOrder   int     `json:"sortOrder"`
}

// PatchPermissionRequest 部分更新权限请求
type PatchPermissionRequest struct {
	Code        *string `json:"code,omitempty"`
	Group       *string `json:"group,omitempty"`
	Description *string `json:"description,omitempty"`
	SortOrder   *int    `json:"sortOrder,omitempty"`
}

// CreatePermission 创建权限
// @Summary      创建权限
// @Description  管理员创建新的系统权限，权限码格式为 module.resource.action
// @Tags         Admin - Permissions
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                      true  "Bearer {token}"
// @Param        request        body      CreatePermissionRequest      true  "创建权限请求"
// @Success      201            {object}  model.APIResponse[Permission]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/permissions [post]
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	permission := &model.Permission{
		Method:      req.Method,
		Path:        req.Path,
		Code:        req.Code,
		Group:       req.Group,
		Description: req.Description,
		ParentID:    req.ParentID,
		SortOrder:   req.SortOrder,
		IsSystem:    req.IsSystem,
	}

	if err := h.permissionSvc.CreatePermission(c.Request.Context(), permission); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, permission)
}

// UpdatePermission 全量更新权限
// @Summary      全量更新权限
// @Description  管理员全量更新指定权限的信息（权限码创建后不可修改）
// @Tags         Admin - Permissions
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                      true  "Bearer {token}"
// @Param        id             path      uint                        true  "权限ID"
// @Param        request        body      UpdatePermissionRequest      true  "更新权限请求"
// @Success      200            {object}  model.APIResponse[Permission]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      404            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/permissions/{id} [put]
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req UpdatePermissionRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	// 获取现有权限以保留不可修改的字段
	existing, err := h.permissionSvc.GetPermission(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	permission := &model.Permission{
		Base:        model.Base{ID: id},
		Method:      existing.Method,
		Path:        existing.Path,
		Code:        existing.Code, // 权限码不可修改
		Group:       req.Group,
		Description: req.Description,
		ParentID:    req.ParentID,
		SortOrder:   req.SortOrder,
		IsSystem:    existing.IsSystem, // 系统标记不可修改
	}

	if err := h.permissionSvc.UpdatePermission(c.Request.Context(), permission); err != nil {
		respondError(c, err)
		return
	}

	updatedPermission, err := h.permissionSvc.GetPermission(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondUpdated(c, updatedPermission)
}

// PatchPermission 部分更新权限
// @Summary      部分更新权限
// @Description  管理员部分更新指定权限的信息（只更新提供的字段，权限码创建后不可修改）
// @Tags         Admin - Permissions
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                      true  "Bearer {token}"
// @Param        id             path      uint                        true  "权限ID"
// @Param        request        body      PatchPermissionRequest       true  "部分更新权限请求"
// @Success      200            {object}  model.APIResponse[Permission]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      404            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/permissions/{id} [patch]
func (h *PermissionHandler) PatchPermission(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req PatchPermissionRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	// 构建更新字段映射
	updates := make(map[string]interface{})
	if req.Code != nil {
		updates["code"] = *req.Code
	}
	if req.Group != nil {
		updates["group"] = *req.Group
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.SortOrder != nil {
		updates["sortOrder"] = float64(*req.SortOrder)
	}

	if len(updates) == 0 {
		respondBadRequest(c, "没有提供要更新的字段")
		return
	}

	updatedPermission, err := h.permissionSvc.PartialUpdatePermission(c.Request.Context(), id, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	respondUpdated(c, updatedPermission)
}

// DeletePermission 删除权限（软删除）
// @Summary      删除权限
// @Description  管理员删除指定的系统权限（系统权限不可删除，被角色引用的权限需要先解除引用）
// @Tags         Admin - Permissions
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      uint    true  "权限ID"
// @Param        force          query     bool    false "强制删除（忽略引用检查）"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      404            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/permissions/{id} [delete]
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	// 检查是否强制删除
	force := c.Query("force") == "true"

	var err error
	if force {
		err = h.permissionSvc.DeletePermissionForce(c.Request.Context(), id)
	} else {
		err = h.permissionSvc.DeletePermission(c.Request.Context(), id)
	}

	if err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

// GetRolePermissions 获取角色的权限列// @Summary      获取角色的权限列// @Description  管理员获取指定角色的所有权限列// @Tags         Admin - Permissions
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      uint    true  "角色ID"
// @Success      200            {object}  model.APIResponse[[]model.Permission]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/roles/{id}/permissions [get]
func (h *PermissionHandler) GetRolePermissions(c *gin.Context) {
	roleID, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	permissions, err := h.permissionSvc.ListPermissionsByRoleID(c.Request.Context(), roleID)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[[]model.Permission]{
		Success: true,
		Code:    http.StatusOK,
		Message: "成功",
		Data:    ensureSlice(permissions),
	})
}

// GetUserPermissions 获取用户的权限列// @Summary      获取用户的权限列// @Description  管理员获取指定用户的所有权限列// @Tags         Admin - Permissions
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      uint    true  "用户ID"
// @Success      200            {object}  model.APIResponse[[]model.Permission]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/users/{id}/permissions [get]
func (h *PermissionHandler) GetUserPermissions(c *gin.Context) {
	userID, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	permissions, err := h.permissionSvc.ListPermissionsByUserID(c.Request.Context(), userID)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[[]model.Permission]{
		Success: true,
		Code:    http.StatusOK,
		Message: "成功",
		Data:    ensureSlice(permissions),
	})
}

// GetPermissionGroups 获取所有权限分组列表
// @Summary      获取权限分组列表
// @Description  管理员获取系统中所有权限分组的列表
// @Tags         Admin - Permissions
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.APIResponse[[]string]
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/permissions/groups [get]
func (h *PermissionHandler) GetPermissionGroups(c *gin.Context) {
	groups, err := h.permissionSvc.ListPermissionGroups(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, ensureSlice(groups))
}

// GetPermissionTree 获取权限树形结构
// @Summary      获取权限树形结构
// @Description  管理员获取按父子关系组织的权限树形结构，用于角色权限配置页面
// @Tags         Admin - Permissions
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.APIResponse[[]model.PermissionTreeNode]
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/permissions/tree [get]
func (h *PermissionHandler) GetPermissionTree(c *gin.Context) {
	tree, err := h.permissionSvc.GetPermissionTree(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}

	// 确保返回空数组而非 nil
	if tree == nil {
		tree = make([]*model.PermissionTreeNode, 0)
	}

	respondSuccess(c, tree)
}

// GetPermissionTreeByGroup 获取按分组组织的权限树形结构
// @Summary      获取按分组组织的权限树
// @Description  管理员获取按分组分类的权限树形结构，便于按模块查看权限
// @Tags         Admin - Permissions
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.APIResponse[[]model.PermissionGroup]
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/permissions/tree/grouped [get]
func (h *PermissionHandler) GetPermissionTreeByGroup(c *gin.Context) {
	tree, err := h.permissionSvc.GetPermissionTreeByGroup(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}

	// 确保返回空数组而非 nil
	if tree == nil {
		tree = make([]model.PermissionGroup, 0)
	}

	respondSuccess(c, tree)
}

// GetCurrentUserPermissions 获取当前管理员的权限列表
// @Summary      获取当前用户权限
// @Description  返回当前登录管理员拥有的权限码列表，超级管理员返回 ['*']
// @Tags         Admin - Permissions
// @Security     BearerAuth
// @Success      200  {object}  model.APIResponse[[]string]
// @Router       /admin/permissions/me [get]
// @Router       /admin/me/permissions [get]
func (h *PermissionHandler) GetCurrentUserPermissions(c *gin.Context) {
	userIDVal, ok := c.Get(middleware.UserIDKey)
	if !ok {
		writeJSONError(c, http.StatusUnauthorized, "未登录")
		return
	}
	userID, _ := userIDVal.(uint64)

	// 检查是否为超级管理员（Requirements 5.3）
	if h.roleSvc != nil {
		isSuperAdmin, err := h.roleSvc.CheckUserIsSuperAdmin(c.Request.Context(), userID)
		if err == nil && isSuperAdmin {
			// 超级管理员返回 ['*']
			writeJSON(c, http.StatusOK, model.APIResponse[[]string]{
				Success: true,
				Code:    http.StatusOK,
				Message: "OK",
				Data:    []string{"*"},
			})
			return
		}
	}

	perms, err := h.permissionSvc.ListPermissionsByUserID(c.Request.Context(), userID)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Extract codes
	var codes []string
	for _, p := range perms {
		if p.Code != "" {
			codes = append(codes, p.Code)
		}
	}

	writeJSON(c, http.StatusOK, model.APIResponse[[]string]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    ensureSlice(codes),
	})
}
