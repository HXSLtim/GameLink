# GameLink 改进项目 - 文件结构规划

**本文档展示**: 所有需要新增和修改的文件  
**使用方式**: 作为开发检查清单使用

---

## 📁 后端文件结构

### 新增数据模型 (6个文件)

```
backend/internal/model/
├── ✨ dispute.go          # 争议/投诉系统 (约150行)
├── ✨ ticket.go           # 客服工单系统 (约180行)
├── ✨ notification.go     # 站内通知系统 (约80行)
├── ✨ chat.go             # 聊天消息系统 (约60行)
├── ✨ favorite.go         # 用户收藏系统 (约40行)
└── ✨ tag.go              # 陪玩师标签系统 (约70行)
```

### 修改现有模型 (3个文件)

```
backend/internal/model/
├── 🔧 user.go            # 添加10+字段和关联
├── 🔧 player.go          # 添加8+字段和关联
└── 🔧 order.go           # 添加6+字段和关联
```

### Repository 层 (6个目录, 12个文件)

```
backend/internal/repository/
├── ✨ dispute/
│   ├── repository.go       # 接口实现 (约200行)
│   └── repository_test.go  # 单元测试 (约150行)
├── ✨ ticket/
│   ├── repository.go       # 接口实现 (约200行)
│   └── repository_test.go  # 单元测试 (约150行)
├── ✨ notification/
│   ├── repository.go       # 接口实现 (约150行)
│   └── repository_test.go  # 单元测试 (约100行)
├── ✨ chat/
│   ├── repository.go       # 接口实现 (约150行)
│   └── repository_test.go  # 单元测试 (约100行)
├── ✨ favorite/
│   ├── repository.go       # 接口实现 (约120行)
│   └── repository_test.go  # 单元测试 (约80行)
└── ✨ tag/
    ├── repository.go       # 接口实现 (约120行)
    └── repository_test.go  # 单元测试 (约80行)
```

### Service 层 (8个目录, 16+个文件)

```
backend/internal/service/
├── ✨ dispute/
│   ├── service.go          # 业务逻辑 (约300行)
│   └── service_test.go     # 单元测试 (约200行)
├── ✨ ticket/
│   ├── service.go          # 业务逻辑 (约300行)
│   └── service_test.go     # 单元测试 (约200行)
├── ✨ notification/
│   ├── service.go          # 业务逻辑 (约200行)
│   └── service_test.go     # 单元测试 (约150行)
├── ✨ chat/
│   ├── service.go          # 业务逻辑 (约200行)
│   ├── hub.go              # WebSocket Hub (约150行)
│   └── service_test.go     # 单元测试 (约150行)
├── ✨ favorite/
│   ├── service.go          # 业务逻辑 (约150行)
│   └── service_test.go     # 单元测试 (约100行)
├── ✨ upload/
│   ├── service.go          # 文件上传 (约200行)
│   ├── local.go            # 本地存储 (约100行)
│   ├── oss.go              # OSS存储 (约150行)
│   └── service_test.go     # 单元测试 (约150行)
└── payment/
    ├── 🔧 payment.go       # 现有文件增强
    ├── ✨ alipay.go        # 支付宝集成 (约300行)
    └── ✨ wechat.go        # 微信支付集成 (约300行)
```

### Handler 层 (15+个文件)

```
backend/internal/handler/
├── user/
│   ├── ✨ dispute.go          # 争议处理 (约200行)
│   ├── ✨ ticket.go           # 工单系统 (约200行)
│   ├── ✨ notification.go     # 通知管理 (约150行)
│   └── ✨ favorite.go         # 收藏管理 (约120行)
├── player/
│   └── ✨ online.go           # 在线状态 (约100行)
├── ✨ websocket/
│   ├── chat.go                # 聊天WebSocket (约200行)
│   └── notification.go        # 通知WebSocket (约150行)
└── ✨ upload/
    └── upload.go              # 文件上传 (约150行)
```

### 调度器和中间件 (4个文件)

```
backend/internal/
├── scheduler/
│   ├── 🔧 settlement_scheduler.go  # 现有文件
│   └── ✨ order_scheduler.go       # 订单调度 (约200行)
└── middleware/
    └── ✨ prometheus.go            # Prometheus监控 (约150行)
```

---

## 📁 前端文件结构

### 用户端页面 (7个目录, 14个文件)

