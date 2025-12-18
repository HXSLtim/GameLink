# 测试任务单：分流规则管理模块全量测试

**任务编号**: TEST-2024-M02  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/biz/routing-rule | #createBtn | 新增规则 | POST /api/v1/admin/routing-rules | P0 | ☐ |
| /admin/biz/routing-rule | #editBtn | 编辑 | PUT /api/v1/admin/routing-rules/:id | P0 | ☐ |
| /admin/biz/routing-rule | #deleteBtn | 删除 | DELETE /api/v1/admin/routing-rules/:id | P0 | ☐ |
| /admin/biz/routing-rule | #toggleBtn | 启用/禁用开关 | POST /api/v1/admin/routing-rules/:id/toggle | P0 | ☐ |
| /admin/biz/routing-rule | #historyBtn | 历史 | GET /api/v1/admin/routing-rules/:id/history | P1 | ☐ |
| /admin/biz/routing-rule | #testBtn | 测试规则 | POST /api/v1/admin/routing-rules/test | P1 | ☐ |
| /admin/biz/routing-rule | #exportBtn | 导出数据 | 前端CSV导出 | P1 | ☐ |
| /admin/biz/routing-rule | #searchBtn | 搜索 | GET /api/v1/admin/routing-rules | P1 | ☐ |
| /admin/biz/routing-rule | #refreshBtn | 刷新 | GET /api/v1/admin/routing-rules | P1 | ☐ |

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

-- 查看现有分流规则
SELECT id, name, priority, status, target_entity_id, conditions, created_at 
FROM routing_rules ORDER BY priority ASC;

-- 查看收款主体（用于关联）
SELECT id, name FROM collection_entities WHERE status = 'active';
```

### 测试账号
- **管理员**: admin@gameLink.com / Admin2025@Pass#

---

## 五、逐个按钮测试记录

### 按钮1: #createBtn 新增规则

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 点击"新增规则"按钮
4. 填写表单：
   - 规则名称: 大额订单分流
   - 优先级: 10
   - 目标收款主体ID: 1
   - 条件类型: 金额范围
   - 条件值: 1000-5000
5. 点击确定

**Evidence收集**:
- [ ] 截图1: 新增规则弹窗表单
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/routing-rules
  - Payload: `{"name":"大额订单分流","priority":10,"targetEntityId":1,"conditions":[{"field":"order_amount","operator":"eq","value":"1000-5000"}]}`
  - Status: 200
- [ ] 截图3: docker logs处理记录
  ```bash
  docker logs gamelink-backend --tail=20 | findstr "routing-rules"
  ```
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM routing_rules WHERE name = '大额订单分流';
  ```
- [ ] 截图5: 列表刷新显示新规则

**异常场景测试**:
- [ ] 场景A: 名称为空提交 → 预期: 前端校验提示
- [ ] 场景B: 优先级重复 → 预期: 后端返回错误或允许
- [ ] 场景C: 目标主体ID不存在 → 预期: 后端返回错误

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮2: #editBtn 编辑

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某规则行的"编辑"按钮
4. 修改规则名称为"大额订单分流-已修改"
5. 点击确定

**Evidence收集**:
- [ ] 截图1: 编辑弹窗（含原数据回显）
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/routing-rules/1
  - Payload: 修改后的数据
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM routing_rules WHERE id = 1;
  ```
- [ ] 截图5: 列表显示更新后的规则名

**异常场景测试**:
- [ ] 场景A: 修改为无效的条件值 → 预期: 校验失败
- [ ] 场景B: 编辑不存在的规则 → 预期: 404错误
- [ ] 场景C: 并发编辑同一规则 → 预期: 后提交者覆盖

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮3: #deleteBtn 删除

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某规则行的"删除"按钮
4. 确认删除

**Evidence收集**:
- [ ] 截图1: 删除确认（如有）
- [ ] 截图2: Network请求详情
  - URL: DELETE /api/v1/admin/routing-rules/1
  - Status: 200
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM routing_rules WHERE id = 1;
  -- 应返回空或deleted_at不为空
  ```
