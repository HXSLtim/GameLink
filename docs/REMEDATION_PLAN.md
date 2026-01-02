# GameLink 项目补全行动计划

> 生成时间: 2026-01-03
> 基于审查报告: 2026-01-03 代码审查

---

## 执行摘要

| 阶段 | 时间 | 任务数 | 优先级 | 目标 |
|------|------|--------|--------|------|
| **Phase 1** | 1-2周 | 5个 | P0 紧急 | 修复高危安全漏洞和关键测试问题 |
| **Phase 2** | 2-4周 | 6个 | P1 重要 | 完善安全机制和代码质量 |
| **Phase 3** | 1-3月 | 6个 | P2 优化 | 架构改进和长期优化 |

**总计**: 17 个任务，预计 1-3 个月完成

---

## Phase 1: 紧急修复 (P0) - 1-2周

### 🔴 P0-SEC-001: 启用 CSRF 保护

**问题描述**: CSRF 中间件已实现但未启用，所有状态改变操作易受攻击。

**CVSS 评分**: 6.5 (Medium)
**影响范围**: 所有 POST/PUT/DELETE 请求

#### 修复步骤

1. **后端修改**: `api/internal/router/router.go`

```go
// 在需要认证的路由组中启用 CSRF
func RegisterRoutes(r *gin.Engine, services *services.All) {
    // 公开路由（不需要 CSRF）
    public := r.Group("/api/v1")
    {
        public.POST("/auth/login", authHandler.Login)
        public.GET("/healthz", healthHandler.Check)
    }

    // 认证路由（启用 CSRF）
    auth := r.Group("/api/v1")
    auth.Use(middleware.JWTAuth())
    auth.Use(middleware.CSRF(middleware.CSRFConfig{
        TokenLength:    32,
        CookieName:     "_csrf",
        HeaderName:     "X-CSRF-Token",
        CookieSecure:   config.IsProduction(), // 生产环境必须 HTTPS
        CookieHTTPOnly: false,                  // 允许 JS 读取
        CookieSameSite: http.SameSiteStrictMode,
    }))
    {
        auth.POST("/orders", orderHandler.CreateOrder)
        auth.PUT("/orders/:id", orderHandler.UpdateOrder)
        // ... 其他需要 CSRF 保护的接口
    }
}
```

2. **前端修改**: `admin/src/api/client.ts`

```typescript
import Cookies from 'js-cookie';

// 获取 CSRF token
const getCsrfToken = (): string => {
    return Cookies.get('_csrf') || '';
};

// 请求拦截器添加 CSRF token
apiClient.interceptors.request.use((config) => {
    const csrfToken = getCsrfToken();
    if (csrfToken && ['post', 'put', 'delete', 'patch'].includes(config.method?.toLowerCase() || '')) {
        config.headers['X-CSRF-Token'] = csrfToken;
    }
    return config;
});

// 响应拦截器刷新 CSRF token
apiClient.interceptors.response.use(
    (response) => {
        const newCsrfToken = response.headers['x-csrf-token'];
        if (newCsrfToken) {
            Cookies.set('_csrf', newCsrfToken);
        }
        return response;
    }
);
```

3. **测试**: `api/internal/handler/middleware/csrf_test.go`

```go
func TestCSRFMiddleware(t *testing.T) {
    // 测试无 token 拒绝
    // 测试有效 token 通过
    // 测试 token 刷新
}
```

**验收标准**:
- [ ] 所有 POST/PUT/DELETE 请求需要 CSRF token
- [ ] 前端自动获取和发送 CSRF token
- [ ] 测试覆盖 100%

---

### 🔴 P0-SEC-002: 移除硬编码密钥

**问题描述**: 开发配置文件包含硬编码密钥，存在泄露风险。

**CVSS 评分**: 7.5 (High)
**影响**: 密钥泄露可导致所有加密通信被解密

#### 修复步骤

1. **创建 `.env.example`**: `api/configs/.env.example`

```bash
# 加密配置
CRYPTO_SECRET_KEY=your-32-character-secret-key-here
CRYPTO_IV=your-16-character-iv-here

# JWT 配置
JWT_SECRET_KEY=your-jwt-secret-key-min-32-chars

# 数据库
DB_PASSWORD=your-database-password

# Redis
REDIS_PASSWORD=your-redis-password
```

2. **更新 `.gitignore`**

```gitignore
# 敏感配置文件
configs/config.*.yaml
.env
.env.local
.env.production
*.key
*.pem
```

3. **修改配置加载**: `api/pkg/config/config.go`

```go
func Load() *AppConfig {
    cfg := &AppConfig{}

    // 从环境变量读取敏感配置
    if cryptoKey := os.Getenv("CRYPTO_SECRET_KEY"); cryptoKey != "" {
        cfg.Crypto.SecretKey = cryptoKey
    }
    if jwtKey := os.Getenv("JWT_SECRET_KEY"); jwtKey != "" {
        cfg.Auth.JWTSecret = jwtKey
    }

    // 生产环境强制验证
    if os.Getenv("APP_ENV") == "production" {
        if err := cfg.Validate(); err != nil {
            log.Fatalf("配置验证失败: %v", err)
        }
    }

    return cfg
}

func (cfg *AppConfig) Validate() error {
    if cfg.Crypto.Enabled {
        if len(cfg.Crypto.SecretKey) < 32 {
            return fmt.Errorf("CRYPTO_SECRET_KEY 必须至少 32 字符")
        }
        if len(cfg.Crypto.IV) < 16 {
            return fmt.Errorf("CRYPTO_IV 必须至少 16 字符")
        }
    }
    if len(cfg.Auth.JWTSecret) < 32 {
        return fmt.Errorf("JWT_SECRET_KEY 必须至少 32 字符")
    }
    return nil
}
```

