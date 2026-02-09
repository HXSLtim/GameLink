<template>
  <view v-if="visible" class="offline-banner" :class="[`offline-banner--${type}`]">
    <uv-icon :name="icon" size="14" :color="iconColor"></uv-icon>
    <text class="offline-text">{{ message }}</text>
    <view v-if="showAction" class="action-btn" @tap="$emit('action')">
      <text>{{ actionText }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  visible: boolean
  type?: 'warning' | 'error' | 'info'
  message?: string
  icon?: string
  showAction?: boolean
  actionText?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'warning',
  message: '网络不可用',
  icon: 'wifi-off',
  showAction: true,
  actionText: '刷新',
})

defineEmits<{
  action: []
}>()

const iconColor = computed(() => {
  const colors = {
    warning: 'var(--color-warning)',
    error: 'var(--color-error)',
    info: 'var(--color-info)',
  }
  return colors[props.type]
})
</script>

<style lang="scss" scoped>
.offline-banner {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-xs);
  padding: var(--spacing-xs) var(--spacing-md);
  background: var(--color-bg-card);
  border-bottom: 1rpx solid var(--color-border);
  
  &--warning {
    border-left: 4rpx solid var(--color-warning);
    .offline-text { color: var(--color-text); }
  }
  
  &--error {
    border-left: 4rpx solid var(--color-error);
    .offline-text { color: var(--color-text); }
  }
  
  &--info {
    border-left: 4rpx solid var(--color-info);
    .offline-text { color: var(--color-text); }
  }
}

.offline-text {
  font-size: var(--font-sm);
  font-weight: 500;
}

.action-btn {
  margin-left: var(--spacing-xs);
  padding: 2rpx var(--spacing-sm);
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
  background: var(--color-bg-secondary);
  cursor: pointer;
  @include press-effect;
  
  text {
    font-size: var(--font-xs);
    color: var(--color-text-secondary);
    font-weight: 600;
  }
  
  &:active {
    transform: scale(0.95);
  }
}
</style>
