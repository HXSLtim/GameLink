<template>
  <view class="payment-result-page page-container">
    <!-- 结果卡片 -->
    <ResultCard
      :type="resultType"
      :title="resultTitle"
      :description="resultDesc"
      :amount="orderInfo.amount"
    />

    <!-- 订单信息 -->
    <GlCard v-if="orderInfo.orderNo" :shadow="false" bordered class="order-info-card">
      <view class="info-row">
        <text class="info-label">订单编号</text>
        <view class="info-value-wrap">
          <text class="info-value">{{ orderInfo.orderNo }}</text>
          <text class="copy-btn" @tap="copyOrderNo">复制</text>
        </view>
      </view>
      <view class="info-row">
        <text class="info-label">支付方式</text>
        <text class="info-value">{{ paymentMethodText }}</text>
      </view>
      <view class="info-row">
        <text class="info-label">支付时间</text>
        <text class="info-value">{{ formatTime(orderInfo.paidAt) }}</text>
      </view>
    </GlCard>

    <!-- 操作按钮 -->
    <view class="action-buttons">
      <GlButton type="primary" block round @click="goToOrderDetail">查看订单</GlButton>
      <GlButton type="default" plain block round @click="goToHome">返回首页</GlButton>
    </view>

    <!-- 温馨提示 -->
    <GlCard title="温馨提示" :shadow="false" bordered class="tips-card">
      <view class="tips-list">
        <text class="tips-item">• 支付成功后，陪玩师将尽快与您联系</text>
        <text class="tips-item">• 如有问题，请联系在线客服</text>
        <text class="tips-item">• 您可以在「我的订单」中查看订单状态</text>
      </view>
    </GlCard>

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
// Pattern 组件
import GlCard from '@/components/gl/Card/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
// Business 组件
import ResultCard from '@/components/ResultCard/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// Composables
import { usePaymentResult } from '@/composables/usePaymentResult'

const {
  resultType,
  orderInfo,
  resultTitle,
  resultDesc,
  paymentMethodText,
  formatTime,
  init,
  copyOrderNo,
  goToOrderDetail,
  goToHome,
} = usePaymentResult()

onLoad((options) => {
  init(options)
})
</script>

<style lang="scss" scoped>
.payment-result-page {
  min-height: 100vh;
  min-height: 100dvh;
  background: var(--color-bg);
  box-sizing: border-box;
  padding-bottom: 100rpx;
}

.order-info-card {
  margin: 0 24rpx 24rpx;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20rpx 0;
  border-bottom: 1rpx solid var(--color-border);
  
  &:last-child {
    border-bottom: none;
  }
}

.info-label {
  font-size: 28rpx;
  color: var(--color-text-secondary);
}

.info-value-wrap {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.info-value {
  font-size: 28rpx;
  color: var(--color-text);
}

.copy-btn {
  font-size: 24rpx;
  color: var(--color-primary);
}

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
  padding: 24rpx;
}

.tips-card {
  margin: 0 24rpx;
}

.tips-list {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.tips-item {
  font-size: 26rpx;
  color: var(--color-text-secondary);
  line-height: 1.5;
}
</style>
