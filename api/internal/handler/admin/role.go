package admin

import (
	"errors"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	roleservice "gamelink/internal/service/admin"
	"gamelink/internal/service/audit"
	"gamelink/pkg/apierr"
)

// RoleHandler 角色管理处理器
type RoleHandler struct {
	roleSvc  *roleservice.RoleService
	auditSvc *audit.Service
}

// NewRoleHandler 创建角色处理器实例
func NewRoleHandler(roleSvc *roleservice.RoleService) *RoleHandler {
	return &RoleHandler{roleSvc: roleSvc}
}

// NewRoleHandlerWithAudit 创建带审计服务的角色处理器实例
func NewRoleHandlerWithAudit(roleSvc *roleservice.RoleService, auditSvc *audit.Service) *RoleHandler {
	return &RoleHandler{roleSvc: roleSvc, auditSvc: auditSvc}
}

// logRolePermissionChange 记录角色权限变更审计日志
func (h *RoleHandler) logRolePermissionChange(c *gin.Context, action model.AuditAction, roleID uint64, roleName string, beforeData, afterData any) {
	if h.auditSvc == nil {
		return
	}

	// 获取操作者信息
	operatorID, _ := middleware.GetUserID(c)
	operatorName, _ := middleware.GetUserRole(c)

	// 获取请求信息
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	requestID := c.GetHeader("X-Request-ID")

	h.auditSvc.LogRoleChange(
		c.Request.Context(),
		operatorID,
		operatorName,
		action,
		roleID,
		roleName,
		beforeData,
		afterData,
		ipAddress,
		userAgent,
		requestID,
	)
}

// logUserRoleChange 记录用户角色变更审计日志
func (h *RoleHandler) logUserRoleChange(c *gin.Context, action model.AuditAction, userID uint64, userName string, beforeData, afterData any) {
	if h.auditSvc == nil {
		return
	}

	// 获取操作者信息
	operatorID, _ := middleware.GetUserID(c)
	operatorName, _ := middleware.GetUserRole(c)

	// 获取请求信息
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	requestID := c.GetHeader("X-Request-ID")

	h.auditSvc.LogUserRoleChange(
		c.Request.Context(),
		operatorID,
		operatorName,
		action,
		userID,
		userName,
		beforeData,
		afterData,
		ipAddress,
		userAgent,
		requestID,
	)
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

// UpdateUserRolesRequest 更新用户角色请求（用于 PUT /users/:id/roles）
type UpdateUserRolesRequest struct {
	RoleIDs []uint64 `json:"roleIds" binding:"required"`
}

// BatchAssignRolesRequest 批量用户角色分配请求
type BatchAssignRolesRequest struct {
	UserIDs []uint64 `json:"userIds" binding:"required,min=1"`
	RoleIDs []uint64 `json:"roleIds" binding:"required"`
}

// BatchAssignRolesResult 批量分配结果
type BatchAssignRolesResult struct {
	SuccessCount int                     `json:"successCount"`
	FailedCount  int                     `json:"failedCount"`
	FailedUsers  []BatchAssignFailedUser `json:"failedUsers,omitempty"`
}

// BatchAssignFailedUser 批量分配失败的用户信息
type BatchAssignFailedUser struct {
	UserID uint64 `json:"userId"`
	Reason string `json:"reason"`
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
// @Param        keyword           query     string  false  "关键词搜索"
// @Param        isSystem          query     bool    false  "是否为系统角色"
// @Success      200               {object}  gin.H
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
		respondError(c, err)
		return
	}

	// 使用标准分页响应格式
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	respondList(c, roles, &model.Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
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
// @Success      200               {object}  gin.H
// @Failure      400               {object}  model.ErrorResponse
// @Failure      401               {object}  model.ErrorResponse
// @Failure      404               {object}  model.ErrorResponse
// @Failure      500               {object}  model.ErrorResponse
// @Router       /admin/roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	withPermissions := c.Query("with_permissions") == "true"

	var role *model.RoleModel
	var err error
	if withPermissions {
		role, err = h.roleSvc.GetRoleWithPermissions(c.Request.Context(), id)
	} else {
		role, err = h.roleSvc.GetRole(c.Request.Context(), id)
	}

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("角色不存在"))
		} else {
			respondError(c, err)
		}
		return
	}

	respondSuccess(c, role)
}

