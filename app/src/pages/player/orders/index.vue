<template>
  <BasePageLayout
    class="player-orders-page"
    :scroll="false"
    padding="0"
    title="订单管理"
    :show-back="true"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
  >
    <template #tabs>
      <!-- 状态标签 -->
      <TabsBar
        v-model="currentTab"
        :tabs="tabs"
        scrollable
        @change="handleTabChange"
      />
    </template>

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
      <view class="player-orders-grid">
        <ListItem
          v-for="(order, index) in orders"
          :key="order.id"
          :index="index"
        >
          <OrderCard
            :order="order"
            view-mode="player"
            @person-click="() => {}"
            @action="(action) => handleAction(order, action)"
          />
        </ListItem>
      </view>
    </InfiniteList>

  </BasePageLayout>
</template>

<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import type { TabItem } from '@/types/ui'
import type { OrderTabKey } from '@/types/order'
// Pattern 组件
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
import TabsBar from '@/components/TabsBar/index.vue'
import InfiniteList from '@/components/InfiniteList/index.vue'
import ListItem from '@/components/ListItem/index.vue'
// Business 组件
import OrderCard from '@/components/OrderCard/index.vue'
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

const handleTabChange = (key: string, _tab: TabItem) => {
  switchTab(key as OrderTabKey)
}
</script>

<style lang="scss" scoped>
.player-orders-grid {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);

  @include desktop {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    column-gap: var(--spacing-sm);
    row-gap: var(--spacing-sm);
  }

  @include desktop-lg {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  :deep(.list-item) {
    margin-bottom: 0;
  }
}
</style>
