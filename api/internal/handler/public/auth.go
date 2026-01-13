// Package public provides public API handlers that don't require authentication.
package public

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	authService "gamelink/internal/service/auth"
	"gamelink/pkg/apierr"
)

// AuthHandler 公共认证处理器
type AuthHandler struct {
	wechatSvc *authService.WeChatAuthService
}

// NewAuthHandler 创建公共认证处理器
func NewAuthHandler(wechatSvc *authService.WeChatAuthService) *AuthHandler {
	return &AuthHandler{
		wechatSvc: wechatSvc,
	}
}

// WeChatLogin 微信小程序登录
// @Summary 微信小程序登录
// @Description 使用微信登录凭证进行登录，返回访问令牌
// @Tags 公共-认证
// @Accept json
// @Produce json
// @Param request body authService.WeChatLoginRequest true "登录请求"
// @Success 200 {object} authService.WeChatLoginResponse
// @Failure 400 {object}  apierr.APIError
// @Failure 500 {object}  apierr.APIError
// @Router /public/auth/wechat/login [post]
func (h *AuthHandler) WeChatLogin(c *gin.Context) {
	var req authService.WeChatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("请求参数错误"))
		return
	}

	result, err := h.wechatSvc.WeChatLogin(c.Request.Context(), req)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, result)
}

// RefreshToken 刷新访问令牌
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 公共-认证
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "刷新请求"
// @Success 200 {object} authService.WeChatLoginResponse
// @Failure 400 {object}  apierr.APIError
// @Failure 401 {object}  apierr.APIError
// @Router /public/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("请求参数错误"))
		return
	}

	result, err := h.wechatSvc.RefreshAccessToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, result)
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RegisterRoutes 注册公共认证路由
func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/wechat/login", h.WeChatLogin)
		auth.POST("/refresh", h.RefreshToken)
	}
}

// HealthCheck 健康检查（公共端点）
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
