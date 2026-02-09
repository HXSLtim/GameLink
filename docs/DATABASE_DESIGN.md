# GameLink 数据库设计文档

> 版本: 1.0
> 最后更新: 2025-02-09
> 作者: Database-Architect

## 目录

1. [数据库概述](#数据库概述)
2. [设计原则](#设计原则)
3. [表分类与功能](#表分类与功能)
4. [核心表结构](#核心表结构)
5. [索引设计](#索引设计)
6. [外键约束](#外键约束)
7. [软删除策略](#软删除策略)
8. [数据迁移](#数据迁移)
9. [种子数据](#种子数据)
10. [性能优化](#性能优化)

---

## 数据库概述

### 技术选型
- **数据库**: PostgreSQL 16+
- **ORM**: GORM (Go)
- **字符集**: UTF8
- **时区**: UTC

### 数据规模
- **总表数**: 80+ 张表
- **核心业务表**: ~30 张
- **系统管理表**: ~20 张
- **统计监控表**: ~15 张
- **关联表**: ~15 张

### 设计特点
- ✅ 统一的基础模型（Base）
- ✅ 软删除支持（gorm.DeletedAt）
- ✅ 扩展字段支持（ExtJSON JSONB）
- ✅ 完整的时间戳审计（created_at, updated_at）
- ✅ 复合索引优化查询性能

---

## 设计原则

### 1. 命名规范

**表名**: 小写蛇形命名（snake_case）
```sql
chat_group_members
player_certifications
order_service_assignments
```

**列名**: 小写蛇形命名
```go
PlayerID        -> player_id
CreatedAt       -> created_at
OrderNo         -> order_no
```

**JSON标签**: 驼峰命名（camelCase）
```go
PlayerID    `json:"playerId"`
CreatedAt   `json:"createdAt"`
OrderNo     `json:"orderNo"`
```

### 2. 基础模型

所有表继承 `Base` 结构：
```go
type Base struct {
    ID        uint64         `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at;index"`
    UpdatedAt time.Time      `json:"updatedAt" gorm:"column:updated_at"`
    DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"column:deleted_at;index"`
    ExtJSON   string         `json:"extJson,omitempty" gorm:"column:ext_json;type:jsonb;default:'{}'"`
}
```

**字段说明**:
- `id`: 主键（自增 uint64）
- `created_at`: 创建时间（索引）
- `updated_at`: 更新时间（自动维护）
- `deleted_at`: 软删除时间（索引，NULL表示未删除）
- `ext_json`: 扩展字段（JSONB类型，用于灵活扩展）

### 3. 枚举类型

使用字符串常量定义枚举：
```go
type OrderStatus string

const (
    OrderStatusPending    OrderStatus = "pending"
    OrderStatusConfirmed  OrderStatus = "confirmed"
    OrderStatusInProgress OrderStatus = "in_progress"
    OrderStatusCompleted  OrderStatus = "completed"
    OrderStatusCanceled   OrderStatus = "canceled"
    OrderStatusRefunded   OrderStatus = "refunded"
)
```

### 4. 外键约束

```go
// 级联删除（ON DELETE CASCADE）
User     *User `json:"-" gorm:"foreignKey:UserID;references:ID"`

// 设置NULL（ON DELETE SET NULL）
Player   *Player `json:"-" gorm:"foreignKey:PlayerID;references:ID"`

// 限制删除（ON DELETE RESTRICT）
Order    *Order `json:"-" gorm:"foreignKey:OrderID;references:ID"`
```

---

## 表分类与功能

### 1. 用户与权限系统（10张表）

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| `users` | 用户主表 | phone, email, password_hash, role, status |
| `roles` | 角色定义 | slug, name, is_system |
| `permissions` | 权限定义 | method, path, code, group |
| `user_roles` | 用户角色关联 | user_id, role_id |
| `role_permissions` | 角色权限关联 | role_id, permission_id |
| `user_tags` | 用户标签定义 | name, color, sort_order |
| `user_tag_relations` | 用户标签关联 | user_id, tag_id |
| `user_login_history` | 登录历史 | user_id, ip, user_agent, success |
| `user_behavior` | 用户行为追踪 | user_id, action, target_type, target_id |
| `permission_audit_logs` | 权限审计日志 | operator_id, target_type, target_id, action |

**关系图**:
```
User (1) ── (N) UserRole
                 (N) (1) Role
                         (1) (N) RolePermission
                                    (N) (1) Permission
```

### 2. 陪玩师体系（8张表）

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| `players` | 陪玩师档案 | user_id, nickname, bio, rating_average, hourly_rate_cents, verification_status |
| `player_games` | 陪玩师游戏关联 | player_id, game_id, skill_level |
| `player_skill_tags` | 技能标签 | player_id, tag_name, level |
| `game_ranks` | 游戏段位配置 | game_id, rank_name, level, icon_url |
| `player_rank_records` | 段位认证记录 | player_id, game_rank_id, status, verified_at |
| `player_certifications` | 实名认证 | player_id, real_name, id_card, status, verified_at |
| `player_schedules` | 排班表 | player_id, date, start_time, end_time, status |
| `player_presence` | 在线状态 | player_id, status, current_game_id, last_online_at |

**关系图**:
```
User (1) ── (1) Player
                  (1) ── (N) PlayerGame
                              (N) (1) Game
                  (1) ── (N) PlayerSkillTag
                  (1) ── (N) PlayerRankRecord
                                      (N) (1) GameRank
                  (1) ── (1) PlayerCertification
```

### 3. 订单系统（7张表）

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| `orders` | 订单主表 | order_no, user_id, item_id, player_id, total_price_cents, status, currency |
| `order_groups` | 主订单（多时段） | group_no, user_id, total_hours, completed_hours, status |
| `order_items` | 订单明细（槽位） | order_id, item_id, slot, unit_price_cents, player_id, status |
| `order_players` | 订单陪玩师关联 | order_id, order_item_id, player_id, status, joined_at |
| `payments` | 支付记录 | order_id, user_id, method, amount_cents, status, provider_trade_no |
| `order_disputes` | 争议处理 | order_id, reporter_id, reason, status, result |
| `order_timeout_logs` | 超时日志 | order_id, timeout_type, timeout_at, action, remark |

**关系图**:
```
User (1) ── (N) Order (1) ── (N) Payment
                  (1) ── (1) ServiceItem
                  (1) ── (N) OrderItem (1) ── (1) OrderPlayer (1) ── (1) Player
                  (N) ── (1) OrderGroup
                  (1) ── (N) OrderDispute
```

### 4. 游戏与服务（6张表）

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| `games` | 游戏配置 | key, name, category_id, icon_url, is_active, sort_order |
| `game_categories` | 游戏分类 | name, description, icon_url, sort_order |
| `service_items` | 服务项目 | item_code, name, sub_category, game_id, base_price_cents, commission_rate |
| `coupons` | 用户优惠券 | template_id, user_id, code, discount_type, value, status |
| `coupon_templates` | 优惠券模板 | name, discount_type, value, min_amount, max_discount, total_count |
| `recharge_options` | 充值档位 | amount_cents, bonus_cents, sort_order |
| `recharge_records` | 充值记录 | user_id, option_id, amount_cents, bonus_cents, payment_method |

**关系图**:
```
GameCategory (1) ── (N) Game (1) ── (N) ServiceItem
CouponTemplate (1) ── (N) Coupon (N) ── (1) User
RechargeOption (1) ── (N) RechargeRecord (N) ── (1) User
```

### 5. 评价与社交（10张表）

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| `reviews` | 评价主表 | order_id, user_id, player_id, score, content, status |
| `review_reports` | 评价举报 | review_id, reporter_id, reason, status, result |
| `review_replies` | 评价回复 | review_id, player_id, content, reply_count, status |
| `review_display_settings` | 评价展示设置 | player_id, hide_low_rating, display_condition |
| `favorites` | 收藏 | user_id, target_type, target_id (player/game/service) |
| `user_blocks` | 用户拉黑 | blocker_id, blocked_id, reason, canceled_at |
| `chat_groups` | 聊天群组 | group_name, group_type, created_by, max_members, is_active |
| `chat_group_members` | 群成员 | group_id, user_id, role, nickname, is_active |
| `chat_messages` | 聊天消息 | group_id, sender_id, content, message_type, audit_status |
| `chat_reports` | 聊天举报 | message_id, reporter_id, reason, status, result |

**关系图**:
```
User (1) ── (N) Review (1) ── (1) Order
                (1) ── (1) Player
                (1) ── (N) ReviewReply
                (1) ── (N) ReviewReport

User (1) ── (N) Favorite
User (1) ── (N) UserBlock

ChatGroup (1) ── (N) ChatGroupMember (N) ── (1) User
ChatGroup (1) ── (N) ChatMessage (N) ── (1) User
```

### 6. 内容社区（4张表）

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| `feeds` | 动态/帖子 | author_id, content, category_id, visibility, moderation_status |
| `feed_images` | 动态图片 | feed_id, url, order, width, height, size_bytes |
| `feed_reports` | 动态举报 | feed_id, reporter_id, reason, status, result |
| `content_categories` | 内容分类 | name, description, icon_url, sort_order, status |

### 7. 财务系统（8张表）

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| `wallets` | 用户钱包 | user_id, balance_cents, frozen_cents, total_recharge_cents |
| `withdraws` | 提现记录 | user_id, amount_cents, status, bank_name, account_number, rejected_at |
| `commission_rules` | 抽成规则 | name, type, rate, game_id, player_id, service_type, is_active |
| `commission_records` | 抽成记录 | order_id, player_id, commission_cents, settlement_month, status |
| `monthly_settlements` | 月度结算 | player_id, settlement_month, total_orders, total_income, status |
| `collection_entities` | 收款主体 | name, type, bank_name, account_number, is_active |
| `settlement_companies` | 结算公司 | name, bank_name, account_number, commission_rate |
| `refund_records` | 退款记录 | payment_id, order_id, user_id, amount_cents, reason, status |

### 8. 团队与活动（8张表）

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| `teams` | 团队 | leader_id, name, description, max_members, current_order_id |
| `team_members` | 团队成员 | team_id, player_id, role, status, joined_at |
| `team_invites` | 团队邀请 | team_id, player_id, inviter_id, status, expires_at |
| `activities` | 活动 | title, type, start_time, end_time, status, rules |
| `activity_rewards` | 活动奖励 | activity_id, reward_type, value, stock |
| `activity_participations` | 活动参与记录 | activity_id, user_id, status, progress, coupon_ids |
| `referrals` | 推荐记录 | referrer_id, referee_id, code_id, status |
| `referral_codes` | 邀请码 | user_id, code, max_uses, used_count, expires_at |
| `referral_rewards` | 推荐奖励 | referral_id, user_id, reward_type, value, status |
| `vip_levels` | VIP等级 | level, name, required_exp, benefits |
| `vip_configs` | VIP系统配置 | key, value, description |

### 9. 通知与监控（9张表）

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| `notification_templates` | 通知模板 | code, title, content, type, channel |
| `user_notifications` | 用户通知 | user_id, template_id, title, content, is_read, related_type, related_id |
| `user_notification_settings` | 通知设置 | user_id, notification_type, enabled, channels |
| `notification_configs` | 通知系统配置 | key, value, description |
| `notification_schedules` | 定时通知任务 | template_code, cron_expression, is_active |
| `alerts` | 监控告警 | level, type, title, message, source, is_read |
| `kpi_targets` | KPI目标 | metric_name, target_value, period, status |
| `operation_logs` | 操作日志 | actor_user_id, entity_type, entity_id, action, changes |
| `sensitive_words` | 敏感词 | word, type, replacement, status |
| `order_timeout_configs` | 订单超时配置 | config_key, config_value, description |
| `order_service_assignments` | 订单客服分配 | order_id, service_user_id, chat_group_id, status |

### 10. 统计与LFG（7张表）

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| `user_statistics` | 用户统计 | user_id, date, orders_count, spent_cents |
| `player_statistics` | 陪玩师统计 | player_id, date, orders_count, earnings_cents |
| `service_item_statistics` | 服务项统计 | item_id, date, orders_count, revenue_cents |
| `game_statistics` | 游戏统计 | game_id, date, orders_count, revenue_cents |
| `platform_statistics` | 平台统计 | date, new_users, total_orders, revenue_cents |
| `lfg_requests` | 快速组队匹配 | user_id, game_id, required_players, title, status, expires_at |
| `tag_thresholds` | 标签阈值配置 | tag_type, tag_name, min_value, benefits |

---

## 核心表结构

### 用户表 (users)

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    phone VARCHAR(32) UNIQUE,
    email VARCHAR(128) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(64),
    nickname VARCHAR(64),
    avatar_url VARCHAR(255),
    role VARCHAR(32) NOT NULL DEFAULT 'user',
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    last_login_at TIMESTAMP,
    ban_reason VARCHAR(500),
    banned_at TIMESTAMP,
    banned_by BIGINT,

    -- VIP相关
    vip_level_id BIGINT,
    vip_unlocked BOOLEAN DEFAULT FALSE,
    vip_exp BIGINT DEFAULT 0,
    total_recharge_cents BIGINT DEFAULT 0,
    vip_unlocked_at TIMESTAMP,
    vip_expire_at TIMESTAMP,
    last_monthly_coupon_at TIMESTAMP,

    -- 微信
    wechat_open_id VARCHAR(64) UNIQUE,
    wechat_union_id VARCHAR(128),

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- 索引
CREATE INDEX idx_users_status_last_login ON users (status, last_login_at);
CREATE INDEX idx_users_email ON users (email) WHERE email IS NOT NULL;
CREATE INDEX idx_users_phone ON users (phone) WHERE phone IS NOT NULL;
```

### 陪玩师表 (players)

```sql
CREATE TABLE players (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    nickname VARCHAR(64),
    bio TEXT,
    rank VARCHAR(32),
    rating_average FLOAT DEFAULT 0 CHECK (rating_average >= 0 AND rating_average <= 5),
    rating_count INT DEFAULT 0,
    order_count INT DEFAULT 0,
    hourly_rate_cents BIGINT,
    main_game_id BIGINT,
    verification_status VARCHAR(32) DEFAULT 'pending',

    -- 在线状态
    online_status VARCHAR(32) DEFAULT 'offline',
    accepting_orders BOOLEAN DEFAULT FALSE,
    last_online_at TIMESTAMP,

    -- 审核相关
    verified_at TIMESTAMP,
    verified_by BIGINT,
    verify_remark VARCHAR(500),
    reject_reason VARCHAR(500),

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- 索引
CREATE INDEX idx_players_user_id ON players (user_id);
CREATE INDEX idx_players_verification ON players (verification_status);
CREATE INDEX idx_players_online_status ON players (online_status, accepting_orders);
```

### 订单表 (orders)

```sql
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    order_no VARCHAR(64) UNIQUE,
    user_id BIGINT NOT NULL,
    item_id BIGINT NOT NULL,
    player_id BIGINT,
    recipient_player_id BIGINT,

    -- 价格
    quantity INT DEFAULT 1,
    unit_price_cents BIGINT NOT NULL,
    total_price_cents BIGINT NOT NULL,
    commission_cents BIGINT DEFAULT 0,
    player_income_cents BIGINT DEFAULT 0,
    currency CHAR(3) DEFAULT 'CNY',

    -- 订单信息
    status VARCHAR(32) DEFAULT 'pending',
    title VARCHAR(128),
    description TEXT,

    -- 护航服务字段
    game_id BIGINT,
    scheduled_start TIMESTAMP,
    scheduled_end TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,

    -- 礼物订单字段
    gift_message TEXT,
    is_anonymous BOOLEAN DEFAULT FALSE,
    delivered_at TIMESTAMP,

    -- 取消/退款
    cancel_reason TEXT,
    refund_amount_cents BIGINT DEFAULT 0,
    refund_reason TEXT,
    refunded_at TIMESTAMP,

    -- 订单拆分
    group_id BIGINT,
    hour_index INT DEFAULT 1,
    is_sub_order BOOLEAN DEFAULT FALSE,
    can_transfer BOOLEAN DEFAULT TRUE,
    transfer_from BIGINT,
    transfer_to BIGINT,
    transfer_note VARCHAR(500),

    -- 多人服务
    required_players INT DEFAULT 1,
    current_players INT DEFAULT 0,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- 索引（Covering Index优化）
CREATE INDEX idx_orders_user_status_created ON orders (user_id, status, created_at DESC)
    INCLUDE (id, player_id, total_price_cents, commission_cents, player_income_cents);
CREATE INDEX idx_orders_player_status ON orders (player_id, status);
CREATE INDEX idx_orders_game_created ON orders (game_id, created_at DESC);
CREATE INDEX idx_orders_group_hour ON orders (group_id, hour_index);
CREATE INDEX idx_orders_transfer ON orders (transfer_from, transfer_to);
```

### 聊天群组表 (chat_groups)

```sql
CREATE TABLE chat_groups (
    id BIGSERIAL PRIMARY KEY,
    group_name VARCHAR(128) NOT NULL,
    group_type VARCHAR(32) NOT NULL, -- public, private, order, team, lfg, custom
    related_order_id BIGINT,
    created_by BIGINT NOT NULL,
    max_members INT DEFAULT 100,
    is_active BOOLEAN DEFAULT TRUE,

    -- 游戏房间字段
    game_id BIGINT,
    room_status VARCHAR(32) DEFAULT 'waiting', -- waiting, ready, in_game, finished, canceled
    is_private BOOLEAN DEFAULT FALSE,
    password VARCHAR(64),
    current_members INT DEFAULT 0,
    related_team_id BIGINT,
    related_lfg_id BIGINT,
    started_at TIMESTAMP,
    finished_at TIMESTAMP,

    -- 语音服务
    voice_enabled BOOLEAN DEFAULT FALSE,
    voice_room_id VARCHAR(128),
    voice_provider VARCHAR(32),
    voice_sdk_app_id BIGINT,
    voice_started_at TIMESTAMP,
    voice_ended_at TIMESTAMP,
    voice_duration INT DEFAULT 0,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- 索引
CREATE INDEX idx_chat_groups_type ON chat_groups (group_type);
CREATE INDEX idx_chat_groups_order ON chat_groups (related_order_id) WHERE related_order_id IS NOT NULL;
CREATE INDEX idx_chat_groups_status ON chat_groups (room_status);
```

---

## 索引设计

### 索引策略

**1. 主键索引**
```sql
-- 所有表默认有主键索引
id BIGSERIAL PRIMARY KEY
```

**2. 唯一索引**
```sql
-- 用户表
CREATE UNIQUE INDEX idx_users_email ON users (email) WHERE email IS NOT NULL;
CREATE UNIQUE INDEX idx_users_phone ON users (phone) WHERE phone IS NOT NULL;

-- 订单表
CREATE UNIQUE INDEX idx_orders_no ON orders (order_no) WHERE order_no != '';
```

**3. 普通索引**
```sql
-- 陪玩师在线状态
CREATE INDEX idx_players_online_status ON players (online_status, accepting_orders);

-- 订单状态
CREATE INDEX idx_orders_status ON orders (status);
```

**4. 复合索引**
```sql
-- 用户订单列表（Covering Index）
CREATE INDEX idx_orders_user_status_created ON orders (user_id, status, created_at DESC)
    INCLUDE (id, player_id, total_price_cents);

-- 陪玩师订单
CREATE INDEX idx_orders_player_status ON orders (player_id, status);

-- 聊天群组成员
CREATE UNIQUE INDEX idx_chat_group_member ON chat_group_members (group_id, user_id);
```

**5. 部分索引**
```sql
-- 只索引非NULL值
CREATE INDEX idx_users_email ON users (email) WHERE email IS NOT NULL;
```

### 索引优化原则

1. **WHERE子句优先**: 为常用查询条件的字段创建索引
2. **ORDER BY支持**: 复合索引包含排序字段
3. **覆盖索引**: INCLUDE常用字段，避免回表查询
4. **选择性高的字段**: 优先为高选择性字段（如status）创建索引

---

## 外键约束

### 级联规则

**1. CASCADE（级联删除）**
```sql
-- 删除用户时，删除其登录历史
ALTER TABLE user_login_history
ADD CONSTRAINT fk_user_login_history_user
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
```

**2. SET NULL（设为NULL）**
```sql
-- 删除陪玩师时，订单的player_id设为NULL
ALTER TABLE orders
ADD CONSTRAINT fk_order_player
FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE SET NULL;
```

**3. RESTRICT（限制删除）**
```sql
-- 有订单时禁止删除用户
ALTER TABLE orders
ADD CONSTRAINT fk_order_user
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT;
```

### 外键命名规范

```
fk_{table}_{referenced_table}
```

例如：
```sql
fk_order_user        -- orders表引用users表
fk_order_player      -- orders表引用players表
fk_chat_member_group  -- chat_group_members表引用chat_groups表
```

---

## 软删除策略

### 实现方式

使用GORM的软删除插件：
```go
import "gorm.io/plugin/soft_delete"

db.Use(soft_delete.New{})
```

**Base模型**:
```go
type Base struct {
    DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}
```

### 查询处理

**自动过滤已删除记录**:
```go
// GORM自动添加 WHERE deleted_at IS NULL
db.Find(&users)
// 实际SQL: SELECT * FROM users WHERE deleted_at IS NULL
```

**包含已删除记录**:
```go
db.Unscoped().Find(&users)
// 实际SQL: SELECT * FROM users
```

**永久删除**:
```go
db.Unscoped().Delete(&user)
// 实际SQL: DELETE FROM users WHERE id = ?
```

### 软删除索引

```sql
CREATE INDEX idx_deleted_at ON users (deleted_at);
```

PostgreSQL会自动创建部分索引，只索引非NULL值。

---

## 数据迁移

### 迁移文件

**位置**: `api/pkg/db/migrate.go`

**迁移流程**:
1. **Phase 1**: 创建基础表（Game, User, Player, Order, Payment）
2. **Phase 2**: 创建依赖表（包含外键的表）
3. **标记版本**: 记录迁移版本到 `seed_metadata` 表

**版本控制**:
```go
const migrateVersion = "2026-02-07-v1"

func isMigrateUpToDate(db *gorm.DB) bool {
    var val string
    db.Raw(`SELECT value FROM seed_metadata WHERE key = 'migrate_version'`).Scan(&val)
    return val == migrateVersion
}
```

**自动迁移**:
```go
func autoMigrate(db *gorm.DB) error {
    // 检查版本
    if isMigrateUpToDate(db) {
        return nil
    }

    // Phase 1: 基础表
    db.AutoMigrate(&model.Game{}, &model.User{}, &model.Player{}, &model.Order{}, &model.Payment{})

    // Phase 2: 依赖表
    db.AutoMigrate(&model.ChatGroup{}, &model.ChatMessage{}, &model.Review{}, /* ... */)

    // 标记版本
    markMigrateVersion(db)
    return nil
}
```

### 数据修复

**数据一致性修复** (`runDataFixups`):
```go
func runDataFixups(db *gorm.DB) error {
    // 标准化订单状态拼写
    db.Exec("UPDATE orders SET status='canceled' WHERE status='cancelled'")

    // 生成订单号
    generateOrderNumbers(db)

    // 确保默认角色
    ensureDefaultRoles(db)

    // 确保超级管理员
    ensureSuperAdmin(db)

    return nil
}
```

### 索引创建

**位置**: `api/pkg/db/migrate.go:ensureIndexes`

**关键索引**:
```sql
-- 订单查询优化
CREATE INDEX idx_orders_user_status_created ON orders (user_id, status, created_at DESC);
CREATE INDEX idx_orders_player_status ON orders (player_id, status);

-- 支付查询优化
CREATE INDEX idx_payments_status_created ON payments (status, created_at DESC);
CREATE INDEX idx_payments_order_created ON payments (order_id, created_at DESC);

-- 佣金查询优化
CREATE INDEX idx_commission_records_player_month ON commission_records (player_id, settlement_month);
```

---

## 种子数据

### 种子数据版本

**当前版本**: `2026-02-07-v4`

**版本检查**:
```go
const seedVersion = "2026-02-07-v4"

func isSeedUpToDate(db *gorm.DB) bool {
    var val string
    db.Raw(`SELECT value FROM seed_metadata WHERE key = 'seed_version'`).Scan(&val)
    return val == seedVersion
}
```

### 种子数据类型

**1. 用户数据** (`seedUsers`)
- 管理员用户（adminA, adminB）
- 陪玩师用户（playerA, playerB, playerC）
- 普通用户（customerA-H）

**2. 陪玩师数据** (`seedPlayers`)
- 昵称、简介、评分、时薪
- 主游戏、认证状态

**3. 游戏数据** (`seedGames`)
- 英雄联盟 (lol)
- 无畏契约 (valorant)
- DOTAP 2 (dota2)

**4. 服务项数据** (`seedServiceItems`)
- 单人护航 (solo)
- 团队护航 (team)
- 礼物 (gift)

**5. 订单数据** (`seedOrders`)
- 各种状态的订单
- 不同游戏、不同陪玩师

**6. 支付数据** (`seedPayments`)
- 支付宝、微信、钱包
- 已支付、待支付、已退款

**7. 评价数据** (`seedReviews`)
- 不同评分、不同状态
- 带图片、被举报

**8. 聊天数据** (`seedChatGroups`, `seedChatMessages`)
- 7种群组类型
- 各种消息类型和审核状态

**9. 内容数据** (`seedContentData`)
- 内容分类
- 动态帖子
- 聊天消息
- 举报数据

**10. 流程数据** (`seedAdditionalFlowOrders`)
- 礼物订单流程
- 团队订单流程
- 支付失败流程

### 种子数据验证

**位置**: `api/pkg/db/seed_flow.go:validateSeedAssociations`

**验证规则** (14+条):
```go
// 1. payment.user_id == order.user_id
// 2. review.user_id == order.user_id
// 3. review.player_id == order.player_id (当存在时)
// 4. gift订单必须使用gift service_item
// 5. team订单必须有order_items/order_players
// 6. activity participation的coupon_ids必须存在
// 7-14. 其他外键约束验证
```

---

## 性能优化

### 1. 查询优化

**避免N+1查询**:
```go
// ❌ 不好：N+1查询
orders := []Order{}
db.Find(&orders)
for _, order := range orders {
    db.Preload("Player").First(&order, order.ID)  // N次查询
}

// ✅ 好：使用Joins或Preload
db.Joins("Player").Find(&orders)
// 或
db.Preload("Player").Find(&orders)
```

**覆盖索引**:
```sql
-- 不需要回表查询
CREATE INDEX idx_orders_user_status_created_covering
ON orders (user_id, status, created_at DESC)
INCLUDE (id, player_id, total_price_cents);

-- 查询时直接从索引获取所有数据
EXPLAIN (ANALYZE)
SELECT id, player_id, total_price_cents
FROM orders
WHERE user_id = ? AND status = 'pending'
ORDER BY created_at DESC;
```

### 2. 批量操作

**批量插入**:
```go
// ✅ 使用CreateInBatches
var users []User
db.CreateInBatches(users, 100)  // 每批100条
```

**批量更新**:
```go
// ✅ 使用Clauses
db.Model(&User{}).Clauses(clause.OnConflict{
    UpdateAll: true,
}).Creates(users)
```

### 3. 连接池配置

**位置**: `api/pkg/db/postgres.go`

```go
func NewPostgresDB(config *config.Database) (*gorm.DB, error) {
    // 连接池配置
    db.DB().SetMaxIdleConns(10)
    db.DB().SetMaxOpenConns(100)
    db.DB().SetConnMaxLifetime(time.Hour)
}
```

### 4. 慢查询优化

**分析慢查询**:
```sql
-- 启用查询日志
SET log_min_duration_statement = 100;  -- 记录超过100ms的查询

-- 查看执行计划
EXPLAIN ANALYZE SELECT * FROM orders WHERE user_id = ?;
```

**优化建议**:
- 为常用查询条件添加索引
- 使用EXPLAIN分析执行计划
- 避免SELECT *，只查询需要的字段
- 合理使用分页（LIMIT + OFFSET）

### 5. 数据归档策略

**聊天消息归档** (30天保留期):
```sql
-- 删除30天前的消息
DELETE FROM chat_messages
WHERE created_at < NOW() - INTERVAL '30 days';
```

**订单数据归档** (1年前):
```sql
-- 归档1年前已完成的订单
DELETE FROM orders
WHERE status = 'completed'
  AND completed_at < NOW() - INTERVAL '1 year';
```

**统计数据聚合**:
```sql
-- 按天聚合统计数据
INSERT INTO platform_statistics (date, new_users, total_orders, revenue_cents)
SELECT
    DATE(created_at) as date,
    COUNT(DISTINCT user_id) as new_users,
    COUNT(*) as total_orders,
    SUM(total_price_cents) as revenue_cents
FROM orders
WHERE created_at >= CURRENT_DATE
GROUP BY DATE(created_at);
```

---

## 附录

### A. 数据库字段类型映射

| Go类型 | PostgreSQL类型 | 说明 |
|--------|---------------|------|
| `int64` | `BIGINT` | 8字节整数 |
| `uint64` | `BIGSERIAL` | 自增主键 |
| `float32` | `REAL` | 4字节浮点数 |
| `float64` | `DOUBLE PRECISION` | 8字节浮点数 |
| `string` | `VARCHAR(n)` | 变长字符串 |
| `string` | `TEXT` | 长文本 |
| `time.Time` | `TIMESTAMP` | 时间戳 |
| `bool` | `BOOLEAN` | 布尔值 |
| `[]byte` | `BYTEA` | 二进制数据 |
| `json.RawMessage` | `JSONB` | JSON数据 |

### B. 数据库大小估算

**预估数据量**（运行1年）:
- 用户: 10,000条
- 陪玩师: 500条
- 订单: 50,000条
- 聊天消息: 1,000,000条
- 评价: 30,000条

**存储空间估算**:
- 用户表: ~5 MB
- 订单表: ~50 MB
- 聊天消息: ~500 MB
- 总计: ~1 GB（不含索引）

### C. 备份策略

**全量备份**: 每天凌晨2点
```bash
pg_dump -U gamelink -d gamelink > backup_$(date +%Y%m%d).sql
```

**增量备份**: 每小时
```bash
pg_dump -U gamelink -d gamelink --schema-only > schema.sql
```

**WAL归档**: 实时归档
```ini
# postgresql.conf
archive_mode = on
archive_command = 'cp %p /wal_archive/%f'
```

### D. 监控指标

**关键监控指标**:
1. 数据库连接数
2. 慢查询数量
3. 表大小增长
4. 索引使用率
5. 死锁发生次数

**监控SQL**:
```sql
-- 表大小
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- 索引使用率
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan as index_scans,
    idx_tup_read as tuples_read,
    idx_tup_fetch as tuples_fetched
FROM pg_stat_user_indexes
ORDER BY idx_scan ASC;
```

---

**文档维护**: Database-Architect
**最后更新**: 2025-02-09
**版本**: 1.0
