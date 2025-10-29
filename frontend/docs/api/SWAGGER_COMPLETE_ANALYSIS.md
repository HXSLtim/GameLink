# GameLink 后端 Swagger 完整分析

> 数据来源: swagger.json  
> API 版本: 0.3.0  
> 生成时间: 2025-10-28  
> Base Path: `/api/v1`

---

## 📊 接口总览

| 模块 | 接口数量 | 状态 |
|------|---------|------|
| 🔐 Auth | 5个 | ✅ 已实现 |
| 👤 Admin/Users | 9个 | ✅ 已实现 |
| 🎮 Admin/Players | 10个 | ⚠️ 部分实现 |
| 📦 Admin/Orders | 9个 | ✅ 已实现 |
| 🎯 Admin/Games | 6个 | ✅ 已实现 |
| 💳 Admin/Payments | 8个 | ⚠️ 模型不完整 |
| ⭐ Admin/Reviews | 6个 | ✅ 已实现 |
| 📈 Admin/Stats | 6个 | ❌ 未实现 |
| **总计** | **49个** | **75%** |

---

## 🔐 Auth 模块 (5个接口)

### 1. POST /auth/login
**用途**: 用户登录  
**请求体**: `handler.loginRequest`
```typescript
{
  username: string;  // 必填：用户名（邮箱或手机号）
  password: string;  // 必填：密码
}
```
**响应**: `handler.loginResponse`
```typescript
{
  token: string;
  expires_at: string;
  user: User;
}
```

### 2. POST /auth/register
**用途**: 用户注册  
**请求体**: `handler.registerRequest`
```typescript
{
  name: string;      // 必填：姓名
  password: string;  // 必填：密码（最少6位）
  email?: string;    // 可选：邮箱
  phone?: string;    // 可选：手机号
}
```
**响应**: `handler.loginResponse`

### 3. POST /auth/logout
**用途**: 用户登出  
**认证**: 需要 Bearer Token  
**响应**: 成功状态

### 4. GET /auth/me
**用途**: 获取当前用户信息  
**认证**: 需要 Bearer Token  
**响应**: `handler.loginResponse`

### 5. POST /auth/refresh
**用途**: 刷新 Token  
**认证**: 需要 Bearer Token  
**响应**: `handler.tokenPayload`

---

## 👤 Admin/Users 模块 (9个接口)

### 1. GET /admin/users
**用途**: 用户列表（分页+筛选）  
**查询参数**:
- `page`: 页码
- `page_size`: 每页数量
- `role[]`: 角色过滤（多值）
- `status[]`: 状态过滤（多值）
- `keyword`: 关键字（匹配 name/email/phone）
- `date_from`, `date_to`: 时间范围

### 2. POST /admin/users
**用途**: 创建用户  
**请求体**: `admin.CreateUserPayload`
```typescript
{
  name: string;        // 必填
  password: string;    // 必填：最少6位
  role: string;        // 必填
  status: string;      // 必填
  email?: string;
  phone?: string;
  avatar_url?: string;
}
```

### 3. GET /admin/users/{id}
**用途**: 获取用户详情

### 4. PUT /admin/users/{id}
**用途**: 更新用户信息  
**请求体**: `admin.UpdateUserPayload`
```typescript
{
  name: string;        // 必填
  role: string;        // 必填
  status: string;      // 必填
  email?: string;
  phone?: string;
  avatar_url?: string;
  password?: string;   // 可选：更新密码
}
```

### 5. DELETE /admin/users/{id}
**用途**: 删除用户

### 6. PUT /admin/users/{id}/role
**用途**: 更新用户角色  
**请求体**: `{ role: string }`

### 7. PUT /admin/users/{id}/status
**用途**: 更新用户状态  
**请求体**: `{ status: string }`

### 8. GET /admin/users/{id}/orders
**用途**: 获取用户的订单列表  
**查询参数**:
- `page`, `page_size`: 分页
- `status[]`: 订单状态
- `date_from`, `date_to`: 时间范围

### 9. GET /admin/users/{id}/logs
**用途**: 获取用户操作日志  
**支持导出**: CSV  
**查询参数**:
- `page`, `page_size`: 分页
- `action`: 动作过滤（create/update/delete）
- `actor_user_id`: 操作者ID
- `date_from`, `date_to`: 时间范围
- `export`: 导出格式（csv）
- `fields`: 导出列
- `header_lang`: 列头语言（en/zh）

---

## 🎮 Admin/Players 模块 (10个接口)

### 1. GET /admin/players
**用途**: 陪玩师列表

### 2. POST /admin/players
**用途**: 创建陪玩师  
**请求体**: `admin.CreatePlayerPayload`
```typescript
{
  user_id: number;              // 必填
  verification_status: string;  // 必填
  nickname?: string;
  bio?: string;
  main_game_id?: number;
  hourly_rate_cents?: number;
}
```

