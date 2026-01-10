# Sanitize Package - XSS 防护工具包

## 概述

`sanitize` 包提供了全面的输入清理功能，用于防止跨站脚本攻击（XSS）和其他注入攻击。

## 主要功能

### 1. 基础转义函数

#### `EscapeString(s string) string`
转义 HTML 特殊字符（`&`, `<`, `>`, `"`, `'`）为 HTML 实体。

```go
import "gamelink/pkg/sanitize"

input := `<script>alert('XSS')</script>`
output := sanitize.EscapeString(input)
// 输出: &lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;
```

#### `UnescapeString(s string) string`
反转义 HTML 实体（仅用于可信数据）。

### 2. 高级清理函数

#### `SanitizeHTML(s string) string`
移除危险的 HTML/JavaScript 模式：
- `<script>` 标签
- 事件处理器（`onclick`, `onload` 等）
- `javascript:` 协议
- `data:text/html` URI

```go
input := `<div onclick="alert('XSS')">Content</div>`
output := sanitize.SanitizeHTML(input)
// 输出: <div>Content</div>
```

#### `ContainsXSSPatterns(s string) bool`
检测字符串是否包含潜在的 XSS 模式。

```go
if sanitize.ContainsXSSPatterns(userInput) {
    // 拒绝或额外处理
}
```

### 3. 特定场景清理函数

#### `SanitizeNickname(nickname string) string`
专门用于清理用户昵称：
- 移除所有 HTML 标签
- 限制长度为 32 个字符
- 转义特殊字符

```go
nickname := `<b>Player</b><script>alert(1)</script>`
cleanNickname := sanitize.SanitizeNickname(nickname)
// 输出: Player
```

#### `SanitizeMessage(content string) string`
用于清理聊天消息：
- 移除危险模式
- 验证 UTF-8 编码
- 限制长度为 5000 个字符
- 转义特殊字符

```go
message := `<img src=x onerror=alert(1)>Hello`
cleanMessage := sanitize.SanitizeMessage(message)
// 输出: <img src=x>Hello
```

#### `SanitizeReview(content string) string`
用于清理评价内容（规则与消息相同）。

#### `SanitizeReport(content string) string`
用于清理举报/申诉内容：
- 限制长度为 1000 个字符
- 移除危险模式
- 验证 UTF-8 编码

### 4. 工具函数

#### `StripTags(s string) string`
移除所有 HTML 标签，仅保留文本内容。

```go
input := `<p>Hello <strong>World</strong></p>`
output := sanitize.StripTags(input)
// 输出: Hello World
```

#### `TruncateString(s string, maxLength int) string`
截断字符串到指定长度，添加省略号。

```go
input := "This is a very long string"
output := sanitize.TruncateString(input, 10)
// 输出: This is a ...
```

#### `ValidateUTF8(s string) bool`
验证字符串是否为有效的 UTF-8 编码。

#### `EscapeAll(data map[string]interface{}) map[string]interface{}`
转义 map 中的所有字符串字段（包括嵌套结构）。

```go
input := map[string]interface{}{
    "name": "John<script>",
    "profile": map[string]interface{}{
        "bio": "<b>Hello</b>",
    },
}
output := sanitize.EscapeAll(input)
// output["name"] = "John&lt;script&gt;"
// output["profile"].(map[string]interface{})["bio"] = "&lt;b&gt;Hello&lt;/b&gt;"
```

## 使用场景

### 在 Handler 层使用

```go
import (
    "github.com/gin-gonic/gin"
    "gamelink/pkg/sanitize"
)

func UpdateNickname(c *gin.Context) {
    var req struct {
        Nickname string `json:"nickname"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        // 错误处理
        return
    }

    // 清理用户输入
    cleanNickname := sanitize.SanitizeNickname(req.Nickname)

    // 使用清理后的数据
    user.Nickname = cleanNickname
    // ... 保存到数据库
}
```

### 批量清理请求体

```go
func CreatePost(c *gin.Context) {
    var req map[string]interface{}

    if err := c.ShouldBindJSON(&req); err != nil {
        // 错误处理
        return
    }

    // 清理所有字符串字段
    cleanReq := sanitize.EscapeAll(req)

    // 使用清理后的数据
    // ...
}
```

### XSS 检测

```go
func CreateUserContent(c *gin.Context) {
    content := c.PostForm("content")

    // 检测是否包含 XSS 模式
    if sanitize.ContainsXSSPatterns(content) {
        c.JSON(400, gin.H{"error": "内容包含非法字符"})
        return
    }

    // 或者直接清理
    cleanContent := sanitize.SanitizeMessage(content)
}
```

## 安全原则

1. **输入验证与清理**：始终在接收用户输入时进行清理
2. **输出编码**：在将数据输出到 HTML 时确保已转义
3. **深度防御**：结合使用多种防护手段
4. **白名单优先**：使用 `StripTags` 等函数移除 HTML，而不是尝试过滤危险标签

## 性能考虑

所有函数都经过优化，可以安全地在热路径中使用。对于高性能场景，参考基准测试：

```bash
go test -bench=. ./pkg/sanitize/...
```

## 测试

包包含全面的单元测试和基准测试：

```bash
# 运行所有测试
go test ./pkg/sanitize/...

# 运行测试并查看覆盖率
go test -cover ./pkg/sanitize/...

# 运行基准测试
go test -bench=. ./pkg/sanitize/...
```

## 最佳实践

### ✅ 推荐

```go
// 1. 在 Handler 层立即清理
nickname := sanitize.SanitizeNickname(req.Nickname)

// 2. 使用特定场景的函数
message := sanitize.SanitizeMessage(req.Message)
review := sanitize.SanitizeReview(req.Review)

// 3. 检测 XSS 攻击
if sanitize.ContainsXSSPatterns(input) {
    // 记录日志或拒绝
}
```

### ❌ 不推荐

```go
// 1. 不要在数据库中存储未清理的用户输入
user.Nickname = req.Nickname // ❌ 危险

// 2. 不要信任前端已清理的数据
// 前端可以被绕过，必须在后端再次清理

// 3. 不要使用正则表达式自行实现清理
// 使用已验证的 sanitize 包函数
```

## 限制

1. **不适用于富文本编辑器**：如果需要允许部分 HTML（如 `<b>`, `<i>`），应使用专门的 HTML 清理库（如 [bluemonday](https://github.com/microcosm-cc/bluemonday)）

2. **不是完整的 CSP 解决方案**：仍需配置 Content-Security-Policy 头

3. **不防御存储型 XSS 的所有场景**：需结合其他安全措施

## 扩展

如需添加新的清理函数，请遵循以下模式：

```go
// Sanitize[Feature] sanitizes [description]
func Sanitize[Feature](input string) string {
    if input == "" {
        return input
    }

    // 1. 移除危险模式
    input = SanitizeHTML(input)

    // 2. 验证编码
    if !ValidateUTF8(input) {
        // 清理无效字符
    }

    // 3. 限制长度
    input = TruncateString(input, maxLength)

    // 4. 转义特殊字符
    return EscapeString(input)
}
```

## 贡献

在提交 PR 前，请确保：
1. 添加相应的单元测试
2. 更新基准测试
3. 更新此文档

## 许可证

MIT License - 详见项目根目录的 LICENSE 文件
