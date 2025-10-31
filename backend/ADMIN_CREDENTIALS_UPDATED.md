# ✅ 管理员账户已更新

## 🔐 唯一的超级管理员账户

系统现在只初始化**一个**超级管理员账户：

```
邮箱：superAdmin@GameLink.com
密码：admin123
```

---

## 📝 更改内容

### 1. 迁移逻辑更新 (`internal/db/migrate.go`)

**修改前**：
```go
email = "admin@gamelink.local"
password = "Admin@123456"
```

**修改后**：
```go
email = "superAdmin@GameLink.com"
password = "admin123"
```

### 2. 种子数据更新 (`internal/db/seed.go`)

**移除了**：
```go
{Key: "adminA", Email: "admin@gamelink.com", Phone: "13800138005", Name: "系统管理员", Role: model.RoleAdmin, Password: "Admin@123456"}
```

**原因**：避免混淆，只保留一个超级管理员账户。

### 3. 文档更新 (`docs/super-admin.md`)

已更新默认凭证说明。

---

## 🚀 如何使用

### 步骤 1：删除旧数据库（如果存在）

```powershell
# 删除旧数据库文件
Remove-Item C:\Users\a2778\Desktop\code\GameLink\backend\var\dev.db -ErrorAction SilentlyContinue
```

### 步骤 2：重启服务

服务会自动创建新数据库并初始化超级管理员账户。

```powershell
cd C:\Users\a2778\Desktop\code\GameLink\backend
$env:ADMIN_AUTH_MODE="jwt"
go run .\cmd\user-service\main.go
```

### 步骤 3：登录

在前端或 Postman 中使用新凭证登录：

```json
{
  "email": "superAdmin@GameLink.com",
  "password": "admin123"
}
```

---

## 🧪 快速测试

运行此脚本验证新账户：

```powershell
cd C:\Users\a2778\Desktop\code\GameLink\backend
.\test-new-admin.ps1
```

---

## 📊 账户对比

| 版本 | 邮箱 | 密码 | 说明 |
|------|------|------|------|
| **旧版** | `admin@gamelink.local` | `Admin@123456` | 已移除 |
| **旧版** | `admin@gamelink.com` | `Admin@123456` | 已移除 |
| **新版** ✅ | `superAdmin@GameLink.com` | `admin123` | **唯一管理员** |

---

## ⚠️ 重要提示

### 开发环境
- ✅ 使用默认凭证：`superAdmin@GameLink.com` / `admin123`
- ✅ 确保设置了 `ADMIN_AUTH_MODE=jwt` 以启用 JWT 认证

### 生产环境
- ⚠️ **必须**设置自定义凭证：
  ```bash
  export APP_ENV=production
  export SUPER_ADMIN_EMAIL="your-admin@company.com"
  export SUPER_ADMIN_PASSWORD="YourStrongPassword123!@#"
  ```
- ⚠️ 不设置会导致服务启动失败（安全保护）

---

## 🐛 故障排除

### Q: 使用新凭证登录失败？
**A**: 可能是因为数据库中还保留着旧账户。解决方法：
1. 停止服务
2. 删除数据库文件：`var/dev.db`
3. 重启服务（会自动重建）

### Q: 提示 403 权限不足？
**A**: 确保：
1. 使用正确的邮箱：`superAdmin@GameLink.com`（区分大小写）
2. 已设置 `ADMIN_AUTH_MODE=jwt`
3. Token 未过期（有效期 1 小时）

### Q: 旧账户还能用吗？
**A**: 不能。旧账户 (`admin@gamelink.local` 和 `admin@gamelink.com`) 已从种子数据中移除。

---

## 🎯 测试检查清单

- [ ] 删除旧数据库
- [ ] 重启服务
- [ ] 使用新凭证登录
- [ ] 访问管理仪表盘
- [ ] 测试其他管理接口

全部通过后，系统就可以正常使用了！

---

**更新时间**: 2025-10-31  
**更改原因**: 统一管理员账户，避免混淆  
**影响范围**: 开发环境和测试环境

