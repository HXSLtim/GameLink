package main

// GameLink Swagger泛型修复示例代码
// 这个文件展示了如何修复不同类型的Swagger注解问题

// ============================================================================
// 示例1: 修复使用 `object` 类型的接口
// ============================================================================

// 问题代码 (internal/handler/player/order.go)
/*
// @Summary      获取可接订单列表
// @Tags         Player - Orders
// @Success      200      {object}  object  // ❌ 错误：使用了object类型
// @Failure      400      {object}  apierr.APIError
func getAvailableOrdersHandler(c *gin.Context, svc *order.OrderService) {
    orders, total, err := svc.GetAvailableOrders(c.Request.Context(), req)
    if err != nil {
        respondAPIError(c, apierr.InternalError("获取可接订单列表失败"))
        return
    }

    respondSuccess(c, "OK", map[string]interface{}{  // ❌ 实际返回具体结构
        "orders": orders,
        "total":  total,
    })
}
*/

// 修复后的代码
/*
// 第一步：定义响应类型
type AvailableOrdersResponse struct {
    Orders []model.Order `json:"orders"`
    Total  int           `json:"total"`
}

// 第二步：修复Swagger注解和实现
// @Summary      获取可接订单列表
// @Tags         Player - Orders  
// @Success      200 {object} model.APIResponse[AvailableOrdersResponse]  // ✅ 正确：指定具体类型
// @Failure      400 {object} model.ErrorResponse
func getAvailableOrdersHandler(c *gin.Context, svc *order.OrderService) {
    orders, total, err := svc.GetAvailableOrders(c.Request.Context(), req)
    if err != nil {
        respondAPIError(c, apierr.InternalError("获取可接订单列表失败"))
        return
    }

    // ✅ 正确：使用泛型APIResponse包装具体类型
    respondJSON(c, http.StatusOK, model.APIResponse[AvailableOrdersResponse]{
        Success: true,
        Code:    http.StatusOK,
        Message: "OK",
        Data: AvailableOrdersResponse{
            Orders: orders,
            Total:  total,
        },
    })
}
*/

// ============================================================================
// 示例2: 修复使用旧响应类型的接口
// ============================================================================

// 问题代码 (internal/handler/admin/user.go)
/*
// @Summary      列出用户
// @Tags         Admin/Users
// @Success      200 {object} model.SuccessResponse  // ❌ 错误：注解与实现不符
// @Failure      404 {object} model.ErrorResponse
func (h *UserHandler) ListUsers(c *gin.Context) {
    users, pagination, err := h.svc.ListUsersWithOptions(c.Request.Context(), opts)
    if err != nil {
        respondAPIError(c, apierr.InternalError("获取用户列表失败"))
        return
    }

    // ❌ 实际代码返回泛型类型，但Swagger注解说返回SuccessResponse
    writeJSON(c, http.StatusOK, model.APIResponse[[]model.User]{
        Success:    true,
        Code:       http.StatusOK,
        Message:    "OK",
        Data:       users,
        Pagination: pagination,
    })
}
*/

// 修复后的代码
/*
// @Summary      列出用户
// @Tags         Admin/Users
// @Success      200 {object} model.APIResponse[[]model.User]  // ✅ 正确：与实际返回类型一致
// @Failure      404 {object} model.ErrorResponse  
func (h *UserHandler) ListUsers(c *gin.Context) {
    users, pagination, err := h.svc.ListUsersWithOptions(c.Request.Context(), opts)
    if err != nil {
        respondAPIError(c, apierr.InternalError("获取用户列表失败"))
        return
    }

    // ✅ 保持不变，确保注解与实际一致
    writeJSON(c, http.StatusOK, model.APIResponse[[]model.User]{
        Success:    true,
        Code:       http.StatusOK,
        Message:    "OK",
        Data:       users,
        Pagination: pagination,
    })
}
*/

// ============================================================================
// 示例3: 修复Auth模块接口
// ============================================================================

// 问题代码 (internal/handler/auth.go)
/*
type loginResponse struct {
    Token     string     `json:"token"`
    ExpiresAt time.Time  `json:"expires_at"`
    User      model.User `json:"user"`
}

// @Summary      登录
// @Tags         Auth
// @Success      200 {object} loginResponse  // ❌ 错误：直接返回自定义类型
// @Failure      401 {object} apierr.APIError
func loginHandler(c *gin.Context, svc *authservice.AuthService) {
    // ... 登录逻辑 ...
    
    // ❌ 错误：直接返回自定义类型，没有用APIResponse包装
    RespondSuccess(c, "登录成功", loginResponse{
        Token:     resp.Token,
        ExpiresAt: resp.ExpiresAt,
        User:      resp.User,
    })
}
*/

// 修复后的代码
/*
// @Summary      登录
// @Tags         Auth  
// @Success      200 {object} model.APIResponse[loginResponse]  // ✅ 正确：使用泛型包装
// @Failure      401 {object} model.ErrorResponse
func loginHandler(c *gin.Context, svc *authservice.AuthService) {
    // ... 登录逻辑 ...
    
    // ✅ 正确：使用泛型APIResponse包装
    respondJSON(c, http.StatusOK, model.APIResponse[loginResponse]{
        Success: true,
        Code:    http.StatusOK,
        Message: "登录成功",
        Data: loginResponse{
            Token:     resp.Token,
            ExpiresAt: resp.ExpiresAt,
            User:      resp.User,
        },
    })
}
*/

// ============================================================================
// 示例4: 复杂嵌套类型的处理
// ============================================================================

