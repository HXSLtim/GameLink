<template>
  <!-- PC 端：侧边栏 -->
  <view v-if="isPC" class="sidebar" :class="{ 'theme-dark': isDark }">
    <!-- Logo -->
    <view class="sidebar-logo" @click="switchTo(0)">
      <view class="logo-icon">
        <text>G</text>
      </view>
      <text class="logo-text">GameLink</text>
    </view>
    
    <!-- 主导航 -->
    <view class="sidebar-nav">
      <view 
        v-for="(item, index) in navItems" 
        :key="item.path"
        class="nav-item"
        :class="{ active: currentIndex === index }"
        @click="switchTo(index)"
      >
        <view class="nav-icon">
          <uv-icon 
            :name="item.icon" 
            :size="20" 
            :color="currentIndex === index ? 'var(--color-icon-accent)' : 'var(--color-text-secondary)'"
          ></uv-icon>
          <view v-if="item.badge" class="nav-badge">{{ item.badge > 99 ? '99+' : item.badge }}</view>
        </view>
        <text class="nav-label">{{ item.text }}</text>
      </view>
      
      <!-- 分割线 -->
      <view class="nav-divider"></view>
      
      <!-- 额外导航 -->
      <view 
        v-for="item in extraNavItems" 
        :key="item.path"
        class="nav-item"
        @click="navigateTo(item.path)"
      >
        <view class="nav-icon">
          <uv-icon :name="item.icon" :size="20" color="var(--color-text-secondary)"></uv-icon>
        </view>
        <text class="nav-label">{{ item.text }}</text>
      </view>
    </view>
    
    <!-- 底部用户区 -->
    <view class="sidebar-footer">
      <view class="user-row">
        <view class="user-section" @click="goToProfile">
          <view class="user-avatar" :class="{ online: isLoggedIn }">
            <image v-if="userInfo?.avatar" :src="userInfo.avatar" mode="aspectFill" />
            <text v-else>{{ userInfo?.nickname?.[0] || '游' }}</text>
          </view>
          <view class="user-info">
            <text class="user-name">{{ userInfo?.nickname || '游客' }}</text>
            <text class="user-status">{{ isLoggedIn ? '在线' : '点击登录' }}</text>
          </view>
        </view>
        <view class="footer-actions">
          <view class="action-btn" @click="toggleTheme" :title="isDark ? '切换日间' : '切换夜间'">
            <image
              class="action-icon"
              :src="isDark ? '/static/icons/moon.svg' : '/static/icons/sun.svg'"
              mode="aspectFit"
            />
          </view>
          <view class="action-btn" @click="navigateTo('/pages/settings/index/index')" title="设置">
            <uv-icon name="setting-fill" size="16" color="var(--color-text-secondary)"></uv-icon>
          </view>
        </view>
      </view>
    </view>
  </view>

  <!-- 移动端：底部 TabBar（仅在 showMobileTabBar 为 true 时显示） -->
  <view v-else-if="showMobileTabBar" class="tabbar" :class="{ 'theme-dark': isDark }">
    <view
      v-for="(item, index) in tabItems"
      :key="item.path"
      class="tabbar-item"
      :class="{ active: currentIndex === index }"
      @tap="switchTo(index)"
    >
      <view class="tabbar-icon">
        <!-- 使用本地图标文件 -->
        <image 
          class="icon-img" 
          :src="currentIndex === index ? item.iconActive : item.iconNormal" 
          mode="aspectFit"
        />
        <view v-if="item.badge" class="tabbar-badge">{{ item.badge > 99 ? '99+' : item.badge }}</view>
        <view v-else-if="item.dot" class="tabbar-dot"></view>
      </view>
      <text class="tabbar-text">{{ item.text }}</text>
    </view>
  </view>
  
  <!-- 占位（仅在移动端且显示 TabBar 时） -->
  <view v-if="!isPC && showMobileTabBar" class="tabbar-placeholder"></view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useUserStore } from '@/store/user'
import { useAppStore } from '@/store/app'
import { useTheme } from '@/composables/useTheme'
import type { TabBarItem } from '@/types/ui'

const props = withDefaults(defineProps<{
  current?: number
  showMobileTabBar?: boolean  // 移动端是否显示底部 TabBar
}>(), {
  showMobileTabBar: true,
  current: 0,
})

const userStore = useUserStore()
const appStore = useAppStore()
const { isDark, toggleTheme } = useTheme()

// PC 端检测 - 立即检测以避免闪烁
const getInitialIsPC = () => {
  // #ifdef H5
  if (typeof window !== 'undefined') {
    return window.innerWidth >= 1024
  }
  // #endif
  return false
}

const isPC = ref(getInitialIsPC())

const checkDevice = () => {
  // #ifdef H5
  isPC.value = window.innerWidth >= 1024
  // #endif
}