// CreateRole 创建角色
// @Summary      创建角色
// @Description  管理员创建新的系统角// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                true  "Bearer {token}"
// @Param        request        body      CreateRoleRequest      true  "创建角色请求"
// @Success      201            {object}  gin.H
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "参数验证失败")
		return
	}

	role := &model.RoleModel{
		Slug:        req.Slug,
		Name:        req.Name,
		Description: req.Description,
		IsSystem:    false,
	}

	if err := h.roleSvc.CreateRole(c.Request.Context(), role); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, role)
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
// @Success      200            {object}  gin.H
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      404            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req UpdateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "参数验证失败")
		return
	}

	role := &model.RoleModel{
		Base:        model.Base{ID: id},
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.roleSvc.UpdateRole(c.Request.Context(), role); err != nil {
		respondError(c, err)
		return
	}

	updatedRole, err := h.roleSvc.GetRole(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondUpdated(c, updatedRole)
}

// DeleteRole 删除角色
// @Summary      删除角色
// @Description  管理员删除指定的系统角色
// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      uint    true  "角色ID"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	if err := h.roleSvc.DeleteRole(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}

	respondMsg(c, "角色删除成功")
}

// GetRolePermissionIDs 获取角色的权限ID列表
// @Summary      获取角色的权限ID列表
// @Description  管理员获取指定角色的所有权限ID列表
// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      uint    true  "角色ID"
// @Success      200            {array}   uint64
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      404            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/roles/{id}/permissions [get]
func (h *RoleHandler) GetRolePermissionIDs(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	permissionIDs, err := h.roleSvc.GetRolePermissionIDs(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("角色不存在"))
		} else {
			respondError(c, err)
		}
		return
	}

	respondSuccess(c, ensureSlice(permissionIDs))
}

// AssignPermissions 为角色分配权限（批量替换）
// @Summary      批量分配角色权限
// @Description  管理员为指定角色批量分配权限（替换现有权限，事务保证原子性）
// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                      true  "Bearer {token}"
// @Param        id             path      uint                        true  "角色ID"
// @Param        request        body      AssignPermissionsRequest     true  "分配权限请求"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/roles/{id}/permissions/batch [put]
func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req AssignPermissionsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "参数验证失败")
		return
	}

	// 获取变更前的权限列表（用于审计日志）
	beforePermIDs, _ := h.roleSvc.GetRolePermissionIDs(c.Request.Context(), id)

	// 获取角色信息（用于审计日志）
	role, _ := h.roleSvc.GetRole(c.Request.Context(), id)
	roleName := ""
	if role != nil {
		roleName = role.Name
	}

	if err := h.roleSvc.AssignPermissionsToRole(c.Request.Context(), id, req.PermissionIDs); err != nil {
		respondError(c, err)
		return
	}

	// 失效相关用户缓存
	_ = h.roleSvc.InvalidateRolePermissionsAndPropagateToUsers(c.Request.Context(), id)

	// 记录审计日志
	h.logRolePermissionChange(c, model.AuditActionAssignPermission, id, roleName,
		map[string]any{"permissionIds": beforePermIDs},
		map[string]any{"permissionIds": req.PermissionIDs})

	respondMsg(c, "权限分配成功")
}

// AddPermissionRequest 单个添加权限请求
type AddPermissionRequest struct {
	PermissionID uint64 `json:"permissionId"`
}

// AddPermissionToRole 为角色添加单个权限
// @Summary      为角色添加单个权限
// @Description  管理员为指定角色添加单个权限
// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      uint    true  "角色ID"
// @Param        pid            path      uint    true  "权限ID"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      404            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/roles/{id}/permissions/{pid} [post]
func (h *RoleHandler) AddPermissionToRole(c *gin.Context) {
	roleID, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	permissionID, err := parseUintParam(c, "pid")
	if err != nil {
		respondBadRequest(c, "无效的权限ID")
		return
	}

	// 获取角色信息（用于审计日志）
	role, _ := h.roleSvc.GetRole(c.Request.Context(), roleID)
	roleName := ""
	if role != nil {
		roleName = role.Name
	}

	if err := h.roleSvc.AddPermissionsToRole(c.Request.Context(), roleID, []uint64{permissionID}); err != nil {
		respondError(c, err)
		return
	}

	// 失效相关用户缓存
	_ = h.roleSvc.InvalidateRolePermissionsAndPropagateToUsers(c.Request.Context(), roleID)

	// 记录审计日志
	h.logRolePermissionChange(c, model.AuditActionAddPermission, roleID, roleName,
		nil,
		map[string]any{"permissionId": permissionID})

	respondMsg(c, "权限添加成功")
}

