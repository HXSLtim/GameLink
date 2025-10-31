# 🎮 用户侧和陪玩师侧页面完整指南

**更新时间**: 2025-10-31
**页面类型**: C端用户 + 陪玩师端
**文档类型**: 页面设计与API对应说明

---

## 📑 目录

1. [系统架构说明](#1-系统架构说明)
2. [用户侧页面（C端）](#2-用户侧页面c端)
3. [陪玩师侧页面](#3-陪玩师侧页面)
4. [页面与API对应关系](#4-页面与api对应关系)
5. [路由设计](#5-路由设计)
6. [权限控制](#6-权限控制)
7. [页面开发建议](#7-页面开发建议)

---

## 1. 系统架构说明

### 1.1 三端架构

```
┌─────────────────────────────────────┐
│              前端应用                │
├─────────────────┬───────────────────┤
│   管理端页面     │   用户端页面       │
│   (Admin)      │   (User)         │
│               │                   │
│ - 用户管理     │ - 浏览陪玩师       │
│ - 游戏管理     │ - 下单流程        │
│ - 订单管理     │ - 支付流程        │
│ - 支付管理     │ - 评价系统        │
│ - 评价管理     │                   │
│ - 陪玩师管理   │   陪玩师端页面     │
│ - 权限管理     │   (Player)       │
│ - 统计报表     │                   │
│               │ - 资料管理        │
└───────────────┴───────────────────┘
           │
           ▼
┌─────────────────────────────────────┐
│              后端API                │
├─────────────────┬───────────────────┤
│   管理端API     │   用户端API       │
│   /admin/*     │   /user/*        │
│               │   /player/*      │
└───────────────┴───────────────────┘
```

### 1.2 角色说明

| 角色 | 描述 | 页面路径 | 认证 |
|------|------|----------|------|
| **C端用户** | 需要陪玩服务的普通用户 | `/user/*` | 需要 |
| **陪玩师** | 提供陪玩服务的玩家 | `/player/*` | 需要 |
| **管理员** | 平台管理人员 | `/admin/*` | 需要 |

---

## 2. 用户侧页面（C端）

### 2.1 首页/陪玩师列表 (`/user/players`)

**后端API**: `GET /api/v1/user/players`

#### 功能特性
- ✅ 陪玩师列表展示
- ✅ 搜索筛选 (游戏、价格、评分、在线状态)
- ✅ 排序功能 (价格、评分、订单数)
- ✅ 分页功能
- ✅ 陪玩师详情查看

#### 页面结构
```
顶部: 搜索栏 + 筛选器
主体: 陪玩师卡片列表
  ├─ 陪玩师头像
  ├─ 昵称 + 等级
  ├─ 主玩游戏
  ├─ 技能标签
  ├─ 评分 + 评价数
  ├─ 价格 (每小时)
  └─ 在线状态
底部: 分页组件
```

#### 组件依赖
- SearchBar (搜索栏)
- FilterPanel (筛选面板)
- PlayerCard (陪玩师卡片)
- Pagination (分页)
- Empty (空状态)
- Skeleton (加载骨架)

#### API调用
```typescript
// 获取陪玩师列表
GET /api/v1/user/players?gameId=&minPrice=&maxPrice=&minRating=&onlineOnly=&sortBy=&page=&pageSize=

// 响应数据
{
  success: true,
  data: {
    players: [
      {
        id: 1,
        nickname: "陪玩小王",
        avatar: "https://...",
        mainGame: { id: 1, name: "王者荣耀" },
        skillTags: ["ADC", "辅助"],
        rating: 4.8,
        reviewCount: 120,
        hourlyRateCents: 5000,
        onlineStatus: "online",
        isVerified: true
      }
    ],
    pagination: {
      current: 1,
      pageSize: 20,
      total: 100
    }
  }
}
```

---

### 2.2 陪玩师详情页 (`/user/players/:id`)

**后端API**: `GET /api/v1/user/players/:id`

#### 功能特性
- ✅ 陪玩师详细信息
- ✅ 作品展示
- ✅ 评价列表
- ✅ 可预约时间段
- ✅ 在线咨询
- ✅ 立即预约按钮

#### 页面结构
```
顶部: 陪玩师基本信息
├─ 头像 + 昵称
├─ 等级 + 认证标识
├─ 主玩游戏
└─ 在线状态

主体: 详细信息
├─ 技能介绍
├─ 服务价格
├─ 可预约时间
├─ 作品展示 (轮播图)
├─ 评价列表 (分页)
└─ 其他陪玩师推荐

底部: 操作栏
├─ 在线咨询
└─ 立即预约
```

#### 组件依赖
- Avatar (头像)
- Tag (标签)
- Rating (评分)
- Carousel (轮播)
- Timeline (时间线)
- CommentList (评价列表)
- ButtonGroup (操作按钮)
- FloatButton (悬浮按钮)

#### API调用
```typescript
// 获取陪玩师详情
GET /api/v1/user/players/:id

// 响应数据
{
  success: true,
  data: {
    id: 1,
    nickname: "陪玩小王",
    avatar: "https://...",
    bio: "专业陪玩，5年经验",
    level: 5,
    mainGame: { id: 1, name: "王者荣耀", icon: "..." },
    skillTags: ["ADC", "辅助", "上单"],
    rating: 4.8,
    reviewCount: 120,
    hourlyRateCents: 5000,
    onlineStatus: "online",
    isVerified: true,
    works: [
      { id: 1, title: "精彩操作集锦", url: "https://...", thumbnail: "..." }
    ],
    schedules: [
      { date: "2025-11-01", slots: ["10:00-12:00", "14:00-16:00", "20:00-22:00"] }
    ],
    reviews: [
      {
        id: 1,
        user: { nickname: "玩家A" },
        rating: 5,
        comment: "技术很好，态度也很棒",
        createdAt: "2025-10-30"
      }
    ]
  }
}
```

---

### 2.3 创建订单页 (`/user/orders/create`)

**后端API**: `POST /api/v1/user/orders`

#### 功能特性
- ✅ 选择陪玩师
- ✅ 选择游戏
- ✅ 选择时长
- ✅ 选择时间段
- ✅ 填写需求
- ✅ 确认订单信息
- ✅ 选择支付方式

#### 页面结构
```
步骤条: 1.选择陪玩师 → 2.选择服务 → 3.确认信息 → 4.支付

步骤1: 选择陪玩师
├─ 搜索陪玩师
└─ 选择陪玩师卡片

步骤2: 选择服务
├─ 游戏选择
├─ 时长选择 (1h, 2h, 3h, 5h, 自定义)
├─ 时间段选择 (日期 + 时间)
└─ 需求描述

步骤3: 确认信息
├─ 订单详情
├─ 服务条款
└─ 优惠券选择

步骤4: 支付
├─ 支付方式选择
└─ 确认支付
```

#### 组件依赖
- Steps (步骤条)
- SearchBar (搜索)
- PlayerCard (陪玩师卡片)
- Select (选择器)
- DatePicker (日期选择)
- TimePicker (时间选择)
- Input.TextArea (需求描述)
- PriceDetail (价格明细)
- PaymentMethod (支付方式)

#### API调用
```typescript
// 创建订单
POST /api/v1/user/orders
{
  playerId: 1,
  gameId: 1,
  duration: 120, // 分钟
  startTime: "2025-11-01 14:00",
  requirements: "希望耐心一点，帮忙上分",
  couponCode: "NEWUSER10"
}

// 响应数据
{
  success: true,
  data: {
    orderId: 10001,
    amountCents: 10000,
    paymentUrl: "https://payment.example.com/...",
    expireAt: "2025-11-01 14:05"
  }
}
```

---

### 2.4 我的订单列表 (`/user/orders`)

**后端API**: `GET /api/v1/user/orders`

#### 功能特性
- ✅ 订单列表展示
- ✅ 状态筛选 (待支付、已支付、进行中、已完成、已取消)
- ✅ 搜索订单
- ✅ 订单详情
- ✅ 评价已完成订单

#### 页面结构
```
顶部: 订单状态筛选 + 搜索
主体: 订单卡片列表
  ├─ 订单编号 + 状态
  ├─ 陪玩师信息 (头像 + 昵称)
  ├─ 游戏信息
  ├─ 服务时间
  ├─ 订单金额
  └─ 操作按钮
底部: 分页组件
```

#### 组件依赖
- Tabs (状态筛选)
- SearchBar (搜索)
- OrderCard (订单卡片)
- Badge (状态标识)
- Button (操作按钮)
- Pagination (分页)
- Empty (空状态)

#### API调用
```typescript
// 获取我的订单列表
GET /api/v1/user/orders?status=&page=&pageSize=

// 响应数据
{
  success: true,
  data: {
    orders: [
      {
        id: 10001,
        status: "pending",
        player: {
          id: 1,
          nickname: "陪玩小王",
          avatar: "https://..."
        },
        game: { name: "王者荣耀" },
        duration: 120,
        startTime: "2025-11-01 14:00",
        amountCents: 10000,
        createdAt: "2025-10-31",
        canReview: false
      }
    ],
    pagination: { current: 1, pageSize: 20, total: 50 }
  }
}
```

---

### 2.5 订单详情页 (`/user/orders/:id`)

**后端API**: `GET /api/v1/user/orders/:id`

#### 功能特性
- ✅ 订单详细信息
- ✅ 订单状态流转
- ✅ 陪玩师联系
- ✅ 订单操作 (取消、确认完成)
- ✅ 支付信息
- ✅ 评价功能

#### 页面结构
```
顶部: 面包屑 + 状态流转
主体: 详细信息
├─ 基本信息卡片
│  ├─ 订单编号
│  ├─ 订单状态
│  ├─ 创建时间
│  └─ 支付时间
├─ 服务信息卡片
│  ├─ 陪玩师信息
│  ├─ 游戏信息
│  ├─ 服务时间
│  └─ 需求描述
├─ 支付信息卡片
│  ├─ 订单金额
│  ├─ 优惠金额
│  ├─ 实际支付
│  └─ 支付方式
└─ 评价卡片 (已完成订单)
底部: 操作按钮
├─ 取消订单 (待支付状态)
├─ 确认完成 (进行中状态)
└─ 立即评价 (已完成状态)
```

#### 组件依赖
- Breadcrumb (面包屑)
- Steps (状态流转)
- Card (卡片)
- Description (描述列表)
- Button (操作按钮)
- Modal (弹窗)
- Rating (评分)
- Comment (评价)

#### API调用
```typescript
// 获取订单详情
GET /api/v1/user/orders/:id

// 响应数据
{
  success: true,
  data: {
    id: 10001,
    status: "pending",
    player: { id: 1, nickname: "陪玩小王", avatar: "...", phone: "..." },
    game: { name: "王者荣耀" },
    duration: 120,
    startTime: "2025-11-01 14:00",
    endTime: "2025-11-01 16:00",
    requirements: "希望耐心一点",
    amountCents: 10000,
    discountCents: 1000,
    paidAmountCents: 9000,
    paymentMethod: "alipay",
    paymentTime: "2025-10-31 10:00",
    review: null,
    timeline: [
      { time: "2025-10-31 10:00", status: "created", message: "订单已创建" },
      { time: "2025-10-31 10:05", status: "paid", message: "支付成功" }
    ]
  }
}
```

---

### 2.6 支付页面 (`/user/payments/:id`)

**后端API**: `GET /api/v1/user/payments/:id`, `POST /api/v1/user/payments`

#### 功能特性
- ✅ 订单支付
- ✅ 支付方式选择
- ✅ 支付状态查询
- ✅ 支付成功/失败处理

#### 页面结构
```
顶部: 订单信息
主体: 支付方式选择
├─ 支付宝
├─ 微信支付
├─ 银行卡
└─ 其他

底部: 确认支付按钮
```

#### 组件依赖
- PaymentMethod (支付方式)
- QRCode (二维码)
- Countdown (倒计时)
- Button (按钮)

---

### 2.7 我的评价 (`/user/reviews`)

**后端API**: `GET /api/v1/user/reviews/my`

#### 功能特性
- ✅ 我的评价列表
- ✅ 评价详情
- ✅ 追加评价
- ✅ 删除评价

#### 页面结构
```
顶部: 统计信息 (总数、平均分)
主体: 评价列表
  ├─ 陪玩师信息
  ├─ 游戏信息
  ├─ 评分
  ├─ 评价内容
  ├─ 评价时间
  └─ 操作按钮
底部: 分页组件
```

#### 组件依赖
- Statistic (统计)
- ReviewCard (评价卡片)
- Rating (评分)
- Pagination (分页)
- Empty (空状态)

#### API调用
```typescript
// 获取我的评价
GET /api/v1/user/reviews/my?page=&pageSize=

// 响应数据
{
  success: true,
  data: {
    reviews: [
      {
        id: 1,
        orderId: 10001,
        player: { nickname: "陪玩小王" },
        game: { name: "王者荣耀" },
        rating: 5,
        comment: "技术很好，态度也很棒",
        createdAt: "2025-10-31"
      }
    ],
    pagination: { current: 1, pageSize: 20, total: 10 }
  }
}
```

---

## 3. 陪玩师侧页面

### 3.1 陪玩师资料管理 (`/player/profile`)

**后端API**: `GET /player/profile`, `PUT /player/profile`

#### 功能特性
- ✅ 个人信息管理
- ✅ 技能标签设置
- ✅ 服务价格设置
- ✅ 头像上传
- ✅ 简介编辑
- ✅ 主玩游戏设置

#### 页面结构
```
顶部: 资料完成度
主体: 资料表单
├─ 基本信息
│  ├─ 头像上传
│  ├─ 昵称
│  ├─ 真实姓名 (实名认证)
│  ├─ 手机号
│  └─ 邮箱
├─ 服务信息
│  ├─ 主玩游戏 (可多选)
│  ├─ 技能标签 (可自定义)
│  ├─ 服务价格 (每小时)
│  └─ 服务时间设置
├─ 个人介绍
│  ├─ 个人简介
│  ├─ 服务说明
│  └─ 作品展示 (图片/视频)
└─ 认证信息
   ├─ 实名认证
   └─ 游戏账号认证
底部: 保存按钮
```

#### 组件依赖
- Upload (头像上传)
- Form (表单)
- Input (输入框)
- Select (选择器)
- Tag (标签)
- InputNumber (价格)
- TextArea (简介)
- Upload (作品展示)
- Button (按钮)

#### API调用
```typescript
// 获取陪玩师资料
GET /api/v1/player/profile

// 更新陪玩师资料
PUT /api/v1/player/profile
{
  nickname: "陪玩小王",
  bio: "专业陪玩，5年经验",
  mainGameIds: [1, 2, 3],
  skillTags: ["ADC", "辅助", "上单"],
  hourlyRateCents: 5000,
  serviceHours: {
    monday: ["09:00-12:00", "14:00-18:00"],
    tuesday: ["09:00-12:00", "14:00-18:00"]
  }
}
```

---

### 3.2 订单大厅 (`/player/orders/available`)

**后端API**: `GET /api/v1/player/orders/available`

#### 功能特性
- ✅ 可接订单列表
- ✅ 订单筛选 (游戏、价格范围)
- ✅ 订单详情查看
- ✅ 一键接单

#### 页面结构
```
顶部: 筛选器 (游戏、价格)
主体: 订单卡片列表
  ├─ 订单编号
  ├─ 用户信息 (匿名)
  ├─ 游戏信息
  ├─ 服务时间
  ├─ 订单金额
  ├─ 用户需求
  └─ 操作按钮 (查看详情、接单)
底部: 分页组件
```

#### 组件依赖
- FilterPanel (筛选面板)
- OrderCard (订单卡片)
- Button (操作按钮)
- Modal (详情弹窗)
- Pagination (分页)

#### API调用
```typescript
// 获取可接订单列表
GET /api/v1/player/orders/available?gameId=&minPrice=&maxPrice=&page=&pageSize=

// 响应数据
{
  success: true,
  data: {
    orders: [
      {
        id: 10001,
        game: { name: "王者荣耀" },
        duration: 120,
        startTime: "2025-11-01 14:00",
        amountCents: 10000,
        requirements: "希望耐心一点",
        userLevel: 5,
        createdAt: "2025-10-31"
      }
    ],
    pagination: { current: 1, pageSize: 20, total: 50 }
  }
}
```

---

### 3.3 我的订单 (`/player/orders`)

**后端API**: `GET /api/v1/player/orders/my`

#### 功能特性
- ✅ 已接订单列表
- ✅ 状态筛选 (待开始、进行中、已完成)
- ✅ 订单详情
- ✅ 开始服务
- ✅ 完成服务

#### 页面结构
```
顶部: 状态筛选 + 统计
主体: 订单列表
  ├─ 订单信息
  ├─ 用户信息
  ├─ 服务时间
  ├─ 订单金额
  └─ 操作按钮
底部: 分页组件
```

#### 组件依赖
- Tabs (状态筛选)
- Statistic (统计)
- OrderCard (订单卡片)
- Button (操作按钮)
- Pagination (分页)

#### API调用
```typescript
// 获取我的订单
GET /api/v1/player/orders/my?status=&page=&pageSize=

// 开始订单
POST /api/v1/player/orders/:id/start

// 完成订单
POST /api/v1/player/orders/:id/complete
```

---

### 3.4 收益概览 (`/player/earnings/summary`)

**后端API**: `GET /api/v1/player/earnings/summary`

#### 功能特性
- ✅ 总收益统计
- ✅ 本月收益
- ✅ 本周收益
- ✅ 今日收益
- ✅ 可提现余额
- ✅ 收益趋势图

#### 页面结构
```
顶部: 收益统计卡片
├─ 总收益
├─ 本月收益
├─ 本周收益
├─ 今日收益
└─ 可提现余额

主体: 收益趋势图
└─ 月度收益趋势

底部: 快捷操作
├─ 申请提现
└─ 收益明细
```

#### 组件依赖
- Statistic (统计卡片)
- Chart (趋势图)
- Card (卡片)
- Button (快捷操作)

#### API调用
```typescript
// 获取收益概览
GET /api/v1/player/earnings/summary

// 响应数据
{
  success: true,
  data: {
    totalEarnings: 500000, // 分
    monthlyEarnings: 80000,
    weeklyEarnings: 20000,
    dailyEarnings: 5000,
    withdrawableBalance: 450000,
    trend: [
      { month: "2025-05", amount: 30000 },
      { month: "2025-06", amount: 45000 }
    ]
  }
}
```

---

### 3.5 收益趋势 (`/player/earnings/trend`)

**后端API**: `GET /api/v1/player/earnings/trend`

#### 功能特性
- ✅ 收益趋势图
- ✅ 时间范围选择 (7天、30天、90天)
- ✅ 收益明细
- ✅ 收益统计

#### 页面结构
```
顶部: 时间范围选择
主体: 趋势图表
└─ 收益趋势图

底部: 收益明细表
├─ 日期
├─ 订单数
├─ 收益金额
└─ 累计收益
```

#### 组件依赖
- DatePicker (时间选择)
- Chart (图表)
- Table (明细表)
- Statistic (统计)

#### API调用
```typescript
// 获取收益趋势
GET /api/v1/player/earnings/trend?days=30
```

---

### 3.6 提现申请 (`/player/earnings/withdraw`)

**后端API**: `POST /api/v1/player/earnings/withdraw`

#### 功能特性
- ✅ 提现金额输入
- ✅ 提现方式选择 (银行卡、支付宝、微信)
- ✅ 提现记录查询

#### 页面结构
```
顶部: 可提现余额
主体: 提现表单
├─ 提现金额
├─ 提现方式
└─ 账户信息

底部: 提交申请

提现记录: 历史提现列表
```

#### 组件依赖
- Statistic (余额显示)
- Form (提现表单)
- InputNumber (金额)
- Select (提现方式)
- Table (提现记录)
- Card (卡片)

#### API调用
```typescript
// 申请提现
POST /api/v1/player/earnings/withdraw
{
  amountCents: 100000,
  method: "bank",
  accountInfo: {
    bankName: "中国银行",
    accountNumber: "6222 **** **** 1234",
    accountName: "陪玩小王"
  }
}

// 获取提现记录
GET /api/v1/player/earnings/withdraw-history
```

---

### 3.7 在线状态管理 (`/player/status`)

**后端API**: `PUT /api/v1/player/status`

#### 功能特性
- ✅ 在线/离线状态切换
- ✅ 服务时间设置
- ✅ 状态说明

#### 页面结构
```
顶部: 当前状态展示
主体: 状态设置
├─ 在线/离线开关
├─ 服务时间设置
└─ 状态说明

底部: 保存按钮
```

#### 组件依赖
- Switch (在线开关)
- TimeRange (时间范围)
- Input (状态说明)
- Button (保存)

#### API调用
```typescript
// 设置在线状态
PUT /api/v1/player/status
{
  online: true,
  serviceHours: {
    monday: ["09:00-12:00", "14:00-18:00"]
  },
  statusMessage: "在线，可接单"
}
```

---

## 4. 页面与API对应关系

### 4.1 用户侧API映射

| 页面路径 | 后端API | 方法 | 说明 |
|----------|---------|------|------|
| `/user/players` | `/api/v1/user/players` | GET | 获取陪玩师列表 |
| `/user/players/:id` | `/api/v1/user/players/:id` | GET | 获取陪玩师详情 |
| `/user/orders/create` | `/api/v1/user/orders` | POST | 创建订单 |
| `/user/orders` | `/api/v1/user/orders` | GET | 获取我的订单列表 |
| `/user/orders/:id` | `/api/v1/user/orders/:id` | GET | 获取订单详情 |
| `/user/orders/:id/cancel` | `/api/v1/user/orders/:id/cancel` | PUT | 取消订单 |
| `/user/orders/:id/complete` | `/api/v1/user/orders/:id/complete` | PUT | 确认完成 |
| `/user/payments` | `/api/v1/user/payments` | POST | 创建支付 |
| `/user/payments/:id` | `/api/v1/user/payments/:id` | GET | 查询支付状态 |
| `/user/reviews/my` | `/api/v1/user/reviews/my` | GET | 获取我的评价 |

### 4.2 陪玩师侧API映射

| 页面路径 | 后端API | 方法 | 说明 |
|----------|---------|------|------|
| `/player/profile` | `/api/v1/player/profile` | GET/PUT | 获取/更新资料 |
| `/player/apply` | `/api/v1/player/apply` | POST | 申请成为陪玩师 |
| `/player/status` | `/api/v1/player/status` | PUT | 设置在线状态 |
| `/player/orders/available` | `/api/v1/player/orders/available` | GET | 获取可接订单 |
| `/player/orders/my` | `/api/v1/player/orders/my` | GET | 获取我的订单 |
| `/player/orders/:id/accept` | `/api/v1/player/orders/:id/accept` | POST | 接单 |
| `/player/orders/:id/complete` | `/api/v1/player/orders/:id/complete` | PUT | 完成订单 |
| `/player/earnings/summary` | `/api/v1/player/earnings/summary` | GET | 获取收益概览 |
| `/player/earnings/trend` | `/api/v1/player/earnings/trend` | GET | 获取收益趋势 |
| `/player/earnings/withdraw` | `/api/v1/player/earnings/withdraw` | POST | 申请提现 |
| `/player/earnings/withdraw-history` | `/api/v1/player/earnings/withdraw-history` | GET | 获取提现记录 |

---

## 5. 路由设计

### 5.1 路由结构

```typescript
// 用户侧路由
{
  path: '/user',
  element: <UserLayout />, // 用户端布局
  children: [
    { path: 'players', element: <PlayerListPage /> },
    { path: 'players/:id', element: <PlayerDetailPage /> },
    { path: 'orders/create', element: <CreateOrderPage /> },
    { path: 'orders', element: <OrderListPage /> },
    { path: 'orders/:id', element: <OrderDetailPage /> },
    { path: 'payments/:id', element: <PaymentPage /> },
    { path: 'reviews', element: <ReviewListPage /> },
  ],
}

// 陪玩师侧路由
{
  path: '/player',
  element: <PlayerLayout />, // 陪玩师端布局
  children: [
    { path: 'profile', element: <ProfilePage /> },
    { path: 'status', element: <StatusPage /> },
    { path: 'orders/available', element: <AvailableOrdersPage /> },
    { path: 'orders', element: <MyOrdersPage /> },
    { path: 'earnings/summary', element: <EarningsSummaryPage /> },
    { path: 'earnings/trend', element: <EarningsTrendPage /> },
    { path: 'earnings/withdraw', element: <WithdrawPage /> },
    { path: 'earnings/withdraw-history', element: <WithdrawHistoryPage /> },
  ],
}
```

### 5.2 布局设计

#### 用户端布局 (`UserLayout`)
```
┌─────────────────────────────────┐
│           Header                │ ← 导航 + 用户菜单
├─────────────────────────────────┤
│                                 │
│                                 │
│           Main Content           │ ← 页面主体
│                                 │
│                                 │
├─────────────────────────────────┤
│           Footer                │ ← 底部信息 (可选)
└─────────────────────────────────┘
```

#### 陪玩师端布局 (`PlayerLayout`)
```
┌─────────────────────────────────┐
│           Header                │ ← 导航 + 在线状态
├─────┬───────────────────────────┤
│     │                           │
│ Side │      Main Content        │ ← 侧边栏 + 主体
│ bar  │                           │
│     │                           │
└─────┴───────────────────────────┘
```

---

## 6. 权限控制

### 6.1 角色判断

```typescript
interface User {
  id: number
  role: 'user' | 'player' | 'admin'
  isPlayerVerified: boolean
}
```

### 6.2 页面访问控制

```typescript
// 用户侧页面 - 需要 user 或 player 角色
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

// 陪玩师侧页面 - 需要 player 角色且已认证
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
```

### 6.3 功能权限控制

```typescript
// 接单按钮 - 只有陪玩师可以操作
{user?.role === 'player' && (
  <Button type="primary" onClick={handleAcceptOrder}>
    接单
  </Button>
)}

// 评价按钮 - 只有订单完成且未评价
{order.status === 'completed' && !order.hasReview && (
  <Button onClick={handleReview}>
    评价
  </Button>
)}
```

---

## 7. 页面开发建议

### 7.1 开发优先级

#### Phase 1: 基础功能 (MVP)
1. **用户侧**
   - 陪玩师列表页
   - 陪玩师详情页
   - 创建订单页
   - 我的订单页

2. **陪玩师侧**
   - 资料管理页
   - 订单大厅页
   - 我的订单页

#### Phase 2: 增强功能
1. **用户侧**
   - 订单详情页
   - 支付页面
   - 评价页面

2. **陪玩师侧**
   - 收益概览页
   - 提现功能

#### Phase 3: 完整功能
1. **用户侧**
   - 订单操作 (取消、完成)
   - 评价管理

2. **陪玩师侧**
   - 收益趋势
   - 在线状态管理
   - 提现记录

### 7.2 组件复用

#### 通用组件
- `PlayerCard` - 陪玩师卡片 (用户端、陪玩师端通用)
- `OrderCard` - 订单卡片 (用户端、陪玩师端通用)
- `RatingDisplay` - 评分展示
- `PriceTag` - 价格标签
- `StatusBadge` - 状态标识

#### 页面专属组件
- `UserLayout` - 用户端布局
- `PlayerLayout` - 陪玩师端布局
- `OrderTimeline` - 订单时间线
- `EarningsChart` - 收益图表
- `WithdrawForm` - 提现表单

### 7.3 数据流设计

#### 状态管理 (Context)
```typescript
// 用户状态
interface UserContextType {
  user: User | null
  isPlayer: boolean
  isPlayerVerified: boolean
  login: (credentials) => Promise<void>
  logout: () => void
}

// 订单状态
interface OrderContextType {
  orders: Order[]
  currentOrder: Order | null
  fetchOrders: (params) => Promise<void>
  createOrder: (data) => Promise<void>
  cancelOrder: (id) => Promise<void>
  completeOrder: (id) => Promise<void>
}
```

#### API调用封装
```typescript
// 用户端API
export const userApi = {
  // 陪玩师相关
  getPlayers: (params) => client.get('/user/players', { params }),
  getPlayerDetail: (id) => client.get(`/user/players/${id}`),

  // 订单相关
  createOrder: (data) => client.post('/user/orders', data),
  getMyOrders: (params) => client.get('/user/orders', { params }),
  getOrderDetail: (id) => client.get(`/user/orders/${id}`),
  cancelOrder: (id, reason) => client.put(`/user/orders/${id}/cancel`, { reason }),
  completeOrder: (id) => client.put(`/user/orders/${id}/complete`),

  // 评价相关
  getMyReviews: (params) => client.get('/user/reviews/my', { params }),
  createReview: (data) => client.post('/user/reviews', data),
}

// 陪玩师端API
export const playerApi = {
  // 资料相关
  getProfile: () => client.get('/player/profile'),
  updateProfile: (data) => client.put('/player/profile', data),
  applyAsPlayer: (data) => client.post('/player/apply', data),
  setStatus: (data) => client.put('/player/status', data),

  // 订单相关
  getAvailableOrders: (params) => client.get('/player/orders/available', { params }),
  getMyOrders: (params) => client.get('/player/orders/my', { params }),
  acceptOrder: (id) => client.post(`/player/orders/${id}/accept`),
  completeOrder: (id) => client.put(`/player/orders/${id}/complete`),

  // 收益相关
  getEarningsSummary: () => client.get('/player/earnings/summary'),
  getEarningsTrend: (days) => client.get('/player/earnings/trend', { params: { days } }),
  requestWithdraw: (data) => client.post('/player/earnings/withdraw', data),
  getWithdrawHistory: (params) => client.get('/player/earnings/withdraw-history', { params }),
}
```

### 7.4 样式设计

#### 主题色彩
```css
/* 用户端 */
:root {
  --primary-color: #1890ff; /* 主色调 */
  --success-color: #52c41a; /* 成功色 */
  --warning-color: #faad14; /* 警告色 */
  --error-color: #f5222d;   /* 错误色 */
}

/* 陪玩师端 */
.player-theme {
  --primary-color: #722ed1; /* 紫色主题 */
}
```

#### 布局样式
```css
/* 用户端 - 居中布局 */
.user-layout {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
}

/* 陪玩师端 - 侧边栏布局 */
.player-layout {
  display: flex;
  min-height: 100vh;
}

.player-sidebar {
  width: 240px;
  background: #fff;
}

.player-content {
  flex: 1;
  padding: 24px;
  background: #f5f5f5;
}
```

---

## 📚 相关文档

- [前端开发完整指南](./FRONTEND_DEVELOPMENT_COMPLETE_GUIDE.md)
- [前端页面结构文档](./FRONTEND_PAGES_STRUCTURE.md)
- [后端API文档](../../backend/docs/swagger.yaml)
- [组件库文档](./组件库文档.md)

---

**文档维护者**: GameLink Frontend Team
**最后更新**: 2025-10-31
**版本**: v1.0
**用户侧页面**: 7个
**陪玩师侧页面**: 7个
**API对应**: 20个接口
