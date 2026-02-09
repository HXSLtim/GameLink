<template>
  <SectionCard :title="title">
    <template #extra>
      <slot name="extra">
        <text v-if="subtitle" class="stats-date">{{ subtitle }}</text>
      </slot>
    </template>
    
    <view class="stats-grid" :style="gridStyle">
      <view 
        v-for="item in items" 
        :key="item.label" 
        class="stat-item"
        :class="{ 'stat-item--clickable': !!item.onClick }"
        @tap="item.onClick?.()"
      >
        <text class="stat-value" :class="{ highlight: item.highlight }">{{ item.value }}</text>
        <text class="stat-label">{{ item.label }}</text>
        <text v-if="item.unit" class="stat-unit">{{ item.unit }}</text>
      </view>
    </view>
  </SectionCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import SectionCard from '@/components/SectionCard/index.vue'
import type { StatItem } from '@/types/ui'

interface Props {
  title: string
  subtitle?: string
  items: StatItem[]
  /** 网格列数，默认自适应（≤3 项等分，>3 项 4 列） */
  columns?: number
}

const props = withDefaults(defineProps<Props>(), {
  columns: 0,
})

const gridStyle = computed(() => {
  const cols = props.columns > 0
    ? props.columns
    : Math.min(props.items.length, 4)
  return { gridTemplateColumns: `repeat(${cols}, 1fr)` }
})
</script>

<style lang="scss" scoped>
.stats-date {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}

.stats-grid {
  display: grid;
  gap: var(--spacing-xs);
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-xs);
  padding: var(--spacing-xs) 0;

  &--clickable {
    cursor: pointer;
    @include press-effect;
  }
}

.stat-value {
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
  
  &.highlight {
    color: var(--color-primary);
    font-weight: 700;
  }
}

.stat-label {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}

.stat-unit {
  font-size: var(--font-xs);
  color: var(--color-text-placeholder);
  margin-top: -4rpx;
}
</style>
