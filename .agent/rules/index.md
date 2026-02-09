---
trigger: always_on
---

# GameLink 项目文档

> 📚 **最后更新**: 2026-01-15
>
> 开发人员入口 - 所有文档位于 `.kiro/steering/`

---

## 快速开始

| 我要... | 查看文档 |
|---------|----------|
| 🚀 **快速了解项目** | [quickstart.md](../../.kiro/steering/quickstart.md) |
| 📋 **理解业务规则** | [01-product.md](../../.kiro/steering/01-product.md) |
| 🏗️ **开发新功能** | [03-project-structure.md](../../.kiro/steering/03-project-structure.md) → [04-data-models.md](../../.kiro/steering/04-data-models.md) |
| 🧪 **编写测试** | [05-testing-standard.md](../../.kiro/steering/05-testing-standard.md) |
| 📊 **查看进度** | [06-progress.md](../../.kiro/steering/06-progress.md) |
| 📱 **开发客户端** | [client/web-pages.md](../../.kiro/steering/client/web-pages.md) |

---

## 文档目录

### 核心文档（必读）

| 文件 | 说明 |
|------|------|
| [quickstart.md](../../.kiro/steering/quickstart.md) | 快速参考指南（命名规范、业务规则速查） |
| [01-product.md](../../.kiro/steering/01-product.md) | 产品概述、商业模式、功能清单 |
| [02-tech-stack.md](../../.kiro/steering/02-tech-stack.md) | 技术栈、CI/CD、部署配置 |
| [03-project-structure.md](../../.kiro/steering/03-project-structure.md) | 代码结构、命名规范、目录布局 |
| [04-data-models.md](../../.kiro/steering/04-data-models.md) | **核心数据模型**（含业务逻辑） |
| [05-testing-standard.md](../../.kiro/steering/05-testing-standard.md) | 测试规范、检查清单 |
| [06-progress.md](../../.kiro/steering/06-progress.md) | 项目进度、模块状态 |

### 数据模型附录

| 文件 | 说明 |
|------|------|
| [04a-marketing-models.md](../../.kiro/steering/04a-marketing-models.md) | 营销模块（VIP/优惠券/活动） |
| [04b-team-models.md](../../.kiro/steering/04b-team-models.md) | 团队系统 |
| [04c-enums-indexes.md](../../.kiro/steering/04c-enums-indexes.md) | 枚举类型、数据库索引 |
| [04d-notification-models.md](../../.kiro/steering/04d-notification-models.md) | 通知系统 |

### 测试文档

| 文件 | 说明 |
|------|------|
| [07-integration-test-plan.md](../../.kiro/steering/07-integration-test-plan.md) | 集成测试计划 |

### 客户端文档 (client/)

| 文件 | 说明 |
|------|------|
| [client/apps-roadmap.md](../../.kiro/steering/client/apps-roadmap.md) | 客户端应用路线图 |
| [client/web-pages.md](../../.kiro/steering/client/web-pages.md) | Web 页面规范 |
| [client/web-design-system.md](../../.kiro/steering/client/web-design-system.md) | 设计系统 |
| [client/web-color-palette.md](../../.kiro/steering/client/web-color-palette.md) | 配色方案 |
| [client/web-pwa-checklist.md](../../.kiro/steering/client/web-pwa-checklist.md) | PWA 检查清单 |
| [client/miniprogram-design.md](../../.kiro/steering/client/miniprogram-design.md) | 小程序设计 |

### 辅助工具

| 文件 | 说明 |
|------|------|
| [ui-ux-pro-max.md](../../.kiro/steering/ui-ux-pro-max.md) | UI/UX 设计智能工具 |

---

## 目录结构

```
.kiro/steering/
├── 00-index.md                 # 本文件 - 开发入口
├── quickstart.md               # 快速参考
├── 01-product.md               # 产品概述
├── 02-tech-stack.md            # 技术栈
├── 03-project-structure.md     # 项目结构
├── 04-data-models.md           # 核心数据模型
├── 04a-marketing-models.md     # 营销模块模型
├── 04b-team-models.md          # 团队模块模型
├── 04c-enums-indexes.md        # 枚举与索引
├── 04d-notification-models.md  # 通知模块模型
├── 05-testing-standard.md      # 测试规范
├── 06-progress.md              # 项目进度
├── 07-integration-test-plan.md # 集成测试计划
├── ui-ux-pro-max.md            # UI/UX 工具
└── client/                     # 客户端文档
    ├── apps-roadmap.md
    ├── web-pages.md
    ├── web-design-system.md
    ├── web-color-palette.md
    ├── web-pwa-checklist.md
    └── miniprogram-design.md
```

---

## 文档维护规则

| 文档 | 更新时机 | 负责人 |
|------|----------|--------|
| 04-data-models.md | **每次模型变更后**（强制） | Backend Lead |
| 04c-enums-indexes.md | 每次枚举/索引变更 | Backend Lead |
| 06-progress.md | 每周或模块完成时 | Project Lead |
| 01-product.md | 每次功能上线后 | PM |

---

## 命名规范速查

### Go 后端

| 类型 | 规范 | 示例 |
|------|------|------|
| 文件 | camelCase | `userService.go` |
| 类型 | PascalCase | `UserService` |
| 导出函数 | PascalCase | `CreateOrder` |
| 私有函数 | camelCase | `validateInput` |

### React 前端

| 类型 | 规范 | 示例 |
|------|------|------|
| 组件 | PascalCase | `UserTable.tsx` |
| 工具函数 | camelCase | `formatDate.ts` |
| 常量 | UPPER_SNAKE_CASE | `API_BASE_URL` |

### 文档文件

| 类型 | 规范 | 示例 |
|------|------|------|
| 所有文档 | kebab-case | `web-design-system.md` |

---

> ⚠️ **重要**: 所有代码修改前，必须先查阅相关文档，确保理解业务规则和技术规范。