### 3. GET /admin/players/{id}
**用途**: 获取陪玩师详情

### 4. PUT /admin/players/{id}
**用途**: 更新陪玩师信息  
**请求体**: `admin.UpdatePlayerPayload`
```typescript
{
  verification_status: string;  // 必填
  nickname?: string;
  bio?: string;
  main_game_id?: number;
  hourly_rate_cents?: number;
}
```

### 5. DELETE /admin/players/{id}
**用途**: 删除陪玩师

### 6. PUT /admin/players/{id}/games
**用途**: 更新陪玩师主游戏  
**请求体**: `{ main_game_id: number }`

### 7. PUT /admin/players/{id}/skill-tags
**用途**: 更新陪玩师技能标签  
**请求体**: `admin.SkillTagsBody`
```typescript
{
  tags: string[];  // 必填：标签数组
}
```

### 8. PUT /admin/players/{id}/verification
**用途**: 更新陪玩师认证状态  
**请求体**: `{ verification_status: string }`

### 9. GET /admin/players/{id}/reviews
**用途**: 获取陪玩师的评价列表  
**查询参数**: `page`, `page_size`

### 10. GET /admin/players/{id}/logs
**用途**: 获取陪玩师操作日志  
**支持导出**: CSV

---

## 📦 Admin/Orders 模块 (9个接口)

### 1. GET /admin/orders
**用途**: 订单列表（分页+筛选）  
**查询参数**:
- `page`, `page_size`: 分页
- `status[]`: 订单状态（多值）
- `user_id`: 用户ID
- `player_id`: 陪玩师ID
- `game_id`: 游戏ID
- `date_from`, `date_to`: 时间范围

### 2. POST /admin/orders
**用途**: 创建订单  
**请求体**: `admin.CreateOrderPayload`
```typescript
{
  user_id: number;         // 必填
  game_id: number;         // 必填
  price_cents: number;     // 必填
  currency: string;        // 必填
  player_id?: number;      // 可选：指定陪玩师
  title?: string;
  description?: string;
  scheduled_start?: string;
  scheduled_end?: string;
}
```

### 3. GET /admin/orders/{id}
**用途**: 获取订单详情

### 4. PUT /admin/orders/{id}
**用途**: 更新订单信息  
**请求体**: `admin.UpdateOrderPayload`
```typescript
{
  price_cents: number;     // 必填
  currency: string;        // 必填
  status: string;          // 必填
  scheduled_start?: string;
  scheduled_end?: string;
  cancel_reason?: string;
}
```

### 5. DELETE /admin/orders/{id}
**用途**: 删除订单

### 6. POST /admin/orders/{id}/assign
**用途**: 指派陪玩师  
**请求体**: `admin.AssignOrderPayload`
```typescript
{
  player_id: number;  // 必填
}
```

### 7. POST /admin/orders/{id}/review
**用途**: 审核订单（通过/拒绝）  
**请求体**: `admin.ReviewOrderPayload`
```typescript
{
  approved?: boolean;  // true=通过，false=拒绝
  reason?: string;     // 审核理由
}
```

### 8. POST /admin/orders/{id}/cancel
**用途**: 取消订单  
**请求体**: `admin.CancelOrderPayload`
```typescript
{
  reason?: string;  // 取消原因
}
```

### 9. GET /admin/orders/{id}/logs
**用途**: 获取订单操作日志  
**支持导出**: CSV  
**动作类型**: create, assign_player, update_status, cancel, delete

---

## 🎯 Admin/Games 模块 (6个接口)

### 1. GET /admin/games
**用途**: 游戏列表  
**查询参数**: `page`, `page_size`

### 2. POST /admin/games
**用途**: 创建游戏  
**请求体**: `admin.GamePayload`
```typescript
{
  key: string;         // 必填：游戏唯一标识
  name: string;        // 必填：游戏名称
  category?: string;   // 游戏分类
  description?: string;
  icon_url?: string;
}
```

### 3. GET /admin/games/{id}
**用途**: 获取游戏详情

### 4. PUT /admin/games/{id}
**用途**: 更新游戏信息  
**请求体**: `admin.GamePayload`

### 5. DELETE /admin/games/{id}
**用途**: 删除游戏

### 6. GET /admin/games/{id}/logs
**用途**: 获取游戏操作日志  
**支持导出**: CSV

---

## 💳 Admin/Payments 模块 (8个接口)

### 1. GET /admin/payments
**用途**: 支付列表（分页+筛选）  
**查询参数**:
- `page`, `page_size`: 分页
- `status[]`: 支付状态
- `method[]`: 支付方式
- `user_id`: 用户ID
- `order_id`: 订单ID
- `date_from`, `date_to`: 时间范围

