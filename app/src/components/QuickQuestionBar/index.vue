<template>
  <view class="quick-section">
    <text class="section-title">快捷问题</text>
    <scroll-view class="quick-scroll" scroll-x :show-scrollbar="false">
      <view class="quick-list">
        <view
          v-for="item in questions"
          :key="item.id"
          class="quick-item"
          @tap="$emit('select', item)"
        >
          <view class="quick-icon"><uv-icon :name="item.icon" size="16" color="var(--color-text-secondary)" /></view>
          <text class="quick-text">{{ item.text }}</text>
          <text class="quick-arrow">›</text>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import type { QuickQuestion } from '@/types/CustomerService'

interface Props {
  questions: QuickQuestion[]
}

defineProps<Props>()

defineEmits<{
  select: [question: QuickQuestion]
}>()
</script>

<style lang="scss" scoped>
.quick-section {
  padding: var(--spacing-sm) var(--spacing-md);
  background: var(--color-bg);
}

.section-title {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
  margin-bottom: var(--spacing-xs);
}

.quick-scroll {
  white-space: nowrap;
  @include hide-scrollbar;
}

.quick-list {
  display: inline-flex;
  gap: var(--spacing-xs);
  padding-bottom: var(--spacing-xs);
}

.quick-item {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-xs);
  padding: 2rpx var(--spacing-md);
  background: var(--color-bg-card);
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
  flex-shrink: 0;
  cursor: pointer;
  @include press-effect;
}

.quick-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.quick-text {
  font-size: var(--font-xs);
  color: var(--color-text);
}

.quick-arrow {
  font-size: var(--font-xs);
  color: var(--color-text-placeholder);
}
</style>
