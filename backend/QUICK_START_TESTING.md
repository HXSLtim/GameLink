# 🚀 测试快速开始指南

> 5分钟开始为GameLink编写测试,达到100%覆盖率

## ⚡ 快速命令

### 查看当前进度
```bash
# 显示测试覆盖率进度
./scripts/batch_test.sh progress

# 或手动查看
go test ./... -cover
```

### 为新包生成测试
```bash
# 方法1: 使用批量脚本
./scripts/batch_test.sh generate internal/cache

# 方法2: 直接使用生成工具
go run scripts/generate_test_skeleton.go internal/cache
```

### 运行测试
```bash
# 测试单个包
./scripts/batch_test.sh test internal/cache

# 测试所有包
./scripts/batch_test.sh all

# 生成HTML报告
./scripts/batch_test.sh report internal/cache
```

## 📝 第一个测试示例

### 步骤1: 选择要测试的文件
例如: `internal/cache/memory.go`

### 步骤2: 生成测试骨架
```bash
go run scripts/generate_test_skeleton.go internal/cache
```

这会创建 `internal/cache/memory_test.go`

### 步骤3: 实现测试
打开生成的测试文件,实现具体测试:

```go
package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/**
 * Test_MemoryCache_SetAndGet 测试内存缓存的设置和获取
 */
func Test_MemoryCache_SetAndGet(t *testing.T) {
	// Arrange - 准备
	cache := NewMemory()
	ctx := context.Background()
	key := "test_key"
	value := "test_value"

	// Act - 执行
	err := cache.Set(ctx, key, value, time.Minute)
	assert.NoError(t, err)

	result, ok, err := cache.Get(ctx, key)

	// Assert - 验证
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, value, result)
}

/**
 * Test_MemoryCache_TTL 测试TTL过期
 */
func Test_MemoryCache_TTL(t *testing.T) {
	// Arrange
	cache := NewMemory()
	ctx := context.Background()

	// Act
	cache.Set(ctx, "key", "value", 100*time.Millisecond)
	time.Sleep(200 * time.Millisecond)

	// Assert
	_, ok, _ := cache.Get(ctx, "key")
	assert.False(t, ok, "键应该已过期")
}
```

### 步骤4: 运行测试
```bash
cd internal/cache
go test -v -cover
```

### 步骤5: 检查覆盖率
```bash
go test -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## 🎯 今天就开始!

### 立即完成基础设施层 (2-3天)

```bash
# Day 1: Cache包
go run scripts/generate_test_skeleton.go internal/cache
# 实现测试... (参考 internal/apierr/errors_test.go)
./scripts/batch_test.sh test internal/cache

# Day 2: Config包
go run scripts/generate_test_skeleton.go internal/config
# 实现测试...
./scripts/batch_test.sh test internal/config

# Day 3: DB包
go run scripts/generate_test_skeleton.go internal/db
# 实现测试...
./scripts/batch_test.sh test internal/db
```

### 每天的工作流程 (30分钟-2小时)

1. **早上**: 选择1-2个包
2. **生成**: 使用工具生成测试骨架 (5分钟)
3. **实现**: 编写测试逻辑 (30-60分钟)
4. **验证**: 运行测试并检查覆盖率 (5分钟)
5. **提交**: 提交代码 (5分钟)

## 📚 参考资料

### 必读文档
1. **TEST_IMPLEMENTATION_GUIDE.md** - 完整实施指南 ⭐
2. **TEST_COVERAGE_ANALYSIS.md** - 项目分析
3. **internal/apierr/errors_test.go** - 100%覆盖率示例

### 测试模式速查

#### 基本测试结构 (AAA模式)
```go
func Test_FunctionName(t *testing.T) {
	// Arrange - 准备数据
	input := "test"

	// Act - 执行函数
	result, err := FunctionName(input)

	// Assert - 验证结果
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
```

#### 表驱动测试
```go
func Test_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "test@example.com", false},
		{"invalid", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
```

#### Mock外部依赖
```go
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Get(id int) (*User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func Test_WithMock(t *testing.T) {
	mockRepo := new(MockRepository)
	mockRepo.On("Get", 1).Return(&User{ID: 1}, nil)

	service := NewService(mockRepo)
	user, err := service.GetUser(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	mockRepo.AssertExpectations(t)
}
```

## 💡 实用技巧

### 快速达到高覆盖率
1. ✅ 先测试主要流程(happy path)
2. ✅ 再测试错误处理
3. ✅ 最后测试边界条件

### 常见陷阱避免
1. ❌ 不要测试私有函数的实现细节
2. ❌ 不要依赖测试执行顺序
3. ❌ 不要在测试中使用Sleep (除非必要)
4. ❌ 不要共享测试状态

### 提高效率
1. ✅ 使用测试生成工具
2. ✅ 复制粘贴相似的测试模式
3. ✅ 使用表驱动测试减少重复
4. ✅ 并行开发多个包的测试

## 🎉 完成一个包的检查清单

完成每个包后,确认:
- [ ] 所有导出函数都有测试
- [ ] 测试覆盖率 >= 目标值
- [ ] 所有测试都能通过
- [ ] 没有被skip的测试
- [ ] 包含边界条件测试
- [ ] 包含错误处理测试
- [ ] 代码有清晰的注释
- [ ] 提交代码到版本控制

## 🔗 快速链接

- [完整实施指南](TEST_IMPLEMENTATION_GUIDE.md)
- [项目分析报告](TEST_COVERAGE_ANALYSIS.md)
- [工作总结](TEST_COMPLETION_SUMMARY.md)
- [示例代码](internal/apierr/errors_test.go)

## ❓ 需要帮助?

1. 查看示例: `internal/apierr/errors_test.go`
2. 阅读指南: `TEST_IMPLEMENTATION_GUIDE.md`
3. 使用工具: `scripts/batch_test.sh`

---

**现在就开始编写测试,让代码更可靠! 🚀**

```bash
# 第一步
./scripts/batch_test.sh progress

# 第二步
./scripts/batch_test.sh generate internal/cache

# 第三步 - 开始编写测试!
```
