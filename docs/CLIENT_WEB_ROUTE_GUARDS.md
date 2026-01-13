# GameLink Desktop Client - Route Guards & Navigation

> **Route Protection & Navigation System** - JWT authentication, role-based access control, and navigation patterns

---

## Document Overview

```
┌─────────────────────────────────────────────────────────────┐
│          GAMELINK CLIENT - ROUTE GUARDS & NAVIGATION       │
├─────────────────────────────────────────────────────────────┤
│  User Roles: Guest | User | Player                         │
│  Auth Method: JWT (Access + Refresh Token)                 │
│  Router: Vue Router 4 / React Router 6                     │
│  State: Pinia / Zustand                                    │
└─────────────────────────────────────────────────────────────┘
```

---

## Table of Contents

1. [5-Layer Guard Architecture](#1-5-layer-guard-architecture)
2. [Route Configuration](#2-route-configuration)
3. [Authentication Flow](#3-authentication-flow)
4. [Role-Based Access Control](#4-role-based-access-control)
5. [Navigation Patterns](#5-navigation-patterns)
6. [Security Best Practices](#6-security-best-practices)

---

## 1. 5-Layer Guard Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    LAYER 1: GLOBAL GUARD                        │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  • Token refresh (auto-refresh before expiration)       │   │
│  │  • Authentication check (is user logged in?)           │   │
│  │  • Role verification (does user have required role?)   │   │
│  │  • Page permission check (can user access this page?)  │   │
│  │  • Redirect handling (save/restore redirect intent)    │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    LAYER 2: AUTH GUARD                          │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  • Login state management                              │   │
│  │  • Logout flow (clear tokens, redirect to login)       │   │
│  │  • Session timeout handling                            │   │
│  │  • Concurrent login prevention (optional)              │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    LAYER 3: ROLE GUARD                          │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  • Guest routes (no auth required)                      │   │
│  │  │   - /login, /register, /forgot-password             │   │
│  │  • User routes (requires auth)                         │   │
│  │  │   - /orders, /profile, /chat                        │   │
│  │  • Player routes (requires player role)                │   │
│  │  │   - /player/dashboard, /player/schedule             │   │
│  │  • Admin routes (requires admin role)                  │   │
│  │      - /admin/*                                        │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    LAYER 4: PAGE GUARD                          │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  • Data fetching (load page data before rendering)     │   │
│  │  • Loading states (skeleton screens)                   │   │
│  │  • Error handling (404, 403, 500)                      │   │
│  │  • Page permissions (fine-grained access control)      │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    LAYER 5: COMPONENT GUARD                     │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  • Permission-based UI rendering                       │   │
│  │  • Button/feature visibility by permission             │   │
│  │  • v-permission directive (Vue)                        │   │
│  │  • <Protected> component (React)                       │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 2. Route Configuration

### 2.1 Route Meta Schema

```typescript
interface RouteMeta {
  // Authentication
  requiresAuth?: boolean        // Default: false
  guestOnly?: boolean           // Redirect authenticated users away

  // Role-based access
  roles?: UserRole[]            // ['guest'] | ['user'] | ['player'] | ['admin']

  // Permissions
  permissions?: Permission[]    // Fine-grained permissions: ['order.create', 'player.update']

  // Layout
  layout?: 'default' | 'auth' | 'player' | 'admin'

  // Page title
  title?: string                // Page title (browser tab)

  // SEO
  description?: string          // Meta description
}

enum UserRole {
  GUEST = 'guest',
  USER = 'user',
  PLAYER = 'player',
  ADMIN = 'admin'
}
```

### 2.2 Complete Route Configuration (Vue Router 4)

```typescript
// router/routes.ts
import type { RouteRecordRaw } from 'vue-router'

export const routes: RouteRecordRaw[] = [
  // ========================================
  // PUBLIC / GUEST ROUTES (No auth required)
  // ========================================
  {
    path: '/login',
    name: 'login',
    component: () => import('@/pages/auth/LoginPage.vue'),
    meta: {
      guestOnly: true,           // Redirect authenticated users to home
      layout: 'auth',
      title: '登录 - GameLink'
    }
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('@/pages/auth/RegisterPage.vue'),
    meta: {
      guestOnly: true,
      layout: 'auth',
      title: '注册 - GameLink'
    }
  },
  {
    path: '/forgot-password',
    name: 'forgot-password',
    component: () => import('@/pages/auth/ForgotPasswordPage.vue'),
    meta: {
      guestOnly: true,
      layout: 'auth',
      title: '忘记密码 - GameLink'
    }
  },

  // ========================================
  // USER ROUTES (Requires auth)
  // ========================================
  {
    path: '/',
    name: 'home',
    component: () => import('@/pages/home/HomePage.vue'),
    meta: {
      requiresAuth: true,
      roles: [UserRole.USER, UserRole.PLAYER],
      layout: 'default',
      title: '首页 - GameLink'
    }
  },
  {
    path: '/players',
    name: 'players',
    component: () => import('@/pages/players/PlayerListPage.vue'),
    meta: {
      requiresAuth: true,
      roles: [UserRole.USER],
      layout: 'default',
      title: '陪玩师列表 - GameLink'
    }
  },
  {
    path: '/players/:id',
    name: 'player-detail',
    component: () => import('@/pages/players/PlayerDetailPage.vue'),
    meta: {
      requiresAuth: true,
      roles: [UserRole.USER],
      layout: 'default',
      title: '陪玩师详情 - GameLink'
    }
  },

  // Order routes
  {
    path: '/orders',
    name: 'orders',
    component: () => import('@/pages/orders/OrderListPage.vue'),
    meta: {
      requiresAuth: true,
      roles: [UserRole.USER],
      layout: 'default',
      title: '我的订单 - GameLink'
    }
  },
  {
    path: '/orders/create',
    name: 'order-create',
    component: () => import('@/pages/orders/OrderCreatePage.vue'),
    meta: {
      requiresAuth: true,
      roles: [UserRole.USER],
      permissions: ['order.create'],
      layout: 'default',
      title: '创建订单 - GameLink'
    }
  },
  {
    path: '/orders/:id',
    name: 'order-detail',
    component: () => import('@/pages/orders/OrderDetailPage.vue'),
    meta: {
      requiresAuth: true,
      roles: [UserRole.USER],
      layout: 'default',
      title: '订单详情 - GameLink'
    }
  },

  // Chat routes
  {
    path: '/chat',
    name: 'chat',
    component: () => import('@/pages/chat/ChatPage.vue'),
    meta: {
      requiresAuth: true,
      roles: [UserRole.USER, UserRole.PLAYER],
      layout: 'default',
      title: '消息中心 - GameLink'
    }
  },
  {
    path: '/chat/:userId',
    name: 'chat-conversation',
    component: () => import('@/pages/chat/ChatConversationPage.vue'),
    meta: {
      requiresAuth: true,
      roles: [UserRole.USER, UserRole.PLAYER],
      layout: 'default',
      title: '聊天 - GameLink'
    }
  },

  // Profile routes
  {
    path: '/profile',
    name: 'profile',
    component: () => import('@/pages/profile/ProfilePage.vue'),
    meta: {
      requiresAuth: true,
      roles: [UserRole.USER, UserRole.PLAYER],
      layout: 'default',
      title: '个人中心 - GameLink'
    }
  },
  {
    path: '/wallet',
    name: 'wallet',
    component: () => import('@/pages/profile/WalletPage.vue'),
    meta: {
      requiresAuth: true,
      roles: [UserRole.USER, UserRole.PLAYER],
      layout: 'default',
      title: '我的钱包 - GameLink'
    }
  },

  // ========================================
  // PLAYER ROUTES (Requires player role)
  // ========================================
  {
    path: '/player/dashboard',
    name: 'player-dashboard',
    component: () => import('@/pages/player/DashboardPage.vue'),
    meta: {
      requiresAuth: true,
      roles: [UserRole.PLAYER],
      layout: 'player',
      title: '陪玩师后台 - GameLink'
    }
  },
  {
    path: '/player/schedule',
    name: 'player-schedule',
    component: () => import('@/pages/player/SchedulePage.vue'),
    meta: {
      requiresAuth: true,
      roles: [UserRole.PLAYER],
      layout: 'player',
      title: '排班管理 - GameLink'
    }
  },
  {
    path: '/player/earnings',
    name: 'player-earnings',
    component: () => import('@/pages/player/EarningsPage.vue'),
    meta: {
      requiresAuth: true,
      roles: [UserRole.PLAYER],
      layout: 'player',
      title: '收益管理 - GameLink'
    }
  },

  // ========================================
  // ADMIN ROUTES (Requires admin role)
  // ========================================
  {
    path: '/admin',
    name: 'admin',
    component: () => import('@/pages/admin/DashboardPage.vue'),
    meta: {
      requiresAuth: true,
      roles: [UserRole.ADMIN],
      layout: 'admin',
      title: '管理后台 - GameLink'
    }
  },
  {
    path: '/admin/users',
    name: 'admin-users',
    component: () => import('@/pages/admin/user/UserListPage.vue'),
    meta: {
      requiresAuth: true,
      roles: [UserRole.ADMIN],
      permissions: ['user.read'],
      layout: 'admin',
      title: '用户管理 - GameLink'
    }
  },
  {
    path: '/admin/players',
    name: 'admin-players',
    component: () => import('@/pages/admin/player/PlayerListPage.vue'),
    meta: {
      requiresAuth: true,
      roles: [UserRole.ADMIN],
      permissions: ['player.read'],
      layout: 'admin',
      title: '陪玩师管理 - GameLink'
    }
  },
  {
    path: '/admin/orders',
    name: 'admin-orders',
    component: () => import('@/pages/admin/order/OrderListPage.vue'),
    meta: {
      requiresAuth: true,
      roles: [UserRole.ADMIN],
      permissions: ['order.read'],
      layout: 'admin',
      title: '订单管理 - GameLink'
    }
  },

  // ========================================
  // ERROR PAGES
  // ========================================
  {
    path: '/403',
    name: 'forbidden',
    component: () => import('@/pages/error/ForbiddenPage.vue'),
    meta: {
      title: '无权访问 - GameLink'
    }
  },
  {
    path: '/404',
    name: 'not-found',
    component: () => import('@/pages/error/NotFoundPage.vue'),
    meta: {
      title: '页面不存在 - GameLink'
    }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404'
  }
]
```

---

## 3. Authentication Flow

### 3.1 Token Management Strategy

```typescript
// stores/auth.ts (Pinia)
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  // ========================================
  // STATE
  // ========================================
  const accessToken = ref<string | null>(null)
  const refreshTokenValue = ref<string | null>(null)
  const user = ref<User | null>(null)
  const roles = ref<UserRole[]>([])

  // Token expiration (decoded from JWT)
  const expiresAt = ref<number | null>(null)

  // Concurrent refresh lock (prevent multiple refresh calls)
  let refreshPromise: Promise<void> | null = null

  // ========================================
  // COMPUTED
  // ========================================
  const isAuthenticated = computed(() => {
    return !!accessToken.value && !!user.value
  })

  const hasRole = computed(() => {
    return (role: UserRole) => roles.value.includes(role)
  })

  const hasPermission = computed(() => {
    return (permission: Permission) => {
      // Implement permission check logic
      // Can be based on user.permissions array or role-based rules
      return user.value?.permissions?.includes(permission) ?? false
    }
  })

  const isTokenExpiringSoon = computed(() => {
    if (!expiresAt.value) return true
    // Refresh if token expires in less than 5 minutes
    return Date.now() > expiresAt.value - 5 * 60 * 1000
  })

  // ========================================
  // ACTIONS
  // ========================================
  async function login(credentials: LoginCredentials) {
    const response = await api.post('/auth/login', credentials)

    // Store tokens
    accessToken.value = response.data.accessToken
    refreshTokenValue.value = response.data.refreshToken

    // Store refresh token in localStorage (survives page reload)
    localStorage.setItem('refresh_token', response.data.refreshToken)

    // Decode JWT to get expiration
    const decoded = decodeJWT(response.data.accessToken)
    expiresAt.value = decoded.exp * 1000

    // Fetch user profile
    await fetchUser()
  }

  async function logout() {
    try {
      // Call backend logout endpoint (invalidate refresh token)
      await api.post('/auth/logout', {
        refreshToken: refreshTokenValue.value
      })
    } finally {
      // Clear state regardless of API call success
      clearAuth()
      router.push({ name: 'login' })
    }
  }

  async function refreshToken() {
    // Prevent concurrent refresh attempts
    if (refreshPromise) {
      return refreshPromise
    }

    refreshPromise = (async () => {
      try {
        const response = await api.post('/auth/refresh', {
          refreshToken: refreshTokenValue.value
        })

        accessToken.value = response.data.accessToken
        refreshTokenValue.value = response.data.refreshToken

        // Update localStorage
        localStorage.setItem('refresh_token', response.data.refreshToken)

        // Update expiration
        const decoded = decodeJWT(response.data.accessToken)
        expiresAt.value = decoded.exp * 1000
      } catch (error) {
        // Refresh failed - user must login again
        clearAuth()
        router.push({ name: 'login', query: { sessionExpired: 'true' } })
        throw error
      } finally {
        refreshPromise = null
      }
    })()

    return refreshPromise
  }

  async function fetchUser() {
    if (!accessToken.value) return

    const response = await api.get('/users/me')
    user.value = response.data
    roles.value = response.data.roles
  }

  function clearAuth() {
    accessToken.value = null
    refreshTokenValue.value = null
    user.value = null
    roles.value = []
    expiresAt.value = null
    localStorage.removeItem('refresh_token')
  }

  // Initialize from localStorage on app load
  function initialize() {
    const storedRefreshToken = localStorage.getItem('refresh_token')
    if (storedRefreshToken) {
      refreshTokenValue.value = storedRefreshToken
      // Try to refresh tokens on app load
      return refreshToken()
    }
  }

  return {
    // State
    accessToken,
    user,
    roles,

    // Computed
    isAuthenticated,
    hasRole,
    hasPermission,
    isTokenExpiringSoon,

    // Actions
    login,
    logout,
    refreshToken,
    fetchUser,
    clearAuth,
    initialize
  }
})
```

### 3.2 Global Router Guard

```typescript
// router/guards.ts
import type { Router } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { UserRole } from '@/types/user'

export function setupRouterGuards(router: Router) {
  router.beforeEach(async (to, from, next) => {
    const authStore = useAuthStore()

    // ========================================
    // STEP 1: Token Refresh (if expiring soon)
    // ========================================
    if (authStore.isTokenExpiringSoon && authStore.refreshTokenValue) {
      try {
        await authStore.refreshToken()
      } catch {
        // Refresh failed - clear auth and redirect to login
        authStore.clearAuth()
        return next({
          name: 'login',
          query: { sessionExpired: 'true', redirect: to.fullPath }
        })
      }
    }

    // ========================================
    // STEP 2: Guest-only Routes (redirect if authenticated)
    // ========================================
    if (to.meta.guestOnly) {
      if (authStore.isAuthenticated) {
        // User is already authenticated, redirect to home
        return next({ name: 'home' })
      }
      return next()
    }

    // ========================================
    // STEP 3: Authentication Check
    // ========================================
    if (to.meta.requiresAuth && !authStore.isAuthenticated) {
      // Save redirect intent
      return next({
        name: 'login',
        query: { redirect: to.fullPath }
      })
    }

    // ========================================
    // STEP 4: Role Check
    // ========================================
    if (to.meta.roles && to.meta.roles.length > 0) {
      const hasRequiredRole = to.meta.roles.some(role =>
        authStore.hasRole(role)
      )

      if (!hasRequiredRole) {
        // User doesn't have required role
        return next({ name: 'forbidden' })
      }
    }

    // ========================================
    // STEP 5: Permission Check
    // ========================================
    if (to.meta.permissions && to.meta.permissions.length > 0) {
      const hasRequiredPermission = to.meta.permissions.every(permission =>
        authStore.hasPermission(permission)
      )

      if (!hasRequiredPermission) {
        // User doesn't have required permissions
        return next({ name: 'forbidden' })
      }
    }

    // ========================================
    // STEP 6: Page Title
    // ========================================
    if (to.meta.title) {
      document.title = to.meta.title
    }

    // All checks passed - proceed to route
    next()
  })

  // ========================================
  // Navigation Error Handler
  // ========================================
  router.onError((error) => {
    console.error('Router error:', error)

    // Handle navigation errors (e.g., failed async component load)
    if (error.name === 'ChunkLoadError') {
      // Lazy-loaded chunk failed - suggest refresh
      console.error('Failed to load page chunk. Please refresh.')
    }
  })
}
```

### 3.3 Redirect Handler

```typescript
// composables/useRedirect.ts
import { useRouter, useRoute } from 'vue-router'
import { watchEffect } from 'vue'

export function useRedirect() {
  const router = useRouter()
  const route = useRoute()

  // Handle redirect after login
  function handlePostLoginRedirect() {
    const redirect = route.query.redirect as string | undefined

    if (redirect) {
      router.push(redirect)
    } else {
      // Default redirect based on user role
      const authStore = useAuthStore()

      if (authStore.hasRole(UserRole.ADMIN)) {
        router.push({ name: 'admin' })
      } else if (authStore.hasRole(UserRole.PLAYER)) {
        router.push({ name: 'player-dashboard' })
      } else {
        router.push({ name: 'home' })
      }
    }
  }

  return {
    handlePostLoginRedirect
  }
}
```

---

## 4. Role-Based Access Control

### 4.1 Permission Directive (Vue 3)

```typescript
// directives/permission.ts
import type { Directive, DirectiveBinding } from 'vue'
import { useAuthStore } from '@/stores/auth'

export const permission: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    const { value } = binding
    const authStore = useAuthStore()

    if (value && typeof value === 'string') {
      const hasPermission = authStore.hasPermission(value)

      if (!hasPermission) {
        // Remove element from DOM
        el.parentNode?.removeChild(el)
      }
    } else {
      throw new Error('v-permission requires a permission string value')
    }
  }
}

// Register directive
// main.ts
app.directive('permission', permission)
```

### 4.2 Permission Component

```vue
<!-- components/common/Protected.vue -->
<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

interface Props {
  permission?: string
  role?: UserRole
  roles?: UserRole[]
  fallback?: 'hide' | 'disable'
}

const props = withDefaults(defineProps<Props>(), {
  fallback: 'hide'
})

const authStore = useAuthStore()

const isAllowed = computed(() => {
  if (props.permission) {
    return authStore.hasPermission(props.permission)
  }

  if (props.roles) {
    return props.roles.some(role => authStore.hasRole(role))
  }

  if (props.role) {
    return authStore.hasRole(props.role)
  }

  return true
})
</script>

<template>
  <div v-if="isAllowed" v-bind="$attrs">
    <slot />
  </div>
  <div v-else-if="fallback === 'disable'" class="disabled" v-bind="$attrs">
    <slot />
  </div>
</template>

<style scoped>
.disabled {
  opacity: 0.5;
  pointer-events: none;
  cursor: not-allowed;
}
</style>
```

### 4.3 Usage Examples

```vue
<template>
  <!-- Using directive -->
  <button v-permission="'order.create'">
    创建订单
  </button>

  <!-- Using component -->
  <Protected permission="order.delete">
    <button @click="deleteOrder">删除订单</button>
  </Protected>

  <!-- Using component with role -->
  <Protected :roles="[UserRole.ADMIN, UserRole.MODERATOR]">
    <button>管理操作</button>
  </Protected>

  <!-- Using component with disabled fallback -->
  <Protected permission="order.update" fallback="disable">
    <button>编辑订单</button>
  </Protected>
</template>
```

---

## 5. Navigation Patterns

### 5.1 Programmatic Navigation

```typescript
// Navigate with role-based redirect
import { useRouter } from 'vue-router'

function navigateAfterLogin() {
  const router = useRouter()
  const authStore = useAuthStore()

  if (authStore.hasRole(UserRole.ADMIN)) {
    router.push({ name: 'admin' })
  } else if (authStore.hasRole(UserRole.PLAYER)) {
    router.push({ name: 'player-dashboard' })
  } else {
    router.push({ name: 'home' })
  }
}

// Navigate with permission check
function navigateToAdmin() {
  const router = useRouter()
  const authStore = useAuthStore()

  if (authStore.hasPermission('admin.access')) {
    router.push({ name: 'admin' })
  } else {
    // Show notification: no permission
    showError('您没有权限访问管理后台')
  }
}
```

### 5.2 Navigation Guards in Components

```vue
<script setup lang="ts">
import { onBeforeRouteLeave, onBeforeRouteUpdate } from 'vue-router'
import { ref } from 'vue'

const hasUnsavedChanges = ref(false)

// Warn before leaving with unsaved changes
onBeforeRouteLeave((to, from, next) => {
  if (hasUnsavedChanges.value) {
    const answer = window.confirm(
      '您有未保存的更改，确定要离开吗？'
    )

    if (answer) {
      next()
    } else {
      next(false)
    }
  } else {
    next()
  }
})

// Handle route updates (same route, different params)
onBeforeRouteUpdate((to, from, next) => {
  // Reload data when params change
  fetchData(to.params.id)
  next()
})
</script>
```

### 5.3 Back Navigation with History

```typescript
// composables/useNavigation.ts
import { useRouter } from 'vue-router'
import { computed } from 'vue'

export function useNavigation() {
  const router = useRouter()

  const canGoBack = computed(() => {
    return window.history.state.back !== null
  })

  function goBack(fallback = { name: 'home' }) {
    if (canGoBack.value) {
      router.back()
    } else {
      router.push(fallback)
    }
  }

  return {
    canGoBack,
    goBack
  }
}
```

---

## 6. Security Best Practices

### 6.1 Token Storage

```typescript
// ✅ GOOD: Access token in memory, refresh token in localStorage
const accessToken = ref<string | null>(null)  // Memory (cleared on refresh)
const refreshToken = localStorage.getItem('refresh_token')  // Survives refresh

// ❌ BAD: Access token in localStorage (vulnerable to XSS)
const accessToken = localStorage.getItem('access_token')
```

### 6.2 XSS Prevention

```typescript
// Sanitize user input before rendering
import DOMPurify from 'dompurify'

function renderUserContent(content: string) {
  return DOMPurify.sanitize(content)
}
```

### 6.3 CSRF Protection

```typescript
// Include CSRF token in POST requests
const csrfToken = getCsrfToken() // From meta tag or cookie

api.post('/orders', data, {
  headers: {
    'X-CSRF-Token': csrfToken
  }
})
```

### 6.4 Token Expiration Handling

```typescript
// Proactive refresh (before expiration)
if (isTokenExpiringSoon(token)) {
  await refreshToken()
}

// Handle expired token during API call
api.interceptors.response.use(
  response => response,
  async error => {
    if (error.response?.status === 401) {
      // Token expired - attempt refresh
      try {
        await authStore.refreshToken()
        // Retry original request
        return api.request(error.config)
      } catch {
        // Refresh failed - redirect to login
        authStore.clearAuth()
        router.push({ name: 'login', query: { sessionExpired: 'true' } })
      }
    }
    return Promise.reject(error)
  }
)
```

### 6.5 Security Headers

```typescript
// Set security meta tags
function setSecurityHeaders() {
  const metaTags = [
    { name: 'X-Content-Type-Options', content: 'nosniff' },
    { name: 'X-Frame-Options', content: 'DENY' },
    { name: 'X-XSS-Protection', content: '1; mode=block' },
    { name: 'Referrer-Policy', content: 'strict-origin-when-cross-origin' },
    { name: 'Content-Security-Policy', content: "default-src 'self'" }
  ]

  metaTags.forEach(tag => {
    const meta = document.createElement('meta')
    Object.assign(meta, tag)
    document.head.appendChild(meta)
  })
}
```

---

**Last Updated**: 2025-01-13
**Version**: 1.0.0
**Status**: Complete
