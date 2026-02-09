<template>
  <view class="action-bar">
    <view class="total-info">
      <text class="total-label">实付</text>
      <PriceTag class="total-value" :amount="total" amount-unit="yuan" size="small" />
      <view v-if="bonus > 0" class="bonus-text">
        <text>（含赠送</text>
        <PriceTag :amount="bonus" amount-unit="yuan" size="small" :show-decimal="false" />
        <text>）</text>
      </view>
    </view>
    <GlButton
      type="primary"
      size="large"
      round
      :disabled="disabled"
      :loading="loading"
      @click="$emit('submit')"
    >
      立即充值
    </GlButton>
  </view>
</template>

<script setup lang="ts">
import GlButton from '@/components/gl/Button/index.vue'
import PriceTag from '@/components/PriceTag/index.vue'

interface Props {
  total: number
  bonus: number
  disabled?: boolean
  loading?: boolean
}

withDefaults(defineProps<Props>(), {
  disabled: false,
  loading: false,
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
  padding: 20rpx 24rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background: var(--color-bg-card);
  border-top: 1rpx solid var(--color-border);
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 100;
}

.total-info {
  display: flex;
  align-items: baseline;
  gap: 8rpx;
}

.total-label {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}

.total-value {
  font-size: var(--font-lg);
  font-weight: 800;
  color: var(--color-primary);
}

.bonus-text {
  display: inline-flex;
  align-items: baseline;
  gap: 4rpx;
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}

.total-value :deep(.price-tag) {
  color: var(--color-primary);
}

.bonus-text :deep(.price-tag) {
  color: var(--color-text-secondary);
}
</style>
