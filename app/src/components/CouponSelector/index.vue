<template>
  <GlCard :shadow="false" bordered class="coupon-card" @click="$emit('click')">
    <view class="coupon-row">
      <text class="coupon-icon">🎫</text>
      <text class="coupon-label">优惠券</text>
      <view class="coupon-right">
        <text v-if="selectedCoupon" class="coupon-value">-¥{{ selectedCoupon.discount }}</text>
        <text v-else class="coupon-placeholder">
          {{ availableCount ? `${availableCount}张可用` : '暂无可用' }}
        </text>
        <uv-icon name="arrow-right" size="14" color="var(--color-text-secondary)"></uv-icon>
      </view>
    </view>
  </GlCard>
</template>

<script setup lang="ts">
import GlCard from '@/components/gl/Card/index.vue'

export interface Coupon {
  id: number
  name: string
  discount: number
  minAmount?: number
  expireAt?: string
}

interface Props {
  selectedCoupon?: Coupon | null
  availableCount?: number
}

withDefaults(defineProps<Props>(), {
  availableCount: 0,
})

defineEmits<{
  click: []
}>()
</script>

<style lang="scss" scoped>
.coupon-card {
  cursor: pointer;
}

.coupon-row {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.coupon-icon {
  font-size: 32rpx;
}

.coupon-label {
  flex: 1;
  font-size: 28rpx;
  color: var(--color-text);
}

.coupon-right {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.coupon-value {
  font-size: 28rpx;
  font-weight: 600;
  color: var(--color-error);
}

.coupon-placeholder {
  font-size: 26rpx;
  color: var(--color-text-secondary);
}
</style>
