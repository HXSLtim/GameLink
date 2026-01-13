package user

import (
	"unicode"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"gamelink/internal/handler/resp"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
)

// ChangePasswordHandler 修改密码处理器
type ChangePasswordHandler struct {
	userRepo repository.UserRepository
}

// NewChangePasswordHandler 创建修改密码处理器
func NewChangePasswordHandler(userRepo repository.UserRepository) *ChangePasswordHandler {
	return &ChangePasswordHandler{userRepo: userRepo}
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// ChangePasswordResponse 修改密码响应
type ChangePasswordResponse struct {
	Message string `json:"message" example:"密码修改成功"`
}

// ChangePassword 修改密码
// @Summary      修改密码
// @Description  修改当前用户的登录密码
// @Tags         User - Auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      ChangePasswordRequest  true  "密码信息"
// @Success      200      {object}  ChangePasswordResponse
// @Failure      400      {object}  apierr.APIError
// @Failure      401      {object}  apierr.APIError
// @Failure      500      {object}  apierr.APIError
// @Router       /auth/change-password [post]
func (h *ChangePasswordHandler) ChangePassword(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("无效的请求参数").WithDetails(err.Error()))
		return
	}

	ctx := c.Request.Context()

	// 获取用户
	user, err := h.userRepo.Get(ctx, userID)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取用户信息失败").WithDetails(err.Error()))
		return
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		resp.Error(c, apierr.BadRequest("原密码错误"))
		return
	}

	// 验证新密码强度
	if !isStrongPassword(req.NewPassword) {
		resp.Error(c, apierr.BadRequest("密码不符合安全要求：需包含大小写字母、数字和特殊符号"))
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		resp.Error(c, apierr.InternalError("密码加密失败"))
		return
	}

	// 更新密码
	if err := h.userRepo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		resp.Error(c, apierr.InternalError("更新密码失败").WithDetails(err.Error()))
		return
	}

	resp.Success[any](c, "密码修改成功", nil)
}

// isStrongPassword 检查密码强度
// 要求：8+ 字符，包含大小写字母、数字和特殊符号
func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

// RegisterChangePasswordRoutes 注册修改密码路由
// 注意：这个路由注册在 /auth 组下，不是 /user 组
func RegisterChangePasswordRoutes(rg *gin.RouterGroup, userRepo repository.UserRepository, authMiddleware gin.HandlerFunc) {
	h := NewChangePasswordHandler(userRepo)
	authGroup := rg.Group("/auth")
	authGroup.Use(authMiddleware)
	authGroup.POST("/change-password", h.ChangePassword)
}
