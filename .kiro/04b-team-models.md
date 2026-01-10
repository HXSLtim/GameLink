# 数据模型 - 团队模块

> 陪玩师团队系统

## 相关文档

- [04-data-models.md](./04-data-models.md) - 核心模块
- [04a-marketing-models.md](./04a-marketing-models.md) - 营销模块
- [04c-enums-indexes.md](./04c-enums-indexes.md) - 枚举类型和数据库索引
- [04d-notification-models.md](./04d-notification-models.md) - 通知系统

---

## 团队模块

### Team（陪玩师团队）

| 字段 | 类型 | 说明 |
|------|------|------|
| Name | string | 团队名称 |
| Description | string | 团队简介 |
| AvatarURL | string | 团队头像 |
| LeaderID | uint64 | 队长ID（Player.ID） |
| Status | TeamStatus | 状态：active/busy/inactive |
| MaxMembers | int | 最大成员数（默认5） |
| MemberCount | int | 当前成员数 |
| IncomeShareType | string | 收入分配方式：equal（平分）/custom（自定义） |
| LeaderBonusRate | float64 | 队长额外抽成比例（预留） |
| TotalOrderCount | int | 累计接单数 |
| TotalIncomeCents | int64 | 累计收入（分） |
| CurrentOrderID | *uint64 | 当前进行中的订单ID |

### TeamMember（团队成员）

| 字段 | 类型 | 说明 |
|------|------|------|
| TeamID | uint64 | 团队ID |
| PlayerID | uint64 | 陪玩师ID |
| Role | TeamMemberRole | 角色：leader/member |
| Status | TeamMemberStatus | 状态：active/left/kicked |
| JoinedAt | time.Time | 加入时间 |
| LeftAt | *time.Time | 离队时间 |
| SortOrder | int | 排序（队长继承顺序） |
| OrderCount | int | 在团队内接单数 |
| IncomeCents | int64 | 在团队内收入（分） |

### TeamInvite（团队邀请）- 预留

| 字段 | 类型 | 说明 |
|------|------|------|
| TeamID | uint64 | 团队ID |
| PlayerID | uint64 | 被邀请的陪玩师ID |
| InviterID | uint64 | 邀请人ID |
| Status | TeamInviteStatus | 状态：pending/accepted/rejected/expired |
| ExpireAt | time.Time | 过期时间 |
| Message | string | 邀请留言 |

---

## 团队业务流程

### 创建团队

```
陪玩师创建团队 → 自动成为队长
  ├── Team.LeaderID = 创建者ID
  ├── Team.MemberCount = 1
  └── TeamMember (Role=leader, SortOrder=0)
```

### 加入团队

```
陪玩师申请/被邀请 → 队长同意
  ├── 检查：陪玩师未加入其他团队
  ├── 创建 TeamMember (Role=member)
  └── Team.MemberCount++
```

### 退出团队

```
成员主动退出 / 被队长踢出
  ├── 检查：Team.CurrentOrderID == nil（无进行中订单）
  ├── TeamMember.Status = left/kicked
  ├── TeamMember.LeftAt = now
  └── Team.MemberCount--
```

### 队长离队

```
队长退出 → 顺位继承
  ├── 按 SortOrder 找下一个 active 成员
  ├── 新队长 Role = leader
  └── Team.LeaderID = 新队长ID
```

### 最后一人离队

```
MemberCount == 0 → 团队自动销毁
  └── Team 软删除
```

### 团队接单

```
【团队接单规则】
- 团队人数必须等于订单座位数（RequiredPlayers）
- 团队人数不足或超出时，不能接单
- 团队接单是整体行为，不支持部分成员接单

【人数匹配检查】
团队人数 == 订单座位数 → 可以接单
团队人数 != 订单座位数 → 拒绝接单，提示"团队人数不匹配"

【接单流程】
团队接受订单 → 所有成员进入忙碌状态
  ├── Team.Status = busy
  ├── Team.CurrentOrderID = 订单ID
  ├── Order.Status = in_progress（订单状态变更）
  ├── Order.CurrentPlayers = RequiredPlayers
  └── 为每个成员创建 OrderPlayer 记录
```

