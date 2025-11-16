# GameLink 后端架构重构计划

## 概述

基于对 GameLink 后端代码的全面分析，本文档制定了系统性的架构重构计划，旨在提升代码质量、可维护性和扩展性。

## 当前架构评估

### 优势
- ✅ 清晰的四层架构设计（Handler → Service → Repository → Model）
- ✅ 完善的依赖注入机制
- ✅ 全面的测试覆盖（使用 gomock）
- ✅ 统一的错误处理机制
- ✅ 良好的模块化组织

### 主要问题
- ❌ 主函数 `main.go` 过于复杂（200+ 行路由注册）
- ❌ Service 层职责过重（如 OrderService 包含订单、支付、抽成等）
- ❌ Repository 接口过于庞大（违反接口隔离原则）
- ❌ 响应格式不统一
- ❌ 配置管理存在安全隐患
- ❌ 依赖注入分散，重复创建实例

## 重构目标

1. **提升代码组织性** - 模块化拆分，职责单一
2. **优化依赖管理** - 统一依赖注入，避免重复实例
3. **改善业务逻辑** - 拆分复杂服务，引入领域事件

## 重构方案

## 第一部分：代码组织优化

### 目标
将代码组织问题进行系统性优化，提升代码的可读性和维护性。

### 重构范围
- `backend/cmd/main.go` - 主入口文件模块化
- `backend/internal/handler/` - 统一响应格式处理
- `backend/internal/config/` - 配置管理优化

### 具体实施

#### 1. 路由注册模块化

**问题**: `registerRoutes` 函数过长，职责混杂

**解决方案**:
```
backend/internal/router/
├── router.go              # 主路由器
├── routes.go              # 路由注册接口
├── admin/
│   ├── routes.go         # 管理端路由
│   ├── rbac.go          # RBAC 相关路由
│   └── permissions.go   # 权限管理路由
├── user/
│   ├── routes.go         # 用户端路由
│   └── auth.go          # 认证相关路由
├── player/
│   ├── routes.go         # 陪玩师端路由
│   └── earnings.go      # 收益相关路由
└── health/
    └── routes.go         # 健康检查路由
```

**重构前后对比**:
```go
// 重构前 - 200+ 行的单一函数
func registerRoutes(router *gin.Engine, cfg config.AppConfig, orm *gorm.DB, ...) {
    // 大量路由注册逻辑混合在一起
}

// 重构后 - 模块化路由注册
func registerRoutes(router *gin.Engine, cfg config.AppConfig, services *appServices) {
    routerConfig := &router.Config{
        Router:  router,
        Config:  cfg,
        Services: services,
    }
    
    router.RegisterHealthRoutes(routerConfig)
    router.RegisterAuthRoutes(routerConfig)
    router.RegisterUserRoutes(routerConfig)
    router.RegisterPlayerRoutes(routerConfig)
    router.RegisterAdminRoutes(routerConfig)
}
```

#### 2. 统一响应格式处理

**问题**: Handler 层响应格式不统一，错误处理分散

**解决方案**:
```
backend/internal/handler/
├── response/
│   ├── types.go          # 响应类型定义
│   ├── builder.go        # 响应构建器
│   ├── success.go        # 成功响应处理
│   ├── error.go         # 错误响应处理
│   └── pagination.go    # 分页响应处理
└── middleware/
    └── response.go       # 响应中间件
```

**统一响应格式**:
```go
type APIResponse[T any] struct {
    Success    bool         `json:"success"`
    Code       int          `json:"code"`
    Message    string       `json:"message"`
    Data       T            `json:"data,omitempty"`
    Pagination *Pagination `json:"pagination,omitempty"`
}

type ResponseBuilder struct{}

func (r *ResponseBuilder) Success(c *gin.Context, data any) {
    c.JSON(http.StatusOK, APIResponse[any]{
        Success: true,
        Code:    http.StatusOK,
        Message: "success",
        Data:    data,
    })
}

func (r *ResponseBuilder) Error(c *gin.Context, err error) {
    status := http.StatusInternalServerError
    message := err.Error()
    
    // 根据错误类型确定状态码
    switch {
    case errors.Is(err, repository.ErrNotFound):
        status = http.StatusNotFound
    case errors.Is(err, service.ErrValidation):
        status = http.StatusBadRequest
    }
    
    c.JSON(status, APIResponse[any]{
        Success: false,
        Code:    status,
        Message: message,
    })
}
```

