# GameLink 前端 API 接口需求文档

**版本**: v1.0.0  
**更新时间**: 2025-01-05  
**对应后端模型**: User, Player, Order, Game, Payment, Review

---

## 📋 目录

1. [认证接口](#1-认证接口)
2. [用户管理接口](#2-用户管理接口)
3. [陪玩师管理接口](#3-陪玩师管理接口)
4. [订单管理接口](#4-订单管理接口)
5. [游戏管理接口](#5-游戏管理接口)
6. [支付管理接口](#6-支付管理接口)
7. [评价管理接口](#7-评价管理接口)
8. [数据统计接口](#8-数据统计接口)
9. [权限管理接口](#9-权限管理接口)
10. [系统设置接口](#10-系统设置接口)

---

## 1. 认证接口

### 1.1 用户登录

```http
POST /api/auth/login
```

**请求体**:

```typescript
{
  username: string; // 手机号或邮箱
  password: string;
}
```

**响应**:

```typescript
{
  token: string;
  user: {
    id: number;
    name: string;
    email?: string;
    phone?: string;
    avatar_url?: string;
    role: 'user' | 'player' | 'admin';
    status: 'active' | 'suspended' | 'banned';
  };
}
```

### 1.2 用户登出

```http
POST /api/auth/logout
```

**Headers**: `Authorization: Bearer {token}`

**响应**: `204 No Content`

### 1.3 刷新令牌

```http
POST /api/auth/refresh
```

**Headers**: `Authorization: Bearer {token}`

**响应**:

```typescript
{
  token: string;
}
```

### 1.4 获取当前用户信息

```http
GET /api/auth/me
```

**Headers**: `Authorization: Bearer {token}`

**响应**:

```typescript
{
  id: number;
  name: string;
  email?: string;
  phone?: string;
  avatar_url?: string;
  role: 'user' | 'player' | 'admin';
  status: 'active' | 'suspended' | 'banned';
  last_login_at?: string;
}
```

---

## 2. 用户管理接口

### 2.1 获取用户列表

```http
GET /api/users
```

**Query 参数**:

```typescript
{
  page?: number;          // 页码，默认 1
  page_size?: number;     // 每页数量，默认 10
  keyword?: string;       // 搜索关键词（姓名/手机/邮箱）
  role?: 'user' | 'player' | 'admin';  // 角色筛选
  status?: 'active' | 'suspended' | 'banned';  // 状态筛选
  created_start?: string; // 注册开始时间 ISO8601
  created_end?: string;   // 注册结束时间 ISO8601
  sort_by?: 'created_at' | 'last_login_at';  // 排序字段
  sort_order?: 'asc' | 'desc';  // 排序方向
}
```

**响应**:

```typescript
{
  list: Array<{
    id: number;
    name: string;
    phone?: string;
    email?: string;
    avatar_url?: string;
    role: 'user' | 'player' | 'admin';
    status: 'active' | 'suspended' | 'banned';
    last_login_at?: string;
    created_at: string;
    updated_at: string;
  }>;
  total: number;
  page: number;
  page_size: number;
}
```

### 2.2 获取用户详情

```http
GET /api/users/:id
```

**路径参数**: `id` - 用户ID

**响应**:

```typescript
{
  id: number;
  name: string;
  phone?: string;
  email?: string;
  avatar_url?: string;
  role: 'user' | 'player' | 'admin';
  status: 'active' | 'suspended' | 'banned';
  last_login_at?: string;
  created_at: string;
  updated_at: string;

  // 统计信息
  stats: {
    order_count: number;      // 订单数量
    total_spent_cents: number; // 总消费（分）
    review_count: number;      // 评价数量
  };

  // 陪玩师信息（仅 role=player 时存在）
  player?: {
    id: number;
    user_id: number;
    nickname?: string;
    bio?: string;
    rating_average: number;
    rating_count: number;
    hourly_rate_cents: number;
    main_game_id?: number;
    verification_status: 'pending' | 'verified' | 'rejected';
    created_at: string;
    updated_at: string;
  };
}
```

### 2.3 更新用户状态

```http
PUT /api/users/:id/status
```

**请求体**:

```typescript
{
  status: 'active' | 'suspended' | 'banned';
  reason?: string;  // 暂停/封禁原因
}
```

**响应**:

```typescript
{
  id: number;
  status: 'active' | 'suspended' | 'banned';
  updated_at: string;
}
```

### 2.4 更新用户角色

```http
PUT /api/users/:id/role
```

**请求体**:

```typescript
{
  role: 'user' | 'player' | 'admin';
}
```

**响应**:

```typescript
{
  id: number;
  role: 'user' | 'player' | 'admin';
  updated_at: string;
}
```

### 2.5 获取用户订单列表

```http
GET /api/users/:id/orders
```

**Query 参数**: 同 [4.1 订单列表](#41-获取订单列表)

---

## 3. 陪玩师管理接口

### 3.1 获取陪玩师列表

```http
GET /api/players
```

**Query 参数**:

```typescript
{
  page?: number;
  page_size?: number;
  keyword?: string;  // 昵称/姓名搜索
  verification_status?: 'pending' | 'verified' | 'rejected';
  game_id?: number;  // 按游戏筛选
  min_rating?: number;  // 最低评分
  sort_by?: 'rating_average' | 'rating_count' | 'hourly_rate_cents' | 'created_at';
  sort_order?: 'asc' | 'desc';
}
```

**响应**:

```typescript
{
  list: Array<{
    id: number;
    user_id: number;
    user: {
      id: number;
      name: string;
      avatar_url?: string;
      status: 'active' | 'suspended' | 'banned';
    };
    nickname?: string;
    bio?: string;
    rating_average: number;
    rating_count: number;
    hourly_rate_cents: number;
    main_game_id?: number;
    verification_status: 'pending' | 'verified' | 'rejected';
    created_at: string;
    updated_at: string;
  }>;
  total: number;
  page: number;
  page_size: number;
}
```

### 3.2 获取陪玩师详情

```http
GET /api/players/:id
```

**响应**:

```typescript
{
  id: number;
  user_id: number;
  user: {
    id: number;
    name: string;
    phone?: string;
    email?: string;
    avatar_url?: string;
    status: 'active' | 'suspended' | 'banned';
  };
  nickname?: string;
  bio?: string;
  rating_average: number;
  rating_count: number;
  hourly_rate_cents: number;
  main_game_id?: number;
  verification_status: 'pending' | 'verified' | 'rejected';
  created_at: string;
  updated_at: string;

  // 关联游戏
  games: Array<{
    id: number;
    name: string;
    icon?: string;
    is_main: boolean;
  }>;

  // 技能标签
  skill_tags: string[];

  // 统计信息
  stats: {
    order_count: number;
    completed_count: number;
    total_earned_cents: number;
  };
}
```

### 3.3 更新陪玩师认证状态

```http
PUT /api/players/:id/verification
```

**请求体**:

```typescript
{
  status: 'verified' | 'rejected';
  reason?: string;  // 拒绝原因
}
```

**响应**:

```typescript
{
  id: number;
  verification_status: 'verified' | 'rejected';
  updated_at: string;
}
```

### 3.4 更新陪玩师游戏

```http
PUT /api/players/:id/games
```

**请求体**:

```typescript
{
  game_ids: number[];
  main_game_id?: number;
}
```

**响应**:

```typescript
{
  id: number;
  games: Array<{
    id: number;
    name: string;
    is_main: boolean;
  }>;
}
```

### 3.5 更新陪玩师技能标签

```http
PUT /api/players/:id/skill-tags
```

**请求体**:

```typescript
{
  tags: string[];
}
```

**响应**:

```typescript
{
  id: number;
  skill_tags: string[];
}
```

---

## 4. 订单管理接口

### 4.1 获取订单列表

```http
GET /api/orders
```

**Query 参数**:

```typescript
{
  page?: number;
  page_size?: number;
  keyword?: string;  // 订单号/标题搜索
  status?: 'pending' | 'paid' | 'accepted' | 'in_progress' | 'pending_review' | 'completed' | 'cancelled';
  review_status?: 'pending' | 'approved' | 'rejected';
  game_id?: number;
  user_id?: number;
  player_id?: number;
  created_start?: string;
  created_end?: string;
  sort_by?: 'created_at' | 'price_cents' | 'scheduled_start';
  sort_order?: 'asc' | 'desc';
}
```

**响应**:

```typescript
{
  list: Array<{
    id: number;
    order_no: string;
    user_id: number;
    player_id?: number;
    game_id: number;
    title: string;
    status:
      | 'pending'
      | 'paid'
      | 'accepted'
      | 'in_progress'
      | 'pending_review'
      | 'completed'
      | 'cancelled';
    review_status?: 'pending' | 'approved' | 'rejected';
    price_cents: number;
    currency: string;
    scheduled_start?: string;
    scheduled_end?: string;
    created_at: string;
    updated_at: string;

    // 关联信息
    user: {
      id: number;
      name: string;
      avatar_url?: string;
    };
    player?: {
      id: number;
      nickname?: string;
      avatar_url?: string;
    };
    game: {
      id: number;
      name: string;
      icon?: string;
    };
  }>;
  total: number;
  page: number;
  page_size: number;
}
```

### 4.2 获取订单详情

```http
GET /api/orders/:id
```

**响应**:

```typescript
{
  id: number;
  order_no: string;
  user_id: number;
  player_id?: number;
  game_id: number;
  title: string;
  description?: string;
  status: 'pending' | 'paid' | 'accepted' | 'in_progress' | 'pending_review' | 'completed' | 'cancelled';
  review_status?: 'pending' | 'approved' | 'rejected';
  price_cents: number;
  currency: string;
  scheduled_start?: string;
  scheduled_end?: string;
  cancel_reason?: string;
  created_at: string;
  updated_at: string;

  // 时间节点
  paid_at?: string;
  accepted_at?: string;
  started_at?: string;
  completed_at?: string;
  cancelled_at?: string;

  // 关联信息
  user: {
    id: number;
    name: string;
    phone?: string;
    avatar_url?: string;
  };
  player?: {
    id: number;
    user_id: number;
    nickname?: string;
    rating_average: number;
    rating_count: number;
    avatar_url?: string;
  };
  game: {
    id: number;
    name: string;
    icon?: string;
    category: string;
  };

  // 操作日志
  logs: Array<{
    id: number;
    order_id: number;
    action: 'create' | 'pay' | 'accept' | 'start' | 'submit_review' | 'approve' | 'reject' | 'complete' | 'cancel' | 'request_refund' | 'refund';
    content: string;
    operator: string;
    operator_role: string;
    status_before?: string;
    status_after?: string;
    created_at: string;
  }>;

  // 审核记录
  reviews: Array<{
    id: number;
    order_id: number;
    status: 'pending' | 'approved' | 'rejected';
    reviewer: string;
    reason?: string;
    created_at: string;
  }>;
}
```

### 4.3 审核订单

```http
POST /api/orders/:id/review
```

**请求体**:

```typescript
{
  result: 'approved' | 'rejected';
  reason?: string;  // 拒绝原因
  note?: string;    // 备注
}
```

**响应**:

```typescript
{
  id: number;
  review_status: 'approved' | 'rejected';
  updated_at: string;
}
```

### 4.4 取消订单

```http
POST /api/orders/:id/cancel
```

**请求体**:

```typescript
{
  reason: string;
}
```

**响应**:

```typescript
{
  id: number;
  status: 'cancelled';
  cancel_reason: string;
  cancelled_at: string;
}
```

---

## 5. 游戏管理接口

### 5.1 获取游戏列表

```http
GET /api/games
```

**Query 参数**:

```typescript
{
  page?: number;
  page_size?: number;
  keyword?: string;
  category?: string;
  status?: 'active' | 'inactive';
  sort_by?: 'name' | 'player_count' | 'created_at';
  sort_order?: 'asc' | 'desc';
}
```

**响应**:

```typescript
{
  list: Array<{
    id: number;
    name: string;
    name_en?: string;
    icon?: string;
    banner?: string;
    category: string;
    tags: string[];
    status: 'active' | 'inactive';
    player_count: number; // 陪玩师数量
    order_count: number; // 订单数量
    description?: string;
    created_at: string;
    updated_at: string;
  }>;
  total: number;
  page: number;
  page_size: number;
}
```

### 5.2 获取游戏详情

```http
GET /api/games/:id
```

**响应**:

```typescript
{
  id: number;
  name: string;
  name_en?: string;
  icon?: string;
  banner?: string;
  category: string;
  tags: string[];
  status: 'active' | 'inactive';
  description?: string;
  created_at: string;
  updated_at: string;

  // 统计信息
  stats: {
    player_count: number;
    order_count: number;
    total_revenue_cents: number;
  };
}
```

### 5.3 创建游戏

```http
POST /api/games
```

**请求体**:

```typescript
{
  name: string;
  name_en?: string;
  icon?: string;      // 图片URL
  banner?: string;    // 横幅图URL
  category: string;
  tags?: string[];
  description?: string;
  status?: 'active' | 'inactive';
}
```

**响应**:

```typescript
{
  id: number;
  name: string;
  // ... 其他字段
  created_at: string;
}
```

### 5.4 更新游戏

```http
PUT /api/games/:id
```

**请求体**: 同 5.3

**响应**: 同 5.3

### 5.5 删除游戏

```http
DELETE /api/games/:id
```

**响应**: `204 No Content`

### 5.6 更新游戏状态

```http
PUT /api/games/:id/status
```

**请求体**:

```typescript
{
  status: 'active' | 'inactive';
}
```

**响应**:

```typescript
{
  id: number;
  status: 'active' | 'inactive';
  updated_at: string;
}
```

### 5.7 获取游戏分类列表

```http
GET /api/games/categories
```

**响应**:

```typescript
{
  categories: Array<{
    name: string;
    count: number;
  }>;
}
```

---

## 6. 支付管理接口

### 6.1 获取支付记录列表

```http
GET /api/payments
```

**Query 参数**:

```typescript
{
  page?: number;
  page_size?: number;
  order_id?: number;
  user_id?: number;
  method?: 'wechat' | 'alipay' | 'balance';
  status?: 'pending' | 'paid' | 'failed' | 'refunded';
  created_start?: string;
  created_end?: string;
}
```

**响应**:

```typescript
{
  list: Array<{
    id: number;
    order_id: number;
    user_id: number;
    method: 'wechat' | 'alipay' | 'balance';
    amount_cents: number;
    currency: string;
    status: 'pending' | 'paid' | 'failed' | 'refunded';
    provider_trade_no?: string;
    paid_at?: string;
    refunded_at?: string;
    created_at: string;
    updated_at: string;

    // 关联信息
    order: {
      id: number;
      order_no: string;
      title: string;
    };
    user: {
      id: number;
      name: string;
    };
  }>;
  total: number;
  page: number;
  page_size: number;
}
```

### 6.2 获取支付详情

```http
GET /api/payments/:id
```

**响应**:

```typescript
{
  id: number;
  order_id: number;
  user_id: number;
  method: 'wechat' | 'alipay' | 'balance';
  amount_cents: number;
  currency: string;
  status: 'pending' | 'paid' | 'failed' | 'refunded';
  provider_trade_no?: string;
  provider_raw?: object;
  paid_at?: string;
  refunded_at?: string;
  created_at: string;
  updated_at: string;

  order: {
    id: number;
    order_no: string;
    title: string;
    status: string;
  };
  user: {
    id: number;
    name: string;
    phone?: string;
  };
}
```

### 6.3 退款处理

```http
POST /api/payments/:id/refund
```

**请求体**:

```typescript
{
  reason: string;
  amount_cents?: number;  // 部分退款金额，不填则全额退款
}
```

**响应**:

```typescript
{
  id: number;
  status: 'refunded';
  refunded_at: string;
}
```

---

## 7. 评价管理接口

### 7.1 获取评价列表

```http
GET /api/reviews
```

**Query 参数**:

```typescript
{
  page?: number;
  page_size?: number;
  order_id?: number;
  user_id?: number;
  player_id?: number;
  score?: 1 | 2 | 3 | 4 | 5;
  created_start?: string;
  created_end?: string;
}
```

**响应**:

```typescript
{
  list: Array<{
    id: number;
    order_id: number;
    user_id: number;
    player_id: number;
    score: 1 | 2 | 3 | 4 | 5;
    content?: string;
    created_at: string;

    user: {
      id: number;
      name: string;
      avatar_url?: string;
    };
    player: {
      id: number;
      nickname?: string;
    };
    order: {
      id: number;
      order_no: string;
      title: string;
    };
  }>;
  total: number;
  page: number;
  page_size: number;
}
```

---

## 8. 数据统计接口

### 8.1 Dashboard 概览数据

```http
GET /api/stats/dashboard
```

**响应**:

```typescript
{
  // 关键指标
  metrics: {
    total_users: number;
    total_players: number;
    total_orders: number;
    total_revenue_cents: number;

    // 增长数据（与上期对比）
    user_growth_rate: number; // 用户增长率 %
    player_growth_rate: number; // 陪玩师增长率 %
    order_growth_rate: number; // 订单增长率 %
    revenue_growth_rate: number; // 收入增长率 %
  }

  // 今日数据
  today: {
    new_users: number;
    new_orders: number;
    revenue_cents: number;
    active_players: number;
  }

  // 订单状态分布
  order_status_distribution: {
    pending: number;
    in_progress: number;
    completed: number;
    cancelled: number;
  }
}
```

### 8.2 收入趋势

```http
GET /api/stats/revenue-trend
```

**Query 参数**:

```typescript
{
  start_date: string;  // ISO8601
  end_date: string;
  granularity?: 'day' | 'week' | 'month';  // 默认 day
}
```

**响应**:

```typescript
{
  data: Array<{
    date: string;
    revenue_cents: number;
    order_count: number;
  }>;
}
```

### 8.3 用户增长趋势

```http
GET /api/stats/user-growth
```

**Query 参数**: 同 8.2

**响应**:

```typescript
{
  data: Array<{
    date: string;
    new_users: number;
    new_players: number;
    total_users: number;
    total_players: number;
  }>;
}
```

### 8.4 订单统计

```http
GET /api/stats/orders
```

**Query 参数**:

```typescript
{
  start_date: string;
  end_date: string;
  group_by?: 'game' | 'player' | 'status';
}
```

**响应**:

```typescript
{
  data: Array<{
    key: string; // 游戏名/陪玩师ID/状态
    count: number;
    revenue_cents: number;
  }>;
  total_count: number;
  total_revenue_cents: number;
}
```

### 8.5 Top 陪玩师排行

```http
GET /api/stats/top-players
```

**Query 参数**:

```typescript
{
  limit?: number;  // 默认 10
  sort_by?: 'revenue' | 'orders' | 'rating';
  start_date?: string;
  end_date?: string;
}
```

**响应**:

```typescript
{
  list: Array<{
    player_id: number;
    nickname?: string;
    avatar_url?: string;
    order_count: number;
    revenue_cents: number;
    rating_average: number;
  }>;
}
```

---

## 9. 权限管理接口

### 9.1 获取角色列表

```http
GET /api/roles
```

**响应**:

```typescript
{
  list: Array<{
    id: number;
    name: string;
    code: string;
    description?: string;
    permissions: string[]; // 权限代码列表
    created_at: string;
  }>;
}
```

### 9.2 获取权限列表

```http
GET /api/permissions
```

**响应**:

```typescript
{
  list: Array<{
    id: number;
    name: string;
    code: string;
    category: string;
    description?: string;
  }>;
}
```

### 9.3 更新角色权限

```http
PUT /api/roles/:id/permissions
```

**请求体**:

```typescript
{
  permission_codes: string[];
}
```

**响应**:

```typescript
{
  id: number;
  permissions: string[];
  updated_at: string;
}
```

### 9.4 获取操作日志

```http
GET /api/audit-logs
```

**Query 参数**:

```typescript
{
  page?: number;
  page_size?: number;
  user_id?: number;
  action?: string;
  resource_type?: string;
  start_date?: string;
  end_date?: string;
}
```

**响应**:

```typescript
{
  list: Array<{
    id: number;
    user_id: number;
    user_name: string;
    action: string;
    resource_type: string;
    resource_id?: number;
    ip_address?: string;
    user_agent?: string;
    created_at: string;
  }>;
  total: number;
  page: number;
  page_size: number;
}
```

---

## 10. 系统设置接口

### 10.1 获取系统配置

```http
GET /api/settings
```

**响应**:

```typescript
{
  platform: {
    name: string;
    logo?: string;
    contact_email?: string;
    contact_phone?: string;
  };
  commission: {
    platform_rate: number;  // 平台抽成比例 0-1
    min_withdrawal_cents: number;  // 最低提现金额（分）
  };
  order: {
    auto_cancel_minutes: number;  // 自动取消未支付订单（分钟）
    max_duration_hours: number;   // 最大服务时长（小时）
  };
  maintenance: {
    enabled: boolean;
    message?: string;
    start_time?: string;
    end_time?: string;
  };
}
```

### 10.2 更新系统配置

```http
PUT /api/settings
```

**请求体**: 同 10.1 响应格式

**响应**: 同 10.1

---

## 📌 通用规范

### 请求头

所有需要认证的接口必须携带：

```
Authorization: Bearer {token}
Content-Type: application/json
```

### 响应格式

#### 成功响应

```typescript
{
  // 直接返回数据
}
```

#### 错误响应

```typescript
{
  error: {
    code: string;       // 错误代码
    message: string;    // 错误信息
    details?: any;      // 详细信息
  }
}
```

### HTTP 状态码

- `200` - 成功
- `201` - 创建成功
- `204` - 成功（无内容）
- `400` - 请求参数错误
- `401` - 未认证
- `403` - 无权限
- `404` - 资源不存在
- `409` - 资源冲突
- `422` - 验证失败
- `500` - 服务器错误

### 分页参数

所有列表接口统一使用：

```typescript
{
  page: number; // 页码，从 1 开始
  page_size: number; // 每页数量
}
```

### 时间格式

所有时间字段使用 ISO8601 格式：

```
2025-01-05T10:30:00Z
```

### 金额单位

所有金额统一使用 **分（cents）** 为单位，避免浮点数精度问题。

---

## 🔄 接口优先级

### 🔴 高优先级（立即需要）

1. **认证接口** (1.1-1.4)
2. **用户管理** (2.1-2.2)
3. **订单管理** (4.1-4.3)
4. **Dashboard统计** (8.1)

### 🟡 中优先级（2周内）

1. **陪玩师管理** (3.1-3.3)
2. **游戏管理** (5.1-5.6)
3. **支付管理** (6.1-6.3)
4. **数据趋势** (8.2-8.5)

### 🟢 低优先级（后续迭代）

1. **评价管理** (7.1)
2. **权限管理** (9.1-9.4)
3. **系统设置** (10.1-10.2)

---

## 📝 备注

1. 所有 ID 使用 `uint64`（Go）/ `number`（TypeScript）
2. 枚举值统一使用小写字符串
3. 可选字段使用 `?` 标记
4. 软删除的数据不应出现在列表中
5. 分页最大 `page_size` 建议限制为 100

---

**文档版本**: v1.0.0  
**维护人**: 前端开发团队  
**最后更新**: 2025-01-05