### 2. POST /admin/payments
**用途**: 创建支付记录  
**请求体**: `admin.CreatePaymentPayload` ⚠️ **空对象**

### 3. GET /admin/payments/{id}
**用途**: 获取支付详情

### 4. PUT /admin/payments/{id}
**用途**: 更新支付信息  
**请求体**: `admin.UpdatePaymentPayload` ⚠️ **空对象**

### 5. DELETE /admin/payments/{id}
**用途**: 删除支付记录

### 6. POST /admin/payments/{id}/capture
**用途**: 确认支付入账  
**请求体**: `admin.CapturePaymentPayload` ⚠️ **空对象**

### 7. POST /admin/payments/{id}/refund
**用途**: 退款处理  
**请求体**: `admin.RefundPaymentPayload` ⚠️ **空对象**

### 8. GET /admin/payments/{id}/logs
**用途**: 获取支付操作日志  
**支持导出**: CSV  
**动作类型**: create, capture, update_status, refund, delete

---

## ⭐ Admin/Reviews 模块 (6个接口)

### 1. GET /admin/reviews
**用途**: 评价列表  
**查询参数**:
- `page`, `page_size`: 分页
- `order_id`: 订单ID
- `user_id`: 用户ID
- `player_id`: 陪玩师ID
- `date_from`, `date_to`: 时间范围

### 2. POST /admin/reviews
**用途**: 创建评价  
**请求体**: `admin.CreateReviewPayload`
```typescript
{
  user_id: number;    // 必填
  player_id: number;  // 必填
  order_id: number;   // 必填
  score: number;      // 必填：评分
  content?: string;   // 评价内容
}
```

### 3. GET /admin/reviews/{id}
**用途**: 获取评价详情

### 4. PUT /admin/reviews/{id}
**用途**: 更新评价  
**请求体**: `admin.UpdateReviewPayload`
```typescript
{
  score: number;      // 必填
  content?: string;
}
```

### 5. DELETE /admin/reviews/{id}
**用途**: 删除评价

### 6. GET /admin/reviews/{id}/logs
**用途**: 获取评价操作日志  
**支持导出**: CSV

---

## 📈 Admin/Stats 模块 (6个接口)

### 1. GET /admin/stats/dashboard
**用途**: Dashboard 总览  
**前端状态**: ❌ 未实现

### 2. GET /admin/stats/orders
**用途**: 订单状态统计  
**前端状态**: ✅ 已实现

### 3. GET /admin/stats/revenue-trend
**用途**: 收入趋势（日）  
**查询参数**: `days` - 天数（默认7天）  
**前端状态**: ❌ 未实现

### 4. GET /admin/stats/user-growth
**用途**: 用户增长（日）  
**查询参数**: `days` - 天数（默认7天）  
**前端状态**: ❌ 未实现

### 5. GET /admin/stats/top-players
**用途**: TOP 陪玩师排行  
**查询参数**: `limit` - 数量（默认10）  
**前端状态**: ❌ 未实现

### 6. GET /admin/stats/audit/overview
**用途**: 审计总览（按实体/动作汇总）  
**查询参数**: `date_from`, `date_to`  
**前端状态**: ❌ 未实现

### 7. GET /admin/stats/audit/trend
**用途**: 审计趋势（日）  
**查询参数**:
- `date_from`, `date_to`: 时间范围
- `entity`: 实体类型（order/payment/player/game/review/user）
- `action`: 动作  
**前端状态**: ❌ 未实现

---

## 🚨 关键发现

### 1. Payment 模块数据模型不完整
⚠️ **问题**: 所有 Payment Payload 都是空对象
```json
"admin.CreatePaymentPayload": {
    "type": "object"
}
```

**影响**:
- 无法知道创建/更新支付需要哪些字段
- 前端无法实现完整的支付功能

**建议**: 后端需要补充以下 Payload 定义：
- `CreatePaymentPayload`: 至少需要 `order_id`, `amount_cents`, `currency`, `method`
- `UpdatePaymentPayload`: 至少需要 `status`
- `CapturePaymentPayload`: 可能需要 `captured_amount_cents`
- `RefundPaymentPayload`: 至少需要 `refund_amount_cents`, `reason`

### 2. 用户端接口完全缺失
⚠️ **严重问题**: 没有任何用户端接口

**缺失的核心功能**:
- 用户浏览游戏/陪玩师
- 用户创建订单
- 用户查看自己的订单
- 用户提交评价
- 用户个人中心

**建议**: 需要开发完整的用户端 API（约30-40个接口）

### 3. Stats 模块前端未实现
❌ **问题**: 6个统计接口中，只实现了 `orders` 统计

**缺失功能**:
- Dashboard 总览（最重要）
- 收入趋势图表
- 用户增长图表
- TOP 陪玩师榜单
- 审计统计和趋势

