<template>
  <GlCard title="费用明细" :shadow="false" bordered class="section-card">
    <view class="fee-list">
      <view v-for="item in fees" :key="item.label" class="fee-row">
        <text class="fee-label">{{ item.label }}</text>
        <text class="fee-value" :class="{ discount: item.isDiscount }">
          {{ item.isDiscount ? '-' : '' }}¥{{ item.value.toFixed(2) }}
        </text>
      </view>
    </view>
    <view class="fee-total">
      <text class="total-label">实付金额</text>
      <text class="total-value">¥{{ total.toFixed(2) }}</text>
    </view>
  </GlCard>
</template>

<script setup lang="ts">
import GlCard from '@/components/gl/Card/index.vue'

export interface FeeItem {
  label: string
  value: number
  isDiscount?: boolean
}

interface Props {
  fees: FeeItem[]
  total: number
}

defineProps<Props>()
</script>

<style lang="scss" scoped>
.section-card {
  margin: 20rpx 24rpx;
}

.fee-list {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
  padding-bottom: 24rpx;
  border-bottom: 1rpx solid var(--color-border);
}

.fee-row {
  display: flex;
  justify-content: space-between;
}

.fee-label {
  font-size: 28rpx;
  color: var(--color-text-secondary);
}

.fee-value {
  font-size: 28rpx;
  color: var(--color-text);
  
  &.discount {
    color: var(--color-error);
  }
}

.fee-total {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 24rpx;
}

.total-label {
  font-size: 30rpx;
  font-weight: 600;
  color: var(--color-text);
}

.total-value {
  font-size: 40rpx;
  font-weight: 700;
  color: var(--color-primary);
}
</style>
