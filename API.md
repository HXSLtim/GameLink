# 📡 GameLink API 文档

本文档提供 GameLink 平台的完整 RESTful API 接口说明。

---

## 📋 目录

- [API 概览](#api-概览)
- [认证与授权](#认证与授权)
- [通用规范](#通用规范)
- [用户管理](#用户管理)
- [陪玩师管理](#陪玩师管理)
- [游戏管理](#游戏管理)
- [订单管理](#订单管理)
- [支付管理](#支付管理)
- [评价管理](#评价管理)
- [聊天通讯](#聊天通讯)
- [通知管理](#通知管理)
- [文件上传](#文件上传)
- [WebSocket 接口](#websocket-接口)
- [错误代码](#错误代码)

---

## 🎯 API 概览

### 基础信息
- **Base URL**: `https://api.gamelink.com/api/v1`
- **协议**: HTTPS
- **数据格式**: JSON
- **字符编码**: UTF-8

### 接口特性
- ✅ RESTful 设计风格
- ✅ JWT 认证机制
- ✅ RBAC 权限控制
- ✅ 请求限流保护
- ✅ 参数验证
- ✅ 错误统一格式
- ✅ API 版本控制

### 在线文档
- 📚 [Swagger UI](https://api.gamelink.com/swagger/index.html)
- 📖 [API 文档](https://api.gamelink.com/docs)

---

## 🔐 认证与授权

### JWT 认证
所有需要认证的接口都需要在请求头中携带 JWT Token：

```http
Authorization: Bearer <your-jwt-token>
```

### 获取 Token
```http
POST /auth/login
Content-Type: application/json

{
  "username": "user@example.com",
  "password": "password123"
}
```

**响应示例:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400,
    "user": {
      "id": 1,
      "username": "user@example.com",
      "roles": ["user"]
    }
  }
}
```

### 刷新 Token
```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "your-refresh-token"
}
```

### 权限角色
- **user**: 普通用户 - 可下单、支付、评价
- **player**: 陪玩师 - 可接单、管理服务、查看收益
- **admin**: 管理员 - 全部权限

---

## 📜 通用规范

### 请求格式
```http
GET /api/v1/resource?page=1&size=20&sort=created_at:desc
Content-Type: application/json
Authorization: Bearer <token>
```

### 响应格式
```json
{
  "success": true,
  "data": {},
  "message": "操作成功",
  "trace_id": "abc123def456",
  "timestamp": "2025-11-13T10:30:00Z"
}
```

### 错误响应格式
```json
{
  "success": false,
  "error": {
    "code": "INVALID_PARAMETER",
    "message": "参数验证失败",
    "details": {
      "field": "email",
      "reason": "邮箱格式不正确"
    }
  },
  "trace_id": "abc123def456",
  "timestamp": "2025-11-13T10:30:00Z"
}
```

### 分页格式
```json
{
  "success": true,
  "data": {
    "items": [],
    "pagination": {
      "page": 1,
      "size": 20,
      "total": 100,
      "pages": 5
    }
  }
}
```

### 状态码说明
| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 204 | 删除成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 422 | 参数验证失败 |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |

---

## 👤 用户管理

### 用户注册
```http
POST /auth/register
```

**请求参数:**
```json
{
  "username": "testuser",
  "email": "user@example.com",
  "password": "password123",
  "phone": "13800138000",
  "role": "user"
}
```

**响应示例:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "user@example.com",
    "phone": "13800138000",
    "role": "user",
    "status": "active",
    "created_at": "2025-11-13T10:30:00Z"
  }
}
```

### 用户登录
```http
POST /auth/login
```

**请求参数:**
```json
{
  "username": "user@example.com",
  "password": "password123"
}
```

### 获取当前用户信息
```http
GET /auth/me
Authorization: Bearer <token>
```

### 更新用户资料
```http
PUT /user/profile
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "nickname": "游戏达人",
  "avatar": "https://example.com/avatar.jpg",
  "bio": "资深游戏玩家",
  "birthday": "1990-01-01",
  "gender": "male"
}
```

### 修改密码
```http
PUT /user/password
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "old_password": "oldpassword123",
  "new_password": "newpassword123"
}
```

### 重置密码
```http
POST /auth/password/reset
```

**请求参数:**
```json
{
  "email": "user@example.com"
}
```

---

## 🎮 陪玩师管理

### 申请成为陪玩师
```http
POST /player/apply
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "real_name": "张三",
  "id_card": "110101199001011234",
  "phone": "13800138000",
  "qq": "123456789",
  "wechat": "wx123456",
  "experience": "3年游戏经验",
  "introduction": "专业陪玩，技术过硬",
  "games": [
    {
      "game_id": 1,
      "level": "王者",
      "price_per_hour": 5000,
      "tags": ["技术", "娱乐"]
    }
  ]
}
```

### 获取陪玩师列表
```http
GET /players?page=1&size=20&game_id=1&level=王者&price_min=1000&price_max=10000
```

**查询参数:**
- `page`: 页码 (默认: 1)
- `size`: 每页数量 (默认: 20)
- `game_id`: 游戏ID
- `level`: 段位
- `price_min`: 最低价格(分)
- `price_max`: 最高价格(分)
- `online_only`: 仅在线 (true/false)

### 获取陪玩师详情
```http
GET /players/{player_id}
```

### 更新陪玩师资料
```http
PUT /player/profile
Authorization: Bearer <token>
```

### 更新陪玩师状态
```http
PUT /player/status
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "status": "online",  // online, offline, busy
  "location": "北京市",
  "available_games": [1, 2, 3]
}
```

### 陪玩师收益统计
```http
GET /player/earnings?start_date=2025-11-01&end_date=2025-11-30
Authorization: Bearer <token>
```

---

## 🎯 游戏管理

### 获取游戏列表
```http
GET /games?status=active
```

**响应示例:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "王者荣耀",
      "icon": "https://example.com/game1.jpg",
      "category": "MOBA",
      "status": "active",
      "sort_order": 1,
      "created_at": "2025-11-13T10:30:00Z"
    }
  ]
}
```

### 获取游戏详情
```http
GET /games/{game_id}
```

---

## 📋 订单管理

### 创建订单
```http
POST /user/orders
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "game_id": 1,
  "service_type": "accompany",  // accompany, teach, team
  "player_id": 1,
  "duration": 60,  // 分钟
  "scheduled_start": "2025-11-13T20:00:00Z",
  "requirements": "需要一个技术好的陪玩师",
  "gift_message": "送给朋友的礼物"
}
```

### 获取订单列表
```http
GET /user/orders?status=pending&page=1&size=20
Authorization: Bearer <token>
```

**查询参数:**
- `status`: 订单状态 (pending, confirmed, in_progress, completed, cancelled)
- `game_id`: 游戏ID
- `start_date`: 开始日期
- `end_date`: 结束日期

### 获取订单详情
```http
GET /user/orders/{order_id}
Authorization: Bearer <token>
```

### 取消订单
```http
PUT /user/orders/{order_id}/cancel
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "reason": "临时有事，无法参加"
}
```

### 确认订单完成
```http
PUT /user/orders/{order_id}/complete
Authorization: Bearer <token>
```

### 评价订单
```http
POST /user/orders/{order_id}/review
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "rating": 5,
  "comment": "陪玩师技术很好，服务态度也很棒",
  "tags": ["技术好", "耐心", "准时"]
}
```

### 陪玩师订单操作
```http
PUT /player/orders/{order_id}/accept
Authorization: Bearer <token>
```

**可选操作:**
- `accept`: 接受订单
- `reject`: 拒绝订单
- `start`: 开始服务
- `complete`: 完成服务

---

## 💳 支付管理

### 创建支付
```http
POST /payments
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "order_id": 1,
  "payment_method": "alipay",  // alipay, wechat, balance
  "amount": 5000,  // 分
  "return_url": "https://example.com/payment/return",
  "notify_url": "https://example.com/payment/notify"
}
```

### 支付回调
```http
POST /payments/{payment_id}/callback
```

### 获取支付状态
```http
GET /payments/{payment_id}
Authorization: Bearer <token>
```

### 申请退款
```http
POST /payments/{payment_id}/refund
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "reason": "服务未达到预期",
  "amount": 2500
}
```

---

## ⭐ 评价管理

### 获取评价列表
```http
GET /reviews?target_type=player&target_id=1&page=1&size=20
```

**查询参数:**
- `target_type`: 评价目标类型 (player, user)
- `target_id`: 目标ID
- `rating`: 评分筛选
- `has_comment`: 是否有评论

### 添加评价回复
```http
POST /reviews/{review_id}/reply
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "content": "感谢您的好评，会继续提供优质服务"
}
```

---

## 💬 聊天通讯

### 获取聊天室列表
```http
GET /chat/rooms
Authorization: Bearer <token>
```

### 获取聊天消息
```http
GET /chat/rooms/{room_id}/messages?page=1&size=50
Authorization: Bearer <token>
```

### 发送消息
```http
POST /chat/rooms/{room_id}/messages
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "content": "你好，什么时候开始？",
  "message_type": "text"  // text, image, file
}
```

### 上传聊天图片
```http
POST /chat/upload
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: <image>
```

---

## 🔔 通知管理

### 获取通知列表
```http
GET /notifications?page=1&size=20&unread_only=true
Authorization: Bearer <token>
```

### 标记通知已读
```http
PUT /notifications/{notification_id}/read
Authorization: Bearer <token>
```

### 批量标记已读
```http
POST /notifications/read
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "notification_ids": [1, 2, 3]
}
```

### 删除通知
```http
DELETE /notifications/{notification_id}
Authorization: Bearer <token>
```

---

## 📁 文件上传

### 上传头像
```http
POST /upload/avatar
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: <image>
```

### 上传陪玩师证书
```http
POST /upload/certificate
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: <image>
type: "id_card"  // id_card, skill_certificate
```

### 获取上传信息
```http
GET /upload/info/{file_id}
Authorization: Bearer <token>
```

---

## 🔌 WebSocket 接口

### 连接WebSocket
```
ws://api.gamelink.com/ws/chat?token=<jwt-token>
```

### 消息格式
```json
{
  "type": "message",
  "data": {
    "room_id": 1,
    "content": "Hello",
    "message_type": "text"
  },
  "timestamp": "2025-11-13T10:30:00Z"
}
```

### 消息类型
- `message`: 聊天消息
- `notification`: 通知消息
- `order_update`: 订单状态更新
- `system`: 系统消息
- `heartbeat`: 心跳消息

---

## 🚫 错误代码

### 通用错误
| 错误代码 | 说明 |
|----------|------|
| SUCCESS | 操作成功 |
| INVALID_PARAMETER | 参数无效 |
| MISSING_PARAMETER | 缺少必需参数 |
| UNAUTHORIZED | 未授权 |
| FORBIDDEN | 禁止访问 |
| NOT_FOUND | 资源不存在 |
| CONFLICT | 资源冲突 |
| RATE_LIMITED | 请求频率限制 |
| INTERNAL_ERROR | 服务器内部错误 |

### 认证错误
| 错误代码 | 说明 |
|----------|------|
| INVALID_CREDENTIALS | 用户名或密码错误 |
| TOKEN_EXPIRED | Token 已过期 |
| TOKEN_INVALID | Token 无效 |
| ACCOUNT_LOCKED | 账户被锁定 |
| ACCOUNT_INACTIVE | 账户未激活 |

### 业务错误
| 错误代码 | 说明 |
|----------|------|
| PLAYER_NOT_FOUND | 陪玩师不存在 |
| ORDER_NOT_FOUND | 订单不存在 |
| ORDER_STATUS_INVALID | 订单状态无效 |
| PAYMENT_FAILED | 支付失败 |
| INSUFFICIENT_BALANCE | 余额不足 |
| FILE_TOO_LARGE | 文件过大 |
| UNSUPPORTED_FILE_TYPE | 不支持的文件类型 |

---

## 🔧 SDK 和工具

### JavaScript SDK
```bash
npm install gamelink-sdk
```

```javascript
import { GameLinkAPI } from 'gamelink-sdk';

const api = new GameLinkAPI({
  baseURL: 'https://api.gamelink.com/api/v1',
  token: 'your-jwt-token'
});

// 获取用户信息
const user = await api.auth.me();

// 创建订单
const order = await api.orders.create({
  game_id: 1,
  service_type: 'accompany',
  duration: 60
});
```

### Python SDK
```bash
pip install gamelink-sdk
```

```python
from gamelink_sdk import GameLinkAPI

api = GameLinkAPI(
    base_url='https://api.gamelink.com/api/v1',
    token='your-jwt-token'
)

# 获取用户信息
user = api.auth.me()

# 创建订单
order = api.orders.create({
    'game_id': 1,
    'service_type': 'accompany',
    'duration': 60
})
```

---

## 📊 限制说明

### 请求频率限制
- **普通接口**: 100 次/分钟
- **上传接口**: 10 次/分钟
- **登录接口**: 5 次/分钟
- **注册接口**: 3 次/分钟

### 文件上传限制
- **头像大小**: 最大 2MB
- **证书大小**: 最大 5MB
- **聊天图片**: 最大 10MB
- **支持格式**: JPG, PNG, GIF

### 分页限制
- **最大页大小**: 100
- **默认页大小**: 20

---

## 🔒 安全说明

### 数据加密
- 所有 HTTPS 通信使用 TLS 1.2+
- 敏感数据传输加密
- 密码存储使用 bcrypt 哈希

### 安全防护
- SQL 注入防护
- XSS 攻击防护
- CSRF 攻击防护
- 文件上传安全检查

---

## 📞 技术支持

### 联系方式
- **API 技术支持**: api-support@gamelink.com
- **开发者社区**: https://community.gamelink.com
- **问题反馈**: https://github.com/your-org/GameLink/issues

### 文档更新
- **API 版本**: v1.0
- **最后更新**: 2025-11-13
- **更新频率**: 随版本更新

---

*本文档持续更新中，最新版本请访问: https://docs.gamelink.com/api*