# GameLink 前端重构计划

> 技术栈：uni-app + Vue 3 + TypeScript + Vite + Pinia + uv-ui + SCSS  
> 工作区：`app/src/`  
> 最后更新：2026-01-31  
> **状态：✅ 开发完成，待联调**

**注意**：原生微信小程序目录 `miniprogram/` 已删除，项目完全迁移至 uni-app。

**相关文档**：
- [组件化封装计划](./COMPONENTIZATION_PLAN.md) - 组件二次封装与分层架构
- [页面组件使用计划](./PAGE_COMPONENT_USAGE_PLAN.md) - 各页面组件使用规划
- [组件职责文档](./COMPONENTS.md) - 现有组件 API 与使用说明

---

## 一、已完成内容

### Phase 1 - 基础架构 ✅

- [x] uni-app 项目初始化
- [x] 依赖集成：@climblee/uv-ui、pinia、sass、dayjs
- [x] 主题系统：日间亮绿色 #00D26A，夜间紫色 #7C3AED
- [x] API 封装：`src/api/request.ts`（统一请求、Token 刷新、错误处理）
- [x] 认证 API：`src/api/auth.ts`（登录、注册、微信登录）

### Phase 2 - 认证 + 组件 ✅

- [x] 登录页：`src/pages/auth/login/index.vue`（微信一键登录 + 账号密码）
- [x] 注册页：`src/pages/auth/register/index.vue`（角色选择：用户/陪玩师）
- [x] 首页：`src/pages/index/index.vue`（引导页 + 快捷入口）
- [x] 公共组件 (7个)：
  - GlTag - 标签/状态
  - GlEmpty - 空状态
  - PlayerCard - 陪玩师卡片
  - OrderCard - 订单卡片
  - RatingStars - 评分星星
  - PriceTag - 价格标签
  - ImageUploader - 图片上传

### Phase 7 - 补充页面 ✅

| 页面 | 路径 | 功能 |
|------|------|------|
| 公共频道列表 | `pages/channel/list` | Discord 风格频道入口 |
| 支付结果页 | `pages/payment/result` | 支付成功/失败/处理中 |
| 协议页 | `pages/agreement` | 用户协议/隐私政策/陪玩师协议 |
| 帮助中心 | `pages/help` | FAQ + 分类搜索 |
| 在线客服 | `pages/service` | 快捷问题 + 实时聊天 |

### Phase 8 - 性能优化 + 发布准备 ✅

- [x] LazyImage - 图片懒加载（骨架屏 + 自动重试）
- [x] VirtualList - 虚拟滚动列表
- [x] ErrorBoundary - 统一错误处理组件
- [x] PrivacyPopup - 小程序隐私授权弹窗
- [x] manifest.json - 完整配置（权限、隐私检查）
- [x] 环境变量配置（.env.development / .env.production）

---

## 二、待完成任务

### Phase 3 - 剩余组件 (8个) ✅

| 组件 | 说明 | 状态 |
|------|------|------|
| GameSelector | 游戏选择器 | ✅ 已完成 |
| RankSelector | 段位选择器 | ✅ 已完成 |
| CustomTabBar | 底部导航栏 | ✅ 已完成 |
| NavBar | 顶部导航栏 | ✅ 已完成 |
| LoadMore | 加载更多 | ✅ 已完成 |
| Skeleton | 骨架屏 | ✅ 已完成 |
| ChatMessageBubble | 聊天气泡 | ✅ 已完成 |
| ThemeToggle | 主题切换 | ✅ 已完成 |

### Phase 4 - 用户端页面 (15个)

