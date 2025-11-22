# 📅 GameLink 100%测试覆盖率 - 周实施计划

**计划周期**: 2025-11-23 至 2026-02-15 (12周)
**版本**: 1.0
**状态**: 执行中

---

## 📊 周计划概览

| 周次 | 日期范围 | 核心任务 | 预期覆盖率 | 负责人 |
|-----|---------|---------|-----------|--------|
| **Week 1** | 11/23-11/29 | 基础设施搭建 | 0% → 5% | 测试架构师 |
| **Week 2** | 11/30-12/06 | Auth模块测试 | 5% → 15% | 工程师A |
| **Week 3** | 12/07-12/13 | Payment模块测试 | 15% → 25% | 工程师B |
| **Week 4** | 12/14-12/20 | Order模块测试 | 25% → 35% | 工程师A |
| **Week 5** | 12/21-12/27 | Commission+Earnings | 35% → 45% | 工程师B |
| **Week 6** | 12/28-01/03 | Player+User模块 | 45% → 55% | 工程师A |
| **Week 7** | 01/04-01/10 | Permission+Role | 55% → 65% | 工程师B |
| **Week 8** | 01/11-01/17 | 其他模块 | 65% → 70% | 工程师A |
| **Week 9** | 01/18-01/24 | Repository层-Part1 | 70% → 80% | 两人协作 |
| **Week 10** | 01/25-01/31 | Repository层-Part2 | 80% → 85% | 两人协作 |
| **Week 11** | 02/01-02/07 | Handler层测试 | 85% → 95% | 两人协作 |
| **Week 12** | 02/08-02/15 | 边缘覆盖+优化 | 95% → 100% | 测试架构师 |

---

## 🔍 Week 1: 基础设施搭建 (11/23-11/29)

### 每日任务分解

#### Day 1 (周一): 测试框架配置
**负责人**: 测试架构师
**时间**: 8小时

**任务清单**:
- [ ] 安装测试依赖包 (2小时)
  ```bash
  go get github.com/stretchr/testify@v1.8.4
  go get github.com/golang/mock/mockgen@v1.6.0
  go get github.com/h2non/gock@v1.2.0
  go get github.com/bxcodec/faker/v3@v3.8.1
  go get github.com/DATA-DOG/go-sqlmock@v1.5.0
  ```

- [ ] 更新go.mod (1小时)
- [ ] 创建测试工具包目录 (1小时)
  ```
  internal/testutil/
  ├── db.go              # 数据库测试工具
  ├── mock.go            # Mock生成工具
  ├── assert.go          # 自定义断言
  ├── fixture.go         # 测试数据工厂
  └── coverage.go        # 覆盖率工具
  ```

- [ ] 编写测试工具代码 (4小时)

**交付物**:
- ✅ 更新的go.mod和go.sum
- ✅ internal/testutil/ 工具包

---

#### Day 2 (周二): 测试工具实现 (Part 1)
**负责人**: 测试架构师
**时间**: 8小时

**任务清单**:
- [ ] 实现DB测试工具 (4小时)
  ```go
  // internal/testutil/db.go
  package testutil
  
  import (
      "gorm.io/driver/sqlite"
      "gorm.io/gorm"
  )
  
  // NewMemoryDB 创建内存数据库用于测试
  func NewMemoryDB(t *testing.T) *gorm.DB {
      db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
      if err != nil {
          t.Fatalf("failed to create memory db: %v", err)
      }
      return db
  }
  
  // MigrateTables 自动迁移表结构
  func MigrateTables(t *testing.T, db *gorm.DB, models ...interface{}) {
      if err := db.AutoMigrate(models...); err != nil {
          t.Fatalf("failed to migrate tables: %v", err)
      }
  }
  ```

- [ ] 实现Mock工具 (2小时)
- [ ] 实现Fixture工厂 (2小时)

**交付物**:
- ✅ 完整的DB测试工具
- ✅ Fixture数据工厂

---

#### Day 3 (周三): 测试工具实现 (Part 2)
**负责人**: 测试架构师
**时间**: 8小时

