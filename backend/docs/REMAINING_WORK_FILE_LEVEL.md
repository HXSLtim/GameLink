# 测试覆盖率提升 - 精确到文件级别的剩余工作清单

**生成时间**: 2025-11-08  
**当前覆盖率**: 35.5%  
**目标覆盖率**: 80.0%  
**剩余工作量**: 44.5%

---

## 📋 文件级别执行清单

### 优先级1: Admin Handler层 (0% → 50%) [预计+8%总体覆盖率]

#### 1.1 `internal/handler/admin/game_test.go` (新建)
**当前**: 0% (6个方法)  
**目标**: 80%  
**需要测试的方法**:
- [ ] `TestGameHandler_ListGames` - 列表查询
- [ ] `TestGameHandler_GetGame` - 详情查询
- [ ] `TestGameHandler_CreateGame` - 创建游戏
- [ ] `TestGameHandler_UpdateGame` - 更新游戏
- [ ] `TestGameHandler_DeleteGame` - 删除游戏
- [ ] `TestGameHandler_ListGameLogs` - 操作日志
- [ ] `TestGameHandler_ListGames_Pagination` - 分页测试
- [ ] `TestGameHandler_CreateGame_Validation` - 参数验证
- [ ] `TestGameHandler_GetGame_NotFound` - 404处理

**预计时间**: 1.5小时

---

#### 1.2 `internal/handler/admin/user_test.go` (新建)
**当前**: 0% (11个方法)  
**目标**: 70%  
**需要测试的方法**:
- [ ] `TestUserHandler_ListUsers` - 用户列表
- [ ] `TestUserHandler_GetUser` - 用户详情
- [ ] `TestUserHandler_CreateUser` - 创建用户
- [ ] `TestUserHandler_UpdateUser` - 更新用户
- [ ] `TestUserHandler_DeleteUser` - 删除用户
- [ ] `TestUserHandler_UpdateUserStatus` - 状态更新
- [ ] `TestUserHandler_UpdateUserRole` - 角色更新
- [ ] `TestUserHandler_ListUserOrders` - 用户订单列表
- [ ] `TestUserHandler_CreateUserWithPlayer` - 创建用户+陪玩师
- [ ] `TestUserHandler_ListUserLogs` - 操作日志
- [ ] `TestUserHandler_ListUsers_Pagination` - 分页
- [ ] `TestUserHandler_CreateUser_Validation` - 参数验证

**预计时间**: 2小时

---

#### 1.3 `internal/handler/admin/player_test.go` (新建)
**当前**: 0% (9个方法)  
**目标**: 70%  
**需要测试的方法**:
- [ ] `TestPlayerHandler_ListPlayers` - 陪玩师列表
- [ ] `TestPlayerHandler_GetPlayer` - 陪玩师详情
- [ ] `TestPlayerHandler_CreatePlayer` - 创建陪玩师
- [ ] `TestPlayerHandler_UpdatePlayer` - 更新陪玩师
- [ ] `TestPlayerHandler_DeletePlayer` - 删除陪玩师
- [ ] `TestPlayerHandler_UpdatePlayerVerification` - 认证状态
- [ ] `TestPlayerHandler_UpdatePlayerGames` - 游戏列表
- [ ] `TestPlayerHandler_UpdatePlayerSkillTags` - 技能标签
- [ ] `TestPlayerHandler_ListPlayerLogs` - 操作日志
- [ ] `TestPlayerHandler_ListPlayers_Pagination` - 分页

**预计时间**: 1.5小时

---

