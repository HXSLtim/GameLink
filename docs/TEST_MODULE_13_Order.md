# 测试任务单：订单管理模块全量测试

**任务编号**: TEST-2024-M13  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/biz/order | #detailBtn | 详情 | GET /api/v1/admin/orders/:id | P0 | ☐ |
| /admin/biz/order | #cancelBtn | 取消 | PUT /api/v1/admin/orders/:id/cancel | P0 | ☐ |
| /admin/biz/order | #refundBtn | 退款 | POST /api/v1/admin/orders/:id/refund | P0 | ☐ |
| /admin/biz/order | #batchCancelBtn | 批量取消 | PUT /api/v1/admin/orders/batch/cancel | P1 | ☐ |
| /admin/biz/order | #batchCompleteBtn | 批量完成 | PUT /api/v1/admin/orders/batch/complete | P1 | ☐ |
| /admin/biz/order | #exportBtn | 导出数据 | 前端CSV导出 | P2 | ☐ |
| /admin/biz/order | #searchBtn | 搜索 | GET /api/v1/admin/orders?orderNo=xxx | P1 | ☐ |
| /admin/biz/order | #statusFilter | 状态筛选 | GET /api/v1/admin/orders?status=xxx | P1 | ☐ |
| /admin/biz/order | #refreshBtn | 刷新 | GET /api/v1/admin/orders | P1 | ☐ |

**重要**: 以上9个按钮，必须全部测试完成，少一个 = 任务未完成

---

## 二、测试标准（参考22项清单）

每个按钮必须提供：
1. ✅ 按钮静态截图（Evidence-01）
2. ✅ Network请求截图（Evidence-02）
3. ✅ docker logs截图（Evidence-03）
4. ✅ 数据库验证SQL截图（Evidence-04）
5. ✅ 异常场景测试结果（Evidence-05）
6. ✅ 完整操作录像（asciinema或录屏）

---

## 三、Docker环境检查

```bash
docker compose -f docker-compose.prod.yml ps
```

---

## 四、测试数据准备

### 数据库种子数据验证
```sql
-- 连接数据库
docker exec -it gamelink-postgres psql -U gamelink -d gamelink

-- 查看订单统计
SELECT status, COUNT(*) FROM orders GROUP BY status;

-- 查看订单详情
SELECT o.id, o.order_no, o.status, o.total_price_cents, 
       u.name as user_name, p.nickname as player_name, g.name as game_name
FROM orders o
LEFT JOIN users u ON o.user_id = u.id
LEFT JOIN players p ON o.player_id = p.id
LEFT JOIN games g ON o.game_id = g.id
LIMIT 10;
```

### 订单状态说明
- `pending`: 待确认
- `confirmed`: 已确认
- `in_progress`: 进行中
- `completed`: 已完成
- `cancelled`: 已取消
- `refunded`: 已退款

---

## 五、逐个按钮测试记录

### 按钮1: #detailBtn 详情

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 点击某订单行的"详情"按钮
4. 观察详情抽屉

**Evidence收集**:
- [ ] 截图1: 详情按钮
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/orders/:id
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM orders WHERE id = :order_id;
  ```
- [ ] 截图5: 详情抽屉显示完整信息（订单信息、用户信息、陪玩师信息、订单进度）

**异常场景测试**:
- [ ] 场景A: 查看不存在的订单 → 预期: 404错误
- [ ] 场景B: 关闭后重新打开 → 预期: 重新加载数据
- [ ] 场景C: 网络中断时查看 → 预期: 错误提示

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮2: #cancelBtn 取消订单

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击待确认/已确认订单的"取消"按钮
4. 确认取消

**Evidence收集**:
- [ ] 截图1: 取消确认弹窗
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/orders/:id/cancel
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, order_no, status, cancel_reason, cancelled_at 
  FROM orders WHERE id = :order_id;
  ```
- [ ] 截图5: 列表状态变为"已取消"

**异常场景测试**:
- [ ] 场景A: 取消已完成订单 → 预期: 按钮不显示
- [ ] 场景B: 取消已取消订单 → 预期: 按钮不显示
- [ ] 场景C: 取消进行中订单 → 预期: 按钮不显示或警告

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮3: #refundBtn 退款

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某订单的"退款"按钮
4. 填写退款金额和原因
5. 确认退款

**Evidence收集**:
- [ ] 截图1: 退款弹窗（含订单信息概览）
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/orders/:id/refund
  - Payload: `{"amount_cents":3000,"reason":"用户申请退款"}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  -- 验证订单状态
  SELECT id, status FROM orders WHERE id = :order_id;
  -- 验证退款记录
  SELECT * FROM refunds WHERE order_id = :order_id;
  -- 验证用户钱包
  SELECT * FROM wallets WHERE user_id = :user_id;
  ```
- [ ] 截图5: 订单状态变为"已退款"

**异常场景测试**:
- [ ] 场景A: 退款金额超过订单金额 → 预期: 校验提示
- [ ] 场景B: 退款金额为0 → 预期: 校验提示
- [ ] 场景C: 重复退款 → 预期: 拒绝或提示

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮4: #batchCancelBtn 批量取消

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击"批量取消"按钮
4. 选择目标（选中的订单/按状态/全部可取消）
5. 填写取消原因
6. 确认取消

**Evidence收集**:
- [ ] 截图1: 批量取消弹窗
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/orders/batch/cancel
  - Payload: `{"orderIds":[1,2,3],"reason":"批量取消"}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, status FROM orders WHERE id IN (1,2,3);
  ```
- [ ] 截图5: 所有选中订单状态变为"已取消"

**异常场景测试**:
- [ ] 场景A: 未选中任何订单 → 预期: 提示选择
- [ ] 场景B: 包含不可取消订单 → 预期: 跳过或提示
- [ ] 场景C: 按状态筛选取消 → 预期: 正确筛选并取消

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮5: #batchCompleteBtn 批量完成

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击"批量完成"按钮
4. 选择目标
5. 确认完成

**Evidence收集**:
- [ ] 截图1: 批量完成弹窗
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/orders/batch/complete
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, status, completed_at FROM orders WHERE id IN (:ids);
  ```
- [ ] 截图5: 所有选中订单状态变为"已完成"

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮6-9: 导出/搜索/筛选/刷新

**导出 #exportBtn**:
- [ ] 点击导出数据
- [ ] 验证CSV文件下载
- [ ] 验证数据完整性（订单号、用户、陪玩师、金额、状态、时间）

**搜索 #searchBtn**:
- [ ] 输入订单号搜索
- [ ] 验证: GET /api/v1/admin/orders?orderNo=xxx

**状态筛选 #statusFilter**:
- [ ] 选择特定状态
- [ ] 验证: GET /api/v1/admin/orders?status=pending

**刷新 #refreshBtn**:
- [ ] 点击刷新按钮
- [ ] 验证: GET /api/v1/admin/orders

---

## 六、全量测试完整性自查

- [ ] 所有P0按钮已测试（3个）
- [ ] 所有P1按钮已测试（5个）
- [ ] 所有P2按钮已测试（1个）
- [ ] 每个按钮提供5张截图+日志
- [ ] 每个按钮测试异常场景≥3个

---

## 七、质量承诺

我承诺以上测试内容真实完整，所有按钮均已按22项清单验证。

**测试人签字**: ___________  
**日期**: ___________

---

## 八、组长审核意见

**审核结果**: ☐ 通过 ☐ 打回重做  
**审核人**: ___________  
**日期**: ___________

---

**文档版本**: v1.0  
**发布日期**: 2024-12-18
