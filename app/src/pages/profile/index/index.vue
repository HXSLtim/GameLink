<template>
  <BasePageLayout
    class="profile-page"
    padding="0"
    title="我的"
    :show-back="false"
    :show-tab-bar="true"
    :tab-bar-current="3"
  >

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
      :counts="orderCounts"
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

  </BasePageLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import ProfileHeader from '@/components/ProfileHeader/index.vue'
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
import OrderQuickEntry from '@/components/OrderQuickEntry/index.vue'
import MenuList from '@/components/MenuList/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
// Business 组件
import ThemeToggle from '@/components/ThemeToggle/index.vue'
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
  padding-bottom: calc(110rpx + env(safe-area-inset-bottom));

  @include desktop {
    padding-bottom: 0;
  }
}

.action-section {
  padding: 48rpx 32rpx 80rpx;
}
</style>
