# Steering 文档导航

> 📚 **最后更新**: 2025-01-06

## 文档结构

本目录包含 GameLink 项目的核心规范文档，按使用频率分为三类：

### 📖 核心文档（高频使用）

| 文档 | 用途 | 何时查看 |
|------|------|----------|
| [01-product.md](../01-product.md) | 产品概述、商业规则、功能清单 | 需求理解、业务咨询 |
| [03-project-structure.md](../03-project-structure.md) | 代码结构、命名规范、目录布局 | 新功能开发、代码审查 |
| [04-data-models.md](../04-data-models.md) | **数据模型权威参考**（含核心业务逻辑） | 数据库变更、业务逻辑实现 |
| [05-testing-standard.md](../05-testing-standard.md) | 测试规范、检查清单 | 测试编写、质量验收 |

### 🛠️ 参考文档（中频使用）

| 文档 | 用途 | 何时查看 |
|------|------|----------|
| [02-tech-stack.md](../02-tech-stack.md) | 技术选型、CI/CD 配置 | 环境搭建、部署配置 |
| [06-project-management.md](../06-project-management.md) | 模块状态、文档维护规则 | 进度查询、文档更新 |

### 📊 项目进度

| 文档 | 用途 | 何时查看 |
|------|------|----------|
| [../docs/PROGRESS.md](../docs/PROGRESS.md) | 项目进度追踪 | 版本发布、进度汇报 |

### 📋 附录文档（低频使用）

| 文档 | 用途 | 何时查看 |
|------|------|----------|
| [04a-marketing-models.md](../04a-marketing-models.md) | 营销模块数据模型（VIP/优惠券/活动） | 营销功能开发 |
| [04b-team-models.md](../04b-team-models.md) | 团队系统数据模型 | 团队功能开发 |
| [04c-enums-indexes.md](../04c-enums-indexes.md) | 枚举类型、数据库索引定义 | 数据库迁移、性能优化 |
| [04d-notification-models.md](../04d-notification-models.md) | 通知系统数据模型 | 通知功能开发 |
| [07-integration-test-plan.md](../07-integration-test-plan.md) | 集成测试计划 | 测试计划编写 |

---

## 快速查找指南

### 我要...

🤔 **理解业务规则**
→ [01-product.md](../01-product.md) 的"商业模式"和"抽成计算"章节

🏗️ **开始新功能开发**
1. 先看 [03-project-structure.md](../03-project-structure.md) 了解代码位置
2. 再看 [04-data-models.md](../04-data-models.md) 理解相关数据模型
3. 参阅 [05-testing-standard.md](../05-testing-standard.md) 编写测试

🔍 **查找某个数据模型的字段**
→ [04-data-models.md](../04-data-models.md) 或附录文档（04a/04b/04d）

🧪 **验证功能是否正确实现**
→ [05-testing-standard.md](../05-testing-standard.md) 的"测试检查清单"

🚀 **部署到生产环境**
→ [02-tech-stack.md](../02-tech-stack.md) 的"CI/CD 流程"和"部署脚本"章节

📊 **检查项目进度**
→ [06-project-management.md](../06-project-management.md) 或 [../docs/PROGRESS.md](../docs/PROGRESS.md)

---

## 文档维护规则

### 更新频率要求

| 文档 | 更新频率 | 负责人 |
|------|----------|--------|
| 01-product.md | 每次功能上线后 | PM/Backend Lead |
| 04-data-models.md | **每次模型变更后**（强制） | Backend Lead |
| 04c-enums-indexes.md | 每次枚举/索引变更 | Backend Lead |
| 06-project-management.md | 每周或模块完成时 | Project Lead |
| ../docs/PROGRESS.md | 每次版本发布 | Project Lead |

### 必读提醒

⚠️ **重要**: 所有代码修改前，必须先查阅相关 steering 文档，确保理解业务规则和技术规范。

💡 **提示**: 文档中 `> **注意**` 或 `> **重要**` 标记的内容为关键信息，开发前必读。

---

## 文档统计

| 类型 | 文档数 | 总行数（估算） |
|------|--------|----------------|
| 核心文档 | 4 | ~800 行 |
| 参考文档 | 3 | ~600 行 |
| 附录文档 | 5 | ~1200 行 |
| **总计** | **12** | **~2600 行** |

> 💡 **优化建议**: 为减少上下文消耗，AI 助手应优先加载核心文档，按需加载附录文档。
