# 测试任务单：提现分流统计模块全量测试

**任务编号**: TEST-2024-M05  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/finance/withdraw-routing | #pageLoad | 页面加载 | GET /api/v1/admin/withdrawals/routing-stats | P0 | ☐ |
| /admin/finance/withdraw-routing | #dateRangeBtn | 日期范围选择 | GET /api/v1/admin/withdrawals/by-company | P0 | ☐ |
| /admin/finance/withdraw-routing | #searchBtn | 搜索 | GET /api/v1/admin/withdrawals/by-company | P0 | ☐ |
| /admin/finance/withdraw-routing | #monthlyReportBtn | 月报 | GET /api/v1/admin/withdrawals/routing-report | P1 | ☐ |
| /admin/finance/withdraw-routing | #exportBtn | 导出数据 | 前端CSV导出 | P1 | ☐ |
| /admin/finance/withdraw-routing | #refreshBtn | 刷新 | GET /api/v1/admin/withdrawals/routing-stats + by-company | P1 | ☐ |

**重要**: 以上6个功能点，必须全部测试完成，少一个 = 任务未完成

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

-- 查看提现记录
SELECT w.id, w.player_id, w.amount_cents/100.0 as amount, w.status, 
       w.settlement_company_id, sc.name as company_name, w.created_at
FROM withdraws w
LEFT JOIN settlement_companies sc ON w.settlement_company_id = sc.id
ORDER BY w.created_at DESC LIMIT 20;

-- 按公司统计提现
SELECT sc.id, sc.name, 
       COUNT(w.id) as withdraw_count, 
       SUM(w.amount_cents)/100.0 as total_amount
FROM settlement_companies sc
LEFT JOIN withdraws w ON sc.id = w.settlement_company_id
GROUP BY sc.id, sc.name;

-- 按状态统计
SELECT status, COUNT(*) as count, SUM(amount_cents)/100.0 as total_amount
FROM withdraws GROUP BY status;
```

### 测试账号
- **管理员**: admin@gameLink.com / Admin2025@Pass#

---

## 五、逐个按钮测试记录

### 功能1: #pageLoad 页面加载

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 访问页面 /admin/finance/withdraw-routing
4. 观察页面加载和数据展示

**Evidence收集**:
- [ ] 截图1: 页面完整展示（统计卡片+公司统计+明细列表）
- [ ] 截图2: Network请求详情
  - 请求1: GET /api/v1/admin/withdrawals/routing-stats
  - 请求2: GET /api/v1/admin/withdrawals/by-company
  - 两个请求都应返回200
- [ ] 截图3: docker logs处理记录
  ```bash
  docker logs gamelink-backend --tail=30 | findstr "withdrawals"
  ```
- [ ] 截图4: 数据库验证（统计数据）
  ```sql
  -- 验证总提现金额
  SELECT SUM(amount_cents)/100.0 as total FROM withdraws;
  -- 验证总笔数
  SELECT COUNT(*) as count FROM withdraws;
  -- 验证已完成金额
  SELECT SUM(amount_cents)/100.0 FROM withdraws WHERE status = 'completed';
  ```
- [ ] 截图5: 统计卡片数据与数据库一致

**异常场景测试**:
- [ ] 场景A: 无提现数据时 → 预期: 统计显示0，列表显示"暂无数据"
- [ ] 场景B: 后端服务不可用 → 预期: 错误提示
- [ ] 场景C: 网络慢时 → 预期: 显示loading状态

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 功能2: #dateRangeBtn 日期范围选择

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击日期范围选择器
4. 选择日期范围: 2024-12-01 至 2024-12-31
5. 观察数据刷新

**Evidence收集**:
- [ ] 截图1: 日期范围选择器
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/withdrawals/by-company?dateFrom=2024-12-01&dateTo=2024-12-31
  - URL: GET /api/v1/admin/withdrawals/routing-stats?dateFrom=2024-12-01&dateTo=2024-12-31
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT COUNT(*), SUM(amount_cents)/100.0 
  FROM withdraws 
  WHERE created_at >= '2024-12-01' AND created_at < '2025-01-01';
  ```
- [ ] 截图5: 统计卡片和列表数据按日期过滤

**异常场景测试**:
- [ ] 场景A: 开始日期大于结束日期 → 预期: 校验失败或自动交换
- [ ] 场景B: 选择无数据的日期范围 → 预期: 显示0和空列表
- [ ] 场景C: 清空日期范围 → 预期: 显示全部数据
- [ ] 场景D: 选择未来日期 → 预期: 允许选择，返回空数据

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 功能3: #searchBtn 搜索

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 输入结算公司ID: 1
4. 选择状态: 已完成
5. 点击搜索

