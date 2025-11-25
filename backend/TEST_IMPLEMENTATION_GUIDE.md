# 测试实施指南 - 达到100%覆盖率

## 📋 快速开始

### 已完成的包
✅ **internal/apierr** - **100%覆盖率** ⭐示例参考

### 工具和资源

#### 1. 测试骨架生成工具
```bash
# 为指定包生成测试骨架
go run scripts/generate_test_skeleton.go internal/cache
go run scripts/generate_test_skeleton.go internal/config
go run scripts/generate_test_skeleton.go internal/db
```

#### 2. 查看当前覆盖率
```bash
# 测试单个包
cd internal/<package_name>
go test -cover -v

# 生成详细覆盖率报告
go test -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

#### 3. 测试所有包
```bash
# 在项目根目录
go test ./... -cover

# 生成总体覆盖率报告
go test ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -func=coverage.out | grep total
```

## 🎯 测试策略 - 按优先级执行

### 第一阶段: 基础设施层 (本周完成)

#### 1. internal/cache (高优先级)
**文件**: cache.go, lock.go, memory.go, redis.go

**测试要点**:
- ✅ 内存缓存的CRUD操作
- ✅ TTL过期机制
- ✅ 并发安全性
- ✅ Redis连接和操作
- ✅ 分布式锁的获取/释放
- ✅ 锁重试机制

**示例测试结构**:
```go
func Test_MemoryCache_SetAndGet(t *testing.T) {
	cache := NewMemory()
	ctx := context.Background()

	// 测试设置和获取
	err := cache.Set(ctx, "key1", "value1", 5*time.Second)
	assert.NoError(t, err)

	val, ok, err := cache.Get(ctx, "key1")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "value1", val)
}

