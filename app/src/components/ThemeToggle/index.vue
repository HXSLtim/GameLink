<template>
  <view class="theme-toggle" :class="{ 'theme-dark': isDarkMode }" @tap="toggleTheme">
    <view class="toggle-track">
      <view class="toggle-icon sun"><image class="toggle-svg" src="/static/icons/sun.svg" mode="aspectFit" /></view>
      <view class="toggle-icon moon"><image class="toggle-svg" src="/static/icons/moon.svg" mode="aspectFit" /></view>
      <view class="toggle-thumb"></view>
    </view>
    <text v-if="showLabel" class="toggle-label">{{ isDarkMode ? '夜间模式' : '日间模式' }}</text>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useTheme } from '@/composables/useTheme'

withDefaults(defineProps<{
  showLabel?: boolean
}>(), {
  showLabel: false,
})

const { themeMode, isDark, toggleTheme: toggle } = useTheme()

const isDarkMode = computed(() => isDark.value || themeMode.value === 'dark')

const toggleTheme = () => {
  toggle()
}
</script>

<style lang="scss" scoped>
.theme-toggle {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-sm);
  cursor: pointer;
  @include press-effect;
}

.toggle-track {
  position: relative;
  width: 92rpx;
  height: 40rpx;
  background: var(--color-bg-card);
  border-radius: var(--radius-sm);
  padding: 4rpx;
  border: 1rpx solid var(--color-border);
  transition: background 0.3s ease;
}

.theme-dark .toggle-track {
  background: var(--color-bg-secondary);
}

.toggle-icon {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: opacity 0.3s ease;
  
  &.sun {
    left: 8rpx;
    opacity: 1;
  }
  
  &.moon {
    right: 8rpx;
    opacity: 0.3;
  }
}

.toggle-svg {
  width: 24rpx;
  height: 24rpx;
}

.theme-dark .toggle-icon {
  &.sun {
    opacity: 0.3;
  }
  
  &.moon {
    opacity: 1;
  }
}

.toggle-thumb {
  position: absolute;
  top: 4rpx;
  left: 4rpx;
  width: 32rpx;
  height: 32rpx;
  background: var(--color-primary);
  border-radius: var(--radius-sm);
  box-shadow: none;
  transition: transform 0.3s ease;
}

.theme-dark .toggle-thumb {
  transform: translateX(52rpx);
}

.toggle-label {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}
</style>
