# GameLink 项目文档

> 📚 **最后更新**: 2026-01-21
>
> 开发人员入口 - 所有文档位于 `.kiro/steering/`

---

## 快速开始

| 我要... | 查看文档 |
|---------|----------|
| 🚀 **快速了解项目** | [quickstart.md](../../.kiro/steering/quickstart.md) |
| 📱 **开发小程序** | [miniprogram-development-plan.md](../../.kiro/steering/client/miniprogram-development-plan.md) ⭐ |
| 📋 **理解业务规则** | [01-product.md](../../.kiro/steering/01-product.md) |
| 🏗️ **开发新功能** | [03-project-structure.md](../../.kiro/steering/03-project-structure.md) → [04-data-models.md](../../.kiro/steering/04-data-models.md) |
| 🧪 **编写测试** | [05-testing-standard.md](../../.kiro/steering/05-testing-standard.md) |
| 📊 **查看进度** | [06-progress.md](../../.kiro/steering/06-progress.md) |

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

---

## 📱 小程序开发

### 开发指南 ⭐

| 文件 | 说明 |
|------|------|
| **[miniprogram-development-plan.md](../../.kiro/steering/client/miniprogram-development-plan.md)** | **小程序开发计划** - 技术栈、架构、开发路线图 |
| [miniprogram-design.md](../../.kiro/steering/client/miniprogram-design.md) | 小程序 UI/UX 设计规范（Discord Dark Theme） |
| [apps-roadmap.md](../../.kiro/steering/client/apps-roadmap.md) | 客户端产品规划（用户流程、页面设计） |

### 业务流程指南 (client/) ⭐ **前端必读**

| 文件 | 说明 |
|------|------|
| [09-auth-flow-guide.md](../../.kiro/steering/client/09-auth-flow-guide.md) | **用户认证流程** - 登录/注册/Token管理 |
| [10-player-flow-guide.md](../../.kiro/steering/client/10-player-flow-guide.md) | **陪玩师流程** - 申请/接单/收益/提现 |
| [11-payment-wallet-guide.md](../../.kiro/steering/client/11-payment-wallet-guide.md) | **支付钱包流程** - 支付/充值/退款 |
| [12-vip-marketing-guide.md](../../.kiro/steering/client/12-vip-marketing-guide.md) | **VIP营销流程** - VIP/优惠券/活动/推荐 |
| [13-chat-review-dispute-guide.md](../../.kiro/steering/client/13-chat-review-dispute-guide.md) | **聊天评价争议** - WebSocket/评价/争议/通知 |

### Web 端参考文档 (client/)

| 文件 | 说明 |
|------|------|
| [web-pages.md](../../.kiro/steering/client/web-pages.md) | Web 页面规范 |
| [web-design-system.md](../../.kiro/steering/client/web-design-system.md) | 设计系统 |
| [web-color-palette.md](../../.kiro/steering/client/web-color-palette.md) | 配色方案 |
| [web-pwa-checklist.md](../../.kiro/steering/client/web-pwa-checklist.md) | PWA 检查清单 |

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
    ├── miniprogram-development-plan.md  # 小程序开发计划 ⭐
    ├── miniprogram-design.md             # 小程序设计规范
    ├── apps-roadmap.md                  # 产品规划
    ├── 09-13 业务流程指南               # 认证/陪玩师/支付/营销/聊天
    └── web-*                             # Web 端参考文档
```

---

## 文档维护规则

| 文档 | 更新时机 | 负责人 |
|------|----------|--------|
| 04-data-models.md | **每次模型变更后**（强制） | Backend Lead |
| 04c-enums-indexes.md | 每次枚举/索引变更 | Backend Lead |
| 06-progress.md | 每周或模块完成时 | Project Lead |
| 01-product.md | 每次功能上线后 | PM |
| miniprogram-development-plan.md | 小程序开发里程碑完成时 | Frontend Lead |

---

## 命名规范速查

### Go 后端

| 类型 | 规范 | 示例 |
|------|------|------|
| 文件 | camelCase | `userService.go` |
| 类型 | PascalCase | `UserService` |
| 导出函数 | PascalCase | `CreateOrder` |
| 私有函数 | camelCase | `validateInput` |

### 前端/小程序

| 类型 | 规范 | 示例 |
|------|------|------|
| 组件 | kebab-case | `gl-button/` |
| 页面 | kebab-case | `pages/player-detail/` |
| TS 文件 | camelCase | `request.ts` |
| Less 文件 | kebab-case | `variables.less` |
| 类名 | kebab-case + BEM | `.player-card__name` |
| 变量 | camelCase | `userInfo` |
| 常量 | UPPER_SNAKE_CASE | `API_BASE_URL` |

### 文档文件

| 类型 | 规范 | 示例 |
|------|------|------|
| 所有文档 | kebab-case | `miniprogram-development-plan.md` |

---

## 小程序快速开始

### 技术栈
- **框架**: 原生微信小程序
- **语言**: TypeScript
- **样式**: Less
- **渲染**: Skyline + glass-easel

### 目录结构
```
app/miniprogram/
├── pages/              # 页面
│   ├── index/          # 首页（陪玩师列表）
│   ├── category/       # 游戏分类
│   ├── message/        # 消息列表
│   ├── profile/        # 个人中心
│   ├── player/         # 陪玩师详情
│   └── order/          # 订单相关
├── components/         # 组件库
│   ├── gl-*            # 基础组件
│   ├── player-card     # 陪玩师卡片
│   └── game-card       # 游戏卡片
├── utils/              # 工具函数
│   ├── request.ts      # API 请求
│   ├── auth.ts         # 认证工具
│   └── storage.ts      # 存储工具
└── styles/             # 全局样式
    └── variables.less  # Discord 风格变量
```

### 开发命令
```bash
# 安装依赖
npm install

# 开发模式（微信开发者工具）
# 1. 打开微信开发者工具
# 2. 导入项目，选择 app/miniprogram 目录
# 3. 开启"不校验合法域名"
```

---

> ⚠️ **重要**: 所有代码修改前，必须先查阅相关文档，确保理解业务规则和技术规范。