**任务清单**:
- [ ] 实现自定义断言 (3小时)
  ```go
  // internal/testutil/assert.go
  package testutil
  
  import (
      "github.com/stretchr/testify/assert"
      "gamelink/internal/model"
  )
  
  // AssertUserEqual 断言用户相等
  func AssertUserEqual(t *testing.T, expected, actual *model.User) {
      assert.Equal(t, expected.ID, actual.ID)
      assert.Equal(t, expected.Email, actual.Email)
      assert.Equal(t, expected.Name, actual.Name)
      assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, 0)
  }
  ```

- [ ] 实现覆盖率工具 (3小时)
- [ ] 编写测试工具文档 (2小时)

**交付物**:
- ✅ 自定义断言库
- ✅ 覆盖率监控工具
- ✅ 测试工具使用文档

---

#### Day 4 (周四): CI/CD集成
**负责人**: 测试架构师
**时间**: 8小时

**任务清单**:
- [ ] 创建GitHub Actions工作流 (4小时)
  ```yaml
  # .github/workflows/test.yml
  name: Test and Coverage
  
  on:
    push:
      branches: [ main, develop ]
    pull_request:
      branches: [ main, develop ]
  
  jobs:
    test:
      runs-on: ubuntu-latest
      
      services:
        postgres:
          image: postgres:14
          env:
            POSTGRES_PASSWORD: postgres
          options: >-
            --health-cmd pg_isready
            --health-interval 10s
            --health-timeout 5s
            --health-retries 5
          ports:
            - 5432:5432
      
      steps:
        - uses: actions/checkout@v3
        
        - name: Set up Go
          uses: actions/setup-go@v4
          with:
            go-version: '1.21'
        
        - name: Run tests
          run: |
            go test ./... -coverprofile=coverage.out -covermode=atomic
            go tool cover -func=coverage.out
        
        - name: Upload coverage
          uses: codecov/codecov-action@v3
          with:
            file: ./coverage.out
            fail_ci_if_error: true
            threshold: 80%
  ```

- [ ] 配置Makefile (2小时)
  ```makefile
  .PHONY: test
  test:
  	go test ./... -coverprofile=coverage.out -covermode=atomic
  
  .PHONY: test-coverage
  test-coverage: test
  	go tool cover -html=coverage.out -o coverage.html
  
  .PHONY: test-watch
  test-watch:
  	gow test ./... -cover
  ```

- [ ] 配置Codecov (2小时)
  - 注册Codecov账号
  - 配置GitHub集成
  - 设置覆盖率门槛

**交付物**:
- ✅ GitHub Actions工作流
- ✅ Makefile测试命令
- ✅ Codecov集成

---

#### Day 5 (周五): Mock生成与文档
**负责人**: 测试架构师
**时间**: 8小时

**任务清单**:
- [ ] 生成Repository接口Mock (4小时)
  ```bash
  # 安装mockery
  go install github.com/vektra/mockery/v2@latest
  
  # 生成所有Repository Mock
  cd internal/repository
  mockery --all --output=mocks --outpkg=mocks
  ```

- [ ] 生成Service接口Mock (2小时)
- [ ] 编写测试指南文档 (2小时)
  ```markdown
  # 测试指南
  
  ## 如何编写测试
  
  1. 使用testify suite
  2. 每个测试必须包含Arrange/Act/Assert
  3. Mock必须精确匹配
  4. 覆盖率必须100%
  
  ## 如何运行测试
  
  ```bash
  # 运行所有测试
  make test
  
  # 查看覆盖率报告
  make test-coverage
  
  # 持续测试
  make test-watch
  ```
  
  ## 测试规范
  
  - 测试文件必须以`_test.go`结尾
  - 测试函数必须以`Test`开头
  - 使用表格驱动测试覆盖多个场景
  ```

**交付物**:
- ✅ 生成的Mock文件
- ✅ 测试指南文档
- ✅ 团队培训材料

---

## 🔍 Week 2: Auth模块测试 (11/30-12/06)

### 本周目标
- 覆盖所有Auth模块代码
- 达到100%覆盖率
- 编写50+测试用例

### 模块分解

