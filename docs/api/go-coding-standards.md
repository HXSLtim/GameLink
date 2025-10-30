# Go 代码编写规范

## 📋 概述

本文档定义了 GameLink 项目中 Go 语言的代码编写规范，旨在保证代码质量、可维护性和团队协作效率。

## 🎯 核心原则

1. **简洁性**: 代码应该简洁明了，避免过度设计
2. **可读性**: 代码应该易于阅读和理解
3. **可维护性**: 代码应该易于修改和扩展
4. **性能**: 在保证代码质量的前提下追求性能
5. **一致性**: 团队成员应该遵循统一的编码风格

## 📝 命名规范

### 包命名
```go
// ✅ 好的包名 - 简短、清晰、全小写
package user
package order
package payment
package cache

// ❌ 避免的包名
package userservice     // 太长
package UserService     // 大写字母
package user_service    // 下划线
package util            // 太通用
```

### 变量和常量命名
```go
// ✅ 好的命名
var (
    userService    UserService
    orderCache     *redis.Client
    maxRetryCount  int = 3
)

const (
    DefaultTimeout = 30 * time.Second
    MaxPageSize    = 100
    APIVersion     = "v1"
)

// ❌ 避免的命名
var (
    us             UserService // 缩写不清晰
    cacheRedis     *redis.Client // 重复
    retry_count    int          // 下划线
)

// 常量使用驼峰命名或全大写+下划线
const (
    orderStatusPending    = "pending"
    OrderStatusProcessing = "processing"
    ORDER_STATUS_COMPLETED = "completed" // 不推荐，除非是导出的常量
)
```

### 函数和方法命名
```go
// ✅ 好的函数名 - 动词开头，清晰表达意图
func CreateUser(user *User) error
func GetOrderByID(id int64) (*Order, error)
func ValidatePayment(payment Payment) bool

// ✅ 好的方法名
func (u *UserService) Register(ctx context.Context, req *RegisterRequest) error
func (o *Order) UpdateStatus(status OrderStatus) error
func (r *OrderRepository) FindByID(ctx context.Context, id int64) (*Order, error)

// ❌ 避免的命名
func user(user *User) error           // 名字太短
func GetUser(user *User) error        // 参数命名不清晰
func Create() error                   // 创建什么？
func UpdateOrderStatusToCompleted()   // 太长
```

### 接口命名
```go
// ✅ 好的接口命名
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id int64) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int64) error
}

type PaymentService interface {
    ProcessPayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error)
    RefundPayment(ctx context.Context, paymentID string) error
}

// ✅ 接口名以 -er 结尾（当只有一个方法时）
type Reader interface {
    Read([]byte) (int, error)
}

type Writer interface {
    Write([]byte) (int, error)
}
```

### 结构体命名
```go
// ✅ 好的结构体命名 - 首字母大写，驼峰命名
type UserService struct {
    repo   UserRepository
    cache  CacheInterface
    logger Logger
}

type OrderRequest struct {
    UserID    int64     `json:"user_id" binding:"required"`
    GameID    int64     `json:"game_id" binding:"required"`
    Price     float64   `json:"price" binding:"required,gt=0"`
    StartTime time.Time `json:"start_time" binding:"required"`
}

// ✅ 私有结构体（小写开头）
type orderCache struct {
    client *redis.Client
    ttl    time.Duration
}
```

## 📚 注释规范

### 包注释
```go
// Package user 提供用户管理相关的功能，包括用户注册、登录、信息管理等。
package user

// Package order 实现订单管理系统，支持订单创建、状态跟踪、分发等功能。
package order
```

### 函数注释
```go
// CreateUser 创建新用户。
//
// 参数：
//   - ctx: 上下文，用于控制请求超时和取消
//   - user: 用户信息，包含必要字段
//
// 返回：
//   - error: 创建失败时返回错误信息
//
// 示例：
//   user := &User{Name: "张三", Phone: "13800138000"}
//   err := CreateUser(ctx, user)
func CreateUser(ctx context.Context, user *User) error {
    // 实现逻辑
}

// GetOrderByID 根据订单ID获取订单信息。
//
// 如果订单不存在，返回 ErrOrderNotFound 错误。
func GetOrderByID(ctx context.Context, id int64) (*Order, error) {
    // 实现逻辑
}
```

