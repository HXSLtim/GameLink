# GitHub Issues 模板

## Issue 1: [测试] Admin Handler层核心测试

**分配给**: @成员A  
**优先级**: P0 (最高)  
**预计时间**: 13.5小时  
**预计覆盖率提升**: +8%

### 任务清单

#### 文件1: `internal/handler/admin/game_test.go` (新建) - 1.5h
- [ ] `TestGameHandler_ListGames` - 列表查询
- [ ] `TestGameHandler_GetGame` - 详情查询
- [ ] `TestGameHandler_CreateGame` - 创建游戏
- [ ] `TestGameHandler_UpdateGame` - 更新游戏
- [ ] `TestGameHandler_DeleteGame` - 删除游戏
- [ ] `TestGameHandler_ListGameLogs` - 操作日志
- [ ] `TestGameHandler_ListGames_Pagination` - 分页测试
- [ ] `TestGameHandler_CreateGame_Validation` - 参数验证
- [ ] `TestGameHandler_GetGame_NotFound` - 404处理

#### 文件2: `internal/handler/admin/user_test.go` (新建) - 2h
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

#### 文件3: `internal/handler/admin/player_test.go` (新建) - 1.5h
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

#### 文件4: `internal/handler/admin/order_test.go` (新建) - 3h
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

#### 文件5: `internal/handler/admin/payment_test.go` (新建) - 1.5h
- [ ] `TestPaymentHandler_CreatePayment` - 创建支付
- [ ] `TestPaymentHandler_CapturePayment` - 捕获支付
- [ ] `TestPaymentHandler_ListPayments` - 支付列表
- [ ] `TestPaymentHandler_GetPayment` - 支付详情
- [ ] `TestPaymentHandler_UpdatePayment` - 更新支付
- [ ] `TestPaymentHandler_DeletePayment` - 删除支付
- [ ] `TestPaymentHandler_RefundPayment` - 退款支付
- [ ] `TestPaymentHandler_ListPaymentLogs` - 操作日志

#### 文件6: `internal/handler/admin/review_test.go` (新建) - 1h
- [ ] `TestReviewHandler_ListReviews` - 评价列表
- [ ] `TestReviewHandler_GetReview` - 评价详情
- [ ] `TestReviewHandler_CreateReview` - 创建评价
- [ ] `TestReviewHandler_UpdateReview` - 更新评价
- [ ] `TestReviewHandler_DeleteReview` - 删除评价
- [ ] `TestReviewHandler_ListPlayerReviews` - 陪玩师评价
- [ ] `TestReviewHandler_ListReviewLogs` - 操作日志

#### 文件7: `internal/handler/admin/role_test.go` (新建) - 1h
- [ ] `TestRoleHandler_ListRoles` - 角色列表
- [ ] `TestRoleHandler_GetRole` - 角色详情
- [ ] `TestRoleHandler_CreateRole` - 创建角色
- [ ] `TestRoleHandler_UpdateRole` - 更新角色
- [ ] `TestRoleHandler_DeleteRole` - 删除角色
- [ ] `TestRoleHandler_AssignPermissions` - 分配权限
- [ ] `TestRoleHandler_AssignRolesToUser` - 分配角色给用户
- [ ] `TestRoleHandler_GetUserRoles` - 获取用户角色

#### 文件8: `internal/handler/admin/permission_test.go` (新建) - 1h
- [ ] `TestPermissionHandler_ListPermissions` - 权限列表
- [ ] `TestPermissionHandler_GetPermission` - 权限详情
- [ ] `TestPermissionHandler_CreatePermission` - 创建权限
- [ ] `TestPermissionHandler_UpdatePermission` - 更新权限
- [ ] `TestPermissionHandler_DeletePermission` - 删除权限
- [ ] `TestPermissionHandler_GetRolePermissions` - 角色权限
- [ ] `TestPermissionHandler_GetUserPermissions` - 用户权限
- [ ] `TestPermissionHandler_GetPermissionGroups` - 权限分组

#### 文件9: `internal/handler/admin/helpers_test.go` (新建) - 1h
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

