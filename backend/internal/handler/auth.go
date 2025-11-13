package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/auth"
	"gamelink/internal/model"
	"gamelink/internal/service"
	authservice "gamelink/internal/service/auth"
)

// RegisterAuthRoutes registers authentication endpoints under the given router group.
// Routes:
// POST /auth/login    -> body {username, password}
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

// Login
// @Summary      登录
// @Description  用户名（邮箱或手机号）+ 密码登录，返回 JWT
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      loginRequest  true  "登录凭据"
// @Success      200      {object}  loginResponse
// @Failure      400      {object}  map[string]any
// @Failure      401      {object}  map[string]any
// @Router       /auth/login [post]
func loginHandler(c *gin.Context, svc *authservice.AuthService) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, ErrInvalidJSONPayload)
		return
	}
	resp, err := svc.Login(c.Request.Context(), authservice.LoginRequest{Username: req.Username, Password: req.Password})
	if err != nil {
		status := http.StatusUnauthorized
		switch {
		case errors.Is(err, service.ErrUserDisabled):
			status = http.StatusForbidden
		case errors.Is(err, service.ErrInvalidCredentials):
			// keep unauthorized status
		}
		respondJSON(c, status, model.APIResponse[any]{Success: false, Code: status, Message: err.Error()})
		return
	}
	respondJSON(c, http.StatusOK, model.APIResponse[loginResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    loginResponse{Token: resp.Token, ExpiresAt: resp.ExpiresAt, User: resp.User},
	})
}

// Register
// @Summary      注册
// @Description  邮箱或手机号 + 密码注册，默认角色为 user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      registerRequest  true  "注册信息"
// @Success      200      {object}  loginResponse
// @Failure      400      {object}  map[string]any
// @Router       /auth/register [post]
func registerHandler(c *gin.Context, svc *authservice.AuthService) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, ErrInvalidJSONPayload)
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
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(c, http.StatusOK, model.APIResponse[loginResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    loginResponse{Token: resp.Token, ExpiresAt: resp.ExpiresAt, User: resp.User},
	})
}

// Me
// @Summary      获取当前用户信息
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  loginResponse
// @Failure      401  {object}  map[string]any
// @Router       /auth/me [get]
func meHandler(c *gin.Context, svc *authservice.AuthService) {
	user, err := svc.Me(c.Request.Context(), c.GetHeader("Authorization"))
	if err != nil {
		status := http.StatusUnauthorized
		if err == service.ErrUserDisabled {
			status = http.StatusForbidden
		}
		respondError(c, status, err.Error())
		return
	}
	respondJSON(c, http.StatusOK, model.APIResponse[loginResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    loginResponse{Token: "", ExpiresAt: time.Time{}, User: *user},
	})
}

// Refresh
// @Summary      刷新 Token
// @Description  使用 Authorization: Bearer <token> 刷新 JWT
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  tokenPayload
// @Failure      401  {object}  map[string]any
// @Router       /auth/refresh [post]
func refreshHandler(c *gin.Context, svc *authservice.AuthService) {
	token, err := auth.ExtractTokenFromHeader(c.GetHeader("Authorization"))
	if err != nil {
		respondError(c, http.StatusUnauthorized, err.Error())
		return
	}
	newToken, err := svc.RefreshToken(c.Request.Context(), token)
	if err != nil {
		status := http.StatusUnauthorized
		if err == service.ErrUserDisabled {
			status = http.StatusForbidden
		}
		respondError(c, status, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[tokenPayload]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    tokenPayload{Token: newToken},
	})
}

// Logout
// @Summary      登出（前端丢弃 Token）
// @Tags         Auth
// @Security     BearerAuth
// @Success      200  {object}  map[string]any
// @Router       /auth/logout [post]
func logoutHandler(c *gin.Context) {
	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "logged out",
	})
}

// local helpers for uniform envelope
func respondJSON[T any](c *gin.Context, status int, payload model.APIResponse[T]) {
	c.JSON(status, payload)
}
func respondError(c *gin.Context, status int, msg string) {
	respondJSON(c, status, model.APIResponse[any]{Success: false, Code: status, Message: msg})
}
