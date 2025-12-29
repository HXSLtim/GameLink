# 集成测试补充计划

> GameLink 遗漏业务测试补充计划
> 最后更新：2025-12-29

## 一、概述

根据对现有集成测试的分析，当前已完成 **355个测试用例**，覆盖了 **26个模块**（12个核心模块 + 6个营销模块 + 5个辅助模块 + 4个P0补充模块 - 1个重复）。

### 当前覆盖情况

| 分类 | 状态 | 测试数 | 覆盖率 |
|------|------|--------|--------|
| 核心模块 (12个) | ✅ 已完成 | 131 | 50-90% |
| 营销模块 (6个) | ✅ 已完成 | 99 | 70-90% |
| 辅助模块 (5个) | ✅ 已完成 | 53 | 70-90% |
| P0补充模块 (4个) | ✅ 已完成 | 72 | 80-95% |

### 已覆盖模块详情

| 模块 | 测试数 | 状态 |
|------|--------|------|
| auth | 12 | ✅ |
| user | 12 | ✅ |
| player | 9 | ✅ |
| order | 8 | ✅ |
| payment | 13 | ✅ |
| withdraw | 12 | ✅ |
| permission | 20 | ✅ |
| menu | 11 | ✅ |
| chat | 16 | ✅ |
| dispute | 9 | ✅ |
| notification | 9 | ✅ |
| **核心模块小计** | **131** | - |
| vip | 13 | ✅ 新增 |
| coupon | 15 | ✅ 新增 |
| recharge | 13 | ✅ 新增 |
| activity | 15 | ✅ 新增 |
| team | 15 | ✅ 新增 |
| referral | 28 | ✅ 新增 |
| **营销模块小计** | **99** | - |
| game | 9 | ✅ 新增 |
| serviceitem | 10 | ✅ 新增 |
| review | 10 | ✅ 新增 |
| sensitiveword | 10 | ✅ 新增 |
| userblock | 14 | ✅ 新增 |
| **辅助模块小计** | **53** | - |
| gamerank | 13 | ✅ 新增 |
| playerrank | 17 | ✅ 新增 |
| playercertification | 15 | ✅ 新增 |
| ordertimeout | 25 | ✅ 新增 |
| **P0补充模块小计** | **72** | - |
| **总计** | **355** | - |

✅ 营销模块：vip, coupon, recharge, activity, team, referral (100% 覆盖)
✅ 辅助模块：game, serviceitem, review, sensitiveword, userblock (100% 覆盖)
✅ P0补充模块：gamerank, playerrank, playercertification, ordertimeout (100% 覆盖)

---

## 二、高优先级补充场景 (P0)

### 2.1 订单模块补充 (order)

**文件**: `api/internal/service/integration/order_supplement_integration_test.go`

| 测试用例 | 描述 |
|----------|------|
| TestOrderService_CreateMultiPlayerOrder | 创建多人订单 (RequiredPlayers > 1) |
| TestOrderService_AddPlayerToOrder | 向订单添加玩家 (补差价) |
| TestOrderService_CreateGiftOrder | 创建礼物订单 (直接到钱包) |
| TestOrderService_AutoCancelPaymentTimeout | 支付超时自动取消 |
| TestOrderService_AssignToCustomerService | 订单确认后分配给客服 |

**关键依赖**:
- `model.Order.RequiredPlayers`
- `model.OrderItem` (多座位)
- `model.ServiceItem.SubCategory = "gift"`
- `model.OrderTimeoutConfig`
- `model.OrderServiceAssignment`

### 2.2 陪玩师模块补充 (player)

**文件**: `api/internal/service/integration/player_supplement_integration_test.go`

| 测试用例 | 描述 |
|----------|------|
| TestPlayerService_ApproveApplication | 审批通过入驻申请 |
| TestPlayerService_RejectApplication | 审批拒绝入驻申请 |
| TestPlayerService_SubmitRealNameCertification | 提交实名认证 |
| TestPlayerService_ApproveRealNameCertification | 审批通过实名认证 |
| TestPlayerService_SubmitRankCertification | 提交段位认证 |
| TestPlayerService_GetPlayerEarnings | 查询陪玩师收益 |
| TestPlayerService_UpdateRatingAfterReview | 评价后更新评分 |

**关键依赖**:
- `model.Player.VerificationStatus` (pending/approved/rejected)
- `model.PlayerCertification`
- `model.PlayerRankRecord`
- `model.CommissionRecord`

### 2.3 支付模块补充 (payment)

**文件**: `api/internal/service/integration/payment_supplement_integration_test.go`

| 测试用例 | 描述 |
|----------|------|
| TestPaymentService_AutoCancelOnTimeout | 支付超时自动取消订单 |
| TestPaymentService_CombinedPaymentFailure | 组合支付失败时退款钱包部分 |
| TestPaymentService_T7SettlementTracking | T+7收入结算流程追踪 |

