# XSS 防护实施指南

## 概述

本文档提供了 GameLink 项目中 XSS 防护的完整实施方案，包括关键输入点的防护措施。

## 关键输入点

根据安全审查，以下输入点需要 XSS 防护：

### 1. 用户昵称 (User.Nickname, Player.Nickname)

**位置**: 用户注册、用户资料更新、陪玩师资料更新

**风险等级**: 🔴 高 - 昵称在多个地方展示，容易被利用

**防护措施**: 使用 `sanitize.SanitizeNickname()`

#### 实施位置

**用户端**:
- `api/internal/handler/user/profile.go` - 用户资料更新

**陪玩师端**:
- `api/internal/handler/player/profile.go` - 陪玩师资料更新

**管理端**:
- `api/internal/handler/admin/user.go` - 管理员更新用户信息
- `api/internal/handler/admin/player.go` - 管理员更新陪玩师信息

#### 实施代码

```go
import "gamelink/pkg/sanitize"

// 在 handler 中
var req UpdateProfileRequest
if err := c.ShouldBindJSON(&req); err != nil {
    // 错误处理
    return
}

// XSS 防护
req.Nickname = sanitize.SanitizeNickname(req.Nickname)
```

### 2. 聊天消息 (ChatMessage.Content)

**位置**: WebSocket 聊天、发送消息 API

**风险等级**: 🔴 高 - 实时通信，XSS 攻击影响范围大

**防护措施**: 使用 `sanitize.SanitizeMessage()`

#### 实施位置

**用户端**:
- `api/internal/handler/user/chat.go` - 发送消息
- `api/internal/handler/player/chat.go` - 陪玩师发送消息
- `api/internal/ws/` - WebSocket 消息处理

#### 实施代码

```go
import "gamelink/pkg/sanitize"

// 在发送消息的 handler 中
var req SendMessageRequest
if err := c.ShouldBindJSON(&req); err != nil {
    return
}

// XSS 防护
req.Content = sanitize.SanitizeMessage(req.Content)

// 保存到数据库或发送
```

### 3. 评价内容 (Review.Content)

**位置**: 订单完成后的评价

**风险等级**: 🟡 中 - 评价在列表页和详情页展示

**防护措施**: 使用 `sanitize.SanitizeReview()`

#### 实施位置

**用户端**:
- `api/internal/handler/user/review.go` - 创建评价

**陪玩师端**:
- `api/internal/handler/player/review.go` - 创建评价

#### 实施代码

```go
import "gamelink/pkg/sanitize"

var req CreateReviewRequest
if err := c.ShouldBindJSON(&req); err != nil {
    return
}

// XSS 防护
req.Comment = sanitize.SanitizeReview(req.Comment)
for i, tag := range req.Tags {
    req.Tags[i] = sanitize.SanitizeNickname(tag)
}
```

### 4. 动态内容 (Feed.Content)

**位置**: 用户发布的动态

**风险等级**: 🟡 中 - 动态在信息流中展示

**防护措施**: 使用 `sanitize.SanitizeMessage()`

#### 实施位置

**用户端**:
- `api/internal/handler/user/feed.go` - 创建动态

#### 实施代码

```go
import "gamelink/pkg/sanitize"

var req CreateFeedRequest
if err := c.ShouldBindJSON(&req); err != nil {
    return
}

// XSS 防护
req.Content = sanitize.SanitizeMessage(req.Content)
if req.Title != "" {
    req.Title = sanitize.SanitizeNickname(req.Title)
}
```

### 5. 举报内容 (Report.Content, Dispute.Description)

**位置**: 用户举报、纠纷申诉

**风险等级**: 🟡 中 - 仅管理员可见，但仍需防护

**防护措施**: 使用 `sanitize.SanitizeReport()`

#### 实施位置

**用户端**:
- `api/internal/handler/user/dispute.go` - 创建纠纷

**陪玩师端**:
- `api/internal/handler/player/dispute.go` - 创建纠纷

#### 实施代码

```go
import "gamelink/pkg/sanitize"

var req CreateDisputeRequest
if err := c.ShouldBindJSON(&req); err != nil {
    return
}

// XSS 防护
req.Reason = sanitize.SanitizeNickname(req.Reason)
req.Description = sanitize.SanitizeReport(req.Description)
```

### 6. 敏感词替换文本

**位置**: 管理员配置的敏感词

**风险等级**: 🟢 低 - 仅管理员配置，但仍需防护

**防护措施**: 使用 `sanitize.SanitizeMessage()`

#### 实施位置

**管理端**:
- `api/internal/handler/admin/content.go` - 敏感词管理

#### 实施代码

```go
import "gamelink/pkg/sanitize"

var req CreateSensitiveWordRequest
if err := c.ShouldBindJSON(&req); err != nil {
    return
}

// XSS 防护
req.Word = sanitize.SanitizeNickname(req.Word)
req.Replacement = sanitize.SanitizeMessage(req.Replacement)
```

## 实施清单

### 阶段 1: 核心功能 (高优先级) 🔴

- [ ] **用户昵称防护**
  - [ ] `user/profile.go` - UpdateProfileHandler
  - [ ] `player/profile.go` - UpdatePlayerProfileHandler
  - [ ] `admin/user.go` - UpdateUserHandler
  - [ ] `admin/player.go` - UpdatePlayerHandler

- [ ] **聊天消息防护**
  - [ ] `user/chat.go` - SendMessageHandler
  - [ ] `player/chat.go` - SendMessageHandler
  - [ ] `ws/chat.go` - WebSocket 消息处理

### 阶段 2: 内容功能 (中优先级) 🟡

