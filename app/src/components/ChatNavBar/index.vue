<template>
  <view class="nav-bar">
    <view class="nav-back" @tap="$emit('back')">
      <uv-icon name="arrow-left" size="20" color="var(--color-text)"></uv-icon>
    </view>
    <view class="nav-center">
      <text class="nav-title">{{ name }}</text>
      <text v-if="type === 'private'" class="online-status" :class="{ online: isOnline }">
        {{ isOnline ? '在线' : '离线' }}
      </text>
      <text v-else-if="memberCount" class="member-count">
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
interface Props {
  name: string
  type?: 'private' | 'order' | 'public'
  isOnline?: boolean
  memberCount?: number
}

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
  padding: 16rpx 24rpx;
  padding-top: calc(16rpx + env(safe-area-inset-top));
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
}

.nav-center {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 0;
}

.nav-title {
  font-size: 34rpx;
  font-weight: 700;
  color: var(--color-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 400rpx;
}

.online-status {
  font-size: 22rpx;
  color: var(--color-text-secondary);
  margin-top: 4rpx;
  
  &.online {
    color: var(--color-primary);
  }
}

.member-count {
  font-size: 22rpx;
  color: var(--color-text-secondary);
  margin-top: 4rpx;
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
}
</style>