#### 1.4 `internal/handler/admin/order_test.go` (新建)
**当前**: 0% (16个方法)  
**目标**: 60%  
**需要测试的方法**:
- [ ] `TestOrderHandler_CreateOrder` - 创建订单
- [ ] `TestOrderHandler_AssignOrder` - 分配订单
- [ ] `TestOrderHandler_ConfirmOrder` - 确认订单
- [ ] `TestOrderHandler_StartOrder` - 开始订单
- [ ] `TestOrderHandler_CompleteOrder` - 完成订单
- [ ] `TestOrderHandler_RefundOrder` - 退款订单
- [ ] `TestOrderHandler_ListOrders` - 订单列表
- [ ] `TestOrderHandler_GetOrder` - 订单详情
- [ ] `TestOrderHandler_UpdateOrder` - 更新订单
- [ ] `TestOrderHandler_DeleteOrder` - 删除订单
- [ ] `TestOrderHandler_GetOrderTimeline` - 订单时间线
- [ ] `TestOrderHandler_ListOrderPayments` - 支付列表
- [ ] `TestOrderHandler_ListOrderRefunds` - 退款列表
- [ ] `TestOrderHandler_ListOrderReviews` - 评价列表
- [ ] `TestOrderHandler_ListOrderLogs` - 操作日志
- [ ] `TestOrderHandler_CancelOrder` - 取消订单
- [ ] `TestOrderHandler_ReviewOrder` - 评价订单

**预计时间**: 3小时

---

#### 1.5 `internal/handler/admin/payment_test.go` (新建，从order.go拆分)
**当前**: 0% (8个方法)  
**目标**: 60%  
**需要测试的方法**:
- [ ] `TestPaymentHandler_CreatePayment` - 创建支付
- [ ] `TestPaymentHandler_CapturePayment` - 捕获支付
- [ ] `TestPaymentHandler_ListPayments` - 支付列表
- [ ] `TestPaymentHandler_GetPayment` - 支付详情
- [ ] `TestPaymentHandler_UpdatePayment` - 更新支付
- [ ] `TestPaymentHandler_DeletePayment` - 删除支付
- [ ] `TestPaymentHandler_RefundPayment` - 退款支付
- [ ] `TestPaymentHandler_ListPaymentLogs` - 操作日志

**预计时间**: 1.5小时

---

#### 1.6 `internal/handler/admin/review_test.go` (新建)
**当前**: 0% (7个方法)  
**目标**: 60%  
**需要测试的方法**:
- [ ] `TestReviewHandler_ListReviews` - 评价列表
- [ ] `TestReviewHandler_GetReview` - 评价详情
- [ ] `TestReviewHandler_CreateReview` - 创建评价
- [ ] `TestReviewHandler_UpdateReview` - 更新评价
- [ ] `TestReviewHandler_DeleteReview` - 删除评价
- [ ] `TestReviewHandler_ListPlayerReviews` - 陪玩师评价
- [ ] `TestReviewHandler_ListReviewLogs` - 操作日志

**预计时间**: 1小时

---

#### 1.7 `internal/handler/admin/role_test.go` (新建)
**当前**: 0% (9个方法)  
**目标**: 60%  
**需要测试的方法**:
- [ ] `TestRoleHandler_ListRoles` - 角色列表
- [ ] `TestRoleHandler_GetRole` - 角色详情
- [ ] `TestRoleHandler_CreateRole` - 创建角色
- [ ] `TestRoleHandler_UpdateRole` - 更新角色
- [ ] `TestRoleHandler_DeleteRole` - 删除角色
- [ ] `TestRoleHandler_AssignPermissions` - 分配权限
- [ ] `TestRoleHandler_AssignRolesToUser` - 分配角色给用户
- [ ] `TestRoleHandler_GetUserRoles` - 获取用户角色

**预计时间**: 1小时

---

#### 1.8 `internal/handler/admin/permission_test.go` (新建)
**当前**: 0% (8个方法)  
**目标**: 60%  
**需要测试的方法**:
- [ ] `TestPermissionHandler_ListPermissions` - 权限列表
- [ ] `TestPermissionHandler_GetPermission` - 权限详情
- [ ] `TestPermissionHandler_CreatePermission` - 创建权限
- [ ] `TestPermissionHandler_UpdatePermission` - 更新权限
- [ ] `TestPermissionHandler_DeletePermission` - 删除权限
- [ ] `TestPermissionHandler_GetRolePermissions` - 角色权限
- [ ] `TestPermissionHandler_GetUserPermissions` - 用户权限
- [ ] `TestPermissionHandler_GetPermissionGroups` - 权限分组

**预计时间**: 1小时

---