```
frontend/src/pages/UserPortal/
├── ✨ Home/
│   ├── index.tsx              # 用户首页 (约300行)
│   └── Home.module.less       # 样式文件 (约200行)
├── ✨ GameList/
│   ├── index.tsx              # 游戏列表 (约250行)
│   └── GameList.module.less   # 样式文件 (约150行)
├── ✨ PlayerList/
│   ├── index.tsx              # 陪玩师列表 (约400行)
│   └── PlayerList.module.less # 样式文件 (约250行)
├── ✨ PlayerDetail/
│   ├── index.tsx              # 陪玩师详情 (约400行)
│   └── PlayerDetail.module.less # 样式文件 (约300行)
├── ✨ OrderCreate/
│   ├── index.tsx              # 创建订单 (约350行)
│   └── OrderCreate.module.less # 样式文件 (约200行)
├── ✨ MyOrders/
│   ├── index.tsx              # 我的订单 (约450行)
│   └── MyOrders.module.less   # 样式文件 (约300行)
└── ✨ Profile/
    ├── index.tsx              # 个人中心 (约500行)
    └── Profile.module.less    # 样式文件 (约350行)
```

### 陪玩师端页面 (7个目录, 14个文件)

```
frontend/src/pages/PlayerPortal/
├── ✨ Dashboard/
│   ├── index.tsx              # 工作台 (约400行)
│   └── Dashboard.module.less  # 样式文件 (约300行)
├── ✨ Orders/
│   ├── index.tsx              # 订单管理 (约450行)
│   └── Orders.module.less     # 样式文件 (约300行)
├── ✨ Earnings/
│   ├── index.tsx              # 收益管理 (约500行)
│   └── Earnings.module.less   # 样式文件 (约350行)
├── ✨ Services/
│   ├── index.tsx              # 服务管理 (约350行)
│   └── Services.module.less   # 样式文件 (约250行)
├── ✨ Profile/
│   ├── index.tsx              # 资料管理 (约400行)
│   └── Profile.module.less    # 样式文件 (约300行)
├── ✨ Reviews/
│   ├── index.tsx              # 评价管理 (约300行)
│   └── Reviews.module.less    # 样式文件 (约200行)
└── ✨ Schedule/
    ├── index.tsx              # 时间管理 (约350行)
    └── Schedule.module.less   # 样式文件 (约250行)
```

### 通用组件 (8个目录, 24个文件)

```
frontend/src/components/
├── ✨ GameCard/
│   ├── index.ts
│   ├── GameCard.tsx           # 游戏卡片 (约150行)
│   └── GameCard.module.less   # 样式 (约100行)
├── ✨ PlayerCard/
│   ├── index.ts
│   ├── PlayerCard.tsx         # 陪玩师卡片 (约200行)
│   └── PlayerCard.module.less # 样式 (约150行)
├── ✨ OrderStatusBadge/
│   ├── index.ts
│   ├── OrderStatusBadge.tsx   # 订单状态 (约100行)
│   └── OrderStatusBadge.module.less # 样式 (约80行)
├── ✨ ChatWindow/
│   ├── index.ts
│   ├── ChatWindow.tsx         # 聊天窗口 (约500行)
│   └── ChatWindow.module.less # 样式 (约400行)
├── ✨ DisputeModal/
│   ├── index.ts
│   ├── DisputeModal.tsx       # 争议弹窗 (约300行)
│   └── DisputeModal.module.less # 样式 (约200行)
├── ✨ TicketModal/
│   ├── index.ts
│   ├── TicketModal.tsx        # 工单弹窗 (约300行)
│   └── TicketModal.module.less # 样式 (约200行)
├── ✨ NotificationBell/
│   ├── index.ts
│   ├── NotificationBell.tsx   # 通知铃铛 (约200行)
│   └── NotificationBell.module.less # 样式 (约150行)
└── ✨ FavoriteButton/
    ├── index.ts
    ├── FavoriteButton.tsx     # 收藏按钮 (约150行)
    └── FavoriteButton.module.less # 样式 (约100行)
```

### 前端服务层 (8个文件)

```
frontend/src/services/
├── api/
│   ├── ✨ dispute.ts          # 争议API (约100行)
│   ├── ✨ ticket.ts           # 工单API (约120行)
│   ├── ✨ notification.ts     # 通知API (约100行)
│   ├── ✨ favorite.ts         # 收藏API (约80行)
│   ├── ✨ chat.ts             # 聊天API (约80行)
│   └── ✨ earnings.ts         # 收益API (约100行)
└── ✨ websocket/
    └── chat.ts                 # WebSocket服务 (约200行)
```

### 前端类型定义 (6个文件)

```
frontend/src/types/
├── ✨ dispute.ts              # 争议类型 (约80行)
├── ✨ ticket.ts               # 工单类型 (约100行)
├── ✨ notification.ts         # 通知类型 (约50行)
├── ✨ favorite.ts             # 收藏类型 (约40行)
├── ✨ chat.ts                 # 聊天类型 (约60行)
└── 🔧 player.ts               # 陪玩师类型增强 (新增约50行)
```

---

## 📊 文件统计

### 后端统计

