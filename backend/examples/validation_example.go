//go:build ignore
// +build ignore

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler"
	"gamelink/internal/handler/middleware"
	"gamelink/internal/service"
)

// 这是一个示例，展示如何使用输入验证中间件
// 在实际项目中，你需要修改现有的handler来使用这些验证

// ExampleValidationHandler 展示验证中间件的使用
type ExampleValidationHandler struct {
	svc *service.AdminService
}

func NewExampleValidationHandler(svc *service.AdminService) *ExampleValidationHandler {
	return &ExampleValidationHandler{svc: svc}
}

// RegisterRoutes 注册路由（包含验证中间件）
func (h *ExampleValidationHandler) RegisterRoutes(router gin.IRouter) {
	users := router.Group("/users")
	{
		// 使用验证中间件
		users.POST("",
			middleware.ValidateJSON(&handler.CreateUserRequest{}),
			h.CreateUser,
		)

		users.PUT("/:id",
			middleware.ValidateJSON(&handler.UpdateUserRequest{}),
			h.UpdateUser,
		)
	}

	games := router.Group("/games")
	{
		games.POST("",
			middleware.ValidateJSON(&handler.CreateGameRequest{}),
			h.CreateGame,
		)

		games.PUT("/:id",
			middleware.ValidateJSON(&handler.UpdateGameRequest{}),
			h.UpdateGame,
		)

		// 查询参数验证示例
		games.GET("",
			middleware.ValidateQuery(map[string]string{
				"page":      "min:1",
				"page_size": "min:1,max:100",
				"keyword":   "max:100",
			}),
			h.ListGames,
		)
	}
}

// CreateUser 处理创建用户请求
func (h *ExampleValidationHandler) CreateUser(c *gin.Context) {
	var req handler.CreateUserRequest

	// 从Context中获取验证后的请求
	if !middleware.GetValidatedRequest(c, &req) {
		// 如果获取失败，说明验证中间件已经处理了错误响应
		return
	}

	// 在这里，req已经被验证过了，可以安全使用
	// 例如：req.Name, req.Email, req.Password 等都已经通过了验证

	// 调用服务层创建用户
	user, err := h.svc.CreateUser(c.Request.Context(), service.CreateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
		Role:     req.Role,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"code":    http.StatusCreated,
		"message": "用户创建成功",
		"data":    user,
	})
}

// UpdateUser 处理更新用户请求
func (h *ExampleValidationHandler) UpdateUser(c *gin.Context) {
	var req handler.UpdateUserRequest

	if !middleware.GetValidatedRequest(c, &req) {
		return
	}

	// 获取用户ID（从URL参数）
	userID := c.Param("id")
	// 这里应该将userID转换为uint64类型
	// 为了示例简化，这里省略了类型转换和错误处理

	// 调用服务层更新用户
	user, err := h.svc.UpdateUser(c.Request.Context(), 0, service.UpdateUserInput{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": "用户更新成功",
		"data":    user,
	})
}

// ListGames 处理游戏列表查询
func (h *ExampleValidationHandler) ListGames(c *gin.Context) {
	var req handler.ListQueryRequest

	// 绑定查询参数到结构体
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": "查询参数错误: " + err.Error(),
		})
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// 调用服务层获取游戏列表
	games, pagination, err := h.svc.ListGamesPaged(c.Request.Context(), req.Page, req.PageSize)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"code":       http.StatusOK,
		"message":    "获取成功",
		"data":       games,
		"pagination": pagination,
	})
}

// CreateGame 处理创建游戏请求
func (h *ExampleValidationHandler) CreateGame(c *gin.Context) {
	var req handler.CreateGameRequest

	if !middleware.GetValidatedRequest(c, &req) {
		return
	}

	game, err := h.svc.CreateGame(c.Request.Context(), service.CreateGameInput{
		Key:         req.Key,
		Name:        req.Name,
		Category:    req.Category,
		IconURL:     req.IconURL,
		Description: req.Description,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"code":    http.StatusCreated,
		"message": "游戏创建成功",
		"data":    game,
	})
}

// UpdateGame 处理更新游戏请求
func (h *ExampleValidationHandler) UpdateGame(c *gin.Context) {
	var req handler.UpdateGameRequest

	if !middleware.GetValidatedRequest(c, &req) {
		return
	}

	// 获取游戏ID（从URL参数）
	gameID := c.Param("id")
	// 这里应该将gameID转换为uint64类型

	game, err := h.svc.UpdateGame(c.Request.Context(), 0, service.UpdateGameInput{
		Name:        req.Name,
		Category:    req.Category,
		IconURL:     req.IconURL,
		Description: req.Description,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": "游戏更新成功",
		"data":    game,
	})
}

// 在实际项目中，你需要修改现有的admin路由注册代码：
//
// 原来的代码：
// admin.RegisterRoutes(api, adminService)
//
// 修改为：
// exampleHandler := NewExampleValidationHandler(adminService)
// exampleHandler.RegisterRoutes(api.Group("/validated"))
