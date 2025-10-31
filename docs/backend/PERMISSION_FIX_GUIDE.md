# 权限不足问题修复指南

## 🔍 问题诊断

你收到 **403 权限不足** 错误的原因是：

**你登录的用户不是超级管理员**

- 你登录的账户：`admin@gamelink.com`（user_id=7）
- 超级管理员账户：`admin@gamelink.local`（user_id=1）

注意区别：`.com` vs `.local`

---

## ✅ 解决方案（选择一种）

### 方案 1：使用正确的超级管理员账户（最简单）

**重新登录使用这个账户**：

```json
{
  "email": "admin@gamelink.local",
  "password": "Admin@123456"
}
```

这是真正的超级管理员账户，拥有所有权限。

---

### 方案 2：升级当前账户为超级管理员

如果你想继续使用 `admin@gamelink.com`，需要给它分配超级管理员角色。

#### 步骤 1：找到数据库文件

数据库通常在以下位置之一：
- `backend/gamelink.db`
- `backend/data/gamelink.db`
- 配置文件中指定的位置

#### 步骤 2：执行 SQL

使用 SQLite 工具或在线工具执行：

```sql
-- 给 admin@gamelink.com 分配 super_admin 角色
INSERT OR IGNORE INTO user_roles (user_id, role_id)
SELECT 
    (SELECT id FROM users WHERE email = 'admin@gamelink.com') as user_id,
    (SELECT id FROM roles WHERE slug = 'super_admin') as role_id;
```

#### 步骤 3：重启服务并重新登录

---

### 方案 3：使用自动修复脚本

运行提供的修复脚本：

```powershell
cd C:\Users\a2778\Desktop\code\GameLink\backend
.\fix-admin-roles.ps1
```

这会自动给所有 `role='admin'` 的用户分配 `super_admin` 角色。

---

## 🔒 RBAC 权限系统说明

### 双重角色系统

系统目前有两套角色机制：

1. **旧系统**（`users.role` 字段）：
   - 直接存储在用户表的 `role` 字段
   - 值：`admin`, `user`, `player`
   - **不支持细粎度权限控制**

2. **新 RBAC 系统**（`user_roles` 表）：
   - 通过 `user_roles` 表关联用户和角色
   - 支持多角色、细粒度权限
   - 角色：`super_admin`, `admin`, `user`, `player`
   - **这是权限中间件使用的系统**

### 权限检查逻辑

```
请求 /api/v1/admin/* 
  ↓
1. JWT 认证（验证 Token）✓
  ↓
2. 检查用户是否为 super_admin（查询 user_roles 表）
  ↓
  如果是 → 放行 ✓
  如果不是 → 继续检查
  ↓
3. 检查用户是否有特定权限（查询 role_permissions 表）
  ↓
  有权限 → 放行 ✓
  无权限 → 返回 403 ✗
```

### 为什么会出现 403

你的账户（`admin@gamelink.com`）：
- ✅ 有旧的 `role='admin'` 字段
- ❌ **没有 RBAC 的 `super_admin` 角色**
- ❌ **没有具体的页面权限**

所以被权限中间件拦截了。

---

## 🧪 验证修复

修复后，测试是否能访问管理接口：

```bash
# 登录
POST http://localhost:8080/api/v1/auth/login
{
  "email": "admin@gamelink.local",  // 或 admin@gamelink.com（如果已修复）
  "password": "Admin@123456"
}

# 获取 Token 后测试
GET http://localhost:8080/api/v1/admin/stats/dashboard
Authorization: Bearer <your_token>
```

**期望响应**：
```json
{
  "success": true,
  "code": 200,
  "data": {
    "totalUsers": 16,
    "totalPlayers": 6,
    ...
  }
}
```

---

## 📋 所有管理员账户

系统中的管理员账户：

| 邮箱 | 密码 | User ID | RBAC 角色 | 说明 |
|------|------|---------|-----------|------|
| `admin@gamelink.local` | `Admin@123456` | 1 | ✅ super_admin | **真正的超级管理员** |
| `admin@gamelink.com` | `password123` | 7 | ❌ 无 | 种子数据，需要手动分配 |

---

## 🐛 长期修复建议

修改种子数据逻辑，确保所有 `role='admin'` 的用户自动获得 RBAC 的 `super_admin` 角色：

```go
// 在 internal/db/seed.go 中添加
func assignAdminRoles(db *gorm.DB) error {
    var adminUsers []model.User
    if err := db.Where("role = ?", "admin").Find(&adminUsers).Error; err != nil {
        return err
    }
    
    var superAdminRole model.RoleModel
    if err := db.Where("slug = ?", "super_admin").First(&superAdminRole).Error; err != nil {
        return err
    }
    
    for _, user := range adminUsers {
        var exists model.UserRole
        err := db.Where("user_id = ? AND role_id = ?", user.ID, superAdminRole.ID).
            First(&exists).Error
        
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // 分配角色
            db.Create(&model.UserRole{
                UserID: user.ID,
                RoleID: superAdminRole.ID,
            })
        }
    }
    
    return nil
}
```

---

## ✨ 总结

**最快解决方法**：使用 `admin@gamelink.local` / `Admin@123456` 登录

**如果你已经在使用 `admin@gamelink.com`**：运行 SQL 或修复脚本给它分配 `super_admin` 角色

修复后，所有管理接口都能正常访问了！

