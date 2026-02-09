# GameLink Uniapp PRD 与页面规划（用户端/陪玩师端）

**版本**: v1.1  
**更新时间**: 2026-01-29  
**适用范围**: Uniapp（H5 + 小程序），用户端与陪玩师端

---

## 1. 背景与目标

GameLink 前端已统一迁移为 **Uniapp（Vue 3 + TypeScript + Vite）**，用于跨端覆盖 **H5 与小程序**。  
目标是在单一代码体系下完成用户端与陪玩师端核心业务闭环，并保证跨端体验一致。

**核心目标**:
- 用户端完成 “浏览陪玩师 → 下单 → 支付 → 服务完成 → 评价” 闭环
- 陪玩师端完成 “接单 → 服务 → 完成 → 收益管理” 闭环
- H5 与小程序体验一致，兼容登录/支付/通知差异

**不在本 PRD 范围**: 管理端（独立 Web 项目）

---

## 2. 用户角色

| 角色 | 说明 | 主要目标 |
| --- | --- | --- |
| 普通用户 | 寻找陪玩服务 | 浏览、下单、支付、评价 |
| 陪玩师 | 提供陪玩服务 | 接单、履约、收益管理 |

---

## 3. 核心业务流程

**用户侧流程**  
浏览 → 详情 → 下单 → 支付 → 服务进行中 → 完成确认 → 评价

**陪玩师侧流程**  
认证 → 接单 → 服务 → 完成 → 收益/统计

---

## 4. 功能范围（按模块）

### 4.1 用户端
- **首页/发现**: 游戏推荐、热门陪玩师、搜索与筛选
- **陪玩师列表**: 多维筛选（游戏/价格/评分/在线）
- **陪玩师详情**: 详情、评价、收藏、下单入口
- **下单/订单**: 创建订单、订单列表、订单详情、取消/完成、评价
- **支付/钱包**: 支付流程、充值、钱包明细
- **消息与聊天**: 消息列表、聊天窗口、通知中心、公共频道/聊天室
- **个人中心**: 资料、设置、头像上传、账号安全

### 4.2 陪玩师端
- **工作台**: 今日数据、快捷入口
- **订单管理**: 我的订单、接单、完成
- **收益管理**: 收益概览、提现入口
- **服务管理**: 服务配置（游戏+段位）
- **陪玩认证**: 认证状态、提交资料

### 4.3 通用能力
- 登录/注册（含小程序登录）
- WebSocket 实时消息
- 图片上传与预览
- 错误提示与网络异常处理

---

## 5. 非功能需求

**兼容性**
- H5 + 微信小程序（优先）
- API 域名白名单与小程序安全域配置

**性能**
- 列表分页加载
- 图片懒加载与缓存
- 首屏加载 < 2s（目标）

**安全**
- 登录态管理（Token / 小程序登录态）
- 重要操作防重
- 敏感信息不落地

---

## 6. API 依赖（对接层面）

以下为**必须联调**的 API 模块（名称以业务含义为准）：

**用户端**  
- 陪玩师列表/详情/评价/收藏  
- 订单创建/详情/取消/完成/评价  
- 支付/钱包/充值  
- 用户资料/头像上传/设置  
- 聊天/消息/通知（含公共频道列表/加入/离开/已读）  

**陪玩师端**  
- 订单管理（接单/完成）  
- 认证提交/状态  
- 服务管理  
- 统计与收益  

**跨端特殊说明**  
- 小程序登录：`uni.login` → 后端换取 openid/token  
- 小程序支付：需要单独适配微信支付参数  
- H5 支付：需支持 web 支付或跳转  
- WebSocket：H5 与小程序需兼容鉴权方式

---

## 7. 页面规划（以 Uniapp `pages.json` 为准）

### 7.1 TabBar 页面（使用 CustomTabBar）

| 页面路径 | 名称 | 角色 | 核心功能 | API 依赖 | 状态 |
| --- | --- | --- | --- | --- | --- |
| `pages/index/index` | 首页 | 用户 | 热门游戏/陪玩师、搜索入口 | getHotGames 等 | 已存在 |
| `pages/player/list/index` | 陪玩师列表 | 用户 | 筛选、排序 | getPlayerList | 已对接 |
| `pages/message/list/index` | 消息列表 | 用户 | 通知/聊天列表 | 通知/聊天接口 | 已存在 |
| `pages/profile/index/index` | 我的 | 用户 | 资料/钱包入口 | getUserProfile/getWalletInfo | 已对接 |

### 7.2 用户侧功能页面

