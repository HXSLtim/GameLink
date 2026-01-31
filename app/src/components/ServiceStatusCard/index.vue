<template>
  <view class="service-status">
    <GlAvatar :src="avatar" text="👩‍💼" :size="80" />
    <view class="status-info">
      <text class="status-title">{{ title }}</text>
      <GlTag :type="isOnline ? 'success' : 'default'" size="small">
        {{ isOnline ? '在线' : '离线' }}
      </GlTag>
    </view>
  </view>
</template>

<script setup lang="ts">
import GlAvatar from '@/components/gl/Avatar/index.vue'
import GlTag from '@/components/gl/Tag/index.vue'

interface Props {
  title?: string
  isOnline?: boolean
  avatar?: string
}

withDefaults(defineProps<Props>(), {
  title: 'GameLink 官方客服',
  isOnline: true,
})
</script>

<style lang="scss" scoped>
.service-status {
  display: flex;
  align-items: center;
  gap: 16rpx;
  padding: 16rpx 24rpx;
  background: var(--color-bg-card);
  border-bottom: 1rpx solid var(--color-border);
}

.status-info {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.status-title {
  font-size: 30rpx;
  font-weight: 500;
  color: var(--color-text);
}
</style>