### 结构体和字段注释
```go
// UserService 用户服务，提供用户相关的业务逻辑。
type UserService struct {
    repo   UserRepository    // 用户仓储接口
    cache  CacheInterface    // 缓存接口
    logger Logger           // 日志记录器
    config *Config          // 配置信息
}

// OrderRequest 创建订单的请求参数。
type OrderRequest struct {
    UserID    int64     `json:"user_id" binding:"required"`    // 用户ID
    GameID    int64     `json:"game_id" binding:"required"`    // 游戏ID
    Price     float64   `json:"price" binding:"required,gt=0"` // 订单价格，必须大于0
    StartTime time.Time `json:"start_time" binding:"required"` // 开始时间
    GameLevel string    `json:"game_level"`                   // 游戏段位
}
```

## 🏗 代码组织

### 文件组织
```
user-service/
├── main.go              # 应用入口
├── config/
│   └── config.go        # 配置管理
├── handler/
│   ├── user_handler.go  # HTTP处理器
│   └── auth_handler.go  # 认证处理器
├── service/
│   ├── user_service.go  # 业务逻辑
│   └── auth_service.go  # 认证逻辑
├── repository/
│   ├── user_repo.go     # 数据访问层
│   └── user_repo_test.go # 测试文件
├── model/
│   ├── user.go          # 数据模型
│   └── user_request.go  # 请求模型
├── middleware/
│   ├── auth.go          # 认证中间件
│   └── cors.go          # CORS中间件
└── utils/
    ├── validator.go     # 验证工具
    └── response.go      # 响应工具
```

### 导入顺序
```go
import (
    // 标准库
    "context"
    "fmt"
    "time"

    // 第三方库
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "go.uber.org/zap"

    // 本项目包
    "github.com/gamelink/internal/config"
    "github.com/gamelink/internal/model"
    "github.com/gamelink/pkg/logger"
)
```

## 🔧 错误处理

### 错误定义
```go
// ✅ 自定义错误类型
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

// ✅ 使用 errors.New 和 fmt.Errorf
var (
    ErrUserNotFound    = errors.New("user not found")
    ErrInvalidPassword = errors.New("invalid password")
    ErrOrderExpired    = errors.New("order has expired")
)
```

### 错误处理最佳实践
```go
// ✅ 好的错误处理
func GetUser(ctx context.Context, id int64) (*User, error) {
    if id <= 0 {
        return nil, fmt.Errorf("invalid user id: %d", id)
    }

    user, err := userRepo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    if user == nil {
        return nil, ErrUserNotFound
    }

    return user, nil
}

// ✅ 使用错误包装
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) error {
    if err := s.validateRequest(req); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    if err := s.repo.Create(ctx, req.ToUser()); err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }

    return nil
}

// ❌ 避免的错误处理
func GetUser(id int64) (*User, error) {
    user, err := userRepo.GetByID(id)
    if err != nil {
        return nil, err // 没有包装错误信息
    }
    return user, nil
}
```

## 🎯 函数设计

### 函数长度
```go
// ✅ 好的函数 - 简短、单一职责
func (s *UserService) Register(ctx context.Context, req *RegisterRequest) error {
    // 验证请求参数
    if err := s.validateRegisterRequest(req); err != nil {
        return err
    }

    // 检查用户是否已存在
    if exists, err := s.repo.UserExists(ctx, req.Phone); err != nil {
        return err
    } else if exists {
        return ErrUserAlreadyExists
    }

    // 创建用户
    user := req.ToUser()
    if err := s.repo.Create(ctx, user); err != nil {
        return err
    }

    // 发送欢迎消息
    if err := s.sendWelcomeMessage(ctx, user); err != nil {
        s.logger.Warn("failed to send welcome message", zap.Error(err))
    }

    return nil
}

// ✅ 提取辅助函数
func (s *UserService) validateRegisterRequest(req *RegisterRequest) error {
    if req.Phone == "" {
        return ErrPhoneRequired
    }
    if !isValidPhone(req.Phone) {
        return ErrInvalidPhone
    }
    if len(req.Password) < 6 {
        return ErrPasswordTooShort
    }
    return nil
}
```

### 参数设计
```go
// ✅ 使用结构体作为参数（参数较多时）
type CreateOrderRequest struct {
    UserID    int64     `json:"user_id" binding:"required"`
    GameID    int64     `json:"game_id" binding:"required"`
    Price     float64   `json:"price" binding:"required,gt=0"`
    StartTime time.Time `json:"start_time" binding:"required"`
    Duration  int       `json:"duration" binding:"required,min=1"`
}

func (s *OrderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*Order, error) {
    // 实现
}

// ✅ 第一个参数总是 context.Context
func (s *UserService) GetByID(ctx context.Context, id int64) (*User, error) {
    // 实现
}

// ✅ 返回值顺序：(result, error)
func (s *UserService) ValidateUser(ctx context.Context, user *User) (bool, error) {
    // 实现
}
```