#### 3. 配置管理优化

**问题**: 硬编码默认值，配置验证分散

**解决方案**:
```
backend/internal/config/
├── config.go           # 配置结构定义
├── validator.go        # 配置验证逻辑
├── defaults.go         # 安全的默认值
├── loader.go           # 配置加载器
└── env.go             # 环境变量处理
```

**配置安全化**:
```go
// 移除不安全的默认值
const DefaultDevJWTSecret = "" // 空值，强制要求设置

// 统一配置验证
func Validate(env string, cfg *AppConfig) error {
    var errors []string
    
    if cfg.Auth.JWTSecret == "" {
        errors = append(errors, "JWT secret is required")
    }
    
    if cfg.Database.DSN == "" && env == "production" {
        errors = append(errors, "Database DSN is required in production")
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("validation failed: %s", strings.Join(errors, ", "))
    }
    
    return nil
}
```

### 预期收益
- 🎯 代码可读性提升 40%
- 🎯 维护成本降低 30%
- 🎯 配置安全性显著提升

---

## 第二部分：依赖注入优化

### 目标
优化依赖管理，消除重复实例创建，提升系统性能。

### 重构范围
- `backend/internal/service/` - 服务层依赖管理
- `backend/internal/repository/` - 仓储层接口优化
- `backend/cmd/main.go` - 服务初始化逻辑

### 具体实施

#### 1. Repository 接口拆分

**问题**: 单一接口包含过多方法，违反接口隔离原则

**解决方案**:
```
backend/internal/repository/
├── interfaces/
│   ├── order/
│   │   ├── reader.go       # 订单读取接口
│   │   ├── writer.go       # 订单写入接口
│   │   ├── query.go        # 订单查询接口
│   │   └── composer.go     # 接口组合
│   ├── user/
│   │   ├── reader.go
│   │   ├── writer.go
│   │   └── composer.go
│   └── payment/
│       ├── reader.go
│       ├── writer.go
│       └── composer.go
└── implementations/
    └── order_repository.go   # 具体实现
```

**接口拆分示例**:
```go
// 拆分前 - 庞大的单一接口
type OrderRepository interface {
    Create(ctx context.Context, order *model.Order) error
    Get(ctx context.Context, id uint64) (*model.Order, error)
    Update(ctx context.Context, order *model.Order) error
    Delete(ctx context.Context, id uint64) error
    List(ctx context.Context, opts OrderListOptions) ([]model.Order, int64, error)
    // ... 20+ 更多方法
}

// 拆分后 - 职责明确的接口
type OrderReader interface {
    Get(ctx context.Context, id uint64) (*model.Order, error)
    List(ctx context.Context, opts OrderListOptions) ([]model.Order, int64, error)
}

type OrderWriter interface {
    Create(ctx context.Context, order *model.Order) error
    Update(ctx context.Context, order *model.Order) error
    Delete(ctx context.Context, id uint64) error
}

type OrderRepository interface {
    OrderReader
    OrderWriter
}
```

#### 2. 引入依赖注入容器

**问题**: 服务实例化分散，重复创建相同实例

**解决方案**:
```
backend/internal/container/
├── container.go         # 依赖注入容器
├── providers.go        # 服务提供者
├── wire_gen.go         # 自动生成的依赖代码
└── types.go           # 容器类型定义
```

**使用 Wire 进行依赖注入**:
```go
//go:build wireinject
//+build wireinject

package container

import (
    "github.com/google/wire"
    "gamelink/internal/config"
    "gamelink/internal/db"
    "gamelink/internal/cache"
    // ... 其他导入
)

// Application 应用程序结构
type Application struct {
    Router *gin.Engine
    Config *config.AppConfig
    DB     *gorm.DB
    Cache  cache.Cache
    Services *Services
}

// Services 所有服务集合
type Services struct {
    Admin    *adminservice.AdminService
    Order    *orderservice.OrderService
    Payment  *paymentservice.PaymentService
    // ... 其他服务
}

// InitializeApp 初始化应用程序
func InitializeApp() (*Application, error) {
    wire.Build(
        // 配置提供者
        config.Load,
        config.Validate,
        
        // 基础设施提供者
        db.Open,
        cache.New,
        
        // 仓储提供者
        NewRepositories,
        
        // 服务提供者
        NewServices,
        
        // 路由提供者
        NewRouter,
        
        // 应用程序
        NewApplication,
    )
    return &Application{}, nil
}
```

