<template>
  <view class="player-orders-page page-container">
    <!-- 顶部导航 -->
    <NavBar title="订单管理" @back="goBack" />

    <!-- 状态标签 -->
    <TabsBar
      v-model="currentTab"
      :tabs="tabs"
      scrollable
      @change="switchTab"
    />

    <!-- 订单列表 -->
    <InfiniteList
      :state="pageState"
      :loading="loadingMore"
      :no-more="noMore"
      :empty-title="getEmptyTitle()"
      :empty-desc="getEmptyDesc()"
      padding="24rpx"
      @load-more="loadMore"
      @retry="refresh"
    >
      <ListItem
        v-for="(order, index) in orders"
        :key="order.id"
        :index="index"
      >
        <PlayerOrderCard
          :order="order"
          @user-click="() => {}"
          @action="(action) => handleAction(order, action)"
        />
      </ListItem>
    </InfiniteList>

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import TabsBar from '@/components/TabsBar/index.vue'
import InfiniteList from '@/components/InfiniteList/index.vue'
import ListItem from '@/components/ListItem/index.vue'
// Business 组件
import PlayerOrderCard from '@/components/PlayerOrderCard/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// Composables
import { usePlayerOrders } from '@/composables/usePlayerOrders'

const {
  orders,
  pageState,
  loadingMore,
  noMore,
  tabs,
  currentTab,
  loadMore,
  refresh,
  switchTab,
  getEmptyTitle,
  getEmptyDesc,
  handleAction,
  goBack,
} = usePlayerOrders()

onShow(() => {
  refresh()
})
</script>

<style lang="scss" scoped>
.player-orders-page {
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
</style>
