<template>
  <view class="quick-actions">
    <view 
      v-for="item in items" 
      :key="item.key" 
      class="quick-item"
      @tap="$emit('click', item.key)"
    >
      <view class="quick-icon"><uv-icon :name="item.icon" size="24" color="var(--color-text-secondary)" /></view>
      <text class="quick-label">{{ item.label }}</text>
      <view v-if="item.badge" class="quick-badge">{{ item.badge > 99 ? '99+' : item.badge }}</view>
    </view>
  </view>
</template>

<script setup lang="ts">
import type { QuickActionItem } from '@/types/ui'

interface Props {
  items: QuickActionItem[]
}

defineProps<Props>()

defineEmits<{
  click: [key: string]
}>()
</script>

<style lang="scss" scoped>
.quick-actions {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--spacing-xs);
  padding: var(--spacing-sm);
  margin: 0 var(--spacing-md) var(--spacing-sm);
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1rpx solid var(--color-border);
}

.quick-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-xs);
  padding: var(--spacing-xs) 0;
  position: relative;
  transition: all 0.2s;
  cursor: pointer;
  @include press-effect;
}

.quick-icon {
  @include flex-center;
  width: 56rpx;
  height: 56rpx;
  background: var(--color-bg-secondary);
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
}

.quick-label {
  font-size: var(--font-xs);
  color: var(--color-text);
  font-weight: 500;
}

.quick-badge {
  position: absolute;
  top: 4rpx;
  right: 8rpx;
  min-width: 28rpx;
  height: 28rpx;
  padding: 0 var(--spacing-xs);
  background: var(--color-error);
  border-radius: var(--radius-full);
  font-size: var(--font-xs);
  font-weight: 600;
  color: #FFFFFF;
  text-align: center;
  line-height: 28rpx;
}
</style>
