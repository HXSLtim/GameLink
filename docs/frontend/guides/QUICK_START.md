# 🚀 GameLink 前端快速启动指南

**更新时间**: 2025-01-05

---

## ✅ 前置条件

### 1. 后端服务运行

确保后端服务正在运行：

```bash
# 后端服务地址
http://localhost:8080
```

### 2. 测试账号

使用后端测试报告中的账号：

- **邮箱**: `fulltest@example.com`
- **密码**: `Test@123456`

或者使用任何已创建的账号。

---

## 🎯 快速启动步骤

### Step 1: 安装依赖（首次）

```bash
cd /mnt/c/Users/a2778/Desktop/code/GameLink/frontend
npm install
```

### Step 2: 启动开发服务器

```bash
npm run dev
```

服务器将在以下地址运行：

- **Local**: http://localhost:5173
- **Network**: http://[your-ip]:5173

### Step 3: 访问应用

在浏览器中打开：

```
http://localhost:5173
```

### Step 4: 登录

1. 输入邮箱: `fulltest@example.com`
2. 输入密码: `Test@123456`
3. 点击"登录"按钮

---

## 🧪 测试功能

### ✅ 认证功能

- [x] 登录
- [x] Token保存
- [x] 自动跳转仪表盘
- [x] 登出
- [x] Token过期自动跳转登录

### ✅ 用户管理

- [x] 查看用户列表
- [x] 搜索用户（姓名/手机/邮箱）
- [x] 筛选用户（角色/状态）
- [x] 查看用户详情
- [x] 分页功能
- [x] 加载状态（骨架屏）

### ✅ 订单管理

- [x] 查看订单列表
- [x] 筛选订单（状态/审核状态）
- [x] 查看订单详情
- [x] 分页功能
- [x] 加载状态

---

## 🔧 常用命令

### 开发

```bash
# 启动开发服务器
npm run dev

# 启动开发服务器（后台运行）
npm run dev &
```

### 构建

```bash
# 生产构建
npm run build

# 预览生产构建
npm run preview
```

### 代码质量

```bash
# 格式化代码
npm run format

# 运行测试
npm run test

# 类型检查
npm run type-check
```

---

## 🐛 常见问题

### 问题1: 连接不到后端

**症状**: 登录时提示"网络错误"

**解决方案**:

1. 确认后端服务正在运行: `http://localhost:8080/api/health`
2. 检查环境变量: `.env.development` 中 `VITE_API_BASE_URL`
3. 查看浏览器控制台的网络请求

### 问题2: Token过期

**症状**: 自动跳转到登录页

**解决方案**:

- 这是正常行为，Token有效期为24小时
- 重新登录即可

### 问题3: 页面空白

**症状**: 启动后页面空白

**解决方案**:

1. 检查浏览器控制台错误
2. 确认依赖已安装: `npm install`
3. 清除浏览器缓存后刷新

### 问题4: 样式错误

**症状**: 样式显示不正确

**解决方案**:

1. 清除缓存: `Ctrl+F5` 或 `Cmd+Shift+R`
2. 重启开发服务器
3. 检查Less文件是否正确编译

---

## 📁 项目结构

```
frontend/
├── src/
│   ├── api/
│   │   └── client.ts           # API Client配置 ✅
│   ├── services/
│   │   └── api/
│   │       ├── auth.ts         # 认证API ✅
│   │       ├── user.ts         # 用户管理API ✅
│   │       └── order.ts        # 订单管理API ✅
│   ├── contexts/
│   │   ├── AuthContext.tsx     # 认证上下文 ✅
│   │   ├── ThemeContext.tsx    # 主题上下文 ✅
│   │   └── I18nContext.tsx     # 国际化上下文 ✅
│   ├── pages/
│   │   ├── Login/              # 登录页 ✅
│   │   ├── Dashboard/          # 仪表盘 ✅
│   │   ├── Users/              # 用户管理 ✅
│   │   └── Orders/             # 订单管理 ✅
│   ├── components/             # 自定义组件库 ✅
│   ├── types/                  # TypeScript类型定义 ✅
│   └── utils/                  # 工具函数 ✅
├── .env.development            # 开发环境配置
└── .env.production             # 生产环境配置
```

