# GameLink Admin Panel

GameLink 游戏陪玩平台 - 管理后台面板

## 项目概述

GameLink 是一个现代化的游戏陪玩管理平台，连接用户与游戏陪练师进行付费游戏会话。管理后台提供全方位的运营管理功能，包括用户管理、订单处理、陪练师管理、客服聊天、营销活动等。

### 技术栈

| 类别 | 技术 | 版本 |
|------|------|------|
| 框架 | React | 19.2.0 |
| 语言 | TypeScript | 5.9+ |
| 构建 | Vite | 7.2+ |
| UI库 | Ant Design | 6.0+ |
| 状态管理 | Zustand | 5.0+ |
| 路由 | React Router | 7.9+ |
| HTTP | Axios | 1.13+ |
| 实时通信 | Socket.IO | 4.8+ |
| 图表 | Recharts | 3.5+ |
| 动画 | Framer Motion | 12.23+ |
| 测试 | Vitest + Playwright | - |

---

## 快速开始

### 前置要求

- Node.js 18+
- npm 或 pnpm

### 安装依赖

```bash
cd admin
npm install
```

### 开发模式

```bash
npm run dev
```

访问: http://localhost:5173

默认管理员账号:
- 用户名: `admin`
- 密码: `admin123`

---

## 开发指南

### 目录结构

```
admin/src/
├── api/              # API 客户端
│   ├── auth.ts      # 认证接口
│   ├── admin.ts     # 管理员接口
│   └── client.ts    # 客户端接口
├── components/      # 公共组件
│   ├── PermissionGuard.tsx  # 权限守卫
│   ├── SearchTable.tsx      # 搜索表格
│   ├── PageContainer.tsx    # 页面容器
│   └── ...
├── pages/           # 页面组件
│   ├── admin/       # 管理员页面
│   ├── player/      # 陪练师页面
│   ├── user/        # 用户页面
│   └── login.tsx    # 登录页
├── stores/          # Zustand 状态管理
│   ├── authStore.ts
│   ├── chatStore.ts
│   └── ...
├── hooks/           # 自定义 Hooks
│   ├── usePermission.ts
│   ├── useWebSocket.ts
│   └── ...
├── utils/           # 工具函数
├── config/          # 配置文件
│   ├── adminRoutes.ts  # 路由配置
│   └── debug.ts        # 调试配置
└── types/           # TypeScript 类型定义
```

### 核心概念

#### 1. 权限系统

权限格式: `模块.资源.操作`

```typescript
// 权限示例
'user.list.view'      // 查看用户列表
'order.edit'          // 编辑订单
'player.rank.update'  // 更新陪练师等级
'*'                   // 超级管理员（所有权限）
```

#### 2. 权限守卫组件

```tsx
import { PermissionGuard, PermissionButton } from '@/components';

// 条件渲染子组件
<PermissionGuard permission="user.delete" mode="any">
  <Button danger>删除用户</Button>
</PermissionGuard>

// 权限按钮（无权限时隐藏）
<PermissionButton permission="order.edit">
  <Button>编辑订单</Button>
</PermissionButton>

// 无权限时禁用按钮
<PermissionButton permission="order.delete" disableOnNoPermission>
  <Button>删除订单</Button>
</PermissionButton>
```

#### 3. 搜索表格组件

```tsx
import { SearchTable } from '@/components';

<SearchTable
  columns={columns}
  dataSource={data}
  loading={loading}
  searchFields={[
    { name: 'keyword', label: '搜索', type: 'input' },
    { name: 'status', label: '状态', type: 'select', options: statusOptions },
  ]}
  onSearch={handleSearch}
  onRefresh={fetchData}
/>
```

#### 4. 状态管理 (Zustand)

```typescript
import { useAuthStore } from '@/stores/modules/authStore';

// 组件中使用
const { user, permissions, login, logout } = useAuthStore();
```

---

## 测试

### 单元测试

```bash
# 运行测试 (watch 模式)
npm run test

# 运行一次
npm run test:run

# UI 模式
npm run test:ui
```

### E2E 测试

```bash
# 安装 Playwright 浏览器
npm run test:e2e:install

# 运行 E2E 测试
npm run test:e2e

# 交互模式
npm run test:e2e:ui

# 调试模式
npm run test:e2e:debug

# 查看测试报告
npm run test:e2e:report
```

