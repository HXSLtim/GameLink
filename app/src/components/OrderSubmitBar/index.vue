<template>
  <view class="action-bar">
    <view class="total-info">
      <text class="total-label">合计：</text>
      <text class="total-price">¥{{ total.toFixed(2) }}</text>
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
  padding: 20rpx 24rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background: var(--color-bg-card);
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
  gap: 8rpx;
}

.total-label {
  font-size: 28rpx;
  color: var(--color-text-secondary);
}

.total-price {
  font-size: 40rpx;
  font-weight: 800;
  color: var(--color-primary);
}
</style>
