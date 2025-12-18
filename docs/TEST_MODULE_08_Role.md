# 测试任务单：角色管理模块全量测试

**任务编号**: TEST-2024-M08  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/sys/role | #createBtn | 新增角色 | POST /api/v1/admin/roles | P0 | ☐ |
| /admin/sys/role | #editBtn | 编辑 | PUT /api/v1/admin/roles/:id | P0 | ☐ |
| /admin/sys/role | #deleteBtn | 删除 | DELETE /api/v1/admin/roles/:id | P0 | ☐ |
| /admin/sys/role | #permissionBtn | 权限 | GET /api/v1/admin/roles/:id/permissions | P0 | ☐ |
| /admin/sys/role | #searchBtn | 搜索 | GET /api/v1/admin/roles?keyword=xxx | P1 | ☐ |
| /admin/sys/role | #refreshBtn | 刷新 | GET /api/v1/admin/roles | P1 | ☐ |

**重要**: 以上6个按钮，必须全部测试完成，少一个 = 任务未完成

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
- gamelink-backend: Up (healthy)
- gamelink-frontend: Up (healthy)
- gamelink-postgres: Up (healthy)
- gamelink-redis: Up (healthy)

**将结果截图贴在此处**:

---

## 四、测试数据准备

### 数据库种子数据验证
```sql
-- 连接数据库
docker exec -it gamelink-postgres psql -U gamelink -d gamelink

-- 查看现有角色
SELECT id, name, slug, description, is_system, created_at FROM roles ORDER BY id;

-- 查看角色权限关联
SELECT r.name as role_name, COUNT(rp.permission_id) as permission_count 
FROM roles r 
LEFT JOIN role_permissions rp ON r.id = rp.role_id 
GROUP BY r.id, r.name;

-- 查看角色用户数
SELECT r.name as role_name, COUNT(ur.user_id) as user_count 
FROM roles r 
LEFT JOIN user_roles ur ON r.id = ur.role_id 
GROUP BY r.id, r.name;
```

### 测试账号
- **管理员**: admin@gameLink.com / Admin2025@Pass#

### 系统角色说明
- `superAdmin`: 超级管理员（系统角色，不可删除）
- `admin`: 管理员（系统角色，不可删除）
- `operator`: 运营人员
- `customer_service`: 客服人员

---

## 五、逐个按钮测试记录

### 按钮1: #createBtn 新增角色

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 点击"新增角色"按钮
4. 填写表单：
   - 角色名称: 测试运营角色
   - 角色编码: test_operator
   - 描述: 测试用角色
5. 点击确定
6. 监控容器日志

**Evidence收集**:
- [ ] 截图1: 新增角色弹窗表单
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/roles
  - Payload: `{"name":"测试运营角色","slug":"test_operator","description":"测试用角色"}`
  - Status: 200
- [ ] 截图3: docker logs gamelink-backend 处理记录
  ```bash
  docker logs gamelink-backend --tail=20 | findstr "roles"
  ```
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM roles WHERE slug = 'test_operator';
  ```
- [ ] 截图5: 列表刷新显示新角色

**异常场景测试**:
- [ ] 场景A: 名称为空提交 → 预期: 前端校验提示"请输入角色名称"
- [ ] 场景B: 编码格式错误(含大写) → 预期: 前端校验提示"只能输入小写字母和下划线"
- [ ] 场景C: 重复编码提交 → 预期: 后端返回错误"角色编码已存在"

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮2: #editBtn 编辑

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某非系统角色行的"编辑"按钮
4. 修改角色名称为"测试运营角色-已修改"
5. 点击确定

**Evidence收集**:
- [ ] 截图1: 编辑弹窗（含原数据回显）
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/roles/:id
  - Payload: `{"name":"测试运营角色-已修改","slug":"test_operator","description":"..."}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM roles WHERE slug = 'test_operator';
  -- 验证 updated_at 已更新
  ```
- [ ] 截图5: 列表显示更新后的角色名

**异常场景测试**:
- [ ] 场景A: 系统角色编辑编码 → 预期: 编码字段禁用
- [ ] 场景B: 修改为已存在的编码 → 预期: 错误提示
- [ ] 场景C: 编辑不存在的角色 → 预期: 404错误

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮3: #deleteBtn 删除

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某非系统角色行的"删除"按钮
4. 确认删除弹窗点击"确定"

**Evidence收集**:
- [ ] 截图1: 删除确认弹窗
- [ ] 截图2: Network请求详情
  - URL: DELETE /api/v1/admin/roles/:id
  - Status: 200
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  -- 验证角色已删除
  SELECT * FROM roles WHERE slug = 'test_operator';
  -- 验证权限关联已清理
  SELECT * FROM role_permissions WHERE role_id = :deleted_id;
  -- 验证用户角色关联已清理
  SELECT * FROM user_roles WHERE role_id = :deleted_id;
  ```
