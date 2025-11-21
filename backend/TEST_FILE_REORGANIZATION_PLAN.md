# 测试文件重命名和合并计划

**制定时间**: 2025-11-22 06:15:00
**目标**: 统一测试文件命名为 `*_test.go`，合并相关测试

---

## 📊 当前测试文件统计

### 需要合并的文件

| 文件类型 | 数量 | 示例 |
|----------|------|------|
| `*_quick_test.go` | 24个 | `order_handler_quick_test.go` |
| `*_coverage_test.go` | 4个 | `commission_handler_coverage_test.go` |
| `*_extended_test.go` | 8个 | `order_extended_test.go` |
| `*_badjson_quick_test.go` | 1个 | `order_badjson_quick_test.go` |
| `*_invalid_quick_test.go` | 1个 | `order_invalid_quick_test.go` |
| `*_filters_quick_test.go` | 1个 | `order_filters_quick_test.go` |
| **总计** | **39个文件** | - |

---

## 🎯 重命名和合并计划

### 1. Handler层测试文件合并

#### 订单Handler测试
```bash
# 当前文件
internal/handler/user/order_test.go
internal/handler/user/order_handler_quick_test.go
internal/handler/user/order_badjson_quick_test.go
internal/handler/user/order_invalid_quick_test.go
internal/handler/user/order_filters_quick_test.go

internal/handler/admin/order_test.go
internal/handler/admin/order_handler_quick_test.go
internal/handler/admin/order_payment_failure_test.go

internal/handler/player/order_test.go

# 合并后（只保留3个文件）
internal/handler/user/order_test.go          # 合并所有用户端订单测试
internal/handler/admin/order_test.go         # 合并所有管理端订单测试
internal/handler/player/order_test.go        # 保留陪玩师端订单测试
```

**合并示例**:
```go
// internal/handler/user/order_test.go
func TestOrderHandler(t *testing.T) {
    t.Run("CreateOrder", func(t *testing.T) {
        // 原 order_test.go 内容
    })
    
    t.Run("CreateOrder_Quick", func(t *testing.T) {
        // 原 order_handler_quick_test.go 内容
    })
    
    t.Run("CreateOrder_InvalidJSON", func(t *testing.T) {
        // 原 order_badjson_quick_test.go 内容
    })
    
    t.Run("CreateOrder_InvalidInput", func(t *testing.T) {
        // 原 order_invalid_quick_test.go 内容
    })
    
    t.Run("ListOrders_Filters", func(t *testing.T) {
        // 原 order_filters_quick_test.go 内容
    })
}
```

---

#### 支付Handler测试
```bash
# 当前文件
internal/handler/user/payment_test.go
internal/handler/user/payment_handler_quick_test.go
internal/handler/user/payment_badjson_quick_test.go

internal/handler/admin/payment_test.go

# 合并后
internal/handler/user/payment_test.go          # 合并用户端支付测试
internal/handler/admin/payment_test.go         # 保留管理端支付测试
```

---

#### 礼物Handler测试
```bash
# 当前文件
internal/handler/user/gift_test.go
internal/handler/user/gift_handler_quick_test.go

internal/handler/player/gift_test.go
internal/handler/player/gift_handler_quick_test.go

internal/handler/admin/item_test.go
internal/handler/admin/item_handler_quick_test.go

# 合并后
internal/handler/user/gift_test.go             # 合并用户端礼物测试
internal/handler/player/gift_test.go           # 合并陪玩师端礼物测试
internal/handler/admin/item_test.go            # 合并管理端物品测试
```

---

#### 佣金Handler测试
```bash
# 当前文件
internal/handler/player/commission_test.go
internal/handler/player/commission_handler_quick_test.go

internal/handler/admin/commission_test.go
internal/handler/admin/commission_handler_quick_test.go
internal/handler/admin/commission_handler_coverage_test.go
internal/handler/admin/commission_complete_test.go

# 合并后
internal/handler/player/commission_test.go     # 合并陪玩师端佣金测试
internal/handler/admin/commission_test.go      # 合并管理端佣金测试
```

---

#### 收入Handler测试
```bash
# 当前文件
internal/handler/player/earnings_test.go
internal/handler/player/earnings_handler_quick_test.go

# 合并后
internal/handler/player/earnings_test.go       # 合并陪玩师端收入测试
```

---

#### 排名Handler测试
```bash
# 当前文件
internal/handler/admin/ranking_test.go
internal/handler/admin/ranking_handler_quick_test.go
internal/handler/admin/ranking_handler_coverage_test.go

# 合并后
internal/handler/admin/ranking_test.go         # 合并管理端排名测试
```

