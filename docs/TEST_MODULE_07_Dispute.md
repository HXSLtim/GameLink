# 测试任务单：纠纷管理模块全量测试

**任务编号**: TEST-2024-M07  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/biz/dispute | #statsLoad | 统计加载 | GET /api/v1/admin/disputes/stats | P0 | ☐ |
| /admin/biz/dispute | #searchBtn | 搜索(订单号) | GET /api/v1/admin/disputes | P0 | ☐ |
| /admin/biz/dispute | #filterBtn | 状态筛选 | GET /api/v1/admin/disputes?status=X | P0 | ☐ |
| /admin/biz/dispute | #detailBtn | 详情 | 前端抽屉展示 | P0 | ☐ |
| /admin/biz/dispute | #assignBtn | 分配 | POST /api/v1/admin/disputes/:id/assign | P0 | ☐ |
| /admin/biz/dispute | #resolveBtn | 处理 | POST /api/v1/admin/disputes/:id/resolve | P0 | ☐ |
| /admin/biz/dispute | #rollbackBtn | 回滚 | POST /api/v1/admin/disputes/:id/rollback | P1 | ☐ |
| /admin/biz/dispute | #exportBtn | 导出数据 | 前端CSV导出 | P1 | ☐ |
| /admin/biz/dispute | #refreshBtn | 刷新 | GET /api/v1/admin/disputes | P1 | ☐ |

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

## 三、Docker环境检查（执行后贴结果）

```bash
docker compose -f docker-compose.prod.yml ps
```

**预期结果**: 所有容器状态为"Up (healthy)"

**将结果截图贴在此处**:

---

## 四、测试数据准备

### 数据库种子数据验证
```sql
-- 连接数据库
docker exec -it gamelink-postgres psql -U gamelink -d gamelink

-- 查看纠纷列表
SELECT od.id, od.order_id, o.order_no, od.user_id, od.player_id, 
       od.reason, od.status, od.assigned_to_user_id, od.created_at
FROM order_disputes od
LEFT JOIN orders o ON od.order_id = o.id
ORDER BY od.created_at DESC;

-- 查看纠纷统计
SELECT status, COUNT(*) as count FROM order_disputes GROUP BY status;

-- 查看管理员列表（用于分配）
SELECT id, name, email FROM users WHERE role = 'admin';
```

### 测试账号
- **管理员**: 使用 `.env` 中的 `SUPER_ADMIN_EMAIL` / `SUPER_ADMIN_PASSWORD`

---

## 五、逐个按钮测试记录

### 功能1: #statsLoad 统计加载

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 访问页面 /admin/biz/dispute
4. 观察统计卡片加载

**Evidence收集**:
- [ ] 截图1: 统计卡片（待处理、处理中、已解决、总计）
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/disputes/stats
  - Response: `{"pending":X,"assigned":X,"mediating":X,"resolved":X,"total":X}`
  - Status: 200
- [ ] 截图3: docker logs处理记录
  ```bash
  docker logs gamelink-backend --tail=20 | findstr "disputes/stats"
  ```
- [ ] 截图4: 数据库验证
  ```sql
  SELECT status, COUNT(*) FROM order_disputes GROUP BY status;
  SELECT COUNT(*) as total FROM order_disputes;
  ```
- [ ] 截图5: 统计数据与数据库一致

**异常场景测试**:
- [ ] 场景A: 无纠纷数据时 → 预期: 所有统计显示0
- [ ] 场景B: 后端服务不可用 → 预期: 错误提示
- [ ] 场景C: 统计API返回错误 → 预期: 显示默认值0

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 功能2: #searchBtn 搜索(订单号)

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 在订单号搜索框输入订单号（如"ORD202412"）
4. 点击搜索

