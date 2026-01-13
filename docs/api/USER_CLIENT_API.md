# 用户端 API 接口文档

> 面向 Client (Desktop PWA) 和 App (小程序) 的用户端接口规范

## 通用说明

### 基础路径
```
/api/v1/user
```

### 认证方式
所有接口（除登录/注册外）需要在 Header 中携带 JWT Token：
```
Authorization: Bearer {token}
```

### 统一响应格式
```json
{
  "success": true,
  "code": 200,
  "message": "OK",
  "data": { ... }
}
```

### 错误响应格式
```json
{
  "success": false,
  "code": 400,
  "message": "错误描述",
  "details": "详细错误信息（可选）"
}
```

### 金额单位
所有金额字段均以 **分 (cents)** 为单位，前端需自行转换为元显示。

---

## 1. 钱包模块 (`/wallet`)

### 1.1 获取钱包余额

**Endpoint:** `GET /user/wallet/balance`

**描述:** 获取当前用户的钱包余额信息

**请求头:**
```
Authorization: Bearer {token}
```

**响应:**
```json
{
  "success": true,
  "code": 200,
  "message": "OK",
  "data": {
    "balanceCents": 125050,
    "frozenCents": 10000
  }
}
```

**字段说明:**
| 字段 | 类型 | 说明 |
|------|------|------|
| balanceCents | int64 | 可用余额（分） |
| frozenCents | int64 | 冻结金额（分），T+7 结算期内的收益 |

---

### 1.2 获取交易记录

**Endpoint:** `GET /user/wallet/transactions`

**描述:** 获取用户的钱包交易记录（充值、支付、退款、提现）

**请求参数 (Query):**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| pageSize | int | 否 | 每页数量，默认 20，最大 100 |
| type | string | 否 | 交易类型过滤：recharge/payment/refund/withdrawal |
| status | string | 否 | 状态过滤：pending/success/failed |

**响应:**
```json
{
  "success": true,
  "code": 200,
  "message": "OK",
  "data": {
    "items": [
      {
        "id": 1001,
        "type": "recharge",
        "amountCents": 50000,
        "status": "success",
        "description": "支付宝充值",
        "balanceAfterCents": 125050,
        "createdAt": "2024-01-15T10:30:00Z"
      },
      {
        "id": 1002,
        "type": "payment",
        "amountCents": -20000,
        "status": "success",
        "description": "订单支付 #12345",
        "orderId": 12345,
        "balanceAfterCents": 105050,
        "createdAt": "2024-01-15T14:20:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "pageSize": 20,
      "total": 45,
      "totalPages": 3
    }
  }
}
```

**交易类型 (type):**
| 值 | 说明 |
|------|------|
| recharge | 充值 |
| payment | 订单支付 |
| refund | 退款 |
| withdrawal | 提现 |

**状态 (status):**
| 值 | 说明 |
|------|------|
| pending | 处理中 |
| success | 成功 |
| failed | 失败 |

---

### 1.3 钱包充值

**Endpoint:** `POST /user/wallet/recharge`

**描述:** 发起钱包充值请求

**请求体:**
```json
{
  "amountCents": 10000,
  "method": "alipay"
}
```

**请求字段:**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| amountCents | int64 | 是 | 充值金额（分），最小 1 |
| method | string | 是 | 支付方式：wechat / alipay |

**响应:**
```json
{
  "success": true,
  "code": 200,
  "message": "OK",
  "data": {
    "orderId": 12345,
    "paymentId": 67890,
    "balanceCents": 135050,
    "payUrl": "https://pay.example.com/xxx"
  }
}
```

**错误码:**
| code | message | 说明 |
|------|---------|------|
| 400 | invalid amount | 金额无效 |
| 401 | unauthorized | 未登录 |
| 500 | recharge failed | 充值失败 |

---

### 1.4 申请提现

**Endpoint:** `POST /user/wallet/withdraw`

**描述:** 申请提现到银行卡

**请求体:**
```json
{
  "amountCents": 50000,
  "bankCardId": 123
}
```

**请求字段:**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| amountCents | int64 | 是 | 提现金额（分） |
| bankCardId | uint64 | 是 | 绑定的银行卡 ID |

**响应:**
```json
{
  "success": true,
  "code": 200,
  "message": "提现申请已提交",
  "data": {
    "withdrawId": 789,
    "status": "pending",
    "estimatedArrival": "2024-01-18T00:00:00Z"
  }
}
```

