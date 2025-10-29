# 📚 GameLink 项目文档

欢迎来到 GameLink 项目文档中心！这里包含了项目的所有技术文档、开发指南和规范标准。

## 🗂️ 文档分类

### 📋 项目概述
- **[主 README](../README.md)** - 项目介绍和快速开始指南
- **[CLAUDE.md](../CLAUDE.md)** - Claude Code 开发配置和工作规范
- **[CONTRIBUTING.md](../CONTRIBUTING.md)** - 贡献指南和开发流程
- **[AGENTS.md](../AGENTS.md)** - AI 开发代理使用指南

### 🔄 迁移和更新
- **[CamelCase 迁移报告](CAMELCASE_MIGRATION_REPORT.md)** - snake_case 到 camelCase 迁移完整记录
- **[优化指南](../optimization_guide.md)** - 性能优化最佳实践

### 🛠️ 开发规范
- **[Go 编码规范](go-coding-standards.md)** - Go 语言编码标准和最佳实践
- **[API 设计规范](api-design-standards.md)** - RESTful API 设计原则和规范

### 📁 前端文档
详见 `frontend/docs/` 目录：
- **[前端技术文档](../frontend/docs/TECHNICAL_DOCUMENTATION.md)** - 前端架构和技术栈详解
- **[开发者指南](../frontend/docs/DEVELOPER_GUIDE.md)** - 前端开发详细指南
- **[用户文档](../frontend/docs/USER_DOCUMENTATION.md)** - 用户操作手册

### 🗄️ 后端文档
详见 `backend/docs/` 目录：
- **[数据库种子文档](../backend/docs/database-seed.md)** - 数据库初始化和种子数据
- **[加密中间件文档](../backend/docs/crypto-middleware.md)** - 请求加密和解密机制

## 🎯 快速导航

### 🚀 新手入门
1. 阅读 [主 README](../README.md) 了解项目概况
2. 搭建开发环境（后端 + 前端）
3. 查看 [CamelCase 迁移报告](CAMELCASE_MIGRATION_REPORT.md) 了解最新规范
4. 开始开发：参考 [Go 编码规范](go-coding-standards.md) 和 [前端开发者指南](../frontend/docs/DEVELOPER_GUIDE.md)

### 🔧 日常开发
- **代码规范**: [Go 编码规范](go-coding-standards.md) | [前端编码规范](../frontend/docs/design/CODING_STANDARDS.md)
- **API 开发**: [API 设计规范](api-design-standards.md)
- **测试指南**: 主 README 中的测试章节
- **部署指南**: 主 README 中的部署章节

### 📊 项目维护
- **性能优化**: [优化指南](../optimization_guide.md)
- **架构决策**: 各技术文档中的设计说明
- **更新日志**: Git commit history 和 releases

## 🏗️ 项目架构概览

```
GameLink Documentation Structure
├── 📋 项目文档 (根目录)
│   ├── README.md              # 项目主文档
│   ├── CLAUDE.md              # AI 开发配置
│   ├── CONTRIBUTING.md        # 贡献指南
│   └── optimization_guide.md  # 性能优化指南
├── 🔄 迁移报告 (docs/)
│   └── CAMELCASE_MIGRATION_REPORT.md  # 命名规范迁移
├── 🛠️ 开发规范 (docs/)
│   ├── go-coding-standards.md        # Go 编码规范
│   └── api-design-standards.md       # API 设计规范
├── 📱 前端文档 (frontend/docs/)
│   ├── TECHNICAL_DOCUMENTATION.md    # 技术文档
│   ├── DEVELOPER_GUIDE.md           # 开发指南
│   ├── USER_DOCUMENTATION.md        # 用户文档
│   ├── design/                      # 设计系统
│   ├── features/                    # 功能实现指南
│   ├── api/                        # API 集成文档
│   └── archive/                     # 归档报告
└── 🖥️ 后端文档 (backend/docs/)
    ├── database-seed.md            # 数据库种子
    ├── crypto-middleware.md        # 加密中间件
    └── go-coding-standards.md      # Go 编码规范
```

## 📝 文档维护

### 📖 文档规范
- 所有文档使用 Markdown 格式
- 标题使用标准的层级结构（# ## ### ####）
- 代码块标明语言类型
- 使用表格和列表提高可读性
- 添加适当的 emoji 增强视觉效果

### 🔄 更新频率
- **项目文档**: 每次重大更新后
- **API 文档**: 每次接口变更后
- **迁移报告**: 每次重大迁移后
- **规范文档**: 每季度审查更新

### 🤝 贡献指南
欢迎提交文档改进建议：
1. Fork 项目仓库
2. 创建文档改进分支
3. 提交 Pull Request
4. 等待代码审查和合并

## 📞 联系方式

如有文档相关的问题或建议，请联系：

- **项目维护**: GameLink 开发团队
- **技术支持**: dev@gamelink.com
- **文档反馈**: docs@gamelink.com
- **问题报告**: [GitHub Issues](https://github.com/your-org/gamelink/issues)

---

**最后更新**: 2025年1月29日
**文档版本**: 1.0.0
**维护团队**: GameLink 开发团队

<div align="center">

**📚 知识共享，共同进步！**

Made with ❤️ by GameLink Team

</div>