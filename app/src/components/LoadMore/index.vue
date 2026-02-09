<template>
  <view v-if="loading || noMore" class="load-more" :class="{ 'load-more-loading': loading }">
    <view v-if="loading" class="load-more-spinner">
      <view class="spinner"></view>
      <text class="load-more-text">{{ loadingText }}</text>
    </view>
    <view v-else-if="noMore" class="load-more-end">
      <view class="load-more-line"></view>
      <text class="load-more-text">{{ noMoreText }}</text>
      <view class="load-more-line"></view>
    </view>
    <!-- 空闲状态不显示按钮，依赖滚动到底部自动加载 -->
  </view>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  loading?: boolean
  noMore?: boolean
  loadingText?: string
  noMoreText?: string
}>(), {
  loading: false,
  noMore: false,
  loadingText: '加载中...',
  noMoreText: '没有更多了',
})
</script>

<style lang="scss" scoped>
.load-more {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-sm) 0;
  
  &-spinner {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    
    .spinner {
      width: 32rpx;
      height: 32rpx;
      border: 3rpx solid var(--color-border);
      border-top-color: var(--color-primary);
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }
  }
  
  &-end {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    width: 100%;
    padding: 0 var(--spacing-lg);
  }
  
  &-line {
    flex: 1;
    height: 1rpx;
    background: var(--color-border);
  }
  
  &-text {
    font-size: var(--font-sm);
    color: var(--color-text-secondary);
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
