# GameLink Client 端开发任务提示词

> 用于派发 AI 助手完成 Client PWA 前端开发任务

---

## 📊 现状分析（2025-01-14 更新）

### 已完成功能 ✅

| 模块 | 页面/组件 | 状态 | 说明 |
|------|-----------|------|------|
| **认证** | login-page.tsx | ✅ 完成 | 登录/注册切换，已连接真实 API |
| **首页** | home-page.tsx | ✅ 完成 | Hero、游戏分类、推荐陪玩师 |
| **陪玩师列表** | player-list-page.tsx | ✅ 完成 | 筛选、排序、卡片列表，已连接 API |
| **订单列表** | order-list-page.tsx | ✅ 完成 | Tab 筛选、订单卡片，已连接 API |
| **聊天列表** | chat-list-page.tsx | ✅ 完成 | 会话列表、未读数，已连接 API |
| **聊天详情** | chat-room-page.tsx | ✅ 完成 | 消息列表、发送消息 |
| **个人中心** | profile-page.tsx | ✅ 完成 | 用户信息、钱包、VIP、订单统计 |
| **状态管理** | 所有 stores | ✅ 完成 | auth/player/order/wallet/vip/chat/favorite/notification/theme |
| **HTTP 客户端** | http.ts | ✅ 完成 | Axios 封装、Token 拦截、响应解包 |
| **国际化** | i18n | ✅ 完成 | zh-CN/en-US 双语支持 |
| **主题** | theme-store | ✅ 完成 | 深色/浅色切换 |

### 待完善功能 🔧

| 模块 | 问题 | 优先级 |
|------|------|--------|
| **陪玩师详情** | 缺少独立详情页，点击直接跳转聊天 | P1 |
| **下单流程** | 缺少完整下单页面（选服务、选时长、支付） | P1 |
| **订单详情** | 缺少订单详情页 | P1 |
| **钱包页面** | 缺少独立钱包页面（充值、提现、交易记录） | P2 |
| **VIP 页面** | 缺少 VIP 详情/升级页面 | P2 |
| **个人资料编辑** | 缺少编辑页面（修改头像、昵称） | P2 |
| **修改密码** | 缺少修改密码页面 | P2 |
| **收藏功能** | store 存在但页面未实现 | P3 |
| **通知中心** | store 存在但页面未实现 | P3 |

### PWA 功能 ⬜ 未实现

| 功能 | 状态 |
|------|------|
| Web App Manifest | ⬜ 未配置 |
| Service Worker | ⬜ 未实现 |
| 离线支持 | ⬜ 未实现 |
| 安装提示 | ⬜ 未实现 |
| Push 通知 | ⬜ 未实现 |

---

## 🎯 开发任务优先级

### P1 - 核心业务流程（必须完成）

1. **陪玩师详情页** - 用户查看陪玩师完整信息
2. **下单流程** - 选择服务 → 确认订单 → 支付
3. **订单详情页** - 查看订单状态、操作（取消/评价）

### P2 - 用户体验增强

4. **钱包页面** - 充值、提现、交易记录
5. **VIP 页面** - VIP 等级、权益、升级
6. **个人资料编辑** - 修改头像、昵称
7. **修改密码** - 安全设置

### P3 - 功能完善

8. **收藏列表** - 收藏的陪玩师
9. **通知中心** - 系统通知、订单通知
10. **PWA 功能** - 离线支持、安装提示

---

## 项目背景（每次任务必带）

```
你正在开发 GameLink 游戏陪玩平台的 Client 端（桌面 PWA）。

技术栈：
- React 19 + TypeScript + Vite
- UI 组件库：shadcn/ui + Tailwind CSS
- 状态管理：Zustand
- 路由：React Router
- HTTP：Axios（已封装在 client/src/lib/http.ts）
- 国际化：i18next
- 图标：Lucide React

项目位置：client/
后端 API：http://localhost:8080/api/v1

代码规范：
- 组件使用 PascalCase 命名
- 文件使用 kebab-case
- 使用 TypeScript 严格模式
- 遵循 shadcn/ui 组件模式
- 支持深色/浅色主题切换
- 金额单位是分（cents），显示时除以 100

现有代码结构：
- client/src/pages/ - 页面组件
- client/src/components/ - 公共组件
- client/src/stores/modules/ - Zustand stores
- client/src/lib/http.ts - HTTP 客户端（已封装响应解包）
- client/src/locales/ - 国际化文件
```

---

## 任务模板

### 任务 1：陪玩师详情页（P1）

