# API 设计规范

## 📋 概述

本文档定义了 GameLink 陪玩管理平台的 API 设计规范，包括RESTful API、WebSocket API、GraphQL API等的设计原则、命名规范、数据格式、错误处理等。

## 🎯 设计原则

### 1. RESTful 设计
- 遵循 REST 架构风格
- 使用合适的 HTTP 方法
- 资源导向的 URL 设计
- 无状态通信
- 统一的接口设计

### 2. 一致性
- 统一的 URL 命名规范
- 统一的响应格式
- 统一的错误处理
- 统一的分页机制

### 3. 可扩展性
- 版本控制机制
- 向后兼容性保证
- 灵活的查询参数
- 可选的响应字段

### 4. 安全性
- HTTPS 强制使用
- 身份认证和授权
- 输入验证和过滤
- 防止常见攻击

## 🌐 URL 设计规范

### 基础格式
```
https://api.gamelink.com/v1/{resource}/{id}/{sub-resource}/{sub-id}
```

### 命名规范
```http
# ✅ 好的URL设计
GET    /api/v1/users                    # 获取用户列表
GET    /api/v1/users/{id}               # 获取特定用户
POST   /api/v1/users                    # 创建用户
PUT    /api/v1/users/{id}               # 更新用户
DELETE /api/v1/users/{id}               # 删除用户

GET    /api/v1/users/{id}/orders        # 获取用户的订单
POST   /api/v1/users/{id}/orders        # 为用户创建订单
GET    /api/v1/orders/{id}/items        # 获取订单项目

# ❌ 避免的URL设计
GET    /api/v1/getAllUsers              # 动词在URL中
POST   /api/v1/users/{id}/deleteUser    # 重复的资源名
GET    /api/v1/users/{userId}/orders    # 不一致的命名
GET    /api/v1/user-orders              # 复合名词
```

### 复杂查询
```http
# ✅ 使用查询参数处理复杂过滤
GET /api/v1/orders?status=pending&user_id=123&game_id=456&created_after=2024-01-01

# ✅ 搜索功能
GET /api/v1/users/search?q=张三&fields=name,phone

# ✅ 排序
GET /api/v1/orders?sort=created_at:desc,price:asc

# ✅ 分页
GET /api/v1/orders?page=1&page_size=20

# ✅ 字段选择
GET /api/v1/users?fields=id,name,phone,avatar
```

## 📊 HTTP 方法使用

### 标准方法
| 方法 | 用途 | 是否幂等 | 是否安全 |
|------|------|----------|----------|
| GET  | 获取资源 | ✅ | ✅ |
| POST | 创建资源 | ❌ | ❌ |
| PUT  | 完整更新资源 | ✅ | ❌ |
| PATCH | 部分更新资源 | ❌ | ❌ |
| DELETE | 删除资源 | ✅ | ❌ |

### 使用示例
```http
# 资源集合
GET    /api/v1/users              # 获取用户列表
POST   /api/v1/users              # 创建新用户

# 单个资源
GET    /api/v1/users/{id}         # 获取特定用户
PUT    /api/v1/users/{id}         # 完整更新用户
PATCH  /api/v1/users/{id}         # 部分更新用户
DELETE /api/v1/users/{id}         # 删除用户

# 子资源
GET    /api/v1/users/{id}/orders  # 获取用户订单
POST   /api/v1/users/{id}/orders  # 为用户创建订单

# 自定义动作
POST   /api/v1/orders/{id}/cancel     # 取消订单
POST   /api/v1/orders/{id}/confirm    # 确认订单
POST   /api/v1/payments/{id}/refund   # 申请退款
```

## 📦 请求格式

### Content-Type
```http
# JSON格式
Content-Type: application/json

# 文件上传
Content-Type: multipart/form-data

# 表单提交
Content-Type: application/x-www-form-urlencoded

# GraphQL
Content-Type: application/graphql
```

### 请求头规范
```http
# 必需的请求头
Authorization: Bearer {jwt_token}
Content-Type: application/json
Accept: application/json
User-Agent: GameLink/1.0.0 (iOS)

# 可选的请求头
X-Request-ID: {unique_request_id}
X-Client-Version: 1.0.0
X-Device-ID: {device_identifier}
X-Platform: ios|android|web
```

### 分页参数
```json
{
  "page": 1,
  "page_size": 20,
  "total": 100,
  "total_pages": 5,
  "has_next": true,
  "has_prev": false
}
```

## 📤 响应格式

### 成功响应
```json
{
  "success": true,
  "code": 200,
  "message": "OK",
  "data": {
    // 响应数据
  },
  "meta": {
    "timestamp": "2024-01-01T12:00:00Z",
    "request_id": "req_123456789",
    "version": "v1"
  }
}
```

