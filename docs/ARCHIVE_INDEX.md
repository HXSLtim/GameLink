# 📦 GameLink 文档归档索引

**更新日期**: 2026-01-03

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
├── backend/           # 后端相关归档（23个文档）
├── frontend/          # 前端相关归档（15个文档）
├── features/          # 功能设计文档（13个文档）
├── implementation/    # 实现报告（16个文档）
├── migration/         # 代码迁移文档（6个文档）
├── plans/             # 项目计划（7个文档）
├── summaries/         # 文档总结（7个文档）
├── README.md          # 归档说明
└── ARCHIVE_INDEX.md   # 详细归档索引
```

---

## 📊 归档统计

| 分类 | 数量 | 说明 |
|------|------|------|
| backend | 23 | 后端相关文档（RBAC、测试覆盖率等） |
| frontend | 15 | 前端相关文档（UI组件、导航等） |
| features | 13 | 功能设计文档（游戏、财务、用户端等） |
| implementation | 16 | 实现报告（CRUD、通知系统等） |
| migration | 6 | 代码迁移文档（camelCase迁移等） |
| plans | 7 | 项目计划（架构重构、改进路线图等） |
| summaries | 7 | 文档总结（覆盖率、项目状态等） |
| **总计** | **92** | **所有归档文档** |

---

## 📖 详细索引

完整的归档文档列表请查看：[docs/archive/ARCHIVE_INDEX.md](./archive/ARCHIVE_INDEX.md)

---

## 📋 维护规则

1. **核心文档**: 保留在 `/docs/` 根目录
2. **实现文档**: 完成后移至 `archive/implementation/`
3. **功能文档**: 移至 `archive/features/`
4. **后端文档**: 移至 `archive/backend/`
5. **前端文档**: 移至 `archive/frontend/`
6. **迁移文档**: 移至 `archive/migration/`
7. **计划文档**: 移至 `archive/plans/`
8. **总结文档**: 移至 `archive/summaries/`
9. **临时报告**: 定期清理或归档
10. **代码目录**: 不保留 `.md` 文档（除 README）

---

## 🔄 最近归档 (2026-01-03)

**重大重组**：
- 创建新的 `migration/` 子目录，归档 6 个 camelCase 迁移相关文档
- 将 5 个 RBAC 权限系统文档移至 `backend/`
- 将 9 个前端 UI/UX 文档移至 `frontend/`
- 将 5 个功能设计文档移至 `features/`
- 将 2 个实现/状态文档移至对应目录
- 更新归档索引，从 27 个根文件整理到 7 个分类子目录
- **归档文档总数**: 92 个
- **根目录清理**: 仅保留 README.md 和 ARCHIVE_INDEX.md