func Test_MemoryCache_TTLExpiry(t *testing.T) {
	cache := NewMemory()
	ctx := context.Background()

	// 设置短TTL
	cache.Set(ctx, "key1", "value1", 100*time.Millisecond)

	// 等待过期
	time.Sleep(200*time.Millisecond)

	// 验证已过期
	_, ok, _ := cache.Get(ctx, "key1")
	assert.False(t, ok)
}
```

**预计工作量**: 2-3小时

#### 2. internal/config (高优先级)
**文件**: database.go, defaults.go, env.go, validate.go

**测试要点**:
- ✅ 配置加载和解析
- ✅ 环境变量读取
- ✅ 默认值设置
- ✅ 配置验证
- ✅ 错误处理

**预计工作量**: 3-4小时

#### 3. internal/db (高优先级)
**文件**: db.go, migrate.go, postgres.go, sqlite.go, seed.go

**测试要点**:
- ✅ 数据库连接
- ✅ 迁移执行
- ✅ SQLite/PostgreSQL特定功能
- ✅ 种子数据生成
- ✅ 错误处理

**预计工作量**: 4-5小时

### 第二阶段: 数据访问层 (第2周)

#### 1. internal/repository/* (高优先级)
**当前状态**: 部分有测试

**需要补充的子包**:
- chat/
- commission/
- common/
- dispute/
- feed/
- notification/
- operation_log/
- payment/
- player/
- player_tag/
- ranking/
- review/
- reviewreply/
- role/
- stats/
- withdraw/

**测试模式** (使用testify/mock):
```go
func Test_UserRepository_Create(t *testing.T) {
	// 使用内存数据库或mock
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := NewUserRepository(db)

	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
	}

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}
```

**预计工作量**: 每个子包1-2小时 × 15 = 15-30小时

### 第三阶段: 业务逻辑层 (第3-4周)

#### 1. internal/service/* (高优先级)
**包含**: admin, assignment, auth, chat, commission, earnings等

**测试要点**:
- ✅ 业务逻辑正确性
- ✅ 事务处理
- ✅ 错误处理
- ✅ 权限检查
- ✅ Mock外部依赖

**测试模式**:
```go
func Test_OrderService_CreateOrder(t *testing.T) {
	// Mock仓储层
	mockRepo := new(MockOrderRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewOrderService(mockRepo, mockUserRepo)

	// 设置mock期望
	mockUserRepo.On("GetByID", mock.Anything, uint(1)).
		Return(&model.User{ID: 1}, nil)
	mockRepo.On("Create", mock.Anything, mock.Anything).
		Return(nil)

	// 测试创建订单
	order, err := service.CreateOrder(ctx, &CreateOrderRequest{
		UserID: 1,
		GameID: 1,
	})

	assert.NoError(t, err)
	assert.NotNil(t, order)
	mockRepo.AssertExpectations(t)
}
```

**预计工作量**: 每个服务2-4小时 × 20 = 40-80小时

### 第四阶段: HTTP处理层 (第5-6周)

#### 1. internal/handler/* (高优先级)
**包含**: admin/, player/, user/, notification/

**测试要点**:
- ✅ 请求解析
- ✅ 响应格式
- ✅ HTTP状态码
- ✅ 错误处理
- ✅ 认证授权

**测试模式** (使用httptest):
```go
func Test_UserHandler_Register(t *testing.T) {
	// 设置测试环境
	router := gin.New()
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	router.POST("/register", handler.Register)

	// 准备请求
	reqBody := `{"username":"test","email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "success", response["status"])
}
```

**预计工作量**: 每个handler文件1-2小时 × 40 = 40-80小时

#### 2. internal/handler/middleware/* (中优先级)
**文件**: cors.go, csrf.go, jwt.go, permission.go, rate_limit.go等

**测试要点**:
- ✅ 中间件逻辑
- ✅ 请求链处理
- ✅ 错误处理

**预计工作量**: 每个中间件0.5-1小时 × 12 = 6-12小时

### 第五阶段: 支撑系统 (第7周)

#### 1. internal/metrics (中优先级)
#### 2. internal/logging (中优先级)
#### 3. internal/router (中优先级)
#### 4. internal/scheduler (低优先级)
#### 5. internal/lifecycle (低优先级)
#### 6. internal/container (低优先级)

**预计工作量**: 总计15-20小时

## 🛠️ 测试编写最佳实践

### 1. 测试命名规范
```go
// 格式: Test_<FunctionName>_<Scenario>
Test_CreateUser_Success
Test_CreateUser_DuplicateEmail
Test_CreateUser_InvalidInput
```

### 2. AAA模式 (Arrange-Act-Assert)
```go
func Test_Example(t *testing.T) {
	// Arrange - 准备测试数据
	user := &User{Name: "test"}

	// Act - 执行被测试的代码
	result, err := CreateUser(user)

	// Assert - 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
```

### 3. 表驱动测试
```go
func Test_ValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "test@example.com", false},
		{"invalid email", "invalid", true},
		{"empty email", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
```

### 4. Mock外部依赖
```go
// 使用 testify/mock
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetByID(ctx context.Context, id uint) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}
```

### 5. 测试清理
```go
func Test_Example(t *testing.T) {
	// 设置
	db := setupTestDB(t)
	defer cleanupTestDB(t, db) // 确保清理

	// 测试逻辑...
}
```

## 📊 覆盖率目标

| 层级 | 目标覆盖率 | 优先级 |
|------|-----------|--------|
| 基础设施层 (apierr, config, cache, db) | 95%+ | 🔴 高 |
| 数据访问层 (repository) | 85%+ | 🔴 高 |
| 业务逻辑层 (service) | 90%+ | 🔴 高 |
| HTTP处理层 (handler) | 80%+ | 🟡 中 |
| 中间件层 (middleware) | 75%+ | 🟡 中 |
| 支撑系统 | 70%+ | 🟢 低 |
| **总体目标** | **100%** | ⭐ |

## 🚀 执行计划

### Week 1: 基础设施层
- [x] internal/apierr (已完成 100%)
- [ ] internal/cache
- [ ] internal/config
- [ ] internal/db

### Week 2-3: 数据访问层
- [ ] 完成所有 internal/repository 子包测试
- [ ] 补充现有测试的覆盖率

### Week 4-5: 业务逻辑层
- [ ] 完成所有 internal/service 子包测试

### Week 6-7: HTTP和支撑层
- [ ] internal/handler 所有处理器
- [ ] internal/handler/middleware 所有中间件
- [ ] 支撑系统 (metrics, logging, router等)

### Week 8: 收尾和验证
- [ ] 检查并补充遗漏的测试
- [ ] 验证总体覆盖率达到100%
- [ ] 生成最终覆盖率报告
- [ ] 代码审查和质量检查

## 📝 每日工作流程

1. **选择包**: 按优先级选择要测试的包
2. **生成骨架**: 使用测试生成工具创建基础测试文件
3. **实现测试**: 根据业务逻辑编写具体测试
4. **运行测试**: `go test -cover`
5. **检查覆盖率**: 确保达到目标覆盖率
6. **提交代码**: 使用规范的commit message
7. **更新进度**: 在TEST_COVERAGE_ANALYSIS.md中更新状态

## 🎯 质量检查清单

每个包完成测试后,确保:
- [ ] 所有公共函数都有测试
- [ ] 测试覆盖率达到目标
- [ ] 包含边界条件测试
- [ ] 包含错误处理测试
- [ ] 使用AAA模式组织测试
- [ ] 测试命名清晰规范
- [ ] Mock了外部依赖
- [ ] 所有测试都能通过
- [ ] 没有t.Skip()的未完成测试
- [ ] 代码有适当的注释

## 💡 常见问题解决

### 1. 数据库测试
```go
// 使用SQLite内存数据库
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// 运行迁移
	err = db.AutoMigrate(&User{}, &Order{})
	require.NoError(t, err)

	return db
}
```

### 2. 时间相关测试
```go
// 使用时间冻结
func Test_TimeDependent(t *testing.T) {
	now := time.Now()
	mockTime := func() time.Time { return now }

	// 使用mockTime代替time.Now()
}
```

### 3. 并发测试
```go
func Test_Concurrent(t *testing.T) {
	var wg sync.WaitGroup
	cache := NewMemory()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", n)
			cache.Set(context.Background(), key, "value", time.Minute)
		}(i)
	}

	wg.Wait()
	// 验证结果
}
```

## 📚 参考资源

- [internal/apierr/errors_test.go](internal/apierr/errors_test.go) - 100%覆盖率示例
- [TEST_COVERAGE_ANALYSIS.md](TEST_COVERAGE_ANALYSIS.md) - 详细分析报告
- [Go Testing文档](https://golang.org/pkg/testing/)
- [Testify文档](https://github.com/stretchr/testify)

## 🎉 成功标准

当所有以下条件满足时,项目测试达标:
- ✅ 总体覆盖率 >= 100%
- ✅ 核心模块覆盖率 >= 90%
- ✅ 所有测试都能通过
- ✅ 没有被跳过的测试
- ✅ 通过代码审查
- ✅ CI/CD集成完成
