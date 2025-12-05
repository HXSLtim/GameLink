# 设计文档 - 内容管理模块

## 概述

内容管理模块负责管理 GameLink 平台的用户生成内容（UGC），包括动态内容、聊天消息等。该模块采用前后端分离架构，通过敏感词过滤、内容审核、举报处理等机制，确保平台内容的健康性和安全性。

## 架构

### 系统架构图

```
┌─────────────────────────────────────────────────────────────┐
│                        前端层 (React)                        │
├─────────────────────────────────────────────────────────────┤
│  动态审核  │  聊天监控  │  举报管理  │  敏感词管理  │  统计分析  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓ HTTP/HTTPS
┌─────────────────────────────────────────────────────────────┐
│                      API 层 (Go Backend)                     │
├─────────────────────────────────────────────────────────────┤
│  动态API  │  聊天API  │  举报API  │  敏感词API  │  统计API    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                      业务逻辑层                              │
├─────────────────────────────────────────────────────────────┤
│  内容服务  │  审核服务  │  举报服务  │  敏感词服务  │  统计服务  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                      数据存储层                              │
├─────────────────────────────────────────────────────────────┤
│  PostgreSQL/SQLite  │  Redis (缓存)  │  MongoDB (日志)       │
│  (内容数据)         │                │                       │
└─────────────────────────────────────────────────────────────┘
```

## 组件和接口

### 前端组件

#### 1. FeedList 组件
```typescript
interface FeedListProps {
  onViewDetail: (feedId: number) => void;
  onApprove: (feedId: number) => void;
  onReject: (feedId: number) => void;
  onDelete: (feedId: number) => void;
}
```

#### 2. ChatMonitor 组件
```typescript
interface ChatMonitorProps {
  onViewGroup: (groupId: number) => void;
  onDeleteMessage: (messageId: number) => void;
  onMuteUser: (userId: number, duration: number) => void;
}
```

#### 3. ReportManager 组件
```typescript
interface ReportManagerProps {
  onViewDetail: (reportId: number) => void;
  onProcess: (reportId: number, action: ReportAction) => void;
}
```

### 后端接口

```go
type ContentService interface {
    ListFeeds(ctx context.Context, req *ListFeedsRequest) (*ListFeedsResponse, error)
    GetFeed(ctx context.Context, feedID int64) (*Feed, error)
    ApproveFeed(ctx context.Context, feedID int64) error
    RejectFeed(ctx context.Context, feedID int64, reason string) error
    DeleteFeed(ctx context.Context, feedID int64) error
    ListChatMessages(ctx context.Context, req *ListMessagesRequest) (*ListMessagesResponse, error)
    DeleteMessage(ctx context.Context, messageID int64) error
    MuteUser(ctx context.Context, userID int64, duration time.Duration) error
}
```

## 数据模型

### Feed (动态)
```typescript
interface Feed {
  id: number;
  userID: number;
  content: string;
  images: string[];
  visibility: FeedVisibility;
  status: FeedStatus;
  likesCount: number;
  commentsCount: number;
  createdAt: Date;
  updatedAt: Date;
}

enum FeedVisibility {
  PUBLIC = 'public',
  PRIVATE = 'private',
  FRIENDS = 'friends'
}

enum FeedStatus {
  PENDING = 'pending',
  APPROVED = 'approved',
  REJECTED = 'rejected',
  DELETED = 'deleted'
}
```

### ChatMessage (聊天消息)
```typescript
interface ChatMessage {
  id: number;
  senderID: number;
  receiverID?: number;
  groupID?: number;
  content: string;
  type: MessageType;
  isDeleted: boolean;
  createdAt: Date;
}

enum MessageType {
  TEXT = 'text',
  IMAGE = 'image',
  VOICE = 'voice',
  VIDEO = 'video'
}
```

### Report (举报)
```typescript
interface Report {
  id: number;
  reporterID: number;
  contentType: ContentType;
  contentID: number;
  reason: string;
  status: ReportStatus;
  createdAt: Date;
  processedAt?: Date;
}

enum ContentType {
  FEED = 'feed',
  CHAT_MESSAGE = 'chat_message',
  REVIEW = 'review'
}

enum ReportStatus {
  PENDING = 'pending',
  PROCESSED = 'processed',
  REJECTED = 'rejected'
}
```

## 正确性属性

### 属性 1：动态记录完整性
*对于任何*动态列表请求，返回的所有动态记录必须包含完整的必填字段（用户ID、内容、可见性、状态），且字段值不为空或null
**验证：需求 2.1**

### 属性 2：审核状态转换合法性
*对于任何*动态审核操作，状态转换必须遵循合法路径：pending → approved 或 pending → rejected，已审核的动态不能再次审核
**验证：需求 1.3, 1.4**

### 属性 3：敏感词检测准确性
*对于任何*包含敏感词的内容，系统必须检测到敏感词并标记内容为待审核状态
**验证：需求 1.5, 4.5**

### 属性 4：举报内容关联性
*对于任何*举报记录，被举报的内容必须存在且contentID必须有效
**验证：需求 5.1, 5.2**

### 属性 5：禁言时间有效性
*对于任何*禁言操作，禁言时长必须大于0且不超过30天
**验证：需求 3.5**

### 属性 6：批量操作原子性
*对于任何*批量审核操作，要么所有内容都成功审核，要么所有内容都保持原状态
**验证：需求 2.5**

### 属性 7：权限验证一致性
*对于任何*需要权限的操作，系统必须先验证用户权限，验证失败时必须返回403错误且不执行任何业务逻辑
**验证：需求 10.1, 10.2, 10.3, 10.4**

### 属性 8：操作日志完整性
*对于任何*内容相关操作（创建、审核、删除、举报），系统必须记录操作日志，包含操作类型、操作人、操作时间、操作前状态和操作后状态
**验证：需求 9.1**

### 属性 9：内容分类唯一性
*对于任何*内容分类添加操作，新添加的分类名称必须在分类列表中不存在
**验证：需求 6.3**

### 属性 10：敏感词唯一性
*对于任何*敏感词添加操作，新添加的敏感词必须在敏感词库中不存在
**验证：需求 4.3**

## 错误处理

### 错误分类
- **400 Bad Request**: 请求参数错误
- **401 Unauthorized**: 未认证
- **403 Forbidden**: 无权限
- **404 Not Found**: 内容不存在
- **500 Internal Server Error**: 服务器内部错误

## 测试策略

### 单元测试
使用 Vitest + Testing Library 进行组件和工具函数测试

### 属性测试
使用 fast-check 验证系统的通用正确性属性，每个属性测试运行至少100次迭代

### 集成测试
测试内容审核完整流程、举报处理流程、敏感词过滤流程

### E2E测试
使用 Playwright 进行端到端测试