| 页面 | 路径 | 状态 |
|------|------|------|
| 陪玩师列表 | `pages/player/list` | ✅ 已完成 |
| 陪玩师详情 | `pages/player/detail` | ✅ 已完成 |
| 下单页 | `pages/order/create` | ✅ 已完成 |
| 订单列表 | `pages/order/list` | ✅ 已完成 |
| 订单详情 | `pages/order/detail` | ✅ 已完成 |
| 钱包 | `pages/wallet/index` | ✅ 已完成 |
| 充值 | `pages/wallet/recharge` | ✅ 已完成 |
| 消息中心 | `pages/message/list` | ✅ 已完成 |
| 聊天详情 | `pages/message/chat` | ✅ 已完成 |
| 个人中心 | `pages/profile/index` | ✅ 已完成 |
| 编辑资料 | `pages/profile/edit` | ✅ 已完成 |
| 收藏列表 | `pages/favorite/list` | ✅ 已完成 |
| 游戏列表 | `pages/game/list` | ✅ 已完成 |
| 评价列表 | `pages/review/list` | ✅ 已完成 |
| 设置 | `pages/settings/index` | ✅ 已完成 |

### Phase 5 - 陪玩师端页面 (5个核心)

| 页面 | 路径 | 状态 |
|------|------|------|
| 工作台 | `pages/player/dashboard` | ✅ 已完成 |
| 订单管理 | `pages/player/orders` | ✅ 已完成 |
| 收益中心 | `pages/player/earnings` | ✅ 已完成 |
| 服务管理 | `pages/player/services` | ✅ 已完成 |
| 陪玩认证 | `pages/player/certification` | ✅ 已完成 |

### Phase 6 - 聊天系统（Discord 风格）

**架构说明**：类似 Discord 的频道系统，平台可监控所有聊天内容

```
GameLink 平台
├── 📢 公共频道 - 大厅、游戏讨论区
├── 🔒 订单频道 - 用户+陪玩师+客服（自动创建/销毁）
└── 💬 私信 (DM) - 两人对话（后台可见）
```

| 群组类型 | 说明 | 监控 |
|----------|------|------|
| `public` | 公共聊天室，所有用户可加入 | ✅ |
| `order` | 订单群组，订单完成后销毁 | ✅ |
| `private` | 私聊，两人对话 | ✅ 后台可查看 |

---

## 三、API 接口参考

### 3.1 登录认证（微信一键）

```typescript
// 小程序登录
POST /api/v1/public/auth/wechat/login
// 兼容路径
POST /api/v1/auth/wechat/login

// 请求体
{
  code: string,           // wx.login 获取
  encryptedData?: string, // 可选：获取用户信息
  iv?: string,
  referralCode?: string   // 可选：邀请码
}

// 返回
{
  accessToken: string,
  refreshToken: string,
  user: { ... }
}

// Token 刷新
POST /api/v1/public/auth/refresh
POST /api/v1/auth/wechat/refresh  // 兼容
```

### 3.2 图片上传（COS 签名 URL）

```typescript
// 头像
POST /api/v1/user/avatar
FormData: { avatar: File }

// 认证材料
POST /api/v1/certification/upload/image
FormData: { image: File, type: 'id-card' | 'skill-proof' }

// 聊天图片
POST /api/v1/chat/upload/image
FormData: { image: File }

// 评价图片（多文件）
POST /api/v1/review/upload/images
FormData: { images: File[] }

// ⚠️ 返回的 URL 是签名 URL，默认 1 小时过期
// 可通过 OSS_SIGNED_URL_TTL_SECONDS 调整
```

### 3.3 支付接口（当前 Mock）

```typescript
// 创建支付
POST /api/v1/user/payments

// 请求体
{
  orderId: number,
  method: 'wechat' | 'alipay' | 'wallet' | 'combined',
  requestId: string,  // UUID 防重复
  // 组合支付时：
  walletAmountCents?: number,
  thirdPartyMethod?: 'wechat' | 'alipay'
}

// 返回 payInfo（Mock 模式自动标记成功）
// wechat: { prepay_id, code_url }
// alipay: { trade_no, qr_code }
// wallet: { status: 'paid', newBalance }

// 回调接口（暂不对接）
POST /api/v1/public/payments/wechat/notify
POST /api/v1/public/payments/alipay/notify
```

### 3.4 WebSocket（聊天 & 通知）

