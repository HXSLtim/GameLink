# 📋 GameLink 后端项目规范指南

**项目版本**: v1.0
**最后更新**: 2025年11月2日
**维护团队**: GameLink 开发团队

---

## 🎯 快速导航

### 📚 核心文档
| 文档 | 描述 | 重要性 |
|------|------|--------|
| [CODING_STANDARDS.md](./CODING_STANDARDS.md) | 代码编写规范 | ⭐⭐⭐ 必须 |
| [docs/backend/README.md](./docs/backend/) | 后端完整文档 | ⭐⭐⭐ 必须 |
| [docs/backend/QUICK_START_UNIFIED.md](./docs/backend/) | 快速开始指南 | ⭐⭐ 推荐 |

### 📊 重构报告（已归档）
| 报告 | 描述 | 状态 |
|------|------|------|
| [archive/docs/refactoring-reports/REPOSITORY_REFACTORING_COMPLETE.md](./archive/docs/refactoring-reports/) | Repository层重构完成报告 | ✅ 完成 |
| [archive/docs/refactoring-reports/COMPLETE_REFACTORING_STATUS_REPORT.md](./archive/docs/refactoring-reports/) | 完整重构状态报告 | ✅ 完成 |
| [archive/docs/refactoring-reports/DIRECTORY_STRUCTURE_COMPLIANCE_REPORT.md](./archive/docs/refactoring-reports/) | 目录规范符合性评估 | ✅ 完成 |

---

## 🏗️ 项目架构概览

### 四层架构
```
🌐 HTTP请求
   ↓
🎯 Handler层 (API处理)
   ↓
💼 Service层 (业务逻辑)
   ↓
🗄️ Repository层 (数据访问)
   ↓
📊 Model层 (数据模型)
```

### 目录结构
```
backend/
├── 📁 cmd/              # 应用程序入口
├── 📁 internal/         # 内部包
│   ├── 📁 model/        # 数据模型
│   ├── 📁 repository/   # 数据访问层
│   ├── 📁 service/      # 业务逻辑层
│   ├── 📁 handler/      # API处理层
│   └── 📁 middleware/   # 中间件
├── 📁 docs/             # 项目文档
├── 📁 configs/          # 配置文件
├── 📁 scripts/          # 脚本文件
├── 📄 CODING_STANDARDS.md # 代码规范 ⭐
└── 📄 PROJECT_GUIDELINES.md # 项目指南 ⭐
```

---

## 📝 开发规范要点

### 1. 文件命名规范
```go
✅ 推荐：简洁、明确
├── repository/user/repository.go
├── service/auth/auth.go
├── handler/admin/user.go
└── user_test.go

❌ 避免：冗余后缀
├── repository/user/user_gorm_repository.go
├── service/auth/auth_service.go
└── handler/admin/user_handler.go
```

### 2. 代码分层原则
- **Model**: 纯数据模型，无业务逻辑
- **Repository**: 纯数据访问，无业务逻辑
- **Service**: 业务逻辑，无HTTP处理
- **Handler**: HTTP处理，无业务逻辑

### 3. 测试规范
- 测试文件与源文件同目录
- 测试包名与源码包名一致
- 测试函数命名：`TestFunction_Scenario`

### 4. Git提交规范
```bash
✅ 推荐：<类型>(<范围>): <描述>
feat(service): 添加用户创建功能
fix(handler): 修复订单状态更新错误
docs(readme): 更新安装说明
```

---

## 🚀 快速开始

### 环境要求
- Go 1.21+
- MySQL 8.0+
- Redis 6.0+

### 启动步骤
```bash
# 1. 安装依赖
go mod tidy

# 2. 配置环境
cp configs/config.example.yaml configs/config.yaml

# 3. 初始化数据库
go run cmd/migrate/main.go up

# 4. 启动服务
go run cmd/main.go

# 5. 验证服务
curl http://localhost:8080/health
```

### 开发工具
```bash
# 代码检查
golangci-lint run

# 热重载
air

# 测试
go test ./...
```

---

## 📊 当前项目状态

### ✅ 已完成的重构
- **Repository层重构** (100%) - 命名统一，结构清晰
- **Handler层整合** (90%) - 按用户角色组织
- **代码规范制定** (100%) - 完整的编码标准

### ⚠️ 需要注意的问题
- **Service层命名** - 少量文件仍有冗余后缀
- **空目录清理** - 重构遗留的空目录需清理
- **测试覆盖率** - 部分模块测试覆盖不足

### 🎯 下一步计划
1. 完成Service层最终优化
2. 清理空目录和遗留文件
3. 提升测试覆盖率
4. 完善API文档

---

## 🔧 开发工具配置

### 必需工具
```bash
# 安装开发工具
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/cosmtrek/air@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

### IDE配置
推荐使用 VS Code + Go 插件：
```json
{
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.testOnSave": true,
    "go.coverOnSave": true
}
```

---

## 📞 获取帮助

### 📚 文档资源
- **完整文档**: [docs/backend/](./docs/backend/)
- **API文档**: 启动服务后访问 `/swagger/index.html`
- **重构报告**: [archive/docs/refactoring-reports/](./archive/docs/refactoring-reports/)

### 🐛 问题反馈
- 技术问题：查看 [docs/backend/TROUBLESHOOTING.md](./docs/backend/)
- 功能建议：创建 GitHub Issue
- 紧急问题：联系团队负责人

### 🤝 团队协作
- 代码审查：严格按照清单检查
- 规范遵守：必须遵循 [CODING_STANDARDS.md](./CODING_STANDARDS.md)
- 文档更新：重要变更需同步更新文档

---

## 📋 开发检查清单

### 编码前
- [ ] 阅读相关需求文档
- [ ] 了解现有代码结构
- [ ] 创建功能分支

### 编码中
- [ ] 遵循代码规范
- [ ] 编写单元测试
- [ ] 考虑错误处理
- [ ] 添加必要注释

### 编码后
- [ ] 运行测试：`go test ./...`
- [ ] 代码检查：`golangci-lint run`
- [ ] 格式化：`go fmt ./...`
- [ ] 构建验证：`go build ./...`

### 提交前
- [ ] 更新相关文档
- [ ] 提交信息规范
- [ ] 推送分支
- [ ] 创建Pull Request

---

## ✨ 总结

GameLink 后端项目采用标准的四层架构，代码组织清晰，文档完善。所有开发人员应该：

1. **严格遵循代码规范** - [CODING_STANDARDS.md](./CODING_STANDARDS.md)
2. **保持架构清晰** - 遵循分层原则
3. **编写充分测试** - 保证代码质量
4. **及时更新文档** - 保持文档同步

**让我们一起构建高质量、可维护的代码库！** 🚀