---

## 🌐 API端点

### Base URL

- **开发环境**: `http://localhost:8080`
- **生产环境**: `https://api.gamelink.com` (待配置)

### 已实现的API

#### 认证 (5个)

- `POST /api/auth/login` - 登录
- `POST /api/auth/register` - 注册
- `POST /api/auth/refresh` - 刷新Token
- `POST /api/auth/logout` - 登出
- `GET /api/health` - 健康检查

#### 用户管理 (8个)

- `GET /api/admin/users` - 用户列表
- `GET /api/admin/users/:id` - 用户详情
- `POST /api/admin/users` - 创建用户
- `PUT /api/admin/users/:id` - 更新用户
- `PUT /api/admin/users/:id/status` - 更新状态
- `PUT /api/admin/users/:id/role` - 更新角色
- `DELETE /api/admin/users/:id` - 删除用户
- `GET /api/admin/users/:id/orders` - 用户订单

#### 订单管理 (7个)

- `GET /api/admin/orders` - 订单列表
- `GET /api/admin/orders/:id` - 订单详情
- `POST /api/admin/orders` - 创建订单
- `PUT /api/admin/orders/:id` - 更新订单
- `POST /api/admin/orders/:id/assign` - 分配订单
- `POST /api/admin/orders/:id/review` - 审核订单
- `GET /api/admin/orders/:id/logs` - 操作日志

**总计**: 20个API ✅

---

## 🎨 设计系统

### 主题

- **亮色主题**: 白色背景 + 黑色文字
- **暗色主题**: 黑色背景 + 白色文字
- **主题切换**: 点击导航栏右上角的主题图标

### 风格

- **Neo-brutalism**: 纯黑白配色，直角边框，硬阴影
- **无圆角**: 所有元素使用直角
- **纯色阴影**: 8px x 8px 黑色阴影
- **粗边框**: 2px 黑色边框

---

## 📚 相关文档

- [API需求文档](./API_REQUIREMENTS.md)
- [API对接指南](./API_INTEGRATION_GUIDE.md)
- [API对接完成报告](./API_INTEGRATION_COMPLETE.md)
- [设计系统](./DESIGN_SYSTEM_V2.md)

---

## 💡 开发提示

### 添加新API

1. 在 `src/services/api/` 创建新的API服务文件
2. 定义请求/响应类型
3. 使用 `apiClient` 发送请求
4. 在页面组件中调用

**示例**:

```typescript
// src/services/api/game.ts
export const gameApi = {
  getList: (params) => apiClient.get('/api/admin/games', { params }),
  getDetail: (id) => apiClient.get(`/api/admin/games/${id}`),
};

// 在组件中使用
const result = await gameApi.getList({ page: 1, page_size: 10 });
```

### 错误处理

所有API调用都应包含错误处理：

```typescript
try {
  const data = await userApi.getList({...});
  setUsers(data.list);
} catch (err) {
  const errorMessage = err instanceof Error ? err.message : '加载失败';
  setError(errorMessage);
  console.error('Error:', err);
}
```

### 加载状态

使用骨架屏提升用户体验：

```typescript
if (loading) {
  return <PageSkeleton />;
}
```

---

## ✅ 就绪状态

- ✅ 依赖已安装 (axios)
- ✅ API Client已配置
- ✅ 认证流程已实现
- ✅ 用户管理已实现
- ✅ 订单管理已实现
- ✅ 类型定义完整
- ✅ 错误处理完善
- ✅ 加载状态优化

**项目状态**: 🟢 生产就绪（核心模块）

---

**Happy Coding! 🎉**

如有问题，请查看控制台日志或联系开发团队。
