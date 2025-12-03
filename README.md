# 🎮 GameLink - 现代化游戏陪玩管理平台

[![Go Version](https://img.shields.io/badge/Go-1.25.3+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![React Version](https://img.shields.io/badge/React-18+-61DAFB?style=flat&logo=react)](https://reactjs.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen)](https://github.com/HXSLtim/GameLink/actions)
[![Coverage](https://img.shields.io/badge/Coverage-76.4%25-yellow)](backend/LATEST_COVERAGE_REPORT.md)

**Go + React 全栈项目 | 智能订单分发 | 多角色管理 | 实时通讯**

---

## 🌟 项目简介

GameLink 是一个现代化的游戏陪玩管理平台，采用 Go 后端 + React 前端的架构，为游戏陪玩服务提供高效的订单分发、用户管理和陪玩师管理功能。

### 核心功能
- 🎯 **智能订单分发** - 自动匹配用户与陪玩师，支持抢单池和客服指派
- 👥 **多角色管理** - 用户/陪玩师/管理员权限体系
- 💬 **实时通讯** - WebSocket 即时通讯，支持群聊和私聊
- 💳 **完整支付** - 订单支付、退款、收益结算一体化
- 📊 **数据监控** - 实时订单状态、收益统计、系统指标
- 🔐 **安全认证** - JWT + RBAC 权限控制

---

## 🚀 快速开始

### 环境要求
- **Go**: 1.25.3+
- **Node.js**: 18+
- **MySQL**: 8.0+
- **Redis**: 6.0+

### 一键启动
```bash
# 克隆项目
git clone https://github.com/HXSLtim/GameLink.git
cd GameLink

# 使用快速启动脚本
./scripts/quick-start.sh
```

### 手动启动

#### 1. 后端服务
```bash
cd backend
go mod download
make run CMD=user-service
```

#### 2. 前端应用
```bash
cd frontend
npm install
npm run dev
```

#### 3. 访问应用
- 🌐 **前端应用**: http://localhost:5173
- 🔌 **后端API**: http://localhost:8080
- 📚 **API文档**: http://localhost:8080/swagger/index.html

---

## 📊 项目概览

### 系统架构
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   前端应用       │    │   后端API服务    │    │   数据存储       │
│                │    │                │    │                │
│ • React 18     │◄──►│ • Go 1.25.3    │◄──►│ • MySQL        │
│ • TypeScript   │    │ • Gin + GORM   │    │ • Redis        │
│ • WebSocket    │    │ • JWT Auth     │    │ • 文件存储      │
│ • 响应式设计     │    │ • Swagger API  │    │                │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### 技术栈
**后端**: Go 1.25.3, Gin, GORM, Redis, JWT, WebSocket, 纯Go SQLite驱动
**前端**: React 18, TypeScript, Vite, Less, WebSocket
**数据库**: MySQL, Redis, SQLite（测试环境）
**测试**: Go testing + Testify, Vitest, Playwright, 并发安全测试
**安全**: JWT + RBAC, CSRF防护, SQL注入防护, 输入验证
**错误处理**: 分层错误机制（Repository→Service→API标准化）

### 项目状态
- ✅ **后端完成度**: 85%
- ⏳ **前端完成度**: 70%
- 📈 **测试覆盖率**: 76.4%
- 📚 **文档完整性**: 95%
- 🛡️ **安全等级**: 高 (SQL注入漏洞已修复，JWT认证增强)
- ⚡ **性能优化**: 并发处理优化，纯Go测试环境

---

## 📁 项目结构

```
GameLink/
├── backend/                 # Go 后端服务
│   ├── cmd/                # 应用入口
│   ├── internal/           # 内部模块
│   ├── configs/            # 配置文件
│   └── docs/               # API 文档
├── frontend/               # React 前端应用
│   ├── src/                # 源代码
│   ├── public/             # 静态资源
│   └── docs/               # 前端文档
├── docs/                   # 项目文档
└── scripts/                # 部署脚本
```

---

## 🎯 功能特色

### 三端架构
- **用户端**: 首页、游戏列表、陪玩师浏览、订单创建、支付、评价
- **陪玩师端**: 工作台、订单管理、收益管理、服务管理、车队功能
- **管理后台**: 仪表盘、用户管理、订单监控、财务管理、系统设置

### 核心业务流程
1. **订单创建** - 用户选择服务，创建订单进入订单池
2. **智能分发** - 陪玩师抢单或客服指派
3. **服务执行** - 实时通讯，进度跟踪
4. **支付结算** - 自动分账，收益提现
5. **评价反馈** - 双向评价，信用积累

### 技术特色
- **测试体系**：表驱动测试、并发安全测试、边界值测试、纯Go测试环境
- **错误处理**：分层错误机制，Repository→Service→API标准化，支持错误码和详情
- **并发安全**：订单接单原子性操作，分布式锁机制，压力测试验证
- **安全加固**：SQL注入防护、JWT令牌管理、CSRF防护、输入验证

---

## 📚 文档导航

### 📋 开发文档
- **[开发指南](DEVELOPMENT.md)** - 详细开发环境搭建和规范
- **[部署指南](DEPLOYMENT.md)** - 生产环境部署和运维
- **[架构设计](ARCHITECTURE.md)** - 系统架构和设计理念
- **[API 文档](API.md)** - RESTful API 接口文档
- **[测试指南](backend/TESTING_GUIDE.md)** - 单元测试和集成测试最佳实践
- **[错误处理规范](backend/ERROR_HANDLING.md)** - 分层错误处理机制说明

### 🎯 功能指南
- **[前端开发完整指南](frontend/docs/FRONTEND_DEVELOPMENT_COMPLETE_GUIDE.md)**
- **[页面结构说明](frontend/docs/FRONTEND_PAGES_STRUCTURE.md)**
- **[用户端页面设计](frontend/docs/USER_FACING_PAGES_GUIDE.md)**

### 📊 项目报告
- **[项目状态报告](docs/PROJECT_STATUS_FINAL_REPORT.md)** - 整体项目进展和质量评估
- **[测试覆盖率报告](backend/LATEST_COVERAGE_REPORT.md)** - 单元测试和集成测试覆盖率分析
- **[用户接口设计报告](backend/USER_INTERFACE_INTEGRITY_REPORT.md)** - 三端接口完整性和一致性分析
- **[安全修复总结](SECURITY_FIXES_SUMMARY_2025-12-01.md)** - SQL注入、JWT认证等安全漏洞修复记录
- **[并发处理优化报告](backend/CONCURRENCY_FIX_ORDER_ACCEPT.md)** - 订单接单并发安全性分析和优化
- **[性能优化报告](backend/PERFORMANCE_FIX_PLAYER_LOOKUP.md)** - 数据库查询和缓存策略优化

---

## 🔧 开发工具

### 快捷命令
```bash
# 后端开发
cd backend
make lint          # 代码检查
make test           # 运行测试
make swagger        # 生成API文档

# 前端开发
cd frontend
npm run lint        # 代码检查
npm run test        # 运行测试
npm run build       # 构建生产版本
```

### 测试
```bash
# 运行所有测试
make test

# 生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

---

## 🛠️ 部署

### Docker 部署
```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps
```

### 生产部署
详细部署指南请参考 [DEPLOYMENT.md](DEPLOYMENT.md)

---

## 🤝 贡献指南

### 参与方式
1. Fork 项目到你的 GitHub 账户
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交代码 (`git commit -m 'feat: add AmazingFeature'`)
4. 推送分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 开发规范
- 遵循 [Go 编码规范](https://golang.org/doc/effective_go.html)
- 遵循 [TypeScript 编码规范](https://www.typescriptlang.org/docs/)
- 添加必要的测试用例
- 更新相关文档
- 通过所有 CI 检查

### 问题反馈
- 📋 **功能建议**: [Issues](https://github.com/HXSLtim/GameLink/issues)
- 🐛 **Bug报告**: [Issues](https://github.com/HXSLtim/GameLink/issues)
- 💬 **技术讨论**: [Discussions](https://github.com/HXSLtim/GameLink/discussions)

---

## 📞 联系我们

### 🏢 团队信息
- **项目负责人**: GameLink开发团队
- **技术支持**: a2778978136@63.com
- **商务合作**: business@gamelink.com

### 📱 更多资源
- **项目仓库**: https://github.com/HXSLtim/GameLink.git
- **技术博客**: https://blog.gamelink.com
- **在线演示**: https://demo.gamelink.com

---

## 📄 开源协议

本项目采用 [MIT License](LICENSE) 开源协议。

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给我们一个Star！**

**🚀 让我们一起构建更好的游戏陪玩生态！**

*最后更新: 2025-12-02 | 版本: v2.2*

</div>
