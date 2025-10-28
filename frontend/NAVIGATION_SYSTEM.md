# GameLink 导航系统实现文档

**实现日期**: 2025-10-28  
**设计风格**: Neo-brutalism 黑白极简  
**状态**: ✅ 完成

---

## 📋 实现内容

### 1. 布局系统 ✅

#### Header（顶部导航栏）

**文件**: `src/components/Layout/Header.tsx`

**功能**:

- ⚫⚪ 黑白品牌 Logo（GameLink）
- 🍔 汉堡菜单按钮（切换侧边栏）
- 👤 用户信息显示
- 📋 用户下拉菜单
- 🚪 退出登录功能

**设计特点**:

- 2px 黑色边框
- 固定在顶部（sticky）
- 高度 64px
- 悬停动画效果

**代码示例**:

```tsx
<Header
  user={{ username: 'admin', role: 'admin' }}
  onLogout={handleLogout}
  onToggleSidebar={handleToggle}
/>
```

#### Sidebar（侧边栏）

**文件**: `src/components/Layout/Sidebar.tsx`

**功能**:

- 📱 可折叠设计
- 🎯 高亮当前路由
- 🔗 路由链接集成
- 📱 响应式适配

**菜单样式**:

- 2px 黑色边框
- 实体阴影效果
- 激活状态：黑底白字
- 悬停动画

**代码示例**:

```tsx
const menuItems = [
  {
    key: 'dashboard',
    label: '仪表盘',
    icon: <DashboardIcon />,
    path: '/dashboard',
  },
];

<Sidebar menuItems={menuItems} collapsed={false} />;
```

#### Layout（主布局）

**文件**: `src/components/Layout/Layout.tsx`

**功能**:

- 🎨 组合 Header + Sidebar + Content
- 📐 响应式布局
- 🎯 统一管理

**布局结构**:

```
┌─────────────────────────────────┐
│          Header (64px)          │
├──────────┬──────────────────────┤
│          │                      │
│ Sidebar  │      Content         │
│ (240px)  │      (flex: 1)       │
│          │                      │
└──────────┴──────────────────────┘
```

### 2. 面包屑组件 ✅

**文件**: `src/components/Breadcrumb/Breadcrumb.tsx`

**功能**:

- 🍞 显示当前页面路径
- 🔗 支持点击跳转
- ➡️ 自定义分隔符
- 🎨 黑白风格设计

**使用示例**:

```tsx
const items = [
  { label: '首页', path: '/' },
  { label: '用户管理', path: '/users' },
  { label: '用户详情' }, // 当前页，不可点击
];

<Breadcrumb items={items} />;
```

**样式特点**:

- 当前页面：黑色边框高亮
- 链接悬停：下划线动画
- 箭头分隔符

### 3. 路由系统 ✅

**文件**: `src/router/index.tsx`

**路由配置**:

```tsx
{
  '/login': 登录页（无保护）,
  '/': {
    element: ProtectedRoute,  // 路由守卫
    children: [
      '/dashboard': 仪表盘,
      '/users': 用户管理,
      '/orders': 订单管理,
      '/settings': 系统设置,
    ]
  }
}
```

**特性**:

- ✅ 嵌套路由
- ✅ 路由守卫
- ✅ 自动重定向
- ✅ 404 处理

### 4. 路由守卫 ✅

**文件**: `src/router/ProtectedRoute.tsx`

**功能**:

- 🔒 检查用户登录状态
- 🔄 未登录自动跳转登录页
- ⏳ 加载状态显示
- 🎯 基于 AuthContext

**实现逻辑**:

```tsx
const ProtectedRoute = () => {
  const { user, loading } = useAuth();

  if (loading) return <Loading />;
  if (!user) return <Navigate to="/login" />;
  return <Outlet />;
};
```

### 5. Dashboard 仪表盘 ✅

**文件**: `src/pages/Dashboard/Dashboard.tsx`

**功能**:

- 📊 统计卡片展示
- 📈 数据趋势显示
- 🎨 黑白卡片设计
- 🍞 面包屑导航

**统计模块**:

- 👥 总用户数
- 📦 总订单数
- 💰 总收入
- 📱 活跃用户

**设计特点**:

- 网格布局（Grid）
- 实体阴影卡片
- 悬停动画效果
- 图标 + 数据展示

---

## 🎨 设计风格

### 全局特点

