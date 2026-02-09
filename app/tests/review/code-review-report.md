# 代码审查报告

**审查人**: Mobile-Lead
**审查时间**: 2026-02-09
**审查范围**: 移动端 API 相关代码
**目的**: 为集成测试做准备，发现潜在问题

---

## 📋 审查项目

### 1. API 类型定义 ✅

**文件**: `app/src/types/`

**用户类型** (`user.ts`):
```typescript
export type UserRole = 'user' | 'player' | 'admin'
export type AppUserRole = Exclude<UserRole, 'admin'>
```
✅ **状态**: 清晰明确
✅ **建议**: 无

**订单类型** (`order.ts`):
```typescript
export type OrderStatus =
  | 'pending' | 'confirmed' | 'in_progress' | 'completed'
  | 'canceled' | 'refunding' | 'refunded' | 'disputed'

export type OrderPaymentMethod = 'wechat' | 'alipay' | 'wallet' | 'combined'
export type OrderPaymentStatus = 'success' | 'pending' | 'failed' | 'paid'
```
✅ **状态**: 类型定义完整
⚠️ **注意**: `OrderPaymentStatus` 中 `'success'` 和 `'paid'` 语义重叠，建议统一

**潜在问题**:
- `OrderPaymentStatus` 有两个相似的状态 (`'success'` 和 `'paid'`)
- **建议**: 统一使用 `'paid'` 表示已支付，`'success'` 用于其他场景

---

### 2. API 请求参数 ✅

**文件**: `app/src/api/auth.ts`

**登录请求**:
```typescript
export interface LoginRequest {
  phone?: string   // 手机号
  email?: string   // 邮箱
  password: string
}
```
✅ **状态**: 正确
✅ **说明**: 支持手机号或邮箱登录，与后端对齐

**微信登录请求**:
```typescript
export interface WeChatLoginRequest {
  code: string
  encryptedData?: string
  iv?: string
  referralCode?: string
}
```
✅ **状态**: 完整
✅ **说明**: 包含可选的加密数据和推荐码

**潜在问题**: 无

---

### 3. API 响应类型 ✅

**文件**: `app/src/api/auth.ts`

**登录响应**:
```typescript
export interface LoginResponse {
  token: string
  expires_at: string
  user: UserInfo
  // 兼容字段
  accessToken?: string
  refreshToken?: string
}
```
✅ **状态**: 正确
⚠️ **注意**: 字段名 `token` 和 `accessToken` 并存，使用 `normalizeUserInfo` 处理

**已处理**: `app/src/store/user.ts` 中的 `normalizeUserInfo` 函数处理了字段差异：
```typescript
const accessToken = data.accessToken || data.token || ''
```
✅ **状态**: 已有兼容处理

---

### 4. 错误处理 ✅

**文件**: `app/src/api/request.ts`

**错误类型**:
```typescript
export class ApiError extends Error {
  code: number
  constructor(message: string, code: number) {
    super(message)
    this.code = code
    this.name = 'ApiError'
  }
}
```
✅ **状态**: 清晰的错误类型定义

**HTTP 状态码处理**:
```typescript
if (statusCode === 401) {
  userStore.logout()
  reject(new ApiError('登录已过期，请重新登录', 401))
  return
}

if (statusCode === 403) {
  uni.showToast({ title: '无权限访问', icon: 'none' })
  reject(new ApiError('无权限访问', 403))
  return
}
```
✅ **状态**: 正确处理 401 和 403
✅ **说明**: 401 自动登出，403 显示提示

**潜在问题**: 无

---

### 5. 请求配置 ✅

**文件**: `app/src/api/request.ts`

**动态 BaseURL**:
```typescript
function resolveBaseUrl(): string {
  const envUrl = import.meta.env.VITE_API_BASE_URL
  if (envUrl) return envUrl

  // #ifdef H5
  const host = window.location.hostname || 'localhost'
  return `http://${host}:8080/api/v1`
  // #endif

  return 'http://localhost:8080/api/v1'
}
```
✅ **状态**: 灵活的配置
⚠️ **注意**: 默认端口 8080 与 BUG-001 冲突

**需要修改**:
```typescript
// 建议改为
return 'http://localhost:8000/api/v1'  // 使用后端专用端口
```

---

### 6. 分页参数 ✅

**文件**: `app/src/api/` (各模块)

**示例**:
```typescript
export interface OrderListParams {
  page?: number
  page_size?: number
  status?: OrderStatus | 'all'
}
```
✅ **状态**: 已对齐（之前任务 #16 完成）
✅ **说明**: 使用 `page_size` 而非 `pageSize`

**潜在问题**: 无

---

### 7. WebSocket 连接 ✅

**文件**: `app/src/composables/useWebSocket.ts`

**连接配置**:
```typescript
const wsUrl = computed(() => {
  const baseUrl = BASE_URL.replace('http://', 'ws://').replace('https://', 'wss://')
  return `${baseUrl}/ws`
})
```
✅ **状态**: 正确的 WebSocket URL 转换

**认证处理**:
```typescript
const connect = () => {
  const token = useUserStore().token
  ws.value = new WebSocket(`${wsUrl.value}?token=${token}`)
  // ...
}
```
✅ **状态**: Token 正确传递

**潜在问题**: 无

---

## 📊 审查总结

### ✅ 优点

1. **类型定义完整**: 所有 API 请求和响应都有明确的类型定义
2. **错误处理完善**: HTTP 状态码和业务错误都有处理
3. **字段兼容性好**: 使用 `normalizeUserInfo` 处理后端字段差异
4. **分页对齐**: 已统一使用 `page_size`
5. **WebSocket 配置正确**: URL 转换和认证处理正确

### ⚠️ 需要注意

1. **端口配置**: 默认 8080 端口需要改为 8000
2. **订单状态**: `OrderPaymentStatus` 中 `'success'` 和 `'paid'` 语义重叠

### 🔧 建议修改

#### 修改 1: API BaseURL (优先级 P0)

**文件**: `app/src/api/request.ts`

**当前代码**:
```typescript
return 'http://localhost:8080/api/v1'
```

**建议修改**:
```typescript
return 'http://localhost:8000/api/v1'
```

**原因**: 后端服务器改用 8000 端口，避免与前端 Vite (8080) 冲突

#### 修改 2: 订单支付状态 (优先级 P2)

**文件**: `app/src/types/order.ts`

**当前代码**:
```typescript
export type OrderPaymentStatus = 'success' | 'pending' | 'failed' | 'paid'
```

**建议修改**:
```typescript
export type OrderPaymentStatus = 'paid' | 'pending' | 'failed' | 'refunded'
```

**原因**: `'success'` 和 `'paid'` 语义重叠，统一使用 `'paid'` 表示已支付

---

## 🎯 测试准备就绪

基于代码审查结果，移动端代码质量良好，主要问题是端口配置。

**可以立即进行的测试**:
1. ✅ API 类型定义正确
2. ✅ 错误处理完善
3. ✅ 分页参数对齐
4. ✅ WebSocket 配置正确

**需要修改后**:
1. ⏸️ API BaseURL (等待 BUG-001 解决)
2. ⏸️ 开始实际 API 测试

---

**审查结论**: 代码质量良好，可以开始集成测试 (等待端口问题解决)
