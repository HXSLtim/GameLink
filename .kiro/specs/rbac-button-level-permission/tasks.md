# 实现计划 - RBAC 按钮级权限管理模块

- [x] 1. 扩展权限数据模型





  - [x] 1.1 扩展 Permission 模型


    - 添加 IsSystem 字段标记系统权限（不可删除）
    - 添加 DeletedAt 字段支持软删除
    - 添加 ParentID 和 SortOrder 字段支持树形结构
    - _Requirements: 1.1, 1.5_
  - [x] 1.2 扩展 Role 模型


    - 添加 ParentID 字段支持角色继承
    - 添加 Priority 字段用于权限冲突解决（子角色覆盖父角色）
    - 添加 Level 字段记录继承层级（最大 5 层）
    - 添加 DeletedAt 字段支持软删除
    - _Requirements: 10.2, 10.4_
  - [x] 1.3 创建 PermissionAuditLog 模型


    - 包含 BeforeData/AfterData 字段存储操作前后数据快照
    - 包含 IPAddress/UserAgent/RequestID 字段
    - _Requirements: 6.1, 6.2_
  - [ ]* 1.4 编写模型单元测试
    - 测试软删除功能
    - 测试继承层级计算
    - _Requirements: 1.1, 10.2_

- [-] 2. 实现权限树形结构服务



  - [x] 2.1 实现 GetPermissionTree 方法


    - 使用递归 CTE 查询避免 N+1 问题
    - 按 Group 分组返回树形结构
    - _Requirements: 2.1_
  - [-] 2.2 实现权限码格式验证

    - 验证格式：module.resource.action（三段式）
    - 创建后权限码不可修改
    - _Requirements: 1.2, 1.3_
  - [ ] 2.3 实现权限删除引用检查
    - 检查是否被角色引用
    - 系统权限（IsSystem=true）不可删除
    - _Requirements: 1.5_
  - [ ]* 2.4 编写权限树属性测试
    - **Property 1: 权限码格式验证**
    - **Validates: Requirements 1.3**

- [ ] 3. 实现角色继承功能
  - [ ] 3.1 实现 SetRoleParent 方法
    - 设置角色继承关系
    - 自动计算 Level 字段
    - 限制最大继承深度为 5 层
    - _Requirements: 10.2_
  - [ ] 3.2 实现 GetRoleInheritanceChain 方法
    - 使用递归 CTE 查询继承链
    - 避免 N+1 查询问题
    - _Requirements: 10.2, 10.3_
  - [ ] 3.3 实现 ValidateNoCircularInheritance 方法
    - 检测循环继承
    - 拒绝会导致循环的设置
    - _Requirements: 10.5_
  - [ ] 3.4 实现权限冲突解决策略
    - 子角色权限覆盖父角色（按 Priority 排序）
    - 合并时去重
    - _Requirements: 10.4_
  - [ ]* 3.5 编写角色继承属性测试
    - **Property 19: 循环继承检测**
    - **Validates: Requirements 10.5**

- [ ] 4. Checkpoint - 确保所有测试通过
  - 确保所有测试通过，如有问题请询问用户

- [ ] 5. 实现权限审计日志服务
  - [ ] 5.1 实现异步审计日志写入
    - 使用 channel 缓冲队列
    - 后台 goroutine 批量写入
    - 避免阻塞主流程
    - _Requirements: 6.1, 6.2_
  - [ ] 5.2 实现审计日志查询
    - 支持按时间范围筛选
    - 支持按操作类型筛选
    - 支持按操作者筛选
    - _Requirements: 6.3_
  - [ ] 5.3 实现审计日志导出
    - 导出 CSV 格式
    - 包含完整操作记录
    - _Requirements: 6.5_
  - [ ] 5.4 实现审计日志归档策略
    - 在线保留 90 天
    - 归档保留 365 天
    - _Requirements: 6.4_
  - [ ]* 5.5 编写审计日志属性测试
    - **Property 11: 审计日志完整性**
    - **Validates: Requirements 6.1, 6.2**

