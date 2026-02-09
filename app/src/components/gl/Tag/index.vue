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
  border-radius: var(--radius-sm);
  font-weight: 500;
  border: 1rpx solid transparent;
  
  // 尺寸
  &--mini {
    height: 32rpx;
    padding: 0 var(--spacing-xs);
    font-size: var(--font-xs);
    border-radius: var(--radius-sm);
  }
  
  &--small {
    height: 40rpx;
    padding: 0 var(--spacing-sm);
    font-size: var(--font-sm);
  }
  
  &--medium {
    height: 48rpx;
    padding: 0 var(--spacing-md);
    font-size: var(--font-md);
  }
  
  // 类型 - 使用 CSS 变量适配日/夜主题
  &--primary {
    background: var(--color-primary-tint);
    color: var(--color-primary);
    border-color: var(--color-primary-tint-border);
  }
  
  &--success {
    background: var(--color-success-tint);
    color: var(--color-success);
    border-color: var(--color-success-tint-border);
  }
  
  &--warning {
    background: var(--color-warning-tint);
    color: var(--color-warning);
    border-color: var(--color-warning-tint-border);
  }
  
  &--error {
    background: var(--color-error-tint);
    color: var(--color-error);
    border-color: var(--color-error-tint-border);
  }
  
  &--info {
    background: var(--color-info-tint);
    color: var(--color-info);
    border-color: var(--color-info-tint-border);
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
  transition: all 0.15s ease;
  
  &:active {
    opacity: 1;
    background: var(--color-bg-secondary);
  }
}
</style>
