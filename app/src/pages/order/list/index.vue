<template>
  <view class="order-list-page page-container">
    <!-- 顶部导航 -->
    <NavBar title="我的订单" @back="goBack" />

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
      :error-message="errorMessage"
      empty-title="暂无订单"
      empty-desc="快去找个陪玩师一起玩吧"
      padding="24rpx"
      @load-more="loadMore"
      @retry="refresh"
    >
      <template #default>
        <ListItem
          v-for="(order, index) in orders"
          :key="order.id"
          :index="index"
          @click="goToDetail(order.id)"
        >
          <OrderCard
            :order="order"
            @action="handleOrderAction"
          />
        </ListItem>
      </template>
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
import OrderCard from '@/components/OrderCard/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// Composables
import { useOrderList } from '@/composables/useOrderList'

// 使用订单列表 Hook
const {
  orders,
  pageState,
  errorMessage,
  loadingMore,
  noMore,
  currentTab,
  tabs,
  loadMore,
  refresh,
  switchTab,
  handleOrderAction,
  goToDetail,
  goBack,
} = useOrderList()

onShow(() => {
  refresh()
})
</script>

<style lang="scss" scoped>
.order-list-page {
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
</style>
