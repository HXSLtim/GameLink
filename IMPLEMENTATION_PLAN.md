# 🎮 GameLink 功能实装计划
## 参考设计：KOOK / Discord 风格

**创建日期**: 2025-11-16
**目标**: 实装用户端和陪玩师端，参考KOOK/Discord设计风格

---

## 📊 现状分析

### ✅ 已完成部分

**后端 (100%)**:
- ✅ 75+ API接口全部实现
- ✅ 用户端接口 (12个)
- ✅ 陪玩师端接口 (12个)
- ✅ 管理端接口 (30+个)
- ✅ 实时聊天接口
- ✅ 支付、订单、评价等完整业务流程

**前端管理端 (100%)**:
- ✅ 17个管理端页面
- ✅ 用户管理、订单管理、游戏管理
- ✅ 数据统计、权限管理
- ✅ 95/100代码质量

### ❌ 缺失部分

**用户端 (0%)**:
- ❌ 用户主界面
- ❌ 陪玩师浏览页面
- ❌ 订单创建/管理
- ❌ 支付界面
- ❌ 聊天界面
- ❌ 个人中心

**陪玩师端 (0%)**:
- ❌ 陪玩师工作台
- ❌ 接单界面
- ❌ 订单管理
- ❌ 收益查看
- ❌ 个人资料设置

---

## 🎨 设计参考：KOOK / Discord

### 核心设计理念

#### 1. **三栏布局** (Discord/KOOK经典布局)

```
┌────────────────────────────────────────────────────────┐
│  [Logo]   │        主要内容区域         │    侧边栏     │
│           │                            │              │
│  导航区    │    - 陪玩师列表            │   - 筛选     │
│           │    - 订单管理              │   - 详情     │
│  - 首页    │    - 聊天窗口              │   - 推荐     │
│  - 陪玩师  │    - 个人中心              │              │
│  - 订单    │                            │              │
│  - 消息    │                            │              │
│  - 设置    │                            │              │
│           │                            │              │
└────────────────────────────────────────────────────────┘
```

#### 2. **色彩系统** (参考Discord)

- **深色主题** (默认)：
  - 背景: `#36393f` (深灰)
  - 卡片: `#2f3136` (更深)
  - 强调色: `#5865f2` (Discord蓝) / `#6DD230` (KOOK绿)
  - 文字: `#dcddde` (浅灰)

- **浅色主题**：
  - 背景: `#ffffff`
  - 卡片: `#f2f3f5`
  - 强调色: 同深色主题
  - 文字: `#2e3338`

#### 3. **组件风格**

**卡片设计**：
```tsx
// 陪玩师卡片 - Discord风格
<Card>
  <Avatar size="large" />
  <PlayerInfo>
    <Name>玩家昵称</Name>
    <Status online>在线</Status>
    <Tags>
      <Tag>王者荣耀</Tag>
      <Tag>陪玩上分</Tag>
    </Tags>
    <Rating>⭐ 4.9 (156评价)</Rating>
  </PlayerInfo>
  <ActionButtons>
    <Button type="primary">立即下单</Button>
    <Button>查看详情</Button>
  </ActionButtons>
</Card>
```

**实时状态指示器**：
- 🟢 在线
- 🟡 忙碌
- ⚫ 离线
- 🔴 游戏中

---

## 🚀 实装计划

### Phase 1: 基础架构 (1-2天)

#### 1.1 项目结构调整

创建三端独立项目：

```
GameLink/
├── frontend-admin/     # 管理端 (已有)
├── frontend-user/      # 用户端 (新建)
└── frontend-player/    # 陪玩师端 (新建)
```

或使用单体应用 + 路由区分：

```
frontend/
├── src/
│   ├── apps/
│   │   ├── admin/      # 管理端
│   │   ├── user/       # 用户端
│   │   └── player/     # 陪玩师端
│   ├── shared/         # 共享组件
│   └── layouts/        # 布局组件
```

**推荐方案**: **单体应用** (便于共享组件和状态)

#### 1.2 WebSocket实时通讯

