package handler

import (
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/apierr"
	"gamelink/internal/auth"
	"gamelink/internal/model"
	"gamelink/internal/service"
	authservice "gamelink/internal/service/auth"
)

// RegisterAuthRoutes registers authentication endpoints under the given router group.
// Routes:
// POST /auth/login    -> body {username, password}
// POST /auth/register  -> body {phone, email, password, name}
// POST /auth/refresh  -> Authorization: Bearer <token>
// POST /auth/logout   -> stateless logout, client discards token
// GET  /auth/me       -> return current user info (JWT required)
func RegisterAuthRoutes(router gin.IRouter, svc *authservice.AuthService) {
	auth := router.Group("/auth")
	auth.POST("/login", func(c *gin.Context) { loginHandler(c, svc) })
	auth.POST("/register", func(c *gin.Context) { registerHandler(c, svc) })
	auth.POST("/refresh", func(c *gin.Context) { refreshHandler(c, svc) })
	auth.POST("/logout", logoutHandler)

	auth.GET("/me", func(c *gin.Context) { meHandler(c, svc) })
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token     string     `json:"token"`
	ExpiresAt time.Time  `json:"expires_at"`
	User      model.User `json:"user"`
}

type tokenPayload struct {
	Token string `json:"token"`
}

type registerRequest struct {
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

// Login handles user login
// @Summary      登录
// @Description  用户名（邮箱或手机号）+ 密码登录，返回 JWT
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      loginRequest  true  "登录凭据"
// @Success      200      {object}  loginResponse
// @Failure      400      {object}  apierr.APIError
// @Failure      401      {object}  apierr.APIError
// @Router       /auth/login [post]
func loginHandler(c *gin.Context, svc *authservice.AuthService) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "无效的请求格式: "+err.Error())
		return
	}

	resp, err := svc.Login(c.Request.Context(), authservice.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			RespondAPIError(c, apierr.Unauthorized("用户名或密码错误"))
		case service.ErrUserDisabled:
			RespondAPIError(c, apierr.Forbidden("账号已被禁用"))
		default:
			RespondAPIError(c, apierr.Unauthorized("登录失败: "+err.Error()))
		}
		return
	}

	RespondSuccess(c, "登录成功", loginResponse{
		Token:     resp.Token,
		ExpiresAt: resp.ExpiresAt,
		User:      resp.User,
	})
}

// Register handles user registration
// @Summary      注册
// @Description  邮箱或手机号 + 密码注册，默认角色为 user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      registerRequest  true  "注册信息"
// @Success      200      {object}  loginResponse
// @Failure      400      {object}  apierr.APIError
// @Router       /auth/register [post]
func registerHandler(c *gin.Context, svc *authservice.AuthService) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "无效的请求格式: "+err.Error())
		return
	}

	resp, err := svc.Register(c.Request.Context(), authservice.RegisterRequest{
		Phone:    req.Phone,
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Role:     model.RoleUser,
	})
	if err != nil {
		RespondAPIError(c, apierr.BadRequest("注册失败: "+err.Error()))
		return
	}

	RespondSuccess(c, "登录成功", loginResponse{
		Token:     resp.Token,
		ExpiresAt: resp.ExpiresAt,
		User:      resp.User,
	})
}

// Me returns current user info
// @Summary      获取当前用户信息
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  loginResponse
// @Failure      401  {object}  apierr.APIError
// @Router       /auth/me [get]
func meHandler(c *gin.Context, svc *authservice.AuthService) {
	user, err := svc.Me(c.Request.Context(), c.GetHeader("Authorization"))
	if err != nil {
		switch err {
		case service.ErrUserDisabled:
			RespondAPIError(c, apierr.Forbidden("账号已被禁用"))
		default:
			RespondAPIError(c, apierr.Unauthorized("认证失败: "+err.Error()))
		}
		return
	}

	RespondSuccess(c, "登录成功", loginResponse{
		Token:     "",
		ExpiresAt: time.Time{},
		User:      *user,
	})
}

// Refresh refreshes JWT token
// @Summary      刷新 Token
// @Description  使用 Authorization: Bearer <token> 刷新 JWT
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  tokenPayload
// @Failure      401  {object}  apierr.APIError
// @Router       /auth/refresh [post]
func refreshHandler(c *gin.Context, svc *authservice.AuthService) {
	token, err := auth.ExtractTokenFromHeader(c.GetHeader("Authorization"))
	if err != nil {
		RespondAPIError(c, apierr.Unauthorized("无效的Token格式: "+err.Error()))
		return
	}

	newToken, err := svc.RefreshToken(c.Request.Context(), token)
	if err != nil {
		switch err {
		case service.ErrUserDisabled:
			RespondAPIError(c, apierr.Forbidden("账号已被禁用"))
		default:
			RespondAPIError(c, apierr.Unauthorized("刷新Token失败: "+err.Error()))
		}
		return
	}

	RespondSuccess(c, "刷新成功", tokenPayload{Token: newToken})
}

// Logout handles user logout (stateless, client discards token)
// @Summary      登出（前端丢弃 Token）
// @Tags         Auth
// @Security     BearerAuth
// @Success      200  {object}  apierr.APIError
// @Router       /auth/logout [post]
func logoutHandler(c *gin.Context) {
	RespondSuccess(c, "登出成功", gin.H{"message": "logged out"})
}
