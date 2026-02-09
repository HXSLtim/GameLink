# GameLink 完整技术架构文档

**版本**: 1.0
**更新日期**: 2026-02-09
**目标读者**: 技术负责人、架构师、高级开发者

---

## 目录

1. [系统架构概览](#1-系统架构概览)
2. [后端架构详解](#2-后端架构详解)
3. [管理后台架构](#3-管理后台架构)
4. [小程序/H5架构](#4-小程序h5架构)
5. [数据流与状态管理](#5-数据流与状态管理)
6. [安全机制](#6-安全机制)
7. [性能优化](#7-性能优化)
8. [部署架构](#8-部署架构)
9. [监控与日志](#9-监控与日志)

---

## 1. 系统架构概览

### 1.1 整体架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                          客户端层                                │
├──────────────────┬──────────────────┬───────────────────────────┤
│   管理后台        │   小程序/H5       │       未来扩展            │
│  (React 19)      │   (Vue 3)        │    (iOS/Android)         │
│  Port: 5173      │   Port: 5174     │                           │
└────────┬─────────┴────────┬─────────┴───────────────────────────┘
         │                  │
         └──────────┬───────┘
                    │
         ┌──────────▼──────────┐
         │     Nginx            │
         │   (反向代理/SSL)     │
         │   Port: 443/80      │
         └──────────┬──────────┘
                    │
    ┌───────────────┼───────────────┬────────────────┐
    │               │               │                │
┌───▼────┐   ┌─────▼─────┐   ┌────▼────┐   ┌─────▼─────┐
│ Go API │   │WebSocket  │   │ Cron Job│   │  Worker   │
│ (Gin)  │   │(Socket.IO)│   │ (定时)  │   │ (异步任务) │
│ :8080  │   │           │   │         │   │           │
└───┬────┘   └─────┬─────┘   └─────────┘   └───────────┘
    │              │
    └──────────┬───┘
               │
    ┌──────────▼──────────┐
    │   PostgreSQL        │
    │   (主数据库)         │
    │   Port: 5432        │
    └──────────┬──────────┘
               │
    ┌──────────▼──────────┐
    │   Redis             │
    │   (缓存/队列/会话)   │
    │   Port: 6379        │
    └─────────────────────┘
```

### 1.2 技术栈矩阵

| 层级 | 管理后台 | 小程序/H5 | 后端 |
|------|---------|----------|------|
| **框架** | React 19 | Vue 3.4 | Go 1.24 + Gin |
| **语言** | TypeScript 5.9 | TypeScript | Go |
| **构建** | Vite 7.2 | Vite | Go Build |
| **状态管理** | Zustand | Pinia | - |
| **路由** | React Router 7 | uni-app Router | Gin Router |
| **UI库** | Ant Design 6 | 自定义组件库 | - |
| **HTTP** | Axios | uni.request | Gin |
| **WebSocket** | Socket.IO Client | uni.connectSocket | Socket.IO |
| **测试** | Vitest + Playwright | - | Go Test |
| **组件数** | 20+ 通用组件 | 133 组件 | - |
| **页面数** | 40+ 页面 | 28 页面 | - |

---

## 2. 后端架构详解

### 2.1 分层架构

```
┌─────────────────────────────────────────────┐
│           Handler Layer (HTTP)              │
│  - 参数验证                                  │
│  - 路由处理                                  │
│  - 响应封装                                  │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│          Service Layer (业务逻辑)            │
│  - 业务规则                                  │
│  - 事务管理                                  │
│  - 缓存策略                                  │
│  - 外部服务调用                              │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│       Repository Layer (数据访问)            │
│  - 数据库操作                                │
│  - 缓存操作                                  │
│  - 数据转换                                  │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│          Model Layer (数据模型)              │
│  - 数据结构定义                              │
│  - 数据库映射                                │
│  - 验证规则                                  │
└─────────────────────────────────────────────┘
```

### 2.2 目录结构详解

```
api/
├── cmd/                          # 应用入口
│   └── main.go                   # 主入口
│
├── internal/                     # 内部代码
│   ├── handler/                  # HTTP 处理器 (按用户角色分组)
│   │   ├── admin/                # 管理员接口 (120+ 文件)
│   │   │   ├── user.go           # 用户管理
│   │   │   ├── order.go          # 订单管理
│   │   │   ├── player.go         # 陪玩师管理
│   │   │   ├── system.go         # 系统配置
│   │   │   └── ...               # 其他管理功能
│   │   ├── user/                 # 用户接口 (15+ 文件)
│   │   │   ├── order.go          # 订单操作
│   │   │   ├── chat.go           # 聊天功能
│   │   │   ├── wallet.go         # 钱包功能
│   │   │   └── ...               # 其他用户功能
│   │   ├── player/               # 陪玩师接口 (10+ 文件)
│   │   │   ├── service.go        # 服务管理
│   │   │   ├── review.go         # 评价管理
│   │   │   └── ...               # 其他陪玩师功能
│   │   ├── public/               # 公开接口
│   │   │   ├── player.go         # 陪玩师列表
│   │   │   └── auth.go           # 认证
│   │   └── middleware/           # 中间件
│   │       ├── jwtAuth.go        # JWT 认证
│   │       ├── cors.go           # 跨域处理
│   │       ├── rateLimiter.go    # 限流
│   │       └── signature.go      # 签名验证
│   │
│   ├── service/                  # 业务逻辑层 (57 个模块)
│   │   ├── order/                # 订单服务
│   │   │   ├── creation.go       # 订单创建
│   │   │   ├── transfer.go       # 订单流转
│   │   │   └── order.go          # 订单核心
│   │   ├── payment/              # 支付服务
│   │   │   ├── payment.go        # 支付处理
│   │   │   └── wechatProvider.go # 微信支付
│   │   ├── chat/                 # 聊天服务
│   │   │   ├── service.go        # 聊天核心
│   │   │   └── message.go        # 消息处理
│   │   ├── player/               # 陪玩师服务
│   │   ├── user/                 # 用户服务
│   │   ├── gameroom/             # 游戏房间
│   │   ├── content/              # 内容审核
│   │   ├── notification/         # 通知服务
│   │   └── ...                   # 其他服务
│   │
│   ├── repository/               # 数据访问层 (56 个模块)
│   │   ├── interfaces.go        # 接口定义
│   │   ├── chat/                 # 聊天数据
│   │   ├── order/                # 订单数据
│   │   ├── player/               # 陪玩师数据
│   │   └── ...                   # 其他数据
│   │
│   ├── model/                    # 数据模型 (67 个)
│   │   ├── user.go               # 用户模型
│   │   ├── order.go              # 订单模型
│   │   ├── player.go             # 陪玩师模型
│   │   ├── chat.go               # 聊天模型
│   │   └── ...                   # 其他模型
│   │
│   └── router/                   # 路由注册
│       ├── router.go             # 主路由
│       ├── adminRoutes.go        # 管理员路由
│       ├── userRoutes.go         # 用户路由
│       ├── playerRoutes.go       # 陪玩师路由
│       ├── publicRoutes.go       # 公开路由
│       └── services.go           # 微服务路由
│
├── pkg/                          # 公共包
│   ├── auth/                     # 认证
│   │   ├── jwt.go                # JWT 处理
│   │   └── wechat.go             # 微信认证
│   ├── config/                   # 配置
│   │   └── config.go             # 配置加载
│   ├── db/                       # 数据库
│   │   ├── postgres.go           # PostgreSQL 连接
│   │   ├── migrate.go            # 数据库迁移
│   │   └── seed.go               # 种子数据
│   ├── scheduler/                # 定时任务
│   │   ├── chatRetention.go      # 聊天保留
│   │   └── settlementScheduler.go # 结算调度
│   └── oss/                      # 对象存储
│       └── service.go            # OSS 服务
│
├── ws/                           # WebSocket
│   ├── client.go                 # 客户端管理
│   ├── hub.go                    # 连接中心
│   ├── message.go                # 消息处理
│   └── redisPubsub.go            # Redis 发布订阅
│
├── .golangci.yml                 # 代码检查配置
├── go.mod                        # Go 模块
├── go.sum                        # 依赖锁定
└── main.go                       # 应用入口
```

### 2.3 依赖注入

使用 `fx` 框架进行依赖注入：

```go
// 应用构建
func NewApp(
    lifecycle fx.Lifecycle,
    config *config.Config,
    db *gorm.DB,
    redis *redis.Client,
) *gin.Engine {
    // 创建依赖
    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userService)

    // 注册路由
    router := gin.Default()
    router.POST("/users", userHandler.Create)

    return router
}
```

### 2.4 中间件链

```
Request
  ↓
CORS (跨域处理)
  ↓
Signature (签名验证) - 可选
  ↓
RateLimiter (限流)
  ↓
JWT Auth (身份认证)
  ↓
Permission (权限检查) - 管理端
  ↓
Handler (业务处理)
  ↓
Response
```

### 2.5 WebSocket 架构

```
┌─────────────────────────────────────────────┐
│           WebSocket Hub                     │
│  - 连接管理                                  │
│  - 房间管理                                  │
│  - 消息广播                                  │
└──────────────────┬──────────────────────────┘
                   │
    ┌──────────────┼──────────────┐
    │              │              │
┌───▼────┐   ┌────▼────┐   ┌────▼────┐
│ 聊天室  │   │ 订单更新 │   │ 系统通知│
│        │   │         │   │         │
└────────┘   └─────────┘   └─────────┘
     │              │              │
     └──────────────┼──────────────┘
                    │
         ┌──────────▼──────────┐
         │   Redis Pub/Sub      │
         │   (跨实例消息同步)    │
         └─────────────────────┘
```

---

## 3. 管理后台架构

### 3.1 整体架构

```
┌─────────────────────────────────────────────┐
│            Presentation Layer                │
│  - 40+ 页面组件                              │
│  - 20+ 通用组件                              │
│  - Ant Design 6                             │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│            State Management                 │
│  - Zustand (全局状态)                        │
│  - Context (权限、主题)                      │
│  - React Hooks (本地状态)                    │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│            Business Logic                   │
│  - Custom Hooks                             │
│  - Utility Functions                        │
│  - Data Transformation                       │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│            Data Layer                        │
│  - API Client (Axios)                        │
│  - WebSocket Client                         │
│  - Local Storage                            │
└─────────────────────────────────────────────┘
```

### 3.2 目录结构

```
admin/src/
├── pages/                        # 页面 (40+)
│   ├── admin/                    # 管理端页面
│   │   ├── Dashboard/            # 仪表盘
│   │   ├── User/                 # 用户管理
│   │   ├── Role/                 # 角色管理
│   │   ├── Game/                 # 游戏管理
│   │   ├── Order/                # 订单管理
│   │   ├── Player/               # 陪玩师管理
│   │   ├── Commission/           # 佣金管理
│   │   ├── Withdraw/             # 提现管理
│   │   ├── Content/              # 内容管理
│   │   │   ├── Feeds.tsx         # 动态审核
│   │   │   ├── ChatMonitor.tsx   # 聊天监控
│   │   │   └── Reports.tsx       # 举报管理
│   │   ├── Review/               # 评价管理
│   │   │   ├── index.tsx         # 评价列表
│   │   │   ├── Moderation.tsx    # 评价审核
│   │   │   └── SensitiveWords.tsx # 敏感词
│   │   ├── Monitor/              # 监控中心
│   │   │   ├── Realtime.tsx      # 实时监控
│   │   │   ├── Analytics.tsx     # 运营分析
│   │   │   └── KPI.tsx           # KPI 仪表板
│   │   ├── VIP/                  # VIP 管理
│   │   ├── Coupon/               # 优惠券管理
│   │   ├── Activity/             # 活动管理
│   │   ├── Referral/             # 推荐管理
│   │   ├── Team/                 # 团队管理
│   │   ├── PaymentRecords/       # 支付记录
│   │   ├── Settlement/           # 结算公司
│   │   └── ...                   # 其他页面
│   ├── adminChat/                # 聊天管理
│   │   ├── rooms/                # 聊天室管理
│   │   └── records/              # 聊天记录
│   ├── biz/                      # 业务页面
│   │   └── service/              # 服务项目
│   ├── sys/                      # 系统页面
│   │   ├── setting/              # 系统设置
│   │   ├── log/                  # 审计日志
│   │   ├── menu/                 # 菜单管理
│   │   ├── permission/           # 权限管理
│   │   └── user-role/            # 用户角色
│   ├── auth/                     # 认证页面
│   │   └── Login.tsx             # 登录页
│   ├── Forbidden.tsx             # 403 页面
│   └── NotFound.tsx              # 404 页面
│
├── components/                   # 组件 (20+)
│   ├── common/                   # 通用组件
│   │   ├── StateContainer/       # 状态容器
│   │   ├── BatchActions/         # 批量操作
│   │   ├── SearchFilters/        # 搜索筛选
│   │   └── ...                   # 其他组件
│   ├── PermissionGuard/          # 权限守卫
│   │   ├── index.tsx             # 权限包装器
│   │   └── withPermission.tsx    # HOC
│   ├── SearchTable/              # 搜索表格
│   ├── PageContainer/            # 页面容器
│   ├── PermissionTree/           # 权限树
│   ├── UserSelector/             # 用户选择器
│   ├── AmountDisplay/            # 金额显示
│   ├── DataExport/               # 数据导出
│   ├── AuditLog/                 # 操作日志
│   ├── DateRangePicker/          # 日期选择
│   └── ...                       # 其他组件
│
├── api/                          # API 封装
│   ├── admin.ts                  # 管理端 API (890 行)
│   │   ├── 接口定义               │
│   │   ├── 类型定义               │
│   │   └── API 函数              │
│   ├── client.ts                 # Axios 客户端
│   └── ...                       # 其他 API
│
├── router/                       # 路由
│   ├── index.tsx                 # 路由主入口
│   ├── routes.tsx                # 路由配置
│   ├── componentMap.tsx          # 组件映射 (230 行)
│   ├── Guard.tsx                 # 路由守卫
│   └── types.ts                  # 路由类型
│
├── context/                      # Context
│   ├── AdminContext.tsx          # 管理员上下文
│   ├── ThemeContext.tsx          # 主题上下文
│   ├── permissionEvents.ts       # 权限事件
│   └── useAdmin.ts               # Admin Hook
│
├── utils/                        # 工具函数
│   ├── logger.ts                 # 日志工具
│   ├── menuPermission.ts         # 菜单权限
│   ├── dynamicRoutes.ts          # 动态路由
│   ├── permission.ts             # 权限工具
│   └── websocket/                # WebSocket 工具
│       ├── manager.ts            # WS 管理器
│       ├── types.ts              # WS 类型
│       └── messageHandler.ts     # 消息处理
│
├── layouts/                      # 布局
│   └── AdminLayout/              # 管理后台布局
│       ├── index.tsx             # 主布局
│       ├── Sidebar/              # 侧边栏
│       ├── Header/               # 顶部栏
│       └── Breadcrumb/           # 面包屑
│
├── config/                       # 配置
│   ├── adminRoutes.ts            # 路由配置 (900 行)
│   └── debug.ts                  # 调试配置
│
├── hooks/                        # Hooks
│   ├── usePermission.ts          # 权限 Hook
│   ├── useWebSocket.ts           # WebSocket Hook
│   └── ...                       # 其他 Hooks
│
├── styles/                       # 样式
│   ├── global.less               # 全局样式
│   ├── variables.less            # 样式变量
│   └── themes/                   # 主题
│
├── App.tsx                       # 应用根组件
├── main.tsx                      # 应用入口
└── vite.config.ts                # Vite 配置
```

### 3.3 动态路由系统

**基于后端菜单的动态路由生成**:

```typescript
// 1. 登录后获取菜单和权限
const { menus, permissions } = useAdmin();

// 2. 根据菜单生成路由
const dynamicRoutes = generateRoutesFromMenus(menus);

// 3. 合并静态路由和动态路由
const finalRoutes = [...staticRoutes, ...dynamicRoutes];

// 4. 渲染路由
const element = useRoutes(finalRoutes);
```

**路由配置示例**:
```typescript
{
  path: '/admin',
  element: <AdminLayout />,
  meta: { roles: ['ADMIN'], requiresAuth: true },
  children: [
    {
      path: 'sys/user',
      element: <UserPage />,
      meta: {
        title: '用户管理',
        permission: 'admin.users.list'
      }
    }
  ]
}
```

### 3.4 权限系统

**权限码格式**: `模块.资源.操作`

```
admin.users.list      # 查看用户列表
admin.users.create    # 创建用户
admin.users.delete    # 删除用户
admin.orders.view     # 查看订单
*                     # 超级管理员(所有权限)
```

**权限检查**:
```typescript
// 组件级
<PermissionGuard permission="admin.users.delete">
  <Button danger>删除</Button>
</PermissionGuard>

// Hook 级
const { hasPermission } = usePermission();
if (hasPermission('admin.users.edit')) {
  // 有权限
}

// 路由级
{
  path: 'sys/user',
  element: <UserPage />,
  meta: { permission: 'admin.users.list' }
}
```

**跨标签页权限同步**:
```typescript
// 权限变更后触发事件
localStorage.setItem('permission_change_timestamp', Date.now().toString());

// 其他标签页监听变化
window.addEventListener('storage', (e) => {
  if (e.key === 'permission_change_timestamp') {
    refreshMenus();
  }
});
```

---

## 4. 小程序/H5架构

### 4.1 整体架构

```
┌─────────────────────────────────────────────┐
│          Presentation Layer                  │
│  - 28 页面                                   │
│  - 133 组件                                  │
│  - uni-app 框架                              │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│          Business Logic Layer                │
│  - 38 Composables                            │
│  - 页面级逻辑                                 │
│  - 跨页面逻辑                                 │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│          State Management                    │
│  - Pinia Stores                              │
│  - 本地状态 (ref/reactive)                   │
│  - 持久化状态                                │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│          Data Layer                           │
│  - API 模块 (14 个)                          │
│  - WebSocket 连接                            │
│  - 本地存储                                  │
└─────────────────────────────────────────────┘
```

### 4.2 目录结构

```
app/src/
├── pages/                        # 页面 (28 个)
│   ├── auth/                     # 认证
│   │   ├── login/                # 登录
│   │   └── register/             # 注册
│   ├── index/                    # 首页
│   ├── player/                   # 陪玩师
│   │   ├── list/                 # 陪玩师列表
│   │   ├── detail/               # 陪玩师详情
│   │   ├── services/             # 服务管理
│   │   ├── orders/               # 订单管理
│   │   ├── certification/        # 认证
│   │   ├── dashboard/            # 仪表盘
│   │   └── earnings/             # 收益
│   ├── order/                    # 订单
│   │   ├── create/               # 创建订单
│   │   ├── detail/               # 订单详情
│   │   └── list/                 # 订单列表
│   ├── message/                  # 消息
│   │   ├── list/                 # 消息列表
│   │   └── chat/                 # 聊天
│   ├── profile/                  # 个人中心
│   │   ├── index/                # 个人主页
│   │   └── edit/                 # 编辑资料
│   ├── wallet/                   # 钱包
│   │   ├── index/                # 钱包首页
│   │   └── recharge/             # 充值
│   ├── review/                   # 评价
│   │   └── list/                 # 评价列表
│   ├── channel/                  # 频道
│   │   └── list/                 # 频道列表
│   ├── favorite/                 # 收藏
│   │   └── list/                 # 收藏列表
│   ├── help/                     # 帮助
│   ├── service/                  # 服务
│   └── settings/                 # 设置
│
├── components/                   # 组件 (133 个)
│   ├── gl/                       # 基础组件
│   │   ├── Button/               # 按钮
│   │   ├── Card/                 # 卡片
│   │   ├── Tag/                  # 标签
│   │   ├── Avatar/               # 头像
│   │   ├── Empty/                # 空状态
│   │   ├── Icon/                 # 图标
│   │   ├── Loading/              # 加载
│   │   ├── NavBar/               # 导航栏
│   │   ├── TabBar/               # 标签栏
│   │   └── ...                   # 其他基础组件
│   ├── gl-button.vue             # 按钮组件
│   ├── gl-card.vue               # 卡片组件
│   ├── PlayerCard.vue            # 陪玩师卡片
│   ├── OrderCard.vue             # 订单卡片
│   ├── ReviewCard.vue            # 评价卡片
│   ├── ChatMessageBubble.vue     # 聊天气泡
│   ├── ChatInputBar.vue          # 聊天输入
│   ├── PaymentMethodSelector.vue # 支付方式
│   ├── QuantitySelector.vue      # 数量选择
│   ├── SchedulePicker.vue        # 时间选择
│   ├── AmountSelector.vue        # 金额选择
│   ├── CouponSelector.vue        # 优惠券选择
│   ├── FilterPanel.vue           # 筛选面板
│   ├── SearchBar.vue             # 搜索栏
│   ├── QuickActions.vue          # 快捷操作
│   ├── ResultCard.vue            # 结果卡片
│   ├── WalletBalanceCard.vue     # 余额卡片
│   ├── TransactionItem.vue       # 交易记录
│   ├── EarningsSummaryCard.vue   # 收益汇总
│   ├── EarningsItem.vue          # 收益项
│   ├── OrderStatusCard.vue       # 订单状态
│   ├── OrderActionBar.vue        # 订单操作栏
│   └── ...                       # 其他组件
│
├── composables/                  # Composables (38 个)
│   ├── useLogin.ts               # 登录逻辑
│   ├── useRegister.ts            # 注册逻辑
│   ├── usePlayerList.ts          # 陪玩师列表
│   ├── usePlayerDetail.ts        # 陪玩师详情
│   ├── usePlayerServices.ts      # 陪玩师服务
│   ├── usePlayerOrders.ts        # 陪玩师订单
│   ├── usePlayerEarnings.ts      # 陪玩师收益
│   ├── usePlayerCertification.ts # 认证
│   ├── usePlayerDashboard.ts     # 陪玩师工作台
│   ├── useOrderCreate.ts         # 创建订单
│   ├── useOrderDetail.ts         # 订单详情
│   ├── useOrderList.ts           # 订单列表
│   ├── usePaymentResult.ts       # 支付结果
│   ├── useWallet.ts              # 钱包
│   ├── useRecharge.ts            # 充值
│   ├── useChatRoom.ts            # 聊天室
│   ├── useMessageList.ts         # 消息列表
│   ├── useCustomerService.ts     # 客服
│   ├── useFavoriteList.ts        # 收藏列表
│   ├── useGameList.ts            # 游戏列表
│   ├── useReviewList.ts          # 评价列表
│   ├── useProfile.ts             # 个人资料
│   ├── useProfileEdit.ts         # 编辑资料
│   ├── useSettings.ts            # 设置
│   ├── useTheme.ts               # 主题
│   ├── useToast.ts               # 提示
│   ├── useWebSocket.ts           # WebSocket
│   ├── useChannelList.ts         # 频道列表
│   ├── useHelp.ts                # 帮助
│   ├── usePageState.ts           # 页面状态
│   ├── usePagination.ts          # 分页
│   └── ...                       # 其他 Composables
│
├── api/                          # API 封装 (14 个模块)
│   ├── request.ts                # 请求封装
│   ├── auth.ts                   # 认证 API
│   ├── user.ts                   # 用户 API
│   ├── player.ts                 # 陪玩师 API
│   ├── order.ts                  # 订单 API
│   ├── payment.ts                # 支付 API
│   ├── wallet.ts                 # 钱包 API
│   ├── chat.ts                   # 聊天 API
│   ├── review.ts                 # 评价 API
│   ├── upload.ts                 # 上传 API
│   └── ...                       # 其他 API
│
├── store/                        # Pinia Store
│   ├── index.ts                  # Store 入口
│   └── user.ts                   # 用户状态
│
├── styles/                       # 样式
│   ├── index.scss                # 样式入口
│   ├── variables.scss            # 样式变量
│   └── mixins.scss               # 样式混合
│
├── data/                         # 数据
│   ├── agreements.ts            # 协议文本
│   └── ...                       # 其他数据
│
├── types/                        # 类型定义
│   └── ...                       # 类型文件
│
├── App.vue                       # 应用根组件
├── main.ts                       # 应用入口
├── pages.json                    # 页面配置
├── manifest.json                 # 应用配置
├── uni.scss                      # uni-app 样式
└── vite.config.ts                # Vite 配置
```

### 4.3 组件体系

**三层组件架构**:

1. **基础组件 (gl/)**: 通用 UI 组件，无业务逻辑
   - gl-button, gl-card, gl-tag, gl-avatar
   - 可在任何页面使用

2. **业务组件**: 包含业务逻辑的组件
   - PlayerCard, OrderCard, ReviewCard
   - 特定业务场景使用

3. **模式组件**: 复杂的交互模式
   - FilterPanel, SearchBar, QuickActions
   - 跨页面复用

### 4.4 Composables 模式

**职责**: 封装可复用的业务逻辑

```typescript
// 示例: usePlayerList
export function usePlayerList() {
  const players = ref<Player[]>([]);
  const loading = ref(false);
  const hasMore = ref(true);

  const loadPlayers = async (page: number) => {
    loading.value = true;
    try {
      const res = await playerApi.getPlayers({ page });
      players.value.push(...res.data);
      hasMore.value = res.data.length > 0;
    } finally {
      loading.value = false;
    }
  };

  const loadMore = () => {
    if (!loading.value && hasMore.value) {
      loadPlayers(players.value.length / 20 + 1);
    }
  };

  return { players, loading, hasMore, loadMore };
}
```

---

## 5. 数据流与状态管理

### 5.1 后端数据流

```
Client Request
  ↓
[Middleware Chain]
  ↓
Handler (参数验证)
  ↓
Service (业务逻辑)
  ↓
Repository (数据访问)
  ↓
Database (PostgreSQL)
  ↓
Cache (Redis)
  ↓
Response
```

### 5.2 前端状态管理

**管理后台 (Zustand + Context)**:

```
┌─────────────────────────────────────┐
│      Zustand Stores (全局状态)       │
│  - 用户状态                          │
│  - 主题状态                          │
└──────────────────┬──────────────────┘
                   │
┌──────────────────▼──────────────────┐
│      Context (权限/菜单)             │
│  - AdminContext                      │
│  - ThemeContext                      │
└──────────────────┬──────────────────┘
                   │
┌──────────────────▼──────────────────┐
│      React Hooks (本地状态)          │
│  - useState                          │
│  - useEffect                         │
│  - useContext                        │
└─────────────────────────────────────┘
```

**小程序/H5 (Pinia)**:

```
┌─────────────────────────────────────┐
│      Pinia Stores (全局状态)         │
│  - user store                        │
│  - app store                         │
└──────────────────┬──────────────────┘
                   │
┌──────────────────▼──────────────────┐
│      Composables (业务逻辑)          │
│  - usePlayerList                     │
│  - useOrderCreate                    │
└──────────────────┬──────────────────┘
                   │
┌──────────────────▼──────────────────┐
│      Vue 3 Reactivity (本地状态)     │
│  - ref                               │
│  - reactive                          │
│  - computed                          │
└─────────────────────────────────────┘
```

### 5.3 数据同步策略

**服务端 → 客户端**:
1. **轮询**: 定时拉取数据 (低频数据)
2. **WebSocket**: 实时推送 (高频数据)
3. **事件驱动**: 用户操作触发更新

**客户端 → 服务端**:
1. **直接调用**: API 请求
2. **乐观更新**: 先更新 UI，后同步服务端
3. **队列缓冲**: 批量提交

---

## 6. 安全机制

### 6.1 认证

**JWT Token 流程**:
```
1. 用户登录
2. 验证密码
3. 生成 JWT Token
   - Header: 算法类型
   - Payload: 用户ID、角色、过期时间
   - Signature: 签名
4. 返回 Token 给客户端
5. 客户端存储 (localStorage)
6. 后续请求携带 Token
7. 服务端验证 Token
```

**Token 刷新**:
```typescript
// Access Token: 15 分钟
// Refresh Token: 7 天

// Access Token 过期
// 1. 使用 Refresh Token 获取新的 Access Token
// 2. 更新本地存储
// 3. 重试原请求
```

### 6.2 授权

**RBAC 模型**:
```
User (用户)
  ↓
Role (角色: user/player/admin)
  ↓
Permission (权限: admin.users.list)
  ↓
Resource (资源: 用户、订单、陪玩师)
  ↓
Action (操作: list/create/update/delete)
```

**权限检查**:
```go
// 后端中间件
func RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        user := getCurrentUser(c)
        if !hasPermission(user, permission) {
            c.JSON(403, gin.H{"error": "Permission denied"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 6.3 数据加密

**敏感字段加密**:
```go
// 使用 AES-256 加密
type User struct {
    ID           int    `json:"id"`
    Name         string `json:"name"`
    PhoneEncrypted string `json:"-" gorm:"column:phone_encrypted"`
    IDCardEncrypted string `json:"-" gorm:"column:idcard_encrypted"`
}

// 加密工具
func Encrypt(text string) (string, error) {
    // AES-256-GCM 加密
}
```

**传输加密**:
- HTTPS (TLS 1.3)
- API 签名验证 (管理端)
- WebSocket 加密 (WSS)

### 6.4 防护措施

**限流**:
```go
// 基于 Redis 的滑动窗口限流
func RateLimiter(maxRequests int, duration time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := fmt.Sprintf("rate_limit:%s", c.ClientIP())
        count, _ := redis.Incr(key)
        if count == 1 {
            redis.Expire(key, duration)
        }
        if count > maxRequests {
            c.JSON(429, gin.H{"error": "Too many requests"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

**SQL 注入防护**:
- 使用参数化查询 (GORM)
- 输入验证

**XSS 防护**:
- React 自动转义
- Content Security Policy

**CSRF 防护**:
- SameSite Cookie
- CSRF Token (管理端)

---

## 7. 性能优化

### 7.1 后端优化

**数据库优化**:
1. **索引优化**
   ```sql
   CREATE INDEX idx_orders_user_id ON orders(user_id);
   CREATE INDEX idx_orders_status ON orders(status);
   CREATE INDEX idx_orders_created_at ON orders(created_at DESC);
   ```

2. **查询优化**
   ```go
   // 使用预加载 (Eager Loading)
   db.Preload("Player").Preload("Game").Find(&orders)

   // 只查询需要的字段
   db.Select("id", "name", "avatar_url").Find(&users)

   // 分页查询
   db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&data)
   ```

3. **连接池配置**
   ```go
   db.DB().SetMaxIdleConns(10)
   db.DB().SetMaxOpenConns(100)
   db.DB().SetConnMaxLifetime(time.Hour)
   ```

**缓存策略**:
```go
// 多级缓存
1. 查询 Redis
2. Redis 未命中，查询 PostgreSQL
3. 回写 Redis (设置过期时间)

// 缓存失效
- 数据更新时删除缓存
- 设置合理的过期时间
- 使用缓存预热
```

**异步处理**:
```go
// 使用 Goroutine 处理耗时任务
go func() {
    // 发送通知
    sendNotification(order)
    // 更新统计
    updateStats(order)
}()
```

### 7.2 前端优化

**代码分割**:
```typescript
// 路由懒加载
const Dashboard = lazy(() => import('@/pages/admin/Dashboard'));

// 组件懒加载
const HeavyComponent = lazy(() => import('./HeavyComponent'));
```

**资源优化**:
```typescript
// Vite 配置
export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          'react-vendor': ['react', 'react-dom'],
          'antd-vendor': ['antd', '@ant-design/icons'],
          'charts-vendor': ['recharts'],
        }
      }
    }
  }
});
```

**图片优化**:
```typescript
// 使用 WebP 格式
// 响应式图片
// 懒加载
<img loading="lazy" src="..." />
```

**请求优化**:
```typescript
// 请求合并
// 请求缓存
// 防抖/节流
const debouncedSearch = debounce(handleSearch, 300);
```

### 7.3 小程序优化

**分包加载**:
```json
{
  "subPackages": [
    {
      "root": "pages/player",
      "pages": ["list", "detail"]
    }
  ]
}
```

**资源优化**:
- 图片压缩
- 代码压缩
- 减少包大小

---

## 8. 部署架构

### 8.1 开发环境

```
┌─────────────────────────────────────┐
│         Developer Machine            │
│                                      │
│  ┌────────┐  ┌────────┐  ┌────────┐ │
│  │  Go    │  │ React  │  │ Vue   │ │
│  │ :8080  │  │ :5173  │  │ :5174  │ │
│  └────────┘  └────────┘  └────────┘ │
│       │            │           │     │
└───────┼────────────┼───────────┼─────┘
        │            │           │
    ┌───▼────────────▼───────────▼────┐
    │         Docker Compose           │
    │  ┌──────────┐  ┌─────────────┐  │
    │  │PostgreSQL│  │    Redis    │  │
    │  │  :5432   │  │    :6379    │  │
    │  └──────────┘  └─────────────┘  │
    └─────────────────────────────────┘
```

### 8.2 生产环境

```
                    ┌─────────────┐
                    │   用户       │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │   CDN       │
                    │  (静态资源)  │
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
         ┌────▼────┐  ┌────▼────┐  ┌──▼─────┐
         │ Nginx 1 │  │ Nginx 2 │  │ Nginx 3│
         │ (负载均衡)│  │ (负载均衡)│  │ (负载均衡)│
         └────┬────┘  └────┬────┘  └───┬────┘
              │            │            │
              └────────┬───┴────────────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
    ┌───▼────┐   ┌────▼────┐   ┌────▼────┐
    │ API 1  │   │ API 2   │   │ API 3   │
    │ :8080  │   │ :8080   │   │ :8080   │
    └───┬────┘   └────┬────┘   └────┬────┘
        │              │              │
        └──────────────┼──────────────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
    ┌───▼─────────┐ ┌─▼──────────┐ ┌─▼──────────┐
    │ PostgreSQL  │ │   Redis    │ │  OSS       │
    │ (主从复制)  │ │  (集群)    │ │ (对象存储)  │
    └─────────────┘ └────────────┘ └────────────┘
```

### 8.3 Docker 部署

**后端 Dockerfile**:
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o api cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api .
CMD ["./api"]
```

**前端 Dockerfile**:
```dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
```

---

## 9. 监控与日志

### 9.1 日志系统

**后端日志**:
```
api/logs/
├── app.log           # 应用日志
├── error.log         # 错误日志
├── access.log        # 访问日志
└── sql.log           # SQL 日志
```

**日志级别**:
```go
logger.Debug("详细调试信息")
logger.Info("一般信息")
logger.Warn("警告信息")
logger.Error("错误信息")
logger.Fatal("致命错误")
```

**前端日志**:
```typescript
// 开发环境
console.log('Debug info');

// 生产环境
// 发送到日志服务
logger.log('User action', { action: 'click', target: 'button' });
```

### 9.2 监控指标

**系统指标**:
- CPU 使用率
- 内存使用率
- 磁盘 I/O
- 网络流量

**应用指标**:
- 请求响应时间
- 请求成功率
- 错误率
- 并发连接数

**业务指标**:
- 活跃用户数
- 订单量
- 支付成功率
- 聊天消息数

### 9.3 告警机制

**告警规则**:
```yaml
- name: HighErrorRate
  condition: error_rate > 5%
  duration: 5m
  action: 发送钉钉通知

- name: SlowResponse
  condition: avg_response_time > 2s
  duration: 10m
  action: 发送邮件

- name: HighCPU
  condition: cpu_usage > 80%
  duration: 15m
  action: 发送短信
```

---

## 10. 总结

### 10.1 架构亮点

1. **分层清晰**: Handler → Service → Repository → Model
2. **权限完善**: RBAC + 动态路由 + 组件级权限
3. **组件复用**: 133 个小程序组件，20+ 管理后台组件
4. **实时通信**: WebSocket + Redis Pub/Sub
5. **性能优化**: 多级缓存 + 代码分割 + 懒加载
6. **安全可靠**: JWT + HTTPS + 数据加密 + 限流
7. **易于扩展**: 微服务架构 + Docker 容器化

### 10.2 技术债务

1. **测试覆盖率**: 需要提高到 80%+
2. **文档完善**: API 文档、组件文档
3. **性能监控**: 需要完整的监控系统
4. **错误追踪**: 需要集成 Sentry
5. **自动化测试**: 需要 E2E 测试

### 10.3 未来规划

1. **微服务拆分**: 按业务领域拆分服务
2. **消息队列**: 引入 Kafka/RabbitMQ
3. **搜索引擎**: 集成 Elasticsearch
4. **分布式追踪**: Jaeger/Zipkin
5. **服务网格**: Istio
6. **CI/CD**: 完善 Jenkins/GitHub Actions

---

**文档维护**: 技术负责人
**更新频率**: 每季度更新一次
**反馈渠道**: GitHub Issues
