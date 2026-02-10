# GameLink 单元测试指南

## 目标

提升代码质量和可测试性，通过 Mock Repository 实现 Service 层的独立测试（无需数据库）。

## 架构优势

当前项目已具备良好的测试基础：

✅ **Repository 接口完整定义** (`api/internal/repository/interfaces.go`)  
✅ **Service 层依赖接口注入**  
✅ **Mock Repository 已存在** (`api/internal/repository/mocks/mocks.go`)

## 编写单元测试步骤

### 1. 测试文件命名

```
api/internal/service/admin/user_test.go
```

### 2. 使用 Mock Repository

```go
package admin

import (
    "context"
    "testing"
    
    "gamelink/internal/model"
    "gamelink/internal/repository"
    "gamelink/internal/repository/mocks"
)

func TestGetUser(t *testing.T) {
    // Arrange: 创建 Mock Repository
    mockUserRepo := &mocks.MockUserRepository{
        GetFunc: func(ctx context.Context, id uint64) (*model.User, error) {
            if id == 1 {
                return &model.User{
                    ID:    1,
                    Name:  "Test User",
                    Email: "test@example.com",
                    Role:  model.RoleUser,
                }, nil
            }
            return nil, repository.ErrNotFound
        },
    }
    
    // 注入 Mock 依赖
    svc := &AdminService{
        users: mockUserRepo,
    }
    
    // Act: 调用被测试方法
    user, err := svc.GetUser(context.Background(), 1)
    
    // Assert: 验证结果
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if user.Name != "Test User" {
        t.Errorf("expected name 'Test User', got %s", user.Name)
    }
}
```

### 3. 测试驱动开发 (TDD)

**推荐流程：**

1. **Red**: 先写测试，测试失败
2. **Green**: 实现最简单的代码让测试通过
3. **Refactor**: 重构优化代码

## 测试模板

### 基本 CRUD 测试

```go
func TestCreateUser(t *testing.T) {
    tests := []struct {
        name          string
        input         CreateUserInput
        mockSetup     func(*mocks.MockUserRepository)
        expectedError error
    }{
        {
            name: "successful creation",
            input: CreateUserInput{
                Name:     "New User",
                Email:    "new@example.com",
                Password: "ValidPass123",
                Role:     model.RoleUser,
                Status:   model.UserStatusActive,
            },
            mockSetup: func(m *mocks.MockUserRepository) {
                m.FindByEmailFunc = func(ctx context.Context, email string) (*model.User, error) {
                    return nil, repository.ErrNotFound // 无重复
                }
                m.CreateFunc = func(ctx context.Context, user *model.User) error {
                    user.ID = 123 // 模拟数据库自增
                    return nil
                }
            },
            expectedError: nil,
        },
        {
            name: "duplicate email",
            input: CreateUserInput{
                Name:  "User",
                Email: "existing@example.com",
            },
            mockSetup: func(m *mocks.MockUserRepository) {
                m.FindByEmailFunc = func(ctx context.Context, email string) (*model.User, error) {
                    return &model.User{ID: 1}, nil // 模拟重复
                }
            },
            expectedError: apierr.ErrDuplicateEmail,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := &mocks.MockUserRepository{}
            tt.mockSetup(mockRepo)
            
            svc := &AdminService{users: mockRepo}
            _, err := svc.CreateUser(context.Background(), tt.input)
            
            // 验证错误类型
            if tt.expectedError != nil && !apierr.IsType(err, tt.expectedError) {
                t.Errorf("expected error %v, got %v", tt.expectedError, err)
            }
        })
    }
}
```

### 列表/分页测试

```go
func TestListUsersWithOptions(t *testing.T) {
    mockUsers := []model.User{
        {ID: 1, Name: "Alice", Role: model.RoleUser},
        {ID: 2, Name: "Bob", Role: model.RolePlayer},
    }
    
    mockRepo := &mocks.MockUserRepository{
        ListWithFiltersFunc: func(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
            return mockUsers, int64(len(mockUsers)), nil
        },
    }
    
    svc := &AdminService{users: mockRepo}
    users, pagination, err := svc.ListUsersWithOptions(context.Background(), UserListOptions{
        Page:     1,
        PageSize: 10,
    })
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(users) != 2 {
        t.Errorf("expected 2 users, got %d", len(users))
    }
    if pagination.Total != 2 {
        t.Errorf("expected total 2, got %d", pagination.Total)
    }
}
```

## 运行测试

```bash
# 运行单个测试
go test -v ./internal/service/admin -run TestGetUser

# 运行所有测试
go test ./...

# 查看覆盖率
go test -cover ./internal/service/admin

# 生成覆盖率报告
go test -coverprofile=coverage.out ./internal/service/admin
go tool cover -html=coverage.out
```

## 最佳实践

### ✅ DO

1. **使用 Table-Driven Tests**：一次测试多个场景
2. **Mock 外部依赖**：数据库、API、时间等
3. **测试边界条件**：空值、极大值、错误情况
4. **清晰的命名**：`TestMethodName_Scenario_ExpectedBehavior`
5. **独立测试**：每个测试互不依赖

### ❌ DON'T

1. **不要依赖真实数据库**：使用 Mock
2. **不要测试第三方库**：只测试自己的代码
3. **不要忽略错误**：明确断言错误类型
4. **不要过度 Mock**：只 Mock 必要的依赖
5. **不要写复杂测试**：一个测试只验证一个行为

## Mock Repository 现状

项目已有完整的 Mock 实现：

- `MockUserRepository`
- `MockPlayerRepository`
- `MockOrderRepository`
- `MockPaymentRepository`
- `MockGameRepository`
- ...等 50+ Mock

**位置**: `api/internal/repository/mocks/mocks.go`

## 当前测试覆盖率

```bash
# 查看当前覆盖率
go test -cover ./internal/service/admin
```

**目标**: 从当前提升至 **70%+** 覆盖率

## 示例：完整测试套件

参考文件：`api/internal/service/admin/user_test.go`

该文件包含：
- ✅ 获取用户测试 (`TestGetUser`)
- ✅ 创建用户测试 (`TestCreateUser`)
- ✅ 列表查询测试 (`TestListUsersWithOptions`)
- ✅ 边界条件测试（未找到、重复等）

## 持续集成

测试应在 CI/CD 流程中自动运行：

```yaml
# .github/workflows/ci.yml
- name: Run tests
  run: go test -v -race -coverprofile=coverage.out ./...

- name: Check coverage
  run: |
    coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    echo "Coverage: $coverage"
    # 可设置最低覆盖率要求
```

## 下一步

1. **补充现有测试**：为核心 Service 方法添加测试
2. **提升覆盖率**：重点覆盖复杂业务逻辑
3. **集成测试**：必要时添加集成测试（使用测试数据库）
4. **性能测试**：使用 `testing.B` 进行基准测试

## 资源

- [Go Testing 官方文档](https://pkg.go.dev/testing)
- [Table-Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Testify 断言库](https://github.com/stretchr/testify)

---

**维护者**: Backend Team  
**最后更新**: 2026-02-10  
**版本**: 1.0
