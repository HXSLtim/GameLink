# 权限码列表文档

本文档列出了 GameLink 平台所有的权限码定义，用于 RBAC 按钮级权限控制。

## 权限码格式

权限码采用三段式命名：`{module}.{resource}.{action}`

- **module**: 模块名（admin, review, content, commission, notification, wallet, service_item, dispute, system）
- **resource**: 资源名（permissions, roles, users, orders, games 等）
- **action**: 操作名（read, create, update, delete, assign, export 等）

## 权限分组

| 分组标识 | 分组名称 | 描述 | 模块 |
|---------|---------|------|------|
| /admin/permissions | 权限管理 | 管理系统权限定义 | admin |
| /admin/roles | 角色管理 | 管理系统角色和权限分配 | admin |
| /admin/users | 用户管理 | 管理平台用户 | admin |
| /admin/players | 陪玩师管理 | 管理陪玩师资料和认证 | admin |
| /admin/games | 游戏管理 | 管理游戏列表 | admin |
| /admin/orders | 订单管理 | 管理订单和订单状态 | admin |
| /admin/payments | 支付管理 | 管理支付和退款 | admin |
| /admin/withdraws | 提现管理 | 管理提现申请和审批 | admin |
| /admin/audit | 审计日志 | 查看和导出审计日志 | admin |
| /admin/stats | 统计数据 | 查看平台统计数据 | admin |
| /admin/menus | 菜单管理 | 管理后台菜单 | admin |
| /admin/reviews | 评价管理 | 管理用户评价和审核 | review |
| /admin/content | 内容管理 | 管理动态、聊天和举报 | content |
| /admin/commission | 佣金管理 | 管理佣金结算 | commission |
| /admin/notification | 通知管理 | 管理系统通知 | notification |
| /admin/wallet | 钱包管理 | 管理用户钱包 | wallet |
| /admin/service-items | 服务项目管理 | 管理服务项目 | service_item |
| /admin/disputes | 纠纷管理 | 管理订单纠纷 | dispute |
| /admin/system | 系统管理 | 系统配置和监控 | system |

## 权限码详细列表

### Admin 模块 - 管理后台权限

#### 权限管理 (admin.permissions.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `admin.permissions.read` | 查看权限列表/详情/树/分组 | `PermCodeAdminPermissionsRead` |
| `admin.permissions.create` | 创建权限 | `PermCodeAdminPermissionsCreate` |
| `admin.permissions.update` | 更新权限 | `PermCodeAdminPermissionsUpdate` |
| `admin.permissions.delete` | 删除权限 | `PermCodeAdminPermissionsDelete` |

#### 角色管理 (admin.roles.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `admin.roles.read` | 查看角色列表/详情/权限 | `PermCodeAdminRolesRead` |
| `admin.roles.create` | 创建角色 | `PermCodeAdminRolesCreate` |
| `admin.roles.update` | 更新角色 | `PermCodeAdminRolesUpdate` |
| `admin.roles.delete` | 删除角色 | `PermCodeAdminRolesDelete` |
| `admin.roles.assign` | 分配角色权限 | `PermCodeAdminRolesAssign` |

#### 用户管理 (admin.users.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `admin.users.read` | 查看用户列表/详情/订单/日志/角色/权限 | `PermCodeAdminUsersRead` |
| `admin.users.create` | 创建用户 | `PermCodeAdminUsersCreate` |
| `admin.users.update` | 更新用户 | `PermCodeAdminUsersUpdate` |
| `admin.users.delete` | 删除用户 | `PermCodeAdminUsersDelete` |
| `admin.users.assign` | 分配用户角色 | `PermCodeAdminUsersAssign` |
| `admin.users.status` | 更新用户状态 | `PermCodeAdminUsersStatus` |
| `admin.users.points` | 管理用户积分 | `PermCodeAdminUsersPoints` |