```
新增文件: 45个
├── Model层: 6个 (约580行代码)
├── Repository层: 12个 (约2,260行代码 含测试)
├── Service层: 16个 (约4,000行代码 含测试)
├── Handler层: 8个 (约1,800行代码)
├── WebSocket: 2个 (约350行代码)
└── 其他: 3个 (约550行代码)

修改文件: 3个
├── user.go (新增约50行)
├── player.go (新增约40行)
└── order.go (新增约30行)

总计代码量: 约9,000行
├── 业务代码: 约6,500行
└── 测试代码: 约2,500行
```

### 前端统计

```
新增文件: 80个
├── 用户端页面: 14个 (约4,500行)
├── 陪玩师端页面: 14个 (约5,000行)
├── 通用组件: 24个 (约3,500行)
├── 服务层: 8个 (约880行)
├── 类型定义: 6个 (约480行)
└── 其他: 14个 (约1,000行)

修改文件: 5个
├── router/index.tsx (新增路由)
├── types/player.ts (类型增强)
└── 其他配置文件

总计代码量: 约15,000行
├── 业务代码: 约12,000行
├── 样式代码: 约2,500行
└── 测试代码: 约500行
```

### 总计

```
总文件数: 125+个
总代码行数: 24,000+行
├── 后端: 9,000行
└── 前端: 15,000行

预估工时: 1,160小时
├── 后端开发: 480小时
├── 前端开发: 480小时
├── 测试: 80小时
├── 项目管理: 120小时
```

---

## 🎯 开发优先级

### P0 - 第一优先级 (Week 1-2)

**后端**:
```
✅ dispute.go          # 争议系统核心
✅ ticket.go           # 工单系统核心
✅ notification.go     # 通知系统核心
✅ 对应的Repository和Service层
✅ 支付系统真实集成
```

### P1 - 第二优先级 (Week 3-4)

**前端**:
```
✅ 用户端7个页面      # 核心用户体验
✅ 陪玩师端7个页面    # 核心陪玩师体验
✅ GameCard组件
✅ PlayerCard组件
```

### P2 - 第三优先级 (Week 5)

**通用功能**:
```
✅ ChatWindow组件     # 聊天功能
✅ WebSocket服务      # 实时通信
✅ NotificationBell   # 通知推送
✅ 争议和工单前端页面
```

### P3 - 第四优先级 (Week 5-6)

**优化和增强**:
```
✅ favorite.go和tag.go    # 收藏和标签
✅ 其他通用组件
✅ 性能优化
✅ 测试补充
```

---

## 📋 每周文件创建清单

### Week 1: 后端数据层
```
Day 1:
├── ✅ dispute.go
├── ✅ ticket.go
├── ✅ notification.go
├── ✅ chat.go
├── ✅ favorite.go
└── ✅ tag.go

Day 2:
├── ✅ user.go (修改)
├── ✅ player.go (修改)
├── ✅ order.go (修改)
└── ✅ 数据库迁移

Day 3-4:
├── ✅ 6个Repository文件
└── ✅ 6个Repository测试文件

Day 5-7:
├── ✅ 8个Service文件
└── ✅ 8个Service测试文件
```

### Week 2: 后端API层
```
Day 8-10:
├── ✅ 8个Handler文件
├── ✅ 2个WebSocket文件
└── ✅ upload.go

Day 11-12:
├── ✅ alipay.go
├── ✅ wechat.go
└── ✅ 支付测试

Day 13-14:
├── ✅ order_scheduler.go
├── ✅ prometheus.go
└── ✅ API文档更新
```

### Week 3: 用户端前端
```
Day 15-21:
├── ✅ Home/ (2个文件)
├── ✅ GameList/ (2个文件)
├── ✅ PlayerList/ (2个文件)
├── ✅ PlayerDetail/ (2个文件)
├── ✅ OrderCreate/ (2个文件)
├── ✅ MyOrders/ (2个文件)
├── ✅ Profile/ (2个文件)
└── ✅ 相关API服务和类型定义
```

### Week 4: 陪玩师端前端
```
Day 22-28:
├── ✅ Dashboard/ (2个文件)
├── ✅ Orders/ (2个文件)
├── ✅ Earnings/ (2个文件)
├── ✅ Services/ (2个文件)
├── ✅ Profile/ (2个文件)
├── ✅ Reviews/ (2个文件)
├── ✅ Schedule/ (2个文件)
└── ✅ 相关API服务和类型定义
```

### Week 5: 通用功能
```
Day 29-35:
├── ✅ GameCard/ (3个文件)
├── ✅ PlayerCard/ (3个文件)
├── ✅ ChatWindow/ (3个文件)
├── ✅ NotificationBell/ (3个文件)
├── ✅ DisputeModal/ (3个文件)
├── ✨ TicketModal/ (3个文件)
├── ✅ OrderStatusBadge/ (3个文件)
├── ✅ FavoriteButton/ (3个文件)
└── ✅ WebSocket服务集成
```

