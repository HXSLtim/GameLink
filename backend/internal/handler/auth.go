package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/auth"
	"gamelink/internal/model"
	"gamelink/internal/service"
)

// RegisterAuthRoutes registers authentication endpoints under the given router group.
// Routes:
// POST /auth/login    -> body {username, password}
// POST /auth/refresh  -> Authorization: Bearer <token>
// POST /auth/logout   -> stateless logout, client discards token
func RegisterAuthRoutes(router gin.IRoutes, svc *service.AuthService) {
	router.POST("/auth/login", func(c *gin.Context) { loginHandler(c, svc) })
	router.POST("/auth/refresh", func(c *gin.Context) { refreshHandler(c, svc) })
	router.POST("/auth/logout", logoutHandler)
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

// tokenPayload is used for Swagger response schema.
// tokenPayload is used for Swagger response schema.
// nolint:unused
type tokenPayload struct {
	Token string `json:"token"`
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
func loginHandler(c *gin.Context, svc *service.AuthService) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid json payload")
		return
	}
	resp, err := svc.Login(c.Request.Context(), service.LoginRequest{Username: req.Username, Password: req.Password})
	if err != nil {
		respondJSON(c, http.StatusUnauthorized, model.APIResponse[any]{Success: false, Code: http.StatusUnauthorized, Message: err.Error()})
		return
	}
	respondJSON(c, http.StatusOK, model.APIResponse[loginResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    loginResponse{Token: resp.Token, ExpiresAt: resp.ExpiresAt, User: resp.User},
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
func refreshHandler(c *gin.Context, svc *service.AuthService) {
	token, err := auth.ExtractTokenFromHeader(c.GetHeader("Authorization"))
	if err != nil {
		respondError(c, http.StatusUnauthorized, err.Error())
		return
	}
	newToken, err := svc.RefreshToken(c.Request.Context(), token)
	if err != nil {
		respondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[map[string]any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    map[string]any{"token": newToken},
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
