# Mock 登录使用指南

**更新日期**: 2025-10-28  
**功能**: Mock 模拟登录（无需后端）  
**状态**: ✅ 已启用

---

## 🎯 功能说明

现在登录功能已更新为 **Mock 模式**，可以使用**任意用户名和密码**直接登录进入系统，无需连接后端 API。

---

## 🚀 快速开始

### 1. 启动项目

```bash
npm run dev
```

访问: http://localhost:5173

### 2. 登录测试

在登录页面输入：

#### 👤 方案一：管理员账号

- **用户名**: `admin`
- **密码**: 任意（至少6个字符）

登录后用户角色显示为 `admin`

#### 👤 方案二：普通用户

- **用户名**: 任意（至少3个字符）
- **密码**: 任意（至少6个字符）

登录后用户角色显示为 `user`

#### 示例

```
用户名: test
密码: 123456
```

```
用户名: alice
密码: password
```

```
用户名: admin
密码: admin123
```

**✨ 所有账号都可以成功登录！**

---

## 📝 Mock 登录实现

### 代码位置

**文件**: `src/contexts/AuthContext.tsx`

### 核心功能

```typescript
// 登录方法 (Mock 实现)
const login = async (username: string, password: string) => {
  // 1. 模拟网络延迟 800ms
  await new Promise((resolve) => setTimeout(resolve, 800));

  // 2. 基础验证
  if (!username || !password) {
    throw new Error('用户名和密码不能为空');
  }

  // 3. 生成 Mock Token
  const mockToken = `mock-token-${Date.now()}-${randomString}`;

  // 4. 创建 Mock 用户信息
  const mockUser = {
    id: Math.floor(Math.random() * 1000),
    username: username,
    email: `${username}@gamelink.com`,
    role: username === 'admin' ? 'admin' : 'user',
    avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=${username}`,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  };

  // 5. 保存到 localStorage
  storage.setItem(STORAGE_KEYS.token, mockToken);
  storage.setItem(STORAGE_KEYS.user, mockUser);

  // 6. 更新状态
  setToken(mockToken);
  setUser(mockUser);
};
```

### Mock 用户数据结构

```typescript
{
  id: 123,                                    // 随机生成
  username: "admin",                          // 你输入的用户名
  email: "admin@gamelink.com",               // 自动生成
  role: "admin",                              // admin 用户名 = admin 角色
  avatar: "https://api.dicebear.com/...",    // 随机头像
  createdAt: "2025-10-28T10:30:00.000Z",
  updatedAt: "2025-10-28T10:30:00.000Z"
}
```

---

## 🔐 登录规则

### 验证规则

1. **用户名**: 必填，最少 3 个字符
2. **密码**: 必填，最少 6 个字符

### 角色分配

- 用户名 = `admin` → 角色 = `admin`
- 其他用户名 → 角色 = `user`

### Token 生成

```
mock-token-1698483000000-a1b2c3d4e
         └─时间戳─┘  └─随机字符串─┘
```

---

## 💾 数据持久化

### 存储位置

登录信息保存在浏览器的 `localStorage` 中：

| Key              | Value                 | 说明         |
| ---------------- | --------------------- | ------------ |
| `gamelink_token` | mock-token-xxx        | Mock Token   |
| `gamelink_user`  | { id, username, ... } | 用户信息对象 |

### 查看存储

打开浏览器开发者工具：

```
Application → Storage → Local Storage → http://localhost:5173
```

你会看到：

```
gamelink_token: "mock-token-1698483000000-a1b2c3d4e"
gamelink_user: {"id":123,"username":"admin",...}
```

### 自动恢复

刷新页面后，系统会自动从 localStorage 恢复登录状态，无需重新登录。

---

## 🔄 登录流程

### 完整流程

```
1. 用户访问 / 或 /dashboard
      ↓
2. ProtectedRoute 检查登录状态
      ↓
3. 未登录 → 重定向到 /login
      ↓
4. 用户输入用户名和密码
      ↓
5. 点击"登录"按钮
      ↓
6. AuthContext.login(username, password)
      ↓
7. 验证通过 → 生成 Mock 数据
      ↓
8. 保存到 localStorage
      ↓
9. 更新 React Context 状态
      ↓
10. navigate('/dashboard')
      ↓
11. 显示仪表盘（带 Header + Sidebar）
```

### 页面刷新流程

```
1. 页面刷新
      ↓
2. AuthContext 初始化
      ↓
3. 从 localStorage 读取 token 和 user
      ↓
4. 如果存在 → 恢复登录状态
      ↓