### 4. 操作日志导出功能未使用
✅ **发现**: 各模块都提供了 CSV 导出功能，但前端未实现

**支持导出的接口**:
- 用户日志
- 陪玩师日志
- 订单日志
- 游戏日志
- 支付日志
- 评价日志

---

## 📋 前端待实现接口清单

### 高优先级 (9个)

#### Stats 模块 (5个)
```typescript
// 创建 src/services/api/stats.ts

export const statsApi = {
  // 1. Dashboard 总览 ⭐ 最高优先级
  getDashboard: (): Promise<DashboardStats> => {
    return apiClient.get('/api/v1/admin/stats/dashboard');
  },

  // 2. 收入趋势
  getRevenueTrend: (days?: number): Promise<RevenueTrendData> => {
    return apiClient.get('/api/v1/admin/stats/revenue-trend', { 
      params: { days } 
    });
  },

  // 3. 用户增长
  getUserGrowth: (days?: number): Promise<UserGrowthData> => {
    return apiClient.get('/api/v1/admin/stats/user-growth', { 
      params: { days } 
    });
  },

  // 4. TOP 陪玩师
  getTopPlayers: (limit?: number): Promise<TopPlayer[]> => {
    return apiClient.get('/api/v1/admin/stats/top-players', { 
      params: { limit } 
    });
  },

  // 5. 审计总览
  getAuditOverview: (params?: { 
    date_from?: string; 
    date_to?: string; 
  }): Promise<AuditOverview> => {
    return apiClient.get('/api/v1/admin/stats/audit/overview', { params });
  },

  // 6. 审计趋势
  getAuditTrend: (params?: {
    date_from?: string;
    date_to?: string;
    entity?: 'order' | 'payment' | 'player' | 'game' | 'review' | 'user';
    action?: string;
  }): Promise<AuditTrendData> => {
    return apiClient.get('/api/v1/admin/stats/audit/trend', { params });
  },
};
```

#### Player 模块 (4个)
```typescript
// 补充到 src/services/api/user.ts (playerApi)

// 1. 更新主游戏
updateMainGame: (id: number, main_game_id: number): Promise<Player> => {
  return apiClient.put(`/api/v1/admin/players/${id}/games`, { main_game_id });
},

// 2. 更新技能标签
updateSkillTags: (id: number, tags: string[]): Promise<Player> => {
  return apiClient.put(`/api/v1/admin/players/${id}/skill-tags`, { tags });
},

// 3. 获取陪玩师游戏列表（暂无，需要后端补充）
// getPlayerGames: (id: number): Promise<Game[]> => {
//   return apiClient.get(`/api/v1/admin/players/${id}/games`);
// },
```

### 中优先级 (导出功能)

#### 操作日志导出 (6个模块)
```typescript
// 为各个 API 模块添加日志导出功能

// 示例：订单日志导出
exportOrderLogs: (
  id: number,
  params: {
    export: 'csv';
    fields?: string;
    header_lang?: 'en' | 'zh';
    // ... 其他筛选参数
  }
): Promise<Blob> => {
  return apiClient.get(`/api/v1/admin/orders/${id}/logs`, {
    params,
    responseType: 'blob',
  });
},
```

---

## 📊 前后端同步度评估

| 模块 | 后端接口 | 前端实现 | 同步度 | 评级 |
|------|---------|---------|--------|------|
| Auth | 5 | 4 | 80% | ⭐⭐⭐⭐ |
| Users | 9 | 7 | 78% | ⭐⭐⭐⭐ |
| Players | 10 | 6 | 60% | ⭐⭐⭐ |
| Orders | 9 | 9 | 100% | ⭐⭐⭐⭐⭐ |
| Games | 6 | 6 | 100% | ⭐⭐⭐⭐⭐ |
| Payments | 8 | 6 | 75% | ⭐⭐⭐ |
| Reviews | 6 | 6 | 100% | ⭐⭐⭐⭐⭐ |
| Stats | 6 | 1 | 17% | ⭐ |
| **总计** | **49** | **39** | **80%** | ⭐⭐⭐⭐ |

---

## 🎯 下一步行动建议

### 阶段 1: 补全管理端 (2-3天)
1. ✅ 实现 Stats 模块的 5 个接口
2. ✅ 补充 Player 模块的技能标签和主游戏管理
3. ✅ 实现操作日志导出功能

### 阶段 2: 完善 Payment 模块 (1天)
1. 🔴 **后端**: 补充 Payment Payload 定义
2. ✅ **前端**: 实现完整的支付流程

### 阶段 3: 开发用户端 (7-10天)
1. 🔴 **后端**: 开发用户端完整接口（30-40个）
2. ✅ **前端**: 开发用户端页面和 API 调用

---

