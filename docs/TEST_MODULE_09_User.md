# 测试任务单：用户管理模块全量测试

**任务编号**: TEST-2024-M09  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/sys/user | #createBtn | 新增用户 | POST /api/v1/admin/users | P0 | ☐ |
| /admin/sys/user | #editBtn | 编辑 | PUT /api/v1/admin/users/:id | P0 | ☐ |
| /admin/sys/user | #deleteBtn | 删除 | DELETE /api/v1/admin/users/:id | P0 | ☐ |
| /admin/sys/user | #banBtn | 封禁/解封 | PUT /api/v1/admin/users/:id/status | P0 | ☐ |
| /admin/sys/user | #detailBtn | 详情 | GET /api/v1/admin/users/:id | P0 | ☐ |
| /admin/sys/user | #batchDeleteBtn | 批量删除 | DELETE /api/v1/admin/users/batch | P1 | ☐ |
| /admin/sys/user | #batchRoleBtn | 批量修改角色 | PUT /api/v1/admin/users/batch/role | P1 | ☐ |
| /admin/sys/user | #batchStatusBtn | 批量修改状态 | PUT /api/v1/admin/users/batch/status | P1 | ☐ |
| /admin/sys/user | #batchNotifyBtn | 批量发送通知 | POST /api/v1/admin/notifications/batch | P1 | ☐ |
| /admin/sys/user | #batchPointsBtn | 批量增加积分 | POST /api/v1/admin/users/batch/points | P1 | ☐ |
| /admin/sys/user | #exportBtn | 导出数据 | 前端CSV导出 | P2 | ☐ |
| /admin/sys/user | #searchBtn | 搜索 | GET /api/v1/admin/users?keyword=xxx | P1 | ☐ |
| /admin/sys/user | #refreshBtn | 刷新 | GET /api/v1/admin/users | P1 | ☐ |

**重要**: 以上13个按钮，必须全部测试完成，少一个 = 任务未完成

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
# 在测试开始前执行
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

-- 查看用户统计
SELECT role, status, COUNT(*) FROM users GROUP BY role, status;

-- 查看用户钱包
SELECT u.id, u.name, w.balance_cents FROM users u LEFT JOIN wallets w ON u.id = w.user_id LIMIT 10;
```

### 测试账号
- **管理员**: admin@gameLink.com / Admin2025@Pass#

---

## 五、逐个按钮测试记录

### 按钮1: #createBtn 新增用户

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 点击"新增用户"按钮
4. 填写表单：
   - 用户名: 测试用户001
   - 邮箱: test001@test.com
   - 手机号: 13800138001
   - 密码: Test@123456
   - 角色: 普通用户
   - 状态: 正常
5. 点击保存

**Evidence收集**:
- [ ] 截图1: 新增用户弹窗表单
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/users
  - Payload: `{"name":"测试用户001","email":"test001@test.com",...}`
  - Status: 200
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM users WHERE email = 'test001@test.com';
  ```
- [ ] 截图5: 列表刷新显示新用户

**异常场景测试**:
- [ ] 场景A: 邮箱格式错误 → 预期: 前端校验提示
- [ ] 场景B: 重复邮箱 → 预期: 后端返回"邮箱已存在"
- [ ] 场景C: 密码强度不足 → 预期: 校验提示

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮2: #editBtn 编辑

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某用户行的"编辑"按钮
4. 修改用户名
5. 点击保存

**Evidence收集**:
- [ ] 截图1: 编辑弹窗（含原数据回显）
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/users/:id
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM users WHERE id = :user_id;
  ```
- [ ] 截图5: 列表显示更新后的数据

**异常场景测试**:
- [ ] 场景A: 修改为已存在的邮箱 → 预期: 错误提示
- [ ] 场景B: 清空必填字段 → 预期: 校验提示
- [ ] 场景C: 修改密码 → 预期: 密码正确更新

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮3: #deleteBtn 删除

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某用户行的"删除"按钮
4. 确认删除

**Evidence收集**:
- [ ] 截图1: 删除确认弹窗
- [ ] 截图2: Network请求详情
  - URL: DELETE /api/v1/admin/users/:id
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM users WHERE id = :deleted_id;
  -- 验证关联数据清理
  SELECT * FROM wallets WHERE user_id = :deleted_id;
  ```
