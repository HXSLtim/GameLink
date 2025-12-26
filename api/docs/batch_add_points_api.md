# 批量增加积分API使用指南

## API端点
```
POST /api/v1/admin/users/batch/points
```

## 认证
需要管理员权限，Bearer Token认证

## 功能概述
支持三种批量增加积分的模式：
1. **指定用户列表模式** (target=users): 为指定的用户ID列表增加积分
2. **按角色筛选模式** (target=role): 为特定角色的所有用户增加积分
3. **全体用户模式** (target=all): 为系统中所有用户增加积分

## 请求参数说明

### 公共参数
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| target | string | 是 | 目标类型，可选值：users, role, all |
| cents | int64 | 是 | 增加的积分金额（分），范围：1-1000000（即0.01元-10000元） |
| reason | string | 是 | 增加原因，最大200字符 |
| type | string | 是 | 积分类型，可选值：admin, activity, compensation |

### 模式特定参数
| 参数 | 类型 | 使用场景 | 说明 |
|------|------|----------|------|
| userIds | []uint64 | target=users时必填 | 用户ID列表，最多1000个 |
| roles | []string | target=role时必填 | 角色列表，可选值：user, player, admin |

## 使用示例

### 1. 指定用户列表模式
为指定的用户增加积分（500分 = 5元）：

```bash
curl -X POST http://localhost:8080/api/v1/admin/users/batch/points \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target": "users",
    "userIds": [1, 2, 3, 100, 101],
    "cents": 500,
    "reason": "新年活动奖励",
    "type": "activity"
  }'
```

**响应示例：**
```json
{
  "success": true,
  "message": "批量增加积分成功，共增加2500分（25.00元）",
  "successCount": 5,
  "failedCount": 0
}
```

### 2. 按角色筛选模式
为所有player角色用户增加积分（1000分 = 10元）：

```bash
curl -X POST http://localhost:8080/api/v1/admin/users/batch/points \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target": "role",
    "roles": ["player"],
    "cents": 1000,
    "reason": "陪玩师月度奖励",
    "type": "admin"
  }'
```

**为多个角色增加积分（200分 = 2元）：**
```bash
curl -X POST http://localhost:8080/api/v1/admin/users/batch/points \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target": "role",
    "roles": ["user", "player"],
    "cents": 200,
    "reason": "系统升级补偿",
    "type": "compensation"
  }'
```

**响应示例：**
```json
{
  "success": true,
  "message": "批量增加积分成功，共增加150000分（1500.00元）",
  "successCount": 150,
  "failedCount": 0
}
```

### 3. 全体用户模式
为系统中所有用户增加积分（100分 = 1元）：

```bash
curl -X POST http://localhost:8080/api/v1/admin/users/batch/points \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target": "all",
    "cents": 100,
    "reason": "平台周年庆典",
    "type": "activity"
  }'
```

**响应示例：**
```json
{
  "success": true,
  "message": "批量增加积分成功，共增加50000分（500.00元）",
  "successCount": 500,
  "failedCount": 0
}
```

## 错误处理

### 常见错误

#### 1. 参数验证失败
```json
{
  "success": false,
  "message": "参数验证失败",
  "error": "Key: 'BatchAddPointsRequest.Target' Error:Field validation for 'Target' failed on the 'oneof' tag"
}
```

#### 2. 用户数量超限
```json
{
  "success": false,
  "message": "批量增加积分失败",
  "error": "符合条件的用户超过1000个，请缩小范围或分批操作"
}
```

#### 3. 缺少必要参数
当target=users时未提供userIds：
```json
{
  "success": false,
  "message": "批量增加积分失败",
  "error": "target为users时，userIds不能为空"
}
```

当target=role时未提供roles：
```json
{
  "success": false,
  "message": "批量增加积分失败",
  "error": "target为role时，roles不能为空"
}
```

#### 4. 部分失败
```json
{
  "success": true,
  "message": "批量增加积分完成，成功98个，失败2个",
  "successCount": 98,
  "failedCount": 2
}
```

## 业务逻辑说明

### 钱包自动创建
- 系统会自动为没有钱包的用户创建钱包记录（初始余额为0）
- 如果钱包创建失败，该用户会被计入失败数量，但不影响其他用户的处理

### 事务处理
- 所有积分增加操作在数据库事务中执行
- 单个用户的失败不会回滚整个批量操作
- 失败的用户会被记录到failedCount中

### 操作日志
- 每次批量操作都会记录详细的操作日志
- 日志包含：操作者ID、目标类型、影响用户数、积分数量、操作原因等
- 日志异步记录，不影响主操作性能

## 安全限制

1. **用户数量限制**: 单次操作最多1000个用户
2. **积分数量限制**: 单次最多增加1000000分（即10000元）
3. **角色验证**: 只能为有效的角色（user, player, admin）增加积分
4. **权限要求**: 需要管理员权限才能执行批量操作

## 最佳实践

1. **测试环境验证**: 在生产环境执行前，先在测试环境验证参数和结果
2. **分批处理大量用户**: 如果需要为超过1000个用户增加积分，建议分批执行
3. **记录操作原因**: 提供清晰的reason说明，便于后续审计
4. **监控失败数量**: 关注failedCount，如果失败率过高需要调查原因
5. **使用合适的type**: 根据实际场景选择admin（管理员发放）、activity（活动奖励）或compensation（补偿）

## 常见使用场景

### 场景1：新用户注册奖励
```json
{
  "target": "users",
  "userIds": [新注册的用户ID列表],
  "cents": 100,
  "reason": "新用户注册奖励",
  "type": "activity"
}
```
注：100分 = 1元

### 场景2：陪玩师月度绩效奖励
```json
{
  "target": "role",
  "roles": ["player"],
  "cents": 50000,
  "reason": "2025年1月绩效奖励",
  "type": "admin"
}
```
注：50000分 = 500元

### 场景3：系统维护补偿
```json
{
  "target": "all",
  "cents": 5000,
  "reason": "2025-01-15系统维护补偿",
  "type": "compensation"
}
```
注：5000分 = 50元

### 场景4：特定活动参与奖励
```json
{
  "target": "users",
  "userIds": [活动参与用户ID列表],
  "cents": 20000,
  "reason": "春节答题活动奖励",
  "type": "activity"
}
```
注：20000分 = 200元

## 技术实现说明

### 查询优化
- 按角色查询时使用索引优化：`WHERE role IN ?`
- 全体用户查询时只选择ID字段：`SELECT id FROM users`

### 并发处理
- 使用数据库事务保证数据一致性
- 每个用户的钱包更新使用FirstOrCreate模式避免并发问题

### 性能考虑
- 1000个用户的批量操作通常在2-5秒内完成
- 操作日志异步记录，不阻塞主流程
- 建议在低峰期执行大规模批量操作
