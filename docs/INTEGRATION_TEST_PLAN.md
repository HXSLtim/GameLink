# GameLink 集成测试规划

> **文档版本**: v2.0
> **创建日期**: 2025-12-28
> **更新日期**: 2025-12-29
> **状态**: 基于PRD全面规划

---

## 一、概述

### 1.1 测试目标

集成测试旨在验证服务层（Service Layer）各模块之间的交互是否正确，确保业务逻辑在真实数据库环境下的正确性。

### 1.2 测试范围

- **测试层级**: Service 层集成测试
- **测试环境**: PostgreSQL 测试数据库
- **测试方式**: 真实数据库 + 真实 Repository + 真实 Service

### 1.3 当前状态

| 模块分类 | 总数 | 已有集成测试 | 缺失 | 覆盖率 |
|---------|------|-------------|------|--------|
| 核心业务模块 | 19 | 10 | 9 | 53% |
| 营销模块 | 6 | 6 | 0 | 100% |
| 辅助模块 | 7 | 1 | 6 | 14% |
| **总计** | **32** | **17** | **15** | **53%** |

### 1.4 已完成的集成测试

✅ **Phase 1 P0 核心业务流程测试** (已完成编写，待运行验证):
- `auth` - 认证服务 (12个测试用例)
- `player` - 陪玩师服务 (9个测试用例)
- `payment` - 支付服务 (13个测试用例)
- `withdraw` - 提现服务 (12个测试用例)
- `order` - 订单服务 (8个测试用例)
- `dispute` - 争议服务 (9个测试用例)
- `chat` - 聊天服务 (17个测试用例)
- `notification` - 通知服务 (4个测试用例)
- `menu` - 菜单服务 (2个测试用例)
- `permission` - 权限服务 (4个测试用例)

✅ **营销模块测试**:
- `vip` - VIP会员服务 (测试文件已创建)
- `coupon` - 优惠券服务 (测试文件已创建)
- `recharge` - 充值服务 (测试文件已创建)
- `activity` - 活动服务 (测试文件已创建)
- `team` - 团队服务 (测试文件已创建)
- `referral` - 推荐服务 (测试文件已创建)

**总计**: 约90+个测试用例已完成框架搭建

⏳ **待实现**:
- review - 评价服务 (需要重构以匹配新API)
- commission - 佣金服务 (需要重构以匹配新API)
- userblock - 用户拉黑服务 (需要新建)
- user - 用户服务 (需要新建)
- kpi - KPI服务
- statistics - 统计服务
- wallet - 钱包服务

---

## 二、基于PRD的完整测试场景规划

### 2.1 用户端功能测试场景

根据PRD §4.1，用户端包含以下功能模块：

#### 2.1.1 首页浏览

| 测试场景 | 描述 | 优先级 |
|---------|------|--------|
| 获取游戏列表 | 返回所有活跃游戏，按分类排序 | P0 |
| 获取陪玩师列表 | 支持分页、按游戏/段位筛选 | P0 |
| 搜索陪玩师 | 按昵称/ID搜索 | P1 |
| 获取排行榜 | 按评分/接单量/收入排名 | P2 |

**集成测试位置**: `game_integration_test.go`, `player_integration_test.go`, `ranking_integration_test.go`

#### 2.1.2 订单功能

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 创建solo订单 | 单人陪玩订单 | P0 | 需指定游戏、服务项、陪玩师 |
| 创建team订单 | 多人陪玩订单 | P0 | 需指定座位数量 |
| 创建gift订单 | 礼物订单（直接付款） | P0 | 无需服务完成，立即完成 |
| 订单支付 | 第三方支付 | P0 | 支付后订单状态变更为confirmed |
| 钱包支付 | 余额支付 | P0 | 余额不足时需组合支付 |
| 组合支付 | 钱包+第三方 | P1 | 先扣钱包，剩余第三方支付 |
| 我的订单列表 | 分页查询我的订单 | P0 | 支持状态筛选 |
| 订单详情 | 查看订单完整信息 | P0 | 包含陪玩师、支付、评价、时间线 |
| 取消订单 | 用户取消pending订单 | P0 | 取消后订单状态为canceled |
| 申请退款 | 订单完成后7天内可申请 | P0 | 超过7天不可申请 |
| 订单补差价 | 升级服务或加座位 | P1 | 生成补差订单 |
| 订单超时处理 | 支付超时/接单超时自动取消 | P1 | 定时任务处理 |

**集成测试位置**: `order_integration_test.go`

#### 2.1.3 评价功能

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 提交评价 | 订单完成后评价 | P0 | 只有已完成订单可评价 |
| 一键评价 | 快速5星好评 | P1 | 系统预设好评内容 |
| 单独评价 | 分别评分（态度/技术/沟通） | P1 | 多维度评分 |
| 修改评价 | 评价后可修改 | P1 | 仅限评价者本人 |
| 评价回复 | 陪玩师回复评价 | P2 | |
| 评价申诉 | 对不当评价申诉 | P2 | 管理员审核 |

**集成测试位置**: `review_integration_test.go`

#### 2.1.4 钱包功能

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 查询余额 | 查看可用余额和冻结金额 | P0 | FrozenCents为T+7冻结 |
| 充值钱包 | 第三方支付充值 | P0 | 支持赠送优惠 |
| 交易记录 | 分页查询交易记录 | P1 | 包含充值/支付/退款/提现 |
| T+7解冻 | 订单完成7天后解冻 | P0 | 定时任务处理 |

**集成测试位置**: `wallet_integration_test.go`, `payment_integration_test.go`