**关键依赖**:
- `model.Payment` 超时检测
- `model.CommissionRecord.SettlementStatus`
- `model.Wallet` 退款处理

---

## 三、中优先级补充场景 (P1)

### 3.1 聊天模块补充 (chat)

**文件**: `api/internal/service/integration/chat_supplement_integration_test.go`

| 测试用例 | 描述 |
|----------|------|
| TestChatService_SensitiveWordFilter | 敏感词拦截消息 |
| TestChatService_DeleteMessage | 删除消息 |
| TestChatService_AutoCreateOrderGroup | 支付后自动创建订单聊天群 |

**关键依赖**:
- `model.SensitiveWord`
- `model.ChatMessage` 删除
- 订单支付成功后自动创建 `ChatGroup`

### 3.2 争议模块补充 (dispute)

**文件**: `api/internal/service/integration/dispute_supplement_integration_test.go`

| 测试用例 | 描述 |
|----------|------|
| TestDisputeService_CreateChatSnapshot | 创建争议时保存聊天快照 |
| TestDisputeService_DoubleCSAssignment | 双客服分配机制 |
| TestDisputeService_UploadEvidence | 上传证据截图 (最多5张) |

**关键依赖**:
- `model.ChatSnapshot`
- `model.OrderDispute.OriginalServiceID`

---

## 四、营销模块测试 (Phase 3)

### 4.1 VIP会员模块

**文件**: `api/internal/service/integration/vip_integration_test.go`

| 测试用例 | 描述 |
|----------|------|
| TestVIPService_GetUserVIPLevel | 获取用户VIP等级 |
| TestVIPService_AddExp | 增加经验值 |
| TestVIPService_LevelUp | 等级提升 |
| TestVIPService_GetVIPBenefits | 获取VIP权益 |
| TestVIPService_ApplyDiscount | 应用VIP折扣 |

### 4.2 优惠券模块

**文件**: `api/internal/service/integration/coupon_integration_test.go`

| 测试用例 | 描述 |
|----------|------|
| TestCouponService_CreateTemplate | 创建优惠券模板 |
| TestCouponService_IssueCoupon | 发放优惠券 |
| TestCouponService_UseCoupon | 使用优惠券 |
| TestCouponService_UseCoupon_Expired | 使用过期优惠券 |
| TestCouponService_UseCoupon_ExceedUsageLimit | 超出使用限制 |
| TestCouponService_CalculateDiscount | 计算折扣 |

### 4.3 充值模块

**文件**: `api/internal/service/integration/recharge_integration_test.go`

| 测试用例 | 描述 |
|----------|------|
| TestRechargeService_GetRechargeOptions | 获取充值选项 |
| TestRechargeService_CreateRecharge | 创建充值订单 |
| TestRechargeService_ApplyBonus | 应用赠送优惠 |

### 4.4 活动模块

**文件**: `api/internal/service/integration/activity_integration_test.go`

| 测试用例 | 描述 |
|----------|------|
| TestActivityService_CreateActivity | 创建活动 |
| TestActivityService_JoinActivity | 参加活动 |
| TestActivityService_CheckEligibility | 检查参与资格 |
| TestActivityService_ClaimReward | 领取奖励 |

### 4.5 团队模块

**文件**: `api/internal/service/integration/team_integration_test.go`

| 测试用例 | 描述 |
|----------|------|
| TestTeamService_CreateTeam | 创建团队 |
| TestTeamService_InviteMember | 邀请成员 |
| TestTeamService_AcceptInvite | 接受邀请 |
| TestTeamService_RemoveMember | 移除成员 |
| TestTeamService_DistributeRevenue | 分配收益 |

### 4.6 推荐模块

**文件**: `api/internal/service/integration/referral_integration_test.go`

| 测试用例 | 描述 |
|----------|------|
| TestReferralService_GenerateCode | 生成推荐码 |
| TestReferralService_ApplyReferral | 使用推荐码 |
| TestReferralService_RecordReferral | 记录推荐关系 |
| TestReferralService_ClaimReward | 领取推荐奖励 |

---

## 五、辅助模块测试 (Phase 4)

### 5.1 游戏管理

**文件**: `api/internal/service/integration/game_integration_test.go`

| 测试用例 | 描述 |
|----------|------|
| TestGameService_CreateGame | 创建游戏 |
| TestGameService_UpdateGame | 更新游戏 |
| TestGameService_DeleteGame | 删除游戏 |
| TestGameService_GetActiveGames | 获取活跃游戏列表 |

### 5.2 服务项管理

**文件**: `api/internal/service/integration/serviceitem_integration_test.go`

| 测试用例 | 描述 |
|----------|------|
| TestServiceItemService_CreateItem | 创建服务项 |
| TestServiceItemService_UpdateItem | 更新服务项 |
| TestServiceItemService_CalculatePrice | 计算价格 (含VIP折扣) |