#### 1.9 `internal/handler/admin/helpers_test.go` (新建)
**当前**: 0% (12个辅助函数)  
**目标**: 80%  
**需要测试的方法**:
- [ ] `TestParseUintParam` - 解析uint参数
- [ ] `TestQueryIntDefault` - 查询int默认值
- [ ] `TestQueryUint64Ptr` - 查询uint64指针
- [ ] `TestQueryTimePtr` - 查询时间指针
- [ ] `TestParseCSVParams` - 解析CSV参数
- [ ] `TestWriteJSON` - 写入JSON
- [ ] `TestWriteJSONError` - 写入JSON错误
- [ ] `TestParsePagination` - 解析分页
- [ ] `TestBuildOrderListOptions` - 构建订单列表选项
- [ ] `TestBuildPaymentListOptions` - 构建支付列表选项
- [ ] `TestBuildUserListOptions` - 构建用户列表选项
- [ ] `TestNormalizeOrderStatus` - 标准化订单状态

**预计时间**: 1小时

---

#### 1.10 `internal/handler/admin/stats_handler_test.go` (新建)
**当前**: 0% (7个方法)  
**目标**: 50%  
**需要测试的方法**:
- [ ] `TestStatsHandler_Dashboard` - 仪表盘
- [ ] `TestStatsHandler_RevenueTrend` - 收入趋势
- [ ] `TestStatsHandler_UserGrowth` - 用户增长
- [ ] `TestStatsHandler_OrdersSummary` - 订单摘要
- [ ] `TestStatsHandler_TopPlayers` - 顶级陪玩师
- [ ] `TestStatsHandler_AuditOverview` - 审计概览
- [ ] `TestStatsHandler_AuditTrend` - 审计趋势

**预计时间**: 1小时

---

#### 1.11 `internal/handler/admin/system_handler_test.go` (新建)
**当前**: 0% (5个方法)  
**目标**: 50%  
**需要测试的方法**:
- [ ] `TestSystemInfoHandler_Config` - 配置信息
- [ ] `TestSystemInfoHandler_DBStatus` - 数据库状态
- [ ] `TestSystemInfoHandler_CacheStatus` - 缓存状态
- [ ] `TestSystemInfoHandler_Resources` - 资源信息
- [ ] `TestSystemInfoHandler_Version` - 版本信息

**预计时间**: 0.5小时

---

#### 1.12 其他Admin Handler文件 (可选，低优先级)
- `internal/handler/admin/commission_test.go` - 佣金管理 (4个方法)
- `internal/handler/admin/dashboard_test.go` - 仪表盘 (4个方法)
- `internal/handler/admin/item_test.go` - 服务项管理 (7个方法)
- `internal/handler/admin/ranking_test.go` - 排行榜佣金 (5个方法)
- `internal/handler/admin/stats_test.go` - 统计分析 (4个方法)
- `internal/handler/admin/withdraw_test.go` - 提现管理 (6个方法)

**预计时间**: 2小时 (可选)

**Admin Handler层总计**: 13-15小时

---

### 优先级2: Service层完善 (平均64% → 78%) [预计+8%总体覆盖率]

#### 2.1 `internal/service/admin/admin_test.go` (增强)
**当前**: 40.7% (56个方法，约23个已测试)  
**目标**: 70%  
**需要新增的测试方法**:
- [ ] `TestService_GetOrderPayments` - 获取订单支付列表
- [ ] `TestService_GetOrderRefunds` - 获取订单退款列表
- [ ] `TestService_GetOrderReviews` - 获取订单评价列表
- [ ] `TestService_GetOrderTimeline` - 获取订单时间线
- [ ] `TestService_ListOperationLogs` - 操作日志列表
- [ ] `TestService_UpdateOrder` - 更新订单 (边界条件)
- [ ] `TestService_UpdatePayment` - 更新支付 (边界条件)
- [ ] `TestService_ListUsersWithOptions` - 带选项的用户列表
- [ ] `TestService_UpdatePlayerSkillTags` - 更新技能标签 (需要TxManager)
- [ ] `TestService_RegisterUserAndPlayer` - 注册用户和陪玩师 (需要TxManager)

**预计时间**: 2小时

---

