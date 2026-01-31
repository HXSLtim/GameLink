<template>
  <view class="review-list-page page-container">
    <!-- 顶部导航 -->
    <NavBar title="我的评价" @back="goBack" />

    <!-- 标签切换 -->
    <TabsBar
      v-model="currentTab"
      :tabs="tabItems"
      stretch
      @change="switchTab"
    />

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
    </InfiniteList>

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import TabsBar from '@/components/TabsBar/index.vue'
import InfiniteList from '@/components/InfiniteList/index.vue'
import ListItem from '@/components/ListItem/index.vue'
// Business 组件
import ReviewCard from '@/components/ReviewCard/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
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

// 转换为 TabsBar 格式
const tabItems = computed(() => 
  tabs.map(t => ({ key: t.value, label: t.label, badge: t.count || undefined }))
)

onShow(() => {
  refresh()
})
</script>

<style lang="scss" scoped>
.review-list-page {
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
