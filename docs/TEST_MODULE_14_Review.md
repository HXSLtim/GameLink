# 测试任务单：评价管理模块全量测试

**任务编号**: TEST-2024-M14  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/biz/review | #detailBtn | 详情 | GET /api/v1/admin/reviews/:id | P0 | ☐ |
| /admin/biz/review | #approveBtn | 批准 | PUT /api/v1/admin/reviews/:id/approve | P0 | ☐ |
| /admin/biz/review | #rejectBtn | 拒绝 | PUT /api/v1/admin/reviews/:id/reject | P0 | ☐ |
| /admin/biz/review | #deleteBtn | 删除 | DELETE /api/v1/admin/reviews/:id | P0 | ☐ |
| /admin/biz/review | #batchApproveBtn | 批量批准 | PUT /api/v1/admin/reviews/batch/approve | P1 | ☐ |
| /admin/biz/review | #searchBtn | 搜索 | GET /api/v1/admin/reviews?keyword=xxx | P1 | ☐ |
| /admin/biz/review | #statusFilter | 状态筛选 | GET /api/v1/admin/reviews?status=xxx | P1 | ☐ |
| /admin/biz/review | #ratingFilter | 评分筛选 | GET /api/v1/admin/reviews?rating=xxx | P1 | ☐ |
| /admin/biz/review | #refreshBtn | 刷新 | GET /api/v1/admin/reviews | P1 | ☐ |

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

-- 查看评价统计
SELECT status, COUNT(*) FROM reviews GROUP BY status;

-- 查看评价详情
SELECT r.id, r.rating, r.comment, r.status, r.is_reported,
       u.name as reviewer_name, p.nickname as player_name
FROM reviews r
LEFT JOIN users u ON r.reviewer_id = u.id
LEFT JOIN players p ON r.player_id = p.id
LIMIT 10;
```

### 评价状态说明
- `pending`: 待审核
- `approved`: 已通过
- `rejected`: 已拒绝
- `deleted`: 已删除

---

## 五、逐个按钮测试记录

### 按钮1: #detailBtn 详情

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 点击某评价行的"详情"按钮
4. 观察详情抽屉

**Evidence收集**:
- [ ] 截图1: 详情按钮
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/reviews/:id
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM reviews WHERE id = :review_id;
  ```
- [ ] 截图5: 详情抽屉显示完整信息（评价内容、评分、图片、操作日志）

**异常场景测试**:
- [ ] 场景A: 查看不存在的评价 → 预期: 404错误
- [ ] 场景B: 查看被举报的评价 → 预期: 显示举报标记
- [ ] 场景C: 查看带图片的评价 → 预期: 图片正常显示

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮2: #approveBtn 批准

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击待审核评价的"批准"按钮
4. 填写批准原因（可选）
5. 确认批准

**Evidence收集**:
- [ ] 截图1: 批准弹窗
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/reviews/:id/approve
  - Payload: `{"reason":"批准评价"}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, status, approved_at, approved_by FROM reviews WHERE id = :review_id;
  ```
- [ ] 截图5: 列表状态变为"已通过"

**异常场景测试**:
- [ ] 场景A: 批准已批准的评价 → 预期: 按钮不显示
- [ ] 场景B: 批准已拒绝的评价 → 预期: 按钮不显示
- [ ] 场景C: 取消批准 → 预期: 无请求发送

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮3: #rejectBtn 拒绝

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击待审核评价的"拒绝"按钮
4. 填写拒绝原因
5. 确认拒绝

**Evidence收集**:
- [ ] 截图1: 拒绝弹窗
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/reviews/:id/reject
  - Payload: `{"reason":"内容不当"}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, status, reject_reason FROM reviews WHERE id = :review_id;
  ```
- [ ] 截图5: 列表状态变为"已拒绝"

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮4: #deleteBtn 删除

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某评价的"删除"按钮
4. 确认删除

**Evidence收集**:
- [ ] 截图1: 删除确认弹窗
- [ ] 截图2: Network请求详情
  - URL: DELETE /api/v1/admin/reviews/:id
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM reviews WHERE id = :deleted_id;
  ```
- [ ] 截图5: 列表中评价已消失或状态变为"已删除"

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮5: #batchApproveBtn 批量批准

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 勾选多个待审核评价
4. 点击"批量批准"按钮
5. 确认批准

**Evidence收集**:
- [ ] 截图1: 批量批准弹窗
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/reviews/batch/approve
  - Payload: `{"ids":[1,2,3]}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, status FROM reviews WHERE id IN (1,2,3);
  ```
- [ ] 截图5: 所有选中评价状态变为"已通过"

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮6-9: 搜索/筛选/刷新

**搜索 #searchBtn**:
- [ ] 输入订单ID或评价内容搜索
- [ ] 验证: GET /api/v1/admin/reviews?keyword=xxx

**状态筛选 #statusFilter**:
- [ ] 选择特定状态
- [ ] 验证: GET /api/v1/admin/reviews?status=pending

**评分筛选 #ratingFilter**:
- [ ] 选择特定评分
- [ ] 验证: GET /api/v1/admin/reviews?rating=5

**刷新 #refreshBtn**:
- [ ] 点击刷新按钮
- [ ] 验证: GET /api/v1/admin/reviews

---

## 六、全量测试完整性自查

- [ ] 所有P0按钮已测试（4个）
- [ ] 所有P1按钮已测试（5个）
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