#### Day 1-2: JWT核心测试
**负责人**: 工程师A
**测试文件**: `internal/auth/jwt_test.go`

**测试用例清单**:
- [ ] TestGenerateToken_Success
- [ ] TestGenerateToken_EmptyClaims
- [ ] TestValidateToken_Success
- [ ] TestValidateToken_Expired
- [ ] TestValidateToken_InvalidSignature
- [ ] TestValidateToken_Malformed
- [ ] TestRefreshToken_Success
- [ ] TestRefreshToken_Expired
- [ ] TestTokenConcurrency_Safe

**预期覆盖率**: 100%

#### Day 3-4: 密码安全测试
**负责人**: 工程师A
**测试文件**: `internal/auth/password_test.go`

**测试用例清单**:
- [ ] TestHashPassword_Success
- [ ] TestHashPassword_VerifySuccess
- [ ] TestHashPassword_VerifyFail
- [ ] TestHashPassword_Cost
- [ ] TestVerifyPassword_Correct
- [ ] TestVerifyPassword_Incorrect
- [ ] TestPasswordConcurrency_Safe

**预期覆盖率**: 100%

#### Day 5: 权限检查测试
**负责人**: 工程师A
**测试文件**: `internal/auth/permission_test.go`

**测试用例清单**:
- [ ] TestCheckPermission_Allowed
- [ ] TestCheckPermission_Denied
- [ ] TestCheckPermission_InvalidToken
- [ ] TestCheckPermission_EmptyToken
- [ ] TestCheckPermission_Expired

**预期覆盖率**: 100%

### 交付标准
- [ ] 所有测试通过
- [ ] 覆盖率100%
- [ ] Mock使用正确
- [ ] Code Review通过

---

## 💳 Week 3: Payment模块测试 (12/07-12/13)

### 本周目标
- 覆盖Payment Service所有代码
- 达到100%覆盖率
- 编写40+测试用例

### 模块分解

#### Day 1-2: 创建支付测试
**负责人**: 工程师B
**测试文件**: `internal/service/payment/create_test.go`

**测试场景**:
- [ ] 创建支付成功
- [ ] 订单不存在
- [ ] 订单已支付
- [ ] 订单已取消
- [ ] 金额不匹配
- [ ] 用户不匹配

#### Day 3: 支付查询测试
**负责人**: 工程师B
**测试文件**: `internal/service/payment/query_test.go`

**测试场景**:
- [ ] 查询支付成功
- [ ] 支付不存在
- [ ] 无权限查询

#### Day 4: 支付取消测试
**负责人**: 工程师B
**测试文件**: `internal/service/payment/cancel_test.go`

**测试场景**:
- [ ] 取消支付成功
- [ ] 支付已处理
- [ ] 支付已取消
- [ ] 无权限取消

#### Day 5: 回调处理测试
**负责人**: 工程师B
**测试文件**: `internal/service/payment/callback_test.go`

**测试场景**:
- [ ] 回调成功
- [ ] 回调签名错误
- [ ] 重复回调
- [ ] 回调超时

### 交付标准
- [ ] 所有测试通过
- [ ] 覆盖率100%
- [ ] 并发测试覆盖
- [ ] Code Review通过

---

## 📦 Week 4: Order模块测试 (12/14-12/20)

### 本周目标
- 覆盖Order Service所有代码
- 达到100%覆盖率
- 编写60+测试用例

### 模块分解

#### Day 1: 订单创建测试
**负责人**: 工程师A
**测试文件**: `internal/service/order/create_test.go`

**测试场景**:
- [ ] 创建订单成功
- [ ] 陪玩师不存在
- [ ] 游戏不存在
- [ ] 参数验证失败
- [ ] 并发创建订单

#### Day 2: 订单状态流转测试
**负责人**: 工程师A
**测试文件**: `internal/service/order/status_test.go`

**测试场景**:
- [ ] 所有合法状态流转 (pending→confirmed→in_progress→completed)
- [ ] 所有非法状态流转 (pending→completed)
- [ ] 状态流转通知

#### Day 3: 订单取消测试
**负责人**: 工程师A
**测试文件**: `internal/service/order/cancel_test.go`