**文档版本**: v2.0  
**数据来源**: swagger.json  
**最后更新**: 2025-10-28



> 数据来源: swagger.json  
> API 版本: 0.3.0  
> 生成时间: 2025-10-28  
> Base Path: `/api/v1`

---

## 📊 接口总览

| 模块 | 接口数量 | 状态 |
|------|---------|------|
| 🔐 Auth | 5个 | ✅ 已实现 |
| 👤 Admin/Users | 9个 | ✅ 已实现 |
| 🎮 Admin/Players | 10个 | ⚠️ 部分实现 |
| 📦 Admin/Orders | 9个 | ✅ 已实现 |
| 🎯 Admin/Games | 6个 | ✅ 已实现 |
| 💳 Admin/Payments | 8个 | ⚠️ 模型不完整 |
| ⭐ Admin/Reviews | 6个 | ✅ 已实现 |
| 📈 Admin/Stats | 6个 | ❌ 未实现 |
| **总计** | **49个** | **75%** |

---

## 🔐 Auth 模块 (5个接口)

### 1. POST /auth/login
**用途**: 用户登录  
**请求体**: `handler.loginRequest`
```typescript
{
  username: string;  // 必填：用户名（邮箱或手机号）
  password: string;  // 必填：密码
}
```
**响应**: `handler.loginResponse`
```typescript
{
  token: string;
  expires_at: string;
  user: User;
}
```

### 2. POST /auth/register
**用途**: 用户注册  
**请求体**: `handler.registerRequest`
```typescript
{
  name: string;      // 必填：姓名
  password: string;  // 必填：密码（最少6位）
  email?: string;    // 可选：邮箱
  phone?: string;    // 可选：手机号
}
```
**响应**: `handler.loginResponse`

### 3. POST /auth/logout
**用途**: 用户登出  
**认证**: 需要 Bearer Token  
**响应**: 成功状态

### 4. GET /auth/me
**用途**: 获取当前用户信息  
**认证**: 需要 Bearer Token  
**响应**: `handler.loginResponse`

### 5. POST /auth/refresh
**用途**: 刷新 Token  
**认证**: 需要 Bearer Token  
**响应**: `handler.tokenPayload`

---

## 👤 Admin/Users 模块 (9个接口)

### 1. GET /admin/users
**用途**: 用户列表（分页+筛选）  
**查询参数**:
- `page`: 页码
- `page_size`: 每页数量
- `role[]`: 角色过滤（多值）
- `status[]`: 状态过滤（多值）
- `keyword`: 关键字（匹配 name/email/phone）
- `date_from`, `date_to`: 时间范围

### 2. POST /admin/users
**用途**: 创建用户  
**请求体**: `admin.CreateUserPayload`
```typescript
{
  name: string;        // 必填
  password: string;    // 必填：最少6位
  role: string;        // 必填
  status: string;      // 必填
  email?: string;
  phone?: string;
  avatar_url?: string;
}
```

### 3. GET /admin/users/{id}
**用途**: 获取用户详情

### 4. PUT /admin/users/{id}
**用途**: 更新用户信息  
**请求体**: `admin.UpdateUserPayload`
```typescript
{
  name: string;        // 必填
  role: string;        // 必填
  status: string;      // 必填
  email?: string;
  phone?: string;
  avatar_url?: string;
  password?: string;   // 可选：更新密码
}
```

### 5. DELETE /admin/users/{id}
**用途**: 删除用户

### 6. PUT /admin/users/{id}/role
**用途**: 更新用户角色  
**请求体**: `{ role: string }`

### 7. PUT /admin/users/{id}/status
**用途**: 更新用户状态  
**请求体**: `{ status: string }`

### 8. GET /admin/users/{id}/orders
**用途**: 获取用户的订单列表  
**查询参数**:
- `page`, `page_size`: 分页
- `status[]`: 订单状态
- `date_from`, `date_to`: 时间范围

### 9. GET /admin/users/{id}/logs
**用途**: 获取用户操作日志  
**支持导出**: CSV  
**查询参数**:
- `page`, `page_size`: 分页
- `action`: 动作过滤（create/update/delete）
- `actor_user_id`: 操作者ID
- `date_from`, `date_to`: 时间范围
- `export`: 导出格式（csv）
- `fields`: 导出列
- `header_lang`: 列头语言（en/zh）

---

## 🎮 Admin/Players 模块 (10个接口)

### 1. GET /admin/players
**用途**: 陪玩师列表

### 2. POST /admin/players
**用途**: 创建陪玩师  
**请求体**: `admin.CreatePlayerPayload`
```typescript
{
  user_id: number;              // 必填
  verification_status: string;  // 必填
  nickname?: string;
  bio?: string;
  main_game_id?: number;
  hourly_rate_cents?: number;
}
```

