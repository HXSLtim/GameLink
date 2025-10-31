# ✅ 正确的超级管理员凭证

## 🔐 请使用此账户登录

```
邮箱：admin@gamelink.local
密码：Admin@123456
```

**注意事项**：
- ✅ 邮箱后缀是 `.local`（不是 `.com`）
- ✅ 密码首字母大写：`Admin@123456`（不是 `admin123`）

---

## 📊 所有管理员账户对比

| 邮箱 | 密码 | User ID | RBAC 角色 | 能否访问管理功能 |
|------|------|---------|-----------|------------------|
| `admin@gamelink.local` | `Admin@123456` | 1 | ✅ super_admin | ✅ 是 |
| `admin@gamelink.com` | `Admin@123456` | 7 | ❌ 无 | ❌ 否（403错误）|

---

## 🧪 快速测试

运行此脚本验证：

```powershell
cd C:\Users\a2778\Desktop\code\GameLink\backend

# 测试登录
$body = '{"email":"admin@gamelink.local","password":"Admin@123456"}'
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method POST -Body ([System.Text.Encoding]::UTF8.GetBytes($body)) -ContentType "application/json; charset=utf-8"

Write-Host "Login successful!" -ForegroundColor Green
Write-Host "User: $($response.data.user.name)" -ForegroundColor Cyan
Write-Host "Token obtained!" -ForegroundColor Green

# 测试管理接口
$headers = @{"Authorization" = "Bearer $($response.data.accessToken)"}
$dashboard = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/admin/stats/dashboard" -Headers $headers

Write-Host "Dashboard access successful!" -ForegroundColor Green
Write-Host "Total Users: $($dashboard.data.totalUsers)" -ForegroundColor Cyan
```

---

## 💡 为什么会有两个管理员？

1. **`admin@gamelink.local`**（User ID 1）
   - 在 `migrate.go` 中自动创建
   - 系统初始化时的第一个用户
   - 自动分配 RBAC 的 `super_admin` 角色
   - **这是真正的超级管理员**

2. **`admin@gamelink.com`**（User ID 7）
   - 在 `seed.go` 中作为测试数据创建
   - 只有旧的 `role='admin'` 字段
   - 没有 RBAC 系统的角色分配
   - **需要手动分配权限才能使用**

---

## 🔧 如果必须使用 admin@gamelink.com

如果你想使用 `admin@gamelink.com` 账户，需要给它分配 super_admin 角色。

在数据库中执行：

```sql
INSERT OR IGNORE INTO user_roles (user_id, role_id)
VALUES (
    (SELECT id FROM users WHERE email = 'admin@gamelink.com'),
    (SELECT id FROM roles WHERE slug = 'super_admin')
);
```

然后重启服务并重新登录。

---

## ✅ 推荐做法

**直接使用 `admin@gamelink.local` 账户**，它已经配置好了所有权限！

```
Email: admin@gamelink.local
Password: Admin@123456
```

登录后，所有 403 错误都会消失。

