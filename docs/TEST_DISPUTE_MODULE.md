# 测试任务单：纠纷管理模块全量测试

**任务编号**: TEST-2024-DISPUTE-001  
**测试环境**: Docker生产环境 (http://localhost)  
**测试日期**: 2024-12-18  
**模块路径**: `/admin/biz/dispute`

---

## 一、环境检查

### 1.1 容器状态检查
```bash
docker compose -f docker-compose.prod.yml ps
```

**预期结果**: 所有容器状态为 "Up (healthy)"
| 容器名称 | 状态 | 端口 |
|---------|------|------|
| gamelink-backend | Up (healthy) | 8081→8080 |
| gamelink-frontend | Up (healthy) | 80, 443 |
| gamelink-postgres | Up (healthy) | 5432 |
| gamelink-redis | Up (healthy) | 6379 |

### 1.2 测试账号
- **管理员**: `admin@gameLink.com` / `Admin2025@Pass#`

---

## 二、测试范围（必须100%覆盖）

| 序号 | 按钮名称 | 按钮位置 | 关联API | 优先级 | 测试状态 |
|-----|---------|---------|---------|--------|---------|
| 1 | 搜索 | 搜索栏 | GET /api/v1/admin/disputes | P0 | ☐ |
| 2 | 重置 | 搜索栏 | - | P1 | ☐ |
| 3 | 刷新 | 工具栏 | GET /api/v1/admin/disputes | P1 | ☐ |
| 4 | 导出数据 | 工具栏 | 前端导出CSV | P1 | ☐ |
| 5 | 详情 | 操作列 | GET /api/v1/admin/disputes/:id | P0 | ☐ |
| 6 | 分配 | 操作列(待处理) | POST /api/v1/admin/disputes/:id/assign | P0 | ☐ |
| 7 | 处理 | 操作列(已分配) | POST /api/v1/admin/disputes/:id/resolve | P0 | ☐ |
| 8 | 回滚 | 操作列(已分配) | POST /api/v1/admin/disputes/:id/rollback | P0 | ☐ |

---

## 三、数据库种子数据验证

### 3.1 验证纠纷数据存在
```bash
docker exec -it gamelink-postgres psql -U gamelink -d gamelink -c "SELECT id, order_id, status, reason FROM order_disputes;"
```

**预期结果**: 应有5条测试数据
| id | order_id | status | reason |
|----|----------|--------|--------|
| 1 | 1 | pending | 服务态度问题 |
| 2 | 2 | assigned | 技术水平不符 |
| 3 | 3 | mediating | 订单时间争议 |
| 4 | 4 | resolved | 退款金额争议 |
| 5 | 5 | resolved | 服务中断问题 |

---

## 四、逐个按钮测试记录

### 按钮1: 页面加载 + 统计数据

**测试目标**: 验证页面加载时统计卡片显示真实数据

**测试步骤**:
1. 登录管理后台
2. 导航到 业务管理 → 纠纷管理
3. 观察页面顶部4个统计卡片

**API验证**:
- [ ] 请求发送: `GET /api/v1/admin/disputes/stats` ✓/✗
- [ ] 响应状态: HTTP 200, code: 200 ✓/✗
- [ ] 响应数据结构:
```json
{
  "success": true,
  "code": 200,
  "data": {
    "pending": <number>,
    "assigned": <number>,
    "mediating": <number>,
    "resolved": <number>,
    "total": <number>
  }
}
```

**后端日志验证**:
```bash
docker logs gamelink-backend --tail=50 | grep -i "disputes/stats"
```

**数据库验证**:
```sql
SELECT status, COUNT(*) FROM order_disputes GROUP BY status;
```

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮2: 搜索（状态筛选）

**测试目标**: 验证按状态筛选纠纷列表

**测试步骤**:
1. 在状态下拉框选择"待处理"
2. 点击搜索按钮
3. 验证列表只显示待处理状态的纠纷

**API验证**:
- [ ] 请求发送: `GET /api/v1/admin/disputes?status=pending` ✓/✗
- [ ] 请求参数: `{ page: 1, pageSize: 10, status: "pending" }` ✓/✗
- [ ] 响应状态: HTTP 200 ✓/✗

**后端日志验证**:
```bash
docker logs gamelink-backend --tail=30 | grep "GET.*disputes"
```

**数据库验证**:
```sql
SELECT COUNT(*) FROM order_disputes WHERE status = 'pending';
```

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮3: 搜索（订单号筛选）

**测试目标**: 验证按订单号筛选纠纷

**测试步骤**:
1. 在订单号输入框输入已知订单号
2. 点击搜索按钮
3. 验证列表显示对应纠纷

**API验证**:
- [ ] 请求发送: `GET /api/v1/admin/disputes?orderNo=xxx` ✓/✗
- [ ] 响应数据包含匹配的纠纷 ✓/✗

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮4: 详情

**测试目标**: 验证查看纠纷详情抽屉

**测试步骤**:
1. 点击任意纠纷的"详情"按钮
2. 验证抽屉打开并显示完整信息

**API验证**:
- [ ] 请求发送: `GET /api/v1/admin/disputes/:id` ✓/✗
- [ ] 响应包含: 纠纷ID、订单号、用户、陪玩师、原因、状态等 ✓/✗

**页面验证**:
- [ ] 抽屉正常打开 ✓/✗
- [ ] 基本信息显示完整 ✓/✗
- [ ] 处理进度时间线显示 ✓/✗

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮5: 分配（核心功能）

**测试目标**: 验证将待处理纠纷分配给客服

**前置条件**: 存在状态为"pending"的纠纷

**测试步骤**:
1. 找到状态为"待处理"的纠纷
2. 点击"分配"按钮
3. 在弹窗中选择客服人员
4. 点击确定

**API验证**:
- [ ] 请求发送: `POST /api/v1/admin/disputes/:id/assign` ✓/✗
- [ ] 请求参数: `{ assignedToUserId: <number>, source: "manual" }` ✓/✗
- [ ] 响应状态: HTTP 200 ✓/✗

**后端日志验证**:
```bash
docker logs gamelink-backend --tail=30 | grep -i "assign"
```

**数据库验证**:
```sql
SELECT id, status, assigned_to_user_id, assigned_at 
FROM order_disputes 
WHERE id = <测试ID>;
```

**预期结果**:
- status 变为 "assigned"
- assigned_to_user_id 有值
- assigned_at 有时间戳

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮6: 处理（核心功能）

**测试目标**: 验证处理已分配的纠纷

**前置条件**: 存在状态为"assigned"或"mediating"的纠纷

**测试步骤**:
1. 找到状态为"已分配"的纠纷
2. 点击"处理"按钮
3. 选择解决方案（全额退款/部分退款/重新分配/驳回）
4. 填写退款金额（如适用）
5. 填写处理备注
6. 点击确定

**API验证**:
- [ ] 请求发送: `POST /api/v1/admin/disputes/:id/resolve` ✓/✗
- [ ] 请求参数:
```json
{
  "resolution": "refund|partial|reassign|reject",
  "resolutionAmount": <number>,
  "resolutionNotes": "<string>"
}
```
- [ ] 响应状态: HTTP 200 ✓/✗

**后端日志验证**:
```bash
docker logs gamelink-backend --tail=30 | grep -i "resolve"
```

**数据库验证**:
```sql
SELECT id, status, resolution, resolution_amount, resolution_notes, resolved_at 
FROM order_disputes 
WHERE id = <测试ID>;
```

**预期结果**:
- status 变为 "resolved"
- resolution 有值
- resolved_at 有时间戳

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮7: 回滚

**测试目标**: 验证回滚已分配的纠纷

**前置条件**: 存在状态为"assigned"的纠纷

**测试步骤**:
1. 找到状态为"已分配"的纠纷
2. 点击"回滚"按钮
3. 确认回滚操作

**API验证**:
- [ ] 请求发送: `POST /api/v1/admin/disputes/:id/rollback` ✓/✗
- [ ] 请求参数: `{ rollbackReason: "<string>" }` ✓/✗
- [ ] 响应状态: HTTP 200 ✓/✗

**数据库验证**:
```sql
SELECT id, status, assigned_to_user_id, rolled_back_at 
FROM order_disputes 
WHERE id = <测试ID>;
```

**预期结果**:
- status 变回 "pending"
- assigned_to_user_id 变为 NULL
- rolled_back_at 有时间戳

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮8: 导出数据

**测试目标**: 验证导出纠纷列表为CSV

**测试步骤**:
1. 点击"导出数据"按钮
2. 验证CSV文件下载

**验证项**:
- [ ] 文件成功下载 ✓/✗
- [ ] 文件名包含 "disputes" ✓/✗
- [ ] CSV内容包含所有列 ✓/✗
- [ ] 数据与页面显示一致 ✓/✗

**测试结果**: ☐ 通过 ☐ 失败

---

## 五、异常场景测试

| 场景 | 操作 | 预期结果 | 实际结果 | 状态 |
|------|------|---------|---------|------|
| 无权限访问 | 用普通用户访问纠纷页面 | 403或重定向 | | ☐ |
| 分配不存在的纠纷 | 手动构造请求 | 404错误 | | ☐ |
| 重复分配 | 对已分配纠纷再次分配 | 400错误提示 | | ☐ |
| 处理已解决纠纷 | 对resolved状态纠纷点处理 | 400错误提示 | | ☐ |
| 快速连续点击 | 连续点击分配按钮5次 | 只执行一次 | | ☐ |

---

## 六、容器日志检查

### 6.1 后端错误日志
```bash
docker logs gamelink-backend --since="1h" | grep -i "error\|panic\|fatal"
```
**预期**: 无错误日志

### 6.2 数据库连接日志
```bash
docker logs gamelink-backend --since="1h" | grep -i "database\|postgres\|connection"
```
**预期**: 无连接失败

---

## 七、测试总结

### 7.1 测试覆盖率
- 总按钮数: 8
- 已测试: ___
- 覆盖率: ___%

### 7.2 发现问题

| 问题编号 | 问题描述 | 严重程度 | 状态 |
|---------|---------|---------|------|
| | | | |

### 7.3 测试结论
- [ ] 全部通过
- [ ] 部分通过，存在问题
- [ ] 测试失败

---

**测试人签字**: ___________  
**日期**: ___________

**审核人签字**: ___________  
**日期**: ___________
