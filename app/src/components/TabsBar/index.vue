<template>
  <scroll-view
    class="tabs-bar"
    :class="{ 
      'tabs-bar--scrollable': scrollable,
      'tabs-bar--stretch': stretch,
      'tabs-bar--card': type === 'card'
    }"
    :scroll-x="scrollable"
    :scroll-with-animation="true"
    :scroll-into-view="scrollIntoView"
  >
    <view class="tabs-container">
      <view
        v-for="tab in tabs"
        :key="tab.key"
        :id="`tab-${tab.key}`"
        class="tab-item"
        :class="{ 
          'tab-item--active': modelValue === tab.key,
          'tab-item--disabled': tab.disabled
        }"
        @tap="handleChange(tab)"
      >
        <text class="tab-label">{{ tab.label }}</text>
        <view v-if="tab.badge" class="tab-badge">
          {{ formatBadge(tab.badge) }}
        </view>
      </view>
    </view>
  </scroll-view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

export interface TabItem {
  key: string
  label: string
  badge?: number | string
  disabled?: boolean
}

interface Props {
  tabs: TabItem[]
  modelValue: string
  scrollable?: boolean
  stretch?: boolean
  type?: 'line' | 'card'
}

const props = withDefaults(defineProps<Props>(), {
  scrollable: false,
  stretch: false,
  type: 'line',
})

const emit = defineEmits<{
  'update:modelValue': [key: string]
  change: [key: string, tab: TabItem]
}>()

const scrollIntoView = computed(() => {
  return props.scrollable ? `tab-${props.modelValue}` : ''
})

const formatBadge = (badge: number | string) => {
  if (typeof badge === 'number' && badge > 99) {
    return '99+'
  }
  return badge
}

const handleChange = (tab: TabItem) => {
  if (tab.disabled) return
  emit('update:modelValue', tab.key)
  emit('change', tab.key, tab)
}
</script>

<style lang="scss" scoped>
.tabs-bar {
  background: var(--color-bg-card);
  border-bottom: 1rpx solid var(--color-border);
  white-space: nowrap;
  
  &--scrollable {
    .tabs-container {
      display: inline-flex;
      padding: 0 16rpx;
    }
  }
  
  &--stretch {
    .tabs-container {
      display: flex;
    }
    
    .tab-item {
      flex: 1;
      justify-content: center;
    }
  }
  
  &--card {
    background: transparent;
    border-bottom: none;
    padding: 16rpx;
    
    .tabs-container {
      display: flex;
      gap: 16rpx;
      background: var(--color-bg-secondary);
      border-radius: 16rpx;
      padding: 8rpx;
    }
    
    .tab-item {
      flex: 1;
      justify-content: center;
      padding: 16rpx 24rpx;
      border-radius: 12rpx;
      
      &--active {
        background: var(--color-bg-card);
        box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.1);
        
        &::after {
          display: none;
        }
      }
    }
  }
}

.tabs-container {
  display: flex;
  min-width: 100%;
}

.tab-item {
  display: flex;
  align-items: center;
  gap: 8rpx;
  padding: 24rpx 32rpx;
  position: relative;
  transition: all 0.2s;
  flex-shrink: 0;
  
  &--active {
    .tab-label {
      color: var(--color-primary);
      font-weight: 600;
    }
    
    &::after {
      content: '';
      position: absolute;
      bottom: 0;
      left: 50%;
      transform: translateX(-50%);
      width: 48rpx;
      height: 6rpx;
      background: var(--color-primary);
      border-radius: 3rpx;
    }
  }
  
  &--disabled {
    opacity: 0.5;
    pointer-events: none;
  }
  
  &:active {
    opacity: 0.7;
  }
}

.tab-label {
  font-size: 28rpx;
  color: var(--color-text-secondary);
  transition: all 0.2s;
}

.tab-badge {
  min-width: 32rpx;
  height: 32rpx;
  padding: 0 8rpx;
  background: var(--color-error, #EF4444);
  color: #fff;
  font-size: 20rpx;
  font-weight: 600;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
