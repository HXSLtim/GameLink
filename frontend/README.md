# GameLink Frontend

GameLink 是一个现代化的游戏陪玩管理平台前端，采用 React + TypeScript + Vite 技术栈。

## 技术栈

- **框架**: React 18 + TypeScript
- **构建工具**: Vite
- **路由**: React Router v6
- **状态管理**: Zustand
- **UI 组件**: 自研组件库
- **样式**: Less
- **测试**: Vitest
- **代码规范**: ESLint + Prettier + Husky

## 项目结构

```
src/
├── api/                      # API 层（三层架构）
│   ├── interface/           # 类型定义层
│   ├── modules/             # 业务逻辑层
│   ├── request/             # HTTP 请求层
│   ├── client.ts            # API 客户端配置
│   └── index.ts             # 统一导出
│
├── apps/                     # 三端应用
│   ├── admin/               # 管理后台
│   ├── player/              # 陪玩师端
│   └── user/                # 用户端
│
├── router/                   # 路由配置
│   ├── index.tsx            # 路由入口
│   ├── routes.tsx           # 整合路由配置
│   └── ProtectedRoute.tsx   # 路由守卫
│
├── shared/                   # 共享资源
│   ├── components/          # 通用组件
│   ├── hooks/               # 通用 hooks
│   ├── utils/               # 工具函数
│   ├── types/               # 通用类型
│   ├── styles/              # 全局样式
│   └── theme/               # 主题配置
│
└── stores/                   # 状态管理
```

## 快速开始

### 安装依赖

```bash
npm install
```

### 开发环境

```bash
npm run dev
```

### 构建生产版本

```bash
npm run build
```

### 代码检查

```bash
# 运行 ESLint
npm run lint

# 自动修复 ESLint 问题
npm run lint:fix

# 运行 Prettier
npm run format

# 检查 Prettier 格式
npm run format:check

# 运行 TypeScript 类型检查
npm run typecheck
```

### 测试

```bash
# 运行测试
npm test

# 运行测试一次
npm run test:run

# 运行测试覆盖率
npm run test:coverage
```

### Git 提交

项目使用 Husky + lint-staged + commitlint 进行代码提交规范检查。

```bash
# 使用 commitizen 进行规范提交
npm run commit
```

提交信息格式遵循 [Conventional Commits](https://conventionalcommits.org/) 规范：

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

### 代码验证

在提交代码前，建议运行完整的验证流程：

```bash
npm run validate
```

这会依次运行：类型检查 → ESLint → Prettier 格式检查 → 单元测试

## 环境变量

复制 `.env.example` 文件为 `.env` 并配置相关变量：

```bash
cp .env.example .env
```

## 别名路径

项目配置了以下路径别名：

- `@` → `src/`
- `@api` → `src/api/`
- `@apps` → `src/apps/`
- `@router` → `src/router/`
- `@shared` → `src/shared/`
- `@stores` → `src/stores/`

## 代码规范

- 使用 ESLint 进行代码检查
- 使用 Prettier 进行代码格式化
- 使用 Husky 在提交前自动检查代码
- 使用 Conventional Commits 规范提交信息

## 浏览器支持

- Chrome >= 87
- Firefox >= 78
- Safari >= 14
- Edge >= 88