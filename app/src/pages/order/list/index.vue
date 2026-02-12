<template>
  <BasePageLayout
    class="order-list-page"
    :scroll="false"
    padding="0"
    title="我的订单"
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
      :error-message="errorMessage"
      empty-title="暂无订单"
      empty-desc="快去找个陪玩师一起玩吧"
      padding="24rpx"
      @load-more="loadMore"
      @retry="refresh"
    >
      <template #default>
        <view class="order-grid">
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
        </view>
      </template>
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

const handleTabChange = (key: string, _tab: TabItem) => {
  switchTab(key as OrderTabKey)
}
</script>

<style lang="scss" scoped>
.order-list-page {
  padding-bottom: calc(110rpx + env(safe-area-inset-bottom));

  @include desktop {
    padding-bottom: 0;
  }
}

.order-grid {
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
