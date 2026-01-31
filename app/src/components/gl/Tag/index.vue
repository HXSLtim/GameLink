<template>
  <view
    class="gl-tag"
    :class="[
      `gl-tag--${type}`,
      `gl-tag--${size}`,
      {
        'gl-tag--plain': plain,
        'gl-tag--round': round,
        'gl-tag--closable': closable,
      }
    ]"
    :style="customStyle"
    @tap="handleClick"
  >
    <uv-icon v-if="icon" :name="icon" :size="iconSize" class="gl-tag__icon" />
    <text class="gl-tag__text">
      <slot>{{ text }}</slot>
    </text>
    <view v-if="closable" class="gl-tag__close" @tap.stop="handleClose">
      <uv-icon name="close" :size="closeIconSize" />
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type TagType = 'primary' | 'success' | 'warning' | 'error' | 'info' | 'default'
type TagSize = 'mini' | 'small' | 'medium'

interface Props {
  type?: TagType
  size?: TagSize
  text?: string
  icon?: string
  plain?: boolean
  round?: boolean
  closable?: boolean
  customStyle?: string | Record<string, string>
}

const props = withDefaults(defineProps<Props>(), {
  type: 'default',
  size: 'small',
  plain: false,
  round: false,
  closable: false,
})

const emit = defineEmits<{
  click: [e: Event]
  close: [e: Event]
}>()

const iconSize = computed(() => {
  const sizes = { mini: 10, small: 12, medium: 14 }
  return sizes[props.size]
})

const closeIconSize = computed(() => {
  const sizes = { mini: 10, small: 12, medium: 14 }
  return sizes[props.size]
})

const handleClick = (e: Event) => {
  emit('click', e)
}

const handleClose = (e: Event) => {
  emit('close', e)
}
</script>

<style lang="scss" scoped>
.gl-tag {
  display: inline-flex;
  align-items: center;
  gap: 4rpx;
  border-radius: 8rpx;
  font-weight: 500;
  border: 1rpx solid transparent;
  
  // 尺寸
  &--mini {
    height: 32rpx;
    padding: 0 8rpx;
    font-size: 20rpx;
    border-radius: 6rpx;
  }
  
  &--small {
    height: 40rpx;
    padding: 0 12rpx;
    font-size: 22rpx;
  }
  
  &--medium {
    height: 48rpx;
    padding: 0 16rpx;
    font-size: 24rpx;
  }
  
  // 类型
  &--primary {
    background: rgba(0, 210, 106, 0.1);
    color: var(--color-primary);
    border-color: rgba(0, 210, 106, 0.2);
  }
  
  &--success {
    background: rgba(16, 185, 129, 0.1);
    color: var(--color-success);
    border-color: rgba(16, 185, 129, 0.2);
  }
  
  &--warning {
    background: rgba(245, 158, 11, 0.1);
    color: var(--color-warning);
    border-color: rgba(245, 158, 11, 0.2);
  }
  
  &--error {
    background: rgba(239, 68, 68, 0.1);
    color: var(--color-error);
    border-color: rgba(239, 68, 68, 0.2);
  }
  
  &--info {
    background: rgba(59, 130, 246, 0.1);
    color: var(--color-info, #3B82F6);
    border-color: rgba(59, 130, 246, 0.2);
  }
  
  &--default {
    background: var(--color-bg-secondary);
    color: var(--color-text-secondary);
    border-color: var(--color-border);
  }
  
  // Plain
  &--plain {
    background: transparent !important;
  }
  
  // 圆角
  &--round {
    border-radius: 999rpx;
  }
}

.gl-tag__icon {
  flex-shrink: 0;
}

.gl-tag__text {
  line-height: 1;
}

.gl-tag__close {
  margin-left: 4rpx;
  padding: 4rpx;
  border-radius: 50%;
  opacity: 0.7;
  
  &:active {
    opacity: 1;
    background: rgba(0, 0, 0, 0.1);
  }
}
</style>
