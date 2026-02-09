<template>
  <view
    class="gl-button"
    :class="[
      `gl-button--${type}`,
      `gl-button--${size}`,
      {
        'gl-button--plain': plain,
        'gl-button--round': round,
        'gl-button--block': block,
        'gl-button--disabled': disabled,
        'gl-button--loading': loading,
      }
    ]"
    :style="customStyle"
    @tap="handleClick"
  >
    <uv-loading-icon v-if="loading" :size="loadingSize" :color="loadingColor" mode="circle" />
    <uv-icon v-else-if="icon" :name="icon" :size="iconSize" :color="iconColor" />
    <text v-if="$slots.default || text" class="gl-button__text">
      <slot>{{ text }}</slot>
    </text>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type ButtonType = 'primary' | 'success' | 'warning' | 'error' | 'info' | 'default'
type ButtonSize = 'mini' | 'small' | 'medium' | 'large'

interface Props {
  type?: ButtonType
  size?: ButtonSize
  text?: string
  icon?: string
  plain?: boolean
  round?: boolean
  block?: boolean
  disabled?: boolean
  loading?: boolean
  customStyle?: string | Record<string, string>
}

const props = withDefaults(defineProps<Props>(), {
  type: 'default',
  size: 'medium',
  plain: false,
  round: false,
  block: false,
  disabled: false,
  loading: false,
})

const emit = defineEmits<{
  click: [e: Event]
}>()

const iconSize = computed(() => {
  const sizes = { mini: 12, small: 14, medium: 16, large: 18 }
  return sizes[props.size]
})

const loadingSize = computed(() => {
  const sizes = { mini: 12, small: 14, medium: 16, large: 18 }
  return sizes[props.size]
})

const iconColor = computed(() => {
  if (props.plain) {
    const colors: Record<ButtonType, string> = {
      primary: 'var(--color-primary)',
      success: 'var(--color-success)',
      warning: 'var(--color-warning)',
      error: 'var(--color-error)',
      info: 'var(--color-info)',
      default: 'var(--color-text)',
    }
    return colors[props.type]
  }
  return props.type === 'default' ? 'var(--color-text)' : '#fff'
})

const loadingColor = computed(() => iconColor.value)

const handleClick = (e: Event) => {
  if (props.disabled || props.loading) return
  emit('click', e)
}
</script>

<style lang="scss" scoped>
.gl-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-xs);
  border-radius: var(--radius-md);
  font-weight: 600;
  transition: all 0.2s;
  border: 1rpx solid transparent;
  cursor: pointer;
  
  &:active:not(.gl-button--disabled):not(.gl-button--loading) {
    transform: scale(0.96);
    opacity: 0.9;
  }
  
  @include hover-supported {
    &:hover:not(.gl-button--disabled):not(.gl-button--loading) {
      filter: brightness(1.02);
    }
  }
  
  // 尺寸
  &--mini {
    height: 48rpx;
    padding: 0 var(--spacing-sm);
    font-size: var(--font-xs);
    border-radius: var(--radius-sm);
  }
  
  &--small {
    height: 60rpx;
    padding: 0 var(--spacing-md);
    font-size: var(--font-sm);
  }
  
  &--medium {
    height: 80rpx;
    padding: 0 var(--spacing-lg);
    font-size: var(--font-md);
  }
  
  &--large {
    height: 96rpx;
    padding: 0 var(--spacing-xl);
    font-size: var(--font-base);
    border-radius: var(--radius-md);
  }
  
  // 类型
  &--primary {
    background: var(--color-primary);
    color: #fff;
  }
  
  &--success {
    background: var(--color-success);
    color: #fff;
  }
  
  &--warning {
    background: var(--color-warning);
    color: #fff;
  }
  
  &--error {
    background: var(--color-error);
    color: #fff;
  }
  
  &--info {
    background: var(--color-info);
    color: #fff;
  }
  
  &--default {
    background: var(--color-bg-secondary);
    color: var(--color-text);
    border-color: var(--color-border);
  }
  
  // Plain 样式
  &--plain {
    background: transparent !important;
    box-shadow: none !important;
    
    &.gl-button--primary {
      border-color: var(--color-primary);
      color: var(--color-primary);
    }
    &.gl-button--success {
      border-color: var(--color-success);
      color: var(--color-success);
    }
    &.gl-button--warning {
      border-color: var(--color-warning);
      color: var(--color-warning);
    }
    &.gl-button--error {
      border-color: var(--color-error);
      color: var(--color-error);
    }
    &.gl-button--info {
      border-color: var(--color-info);
      color: var(--color-info);
    }
    &.gl-button--default {
      border-color: var(--color-border);
      color: var(--color-text);
    }
  }
  
  // 圆角
  &--round {
    border-radius: 999rpx;
  }
  
  // 块级
  &--block {
    display: flex;
    width: 100%;
  }
  
  // 禁用
  &--disabled {
    opacity: 0.5;
    cursor: not-allowed;
    pointer-events: none;
  }
  
  // 加载中
  &--loading {
    opacity: 0.7;
    pointer-events: none;
  }
}

.gl-button__text {
  line-height: 1;
}
</style>
