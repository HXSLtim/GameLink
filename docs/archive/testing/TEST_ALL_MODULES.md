# 全量测试任务单 - 管理后台各模块

**任务编号**: TEST-2024-FULL  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 环境准备

### 1. 容器状态检查
```bash
docker compose -f docker-compose.prod.yml ps
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
```

**预期结果**: 所有容器状态为 "Up (healthy)"
- gamelink-backend: 8081→8080
- gamelink-frontend: 80, 443
- gamelink-postgres: 5432
- gamelink-redis: 6379

### 2. 测试账号
- **管理员**: 使用 `.env` 中的 `SUPER_ADMIN_EMAIL` / `SUPER_ADMIN_PASSWORD`（由 `docker-compose.prod.yml` 注入后端容器）
- **测试地址**: http://localhost/admin

### 3. 日志监控命令
```bash
# 后端日志
docker logs -f gamelink-backend --tail=50

# 检查错误
docker logs gamelink-backend | findstr /i "error"

# 数据库连接
docker exec -it gamelink-postgres psql -U gamelink -d gamelink
```

---

# 模块一：用户标签管理 (UserTag)

**页面路径**: `/admin/sys/user-tag`

## 测试范围（必须100%覆盖）

| 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|--------|---------|---------|--------|---------|
| #createBtn | 新增标签 | POST /api/v1/admin/user-tags | P0 | ☐ |
| #editBtn | 编辑 | PUT /api/v1/admin/user-tags/:id | P0 | ☐ |
| #deleteBtn | 删除 | DELETE /api/v1/admin/user-tags/:id | P0 | ☐ |
| #viewUsersBtn | 查看用户 | GET /api/v1/admin/user-tags/:id/users | P1 | ☐ |
| #searchBtn | 搜索 | 前端过滤 | P1 | ☐ |
| #exportBtn | 导出数据 | 前端导出CSV | P1 | ☐ |
| #refreshBtn | 刷新 | GET /api/v1/admin/user-tags | P1 | ☐ |

## 数据库验证SQL
```sql
-- 查看所有标签
SELECT * FROM user_tags ORDER BY id DESC LIMIT 10;

-- 查看标签用户关联
SELECT ut.name, COUNT(utr.user_id) as user_count 
FROM user_tags ut 
LEFT JOIN user_tag_relations utr ON ut.id = utr.tag_id 
GROUP BY ut.id, ut.name;
```

## 按钮测试记录

### 按钮1: 新增标签
**测试步骤**:
1. 点击"新增标签"按钮
2. 填写：名称=测试标签, 颜色=#ff4d4f, 描述=测试用
3. 点击确定

**Evidence收集**:
- [ ] 截图1: 弹窗表单
- [ ] 截图2: Network请求 POST /api/v1/admin/user-tags
- [ ] 截图3: docker logs gamelink-backend 处理记录
- [ ] 截图4: 数据库 SELECT * FROM user_tags WHERE name='测试标签'
- [ ] 截图5: 列表刷新显示新标签

**异常场景测试**:
- [ ] 场景A: 名称为空提交
- [ ] 场景B: 重复名称提交
- [ ] 场景C: 快速连续点击确定

---

# 模块二：分流规则管理 (RoutingRule)

**页面路径**: `/admin/biz/routing-rule`

## 测试范围

| 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|--------|---------|---------|--------|---------|
| #createBtn | 新增规则 | POST /api/v1/admin/routing-rules | P0 | ☐ |
| #editBtn | 编辑 | PUT /api/v1/admin/routing-rules/:id | P0 | ☐ |
| #deleteBtn | 删除 | DELETE /api/v1/admin/routing-rules/:id | P0 | ☐ |
| #toggleBtn | 启用/禁用 | POST /api/v1/admin/routing-rules/:id/toggle | P0 | ☐ |
| #historyBtn | 历史 | GET /api/v1/admin/routing-rules/:id/history | P1 | ☐ |
| #testBtn | 测试规则 | POST /api/v1/admin/routing-rules/test | P1 | ☐ |
| #exportBtn | 导出数据 | 前端导出CSV | P1 | ☐ |

## 数据库验证SQL
```sql
-- 查看分流规则
SELECT * FROM routing_rules ORDER BY priority ASC;

-- 查看规则历史
SELECT * FROM routing_rule_histories WHERE rule_id = 1 ORDER BY created_at DESC;
```

## 按钮测试记录

### 按钮1: 新增规则
**测试步骤**:
1. 点击"新增规则"
2. 填写：名称=测试规则, 优先级=10, 目标主体ID=1, 条件类型=金额范围, 条件值=100-500
3. 点击确定

**Evidence收集**:
- [ ] 截图1-5 (同上模板)