### 参考文档
- 详细任务: `backend/REMAINING_WORK_FILE_LEVEL.md`
- 测试规范: `.cursor/rules/backend-testing.mdc`
- 示例代码: `backend/internal/handler/health_test.go`

### 验收标准
- [ ] 所有测试通过
- [ ] 每个文件覆盖率 ≥ 60%
- [ ] 代码审查通过
- [ ] 无编译错误

---

## Issue 2: [测试] Service层测试增强

**分配给**: @成员B  
**优先级**: P0 (最高)  
**预计时间**: 5.5小时  
**预计覆盖率提升**: +8%

### 任务清单

#### 文件1: `internal/service/admin/admin_test.go` (增强) - 2h
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

#### 文件2: `internal/service/role/role_test.go` (增强) - 1.5h
- [ ] `TestRoleService_ListRolesPagedWithFilter` - 带过滤的分页列表
- [ ] `TestRoleService_GetRoleWithPermissions` - 获取角色及权限
- [ ] `TestRoleService_CreateRole_Validation` - 创建角色验证
- [ ] `TestRoleService_UpdateRole_SystemRole` - 更新系统角色
- [ ] `TestRoleService_AssignPermissionsToRole` - 分配权限给角色
- [ ] `TestRoleService_RemovePermissionsFromRole` - 移除角色权限
- [ ] `TestRoleService_AssignRolesToUser` - 分配角色给用户
- [ ] `TestRoleService_RemoveRolesFromUser` - 移除用户角色

#### 文件3: `internal/service/player/player_test.go` (增强) - 1h
- [ ] `TestPlayerService_GetPlayerOrderCount` - 获取订单数量
- [ ] `TestPlayerService_GetPlayerStats` - 获取统计数据
- [ ] `TestPlayerService_GetPlayerReviews` - 获取评价列表
- [ ] `TestPlayerService_CalculateGoodRatio` - 计算好评率
- [ ] `TestPlayerService_CalculateAvgResponseTime` - 计算平均响应时间
- [ ] `TestPlayerService_CalculateRepeatRate` - 计算复购率

#### 文件4: `internal/service/order/order_test.go` (增强) - 1h
- [ ] `TestOrderService_GetOrderPayments` - 获取订单支付
- [ ] `TestOrderService_GetOrderRefunds` - 获取订单退款
- [ ] `TestOrderService_GetOrderReviews` - 获取订单评价
- [ ] `TestOrderService_GetOrderTimeline` - 获取订单时间线
- [ ] `TestOrderService_CancelOrder_EdgeCases` - 取消订单边界条件
- [ ] `TestOrderService_RefundOrder_EdgeCases` - 退款订单边界条件

### 参考文档
- 现有测试: `backend/internal/service/admin/admin_test.go`
- 示例代码: `backend/internal/service/earnings/earnings_test.go`
- 测试规范: `.cursor/rules/backend-testing.mdc`

### 验收标准
- [ ] 所有测试通过
- [ ] Admin Service覆盖率 ≥ 70%
- [ ] Role Service覆盖率 ≥ 80%
- [ ] Player Service覆盖率 ≥ 80%
- [ ] Order Service覆盖率 ≥ 85%

---

## Issue 3: [测试] User/Player Handler增强

**分配给**: @成员A 或 @成员B  
**优先级**: P1 (高)  
**预计时间**: 3.5小时  
**预计覆盖率提升**: +5%

### 任务清单

#### 文件1: `internal/handler/user/order_test.go` (增强) - 1h
- [ ] `TestUserOrderHandler_CreateOrder_Validation` - 创建订单验证
- [ ] `TestUserOrderHandler_CancelOrder` - 取消订单
- [ ] `TestUserOrderHandler_GetOrderDetails` - 获取订单详情
- [ ] `TestUserOrderHandler_ListOrders_Filter` - 列表过滤
- [ ] `TestUserOrderHandler_ListOrders_Pagination` - 分页测试

#### 文件2: `internal/handler/user/payment_test.go` (增强) - 0.5h
- [ ] `TestUserPaymentHandler_CreatePayment_Validation` - 创建支付验证
- [ ] `TestUserPaymentHandler_GetPaymentStatus` - 获取支付状态
- [ ] `TestUserPaymentHandler_ListPayments_Filter` - 列表过滤