- [ ] 6. 增强权限中间件
  - [ ] 6.1 支持按权限码检查
    - 除了 method+path，支持 code 检查
    - _Requirements: 4.1_
  - [ ] 6.2 支持多权限检查模式
    - any 模式：任一权限满足
    - all 模式：全部权限满足
    - except 模式：排除某些权限
    - _Requirements: 4.3_
  - [ ] 6.3 实现超级管理员绕过
    - 自动跳过所有权限检查
    - _Requirements: 4.4_
  - [ ] 6.4 实现权限白名单
    - 配置无需权限检查的接口
    - _Requirements: 4.1_
  - [ ]* 6.5 编写中间件属性测试
    - **Property 7: 超级管理员权限绕过**
    - **Property 8: API 权限验证一致性**
    - **Validates: Requirements 4.1, 4.4**

- [ ] 7. 实现权限缓存管理
  - [ ] 7.1 实现缓存 Key 设计
    - 包含版本号便于快速刷新
    - 使用 Hash 结构存储用户权限
    - _Requirements: 5.1, 5.2_
  - [ ] 7.2 实现缓存 TTL 随机抖动
    - 防止缓存雪崩
    - 默认 30 分钟 + 10% 随机抖动
    - _Requirements: 5.5_
  - [ ] 7.3 实现缓存预热机制
    - 系统启动时预热权限树
    - 预热系统角色权限
    - _Requirements: 5.5_
  - [ ] 7.4 实现缓存失效传播
    - 角色权限变更时失效所有相关用户缓存
    - 用户角色变更时失效用户缓存
    - _Requirements: 5.1, 5.2_
  - [ ]* 7.5 编写缓存属性测试
    - **Property 9: 权限缓存失效传播**
    - **Validates: Requirements 5.1, 5.2**

- [ ] 8. Checkpoint - 确保所有测试通过
  - 确保所有测试通过，如有问题请询问用户

- [ ] 9. 实现权限管理 API
  - [ ] 9.1 实现权限列表接口
    - GET /api/admin/permissions（分页、搜索、筛选）
    - GET /api/admin/permissions/:id（详情）
    - _Requirements: 1.1_
  - [ ] 9.2 实现权限树接口
    - GET /api/admin/permissions/tree
    - GET /api/admin/permissions/groups
    - _Requirements: 2.1_
  - [ ] 9.3 实现权限 CRUD 接口
    - POST /api/admin/permissions（创建）
    - PUT /api/admin/permissions/:id（全量更新）
    - PATCH /api/admin/permissions/:id（部分更新）
    - DELETE /api/admin/permissions/:id（软删除）
    - _Requirements: 1.2, 1.4, 1.5_
  - [ ]* 9.4 更新 Swagger 文档
    - _Requirements: 1.1_

- [ ] 10. 实现角色权限分配 API
  - [ ] 10.1 实现角色权限查询接口
    - GET /api/admin/roles/:id/permissions
    - _Requirements: 2.4_
  - [ ] 10.2 实现批量权限分配接口
    - PUT /api/admin/roles/:id/permissions/batch（事务保证原子性）
    - POST /api/admin/roles/:id/permissions/:pid（单个添加）
    - DELETE /api/admin/roles/:id/permissions/:pid（单个移除）
    - _Requirements: 2.2, 2.3_
  - [ ] 10.3 实现分配后缓存失效
    - 自动失效相关用户缓存
    - 记录审计日志
    - _Requirements: 5.2, 6.1_
  - [ ]* 10.4 编写角色权限分配属性测试
    - **Property 5: 权限分配持久化**
    - **Validates: Requirements 2.3**

- [ ] 11. 实现用户角色分配 API
  - [ ] 11.1 实现用户角色查询接口
    - GET /api/admin/users/:id/roles
    - GET /api/admin/users/:id/permissions（有效权限）
    - _Requirements: 10.3_
  - [ ] 11.2 实现用户角色分配接口
    - PUT /api/admin/users/:id/roles
    - 支持批量用户角色分配
    - _Requirements: 9.1, 9.2_
  - [ ] 11.3 实现分配后缓存失效
    - 自动失效用户缓存
    - 记录审计日志
    - _Requirements: 5.1, 6.2_
  - [ ]* 11.4 编写用户角色分配属性测试
    - **Property 17: 多角色权限合并**
    - **Validates: Requirements 10.1**

