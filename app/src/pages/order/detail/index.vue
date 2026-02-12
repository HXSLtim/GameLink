<template>
  <BasePageLayout
    class="order-detail-page"
    :scroll="!loading"
    title="订单详情"
    :show-back="true"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
  >
    <template #nav>
      <!-- 顶部导航 -->
      <NavBar title="订单详情" @back="goBack" />
    </template>

    <!-- 加载状态 -->
    <view v-if="loading" class="loading-wrap">
      <Skeleton :rows="6" />
    </view>

    <!-- 内容区域 -->
    <template v-else>
      <!-- 订单状态卡片 -->
      <OrderStatusCard :status="order.status" :countdown="countdown" />

      <!-- 陪玩师信息 -->
      <PlayerCard
        :player="order.player"
        variant="compact"
        :show-order-count="true"
        :show-online-tag="false"
        :show-arrow="true"
        :clickable="true"
        @click="goToPlayer"
      />

      <!-- 服务信息 -->
      <OrderInfoSection title="服务信息" :items="serviceInfo" />

      <!-- 费用明细 -->
      <OrderFeeSection :fees="feeItems" :total="order.totalAmount || 0" />

      <!-- 订单信息 -->
      <OrderInfoSection title="订单信息" :items="orderInfo" />

      <!-- 评价区域 -->
      <OrderReviewSection
        v-if="order.status === 'completed' && order.review"
        :rating="order.review.rating"
        :content="order.review.content || ''"
      />

      <!-- 底部占位 -->
      <view class="bottom-placeholder"></view>
    </template>

    <template #footer>
      <!-- 底部操作栏 -->
      <OrderActionBar :status="order.status" :has-review="!!order.review" @action="handleAction" />

      <!-- 评价弹窗 -->
      <ReviewModal
        :show="showReviewModal"
        v-model:rating="reviewForm.rating"
        v-model:content="reviewForm.content"
        v-model:tags="reviewForm.tags"
        :available-tags="reviewTags"
        :loading="reviewLoading"
        @close="showReviewModal = false"
        @submit="submitReview"
      />
    </template>
  </BasePageLayout>
</template>

<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
// Business 组件
import OrderStatusCard from '@/components/OrderStatusCard/index.vue'
import PlayerCard from '@/components/PlayerCard/index.vue'
import OrderInfoSection from '@/components/OrderInfoSection/index.vue'
import OrderFeeSection from '@/components/OrderFeeSection/index.vue'
import OrderActionBar from '@/components/OrderActionBar/index.vue'
import OrderReviewSection from '@/components/OrderReviewSection/index.vue'
import ReviewModal from '@/components/ReviewModal/index.vue'
import Skeleton from '@/components/Skeleton/index.vue'
// Composables
import { useOrderDetail } from '@/composables/useOrderDetail'

const {
  loading,
  showReviewModal,
  reviewLoading,
  countdown,
  order,
  serviceInfo,
  orderInfo,
  feeItems,
  reviewForm,
  reviewTags,
  submitReview,
  loadOrderDetail,
  handleAction,
  goBack,
  goToPlayer,
} = useOrderDetail()

onLoad((options) => {
  const id = Number(options?.id)
  if (id) {
    loadOrderDetail(id)
  }
})
</script>

<style lang="scss" scoped>
.loading-wrap {
  flex: 1;
  padding: 24rpx;
}


.bottom-placeholder {
  height: 160rpx;
}
</style>
