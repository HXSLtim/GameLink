package admin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service"
)

// RoleHandler 角色管理处理器�?type RoleHandler struct {
	roleSvc *service.RoleService
}

// NewRoleHandler 创建角色处理器实例�?func NewRoleHandler(roleSvc *service.RoleService) *RoleHandler {
	return &RoleHandler{roleSvc: roleSvc}
}

// ListRoles 获取角色列表�?func (h *RoleHandler) ListRoles(c *gin.Context) {
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

// GetRole 获取角色详情�?func (h *RoleHandler) GetRole(c *gin.Context) {
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
			writeJSONError(c, http.StatusNotFound, "角色不存�?)
		} else {
			writeJSONError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[*model.RoleModel]{
		Success: true,
		Code:    http.StatusOK,
		Message: "成功",
		Data:    role,
	})
}

// CreateRole 创建角色�?func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req struct {
		Slug        string `json:"slug" binding:"required,max=64"`
		Name        string `json:"name" binding:"required,max=128"`
		Description string `json:"description" binding:"max=255"`
	}

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

	writeJSON(c, http.StatusCreated, model.APIResponse[*model.RoleModel]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "角色创建成功",
		Data:    role,
	})
}

// UpdateRole 更新角色�?func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	var req struct {
		Name        string `json:"name" binding:"omitempty,max=128"`
		Description string `json:"description" binding:"max=255"`
	}

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

	writeJSON(c, http.StatusOK, model.APIResponse[*model.RoleModel]{
		Success: true,
		Code:    http.StatusOK,
		Message: "角色更新成功",
		Data:    updatedRole,
	})
}

// DeleteRole 删除角色�?func (h *RoleHandler) DeleteRole(c *gin.Context) {
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

// AssignPermissions 为角色分配权限�?func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	var req struct {
		PermissionIDs []uint64 `json:"permissionIds" binding:"required"`
	}

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

// AssignRolesToUser 为用户分配角色�?func (h *RoleHandler) AssignRolesToUser(c *gin.Context) {
	var req struct {
		UserID  uint64   `json:"userId" binding:"required"`
		RoleIDs []uint64 `json:"roleIds" binding:"required"`
	}

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

// GetUserRoles 获取用户的角色列表�?func (h *RoleHandler) GetUserRoles(c *gin.Context) {
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