// #ifdef H5
if (typeof window !== 'undefined') {
  window.addEventListener('resize', checkDevice)
}
// #endif

// 用户状态
const isLoggedIn = computed(() => userStore.isLoggedIn)
const userInfo = computed(() => userStore.userInfo)
const unreadCount = computed(() => appStore.unreadCount || 0)

// 当前选中索引
const currentIndex = ref(props.current)

// 监听 props
watch(() => props.current, (val) => {
  currentIndex.value = val
})

// TabBar 项（移动端底部）
const tabItems = computed<TabBarItem[]>(() => [
  { 
    path: '/pages/index/index', 
    text: '首页', 
    icon: 'home',
    iconNormal: '/static/icons/home.svg',
    iconActive: '/static/icons/home-active.svg'
  },
  { 
    path: '/pages/player/list/index', 
    text: '陪玩', 
    icon: 'grid',
    iconNormal: '/static/icons/player.svg',
    iconActive: '/static/icons/player-active.svg'
  },
  { 
    path: '/pages/message/list/index', 
    text: '消息', 
    icon: 'chat',
    iconNormal: '/static/icons/message.svg',
    iconActive: '/static/icons/message-active.svg',
    badge: unreadCount.value || undefined 
  },
  { 
    path: '/pages/profile/index/index', 
    text: '我的', 
    icon: 'account',
    iconNormal: '/static/icons/profile.svg',
    iconActive: '/static/icons/profile-active.svg'
  },
])

// PC 端导航项（侧边栏）
const navItems = computed<TabBarItem[]>(() => [
  { path: '/pages/index/index', text: '首页', icon: 'home' },
  { path: '/pages/player/list/index', text: '陪玩', icon: 'grid' },
  { path: '/pages/message/list/index', text: '消息', icon: 'chat', badge: unreadCount.value || undefined },
  { path: '/pages/profile/index/index', text: '我的', icon: 'account' },
])

// PC 端额外导航项
const extraNavItems: TabBarItem[] = [
  { path: '/pages/channel/list/index', text: '频道', icon: 'more-circle' },
  { path: '/pages/order/list/index', text: '订单', icon: 'file-text' },
  { path: '/pages/wallet/index/index', text: '钱包', icon: 'red-packet' },
]

// 根据当前路由更新索引
const updateCurrentIndex = () => {
  const pages = getCurrentPages()
  if (pages.length > 0) {
    const currentPage = pages[pages.length - 1]
    const route = '/' + currentPage.route
    // 检查 tabItems（移动端）和 navItems（PC端）
    const items = isPC.value ? navItems.value : tabItems.value
    const index = items.findIndex(item => item.path === route)
    if (index >= 0) {
      currentIndex.value = index
    }
  }
}

// 生命周期钩子 - 放在 updateCurrentIndex 之后
onMounted(() => {
  checkDevice()
  updateCurrentIndex()
})

// 页面显示时更新索引
onShow(() => {
  updateCurrentIndex()
})

// 切换 Tab
const switchTo = (index: number) => {
  if (currentIndex.value === index) return
  
  const item = tabItems.value[index]
  if (!item) return
  
  currentIndex.value = index
  
  uni.switchTab({
    url: item.path,
    fail: () => {
      uni.navigateTo({ url: item.path })
    }
  })
}

// 导航到非 TabBar 页面
const navigateTo = (path: string) => {
  uni.navigateTo({ url: path })
}

// 跳转到个人中心
const goToProfile = () => {
  if (isLoggedIn.value) {
    switchTo(3) // 我的 Tab
  } else {
    uni.navigateTo({ url: '/pages/auth/login/index' })
  }
}
</script>