4. **前端密钥移除**: `admin/src/utils/crypto.ts`

```typescript
// ❌ 删除硬编码密钥
// const SECRET_KEY = "H/oguKMv23lWlivgq8snNZmTzSUp6KSHZnEEo1c0Ook=";

// ✅ 前端不处理加密，仅由后端处理
export const decryptData = (data: string): string => {
    // 抛出错误，不应在前端解密
    throw new Error('前端不应处理解密，请由后端处理');
};
```

**验收标准**:
- [ ] 敏感文件在 .gitignore 中
- [ ] 环境变量正确配置
- [ ] 生产环境强制验证
- [ ] 前端不再存储加密密钥

---

### 🔴 P0-SEC-003: 添加 API 速率限制

**问题描述**: API 接口缺少速率限制，易受暴力破解和 DoS 攻击。

**CVSS 评分**: 5.0 (Medium)

#### 修复步骤

1. **安装依赖**: `api/go.mod`

```bash
go get github.com/ulule/limiter/v3
go get github.com/ulule/limiter/v3/drivers/store/redis
```

2. **创建速率限制器**: `api/internal/handler/middleware/rateLimiter.go`

```go
package middleware

import (
    "github.com/ulule/limiter/v3"
    "github.com/ulule/limiter/v3/drivers/store/redis"
    "net/http"
    "time"
)

// 速率限制配置
var RateLimits = struct {
    Login    limiter.Rate
    SMS      limiter.Rate
    API      limiter.Rate
    Upload   limiter.Rate
}{
    Login: limiter.Rate{
        Period: 15 * time.Minute,
        Limit:  5, // 15分钟内5次
    },
    SMS: limiter.Rate{
        Period: 1 * time.Hour,
        Limit:  3, // 1小时内3次
    },
    API: limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  60, // 每分钟60次
    },
    Upload: limiter.Rate{
        Period: 1 * time.Hour,
        Limit:  10, // 1小时内10次
    },
}

// RateLimitMiddleware 创建速率限制中间件
func RateLimitMiddleware(store *redis.Store, rate limiter.Rate) gin.HandlerFunc {
    instance := limiter.New(store, rate)

    return func(c *gin.Context) {
        key := c.ClientIP() // 使用 IP 作为限制键
        context, err := instance.Get(c.Request.Context(), key)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "success": false,
                "code":    500,
                "message": "速率限制检查失败",
            })
            c.Abort()
            return
        }

        // 设置速率限制响应头
        c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
        c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
        c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

        if context.Reached {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "success": false,
                "code":    429,
                "message": "请求过于频繁，请稍后再试",
                "data": gin.H{
                    "retry_after": int(context.Reset.Sub(time.Now()).Seconds()),
                },
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
```

3. **应用到路由**: `api/internal/router/router.go`

```go
func RegisterRoutes(r *gin.Engine, services *services.All) {
    // 创建 Redis store
    store, err := redis.NewStore(redis.NewClient(config.Config.Redis))
    if err != nil {
        log.Fatalf("创建 Redis store 失败: %v", err)
    }

    // 公开路由
    public := r.Group("/api/v1")
    {
        // 登录接口严格限制
        public.POST("/auth/login",
            middleware.RateLimitMiddleware(store, middleware.RateLimits.Login),
            authHandler.Login)

        // 短信验证码严格限制
        public.POST("/auth/sms",
            middleware.RateLimitMiddleware(store, middleware.RateLimits.SMS),
            authHandler.SendSMS)
    }

    // 认证路由
    auth := r.Group("/api/v1")
    auth.Use(middleware.JWTAuth())
    {
        // 通用 API 限制
        auth.GET("/orders",
            middleware.RateLimitMiddleware(store, middleware.RateLimits.API),
            orderHandler.ListOrders)

        // 上传接口严格限制
        auth.POST("/upload",
            middleware.RateLimitMiddleware(store, middleware.RateLimits.Upload),
            uploadHandler.Upload)
    }
}
```

4. **前端处理**: `admin/src/api/client.ts`

```typescript
// 响应拦截器处理 429 错误
apiClient.interceptors.response.use(
    (response) => response,
    async (error) => {
        if (error.response?.status === 429) {
            const retryAfter = error.response.data?.data?.retry_after;
            if (retryAfter) {
                message.warning(`请求过于频繁，请在 ${retryAfter} 秒后重试`);
            } else {
                message.warning('请求过于频繁，请稍后再试');
            }
        }
        return Promise.reject(error);
    }
);
```

**验收标准**:
- [ ] 登录接口限制 5次/15分钟
- [ ] 短信接口限制 3次/小时
- [ ] 通用 API 限制 60次/分钟
- [ ] 返回 429 状态码和重试时间
- [ ] 前端正确处理速率限制

---

### 🔴 P0-CODE-001: 修复 React Hooks 依赖警告