**异常场景测试**:
- [ ] 场景A: 优先级重复
- [ ] 场景B: 目标主体ID不存在
- [ ] 场景C: 条件值格式错误

---

# 模块三：结算公司管理 (SettlementCompany)

**页面路径**: `/admin/finance/settlement-company`

## 测试范围

| 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|--------|---------|---------|--------|---------|
| #createBtn | 新增公司 | POST /api/v1/admin/settlement-companies | P0 | ☐ |
| #editBtn | 编辑 | PUT /api/v1/admin/settlement-companies/:id | P0 | ☐ |
| #toggleBtn | 启用/禁用 | POST /api/v1/admin/settlement-companies/:id/toggle | P0 | ☐ |
| #exportBtn | 导出数据 | 前端导出CSV | P1 | ☐ |

## 数据库验证SQL
```sql
-- 查看结算公司
SELECT * FROM settlement_companies ORDER BY id DESC;

-- 查看公司关联陪玩师数
SELECT sc.name, COUNT(p.id) as player_count 
FROM settlement_companies sc 
LEFT JOIN players p ON sc.id = p.settlement_company_id 
GROUP BY sc.id, sc.name;
```

## 按钮测试记录

### 按钮1: 新增公司
**测试步骤**:
1. 点击"新增公司"
2. 填写：公司名称=测试公司, 信用代码=91110000MA00ABCD12, 联系人=张三, 联系电话=13800138000
3. 点击确定

**Evidence收集**:
- [ ] 截图1-5 (同上模板)

**异常场景测试**:
- [ ] 场景A: 信用代码不是18位
- [ ] 场景B: 信用代码重复
- [ ] 场景C: 必填字段为空

---

# 模块四：排行榜抽成配置 (RankingCommission)

**页面路径**: `/admin/finance/ranking-commission`

## 测试范围

| 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|--------|---------|---------|--------|---------|
| #createBtn | 新增配置 | POST /api/v1/admin/ranking-commission/configs | P0 | ☐ |
| #editBtn | 编辑 | PUT /api/v1/admin/ranking-commission/configs/:id | P0 | ☐ |
| #deleteBtn | 删除 | DELETE /api/v1/admin/ranking-commission/configs/:id | P0 | ☐ |
| #addRuleBtn | 添加规则 | 表单内操作 | P1 | ☐ |
| #exportBtn | 导出数据 | 前端导出CSV | P1 | ☐ |

## 数据库验证SQL
```sql
-- 查看抽成配置
SELECT * FROM ranking_commission_configs ORDER BY id DESC;

-- 查看配置详情（含规则）
SELECT id, name, ranking_type, month, rules, is_active FROM ranking_commission_configs;
```

## 按钮测试记录

### 按钮1: 新增配置
**测试步骤**:
1. 点击"新增配置"
2. 填写：名称=2024年12月收入排行, 排行类型=收入排行, 月份=2024-12
3. 添加规则：1-10名 5%, 11-50名 3%
4. 点击确定

**Evidence收集**:
- [ ] 截图1-5 (同上模板)

**异常场景测试**:
- [ ] 场景A: 同月份同类型重复创建
- [ ] 场景B: 规则排名范围重叠
- [ ] 场景C: 抽成比例超过100%

---

# 模块五：提现分流统计 (WithdrawRouting)

**页面路径**: `/admin/finance/withdraw-routing`

## 测试范围

| 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|--------|---------|---------|--------|---------|
| #searchBtn | 搜索 | GET /api/v1/admin/withdrawals/by-company | P0 | ☐ |
| #dateRangeBtn | 日期筛选 | GET /api/v1/admin/withdrawals/by-company | P0 | ☐ |
| #monthlyReportBtn | 月报 | GET /api/v1/admin/withdrawals/routing-report | P1 | ☐ |
| #exportBtn | 导出数据 | 前端导出CSV | P1 | ☐ |
| #refreshBtn | 刷新 | GET /api/v1/admin/withdrawals/routing-stats | P1 | ☐ |

## 数据库验证SQL
```sql
-- 查看提现记录
SELECT w.*, sc.name as company_name 
FROM withdraws w 
LEFT JOIN settlement_companies sc ON w.settlement_company_id = sc.id 
ORDER BY w.created_at DESC LIMIT 10;

-- 按公司统计
SELECT sc.name, COUNT(w.id) as count, SUM(w.amount_cents)/100.0 as total_amount 
FROM withdraws w 
JOIN settlement_companies sc ON w.settlement_company_id = sc.id 
GROUP BY sc.id, sc.name;
```

## 按钮测试记录

### 按钮1: 日期筛选
**测试步骤**:
1. 选择日期范围：2024-12-01 至 2024-12-31
2. 观察列表和统计数据变化

**Evidence收集**:
- [ ] 截图1-5 (同上模板)

