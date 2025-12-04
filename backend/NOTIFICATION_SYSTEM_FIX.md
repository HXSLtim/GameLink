# 通知系统问题修复说明

## 已修复的问题

### 1. 前端代理配置错误
- **问题**: Vite代理配置为 `/api`，导致请求转发到后端时丢失了 `/api/v1` 前缀
- **修复**: 将代理配置改为 `/api/v1`，确保请求正确转发到后端
- **文件**: `frontend/vite.config.ts`

### 2. 前端API路径错误
- **问题**: 前端使用的API路径包含错误的 `/user/` 前缀
- **修复**: 移除路径中的 `/user/` 前缀，改为直接使用 `/notifications`
- **文件**: `frontend/src/api/user.ts`

### 3. 参数名不匹配
- **问题**: 前端使用 `page_size`，后端期望 `pageSize` (camelCase)
- **修复**: 修改前端参数名为 `pageSize`
- **文件**: `frontend/src/pages/admin/Notifications/index.tsx`

### 4. 未实现的API功能
- **问题**: 前端实现了删除通知功能，但后端没有对应的API
- **修复**: 移除前端的删除按钮和相关代码
- **文件**: `frontend/src/pages/admin/Notifications/index.tsx`

### 5. API方法错误
- **问题**: 标记已读使用PUT方法，但后端实现为POST
- **修复**: 改为使用POST方法与后端保持一致
- **文件**: `frontend/src/api/user.ts`

## 通知系统功能

### 已实现的后端功能
- `GET /api/v1/notifications` - 获取通知列表（支持分页、未读过滤）
- `POST /api/v1/notifications/read` - 批量标记通知为已读
- `GET /api/v1/notifications/unread-count` - 获取未读通知数量

### 前端功能
- 通知列表展示，支持分页加载
- 标记单条通知为已读
- 标记所有通知为已读
- 未读通知高亮显示
- 时间格式化显示
- 空状态处理

## 测试步骤

1. 确保后端服务已启动并运行在 `http://localhost:8080`
2. 确保前端服务已启动并运行在 `http://localhost:5173`
3. 登录系统，查看系统右上角的通知图标
4. 访问通知中心页面查看通知列表
5. 测试标记已读和全部标记已读功能

## API请求示例

### 获取通知列表
```bash
curl -X GET "http://localhost:8080/api/v1/notifications?page=1&pageSize=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 标记通知为已读
```bash
curl -X POST "http://localhost:8080/api/v1/notifications/read" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"ids": [1, 2, 3]}'
```

### 获取未读通知数量
```bash
curl -X GET "http://localhost:8080/api/v1/notifications/unread-count" \
  -H "Authorization: Bearer YOUR_TOKEN"
```
