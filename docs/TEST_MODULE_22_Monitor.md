# 测试任务单：系统监控模块全量测试

**任务编号**: TEST-2024-022  
**测试环境**: Docker 生产环境 (docker-compose.prod.yml)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、模块概述

系统监控模块包含三个子页面：
1. **实时监控 (Realtime)** - 系统运行状态、在线用户、订单队列、异常告警
2. **KPI 仪表板 (KPI)** - 业务 KPI 指标概览、趋势图表和目标管理
3. **运营分析 (Analytics)** - 用户活跃度、留存率、付费分析、转化漏斗

**页面路径**:
- `/admin/monitor/realtime` - 实时监控
- `/admin/monitor/kpi` - KPI 仪表板
- `/admin/monitor/analytics` - 运营分析

**关联 API**:
- `GET /api/v1/admin/monitor/system-status` - 系统状态
- `GET /api/v1/admin/monitor/online-users` - 在线用户
- `GET /api/v1/admin/monitor/order-queue` - 订单队列
- `GET /api/v1/admin/monitor/alerts` - 告警列表
- `PUT /api/v1/admin/monitor/alerts/:id/read` - 标记告警已读
- `PUT /api/v1/admin/monitor/alerts/batch-read` - 批量标记已读
- `GET /api/v1/admin/monitor/kpi/overview` - KPI 概览
- `GET /api/v1/admin/monitor/kpi/trend/:metric` - KPI 趋势
- `GET /api/v1/admin/monitor/kpi/targets` - KPI 目标列表
- `POST /api/v1/admin/monitor/kpi/targets` - 创建 KPI 目标
- `PUT /api/v1/admin/monitor/kpi/targets/:id` - 更新 KPI 目标
- `DELETE /api/v1/admin/monitor/kpi/targets/:id` - 删除 KPI 目标
- `GET /api/v1/admin/monitor/analytics/active-users` - 活跃用户
- `GET /api/v1/admin/monitor/analytics/retention` - 留存率
- `GET /api/v1/admin/monitor/analytics/payment` - 付费分析
- `GET /api/v1/admin/monitor/analytics/funnel` - 转化漏斗
- `WebSocket /ws/monitor` - 实时监控数据推送

---

## 二、测试范围（必须 100% 覆盖）

### 2.1 实时监控页面 (Realtime)

| 页面路径 | 按钮/交互 ID | 按钮文案 | 关联 API | 优先级 | 测试状态 |
|---------|-------------|---------|---------|--------|---------|
| /admin/monitor/realtime | #refreshBtn | 刷新 | GET /monitor/* | P0 | ☐ |
| /admin/monitor/realtime | #markReadBtn | 标记已读 | PUT /alerts/:id/read | P0 | ☐ |
| /admin/monitor/realtime | #markAllReadBtn | 全部标记已读 | PUT /alerts/batch-read | P0 | ☐ |
| /admin/monitor/realtime | WebSocket | 实时数据连接 | WS /ws/monitor | P0 | ☐ |

### 2.2 KPI 仪表板页面 (KPI)

| 页面路径 | 按钮/交互 ID | 按钮文案 | 关联 API | 优先级 | 测试状态 |
|---------|-------------|---------|---------|--------|---------|
| /admin/monitor/kpi | #periodSelect | 周期选择 | GET /kpi/overview | P0 | ☐ |
| /admin/monitor/kpi | #dateRangePicker | 日期范围 | GET /kpi/overview | P1 | ☐ |
| /admin/monitor/kpi | #compareSelect | 对比选择 | GET /kpi/overview | P1 | ☐ |
| /admin/monitor/kpi | #metricCard | 指标卡片点击 | GET /kpi/trend/:metric | P0 | ☐ |
| /admin/monitor/kpi | #createTargetBtn | 新建目标 | POST /kpi/targets | P0 | ☐ |
| /admin/monitor/kpi | #editTargetBtn | 编辑目标 | PUT /kpi/targets/:id | P0 | ☐ |
| /admin/monitor/kpi | #deleteTargetBtn | 删除目标 | DELETE /kpi/targets/:id | P0 | ☐ |
| /admin/monitor/kpi | #targetModalOk | 目标弹窗确认 | POST/PUT /kpi/targets | P0 | ☐ |
| /admin/monitor/kpi | #targetModalCancel | 目标弹窗取消 | - | P2 | ☐ |

