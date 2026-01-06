package admin

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"gamelink/internal/model"
	adminservice "gamelink/internal/service/admin"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type BatchSyncHandler struct {
	menuSvc *adminservice.MenuService
	permSvc *adminservice.PermissionService
	roleSvc *adminservice.RoleService
	db      *gorm.DB
}

func NewBatchSyncHandler(menuSvc *adminservice.MenuService, permSvc *adminservice.PermissionService, roleSvc *adminservice.RoleService, db *gorm.DB) *BatchSyncHandler {
	return &BatchSyncHandler{menuSvc: menuSvc, permSvc: permSvc, roleSvc: roleSvc, db: db}
}

type MenuSyncItem struct {
	Name        string         `json:"name"`
	Path        string         `json:"path"`
	Component   string         `json:"component"`
	Icon        string         `json:"icon,omitempty"`
	Order       int            `json:"order"`
	Hidden      bool           `json:"hidden,omitempty"`
	Visible     *bool          `json:"visible,omitempty"`
	Permission  string         `json:"permission,omitempty"`
	Redirect    string         `json:"redirect,omitempty"`
	Description string         `json:"description,omitempty"`
	Children    []MenuSyncItem `json:"children,omitempty"`
}

type PermissionSyncItem struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Group       string `json:"group"`
	Description string `json:"description"`
}

type BatchSyncRequest struct {
	Menus                       []MenuSyncItem       `json:"menus"`
	Permissions                 []PermissionSyncItem `json:"permissions"`
	AssignSuperAdminPermissions bool                 `json:"assignSuperAdminPermissions"`
}

type BatchSyncResult struct {
	Success bool     `json:"success"`
	Created int      `json:"created"`
	Updated int      `json:"updated"`
	Skipped int      `json:"skipped"`
	Errors  []string `json:"errors"`
}

type SuperAdminAssignResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type BatchSyncResponse struct {
	Success          bool                    `json:"success"`
	MenuSync         *BatchSyncResult        `json:"menuSync,omitempty"`
	PermissionSync   *BatchSyncResult        `json:"permissionSync,omitempty"`
	SuperAdminAssign *SuperAdminAssignResult `json:"superAdminAssign,omitempty"`
	Errors           []string                `json:"errors"`
}

func (h *BatchSyncHandler) BatchSync(c *gin.Context) {
	var req BatchSyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}
	ctx := c.Request.Context()
	response := BatchSyncResponse{Success: true, Errors: []string{}}

	if len(req.Permissions) > 0 {
		permResult := h.syncPermissions(ctx, req.Permissions)
		response.PermissionSync = permResult
		if !permResult.Success {
			response.Errors = append(response.Errors, permResult.Errors...)
		}
	}

	if len(req.Menus) > 0 {
		menuResult := h.syncMenus(ctx, req.Menus, nil)
		response.MenuSync = menuResult
		if !menuResult.Success {
			response.Errors = append(response.Errors, menuResult.Errors...)
		}
	}

	if req.AssignSuperAdminPermissions {
		assignResult := h.assignSuperAdminPermissions(ctx)
		response.SuperAdminAssign = assignResult
		if !assignResult.Success {
			response.Errors = append(response.Errors, assignResult.Message)
		}
	}

	if len(response.Errors) > 0 {
		response.Success = false
	}

	// Track initialization state in database if sync was successful
	if response.Success && h.db != nil {
		if err := h.recordInitState(ctx, req, c); err != nil {
			// Log the error but don't fail the sync
			fmt.Printf("warning: failed to record init state: %v\n", err)
		}
	}

	respondSuccessWithMsg(c, "sync complete", response)
}

