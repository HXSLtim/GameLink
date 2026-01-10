package user

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	authService "gamelink/internal/service/auth"
	"gamelink/pkg/apierr"
	"gamelink/pkg/auth"
)

// RoleHandler 角色切换处理器
type RoleHandler struct {
	roleSvc *authService.RoleService
}

// NewRoleHandler 创建角色切换处理器
func NewRoleHandler(roleSvc *authService.RoleService) *RoleHandler {
	return &RoleHandler{
		roleSvc: roleSvc,
	}
}

// GetAvailableRoles 获取可用角色列表
// @Summary 获取可用角色列表
// @Description 获取当前用户可切换的角色列表
// @Tags 用户-角色
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} resp.Response{data=authService.AvailableRolesResponse}
// @Failure 401 {object} resp.Response
// @Router /user/roles [get]
func (h *RoleHandler) GetAvailableRoles(c *gin.Context) {
	userID, currentRole, err := extractUserContext(c)
	if err != nil {
		resp.Error(c, err)
		return
	}

	result, err := h.roleSvc.GetAvailableRoles(c.Request.Context(), userID, currentRole)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, result)
}

// SwitchRole 切换角色
// @Summary 切换角色
// @Description 切换当前用户的激活角色（用户/陪玩师）
// @Tags 用户-角色
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body authService.RoleSwitchRequest true "切换请求"
// @Success 200 {object} resp.Response{data=authService.RoleSwitchResponse}
// @Failure 400 {object} resp.Response
// @Failure 401 {object} resp.Response
// @Failure 403 {object} resp.Response
// @Router /user/switch-role [post]
func (h *RoleHandler) SwitchRole(c *gin.Context) {
	userID, _, err := extractUserContext(c)
	if err != nil {
		resp.Error(c, err)
		return
	}

	var req authService.RoleSwitchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("请求参数错误"))
		return
	}

	result, err := h.roleSvc.SwitchRole(c.Request.Context(), userID, req.TargetRole)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, result)
}

// extractUserContext 从上下文提取用户信息
func extractUserContext(c *gin.Context) (uint64, string, error) {
	// 从 JWT claims 中获取用户 ID
	claims, exists := c.Get("claims")
	if !exists {
		return 0, "", apierr.Unauthorized("未登录")
	}

	// 尝试解析为 CustomClaims（小程序 Token）
	if customClaims, ok := claims.(*auth.CustomClaims); ok {
		return customClaims.UserID, customClaims.CurrentRole, nil
	}

	// 尝试解析为标准 Claims（管理后台 Token）
	if stdClaims, ok := claims.(*auth.Claims); ok {
		return stdClaims.UserID, stdClaims.Role, nil
	}

	return 0, "", apierr.Unauthorized("无效的认证信息")
}

// RegisterRoutes 注册角色相关路由
func (h *RoleHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/roles", h.GetAvailableRoles)
	rg.POST("/switch-role", h.SwitchRole)
}
