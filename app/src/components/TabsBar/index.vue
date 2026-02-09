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
import type { TabItem } from '@/types/ui'

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
  type: 'card',
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
  background: transparent;
  border-bottom: none;
  white-space: nowrap;
  padding: var(--spacing-xs) 0;
  
  &--scrollable {
    @include hide-scrollbar;
    .tabs-container {
      display: inline-flex;
      padding: 0 var(--spacing-sm);
      gap: var(--spacing-xs);
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
    padding: 0 var(--spacing-sm);
    
    .tabs-container {
      display: flex;
      gap: var(--spacing-xs);
      background: transparent;
      padding: 0;
    }
    
    .tab-item {
      flex: 1;
      justify-content: center;
      padding: var(--spacing-xs) var(--spacing-md);
      border-radius: var(--radius-full);
      
      &--active {
        background: var(--color-bg-secondary);
        box-shadow: none;
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
  gap: var(--spacing-xs);
  padding: var(--spacing-xs) var(--spacing-lg);
  position: relative;
  transition: all 0.2s;
  flex-shrink: 0;
  @include press-effect;
  border-radius: var(--radius-full);
  border: 1rpx solid var(--color-border);
  background: var(--color-bg-card);
  
  &--active {
    .tab-label {
      color: var(--color-text);
      font-weight: 600;
    }
    background: var(--color-bg-secondary);
    border-color: var(--color-border);
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
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  transition: all 0.2s;
}

.tab-badge {
  min-width: 32rpx;
  height: 32rpx;
  padding: 0 var(--spacing-xs);
  background: var(--color-error);
  color: #fff;
  font-size: var(--font-xs);
  font-weight: 600;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
