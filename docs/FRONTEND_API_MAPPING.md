# GameLink 前端 API 映射基线（规划 vs 现状）

更新时间：2026-02-12

## 1) 使用说明

本文件用于统一前端联调口径，避免“规划接口名”和“当前真实接口”混用。

- **规划接口**：产品/设计文档中的抽象路径（如 `/orders`）
- **真实接口**：当前前端实际调用的后端路径（如 `/user/orders`）
- **处理策略**：当前建议执行方式（直接使用/别名映射/待后端统一）

---

## 2) 模块映射总表（高优先）

### 2.1 认证（auth）

| 规划接口 | 真实接口（当前） | 方法 | 前端文件 | 处理策略 |
|---|---|---|---|---|
| `/auth/wechat-login` | `/public/auth/wechat/login` | POST | `app/src/api/auth.ts` | 以真实接口为准 |
| `/auth/login` | `/auth/login` | POST | `app/src/api/auth.ts` | 一致 |
| `/auth/refresh-token` | `/auth/refresh` | POST | `app/src/api/auth.ts` | 文档别名说明 |

### 2.2 首页/通用

| 规划接口 | 真实接口（当前） | 方法 | 前端文件 | 处理策略 |
|---|---|---|---|---|
| `/home/banners` | `/public/banners` | GET | `app/src/api/banner.ts` | 以真实接口为准 |
| `/games/hot` | `/public/games`（`page_size` 限制） | GET | `app/src/api/game.ts` | 文档注明“热门由查询参数实现” |
| `/notifications/unread-count` | `/user/notifications/unread-count` | GET | `app/src/api/notification.ts` | 统一走 `user` 前缀 |

### 2.3 陪玩师（public/player）

| 规划接口 | 真实接口（当前） | 方法 | 前端文件 | 处理策略 |
|---|---|---|---|---|
| `/players` | `/public/players` | GET | `app/src/api/publicPlayer.ts` | 以真实接口为准 |
| `/players/:id` | `/public/players/:id` | GET | `app/src/api/publicPlayer.ts` | 一致（加 `public`） |
| `/players/:id/services` | `/public/players/:id/services` | GET | `app/src/api/publicPlayer.ts` | 一致（加 `public`） |
| `/players/:id/reviews` | `/public/players/:id/reviews` | GET | `app/src/api/publicPlayer.ts` | 一致（加 `public`） |
| `/player/status` | `/player/online-status` | PUT | `app/src/api/player.ts` | 文档别名说明 |
| `/player/stats/today` | `/player/stats/today` | GET | `app/src/api/player.ts` | 一致 |
| `/player/stats/overview` | `/player/stats/overview` | GET | `app/src/api/player.ts` | 一致 |
| `/player/orders/:id/finish` | `/player/orders/:id/complete` | PUT | `app/src/api/order.ts` | 以真实接口为准 |

### 2.4 订单（order）

| 规划接口 | 真实接口（当前） | 方法 | 前端文件 | 处理策略 |
|---|---|---|---|---|
| `/orders` | `/user/orders` | GET/POST | `app/src/api/order.ts` | 统一“用户态用 `/user` 前缀” |
| `/orders/:id` | `/user/orders/:id` | GET | `app/src/api/order.ts` | 同上 |
| `/orders/:id/cancel` | `/user/orders/:id/cancel` | PUT | `app/src/api/order.ts` | 方法以真实接口为准 |
| `/orders/:id/complete` | `/user/orders/:id/complete` | PUT | `app/src/api/order.ts` | 方法以真实接口为准 |
| `/orders/:id/refund` | `/user/orders/:id/refund` | POST | `app/src/api/order.ts` | 一致（仅前缀差异） |
| `/orders/:id/review` | `/user/reviews` | POST | `app/src/api/order.ts` | 文档备注“独立评价接口” |
| `/orders/:id/pay` | `/user/payments` | POST | `app/src/api/wallet.ts` | 文档改为支付模块接口 |

### 2.5 钱包（wallet）