5. 如果不存在 → 保持未登录
```

---

## 🚪 退出登录

### 操作步骤

1. 点击右上角用户名
2. 下拉菜单中点击"退出登录"
3. 清除 localStorage 数据
4. 跳转到登录页

### 实现代码

```typescript
const logout = () => {
  // 清除存储
  storage.removeItem(STORAGE_KEYS.token);
  storage.removeItem(STORAGE_KEYS.user);

  // 清除状态
  setToken(null);
  setUser(null);

  console.log('✅ 退出登录成功');
};
```

---

## 🎨 用户界面显示

### Header 显示

登录后，Header 右上角显示：

```
┌──────────────────┐
│ 👤 admin    ▼   │  ← 点击展开下拉菜单
└──────────────────┘
```

展开后：

```
┌──────────────────┐
│ admin            │
│ admin            │  ← 用户名
├──────────────────┤
│ 🚪 退出登录      │  ← 点击退出
└──────────────────┘
```

### Sidebar 菜单

根据角色显示不同菜单（当前所有角色看到相同菜单）：

- 🏠 仪表盘
- 👥 用户管理
- 📦 订单管理
- ⚙️ 系统设置

---

## 🔍 调试信息

### Console 日志

登录成功后，控制台会输出：

```javascript
✅ Mock 登录成功: {
  username: "admin",
  token: "mock-token-1698483000000-a1b2c3d4e"
}
```

退出登录后：

```javascript
✅ 退出登录成功
```

### 网络请求

Mock 模式下**不会发送任何网络请求**，所有操作都在前端完成。

---

## 📊 与真实登录的区别

| 特性       | Mock 登录       | 真实登录               |
| ---------- | --------------- | ---------------------- |
| 后端 API   | ❌ 不需要       | ✅ 需要                |
| 网络请求   | ❌ 无           | ✅ 有                  |
| Token 验证 | ❌ 不验证       | ✅ 验证                |
| 用户数据   | 🔧 前端生成     | 📡 后端返回            |
| 角色权限   | 🔧 简单判断     | 📡 后端控制            |
| 数据持久化 | 💾 localStorage | 💾 localStorage + 后端 |

---

## 🔄 切换到真实登录

### 需要的步骤

1. **后端准备**
   - 实现登录 API: `POST /api/v1/auth/login`
   - 实现获取用户信息 API: `GET /api/v1/auth/me`

2. **前端修改**

```typescript
// src/contexts/AuthContext.tsx

// 替换 Mock 实现
const login = async (username: string, password: string) => {
  setLoginLoading(true);

  try {
    // 调用真实 API
    const response = await fetch('/api/v1/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    });

    if (!response.ok) {
      throw new Error('登录失败');
    }

    const { token } = await response.json();

    // 保存 token
    storage.setItem(STORAGE_KEYS.token, token);
    setToken(token);

    // 获取用户信息
    const userResponse = await fetch('/api/v1/auth/me', {
      headers: { Authorization: `Bearer ${token}` },
    });

    const user = await userResponse.json();
    storage.setItem(STORAGE_KEYS.user, user);
    setUser(user);
  } finally {
    setLoginLoading(false);
  }
};
```

---

## 💡 常见问题

### Q1: 刷新页面后还是登录状态吗？

**A**: ✅ 是的。登录信息保存在 localStorage 中，刷新页面会自动恢复。

### Q2: 可以同时多个账号登录吗？

**A**: ❌ 不可以。新登录会覆盖旧的登录信息（localStorage 只保存一份）。

### Q3: 密码会被验证吗？

**A**: ❌ Mock 模式下不会验证密码，任何密码都可以登录（只要长度 ≥ 6）。

### Q4: 如何清除登录状态？

**A**:

- 方法1: 点击"退出登录"
- 方法2: 清除浏览器 localStorage
- 方法3: 打开开发者工具执行: `localStorage.clear()`

### Q5: Token 会过期吗？

**A**: ❌ Mock Token 不会过期，永久有效（除非手动清除）。

### Q6: 可以自定义用户信息吗？

**A**: ✅ 可以。修改 `AuthContext.tsx` 中的 `mockUser` 对象。

---

## 🎯 最佳实践

### 开发阶段

```typescript
// ✅ 使用 Mock 登录
// 优点：
// - 不依赖后端
// - 快速开发
// - 随时可登录

// 建议的测试账号：
admin / 123456; // 管理员
user / 123456; // 普通用户
test / 123456; // 测试用户
```

### 测试阶段

```typescript
// 测试不同角色
admin → 测试管理员功能
user  → 测试普通用户功能

// 测试长用户名
verylongusername123 / 123456
```

### 生产环境

```typescript
// ❌ 不要使用 Mock 登录
// ✅ 切换到真实登录 API
```

---

## 📚 相关文档

- [导航系统文档](./NAVIGATION_SYSTEM.md)
- [设计系统文档](./DESIGN_SYSTEM_V2.md)
- [路由守卫文档](./NAVIGATION_SYSTEM.md#4-路由守卫)

---

## ✅ 测试清单

- [x] 任意用户名密码可登录
- [x] admin 用户名显示 admin 角色
- [x] 其他用户名显示 user 角色
- [x] 登录后跳转到 dashboard
- [x] Header 显示用户信息
- [x] 刷新页面保持登录状态
- [x] 退出登录功能正常
- [x] localStorage 正确存储数据
- [x] 控制台输出正确日志

---

**更新者**: GameLink Frontend Team  
**模式**: Mock 登录（开发模式）  
**最后更新**: 2025-10-28

---

<div align="center">

## 🎉 现在可以用任意账号登录了！

**No Backend Required**

⚫⚪ **简单 · 快速 · 开箱即用**

</div>