### 2.3 运营分析页面 (Analytics)

| 页面路径 | 按钮/交互 ID | 按钮文案 | 关联 API | 优先级 | 测试状态 |
|---------|-------------|---------|---------|--------|---------|
| /admin/monitor/analytics | #dateRangePicker | 日期范围选择 | GET /analytics/* | P0 | ☐ |
| /admin/monitor/analytics | #granularitySelect | 粒度选择 | GET /analytics/* | P0 | ☐ |

**重要**: 以上 15 个按钮/交互，必须全部测试完成，少一个 = 任务未完成

---

## 三、Docker 环境检查

### 3.1 测试前环境检查

```powershell
# 检查容器状态
docker compose -f docker-compose.prod.yml ps

# 检查容器健康状态
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# 检查后端日志
docker logs gamelink-backend --tail=20

# 检查 Redis 连接（监控数据缓存）
docker exec -it gamelink-redis redis-cli -a TvXYJ305HNhsnIpQ PING
```

**预期结果**: 所有容器状态为 "Up"，PING 返回 PONG

### 3.2 数据库表结构确认

```sql
-- 连接数据库
docker exec -it gamelink-postgres psql -U gamelink -d gamelink

-- 检查 KPI 目标表
\d kpi_targets

-- 检查告警表
\d alerts

-- 检查系统监控相关表
\d system_metrics
```

---

## 四、逐个按钮测试记录

---

### 子模块 A：实时监控 (Realtime)

---

