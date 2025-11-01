package admin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service"
)

// PermissionHandler 权限管理处理器。
type PermissionHandler struct {
	permissionSvc *service.PermissionService
}

// NewPermissionHandler 创建权限处理器实例。
func NewPermissionHandler(permissionSvc *service.PermissionService) *PermissionHandler {
	return &PermissionHandler{permissionSvc: permissionSvc}
}

// ListPermissions 获取权限列表。
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	keyword := c.Query("keyword")
	method := c.Query("method")
	group := c.Query("group")

	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	var permissions []model.Permission
	var total int64
	var err error

	// 如果有过滤条件，使用过滤查询
	if keyword != "" || method != "" || group != "" {
		permissions, total, err = h.permissionSvc.ListPermissionsPagedWithFilter(c.Request.Context(), page, pageSize, keyword, method, group)
	} else {
		permissions, total, err = h.permissionSvc.ListPermissionsPaged(c.Request.Context(), page, pageSize)
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
			"items":      ensureSlice(permissions),
			"page":       page,
			"pageSize":   pageSize,
			"totalCount": total,
		},
	})
}

// GetPermission 获取权限详情。
func (h *PermissionHandler) GetPermission(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "无效的权限ID")
		return
	}

	permission, err := h.permissionSvc.GetPermission(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeJSONError(c, http.StatusNotFound, "权限不存在")
		} else {
			writeJSONError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[*model.Permission]{
		Success: true,
		Code:    http.StatusOK,
		Message: "成功",
		Data:    permission,
	})
}

// CreatePermission 创建权限。
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req struct {
		Method      model.HTTPMethod `json:"method" binding:"required"`
		Path        string           `json:"path" binding:"required,max=255"`
		Code        string           `json:"code" binding:"max=128"`
		Group       string           `json:"group" binding:"max=64"`
		Description string           `json:"description" binding:"max=255"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	permission := &model.Permission{
		Method:      req.Method,
		Path:        req.Path,
		Code:        req.Code,
		Group:       req.Group,
		Description: req.Description,
	}

	if err := h.permissionSvc.CreatePermission(c.Request.Context(), permission); err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusCreated, model.APIResponse[*model.Permission]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "权限创建成功",
		Data:    permission,
	})
}

// UpdatePermission 更新权限。
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "无效的权限ID")
		return
	}

	var req struct {
		Code        string `json:"code" binding:"max=128"`
		Group       string `json:"group" binding:"max=64"`
		Description string `json:"description" binding:"max=255"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	permission := &model.Permission{
		Base:        model.Base{ID: id},
		Code:        req.Code,
		Group:       req.Group,
		Description: req.Description,
	}

	if err := h.permissionSvc.UpdatePermission(c.Request.Context(), permission); err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	updatedPermission, err := h.permissionSvc.GetPermission(c.Request.Context(), id)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[*model.Permission]{
		Success: true,
		Code:    http.StatusOK,
		Message: "权限更新成功",
		Data:    updatedPermission,
	})
}

// DeletePermission 删除权限。
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "无效的权限ID")
		return
	}

	if err := h.permissionSvc.DeletePermission(c.Request.Context(), id); err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "权限删除成功",
		Data:    nil,
	})
}

// GetRolePermissions 获取角色的权限列表。
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

// GetUserPermissions 获取用户的权限列表。
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

// GetPermissionGroups 获取所有权限分组列表。
func (h *PermissionHandler) GetPermissionGroups(c *gin.Context) {
	groups, err := h.permissionSvc.ListPermissionGroups(c.Request.Context())
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[[]string]{
		Success: true,
		Code:    http.StatusOK,
		Message: "成功",
		Data:    ensureSlice(groups),
	})
}