| 页面路径 | 名称 | 核心功能 | API 依赖 | 状态 |
| --- | --- | --- | --- | --- |
| `pages/player/detail/index` | 陪玩师详情 | 详情/评价/收藏/下单入口 | getPlayerDetail/getPlayerReviews/addFavorite/removeFavorite | 已对接 |
| `pages/order/create/index` | 下单 | 创建订单 | createOrder/getPlayerDetail | 已对接 |
| `pages/order/list/index` | 我的订单 | 列表筛选 | getOrderList | 已存在 |
| `pages/order/detail/index` | 订单详情 | 取消/完成/支付/评价 | getOrderDetail/cancelOrder/completeOrder/payOrder/submitReview | 已对接 |
| `pages/wallet/index/index` | 钱包 | 余额、明细入口 | getWalletInfo | 已对接 |
| `pages/wallet/recharge/index` | 充值 | 充值/支付 | recharge | 已对接 |
| `pages/profile/edit/index` | 编辑资料 | 修改资料/上传头像 | getUserProfile/updateUserProfile/uploadAvatar | 已对接 |
| `pages/settings/index/index` | 设置 | 主题/通知/隐私 | userSettings/notificationSettings | 已存在 |
| `pages/review/list/index` | 我的评价 | 评价列表 | getMyReviews | 已存在 |
| `pages/game/list/index` | 游戏列表 | 游戏筛选 | getGameList | 已存在 |
| `pages/favorite/list/index` | 收藏 | 收藏列表 | getFavorites | 已存在 |
| `pages/message/chat/index` | 聊天页 | 发送/接收消息 | WebSocket/聊天 API | 已存在 |

### 7.3 认证与账号

| 页面路径 | 名称 | 核心功能 | API 依赖 | 状态 |
| --- | --- | --- | --- | --- |
| `pages/auth/login/index` | 登录 | 手机号/小程序登录 | auth/login/wechat | 已存在 |
| `pages/auth/register/index` | 注册 | 注册入口 | auth/register | 已存在 |

### 7.4 陪玩师端

| 页面路径 | 名称 | 核心功能 | API 依赖 | 状态 |
| --- | --- | --- | --- | --- |
| `pages/player/dashboard/index` | 工作台 | 今日数据、快捷入口 | player/stats | 已存在 |
| `pages/player/orders/index` | 陪玩师订单 | 接单/完成 | getPlayerOrders/acceptPlayerOrder/completePlayerOrder | 已对接 |
| `pages/player/earnings/index` | 我的收益 | 收益概览 | getEarningsSummary | 已存在 |
| `pages/player/services/index` | 服务管理 | 服务增删改 | player/services API | 已存在 |
| `pages/player/certification/index` | 陪玩认证 | 认证提交 | getCertificationStatus/submitCertification/uploadIdCardImage | 已对接 |

### 7.5 补充页面（已完成）

| 页面路径 | 名称 | 说明 | 状态 |
| --- | --- | --- | --- |
| `pages/channel/list/index` | 公共频道列表 | Discord 风格频道入口 | ✅ 已完成 |
| `pages/payment/result/index` | 支付结果页 | 支付成功/失败/处理中 | ✅ 已完成 |
| `pages/agreement/index` | 协议页 | 用户协议/隐私政策/陪玩师协议/充值协议 | ✅ 已完成 |
| `pages/help/index` | 帮助中心 | FAQ + 分类搜索 | ✅ 已完成 |
| `pages/service/index` | 在线客服 | 快捷问题 + 实时聊天 | ✅ 已完成 |
| - | 频道聊天 | 复用 `pages/message/chat/index` | ✅ 已支持 |

### 7.6 聊天系统架构（Discord 风格）

#### 设计理念

采用类似 **Discord 频道** 的架构，支持多种群组类型，**所有聊天内容均可被后台监控**，用于合规审查、防止线下交易等。

#### 群组类型

| 类型 | groupType | 说明 | 监控 |
|------|-----------|------|------|
| 公共频道 | `public` | 游戏主题频道，所有用户可加入 | ✅ 后台可查 |
| 订单群 | `order` | 订单创建后自动生成，用户+陪玩师 | ✅ 后台可查 |
| 私人群 | `private` | 用户发起私聊，支持多人 | ✅ 后台可查 |

#### 核心特性

1. **允许私聊但可监控**
   - 用户可自由创建私聊群组
   - 管理员/客服可在后台查看所有聊天记录
   - 用于审查违规内容、防止线下交易

2. **频道层级**
   - 公共频道按游戏分类（如：王者荣耀频道、LOL频道）
   - 订单群自动创建/归档
   - 私人群用户自主管理

3. **消息类型支持**
   - 文本消息
   - 图片消息（支持审核）
   - 系统消息（订单状态变更等）