// 对于复杂的数据结构，应该定义清晰的类型
type OrderDetailResponse struct {
    Order         model.Order         `json:"order"`
    Player        model.Player        `json:"player"`
    User          model.User          `json:"user"`
    Payment       *model.Payment      `json:"payment,omitempty"`
    Reviews       []model.Review      `json:"reviews"`
    Timeline      []TimelineItem      `json:"timeline"`
}

type TimelineItem struct {
    Time    time.Time `json:"time"`
    Action  string    `json:"action"`
    Message string    `json:"message"`
}

// @Summary      获取订单详情
// @Tags         Orders
// @Success      200 {object} model.APIResponse[OrderDetailResponse]  // ✅ 复杂类型也清晰定义
// @Failure      404 {object} model.ErrorResponse
func getOrderDetail(c *gin.Context) {
    // ... 获取订单详情逻辑 ...
    
    respondJSON(c, http.StatusOK, model.APIResponse[OrderDetailResponse]{
        Success: true,
        Code:    http.StatusOK,
        Message: "OK",
        Data:    orderDetail,  // OrderDetailResponse类型
    })
}

// ============================================================================
// 示例5: 列表和分页响应的处理
// ============================================================================

// 列表响应类型
type UserListResponse struct {
    Users      []model.User `json:"users"`
    Pagination Pagination   `json:"pagination"`
}

type Pagination struct {
    Page       int  `json:"page"`
    PageSize   int  `json:"page_size"`
    Total      int  `json:"total"`
    TotalPages int  `json:"total_pages"`
}

// @Summary      获取用户列表
// @Tags         Users
// @Success      200 {object} model.APIResponse[UserListResponse]  // ✅ 列表类型清晰定义
// @Failure      500 {object} model.ErrorResponse
func listUsers(c *gin.Context) {
    users, pagination, err := svc.ListUsers(c.Request.Context(), page, pageSize)
    if err != nil {
        respondAPIError(c, apierr.InternalError("获取用户列表失败"))
        return
    }
    
    respondJSON(c, http.StatusOK, model.APIResponse[UserListResponse]{
        Success: true,
        Code:    http.StatusOK,
        Message: "OK",
        Data: UserListResponse{
            Users:      users,
            Pagination: pagination,
        },
    })
}

// ============================================================================
// 示例6: 错误响应的统一处理
// ============================================================================

// 统一的错误响应处理函数
func respondAPIError(c *gin.Context, apiErr *apierr.APIError) {
    respondJSON(c, apiErr.Code, model.APIResponse[any]{
        Success: false,
        Code:    apiErr.Code,
        Message: apiErr.Message,
        TraceID: c.GetString("trace_id"),
    })
}

// 使用示例
// @Summary      创建用户
// @Tags         Users
// @Success      201 {object} model.APIResponse[model.User]  // ✅ 成功响应
// @Failure      400 {object} model.ErrorResponse  // ✅ 统一的错误响应
// @Failure      409 {object} model.ErrorResponse  // ✅ 冲突错误
// @Failure      500 {object} model.ErrorResponse  // ✅ 服务器错误
func createUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondAPIError(c, apierr.BadRequest("无效的请求参数"))
        return
    }
    
    user, err := svc.CreateUser(c.Request.Context(), req)
    if err != nil {
        if errors.Is(err, ErrUserExists) {
            respondAPIError(c, apierr.Conflict("用户已存在"))
            return
        }
        respondAPIError(c, apierr.InternalError("创建用户失败"))
        return
    }
    
    respondJSON(c, http.StatusCreated, model.APIResponse[model.User]{
        Success: true,
        Code:    http.StatusCreated,
        Message: "创建成功",
        Data:    user,
    })
}

// ============================================================================
// 示例7: 工具函数的统一
// ============================================================================

// 建议的统一响应工具函数
func respondJSON[T any](c *gin.Context, status int, payload model.APIResponse[T]) {
    c.JSON(status, payload)
}

func respondSuccess[T any](c *gin.Context, message string, data T) {
    respondJSON(c, http.StatusOK, model.APIResponse[T]{
        Success: true,
        Code:    http.StatusOK,
        Message: message,
        Data:    data,
    })
}

func respondError(c *gin.Context, code int, message string) {
    respondJSON(c, code, model.APIResponse[any]{
        Success: false,
        Code:    code,
        Message: message,
    })
}

// 使用工具函数的示例
// @Summary      获取用户信息
// @Tags         Users
// @Success      200 {object} model.APIResponse[model.User]
// @Failure      404 {object} model.ErrorResponse
func getUser(c *gin.Context) {
    userID, err := parseUintParam(c, "id")
    if err != nil {
        respondError(c, http.StatusBadRequest, "无效的用户ID")
        return
    }
    
    user, err := svc.GetUser(c.Request.Context(), userID)
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            respondAPIError(c, apierr.NotFound("用户不存在"))
            return
        }
        respondAPIError(c, apierr.InternalError("获取用户信息失败"))
        return
    }
    
    respondSuccess(c, "OK", user)
}

// ============================================================================
// 总结：修复原则和最佳实践
// ============================================================================

/*
修复原则：
1. 统一使用 model.APIResponse[T] 格式
2. Swagger注解必须与实际返回类型完全一致  
3. 避免使用 object 作为响应类型
4. 复杂响应结构应该定义专门的类型
5. 保持代码和文档的同步更新

最佳实践：
1. 先定义响应类型，再写Swagger注解
2. 使用工具函数统一响应格式
3. 错误响应也要使用统一的APIResponse格式
4. 定期检查和更新Swagger文档
5. 在代码审查中检查Swagger注解的准确性
*/