<template>
  <view class="detail-item" @tap="$emit('click')">
    <view class="detail-icon">
      <text>{{ icon }}</text>
    </view>
    <view class="detail-info">
      <text class="detail-title">{{ title }}</text>
      <text class="detail-desc">{{ description }}</text>
      <text class="detail-time">{{ formattedTime }}</text>
    </view>
    <view class="detail-amount" :class="{ negative: amount < 0 }">
      <text>{{ amount > 0 ? '+' : '' }}¥{{ formattedAmount }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  type: string
  title: string
  description?: string
  amount: number // 分
  createdAt: string
}

const props = defineProps<Props>()

defineEmits<{
  click: []
}>()

const icons: Record<string, string> = {
  order: '💰',
  withdraw: '💳',
  bonus: '🎁',
  refund: '↩️',
}

const icon = computed(() => icons[props.type] || '💰')

const formattedAmount = computed(() => (Math.abs(props.amount) / 100).toFixed(2))

const formattedTime = computed(() => {
  const date = new Date(props.createdAt)
  return `${date.getMonth() + 1}/${date.getDate()} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
})
</script>

<style lang="scss" scoped>
.detail-item {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 24rpx 0;
  border-bottom: 1rpx solid var(--color-border);
  
  &:last-child {
    border-bottom: none;
  }
}

.detail-icon {
  width: 72rpx;
  height: 72rpx;
  background: var(--color-bg-secondary);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32rpx;
  flex-shrink: 0;
}

.detail-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.detail-title {
  font-size: 30rpx;
  font-weight: 500;
  color: var(--color-text);
}

.detail-desc {
  font-size: 24rpx;
  color: var(--color-text-secondary);
}

.detail-time {
  font-size: 22rpx;
  color: var(--color-text-placeholder);
}

.detail-amount {
  flex-shrink: 0;
  font-size: 32rpx;
  font-weight: 600;
  color: var(--color-primary);
  
  &.negative {
    color: var(--color-text);
  }
}
</style>
