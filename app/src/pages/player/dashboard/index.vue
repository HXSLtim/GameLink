<template>
  <view class="dashboard-page page-container">
    <!-- 顶部导航 -->
    <NavBar title="工作台" :show-back="false">
      <template #right>
        <uv-icon name="setting-fill" size="22" color="var(--color-text)" @click="goToSettings"></uv-icon>
      </template>
    </NavBar>

    <scroll-view 
      class="content-scroll" 
      scroll-y 
      @refresherrefresh="onRefresh" 
      :refresher-enabled="true" 
      :refresher-triggered="refreshing"
    >
      <!-- 状态卡片 -->
      <DashboardStatusCard
        :status="workStatus"
        :nickname="playerInfo.nickname"
        :avatar="playerInfo.avatar"
        @toggle="toggleWorkStatus"
      />

      <!-- 今日数据 -->
      <StatsCard
        title="今日数据"
        :subtitle="formatDate(new Date())"
        :items="statItems"
      />

      <!-- 快捷入口 -->
      <QuickActions
        :items="quickActions"
        @click="handleQuickAction"
      />

      <!-- 功能列表 -->
      <MenuList
        :items="menuItems"
        @click="handleMenuClick"
      />

      <view class="bottom-placeholder"></view>
    </scroll-view>

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import MenuList from '@/components/MenuList/index.vue'
// Business 组件
import DashboardStatusCard from '@/components/DashboardStatusCard/index.vue'
import StatsCard from '@/components/StatsCard/index.vue'
import QuickActions from '@/components/QuickActions/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// Composables
import { usePlayerDashboard } from '@/composables/usePlayerDashboard'

const {
  refreshing,
  workStatus,
  playerInfo,
  statItems,
  quickActions,
  menuItems,
  onRefresh,
  toggleWorkStatus,
  handleQuickAction,
  handleMenuClick,
  goToSettings,
  formatDate,
  init,
} = usePlayerDashboard()

onMounted(init)

onShow(() => {
  init()
})
</script>

<style lang="scss" scoped>
.dashboard-page {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  box-sizing: border-box;
  
  @include desktop {
    height: 100vh;
    min-height: auto;
    overflow: hidden;
  }
}

.content-scroll {
  flex: 1;
  overflow-y: auto;
}

.bottom-placeholder {
  height: 100rpx;
}
</style>
