# 测试任务单：用户行为分析模块全量测试

**任务编号**: TEST-2024-M06  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/monitor/user-behavior | #loadStats | 页面加载-核心指标 | GET /api/v1/admin/users/behavior/stats | P0 | ☐ |
| /admin/monitor/user-behavior | #loadTrend | 页面加载-趋势数据 | GET /api/v1/admin/users/behavior/trend | P0 | ☐ |
| /admin/monitor/user-behavior | #loadDistribution | 页面加载-分布数据 | GET /api/v1/admin/users/behavior/distribution | P0 | ☐ |
| /admin/monitor/user-behavior | #daySelect | 天数选择(7/14/30天) | GET /api/v1/admin/users/behavior/trend?days=X | P1 | ☐ |

**重要**: 以上4个功能点，必须全部测试完成，少一个 = 任务未完成

---

## 二、测试标准（参考22项清单）

每个功能必须提供：
1. ✅ 功能静态截图（Evidence-01）
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

-- 查看用户登录历史
SELECT DATE(created_at) as date, COUNT(DISTINCT user_id) as dau
FROM user_login_histories
WHERE created_at > NOW() - INTERVAL '30 days'
GROUP BY DATE(created_at)
ORDER BY date DESC;

-- 查看用户行为记录
SELECT action_type, COUNT(*) as count
FROM user_behaviors
WHERE created_at > NOW() - INTERVAL '7 days'
GROUP BY action_type;

-- 查看用户总数和新增
SELECT 
  COUNT(*) as total_users,
  COUNT(*) FILTER (WHERE created_at > NOW() - INTERVAL '1 day') as new_today,
  COUNT(*) FILTER (WHERE created_at > NOW() - INTERVAL '30 days') as new_month
FROM users;
```

### 测试账号
- **管理员**: 使用 `.env` 中的 `SUPER_ADMIN_EMAIL` / `SUPER_ADMIN_PASSWORD`

---

## 五、逐个功能测试记录

### 功能1: #loadStats 页面加载-核心指标

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 访问页面 /admin/monitor/user-behavior
4. 观察核心指标卡片加载

**Evidence收集**:
- [ ] 截图1: 核心指标卡片（DAU、MAU、平均在线时长、人均消费、新增用户、活跃率）
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/users/behavior/stats
  - Response: `{"dau":X,"mau":X,"avgOnlineTime":X,"avgSpending":X,"newUsers":X,"activeRate":X}`
  - Status: 200
- [ ] 截图3: docker logs处理记录
  ```bash
  docker logs gamelink-backend --tail=20 | findstr "behavior"
  ```
- [ ] 截图4: 数据库验证
  ```sql
  -- DAU验证
  SELECT COUNT(DISTINCT user_id) as dau 
  FROM user_login_histories 
  WHERE DATE(created_at) = CURRENT_DATE;
  
  -- MAU验证
  SELECT COUNT(DISTINCT user_id) as mau 
  FROM user_login_histories 
  WHERE created_at > NOW() - INTERVAL '30 days';
  
  -- 新增用户验证
  SELECT COUNT(*) as new_users 
  FROM users 
  WHERE DATE(created_at) = CURRENT_DATE;
  ```
- [ ] 截图5: 指标数据与数据库一致

**异常场景测试**:
- [ ] 场景A: 无用户数据时 → 预期: 所有指标显示0
- [ ] 场景B: 后端服务不可用 → 预期: 错误提示，显示loading或默认值
- [ ] 场景C: 网络慢时 → 预期: 显示Spin loading状态

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 功能2: #loadTrend 页面加载-趋势数据

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 观察趋势数据表格加载（默认7天）

**Evidence收集**:
- [ ] 截图1: 趋势数据表格（日期、DAU、新增用户、订单数）
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/users/behavior/trend?days=7
  - Response: `[{"date":"2024-12-18","dau":X,"newUsers":X,"orders":X},...]`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  -- 最近7天DAU趋势
  SELECT DATE(created_at) as date, COUNT(DISTINCT user_id) as dau
  FROM user_login_histories
  WHERE created_at > NOW() - INTERVAL '7 days'
  GROUP BY DATE(created_at)
  ORDER BY date DESC;
  
  -- 最近7天新增用户
  SELECT DATE(created_at) as date, COUNT(*) as new_users
  FROM users
  WHERE created_at > NOW() - INTERVAL '7 days'
  GROUP BY DATE(created_at)
  ORDER BY date DESC;
  
  -- 最近7天订单数
  SELECT DATE(created_at) as date, COUNT(*) as orders
  FROM orders
  WHERE created_at > NOW() - INTERVAL '7 days'
  GROUP BY DATE(created_at)
  ORDER BY date DESC;
  ```
