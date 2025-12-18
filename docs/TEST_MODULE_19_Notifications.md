# 测试任务单：通知管理模块全量测试

**任务编号**: TEST-2024-M19  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/sys/notifications | #markReadBtn | 标记已读 | PUT /api/v1/notifications/:id/read | P0 | ☐ |
| /admin/sys/notifications | #markAllReadBtn | 全部已读 | PUT /api/v1/notifications/read-all | P0 | ☐ |
| /admin/sys/notifications | #deleteBtn | 删除 | DELETE /api/v1/notifications/:id | P0 | ☐ |
| /admin/sys/notifications | #loadMoreBtn | 加载更多 | GET /api/v1/notifications?page=2 | P1 | ☐ |

**重要**: 以上4个按钮，必须全部测试完成

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

-- 查看通知列表
SELECT id, title, message, is_read, created_at 
FROM notifications 
WHERE user_id = :admin_user_id 
ORDER BY created_at DESC LIMIT 20;

-- 查看未读通知数
SELECT COUNT(*) as unread_count 
FROM notifications 
WHERE user_id = :admin_user_id AND is_read = false;
```

---

## 四、逐个按钮测试记录

### 按钮1: #markReadBtn 标记已读

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 找到一条未读通知（左侧有蓝色边框）
4. 点击"✓"按钮标记已读

**Evidence收集**:
- [ ] 截图1: 未读通知卡片
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/notifications/:id/read
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, is_read, read_at FROM notifications WHERE id = :notification_id;
  ```
- [ ] 截图5: 通知卡片样式变为已读状态

**异常场景测试**:
- [ ] 场景A: 标记已读的通知 → 预期: 按钮不显示
- [ ] 场景B: 网络中断时标记 → 预期: 错误提示
- [ ] 场景C: 标记不存在的通知 → 预期: 404错误

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮2: #markAllReadBtn 全部已读

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击"全部已读"按钮

**Evidence收集**:
- [ ] 截图1: 全部已读按钮
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/notifications/read-all
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT COUNT(*) FROM notifications 
  WHERE user_id = :admin_user_id AND is_read = false;
  -- 应该返回 0
  ```
- [ ] 截图5: 所有通知变为已读状态

**异常场景测试**:
- [ ] 场景A: 无未读通知时点击 → 预期: 正常完成
- [ ] 场景B: 大量未读通知 → 预期: 正常处理
- [ ] 场景C: 网络中断时点击 → 预期: 错误提示

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮3: #deleteBtn 删除

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某通知的删除按钮
4. 确认删除

**Evidence收集**:
- [ ] 截图1: 删除确认弹窗
- [ ] 截图2: Network请求详情
  - URL: DELETE /api/v1/notifications/:id
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM notifications WHERE id = :deleted_id;
  ```
- [ ] 截图5: 通知从列表中消失

**异常场景测试**:
- [ ] 场景A: 取消删除 → 预期: 无请求发送
- [ ] 场景B: 删除不存在的通知 → 预期: 404错误
- [ ] 场景C: 网络中断时删除 → 预期: 错误提示

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮4: #loadMoreBtn 加载更多

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 滚动到底部
4. 点击"加载更多"按钮

**Evidence收集**:
- [ ] 截图1: 加载更多按钮
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/notifications?page=2&page_size=20
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
- [ ] 截图5: 列表追加新数据

**异常场景测试**:
- [ ] 场景A: 无更多数据 → 预期: 按钮消失
- [ ] 场景B: 网络中断时加载 → 预期: 错误提示
- [ ] 场景C: 快速连续点击 → 预期: 防抖生效

**测试结果**: ☐ 通过 ☐ 失败

---

## 五、全量测试完整性自查

- [ ] 所有P0按钮已测试（3个）
- [ ] 所有P1按钮已测试（1个）
- [ ] 每个按钮提供5张截图+日志
- [ ] 每个按钮测试异常场景≥3个

---

## 六、质量承诺

我承诺以上测试内容真实完整，所有按钮均已按22项清单验证。

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