### 3. GET /admin/players/{id}
**用途**: 获取陪玩师详情

### 4. PUT /admin/players/{id}
**用途**: 更新陪玩师信息  
**请求体**: `admin.UpdatePlayerPayload`
```typescript
{
  verification_status: string;  // 必填
  nickname?: string;
  bio?: string;
  main_game_id?: number;
  hourly_rate_cents?: number;
}
```

### 5. DELETE /admin/players/{id}
**用途**: 删除陪玩师

### 6. PUT /admin/players/{id}/games
**用途**: 更新陪玩师主游戏  
**请求体**: `{ main_game_id: number }`

### 7. PUT /admin/players/{id}/skill-tags
**用途**: 更新陪玩师技能标签  
**请求体**: `admin.SkillTagsBody`
```typescript
{
  tags: string[];  // 必填：标签数组
}
```

### 8. PUT /admin/players/{id}/verification
**用途**: 更新陪玩师认证状态  
**请求体**: `{ verification_status: string }`

### 9. GET /admin/players/{id}/reviews
**用途**: 获取陪玩师的评价列表  
**查询参数**: `page`, `page_size`

### 10. GET /admin/players/{id}/logs
**用途**: 获取陪玩师操作日志  
**支持导出**: CSV

---

## 📦 Admin/Orders 模块 (9个接口)

### 1. GET /admin/orders
**用途**: 订单列表（分页+筛选）  
**查询参数**:
- `page`, `page_size`: 分页
- `status[]`: 订单状态（多值）
- `user_id`: 用户ID
- `player_id`: 陪玩师ID
- `game_id`: 游戏ID
- `date_from`, `date_to`: 时间范围

### 2. POST /admin/orders
**用途**: 创建订单  
**请求体**: `admin.CreateOrderPayload`
```typescript
{
  user_id: number;         // 必填
  game_id: number;         // 必填
  price_cents: number;     // 必填
  currency: string;        // 必填
  player_id?: number;      // 可选：指定陪玩师
  title?: string;
  description?: string;
  scheduled_start?: string;
  scheduled_end?: string;
}
```

### 3. GET /admin/orders/{id}
**用途**: 获取订单详情

### 4. PUT /admin/orders/{id}
**用途**: 更新订单信息  
**请求体**: `admin.UpdateOrderPayload`
```typescript
{
  price_cents: number;     // 必填
  currency: string;        // 必填
  status: string;          // 必填
  scheduled_start?: string;
  scheduled_end?: string;
  cancel_reason?: string;
}
```

### 5. DELETE /admin/orders/{id}
**用途**: 删除订单

### 6. POST /admin/orders/{id}/assign
**用途**: 指派陪玩师  
**请求体**: `admin.AssignOrderPayload`
```typescript
{
  player_id: number;  // 必填
}
```

### 7. POST /admin/orders/{id}/review
**用途**: 审核订单（通过/拒绝）  
**请求体**: `admin.ReviewOrderPayload`
```typescript
{
  approved?: boolean;  // true=通过，false=拒绝
  reason?: string;     // 审核理由
}
```

### 8. POST /admin/orders/{id}/cancel
**用途**: 取消订单  
**请求体**: `admin.CancelOrderPayload`
```typescript
{
  reason?: string;  // 取消原因
}
```

### 9. GET /admin/orders/{id}/logs
**用途**: 获取订单操作日志  
**支持导出**: CSV  
**动作类型**: create, assign_player, update_status, cancel, delete

---

## 🎯 Admin/Games 模块 (6个接口)

### 1. GET /admin/games
**用途**: 游戏列表  
**查询参数**: `page`, `page_size`

### 2. POST /admin/games
**用途**: 创建游戏  
**请求体**: `admin.GamePayload`
```typescript
{
  key: string;         // 必填：游戏唯一标识
  name: string;        // 必填：游戏名称
  category?: string;   // 游戏分类
  description?: string;
  icon_url?: string;
}
```

### 3. GET /admin/games/{id}
**用途**: 获取游戏详情

### 4. PUT /admin/games/{id}
**用途**: 更新游戏信息  
**请求体**: `admin.GamePayload`

### 5. DELETE /admin/games/{id}
**用途**: 删除游戏

### 6. GET /admin/games/{id}/logs
**用途**: 获取游戏操作日志  
**支持导出**: CSV

---

## 💳 Admin/Payments 模块 (8个接口)

### 1. GET /admin/payments
**用途**: 支付列表（分页+筛选）  
**查询参数**:
- `page`, `page_size`: 分页
- `status[]`: 支付状态
- `method[]`: 支付方式
- `user_id`: 用户ID
- `order_id`: 订单ID
- `date_from`, `date_to`: 时间范围

### 2. POST /admin/payments
**用途**: 创建支付记录  
**请求体**: `admin.CreatePaymentPayload` ⚠️ **空对象**