**问题描述**: 20+ 条 React Hooks 依赖警告，可能导致闭包陷阱。

**影响**: 使用过期的值，导致 UI 不正确

#### 修复步骤

1. **运行检查**: `admin/`

```bash
cd admin
npm run lint 2>&1 | grep "React Hook" > hooks-warnings.txt
```

2. **批量修复**: 常见模式

```typescript
// ❌ 错误：缺少依赖
const handleRefresh = useCallback(() => {
    loadData();
    message.success('刷新成功');
}, [loadData]);  // 缺少 message

// ✅ 修复 1：添加完整依赖
const handleRefresh = useCallback(() => {
    loadData();
    message.success('刷新成功');
}, [loadData, message]);

// ✅ 修复 2：如果确定不需要，使用 eslint-disable
const handleRefresh = useCallback(() => {
    loadData();
    message.success('刷新成功');
    // eslint-disable-next-line react-hooks/exhaustive-deps
}, [loadData]);

// ✅ 修复 3：重构为不需要依赖
const handleRefresh = useCallback(() => {
    loadData();
}, [loadData]);

// 在组件中处理
const onRefresh = () => {
    handleRefresh();
    message.success('刷新成功');
};
```

3. **修复文件清单**:

| 文件 | 警告数 | 修复方式 |
|------|--------|----------|
| `pages/admin/Commission/index.tsx` | 2 | 添加依赖 |
| `pages/admin/Order/index.tsx` | 5 | 添加依赖 |
| `pages/admin/Player/index.tsx` | 3 | 添加依赖 |
| `hooks/usePermission.ts` | 1 | 添加依赖 |
| `hooks/useWebSocket.ts` | 2 | 添加依赖 |

**验收标准**:
- [ ] `npm run lint` 无 React Hooks 警告
- [ ] 所有功能正常工作
- [ ] 回归测试通过

---

### 🔴 P0-TEST-001: 补充 Handler 层测试

**问题描述**: Handler 层覆盖率仅 2.3%，严重不足。

**当前覆盖率**:
- Handler: 2.3%
- 目标: 30%

#### 修复步骤

1. **测试框架**: `api/internal/handler/handler_test.go`

```go
package handler_test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// TestHelper 测试辅助结构
type TestHelper struct {
    router *gin.Engine
    tokens map[string]string // userID -> token
}

// SetupTest 设置测试环境
func SetupTest() *TestHelper {
    gin.SetMode(gin.TestMode)
    router := gin.New()

    // 注册中间件
    router.Use(gin.Recovery())

    // TODO: 注册路由
    // router = handler.RegisterRoutes(router, services)

    return &TestHelper{
        router: router,
        tokens: make(map[string]string),
    }
}

// MakeRequest 发送测试请求
func (th *TestHelper) MakeRequest(method, path string, body interface{}, token string) *httptest.ResponseRecorder {
    var bodyReader *bytes.Buffer
    if body != nil {
        bodyBytes, _ := json.Marshal(body)
        bodyReader = bytes.NewBuffer(bodyBytes)
    }

    req, _ := http.NewRequest(method, path, bodyReader)
    req.Header.Set("Content-Type", "application/json")

    if token != "" {
        req.Header.Set("Authorization", "Bearer "+token)
    }

    w := httptest.NewRecorder()
    th.router.ServeHTTP(w, req)
    return w
}

// AssertSuccess 断言成功响应
func (th *TestHelper) AssertSuccess(t *testing.T, w *httptest.ResponseRecorder, expectedData interface{}) {
    assert.Equal(t, http.StatusOK, w.Code)

    var resp struct {
        Success bool        `json:"success"`
        Data    interface{} `json:"data"`
    }
    err := json.Unmarshal(w.Body.Bytes(), &resp)
    require.NoError(t, err)
    assert.True(t, resp.Success)

    if expectedData != nil {
        data, _ := json.Marshal(resp.Data)
        expected, _ := json.Marshal(expectedData)
        assert.JSONEq(t, string(expected), string(data))
    }
}

// AssertError 断言错误响应
func (th *TestHelper) AssertError(t *testing.T, w *httptest.ResponseRecorder, expectedCode int, expectedMessage string) {
    assert.Equal(t, expectedCode, w.Code)

    var resp struct {
        Success bool   `json:"success"`
        Code    int    `json:"code"`
        Message string `json:"message"`
    }
    err := json.Unmarshal(w.Body.Bytes(), &resp)
    require.NoError(t, err)
    assert.False(t, resp.Success)

    if expectedMessage != "" {
        assert.Contains(t, resp.Message, expectedMessage)
    }
}
```

2. **订单 Handler 测试**: `api/internal/handler/admin/order_handler_test.go`

