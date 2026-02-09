# 测试任务单：告警管理模块全量测试

**任务编号**: TEST-2024-M18  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/sys/alert | #createRuleBtn | 新增规则 | POST /api/v1/admin/alerts/rules | P0 | ☐ |
| /admin/sys/alert | #editRuleBtn | 编辑 | PUT /api/v1/admin/alerts/rules/:id | P0 | ☐ |
| /admin/sys/alert | #deleteRuleBtn | 删除 | DELETE /api/v1/admin/alerts/rules/:id | P0 | ☐ |
| /admin/sys/alert | #statusSwitch | 启用/禁用 | PUT /api/v1/admin/alerts/rules/:id/status | P0 | ☐ |
| /admin/sys/alert | #acknowledgeBtn | 确认告警 | PUT /api/v1/admin/alerts/:id/acknowledge | P0 | ☐ |
| /admin/sys/alert | #resolveBtn | 解决告警 | PUT /api/v1/admin/alerts/:id/resolve | P0 | ☐ |
| /admin/sys/alert | #tabSwitch | Tab切换 | 前端切换 | P1 | ☐ |
| /admin/sys/alert | #refreshBtn | 刷新 | GET /api/v1/admin/alerts | P1 | ☐ |

**重要**: 以上8个按钮，必须全部测试完成

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

-- 查看告警规则
SELECT id, name, type, level, metric, threshold, is_active FROM alert_rules;

-- 查看告警记录
SELECT id, rule_name, level, status, triggered_at FROM alert_records ORDER BY triggered_at DESC LIMIT 10;

-- 查看告警统计
SELECT status, COUNT(*) FROM alert_records GROUP BY status;
```

### 告警类型说明
- `system`: 系统告警（CPU、内存、磁盘等）
- `business`: 业务告警（订单异常、支付失败等）
- `security`: 安全告警（异常登录、攻击检测等）

### 告警级别说明
- `critical`: 严重（需立即处理）
- `warning`: 警告（需关注）
- `info`: 信息（仅通知）

---

## 四、逐个按钮测试记录

### 按钮1: #createRuleBtn 新增规则

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 切换到"告警规则"Tab
4. 点击"新增规则"按钮
5. 填写表单：
   - 规则名称: 测试CPU告警
   - 告警类型: 系统告警
   - 告警级别: 警告
   - 监控指标: cpu_usage
   - 条件: >
   - 阈值: 80
   - 持续时间: 5分钟
   - 通知渠道: 邮件、短信
6. 点击保存

**Evidence收集**:
- [ ] 截图1: 新增规则弹窗表单
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/alerts/rules
  - Payload: `{"name":"测试CPU告警","type":"system","level":"warning",...}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM alert_rules WHERE name = '测试CPU告警';
  ```
- [ ] 截图5: 规则列表显示新规则

**异常场景测试**:
- [ ] 场景A: 名称为空 → 预期: 校验提示
- [ ] 场景B: 阈值为负数 → 预期: 校验提示
- [ ] 场景C: 未选择通知渠道 → 预期: 校验提示

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮2: #editRuleBtn 编辑规则

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某规则的"编辑"按钮
4. 修改阈值
5. 点击保存

**Evidence收集**:
- [ ] 截图1: 编辑弹窗（含原数据回显）
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/alerts/rules/:id
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
- [ ] 截图5: 规则列表显示更新后的数据

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮3: #deleteRuleBtn 删除规则

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某规则的"删除"按钮
4. 确认删除

**Evidence收集**:
- [ ] 截图1: 删除确认弹窗
- [ ] 截图2: Network请求详情
  - URL: DELETE /api/v1/admin/alerts/rules/:id
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
- [ ] 截图5: 规则列表中规则已消失

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮4: #statusSwitch 启用/禁用

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某规则的状态开关
4. 观察状态变化

**Evidence收集**:
- [ ] 截图1: 状态开关
- [ ] 截图2: Network请求详情
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
- [ ] 截图5: 开关状态已切换

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮5: #acknowledgeBtn 确认告警

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 在告警记录Tab找到"触发中"的告警
4. 点击"确认"按钮

**Evidence收集**:
- [ ] 截图1: 确认按钮
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/alerts/:id/acknowledge
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, status, acknowledged_at FROM alert_records WHERE id = :alert_id;
  ```
- [ ] 截图5: 告警状态变为"已确认"

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮6: #resolveBtn 解决告警

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 找到"触发中"或"已确认"的告警
4. 点击"解决"按钮

**Evidence收集**:
- [ ] 截图1: 解决按钮
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/alerts/:id/resolve
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, status, resolved_at FROM alert_records WHERE id = :alert_id;
  ```
- [ ] 截图5: 告警状态变为"已解决"

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮7-8: Tab切换/刷新

**Tab切换 #tabSwitch**:
- [ ] 切换到"告警记录"Tab
- [ ] 切换到"告警规则"Tab
- [ ] 验证数据正确加载

**刷新 #refreshBtn**:
- [ ] 点击刷新按钮
- [ ] 验证数据重新加载

---

## 五、全量测试完整性自查

- [ ] 所有P0按钮已测试（6个）
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
