# 测试任务单：结算公司管理模块全量测试

**任务编号**: TEST-2024-M03  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/finance/settlement-company | #createBtn | 新增公司 | POST /api/v1/admin/settlement-companies | P0 | ☐ |
| /admin/finance/settlement-company | #editBtn | 编辑 | PUT /api/v1/admin/settlement-companies/:id | P0 | ☐ |
| /admin/finance/settlement-company | #toggleBtn | 启用/禁用开关 | POST /api/v1/admin/settlement-companies/:id/toggle | P0 | ☐ |
| /admin/finance/settlement-company | #searchBtn | 搜索 | GET /api/v1/admin/settlement-companies | P1 | ☐ |
| /admin/finance/settlement-company | #exportBtn | 导出数据 | 前端CSV导出 | P1 | ☐ |
| /admin/finance/settlement-company | #refreshBtn | 刷新 | GET /api/v1/admin/settlement-companies | P1 | ☐ |

**重要**: 以上6个按钮，必须全部测试完成，少一个 = 任务未完成

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

-- 查看现有结算公司
SELECT id, name, credit_code, contact_name, contact_phone, status, created_at 
FROM settlement_companies ORDER BY id;

-- 查看公司关联陪玩师数
SELECT sc.id, sc.name, COUNT(p.id) as player_count 
FROM settlement_companies sc 
LEFT JOIN players p ON sc.id = p.settlement_company_id 
GROUP BY sc.id, sc.name;
```

### 测试账号
- **管理员**: admin@gameLink.com / Admin2025@Pass#

---

## 五、逐个按钮测试记录

### 按钮1: #createBtn 新增公司

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 点击"新增公司"按钮
4. 填写表单：
   - 公司名称: 测试结算公司有限公司
   - 统一社会信用代码: 91110000MA00TEST01
   - 联系人: 张三
   - 联系电话: 13800138001
   - 地址: 北京市朝阳区测试路1号
   - 银行名称: 中国工商银行
   - 银行账号: 6222021234567890123
   - 开户支行: 北京朝阳支行
5. 点击确定

**Evidence收集**:
- [ ] 截图1: 新增公司弹窗表单
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/settlement-companies
  - Payload: 完整表单数据
  - Status: 200
- [ ] 截图3: docker logs处理记录
  ```bash
  docker logs gamelink-backend --tail=20 | findstr "settlement"
  ```
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM settlement_companies WHERE credit_code = '91110000MA00TEST01';
  ```
- [ ] 截图5: 列表刷新显示新公司

**异常场景测试**:
- [ ] 场景A: 公司名称为空 → 预期: 前端校验提示"请输入公司名称"
- [ ] 场景B: 信用代码不是18位 → 预期: 校验提示"必须为18位"
- [ ] 场景C: 信用代码重复 → 预期: 后端返回错误"信用代码已存在"
- [ ] 场景D: 快速连续点击确定 → 预期: 防抖生效

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮2: #editBtn 编辑

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某公司行的"编辑"按钮
4. 修改联系人为"李四"
5. 点击确定