```
任务：实现陪玩师详情页

当前状态：
- player-store.ts 已有 fetchPlayerById 方法
- 点击陪玩师卡片目前直接跳转聊天页

需要完成：
1. 创建详情页 (client/src/pages/player/player-detail-page.tsx)
   - 路由：/players/:id
   - 陪玩师基本信息（头像、昵称、简介、在线状态）
   - 游戏标签、评分、接单数
   - 服务项目列表（游戏、段位、价格）
   - 评价列表（评分、内容、时间）
   - 底部固定：价格 + 立即预约按钮

2. 服务项目组件 (client/src/components/player/service-item-list.tsx)
   - 显示可选服务项目
   - 选择数量/时长
   - 实时计算总价

3. 更新路由 (client/src/App.tsx)
   - 添加 /players/:id 路由

4. 更新陪玩师卡片点击行为
   - 从跳转聊天改为跳转详情页

API 接口：
- GET /api/v1/public/players/:id - 陪玩师详情
- GET /api/v1/public/service-items?playerId=xxx - 服务项目列表

参考现有代码：
- client/src/pages/player/player-list-page.tsx
- client/src/stores/modules/player-store.ts
```

---

### 任务 2：下单流程（P1）

```
任务：实现完整下单流程

当前状态：
- order-store.ts 已有 createOrder 方法
- 缺少下单页面和支付页面

需要完成：
1. 下单确认页 (client/src/pages/order/create-order-page.tsx)
   - 路由：/orders/create?playerId=xxx&serviceItemId=xxx
   - 显示陪玩师信息
   - 显示服务项目
   - 选择数量/时长
   - 填写备注（可选）
   - 价格明细（单价 × 数量 = 总价）
   - 提交订单按钮

2. 支付页面 (client/src/pages/order/payment-page.tsx)
   - 路由：/orders/:id/pay
   - 显示订单信息
   - 选择支付方式（钱包余额/微信/支付宝）
   - 钱包余额不足提示
   - 确认支付按钮
   - 支付成功跳转订单详情

3. 更新 order-store.ts
   - 添加 payOrder 方法

API 接口：
- POST /api/v1/user/orders - 创建订单
- POST /api/v1/user/orders/:id/pay - 支付订单
- GET /api/v1/user/wallet/balance - 获取钱包余额

参考现有代码：
- client/src/stores/modules/order-store.ts
- client/src/stores/modules/wallet-store.ts
```

---

### 任务 3：订单详情页（P1）

```
任务：实现订单详情页

当前状态：
- order-store.ts 已有 fetchOrderById 方法
- order-list-page.tsx 已有订单列表

需要完成：
1. 订单详情页 (client/src/pages/order/order-detail-page.tsx)
   - 路由：/orders/:id
   - 订单状态（带状态图标和颜色）
   - 订单时间线（创建 → 支付 → 确认 → 完成）
   - 陪玩师信息卡片
   - 服务项目信息
   - 价格明细
   - 操作按钮：
     - 待支付：去支付、取消订单
     - 进行中：联系陪玩师、申请退款
     - 已完成：评价、再次下单
     - 已取消/已退款：再次下单

2. 订单状态时间线组件 (client/src/components/order/order-timeline.tsx)

3. 更新订单列表卡片
   - 点击跳转详情页

API 接口：
- GET /api/v1/user/orders/:id - 订单详情
- POST /api/v1/user/orders/:id/cancel - 取消订单
- POST /api/v1/user/orders/:id/review - 提交评价

参考现有代码：
- client/src/pages/order/order-list-page.tsx
- client/src/stores/modules/order-store.ts
```

---

### 任务 4：钱包页面（P2）

```
任务：实现钱包页面

当前状态：
- wallet-store.ts 已完成，包含 fetchWallet、fetchTransactions、recharge、withdraw
- profile-page.tsx 显示余额但无详情页

需要完成：
1. 钱包页面 (client/src/pages/wallet/wallet-page.tsx)
   - 路由：/wallet
   - 余额卡片（可用余额、冻结余额）
   - 快捷操作（充值、提现）
   - 交易记录列表（Tab：全部/充值/支付/退款/提现）
   - 交易记录卡片（类型图标、金额、状态、时间）

2. 充值页面 (client/src/pages/wallet/recharge-page.tsx)
   - 路由：/wallet/recharge
   - 充值档位选择
   - 自定义金额输入
   - 支付方式选择
   - 确认充值

3. 提现页面 (client/src/pages/wallet/withdraw-page.tsx)
   - 路由：/wallet/withdraw
   - 可提现余额
   - 提现金额输入
   - 银行卡选择
   - 确认提现

API 接口：
- GET /api/v1/user/wallet/balance - 余额
- GET /api/v1/user/wallet/transactions - 交易记录
- POST /api/v1/user/wallet/recharge - 充值
- POST /api/v1/user/wallet/withdraw - 提现
- GET /api/v1/user/recharge/options - 充值档位

参考现有代码：
- client/src/stores/modules/wallet-store.ts
- client/src/pages/profile/profile-page.tsx
```

---

### 任务 5：VIP 页面（P2）

