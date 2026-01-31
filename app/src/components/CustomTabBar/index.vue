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
            :color="currentIndex === index ? '#fff' : 'var(--color-text-secondary)'"
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
            <uv-icon :name="isDark ? 'eye-off' : 'eye'" size="16" :color="isDark ? '#FFD700' : 'var(--color-text-secondary)'"></uv-icon>
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

// TabBar 项配置
interface TabItem {
  path: string
  text: string
  icon: string
  iconNormal?: string
  iconActive?: string
  badge?: number
  dot?: boolean
}

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
const tabItems = computed<TabItem[]>(() => [
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
const navItems = computed<TabItem[]>(() => [
  { path: '/pages/index/index', text: '首页', icon: 'home' },
  { path: '/pages/player/list/index', text: '陪玩', icon: 'grid' },
  { path: '/pages/message/list/index', text: '消息', icon: 'chat', badge: unreadCount.value || undefined },
  { path: '/pages/profile/index/index', text: '我的', icon: 'account' },
])

// PC 端额外导航项
const extraNavItems: TabItem[] = [
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
  z-index: 9998;
  display: flex;
  align-items: center;
  height: 110rpx;
  background: var(--color-bg-card, #FFFFFF);
  border-top: 1rpx solid var(--color-border, #E5E5E5);
  padding-bottom: env(safe-area-inset-bottom);
  transition: background-color 0.3s, border-color 0.3s;
  
  // PC 端隐藏移动端 TabBar
  @media screen and (min-width: 1024px) {
    display: none !important;
  }
  
  &.theme-dark {
    background: var(--color-bg-card, #1A1A2E);
    border-top-color: var(--color-border, #2D2D4A);
    
    .tabbar-text {
      color: var(--color-text-secondary, #94A3B8);
    }
    
    .tabbar-item.active .tabbar-text {
      color: var(--color-primary, #8B5CF6);
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
  transition: transform 0.2s;
  
  &:active {
    transform: scale(0.92);
  }
  
  &.active {
    .tabbar-text {
      color: var(--color-primary, #00D26A);
      font-weight: 600;
    }
    
    .tabbar-icon {
      transform: scale(1.1);
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
  transition: transform 0.2s cubic-bezier(0.34, 1.56, 0.64, 1);
  
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
  background: linear-gradient(135deg, #EF4444 0%, #DC2626 100%);
  border-radius: 16rpx;
  font-size: 20rpx;
  color: #FFFFFF;
  text-align: center;
  line-height: 32rpx;
  box-shadow: 0 2rpx 8rpx rgba(239, 68, 68, 0.4);
}

.tabbar-dot {
  position: absolute;
  top: -4rpx;
  right: -4rpx;
  width: 16rpx;
  height: 16rpx;
  background: #EF4444;
  border-radius: 50%;
  box-shadow: 0 0 8rpx rgba(239, 68, 68, 0.6);
}

.tabbar-text {
  font-size: 22rpx;
  color: var(--color-text-secondary, #666666);
  margin-top: 6rpx;
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
  z-index: 9999; // 确保在最上层
  transition: background-color 0.3s, border-color 0.3s;
  
  &.theme-dark {
    background: var(--color-bg-card, #1A1A2E);
    border-color: var(--color-border, #2D2D4A);
    
    .nav-item:hover {
      background: var(--color-bg-hover, #252542);
    }
    
    .nav-item.active {
      background: linear-gradient(135deg, var(--color-primary, #8B5CF6) 0%, var(--color-primary-dark, #7C3AED) 100%);
    }
    
    .logo-icon {
      background: linear-gradient(135deg, var(--color-primary, #8B5CF6) 0%, var(--color-primary-dark, #7C3AED) 100%);
    }
  }
}

.sidebar-logo {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px;
  cursor: pointer;
  transition: opacity 0.2s;
  
  &:hover {
    opacity: 0.8;
  }
  
  .logo-icon {
    width: 40px;
    height: 40px;
    background: linear-gradient(135deg, var(--color-primary, #00D26A) 0%, var(--color-primary-dark, #00B359) 100%);
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    
    text {
      font-size: 20px;
      font-weight: 700;
      color: #fff;
    }
  }
  
  .logo-text {
    font-size: 18px;
    font-weight: 700;
    color: var(--color-text);
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
  margin-bottom: 4px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
  
  &:hover {
    background: var(--color-bg-secondary, #F1F5F9);
  }
  
  &.active {
    background: linear-gradient(135deg, var(--color-primary, #00D26A) 0%, var(--color-primary-dark, #00B359) 100%);
    
    .nav-label {
      color: #fff;
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
  background: #EF4444;
  border-radius: 9px;
  font-size: 11px;
  font-weight: 600;
  color: #fff;
  text-align: center;
  line-height: 18px;
}

.nav-label {
  font-size: 14px;
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
  border-radius: 10px;
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
  background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-primary-dark, #00B85A) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  flex-shrink: 0;
  
  text {
    color: #fff;
    font-size: 14px;
    font-weight: 600;
  }
  
  image {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  text {
    font-size: 16px;
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
    background: #10B981;
    border: 2px solid var(--color-bg-card);
    border-radius: 50%;
  }
}

.user-info {
  flex: 1;
  min-width: 0;
  
  .user-name {
    display: block;
    font-size: 14px;
    font-weight: 600;
    color: var(--color-text);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  
  .user-status {
    display: block;
    font-size: 12px;
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
</style>
