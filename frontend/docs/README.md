# 📚 GameLink 前端项目文档中心

## 🎯 快速导航

| 文档类型    | 数量 | 路径                                   |
| ----------- | ---- | -------------------------------------- |
| 🔐 加密文档 | 5    | [crypto/](./crypto/)                   |
| 📡 API文档  | 3    | [api/](./api/)                         |
| 🎨 设计规范 | 3    | [design/](./design/)                   |
| 📖 开发指南 | 3    | [guides/](./guides/)                   |
| ⚡ 功能文档 | 6    | [features/](./features/)               |
| 🔧 重构文档 | 4    | [refactoring/](./refactoring/)         |
| 📊 历史报告 | 21   | [archive/reports/](./archive/reports/) |

---

## 🔐 加密中间件（最新）

> **推荐阅读顺序**: README → INTEGRATION → MIDDLEWARE → EXAMPLES → ENV_CONFIG

- **[README.md](./crypto/README.md)** ⭐ 快速入门（5分钟）
- **[INTEGRATION.md](./crypto/INTEGRATION.md)** 完整集成指南
- **[MIDDLEWARE.md](./crypto/MIDDLEWARE.md)** 技术实现细节
- **[EXAMPLES.md](./crypto/EXAMPLES.md)** 代码使用示例
- **[ENV_CONFIG.md](./crypto/ENV_CONFIG.md)** 环境变量配置

### 快速开始

```typescript
// 自动加密，无需修改代码
await authApi.login({ username: 'admin', password: '123456' });
```

---

## 📡 API 文档

- **[API_REQUIREMENTS.md](./api/API_REQUIREMENTS.md)** API 需求规范
- **[INTEGRATION.md](./api/INTEGRATION.md)** API 集成指南
- **[backend-models.md](./api/backend-models.md)** 后端接口模型

---

## 🎨 设计规范

- **[DESIGN_SYSTEM.md](./design/DESIGN_SYSTEM.md)** 设计系统规范
- **[DESIGN_SYSTEM_V2.md](./design/DESIGN_SYSTEM_V2.md)** 设计系统 V2
- **[CODING_STANDARDS.md](./design/CODING_STANDARDS.md)** 编码规范

---

## 📖 开发指南

- **[QUICK_START.md](./guides/QUICK_START.md)** 快速开始
- **[MIGRATION_GUIDE.md](./guides/MIGRATION_GUIDE.md)** 迁移指南
- **[MVP_DEVELOPMENT_PLAN.md](./guides/MVP_DEVELOPMENT_PLAN.md)** MVP 开发计划

---

## ⚡ 功能文档

### 主题相关

- **[THEME_RIPPLE_EFFECT.md](./features/THEME_RIPPLE_EFFECT.md)** 主题涟漪效果
- **[THEME_TOGGLE_GUIDE.md](./features/THEME_TOGGLE_GUIDE.md)** 主题切换指南

### 其他功能

- **[NAVIGATION_SYSTEM.md](./features/NAVIGATION_SYSTEM.md)** 导航系统
- **[I18N_IMPLEMENTATION_GUIDE.md](./features/I18N_IMPLEMENTATION_GUIDE.md)** 国际化实现
- **[FIGMA_TO_CODE_GUIDE.md](./features/FIGMA_TO_CODE_GUIDE.md)** Figma 转代码
- **[MOCK_LOGIN_GUIDE.md](./features/MOCK_LOGIN_GUIDE.md)** Mock 登录

---

## 🔧 重构文档

- **[COLOR_VARIABLES_REFACTOR.md](./refactoring/COLOR_VARIABLES_REFACTOR.md)** 颜色变量重构
- **[IMPORT_PATH_GUIDE.md](./refactoring/IMPORT_PATH_GUIDE.md)** 导入路径指南
- **[PATH_ALIAS_FIX.md](./refactoring/PATH_ALIAS_FIX.md)** 路径别名修复
- **[HMR_AND_THEME_FIX.md](./refactoring/HMR_AND_THEME_FIX.md)** HMR 和主题修复

---

## 📊 历史报告（归档）

查看 [archive/reports/](./archive/reports/) 目录，包含 21 个历史报告文档。

---

## 🎯 推荐阅读路径

### 新成员入职（1小时）

1. [快速开始](./guides/QUICK_START.md) - 10分钟
2. [设计系统](./design/DESIGN_SYSTEM.md) - 20分钟
3. [编码规范](./design/CODING_STANDARDS.md) - 15分钟
4. [加密快速入门](./crypto/README.md) - 5分钟
5. [API 需求](./api/API_REQUIREMENTS.md) - 10分钟

### 功能开发（2小时）

1. [API 集成指南](./api/INTEGRATION.md)
2. [加密集成指南](./crypto/INTEGRATION.md)
3. [主题切换指南](./features/THEME_TOGGLE_GUIDE.md)
4. [国际化实现](./features/I18N_IMPLEMENTATION_GUIDE.md)

### 深入学习（半天）

1. 阅读所有设计规范
2. 阅读所有 API 文档
3. 阅读所有加密文档
4. 阅读重点功能文档

---

## 🔍 按主题查找

### 加密相关

- 入门: [crypto/README.md](./crypto/README.md)
- 集成: [crypto/INTEGRATION.md](./crypto/INTEGRATION.md)
- 配置: [crypto/ENV_CONFIG.md](./crypto/ENV_CONFIG.md)

### 样式相关

- 设计系统: [design/DESIGN_SYSTEM.md](./design/DESIGN_SYSTEM.md)
- 主题切换: [features/THEME_TOGGLE_GUIDE.md](./features/THEME_TOGGLE_GUIDE.md)
- 颜色变量: [refactoring/COLOR_VARIABLES_REFACTOR.md](./refactoring/COLOR_VARIABLES_REFACTOR.md)

### 路由相关

- 导航系统: [features/NAVIGATION_SYSTEM.md](./features/NAVIGATION_SYSTEM.md)
- 路径别名: [refactoring/IMPORT_PATH_GUIDE.md](./refactoring/IMPORT_PATH_GUIDE.md)

### 国际化

- 实现指南: [features/I18N_IMPLEMENTATION_GUIDE.md](./features/I18N_IMPLEMENTATION_GUIDE.md)

---

## 📂 文档目录结构

```
docs/
├── README.md              # 本文档
├── crypto/               # 加密文档 (5个)
├── api/                  # API文档 (3个)
├── design/               # 设计规范 (3个)
├── guides/               # 开发指南 (3个)
├── features/             # 功能文档 (6个)
├── refactoring/          # 重构文档 (4个)
└── archive/              # 历史归档
    └── reports/          # 历史报告 (21个)
```

---

## 🆕 最新更新

- **2025-01-28** - 完成文档归档整理
- **2025-01-28** - 添加加密中间件文档
- **2025-01-28** - 添加环境变量配置文档

---

## 📝 文档贡献

发现文档问题或有改进建议？

1. 直接修改对应文档
2. 提交 Pull Request
3. 更新此索引（如有新文档）

---

## 🔗 外部资源

- [React 官方文档](https://react.dev/)
- [TypeScript 官方文档](https://www.typescriptlang.org/)
- [Vite 官方文档](https://vitejs.dev/)
- [Arco Design](https://arco.design/)

---

**最后更新**: 2025-01-28  
**文档总数**: 45+  
**维护者**: GameLink Team