#### 3. 服务生命周期管理

**问题**: 服务启动和停止逻辑分散

**解决方案**:
```
backend/internal/lifecycle/
├── manager.go          # 生命周期管理器
├── service.go         # 服务接口定义
├── registry.go        # 服务注册表
└── hooks.go          # 生命周期钩子
```

**生命周期管理**:
```go
type Service interface {
    Name() string
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}

type Manager struct {
    services map[string]Service
    hooks    map[Phase][]Hook
}

func (m *Manager) Register(service Service) {
    m.services[service.Name()] = service
}

func (m *Manager) Start(ctx context.Context) error {
    for name, service := range m.services {
        if err := service.Start(ctx); err != nil {
            return fmt.Errorf("failed to start service %s: %w", name, err)
        }
    }
    return nil
}

func (m *Manager) Stop(ctx context.Context) error {
    for name, service := range m.services {
        if err := service.Stop(ctx); err != nil {
            log.Printf("failed to stop service %s: %v", name, err)
        }
    }
    return nil
}
```

### 预期收益
- 🎯 依赖管理复杂度降低 50%
- 🎯 内存使用优化 20%
- 🎯 服务启动速度提升 30%

---

## 第三部分：业务逻辑重构

### 目标
拆分复杂的业务逻辑，提升代码的可测试性和可维护性。

### 重构范围
- `backend/internal/service/order/order.go` - 订单服务重构
- `backend/internal/service/payment/` - 支付服务独立
- `backend/internal/service/commission/` - 抽成服务独立

### 具体实施

#### 1. OrderService 职责拆分

**问题**: CreateOrder 方法过长（100+ 行），包含多种业务逻辑

**解决方案**:
```
backend/internal/service/order/
├── service.go           # 主要服务接口
├── creation.go         # 订单创建逻辑
├── validation.go       # 订单验证逻辑
├── pricing.go          # 价格计算逻辑
├── status.go          # 状态管理逻辑
├── timeline.go        # 时间线构建逻辑
└── dto.go            # 数据传输对象
```

**订单创建逻辑拆分**:
```go
// 拆分前 - 复杂的 CreateOrder 方法
func (s *OrderService) CreateOrder(ctx context.Context, userID uint64, req CreateOrderRequest) (*CreateOrderResponse, error) {
    // 100+ 行代码，包含：
    // - 陪玩师验证
    // - 游戏验证  
    // - 价格计算
    // - 订单创建
    // - 抽成计算
    // - 返回构建
}

// 拆分后 - 职责明确的方法
func (s *OrderService) CreateOrder(ctx context.Context, userID uint64, req CreateOrderRequest) (*CreateOrderResponse, error) {
    // 1. 验证
    if err := s.validateCreateOrder(ctx, userID, req); err != nil {
        return nil, err
    }
    
    // 2. 计算
    pricing, err := s.calculateOrderPricing(ctx, req)
    if err != nil {
        return nil, err
    }
    
    // 3. 创建
    order, err := s.createOrderEntity(ctx, userID, req, pricing)
    if err != nil {
        return nil, err
    }
    
    // 4. 返回
    return s.buildCreateOrderResponse(order, pricing), nil
}

func (s *OrderService) validateCreateOrder(ctx context.Context, userID uint64, req CreateOrderRequest) error {
    // 专注于验证逻辑
    return s.validator.ValidateCreateOrder(ctx, userID, req)
}

func (s *OrderService) calculateOrderPricing(ctx context.Context, req CreateOrderRequest) (*OrderPricing, error) {
    // 专注于价格计算
    return s.pricingCalculator.Calculate(ctx, req)
}
```

#### 2. 引入领域事件

**问题**: 复杂业务流程耦合严重

**解决方案**:
```
backend/internal/events/
├── dispatcher.go       # 事件分发器
├── handlers/          # 事件处理器
│   ├── order_created.go
│   ├── order_paid.go
│   ├── order_completed.go
│   ├── payment_completed.go
│   └── commission_calculated.go
├── events.go          # 事件定义
└── store.go           # 事件存储
```