**异常场景测试**:
- [ ] 场景A: 开始日期大于结束日期
- [ ] 场景B: 选择无数据的日期范围
- [ ] 场景C: 清空日期后刷新

---

# 模块六：用户行为分析 (UserBehavior)

**页面路径**: `/admin/monitor/user-behavior`

## 测试范围

| 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|--------|---------|---------|--------|---------|
| #loadStats | 页面加载 | GET /api/v1/admin/users/behavior/stats | P0 | ☐ |
| #loadTrend | 趋势数据 | GET /api/v1/admin/users/behavior/trend | P0 | ☐ |
| #loadDistribution | 分布数据 | GET /api/v1/admin/users/behavior/distribution | P0 | ☐ |
| #daySelect | 天数选择 | GET /api/v1/admin/users/behavior/trend?days=X | P1 | ☐ |

## 数据库验证SQL
```sql
-- 用户行为统计
SELECT action_type, COUNT(*) as count 
FROM user_behaviors 
WHERE created_at > NOW() - INTERVAL '7 days' 
GROUP BY action_type;

-- 登录历史
SELECT DATE(created_at) as date, COUNT(DISTINCT user_id) as dau 
FROM user_login_histories 
WHERE created_at > NOW() - INTERVAL '7 days' 
GROUP BY DATE(created_at) 
ORDER BY date DESC;
```

## 按钮测试记录

### 按钮1: 天数选择
**测试步骤**:
1. 选择"最近7天"
2. 切换到"最近14天"
3. 切换到"最近30天"
4. 观察趋势数据变化

**Evidence收集**:
- [ ] 截图1-5 (同上模板)

**异常场景测试**:
- [ ] 场景A: 后端服务不可用时的错误提示
- [ ] 场景B: 数据为空时的展示

---

# 模块七：纠纷管理 (Dispute)

**页面路径**: `/admin/biz/dispute`

## 测试范围

| 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|--------|---------|---------|--------|---------|
| #searchBtn | 搜索 | GET /api/v1/admin/disputes | P0 | ☐ |
| #filterBtn | 状态筛选 | GET /api/v1/admin/disputes?status=X | P0 | ☐ |
| #detailBtn | 详情 | GET /api/v1/admin/disputes/:id | P0 | ☐ |
| #assignBtn | 分配 | POST /api/v1/admin/disputes/:id/assign | P0 | ☐ |
| #resolveBtn | 处理 | POST /api/v1/admin/disputes/:id/resolve | P0 | ☐ |
| #rollbackBtn | 回滚 | POST /api/v1/admin/disputes/:id/rollback | P1 | ☐ |
| #exportBtn | 导出数据 | 前端导出CSV | P1 | ☐ |
| #statsLoad | 统计加载 | GET /api/v1/admin/disputes/stats | P0 | ☐ |

## 数据库验证SQL
```sql
-- 查看纠纷列表
SELECT od.*, o.order_no, u.name as user_name 
FROM order_disputes od 
LEFT JOIN orders o ON od.order_id = o.id 
LEFT JOIN users u ON od.user_id = u.id 
ORDER BY od.created_at DESC LIMIT 10;

-- 纠纷统计
SELECT status, COUNT(*) as count FROM order_disputes GROUP BY status;
```

## 按钮测试记录

### 按钮1: 分配纠纷
**测试步骤**:
1. 找到状态为"待处理"的纠纷
2. 点击"分配"按钮
3. 选择处理人
4. 点击确定

**Evidence收集**:
- [ ] 截图1-5 (同上模板)

**异常场景测试**:
- [ ] 场景A: 分配给不存在的用户
- [ ] 场景B: 重复分配同一纠纷
- [ ] 场景C: 分配已解决的纠纷

---

## 全量测试完整性自查

- [ ] 模块一：用户标签管理 - 7个按钮全部测试
- [ ] 模块二：分流规则管理 - 7个按钮全部测试
- [ ] 模块三：结算公司管理 - 4个按钮全部测试
- [ ] 模块四：排行榜抽成配置 - 5个按钮全部测试
- [ ] 模块五：提现分流统计 - 5个按钮全部测试
- [ ] 模块六：用户行为分析 - 4个按钮全部测试
- [ ] 模块七：纠纷管理 - 8个按钮全部测试

**总计**: 40个按钮/功能点

---

## 质量承诺

我承诺以上测试内容真实完整，所有按钮均已按22项清单验证。如有遗漏，愿意承担测试质量责任。

**测试人签字**: ___________  
**日期**: ___________

---

## 组长审核意见

**审核结果**: ☐ 通过 ☐ 打回重做  
**打回原因**: （如有）  
**审核人**: ___________  
**日期**: ___________

---

**文档版本**: v1.0  
**发布日期**: 2024-12-18  
