# 快速参考 - GameLink 项目

> 🚀 **面向 AI 助手的精简版项目指南**
>
> 完整文档见 [00-INDEX.md](./00-INDEX.md)

---

## 🎯 项目概览

**GameLink** = 游戏陪玩服务平台（用户下单 + 陪玩师接单 + 平台管理）

| 维度 | 说明 |
|------|------|
| **核心业务** | 用户支付 → 陪玩师服务 → 平台抽成（15-25%） |
| **用户角色** | User（用户）、Player（陪玩师）、Admin（管理员） |
| **技术栈** | Go (Gin) + React 19 (Ant Design 6) + PostgreSQL + Redis |
| **完成度** | 后端 100% (36/36 模块) | 前端 75% | 测试覆盖率 ~80% |

---

## 📂 关键目录（必记）

```
GameLink/
├── api/                      # Go 后端
│   ├── internal/
│   │   ├── handler/          # HTTP 处理器（按模块分：admin/user/player）
│   │   ├── service/          # 业务逻辑 ⭐ 核心层
│   │   ├── repository/       # 数据访问
│   │   └── model/            # 数据模型（见 04-data-models.md）
│   └── pkg/                  # 公共包（auth/cache/config/db）
├── admin/                    # React 管理后台
│   └── src/
│       ├── pages/admin/      # 页面组件
│       ├── api/              # API 客户端
│       └── components/       # 公共组件
└── .kiro/steering/           # 📚 本文档目录
```

---

## 🔥 核心业务规则（必读）

### 抽成计算（三维度）

```
最终抽成 = 项目基础抽成 - 陪玩师个人调整 - 排名减免

示例：¥100 订单，20% 基础抽成
  → 项目设置 20%（默认）
  → 陪玩师个人可能设置 15-25%（覆盖项目）
  → 排名前10可能减免 5%（激励机制）
  → 最终：15% 抽成 = ¥15 平台收入，¥85 陪玩师收入
```

> 详见：[04-data-models.md](./04-data-models.md) 的 CommissionRule 模型

### 订单类型

| 类型 | 说明 | 支付时机 |
|------|------|----------|
| **solo** | 单陪玩师 | 下单即支付 |
| **team** | 多陪玩师（2+） | 所有陪玩师确认后支付 |
| **gift** | 直接打赏（无服务） | 立即支付，无 T+7 |

### 收入结算（T+7 规则）

```
订单完成 → 收益入账到 FrozenCents（冻结）
        ↓
      7 天等待期
        ↓
    无纠纷/退款 → FrozenCents → BalanceCents（可提现）
```

> 详见：[04-data-models.md](./04-data-models.md) 的 Wallet 模型

---

## 💻 命名规范（速查）

### Go 后端

| 类型 | 规范 | 示例 |
|------|------|------|
| 文件 | **camelCase**（小驼峰） | `userService.go` ❌ `user_service.go` |
| 类型 | PascalCase | `UserService` |
| 函数（导出） | PascalCase | `CreateOrder` |
| 函数（私有） | camelCase | `validateInput` |
| 测试文件 | `*_test.go` | `userService_test.go` |

### React 前端

| 类型 | 规范 | 示例 |
|------|------|------|
| 组件 | PascalCase | `UserTable.tsx` |
| 工具函数 | camelCase | `formatDate.ts` |
| 类型/接口 | PascalCase | `UserResponse` |
| 常量 | UPPER_SNAKE_CASE | `API_BASE_URL` |

---

## 🧪 测试规范（三层次）

> ⚠️ **仅 Level 1（前端渲染）是不够的！**

| 层次 | 内容 | 验证方法 |
|------|------|----------|
| **Level 1** | 页面渲染 | 组件挂载、按钮可点击 |
| **Level 2** ⭐ | 前后端联调 | 请求/响应格式、HTTP 状态码 |
| **Level 3** ⭐ | 数据逻辑 | 数据库变更、业务规则正确性 |

**测试检查清单**：
- [ ] 请求发送正确（URL/Method/参数）
- [ ] 响应格式符合文档（`{ success, data, code, message }`）
- [ ] 数据库变更正确（表记录更新）
- [ ] 异常场景处理（403/500/超时）

> 完整清单：[05-testing-standard.md](./05-testing-standard.md)

---

## 🚨 常见陷阱（避免）

### ❌ 错误做法

1. **命名使用 snake_case** → Go 文件应使用 `userService.go` 而非 `user_service.go`
2. **只测前端不测后端** → Level 1 测试不足，必须验证数据库
3. **忽略文档直接开发** → 导致业务逻辑错误（如抽成计算）
4. **忘记更新 steering 文档** → 模型变更后必须同步更新 `04-data-models.md`

### ✅ 正确流程

```
1. 查阅 steering 文档（理解业务规则）
   ↓
2. 编写代码（遵循命名规范）
   ↓
3. 编写测试（Level 2 + Level 3）
   ↓
4. 更新文档（如有模型变更）
   ↓
5. 提交代码（Conventional Commits）
```

---

## 📦 CI/CD 质量门

| 检查项 | 阈值 | 失败后果 |
|--------|------|----------|
| 测试覆盖率 | **≥70%** | 构建失败 ❌ |
| Linter（Go） | 无错误 | 构建失败 ❌ |
| Linter（前端） | 无错误 | 构建失败 ❌ |
| 安全扫描 | 无高危漏洞 | 构建失败 ❌ |

> 详情：[02-tech-stack.md](./02-tech-stack.md)

---

## 🔍 如何快速找到...

### 某个数据模型的定义？

**核心模块** → [04-data-models.md](./04-data-models.md)
**营销模块**（VIP/优惠券） → [04a-marketing-models.md](./04a-marketing-models.md)
**团队模块** → [04b-team-models.md](./04b-team-models.md)
**通知模块** → [04d-notification-models.md](./04d-notification-models.md)

### 某个功能放在哪个目录？

**后端 Handler**：`api/internal/handler/{admin|user|player}/`
**后端 Service**：`api/internal/service/`
**前端页面**：`admin/src/pages/admin/{Module}/`

### 测试用例怎么写？

参考：`api/internal/service/integration/testdb.go` 中的 30+ 测试辅助函数

---

## 📚 文档导航

| 需求 | 文档 |
|------|------|
| **开发新功能** | 03-project-structure.md → 04-data-models.md → 05-testing-standard.md |
| **查询进度** | 06-project-management.md 或 [../docs/PROGRESS.md](../docs/PROGRESS.md) |
| **部署上线** | 02-tech-stack.md（CI/CD 章节） |
| **数据库变更** | 04c-enums-indexes.md（枚举/索引） |

> 📖 **完整文档索引**：[00-INDEX.md](./00-INDEX.md)
