<template>
  <view class="balance-info">
    <text class="balance-label">{{ label }}</text>
    <PriceTag
      class="balance-value"
      :amount="amount"
      :amount-unit="amountUnit"
      size="large"
    />
  </view>
</template>

<script setup lang="ts">
import PriceTag from '@/components/PriceTag/index.vue'

interface Props {
  label?: string
  amount: number
  amountUnit?: 'cents' | 'yuan'
}

withDefaults(defineProps<Props>(), {
  label: '当前余额',
  amountUnit: 'yuan',
})
</script>

<style lang="scss" scoped>
.balance-info {
  display: flex;
  align-items: baseline;
  gap: var(--spacing-md);
  padding: var(--spacing-md) var(--spacing-md);
}

.balance-label {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}

.balance-value {
  font-size: var(--font-xl);
  font-weight: 700;
  color: var(--color-text);
}

.balance-value :deep(.price-tag) {
  color: var(--color-text);
}

.balance-value :deep(.price-tag .amount) {
  font-weight: 700;
}
</style>
