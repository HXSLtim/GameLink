# 测试任务单：用户标签管理模块全量测试

**任务编号**: TEST-2024-M01  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/sys/user-tag | #createBtn | 新增标签 | POST /api/v1/admin/user-tags | P0 | ☐ |
| /admin/sys/user-tag | #editBtn | 编辑 | PUT /api/v1/admin/user-tags/:id | P0 | ☐ |
| /admin/sys/user-tag | #deleteBtn | 删除 | DELETE /api/v1/admin/user-tags/:id | P0 | ☐ |
| /admin/sys/user-tag | #viewUsersBtn | 查看用户(X人) | GET /api/v1/admin/user-tags/:id/users | P1 | ☐ |
| /admin/sys/user-tag | #searchBtn | 搜索 | 前端过滤(无API) | P1 | ☐ |
| /admin/sys/user-tag | #exportBtn | 导出数据 | 前端CSV导出 | P1 | ☐ |
| /admin/sys/user-tag | #refreshBtn | 刷新 | GET /api/v1/admin/user-tags | P1 | ☐ |

**重要**: 以上7个按钮，必须全部测试完成，少一个 = 任务未完成

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

-- 查看现有标签
SELECT id, name, color, description, created_at FROM user_tags ORDER BY id;

-- 查看标签用户关联数
SELECT ut.id, ut.name, COUNT(utr.user_id) as user_count 
FROM user_tags ut 
LEFT JOIN user_tag_relations utr ON ut.id = utr.tag_id 
GROUP BY ut.id, ut.name;
```

### 测试账号
- **管理员**: 使用 `.env` 中的 `SUPER_ADMIN_EMAIL` / `SUPER_ADMIN_PASSWORD`

---

## 五、逐个按钮测试记录

### 按钮1: #createBtn 新增标签

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 点击"新增标签"按钮
4. 填写表单：
   - 标签名称: 测试VIP用户
   - 标签颜色: #ff4d4f (红色)
   - 描述: 测试用标签
5. 点击确定
6. 监控容器日志

**Evidence收集**:
- [ ] 截图1: 新增标签弹窗表单
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/user-tags
  - Payload: `{"name":"测试VIP用户","color":"#ff4d4f","description":"测试用标签"}`
  - Status: 200
- [ ] 截图3: docker logs gamelink-backend 处理记录
  ```bash
  docker logs gamelink-backend --tail=20 | findstr "user-tags"
  ```
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM user_tags WHERE name = '测试VIP用户';
  ```
- [ ] 截图5: 列表刷新显示新标签

**异常场景测试**:
- [ ] 场景A: 名称为空提交 → 预期: 前端校验提示"请输入标签名称"
- [ ] 场景B: 重复名称提交 → 预期: 后端返回错误"标签名称已存在"
- [ ] 场景C: 快速连续点击确定5次 → 预期: 防抖生效，只创建1条

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮2: #editBtn 编辑

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某标签行的"编辑"按钮
4. 修改标签名称为"测试VIP用户-已修改"
5. 点击确定

**Evidence收集**:
- [ ] 截图1: 编辑弹窗（含原数据回显）
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/user-tags/1
  - Payload: `{"name":"测试VIP用户-已修改","color":"#ff4d4f","description":"..."}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM user_tags WHERE id = 1;
  -- 验证 updated_at 已更新
  ```
- [ ] 截图5: 列表显示更新后的标签名

**异常场景测试**:
- [ ] 场景A: 修改为已存在的名称 → 预期: 错误提示
- [ ] 场景B: 编辑不存在的标签(手动改URL) → 预期: 404错误
- [ ] 场景C: 并发编辑同一标签 → 预期: 后提交者覆盖

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮3: #deleteBtn 删除

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某标签行的"删除"按钮
4. 确认删除弹窗点击"确定"

**Evidence收集**:
- [ ] 截图1: 删除确认弹窗
- [ ] 截图2: Network请求详情
  - URL: DELETE /api/v1/admin/user-tags/1
  - Status: 200
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  -- 验证标签已删除
  SELECT * FROM user_tags WHERE id = 1;
  -- 验证关联关系已清理
  SELECT * FROM user_tag_relations WHERE tag_id = 1;
  ```
- [ ] 截图5: 列表中标签已消失

**异常场景测试**:
- [ ] 场景A: 删除有用户关联的标签 → 预期: 提示"删除后所有用户的该标签将被移除"
- [ ] 场景B: 删除不存在的标签 → 预期: 404错误
- [ ] 场景C: 取消删除 → 预期: 无任何请求发送

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮4: #viewUsersBtn 查看用户(X人)

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某标签的用户数链接（如"5人"）
4. 观察弹窗加载用户列表

**Evidence收集**:
- [ ] 截图1: 用户数链接可点击
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/user-tags/1/users?page=1&page_size=10
  - Response: 用户列表数据
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT u.id, u.name, u.email 
  FROM users u 
  JOIN user_tag_relations utr ON u.id = utr.user_id 
  WHERE utr.tag_id = 1;
  ```
- [ ] 截图5: 弹窗显示用户列表

**异常场景测试**:
- [ ] 场景A: 标签无关联用户 → 预期: 显示"暂无数据"
- [ ] 场景B: 翻页测试 → 预期: 分页正常
- [ ] 场景C: 关闭弹窗后重新打开 → 预期: 重新加载数据

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮5: #searchBtn 搜索

**测试步骤**:
1. 在搜索框输入关键词"VIP"
2. 点击搜索或按回车
3. 观察列表过滤结果

**Evidence收集**:
- [ ] 截图1: 搜索框输入状态
- [ ] 截图2: 无Network请求（前端过滤）
- [ ] 截图3: 不适用
- [ ] 截图4: 不适用
- [ ] 截图5: 列表只显示包含"VIP"的标签

**异常场景测试**:
- [ ] 场景A: 搜索不存在的关键词 → 预期: 显示空列表
- [ ] 场景B: 清空搜索框 → 预期: 显示全部标签
- [ ] 场景C: 特殊字符搜索 → 预期: 正常过滤，无报错

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮6: #exportBtn 导出数据

**测试步骤**:
1. 点击"导出数据"按钮
2. 观察文件下载

**Evidence收集**:
- [ ] 截图1: 导出按钮可点击
- [ ] 截图2: 无Network请求（前端导出）
- [ ] 截图3: 不适用
- [ ] 截图4: 不适用
- [ ] 截图5: 下载的CSV文件内容正确

**特别验证**:
- [ ] CSV文件名格式: user_tags_YYYYMMDD.csv
- [ ] CSV包含列: ID, 标签名称, 颜色, 描述, 用户数, 创建时间
- [ ] 数据与页面列表一致

**异常场景测试**:
- [ ] 场景A: 空列表导出 → 预期: 只有表头的CSV
- [ ] 场景B: 大量数据导出(100+) → 预期: 正常导出，无卡顿

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮7: #refreshBtn 刷新

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击刷新按钮
4. 观察列表重新加载

**Evidence收集**:
- [ ] 截图1: 刷新按钮
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/user-tags
  - Status: 200
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库当前数据
  ```sql
  SELECT COUNT(*) FROM user_tags;
  ```
- [ ] 截图5: 列表数据与数据库一致

**异常场景测试**:
- [ ] 场景A: 后端服务不可用时刷新 → 预期: 错误提示"加载失败"
- [ ] 场景B: 快速连续点击刷新 → 预期: 防抖生效

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

## 六、全量测试完整性自查

- [ ] 所有P0按钮已测试（3个）
- [ ] 所有P1按钮已测试（4个）
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
