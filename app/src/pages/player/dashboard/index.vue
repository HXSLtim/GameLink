<template>
  <BasePageLayout
    class="dashboard-page"
    padding="0"
    title="工作台"
    :show-back="false"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
    :refresher-enabled="true"
    :refresher-triggered="refreshing"
    @refresherrefresh="onRefresh"
  >
    <template #nav>
      <!-- 顶部导航 -->
      <NavBar title="工作台" :show-back="false">
        <template #right>
          <uv-icon name="setting-fill" size="22" color="var(--color-text)" @click="goToSettings"></uv-icon>
        </template>
      </NavBar>
    </template>

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

  </BasePageLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
import MenuList from '@/components/MenuList/index.vue'
// Business 组件
import DashboardStatusCard from '@/components/DashboardStatusCard/index.vue'
import StatsCard from '@/components/StatsCard/index.vue'
import QuickActions from '@/components/QuickActions/index.vue'
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
.bottom-placeholder {
  height: 100rpx;
}
</style>