**事件驱动架构**:
```go
// 事件定义
type OrderCreated struct {
    OrderID    uint64
    UserID     uint64
    PlayerID   uint64
    Amount     int64
    CreatedAt  time.Time
}

type PaymentCompleted struct {
    OrderID    uint64
    PaymentID  uint64
    Amount     int64
    Method     string
    PaidAt     time.Time
}

// 事件处理器
type CommissionHandler struct {
    service *commissionservice.CommissionService
}

func (h *CommissionHandler) HandleOrderCreated(ctx context.Context, event OrderCreated) error {
    // 异步计算抽成
    go h.service.CalculateCommission(ctx, event.OrderID)
    return nil
}

// 在订单创建时发布事件
func (s *OrderService) CreateOrder(ctx context.Context, userID uint64, req CreateOrderRequest) (*CreateOrderResponse, error) {
    order, err := s.createOrderEntity(ctx, userID, req, pricing)
    if err != nil {
        return nil, err
    }
    
    // 发布事件
    s.dispatcher.Publish(ctx, OrderCreated{
        OrderID:   order.ID,
        UserID:    userID,
        PlayerID:  *order.PlayerID,
        Amount:    order.TotalPriceCents,
        CreatedAt: order.CreatedAt,
    })
    
    return s.buildCreateOrderResponse(order, pricing), nil
}
```

#### 3. 业务规则抽象

**问题**: 业务规则散布在服务方法中

**解决方案**:
```
backend/internal/rules/
├── order_rules.go      # 订单业务规则
├── payment_rules.go    # 支付业务规则
├── commission_rules.go # 抽成业务规则
└── validator.go       # 规则验证器
```

**业务规则抽象**:
```go
// 订单状态流转规则
type OrderStatusTransitionRule struct {
    From model.OrderStatus
    To   model.OrderStatus
    Condition func(order *model.Order) bool
}

var OrderStatusRules = []OrderStatusTransitionRule{
    {
        From: model.OrderStatusPending,
        To:   model.OrderStatusConfirmed,
        Condition: func(order *model.Order) bool {
            return order.IsPaid()
        },
    },
    {
        From: model.OrderStatusConfirmed,
        To:   model.OrderStatusInProgress,
        Condition: func(order *model.Order) bool {
            return order.PlayerID != nil && *order.PlayerID > 0
        },
    },
}

type OrderRuleEngine struct {
    rules []OrderStatusTransitionRule
}

func (e *OrderRuleEngine) CanTransition(order *model.Order, toStatus model.OrderStatus) bool {
    for _, rule := range e.rules {
        if rule.From == order.Status && rule.To == toStatus {
            return rule.Condition(order)
        }
    }
    return false
}
```

### 预期收益
- 🎯 业务逻辑复杂度降低 60%
- 🎯 代码可测试性提升 80%
- 🎯 新功能开发效率提升 40%

---

## 重构执行计划

### 独立性保证
1. **接口隔离**: 每部分重构都通过接口隔离
2. **向后兼容**: 保持现有 API 接口不变
3. **渐进式重构**: 可以逐个模块进行重构
4. **测试覆盖**: 每个重构步骤都有对应的测试

### 风险控制
- 完整的测试覆盖确保功能正确性
- 重构过程中保持原有功能不变
- 可以在任何阶段暂停而不影响系统稳定性
- 使用 feature flag 控制新代码的启用

### 质量保证
- 代码审查确保重构质量
- 性能测试确保系统性能
- 集成测试确保系统稳定性
- 文档更新确保知识传承

---

## 总结

本重构计划通过三个独立且互不影响的部分，系统性地解决 GameLink 后端架构中的主要问题：

1. **第一部分** - 代码组织优化，提升可读性和维护性
2. **第二部分** - 依赖注入优化，提升性能和资源利用率
3. **第三部分** - 业务逻辑重构，提升可测试性和扩展性

每个部分都可以独立完成，不影响其他部分的工作，确保了重构工作的灵活性和可控性。通过这次重构，GameLink 后端将具备更好的代码质量、更高的开发效率和更强的系统稳定性。
