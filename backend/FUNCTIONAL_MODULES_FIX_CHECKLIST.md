# GameLink 功能模块代码整改清单

**整改目标**: 按照功能模块对代码进行系统性审查和整改  
**整改日期**: 2025-11-22  
**执行人**: AI Code Review Agent

---

## 📋 整改清单总览

### 🔴 高优先级 (安全性 & 性能问题)
- [ ] **1. JWT安全加固** - 移除硬编码密钥
- [ ] **2. 输入验证增强** - 防止XSS和SQL注入
- [ ] **3. 数据库查询优化** - 添加索引和缓存
- [ ] **4. 错误处理统一** - 标准化错误响应

### 🟡 中优先级 (功能完善 & 可维护性)
- [ ] **5. 订单状态机优化** - 集中状态转换逻辑
- [ ] **6. 支付幂等性** - 防止重复支付
- [ ] **7. 缓存策略实施** - Redis缓存层
- [ ] **8. 代码复用提升** - 泛型Repository

### 🟢 低优先级 (监控 & 开发体验)
- [ ] **9. 监控指标添加** - Prometheus集成
- [ ] **10. 日志规范化** - 结构化日志
- [ ] **11. 测试覆盖率提升** - 补充单元测试

---

## 🔴 高优先级整改任务

### 1. JWT安全加固

#### 1.1 移除硬编码JWT密钥
**文件**: `internal/handler/middleware/jwt_auth.go`
**位置**: 第34行
**问题**: 开发环境硬编码默认密钥

**整改前**:
```go
if secretKey == "" {
    if os.Getenv("APP_ENV") == "production" {
        return func(c *gin.Context) {
            c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
                "success": false,
                "code":    http.StatusServiceUnavailable,
                "message": "jwt not configured",
            })
        }
    }
    // 开发环境使用默认值
    secretKey = "gamelink-default-secret-key-change-in-production"
}
```

**整改后**:
```go
if secretKey == "" {
    logging.Error("JWT_SECRET_KEY not configured")
    return func(c *gin.Context) {
        c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
            "success": false,
            "code":    http.StatusServiceUnavailable,
            "message": "认证服务配置错误，请联系管理员",
        })
    }
}

// 验证密钥长度
if len(secretKey) < 32 {
    logging.Error("JWT_SECRET_KEY too short, must be at least 32 characters")
    return func(c *gin.Context) {
        c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
            "success": false,
            "code":    http.StatusServiceUnavailable,
            "message": "认证服务配置错误，请联系管理员",
        })
    }
}
```

**验证方式**:
```bash
# 不设置JWT_SECRET_KEY应该无法启动
unset JWT_SECRET_KEY
go run cmd/main.go
# 应该看到错误日志和503响应

# 设置正确的密钥
export JWT_SECRET_KEY="your-32-characters-or-longer-secret-key-here"
go run cmd/main.go
# 应该正常启动
```

---

#### 1.2 添加Token自动刷新机制
**文件**: `internal/handler/middleware/jwt_auth.go`
**位置**: 第98-100行
**问题**: 只提示刷新，不自动刷新

**整改前**:
```go
// 检查Token剩余时间，如果快要过期，在响应头中提示前端刷新Token
remainingTime := auth.GetTokenRemainingTime(claims)
if remainingTime < 1*time.Hour {
    c.Header("X-Token-Refresh-Recommendation", "true")
}
```

**整改后**:
```go
// 检查Token剩余时间
remainingTime := auth.GetTokenRemainingTime(claims)

// 如果Token即将过期（15分钟内），自动刷新
if remainingTime < 15*time.Minute {
    newToken, err := jwtManager.RefreshToken(claims)
    if err == nil {
        // 在响应头中返回新Token
        c.Header("X-Refreshed-Token", newToken)
        
        // 更新Context中的Token信息
        newClaims, _ := jwtManager.VerifyToken(newToken)
        if newClaims != nil {
            c.Set("jwt_claims", newClaims)
            c.Set("user_id", newClaims.UserID)
            c.Set("user_role", newClaims.Role)
        }
        
        logging.Debug("Token auto-refreshed", "user_id", claims.UserID)
    } else {
        logging.Warn("Failed to refresh token", "error", err, "user_id", claims.UserID)
    }
} else if remainingTime < 1*time.Hour {
    // 仍然保留提示，让前端可以主动刷新
    c.Header("X-Token-Refresh-Recommendation", "true")
}
```

---

### 2. 输入验证增强

#### 2.1 添加XSS过滤
**文件**: `internal/service/user/user.go` (新建或修改现有文件)
**问题**: 用户输入直接存储，没有XSS过滤

**整改内容**:
```go
package user

import (
    "context"
    "github.com/microcosm-cc/bluemonday"
    "gamelink/internal/model"
)

// XSS净化策略
var xssPolicy = bluemonday.UGCPolicy()

func sanitizeInput(input string) string {
    return xssPolicy.Sanitize(input)
}

type UpdateProfileRequest struct {
    Name     string `json:"name" binding:"required,max=64"`
    Bio      string `json:"bio" binding:"max=500"`
    AvatarURL string `json:"avatarUrl" binding:"url,max=255"`
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uint64, req UpdateProfileRequest) error {
    user, err := s.userRepo.Get(ctx, userID)
    if err != nil {
        return err
    }
    
    // 净化用户输入，防止XSS攻击
    user.Name = sanitizeInput(req.Name)
    user.Bio = sanitizeInput(req.Bio)
    user.AvatarURL = req.AvatarURL // URL已经在binding中验证
    
    return s.userRepo.Update(ctx, user)
}
```

**需要安装依赖**:
```bash
go get github.com/microcosm-cc/bluemonday
```

---

#### 2.2 增强邮箱验证
**文件**: `internal/service/auth/auth.go`
**位置**: 第271-277行
**问题**: 邮箱验证过于简单

**整改前**:
```go
func isValidEmail(email string) bool {
    if email == "" {
        return false
    }
    _, err := mail.ParseAddress(email)
    return err == nil
}
```

**整改后**:
```go
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
    if email == "" || len(email) > 128 {
        return false
    }
    
    // 基本格式验证
    _, err := mail.ParseAddress(email)
    if err != nil {
        return false
    }
    
    // 正则表达式验证
    if !emailRegex.MatchString(email) {
        return false
    }
    
    // 检查常见临时邮箱域名（可选）
    disposableDomains := []string{"tempmail.com", "10minutemail.com", "guerrillamail.com"}
    parts := strings.Split(email, "@")
    if len(parts) == 2 {
        domain := strings.ToLower(parts[1])
        for _, disposable := range disposableDomains {
            if domain == disposable || strings.HasSuffix(domain, "."+disposable) {
                return false // 拒绝临时邮箱
            }
        }
    }
    
    return true
}
```

**需要添加的导入**:
```go
import (
    "regexp"
    "strings"
)
```

