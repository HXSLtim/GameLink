# 测试任务单：提现管理模块全量测试

**任务编号**: TEST-2024-M15  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/finance/withdraw | #detailBtn | 详情 | GET /api/v1/admin/withdraws/:id | P0 | ☐ |
| /admin/finance/withdraw | #approveBtn | 批准 | PUT /api/v1/admin/withdraws/:id/approve | P0 | ☐ |
| /admin/finance/withdraw | #rejectBtn | 拒绝 | PUT /api/v1/admin/withdraws/:id/reject | P0 | ☐ |
| /admin/finance/withdraw | #completeBtn | 完成打款 | PUT /api/v1/admin/withdraws/:id/complete | P0 | ☐ |
| /admin/finance/withdraw | #batchApproveBtn | 批量批准 | PUT /api/v1/admin/withdraws/batch/approve | P1 | ☐ |
| /admin/finance/withdraw | #batchRejectBtn | 批量拒绝 | PUT /api/v1/admin/withdraws/batch/reject | P1 | ☐ |
| /admin/finance/withdraw | #batchCompleteBtn | 批量完成打款 | PUT /api/v1/admin/withdraws/batch/complete | P1 | ☐ |
| /admin/finance/withdraw | #exportBtn | 导出数据 | 前端CSV导出 | P2 | ☐ |
| /admin/finance/withdraw | #searchBtn | 搜索 | GET /api/v1/admin/withdraws?playerId=xxx | P1 | ☐ |
| /admin/finance/withdraw | #statusFilter | 状态筛选 | GET /api/v1/admin/withdraws?status=xxx | P1 | ☐ |
| /admin/finance/withdraw | #refreshBtn | 刷新 | GET /api/v1/admin/withdraws | P1 | ☐ |

**重要**: 以上11个按钮，必须全部测试完成，少一个 = 任务未完成

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

-- 查看提现统计
SELECT status, COUNT(*), SUM(amount_cents) as total_amount FROM withdraws GROUP BY status;

-- 查看提现详情
SELECT w.id, w.amount_cents, w.status, w.bank_name, w.bank_account,
       p.nickname as player_name
FROM withdraws w
LEFT JOIN players p ON w.player_id = p.id
LIMIT 10;
```

### 提现状态说明
- `pending`: 待审核
- `approved`: 已批准（待打款）
- `rejected`: 已拒绝
- `completed`: 已完成

---

## 四、逐个按钮测试记录

### 按钮1: #detailBtn 详情

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 点击某提现行的"详情"按钮
4. 观察详情弹窗

**Evidence收集**:
- [ ] 截图1: 详情按钮
- [ ] 截图2: Network请求详情
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM withdraws WHERE id = :withdraw_id;
  ```
- [ ] 截图5: 详情弹窗显示完整信息（陪玩师、金额、银行信息、状态）

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮2: #approveBtn 批准

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击待审核提现的"批准"按钮
4. 确认批准

**Evidence收集**:
- [ ] 截图1: 批准确认弹窗
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/withdraws/:id/approve
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, status, processed_at FROM withdraws WHERE id = :withdraw_id;
  ```
- [ ] 截图5: 列表状态变为"已批准"

**异常场景测试**:
- [ ] 场景A: 批准已批准的提现 → 预期: 按钮不显示
- [ ] 场景B: 批准已完成的提现 → 预期: 按钮不显示
- [ ] 场景C: 取消批准 → 预期: 无请求发送

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮3: #rejectBtn 拒绝

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击待审核提现的"拒绝"按钮
4. 填写拒绝原因
5. 确认拒绝

**Evidence收集**:
- [ ] 截图1: 拒绝弹窗
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/withdraws/:id/reject
  - Payload: `{"reason":"银行信息有误"}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, status, reject_reason FROM withdraws WHERE id = :withdraw_id;
  -- 验证用户钱包余额恢复
  SELECT * FROM wallets WHERE user_id = :player_user_id;
  ```
- [ ] 截图5: 列表状态变为"已拒绝"

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮4: #completeBtn 完成打款

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击已批准提现的"完成打款"按钮
4. 确认完成

**Evidence收集**:
- [ ] 截图1: 完成确认弹窗
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/withdraws/:id/complete
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, status, completed_at FROM withdraws WHERE id = :withdraw_id;
  ```
- [ ] 截图5: 列表状态变为"已完成"

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮5-7: 批量操作

**批量批准 #batchApproveBtn**:
- [ ] 选中多个待审核提现或选择全部
- [ ] 验证: 所有选中提现状态变为"已批准"

**批量拒绝 #batchRejectBtn**:
- [ ] 选中多个待审核提现或选择全部
- [ ] 填写拒绝原因
- [ ] 验证: 所有选中提现状态变为"已拒绝"

**批量完成打款 #batchCompleteBtn**:
- [ ] 选中多个已批准提现或选择全部
- [ ] 验证: 所有选中提现状态变为"已完成"

---

### 按钮8-11: 导出/搜索/筛选/刷新

**导出 #exportBtn**:
- [ ] 点击导出数据
- [ ] 验证CSV文件下载
- [ ] 验证数据完整性

**搜索 #searchBtn**:
- [ ] 输入陪玩师ID搜索
- [ ] 验证: GET /api/v1/admin/withdraws?playerId=xxx

**状态筛选 #statusFilter**:
- [ ] 选择特定状态
- [ ] 验证: GET /api/v1/admin/withdraws?status=pending

**刷新 #refreshBtn**:
- [ ] 点击刷新按钮
- [ ] 验证: GET /api/v1/admin/withdraws

---

## 五、全量测试完整性自查

- [ ] 所有P0按钮已测试（4个）
- [ ] 所有P1按钮已测试（6个）
- [ ] 所有P2按钮已测试（1个）
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