- [ ] 截图5: 列表中规则已消失

**异常场景测试**:
- [ ] 场景A: 删除启用中的规则 → 预期: 允许删除或提示先禁用
- [ ] 场景B: 删除不存在的规则 → 预期: 404错误
- [ ] 场景C: 删除后立即刷新 → 预期: 规则不再显示

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮4: #toggleBtn 启用/禁用开关

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某规则的启用/禁用开关
4. 观察状态变化

**Evidence收集**:
- [ ] 截图1: 开关当前状态
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/routing-rules/1/toggle
  - Payload: `{"enabled":false}` 或 `{"enabled":true}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, name, status FROM routing_rules WHERE id = 1;
  -- status 应为 'active' 或 'inactive'
  ```
- [ ] 截图5: 开关状态已切换

**异常场景测试**:
- [ ] 场景A: 快速连续切换 → 预期: 防抖生效，状态正确
- [ ] 场景B: 禁用唯一启用的规则 → 预期: 允许或警告
- [ ] 场景C: 后端服务不可用时切换 → 预期: 错误提示

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮5: #historyBtn 历史

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某规则的"历史"按钮
4. 观察历史记录弹窗

**Evidence收集**:
- [ ] 截图1: 历史按钮
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/routing-rules/1/history
  - Response: 历史记录列表
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM routing_rule_histories WHERE rule_id = 1 ORDER BY created_at DESC;
  ```
- [ ] 截图5: 历史记录弹窗显示时间线

**异常场景测试**:
- [ ] 场景A: 无历史记录的规则 → 预期: 显示空状态
- [ ] 场景B: 大量历史记录 → 预期: 分页或滚动加载

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮6: #testBtn 测试规则

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击"测试规则"按钮
4. 填写测试参数：
   - 订单金额: 2000
   - 游戏ID: 1
5. 点击确定

**Evidence收集**:
- [ ] 截图1: 测试规则弹窗
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/routing-rules/test
  - Payload: `{"amount":2000,"gameId":1}`
  - Response: 匹配结果
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 验证匹配逻辑
  ```sql
  -- 查看应该匹配的规则
  SELECT * FROM routing_rules 
  WHERE status = 'active' 
  AND conditions::text LIKE '%order_amount%'
  ORDER BY priority ASC;
  ```
- [ ] 截图5: 测试结果显示匹配的规则和目标主体

**异常场景测试**:
- [ ] 场景A: 无匹配规则 → 预期: 显示"默认规则"
- [ ] 场景B: 金额为负数 → 预期: 校验失败
- [ ] 场景C: 多规则匹配 → 预期: 返回优先级最高的

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮7: #exportBtn 导出数据

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
- [ ] CSV包含列: ID, 规则名称, 优先级, 目标主体, 状态, 创建时间
- [ ] 数据与页面列表一致

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮8: #searchBtn 搜索

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 输入关键词"大额"
4. 选择状态"启用"
5. 点击搜索

**Evidence收集**:
- [ ] 截图1: 搜索条件
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/routing-rules?keyword=大额&status=active
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM routing_rules 
  WHERE name LIKE '%大额%' AND status = 'active';
  ```
- [ ] 截图5: 列表显示过滤结果

**异常场景测试**:
- [ ] 场景A: 搜索无结果 → 预期: 显示空列表
- [ ] 场景B: 清空搜索条件 → 预期: 显示全部
- [ ] 场景C: SQL注入尝试 → 预期: 安全过滤

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮9: #refreshBtn 刷新

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击刷新按钮

**Evidence收集**:
- [ ] 截图1: 刷新按钮
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/routing-rules
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库当前数据
- [ ] 截图5: 列表数据刷新

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

## 六、全量测试完整性自查

- [ ] 所有P0按钮已测试（4个）
- [ ] 所有P1按钮已测试（5个）
- [ ] 每个按钮提供5张截图+日志
- [ ] 每个按钮测试异常场景≥3个
- [ ] 所有截图有明确的文件名
- [ ] 日志文件已打包

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