### Week 6: 测试和优化
```
Day 36-42:
├── ✅ 所有单元测试文件
├── ✅ 集成测试文件
├── ✅ E2E测试文件
└── ✅ 文档更新
```

---

## 🔍 文件详细信息

### 数据模型文件详情

#### dispute.go
```
文件: backend/internal/model/dispute.go
代码行数: ~150行
主要内容:
  - DisputeStatus枚举 (4个状态)
  - DisputeType枚举 (3个类型)
  - Dispute结构体 (20+字段)
  - 关联关系定义
依赖:
  - model.Base
  - model.User
  - model.Order
测试覆盖: 需要覆盖所有字段和方法
```

#### ticket.go
```
文件: backend/internal/model/ticket.go
代码行数: ~180行
主要内容:
  - TicketStatus枚举 (4个状态)
  - TicketPriority枚举 (4个优先级)
  - TicketCategory枚举 (5个分类)
  - Ticket结构体 (15+字段)
  - TicketMessage结构体 (6+字段)
  - 关联关系定义
依赖:
  - model.Base
  - model.User
测试覆盖: 需要覆盖所有字段和方法
```

### 核心页面详情

#### UserPortal/Home/index.tsx
```
文件: frontend/src/pages/UserPortal/Home/index.tsx
代码行数: ~300行
主要功能:
  - Hero Banner展示
  - 热门游戏展示
  - 推荐陪玩师展示
  - 快速下单入口
依赖组件:
  - GameCard
  - PlayerCard
  - Button
API依赖:
  - gameApi.list()
  - playerApi.getRecommended()
```

#### PlayerPortal/Dashboard/index.tsx
```
文件: frontend/src/pages/PlayerPortal/Dashboard/index.tsx
代码行数: ~400行
主要功能:
  - 今日数据统计
  - 在线状态切换
  - 待接单列表
  - 进行中订单
  - 快捷操作
依赖组件:
  - Card
  - Tabs
  - OrderCard
API依赖:
  - playerApi.getStats()
  - playerApi.updateOnlineStatus()
  - orderApi.getPendingOrders()
```

---

## ✅ 验收标准

### 每个文件的验收标准

#### 数据模型文件
- [x] 包含完整的字段定义
- [x] GORM标签正确
- [x] JSON标签使用camelCase
- [x] 关联关系定义完整
- [x] 包含必要的方法
- [x] 有单元测试

#### Repository文件
- [x] 接口定义完整
- [x] 实现所有CRUD方法
- [x] Context作为第一参数
- [x] 错误处理完善
- [x] 有单元测试
- [x] 测试覆盖率 >= 80%

#### Service文件
- [x] 业务逻辑清晰
- [x] 错误包装完善
- [x] 依赖注入正确
- [x] 日志记录完整
- [x] 有单元测试
- [x] 测试覆盖率 >= 80%

#### Handler文件
- [x] HTTP处理正确
- [x] 参数验证完整
- [x] 错误响应统一
- [x] Swagger文档完整
- [x] 有单元测试

#### 前端页面文件
- [x] TypeScript类型完整
- [x] 组件结构清晰
- [x] API对接正确
- [x] 错误处理完善
- [x] 响应式设计
- [x] 有组件测试

#### 前端组件文件
- [x] Props类型定义
- [x] 组件复用性好
- [x] 样式模块化
- [x] 有单元测试
- [x] 文档注释完整

---

## 🚀 开始开发

### 第一步: 准备工作
```bash
# 1. 查看总索引
cat README_IMPROVEMENT.md

# 2. 查看执行摘要
cat IMPROVEMENT_SUMMARY.md

# 3. 运行结构搭建脚本
.\scripts\setup-improvement-structure.ps1
```

### 第二步: 查看详细计划
```bash
# 查看完整方案
cat GAMELINK_IMPROVEMENT_PLAN.md

# 或者分章节查看
# 第1章: 数据模型
# 第2章: 后端API
# 第3章: 前端页面
# 等等
```

### 第三步: 开始编码
```bash
# 从第一个文件开始
cd backend/internal/model
code dispute.go

# 参考代码模板
# 在GAMELINK_IMPROVEMENT_PLAN.md中搜索 "dispute.go"
```

---

## 📞 支持

需要帮助? 查看:
1. README_IMPROVEMENT.md - 总索引
2. IMPROVEMENT_GUIDE.md - 使用指南
3. GAMELINK_IMPROVEMENT_PLAN.md - 详细方案

**让我们开始这个激动人心的改进之旅! 🚀**

---

**文档版本**: v1.0  
**创建时间**: 2025年11月7日  
**文件总数**: 125+个  
**代码总量**: 24,000+行