#### 2.1.5 聊天功能

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 公共聊天室 | 发送消息到公共聊天室 | P1 | 需审核后显示 |
| 订单群聊 | 创建订单聊天组 | P0 | 用户+陪玩师自动加入 |
| 发送消息 | 发送文本/图片消息 | P0 | 订单群聊消息免审 |
| 敏感词过滤 | 敏感词自动拦截 | P0 | 触发敏感词拒绝发送 |
| 查看聊天记录 | 分页查询历史消息 | P1 | 仅显示已通过审核的 |
| 标记已读 | 标记消息为已读 | P2 | 更新LastReadMessageID |
| 拉黑用户 | 拉黑后无法互相发消息 | P1 | 双向消息屏蔽 |

**集成测试位置**: `chat_integration_test.go`

### 2.2 陪玩师端功能测试场景

根据PRD §4.2，陪玩师端包含以下功能模块：

#### 2.2.1 入驻认证

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 创建入驻申请 | 提交入驻信息 | P0 | 需填写昵称/擅长游戏/段位 |
| 实名认证 | 提交身份证信息 | P0 | 管理员审核 |
| 段位认证 | 上传段位截图 | P0 | 管理员审核 |
| 审批通过 | 管理员批准入驻 | P0 | 状态变为verified |
| 审批拒绝 | 管理员拒绝入驻 | P0 | 可重新申请 |
| 修改资料 | 更新个人信息 | P1 | 需重新审核敏感字段 |

**集成测试位置**: `player_integration_test.go`

#### 2.2.2 订单管理

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 查看待接订单 | 查看confirmed状态订单 | P0 | 可抢单池 |
| 接单 | 陪玩师接受订单 | P0 | 接单后状态变in_progress |
| 拒单 | 陪玩师拒绝订单 | P1 | 订单回到待接池 |
| 我的订单 | 查看自己的订单列表 | P0 | 支持状态筛选 |
| 完成订单 | 陪玩师标记服务完成 | P0 | 完成后用户可评价 |

**集成测试位置**: `order_integration_test.go`

#### 2.2.3 收益管理

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 收益统计 | 查看总收入/本月收入 | P0 | 区分可用/冻结余额 |
| 创建提现申请 | 申请提现到支付宝/微信 | P0 | 只能提现可用余额 |
| 提现审核 | 管理员审核提现 | P0 | 审核通过后打款 |
| 提现记录 | 查看提现历史 | P1 | 分页查询 |

**集成测试位置**: `withdraw_integration_test.go`

#### 2.2.4 状态管理

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 设置在线状态 | 在线/离线/忙碌 | P0 | 影响是否可以接单 |
| 接单开关 | 开启/关闭接单 | P0 | 关闭后不显示在可接单列表 |

**集成测试位置**: `player_integration_test.go`

#### 2.2.5 团队功能

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 创建团队 | 创建陪玩团队 | P2 | 需团队名称/描述 |
| 邀请成员 | 邀请陪玩师加入 | P2 | 被邀请者接受才能加入 |
| 接受邀请 | 接受团队邀请 | P2 | 成为团队成员 |
| 团队接单 | 团队统一接单 | P2 | 收益按比例分配 |
| 移除成员 | 队长移除成员 | P2 | |
| 收益分配 | 团队订单收益分配 | P2 | 按配置比例分配 |

**集成测试位置**: `team_integration_test.go`

### 2.3 管理后台功能测试场景

根据PRD §4.3，管理后台包含以下功能模块：

#### 2.3.1 仪表盘

| 测试场景 | 描述 | 优先级 |
|---------|------|--------|
| 数据概览 | 查看关键指标 | P1 |
| 趋势图表 | 订单/收入趋势 | P1 |
| 实时监控 | 实时订单/用户统计 | P2 |

**集成测试位置**: `statistics_integration_test.go`

#### 2.3.2 用户管理

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 用户列表 | 分页查询用户 | P0 | 支持按角色/状态筛选 |
| 用户详情 | 查看用户完整信息 | P0 | 包含订单/钱包/评价 |
| 封禁用户 | 封禁违规用户 | P0 | 封禁后无法登录/下单 |
| 解封用户 | 解封用户 | P0 | 恢复正常使用 |
| 添加标签 | 为用户打标签 | P1 | 用于用户分类 |
| 移除标签 | 移除用户标签 | P1 | |

**集成测试位置**: `user_integration_test.go`

#### 2.3.3 陪玩师管理

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 陪玩师列表 | 分页查询陪玩师 | P0 | 支持按认证状态筛选 |
| 陪玩师详情 | 查看陪玩师完整信息 | P0 | 包含认证/订单/收益 |
| 入驻审核 | 审批入驻申请 | P0 | 通过/拒绝 |
| 段位审核 | 审核段位认证 | P0 | 通过/拒绝 |
| 实名审核 | 审核实名认证 | P0 | 查看身份证信息 |
| 调整段位 | 手动调整陪玩师段位 | P1 | 影响定价 |

**集成测试位置**: `player_integration_test.go`

#### 2.3.4 订单管理

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 订单列表 | 分页查询所有订单 | P0 | 支持多条件筛选 |
| 订单详情 | 查看订单完整信息 | P0 | 包含所有关联数据 |
| 退款审核 | 审核退款申请 | P0 | 同意/拒绝退款 |
| 争议处理 | 处理订单争议 | P0 | 分配客服/解决争议 |
| 手动完成 | 管理员手动完成订单 | P1 | 异常情况处理 |

**集成测试位置**: `order_integration_test.go`, `dispute_integration_test.go`

#### 2.3.5 财务管理

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 提现审核 | 审核提现申请 | P0 | 同意/拒绝提现 |
| 佣金配置 | 配置佣金规则 | P1 | 默认/个人/排名 |
| 结算管理 | 查看结算报表 | P1 | 按公司/时间段 |
| 收入统计 | 平台收入统计 | P1 | 日/月/年维度 |

**集成测试位置**: `withdraw_integration_test.go`, `commission_integration_test.go`

