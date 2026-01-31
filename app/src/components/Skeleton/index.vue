<template>
  <view class="skeleton" :class="{ 'skeleton-animate': animate }">
    <!-- 预设模板 -->
    <template v-if="type === 'card'">
      <view class="skeleton-card">
        <view class="skeleton-avatar"></view>
        <view class="skeleton-content">
          <view class="skeleton-title"></view>
          <view class="skeleton-desc"></view>
          <view class="skeleton-desc short"></view>
        </view>
      </view>
    </template>
    
    <template v-else-if="type === 'list'">
      <view v-for="i in rows" :key="i" class="skeleton-list-item">
        <view class="skeleton-avatar small"></view>
        <view class="skeleton-content">
          <view class="skeleton-title"></view>
          <view class="skeleton-desc"></view>
        </view>
      </view>
    </template>
    
    <template v-else-if="type === 'player'">
      <view class="skeleton-player">
        <view class="skeleton-player-cover"></view>
        <view class="skeleton-player-info">
          <view class="skeleton-avatar"></view>
          <view class="skeleton-title"></view>
          <view class="skeleton-tags">
            <view class="skeleton-tag"></view>
            <view class="skeleton-tag"></view>
          </view>
          <view class="skeleton-price"></view>
        </view>
      </view>
    </template>
    
    <!-- 自定义布局 -->
    <template v-else>
      <slot>
        <view v-for="i in rows" :key="i" class="skeleton-row" :style="{ width: getRowWidth(i) }"></view>
      </slot>
    </template>
  </view>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  type?: 'default' | 'card' | 'list' | 'player'
  rows?: number
  animate?: boolean
  rowWidths?: string[]
}>(), {
  type: 'default',
  rows: 3,
  animate: true,
  rowWidths: () => ['100%', '80%', '60%'],
})

const getRowWidth = (index: number) => {
  const widths = props.rowWidths
  return widths[(index - 1) % widths.length] || '100%'
}
</script>

<style lang="scss" scoped>
.skeleton {
  &-animate {
    .skeleton-avatar,
    .skeleton-title,
    .skeleton-desc,
    .skeleton-row,
    .skeleton-tag,
    .skeleton-price,
    .skeleton-player-cover {
      background: linear-gradient(
        90deg,
        var(--bg-secondary, #F0F0F0) 25%,
        var(--bg-card, #E8E8E8) 50%,
        var(--bg-secondary, #F0F0F0) 75%
      );
      background-size: 200% 100%;
      animation: shimmer 1.5s infinite;
    }
  }
}

.skeleton-avatar {
  width: 96rpx;
  height: 96rpx;
  border-radius: 50%;
  background: var(--bg-secondary, #F0F0F0);
  flex-shrink: 0;
  
  &.small {
    width: 80rpx;
    height: 80rpx;
  }
}

.skeleton-title {
  height: 32rpx;
  background: var(--bg-secondary, #F0F0F0);
  border-radius: 8rpx;
  width: 60%;
}

.skeleton-desc {
  height: 24rpx;
  background: var(--bg-secondary, #F0F0F0);
  border-radius: 6rpx;
  width: 100%;
  margin-top: 16rpx;
  
  &.short {
    width: 40%;
  }
}

.skeleton-row {
  height: 28rpx;
  background: var(--bg-secondary, #F0F0F0);
  border-radius: 8rpx;
  margin-bottom: 16rpx;
  
  &:last-child {
    margin-bottom: 0;
  }
}

.skeleton-card {
  display: flex;
  gap: 24rpx;
  padding: 24rpx;
  background: var(--bg-card, #FFFFFF);
  border-radius: 16rpx;
}

.skeleton-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.skeleton-list-item {
  display: flex;
  gap: 20rpx;
  padding: 20rpx 0;
  border-bottom: 1rpx solid var(--border, #E5E5E5);
  
  &:last-child {
    border-bottom: none;
  }
}

.skeleton-player {
  background: var(--bg-card, #FFFFFF);
  border-radius: 16rpx;
  overflow: hidden;
  
  &-cover {
    width: 100%;
    height: 300rpx;
    background: var(--bg-secondary, #F0F0F0);
  }
  
  &-info {
    padding: 24rpx;
  }
}

.skeleton-tags {
  display: flex;
  gap: 16rpx;
  margin-top: 16rpx;
}

.skeleton-tag {
  width: 100rpx;
  height: 40rpx;
  background: var(--bg-secondary, #F0F0F0);
  border-radius: 8rpx;
}

.skeleton-price {
  width: 120rpx;
  height: 36rpx;
  background: var(--bg-secondary, #F0F0F0);
  border-radius: 8rpx;
  margin-top: 16rpx;
}

@keyframes shimmer {
  0% {
    background-position: -200% 0;
  }
  100% {
    background-position: 200% 0;
  }
}
</style>
