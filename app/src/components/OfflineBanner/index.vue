<template>
  <view v-if="visible" class="offline-banner" :class="[`offline-banner--${type}`]">
    <uv-icon :name="icon" size="14" :color="iconColor"></uv-icon>
    <text class="offline-text">{{ message }}</text>
    <view v-if="showAction" class="action-btn" @tap="$emit('action')">
      <text>{{ actionText }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  visible: boolean
  type?: 'warning' | 'error' | 'info'
  message?: string
  icon?: string
  showAction?: boolean
  actionText?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'warning',
  message: '网络不可用',
  icon: 'wifi-off',
  showAction: true,
  actionText: '刷新',
})

defineEmits<{
  action: []
}>()

const iconColor = computed(() => {
  const colors = {
    warning: '#F59E0B',
    error: '#EF4444',
    info: '#3B82F6',
  }
  return colors[props.type]
})
</script>

<style lang="scss" scoped>
.offline-banner {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12rpx;
  padding: 12rpx 20rpx;
  
  &--warning {
    background: linear-gradient(135deg, #FEF3C7 0%, #FDE68A 100%);
    border-bottom: 1rpx solid #F59E0B;
    
    .offline-text { color: #92400E; }
    .action-btn { background: linear-gradient(135deg, #F59E0B 0%, #D97706 100%); }
  }
  
  &--error {
    background: linear-gradient(135deg, #FEE2E2 0%, #FECACA 100%);
    border-bottom: 1rpx solid #EF4444;
    
    .offline-text { color: #991B1B; }
    .action-btn { background: linear-gradient(135deg, #EF4444 0%, #DC2626 100%); }
  }
  
  &--info {
    background: linear-gradient(135deg, #DBEAFE 0%, #BFDBFE 100%);
    border-bottom: 1rpx solid #3B82F6;
    
    .offline-text { color: #1E40AF; }
    .action-btn { background: linear-gradient(135deg, #3B82F6 0%, #2563EB 100%); }
  }
}

.offline-text {
  font-size: 24rpx;
  font-weight: 500;
}

.action-btn {
  margin-left: 12rpx;
  padding: 8rpx 16rpx;
  border-radius: 16rpx;
  
  text {
    font-size: 22rpx;
    color: #fff;
    font-weight: 600;
  }
  
  &:active {
    transform: scale(0.95);
  }
}
</style>
