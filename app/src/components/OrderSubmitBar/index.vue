<template>
  <view class="action-bar">
    <view class="total-info">
      <text class="total-label">合计：</text>
      <PriceTag class="total-price" :amount="total" amount-unit="yuan" size="small" />
    </view>
    <GlButton 
      type="primary" 
      size="large" 
      round 
      :disabled="disabled"
      :loading="loading"
      @click="$emit('submit')"
    >
      {{ loading ? '提交中...' : buttonText }}
    </GlButton>
  </view>
</template>

<script setup lang="ts">
import GlButton from '@/components/gl/Button/index.vue'
import PriceTag from '@/components/PriceTag/index.vue'

interface Props {
  total: number
  disabled?: boolean
  loading?: boolean
  buttonText?: string
}

withDefaults(defineProps<Props>(), {
  disabled: false,
  loading: false,
  buttonText: '提交订单',
})

defineEmits<{
  submit: []
}>()
</script>

<style lang="scss" scoped>
.action-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-sm) var(--spacing-md);
  padding-bottom: calc(var(--spacing-sm) + env(safe-area-inset-bottom));
  background: var(--color-bg);
  border-top: 1rpx solid var(--color-border);
  position: sticky;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 100;
}

.total-info {
  display: flex;
  align-items: baseline;
  gap: var(--spacing-xs);
}

.total-label {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}

.total-price {
  font-size: var(--font-lg);
  font-weight: 700;
  color: var(--color-primary);
}

.total-price :deep(.price-tag) {
  color: var(--color-primary);
}
</style>