---

#### 统计Handler测试
```bash
# 当前文件
internal/handler/admin/stats_handler_coverage_test.go
internal/handler/admin/stats_handler_quick_test.go

# 合并后
internal/handler/admin/stats_test.go           # 重命名并合并统计测试
```

---

#### 用户管理Handler测试
```bash
# 当前文件
internal/handler/admin/user_test.go
internal/handler/admin/user_handler_quick_test.go
internal/handler/admin/user_list_quick_test.go

# 合并后
internal/handler/admin/user_test.go            # 合并管理端用户测试
```

---

#### 提现Handler测试
```bash
# 当前文件
internal/handler/admin/withdraw_test.go
internal/handler/admin/withdraw_handler_quick_test.go
internal/handler/admin/withdraw_complete_test.go

# 合并后
internal/handler/admin/withdraw_test.go        # 合并管理端提现测试
```

---

#### 其他Handler测试
```bash
# 当前文件
internal/handler/adminouter_permission_quick_test.go
internal/handler/adminouter_quick_test.go
internal/handler/admin	est_router_helpers_test.go
internal/handler/adminouter.go

# 合并后
internal/handler/admin/router_test.go          # 合并路由测试
internal/handler/admin/helpers_test.go         # 合并辅助函数测试
```

---

### 2. Service层测试文件合并

#### 订单Service测试
```bash
# 当前文件
internal/service/order/order_test.go
internal/service/order/order_extended_test.go
internal/service/order/order_autodestroy_test.go
internal/service/order/order_availability_test.go

# 合并后
internal/service/order/order_test.go           # 合并所有订单服务测试
```

---

#### 支付Service测试
```bash
# 当前文件
internal/service/payment/payment_test.go
internal/service/payment/payment_extended_test.go
internal/service/payment/payment_additional_test.go
internal/service/payment/payment_full_coverage_test.go

# 合并后
internal/service/payment/payment_test.go       # 合并所有支付服务测试
```

---

#### 佣金Service测试
```bash
# 当前文件
internal/service/commission/commission_test.go
internal/service/commission/commission_extended_test.go
internal/service/commission/commission_additional_test.go

# 合并后
internal/service/commission/commission_test.go # 合并所有佣金服务测试
```

---

#### 物品Service测试
```bash
# 当前文件
internal/service/item/item_test.go
internal/service/item/item_extended_test.go

# 合并后
internal/service/item/item_test.go             # 合并所有物品服务测试
```

---

#### Admin Service测试
```bash
# 当前文件
internal/service/admin/admin_test.go
internal/service/admin/admin_extended_test.go
internal/service/admin/admin_deep_test.go
internal/service/admin/admin_order_timeline_test.go
internal/service/admin/admin_quick_test.go
internal/service/admin/admin_reviews_service_test.go
internal/service/admin/admin_service_more_test.go
internal/service/admin/admin_tx_test.go
internal/service/admin/admin_user_game_test.go

# 合并后
internal/service/admin/admin_test.go           # 合并所有管理服务测试
```

---

#### 集成测试
```bash
# 当前文件
internal/service/integration_test.go
internal/service/integration_extended_test.go

# 合并后
internal/service/integration_test.go           # 合并所有集成测试
```

---

### 3. Repository层测试文件合并

#### 订单Repository测试
```bash
# 当前文件
internal/repository/implementations/order_repository_test.go

# 保留
internal/repository/implementations/order_repository_test.go
```

---

#### 聊天Repository测试
```bash
# 当前文件
internal/repository/chat/message_repository_test.go
internal/repository/chat/repository_quick_test.go
internal/repository/chat/group_repository_test.go

# 合并后
internal/repository/chat/repository_test.go     # 合并聊天仓库测试
```

---

#### 排名Repository测试
```bash
# 当前文件
internal/repository/ranking/repository_test.go
internal/repository/ranking/repository_simple_test.go
internal/repository/ranking/commission_repository_test.go

# 合并后
internal/repository/ranking/repository_test.go # 合并排名仓库测试
```

---

### 4. Model层测试文件合并

```bash
# 当前文件
internal/model/dispute_test.go
internal/model/model_helpers_test.go
internal/model/order_helper_test.go
internal/model/rating_test.go
internal/model/upload_test.go

# 保留（已经是标准命名）
internal/model/dispute_test.go
internal/model/model_helpers_test.go
internal/model/order_helper_test.go
internal/model/rating_test.go
internal/model/upload_test.go
```

