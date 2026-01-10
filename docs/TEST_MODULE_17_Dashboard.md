# 测试任务单：仪表盘模块全量测试

**任务编号**: TEST-2024-M17  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/dashboard | #daysPicker | 时间范围选择 | GET /api/v1/admin/dashboard/stats | P0 | ☐ |
| /admin/dashboard | #viewAllOrders | 查看全部订单 | 跳转 /admin/biz/order | P1 | ☐ |
| #statsCards | 统计卡片 | GET /api/v1/admin/dashboard/stats | P0 | ☐ |
| #revenueTrend | 收入趋势图 | GET /api/v1/admin/dashboard/revenue-trend | P0 | ☐ |
| #userGrowth | 用户增长图 | GET /api/v1/admin/dashboard/user-growth | P0 | ☐ |
| #orderStatus | 订单状态分布 | GET /api/v1/admin/dashboard/stats | P1 | ☐ |
| #paymentStatus | 支付状态分布 | GET /api/v1/admin/dashboard/stats | P1 | ☐ |
| #recentOrders | 最新订单列表 | GET /api/v1/admin/orders?page=1&page_size=5 | P1 | ☐ |
| #topPlayers | 热门陪玩列表 | GET /api/v1/admin/players/top | P1 | ☐ |

**重要**: 以上9个功能点，必须全部测试完成

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

-- 验证统计数据
SELECT COUNT(*) as total_orders FROM orders;
SELECT COUNT(*) as total_users FROM users;
SELECT COUNT(*) as total_players FROM players;
SELECT SUM(total_price_cents) as total_revenue FROM orders WHERE status = 'completed';

-- 验证订单状态分布
SELECT status, COUNT(*) FROM orders GROUP BY status;

-- 验证热门陪玩
SELECT p.id, p.nickname, p.rating_average, p.rating_count 
FROM players p 
ORDER BY p.rating_average DESC, p.rating_count DESC 
LIMIT 5;
```

---

## 四、逐个功能点测试记录

### 功能1: #daysPicker 时间范围选择

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 切换时间范围（近7天/近30天/近90天）
4. 观察数据变化

**Evidence收集**:
- [ ] 截图1: 时间选择器
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/dashboard/stats?days=30
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证（按时间范围查询）
- [ ] 截图5: 图表数据正确更新

**异常场景测试**:
- [ ] 场景A: 切换到无数据的时间范围 → 预期: 显示空数据
- [ ] 场景B: 快速连续切换 → 预期: 防抖生效
- [ ] 场景C: 网络中断时切换 → 预期: 错误提示

**测试结果**: ☐ 通过 ☐ 失败

---

### 功能2: #statsCards 统计卡片

**测试步骤**:
1. 页面加载后观察统计卡片
2. 验证数据准确性

**Evidence收集**:
- [ ] 截图1: 四个统计卡片（总订单数、交易总额、总用户数、总陪玩数）
- [ ] 截图2: Network请求详情
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT 
    (SELECT COUNT(*) FROM orders) as total_orders,
    (SELECT SUM(total_price_cents) FROM orders WHERE payment_status = 'paid') as total_paid,
    (SELECT COUNT(*) FROM users) as total_users,
    (SELECT COUNT(*) FROM players) as total_players;
  ```
- [ ] 截图5: 数据与数据库一致

**测试结果**: ☐ 通过 ☐ 失败

---

### 功能3: #revenueTrend 收入趋势图

**测试步骤**:
1. 观察收入趋势折线图
2. 验证数据点准确性

**Evidence收集**:
- [ ] 截图1: 收入趋势图
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/dashboard/revenue-trend?days=7
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT DATE(created_at) as date, SUM(total_price_cents) as revenue
  FROM orders WHERE payment_status = 'paid'
  AND created_at >= NOW() - INTERVAL '7 days'
  GROUP BY DATE(created_at) ORDER BY date;
  ```
- [ ] 截图5: 图表数据与数据库一致

**测试结果**: ☐ 通过 ☐ 失败

---

### 功能4: #userGrowth 用户增长图

**测试步骤**:
1. 观察用户增长折线图
2. 验证数据点准确性

**Evidence收集**:
- [ ] 截图1: 用户增长图
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/dashboard/user-growth?days=7
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT DATE(created_at) as date, COUNT(*) as new_users
  FROM users WHERE created_at >= NOW() - INTERVAL '7 days'
  GROUP BY DATE(created_at) ORDER BY date;
  ```
- [ ] 截图5: 图表数据与数据库一致

**测试结果**: ☐ 通过 ☐ 失败

---

### 功能5-6: 状态分布饼图

**订单状态分布 #orderStatus**:
- [ ] 饼图正确显示各状态占比
- [ ] 数据与数据库一致

**支付状态分布 #paymentStatus**:
- [ ] 饼图正确显示各状态占比
- [ ] 数据与数据库一致

---

### 功能7-9: 列表数据

**最新订单 #recentOrders**:
- [ ] 显示最近5条订单
- [ ] 点击"查看全部"跳转正确

**热门陪玩 #topPlayers**:
- [ ] 显示评分最高的5个陪玩师
- [ ] 排名、评分数据正确

---

## 五、全量测试完整性自查

- [ ] 所有P0功能点已测试（5个）
- [ ] 所有P1功能点已测试（4个）
- [ ] 每个功能点提供5张截图+日志
- [ ] 每个功能点测试异常场景≥3个

---

## 六、质量承诺

我承诺以上测试内容真实完整，所有功能点均已验证。

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
