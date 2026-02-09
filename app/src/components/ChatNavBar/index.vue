<template>
  <view class="nav-bar">
    <view class="nav-back" @tap="$emit('back')">
      <uv-icon name="arrow-left" size="20" color="var(--color-text)"></uv-icon>
    </view>
    <view class="nav-center">
      <view class="nav-title-row">
        <text class="nav-title">{{ name }}</text>
        <!-- 私聊在线状态 -->
        <view v-if="type === 'private'" class="online-dot" :class="{ online: isOnline }" />
      </view>
      <text v-if="type === 'private'" class="nav-subtitle">
        {{ isOnline ? '在线' : '离线' }}
      </text>
      <text v-else-if="memberCount" class="nav-subtitle">
        {{ memberCount }}人
      </text>
    </view>
    <view class="nav-actions">
      <view class="action-btn" @tap="$emit('menu')">
        <uv-icon name="more-dot-fill" size="20" color="var(--color-text)"></uv-icon>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import type { ChatInfo } from '@/types/message'

type Props = Pick<ChatInfo, 'name'> & Partial<Pick<ChatInfo, 'type' | 'isOnline' | 'memberCount'>>

withDefaults(defineProps<Props>(), {
  type: 'private',
  isOnline: false,
})

defineEmits<{
  back: []
  menu: []
}>()
</script>

<style lang="scss" scoped>
.nav-bar {
  display: flex;
  align-items: center;
  padding: var(--spacing-sm) var(--spacing-md);
  padding-top: calc(var(--spacing-sm) + env(safe-area-inset-top));
  background: var(--color-bg-card);
  border-bottom: 1rpx solid var(--color-border);
  min-height: 88rpx;
}

.nav-back {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background 0.2s;

  &:hover {
    background: var(--color-bg-secondary);
  }

  &:active {
    transform: scale(0.92);
  }
}

.nav-center {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 0;
}

.nav-title-row {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.nav-title {
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 400rpx;
}

.online-dot {
  width: 12rpx;
  height: 12rpx;
  border-radius: 50%;
  background: var(--color-text-placeholder);
  flex-shrink: 0;
  transition: background 0.3s;

  &.online {
    background: #34C759;
    box-shadow: 0 0 6rpx rgba(52, 199, 89, 0.5);
  }
}

.nav-subtitle {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
  margin-top: 2rpx;
}

.nav-actions {
  width: 64rpx;
  display: flex;
  justify-content: flex-end;
}

.action-btn {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background 0.2s;

  &:hover {
    background: var(--color-bg-secondary);
  }

  &:active {
    transform: scale(0.92);
  }
}
</style>
