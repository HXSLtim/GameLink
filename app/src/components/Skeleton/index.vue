<template>
  <view class="skeleton" :class="{ 'skeleton-animate': animate }">
    <!-- 预设模板: 卡片模式 -->
    <template v-if="type === 'card'">
      <view class="skeleton-card">
        <view class="skeleton-avatar"></view>
        <view class="skeleton-content">
          <view class="skeleton-title" style="width: 40%"></view>
          <view class="skeleton-row" style="width: 90%"></view>
          <view class="skeleton-row" style="width: 60%"></view>
        </view>
      </view>
    </template>

    <!-- 预设模板: 列表模式 -->
    <template v-else-if="type === 'list'">
      <view v-for="i in rows" :key="i" class="skeleton-list-item">
        <view class="skeleton-avatar small"></view>
        <view class="skeleton-content">
          <view class="skeleton-title" style="width: 30%"></view>
          <view class="skeleton-row" style="width: 100%"></view>
        </view>
      </view>
    </template>

    <!-- 预设模板: 陪玩师列表 -->
    <template v-else-if="type === 'player'">
      <view class="skeleton-grid">
        <view v-for="i in 4" :key="i" class="skeleton-player-card">
          <view class="skeleton-cover"></view>
          <view class="skeleton-info">
            <view class="skeleton-title" style="width: 70%"></view>
            <view class="skeleton-row" style="width: 40%"></view>
            <view class="skeleton-row" style="width: 100%; height: 40rpx; margin-top: 16rpx"></view>
          </view>
        </view>
      </view>
    </template>

    <!-- 预设模板: 游戏图标 -->
    <template v-else-if="type === 'game'">
      <view class="skeleton-game-scroll">
        <view v-for="i in 5" :key="i" class="skeleton-game-item">
          <view class="skeleton-icon"></view>
          <view class="skeleton-text-xs"></view>
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
// 骨架屏动画 Mixin
@mixin skeleton-shimmer {
  background: #f2f3f5;
  background-image: linear-gradient(
    90deg,
    #f2f3f5 25%,
    #e6e8eb 37%,
    #f2f3f5 63%
  );
  background-size: 400% 100%;
  animation: skeleton-loading 1.4s ease infinite;
}

@keyframes skeleton-loading {
  0% { background-position: 100% 50%; }
  100% { background-position: 0 50%; }
}

.skeleton {
  width: 100%;
  box-sizing: border-box;

  &-animate {
    .skeleton-avatar,
    .skeleton-title,
    .skeleton-row,
    .skeleton-cover,
    .skeleton-icon,
    .skeleton-text-xs {
      @include skeleton-shimmer;
    }
  }
}

// 基础元素
.skeleton-avatar {
  width: 96rpx;
  height: 96rpx;
  border-radius: $gl-radius-circle;
  flex-shrink: 0;

  &.small {
    width: 80rpx;
    height: 80rpx;
  }
}

.skeleton-title {
  height: 32rpx;
  border-radius: $gl-radius-sm;
  margin-bottom: $gl-spacing-sm;
}

.skeleton-row {
  height: 24rpx;
  border-radius: $gl-radius-sm;
  margin-bottom: $gl-spacing-sm;

  &:last-child {
    margin-bottom: 0;
  }
}

.skeleton-cover {
  width: 100%;
  height: 320rpx;
  border-radius: $gl-radius-md $gl-radius-md 0 0;
}

.skeleton-icon {
  width: 88rpx;
  height: 88rpx;
  border-radius: $gl-radius-md;
  margin-bottom: $gl-spacing-xs;
}

.skeleton-text-xs {
  width: 60rpx;
  height: 20rpx;
  border-radius: $gl-radius-sm;
}

// 布局容器
.skeleton-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.skeleton-card {
  display: flex;
  gap: $gl-spacing-md;
  padding: $gl-spacing-md;
  background: #fff;
  border-radius: $gl-radius-md;
  margin-bottom: $gl-spacing-md;
}

.skeleton-list-item {
  display: flex;
  gap: $gl-spacing-md;
  padding: $gl-spacing-md 0;
  border-bottom: 1rpx solid rgba(0,0,0,0.05);
}

// 陪玩师网格
.skeleton-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: $gl-spacing-md;
  padding: $gl-spacing-md;
}

.skeleton-player-card {
  background: #fff;
  border-radius: $gl-radius-md;
  overflow: hidden;
  box-shadow: $gl-shadow-sm;

  .skeleton-info {
    padding: $gl-spacing-sm;
  }
}

// 游戏横滚
.skeleton-game-scroll {
  display: flex;
  gap: $gl-spacing-md;
  padding: $gl-spacing-md;
  overflow: hidden;
}

.skeleton-game-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex-shrink: 0;
}
</style>
