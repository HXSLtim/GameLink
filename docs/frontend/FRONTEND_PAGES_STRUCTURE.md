# 📄 前端页面结构完整文档

**更新时间**: 2025-10-31
**页面总数**: 22个
**文档类型**: 页面结构说明

---

## 📑 目录

1. [页面概览](#1-页面概览)
2. [认证页面](#2-认证页面)
3. [管理页面](#3-管理页面)
4. [演示页面](#4-演示页面)
5. [页面关联关系](#5-页面关联关系)
6. [路由配置](#6-路由配置)
7. [组件依赖](#7-组件依赖)
8. [页面开发规范](#8-页面开发规范)

---

## 1. 页面概览

### 1.1 页面分类

| 分类 | 数量 | 说明 |
|------|------|------|
| 认证页面 | 2 | 登录、注册 |
| 管理页面 | 17 | 各种管理功能 |
| 演示页面 | 3 | 组件、缓存演示 |

### 1.2 认证状态

- **需要认证**: 17个页面
- **无需认证**: 5个页面 (登录、注册、组件演示、缓存演示)

### 1.3 布局类型

- **MainLayout**: 主布局（包含Header、Sidebar、Content）
- **独立布局**: 无需认证的页面

---

## 2. 认证页面

### 2.1 登录页 (`/login`)

**文件路径**: `src/pages/Login/Login.tsx`
**依赖文件**: `src/pages/Login/README.md`

#### 功能特性
- ✅ 用户名/邮箱登录
- ✅ 密码登录
- ✅ 记住登录状态
- ✅ 加密传输 (Crypto Middleware)
- ✅ 错误提示
- ✅ 加载状态
- ✅ 自动跳转

#### API调用
```typescript
authApi.login({ username, password })
```

#### 组件依赖
- Form (表单)
- Input (输入框)
- Button (按钮)
- Checkbox (记住我)
- Message (消息提示)

#### 路由守卫
- **无需认证**
- 已登录用户访问自动重定向到 `/dashboard`

---

### 2.2 注册页 (`/register`)

**文件路径**: `src/pages/Register/Register.tsx`

#### 功能特性
- ✅ 用户信息注册
- ✅ 密码确认
- ✅ 表单验证
- ✅ 协议同意
- ✅ 自动登录
- ✅ 错误处理

#### API调用
```typescript
authApi.register({ name, email, password })
```

#### 组件依赖
- Form (表单)
- Input (输入框)
- Button (按钮)
- Checkbox (协议)
- Message (消息提示)

#### 路由守卫
- **无需认证**
- 已登录用户访问自动重定向到 `/dashboard`

---

## 3. 管理页面

### 3.1 仪表盘 (`/dashboard`)

**文件路径**: `src/pages/Dashboard/Dashboard.tsx`

#### 功能特性
- ✅ 数据统计卡片
- ✅ 图表展示
- ✅ 快捷操作
- ✅ 实时数据更新

#### 展示内容
- 用户总数
- 订单数量
- 收入统计
- 待处理事项

#### 组件依赖
- Card (卡片)
- Statistic (统计)
- Chart (图表)
- Row/Col (栅格)
- Button (按钮)

---

### 3.2 订单管理

#### 3.2.1 订单列表 (`/orders`)

**文件路径**: `src/pages/Orders/OrderList.tsx`

##### 功能特性
- ✅ 订单列表展示
- ✅ 分页功能
- ✅ 搜索筛选 (状态、时间、用户)
- ✅ 排序功能
- ✅ 批量操作
- ✅ 导出功能
- ✅ 创建订单按钮

##### 操作按钮
- 查看详情
- 编辑订单
- 取消订单
- 删除订单
- 导出数据

##### 组件依赖
- DataTable (数据表格)
- Button (按钮)
- Modal (弹窗)
- Select (选择器)
- DatePicker (日期选择器)
- Input (搜索框)

##### API调用
```typescript
orderApi.list(params)
orderApi.create(data)
orderApi.update(id, data)
orderApi.delete(id)
```

---

#### 3.2.2 订单详情 (`/orders/:id`)

**文件路径**: `src/pages/Orders/OrderDetail.tsx`

##### 功能特性
- ✅ 订单详细信息
- ✅ 状态流转
- ✅ 操作日志
- ✅ 相关文档
- ✅ 支付信息
- ✅ 用户信息

##### 页面结构
```
顶部: 面包屑导航 + 操作按钮
主体:
  - 基本信息卡片
  - 状态流转卡片
  - 用户信息卡片
  - 陪玩师信息卡片
  - 支付信息卡片
  - 操作日志卡片
  - 评价信息卡片
```

##### 组件依赖
- Card (卡片)
- Description (描述列表)
- Steps (步骤条)
- Timeline (时间线)
- Tag (标签)
- Button (按钮)
- Modal (弹窗)

##### API调用
```typescript
orderApi.get(id)
orderApi.updateStatus(id, status)
```

---

#### 3.2.3 订单表单弹窗 (`/orders/OrderFormModal.tsx`)

**文件路径**: `src/pages/Orders/OrderFormModal.tsx`

##### 功能特性
- ✅ 创建订单
- ✅ 编辑订单
- ✅ 表单验证
- ✅ 动态表单

##### 组件依赖
- Modal (弹窗)
- Form (表单)
- Select (选择器)
- Input (输入框)
- DatePicker (日期选择器)
- InputNumber (数字输入)
- Button (按钮)

---

### 3.3 游戏管理

#### 3.3.1 游戏列表 (`/games`)

**文件路径**: `src/pages/Games/GameList.tsx`

##### 功能特性
- ✅ 游戏列表展示
- ✅ 搜索筛选
- ✅ 分类管理
- ✅ 图标展示
- ✅ 状态管理

##### 组件依赖
- DataTable (数据表格)
- Image (图片)
- Tag (标签)
- Button (按钮)
- Modal (弹窗)

##### API调用
```typescript
gameApi.list(params)
gameApi.create(data)
gameApi.update(id, data)
gameApi.delete(id)
```

---

#### 3.3.2 游戏详情 (`/games/:id`)

**文件路径**: `src/pages/Games/GameDetail.tsx`

##### 功能特性
- ✅ 游戏详细信息
- ✅ 游戏统计
- ✅ 陪玩师数量
- ✅ 订单数量

##### 组件依赖
- Card (卡片)
- Description (描述列表)
- Statistic (统计)
- Image (图片)

##### API调用
```typescript
gameApi.get(id)
```

---

#### 3.3.3 游戏表单弹窗 (`/games/GameFormModal.tsx`)

**文件路径**: `src/pages/Games/GameFormModal.tsx`

##### 功能特性
- ✅ 创建游戏
- ✅ 编辑游戏
- ✅ 图标上传
- ✅ 分类选择

##### 组件依赖
- Modal (弹窗)
- Form (表单)
- Input (输入框)
- Select (选择器)
- Upload (上传)
- Button (按钮)

---

### 3.4 用户管理

#### 3.4.1 用户列表 (`/users`)

**文件路径**: `src/pages/Users/UserList.tsx`

##### 功能特性
- ✅ 用户列表展示
- ✅ 搜索筛选 (姓名、邮箱、角色)
- ✅ 状态筛选
- ✅ 批量操作
- ✅ 创建用户按钮

##### 操作按钮
- 查看详情
- 编辑用户
- 删除用户
- 更新状态
- 分配角色

##### 组件依赖
- DataTable (数据表格)
- Button (按钮)
- Modal (弹窗)
- Select (选择器)
- Tag (标签)
- Avatar (头像)

##### API调用
```typescript
userApi.list(params)
userApi.create(data)
userApi.update(id, data)
userApi.delete(id)
userApi.updateStatus(id, status)
```

---

#### 3.4.2 用户详情 (`/users/:id`)

**文件路径**: `src/pages/Users/UserDetail.tsx`

##### 功能特性
- ✅ 用户详细信息
- ✅ 订单历史
- ✅ 评价记录
- ✅ 支付记录

##### 页面结构
```
顶部: 面包屑导航 + 操作按钮
主体:
  - 基本信息卡片
  - 订单历史卡片
  - 评价记录卡片
  - 支付记录卡片
```

##### 组件依赖
- Card (卡片)
- Description (描述列表)
- Tabs (标签页)
- Table (表格)
- Avatar (头像)
- Tag (标签)

##### API调用
```typescript
userApi.get(id)
userApi.getOrders(id)
userApi.getReviews(id)
userApi.getPayments(id)
```

---

#### 3.4.3 用户表单弹窗 (`/users/UserFormModal.tsx`)

**文件路径**: `src/pages/Users/UserFormModal.tsx`

##### 功能特性
- ✅ 创建用户
- ✅ 编辑用户
- ✅ 密码设置
- ✅ 角色分配

##### 组件依赖
- Modal (弹窗)
- Form (表单)
- Input (输入框)
- Select (选择器)
- Radio (单选框)
- Switch (开关)
- Button (按钮)

---

### 3.5 支付管理

#### 3.5.1 支付列表 (`/payments`)

**文件路径**: `src/pages/Payments/PaymentList.tsx`

##### 功能特性
- ✅ 支付列表展示
- ✅ 搜索筛选
- ✅ 状态筛选
- ✅ 时间范围筛选

##### 组件依赖
- DataTable (数据表格)
- Tag (标签)
- Statistic (统计)
- Select (选择器)
- DatePicker (日期选择器)

##### API调用
```typescript
paymentApi.list(params)
```

---

#### 3.5.2 支付详情 (`/payments/:id`)

**文件路径**: `src/pages/Payments/PaymentDetailPage.tsx`

##### 功能特性
- ✅ 支付详细信息
- ✅ 交易记录
- ✅ 退款处理

##### 组件依赖
- Card (卡片)
- Description (描述列表)
- Timeline (时间线)
- Button (按钮)
- Tag (标签)

##### API调用
```typescript
paymentApi.get(id)
paymentApi.refund(id, amount)
```

---

### 3.6 陪玩师管理

#### 3.6.1 陪玩师列表 (`/players`)

**文件路径**: `src/pages/Players/PlayerList.tsx`

##### 功能特性
- ✅ 陪玩师列表展示
- ✅ 搜索筛选
- ✅ 技能标签
- ✅ 验证状态
- ✅ 在线状态

##### 组件依赖
- DataTable (数据表格)
- Tag (标签)
- Badge (徽章)
- Avatar (头像)
- Switch (开关)
- Button (按钮)
- Modal (弹窗)

##### API调用
```typescript
playerApi.list(params)
playerApi.update(id, data)
playerApi.updateStatus(id, status)
```

---

#### 3.6.2 陪玩师表单弹窗 (`/players/PlayerFormModal.tsx`)

**文件路径**: `src/pages/Players/PlayerFormModal.tsx`

##### 功能特性
- ✅ 创建陪玩师
- ✅ 编辑陪玩师
- ✅ 技能标签
- ✅ 头像上传

##### 组件依赖
- Modal (弹窗)
- Form (表单)
- Input (输入框)
- Select (选择器)
- Tag (标签)
- Upload (上传)
- Button (按钮)

---

### 3.7 评价管理

#### 3.7.1 评价列表 (`/reviews`)

**文件路径**: `src/pages/Reviews/ReviewList.tsx`

##### 功能特性
- ✅ 评价列表展示
- ✅ 搜索筛选
- ✅ 评分筛选
- ✅ 审核功能

##### 组件依赖
- DataTable (数据表格)
- Rating (评分)
- Tag (标签)
- Button (按钮)
- Modal (弹窗)

##### API调用
```typescript
reviewApi.list(params)
reviewApi.approve(id)
reviewApi.reject(id)
```

---

#### 3.7.2 评价表单弹窗 (`/reviews/ReviewFormModal.tsx`)

**文件路径**: `src/pages/Reviews/ReviewFormModal.tsx`

##### 功能特性
- ✅ 查看评价详情
- ✅ 评价审核
- ✅ 回复评价

##### 组件依赖
- Modal (弹窗)
- Form (表单)
- Rating (评分)
- Input (输入框)
- Button (按钮)

---

### 3.8 报表统计

#### 3.8.1 报表仪表盘 (`/reports`)

**文件路径**: `src/pages/Reports/ReportDashboard.tsx`

##### 功能特性
- ✅ 数据统计
- ✅ 图表分析
- ✅ 趋势预测
- ✅ 导出报表
- ✅ 时间范围选择

##### 展示内容
- 收入趋势图
- 订单趋势图
- 用户增长图
- 游戏排行榜
- 陪玩师排行榜

##### 组件依赖
- Card (卡片)
- Chart (图表)
- Statistic (统计)
- DatePicker (日期选择器)
- Select (选择器)
- Button (按钮)

##### API调用
```typescript
reportApi.getRevenueTrend(params)
reportApi.getOrderTrend(params)
reportApi.getUserGrowth(params)
reportApi.getTopGames(params)
reportApi.getTopPlayers(params)
```

---

### 3.9 权限管理

#### 3.9.1 权限列表 (`/permissions`)

**文件路径**: `src/pages/Permissions/PermissionList.tsx`

##### 功能特性
- ✅ 权限列表展示
- ✅ 角色分配
- ✅ 权限继承
- ✅ 搜索筛选

##### 组件依赖
- Tree (树形控件)
- Table (表格)
- Button (按钮)
- Modal (弹窗)
- Checkbox (复选框)

##### API调用
```typescript
permissionApi.list()
permissionApi.assignToRole(roleId, permissionIds)
```

---

### 3.10 系统设置

#### 3.10.1 设置仪表盘 (`/settings`)

**文件路径**: `src/pages/Settings/SettingsDashboard.tsx`

##### 功能特性
- ✅ 系统配置
- ✅ 用户偏好
- ✅ 主题切换
- ✅ 语言设置

##### 设置项
- 基本设置
- 安全设置
- 通知设置
- 主题设置

##### 组件依赖
- Tabs (标签页)
- Form (表单)
- Switch (开关)
- Select (选择器)
- Button (按钮)

##### API调用
```typescript
settingsApi.get()
settingsApi.update(data)
```

---

## 4. 演示页面

### 4.1 组件演示 (`/showcase`)

**文件路径**: `src/pages/ComponentsDemo/ComponentsDemo.tsx`

#### 功能特性
- ✅ 展示所有组件
- ✅ 代码示例
- ✅ Props演示
- ✅ 交互演示

#### 展示内容
- 按钮组件演示
- 输入框组件演示
- 表格组件演示
- 弹窗组件演示
- 表单组件演示

#### 组件依赖
- 所有基础组件

#### 路由守卫
- **无需认证**

---

### 4.2 缓存演示 (`/cache-demo`)

**文件路径**: `src/pages/CacheDemo/index.tsx`

#### 子页面
- `/cache-demo` - 重定向到A页面
- `/cache-demo/a` - 页面A
- `/cache-demo/b` - 页面B

#### 功能特性
- ✅ 路由缓存演示
- ✅ 状态保持演示
- ✅ 页面切换动画

#### 组件依赖
- RouteCache (路由缓存)
- Button (按钮)
- Input (输入框)

#### 路由守卫
- **无需认证**

---

## 5. 页面关联关系

### 5.1 列表 → 详情

```
用户列表 (/users) → 用户详情 (/users/:id)
订单列表 (/orders) → 订单详情 (/orders/:id)
游戏列表 (/games) → 游戏详情 (/games/:id)
支付列表 (/payments) → 支付详情 (/payments/:id)
```

### 5.2 列表 → 表单弹窗

```
用户列表 → 用户表单弹窗
订单列表 → 订单表单弹窗
游戏列表 → 游戏表单弹窗
陪玩师列表 → 陪玩师表单弹窗
评价列表 → 评价表单弹窗
```

### 5.3 面包屑导航

```
/users → 面包屑: 首页 / 用户管理
/users/123 → 面包屑: 首页 / 用户管理 / 用户详情
/orders → 面包屑: 首页 / 订单管理
/orders/456 → 面包屑: 首页 / 订单管理 / 订单详情
```

---

## 6. 路由配置

### 6.1 完整路由列表

```typescript
// 公开路由
{ path: '/login', element: <Login /> }
{ path: '/register', element: <Register /> }

// 演示路由
{ path: '/showcase', element: <ComponentsDemo /> }
{
  path: '/cache-demo',
  children: [
    { index: true, element: <Navigate to="/cache-demo/a" replace /> },
    { path: 'a', element: <CachePageA /> },
    { path: 'b', element: <CachePageB /> },
  ],
}

// 受保护路由
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
}
```

### 6.2 路由守卫

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

---

## 7. 组件依赖

### 7.1 基础组件

| 组件 | 用途 | 页面使用频率 |
|------|------|--------------|
| Button | 操作按钮 | 高 |
| Input | 输入框 | 高 |
| Select | 选择器 | 高 |
| Modal | 弹窗 | 高 |
| Table | 表格 | 高 |
| Card | 卡片 | 高 |
| Form | 表单 | 高 |
| DatePicker | 日期选择器 | 中 |
| Upload | 上传 | 中 |
| Image | 图片 | 中 |
| Tag | 标签 | 中 |
| Avatar | 头像 | 中 |
| Badge | 徽章 | 中 |
| Rating | 评分 | 低 |
| Switch | 开关 | 低 |
| Checkbox | 复选框 | 低 |
| Radio | 单选框 | 低 |

### 7.2 布局组件

| 组件 | 用途 | 页面使用频率 |
|------|------|--------------|
| Layout | 布局容器 | 高 |
| Header | 顶部导航 | 高 |
| Sidebar | 侧边栏 | 高 |
| Breadcrumb | 面包屑 | 高 |
| Tabs | 标签页 | 中 |

### 7.3 数据组件

| 组件 | 用途 | 页面使用频率 |
|------|------|--------------|
| DataTable | 数据表格 | 高 |
| Table | 表格 | 中 |
| Pagination | 分页 | 高 |
| Statistic | 统计 | 中 |
| Chart | 图表 | 中 |
| Description | 描述列表 | 中 |
| Timeline | 时间线 | 低 |
| Steps | 步骤条 | 低 |

### 7.4 反馈组件

| 组件 | 用途 | 页面使用频率 |
|------|------|--------------|
| Message | 消息提示 | 高 |
| Notification | 通知 | 中 |
| Modal | 模态框 | 高 |
| Popconfirm | 气泡确认 | 中 |
| Spin | 加载 | 高 |
| Empty | 空状态 | 中 |

---

## 8. 页面开发规范

### 8.1 页面文件结构

```
pages/
└── ModuleName/
    ├── index.ts                 # 导出文件
    ├── PageName.tsx             # 主页面组件
    ├── PageName.less            # 页面样式
    ├── PageName.test.tsx        # 测试文件
    ├── hooks/                   # 页面专属Hook
    │   └── usePageName.ts
    ├── components/              # 页面专属组件
    │   └── ComponentName/
    ├── utils/                   # 页面专属工具
    │   └── pageUtils.ts
    └── types/                   # 页面专属类型
        └── pageName.ts
```

### 8.2 页面组件规范

#### 8.2.1 页面结构

```tsx
import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Table, Button } from '@/components'
import { pageNameApi } from '@/api/pageName'

interface Props {}

const PageName: React.FC<Props> = () => {
  // 状态管理
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(false)

  // 路由参数
  const { id } = useParams()
  const navigate = useNavigate()

  // 生命周期
  useEffect(() => {
    fetchData()
  }, [id])

  // 数据获取
  const fetchData = async () => {
    try {
      setLoading(true)
      const result = await pageNameApi.get(id)
      setData(result)
    } catch (error) {
      // 错误处理
    } finally {
      setLoading(false)
    }
  }

  // 事件处理
  const handleAction = () => {
    // 处理逻辑
  }

  // 渲染
  return (
    <div className="page-name">
      <Card>
        {/* 页面内容 */}
      </Card>
    </div>
  )
}

export default PageName
```

#### 8.2.2 页面命名规范

- **PascalCase**: `UserList.tsx`, `OrderDetail.tsx`
- **描述性**: 名称应能准确描述页面功能
- **统一性**: 保持命名风格一致

### 8.3 样式规范

#### 8.3.1 样式文件

```less
// PageName.less
.page-name {
  // 页面容器
  .page-header {
    // 页面头部
    margin-bottom: @spacing-lg;
  }

  .page-content {
    // 页面内容
    .content-section {
      // 内容区块
      margin-bottom: @spacing-lg;
    }
  }
}
```

#### 8.3.2 BEM命名

```css
.page-name {}
.page-name__header {}
.page-name__content {}
.page-name--highlighted {}
```

### 8.4 状态管理

#### 8.4.1 本地状态

```tsx
const [data, setData] = useState([])
const [loading, setLoading] = useState(false)
const [pagination, setPagination] = useState({
  current: 1,
  pageSize: 10,
  total: 0,
})
```

#### 8.4.2 全局状态 (Context)

```tsx
// 使用AuthContext
const { user, logout } = useAuth()

// 使用ThemeContext
const { theme, toggleTheme } = useTheme()
```

### 8.5 错误处理

#### 8.5.1 API错误处理

```tsx
try {
  const result = await apiCall()
} catch (error) {
  message.error(error.message)
  // 或
  notification.error({
    message: '操作失败',
    description: error.message,
  })
}
```

#### 8.5.2 边界错误处理

```tsx
class PageErrorBoundary extends React.Component {
  // 错误边界组件
}
```

---

## 📚 相关文档

- [前端开发完整指南](./FRONTEND_DEVELOPMENT_COMPLETE_GUIDE.md)
- [组件库文档](./组件库文档.md)
- [开发指南](./DEVELOPER_GUIDE.md)
- [技术文档](./TECHNICAL_DOCUMENTATION.md)

---

**文档维护者**: GameLink Frontend Team
**最后更新**: 2025-10-31
**版本**: v1.0
**页面总数**: 22个
**完成度**: 100%