### 3. GET /admin/payments/{id}
**用途**: 获取支付详情

### 4. PUT /admin/payments/{id}
**用途**: 更新支付信息  
**请求体**: `admin.UpdatePaymentPayload` ⚠️ **空对象**

### 5. DELETE /admin/payments/{id}
**用途**: 删除支付记录

### 6. POST /admin/payments/{id}/capture
**用途**: 确认支付入账  
**请求体**: `admin.CapturePaymentPayload` ⚠️ **空对象**

### 7. POST /admin/payments/{id}/refund
**用途**: 退款处理  
**请求体**: `admin.RefundPaymentPayload` ⚠️ **空对象**

### 8. GET /admin/payments/{id}/logs
**用途**: 获取支付操作日志  
**支持导出**: CSV  
**动作类型**: create, capture, update_status, refund, delete

---

## ⭐ Admin/Reviews 模块 (6个接口)

### 1. GET /admin/reviews
**用途**: 评价列表  
**查询参数**:
- `page`, `page_size`: 分页
- `order_id`: 订单ID
- `user_id`: 用户ID
- `player_id`: 陪玩师ID
- `date_from`, `date_to`: 时间范围

### 2. POST /admin/reviews
**用途**: 创建评价  
**请求体**: `admin.CreateReviewPayload`
```typescript
{
  user_id: number;    // 必填
  player_id: number;  // 必填
  order_id: number;   // 必填
  score: number;      // 必填：评分
  content?: string;   // 评价内容
}
```

### 3. GET /admin/reviews/{id}
**用途**: 获取评价详情

### 4. PUT /admin/reviews/{id}
**用途**: 更新评价  
**请求体**: `admin.UpdateReviewPayload`
```typescript
{
  score: number;      // 必填
  content?: string;
}
```

### 5. DELETE /admin/reviews/{id}
**用途**: 删除评价

### 6. GET /admin/reviews/{id}/logs
**用途**: 获取评价操作日志  
**支持导出**: CSV

---

## 📈 Admin/Stats 模块 (6个接口)

### 1. GET /admin/stats/dashboard
**用途**: Dashboard 总览  
**前端状态**: ❌ 未实现

### 2. GET /admin/stats/orders
**用途**: 订单状态统计  
**前端状态**: ✅ 已实现

### 3. GET /admin/stats/revenue-trend
**用途**: 收入趋势（日）  
**查询参数**: `days` - 天数（默认7天）  
**前端状态**: ❌ 未实现

### 4. GET /admin/stats/user-growth
**用途**: 用户增长（日）  
**查询参数**: `days` - 天数（默认7天）  
**前端状态**: ❌ 未实现

### 5. GET /admin/stats/top-players
**用途**: TOP 陪玩师排行  
**查询参数**: `limit` - 数量（默认10）  
**前端状态**: ❌ 未实现

### 6. GET /admin/stats/audit/overview
**用途**: 审计总览（按实体/动作汇总）  
**查询参数**: `date_from`, `date_to`  
**前端状态**: ❌ 未实现

### 7. GET /admin/stats/audit/trend
**用途**: 审计趋势（日）  
**查询参数**:
- `date_from`, `date_to`: 时间范围
- `entity`: 实体类型（order/payment/player/game/review/user）
- `action`: 动作  
**前端状态**: ❌ 未实现

---

## 🚨 关键发现

### 1. Payment 模块数据模型不完整
⚠️ **问题**: 所有 Payment Payload 都是空对象
```json
"admin.CreatePaymentPayload": {
    "type": "object"
}
```

**影响**:
- 无法知道创建/更新支付需要哪些字段
- 前端无法实现完整的支付功能

**建议**: 后端需要补充以下 Payload 定义：
- `CreatePaymentPayload`: 至少需要 `order_id`, `amount_cents`, `currency`, `method`
- `UpdatePaymentPayload`: 至少需要 `status`
- `CapturePaymentPayload`: 可能需要 `captured_amount_cents`
- `RefundPaymentPayload`: 至少需要 `refund_amount_cents`, `reason`

### 2. 用户端接口完全缺失
⚠️ **严重问题**: 没有任何用户端接口

**缺失的核心功能**:
- 用户浏览游戏/陪玩师
- 用户创建订单
- 用户查看自己的订单
- 用户提交评价
- 用户个人中心

**建议**: 需要开发完整的用户端 API（约30-40个接口）

### 3. Stats 模块前端未实现
❌ **问题**: 6个统计接口中，只实现了 `orders` 统计

**缺失功能**:
- Dashboard 总览（最重要）
- 收入趋势图表
- 用户增长图表
- TOP 陪玩师榜单
- 审计统计和趋势