#### 陪玩师管理 (admin.players.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `admin.players.read` | 查看陪玩师列表/详情/日志 | `PermCodeAdminPlayersRead` |
| `admin.players.create` | 创建陪玩师 | `PermCodeAdminPlayersCreate` |
| `admin.players.update` | 更新陪玩师 | `PermCodeAdminPlayersUpdate` |
| `admin.players.delete` | 删除陪玩师 | `PermCodeAdminPlayersDelete` |
| `admin.players.verify` | 审核陪玩师认证 | `PermCodeAdminPlayersVerify` |

#### 游戏管理 (admin.games.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `admin.games.read` | 查看游戏列表/详情/日志 | `PermCodeAdminGamesRead` |
| `admin.games.create` | 创建游戏 | `PermCodeAdminGamesCreate` |
| `admin.games.update` | 更新游戏 | `PermCodeAdminGamesUpdate` |
| `admin.games.delete` | 删除游戏 | `PermCodeAdminGamesDelete` |

#### 订单管理 (admin.orders.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `admin.orders.read` | 查看订单列表/详情/日志/时间线/支付/退款/评价 | `PermCodeAdminOrdersRead` |
| `admin.orders.create` | 创建订单 | `PermCodeAdminOrdersCreate` |
| `admin.orders.update` | 更新订单 | `PermCodeAdminOrdersUpdate` |
| `admin.orders.delete` | 删除订单 | `PermCodeAdminOrdersDelete` |
| `admin.orders.assign` | 指派订单 | `PermCodeAdminOrdersAssign` |
| `admin.orders.cancel` | 取消订单 | `PermCodeAdminOrdersCancel` |
| `admin.orders.refund` | 退款订单 | `PermCodeAdminOrdersRefund` |
| `admin.orders.confirm` | 确认订单 | `PermCodeAdminOrdersConfirm` |
| `admin.orders.start` | 开始订单 | `PermCodeAdminOrdersStart` |
| `admin.orders.complete` | 完成订单 | `PermCodeAdminOrdersComplete` |

#### 支付管理 (admin.payments.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `admin.payments.read` | 查看支付记录 | `PermCodeAdminPaymentsRead` |
| `admin.payments.refund` | 处理退款 | `PermCodeAdminPaymentsRefund` |
| `admin.payments.export` | 导出支付记录 | `PermCodeAdminPaymentsExport` |

#### 提现管理 (admin.withdraws.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `admin.withdraws.read` | 查看提现记录 | `PermCodeAdminWithdrawsRead` |
| `admin.withdraws.approve` | 审批提现 | `PermCodeAdminWithdrawsApprove` |
| `admin.withdraws.reject` | 拒绝提现 | `PermCodeAdminWithdrawsReject` |
| `admin.withdraws.export` | 导出提现记录 | `PermCodeAdminWithdrawsExport` |

#### 审计日志 (admin.audit.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `admin.audit.read` | 查看审计日志 | `PermCodeAdminAuditRead` |
| `admin.audit.export` | 导出审计日志 | `PermCodeAdminAuditExport` |

#### 统计数据 (admin.stats.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `admin.stats.read` | 查看统计数据 | `PermCodeAdminStatsRead` |
| `admin.stats.export` | 导出统计数据 | `PermCodeAdminStatsExport` |

#### 菜单管理 (admin.menus.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `admin.menus.read` | 查看菜单列表 | `PermCodeAdminMenusRead` |
| `admin.menus.create` | 创建菜单 | `PermCodeAdminMenusCreate` |
| `admin.menus.update` | 更新菜单 | `PermCodeAdminMenusUpdate` |
| `admin.menus.delete` | 删除菜单 | `PermCodeAdminMenusDelete` |

### Review 模块 - 评价管理权限

