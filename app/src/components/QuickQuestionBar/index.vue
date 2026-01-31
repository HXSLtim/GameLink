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
          <text class="quick-icon">{{ item.icon }}</text>
          <text class="quick-text">{{ item.text }}</text>
          <text class="quick-arrow">›</text>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
export interface QuickQuestion {
  id: number
  icon: string
  text: string
  answer: string
}

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
  padding: 16rpx 24rpx;
  background: var(--color-bg);
}

.section-title {
  font-size: 26rpx;
  color: var(--color-text-secondary);
  margin-bottom: 12rpx;
}

.quick-scroll {
  white-space: nowrap;
}

.quick-list {
  display: inline-flex;
  gap: 16rpx;
  padding-bottom: 8rpx;
}

.quick-item {
  display: inline-flex;
  align-items: center;
  gap: 8rpx;
  padding: 12rpx 20rpx;
  background: var(--color-bg-card);
  border-radius: 24rpx;
  flex-shrink: 0;
}

.quick-icon {
  font-size: 24rpx;
}

.quick-text {
  font-size: 26rpx;
  color: var(--color-text);
}

.quick-arrow {
  font-size: 24rpx;
  color: var(--color-text-placeholder);
}
</style>
