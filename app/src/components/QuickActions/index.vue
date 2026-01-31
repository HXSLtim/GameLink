<template>
  <view class="quick-actions">
    <view 
      v-for="item in items" 
      :key="item.key" 
      class="quick-item"
      @tap="$emit('click', item.key)"
    >
      <view class="quick-icon">{{ item.icon }}</view>
      <text class="quick-label">{{ item.label }}</text>
      <view v-if="item.badge" class="quick-badge">{{ item.badge > 99 ? '99+' : item.badge }}</view>
    </view>
  </view>
</template>

<script setup lang="ts">
export interface QuickActionItem {
  key: string
  icon: string
  label: string
  badge?: number
}

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
  gap: 16rpx;
  padding: 24rpx;
  margin: 0 24rpx 20rpx;
  background: var(--color-bg-card);
  border-radius: 20rpx;
  border: 2rpx solid var(--color-border);
}

.quick-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12rpx;
  padding: 20rpx 0;
  position: relative;
  transition: all 0.2s;
  
  &:active {
    transform: scale(0.95);
  }
}

.quick-icon {
  font-size: 48rpx;
}

.quick-label {
  font-size: 24rpx;
  color: var(--color-text);
  font-weight: 500;
}

.quick-badge {
  position: absolute;
  top: 8rpx;
  right: 16rpx;
  min-width: 32rpx;
  height: 32rpx;
  padding: 0 8rpx;
  background: var(--color-error);
  border-radius: 16rpx;
  font-size: 20rpx;
  font-weight: 600;
  color: #FFFFFF;
  text-align: center;
  line-height: 32rpx;
}
</style>
