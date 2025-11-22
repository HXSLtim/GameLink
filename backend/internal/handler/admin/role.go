package admin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	roleservice "gamelink/internal/service/role"
)

// RoleHandler 角色管理处理type RoleHandler struct {
	roleSvc *roleservice.RoleService
}

// NewRoleHandler 创建角色处理器实func NewRoleHandler(roleSvc *roleservice.RoleService) *RoleHandler {
	return &RoleHandler{roleSvc: roleSvc}
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Slug        string `json:"slug" binding:"required,max=64"`
	Name        string `json:"name" binding:"required,max=128"`
	Description string `json:"description" binding:"max=255"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        string `json:"name" binding:"omitempty,max=128"`
	Description string `json:"description" binding:"max=255"`
}

// AssignPermissionsRequest 分配权限请求
type AssignPermissionsRequest struct {
	PermissionIDs []uint64 `json:"permissionIds" binding:"required"`
}

// AssignRolesToUserRequest 分配角色请求
type AssignRolesToUserRequest struct {
	UserID  uint64   `json:"userId" binding:"required"`
	RoleIDs []uint64 `json:"roleIds" binding:"required"`
}

// ListRoles 获取角色列表
// @Summary      获取角色列表
// @Description  管理员获取系统角色列表，支持分页和过// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization     header    string  true  "Bearer {token}"
// @Param        page              query     int     false  "页码" default(1)
// @Param        pageSize          query     int     false  "每页数量" default(10)
// @Param        with_permissions  query     bool    false  "是否包含权限信息"
// @Param        keyword           query     string  false  "关键词搜
// @Param        isSystem          query     bool    false  "是否为系统角
// @Success      200               {object}  model.APIResponse[gin.H]
// @Failure      400               {object}  model.ErrorResponse
// @Failure      401               {object}  model.ErrorResponse
// @Failure      500               {object}  model.ErrorResponse
// @Router       /admin/roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	withPermissions := c.Query("with_permissions") == "true"
	keyword := c.Query("keyword")
	isSystemStr := c.Query("isSystem")

	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	var roles []model.RoleModel
	var total int64
	var err error

	if withPermissions {
		roles, err = h.roleSvc.ListRolesWithPermissions(c.Request.Context())
		total = int64(len(roles))
	} else {
		// 解析 isSystem 参数
		var isSystem *bool
		if isSystemStr != "" {
			val := isSystemStr == "true"
			isSystem = &val
		}

		// 如果有过滤条件，使用过滤查询
		if keyword != "" || isSystem != nil {
			roles, total, err = h.roleSvc.ListRolesPagedWithFilter(c.Request.Context(), page, pageSize, keyword, isSystem)
		} else {
			roles, total, err = h.roleSvc.ListRolesPaged(c.Request.Context(), page, pageSize)
		}
	}

	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "成功",
		Data: gin.H{
			"items":      ensureSlice(roles),
			"page":       page,
			"pageSize":   pageSize,
			"totalCount": total,
		},
	})
}

// GetRole 获取角色详情
// @Summary      获取角色详情
// @Description  管理员根据ID获取指定角色的详细信// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization     header    string  true  "Bearer {token}"
// @Param        id                path      uint    true  "角色ID"
// @Param        with_permissions  query     bool    false  "是否包含权限信息"
// @Success      200               {object}  model.APIResponse[gin.H]
// @Failure      400               {object}  model.ErrorResponse
// @Failure      401               {object}  model.ErrorResponse
// @Failure      404               {object}  model.ErrorResponse
// @Failure      500               {object}  model.ErrorResponse
// @Router       /admin/roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	withPermissions := c.Query("with_permissions") == "true"

	var role *model.RoleModel
	if withPermissions {
		role, err = h.roleSvc.GetRoleWithPermissions(c.Request.Context(), id)
	} else {
		role, err = h.roleSvc.GetRole(c.Request.Context(), id)
	}

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeJSONError(c, http.StatusNotFound, "角色不存)
		} else {
			writeJSONError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[model.RoleModel]{
		Success: true,
		Code:    http.StatusOK,
		Message: "成功",
		Data:    role,
	})
}

// CreateRole 创建角色
// @Summary      创建角色
// @Description  管理员创建新的系统角// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                true  "Bearer {token}"
// @Param        request        body      CreateRoleRequest      true  "创建角色请求"
// @Success      201            {object}  model.APIResponse[gin.H]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	role := &model.RoleModel{
		Slug:        req.Slug,
		Name:        req.Name,
		Description: req.Description,
		IsSystem:    false,
	}

	if err := h.roleSvc.CreateRole(c.Request.Context(), role); err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusCreated, model.APIResponse[model.RoleModel]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "角色创建成功",
		Data:    role,
	})
}

// UpdateRole 更新角色
// @Summary      更新角色
// @Description  管理员更新指定角色的信息
// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                true  "Bearer {token}"
// @Param        id             path      uint                  true  "角色ID"
// @Param        request        body      UpdateRoleRequest      true  "更新角色请求"
// @Success      200            {object}  model.APIResponse[gin.H]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      404            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	var req UpdateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	role := &model.RoleModel{
		Base:        model.Base{ID: id},
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.roleSvc.UpdateRole(c.Request.Context(), role); err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	updatedRole, err := h.roleSvc.GetRole(c.Request.Context(), id)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[model.RoleModel]{
		Success: true,
		Code:    http.StatusOK,
		Message: "角色更新成功",
		Data:    updatedRole,
	})
}

// DeleteRole 删除角色
// @Summary      删除角色
// @Description  管理员删除指定的系统角色
// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      uint    true  "角色ID"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	if err := h.roleSvc.DeleteRole(c.Request.Context(), id); err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "角色删除成功",
		Data:    nil,
	})
}

// AssignPermissions 为角色分配权// @Summary      为角色分配权// @Description  管理员为指定角色分配多个权限
// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                      true  "Bearer {token}"
// @Param        id             path      uint                        true  "角色ID"
// @Param        request        body      AssignPermissionsRequest     true  "分配权限请求"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/roles/{id}/permissions [post]
func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	var req AssignPermissionsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	if err := h.roleSvc.AssignPermissionsToRole(c.Request.Context(), id, req.PermissionIDs); err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "权限分配成功",
		Data:    nil,
	})
}

// AssignRolesToUser 为用户分配角// @Summary      为用户分配角// @Description  管理员为指定用户分配多个角色
// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                    true  "Bearer {token}"
// @Param        request        body      AssignRolesToUserRequest    true  "分配角色请求"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/users/roles [post]
func (h *RoleHandler) AssignRolesToUser(c *gin.Context) {
	var req AssignRolesToUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	if err := h.roleSvc.AssignRolesToUser(c.Request.Context(), req.UserID, req.RoleIDs); err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "角色分配成功",
		Data:    nil,
	})
}

// GetUserRoles 获取用户的角色列// @Summary      获取用户的角色列// @Description  管理员获取指定用户的角色列表
// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        user_id         path      uint    true  "用户ID"
// @Success      200            {object}  model.APIResponse[[]model.RoleModel]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/users/{user_id}/roles [get]
func (h *RoleHandler) GetUserRoles(c *gin.Context) {
	userID, err := parseUintParam(c, "user_id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	roles, err := h.roleSvc.ListRolesByUserID(c.Request.Context(), userID)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[[]model.RoleModel]{
		Success: true,
		Code:    http.StatusOK,
		Message: "成功",
		Data:    ensureSlice(roles),
	})
}