```typescript
// services/websocket.ts
class WebSocketService {
  private ws: WebSocket | null = null;

  connect(userId: string) {
    this.ws = new WebSocket(`ws://localhost:8080/ws?userId=${userId}`);
    this.setupListeners();
  }

  // 订单状态更新
  onOrderUpdate(callback: (order: Order) => void) {
    this.on('order_update', callback);
  }

  // 新消息通知
  onNewMessage(callback: (message: Message) => void) {
    this.on('new_message', callback);
  }
}
```

#### 1.3 共享组件库

参考Discord的组件：
- `<DiscordLayout />` - 三栏布局
- `<PlayerCard />` - 陪玩师卡片
- `<ChatWindow />` - 聊天窗口
- `<OrderCard />` - 订单卡片
- `<StatusBadge />` - 状态徽章
- `<Avatar />` - 头像组件

---

### Phase 2: 用户端实装 (3-4天)

#### 2.1 核心页面

##### **首页 (Discover)**
```
功能：
- 精选陪玩师推荐
- 游戏分类导航
- 热门订单展示
- 实时活动公告

参考：Discord服务器发现页
```

##### **陪玩师浏览 (Players)**
```
功能：
- 筛选器 (游戏、价格、段位、在线状态)
- 陪玩师列表 (卡片或列表视图)
- 排序 (评分、价格、接单量)
- 实时在线状态

参考：KOOK的语音频道成员列表
```

##### **订单管理 (Orders)**
```
功能：
- 我的订单列表
- 订单详情查看
- 订单状态追踪
- 支付/取消/完成操作

参考：Discord的好友/消息列表
```

##### **聊天 (Messages)**
```
功能：
- 聊天列表 (左侧)
- 聊天窗口 (中间)
- 用户信息 (右侧)
- 实时消息通知

参考：Discord聊天界面 (几乎完全一样)
```

##### **个人中心 (Profile)**
```
功能：
- 个人信息编辑
- 订单历史
- 收藏的陪玩师
- 我的评价

参考：Discord用户设置
```

#### 2.2 路由配置

```typescript
// apps/user/routes.tsx
const userRoutes = [
  { path: '/user', element: <UserLayout />, children: [
    { path: 'discover', element: <DiscoverPage /> },
    { path: 'players', element: <PlayersPage /> },
    { path: 'players/:id', element: <PlayerDetailPage /> },
    { path: 'orders', element: <OrdersPage /> },
    { path: 'orders/create', element: <CreateOrderPage /> },
    { path: 'messages', element: <MessagesPage /> },
    { path: 'profile', element: <ProfilePage /> },
  ]},
];
```

---

### Phase 3: 陪玩师端实装 (3-4天)

#### 3.1 核心页面

##### **工作台 (Dashboard)**
```
功能：
- 今日数据概览 (接单量、收益、好评率)
- 待处理订单
- 收益趋势图
- 快捷操作

参考：Discord服务器管理面板
```

##### **订单管理 (Orders)**
```
功能：
- 可接订单列表
- 进行中的订单
- 历史订单
- 一键接单/完成

参考：KOOK的消息管理
```

##### **收益中心 (Earnings)**
```
功能：
- 收益统计
- 提现管理
- 交易记录
- 佣金明细

参考：独立设计 (图表展示)
```

##### **个人资料 (Profile)**
```
功能：
- 基本信息设置
- 游戏段位认证
- 服务价格设置
- 工作时间设置
- 个人展示页预览

参考：Discord的个人资料编辑
```

#### 3.2 实时功能

```typescript
// 新订单通知
useEffect(() => {
  wsService.onNewOrder((order) => {
    notification.info({
      message: '新订单',
      description: `${order.game} - ${order.user}`,
      btn: <Button onClick={() => acceptOrder(order)}>接单</Button>
    });
  });
}, []);
```

---

### Phase 4: UI组件库 (贯穿整个开发)

#### 4.1 Discord风格组件

```typescript
// components/discord/

// 布局组件
<DiscordLayout>
  <DiscordSidebar />
  <DiscordMain />
  <DiscordPanel />
</DiscordLayout>

// 消息组件
<MessageList>
  <Message
    avatar="..."
    username="玩家A"
    timestamp="10:30"
    content="你好，什么时候可以开始？"
  />
</MessageList>

// 状态组件
<StatusIndicator status="online" />
<ActivityStatus activity="正在游戏中：王者荣耀" />

// 卡片组件
<ServerCard
  name="王者荣耀陪玩"
  memberCount={1234}
  onlineCount={567}