```go
package admin_test

import (
    "net/http"
    "testing"
    "time"

    "gamelink/internal/handler"
    "gamelink/internal/model"
)

func TestOrderHandler_ListOrders(t *testing.T) {
    th := handler.SetupTest()
    adminToken := th.CreateAdminUser()

    t.Run("成功获取订单列表", func(t *testing.T) {
        // 创建测试订单
        order := th.CreateTestOrder(model.OrderStatusPending)

        w := th.MakeRequest("GET", "/api/v1/admin/orders?page=1&page_size=10", nil, adminToken)
        th.AssertSuccess(t, w, map[string]interface{}{
            "total": 1,
            "items": []interface{}{
                map[string]interface{}{
                    "id":     float64(order.ID),
                    "status": model.OrderStatusPending,
                },
            },
        })
    })

    t.Run("分页参数验证", func(t *testing.T) {
        w := th.MakeRequest("GET", "/api/v1/admin/orders?page=invalid", nil, adminToken)
        th.AssertError(t, w, http.StatusBadRequest, "invalid page")
    })

    t.Run("状态过滤", func(t *testing.T) {
        th.CreateTestOrder(model.OrderStatusConfirmed)
        th.CreateTestOrder(model.OrderStatusCompleted)

        w := th.MakeRequest("GET", "/api/v1/admin/orders?status=confirmed", nil, adminToken)
        th.AssertSuccess(t, w, map[string]interface{}{
            "total": 1,
        })
    })

    t.Run("时间范围过滤", func(t *testing.T) {
        now := time.Now()
        yesterday := now.Add(-24 * time.Hour)

        w := th.MakeRequest("GET",
            "/api/v1/admin/orders?start_time="+yesterday.Format(time.RFC3339)+
                "&end_time="+now.Format(time.RFC3339),
            nil, adminToken)
        th.AssertSuccess(t, w, nil)
    })

    t.Run("未授权访问", func(t *testing.T) {
        w := th.MakeRequest("GET", "/api/v1/admin/orders", nil, "")
        th.AssertError(t, w, http.StatusUnauthorized, "unauthorized")
    })

    t.Run("权限不足", func(t *testing.T) {
        userToken := th.CreateRegularUser()
        w := th.MakeRequest("GET", "/api/v1/admin/orders", nil, userToken)
        th.AssertError(t, w, http.StatusForbidden, "permission denied")
    })
}

func TestOrderHandler_GetOrderDetail(t *testing.T) {
    th := handler.SetupTest()
    adminToken := th.CreateAdminUser()

    t.Run("成功获取订单详情", func(t *testing.T) {
        order := th.CreateTestOrder(model.OrderStatusPending)

        w := th.MakeRequest("GET", "/api/v1/admin/orders/"+strconv.FormatUint(order.ID, 10), nil, adminToken)
        th.AssertSuccess(t, w, map[string]interface{}{
            "id":     float64(order.ID),
            "status": model.OrderStatusPending,
        })
    })

    t.Run("订单不存在", func(t *testing.T) {
        w := th.MakeRequest("GET", "/api/v1/admin/orders/999999", nil, adminToken)
        th.AssertError(t, w, http.StatusNotFound, "order not found")
    })
}

func TestOrderHandler_UpdateOrderStatus(t *testing.T) {
    th := handler.SetupTest()
    adminToken := th.CreateAdminUser()

    t.Run("成功更新状态", func(t *testing.T) {
        order := th.CreateTestOrder(model.OrderStatusPending)

        payload := map[string]interface{}{
            "status": model.OrderStatusConfirmed,
        }
        w := th.MakeRequest("PUT",
            "/api/v1/admin/orders/"+strconv.FormatUint(order.ID, 10)+"/status",
            payload, adminToken)
        th.AssertSuccess(t, w, nil)

        // 验证状态已更新
        w = th.MakeRequest("GET", "/api/v1/admin/orders/"+strconv.FormatUint(order.ID, 10), nil, adminToken)
        th.AssertSuccess(t, w, map[string]interface{}{
            "status": model.OrderStatusConfirmed,
        })
    })

    t.Run("无效状态转换", func(t *testing.T) {
        order := th.CreateTestOrder(model.OrderStatusCompleted)

        payload := map[string]interface{}{
            "status": model.OrderStatusPending, // 不允许从已完成回到待处理
        }
        w := th.MakeRequest("PUT",
            "/api/v1/admin/orders/"+strconv.FormatUint(order.ID, 10)+"/status",
            payload, adminToken)
        th.AssertError(t, w, http.StatusBadRequest, "invalid status transition")
    })

    t.Run("缺少操作原因", func(t *testing.T) {
        order := th.CreateTestOrder(model.OrderStatusPending)

        payload := map[string]interface{}{
            "status": model.OrderStatusCancelled,
            // 缺少 reason
        }
        w := th.MakeRequest("PUT",
            "/api/v1/admin/orders/"+strconv.FormatUint(order.ID, 10)+"/status",
            payload, adminToken)
        th.AssertError(t, w, http.StatusBadRequest, "reason is required")
    })
}
```

3. **用户 Handler 测试**: `api/internal/handler/admin/user_handler_test.go`

4. **陪玩师 Handler 测试**: `api/internal/handler/admin/player_handler_test.go`

5. **运行测试**:

