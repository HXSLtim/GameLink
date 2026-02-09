# API 速率限制中间件文档

## 概述

本项目实现了基于 Redis 的 API 速率限制中间件，使用 `github.com/ulule/limiter/v3` 库实现分布式限流功能。

## 功能特性

- ✅ **Redis 后端存储**：使用 Redis 存储计数器，支持分布式部署
- ✅ **标准响应头**：返回 `X-RateLimit-*` 响应头，符合 RFC 6585 规范
- ✅ **429 状态码**：超限返回 HTTP 429 (Too Many Requests) 状态码
- ✅ **灵活配置**：支持按路由配置不同的限流规则
- ✅ **用户/IP 限流**：支持基于用户 ID 或 IP 地址的限流
- ✅ **故障降级**：Redis 连接失败时自动降级（fail-open），不影响业务

## 限流配置

### 默认配置

| 接口类型 | 限制 | 时间窗口 | 说明 |
|---------|------|---------|------|
| 登录接口 | 5 次 | 15 分钟 | 防止暴力破解 |
| 短信接口 | 3 次 | 1 小时 | 防止短信轰炸 |
| 通用 API | 60 次 | 1 分钟 | 防止 API 滥用 |
| 上传接口 | 10 次 | 1 小时 | 防止存储滥用 |

### 配置结构

```go
type RedisRateLimiterConfig struct {
    LoginRequests    int           // 登录接口请求数限制
    LoginWindow      time.Duration // 登录接口时间窗口
    SMSRequests      int           // 短信接口请求数限制
    SMSWindow        time.Duration // 短信接口时间窗口
    APIRequests      int           // 通用 API 请求数限制
    APIWindow        time.Duration // 通用 API 时间窗口
    UploadRequests   int           // 上传接口请求数限制
    UploadWindow     time.Duration // 上传接口时间窗口
    Enabled          bool          // 是否启用限流
}
```

## 使用方法

### 1. 创建速率限制器

```go
import (
    "gamelink/internal/handler/middleware"
    "gamelink/pkg/cache"
)

// 在路由初始化时创建限流器
config := middleware.DefaultRedisRateLimiterConfig()
redisRateLimiter, err := middleware.NewRedisRateLimiter(cacheClient, config)
if err != nil {
    log.Fatal("Failed to create rate limiter:", err)
}
```

### 2. 应用到路由

#### 登录接口限流

```go
import (
    "gamelink/internal/handler"
    "gamelink/internal/handler/middleware"
)

// 方式 1: 在注册路由时应用
func RegisterAuthRoutes(router gin.IRouter, svc *authservice.AuthService, rateLimitMiddleware ...gin.HandlerFunc) {
    auth := router.Group("/auth")

    // 应用登录限流
    if len(rateLimitMiddleware) > 0 && rateLimitMiddleware[0] != nil {
        auth.POST("/login", rateLimitMiddleware[0], func(c *gin.Context) {
            loginHandler(c, svc)
        })
    } else {
        auth.POST("/login", func(c *gin.Context) {
            loginHandler(c, svc)
        })
    }

    // 其他路由...
}

// 在 router.go 中调用
handler.RegisterAuthRoutes(api, r.authSvc, redisRateLimiter.LoginRateLimit())
```

#### 短信接口限流

```go
// 短信发送接口
api.POST("/auth/sms", redisRateLimiter.SMSRateLimit(), func(c *gin.Context) {
    // 处理短信发送逻辑
})
```

#### 上传接口限流

```go
// 文件上传接口
api.POST("/upload", redisRateLimiter.UploadRateLimit(), func(c *gin.Context) {
    // 处理文件上传逻辑
})
```

#### 通用 API 限流

```go
// 通用 API 接口
api.GET("/api/v1/*path", redisRateLimiter.APIRateLimit(), func(c *gin.Context) {
    // 处理 API 请求
})
```

#### 自定义限流规则

```go
// 创建自定义限流器: 100 次/小时
customLimiter := redisRateLimiter.CustomRateLimit(100, time.Hour, "custom_action")

api.POST("/custom/action", customLimiter, func(c *gin.Context) {
    // 处理自定义操作
})
```

## 响应格式

### 成功响应

当请求在限流范围内时，返回正常的响应，并包含以下响应头：

```
HTTP/1.1 200 OK
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 59
X-RateLimit-Reset: 1672531200
```

### 超限响应

当请求超过限流阈值时，返回 429 状态码：

```
HTTP/1.1 429 Too Many Requests
Content-Type: application/json
Retry-After: 60
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1672531200

{
  "success": false,
  "code": 429,
  "message": "Rate limit exceeded: 60 requests per 1m0s",
  "requestId": "req_1234567890"
}
```

## 响应头说明

| 响应头 | 类型 | 说明 |
|-------|------|------|
| `X-RateLimit-Limit` | int | 请求限制总数 |
| `X-RateLimit-Remaining` | int | 剩余可用请求数 |
| `X-RateLimit-Reset` | int64 | 限流窗口重置时间（Unix 时间戳） |
| `Retry-After` | int | 建议重试等待时间（秒） |

## 前端集成

### 处理 429 响应

```typescript
// 使用 Axios 拦截器处理 429 响应
import axios from 'axios';

const api = axios.create({
  baseURL: '/api/v1'
});

// 响应拦截器
api.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 429) {
      const retryAfter = error.response.headers['retry-after'];
      const resetTime = error.response.headers['x-ratelimit-reset'];

      // 显示友好的限流提示
      const message = error.response.data.message ||
        `请求过于频繁，请在 ${retryAfter} 秒后重试`;

      // 使用通知组件显示
      notification.warning({
        message: '请求受限',
        description: message,
        duration: retryAfter ? parseInt(retryAfter) : 5
      });

      // 或者自动重试
      if (retryAfter) {
        return new Promise((resolve) => {
          setTimeout(() => {
            resolve(api.request(error.config));
          }, parseInt(retryAfter) * 1000);
        });
      }
    }
    return Promise.reject(error);
  }
);

export default api;
```

