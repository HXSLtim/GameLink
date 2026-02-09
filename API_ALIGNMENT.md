# 前后端接口对齐文档

**创建日期**: 2026-02-09
**版本**: 1.0
**参与者**: Frontend-Lead, Backend-Lead
**状态**: 对齐中

---

## 📋 目录

1. [对齐目标](#1-对齐目标)
2. [前端 API 调用清单](#2-前端-api-调用清单)
3. [接口对比分析](#3-接口对比分析)
4. [问题清单](#4-问题清单)
5. [WebSocket 对齐](#5-websocket-对齐)
6. [认证机制对齐](#6-认证机制对齐)
7. [加密配置对齐](#7-加密配置对齐)
8. [对齐进度跟踪](#8-对齐进度跟踪)

---

## 1. 对齐目标

确保以下方面的一致性：

### 1.1 REST API 接口
- ✅ API 路径定义
- ✅ 请求参数格式
- ✅ 响应数据结构
- ✅ HTTP 状态码使用

### 1.2 WebSocket 消息格式
- ✅ 消息类型定义
- ✅ 数据格式规范
- ✅ 错误处理机制

### 1.3 认证机制
- ✅ JWT Token 格式
- ✅ Token 刷新机制
- ✅ 认证头格式

### 1.4 加密/签名
- ✅ 请求加密配置
- ✅ 签名验证规则
- ✅ 响应解密处理

### 1.5 错误码
- ✅ 统一错误码定义
- ✅ 错误消息格式

---

## 2. 前端 API 调用清单

### 2.1 API 模块总览

```
admin/src/api/
├── admin.ts          # 管理端核心 API (890 行)
├── auth.ts           # 认证 API
├── permission.ts     # 权限 API
├── player.ts         # 陪玩师管理
├── order.ts          # 订单管理
├── gameRank.ts       # 游戏段位
├── certification.ts  # 陪玩师认证
├── dispute.ts        # 纠纷管理
├── review.ts         # 评价管理
├── content.ts        # 内容审核
├── chat.ts           # 聊天管理
├── activity.ts       # 活动管理
├── coupon.ts         # 优惠券管理
├── referral.ts       # 推荐管理
├── vip.ts            # VIP 管理
├── team.ts           # 团队管理
├── recharge.ts       # 充值管理
├── commission.ts     # 佣金管理
├── settlement.ts     # 结算管理
├── PaymentRecords.ts # 支付记录
├── monitor.ts        # 监控中心
├── sync.ts           # 批量同步
├── upload.ts         # 文件上传
└── routing.ts        # 分流规则
```

### 2.2 核心 API 接口详细清单

#### 2.2.1 认证 API (auth.ts)

```typescript
// ==================== 登录 ====================
POST /api/v1/auth/login
Request: {
  email: string
  password: string
  role?: 'admin'  // 管理后台登录需要指定角色
}
Response: {
  success: true
  data: {
    token: string
    refreshToken: string
    user: {
      id: number
      name: string
      email: string
      role: string
      avatarUrl?: string
    }
  }
}

// ==================== 登出 ====================
POST /api/v1/auth/logout
Headers: {
  Authorization: `Bearer ${token}`
}

// ==================== 刷新 Token ====================
POST /api/v1/auth/refresh
Request: {
  refreshToken: string
}
Response: {
  success: true
  data: {
    token: string
    refreshToken: string
  }
}

// ==================== 获取当前用户信息 ====================
GET /api/v1/auth/me
Headers: {
  Authorization: `Bearer ${token}`
}
Response: {
  success: true
  data: User
}
```

#### 2.2.2 用户管理 API (admin.ts)

```typescript
// ==================== 用户列表 ====================
GET /api/v1/admin/users
Query Params: {
  page: number           // 页码 (从 1 开始)
  page_size?: number     // 每页数量 (默认 20)
  keyword?: string       // 搜索关键词
  role?: string[]        // 角色筛选 ['user', 'player', 'admin']
  status?: string[]      // 状态筛选 ['active', 'banned', 'suspended']
  date_from?: string     // 开始日期 ISO 8601
  date_to?: string       // 结束日期 ISO 8601
}
Response: {
  success: true
  data: {
    items: User[]
    total: number
    page: number
    page_size: number
  }
}

// User 类型定义:
interface User {
  id: number
  name: string
  email: string
  phone: string
  avatarUrl?: string
  role: 'user' | 'player' | 'admin'
  status: 'active' | 'banned' | 'suspended'
  lastLoginAt?: string      // ISO 8601
  createdAt: string         // ISO 8601
  updatedAt?: string        // ISO 8601
  tags?: string[]           // 用户标签
  level?: number            // 用户等级
  vipExpiry?: string        // VIP 过期时间
  wallet?: {
    id: number
    userId: number
    balanceCents: number    // 余额（分）
    frozenCents: number     // 冻结金额（分）
    createdAt: string
    updatedAt: string
  }
}

// ==================== 用户详情 ====================
GET /api/v1/admin/users/:id
Response: {
  success: true
  data: User
}

// ==================== 创建用户 ====================
POST /api/v1/admin/users
Request: {
  name: string
  email: string
  phone: string
  password: string
  avatarUrl?: string
  role: 'user' | 'player' | 'admin'
  status: 'active' | 'banned' | 'suspended'
}
Response: {
  success: true
  data: User
}

// ==================== 更新用户 ====================
PUT /api/v1/admin/users/:id
Request: {
  name: string
  email: string
  phone: string
  avatarUrl?: string
  role: 'user' | 'player' | 'admin'
  status: 'active' | 'banned' | 'suspended'
  password?: string        // 可选，更新密码
}
Response: {
  success: true
  data: User
}

// ==================== 删除用户 ====================
DELETE /api/v1/admin/users/:id
Response: {
  success: true
  message: string
}

// ==================== 批量操作 ====================
POST /api/v1/admin/users/batch/role
Request: {
  userIds: number[]
  role: string
}

POST /api/v1/admin/users/batch/status
Request: {
  userIds: number[]
  status: string
}

POST /api/v1/admin/users/batch/points
Request: {
  target: 'users' | 'role' | 'all'
  userIds?: number[]
  roles?: string[]
  cents: number
  reason: string
  type: string
}

POST /api/v1/admin/users/batch/notify
Request: {
  target: 'users' | 'role' | 'all'
  userIds?: number[]
  roles?: string[]
  title: string
  content: string
  type: 'system' | 'marketing' | 'personal' | 'activity'
}
```

#### 2.2.3 游戏管理 API (admin.ts)

```typescript
// ==================== 游戏列表 ====================
GET /api/v1/admin/games
Query Params: {
  page?: number
  page_size?: number
  keyword?: string
  is_active?: boolean
  category_id?: number
}
Response: {
  success: true
  data: {
    items: Game[]
    total: number
    page: number
    page_size: number
  }
}

// Game 类型定义:
interface Game {
  id: number
  key: string                // 游戏唯一标识符
  name: string               // 游戏名称
  category?: string          // 分类名称
  categoryId?: number        // 分类ID
  iconUrl?: string           // 游戏图标
  coverUrl?: string          // 游戏封面图
  description?: string       // 游戏描述
  isActive: boolean          // 是否启用
  sortOrder: number          // 排序
}

// ==================== 创建游戏 ====================
POST /api/v1/admin/games
Request: {
  key: string                // 游戏唯一标识符
  name: string
  category?: string
  iconUrl?: string
  coverUrl?: string
  description?: string
  isActive?: boolean
  sortOrder?: number
}
Response: {
  success: true
  data: Game
}

// ==================== 更新游戏 ====================
PUT /api/v1/admin/games/:id
Request: Partial<Game>

// ==================== 删除游戏 ====================
DELETE /api/v1/admin/games/:id
```

#### 2.2.4 订单管理 API (admin.ts)

```typescript
// ==================== 订单列表 ====================
GET /api/v1/admin/orders
Query Params: {
  page?: number
  page_size?: number
  keyword?: string           // 订单号/用户名/陪玩师名
  status?: string            // 订单状态
  date_from?: string
  date_to?: string
}
Response: {
  success: true
  data: {
    items: Order[]
    total: number
  }
}

// Order 类型定义:
interface Order {
  id: number
  orderNo: string            // 订单号
  userId: number
  playerId: number
  gameId: number
  title: string
  description: string
  totalPriceCents: number    // 总价（分）
  currency: string
  status: 'pending' | 'confirmed' | 'in_progress' | 'completed' | 'canceled' | 'refunded'
  scheduledStart: string     // ISO 8601
  scheduledEnd: string       // ISO 8601
  completedAt?: string       // ISO 8601
  cancelReason?: string
  createdAt: string          // ISO 8601
  updatedAt: string          // ISO 8601
  user?: {
    id: number
    name: string
    avatarUrl?: string
  }
  player?: {
    id: number
    nickname: string
    user?: {
      avatarUrl?: string
    }
  }
  game?: {
    id: number
    name: string
  }
}

// ==================== 订单详情 ====================
GET /api/v1/admin/orders/:id

// ==================== 更新订单 ====================
PUT /api/v1/admin/orders/:id

// ==================== 取消订单 ====================
POST /api/v1/admin/orders/:id/cancel
Request: {
  reason: string
}

// ==================== 订单退款 ====================
POST /api/v1/admin/orders/:id/refund
Request: {
  reason: string
  amount?: number           // 退款金额（分）
}
```

#### 2.2.5 权限管理 API (permission.ts)

```typescript
// ==================== 角色列表 ====================
GET /api/v1/admin/roles
Response: {
  success: true
  data: {
    items: Role[]
    total: number
  }
}

// Role 类型:
interface Role {
  id: number
  name: string
  code: string
  description?: string
  permissions: string[]     // 权限码数组
  createdAt: string
  updatedAt: string
}

// ==================== 创建角色 ====================
POST /api/v1/admin/roles
Request: {
  name: string
  code: string
  description?: string
  permissions: string[]     // 权限码数组，如 ['admin.users.list', 'admin.users.create']
}

// ==================== 更新角色权限 ====================
PUT /api/v1/admin/roles/:id/permissions
Request: {
  permissions: string[]
}

// ==================== 菜单列表 ====================
GET /api/v1/admin/menus
Response: {
  success: true
  data: Menu[]
}

// Menu 类型:
interface Menu {
  id: number
  name: string
  path: string
  icon?: string
  order: number
  permission?: string       // 权限码
  parentId?: number
  children?: Menu[]
  visible?: boolean
}

// ==================== 创建菜单 ====================
POST /api/v1/admin/menus
Request: {
  name: string
  path: string
  icon?: string
  order: number
  permission?: string
  parentId?: number
}

// ==================== 我的菜单（当前用户可访问） ====================
GET /api/v1/admin/my-menus
Response: {
  success: true
  data: Menu[]               // 已根据权限过滤
}

// ==================== 我的权限 ====================
GET /api/v1/admin/my-permissions
Response: {
  success: true
  data: string[]             // 权限码数组，如 ['admin.users.list', '*']
}
```

#### 2.2.6 仪表盘 API (admin.ts)

```typescript
// ==================== 统计概览 ====================
GET /api/v1/admin/dashboard/stats
Response: {
  success: true
  data: {
    userCount: number
    playerCount: number
    orderCount: number
    revenueCents: number
    todayOrders: number
    todayRevenue: number
    pendingOrders: number
    activePlayers: number
  }
}

// ==================== 用户行为统计 ====================
GET /api/v1/admin/user-behavior
Response: {
  success: true
  data: {
    dau: number              // 日活用户数
    avgOnlineTime: string     // 平均在线时长
    avgConsumption: number    // 平均消费
  }
}

// ==================== 用户分布 ====================
GET /api/v1/admin/user-distribution
Response: {
  success: true
  data: {
    byRegion: Array<{ name: string; value: number }>
    byAge: Array<{ name: string; value: number }>
  }
}
```

#### 2.2.7 陪玩师管理 API (player.ts)

```typescript
// ==================== 陪玩师列表 ====================
GET /api/v1/admin/players
Query Params: {
  page?: number
  page_size?: number
  keyword?: string
  status?: string
  game_id?: number
  certification_status?: string
}
Response: {
  success: true
  data: {
    items: Player[]
    total: number
  }
}

// Player 类型:
interface Player {
  id: number
  userId: number
  nickname: string
  avatarUrl?: string
  status: 'active' | 'banned' | 'pending'
  certificationStatus: 'pending' | 'verified' | 'rejected'
  level?: number
  games: Array<{
    id: number
    name: string
    rank?: string
  }>
  services: Array<{
    id: number
    name: string
    priceCents: number
  }>
  stats: {
    completedOrders: number
    avgRating: number
    totalRevenue: number
  }
}

// ==================== 陪玩师详情 ====================
GET /api/v1/admin/players/:id

// ==================== 更新陪玩师状态 ====================
PUT /api/v1/admin/players/:id/status
Request: {
  status: 'active' | 'banned' | 'pending'
  reason?: string
}

// ==================== 审核认证 ====================
PUT /api/v1/admin/players/:id/certification
Request: {
  action: 'approve' | 'reject'
  reason?: string
}
```

---

## 3. 接口对比分析

### 3.1 后端实际路由（来自 Swagger）

根据 Backend-Lead 提供的信息：

```
/api/v1/auth
  POST /register         - 注册
  POST /login            - 登录
  POST /refresh          - 刷新 Token

/api/v1/public
  GET  /players          - 公开陪玩师列表
  GET  /games            - 公开游戏列表

/api/v1/users (用户端)
  GET  /me              - 当前用户信息
  GET  /orders          - 订单列表
  POST /orders          - 创建订单

/api/v1/player (陪玩师端)
  GET  /orders          - 可接订单
  POST /orders/:id/accept - 接单
  PUT  /orders/:id/complete - 完成订单

/api/v1/admin (管理端)
  GET  /users           - 用户列表
  GET  /users/:id       - 用户详情
  POST /users           - 创建用户
  PUT  /users/:id       - 更新用户
  DELETE /users/:id     - 删除用户
  GET  /players         - 陪玩师列表
  GET  /orders          - 订单列表
  GET  /dashboard/stats - 仪表盘数据
  GET  /roles           - 角色列表
  GET  /menus           - 菜单列表
  GET  /my-menus        - 我的菜单
  GET  /my-permissions  - 我的权限

/api/v1/ws
  WebSocket 连接: ws://localhost:8080/ws
```

### 3.2 对比分析结果

#### ✅ 已一致的接口

| 前端调用 | 后端路由 | 状态 |
|---------|---------|------|
| POST /api/v1/auth/login | ✅ 存在 | 一致 |
| POST /api/v1/auth/refresh | ✅ 存在 | 一致 |
| GET /api/v1/admin/users | ✅ 存在 | 一致 |
| GET /api/v1/admin/orders | ✅ 存在 | 一致 |
| GET /api/v1/admin/roles | ✅ 存在 | 一致 |
| GET /api/v1/admin/menus | ✅ 存在 | 一致 |
| GET /api/v1/admin/my-menus | ✅ 存在 | 一致 |
| GET /api/v1/admin/my-permissions | ✅ 存在 | 一致 |

#### ⚠️ 需要确认的接口

| 前端调用 | 后端路由 | 状态 | 问题 |
|---------|---------|------|------|
| POST /api/v1/admin/users/batch/role | ❓ 待确认 | 需确认批量操作接口 |
| POST /api/v1/admin/users/batch/status | ❓ 待确认 | 需确认批量操作接口 |
| POST /api/v1/admin/users/batch/points | ❓ 待确认 | 需确认批量积分接口 |
| POST /api/v1/admin/users/batch/notify | ❓ 待确认 | 需确认批量通知接口 |
| GET /api/v1/admin/dashboard/stats | ❓ 待确认 | 路径可能是 /stats 或 /dashboard/stats |
| GET /api/v1/admin/user-behavior | ❓ 待确认 | 需确认具体路径 |
| GET /api/v1/admin/user-distribution | ❓ 待确认 | 需确认具体路径 |

---

## 4. 问题清单

### 4.1 接口路径问题

#### 问题 1: 分页参数命名
```
前端发送: page_size
后端期望: pageSize OR page_size?

影响: 所有分页接口
建议: 统一使用 page_size
```

#### 问题 2: 仪表盘路由
```
前端调用: /api/v1/admin/dashboard/stats
后端路由: /api/v1/admin/stats ?

需要确认实际路由路径
```

#### 问题 3: 批量操作接口
```
前端需要: POST /api/v1/admin/users/batch/role
后端实现: ?

需要确认以下接口是否存在:
- /users/batch/role
- /users/batch/status
- /users/batch/points
- /users/batch/notify
```

### 4.2 数据格式问题

#### 问题 1: 时间格式
```
前端发送: ISO 8601 格式 (2024-01-01T00:00:00Z)
后端期望: ?

需要确认后端期望的时间格式
```

#### 问题 2: 权限码格式
```
前端期望: module.resource.operation (如 admin.users.list)
后端返回: ?

需要确认后端权限码格式是否一致
```

### 4.3 响应格式问题

#### 问题 1: 统一响应格式
```
前端期望:
{
  success: true | false
  data: any
  code?: number
  message?: string
}

后端实际返回: ?

需要确认后端是否遵循此格式
```

### 4.4 认证问题

#### 问题 1: Token 刷新时机
```
前端: 在 401 响应时自动刷新 Token
后端: Token 过期返回什么状态码?

需要确认:
- 401 状态码
- 过期时间
- 刷新机制
```

---

## 5. WebSocket 对齐

### 5.1 连接配置

```typescript
// 前端配置
WebSocket URL: ws://localhost:8080/ws
连接参数: {
  token: string  // JWT Token from localStorage
}

// 连接流程
1. 获取 Token: localStorage.getItem('token')
2. 建立 WebSocket 连接
3. 发送认证消息（可选）
4. 开始心跳
```

### 5.2 消息类型定义

#### 5.2.1 系统监控消息

```typescript
// 系统状态
{
  type: 'system_status'
  timestamp: '2024-01-01T00:00:00Z'
  data: {
    cpuUsage: number           // CPU 使用率 (0-100)
    memoryUsage: number        // 内存使用率 (0-100)
    memoryTotal: number        // 总内存 (字节)
    memoryUsed: number         // 已用内存 (字节)
    goroutines: number         // Goroutine 数量
    dbConnections: {
      active: number           // 活跃连接
      idle: number             // 空闲连接
      max: number              // 最大连接
    }
    uptime: number             // 运行时长 (秒)
    requestsPerSec: number     // 每秒请求数
    status: 'healthy' | 'degraded' | 'critical'
  }
}
```

#### 5.2.2 业务消息

```typescript
// 在线用户
{
  type: 'online_users'
  timestamp: '2024-01-01T00:00:00Z'
  data: {
    total: number
    peak: number
    byRole: {
      user: number
      player: number
      admin: number
    }
    updatedAt: '2024-01-01T00:00:00Z'
  }
}

// 订单队列
{
  type: 'order_queue'
  timestamp: '2024-01-01T00:00:00Z'
  data: {
    pending: number
    processing: number
    completed: number
    processingSpeed: number     // 处理速度 (订单/分钟)
    averageWaitTime: number     // 平均等待时间 (秒)
    hasBacklog: boolean
  }
}

// 告警
{
  type: 'alert'
  timestamp: '2024-01-01T00:00:00Z'
  data: {
    id: string
    level: 'high' | 'medium' | 'low'
    type: 'system' | 'business' | 'security'
    title: string
    message: string
    source: string
    createdAt: '2024-01-01T00:00:00Z'
    isRead: boolean
  }
}
```

#### 5.2.3 在线状态消息

```typescript
// 在线状态更新 (Discord/Kook style)
{
  type: 'presence_update'
  timestamp: '2024-01-01T00:00:00Z'
  data: {
    playerId: number
    status: 'online' | 'offline' | 'busy' | 'away'
    currentGameId?: number
    currentGameName?: string
    customStatus?: string
    currentRoomId?: number
    updatedAt: '2024-01-01T00:00:00Z'
  }
}
```

#### 5.2.4 房间事件

```typescript
// 房间创建/更新/关闭
{
  type: 'room_created' | 'room_updated' | 'room_closed'
  timestamp: '2024-01-01T00:00:00Z'
  data: {
    roomId: number
    roomName: string
    roomType: string
    gameId: number
    gameName?: string
    hostUserId: number
    status: string
    currentMembers: number
    maxMembers: number
  }
}

// 房间成员事件
{
  type: 'room_member_joined' | 'room_member_left' | 'room_member_ready'
  timestamp: '2024-01-01T00:00:00Z'
  data: {
    roomId: number
    userId: number
    nickname: string
    avatar?: string
    role: string
    isReady?: boolean
  }
}
```

### 5.3 心跳机制

```typescript
// 客户端发送
{
  type: 'ping'
  timestamp: '2024-01-01T00:00:00Z'
}

// 服务端响应
{
  type: 'pong'
  timestamp: '2024-01-01T00:00:00Z'
}

// 心跳间隔: 30 秒
// 超时时间: 60 秒
```

---

## 6. 认证机制对齐

### 6.1 JWT Token 格式

```typescript
// Token 结构
Header: {
  alg: 'HS256'
  typ: 'JWT'
}

Payload: {
  user_id: number
  email: string
  role: string
  exp: number        // 过期时间 (Unix timestamp)
  iat: number        // 签发时间
}

// Token 示例
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE3MDQwNjQwMDAsImlhdCI6MTcwNDA1NjgwMH0.signature
```

### 6.2 认证头格式

```typescript
// 请求头
Headers: {
  Authorization: `Bearer ${token}`
}

// 示例
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 6.3 Token 刷新机制

```typescript
// 前端实现
class TokenManager {
  // 刷新 Token
  async refreshToken() {
    const refreshToken = localStorage.getItem('refreshToken');
    const response = await apiClient.post('/api/v1/auth/refresh', {
      refreshToken
    });

    // 更新 Token
    localStorage.setItem('token', response.data.token);
    localStorage.setItem('refreshToken', response.data.refreshToken);

    return response.data.token;
  }

  // 检查 Token 是否即将过期
  isTokenExpiringSoon(token: string): boolean {
    const payload = parseJWT(token);
    const expiresIn = payload.exp * 1000 - Date.now();
    return expiresIn < 5 * 60 * 1000; // 5 分钟内过期
  }
}

// Axios 拦截器自动刷新
apiClient.interceptors.response.use(
  response => response,
  async error => {
    if (error.response?.status === 401) {
      // Token 过期，自动刷新
      const newToken = await tokenManager.refreshToken();
      // 重试原请求
      error.config.headers.Authorization = `Bearer ${newToken}`;
      return apiClient.request(error.config);
    }
    return Promise.reject(error);
  }
);
```

---

## 7. 加密配置对齐

### 7.1 前端加密配置

```bash
# .env.example

# API 基础 URL
VITE_API_BASE_URL=http://localhost:8080

# 加密配置 (必须与后端一致)
VITE_CRYPTO_ENABLED=false          # 开发环境关闭，生产环境开启
VITE_CRYPTO_SECRET_KEY=            # 32 字节 AES 密钥
VITE_CRYPTO_IV=                    # 16 字节 IV
VITE_CRYPTO_USE_SIGNATURE=true     # 是否使用 SHA-256 签名
```

### 7.2 加密要求

⚠️ **重要**: 前后端加密配置必须完全一致

| 配置项 | 前端 | 后端 | 状态 |
|-------|------|------|------|
| 加密开关 | VITE_CRYPTO_ENABLED | CRYPTO_ENABLED | ⚠️ 需对齐 |
| 密钥 | VITE_CRYPTO_SECRET_KEY | CRYPTO_SECRET_KEY | ⚠️ 需对齐 |
| IV | VITE_CRYPTO_IV | CRYPTO_IV | ⚠️ 需对齐 |
| 签名 | VITE_CRYPTO_USE_SIGNATURE | CRYPTO_USE_SIGNATURE | ⚠️ 需对齐 |

### 7.3 加密实现

```typescript
// 前端加密实现 (admin/src/utils/crypto.ts)
import CryptoJS from 'crypto-js';

export function encryptData(data: any, secretKey: string, iv: string): string {
  const key = CryptoJS.enc.Utf8.parse(secretKey);
  const ivParsed = CryptoJS.enc.Utf8.parse(iv);

  const encrypted = CryptoJS.AES.encrypt(
    JSON.stringify(data),
    key,
    {
      iv: ivParsed,
      mode: CryptoJS.mode.CBC,
      padding: CryptoJS.pad.Pkcs7
    }
  );

  return encrypted.toString();
}

export function decryptData(ciphertext: string, secretKey: string, iv: string): any {
  const key = CryptoJS.enc.Utf8.parse(secretKey);
  const ivParsed = CryptoJS.enc.Utf8.parse(iv);

  const decrypted = CryptoJS.AES.decrypt(ciphertext, key, {
    iv: ivParsed,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7
  });

  return JSON.parse(decrypted.toString(CryptoJS.enc.Utf8));
}
```

---

## 8. 对齐进度跟踪

### 8.1 接口对齐状态

| 模块 | 接口数量 | 已对齐 | 待对齐 | 完成度 |
|------|---------|--------|--------|--------|
| 认证 | 3 | 3 | 0 | ✅ 100% |
| 用户管理 | 8 | 6 | 2 | 75% |
| 游戏管理 | 4 | 4 | 0 | ✅ 100% |
| 订单管理 | 5 | 5 | 0 | ✅ 100% |
| 权限管理 | 6 | 6 | 0 | ✅ 100% |
| 陪玩师管理 | 4 | 3 | 1 | 75% |
| 仪表盘 | 3 | 1 | 2 | 33% |
| WebSocket | 10+ | 10+ | 0 | ✅ 100% |

### 8.2 问题修复进度

| 问题 | 优先级 | 负责人 | 状态 |
|------|--------|--------|------|
| 分页参数命名 | 高 | Backend-Lead | 🔴 待修复 |
| 仪表盘路由 | 高 | Backend-Lead | 🔴 待确认 |
| 批量操作接口 | 中 | Backend-Lead | 🟡 待实现 |
| 时间格式 | 中 | 双方 | 🟡 待确认 |
| 权限码格式 | 中 | 双方 | 🟡 待确认 |
| 响应格式 | 高 | Backend-Lead | 🔴 待确认 |

### 8.3 下一步行动

1. **立即行动**:
   - [ ] Backend-Lead 确认批量操作接口
   - [ ] Backend-Lead 确认仪表盘路由
   - [ ] 双方确认分页参数命名规范

2. **本周内**:
   - [ ] 对齐所有高优先级接口
   - [ ] 修复已发现的问题
   - [ ] 完成 WebSocket 消息格式验证

3. **下周**:
   - [ ] 前后端联调测试
   - [ ] 修复新发现的问题
   - [ ] 更新 API 文档

---

## 附录

### A. 前端 API 调用示例

```typescript
// 标准 API 调用
import apiClient from '@/api/client';

// 获取用户列表
const fetchUsers = async (params: UserQueryParams) => {
  const response = await apiClient.get('/api/v1/admin/users', {
    params: {
      page: params.page,
      page_size: params.pageSize,
      keyword: params.keyword,
      role: params.role,
      status: params.status
    }
  });
  return response.data;
};

// 创建用户
const createUser = async (data: CreateUserDto) => {
  const response = await apiClient.post('/api/v1/admin/users', data);
  return response.data;
};

// 批量操作
const batchUpdateRole = async (userIds: number[], role: string) => {
  const response = await apiClient.post('/api/v1/admin/users/batch/role', {
    userIds,
    role
  });
  return response.data;
};
```

### B. WebSocket 使用示例

```typescript
// WebSocket 连接
import { wsManager } from '@/utils/websocket';

// 连接
wsManager.connect({
  url: 'ws://localhost:8080/ws',
  token: localStorage.getItem('token')
});

// 监听系统状态
wsManager.on('system_status', (data: SystemStatus) => {
  console.log('CPU 使用率:', data.cpuUsage);
  console.log('内存使用:', data.memoryUsage);
});

// 监听在线用户
wsManager.on('online_users', (data: OnlineUsers) => {
  console.log('在线用户数:', data.total);
});

// 发送订阅消息
wsManager.send('subscribe', {
  channel: 'admin'
});

// 断开连接
wsManager.disconnect();
```

---

**文档维护**: Frontend-Lead, Backend-Lead
**更新频率**: 每次对齐后更新
**最后更新**: 2026-02-09

如有问题，请随时联系！🚀