**Evidence收集**:
- [ ] 截图1: 编辑弹窗（含原数据回显）
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/settlement-companies/1
  - Payload: 修改后的数据
  - **注意**: 信用代码字段应为disabled不可编辑
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, name, contact_name, updated_at FROM settlement_companies WHERE id = 1;
  -- 验证 contact_name = '李四', updated_at 已更新
  ```
- [ ] 截图5: 列表显示更新后的联系人

**异常场景测试**:
- [ ] 场景A: 清空必填字段 → 预期: 校验失败
- [ ] 场景B: 编辑不存在的公司 → 预期: 404错误
- [ ] 场景C: 信用代码字段是否可编辑 → 预期: 不可编辑（disabled）

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮3: #toggleBtn 启用/禁用开关

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某公司的启用/禁用开关
4. 观察状态变化

**Evidence收集**:
- [ ] 截图1: 开关当前状态（启用/禁用）
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/settlement-companies/1/toggle
  - Payload: `{"enabled":false}` 或 `{"enabled":true}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, name, status FROM settlement_companies WHERE id = 1;
  -- status 应为 'active' 或 'inactive'
  ```
- [ ] 截图5: 开关状态已切换，统计卡片"启用公司"数量变化

**异常场景测试**:
- [ ] 场景A: 禁用有关联陪玩师的公司 → 预期: 允许禁用或警告提示
- [ ] 场景B: 快速连续切换 → 预期: 防抖生效，状态正确
- [ ] 场景C: 后端服务不可用时切换 → 预期: 错误提示，开关回滚

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮4: #searchBtn 搜索

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 输入关键词"测试"
4. 选择状态"启用"
5. 点击搜索

**Evidence收集**:
- [ ] 截图1: 搜索条件填写
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/settlement-companies?keyword=测试&status=active&page=1&pageSize=10
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM settlement_companies 
  WHERE (name LIKE '%测试%' OR credit_code LIKE '%测试%') 
  AND status = 'active';
  ```
- [ ] 截图5: 列表显示过滤结果

**异常场景测试**:
- [ ] 场景A: 搜索无结果 → 预期: 显示空列表"暂无数据"
- [ ] 场景B: 清空搜索条件 → 预期: 显示全部公司
- [ ] 场景C: 特殊字符搜索 → 预期: 安全过滤，无报错

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮5: #exportBtn 导出数据

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
- [ ] CSV文件名格式: settlement_companies_YYYYMMDD.csv
- [ ] CSV包含列: ID, 公司名称, 统一社会信用代码, 联系人, 联系电话, 银行名称, 陪玩师数, 状态, 创建时间
- [ ] 数据与页面列表一致
- [ ] 敏感信息（银行账号）是否脱敏

**异常场景测试**:
- [ ] 场景A: 空列表导出 → 预期: 只有表头的CSV
- [ ] 场景B: 大量数据导出 → 预期: 正常导出

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮6: #refreshBtn 刷新

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击刷新按钮

**Evidence收集**:
- [ ] 截图1: 刷新按钮
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/settlement-companies
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库当前数据
  ```sql
  SELECT COUNT(*) FROM settlement_companies;
  ```
- [ ] 截图5: 列表数据刷新，统计卡片更新

**异常场景测试**:
- [ ] 场景A: 后端服务不可用时刷新 → 预期: 错误提示
- [ ] 场景B: 快速连续点击刷新 → 预期: 防抖生效

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

## 六、统计卡片验证

页面顶部有3个统计卡片，需验证数据准确性：

| 卡片名称 | 验证SQL | 预期值 | 实际值 | 结果 |
|---------|---------|--------|--------|------|
| 公司总数 | `SELECT COUNT(*) FROM settlement_companies;` | | | ☐ |
| 启用公司 | `SELECT COUNT(*) FROM settlement_companies WHERE status='active';` | | | ☐ |
| 关联陪玩师 | `SELECT COUNT(*) FROM players WHERE settlement_company_id IS NOT NULL;` | | | ☐ |

---

## 七、全量测试完整性自查

- [ ] 所有P0按钮已测试（3个）
- [ ] 所有P1按钮已测试（3个）
- [ ] 每个按钮提供5张截图+日志
- [ ] 每个按钮测试异常场景≥3个
- [ ] 统计卡片数据验证通过
- [ ] 所有截图有明确的文件名
- [ ] 日志文件已打包

---

## 八、质量承诺

我承诺以上测试内容真实完整，所有按钮均已按22项清单验证。

**测试人签字**: ___________  
**日期**: ___________

---

## 九、组长审核意见

**审核结果**: ☐ 通过 ☐ 打回重做  
**打回原因**: （如有）  
**审核人**: ___________  
**日期**: ___________

---

**文档版本**: v1.0  
**发布日期**: 2024-12-18
