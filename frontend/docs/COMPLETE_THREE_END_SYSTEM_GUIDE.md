# 🏗️ GameLink 三端页面体系完整指南

**更新时间**: 2025-10-31
**系统类型**: 管理端 + 用户端 + 陪玩师端
**文档类型**: 三端页面体系总结

---

## 📑 目录

1. [三端系统架构](#1-三端系统架构)
2. [页面总数统计](#2-页面总数统计)
3. [管理端页面](#3-管理端页面)
4. [用户端页面](#4-用户端页面)
5. [陪玩师端页面](#5-陪玩师端页面)
6. [跨端功能对比](#6-跨端功能对比)
7. [路由体系设计](#7-路由体系设计)
8. [权限体系](#8-权限体系)
9. [开发建议](#9-开发建议)

---

## 1. 三端系统架构

### 1.1 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                    GameLink 前端应用                         │
├─────────────────────────────────┬───────────────────────────┤
│        管理端 (Admin)            │        用户端 (User)       │
│        /admin/*                 │        /user/*           │
│                                 │                           │
│ - 后台管理系统                   │ - 陪玩服务平台           │
│ - 数据管理                       │ - 浏览陪玩师             │
│ - 用户管理                       │ - 下单流程               │
│ - 订单管理                       │ - 支付流程               │
│ - 支付管理                       │ - 评价系统               │
│ - 陪玩师管理                     │                           │
│ - 评价管理                       │   陪玩师端 (Player)      │
│ - 权限管理                       │   /player/*             │
│ - 统计报表                       │                           │
│                                 │ - 陪玩师工作台           │
└─────────────────────────────────┴───────────────────────────┘
           │                                   │
           └─────────────────┬─────────────────┘
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                   GameLink 后端API                          │
├─────────────────┬───────────────────┬─────────────────────────┤
│   管理端API      │    用户端API      │      陪玩师端API       │
│   /admin/*      │    /user/*       │      /player/*         │
│                 │                  │                        │
│ ✓ 用户管理      │ ✓ 陪玩师列表     │ ✓ 资料管理            │
│ ✓ 游戏管理      │ ✓ 订单创建       │ ✓ 订单管理            │
│ ✓ 订单管理      │ ✓ 订单查询       │ ✓ 收益管理            │
│ ✓ 支付管理      │ ✓ 支付处理       │ ✓ 提现功能            │
│ ✓ 陪玩师管理    │ ✓ 评价管理       │ ✓ 状态管理            │
│ ✓ 评价管理      │                  │                        │
│ ✓ 权限管理      │                  │                        │
│ ✓ 统计报表      │                  │                        │
└─────────────────┴──────────────────┴────────────────────────┘
```

### 1.2 三端角色说明

| 角色 | 描述 | 页面前缀 | 用户特征 | 主要需求 |
|------|------|----------|----------|----------|
| **管理员** | 平台运营人员 | `/admin` | 系统管理者 | 数据统计、用户管理、订单管理 |
| **C端用户** | 需要陪玩服务的普通用户 | `/user` | 消费者 | 浏览陪玩师、下单、支付、评价 |
| **陪玩师** | 提供陪玩服务的玩家 | `/player` | 服务提供者 | 管理资料、接单、收益管理、提现 |

### 1.3 技术实现说明

```
前端: React 18 + TypeScript + Vite
├─ 管理端: 基于React Router的SPA应用
├─ 用户端: 基于React Router的SPA应用
└─ 陪玩师端: 基于React Router的SPA应用

后端: Go + Gin + GORM
├─ 管理API: /admin/*
├─ 用户API: /user/*
└─ 陪玩师API: /player/*
```

---

## 2. 页面总数统计

### 2.1 按端分类统计

| 端别 | 页面数量 | 占比 | 状态 |
|------|----------|------|------|
| **管理端** | 17页 | 48.6% | ✅ 已实现 |
| **用户端** | 7页 | 20.0% | 📋 待开发 |
| **陪玩师端** | 7页 | 20.0% | 📋 待开发 |
| **共享页面** | 4页 | 11.4% | ✅ 已有 |
| **总计** | **35页** | **100%** | **37.1%完成** |

### 2.2 按功能分类统计

| 功能模块 | 管理端 | 用户端 | 陪玩师端 | 总计 |
|----------|--------|--------|----------|------|
| 认证相关 | 2页 | 0页 | 0页 | 2页 |
| 仪表盘/首页 | 1页 | 0页 | 0页 | 1页 |
| 用户/陪玩师 | 2页 | 2页 | 1页 | 5页 |
| 游戏管理 | 2页 | 0页 | 0页 | 2页 |
| 订单管理 | 3页 | 2页 | 2页 | 7页 |
| 支付管理 | 2页 | 1页 | 0页 | 3页 |
| 评价管理 | 2页 | 1页 | 0页 | 3页 |
| 收益管理 | 0页 | 0页 | 3页 | 3页 |
| 权限管理 | 1页 | 0页 | 0页 | 1页 |
| 统计报表 | 1页 | 0页 | 0页 | 1页 |
| 设置页面 | 1页 | 0页 | 1页 | 2页 |
| 演示页面 | 2页 | 0页 | 0页 | 2页 |
| **总计** | **17页** | **7页** | **7页** | **31页** |

### 2.3 完成度统计

| 端别 | 已实现 | 待开发 | 完成度 |
|------|--------|--------|--------|
| 管理端 | 17页 | 0页 | ✅ 100% |
| 用户端 | 0页 | 7页 | 📋 0% |
| 陪玩师端 | 0页 | 7页 | 📋 0% |
| 共享页面 | 4页 | 0页 | ✅ 100% |
| **总计** | **21页** | **14页** | **60.0%** |

---

## 3. 管理端页面

### 3.1 页面列表 (17页)

#### 认证页面 (2页)
| 页面 | 路径 | 文件路径 | 状态 |
|------|------|----------|------|
| 登录页 | `/login` | `src/pages/Login/Login.tsx` | ✅ |
| 注册页 | `/register` | `src/pages/Register/Register.tsx` | ✅ |

#### 核心管理页面 (15页)
| 页面 | 路径 | 文件路径 | 状态 |
|------|------|----------|------|
| 仪表盘 | `/dashboard` | `src/pages/Dashboard/Dashboard.tsx` | ✅ |
| 用户管理-列表 | `/users` | `src/pages/Users/UserList.tsx` | ✅ |
| 用户管理-详情 | `/users/:id` | `src/pages/Users/UserDetail.tsx` | ✅ |
| 用户管理-表单 | 弹窗 | `src/pages/Users/UserFormModal.tsx` | ✅ |
| 游戏管理-列表 | `/games` | `src/pages/Games/GameList.tsx` | ✅ |
| 游戏管理-详情 | `/games/:id` | `src/pages/Games/GameDetail.tsx` | ✅ |
| 游戏管理-表单 | 弹窗 | `src/pages/Games/GameFormModal.tsx` | ✅ |
| 订单管理-列表 | `/orders` | `src/pages/Orders/OrderList.tsx` | ✅ |
| 订单管理-详情 | `/orders/:id` | `src/pages/Orders/OrderDetail.tsx` | ✅ |
| 订单管理-表单 | 弹窗 | `src/pages/Orders/OrderFormModal.tsx` | ✅ |
| 支付管理-列表 | `/payments` | `src/pages/Payments/PaymentList.tsx` | ✅ |
| 支付管理-详情 | `/payments/:id` | `src/pages/Payments/PaymentDetailPage.tsx` | ✅ |
| 陪玩师管理-列表 | `/players` | `src/pages/Players/PlayerList.tsx` | ✅ |
| 陪玩师管理-表单 | 弹窗 | `src/pages/Players/PlayerFormModal.tsx` | ✅ |
| 评价管理-列表 | `/reviews` | `src/pages/Reviews/ReviewList.tsx` | ✅ |
| 评价管理-表单 | 弹窗 | `src/pages/Reviews/ReviewFormModal.tsx` | ✅ |
| 报表统计 | `/reports` | `src/pages/Reports/ReportDashboard.tsx` | ✅ |
| 权限管理 | `/permissions` | `src/pages/Permissions/PermissionList.tsx` | ✅ |
| 系统设置 | `/settings` | `src/pages/Settings/SettingsDashboard.tsx` | ✅ |

#### 演示页面 (2页)
| 页面 | 路径 | 文件路径 | 状态 |
|------|------|----------|------|
| 组件演示 | `/showcase` | `src/pages/ComponentsDemo/ComponentsDemo.tsx` | ✅ |
| 缓存演示 | `/cache-demo` | `src/pages/CacheDemo/index.tsx` | ✅ |

### 3.2 功能特性

#### ✅ 已实现功能
- 用户增删改查
- 游戏信息管理
- 订单全流程管理
- 支付记录查看
- 陪玩师信息管理
- 评价审核
- 权限控制 (RBAC)
- 数据统计报表
- 系统设置

#### 📊 数据展示
- 实时数据统计
- 图表可视化 (收入趋势、用户增长)
- 数据导出功能
- 操作日志记录

---

## 4. 用户端页面

### 4.1 页面列表 (7页)

| 页面 | 路径 | 后端API | 状态 | 优先级 |
|------|------|---------|------|--------|
| 陪玩师列表 | `/user/players` | `GET /user/players` | 📋 待开发 | P0 |
| 陪玩师详情 | `/user/players/:id` | `GET /user/players/:id` | 📋 待开发 | P0 |
| 创建订单 | `/user/orders/create` | `POST /user/orders` | 📋 待开发 | P0 |
| 我的订单 | `/user/orders` | `GET /user/orders` | 📋 待开发 | P0 |
| 订单详情 | `/user/orders/:id` | `GET /user/orders/:id` | 📋 待开发 | P1 |
| 支付页面 | `/user/payments/:id` | `GET /user/payments/:id` | 📋 待开发 | P1 |
| 我的评价 | `/user/reviews` | `GET /user/reviews/my` | 📋 待开发 | P2 |

### 4.2 页面详细设计

#### 4.2.1 陪玩师列表页 (`/user/players`)

**功能描述**: 用户浏览可预约的陪玩师列表

**核心功能**:
- ✅ 陪玩师卡片展示
- ✅ 搜索筛选 (游戏、价格、评分)
- ✅ 在线状态筛选
- ✅ 排序功能
- ✅ 分页加载

**页面结构**:
```
顶部: 搜索栏 + 筛选器
主体: 陪玩师卡片网格布局 (每行3-4个)
底部: 分页组件
```

**技术实现**:
```typescript
// 组件层级
UserLayout
└─ PlayerListPage
   ├─ SearchBar
   ├─ FilterPanel
   ├─ PlayerCardGrid
   │  └─ PlayerCard[]
   └─ Pagination
```

**API调用**:
```typescript
GET /api/v1/user/players?gameId=&minPrice=&maxPrice=&minRating=&onlineOnly=&sortBy=&page=&pageSize=
```

---

#### 4.2.2 陪玩师详情页 (`/user/players/:id`)

**功能描述**: 查看陪玩师详细信息和评价

**核心功能**:
- ✅ 陪玩师基本信息
- ✅ 技能标签展示
- ✅ 作品展示
- ✅ 评价列表
- ✅ 立即预约按钮

**页面结构**:
```
顶部: 陪玩师基本信息卡片
主体: 详细信息
├─ 技能介绍卡片
├─ 服务价格卡片
├─ 作品展示卡片
├─ 评价列表卡片
└─ 推荐陪玩师卡片
底部: 悬浮操作按钮 (在线咨询、立即预约)
```

---

#### 4.2.3 创建订单页 (`/user/orders/create`)

**功能描述**: 创建陪玩服务订单

**核心功能**:
- ✅ 选择陪玩师
- ✅ 选择游戏
- ✅ 选择时长
- ✅ 选择时间段
- ✅ 填写需求
- ✅ 确认订单

**页面结构**:
```
步骤条: 1.选择陪玩师 → 2.选择服务 → 3.确认信息 → 4.支付

每个步骤对应一个表单页面
```

**状态流转**:
```
步骤1 → 步骤2 → 步骤3 → 步骤4 → 创建订单成功
```

---

#### 4.2.4 我的订单页 (`/user/orders`)

**功能描述**: 查看用户的订单历史

**核心功能**:
- ✅ 订单列表
- ✅ 状态筛选
- ✅ 搜索订单
- ✅ 快速操作

**页面结构**:
```
顶部: 状态筛选标签 + 搜索框
主体: 订单卡片列表
底部: 分页组件
```

**订单状态**:
- 待支付 (pending)
- 已支付 (paid)
- 进行中 (in_progress)
- 已完成 (completed)
- 已取消 (cancelled)

---

#### 4.2.5 订单详情页 (`/user/orders/:id`)

**功能描述**: 查看订单详细信息

**核心功能**:
- ✅ 订单详情
- ✅ 陪玩师信息
- ✅ 订单时间线
- ✅ 订单操作

**页面结构**:
```
顶部: 面包屑导航
主体: 详细信息卡片
├─ 基本信息
├─ 服务信息
├─ 支付信息
└─ 评价信息
底部: 操作按钮
```

---

#### 4.2.6 支付页面 (`/user/payments/:id`)

**功能描述**: 完成订单支付

**核心功能**:
- ✅ 订单信息确认
- ✅ 支付方式选择
- ✅ 支付状态查询

**页面结构**:
```
顶部: 订单信息
主体: 支付方式选择
底部: 确认支付按钮
```

---

#### 4.2.7 我的评价页 (`/user/reviews`)

**功能描述**: 查看用户的历史评价

**核心功能**:
- ✅ 评价列表
- ✅ 评价统计
- ✅ 追加评价

**页面结构**:
```
顶部: 评价统计 (总数、平均分)
主体: 评价列表
底部: 分页组件
```

---

## 5. 陪玩师端页面

### 5.1 页面列表 (7页)

| 页面 | 路径 | 后端API | 状态 | 优先级 |
|------|------|---------|------|--------|
| 陪玩师资料 | `/player/profile` | `GET/PUT /player/profile` | 📋 待开发 | P0 |
| 在线状态 | `/player/status` | `PUT /player/status` | 📋 待开发 | P1 |
| 申请成为陪玩师 | `/player/apply` | `POST /player/apply` | 📋 待开发 | P1 |
| 订单大厅 | `/player/orders/available` | `GET /player/orders/available` | 📋 待开发 | P0 |
| 我的订单 | `/player/orders` | `GET /player/orders/my` | 📋 待开发 | P0 |
| 收益概览 | `/player/earnings/summary` | `GET /player/earnings/summary` | 📋 待开发 | P1 |
| 收益趋势 | `/player/earnings/trend` | `GET /player/earnings/trend` | 📋 待开发 | P2 |
| 提现申请 | `/player/earnings/withdraw` | `POST /player/earnings/withdraw` | 📋 待开发 | P1 |
| 提现记录 | `/player/earnings/withdraw-history` | `GET /player/earnings/withdraw-history` | 📋 待开发 | P2 |

### 5.2 页面详细设计

#### 5.2.1 陪玩师资料页 (`/player/profile`)

**功能描述**: 管理陪玩师个人资料

**核心功能**:
- ✅ 个人信息管理
- ✅ 技能标签设置
- ✅ 服务价格设置
- ✅ 头像上传
- ✅ 简介编辑

**页面结构**:
```
顶部: 资料完成度指示器
主体: 资料表单
├─ 基本信息 (头像、昵称、联系方式)
├─ 服务信息 (游戏、技能、价格)
├─ 个人介绍 (简介、作品)
└─ 认证信息 (实名认证、游戏认证)
底部: 保存按钮
```

---

#### 5.2.2 在线状态页 (`/player/status`)

**功能描述**: 设置在线/离线状态

**核心功能**:
- ✅ 在线状态切换
- ✅ 服务时间设置
- ✅ 状态说明

**页面结构**:
```
顶部: 当前状态显示
主体: 状态设置表单
├─ 在线/离线开关
├─ 服务时间设置
└─ 状态说明
底部: 保存按钮
```

---

#### 5.2.3 申请成为陪玩师页 (`/player/apply`)

**功能描述**: 普通用户申请成为陪玩师

**核心功能**:
- ✅ 申请表单
- ✅ 资料填写
- ✅ 提交审核

**页面结构**:
```
顶部: 申请说明
主体: 申请表单
底部: 提交申请按钮
```

---

#### 5.2.4 订单大厅页 (`/player/orders/available`)

**功能描述**: 查看可接订单

**核心功能**:
- ✅ 可接订单列表
- ✅ 订单筛选
- ✅ 快速接单

**页面结构**:
```
顶部: 筛选器 (游戏、价格)
主体: 订单卡片列表
底部: 分页组件
```

---

#### 5.2.5 我的订单页 (`/player/orders`)

**功能描述**: 查看已接订单

**核心功能**:
- ✅ 已接订单列表
- ✅ 状态筛选
- ✅ 订单操作

**页面结构**:
```
顶部: 状态筛选标签
主体: 订单列表
底部: 分页组件
```

---

#### 5.2.6 收益概览页 (`/player/earnings/summary`)

**功能描述**: 查看收益统计

**核心功能**:
- ✅ 总收益统计
- ✅ 本月/本周/今日收益
- ✅ 可提现余额
- ✅ 收益趋势

**页面结构**:
```
顶部: 收益统计卡片 (4个)
主体: 收益趋势图表
底部: 快捷操作按钮 (申请提现、查看明细)
```

---

#### 5.2.7 提现申请页 (`/player/earnings/withdraw`)

**功能描述**: 申请提现

**核心功能**:
- ✅ 提现金额输入
- ✅ 提现方式选择
- ✅ 提现记录查询

**页面结构**:
```
顶部: 可提现余额
主体: 提现表单
├─ 提现金额
├─ 提现方式 (银行卡、支付宝、微信)
└─ 账户信息
底部: 提交申请

提现记录: 历史提现列表
```

---

## 6. 跨端功能对比

### 6.1 订单相关功能

| 功能 | 管理端 | 用户端 | 陪玩师端 |
|------|--------|--------|----------|
| 创建订单 | ✅ 创建订单 | ✅ 创建订单 | ❌ 无 |
| 查看订单列表 | ✅ 查看所有 | ✅ 查看我的 | ✅ 查看我的 |
| 订单详情 | ✅ 查看所有 | ✅ 查看我的 | ✅ 查看我的 |
| 订单操作 | ✅ 所有操作 | ✅ 用户操作 | ✅ 陪玩师操作 |
| 订单统计 | ✅ 全部统计 | ❌ 无 | ❌ 无 |

### 6.2 用户相关功能

| 功能 | 管理端 | 用户端 | 陪玩师端 |
|------|--------|--------|----------|
| 用户管理 | ✅ 管理所有用户 | ❌ 无 | ✅ 管理自己的资料 |
| 陪玩师管理 | ✅ 管理所有陪玩师 | ✅ 浏览陪玩师 | ✅ 管理自己的资料 |
| 认证信息 | ✅ 审核认证 | ❌ 无 | ✅ 提交认证 |

### 6.3 支付相关功能

| 功能 | 管理端 | 用户端 | 陪玩师端 |
|------|--------|--------|----------|
| 创建支付 | ✅ 管理员可创建 | ✅ 用户创建 | ❌ 无 |
| 查看支付 | ✅ 查看所有 | ✅ 查看我的 | ❌ 无 |
| 退款处理 | ✅ 处理退款 | ❌ 无 | ❌ 无 |
| 收益管理 | ✅ 平台收益 | ❌ 无 | ✅ 个人收益 |

### 6.4 评价相关功能

| 功能 | 管理端 | 用户端 | 陪玩师端 |
|------|--------|--------|----------|
| 查看评价 | ✅ 查看所有 | ✅ 查看我的 | ✅ 查看对我的评价 |
| 创建评价 | ❌ 无 | ✅ 创建评价 | ❌ 无 |
| 审核评价 | ✅ 审核评价 | ❌ 无 | ❌ 无 |

---

## 7. 路由体系设计

### 7.1 完整路由结构

```typescript
export const router = createBrowserRouter([
  // 公开路由
  { path: '/login', element: <LoginPage /> },
  { path: '/register', element: <RegisterPage /> },

  // 用户端路由
  {
    path: '/user',
    element: <UserLayout />,
    children: [
      { path: 'players', element: <PlayerListPage /> },
      { path: 'players/:id', element: <PlayerDetailPage /> },
      { path: 'orders/create', element: <CreateOrderPage /> },
      { path: 'orders', element: <OrderListPage /> },
      { path: 'orders/:id', element: <OrderDetailPage /> },
      { path: 'payments/:id', element: <PaymentPage /> },
      { path: 'reviews', element: <ReviewListPage /> },
    ],
  },

  // 陪玩师端路由
  {
    path: '/player',
    element: <PlayerLayout />,
    children: [
      { path: 'profile', element: <ProfilePage /> },
      { path: 'status', element: <StatusPage /> },
      { path: 'apply', element: <ApplyPage /> },
      { path: 'orders/available', element: <AvailableOrdersPage /> },
      { path: 'orders', element: <MyOrdersPage /> },
      { path: 'earnings/summary', element: <EarningsSummaryPage /> },
      { path: 'earnings/trend', element: <EarningsTrendPage /> },
      { path: 'earnings/withdraw', element: <WithdrawPage /> },
      { path: 'earnings/withdraw-history', element: <WithdrawHistoryPage /> },
    ],
  },

  // 管理端路由
  {
    path: '/admin',
    element: <AdminLayout />,
    children: [
      { path: 'dashboard', element: <DashboardPage /> },
      { path: 'users', element: <UserListPage /> },
      { path: 'users/:id', element: <UserDetailPage /> },
      { path: 'games', element: <GameListPage /> },
      { path: 'games/:id', element: <GameDetailPage /> },
      { path: 'orders', element: <OrderListPage /> },
      { path: 'orders/:id', element: <OrderDetailPage /> },
      { path: 'payments', element: <PaymentListPage /> },
      { path: 'payments/:id', element: <PaymentDetailPage /> },
      { path: 'players', element: <PlayerListPage /> },
      { path: 'reviews', element: <ReviewListPage /> },
      { path: 'reports', element: <ReportDashboardPage /> },
      { path: 'permissions', element: <PermissionListPage /> },
      { path: 'settings', element: <SettingsPage /> },
    ],
  },

  // 演示路由
  { path: '/showcase', element: <ComponentsDemoPage /> },

  // 404
  { path: '*', element: <Navigate to="/" replace /> },
])
```

### 7.2 路由守卫设计

```typescript
// 通用路由守卫
const ProtectedRoute = ({ children }) => {
  const { user } = useAuth()

  if (!user) {
    return <Navigate to="/login" />
  }

  return children
}

// 用户端路由守卫
const UserRoute = ({ children }) => {
  const { user } = useAuth()

  if (!user) {
    return <Navigate to="/login" />
  }

  if (user.role === 'admin') {
    return <Navigate to="/admin/dashboard" />
  }

  return children
}

// 陪玩师端路由守卫
const PlayerRoute = ({ children }) => {
  const { user } = useAuth()

  if (!user) {
    return <Navigate to="/login" />
  }

  if (user.role !== 'player' || !user.isPlayerVerified) {
    return <Navigate to="/player/apply" />
  }

  return children
}

// 管理端路由守卫
const AdminRoute = ({ children }) => {
  const { user } = useAuth()

  if (!user) {
    return <Navigate to="/login" />
  }

  if (user.role !== 'admin') {
    return <Navigate to="/user/players" />
  }

  return children
}
```

---

## 8. 权限体系

### 8.1 角色定义

```typescript
interface User {
  id: number
  username: string
  email: string
  role: 'admin' | 'user' | 'player'
  isPlayerVerified: boolean // 只有player角色才有效
  permissions: string[]
}
```

### 8.2 权限检查

```typescript
// 检查用户角色
const hasRole = (requiredRole: string) => {
  const { user } = useAuth()
  return user?.role === requiredRole
}

// 检查是否认证
const isAuthenticated = () => {
  const { user } = useAuth()
  return !!user
}

// 检查是否陪玩师且已认证
const isVerifiedPlayer = () => {
  const { user } = useAuth()
  return user?.role === 'player' && user?.isPlayerVerified
}
```

### 8.3 页面访问控制

```typescript
// 根据用户角色自动跳转
const AutoRedirect = () => {
  const { user } = useAuth()

  if (!user) {
    return <Navigate to="/login" />
  }

  switch (user.role) {
    case 'admin':
      return <Navigate to="/admin/dashboard" />
    case 'player':
      if (user.isPlayerVerified) {
        return <Navigate to="/player/dashboard" />
      } else {
        return <Navigate to="/player/apply" />
      }
    default:
      return <Navigate to="/user/players" />
  }
}
```

### 8.4 组件访问控制

```typescript
// 根据角色显示/隐藏组件
{user?.role === 'admin' && (
  <AdminOnlyComponent />
)}

{isVerifiedPlayer() && (
  <PlayerFeature />
)}
```

---

## 9. 开发建议

### 9.1 开发优先级

#### Phase 1: 用户端MVP (4周)
**目标**: 让用户能够浏览陪玩师并下单

1. **用户端基础页面** (2周)
   - 陪玩师列表页 (`/user/players`)
   - 陪玩师详情页 (`/user/players/:id`)
   - 创建订单页 (`/user/orders/create`)
   - 我的订单页 (`/user/orders`)

2. **支付集成** (1周)
   - 支付页面
   - 支付回调处理

3. **基本陪玩师功能** (1周)
   - 资料管理页 (`/player/profile`)
   - 订单大厅页 (`/player/orders/available`)
   - 我的订单页 (`/player/orders`)

#### Phase 2: 陪玩师功能完善 (2周)
1. **收益管理**
   - 收益概览页
   - 提现申请页

2. **状态管理**
   - 在线状态设置
   - 申请成为陪玩师

#### Phase 3: 用户端增强 (2周)
1. **订单管理**
   - 订单详情页
   - 订单操作 (取消、完成)

2. **评价系统**
   - 我的评价页
   - 评价功能

### 9.2 组件复用策略

#### 通用组件库
```typescript
// 跨端通用组件
export const UniversalComponents = {
  // 订单相关
  OrderCard: '订单卡片',
  OrderStatus: '订单状态',
  OrderTimeline: '订单时间线',

  // 陪玩师相关
  PlayerCard: '陪玩师卡片',
  PlayerRating: '陪玩师评分',
  PlayerAvatar: '陪玩师头像',

  // 通用组件
  SearchBar: '搜索栏',
  FilterPanel: '筛选面板',
  Pagination: '分页',
  Empty: '空状态',
  Loading: '加载中',
}
```

#### 布局组件
```typescript
// 三端专用布局
export const Layouts = {
  UserLayout: '用户端布局',      // 顶部导航 + 主体内容
  PlayerLayout: '陪玩师端布局',   // 侧边栏 + 主体内容
  AdminLayout: '管理端布局',      // 侧边栏 + 顶部 + 主体内容
}
```

### 9.3 状态管理方案

```typescript
// 全局状态设计
interface GlobalState {
  // 用户状态
  auth: {
    user: User | null
    token: string | null
    loading: boolean
  }

  // 用户端状态
  user: {
    players: Player[]
    orders: Order[]
    filters: FilterState
    pagination: Pagination
  }

  // 陪玩师端状态
  player: {
    profile: PlayerProfile | null
    availableOrders: Order[]
    myOrders: Order[]
    earnings: Earnings | null
  }

  // 管理端状态
  admin: {
    users: User[]
    games: Game[]
    orders: Order[]
    stats: Stats
  }
}
```

### 9.4 API调用封装

```typescript
// 按端分类的API
export const api = {
  // 用户端API
  user: {
    players: userPlayerApi,
    orders: userOrderApi,
    payments: userPaymentApi,
    reviews: userReviewApi,
  },

  // 陪玩师端API
  player: {
    profile: playerProfileApi,
    orders: playerOrderApi,
    earnings: playerEarningsApi,
    withdraw: playerWithdrawApi,
  },

  // 管理端API
  admin: {
    users: adminUserApi,
    games: adminGameApi,
    orders: adminOrderApi,
    payments: adminPaymentApi,
    players: adminPlayerApi,
    reviews: adminReviewApi,
    stats: adminStatsApi,
  },
}
```

### 9.5 代码组织建议

```
src/
├── pages/                    # 页面组件
│   ├── user/                # 用户端页面
│   │   ├── players/
│   │   ├── orders/
│   │   ├── payments/
│   │   └── reviews/
│   ├── player/              # 陪玩师端页面
│   │   ├── profile/
│   │   ├── orders/
│   │   ├── earnings/
│   │   └── withdraw/
│   └── admin/               # 管理端页面
│       ├── users/
│       ├── games/
│       ├── orders/
│       └── ...
├── components/              # 通用组件
│   ├── user/               # 用户端专用组件
│   ├── player/             # 陪玩师端专用组件
│   ├── admin/              # 管理端专用组件
│   └── shared/             # 跨端通用组件
├── layouts/                # 布局组件
│   ├── UserLayout.tsx
│   ├── PlayerLayout.tsx
│   └── AdminLayout.tsx
├── hooks/                  # 自定义Hooks
│   ├── useAuth.ts
│   ├── useUser.ts
│   ├── usePlayer.ts
│   └── useAdmin.ts
├── api/                    # API调用
│   ├── user/
│   ├── player/
│   └── admin/
└── types/                  # 类型定义
    ├── user.ts
    ├── player.ts
    └── admin.ts
```

### 9.6 测试策略

#### 单元测试
- ✅ 通用组件测试
- ✅ 工具函数测试
- ✅ API调用测试

#### 集成测试
- ✅ 页面功能测试
- ✅ 路由跳转测试
- ✅ 用户交互测试

#### E2E测试
- ✅ 核心用户流程测试
- ✅ 关键业务场景测试

---

## 📚 附录

### A. 页面开发检查清单

#### 用户端页面检查清单
- [ ] 陪玩师列表页 (搜索、筛选、分页)
- [ ] 陪玩师详情页 (信息展示、预约入口)
- [ ] 创建订单页 (表单、验证、API调用)
- [ ] 我的订单页 (列表、状态、操作)
- [ ] 订单详情页 (详情、操作、时间线)
- [ ] 支付页面 (支付方式、状态查询)
- [ ] 我的评价页 (列表、统计)

#### 陪玩师端页面检查清单
- [ ] 资料管理页 (表单、验证、文件上传)
- [ ] 在线状态页 (开关、时间设置)
- [ ] 申请成为陪玩师页 (表单、审核流程)
- [ ] 订单大厅页 (列表、筛选、接单)
- [ ] 我的订单页 (列表、状态、操作)
- [ ] 收益概览页 (统计、图表)
- [ ] 提现申请页 (表单、记录)

### B. API开发清单

#### 用户端API
- [ ] GET /user/players - 获取陪玩师列表
- [ ] GET /user/players/:id - 获取陪玩师详情
- [ ] POST /user/orders - 创建订单
- [ ] GET /user/orders - 获取我的订单
- [ ] GET /user/orders/:id - 获取订单详情
- [ ] PUT /user/orders/:id/cancel - 取消订单
- [ ] PUT /user/orders/:id/complete - 确认完成
- [ ] POST /user/payments - 创建支付
- [ ] GET /user/payments/:id - 查询支付状态
- [ ] GET /user/reviews/my - 获取我的评价

#### 陪玩师端API
- [ ] GET /player/profile - 获取资料
- [ ] PUT /player/profile - 更新资料
- [ ] POST /player/apply - 申请成为陪玩师
- [ ] PUT /player/status - 设置状态
- [ ] GET /player/orders/available - 获取可接订单
- [ ] GET /player/orders/my - 获取我的订单
- [ ] POST /player/orders/:id/accept - 接单
- [ ] PUT /player/orders/:id/complete - 完成订单
- [ ] GET /player/earnings/summary - 获取收益概览
- [ ] GET /player/earnings/trend - 获取收益趋势
- [ ] POST /player/earnings/withdraw - 申请提现
- [ ] GET /player/earnings/withdraw-history - 获取提现记录

---

**文档维护者**: GameLink Frontend Team
**最后更新**: 2025-10-31
**版本**: v1.0
**管理端页面**: 17页 (✅100%)
**用户端页面**: 7页 (📋0%)
**陪玩师端页面**: 7页 (📋0%)
**总计**: 31页 (60%完成)
