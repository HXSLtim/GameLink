<template>
  <view class="result-card" :class="type">
    <view class="result-icon">{{ icon }}</view>
    <text class="result-title">{{ title }}</text>
    <text class="result-desc">{{ description }}</text>
    
    <view v-if="amount" class="amount-info">
      <text class="amount-label">{{ amountLabel }}</text>
      <text class="amount-value">¥{{ formattedAmount }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type ResultType = 'success' | 'pending' | 'failed' | 'warning'

interface Props {
  type: ResultType
  title: string
  description?: string
  amount?: number // 分
  amountLabel?: string
}

const props = withDefaults(defineProps<Props>(), {
  amountLabel: '支付金额',
})

const icons: Record<ResultType, string> = {
  success: '✅',
  pending: '⏳',
  failed: '❌',
  warning: '⚠️',
}

const icon = computed(() => icons[props.type])
const formattedAmount = computed(() => (props.amount! / 100).toFixed(2))
</script>

<style lang="scss" scoped>
.result-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 64rpx 48rpx;
  margin: 24rpx;
  background: var(--color-bg-card);
  border-radius: 24rpx;
  border: 2rpx solid var(--color-border);
  
  &.success {
    background: linear-gradient(135deg, rgba(0, 210, 106, 0.1) 0%, rgba(0, 210, 106, 0.05) 100%);
    border-color: var(--color-primary);
  }
  
  &.pending {
    background: linear-gradient(135deg, rgba(245, 158, 11, 0.1) 0%, rgba(245, 158, 11, 0.05) 100%);
    border-color: #F59E0B;
  }
  
  &.failed {
    background: linear-gradient(135deg, rgba(239, 68, 68, 0.1) 0%, rgba(239, 68, 68, 0.05) 100%);
    border-color: #EF4444;
  }
}

.result-icon {
  font-size: 96rpx;
  margin-bottom: 24rpx;
}

.result-title {
  font-size: 40rpx;
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: 12rpx;
}

.result-desc {
  font-size: 28rpx;
  color: var(--color-text-secondary);
  text-align: center;
}

.amount-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 32rpx;
  padding-top: 32rpx;
  border-top: 1rpx solid var(--color-border);
  width: 100%;
}

.amount-label {
  font-size: 26rpx;
  color: var(--color-text-secondary);
  margin-bottom: 8rpx;
}

.amount-value {
  font-size: 48rpx;
  font-weight: 800;
  color: var(--color-primary);
}
</style>