```bash
cd api
# 运行 handler 测试
go test ./internal/handler/... -v -cover

# 生成覆盖率报告
go test ./internal/handler/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

**验收标准**:
- [ ] Handler 层覆盖率达到 30%
- [ ] 核心接口测试完整
- [ ] 边界条件测试覆盖
- [ ] 错误场景测试覆盖

---

## Phase 2: 重要改进 (P1) - 2-4周

### 🟡 P1-SEC-004: JWT 重放攻击防护

**问题**: JWT 缺少 JTI 声明，无法撤销令牌。

#### 修复步骤

1. **修改 Claims 结构**: `api/pkg/auth/jwt.go`

```go
type Claims struct {
    UserID    uint64 `json:"user_id"`
    Role      string `json:"role"`
    JTI       string `json:"jti"`        // JWT ID (唯一标识符)
    SessionID string `json:"session_id"` // 会话标识
    jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token
func (m *JWTManager) GenerateToken(userID uint64, role string) (string, error) {
    jti := generateJTI() // 生成唯一 ID
    sessionID := generateSessionID()

    claims := Claims{
        UserID:    userID,
        Role:      role,
        JTI:       jti,
        SessionID: sessionID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.expiration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(m.secretKey))
}

// RevokeToken 撤销 token
func (m *JWTManager) RevokeToken(ctx context.Context, jti string) error {
    // 在 Redis 中存储撤销的 JTI，TTL 设置为 token 过期时间
    return m.redis.Set(ctx, "revoked:jti:"+jti, "1", m.expiration).Err()
}

// IsTokenRevoked 检查 token 是否已撤销
func (m *JWTManager) IsTokenRevoked(ctx context.Context, jti string) (bool, error) {
    exists, err := m.redis.Exists(ctx, "revoked:jti:"+jti).Result()
    return exists > 0, err
}

// Logout 用户登出，撤销当前 token
func (m *JWTManager) Logout(ctx context.Context, tokenString string) error {
    claims, err := m.ValidateToken(tokenString)
    if err != nil {
        return err
    }

    return m.RevokeToken(ctx, claims.JTI)
}
```

2. **认证中间件检查撤销**: `api/internal/handler/middleware/jwtAuth.go`

```go
func JWTAuth(jwtManager *auth.JWTManager) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractToken(c)
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
            c.Abort()
            return
        }

        claims, err := jwtManager.ValidateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        // 检查 token 是否已撤销
        revoked, err := jwtManager.IsTokenRevoked(c.Request.Context(), claims.JTI)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "token check failed"})
            c.Abort()
            return
        }
        if revoked {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "token has been revoked"})
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("role", claims.Role)
        c.Set("session_id", claims.SessionID)
        c.Next()
    }
}
```

**验收标准**:
- [ ] JWT 包含 JTI 声明
- [ ] Redis 存储撤销列表
- [ ] Logout 正确撤销 token
- [ ] 中间件检查撤销状态

---

### 🟡 P1-SEC-005: Token 存储迁移

**问题**: Token 存储在 localStorage，易被 XSS 窃取。

#### 修复步骤

1. **后端设置 httpOnly Cookie**: `api/internal/handler/auth.go`

```go
func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondAPIError(c, apierr.BadRequest("invalid request"))
        return
    }

    user, err := h.authService.Login(c.Request.Context(), req)
    if err != nil {
        respondAPIError(c, err)
        return
    }

    token, err := h.jwtManager.GenerateToken(user.ID, user.Role)
    if err != nil {
        respondAPIError(c, apierr.InternalError("failed to generate token"))
        return
    }

    // 设置 httpOnly Cookie
    c.SetSameSite(http.SameSiteStrictMode)
    c.SetCookie(
        "auth_token",           // 名称
        token,                  // 值
        3600*24,                // 24小时
        "/",                    // 路径
        "",                     // 域名（生产环境设置具体域名）
        config.IsProduction(),  // Secure (HTTPS only)
        true,                   // HttpOnly (防止 JS 访问)
    )

    resp.OK(c, gin.H{
        "user": user,
        // 不返回 token，token 在 Cookie 中
    })
}

func (h *AuthHandler) Logout(c *gin.Context) {
    // 获取当前 token
    token, err := c.Cookie("auth_token")
    if err == nil {
        // 撤销 token
        h.jwtManager.Logout(c.Request.Context(), token)
    }

    // 清除 Cookie
    c.SetSameSite(http.SameSiteStrictMode)
    c.SetCookie(
        "auth_token",
        "",
        -1,     // 立即过期
        "/",
        "",
        config.IsProduction(),
        true,
    )

    resp.OK(c, nil)
}
```

2. **前端移除 localStorage 管理**: `admin/src/api/client.ts`

```typescript
// ❌ 删除
// const token = localStorage.getItem('token');

// ✅ Cookie 由浏览器自动发送
// 无需手动管理 token

// 登录响应处理
const handleLoginResponse = (response: LoginResponse) => {
    // Token 已在 Cookie 中，只需存储用户信息
  sessionStorage.setItem('user', JSON.stringify(response.user));
  sessionStorage.setItem('user_role', response.user.role);
};

// 登出
const logout = async () => {
    try {
        await apiClient.post('/auth/logout');
    } finally {
        // 清除会话信息（不影响 Cookie）
        sessionStorage.clear();
        // 跳转登录页
        window.location.href = '/login';
    }
};
```

**验收标准**:
- [ ] Token 仅存储在 httpOnly Cookie
- [ ] 前端无法通过 JS 访问 Token
- [ ] 登出正确撤销 Token 和清除 Cookie

---

### 🟡 P1-SEC-006: 加密签名时间戳验证

**问题**: 加密签名的时间戳未验证有效期，可重放。

#### 修复步骤

**修改**: `api/internal/handler/middleware/crypto.go`

```go
const TimestampTolerance = 5 * 60 // 5分钟容差