1. **纯黑白配色**
   - 主色：#000000（纯黑）
   - 背景：#ffffff（纯白）
   - 辅助：灰度系列

2. **Neo-brutalism 元素**
   - 2px 粗边框
   - 实体阴影（8px × 8px）
   - 无圆角设计
   - 直线条风格

3. **交互动画**
   - 悬停上浮效果
   - 阴影增大动画
   - 平滑过渡（200ms）

---

## 📂 文件结构

```
src/
├── components/
│   ├── Layout/
│   │   ├── Header.tsx              # 顶部导航栏
│   │   ├── Header.module.less
│   │   ├── Sidebar.tsx             # 侧边栏
│   │   ├── Sidebar.module.less
│   │   ├── Layout.tsx              # 主布局
│   │   ├── Layout.module.less
│   │   └── index.ts
│   ├── Breadcrumb/
│   │   ├── Breadcrumb.tsx          # 面包屑
│   │   ├── Breadcrumb.module.less
│   │   └── index.ts
│   └── index.ts                    # 组件统一导出
├── router/
│   ├── index.tsx                   # 路由配置
│   ├── ProtectedRoute.tsx          # 路由守卫
│   └── layouts/
│       ├── MainLayout.tsx          # 主布局配置
│       └── index.ts
├── pages/
│   ├── Login/                      # 登录页
│   │   ├── Login.tsx
│   │   ├── Login.module.less
│   │   └── index.ts
│   └── Dashboard/                  # 仪表盘
│       ├── Dashboard.tsx
│       ├── Dashboard.module.less
│       └── index.ts
└── App.tsx                         # 路由Provider
```

---

## 🚀 使用指南

### 1. 访问登录页

```
http://localhost:5173/login
```

**测试账号**:

- 用户名: `admin`
- 密码: `admin123`

### 2. 登录后跳转

自动跳转到仪表盘：

```
http://localhost:5173/dashboard
```

### 3. 侧边栏菜单

- 🏠 仪表盘 `/dashboard`
- 👥 用户管理 `/users`（占位）
- 📦 订单管理 `/orders`（占位）
- ⚙️ 系统设置 `/settings`（占位）

### 4. 添加新路由

#### Step 1: 创建页面组件

```tsx
// src/pages/Users/Users.tsx
export const Users = () => {
  return (
    <div>
      <Breadcrumb items={[{ label: '首页', path: '/' }, { label: '用户管理' }]} />
      <h1>用户管理</h1>
    </div>
  );
};
```

#### Step 2: 添加路由

```tsx
// src/router/index.tsx
{
  path: 'users',
  element: <Users />,
}
```

#### Step 3: 添加菜单项

```tsx
// src/router/layouts/MainLayout.tsx
const menuItems: MenuItem[] = [
  // ...existing items
  {
    key: 'users',
    label: '用户管理',
    icon: <UsersIcon />,
    path: '/users',
  },
];
```

---

## 🎯 核心特性

### 1. 响应式设计

**桌面端** (> 768px):

- Sidebar: 240px 宽度
- 可折叠为 64px

**移动端** (<= 768px):

- Sidebar: 固定定位
- 折叠时隐藏（translateX(-100%)）
- 显示时覆盖内容区

### 2. 状态管理

使用 `AuthContext` 管理：

- ✅ 用户登录状态
- ✅ 用户信息
- ✅ 登录/退出方法

### 3. 路由保护

**保护规则**:

1. 未登录 → 重定向到 `/login`
2. 已登录访问 `/login` → 保持（可优化）
3. 404 路由 → 重定向到 `/dashboard`

### 4. 动画效果

**Header**:

- 用户菜单下拉动画（slideDown 0.2s）
- 按钮悬停背景色变化

**Sidebar**:

- 菜单项悬停：阴影 + 上浮
- 激活状态：黑底白字 + 左侧白色标记
- 宽度过渡动画（300ms）

**Breadcrumb**:

- 链接悬停：下划线展开
- 当前页：边框高亮

**Dashboard**:

- 卡片悬停：阴影增大 + 上浮
- 统计数字：大字号显示
- 趋势指标：上升↑/下降↓

---

## 🔧 技术实现

### 依赖包

```json
{
  "react-router-dom": "^6.27.0", // 路由管理
  "react": "^18.3.0", // React 核心
  "typescript": "^5.6.0" // TypeScript
}
```

### 核心 API

**React Router**:

