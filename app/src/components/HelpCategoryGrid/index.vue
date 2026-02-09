<template>
  <view class="category-section">
    <view class="category-grid">
      <view
        v-for="cat in categories"
        :key="cat.id"
        class="category-item"
        :class="{ active: selectedId === cat.id }"
        @tap="$emit('select', cat.id)"
      >
        <view class="category-icon"><uv-icon :name="cat.icon" size="22" color="var(--color-text-secondary)" /></view>
        <text class="category-name">{{ cat.name }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import type { HelpCategory } from '@/types/help'

interface Props {
  categories: HelpCategory[]
  selectedId: string | null
}

defineProps<Props>()

defineEmits<{
  select: [id: number]
}>()
</script>

<style lang="scss" scoped>
.category-section {
  padding: var(--spacing-md);
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--spacing-sm);
}

.category-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-xs);
  padding: var(--spacing-sm) var(--spacing-sm);
  background: var(--color-bg-card);
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
  transition: all 0.2s;
  cursor: pointer;
  @include press-effect;

  &:active,
  &.active {
    border-color: var(--color-border);
    background: var(--color-bg-secondary);
  }
}

.category-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.category-name {
  font-size: var(--font-xs);
  color: var(--color-text);
  font-weight: 500;
}
</style>
