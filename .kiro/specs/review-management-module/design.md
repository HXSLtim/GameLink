# 设计文档 - 评价管理模块

## 概述

评价管理模块负责管理 GameLink 平台的用户评价内容，包括评价审核、举报处理、敏感词过滤等功能。该模块采用前后端分离架构，确保评价内容的健康性和真实性，维护良好的平台评价环境。

## 架构

### 系统架构图

```
┌─────────────────────────────────────────────────────────────┐
│                        前端层 (React)                        │
├─────────────────────────────────────────────────────────────┤
│  评价列表  │  评价审核  │  举报管理  │  敏感词管理  │  统计分析  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓ HTTP/HTTPS
┌─────────────────────────────────────────────────────────────┐
│                      API 层 (Go Backend)                     │
├─────────────────────────────────────────────────────────────┤
│  评价API  │  审核API  │  举报API  │  敏感词API  │  统计API    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                      业务逻辑层                              │
├─────────────────────────────────────────────────────────────┤
│  评价服务  │  审核服务  │  举报服务  │  敏感词服务  │  统计服务  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                      数据存储层                              │
├─────────────────────────────────────────────────────────────┤
│  PostgreSQL/SQLite  │  Redis (缓存)  │  MongoDB (日志)       │
│  (评价数据)         │                │                       │
└─────────────────────────────────────────────────────────────┘
```

## 组件和接口

### 前端组件

#### 1. ReviewList 组件
```typescript
interface ReviewListProps {
  onViewDetail: (reviewId: number) => void;
  onApprove: (reviewId: number) => void;
  onReject: (reviewId: number) => void;
}
```

#### 2. ReviewDetail 组件
```typescript
interface ReviewDetailProps {
  reviewId: number;
  onReply: (content: string) => Promise<void>;
  onDelete: () => Promise<void>;
}
```

#### 3. SensitiveWordManager 组件
```typescript
interface SensitiveWordManagerProps {
  onAdd: (word: SensitiveWord) => Promise<void>;
  onUpdate: (id: number, word: SensitiveWord) => Promise<void>;
  onDelete: (id: number) => Promise<void>;
}
```

### 后端接口

```go
type ReviewService interface {
    ListReviews(ctx context.Context, req *ListReviewsRequest) (*ListReviewsResponse, error)
    GetReview(ctx context.Context, reviewID int64) (*Review, error)
    ApproveReview(ctx context.Context, reviewID int64) error
    RejectReview(ctx context.Context, reviewID int64, reason string) error
    DeleteReview(ctx context.Context, reviewID int64) error
    ReplyReview(ctx context.Context, reviewID int64, content string) error
}
```

## 数据模型

### Review (评价)
```typescript
interface Review {
  id: number;
  orderID: number;
  reviewerID: number;
  playerID: number;
  score: number;              // 1-5分
  content: string;
  images: string[];
  status: ReviewStatus;
  isReported: boolean;
  createdAt: Date;
  updatedAt: Date;
}

enum ReviewStatus {
  PENDING = 'pending',
  APPROVED = 'approved',
  REJECTED = 'rejected',
  DELETED = 'deleted'
}
```

### SensitiveWord (敏感词)
```typescript
interface SensitiveWord {
  id: number;
  word: string;
  category: SensitiveWordCategory;
  severity: SensitiveWordSeverity;
  createdAt: Date;
}

enum SensitiveWordCategory {
  POLITICAL = 'political',
  PORNOGRAPHIC = 'pornographic',
  VIOLENT = 'violent',
  ADVERTISING = 'advertising',
  OTHER = 'other'
}

enum SensitiveWordSeverity {
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high'
}
```

## 正确性属性

### 属性 1：评价记录完整性
*对于任何*评价列表请求，返回的所有评价记录必须包含完整的必填字段（订单ID、评价者ID、被评价者ID、评分、内容、状态），且字段值不为空或null
**验证：需求 1.1**

### 属性 2：评分范围约束
*对于任何*评价记录，评分必须在1到5之间（包含1和5）
**验证：需求 1.1**

### 属性 3：审核状态转换合法性
*对于任何*评价审核操作，状态转换必须遵循合法路径：pending → approved 或 pending → rejected，已审核的评价不能再次审核
**验证：需求 2.2, 2.3**

### 属性 4：敏感词唯一性
*对于任何*敏感词添加操作，新添加的敏感词必须在敏感词库中不存在
**验证：需求 5.3**

### 属性 5：举报评价标记
*对于任何*被举报的评价，isReported字段必须为true
**验证：需求 1.5, 3.1**

### 属性 6：评价排序一致性
*对于任何*评价列表查询，当按创建时间降序排列时，第一条记录的创建时间必须大于或等于最后一条记录的创建时间
**验证：需求 1.1**

### 属性 7：敏感词检测准确性
*对于任何*包含敏感词的评价内容，系统必须检测到敏感词并标记评价为待审核状态
**验证：需求 2.4, 5.5**

### 属性 8：权限验证一致性
*对于任何*需要权限的操作，系统必须先验证用户权限，验证失败时必须返回403错误且不执行任何业务逻辑
**验证：需求 10.1, 10.2, 10.3, 10.4**

### 属性 9：操作日志完整性
*对于任何*评价相关操作（创建、审核、删除、回复），系统必须记录操作日志，包含操作类型、操作人、操作时间、操作前状态和操作后状态
**验证：需求 9.1**

### 属性 10：批量操作原子性
*对于任何*批量审核操作，要么所有评价都成功审核，要么所有评价都保持原状态
**验证：需求 2.5**

## 错误处理

### 错误分类
- **400 Bad Request**: 请求参数错误（评分超出范围、内容为空等）
- **401 Unauthorized**: 未认证
- **403 Forbidden**: 无权限
- **404 Not Found**: 评价不存在
- **500 Internal Server Error**: 服务器内部错误

## 测试策略

### 单元测试
使用 Vitest + Testing Library 进行组件和工具函数测试

### 属性测试
使用 fast-check 验证系统的通用正确性属性，每个属性测试运行至少100次迭代

### 集成测试
测试评价审核完整流程、举报处理流程、敏感词过滤流程

### E2E测试
使用 Playwright 进行端到端测试
