package admin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	permissionservice "gamelink/internal/service/admin"
)

// Permission 权限模型（类型别名）
type Permission = model.Permission

// PermissionHandler 权限管理处理器
type PermissionHandler struct {
	permissionSvc *permissionservice.PermissionService
}

// NewPermissionHandler 创建权限处理器实例
func NewPermissionHandler(permissionSvc *permissionservice.PermissionService) *PermissionHandler {
	return &PermissionHandler{permissionSvc: permissionSvc}
}

// ListPermissions 获取权限列表
// @Summary      获取权限列表
// @Description  管理员获取系统权限列表，支持分页和过// @Tags         Admin - Permissions
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        page           query     int     false  "页码" default(1)
// @Param        pageSize       query     int     false  "每页数量" default(10)
// @Param        keyword        query     string  false  "关键词搜索"
// @Param        method         query     string  false  "HTTP方法过滤"
// @Param        group          query     string  false  "权限分组过滤"
// @Success      200            {object}  model.APIResponse[gin.H]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/permissions [get]
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	// TODO: 实现keyword, method, group过滤功能
	_ = c.Query("keyword")
	_ = c.Query("method")
	_ = c.Query("group")

	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	var permissions []model.Permission
	var total int64
	var err error

	// 直接调用ListPermissionsPaged
	permissions, total, err = h.permissionSvc.ListPermissionsPaged(c.Request.Context(), page, pageSize)

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

// GetPermission 获取权限详情
// @Summary      获取权限详情
// @Description  管理员根据ID获取指定权限的详细信// @Tags         Admin - Permissions
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

// CreatePermission 创建权限
// @Summary      创建权限
// @Description  管理员创建新的系统权// @Tags         Admin - Permissions
// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Method      model.HTTPMethod `json:"method" binding:"required"`
	Path        string           `json:"path" binding:"required,max=255"`
	Code        string           `json:"code" binding:"max=128"`
	Group       string           `json:"group" binding:"max=64"`
	Description string           `json:"description" binding:"max=255"`
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	Code        string `json:"code" binding:"max=128"`
	Group       string `json:"group" binding:"max=64"`
	Description string `json:"description" binding:"max=255"`
}

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

// UpdatePermission 更新权限
// @Summary      更新权限
// @Description  管理员更新指定权限的信息
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

// DeletePermission 删除权限
// @Summary      删除权限
// @Description  管理员删除指定的系统权限
// @Tags         Admin - Permissions
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      uint    true  "权限ID"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/permissions/{id} [delete]
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

// GetPermissionGroups 获取所有权限分组列// @Summary      获取权限分组列表
// @Description  管理员获取系统中所有权限分组的列表
// @Tags         Admin - Permissions
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.APIResponse[[]string]
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/permission-groups [get]
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

// GetCurrentUserPermissions 获取当前管理员的权限列表
// @Summary      获取当前用户权限
// @Description  返回当前登录管理员拥有的权限码列表
// @Tags         Admin - Permissions
// @Security     BearerAuth
// @Success      200  {object}  model.APIResponse[[]model.Permission]
// @Router       /admin/permissions/me [get]
func (h *PermissionHandler) GetCurrentUserPermissions(c *gin.Context) {
	userIDVal, ok := c.Get(middleware.UserIDKey)
	if !ok {
		writeJSONError(c, http.StatusUnauthorized, "未登录")
		return
	}
	userID, _ := userIDVal.(uint64)
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