- [ ] 12. 实现当前用户权限 API
  - [ ] 12.1 实现当前用户权限接口
    - GET /api/admin/me/permissions
    - 超级管理员返回 ['*']
    - _Requirements: 5.3_
  - [ ] 12.2 实现当前用户菜单接口
    - GET /api/admin/me/menus
    - 根据权限自动过滤菜单
    - _Requirements: 8.1_
  - [ ]* 12.3 编写当前用户权限属性测试
    - **Property 10: 登录权限完整性**
    - **Validates: Requirements 5.3**

- [ ] 13. Checkpoint - 确保所有测试通过
  - 确保所有测试通过，如有问题请询问用户

- [ ] 14. 前端权限 API 集成
  - [ ] 14.1 添加权限管理 API 调用方法
    - 权限 CRUD API
    - 权限树 API
    - _Requirements: 1.1, 2.1_
  - [ ] 14.2 添加角色权限分配 API 调用方法
    - 角色权限查询
    - 批量权限分配
    - _Requirements: 2.2, 2.3_
  - [ ] 14.3 添加用户角色分配 API 调用方法
    - 用户角色查询
    - 用户角色分配
    - _Requirements: 9.1, 9.2_
  - [ ] 14.4 添加审计日志 API 调用方法
    - 审计日志查询
    - 审计日志导出
    - _Requirements: 6.3, 6.5_
  - [ ] 14.5 创建 TypeScript 类型定义
    - Permission 类型
    - Role 类型
    - AuditLog 类型
    - _Requirements: 1.1_

- [ ] 15. 增强 PermissionGuard 组件
  - [ ] 15.1 添加加载状态支持
    - loading 属性显示加载状态
    - 避免权限闪烁
    - _Requirements: 3.3_
  - [ ] 15.2 添加禁用模式支持
    - disabled 属性禁用模式
    - tooltip 提示无权限原因
    - _Requirements: 3.2_
  - [ ] 15.3 添加 fallback 支持
    - 自定义无权限显示内容
    - _Requirements: 3.2_
  - [ ]* 15.4 编写组件测试
    - **Property 6: 前端权限检查一致性**
    - **Validates: Requirements 3.1, 3.4**

- [ ] 16. 增强 usePermission Hook
  - [ ] 16.1 优化性能
    - 使用 useMemo 避免不必要的重渲染
    - 减少依赖项
    - _Requirements: 3.1_
  - [ ] 16.2 添加批量权限检查
    - usePermissions 支持批量检查
    - 返回每个权限的检查结果
    - _Requirements: 3.4_
  - [ ] 16.3 添加动态检查功能
    - usePermissionChecker 返回检查函数
    - 支持运行时动态检查
    - _Requirements: 3.4_
  - [ ]* 16.4 编写 Hook 测试
    - _Requirements: 3.1, 3.4_

- [ ] 17. 实现权限管理页面
  - [ ] 17.1 实现权限列表页面
    - 分页、搜索、筛选
    - 按分组展示
    - _Requirements: 1.1_
  - [ ] 17.2 实现权限创建/编辑表单
    - 权限码格式验证
    - 编辑时权限码不可修改
    - _Requirements: 1.2, 1.4_
  - [ ] 17.3 实现权限删除确认
    - 显示引用警告
    - 系统权限不可删除
    - _Requirements: 1.5_
  - [ ]* 17.4 编写页面测试
    - _Requirements: 1.1_

- [ ] 18. 实现角色权限配置页面
  - [ ] 18.1 实现权限树组件
    - 支持虚拟滚动（性能优化）
    - 支持全选/反选
    - _Requirements: 2.1_
  - [ ] 18.2 实现父子节点联动选择
    - 选中父节点自动选中子节点
    - 取消父节点自动取消子节点
    - _Requirements: 2.2_
  - [ ] 18.3 实现角色权限配置
    - 高亮已分配权限
    - 系统角色显示特殊提示
    - _Requirements: 2.4, 2.5_
  - [ ]* 18.4 编写页面测试
    - **Property 4: 角色权限树形选择一致性**
    - **Validates: Requirements 2.2**