```
任务：实现 VIP 页面

当前状态：
- vip-store.ts 已完成，包含 fetchVipInfo
- profile-page.tsx 显示 VIP 等级但无详情页

需要完成：
1. VIP 页面 (client/src/pages/vip/vip-page.tsx)
   - 路由：/vip
   - 当前等级卡片（等级图标、名称、经验进度条）
   - 等级权益列表
   - 所有等级展示（当前等级高亮）
   - 升级说明（消费/充值门槛）

2. VIP 等级卡片组件 (client/src/components/vip/vip-level-card.tsx)

API 接口：
- GET /api/v1/user/vip/info - VIP 信息
- GET /api/v1/user/vip/levels - 等级列表
- GET /api/v1/user/vip/threshold - 解锁门槛

参考现有代码：
- client/src/stores/modules/vip-store.ts
- client/src/pages/profile/profile-page.tsx
```

---

### 任务 6：个人资料编辑（P2）

```
任务：实现个人资料编辑页面

当前状态：
- auth-store.ts 已有 updateProfile 方法
- profile-page.tsx 有编辑按钮但无功能

需要完成：
1. 编辑资料页 (client/src/pages/profile/edit-profile-page.tsx)
   - 路由：/profile/edit
   - 头像上传（点击更换）
   - 昵称输入
   - 保存按钮

2. 头像上传组件 (client/src/components/profile/avatar-upload.tsx)
   - 点击选择图片
   - 图片预览
   - 上传进度

API 接口：
- PUT /api/v1/user/profile - 更新资料
- POST /api/v1/user/upload/avatar - 上传头像（如有）

参考现有代码：
- client/src/stores/modules/auth-store.ts
- client/src/pages/profile/profile-page.tsx
```

---

### 任务 7：修改密码（P2）

```
任务：实现修改密码页面

当前状态：
- 后端已有 POST /api/v1/auth/change-password 接口
- profile-page.tsx 安全 Tab 有入口但无功能

需要完成：
1. 修改密码页 (client/src/pages/profile/change-password-page.tsx)
   - 路由：/profile/change-password
   - 当前密码输入
   - 新密码输入
   - 确认新密码
   - 密码强度提示
   - 提交按钮

API 接口：
- POST /api/v1/auth/change-password
  - 请求：{ oldPassword, newPassword }
  - 密码要求：8+ 字符，包含大小写、数字、特殊符号

参考现有代码：
- client/src/pages/auth/login-page.tsx（表单样式）
```

---

### 任务 8：PWA 功能（P3）

```
任务：实现 PWA 功能

当前状态：
- 无 PWA 配置

需要完成：
1. Web App Manifest (client/public/manifest.json)
   - name: "GameLink - 游戏陪玩平台"
   - short_name: "GameLink"
   - 图标：192x192, 512x512
   - theme_color: "#5865F2"
   - background_color: "#313338"
   - display: "standalone"

2. Service Worker (使用 vite-plugin-pwa)
   - 安装配置 vite-plugin-pwa
   - 静态资源缓存
   - API 缓存策略（stale-while-revalidate）

3. 安装提示组件 (client/src/components/pwa/install-prompt.tsx)
   - 检测 beforeinstallprompt 事件
   - 自定义安装提示 UI
   - 安装按钮

4. 离线状态提示 (client/src/components/pwa/offline-banner.tsx)
   - 检测网络状态
   - 离线时显示提示条

参考：.kiro/CLIENT_WEB_PWA_CHECKLIST.md
```

---

## 通用注意事项

```
开发时请注意：

1. 主题支持
   - 所有组件必须支持深色/浅色主题
   - 使用 Tailwind dark: 前缀或 CSS 变量

2. 响应式设计
   - 桌面端优先，适配平板
   - 最小宽度 320px

3. 错误处理
   - API 错误统一处理（http.ts 已配置 401 自动登出）
   - 显示友好的错误提示（使用 toast）

4. 加载状态
   - 使用 Skeleton 组件
   - 按钮 loading 状态

5. 国际化
   - 所有文案使用 i18next
   - 新增文案需同时更新 zh-CN 和 en-US

6. 类型安全
   - 定义完整的 TypeScript 类型
   - 复用 stores 中已定义的类型

7. 金额处理
   - 后端金额单位是分（cents）
   - 显示时除以 100，保留 2 位小数
   - 使用：(cents / 100).toFixed(2)
```

---

## API 响应格式

```typescript
// 统一响应格式（http.ts 已自动解包 data）
interface ApiResponse<T> {
  success: boolean;
  code: number;
  message: string;
  data: T;  // http.ts 返回的就是这个 data
}

// 分页响应
interface PaginatedResponse<T> {
  items: T[];
  pagination: {
    page: number;
    pageSize: number;
    total: number;
    totalPages: number;
  };
}
```

---

## 快速开始命令

```bash
# 启动开发服务器
cd client && npm run dev

# 构建生产版本
cd client && npm run build

# 类型检查
cd client && npm run type-check

# 代码检查
cd client && npm run lint
```
