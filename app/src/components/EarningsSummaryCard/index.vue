<template>
  <view class="summary-card">
    <view class="summary-header">
      <text class="summary-label">可提现金额（元）</text>
      <GlButton type="primary" size="mini" round @click="$emit('withdraw')">提现</GlButton>
    </view>
    <view class="summary-balance">
      <text class="balance-value">{{ formatMoney(withdrawable) }}</text>
    </view>
    <view class="summary-row">
      <view class="summary-item">
        <text class="item-label">待结算</text>
        <text class="item-value">¥{{ formatMoney(pending) }}</text>
      </view>
      <view class="summary-item">
        <text class="item-label">已提现</text>
        <text class="item-value">¥{{ formatMoney(withdrawn) }}</text>
      </view>
      <view class="summary-item">
        <text class="item-label">累计收益</text>
        <text class="item-value highlight">¥{{ formatMoney(total) }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import GlButton from '@/components/gl/Button/index.vue'

interface Props {
  withdrawable: number // 分
  pending: number // 分
  withdrawn: number // 分
  total: number // 分
}

defineProps<Props>()

defineEmits<{
  withdraw: []
}>()

const formatMoney = (cents: number) => (cents / 100).toFixed(2)
</script>

<style lang="scss" scoped>
.summary-card {
  margin: 24rpx;
  padding: 32rpx;
  background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-primary-light, #4ADE80) 100%);
  border-radius: 24rpx;
  color: #FFFFFF;
}

.summary-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16rpx;
}

.summary-label {
  font-size: 26rpx;
  opacity: 0.85;
}

.summary-balance {
  margin-bottom: 32rpx;
}

.balance-value {
  font-size: 64rpx;
  font-weight: 800;
}

.summary-row {
  display: flex;
  padding-top: 24rpx;
  border-top: 1rpx solid rgba(255, 255, 255, 0.2);
}

.summary-item {
  flex: 1;
  text-align: center;
}

.item-label {
  display: block;
  font-size: 24rpx;
  opacity: 0.7;
  margin-bottom: 8rpx;
}

.item-value {
  font-size: 28rpx;
  font-weight: 600;
  
  &.highlight {
    font-weight: 700;
  }
}
</style>
