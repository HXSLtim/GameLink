<template>
  <view class="order-card">
    <!-- 订单头部 -->
    <view class="order-header">
      <text class="order-no">订单号：{{ order.orderNo }}</text>
      <GlTag :type="statusConfig.type" size="small">{{ statusConfig.text }}</GlTag>
    </view>
    
    <!-- 用户信息 -->
    <view class="order-user" @tap="$emit('user-click')">
      <GlAvatar :src="order.user?.avatar" :text="order.user?.nickname" size="medium" />
      <view class="user-info">
        <text class="user-name">{{ order.user?.nickname }}</text>
        <text class="order-desc">{{ order.gameName }} · {{ order.serviceName }}</text>
      </view>
    </view>
    
    <!-- 订单信息 -->
    <view class="order-info">
      <view class="info-item">
        <text class="info-label">数量</text>
        <text class="info-value">{{ order.quantity }}{{ order.unit }}</text>
      </view>
      <view class="info-item">
        <text class="info-label">收益</text>
        <text class="info-value earnings">¥{{ (order.earnings / 100).toFixed(2) }}</text>
      </view>
      <view class="info-item">
        <text class="info-label">下单时间</text>
        <text class="info-value">{{ formattedTime }}</text>
      </view>
    </view>
    
    <!-- 备注 -->
    <view v-if="order.remark" class="order-remark">
      <text class="remark-label">备注：</text>
      <text class="remark-content">{{ order.remark }}</text>
    </view>
    
    <!-- 操作按钮 -->
    <view class="order-actions">
      <GlButton
        v-for="action in actions"
        :key="action.key"
        :type="action.primary ? 'primary' : 'default'"
        :plain="!action.primary"
        size="small"
        @click="$emit('action', action.key)"
      >
        {{ action.label }}
      </GlButton>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlAvatar from '@/components/gl/Avatar/index.vue'
import GlTag from '@/components/gl/Tag/index.vue'
import GlButton from '@/components/gl/Button/index.vue'

export interface PlayerOrderUser {
  id: number
  nickname: string
  avatar?: string
}

export interface PlayerOrderData {
  id: number
  orderNo: string
  status: string
  user: PlayerOrderUser
  gameName: string
  serviceName: string
  quantity: number
  unit: string
  earnings: number
  remark?: string
  createdAt: string
}

interface Props {
  order: PlayerOrderData
}

const props = defineProps<Props>()

defineEmits<{
  'user-click': []
  action: [key: string]
}>()

const statusConfigs: Record<string, { text: string; type: 'primary' | 'success' | 'warning' | 'error' | 'default' }> = {
  pending: { text: '待接单', type: 'warning' },
  confirmed: { text: '待服务', type: 'primary' },
  in_progress: { text: '服务中', type: 'primary' },
  completed: { text: '已完成', type: 'success' },
  canceled: { text: '已取消', type: 'default' },
  refunded: { text: '已退款', type: 'default' },
}

const statusConfig = computed(() => statusConfigs[props.order.status] || { text: props.order.status, type: 'default' as const })

const formattedTime = computed(() => {
  const date = new Date(props.order.createdAt)
  return `${date.getMonth() + 1}/${date.getDate()} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
})

interface ActionItem {
  key: string
  label: string
  primary?: boolean
}

const actions = computed((): ActionItem[] => {
  switch (props.order.status) {
    case 'pending':
      return [
        { key: 'reject', label: '拒绝' },
        { key: 'accept', label: '接单', primary: true },
      ]
    case 'confirmed':
      return [
        { key: 'contact', label: '联系用户' },
        { key: 'start', label: '开始服务', primary: true },
      ]
    case 'in_progress':
      return [
        { key: 'contact', label: '联系用户' },
        { key: 'complete', label: '完成服务', primary: true },
      ]
    default:
      return [{ key: 'detail', label: '查看详情' }]
  }
})
</script>

<style lang="scss" scoped>
.order-card {
  background: var(--color-bg-card);
  border-radius: 20rpx;
  padding: 24rpx;
  margin-bottom: 20rpx;
  border: 2rpx solid var(--color-border);
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 20rpx;
  border-bottom: 1rpx solid var(--color-border);
  margin-bottom: 20rpx;
}

.order-no {
  font-size: 24rpx;
  color: var(--color-text-secondary);
}

.order-user {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-bottom: 20rpx;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  display: block;
  font-size: 30rpx;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 4rpx;
}

.order-desc {
  font-size: 24rpx;
  color: var(--color-text-secondary);
}

.order-info {
  display: flex;
  gap: 32rpx;
  padding: 20rpx;
  background: var(--color-bg-secondary);
  border-radius: 12rpx;
  margin-bottom: 16rpx;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.info-label {
  font-size: 22rpx;
  color: var(--color-text-secondary);
}

.info-value {
  font-size: 28rpx;
  color: var(--color-text);
  font-weight: 500;
  
  &.earnings {
    color: var(--color-primary);
    font-weight: 700;
  }
}

.order-remark {
  padding: 16rpx;
  background: var(--color-bg-secondary);
  border-radius: 12rpx;
  margin-bottom: 16rpx;
}

.remark-label {
  font-size: 24rpx;
  color: var(--color-text-secondary);
}

.remark-content {
  font-size: 26rpx;
  color: var(--color-text);
}

.order-actions {
  display: flex;
  justify-content: flex-end;
  gap: 16rpx;
  padding-top: 16rpx;
  border-top: 1rpx solid var(--color-border);
}
</style>