#### 评价管理 (review.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `review.list` | 查看评价列表 | `PermCodeReviewList` |
| `review.get` | 查看评价详情 | `PermCodeReviewGet` |
| `review.pending` | 查看待审核评价 | `PermCodeReviewPending` |
| `review.logs` | 查看评价操作日志 | `PermCodeReviewLogs` |
| `review.player` | 查看陪玩师评价 | `PermCodeReviewPlayer` |
| `review.order` | 查看订单评价 | `PermCodeReviewOrder` |
| `review.approve` | 批准评价 | `PermCodeReviewApprove` |
| `review.reject` | 拒绝评价 | `PermCodeReviewReject` |
| `review.batch_approve` | 批量批准评价 | `PermCodeReviewBatchApprove` |
| `review.batch_reject` | 批量拒绝评价 | `PermCodeReviewBatchReject` |
| `review.delete` | 删除评价 | `PermCodeReviewDelete` |
| `review.update` | 更新评价 | `PermCodeReviewUpdate` |
| `review.create` | 创建评价 | `PermCodeReviewCreate` |
| `review.stats` | 查看评价统计 | `PermCodeReviewStats` |
| `review.trend` | 查看评价趋势 | `PermCodeReviewTrend` |
| `review.top_players` | 查看陪玩师排行榜 | `PermCodeReviewTopPlayers` |
| `review.game_stats` | 查看游戏评价统计 | `PermCodeReviewGameStats` |
| `review.export` | 导出评价统计 | `PermCodeReviewExport` |
| `review.detect_sensitive` | 检测敏感词 | `PermCodeReviewDetectSensitive` |

#### 评价举报管理 (review_report.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `review_report.list` | 查看评价举报列表 | `PermCodeReviewReportList` |
| `review_report.get` | 查看举报详情 | `PermCodeReviewReportGet` |
| `review_report.create` | 创建评价举报 | `PermCodeReviewReportCreate` |
| `review_report.handle` | 处理评价举报 | `PermCodeReviewReportHandle` |

#### 评价回复管理 (review_reply.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `review_reply.update` | 更新评价回复 | `PermCodeReviewReplyUpdate` |
| `review_reply.delete` | 删除评价回复 | `PermCodeReviewReplyDelete` |

#### 评价展示设置 (review_settings.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `review_settings.get` | 查看评价展示设置 | `PermCodeReviewSettingsGet` |
| `review_settings.update` | 更新评价展示设置 | `PermCodeReviewSettingsUpdate` |

#### 敏感词管理 (sensitive_word.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `sensitive_word.list` | 查看敏感词列表 | `PermCodeSensitiveWordList` |
| `sensitive_word.create` | 添加敏感词 | `PermCodeSensitiveWordCreate` |
| `sensitive_word.update` | 更新敏感词 | `PermCodeSensitiveWordUpdate` |
| `sensitive_word.delete` | 删除敏感词 | `PermCodeSensitiveWordDelete` |

#### 操作日志 (operation_log.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `operation_log.list` | 查看操作日志 | `PermCodeOperationLogList` |
| `operation_log.export` | 导出操作日志 | `PermCodeOperationLogExport` |

### Content 模块 - 内容管理权限

#### 动态管理 (content.feed.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `content.feed.list` | 查看动态列表 | `PermCodeContentFeedList` |
| `content.feed.get` | 查看动态详情 | `PermCodeContentFeedGet` |
| `content.feed.approve` | 批准动态 | `PermCodeContentFeedApprove` |
| `content.feed.reject` | 拒绝动态 | `PermCodeContentFeedReject` |
| `content.feed.batch_approve` | 批量批准动态 | `PermCodeContentFeedBatchApprove` |
| `content.feed.batch_reject` | 批量拒绝动态 | `PermCodeContentFeedBatchReject` |
| `content.feed.delete` | 删除动态 | `PermCodeContentFeedDelete` |

#### 聊天监控 (content.chat.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `content.chat.list` | 查看聊天消息 | `PermCodeContentChatList` |
| `content.chat.delete` | 删除聊天消息 | `PermCodeContentChatDelete` |
| `content.chat.mute` | 禁言用户 | `PermCodeContentChatMute` |
| `content.chat.unmute` | 解除禁言 | `PermCodeContentChatUnmute` |

