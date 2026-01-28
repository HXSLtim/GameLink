<template>
  <view class="empty-state">
    <view class="empty-icon">
      <text>{{ icon }}</text>
    </view>
    <text class="empty-title">{{ title }}</text>
    <text v-if="description" class="empty-desc">{{ description }}</text>
    <view v-if="actionText" class="empty-action">
      <button class="btn-action" @click="$emit('action')">
        {{ actionText }}
      </button>
    </view>
    <slot></slot>
  </view>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  icon?: string
  title?: string
  description?: string
  actionText?: string
}>(), {
  icon: '📭',
  title: '暂无数据',
  description: '',
  actionText: '',
})

defineEmits<{
  action: []
}>()
</script>

<style lang="scss" scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80rpx 48rpx;
  
  .empty-icon {
    font-size: 96rpx;
    margin-bottom: 24rpx;
  }
  
  .empty-title {
    font-size: 32rpx;
    font-weight: 500;
    color: var(--color-text);
    margin-bottom: 12rpx;
  }
  
  .empty-desc {
    font-size: 26rpx;
    color: var(--color-text-secondary);
    text-align: center;
    line-height: 1.6;
    margin-bottom: 32rpx;
  }
  
  .btn-action {
    padding: 16rpx 48rpx;
    background: var(--color-primary);
    color: #FFFFFF;
    font-size: 28rpx;
    border-radius: 40rpx;
    border: none;
    
    &::after {
      border: none;
    }
  }
}
</style>