<style lang="scss" scoped>
// ============================================
// 移动端底部 TabBar
// ============================================
.tabbar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  height: 100rpx;
  background: var(--color-bg-card);
  border-top: 1rpx solid var(--color-border, #E5E5E5);
  padding-bottom: env(safe-area-inset-bottom);
  transition: background-color 0.3s ease, border-color 0.3s ease;
  // 移动端轻微投影，强化层级
  box-shadow: 0 -1px 4px rgba(0, 0, 0, 0.04);
  
  // PC 端隐藏移动端 TabBar
  @media screen and (min-width: 1024px) {
    display: none !important;
  }
  
  &.theme-dark {
    background: var(--color-bg-card, #1C1D20);
    border-top-color: var(--color-border, #2A2B30);
    box-shadow: 0 -1px 4px rgba(0, 0, 0, 0.2);
    
    .tabbar-text {
      color: var(--color-text-secondary, #94A3B8);
    }
    
    .tabbar-item.active .tabbar-text {
      color: var(--color-icon-accent, #5865F2);
    }
  }
}

.tabbar-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  transition: background 0.2s;
  
  &.active {
    .tabbar-text {
      color: var(--color-icon-accent, #7ACC35);
      font-weight: 600;
    }
  }
}

.tabbar-icon {
  position: relative;
  width: 52rpx;
  height: 52rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  
  // 图片图标样式
  .icon-img {
    width: 48rpx;
    height: 48rpx;
    transition: transform 0.2s;
  }
}

.tabbar-badge {
  position: absolute;
  top: -8rpx;
  right: -20rpx;
  min-width: 32rpx;
  height: 32rpx;
  padding: 0 8rpx;
  background: var(--color-error);
  border-radius: var(--radius-full);
  font-size: var(--font-xs);
  color: #FFFFFF;
  text-align: center;
  line-height: 32rpx;
}

.tabbar-dot {
  position: absolute;
  top: -4rpx;
  right: -4rpx;
  width: 16rpx;
  height: 16rpx;
  background: var(--color-error);
  border-radius: var(--radius-full);
}

.tabbar-text {
  font-size: var(--font-xs);
  color: var(--color-text-secondary, #666666);
  margin-top: var(--spacing-xs);
  transition: color 0.3s, font-weight 0.2s;
}

.tabbar-placeholder {
  // 不使用占位符，由页面自行处理 padding-bottom
  display: none;
  height: 0;
}

// ============================================
// PC 端侧边栏
// ============================================
.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  width: 240px;
  display: flex;
  flex-direction: column;
  background: var(--color-bg-card, #FFFFFF);
  border-right: 1px solid var(--color-border, #E5E5E5);
  z-index: 100;
  transition: background-color 0.3s ease, border-color 0.3s ease;
  
  &.theme-dark {
    background: var(--color-bg-card, #1C1D20);
    border-color: var(--color-border, #2A2B30);
    
    .nav-item:hover {
      background: var(--color-bg-hover, #252542);
    }
    
    .nav-item.active {
      background: var(--color-bg-secondary);
    }
    
    .logo-icon {
      background: var(--color-bg-secondary);
    }
  }
}

.sidebar-logo {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  cursor: pointer;
  transition: opacity 0.2s ease;
  border-bottom: 1px solid var(--color-divider);
  
  &:hover {
    opacity: 0.8;
  }
  
  .logo-icon {
    width: 36px;
    height: 36px;
    background: var(--color-primary);
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    
    text {
      font-size: var(--font-base);
      font-weight: 700;
      color: #FFFFFF;
    }
  }
  
  .logo-text {
    font-size: var(--font-md);
    font-weight: 700;
    color: var(--color-text);
    letter-spacing: -0.3px;
  }
}

.sidebar-nav {
  flex: 1;
  padding: 8px 12px;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  margin-bottom: 2px;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.15s ease;
  position: relative;
  
  &:hover {
    background: var(--color-bg-secondary, #F1F5F9);
  }
  
  &.active {
    background: var(--color-bg-secondary);
    
    // 左侧主题色指示条
    &::before {
      content: '';
      position: absolute;
      left: 0;
      top: 8px;
      bottom: 8px;
      width: 3px;
      border-radius: 0 3px 3px 0;
      background: var(--color-icon-accent);
    }
    
    .nav-label {
      color: var(--color-text);
      font-weight: 600;
    }
  }
}

.nav-icon {
  position: relative;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-badge {
  position: absolute;
  top: -6px;
  right: -10px;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  background: var(--color-error);
  border-radius: 9px;
  font-size: var(--font-xs);
  font-weight: 600;
  color: #fff;
  text-align: center;
  line-height: 18px;
}

.nav-label {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  transition: color 0.2s;
}

.nav-divider {
  height: 1px;
  margin: 12px 16px;
  background: var(--color-border);
}

.sidebar-footer {
  padding: 12px;
  border-top: 1px solid var(--color-border);
}

.user-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.user-section {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s;
  flex: 1;
  min-width: 0;
  
  &:hover {
    background: var(--color-bg-secondary);
  }
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  flex-shrink: 0;
  
  image {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  text {
    font-size: var(--font-sm);
    font-weight: 600;
    color: var(--color-text-secondary);
  }
  
  &.online::after {
    content: '';
    position: absolute;
    bottom: 2px;
    right: 2px;
    width: 10px;
    height: 10px;
    background: var(--color-success);
    border: 2px solid var(--color-bg-card);
    border-radius: 50%;
  }
}

.user-info {
  flex: 1;
  min-width: 0;
  
  .user-name {
    display: block;
    font-size: var(--font-sm);
    font-weight: 600;
    color: var(--color-text);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  
  .user-status {
    display: block;
    font-size: var(--font-xs);
    color: var(--color-text-secondary);
  }
}

.footer-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.action-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  
  &:hover {
    background: var(--color-bg-secondary);
    transform: scale(1.05);
  }
}

.action-icon {
  width: 18px;
  height: 18px;
}
</style>