- [ ] 截图5: 趋势数据与数据库一致

**异常场景测试**:
- [ ] 场景A: 无趋势数据时 → 预期: 显示"暂无数据"
- [ ] 场景B: 部分日期无数据 → 预期: 该日期显示0或不显示

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 功能3: #loadDistribution 页面加载-分布数据

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 观察分布数据卡片加载（地域分布、年龄分布、设备分布）

**Evidence收集**:
- [ ] 截图1: 三个分布卡片（地域、年龄、设备）
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/users/behavior/distribution
  - Response: 
    ```json
    {
      "regions": [{"name":"北京","count":X},...],
      "ageGroups": [{"range":"18-25","count":X},...],
      "devices": [{"type":"iOS","count":X},...]
    }
    ```
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  -- 地域分布（如果有region字段）
  SELECT region, COUNT(*) as count 
  FROM users 
  WHERE region IS NOT NULL
  GROUP BY region 
  ORDER BY count DESC LIMIT 5;
  
  -- 设备分布（如果有device字段）
  SELECT device_type, COUNT(*) as count 
  FROM user_login_histories 
  GROUP BY device_type;
  ```
- [ ] 截图5: 分布数据与数据库一致

**异常场景测试**:
- [ ] 场景A: 无分布数据时 → 预期: 显示"暂无数据"
- [ ] 场景B: 只有部分分布数据 → 预期: 有数据的显示，无数据的显示空

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 功能4: #daySelect 天数选择(7/14/30天)

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击天数选择下拉框
4. 依次选择"最近7天"、"最近14天"、"最近30天"
5. 观察趋势数据变化

**Evidence收集**:
- [ ] 截图1: 天数选择下拉框
- [ ] 截图2: Network请求详情（选择14天时）
  - URL: GET /api/v1/admin/users/behavior/trend?days=14
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证（14天数据）
  ```sql
  SELECT DATE(created_at) as date, COUNT(DISTINCT user_id) as dau
  FROM user_login_histories
  WHERE created_at > NOW() - INTERVAL '14 days'
  GROUP BY DATE(created_at)
  ORDER BY date DESC;
  ```
- [ ] 截图5: 趋势表格行数变化（7天→14天→30天）

**异常场景测试**:
- [ ] 场景A: 快速切换天数 → 预期: 防抖生效，只发送最后一次请求
- [ ] 场景B: 选择30天但只有7天数据 → 预期: 只显示有数据的7天
- [ ] 场景C: 后端返回错误 → 预期: 错误提示，保持原数据

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

## 六、核心指标验证

| 指标名称 | 验证SQL | 预期值 | 实际值 | 结果 |
|---------|---------|--------|--------|------|
| DAU | `SELECT COUNT(DISTINCT user_id) FROM user_login_histories WHERE DATE(created_at)=CURRENT_DATE;` | | | ☐ |
| MAU | `SELECT COUNT(DISTINCT user_id) FROM user_login_histories WHERE created_at>NOW()-INTERVAL '30 days';` | | | ☐ |
| 新增用户 | `SELECT COUNT(*) FROM users WHERE DATE(created_at)=CURRENT_DATE;` | | | ☐ |
| 活跃率 | DAU / 总用户数 * 100% | | | ☐ |

---

## 七、趋势数据验证（最近3天示例）

| 日期 | DAU(页面) | DAU(数据库) | 新增(页面) | 新增(数据库) | 订单(页面) | 订单(数据库) | 结果 |
|------|----------|------------|-----------|-------------|-----------|-------------|------|
| 12-18 | | | | | | | ☐ |
| 12-17 | | | | | | | ☐ |
| 12-16 | | | | | | | ☐ |

---

## 八、全量测试完整性自查

- [ ] 所有P0功能已测试（3个）
- [ ] 所有P1功能已测试（1个）
- [ ] 每个功能提供5张截图+日志
- [ ] 每个功能测试异常场景≥2个
- [ ] 核心指标数据验证通过
- [ ] 趋势数据验证通过
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