#### 2.2 `internal/service/role/role_test.go` (增强)
**当前**: 55.5%  
**目标**: 80%  
**需要新增的测试方法**:
- [ ] `TestRoleService_ListRolesPagedWithFilter` - 带过滤的分页列表
- [ ] `TestRoleService_GetRoleWithPermissions` - 获取角色及权限
- [ ] `TestRoleService_CreateRole_Validation` - 创建角色验证
- [ ] `TestRoleService_UpdateRole_SystemRole` - 更新系统角色
- [ ] `TestRoleService_AssignPermissionsToRole` - 分配权限给角色
- [ ] `TestRoleService_RemovePermissionsFromRole` - 移除角色权限
- [ ] `TestRoleService_AssignRolesToUser` - 分配角色给用户
- [ ] `TestRoleService_RemoveRolesFromUser` - 移除用户角色

**预计时间**: 1.5小时

---

#### 2.3 `internal/service/player/player_test.go` (增强)
**当前**: 66.8%  
**目标**: 80%  
**需要新增的测试方法**:
- [ ] `TestPlayerService_GetPlayerOrderCount` - 获取订单数量
- [ ] `TestPlayerService_GetPlayerStats` - 获取统计数据
- [ ] `TestPlayerService_GetPlayerReviews` - 获取评价列表
- [ ] `TestPlayerService_CalculateGoodRatio` - 计算好评率
- [ ] `TestPlayerService_CalculateAvgResponseTime` - 计算平均响应时间
- [ ] `TestPlayerService_CalculateRepeatRate` - 计算复购率

**预计时间**: 1小时

---

#### 2.4 `internal/service/order/order_test.go` (增强)
**当前**: 67.8%  
**目标**: 85%  
**需要新增的测试方法**:
- [ ] `TestOrderService_GetOrderPayments` - 获取订单支付
- [ ] `TestOrderService_GetOrderRefunds` - 获取订单退款
- [ ] `TestOrderService_GetOrderReviews` - 获取订单评价
- [ ] `TestOrderService_GetOrderTimeline` - 获取订单时间线
- [ ] `TestOrderService_CancelOrder_EdgeCases` - 取消订单边界条件
- [ ] `TestOrderService_RefundOrder_EdgeCases` - 退款订单边界条件

**预计时间**: 1小时

**Service层总计**: 5.5小时

---

### 优先级3: User/Player Handler层 (39% → 70%) [预计+5%总体覆盖率]

#### 3.1 `internal/handler/user/order_test.go` (增强)
**当前**: 39%  
**目标**: 70%  
**需要新增的测试方法**:
- [ ] `TestUserOrderHandler_CreateOrder_Validation` - 创建订单验证
- [ ] `TestUserOrderHandler_CancelOrder` - 取消订单
- [ ] `TestUserOrderHandler_GetOrderDetails` - 获取订单详情
- [ ] `TestUserOrderHandler_ListOrders_Filter` - 列表过滤
- [ ] `TestUserOrderHandler_ListOrders_Pagination` - 分页测试

**预计时间**: 1小时

---

#### 3.2 `internal/handler/user/payment_test.go` (增强)
**当前**: 39%  
**目标**: 70%  
**需要新增的测试方法**:
- [ ] `TestUserPaymentHandler_CreatePayment_Validation` - 创建支付验证
- [ ] `TestUserPaymentHandler_GetPaymentStatus` - 获取支付状态
- [ ] `TestUserPaymentHandler_ListPayments_Filter` - 列表过滤

**预计时间**: 0.5小时

---

#### 3.3 `internal/handler/user/review_test.go` (增强)
**当前**: 39%  
**目标**: 70%  
**需要新增的测试方法**:
- [ ] `TestUserReviewHandler_CreateReview_Validation` - 创建评价验证
- [ ] `TestUserReviewHandler_UpdateReview` - 更新评价
- [ ] `TestUserReviewHandler_ListReviews_Filter` - 列表过滤

**预计时间**: 0.5小时

---

