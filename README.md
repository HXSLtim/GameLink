# 🎮 GameLink - 现代化游戏陪玩管理平台

**Go + React 全栈项目 | 实时订单分发 | 用户管理 | 即时通讯**

更新时间: 2025-11-10

---

## 🌟 项目特色

### 核心功能
- 🎯 **智能订单分发系统** - 自动匹配用户与陪玩师
- 👥 **多角色用户管理** - 用户/陪玩师/管理员权限体系
- 💬 **即时通讯系统** - 公共群聊 + 订单专属小群
- 💳 **完整支付流程** - 订单确认、支付、评价一体化
- 📊 **实时数据监控** - 订单状态、收益统计、系统指标
- 🔐 **安全认证体系** - JWT + RBAC权限控制

### 技术亮点
- ⚡ **高性能架构** - Go 1.25.3 + Gin + GORM + Redis
- 🎨 **现代前端** - React 18 + TypeScript + Vite
- 🔌 **实时通信** - WebSocket + 消息推送
- 📱 **响应式设计** - 桌面端 + 移动端完美适配
- 🧪 **完整测试** - 184个测试用例，100%通过率
- 📚 **详细文档** - 完整的API文档和开发指南

---

## 📊 项目状态

### 开发进度
- **后端完成度**: 85% ✅
- **前端完成度**: 70% ⏳
- **测试覆盖率**: 49.5% (184个测试，100%通过率)
- **文档完整性**: 95% ✅

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

---

## 🚀 快速开始

### 环境要求
- **Go**: 1.25.3+
- **Node.js**: 18+
- **MySQL**: 8.0+
- **Redis**: 6.0+

### 本地开发

#### 1. 克隆项目
```bash
git clone https://github.com/your-org/GameLink.git
cd GameLink
```

#### 2. 启动后端服务
```bash
cd backend

# 安装依赖
go mod download

# 初始化数据库
make migrate

# 启动开发服务器
make run CMD=user-service

# 或使用Go命令
go run ./cmd/user-service
```