- [ ] 截图5: 列表中用户已消失

**异常场景测试**:
- [ ] 场景A: 删除管理员账号 → 预期: 拒绝或警告
- [ ] 场景B: 删除有订单的用户 → 预期: 提示关联数据
- [ ] 场景C: 取消删除 → 预期: 无请求发送

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮4: #banBtn 封禁/解封

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某正常用户的"封禁"按钮
4. 确认操作

**Evidence收集**:
- [ ] 截图1: 封禁确认弹窗
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/users/:id/status
  - Payload: `{"status":"banned"}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, name, status FROM users WHERE id = :user_id;
  ```
- [ ] 截图5: 用户状态变为"已封禁"

**异常场景测试**:
- [ ] 场景A: 封禁管理员 → 预期: 拒绝操作
- [ ] 场景B: 解封已封禁用户 → 预期: 状态恢复正常
- [ ] 场景C: 封禁自己 → 预期: 拒绝操作

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮5: #detailBtn 详情

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某用户的"详情"按钮
4. 观察详情抽屉

**Evidence收集**:
- [ ] 截图1: 详情按钮
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/users/:id
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM users WHERE id = :user_id;
  ```
- [ ] 截图5: 详情抽屉显示完整信息

**异常场景测试**:
- [ ] 场景A: 查看登录历史Tab → 预期: 加载登录记录
- [ ] 场景B: 查看操作日志Tab → 预期: 加载操作记录
- [ ] 场景C: 关闭后重新打开 → 预期: 重新加载数据

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮6-10: 批量操作按钮

**批量删除 #batchDeleteBtn**:
- [ ] 选中多个用户后点击批量删除
- [ ] 验证: DELETE /api/v1/admin/users/batch
- [ ] 数据库验证所有选中用户已删除

**批量修改角色 #batchRoleBtn**:
- [ ] 选中多个用户后点击批量修改角色
- [ ] 选择目标角色
- [ ] 验证: PUT /api/v1/admin/users/batch/role
- [ ] 数据库验证角色已更新

**批量修改状态 #batchStatusBtn**:
- [ ] 选中多个用户后点击批量修改状态
- [ ] 选择目标状态
- [ ] 验证: PUT /api/v1/admin/users/batch/status
- [ ] 数据库验证状态已更新

**批量发送通知 #batchNotifyBtn**:
- [ ] 选中用户或选择全部
- [ ] 填写通知标题和内容
- [ ] 验证: POST /api/v1/admin/notifications/batch
- [ ] 数据库验证通知已创建

**批量增加积分 #batchPointsBtn**:
- [ ] 选中用户或选择全部
- [ ] 填写积分数量和原因
- [ ] 验证: POST /api/v1/admin/users/batch/points
- [ ] 数据库验证钱包余额已更新

---

### 按钮11-13: 搜索/导出/刷新

**搜索 #searchBtn**:
- [ ] 输入关键词搜索
- [ ] 验证: GET /api/v1/admin/users?keyword=xxx
- [ ] 验证筛选条件（角色、状态、日期范围）

**导出 #exportBtn**:
- [ ] 点击导出数据
- [ ] 验证CSV文件下载
- [ ] 验证数据完整性

**刷新 #refreshBtn**:
- [ ] 点击刷新按钮
- [ ] 验证: GET /api/v1/admin/users
- [ ] 验证数据与数据库一致

---

## 六、全量测试完整性自查

- [ ] 所有P0按钮已测试（5个）
- [ ] 所有P1按钮已测试（7个）
- [ ] 所有P2按钮已测试（1个）
- [ ] 每个按钮提供5张截图+日志
- [ ] 每个按钮测试异常场景≥3个

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
