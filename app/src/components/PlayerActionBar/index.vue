<template>
  <view class="action-bar">
    <view class="action-left">
      <view class="action-icon-btn" @tap="$emit('favorite')">
        <uv-icon :name="isFavorite ? 'heart-fill' : 'heart'" size="22" :color="isFavorite ? 'var(--color-error)' : 'var(--color-text-secondary)'" class="action-icon" />
        <text class="action-text">收藏</text>
      </view>
      <view class="action-icon-btn" @tap="$emit('chat')">
        <uv-icon name="chat" size="22" color="var(--color-text-secondary)" class="action-icon" />
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
  gap: var(--spacing-sm);
  padding: var(--spacing-sm) var(--spacing-md);
  padding-bottom: calc(var(--spacing-sm) + env(safe-area-inset-bottom));
  background: var(--color-bg-card);
  border-top: 1rpx solid var(--color-border);
}

.action-left {
  display: flex;
  gap: var(--spacing-sm);
}

.action-icon-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-xs);
  padding: var(--spacing-xs) var(--spacing-sm);
  background: var(--color-bg-secondary);
  border: 1rpx solid var(--color-border);
  border-radius: var(--radius-md);
  transition: all 0.2s;
  cursor: pointer;
  @include press-effect;
}

.action-icon {
  flex-shrink: 0;
}

.action-text {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}

.action-right {
  flex: 1;
}
</style>
