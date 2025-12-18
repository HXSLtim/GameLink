# 测试任务单：佣金管理模块全量测试

**任务编号**: TEST-2024-M16  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/finance/commission | #createRuleBtn | 新增规则 | POST /api/v1/admin/commission/rules | P0 | ☐ |
| /admin/finance/commission | #editRuleBtn | 编辑 | PUT /api/v1/admin/commission/rules/:id | P0 | ☐ |
| /admin/finance/commission | #triggerSettlementBtn | 触发月度结算 | POST /api/v1/admin/settlement/trigger | P0 | ☐ |
| /admin/finance/commission | #monthPicker | 月份选择 | GET /api/v1/admin/stats/platform?month=xxx | P1 | ☐ |
| /admin/finance/commission | #refreshBtn | 刷新 | GET /api/v1/admin/stats/platform | P1 | ☐ |

**重要**: 以上5个按钮，必须全部测试完成，少一个 = 任务未完成

---

## 二、Docker环境检查

```bash
docker compose -f docker-compose.prod.yml ps
```

---

## 三、测试数据准备

### 数据库种子数据验证
```sql
-- 连接数据库
docker exec -it gamelink-postgres psql -U gamelink -d gamelink

-- 查看佣金规则
SELECT id, name, rate_percent, min_amount_cents, max_amount_cents, is_default 
FROM commission_rules;

-- 查看结算记录
SELECT id, month, total_revenue_cents, total_commission_cents, status, created_at 
FROM settlements ORDER BY month DESC LIMIT 5;

-- 查看平台统计
SELECT 
  COUNT(*) as total_orders,
  SUM(total_price_cents) as total_revenue,
  SUM(CASE WHEN status = 'completed' THEN total_price_cents ELSE 0 END) as completed_revenue
FROM orders;
```

---

## 四、逐个按钮测试记录

### 按钮1: #createRuleBtn 新增规则

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 点击"新增规则"按钮
4. 填写表单：
   - 规则名称: 测试抽成规则
   - 抽成比例: 25%
   - 最低订单金额: 50元
   - 最高订单金额: 1000元
   - 是否默认: 否
5. 点击保存

**Evidence收集**:
- [ ] 截图1: 新增规则弹窗表单
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/commission/rules
  - Payload: `{"name":"测试抽成规则","ratePercent":25,...}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM commission_rules WHERE name = '测试抽成规则';
  ```
- [ ] 截图5: 规则列表显示新规则

**异常场景测试**:
- [ ] 场景A: 名称为空 → 预期: 校验提示
- [ ] 场景B: 抽成比例超过100% → 预期: 校验提示
- [ ] 场景C: 最低金额大于最高金额 → 预期: 校验提示

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮2: #editRuleBtn 编辑规则

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某规则的"编辑"按钮
4. 修改抽成比例
5. 点击保存

**Evidence收集**:
- [ ] 截图1: 编辑弹窗（含原数据回显）
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/commission/rules/:id
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM commission_rules WHERE id = :rule_id;
  ```
- [ ] 截图5: 规则列表显示更新后的数据

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮3: #triggerSettlementBtn 触发月度结算

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 选择要结算的月份
4. 点击"触发月度结算"按钮
5. 确认操作

**Evidence收集**:
- [ ] 截图1: 结算确认弹窗
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/settlement/trigger
  - Payload: `{"month":"2024-12"}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  -- 验证结算记录
  SELECT * FROM settlements WHERE month = '2024-12';
  -- 验证陪玩师收益
  SELECT p.nickname, pe.amount_cents, pe.commission_cents 
  FROM player_earnings pe
  JOIN players p ON pe.player_id = p.id
  WHERE pe.month = '2024-12';
  ```
- [ ] 截图5: 统计数据更新

**异常场景测试**:
- [ ] 场景A: 重复触发同月结算 → 预期: 提示已结算或覆盖
- [ ] 场景B: 结算未来月份 → 预期: 拒绝或警告
- [ ] 场景C: 无订单月份结算 → 预期: 正常完成，金额为0

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮4-5: 月份选择/刷新

**月份选择 #monthPicker**:
- [ ] 选择不同月份
- [ ] 验证: GET /api/v1/admin/stats/platform?month=2024-11
- [ ] 验证统计数据正确切换

**刷新 #refreshBtn**:
- [ ] 点击刷新按钮
- [ ] 验证: GET /api/v1/admin/stats/platform
- [ ] 验证数据与数据库一致

---

## 五、全量测试完整性自查

- [ ] 所有P0按钮已测试（3个）
- [ ] 所有P1按钮已测试（2个）
- [ ] 每个按钮提供5张截图+日志
- [ ] 每个按钮测试异常场景≥3个

---

## 六、质量承诺

我承诺以上测试内容真实完整，所有按钮均已按22项清单验证。

**测试人签字**: ___________  
**日期**: ___________

---

## 七、组长审核意见

**审核结果**: ☐ 通过 ☐ 打回重做  
**审核人**: ___________  
**日期**: ___________

---

**文档版本**: v1.0  
**发布日期**: 2024-12-18