---

## 构建

### 开发构建

```bash
npm run build
```

输出: `dist/`

### Bundle 分析

```bash
npm run build:analyze
```

生成可视化报告，帮助分析包体积。

---

## 部署

### 环境变量

创建 `.env.production`:

```bash
# API 地址
VITE_API_BASE_URL=https://api.gamelink.com

# WebSocket 地址
VITE_WS_URL=wss://api.gamelink.com

# 加密配置 (生产环境必填)
VITE_CRYPTO_SECRET_KEY=your-32-char-secret-key
VITE_CRYPTO_IV=your-16-char-iv
```

### Docker 部署

```bash
# 构建镜像
docker build -t gamelink-admin .

# 运行容器
docker run -p 80:80 gamelink-admin
```

### Nginx 配置

```nginx
server {
    listen 80;
    root /usr/share/nginx/html;
    index index.html;

    # SPA 路由支持
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API 代理
    location /api/ {
        proxy_pass http://backend:8080;
    }
}
```

---

## 脚本命令

| 命令 | 说明 |
|------|------|
| `npm run dev` | 启动开发服务器 |
| `npm run build` | 生产构建 |
| `npm run build:analyze` | 构建并分析 |
| `npm run lint` | 代码检查 |
| `npm run test` | 单元测试 |
| `npm run test:e2e` | E2E 测试 |
| `npm run preview` | 预览构建 |

---

## 功能模块

| 模块 | 路径 | 说明 |
|------|------|------|
| 仪表盘 | `/admin/dashboard` | 数据概览 |
| 用户管理 | `/admin/sys/user` | 用户 CRUD |
| 陪练管理 | `/admin/player` | 陪练师管理 |
| 订单管理 | `/admin/order` | 订单处理 |
| 支付管理 | `/admin/payment` | 支付记录 |
| 纠纷管理 | `/admin/dispute` | 纠纷处理 |
| 实时监控 | `/admin/monitor` | 系统监控大屏 |
| 聊天管理 | `/admin/chat` | 聊天室和记录 |
| 佣金管理 | `/admin/commission` | 佣金规则 |
| VIP管理 | `/admin/vip` | 会员等级 |
| 优惠券 | `/admin/coupon` | 优惠券管理 |
| 活动管理 | `/admin/activity` | 营销活动 |
| 推荐管理 | `/admin/referral` | 推荐系统 |
| 团队管理 | `/admin/team` | 团队管理 |
| 提现管理 | `/admin/withdraw` | 提现审核 |
| 充值管理 | `/admin/recharge` | 充值记录 |
| 结算管理 | `/admin/settlement` | 收入结算 |
| 游戏管理 | `/admin/game` | 游戏配置 |
| 服务项 | `/admin/service` | 服务定价 |
| 路由规则 | `/admin/routing` | 订单分配 |
| 系统设置 | `/admin/sys/*` | 菜单/角色/权限 |

---

## 开发规范

### 代码风格

- 使用 ESLint + Prettier
- 遵循 TypeScript 严格模式
- 组件使用 PascalCase
- 工具函数使用 camelCase

### Git 提交规范

```
feat: 新功能
fix: 修复 bug
docs: 文档更新
refactor: 代码重构
test: 测试相关
chore: 构建/工具变更
```

### 组件开发规范

```tsx
// 1. 使用函数组件 + Hooks
// 2. Props 定义接口
interface UserCardProps {
  user: User;
  onEdit?: (id: number) => void;
}

// 3. 组件导出
export function UserCard({ user, onEdit }: UserCardProps) {
  // ...
}

// 4. 默认导出用于懒加载
export default UserCard;
```

---

## 故障排查

### 常见问题

**Q: 开发环境跨域错误**
A: 检查 `vite.config.ts` 中的 proxy 配置

**Q: 权限不生效**
A: 确认用户已登录且权限已加载，检查 PermissionGuard 配置

**Q: WebSocket 连接失败**
A: 检查 `VITE_WS_URL` 环境变量配置

**Q: 构建后页面空白**
A: 检查路由配置是否使用 HashRouter，或 Nginx 配置是否正确

---

## 相关文档

- [组件使用文档](./docs/COMPONENTS.md)
- [项目总文档](../CLAUDE.md)
- [业务规则](../.kiro/steering/QUICKSTART.md)