### React 组件示例

```typescript
import { message } from 'antd';
import api from '@/api';

export const LoginForm = () => {
  const [loading, setLoading] = useState(false);

  const handleLogin = async (values: LoginValues) => {
    setLoading(true);
    try {
      const response = await api.post('/auth/login', values);
      message.success('登录成功');
      // 处理登录成功逻辑
    } catch (error: any) {
      if (error.response?.status === 429) {
        const retryAfter = error.response.headers['retry-after'];
        message.error(`登录尝试过于频繁，请在 ${retryAfter} 秒后重试`);
      } else {
        message.error('登录失败，请检查用户名和密码');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <Form onFinish={handleLogin}>
      {/* 登录表单 */}
    </Form>
  );
};
```

## 客户端 IP 获取

限流器会按以下优先级获取客户端真实 IP：

1. **X-Forwarded-For** 头（取第一个 IP）
2. **X-Real-IP** 头
3. **RemoteAddr**（直连 IP）

### Nginx 配置示例

```nginx
location /api/ {
    proxy_pass http://backend;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

## 测试

### 运行测试

```bash
cd api
go test -v ./internal/handler/middleware -run TestRedisRateLimiter
```

### 测试覆盖

测试套件包含以下测试用例：

- ✅ 限流器创建测试
- ✅ 登录接口限流测试（5次/15分钟）
- ✅ 短信接口限流测试（3次/小时）
- ✅ 通用 API 限流测试（60次/分钟）
- ✅ 上传接口限流测试（10次/小时）
- ✅ 限流禁用测试
- ✅ 自定义限流规则测试
- ✅ 用户级限流测试
- ✅ 客户端 IP 提取测试
- ✅ 响应头验证测试
- ✅ 默认配置测试

## 部署注意事项

### 1. Redis 配置

确保 Redis 实例可用且配置正确：

```yaml
# configs/config.production.yaml
cache:
  type: "redis"
  redis:
    addr: "localhost:6379"
    password: "your-redis-password"
    db: 0
```

### 2. 环境变量

生产环境建议设置以下环境变量：

```bash
# Redis 配置
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=your-password
REDIS_DB=0

# 限流配置（可选，覆盖默认值）
RATE_LIMIT_ENABLED=true
```

### 3. 监控

建议监控以下指标：

- 限流触发频率（429 响应数）
- Redis 连接状态
- 平均剩余配额
- 限流器响应时间

## 性能考虑

### Redis 连接池

限流器复用现有的 Redis 连接池（来自 `cache.Cache`），无需额外配置。

### 内存使用

每个限流 key 在 Redis 中占用约 100-200 字节，根据用户量和限流规则计算：

```
内存估算 = (活跃用户数 + 活跃IP数) × 限流规则数 × 150 字节
```

例如：
- 10,000 活跃用户
- 4 种限流规则
- 预计内存：10,000 × 4 × 150 ≈ 6 MB

### TTL 自动清理

Redis key 自动过期，无需手动清理：

- 登录限流：15 分钟
- 短信限流：1 小时
- API 限流：1 分钟
- 上传限流：1 小时

## 故障处理

### Redis 连接失败

当 Redis 连接失败时，限流器采用 **fail-open** 策略：

1. 记录错误日志
2. **不阻止请求**（允许通过）
3. 返回正常的业务响应

这确保了 Redis 故障不会影响核心业务功能。

### 日志示例

```
[ERROR] Redis rate limiter error: connection refused
[WARN] Rate limiting disabled for this request due to Redis error
```

## 常见问题

### Q1: 如何调整限流阈值？

修改 `DefaultRedisRateLimiterConfig()` 或在创建时传入自定义配置：

```go
config := middleware.RedisRateLimiterConfig{
    Enabled:       true,
    LoginRequests: 10,  // 改为 10 次
    LoginWindow:   30 * time.Minute,  // 改为 30 分钟
    // ...
}
limiter, _ := middleware.NewRedisRateLimiter(cacheClient, config)
```

### Q2: 如何为特定路由添加限流？

使用 `CustomRateLimit` 方法：

```go
api.POST("/admin/action", redisRateLimiter.CustomRateLimit(
    100,            // 100 次
    time.Hour,      // 每小时
    "admin_action", // 前缀
), handler)
```

### Q3: 限流 key 的格式是什么？

```
ratelimit:{prefix}:{type}:{identifier}

示例：
- ratelimit:login:ip:192.168.1.100
- ratelimit:sms:user:12345
- ratelimit:api:ip:10.0.0.1
```

### Q4: 如何临时禁用限流？

方式 1：配置中禁用
```go
config.Enabled = false
```

方式 2：环境变量
```bash
RATE_LIMIT_ENABLED=false
```

### Q5: 不同环境的限流配置如何管理？

建议使用不同的配置文件：

```yaml
# configs/config.development.yaml
rate_limit:
  enabled: false
  login_requests: 1000

# configs/config.production.yaml
rate_limit:
  enabled: true
  login_requests: 5
```

## 相关文件

- **实现**: `api/internal/handler/middleware/redisRateLimiter.go`
- **测试**: `api/internal/handler/middleware/redisRateLimiter_test.go`
- **Cache**: `api/pkg/cache/redis.go`

## 参考资源

- [RFC 6585: HTTP 429 Status Code](https://tools.ietf.org/html/rfc6585)
- [ulule/limiter 文档](https://github.com/ulule/limiter)
- [Redis 速率限制最佳实践](https://redis.io/docs/manual/patterns/rate-limit/)