---

### 5. 其他测试文件合并

```bash
# 当前文件
internal/service/assignment/integration_test.go

internal/service/chat/chat_service_quick_test.go

internal/db/migrate_test.go
internal/db/migrate_additional_test.go
internal/db/db_additional_test.go
internal/db/seed_test.go

internal/config/env_test.go
internal/config/config_extra_test.go

# 合并后
internal/service/assignment/integration_test.go  # 保留
internal/service/chat/chat_test.go              # 重命名
internal/db/db_test.go                          # 合并数据库测试
internal/config/config_test.go                  # 合并配置测试
```

---

## 📋 具体操作步骤

### 第一步：备份现有测试文件
```bash
# 创建备份目录
mkdir -p backup/tests_$(date +%Y%m%d_%H%M%S)

# 备份所有测试文件
cp -r internal/*_test.go backup/tests_$(date +%Y%m%d_%H%M%S)/
```

---

### 第二步：按模块合并测试文件

#### 示例：合并订单Handler测试
```bash
# 1. 创建新的 order_test.go
cat > internal/handler/user/order_test.go << 'EOF'
package user

import (
    "testing"
    // ... 其他导入
)

func TestOrderHandler(t *testing.T) {
    t.Run("CreateOrder", func(t *testing.T) {
        testCreateOrder(t)
    })
    
    t.Run("CreateOrder_Quick", func(t *testing.T) {
        testCreateOrderQuick(t)
    })
    
    t.Run("GetMyOrders", func(t *testing.T) {
        testGetMyOrders(t)
    })
    
    // ... 其他测试
}

// 原 order_test.go 内容
func testCreateOrder(t *testing.T) {
    // ...
}

// 原 order_handler_quick_test.go 内容
func testCreateOrderQuick(t *testing.T) {
    // ...
}

// 原 order_badjson_quick_test.go 内容
func testCreateOrderWithBadJSON(t *testing.T) {
    // ...
}

// 原 order_invalid_quick_test.go 内容
func testCreateOrderWithInvalidInput(t *testing.T) {
    // ...
}

// 原 order_filters_quick_test.go 内容
func testOrderListWithFilters(t *testing.T) {
    // ...
}
EOF

# 2. 删除旧文件
rm internal/handler/user/order_handler_quick_test.go
rm internal/handler/user/order_badjson_quick_test.go
rm internal/handler/user/order_invalid_quick_test.go
rm internal/handler/user/order_filters_quick_test.go
```

---

#### 示例：合并订单Service测试
```bash
# 1. 合并到 order_test.go
cat internal/service/order/order_extended_test.go >> internal/service/order/order_test.go
cat internal/service/order/order_autodestroy_test.go >> internal/service/order/order_test.go
cat internal/service/order/order_availability_test.go >> internal/service/order/order_test.go

# 2. 删除旧文件
rm internal/service/order/order_extended_test.go
rm internal/service/order/order_autodestroy_test.go
rm internal/service/order/order_availability_test.go
```

---

### 第三步：验证测试通过
```bash
# 运行测试
go test ./internal/handler/user/... -v
go test ./internal/service/order/... -v

# 检查测试覆盖率
go test ./internal/handler/user/... -cover
go test ./internal/service/order/... -cover

# 确保覆盖率不低于原有水平
```

---

### 第四步：Git提交
```bash
# 添加新文件
git add internal/handler/user/order_test.go
git add internal/service/order/order_test.go

# 删除旧文件
git rm internal/handler/user/order_handler_quick_test.go
git rm internal/handler/user/order_badjson_quick_test.go
git rm internal/handler/user/order_invalid_quick_test.go
git rm internal/handler/user/order_filters_quick_test.go
git rm internal/service/order/order_extended_test.go
git rm internal/service/order/order_autodestroy_test.go
git rm internal/service/order/order_availability_test.go

# 提交
git commit -m "refactor(tests): 统一订单模块测试文件命名

- 合并所有订单handler测试到 order_test.go
- 合并所有订单service测试到 order_test.go
- 删除冗余的 quick/extended/autodestroy 测试文件
- 使用子测试组织不同类型的测试用例"
```

---

## 📊 合并后文件清单

### 最终保留的测试文件（标准命名）