**Evidence收集**:
- [ ] 截图1: 搜索条件填写
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/withdrawals/by-company?settlementCompanyId=1&status=completed
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM withdraws 
  WHERE settlement_company_id = 1 AND status = 'completed'
  ORDER BY created_at DESC;
  ```
- [ ] 截图5: 列表显示过滤结果

**异常场景测试**:
- [ ] 场景A: 搜索无结果 → 预期: 显示空列表
- [ ] 场景B: 公司ID不存在 → 预期: 返回空列表
- [ ] 场景C: 清空搜索条件 → 预期: 显示全部数据

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 功能4: #monthlyReportBtn 月报

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击"月报"按钮
4. 观察报表生成

**Evidence收集**:
- [ ] 截图1: 月报按钮
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/withdrawals/routing-report?reportType=monthly&year=2024&month=12
  - Response: 报表数据
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证（当月数据）
  ```sql
  SELECT sc.name, COUNT(w.id), SUM(w.amount_cents)/100.0
  FROM withdraws w
  JOIN settlement_companies sc ON w.settlement_company_id = sc.id
  WHERE EXTRACT(YEAR FROM w.created_at) = 2024 
    AND EXTRACT(MONTH FROM w.created_at) = 12
  GROUP BY sc.id, sc.name;
  ```
- [ ] 截图5: 报表生成成功提示或下载

**异常场景测试**:
- [ ] 场景A: 当月无数据 → 预期: 生成空报表或提示
- [ ] 场景B: 后端生成报表超时 → 预期: 超时提示
- [ ] 场景C: 快速连续点击 → 预期: 防抖生效

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 功能5: #exportBtn 导出数据

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
- [ ] CSV文件名格式: withdraw_routing_YYYYMMDD.csv
- [ ] CSV包含列: ID, 陪玩师, 金额, 结算公司, 状态, 创建时间
- [ ] 数据与页面列表一致
- [ ] 金额格式正确（元，保留2位小数）

**异常场景测试**:
- [ ] 场景A: 空列表导出 → 预期: 只有表头的CSV
- [ ] 场景B: 大量数据导出 → 预期: 正常导出

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 功能6: #refreshBtn 刷新

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击刷新按钮

**Evidence收集**:
- [ ] 截图1: 刷新按钮
- [ ] 截图2: Network请求详情
  - 请求1: GET /api/v1/admin/withdrawals/routing-stats
  - 请求2: GET /api/v1/admin/withdrawals/by-company
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库当前数据
- [ ] 截图5: 统计卡片和列表数据刷新

**异常场景测试**:
- [ ] 场景A: 后端服务不可用时刷新 → 预期: 错误提示
- [ ] 场景B: 快速连续点击刷新 → 预期: 防抖生效

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

## 六、统计卡片验证

页面顶部有4个统计卡片，需验证数据准确性：

| 卡片名称 | 验证SQL | 预期值 | 实际值 | 结果 |
|---------|---------|--------|--------|------|
| 总提现金额 | `SELECT SUM(amount_cents)/100.0 FROM withdraws;` | | | ☐ |
| 总提现笔数 | `SELECT COUNT(*) FROM withdraws;` | | | ☐ |
| 已完成金额 | `SELECT SUM(amount_cents)/100.0 FROM withdraws WHERE status='completed';` | | | ☐ |
| 待处理金额 | `SELECT SUM(amount_cents)/100.0 FROM withdraws WHERE status='pending';` | | | ☐ |

## 七、按公司统计卡片验证

| 公司名称 | 验证SQL | 预期金额 | 预期笔数 | 实际值 | 结果 |
|---------|---------|---------|---------|--------|------|
| 公司A | `SELECT SUM(amount_cents)/100.0, COUNT(*) FROM withdraws WHERE settlement_company_id=1;` | | | | ☐ |
| 公司B | `SELECT SUM(amount_cents)/100.0, COUNT(*) FROM withdraws WHERE settlement_company_id=2;` | | | | ☐ |

---

## 八、全量测试完整性自查

- [ ] 所有P0功能已测试（3个）
- [ ] 所有P1功能已测试（3个）
- [ ] 每个功能提供5张截图+日志
- [ ] 每个功能测试异常场景≥3个
- [ ] 统计卡片数据验证通过
- [ ] 按公司统计数据验证通过
- [ ] 所有截图有明确的文件名
- [ ] 日志文件已打包

---

## 九、质量承诺

我承诺以上测试内容真实完整，所有功能均已按22项清单验证。

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
