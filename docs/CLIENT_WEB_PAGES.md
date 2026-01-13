# GameLink Client 页面功能文档

> **PWA Client 端** - 完整页面功能规格

---

## 文档概述

```
┌─────────────────────────────────────────────────────────────┐
│              GAMELINK CLIENT PWA - 页面功能文档              │
├─────────────────────────────────────────────────────────────┤
│  项目类型: 游戏陪玩平台 PWA                                 │
│  用户角色: Guest (游客) | User (用户) | Player (陪玩师)     │
│  后端状态: ✅ 100% 完成 (36/36 模块)                        │
│  设计系统: Kook Day + Discord Night 双主题                  │
└─────────────────────────────────────────────────────────────┘
```

---

## 目录

1. [认证页面模块](#1-认证页面模块)
2. [首页与发现](#2-首页与发现)
3. [陪玩师相关](#3-陪玩师相关)
4. [订单系统](#4-订单系统)
5. [聊天通讯](#5-聊天通讯)
6. [支付与钱包](#6-支付与钱包)
7. [个人中心](#7-个人中心)
8. [营销功能](#8-营销功能)
9. [陪玩师专用](#9-陪玩师专用)
10. [设置与辅助](#10-设置与辅助)

---

## 1. 认证页面模块

### 1.1 登录页 `/login`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/login` |
| **访问角色** | Guest |
| **优先级** | P0 (核心) |

**功能点:**
- 手机号/邮箱登录
- 验证码登录
- 密码登录
- 记住登录状态
- 第三方登录预留入口

**API 依赖:**
```
POST /api/v1/auth/login
POST /api/v1/auth/send-code
```

**组件:**
- `LoginForm` - 登录表单组件
- `CodeInput` - 验证码输入框
- `SocialLogin` - 第三方登录按钮组

---

### 1.2 注册页 `/register`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/register` |
| **访问角色** | Guest |
| **优先级** | P0 (核心) |

**功能点:**
- 手机号注册
- 验证码验证
- 设置密码
- 用户协议勾选
- 注册后自动登录

**API 依赖:**
```
POST /api/v1/auth/register
POST /api/v1/auth/send-code
```

**组件:**
- `RegisterForm` - 注册表单
- `PasswordStrength` - 密码强度指示器
- `AgreementCheckbox` - 协议勾选框

---

### 1.3 忘记密码 `/forgot-password`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/forgot-password` |
| **访问角色** | Guest |
| **优先级** | P1 |

**功能点:**
- 手机号验证
- 新密码设置
- 验证码验证

**API 依赖:**
```
POST /api/v1/auth/password/reset/send
POST /api/v1/auth/password/reset/confirm
```

**组件:**
- `ResetPasswordForm` - 重置密码表单

---

## 2. 首页与发现

### 2.1 首页 `/`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/` |
| **访问角色** | Guest / User |
| **优先级** | P0 (核心) |

**功能点:**
- 轮播 Banner (活动/推荐)
- 热门陪玩师推荐 (在线优先)
- 游戏分类筛选
- 搜索入口
- 快速下单入口

**API 依赖:**
```
GET /api/v1/players (onlineOnly=true, sortBy=rating)
GET /api/v1/games
GET /api/v1/activities (active)
```

**组件:**
- `HeroBanner` - 轮播 Banner
- `PlayerCard` - 陪玩师卡片
- `GameFilter` - 游戏分类筛选器
- `QuickSearch` - 快速搜索栏

---

### 2.2 陪玩师列表 `/players`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/players` |
| **访问角色** | Guest / User |
| **优先级** | P0 (核心) |

**功能点:**
- 陪玩师列表展示
- 多维度筛选 (游戏/价格/评分/排序)
- 分页加载
- 无限滚动
- 收藏快捷操作

**API 依赖:**
```
GET /api/v1/players
  ?gameId={id}
  &minPrice={cents}
  &maxPrice={cents}
  &minRating={float}
  &onlineOnly={bool}
  &sortBy={price|rating|orders}
  &page={int}
  &pageSize={int}
```

**组件:**
- `PlayerCard` - 陪玩师卡片
- `FilterSidebar` - 筛选侧边栏
- `SortSelector` - 排序选择器
- `InfiniteScroll` - 无限滚动容器

---

### 2.3 陪玩师详情 `/player/:id`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/player/:id` |
| **访问角色** | Guest / User |
| **优先级** | P0 (核心) |

**功能点:**
- 陪玩师完整信息展示
- 服务项目列表
- 历史评价展示
- 在线状态显示
- 立即预约/收藏按钮
- 相似陪玩师推荐

**API 依赖:**
```
GET /api/v1/players/{id}
GET /api/v1/players/{id}/reviews
GET /api/v1/players (similar,推荐)
GET /api/v1/favorites/players/{id}/check (收藏状态)
```

**组件:**
- `PlayerProfileHeader` - 陪玩师头部信息
- `ServiceItemList` - 服务项目列表
- `ReviewList` - 评价列表
- `BookingButton` - 预约按钮组
- `SimilarPlayers` - 相似推荐

---

## 3. 陪玩师相关

### 3.1 我的收藏 `/favorites`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/favorites` |
| **访问角色** | User |
| **优先级** | P1 |

**功能点:**
- 收藏的陪玩师列表
- 快速预约
- 取消收藏
- 空状态提示

**API 依赖:**
```
GET /api/v1/favorites/players
DELETE /api/v1/favorites/players/{id}
GET /api/v1/favorites/players/{id}/check
```

**组件:**
- `FavoritePlayerCard` - 收藏卡片
- `EmptyState` - 空状态组件

---

## 4. 订单系统

### 4.1 创建订单 `/order/create`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/order/create` |
| **访问角色** | User |
| **优先级** | P0 (核心) |

**功能点:**
- 选择陪玩师
- 选择游戏/服务
- 选择时间段
- 确认价格
- 使用优惠券
- VIP 折扣显示
- 提交订单

**API 依赖:**
```
GET /api/v1/players/{id}
GET /api/v1/items (服务项目)
GET /api/v1/user/coupons (可用优惠券)
GET /api/v1/user/vip/status
POST /api/v1/orders
```

**组件:**
- `OrderWizard` - 订单创建向导
- `PlayerSelector` - 陪玩师选择器
- `ServiceSelector` - 服务选择器
- `TimeSlotPicker` - 时间段选择器
- `PriceSummary` - 价格汇总
- `CouponSelector` - 优惠券选择器

---

### 4.2 订单详情 `/order/:id`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/order/:id` |
| **访问角色** | User |
| **优先级** | P0 (核心) |

**功能点:**
- 订单完整信息
- 实时状态更新
- 倒计时 (超时自动取消)
- 进入聊天室按钮
- 取消订单
- 确认完成
- 申请退款
- 评价入口

**API 依赖:**
```
GET /api/v1/orders/{id}
PUT /api/v1/orders/{id}/cancel
PUT /api/v1/orders/{id}/complete
GET /api/v1/chat/groups (订单群聊)
```

**组件:**
- `OrderDetailHeader` - 订单头部
- `OrderTimeline` - 订单时间线
- `OrderActions` - 订单操作按钮
- `CountdownTimer` - 倒计时器
- `StatusBadge` - 状态徽章

---

### 4.3 我的订单 `/orders`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/orders` |
| **访问角色** | User |
| **优先级** | P0 (核心) |

**功能点:**
- 订单列表 (全部/进行中/已完成/已取消)
- 状态筛选
- 订单搜索
- 下拉刷新
- 分页加载

**API 依赖:**
```
GET /api/v1/orders
  ?status={pending|confirmed|in_progress|completed|canceled}
  &page={int}
  &pageSize={int}
```

**组件:**
- `OrderCard` - 订单卡片
- `OrderTabs` - 订单状态标签
- `PullToRefresh` - 下拉刷新

---

### 4.4 评价订单 `/order/:id/review`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/order/:id/review` |
| **访问角色** | User |
| **优先级** | P1 |

**功能点:**
- 星级评分
- 标签选择 (服务态度/技术/声音等)
- 文字评价
- 图片上传
- 匿名选项

**API 依赖:**
```
GET /api/v1/orders/{id}
POST /api/v1/reviews
```

**组件:**
- `StarRating` - 星级评分
- `TagSelector` - 标签选择器
- `ImageUploader` - 图片上传器
- `ReviewForm` - 评价表单

---

## 5. 聊天通讯

### 5.1 聊天列表 `/chat`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/chat` |
| **访问角色** | User |
| **优先级** | P0 (核心) |

**功能点:**
- 所有聊天群组列表
- 未读消息数显示
- 最后消息预览
- 在线状态
- 消息时间
- 左滑删除/置顶

**API 依赖:**
```
GET /api/v1/chat/groups
  ?page={int}
  &pageSize={int}
```

**组件:**
- `ChatGroupCard` - 聊天群组卡片
- `UnreadBadge` - 未读徽章
- `SwipeActions` - 滑动操作

---

### 5.2 聊天室 `/chat/:groupId`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/chat/:groupId` |
| **访问角色** | User |
| **优先级** | P0 (核心) |

**功能点:**
- 实时消息展示 (WebSocket)
- 消息发送 (文本/图片)
- 消息回复
- 消息举报
- 输入状态提示
- 消息已读状态
- 图片预览
- 消息分页加载

**API 依赖:**
```
GET /api/v1/chat/groups/{id}/messages
  ?page={int}
  &pageSize={int}
  &beforeId={uint64}
  &afterId={uint64}
POST /api/v1/chat/groups/{id}/messages
POST /api/v1/chat/messages/{id}/report
WebSocket: ws://host/ws/chat/{groupId}
```

**组件:**
- `MessageBubble` - 消息气泡
- `MessageInput` - 消息输入框
- `ImagePreview` - 图片预览
- `ReplyPreview` - 回复预览
- `TypingIndicator` - 输入指示器

---

## 6. 支付与钱包

### 6.1 支付页面 `/payment/:orderId`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/payment/:orderId` |
| **访问角色** | User |
| **优先级** | P0 (核心) |

**功能点:**
- 订单信息确认
- 支付方式选择 (微信/支付宝/余额)
- 余额支付
- VIP 折扣显示
- 支付密码输入
- 支付结果处理

**API 依赖:**
```
GET /api/v1/orders/{id}
GET /api/v1/user/wallet
POST /api/v1/payments
GET /api/v1/payments/{id}/status
```

**组件:**
- `PaymentMethods` - 支付方式选择
- `OrderSummary` - 订单汇总
- `PasswordInput` - 密码输入
- `PaymentResult` - 支付结果

---

### 6.2 我的钱包 `/wallet`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/wallet` |
| **访问角色** | User |
| **优先级** | P1 |

**功能点:**
- 余额显示
- 充值入口
- 交易记录
- 收入/支出筛选
- 提现入口 (陪玩师)

**API 依赖:**
```
GET /api/v1/user/wallet
GET /api/v1/user/wallet/transactions
  ?type={recharge|consume|income}
  &page={int}
```

**组件:**
- `WalletBalance` - 钱包余额卡片
- `TransactionList` - 交易记录列表
- `QuickActions` - 快捷操作按钮

---

### 6.3 充值 `/wallet/recharge`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/wallet/recharge` |
| **访问角色** | User |
| **优先级** | P1 |

**功能点:**
- 充值金额选择
- 自定义金额输入
- 充值活动展示
- 支付方式选择
- 充值记录

**API 依赖:**
```
GET /api/v1/user/recharge/plans
GET /api/v1/user/recharge/activities
POST /api/v1/user/recharge
```

**组件:**
- `AmountSelector` - 金额选择器
- `ActivityBanner` - 活动横幅
- `RechargeRecord` - 充值记录

---

## 7. 个人中心

### 7.1 个人中心 `/profile`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/profile` |
| **访问角色** | User |
| **优先级** | P0 (核心) |

**功能点:**
- 用户信息展示
- VIP 状态
- 优惠券数量
- 订单统计
- 快捷功能入口
- 设置入口

**API 依赖:**
```
GET /api/v1/user/me
GET /api/v1/user/vip/status
GET /api/v1/user/coupons/count
GET /api/v1/orders/stats
```

**组件:**
- `UserProfileHeader` - 用户头部信息
- `VIPBadge` - VIP 徽章
- `QuickActions` - 快捷操作
- `MenuList` - 菜单列表

---

### 7.2 个人资料编辑 `/profile/edit`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/profile/edit` |
| **访问角色** | User |
| **优先级** | P1 |

**功能点:**
- 昵称编辑
- 头像上传
- 性别设置
- 生日设置
- 个人简介
- 地区设置

**API 依赖:**
```
GET /api/v1/user/me
PUT /api/v1/user/profile
POST /api/v1/upload/avatar
```

**组件:**
- `ProfileForm` - 资料表单
- `AvatarUploader` - 头像上传器
- `RegionPicker` - 地区选择器

---

### 7.3 地址管理 `/profile/addresses`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/profile/addresses` |
| **访问角色** | User |
| **优先级** | P2 |

**功能点:**
- 地址列表
- 添加地址
- 编辑地址
- 删除地址
- 默认地址设置

**API 依赖:**
```
GET /api/v1/user/addresses
POST /api/v1/user/addresses
PUT /api/v1/user/addresses/{id}
DELETE /api/v1/user/addresses/{id}
```

**组件:**
- `AddressCard` - 地址卡片
- `AddressForm` - 地址表单

---

## 8. 营销功能

### 8.1 优惠券中心 `/coupons`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/coupons` |
| **访问角色** | User |
| **优先级** | P1 |

**功能点:**
- 可用优惠券列表
- 已使用/已过期筛选
- 领取优惠券
- 优惠券详情
- 使用说明

**API 依赖:**
```
GET /api/v1/user/coupons
  ?status={available|used|expired}
POST /api/v1/user/coupons/{id}/claim
```

**组件:**
- `CouponCard` - 优惠券卡片
- `CouponDetail` - 优惠券详情

---

### 8.2 VIP 中心 `/vip`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/vip` |
| **访问角色** | User |
| **优先级** | P1 |

**功能点:**
- VIP 等级展示
- VIP 权益说明
- 解锁条件展示
- 当前进度
- 升级按钮
- 充值入口

**API 依赖:**
```
GET /api/v1/user/vip/status
GET /api/v1/user/vip/levels
GET /api/v1/user/vip/threshold
```

**组件:**
- `VIPLevelCard` - VIP 等级卡片
- `BenefitList` - 权益列表
- `ProgressBar` - 进度条

---

### 8.3 活动中心 `/activities`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/activities` |
| **访问角色** | Guest / User |
| **优先级** | P1 |

**功能点:**
- 活动列表
- 活动详情
- 活动参与
- 分享功能

**API 依赖:**
```
GET /api/v1/activities
GET /api/v1/activities/{id}
POST /api/v1/activities/{id}/join
```

**组件:**
- `ActivityCard` - 活动卡片
- `ActivityDetail` - 活动详情

---

### 8.4 推荐有礼 `/referral`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/referral` |
| **访问角色** | User |
| **优先级** | P2 |

**功能点:**
- 推荐码展示
- 推荐记录
- 奖励记录
- 分享推荐码
- 推荐规则说明

**API 依赖:**
```
GET /api/v1/user/referral/code
GET /api/v1/user/referral/records
GET /api/v1/user/referral/rewards
```

**组件:**
- `ReferralCodeCard` - 推荐码卡片
- `RecordList` - 记录列表
- `ShareButton` - 分享按钮

---

## 9. 陪玩师专用

### 9.1 申请成为陪玩师 `/player/apply`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/player/apply` |
| **访问角色** | User (非陪玩师) |
| **优先级** | P1 |

**功能点:**
- 申请表单填写
- 真实姓名
- 身份证号
- 擅长游戏选择
- 段位证明上传
- 个人简介
- 审核状态展示

**API 依赖:**
```
POST /api/v1/player/apply
GET /api/v1/player/profile (检查状态)
GET /api/v1/games
```

**组件:**
- `ApplicationForm` - 申请表单
- `GameSelector` - 游戏选择器
- `ImageUploader` - 图片上传器
- `StatusTracker` - 状态追踪

---

### 9.2 陪玩师工作台 `/player/dashboard`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/player/dashboard` |
| **访问角色** | Player |
| **优先级** | P0 (陪玩师核心) |

**功能点:**
- 今日收益统计
- 订单统计
- 在线状态切换
- 接单开关
- 待处理订单
- 快捷入口

**API 依赖:**
```
GET /api/v1/player/profile
GET /api/v1/player/earnings/today
GET /api/v1/player/orders?status=pending
PUT /api/v1/player/status
```

**组件:**
- `EarningsCard` - 收益卡片
- `OrderStats` - 订单统计
- `OnlineToggle` - 在线状态开关
- `QuickActions` - 快捷操作

---

### 9.3 陪玩师资料编辑 `/player/profile/edit`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/player/profile/edit` |
| **访问角色** | Player |
| **优先级** | P1 |

**功能点:**
- 基础信息编辑
- 擅长游戏管理
- 服务时段设置
- 价格查看 (系统统一定价)
- 相册管理
- 语音介绍

**API 依赖:**
```
GET /api/v1/player/profile
PUT /api/v1/player/profile
POST /api/v1/upload/voice
POST /api/v1/upload/gallery
```

**组件:**
- `PlayerProfileForm` - 陪玩师资料表单
- `GameRankSelector` - 游戏段位选择器
- `TimeSlotPicker` - 时段选择器
- `VoiceRecorder` - 语音录制器

---

### 9.4 陪玩师订单 `/player/orders`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/player/orders` |
| **访问角色** | Player |
| **优先级** | P0 (陪玩师核心) |

**功能点:**
- 订单列表
- 状态筛选
- 接单/拒单
- 开始服务
- 完成服务
- 订单详情

**API 依赖:**
```
GET /api/v1/player/orders
  ?status={pending|confirmed|in_progress|completed}
PUT /api/v1/player/orders/{id}/accept
PUT /api/v1/player/orders/{id}/reject
PUT /api/v1/player/orders/{id}/start
PUT /api/v1/player/orders/{id}/complete
```

**组件:**
- `PlayerOrderCard` - 订单卡片
- `OrderActions` - 订单操作按钮
- `AcceptModal` - 接单确认弹窗

---

### 9.5 收益中心 `/player/earnings`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/player/earnings` |
| **访问角色** | Player |
| **优先级** | P0 (陪玩师核心) |

**功能点:**
- 收益总览
- 收益明细
- 提现申请
- 提现记录
- 抽成说明

**API 依赖:**
```
GET /api/v1/player/earnings
GET /api/v1/player/earnings/records
GET /api/v1/player/commission
POST /api/v1/player/withdraw
```

**组件:**
- `EarningsSummary` - 收益汇总
- `EarningsChart` - 收益图表
- `WithdrawForm` - 提现表单
- `RecordList` - 记录列表

---

## 10. 设置与辅助

### 10.1 设置 `/settings`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/settings` |
| **访问角色** | User |
| **优先级** | P1 |

**功能点:**
- 账号设置
- 通知设置
- 隐私设置
- 主题切换
- 语言设置
- 关于我们
- 退出登录

**API 依赖:**
```
GET /api/v1/user/settings
PUT /api/v1/user/settings
POST /api/v1/auth/logout
```

**组件:**
- `SettingsList` - 设置列表
- `SwitchItem` - 开关项
- `ThemeSelector` - 主题选择器

---

### 10.2 黑名单管理 `/settings/blocked`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/settings/blocked` |
| **访问角色** | User |
| **优先级** | P2 |

**功能点:**
- 已拉黑用户列表
- 解除拉黑
- 添加拉黑

**API 依赖:**
```
GET /api/v1/user/blocks
POST /api/v1/user/blocks
DELETE /api/v1/user/blocks/{id}
```

**组件:**
- `BlockedUserCard` - 拉黑用户卡片
- `UnblockButton` - 解除按钮

---

### 10.3 争议申诉 `/dispute/:orderId`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/dispute/:orderId` |
| **访问角色** | User / Player |
| **优先级** | P2 |

**功能点:**
- 争议类型选择
- 问题描述
- 证据上传
- 提交申诉
- 申诉进度查询

**API 依赖:**
```
GET /api/v1/dispute/order/{orderId}
POST /api/v1/dispute
GET /api/v1/dispute/{id}
```

**组件:**
- `DisputeForm` - 申诉表单
- `EvidenceUploader` - 证据上传器
- `ProgressTracker` - 进度追踪

---

### 10.4 帮助中心 `/help`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/help` |
| **访问角色** | Guest / User |
| **优先级** | P2 |

**功能点:**
- 常见问题
- 问题分类
- 搜索功能
- 联系客服

**组件:**
- `FAQList` - 常见问题列表
- `CategoryNav` - 分类导航
- `ContactSupport` - 联系客服

---

### 10.5 消息通知 `/notifications`

| 属性 | 值 |
|:-----|:---|
| **路由** | `/notifications` |
| **访问角色** | User |
| **优先级** | P1 |

**功能点:**
- 通知列表
- 已读/未读筛选
- 批量已读
- 通知详情
- 跳转对应页面

**组件:**
- `NotificationCard` - 通知卡片
- `EmptyState` - 空状态

---

## 页面优先级汇总

```
┌────────────────────────────────────────────────────────────┐
│                   优先级分布统计                            │
├────────────────────────────────────────────────────────────┤
│  P0 (核心 - 必须实现):  15 页                              │
│  P1 (重要 - 第一批):    13 页                              │
│  P2 (常规 - 后续迭代):  6 页                               │
├────────────────────────────────────────────────────────────┤
│  总计: 34 个主要页面                                       │
└────────────────────────────────────────────────────────────┘
```

### P0 核心页面 (MVP 必备)

| 页面 | 用户角色 | 说明 |
|:-----|:---------|:-----|
| `/login` | Guest | 登录入口 |
| `/register` | Guest | 注册入口 |
| `/` | Guest/User | 首页 |
| `/players` | Guest/User | 陪玩师列表 |
| `/player/:id` | Guest/User | 陪玩师详情 |
| `/order/create` | User | 创建订单 |
| `/order/:id` | User | 订单详情 |
| `/orders` | User | 我的订单 |
| `/chat` | User | 聊天列表 |
| `/chat/:groupId` | User | 聊天室 |
| `/payment/:orderId` | User | 支付页面 |
| `/profile` | User | 个人中心 |
| `/player/dashboard` | Player | 陪玩师工作台 |
| `/player/orders` | Player | 陪玩师订单 |
| `/player/earnings` | Player | 收益中心 |

---

## 权限矩阵

```
┌─────────────────────────────────────────────────────────────────────┐
│                        页面访问权限矩阵                              │
├──────────┬───────────────────────────────────────────────────────────┤
│  页面     │  Guest  │  User   │  Player  │  说明                    │
├──────────┼─────────┼─────────┼──────────┼───────────────────────────┤
│ /login   │    ✅   │    ✅   │    ✅    │  已登录自动跳转首页         │
│ /register│    ✅   │    ✅   │    ✅    │  已登录自动跳转首页         │
│ /        │    ✅   │    ✅   │    ✅    │  首页                      │
│ /players │    ✅   │    ✅   │    ✅    │  陪玩师列表                │
│ /order/* │    ❌   │    ✅   │    ✅    │  需要登录                  │
│ /chat/*  │    ❌   │    ✅   │    ✅    │  需要登录                  │
│ /wallet  │    ❌   │    ✅   │    ✅    │  需要登录                  │
│ /player/dashboard │ ❌ │ ❌ │ ✅ │ 仅陪玩师                   │
│ /player/earnings │ ❌ │ ❌ │ ✅ │ 仅陪玩师                   │
│ /admin/* │    ❌   │    ❌   │    ❌    │  管理员专用 (跳转管理端)   │
└──────────┴─────────┴─────────┴──────────┴───────────────────────────┘
```

---

## 组件复用策略

### 核心组件库

```
┌─────────────────────────────────────────────────────────────┐
│                    组件复用优先级                            │
├─────────────────────────────────────────────────────────────┤
│  🔥 高复用 (5+ 处):                                         │
│    - PlayerCard        (首页/列表/详情/收藏/订单)           │
│    - OrderCard         (订单列表/详情/工作台)               │
│    - StatusBadge       (所有状态展示)                       │
│    - EmptyState        (所有空状态)                         │
│    - Loading           (所有加载状态)                       │
├─────────────────────────────────────────────────────────────┤
│  ⚡ 中复用 (3-4 处):                                        │
│    - GameFilter        (列表/详情/申请)                     │
│    - ReviewList        (详情/评价页)                        │
│    - PriceSummary      (订单/支付)                          │
│    - MessageBubble     (所有聊天场景)                       │
├─────────────────────────────────────────────────────────────┤
│  📄 低复用 (1-2 处):                                        │
│    - ApplicationForm  (申请页专用)                          │
│    - WithdrawForm     (提现专用)                            │
│    - DisputeForm      (申诉专用)                            │
└─────────────────────────────────────────────────────────────┘
```

---

## 路由配置建议

```typescript
// 路由结构示例
const routes = [
  // 公开页面
  { path: '/', component: HomePage, meta: { auth: false } },
  { path: '/login', component: LoginPage, meta: { guest: true } },
  { path: '/register', component: RegisterPage, meta: { guest: true } },

  // 需要登录的页面
  { path: '/orders', component: MyOrdersPage, meta: { auth: true } },
  { path: '/order/:id', component: OrderDetailPage, meta: { auth: true } },
  { path: '/order/create', component: CreateOrderPage, meta: { auth: true } },

  // 陪玩师专用
  { path: '/player/dashboard', component: PlayerDashboard, meta: { auth: true, role: 'player' } },

  // 重定向
  { path: '/admin', redirect: 'https://admin.gamelink.com' }
]
```

---

## API 对接清单

### 认证相关
```
POST /api/v1/auth/login
POST /api/v1/auth/register
POST /api/v1/auth/logout
POST /api/v1/auth/send-code
POST /api/v1/auth/refresh
```

### 陪玩师相关
```
GET /api/v1/players
GET /api/v1/players/:id
GET /api/v1/favorites/players
POST /api/v1/favorites/players/:id
DELETE /api/v1/favorites/players/:id
```

### 订单相关
```
POST /api/v1/orders
GET /api/v1/orders
GET /api/v1/orders/:id
PUT /api/v1/orders/:id/cancel
PUT /api/v1/orders/:id/complete
```

### 聊天相关
```
GET /api/v1/chat/groups
GET /api/v1/chat/groups/:id/messages
POST /api/v1/chat/groups/:id/messages
WebSocket: /ws/chat/:groupId
```

---

**文档版本**: 1.0.0
**创建日期**: 2025-01-11
**状态**: 待评审
**负责人**: Super Dev Team (PM + UX)
