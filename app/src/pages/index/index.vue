<template>
  <view class="home-page" :class="{ 'theme-dark': isDark }">
    <!-- 顶部状态栏占位 -->
    <view class="status-bar safe-area-top"></view>
    
    <!-- 头部 -->
    <view class="header">
      <view class="header-left">
        <text class="app-name">GameLink</text>
      </view>
      <view class="header-right">
        <view class="theme-toggle" @click="toggleTheme">
          <text>{{ isDark ? '🌙' : '☀️' }}</text>
        </view>
      </view>
    </view>

    <!-- 未登录状态 -->
    <view v-if="!isLoggedIn" class="welcome-section">
      <view class="welcome-content">
        <image class="welcome-image" src="/static/logo.png" mode="aspectFit" />
        <text class="welcome-title">欢迎来到 GameLink</text>
        <text class="welcome-desc">专业游戏陪玩平台，让游戏更有趣</text>
        
        <view class="welcome-actions">
          <button class="btn-primary" @click="goToLogin">
            立即登录
          </button>
          <button class="btn-outline" @click="goToRegister">
            注册账号
          </button>
        </view>
      </view>
    </view>

    <!-- 已登录状态 -->
    <view v-else class="main-content">
      <!-- 用户信息卡片 -->
      <view class="user-card card">
        <view class="user-avatar">
          <image 
            v-if="userInfo?.avatar" 
            :src="userInfo.avatar" 
            mode="aspectFill" 
          />
          <text v-else class="avatar-placeholder">{{ userInfo?.nickname?.[0] || '用' }}</text>
        </view>
        <view class="user-info">
          <text class="user-name">{{ userInfo?.nickname || '用户' }}</text>
          <text class="user-role">{{ userInfo?.role === 'player' ? '陪玩师' : '用户' }}</text>
        </view>
        <view class="user-action" @click="handleLogout">
          <text class="logout-text">退出</text>
        </view>
      </view>

      <!-- 快捷入口 -->
      <view class="quick-actions">
        <text class="section-title">快捷入口</text>
        <view class="action-grid">
          <view class="action-item" @click="goTo('/pages/user/players/index')">
            <view class="action-icon bg-primary">找陪玩</view>
            <text class="action-name">找陪玩</text>
          </view>
          <view class="action-item" @click="goTo('/pages/user/order/list/index')">
            <view class="action-icon bg-info">我的订单</view>
            <text class="action-name">我的订单</text>
          </view>
          <view class="action-item" @click="goTo('/pages/user/wallet/index')">
            <view class="action-icon bg-warning">我的钱包</view>
            <text class="action-name">我的钱包</text>
          </view>
          <view class="action-item" @click="goTo('/pages/user/messages/index')">
            <view class="action-icon bg-error">消息</view>
            <text class="action-name">消息中心</text>
          </view>
        </view>
      </view>

      <!-- 陪玩师快捷入口 -->
      <view v-if="isPlayer" class="quick-actions">
        <text class="section-title">陪玩师工作台</text>
        <view class="action-grid">
          <view class="action-item" @click="goTo('/pages/player/dashboard/index')">
            <view class="action-icon bg-primary">工作台</view>
            <text class="action-name">工作台</text>
          </view>
          <view class="action-item" @click="goTo('/pages/player/orders/hall/index')">
            <view class="action-icon bg-success">接单大厅</view>
            <text class="action-name">接单大厅</text>
          </view>
          <view class="action-item" @click="goTo('/pages/player/earnings/index')">
            <view class="action-icon bg-warning">我的收益</view>
            <text class="action-name">我的收益</text>
          </view>
          <view class="action-item" @click="goTo('/pages/player/certification/index')">
            <view class="action-icon bg-info">认证中心</view>
            <text class="action-name">认证中心</text>
          </view>
        </view>
      </view>

      <!-- 提示信息 -->
      <view class="tip-card card">
        <text class="tip-title">温馨提示</text>
        <text class="tip-content">
          当前为开发预览版，更多功能正在开发中...
        </text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useUserStore } from '@/store/user'
import { useTheme } from '@/composables/useTheme'

const userStore = useUserStore()
const { isDark, toggleTheme } = useTheme()

// 计算属性
const isLoggedIn = computed(() => userStore.isLoggedIn)
const isPlayer = computed(() => userStore.isPlayer)
const userInfo = computed(() => userStore.userInfo)

// 跳转到登录
function goToLogin() {
  uni.navigateTo({ url: '/pages/auth/login/index' })
}

// 跳转到注册
function goToRegister() {
  uni.navigateTo({ url: '/pages/auth/register/index' })
}

