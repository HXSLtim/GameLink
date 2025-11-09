# GameLink Cursor 规则生成报告

## 📊 生成概览

**生成日期**: 2025年11月5日  
**项目**: GameLink 游戏陪玩管理平台  
**规则总数**: 14个  
**总大小**: ~140KB

## ✅ 已生成的规则文件

### 1. 项目级规则 (2个)

| 文件名 | 大小 | 应用方式 | 说明 |
|--------|------|----------|------|
| `project-overview.mdc` | 2.4KB | alwaysApply | 项目概览、技术栈、核心业务模型 |
| `git-workflow.mdc` | 7.8KB | alwaysApply | Git工作流程、提交规范、分支管理 |

### 2. 后端规则 (7个)

| 文件名 | 大小 | 应用范围 | 说明 |
|--------|------|----------|------|
| `backend-architecture.mdc` | 4.7KB | `backend/**/*.go` | 四层架构、依赖注入、目录结构 |
| `backend-naming.mdc` | 6.6KB | `backend/**/*.go` | 命名规范、文件组织、JSON标签 |
| `backend-error-handling.mdc` | 7.6KB | `backend/**/*.go` | 错误定义、包装、Handler响应 |
| `backend-testing.mdc` | 9.4KB | `backend/**/*_test.go` | 测试组织、Mock使用、覆盖率 |
| `backend-api-design.mdc` | 11.9KB | `backend/internal/handler/**/*.go` | RESTful设计、请求响应、Swagger |
| `backend-context.mdc` | 9.5KB | `backend/**/*.go` | Context传递、超时控制、值传递 |
| `database-design.mdc` | 13.1KB | `backend/internal/model/**/*.go`<br>`backend/internal/repository/**/*.go` | GORM使用、关联关系、事务处理 |

### 3. 前端规则 (2个)

| 文件名 | 大小 | 应用范围 | 说明 |
|--------|------|----------|------|
| `frontend-api-patterns.mdc` | 10.7KB | `frontend/src/**/*.{ts,tsx}` | API客户端、服务层、错误处理 |
| `frontend-typescript-react.mdc` | 11.8KB | `frontend/src/**/*.{ts,tsx}` | 组件定义、Hooks、类型系统 |

### 4. 通用规则 (3个)

| 文件名 | 大小 | 应用方式 | 说明 |
|--------|------|----------|------|
| `security.mdc` | 12.0KB | alwaysApply | 认证授权、密码安全、输入验证 |
| `logging.mdc` | 9.9KB | description | 结构化日志、性能监控、日志配置 |
| `performance.mdc` | 11.7KB | description | 数据库优化、缓存策略、并发处理 |

### 5. 文档 (1个)

| 文件名 | 大小 | 说明 |
|--------|------|------|
| `README.md` | 5.4KB | 规则索引和使用指南 |

## 🎯 规则覆盖的主要领域

### 后端 (Go)

✅ **架构设计**
- 四层架构 (Model → Repository → Service → Handler)
- 依赖注入
- 接口设计

✅ **编码规范**
- 文件和包命名
- 结构体和接口命名
- 方法和变量命名
- JSON 标签规范 (camelCase)

✅ **错误处理**
- 错误定义和包装
- 使用 `fmt.Errorf` 和 `%w`
- 统一的错误响应

✅ **测试**
- AAA 模式
- 表驱动测试
- Mock 使用
- 测试覆盖率要求

✅ **API 设计**
- RESTful 路由
- 请求/响应结构
- HTTP 状态码
- Swagger 文档

✅ **数据库**
- GORM 模型定义
- Repository 模式
- 事务处理
- 性能优化

✅ **Context**
- Context 传递
- 超时控制
- 值传递

### 前端 (React + TypeScript)

✅ **组件开发**
- 函数组件
- Props 类型定义
- Hooks 使用
- 性能优化 (memo, useMemo, useCallback)

✅ **API 集成**
- Axios 客户端
- 服务层组织
- 错误处理
- 自定义 Hooks

✅ **类型系统**
- Interface 定义
- 泛型使用
- 类型继承和扩展

