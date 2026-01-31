<template>
  <view class="order-card" @tap="$emit('click')">
    <!-- 订单头部 -->
    <view class="order-header">
      <text class="order-no">订单号：{{ order.orderNo }}</text>
      <view class="order-status" :class="`order-status--${statusInfo.class}`">
        {{ statusInfo.text }}
      </view>
    </view>
    
    <!-- 陪玩师信息 -->
    <view class="order-player">
      <GlAvatar
        :src="order.player.avatar"
        :text="order.player.nickname"
        size="medium"
        status="online"
      />
      <view class="player-info">
        <text class="player-name">{{ order.player.nickname }}</text>
        <text class="service-name">{{ order.serviceName }}</text>
      </view>
    </view>
    
    <!-- 订单信息 -->
    <view class="order-info">
      <view class="info-item">
        <text class="info-label">游戏</text>
        <text class="info-value">{{ order.gameName }}</text>
      </view>
      <view class="info-item">
        <text class="info-label">数量</text>
        <text class="info-value">{{ order.quantity }}{{ order.unit }}</text>
      </view>
      <view class="info-item">
        <text class="info-label">金额</text>
        <text class="info-value info-value--price">¥{{ order.totalAmount.toFixed(2) }}</text>
      </view>
    </view>
    
    <!-- 订单时间 -->
    <view class="order-time">
      <text>{{ formattedTime }}</text>
    </view>
    
    <!-- 操作按钮 -->
    <view class="order-actions">
      <template v-for="action in availableActions" :key="action.key">
        <GlButton
          :type="action.type"
          :plain="action.plain"
          size="small"
          round
          @click.stop="handleAction(action.key)"
        >
          {{ action.label }}
        </GlButton>
      </template>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlAvatar from '@/components/gl/Avatar/index.vue'
import GlButton from '@/components/gl/Button/index.vue'

export type OrderStatus = 'pending' | 'confirmed' | 'in_progress' | 'completed' | 'canceled' | 'refunded' | 'disputed'

export interface OrderPlayer {
  id: number
  nickname: string
  avatar?: string
}

export interface Order {
  id: number
  orderNo: string
  status: OrderStatus
  player: OrderPlayer
  gameName: string
  serviceName: string
  quantity: number
  unit: string
  totalAmount: number
  createdAt: string
  reviewed?: boolean
}

interface Props {
  order: Order
}

const props = defineProps<Props>()

const emit = defineEmits<{
  click: []
  action: [key: string, order: Order]
}>()

// 状态映射
const statusMap: Record<OrderStatus, { text: string; class: string }> = {
  pending: { text: '待支付', class: 'warning' },
  confirmed: { text: '待服务', class: 'info' },
  in_progress: { text: '进行中', class: 'primary' },
  completed: { text: '已完成', class: 'success' },
  canceled: { text: '已取消', class: 'default' },
  refunded: { text: '已退款', class: 'default' },
  disputed: { text: '争议中', class: 'error' },
}

const statusInfo = computed(() => statusMap[props.order.status] || { text: props.order.status, class: 'default' })

// 格式化时间
const formattedTime = computed(() => {
  const date = new Date(props.order.createdAt)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hour = String(date.getHours()).padStart(2, '0')
  const minute = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day} ${hour}:${minute}`
})

// 可用操作按钮
interface ActionButton {
  key: string
  label: string
  type: 'primary' | 'default'
  plain: boolean
}

const availableActions = computed((): ActionButton[] => {
  const { status, reviewed } = props.order
  const actions: ActionButton[] = []
  
  switch (status) {
    case 'pending':
      actions.push({ key: 'cancel', label: '取消订单', type: 'default', plain: true })
      actions.push({ key: 'pay', label: '去支付', type: 'primary', plain: false })
      break
    case 'confirmed':
    case 'in_progress':
      actions.push({ key: 'contact', label: '联系陪玩', type: 'default', plain: true })
      if (status === 'in_progress') {
        actions.push({ key: 'complete', label: '确认完成', type: 'primary', plain: false })
      }
      break
    case 'completed':
      if (!reviewed) {
        actions.push({ key: 'review', label: '去评价', type: 'primary', plain: false })
      } else {
        actions.push({ key: 'reorder', label: '再来一单', type: 'default', plain: true })
      }
      break
    case 'canceled':
    case 'refunded':
      actions.push({ key: 'reorder', label: '再来一单', type: 'default', plain: true })
      break
    case 'disputed':
      actions.push({ key: 'viewDispute', label: '查看进度', type: 'default', plain: true })
      break
  }
  
  return actions
})

const handleAction = (key: string) => {
  emit('action', key, props.order)
}
</script>

<style lang="scss" scoped>
.order-card {
  background: var(--color-bg-card);
  border-radius: 24rpx;
  padding: 28rpx;
  margin-bottom: 20rpx;
  border: 2rpx solid var(--color-border);
  transition: all 0.2s;
  
  &:active {
    transform: scale(0.99);
    border-color: var(--color-primary);
  }
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 16rpx;
  border-bottom: 1rpx solid var(--color-border);
  margin-bottom: 20rpx;
}

.order-no {
  font-size: 24rpx;
  color: var(--color-text-secondary);
  font-weight: 500;
}

.order-status {
  font-size: 26rpx;
  font-weight: 600;
  padding: 6rpx 16rpx;
  border-radius: 12rpx;
  
  &--warning { color: var(--color-warning); background: rgba(245, 158, 11, 0.1); }
  &--info { color: var(--color-info, #3B82F6); background: rgba(59, 130, 246, 0.1); }
  &--primary { color: var(--color-primary); background: rgba(0, 210, 106, 0.1); }
  &--success { color: var(--color-success); background: rgba(34, 197, 94, 0.1); }
  &--default { color: var(--color-text-placeholder); background: rgba(156, 163, 175, 0.1); }
  &--error { color: var(--color-error); background: rgba(239, 68, 68, 0.1); }
}

.order-player {
  display: flex;
  align-items: center;
  gap: 20rpx;
  margin-bottom: 20rpx;
}

.player-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.player-name {
  font-size: 32rpx;
  font-weight: 600;
  color: var(--color-text);
}

.service-name {
  font-size: 26rpx;
  color: var(--color-text-secondary);
}

.order-info {
  display: flex;
  gap: 32rpx;
  padding: 20rpx;
  background: var(--color-bg-secondary);
  border-radius: 16rpx;
  margin-bottom: 20rpx;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.info-label {
  font-size: 22rpx;
  color: var(--color-text-placeholder);
}

.info-value {
  font-size: 28rpx;
  font-weight: 500;
  color: var(--color-text);
  
  &--price {
    font-weight: 700;
    color: var(--color-primary);
    font-size: 32rpx;
  }
}

.order-time {
  margin-bottom: 16rpx;
  
  text {
    font-size: 24rpx;
    color: var(--color-text-placeholder);
  }
}

.order-actions {
  display: flex;
  justify-content: flex-end;
  gap: 20rpx;
  padding-top: 20rpx;
  border-top: 1rpx solid var(--color-border);
}
</style>