- [ ] 19. 实现用户角色分配页面
  - [ ] 19.1 实现用户列表显示当前角色
    - 显示用户已分配的角色
    - _Requirements: 9.1_
  - [ ] 19.2 实现角色分配弹窗
    - 多选角色
    - 操作预览
    - _Requirements: 9.2_
  - [ ] 19.3 实现批量角色分配
    - 批量选择用户
    - 显示操作结果摘要
    - _Requirements: 9.1, 9.3_
  - [ ] 19.4 实现有效权限查看
    - 显示合并后的完整权限列表
    - 显示权限来源
    - _Requirements: 10.3_
  - [ ]* 19.5 编写页面测试
    - **Property 16: 批量操作结果报告**
    - **Validates: Requirements 9.2, 9.3, 9.4**

- [ ] 20. 实现审计日志页面
  - [ ] 20.1 实现审计日志列表
    - 分页显示
    - 显示操作前后数据对比
    - _Requirements: 6.3_
  - [ ] 20.2 实现筛选功能
    - 按时间范围筛选
    - 按操作类型筛选
    - 按操作者筛选
    - _Requirements: 6.3_
  - [ ] 20.3 实现导出功能
    - 导出 CSV 格式
    - _Requirements: 6.5_
  - [ ]* 20.4 编写页面测试
    - **Property 12: 审计日志过滤正确性**
    - **Validates: Requirements 6.3**

- [ ] 21. Checkpoint - 确保所有测试通过
  - 确保所有测试通过，如有问题请询问用户

- [ ] 22. 实现菜单权限联动
  - [ ] 22.1 实现菜单权限过滤
    - 根据用户权限过滤菜单
    - 无子菜单权限时隐藏父菜单
    - _Requirements: 8.1, 8.2_
  - [ ] 22.2 实现路由权限守卫
    - 无权限重定向到 403 页面
    - _Requirements: 8.3_
  - [ ] 22.3 实现 403 页面
    - 显示无权限提示
    - 提供申请权限入口（可选）
    - _Requirements: 8.3_
  - [ ] 22.4 实现权限变更后菜单更新
    - 监听权限变更
    - 自动刷新菜单
    - _Requirements: 8.4_
  - [ ]* 22.5 编写菜单权限测试
    - **Property 14: 菜单权限过滤**
    - **Property 15: 路由权限保护**
    - **Validates: Requirements 8.1, 8.3**

- [ ] 23. 集成测试 - 权限验证流程
  - [ ] 23.1 测试用户权限查询
    - 使用独立事务隔离测试数据
    - _Requirements: 4.1_
  - [ ] 23.2 测试 API 权限验证
    - 测试有权限和无权限场景
    - _Requirements: 4.1, 4.2_
  - [ ] 23.3 测试超级管理员绕过
    - _Requirements: 4.4_
  - [ ] 23.4 测试权限缓存失效
    - _Requirements: 5.1, 5.2_
  - [ ] 23.5 测试多角色权限合并
    - _Requirements: 10.1_

- [ ] 24. 集成测试 - 审计日志
  - [ ] 24.1 测试权限分配审计日志
    - 验证日志包含完整信息
    - _Requirements: 6.1_
  - [ ] 24.2 测试角色分配审计日志
    - _Requirements: 6.2_
  - [ ] 24.3 测试审计日志查询筛选
    - _Requirements: 6.3_
  - [ ] 24.4 测试审计日志导出
    - _Requirements: 6.5_

- [ ] 25. 种子数据更新
  - [ ] 25.1 定义权限码常量
    - 按模块分组定义
    - 使用代码而非 SQL 管理
    - _Requirements: 1.1_
  - [ ] 25.2 配置默认角色权限
    - superAdmin: 所有权限
    - admin: 管理权限
    - player: 陪玩师权限
    - user: 用户权限
    - _Requirements: 2.3_
  - [ ] 25.3 标记系统权限和角色
    - IsSystem = true
    - _Requirements: 1.5_
  - [ ] 25.4 确保种子数据幂等执行
    - _Requirements: 1.1_

- [ ] 26. 文档更新
  - [ ] 26.1 更新 Swagger API 文档
    - _Requirements: 1.1_
  - [ ] 26.2 编写权限码列表文档
    - _Requirements: 1.1_
  - [ ] 26.3 编写前端权限组件使用指南
    - _Requirements: 3.1_
  - [ ] 26.4 编写角色权限配置指南
    - _Requirements: 2.1_

- [ ] 27. Final Checkpoint - 确保所有测试通过
  - 确保所有测试通过，如有问题请询问用户