**测试场景**:
- [ ] 取消订单成功
- [ ] 订单已支付取消
- [ ] 订单已完成取消
- [ ] 无权限取消

#### Day 4: 订单查询测试
**负责人**: 工程师A
**测试文件**: `internal/service/order/query_test.go`

**测试场景**:
- [ ] 查询用户订单
- [ ] 查询陪玩师订单
- [ ] 分页查询
- [ ] 条件查询

#### Day 5: 订单超时测试
**负责人**: 工程师A
**测试文件**: `internal/service/order/timeout_test.go`

**测试场景**:
- [ ] 订单自动取消
- [ ] 订单超时提醒
- [ ] 订单超时赔付

### 交付标准
- [ ] 所有测试通过
- [ ] 覆盖率100%
- [ ] 状态机测试完整
- [ ] Code Review通过

---

## 💰 Week 5: Commission + Earnings (12/21-12/27)

### 本周目标
- 覆盖Commission和Earnings模块
- 达到100%覆盖率
- 编写50+测试用例

### 模块分解

#### Day 1-2: Commission规则测试
**负责人**: 工程师B
**测试文件**: `internal/service/commission/rule_test.go`

**测试场景**:
- [ ] 创建抽成规则
- [ ] 更新抽成规则
- [ ] 删除抽成规则
- [ ] 规则计算验证

#### Day 3: Commission结算测试
**负责人**: 工程师B
**测试文件**: `internal/service/commission/settlement_test.go`

**测试场景**:
- [ ] 月度结算
- [ ] 结算统计
- [ ] 结算触发

#### Day 4: Earnings收益测试
**负责人**: 工程师B
**测试文件**: `internal/service/earnings/*_test.go`

**测试场景**:
- [ ] 收益计算
- [ ] 收益趋势
- [ ] 提现申请

#### Day 5: 财务数据一致性测试
**负责人**: 工程师B
**测试文件**: `internal/service/commission/consistency_test.go`

**测试场景**:
- [ ] 订单金额与抽成一致性
- [ ] 结算金额与收益一致性
- [ ] 并发操作数据一致性

### 交付标准
- [ ] 所有测试通过
- [ ] 覆盖率100%
- [ ] 财务计算准确
- [ ] Code Review通过

---

## 👥 Week 6: Player + User模块 (12/28-01/03)

### 本周目标
- 覆盖Player和User模块
- 达到100%覆盖率
- 编写50+测试用例

### 模块分解

#### Day 1-2: Player模块测试
**负责人**: 工程师A
**测试文件**: `internal/service/player/*_test.go`

**测试场景**:
- [ ] 申请成为陪玩师
- [ ] 陪玩师资料管理
- [ ] 在线状态管理
- [ ] 标签管理

#### Day 3-4: User模块测试
**负责人**: 工程师A
**测试文件**: `internal/service/user/*_test.go`

**测试场景**:
- [ ] 用户注册
- [ ] 用户登录
- [ ] 用户资料更新
- [ ] 密码修改

#### Day 5: 角色权限集成测试
**负责人**: 工程师A
**测试文件**: `internal/service/user/role_test.go`

**测试场景**:
- [ ] 用户角色分配
- [ ] 权限检查
- [ ] 角色权限变更

### 交付标准
- [ ] 所有测试通过
- [ ] 覆盖率100%
- [ ] 用户流程完整
- [ ] Code Review通过

---

## 🔐 Week 7: Permission + Role模块 (01/04-01/10)

### 本周目标
- 覆盖Permission和Role模块
- 达到100%覆盖率
- 编写40+测试用例

### 模块分解

#### Day 1-2: Permission模块测试
**负责人**: 工程师B
**测试文件**: `internal/service/permission/*_test.go`

**测试场景**:
- [ ] 创建权限
- [ ] 更新权限
- [ ] 删除权限
- [ ] 权限查询

#### Day 3-4: Role模块测试
**负责人**: 工程师B
**测试文件**: `internal/service/role/*_test.go`

