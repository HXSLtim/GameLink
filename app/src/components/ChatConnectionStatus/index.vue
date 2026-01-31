<template>
  <view v-if="status !== 'connected'" class="connection-status" :class="status">
    <uv-icon 
      :name="status === 'connecting' ? 'reload' : 'wifi-off'" 
      size="14" 
      color="#fff"
    ></uv-icon>
    <text>{{ statusText }}</text>
    <view v-if="status === 'disconnected'" class="reconnect-btn" @tap="$emit('reconnect')">
      重新连接
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type WsStatus = 'connecting' | 'connected' | 'disconnected'

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
  gap: 12rpx;
  padding: 16rpx 24rpx;
  font-size: 24rpx;
  color: #fff;
  
  &.connecting {
    background: #F59E0B;
  }
  
  &.disconnected {
    background: #EF4444;
  }
}

.reconnect-btn {
  margin-left: 16rpx;
  padding: 8rpx 16rpx;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 8rpx;
  font-size: 22rpx;
}
</style>
