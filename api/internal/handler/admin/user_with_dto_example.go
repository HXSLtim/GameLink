package admin

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/admin/dto"
	adminservice "gamelink/internal/service/admin"
	"gamelink/pkg/apierr"
)

// UserWithDTOHandler 使用DTO模式的用户处理器示例
// 这是重构后的版本，展示如何使用DTO分离API层和业务层
type UserWithDTOHandler struct {
	svc *adminservice.AdminService
}

// NewUserWithDTOHandler 创建Handler
func NewUserWithDTOHandler(svc *adminservice.AdminService) *UserWithDTOHandler {
	return &UserWithDTOHandler{svc: svc}
}

// CreateUser
// @Summary      创建用户 (DTO模式)
// @Description  使用DTO进行请求和响应处理
// @Tags         Admin/Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  dto.CreateUserRequest  true  "用户信息"
// @Success      201  {object}  dto.UserResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/users [post]
//
// CreateUser 创建用户的DTO版本
func (h *UserWithDTOHandler) CreateUser(c *gin.Context) {
	// Step 1: 绑定并验证请求DTO
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apierr.BadRequest("invalid request").WithDetails(err.Error()))
		return
	}

	// Step 2: DTO → Service Input 转换
	input := dto.ToCreateUserInput(&req)

	// Step 3: 调用Service层
	user, err := h.svc.CreateUser(c.Request.Context(), input)
	if err != nil {
		respondError(c, err)
		return
	}

	// Step 4: Entity → Response DTO 转换
	response := dto.ToUserResponse(user)

	// Step 5: 返回响应
	respondCreated(c, response)
}

// UpdateUser
// @Summary      更新用户 (DTO模式)
// @Tags         Admin/Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                     true  "用户ID"
// @Param        request  body  dto.UpdateUserRequest   true  "用户信息"
// @Success      200  {object}  dto.UserResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/users/{id} [put]
//
// UpdateUser 更新用户的DTO版本
func (h *UserWithDTOHandler) UpdateUser(c *gin.Context) {
	// Step 1: 解析路径参数
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	// Step 2: 绑定并验证请求DTO
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apierr.BadRequest("invalid request").WithDetails(err.Error()))
		return
	}

	// Step 3: DTO → Service Input 转换
	input := dto.ToUpdateUserInput(&req)

	// Step 4: 调用Service层
	user, err := h.svc.UpdateUser(c.Request.Context(), id, input)
	if err != nil {
		respondError(c, err)
		return
	}

	// Step 5: Entity → Response DTO 转换
	response := dto.ToUserResponse(user)

	// Step 6: 返回响应
	respondSuccess(c, response)
}

// GetUser
// @Summary      获取用户 (DTO模式)
// @Tags         Admin/Users
// @Security     BearerAuth
// @Param        id   path      int  true  "用户ID"
// @Produce      json
// @Success      200  {object}  dto.UserResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/users/{id} [get]
//
// GetUser 获取用户的DTO版本
func (h *UserWithDTOHandler) GetUser(c *gin.Context) {
	// Step 1: 解析路径参数
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	// Step 2: 调用Service层
	user, err := h.svc.GetUser(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	// Step 3: Entity → Response DTO 转换（带脱敏）
	response := dto.ToUserResponseWithConfig(user, dto.MapperConfig{
		MaskSensitive: false, // 管理后台通常不脱敏，用户端API需要脱敏
	})

	// Step 4: 返回响应
	respondSuccess(c, response)
}

// ListUsers
// @Summary      列出用户 (DTO模式)
// @Tags         Admin/Users
// @Security     BearerAuth
// @Param        page     query     int       false  "页码"
// @Param        pageSize query     int       false  "每页数量"
// @Param        role     query     string    false  "角色筛选"
// @Param        status   query     string    false  "状态筛选"
// @Produce      json
// @Success      200  {object}  dto.UserListResponse
// @Router       /admin/users [get]
//
// ListUsers 列出用户的DTO版本
func (h *UserWithDTOHandler) ListUsers(c *gin.Context) {
	// Step 1: 构建查询选项（这部分可以提取为DTO）
	opts, ok := buildUserListOptions(c)
	if !ok {
		return
	}

	// Step 2: 调用Service层
	users, pagination, err := h.svc.ListUsersWithOptions(c.Request.Context(), opts)
	if err != nil {
		respondError(c, err)
		return
	}

	// Step 3: Entity List → Response DTO 转换
	response := dto.ToUserListResponseWithConfig(
		users,
		pagination.Total,
		opts.Page,
		opts.PageSize,
		dto.MapperConfig{MaskSensitive: false},
	)

	// Step 4: 返回响应
	c.JSON(200, response)
}

// ==================== 对比：使用DTO前后的区别 ====================

/*
❌ 使用DTO前的问题：

1. Handler 直接操作 model.User
   - API 变更 = 数据库变更
   - 无法控制响应字段
   - 敏感字段可能泄漏

2. 验证逻辑混乱
   - 验证分散在 Handler/Service
   - 难以复用
   - 错误信息不统一

3. 难以测试
   - Handler 需要构造完整的 model.User
   - 测试依赖数据库结构


✅ 使用DTO后的优势：

1. API 层和数据层解耦
   - API 字段变更不影响数据库
   - 可以隐藏/转换字段
   - 支持版本演进

2. 清晰的数据流
   Request DTO → Service Input → Entity → Response DTO
   
3. 更好的验证
   - binding tags 统一验证规则
   - 类型安全
   - 自动生成API文档

4. 安全性提升
   - 可控制的字段暴露
   - 支持数据脱敏
   - 防止质量注入

5. 易于测试
   - Mock Request DTO 即可
   - 不依赖数据库结构
   - 单元测试更简单
*/