---

### 3. 数据库查询优化

#### 3.1 为用户表添加索引
**文件**: `internal/model/user.go`
**问题**: 常用查询字段缺少索引

**整改前**:
```go
type User struct {
    Base
    Phone        string     `json:"phone,omitempty" gorm:"size:32;uniqueIndex"`
    Email        string     `json:"email,omitempty" gorm:"size:128;uniqueIndex"`
    Name         string     `json:"name" gorm:"size:64"`
    Status       UserStatus `json:"status" gorm:"size:32;index"`
    LastLoginAt  *time.Time `json:"lastLoginAt,omitempty" gorm:"column:last_login_at"`
}
```

**整改后**:
```go
type User struct {
    Base
    Phone        string     `json:"phone,omitempty" gorm:"size:32;uniqueIndex"`
    Email        string     `json:"email,omitempty" gorm:"size:128;uniqueIndex"`
    Name         string     `json:"name" gorm:"size:64;index"` // 添加索引，用于搜索
    Status       UserStatus `json:"status" gorm:"size:32;index;index:idx_status_last_login,priority:1"`
    LastLoginAt  *time.Time `json:"lastLoginAt,omitempty" gorm:"column:last_login_at;index:idx_status_last_login,priority:2"`
    
    // 添加复合索引，优化状态+时间的查询
    // 例如：查询最近登录的活跃用户
}
```

**需要创建的数据库迁移**:
```sql
-- 为name字段添加索引
CREATE INDEX idx_users_name ON users(name);

-- 创建复合索引
CREATE INDEX idx_users_status_last_login ON users(status, last_login_at);
```

---

#### 3.2 为订单表添加复合索引
**文件**: `internal/model/order.go`
**问题**: 多条件查询性能差

**整改前**:
```go
type Order struct {
    Base
    OrderNo           string      `gorm:"size:64;uniqueIndex"`
    UserID            uint64      `gorm:"not null;index"`
    ItemID            uint64      `gorm:"not null;index"`
    PlayerID          *uint64     `gorm:"index"`
    Status            OrderStatus `gorm:"size:32;index;default:'pending'"`
    CreatedAt         time.Time
}
```

**整改后**:
```go
type Order struct {
    Base
    OrderNo           string      `gorm:"size:64;uniqueIndex"`
    UserID            uint64      `gorm:"not null;index;index:idx_user_status_created,priority:1"`
    ItemID            uint64      `gorm:"not null;index"`
    PlayerID          *uint64     `gorm:"index;index:idx_player_status,priority:2"`
    Status            OrderStatus `gorm:"size:32;index;default:'pending';index:idx_user_status_created,priority:2;index:idx_player_status,priority:1"`
    CreatedAt         time.Time   `gorm:"index:idx_user_status_created,priority:3"`
    
    // 复合索引说明：
    // idx_user_status_created: (user_id, status, created_at) - 优化用户订单列表查询
    // idx_player_status: (status, player_id) - 优化陪玩师订单列表查询
}
```

**需要创建的数据库迁移**:
```sql
-- 创建复合索引，优化用户订单查询
CREATE INDEX idx_orders_user_status_created ON orders(user_id, status, created_at DESC);

-- 创建复合索引，优化陪玩师订单查询
CREATE INDEX idx_orders_player_status ON orders(status, player_id, created_at DESC);

-- 为RecipientPlayerID添加索引（礼物订单）
CREATE INDEX idx_orders_recipient_player ON orders(recipient_player_id) WHERE recipient_player_id IS NOT NULL;
```

---

#### 3.3 为用户Repository添加Redis缓存
**文件**: `internal/repository/user/repository.go`
**问题**: 频繁查询用户信息没有缓存

**整改前**:
```go
type gormUserRepository struct {
    db *gorm.DB
}

func (r *gormUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
    var user model.User
    if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, repository.ErrNotFound
        }
        return nil, err
    }
    return &user, nil
}
```

**整改后**:
```go
type gormUserRepository struct {
    db    *gorm.DB
    cache cache.Cache
}

func NewUserRepository(db *gorm.DB, cache cache.Cache) repository.UserRepository {
    return &gormUserRepository{
        db:    db,
        cache: cache,
    }
}

func (r *gormUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)
    
    // 先查缓存
    var user model.User
    if err := r.cache.Get(ctx, cacheKey, &user); err == nil {
        if user.ID == 0 {
            return nil, repository.ErrNotFound // 缓存的空值标记
        }
        return &user, nil
    }
    
    // 查询数据库
    if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            // 缓存空值，防止缓存穿透
            _ = r.cache.Set(ctx, cacheKey, model.User{}, 5*time.Minute)
            return nil, repository.ErrNotFound
        }
        return nil, err
    }
    
    // 写入缓存
    _ = r.cache.Set(ctx, cacheKey, user, 15*time.Minute)
    
    return &user, nil
}

func (r *gormUserRepository) Update(ctx context.Context, user *model.User) error {
    tx := r.db.WithContext(ctx).Model(user).Updates(map[string]any{
        "phone":         user.Phone,
        "email":         user.Email,
        "name":          user.Name,
        "avatar_url":    user.AvatarURL,
        "role":          user.Role,
        "status":        user.Status,
        "password_hash": user.PasswordHash,
        "last_login_at": user.LastLoginAt,
    })
    if tx.Error != nil {
        return tx.Error
    }
    if tx.RowsAffected == 0 {
        return repository.ErrNotFound
    }
    
    // 更新缓存
    cacheKey := fmt.Sprintf("user:%d", user.ID)
    _ = r.cache.Delete(ctx, cacheKey)
    
    return nil
}

func (r *gormUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    cacheKey := fmt.Sprintf("user:email:%s", email)
    
    var user model.User
    if err := r.cache.Get(ctx, cacheKey, &user); err == nil {
        if user.ID == 0 {
            return nil, repository.ErrNotFound
        }
        return &user, nil
    }
    
    if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            _ = r.cache.Set(ctx, cacheKey, model.User{}, 5*time.Minute)
            return nil, repository.ErrNotFound
        }
        return nil, err
    }
    
    _ = r.cache.Set(ctx, cacheKey, user, 15*time.Minute)
    return &user, nil
}
```

**需要添加的导入**:
```go
import (
    "context"
    "fmt"
    "time"
    "gamelink/internal/cache"
)
```

---

### 4. 错误处理统一

#### 4.1 创建统一的错误包
**文件**: `internal/apierr/errors.go` (新建)
**问题**: 错误处理分散，没有统一格式

