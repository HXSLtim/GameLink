# 测试任务单：权限管理模块全量测试

**任务编号**: TEST-2024-M21  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/sys/permission | #viewPermissionsBtn | 查看权限树 | GET /api/v1/admin/permissions | P0 | ☐ |
| /admin/sys/permission | #createPermissionBtn | 新增权限 | POST /api/v1/admin/permissions | P0 | ☐ |
| /admin/sys/permission | #editPermissionBtn | 编辑 | PUT /api/v1/admin/permissions/:id | P0 | ☐ |
| /admin/sys/permission | #deletePermissionBtn | 删除 | DELETE /api/v1/admin/permissions/:id | P0 | ☐ |
| /admin/sys/permission | #searchBtn | 搜索 | GET /api/v1/admin/permissions?keyword=xxx | P1 | ☐ |
| /admin/sys/permission | #refreshBtn | 刷新 | GET /api/v1/admin/permissions | P1 | ☐ |

**重要**: 以上6个按钮，必须全部测试完成

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

-- 查看权限列表
SELECT id, name, slug, parent_id, type, created_at FROM permissions ORDER BY id;

-- 查看权限树结构
WITH RECURSIVE perm_tree AS (
  SELECT id, name, slug, parent_id, 0 as level
  FROM permissions WHERE parent_id IS NULL
  UNION ALL
  SELECT p.id, p.name, p.slug, p.parent_id, pt.level + 1
  FROM permissions p
  JOIN perm_tree pt ON p.parent_id = pt.id
)
SELECT * FROM perm_tree ORDER BY level, id;

-- 查看角色权限关联
SELECT r.name as role_name, COUNT(rp.permission_id) as permission_count
FROM roles r
LEFT JOIN role_permissions rp ON r.id = rp.role_id
GROUP BY r.id, r.name;
```

### 权限类型说明
- `menu`: 菜单权限
- `button`: 按钮权限
- `api`: API权限

---

## 四、逐个按钮测试记录

### 按钮1: #viewPermissionsBtn 查看权限树

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 进入权限管理页面
4. 观察权限树加载

**Evidence收集**:
- [ ] 截图1: 权限树展示
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/permissions
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT COUNT(*) FROM permissions;
  ```
- [ ] 截图5: 权限树结构与数据库一致

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮2: #createPermissionBtn 新增权限

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击"新增权限"按钮
4. 填写表单：
   - 权限名称: 测试权限
   - 权限标识: test.permission
   - 父级权限: 选择或留空
   - 权限类型: 按钮
5. 点击保存

**Evidence收集**:
- [ ] 截图1: 新增权限弹窗
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/permissions
  - Payload: `{"name":"测试权限","slug":"test.permission","type":"button"}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM permissions WHERE slug = 'test.permission';
  ```
- [ ] 截图5: 权限树显示新权限

**异常场景测试**:
- [ ] 场景A: 名称为空 → 预期: 校验提示
- [ ] 场景B: 标识重复 → 预期: 后端返回错误
- [ ] 场景C: 标识格式错误 → 预期: 校验提示

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮3: #editPermissionBtn 编辑权限

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某权限的"编辑"按钮
4. 修改权限名称
5. 点击保存

**Evidence收集**:
- [ ] 截图1: 编辑弹窗（含原数据回显）
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/permissions/:id
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
- [ ] 截图5: 权限树显示更新后的数据

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮4: #deletePermissionBtn 删除权限

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某权限的"删除"按钮
4. 确认删除

**Evidence收集**:
- [ ] 截图1: 删除确认弹窗
- [ ] 截图2: Network请求详情
  - URL: DELETE /api/v1/admin/permissions/:id
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM permissions WHERE id = :deleted_id;
  -- 验证子权限处理
  SELECT * FROM permissions WHERE parent_id = :deleted_id;
  -- 验证角色权限关联清理
  SELECT * FROM role_permissions WHERE permission_id = :deleted_id;
  ```
- [ ] 截图5: 权限树中权限已消失

**异常场景测试**:
- [ ] 场景A: 删除有子权限的权限 → 预期: 提示或级联删除
- [ ] 场景B: 删除已分配给角色的权限 → 预期: 提示或清理关联
- [ ] 场景C: 取消删除 → 预期: 无请求发送

**测试结果**: ☐ 通过 ☐ 失败

---

### 按钮5-6: 搜索/刷新

**搜索 #searchBtn**:
- [ ] 输入权限名称或标识搜索
- [ ] 验证: GET /api/v1/admin/permissions?keyword=xxx
- [ ] 验证搜索结果正确

**刷新 #refreshBtn**:
- [ ] 点击刷新按钮
- [ ] 验证: GET /api/v1/admin/permissions
- [ ] 验证数据与数据库一致

---

## 五、全量测试完整性自查

- [ ] 所有P0按钮已测试（4个）
- [ ] 所有P1按钮已测试（2个）
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
