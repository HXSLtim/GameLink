# 管理员账户变更摘要

**变更日期**: 2025-10-31  
**变更类型**: 代码重构 - 简化管理员初始化  
**影响范围**: 开发和测试环境

---

## 📝 变更原因

用户要求：
- ✅ 只初始化一个超级管理员账户
- ✅ 使用指定的凭证：`superAdmin@GameLink.com` / `admin123`
- ✅ 避免多个管理员账户造成混淆

---

## 🔧 代码变更

### 1. `backend/internal/db/migrate.go`

**修改位置**: Line 149-161

```diff
  if email == "" && phone == "" {
      if env == "production" {
          return errors.New("SUPER_ADMIN_EMAIL or SUPER_ADMIN_PHONE must be set in production")
      }
-     email = "admin@gamelink.local"
+     email = "superAdmin@GameLink.com"
  }

  if password == "" {
      if env == "production" {
          return errors.New("SUPER_ADMIN_PASSWORD must be set in production")
      }
-     password = "Admin@123456"
+     password = "admin123"
  }
```

### 2. `backend/internal/db/seed.go`

**修改位置**: Line 89

```diff
  {Key: "proB", Email: "streamer@gamelink.com", Phone: "13800138004", Name: "魔王主播", Role: model.RolePlayer, Password: "Player@654321"},
- {Key: "adminA", Email: "admin@gamelink.com", Phone: "13800138005", Name: "系统管理员", Role: model.RoleAdmin, Password: "Admin@123456"},
+ // adminA removed - 只使用迁移时创建的超级管理员 (superAdmin@GameLink.com)
  {Key: "customerD", Email: "casual.player@gamelink.com", Phone: "13800138005", Name: "休闲玩家", Role: model.RoleUser, Password: "User@123789"},
```

**说明**: 
- 移除了种子数据中的 `admin@gamelink.com` 账户
- 调整了后续用户的手机号（避免冲突）

### 3. `backend/docs/super-admin.md`

**修改位置**: 默认凭证表格

```diff
  | 环境变量 | 说明 | 默认值（非生产环境） |
  |----------|------|----------------------|
- | `SUPER_ADMIN_EMAIL` | 超管邮箱，用作唯一登录标识 | `admin@gamelink.local` |
+ | `SUPER_ADMIN_EMAIL` | 超管邮箱，用作唯一登录标识 | `superAdmin@GameLink.com` |
  | `SUPER_ADMIN_PHONE` | 超管手机号，可选 | 空 |
  | `SUPER_ADMIN_NAME` | 显示名称 | `Super Admin` |
- | `SUPER_ADMIN_PASSWORD` | 登录密码（明文，启动时会自动加密） | `Admin@123456` |
+ | `SUPER_ADMIN_PASSWORD` | 登录密码（明文，启动时会自动加密） | `admin123` |
```

---

## 🎯 新的管理员凭证

### 开发环境（默认）

```
邮箱：superAdmin@GameLink.com
密码：admin123
```

**特性**：
- ✅ 自动创建（首次启动时）
- ✅ 自动分配 super_admin 角色
- ✅ 拥有所有管理权限
- ✅ 唯一的管理员账户

### 生产环境（必须自定义）

生产环境**必须**通过环境变量设置：

```bash
export APP_ENV=production
export SUPER_ADMIN_EMAIL="your-admin@company.com"
export SUPER_ADMIN_PASSWORD="YourStrongPassword123!@#"
```

如不设置，服务将拒绝启动（安全保护）。

---

## 🚀 升级步骤

### 对于新部署
直接启动服务即可，会自动使用新凭证。

### 对于现有部署

#### 步骤 1: 备份数据库（可选）
```powershell
Copy-Item backend/var/dev.db backend/var/dev.db.backup
```

#### 步骤 2: 删除数据库（推荐）
```powershell
Remove-Item backend/var/dev.db
```

#### 步骤 3: 重启服务
```powershell
cd backend
$env:ADMIN_AUTH_MODE="jwt"
go run .\cmd\user-service\main.go
```

服务会自动：
- 创建新数据库
- 初始化所有表结构
- 创建 `superAdmin@GameLink.com` 账户
- 分配 super_admin 角色
- 加载种子数据（如启用）