**整改内容**:
```go
package apierr

import (
    "errors"
    "net/http"
)

// 通用错误定义
var (
    ErrNotFound         = errors.New("资源不存在")
    ErrUnauthorized     = errors.New("未授权访问")
    ErrForbidden        = errors.New("权限不足")
    ErrInvalidInput     = errors.New("输入参数无效")
    ErrInternal         = errors.New("服务器内部错误")
    ErrConflict         = errors.New("资源冲突")
    ErrTooManyRequests  = errors.New("请求过于频繁")
)

// APIError 统一的API错误响应
type APIError struct {
    Code       int         `json:"code"`
    Message    string      `json:"message"`
    Details    string      `json:"details,omitempty"`
    Field      string      `json:"field,omitempty"`
    RequestID  string      `json:"requestId,omitempty"`
}

func (e APIError) Error() string {
    return e.Message
}

// New 创建新的API错误
func New(code int, message string) *APIError {
    return &APIError{
        Code:    code,
        Message: message,
    }
}

// WithDetails 添加错误详情
func (e *APIError) WithDetails(details string) *APIError {
    e.Details = details
    return e
}

// WithField 添加错误字段
func (e *APIError) WithField(field string) *APIError {
    e.Field = field
    return e
}

// WithRequestID 添加请求ID
func (e *APIError) WithRequestID(requestID string) *APIError {
    e.RequestID = requestID
    return e
}

// 常用错误构造函数
func BadRequest(message string) *APIError {
    return New(http.StatusBadRequest, message)
}

func Unauthorized(message string) *APIError {
    return New(http.StatusUnauthorized, message)
}

func Forbidden(message string) *APIError {
    return New(http.StatusForbidden, message)
}

func NotFound(message string) *APIError {
    return New(http.StatusNotFound, message)
}

func Conflict(message string) *APIError {
    return New(http.StatusConflict, message)
}

func InternalError(message string) *APIError {
    return New(http.StatusInternalServerError, message)
}

// ValidationError 参数验证错误
type ValidationError struct {
    APIError
    Field   string `json:"field"`
    Value   string `json:"value"`
}

func NewValidationError(field, message string) *ValidationError {
    return &ValidationError{
        APIError: APIError{
            Code:    http.StatusBadRequest,
            Message: message,
        },
        Field: field,
    }
}
```

---

#### 4.2 在Handler中使用统一错误
**文件**: `internal/handler/auth.go`
**位置**: 错误处理部分

**整改前**:
```go
func loginHandler(c *gin.Context, svc *authservice.AuthService) {
    var req loginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondError(c, http.StatusBadRequest, ErrInvalidJSONPayload)
        return
    }
    
    resp, err := svc.Login(c.Request.Context(), authservice.LoginRequest{Username: req.Username, Password: req.Password})
    if err != nil {
        status := http.StatusUnauthorized
        switch err {
        case service.ErrInvalidCredentials:
            status = http.StatusUnauthorized
        case service.ErrUserDisabled:
            status = http.StatusForbidden
        default:
            status = http.StatusUnauthorized
        }
        respondJSON(c, status, model.APIResponse[any]{Success: false, Code: status, Message: err.Error()})
        return
    }
    // ...
}
```

**整改后**:
```go
func loginHandler(c *gin.Context, svc *authservice.AuthService) {
    var req loginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        apiErr := apierr.BadRequest("无效的请求格式").WithDetails(err.Error())
        respondAPIError(c, apiErr)
        return
    }
    
    resp, err := svc.Login(c.Request.Context(), authservice.LoginRequest{Username: req.Username, Password: req.Password})
    if err != nil {
        var apiErr *apierr.APIError
        
        switch err {
        case service.ErrInvalidCredentials:
            apiErr = apierr.Unauthorized("用户名或密码错误")
        case service.ErrUserDisabled:
            apiErr = apierr.Forbidden("账号已被禁用")
        case service.ErrNotFound:
            apiErr = apierr.NotFound("用户不存在")
        default:
            apiErr = apierr.InternalError("登录失败，请稍后重试").WithDetails(err.Error())
        }
        
        respondAPIError(c, apiErr)
        return
    }
    
    respondJSON(c, http.StatusOK, model.APIResponse[loginResponse]{
        Success: true,
        Code:    http.StatusOK,
        Message: "登录成功",
        Data:    loginResponse{Token: resp.Token, ExpiresAt: resp.ExpiresAt, User: resp.User},
    })
}

// 统一的错误响应函数
func respondAPIError(c *gin.Context, apiErr *apierr.APIError) {
    // 添加请求ID（如果中间件设置了）
    if requestID, exists := c.Get("request_id"); exists {
        if reqID, ok := requestID.(string); ok {
            apiErr.RequestID = reqID
        }
    }
    
    c.JSON(apiErr.Code, model.APIResponse[any]{
        Success: false,
        Code:    apiErr.Code,
        Message: apiErr.Message,
    })
}
```

---

## 🟡 中优先级整改任务

### 5. 订单状态机优化

#### 5.1 集中状态转换逻辑
**文件**: `internal/model/order.go`
**问题**: 状态转换逻辑分散在各个服务中

**整改内容**:
```go
package model

import "fmt"

// OrderStatus defines lifecycle states for an order.
type OrderStatus string

// OrderStatus values define the lifecycle of an order.
const (
    OrderStatusPending    OrderStatus = "pending"
    OrderStatusConfirmed  OrderStatus = "confirmed"
    OrderStatusInProgress OrderStatus = "in_progress"
    OrderStatusCompleted  OrderStatus = "completed"
    OrderStatusCanceled   OrderStatus = "canceled"
    OrderStatusRefunded   OrderStatus = "refunded"
)

// CanTransitionTo 检查状态转换是否合法
func (s OrderStatus) CanTransitionTo(target OrderStatus) bool {
    // 定义状态机转换规则
    transitions := map[OrderStatus][]OrderStatus{
        OrderStatusPending: {
            OrderStatusConfirmed,  // 支付成功后确认
            OrderStatusCanceled,   // 用户或系统取消
        },
        OrderStatusConfirmed: {
            OrderStatusInProgress, // 陪玩师开始服务
            OrderStatusCanceled,   // 开始前取消
            OrderStatusRefunded,   // 退款
        },
        OrderStatusInProgress: {
            OrderStatusCompleted,  // 服务完成
            OrderStatusCanceled,   // 服务中取消
            OrderStatusRefunded,   // 退款
        },
        OrderStatusCompleted: {}, // 终态，不可转换
        OrderStatusCanceled:  {}, // 终态，不可转换
        OrderStatusRefunded:  {}, // 终态，不可转换
    }
    
    allowed, ok := transitions[s]
    if !ok {
        return false
    }
    
    for _, status := range allowed {
        if status == target {
            return true
        }
    }
    return false
}

// GetStatusName 获取状态中文名称
func (s OrderStatus) GetStatusName() string {
    names := map[OrderStatus]string{
        OrderStatusPending:    "待确认",
        OrderStatusConfirmed:  "已确认",
        OrderStatusInProgress: "进行中",
        OrderStatusCompleted:  "已完成",
        OrderStatusCanceled:   "已取消",
        OrderStatusRefunded:   "已退款",
    }
    
    if name, ok := names[s]; ok {
        return name
    }
    return string(s)
}

// IsTerminal 检查是否为终态
func (s OrderStatus) IsTerminal() bool {
    switch s {
    case OrderStatusCompleted, OrderStatusCanceled, OrderStatusRefunded:
        return true
    default:
        return false
    }
}

// ValidateStatusTransition 验证状态转换（带错误信息）
func (s OrderStatus) ValidateStatusTransition(target OrderStatus) error {
    if !s.CanTransitionTo(target) {
        return fmt.Errorf("无效的状态转换: %s -> %s", s.GetStatusName(), target.GetStatusName())
    }
    return nil
}
```