#### 文件3: `internal/handler/user/review_test.go` (增强) - 0.5h
- [ ] `TestUserReviewHandler_CreateReview_Validation` - 创建评价验证
- [ ] `TestUserReviewHandler_UpdateReview` - 更新评价
- [ ] `TestUserReviewHandler_ListReviews_Filter` - 列表过滤

#### 文件4: `internal/handler/user/player_test.go` (增强) - 0.5h
- [ ] `TestUserPlayerHandler_SearchPlayers_Filter` - 搜索过滤
- [ ] `TestUserPlayerHandler_GetPlayerDetails` - 获取详情
- [ ] `TestUserPlayerHandler_ListPlayers_Pagination` - 分页测试

#### 文件5: `internal/handler/player/*_test.go` (增强) - 1h
- [ ] `TestPlayerOrderHandler_*` - 订单相关测试
- [ ] `TestPlayerProfileHandler_*` - 资料相关测试
- [ ] `TestPlayerEarningsHandler_*` - 收益相关测试

### 验收标准
- [ ] 所有测试通过
- [ ] 每个文件覆盖率 ≥ 70%

---

## Issue 4: [测试] Repository层补充

**分配给**: @成员B  
**优先级**: P1 (高)  
**预计时间**: 1.5小时  
**预计覆盖率提升**: +3%

### 任务清单

#### 文件1: `internal/repository/commission/repository_test.go` (增强) - 0.5h
- [ ] `TestCommissionRepository_GetRuleForOrder_EdgeCases` - 边界条件
- [ ] `TestCommissionRepository_GetSettlement_EdgeCases` - 边界条件
- [ ] `TestCommissionRepository_UpdateRecord_EdgeCases` - 边界条件

#### 文件2: `internal/repository/serviceitem/repository_test.go` (增强) - 0.5h
- [ ] `TestServiceItemRepository_List_WithFilters` - 带过滤的列表
- [ ] `TestServiceItemRepository_GetGameServices_EdgeCases` - 边界条件

#### 文件3: `internal/repository/permission/repository_test.go` (增强) - 0.5h
- [ ] `TestPermissionRepository_GetBySlug` - 通过slug获取
- [ ] `TestPermissionRepository_ListByGroup` - 按分组列表
- [ ] `TestPermissionRepository_GetUserPermissions` - 用户权限

### 验收标准
- [ ] 所有测试通过
- [ ] 每个文件覆盖率 ≥ 85%

---

## Issue 5: [测试] 小模块批量测试

**分配给**: @成员C  
**优先级**: P1 (高)  
**预计时间**: 3.5小时  
**预计覆盖率提升**: +5%

### 任务清单

#### 文件1: `internal/cache/redis_test.go` (新建) - 0.5h
- [ ] `TestRedisCache_Get` - 获取缓存
- [ ] `TestRedisCache_Set` - 设置缓存
- [ ] `TestRedisCache_Delete` - 删除缓存
- [ ] `TestRedisCache_Close` - 关闭连接

#### 文件2: `internal/auth/jwt_test.go` (增强) - 0.5h
- [ ] `TestJWTManager_RefreshToken` - 刷新Token
- [ ] `TestJWTManager_IsTokenExpired` - Token过期检查
- [ ] `TestJWTManager_GetTokenRemainingTime` - Token剩余时间

#### 文件3: `internal/db/db_test.go` (新建) - 0.5h
- [ ] `TestOpen_Postgres` - PostgreSQL连接
- [ ] `TestOpen_SQLite` - SQLite连接
- [ ] `TestOpen_Error` - 连接错误处理

#### 文件4: `internal/db/seed_test.go` (新建) - 0.5h
- [ ] `TestApplySeeds` - 应用种子数据
- [ ] `TestSeedGames` - 种子游戏数据
- [ ] `TestSeedUser` - 种子用户数据
- [ ] `TestSeedPlayer` - 种子陪玩师数据