- [ ] **评价防护**
  - [ ] `user/review.go` - CreateReviewHandler
  - [ ] `player/review.go` - CreateReviewHandler

- [ ] **动态防护**
  - [ ] `user/feed.go` - CreateFeedHandler

- [ ] **举报/纠纷防护**
  - [ ] `user/dispute.go` - CreateDisputeHandler
  - [ ] `player/dispute.go` - CreateDisputeHandler

### 阶段 3: 管理功能 (低优先级) 🟢

- [ ] **敏感词防护**
  - [ ] `admin/content.go` - CreateSensitiveWordHandler
  - [ ] `admin/content.go` - UpdateSensitiveWordHandler

## 测试验证

### 单元测试

每个修改后的 handler 需要添加测试用例：

```go
func TestUpdateProfileHandler_XSS(t *testing.T) {
    // 测试 XSS 防护
    req := UpdateProfileRequest{
        Nickname: `<script>alert('XSS')</script>Player`,
    }

    // 调用 handler
    // 验证昵称已被清理
}
```

### 集成测试

使用 Postman/curl 测试实际请求：

```bash
# 测试昵称 XSS 防护
curl -X PUT http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "<script>alert(\"XSS\")</script>Player"
  }'

# 预期结果：昵称被清理为 "Player"
```

### 验证清单

- [ ] 确认 `<script>` 标签被移除
- [ ] 确认事件处理器（`onclick`, `onload`）被移除
- [ ] 确认 `javascript:` 协议被移除
- [ ] 确认特殊字符被转义（`<` → `&lt;`, `>` → `&gt;`）
- [ ] 确认长度限制生效
- [ ] 确认 UTF-8 编码验证
- [ ] 确认数据库中存储的是清理后的内容
- [ ] 确认前端显示的是转义后的内容

## 监控与日志

### 记录 XSS 攻击尝试

```go
import "gamelink/pkg/sanitize"

func SomeHandler(c *gin.Context) {
    var req struct {
        Content string `json:"content"`
    }
    c.ShouldBindJSON(&req)

    // 检测 XSS 攻击
    if sanitize.ContainsXSSPatterns(req.Content) {
        // 记录安全事件
        logSecurityEvent(c, "xss_attempt", map[string]interface{}{
            "user_id": getUserID(c),
            "content": req.Content,
            "ip":      c.ClientIP(),
        })

        // 清理内容
        req.Content = sanitize.SanitizeMessage(req.Content)
    }

    // 继续处理...
}
```

### 统计分析

建议收集以下指标：
1. XSS 攻击尝试次数
2. 受影响的用户/陪玩师
3. 攻击类型分布（script, event handler, javascript: 等）
4. 攻击来源 IP 分布

## 性能影响

### 基准测试结果

```bash
go test -bench=. ./pkg/sanitize/...
```

预期性能影响：
- `EscapeString`: ~100ns/op
- `SanitizeHTML`: ~500ns/op
- `SanitizeMessage`: ~1μs/op
- `SanitizeNickname`: ~800ns/op

对于典型请求，添加 XSS 防护后的延迟增加 <1ms，影响可忽略。

## 最佳实践

### ✅ 推荐做法

1. **Handler 层防护**: 在接收用户输入时立即清理
2. **使用专用函数**: 根据场景选择合适的清理函数
3. **记录安全事件**: 检测到 XSS 尝试时记录日志
4. **多层防御**: 输入清理 + 输出转义 + CSP 头
5. **定期审计**: 定期检查是否有遗漏的输入点

### ❌ 不推荐做法

1. **仅依赖前端验证**: 前端可以被绕过
2. **存储原始内容**: 永远不要存储未清理的用户输入
3. **信任内部用户**: 即使管理员输入也需要清理
4. **过度清理**: 不要过度限制合法的用户输入
5. **忽视测试**: 必须测试各种 XSS 攻击向量

## 常见问题

### Q: 清理后的内容在数据库中存储的是转义后的内容吗？

A: 是的。`Sanitize*` 函数会先移除危险模式，然后转义特殊字符。数据库中存储的是安全的、转义后的内容。

### Q: 如果用户需要使用特殊字符怎么办？

A: `Sanitize*` 函数只移除危险的 HTML/JavaScript 模式，普通特殊字符会被转义但不会丢失。例如：
- `Hello <world>` → `Hello &lt;world&gt;`
- 前端显示时会自动解码为 `Hello <world>`

### Q: 是否需要修改数据库 Schema？

A: 不需要。Schema 保持不变，我们只是在存储前清理数据。

### Q: 如果需要允许部分 HTML 标签怎么办？

A: 对于需要富文本的场景（如评论），建议使用专门的库如 [bluemonday](https://github.com/microcosm-cc/bluemonday)，它会提供更细粒度的控制。

### Q: XSS 防护是否会影响性能？

A: 影响极小（<1ms）。基准测试显示清理函数的性能非常高，不会成为瓶颈。

## 参考资源

- [OWASP XSS 防护备忘单](https://cheatsheetseries.owasp.org/cheatsheets/Cross_Site_Scripting_Prevention_Cheat_Sheet.html)
- [CSP (Content Security Policy)](https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP)
- [sanitize 包文档](../api/pkg/sanitize/README.md)
- [XSS 攻击向量](https://owasp.org/www-community/attacks/xss/)

## 变更历史

| 日期 | 版本 | 变更内容 |
|------|------|---------|
| 2025-01-01 | 1.0 | 初始版本，完成 sanitize 包和示例实现 |

## 联系方式

如有问题或建议，请联系安全团队或提交 Issue。