#### API 接口

| 功能 | 接口 | 说明 |
|------|------|------|
| 公共频道列表 | `GET /public/chat/public-channels` | 获取所有公共频道 |
| 创建群组 | `POST /user/chat/groups` | 创建私聊/订单群 |
| 群组详情 | `GET /user/chat/groups/:id` | 获取群组信息 |
| 加入频道 | `POST /user/chat/groups/:id/join` | 加入公共频道 |
| 离开频道 | `POST /user/chat/groups/:id/leave` | 离开频道 |
| 标记已读 | `POST /user/chat/groups/:id/read` | 传 `messageId` |

#### 后台监控能力

- 管理端可查看所有群组的聊天记录
- 敏感词自动检测与告警
- 可对违规用户进行禁言/封禁
- 聊天记录保留策略（按需配置）

#### 页面规划

| 页面 | 说明 | 优先级 |
|------|------|--------|
| 公共频道列表 | 展示所有公共频道，支持加入/离开 | P0 |
| 聊天页 | 复用 `pages/message/chat/index`，根据 groupType 展示不同 UI | 已完成 |

#### 聊天页复用策略

`pages/message/chat/index` 支持三种模式：
- **公共频道模式**：显示成员列表、频道信息
- **订单群模式**：显示订单状态、快捷操作
- **私聊模式**：标准聊天界面

---

## 8. 页面字段清单（核心页面）

> 说明：字段为前端展示所需的最低集合，具体字段以 API 返回为准。

### 8.1 首页 `pages/index/index`
- 热门游戏：`id`, `name`, `iconUrl`, `category`, `playerCount`
- 推荐陪玩师：`id`, `nickname`, `avatarUrl`, `rank`, `ratingAverage`, `hourlyRateCents`, `isOnline`
- 搜索/筛选：`keyword`, `gameId`, `minPrice`, `maxPrice`, `rating`, `onlineOnly`

### 8.2 陪玩师列表 `pages/player/list/index`
- 陪玩师卡片：`id`, `nickname`, `avatarUrl`, `bio`, `rank`, `ratingAverage`, `ratingCount`, `hourlyRateCents`, `mainGame`, `isOnline`, `orderCount`
- 筛选条件：`gameId`, `minPrice`, `maxPrice`, `minRating`, `onlineOnly`, `sortBy`, `page`, `pageSize`

### 8.3 陪玩师详情 `pages/player/detail/index`
- 基础信息：`id`, `nickname`, `avatarUrl`, `bio`, `rank`, `ratingAverage`, `ratingCount`, `hourlyRateCents`, `isOnline`
- 服务/标签：`tags`, `mainGame`, `rankLevel`
- 统计信息：`orderCount`, `goodRatio`, `avgResponseMin`
- 评价列表：`reviewId`, `rating`, `comment`, `userNickname`, `userAvatarUrl`, `createdAt`
- 收藏状态：`isFavorite`

### 8.4 下单页 `pages/order/create/index`
- 陪玩师信息：`playerId`, `nickname`, `avatarUrl`, `hourlyRateCents`
- 游戏/段位：`gameId`, `gameName`, `rankId`, `rankName`
- 预约信息：`scheduledStart`, `durationHours`
- 价格信息：`unitPriceCents`, `totalPriceCents`
- 备注：`userNotes`

### 8.5 订单列表 `pages/order/list/index`
- 订单卡片：`orderId`, `orderNo`, `title`, `status`, `totalPriceCents`, `createdAt`, `scheduledStart`
- 陪玩师信息：`playerId`, `playerNickname`, `playerAvatar`
- 游戏信息：`gameName`
- 操作状态：`canPay`, `canCancel`, `canComplete`, `canReview`

### 8.6 订单详情 `pages/order/detail/index`
- 订单信息：`orderNo`, `status`, `title`, `description`, `scheduledStart`, `scheduledEnd`, `startedAt`, `completedAt`, `createdAt`
- 金额信息：`unitPriceCents`, `totalPriceCents`, `refundAmountCents`
- 陪玩师信息：`playerId`, `playerNickname`, `playerAvatar`, `rank`
- 支付信息：`paymentId`, `status`, `method`, `paidAt`
- 时间线：`time`, `status`, `message`
- 评价信息：`rating`, `comment`, `createdAt`

### 8.7 消息列表 `pages/message/list/index`
- 会话列表：`groupId`, `groupName`, `groupType`（public/order/private）, `avatarUrl`, `lastMessage`, `lastMessageAt`, `unreadCount`, `isMuted`
- 通知列表：`id`, `title`, `message`, `priority`, `createdAt`, `isRead`
- Tab 切换：聊天 / 通知

