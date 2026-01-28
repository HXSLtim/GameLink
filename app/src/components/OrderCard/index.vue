<template>
  <view class="order-card" @click="$emit('click', order)">
    <!-- 头部：订单号和状态 -->
    <view class="order-header">
      <text class="order-no">订单号：{{ order.orderNo }}</text>
      <StatusBadge :status="order.status" size="small" />
    </view>

    <!-- 内容 -->
    <view class="order-content">
      <!-- 陪玩师/用户信息 -->
      <view class="order-user">
        <view class="user-avatar">
          <image v-if="targetUser?.avatar" :src="targetUser.avatar" mode="aspectFill" />
          <text v-else class="avatar-placeholder">{{ targetUser?.nickname?.[0] || '?' }}</text>
        </view>
        <view class="user-info">
          <text class="user-name">{{ targetUser?.nickname || '未知用户' }}</text>
          <text class="service-name">{{ order.serviceName || '游戏陪玩' }}</text>
        </view>
      </view>

      <!-- 订单详情 -->
      <view class="order-details">
        <view class="detail-item">
          <text class="detail-label">游戏</text>
          <text class="detail-value">{{ order.gameName || '-' }}</text>
        </view>
        <view class="detail-item">
          <text class="detail-label">数量</text>
          <text class="detail-value">{{ order.quantity }}{{ order.unit || '局' }}</text>
        </view>
        <view class="detail-item">
          <text class="detail-label">金额</text>
          <text class="detail-value price">¥{{ formatPrice(order.totalAmount) }}</text>
        </view>
      </view>
    </view>

    <!-- 底部：时间和操作 -->
    <view class="order-footer">
      <text class="order-time">{{ formatTime(order.createdAt) }}</text>
      <view class="order-actions">
        <slot name="actions">
          <!-- 默认操作按钮 -->
          <view v-if="showDefaultActions" class="action-buttons">
            <button 
              v-if="canCancel" 
              class="btn-action btn-cancel"
              @click.stop="$emit('cancel', order)"
            >
              取消订单
            </button>
            <button 
              v-if="canPay" 
              class="btn-action btn-primary"
              @click.stop="$emit('pay', order)"
            >
              去支付
            </button>
            <button 
              v-if="canReview" 
              class="btn-action btn-primary"
              @click.stop="$emit('review', order)"
            >
              去评价
            </button>
            <button 
              v-if="canAccept" 
              class="btn-action btn-primary"
              @click.stop="$emit('accept', order)"
            >
              接单
            </button>
            <button 
              v-if="canComplete" 
              class="btn-action btn-primary"
              @click.stop="$emit('complete', order)"
            >
              完成订单
            </button>
          </view>
        </slot>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import StatusBadge from '@/components/StatusBadge/index.vue'
import dayjs from 'dayjs'

export interface OrderUser {
  id: number
  nickname: string
  avatar?: string
}

export interface Order {
  id: number
  orderNo: string
  status: string
  gameName?: string
  serviceName?: string
  quantity: number
  unit?: string
  totalAmount: number
  createdAt: string
  user?: OrderUser
  player?: OrderUser
}

const props = withDefaults(defineProps<{
  order: Order
  viewAs?: 'user' | 'player'  // 以用户视角还是陪玩师视角
  showDefaultActions?: boolean
}>(), {
  viewAs: 'user',
  showDefaultActions: true,
})

defineEmits<{
  click: [order: Order]
  cancel: [order: Order]
  pay: [order: Order]
  review: [order: Order]
  accept: [order: Order]
  complete: [order: Order]
}>()

// 目标用户（用户视角看陪玩师，陪玩师视角看用户）
const targetUser = computed(() => {
  return props.viewAs === 'user' ? props.order.player : props.order.user
})

// 操作按钮显示逻辑
const canCancel = computed(() => {
  return props.order.status === 'pending' && props.viewAs === 'user'
})

const canPay = computed(() => {
  return props.order.status === 'pending' && props.viewAs === 'user'
})

const canReview = computed(() => {
  return props.order.status === 'completed' && props.viewAs === 'user'
})

const canAccept = computed(() => {
  return props.order.status === 'confirmed' && props.viewAs === 'player'
})

const canComplete = computed(() => {
  return props.order.status === 'in_progress' && props.viewAs === 'player'
})

// 格式化价格
function formatPrice(amount: number): string {
  return (amount / 100).toFixed(2)
}

// 格式化时间
function formatTime(time: string): string {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}
</script>

<style lang="scss" scoped>
.order-card {
  background: var(--color-bg-card);
  border-radius: 24rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
}

// 头部
.order-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 20rpx;
  border-bottom: 2rpx solid var(--color-divider);
  
  .order-no {
    font-size: 24rpx;
    color: var(--color-text-secondary);
  }
}

// 内容
.order-content {
  padding: 20rpx 0;
}

.order-user {
  display: flex;
  align-items: center;
  margin-bottom: 20rpx;
  
  .user-avatar {
    width: 80rpx;
    height: 80rpx;
    border-radius: 50%;
    overflow: hidden;
    background: var(--color-bg-secondary);
    margin-right: 20rpx;
    
    image {
      width: 100%;
      height: 100%;
    }
    
    .avatar-placeholder {
      width: 100%;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 32rpx;
      color: #FFFFFF;
      background: var(--color-primary);
    }
  }
  
  .user-info {
    .user-name {
      display: block;
      font-size: 28rpx;
      font-weight: 500;
      color: var(--color-text);
      margin-bottom: 4rpx;
    }
    
    .service-name {
      font-size: 24rpx;
      color: var(--color-text-secondary);
    }
  }
}

.order-details {
  display: flex;
  gap: 32rpx;
  
  .detail-item {
    .detail-label {
      display: block;
      font-size: 22rpx;
      color: var(--color-text-placeholder);
      margin-bottom: 4rpx;
    }
    
    .detail-value {
      font-size: 26rpx;
      color: var(--color-text);
      
      &.price {
        color: var(--color-primary);
        font-weight: 600;
      }
    }
  }
}

// 底部
.order-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 20rpx;
  border-top: 2rpx solid var(--color-divider);
  
  .order-time {
    font-size: 22rpx;
    color: var(--color-text-placeholder);
  }
}

.action-buttons {
  display: flex;
  gap: 16rpx;
}

.btn-action {
  padding: 12rpx 24rpx;
  font-size: 24rpx;
  border-radius: 32rpx;
  border: none;
  
  &::after {
    border: none;
  }
  
  &.btn-primary {
    background: var(--color-primary);
    color: #FFFFFF;
  }
  
  &.btn-cancel {
    background: var(--color-bg-secondary);
    color: var(--color-text-secondary);
  }
}
</style>