**测试场景**:
- [ ] 创建角色
- [ ] 更新角色
- [ ] 删除角色
- [ ] 角色分配权限

#### Day 5: RBAC集成测试
**负责人**: 工程师B
**测试文件**: `internal/service/auth/rbac_test.go`

**测试场景**:
- [ ] 角色权限验证
- [ ] 用户角色权限链
- [ ] 权限继承

### 交付标准
- [ ] 所有测试通过
- [ ] 覆盖率100%
- [ ] RBAC逻辑完整
- [ ] Code Review通过

---

## 📢 Week 8: 其他模块 (01/11-01/17)

### 本周目标
- 覆盖Notification、Gift、Review等模块
- 达到100%覆盖率
- 编写40+测试用例

### 模块分解

#### Day 1: Notification模块
**负责人**: 工程师A
**测试文件**: `internal/service/notification/*_test.go`

#### Day 2: Gift模块
**负责人**: 工程师A
**测试文件**: `internal/service/gift/*_test.go`

#### Day 3: Review模块
**负责人**: 工程师A
**测试文件**: `internal/service/review/*_test.go`

#### Day 4: Feed模块
**负责人**: 工程师A
**测试文件**: `internal/service/feed/*_test.go`

#### Day 5: 模块集成测试
**负责人**: 工程师A
**测试文件**: `internal/service/integration/*_test.go`

### 交付标准
- [ ] 所有测试通过
- [ ] 覆盖率100%
- [ ] Service层全部完成
- [ ] Code Review通过

---

## 🗄️ Week 9-10: Repository层测试 (01/18-01/31)

### 双周目标
- 覆盖所有Repository层代码
- 达到100%覆盖率
- 编写100+测试用例

### 分工策略

**工程师A**:
- User、Player、Order、Payment Repository
- 重点: 业务相关Repository

**工程师B**:
- Auth、Permission、Role、Commission Repository
- 重点: 系统相关Repository

### 测试策略

```go
// 使用SQLMock测试SQL查询
func (s *UserRepositoryTestSuite) TestListWithComplexQuery() {
    // Mock SQL
    s.sqlMock.ExpectQuery("SELECT \* FROM users").
        WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
        WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).
            AddRow(1, "test@example.com"))
    
    // 执行查询
    users, err := s.repository.List(s.ctx, repository.UserListOptions{
        Email: "test@example.com",
    })
    
    // 验证
    s.NoError(err)
    s.Len(users, 1)
    s.NoError(s.sqlMock.ExpectationsWereMet())
}
```

**关键测试点**:
- [ ] SQL查询正确性
- [ ] 事务提交/回滚
- [ ] 连接池管理
- [ ] 索引使用
- [ ] 批量操作

### 交付标准
- [ ] 所有Repository测试通过
- [ ] 覆盖率100%
- [ ] SQLMock验证正确
- [ ] Code Review通过

---

## 🌐 Week 11: Handler层测试 (02/01-02/07)

### 本周目标
- 覆盖所有Handler层代码
- 达到100%覆盖率
- 编写80+测试用例

### 测试策略

```go
// 测试所有API端点
func (s *UserHandlerTestSuite) TestAPIEndpoints() {
    tests := []struct {
        name           string
        method         string
        path           string
        body           interface{}
        setupMock      func()
        expectedStatus int
        expectedBody   string
    }{
        {
            name:   "创建用户成功",
            method: "POST",
            path:   "/api/v1/users",
            body: map[string]string{
                "email": "test@example.com",
                "name":  "Test User",
            },
            setupMock: func() {
                s.mockSvc.On("CreateUser", mock.Anything, mock.Anything).Return(&model.User{ID: 1}, nil)
            },
            expectedStatus: http.StatusCreated,
            expectedBody:   `"id":1`,
        },
        // ... 更多测试用例
    }
    
    for _, tt := range tests {
        s.Run(tt.name, func() {
            tt.setupMock()
            
            body, _ := json.Marshal(tt.body)
            req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader(body))
            w := httptest.NewRecorder()
            
            s.router.ServeHTTP(w, req)
            
            s.Equal(tt.expectedStatus, w.Code)
            s.Contains(w.Body.String(), tt.expectedBody)
        })
    }
}
```