- [ ] 截图5: 列表中角色已消失

**异常场景测试**:
- [ ] 场景A: 删除系统角色 → 预期: 提示"系统角色不可删除"
- [ ] 场景B: 删除有用户关联的角色 → 预期: 提示确认或拒绝
- [ ] 场景C: 取消删除 → 预期: 无任何请求发送

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮4: #permissionBtn 权限配置

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某角色行的"权限"按钮
4. 观察跳转到权限配置页面

**Evidence收集**:
- [ ] 截图1: 权限按钮可点击
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/roles/:id/permissions
  - Response: 权限树数据
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT p.id, p.name, p.slug 
  FROM permissions p 
  JOIN role_permissions rp ON p.id = rp.permission_id 
  WHERE rp.role_id = :role_id;
  ```
- [ ] 截图5: 权限配置页面正确显示已有权限

**异常场景测试**:
- [ ] 场景A: 超级管理员角色 → 预期: 显示"拥有所有权限"提示
- [ ] 场景B: 无权限角色 → 预期: 权限树全部未勾选
- [ ] 场景C: 保存权限变更 → 预期: 权限正确更新

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮5: #searchBtn 搜索

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 在搜索框输入关键词"admin"
4. 点击搜索或按回车

**Evidence收集**:
- [ ] 截图1: 搜索框输入状态
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/roles?keyword=admin
  - Status: 200
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM roles WHERE name ILIKE '%admin%' OR slug ILIKE '%admin%';
  ```
- [ ] 截图5: 列表只显示匹配的角色

**异常场景测试**:
- [ ] 场景A: 搜索不存在的关键词 → 预期: 显示空列表
- [ ] 场景B: 清空搜索框 → 预期: 显示全部角色
- [ ] 场景C: 特殊字符搜索 → 预期: 正常过滤，无报错

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮6: #refreshBtn 刷新

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击刷新按钮
4. 观察列表重新加载

**Evidence收集**:
- [ ] 截图1: 刷新按钮
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/roles
  - Status: 200
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库当前数据
  ```sql
  SELECT COUNT(*) FROM roles;
  ```
- [ ] 截图5: 列表数据与数据库一致

**异常场景测试**:
- [ ] 场景A: 后端服务不可用时刷新 → 预期: 错误提示"加载失败"
- [ ] 场景B: 快速连续点击刷新 → 预期: 防抖生效

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

## 六、全量测试完整性自查

- [ ] 所有P0按钮已测试（4个）
- [ ] 所有P1按钮已测试（2个）
- [ ] 每个按钮提供5张截图+日志
- [ ] 每个按钮测试异常场景≥3个
- [ ] 所有截图有明确的文件名（btnName_stepNumber.png）
- [ ] 日志文件已打包（logs.tar.gz）

---

## 七、质量承诺

我承诺以上测试内容真实完整，所有按钮均已按22项清单验证。如有遗漏，愿意承担测试质量责任。

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
**监督人**: [组长姓名]  
**批准人**: [技术总监姓名]
