<template>
  <view
    class="gl-avatar"
    :class="[
      `gl-avatar--${size}`,
      `gl-avatar--${shape}`,
      { 'gl-avatar--bordered': bordered }
    ]"
    :style="[sizeStyle, customStyle]"
    @tap="$emit('click')"
  >
    <image
      v-if="src"
      :src="src"
      class="gl-avatar__image"
      mode="aspectFill"
      @error="handleError"
    />
    <view v-else class="gl-avatar__placeholder">
      <uv-icon v-if="icon" :name="icon" :size="iconSize" color="var(--color-text-secondary)" />
      <text v-else-if="text" class="gl-avatar__text">{{ displayText }}</text>
      <uv-icon v-else name="account" :size="iconSize" color="var(--color-text-secondary)" />
    </view>
    
    <!-- 徽章 -->
    <view v-if="badge" class="gl-avatar__badge">
      {{ badge }}
    </view>
    
    <!-- 状态点 -->
    <view
      v-if="status"
      class="gl-avatar__status"
      :class="`gl-avatar__status--${status}`"
    />
  </view>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

type AvatarSize = 'mini' | 'small' | 'medium' | 'large' | 'xlarge' | number
type AvatarShape = 'circle' | 'square'
type AvatarStatus = 'online' | 'offline' | 'busy'

interface Props {
  src?: string
  size?: AvatarSize
  shape?: AvatarShape
  icon?: string
  text?: string
  bordered?: boolean
  badge?: string | number
  status?: AvatarStatus
  customStyle?: string | Record<string, string>
}

const props = withDefaults(defineProps<Props>(), {
  size: 'medium',
  shape: 'circle',
  bordered: false,
})

defineEmits<{
  click: []
  error: [e: Event]
}>()

const fallbackSrc = ref('')

const sizeMap: Record<string, number> = {
  mini: 48,
  small: 64,
  medium: 80,
  large: 100,
  xlarge: 120,
}

const numericSize = computed(() => {
  if (typeof props.size === 'number') return props.size
  return sizeMap[props.size] || 80
})

const sizeStyle = computed(() => {
  const s = numericSize.value
  return {
    width: `${s}rpx`,
    height: `${s}rpx`,
  }
})

const iconSize = computed(() => {
  return Math.floor(numericSize.value * 0.5)
})

const displayText = computed(() => {
  if (!props.text) return ''
  return props.text.slice(0, 2).toUpperCase()
})

const handleError = (e: Event) => {
  fallbackSrc.value = ''
}
</script>

<style lang="scss" scoped>
.gl-avatar {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  background: var(--color-bg-secondary);
  flex-shrink: 0;
  
  &--circle {
    border-radius: 50%;
  }
  
  &--square {
    border-radius: 12rpx;
  }
  
  &--bordered {
    border: 4rpx solid var(--color-bg-card);
    box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.1);
  }
}

.gl-avatar__image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.gl-avatar__placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
}

.gl-avatar__text {
  font-weight: 600;
  color: var(--color-text-secondary);
}

.gl-avatar__badge {
  position: absolute;
  top: -4rpx;
  right: -4rpx;
  min-width: 32rpx;
  height: 32rpx;
  padding: 0 8rpx;
  background: var(--color-error);
  color: #fff;
  font-size: 20rpx;
  font-weight: 600;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2rpx solid var(--color-bg-card);
}

.gl-avatar__status {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 20rpx;
  height: 20rpx;
  border-radius: 50%;
  border: 3rpx solid var(--color-bg-card);
  
  &--online {
    background: var(--color-success, #10B981);
  }
  
  &--offline {
    background: var(--color-text-placeholder, #9CA3AF);
  }
  
  &--busy {
    background: var(--color-warning, #F59E0B);
  }
}
</style>
