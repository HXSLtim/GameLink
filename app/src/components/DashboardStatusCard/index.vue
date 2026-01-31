<template>
  <view class="status-card" :class="status">
    <view class="status-left">
      <view class="avatar-wrap">
        <GlAvatar :src="avatar" :text="nickname" :size="80" :status="status === 'online' ? 'online' : status === 'busy' ? 'busy' : undefined" bordered />
      </view>
      <view class="status-info">
        <text class="player-name">{{ nickname }}</text>
        <text class="status-text">{{ statusText }}</text>
      </view>
    </view>
    <GlButton :type="status === 'online' ? 'default' : 'primary'" size="small" round @click="$emit('toggle')">
      {{ status === 'online' ? '下线' : '上线接单' }}
    </GlButton>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlAvatar from '@/components/gl/Avatar/index.vue'
import GlButton from '@/components/gl/Button/index.vue'

interface Props {
  status: 'online' | 'busy' | 'offline'
  nickname: string
  avatar?: string
}

const props = defineProps<Props>()

defineEmits<{
  toggle: []
}>()

const statusText = computed(() => {
  const texts = {
    online: '在线接单中',
    busy: '忙碌中',
    offline: '已下线',
  }
  return texts[props.status]
})
</script>

<style lang="scss" scoped>
.status-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 32rpx;
  margin: 24rpx;
  border-radius: 24rpx;
  background: linear-gradient(135deg, var(--color-bg-secondary) 0%, var(--color-bg-card) 100%);
  border: 2rpx solid var(--color-border);
  
  &.online {
    background: linear-gradient(135deg, rgba(0, 210, 106, 0.1) 0%, rgba(0, 210, 106, 0.05) 100%);
    border-color: var(--color-primary);
  }
  
  &.busy {
    background: linear-gradient(135deg, rgba(245, 158, 11, 0.1) 0%, rgba(245, 158, 11, 0.05) 100%);
    border-color: #F59E0B;
  }
}

.status-left {
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.avatar-wrap {
  flex-shrink: 0;
}

.status-info {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.player-name {
  font-size: 32rpx;
  font-weight: 700;
  color: var(--color-text);
}

.status-text {
  font-size: 26rpx;
  color: var(--color-text-secondary);
}
</style>