### 按钮 A1: #refreshBtn 刷新

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0 > $null`
2. 打开开发者工具 → Network 面板
3. 访问 `/admin/monitor/realtime`
4. 点击「刷新」按钮
5. 监控所有容器日志

**Evidence 收集**:
- [ ] 截图1: 按钮点击前页面状态
- [ ] 截图2: Network 请求详情（system-status, online-users, order-queue, alerts）
- [ ] 截图3: docker logs gamelink-backend 处理记录
- [ ] 截图4: 数据刷新后页面状态
- [ ] 截图5: Console 无错误

**数据库验证 SQL**:
```sql
-- 验证系统指标记录
SELECT * FROM system_metrics ORDER BY created_at DESC LIMIT 5;
```

**异常场景测试**:
- [ ] 场景A: 后端服务重启中点击刷新
- [ ] 场景B: Redis 缓存不可用时刷新
- [ ] 场景C: 快速连续点击 5 次

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮 A2: #markReadBtn 标记已读

**测试步骤**:
1. 确保有未读告警存在
2. 打开开发者工具 → Network 面板
3. 点击某条告警的「标记已读」按钮
4. 验证请求和响应

**Evidence 收集**:
- [ ] 截图1: 未读告警列表
- [ ] 截图2: Network 请求详情（PUT /alerts/:id/read）
- [ ] 截图3: docker logs gamelink-backend 处理记录
- [ ] 截图4: 数据库验证告警状态
- [ ] 截图5: 告警状态更新后页面

**数据库验证 SQL**:
```sql
-- 验证告警已读状态
SELECT id, title, is_read, updated_at FROM alerts WHERE id = '<告警ID>';
```

**异常场景测试**:
- [ ] 场景A: 告警已被其他管理员标记已读
- [ ] 场景B: 告警 ID 不存在
- [ ] 场景C: 无权限标记告警

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮 A3: #markAllReadBtn 全部标记已读

**测试步骤**:
1. 确保有多条未读告警
2. 打开开发者工具 → Network 面板
3. 点击「全部标记已读」按钮
4. 验证批量更新

**Evidence 收集**:
- [ ] 截图1: 多条未读告警
- [ ] 截图2: Network 请求详情（PUT /alerts/batch-read）
- [ ] 截图3: docker logs gamelink-backend 处理记录
- [ ] 截图4: 数据库验证批量更新
- [ ] 截图5: 所有告警已读状态

**数据库验证 SQL**:
```sql
-- 验证所有告警已读
SELECT COUNT(*) as unread_count FROM alerts WHERE is_read = false;
```

**异常场景测试**:
- [ ] 场景A: 无未读告警时点击
- [ ] 场景B: 批量更新部分失败
- [ ] 场景C: 并发标记已读

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 交互 A4: WebSocket 实时数据连接

**测试步骤**:
1. 打开开发者工具 → Network → WS 面板
2. 访问 `/admin/monitor/realtime`
3. 观察 WebSocket 连接建立
4. 验证实时数据推送

**Evidence 收集**:
- [ ] 截图1: WebSocket 连接状态（Connected）
- [ ] 截图2: WS 消息列表（system_status, online_users 等）
- [ ] 截图3: docker logs gamelink-backend WebSocket 日志
- [ ] 截图4: 页面实时数据更新
- [ ] 截图5: 连接断开重连测试

**异常场景测试**:
- [ ] 场景A: 网络断开后自动重连
- [ ] 场景B: 后端重启后重连
- [ ] 场景C: Token 过期时 WebSocket 行为

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 子模块 B：KPI 仪表板 (KPI)

---

### 按钮 B1: #periodSelect 周期选择

**测试步骤**:
1. 打开开发者工具 → Network 面板
2. 访问 `/admin/monitor/kpi`
3. 切换周期选择（今日/本周/本月/自定义）
4. 验证数据重新加载

**Evidence 收集**:
- [ ] 截图1: 默认周期状态
- [ ] 截图2: Network 请求详情（GET /kpi/overview?period=xxx）
- [ ] 截图3: docker logs gamelink-backend 处理记录
- [ ] 截图4: 不同周期数据对比
- [ ] 截图5: 图表数据更新

**异常场景测试**:
- [ ] 场景A: 快速切换周期
- [ ] 场景B: 无数据周期显示
- [ ] 场景C: 自定义日期范围过大

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮 B2: #metricCard 指标卡片点击

**测试步骤**:
1. 打开开发者工具 → Network 面板
2. 点击不同的 KPI 指标卡片（GMV/订单数/新用户等）
3. 验证趋势图表更新

**Evidence 收集**:
- [ ] 截图1: 默认选中指标
- [ ] 截图2: Network 请求详情（GET /kpi/trend/:metric）
- [ ] 截图3: docker logs gamelink-backend 处理记录
- [ ] 截图4: 趋势图表数据更新
- [ ] 截图5: 卡片选中状态变化

**异常场景测试**:
- [ ] 场景A: 快速切换指标
- [ ] 场景B: 指标无数据
- [ ] 场景C: 趋势数据加载超时

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮 B3: #createTargetBtn 新建目标

**测试步骤**:
1. 打开开发者工具 → Network 面板
2. 点击「新建目标」按钮
3. 填写目标表单（指标、周期类型、目标值、有效期）
4. 点击确认提交

**Evidence 收集**:
- [ ] 截图1: 新建目标弹窗
- [ ] 截图2: Network 请求详情（POST /kpi/targets）
- [ ] 截图3: docker logs gamelink-backend 处理记录
- [ ] 截图4: 数据库新增记录
- [ ] 截图5: 目标列表更新

**数据库验证 SQL**:
```sql
-- 验证新建目标
SELECT * FROM kpi_targets ORDER BY created_at DESC LIMIT 1;
```

**异常场景测试**:
- [ ] 场景A: 必填字段为空提交
- [ ] 场景B: 目标值为负数
- [ ] 场景C: 有效期结束日期早于开始日期

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮 B4: #editTargetBtn 编辑目标

**测试步骤**:
1. 在目标列表中点击编辑按钮
2. 修改目标信息
3. 点击确认保存

**Evidence 收集**:
- [ ] 截图1: 编辑弹窗（数据回填）
- [ ] 截图2: Network 请求详情（PUT /kpi/targets/:id）
- [ ] 截图3: docker logs gamelink-backend 处理记录
- [ ] 截图4: 数据库更新记录
- [ ] 截图5: 目标列表更新

**数据库验证 SQL**:
```sql
-- 验证目标更新
SELECT * FROM kpi_targets WHERE id = <目标ID>;
```

**异常场景测试**:
- [ ] 场景A: 目标已被删除
- [ ] 场景B: 并发编辑同一目标
- [ ] 场景C: 修改为无效值

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮 B5: #deleteTargetBtn 删除目标

**测试步骤**:
1. 在目标列表中点击删除按钮
2. 确认删除弹窗
3. 验证删除结果

**Evidence 收集**:
- [ ] 截图1: 删除确认弹窗
- [ ] 截图2: Network 请求详情（DELETE /kpi/targets/:id）
- [ ] 截图3: docker logs gamelink-backend 处理记录
- [ ] 截图4: 数据库记录删除
- [ ] 截图5: 目标列表更新

**数据库验证 SQL**:
```sql
-- 验证目标删除（软删除检查）
SELECT * FROM kpi_targets WHERE id = <目标ID>;
```

**异常场景测试**:
- [ ] 场景A: 目标已被删除
- [ ] 场景B: 取消删除操作
- [ ] 场景C: 无权限删除

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 子模块 C：运营分析 (Analytics)

---

### 按钮 C1: #dateRangePicker 日期范围选择

**测试步骤**:
1. 打开开发者工具 → Network 面板
2. 访问 `/admin/monitor/analytics`
3. 修改日期范围
4. 验证所有图表数据更新

**Evidence 收集**:
- [ ] 截图1: 默认日期范围
- [ ] 截图2: Network 请求详情（4 个并行请求）
- [ ] 截图3: docker logs gamelink-backend 处理记录
- [ ] 截图4: 图表数据更新
- [ ] 截图5: 统计卡片数据更新

**异常场景测试**:
- [ ] 场景A: 选择未来日期
- [ ] 场景B: 日期范围超过 1 年
- [ ] 场景C: 无数据日期范围

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮 C2: #granularitySelect 粒度选择

**测试步骤**:
1. 打开开发者工具 → Network 面板
2. 切换粒度（按天/按周/按月）
3. 验证数据聚合变化

**Evidence 收集**:
- [ ] 截图1: 默认粒度（按天）
- [ ] 截图2: Network 请求详情（granularity 参数）
- [ ] 截图3: docker logs gamelink-backend 处理记录
- [ ] 截图4: 趋势图表数据点变化
- [ ] 截图5: 留存矩阵数据变化

**异常场景测试**:
- [ ] 场景A: 快速切换粒度
- [ ] 场景B: 按月粒度数据不足
- [ ] 场景C: 数据聚合计算超时

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

## 五、特殊测试场景

### 5.1 WebSocket 连接稳定性测试

```powershell
# 模拟后端重启
docker restart gamelink-backend