/>
```

#### 4.2 主题系统

```typescript
// theme/discord.ts
export const discordTheme = {
  dark: {
    background: {
      primary: '#36393f',
      secondary: '#2f3136',
      tertiary: '#202225',
    },
    text: {
      primary: '#dcddde',
      secondary: '#b9bbbe',
      muted: '#72767d',
    },
    accent: {
      primary: '#5865f2',
      success: '#3ba55d',
      warning: '#faa61a',
      danger: '#ed4245',
    },
  },
  light: {
    // 浅色主题配置
  },
};
```

---

## 📱 移动端适配

参考Discord移动端：
- 底部导航栏
- 滑动抽屉
- 手势交互
- 响应式布局

---

## 🛠️ 技术栈

### 新增技术

| 技术 | 用途 | 参考 |
|------|------|------|
| **Socket.IO** | WebSocket实时通讯 | Discord实时消息 |
| **Framer Motion** | 动画效果 | Discord流畅动画 |
| **React Query** | 数据缓存和同步 | 实时数据更新 |
| **Zustand** | 全局状态管理 | 替代Context |
| **Day.js** | 时间处理 | 消息时间戳 |
| **React Virtualized** | 长列表优化 | 陪玩师列表 |

### 已有技术 (保持)

- React 18
- TypeScript
- Arco Design (基础组件)
- Vite
- React Router

---

## 📦 实施步骤

### Week 1: 基础架构

- [ ] Day 1-2: 创建用户端和陪玩师端项目结构
- [ ] Day 3-4: 实装WebSocket服务
- [ ] Day 5: 搭建Discord风格布局组件
- [ ] Day 6-7: 创建共享组件库

### Week 2: 用户端核心功能

- [ ] Day 1-2: 首页 + 陪玩师浏览
- [ ] Day 3-4: 订单管理 + 支付流程
- [ ] Day 5-6: 聊天功能
- [ ] Day 7: 个人中心

### Week 3: 陪玩师端核心功能

- [ ] Day 1-2: 工作台 + 订单管理
- [ ] Day 3-4: 收益中心
- [ ] Day 5-6: 个人资料设置
- [ ] Day 7: 测试和优化

### Week 4: 优化和发布

- [ ] Day 1-3: UI/UX优化
- [ ] Day 4-5: 性能优化
- [ ] Day 6: 移动端适配
- [ ] Day 7: 发布准备

---

## 🎯 优先级

### P0 (必须)
1. ✅ 用户端陪玩师浏览 + 下单
2. ✅ 陪玩师端接单 + 完成订单
3. ✅ 基础聊天功能
4. ✅ 支付流程

### P1 (重要)
1. 实时通知
2. 个人中心
3. 收益管理
4. 评价系统

### P2 (优化)
1. 深色/浅色主题切换
2. 移动端适配
3. 动画效果
4. 性能优化

---

## 📊 成功指标

### 用户体验
- ⏱️ 页面加载时间 < 2秒
- 📱 移动端完全适配
- 🎨 UI一致性 ≥ 95%
- ✨ 动画流畅度 ≥ 60fps

### 功能完整性
- ✅ 核心业务流程 100%
- ✅ API集成 100%
- ✅ 实时功能 ≥ 90%

### 代码质量
- 📝 TypeScript覆盖率 ≥ 95%
- 🧪 单元测试覆盖率 ≥ 80%
- 📚 组件文档 100%

---

## 🎨 设计资源

### 参考链接
- [Discord设计系统](https://discord.com/branding)
- [KOOK官网](https://www.kookapp.cn)
- [Figma社区Discord克隆](https://www.figma.com/community/file/discord-ui)

### 色板
```css
/* Discord经典色 */
--discord-blurple: #5865f2;
--discord-green: #3ba55d;
--discord-yellow: #faa61a;
--discord-red: #ed4245;

/* KOOK绿色系 */
--kook-primary: #6DD230;
--kook-secondary: #53C41A;
```

---

## 📝 开发规范

### 组件命名
```
User开头 - 用户端组件
Player开头 - 陪玩师端组件
Discord开头 - Discord风格组件
Shared开头 - 共享组件
```

### 文件结构
```
features/
  user/
    discover/
      DiscoverPage.tsx
      components/
        FeaturedPlayers.tsx
        GameCategories.tsx
      hooks/
        useDiscoverData.ts
      styles/
        discover.module.less
```

---

## 🚀 开始实施

**立即开始**:
1. 创建用户端和陪玩师端项目结构
2. 搭建Discord风格的布局组件
3. 实装首个核心页面（陪玩师浏览）

**需要确认**:
- [ ] 是否使用单体应用还是分离项目？
- [ ] 是否需要立即实装WebSocket？
- [ ] 优先实装用户端还是陪玩师端？

---

**准备好开始实装了吗？我建议从用户端的陪玩师浏览页面开始！** 🎮
