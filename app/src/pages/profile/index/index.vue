<template>
  <view class="profile-page page-container">
    <!-- 用户信息卡片 -->
    <ProfileHeader
      :avatar="userInfo.avatar"
      :nickname="userInfo.nickname"
      :user-id="userInfo.id"
      :is-logged-in="isLoggedIn"
      :is-player="userInfo.role === 'player'"
      :order-count="userInfo.orderCount"
      :favorite-count="userInfo.favoriteCount"
      :balance="userInfo.balance"
      @edit="goToEdit"
      @login="goToLogin"
      @stat-click="handleStatClick"
    />
    
    <!-- 订单入口 - 仅登录用户显示 -->
    <OrderQuickEntry
      v-if="isLoggedIn"
      :pending-count="orderCounts.pending"
      :in-progress-count="orderCounts.inProgress"
      :to-review-count="orderCounts.toReview"
      @click="handleOrderClick"
      @view-all="handleViewAllOrders"
    />
    
    <!-- 功能菜单 - 登录用户专属 -->
    <MenuList
      v-if="isLoggedIn"
      :items="userMenuItems"
      @click="handleMenuClick"
    />
    
    <!-- 设置菜单 -->
    <MenuList
      :items="settingsMenuItems"
      @click="handleMenuClick"
    />
    
    <!-- 主题切换 -->
    <MenuList :items="[themeMenuItem]">
      <template #theme>
        <ThemeToggle />
      </template>
    </MenuList>
    
    <!-- 登录/退出按钮 -->
    <view class="action-section">
      <GlButton
        v-if="isLoggedIn"
        type="error"
        plain
        block
        round
        size="large"
        @click="handleLogout"
      >
        退出登录
      </GlButton>
      <GlButton
        v-else
        type="primary"
        block
        round
        size="large"
        @click="goToLogin"
      >
        立即登录
      </GlButton>
    </view>
    
    <!-- 自定义 TabBar -->
    <CustomTabBar :current="3" />
  </view>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import ProfileHeader from '@/components/ProfileHeader/index.vue'
import OrderQuickEntry from '@/components/OrderQuickEntry/index.vue'
import MenuList from '@/components/MenuList/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
// Business 组件
import ThemeToggle from '@/components/ThemeToggle/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// Composables
import { useProfile } from '@/composables/useProfile'

const {
  isLoggedIn,
  userInfo,
  orderCounts,
  userMenuItems,
  settingsMenuItems,
  themeMenuItem,
  loadUserData,
  handleMenuClick,
  handleOrderClick,
  handleViewAllOrders,
  handleStatClick,
  goToEdit,
  goToLogin,
  handleLogout,
} = useProfile()

onMounted(() => {
  loadUserData()
})

onShow(() => {
  loadUserData()
})
</script>

<style lang="scss" scoped>
.profile-page {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  box-sizing: border-box;
  padding-bottom: calc(110rpx + env(safe-area-inset-bottom));
  
  @include desktop {
    height: 100vh;
    min-height: auto;
    padding-bottom: 0;
    overflow: hidden;
  }
}

.action-section {
  padding: 48rpx 32rpx 80rpx;
}
</style>
