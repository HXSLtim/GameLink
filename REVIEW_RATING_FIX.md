# 评价管理评分显示修复

## 问题描述

评价管理页面的评分字段显示为 `undefined`，无法正常显示评价的星级评分。

## 原因分析

后端返回的字段名与前端期望的字段名不匹配：

**后端返回的数据：**
```json
{
  "id": 1,
  "order_id": 1,
  "user_id": 2,
  "player_id": 1,
  "score": 5,
  "content": "很满意的陪玩体验..."
}
```

**前端期望的字段：**
```typescript
interface Review {
  order_id: number;
  reviewer_id: number;  // ❌ 后端返回的是 user_id
  player_id: number;
  rating: number;       // ❌ 后端返回的是 score
  comment?: string;     // ❌ 后端返回的是 content
}
```

**字段不匹配导致的问题：**
- `record.rating` 访问不到数据 → `undefined`
- `record.reviewer_id` 访问不到数据 → `undefined`
- `record.comment` 访问不到数据 → `undefined`

## 解决方案

修改后端 `Review` 模型的 JSON 标签，使其与前端期望的字段名完全匹配。

**修改文件：** `backend/internal/model/review.go`

### 修改前

```go
type Review struct {
	Base
	OrderID  uint64 `json:"order_id" gorm:"index"`
	UserID   uint64 `json:"user_id" gorm:"index"`
	PlayerID uint64 `json:"player_id" gorm:"index"`
	Score    Rating `json:"score" gorm:"type:tinyint"`
	Content  string `json:"content,omitempty" gorm:"type:text"`
}
```

### 修改后

```go
type Review struct {
	Base
	OrderID  uint64 `json:"order_id" gorm:"index"`
	UserID   uint64 `json:"reviewer_id" gorm:"column:user_id;index"` // 前端使用reviewer_id
	PlayerID uint64 `json:"player_id" gorm:"index"`
	Score    Rating `json:"rating" gorm:"column:score;type:tinyint"`   // 前端使用rating
	Content  string `json:"comment,omitempty" gorm:"column:content;type:text"` // 前端使用comment
}
```

### 关键变更

1. **UserID 字段**
   - JSON 标签：`user_id` → `reviewer_id`
   - 添加 `column:user_id` 指定数据库列名
   - 原因：前端使用 `reviewer_id` 更语义化

2. **Score 字段**
   - JSON 标签：`score` → `rating`
   - 添加 `column:score` 指定数据库列名
   - 原因：前端使用 `rating` 符合评分语义

3. **Content 字段**
   - JSON 标签：`content` → `comment`
   - 添加 `column:content` 指定数据库列名
   - 原因：前端使用 `comment` 更符合评论语义

## 修改后的 API 响应

```json
{
  "success": true,
  "code": 200,
  "data": [
    {
      "id": 1,
      "order_id": 1,
      "reviewer_id": 2,    // ✓ 改为 reviewer_id
      "player_id": 1,
      "rating": 5,         // ✓ 改为 rating
      "comment": "很满意的陪玩体验..." // ✓ 改为 comment
    }
  ]
}
```

## 前端显示效果

修复后，评价管理页面应该能正确显示：

1. **评分列** - 显示星级评分（如：⭐⭐⭐⭐⭐ 非常满意）
2. **评价人列** - 显示评价人ID和姓名
3. **评价内容列** - 显示评价文字内容

## 测试验证

### 测试步骤

1. 重启后端服务：
```bash
cd backend
make run CMD=user-service
```

2. 访问评价管理页面：`http://localhost:5173/reviews`

### 预期结果

- ✅ 评分列正确显示星级（1-5星）
- ✅ 评分颜色正确：
  - 5星：绿色（非常满意）
  - 4星：绿色（满意）
  - 3星：蓝色（一般）
  - 2星：橙色（较差）
  - 1星：红色（非常差）
- ✅ 评价人列显示正确的用户ID
- ✅ 评价内容列显示完整的评论文字

## 技术说明

### GORM column 标签的作用

当 JSON 字段名与数据库列名不同时，需要使用 `column` 标签指定数据库列名：

```go
UserID uint64 `json:"reviewer_id" gorm:"column:user_id"`
```

这样：
- JSON 序列化时使用 `reviewer_id`
- 数据库查询时使用 `user_id` 列

### 为什么不修改前端？

前端代码已经在多处使用了 `reviewer_id`、`rating`、`comment` 这些字段名：
- ReviewList.tsx（列表页面）
- ReviewFormModal.tsx（表单模态框）
- ReviewDetail.tsx（详情页面）
- 可能还有其他组件

修改前端需要改动更多文件，而修改后端只需改一个模型文件，风险更小。

### 命名一致性问题

注意到前端 Review 类型混用了两种命名风格：
- snake_case: `order_id`, `reviewer_id`, `player_id`
- camelCase: `rating`, `comment`

**建议未来优化：**
统一改为 camelCase（符合 JavaScript/TypeScript 规范）：
```typescript
interface Review {
  orderId: number;
  reviewerId: number;
  playerId: number;
  rating: number;
  comment?: string;
}
```

但这需要：
1. 修改后端 JSON 标签为 camelCase
2. 修改所有使用 Review 的前端组件
3. 确保不影响其他功能

## 相关文件

### 后端
- `backend/internal/model/review.go` - Review 模型定义（已修改）

### 前端
- `frontend/src/types/review.ts` - Review 类型定义
- `frontend/src/pages/Reviews/ReviewList.tsx` - 评价列表页面
- `frontend/src/pages/Reviews/ReviewFormModal.tsx` - 评价表单
- `frontend/src/services/api/review.ts` - 评价 API 服务

## 总结

通过修改后端 Review 模型的 JSON 标签，使其与前端期望的字段名匹配，成功修复了评分显示问题。

修改要点：
- `user_id` → `reviewer_id`（更语义化）
- `score` → `rating`（符合评分语义）
- `content` → `comment`（符合评论语义）

所有修改都通过 GORM 的 `column` 标签保持了数据库列名不变，确保数据库兼容性。