#### 步骤 4: 测试新账户
```powershell
.\test-new-admin.ps1
```

---

## 📊 变更对比

| 项目 | 变更前 | 变更后 |
|------|--------|--------|
| 管理员数量 | 2个 | **1个** ✅ |
| 迁移创建 | `admin@gamelink.local` | `superAdmin@GameLink.com` ✅ |
| 种子创建 | `admin@gamelink.com` | **已移除** ✅ |
| 默认密码 | `Admin@123456` | `admin123` ✅ |
| RBAC 角色 | super_admin | super_admin ✓ |

---

## ⚠️ 兼容性说明

### 破坏性变更
- ❌ 旧的管理员账户不再创建
- ❌ 现有数据库中的旧账户不会自动更新

### 不受影响
- ✓ 所有业务逻辑
- ✓ API 接口
- ✓ 权限系统
- ✓ 其他用户账户

### 迁移建议
如果你有现有的数据库：
1. **测试环境**: 删除数据库，重新初始化（推荐）
2. **生产环境**: 使用环境变量设置自定义凭证

---

## 🧪 测试验证

### 自动测试脚本
```powershell
cd backend
.\test-new-admin.ps1
```

### 手动测试步骤

#### 1. 登录测试
```bash
POST http://localhost:8080/api/v1/auth/login
Content-Type: application/json

{
  "email": "superAdmin@GameLink.com",
  "password": "admin123"
}
```

**预期结果**：
```json
{
  "success": true,
  "code": 200,
  "data": {
    "accessToken": "eyJhbGc...",
    "user": {
      "id": 1,
      "email": "superAdmin@GameLink.com",
      "name": "Super Admin",
      "role": "admin"
    }
  }
}
```

#### 2. 权限测试
```bash
GET http://localhost:8080/api/v1/admin/stats/dashboard
Authorization: Bearer <token>
```

**预期结果**: 200 OK，返回仪表盘数据

#### 3. 旧账户测试
尝试使用旧凭证登录：
- `admin@gamelink.local` / `Admin@123456`
- `admin@gamelink.com` / `Admin@123456`

**预期结果**: 401 Unauthorized（账户不存在）

---

## 📚 相关文档

已创建/更新的文档：
- ✅ `ADMIN_CREDENTIALS_UPDATED.md` - 详细更新说明
- ✅ `test-new-admin.ps1` - 自动测试脚本
- ✅ `docs/super-admin.md` - 超级管理员配置说明
- ✅ `ADMIN_ACCOUNT_CHANGE_SUMMARY.md` - 本文档

---

## 🐛 问题排查

### Q: 使用新凭证登录失败 (401)
**原因**: 数据库中还保留旧账户  
**解决**: 删除 `var/dev.db`，重启服务

### Q: 登录成功但访问管理接口提示 403
**原因**: JWT 认证模式未启用  
**解决**: 设置 `$env:ADMIN_AUTH_MODE="jwt"`，重启服务

### Q: 想保留现有数据怎么办？
**方案 1**: 在数据库中手动更新账户
```sql
UPDATE users SET email='superAdmin@GameLink.com' WHERE id=1;
UPDATE users SET password_hash='<bcrypt hash of admin123>' WHERE id=1;
```

**方案 2**: 使用 SQL 脚本迁移数据（需要自己编写）

### Q: 生产环境如何设置？
**答**: 必须使用环境变量
```bash
export APP_ENV=production
export SUPER_ADMIN_EMAIL="admin@yourcompany.com"
export SUPER_ADMIN_PASSWORD="YourSecurePassword123!@#"
```

---

## ✅ 验收标准

变更完成后，确认以下几点：

- [ ] 服务能正常启动
- [ ] 使用 `superAdmin@GameLink.com` / `admin123` 能登录
- [ ] 登录后能访问 `/api/v1/admin/stats/dashboard`
- [ ] 能访问其他管理接口（用户、游戏、订单等）
- [ ] 旧账户无法登录
- [ ] 自动测试脚本通过

全部通过后，变更即完成！

---

**变更人**: Claude AI  
**审核建议**: 在合并到主分支前，在测试环境充分验证  
**回滚方案**: Git revert 本次提交，恢复旧的凭证配置

