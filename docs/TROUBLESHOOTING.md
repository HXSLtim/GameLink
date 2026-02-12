# GameLink 故障排查手册

> **版本**: v2.0
> **最后更新**: 2026-02-11
> **适用**: 开发/测试/生产环境

---

## 目录

1. [启动问题](#1-启动问题)
2. [运行时问题](#2-运行时问题)
3. [性能问题](#3-性能问题)
4. [数据库问题](#4-数据库问题)
5. [前端问题](#5-前端问题)
6. [第三方服务问题](#6-第三方服务问题)
7. [WebSocket 问题](#7-websocket-问题)
8. [部署问题](#8-部署问题)

---

## 1. 启动问题

### 1.1 数据库连接失败

**症状**:
```
failed to connect to database: connection refused
```

**原因分析**:
- PostgreSQL 未启动
- 连接地址/端口配置错误
- 数据库用户密码错误
- 防火墙阻止连接

**诊断步骤**:
```bash
# 1. 检查 PostgreSQL 是否运行
sudo systemctl status postgresql
# 或 Docker 环境
docker compose ps postgres

# 2. 测试数据库连接
psql -h localhost -U gamelink -d gamelink_dev

# 3. 检查防火墙
sudo ufw status
sudo iptables -L -n

# 4. 查看监听端口
sudo netstat -tlnp | grep 5432
```

**解决方案**:
```bash
# 启动 PostgreSQL
sudo systemctl start postgresql

# Docker 环境
docker compose up -d postgres

# 修改 pg_hba.conf 允许连接
sudo nano /etc/postgresql/16/main/pg_hba.conf
# 添加: host all all 0.0.0.0/0 md5

# 重启 PostgreSQL
sudo systemctl restart postgresql
```

---

### 1.2 Redis 连接失败

**症状**:
```
redis: connection refused
```

**原因分析**:
- Redis 未启动
- 密码配置错误
- 端口被占用

**诊断步骤**:
```bash
# 1. 检查 Redis 是否运行
sudo systemctl status redis
# 或 Docker
docker compose ps redis

# 2. 测试连接
redis-cli -h localhost -p 6379 -a your_password ping

# 3. 查看端口占用
sudo netstat -tlnp | grep 6379
```

**解决方案**:
```bash
# 启动 Redis
sudo systemctl start redis

# Docker 环境
docker compose up -d redis

# 修改 redis.conf 绑定地址
bind 0.0.0.0
```

---

### 1.3 端口占用

**症状**:
```
bind: address already in use
```

**原因分析**: 端口已被其他进程占用

**诊断步骤**:
```bash
# 查看端口占用
# Linux/Mac
sudo lsof -i :8080

# Windows
netstat -ano | findstr :8080
```

**解决方案**:
```bash
# 方法1: 杀死占用进程
sudo kill -9 <PID>

# 方法2: 修改应用端口
# 编辑 .env
PORT=8081

# 方法3: 修改管理后台端口
# 编辑 admin/vite.config.ts
server: {
  port: 5174
}
```

---

### 1.4 环境变量配置错误

**症状**: 应用启动但行为异常

**诊断步骤**:
```bash
# 检查环境变量文件
cat .env

# 验证必需变量
env | grep -E "DB_|REDIS_|JWT_"
```

**解决方案**:
```bash
# 复制模板重新配置
cp .env.example .env

# 必需配置项
APP_ENV=development
DB_HOST=localhost
DB_NAME=gamelink_dev
DB_USER=gamelink
DB_PASSWORD=your_password
REDIS_HOST=localhost
REDIS_PASSWORD=your_password
JWT_SECRET_KEY=your_secret_key_min_32_chars
```

---

## 2. 运行时问题

### 2.1 API 请求失败 (404)

**症状**:
```
GET /api/v1/users - 404 Not Found
```

**原因分析**:
- 路由未注册
- 路径拼写错误
- 中间件拦截

**诊断步骤**:
```bash
# 1. 检查 Swagger 文档
# 访问 http://localhost:8080/swagger/index.html

# 2. 查看后端日志
journalctl -u gamelink-api -f
# 或开发环境
tail -f api/logs/app.log

# 3. 测试路由
curl http://localhost:8080/api/v1/users
```

**解决方案**:
```bash
# 1. 确认路由已注册
# 检查 api/internal/router/router.go

# 2. 检查 Handler 方法
# 确认方法已导出（首字母大写）

# 3. 检查中间件
# 确认没有错误拦截请求
```

---

### 2.2 认证失败 (401)

**症状**:
```
GET /api/v1/admin/users - 401 Unauthorized
```

**原因分析**:
- Token 过期
- Token 格式错误
- JWT 密钥不匹配

**诊断步骤**:
```bash
# 1. 检查 Token
# 浏览器 DevTools -> Application -> Local Storage

# 2. 解码 JWT
# 访问 https://jwt.io/
# 检查 exp 是否过期

# 3. 查看后端日志
# 搜索 "jwt" 或 "auth" 关键词
```

**解决方案**:
```bash
# 前端: 重新登录获取新 Token
# 后端: 检查 JWT_SECRET_KEY 配置

# Token 刷新
POST /api/v1/auth/refresh
{
  "refreshToken": "your_refresh_token"
}
```

---

### 2.3 权限不足 (403)

**症状**:
```
DELETE /api/v1/admin/users/1 - 403 Forbidden
```

**原因分析**:
- 用户角色无该权限
- RBAC 配置错误

**诊断步骤**:
```bash
# 1. 检查用户角色
GET /api/v1/auth/me

# 2. 检查角色权限
GET /api/v1/admin/roles

# 3. 检查所需权限码
# 查看对应 Handler 的权限注解
```

**解决方案**:
```sql
-- 给角色分配权限
INSERT INTO role_permissions (role_id, permission_id)
VALUES (1, (SELECT id FROM permissions WHERE code = 'admin.users.delete'));

-- 或直接授予超级管理员权限
UPDATE users SET role = 'admin' WHERE email = 'your@email.com';
```

---

### 2.4 请求超时

**症状**: 请求长时间无响应

**原因分析**:
- 数据库查询慢
- 外部 API 调用超时
- 死锁

**诊断步骤**:
```bash
# 1. 查看应用日志
# 搜索 "timeout" 或 "slow query"

# 2. 检查数据库慢查询
SELECT * FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;

# 3. 检查锁等待
SELECT * FROM pg_stat_activity
WHERE state = 'active'
AND wait_event_type = 'Lock';
```

**解决方案**:
```sql
-- 1. 终止长时间运行的查询
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE state = 'active'
AND query_start < now() - interval '5 minutes';

-- 2. 添加索引
CREATE INDEX CONCURRENTLY idx_orders_user_status
ON orders (user_id, status);

-- 3. 优化查询
-- 使用 LIMIT、避免 SELECT * 等
```

---

## 3. 性能问题

### 3.1 API 响应慢

**症状**: 接口响应时间 > 1s

**诊断步骤**:
```bash
# 1. 使用 pprof 性能分析
# 添加到 main.go
import _ "net/http/pprof"

# 访问
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30

# 2. 查看 CPU 分析
go tool pprof -http=:8080 cpu.prof

# 3. 查看内存分析
go tool pprof -http=:8080 heap.prof
```

**优化方案**:
```go
// 1. 添加缓存
func (s *orderService) GetOrders(ctx context.Context, userID int) ([]*Order, error) {
    // 先查缓存
    cacheKey := fmt.Sprintf("orders:user:%d", userID)
    if cached := s.redis.Get(ctx, cacheKey); cached != nil {
        return cached, nil
    }

    // 查数据库
    orders, err := s.repo.GetByUser(ctx, userID)
    if err != nil {
        return nil, err
    }

    // 写缓存
    s.redis.Set(ctx, cacheKey, orders, 5*time.Minute)
    return orders, nil
}

// 2. 使用 goroutine 并发
func (s *userService) GetUserWithOrders(ctx context.Context, userID int) (*UserDetail, error) {
    var user *User
    var orders []Order
    var err error

    // 使用 errgroup
    g, ctx := errgroup.WithContext(ctx)

    g.Go(func() error {
        user, err = s.userRepo.Get(ctx, userID)
        return err
    })

    g.Go(func() error {
        orders, err = s.orderRepo.GetByUser(ctx, userID)
        return err
    })

    if err := g.Wait(); err != nil {
        return nil, err
    }

    return &UserDetail{User: user, Orders: orders}, nil
}
```

---

### 3.2 内存占用高

**症状**: 应用内存持续增长

**诊断步骤**:
```bash
# 1. 查看进程内存
top -p <pid>

# 2. Go 内存统计
curl http://localhost:8080/debug/pprof/heap > heap.prof
go tool pprof -http=:8080 heap.prof

# 3. 查看垃圾回收
go tool pprof http://localhost:8080/debug/pprof/heap
```

**解决方案**:
```go
// 1. 释放大对象
func processLargeData() {
    data := make([]byte, 100*1024*1024) // 100MB
    // 使用完后置空
    data = nil
    runtime.GC() // 手动触发 GC
}

// 2. 使用对象池
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func useBuffer() *bytes.Buffer {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()
    return buf
}

// 3. 避免 goroutine 泄漏
// 使用 context 控制 goroutine 生命周期
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case task := <-tasks:
            process(task)
        }
    }
}
```

---

### 3.3 数据库连接池耗尽

**症状**:
```
connection pool exhausted
```

**诊断步骤**:
```bash
# 1. 查看当前连接数
SELECT count(*) FROM pg_stat_activity;

# 2. 查看连接详情
SELECT datname, usename, state, count(*)
FROM pg_stat_activity
GROUP BY datname, usename, state;

# 3. 查看最大连接数
SHOW max_connections;
```

**解决方案**:
```go
// 调整连接池配置
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    &gorm.Config{
        ConnPool: &stdsql.DB{
            MaxIdleConns: 10,  // 空闲连接数
            MaxOpenConns: 100, // 最大连接数
            ConnMaxLifetime: time.Hour, // 连接最大生存时间
        },
    },
})
```

---

## 4. 数据库问题

### 4.1 慢查询

**症状**: 查询响应时间长

**诊断步骤**:
```sql
-- 1. 启用慢查询日志
ALTER SYSTEM SET log_min_duration_statement = 1000; -- 1秒
SELECT pg_reload_conf();

-- 2. 查看慢查询
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 20;

-- 3. 使用 EXPLAIN 分析
EXPLAIN ANALYZE
SELECT * FROM orders
WHERE user_id = 1 AND status = 'pending'
ORDER BY created_at DESC;
```

**解决方案**:
```sql
-- 1. 添加索引
CREATE INDEX CONCURRENTLY idx_orders_user_status_created
ON orders (user_id, status, created_at DESC);

-- 2. 使用覆盖索引
CREATE INDEX idx_orders_user_status_covering
ON orders (user_id, status, created_at DESC)
INCLUDE (id, total_price_cents);

-- 3. 优化查询
-- 避免 SELECT *
-- 使用 LIMIT 限制返回行数
-- 避免 LIKE '%keyword%'，使用全文索引
```

---

### 4.2 死锁

**症状**: 事务长时间等待

**诊断步骤**:
```sql
-- 查看锁等待
SELECT
    l.locktype,
    l.relation::regclass,
    l.mode,
    a.usename,
    a.query,
    a.query_start
FROM pg_locks l
JOIN pg_stat_activity a ON l.pid = a.pid
WHERE NOT granted;
```

**解决方案**:
```sql
-- 1. 终止阻塞进程
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE state = 'active'
AND pid IN (
    SELECT blocking_pid
    FROM pg_stat_activity
);

-- 2. 优化事务
-- 尽量缩短事务时间
-- 保持一致的访问顺序

-- 3. 应用层优化
func (s *orderService) UpdateOrder(ctx context.Context, id int, updates map[string]interface{}) error {
    // 使用行级锁
    return s.db.Transaction(func(tx *gorm.DB) error {
        var order Order
        // FOR UPDATE 行级锁
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, id).Error; err != nil {
            return err
        }
        return tx.Model(&order).Updates(updates).Error
    })
}
```

---

## 5. 前端问题

### 5.1 白屏

**症状**: 页面加载后白屏

**诊断步骤**:
```bash
# 1. 查看浏览器控制台
# F12 -> Console

# 2. 查看网络请求
# F12 -> Network

# 3. 检查构建日志
npm run build
```

**解决方案**:
```javascript
// 1. 检查路由配置
// 确保路由与实际组件匹配

// 2. 检查环境变量
// console.log(import.meta.env)

// 3. 检查 API 地址
// vite.config.ts -> server -> proxy
```

---

### 5.2 API 跨域

**症状**:
```
Access to XMLHttpRequest has been blocked by CORS policy
```

**解决方案**:
```go
// 后端配置 CORS
router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5173", "https://admin.yourdomain.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}))
```

---

## 6. 第三方服务问题

### 6.1 微信支付回调失败

**症状**: 支付成功但订单状态未更新

**诊断步骤**:
```bash
# 1. 检查回调 URL 可访问性
curl https://api.yourdomain.com/api/v1/payments/wechat/callback

# 2. 查看后端日志
# 搜索 "wechat" 或 "callback"

# 3. 检查微信商户平台
# 查看回调日志
```

**解决方案**:
```go
// 1. 确保回调 URL 可被公网访问
// 不能是 localhost

// 2. 添加日志记录回调
func (h *PaymentHandler) WechatCallback(c *gin.Context) {
    log.Info("Wechat callback received", "body", c.Request.Body)
    // ...
}

// 3. 处理重复回调
// 使用幂等性检查
if payment.Status == "paid" {
    return c.JSON(200, gin.H{"code": "SUCCESS"})
}
```

---

## 7. WebSocket 问题

### 7.1 连接频繁断开

**症状**: WebSocket 不断重连

**诊断步骤**:
```javascript
// 前端添加事件监听
ws.addEventListener('close', (event) => {
    console.log('WebSocket closed:', event.code, event.reason);
});

ws.addEventListener('error', (error) => {
    console.log('WebSocket error:', error);
});
```

**解决方案**:
```go
// 1. 调整心跳间隔
// 前端: 30秒发送一次 ping
setInterval(() => {
    ws.send(JSON.stringify({ type: 'ping' }));
}, 30000);

// 后端: 60秒超时
func (h *Hub) readPump(ws *WebSocket) {
    ws.SetReadDeadline(time.Now().Add(60 * time.Second))
}

// 2. 添加断线重连
let reconnectAttempts = 0;
const maxReconnectAttempts = 5;

function connect() {
    ws = new WebSocket(url);

    ws.addEventListener('close', () => {
        if (reconnectAttempts < maxReconnectAttempts) {
            setTimeout(() => {
                reconnectAttempts++;
                connect();
            }, 5000);
        }
    });
}
```

---

## 8. 部署问题

### 8.1 Docker 镜像构建失败

**症状**:
```
ERROR [build] failed to compute cache key
```

**解决方案**:
```bash
# 1. 清理缓存
docker builder prune

# 2. 使用 --no-cache
docker build --no-cache -t gamelink-api .

# 3. 检查 Dockerfile
# 确保依赖文件存在
```

---

## 附录：快速诊断脚本

```bash
#!/bin/bash
# health_check.sh - 快速健康检查

echo "=== GameLink 健康检查 ==="

# API 检查
echo "1. API 服务..."
curl -sf http://localhost:8080/health || echo "❌ API 不可用"

# 数据库检查
echo "2. PostgreSQL..."
pg_isready -h localhost -p 5432 || echo "❌ PostgreSQL 不可用"

# Redis 检查
echo "3. Redis..."
redis-cli -h localhost -r 0 ping || echo "❌ Redis 不可用"

# 磁盘空间
echo "4. 磁盘空间..."
df -h | grep -E '(Filesystem|/$)'

# 内存使用
echo "5. 内存使用..."
free -h

# 进程状态
echo "6. 应用进程..."
ps aux | grep gamelink-api | grep -v grep
```

---

**文档维护**: 产品经理 + DevOps
**更新频率**: 发现新问题时更新