**Evidence收集**:
- [ ] 截图1: 搜索框输入订单号
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/disputes?orderNo=ORD202412&page=1&pageSize=10
  - Response: 过滤后的纠纷列表
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT od.*, o.order_no 
  FROM order_disputes od
  JOIN orders o ON od.order_id = o.id
  WHERE o.order_no LIKE '%ORD202412%';
  ```
- [ ] 截图5: 列表显示过滤结果

**异常场景测试**:
- [ ] 场景A: 搜索不存在的订单号 → 预期: 显示空列表
- [ ] 场景B: 清空搜索框 → 预期: 显示全部纠纷
- [ ] 场景C: 特殊字符搜索 → 预期: 安全过滤

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 功能3: #filterBtn 状态筛选

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 选择状态"待处理"
4. 观察列表过滤

**Evidence收集**:
- [ ] 截图1: 状态筛选下拉框
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/disputes?status=pending&page=1&pageSize=10
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM order_disputes WHERE status = 'pending';
  ```
- [ ] 截图5: 列表只显示待处理状态的纠纷

**异常场景测试**:
- [ ] 场景A: 筛选无结果的状态 → 预期: 显示空列表
- [ ] 场景B: 清空筛选条件 → 预期: 显示全部
- [ ] 场景C: 组合搜索+筛选 → 预期: 两个条件同时生效

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 功能4: #detailBtn 详情

**测试步骤**:
1. 点击某纠纷行的"详情"按钮
2. 观察详情抽屉展示

**Evidence收集**:
- [ ] 截图1: 详情按钮
- [ ] 截图2: 无Network请求（前端展示已有数据）或有详情API请求
- [ ] 截图3: 不适用
- [ ] 截图4: 数据库验证
  ```sql
  SELECT od.*, o.order_no, o.total_price_cents,
         u.name as user_name, p.nickname as player_name
  FROM order_disputes od
  LEFT JOIN orders o ON od.order_id = o.id
  LEFT JOIN users u ON od.user_id = u.id
  LEFT JOIN players p ON od.player_id = p.id
  WHERE od.id = 1;
  ```
- [ ] 截图5: 详情抽屉显示完整信息（纠纷ID、订单号、用户、陪玩师、原因、描述、处理进度）

**异常场景测试**:
- [ ] 场景A: 关闭抽屉后重新打开 → 预期: 正常显示
- [ ] 场景B: 查看已解决纠纷详情 → 预期: 显示解决方案和处理备注

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 功能5: #assignBtn 分配

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 找到状态为"待处理"的纠纷
4. 点击"分配"按钮
5. 选择处理人
6. 点击确定

**Evidence收集**:
- [ ] 截图1: 分配弹窗（处理人下拉框）
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/disputes/1/assign
  - Payload: `{"assignedToUserId":2,"source":"manual"}`
  - Status: 200
- [ ] 截图3: docker logs处理记录
  ```bash
  docker logs gamelink-backend --tail=20 | findstr "assign"
  ```
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, status, assigned_to_user_id, assigned_at 
  FROM order_disputes WHERE id = 1;
  -- status 应变为 'assigned'
  -- assigned_to_user_id 应为选择的用户ID
  -- assigned_at 应有值
  ```
- [ ] 截图5: 列表状态变为"已分配"，统计卡片更新

**异常场景测试**:
- [ ] 场景A: 分配给不存在的用户 → 预期: 后端返回错误
- [ ] 场景B: 重复分配同一纠纷 → 预期: 允许重新分配或提示已分配
- [ ] 场景C: 分配已解决的纠纷 → 预期: 按钮不可用或提示错误
- [ ] 场景D: 快速连续点击确定 → 预期: 防抖生效

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 功能6: #resolveBtn 处理

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 找到状态为"已分配"或"处理中"的纠纷
4. 点击"处理"按钮
5. 填写处理信息：
   - 解决方案: 全额退款
   - 退款金额: 100.00
   - 处理备注: 测试处理
6. 点击确定

**Evidence收集**:
- [ ] 截图1: 处理弹窗（解决方案、退款金额、处理备注）
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/disputes/1/resolve
  - Payload: `{"resolution":"refund","resolutionAmount":10000,"resolutionNotes":"测试处理"}`
  - Status: 200
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, status, resolution, resolution_amount_cents, 
         resolution_notes, resolved_at 
  FROM order_disputes WHERE id = 1;
  -- status 应变为 'resolved'
  -- resolution 应为 'refund'
  -- resolution_amount_cents 应为 10000
  -- resolved_at 应有值
  ```
