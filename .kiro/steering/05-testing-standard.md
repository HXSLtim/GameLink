# 测试规范

## 核心原则

测试不是"看界面"，是"验逻辑"。必须从浏览器到数据库的每个环节都有证据证明正常工作。

## 测试三层次

| 层次 | 内容 | 重要性 |
|------|------|--------|
| Level 1 | 前端表现层（页面渲染、按钮可交互） | ⚠️ 仅此不合格 |
| Level 2 | 前后端联调层（请求/响应/数据结构） | ⭐ 重点 |
| Level 3 | 数据与逻辑层（数据库变更、业务逻辑、异常处理） | ⭐ 重点 |

## Docker 环境测试流程

### 1. 容器状态检查

```bash
docker compose ps                    # 所有容器运行状态
docker stats                         # CPU/内存占用
docker ps --format "table {{.Names}}\t{{.RestartCount}}"  # 重启次数应为0
```

### 2. 日志监控

```bash
docker logs -f gamelink-backend      # 后端日志
docker logs gamelink-backend | grep -i error  # 检查错误
```

### 3. 数据库验证

```bash
docker exec -it gamelink-postgres psql -U gamelink -d gamelink
# 执行 SQL 查询验证数据
```

## 测试检查清单

### 基础功能测试
- [ ] 正常流程走完（请求→响应→数据→反馈）
- [ ] 接口 URL、Method 正确
- [ ] 请求参数完整准确
- [ ] 响应数据符合文档
- [ ] 数据库数据一致
- [ ] 页面渲染正确

### 异常场景测试（必测）
- [ ] 网络中断/慢网络
- [ ] 请求超时
- [ ] 后端返回错误码（403/500）
- [ ] 参数缺失或格式错误
- [ ] 业务逻辑异常
- [ ] 多次快速点击（防抖测试）

## 问题发现后的五问深究法

1. **现象是什么？**（截图+文字描述）
2. **哪一层的问题？**（前端/后端/网络/数据）
3. **具体哪个环节出错？**（请求没发/参数错误/返回异常/渲染失败）
4. **错误根因是什么？**（代码第几行？逻辑哪里不对？）
5. **如何修复和验证？**（修改方案+验证步骤）

## 测试报告模板

```markdown
## 功能测试报告：[功能名称]

### 1. 测试场景
[描述测试的功能]

### 2. 联调验证
- [ ] 请求发送：[METHOD] [URL] ✓/✗
- [ ] 请求参数：{...} ✓/✗
- [ ] 响应状态：[HTTP状态码], code: [业务码] ✓/✗
- [ ] 数据库：[表名]变更记录 ✓/✗
- [ ] 页面反馈：[描述] ✓/✗

### 3. 异常测试
| 场景 | 预期 | 实际 | 结果 |
|------|------|------|------|
| ... | ... | ... | ✓/✗ |

### 4. 容器日志
[关键日志截取]
```

## 禁止行为

- ❌ 只看 UI 不抓包
- ❌ 只测成功场景
- ❌ 发现问题不追根
- ❌ 不验证数据库
- ❌ 不监控容器日志就提交测试通过

## 按钮测试要点

每个按钮必须验证：

1. **静态检查**：按钮可见性、disabled 状态
2. **点击事件**：Network 请求、参数正确性
3. **后端处理**：容器日志、业务逻辑执行
4. **数据持久化**：数据库写入、缓存更新
5. **响应返回**：HTTP 状态码、业务 code、前端处理

## 单个按钮测试清单

```
□ 按钮可见性测试
□ 发送网络请求（Network 面板）
□ 请求到达后端（docker logs）
□ 业务逻辑执行无异常
□ 数据库数据正确写入
□ 后端返回响应正确
□ 前端正确处理响应
□ 异常场景测试（容器停止、网络断开）
```

## 测试覆盖率目标

- 当前：~80%（服务层平均）
- 目标：80%+
- 单元测试：表驱动测试，mock 依赖
- 集成测试：真实数据库（PostgreSQL），测试 fixtures
- 并发测试：Race detector，压力测试
- **E2E测试**：Playwright自动化测试覆盖关键业务流程（127个测试用例）

