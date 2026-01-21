# GameLink 微信小程序开发计划

> **版本**: v1.0 | **日期**: 2026-01-21 | **状态**: 进行中
>
> 原生微信小程序技术栈，基于 Discord Dark Theme 设计规范

---

## 📑 目录

1. [技术架构](#1-技术架构)
2. [当前实现状态](#2-当前实现状态)
3. [开发里程碑](#3-开发里程碑)
4. [页面与组件清单](#4-页面与组件清单)
5. [API 对接方案](#5-api-对接方案)
6. [开发规范](#6-开发规范)
7. [测试与发布](#7-测试与发布)

---

## 1. 技术架构

### 1.1 技术栈

| 层次 | 技术 | 版本 | 说明 |
|------|------|------|------|
| **框架** | 原生微信小程序 | - | 性能最优 |
| **语言** | TypeScript | 5.x | 类型安全 |
| **样式** | Less | 4.x | CSS 预处理 |
| **渲染** | Skyline | - | 高性能渲染引擎 |
| **组件框架** | glass-easel | - | 组件化开发 |
| **状态管理** | 全局数据 + Storage | - | 轻量级状态管理 |
| **HTTP** | wx.request 封装 | - | API 请求 |
| **加密** | crypto-js | 4.2+ | AES-256-CBC |
| **WebSocket** | wx.connectSocket | - | 实时通讯 |

### 1.2 目录结构

```
app/miniprogram/
├── app.json                  # 小程序配置
├── app.ts                    # 小程序逻辑
├── app.less                  # 全局样式
├── theme.json                # 主题配置
├── sitemap.json              # 站点地图
├── pages/                    # 页面
│   ├── index/                # 首页（陪玩师列表）
│   ├── category/             # 游戏分类
│   ├── message/              # 消息列表
│   ├── profile/              # 个人中心
│   ├── player/               # 陪玩师详情
│   ├── login/                # 登录页
│   └── order/                # 订单相关
│       ├── create/           # 创建订单
│       ├── detail/           # 订单详情
│       └── list/             # 订单列表
├── components/               # 组件库
│   ├── gl-page/              # 页面容器
│   ├── gl-button/            # 按钮
│   ├── gl-card/              # 卡片
│   ├── gl-avatar/            # 头像
│   ├── gl-tag/               # 标签
│   ├── gl-navbar/            # 导航栏
│   ├── gl-search/            # 搜索框
│   ├── gl-loading/           # 加载状态
│   ├── gl-empty/             # 空状态
│   ├── gl-section/           # 区块标题
│   ├── gl-icon/              # 图标
│   ├── player-card/          # 陪玩师卡片
│   ├── game-card/            # 游戏卡片
│   ├── message-item/         # 消息项
│   ├── dev-toolbar/          # 开发工具栏
│   ├── user-menu/            # 用户菜单
│   ├── user-profile-header/  # 用户资料头部
│   └── user-stats/           # 用户统计
├── custom-tab-bar/           # 自定义 TabBar
├── utils/                    # 工具函数
│   ├── request.ts            # API 请求封装
│   ├── auth.ts               # 认证工具
│   ├── storage.ts            # 存储工具
│   └── theme.ts              # 主题切换
├── styles/                   # 全局样式
│   └── variables.less        # Discord 风格变量
├── config/                   # 配置文件
│   └── index.ts              # 配置入口
├── typings/                  # 类型定义
│   └── index.d.ts            # 全局类型
└── assets/                   # 静态资源
    └── icons/                # SVG 图标
```

### 1.3 已实现的基础设施

| 模块 | 状态 | 说明 |
|------|:---:|:------|
| **请求封装** | ✅ | Token 管理、自动刷新、错误处理 |
| **存储工具** | ✅ | StorageKeys、get/set/remove |
| **主题系统** | ✅ | 用户/陪玩师模式切换 |
| **基础组件** | ✅ | gl-page、gl-button、gl-card 等 |
| **自定义 TabBar** | ✅ | 支持 SVG 图标，双模式切换 |
| **Discord 样式** | ✅ | variables.less 变量定义 |

---

## 2. 当前实现状态

### 2.1 已实现页面

| 页面 | 路径 | 状态 | 完成度 |
|:---|:---|:---:|:---:|
| 首页 | pages/index/ | ⚠️ 框架 | 30% |
| 游戏分类 | pages/category/ | ⚠️ 框架 | 20% |
| 消息列表 | pages/message/ | ⚠️ 框架 | 20% |
| 个人中心 | pages/profile/ | ⚠️ 框架 | 30% |
| 陪玩师详情 | pages/player/ | ⚠️ 框架 | 20% |
| 登录 | pages/login/ | ⚠️ 框架 | 40% |
| 创建订单 | pages/order/create/ | ❌ 未开始 | 0% |
| 订单详情 | pages/order/detail/ | ❌ 未开始 | 0% |
| 订单列表 | pages/order/list/ | ❌ 未开始 | 0% |

### 2.2 已实现组件

| 组件 | 路径 | 状态 |
|:---|:---:|:---:|
| 页面容器 | components/gl-page/ | ✅ |
| 按钮 | components/gl-button/ | ✅ |
| 卡片 | components/gl-card/ | ✅ |
| 头像 | components/gl-avatar/ | ✅ |
| 标签 | components/gl-tag/ | ✅ |
| 导航栏 | components/gl-navbar/ | ✅ |
| 搜索框 | components/gl-search/ | ✅ |
| 加载 | components/gl-loading/ | ✅ |
| 空状态 | components/gl-empty/ | ✅ |
| 区块标题 | components/gl-section/ | ✅ |
| 图标 | components/gl-icon/ | ✅ |
| 陪玩师卡片 | components/player-card/ | ✅ |
| 游戏卡片 | components/game-card/ | ✅ |
| 消息项 | components/message-item/ | ✅ |

### 2.3 已实现工具函数

| 模块 | 文件 | 功能 |
|:---|:---|:---|
| 请求封装 | utils/request.ts | HTTP 请求、Token 刷新、错误处理 |
| 认证工具 | utils/auth.ts | 登录状态、Token 管理 |
| 存储工具 | utils/storage.ts | StorageKeys、get/set/remove |
| 主题切换 | utils/theme.ts | 用户/陪玩师模式切换 |

---

## 3. 开发里程碑

### Phase 1: 基础功能完善（4 周）

#### Week 1: 认证与用户状态

- [x] ~~登录页面框架~~
- [ ] 手机号 + 验证码登录
- [ ] JWT Token 管理
- [ ] 自动登录
- [ ] 用户模式切换（用户 ↔ 陪玩师）
- [ ] 主题样式切换

#### Week 2: 首页与游戏分类

**首页**
- [ ] 游戏分类入口
- [ ] 陪玩师推荐列表（智能排序）
- [ ] 筛选器（游戏/段位/价格）
- [ ] 排序方式（评分/价格/接单量）
- [ ] 无限滚动加载
- [ ] 下拉刷新

**游戏分类**
- [ ] 游戏列表展示
- [ ] 游戏卡片（在线人数）
- [ ] 进入游戏详情

#### Week 3: 陪玩师详情与下单

**陪玩师详情**
- [ ] 基本信息展示
- [ ] 游戏段位展示
- [ ] 服务项目列表
- [ ] 数据统计（评分/接单量/好评率）
- [ ] 用户评价列表
- [ ] 收藏功能

**创建订单**
- [ ] 服务类型选择（陪玩/陪练/教学）
- [ ] 游戏/段位选择
- [ ] 时长选择（1小时/2小时/包夜）
- [ ] 特殊需求备注
- [ ] 订单金额确认
- [ ] 优惠券抵扣
- [ ] 支付方式选择（余额/微信）

#### Week 4: 订单管理

**订单列表**
- [ ] Tab 切换（待支付/待接单/进行中/待评价/已完成）
- [ ] 订单卡片展示
- [ ] 订单状态筛选

**订单详情**
- [ ] 订单状态进度条
- [ ] 陪玩师信息
- [ ] 服务信息
- [ ] 服务倒计时
- [ ] 聊天室入口
- [ ] 确认完成按钮

### Phase 2: 陪玩师功能（3 周）

#### Week 5: 陪玩师工作台

- [ ] 接单开关（ON/OFF）
- [ ] 在线状态选择（在线/忙碌/离线）
- [ ] 今日统计展示
- [ ] 快捷入口

#### Week 6: 抢单池与接单

**抢单池**
- [ ] 筛选器（游戏/段位/最低价格）
- [ ] 实时订单列表（WebSocket）
- [ ] 订单卡片（收入透明展示）
- [ ] 接单操作
- [ ] 震动提醒
- [ ] 空状态

**我的订单（陪玩师）**
- [ ] Tab 切换（进行中/待完成/待评价/已完成）
- [ ] 订单详情（收入明细）
- [ ] 完成服务按钮

#### Week 7: 收益与提现

**收益中心**
- [ ] 收入总览（今日/本周/本月/累计）
- [ ] 收入趋势图
- [ ] 资产详情（可提现/冻结/T+7）
- [ ] 提现申请

### Phase 3: 聊天与评价（2 周）

#### Week 8: 聊天室

- [ ] WebSocket 连接管理
- [ ] 消息列表展示
- [ ] 系统消息（订单状态变更）
- [ ] 用户消息
- [ ] 输入框（仅文字）
- [ ] 消息推送

#### Week 9: 评价系统

- [ ] 提交评价（星级/标签/文字）
- [ ] 修改评价（24小时内）
- [ ] 评价列表（筛选）
- [ ] 标签统计

### Phase 4: 高级功能（3 周）

#### Week 10: 钱包与充值

- [ ] 余额展示
- [ ] 充值接口
- [ ] 充值记录
- [ ] 交易记录

#### Week 11: 排行榜

- [ ] 评分排行
- [ ] 接单量排行
- [ ] 收入排行

#### Week 12: 个人中心完善

**用户中心**
- [ ] 我的优惠券
- [ ] 收藏的陪玩师
- [ ] 设置

**陪玩师中心**
- [ ] 认证流程（实名/段位）
- [ ] 服务说明编辑
- [ ] 接单时间设置

### Phase 5: 优化与发布（2 周）

#### Week 13: 性能优化

- [ ] 图片懒加载
- [ ] 列表虚拟滚动
- [ ] 缓存策略
- [ ] 骨架屏
- [ ] 加载状态

#### Week 14: 测试与发布

- [ ] 端到端流程测试
- [ ] 异常处理
- [ ] 微信审核准备
- [ ] 发布上线

---

## 4. 页面与组件清单

### 4.1 页面列表

#### 用户模式页面

| 页面 | 路径 | 优先级 | 依赖 |
|:---|:---|:---:|:---|
| 首页 | pages/index/ | P0 | - |
| 游戏分类 | pages/category/ | P1 | - |
| 陪玩师详情 | pages/player/ | P0 | 首页 |
| 创建订单 | pages/order/create/ | P0 | 陪玩师详情 |
| 订单列表 | pages/order/list/ | P0 | - |
| 订单详情 | pages/order/detail/ | P0 | 订单列表 |
| 聊天室 | pages/chat/ | P1 | 订单详情 |
| 个人中心 | pages/profile/ | P1 | - |
| 钱包 | pages/wallet/ | P2 | - |
| 优惠券 | pages/coupon/ | P2 | - |
| 排行榜 | pages/ranking/ | P2 | - |

#### 陪玩师模式页面

| 页面 | 路径 | 优先级 | 依赖 |
|:---|:---|:---:|:---|
| 工作台 | pages/workbench/ | P0 | - |
| 抢单池 | pages/order-pool/ | P0 | - |
| 我的订单 | pages/player-orders/ | P0 | - |
| 收益中心 | pages/earnings/ | P1 | - |
| 提现申请 | pages/withdraw/ | P1 | 收益中心 |
| 数据分析 | pages/analytics/ | P2 | - |
| 认证流程 | pages/certification/ | P1 | - |

#### 通用页面

| 页面 | 路径 | 优先级 |
|:---|:---|:---:|
| 登录 | pages/login/ | P0 |
| 注册 | pages/register/ | P0 |
| 设置 | pages/settings/ | P2 |

### 4.2 组件清单

#### 基础组件 (gl-*)

| 组件 | 状态 | 说明 |
|:---|:---:|:---|
| gl-page | ✅ | 页面容器（状态栏/TabBar 占位） |
| gl-button | ✅ | 按钮（primary/secondary/outline/ghost） |
| gl-card | ✅ | 卡片容器 |
| gl-avatar | ✅ | 头像（带在线状态） |
| gl-tag | ✅ | 标签 |
| gl-navbar | ✅ | 导航栏 |
| gl-search | ✅ | 搜索框 |
| gl-loading | ✅ | 加载状态 |
| gl-empty | ✅ | 空状态 |
| gl-section | ✅ | 区块标题 |
| gl-icon | ✅ | 图标 |

#### 业务组件

| 组件 | 状态 | 说明 |
|:---|:---|:---|
| player-card | ✅ | 陪玩师卡片 |
| game-card | ✅ | 游戏分类卡片 |
| message-item | ✅ | 消息项 |
| order-card | ❌ | 订单卡片 |
| review-item | ❌ | 评价项 |
| earnings-card | ❌ | 收益卡片 |
| ranking-item | ❌ | 排行榜项 |

---

## 5. API 对接方案

### 5.1 已有后端 API

后端 API 已 100% 完成，详见业务流程文档：
- [09-auth-flow-guide.md](../../.kiro/steering/client/09-auth-flow-guide.md)
- [10-player-flow-guide.md](../../.kiro/steering/client/10-player-flow-guide.md)
- [11-payment-wallet-guide.md](../../.kiro/steering/client/11-payment-wallet-guide.md)

### 5.2 API 模块化

建议在 `utils/` 下创建 API 模块：

```typescript
// utils/api/user.ts
export const userAPI = {
  // 陪玩师
  getPlayers: (params: PlayerListParams) =>
    http.get('/players', params),

  getPlayerDetail: (id: number) =>
    http.get(`/players/${id}`),

  // 订单
  createOrder: (data: CreateOrderDTO) =>
    http.post('/orders', data),

  getMyOrders: (params: OrderQueryParams) =>
    http.get('/user/orders', params),

  // 支付
  createPayment: (data: CreatePaymentDTO) =>
    http.post('/payments', data),
};

// utils/api/player.ts
export const playerAPI = {
  // 抢单
  getOrderPool: () =>
    http.get('/orders/available'),

  acceptOrder: (orderId: string) =>
    http.post(`/orders/${orderId}/accept`),

  // 收益
  getEarnings: () =>
    http.get('/player/earnings'),

  withdraw: (data: WithdrawDTO) =>
    http.post('/withdraws', data),
};
```

### 5.3 WebSocket 实现

```typescript
// utils/websocket.ts
class WebSocketManager {
  private ws: WechatMiniprogram.SocketTask | null = null

  connect(url: string) {
    this.ws = wx.connectSocket({
      url,
      header: {
        'Authorization': `Bearer ${getToken()}`
      }
    })

    wx.onSocketMessage((res) => {
      const message = JSON.parse(res.data as string)
      // 处理消息
    })
  }
}

export const wsManager = new WebSocketManager()
```

---

## 6. 开发规范

### 6.1 命名规范

| 类型 | 规范 | 示例 |
|:---|:---|:---|
| 页面目录 | kebab-case | `pages/player-detail/` |
| 组件目录 | kebab-case | `components/gl-button/` |
| TS 文件 | camelCase | `request.ts` |
| Less 文件 | kebab-case | `variables.less` |
| 类名 | kebab-case + BEM | `.player-card__name` |
| 变量 | camelCase | `userInfo` |
| 常量 | UPPER_SNAKE_CASE | `API_BASE_URL` |

### 6.2 样式规范

```less
// 使用 Discord 风格变量
@import './variables.less';

.player-card {
  background-color: @bg-secondary;
  border-radius: @radius-lg;
  padding: @card-padding;

  &__name {
    color: @text-header;
    font-size: @font-size-md;
  }

  &--online {
    .status-dot(@color-success);
  }
}
```

### 6.3 组件开发规范

```typescript
// 组件 TypeScript 文件
Component({
  properties: {
    title: String,
    data: Object,
  },

  methods: {
    onTap() {
      this.triggerEvent('tap', { data: this.data })
    }
  }
})
```

---

## 7. 测试与发布

### 7.1 测试计划

| 测试类型 | 工具 | 覆盖范围 |
|:---|:---|:---|
| 单元测试 | Jest | 工具函数 |
| 组件测试 | 微信开发者工具 | 组件渲染 |
| 集成测试 | 微信开发者工具 | 端到端流程 |

### 7.2 发布清单

- [ ] 代码审查通过
- [ ] 所有测试通过
- [ ] 性能优化完成
- [ ] 无障碍访问（可选）
- [ ] 微信审核准备
  - [ ] 类目选择
  - [ ] 功能描述
  - [ ] 测试账号
- [ ] 隐私协议更新
- [ ] 用户协议更新

---

## 📌 附录

### A. 已完成的文档

- [miniprogram-design.md](../../.kiro/steering/client/miniprogram-design.md) - UI/UX 设计规范
- [apps-roadmap.md](../../.kiro/steering/client/apps-roadmap.md) - 产品规划

### B. 业务流程文档

- [09-auth-flow-guide.md](../../.kiro/steering/client/09-auth-flow-guide.md)
- [10-player-flow-guide.md](../../.kiro/steering/client/10-player-flow-guide.md)
- [11-payment-wallet-guide.md](../../.kiro/steering/client/11-payment-wallet-guide.md)
- [12-vip-marketing-guide.md](../../.kiro/steering/client/12-vip-marketing-guide.md)
- [13-chat-review-dispute-guide.md](../../.kiro/steering/client/13-chat-review-dispute-guide.md)

### C. 参考资料

- [微信小程序官方文档](https://developers.weixin.qq.com/miniprogram/dev/framework/)
- [Skyline 渲染引擎](https://developers.weixin.qq.com/miniprogram/dev/framework/skyline/)
- [glass-easel 组件框架](https://github.com/wechat-miniprogram/glass-easel-component-framework)

---

## ✍️ 文档变更记录

| 版本 | 日期 | 作者 | 变更说明 |
|------|------|------|----------|
| v1.0 | 2026-01-21 | Super Dev | 初始版本 |

---

**文档状态**: ✅ 已完成
**下一步**: 开始 Phase 1 - Week 1 开发