```typescript
// 连接地址
ws(s)://{host}/api/v1/ws

// 鉴权方式（三选一）
1. Header: Authorization: Bearer <token>
2. Query: ?token=<token> 或 ?access_token=<token>
3. Cookie: auth_token（需启用 USE_COOKIE_AUTH）

// 消息格式
{
  type: 'chat_message' | 'order_status' | 'order_new' | 'notification',
  timestamp: number,
  data: { ... }
}

// 消息类型说明
// chat_message  - 聊天消息
// order_status  - 订单状态变化
// order_new     - 新订单通知（推给陪玩师）
// notification  - 站内通知
```

### 3.5 消息推送（当前支持）

| 方式 | 状态 |
|------|------|
| 站内通知 + WS 推送 | ✅ 已实现 |
| 小程序模板消息 | ❌ 未接入 |
| App Push | ❌ 未接入 |
| 短信通知 | ❌ 未接入 |

---

## 四、关键文件说明

| 文件 | 作用 |
|------|------|
| `src/store/user.ts` | 用户状态管理（Token、用户信息、登录/登出） |
| `src/composables/useTheme.ts` | 主题切换 Hook |
| `src/styles/variables.scss` | 主题色彩变量 |
| `src/pages.json` | 页面路由配置 |
| `vite.config.ts` | Vite 配置（别名 @、SCSS 全局变量） |
| `src/api/request.ts` | 统一请求封装 |
| `src/api/auth.ts` | 认证相关 API |
| `src/api/chat.ts` | 聊天相关 API（Discord 风格） |
| `src/api/player.ts` | 陪玩师相关 API（服务、排班、统计） |
| `src/api/notification.ts` | 通知相关 API |
| `src/api/settings.ts` | 用户设置 API |
| `src/api/order.ts` | 订单相关 API |
| `src/api/wallet.ts` | 钱包相关 API |
| `src/api/game.ts` | 游戏相关 API |
| `src/api/favorite.ts` | 收藏相关 API |
| `src/store/user.ts` | 用户状态管理 |
| `src/store/app.ts` | 应用全局状态（主题、网络） |
| `src/store/player.ts` | 陪玩师状态管理 |
| `src/composables/useAuth.ts` | 认证相关 Hook |
| `src/composables/useTheme.ts` | 主题切换 Hook |
| `src/composables/useWebSocket.ts` | WebSocket 连接管理 |
| `src/composables/useLoading.ts` | 加载状态管理 |
| `src/composables/usePagination.ts` | 分页管理 |
| `src/utils/format.ts` | 格式化函数（金额、时间、脱敏） |
| `src/utils/validate.ts` | 验证函数（手机、邮箱、身份证） |
| `src/utils/storage.ts` | 本地存储工具（带过期时间） |
| `src/utils/index.ts` | 工具函数（防抖、节流、UUID、深拷贝）|
| `src/api/publicPlayer.ts` | 公开陪玩师 API（列表、详情、搜索） |
| `src/api/user.ts` | 用户信息 API（资料、头像） |
| `src/api/certification.ts` | 陪玩认证 API |

---

## 五、后端 API 对接状态 ✅

后端已实现以下接口，前端已完成对接：

### 5.1 聊天系统（Discord 风格）

| 接口 | 路径 | 前端对接 |
|------|------|----------|
| 创建聊天群组 | `POST /user/chat/groups` | ✅ |
| 获取群组详情 | `GET /user/chat/groups/:id` | ✅ |
| 加入群组 | `POST /user/chat/groups/:id/join` | ✅ |
| 离开群组 | `POST /user/chat/groups/:id/leave` | ✅ |
| 标记已读 | `POST /user/chat/groups/:id/read` | ✅ |
| 公共频道列表 | `GET /public/chat/public-channels` | ✅ |

### 5.2 陪玩师相关

