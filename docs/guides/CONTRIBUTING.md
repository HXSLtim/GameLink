# 贡献指南

感谢您对 GameLink 项目的关注！我们欢迎所有形式的贡献。本指南将帮助您了解如何参与项目开发。

## 🤝 贡献类型

我们欢迎以下类型的贡献：

- 🐛 **Bug报告**: 发现并报告问题
- 💡 **功能建议**: 提出新功能想法
- 📝 **文档改进**: 完善项目文档
- 🔧 **代码贡献**: 修复Bug或实现新功能
- 🧪 **测试用例**: 增加测试覆盖率
- 🎨 **UI/UX改进**: 改善用户界面和体验
- 📊 **性能优化**: 提升系统性能

## 🚀 开始贡献

### 1. Fork项目
```bash
# 在GitHub上Fork项目到您的账户
# 然后克隆到本地
git clone https://github.com/YOUR_USERNAME/gamelink.git
cd gamelink
```

### 2. 设置开发环境
```bash
# 添加上游仓库
git remote add upstream https://github.com/ORIGINAL_OWNER/gamelink.git

# 安装Go依赖
cd backend
go mod download

# 安装前端依赖
cd ../frontend/user-app
npm install
```

### 3. 创建开发分支
```bash
# 确保在最新的主分支
git checkout main
git pull upstream main

# 创建功能分支
git checkout -b feature/your-feature-name
# 或修复分支
git checkout -b fix/issue-number-description
```

## 📋 开发流程

### 代码风格

请遵循项目的代码规范：

#### Go代码规范
- 遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- 使用 `golangci-lint` 检查代码质量
- 函数和方法必须有注释
- 单元测试覆盖率 > 80%

#### 前端代码规范
- 使用 TypeScript 进行类型检查
- 遵循 ESLint 和 Prettier 配置
- 组件必须有 PropTypes 或 TypeScript 接口
- 使用语义化的变量和函数命名

### 提交规范

我们使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```bash
# 提交格式
<type>(<scope>): <subject>

<body>

<footer>
```

#### 提交类型
- `feat`: 新功能
- `fix`: Bug修复
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建工具、依赖更新等

#### 提交示例
```bash
feat(user): add user registration feature

- Add phone verification
- Add password encryption
- Add user profile creation

Closes #123
```

### 代码审查

所有代码提交都需要经过代码审查：

1. **创建Pull Request**
   - 使用清晰的标题和描述
   - 关联相关的Issue
   - 添加必要的标签

2. **自检清单**
   - [ ] 代码通过所有测试
   - [ ] 代码符合项目规范
   - [ ] 添加了必要的测试
   - [ ] 更新了相关文档

3. **响应反馈**
   - 及时回应审查意见
   - 根据建议修改代码
   - 保持友好和专业的沟通

## 🐛 报告Bug

### Bug报告模板
使用以下模板报告Bug：

```markdown
**Bug描述**
简要描述遇到的问题

**复现步骤**
1. 进入 '...'
2. 点击 '....'
3. 滚动到 '....'
4. 看到错误

**期望行为**
描述您期望发生的情况

**实际行为**
描述实际发生的情况

**截图**
如果适用，添加截图来帮助解释问题

**环境信息**
- 操作系统: [e.g. iOS 15.0, Windows 11]
- 浏览器: [e.g. Chrome 108.0]
- 应用版本: [e.g. v1.2.3]

**附加信息**
添加任何其他关于问题的信息
```

## 💡 功能请求

### 功能请求模板
```markdown
**功能描述**
简要描述您希望添加的功能

**问题背景**
描述这个功能要解决的问题

**解决方案**
描述您希望的解决方案

**替代方案**
描述您考虑过的其他解决方案

**附加信息**
添加任何其他关于功能请求的信息
```

## 📝 文档贡献

### 文档类型
- API文档
- 用户指南
- 开发文档
- 部署指南
- 架构文档

### 文档规范
- 使用清晰的标题结构
- 添加代码示例
- 包含图表和截图
- 保持文档更新

## 🧪 测试指南

### 运行测试
```bash
# 后端测试
cd backend
go test ./...

# 前端测试
cd frontend/user-app
npm test

# 端到端测试
npm run test:e2e
```

### 编写测试
- 为新功能编写单元测试
- 确保测试覆盖率
- 添加集成测试
- 编写性能测试

## 📦 发布流程

### 版本号规范
我们使用 [Semantic Versioning](https://semver.org/)：
- `MAJOR.MINOR.PATCH`
- `MAJOR`: 不兼容的API变更
- `MINOR`: 向后兼容的功能新增
- `PATCH`: 向后兼容的问题修正

### 发布检查清单
- [ ] 所有测试通过
- [ ] 文档已更新
- [ ] 变更日志已更新
- [ ] 版本号已更新
- [ ] 标签已创建

## 🏷 标签和里程碑

### 标签使用
- `bug`: Bug修复
- `enhancement`: 功能增强
- `documentation`: 文档相关
- `good first issue`: 适合新手的问题
- `help wanted`: 需要帮助

### 里程碑
- `v1.0.0`: 第一个正式版本
- `v1.1.0`: 功能增强版本
- `v2.0.0`: 重大版本更新

## 📞 联系方式

- **GitHub Issues**: [项目Issues页面](https://github.com/your-org/gamelink/issues)
- **讨论区**: [GitHub Discussions](https://github.com/your-org/gamelink/discussions)
 - **邮箱**: a2778978136@163.com
- **微信群**: 扫描二维码加入开发交流群

## 🏆 贡献者

我们感谢所有贡献者的努力！贡献者将被列入：

- README贡献者列表
- 发布公告
- 项目网站展示

## 📜 行为准则

### 我们的承诺
为了营造开放和友好的环境，我们承诺：

- 使用友好和包容的语言
- 尊重不同的观点和经验
- 优雅地接受建设性批评
- 关注对社区最有利的事情
- 对其他社区成员表示同理心

### 不可接受的行为
- 使用性别化语言或图像
- 人身攻击或政治攻击
- 公开或私下骚扰
- 未经明确许可发布他人的私人信息
- 其他在专业环境中可能被认为不当的行为

## 🎁 致谢

再次感谢您的贡献！您的参与让 GameLink 变得更好。

---

如果您有任何问题或建议，请随时联系我们。期待您的贡献！