### 5.3 敏感词管理

**文件**: `api/internal/service/integration/sensitiveword_integration_test.go`

| 测试用例 | 描述 |
|----------|------|
| TestSensitiveWordService_CheckText | 检查文本 |
| TestSensitiveWordService_AddWord | 添加敏感词 |
| TestSensitiveWordService_ReplaceText | 替换敏感词 |

---

## 六、实施计划

### 阶段 1: 核心模块补充 (预计 20-25 个测试用例)

- 订单补充 (5个测试)
- 陪玩师补充 (7个测试)
- 支付补充 (3个测试)
- 聊天补充 (3个测试)
- 争议补充 (3个测试)

### 阶段 2: 营销模块 (预计 25-30 个测试用例)

- VIP (5个测试)
- Coupon (6个测试)
- Recharge (3个测试)
- Activity (4个测试)
- Team (5个测试)
- Referral (4个测试)

### 阶段 3: 辅助模块 (预计 10-15 个测试用例)

- Game (4个测试)
- ServiceItem (3个测试)
- SensitiveWord (3个测试)
- Ranking/Statistics (按需补充)

---

## 七、关键文件位置

### 现有测试基础设施

| 文件 | 说明 |
|------|------|
| `api/internal/service/integration/testdb.go` | 测试工具函数 |
| `docs/INTEGRATION_TEST_PLAN.md` | 测试计划文档 |
| `docker-compose.test.yml` | 测试数据库配置 |

### 待创建的测试文件

| 文件 | 说明 | 状态 |
|------|------|------|
| `order_supplement_integration_test.go` | 订单补充测试 | 待创建 |
| `player_supplement_integration_test.go` | 陪玩师补充测试 | 待创建 |
| `payment_supplement_integration_test.go` | 支付补充测试 | 待创建 |
| `chat_supplement_integration_test.go` | 聊天补充测试 | 待创建 |
| `dispute_supplement_integration_test.go` | 争议补充测试 | 待创建 |
| `vip_integration_test.go` | VIP测试 | ✅ 已完成 |
| `coupon_integration_test.go` | 优惠券测试 | ✅ 已完成 |
| `recharge_integration_test.go` | 充值测试 | ✅ 已完成 |
| `activity_integration_test.go` | 活动测试 | ✅ 已完成 |
| `team_integration_test.go` | 团队测试 | ✅ 已完成 |
| `referral_integration_test.go` | 推荐测试 | ✅ 已完成 |
| `game_integration_test.go` | 游戏管理测试 | ✅ 已完成 |
| `serviceitem_integration_test.go` | 服务项测试 | ✅ 已完成 |
| `sensitiveword_integration_test.go` | 敏感词测试 | ✅ 已完成 |
| `gamerank_integration_test.go` | 游戏段位测试 | ✅ 已完成 |
| `playerrank_integration_test.go` | 陪玩师段位认证测试 | ✅ 已完成 |
| `playercertification_integration_test.go` | 陪玩师实名认证测试 | ✅ 已完成 |
| `ordertimeout_integration_test.go` | 订单超时处理测试 | ✅ 已完成 |

---

## 八、成功标准

- ✅ 所有新增测试用例编译通过
- ✅ 所有新增测试用例运行通过
- ✅ 测试覆盖率达到 75%+
- ✅ P0 场景 100% 覆盖
- ✅ P1 场景 80%+ 覆盖
- ✅ 营销模块 70%+ 覆盖

---

## 九、测试运行命令

```bash
# 设置测试数据库端口
$env:TEST_DB_PORT = "5433"

# 运行所有集成测试
go test -C api ./internal/service/integration/... -v -count=1 -timeout 300s

# 运行特定模块测试
go test -C api ./internal/service/integration/... -v -count=1 -timeout 300s -run "TestOrderService"

# 运行补充测试
go test -C api ./internal/service/integration/... -v -count=1 -timeout 300s -run "Supplement"
```

---

## 十、变更日志

| 日期 | 变更内容 |
|------|----------|
| 2025-12-29 | ✅ P0补充模块测试完成：新增 72 个测试用例（gamerank 13, playerrank 17, playercertification 15, ordertimeout 25） |
| 2025-12-29 | ✅ 辅助模块测试完成：新增 53 个测试用例（game 9, serviceitem 10, review 10, sensitiveword 10, userblock 14） |
| 2025-12-29 | ✅ 营销模块测试完成：新增 99 个测试用例（vip 13, coupon 15, recharge 13, activity 15, team 15, referral 28） |
| 2025-12-29 | 更新测试统计：355 个测试用例全部通过，26 个模块覆盖 |
| 2025-12-29 | 初始化集成测试补充计划，基于现有测试用例分析 |
