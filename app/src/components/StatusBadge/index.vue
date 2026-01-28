<template>
  <view class="status-badge" :class="[`status-${type}`, `size-${size}`]">
    <view v-if="dot" class="status-dot"></view>
    <text class="status-text">{{ text || statusText }}</text>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  status?: string
  text?: string
  type?: 'success' | 'warning' | 'error' | 'info' | 'default'
  size?: 'small' | 'medium' | 'large'
  dot?: boolean
}>(), {
  type: 'default',
  size: 'medium',
  dot: false,
})

// 状态文本映射
const statusTextMap: Record<string, string> = {
  // 订单状态
  pending: '待支付',
  confirmed: '待接单',
  in_progress: '进行中',
  completed: '已完成',
  canceled: '已取消',
  refunded: '已退款',
  disputed: '争议中',
  
  // 认证状态
  verified: '已认证',
  rejected: '已拒绝',
  revoked: '已撤销',
  expired: '已过期',
  
  // 在线状态
  online: '在线',
  busy: '忙碌',
  offline: '离线',
  
  // 通用
  active: '正常',
  inactive: '停用',
}

// 状态类型映射
const statusTypeMap: Record<string, string> = {
  // 订单状态
  pending: 'warning',
  confirmed: 'info',
  in_progress: 'success',
  completed: 'default',
  canceled: 'default',
  refunded: 'default',
  disputed: 'error',
  
  // 认证状态
  verified: 'success',
  rejected: 'error',
  revoked: 'warning',
  expired: 'default',
  
  // 在线状态
  online: 'success',
  busy: 'warning',
  offline: 'default',
  
  // 通用
  active: 'success',
  inactive: 'default',
}

const statusText = computed(() => {
  if (props.text) return props.text
  return statusTextMap[props.status || ''] || props.status || ''
})

const computedType = computed(() => {
  if (props.type !== 'default') return props.type
  return statusTypeMap[props.status || ''] || 'default'
})
</script>

<style lang="scss" scoped>
.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 4rpx 16rpx;
  border-radius: 8rpx;
  
  .status-dot {
    width: 12rpx;
    height: 12rpx;
    border-radius: 50%;
    margin-right: 8rpx;
  }
  
  .status-text {
    font-weight: 500;
  }
}

// 尺寸
.size-small {
  padding: 2rpx 12rpx;
  .status-text { font-size: 20rpx; }
  .status-dot { width: 8rpx; height: 8rpx; }
}

.size-medium {
  padding: 4rpx 16rpx;
  .status-text { font-size: 24rpx; }
  .status-dot { width: 12rpx; height: 12rpx; }
}

.size-large {
  padding: 8rpx 20rpx;
  .status-text { font-size: 28rpx; }
  .status-dot { width: 14rpx; height: 14rpx; }
}

// 类型
.status-default {
  background: rgba(156, 163, 175, 0.1);
  .status-text { color: #6B7280; }
  .status-dot { background: #6B7280; }
}

.status-success {
  background: rgba(34, 197, 94, 0.1);
  .status-text { color: #22C55E; }
  .status-dot { background: #22C55E; }
}

.status-warning {
  background: rgba(245, 158, 11, 0.1);
  .status-text { color: #F59E0B; }
  .status-dot { background: #F59E0B; }
}

.status-error {
  background: rgba(239, 68, 68, 0.1);
  .status-text { color: #EF4444; }
  .status-dot { background: #EF4444; }
}

.status-info {
  background: rgba(59, 130, 246, 0.1);
  .status-text { color: #3B82F6; }
  .status-dot { background: #3B82F6; }
}
</style>