### 列表响应
```json
{
  "success": true,
  "code": 200,
  "message": "OK",
  "data": [
    {
      "id": 1,
      "name": "张三",
      "phone": "13800138000"
    },
    {
      "id": 2,
      "name": "李四",
      "phone": "13800138001"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5,
    "has_next": true,
    "has_prev": false
  },
  "meta": {
    "timestamp": "2024-01-01T12:00:00Z",
    "request_id": "req_123456789",
    "version": "v1"
  }
}
```

### 错误响应
```json
{
  "success": false,
  "code": 400,
  "message": "请求参数错误",
  "error": {
    "type": "VALIDATION_ERROR",
    "details": [
      {
        "field": "phone",
        "message": "手机号格式不正确"
      },
      {
        "field": "password",
        "message": "密码长度不能少于6位"
      }
    ]
  },
  "meta": {
    "timestamp": "2024-01-01T12:00:00Z",
    "request_id": "req_123456789",
    "version": "v1"
  }
}
```

## 🔢 数据类型规范

### 基本类型
```json
{
  "string_field": "string_value",
  "integer_field": 123,
  "float_field": 123.45,
  "boolean_field": true,
  "null_field": null,
  "timestamp_field": "2024-01-01T12:00:00Z",
  "date_field": "2024-01-01",
  "uuid_field": "550e8400-e29b-41d4-a716-446655440000"
}
```

### 枚举类型
```json
{
  "order_status": "pending",
  "user_type": "player",
  "payment_method": "wechat"
}

// 枚举值定义
// order_status: pending, confirmed, in_progress, completed, cancelled
// user_type: user, player, admin
// payment_method: wechat, alipay, bank_card
```

### 复杂对象
```json
{
  "user": {
    "id": 123,
    "name": "张三",
    "phone": "13800138000",
    "avatar": "https://cdn.gamelink.com/avatars/123.jpg",
    "profile": {
      "age": 25,
      "gender": "male",
      "location": "北京"
    },
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

## ❌ 错误处理规范

### HTTP 状态码
| 状态码 | 含义 | 使用场景 |
|--------|------|----------|
| 200 | OK | 请求成功 |
| 201 | Created | 资源创建成功 |
| 204 | No Content | 删除成功，无内容返回 |
| 400 | Bad Request | 请求参数错误 |
| 401 | Unauthorized | 未认证 |
| 403 | Forbidden | 无权限 |
| 404 | Not Found | 资源不存在 |
| 409 | Conflict | 资源冲突 |
| 422 | Unprocessable Entity | 请求格式正确但语义错误 |
| 429 | Too Many Requests | 请求频率限制 |
| 500 | Internal Server Error | 服务器内部错误 |
| 502 | Bad Gateway | 网关错误 |
| 503 | Service Unavailable | 服务不可用 |

### 业务错误码
```json
{
  "success": false,
  "code": 400,
  "message": "业务错误",
  "error": {
    "type": "BUSINESS_ERROR",
    "code": "USER_ALREADY_EXISTS",
    "message": "用户已存在"
  }
}

// 常见业务错误码
// USER_NOT_FOUND: 用户不存在
// USER_ALREADY_EXISTS: 用户已存在
// INVALID_PASSWORD: 密码错误
// ORDER_NOT_FOUND: 订单不存在
// ORDER_STATUS_INVALID: 订单状态无效
// PAYMENT_FAILED: 支付失败
// INSUFFICIENT_BALANCE: 余额不足
```

### 验证错误
```json
{
  "success": false,
  "code": 422,
  "message": "数据验证失败",
  "error": {
    "type": "VALIDATION_ERROR",
    "details": [
      {
        "field": "phone",
        "code": "INVALID_FORMAT",
        "message": "手机号格式不正确",
        "value": "123456"
      },
      {
        "field": "password",
        "code": "TOO_SHORT",
        "message": "密码长度不能少于6位",
        "value": "123"
      }
    ]
  }
}
```

## 🔐 认证和授权

### JWT Token 格式
```json
// Header
{
  "alg": "HS256",
  "typ": "JWT"
}

// Payload
{
  "sub": "123",
  "user_id": 123,
  "user_type": "user",
  "exp": 1640995200,
  "iat": 1640908800,
  "iss": "gamelink"
}
```

### 认证流程
```http
# 1. 用户登录
POST /api/v1/auth/login
Content-Type: application/json

{
  "phone": "13800138000",
  "password": "password123"
}