### 服务层覆盖率详情（2025-12-25）

| 模块 | 覆盖率 | 状态 |
|------|--------|------|
| menu | 100.0% | ✅ |
| handler/resp | 96.0% | ✅ |
| item | 90.9% | ✅ |
| sensitiveword | 88.9% | ✅ |
| audit | 86.9% | ✅ |
| settlementcompany | 86.5% | ✅ |
| player | 86.8% | ✅ |
| role | 84.9% | ✅ |
| user | 84.5% | ✅ |
| auth | 84.3% | ✅ |
| permission | 84.2% | ✅ |
| withdraw | 81.4% | ✅ |
| review | 80.3% | ✅ |
| wallet | 80.0% | ✅ |
| order | 79.2% | ⚠️ |
| content | 76.2% | ⚠️ |
| payment | 75.8% | ⚠️ |
| routingrule | 50.7% | ⚠️ |
| admin | 1.2% | ❌ |

### 待提升模块

- `routingrule`: RoutingEngine 复杂接口需要更多测试
- `admin`: 需要添加单元测试
- `order`: 接近 80%，可继续提升
- `payment`: 接近 80%，可继续提升

## E2E 测试（前端自动化）

### 概述
使用 Playwright 进行管理后台的端到端测试，覆盖关键业务流程。

### 测试覆盖

| 模块 | 测试文件 | 测试数量 | 状态 |
|------|---------|---------|------|
| 认证 | `auth.spec.ts` | 19 | ✅ |
| 用户管理 | `user-management.spec.ts` | 25 | ✅ |
| 订单管理 | `order-management.spec.ts` | 20 | ✅ |
| 支付管理 | `payment-management.spec.ts` | 22 | ✅ |
| 陪玩师管理 | `player-management.spec.ts` | 28 | ✅ |
| **总计** | - | **114** | ✅ |

### 运行 E2E 测试

```bash
cd admin

# 安装 Playwright 浏览器
npm run test:e2e:install

# 运行所有 E2E 测试（headless）
npm run test:e2e

# 运行测试并显示浏览器
npm run test:e2e:headed

# 调试模式
npm run test:e2e:debug

# 查看测试报告
npm run test:e2e:report
```

### E2E 测试特点

1. **Page Object Model**: 使用页面对象模式，提高测试可维护性
2. **自动清理**: 测试数据自动创建和清理，避免污染环境
3. **并发执行**: 支持多线程并行测试，提升执行速度
4. **详细报告**: 失败时自动截图和录屏，便于调试
5. **API 辅助**: 提供 API 辅助函数用于测试数据准备

### Page Objects

- `LoginPage`: 登录页面的所有交互
- `UserManagementPage`: 用户管理 CRUD 操作
- `OrderManagementPage`: 订单查看、取消、退款
- `PaymentManagementPage`: 支付记录查看、退款处理
- `PlayerManagementPage`: 陪玩师管理、审核流程

### 测试数据管理

- 使用 Fixtures 生成唯一测试数据（带时间戳）
- 每个测试独立运行，互不影响
- `afterEach` 钩子自动清理测试数据
- 支持通过 API 直接创建测试数据

### 环境要求

- 后端 API 运行在 `http://localhost:8080`
- 管理后台运行在 `http://localhost:5173`（自动启动）
- 管理员账号：`admin` / `admin123`（可通过环境变量配置）

### CI/CD 集成

测试已集成到 CI/CD 流程，并执行以下**质量门策略**：

**关键步骤（失败会阻塞构建）**:
- Linter 检查（后端 golangci-lint，前端 ESLint）
- 单元测试执行（后端 + 前端）
- 覆盖率检查（低于 70% 构建失败）
- 类型检查（前端 TypeScript）

**非关键步骤（失败不阻塞构建）**:
- 覆盖率报告上传（Codecov）
- SARIF 结果上传（GitHub Security）

E2E 测试可集成到 CI/CD 流程：

```yaml
- name: Run E2E tests
  run: |
    cd admin
    npm run test:e2e
```

详细文档：[admin/tests/e2e/README.md](../../admin/tests/e2e/README.md)