| 接口 | 路径 | 前端对接 |
|------|------|----------|
| 服务列表 | `GET /player/services` | ✅ |
| 创建服务 | `POST /player/services` | ✅ |
| 更新服务 | `PUT /player/services/:id` | ✅ |
| 删除服务 | `DELETE /player/services/:id` | ✅ |
| 切换服务状态 | `PUT /player/services/:id/status` | ✅ |
| 获取排班 | `GET /player/schedule` | ✅ |
| 更新排班 | `PUT /player/schedule` | ✅ |
| 今日统计 | `GET /player/stats/today` | ✅ |
| 总览统计 | `GET /player/stats/overview` | ✅ |
| 公开评价 | `GET /public/players/:id/reviews` | ✅ |

### 5.3 用户设置

| 接口 | 路径 | 前端对接 |
|------|------|----------|
| 获取设置 | `GET /user/settings` | ✅ |
| 更新设置 | `PUT /user/settings` | ✅ |
| 通知设置获取 | `GET /user/notification-settings` | ✅ |
| 通知设置更新 | `PUT /user/notification-settings` | ✅ |

### 5.4 通知系统

| 接口 | 路径 | 前端对接 |
|------|------|----------|
| 通知列表 | `GET /user/notifications` | ✅ |
| 未读数量 | `GET /user/notifications/unread-count` | ✅ |
| 标记已读 | `POST /user/notifications/:id/read` | ✅ |
| 全部已读 | `POST /user/notifications/read-all` | ✅ |

### 5.5 公开陪玩师 API

| 接口 | 路径 | 前端对接 |
|------|------|----------|
| 陪玩师列表 | `GET /public/players` | ✅ |
| 陪玩师详情 | `GET /public/players/:id` | ✅ |
| 陪玩师服务 | `GET /public/players/:id/services` | ✅ |
| 推荐陪玩师 | `GET /public/players/recommended` | ✅ |
| 热门陪玩师 | `GET /public/players/hot` | ✅ |
| 搜索陪玩师 | `GET /public/search/players` | ✅ |

### 5.6 用户资料 API

| 接口 | 路径 | 前端对接 |
|------|------|----------|
| 获取用户资料 | `GET /user/profile` | ✅ |
| 更新用户资料 | `PUT /user/profile` | ✅ |
| 上传头像 | `POST /user/avatar` | ✅ |
| 修改密码 | `PUT /user/password` | ✅ |

### 5.7 陪玩认证 API

| 接口 | 路径 | 前端对接 |
|------|------|----------|
| 获取认证状态 | `GET /player/certification` | ✅ |
| 提交认证 | `POST /player/certification` | ✅ |
| 上传身份证 | `POST /certification/upload/image` | ✅ |

### 5.8 陪玩师订单管理 API

| 接口 | 路径 | 前端对接 |
|------|------|----------|
| 陪玩师订单列表 | `GET /player/orders` | ✅ |
| 接单 | `POST /player/orders/:id/accept` | ✅ |
| 完成订单 | `POST /player/orders/:id/complete` | ✅ |

---

## 六、页面 API 对接完成状态

### 6.1 用户端核心页面

| 页面 | API 对接 |
|------|----------|
| 陪玩师列表 | ✅ getPlayerList, getHotGames |
| 陪玩师详情 | ✅ getPlayerDetail, getPlayerReviews, addFavorite, removeFavorite |
| 下单页 | ✅ getPlayerDetail, createOrder |
| 订单列表 | ✅ getOrders |
| 订单详情 | ✅ getOrderDetail, cancelOrder, completeOrder, payOrder, submitReview |
| 钱包 | ✅ getWalletInfo, getTransactions |
| 充值 | ✅ getWalletInfo, recharge |
| 消息中心 | ✅ getChatGroups, getNotifications, getUnreadCount |
| 聊天详情 | ✅ getChatGroupDetail, getChatMessages, sendChatMessage, markMessagesRead |
| 个人中心 | ✅ getUserProfile, getWalletInfo |
| 编辑资料 | ✅ getUserProfile, updateUserProfile, uploadAvatar |
| 收藏列表 | ✅ getFavorites, batchRemoveFavorites |
| 游戏列表 | ✅ getGameCategories, getGames |
| 设置 | ✅ getUserSettings, updateUserSettings, getNotificationSettings, updateNotificationSettings |