---

#### 5.2 在订单服务中使用状态机
**文件**: `internal/service/order/order.go`
**问题**: 状态转换逻辑硬编码

**整改前**:
```go
func (s *OrderService) CancelOrder(ctx context.Context, userID uint64, orderID uint64, req CancelOrderRequest) error {
    order, err := s.orders.Get(ctx, orderID)
    if err != nil {
        return err
    }
    
    // 硬编码的状态检查
    if order.Status != model.OrderStatusPending && order.Status != model.OrderStatusConfirmed {
        return ErrInvalidTransition
    }
    
    order.Status = model.OrderStatusCanceled
    order.CancelReason = req.Reason
    
    return s.orders.Update(ctx, order)
}
```

**整改后**:
```go
func (s *OrderService) CancelOrder(ctx context.Context, userID uint64, orderID uint64, req CancelOrderRequest) error {
    order, err := s.orders.Get(ctx, orderID)
    if err != nil {
        return err
    }
    
    // 权限检查
    if order.UserID != userID {
        return ErrUnauthorized
    }
    
    // 使用状态机验证转换
    if err := order.Status.ValidateStatusTransition(model.OrderStatusCanceled); err != nil {
        return fmt.Errorf("%w: %s", ErrInvalidTransition, err.Error())
    }
    
    // 记录状态转换日志
    if err := s.logStatusTransition(ctx, order.ID, order.Status, model.OrderStatusCanceled, req.Reason); err != nil {
        logging.Warn("Failed to log status transition", "error", err, "order_id", order.ID)
    }
    
    order.Status = model.OrderStatusCanceled
    order.CancelReason = req.Reason
    canceledAt := time.Now()
    order.CompletedAt = &canceledAt
    
    return s.orders.Update(ctx, order)
}

// logStatusTransition 记录状态转换日志
type statusTransitionLog struct {
    OrderID     uint64                `json:"orderId"`
    FromStatus  model.OrderStatus     `json:"fromStatus"`
    ToStatus    model.OrderStatus     `json:"toStatus"`
    Reason      string                `json:"reason"`
    ActorUserID uint64                `json:"actorUserId"`
    Timestamp   time.Time             `json:"timestamp"`
}

func (s *OrderService) logStatusTransition(ctx context.Context, orderID uint64, from, to model.OrderStatus, reason string) error {
    actorUserID := getActorUserIDFromContext(ctx) // 从Context获取操作人ID
    
    log := statusTransitionLog{
        OrderID:     orderID,
        FromStatus:  from,
        ToStatus:    to,
        Reason:      reason,
        ActorUserID: actorUserID,
        Timestamp:   time.Now(),
    }
    
    // 记录到操作日志表
    operationLog := &model.OperationLog{
        EntityType: "order",
        EntityID:   orderID,
        Action:     "update_status",
        ActorUserID: &actorUserID,
        Reason:     fmt.Sprintf("状态转换: %s -> %s", from.GetStatusName(), to.GetStatusName()),
        MetadataJSON: mustMarshalJSON(log),
    }
    
    return s.operationLogRepo.Append(ctx, operationLog)
}

// mustMarshalJSON 安全的JSON序列化
func mustMarshalJSON(v interface{}) []byte {
    data, _ := json.Marshal(v)
    return data
}
```

---

### 6. 支付幂等性

#### 6.1 实现支付幂等性
**文件**: `internal/service/payment/payment.go`
**问题**: 重复提交可能导致重复支付

**整改内容**:
```go
package payment

import (
    "context"
    "fmt"
    "time"
    
    "gamelink/internal/model"
    "gamelink/internal/repository"
    "gamelink/internal/cache"
)

type PaymentService struct {
    payments repository.PaymentRepository
    orders   repository.OrderRepository
    cache    cache.Cache
}

type CreatePaymentRequest struct {
    OrderID         uint64 `json:"orderId" binding:"required"`
    Method          model.PaymentMethod `json:"method" binding:"required"`
    AmountCents     int64  `json:"amountCents" binding:"required,min=1"`
    IdempotencyKey  string `json:"idempotencyKey" binding:"required"` // 幂等性Key
    ClientIP        string `json:"-"` // 客户端IP（从Context获取）
}

func (s *PaymentService) CreatePayment(ctx context.Context, req CreatePaymentRequest) (*model.Payment, error) {
    // 验证订单存在且状态正确
    order, err := s.orders.Get(ctx, req.OrderID)
    if err != nil {
        return nil, fmt.Errorf("订单不存在: %w", err)
    }
    
    if order.Status != model.OrderStatusPending {
        return nil, fmt.Errorf("订单状态不正确，无法支付: %s", order.Status.GetStatusName())
    }
    
    // 验证支付金额
    if req.AmountCents != order.TotalPriceCents {
        return nil, fmt.Errorf("支付金额与订单金额不匹配: expected %d, got %d", 
            order.TotalPriceCents, req.AmountCents)
    }
    
    // 检查幂等性
    cacheKey := fmt.Sprintf("payment:idempotency:%s", req.IdempotencyKey)
    var existingPayment model.Payment
    
    if err := s.cache.Get(ctx, cacheKey, &existingPayment); err == nil {
        // 已处理过，直接返回
        logging.Info("Duplicate payment request detected", 
            "idempotency_key", req.IdempotencyKey,
            "payment_id", existingPayment.ID)
        return &existingPayment, nil
    }
    
    // 检查是否已有相同订单的待支付记录
    pendingPayments, err := s.payments.List(ctx, repository.PaymentListOptions{
        OrderID: &req.OrderID,
        Status:  &model.PaymentStatusPending,
    })
    if err != nil {
        return nil, err
    }
    
    if len(pendingPayments) > 0 {
        return nil, fmt.Errorf("订单已有待支付记录，请使用现有支付")
    }
    
    // 创建支付记录
    payment := &model.Payment{
        OrderID:         req.OrderID,
        UserID:          order.UserID,
        Method:          req.Method,
        AmountCents:     req.AmountCents,
        Currency:        order.Currency,
        Status:          model.PaymentStatusPending,
        IdempotencyKey:  req.IdempotencyKey,
        ClientIP:        req.ClientIP,
    }
    
    if err := s.payments.Create(ctx, payment); err != nil {
        return nil, fmt.Errorf("创建支付记录失败: %w", err)
    }
    
    // 写入幂等性缓存（24小时有效期）
    _ = s.cache.Set(ctx, cacheKey, payment, 24*time.Hour)
    
    // 发送支付确认消息到消息队列（异步处理）
    go s.processPaymentAsync(payment.ID)
    
    return payment, nil
}

// processPaymentAsync 异步处理支付
func (s *PaymentService) processPaymentAsync(paymentID uint64) {
    ctx := context.Background()
    
    // 模拟支付处理
    time.Sleep(2 * time.Second)
    
    payment, err := s.payments.Get(ctx, paymentID)
    if err != nil {
        logging.Error("Failed to get payment for processing", "error", err, "payment_id", paymentID)
        return
    }
    
    // 调用第三方支付网关（这里模拟成功）
    // ... 实际支付逻辑
    
    // 更新支付状态
    payment.Status = model.PaymentStatusPaid
    now := time.Now()
    payment.PaidAt = &now
    payment.ProviderTradeNo = generateTradeNo()
    
    if err := s.payments.Update(ctx, payment); err != nil {
        logging.Error("Failed to update payment status", "error", err, "payment_id", paymentID)
        return
    }
    
    // 更新订单状态
    order, err := s.orders.Get(ctx, payment.OrderID)
    if err != nil {
        logging.Error("Failed to get order for payment", "error", err, "order_id", payment.OrderID)
        return
    }
    
    order.Status = model.OrderStatusConfirmed
    if err := s.orders.Update(ctx, order); err != nil {
        logging.Error("Failed to update order status", "error", err, "order_id", order.ID)
        return
    }
    
    logging.Info("Payment processed successfully", "payment_id", paymentID, "order_id", order.ID)
}

// generateTradeNo 生成交易号
func generateTradeNo() string {
    return fmt.Sprintf("PAY%s%d", time.Now().Format("20060102150405"), time.Now().UnixNano())
}
```