#### 2.3.6 内容管理

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 游戏管理 | 创建/编辑/删除游戏 | P0 | |
| 段位管理 | 为游戏配置段位 | P0 | 影响定价 |
| 服务项目管理 | 创建/编辑服务项 | P0 | 配置价格/佣金 |
| 敏感词管理 | 添加/删除敏感词 | P1 | 用于聊天审核 |

**集成测试位置**: `game_integration_test.go`, `serviceitem_integration_test.go`, `sensitiveword_integration_test.go`

#### 2.3.7 系统管理

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 角色管理 | 创建/编辑/删除角色 | P0 | RBAC基础 |
| 权限管理 | 配置角色权限 | P0 | API级别权限控制 |
| 菜单管理 | 配置后台菜单 | P1 | 按角色显示 |
| 操作日志 | 查看操作日志 | P1 | 审计追踪 |

**集成测试位置**: `permission_integration_test.go`, `menu_integration_test.go`

### 2.4 营销功能测试场景

#### 2.4.1 VIP会员

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 获取VIP等级 | 查看当前VIP等级 | P1 | 基于经验值 |
| 增加经验值 | 消费增加经验 | P1 | |
| 等级提升 | 达到阈值自动升级 | P1 | |
| VIP权益 | 查看VIP专属权益 | P1 | 折扣/特权 |
| VIP折扣 | 下单应用VIP折扣 | P1 | 折扣比例按等级 |

**集成测试位置**: `vip_integration_test.go`

#### 2.4.2 优惠券

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 创建优惠券模板 | 配置优惠券规则 | P1 | 金额/门槛/有效期 |
| 发放优惠券 | 给用户发放优惠券 | P1 | |
| 使用优惠券 | 下单使用优惠券 | P1 | 抵扣金额 |
| 过期优惠券 | 使用过期优惠券 | P1 | 应拒绝 |
| 超次数使用 | 超出使用限制 | P1 | 应拒绝 |
| 我的优惠券 | 查看可用优惠券 | P2 | |

**集成测试位置**: `coupon_integration_test.go`

#### 2.4.3 充值

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 获取充值选项 | 查看充值档位 | P1 | |
| 创建充值订单 | 创建充值记录 | P1 | |
| 应用赠送优惠 | 充值赠送活动 | P1 | |
| 充值记录 | 查看充值历史 | P2 | |

**集成测试位置**: `recharge_integration_test.go`

#### 2.4.4 活动

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 创建活动 | 配置活动规则 | P1 | |
| 参加活动 | 用户参与活动 | P1 | |
| 检查资格 | 验证参与条件 | P1 | |
| 领取奖励 | 完成任务领奖励 | P1 | |
| 活动列表 | 查看可用活动 | P2 | |
| 活动过期 | 检查活动过期 | P1 | |

**集成测试位置**: `activity_integration_test.go`

#### 2.4.5 推荐

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 生成推荐码 | 用户获取推荐码 | P2 | |
| 使用推荐码 | 新用户注册使用 | P2 | |
| 记录推荐关系 | 建立推荐绑定 | P2 | |
| 领取奖励 | 达成条件领奖励 | P2 | |
| 推荐列表 | 查看我的推荐 | P2 | |

**集成测试位置**: `referral_integration_test.go`

### 2.5 核心业务规则测试场景

#### 2.5.1 佣金计算

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 默认佣金 | 使用服务项默认佣金 | P0 | 默认20% |
| 个人佣金 | 陪玩师个人专属佣金 | P0 | 覆盖默认佣金 |
| 排名佣金 | 上月排名减免佣金 | P0 | 叠加减免 |
| 三维度计算 | 综合计算最终佣金 | P0 | base - ranking_discount |
| 佣金记录 | 记录每笔佣金 | P1 | 用于结算对账 |

**业务逻辑**:
```
Final Commission Rate = Base Rate - Ranking Discount
Player Income = Order Amount * (1 - Final Commission Rate / 100)

Example:
  Order: ¥100
  Base Rate: 20% (service item default)
  Player Rate: 18% (individual special)
  Ranking: Top 10, -5% discount
  Final: 18% - 5% = 13%
  Player Income: ¥100 * (1 - 0.13) = ¥87
```

**集成测试位置**: `commission_integration_test.go`

#### 2.5.2 订单状态流转

| 测试场景 | 描述 | 优先级 |
|---------|------|--------|
| pending → confirmed | 支付成功 | P0 |
| confirmed → in_progress | 陪玩师接单 | P0 |
| in_progress → completed | 服务完成 | P0 |
| pending → canceled | 用户取消 | P0 |
| confirmed → canceled | 支付后取消 | P1 |
| completed → disputed | 发起争议 | P0 |
| disputed → refunded | 争议退款 | P0 |
| disputed → completed | 争议驳回 | P0 |

**状态图**:
```
                    ┌─────────────┐
                    │   pending   │
                    └──────┬──────┘
                           │ 支付
                           ▼
                    ┌─────────────┐
           取消     │  confirmed  │
          ┌──────────┴──────┬──────┘
          │               │ 接单
          ▼               ▼
    ┌─────────┐    ┌──────────────┐
    │canceled │    │ in_progress  │
    └─────────┘    └───────┬──────┘
                           │ 完成
                           ▼
                    ┌─────────────┐
                    │  completed  │◄─────┐
                    └──────┬──────┘      │
                           │ 争议        │ 驳回
                           ▼             │
                    ┌─────────────┐      │
                    │  disputed   │───┐──┘
                    └──────┬──────┘   │
                           │ 退款      │
                           ▼           │
                    ┌─────────────┐   │
                    │  refunded   │◄──┘
                    └─────────────┘
```

**集成测试位置**: `order_integration_test.go`

