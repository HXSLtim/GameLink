# 📦 GameLink 文档归档索引

**更新日期**: 2025-12-16

---

## 📚 文档结构

| 目录 | 说明 |
|------|------|
| `/docs/` | 核心文档（PRD、开发、进度） |
| `/docs/api/` | API 规范文档 |
| `/docs/guides/` | 开发指南 |
| `/docs/standards/` | 编码规范 |
| `/docs/archive/` | 历史归档文档 |

---

## 🎯 核心文档

| 文档 | 说明 |
|------|------|
| `PRD.md` | 产品需求文档 |
| `DEVELOPMENT.md` | 开发技术文档 |
| `PROGRESS.md` | 开发进度追踪 |
| `README.md` | 项目说明 |
| `INDEX.md` | 文档导航 |
| `CI-CD.md` | CI/CD 配置 |
| `COMPONENT_LIBRARY.md` | 组件库文档 |
| `UI_DESIGN_SPEC.md` | UI 设计规范 |
| `INTERNATIONALIZATION.md` | 国际化方案 |
| `即时通讯系统设计文档.md` | IM 系统设计 |

---

## 📦 归档目录结构

```
docs/archive/
├── implementation/     # 功能实现文档
├── features/          # 功能设计文档
├── plans/             # 规划文档
├── reports/           # 报告文档
├── summaries/         # 总结文档
├── frontend/          # 前端相关归档
├── backend/           # 后端相关归档
├── coverage/          # 测试覆盖率报告
└── temp-reports/      # 临时报告
```

### implementation/ - 实现文档
- `NOTIFICATION_SYSTEM_FIX.md` - 通知系统修复
- `USER_MANAGEMENT_IMPLEMENTATION.md` - 用户管理实现
- `USER_MANAGEMENT_IMPLEMENTATION_COMPLETE.md` - 用户管理完成报告
- `FINANCIAL_MANAGEMENT_IMPLEMENTATION.md` - 财务管理实现
- `FINANCIAL_MANAGEMENT_SUMMARY.md` - 财务管理总结
- `USER_SIDE_IMPLEMENTATION.md` - 用户端实现
- `WORKFLOW_C_IMPLEMENTATION_GUIDE.md` - 工作流实现指南
- `WORKFLOW_C_IMPLEMENTATION_SUMMARY.md` - 工作流实现总结
- `WORKFLOW_C_QUICK_REFERENCE.md` - 工作流快速参考

### features/ - 功能设计
- `FINANCIAL_MANAGEMENT_DESIGN.md` - 财务管理设计
- `FINANCIAL_MANAGEMENT_REQUIREMENTS.md` - 财务管理需求
- `USER_SIDE_PLANNING.md` - 用户端规划
- `USER_SIDE_QUICKSTART.md` - 用户端快速开始
- `USER_FLOW_PROTOTYPE.md` - 用户流程原型

### plans/ - 规划文档
- `BACKEND_ARCHITECTURE_REFACTOR_PLAN.md` - 后端架构重构计划
- `WORK_DISTRIBUTION_PLAN.md` - 工作分配计划

### reports/ - 报告文档
- `DOCUMENTATION_SUMMARY.md` - 文档总结
- `DOCUMENT_CLEANUP_REPORT.md` - 文档清理报告

### frontend/ - 前端归档
- `BUILD_OPTIMIZATION.md` - 构建优化
- `OPTIMIZATION_SETUP.md` - 优化配置

---

## 📋 维护规则

1. **核心文档**: 保留在 `/docs/` 根目录
2. **实现文档**: 完成后移至 `archive/implementation/`
3. **临时报告**: 定期清理或归档
4. **代码目录**: 不保留 `.md` 文档（除 README）

---

## 🔄 最近归档 (2025-12-16)

- 清理 backend 测试产物（.exe, .txt）
- 归档 backend 实现文档 (3个)
- 归档 frontend 优化文档 (2个)
- 归档 docs 实现/规划文档 (9个)
