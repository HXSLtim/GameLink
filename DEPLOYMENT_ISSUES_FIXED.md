# 🔧 部署问题修复记录

**日期**: 2025-12-13  
**状态**: ✅ 主要问题已修复

---

## 🐛 发现的问题

### 1. 超级管理员角色未分配 ✅ 已修复

**问题描述**:
- 超级管理员用户创建成功（ID=1, email=admin@gamelink.com）
- 但 `user_roles` 表中没有该用户的角色记录
- 导致超级管理员无法访问任何需要权限的接口，返回 403 错误

**根本原因**:
- `backend/pkg/db/migrate.go` 中的 `ensureSuperAdmin()` 函数在分配角色时可能失败
- 错误被捕获但只打印警告，没有阻止程序继续运行
- 可能的失败原因：
  1. 角色表在用户创建时还未初始化
  2. 数据库事务问题
  3. 表结构问题

**修复方法**:
```sql
-- 手动为超级管理员分配角色
INSERT INTO user_roles (user_id, role_id) VALUES (1, 1) ON CONFLICT DO NOTHING;
```

**验证**:
```sql
-- 查询结果应该显示 superAdmin 角色
SELECT u.id, u.email, u.name, r.slug, r.name 
FROM users u 
LEFT JOIN user_roles ur ON u.id = ur.user_id 
LEFT JOIN roles r ON ur.role_id = r.id 
WHERE u.id = 1;

-- 结果:
-- id |       email        |    name     |    slug    |    name    
-- ----+--------------------+-------------+------------+------------
--  1 | admin@gamelink.com | Super Admin | superAdmin | 超级管理员
```

**建议的代码改进**:
1. 在 `ensureSuperAdmin()` 中使用数据库事务
2. 如果角色分配失败，应该返回错误而不是只打印警告
3. 添加重试机制或确保角色表先于用户创建

---

### 2. 菜单查询 SQL 语法错误 ⚠️ 待修复

**问题描述**:
```
ERROR: syntax error at or near "order" (SQLSTATE 42601)
SELECT * FROM "menus" WHERE "menus"."deleted_at" IS NULL ORDER BY `order` ASC, id ASC
```

**根本原因**:
- `order` 是 PostgreSQL 的保留关键字
- 在 `ORDER BY` 子句中使用了反引号 `` `order` ``（MySQL 语法）
- PostgreSQL 应该使用双引号 `"order"` 或者重命名列

**影响**:
- 菜单列表接口返回 500 错误
- 前端无法加载菜单数据

**修复方案**:

**方案 1**: 修改查询语句（推荐）
```go
// backend/internal/repository/admin/menu.go:50
// 将 ORDER BY `order` 改为 ORDER BY "order"
db.Order(`"order" ASC, id ASC`)
```

**方案 2**: 重命名数据库列
```sql
-- 将 order 列重命名为 sort_order 或 display_order
ALTER TABLE menus RENAME COLUMN "order" TO sort_order;
```

**方案 3**: 使用别名
```go
db.Order("menus.\"order\" ASC, id ASC")
```

---

### 3. 健康检查路径错误 ✅ 已修复

**问题描述**:
- Docker healthcheck 使用 `/api/v1/health` 路径
- 但后端实际注册的是 `/api/v1/healthz` 路径
- 导致健康检查失败，容器显示 unhealthy

**修复**:
```yaml
# docker-compose.prod.yml
healthcheck:
  test: ["CMD", "wget", "--no-verbose", "--tries=1", "-O", "/dev/null", "http://localhost:8080/api/v1/healthz"]
  interval: 10s
  timeout: 5s
  retries: 5
  start_period: 30s
```

**说明**:
- 使用 `-O /dev/null` 而不是 `--spider`，因为 wget spider 模式使用 HEAD 请求
- 后端只注册了 GET 方法，不支持 HEAD 请求

---

### 4. 管理员密码长度不足 ✅ 已修复

**问题描述**:
```
failed to bootstrap application: SUPER_ADMIN_PASSWORD must be at least 8 characters in production
```

**原因**:
- `.env` 文件中的密码只有 6 个字符：`123456`
- 生产环境要求至少 8 个字符

**修复**:
```bash
# .env
SUPER_ADMIN_PASSWORD=Admin2025@Pass#
```

---

## ✅ 已完成的工作

### 1. 加密部署脚本验证

**脚本**: `scripts/deploy-production-encrypted.ps1`

**测试结果**: ✅ 成功