#### 2.5.3 收入结算 (T+7)

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 订单完成冻结 | 收入进入FrozenCents | P0 | 订单完成时触发 |
| T+7解冻 | 7天后解冻到BalanceCents | P0 | 定时任务处理 |
| 争议冻结 | 争议期间继续冻结 | P0 | 直到争议解决 |
| 退款扣除 | 退款从FrozenCents扣 | P0 | 不足则从Balance扣 |
| 提现限制 | 只能提现BalanceCents | P0 | Frozen不可提 |

**资金流转**:
```
Order Complete (day 0)
  ↓
FrozenCents += player_income (冻结期开始)
  ↓
  ... 7 days waiting period ...
  ↓
Day 7: No dispute?
  YES → FrozenCents → BalanceCents (可提现)
  NO  → FrozenCents frozen until dispute resolved
  ↓
Dispute Resolution:
  Refund → Deduct from FrozenCents
  Reject → FrozenCents → BalanceCents
```

**集成测试位置**: `wallet_integration_test.go`, `order_integration_test.go`

#### 2.5.4 争议处理

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 用户发起争议 | 订单完成7天内可发起 | P0 | |
| 陪玩师发起争议 | 陪玩师也可发起 | P0 | |
| 双客服机制 | 分配独立客服 | P0 | 避免原订单客服 |
| 30分钟SLA | 客服30分钟内响应 | P0 | 超时标记SLA breach |
| 退款处理 | 同意退款 | P0 | 订单状态变refunded |
| 驳回争议 | 拒绝退款 | P0 | 订单恢复completed |
| 聊天快照 | 保存聊天记录作为证据 | P1 | |

**争议流程**:
```
1. User/Player creates dispute (within 7 days of order completion)
2. System assigns:
   - Original CS (from order, if any)
   - Independent CS (for fairness)
3. CS has 30 minutes to respond (SLA)
4. CS resolution:
   - Refund: order → refunded, money returned to user
   - Reject: order → completed, player keeps income
5. Evidence:
   - Chat snapshot
   - Order details
   - Screenshots
```

**集成测试位置**: `dispute_integration_test.go`

#### 2.5.5 用户拉黑

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 拉黑用户 | 用户拉黑陪玩师 | P1 | |
| 拉黑效果 | 无法互相发消息 | P0 | 双向屏蔽 |
| 隐藏列表 | 互相不可见 | P0 | 搜索/列表中隐藏 |
| 订单隔离 | 已有订单可继续 | P0 | 新订单不可下单给对方 |
| 方向性 | A拉黑B ≠ B拉黑A | P1 | 独立记录 |

**拉黑规则**:
```
User A blocks Player B:
  ├─ A cannot send messages to B
  ├─ B cannot send messages to A
  ├─ A cannot see B in player list
  ├─ B cannot see A in user list
  ├─ A cannot place new orders to B
  ├─ Existing orders continue normally
  └─ Directional: A blocks B ≠ B blocks A
```

**集成测试位置**: `userblock_integration_test.go`

#### 2.5.6 敏感词过滤

| 测试场景 | 描述 | 优先级 | 业务规则 |
|---------|------|--------|----------|
| 检测敏感词 | 文本匹配敏感词 | P0 | |
| 完全匹配 | 精确匹配 | P0 | |
| 模糊匹配 | 包含匹配 | P1 | |
| 替换处理 | 替换为*** | P1 | |
| 拒绝发送 | 直接拒绝 | P0 | 高敏感度词汇 |
| 分类管理 | 按类别管理 | P2 | 政治/色情/广告 |

**集成测试位置**: `sensitiveword_integration_test.go`, `chat_integration_test.go`

---

## 三、测试场景按模块分类

### 3.1 用户管理 (user)

| 测试用例 | 描述 | 优先级 | 状态 |
|---------|------|--------|------|
| `TestUserService_Register` | 用户注册 | P0 | ⏳ 待实现 |
| `TestUserService_Register_DuplicatePhone` | 重复手机号 | P0 | ⏳ 待实现 |
| `TestUserService_Register_DuplicateEmail` | 重复邮箱 | P0 | ⏳ 待实现 |
| `TestUserService_Login` | 用户登录 | P0 | ✅ 已在auth |
| `TestUserService_UpdateProfile` | 更新资料 | P1 | ⏳ 待实现 |
| `TestUserService_ChangePassword` | 修改密码 | P1 | ⏳ 待实现 |
| `TestUserService_GetUserByID` | 获取用户信息 | P1 | ⏳ 待实现 |
| `TestUserService_ListUsers` | 用户列表 | P1 | ⏳ 待实现 |
| `TestUserService_BanUser` | 封禁用户 | P0 | ⏳ 待实现 |
| `TestUserService_UnbanUser` | 解封用户 | P0 | ⏳ 待实现 |
| `TestUserService_AddUserTags` | 添加标签 | P2 | ⏳ 待实现 |
| `TestUserService_RemoveUserTags` | 移除标签 | P2 | ⏳ 待实现 |

**预估工作量**: 4-5小时

### 3.2 陪玩师管理 (player)

| 测试用例 | 描述 | 优先级 | 状态 |
|---------|------|--------|------|
| `TestPlayerService_CreateApplication` | 创建入驻申请 | P0 | ⏳ 待实现 |
| `TestPlayerService_ApproveApplication` | 审批通过 | P0 | ⏳ 待实现 |
| `TestPlayerService_RejectApplication` | 审批拒绝 | P0 | ⏳ 待实现 |
| `TestPlayerService_UpdateProfile` | 更新资料 | P1 | ⏳ 待实现 |
| `TestPlayerService_SetOnlineStatus` | 设置在线状态 | P0 | ⏳ 待实现 |
| `TestPlayerService_SetAcceptingOrders` | 接单开关 | P0 | ⏳ 待实现 |
| `TestPlayerService_GetPlayerByID` | 获取详情 | P1 | ✅ 已有 |
| `TestPlayerService_ListPlayers` | 陪玩师列表 | P1 | ✅ 已有 |
| `TestPlayerService_SearchPlayers` | 搜索 | P2 | ⏳ 待实现 |
| `TestPlayerService_GetPlayerEarnings` | 获取收益 | P1 | ⏳ 待实现 |
| `TestPlayerService_RankCertification` | 段位认证 | P0 | ⏳ 待实现 |
| `TestPlayerService_RealNameCertification` | 实名认证 | P0 | ⏳ 待实现 |

