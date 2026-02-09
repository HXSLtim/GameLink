<template>
  <view v-if="status !== 'connected'" class="connection-status" :class="status">
    <view class="status-dot" :class="{ spin: status === 'connecting' }" />
    <text class="status-text">{{ statusText }}</text>
    <view v-if="status === 'disconnected'" class="reconnect-btn" @tap="$emit('reconnect')">
      重新连接
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { WsStatus } from '@/types/message'

interface Props {
  status: WsStatus
}

const props = defineProps<Props>()

defineEmits<{
  reconnect: []
}>()

const statusText = computed(() => {
  const texts: Record<WsStatus, string> = {
    connecting: '正在连接...',
    connected: '已连接',
    disconnected: '连接断开',
  }
  return texts[props.status]
})
</script>

<style lang="scss" scoped>
.connection-status {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
  padding: 10rpx var(--spacing-md);
  font-size: var(--font-xs);
  color: #fff;
  
  &.connecting {
    background: var(--color-warning);
  }
  
  &.disconnected {
    background: var(--color-error);
  }
}

.status-dot {
  width: 12rpx;
  height: 12rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.8);

  &.spin {
    animation: blink 1s ease-in-out infinite;
  }
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}

.status-text {
  font-weight: 500;
}

.reconnect-btn {
  margin-left: var(--spacing-xs);
  padding: 4rpx var(--spacing-sm);
  background: rgba(255, 255, 255, 0.2);
  border-radius: 100rpx;
  font-size: var(--font-xs);
  cursor: pointer;
  transition: background 0.2s;

  &:hover {
    background: rgba(255, 255, 255, 0.35);
  }
  
  &:active {
    transform: scale(0.95);
  }
}
</style>