---

## 2. VIP 模块 (`/vip`)

### 2.1 获取 VIP 等级列表

**Endpoint:** `GET /user/vip/levels`

**描述:** 获取所有启用的 VIP 等级列表

**响应:**
```json
{
  "success": true,
  "code": 200,
  "message": "OK",
  "data": [
    {
      "id": 1,
      "slug": "bronze",
      "name": "青铜会员",
      "level": 1,
      "minExp": 0,
      "maxExp": 999,
      "discountRate": 0,
      "monthlyTickets": 0,
      "icon": "https://cdn.example.com/vip/bronze.png",
      "color": "#CD7F32",
      "benefits": ["专属客服", "生日礼包"],
      "isActive": true
    },
    {
      "id": 2,
      "slug": "silver",
      "name": "白银会员",
      "level": 2,
      "minExp": 1000,
      "maxExp": 4999,
      "discountRate": 5,
      "monthlyTickets": 1,
      "icon": "https://cdn.example.com/vip/silver.png",
      "color": "#C0C0C0",
      "benefits": ["专属客服", "生日礼包", "优先匹配"],
      "isActive": true
    }
  ]
}
```

---

### 2.2 获取用户 VIP 信息

**Endpoint:** `GET /user/vip/info`

**描述:** 获取当前用户的 VIP 等级、经验值和权益信息

**响应:**
```json
{
  "success": true,
  "code": 200,
  "message": "OK",
  "data": {
    "vipUnlocked": true,
    "currentLevel": {
      "id": 2,
      "slug": "silver",
      "name": "白银会员",
      "level": 2,
      "icon": "https://cdn.example.com/vip/silver.png",
      "color": "#C0C0C0"
    },
    "currentExp": 2450,
    "nextLevelExp": 5000,
    "expProgress": 0.49,
    "vipUnlockedAt": "2024-01-01T00:00:00Z",
    "vipExpireAt": null,
    "benefits": [
      { "id": "b1", "name": "专属客服", "icon": "customer-service" },
      { "id": "b2", "name": "生日礼包", "icon": "gift" },
      { "id": "b3", "name": "优先匹配", "icon": "thunderbolt" }
    ],
    "monthlyTicketsRemaining": 1,
    "discountRate": 5
  }
}
```

**字段说明:**
| 字段 | 类型 | 说明 |
|------|------|------|
| vipUnlocked | bool | VIP 是否已解锁 |
| currentLevel | object | 当前 VIP 等级信息 |
| currentExp | int64 | 当前经验值（累计消费，分） |
| nextLevelExp | int64 | 下一等级所需经验值 |
| expProgress | float | 升级进度 (0-1) |
| vipExpireAt | string/null | VIP 过期时间，null 表示永久 |
| benefits | array | 当前等级权益列表 |
| monthlyTicketsRemaining | int | 本月剩余优惠券数量 |
| discountRate | int | 折扣比例 (%) |

---

### 2.3 获取 VIP 解锁门槛

**Endpoint:** `GET /user/vip/threshold`

**描述:** 获取 VIP 解锁所需的消费/充值门槛

**响应:**
```json
{
  "success": true,
  "code": 200,
  "message": "OK",
  "data": {
    "consumeThreshold": 10000,
    "rechargeThreshold": 5000
  }
}
```

**说明:** 满足任一条件即可解锁 VIP（单位：分）

---

## 3. 订单模块 (`/orders`)

### 3.1 获取订单统计

**Endpoint:** `GET /user/orders/stats`

**描述:** 获取用户订单统计数据，用于个人中心展示

**响应:**
```json
{
  "success": true,
  "code": 200,
  "message": "OK",
  "data": {
    "totalCount": 45,
    "monthlyCount": 12,
    "monthlyChange": 2,
    "pendingCount": 1,
    "inProgressCount": 2,
    "completedCount": 38,
    "canceledCount": 4,
    "totalSpentCents": 450000,
    "avgOrderAmountCents": 10000
  }
}
```

**字段说明:**
| 字段 | 类型 | 说明 |
|------|------|------|
| totalCount | int | 总订单数 |
| monthlyCount | int | 本月订单数 |
| monthlyChange | int | 相比上月变化 |
| pendingCount | int | 待处理订单数（可用于角标） |
| inProgressCount | int | 进行中订单数 |
| completedCount | int | 已完成订单数 |
| canceledCount | int | 已取消订单数 |
| totalSpentCents | int64 | 累计消费（分） |
| avgOrderAmountCents | int64 | 平均订单金额（分） |