### 4. 操作日志导出功能未使用
✅ **发现**: 各模块都提供了 CSV 导出功能，但前端未实现

**支持导出的接口**:
- 用户日志
- 陪玩师日志
- 订单日志
- 游戏日志
- 支付日志
- 评价日志

---

## 📋 前端待实现接口清单

### 高优先级 (9个)

#### Stats 模块 (5个)
```typescript
// 创建 src/services/api/stats.ts

export const statsApi = {
  // 1. Dashboard 总览 ⭐ 最高优先级
  getDashboard: (): Promise<DashboardStats> => {
    return apiClient.get('/api/v1/admin/stats/dashboard');
  },

  // 2. 收入趋势
  getRevenueTrend: (days?: number): Promise<RevenueTrendData> => {
    return apiClient.get('/api/v1/admin/stats/revenue-trend', { 
      params: { days } 
    });
  },

  // 3. 用户增长
  getUserGrowth: (days?: number): Promise<UserGrowthData> => {
    return apiClient.get('/api/v1/admin/stats/user-growth', { 
      params: { days } 
    });
  },

  // 4. TOP 陪玩师
  getTopPlayers: (limit?: number): Promise<TopPlayer[]> => {
    return apiClient.get('/api/v1/admin/stats/top-players', { 
      params: { limit } 
    });
  },

  // 5. 审计总览
  getAuditOverview: (params?: { 
    date_from?: string; 
    date_to?: string; 
  }): Promise<AuditOverview> => {
    return apiClient.get('/api/v1/admin/stats/audit/overview', { params });
  },

  // 6. 审计趋势
  getAuditTrend: (params?: {
    date_from?: string;
    date_to?: string;
    entity?: 'order' | 'payment' | 'player' | 'game' | 'review' | 'user';
    action?: string;
  }): Promise<AuditTrendData> => {
    return apiClient.get('/api/v1/admin/stats/audit/trend', { params });
  },
};
```

#### Player 模块 (4个)
```typescript
// 补充到 src/services/api/user.ts (playerApi)

// 1. 更新主游戏
updateMainGame: (id: number, main_game_id: number): Promise<Player> => {
  return apiClient.put(`/api/v1/admin/players/${id}/games`, { main_game_id });
},

// 2. 更新技能标签
updateSkillTags: (id: number, tags: string[]): Promise<Player> => {
  return apiClient.put(`/api/v1/admin/players/${id}/skill-tags`, { tags });
},

// 3. 获取陪玩师游戏列表（暂无，需要后端补充）
// getPlayerGames: (id: number): Promise<Game[]> => {
//   return apiClient.get(`/api/v1/admin/players/${id}/games`);
// },
```

### 中优先级 (导出功能)

#### 操作日志导出 (6个模块)
```typescript
// 为各个 API 模块添加日志导出功能

// 示例：订单日志导出
exportOrderLogs: (
  id: number,
  params: {
    export: 'csv';
    fields?: string;
    header_lang?: 'en' | 'zh';
    // ... 其他筛选参数
  }
): Promise<Blob> => {
  return apiClient.get(`/api/v1/admin/orders/${id}/logs`, {
    params,
    responseType: 'blob',
  });
},
```

---

## 📊 前后端同步度评估

| 模块 | 后端接口 | 前端实现 | 同步度 | 评级 |
|------|---------|---------|--------|------|
| Auth | 5 | 4 | 80% | ⭐⭐⭐⭐ |
| Users | 9 | 7 | 78% | ⭐⭐⭐⭐ |
| Players | 10 | 6 | 60% | ⭐⭐⭐ |
| Orders | 9 | 9 | 100% | ⭐⭐⭐⭐⭐ |
| Games | 6 | 6 | 100% | ⭐⭐⭐⭐⭐ |
| Payments | 8 | 6 | 75% | ⭐⭐⭐ |
| Reviews | 6 | 6 | 100% | ⭐⭐⭐⭐⭐ |
| Stats | 6 | 1 | 17% | ⭐ |
| **总计** | **49** | **39** | **80%** | ⭐⭐⭐⭐ |

---

## 🎯 下一步行动建议

### 阶段 1: 补全管理端 (2-3天)
1. ✅ 实现 Stats 模块的 5 个接口
2. ✅ 补充 Player 模块的技能标签和主游戏管理
3. ✅ 实现操作日志导出功能

### 阶段 2: 完善 Payment 模块 (1天)
1. 🔴 **后端**: 补充 Payment Payload 定义
2. ✅ **前端**: 实现完整的支付流程

### 阶段 3: 开发用户端 (7-10天)
1. 🔴 **后端**: 开发用户端完整接口（30-40个）
2. ✅ **前端**: 开发用户端页面和 API 调用

---

**文档版本**: v2.0  
**数据来源**: swagger.json  
**最后更新**: 2025-10-28