- [ ] 截图5: 列表状态变为"已解决"，统计卡片更新

**异常场景测试**:
- [ ] 场景A: 解决方案为空 → 预期: 校验失败
- [ ] 场景B: 处理备注为空 → 预期: 校验失败
- [ ] 场景C: 退款金额超过订单金额 → 预期: 校验失败或警告
- [ ] 场景D: 处理待处理状态的纠纷 → 预期: 按钮不可用（需先分配）

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 功能7: #rollbackBtn 回滚

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 找到状态为"已分配"或"处理中"的纠纷
4. 点击"回滚"按钮
5. 确认回滚

**Evidence收集**:
- [ ] 截图1: 回滚确认弹窗
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/disputes/1/rollback
  - Payload: `{"rollbackReason":"管理员手动回滚"}`
  - Status: 200
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, status, assigned_to_user_id, assigned_at 
  FROM order_disputes WHERE id = 1;
  -- status 应变回 'pending'
  -- assigned_to_user_id 应为 NULL
  ```
- [ ] 截图5: 列表状态变为"待处理"，统计卡片更新

**异常场景测试**:
- [ ] 场景A: 回滚待处理的纠纷 → 预期: 按钮不可用
- [ ] 场景B: 回滚已解决的纠纷 → 预期: 按钮不可用或提示错误
- [ ] 场景C: 取消回滚 → 预期: 无请求发送

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 功能8: #exportBtn 导出数据

**测试步骤**:
1. 点击"导出数据"按钮
2. 观察文件下载

**Evidence收集**:
- [ ] 截图1: 导出按钮
- [ ] 截图2: 无Network请求（前端导出）
- [ ] 截图3: 不适用
- [ ] 截图4: 不适用
- [ ] 截图5: 下载的CSV文件内容

**特别验证**:
- [ ] CSV包含列: ID, 订单号, 用户, 陪玩师, 纠纷原因, 状态, 处理人, 解决方案, 创建时间, 解决时间
- [ ] 数据与页面列表一致

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 功能9: #refreshBtn 刷新

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击刷新按钮

**Evidence收集**:
- [ ] 截图1: 刷新按钮
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/disputes
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库当前数据
- [ ] 截图5: 列表和统计卡片数据刷新

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

## 六、统计卡片验证

| 卡片名称 | 验证SQL | 预期值 | 实际值 | 结果 |
|---------|---------|--------|--------|------|
| 待处理 | `SELECT COUNT(*) FROM order_disputes WHERE status='pending';` | | | ☐ |
| 处理中 | `SELECT COUNT(*) FROM order_disputes WHERE status IN ('assigned','mediating');` | | | ☐ |
| 已解决 | `SELECT COUNT(*) FROM order_disputes WHERE status='resolved';` | | | ☐ |
| 总计 | `SELECT COUNT(*) FROM order_disputes;` | | | ☐ |

---

## 七、状态流转验证

| 初始状态 | 操作 | 预期状态 | 实际状态 | 结果 |
|---------|------|---------|---------|------|
| pending | 分配 | assigned | | ☐ |
| assigned | 处理 | resolved | | ☐ |
| assigned | 回滚 | pending | | ☐ |
| mediating | 处理 | resolved | | ☐ |
| mediating | 回滚 | pending | | ☐ |
| resolved | 分配 | 不允许 | | ☐ |

---

## 八、全量测试完整性自查

- [ ] 所有P0按钮已测试（6个）
- [ ] 所有P1按钮已测试（3个）
- [ ] 每个按钮提供5张截图+日志
- [ ] 每个按钮测试异常场景≥3个
- [ ] 统计卡片数据验证通过
- [ ] 状态流转验证通过
- [ ] 所有截图有明确的文件名
- [ ] 日志文件已打包

---

## 九、质量承诺

我承诺以上测试内容真实完整，所有按钮均已按22项清单验证。

**测试人签字**: ___________  
**日期**: ___________

---

## 十、组长审核意见

**审核结果**: ☐ 通过 ☐ 打回重做  
**打回原因**: （如有）  
**审核人**: ___________  
**日期**: ___________

---

**文档版本**: v1.0  
**发布日期**: 2024-12-18
