# GameLink 前端接口开发需求

## 📊 接口开发状态概览

### 后端接口统计

- **认证接口**: 5个
- **管理端接口**: 37个
- **用户端接口**: 0个（❌ 尚未开发）
- **总计**: 42个接口

---

## 🔴 管理端缺失接口（高优先级）

### 1. 认证模块

#### ❌ 注册接口

```typescript
POST / auth / register;
```

**状态**: 未实现  
**优先级**: 🔴 高  
**说明**: 前端已有 Register 页面，但缺少 API 调用

**需要添加到**: `src/services/api/auth.ts`

```typescript
register: (data: RegisterRequest): Promise<LoginResult> => {
  return apiClient.post('/api/v1/auth/register', data);
};
```

**类型定义**: `src/types/auth.ts`

```typescript
export interface RegisterRequest {
  email: string;
  name: string;
  password: string;
  phone?: string;
}
```

---

### 2. 陪玩师模块

#### ❌ 获取陪玩师的游戏列表

```typescript
GET / admin / players / { id } / games;
```

**状态**: 未实现  
**优先级**: 🔴 高  
**说明**: 陪玩师详情页需要展示其擅长的游戏列表

**需要添加到**: `src/services/api/user.ts` (playerApi)

```typescript
getPlayerGames: (id: number): Promise<Game[]> => {
  return apiClient.get(`/api/v1/admin/players/${id}/games`);
};
```

#### ❌ 管理陪玩师技能标签

```typescript
GET / admin / players / { id } / skill - tags;
POST / admin / players / { id } / skill - tags;
DELETE / admin / players / { id } / skill - tags / { tagId };
```

**状态**: 未实现  
**优先级**: 🟡 中  
**说明**: 陪玩师技能标签管理（如"高胜率"、"温柔"、"幽默"等）

**需要添加到**: `src/services/api/user.ts` (playerApi)

```typescript
getSkillTags: (id: number): Promise<SkillTag[]> => {
  return apiClient.get(`/api/v1/admin/players/${id}/skill-tags`);
},
addSkillTag: (id: number, data: { tag_name: string }): Promise<void> => {
  return apiClient.post(`/api/v1/admin/players/${id}/skill-tags`, data);
},
removeSkillTag: (id: number, tagId: number): Promise<void> => {
  return apiClient.delete(`/api/v1/admin/players/${id}/skill-tags/${tagId}`);
}
```

---

### 3. 统计报表模块

#### ❌ 仪表板统计

```typescript
GET / admin / stats / dashboard;
```

**状态**: 未实现  
**优先级**: 🔴 高  
**说明**: 管理后台首页 Dashboard 需要的综合统计数据

**需要创建**: `src/services/api/stats.ts`

```typescript
export const statsApi = {
  getDashboard: (): Promise<DashboardStats> => {
    return apiClient.get('/api/v1/admin/stats/dashboard');
  },
};
```

**预期返回数据结构**:

```typescript
interface DashboardStats {
  total_users: number;
  total_orders: number;
  total_revenue_cents: number;
  active_players: number;
  today_stats: {
    new_users: number;
    new_orders: number;
    revenue_cents: number;
  };
  // ...更多统计字段
}
```

#### ❌ 审计统计（概览）

```typescript
GET / admin / stats / audit / overview;
```

**状态**: 未实现  
**优先级**: 🟡 中  
**说明**: 审计模块的统计数据概览

#### ❌ 审计趋势

```typescript
GET / admin / stats / audit / trend;
```

**状态**: 未实现  
**优先级**: 🟡 中  
**说明**: 审计数据的时间趋势分析

#### ❌ 营收趋势

```typescript
GET / admin / stats / revenue - trend;
```

**状态**: 未实现  
**优先级**: 🔴 高  
**说明**: 营收趋势图表数据

#### ❌ TOP 陪玩师榜单

```typescript
GET / admin / stats / top - players;
```

**状态**: 未实现  
**优先级**: 🟡 中  
**说明**: 根据订单量、评分等维度的陪玩师排行榜