#### 文件5: `internal/logging/logger_test.go` (增强) - 0.5h
- [ ] `TestLogger_WithContext` - 上下文日志
- [ ] `TestLogger_WithFields` - 字段日志
- [ ] `TestLogger_Error` - 错误日志
- [ ] `TestLogger_Warn` - 警告日志

#### 文件6: `internal/metrics/metrics_test.go` (增强) - 0.5h
- [ ] `TestMetrics_Increment` - 增加计数
- [ ] `TestMetrics_RecordDuration` - 记录时长
- [ ] `TestMetrics_RecordGauge` - 记录仪表

#### 文件7: `internal/config/env_test.go` (增强) - 0.5h
- [ ] `TestConfig_LoadFromFile` - 从文件加载
- [ ] `TestConfig_OverrideFromEnv` - 环境变量覆盖
- [ ] `TestConfig_ValidateCrypto` - 加密配置验证

### 参考文档
- 示例代码: `backend/internal/cache/memory_test.go`
- 测试规范: `.cursor/rules/backend-testing.mdc`

### 验收标准
- [ ] 所有测试通过
- [ ] 每个文件覆盖率 ≥ 50%

---

## Issue 6: [测试] Admin Handler补充 (可选)

**分配给**: @成员A  
**优先级**: P2 (中)  
**预计时间**: 2-3小时  
**预计覆盖率提升**: +2%

### 任务清单

#### 文件1: `internal/handler/admin/stats_handler_test.go` (新建) - 1h
- [ ] `TestStatsHandler_Dashboard` - 仪表盘
- [ ] `TestStatsHandler_RevenueTrend` - 收入趋势
- [ ] `TestStatsHandler_UserGrowth` - 用户增长
- [ ] `TestStatsHandler_OrdersSummary` - 订单摘要
- [ ] `TestStatsHandler_TopPlayers` - 顶级陪玩师
- [ ] `TestStatsHandler_AuditOverview` - 审计概览
- [ ] `TestStatsHandler_AuditTrend` - 审计趋势

#### 文件2: `internal/handler/admin/system_handler_test.go` (新建) - 0.5h
- [ ] `TestSystemInfoHandler_Config` - 配置信息
- [ ] `TestSystemInfoHandler_DBStatus` - 数据库状态
- [ ] `TestSystemInfoHandler_CacheStatus` - 缓存状态
- [ ] `TestSystemInfoHandler_Resources` - 资源信息
- [ ] `TestSystemInfoHandler_Version` - 版本信息

#### 文件3: 其他Admin Handler文件 (可选) - 2h
- [ ] `internal/handler/admin/commission_test.go`
- [ ] `internal/handler/admin/dashboard_test.go`
- [ ] `internal/handler/admin/item_test.go`
- [ ] `internal/handler/admin/ranking_test.go`
- [ ] `internal/handler/admin/stats_test.go`
- [ ] `internal/handler/admin/withdraw_test.go`

### 验收标准
- [ ] 所有测试通过
- [ ] 每个文件覆盖率 ≥ 50%

---

## 📊 总体进度跟踪

### 第1天目标 (35.5% → 50%)
- [ ] Issue 2: Service层测试 (成员B)
- [ ] Issue 5: 小模块测试 (成员C)
- [ ] Issue 1: Admin Handler开始 (成员A)

### 第2天目标 (50% → 63%)
- [ ] Issue 1: Admin Handler完成 (成员A)
- [ ] Issue 3: User/Player Handler (成员A或B)

### 第3天目标 (63% → 70%+)
- [ ] Issue 4: Repository层补充 (成员B)
- [ ] Issue 6: Admin Handler补充 (成员A，可选)
- [ ] 代码审查和优化 (全体)

---

## 📝 使用说明

1. **创建Issue**: 将每个Issue复制到GitHub Issues中
2. **分配成员**: 使用 `@成员A`, `@成员B`, `@成员C` 标签
3. **跟踪进度**: 使用复选框跟踪任务完成情况
4. **每日同步**: 在Issue中更新进度和遇到的问题

---

**文档版本**: 1.0  
**最后更新**: 2025-11-08

