# GameLink API 变更日志

> **版本**: v2.0
> **最后更新**: 2026-02-11
> **遵循**: [SemVer](https://semver.org/lang/zh-CN/)

---

## 变更类型说明

| 类型 | 说明 | 示例 |
|------|------|------|
| `Added` | 新增功能 | 新增 API 接口 |
| `Changed` | 功能变更 | 修改请求参数 |
| `Deprecated` | 即将废弃 | 标记为废弃 |
| `Removed` | 已移除 | 删除 API 接口 |
| `Fixed` | 问题修复 | 修复 Bug |
| `Security` | 安全修复 | 安全漏洞修补 |

---

## [Unreleased]

### 待发布

#### Added
- 计划添加批量操作接口（用户管理）

#### Changed
- 计划调整仪表盘接口路径

---

## [5.1.0] - 2026-02-09

### Added
- 小程序组件化完成（133个组件）
- API 路径对齐（16处修正）
- 主题系统完善（CSS 变量 + 暗色模式）

### Changed
- 全量替换 Emoji 为 uv-icon
- 硬编码颜色清理为 CSS 变量
- Mock 数据替换为真实 API 调用

### Fixed
- 修复 API 路径不一致问题（5处）
- 修复参数校验问题（5处）

---

## [5.0.0] - 2026-01-28

### Added
- 后端所有模块完成（57个 Service）
- 管理后台核心功能完成（40+页面）
- WebSocket 实时通讯完成
- 监控系统集成（Prometheus + Grafana）

### Changed
- **[Breaking]** API 响应字段统一为 camelCase
- **[Breaking]** 统一使用 `/api/v1` 前缀

### Fixed
- 修复 JWT Token 刷新问题
- 修复订单状态流转问题

---

## [4.0.0] - 2026-01-15

### Added
- Phase 4 架构改进：
  - UnitOfWork 事务管理
  - AdminDeps 依赖注入结构
  - 统一 DTO 模型

### Changed
- **[Breaking]** Handler 函数签名变更
- **[Breaking]** Service 初始化方式变更

---

## [3.0.0] - 2026-01-01

### Added
- 支付系统（微信、支付宝）
- 钱包系统
- 充值与提现
- 佣金结算系统

### Changed
- **[Breaking]** 订单创建流程变更
- **[Breaking]** 支付回调接口路径变更

---

## [2.0.0] - 2025-12-15

### Added
- 陪玩师管理系统
- 陪玩师认证（实名、段位）
- 服务管理
- 在线状态管理

### Changed
- **[Breaking]** 用户模型拆分为 User + Player

---

## [1.0.0] - 2025-12-01

### Added
- 基础用户系统
- 订单系统
- 聊天系统
- 管理后台框架
- 小程序基础框架

---

## API 接口变更详情

### 认证相关 (/api/v1/auth)

#### v5.1.0
- **[Fixed]** 修复 Token 刷新时返回格式不一致问题

#### v5.0.0
- **[Changed]** `/register` 请求参数变更
  ```diff
  POST /api/v1/auth/register
  {
  -   "user_name": "string",
  +   "name": "string",
      "email": "string",
      "password": "string"
  }
  ```

#### v4.0.0
- **[Changed]** `/login` 响应格式统一为 camelCase
  ```json
  {
    "success": true,
    "data": {
  -   "user_id": 1,
  +   "userId": 1,
  -   "refresh_token": "xxx",
  +   "refreshToken": "xxx",
      "token": "xxx"
    }
  }
  ```

---

### 用户管理 (/api/v1/admin/users)

#### v5.1.0 (计划)
- **[Added]** 批量操作接口（待实现）
  - `POST /users/batch/role`
  - `POST /users/batch/status`
  - `POST /users/batch/points`
  - `POST /users/batch/notify`

#### v5.0.0
- **[Changed]** 响应字段统一为 camelCase
  ```json
  {
    "success": true,
    "data": {
      "items": [
        {
  -       "created_at": "2025-01-01T00:00:00Z",
  -       "updated_at": "2025-01-01T00:00:00Z",
  +       "createdAt": "2025-01-01T00:00:00Z",
  +       "updatedAt": "2025-01-01T00:00:00Z",
          "id": 1,
          "name": "张三"
        }
      ]
    }
  }
  ```

---

### 订单管理 (/api/v1/admin/orders)

#### v5.0.0
- **[Changed]** 订单状态枚举值变更
  ```diff
  - 'created' | 'confirmed' | 'in_progress' | 'completed' | 'canceled'
  + 'pending' | 'confirmed' | 'in_progress' | 'completed' | 'canceled' | 'refunded'
  ```

---

### 支付相关 (/api/v1/payments)

#### v3.0.0
- **[Added]** 微信支付接口
  - `POST /payments/wechat/create`
  - `POST /payments/wechat/callback`

- **[Added]** 支付宝支付接口
  - `POST /payments/alipay/create`
  - `POST /payments/alipay/callback`

---

### WebSocket (/ws)

#### v5.0.0
- **[Added]** 系统监控推送
  ```json
  {
    "type": "system_status",
    "timestamp": "2025-01-01T00:00:00Z",
    "data": {
      "cpuUsage": 45.2,
      "memoryUsage": 62.8,
      "status": "healthy"
    }
  }
  ```

- **[Added]** 在线用户推送
  ```json
  {
    "type": "online_users",
    "timestamp": "2025-01-01T00:00:00Z",
    "data": {
      "total": 1234,
      "peak": 1500
    }
  }
  ```

---

## 废弃接口

### 已废弃（仍在使用，即将移除）

| 接口 | 废弃版本 | 移除版本 | 替代方案 |
|------|----------|----------|----------|
| `GET /api/v1/player/me` | v5.0.0 | v6.0.0 | `GET /api/v1/auth/me` |
| `POST /api/v1/orders/:id/finish` | v5.0.0 | v6.0.0 | `PUT /api/v1/orders/:id/complete` |

### 已移除

| 接口 | 移除版本 | 原因 |
|------|----------|------|
| `GET /api/v1/public/players/random` | v4.0.0 | 功能合并到列表接口 |
| `POST /api/v1/auth/login/admin` | v5.0.0 | 统一使用 `/login` + role 参数 |

---

## 迁移指南

### 从 v4.x 升级到 v5.0.0

#### 1. API 响应字段变更

所有接口响应字段从 snake_case 变更为 camelCase。

**迁移步骤**:

```typescript
// 前端迁移示例

// ❌ 旧代码
const users = response.data.items;
users.forEach(user => {
  console.log(user.created_at);
  console.log(user.wallet.balance_cents);
});

// ✅ 新代码
const users = response.data.items;
users.forEach(user => {
  console.log(user.createdAt);
  console.log(user.wallet.balanceCents);
});
```

#### 2. 统一 API 前缀

所有接口使用 `/api/v1` 前缀。

**迁移步骤**:

```typescript
// ❌ 旧代码
const API_BASE = 'http://localhost:8080';
await fetch(`${API_BASE}/auth/login`);

// ✅ 新代码
const API_BASE = 'http://localhost:8080/api/v1';
await fetch(`${API_BASE}/auth/login`);
```

#### 3. 订单状态变更

订单状态枚举值变更。

**迁移步骤**:

```typescript
// ❌ 旧代码
enum OrderStatus {
  Created = 'created',
  Confirmed = 'confirmed',
  // ...
}

// ✅ 新代码
enum OrderStatus {
  Pending = 'pending',
  Confirmed = 'confirmed',
  // ...
}
```

---

## 版本号说明

GameLink 遵循 [语义化版本 2.0.0](https://semver.org/lang/zh-CN/)：

```
主版本号.次版本号.修订号 (Major.Minor.Patch)

Major: 不兼容的 API 变更
Minor: 向下兼容的功能性新增
Patch: 向下兼容的问题修复
```

**当前版本**: 5.1.0

---

## 变更请求流程

### 如何提交 API 变更

1. **创建功能分支**
   ```bash
   git checkout -b feature/api-change
   ```

2. **修改 API 接口**
   - 更新 Handler 代码
   - 更新 Swagger 注解
   - 添加/更新测试用例

3. **更新本文档**
   - 在 `[Unreleased]` 下添加变更记录
   - 如果是 Breaking Change，添加迁移指南

4. **提交 Pull Request**
   - 标题格式：`[api] 变更描述`
   - 指定审查者：Backend-Lead

5. **审查通过后合并**
   - Squash and Merge
   - 更新版本号

---

**文档维护**: 产品经理 + Backend-Lead
**更新频率**: 每次 API 变更后更新

---

## 历史版本发布时间线

```
2025-12-01: v1.0.0 - 基础框架
2025-12-15: v2.0.0 - 陪玩师系统
2026-01-01: v3.0.0 - 支付结算
2026-01-15: v4.0.0 - 架构改进
2026-01-28: v5.0.0 - 功能完善
2026-02-09: v5.1.0 - 工程优化
2026-02-11: [Unreleased] - 待发布
```