#### 3.4 `internal/handler/user/player_test.go` (增强)
**当前**: 39%  
**目标**: 70%  
**需要新增的测试方法**:
- [ ] `TestUserPlayerHandler_SearchPlayers_Filter` - 搜索过滤
- [ ] `TestUserPlayerHandler_GetPlayerDetails` - 获取详情
- [ ] `TestUserPlayerHandler_ListPlayers_Pagination` - 分页测试

**预计时间**: 0.5小时

---

#### 3.5 `internal/handler/player/*_test.go` (增强)
**当前**: 39%  
**目标**: 70%  
**需要新增的测试方法**:
- [ ] `TestPlayerOrderHandler_*` - 订单相关测试
- [ ] `TestPlayerProfileHandler_*` - 资料相关测试
- [ ] `TestPlayerEarningsHandler_*` - 收益相关测试

**预计时间**: 1小时

**Handler层总计**: 3.5小时

---

### 优先级4: 小模块批量提升 [预计+5%总体覆盖率]

#### 4.1 `internal/cache/redis_test.go` (新建)
**当前**: 0% (4个方法)  
**目标**: 60%  
**需要测试的方法**:
- [ ] `TestRedisCache_Get` - 获取缓存
- [ ] `TestRedisCache_Set` - 设置缓存
- [ ] `TestRedisCache_Delete` - 删除缓存
- [ ] `TestRedisCache_Close` - 关闭连接

**预计时间**: 0.5小时

---

#### 4.2 `internal/auth/jwt_test.go` (增强)
**当前**: 部分覆盖  
**目标**: 85%  
**需要新增的测试方法**:
- [ ] `TestJWTManager_RefreshToken` - 刷新Token (需要查看实际方法签名)
- [ ] `TestJWTManager_IsTokenExpired` - Token过期检查 (如果方法存在)
- [ ] `TestJWTManager_GetTokenRemainingTime` - Token剩余时间 (如果方法存在)

**预计时间**: 0.5小时 (需要先确认方法是否存在)

---

#### 4.3 `internal/db/db_test.go` (新建)
**当前**: 30.9%  
**目标**: 60%  
**需要测试的方法**:
- [ ] `TestOpen_Postgres` - PostgreSQL连接
- [ ] `TestOpen_SQLite` - SQLite连接
- [ ] `TestOpen_Error` - 连接错误处理

**预计时间**: 0.5小时

---

#### 4.4 `internal/db/seed_test.go` (新建)
**当前**: 0%  
**目标**: 50%  
**需要测试的方法**:
- [ ] `TestApplySeeds` - 应用种子数据
- [ ] `TestSeedGames` - 种子游戏数据
- [ ] `TestSeedUser` - 种子用户数据
- [ ] `TestSeedPlayer` - 种子陪玩师数据

**预计时间**: 0.5小时

---

#### 4.5 `internal/logging/logger_test.go` (增强)
**当前**: 29.2%  
**目标**: 60%  
**需要新增的测试方法**:
- [ ] `TestLogger_WithContext` - 上下文日志
- [ ] `TestLogger_WithFields` - 字段日志
- [ ] `TestLogger_Error` - 错误日志
- [ ] `TestLogger_Warn` - 警告日志

**预计时间**: 0.5小时

---

#### 4.6 `internal/metrics/metrics_test.go` (增强)
**当前**: 19.2%  
**目标**: 50%  
**需要新增的测试方法**:
- [ ] `TestMetrics_Increment` - 增加计数
- [ ] `TestMetrics_RecordDuration` - 记录时长
- [ ] `TestMetrics_RecordGauge` - 记录仪表

**预计时间**: 0.5小时

---

#### 4.7 `internal/config/env_test.go` (增强)
**当前**: 61.1%  
**目标**: 75%  
**需要新增的测试方法**:
- [ ] `TestConfig_LoadFromFile` - 从文件加载
- [ ] `TestConfig_OverrideFromEnv` - 环境变量覆盖
- [ ] `TestConfig_ValidateCrypto` - 加密配置验证

**预计时间**: 0.5小时

**小模块总计**: 3.5小时

---

### 优先级5: Repository层补充 [预计+3%总体覆盖率]