func validateSignature(plain []byte, req *CryptoRequest, cfg *CryptoConfig) error {
    if req.Signature == "" || req.Timestamp == 0 {
        return errors.New("缺少签名或时间戳")
    }

    now := time.Now().Unix()
    timestamp := req.Timestamp

    // 验证时间戳在允许范围内（±5分钟）
    if timestamp < now-TimestampTolerance || timestamp > now+TimestampTolerance {
        return fmt.Errorf("时间戳无效或已过期: server_time=%d, request_time=%d",
            now, timestamp)
    }

    // 验证签名
    expected := generateSignature(plain, req.Timestamp, cfg.SecretKey)
    if !strings.EqualFold(expected, req.Signature) {
        return errors.New("签名验证失败")
    }

    // 可选：检查时间戳是否已被使用（防止重放）
    // if isTimestampUsed(timestamp) {
    //     return errors.New("时间戳已使用，可能为重放攻击")
    // }

    return nil
}
```

**验收标准**:
- [ ] 时间戳验证 ±5 分钟
- [ ] 过期时间戳被拒绝
- [ ] 未来时间戳被拒绝
- [ ] 签名验证正确

---

### 🟡 P1-CODE-002: 统一日志管理

**问题**: 217 处 console.log，生产环境泄露信息。

#### 修复步骤

1. **创建日志工具**: `admin/src/utils/logger.ts`

```typescript
interface LogContext {
    [key: string]: unknown;
}

class Logger {
    private isDev = import.meta.env.DEV;

    private format(level: string, message: string, context?: LogContext) {
        const timestamp = new Date().toISOString();
        const logEntry = {
            timestamp,
            level,
            message,
            ...(context && { context }),
        };

        // 生产环境发送到日志服务
        if (!this.isDev) {
            this.sendToLogService(logEntry);
            return;
        }

        // 开发环境输出到控制台
        return logEntry;
    }

    private sendToLogService(logEntry: unknown) {
        // TODO: 实现日志服务上报
        // fetch('/api/v1/logs', { method: 'POST', body: JSON.stringify(logEntry) });
    }

    info(message: string, context?: LogContext) {
        if (this.isDev) {
            console.log('[INFO]', this.format('info', message, context));
        }
    }

    warn(message: string, context?: LogContext) {
        if (this.isDev) {
            console.warn('[WARN]', this.format('warn', message, context));
        }
    }

    error(message: string, error?: Error | unknown, context?: LogContext) {
        const errorContext = {
            ...(context || {}),
            ...(error instanceof Error && {
                error: {
                    name: error.name,
                    message: error.message,
                    stack: this.isDev ? error.stack : undefined,
                },
            }),
        };

        if (this.isDev) {
            console.error('[ERROR]', this.format('error', message, errorContext));
        }
    }
}

export const logger = new Logger();
```

2. **批量替换**:

```bash
# 查找所有 console.log
cd admin
grep -r "console\." src/ --include="*.ts" --include="*.tsx" -l

# 手动替换常见模式
# console.log → logger.info
# console.warn → logger.warn
# console.error → logger.error
```

3. **替换示例**:

```typescript
// ❌ 替换前
console.error('Load orders error:', error);
console.log('User logged in:', user);
console.warn('API rate limit approaching');

// ✅ 替换后
import { logger } from '@/utils/logger';

logger.error('Load orders error', error, { action: 'load_orders' });
logger.info('User logged in', { user_id: user.id });
logger.warn('API rate limit approaching', { remaining: 42 });
```

**验收标准**:
- [ ] 所有 console.log 替换为 logger
- [ ] 生产环境无 console 输出
- [ ] 敏感信息不在日志中
- [ ] 日志包含上下文信息

---

### 🟡 P1-ARCH-001: 添加事务管理

**问题**: 跨 Repository 操作无事务边界，数据一致性风险。

#### 修复步骤

1. **创建工作单元**: `api/pkg/uow/work.go`

```go
package uow

import (
    "context"
    "fmt"

    "gorm.io/gorm"
)

// Work 工作单元
type Work struct {
    db    *gorm.DB
    repos map[string]interface{}
    ctx   context.Context
}

// NewWork 创建工作单元
func NewWork(ctx context.Context, db *gorm.DB) *Work {
    return &Work{
        db:    db,
        repos: make(map[string]interface{}),
        ctx:   ctx,
    }
}

// OrderRepo 获取订单仓储
func (w *Work) OrderRepo() OrderRepository {
    if r, ok := w.repos["order"]; ok {
        return r.(OrderRepository)
    }
    r := NewOrderRepository(w.db)
    w.repos["order"] = r
    return r
}

// PaymentRepo 获取支付仓储
func (w *Work) PaymentRepo() PaymentRepository {
    if r, ok := w.repos["payment"]; ok {
        return r.(PaymentRepository)
    }
    r := NewPaymentRepository(w.db)
    w.repos["payment"] = r
    return r
}

// Commit 提交事务
func (w *Work) Commit(fn func(*Work) error) error {
    return w.db.WithContext(w.ctx).Transaction(func(tx *gorm.DB) error {
        w.db = tx
        return fn(w)
    })
}
```

2. **Service 使用工作单元**: `api/internal/service/order/orderService.go`

```go
type OrderService struct {
    db *gorm.DB
    // 移除独立的 repos，通过 UoW 获取
}

