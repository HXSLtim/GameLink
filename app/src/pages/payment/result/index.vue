<template>
  <BasePageLayout
    class="payment-result-page"
    padding="0"
    title="支付结果"
    :show-back="true"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
  >
    <!-- 结果卡片 -->
    <ResultCard
      :type="resultType"
      :title="resultTitle"
      :description="resultDesc"
      :amount="orderInfo.amount"
    />

    <!-- 订单信息 -->
    <SectionCard v-if="orderInfo.orderNo" margin="0 24rpx 24rpx" padding="24rpx">
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
    </SectionCard>

    <!-- 操作按钮 -->
    <view class="action-buttons">
      <GlButton type="primary" block round @click="goToOrderDetail">查看订单</GlButton>
      <GlButton type="default" plain block round @click="goToHome">返回首页</GlButton>
    </view>

    <!-- 温馨提示 -->
    <SectionCard title="温馨提示" margin="0 24rpx">
      <view class="tips-list">
        <text class="tips-item">• 支付成功后，陪玩师将尽快与您联系</text>
        <text class="tips-item">• 如有问题，请联系在线客服</text>
        <text class="tips-item">• 您可以在「我的订单」中查看订单状态</text>
      </view>
    </SectionCard>

  </BasePageLayout>
</template>

<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
// Pattern 组件
import SectionCard from '@/components/SectionCard/index.vue'
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
// Business 组件
import ResultCard from '@/components/ResultCard/index.vue'
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
  padding-bottom: 100rpx;
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
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}

.info-value-wrap {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.info-value {
  font-size: var(--font-sm);
  color: var(--color-text);
}

.copy-btn {
  font-size: var(--font-xs);
  color: var(--color-primary);
}

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
  padding: 24rpx;
}


.tips-list {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.tips-item {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  line-height: 1.5;
}
</style>
