# 小程序 API 实现 TODO

> 📅 创建时间: 2026-01-10
> 🎯 目标: 实现小程序端所需的全部 API

---

## 📊 实现进度总览

| 优先级 | 模块 | 状态 | 预计工时 |
|--------|------|------|----------|
| P0 | 微信小程序登录 | ✅ 已完成 | 4h |
| P0 | 角色切换机制 | ✅ 已完成 | 2h |
| P1 | 公共 API 层 | ✅ 已完成 | 3h |
| P1 | 在线状态管理 | ✅ 已完成 | 2h |
| P1 | 接单大厅 | ✅ 已有实现 | 0h |
| P2 | 收藏功能 | ✅ 已完成 | 2h |
| P2 | 搜索功能 | ✅ 已完成 | 2h |

**已完成**: ~15h | **剩余**: 0h

---

## P0: 微信小程序登录

### 文件清单
- [ ] `api/internal/handler/public/auth.go` - 公共认证 Handler
- [ ] `api/internal/service/auth/wechat.go` - 微信登录服务
- [ ] `api/internal/model/wechatSession.go` - 微信会话模型（可选）

### API 端点
```
POST /api/v1/public/auth/wechat/login
POST /api/v1/public/auth/phone/login
POST /api/v1/public/auth/refresh
```

### 实现要点
1. 调用微信 `code2Session` 接口获取 `openid` 和 `session_key`
2. 解密 `encryptedData` 获取手机号（可选）
3. 查找或创建用户（通过 `openid` 关联）
4. 处理推荐码逻辑
5. 生成 JWT Token（包含 `is_player` 和 `current_role` 声明）

### 依赖配置
```env
WECHAT_APPID=xxx
WECHAT_SECRET=xxx
```

---

## P0: 角色切换机制

### 文件清单
- [ ] `api/internal/handler/user/role.go` - 角色切换 Handler
- [ ] `api/internal/service/auth/role.go` - 角色切换服务

### API 端点
```
GET  /api/v1/user/roles        - 获取可用角色
POST /api/v1/user/switch-role  - 切换角色
```

### 实现要点
1. 验证用户是否有目标角色权限（`is_player` 字段）
2. 生成新的 JWT Token（更新 `current_role` 声明）
3. 返回新 Token 供前端更新

### JWT Payload 结构
```json
{
  "user_id": 123,
  "role": "user",           // 基础角色
  "current_role": "player", // 当前激活角色
  "is_player": true,
  "exp": 1735689600
}
```

---

## P1: 公共 API 层

### 文件清单
- [ ] `api/internal/handler/public/player.go` - 公开陪玩师列表
- [ ] `api/internal/handler/public/game.go` - 游戏列表
- [ ] `api/internal/handler/public/serviceItem.go` - 服务项目列表

### API 端点
```
GET /api/v1/public/players          - 陪玩师列表（无需认证）
GET /api/v1/public/players/:id      - 陪玩师详情（无需认证）
GET /api/v1/public/games            - 游戏列表
GET /api/v1/public/service-items    - 服务项目列表
```

### 实现要点
1. 复用现有 `user/player.go` 的逻辑
2. 移除认证要求
3. 添加缓存支持（热门数据）

---

## P1: 在线状态管理

### 文件清单
- [ ] `api/internal/handler/player/status.go` - 在线状态 Handler
- [ ] `api/internal/service/player/status.go` - 在线状态服务
- [ ] `api/internal/model/playerStatus.go` - 状态模型（如需持久化）

### API 端点
```
GET  /api/v1/player/online-status   - 获取在线状态
PUT  /api/v1/player/online-status   - 更新在线状态
```

### 实现要点
1. 状态存储在 Redis（实时性要求高）
2. 状态类型: `online` | `busy` | `offline`
3. 自动下线机制（心跳超时）
4. WebSocket 广播状态变更

### Redis Key 设计
```
player:status:{player_id} -> {status, updated_at, current_orders}
player:online:set -> SET of online player IDs
```

---

## P1: 接单大厅

### 文件清单
- [ ] `api/internal/handler/player/availableOrders.go` - 可接订单 Handler

