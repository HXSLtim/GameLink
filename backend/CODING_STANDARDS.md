# 📚 GameLink 后端代码编写规范

**版本**: v1.0
**生效日期**: 2025年11月2日
**适用范围**: GameLink 后端项目全体开发人员

---

## 🎯 目录

1. [项目结构规范](#-项目结构规范)
2. [文件命名规范](#-文件命名规范)
3. [Go代码规范](#-go代码规范)
4. [测试规范](#-测试规范)
5. [Git提交规范](#-git提交规范)
6. [代码审查清单](#-代码审查清单)

---

## 🏗️ 项目结构规范

### 目录结构标准
```
backend/
├── cmd/                    # 应用程序入口点
│   └── main.go            # 主程序入口
├── internal/              # 内部包，不对外暴露
│   ├── model/             # 数据模型
│   ├── repository/        # 数据访问层
│   │   └── {domain}/repository.go
│   ├── service/           # 业务逻辑层
│   │   └── {domain}/{domain}.go
│   ├── handler/           # API处理层
│   │   ├── admin/{domain}.go
│   │   ├── user/{domain}.go
│   │   └── player/{domain}.go
│   ├── middleware/        # 中间件
│   ├── config/            # 配置管理
│   ├── db/                # 数据库相关
│   ├── cache/             # 缓存相关
│   └── utils/             # 工具函数
├── pkg/                   # 公共库，可对外暴露
├── docs/                  # 项目文档
├── configs/               # 配置文件
├── scripts/               # 脚本文件
├── tests/                 # 测试相关
└── go.mod                 # Go模块定义
```

### 架构分层原则
```
请求 → Handler → Service → Repository → Model → 数据库
```

**分层职责：**
- **Model**: 纯数据模型，无业务逻辑
- **Repository**: 纯数据访问，无业务逻辑
- **Service**: 业务逻辑，无HTTP处理
- **Handler**: HTTP处理，无业务逻辑

---

## 📝 文件命名规范

### 基本原则
- ✅ 使用小写字母和下划线
- ✅ 名称清晰表达功能
- ✅ 避免冗余后缀

### 文件命名标准

#### 模型文件
```go
✅ user.go          // 用户模型
✅ order.go         // 订单模型
✅ player.go        // 陪玩师模型

❌ user_model.go   // 冗余后缀
❌ User.go         // 不应使用大写开头
```

#### Repository文件
```go
✅ repository/user/repository.go
✅ repository/order/repository.go

❌ repository/user/user_repository.go    // 冗余后缀
❌ repository/user/user_gorm_repository.go // 技术细节前缀
```

#### Service文件
```go
✅ service/auth/auth.go
✅ service/order/order.go
✅ service/commission/commission.go

❌ service/auth/auth_service.go     // 冗余后缀
❌ service/order/order_service.go   // 冗余后缀
```

#### Handler文件
```go
✅ handler/admin/user.go
✅ handler/user/order.go
✅ handler/player/profile.go

❌ handler/admin/user_handler.go    // 冗余后缀
❌ handler/admin/admin_user.go      // 冗余前缀
```

#### 测试文件
```go
✅ user_test.go          // 对应 user.go
✅ repository_test.go    // 对应 repository.go
✅ auth_test.go          // 对应 auth.go

❌ test_user.go          // 不要前缀
❌ user_tests.go         // 不要复数形式
```

---

## 🐹 Go代码规范

### 包命名规范
```go
// ✅ 推荐：简洁、小写、有意义
package user
package order
package commission
package utils

// ❌ 避免：复杂或无意义的命名
package userservice
package repositoryimpl
package pkg
```

### 接口命名规范
```go
// ✅ 推荐：使用能力描述，以er结尾
type UserReader interface {
    GetUser(ctx context.Context, id uint64) (*User, error)
}

type OrderCreator interface {
    CreateOrder(ctx context.Context, order *Order) error
}

// ✅ 对于Repository可以使用具体名称
type UserRepository interface {
    // 方法定义
}
```

### 结构体命名规范
```go
// ✅ 推荐：大写开头，驼峰命名
type UserService struct {
    repo UserRepository
    cache Cache
}

type OrderRequest struct {
    UserID uint64 `json:"userId"`
    Amount int64  `json:"amount"`
}

// ✅ 私有结构体小写开头
type config struct {
    database DatabaseConfig
    redis    RedisConfig
}
```

### 方法命名规范
```go
// ✅ 推荐：动词+名词，驼峰命名
func (s *UserService) CreateUser(ctx context.Context, user *User) error
func (s *OrderService) ProcessPayment(ctx context.Context, orderID uint64) error
func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*User, error)

// ✅ 私有方法小写开头
func (s *UserService) validateUser(user *User) error
func (s *UserService) hashPassword(password string) string
```

### 变量命名规范
```go
// ✅ 推荐：小写驼峰命名
var userService *UserService
var orderRepository OrderRepository
var databaseConfig DatabaseConfig

// ✅ 常量使用大写+下划线
const (
    DefaultTimeout = 30 * time.Second
    MaxRetries     = 3
    UserStatusActive = "active"
)

// ✅ 短变量可以使用缩写
func (s *UserService) GetUser(id uint64) (*User, error) {
    var u User  // 短变量可以用缩写
    // ...
    return &u, nil
}
```

### 错误处理规范
```go
// ✅ 推荐：明确的错误处理
func (s *UserService) CreateUser(user *User) error {
    if err := s.validateUser(user); err != nil {
        return fmt.Errorf("user validation failed: %w", err)
    }

    if err := s.repo.Create(user); err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }

    return nil
}

// ✅ 使用自定义错误类型
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidInput = errors.New("invalid input")
    ErrPermissionDenied = errors.New("permission denied")
)
```

### 上下文传递规范
```go
// ✅ 推荐：所有公共方法都接收context
func (s *UserService) GetUser(ctx context.Context, id uint64) (*User, error)
func (s *OrderService) CreateOrder(ctx context.Context, order *Order) error
func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*User, error)

// ✅ 在函数间传递context
func (s *UserService) ComplexOperation(ctx context.Context, req *Request) error {
    // 传递context到其他方法
    user, err := s.repo.FindByID(ctx, req.UserID)
    if err != nil {
        return err
    }

    // 继续传递...
    return s.processOrder(ctx, user)
}
```

### 日志记录规范
```go
// ✅ 推荐：使用结构化日志
import "log/slog"

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    logger := slog.With("userId", user.ID, "email", user.Email)

    logger.Info("creating user")

    if err := s.repo.Create(user); err != nil {
        logger.Error("failed to create user", "error", err)
        return err
    }

    logger.Info("user created successfully")
    return nil
}
```

---

## 🧪 测试规范

### 测试文件组织
```go
// ✅ 推荐：测试文件与源文件同目录
// repository/user/repository.go
// repository/user/repository_test.go

// ✅ 包名保持一致
package user  // repository_test.go

// ✅ 导入必要的测试包
import (
    "testing"
    "context"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)
```

### 测试函数命名
```go
// ✅ 推荐：Test + 被测函数名 + 测试场景
func TestUserService_CreateUser(t *testing.T) {
    // 测试正常情况
}

func TestUserService_CreateUser_WithInvalidInput(t *testing.T) {
    // 测试异常情况
}

func TestUserService_CreateUser_WithDuplicateEmail(t *testing.T) {
    // 测试特定场景
}
```

### 测试结构模板
```go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange - 准备测试数据
    tests := []struct {
        name    string
        input   *User
        wantErr bool
    }{
        {
            name: "valid user",
            input: &User{
                Name:  "Test User",
                Email: "test@example.com",
            },
            wantErr: false,
        },
        {
            name: "invalid email",
            input: &User{
                Name:  "Test User",
                Email: "invalid-email",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Act - 执行被测方法
            err := service.CreateUser(context.Background(), tt.input)

            // Assert - 验证结果
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Mock使用规范
```go
// ✅ 推荐：使用接口和mock
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint64) (*User, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*User), args.Error(1)
}

// 测试中使用mock
func TestUserService_GetUser(t *testing.T) {
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)

    // 设置mock期望
    expectedUser := &User{ID: 1, Name: "Test"}
    mockRepo.On("FindByID", mock.Anything, uint64(1)).Return(expectedUser, nil)

    // 执行测试
    user, err := service.GetUser(context.Background(), 1)

    // 验证结果
    assert.NoError(t, err)
    assert.Equal(t, expectedUser, user)
    mockRepo.AssertExpectations(t)
}
```

---

## 🔄 Git提交规范

### 提交消息格式
```bash
# ✅ 推荐：<类型>(<范围>): <描述>

feat(service): 添加用户创建功能
fix(handler): 修复订单状态更新错误
docs(readme): 更新安装说明
style(code): 格式化代码
refactor(repository): 重构用户查询逻辑
test(user): 添加用户服务单元测试
chore(deps): 更新依赖包版本
```

### 提交类型说明
```bash
feat:     新功能
fix:      修复bug
docs:     文档更新
style:    代码格式化（不影响功能）
refactor: 重构代码（既不是新功能也不是修复）
test:     添加或修改测试
chore:    构建过程或辅助工具的变动
perf:     性能优化
ci:       CI/CD配置变更
```

### 提交消息示例
```bash
# ✅ 好的提交消息
feat(handler): 添加管理员用户列表API

- 实现用户分页查询
- 添加权限验证中间件
- 支持按用户名搜索
- 添加完整的单元测试

Closes #123

# ❌ 不好的提交消息
fix bug
update
add stuff
temp commit
```

---

## 👀 代码审查清单

### 功能性检查
- [ ] 代码实现了需求规格
- [ ] 边界条件处理正确
- [ ] 错误处理完善
- [ ] 性能考虑合理
- [ ] 安全性检查通过

### 代码质量检查
- [ ] 命名规范遵循
- [ ] 代码结构清晰
- [ ] 注释充分且准确
- [ ] 代码可读性好
- [ ] 避免代码重复

### 测试检查
- [ ] 单元测试覆盖主要功能
- [ ] 测试用例考虑边界情况
- [ ] 测试可以独立运行
- [ ] Mock使用合理
- [ ] 测试数据有意义

### 安全性检查
- [ ] 输入验证充分
- [ ] SQL注入防护
- [ ] XSS攻击防护
- [ ] 权限检查正确
- [ ] 敏感信息不泄露

### 性能检查
- [ ] 数据库查询优化
- [ ] 避免N+1查询
- [ ] 缓存使用合理
- [ ] 内存泄漏检查
- [ ] 并发安全考虑

---

## 📋 开发工作流程

### 1. 开发前准备
```bash
# 1. 创建功能分支
git checkout -b feature/user-management

# 2. 同步最新代码
git pull origin main

# 3. 安装依赖
go mod tidy
go mod download
```

### 2. 编码阶段
```bash
# 1. 编写代码（遵循规范）
# 2. 本地测试
go test ./...
go build ./...

# 3. 代码格式化
go fmt ./...
go vet ./...

# 4. 提交代码
git add .
git commit -m "feat(service): 添加用户管理功能"
```

### 3. 代码审查
```bash
# 1. 推送分支
git push origin feature/user-management

# 2. 创建Pull Request
# 3. 填写PR描述
# 4. 等待代码审查
# 5. 根据反馈修改代码
```

### 4. 合并代码
```bash
# 1. 通过审查后合并
git checkout main
git pull origin main
git branch -d feature/user-management
```

---

## 🔧 开发工具配置

### Go配置文件
```golangci.yml
# .golangci.yml
linters:
  enable:
    - gofmt
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - structcheck
    - varcheck
    - ineffassign
    - deadcode
    - typecheck
    - gosec
    - misspell
    - unconvert
    - dupl
    - goconst
    - gocyclo

linters-settings:
  gocyclo:
    min-complexity: 10
  dupl:
    threshold: 100
  goconst:
    min-len: 2
    min-occurrences: 2
```

### IDE配置
```json
// .vscode/settings.json
{
    "go.lintTool": "golangci-lint",
    "go.lintFlags": [
        "--fast"
    ],
    "go.testOnSave": true,
    "go.coverOnSave": true,
    "go.coverageDecorator": {
        "type": "gutter",
        "coveredHighlightColor": "rgba(64,128,64,0.5)",
        "uncoveredHighlightColor": "rgba(128,64,64,0.25)"
    },
    "go.useLanguageServer": true,
    "go.formatTool": "goimports"
}
```

---

## 📚 参考资源

### 官方文档
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Testing Documentation](https://golang.org/pkg/testing/)

### 最佳实践
- [Google Go Style Guide](https://google.github.io/styleguide/go/)
- [Uber Go Style Guide](https://github.com/uber-go/guide)
- [Clean Architecture in Go](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

### 工具推荐
- [golangci-lint](https://golangci-lint.run/) - Go代码检查工具
- [gomock](https://github.com/golang/mock) - Mock生成工具
- [swag](https://github.com/swaggo/swag) - Swagger文档生成
- [air](https://github.com/cosmtrek/air) - 热重载工具

---

## 🎯 总结

### 核心原则
1. **简洁性** - 代码应该简洁明了，避免过度设计
2. **可读性** - 代码应该易于理解和维护
3. **一致性** - 遵循统一的命名和结构规范
4. **测试性** - 代码应该易于测试
5. **安全性** - 始终考虑安全性问题

### 持续改进
- 定期review和更新代码规范
- 收集团队反馈，优化规范内容
- 跟进Go语言社区最佳实践
- 使用自动化工具确保规范执行

---

**本规范文档将随项目发展持续更新，请全体开发人员严格遵守并积极参与改进。** 🚀