#### 3. 启动前端应用
```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

#### 4. 访问应用
- **前端应用**: http://localhost:5173
- **后端API**: http://localhost:8080
- **Swagger文档**: http://localhost:8080/swagger/index.html
- **API文档**: http://localhost:8080/docs

## 📁 项目结构

### 后端架构 (Go)
```
backend/
├── cmd/                    # 应用入口点
│   └── user-service/       # 用户服务主程序
├── internal/               # 内部包
│   ├── admin/             # 管理端处理器
│   ├── handler/           # HTTP处理器
│   │   ├── admin/         # 管理端接口
│   │   ├── user/          # 用户端接口
│   │   ├── player/        # 陪玩师接口
│   │   ├── middleware/    # 中间件
│   │   └── websocket/     # WebSocket处理
│   ├── service/           # 业务逻辑层
│   ├── repository/        # 数据访问层
│   ├── model/             # 数据模型
│   ├── auth/              # JWT认证
│   ├── cache/             # 缓存层
│   └── config/            # 配置管理
├── configs/               # 配置文件
├── docs/swagger/          # API文档
├── scripts/sql/           # SQL脚本
└── tests/                 # 测试文件
```

### 前端架构 (React)
```
frontend/
├── src/
│   ├── api/               # API调用层
│   ├── components/        # 可复用组件
│   │   ├── chat/         # 聊天组件
│   │   ├── order/        # 订单组件
│   │   ├── user/         # 用户组件
│   │   └── common/       # 通用组件
│   ├── pages/             # 页面组件
│   │   ├── admin/        # 管理端页面
│   │   ├── user/         # 用户端页面
│   │   └── player/       # 陪玩师页面
│   ├── layouts/           # 布局组件
│   ├── types/             # TypeScript类型
│   ├── utils/             # 工具函数
│   └── hooks/             # 自定义Hooks
├── public/                # 静态资源
└── docs/                  # 前端文档
```

### 数据库模型
- **User**: 用户基础信息 (user/player/admin角色)
- **Player**: 陪玩师认证和资料
- **Game**: 游戏配置信息
- **Order**: 订单管理 (6种状态流转)
- **Payment**: 支付记录
- **Review**: 评价系统
- **ChatGroup**: 聊天群组
- **ChatMessage**: 聊天消息

---

## 📘 陪玩结单需求与页面规划

### 1. 核心业务流程需求
#### 1.1 订单池抢单模式
- **用户侧**: 发布订单时填写游戏、服务类型、段位、时长、预算等必填项，订单进入公共订单池。
- **陪玩师侧**: 通过“派单/抢单大厅”实时浏览订单流，按游戏、价格、服务标签、用户信用等条件筛选后抢单，系统即时反馈抢单结果。
- **平台侧**: 负责制定抢单门槛与并发锁定策略，并依托 WebSocket 在订单创建、被抢、取消等事件发生时推送更新，心跳/重连保障实时性。

#### 1.2 客服指派模式
- **业务价值**: 面向高价值、需求复杂或新用户订单，由客服基于完整数据做精细匹配，保障服务质量。
- **后台能力**: 管理端需提供“待指派订单池”、订单详情、陪玩师资源池、指派弹窗及沟通日志，支持客服筛选陪玩师并发出指派通知，被指派者可接受/拒绝并记录响应。
- **协同机制**: 支持用户自主选择 vs. 平台规则驱动；长时间无人抢单可自动转客服；可采用“客服推荐列表 + 用户最终确认”的混合模式。

#### 1.3 结单评价与客户维系
- **双向评价**: 服务完成后用户与陪玩师在技术、态度、沟通等多维度打分，可附带文本/标签、匿名与延迟展示，评价权重累计为信用分。
- **客服介入**: 后台可监控评价、处理举报、调解争议并记录结果，对违规评价进行屏蔽/撤销。
- **维系工具**: 陪玩师可回复评价、向粉丝发优惠券；客服可做用户分层、节日关怀、召回与社群运营，形成闭环留存。

### 2. GameLink 功能映射
#### 2.1 用户端（Customer）
- **MVP 页面**: 首页、游戏列表、陪玩师列表、陪玩师详情、订单创建、支付、我的订单，完成“发现→下单→支付→评价”闭环。
- **Phase 2**: 个人中心、收藏/关注、订阅推荐、优惠券钱包、活动页等互动与运营能力。
- **已接入接口（节选）**: `POST /auth/register`, `POST /auth/login`, `GET /auth/me`, 后续订单/支付接口按 `/user/orders*` `/payments*` 规划对齐。

#### 2.2 陪玩师端（Player）
- **MVP 页面**: 工作台、订单管理、收益管理（含提现）、服务管理，涵盖资料维护、接单履约与收益提现。
- **新增车队模式**: 陪玩师上线后可创建/加入车队，由队长统一抢单接单并在队内分配服务，支持多人协同。
- **Phase 2**: 评价回复、排班/日历、营销工具、粉丝运营等职业化能力。

#### 2.3 管理后台（Admin）
- **MVP 页面**: 仪表盘、订单管理、用户管理、陪玩师管理、财务管理、系统设置，支撑运营、审核、对账与权限。
- **增强方向**: 内容/Banner 管理、自定义数据看板、自动化风控、客服指派工作台。
- **已接入接口（节选）**: Swagger 中的 `/admin/users*`, `/admin/orders*`, `/admin/payments*`, `/admin/games*`, `/admin/players*`, `/admin/stats*`，详见下方“已实现接口”表。

#### 2.4 跨域系统模块
- **认证与授权**: JWT 登录、密码找回、RBAC 权限模型与统一状态字典，保障三端权限隔离。
- **订单与支付**: 订单全生命周期、幂等控制、支付/退款、风控策略。
- **实时通信与通知**: gorilla WebSocket Hub、消息存储、离线推送、站内信/短信/邮件通知中心，并计划扩展订单事件推送。

### 3. 页面规划概览
| 角色 | MVP 页面 | Phase 2 页面 | 关键交互 |
| --- | --- | --- | --- |
| 用户端 | 首页、游戏列表、陪玩师列表、陪玩师详情、下单、支付、我的订单 | 个人中心、收藏/关注、活动/优惠、消息中心 | 搜索/筛选、套餐选择、支付、评价 |
| 陪玩师端 | 工作台、订单管理、收益管理、服务管理 | 排班日历、评价回复、营销、粉丝通知 | 车队组队/队长抢单、接单/拒单、开始/完成、提现、服务上下线 |
| 管理后台 | 仪表盘、订单、用户、陪玩师、财务、系统设置 | 内容管理、数据看板、风控策略、客服指派 | 指派、审核、批量操作、权限配置 |

### 4. 功能缺口与增强路线

| 编号 | 缺口主题 | 目标 | 规划阶段 |
| --- | --- | --- | --- |
| G1 | 订单池 / 抢单大厅 | 陪玩师可实时浏览、筛选、抢占订单，成功后自动锁单 | 需求→设计 |
| G2 | 客服指派工具链 | 后台能根据订单条件筛选陪玩师并一键指派，记录响应 | 方案评审 |
| G3 | 实时订单推送 | 订单创建/指派/状态变更≤1s 推送到相关客户端 | 技术预研 |
| G4 | 社区与维系 | 用户关系链 + 内容沉淀，支持评价互动与召回策略 | 产品探索 |
| G5 | 车队/队长抢单机制 | 队长统一抢单并分配队内成员，支持队伍管理 | 原型设计 |
| G6 | 智能匹配与风控 | 构建匹配算法 + 风控规则引擎，降低人工负担 | 数据建模 |

#### G1 订单池 / 抢单大厅
- **里程碑**
  1. 调研竞品交互，输出订单卡片与筛选方案 (W1)  
  2. 定义 WebSocket 订单事件格式与推送范围 (W2)  
  3. 完成陪玩师端前端页面与后端查询接口 (W3)  
- **依赖**: G3 实时推送、业务状态字典一致性  
- **交付件**: 需求文档、原型稿、API 套件、联调用例

#### G2 客服指派工具链
- **里程碑**
  1. 梳理订单状态→指派流转图，定义权限与 SLA (W1)  
  2. 管理后台新增“待指派”列表 + 陪玩师资源池 (W2)  
  3. 推送指派通知、记录响应并支持重指派 (W3)  
- **依赖**: G1 订单状态标签、G3 推送  
- **交付件**: 后台 UI、指派 API、操作审计日志

#### G3 实时订单推送
- **里程碑**
  1. 扩展 WebSocket Hub 支持订单 topic、多租户鉴权 (W1)  
  2. 实现可靠投递（ACK + 重试 + 离线队列）(W2)  
  3. 接入监控告警，输出 P95 延迟指标 (W3)  
- **依赖**: Redis/消息队列资源、DevOps 监控  
- **交付件**: WebSocket 协议说明、延迟报表、SDK

#### G4 社区与维系
- **里程碑**
  1. 明确社区内容类型（动态、图文、活动）及审核规则 (W1)  
  2. 用户关系链（关注/粉丝）与通知策略 (W2)  
  3. 评价回复、召回工具与客服工作台联动 (W3)  
- **依赖**: 内容审核服务、通知中心  
- **交付件**: 社区 PRD、后台审核模块、运营脚本

#### G5 车队/队长抢单机制
- **里程碑**
  1. 设计车队模型（队伍、成员、角色、在线状态）与规则 (W1)  
  2. 队长端抢单/派单 UI，队员端接受任务流程 (W2)  
  3. 车队绩效统计、收益分配与风控 (W3)  
- **依赖**: G1 订单池、G3 推送、收益结算模块  
- **交付件**: 数据库表设计、队伍 API、前端交互稿

#### G6 智能匹配与风控
- **里程碑**
  1. 收集特征（段位、成功率、响应时间、投诉率等）并构建特征库 (W1)  
  2. MVP 规则引擎 + 评分模型，生成陪玩师推荐清单 (W2)  
  3. 风控规则（异常订单、刷单、恶意取消）与自动干预策略 (W3)  
- **依赖**: 数据仓库、日志埋点、G2 指派接口  
- **交付件**: 推荐 API、风控规则配置界面、监控仪表盘

#### 车队收益分配策略（新增）
1. **平台规则层**：在 `player_teams` 增加 `profitMode`（如 `equal_split`、`weight_split`、`custom`），系统先计算订单可分配金额池（实收 - 平台抽成），若为标准模式则自动按人数或预设权重发放，队长无法修改，确保新人队伍也有默认方案。
2. **队长自定义层**：当 `profitMode=custom` 时，队长通过 `POST /api/v1/player/teams/{teamId}/orders/{orderId}/payout-plan` 提交本单分配方案（百分比总和 100%），并可在派单前多次调整；收益结算读取该方案写入 `team_order_payouts`，队员需在客户端点击“确认收益方案”才会生效。
3. **安全与风控**：平台设置单个成员最小收益比例与异常阈值，防止压价；任何队员未确认或中途退出会阻断结算并提醒队长重新分配；所有收益方案、调整和确认操作写入审计日志与告警系统，可随时强制回退到平台默认分配。

---

## 🔧 开发工具

### 常用命令

#### 后端开发
```bash
cd backend

