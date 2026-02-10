# 单元测试实战示例

## 问题现状

当前 `api/internal/service/admin/adminService_test.go` 存在编译错误，因为旧的 Mock 代码未及时更新。

## 解决方案

### 方案 1: 重新生成 Mock (推荐)

使用 `mockgen` 工具自动生成最新的 Mock：

```bash
# 安装 mockgen
go install github.com/golang/mock/mockgen@latest

# 为 OrderRepository 生成 Mock
mockgen -source=internal/repository/interfaces/order.go \
        -destination=internal/repository/mocks/order_mock.go \
        -package=mocks

# 为其他 Repository 生成 Mock
mockgen -source=internal/repository/interfaces.go \
        -destination=internal/repository/mocks/repo_mocks.go \
        -package=mocks
```

### 方案 2: 手写简单 Mock (快速验证)

对于简单场景，可以手写 Mock 而不依赖 gomock：

#### 示例：测试 User 创建逻辑

**文件：`internal/service/admin/simple_user_test.go`**

```go
package admin

import (
    "context"
    "testing"
    
    "gamelink/internal/model"
    "gamelink/internal/repository"
)

// 简单的手写 Mock，不依赖 gomock
type simpleUserRepoMock struct {
    users map[uint64]*model.User
}

func (m *simpleUserRepoMock) Get(ctx context.Context, id uint64) (*model.User, error) {
    if user, ok := m.users[id]; ok {
        return user, nil
    }
    return nil, repository.ErrNotFound
}

func (m *simpleUserRepoMock) Create(ctx context.Context, user *model.User) error {
    user.ID = uint64(len(m.users) + 1) // 模拟自增 ID
    m.users[user.ID] = user
    return nil
}

func (m *simpleUserRepoMock) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    for _, u := range m.users {
        if u.Email == email {
            return u, nil
        }
    }
    return nil, repository.ErrNotFound
}

// 其他方法返回默认值或 nil
func (m *simpleUserRepoMock) List(ctx context.Context) ([]model.User, error) {
    return nil, nil
}
func (m *simpleUserRepoMock) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
    return nil, 0, nil
}
func (m *simpleUserRepoMock) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
    return nil, 0, nil
}
func (m *simpleUserRepoMock) Count(ctx context.Context, opts repository.UserListOptions) (int, error) {
    return 0, nil
}
func (m *simpleUserRepoMock) GetByIDs(ctx context.Context, ids []uint64) ([]model.User, error) {
    return nil, nil
}
func (m *simpleUserRepoMock) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
    return nil, repository.ErrNotFound
}
func (m *simpleUserRepoMock) GetByEmail(ctx context.Context, email string) (*model.User, error) {
    return m.FindByEmail(ctx, email)
}
func (m *simpleUserRepoMock) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
    return nil, repository.ErrNotFound
}
func (m *simpleUserRepoMock) Update(ctx context.Context, user *model.User) error {
    return nil
}
func (m *simpleUserRepoMock) UpdatePassword(ctx context.Context, userID uint64, newPassword string) error {
    return nil
}
func (m *simpleUserRepoMock) Delete(ctx context.Context, id uint64) error {
    return nil
}
func (m *simpleUserRepoMock) GetByWeChatOpenID(ctx context.Context, openID string) (*model.User, error) {
    return nil, repository.ErrNotFound
}
func (m *simpleUserRepoMock) GetByWeChatUnionID(ctx context.Context, unionID string) (*model.User, error) {
    return nil, repository.ErrNotFound
}

// 实际测试
func TestSimpleGetUser(t *testing.T) {
    // Arrange: 准备数据
    mockRepo := &simpleUserRepoMock{
        users: map[uint64]*model.User{
            1: {
                ID:    1,
                Name:  "Test User",
                Email: "test@example.com",
                Role:  model.RoleUser,
            },
        },
    }
    
    svc := &AdminService{
        users: mockRepo,
    }
    
    // Act: 执行测试
    user, err := svc.GetUser(context.Background(), 1)
    
    // Assert: 验证结果
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if user.Name != "Test User" {
        t.Errorf("expected name 'Test User', got %s", user.Name)
    }
    
    // 测试不存在的用户
    _, err = svc.GetUser(context.Background(), 999)
    if err != repository.ErrNotFound {
        t.Errorf("expected ErrNotFound, got %v", err)
    }
}
```

### 运行示例测试

```bash
cd api

# 运行单个测试
go test -v ./internal/service/admin -run TestSimpleGetUser

# 查看覆盖率
go test -cover ./internal/service/admin -run TestSimpleGetUser
```

## 推荐工作流

### 短期（当前）

1. **跳过旧测试**：暂时忽略 `adminService_test.go` 的编译错误
2. **创建新测试**：使用简单 Mock 为新功能编写测试
3. **逐步积累**：慢慢建立测试覆盖率

```bash
# 只运行新测试
go test -v ./internal/service/admin -run TestSimple
```

### 中期（1-2周）

1. **修复 Mock 生成**：使用 mockgen 重新生成所有 Mock
2. **修复旧测试**：逐个修复 `adminService_test.go` 中的测试
3. **提升覆盖率**：补充缺失的测试用例

### 长期（1个月+）

1. **引入测试框架**：考虑使用 [testify](https://github.com/stretchr/testify) 简化断言
2. **集成测试**：添加使用真实数据库的集成测试
3. **CI 集成**：在 CI/CD 中强制测试通过

## 测试最佳实践

### ✅ 好的做法

```go
// 1. 清晰的测试命名
func TestGetUser_WhenUserExists_ReturnsUser(t *testing.T) { }

// 2. Arrange-Act-Assert 模式
func TestCreateUser(t *testing.T) {
    // Arrange
    mockRepo := setupMockRepo()
    svc := NewService(mockRepo)
    
    // Act
    result, err := svc.CreateUser(ctx, input)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}

// 3. Table-Driven Tests
func TestValidateEmail(t *testing.T) {
    tests := []struct{
        name  string
        email string
        valid bool
    }{
        {"valid email", "user@example.com", true},
        {"invalid format", "not-an-email", false},
        {"empty", "", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := validateEmail(tt.email)
            if result != tt.valid {
                t.Errorf("expected %v, got %v", tt.valid, result)
            }
        })
    }
}
```

### ❌ 避免的做法

```go
// 1. 不要测试太多东西
func TestEverything(t *testing.T) {
    // 测试创建、更新、删除、查询... ❌
}

// 2. 不要依赖顺序
func TestA(t *testing.T) { /* 创建 */ }
func TestB(t *testing.T) { /* 依赖 TestA 的数据 ❌ */ }

// 3. 不要忽略错误
result, _ := svc.DoSomething() // ❌ 应该检查错误
```

## 总结

当前最务实的做法：

1. **接受现状**：旧测试有编译错误，暂时不修
2. **专注新测试**：为新功能写简单的手写 Mock 测试
3. **逐步改进**：等有时间再系统性修复 Mock 生成问题

**优先级**：功能开发 > 新功能测试 > 修复旧测试

---

**参考资源**

- [Go Testing 官方文档](https://pkg.go.dev/testing)
- [gomock GitHub](https://github.com/golang/mock)
- [testify GitHub](https://github.com/stretchr/testify)
- [Table Driven Tests 模式](https://github.com/golang/go/wiki/TableDrivenTests)

**维护者**: Backend Team  
**最后更新**: 2026-02-10  
**版本**: 1.0