---

### 3.2 获取订单列表

**Endpoint:** `GET /user/orders`

**描述:** 获取当前用户的订单列表

**请求参数 (Query):**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| status | string | 否 | 状态过滤 |
| page | int | 否 | 页码，默认 1 |
| pageSize | int | 否 | 每页数量，默认 20 |

**状态值:**
- pending: 待确认
- confirmed: 已确认
- in_progress: 进行中
- completed: 已完成
- canceled: 已取消
- refunded: 已退款

**响应:**
```json
{
  "success": true,
  "code": 200,
  "message": "OK",
  "data": {
    "items": [
      {
        "id": 12345,
        "title": "王者荣耀陪玩",
        "status": "completed",
        "totalPriceCents": 6000,
        "quantity": 2,
        "player": {
          "id": 100,
          "name": "小明",
          "avatarUrl": "https://..."
        },
        "createdAt": "2024-01-15T10:00:00Z",
        "completedAt": "2024-01-15T12:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "pageSize": 20,
      "total": 45,
      "totalPages": 3
    }
  }
}
```

---

### 3.3 获取订单详情

**Endpoint:** `GET /user/orders/{id}`

**描述:** 获取指定订单的详细信息

**响应:**
```json
{
  "success": true,
  "code": 200,
  "message": "OK",
  "data": {
    "id": 12345,
    "title": "王者荣耀陪玩",
    "status": "completed",
    "type": "solo",
    "unitPriceCents": 3000,
    "quantity": 2,
    "totalPriceCents": 6000,
    "player": {
      "id": 100,
      "name": "小明",
      "avatarUrl": "https://...",
      "rating": 4.8
    },
    "game": {
      "id": 1,
      "name": "王者荣耀",
      "icon": "https://..."
    },
    "serviceItem": {
      "id": 10,
      "name": "上分陪玩",
      "description": "专业上分服务"
    },
    "remark": "希望能带我上王者",
    "createdAt": "2024-01-15T10:00:00Z",
    "confirmedAt": "2024-01-15T10:05:00Z",
    "startedAt": "2024-01-15T10:10:00Z",
    "completedAt": "2024-01-15T12:30:00Z",
    "payment": {
      "id": 67890,
      "method": "wallet",
      "amountCents": 6000,
      "status": "paid",
      "paidAt": "2024-01-15T10:00:00Z"
    },
    "review": {
      "id": 111,
      "rating": 5,
      "content": "很棒的陪玩体验！",
      "createdAt": "2024-01-15T13:00:00Z"
    }
  }
}
```

---

## 4. 用户模块 (`/users`, `/auth`)

### 4.1 获取当前用户信息

**Endpoint:** `GET /auth/me`

**描述:** 获取当前登录用户的详细信息

**响应:**
```json
{
  "success": true,
  "code": 200,
  "message": "OK",
  "data": {
    "id": 1001,
    "phone": "138****8888",
    "email": "user@example.com",
    "name": "游戏达人",
    "avatarUrl": "https://cdn.example.com/avatars/1001.jpg",
    "role": "user",
    "status": "active",
    "vipUnlocked": true,
    "vipLevelId": 2,
    "vipExp": 2450,
    "wallet": {
      "balanceCents": 125050,
      "frozenCents": 0
    },
    "createdAt": "2024-01-01T00:00:00Z",
    "lastLoginAt": "2024-01-15T08:00:00Z"
  }
}
```

---

### 4.2 更新用户资料

**Endpoint:** `PUT /user/profile`

**描述:** 更新当前用户的个人资料

**请求体:**
```json
{
  "name": "新昵称",
  "avatarUrl": "https://cdn.example.com/avatars/new.jpg"
}
```

**请求字段:**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 否 | 昵称，2-32 字符 |
| avatarUrl | string | 否 | 头像 URL |

**响应:**
```json
{
  "success": true,
  "code": 200,
  "message": "资料更新成功",
  "data": {
    "id": 1001,
    "name": "新昵称",
    "avatarUrl": "https://cdn.example.com/avatars/new.jpg",
    "updatedAt": "2024-01-15T10:00:00Z"
  }
}
```

---

### 4.3 修改密码

**Endpoint:** `POST /auth/change-password`

**描述:** 修改当前用户的登录密码