#### ❌ 用户增长趋势

```typescript
GET / admin / stats / user - growth;
```

**状态**: 未实现  
**优先级**: 🟡 中  
**说明**: 用户增长趋势图表数据

---

## 🔵 用户端接口（全部缺失，极高优先级）

### 概述

⚠️ **后端目前没有任何用户端接口！** 所有接口都是 `/admin/` 开头的管理端接口。

用户端前端需要独立的用户端接口，建议后端新增以下模块：

---

### 1. 用户端 - 游戏浏览

#### 🔵 获取游戏列表（用户端）

```typescript
GET / api / v1 / games;
```

**功能**: 用户浏览可预约的游戏列表
**筛选参数**:

- `category`: 游戏分类
- `keyword`: 关键词搜索
- `sort_by`: 排序方式（热度、新上架等）

#### 🔵 获取游戏详情（用户端）

```typescript
GET / api / v1 / games / { id };
```

**功能**: 查看游戏详细信息和可用陪玩师

---

### 2. 用户端 - 陪玩师浏览

#### 🔵 获取陪玩师列表（用户端）

```typescript
GET / api / v1 / players;
```

**功能**: 用户浏览陪玩师列表
**筛选参数**:

- `game_id`: 按游戏筛选
- `min_rating`: 最低评分
- `max_price_cents`: 最高价格
- `skill_tags[]`: 技能标签
- `sort_by`: 排序（评分、价格、订单量）

#### 🔵 获取陪玩师详情（用户端）

```typescript
GET / api / v1 / players / { id };
```

**功能**: 查看陪玩师详细资料、评价、擅长游戏

#### 🔵 获取陪玩师评价（用户端）

```typescript
GET / api / v1 / players / { id } / reviews;
```

**功能**: 查看陪玩师的用户评价列表

---

### 3. 用户端 - 订单管理

#### 🔵 创建订单（用户端）

```typescript
POST / api / v1 / orders;
```

**功能**: 用户预约陪玩服务

#### 🔵 我的订单列表（用户端）

```typescript
GET / api / v1 / me / orders;
```

**功能**: 查看当前用户的所有订单
**筛选参数**:

- `status`: 订单状态
- `date_from`, `date_to`: 时间范围

#### 🔵 订单详情（用户端）

```typescript
GET / api / v1 / orders / { id };
```

**功能**: 查看订单详细信息（需验证订单所属用户）

#### 🔵 取消订单（用户端）

```typescript
POST / api / v1 / orders / { id } / cancel;
```

**功能**: 用户取消自己的订单

#### 🔵 确认完成订单（用户端）

```typescript
POST / api / v1 / orders / { id } / complete;
```

**功能**: 用户确认订单已完成

---

### 4. 用户端 - 评价系统

#### 🔵 提交评价（用户端）

```typescript
POST / api / v1 / orders / { orderId } / review;
```

**功能**: 用户对完成的订单进行评价

**请求体**:

```typescript
{
  rating: number;        // 1-5 星
  comment: string;       // 评价内容
  tags: string[];        // 评价标签（如"准时"、"服务好"）
}
```

#### 🔵 修改评价（用户端）

```typescript
PUT / api / v1 / reviews / { id };
```

**功能**: 用户修改自己的评价

#### 🔵 删除评价（用户端）

```typescript
DELETE / api / v1 / reviews / { id };
```

**功能**: 用户删除自己的评价

---

### 5. 用户端 - 支付系统

#### 🔵 创建支付（用户端）

```typescript
POST / api / v1 / payments;
```

**功能**: 用户为订单发起支付

#### 🔵 支付状态查询（用户端）

```typescript
GET / api / v1 / payments / { id };
```

**功能**: 查询支付状态

#### 🔵 我的支付记录（用户端）

```typescript
GET / api / v1 / me / payments;
```

**功能**: 查看当前用户的支付记录

---

### 6. 用户端 - 个人中心

#### 🔵 获取个人信息（用户端）