#### 举报管理 (content.report.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `content.report.list` | 查看举报列表 | `PermCodeContentReportList` |
| `content.report.get` | 查看举报详情 | `PermCodeContentReportGet` |
| `content.report.process` | 处理举报 | `PermCodeContentReportProcess` |

#### 内容统计 (content.stats)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `content.stats` | 查看内容统计 | `PermCodeContentStats` |

#### 内容分类管理 (content.category.*)

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `content.category.list` | 查看内容分类 | `PermCodeContentCategoryList` |
| `content.category.get` | 查看分类详情 | `PermCodeContentCategoryGet` |
| `content.category.create` | 创建内容分类 | `PermCodeContentCategoryCreate` |
| `content.category.update` | 更新内容分类 | `PermCodeContentCategoryUpdate` |
| `content.category.delete` | 删除内容分类 | `PermCodeContentCategoryDelete` |

### Commission 模块 - 佣金管理权限

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `commission.read` | 查看佣金记录 | `PermCodeCommissionRead` |
| `commission.settle` | 结算佣金 | `PermCodeCommissionSettle` |
| `commission.export` | 导出佣金记录 | `PermCodeCommissionExport` |

### Notification 模块 - 通知管理权限

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `notification.read` | 查看通知 | `PermCodeNotificationRead` |
| `notification.create` | 创建通知 | `PermCodeNotificationCreate` |
| `notification.batch_send` | 批量发送通知 | `PermCodeNotificationBatchSend` |

### Wallet 模块 - 钱包管理权限

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `wallet.read` | 查看钱包 | `PermCodeWalletRead` |
| `wallet.update` | 更新钱包 | `PermCodeWalletUpdate` |
| `wallet.transactions` | 查看交易记录 | `PermCodeWalletTransactions` |

### Service Item 模块 - 服务项目管理权限

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `service_item.read` | 查看服务项目 | `PermCodeServiceItemRead` |
| `service_item.create` | 创建服务项目 | `PermCodeServiceItemCreate` |
| `service_item.update` | 更新服务项目 | `PermCodeServiceItemUpdate` |
| `service_item.delete` | 删除服务项目 | `PermCodeServiceItemDelete` |

### Dispute 模块 - 纠纷管理权限

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `dispute.read` | 查看纠纷 | `PermCodeDisputeRead` |
| `dispute.create` | 创建纠纷 | `PermCodeDisputeCreate` |
| `dispute.resolve` | 解决纠纷 | `PermCodeDisputeResolve` |

### System 模块 - 系统管理权限

| 权限码 | 描述 | 常量名 |
|-------|------|--------|
| `system.config.read` | 查看系统配置 | `PermCodeSystemConfigRead` |
| `system.config.update` | 更新系统配置 | `PermCodeSystemConfigUpdate` |
| `system.monitor.read` | 查看系统监控 | `PermCodeSystemMonitorRead` |

## 使用示例

### 后端使用

```go
import "gamelink/internal/model"

// 检查用户是否有权限
hasPermission := permissionService.CheckUserHasPermissionCode(ctx, userID, model.PermCodeAdminUsersRead)

// 在中间件中使用
router.GET("/users", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users"), handler.ListUsers)
```

### 前端使用

```tsx
import { PermissionGuard } from '@/components/PermissionGuard';

// 使用 PermissionGuard 组件
<PermissionGuard permission="admin.users.create">
  <Button type="primary">创建用户</Button>
</PermissionGuard>

// 使用 usePermission Hook
const { hasPermission } = usePermission('admin.users.delete');
```

## 默认角色权限配置

| 角色 | 权限范围 |
|------|---------|
| superAdmin | 所有权限（`*`） |
| admin | 管理权限（不含系统配置） |
| player | 陪玩师相关权限 |
| user | 用户基础权限 |

## 注意事项

1. 权限码创建后不可修改
2. 系统权限（IsSystem=true）不可删除
3. 超级管理员自动拥有所有权限，无需逐一检查
4. 权限变更后会自动失效相关用户的缓存
