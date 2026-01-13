# GameLink Desktop Client - Page Design Principles

> **Page Component Design & Best Practices** - State management, loading patterns, error handling, and UI/UX standards

---

## Document Overview

```
┌─────────────────────────────────────────────────────────────┐
│          GAMELINK CLIENT - PAGE DESIGN PRINCIPLES          │
├─────────────────────────────────────────────────────────────┤
│  Framework: Vue 3 Composition API / React 19 Hooks          │
│  State: Pinia / Zustand                                    │
│  Design: GameLink Design System (Day/Night themes)         │
│  Spacing: 8pt Grid System                                  │
└─────────────────────────────────────────────────────────────┘
```

---

## Table of Contents

1. [Page State Machine](#1-page-state-machine)
2. [Component Structure](#2-component-structure)
3. [Data Fetching Patterns](#3-data-fetching-patterns)
4. [Loading States](#4-loading-states)
5. [Empty States](#5-empty-states)
6. [Error Handling](#6-error-handling)
7. [Permission-Based UI](#7-permission-based-ui)
8. [Performance Optimization](#8-performance-optimization)
9. [Accessibility Standards](#9-accessibility-standards)
10. [Page Checklist](#10-page-checklist)

---

## 1. Page State Machine

```
┌─────────────────────────────────────────────────────────────────┐
│                     PAGE STATE MACHINE                          │
└─────────────────────────────────────────────────────────────────┘

    Initial            Loading            Loaded
  (首次加载)    →      (数据加载)    →     (有数据)
       │                                      │
       │                                      │
       │                                      ↓
       │                                  Empty
       │                              (无数据)
       │                                      │
       │                                      │
       └──────────────────────────────────────┤
                                              │
                                              ↓
                                          Error
                                       (加载失败)

    States:
    ┌─────────────────────────────────────────────────────────┐
    │  Initial    │ 首次进入页面，尚未开始加载                   │
    │  Loading    │ 正在加载数据（显示骨架屏）                   │
    │  Loaded     │ 数据加载成功，有数据可显示                   │
    │  Empty      │ 数据加载成功，但无数据                       │
    │  Error      │ 数据加载失败（网络错误/服务器错误）            │
    └─────────────────────────────────────────────────────────┘
```

### 1.1 State Implementation

```vue
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

// ========================================
// STATE
// ========================================
const isLoading = ref(true)
const isInitialLoad = ref(true)  // Distinguish first load vs refresh
const error = ref<Error | null>(null)
const data = ref<Player[]>([])

// ========================================
// COMPUTED
// ========================================
const isEmpty = computed(() => {
  return !isLoading.value && !error.value && data.value.length === 0
})

const showContent = computed(() => {
  return !isLoading.value && !error.value && data.value.length > 0
})

// ========================================
// ACTIONS
// ========================================
async function fetchData() {
  isLoading.value = true
  error.value = null

  try {
    const response = await api.get('/players')
    data.value = response.data
  } catch (err) {
    error.value = err as Error
  } finally {
    isLoading.value = false
    isInitialLoad.value = false
  }
}

function retry() {
  fetchData()
}

// ========================================
// LIFECYCLE
// ========================================
onMounted(() => {
  fetchData()
})
</script>
```

### 1.2 Template Pattern

```vue
<template>
  <div class="page-container">
    <!-- Loading State -->
    <div v-if="isLoading && isInitialLoad" class="loading-state">
      <PageSkeleton />
    </div>

    <!-- Refresh Loading (with content overlay) -->
    <div v-else-if="isLoading && !isInitialLoad" class="refresh-loading">
      <LinearProgress />
    </div>

    <!-- Empty State -->
    <div v-else-if="isEmpty" class="empty-state">
      <EmptyState
        icon="🎮"
        title="暂无陪玩师"
        description="当前没有可用的陪玩师，请稍后再来"
      >
        <template #actions>
          <button @click="fetchData">刷新</button>
        </template>
      </EmptyState>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <ErrorState
        :error="error"
        @retry="retry"
      />
    </div>

    <!-- Content State -->
    <div v-else-if="showContent" class="content">
      <PlayerGrid :players="data" />
    </div>
  </div>
</template>
```

---

## 2. Component Structure

### 2.1 Standard Page Layout

```vue
<template>
  <div class="page-container">
    <!-- Page Header -->
    <header class="page-header">
      <div class="page-title">
        <h1>{{ pageTitle }}</h1>
        <p v-if="pageDescription" class="page-description">
          {{ pageDescription }}
        </p>
      </div>

      <!-- Page Actions -->
      <div v-if="$slots.actions" class="page-actions">
        <slot name="actions" />
      </div>
    </header>

    <!-- Page Content -->
    <main class="page-content">
      <slot />
    </main>

    <!-- Page Footer (optional) -->
    <footer v-if="$slots.footer" class="page-footer">
      <slot name="footer" />
    </footer>
  </div>
</template>

<script setup lang="ts">
interface Props {
  pageTitle: string
  pageDescription?: string
}

defineProps<Props>()
</script>

<style scoped>
.page-container {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
}

.page-title h1 {
  font-size: 24px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 8px 0;
}

.page-description {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin: 0;
}

.page-actions {
  display: flex;
  gap: 12px;
}

.page-content {
  min-height: 400px;
}
</style>
```

### 2.2 Page Component Template

```vue
<!-- pages/players/PlayerListPage.vue -->
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import PageLayout from '@/components/layout/PageLayout.vue'
import PlayerGrid from '@/components/player/PlayerGrid.vue'
import PlayerSkeleton from '@/components/player/PlayerSkeleton.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import ErrorState from '@/components/common/ErrorState.vue'

// State
const isLoading = ref(true)
const isInitialLoad = ref(true)
const error = ref<Error | null>(null)
const players = ref<Player[]>([])

// Computed
const isEmpty = computed(() => {
  return !isLoading.value && !error.value && players.value.length === 0
})

const showContent = computed(() => {
  return !isLoading.value && !error.value && players.value.length > 0
})

// Actions
async function fetchPlayers() {
  isLoading.value = true
  error.value = null

  try {
    const response = await api.get('/players', {
      params: route.query
    })
    players.value = response.data
  } catch (err) {
    error.value = err as Error
  } finally {
    isLoading.value = false
    isInitialLoad.value = false
  }
}

// Lifecycle
onMounted(() => {
  fetchPlayers()
})
</script>

<template>
  <PageLayout
    page-title="陪玩师列表"
    page-description="选择您心仪的陪玩师开始游戏"
  >
    <template #actions>
      <button variant="primary" @click="refresh">
        刷新
      </button>
    </template>

    <!-- Loading State -->
    <PlayerSkeleton v-if="isLoading && isInitialLoad" :count="12" />

    <!-- Empty State -->
    <EmptyState
      v-else-if="isEmpty"
      icon="🎮"
      title="暂无陪玩师"
      description="当前没有可用的陪玩师"
    />

    <!-- Error State -->
    <ErrorState
      v-else-if="error"
      :error="error"
      @retry="fetchPlayers"
    />

    <!-- Content -->
    <PlayerGrid v-else-if="showContent" :players="players" />
  </PageLayout>
</template>
```

---

## 3. Data Fetching Patterns

### 3.1 Single Data Fetch

```typescript
// composables/usePlayers.ts
import { ref } from 'vue'
import type { Player } from '@/types/player'

export function usePlayers() {
  const isLoading = ref(false)
  const error = ref<Error | null>(null)
  const players = ref<Player[]>([])

  async function fetchPlayers(params?: PlayerQueryParams) {
    isLoading.value = true
    error.value = null

    try {
      const response = await api.get('/players', { params })
      players.value = response.data
    } catch (err) {
      error.value = err as Error
      throw error
    } finally {
      isLoading.value = false
    }
  }

  return {
    isLoading,
    error,
    players,
    fetchPlayers
  }
}
```

### 3.2 Multiple Data Fetches (Parallel)

```typescript
// composables/useDashboard.ts
import { ref } from 'vue'

export function useDashboard() {
  const isLoading = ref(true)
  const error = ref<Error | null>(null)

  const stats = ref<DashboardStats | null>(null)
  const recentOrders = ref<Order[]>([])
  const upcomingSessions = ref<Session[]>([])

  async function fetchDashboardData() {
    isLoading.value = true
    error.value = null

    try {
      // Fetch all data in parallel
      const [statsRes, ordersRes, sessionsRes] = await Promise.all([
        api.get('/dashboard/stats'),
        api.get('/orders', { params: { limit: 5 } }),
        api.get('/sessions/upcoming')
      ])

      stats.value = statsRes.data
      recentOrders.value = ordersRes.data
      upcomingSessions.value = sessionsRes.data
    } catch (err) {
      error.value = err as Error
    } finally {
      isLoading.value = false
    }
  }

  return {
    isLoading,
    error,
    stats,
    recentOrders,
    upcomingSessions,
    fetchDashboardData
  }
}
```

### 3.3 Paginated Data Fetch

```typescript
// composables/usePaginatedData.ts
import { ref, computed } from 'vue'

export function usePaginatedData<T>(endpoint: string) {
  const isLoading = ref(false)
  const error = ref<Error | null>(null)

  const data = ref<T[]>([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(20)

  const totalPages = computed(() => Math.ceil(total.value / pageSize.value))
  const hasNextPage = computed(() => page.value < totalPages.value)
  const hasPrevPage = computed(() => page.value > 1)

  async function fetchData(params?: Record<string, any>) {
    isLoading.value = true
    error.value = null

    try {
      const response = await api.get(endpoint, {
        params: {
          page: page.value,
          pageSize: pageSize.value,
          ...params
        }
      })

      data.value = response.data.items
      total.value = response.data.total
    } catch (err) {
      error.value = err as Error
    } finally {
      isLoading.value = false
    }
  }

  function nextPage() {
    if (hasNextPage.value) {
      page.value++
      fetchData()
    }
  }

  function prevPage() {
    if (hasPrevPage.value) {
      page.value--
      fetchData()
    }
  }

  function goToPage(newPage: number) {
    if (newPage >= 1 && newPage <= totalPages.value) {
      page.value = newPage
      fetchData()
    }
  }

  function refresh() {
    page.value = 1
    fetchData()
  }

  return {
    isLoading,
    error,
    data,
    total,
    page,
    pageSize,
    totalPages,
    hasNextPage,
    hasPrevPage,
    fetchData,
    nextPage,
    prevPage,
    goToPage,
    refresh
  }
}
```

---

## 4. Loading States

### 4.1 Page Skeleton

```vue
<!-- components/common/PageSkeleton.vue -->
<template>
  <div class="page-skeleton">
    <div class="skeleton-header">
      <div class="skeleton-title"></div>
      <div class="skeleton-description"></div>
    </div>

    <div class="skeleton-content">
      <div
        v-for="i in count"
        :key="i"
        class="skeleton-card"
      ></div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  count?: number
}

withDefaults(defineProps<Props>(), {
  count: 8
})
</script>

<style scoped>
.page-skeleton {
  padding: 24px;
}

.skeleton-header {
  margin-bottom: 24px;
}

.skeleton-title {
  width: 200px;
  height: 32px;
  background: linear-gradient(
    90deg,
    var(--color-bg-secondary) 25%,
    var(--color-border-light) 50%,
    var(--color-bg-secondary) 75%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
}

.skeleton-description {
  width: 400px;
  height: 20px;
  margin-top: 8px;
  background: linear-gradient(
    90deg,
    var(--color-bg-secondary) 25%,
    var(--color-border-light) 50%,
    var(--color-bg-secondary) 75%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
}

.skeleton-content {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.skeleton-card {
  height: 280px;
  background: linear-gradient(
    90deg,
    var(--color-bg-card) 25%,
    var(--color-bg-secondary) 50%,
    var(--color-bg-card) 75%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 8px;
}

@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
</style>
```

### 4.2 Linear Progress (Refresh)

```vue
<!-- components/common/LinearProgress.vue -->
<template>
  <div class="linear-progress">
    <div class="linear-progress-bar"></div>
  </div>
</template>

<style scoped>
.linear-progress {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  z-index: 9999;
}

.linear-progress-bar {
  height: 100%;
  background: var(--color-primary-500);
  animation: progress 1.5s ease-in-out infinite;
}

@keyframes progress {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(100%);
  }
}
</style>
```

---

## 5. Empty States

```vue
<!-- components/common/EmptyState.vue -->
<template>
  <div class="empty-state">
    <div class="empty-state-icon">{{ icon }}</div>

    <h3 class="empty-state-title">{{ title }}</h3>

    <p v-if="description" class="empty-state-description">
      {{ description }}
    </p>

    <div v-if="$slots.actions" class="empty-state-actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  icon: string
  title: string
  description?: string
}

defineProps<Props>()
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 64px 24px;
  text-align: center;
}

.empty-state-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.empty-state-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 8px 0;
}

.empty-state-description {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin: 0 0 24px 0;
}

.empty-state-actions {
  display: flex;
  gap: 12px;
}
</style>
```

---

## 6. Error Handling

### 6.1 Error State Component

```vue
<!-- components/common/ErrorState.vue -->
<template>
  <div class="error-state">
    <div class="error-state-icon">⚠️</div>

    <h3 class="error-state-title">{{ title }}</h3>

    <p class="error-state-message">
      {{ message }}
    </p>

    <div v-if="showDetails" class="error-state-details">
      <code>{{ error?.message }}</code>
    </div>

    <button class="error-state-retry" @click="$emit('retry')">
      重试
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  error: Error | null
  showDetails?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showDetails: false
})

defineEmits<{
  retry: []
}>()

const title = computed(() => {
  if (!props.error) return '未知错误'

  const status = (props.error as any).response?.status

  switch (status) {
    case 401:
      return '登录已过期'
    case 403:
      return '无权访问'
    case 404:
      return '资源不存在'
    case 500:
      return '服务器错误'
    default:
      return '加载失败'
  }
})

const message = computed(() => {
  if (!props.error) return '请稍后再试'

  const status = (props.error as any).response?.status

  switch (status) {
    case 401:
      return '请重新登录'
    case 403:
      return '您没有权限访问此资源'
    case 404:
      return '请求的资源不存在'
    case 500:
      return '服务器内部错误，请稍后再试'
    default:
      return '网络连接异常，请检查网络后重试'
  }
})
</script>

<style scoped>
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 64px 24px;
  text-align: center;
}

.error-state-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.error-state-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 8px 0;
}

.error-state-message {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin: 0 0 24px 0;
}

.error-state-details {
  padding: 12px;
  background: var(--color-bg-deep);
  border-radius: 4px;
  margin-bottom: 24px;
}

.error-state-details code {
  font-size: 12px;
  color: var(--color-error);
  word-break: break-all;
}

.error-state-retry {
  padding: 10px 24px;
  background: var(--color-primary-500);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.error-state-retry:hover {
  background: var(--color-primary-600);
}
</style>
```

### 6.2 Global Error Handler

```typescript
// utils/errorHandler.ts
import { toast } from 'vue3-toastify'

export function handleError(error: any): string {
  console.error('Error:', error)

  const status = error.response?.status
  const message = error.response?.data?.message || error.message

  switch (status) {
    case 401:
      toast.error('登录已过期，请重新登录')
      // Redirect to login after delay
      setTimeout(() => {
        window.location.href = '/login'
      }, 1500)
      break

    case 403:
      toast.error('您没有权限执行此操作')
      break

    case 404:
      toast.error('请求的资源不存在')
      break

    case 429:
      toast.error('操作过于频繁，请稍后再试')
      break

    case 500:
      toast.error('服务器错误，请稍后再试')
      break

    default:
      if (error.message.includes('Network Error')) {
        toast.error('网络连接失败，请检查网络')
      } else {
        toast.error(message || '操作失败，请重试')
      }
  }

  return message
}
```

---

## 7. Permission-Based UI

### 7.1 Permission Store

```typescript
// stores/permission.ts
import { defineStore } from 'pinia'
import { computed } from 'vue'

export const usePermissionStore = defineStore('permission', () => {
  const permissions = ref<string[]>([])

  const hasPermission = computed(() => {
    return (permission: string) => {
      return permissions.value.includes(permission)
    }
  })

  const hasAnyPermission = computed(() => {
    return (perms: string[]) => {
      return perms.some(p => permissions.value.includes(p))
    }
  })

  const hasAllPermissions = computed(() => {
    return (perms: string[]) => {
      return perms.every(p => permissions.value.includes(p))
    }
  })

  function setPermissions(perms: string[]) {
    permissions.value = perms
  }

  return {
    permissions,
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    setPermissions
  }
})
```

### 7.2 Permission Button Component

```vue
<!-- components/common/PermissionButton.vue -->
<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { usePermissionStore } from '@/stores/permission'

interface Props {
  permission?: string
  permissions?: string[]
  role?: UserRole
  disabled?: boolean
  tooltip?: string
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false
})

const authStore = useAuthStore()
const permissionStore = usePermissionStore()

const isAllowed = computed(() => {
  if (props.permission) {
    return permissionStore.hasPermission(props.permission)
  }

  if (props.permissions) {
    return permissionStore.hasAnyPermission(props.permissions)
  }

  if (props.role) {
    return authStore.hasRole(props.role)
  }

  return true
})

const isDisabled = computed(() => {
  return props.disabled || !isAllowed.value
})

const showTooltip = computed(() => {
  return props.tooltip && !isAllowed.value
})
</script>

<template>
  <button
    v-bind="$attrs"
    :disabled="isDisabled"
    :title="showTooltip ? tooltip : undefined"
    class="permission-button"
    :class="{ disabled: isDisabled }"
  >
    <slot />
  </button>
</template>

<style scoped>
.permission-button {
  /* Base button styles */
}

.permission-button.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
```

---

## 8. Performance Optimization

### 8.1 Lazy Loading

```typescript
// Lazy load page components
const routes = [
  {
    path: '/players/:id',
    component: () => import('@/pages/players/PlayerDetailPage.vue')
  }
]
```

### 8.2 List Virtualization

```vue
<!-- Use for long lists -->
<script setup lang="ts">
import { useVirtualList } from '@vueuse/core'

const { list, containerProps, wrapperProps } = useVirtualList(
  longList,
  { itemHeight: 80 }
)
</script>

<template>
  <div v-bind="containerProps" style="height: 600px; overflow: auto;">
    <div v-bind="wrapperProps">
      <div
        v-for="{ data, index } in list"
        :key="index"
        style="height: 80px;"
      >
        {{ data }}
      </div>
    </div>
  </div>
</template>
```

### 8.3 Image Optimization

```vue
<template>
  <!-- Lazy load images -->
  <img
    v-lazy="imageUrl"
    :alt="imageAlt"
    loading="lazy"
  />

  <!-- Use WebP with fallback -->
  <picture>
    <source :srcset="webpUrl" type="image/webp" />
    <img :src="jpgUrl" :alt="imageAlt" loading="lazy" />
  </picture>
</template>
```

---

## 9. Accessibility Standards

### 9.1 Semantic HTML

```vue
<template>
  <article class="player-card">
    <header>
      <img :src="player.avatar" :alt="`头像 - ${player.nickname}`" />
    </header>

    <main>
      <h3>{{ player.nickname }}</h3>
      <p>{{ player.bio }}</p>
    </main>

    <footer>
      <button aria-label="查看详情">详情</button>
    </footer>
  </article>
</template>
```

### 9.2 Keyboard Navigation

```vue
<template>
  <div
    class="card"
    tabindex="0"
    @click="handleClick"
    @keydown.enter="handleClick"
    @keydown.space.prevent="handleClick"
  >
    <!-- Content -->
  </div>
</template>
```

### 9.3 Focus Management

```typescript
// Trap focus in modal
function trapFocus(element: HTMLElement) {
  const focusable = element.querySelectorAll(
    'a, button, input, textarea, select, [tabindex]:not([tabindex="-1"])'
  )

  const first = focusable[0] as HTMLElement
  const last = focusable[focusable.length - 1] as HTMLElement

  element.addEventListener('keydown', (e) => {
    if (e.key === 'Tab') {
      if (e.shiftKey && document.activeElement === first) {
        e.preventDefault()
        last.focus()
      } else if (!e.shiftKey && document.activeElement === last) {
        e.preventDefault()
        first.focus()
      }
    }
  })
}
```

---

## 10. Page Checklist

### Development Checklist

Before marking a page as complete, verify:

```
┌─────────────────────────────────────────────────────────────┐
│                    PAGE COMPLETION CHECKLIST                │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  STATE MANAGEMENT                                           │
│  ☐ Initial state handled                                   │
│  ☐ Loading state implemented (skeleton screen)             │
│  ☐ Empty state shown when no data                          │
│  ☐ Error state with retry button                           │
│  ☐ Content state displays data correctly                   │
│                                                             │
│  DATA FETCHING                                              │
│  ☐ API calls implemented                                   │
│  ☐ Error handling for API failures                         │
│  ☐ Loading indicators during fetch                         │
│  ☐ Pagination (if applicable)                              │
│  ☐ Refresh/retry functionality                             │
│                                                             │
│  ROUTING                                                    │
│  ☐ Route configured with meta (auth, roles)                │
│  ☐ Page title set                                          │
│  ☐ Redirect handling (after login/logout)                  │
│  ☐ Back button works correctly                             │
│                                                             │
│  PERMISSIONS                                                │
│  ☐ Permission checks implemented                           │
│  ☐ Protected UI elements (v-permission or <Protected>)     │
│  ☐ Fallback for unauthorized access                        │
│                                                             │
│  UI/UX                                                      │
│  ☐ Responsive design (mobile/desktop)                      │
│  ☐ Loading skeletons                                       │
│  ☐ Empty states with helpful messages                      │
│  ☐ Error states with clear next steps                      │
│  ☐ Hover states for desktop                                │
│  ☐ Focus states for keyboard navigation                    │
│                                                             │
│  ACCESSIBILITY                                              │
│  ☐ Semantic HTML (header, main, nav, etc.)                 │
│  ☐ Alt text for images                                     │
│  ☐ ARIA labels for interactive elements                    │
│  ☐ Keyboard navigation support                             │
│  ☐ Color contrast WCAG AA (4.5:1)                          │
│  ☐ Touch targets ≥ 44x44px (mobile)                        │
│                                                             │
│  PERFORMANCE                                                │
│  ☐ Lazy loading for components                             │
│  ☐ Image optimization (WebP, lazy load)                    │
│  ☐ No memory leaks (clean up effects)                      │
│  ☐ Efficient re-renders (use computed, keys)               │
│                                                             │
│  TESTING                                                    │
│  ☐ Unit tests for components                               │
│  ☐ API integration tests                                   │
│  ☐ Manual testing (desktop/mobile)                         │
│  ☐ Accessibility audit (axe DevTools)                      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

**Last Updated**: 2025-01-13
**Version**: 1.0.0
**Status**: Complete