### 8.8 公共频道列表（规划页）
- 频道信息：`groupId`, `groupName`, `description`, `avatarUrl`, `currentMembers`, `maxMembers`, `isActive`, `gameId`
- 状态标识：`isJoined`
- 分类：按游戏分组展示

### 8.9 聊天页 `pages/message/chat/index`
- 群组信息：`groupId`, `groupName`, `groupType`（public/order/private）, `memberCount`, `orderId`（订单群专用）
- 成员列表：`userId`, `nickname`, `avatarUrl`, `role`（owner/admin/member）
- 消息列表：`messageId`, `senderId`, `senderName`, `senderAvatar`, `messageType`（text/image/system）, `content`, `imageUrl`, `createdAt`, `replyToId`
- 输入框：`content`, `imageUploadUrl`
- 已读状态：`lastReadMessageId`
- 模式适配：
  - 公共频道：显示成员数、频道描述
  - 订单群：显示订单状态、快捷操作按钮
  - 私聊：标准聊天 UI

### 8.10 个人中心 `pages/profile/index/index`
- 用户信息：`userId`, `nickname/name`, `avatarUrl`, `phone`, `role`
- 钱包入口：`balance`
- 入口卡片：订单/收藏/设置/认证

### 8.11 编辑资料 `pages/profile/edit/index`
- 可编辑字段：`nickname/name`, `avatarUrl`
- 更新结果：`updatedAt`

### 8.12 钱包与充值 `pages/wallet/index/index`, `pages/wallet/recharge/index`
- 钱包：`balance`, `totalSpent`, `totalRecharge`
- 充值：`amountCents`, `method`, `status`

### 8.13 陪玩师工作台 `pages/player/dashboard/index`
- 今日数据：`orderCount`, `earningsCents`, `ratingAverage`
- 快捷入口：订单/收益/服务/认证

### 8.14 陪玩师订单 `pages/player/orders/index`
- 订单卡片：`orderId`, `title`, `status`, `priceCents`, `scheduledStart`, `durationHours`
- 用户信息：`userNickname`, `userAvatar`
- 操作按钮：接单/完成

### 8.15 服务管理 `pages/player/services/index`
- 服务项：`serviceId`, `gameId`, `gameName`, `rankId`, `rankName`, `description`, `isActive`, `updatedAt`

### 8.16 陪玩认证 `pages/player/certification/index`
- 认证状态：`status`, `submittedAt`, `verifiedAt`, `rejectReason`
- 证件信息：`realName`, `idCardFrontUrl`, `idCardBackUrl`

---

## 9. 优先级规划（执行建议）

**P0** ✅ 已完成
- 核心业务闭环联调
- 订单/支付/聊天稳定性
- 公共频道列表页

**P1** ✅ 已完成
- 协议页（用户协议/隐私政策/陪玩师协议/充值协议）
- 帮助中心（FAQ + 分类）
- 客服入口（快捷问题 + 聊天）
- 支付结果页

**P2** ✅ 已完成
- 静态资源：默认头像、空状态 SVG（订单/消息/收藏）
- CustomTabBar 图标：占位文件 + 设计规范文档
- 错误边界组件：`ErrorBoundary`（网络/服务器/权限/404 等）

**P3** ✅ 已完成
- 性能优化：`LazyImage` 图片懒加载、`VirtualList` 虚拟列表
- 小程序发布：`manifest.json` 完整配置、`PrivacyPopup` 隐私授权弹窗

**P4** ✅ 2026-01-29 PRD 符合性修复
- 首页：添加热门游戏 + 推荐陪玩师列表（无论登录与否均可浏览）
- 下单页：添加预约时间（`scheduledStart`）和时长（`durationHours`）字段
- 订单详情：添加订单时间线（`time`, `status`, `message`）
- 消息列表：Tab 修改为「聊天 / 通知」（移除订单 Tab）
- 聊天页：添加 WebSocket 连接状态指示（连接中/已连接/断开）
- 钱包：添加 `totalSpent` / `totalRecharge` 字段展示
- 设置页：隐私字段改为 PRD 规定的 `showOnlineStatus` / `allowStrangerMessage`
- 全局：移除 emoji 图标，统一使用 `uv-icon` 组件

---

## 10. 建议同步事项

- 前端与后端统一 **API 名称映射表**  
- 小程序登录/支付差异化流程明确  
- WebSocket 鉴权方式统一（Header/Query/Cookie）  
- 关键页面联调清单（按上述页面表推进）

---

如需，我可以将此文档加入索引 (`docs/frontend/INDEX.md`) 并同步到 README。 