// 通用跳转
function goTo(url: string) {
  uni.navigateTo({ url })
}

// 退出登录
function handleLogout() {
  uni.showModal({
    title: '提示',
    content: '确定要退出登录吗？',
    success: (res) => {
      if (res.confirm) {
        userStore.logout()
      }
    }
  })
}
</script>

<style lang="scss" scoped>
.home-page {
  min-height: 100vh;
  background: var(--color-bg);
}

.status-bar {
  height: var(--status-bar-height, 44px);
}

// 头部
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24rpx 32rpx;
  background: var(--color-bg-card);
  
  .app-name {
    font-size: 40rpx;
    font-weight: 700;
    color: var(--color-primary);
  }
  
  .theme-toggle {
    width: 64rpx;
    height: 64rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--color-bg-secondary);
    border-radius: 50%;
    font-size: 32rpx;
  }
}

// 欢迎页面
.welcome-section {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 64rpx 48rpx;
  min-height: calc(100vh - 200rpx);
}

.welcome-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  
  .welcome-image {
    width: 200rpx;
    height: 200rpx;
    margin-bottom: 48rpx;
  }
  
  .welcome-title {
    font-size: 44rpx;
    font-weight: 600;
    color: var(--color-text);
    margin-bottom: 16rpx;
  }
  
  .welcome-desc {
    font-size: 28rpx;
    color: var(--color-text-secondary);
    margin-bottom: 64rpx;
  }
}

.welcome-actions {
  display: flex;
  flex-direction: column;
  gap: 24rpx;
  width: 100%;
  max-width: 500rpx;
  
  .btn-primary {
    width: 100%;
    height: 96rpx;
    background: var(--color-primary);
    color: #FFFFFF;
    font-size: 32rpx;
    font-weight: 500;
    border-radius: 48rpx;
    border: none;
    
    &::after {
      border: none;
    }
  }
  
  .btn-outline {
    width: 100%;
    height: 96rpx;
    background: transparent;
    color: var(--color-primary);
    font-size: 32rpx;
    font-weight: 500;
    border-radius: 48rpx;
    border: 2rpx solid var(--color-primary);
    
    &::after {
      border: none;
    }
  }
}

// 主内容
.main-content {
  padding: 24rpx 32rpx;
}

// 卡片
.card {
  background: var(--color-bg-card);
  border-radius: 24rpx;
  padding: 32rpx;
  margin-bottom: 24rpx;
}

// 用户卡片
.user-card {
  display: flex;
  align-items: center;
  
  .user-avatar {
    width: 96rpx;
    height: 96rpx;
    border-radius: 50%;
    overflow: hidden;
    background: var(--color-primary);
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 24rpx;
    
    image {
      width: 100%;
      height: 100%;
    }
    
    .avatar-placeholder {
      font-size: 40rpx;
      color: #FFFFFF;
      font-weight: 600;
    }
  }
  
  .user-info {
    flex: 1;
    
    .user-name {
      display: block;
      font-size: 32rpx;
      font-weight: 600;
      color: var(--color-text);
      margin-bottom: 8rpx;
    }
    
    .user-role {
      font-size: 24rpx;
      color: var(--color-text-secondary);
    }
  }
  
  .logout-text {
    font-size: 28rpx;
    color: var(--color-error);
  }
}

// 快捷入口
.quick-actions {
  margin-bottom: 32rpx;
  
  .section-title {
    display: block;
    font-size: 28rpx;
    color: var(--color-text-secondary);
    margin-bottom: 24rpx;
    padding-left: 8rpx;
  }
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24rpx;
}

.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  
  .action-icon {
    width: 96rpx;
    height: 96rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 24rpx;
    font-size: 24rpx;
    color: #FFFFFF;
    font-weight: 500;
    margin-bottom: 12rpx;
    
    &.bg-primary { background: var(--color-primary); }
    &.bg-success { background: var(--color-success); }
    &.bg-warning { background: var(--color-warning); }
    &.bg-error { background: var(--color-error); }
    &.bg-info { background: var(--color-info); }
  }
  
  .action-name {
    font-size: 24rpx;
    color: var(--color-text);
  }
}

// 提示卡片
.tip-card {
  background: var(--color-bg-secondary);
  
  .tip-title {
    display: block;
    font-size: 28rpx;
    font-weight: 500;
    color: var(--color-text);
    margin-bottom: 12rpx;
  }
  
  .tip-content {
    font-size: 24rpx;
    color: var(--color-text-secondary);
    line-height: 1.6;
  }
}
</style>
