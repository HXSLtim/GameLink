# JWT 认证快速测试指南

## ✅ JWT 认证已启用

环境变量 `ADMIN_AUTH_MODE=jwt` 已设置，服务正在 JWT 认证模式下运行。

---

## 🧪 测试方法

### 方法 1：浏览器测试页面（最简单）

1. 打开测试页面：
   ```
   file:///C:/Users/a2778/Desktop/code/GameLink/backend/test-jwt-auth.html
   ```

2. 点击 **Login** 按钮（已预填超级管理员凭证）

3. 登录成功后，点击 **Test Dashboard Access**

4. 查看结果 ✓

### 方法 2：使用你的前端应用

你的前端应用现在可以正常使用了！

**登录凭证**：
```json
{
  "email": "admin@gamelink.local",
  "password": "Admin@123456"
}
```

登录后，确保在所有管理接口请求中添加 Header：
```
Authorization: Bearer <你的token>
```

### 方法 3：Postman/Thunder Client

**Step 1: 登录**
```
POST http://localhost:8080/api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@gamelink.local",
  "password": "Admin@123456"
}
```

**Step 2: 复制返回的 `accessToken`**

**Step 3: 访问管理接口**
```
GET http://localhost:8080/api/v1/admin/stats/dashboard
Authorization: Bearer <你的token>
```

---

## 🔒 JWT 认证说明

### 当前配置
- **认证模式**: JWT (JSON Web Token)
- **Token 过期时间**: 3600 秒（1小时）
- **刷新机制**: 支持 RefreshToken

### 认证流程
1. 用户通过 `/api/v1/auth/login` 登录
2. 服务器返回 `accessToken` 和 `refreshToken`
3. 客户端在后续请求中携带 `Authorization: Bearer <accessToken>`
4. Token 过期后可使用 `refreshToken` 刷新

### 受保护的路由
所有 `/api/v1/admin/*` 路由都需要 JWT 认证：
- ✅ 用户管理
- ✅ 玩家管理
- ✅ 游戏管理
- ✅ 订单管理
- ✅ 统计分析
- ✅ 角色权限管理

---

## ⚠️ 注意事项

### 开发环境
- 当前配置：`ADMIN_AUTH_MODE=jwt`
- 如果需要禁用认证（仅开发），清除此环境变量

### 生产环境
- **必须**设置 `APP_ENV=production`
- JWT 认证自动启用，无需额外配置
- **必须**设置安全的 `JWT_SECRET`

---

## 🐛 常见问题

### Q: 提示"未授权"或 401 错误
A: 检查：
1. Token 是否过期（有效期 1 小时）
2. Authorization Header 格式：`Bearer <token>`（注意 Bearer 后有空格）
3. Token 是否完整复制（不要截断）

### Q: Token 过期后怎么办？
A: 使用 RefreshToken：
```
POST http://localhost:8080/api/v1/auth/refresh
Content-Type: application/json

{
  "refreshToken": "<你的refresh_token>"
}
```

### Q: 如何在开发中临时禁用认证？
A: 清除环境变量：
```powershell
Remove-Item Env:\ADMIN_AUTH_MODE
Remove-Item Env:\ADMIN_TOKEN
```
然后重启服务

---

## ✨ 测试结果验证

如果一切正常，你应该能看到：

**登录响应**：
```json
{
  "success": true,
  "code": 200,
  "message": "login successful",
  "data": {
    "accessToken": "eyJhbGciOiJIUz...",
    "refreshToken": "eyJhbGciOiJIUz...",
    "expiresIn": 3600,
    "user": {
      "id": 1,
      "name": "Super Admin",
      "email": "admin@gamelink.local",
      "role": "admin"
    }
  }
}
```

**Dashboard 响应**：
```json
{
  "success": true,
  "code": 200,
  "data": {
    "totalUsers": 16,
    "totalPlayers": 6,
    "totalOrders": 11,
    "totalRevenue": 4200.00,
    ...
  }
}
```

---

**服务状态**: 🟢 运行中  
**认证模式**: 🔐 JWT  
**端口**: 8080

现在你可以在前端正常使用了！