**请求体:**
```json
{
  "oldPassword": "OldPass123!",
  "newPassword": "NewPass456!"
}
```

**请求字段:**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| oldPassword | string | 是 | 当前密码 |
| newPassword | string | 是 | 新密码，8+ 字符，需包含大小写、数字、特殊符号 |

**响应:**
```json
{
  "success": true,
  "code": 200,
  "message": "密码修改成功"
}
```

**错误码:**
| code | message | 说明 |
|------|---------|------|
| 400 | invalid old password | 原密码错误 |
| 400 | password does not meet security requirements | 新密码不符合安全要求 |
| 401 | unauthorized | 未登录 |

---

### 4.4 请求密码重置

**Endpoint:** `POST /auth/password-reset/request`

**描述:** 请求发送密码重置邮件

**请求体:**
```json
{
  "email": "user@example.com"
}
```

**响应:**
```json
{
  "success": true,
  "code": 200,
  "message": "If the email is registered, a password reset link will be sent"
}
```

> 注意：无论邮箱是否存在，都返回成功，防止用户枚举攻击

---

### 4.5 确认密码重置

**Endpoint:** `POST /auth/password-reset/confirm`

**描述:** 使用重置令牌设置新密码

**请求体:**
```json
{
  "token": "reset-token-from-email",
  "newPassword": "NewPass456!"
}
```

**响应:**
```json
{
  "success": true,
  "code": 200,
  "message": "Password reset successfully. Please login with your new password."
}
```

---

### 4.6 通知设置（可选）

**Endpoint:** `PUT /user/settings/notifications`

**描述:** 更新用户的通知偏好设置

**请求体:**
```json
{
  "pushEnabled": true,
  "emailEnabled": false,
  "smsEnabled": true,
  "orderNotify": true,
  "promotionNotify": false
}
```

**响应:**
```json
{
  "success": true,
  "code": 200,
  "message": "通知设置已更新",
  "data": {
    "pushEnabled": true,
    "emailEnabled": false,
    "smsEnabled": true,
    "orderNotify": true,
    "promotionNotify": false,
    "updatedAt": "2024-01-15T10:00:00Z"
  }
}
```

---

## 5. 充值档位模块 (`/recharge`)

### 5.1 获取充值档位列表

**Endpoint:** `GET /user/recharge/options`

**描述:** 获取可用的充值档位列表

**响应:**
```json
{
  "success": true,
  "code": 200,
  "message": "OK",
  "data": [
    {
      "id": 1,
      "name": "小额充值",
      "amountCents": 1000,
      "bonusCents": 0,
      "totalCents": 1000,
      "isHot": false,
      "isActive": true
    },
    {
      "id": 2,
      "name": "热门充值",
      "amountCents": 5000,
      "bonusCents": 500,
      "totalCents": 5500,
      "isHot": true,
      "isActive": true
    },
    {
      "id": 3,
      "name": "大额充值",
      "amountCents": 10000,
      "bonusCents": 1500,
      "totalCents": 11500,
      "isHot": false,
      "isActive": true
    }
  ]
}
```

---

### 5.2 创建充值订单

**Endpoint:** `POST /user/recharge/orders`

**描述:** 选择档位创建充值订单

**请求体:**
```json
{
  "optionId": 2,
  "paymentChannel": "alipay",
  "paymentMethod": "app"
}
```

**响应:**
```json
{
  "success": true,
  "code": 201,
  "message": "created",
  "data": {
    "record": {
      "id": 123,
      "userId": 1001,
      "optionId": 2,
      "amountCents": 5000,
      "bonusCents": 500,
      "totalCents": 5500,
      "status": "pending",
      "paymentChannel": "alipay",
      "createdAt": "2024-01-15T10:00:00Z"
    },
    "payInfo": null
  }
}
```

---

## 附录

### A. HTTP 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 500 | 服务器内部错误 |

### B. 分页参数

所有列表接口支持以下分页参数：

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int | 1 | 页码，从 1 开始 |
| pageSize | int | 20 | 每页数量，最大 100 |

### C. 时间格式

所有时间字段使用 ISO 8601 格式：`YYYY-MM-DDTHH:mm:ssZ`

### D. 金额转换

```typescript
// 分转元
const yuan = cents / 100;

// 元转分
const cents = Math.round(yuan * 100);

// 格式化显示
const formatMoney = (cents: number) => `¥${(cents / 100).toFixed(2)}`;
```