| 规划接口 | 真实接口（当前） | 方法 | 前端文件 | 处理策略 |
|---|---|---|---|---|
| `/wallet` | `/user/wallet/balance` | GET | `app/src/api/wallet.ts` | 文档别名说明 |
| `/wallet/transactions` | `/user/wallet/transactions` | GET | `app/src/api/wallet.ts` | 一致（加 `user`） |
| `/wallet/recharge` | `/user/wallet/recharge` | POST | `app/src/api/wallet.ts` | 一致（加 `user`） |
| `/wallet/vip` | `/user/vip/status` | GET | `app/src/api/wallet.ts` | 以真实接口为准 |
| `/wallet/withdraw` | （前端暂未接入统一接口） | POST | - | 待后端确认后接入 |

### 2.6 聊天/频道（chat）

| 规划接口 | 真实接口（当前） | 方法 | 前端文件 | 处理策略 |
|---|---|---|---|---|
| `/chat/groups` | `/user/chat/groups` | GET | `app/src/api/chat.ts` | 一致（加 `user`） |
| `/chat/groups/:id/messages` | `/user/chat/groups/:id/messages` | GET/POST | `app/src/api/chat.ts` | 一致（加 `user`） |
| `/chat/channels` | `/public/chat/public-channels` | GET | `app/src/api/chat.ts` | 以真实接口为准 |
| `/chat/channels/:id/join` | `/user/chat/groups/:id/join` | POST | `app/src/api/chat.ts` | 规划需改为 group 模型 |
| `WebSocket` | `wss://.../ws/chat/:id?token=...` | - | `app/src/composables/useChatRoom.ts` | 保持当前 |

### 2.7 用户与设置（user）

| 规划接口 | 真实接口（当前） | 方法 | 前端文件 | 处理策略 |
|---|---|---|---|---|
| `/user/profile` | `/user/profile`（补充：登录态用户信息也可用 `/auth/me`） | GET/PUT | `app/src/api/user.ts`、`app/src/api/auth.ts` | 以真实接口为准 |
| `/upload/avatar` | `/user/avatar` | POST(upload) | `app/src/api/user.ts` | 一致（加 `user`） |
| 设置类接口 | `/user/settings`、`/user/notification-settings` | GET/PUT | `app/src/api/settings.ts` | 一致（加 `user`） |

### 2.8 认证（certification）

| 规划接口 | 真实接口（当前） | 方法 | 前端文件 | 处理策略 |
|---|---|---|---|---|
| `/certification/status` | `/player/certification/identity` | GET | `app/src/api/certification.ts` | 文档别名说明 |
| `/certification/submit` | `/player/certification/identity` | POST | `app/src/api/certification.ts` | 文档别名说明 |
| `/upload/idcard` | `/certification/upload/image` | POST(upload) | `app/src/api/certification.ts` | 通过 `type/side` 区分 |
| `/upload/rank-screenshot` | `/certification/upload/image` | POST(upload) | `app/src/api/certification.ts` | 通过 `type` 区分 |

---

## 3) 页面链路验收建议（按优先级）

### P0：必须先通
1. 登录：`pages/auth/login/index`
2. 首页浏览：`pages/index/index`（banner/游戏/推荐陪玩师）
3. 陪玩师详情 -> 下单：`pages/player/detail/index` -> `pages/order/create/index`
4. 订单支付/状态回写：`pages/order/detail/index`
5. 钱包流水同步：`pages/wallet/index/index`

### P1：角色侧闭环
1. 陪玩师工作台：`pages/player/dashboard/index`
2. 陪玩师订单处理：`pages/player/orders/index`
3. 服务管理：`pages/player/services/index`
4. 排班管理：`pages/player/schedule/index`

### P2：社交与运营
1. 消息会话 + 聊天室：`pages/message/list/index`、`pages/message/chat/index`
2. 公共频道：`pages/channel/list/index`
3. 个人中心与收藏：`pages/profile/index/index`、`pages/favorite/list/index`

---

## 4) 收敛执行规则（建议）

1. **以真实接口为准**：联调阶段禁止再新增“规划路径”调用。
2. **文档层做别名**：规划文档保留业务可读性，但必须附“真实路径”。
3. **新增接口必须双更**：代码 + 本映射文档同步更新。
4. **模块负责人机制**：每个模块固定 1 名接口对口人，避免多人口径冲突。

---

## 5) 当前已处理的实装项

- 已补齐排班页面与路由，修复工作台“排班设置”入口死链：
  - `app/src/composables/usePlayerSchedule.ts`
  - `app/src/pages/player/schedule/index.vue`
  - `app/src/pages.json`
