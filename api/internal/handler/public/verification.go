package public

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	_ "gamelink/internal/model" // for swagger
	"gamelink/internal/service/sms"
	"gamelink/internal/service/verification"
	"gamelink/pkg/apierr"
)

// VerificationHandler 验证码处理器
type VerificationHandler struct {
	svc *verification.Service
}

// NewVerificationHandler 创建验证码处理器
func NewVerificationHandler(svc *verification.Service) *VerificationHandler {
	return &VerificationHandler{svc: svc}
}

// SendCodeRequest 发送验证码请求
type SendCodeRequest struct {
	Target string `json:"target" binding:"required"`                 // 手机号或邮箱
	Type   string `json:"type" binding:"required,oneof=phone email"` // phone 或 email
}

// VerifyCodeRequest 验证验证码请求
type VerifyCodeRequest struct {
	Target string `json:"target" binding:"required"`                 // 手机号或邮箱
	Code   string `json:"code" binding:"required,len=6"`             // 6位验证码
	Type   string `json:"type" binding:"required,oneof=phone email"` // phone 或 email
}

// SendCodeResponse 发送验证码响应
type SendCodeResponse struct {
	Message    string `json:"message"`              // 提示消息
	Code       string `json:"code,omitempty"`       // 验证码（仅开发环境）
	MasterCode string `json:"masterCode,omitempty"` // 万能验证码（仅开发环境）
}

// VerifyCodeResponse 验证验证码响应
type VerifyCodeResponse struct {
	Verified bool `json:"verified"` // 是否验证成功
}

// MasterCodeResponse 万能验证码响应
type MasterCodeResponse struct {
	MasterCode string `json:"masterCode"` // 万能验证码
	Hint       string `json:"hint"`       // 提示信息
}

// SendCode 发送验证码
// @Summary 发送验证码
// @Description 发送手机或邮箱验证码
// @Tags 公共-验证码
// @Accept json
// @Produce json
// @Param request body SendCodeRequest true "发送验证码请求"
// @Success 200 {object} model.APIResponse[SendCodeResponse] "发送成功"
// @Failure 400 {object} model.ErrorResponse "请求参数错误"
// @Failure 429 {object} model.ErrorResponse "请求过于频繁"
// @Router /public/verification/send [post]
func (h *VerificationHandler) SendCode(c *gin.Context) {
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("invalid request: "+err.Error()))
		return
	}

	codeType := verification.CodeType(req.Type)
	err := h.svc.SendCode(c.Request.Context(), req.Target, codeType)
	if err != nil {
		switch err {
		case verification.ErrRateLimited:
			resp.Error(c, apierr.TooManyRequests("请等待60秒后再试"))
		case sms.ErrSMSDisabled:
			resp.Error(c, apierr.InternalError("短信服务未配置"))
		default:
			resp.Error(c, apierr.InternalError("发送验证码失败: "+err.Error()))
		}
		return
	}

	response := gin.H{
		"message": "验证码已发送",
	}

	if h.svc.GetMasterCode() != "" {
		response["masterCode"] = h.svc.GetMasterCode()
	}

	resp.OK(c, response)
}

// VerifyCode 验证验证码
// @Summary 验证验证码
// @Description 验证手机或邮箱验证码
// @Tags 公共-验证码
// @Accept json
// @Produce json
// @Param request body VerifyCodeRequest true "验证验证码请求"
// @Success 200 {object} model.APIResponse[VerifyCodeResponse] "验证成功"
// @Failure 400 {object} model.ErrorResponse "验证码错误"
// @Failure 429 {object} model.ErrorResponse "尝试次数过多"
// @Router /public/verification/verify [post]
func (h *VerificationHandler) VerifyCode(c *gin.Context) {
	var req VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("invalid request: "+err.Error()))
		return
	}

	codeType := verification.CodeType(req.Type)
	err := h.svc.VerifyCode(c.Request.Context(), req.Target, req.Code, codeType)
	if err != nil {
		switch err {
		case verification.ErrCodeNotFound:
			resp.Error(c, apierr.BadRequest("验证码不存在或已过期"))
		case verification.ErrCodeInvalid:
			resp.Error(c, apierr.BadRequest("验证码错误"))
		case verification.ErrTooManyAttempts:
			resp.Error(c, apierr.TooManyRequests("尝试次数过多，请重新获取验证码"))
		default:
			resp.Error(c, apierr.InternalError("验证失败"))
		}
		return
	}

	resp.OK(c, gin.H{
		"verified": true,
	})
}

// GetMasterCode 获取万能验证码（仅开发环境）
// @Summary 获取万能验证码
// @Description 获取万能验证码（仅开发/测试环境可用）
// @Tags 公共-验证码
// @Produce json
// @Success 200 {object} model.APIResponse[MasterCodeResponse] "万能验证码"
// @Failure 403 {object} model.ErrorResponse "生产环境不可用"
// @Router /public/verification/master-code [get]
func (h *VerificationHandler) GetMasterCode(c *gin.Context) {
	masterCode := h.svc.GetMasterCode()
	if masterCode == "" {
		resp.Error(c, apierr.Forbidden("生产环境不可用"))
		return
	}

	resp.OK(c, gin.H{
		"masterCode": masterCode,
		"hint":       "万能验证码仅在开发/测试环境有效",
	})
}

// RegisterRoutes 注册验证码路由
func (h *VerificationHandler) RegisterRoutes(group *gin.RouterGroup) {
	verificationGroup := group.Group("/verification")
	{
		verificationGroup.POST("/send", h.SendCode)
		verificationGroup.POST("/verify", h.VerifyCode)
		verificationGroup.GET("/master-code", h.GetMasterCode)
	}
}
