<template>
  <view class="record-item">
    <view class="record-icon" :class="record.type">
      <uv-icon :name="iconName" size="20" :color="iconColor"></uv-icon>
    </view>
    <view class="record-info">
      <text class="record-title">{{ record.title }}</text>
      <text class="record-desc">{{ record.description }}</text>
      <text class="record-time">{{ formattedTime }}</text>
    </view>
    <view class="record-amount" :class="{ positive: record.amount > 0 }">
      <text>{{ record.amount > 0 ? '+' : '' }}{{ formattedAmount }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

export type TransactionType = 'recharge' | 'withdraw' | 'payment' | 'refund' | 'earning' | 'bonus'

export interface TransactionData {
  id: number
  type: TransactionType
  title: string
  description: string
  amount: number // 分
  createdAt: string
}

interface Props {
  record: TransactionData
}

const props = defineProps<Props>()

const iconConfigs: Record<TransactionType, { name: string; color: string }> = {
  recharge: { name: 'plus-circle', color: '#10B981' },
  withdraw: { name: 'red-packet', color: '#F59E0B' },
  payment: { name: 'shopping-cart', color: '#3B82F6' },
  refund: { name: 'reload', color: '#8B5CF6' },
  earning: { name: 'money', color: '#10B981' },
  bonus: { name: 'gift', color: '#EC4899' },
}

const iconName = computed(() => iconConfigs[props.record.type]?.name || 'list')
const iconColor = computed(() => iconConfigs[props.record.type]?.color || '#9CA3AF')

const formattedAmount = computed(() => (props.record.amount / 100).toFixed(2))

const formattedTime = computed(() => {
  const date = new Date(props.record.createdAt)
  return `${date.getMonth() + 1}/${date.getDate()} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
})
</script>

<style lang="scss" scoped>
.record-item {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 24rpx 0;
  border-bottom: 1rpx solid var(--color-border);
  
  &:last-child {
    border-bottom: none;
  }
}

.record-icon {
  width: 72rpx;
  height: 72rpx;
  border-radius: 50%;
  background: var(--color-bg-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.record-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.record-title {
  font-size: 30rpx;
  font-weight: 500;
  color: var(--color-text);
}

.record-desc {
  font-size: 24rpx;
  color: var(--color-text-secondary);
}

.record-time {
  font-size: 22rpx;
  color: var(--color-text-placeholder);
}

.record-amount {
  flex-shrink: 0;
  font-size: 32rpx;
  font-weight: 600;
  color: var(--color-text);
  
  &.positive {
    color: var(--color-primary);
  }
}
</style>
