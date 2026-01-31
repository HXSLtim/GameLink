<template>
  <view class="action-bar">
    <view class="action-left">
      <view class="action-icon-btn" @tap="$emit('favorite')">
        <text class="action-icon">{{ isFavorite ? '❤️' : '🤍' }}</text>
        <text class="action-text">收藏</text>
      </view>
      <view class="action-icon-btn" @tap="$emit('chat')">
        <text class="action-icon">💬</text>
        <text class="action-text">聊天</text>
      </view>
    </view>
    <view class="action-right">
      <GlButton 
        :type="isOnline ? 'primary' : 'default'"
        :disabled="!isOnline"
        size="large"
        round
        block
        @click="$emit('order')"
      >
        {{ isOnline ? '立即下单' : '陪玩师离线' }}
      </GlButton>
    </view>
  </view>
</template>

<script setup lang="ts">
import GlButton from '@/components/gl/Button/index.vue'

interface Props {
  isFavorite: boolean
  isOnline: boolean
}

defineProps<Props>()

defineEmits<{
  favorite: []
  chat: []
  order: []
}>()
</script>

<style lang="scss" scoped>
.action-bar {
  display: flex;
  align-items: center;
  gap: 24rpx;
  padding: 20rpx 32rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background: var(--color-bg-card);
  border-top: 1rpx solid var(--color-border);
  box-shadow: 0 -4rpx 20rpx rgba(0, 0, 0, 0.05);
}

.action-left {
  display: flex;
  gap: 32rpx;
}

.action-icon-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4rpx;
  padding: 8rpx 16rpx;
  transition: all 0.2s;
  
  &:active {
    transform: scale(0.9);
  }
}

.action-icon {
  font-size: 40rpx;
}

.action-text {
  font-size: 22rpx;
  color: var(--color-text-secondary);
}

.action-right {
  flex: 1;
}
</style>