# 响应
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 3600,
    "user": {
      "id": 123,
      "name": "张三",
      "phone": "13800138000"
    }
  }
}

# 2. 使用Token访问API
GET /api/v1/users/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

### 权限控制
```http
# 角色权限定义
// user: 查看自己的信息，创建订单
// player: 查看订单信息，接单，管理个人信息
// admin: 管理所有用户和订单，查看统计数据

# 权限检查示例
GET /api/v1/admin/users
Authorization: Bearer {admin_token}
X-User-Role: admin
```

## 🔍 查询和过滤

### 通用查询参数
```http
# 分页
?cursor=next_page_token&limit=20
?page=1&page_size=20

# 排序
?sort=created_at:desc,price:asc
?order_by=-created_at,price

# 过滤
?status=pending&user_id=123&game_id=456
?created_after=2024-01-01&created_before=2024-01-31

# 字段选择
?fields=id,name,phone,avatar
?exclude=password,salt

# 搜索
?q=张三&search_fields=name,phone
```

### 高级查询
```http
# 范围查询
?price[gte]=100&price[lte]=500
?created_at[between]=2024-01-01,2024-01-31

# 数组查询
?status[]=pending&status[]=confirmed
?game_id[]=1&game_id[]=2&game_id[]=3

# 地理位置查询
?location[near]=39.9042,116.4074&location[radius]=1000
```

## 📱 版本控制

### URL版本控制
```http
# 主版本
/api/v1/users
/api/v2/users

# 子版本（可选）
/api/v1.1/users
```

### 请求头版本控制
```http
Accept: application/vnd.gamelink.v1+json
Accept: application/vnd.gamelink.v2+json

API-Version: v1
API-Version: v2
```

### 版本兼容性
- 向后兼容：新版本支持旧版本客户端
- 废弃通知：提前通知API废弃计划
- 迁移指南：提供版本迁移文档

## 🔄 WebSocket API

### 连接格式
```
wss://api.gamelink.com/v1/ws?token={jwt_token}
```

### 消息格式
```json
{
  "type": "message_type",
  "id": "message_id",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    // 消息内容
  }
}
```

### 消息类型
```json
// 客户端发送
{
  "type": "subscribe",
  "data": {
    "channels": ["order_updates", "notifications"]
  }
}

{
  "type": "unsubscribe",
  "data": {
    "channels": ["order_updates"]
  }
}

// 服务端推送
{
  "type": "order_update",
  "data": {
    "order_id": 123,
    "status": "confirmed",
    "message": "订单已确认"
  }
}

{
  "type": "notification",
  "data": {
    "id": 456,
    "title": "新订单通知",
    "content": "您有一个新的订单",
    "type": "new_order"
  }
}
```

## 🚀 性能优化

### 缓存策略
```http
# 缓存控制头
Cache-Control: public, max-age=3600
ETag: "abc123"
Last-Modified: Wed, 01 Jan 2024 12:00:00 GMT

# 条件请求
If-None-Match: "abc123"
If-Modified-Since: Wed, 01 Jan 2024 12:00:00 GMT
```

### 压缩
```http
# 启用压缩
Accept-Encoding: gzip, deflate, br
Content-Encoding: gzip
```

### 限流
```http
# 限流信息
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640995200

# 超限响应
{
  "success": false,
  "code": 429,
  "message": "请求过于频繁",
  "error": {
    "type": "RATE_LIMIT_EXCEEDED",
    "retry_after": 60
  }
}
```

## 📊 API文档

### OpenAPI规范
```yaml
openapi: 3.0.0
info:
  title: GameLink API
  version: 1.0.0
  description: 陪玩管理平台API
servers:
  - url: https://api.gamelink.com/v1
    description: 生产环境
  - url: https://staging-api.gamelink.com/v1
    description: 测试环境

paths:
  /users:
    get:
      summary: 获取用户列表
      parameters:
        - name: page
          in: query
          schema:
            type: integer
            default: 1
        - name: page_size
          in: query
          schema:
            type: integer
            default: 20
            maximum: 100
      responses:
        '200':
          description: 成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/UserListResponse'
```

### 文档生成
- 使用Swagger/OpenAPI规范
- 自动生成文档
- 提供在线测试界面
- 代码示例和SDK

## 🧪 API测试

### 测试工具
- Postman集合
- 自动化测试脚本
- 性能测试工具
- 安全测试工具

### 测试覆盖
- 功能测试
- 边界测试
- 错误处理测试
- 性能测试
- 安全测试

---

遵循这些API设计规范将帮助我们创建一致、可靠、易用的API接口。如有疑问，请与团队讨论并持续改进规范。