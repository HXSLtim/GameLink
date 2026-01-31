<template>
  <GlCard :title="title" :shadow="false" bordered class="stats-card">
    <template #extra>
      <text class="stats-date">{{ subtitle }}</text>
    </template>
    
    <view class="stats-grid">
      <view 
        v-for="item in items" 
        :key="item.label" 
        class="stat-item"
        @tap="item.onClick?.()"
      >
        <text class="stat-value" :class="{ highlight: item.highlight }">{{ item.value }}</text>
        <text class="stat-label">{{ item.label }}</text>
      </view>
    </view>
  </GlCard>
</template>

<script setup lang="ts">
import GlCard from '@/components/gl/Card/index.vue'

export interface StatItem {
  value: string | number
  label: string
  highlight?: boolean
  onClick?: () => void
}

interface Props {
  title: string
  subtitle?: string
  items: StatItem[]
}

defineProps<Props>()
</script>

<style lang="scss" scoped>
.stats-card {
  margin: 0 24rpx 20rpx;
}

.stats-date {
  font-size: 24rpx;
  color: var(--color-text-secondary);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16rpx;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;
  padding: 16rpx 0;
}

.stat-value {
  font-size: 36rpx;
  font-weight: 700;
  color: var(--color-text);
  
  &.highlight {
    color: var(--color-primary);
  }
}

.stat-label {
  font-size: 24rpx;
  color: var(--color-text-secondary);
}
</style>
