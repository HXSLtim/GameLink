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

    <template v-else-if="type === 'game'">
      <view class="skeleton-game">
        <view v-for="i in rows" :key="i" class="skeleton-game-item">
          <view class="skeleton-game-icon"></view>
          <view class="skeleton-game-name"></view>
          <view class="skeleton-game-count"></view>
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
  type?: 'default' | 'card' | 'list' | 'player' | 'game'
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
    .skeleton-player-cover,
    .skeleton-game-icon,
    .skeleton-game-name,
    .skeleton-game-count {
      background: var(--color-bg-secondary);
      animation: shimmer 1.5s ease-in-out infinite;
    }
  }
}

.skeleton-avatar {
  width: 96rpx;
  height: 96rpx;
  border-radius: var(--radius-full);
  background: var(--color-bg-secondary);
  flex-shrink: 0;
  
  &.small {
    width: 80rpx;
    height: 80rpx;
  }
}

.skeleton-title {
  height: 32rpx;
  background: var(--color-bg-secondary);
  border-radius: var(--radius-sm);
  width: 60%;
}

.skeleton-desc {
  height: 24rpx;
  background: var(--color-bg-secondary);
  border-radius: var(--radius-sm);
  width: 100%;
  margin-top: var(--spacing-sm);
  
  &.short {
    width: 40%;
  }
}

.skeleton-row {
  height: var(--font-md);
  background: var(--color-bg-secondary);
  border-radius: var(--radius-sm);
  margin-bottom: var(--spacing-sm);
  
  &:last-child {
    margin-bottom: 0;
  }
}

.skeleton-card {
  display: flex;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
}

.skeleton-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.skeleton-list-item {
  display: flex;
  gap: var(--spacing-md);
  padding: var(--spacing-sm) 0;
  border-bottom: 1rpx solid var(--color-border);
  
  &:last-child {
    border-bottom: none;
  }
}

.skeleton-player {
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  overflow: hidden;
  
  &-cover {
    width: 100%;
    height: 300rpx;
    background: var(--color-bg-secondary);
  }
  
  &-info {
    padding: var(--spacing-md);
  }
}

.skeleton-game {
  display: flex;
  gap: var(--spacing-md);
  padding-right: var(--spacing-md);
}

.skeleton-game-item {
  width: 160rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-xs);
  padding: var(--spacing-md);
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1rpx solid var(--color-border);
}

.skeleton-game-icon {
  width: 80rpx;
  height: 80rpx;
  border-radius: var(--radius-sm);
  background: var(--color-bg-secondary);
}

.skeleton-game-name {
  width: 90rpx;
  height: 24rpx;
  border-radius: var(--radius-sm);
  background: var(--color-bg-secondary);
}

.skeleton-game-count {
  width: 70rpx;
  height: 20rpx;
  border-radius: var(--radius-sm);
  background: var(--color-bg-secondary);
}

.skeleton-tags {
  display: flex;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-sm);
}

.skeleton-tag {
  width: 100rpx;
  height: 40rpx;
  background: var(--color-bg-secondary);
  border-radius: var(--radius-sm);
}

.skeleton-price {
  width: 120rpx;
  height: 36rpx;
  background: var(--color-bg-secondary);
  border-radius: var(--radius-sm);
  margin-top: var(--spacing-sm);
}

@keyframes shimmer {
  0% {
    background-color: var(--color-bg-secondary);
  }
  50% {
    background-color: var(--color-bg-card);
  }
  100% {
    background-color: var(--color-bg-secondary);
  }
}
</style>