```typescript
GET / api / v1 / me;
```

**功能**: 获取当前登录用户的信息

#### 🔵 更新个人信息（用户端）

```typescript
PUT / api / v1 / me;
```

**功能**: 更新个人资料（昵称、头像、手机等）

#### 🔵 修改密码（用户端）

```typescript
PUT / api / v1 / me / password;
```

**功能**: 用户修改密码

#### 🔵 我的收藏（陪玩师）

```typescript
GET / api / v1 / me / favorites;
POST / api / v1 / me / favorites;
DELETE / api / v1 / me / favorites / { playerId };
```

**功能**: 用户收藏喜欢的陪玩师

---

## 📋 前端需要创建的新文件

### 1. 统计 API 服务

**文件**: `src/services/api/stats.ts`

```typescript
// 管理端统计接口
export const statsApi = {
  getDashboard: () => { ... },
  getRevenueTrend: () => { ... },
  getUserGrowth: () => { ... },
  getTopPlayers: () => { ... },
  getAuditOverview: () => { ... },
  getAuditTrend: () => { ... },
}
```

### 2. 用户端 API 服务

**目录**: `src/services/api/client/`

建议将用户端接口单独组织：

```
src/services/api/
├── admin/           # 管理端接口（现有的）
│   ├── auth.ts
│   ├── game.ts
│   ├── order.ts
│   ├── payment.ts
│   ├── review.ts
│   ├── user.ts
│   └── stats.ts (NEW)
├── client/          # 用户端接口（NEW）
│   ├── auth.ts
│   ├── game.ts
│   ├── player.ts
│   ├── order.ts
│   ├── payment.ts
│   ├── review.ts
│   └── profile.ts
└── index.ts
```

### 3. 类型定义更新

**文件**:

- `src/types/stats.ts` (NEW) - 统计数据类型
- `src/types/client/` (NEW) - 用户端类型定义

---

## 🎯 实施优先级

### 阶段 1: 管理端补全（1-2天）

1. ✅ 实现注册接口 API 调用
2. ✅ 实现 Dashboard 统计接口
3. ✅ 实现陪玩师游戏列表接口
4. ⏸️ 实现营收趋势接口

### 阶段 2: 用户端核心功能（3-5天）

1. ✅ 游戏浏览
2. ✅ 陪玩师浏览和搜索
3. ✅ 订单创建和管理
4. ✅ 支付流程

### 阶段 3: 用户端增强功能（2-3天）

1. ✅ 评价系统
2. ✅ 个人中心
3. ✅ 收藏功能

### 阶段 4: 管理端高级统计（2-3天）

1. ✅ 审计统计
2. ✅ 用户增长分析
3. ✅ TOP 榜单

---

## 📝 注意事项

### 1. 权限控制

- 管理端接口（`/admin/*`）需要验证管理员权限
- 用户端接口（`/api/v1/*`）需要验证用户身份
- 部分用户端接口（如游戏列表）可能允许匿名访问

### 2. 数据隔离

- 用户端只能访问自己的订单、支付、评价
- 管理端可以访问所有数据

### 3. API 命名规范

- 管理端：`/admin/{resource}`
- 用户端：`/api/v1/{resource}` 或 `/api/v1/me/{resource}`

### 4. 前端路由规划

```
管理端（Admin）: /admin/*
  - /admin/dashboard
  - /admin/users
  - /admin/orders
  - /admin/games
  - /admin/players
  - /admin/payments
  - /admin/reports

用户端（Client）: /* 或 /app/*
  - /
  - /games
  - /players
  - /orders
  - /profile
  - /favorites
```

---

## ✅ 下一步行动

1. **后端开发**：新增用户端接口（37个）
2. **前端开发**：
   - 补全管理端缺失接口调用（9个）
   - 开发用户端完整 API 层（37个）
   - 开发用户端页面和组件

**预计总工作量**: 10-15 个工作日

---

**文档版本**: v1.0  
**创建时间**: 2025-10-28  
**最后更新**: 2025-10-28