**预估工作量**: 5-6小时

### 3.3 订单管理 (order)

| 测试用例 | 描述 | 优先级 | 状态 |
|---------|------|--------|------|
| `TestOrderService_CreateOrder_Solo` | 创建单人订单 | P0 | ✅ 已有 |
| `TestOrderService_CreateOrder_Team` | 创建多人订单 | P0 | ⏳ 待实现 |
| `TestOrderService_CreateOrder_Gift` | 创建礼物订单 | P0 | ⏳ 待实现 |
| `TestOrderService_CreateOrder_WithCoupon` | 使用优惠券 | P1 | ⏳ 待实现 |
| `TestOrderService_CreateOrder_WithVIPDiscount` | 使用VIP折扣 | P1 | ⏳ 待实现 |
| `TestOrderService_AcceptOrder` | 接单 | P0 | ✅ 已有 |
| `TestOrderService_RejectOrder` | 拒单 | P1 | ⏳ 待实现 |
| `TestOrderService_StartOrder` | 开始服务 | P0 | ⏳ 待实现 |
| `TestOrderService_CompleteOrder` | 完成订单 | P0 | ✅ 已有 |
| `TestOrderService_CancelOrder` | 取消订单 | P0 | ✅ 已有 |
| `TestOrderService_SupplementOrder` | 补差价订单 | P1 | ⏳ 待实现 |
| `TestOrderService_GetMyOrders` | 我的订单 | P0 | ✅ 已有 |
| `TestOrderService_GetOrderDetail` | 订单详情 | P0 | ✅ 已有 |
| `TestOrderService_GetAvailableOrders` | 待接订单池 | P0 | ✅ 已有 |
| `TestOrderService_Timeout_Cancel` | 超时自动取消 | P1 | ⏳ 待实现 |
| `TestOrderService_Timeout_AssignCS` | 超时分配客服 | P1 | ⏳ 待实现 |

**预估工作量**: 6-7小时（含新增场景）

### 3.4 支付管理 (payment)

| 测试用例 | 描述 | 优先级 | 状态 |
|---------|------|--------|------|
| `TestPaymentService_CreatePayment_WeChat` | 微信支付 | P0 | ✅ 已有 |
| `TestPaymentService_CreatePayment_Alipay` | 支付宝支付 | P0 | ⏳ 待实现 |
| `TestPaymentService_CreatePayment_WalletOnly` | 钱包支付 | P0 | ✅ 已有 |
| `TestPaymentService_CreatePayment_Combined` | 组合支付 | P0 | ✅ 已有 |
| `TestPaymentService_WalletInsufficient` | 余额不足 | P0 | ✅ 已有 |
| `TestPaymentService_PaymentCallback_Success` | 支付成功回调 | P0 | ✅ 已有 |
| `TestPaymentService_PaymentCallback_Failed` | 支付失败回调 | P0 | ⏳ 待实现 |
| `TestPaymentService_RefundPayment` | 退款 | P0 | ✅ 已有 |
| `TestPaymentService_Refund_Combined` | 组合支付退款 | P1 | ✅ 已有 |
| `TestPaymentService_CancelPayment` | 取消支付 | P0 | ✅ 已有 |
| `TestPaymentService_GetPaymentStatus` | 支付状态 | P1 | ✅ 已有 |
| `TestPaymentService_CalculateCombined` | 计算组合支付 | P1 | ✅ 已有 |

**预估工作量**: 4-5小时（含新增场景）

### 3.5 提现管理 (withdraw)

| 测试用例 | 描述 | 优先级 | 状态 |
|---------|------|--------|------|
| `TestWithdrawService_CreateWithdraw` | 创建提现申请 | P0 | ⏳ 待实现 |
| `TestWithdrawService_CreateWithdraw_Insufficient` | 余额不足 | P0 | ⏳ 待实现 |
| `TestWithdrawService_CreateWithdraw_ExceedsDailyLimit` | 超出每日限额 | P0 | ⏳ 待实现 |
| `TestWithdrawService_ApproveWithdraw` | 审批通过 | P0 | ⏳ 待实现 |
| `TestWithdrawService_RejectWithdraw` | 审批拒绝 | P0 | ⏳ 待实现 |
| `TestWithdrawService_ProcessWithdraw` | 处理提现 | P1 | ⏳ 待实现 |
| `TestWithdrawService_RouteWithdrawal` | 路由到结算公司 | P0 | ✅ 已有 |
| `TestWithdrawService_GetPlayerBalance` | 获取余额 | P1 | ✅ 已有 |
| `TestWithdrawService_GetWithdrawal` | 获取提现详情 | P1 | ✅ 已有 |
| `TestWithdrawService_CompleteWithdrawalPayment` | 完成打款 | P1 | ✅ 已有 |

**预估工作量**: 4-5小时（含新增场景）

### 3.6 争议管理 (dispute)

