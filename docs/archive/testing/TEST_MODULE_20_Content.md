# 测试任务单：内容管理模块全量测试

**任务编号**: TEST-2024-M20  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、模块说明

内容管理模块包含多个子模块：
- **动态管理 (Feeds)**: 用户发布的动态内容
- **聊天监控 (ChatMonitor)**: 聊天记录监控
- **举报管理 (Reports)**: 用户举报处理
- **分类管理 (Categories)**: 内容分类配置
- **内容统计 (Stats)**: 内容数据统计

---

## 二、测试范围（必须100%覆盖）

### 动态管理 (Feeds)
| 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|--------|---------|---------|--------|---------|
| #viewFeedBtn | 查看 | GET /api/v1/admin/feeds/:id | P0 | ☐ |
| #deleteFeedBtn | 删除 | DELETE /api/v1/admin/feeds/:id | P0 | ☐ |
| #hideFeedBtn | 隐藏 | PUT /api/v1/admin/feeds/:id/hide | P1 | ☐ |
| #searchFeedBtn | 搜索 | GET /api/v1/admin/feeds?keyword=xxx | P1 | ☐ |

### 聊天监控 (ChatMonitor)
| 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|--------|---------|---------|--------|---------|
| #viewChatBtn | 查看聊天 | GET /api/v1/admin/chats/:id | P0 | ☐ |
| #warnUserBtn | 警告用户 | POST /api/v1/admin/users/:id/warn | P1 | ☐ |
| #filterChatBtn | 筛选 | GET /api/v1/admin/chats?filter=xxx | P1 | ☐ |

### 举报管理 (Reports)
| 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|--------|---------|---------|--------|---------|
| #viewReportBtn | 查看 | GET /api/v1/admin/reports/:id | P0 | ☐ |
| #handleReportBtn | 处理 | PUT /api/v1/admin/reports/:id/handle | P0 | ☐ |
| #dismissReportBtn | 驳回 | PUT /api/v1/admin/reports/:id/dismiss | P0 | ☐ |

### 分类管理 (Categories)
| 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|--------|---------|---------|--------|---------|
| #createCategoryBtn | 新增分类 | POST /api/v1/admin/categories | P0 | ☐ |
| #editCategoryBtn | 编辑 | PUT /api/v1/admin/categories/:id | P0 | ☐ |
| #deleteCategoryBtn | 删除 | DELETE /api/v1/admin/categories/:id | P0 | ☐ |

---

## 三、Docker环境检查

```bash
docker compose -f docker-compose.prod.yml ps
```

---

## 四、测试数据准备

### 数据库种子数据验证
```sql
-- 连接数据库
docker exec -it gamelink-postgres psql -U gamelink -d gamelink

-- 查看动态
SELECT id, user_id, content, status, created_at FROM feeds LIMIT 10;

-- 查看举报
SELECT id, reporter_id, target_type, target_id, reason, status FROM reports LIMIT 10;

-- 查看分类
SELECT id, name, type, sort_order, is_active FROM categories;
```

---

## 五、关键按钮测试记录

### 举报处理 #handleReportBtn

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 找到待处理的举报
4. 点击"处理"按钮
5. 选择处理方式（警告/封禁/删除内容）
6. 确认处理

**Evidence收集**:
- [ ] 截图1: 处理弹窗
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/reports/:id/handle
  - Payload: `{"action":"warn","reason":"违规内容"}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, status, handled_at, handled_by FROM reports WHERE id = :report_id;
  ```
- [ ] 截图5: 举报状态变为"已处理"

**异常场景测试**:
- [ ] 场景A: 处理已处理的举报 → 预期: 按钮不显示
- [ ] 场景B: 未选择处理方式 → 预期: 校验提示
- [ ] 场景C: 网络中断时处理 → 预期: 错误提示

**测试结果**: ☐ 通过 ☐ 失败

---

### 分类管理 #createCategoryBtn

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击"新增分类"按钮
4. 填写分类名称、类型、排序
5. 点击保存

**Evidence收集**:
- [ ] 截图1: 新增分类弹窗
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/categories
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM categories WHERE name = '测试分类';
  ```
- [ ] 截图5: 分类列表显示新分类

**测试结果**: ☐ 通过 ☐ 失败

---

## 六、全量测试完整性自查

- [ ] 动态管理所有按钮已测试（4个）
- [ ] 聊天监控所有按钮已测试（3个）
- [ ] 举报管理所有按钮已测试（3个）
- [ ] 分类管理所有按钮已测试（3个）
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
**审核人**: ___________  
**日期**: ___________

---

**文档版本**: v1.0  
**发布日期**: 2024-12-18