## 🏛 并发编程

### Goroutine 使用
```go
// ✅ 好的 Goroutine 使用
func (s *OrderService) ProcessOrders(ctx context.Context, orders []Order) error {
    var wg sync.WaitGroup
    errChan := make(chan error, len(orders))

    // 限制并发数量
    semaphore := make(chan struct{}, 10)

    for _, order := range orders {
        wg.Add(1)
        go func(o Order) {
            defer wg.Done()

            semaphore <- struct{}{} // 获取信号量
            defer func() { <-semaphore }() // 释放信号量

            if err := s.processSingleOrder(ctx, o); err != nil {
                errChan <- err
            }
        }(order)
    }

    wg.Wait()
    close(errChan)

    // 检查是否有错误
    for err := range errChan {
        if err != nil {
            return err
        }
    }

    return nil
}

// ✅ 使用 context 控制 Goroutine 生命周期
func (s *NotificationService) StartWorker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            s.logger.Info("worker stopped")
            return
        case notification := <-s.notificationChan:
            if err := s.sendNotification(ctx, notification); err != nil {
                s.logger.Error("failed to send notification", zap.Error(err))
            }
        }
    }
}
```

### Channel 使用
```go
// ✅ 定义 Channel 类型
type OrderProcessor struct {
    orderChan   chan Order
    resultChan  chan ProcessResult
    workerCount int
}

// ✅ 使用 buffered channel 提高性能
func NewOrderProcessor(bufferSize, workerCount int) *OrderProcessor {
    return &OrderProcessor{
        orderChan:   make(chan Order, bufferSize),
        resultChan:  make(chan ProcessResult, bufferSize),
        workerCount: workerCount,
    }
}

// ✅ 正确关闭 Channel
func (p *OrderProcessor) Stop() {
    close(p.orderChan)  // 关闭输入 channel
    p.wg.Wait()         // 等待所有 worker 完成
    close(p.resultChan) // 关闭输出 channel
}
```

## 🧪 测试规范

### 单元测试
```go
// ✅ 好的测试命名
func TestUserService_CreateUser_Success(t *testing.T) {
    // 准备测试数据
    ctx := context.Background()
    req := &CreateUserRequest{
        Phone:    "13800138000",
        Password: "password123",
        Name:     "张三",
    }

    // 创建 mock 对象
    mockRepo := &MockUserRepository{}
    mockCache := &MockCache{}

    service := NewUserService(mockRepo, mockCache)

    // 设置 mock 期望
    mockRepo.On("UserExists", ctx, req.Phone).Return(false, nil)
    mockRepo.On("Create", ctx, mock.AnythingOfType("*User")).Return(nil)
    mockCache.On("Set", ctx, mock.Anything, mock.Anything).Return(nil)

    // 执行测试
    err := service.CreateUser(ctx, req)

    // 验证结果
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}

// ✅ 测试错误情况
func TestUserService_CreateUser_UserAlreadyExists(t *testing.T) {
    ctx := context.Background()
    req := &CreateUserRequest{
        Phone:    "13800138000",
        Password: "password123",
        Name:     "张三",
    }

    mockRepo := &MockUserRepository{}
    mockCache := &MockCache{}

    service := NewUserService(mockRepo, mockCache)

    // 设置 mock：用户已存在
    mockRepo.On("UserExists", ctx, req.Phone).Return(true, nil)

    err := service.CreateUser(ctx, req)

    assert.Error(t, err)
    assert.Equal(t, ErrUserAlreadyExists, err)
    mockRepo.AssertExpectations(t)
}
```

### 集成测试
```go
// ✅ 集成测试示例
func TestOrderService_Integration(t *testing.T) {
    // 设置测试环境
    db := setupTestDB(t)
    redis := setupTestRedis(t)
    defer cleanupTestDB(db)
    defer cleanupTestRedis(redis)

    // 创建服务实例
    repo := NewOrderRepository(db)
    cache := NewOrderCache(redis)
    service := NewOrderService(repo, cache)

    // 测试完整流程
    ctx := context.Background()

    // 创建订单
    order, err := service.CreateOrder(ctx, &CreateOrderRequest{
        UserID:    1,
        GameID:    1,
        Price:     100.0,
        StartTime: time.Now(),
        Duration:  60,
    })
    assert.NoError(t, err)
    assert.NotNil(t, order)

    // 获取订单
    retrieved, err := service.GetOrderByID(ctx, order.ID)
    assert.NoError(t, err)
    assert.Equal(t, order.ID, retrieved.ID)

    // 更新订单状态
    err = service.UpdateOrderStatus(ctx, order.ID, OrderStatusConfirmed)
    assert.NoError(t, err)
}
```