| 测试用例 | 描述 | 优先级 | 状态 |
|---------|------|--------|------|
| `TestDisputeService_InitiateDispute_ByUser` | 用户发起争议 | P0 | ✅ 已有 |
| `TestDisputeService_InitiateDispute_ByPlayer` | 陪玩师发起 | P0 | ✅ 已有 |
| `TestDisputeService_InitiateDispute_Unauthorized` | 无权限发起 | P0 | ✅ 已有 |
| `TestDisputeService_InitiateDispute_After7Days` | 超过7天 | P0 | ⏳ 待实现 |
| `TestDisputeService_AssignDispute_OriginalCS` | 分配原客服 | P0 | ⏳ 待实现 |
| `TestDisputeService_AssignDispute_IndependentCS` | 分配独立客服 | P0 | ⏳ 待实现 |
| `TestDisputeService_ResolveDispute_Refund` | 退款处理 | P0 | ✅ 已有 |
| `TestDisputeService_ResolveDispute_Reject` | 驳回争议 | P0 | ✅ 已有 |
| `TestDisputeService_CheckSLABreaches` | SLA超时检查 | P0 | ✅ 已有 |
| `TestDisputeService_GetDisputeDetail` | 争议详情 | P1 | ✅ 已有 |
| `TestDisputeService_ListPendingDisputes` | 待处理列表 | P1 | ✅ 已有 |

**预估工作量**: 4-5小时（含新增场景）

### 3.7 评价管理 (review)

| 测试用例 | 描述 | 优先级 | 状态 |
|---------|------|--------|------|
| `TestReviewService_CreateReview` | 创建评价 | P0 | ⏳ 待实现 |
| `TestReviewService_CreateReview_OneClick` | 一键好评 | P1 | ⏳ 待实现 |
| `TestReviewService_CreateReview_BeforeComplete` | 完成前评价 | P0 | ⏳ 待实现 |
| `TestReviewService_UpdateReview` | 修改评价 | P1 | ⏳ 待实现 |
| `TestReviewService_ReplyReview` | 陪玩师回复 | P2 | ⏳ 待实现 |
| `TestReviewService_AppealReview` | 评价申诉 | P2 | ⏳ 待实现 |
| `TestReviewService_GetReviews` | 评价列表 | P1 | ⏳ 待实现 |
| `TestReviewService_GetReviewStats` | 评价统计 | P2 | ⏳ 待实现 |

**预估工作量**: 3-4小时

### 3.8 聊天管理 (chat)

| 测试用例 | 描述 | 优先级 | 状态 |
|---------|------|--------|------|
| `TestChatService_JoinGroup` | 加入群组 | P0 | ✅ 已有 |
| `TestChatService_LeaveGroup` | 离开群组 | P0 | ✅ 已有 |
| `TestChatService_SendMessage_Text` | 发送文本 | P0 | ✅ 已有 |
| `TestChatService_SendMessage_Image` | 发送图片 | P1 | ⏳ 待实现 |
| `TestChatService_SendMessage_SensitiveWord` | 敏感词拦截 | P0 | ⏳ 待实现 |
| `TestChatService_ListMessages` | 消息列表 | P0 | ✅ 已有 |
| `TestChatService_MarkRead` | 标记已读 | P2 | ✅ 已有 |
| `TestChatService_ApproveMessage` | 审核通过 | P1 | ✅ 已有 |
| `TestChatService_RejectMessage` | 审核拒绝 | P1 | ✅ 已有 |
| `TestChatService_ReportMessage` | 举报消息 | P2 | ✅ 已有 |
| `TestChatService_CreateOrderGroup` | 创建订单群 | P0 | ✅ 已有 |
| `TestChatService_UserBlockEffect` | 拉黑效果 | P0 | ⏳ 待实现 |

**预估工作量**: 4-5小时（含新增场景）

### 3.9 用户拉黑 (userblock)

| 测试用例 | 描述 | 优先级 | 状态 |
|---------|------|--------|------|
| `TestUserBlockService_BlockPlayer` | 拉黑陪玩师 | P0 | ⏳ 待实现 |
| `TestUserBlockService_UnblockPlayer` | 取消拉黑 | P0 | ⏳ 待实现 |
| `TestUserBlockService_CheckBlocked` | 检查拉黑状态 | P0 | ⏳ 待实现 |
| `TestUserBlockService_BlockEffect_Message` | 消息屏蔽 | P0 | ⏳ 待实现 |
| `TestUserBlockService_BlockEffect_List` | 列表隐藏 | P0 | ⏳ 待实现 |
| `TestUserBlockService_BlockEffect_Order` | 订单限制 | P0 | ⏳ 待实现 |
| `TestUserBlockService_GetBlockedUsers` | 拉黑列表 | P1 | ⏳ 待实现 |

**预估工作量**: 3-4小时

### 3.10 钱包管理 (wallet)

| 测试用例 | 描述 | 优先级 | 状态 |
|---------|------|--------|------|
| `TestWalletService_GetBalance` | 获取余额 | P0 | ⏳ 待实现 |
| `TestWalletService_FreezeOnOrderComplete` | 完成订单冻结 | P0 | ⏳ 待实现 |
| `TestWalletService_ThawAfter7Days` | T+7解冻 | P0 | ⏳ 待实现 |
| `TestWalletService_FreezeOnDispute` | 争议冻结 | P0 | ⏳ 待实现 |
| `TestWalletService_DeductOnRefund` | 退款扣除 | P0 | ⏳ 待实现 |
| `TestWalletService_GetTransactionHistory` | 交易记录 | P1 | ⏳ 待实现 |
| `TestWalletService_CalculateAvailableBalance` | 可用余额计算 | P0 | ⏳ 待实现 |

**预估工作量**: 3-4小时

### 3.11 佣金管理 (commission)

| 测试用例 | 描述 | 优先级 | 状态 |
|---------|------|--------|------|
| `TestCommissionService_Calculate_Default` | 默认佣金 | P0 | ⏳ 待实现 |
| `TestCommissionService_Calculate_PlayerSpecific` | 个人专属佣金 | P0 | ⏳ 待实现 |
| `TestCommissionService_Calculate_RankingDiscount` | 排名减免 | P0 | ⏳ 待实现 |
| `TestCommissionService_Calculate_Composite` | 综合计算 | P0 | ⏳ 待实现 |
| `TestCommissionService_CreateCommissionRule` | 创建佣金规则 | P1 | ⏳ 待实现 |
| `TestCommissionService_RecordCommission` | 记录佣金 | P1 | ⏳ 待实现 |
| `TestCommissionService_GetCommissionRecords` | 佣金记录 | P2 | ⏳ 待实现 |