// RemovePermissionFromRole 从角色移除单个权限
// @Summary      从角色移除单个权限
// @Description  管理员从指定角色移除单个权限
// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      uint    true  "角色ID"
// @Param        pid            path      uint    true  "权限ID"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      404            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/roles/{id}/permissions/{pid} [delete]
func (h *RoleHandler) RemovePermissionFromRole(c *gin.Context) {
	roleID, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	permissionID, err := parseUintParam(c, "pid")
	if err != nil {
		respondBadRequest(c, "无效的权限ID")
		return
	}

	// 获取角色信息（用于审计日志）
	role, _ := h.roleSvc.GetRole(c.Request.Context(), roleID)
	roleName := ""
	if role != nil {
		roleName = role.Name
	}

	if err := h.roleSvc.RemovePermissionsFromRole(c.Request.Context(), roleID, []uint64{permissionID}); err != nil {
		respondError(c, err)
		return
	}

	// 失效相关用户缓存
	_ = h.roleSvc.InvalidateRolePermissionsAndPropagateToUsers(c.Request.Context(), roleID)

	// 记录审计日志
	h.logRolePermissionChange(c, model.AuditActionRemovePermission, roleID, roleName,
		map[string]any{"permissionId": permissionID},
		nil)

	respondMsg(c, "权限移除成功")
}

// AssignRolesToUser 为用户分配角色
// @Summary      为用户分配角色
// @Description  管理员为指定用户分配多个角色（自动失效缓存并记录审计日志）
// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                    true  "Bearer {token}"
// @Param        request        body      AssignRolesToUserRequest    true  "分配角色请求"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/users/roles [post]
func (h *RoleHandler) AssignRolesToUser(c *gin.Context) {
	var req AssignRolesToUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "参数验证失败")
		return
	}

	// 获取变更前的角色列表（用于审计日志）
	beforeRoles, _ := h.roleSvc.ListRolesByUserID(c.Request.Context(), req.UserID)
	beforeRoleIDs := make([]uint64, len(beforeRoles))
	for i, r := range beforeRoles {
		beforeRoleIDs[i] = r.ID
	}

	// 分配角色（服务层会自动失效用户缓存）
	if err := h.roleSvc.AssignRolesToUser(c.Request.Context(), req.UserID, req.RoleIDs); err != nil {
		respondError(c, err)
		return
	}

	// 记录审计日志
	h.logUserRoleChange(c, model.AuditActionAssign, req.UserID, "",
		map[string]any{"roleIds": beforeRoleIDs},
		map[string]any{"roleIds": req.RoleIDs})

	respondMsg(c, "角色分配成功")
}

// GetUserRoles 获取用户的角色列表
// @Summary      获取用户的角色列表
// @Description  管理员获取指定用户的角色列表
// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id              path      uint    true  "用户ID"
// @Success      200            {array}   model.RoleModel
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/users/{id}/roles [get]
func (h *RoleHandler) GetUserRoles(c *gin.Context) {
	// 支持两种参数名：user_id 和 id
	var userID uint64
	var ok bool

	// 先检查哪个参数存在
	if c.Param("id") != "" {
		userID, ok = ParseIDAndRespond(c, "id")
	} else if c.Param("user_id") != "" {
		userID, ok = ParseIDAndRespond(c, "user_id")
	} else {
		respondBadRequest(c, "缺少用户ID参数")
		return
	}

	if !ok {
		return
	}

	roles, err := h.roleSvc.ListRolesByUserID(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, ensureSlice(roles))
}

// UpdateUserRoles 更新用户的角色（替换现有角色）
// @Summary      更新用户角色
// @Description  管理员更新指定用户的角色（替换现有角色，自动失效缓存并记录审计日志）
// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                    true  "Bearer {token}"
// @Param        id             path      uint                      true  "用户ID"
// @Param        request        body      UpdateUserRolesRequest    true  "更新角色请求"
// @Success      200            {array}   model.RoleModel
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/users/{id}/roles [put]
func (h *RoleHandler) UpdateUserRoles(c *gin.Context) {
	userID, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req UpdateUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "参数验证失败")
		return
	}

	// 获取变更前的角色列表（用于审计日志）
	beforeRoles, _ := h.roleSvc.ListRolesByUserID(c.Request.Context(), userID)
	beforeRoleIDs := make([]uint64, len(beforeRoles))
	for i, r := range beforeRoles {
		beforeRoleIDs[i] = r.ID
	}

	// 分配角色
	if err := h.roleSvc.AssignRolesToUser(c.Request.Context(), userID, req.RoleIDs); err != nil {
		respondError(c, err)
		return
	}

	// 获取更新后的角色列表
	afterRoles, err := h.roleSvc.ListRolesByUserID(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	// 记录审计日志
	h.logUserRoleChange(c, model.AuditActionAssign, userID, "",
		map[string]any{"roleIds": beforeRoleIDs},
		map[string]any{"roleIds": req.RoleIDs})

	respondUpdated(c, ensureSlice(afterRoles))
}

