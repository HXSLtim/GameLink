# 测试任务单：陪玩师管理模块全量测试

**任务编号**: TEST-2024-M11  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/biz/player | #detailBtn | 详情 | GET /api/v1/admin/players/:id | P0 | ☐ |
| /admin/biz/player | #auditBtn | 审核 | PUT /api/v1/admin/players/:id/verify | P0 | ☐ |
| /admin/biz/player | #banBtn | 封禁 | PUT /api/v1/admin/players/:id/verify | P0 | ☐ |
| /admin/biz/player | #unbanBtn | 解封 | PUT /api/v1/admin/players/:id/verify | P0 | ☐ |
| /admin/biz/player | #batchStatusBtn | 批量修改状态 | PUT /api/v1/admin/players/batch/status | P1 | ☐ |
| /admin/biz/player | #batchDeleteBtn | 批量删除 | DELETE /api/v1/admin/players/batch | P1 | ☐ |
| /admin/biz/player | #exportBtn | 导出数据 | 前端CSV导出 | P2 | ☐ |
| /admin/biz/player | #searchBtn | 搜索 | GET /api/v1/admin/players?keyword=xxx | P1 | ☐ |
| /admin/biz/player | #refreshBtn | 刷新 | GET /api/v1/admin/players | P1 | ☐ |

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

---

## 四、测试数据准备

### 数据库种子数据验证
```sql
-- 连接数据库
docker exec -it gamelink-postgres psql -U gamelink -d gamelink

-- 查看陪玩师统计
SELECT verification_status, COUNT(*) FROM players GROUP BY verification_status;

-- 查看陪玩师详情
SELECT p.id, p.nickname, p.verification_status, u.name as user_name, g.name as game_name
FROM players p
LEFT JOIN users u ON p.user_id = u.id
LEFT JOIN games g ON p.main_game_id = g.id
LIMIT 10;
```

### 陪玩师状态说明
- `pending`: 待审核
- `verified`: 已通过
- `rejected`: 已拒绝

---

## 五、逐个按钮测试记录

### 按钮1: #detailBtn 详情

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 点击某陪玩师行的"详情"按钮
4. 观察详情抽屉

**Evidence收集**:
- [ ] 截图1: 详情按钮
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/players/:id
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM players WHERE id = :player_id;
  ```
- [ ] 截图5: 详情抽屉显示完整信息（基本信息、评分、审核信息）

**异常场景测试**:
- [ ] 场景A: 查看不存在的陪玩师 → 预期: 404错误
- [ ] 场景B: 关闭后重新打开 → 预期: 重新加载数据
- [ ] 场景C: 网络中断时查看 → 预期: 错误提示

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮2: #auditBtn 审核

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击待审核陪玩师的"审核"按钮
4. 填写审核备注
5. 点击"通过"或"拒绝"

**Evidence收集**:
- [ ] 截图1: 审核弹窗（含申请人信息）
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/players/:id/verify
  - Payload: `{"status":"verified","remark":"审核通过"}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, verification_status, verified_at, verified_by, verify_remark 
  FROM players WHERE id = :player_id;
  ```
- [ ] 截图5: 列表状态更新为"已通过"

**异常场景测试**:
- [ ] 场景A: 审核通过 → 预期: 状态变为verified
- [ ] 场景B: 审核拒绝 → 预期: 状态变为rejected，记录拒绝原因
- [ ] 场景C: 审核已审核的陪玩师 → 预期: 按钮不显示

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮3: #banBtn 封禁

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击已通过陪玩师的"封禁"按钮
4. 确认操作

**Evidence收集**:
- [ ] 截图1: 封禁确认弹窗
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/players/:id/verify
  - Payload: `{"status":"rejected"}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, verification_status FROM players WHERE id = :player_id;
  ```
- [ ] 截图5: 列表状态变为"已拒绝"

**异常场景测试**:
- [ ] 场景A: 封禁待审核陪玩师 → 预期: 按钮不显示
- [ ] 场景B: 封禁已封禁陪玩师 → 预期: 按钮不显示
- [ ] 场景C: 取消封禁 → 预期: 无请求发送

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮4: #unbanBtn 解封

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击已拒绝陪玩师的"解封"按钮
4. 确认操作

**Evidence收集**:
- [ ] 截图1: 解封确认弹窗
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/players/:id/verify
  - Payload: `{"status":"verified"}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, verification_status FROM players WHERE id = :player_id;
  ```
- [ ] 截图5: 列表状态变为"已通过"

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮5-6: 批量操作

**批量修改状态 #batchStatusBtn**:
- [ ] 选中多个陪玩师或选择筛选条件
- [ ] 选择目标状态
- [ ] 验证: PUT /api/v1/admin/players/batch/status
- [ ] 数据库验证所有选中陪玩师状态已更新

**批量删除 #batchDeleteBtn**:
- [ ] 选中多个陪玩师或选择筛选条件
- [ ] 确认删除
- [ ] 验证: DELETE /api/v1/admin/players/batch
- [ ] 数据库验证所有选中陪玩师已删除

---

### 按钮7-9: 导出/搜索/刷新

**导出 #exportBtn**:
- [ ] 点击导出数据
- [ ] 验证CSV文件下载
- [ ] 验证数据完整性

**搜索 #searchBtn**:
- [ ] 输入关键词搜索
- [ ] 验证: GET /api/v1/admin/players?keyword=xxx
- [ ] 验证状态筛选

**刷新 #refreshBtn**:
- [ ] 点击刷新按钮
- [ ] 验证: GET /api/v1/admin/players
- [ ] 验证数据与数据库一致

---

## 六、全量测试完整性自查

- [ ] 所有P0按钮已测试（4个）
- [ ] 所有P1按钮已测试（4个）
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
**打回原因**: （如有）  
**审核人**: ___________  
**日期**: ___________

---

**文档版本**: v1.0  
**发布日期**: 2024-12-18