### API 端点
```
GET  /api/v1/player/available-orders     - 获取可接订单列表
POST /api/v1/player/orders/:id/accept    - 接单
POST /api/v1/player/orders/:id/reject    - 拒单
```

### 实现要点
1. 查询状态为 `waiting` 的订单
2. 按游戏、服务类型筛选
3. 排序：VIP用户优先、等待时间
4. 接单时检查并发限制

---

## P2: 收藏功能

### 文件清单
- [ ] `api/internal/handler/user/favorite.go` - 收藏 Handler
- [ ] `api/internal/model/favorite.go` - 收藏模型
- [ ] `api/internal/repository/favorite.go` - 收藏仓储

### API 端点
```
GET    /api/v1/user/favorites/players      - 获取收藏列表
POST   /api/v1/user/favorites/players/:id  - 添加收藏
DELETE /api/v1/user/favorites/players/:id  - 取消收藏
```

### 数据模型
```go
type Favorite struct {
    ID        uint64
    UserID    uint64
    PlayerID  uint64
    CreatedAt time.Time
}
```

---

## P2: 搜索功能

### 文件清单
- [ ] `api/internal/handler/public/search.go` - 搜索 Handler
- [ ] `api/internal/service/search/search.go` - 搜索服务

### API 端点
```
GET /api/v1/public/search?q=关键词&type=player|game
```

### 实现要点
1. 支持陪玩师昵称、游戏名称搜索
2. 模糊匹配（LIKE 或全文索引）
3. 结果分组返回

---

## 🔧 路由注册调整

### 新增路由组
```go
// api/internal/router/router.go

// 公共 API（无需认证）
public := v1.Group("/public")
{
    publicHandler.RegisterAuthRoutes(public, authSvc)
    publicHandler.RegisterPlayerRoutes(public, playerSvc)
    publicHandler.RegisterGameRoutes(public, gameSvc)
    publicHandler.RegisterSearchRoutes(public, searchSvc)
}

// 用户端 API
user := v1.Group("/user")
user.Use(authMiddleware)
{
    // ... 现有路由
    userHandler.RegisterRoleRoutes(user, authSvc)
    userHandler.RegisterFavoriteRoutes(user, favoriteSvc)
}

// 陪玩师端 API
player := v1.Group("/player")
player.Use(authMiddleware, playerRoleMiddleware)
{
    // ... 现有路由
    playerHandler.RegisterStatusRoutes(player, statusSvc)
    playerHandler.RegisterAvailableOrderRoutes(player, orderSvc)
}
```

---

## 📝 实现顺序

### Phase 1: 认证基础 (Day 1)
1. [x] 创建 TODO 文档
2. [x] 实现微信登录服务 (`api/internal/service/auth/wechat.go`)
3. [x] 实现公共认证 Handler (`api/internal/handler/public/auth.go`)
4. [x] 实现角色切换 (`api/internal/service/auth/role.go`, `api/internal/handler/user/role.go`)

### Phase 2: 公共 API (Day 1)
5. [x] 创建 public handler 目录
6. [x] 实现公开陪玩师列表 (`api/internal/handler/public/player.go`)
7. [x] 实现游戏/服务项目列表 (`api/internal/handler/public/game.go`, `serviceItem.go`)
8. [x] 注册公共路由 (`api/internal/router/publicRoutes.go`)

### Phase 3: 陪玩师功能 (Day 3)
9. [x] 实现在线状态管理
10. [x] 实现接单大厅（已有实现：`/player/orders/available`）
11. [ ] WebSocket 状态广播（可选优化）

### Phase 4: 辅助功能 (Day 4)
12. [x] 实现收藏功能
13. [x] 实现搜索功能
14. [ ] 集成测试

---

## ✅ 验收标准

- [ ] 所有 API 端点可正常访问
- [ ] JWT Token 包含正确的角色声明
- [ ] 角色切换后 Token 正确更新
- [ ] 公共 API 无需认证即可访问
- [ ] 在线状态实时更新
- [ ] 单元测试覆盖率 > 70%