### 6.2 陪玩师端核心页面

| 页面 | API 对接 |
|------|----------|
| 工作台 | ✅ getTodayStats |
| 订单管理 | ✅ getPlayerOrders, acceptPlayerOrder, completePlayerOrder |
| 收益中心 | ✅ getOverviewStats |
| 服务管理 | ✅ getPlayerServices, createPlayerService, updatePlayerService, deletePlayerService, toggleServiceStatus |
| 陪玩认证 | ✅ getCertificationStatus, submitCertification, uploadIdCardImage, uploadRankScreenshot |

---

## 七、启动命令

```bash
# 进入 app 目录
cd app

# 安装依赖
npm install

# H5 开发
npm run dev:h5

# 微信小程序
npm run dev:mp-weixin

# 类型检查
npm run type-check
```

---

## 八、注意事项

1. ~~`app/miniprogram/` 是原生微信小程序代码~~ **已删除，完全迁移至 uni-app**
2. API 基础 URL 通过环境变量配置：`.env.development` / `.env.production`
3. 微信小程序需要在 `src/manifest.json` 配置 AppID
4. 主题色：日间 `#00D26A` / 夜间 `#7C3AED`
5. **COS 图片 URL 有效期 1 小时**，需注意缓存策略
6. **支付当前是 Mock 模式**，创建后自动标记成功
7. **聊天系统为 Discord 风格**，后台可监控所有消息
8. CustomTabBar 图标需替换为真实 PNG（见 `static/icons/README.md`）

---

## 九、组件清单（19 个）

| 组件 | 功能 |
|------|------|
| ChatMessageBubble | 聊天气泡 |
| GlEmpty | 空状态 |
| ErrorBoundary | 错误边界 |
| GameSelector | 游戏选择器 |
| ImageUploader | 图片上传 |
| LazyImage | 图片懒加载 |
| LoadMore | 加载更多 |
| NavBar | 顶部导航 |
| OrderCard | 订单卡片 |
| PlayerCard | 陪玩师卡片 |
| PriceTag | 价格标签 |
| PrivacyPopup | 隐私授权弹窗 |
| RankSelector | 段位选择器 |
| RatingStars | 评分星星 |
| Skeleton | 骨架屏 |
| GlTag | 标签/状态 |
| CustomTabBar | 底部导航 |
| ThemeToggle | 主题切换 |
| VirtualList | 虚拟列表 |

---

## 十、页面清单（28 个）

### 用户端（20 个）

| 页面 | 路径 |
|------|------|
| 首页 | `pages/index/index` |
| 登录 | `pages/auth/login/index` |
| 注册 | `pages/auth/register/index` |
| 陪玩师列表 | `pages/player/list/index` |
| 陪玩师详情 | `pages/player/detail/index` |
| 下单 | `pages/order/create/index` |
| 订单列表 | `pages/order/list/index` |
| 订单详情 | `pages/order/detail/index` |
| 钱包 | `pages/wallet/index/index` |
| 充值 | `pages/wallet/recharge/index` |
| 消息列表 | `pages/message/list/index` |
| 聊天详情 | `pages/message/chat/index` |
| 个人中心 | `pages/profile/index/index` |
| 编辑资料 | `pages/profile/edit/index` |
| 收藏列表 | `pages/favorite/list/index` |
| 游戏列表 | `pages/game/list/index` |
| 评价列表 | `pages/review/list/index` |
| 设置 | `pages/settings/index/index` |
| 公共频道 | `pages/channel/list/index` |
| 支付结果 | `pages/payment/result/index` |
| 协议 | `pages/agreement/index` |
| 帮助中心 | `pages/help/index` |
| 在线客服 | `pages/service/index` |

### 陪玩师端（5 个）

| 页面 | 路径 |
|------|------|
| 工作台 | `pages/player/dashboard/index` |
| 订单管理 | `pages/player/orders/index` |
| 收益中心 | `pages/player/earnings/index` |
| 服务管理 | `pages/player/services/index` |
| 陪玩认证 | `pages/player/certification/index` |
