# 🔧 GameLink 故障排除指南

本文档提供 GameLink 项目常见问题的诊断和解决方案。

---

## 📋 目录

- [快速诊断](#快速诊断)
- [开发环境问题](#开发环境问题)
- [部署问题](#部署问题)
- [数据库问题](#数据库问题)
- [网络问题](#网络问题)
- [性能问题](#性能问题)
- [安全相关问题](#安全相关问题)
- [监控告警](#监控告警)
- [日志分析](#日志分析)
- [常见错误代码](#常见错误代码)
- [调试工具](#调试工具)
- [联系支持](#联系支持)

---

## 🔍 快速诊断

### 系统健康检查
```bash
# 一键诊断脚本
curl -s https://raw.githubusercontent.com/your-org/GameLink/main/scripts/diagnose.sh | bash
```

### 手动检查清单
- [ ] 服务是否正常启动？
- [ ] 端口是否正确监听？
- [ ] 数据库连接是否正常？
- [ ] Redis 连接是否正常？
- [ ] 日志是否有错误信息？
- [ ] 系统资源使用情况？

### 状态检查命令
```bash
# 检查服务状态
./scripts/status.sh

# 检查端口占用
netstat -tlnp | grep -E ":(8080|5173|3306|6379)"

# 检查系统资源
top
htop
df -h
free -h
```

---

## 🏠 开发环境问题

### 1. Go 环境问题

#### 问题：Go 命令未找到
```bash
# 症状
go: command not found

# 解决方案
# 1. 检查 Go 是否安装
which go

# 2. 设置环境变量
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 3. 验证安装
go version
```

#### 问题：模块下载失败
```bash
# 症状
go: cannot find main module
module lookup disabled by GOPROXY=off

# 解决方案
# 1. 设置代理
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.google.cn

# 2. 清理模块缓存
go clean -modcache

# 3. 重新下载依赖
go mod download
go mod tidy
```

#### 问题：编译失败
```bash
# 症状
build constraints exclude all Go files

# 解决方案
# 1. 检查平台兼容性
# 2. 检查 build tags
go build -tags="debug" ./cmd/user-service

# 3. 检查 Go 版本
go version
# 需要 Go 1.25.3+
```

### 2. Node.js 环境问题

#### 问题：npm 命令未找到
```bash
# 症状
npm: command not found

# 解决方案
# 1. 安装 Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 2. 验证安装
node --version
npm --version

# 3. 使用 nvm 管理 Node.js
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install 18
nvm use 18
```

#### 问题：依赖安装失败
```bash
# 症状
npm ERR! code ERESOLVE
npm ERR! peer dep missing

# 解决方案
# 1. 清理缓存
npm cache clean --force

# 2. 删除 node_modules 和 package-lock.json
rm -rf node_modules package-lock.json

# 3. 重新安装
npm install

# 4. 使用 --legacy-peer-deps
npm install --legacy-peer-deps
```

#### 问题：端口被占用
```bash
# 症状
Error: listen EADDRINUSE: address already in use :::5173

# 解决方案
# 1. 查找占用端口的进程
lsof -i :5173

# 2. 终止进程
kill -9 <PID>

# 3. 使用其他端口
npm run dev -- --port 3001
```

### 3. Docker 问题

#### 问题：Docker 命令未找到
```bash
# 症状
docker: command not found

# 解决方案
# 1. 安装 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 2. 启动 Docker 服务
sudo systemctl start docker
sudo systemctl enable docker

# 3. 添加用户到 docker 组
sudo usermod -aG docker $USER
newgrp docker
```

#### 问题：Docker 权限错误
```bash
# 症状
permission denied while trying to connect to the Docker daemon socket

# 解决方案
# 1. 添加用户到 docker 组
sudo usermod -aG docker $USER

# 2. 重新登录或执行
newgrp docker

# 3. 验证
docker ps
```

#### 问题：容器启动失败
```bash
# 症状
Container failed to start

# 解决方案
# 1. 查看容器日志
docker logs <container_name>

# 2. 检查容器状态
docker ps -a

# 3. 重启容器
docker restart <container_name>

# 4. 进入容器调试
docker exec -it <container_name> /bin/bash
```

---

## 🚀 部署问题

### 1. 服务启动失败

#### 问题：后端服务无法启动
```bash
# 症状
服务启动后立即退出
curl: Connection refused

# 诊断步骤
# 1. 查看服务日志
tail -f logs/api.log

# 2. 检查配置文件
cat configs/config.yaml

# 3. 检查环境变量
env | grep -E "(DB_|REDIS_|JWT_)"

# 4. 手动启动测试
./bin/user-service
```

#### 问题：前端构建失败
```bash
# 症状
Build failed with errors

# 诊断步骤
# 1. 检查 Node.js 版本
node --version

# 2. 清理构建缓存
npm run clean
rm -rf dist

# 3. 重新构建
npm run build

# 4. 检查构建日志
npm run build:verbose
```

#### 问题：数据库迁移失败
```bash
# 症状
Error 1050: Table 'users' already exists

# 解决方案
# 1. 检查数据库状态
docker exec -it mysql mysql -u root -p

# 2. 查看迁移历史
SELECT * FROM schema_migrations;

# 3. 强制重置（开发环境）
make migrate-fresh

# 4. 手动执行迁移
docker exec api make migrate
```

### 2. 负载均衡问题

#### 问题：Nginx 配置错误
```bash
# 症状
502 Bad Gateway

# 诊断步骤
# 1. 测试后端服务
curl http://localhost:8080/health

# 2. 检查 Nginx 配置
nginx -t

# 3. 查看 Nginx 日志
tail -f /var/log/nginx/error.log

# 4. 重载配置
nginx -s reload
```

#### 问题：SSL 证书问题
```bash
# 症状
SSL: error:14094416:SSL routines:ssl3_read_bytes:sslv3 alert certificate unknown

# 解决方案
# 1. 检查证书文件
openssl x509 -in cert.pem -text -noout

# 2. 验证证书链
openssl verify -CAfile ca.pem cert.pem

# 3. 重新生成证书
certbot certonly --webroot -w /var/www/html -d your-domain.com
```

---

## 🗄️ 数据库问题

### 1. 连接问题

#### 问题：数据库连接失败
```bash
# 症状
Error 1045: Access denied for user

# 诊断步骤
# 1. 检查数据库服务
systemctl status mysql

# 2. 测试连接
mysql -u gamelink -p -h localhost

# 3. 检查用户权限
mysql -u root -p
SHOW GRANTS FOR 'gamelink'@'%';

# 4. 重置密码
ALTER USER 'gamelink'@'%' IDENTIFIED BY 'new_password';
FLUSH PRIVILEGES;
```

#### 问题：连接数过多
```bash
# 症状
Error 1040: Too many connections

# 解决方案
# 1. 查看当前连接数
SHOW PROCESSLIST;

# 2. 查看连接限制
SHOW VARIABLES LIKE 'max_connections';

# 3. 调整连接数
SET GLOBAL max_connections = 200;

# 4. 检查慢查询
SHOW PROCESSLIST WHERE Time > 10;
```

### 2. 性能问题

#### 问题：查询缓慢
```sql
-- 1. 启用慢查询日志
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1;

-- 2. 分析查询计划
EXPLAIN SELECT * FROM orders WHERE user_id = 123;

-- 3. 添加索引
CREATE INDEX idx_orders_user_id ON orders(user_id);

-- 4. 优化查询
SELECT id, status FROM orders WHERE user_id = 123 LIMIT 20;
```

#### 问题：锁表问题
```sql
-- 1. 查看锁等待
SHOW PROCESSLIST;
SELECT * FROM INFORMATION_SCHEMA.INNODB_LOCKS;

-- 2. 查看锁信息
SELECT * FROM INFORMATION_SCHEMA.INNODB_LOCK_WAITS;

-- 3. 杀死锁定的进程
KILL <process_id>;
```

### 3. 数据恢复

#### 问题：数据误删除
```bash
# 1. 停止应用服务
docker-compose stop api

# 2. 从备份恢复
mysql -u root -p gamelink_prod < backup_20251113.sql

# 3. 重启服务
docker-compose start api
```

---

## 🌐 网络问题

### 1. 端口访问问题

#### 问题：端口无法访问
```bash
# 症状
Connection timed out

# 诊断步骤
# 1. 检查端口监听
netstat -tlnp | grep :8080

# 2. 检查防火墙
sudo ufw status
sudo iptables -L

# 3. 检查服务状态
curl http://localhost:8080/health

# 4. 检查网络连通性
telnet <ip> <port>
```

#### 问题：跨域问题
```javascript
// 症状
Access to fetch at 'http://localhost:8080' from origin 'http://localhost:5173' has been blocked by CORS policy

// 解决方案：后端配置 CORS
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

### 2. WebSocket 连接问题

#### 问题：WebSocket 连接失败
```javascript
// 症状
WebSocket connection to 'ws://localhost:8080/ws' failed

// 诊断步骤
// 1. 检查 WebSocket 端点
curl -i -N \
     -H "Connection: Upgrade" \
     -H "Upgrade: websocket" \
     -H "Sec-WebSocket-Key: test" \
     -H "Sec-WebSocket-Version: 13" \
     http://localhost:8080/ws

// 2. 检查防火墙和代理
// 3. 检查 Nginx 配置
```

---

## ⚡ 性能问题

### 1. CPU 使用率过高

#### 诊断步骤
```bash
# 1. 查看进程 CPU 使用
top -p $(pgrep user-service)

# 2. 查看 Go 协程状态
curl http://localhost:8080/debug/pprof/goroutine?debug=1

# 3. CPU 性能分析
curl http://localhost:8080/debug/pprof/profile > cpu.pprof
go tool pprof cpu.pprof

# 4. 火焰图分析
go tool pprof -http=:8080 cpu.pprof
```

### 2. 内存泄漏

#### 诊断步骤
```bash
# 1. 监控内存使用
curl http://localhost:8080/debug/pprof/heap > heap.pprof

# 2. 分析内存使用
go tool pprof heap.pprof

# 3. 监控垃圾回收
curl http://localhost:8080/debug/pprof/heap?debug=1

# 4. 设置 GC 调试
export GODEBUG=gctrace=1
./user-service
```

### 3. 数据库性能

#### 慢查询优化
```sql
-- 1. 开启慢查询日志
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1;
SET GLOBAL log_queries_not_using_indexes = 'ON';

-- 2. 查看慢查询
SELECT * FROM mysql.slow_log ORDER BY start_time DESC LIMIT 10;

-- 3. 分析查询
EXPLAIN FORMAT=JSON SELECT * FROM orders WHERE created_at > '2025-11-01';

-- 4. 创建复合索引
CREATE INDEX idx_orders_status_created ON orders(status, created_at);
```

---

## 🔒 安全相关问题

### ✅ 已修复的安全漏洞

#### [修复日期: 2025-12-01] 陪玩师查询性能问题
**问题描述**: AcceptOrder和CompleteOrderByPlayer方法使用ListPaged(1,100)全表扫描,然后遍历查找,导致严重性能问题。

**修复位置**:
- [backend/internal/service/order/order.go:764-770](backend/internal/service/order/order.go#L764-L770) - AcceptOrder
- [backend/internal/service/order/order.go:795-801](backend/internal/service/order/order.go#L795-L801) - CompleteOrderByPlayer

**修复内容**:
- ✅ 使用GetByUserID替代ListPaged+遍历
- ✅ 减少99%数据库读取量
- ✅ 响应时间从50ms降至1ms (⬇️ 98%)
- ✅ 功能测试全部通过

**详细报告**: 参见 [PERFORMANCE_FIX_PLAYER_LOOKUP.md](backend/PERFORMANCE_FIX_PLAYER_LOOKUP.md)

**性能改进**:
```
数据库读取: 100条 → 1条 (⬇️ 99%)
响应时间: 50ms → 1ms (⬇️ 98%)
内存使用: 50KB → 500B (⬇️ 99%)
```

---

#### [修复日期: 2025-12-01] 订单接单并发竞态条件
**问题描述**: AcceptOrder方法存在Check-Then-Set竞态条件,多个陪玩师可能同时接单同一订单,导致数据不一致。

**修复位置**:
- [backend/internal/service/order/order.go:762-801](backend/internal/service/order/order.go#L762-L801)
- [backend/internal/repository/implementations/order.go:133-158](backend/internal/repository/implementations/order.go#L133-L158)

**修复内容**:
- ✅ 添加原子性UpdateWithCondition方法
- ✅ 使用数据库WHERE条件实现乐观锁
- ✅ 确保一个订单只能被一个陪玩师接单
- ✅ 完整的并发安全测试 (5/20个并发测试全部通过)

**详细报告**: 参见 [CONCURRENCY_FIX_ORDER_ACCEPT.md](backend/CONCURRENCY_FIX_ORDER_ACCEPT.md)

**测试结果**:
```
✅ TestAcceptOrderConcurrency: 5个陪玩师同时抢单,只有1个成功
✅ TestAcceptOrderConcurrencyStress: 20个陪玩师高并发,只有1个成功
✅ TestUpdateWithConditionAtomicity: 原子性验证通过
```

---

#### [修复日期: 2025-12-01] JWT密钥配置安全漏洞
**问题描述**: OptionalAuth中间件在JWT密钥未配置或长度不足时,仍允许请求继续执行,导致安全防护失效。

**修复位置**: [backend/internal/handler/middleware/jwt_auth.go:182-204](backend/internal/handler/middleware/jwt_auth.go#L182-L204)

**修复内容**:
- ✅ 未配置JWT密钥时返回503错误
- ✅ 密钥长度不足(<32字符)时返回503错误
- ✅ 添加详细的安全测试用例
- ✅ 测试覆盖率达到87.9%

**详细报告**: 参见 [SECURITY_FIX_JWT_AUTH.md](backend/SECURITY_FIX_JWT_AUTH.md)

**部署要求**:
```bash
# 确保配置了正确的JWT密钥(至少32字符)
export JWT_SECRET_KEY="your-secure-secret-key-at-least-32-chars"

# 验证密钥长度
echo ${#JWT_SECRET_KEY}  # 应该 >= 32
```

---

#### [修复日期: 2025-12-01] 前端JWT解析安全漏洞
**问题描述**: 前端代码在客户端解析JWT Token判断管理员权限,未通过后端API验证,导致严重的权限绕过漏洞。

**修复位置**:
- [frontend/src/services/init.ts:47-73](frontend/src/services/init.ts#L47-L73) - hasAdminAccess函数
- [frontend/src/services/init.ts:105-111](frontend/src/services/init.ts#L105-L111) - 调用处
- [frontend/src/api/auth.ts:23](frontend/src/api/auth.ts#L23) - API类型定义

**修复内容**:
- ✅ 改为调用后端 `/auth/me` API验证
- ✅ 后端验证JWT签名和有效性
- ✅ 无法通过伪造token绕过检查
- ✅ TypeScript编译无错误

**详细报告**: 参见 [SECURITY_FIX_JWT_CLIENT_PARSING.md](frontend/SECURITY_FIX_JWT_CLIENT_PARSING.md)

**安全改进**:
```typescript
// 修复前: 客户端解析JWT (❌ 不安全)
const payload = JSON.parse(atob(token.split('.')[1]));
const role = payload.role;

// 修复后: 后端API验证 (✅ 安全)
const response = await authApi.getMe();
const role = response.data.data.user.role;
```

**攻击防护**:
- ✅ 阻止JWT伪造攻击
- ✅ 防止权限绕过
- ✅ 确保token签名验证
- ✅ 检查token过期时间

---

### 1. 认证问题

#### 问题：JWT Token 无效
```bash
# 症状
401 Unauthorized: Invalid token

# 诊断步骤
# 1. 检查 Token 格式
echo "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" | base64 -d

# 2. 验证 Token
jwtdecode eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c

# 3. 检查密钥配置
grep JWT_SECRET .env
```

#### 问题：JWT密钥配置错误
```bash
# 症状
503 Service Unavailable: 认证服务配置错误,请联系管理员

# 原因
# 1. JWT_SECRET_KEY 环境变量未配置
# 2. JWT_SECRET_KEY 长度不足32字符

# 解决方案
# 1. 生成安全的密钥(至少32字符)
openssl rand -base64 32

# 2. 配置环境变量
export JWT_SECRET_KEY="生成的密钥"

# 3. 验证配置
env | grep JWT_SECRET_KEY
echo ${#JWT_SECRET_KEY}  # 应该 >= 32

# 4. 重启服务
docker-compose restart api
```

#### 问题：权限检查失败
```bash
# 症状
403 Forbidden: Permission denied

# 解决方案
# 1. 检查用户角色
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/v1/auth/me

# 2. 检查权限配置
SELECT * FROM user_roles WHERE user_id = 123;
SELECT * FROM role_permissions WHERE role_id = 1;
```

### 2. 数据泄露

#### 问题：敏感信息泄露
```bash
# 检查日志中的敏感信息
grep -i "password\|secret\|key" logs/*.log

# 清理敏感信息
sed -i 's/password=.*/password=****/' logs/*.log
```

---

## 📊 监控告警

### 1. Prometheus 告警

#### 常见告警处理
```yaml
# 高错误率告警
- alert: HighErrorRate
  expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
  for: 5m

  # 处理步骤
  # 1. 检查服务状态
  # 2. 查看错误日志
  # 3. 检查依赖服务

# 高响应时间告警
- alert: HighResponseTime
  expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
  for: 5m

  # 处理步骤
  # 1. 检查数据库性能
  # 2. 检查网络延迟
  # 3. 分析慢查询
```

### 2. Grafana 仪表盘

#### 关键指标监控
```bash
# API 响应时间
rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])

# 错误率
rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])

# 并发连接数
websocket_connections_current

# 数据库连接数
mysql_global_status_threads_connected
```

---

## 📝 日志分析

### 1. 应用日志

#### 日志级别和格式
```go
// 结构化日志示例
logger.WithFields(logrus.Fields{
    "user_id":    123,
    "order_id":   456,
    "action":     "create_order",
    "duration":   time.Since(start),
}).Info("Order created successfully")

// 错误日志
logger.WithFields(logrus.Fields{
    "error":      err.Error(),
    "stack_trace": debug.Stack(),
}).Error("Failed to create order")
```

#### 日志分析技巧
```bash
# 1. 统计错误类型
grep "ERROR" logs/app.log | awk '{print $4}' | sort | uniq -c

# 2. 查找慢请求
grep "duration" logs/app.log | awk '$4 > 1000'

# 3. 分析访问模式
grep "POST /api/v1/orders" logs/app.log | wc -l

# 4. 实时监控
tail -f logs/app.log | grep "ERROR"
```

### 2. 系统日志

#### 关键系统日志
```bash
# 1. 内核日志
dmesg | grep -i error

# 2. 系统日志
journalctl -u gamelink-api -f

# 3. Nginx 日志
tail -f /var/log/nginx/access.log
tail -f /var/log/nginx/error.log

# 4. Docker 日志
docker logs -f gamelink_api
```

---

## 🚨 常见错误代码

### HTTP 状态码
| 状态码 | 说明 | 解决方案 |
|--------|------|----------|
| 400 | 请求参数错误 | 检查请求参数格式和必填字段 |
| 401 | 未授权 | 检查 JWT Token 是否有效 |
| 403 | 权限不足 | 检查用户角色和权限配置 |
| 404 | 资源不存在 | 检查请求的 URL 和资源ID |
| 409 | 资源冲突 | 检查数据唯一性约束 |
| 422 | 参数验证失败 | 检查字段验证规则 |
| 429 | 请求过于频繁 | 降低请求频率或联系管理员 |
| 500 | 服务器错误 | 查看服务器日志 |

### 业务错误码
| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| USER_NOT_FOUND | 用户不存在 | 检查用户ID或重新注册 |
| INVALID_PASSWORD | 密码错误 | 重置密码 |
| PLAYER_NOT_VERIFIED | 陪玩师未认证 | 完成认证流程 |
| ORDER_STATUS_INVALID | 订单状态错误 | 检查订单当前状态 |
| INSUFFICIENT_BALANCE | 余额不足 | 充值或选择其他支付方式 |
| PAYMENT_FAILED | 支付失败 | 检查支付配置或重试 |
| FILE_TOO_LARGE | 文件过大 | 压缩文件或选择其他文件 |

---

## 🛠️ 调试工具

### 1. 后端调试

#### Delve 调试器
```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 启动调试
dlv debug ./cmd/user-service

# 断点调试
(dlv) break main.go:42
(dlv) continue
(dlv) print user
```

#### pprof 性能分析
```bash
# CPU 分析
curl http://localhost:8080/debug/pprof/profile > cpu.pprof
go tool pprof cpu.pprof

# 内存分析
curl http://localhost:8080/debug/pprof/heap > heap.pprof
go tool pprof heap.pprof

# 协程分析
curl http://localhost:8080/debug/pprof/goroutine > goroutine.pprof
go tool pprof goroutine.pprof
```

### 2. 前端调试

#### Chrome DevTools
```javascript
// 1. 网络请求调试
fetch('/api/v1/users')
  .then(res => res.json())
  .then(data => console.log(data));

// 2. 性能分析
performance.mark('start-operation');
// ... 执行操作
performance.mark('end-operation');
performance.measure('operation-duration', 'start-operation', 'end-operation');

// 3. 内存使用
console.log(performance.memory);
```

#### React DevTools
```javascript
// 组件调试
import { useEffect } from 'react';

function MyComponent() {
  useEffect(() => {
    // 在 React DevTools 中可见
    console.log('Component mounted');
  }, []);
}
```

### 3. 数据库调试

#### MySQL 调试
```sql
-- 1. 查看进程列表
SHOW FULL PROCESSLIST;

-- 2. 查看锁信息
SHOW ENGINE INNODB STATUS;

-- 3. 查看表状态
SHOW TABLE STATUS LIKE 'orders';

-- 4. 分析查询
EXPLAIN ANALYZE SELECT * FROM orders WHERE user_id = 123;
```

---

## 📞 联系支持

### 技术支持团队

| 问题类型 | 联系方式 | 响应时间 |
|----------|----------|----------|
| 紧急故障 | hotline@gamelink.com | 15分钟内 |
| 技术问题 | support@gamelink.com | 2小时内 |
| 功能咨询 | help@gamelink.com | 24小时内 |
| 商务合作 | business@gamelink.com | 48小时内 |

### 社区支持
- **GitHub Issues**: https://github.com/your-org/GameLink/issues
- **开发者论坛**: https://community.gamelink.com
- **知识库**: https://kb.gamelink.com

### 报告问题时请提供

1. **环境信息**
   - 操作系统版本
   - Go/Node.js 版本
   - Docker 版本
   - 浏览器版本

2. **问题描述**
   - 详细错误信息
   - 重现步骤
   - 期望结果

3. **日志信息**
   - 应用日志
   - 系统日志
   - 错误截图

4. **配置信息**
   - 环境变量（隐藏敏感信息）
   - 配置文件
   - 部署架构

---

## 📚 相关文档

- [开发指南](./DEVELOPMENT.md)
- [部署指南](./DEPLOYMENT.md)
- [API 文档](./API.md)
- [架构设计](./ARCHITECTURE.md)

---

*本文档持续更新中，最后更新: 2025-11-13*