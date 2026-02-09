<template>
  <SectionCard title="费用明细" margin="var(--spacing-sm) var(--spacing-md)">
    <view class="fee-list">
      <view v-for="item in fees" :key="item.label" class="fee-row">
        <text class="fee-label">{{ item.label }}</text>
        <view class="fee-value" :class="{ discount: item.isDiscount }">
          <text v-if="item.isDiscount" class="fee-sign">-</text>
          <PriceTag :amount="item.value" amount-unit="yuan" size="small" />
        </view>
      </view>
    </view>
    <view class="fee-total">
      <text class="total-label">实付金额</text>
      <PriceTag class="total-value" :amount="total" amount-unit="yuan" size="small" />
    </view>
  </SectionCard>
</template>

<script setup lang="ts">
import SectionCard from '@/components/SectionCard/index.vue'
import PriceTag from '@/components/PriceTag/index.vue'
import type { FeeItem } from '@/types/order'

interface Props {
  fees: FeeItem[]
  total: number
}

defineProps<Props>()
</script>

<style lang="scss" scoped>
.fee-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
  padding-bottom: var(--spacing-sm);
  border-bottom: 1rpx solid var(--color-border);
}

.fee-row {
  display: flex;
  justify-content: space-between;
}

.fee-label {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}

.fee-value {
  display: inline-flex;
  align-items: baseline;
  gap: 4rpx;
  font-size: var(--font-sm);
  color: var(--color-text);
  
  &.discount {
    color: var(--color-error);
  }
}

.fee-value :deep(.price-tag) {
  color: inherit;
}

.fee-sign {
  font-size: var(--font-sm);
}

.fee-total {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: var(--spacing-sm);
}

.total-label {
  font-size: var(--font-sm);
  font-weight: 600;
  color: var(--color-text);
}

.total-value {
  font-size: var(--font-base);
  font-weight: 700;
  color: var(--color-text);
}

.total-value :deep(.price-tag) {
  color: var(--color-text);
}
</style>
