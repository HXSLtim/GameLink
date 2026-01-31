<template>
  <view class="player-card">
    <GlAvatar :src="player.avatar" :text="player.nickname" :size="80" :status="player.isOnline ? 'online' : undefined" />
    <view class="player-info">
      <text class="player-name">{{ player.nickname }}</text>
      <view class="player-meta">
        <GlTag :type="player.isOnline ? 'success' : 'default'" size="mini">
          {{ player.isOnline ? '在线' : '离线' }}
        </GlTag>
        <text class="rating">⭐ {{ (player.rating || 5).toFixed(1) }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import GlAvatar from '@/components/gl/Avatar/index.vue'
import GlTag from '@/components/gl/Tag/index.vue'

export interface OrderPlayerInfo {
  id: number
  nickname: string
  avatar?: string
  isOnline?: boolean
  rating?: number
}

interface Props {
  player: OrderPlayerInfo
}

defineProps<Props>()
</script>

<style lang="scss" scoped>
.player-card {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 24rpx;
  margin: 24rpx;
  background: var(--color-bg-card);
  border-radius: 20rpx;
  border: 2rpx solid var(--color-border);
}

.player-info {
  flex: 1;
}

.player-name {
  display: block;
  font-size: 32rpx;
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: 8rpx;
}

.player-meta {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.rating {
  font-size: 26rpx;
  color: var(--color-text-secondary);
}
</style>