### 基准测试
```go
// ✅ 基准测试示例
func BenchmarkUserService_CreateUser(b *testing.B) {
    service := setupUserService(b)
    ctx := context.Background()

    req := &CreateUserRequest{
        Phone:    "13800138000",
        Password: "password123",
        Name:     "张三",
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // 每次使用不同的手机号避免冲突
        req.Phone = fmt.Sprintf("138001380%04d", i)
        err := service.CreateUser(ctx, req)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## 📊 性能优化

### 内存优化
```go
// ✅ 使用对象池减少内存分配
var userPool = sync.Pool{
    New: func() interface{} {
        return &User{}
    },
}

func (s *UserService) ProcessUsers(users []*User) error {
    for _, user := range users {
        // 从池中获取对象
        processedUser := userPool.Get().(*User)
        defer userPool.Put(processedUser)

        // 复制数据
        *processedUser = *user

        // 处理逻辑
        if err := s.processUser(processedUser); err != nil {
            return err
        }
    }
    return nil
}

// ✅ 避免不必要的内存分配
func (s *UserService) ValidateUsers(users []User) []error {
    // 预分配切片容量
    errors := make([]error, 0, len(users))

    for _, user := range users {
        if err := s.validateUser(&user); err != nil {
            errors = append(errors, err)
        }
    }
    return errors
}
```

### 数据库优化
```go
// ✅ 批量操作
func (r *OrderRepository) CreateOrders(ctx context.Context, orders []Order) error {
    const batchSize = 100

    for i := 0; i < len(orders); i += batchSize {
        end := i + batchSize
        if end > len(orders) {
            end = len(orders)
        }

        batch := orders[i:end]
        if err := r.createBatch(ctx, batch); err != nil {
            return fmt.Errorf("failed to create batch %d-%d: %w", i, end, err)
        }
    }
    return nil
}

// ✅ 使用连接池
func (r *OrderRepository) GetOrdersWithUser(ctx context.Context, orderIDs []int64) ([]OrderWithUser, error) {
    query := `
        SELECT o.*, u.name as user_name, u.phone as user_phone
        FROM orders o
        LEFT JOIN users u ON o.user_id = u.id
        WHERE o.id IN (?)
    `

    var results []OrderWithUser
    if err := r.db.WithContext(ctx).Raw(query, orderIDs).Scan(&results).Error; err != nil {
        return nil, fmt.Errorf("failed to get orders with user: %w", err)
    }

    return results, nil
}
```

## 🛠 工具配置

### golangci-lint 配置
```yaml
# .golangci.yml
run:
  timeout: 5m
  tests: true

linters-settings:
  govet:
    check-shadowing: true
  golint:
    min-confidence: 0
  gocyclo:
    min-complexity: 15
  maligned:
    suggest-new: true
  dupl:
    threshold: 100
  goconst:
    min-len: 2
    min-occurrences: 2

linters:
  enable:
    - bodyclose
    - deadcode
    - depguard
    - dogsled
    - dupl
    - errcheck
    - gochecknoinits
    - goconst
    - gocritic
    - gocyclo
    - gofmt
    - goimports
    - golint
    - gomnd
    - goprintffuncname
    - gosec
    - gosimple
    - govet
    - ineffassign
    - interfacer
    - lll
    - misspell
    - nakedret
    - rowserrcheck
    - scopelint
    - staticcheck
    - structcheck
    - stylecheck
    - typecheck
    - unconvert
    - unparam
    - unused
    - varcheck
    - whitespace
```

### pre-commit 配置
```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.4.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-added-large-files

  - repo: local
    hooks:
      - id: golangci-lint
        name: golangci-lint
        entry: golangci-lint run --timeout=5m
        language: system
        types: [go]
        pass_filenames: false

      - id: go-test
        name: go test
        entry: go test ./...
        language: system
        types: [go]
        pass_filenames: false

      - id: go-mod-tidy
        name: go mod tidy
        entry: go mod tidy
        language: system
        files: (go\.mod|go\.sum)$
        pass_filenames: false
```

## 📚 参考资料

- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go官方文档](https://golang.org/doc/)
- [Go语言最佳实践](https://github.com/golang/go/wiki/LearnServerProgramming)

---

遵循这些规范将帮助我们创建高质量、可维护的 Go 代码。如有疑问，请与团队讨论并持续改进规范。