# 测试任务单：排行榜抽成配置模块全量测试

**任务编号**: TEST-2024-M04  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/finance/ranking-commission | #createBtn | 新增配置 | POST /api/v1/admin/ranking-commission/configs | P0 | ☐ |
| /admin/finance/ranking-commission | #editBtn | 编辑 | PUT /api/v1/admin/ranking-commission/configs/:id | P0 | ☐ |
| /admin/finance/ranking-commission | #deleteBtn | 删除 | DELETE /api/v1/admin/ranking-commission/configs/:id | P0 | ☐ |
| /admin/finance/ranking-commission | #addRuleBtn | 添加规则(表单内) | 表单内操作 | P1 | ☐ |
| /admin/finance/ranking-commission | #removeRuleBtn | 删除规则(表单内) | 表单内操作 | P1 | ☐ |
| /admin/finance/ranking-commission | #searchBtn | 搜索 | GET /api/v1/admin/ranking-commission/configs | P1 | ☐ |
| /admin/finance/ranking-commission | #exportBtn | 导出数据 | 前端CSV导出 | P1 | ☐ |
| /admin/finance/ranking-commission | #refreshBtn | 刷新 | GET /api/v1/admin/ranking-commission/configs | P1 | ☐ |

**重要**: 以上8个按钮，必须全部测试完成，少一个 = 任务未完成

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

-- 查看现有抽成配置
SELECT id, name, ranking_type, month, rules, is_active, created_at 
FROM ranking_commission_configs ORDER BY month DESC, ranking_type;

-- 查看规则详情（JSON格式）
SELECT id, name, jsonb_array_length(rules) as rule_count, rules 
FROM ranking_commission_configs;
```

### 测试账号
- **管理员**: admin@gameLink.com / Admin2025@Pass#

---

## 五、逐个按钮测试记录

### 按钮1: #createBtn 新增配置

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 点击"新增配置"按钮
4. 填写表单：
   - 配置名称: 2024年12月收入排行抽成
   - 排行类型: 收入排行
   - 月份: 2024-12
   - 状态: 启用
   - 描述: 测试用配置
5. 添加抽成规则：
   - 规则1: 1-10名, 抽成5%
   - 规则2: 11-50名, 抽成3%
   - 规则3: 51-100名, 抽成2%
6. 点击确定

**Evidence收集**:
- [ ] 截图1: 新增配置弹窗表单（含规则列表）
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/ranking-commission/configs
  - Payload: 
    ```json
    {
      "name": "2024年12月收入排行抽成",
      "rankingType": "income",
      "month": "2024-12",
      "isActive": true,
      "description": "测试用配置",
      "rules": [
        {"rankStart": 1, "rankEnd": 10, "commissionRate": 5},
        {"rankStart": 11, "rankEnd": 50, "commissionRate": 3},
        {"rankStart": 51, "rankEnd": 100, "commissionRate": 2}
      ]
    }
    ```
  - Status: 200
- [ ] 截图3: docker logs处理记录
  ```bash
  docker logs gamelink-backend --tail=20 | findstr "ranking-commission"
  ```
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM ranking_commission_configs 
  WHERE month = '2024-12' AND ranking_type = 'income';
  ```
- [ ] 截图5: 列表刷新显示新配置

**异常场景测试**:
- [ ] 场景A: 配置名称为空 → 预期: 前端校验提示
- [ ] 场景B: 同月份同类型重复创建 → 预期: 后端返回错误"该月份已存在同类型配置"
- [ ] 场景C: 规则排名范围重叠(1-10, 5-20) → 预期: 校验失败
- [ ] 场景D: 抽成比例超过100% → 预期: 校验失败
- [ ] 场景E: 无规则提交 → 预期: 校验失败"至少添加一条规则"

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮2: #editBtn 编辑

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某配置行的"编辑"按钮
4. 修改配置名称为"2024年12月收入排行抽成-已修改"
5. 修改第一条规则抽成比例为6%
6. 点击确定

**Evidence收集**:
- [ ] 截图1: 编辑弹窗（含原数据回显，规则列表回显）
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/ranking-commission/configs/1
  - Payload: 修改后的完整数据
  - **注意**: 排行类型和月份字段应为disabled不可编辑
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, name, rules, updated_at FROM ranking_commission_configs WHERE id = 1;
  -- 验证 name 和 rules 已更新
  ```
- [ ] 截图5: 列表显示更新后的配置名

**异常场景测试**:
- [ ] 场景A: 排行类型字段是否可编辑 → 预期: 不可编辑（disabled）
- [ ] 场景B: 月份字段是否可编辑 → 预期: 不可编辑（disabled）
- [ ] 场景C: 删除所有规则后保存 → 预期: 校验失败

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮3: #deleteBtn 删除

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某配置行的"删除"按钮
4. 确认删除

**Evidence收集**:
- [ ] 截图1: 删除确认（如有）
- [ ] 截图2: Network请求详情
  - URL: DELETE /api/v1/admin/ranking-commission/configs/1
  - Status: 200
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM ranking_commission_configs WHERE id = 1;
  -- 应返回空或deleted_at不为空
  ```
