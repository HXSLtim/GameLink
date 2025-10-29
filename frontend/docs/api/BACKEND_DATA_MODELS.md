# 后端数据模型完整文档

> 数据来源: http://localhost:8080/swagger  
> 生成时间: 2025-10-28  
> 模型总数: 24 个

---

## 📋 目录

- [Auth 模块](#auth-模块)
- [User 模块](#user-模块)
- [Player 模块](#player-模块)
- [Order 模块](#order-模块)
- [Game 模块](#game-模块)
- [Payment 模块](#payment-模块)
- [Review 模块](#review-模块)
- [Common 模块](#common-模块)
- [前后端类型对比](#前后端类型对比)

---

## Auth 模块

### loginRequest

**完整名称**: `handler.loginRequest`

**用途**: 用户登录请求

| 字段名   | 类型   | 必填  | 说明   |
| -------- | ------ | ----- | ------ |
| username | string | ✅ 是 | 用户名 |
| password | string | ✅ 是 | 密码   |

**对应前端类型**: `src/types/auth.ts` - `LoginRequest`

```typescript
export interface LoginRequest {
  username: string;
  password: string;
}
```

---

### registerRequest

**完整名称**: `handler.registerRequest`

**用途**: 用户注册请求

| 字段名   | 类型   | 必填  | 说明     |
| -------- | ------ | ----- | -------- |
| name     | string | ✅ 是 | 用户姓名 |
| password | string | ✅ 是 | 密码     |
| email    | string | ⚪ 否 | 邮箱     |
| phone    | string | ⚪ 否 | 手机号   |

**对应前端类型**: `src/types/auth.ts` - `RegisterRequest`

```typescript
export interface RegisterRequest {
  email: string; // ⚠️ 前端为必填，后端为可选
  name: string;
  password: string;
  phone?: string;
}
```

⚠️ **差异**: 前端将 `email` 设为必填，但后端为可选

---

### loginResponse

**完整名称**: `handler.loginResponse`

**用途**: 登录成功响应

| 字段名     | 类型       | 必填  | 说明      |
| ---------- | ---------- | ----- | --------- |
| token      | string     | ⚪ 否 | JWT Token |
| expires_at | string     | ⚪ 否 | 过期时间  |
| user       | model.User | ⚪ 否 | 用户信息  |

**对应前端类型**: `src/types/auth.ts` - `LoginResult`

```typescript
export interface LoginResult {
  token: string;
  expires_at: string;
  user: User;
}
```

---

### tokenPayload

**完整名称**: `handler.tokenPayload`

**用途**: Token 刷新响应

| 字段名 | 类型   | 必填  | 说明           |
| ------ | ------ | ----- | -------------- |
| token  | string | ⚪ 否 | 新的 JWT Token |

---

## User 模块

### User (model.User)

**完整名称**: `model.User`

**用途**: 用户基础模型

| 字段名        | 类型             | 必填  | 说明                           |
| ------------- | ---------------- | ----- | ------------------------------ |
| id            | integer          | ⚪ 否 | 用户ID                         |
| name          | string           | ⚪ 否 | 姓名                           |
| email         | string           | ⚪ 否 | 邮箱                           |
| phone         | string           | ⚪ 否 | 手机号                         |
| avatar_url    | string           | ⚪ 否 | 头像URL                        |
| role          | model.Role       | ⚪ 否 | 角色 (user/player/admin)       |
| status        | model.UserStatus | ⚪ 否 | 状态 (active/suspended/banned) |
| last_login_at | string           | ⚪ 否 | 最后登录时间                   |
| created_at    | string           | ⚪ 否 | 创建时间                       |
| updated_at    | string           | ⚪ 否 | 更新时间                       |

**对应前端类型**: `src/types/user.ts` - `User`

---

### UserStatus (model.UserStatus)

**完整名称**: `model.UserStatus`

**用途**: 用户状态枚举

**可选值**:

- `active` - 活跃
- `suspended` - 暂停
- `banned` - 封禁

**对应前端类型**: `src/types/user.ts` - `UserStatus`

```typescript
export enum UserStatus {
  ACTIVE = 'active',
  SUSPENDED = 'suspended',
  BANNED = 'banned',
}
```

---

### CreateUserPayload

**完整名称**: `admin.CreateUserPayload`

**用途**: 管理员创建用户

| 字段名     | 类型   | 必填  | 说明    |
| ---------- | ------ | ----- | ------- |
| name       | string | ✅ 是 | 姓名    |
| password   | string | ✅ 是 | 密码    |
| role       | string | ✅ 是 | 角色    |
| status     | string | ✅ 是 | 状态    |
| email      | string | ⚪ 否 | 邮箱    |
| phone      | string | ⚪ 否 | 手机号  |
| avatar_url | string | ⚪ 否 | 头像URL |

**对应前端类型**: `src/types/user.ts` - `CreateUserRequest`

---

### UpdateUserPayload

**完整名称**: `admin.UpdateUserPayload`

**用途**: 管理员更新用户

| 字段名     | 类型   | 必填  | 说明             |
| ---------- | ------ | ----- | ---------------- |
| name       | string | ✅ 是 | 姓名             |
| role       | string | ✅ 是 | 角色             |
| status     | string | ✅ 是 | 状态             |
| email      | string | ⚪ 否 | 邮箱             |
| phone      | string | ⚪ 否 | 手机号           |
| avatar_url | string | ⚪ 否 | 头像URL          |
| password   | string | ⚪ 否 | 密码（可选更新） |

**对应前端类型**: `src/types/user.ts` - `UpdateUserRequest`

---

## Player 模块

### CreatePlayerPayload

**完整名称**: `admin.CreatePlayerPayload`

**用途**: 创建陪玩师

| 字段名              | 类型    | 必填  | 说明       |
| ------------------- | ------- | ----- | ---------- |
| user_id             | integer | ✅ 是 | 关联用户ID |
| verification_status | string  | ✅ 是 | 认证状态   |
| nickname            | string  | ⚪ 否 | 昵称       |
| bio                 | string  | ⚪ 否 | 个人简介   |
| main_game_id        | integer | ⚪ 否 | 主玩游戏ID |
| hourly_rate_cents   | integer | ⚪ 否 | 时薪（分） |

**对应前端类型**: `src/types/user.ts` - `CreatePlayerRequest`

---

### UpdatePlayerPayload

**完整名称**: `admin.UpdatePlayerPayload`

**用途**: 更新陪玩师信息

| 字段名              | 类型    | 必填  | 说明       |
| ------------------- | ------- | ----- | ---------- |
| verification_status | string  | ✅ 是 | 认证状态   |
| nickname            | string  | ⚪ 否 | 昵称       |
| bio                 | string  | ⚪ 否 | 个人简介   |
| main_game_id        | integer | ⚪ 否 | 主玩游戏ID |
| hourly_rate_cents   | integer | ⚪ 否 | 时薪（分） |

**对应前端类型**: `src/types/user.ts` - `UpdatePlayerRequest`

---

## Order 模块

### CreateOrderPayload

**完整名称**: `admin.CreateOrderPayload`

**用途**: 创建订单

| 字段名          | 类型    | 必填  | 说明                         |
| --------------- | ------- | ----- | ---------------------------- |
| user_id         | integer | ✅ 是 | 用户ID                       |
| game_id         | integer | ✅ 是 | 游戏ID                       |
| price_cents     | integer | ✅ 是 | 价格（分）                   |
| currency        | string  | ✅ 是 | 货币代码                     |
| player_id       | integer | ⚪ 否 | 陪玩师ID（可预约指定陪玩师） |
| title           | string  | ⚪ 否 | 订单标题                     |
| description     | string  | ⚪ 否 | 订单描述                     |
| scheduled_start | string  | ⚪ 否 | 预约开始时间                 |
| scheduled_end   | string  | ⚪ 否 | 预约结束时间                 |

**对应前端类型**: `src/types/order.ts` - `CreateOrderRequest`

✅ **一致性**: 前端类型与后端完全一致

---

### UpdateOrderPayload

**完整名称**: `admin.UpdateOrderPayload`

**用途**: 更新订单

| 字段名          | 类型    | 必填  | 说明         |
| --------------- | ------- | ----- | ------------ |
| price_cents     | integer | ✅ 是 | 价格（分）   |
| currency        | string  | ✅ 是 | 货币代码     |
| status          | string  | ✅ 是 | 订单状态     |
| scheduled_start | string  | ⚪ 否 | 预约开始时间 |
| scheduled_end   | string  | ⚪ 否 | 预约结束时间 |
| cancel_reason   | string  | ⚪ 否 | 取消原因     |

**对应前端类型**: `src/types/order.ts` - `UpdateOrderRequest`

✅ **一致性**: 前端类型与后端完全一致

---

### AssignOrderPayload

**完整名称**: `admin.AssignOrderPayload`

**用途**: 分配订单给陪玩师

| 字段名    | 类型    | 必填  | 说明     |
| --------- | ------- | ----- | -------- |
| player_id | integer | ✅ 是 | 陪玩师ID |

**对应前端类型**: `src/types/order.ts` - `AssignOrderRequest`

---

### ReviewOrderPayload

**完整名称**: `admin.ReviewOrderPayload`

**用途**: 审核订单

| 字段名   | 类型    | 必填  | 说明                              |
| -------- | ------- | ----- | --------------------------------- |
| approved | boolean | ⚪ 否 | 是否通过（true=通过，false=拒绝） |
| reason   | string  | ⚪ 否 | 审核理由/拒绝原因                 |

**对应前端类型**: `src/types/order.ts` - `ReviewOrderRequest`

✅ **一致性**: 前端类型与后端完全一致

---

### CancelOrderPayload

**完整名称**: `admin.CancelOrderPayload`

**用途**: 取消订单

| 字段名 | 类型   | 必填  | 说明     |
| ------ | ------ | ----- | -------- |
| reason | string | ⚪ 否 | 取消原因 |

**对应前端类型**: `src/types/order.ts` - `CancelOrderRequest`

✅ **一致性**: 前端类型与后端完全一致

---

## Game 模块

### GamePayload

**完整名称**: `admin.GamePayload`

**用途**: 创建/更新游戏

| 字段名      | 类型   | 必填  | 说明         |
| ----------- | ------ | ----- | ------------ |
| key         | string | ✅ 是 | 游戏唯一标识 |
| name        | string | ✅ 是 | 游戏名称     |
| category    | string | ⚪ 否 | 游戏分类     |
| description | string | ⚪ 否 | 游戏描述     |
| icon_url    | string | ⚪ 否 | 游戏图标URL  |

**对应前端类型**:

- `src/types/game.ts` - `CreateGameRequest`
- `src/types/game.ts` - `UpdateGameRequest`

---

## Payment 模块

### CreatePaymentPayload

**完整名称**: `admin.CreatePaymentPayload`

**用途**: 创建支付

⚠️ **注意**: Swagger 中未定义具体字段

---

### UpdatePaymentPayload

**完整名称**: `admin.UpdatePaymentPayload`

**用途**: 更新支付信息

⚠️ **注意**: Swagger 中未定义具体字段

---

### CapturePaymentPayload

**完整名称**: `admin.CapturePaymentPayload`

**用途**: 确认收款

⚠️ **注意**: Swagger 中未定义具体字段

---

### RefundPaymentPayload

**完整名称**: `admin.RefundPaymentPayload`

**用途**: 申请退款

⚠️ **注意**: Swagger 中未定义具体字段

---

## Review 模块

### CreateReviewPayload

**完整名称**: `admin.CreateReviewPayload`

**用途**: 创建评价

| 字段名    | 类型    | 必填  | 说明     |
| --------- | ------- | ----- | -------- |
| user_id   | integer | ✅ 是 | 用户ID   |
| player_id | integer | ✅ 是 | 陪玩师ID |
| order_id  | integer | ✅ 是 | 订单ID   |
| score     | integer | ✅ 是 | 评分     |
| content   | string  | ⚪ 否 | 评价内容 |

**对应前端类型**: `src/types/review.ts` - `CreateReviewRequest`

---

### UpdateReviewPayload

**完整名称**: `admin.UpdateReviewPayload`

**用途**: 更新评价

| 字段名  | 类型    | 必填  | 说明     |
| ------- | ------- | ----- | -------- |
| score   | integer | ✅ 是 | 评分     |
| content | string  | ⚪ 否 | 评价内容 |

**对应前端类型**: `src/types/review.ts` - `UpdateReviewRequest`

---

## Common 模块

### Role (model.Role)

**完整名称**: `model.Role`

**用途**: 用户角色枚举

**可选值**:

- `user` - 普通用户
- `player` - 陪玩师
- `admin` - 管理员

**对应前端类型**: `src/types/user.ts` - `UserRole`

```typescript
export enum UserRole {
  USER = 'user',
  PLAYER = 'player',
  ADMIN = 'admin',
}
```

---

### SkillTagsBody

**完整名称**: `admin.SkillTagsBody`

**用途**: 陪玩师技能标签

| 字段名 | 类型          | 必填  | 说明         |
| ------ | ------------- | ----- | ------------ |
| tags   | array<string> | ✅ 是 | 技能标签数组 |

---

## 前后端类型对比

### ✅ 完全一致的类型

| 模块  | 类型               | 前端文件             | 状态    |
| ----- | ------------------ | -------------------- | ------- |
| Order | CreateOrderRequest | `src/types/order.ts` | ✅ 一致 |
| Order | UpdateOrderRequest | `src/types/order.ts` | ✅ 一致 |
| Order | ReviewOrderRequest | `src/types/order.ts` | ✅ 一致 |
| Order | CancelOrderRequest | `src/types/order.ts` | ✅ 一致 |
| Auth  | LoginRequest       | `src/types/auth.ts`  | ✅ 一致 |
| User  | User               | `src/types/user.ts`  | ✅ 一致 |
| User  | UserStatus         | `src/types/user.ts`  | ✅ 一致 |
| User  | UserRole           | `src/types/user.ts`  | ✅ 一致 |

---

### ⚠️ 存在差异的类型

| 模块    | 类型            | 差异说明                        | 建议             |
| ------- | --------------- | ------------------------------- | ---------------- |
| Auth    | RegisterRequest | 前端 `email` 为必填，后端为可选 | 建议统一为可选   |
| Payment | 所有 Payload    | 后端 Swagger 未定义字段         | 需要后端补充文档 |

---

### ❌ 前端缺失的类型

| 模块   | 后端类型      | 说明           | 优先级 |
| ------ | ------------- | -------------- | ------ |
| Player | SkillTagsBody | 技能标签管理   | 🟡 中  |
| Auth   | tokenPayload  | Token 刷新响应 | 🟢 低  |

---

### 📝 建议操作

#### 1. 修复 RegisterRequest 不一致

**文件**: `src/types/auth.ts`

```typescript
// 建议修改为（与后端一致）
export interface RegisterRequest {
  name: string;
  password: string;
  email?: string; // 改为可选
  phone?: string;
}
```

#### 2. 添加 SkillTagsBody 类型

**文件**: `src/types/user.ts`

```typescript
// 新增
export interface SkillTagsBody {
  tags: string[];
}
```

#### 3. 补充 Payment 相关类型定义

需要后端在 Swagger 中补充以下 Payload 的字段定义：

- CreatePaymentPayload
- UpdatePaymentPayload
- CapturePaymentPayload
- RefundPaymentPayload

---

## 📊 类型同步总结

| 状态        | 数量 | 百分比 |
| ----------- | ---- | ------ |
| ✅ 完全一致 | 8    | 73%    |
| ⚠️ 存在差异 | 2    | 18%    |
| ❌ 缺失     | 1    | 9%     |

**总体评估**: 前后端类型定义基本一致，存在少量差异需要修复。

---

**文档版本**: v1.0  
**生成时间**: 2025-10-28  
**数据来源**: http://localhost:8080/swagger  
**维护**: 需要与后端 Swagger 文档保持同步

> 数据来源: http://localhost:8080/swagger  
> 生成时间: 2025-10-28  
> 模型总数: 24 个

---

## 📋 目录

- [Auth 模块](#auth-模块)
- [User 模块](#user-模块)
- [Player 模块](#player-模块)
- [Order 模块](#order-模块)
- [Game 模块](#game-模块)
- [Payment 模块](#payment-模块)
- [Review 模块](#review-模块)
- [Common 模块](#common-模块)
- [前后端类型对比](#前后端类型对比)

---

## Auth 模块

### loginRequest

**完整名称**: `handler.loginRequest`

**用途**: 用户登录请求

| 字段名   | 类型   | 必填  | 说明   |
| -------- | ------ | ----- | ------ |
| username | string | ✅ 是 | 用户名 |
| password | string | ✅ 是 | 密码   |

**对应前端类型**: `src/types/auth.ts` - `LoginRequest`

```typescript
export interface LoginRequest {
  username: string;
  password: string;
}
```

---

### registerRequest

**完整名称**: `handler.registerRequest`

**用途**: 用户注册请求

| 字段名   | 类型   | 必填  | 说明     |
| -------- | ------ | ----- | -------- |
| name     | string | ✅ 是 | 用户姓名 |
| password | string | ✅ 是 | 密码     |
| email    | string | ⚪ 否 | 邮箱     |
| phone    | string | ⚪ 否 | 手机号   |

**对应前端类型**: `src/types/auth.ts` - `RegisterRequest`

```typescript
export interface RegisterRequest {
  email: string; // ⚠️ 前端为必填，后端为可选
  name: string;
  password: string;
  phone?: string;
}
```

⚠️ **差异**: 前端将 `email` 设为必填，但后端为可选

---

### loginResponse

**完整名称**: `handler.loginResponse`

**用途**: 登录成功响应

| 字段名     | 类型       | 必填  | 说明      |
| ---------- | ---------- | ----- | --------- |
| token      | string     | ⚪ 否 | JWT Token |
| expires_at | string     | ⚪ 否 | 过期时间  |
| user       | model.User | ⚪ 否 | 用户信息  |

**对应前端类型**: `src/types/auth.ts` - `LoginResult`

```typescript
export interface LoginResult {
  token: string;
  expires_at: string;
  user: User;
}
```

---

### tokenPayload

**完整名称**: `handler.tokenPayload`

**用途**: Token 刷新响应

| 字段名 | 类型   | 必填  | 说明           |
| ------ | ------ | ----- | -------------- |
| token  | string | ⚪ 否 | 新的 JWT Token |

---

## User 模块

### User (model.User)

**完整名称**: `model.User`

**用途**: 用户基础模型

| 字段名        | 类型             | 必填  | 说明                           |
| ------------- | ---------------- | ----- | ------------------------------ |
| id            | integer          | ⚪ 否 | 用户ID                         |
| name          | string           | ⚪ 否 | 姓名                           |
| email         | string           | ⚪ 否 | 邮箱                           |
| phone         | string           | ⚪ 否 | 手机号                         |
| avatar_url    | string           | ⚪ 否 | 头像URL                        |
| role          | model.Role       | ⚪ 否 | 角色 (user/player/admin)       |
| status        | model.UserStatus | ⚪ 否 | 状态 (active/suspended/banned) |
| last_login_at | string           | ⚪ 否 | 最后登录时间                   |
| created_at    | string           | ⚪ 否 | 创建时间                       |
| updated_at    | string           | ⚪ 否 | 更新时间                       |

**对应前端类型**: `src/types/user.ts` - `User`

---

### UserStatus (model.UserStatus)

**完整名称**: `model.UserStatus`

**用途**: 用户状态枚举

**可选值**:

- `active` - 活跃
- `suspended` - 暂停
- `banned` - 封禁

**对应前端类型**: `src/types/user.ts` - `UserStatus`

```typescript
export enum UserStatus {
  ACTIVE = 'active',
  SUSPENDED = 'suspended',
  BANNED = 'banned',
}
```

---

### CreateUserPayload

**完整名称**: `admin.CreateUserPayload`

**用途**: 管理员创建用户

| 字段名     | 类型   | 必填  | 说明    |
| ---------- | ------ | ----- | ------- |
| name       | string | ✅ 是 | 姓名    |
| password   | string | ✅ 是 | 密码    |
| role       | string | ✅ 是 | 角色    |
| status     | string | ✅ 是 | 状态    |
| email      | string | ⚪ 否 | 邮箱    |
| phone      | string | ⚪ 否 | 手机号  |
| avatar_url | string | ⚪ 否 | 头像URL |

**对应前端类型**: `src/types/user.ts` - `CreateUserRequest`

---

### UpdateUserPayload

**完整名称**: `admin.UpdateUserPayload`

**用途**: 管理员更新用户

| 字段名     | 类型   | 必填  | 说明             |
| ---------- | ------ | ----- | ---------------- |
| name       | string | ✅ 是 | 姓名             |
| role       | string | ✅ 是 | 角色             |
| status     | string | ✅ 是 | 状态             |
| email      | string | ⚪ 否 | 邮箱             |
| phone      | string | ⚪ 否 | 手机号           |
| avatar_url | string | ⚪ 否 | 头像URL          |
| password   | string | ⚪ 否 | 密码（可选更新） |

**对应前端类型**: `src/types/user.ts` - `UpdateUserRequest`

---

## Player 模块

### CreatePlayerPayload

**完整名称**: `admin.CreatePlayerPayload`

**用途**: 创建陪玩师

| 字段名              | 类型    | 必填  | 说明       |
| ------------------- | ------- | ----- | ---------- |
| user_id             | integer | ✅ 是 | 关联用户ID |
| verification_status | string  | ✅ 是 | 认证状态   |
| nickname            | string  | ⚪ 否 | 昵称       |
| bio                 | string  | ⚪ 否 | 个人简介   |
| main_game_id        | integer | ⚪ 否 | 主玩游戏ID |
| hourly_rate_cents   | integer | ⚪ 否 | 时薪（分） |

**对应前端类型**: `src/types/user.ts` - `CreatePlayerRequest`

---

### UpdatePlayerPayload

**完整名称**: `admin.UpdatePlayerPayload`

**用途**: 更新陪玩师信息

| 字段名              | 类型    | 必填  | 说明       |
| ------------------- | ------- | ----- | ---------- |
| verification_status | string  | ✅ 是 | 认证状态   |
| nickname            | string  | ⚪ 否 | 昵称       |
| bio                 | string  | ⚪ 否 | 个人简介   |
| main_game_id        | integer | ⚪ 否 | 主玩游戏ID |
| hourly_rate_cents   | integer | ⚪ 否 | 时薪（分） |

**对应前端类型**: `src/types/user.ts` - `UpdatePlayerRequest`

---

## Order 模块

### CreateOrderPayload

**完整名称**: `admin.CreateOrderPayload`

**用途**: 创建订单

| 字段名          | 类型    | 必填  | 说明                         |
| --------------- | ------- | ----- | ---------------------------- |
| user_id         | integer | ✅ 是 | 用户ID                       |
| game_id         | integer | ✅ 是 | 游戏ID                       |
| price_cents     | integer | ✅ 是 | 价格（分）                   |
| currency        | string  | ✅ 是 | 货币代码                     |
| player_id       | integer | ⚪ 否 | 陪玩师ID（可预约指定陪玩师） |
| title           | string  | ⚪ 否 | 订单标题                     |
| description     | string  | ⚪ 否 | 订单描述                     |
| scheduled_start | string  | ⚪ 否 | 预约开始时间                 |
| scheduled_end   | string  | ⚪ 否 | 预约结束时间                 |

**对应前端类型**: `src/types/order.ts` - `CreateOrderRequest`

✅ **一致性**: 前端类型与后端完全一致

---

### UpdateOrderPayload

**完整名称**: `admin.UpdateOrderPayload`

**用途**: 更新订单

| 字段名          | 类型    | 必填  | 说明         |
| --------------- | ------- | ----- | ------------ |
| price_cents     | integer | ✅ 是 | 价格（分）   |
| currency        | string  | ✅ 是 | 货币代码     |
| status          | string  | ✅ 是 | 订单状态     |
| scheduled_start | string  | ⚪ 否 | 预约开始时间 |
| scheduled_end   | string  | ⚪ 否 | 预约结束时间 |
| cancel_reason   | string  | ⚪ 否 | 取消原因     |

**对应前端类型**: `src/types/order.ts` - `UpdateOrderRequest`

✅ **一致性**: 前端类型与后端完全一致

---

### AssignOrderPayload

**完整名称**: `admin.AssignOrderPayload`

**用途**: 分配订单给陪玩师

| 字段名    | 类型    | 必填  | 说明     |
| --------- | ------- | ----- | -------- |
| player_id | integer | ✅ 是 | 陪玩师ID |

**对应前端类型**: `src/types/order.ts` - `AssignOrderRequest`

---

### ReviewOrderPayload

**完整名称**: `admin.ReviewOrderPayload`

**用途**: 审核订单

| 字段名   | 类型    | 必填  | 说明                              |
| -------- | ------- | ----- | --------------------------------- |
| approved | boolean | ⚪ 否 | 是否通过（true=通过，false=拒绝） |
| reason   | string  | ⚪ 否 | 审核理由/拒绝原因                 |

**对应前端类型**: `src/types/order.ts` - `ReviewOrderRequest`

✅ **一致性**: 前端类型与后端完全一致

---

### CancelOrderPayload

**完整名称**: `admin.CancelOrderPayload`

**用途**: 取消订单

| 字段名 | 类型   | 必填  | 说明     |
| ------ | ------ | ----- | -------- |
| reason | string | ⚪ 否 | 取消原因 |

**对应前端类型**: `src/types/order.ts` - `CancelOrderRequest`

✅ **一致性**: 前端类型与后端完全一致

---

## Game 模块

### GamePayload

**完整名称**: `admin.GamePayload`

**用途**: 创建/更新游戏

| 字段名      | 类型   | 必填  | 说明         |
| ----------- | ------ | ----- | ------------ |
| key         | string | ✅ 是 | 游戏唯一标识 |
| name        | string | ✅ 是 | 游戏名称     |
| category    | string | ⚪ 否 | 游戏分类     |
| description | string | ⚪ 否 | 游戏描述     |
| icon_url    | string | ⚪ 否 | 游戏图标URL  |

**对应前端类型**:

- `src/types/game.ts` - `CreateGameRequest`
- `src/types/game.ts` - `UpdateGameRequest`

---

## Payment 模块

### CreatePaymentPayload

**完整名称**: `admin.CreatePaymentPayload`

**用途**: 创建支付

⚠️ **注意**: Swagger 中未定义具体字段

---

### UpdatePaymentPayload

**完整名称**: `admin.UpdatePaymentPayload`

**用途**: 更新支付信息

⚠️ **注意**: Swagger 中未定义具体字段

---

### CapturePaymentPayload

**完整名称**: `admin.CapturePaymentPayload`

**用途**: 确认收款

⚠️ **注意**: Swagger 中未定义具体字段

---

### RefundPaymentPayload

**完整名称**: `admin.RefundPaymentPayload`

**用途**: 申请退款

⚠️ **注意**: Swagger 中未定义具体字段

---

## Review 模块

### CreateReviewPayload

**完整名称**: `admin.CreateReviewPayload`

**用途**: 创建评价

| 字段名    | 类型    | 必填  | 说明     |
| --------- | ------- | ----- | -------- |
| user_id   | integer | ✅ 是 | 用户ID   |
| player_id | integer | ✅ 是 | 陪玩师ID |
| order_id  | integer | ✅ 是 | 订单ID   |
| score     | integer | ✅ 是 | 评分     |
| content   | string  | ⚪ 否 | 评价内容 |

**对应前端类型**: `src/types/review.ts` - `CreateReviewRequest`

---

### UpdateReviewPayload

**完整名称**: `admin.UpdateReviewPayload`

**用途**: 更新评价

| 字段名  | 类型    | 必填  | 说明     |
| ------- | ------- | ----- | -------- |
| score   | integer | ✅ 是 | 评分     |
| content | string  | ⚪ 否 | 评价内容 |

**对应前端类型**: `src/types/review.ts` - `UpdateReviewRequest`

---

## Common 模块

### Role (model.Role)

**完整名称**: `model.Role`

**用途**: 用户角色枚举

**可选值**:

- `user` - 普通用户
- `player` - 陪玩师
- `admin` - 管理员

**对应前端类型**: `src/types/user.ts` - `UserRole`

```typescript
export enum UserRole {
  USER = 'user',
  PLAYER = 'player',
  ADMIN = 'admin',
}
```

---

### SkillTagsBody

**完整名称**: `admin.SkillTagsBody`

**用途**: 陪玩师技能标签

| 字段名 | 类型          | 必填  | 说明         |
| ------ | ------------- | ----- | ------------ |
| tags   | array<string> | ✅ 是 | 技能标签数组 |

---

## 前后端类型对比

### ✅ 完全一致的类型

| 模块  | 类型               | 前端文件             | 状态    |
| ----- | ------------------ | -------------------- | ------- |
| Order | CreateOrderRequest | `src/types/order.ts` | ✅ 一致 |
| Order | UpdateOrderRequest | `src/types/order.ts` | ✅ 一致 |
| Order | ReviewOrderRequest | `src/types/order.ts` | ✅ 一致 |
| Order | CancelOrderRequest | `src/types/order.ts` | ✅ 一致 |
| Auth  | LoginRequest       | `src/types/auth.ts`  | ✅ 一致 |
| User  | User               | `src/types/user.ts`  | ✅ 一致 |
| User  | UserStatus         | `src/types/user.ts`  | ✅ 一致 |
| User  | UserRole           | `src/types/user.ts`  | ✅ 一致 |

---

### ⚠️ 存在差异的类型

| 模块    | 类型            | 差异说明                        | 建议             |
| ------- | --------------- | ------------------------------- | ---------------- |
| Auth    | RegisterRequest | 前端 `email` 为必填，后端为可选 | 建议统一为可选   |
| Payment | 所有 Payload    | 后端 Swagger 未定义字段         | 需要后端补充文档 |

---

### ❌ 前端缺失的类型

| 模块   | 后端类型      | 说明           | 优先级 |
| ------ | ------------- | -------------- | ------ |
| Player | SkillTagsBody | 技能标签管理   | 🟡 中  |
| Auth   | tokenPayload  | Token 刷新响应 | 🟢 低  |

---

### 📝 建议操作

#### 1. 修复 RegisterRequest 不一致

**文件**: `src/types/auth.ts`

```typescript
// 建议修改为（与后端一致）
export interface RegisterRequest {
  name: string;
  password: string;
  email?: string; // 改为可选
  phone?: string;
}
```

#### 2. 添加 SkillTagsBody 类型

**文件**: `src/types/user.ts`

```typescript
// 新增
export interface SkillTagsBody {
  tags: string[];
}
```

#### 3. 补充 Payment 相关类型定义

需要后端在 Swagger 中补充以下 Payload 的字段定义：

- CreatePaymentPayload
- UpdatePaymentPayload
- CapturePaymentPayload
- RefundPaymentPayload

---

## 📊 类型同步总结

| 状态        | 数量 | 百分比 |
| ----------- | ---- | ------ |
| ✅ 完全一致 | 8    | 73%    |
| ⚠️ 存在差异 | 2    | 18%    |
| ❌ 缺失     | 1    | 9%     |

**总体评估**: 前后端类型定义基本一致，存在少量差异需要修复。

---

**文档版本**: v1.0  
**生成时间**: 2025-10-28  
**数据来源**: http://localhost:8080/swagger  
**维护**: 需要与后端 Swagger 文档保持同步