### 通用

✅ **安全性**
- JWT 认证
- 密码加密 (bcrypt)
- 输入验证
- CORS 配置
- 防止常见攻击 (SQL注入、XSS)

✅ **日志**
- 结构化日志 (slog)
- 日志级别
- 性能监控
- 错误追踪

✅ **性能**
- 数据库优化
- 缓存策略
- 并发处理
- 前端性能优化

✅ **Git 工作流**
- Conventional Commits
- 分支管理
- PR 规范
- 代码审查

## 📋 规则应用机制

### 自动应用规则 (alwaysApply: true)

以下规则会应用到**所有请求**:

1. `project-overview.mdc` - 项目概览
2. `git-workflow.mdc` - Git 工作流程
3. `security.mdc` - 安全性规范

### 基于文件类型应用 (globs)

规则会根据文件路径自动应用:

**后端 Go 文件** (`backend/**/*.go`):
- backend-architecture.mdc
- backend-naming.mdc
- backend-error-handling.mdc
- backend-context.mdc

**后端测试文件** (`backend/**/*_test.go`):
- backend-testing.mdc

**Handler 文件** (`backend/internal/handler/**/*.go`):
- backend-api-design.mdc

**Model/Repository 文件**:
- database-design.mdc

**前端文件** (`frontend/src/**/*.{ts,tsx}`):
- frontend-api-patterns.mdc
- frontend-typescript-react.mdc

### 手动调用规则 (description)

以下规则可以按需调用:

- `logging.mdc` - 日志记录规范
- `performance.mdc` - 性能优化指南

## 🎓 使用建议

### 1. 开发新功能时

在开始编写代码前,快速浏览相关规则:

**后端功能**:
```
1. 查看 backend-architecture.mdc 确认架构分层
2. 查看 backend-naming.mdc 确认命名规范
3. 查看 backend-api-design.mdc 了解 API 设计规范
```

**前端功能**:
```
1. 查看 frontend-typescript-react.mdc 了解组件规范
2. 查看 frontend-api-patterns.mdc 了解 API 调用规范
```

### 2. 代码审查时

使用规则作为审查清单:

- ✅ 是否遵循命名规范?
- ✅ 是否正确处理错误?
- ✅ 是否添加了测试?
- ✅ 是否考虑了安全性?
- ✅ 是否有性能问题?

### 3. 重构代码时

参考规范进行改进:

- 使用 `performance.mdc` 优化性能
- 使用 `logging.mdc` 改进日志
- 使用 `security.mdc` 加强安全

## 📚 与现有文档的关系

Cursor 规则是对现有文档的**补充和强化**:

| Cursor 规则 | 现有文档 | 关系 |
|-------------|----------|------|
| backend-*.mdc | backend/CODING_STANDARDS.md | 规则是文档的精简版,聚焦核心规范 |
| git-workflow.mdc | 团队 Git 规范 | 规则提供可执行的检查清单 |
| frontend-*.mdc | 前端已有规则 | 规则互补,覆盖更多场景 |

## 🔄 维护建议

### 定期更新

- 每月审查规则是否需要更新
- 新技术栈引入时添加相应规则
- 发现问题时及时修正

### 团队协作

- 规则变更需团队评审
- 重要规则修改需通知全员
- 鼓励团队成员提出改进建议

### 版本控制

- 规则文件纳入 Git 版本控制
- 重大变更记录在 CHANGELOG
- 使用语义化版本号

## ✨ 下一步

1. **熟悉规则**: 团队成员阅读规则文档
2. **实践应用**: 在实际开发中应用规则
3. **收集反馈**: 收集团队反馈,优化规则
4. **持续改进**: 根据项目发展更新规则

## 📞 支持

如有问题或建议:

1. 查看 `.cursor/rules/README.md`
2. 咨询项目负责人
3. 在团队会议上讨论

---

**生成工具**: Cursor AI  
**生成者**: Claude Sonnet 4.5  
**项目**: GameLink 游戏陪玩管理平台  
**版本**: 1.0.0  
**日期**: 2025年11月5日

