# 📚 GameLink 前端开发完整指南

**更新时间**: 2025-10-31
**版本**: v1.0
**适用对象**: 前端开发人员、全栈开发人员、UI/UX设计师

---

## 📑 目录

1. [项目概述](#1-项目概述)
2. [技术栈](#2-技术栈)
3. [项目结构](#3-项目结构)
4. [页面文档](#4-页面文档)
5. [组件库](#5-组件库)
6. [路由系统](#6-路由系统)
7. [API集成](#7-api集成)
8. [开发规范](#8-开发规范)
9. [开发指南](#9-开发指南)
10. [测试指南](#10-测试指南)
11. [部署指南](#11-部署指南)
12. [常见问题](#12-常见问题)

---

## 1. 项目概述

### 1.1 项目简介

GameLink 是一个游戏陪玩服务平台，前端采用 React + TypeScript 开发，提供用户端、管理端和陪玩师端功能。

### 1.2 功能模块

- 🔐 **用户认证**: 登录、注册、JWT Token管理
- 👥 **用户管理**: 用户列表、详情、角色分配
- 🎮 **游戏管理**: 游戏列表、详情、分类管理
- 📦 **订单管理**: 订单创建、状态流转、详情查看
- 💳 **支付管理**: 支付记录、退款处理
- ⭐ **评价管理**: 评价列表、审核
- 📊 **数据统计**: 报表、图表、数据分析
- 🔑 **权限管理**: RBAC权限控制

### 1.3 角色说明

- **普通用户**: 浏览陪玩师、下单、支付、评价
- **陪玩师**: 管理资料、接单、查看收益
- **管理员**: 管理所有数据、查看统计

---

## 2. 技术栈

### 2.1 核心技术

| 技术 | 版本 | 说明 |
|------|------|------|
| React | 18.x | 前端框架 |
| TypeScript | 5.x | 类型系统 |
| Vite | 5.x | 构建工具 |
| React Router | 6.x | 路由管理 |

### 2.2 UI组件

| 技术 | 版本 | 说明 |
|------|------|------|
| Arco Design | Latest | UI组件库 |
| Less | Latest | CSS预处理器 |
| React Icons | Latest | 图标库 |

### 2.3 状态管理

- React Context API (内置)
- React Hooks (useState, useEffect, useContext等)

### 2.4 开发工具

| 技术 | 版本 | 说明 |
|------|------|------|
| Vitest | Latest | 测试框架 |
| ESLint | Latest | 代码检查 |
| Prettier | Latest | 代码格式化 |
| @vitejs/plugin-react | Latest | Vite React插件 |

---

## 3. 项目结构

```
frontend/
├── public/                     # 静态资源
│   ├── favicon.ico
│   └── index.html
├── src/                        # 源代码
│   ├── api/                    # API调用层
│   │   ├── client.ts           # HTTP客户端
│   │   ├── auth.ts             # 认证API
│   │   ├── users.ts            # 用户API
│   │   ├── orders.ts           # 订单API
│   │   ├── games.ts            # 游戏API
│   │   ├── payments.ts         # 支付API
│   │   └── ...
│   ├── components/             # 可复用组件
│   │   ├── Button/             # 按钮组件
│   │   ├── Card/               # 卡片组件
│   │   ├── Table/              # 表格组件
│   │   ├── Modal/              # 模态框组件
│   │   ├── Layout/             # 布局组件
│   │   ├── DataTable/          # 数据表格
│   │   ├── Form/               # 表单组件
│   │   ├── Input/              # 输入框组件
│   │   └── ...
│   ├── contexts/               # React Context
│   │   ├── AuthContext.tsx     # 认证上下文
│   │   └── ThemeContext.tsx    # 主题上下文
│   ├── layouts/                # 页面布局
│   │   └── MainLayout.tsx      # 主布局
│   ├── pages/                  # 页面组件
│   │   ├── Dashboard/          # 仪表盘
│   │   ├── Login/              # 登录页
│   │   ├── Register/           # 注册页
│   │   ├── Users/              # 用户管理
│   │   ├── Games/              # 游戏管理
│   │   ├── Orders/             # 订单管理
│   │   ├── Payments/           # 支付管理
│   │   ├── Players/            # 陪玩师管理
│   │   ├── Reviews/            # 评价管理
│   │   ├── Reports/            # 报表统计
│   │   ├── Permissions/        # 权限管理
│   │   ├── Settings/           # 系统设置
│   │   ├── ComponentsDemo/     # 组件演示
│   │   └── CacheDemo/          # 缓存演示
│   ├── router/                 # 路由配置
│   │   ├── index.tsx           # 路由定义
│   │   ├── ProtectedRoute.tsx  # 路由守卫
│   │   └── layouts/            # 路由布局
│   │       └── MainLayout.tsx
│   ├── services/               # 业务服务层
│   ├── types/                  # TypeScript类型
│   │   ├── api.ts              # API类型
│   │   ├── auth.ts             # 认证类型
│   │   └── ...
│   ├── utils/                  # 工具函数
│   │   ├── crypto.ts           # 加密工具
│   │   ├── format.ts           # 格式化工具
│   │   └── ...
│   ├── styles/                 # 样式文件
│   │   ├── global.less         # 全局样式
│   │   ├── variables.less      # 变量定义
│   │   └── mixins.less         # 混合宏
│   ├── i18n/                   # 国际化
│   ├── hooks/                  # 自定义Hooks
│   ├── main.tsx                # 应用入口
│   └── App.tsx                 # 根组件
├── tests/                      # 测试文件
├── docs/                       # 项目文档
├── .env.example                # 环境变量示例
├── .eslintrc.js                # ESLint配置
├── .prettierrc.js              # Prettier配置
├── vite.config.ts              # Vite配置
└── package.json                # 依赖管理
```

---

## 4. 页面文档

### 4.1 认证页面

#### 4.1.1 登录页 (`/login`)

**路径**: `src/pages/Login/Login.tsx`
**描述**: 用户登录页面

**功能**:
- 用户名/邮箱登录
- 记住登录状态
- 加密传输
- 错误提示

**API调用**:
```typescript
authApi.login({ username, password })
```

**路由守卫**: 无需认证

---

#### 4.1.2 注册页 (`/register`)

**路径**: `src/pages/Register/Register.tsx`
**描述**: 用户注册页面

**功能**:
- 用户信息注册
- 密码强度验证
- 表单校验
- 协议同意

**路由守卫**: 无需认证

---

### 4.2 主应用页面 (需要认证)

#### 4.2.1 仪表盘 (`/dashboard`)

**路径**: `src/pages/Dashboard/Dashboard.tsx`
**描述**: 系统概览页面

**功能**:
- 数据统计卡片
- 图表展示
- 快捷操作
- 实时更新

**组件依赖**:
- Card
- Statistic
- Chart

---

#### 4.2.2 用户管理 (`/users`)

**路径列表**:
- `src/pages/Users/UserList.tsx` - 用户列表 (`/users`)
- `src/pages/Users/UserDetail.tsx` - 用户详情 (`/users/:id`)
- `src/pages/Users/UserFormModal.tsx` - 用户表单弹窗

**功能**:
- 用户列表展示
- 搜索筛选
- 批量操作
- 用户详情查看
- 创建/编辑用户
- 角色分配
- 状态管理

**组件依赖**:
- DataTable
- Modal
- Form
- ActionButtons

**API调用**:
```typescript
userApi.list(params)
userApi.get(id)
userApi.create(data)
userApi.update(id, data)
userApi.delete(id)
```

---

#### 4.2.3 游戏管理 (`/games`)

**路径列表**:
- `src/pages/Games/GameList.tsx` - 游戏列表 (`/games`)
- `src/pages/Games/GameDetail.tsx` - 游戏详情 (`/games/:id`)
- `src/pages/Games/GameFormModal.tsx` - 游戏表单弹窗

**功能**:
- 游戏列表
- 游戏详情
- 创建/编辑游戏
- 分类管理
- 图标上传

**组件依赖**:
- DataTable
- Upload
- Tag
- Image

---

#### 4.2.4 订单管理 (`/orders`)

**路径列表**:
- `src/pages/Orders/OrderList.tsx` - 订单列表 (`/orders`)
- `src/pages/Orders/OrderDetail.tsx` - 订单详情 (`/orders/:id`)
- `src/pages/Orders/OrderFormModal.tsx` - 订单表单弹窗

**功能**:
- 订单列表
- 订单详情
- 状态流转
- 搜索筛选
- 导出功能

**组件依赖**:
- DataTable
- Timeline
- Steps
- Tag

**订单状态**:
- Pending (待处理)
- Accepted (已接单)
- InProgress (进行中)
- Completed (已完成)
- Cancelled (已取消)

---

#### 4.2.5 支付管理 (`/payments`)

**路径列表**:
- `src/pages/Payments/PaymentList.tsx` - 支付列表 (`/payments`)
- `src/pages/Payments/PaymentDetailPage.tsx` - 支付详情 (`/payments/:id`)

**功能**:
- 支付列表
- 支付详情
- 退款处理
- 交易统计

**组件依赖**:
- DataTable
- Description
- Statistic

---

#### 4.2.6 陪玩师管理 (`/players`)

**路径列表**:
- `src/pages/Players/PlayerList.tsx` - 陪玩师列表 (`/players`)
- `src/pages/Players/PlayerFormModal.tsx` - 陪玩师表单弹窗

**功能**:
- 陪玩师列表
- 技能标签
- 验证状态
- 在线状态

**组件依赖**:
- DataTable
- Tag
- Badge
- Avatar

---

#### 4.2.7 评价管理 (`/reviews`)

**路径列表**:
- `src/pages/Reviews/ReviewList.tsx` - 评价列表 (`/reviews`)
- `src/pages/Reviews/ReviewFormModal.tsx` - 评价表单弹窗

**功能**:
- 评价列表
- 评分展示
- 审核功能

**组件依赖**:
- DataTable
- Rating
- ReviewModal

---

#### 4.2.8 报表统计 (`/reports`)

**路径**: `src/pages/Reports/ReportDashboard.tsx`
**描述**: 数据报表页面

**功能**:
- 数据统计
- 图表分析
- 趋势预测
- 导出报表

**组件依赖**:
- Card
- Chart
- Statistic
- DatePicker

---

#### 4.2.9 权限管理 (`/permissions`)

**路径**: `src/pages/Permissions/PermissionList.tsx`
**描述**: 权限管理页面

**功能**:
- 权限列表
- 角色分配
- 权限继承

**组件依赖**:
- Tree
- DataTable
- Modal

---

#### 4.2.10 系统设置 (`/settings`)

**路径**: `src/pages/Settings/SettingsDashboard.tsx`
**描述**: 系统设置页面

**功能**:
- 系统配置
- 用户偏好
- 主题切换

**组件依赖**:
- Form
- Switch
- Select

---

### 4.3 演示页面

#### 4.3.1 组件演示 (`/showcase`)

**路径**: `src/pages/ComponentsDemo/ComponentsDemo.tsx`
**描述**: 组件库演示页面

**功能**:
- 展示所有组件
- 代码示例
- Props演示

**路由守卫**: 无需认证

---

#### 4.3.2 缓存演示 (`/cache-demo`)

**路径**: `src/pages/CacheDemo/index.tsx`
**描述**: 路由缓存演示

**子路由**:
- `/cache-demo/a` - 页面A
- `/cache-demo/b` - 页面B

**路由守卫**: 无需认证

---

## 5. 组件库

### 5.1 基础组件

#### 5.1.1 Button 按钮

**路径**: `src/components/Button/Button.tsx`

**Props**:
```typescript
interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'text' | 'outlined'
  size?: 'small' | 'medium' | 'large'
  block?: boolean
  loading?: boolean
  disabled?: boolean
  icon?: ReactNode
  children?: ReactNode
  onClick?: (event: MouseEvent<HTMLButtonElement>) => void
  type?: 'button' | 'submit' | 'reset'
}
```

**用法示例**:
```tsx
import { Button } from '@/components'

<Button>默认按钮</Button>
<Button variant="outlined">轮廓按钮</Button>
<Button size="small" icon={<Icon />}>带图标</Button>
<Button loading>加载中</Button>
```

---

#### 5.1.2 Input 输入框

**路径**: `src/components/Input/Input.tsx`

**Props**:
```typescript
interface InputProps {
  value?: string
  defaultValue?: string
  placeholder?: string
  disabled?: boolean
  readonly?: boolean
  size?: 'small' | 'medium' | 'large'
  prefix?: ReactNode
  suffix?: ReactNode
  maxLength?: number
  showWordLimit?: boolean
}
```

**用法示例**:
```tsx
import { Input } from '@/components'

<Input placeholder="请输入" />
<Input prefix={<IconUser />} />
<Input showWordLimit maxLength={100} />
```

---

#### 5.1.3 Card 卡片

**路径**: `src/components/Card/Card.tsx`

**Props**:
```typescript
interface CardProps {
  title?: ReactNode
  extra?: ReactNode
  bordered?: boolean
  hoverable?: boolean
  loading?: boolean
  className?: string
  style?: CSSProperties
  children?: ReactNode
}
```

**用法示例**:
```tsx
import { Card } from '@/components'

<Card title="标题" extra={<Button>操作</Button>}>
  内容区域
</Card>
```

---

### 5.2 布局组件

#### 5.2.1 Layout 布局

**路径**: `src/components/Layout/Layout.tsx`

**子组件**:
- Header - 顶部导航
- Sidebar - 侧边栏
- Content - 内容区
- Footer - 底部

**用法示例**:
```tsx
import { Layout } from '@/components'

<Layout>
  <Layout.Header>顶部</Layout.Header>
  <Layout>
    <Layout.Sidebar>侧边</Layout.Sidebar>
    <Layout.Content>内容</Layout.Content>
  </Layout>
</Layout>
```

---

#### 5.2.2 Breadcrumb 面包屑

**路径**: `src/components/Breadcrumb/Breadcrumb.tsx`

**Props**:
```typescript
interface BreadcrumbProps {
  routes: Array<{
    path?: string
    breadcrumbName: string
  }>
}
```

---

### 5.3 数据展示组件

#### 5.3.1 DataTable 数据表格

**路径**: `src/components/DataTable/DataTable.tsx`

**功能**:
- 列定义
- 排序
- 筛选
- 分页
- 批量操作
- 行选择

**用法示例**:
```tsx
import { DataTable } from '@/components'

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id' },
  { title: '姓名', dataIndex: 'name', key: 'name' },
]

<DataTable
  columns={columns}
  dataSource={data}
  loading={loading}
  pagination={pagination}
  onSearch={handleSearch}
/>
```

---

#### 5.3.2 Pagination 分页

**路径**: `src/components/Pagination/Pagination.tsx`

**Props**:
```typescript
interface PaginationProps {
  current: number
  pageSize: number
  total: number
  onChange: (page: number, pageSize: number) => void
  showSizeChanger?: boolean
  showQuickJumper?: boolean
  showTotal?: (total: number, range: [number, number]) => string
}
```

---

#### 5.3.3 Rating 评分

**路径**: `src/components/Rating/Rating.tsx`

**Props**:
```typescript
interface RatingProps {
  value?: number
  defaultValue?: number
  count?: number
  size?: 'small' | 'medium' | 'large'
  allowClear?: boolean
  readOnly?: boolean
  onChange?: (value: number) => void
}
```

---

### 5.4 表单组件

#### 5.4.1 Form 表单

**路径**: `src/components/Form/Form.tsx`

**功能**:
- 表单验证
- 动态表单
- 布局控制
- 提交处理

**用法示例**:
```tsx
import { Form } from '@/components'

<Form
  form={form}
  layout="vertical"
  onFinish={handleSubmit}
>
  <Form.Item name="username" label="用户名" rules={[{ required: true }]}>
    <Input />
  </Form.Item>
</Form>
```

---

#### 5.4.2 FormField 表单字段

**路径**: `src/components/FormField/FormField.tsx`

---

### 5.5 弹窗组件

#### 5.5.1 Modal 模态框

**路径**: `src/components/Modal/Modal.tsx`

**Props**:
```typescript
interface ModalProps {
  visible: boolean
  title?: ReactNode
  width?: number | string
  footer?: ReactNode
  onCancel: () => void
  onOk?: () => void
  children?: ReactNode
}
```

---

#### 5.5.2 DeleteConfirmModal 删除确认

**路径**: `src/components/DeleteConfirmModal/DeleteConfirmModal.tsx`

**功能**:
- 删除确认
- 危险操作提示
- 加载状态

---

#### 5.5.3 ReviewModal 评价弹窗

**路径**: `src/components/ReviewModal/ReviewModal.tsx`

---

### 5.6 反馈组件

#### 5.6.1 Notification 通知

**路径**: `src/components/Notification/Notification.tsx`

**用法示例**:
```tsx
import { notification } from '@/components'

notification.success({
  message: '成功',
  description: '操作成功',
})
```

---

### 5.7 其他组件

#### 5.7.1 Badge 角标

**路径**: `src/components/Badge/Badge.tsx`

---

#### 5.7.2 Tag 标签

**路径**: `src/components/Tag/Tag.tsx`

---

#### 5.7.3 Skeleton 骨架屏

**路径**: `src/components/Skeleton/Skeleton.tsx`

---

#### 5.7.4 RouteCache 路由缓存

**路径**: `src/components/RouteCache/RouteCache.tsx`

**功能**:
- 路由级缓存
- 状态保持
- 性能优化

---

#### 5.7.5 ActionButtons 操作按钮

**路径**: `src/components/ActionButtons/ActionButtons.tsx`

---

#### 5.7.6 BulkActions 批量操作

**路径**: `src/components/BulkActions/BulkActions.tsx`

---

#### 5.7.7 Grid 网格

**路径**: `src/components/Grid/Grid.tsx`

---

#### 5.7.8 Menu 菜单

**路径**: `src/components/Menu/Menu.tsx`

---

#### 5.7.9 Select 选择器

**路径**: `src/components/Select/Select.tsx`

---

#### 5.7.10 Space 间距

**路径**: `src/components/Space/Space.tsx`

---

#### 5.7.11 Tabs 标签页

**路径**: `src/components/Tabs/Tabs.tsx`

---

#### 5.7.12 Table 表格

**路径**: `src/components/Table/Table.tsx`

---

## 6. 路由系统

### 6.1 路由配置

**配置文件**: `src/router/index.tsx`

**使用技术**: React Router v6

**路由类型**:
- BrowserRouter - 浏览器路由
- HashRouter - 哈希路由

### 6.2 路由结构

```typescript
export const router = createBrowserRouter([
  // 公开路由（无需认证）
  { path: '/login', element: <Login /> },
  { path: '/register', element: <Register /> },

  // 演示路由（无需认证）
  { path: '/showcase', element: <ComponentsDemo /> },
  {
    path: '/cache-demo',
    element: <CacheDemo />,
    children: [
      { index: true, element: <Navigate to="/cache-demo/a" replace /> },
      { path: 'a', element: <CachePageA /> },
      { path: 'b', element: <CachePageB /> },
    ],
  },

  // 受保护的路由（需要认证）
  {
    path: '/',
    element: <ProtectedRoute />,
    children: [
      {
        element: <MainLayout />,
        children: [
          { index: true, element: <Navigate to="/dashboard" replace /> },
          { path: 'dashboard', element: <Dashboard /> },
          { path: 'orders', element: <OrderList /> },
          { path: 'orders/:id', element: <OrderDetail /> },
          { path: 'games', element: <GameList /> },
          { path: 'games/:id', element: <GameDetail /> },
          { path: 'players', element: <PlayerList /> },
          { path: 'users', element: <UserList /> },
          { path: 'users/:id', element: <UserDetail /> },
          { path: 'payments', element: <PaymentList /> },
          { path: 'payments/:id', element: <PaymentDetailPage /> },
          { path: 'reviews', element: <ReviewList /> },
          { path: 'reports', element: <ReportDashboard /> },
          { path: 'permissions', element: <PermissionList /> },
          { path: 'settings', element: <SettingsDashboard /> },
          { path: 'components', element: <ComponentsDemo /> },
        ],
      },
    ],
  },

  // 404路由
  { path: '*', element: <Navigate to="/dashboard" replace /> },
])
```

### 6.3 路由守卫

**组件**: `src/router/ProtectedRoute.tsx`

**功能**:
- 检查认证状态
- 未登录重定向到登录页
- Token过期处理
- 权限验证

**实现**:
```typescript
export const ProtectedRoute = () => {
  const { user, loading } = useAuth()
  const location = useLocation()

  if (loading) {
    return <Spin />
  }

  if (!user) {
    return <Navigate to="/login" state={{ from: location }} replace />
  }

  return <Outlet />
}
```

### 6.4 布局组件

**主布局**: `src/router/layouts/MainLayout.tsx`

**包含组件**:
- Header - 顶部导航
- Sidebar - 侧边菜单
- Content - 主内容区
- Footer - 底部

### 6.5 路由缓存

**组件**: `RouteCache`

**功能**:
- 缓存路由组件
- 保持页面状态
- 前进后退优化

**用法**:
```tsx
<RouteCache>
  <YourComponent />
</RouteCache>
```

---

## 7. API集成

### 7.1 API客户端

**配置文件**: `src/api/client.ts`

**使用Axios进行HTTP请求**

```typescript
import axios from 'axios'

const client = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
client.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器
client.interceptors.response.use(
  (response) => response.data,
  (error) => {
    // 错误处理
    return Promise.reject(error)
  }
)
```

### 7.2 加密支持

**配置**: `src/utils/crypto.ts`

**功能**:
- 请求数据加密
- 响应数据解密
- 签名验证

**启用方式**:
```typescript
// .env
VITE_CRYPTO_ENABLED=true
```

### 7.3 API模块

#### 7.3.1 认证API

**文件**: `src/api/auth.ts`

**接口**:
```typescript
export const authApi = {
  login: (data: LoginRequest) => client.post('/auth/login', data),
  register: (data: RegisterRequest) => client.post('/auth/register', data),
  refresh: () => client.post('/auth/refresh'),
  logout: () => client.post('/auth/logout'),
  me: () => client.get('/auth/me'),
}
```

---

#### 7.3.2 用户API

**文件**: `src/api/users.ts`

**接口**:
```typescript
export const userApi = {
  list: (params?: UserListParams) => client.get('/admin/users', { params }),
  get: (id: number) => client.get(`/admin/users/${id}`),
  create: (data: CreateUserRequest) => client.post('/admin/users', data),
  update: (id: number, data: UpdateUserRequest) => client.put(`/admin/users/${id}`, data),
  delete: (id: number) => client.delete(`/admin/users/${id}`),
  updateStatus: (id: number, status: UserStatus) => client.patch(`/admin/users/${id}/status`),
}
```

---

#### 7.3.3 订单API

**文件**: `src/api/orders.ts`

**接口**:
```typescript
export const orderApi = {
  list: (params?: OrderListParams) => client.get('/admin/orders', { params }),
  get: (id: number) => client.get(`/admin/orders/${id}`),
  create: (data: CreateOrderRequest) => client.post('/admin/orders', data),
  update: (id: number, data: UpdateOrderRequest) => client.put(`/admin/orders/${id}`, data),
  delete: (id: number) => client.delete(`/admin/orders/${id}`),
  cancel: (id: number, reason: string) => client.post(`/admin/orders/${id}/cancel`, { reason }),
}
```

---

#### 7.3.4 游戏API

**文件**: `src/api/games.ts`

---

#### 7.3.5 支付API

**文件**: `src/api/payments.ts`

---

#### 7.3.6 评价API

**文件**: `src/api/reviews.ts`

---

### 7.4 错误处理

**统一错误处理**:
```typescript
client.interceptors.response.use(
  (response) => response,
  (error) => {
    const { response } = error

    if (response?.status === 401) {
      // Token过期，跳转到登录页
      redirectToLogin()
    }

    if (response?.status === 403) {
      // 无权限
      message.error('没有权限访问该资源')
    }

    return Promise.reject(error)
  }
)
```

### 7.5 数据类型

**类型定义**: `src/types/api.ts`

**示例**:
```typescript
interface ApiResponse<T> {
  success: boolean
  code: number
  message: string
  data: T
}

interface User {
  id: number
  name: string
  email: string
  role: string
  status: string
  createdAt: string
}
```

---

## 8. 开发规范

### 8.1 编码规范

**参考**: `docs/design/CODING_STANDARDS.md`

#### 8.1.1 文件命名

- **组件**: PascalCase (如: `UserList.tsx`)
- **页面**: PascalCase (如: `Dashboard.tsx`)
- **工具**: camelCase (如: `formatDate.ts`)
- **常量**: UPPER_SNAKE_CASE (如: `API_ENDPOINTS.ts`)
- **类型**: PascalCase + Type后缀 (如: `UserType.ts`)

#### 8.1.2 目录命名

- **页面**: 复数形式 (如: `pages/Users/`)
- **组件**: 复数形式 (如: `components/Buttons/`)
- **hooks**: 以use开头 (如: `hooks/useAuth.ts`)
- **utils**: 小写 (如: `utils/format.ts`)

#### 8.1.3 导入顺序

```typescript
// 1. React相关
import React, { useState, useEffect } from 'react'

// 2. 第三方库
import { Button, Card } from 'arco-design'
import { useParams } from 'react-router-dom'

// 3. 内部模块
import { UserService } from '@/services/user'
import { formatDate } from '@/utils/format'

// 4. 相对导入
import { Component } from './Component'
```

### 8.2 组件规范

#### 8.2.1 组件结构

```tsx
import React from 'react'

interface Props {
  title: string
  size?: 'small' | 'medium' | 'large'
}

// 组件实现
export const Component: React.FC<Props> = ({ title, size = 'medium' }) => {
  // 状态
  const [state, setState] = useState()

  // 副作用
  useEffect(() => {
    // 副作用逻辑
  }, [])

  // 事件处理
  const handleClick = () => {
    // 处理逻辑
  }

  // 渲染
  return (
    <div>
      <h1>{title}</h1>
    </div>
  )
}
```

#### 8.2.2 Props规范

- 所有Props必须有明确的类型定义
- 非必填Props需要提供默认值
- 布尔值Props使用isXxx或hasXxx命名

#### 8.2.3 状态管理

- 优先使用useState管理本地状态
- 复杂状态考虑使用useReducer
- 全局状态使用Context API

### 8.3 样式规范

#### 8.3.1 Less使用

- 使用Less模块化
- 避免全局样式污染
- 使用变量统一管理颜色、字体等

#### 8.3.2 类名规范

```css
/* BEM规范 */
.block {}
.block__element {}
.block--modifier {}

/* 示例 */
.user-card {}
.user-card__header {}
.user-card--highlighted {}
```

#### 8.3.3 样式变量

```less
// variables.less
@primary-color: #1890ff;
@success-color: #52c41a;
@warning-color: #faad14;
@error-color: #f5222d;

@font-size-sm: 12px;
@font-size-base: 14px;
@font-size-lg: 16px;

@spacing-xs: 4px;
@spacing-sm: 8px;
@spacing-base: 16px;
@spacing-lg: 24px;
```

### 8.4 TypeScript规范

#### 8.4.1 类型定义

- 优先使用interface定义对象类型
- 使用type定义联合类型、基础类型
- 避免使用any类型

#### 8.4.2 泛型使用

```typescript
// 通用列表组件
interface ListProps<T> {
  data: T[]
  renderItem: (item: T, index: number) => React.ReactNode
}

function List<T>({ data, renderItem }: ListProps<T>) {
  return (
    <ul>
      {data.map((item, index) => (
        <li key={index}>{renderItem(item, index)}</li>
      ))}
    </ul>
  )
}
```

### 8.5 Git规范

#### 8.5.1 提交信息

```
type(scope): subject

body

footer
```

**类型 (type)**:
- feat: 新功能
- fix: 修复
- docs: 文档
- style: 格式化
- refactor: 重构
- test: 测试
- chore: 构建/工具

**示例**:
```
feat(users): 添加用户列表搜索功能

实现用户列表的搜索和筛选功能
- 添加搜索框组件
- 集成API接口
- 添加加载状态

Closes #123
```

#### 8.5.2 分支命名

- `feature/*` - 新功能
- `fix/*` - 修复
- `hotfix/*` - 紧急修复
- `release/*` - 发布

---

## 9. 开发指南

### 9.1 开发环境搭建

#### 9.1.1 环境要求

```bash
Node.js >= 18.0.0
npm >= 8.0.0 或 yarn >= 1.22.0
Git >= 2.30.0
```

#### 9.1.2 快速开始

```bash
# 1. 克隆项目
git clone https://github.com/your-org/gamelink.git
cd gamelink/frontend

# 2. 安装依赖
npm install
# 或
yarn install

# 3. 配置环境变量
cp .env.example .env.local
# 编辑 .env.local 文件

# 4. 启动开发服务器
npm run dev
# 或
yarn dev

# 5. 访问应用
# http://localhost:5173
```

#### 9.1.3 环境变量

**`.env.example`**:
```bash
# API配置
VITE_API_BASE_URL=http://localhost:8080/api/v1

# 加密配置
VITE_CRYPTO_ENABLED=true
VITE_CRYPTO_SECRET_KEY=your-32-byte-secret-key
VITE_CRYPTO_IV=your-iv-16-byte

# 开发配置
VITE_DEV_TOOLS=true
VITE_MOCK_DATA=false
```

### 9.2 开发流程

#### 9.2.1 创建新页面

1. 在`src/pages/`下创建页面目录
2. 创建页面组件文件
3. 在路由中注册页面
4. 添加页面菜单（如果需要）

**示例**:
```bash
# 创建目录
mkdir src/pages/NewPage

# 创建组件文件
touch src/pages/NewPage/NewPage.tsx
touch src/pages/NewPage/index.ts

# 注册路由 (src/router/index.tsx)
{
  path: 'new-page',
  element: <NewPage />,
}
```

#### 9.2.2 创建新组件

1. 在`src/components/`下创建组件目录
2. 创建组件文件和样式文件
3. 在`src/components/index.ts`中导出
4. 编写测试文件

**示例**:
```bash
# 创建目录
mkdir src/components/NewComponent

# 创建文件
touch src/components/NewComponent/NewComponent.tsx
touch src/components/NewComponent/NewComponent.less
touch src/components/NewComponent/index.ts
touch src/components/NewComponent/NewComponent.test.tsx

# 导出 (src/components/index.ts)
export { default as NewComponent } from './NewComponent'
```

#### 9.2.3 API集成

1. 在`src/api/`下创建或修改API文件
2. 定义接口类型
3. 编写API调用函数
4. 在页面中调用

**示例**:
```typescript
// src/api/example.ts
import { client } from './client'

interface ExampleRequest {
  id: number
}

interface ExampleResponse {
  id: number
  name: string
}

export const exampleApi = {
  get: (id: number) => client.get<ExampleResponse>(`/example/${id}`),
  create: (data: ExampleRequest) => client.post<ExampleResponse>('/example', data),
}
```

### 9.3 调试指南

#### 9.3.1 常用调试工具

- **React DevTools**: 组件调试
- **Redux DevTools**: 状态管理调试
- **Vite DevTools**: 构建工具调试

#### 9.3.2 日志打印

```typescript
// 开发环境下打印
if (import.meta.env.DEV) {
  console.log('Debug info:', data)
}
```

#### 9.3.3 错误边界

```tsx
class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props)
    this.state = { hasError: false }
  }

  static getDerivedStateFromError(error) {
    return { hasError: true }
  }

  componentDidCatch(error, errorInfo) {
    console.error('Error:', error, errorInfo)
  }

  render() {
    if (this.state.hasError) {
      return <h1>Something went wrong.</h1>
    }

    return this.props.children
  }
}
```

### 9.4 性能优化

#### 9.4.1 React性能优化

**使用Memo**:
```tsx
import { memo } from 'react'

interface Props {
  data: Data[]
}

const List = memo<Props>(({ data }) => {
  return (
    <ul>
      {data.map(item => <li key={item.id}>{item.name}</li>)}
    </ul>
  )
})
```

**使用Callback**:
```tsx
import { useCallback } from 'react'

const Parent = () => {
  const handleClick = useCallback((id: number) => {
    console.log('Click:', id)
  }, [])

  return <Child onClick={handleClick} />
}
```

**使用Memoize Expensive Operations**:
```tsx
import { useMemo } from 'react'

const expensiveValue = useMemo(() => {
  return data.filter(item => item.active).sort((a, b) => a.name.localeCompare(b.name))
}, [data])
```

#### 9.4.2 路由优化

**使用懒加载**:
```tsx
import { lazy, Suspense } from 'react'

const Dashboard = lazy(() => import('pages/Dashboard'))

export const router = createBrowserRouter([
  {
    path: '/dashboard',
    element: (
      <Suspense fallback={<Spin />}>
        <Dashboard />
      </Suspense>
    ),
  },
])
```

#### 9.4.3 代码分割

**动态导入**:
```tsx
const loadComponent = async () => {
  const { HeavyComponent } = await import('./HeavyComponent')
  return <HeavyComponent />
}
```

### 9.5 国际化

#### 9.5.1 配置

**使用react-i18next**:
```typescript
// src/i18n/index.ts
import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'

i18n
  .use(initReactI18next)
  .init({
    lng: 'zh-CN',
    fallbackLng: 'en',
    resources: {
      zh: {
        translation: {
          welcome: '欢迎',
        },
      },
      en: {
        translation: {
          welcome: 'Welcome',
        },
      },
    },
  })

export default i18n
```

#### 9.5.2 使用

```tsx
import { useTranslation } from 'react-i18next'

const Component = () => {
  const { t } = useTranslation()

  return <h1>{t('welcome')}</h1>
}
```

---

## 10. 测试指南

### 10.1 测试工具

- **Vitest**: 单元测试框架
- **React Testing Library**: 组件测试
- **MSW**: API Mock

### 10.2 编写测试

#### 10.2.1 组件测试

```tsx
// Button.test.tsx
import { render, screen } from '@testing-library/react'
import { Button } from './Button'

describe('Button', () => {
  it('renders correctly', () => {
    render(<Button>Click me</Button>)
    expect(screen.getByRole('button')).toHaveTextContent('Click me')
  })

  it('calls onClick when clicked', () => {
    const handleClick = jest.fn()
    render(<Button onClick={handleClick}>Click me</Button>)
    screen.getByRole('button').click()
    expect(handleClick).toHaveBeenCalledTimes(1)
  })
})
```

#### 10.2.2 Hook测试

```tsx
// useAuth.test.ts
import { renderHook, act } from '@testing-library/react'
import { useAuth } from './useAuth'

describe('useAuth', () => {
  it('provides authentication state', () => {
    const { result } = renderHook(() => useAuth())
    expect(result.current.user).toBeNull()
  })
})
```

#### 10.2.3 API测试

```typescript
// api.test.ts
import { authApi } from './auth'
import { mockServer } from '@/tests/mock'

describe('authApi', () => {
  beforeAll(() => {
    mockServer.listen()
  })

  afterEach(() => {
    mockServer.resetHandlers()
  })

  afterAll(() => {
    mockServer.close()
  })

  it('logs in successfully', async () => {
    const response = await authApi.login({ username: 'admin', password: '123456' })
    expect(response.success).toBe(true)
  })
})
```

### 10.3 运行测试

```bash
# 运行所有测试
npm test

# 运行测试并生成覆盖率报告
npm run test:coverage

# 监听模式
npm run test:watch

# 运行特定测试文件
npm test Button.test.tsx
```

---

## 11. 部署指南

### 11.1 构建

```bash
# 开发环境构建
npm run build:dev

# 生产环境构建
npm run build

# 预览构建结果
npm run preview
```

### 11.2 部署配置

#### 11.2.1 Nginx配置

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /var/www/gamelink/frontend;
    index index.html;

    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # SPA路由支持
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API代理
    location /api {
        proxy_pass http://backend-server;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

#### 11.2.2 Docker部署

```dockerfile
# Dockerfile
FROM node:18-alpine as builder

WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### 11.3 CI/CD

#### 11.3.1 GitHub Actions

```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'

      - name: Install dependencies
        run: npm ci

      - name: Run tests
        run: npm test

      - name: Build
        run: npm run build

      - name: Deploy
        run: |
          # 部署到服务器
```

---

## 12. 常见问题

### 12.1 开发问题

#### Q: 如何解决路由刷新404问题？
**A**: 配置服务器将所有路由重定向到index.html（Nginx、Apache等）

#### Q: 如何处理Token过期？
**A**: 在响应拦截器中检测401状态码，自动跳转到登录页

#### Q: 如何优化大数据列表渲染？
**A**: 使用虚拟滚动（react-window）或分页加载

#### Q: 如何处理组件重复渲染？
**A**: 使用React.memo、useMemo、useCallback等优化手段

### 12.2 调试问题

#### Q: 如何调试异步代码？
**A**: 使用Chrome DevTools的异步调试功能，或在代码中添加断点

#### Q: 如何查看Redux状态？
**A**: 安装Redux DevTools浏览器扩展

#### Q: 如何调试网络请求？
**A**: 使用浏览器Network面板或安装相关扩展

### 12.3 性能问题

#### Q: 首屏加载慢怎么办？
**A**: 使用路由懒加载、代码分割、资源压缩等手段

#### Q: 内存泄漏如何排查？
**A**: 使用Chrome DevTools的Memory面板，或开启React严格模式

### 12.4 部署问题

#### Q: 静态资源404？
**A**: 检查部署路径配置，确保静态资源路径正确

#### Q: API请求跨域？
**A**: 配置CORS或使用代理

---

## 📚 附录

### A. 相关文档

- [开发指南](./DEVELOPER_GUIDE.md)
- [技术文档](./TECHNICAL_DOCUMENTATION.md)
- [用户文档](./USER_DOCUMENTATION.md)
- [API文档](./api/)
- [加密文档](./crypto/)

### B. 外部资源

- [React 官方文档](https://react.dev/)
- [TypeScript 官方文档](https://www.typescriptlang.org/)
- [Vite 官方文档](https://vitejs.dev/)
- [Arco Design 组件库](https://arco.design/)
- [React Router 文档](https://reactrouter.com/)
- [Axios 文档](https://axios-http.com/)

### C. 工具推荐

- **VSCode Extensions**:
  - ES7+ React/Redux/React-Native snippets
  - TypeScript Importer
  - Prettier - Code formatter
  - ESLint
  - Auto Rename Tag
  - Bracket Pair Colorizer

### D. 参考资料

- [React性能优化指南](https://react.dev/learn/render-and-commit)
- [TypeScript最佳实践](https://typescript-eslint.io/)
- [Vite构建优化](https://vitejs.dev/guide/build.html)

---

**文档维护者**: GameLink Frontend Team
**最后更新**: 2025-10-31
**版本**: v1.0
**反馈**: 如有问题请提交Issue或PR