---

### 7. 缓存策略实施

#### 7.1 实现Redis缓存层
**文件**: `internal/cache/redis_cache.go` (新建)
**问题**: 没有统一的缓存实现

**整改内容**:
```go
package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/redis/go-redis/v9"
    "gamelink/internal/logging"
)

// RedisCache Redis缓存实现
type RedisCache struct {
    client *redis.Client
    prefix string
}

// NewRedisCache 创建Redis缓存实例
func NewRedisCache(client *redis.Client, prefix string) *RedisCache {
    return &RedisCache{
        client: client,
        prefix: prefix,
    }
}

// Get 获取缓存值
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
    fullKey := c.prefix + key
    
    data, err := c.client.Get(ctx, fullKey).Bytes()
    if err != nil {
        if err == redis.Nil {
            return ErrNotFound
        }
        logging.Warn("Redis get failed", "key", fullKey, "error", err)
        return err
    }
    
    if err := json.Unmarshal(data, dest); err != nil {
        logging.Error("Failed to unmarshal cache value", "key", fullKey, "error", err)
        return err
    }
    
    return nil
}

// Set 设置缓存值
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    fullKey := c.prefix + key
    
    data, err := json.Marshal(value)
    if err != nil {
        logging.Error("Failed to marshal cache value", "key", fullKey, "error", err)
        return err
    }
    
    if err := c.client.Set(ctx, fullKey, data, ttl).Err(); err != nil {
        logging.Error("Redis set failed", "key", fullKey, "error", err)
        return err
    }
    
    return nil
}

// Delete 删除缓存
func (c *RedisCache) Delete(ctx context.Context, key string) error {
    fullKey := c.prefix + key
    return c.client.Del(ctx, fullKey).Err()
}

// Exists 检查缓存是否存在
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
    fullKey := c.prefix + key
    count, err := c.client.Exists(ctx, fullKey).Result()
    return count > 0, err
}

// ClearByPattern 按模式清除缓存
func (c *RedisCache) ClearByPattern(ctx context.Context, pattern string) error {
    fullPattern := c.prefix + pattern
    
    var cursor uint64
    for {
        keys, nextCursor, err := c.client.Scan(ctx, cursor, fullPattern, 100).Result()
        if err != nil {
            return err
        }
        
        if len(keys) > 0 {
            if err := c.client.Del(ctx, keys...).Err(); err != nil {
                logging.Error("Failed to delete keys", "keys", keys, "error", err)
            }
        }
        
        cursor = nextCursor
        if cursor == 0 {
            break
        }
    }
    
    return nil
}

// Ping 检查连接
func (c *RedisCache) Ping(ctx context.Context) error {
    return c.client.Ping(ctx).Err()
}

// Cache接口定义
type Cache interface {
    Get(ctx context.Context, key string, dest interface{}) error
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
    ClearByPattern(ctx context.Context, pattern string) error
    Ping(ctx context.Context) error
}

var ErrNotFound = fmt.Errorf("cache: key not found")
```

---

#### 7.2 在Service中集成缓存
**文件**: `internal/service/player/player.go`
**问题**: 陪玩师信息查询频繁，没有缓存

