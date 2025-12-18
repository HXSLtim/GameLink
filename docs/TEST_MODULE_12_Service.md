# 测试任务单：服务项管理模块全量测试

**任务编号**: TEST-2024-M12  
**测试环境**: Docker生产环境 (localhost:80)  
**负责人**: [测试人员姓名]  
**截止日期**: 2024-XX-XX 18:00  
**监督人**: [组长姓名]

---

## 一、测试范围（必须100%覆盖）

| 页面路径 | 按钮ID | 按钮文案 | 关联API | 优先级 | 测试状态 |
|---------|--------|---------|---------|--------|---------|
| /admin/biz/service | #createBtn | 新增服务 | POST /api/v1/admin/services | P0 | ☐ |
| /admin/biz/service | #editBtn | 编辑 | PUT /api/v1/admin/services/:id | P0 | ☐ |
| /admin/biz/service | #deleteBtn | 删除 | DELETE /api/v1/admin/services/:id | P0 | ☐ |
| /admin/biz/service | #statusSwitch | 启用/禁用 | PUT /api/v1/admin/services/:id/status | P0 | ☐ |
| /admin/biz/service | #searchBtn | 搜索 | GET /api/v1/admin/services?keyword=xxx | P1 | ☐ |
| /admin/biz/service | #gameFilter | 游戏筛选 | 前端过滤 | P1 | ☐ |

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

-- 查看服务项统计
SELECT g.name as game_name, COUNT(s.id) as service_count
FROM services s
JOIN games g ON s.game_id = g.id
GROUP BY g.name;

-- 查看服务项详情
SELECT s.id, s.name, s.category, s.base_price, s.max_price, s.is_active, g.name as game_name
FROM services s
LEFT JOIN games g ON s.game_id = g.id
LIMIT 10;
```

### 服务分类说明
- `陪玩`: 陪玩服务
- `代练`: 代练服务
- `陪聊`: 陪聊服务
- `其他`: 其他服务

---

## 五、逐个按钮测试记录

### 按钮1: #createBtn 新增服务

**测试步骤**:
1. 清空日志: `docker logs gamelink-backend --tail=0`
2. 打开开发者工具 → Network面板
3. 点击"新增服务"按钮
4. 填写表单：
   - 服务名称: 测试上分陪玩
   - 所属游戏: 王者荣耀
   - 服务分类: 陪玩
   - 计费单位: 小时
   - 最低价格: 30
   - 最高价格: 100
   - 服务描述: 测试服务描述
5. 点击保存

**Evidence收集**:
- [ ] 截图1: 新增服务弹窗表单
- [ ] 截图2: Network请求详情
  - URL: POST /api/v1/admin/services
  - Payload: `{"name":"测试上分陪玩","gameId":1,"category":"陪玩",...}`
  - Status: 200
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM services WHERE name = '测试上分陪玩';
  ```
- [ ] 截图5: 列表刷新显示新服务

**异常场景测试**:
- [ ] 场景A: 名称为空提交 → 预期: 前端校验提示
- [ ] 场景B: 最低价格大于最高价格 → 预期: 校验提示
- [ ] 场景C: 未选择游戏 → 预期: 校验提示

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮2: #editBtn 编辑

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某服务行的"编辑"按钮
4. 修改服务名称和价格
5. 点击保存

**Evidence收集**:
- [ ] 截图1: 编辑弹窗（含原数据回显）
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/services/:id
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM services WHERE id = :service_id;
  ```
- [ ] 截图5: 列表显示更新后的数据

**异常场景测试**:
- [ ] 场景A: 清空必填字段 → 预期: 校验提示
- [ ] 场景B: 编辑不存在的服务 → 预期: 404错误
- [ ] 场景C: 并发编辑 → 预期: 后提交者覆盖

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮3: #deleteBtn 删除

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某服务行的"删除"按钮
4. 确认删除

**Evidence收集**:
- [ ] 截图1: 删除确认弹窗
- [ ] 截图2: Network请求详情
  - URL: DELETE /api/v1/admin/services/:id
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT * FROM services WHERE id = :deleted_id;
  ```
- [ ] 截图5: 列表中服务已消失

**异常场景测试**:
- [ ] 场景A: 删除有订单关联的服务 → 预期: 提示或拒绝
- [ ] 场景B: 删除不存在的服务 → 预期: 404错误
- [ ] 场景C: 取消删除 → 预期: 无请求发送

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮4: #statusSwitch 启用/禁用

**测试步骤**:
1. 清空日志
2. 打开Network面板
3. 点击某服务的状态开关
4. 观察状态变化

**Evidence收集**:
- [ ] 截图1: 状态开关
- [ ] 截图2: Network请求详情
  - URL: PUT /api/v1/admin/services/:id/status
  - Payload: `{"isActive":false}`
- [ ] 截图3: docker logs处理记录
- [ ] 截图4: 数据库验证
  ```sql
  SELECT id, name, is_active FROM services WHERE id = :service_id;
  ```
- [ ] 截图5: 开关状态已切换

**异常场景测试**:
- [ ] 场景A: 禁用后再启用 → 预期: 状态正确切换
- [ ] 场景B: 快速连续切换 → 预期: 防抖生效
- [ ] 场景C: 网络中断时切换 → 预期: 错误提示，状态回滚

**测试结果**: ☐ 通过 ☐ 失败  
**失败原因**: （如有）

---

### 按钮5-6: 搜索/筛选

**搜索 #searchBtn**:
- [ ] 输入服务名称搜索
- [ ] 验证列表过滤结果

**游戏筛选 #gameFilter**:
- [ ] 选择特定游戏
- [ ] 验证只显示该游戏的服务

---

## 六、全量测试完整性自查

- [ ] 所有P0按钮已测试（4个）
- [ ] 所有P1按钮已测试（2个）
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
