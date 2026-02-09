<template>
  <view
    class="gl-card"
    :class="[
      {
        'gl-card--shadow': shadow,
        'gl-card--bordered': bordered,
        'gl-card--clickable': clickable,
        'gl-card--full': full,
      }
    ]"
    :style="customStyle"
    @tap="handleClick"
  >
    <!-- 头部 -->
    <view v-if="$slots.header || title" class="gl-card__header">
      <slot name="header">
        <view class="gl-card__title">
          <uv-icon v-if="icon" :name="icon" size="18" color="var(--color-text-secondary)" />
          <text>{{ title }}</text>
        </view>
        <view v-if="$slots.extra || extra" class="gl-card__extra">
          <slot name="extra">
            <text>{{ extra }}</text>
          </slot>
        </view>
      </slot>
    </view>
    
    <!-- 内容 -->
    <view class="gl-card__body" :style="bodyStyle">
      <slot></slot>
    </view>
    
    <!-- 底部 -->
    <view v-if="$slots.footer" class="gl-card__footer">
      <slot name="footer"></slot>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  title?: string
  extra?: string
  icon?: string
  shadow?: boolean
  bordered?: boolean
  clickable?: boolean
  full?: boolean
  padding?: string
  customStyle?: string | Record<string, string>
}

const props = withDefaults(defineProps<Props>(), {
  shadow: true,
  bordered: false,
  clickable: false,
  full: false,
  padding: '24rpx',
})

const emit = defineEmits<{
  click: [e: Event]
}>()

const bodyStyle = computed(() => ({
  padding: props.padding,
}))

const handleClick = (e: Event) => {
  if (props.clickable) {
    emit('click', e)
  }
}
</script>

<style lang="scss" scoped>
.gl-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  overflow: hidden;
  border: 1rpx solid transparent;
  
  &--shadow {
    box-shadow: none;
    border-color: var(--color-border);
  }
  
  &--bordered {
    border-color: var(--color-border);
  }
  
  &--clickable {
    transition: all 0.2s;
    cursor: pointer;
    
    &:active {
      transform: translateY(1rpx);
      background: var(--color-bg-secondary);
    }
    
    @include hover-supported {
      &:hover {
        background: var(--color-bg-hover, var(--color-bg-card));
        border-color: var(--color-border);
      }
    }
  }
  
  &--full {
    border-radius: 0;
  }
}

.gl-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-xs) var(--spacing-md);
  border-bottom: 1rpx solid var(--color-border);
}

.gl-card__title {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  font-size: var(--font-sm);
  font-weight: 600;
  color: var(--color-text);
}

.gl-card__extra {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}

.gl-card__body {
  // padding 通过 prop 控制
}

.gl-card__footer {
  padding: var(--spacing-sm) var(--spacing-md);
  border-top: 1rpx solid var(--color-border);
}
</style>