**整改内容**:
```go
package player

import (
    "context"
    "fmt"
    "time"
    
    "gamelink/internal/cache"
    "gamelink/internal/model"
    "gamelink/internal/repository"
)

type PlayerService struct {
    players repository.PlayerRepository
    users   repository.UserRepository
    cache   cache.Cache
}

func NewPlayerService(players repository.PlayerRepository, users repository.UserRepository, cache cache.Cache) *PlayerService {
    return &PlayerService{
        players: players,
        users:   users,
        cache:   cache,
    }
}

type PlayerProfileDTO struct {
    ID                uint64                       `json:"id"`
    UserID            uint64                       `json:"userId"`
    Nickname          string                       `json:"nickname"`
    Bio               string                       `json:"bio"`
    Rank              string                       `json:"rank"`
    RatingAverage     float32                      `json:"ratingAverage"`
    RatingCount       uint32                       `json:"ratingCount"`
    HourlyRateCents   int64                        `json:"hourlyRateCents"`
    MainGameID        uint64                       `json:"mainGameId"`
    VerificationStatus model.VerificationStatus    `json:"verificationStatus"`
    AvatarURL         string                       `json:"avatarUrl"`
    Name              string                       `json:"name"`
}

func (s *PlayerService) GetPlayerProfile(ctx context.Context, playerID uint64) (*PlayerProfileDTO, error) {
    cacheKey := fmt.Sprintf("player:profile:%d", playerID)
    
    // 尝试从缓存获取
    var dto PlayerProfileDTO
    if err := s.cache.Get(ctx, cacheKey, &dto); err == nil {
        return &dto, nil
    }
    
    // 查询数据库
    player, err := s.players.Get(ctx, playerID)
    if err != nil {
        return nil, err
    }
    
    user, err := s.users.Get(ctx, player.UserID)
    if err != nil {
        return nil, err
    }
    
    // 构建DTO
    dto = PlayerProfileDTO{
        ID:                player.ID,
        UserID:            player.UserID,
        Nickname:          player.Nickname,
        Bio:               player.Bio,
        Rank:              player.Rank,
        RatingAverage:     player.RatingAverage,
        RatingCount:       player.RatingCount,
        HourlyRateCents:   player.HourlyRateCents,
        MainGameID:        player.MainGameID,
        VerificationStatus: player.VerificationStatus,
        AvatarURL:         user.AvatarURL,
        Name:              user.Name,
    }
    
    // 缓存结果（5分钟）
    _ = s.cache.Set(ctx, cacheKey, dto, 5*time.Minute)
    
    return &dto, nil
}

func (s *PlayerService) UpdatePlayerProfile(ctx context.Context, playerID uint64, userID uint64, req UpdateProfileRequest) error {
    // 验证权限
    player, err := s.players.Get(ctx, playerID)
    if err != nil {
        return err
    }
    
    if player.UserID != userID {
        return ErrUnauthorized
    }
    
    // 更新玩家信息
    player.Nickname = req.Nickname
    player.Bio = req.Bio
    player.HourlyRateCents = req.HourlyRateCents
    
    if err := s.players.Update(ctx, player); err != nil {
        return err
    }
    
    // 清除缓存
    cacheKey := fmt.Sprintf("player:profile:%d", playerID)
    _ = s.cache.Delete(ctx, cacheKey)
    
    // 清除列表缓存
    _ = s.cache.ClearByPattern(ctx, "players:list:*")
    
    return nil
}

// UpdateProfileRequest 更新陪玩师资料请求
type UpdateProfileRequest struct {
    Nickname        string `json:"nickname" binding:"required,max=64"`
    Bio             string `json:"bio" binding:"max=500"`
    HourlyRateCents int64  `json:"hourlyRateCents" binding:"required,min=0"`
}
```

---

### 8. 代码复用提升

#### 8.1 实现泛型Repository
**文件**: `internal/repository/common/generic_repository.go` (新建)
**问题**: Repository代码重复严重

**整改内容**:
```go
package common

import (
    "context"
    "fmt"
    
    "gorm.io/gorm"
    "gamelink/internal/repository"
)

// GenericRepository 泛型仓储实现
type GenericRepository[T any] struct {
    db *gorm.DB
}

// NewGenericRepository 创建泛型仓储实例
func NewGenericRepository[T any](db *gorm.DB) *GenericRepository[T] {
    return &GenericRepository[T]{db: db}
}

// Create 创建记录
func (r *GenericRepository[T]) Create(ctx context.Context, entity *T) error {
    return r.db.WithContext(ctx).Create(entity).Error
}

// Get 根据ID获取记录
func (r *GenericRepository[T]) Get(ctx context.Context, id uint64) (*T, error) {
    var entity T
    if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, repository.ErrNotFound
        }
        return nil, err
    }
    return &entity, nil
}

// Update 更新记录
func (r *GenericRepository[T]) Update(ctx context.Context, entity *T) error {
    tx := r.db.WithContext(ctx).Model(entity).Updates(entity)
    if tx.Error != nil {
        return tx.Error
    }
    if tx.RowsAffected == 0 {
        return repository.ErrNotFound
    }
    return nil
}

// Delete 软删除记录
func (r *GenericRepository[T]) Delete(ctx context.Context, id uint64) error {
    var entity T
    tx := r.db.WithContext(ctx).Delete(&entity, id)
    if tx.Error != nil {
        return tx.Error
    }
    if tx.RowsAffected == 0 {
        return repository.ErrNotFound
    }
    return nil
}

// List 获取所有记录
func (r *GenericRepository[T]) List(ctx context.Context) ([]T, error) {
    var entities []T
    if err := r.db.WithContext(ctx).Find(&entities).Error; err != nil {
        return nil, err
    }
    return entities, nil
}

// Count 统计记录数
func (r *GenericRepository[T]) Count(ctx context.Context) (int64, error) {
    var count int64
    var entity T
    if err := r.db.WithContext(ctx).Model(&entity).Count(&count).Error; err != nil {
        return 0, err
    }
    return count, nil
}

// Exists 检查记录是否存在
func (r *GenericRepository[T]) Exists(ctx context.Context, id uint64) (bool, error) {
    var count int64
    var entity T
    
    if err := r.db.WithContext(ctx).
        Model(&entity).
        Where("id = ?", id).
        Count(&count).Error; err != nil {
        return false, err
    }
    
    return count > 0, nil
}

// QueryBuilder 查询构建器
type QueryBuilder[T any] struct {
    db    *gorm.DB
    query *gorm.DB
}

// NewQueryBuilder 创建查询构建器
func NewQueryBuilder[T any](db *gorm.DB) *QueryBuilder[T] {
    return &QueryBuilder[T]{
        db:    db,
        query: db.Model(new(T)),
    }
}

// Where 添加条件
func (qb *QueryBuilder[T]) Where(query interface{}, args ...interface{}) *QueryBuilder[T] {
    qb.query = qb.query.Where(query, args...)
    return qb
}

// Order 添加排序
func (qb *QueryBuilder[T]) Order(value interface{}) *QueryBuilder[T] {
    qb.query = qb.query.Order(value)
    return qb
}

// Limit 添加限制
func (qb *QueryBuilder[T]) Limit(limit int) *QueryBuilder[T] {
    qb.query = qb.query.Limit(limit)
    return qb
}

// Offset 添加偏移
func (qb *QueryBuilder[T]) Offset(offset int) *QueryBuilder[T] {
    qb.query = qb.query.Offset(offset)
    return qb
}

// Get 执行查询并获取结果
func (qb *QueryBuilder[T]) Get(ctx context.Context) ([]T, error) {
    var entities []T
    if err := qb.query.WithContext(ctx).Find(&entities).Error; err != nil {
        return nil, err
    }
    return entities, nil
}

// First 执行查询并获取第一条记录
func (qb *QueryBuilder[T]) First(ctx context.Context) (*T, error) {
    var entity T
    if err := qb.query.WithContext(ctx).First(&entity).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, repository.ErrNotFound
        }
        return nil, err
    }
    return &entity, nil
}

// Count 统计查询结果数量
func (qb *QueryBuilder[T]) Count(ctx context.Context) (int64, error) {
    var count int64
    if err := qb.query.WithContext(ctx).Count(&count).Error; err != nil {
        return 0, err
    }
    return count, nil
}
```

---

