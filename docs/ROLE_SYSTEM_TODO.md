# GameLink 角色体系 TODO（2026-02-12）

## 目标
- 在不破坏现有 `customerService` 兼容性的前提下，新增 `csLeader` / `csAgent` 两级客服角色。
- 保持移动端与后台的角色类型对齐，避免联调时出现未识别角色。

## P0（本次已完成）
- [x] 后端新增角色标识：`csLeader`、`csAgent`
- [x] 种子角色新增：客服主管、客服专员
- [x] 默认权限拆分：`csLeader`（主管全量）、`csAgent`（日常处理子集）
- [x] 保留旧角色：`customerService`（兼容存量账号/权限）
- [x] 继承关系落地：`csLeader` 继承 `csAgent`（主管在专员基础上扩展）
- [x] `seedVersion` 递增到 `2026-02-12-v11`
- [x] autoMigrate 默认角色补齐：finance/customerService/csLeader/csAgent
- [x] 小程序端角色类型补齐：`UserRole` 新增服务端角色联合类型

## P1（建议下一步）
- [ ] 管理后台角色管理页增加“客服主管/客服专员”筛选和中文标签统一
- [ ] 新增一组演示账号并绑定 `csLeader`、`csAgent`（用于联调验收）
- [ ] 输出角色权限矩阵文档（按页面按钮级权限）

## P2（迁移收尾）
- [ ] 评估并制定 `customerService` 退役窗口（仅当存量账号完成迁移后）
- [ ] 批量迁移脚本：`customerService -> csLeader/csAgent`
- [ ] 增加 RBAC 回归用例：继承权限、冲突优先级、菜单可见性

## 验证清单
- [ ] 运行后端迁移并确认 `roles` 表含 8 个系统角色
- [ ] 运行 seed 后确认 `role_permissions` 与 `user_roles` 无缺失
- [ ] 后台账号分别验证：`admin`、`finance`、`csLeader`、`csAgent`
- [ ] 小程序登录后，`userStore.userInfo.role` 不出现类型报错
