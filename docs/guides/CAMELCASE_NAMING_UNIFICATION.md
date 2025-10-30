# camelCase 命名统一规范

## 概述

本次修改将评价（Review）模块的字段命名统一为 camelCase（小驼峰命名），符合 JavaScript/TypeScript 的最佳实践。

## 命名规范

### 原则

1. **后端优先**：后端先定义 JSON 标签为 camelCase
2. **前端跟随**：前端 TypeScript 类型定义使用与后端完全一致的 camelCase
3. **数据库不变**：通过 GORM 的 `column` 标签保持数据库列名为 snake_case

### 规范对比

| 类型 | 命名方式 | 示例 | 说明 |
|------|---------|------|------|
| Go 结构体字段 | PascalCase | `OrderID` | Go 语言规范 |
| JSON 字段（API） | camelCase | `orderId` | JavaScript/TypeScript 规范 |
| 数据库列名 | snake_case | `order_id` | 数据库规范 |
| TypeScript 字段 | camelCase | `orderId` | 与 JSON 保持一致 |

## 修改详情

### 1. 后端 Review 模型

**文件：** `backend/internal/model/review.go`

#### 修改前（混合命名）

```go
type Review struct {
    Base
    OrderID  uint64 `json:"order_id" gorm:"index"`
    UserID   uint64 `json:"reviewer_id" gorm:"column:user_id;index"`
    PlayerID uint64 `json:"player_id" gorm:"index"`
    Score    Rating `json:"rating" gorm:"column:score;type:tinyint"`
    Content  string `json:"comment,omitempty" gorm:"column:content;type:text"`
}
```

#### 修改后（完全 camelCase）

```go
type Review struct {
    Base
    OrderID  uint64 `json:"orderId" gorm:"column:order_id;index"`
    UserID   uint64 `json:"reviewerId" gorm:"column:user_id;index"`
    PlayerID uint64 `json:"playerId" gorm:"column:player_id;index"`
    Score    Rating `json:"rating" gorm:"column:score;type:tinyint"`
    Content  string `json:"comment,omitempty" gorm:"column:content;type:text"`
}
```

**字段映射：**

| Go 字段 | JSON 标签（新） | GORM 列名 | 说明 |
|---------|----------------|-----------|------|
| OrderID | `orderId` | `order_id` | 订单ID |
| UserID | `reviewerId` | `user_id` | 评价人ID |
| PlayerID | `playerId` | `player_id` | 陪玩师ID |
| Score | `rating` | `score` | 评分 (1-5) |
| Content | `comment` | `content` | 评价内容 |

### 2. 前端 Review 类型定义

**文件：** `frontend/src/types/review.ts`

#### 修改前（混合命名）

```typescript
export interface Review extends BaseEntity {
  order_id: number;      // snake_case
  reviewer_id: number;   // snake_case
  player_id: number;     // snake_case
  rating: number;        // camelCase
  comment?: string;      // camelCase
  
  reviewer?: {
    id: number;
    name: string;
    avatar_url?: string; // snake_case
  };
  // ...
}
```

#### 修改后（完全 camelCase）

```typescript
export interface Review extends BaseEntity {
  orderId: number;       // camelCase
  reviewerId: number;    // camelCase
  playerId: number;      // camelCase
  rating: number;        // camelCase
  comment?: string;      // camelCase
  
  reviewer?: {
    id: number;
    name: string;
    avatarUrl?: string;  // camelCase
  };
  // ...
}
```

### 3. 前端组件更新

**文件：** `frontend/src/pages/Reviews/ReviewList.tsx`

#### 字段使用更新

```typescript
// 修改前
<div className={styles.orderId}>ID: {record.order_id}</div>
<div className={styles.reviewerId}>ID: {record.reviewer_id}</div>
<div className={styles.playerId}>ID: {record.player_id}</div>

// 修改后
<div className={styles.orderId}>ID: {record.orderId}</div>
<div className={styles.reviewerId}>ID: {record.reviewerId}</div>
<div className={styles.playerId}>ID: {record.playerId}</div>
```

**文件：** `frontend/src/pages/Reviews/ReviewFormModal.tsx`

✅ 无需修改 - 已经使用 `rating` 和 `comment`（camelCase）

## API 响应示例