#### 8.2 使用泛型Repository重构UserRepository
**文件**: `internal/repository/user/repository.go`
**整改后**:
```go
package user

import (
    "context"
    "strings"
    "time"
    
    "gorm.io/gorm"
    "gamelink/internal/model"
    "gamelink/internal/repository"
    "gamelink/internal/repository/common"
)

// gormUserRepository 用户仓储实现
type gormUserRepository struct {
    *common.GenericRepository[model.User]
    db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) repository.UserRepository {
    return &gormUserRepository{
        GenericRepository: common.NewGenericRepository[model.User](db),
        db:                db,
    }
}

// ListWithFilters 带过滤条件的用户列表
func (r *gormUserRepository) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
    qb := common.NewQueryBuilder[model.User](r.db)
    
    // 添加过滤条件
    if len(opts.Roles) > 0 {
        qb = qb.Where("role IN ?", opts.Roles)
    }
    if len(opts.Statuses) > 0 {
        qb = qb.Where("status IN ?", opts.Statuses)
    }
    if opts.DateFrom != nil {
        qb = qb.Where("created_at >= ?", *opts.DateFrom)
    }
    if opts.DateTo != nil {
        qb = qb.Where("created_at <= ?", *opts.DateTo)
    }
    if kw := strings.TrimSpace(opts.Keyword); kw != "" {
        like := "%" + kw + "%"
        qb = qb.Where("name LIKE ? OR email LIKE ? OR phone LIKE ?", like, like, like)
    }
    
    // 统计总数
    total, err := qb.Count(ctx)
    if err != nil {
        return nil, 0, err
    }
    
    // 分页查询
    page := repository.NormalizePage(opts.Page)
    pageSize := repository.NormalizePageSize(opts.PageSize)
    offset := (page - 1) * pageSize
    
    users, err := qb.
        Order("created_at DESC").
        Limit(pageSize).
        Offset(offset).
        Get(ctx)
    
    return users, total, err
}

// GetByPhone 根据手机号获取用户
func (r *gormUserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
    qb := common.NewQueryBuilder[model.User](r.db).
        Where("phone = ?", phone)
    
    return qb.First(ctx)
}

// FindByEmail 根据邮箱获取用户
func (r *gormUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    qb := common.NewQueryBuilder[model.User](r.db).
        Where("email = ?", email)
    
    return qb.First(ctx)
}

// Update 更新用户（保持向后兼容）
func (r *gormUserRepository) Update(ctx context.Context, user *model.User) error {
    tx := r.db.WithContext(ctx).Model(user).Updates(map[string]any{
        "phone":         user.Phone,
        "email":         user.Email,
        "name":          user.Name,
        "avatar_url":    user.AvatarURL,
        "role":          user.Role,
        "status":        user.Status,
        "password_hash": user.PasswordHash,
        "last_login_at": user.LastLoginAt,
        "updated_at":    time.Now(),
    })
    
    if tx.Error != nil {
        return tx.Error
    }
    if tx.RowsAffected == 0 {
        return repository.ErrNotFound
    }
    
    return nil
}
```

---

## 🟢 低优先级整改任务

### 9. 监控指标添加

#### 9.1 集成Prometheus监控
**文件**: `internal/metrics/prometheus.go` (新建)
**问题**: 没有应用性能监控

**整改内容**:
```go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

// 定义指标
var (
    // HTTP请求计数
    HTTPRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "gamelink_http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    // HTTP请求耗时
    HTTPRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "gamelink_http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )
    
    // 数据库查询计数
    DatabaseQueriesTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "gamelink_database_queries_total",
            Help: "Total number of database queries",
        },
        []string{"table", "operation"},
    )
    
    // 数据库查询耗时
    DatabaseQueryDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "gamelink_database_query_duration_seconds",
            Help:    "Database query duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"table", "operation"},
    )
    
    // 缓存命中计数
    CacheHitsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "gamelink_cache_hits_total",
            Help: "Total number of cache hits",
        },
        []string{"cache", "operation"},
    )
    
    // 缓存未命中计数
    CacheMissesTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "gamelink_cache_misses_total",
            Help: "Total number of cache misses",
        },
        []string{"cache", "operation"},
    )
    
    // 活跃连接数
    ActiveConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "gamelink_active_connections",
            Help: "Number of active connections",
        },
    )
    
    // 订单状态分布
    OrdersByStatus = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "gamelink_orders_by_status",
            Help: "Number of orders by status",
        },
        []string{"status"},
    )
)
```

---

#### 9.2 在Middleware中记录指标
**文件**: `internal/handler/middleware/metrics.go` (新建)

**整改内容**:
```go
package middleware

import (
    "strconv"
    "time"
    
    "github.com/gin-gonic/gin"
    "gamelink/internal/metrics"
)

// PrometheusMetrics Prometheus指标中间件
func PrometheusMetrics() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        // 处理请求
        c.Next()
        
        // 记录指标
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())
        path := c.FullPath()
        method := c.Request.Method
        
        // HTTP请求计数
        metrics.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
        
        // HTTP请求耗时
        metrics.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
        
        // 活跃连接数
        metrics.ActiveConnections.Inc()
        defer metrics.ActiveConnections.Dec()
    }
}

// DatabaseMetrics 数据库指标记录
func DatabaseMetrics(table, operation string, duration time.Duration) {
    metrics.DatabaseQueriesTotal.WithLabelValues(table, operation).Inc()
    metrics.DatabaseQueryDuration.WithLabelValues(table, operation).Observe(duration.Seconds())
}

// CacheHit 记录缓存命中
func CacheHit(cache, operation string) {
    metrics.CacheHitsTotal.WithLabelValues(cache, operation).Inc()
}

// CacheMiss 记录缓存未命中
func CacheMiss(cache, operation string) {
    metrics.CacheMissesTotal.WithLabelValues(cache, operation).Inc()
}
```

---

### 10. 日志规范化

#### 10.1 实现结构化日志
**文件**: `internal/logging/logger.go` (重构)

