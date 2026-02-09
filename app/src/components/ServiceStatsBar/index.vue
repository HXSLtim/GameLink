<template>
  <view class="stats-bar">
    <view v-for="item in items" :key="item.label" class="stat-item">
      <text class="stat-value" :class="{ highlight: item.highlight }">{{ item.value }}</text>
      <text class="stat-label">{{ item.label }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import type { StatItem } from '@/types/ui'

interface Props {
  items: StatItem[]
}

defineProps<Props>()
</script>

<style lang="scss" scoped>
.stats-bar {
  display: flex;
  padding: var(--spacing-sm);
  margin: var(--spacing-md);
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1rpx solid var(--color-border);
}

.stat-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-xs);
  @include press-effect;
}

.stat-value {
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);

  &.highlight {
    color: var(--color-primary);
  }
}

.stat-label {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}
</style>
