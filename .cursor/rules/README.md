# GameLink Cursor 规则索引

本目录包含 GameLink 项目的 Cursor AI 编码规则。这些规则帮助 Cursor 更好地理解项目结构和编码规范。

## 📋 规则列表

### 项目级规则 (Always Apply)

1. **project-overview.mdc** - 项目概览
   - 技术栈介绍
   - 项目结构
   - 核心业务模型
   - 开发原则

2. **git-workflow.mdc** - Git 工作流程和提交规范
   - 分支管理
   - 提交消息格式
   - PR 规范
   - 代码审查

### 后端规则 (Backend Go)

3. **backend-architecture.mdc** - 后端架构和分层规范
   - 四层架构设计
   - 各层职责定义
   - 依赖注入原则
   - 目录结构规范

4. **backend-naming.mdc** - 后端命名规范
   - 文件命名
   - 包命名
   - 接口和结构体命名
   - 方法和变量命名
   - JSON 标签规范

5. **backend-error-handling.mdc** - 错误处理规范
   - 错误定义和包装
   - 错误检查模式
   - Handler 层错误响应
   - 日志记录

6. **backend-testing.mdc** - 后端测试规范
   - 测试文件组织
   - 测试命名
   - AAA 模式和表驱动测试
   - Mock 使用
   - 测试覆盖率

7. **backend-api-design.mdc** - API 设计规范
   - RESTful 路由设计
   - 请求/响应格式
   - Handler 实现模式
   - Swagger 文档

8. **backend-context.mdc** - Context 使用规范
   - Context 传递原则
   - GORM 中使用 Context
   - Context 超时控制
   - Context 值传递

9. **database-design.mdc** - 数据库设计和 GORM 使用规范
   - 数据模型定义
   - 关联关系
   - Repository 实现
   - 事务处理
   - 性能优化

### 前端规则 (Frontend React + TypeScript)

10. **frontend-api-patterns.mdc** - 前端 API 和服务层规范
    - API 客户端结构
    - API 服务组织
    - 类型定义
    - 错误处理
    - 自定义 Hooks

11. **frontend-typescript-react.mdc** - TypeScript 和 React 组件规范
    - 组件定义
    - Hooks 使用
    - 类型定义
    - 事件处理
    - 性能优化

### 通用规则

12. **security.mdc** - 安全性规范和最佳实践
    - 认证和授权
    - 密码安全
    - 输入验证
    - CORS 配置
    - 安全响应头

13. **logging.mdc** - 日志记录规范
    - 后端结构化日志 (slog)
    - 前端日志封装
    - 日志级别
    - 性能日志
    - 日志配置

14. **performance.mdc** - 性能优化指南
    - 数据库优化
    - 缓存策略
    - 并发优化
    - React 性能优化
    - 资源优化

## 🎯 规则应用场景

### 编写后端代码时

相关规则会自动应用:
- backend-architecture.mdc (所有 Go 文件)
- backend-naming.mdc (所有 Go 文件)
- backend-error-handling.mdc (所有 Go 文件)
- backend-context.mdc (所有 Go 文件)
- backend-testing.mdc (测试文件)
- backend-api-design.mdc (Handler 文件)
- database-design.mdc (Model 和 Repository 文件)

### 编写前端代码时

相关规则会自动应用:
- frontend-api-patterns.mdc (所有 TS/TSX 文件)
- frontend-typescript-react.mdc (所有 TS/TSX 文件)

### 所有开发场景

这些规则始终应用:
- project-overview.mdc
- git-workflow.mdc
- security.mdc

## 📖 如何使用

### 查看规则

在 Cursor 中,规则会自动应用到相应的文件。你也可以:

1. 在聊天中询问: "@规则名称"
2. 查看 `.cursor/rules/` 目录下的 `.mdc` 文件
3. 使用 Cursor 的规则面板查看

### 手动应用规则

某些规则设置为手动应用,可以在需要时通过以下方式调用:

```
请根据 logging 规则检查我的日志代码
请根据 performance 规则优化这段代码
```

### 修改规则

如果需要修改规则:

1. 编辑对应的 `.mdc` 文件
2. 保存后规则会自动更新
3. 通知团队成员同步更新

## 🔍 规则元数据说明

### alwaysApply: true

规则会应用到所有请求,适用于:
- 项目概览
- Git 工作流程
- 安全性规范

### globs: pattern

规则只应用到匹配的文件,例如:
- `backend/**/*.go` - 所有后端 Go 文件
- `backend/**/*_test.go` - 所有后端测试文件
- `frontend/src/**/*.{ts,tsx}` - 所有前端 TS/TSX 文件

### description

提供规则的简短描述,帮助 AI 理解何时应用该规则

## 📚 参考文档

项目中的其他重要文档:

- [backend/CODING_STANDARDS.md](../../backend/CODING_STANDARDS.md) - 后端详细编码规范
- [backend/PROJECT_GUIDELINES.md](../../backend/PROJECT_GUIDELINES.md) - 后端项目指南
- [CLAUDE.md](../../CLAUDE.md) - Claude AI 项目指南

## 🤝 贡献

如果发现规则需要改进或添加新规则:

1. 创建新的 `.mdc` 文件
2. 添加适当的元数据 (alwaysApply/globs/description)
3. 更新本 README 文件
4. 提交 PR 供团队审查

## ✨ 最佳实践

1. **保持规则简洁** - 规则应该清晰、简洁、易于理解
2. **提供示例** - 使用 ✅ 和 ❌ 标记好的和不好的示例
3. **及时更新** - 当项目规范变化时,及时更新规则
4. **团队共识** - 规则应该反映团队的共识,而非个人偏好
5. **保持一致** - 确保规则之间不冲突,保持一致性

---

**最后更新**: 2025年11月5日
**维护者**: GameLink 开发团队