// recordInitState records the initialization state in the database
func (h *BatchSyncHandler) recordInitState(ctx context.Context, req BatchSyncRequest, c *gin.Context) error {
	// Calculate version hashes
	menuVersion := h.calculateMenuVersion(req.Menus)
	permVersion := h.calculatePermVersion(req.Permissions)

	// Get user ID from context (set by auth middleware)
	var userID uint64 = 0
	if userIDVal, exists := c.Get("userID"); exists {
		userID = userIDVal.(uint64)
	}

	// Get IP address
	ip := c.ClientIP()

	// Count total items
	menuCount := countMenus(req.Menus)
	permCount := len(req.Permissions)

	// Create or update system state
	var existing model.SystemState
	err := h.db.Where("key = ?", model.SystemStateKeyAdminInit).First(&existing).Error

	initData := &model.SystemInitData{
		MenuCount:       menuCount,
		PermissionCount: permCount,
		MenuVersion:     menuVersion,
		PermVersion:     permVersion,
	}

	if err == nil {
		// Update existing record
		existing.Version = menuVersion + ":" + permVersion
		existing.LastSyncAt = time.Now()
		existing.SyncedBy = userID
		existing.SyncedByIP = ip
		if err := existing.SetInitData(initData); err != nil {
			return err
		}
		return h.db.Save(&existing).Error
	} else if err == gorm.ErrRecordNotFound {
		// Create new record
		newState := model.NewAdminInitState(menuCount, permCount, menuVersion, permVersion, userID, ip)
		return h.db.Create(newState).Error
	}
	return err
}

// calculateMenuVersion calculates a hash version of menus
func (h *BatchSyncHandler) calculateMenuVersion(menus []MenuSyncItem) string {
	data, _ := json.Marshal(menus)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])[:16]
}

// calculatePermVersion calculates a hash version of permissions
func (h *BatchSyncHandler) calculatePermVersion(perms []PermissionSyncItem) string {
	data, _ := json.Marshal(perms)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])[:16]
}

// countMenus counts total number of menus including children
func countMenus(menus []MenuSyncItem) int {
	count := len(menus)
	for _, m := range menus {
		count += countMenus(m.Children)
	}
	return count
}

func (h *BatchSyncHandler) syncPermissions(ctx context.Context, permissions []PermissionSyncItem) *BatchSyncResult {
	result := &BatchSyncResult{Success: true, Errors: []string{}}
	existingPerms, err := h.permSvc.ListPermissions(ctx)
	if err != nil {
		result.Success = false
		result.Errors = append(result.Errors, "get permissions failed: "+err.Error())
		return result
	}
	existingMap := make(map[string]*model.Permission)
	for i := range existingPerms {
		existingMap[existingPerms[i].Code] = &existingPerms[i]
	}
	for _, perm := range permissions {
		existing, exists := existingMap[perm.Code]
		permData := &model.Permission{
			Code: perm.Code, Method: model.HTTPMethod(perm.Method), Path: perm.Path, Group: perm.Group, Description: perm.Description,
		}
		if exists {
			if string(existing.Method) != perm.Method || existing.Path != perm.Path || existing.Description != perm.Description {
				permData.ID = existing.ID
				if err := h.permSvc.UpdatePermission(ctx, permData); err != nil {
					result.Errors = append(result.Errors, "update permission failed "+perm.Code+": "+err.Error())
				} else {
					result.Updated++
				}
			} else {
				result.Skipped++
			}
		} else {
			if err := h.permSvc.CreatePermission(ctx, permData); err != nil {
				result.Errors = append(result.Errors, "create permission failed "+perm.Code+": "+err.Error())
			} else {
				result.Created++
			}
		}
	}
	if len(result.Errors) > 0 {
		result.Success = false
	}
	return result
}

func (h *BatchSyncHandler) syncMenus(ctx context.Context, menus []MenuSyncItem, parentID *uint64) *BatchSyncResult {
	result := &BatchSyncResult{Success: true, Errors: []string{}}
	existingMenus, err := h.menuSvc.List(ctx, parentID)
	if err != nil {
		result.Success = false
		result.Errors = append(result.Errors, "get menus failed: "+err.Error())
		return result
	}
	existingMap := make(map[string]*model.Menu)
	for i := range existingMenus {
		existingMap[existingMenus[i].Path] = &existingMenus[i]
	}
	for _, menu := range menus {
		existing, exists := existingMap[menu.Path]
		hidden := menu.Hidden
		if menu.Visible != nil {
			hidden = !*menu.Visible
		}
		menuData := &model.Menu{
			Name: menu.Name, Path: menu.Path, Component: menu.Component, Icon: menu.Icon,
			ParentID: parentID, Order: menu.Order, Hidden: hidden, Permission: menu.Permission,
			Redirect: menu.Redirect, Description: menu.Description,
		}
		var currentMenuID uint64
		if exists {
			if existing.Name != menu.Name || existing.Component != menu.Component || existing.Order != menu.Order || existing.Icon != menu.Icon {
				menuData.ID = existing.ID
				if err := h.menuSvc.Update(ctx, menuData); err != nil {
					result.Errors = append(result.Errors, "update menu failed "+menu.Path+": "+err.Error())
				} else {
					result.Updated++
				}
			} else {
				result.Skipped++
			}
			currentMenuID = existing.ID
		} else {
			if err := h.menuSvc.Create(ctx, menuData); err != nil {
				result.Errors = append(result.Errors, "create menu failed "+menu.Path+": "+err.Error())
				continue
			}
			result.Created++
			currentMenuID = menuData.ID
		}
		if len(menu.Children) > 0 {
			childResult := h.syncMenus(ctx, menu.Children, &currentMenuID)
			result.Created += childResult.Created
			result.Updated += childResult.Updated
			result.Skipped += childResult.Skipped
			result.Errors = append(result.Errors, childResult.Errors...)
		}
	}
	if len(result.Errors) > 0 {
		result.Success = false
	}
	return result
}