# 观察前端 WebSocket 重连行为
# 预期：自动重连，数据恢复
```

### 5.2 大数据量性能测试

```sql
-- 插入大量测试告警
INSERT INTO alerts (title, message, level, type, is_read, created_at)
SELECT 
  '测试告警 ' || generate_series,
  '测试消息内容',
  (ARRAY['high', 'medium', 'low'])[floor(random() * 3 + 1)],
  (ARRAY['system', 'business', 'security'])[floor(random() * 3 + 1)],
  false,
  NOW() - (generate_series || ' minutes')::interval
FROM generate_series(1, 1000);
```

### 5.3 实时数据推送压力测试

```powershell
# 监控 WebSocket 消息频率
# 预期：消息推送间隔合理，不造成前端卡顿
```

---

## 六、全量测试完整性自查

- [ ] 实时监控页面所有交互已测试
- [ ] KPI 仪表板所有按钮已测试
- [ ] 运营分析页面所有交互已测试
- [ ] WebSocket 连接稳定性已验证
- [ ] 每个按钮提供 5 张截图 + 日志
- [ ] 每个按钮测试异常场景 ≥ 3 个
- [ ] 所有截图有明确的文件名
- [ ] 日志文件已打包

---

## 七、质量承诺

我承诺以上测试内容真实完整，所有按钮均已按 22 项清单验证。如有遗漏，愿意承担测试质量责任。

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
**创建日期**: 2024-12-18  
**最后更新**: 2024-12-18