**整改内容**:
```go
package logging

import (
    "context"
    "log/slog"
    "os"
    "runtime"
    "time"
)

var logger *slog.Logger

// InitLogger 初始化日志器
func InitLogger(env string) {
    var handler slog.Handler
    
    opts := &slog.HandlerOptions{
        Level: getLogLevel(),
        AddSource: env == "development",
    }
    
    if env == "development" {
        handler = slog.NewTextHandler(os.Stdout, opts)
    } else {
        handler = slog.NewJSONHandler(os.Stdout, opts)
    }
    
    logger = slog.New(handler)
}

// getLogLevel 获取日志级别
func getLogLevel() slog.Level {
    level := os.Getenv("LOG_LEVEL")
    switch level {
    case "debug":
        return slog.LevelDebug
    case "info":
        return slog.LevelInfo
    case "warn":
        return slog.LevelWarn
    case "error":
        return slog.LevelError
    default:
        return slog.LevelInfo
    }
}

// Debug 记录调试日志
func Debug(msg string, args ...interface{}) {
    log(slog.LevelDebug, msg, args...)
}

// Info 记录信息日志
func Info(msg string, args ...interface{}) {
    log(slog.LevelInfo, msg, args...)
}

// Warn 记录警告日志
func Warn(msg string, args ...interface{}) {
    log(slog.LevelWarn, msg, args...)
}

// Error 记录错误日志
func Error(msg string, args ...interface{}) {
    log(slog.LevelError, msg, args...)
}

// Fatal 记录致命错误并退出
func Fatal(msg string, args ...interface{}) {
    log(slog.LevelError, msg, args...)
    os.Exit(1)
}

// log 内部日志记录函数
func log(level slog.Level, msg string, args ...interface{}) {
    if logger == nil {
        InitLogger("development")
    }
    
    // 添加调用者信息
    if _, file, line, ok := runtime.Caller(2); ok {
        args = append(args, "file", file, "line", line)
    }
    
    logger.Log(context.Background(), level, msg, args...)
}

// WithContext 添加上下文信息
func WithContext(ctx context.Context, args ...interface{}) context.Context {
    return context.WithValue(ctx, logKey{}, args)
}

// FromContext 从上下文获取日志参数
func FromContext(ctx context.Context) []interface{} {
    if args, ok := ctx.Value(logKey{}).([]interface{}); ok {
        return args
    }
    return nil
}

type logKey struct{}

// HTTPRequestLog HTTP请求日志
func HTTPRequestLog(r *http.Request, statusCode int, duration time.Duration, args ...interface{}) {
    Info("http_request",
        append([]interface{}{
            "method", r.Method,
            "path", r.URL.Path,
            "query", r.URL.RawQuery,
            "status", statusCode,
            "duration_ms", duration.Milliseconds(),
            "client_ip", getClientIP(r),
            "user_agent", r.UserAgent(),
        }, args...)...,
    )
}

// getClientIP 获取客户端IP
func getClientIP(r *http.Request) string {
    if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
        return ip
    }
    if ip := r.Header.Get("X-Real-IP"); ip != "" {
        return ip
    }
    return r.RemoteAddr
}
```

---

### 11. 测试覆盖率提升

#### 11.1 补充用户服务测试
**文件**: `internal/service/user/user_test.go` (新建)

**整改内容**:
```go
package user

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "gamelink/internal/model"
    "gamelink/internal/repository/mocks"
)

// Mock定义
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

// 测试更新用户资料
func TestUpdateProfile(t *testing.T) {
    tests := []struct {
        name        string
        userID      uint64
        request     UpdateProfileRequest
        setupMock   func(*MockUserRepository)
        expectedErr error
    }{
        {
            name:   "成功更新",
            userID: 1,
            request: UpdateProfileRequest{
                Name:      "New Name",
                Bio:       "New bio",
                AvatarURL: "https://example.com/avatar.jpg",
            },
            setupMock: func(m *MockUserRepository) {
                m.On("Get", mock.Anything, uint64(1)).Return(&model.User{
                    ID:   1,
                    Name: "Old Name",
                }, nil)
                m.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
            },
            expectedErr: nil,
        },
        {
            name:   "用户不存在",
            userID: 999,
            request: UpdateProfileRequest{
                Name: "New Name",
            },
            setupMock: func(m *MockUserRepository) {
                m.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
            },
            expectedErr: repository.ErrNotFound,
        },
        {
            name:   "XSS攻击防护",
            userID: 1,
            request: UpdateProfileRequest{
                Name: "<script>alert('xss')</script>",
                Bio:  "<img src=x onerror=alert('xss')>",
            },
            setupMock: func(m *MockUserRepository) {
                m.On("Get", mock.Anything, uint64(1)).Return(&model.User{ID: 1}, nil)
                m.On("Update", mock.Anything, mock.MatchedBy(func(user *model.User) bool {
                    // 验证XSS被过滤
                    return user.Name == "" && user.Bio == ""
                })).Return(nil)
            },
            expectedErr: nil,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := new(MockUserRepository)
            tt.setupMock(mockRepo)
            
            svc := &UserService{
                userRepo: mockRepo,
            }
            
            err := svc.UpdateProfile(context.Background(), tt.userID, tt.request)
            
            if tt.expectedErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.expectedErr, err)
            } else {
                assert.NoError(t, err)
            }
            
            mockRepo.AssertExpectations(t)
        })
    }
}

// 性能测试
func BenchmarkUpdateProfile(b *testing.B) {
    mockRepo := new(MockUserRepository)
    mockRepo.On("Get", mock.Anything, uint64(1)).Return(&model.User{ID: 1}, nil)
    mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
    
    svc := &UserService{
        userRepo: mockRepo,
    }
    
    req := UpdateProfileRequest{
        Name:      "Test User",
        Bio:       "Test bio",
        AvatarURL: "https://example.com/avatar.jpg",
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = svc.UpdateProfile(context.Background(), 1, req)
    }
}
```

---

## 📊 整改效果预期

### 性能提升
- **响应时间**: 减少 40-60% (通过缓存和索引优化)
- **数据库负载**: 减少 50-70% (通过缓存和查询优化)
- **并发能力**: 提升 2-3倍 (通过连接池和优化)

### 安全性提升
- **XSS防护**: 100% 用户输入过滤
- **SQL注入**: 完全防止 (参数化查询)
- **认证安全**: JWT密钥管理规范化

### 可维护性提升
- **代码重复**: 减少 60% (泛型Repository)
- **测试覆盖**: 提升到 70%+
- **错误处理**: 100% 统一格式

---

## 🎯 整改实施计划

### 第一阶段 (1-2周): 高优先级整改
- [ ] JWT安全加固
- [ ] 输入验证增强
- [ ] 数据库索引优化
- [ ] 错误处理统一

### 第二阶段 (2-3周): 中优先级整改
- [ ] 订单状态机优化
- [ ] 支付幂等性实现
- [ ] Redis缓存集成
- [ ] 泛型Repository重构

### 第三阶段 (1-2周): 低优先级整改
- [ ] Prometheus监控集成
- [ ] 结构化日志
- [ ] 测试覆盖率提升
- [ ] 性能测试和调优

---

## ✅ 整改验证清单

### 功能验证
- [ ] 用户注册/登录正常
- [ ] 订单创建/支付/完成流程正常
- [ ] 权限控制有效
- [ ] 缓存命中和失效正常

### 性能验证
- [ ] 响应时间符合预期
- [ ] 数据库查询优化有效
- [ ] 缓存命中率 > 80%
- [ ] 并发测试通过

### 安全验证
- [ ] XSS攻击被过滤
- [ ] SQL注入被阻止
- [ ] JWT认证安全
- [ ] 权限控制严格

---

**整改完成日期**: _______________  
**验证人**: _______________  
**状态**: ⏳ 进行中 / ✅ 已完成