func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*Order, error) {
    var order *Order

    // 使用工作单元确保事务
    err := uow.NewWork(ctx, s.db).Commit(func(w *uow.Work) error {
        orderRepo := w.OrderRepo()
        paymentRepo := w.PaymentRepo()

        // 创建订单
        order = &Order{
            UserID:          req.UserID,
            PlayerID:        req.PlayerID,
            TotalPriceCents: req.TotalPriceCents,
            Status:          model.OrderStatusPending,
        }
        if err := orderRepo.Create(ctx, order); err != nil {
            return fmt.Errorf("create order: %w", err)
        }

        // 创建支付记录
        payment := &Payment{
            OrderID:       order.ID,
            AmountCents:   req.TotalPriceCents,
            PaymentMethod: req.PaymentMethod,
            Status:        model.PaymentStatusPending,
        }
        if err := paymentRepo.Create(ctx, payment); err != nil {
            return fmt.Errorf("create payment: %w", err)
        }

        // 记录佣金（在同一个事务中）
        if err := s.recordCommission(ctx, w, order); err != nil {
            return fmt.Errorf("record commission: %w", err)
        }

        return nil
    })

    if err != nil {
        return nil, err
    }

    return order, nil
}

// recordCommission 记录佣金（在事务内）
func (s *OrderService) recordCommission(ctx context.Context, w *uow.Work, order *Order) error {
    commissionRepo := w.CommissionRepo()

    commission := &Commission{
        OrderID:       order.ID,
        PlayerID:      order.PlayerID,
        BaseRate:      0.20, // 从配置获取
        ActualRate:    0.20,
        AmountCents:   int64(float64(order.TotalPriceCents) * 0.20),
        SettlementStatus: model.SettlementStatusPending,
    }

    return commissionRepo.Create(ctx, commission)
}
```

**验收标准**:
- [ ] 所有跨 Repository 操作使用事务
- [ ] 事务失败时数据回滚
- [ ] 测试验证数据一致性

---

### 🟡 P1-TEST-002: 补充前端页面测试

**问题**: 前端页面测试覆盖率仅 6.67% (7/105)。

#### 修复步骤

1. **优先测试页面**:

| 页面 | 优先级 | 测试重点 |
|------|--------|----------|
| `pages/admin/Order/index.tsx` | P0 | 列表加载、筛选、导出、状态更新 |
| `pages/admin/Player/index.tsx` | P0 | 列表加载、审核、排名更新 |
| `pages/admin/User/index.tsx` | P1 | 用户查询、禁用/启用 |
| `pages/admin/Commission/index.tsx` | P1 | 佣金计算、筛选 |
| `pages/auth/Login.tsx` | P0 | 登录流程、错误处理 |

2. **订单页面测试示例**: `admin/src/pages/admin/Order/index.test.tsx`

```typescript
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import OrderPage from './index';
import { BrowserRouter } from 'react-router-dom';

// Mock API
vi.mock('@/api/order', () => ({
    getOrders: vi.fn(),
    updateOrderStatus: vi.fn(),
    exportOrders: vi.fn(),
}));

const mockGetOrders = vi.mocked(getOrders);

describe('OrderPage', () => {
    let queryClient: QueryClient;

    beforeEach(() => {
        queryClient = new QueryClient({
            defaultOptions: {
                queries: { retry: false },
            },
        });
        vi.clearAllMocks();
    });

    const renderWithProviders = (component: React.ReactElement) => {
        return render(
            <QueryClientProvider client={queryClient}>
                <BrowserRouter>{component}</BrowserRouter>
            </QueryClientProvider>
        );
    };

    it('应该成功加载订单列表', async () => {
        const mockOrders = [
            { id: 1, order_no: 'ORD001', status: 'pending', total_price_cents: 10000 },
            { id: 2, order_no: 'ORD002', status: 'confirmed', total_price_cents: 20000 },
        ];
        mockGetOrders.mockResolvedValue({ data: { items: mockOrders, total: 2 } });

        renderWithProviders(<OrderPage />);

        await waitFor(() => {
            expect(screen.getByText('ORD001')).toBeInTheDocument();
            expect(screen.getByText('ORD002')).toBeInTheDocument();
        });
    });

    it('应该显示加载状态', () => {
        mockGetOrders.mockReturnValue(new Promise(() => {})); //永不 resolve

        renderWithProviders(<OrderPage />);

        expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
    });

    it('应该显示错误状态', async () => {
        mockGetOrders.mockRejectedValue(new Error('API Error'));

        renderWithProviders(<OrderPage />);

        await waitFor(() => {
            expect(screen.getByText(/加载失败/)).toBeInTheDocument();
        });
    });

    it('应该支持状态筛选', async () => {
        const user = userEvent.setup();
        mockGetOrders.mockResolvedValue({ data: { items: [], total: 0 } });

        renderWithProviders(<OrderPage />);

        // 选择状态
        const statusSelect = screen.getByLabelText('订单状态');
        await user.click(statusSelect);
        await user.click(screen.getByText('已确认'));

        // 验证 API 调用
        await waitFor(() => {
            expect(mockGetOrders).toHaveBeenCalledWith(
                expect.objectContaining({
                    params: expect.objectContaining({
                        status: 'confirmed',
                    }),
                })
            );
        });
    });

    it('应该支持导出功能', async () => {
        const user = userEvent.setup();
        mockGetOrders.mockResolvedValue({ data: { items: [], total: 0 } });

        renderWithProviders(<OrderPage />);

        const exportButton = screen.getByText('导出');
        await user.click(exportButton);

        // 验证导出 API 调用
        await waitFor(() => {
            expect(exportOrders).toHaveBeenCalled();
        });
    });
});
```

**验收标准**:
- [ ] 核心页面测试覆盖率达到 60%
- [ ] 关键交互流程测试完整
- [ ] 错误场景测试覆盖

---

## Phase 3: 长期优化 (P2) - 1-3个月

### 🟢 P2-ARCH-002: 引入断路器模式

**问题**: 外部服务调用无熔断保护，存在雪崩风险。

#### 实现步骤

```go
// pkg/circuitbreaker/circuitbreaker.go
type CircuitBreaker struct {
    maxFailures  int
    resetTimeout time.Duration
    state        State
    failures     int
    lastFailure  time.Time
    mu           sync.RWMutex
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
    if cb.isOpen() {
        return ErrCircuitBreakerOpen
    }

    err := fn()
    if err != nil {
        cb.recordFailure()
        return err
    }

    cb.recordSuccess()
    return nil
}
```

---

### 🟢 P2-ARCH-003: 实现分布式锁

**问题**: 高并发场景（如抢单）无锁保护。

#### 实现步骤

```go
// pkg/lock/redis_lock.go
type RedisLock struct {
    client *redis.Client
}