func (h *BatchSyncHandler) assignSuperAdminPermissions(ctx context.Context) *SuperAdminAssignResult {
	result := &SuperAdminAssignResult{}
	roles, err := h.roleSvc.ListRoles(ctx)
	if err != nil {
		result.Success = false
		result.Message = "get roles failed: " + err.Error()
		return result
	}

	// 找到 super_admin 和 admin 角色
	var superAdminRole, adminRole *model.RoleModel
	for i := range roles {
		if roles[i].Slug == string(model.RoleSlugSuperAdmin) || roles[i].Name == "super_admin" {
			superAdminRole = &roles[i]
		}
		if roles[i].Slug == string(model.RoleSlugAdmin) || roles[i].Name == "admin" {
			adminRole = &roles[i]
		}
	}

	if superAdminRole == nil {
		result.Success = false
		result.Message = "super admin role not found"
		return result
	}

	permissions, err := h.permSvc.ListPermissions(ctx)
	if err != nil {
		result.Success = false
		result.Message = "get permissions failed: " + err.Error()
		return result
	}
	if len(permissions) == 0 {
		result.Success = false
		result.Message = "no permissions to assign"
		return result
	}

	permissionIDs := make([]uint64, len(permissions))
	for i, p := range permissions {
		permissionIDs[i] = p.ID
	}

	// 为 super_admin 分配所有权限
	if err := h.roleSvc.AssignPermissionsToRole(ctx, superAdminRole.ID, permissionIDs); err != nil {
		result.Success = false
		result.Message = "assign permissions to super_admin failed: " + err.Error()
		return result
	}

	// 为 admin 角色也分配所有权限
	if adminRole != nil {
		if err := h.roleSvc.AssignPermissionsToRole(ctx, adminRole.ID, permissionIDs); err != nil {
			// admin 角色分配失败不影响整体结果，只记录警告
			result.Success = true
			result.Message = fmt.Sprintf("assigned %d permissions to super_admin, but admin role failed: %v", len(permissions), err)
			return result
		}
		result.Success = true
		result.Message = fmt.Sprintf("assigned %d permissions to super_admin and admin roles", len(permissions))
		return result
	}

	result.Success = true
	result.Message = fmt.Sprintf("assigned %d permissions to super_admin", len(permissions))
	return result
}

// RequireSuperAdmin 中间件：检查用户是否是超级管理员
// 用于保护系统初始化等敏感操作
func RequireSuperAdmin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID
		userIDVal, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{
				"success": false,
				"code":    401,
				"message": "未认证",
			})
			c.Abort()
			return
		}
		userID := userIDVal.(uint64)

		// 查询用户的 RBAC 角色
		var userRoles []model.UserRole
		if err := db.Where("user_id = ?", userID).
			Preload("Role").
			Find(&userRoles).Error; err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"code":    500,
				"message": "无法验证用户权限",
			})
			c.Abort()
			return
		}

		// 检查是否有超级管理员角色
		isSuperAdmin := false
		for _, ur := range userRoles {
			if ur.Role.IsSuperAdmin() {
				isSuperAdmin = true
				break
			}
		}

		if !isSuperAdmin {
			c.JSON(403, gin.H{
				"success": false,
				"code":    403,
				"message": "只有超级管理员才能执行此操作",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