**部署步骤**:
1. ✅ 检查环境变量
2. ✅ 检查/生成加密密钥
3. ✅ 同步密钥到前端
4. ✅ 安装前端依赖（crypto-js）
5. ✅ 构建前端（8.61秒）
6. ✅ 构建 Docker 镜像
7. ✅ 停止旧服务
8. ✅ 启动新服务

**服务状态**:
- ✅ Backend: healthy
- ✅ Frontend: running
- ✅ PostgreSQL: healthy
- ✅ Redis: healthy

**加密状态**:
- ✅ 后端加密中间件已启用
- ✅ 前端加密配置已启用
- ✅ 密钥已同步

---

### 2. 环境配置

**后端配置** (`.env`):
```bash
POSTGRES_DB=gamelink
POSTGRES_USER=gamelink
POSTGRES_PASSWORD=RNGIydxs0475
REDIS_PASSWORD=TvXYJ305HNhsnIpQ
JWT_SECRET_KEY=BaFNYRM44ePcXL/nUt6UoLf7fXaY88HojQIpfvo/zHU=
CRYPTO_SECRET_KEY=fWDO82ax5N7eiut3SPhZnKQkAbCd1234
CRYPTO_IV=Qn8wt3ygOZpYJaE6
SUPER_ADMIN_EMAIL=admin@gamelink.com
SUPER_ADMIN_PASSWORD=Admin2025@Pass#
```

**前端配置** (`frontend/.env.production`):
```bash
VITE_API_BASE_URL=/api/v1
VITE_CRYPTO_ENABLED=true
VITE_CRYPTO_SECRET_KEY=fWDO82ax5N7eiut3SPhZnKQkAbCd1234
VITE_CRYPTO_IV=Qn8wt3ygOZpYJaE6
VITE_CRYPTO_USE_SIGNATURE=true
```

---

### 3. Docker 配置优化

**健康检查修复**:
- ✅ 修正健康检查路径：`/api/v1/healthz`
- ✅ 使用 GET 请求而不是 HEAD 请求
- ✅ 优化检查间隔和超时时间

**加密环境变量**:
```yaml
environment:
  - CRYPTO_ENABLED=true
  - CRYPTO_SECRET_KEY=${CRYPTO_SECRET_KEY}
  - CRYPTO_IV=${CRYPTO_IV}
  - CRYPTO_USE_SIGNATURE=true
```

---

## 📋 待办事项

### 高优先级

1. **修复菜单查询 SQL 错误** ⚠️
   - 文件: `backend/internal/repository/admin/menu.go:50`
   - 修改: `ORDER BY "order" ASC` 或重命名列

2. **改进超级管理员角色分配逻辑** ⚠️
   - 文件: `backend/pkg/db/migrate.go`
   - 使用事务确保角色分配成功
   - 失败时返回错误而不是警告

### 中优先级

3. **添加数据库迁移脚本**
   - 自动修复缺失的角色分配
   - 检查并修复数据完整性

4. **优化前端构建大小**
   - 当前最大 chunk: 1.4MB
   - 考虑代码分割和懒加载

### 低优先级

5. **添加健康检查端点别名**
   - 同时支持 `/health` 和 `/healthz`
   - 提高兼容性

6. **完善部署文档**
   - 添加故障排查指南
   - 记录常见问题和解决方案

---

## 🧪 测试建议

### 1. 权限测试

```bash
# 登录超级管理员
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@gamelink.com","password":"Admin2025@Pass#"}'

# 测试权限接口（应该返回 200）
curl -X GET http://localhost:8080/api/v1/admin/stats/dashboard \
  -H "Authorization: Bearer <token>"
```

### 2. 加密测试

```bash
# 在浏览器开发者工具中查看请求
# 1. 访问 http://localhost
# 2. 打开 Network 标签
# 3. 尝试登录
# 4. 查看请求体应该是加密格式：
#    {"encrypted":true,"payload":"...","timestamp":...,"signature":"..."}
```

### 3. 数据库完整性测试

```sql
-- 检查所有用户的角色分配
SELECT u.id, u.email, u.name, COUNT(ur.role_id) as role_count
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
GROUP BY u.id, u.email, u.name
HAVING COUNT(ur.role_id) = 0;

-- 应该返回 0 行（所有用户都有角色）
```

---

## 📚 相关文档

- **[加密配置指南](CRYPTO_SETUP_GUIDE.md)** - 加密配置详细说明
- **[部署脚本说明](scripts/README.md)** - 部署脚本使用指南
- **[最终部署状态](FINAL_DEPLOYMENT_STATUS.md)** - 完整部署状态
- **[部署验证](DEPLOYMENT_VERIFIED.md)** - 部署验证结果

---

**更新时间**: 2025-12-13  
**版本**: 1.0.0