#### 5.1 `internal/repository/commission/repository_test.go` (增强)
**当前**: 76.9%  
**目标**: 90%  
**需要新增的测试方法**:
- [ ] `TestCommissionRepository_GetRuleForOrder_EdgeCases` - 边界条件
- [ ] `TestCommissionRepository_GetSettlement_EdgeCases` - 边界条件
- [ ] `TestCommissionRepository_UpdateRecord_EdgeCases` - 边界条件

**预计时间**: 0.5小时

---

#### 5.2 `internal/repository/serviceitem/repository_test.go` (增强)
**当前**: 78.2%  
**目标**: 90%  
**需要新增的测试方法**:
- [ ] `TestServiceItemRepository_List_WithFilters` - 带过滤的列表
- [ ] `TestServiceItemRepository_GetGameServices_EdgeCases` - 边界条件

**预计时间**: 0.5小时

---

#### 5.3 `internal/repository/permission/repository_test.go` (增强)
**当前**: 63.0%  
**目标**: 85%  
**需要新增的测试方法**:
- [ ] `TestPermissionRepository_GetBySlug` - 通过slug获取
- [ ] `TestPermissionRepository_ListByGroup` - 按分组列表
- [ ] `TestPermissionRepository_GetUserPermissions` - 用户权限

**预计时间**: 0.5小时

**Repository层总计**: 1.5小时

---

## 📊 工作量汇总

| 优先级 | 模块 | 文件数 | 预计时间 | 覆盖率提升 |
|--------|------|--------|----------|-----------|
| P1 | Admin Handler | 12 | 13-15h | +8% |
| P2 | Service层 | 4 | 5.5h | +8% |
| P3 | User/Player Handler | 5 | 3.5h | +5% |
| P4 | 小模块 | 7 | 3.5h | +5% |
| P5 | Repository层 | 3 | 1.5h | +3% |
| **总计** | | **31** | **27-29h** | **+29%** |

**预计最终覆盖率**: 35.5% + 29% = **64.5%**

**注意**: 要达到80%，还需要额外15.5%的提升，可能需要：
- 更深入的边界条件测试
- 集成测试
- 错误场景测试
- 性能测试

---

## 🎯 执行建议

### 阶段1: 快速提升 (8-10小时)
1. Service层完善 (5.5h) → +8%
2. 小模块批量 (3.5h) → +5%
3. Repository补充 (1.5h) → +3%
**结果**: 35.5% → 52% (+16.5%)

### 阶段2: Handler层核心 (10-12小时)
1. Admin Handler核心文件 (10h) → +6%
2. User/Player Handler (3.5h) → +5%
**结果**: 52% → 63% (+11%)

### 阶段3: 最后冲刺 (8-10小时)
1. Admin Handler剩余文件 (5h) → +2%
2. 边界条件和错误场景 (3-5h) → +2%
**结果**: 63% → 67% (+4%)

### 阶段4: 达到80% (额外10-12小时)
1. 集成测试
2. 性能测试
3. 全面错误场景
**结果**: 67% → 80% (+13%)

---

## 📝 测试文件创建模板

每个Handler测试文件应包含：

```go
package admin

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock AdminService
type mockAdminService struct {
    mock.Mock
}

func (m *mockAdminService) ListGames(ctx context.Context) ([]model.Game, error) {
    args := m.Called(ctx)
    return args.Get(0).([]model.Game), args.Error(1)
}

// ... 其他方法mock

func TestGameHandler_ListGames(t *testing.T) {
    // Arrange
    mockService := new(mockAdminService)
    handler := NewGameHandler(mockService)
    
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.GET("/games", handler.ListGames)
    
    mockService.On("ListGames", mock.Anything).Return([]model.Game{}, nil)
    
    // Act
    req := httptest.NewRequest("GET", "/games", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
    mockService.AssertExpectations(t)
}
```

---

## ✅ 检查清单

每完成一个文件，检查：
- [ ] 所有测试通过
- [ ] 覆盖率达标
- [ ] 代码无编译错误
- [ ] 测试命名规范
- [ ] 测试覆盖主要场景
- [ ] 包含错误场景测试

---

**最后更新**: 2025-11-08  
**文档版本**: 1.0