**预估工作量**: 3-4小时

### 3.12 权限管理 (permission/role/menu)

| 测试用例 | 描述 | 优先级 | 状态 |
|---------|------|--------|------|
| `TestPermissionService_CheckPermission` | 检查权限 | P0 | ✅ 已有 |
| `TestPermissionService_AssignRoleToUser` | 分配角色 | P0 | ✅ 已有 |
| `TestPermissionService_RevokeRoleFromUser` | 撤销角色 | P1 | ⏳ 待实现 |
| `TestPermissionService_CreateRole` | 创建角色 | P1 | ⏳ 待实现 |
| `TestPermissionService_AssignPermissionsToRole` | 分配权限 | P1 | ⏳ 待实现 |
| `TestPermissionService_GetUserPermissions` | 获取用户权限 | P1 | ⏳ 待实现 |
| `TestMenuService_GetMenuTree` | 获取菜单树 | P1 | ✅ 已有 |
| `TestMenuService_CreateMenu` | 创建菜单 | P1 | ⏳ 待实现 |

**预估工作量**: 3-4小时（含新增场景）

### 3.13 游戏和服务项管理 (game/service-item)

| 测试用例 | 描述 | 优先级 | 状态 |
|---------|------|--------|------|
| `TestGameService_CreateGame` | 创建游戏 | P0 | ⏳ 待实现 |
| `TestGameService_UpdateGame` | 更新游戏 | P1 | ⏳ 待实现 |
| `TestGameService_GetActiveGames` | 活跃游戏列表 | P1 | ⏳ 待实现 |
| `TestServiceItemService_CreateItem` | 创建服务项 | P0 | ⏳ 待实现 |
| `TestServiceItemService_CalculatePrice` | 计算价格 | P0 | ⏳ 待实现 |
| `TestServiceItemService_GetItemsByGame` | 按游戏查询 | P1 | ⏳ 待实现 |

**预估工作量**: 3小时

### 3.14 营销模块 (vip/coupon/recharge/activity/team/referral)

**状态**: 测试文件已创建，内容待补充

| 模块 | 测试用例数 | 预估工作量 | 状态 |
|------|-----------|-----------|------|
| vip | 5-8 | 2-3h | ⏳ 待补充 |
| coupon | 8-10 | 3-4h | ⏳ 待补充 |
| recharge | 4-6 | 2h | ⏳ 待补充 |
| activity | 6-8 | 3-4h | ⏳ 待补充 |
| team | 7-9 | 3-4h | ⏳ 待补充 |
| referral | 5-7 | 2-3h | ⏳ 待补充 |

**预估总工作量**: 15-20小时

---

## 四、测试执行计划

### 4.1 Phase 1: 核心业务流程补充（第1-2周）

**目标**: 补充核心业务流程缺失的集成测试

| 模块 | 测试用例数 | 预估时间 | 优先级 |
|------|-----------|---------|--------|
| user | 12 | 4-5h | P0 |
| order (补充) | 8 | 3-4h | P0 |
| payment (补充) | 4 | 2h | P1 |
| withdraw (补充) | 6 | 3-4h | P0 |
| dispute (补充) | 3 | 1-2h | P1 |
| review | 8 | 3-4h | P0 |
| userblock | 7 | 3-4h | P0 |
| wallet | 7 | 3-4h | P0 |
| commission | 7 | 3-4h | P0 |
| **小计** | **62** | **25-33h** | - |

### 4.2 Phase 2: 辅助功能模块（第3周）

| 模块 | 测试用例数 | 预估时间 | 优先级 |
|------|-----------|---------|--------|
| permission (补充) | 6 | 2-3h | P1 |
| chat (补充) | 4 | 2h | P1 |
| game | 6 | 3h | P1 |
| service-item | 5 | 3h | P1 |
| **小计** | **21** | **10-12h** | - |

### 4.3 Phase 3: 营销模块完善（第4周）

| 模块 | 测试用例数 | 预估时间 | 优先级 |
|------|-----------|---------|--------|
| vip | 8 | 2-3h | P1 |
| coupon | 10 | 3-4h | P1 |
| recharge | 6 | 2h | P2 |
| activity | 8 | 3-4h | P2 |
| team | 9 | 3-4h | P2 |
| referral | 7 | 2-3h | P2 |
| **小计** | **48** | **15-20h** | - |

### 4.4 总体统计

| 阶段 | 模块数 | 测试用例数 | 预估工作量 |
|------|-------|-----------|-----------|
| Phase 1 | 9 | 62 | 25-33h |
| Phase 2 | 4 | 21 | 10-12h |
| Phase 3 | 6 | 48 | 15-20h |
| **已完成** | **10** | **~90** | **-** |
| **总计** | **29** | **~221** | **50-65h** |

---

## 五、测试基础设施

### 5.1 测试数据库配置

**文件**: [docker-compose.test.yml](../docker-compose.test.yml)

```bash
# 启动测试数据库
docker compose -f docker-compose.test.yml up -d

# 环境变量
TEST_DB_HOST=localhost
TEST_DB_PORT=5432
TEST_DB_USER=gamelink
TEST_DB_PASSWORD=gamelink
TEST_DB_NAME=gamelink_test
```

### 5.2 测试工具函数

**文件**: [testdb.go](../api/internal/service/integration/testdb.go)