### 团队接单定价规则

【服务类型与定价】
| 类型 | 定价方式 | 团队接单 |
|------|----------|----------|
| 陪玩 | 按段位定价（每座位价格可不同） | 成员按段位分配到对应座位 |
| 护航 | 按小时/场次定价（统一价格） | 团队统一价 × 人数 |

【陪玩订单 - 团队接单】
```
- 订单座位有不同段位要求和价格
- 团队成员按 SortOrder 分配到各座位
- 成员段位必须 >= 座位要求段位
- 收入按各自座位价格计算

示例：
订单：王者座位¥50 + 钻石座位¥30
团队：A(王者段位, SortOrder=0), B(钻石段位, SortOrder=1)
分配：A→王者座位, B→钻石座位
收入：A=¥50×80%, B=¥30×80%
```

【护航订单 - 团队接单】
```
- 所有座位统一价格
- 团队人数 = 订单座位数即可接单
- 收入平分或队长自定义分配

示例：
订单：3人护航，¥30/人/小时，2小时
总价：¥30 × 3 × 2 = ¥180
团队收入：¥180 × 80% = ¥144
平分：每人¥48
```

### 订单完成 - 收入分配

```
【分配时机】
订单完成时立即分配（分配的是冻结金额）

【分配流程】
订单完成 → 计算总收入（扣除平台抽成后）
  → 队长选择分配方式
  → 为每个成员创建 CommissionRecord
  → 收入进入各成员 Wallet.FrozenCents
  → T+7 后各自解冻到 BalanceCents

【分配方式】
  ├── IncomeShareType = equal：总收入 ÷ 成员数
  ├── IncomeShareType = custom：队长指定每人金额或座位归属
  └── 更新 OrderPlayer.IncomeCents
```

### 团队状态恢复

```
订单完成/取消 → 恢复空闲
  ├── Team.Status = active
  └── Team.CurrentOrderID = nil
```

---

## 个人陪玩 vs 团队陪玩

| 对比项 | 个人陪玩 | 团队陪玩 |
|--------|----------|----------|
| 下单方式 | 选择多个不同等级陪玩师 | 选择一个团队 |
| 座位价格 | 可以不同（王者¥50 + 钻石¥30） | 陪玩：按座位段位定价；护航：统一价 |
| 收入分配 | 各自按座位价格 | 陪玩：按座位价格；护航：平分或自定义 |
| 接单主体 | 单个陪玩师 | 整个团队 |

---

## 团队接单代码示例

```go
// 团队整体接单，一次性填满所有座位
func TeamJoinOrder(orderID, teamID uint64) error {
    // 获取团队成员
    members := getTeamMembers(teamID)
    
    return db.Transaction(func(tx *gorm.DB) error {
        // 获取所有待匹配的 OrderItem
        var items []OrderItem
        tx.Where("order_id = ? AND status = ?", orderID, "pending").Find(&items)
        
        // 检查人数是否匹配
        if len(members) != len(items) {
            return errors.New("团队人数不匹配")
        }
        
        // 批量分配
        for i, item := range items {
            item.PlayerID = &members[i].ID
            item.Status = "matched"
            tx.Save(&item)
            tx.Create(&OrderPlayer{
                OrderID:     orderID,
                OrderItemID: item.ID,
                PlayerID:    members[i].ID,
                TeamID:      &teamID,
                JoinedAt:    time.Now(),
                Status:      "joined",
            })
        }
        
        // 更新主订单
        tx.Model(&Order{}).Where("id = ?", orderID).Updates(map[string]interface{}{
            "current_players": len(items),
            "status":          "in_progress",
        })
        
        // 更新团队状态
        tx.Model(&Team{}).Where("id = ?", teamID).Updates(map[string]interface{}{
            "status":           "busy",
            "current_order_id": orderID,
        })
        
        return nil
    })
}
```