### 修改前

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "order_id": 1,
      "user_id": 2,
      "player_id": 1,
      "score": 5,
      "content": "很满意的陪玩体验..."
    }
  ]
}
```

### 修改后

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "orderId": 1,
      "reviewerId": 2,
      "playerId": 1,
      "rating": 5,
      "comment": "很满意的陪玩体验..."
    }
  ]
}
```

## 技术要点

### 1. GORM column 标签的作用

当 JSON 字段名与数据库列名不同时，必须使用 `column` 标签：

```go
UserID uint64 `json:"reviewerId" gorm:"column:user_id"`
```

- **json:"reviewerId"** - API 响应时使用 camelCase
- **gorm:"column:user_id"** - 数据库查询时使用 snake_case

### 2. 数据库兼容性

通过 `column` 标签，确保：
- ✅ API 层使用 camelCase（符合前端规范）
- ✅ 数据库层使用 snake_case（符合数据库规范）
- ✅ 无需修改数据库表结构
- ✅ 无需数据迁移

### 3. TypeScript 类型安全

前端使用完全 camelCase 后：
- ✅ 编辑器自动补全正确
- ✅ TypeScript 类型检查生效
- ✅ 重构时可以追踪所有使用点

## 修改文件清单

### 后端
1. ✅ `backend/internal/model/review.go` - Review 模型 JSON 标签

### 前端
1. ✅ `frontend/src/types/review.ts` - Review 类型定义
2. ✅ `frontend/src/pages/Reviews/ReviewList.tsx` - 评价列表组件
3. ✅ `frontend/src/pages/Reviews/ReviewFormModal.tsx` - 无需修改（已是camelCase）

### 文档
1. ✅ `CAMELCASE_NAMING_UNIFICATION.md` - 本文档

## 测试验证

### 后端测试

```bash
# Windows PowerShell
cd backend
go build ./...
go test ./...
```

### 前端测试

访问评价管理页面：`http://localhost:5173/reviews`

**验证项：**
- ✅ 评分列显示星级（⭐⭐⭐⭐⭐）
- ✅ 订单ID 正确显示
- ✅ 评价人ID 正确显示
- ✅ 陪玩师ID 正确显示
- ✅ 评价内容正确显示
- ✅ 编辑功能正常
- ✅ 删除功能正常

### API 测试

```bash
# 获取评价列表
curl http://localhost:8080/api/v1/admin/reviews

# 验证返回的字段名是 camelCase
# orderId, reviewerId, playerId, rating, comment
```

## 未来规范

### 新增模块命名规范

1. **Go 后端**
```go
type Example struct {
    FieldName string `json:"fieldName" gorm:"column:field_name"`
}
```

2. **TypeScript 前端**
```typescript
interface Example {
  fieldName: string;
}
```

### 统一检查清单

创建新 API 时，确保：
- [ ] Go 结构体使用 PascalCase
- [ ] JSON 标签使用 camelCase
- [ ] GORM column 使用 snake_case
- [ ] TypeScript 类型使用 camelCase
- [ ] 前端组件使用 camelCase

## 推广计划

建议将此规范应用到所有模块：

### 阶段1：核心模块（优先）
- [x] Review（评价）- 已完成
- [ ] Order（订单）
- [ ] User（用户）
- [ ] Player（陪玩师）
- [ ] Game（游戏）

### 阶段2：扩展模块
- [ ] Payment（支付）
- [ ] Stats（统计）
- [ ] OperationLog（操作日志）

### 阶段3：关联字段
统一所有嵌套对象的字段名：
```typescript
// 统一前
avatar_url, order_no, game_id

// 统一后
avatarUrl, orderNo, gameId
```

## 优势总结

### 1. 一致性
- 前后端命名规范统一
- 代码可读性提高
- 团队协作更顺畅

### 2. 可维护性
- TypeScript 类型检查生效
- IDE 自动补全准确
- 重构安全性提高

### 3. 最佳实践
- 符合 JavaScript/TypeScript 生态规范
- 符合 RESTful API 设计规范
- 符合现代 Web 开发标准

### 4. 开发效率
- 减少命名转换的心智负担
- 减少因命名错误导致的 bug
- 提高代码审查效率

## 总结

通过本次统一，评价模块完全采用 camelCase 命名规范：

1. **后端**：JSON 标签使用 camelCase，数据库列名保持 snake_case
2. **前端**：TypeScript 类型和组件使用 camelCase
3. **API**：响应数据使用 camelCase

这为其他模块的命名统一提供了标准模板，建议逐步推广到整个项目。

