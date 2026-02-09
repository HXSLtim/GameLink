<template>
  <BasePageLayout
    class="review-list-page"
    :scroll="false"
    padding="0"
    title="我的评价"
    :show-back="true"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
  >
    <template #tabs>
      <!-- 标签切换 -->
      <TabsBar
        v-model="currentTab"
        :tabs="tabs"
        stretch
        @change="switchTab"
      />
    </template>

    <!-- 评价列表 -->
    <InfiniteList
      :state="pageState"
      :loading="loadingMore"
      :no-more="noMore"
      :empty-title="currentTab === 'pending' ? '暂无待评价订单' : '暂无评价记录'"
      :empty-desc="currentTab === 'pending' ? '完成订单后可以评价哦' : '快去下单体验吧'"
      padding="24rpx"
      @load-more="loadMore"
      @retry="refresh"
    >
      <view class="review-grid">
        <ListItem
          v-for="(review, index) in reviews"
          :key="review.id"
          :index="index"
        >
          <ReviewCard
            :review="review"
            :show-actions="currentTab === 'pending'"
            @order-click="goToOrder(review.orderId)"
            @skip="skipReview(review)"
            @write="writeReview(review)"
          />
        </ListItem>
      </view>
    </InfiniteList>

  </BasePageLayout>
</template>

<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
import TabsBar from '@/components/TabsBar/index.vue'
import InfiniteList from '@/components/InfiniteList/index.vue'
import ListItem from '@/components/ListItem/index.vue'
// Business 组件
import ReviewCard from '@/components/ReviewCard/index.vue'
// Composables
import { useReviewList } from '@/composables/useReviewList'

const {
  reviews,
  pageState,
  loadingMore,
  noMore,
  tabs,
  currentTab,
  loadMore,
  refresh,
  switchTab,
  goToOrder,
  skipReview,
  writeReview,
  goBack,
} = useReviewList()

onShow(() => {
  refresh()
})
</script>

<style lang="scss" scoped>
.review-grid {
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