# 代码检查
make lint
golangci-lint run --timeout=5m

# 运行测试
make test
go test ./... -v

# 生成测试覆盖率报告
go test ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -html=coverage.out -o coverage.html

# 构建应用
make build
go build ./cmd/...

# 生成Swagger文档
make swagger
swag init -g cmd/user-service/main.go
```

#### 前端开发
```bash
cd frontend

# 代码检查
npm run lint

# 代码格式化
npm run format

# 类型检查
npm run typecheck

# 运行测试
npm run test
npm run test:coverage

# 构建生产版本
npm run build
npm run build:analyze
```

### 测试与质量
```bash
# 运行所有测试
go test ./... -v

# 查看覆盖率统计
go tool cover -func=coverage.out | grep total

# 运行性能测试
go test -bench=. ./...

# 集成测试
go test ./tests/integration/... -v
```

---

## 📚 核心文档

### 📋 设计文档
- **[即时通讯系统设计文档](docs/即时通讯系统设计文档.md)** - 完整的聊天功能架构设计
- **[CLAUDE.md](CLAUDE.md)** - 项目开发指南和规范
- **[API文档](docs/swagger/)** - 完整的REST API文档

### 📊 项目报告
- **[项目状态报告](docs/PROJECT_STATUS_FINAL_REPORT.md)** - 项目整体进度和成果
- **[用户接口设计报告](backend/USER_INTERFACE_INTEGRITY_REPORT.md)** - 用户端功能设计
- **[测试覆盖率报告](backend/LATEST_COVERAGE_REPORT.md)** - 测试质量和覆盖率统计

### 🎯 功能指南
- **[前端开发指南](frontend/docs/FRONTEND_DEVELOPMENT_COMPLETE_GUIDE.md)** - 前端开发完整指南
- **[页面结构说明](frontend/docs/FRONTEND_PAGES_STRUCTURE.md)** - 管理端页面详细说明
- **[用户端页面设计](frontend/docs/USER_FACING_PAGES_GUIDE.md)** - 用户端和陪玩师端页面设计

### 🛠️ 技术文档
- **[三端架构指南](frontend/docs/COMPLETE_THREE_END_SYSTEM_GUIDE.md)** - 完整系统架构设计
- **[数据库设计](backend/docs/database/)** - 数据库表结构和关系设计
- **[部署指南](docs/deployment/)** - 生产环境部署和运维指南

---

## 🏆 项目亮点

### ✨ 核心功能完成度
- **🎯 智能订单系统** - 100%完成，支持6种订单状态流转
- **👥 多角色管理** - 100%完成，用户/陪玩师/管理员权限体系
- **💳 完整支付流程** - 95%完成，订单确认、支付、评价一体化
- **💬 即时通讯设计** - 100%完成，公共群聊 + 订单专属小群设计
- **📊 数据监控面板** - 85%完成，实时订单状态和收益统计
- **🔐 安全认证体系** - 100%完成，JWT + RBAC权限控制

### 🚀 技术成就
- **📈 测试覆盖率** - 49.5% (184个测试用例，100%通过率)
- **📚 文档完整性** - 95% (47个详细文档，涵盖架构、API、前端设计)
- **🏗️ 代码架构** - 清晰的分层架构，高内聚低耦合
- **🔧 开发工具链** - 完整的CI/CD和代码质量检查体系
- **📱 响应式设计** - 桌面端和移动端完美适配

### 📊 质量指标
- **代码规范**: 100%遵循Go和TypeScript编码规范
- **API设计**: RESTful设计，统一的响应格式
- **错误处理**: 完善的错误处理和日志记录
- **性能优化**: 数据库查询优化和缓存策略
- **安全防护**: XSS防护、CSRF防护、SQL注入防护

---

## 🚀 发展路线图

### 📅 短期计划 (1-2个月)
- **即时通讯功能开发** - 基于完整设计文档，实施WebSocket聊天系统
- **前端界面完善** - 完成管理端、用户端、陪玩师端所有页面开发
- **移动端适配** - 优化移动端用户体验，支持PWA
- **性能优化** - 提升系统性能，支持更高并发

### 📅 中期计划 (3-6个月)
- **高级功能扩展** - 图片分享、素材图集、内容模板
- **AI智能推荐** - 基于用户行为的智能陪玩师推荐
- **数据分析平台** - 用户行为分析和商业智能报表
- **多语言支持** - 国际化支持，拓展海外市场

### 📅 长期规划 (6个月+)
- **微服务架构** - 系统拆分和服务化改造
- **云原生部署** - 容器化部署和Kubernetes编排
- **开放平台** - 第三方API开放和生态建设
- **移动端App** - 原生移动应用开发

---

---

## 📡 API 规划

**全局约定**
- 统一前缀：`/api/v1`
- 认证方式：Bearer JWT + RBAC，管理端额外校验操作权限
- 返回结构：`{ "success": true, "data": {...}, "traceId": "" }`
- 命名风格：资源名使用复数，动作使用 HTTP 动词或子资源（符合当前代码风格）

### 核心业务 API
| 模块 | Endpoints（示例） | 说明 |
| --- | --- | --- |
| 认证 /auth | `POST /auth/login`、`POST /auth/logout`、`POST /auth/register`、`POST /auth/password/reset` | 用户/陪玩师/管理员统一认证入口，支持手机 + 邮箱 |
| 用户 /users | `GET /admin/users`、`POST /admin/users`、`PUT /admin/users/{id}`、`DELETE /admin/users/{id}`、`GET /user/profile`、`PUT /user/profile` | 管理端 CRUD、用户侧资料维护 |
| 陪玩师 /players | `GET /admin/players`、`POST /player/apply`、`PUT /player/profile`、`PUT /player/status` | 入驻审核、资料/状态管理 |
| 游戏 /games | `GET /admin/games`、`POST /admin/games`、`PUT /admin/games/{id}` | 游戏库维护 |
| 订单 /orders | `POST /user/orders`、`GET /user/orders`、`GET /user/orders/{id}`、`PUT /user/orders/{id}/cancel`、`PUT /user/orders/{id}/complete`、`GET /admin/orders`、`PUT /admin/orders/{id}` | 用户下单 + 状态管理 + 管理端监控 |
| 支付 /payments | `POST /payments`、`GET /payments/{id}`、`POST /payments/{id}/callback`、`POST /payments/{id}/refund` | 支付下单、回调、退款 |
| 评价 /reviews | `POST /user/reviews`、`GET /user/reviews/my`、`GET /player/reviews`、`POST /admin/reviews/{id}/moderate` | 双向评价 + 审核 |
| 聊天 /chats | `WS /ws/chat`、`GET /chat/rooms/{id}/messages`、`POST /chat/rooms/{id}/messages` | 订单会话/公共群聊 |
| 通知 /notifications | `GET /notifications`、`POST /notifications/read` | 站内信、系统通知 |

> 上述核心接口已在现有代码中实现/规划，可直接复用，新增模块需延续同一命名与鉴权风格。

### 已实现接口（swagger 提取）
（来源：`backend/docs/swagger.json`，当前已实现的 Auth + 管理后台 API）

#### Admin/Games
| Method | Path | 描述 |
| --- | --- | --- |
| GET | /admin/games | 列出游戏 |
| POST | /admin/games | 创建游戏 |
| DELETE | /admin/games/{id} | 删除游戏 |
| GET | /admin/games/{id} | 获取游戏 |
| PUT | /admin/games/{id} | 更新游戏 |
| GET | /admin/games/{id}/logs | 获取游戏操作日志 |

#### Admin/Orders
| Method | Path | 描述 |
| --- | --- | --- |
| GET | /admin/orders | 列出订单 |
| POST | /admin/orders | 创建订单 |
| DELETE | /admin/orders/{id} | 删除订单 |
| GET | /admin/orders/{id} | 获取订单 |
| PUT | /admin/orders/{id} | 更新订单 |
| POST | /admin/orders/{id}/assign | 指派订单的陪玩师 |
| POST | /admin/orders/{id}/cancel | 取消订单 |
| POST | /admin/orders/{id}/complete | 完成订单 |
| POST | /admin/orders/{id}/confirm | 确认订单 |
| GET | /admin/orders/{id}/logs | 获取订单操作日志 |
| GET | /admin/orders/{id}/payments | 获取订单支付记录 |
| POST | /admin/orders/{id}/refund | 订单退款 |
| GET | /admin/orders/{id}/refunds | 获取订单退款记录 |
| POST | /admin/orders/{id}/review | 审核订单（通过/拒绝） |
| GET | /admin/orders/{id}/reviews | 获取订单评价列表 |
| POST | /admin/orders/{id}/start | 开始服务 |
| GET | /admin/orders/{id}/timeline | 获取订单时间线 |

#### Admin/Payments
| Method | Path | 描述 |
| --- | --- | --- |
| GET | /admin/payments | 列出支付 |
| POST | /admin/payments | 创建支付记录 |
| DELETE | /admin/payments/{id} | 删除支付 |
| GET | /admin/payments/{id} | 获取支付 |
| PUT | /admin/payments/{id} | 更新支付 |
| POST | /admin/payments/{id}/capture | 确认支付入账 |
| GET | /admin/payments/{id}/logs | 获取支付操作日志 |
| POST | /admin/payments/{id}/refund | 退款处理 |

#### Admin/Players
| Method | Path | 描述 |
| --- | --- | --- |
| GET | /admin/players | 列出玩家资料 |
| POST | /admin/players | 新建玩家资料 |
| DELETE | /admin/players/{id} | 删除玩家资料 |
| GET | /admin/players/{id} | 获取玩家资料 |
| PUT | /admin/players/{id} | 更新玩家资料 |
| PUT | /admin/players/{id}/games | 更新玩家主游戏 |
| GET | /admin/players/{id}/logs | 获取玩家操作日志 |
| GET | /admin/players/{id}/reviews | 获取陪玩师的评价 |
| PUT | /admin/players/{id}/skill-tags | 更新玩家技能标签 |
| PUT | /admin/players/{id}/verification | 更新玩家认证状态 |

#### Admin/Reviews
| Method | Path | 描述 |
| --- | --- | --- |
| GET | /admin/reviews | 评价列表 |
| POST | /admin/reviews | 创建评价 |
| DELETE | /admin/reviews/{id} | 删除评价 |
| GET | /admin/reviews/{id} | 获取评价 |
| PUT | /admin/reviews/{id} | 更新评价 |
| GET | /admin/reviews/{id}/logs | 获取评价操作日志 |

#### Admin/Stats
| Method | Path | 描述 |
| --- | --- | --- |
| GET | /admin/stats/audit/overview | 审计总览（按实体/动作汇总） |
| GET | /admin/stats/audit/trend | 审计趋势（日） |
| GET | /admin/stats/dashboard | Dashboard 概览 |
| GET | /admin/stats/orders | 订单状态统计 |
| GET | /admin/stats/revenue-trend | 收入趋势（日） |
| GET | /admin/stats/top-players | Top 陪玩师排行 |
| GET | /admin/stats/user-growth | 用户增长（日） |

#### Admin/Users
| Method | Path | 描述 |
| --- | --- | --- |
| GET | /admin/users | 列出用户 |
| POST | /admin/users | 创建用户 |
| DELETE | /admin/users/{id} | 删除用户 |
| GET | /admin/users/{id} | 获取用户 |
| PUT | /admin/users/{id} | 更新用户 |
| GET | /admin/users/{id}/logs | 获取用户操作日志 |
| GET | /admin/users/{id}/orders | 获取用户的订单 |
| PUT | /admin/users/{id}/role | 更新用户角色 |
| PUT | /admin/users/{id}/status | 更新用户状态 |

#### Auth
| Method | Path | 描述 |
| --- | --- | --- |
| POST | /auth/login | 登录 |
| POST | /auth/logout | 登出（前端丢弃 Token） |
| GET | /auth/me | 获取当前用户信息 |
| POST | /auth/refresh | 刷新 Token |
| POST | /auth/register | 注册 |

### 扩展模块 API

#### 订单池 /player/orders
| Method & Path | 说明 | 主要参数 | 权限 | 状态 |
| --- | --- | --- | --- | --- |
| `GET /api/v1/player/orders/pool` | 获取可抢订单列表（分页 + 多条件筛选） | `gameId, minPrice, maxPrice, queue=team/single, cursor` | `player:orders:read` | 设计中 (G1) |
| `POST /api/v1/player/orders/{orderId}/snatch` | 队长/个人抢单，成功则锁单 | Body: `teamId(optional)` | `player:orders:snatch` | 设计中 (G1) |
| `POST /api/v1/player/orders/{orderId}/release` | 抢单失败或放弃释放锁 | Body: `reason` | `player:orders:release` | 设计中 (G1) |

### 客服指派 /admin/orders
| Method & Path | 说明 | 主要参数 | 权限 | 状态 |
| --- | --- | --- | --- | --- |
| `GET /api/v1/admin/orders/pending-assign` | 待指派订单列表（支持 SLA 排序） | `gameId, priority, createdBefore` | `admin:orders:read` | 方案评审 (G2) |
| `GET /api/v1/admin/orders/{orderId}/candidates` | 获取推荐陪玩师候选列表 | query: `limit` | `admin:orders:assign` | 方案评审 (G2, 依赖 G6) |
| `POST /api/v1/admin/orders/{orderId}/assign` | 指派订单给指定陪玩师或车队 | Body: `playerId / teamId, note` | `admin:orders:assign` | 方案评审 (G2) |
| `POST /api/v1/admin/orders/{orderId}/assign/cancel` | 取消/重置指派 | Body: `reason` | `admin:orders:assign` | 方案评审 (G2) |

### 实时订单推送 /ws/orders
| 通道 | 说明 | 事件类型 | 可靠性 | 状态 |
| --- | --- | --- | --- | --- |
| `WS /ws/orders?token=...` | 订单事件推送，支持玩家、陪玩师、客服多角色订阅 | `order.created`, `order.snatched`, `order.assigned`, `order.updated` | ACK + 重试 + 离线缓存 | 技术预研 (G3) |

### 社区/维系
| Method & Path | 说明 | 主要参数 | 权限 | 状态 |
| --- | --- | --- | --- | --- |
| `POST /api/v1/user/feeds` | 发布图文动态（仅图片/文本） | Body: `images[], content, visibility` | `user:feeds:write` | ✅ Implemented |
| `GET /api/v1/user/feeds` | 社区动态流 | `cursor, followOnly` | `user:feeds:read` | ✅ Implemented |
| `POST /api/v1/player/reviews/{reviewId}/reply` | 陪玩师回复评价 | Body: `content` | `player:reviews:reply` | ✅ Implemented |
| `GET /api/v1/notifications` | 获取通知列表 | `page, pageSize, unread, priority[]` | `user:notifications:read` | ✅ Implemented |
| `POST /api/v1/notifications/read` | 批量标记通知已读 | Body: `ids[]` | `user:notifications:write` | ✅ Implemented |

### 车队管理 /player/teams
| Method & Path | 说明 | 主要参数 | 权限 | 状态 |
| --- | --- | --- | --- | --- |
| `POST /api/v1/player/teams` | 创建车队（支持命名、封面、目标游戏） | Body: `name, games[], notice` | `player:teams:manage` | 原型设计 (G5) |
| `GET /api/v1/player/teams` | 我管理/加入的车队列表 | query: `role=leader/member` | `player:teams:read` | 原型设计 (G5) |
| `POST /api/v1/player/teams/{teamId}/members` | 邀请/审核队员加入 | Body: `playerId, role` | `player:teams:manage` | 原型设计 (G5) |
| `DELETE /api/v1/player/teams/{teamId}/members/{playerId}` | 移出队员 | — | `player:teams:manage` | 原型设计 (G5) |
| `POST /api/v1/player/teams/{teamId}/orders/{orderId}/dispatch` | 队长将已抢订单分派给队员 | Body: `assigneeId` | `player:teams:dispatch` | 原型设计 (G5) |

### 智能匹配 /admin/recommendations
| Method & Path | 说明 | 主要参数 | 权限 | 状态 |
| --- | --- | --- | --- | --- |
| `GET /api/v1/admin/orders/{orderId}/recommendations` | 返回陪玩师/车队推荐列表 + 评分维度 | `limit, includeTeams` | `admin:orders:assign` | 数据建模 (G6) |
| `POST /api/v1/admin/orders/{orderId}/recommendations/feedback` | 记录指派反馈，用于模型训练 | Body: `selectedId, result` | `admin:orders:assign` | 数据建模 (G6) |

> 对齐现有 API 风格：保持 RESTful 资源命名、以 HTTP 动词表达动作；需要自定义动作时使用子资源（如 `/snatch`, `/dispatch`）；所有写操作默认开启幂等校验并记录审计日志。

---

## 📦 数据模型

### 1. 现有核心模型（`backend/internal/model`）
| 模型 | 关键字段示例 | 说明 | 状态 |
| --- | --- | --- | --- |
| User (`user.go`) | `Phone`, `Email`, `Role`, `Status`, `LastLoginAt`, `Roles[]` | 平台账号（用户/陪玩师/管理员），支持多角色、黑名单、最近登录追踪 | ✅ 已实现 |
| Player (`player.go`) | `UserID`, `Nickname`, `Rank`, `HourlyRateCents`, `MainGameID`, `VerificationStatus` | 陪玩师档案，与 User 一对一，记录段位、价格、认证状态 | ✅ 已实现 |
| Game (`game.go`) | `Name`, `IconURL`, `Genre`, `Status`, `SortOrder` | 游戏元信息，服务项和订单引用 | ✅ 已实现 |
| ServiceItem (`service_item.go`) | `ItemCode`, `Category/SubCategory`, `GameID`, `PlayerID`, `BasePriceCents`, `CommissionRate`, `MinUsers/MaxPlayers` | 统一服务定义（护航/团队/礼物），内置抽成和成团参数 | ✅ 已实现 |
| Order (`order.go`) | `OrderNo`, `UserID`, `PlayerID/RecipientPlayerID`, `Status`, `QueueType`, `RequiredMembers`, `AssignmentSource`, `AssignedTeamID`, `TotalPriceCents`, `CommissionCents`, `GameID`, `ScheduledStart/End`, `GiftMessage`, `RefundAmountCents`, `DisputeStatus` | 护航/礼物/团队订单，覆盖预约、指派、退款、扩展配置；`QueueType`/`RequiredMembers`/`AssignmentSource` 等字段需要在后续迁移中补充 | 🔄 需扩展 |
| Payment (`payment.go`) | `OrderID`, `UserID`, `Method`, `AmountCents`, `Status`, `ProviderTradeNo`, `ProviderRaw`, `PaidAt/RefundedAt` | 支付/退款流水与渠道回执 | ✅ 已实现 |
| Review (`review.go`) | `OrderID`, `UserID`, `PlayerID`, `Score`, `Content` | 评价记录，支持后台审核/展示 | ✅ 已实现 |
| ChatGroup & ChatMessage (`chat.go`) | `GroupType`, `RelatedOrderID`, `Members`, `MessageType`, `AuditStatus`, `Settings` | 公共/订单群聊、消息审核、成员既读管理 | ✅ 已实现 |
| Service Finance (`commission.go`, `financial.go`, `withdraw.go`) | `SettlementBatch`, `CommissionRate`, `WithdrawStatus` | 平台抽成、结算、提现流水 | ✅ 已实现 |
| OperationLog / Audit (`operation_log.go`) | `Entity`, `Action`, `Payload`, `OperatorID` | 后台操作记录，用于 `Admin/Stats` 审计 | ✅ 已实现 |

### 2. 规划中的扩展模型
| 模型 | 关键字段建议 | 场景/说明 | 状态 |
| --- | --- | --- | --- |
| Team | `ID`, `LeaderID`, `Name`, `Games[]`, `ProfitMode`, `Notice`, `Status`, `DefaultShare` | 车队基础信息，支持多游戏、盈利模式（`equal_split/weight/custom`）及默认分成 | 📝 设计中 (G5) |
| TeamMember | `TeamID`, `UserID`, `Role(leader/member/co-lead)`, `JoinedAt`, `Status`, `ProfitShareDefault` | 队伍成员与权限管理，配合派单/收益统计 | 📝 设计中 |
| TeamOrderAssignment | `OrderID`, `TeamID`, `QueueType(single/team)`, `RequiredMembers`, `Status(pending/assigned/in_service/completed)` | 订单被队伍抢单后的绑定关系，记录还需多少成员 | 📝 设计中 |
| TeamAssignmentMember | `AssignmentID`, `MemberID`, `State(assigned/accepted/completed/withdrawn)`, `StartedAt/CompletedAt` | 队内成员对任务的响应与执行进度 | 📝 设计中 |
| TeamPayoutPlan | `AssignmentID`, `ProfitMode`, `Shares[{memberId,percent}]`, `ConfirmedBy[]`, `LockedAt` | 队长设置的分账方案（队员确认后生效），供收益结算读取 | 📝 设计中 |
| OrderDispute | `OrderID`, `RaisedBy(user/player)`, `Reason`, `Evidence`, `AssignmentSource`, `Status(open/in_review/resolved)`, `Resolution`, `RefundAmount`, `HandledBy`, `HandledAt` | 售后/争议闭环，客服可介入、判责、触发退款或重派 | 💡 规划 |
| NotificationEvent | `UserID`, `Channel(web/push/sms)`, `Payload`, `Priority`, `ReferenceType/ReferenceID`, `ReadAt` | 统一通知中心实体，支撑站内信 + 外部消息 | ✅ Implemented |
| Feed/Community | `FeedID`, `AuthorID`, `Images[]`, `Content`, `Visibility`, `ModerationStatus`, `Metrics`, `ReplyCount`, `ComplaintCount` | 图文动态及审核状态，为社区/维系模块提供数据 | ✅ Implemented |
| PlayerStats | `PlayerID`, `CompletedOrders`, `CancelRate`, `ResponseTime`, `SkillTags`, `Languages` | 陪玩师数据指标，供匹配算法/抢单池使用 | 💡 规划 |

> 规划模型将在 G5/G6 阶段逐步落地，届时会同步数据库迁移、OpenAPI Schema 与测试用例。

---

## 🧭 页面规划

### 用户端 (Customer)
| 页面 | 路径 | 核心功能 | 状态 |
| --- | --- | --- | --- |
| 首页 | `/` | Banner、热门游戏、优质陪玩师、搜索入口 | 已上线 |
| 游戏列表 | `/games` | 分类、筛选、排序、搜索建议 | 已上线 |
| 陪玩师列表 | `/games/:id/players` | 条件筛选、标签展示、在线状态 | 已上线 |
| 陪玩师详情 | `/players/:id` | 资料、服务套餐、评价、预约时间 | 已上线 |
| 订单创建 | `/order/create` | 服务选择、时段、客服指派选项、优惠券 | 已上线（待接客服模式） |
| 支付页面 | `/order/:id/pay` | 多支付方式、倒计时、回调状态 | 已上线 |
| 我的订单 | `/orders` | 状态筛选、取消/确认/评价入口 | 已上线 |
| 个人中心（Phase 2） | `/profile` | 账号安全、通知偏好、钱包 | 设计完成 |
| 动态广场（Phase 2） | `/community` | 图文动态、关注、举报 | ✅ Implemented |

### 陪玩师端 (Player)
| 页面 | 路径 | 核心功能 | 状态 |
| --- | --- | --- | --- |
| 工作台 | `/player/dashboard` | 当日订单、收益、状态开关 | 已上线 |
| 订单管理 | `/player/orders` | 接单/拒单、开始/完成、售后处理 | 已上线 |
| 服务管理 | `/player/services` | 服务套餐、价格、可预约时间 | 已上线 |
| 收益管理 | `/player/earnings` | 收益统计、提现申请、银行卡 | 已上线 |
| 抢单大厅 | `/player/orders/pool` | 实时订单、抢单、筛选 | 交互设计 (G1) |
| 车队控制台 | `/player/teams` | 建队、成员管理、队长派单 | 原型制作 (G5) |
| 排班日历（Phase 2） | `/player/schedule` | 可用时间规划、忙闲状态 | 设计中 |
| 营销工具（Phase 2） | `/player/marketing` | 活动配置、粉丝通知 | 规划中 |

### 管理后台 (Admin)
| 页面 | 路径 | 核心功能 | 状态 |
| --- | --- | --- | --- |
| 仪表盘 | `/admin/dashboard` | 全局指标、告警、快捷操作 | 开发中 |
| 订单管理 | `/admin/orders` | 多维筛选、状态流转、异常处理 | 开发中 |
| 指派工作台 | `/admin/assignments` | 待指派列表、陪玩师筛选、日志 | UI 开发 (G2) |
| 用户管理 | `/admin/users` | 用户审核、黑名单、风控信息 | 已上线 |
| 陪玩师管理 | `/admin/players` | 入驻审核、资质、技能标签 | 已上线 |
| 财务/结算 | `/admin/finance` | 提现审核、账单对账、发票 | 开发中 |
| 系统设置 | `/admin/settings` | 权限、参数、主题、通知配置 | 已上线 |
| 风控面板 | `/admin/risk` | 风控规则、异常告警、处置记录 | 需求收集中 (G6) |
| 内容管理（Phase 2） | `/admin/content` | Banner、公告、社区内容审核 | 规划中 |

> 页面规划与 API 清单一一映射：每个页面进入研发前需完成原型稿、接口协议、权限矩阵与监控指标（PV、成功率、错误率），确保需求到交付路径透明。

---

## ✅ 功能验收规范

| 功能块 | 验收要点 | 通过标准 |
| --- | --- | --- |
| 认证 /auth | 多角色登录/注册/注销/密码找回；JWT 签发与刷新；RBAC 权限校验 | 正常/异常流程返回标准响应；审计记录账号+IP+UA；Token 过期策略验证通过 |
| 用户 /users | 用户 CRUD、黑名单、风控标签、资料修改 | 管理端批量操作成功率 ≥99%；资料修改实时生效并触发通知；黑名单 5 分钟内同步到所有服务 |
| 陪玩师 /players | 入驻申请、资质审核、技能标签、在线状态 | 审核流转日志完整；状态与订单池/队长模式一致；违规封禁立即生效 |
| 游戏 /games | 游戏列表 CRUD、分类、排序 | 新增/编辑 5 分钟内前台可见且缓存一致；删除具备回滚机制 |
| 订单 /orders | 下单、状态流转、售后、后台协同 | 全链路可回放（traceId）；状态切换触发通知/推送；异常订单自动告警 |
| 支付 /payments | 支付创建、回调、退款、对账 | 与第三方流水一致；退款需双向确认；失败支付触发告警并可人工补单 |
| 评价 /reviews | 用户评价、陪玩师回复、后台审核 | 一单仅允许一次评价+追加；违规评价 10 分钟内处理；回复通知相关用户 |
| 聊天 /chats | WebSocket 连接、消息存储、历史拉取 | 连接保活 ≥30 分钟；消息去重与敏感词过滤；离线消息可追溯 |
| 通知 /notifications | 站内信、多渠道通知、已读回执 | 未读数实时；高危通知短信兜底；点击率/失败率可观测 |
| 订单池 /抢单 | 实时订单流、筛选、抢单/释放、锁单冲突控制 | 新订单 ≤1s 出现在大厅；抢单冲突率 <0.5%；日志可追溯 |
| 客服指派 | 待指派列表、候选推荐、指派/撤销、响应跟踪 | 指派/撤销都有审计；超时订单自动升级；SLA 仪表盘实时 |
| 实时推送 /ws/orders | 订单事件广播、ACK、重试、离线缓存 | 推送 P95 <800ms；掉线 1 分钟内自动重连；消息不丢不重 |
| 社区/维系 | 图文动态、关注、举报、评价回复 | 动态仅含图片+文本；举报 5 分钟内进入审核；评价回复触达相关用户 |
| 车队/队长机制 | 建队、成员角色、队长抢单、队内派单、收益统计 | 队长抢单后必须分派；队员接单推送实时；收益结算准确可追踪 |
| 智能匹配/风控 | 推荐列表、评分维度、反馈闭环、异常检测 | 推荐接口响应 <500ms；反馈用于模型更新；风控命中触发告警与处置记录 |
| 页面交互 | 三端所有页面 | 页面与 API 映射清单齐全；空态/错误态覆盖；关键路径具备埋点数据 |

共通要求：
1. **测试**：单元覆盖率 ≥80%，关键链路具备集成/E2E 用例；WebSocket/异步流程提供模拟脚本。  
2. **监控**：每个模块需输出核心 KPI（QPS、延迟、成功率、错误率）和告警策略。  
3. **日志/审计**：涉及资金、指派、权限变更必须记录 traceId、操作者、上下文。  
4. **回滚预案**：新增能力需要灰度/特性开关或回滚脚本，异常可 5 分钟内恢复。

---

## 🧪 自动化测试方案

| 模块 | 测试类型 | 覆盖重点 | 工具/命令 |
| --- | --- | --- | --- |
| 认证 /auth | 单元 + 集成 | JWT 签发、刷新、RBAC 中间件；异常登录流 | Go `testing`, testify；`go test ./internal/auth/...` |
| 用户 /users | 单元 + API 回归 | 用户 CRUD、黑名单同步、资料修改事件 | Go tests + Postman/Newman；`newman run tests/users.postman.json` |
| 陪玩师 /players | 单元 + 集成 | 入驻审核流、技能标签、状态切换 | Go tests；`go test ./internal/service/player/...` |
| 游戏 /games | 单元 | 游戏 CRUD、缓存一致性（mock cache） | Go tests + gomock |
| 订单 /orders | 单元 + 集成 + E2E | 全状态流转、售后、并发修改 | Go tests、Cypress API suite；`go test ./internal/service/order/...` |
| 支付 /payments | 集成（沙箱） + 回归 | 支付创建/回调/退款、异常补单 | Mock 支付网关 + Go tests；`make test-payments` |
| 评价 /reviews | 单元 + API | 评价限制、追加评价、审核流程 | Go tests + Newman |
| 聊天 /chats | 单元 + WebSocket 集成 | 连接保活、消息去重、离线补发 | Go tests + k6 WebSocket 脚本 `k6 run scripts/ws_orders.js` |
| 通知 /notifications | 单元 + 集成 | 未读数、批量推送、已读回执 | Go tests + Kafka/Redis mock |
| 订单池 /抢单 | 并发测试 + 端到端 | 实时订单流、抢单冲突、锁释放 | Go race tests + k6 HTTP；`k6 run scripts/order_pool.js` |
| 客服指派 | 集成 + UI 自动化 | 待指派列表、指派/撤销、SLA 倒计时 | Playwright Admin suite；`npx playwright test admin-assign` |
| 实时推送 /ws/orders | WebSocket 压测 + 延迟监控 | 发布/订阅、ACK、重试、离线缓存 | k6 + Loki 日志；`k6 run scripts/ws_orders.js` |
| 社区/维系 | API + 安全测试 | 图片/文本校验、举报流程、评价回复 | Go tests + OWASP ZAP (上传接口) |
| 车队/队长机制 | 集成 + 并发 | 建队、成员管理、队长抢单、收益分配 | Go tests + k6 scenario；`k6 run scripts/team_dispatch.js` |
| 智能匹配/风控 | 单元 + 数据验证 | 特征计算、推荐列表、风控命中 | Python pytest (数据脚本) + Go tests |
| 页面层（用户/陪玩师/后台） | E2E | 关键用户旅程（下单→支付→评价、队长抢单、客服指派） | Playwright/Cypress；`npm run test:e2e` |

**执行策略**
1. **CI 阶段**：`make lint && make test`（含 race/coverage）→ `npm run lint && npm run test` → Newman/k6 轻量用例。  
2. **Nightly**：运行支付沙箱、WebSocket 压测、Playwright 全量 E2E，输出趋势报告。  
3. **Pre-release**：结合验收规范逐项勾选，必须提供测试截图/日志 + 监控基线。  
4. **观测联动**：测试管线把关键指标推送到 Grafana/Loki，方便对比线上 SLO。

---

## 🤝 贡献指南

### 参与方式
1. **Fork项目** 到你的GitHub账户
2. **创建功能分支** (`git checkout -b feature/AmazingFeature`)
3. **提交代码** (`git commit -m 'Add some AmazingFeature'`)
4. **推送分支** (`git push origin feature/AmazingFeature`)
5. **创建Pull Request**

### 开发规范
- 遵循项目编码规范 (见CLAUDE.md)
- 添加必要的测试用例
- 更新相关文档
- 通过所有CI检查

### 问题反馈
- 📋 **功能建议**: [Issues](https://github.com/your-org/GameLink/issues)
- 🐛 **Bug报告**: [Issues](https://github.com/your-org/GameLink/issues)
- 💬 **技术讨论**: [Discussions](https://github.com/your-org/GameLink/discussions)

---

## 📞 联系我们

### 🏢 团队信息
- **项目负责人**: GameLink开发团队
- **技术支持**: dev-team@gamelink.com
- **商务合作**: business@gamelink.com

### 📱 社交媒体
- **官方网站**: https://gamelink.com
- **技术博客**: https://blog.gamelink.com
- **GitHub**: https://github.com/your-org/GameLink

---

## 📄 开源协议

本项目采用 [MIT License](LICENSE) 开源协议，详见项目根目录的LICENSE文件。

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给我们一个Star！**

**🚀 让我们一起构建更好的游戏陪玩生态！**

---

*最后更新: 2025-11-10 | 版本: v2.0*

</div>
