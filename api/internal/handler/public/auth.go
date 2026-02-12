// Package public provides public API handlers that don't require authentication.
package public

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	authService "gamelink/internal/service/auth"
	"gamelink/internal/service/verification"
	"gamelink/pkg/apierr"
)

// AuthHandler 公共认证处理器
type AuthHandler struct {
	wechatSvc       *authService.WeChatAuthService
	authSvc         *authService.AuthService
	verificationSvc *verification.Service
}

// NewAuthHandler 创建公共认证处理器
func NewAuthHandler(
	wechatSvc *authService.WeChatAuthService,
	authSvc *authService.AuthService,
	verificationSvc *verification.Service,
) *AuthHandler {
	return &AuthHandler{
		wechatSvc:       wechatSvc,
		authSvc:         authSvc,
		verificationSvc: verificationSvc,
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

// PhoneLoginRequest 手机号验证码登录请求
type PhoneLoginRequest struct {
	Phone        string `json:"phone" binding:"required"`
	Code         string `json:"code" binding:"required,len=6"`
	Nickname     string `json:"nickname,omitempty"`
	ReferralCode string `json:"referralCode,omitempty"`
}

// PhoneLogin 手机号验证码登录（新用户自动注册）
// @Summary 手机号验证码登录
// @Description 使用手机号+验证码登录；若用户不存在则自动创建账号并登录
// @Tags 公共-认证
// @Accept json
// @Produce json
// @Param request body PhoneLoginRequest true "登录请求"
// @Success 200 {object} authService.LoginResponse
// @Failure 400 {object} apierr.APIError
// @Failure 429 {object} apierr.APIError
// @Failure 500 {object} apierr.APIError
// @Router /public/auth/phone/login [post]
func (h *AuthHandler) PhoneLogin(c *gin.Context) {
	if h.authSvc == nil || h.verificationSvc == nil {
		resp.Error(c, apierr.InternalError("手机号登录服务未初始化"))
		return
	}

	var req PhoneLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("请求参数错误").WithDetails(err.Error()))
		return
	}

	if err := h.verificationSvc.VerifyCode(c.Request.Context(), req.Phone, req.Code, verification.CodeTypePhone); err != nil {
		switch err {
		case verification.ErrCodeNotFound:
			resp.Error(c, apierr.BadRequest("验证码不存在或已过期"))
		case verification.ErrCodeInvalid:
			resp.Error(c, apierr.BadRequest("验证码错误"))
		case verification.ErrTooManyAttempts:
			resp.Error(c, apierr.TooManyRequests("尝试次数过多，请重新获取验证码"))
		default:
			resp.Error(c, apierr.InternalError("验证码校验失败").WithDetails(err.Error()))
		}
		return
	}

	result, err := h.authSvc.LoginOrRegisterByPhone(c.Request.Context(), authService.PhoneCodeLoginRequest{
		Phone:    req.Phone,
		Nickname: req.Nickname,
	})
	if err != nil {
		if apiErr, ok := err.(*apierr.APIError); ok && apiErr.Code >= 400 && apiErr.Code < 500 {
			resp.Error(c, apiErr)
			return
		}
		if apierr.IsBadRequest(err) || apierr.IsForbidden(err) || apierr.IsUnauthorized(err) || apierr.IsNotFound(err) {
			resp.Error(c, err)
			return
		}
		resp.Error(c, apierr.InternalError("手机号登录失败").WithDetails(err.Error()))
		return
	}

	resp.OK(c, result)
}

// RegisterRoutes 注册公共认证路由
func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/wechat/login", h.WeChatLogin)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/phone/login", h.PhoneLogin)
	}
}

// HealthCheck 健康检查（公共端点）
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
