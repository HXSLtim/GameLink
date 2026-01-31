<template>
  <view class="status-card" :class="statusClass">
    <view class="status-icon">{{ statusIcon }}</view>
    <view class="status-info">
      <text class="status-text">{{ statusText }}</text>
      <text class="status-desc">{{ statusDesc }}</text>
    </view>
    <view v-if="status === 'pending' && countdown > 0" class="countdown">
      <text>剩余支付时间：{{ formattedCountdown }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  status: string
  countdown?: number
}

const props = withDefaults(defineProps<Props>(), {
  countdown: 0,
})

const statusConfig: Record<string, { icon: string; text: string; desc: string; class: string }> = {
  pending: { icon: '⏳', text: '待支付', desc: '请在规定时间内完成支付', class: 'pending' },
  confirmed: { icon: '✅', text: '已支付', desc: '等待陪玩师接单', class: 'confirmed' },
  in_progress: { icon: '🎮', text: '服务中', desc: '陪玩师正在为您服务', class: 'in-progress' },
  completed: { icon: '🎉', text: '已完成', desc: '服务已完成，欢迎评价', class: 'completed' },
  cancelled: { icon: '❌', text: '已取消', desc: '订单已取消', class: 'cancelled' },
  refunding: { icon: '🔄', text: '退款中', desc: '正在处理退款', class: 'refunding' },
  refunded: { icon: '💰', text: '已退款', desc: '退款已到账', class: 'refunded' },
}

const config = computed(() => statusConfig[props.status] || statusConfig.pending)
const statusClass = computed(() => config.value.class)
const statusIcon = computed(() => config.value.icon)
const statusText = computed(() => config.value.text)
const statusDesc = computed(() => config.value.desc)

const formattedCountdown = computed(() => {
  const minutes = Math.floor(props.countdown / 60)
  const seconds = props.countdown % 60
  return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
})
</script>

<style lang="scss" scoped>
.status-card {
  display: flex;
  align-items: center;
  gap: 24rpx;
  padding: 40rpx 32rpx;
  border-radius: 0 0 32rpx 32rpx;
  background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-primary-light, #4ADE80) 100%);
  color: #FFFFFF;
  
  &.pending {
    background: linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%);
  }
  
  &.cancelled, &.refunded {
    background: linear-gradient(135deg, #6B7280 0%, #9CA3AF 100%);
  }
  
  &.refunding {
    background: linear-gradient(135deg, #3B82F6 0%, #60A5FA 100%);
  }
}

.status-icon {
  font-size: 56rpx;
}

.status-info {
  flex: 1;
}

.status-text {
  display: block;
  font-size: 36rpx;
  font-weight: 700;
  margin-bottom: 8rpx;
}

.status-desc {
  font-size: 26rpx;
  opacity: 0.9;
}

.countdown {
  padding: 12rpx 24rpx;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 24rpx;
  font-size: 24rpx;
}
</style>
