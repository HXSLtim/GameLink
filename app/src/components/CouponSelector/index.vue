<template>
  <GlCard :shadow="false" bordered clickable class="coupon-card" @click="$emit('click')">
    <view class="coupon-row">
      <text class="coupon-icon">🎫</text>
      <text class="coupon-label">优惠券</text>
      <view class="coupon-right">
        <view v-if="selectedCoupon" class="coupon-value">
          <text>-</text>
          <PriceTag
            :amount="selectedCoupon.discount"
            amount-unit="yuan"
            size="small"
            :show-decimal="false"
          />
        </view>
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
import PriceTag from '@/components/PriceTag/index.vue'
import type { Coupon } from '@/types/coupon'

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
}

.coupon-row {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.coupon-icon {
  font-size: var(--font-base);
}

.coupon-label {
  flex: 1;
  font-size: var(--font-md);
  color: var(--color-text);
}

.coupon-right {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.coupon-value {
  display: inline-flex;
  align-items: baseline;
  gap: 4rpx;
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-error);
}

.coupon-value :deep(.price-tag) {
  color: var(--color-error);
}

.coupon-placeholder {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}
</style>