- [ ] 截图5: 列表中配置已消失

**异常场景测试**:
- [ ] 场景A: 删除启用中的配置 → 预期: 允许删除或提示先禁用
- [ ] 场景B: 删除不存在的配置 → 预期: 404错误
- [ ] 场景C: 删除当月正在使用的配置 → 预期: 警告提示

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮4: #addRuleBtn 添加规则(表单内)

**测试步骤**:
1. 打开新增/编辑弹窗
2. 点击"添加规则"按钮
3. 观察新增一行规则输入框

**Evidence收集**:
- [ ] 截图1: 添加规则按钮
- [ ] 截图2: 无Network请求（表单内操作）
- [ ] 截图3: 不适用
- [ ] 截图4: 不适用
- [ ] 截图5: 新增一行规则输入框（起始排名、结束排名、抽成比例）

**异常场景测试**:
- [ ] 场景A: 连续添加多条规则 → 预期: 正常添加
- [ ] 场景B: 添加后不填写直接保存 → 预期: 校验失败"必填"

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮5: #removeRuleBtn 删除规则(表单内)

**测试步骤**:
1. 打开新增/编辑弹窗
2. 点击某规则行的"删除"按钮
3. 观察规则行被移除

**Evidence收集**:
- [ ] 截图1: 删除规则按钮
- [ ] 截图2: 无Network请求（表单内操作）
- [ ] 截图3: 不适用
- [ ] 截图4: 不适用
- [ ] 截图5: 规则行已移除

**异常场景测试**:
- [ ] 场景A: 删除最后一条规则 → 预期: 允许删除，保存时校验失败
- [ ] 场景B: 删除后撤销（关闭弹窗不保存） → 预期: 原数据不变

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮6: #searchBtn 搜索

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 输入月份"2024-12"
4. 选择排行类型"收入排行"
5. 点击搜索

**Evidence收集**:
- [ ] 截图1: 搜索条件填写
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/ranking-commission/configs?month=2024-12&rankingType=income
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM ranking_commission_configs 
  WHERE month = '2024-12' AND ranking_type = 'income';
  ```
- [ ] 截图5: 列表显示过滤结果

**异常场景测试**:
- [ ] 场景A: 搜索无结果 → 预期: 显示空列表
- [ ] 场景B: 清空搜索条件 → 预期: 显示全部配置
- [ ] 场景C: 月份格式错误 → 预期: 校验提示

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
- [ ] CSV包含列: ID, 配置名称, 排行类型, 月份, 规则数, 状态, 创建时间
- [ ] 数据与页面列表一致

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮8: #refreshBtn 刷新

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击刷新按钮

**Evidence收集**:
- [ ] 截图1: 刷新按钮
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/ranking-commission/configs
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库当前数据
- [ ] 截图5: 列表数据刷新

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

## 六、业务逻辑验证

### 规则排名范围验证
| 测试场景 | 输入 | 预期结果 | 实际结果 | 结果 |
|---------|------|---------|---------|------|
| 正常范围 | 1-10, 11-50 | 通过 | | ☐ |
| 范围重叠 | 1-10, 5-20 | 校验失败 | | ☐ |
| 起始>结束 | 10-1 | 校验失败 | | ☐ |
| 负数排名 | -1-10 | 校验失败 | | ☐ |

### 抽成比例验证
| 测试场景 | 输入 | 预期结果 | 实际结果 | 结果 |
|---------|------|---------|---------|------|
| 正常比例 | 5% | 通过 | | ☐ |
| 超过100% | 150% | 校验失败 | | ☐ |
| 负数比例 | -5% | 校验失败 | | ☐ |
| 小数比例 | 5.5% | 通过或校验失败 | | ☐ |

---

## 七、全量测试完整性自查

- [ ] 所有P0按钮已测试（3个）
- [ ] 所有P1按钮已测试（5个）
- [ ] 每个按钮提供5张截图+日志
- [ ] 每个按钮测试异常场景≥3个
- [ ] 业务逻辑验证通过
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
