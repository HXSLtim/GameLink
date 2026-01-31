<template>
  <view class="order-detail-page page-container">
    <!-- 顶部导航 -->
    <NavBar title="订单详情" @back="goBack" />

    <!-- 加载状态 -->
    <view v-if="loading" class="loading-wrap">
      <uv-skeleton rows="6" title loading></uv-skeleton>
    </view>

    <!-- 内容区域 -->
    <scroll-view v-else class="content-scroll" scroll-y>
      <!-- 订单状态卡片 -->
      <OrderStatusCard :status="order.status" :countdown="countdown" />

      <!-- 陪玩师信息 -->
      <GlCard :shadow="false" bordered class="player-card" @tap="goToPlayer">
        <view class="player-row">
          <GlAvatar :src="order.player?.avatar" :text="order.player?.nickname" size="large" />
          <view class="player-info">
            <text class="player-name">{{ order.player?.nickname }}</text>
            <view class="player-meta">
              <text>⭐ {{ order.player?.rating?.toFixed(1) || '5.0' }}</text>
              <text>接单 {{ order.player?.orderCount || 0 }}</text>
            </view>
          </view>
          <uv-icon name="arrow-right" size="20" color="var(--color-text-secondary)"></uv-icon>
        </view>
      </GlCard>

      <!-- 服务信息 -->
      <OrderInfoSection title="服务信息" :items="serviceInfo" />

      <!-- 费用明细 -->
      <OrderFeeSection :fees="feeItems" :total="order.totalAmount || 0" />

      <!-- 订单信息 -->
      <OrderInfoSection title="订单信息" :items="orderInfo" />

      <!-- 评价区域 -->
      <GlCard v-if="order.status === 'completed' && order.review" title="我的评价" :shadow="false" bordered class="review-card">
        <RatingStars :rating="order.review.rating" size="small" />
        <text class="review-text">{{ order.review.content }}</text>
      </GlCard>

      <!-- 底部占位 -->
      <view class="bottom-placeholder"></view>
    </scroll-view>

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

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import GlCard from '@/components/gl/Card/index.vue'
import GlAvatar from '@/components/gl/Avatar/index.vue'
import RatingStars from '@/components/RatingStars/index.vue'
// Business 组件
import OrderStatusCard from '@/components/OrderStatusCard/index.vue'
import OrderInfoSection from '@/components/OrderInfoSection/index.vue'
import OrderFeeSection from '@/components/OrderFeeSection/index.vue'
import OrderActionBar from '@/components/OrderActionBar/index.vue'
import ReviewModal from '@/components/ReviewModal/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
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
.order-detail-page {
  height: 100vh;
  height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  box-sizing: border-box;
  overflow: hidden;
}

.loading-wrap {
  flex: 1;
  padding: 24rpx;
}

.content-scroll {
  flex: 1;
  padding: 24rpx;
  overflow-y: auto;
}

.player-card {
  cursor: pointer;
}

.player-row {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.player-info {
  flex: 1;
  min-width: 0;
}

.player-name {
  display: block;
  font-size: 30rpx;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 8rpx;
}

.player-meta {
  display: flex;
  gap: 16rpx;
  font-size: 24rpx;
  color: var(--color-text-secondary);
}

.review-card {
  margin-top: 24rpx;
}

.review-text {
  display: block;
  font-size: 28rpx;
  color: var(--color-text-secondary);
  margin-top: 12rpx;
  line-height: 1.5;
}

.bottom-placeholder {
  height: 160rpx;
}
</style>