func (l *RedisLock) TryLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
    return l.client.SetNX(ctx, key, "1", ttl).Result()
}

func (l *RedisLock) Unlock(ctx context.Context, key string) error {
    return l.client.Del(ctx, key).Err()
}
```

---

### 🟢 P2-ARCH-004: 实现幂等性

**问题**: 订单/支付操作无幂等键。

#### 实现步骤

```go
// middleware/idempotency.go
func Idempotency(cache cache.Cache) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := c.GetHeader("X-Idempotency-Key")
        if key == "" {
            c.Next()
            return
        }

        // 检查缓存
        if cached, err := cache.Get(c, key); err == nil {
            c.JSON(http.StatusOK, cached)
            c.Abort()
            return
        }

        // 记录响应
        writer := &responseWriter{ResponseWriter: c.Writer}
        c.Writer = writer
        c.Next()

        if c.Writer.Status() < 400 {
            cache.Set(c, key, writer.body, 24*time.Hour)
        }
    }
}
```

---

### 🟢 P2-CODE-003: 清理 TODO 注释

**问题**: 66 处 TODO 未处理。

#### 处理步骤

```bash
# 查找所有 TODO
cd api
grep -r "TODO" --include="*.go" . | wc -l

# 分类处理
# 1. 实现功能 (创建 Issue)
# 2. 删除过时 TODO
# 3. 转换为文档
```

---

### 🟢 P2-CODE-004: 拆分大型组件

**问题**: OrderPage 775 行，职责过多。

#### 拆分方案

```
OrderPage/
├── index.tsx           # 主页面 (200 行)
├── OrderList.tsx       # 列表组件 (150 行)
├── OrderDetail.tsx     # 详情组件 (150 行)
├── OrderFilter.tsx     # 筛选组件 (100 行)
├── OrderActions.tsx    # 操作组件 (100 行)
└── OrderExport.tsx     # 导出组件 (75 行)
```

---

### 🟢 P2-TEST-003: 添加 E2E 测试

**问题**: 缺少端到端测试。

#### 实现步骤

```bash
cd admin
npm install -D @playwright/test
```

```typescript
// e2e/order.spec.ts
import { test, expect } from '@playwright/test';

test('订单管理流程', async ({ page }) => {
    // 登录
    await page.goto('/login');
    await page.fill('[name="email"]', 'admin@example.com');
    await page.fill('[name="password"]', 'password');
    await page.click('button[type="submit"]');

    // 导航到订单页面
    await page.click('text=订单管理');
    await expect(page).toHaveURL(/\/admin\/orders/);

    // 筛选订单
    await page.selectOption('select[name="status"]', 'confirmed');
    await page.click('button:has-text("查询")');

    // 验证结果
    await expect(page.locator('.order-table')).toBeVisible();
});
```

---

## 执行检查清单

### Phase 1 完成标准

- [ ] 所有 P0 任务已完成
- [ ] 安全扫描无高危漏洞
- [ ] 测试覆盖率达到目标
- [ ] 代码审查通过

### Phase 2 完成标准

- [ ] 所有 P1 任务已完成
- [ ] 中危漏洞已修复
- [ ] 前端测试覆盖率 > 50%
- [ ] 性能测试通过

### Phase 3 完成标准

- [ ] 所有 P2 任务已完成
- [ ] 架构优化完成
- [ ] 文档更新完整
- [ ] 技术债务清理

---

## 附录

### A. 工具安装

```bash
# 后端工具
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

# 前端工具
npm install -D @playwright/test
npm install -D vitest @testing-library/react
```

### B. 参考资源

- [OWASP Top 10 2021](https://owasp.org/Top10/)
- [Go Secure Coding](https://cheatsheetseries.owasp.org/cheatsheets/Go_Security_Cheat_Sheet.html)
- [React Testing Library](https://testing-library.com/docs/react-testing-library/intro/)
