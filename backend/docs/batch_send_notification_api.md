# 批量发送通知API使用指南

## API端点
```
POST /api/v1/admin/users/batch/notification
```

## 认证
需要管理员权限，Bearer Token认证

## 功能概述
支持三种批量发送通知的模式：
1. **指定用户列表模式** (target=users): 为指定的用户ID列表发送通知
2. **按角色筛选模式** (target=role): 为特定角色的所有用户发送通知
3. **全体用户模式** (target=all): 为系统中所有用户发送通知

## 请求参数说明

### 公共参数
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| target | string | 是 | 目标类型，可选值：users, role, all |
| title | string | 是 | 通知标题，最大100字符 |
| content | string | 是 | 通知内容，最大500字符 |
| type | string | 是 | 通知类型，可选值：system, marketing, personal, activity |

### 模式特定参数
| 参数 | 类型 | 使用场景 | 说明 |
|------|------|----------|------|
| userIds | []uint64 | target=users时必填 | 用户ID列表，最多1000个 |
| roles | []string | target=role时必填 | 角色列表，可选值：user, player, admin |

### 通知类型说明
- **system**: 系统通知（高优先级）
- **marketing**: 营销通知（普通优先级）
- **personal**: 个人通知（低优先级）
- **activity**: 活动通知（普通优先级）

## 使用示例

### 1. 指定用户列表模式
为指定的用户发送通知：

```bash
curl -X POST http://localhost:8080/api/v1/admin/users/batch/notification \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target": "users",
    "userIds": [1, 2, 3, 100, 101],
    "title": "系统升级通知",
    "content": "系统将于今晚23:00-24:00进行升级维护，届时将无法访问",
    "type": "system"
  }'
```

**响应示例：**
```json
{
  "success": true,
  "message": "批量发送通知成功"
}
```

### 2. 按角色筛选模式
为所有player角色用户发送通知：

```bash
curl -X POST http://localhost:8080/api/v1/admin/users/batch/notification \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target": "role",
    "roles": ["player"],
    "title": "陪玩师接单技巧培训",
    "content": "平台将于本周五举办陪玩师培训，欢迎参加",
    "type": "activity"
  }'
```

**为多个角色发送通知：**
```bash
curl -X POST http://localhost:8080/api/v1/admin/users/batch/notification \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target": "role",
    "roles": ["user", "player"],
    "title": "新年优惠活动",
    "content": "新年期间下单享8折优惠，快来参加吧！",
    "type": "marketing"
  }'
```

### 3. 全体用户模式
为系统中所有用户发送通知：

```bash
curl -X POST http://localhost:8080/api/v1/admin/users/batch/notification \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target": "all",
    "title": "平台周年庆典",
    "content": "感谢大家的支持，平台周年庆典活动即将开启，敬请期待！",
    "type": "activity"
  }'
```

## 错误处理

### 常见错误

#### 1. 参数验证失败
```json
{
  "success": false,
  "message": "参数验证失败",
  "error": "Key: 'BatchSendNotificationRequest.Target' Error:Field validation for 'Target' failed on the 'oneof' tag"
}
```

#### 2. 用户数量超限
```json
{
  "success": false,
  "message": "批量发送通知失败",
  "error": "符合条件的用户超过1000个，请缩小范围或分批操作"
}
```

#### 3. 缺少必要参数
当target=users时未提供userIds：
```json
{
  "success": false,
  "message": "批量发送通知失败",
  "error": "target为users时，userIds不能为空"
}
```

当target=role时未提供roles：
```json
{
  "success": false,
  "message": "批量发送通知失败",
  "error": "target为role时，roles不能为空"
}
```

## 业务逻辑说明

### 通知优先级
- system类型：高优先级（立即推送）
- marketing类型：普通优先级
- personal类型：低优先级
- activity类型：普通优先级

### 发送机制
- 通知通过web渠道发送
- 单个用户通知失败不影响其他用户
- 所有通知异步发送，不阻塞主流程

### 操作日志
- 每次批量操作都会记录详细的操作日志
- 日志包含：操作者ID、目标类型、影响用户数、通知标题、通知类型等
- 日志异步记录，不影响主操作性能

## 安全限制

1. **用户数量限制**: 单次操作最多1000个用户
2. **标题长度限制**: 最多100字符
3. **内容长度限制**: 最多500字符
4. **角色验证**: 只能为有效的角色（user, player, admin）发送通知
5. **权限要求**: 需要管理员权限才能执行批量操作

## 最佳实践

1. **测试环境验证**: 在生产环境执行前，先在测试环境验证参数和效果
2. **分批处理大量用户**: 如果需要为超过1000个用户发送通知，建议分批执行
3. **合理使用通知类型**: 根据实际场景选择合适的通知类型和优先级
4. **简洁明了的标题**: 标题应清晰表达通知主题
5. **避免过度推送**: 合理控制通知频率，避免打扰用户

## 常见使用场景

### 场景1：系统维护通知
```json
{
  "target": "all",
  "title": "系统维护通知",
  "content": "系统将于2025-01-20 23:00进行维护，预计持续1小时",
  "type": "system"
}
```

### 场景2：陪玩师培训通知
```json
{
  "target": "role",
  "roles": ["player"],
  "title": "陪玩师技能培训",
  "content": "本周六下午举办陪玩师技能培训，请准时参加",
  "type": "activity"
}
```

### 场景3：营销活动通知
```json
{
  "target": "role",
  "roles": ["user", "player"],
  "title": "双十一活动预告",
  "content": "双十一活动即将开启，下单立减，先到先得！",
  "type": "marketing"
}
```

### 场景4：重要用户通知
```json
{
  "target": "users",
  "userIds": [1001, 1002, 1003],
  "title": "VIP会员专属福利",
  "content": "您的VIP会员即将到期，续费享8折优惠",
  "type": "personal"
}
```

### 场景5：平台公告
```json
{
  "target": "all",
  "title": "平台规则更新",
  "content": "平台服务条款已更新，请查看最新版本",
  "type": "system"
}
```

## 技术实现说明

### 查询优化
- 按角色查询时使用索引优化：`WHERE role IN ?`
- 全体用户查询时只选择ID字段：`SELECT id FROM users`

### 并发处理
- 通知创建采用异步方式，不阻塞主流程
- 单个通知失败不影响整体操作

### 性能考虑
- 1000个用户的批量通知通常在2-5秒内完成
- 操作日志异步记录
- 建议在低峰期执行大规模通知推送
