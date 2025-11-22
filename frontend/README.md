# GameLink Frontend

GameLink 是一个现代化的游戏陪玩管理平台前端，采用 React + TypeScript + Vite 技术栈。

## 技术栈

- **框架**: React 18 + TypeScript 5.6
- **构建工具**: Vite 5
- **路由**: React Router v6
- **状态管理**: Zustand
- **样式系统**: Tailwind CSS
- **UI 组件**: 自研组件库
- **测试**: Vitest + Testing Library
- **代码规范**: ESLint + Prettier + Husky
- **HTTP 客户端**: Axios
- **加密**: Crypto-JS

## 项目结构

```
frontend/
├── .env.development          # 开发环境配置
├── .env.production           # 生产环境配置
├── .env.example              # 环境变量模板
├── vite.config.ts            # Vite 配置文件
├── tailwind.config.ts        # Tailwind CSS 配置
├── postcss.config.ts         # PostCSS 配置
├── tsconfig.json             # TypeScript 配置
├── package.json              # 依赖配置
├── .gitignore                # Git 忽略文件
├── commitlint.config.js      # 提交规范配置
├── lint-staged.config.js     # lint-staged 配置
├── .eslintrc.json            # ESLint 配置
├── .prettierrc               # Prettier 配置
├── vitest.config.ts          # Vitest 测试配置
├── docs/                     # 项目文档
│   ├── ENVIRONMENT_VARIABLES.md   # 环境变量使用指南
│   └── STYLE_AND_STATE_MANAGEMENT.md # 样式和状态管理指南
│
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
│   ├── LazyRoutes.tsx       # 路由懒加载配置
│   └── routes.tsx           # 整合路由配置
│
├── shared/                   # 共享资源
│   ├── components/          # 通用组件
│   ├── hooks/               # 通用 hooks
│   ├── utils/               # 工具函数
│   ├── types/               # 通用类型
│   └── theme/               # 主题配置
│
├── stores/                   # 状态管理
│   ├── index.ts             # 统一导出
│   ├── useAppStore.ts       # 应用全局状态
│   └── useAuthStore.ts      # 认证状态
│
├── styles/                   # 全局样式
│   └── tailwind.css         # Tailwind CSS 样式
│
├── App.tsx                   # 应用根组件
├── main.tsx                  # 应用入口
├── config.ts                 # 应用配置
└── setupTests.ts             # 测试配置
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

项目使用多环境配置，支持开发环境和生产环境：

```bash
# 开发环境（npm run dev 自动加载 .env.development）
VITE_API_BASE_URL=http://localhost:8080
VITE_CRYPTO_ENABLED=false

# 生产环境（npm run build 自动加载 .env.production）
VITE_API_BASE_URL=https://api.gamelink.com
VITE_CRYPTO_ENABLED=true
```

如需自定义本地配置，可以创建 `.env.local` 文件（不会被 Git 追踪）：

```bash
# 从模板创建（可选）
cp .env.example .env.local
```

详细的环境变量说明请查看：[环境变量配置指南](./docs/ENVIRONMENT_VARIABLES.md)

## 别名路径

项目配置了以下路径别名：

- `@` → `src/`
- `@api` → `src/api/`
- `@apps` → `src/apps/`
- `@router` → `src/router/`
- `@shared` → `src/shared/`
- `@stores` → `src/stores/`

## 样式系统与状态管理

项目集成了现代化的样式系统和状态管理方案：

- **Tailwind CSS**: 原子化样式系统，提供丰富的实用工具类
- **Zustand**: 轻量级、高性能的状态管理库

详细使用指南请查看：[样式系统与状态管理](./docs/STYLE_AND_STATE_MANAGEMENT.md)

## 代码规范

- 使用 ESLint 进行代码检查
- 使用 Prettier 进行代码格式化
- 使用 Husky 在提交前自动检查代码
- 使用 Conventional Commits 规范提交信息
- 使用 TypeScript 严格模式

## 浏览器支持

- Chrome >= 87
- Firefox >= 78
- Safari >= 14
- Edge >= 88

## 项目文档

- [环境变量配置指南](./docs/ENVIRONMENT_VARIABLES.md) - 详细的环境变量说明和使用场景
- [样式系统与状态管理](./docs/STYLE_AND_STATE_MANAGEMENT.md) - Tailwind CSS 和 Zustand 使用指南

## 开发指南

### Bundle 分析

```bash
# 分析生产构建的包大小
npm run build:analyze
```

### 代理配置

开发服务器已配置后端 API 代理（端口 8080）：

```
/api/*  -> http://localhost:8080
```

## 项目维护

项目遵循以下最佳实践：

1. **代码质量**: 通过 Husky + lint-staged 在提交前自动检查
2. **测试覆盖**: 使用 Vitest 进行单元测试和集成测试
3. **类型安全**: TypeScript 严格模式，避免使用 any 类型
4. **环境分离**: 开发环境和生产环境使用不同的配置文件
5. **文档驱动**: 重要配置和架构决策都有文档说明