**已实现的工具函数**:
- `SetupTestDB(t)` - 初始化测试数据库
- `SkipIfNoTestDB(t)` - 跳过测试
- `CreateTestUser/CreateUniqueTestUser` - 创建测试用户
- `CreateTestPlayer` - 创建测试陪玩师
- `CreateTestGame` - 创建测试游戏
- `CreateTestOrder/CreateTestOrderWithDetails` - 创建测试订单
- `CreateTestPayment` - 创建测试支付
- `CreateTestWallet` - 创建测试钱包
- `CreateTestServiceItem` - 创建测试服务项
- `CreateTestReview` - 创建测试评价
- `CreateTestDispute` - 创建测试争议
- `CreateTestChatGroup/CreateTestChatMessage` - 创建测试聊天
- `CreateTestNotification` - 创建测试通知
- `CreateTestVipLevel` - 创建测试VIP
- `CreateTestCouponTemplate` - 创建测试优惠券
- `CreateTestCommissionRule` - 创建测试佣金规则
- `CreateTestSettlementCompany` - 创建测试结算公司
- `CreateTestWithdraw` - 创建测试提现
- `CreateTestRole/CreateTestPermission` - 创建测试角色权限
- `CreateTestMenu` - 创建测试菜单
- `CreateTestSensitiveWord` - 创建测试敏感词
- `AssignRoleToUser/AssignPermissionToRole` - 分配关系

### 5.3 测试命名规范

```go
// 基础测试
func Test{ServiceName}_{MethodName}(t *testing.T) {
    SkipIfNoTestDB(t)
    db := SetupTestDB(t)
    ctx := context.Background()
    // ...
}

// 场景测试
func Test{ServiceName}_{MethodName}_{Scenario}(t *testing.T) {
    // Example: TestOrderService_CreateOrder_WithCoupon
}

// 错误场景
func Test{ServiceName}_{MethodName}_{ErrorType}(t *testing.T) {
    // Example: TestPaymentService_CreatePayment_InsufficientBalance
}
```

---

## 六、执行命令

### 6.1 运行集成测试

```bash
# 启动测试数据库
docker compose -f docker-compose.test.yml up -d

# 运行所有集成测试
go test ./api/internal/service/integration/... -v

# 运行特定模块
go test ./api/internal/service/integration/order_integration_test.go -v

# 运行特定测试用例
go test ./api/internal/service/integration/... -run TestOrderService_CreateOrder -v

# 查看覆盖率
go test ./api/internal/service/integration/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out

# 运行特定模块的覆盖率
go test ./api/internal/service/integration/order_integration_test.go -cover -coverprofile=coverage_order.out
```

### 6.2 CI/CD集成

```yaml
# .github/workflows/ci.yml
- name: Start test database
  run: docker compose -f docker-compose.test.yml up -d

- name: Run integration tests
  run: go test ./api/internal/service/integration/... -v
  env:
    TEST_DB_HOST: localhost
    TEST_DB_PORT: 5432
    TEST_DB_USER: gamelink
    TEST_DB_PASSWORD: gamelink
    TEST_DB_NAME: gamelink_test

- name: Upload coverage
  uses: codecov/codecov-action@v3
  with:
    files: ./coverage.out
```

---

## 七、成功标准

### 7.1 覆盖率目标

| 模块类别 | 目标覆盖率 |
|---------|-----------|
| 核心业务模块 | ≥ 80% |
| 营销模块 | ≥ 70% |
| 辅助模块 | ≥ 60% |
| **整体目标** | **≥ 75%** |

### 7.2 质量标准

- ✅ 所有测试用例必须通过
- ✅ 每个Service至少覆盖：正常场景、边界场景、错误场景
- ✅ 关键业务流程必须有端到端测试
- ✅ PRD中所有P0功能点必须有对应测试
- ✅ 测试执行时间 < 10分钟（单次完整运行）

### 7.3 验收标准

1. **代码覆盖率**: Service层 ≥ 75%
2. **测试通过率**: 100%
3. **PRD覆盖**: P0功能点 100%覆盖
4. **关键路径覆盖**: 订单、支付、评价、争议、佣金 100%
5. **CI/CD集成**: 集成测试自动运行

---

## 八、注意事项

### 8.1 测试隔离

- 每个测试用例独立运行
- 使用 `SetupTestDB` 自动清理数据库
- 避免测试间数据依赖

### 8.2 测试数据

- 使用 `CreateUniqueTestUser` 等函数确保数据唯一性
- 避免硬编码ID，使用动态生成的数据
- 测试数据要有代表性，覆盖边界值

### 8.3 业务规则验证

- 所有测试必须验证业务规则（如T+7、佣金计算等）
- 重点验证PRD中定义的业务逻辑
- 验证状态流转的正确性

### 8.4 Mock策略

- 集成测试使用真实数据库
- 外部服务（如支付网关）使用Mock
- 避免依赖外部不稳定服务

---

## 九、附录

### 9.1 相关文档

- [PRD文档](../docs/PRD.md) - 产品需求文档
- [测试基础设施](../api/internal/service/integration/testdb.go)
- [项目管理规则](../.kiro/steering/06-project-management.md)
- [数据模型文档](../.kiro/steering/04-data-models.md)

### 9.2 参考资料

- [Go Testing指南](https://golang.org/pkg/testing/)
- [Testify断言库](https://github.com/stretchr/testify)
- [GORM测试最佳实践](https://gorm.io/docs/tests.html)

---

## 十、更新日志

| 日期 | 版本 | 变更内容 |
|------|------|----------|
| 2025-12-29 | v2.0 | 基于PRD全面规划测试场景，补充业务规则测试 |
| 2025-12-28 | v1.1 | Phase 1 P0 核心业务流程测试完成编写 |
| 2025-12-28 | v1.0 | 初始版本，完成测试规划 |

---

**文档维护者**: AI Assistant
**最后更新**: 2025-12-29
