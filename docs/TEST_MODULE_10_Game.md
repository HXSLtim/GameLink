# 测试任务单：游戏管理模块全量测试

**任务编号**: TEST-2024-M10  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/biz/game | #createBtn | 新增游戏 | POST /api/v1/admin/games | P0 | ☐ |
| /admin/biz/game | #editBtn | 编辑 | PUT /api/v1/admin/games/:id | P0 | ☐ |
| /admin/biz/game | #deleteBtn | 删除 | DELETE /api/v1/admin/games/:id | P0 | ☐ |
| /admin/biz/game | #batchDeleteBtn | 批量删除 | DELETE /api/v1/admin/games/batch | P1 | ☐ |
| /admin/biz/game | #exportBtn | 导出数据 | 前端CSV导出 | P2 | ☐ |
| /admin/biz/game | #searchBtn | 搜索 | GET /api/v1/admin/games?keyword=xxx | P1 | ☐ |
| /admin/biz/game | #refreshBtn | 刷新 | GET /api/v1/admin/games | P1 | ☐ |

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
docker compose -f docker-compose.prod.yml ps
```

**预期结果**: 所有容器状态为"Up (healthy)"

---

## 四、测试数据准备

### 数据库种子数据验证
```sql
-- 连接数据库
docker exec -it gamelink-postgres psql -U gamelink -d gamelink

-- 查看现有游戏
SELECT id, key, name, category, created_at FROM games ORDER BY id;

-- 查看游戏分类统计
SELECT category, COUNT(*) FROM games GROUP BY category;
```

### 游戏分类说明
- `moba`: MOBA游戏
- `fps`: 射击游戏
- `rpg`: RPG游戏
- `card`: 卡牌游戏
- `casual`: 休闲游戏
- `other`: 其他

---

## 五、逐个按钮测试记录

### 按钮1: #createBtn 新增游戏

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 点击"新增游戏"按钮
4. 填写表单：
   - 游戏Key: test_game_001
   - 游戏名称: 测试游戏001
   - 游戏图标URL: https://example.com/icon.png
   - 分类: MOBA
   - 描述: 测试用游戏
5. 点击确定

**Evidence收集**:
- [ ] 截图1: 新增游戏弹窗表单
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/games
  - Payload: `{"key":"test_game_001","name":"测试游戏001","category":"moba",...}`
  - Status: 200
- [ ] 截图3: docker logs处理记录
  ```bash
  docker logs gamelink-backend --tail=20 | findstr "games"
  ```
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM games WHERE key = 'test_game_001';
  ```
- [ ] 截图5: 列表刷新显示新游戏

**异常场景测试**:
- [ ] 场景A: Key为空提交 → 预期: 前端校验提示"请输入游戏Key"
- [ ] 场景B: 重复Key提交 → 预期: 后端返回错误"游戏Key已存在"
- [ ] 场景C: 名称为空 → 预期: 前端校验提示"请输入游戏名称"

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮2: #editBtn 编辑

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某游戏行的"编辑"按钮
4. 修改游戏名称为"测试游戏001-已修改"
5. 点击确定

**Evidence收集**:
- [ ] 截图1: 编辑弹窗（含原数据回显，Key字段禁用）
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/games/:id
  - Payload: `{"name":"测试游戏001-已修改",...}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM games WHERE key = 'test_game_001';
  -- 验证 updated_at 已更新
  ```
- [ ] 截图5: 列表显示更新后的游戏名

**异常场景测试**:
- [ ] 场景A: 修改为已存在的名称 → 预期: 允许（名称可重复）
- [ ] 场景B: 清空名称 → 预期: 校验提示
- [ ] 场景C: 编辑不存在的游戏 → 预期: 404错误

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮3: #deleteBtn 删除

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某游戏行的"删除"按钮
4. 确认删除弹窗点击"确定"

**Evidence收集**:
- [ ] 截图1: 删除确认弹窗
- [ ] 截图2: Network请求详情
  - URL: DELETE /api/v1/admin/games/:id
  - Status: 200
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  -- 验证游戏已删除
  SELECT * FROM games WHERE key = 'test_game_001';
  -- 验证关联服务项
  SELECT * FROM services WHERE game_id = :deleted_id;
  ```
- [ ] 截图5: 列表中游戏已消失

**异常场景测试**:
- [ ] 场景A: 删除有服务项关联的游戏 → 预期: 提示或级联删除
- [ ] 场景B: 删除不存在的游戏 → 预期: 404错误
- [ ] 场景C: 取消删除 → 预期: 无任何请求发送

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮4: #batchDeleteBtn 批量删除

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击"批量删除"按钮
4. 选择删除目标（选中的游戏/按分类/全部）
5. 确认删除

**Evidence收集**:
- [ ] 截图1: 批量删除弹窗
- [ ] 截图2: Network请求详情
  - URL: DELETE /api/v1/admin/games/batch
  - Payload: `{"ids":[1,2,3]}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT COUNT(*) FROM games WHERE id IN (1,2,3);
  ```
- [ ] 截图5: 列表中游戏已消失

**异常场景测试**:
- [ ] 场景A: 未选中任何游戏 → 预期: 提示"请选择要删除的游戏"
- [ ] 场景B: 按分类删除 → 预期: 该分类所有游戏被删除
- [ ] 场景C: 删除全部 → 预期: 警告确认后删除

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮5: #exportBtn 导出数据

**测试步骤**:
1. 点击"导出数据"按钮
2. 观察文件下载

**Evidence收集**:
- [ ] 截图1: 导出按钮可点击
- [ ] 截图2: 无Network请求（前端导出）
- [ ] 截图5: 下载的CSV文件内容正确

**特别验证**:
- [ ] CSV文件名格式: games_YYYYMMDD.csv
- [ ] CSV包含列: ID, Key, 名称, 分类, 描述, 创建时间, 更新时间
- [ ] 数据与页面列表一致

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮6: #searchBtn 搜索

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 在搜索框输入关键词"王者"
4. 选择分类筛选
5. 点击搜索

**Evidence收集**:
- [ ] 截图1: 搜索框和筛选条件
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/games?keyword=王者&category=moba
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM games WHERE name ILIKE '%王者%' AND category = 'moba';
  ```
- [ ] 截图5: 列表只显示匹配的游戏

**异常场景测试**:
- [ ] 场景A: 搜索不存在的关键词 → 预期: 显示空列表
- [ ] 场景B: 清空搜索条件 → 预期: 显示全部游戏
- [ ] 场景C: 组合筛选 → 预期: 正确过滤

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮7: #refreshBtn 刷新

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击刷新按钮

**Evidence收集**:
- [ ] 截图1: 刷新按钮
- [ ] 截图2: Network请求详情
  - URL: GET /api/v1/admin/games
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库当前数据
  ```sql
  SELECT COUNT(*) FROM games;
  ```
- [ ] 截图5: 列表数据与数据库一致

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

## 六、全量测试完整性自查

- [ ] 所有P0按钮已测试（3个）
- [ ] 所有P1按钮已测试（3个）
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