- `createBrowserRouter` - 创建路由
- `RouterProvider` - 路由提供者
- `Outlet` - 渲染子路由
- `Navigate` - 重定向
- `Link` - 路由链接
- `useNavigate` - 编程式导航
- `useLocation` - 获取当前路由

**AuthContext**:

- `useAuth()` - 获取认证状态
- `login()` - 登录方法
- `logout()` - 退出登录

---

## 📱 用户流程

### 登录流程

```
访问 / 或 /dashboard
     ↓
检查登录状态（ProtectedRoute）
     ↓
  未登录 ───→ 重定向到 /login
     ↓
显示登录页面
     ↓
输入用户名/密码
     ↓
点击登录按钮
     ↓
调用 AuthContext.login()
     ↓
登录成功 → navigate('/dashboard')
     ↓
显示仪表盘页面（带 Layout）
```

### 退出流程

```
点击用户名 → 显示下拉菜单
     ↓
点击"退出登录"
     ↓
调用 AuthContext.logout()
     ↓
清除用户状态
     ↓
navigate('/login')
     ↓
返回登录页面
```

---

## 💡 最佳实践

### 1. 组件复用

```tsx
// ✅ 好的实践 - 使用通用 Layout
<MainLayout>
  <Dashboard />
</MainLayout>

// ❌ 避免 - 每个页面重复 Header/Sidebar
<div>
  <Header />
  <Sidebar />
  <Dashboard />
</div>
```

### 2. 路由守卫

```tsx
// ✅ 好的实践 - 统一守卫
<Route element={<ProtectedRoute />}>
  <Route path="/dashboard" />
  <Route path="/users" />
</Route>;

// ❌ 避免 - 每个组件单独检查
const Dashboard = () => {
  const { user } = useAuth();
  if (!user) return <Navigate to="/login" />;
  // ...
};
```

### 3. 面包屑配置

```tsx
// ✅ 好的实践 - 配置式
const breadcrumbs = [
  { label: '首页', path: '/' },
  { label: '用户管理', path: '/users' },
  { label: '用户详情' },
];

// ❌ 避免 - 硬编码
<Breadcrumb>
  <Link to="/">首页</Link>
  <span>/</span>
  <Link to="/users">用户管理</Link>
</Breadcrumb>;
```

---

## 🎨 样式定制

### 修改 Header 高度

```less
// src/components/Layout/Header.module.less
.header {
  height: 80px; // 从 64px 改为 80px
}
```

### 修改 Sidebar 宽度

```less
// src/components/Layout/Sidebar.module.less
.sidebar {
  width: 280px; // 从 240px 改为 280px
}
```

### 修改颜色主题

```less
// src/styles/variables.less
--color-black: #1a1a1a; // 从纯黑改为深灰
--color-white: #fafafa; // 从纯白改为浅灰
```

---

## 🐛 故障排除

### 问题：刷新页面后路由丢失

**原因**: Vite 开发服务器需要配置 History 模式

**解决方案**:

```ts
// vite.config.ts
export default defineConfig({
  server: {
    historyApiFallback: true, // 已配置
  },
});
```

### 问题：导航栏在移动端不显示

**检查**:

1. 检查 CSS media query
2. 检查 collapsed 状态
3. 检查 z-index 层级

### 问题：面包屑不更新

**原因**: 面包屑需要手动配置

**解决方案**: 每个页面组件中配置面包屑数据

---

## 📚 参考资源

- [React Router 文档](https://reactrouter.com/)
- [设计系统文档](./DESIGN_SYSTEM_V2.md)
- [组件库文档](./src/components/README.md)
- [Neo-brutalism 设计](https://brutalistwebsites.com/)

---

## ✅ 完成清单

- [x] Header 顶部导航栏
- [x] Sidebar 侧边栏
- [x] Layout 主布局
- [x] Breadcrumb 面包屑
- [x] Router 路由配置
- [x] ProtectedRoute 路由守卫
- [x] Dashboard 仪表盘页面
- [x] 登录跳转集成
- [x] 响应式适配
- [x] 黑白风格设计
- [x] 动画效果
- [x] 代码格式化
- [x] 文档编写

---

**实现者**: GameLink Frontend Team  
**设计风格**: Neo-brutalism 黑白极简  
**最后更新**: 2025-10-28

---

<div align="center">

## 🎉 导航系统完成！

**From Login to Dashboard**

⚫⚪ **黑白极简 · 功能完整 · 体验流畅**

</div>
