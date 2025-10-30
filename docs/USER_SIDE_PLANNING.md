# 🎮 GameLink 用户侧前后端开发规划

**规划时间**: 2025年10月30日  
**项目阶段**: 用户侧（C端）开发  
**目标**: 为普通用户和陪玩师提供完整的业务功能

---

## 📋 目录

- [项目概述](#项目概述)
- [业务流程](#业务流程)
- [功能模块规划](#功能模块规划)
- [后端 API 规划](#后端-api-规划)
- [前端页面规划](#前端页面规划)
- [数据模型](#数据模型)
- [开发阶段](#开发阶段)
- [技术架构](#技术架构)

---

## 📖 项目概述

### 业务定位
GameLink 是一个**游戏陪玩服务平台**，连接普通玩家与专业陪玩师。

### 用户角色

| 角色 | 说明 | 主要功能 |
|------|------|----------|
| **普通用户** | 需要陪玩服务的玩家 | 浏览陪玩师、下单、支付、评价 |
| **陪玩师** | 提供陪玩服务的专业玩家 | 展示技能、接单、完成订单、收益管理 |
| **管理员** | 平台运营人员 | 审核、管理、数据分析（已完成） |

---

## 🔄 业务流程

### 核心流程图

```
[普通用户] ─────────────────────────────────────────────> [陪玩师]
    │                                                         │
    │ 1. 浏览陪玩师列表                                      │
    │ 2. 查看陪玩师详情                                      │
    │ 3. 选择游戏 & 时段                                     │
    │ 4. 创建订单 ──────────────────────> 订单通知 ────────> │ 5. 接单
    │                                                         │
    │ 6. 支付订单                                            │ 7. 确认开始
    │                                                         │
    │ 8. 游戏进行中 <────────────────────────────────────────> │
    │                                                         │
    │ 9. 确认完成 <──────────────────────────────────────────> │ 10. 确认完成
    │                                                         │
    │ 11. 评价订单                                           │
    │                                                         │
    v                                                         v
 [订单完成]                                              [获得收益]
```

### 订单状态流转

```
pending（待支付）
    ↓ 用户支付
confirmed（已支付，待接单）
    ↓ 陪玩师接单
in_progress（进行中）
    ↓ 双方确认完成
completed（已完成）
    ↓ 用户评价
[可评价状态]

取消流程：
pending → canceled（未支付可直接取消）
confirmed → canceled（已支付需退款）
    ↓
refunded（已退款）
```

---

## 🎯 功能模块规划

### 1. 用户端功能

#### 1.1 首页模块
- [ ] **游戏分类展示**
  - 热门游戏卡片展示
  - 游戏分类筛选
  - 搜索功能
- [ ] **推荐陪玩师**
  - 热门陪玩师推荐
  - 评分排序
  - 在线状态显示

#### 1.2 陪玩师模块
- [ ] **陪玩师列表**
  - 按游戏筛选
  - 按价格区间筛选
  - 按评分筛选
  - 按在线状态筛选
  - 分页加载
- [ ] **陪玩师详情页**
  - 基本信息（昵称、头像、简介）
  - 擅长游戏、段位信息
  - 服务价格（时薪）
  - 评价列表
  - 历史订单数、好评率
  - 认证状态
  - 在线状态
- [ ] **预约下单**
  - 选择游戏
  - 选择时长
  - 选择时间段
  - 价格计算
  - 备注需求

#### 1.3 订单模块
- [ ] **我的订单列表**
  - 全部订单
  - 待支付
  - 进行中
  - 已完成
  - 已取消/退款
  - 订单搜索
- [ ] **订单详情页**
  - 订单信息
  - 陪玩师信息
  - 订单状态
  - 操作按钮（支付、取消、确认完成、评价）
  - 订单时间轴
  - 联系陪玩师（聊天入口）

#### 1.4 支付模块
- [ ] **支付页面**
  - 订单信息确认
  - 支付方式选择（微信、支付宝）
  - 优惠券选择
  - 价格明细
  - 支付按钮
- [ ] **支付结果页**
  - 支付成功/失败提示
  - 订单详情跳转
  - 继续浏览

#### 1.5 评价模块
- [ ] **订单评价**
  - 星级评分（1-5星）
  - 文字评价
  - 标签选择（技术好、态度好等）
  - 匿名评价选项

#### 1.6 个人中心
- [ ] **个人信息**
  - 头像、昵称编辑
  - 联系方式管理
  - 密码修改
- [ ] **我的钱包**
  - 余额查询
  - 充值功能
  - 消费记录
  - 退款记录
- [ ] **消息通知**
  - 订单通知
  - 系统通知
  - 优惠活动

### 2. 陪玩师端功能

#### 2.1 陪玩师中心
- [ ] **申请成为陪玩师**
  - 填写基本信息
  - 选择擅长游戏
  - 上传段位证明
  - 设置服务价格
  - 提交审核
- [ ] **陪玩师资料管理**
  - 编辑个人简介
  - 更新段位信息
  - 调整价格
  - 设置接单状态（在线/离线）

#### 2.2 接单模块
- [ ] **订单大厅**
  - 待接单订单列表
  - 订单筛选（游戏、价格、时间）
  - 快速接单
- [ ] **我的接单**
  - 进行中的订单
  - 历史订单
  - 订单详情

#### 2.3 收益模块
- [ ] **收益总览**
  - 今日收益
  - 本月收益
  - 累计收益
  - 收益趋势图
- [ ] **提现功能**
  - 余额查询
  - 提现申请
  - 提现记录
  - 提现规则说明

#### 2.4 数据统计
- [ ] **个人数据**
  - 接单量统计
  - 好评率
  - 平均评分
  - 用户复购率

---

## 🔌 后端 API 规划

### API 路由结构

```
/api/v1/
├── user/              # 用户端 API
│   ├── games/         # 游戏相关
│   ├── players/       # 陪玩师相关
│   ├── orders/        # 订单相关
│   ├── payments/      # 支付相关
│   ├── reviews/       # 评价相关
│   └── profile/       # 个人中心
│
├── player/            # 陪玩师端 API
│   ├── profile/       # 陪玩师资料
│   ├── orders/        # 陪玩师订单
│   ├── earnings/      # 收益管理
│   └── stats/         # 数据统计
│
├── auth/              # 认证相关（已有）
│   ├── login
│   ├── register
│   ├── refresh
│   └── logout
│
└── admin/             # 管理端 API（已完成）
```

### 详细 API 列表

#### 1. 用户端 - 游戏相关

```go
// GET /api/v1/user/games
// 获取游戏列表（支持分类筛选）
type GameListRequest struct {
    Category string `form:"category"` // 游戏分类
    Page     int    `form:"page"`
    PageSize int    `form:"pageSize"`
}

type GameListResponse struct {
    Games []GameDTO `json:"games"`
    Total int64     `json:"total"`
}

// GET /api/v1/user/games/:id
// 获取游戏详情
type GameDetailResponse struct {
    Game        GameDTO   `json:"game"`
    PlayerCount int64     `json:"playerCount"` // 该游戏的陪玩师数量
}
```

#### 2. 用户端 - 陪玩师相关

```go
// GET /api/v1/user/players
// 获取陪玩师列表（支持筛选）
type PlayerListRequest struct {
    GameID      uint64  `form:"gameId"`      // 游戏筛选
    MinPrice    int64   `form:"minPrice"`    // 最低价格（分）
    MaxPrice    int64   `form:"maxPrice"`    // 最高价格（分）
    MinRating   float32 `form:"minRating"`   // 最低评分
    OnlineOnly  bool    `form:"onlineOnly"`  // 仅在线
    SortBy      string  `form:"sortBy"`      // 排序：price/rating/orders
    Page        int     `form:"page"`
    PageSize    int     `form:"pageSize"`
}

type PlayerListResponse struct {
    Players []PlayerCardDTO `json:"players"`
    Total   int64           `json:"total"`
}

type PlayerCardDTO struct {
    ID              uint64  `json:"id"`
    UserID          uint64  `json:"userId"`
    Nickname        string  `json:"nickname"`
    AvatarURL       string  `json:"avatarUrl"`
    Bio             string  `json:"bio"`
    Rank            string  `json:"rank"`
    RatingAverage   float32 `json:"ratingAverage"`
    RatingCount     uint32  `json:"ratingCount"`
    HourlyRateCents int64   `json:"hourlyRateCents"`
    MainGame        string  `json:"mainGame"`      // 游戏名称
    IsOnline        bool    `json:"isOnline"`      // 在线状态
    OrderCount      int64   `json:"orderCount"`    // 历史订单数
}

// GET /api/v1/user/players/:id
// 获取陪玩师详情
type PlayerDetailResponse struct {
    Player  PlayerDetailDTO `json:"player"`
    Reviews []ReviewDTO     `json:"reviews"`       // 最新评价
    Stats   PlayerStatsDTO  `json:"stats"`
}

type PlayerDetailDTO struct {
    PlayerCardDTO
    Tags           []string `json:"tags"`          // 服务标签
    GoodRatio      float32  `json:"goodRatio"`     // 好评率
    AvgResponseMin int      `json:"avgResponseMin"` // 平均响应时间（分钟）
}

type PlayerStatsDTO struct {
    TotalOrders     int64 `json:"totalOrders"`
    CompletedOrders int64 `json:"completedOrders"`
    RepeatRate      float32 `json:"repeatRate"` // 复购率
}
```

#### 3. 用户端 - 订单相关

```go
// POST /api/v1/user/orders
// 创建订单
type CreateOrderRequest struct {
    PlayerID       uint64     `json:"playerId" binding:"required"`
    GameID         uint64     `json:"gameId" binding:"required"`
    Title          string     `json:"title" binding:"required,max=128"`
    Description    string     `json:"description"`
    ScheduledStart *time.Time `json:"scheduledStart" binding:"required"`
    DurationHours  float32    `json:"durationHours" binding:"required,min=1,max=24"`
}

type CreateOrderResponse struct {
    OrderID     uint64 `json:"orderId"`
    PriceCents  int64  `json:"priceCents"`
    NeedPayment bool   `json:"needPayment"`
}

// GET /api/v1/user/orders
// 获取我的订单列表
type MyOrderListRequest struct {
    Status   string `form:"status"`   // pending/confirmed/in_progress/completed/canceled
    Page     int    `form:"page"`
    PageSize int    `form:"pageSize"`
}

type MyOrderListResponse struct {
    Orders []OrderCardDTO `json:"orders"`
    Total  int64          `json:"total"`
}

type OrderCardDTO struct {
    ID              uint64      `json:"id"`
    Title           string      `json:"title"`
    PlayerNickname  string      `json:"playerNickname"`
    PlayerAvatar    string      `json:"playerAvatar"`
    GameName        string      `json:"gameName"`
    Status          OrderStatus `json:"status"`
    PriceCents      int64       `json:"priceCents"`
    ScheduledStart  *time.Time  `json:"scheduledStart"`
    CreatedAt       time.Time   `json:"createdAt"`
    CanPay          bool        `json:"canPay"`
    CanCancel       bool        `json:"canCancel"`
    CanComplete     bool        `json:"canComplete"`
    CanReview       bool        `json:"canReview"`
}

// GET /api/v1/user/orders/:id
// 获取订单详情
type OrderDetailResponse struct {
    Order       OrderDetailDTO     `json:"order"`
    Player      PlayerCardDTO      `json:"player"`
    Payment     *PaymentDTO        `json:"payment"`
    Review      *ReviewDTO         `json:"review"`
    Timeline    []OrderTimelineDTO `json:"timeline"`
}

type OrderDetailDTO struct {
    OrderCardDTO
    Description     string     `json:"description"`
    ScheduledEnd    *time.Time `json:"scheduledEnd"`
    StartedAt       *time.Time `json:"startedAt"`
    CompletedAt     *time.Time `json:"completedAt"`
    CancelReason    string     `json:"cancelReason"`
    RefundAmount    int64      `json:"refundAmount"`
    RefundReason    string     `json:"refundReason"`
}

type OrderTimelineDTO struct {
    Time    time.Time `json:"time"`
    Status  string    `json:"status"`
    Message string    `json:"message"`
}

// PUT /api/v1/user/orders/:id/cancel
// 取消订单
type CancelOrderRequest struct {
    Reason string `json:"reason" binding:"required,max=500"`
}

// PUT /api/v1/user/orders/:id/complete
// 确认完成订单（用户侧）
type CompleteOrderRequest struct {
    Confirm bool `json:"confirm"`
}
```

#### 4. 用户端 - 支付相关

```go
// POST /api/v1/user/payments
// 创建支付
type CreatePaymentRequest struct {
    OrderID uint64        `json:"orderId" binding:"required"`
    Method  PaymentMethod `json:"method" binding:"required,oneof=wechat alipay"`
}

type CreatePaymentResponse struct {
    PaymentID uint64                 `json:"paymentId"`
    PayInfo   map[string]interface{} `json:"payInfo"` // 支付参数（对接支付SDK）
}

// GET /api/v1/user/payments/:id
// 查询支付状态
type PaymentStatusResponse struct {
    PaymentID uint64        `json:"paymentId"`
    OrderID   uint64        `json:"orderId"`
    Status    PaymentStatus `json:"status"`
    PaidAt    *time.Time    `json:"paidAt"`
}

// POST /api/v1/user/payments/:id/cancel
// 取消支付
type CancelPaymentResponse struct {
    Success bool `json:"success"`
}
```

#### 5. 用户端 - 评价相关

```go
// POST /api/v1/user/reviews
// 创建评价
type CreateReviewRequest struct {
    OrderID  uint64   `json:"orderId" binding:"required"`
    Rating   int      `json:"rating" binding:"required,min=1,max=5"`
    Comment  string   `json:"comment" binding:"max=500"`
    Tags     []string `json:"tags"`          // 评价标签
    Anonymous bool    `json:"anonymous"`      // 是否匿名
}

type CreateReviewResponse struct {
    ReviewID uint64 `json:"reviewId"`
}

// GET /api/v1/user/reviews/my
// 获取我的评价列表
type MyReviewListResponse struct {
    Reviews []MyReviewDTO `json:"reviews"`
    Total   int64         `json:"total"`
}

type MyReviewDTO struct {
    ReviewDTO
    OrderTitle     string `json:"orderTitle"`
    PlayerNickname string `json:"playerNickname"`
}
```

#### 6. 用户端 - 个人中心

```go
// GET /api/v1/user/profile
// 获取个人信息
type UserProfileResponse struct {
    ID          uint64     `json:"id"`
    Name        string     `json:"name"`
    Phone       string     `json:"phone"`
    Email       string     `json:"email"`
    AvatarURL   string     `json:"avatarUrl"`
    Status      UserStatus `json:"status"`
    CreatedAt   time.Time  `json:"createdAt"`
    IsPlayer    bool       `json:"isPlayer"`    // 是否是陪玩师
    PlayerID    uint64     `json:"playerId"`    // 陪玩师ID（如果是）
}

// PUT /api/v1/user/profile
// 更新个人信息
type UpdateProfileRequest struct {
    Name      string `json:"name" binding:"required,max=64"`
    AvatarURL string `json:"avatarUrl" binding:"url"`
}

// POST /api/v1/user/profile/change-password
// 修改密码
type ChangePasswordRequest struct {
    OldPassword string `json:"oldPassword" binding:"required,min=6"`
    NewPassword string `json:"newPassword" binding:"required,min=6"`
}

// GET /api/v1/user/wallet
// 获取钱包信息
type WalletResponse struct {
    Balance      int64 `json:"balance"`      // 余额（分）
    TotalSpent   int64 `json:"totalSpent"`   // 累计消费
    TotalRecharge int64 `json:"totalRecharge"` // 累计充值
}
```

#### 7. 陪玩师端 - 陪玩师资料

```go
// POST /api/v1/player/apply
// 申请成为陪玩师
type ApplyPlayerRequest struct {
    Nickname        string   `json:"nickname" binding:"required,max=64"`
    Bio             string   `json:"bio" binding:"max=500"`
    MainGameID      uint64   `json:"mainGameId" binding:"required"`
    Rank            string   `json:"rank" binding:"required,max=32"`
    HourlyRateCents int64    `json:"hourlyRateCents" binding:"required,min=1000"`
    Tags            []string `json:"tags"`
    ProofImages     []string `json:"proofImages"` // 段位证明图片
}

type ApplyPlayerResponse struct {
    PlayerID           uint64             `json:"playerId"`
    VerificationStatus VerificationStatus `json:"verificationStatus"`
}

// GET /api/v1/player/profile
// 获取陪玩师资料
type PlayerProfileResponse struct {
    PlayerDetailDTO
    VerificationStatus VerificationStatus `json:"verificationStatus"`
}

// PUT /api/v1/player/profile
// 更新陪玩师资料
type UpdatePlayerProfileRequest struct {
    Nickname        string   `json:"nickname" binding:"required,max=64"`
    Bio             string   `json:"bio" binding:"max=500"`
    Rank            string   `json:"rank"`
    HourlyRateCents int64    `json:"hourlyRateCents" binding:"min=1000"`
    Tags            []string `json:"tags"`
}

// PUT /api/v1/player/status
// 设置在线状态
type SetPlayerStatusRequest struct {
    Online bool `json:"online"`
}
```

#### 8. 陪玩师端 - 订单管理

```go
// GET /api/v1/player/orders/available
// 获取可接订单列表（订单大厅）
type AvailableOrdersRequest struct {
    GameID   uint64 `form:"gameId"`
    Page     int    `form:"page"`
    PageSize int    `form:"pageSize"`
}

type AvailableOrdersResponse struct {
    Orders []AvailableOrderDTO `json:"orders"`
    Total  int64               `json:"total"`
}

type AvailableOrderDTO struct {
    ID             uint64     `json:"id"`
    Title          string     `json:"title"`
    Description    string     `json:"description"`
    GameName       string     `json:"gameName"`
    UserNickname   string     `json:"userNickname"`
    PriceCents     int64      `json:"priceCents"`
    ScheduledStart *time.Time `json:"scheduledStart"`
    DurationHours  float32    `json:"durationHours"`
    CreatedAt      time.Time  `json:"createdAt"`
}

// POST /api/v1/player/orders/:id/accept
// 接单
type AcceptOrderRequest struct {
    Confirm bool `json:"confirm"`
}

type AcceptOrderResponse struct {
    Success bool       `json:"success"`
    Order   OrderDTO   `json:"order"`
}

// GET /api/v1/player/orders/my
// 获取我接的订单
type MyAcceptedOrdersRequest struct {
    Status   string `form:"status"`
    Page     int    `form:"page"`
    PageSize int    `form:"pageSize"`
}

// PUT /api/v1/player/orders/:id/start
// 开始订单
type StartOrderRequest struct {
    Confirm bool `json:"confirm"`
}

// PUT /api/v1/player/orders/:id/complete
// 完成订单（陪玩师侧）
type CompleteOrderByPlayerRequest struct {
    Confirm bool `json:"confirm"`
}
```

#### 9. 陪玩师端 - 收益管理

```go
// GET /api/v1/player/earnings/summary
// 获取收益概览
type EarningsSummaryResponse struct {
    TodayEarnings   int64   `json:"todayEarnings"`   // 今日收益（分）
    MonthEarnings   int64   `json:"monthEarnings"`   // 本月收益
    TotalEarnings   int64   `json:"totalEarnings"`   // 累计收益
    AvailableBalance int64  `json:"availableBalance"` // 可提现余额
    PendingBalance  int64   `json:"pendingBalance"`  // 待结算余额
    WithdrawTotal   int64   `json:"withdrawTotal"`   // 累计提现
}

// GET /api/v1/player/earnings/trend
// 获取收益趋势
type EarningsTrendRequest struct {
    Days int `form:"days" binding:"required,min=7,max=90"` // 7/30/90天
}

type EarningsTrendResponse struct {
    Trend []DailyEarningDTO `json:"trend"`
}

type DailyEarningDTO struct {
    Date       string `json:"date"`       // YYYY-MM-DD
    Earnings   int64  `json:"earnings"`   // 当日收益
    OrderCount int    `json:"orderCount"` // 订单数
}

// POST /api/v1/player/earnings/withdraw
// 申请提现
type WithdrawRequest struct {
    AmountCents int64  `json:"amountCents" binding:"required,min=10000"` // 最低100元
    Method      string `json:"method" binding:"required,oneof=alipay wechat bank"`
    AccountInfo string `json:"accountInfo" binding:"required"` // 账号信息
}

type WithdrawResponse struct {
    WithdrawID uint64 `json:"withdrawId"`
    Status     string `json:"status"`
}

// GET /api/v1/player/earnings/withdraw-history
// 提现记录
type WithdrawHistoryResponse struct {
    Records []WithdrawRecordDTO `json:"records"`
    Total   int64               `json:"total"`
}

type WithdrawRecordDTO struct {
    ID          uint64    `json:"id"`
    AmountCents int64     `json:"amountCents"`
    Method      string    `json:"method"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"createdAt"`
    ProcessedAt *time.Time `json:"processedAt"`
}
```

#### 10. 陪玩师端 - 数据统计

```go
// GET /api/v1/player/stats
// 获取个人统计数据
type PlayerStatsResponse struct {
    TotalOrders      int64   `json:"totalOrders"`      // 总订单数
    CompletedOrders  int64   `json:"completedOrders"`  // 完成订单数
    CanceledOrders   int64   `json:"canceledOrders"`   // 取消订单数
    CompletionRate   float32 `json:"completionRate"`   // 完成率
    AverageRating    float32 `json:"averageRating"`    // 平均评分
    TotalReviews     int64   `json:"totalReviews"`     // 评价总数
    GoodReviews      int64   `json:"goodReviews"`      // 好评数（4-5星）
    GoodRatio        float32 `json:"goodRatio"`        // 好评率
    RepeatCustomers  int64   `json:"repeatCustomers"`  // 复购用户数
    RepeatRate       float32 `json:"repeatRate"`       // 复购率
    TotalWorkHours   float32 `json:"totalWorkHours"`   // 累计工作时长
    AvgOrderDuration float32 `json:"avgOrderDuration"` // 平均订单时长
}

// GET /api/v1/player/stats/recent-orders
// 获取最近订单统计
type RecentOrdersStatsRequest struct {
    Days int `form:"days" binding:"required,min=7,max=30"`
}

type RecentOrdersStatsResponse struct {
    OrdersByDay    []DailyOrderCountDTO `json:"ordersByDay"`
    OrdersByStatus map[string]int64     `json:"ordersByStatus"`
    TopGames       []GameOrderCountDTO  `json:"topGames"`
}

type DailyOrderCountDTO struct {
    Date  string `json:"date"`
    Count int64  `json:"count"`
}

type GameOrderCountDTO struct {
    GameName string `json:"gameName"`
    Count    int64  `json:"count"`
}
```

---

## 🎨 前端页面规划

### 页面路由结构

```
/                          # 首页（游戏 + 推荐陪玩师）
/games                     # 游戏列表
/games/:id                 # 游戏详情
/players                   # 陪玩师列表
/players/:id               # 陪玩师详情页
/players/:id/book          # 预约下单页

/user/                     # 用户中心
├── /orders                # 我的订单
├── /orders/:id            # 订单详情
├── /reviews               # 我的评价
├── /profile               # 个人信息
└── /wallet                # 我的钱包

/player/                   # 陪玩师中心
├── /apply                 # 申请成为陪玩师
├── /dashboard             # 陪玩师工作台
├── /orders                # 我的接单
├── /orders/available      # 订单大厅
├── /earnings              # 收益管理
├── /profile               # 陪玩师资料
└── /stats                 # 数据统计

/payment/:orderId          # 支付页面
/payment/result            # 支付结果页

/login                     # 登录（已有）
/register                  # 注册（已有）
```

### 页面详细设计

#### 1. 首页 `/`

**组件结构**:
```tsx
<HomePage>
  <Header />
  <HeroSection>
    <SearchBar />              {/* 搜索游戏/陪玩师 */}
  </HeroSection>
  
  <HotGamesSection>
    <GameCard />               {/* 游戏卡片列表 */}
  </HotGamesSection>
  
  <FeaturedPlayersSection>
    <PlayerCard />             {/* 推荐陪玩师卡片 */}
  </FeaturedPlayersSection>
  
  <Footer />
</HomePage>
```

**核心功能**:
- 搜索功能（游戏名、陪玩师昵称）
- 游戏分类快速导航
- 热门陪玩师推荐
- Banner 广告轮播

#### 2. 陪玩师列表 `/players`

**组件结构**:
```tsx
<PlayerListPage>
  <FilterSidebar>
    <GameFilter />             {/* 游戏筛选 */}
    <PriceRangeFilter />       {/* 价格区间 */}
    <RatingFilter />           {/* 评分筛选 */}
    <OnlineFilter />           {/* 在线状态 */}
  </FilterSidebar>
  
  <PlayerGrid>
    <SortBar />                {/* 排序选项 */}
    <PlayerCard />             {/* 陪玩师卡片 */}
    <Pagination />
  </PlayerGrid>
</PlayerListPage>
```

**陪玩师卡片内容**:
- 头像、昵称
- 擅长游戏、段位
- 评分、接单数
- 时薪价格
- 在线状态
- 快速预约按钮

#### 3. 陪玩师详情页 `/players/:id`

**组件结构**:
```tsx
<PlayerDetailPage>
  <PlayerHeader>
    <Avatar />
    <BasicInfo />              {/* 昵称、段位、认证 */}
    <OnlineStatus />
  </PlayerHeader>
  
  <PlayerContent>
    <LeftColumn>
      <ProfileSection>
        <Bio />                {/* 个人简介 */}
        <Skills />             {/* 技能标签 */}
        <Games />              {/* 擅长游戏 */}
      </ProfileSection>
      
      <StatsSection>
        <OrderCount />         {/* 订单统计 */}
        <GoodRatio />          {/* 好评率 */}
        <Rating />             {/* 平均评分 */}
      </StatsSection>
      
      <ReviewsSection>
        <ReviewList />         {/* 评价列表 */}
      </ReviewsSection>
    </LeftColumn>
    
    <RightColumn>
      <BookingCard>
        <PriceInfo />          {/* 价格信息 */}
        <TimeSelector />       {/* 时间选择 */}
        <DurationSelector />   {/* 时长选择 */}
        <TotalPrice />         {/* 总价计算 */}
        <BookButton />         {/* 立即预约 */}
      </BookingCard>
    </RightColumn>
  </PlayerContent>
</PlayerDetailPage>
```

#### 4. 订单列表 `/user/orders`

**组件结构**:
```tsx
<OrderListPage>
  <OrderTabs>
    <Tab label="全部" />
    <Tab label="待支付" />
    <Tab label="进行中" />
    <Tab label="已完成" />
    <Tab label="已取消" />
  </OrderTabs>
  
  <OrderList>
    <OrderCard>
      <OrderInfo />            {/* 订单信息 */}
      <PlayerInfo />           {/* 陪玩师信息 */}
      <OrderStatus />          {/* 状态标识 */}
      <ActionButtons />        {/* 操作按钮 */}
    </OrderCard>
  </OrderList>
  
  <Pagination />
</OrderListPage>
```

**订单卡片操作**:
- 待支付：去支付、取消订单
- 进行中：查看详情、联系陪玩师
- 已完成：评价、再次预约
- 已取消：查看原因

#### 5. 订单详情 `/user/orders/:id`

**组件结构**:
```tsx
<OrderDetailPage>
  <OrderStatusBar />           {/* 订单状态进度条 */}
  
  <OrderInfoCard>
    <BasicInfo />              {/* 订单基本信息 */}
    <TimeInfo />               {/* 时间信息 */}
    <PriceInfo />              {/* 价格明细 */}
  </OrderInfoCard>
  
  <PlayerInfoCard>
    <PlayerAvatar />
    <PlayerName />
    <ContactButton />          {/* 联系陪玩师 */}
  </PlayerInfoCard>
  
  <PaymentInfoCard />          {/* 支付信息（如已支付）*/}
  
  <TimelineCard>
    <OrderTimeline />          {/* 订单时间轴 */}
  </TimelineCard>
  
  <ActionBar>
    <PayButton />              {/* 根据状态显示不同按钮 */}
    <CancelButton />
    <CompleteButton />
    <ReviewButton />
  </ActionBar>
</OrderDetailPage>
```

#### 6. 支付页面 `/payment/:orderId`

**组件结构**:
```tsx
<PaymentPage>
  <OrderSummary>
    <OrderInfo />              {/* 订单摘要 */}
    <PriceBreakdown />         {/* 价格明细 */}
  </OrderSummary>
  
  <PaymentMethods>
    <MethodOption value="wechat">
      <WeChatPayIcon />
    </MethodOption>
    <MethodOption value="alipay">
      <AlipayIcon />
    </MethodOption>
  </PaymentMethods>
  
  <CouponSection>
    <CouponSelector />         {/* 优惠券选择（可选）*/}
  </CouponSection>
  
  <TotalAmount />              {/* 实付金额 */}
  
  <PayButton />                {/* 立即支付 */}
</PaymentPage>
```

#### 7. 评价页面 `/user/orders/:id/review`

**组件结构**:
```tsx
<ReviewPage>
  <OrderInfo />                {/* 订单信息 */}
  
  <RatingSection>
    <StarRating />             {/* 星级评分 */}
    <QuickTags />              {/* 快速标签 */}
  </RatingSection>
  
  <CommentSection>
    <Textarea />               {/* 评价内容 */}
    <CharCounter />            {/* 字数统计 */}
  </CommentSection>
  
  <AnonymousOption>
    <Checkbox label="匿名评价" />
  </AnonymousOption>
  
  <SubmitButton />
</ReviewPage>
```

#### 8. 陪玩师工作台 `/player/dashboard`

**组件结构**:
```tsx
<PlayerDashboard>
  <StatsOverview>
    <StatCard label="今日收益" />
    <StatCard label="本月订单" />
    <StatCard label="好评率" />
    <StatCard label="在线时长" />
  </StatsOverview>
  
  <QuickActions>
    <OnlineToggle />           {/* 在线/离线切换 */}
    <ViewOrdersButton />       {/* 查看订单 */}
    <ViewEarningsButton />     {/* 查看收益 */}
  </QuickActions>
  
  <PendingOrders>
    <OrderCard />              {/* 待接订单 */}
  </PendingOrders>
  
  <RecentActivity>
    <ActivityItem />           {/* 最近动态 */}
  </RecentActivity>
</PlayerDashboard>
```

#### 9. 订单大厅 `/player/orders/available`

**组件结构**:
```tsx
<OrderHallPage>
  <FilterBar>
    <GameFilter />
    <PriceFilter />
    <TimeFilter />
  </FilterBar>
  
  <OrderList>
    <AvailableOrderCard>
      <OrderInfo />
      <UserInfo />
      <TimeInfo />
      <PriceInfo />
      <AcceptButton />         {/* 接单按钮 */}
    </AvailableOrderCard>
  </OrderList>
  
  <AutoRefresh />              {/* 自动刷新 */}
</OrderHallPage>
```

#### 10. 收益管理 `/player/earnings`

**组件结构**:
```tsx
<EarningsPage>
  <EarningsSummary>
    <TotalEarnings />          {/* 累计收益 */}
    <AvailableBalance />       {/* 可提现余额 */}
    <PendingBalance />         {/* 待结算 */}
  </EarningsSummary>
  
  <EarningsTrend>
    <LineChart />              {/* 收益趋势图 */}
    <DateRangeSelector />      {/* 时间范围选择 */}
  </EarningsTrend>
  
  <WithdrawSection>
    <WithdrawButton />         {/* 申请提现 */}
    <WithdrawHistory />        {/* 提现记录 */}
  </WithdrawSection>
  
  <EarningsDetail>
    <EarningsTable />          {/* 收益明细表 */}
  </EarningsDetail>
</EarningsPage>
```

---

## 📊 数据模型

### 已有模型（复用）

```go
// 用户模型
type User struct {
    Base
    Phone        string
    Email        string
    PasswordHash string
    Name         string
    AvatarURL    string
    Role         Role
    Status       UserStatus
    LastLoginAt  *time.Time
    Roles        []RoleModel
}

// 陪玩师模型
type Player struct {
    Base
    UserID             uint64
    Nickname           string
    Bio                string
    Rank               string
    RatingAverage      float32
    RatingCount        uint32
    HourlyRateCents    int64
    MainGameID         uint64
    VerificationStatus VerificationStatus
}

// 订单模型
type Order struct {
    Base
    UserID            uint64
    PlayerID          uint64
    GameID            uint64
    Title             string
    Description       string
    Status            OrderStatus
    PriceCents        int64
    Currency          Currency
    ScheduledStart    *time.Time
    ScheduledEnd      *time.Time
    CancelReason      string
    StartedAt         *time.Time
    CompletedAt       *time.Time
    RefundAmountCents int64
    RefundReason      string
    RefundedAt        *time.Time
}

// 支付模型
type Payment struct {
    Base
    OrderID         uint64
    UserID          uint64
    Method          PaymentMethod
    AmountCents     int64
    Currency        Currency
    Status          PaymentStatus
    ProviderTradeNo string
    ProviderRaw     json.RawMessage
    PaidAt          *time.Time
    RefundedAt      *time.Time
}

// 评价模型
type Review struct {
    Base
    OrderID  uint64
    UserID   uint64
    PlayerID uint64
    Score    Rating  // 1-5
    Content  string
}

// 游戏模型
type Game struct {
    Base
    Key         string
    Name        string
    Category    string
    IconURL     string
    Description string
}
```

### 新增模型需求

```go
// 陪玩师在线状态（可以用 Redis 存储）
type PlayerOnlineStatus struct {
    PlayerID   uint64    `json:"playerId"`
    IsOnline   bool      `json:"isOnline"`
    LastActive time.Time `json:"lastActive"`
}

// 提现记录
type Withdrawal struct {
    Base
    PlayerID    uint64          `json:"playerId" gorm:"column:player_id;not null;index"`
    AmountCents int64           `json:"amountCents" gorm:"column:amount_cents"`
    Method      string          `json:"method" gorm:"size:32"`
    Status      WithdrawStatus  `json:"status" gorm:"size:32;index"`
    AccountInfo string          `json:"accountInfo" gorm:"column:account_info;type:text"`
    ProcessedAt *time.Time      `json:"processedAt" gorm:"column:processed_at"`
    Remark      string          `json:"remark,omitempty" gorm:"type:text"`
}

type WithdrawStatus string

const (
    WithdrawPending   WithdrawStatus = "pending"
    WithdrawApproved  WithdrawStatus = "approved"
    WithdrawRejected  WithdrawStatus = "rejected"
    WithdrawCompleted WithdrawStatus = "completed"
)

// 优惠券（可选，后期功能）
type Coupon struct {
    Base
    Code            string     `json:"code" gorm:"size:32;uniqueIndex"`
    DiscountType    string     `json:"discountType" gorm:"size:32"` // percent/fixed
    DiscountValue   int64      `json:"discountValue"`
    MinOrderAmount  int64      `json:"minOrderAmount"`
    MaxDiscount     int64      `json:"maxDiscount"`
    ValidFrom       time.Time  `json:"validFrom"`
    ValidUntil      time.Time  `json:"validUntil"`
    UsageLimit      int        `json:"usageLimit"`
    UsageCount      int        `json:"usageCount"`
}

// 用户优惠券关联
type UserCoupon struct {
    Base
    UserID   uint64     `json:"userId" gorm:"column:user_id;index"`
    CouponID uint64     `json:"couponId" gorm:"column:coupon_id;index"`
    UsedAt   *time.Time `json:"usedAt" gorm:"column:used_at"`
    OrderID  uint64     `json:"orderId" gorm:"column:order_id"`
}

// 陪玩师标签
type PlayerTag struct {
    Base
    PlayerID uint64 `json:"playerId" gorm:"column:player_id;index"`
    Tag      string `json:"tag" gorm:"size:32;index"`
}

// 消息通知（可选）
type Notification struct {
    Base
    UserID   uint64    `json:"userId" gorm:"column:user_id;index"`
    Type     string    `json:"type" gorm:"size:32"`
    Title    string    `json:"title" gorm:"size:128"`
    Content  string    `json:"content" gorm:"type:text"`
    IsRead   bool      `json:"isRead" gorm:"column:is_read;default:false"`
    ReadAt   *time.Time `json:"readAt" gorm:"column:read_at"`
}
```

---

## 🚀 开发阶段

### 第一阶段：基础功能（2周）

**目标**: 实现核心业务流程

#### 后端任务
- [ ] **用户端 - 陪玩师模块**
  - [ ] GET /api/v1/user/players（列表 + 筛选）
  - [ ] GET /api/v1/user/players/:id（详情）
- [ ] **用户端 - 订单模块**
  - [ ] POST /api/v1/user/orders（创建订单）
  - [ ] GET /api/v1/user/orders（我的订单）
  - [ ] GET /api/v1/user/orders/:id（订单详情）
  - [ ] PUT /api/v1/user/orders/:id/cancel（取消订单）
- [ ] **支付模块（Mock 版本）**
  - [ ] POST /api/v1/user/payments（创建支付）
  - [ ] GET /api/v1/user/payments/:id（查询状态）
- [ ] **评价模块**
  - [ ] POST /api/v1/user/reviews（创建评价）
  - [ ] GET /api/v1/user/reviews/my（我的评价）

#### 前端任务
- [ ] **陪玩师列表页** (`/players`)
- [ ] **陪玩师详情页** (`/players/:id`)
- [ ] **预约下单页** (`/players/:id/book`)
- [ ] **订单列表页** (`/user/orders`)
- [ ] **订单详情页** (`/user/orders/:id`)
- [ ] **支付页面** (`/payment/:orderId`)
- [ ] **评价页面** (`/user/orders/:id/review`)

#### 验收标准
- ✅ 用户可以浏览陪玩师并下单
- ✅ 用户可以查看订单状态
- ✅ 用户可以完成支付（Mock）
- ✅ 用户可以评价订单

---

### 第二阶段：陪玩师端（2周）

**目标**: 陪玩师可以接单和管理订单

#### 后端任务
- [ ] **陪玩师认证**
  - [ ] POST /api/v1/player/apply（申请成为陪玩师）
  - [ ] GET /api/v1/player/profile（陪玩师资料）
  - [ ] PUT /api/v1/player/profile（更新资料）
  - [ ] PUT /api/v1/player/status（在线状态）
- [ ] **陪玩师订单**
  - [ ] GET /api/v1/player/orders/available（订单大厅）
  - [ ] POST /api/v1/player/orders/:id/accept（接单）
  - [ ] GET /api/v1/player/orders/my（我的接单）
  - [ ] PUT /api/v1/player/orders/:id/start（开始订单）
  - [ ] PUT /api/v1/player/orders/:id/complete（完成订单）
- [ ] **收益管理**
  - [ ] GET /api/v1/player/earnings/summary（收益概览）
  - [ ] GET /api/v1/player/earnings/trend（收益趋势）

#### 前端任务
- [ ] **申请陪玩师页** (`/player/apply`)
- [ ] **陪玩师工作台** (`/player/dashboard`)
- [ ] **订单大厅** (`/player/orders/available`)
- [ ] **我的接单** (`/player/orders/my`)
- [ ] **收益管理** (`/player/earnings`)
- [ ] **陪玩师资料** (`/player/profile`)

#### 验收标准
- ✅ 用户可以申请成为陪玩师
- ✅ 陪玩师可以在订单大厅接单
- ✅ 陪玩师可以管理订单状态
- ✅ 陪玩师可以查看收益

---

### 第三阶段：增强功能（2周）

**目标**: 完善用户体验和高级功能

#### 后端任务
- [ ] **首页功能**
  - [ ] GET /api/v1/user/home（首页数据）
  - [ ] GET /api/v1/user/games（游戏列表）
  - [ ] GET /api/v1/user/search（搜索功能）
- [ ] **个人中心**
  - [ ] GET /api/v1/user/profile（个人信息）
  - [ ] PUT /api/v1/user/profile（更新信息）
  - [ ] POST /api/v1/user/profile/change-password（修改密码）
  - [ ] GET /api/v1/user/wallet（钱包信息）
- [ ] **提现功能**
  - [ ] POST /api/v1/player/earnings/withdraw（申请提现）
  - [ ] GET /api/v1/player/earnings/withdraw-history（提现记录）
- [ ] **数据统计**
  - [ ] GET /api/v1/player/stats（个人统计）
  - [ ] GET /api/v1/player/stats/recent-orders（最近订单）

#### 前端任务
- [ ] **首页** (`/`)
- [ ] **游戏列表** (`/games`)
- [ ] **个人中心** (`/user/profile`)
- [ ] **钱包管理** (`/user/wallet`)
- [ ] **陪玩师数据统计** (`/player/stats`)
- [ ] **搜索功能**（全局搜索）

#### 验收标准
- ✅ 首页展示热门游戏和推荐陪玩师
- ✅ 搜索功能可用
- ✅ 个人中心功能完善
- ✅ 提现功能可用
- ✅ 数据统计展示

---

### 第四阶段：优化和完善（1周）

**目标**: 性能优化、体验优化、补充文档

#### 优化任务
- [ ] **性能优化**
  - [ ] 列表接口分页优化
  - [ ] 图片 CDN 加速
  - [ ] Redis 缓存热点数据
  - [ ] SQL 查询优化
- [ ] **体验优化**
  - [ ] 骨架屏加载
  - [ ] 图片懒加载
  - [ ] 无限滚动
  - [ ] 乐观更新
- [ ] **移动端适配**
  - [ ] 响应式布局
  - [ ] 触摸手势支持
  - [ ] 移动端导航
- [ ] **文档完善**
  - [ ] API 文档（Swagger）
  - [ ] 用户使用手册
  - [ ] 开发文档更新

#### 验收标准
- ✅ 首屏加载时间 < 2s
- ✅ 移动端体验流畅
- ✅ 文档齐全

---

## 🛠️ 技术架构

### 后端技术栈

```yaml
语言: Go 1.24
框架: Gin
ORM: GORM
数据库: 
  - PostgreSQL (主库)
  - SQLite (开发)
  - Redis (缓存)
认证: JWT
API文档: Swagger
日志: Zap
监控: Prometheus (可选)
```

### 前端技术栈

```yaml
语言: TypeScript 5
框架: React 18
构建: Vite
样式: Less + Arco Design
状态: React Context
路由: React Router v6
HTTP: Fetch API
图表: ECharts (收益趋势等)
支付: 微信/支付宝 SDK
```

### 第三方服务

```yaml
支付:
  - 微信支付
  - 支付宝支付
存储:
  - 阿里云 OSS (图片)
  - CDN 加速
通知:
  - 短信服务 (订单通知)
  - 邮件服务 (可选)
实时通信:
  - WebSocket (聊天，可选)
```

---

## 📐 开发规范

### API 规范

1. **RESTful 设计**
   - GET: 查询
   - POST: 创建
   - PUT: 更新
   - DELETE: 删除

2. **统一响应格式**
```json
{
  "success": true,
  "code": 200,
  "message": "操作成功",
  "data": { }
}
```

3. **错误处理**
```json
{
  "success": false,
  "code": 400,
  "message": "参数错误",
  "errors": ["字段 xxx 不能为空"]
}
```

4. **分页格式**
```json
{
  "items": [],
  "total": 100,
  "page": 1,
  "pageSize": 20,
  "totalPages": 5
}
```

### 前端规范

1. **组件命名**: PascalCase
2. **文件命名**: 组件用 PascalCase，工具用 camelCase
3. **样式**: 使用 CSS Modules
4. **类型**: 全部使用 TypeScript
5. **状态管理**: Context + Hooks
6. **路由**: 统一在 `routes.tsx` 定义

---

## 🔐 安全考虑

### 后端安全

- [ ] **认证授权**
  - JWT Token 验证
  - 角色权限控制
  - 接口访问限流
- [ ] **数据验证**
  - 参数校验
  - SQL 注入防护
  - XSS 防护
- [ ] **支付安全**
  - 签名验证
  - 防重放攻击
  - 金额校验
- [ ] **敏感数据**
  - 密码加密存储
  - 支付信息加密
  - 日志脱敏

### 前端安全

- [ ] XSS 防护（React 自带）
- [ ] CSRF Token
- [ ] 敏感信息不存本地
- [ ] HTTPS 强制

---

## 📊 监控指标

### 业务指标

- 日活用户数（DAU）
- 订单量
- 成交额
- 转化率
- 陪玩师活跃度

### 技术指标

- API 响应时间
- 错误率
- 数据库连接数
- Redis 命中率
- 服务器负载

---

## 📝 总结

### 开发周期

| 阶段 | 时间 | 主要内容 |
|------|------|----------|
| **第一阶段** | 2周 | 用户端基础功能 |
| **第二阶段** | 2周 | 陪玩师端功能 |
| **第三阶段** | 2周 | 增强功能 |
| **第四阶段** | 1周 | 优化完善 |
| **总计** | **7周** | 用户侧完整功能 |

### 核心优先级

**P0 (必须完成)**:
- 陪玩师列表和详情
- 订单创建和管理
- 支付功能
- 陪玩师接单

**P1 (重要)**:
- 评价系统
- 收益管理
- 个人中心

**P2 (优化)**:
- 首页推荐
- 搜索功能
- 数据统计

---

**文档版本**: v1.0  
**创建时间**: 2025-10-30  
**维护人员**: GameLink Team