// BatchAssignRolesToUsers 批量为多个用户分配角色
// @Summary      批量分配用户角色
// @Description  管理员批量为多个用户分配相同的角色（自动失效缓存并记录审计日志）
// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                    true  "Bearer {token}"
// @Param        request        body      BatchAssignRolesRequest   true  "批量分配角色请求"
// @Success      200            {object}  BatchAssignRolesResult
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/users/roles/batch [put]
func (h *RoleHandler) BatchAssignRolesToUsers(c *gin.Context) {
	var req BatchAssignRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "参数验证失败")
		return
	}

	if len(req.UserIDs) == 0 {
		respondBadRequest(c, "用户ID列表不能为空")
		return
	}

	result := BatchAssignRolesResult{
		FailedUsers: make([]BatchAssignFailedUser, 0),
	}

	for _, userID := range req.UserIDs {
		// 获取变更前的角色列表（用于审计日志）
		beforeRoles, _ := h.roleSvc.ListRolesByUserID(c.Request.Context(), userID)
		beforeRoleIDs := make([]uint64, len(beforeRoles))
		for i, r := range beforeRoles {
			beforeRoleIDs[i] = r.ID
		}

		// 分配角色
		if err := h.roleSvc.AssignRolesToUser(c.Request.Context(), userID, req.RoleIDs); err != nil {
			result.FailedCount++
			result.FailedUsers = append(result.FailedUsers, BatchAssignFailedUser{
				UserID: userID,
				Reason: err.Error(),
			})
			continue
		}

		result.SuccessCount++

		// 记录审计日志
		h.logUserRoleChange(c, model.AuditActionBatchAssign, userID, "",
			map[string]any{"roleIds": beforeRoleIDs},
			map[string]any{"roleIds": req.RoleIDs})
	}

	respondSuccessWithMsg(c, "批量角色分配完成", result)
}

// ============================================================================
// 批量角色操作相关请求和响应结构
// ============================================================================

// BatchDeleteRolesRequest 批量删除角色请求
type BatchDeleteRolesRequest struct {
	IDs []uint64 `json:"ids" binding:"required,min=1"`
}

// BatchAssignPermissionsToRolesRequest 批量为角色分配权限请求
type BatchAssignPermissionsToRolesRequest struct {
	Assignments []struct {
		RoleID        uint64   `json:"roleId" binding:"required"`
		PermissionIDs []uint64 `json:"permissionIds"`
	} `json:"assignments" binding:"required,min=1"`
}

// ============================================================================
// 批量角色操作处理器方法
// ============================================================================

// BatchDeleteRoles 批量删除角色
// @Summary      批量删除角色
// @Description  批量删除多个角色（系统角色不可删除），返回成功和失败的数量
// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                       true  "Bearer {token}"
// @Param        request        body      BatchDeleteRolesRequest       true  "批量删除角色请求"
// @Success      200            {object}  roleservice.RoleBatchDeleteResult
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/roles/batch [delete]
func (h *RoleHandler) BatchDeleteRoles(c *gin.Context) {
	var req BatchDeleteRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "参数验证失败: "+err.Error())
		return
	}

	result, err := h.roleSvc.BatchDeleteRoles(c.Request.Context(), req.IDs)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccessWithMsg(c, "批量角色删除完成", result)
}

// BatchAssignPermissionsToRoles 批量为多个角色分配权限
// @Summary      批量为角色分配权限
// @Description  批量为多个角色分配权限（替换现有权限），返回成功和失败的数量
// @Tags         Admin - Roles
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                                true  "Bearer {token}"
// @Param        request        body      BatchAssignPermissionsToRolesRequest    true  "批量分配权限请求"
// @Success      200            {object}  roleservice.RoleBatchPermissionsResult
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/roles/batch/permissions [put]
func (h *RoleHandler) BatchAssignPermissionsToRoles(c *gin.Context) {
	var req BatchAssignPermissionsToRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "参数验证失败: "+err.Error())
		return
	}

	// 转换请求格式
	assignments := make([]roleservice.RolePermissionAssignment, len(req.Assignments))
	for i, a := range req.Assignments {
		assignments[i] = roleservice.RolePermissionAssignment{
			RoleID:        a.RoleID,
			PermissionIDs: a.PermissionIDs,
		}
	}

	result, err := h.roleSvc.BatchAssignPermissionsToRoles(c.Request.Context(), assignments)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccessWithMsg(c, "批量权限分配完成", result)
}