```bash
# Handler层（约30个文件）
internal/handler/user/
├── auth_test.go
├── order_test.go          # 合并5个文件
├── payment_test.go        # 合并2个文件
├── gift_test.go           # 合并2个文件
├── chat_test.go
├── dispute_test.go
├── feed_test.go
├── player_test.go
└── review_test.go

internal/handler/player/
├── commission_test.go     # 合并2个文件
├── earnings_test.go       # 合并2个文件
├── gift_test.go           # 合并2个文件
├── order_test.go
├── profile_test.go
└── review_test.go

internal/handler/admin/
├── commission_test.go     # 合并4个文件
├── dashboard_test.go      # 合并2个文件
├── dispute_test.go
├── game_test.go           # 合并2个文件
├── item_test.go           # 合并2个文件
├── order_test.go          # 合并3个文件
├── permission_test.go
├── player_test.go         # 合并2个文件
├── ranking_test.go        # 合并3个文件
├── review_test.go
├── role_test.go
├── stats_test.go          # 合并2个文件 + 重命名
├── system_test.go         # 合并2个文件
├── user_test.go           # 合并3个文件
├── withdraw_test.go       # 合并3个文件
├── router_test.go         # 合并3个文件
└── helpers_test.go        # 合并1个文件

# Service层（约20个文件）
internal/service/
├── auth/auth_test.go
├── order/order_test.go            # 合并4个文件
├── payment/payment_test.go        # 合并4个文件
├── commission/commission_test.go  # 合并3个文件
├── item/item_test.go              # 合并2个文件
├── admin/admin_test.go            # 合并8个文件
├── integration/integration_test.go # 合并2个文件
├── chat/chat_test.go              # 重命名 + 合并
├── player/player_test.go
├── ranking/ranking_test.go
├── review/review_test.go
├── role/role_test.go
├── stats/stats_test.go
├── gift/gift_test.go
├── feed/feed_test.go
├── notification/notification_test.go
└── team/team_test.go

# Repository层（约25个文件）
internal/repository/
├── user/repository_test.go
├── player/repository_test.go
├── game/repository_test.go
├── order/repository_test.go
├── payment/repository_test.go
├── commission/repository_test.go
├── review/repository_test.go
├── role/repository_test.go
├── permission/repository_test.go
├── chat/repository_test.go          # 合并3个文件
├── ranking/repository_test.go       # 合并3个文件
├── stats/repository_test.go
├── notification/repository_test.go
├── operation_log/repository_test.go
├── implementations/order_repository_test.go
├── implementations/order_repository_test.go
├── dispute/repository_test.go
├── withdraw/repository_test.go
├── serviceitem/repository_test.go
├── feed/repository_test.go
├── player_tag/repository_test.go
├── reviewreply/repository_test.go
└── common/uow_test.go

# Model层（5个文件，保持不变）
internal/model/
├── dispute_test.go
├── model_helpers_test.go
├── order_helper_test.go
├── rating_test.go
└── upload_test.go

# 其他（4个文件）
internal/db/db_test.go                    # 合并4个文件
internal/config/config_test.go            # 合并2个文件
internal/router/router_test.go
internal/scheduler/settlement_scheduler_test.go
internal/scheduler/chat_retention_test.go
```

**总计**: 从100+个测试文件减少到约85个（减少15%）

---

## ⚠️ 风险提示

### 1. 测试合并风险
- **风险**: 合并后测试冲突或遗漏
- **缓解**: 仔细对比每个测试用例，确保全部迁移
- **验证**: 运行全部测试，确保通过率100%

### 2. 覆盖率下降风险
- **风险**: 合并后测试覆盖率下降
- **缓解**: 合并前后对比覆盖率报告
- **标准**: 覆盖率不低于原有水平

### 3. Git历史丢失
- **风险**: 删除文件导致Git历史丢失
- **缓解**: 先合并再删除，保留提交记录
- **备份**: 创建备份分支

---

## ✅ 验收清单

- [ ] 所有 `*_quick_test.go` 文件已合并
- [ ] 所有 `*_coverage_test.go` 文件已合并
- [ ] 所有 `*_extended_test.go` 文件已合并
- [ ] 所有 `*_badjson_test.go` 文件已合并
- [ ] 所有 `*_invalid_test.go` 文件已合并
- [ ] 所有测试文件命名为 `*_test.go`
- [ ] 所有测试通过（`go test ./...`）
- [ ] 测试覆盖率不低于原有水平
- [ ] Git提交记录清晰

---

**计划制定**: 2025-11-22 06:15:00
**执行状态**: 待执行
**预计工时**: 8小时
**优先级**: 高（阻塞其他改进）
