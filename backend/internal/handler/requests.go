package handler

import "time"

// CreateUserRequest 创建用户的请求结构体
//
// 使用方法：
// router.POST("/users", middleware.ValidateJSON(&CreateUserRequest{}), createUserHandler)
//
// 验证规则说明：
// binding:"required" - 必填字段
// validate:"min=1,max=50" - 长度限制
// validate:"email" - 邮箱格式
// validate:"oneof=admin user player" - 枚举值
// validate:"phone" - 自定义手机号验证
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required" validate:"min=1,max=50"`
	Email    string `json:"email" binding:"required" validate:"email,max=100"`
	Phone    string `json:"phone" validate:"phone"`
	Password string `json:"password" binding:"required" validate:"password"`
	Role     string `json:"role" validate:"oneof=admin user player"`
}

// UpdateUserRequest 更新用户的请求结构体
//
// 使用omitempty标签，表示字段可以为空
type UpdateUserRequest struct {
	Name  string `json:"name,omitempty" validate:"min=1,max=50"`
	Email string `json:"email,omitempty" validate:"email,max=100"`
	Phone string `json:"phone,omitempty" validate:"phone"`
}

// CreateGameRequest 创建游戏的请求结构体
type CreateGameRequest struct {
	Key         string `json:"key" binding:"required" validate:"min=1,max=32"`
	Name        string `json:"name" binding:"required" validate:"min=1,max=100"`
	Category    string `json:"category" validate:"max=50"`
	IconURL     string `json:"icon_url" validate:"url,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

// UpdateGameRequest 更新游戏的请求结构体
type UpdateGameRequest struct {
	Name        string `json:"name,omitempty" validate:"min=1,max=100"`
	Category    string `json:"category,omitempty" validate:"max=50"`
	IconURL     string `json:"icon_url,omitempty" validate:"url,max=255"`
	Description string `json:"description,omitempty" validate:"max=1000"`
}

// ListQueryRequest 分页查询的请求结构体
//
// 用于API的查询参数验证
type ListQueryRequest struct {
	Page     int    `form:"page" validate:"min=1"`
	PageSize int    `form:"page_size" validate:"min=1,max=100"`
	Keyword  string `form:"keyword" validate:"max=100"`
	SortBy   string `form:"sort_by" validate:"oneof=created_at updated_at name"`
	SortDesc bool   `form:"sort_desc"`
}

// DateRangeQueryRequest 日期范围查询的请求结构体
type DateRangeQueryRequest struct {
	ListQueryRequest
	StartDate time.Time `form:"start_date" validate:"required"`
	EndDate   time.Time `form:"end_date" validate:"required,gtfield=StartDate"`
}

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required" validate:"min=3,max=50"`
	Password string `json:"password" binding:"required" validate:"min=6"`
}

// 注意：这些请求结构体需要在相应的handler中使用
// 示例用法：
//
// func (h *UserHandler) CreateUser(c *gin.Context) {
//     var req CreateUserRequest
//     if !middleware.GetValidatedRequest(c, &req) {
//         return // 验证失败，中间件已经处理了响应
//     }
//
//     // 使用验证后的数据进行业务逻辑处理
//     user, err := h.svc.CreateUser(c.Request.Context(), service.CreateUserInput{
//         Name:     req.Name,
//         Email:    req.Email,
//         Phone:    req.Phone,
//         Password: req.Password,
//         Role:     req.Role,
//     })
//
//     if err != nil {
//         writeJSONError(c, http.StatusInternalServerError, err.Error())
//         return
//     }
//
//     writeJSON(c, http.StatusCreated, user)
// }