**关键测试点**:
- [ ] 路由正确性
- [ ] 参数绑定
- [ ] 认证中间件
- [ ] 授权检查
- [ ] 响应格式
- [ ] 错误处理

### 交付标准
- [ ] 所有Handler测试通过
- [ ] 覆盖率100%
- [ ] API文档同步
- [ ] Code Review通过

---

## 🎯 Week 12: 边缘覆盖+优化 (02/08-02/15)

### 本周目标
- 达到100%覆盖率
- 优化测试性能
- 完善测试文档
- 团队培训

### 任务分解

#### Day 1-2: 覆盖率补全
**负责人**: 测试架构师

- [ ] 生成覆盖率报告
- [ ] 识别未覆盖代码
- [ ] 补充边缘测试
- [ ] 修复测试间隙

#### Day 3-4: 测试优化
**负责人**: 测试架构师

- [ ] 优化测试执行速度
- [ ] 减少重复测试
- [ ] 提高Mock效率
- [ ] 并行测试配置

#### Day 5: 文档与培训
**负责人**: 测试架构师

- [ ] 编写测试最佳实践文档
- [ ] 准备团队培训材料
- [ ] 进行测试培训
- [ ] 制定测试维护规范

### 交付标准
- [ ] 覆盖率100%
- [ ] 测试执行时间<5分钟
- [ ] 团队培训完成
- [ ] 测试规范文档化

---

## 📈 进度跟踪

### 每日站会
- **时间**: 每天9:00-9:15
- **内容**: 昨日进度、今日计划、阻塞问题
- **参与**: 所有测试相关人员

### 每周评审
- **时间**: 每周五17:00-18:00
- **内容**: 本周总结、下周计划、风险识别
- **参与**: 测试架构师 + 工程师A + 工程师B

### 覆盖率监控
```bash
# 每日覆盖率检查脚本
#!/bin/bash
echo "=== Daily Coverage Check ==="
go test ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -func=coverage.out | grep total

# 按模块检查
echo "=== Module Coverage ==="
for module in model repository service handler; do
    echo "Module: internal/$module"
    go test ./internal/$module/... -cover 2>&1 | grep coverage || echo "No tests"
done
```

---

## 🚨 风险管理

### 风险1: 测试代码质量不高
**应对措施**:
- Code Review必须检查测试代码
- 测试代码覆盖率要求100%
- 定期进行测试重构

### 风险2: 测试执行速度慢
**应对措施**:
- 使用内存数据库
- 并行测试执行
- Mock外部依赖

### 风险3: 团队成员不熟悉测试
**应对措施**:
- Week 1进行培训
- 提供详细模板
- 代码审查时指导

### 风险4: 业务需求变更导致测试失效
**应对措施**:
- 测试代码与业务代码同步审查
- 建立测试维护规范
- 定期重构测试代码

---

## 📞 沟通机制

### 问题升级路径
1. **技术问题**: 工程师 → 测试架构师
2. **进度问题**: 工程师 → 项目经理
3. **资源问题**: 测试架构师 → 技术总监

### 紧急联系
- **测试架构师**: [电话/邮箱]
- **工程师A**: [电话/邮箱]
- **工程师B**: [电话/邮箱]
- **项目经理**: [电话/邮箱]

---

## ✅ 完成标准

### Week 1完成标准
- [ ] 所有测试依赖安装完成
- [ ] 测试工具包实现完成
- [ ] CI/CD配置完成
- [ ] Mock生成完成
- [ ] 团队培训完成

### Week 2-12完成标准
- [ ] 本周测试代码编写完成
- [ ] 本周覆盖率目标达成
- [ ] Code Review通过
- [ ] 测试文档更新
- [ ] 无阻塞问题

### 项目完成标准
- [ ] 整体覆盖率100%
- [ ] 所有测试通过
- [ ] CI/CD稳定运行
- [ ] 团队熟练掌握测试
- [ ] 测试规范文档化

---

**计划版本**: 1.0
**最后更新**: 2025-11-23
**下次更新**: 每周五评审后更新
