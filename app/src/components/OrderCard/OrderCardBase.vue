<template>
  <view class="order-card" @tap="$emit('click')">
    <!-- 订单头部 -->
    <view class="order-header">
      <text class="order-no">订单号：{{ order.orderNo }}</text>
      <view class="order-status" :class="`order-status--${statusInfo.class}`">
        {{ statusInfo.text }}
      </view>
    </view>
    
    <!-- 对方信息 -->
    <view class="order-person" @tap.stop="$emit('person-click')">
      <GlAvatar
        :src="displayPerson?.avatar"
        :text="displayPerson?.nickname"
        size="medium"
      />
      <view class="person-info">
        <text class="person-name">{{ displayPerson?.nickname }}</text>
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
      <view class="info-item info-item--right">
        <text class="info-label">{{ amountLabel }}</text>
        <PriceTag
          class="info-value info-value--price"
          :amount="amountValue"
          :amount-unit="amountUnit"
          size="small"
        />
      </view>
    </view>
    
    <!-- 备注 -->
    <view v-if="showRemark && order.remark" class="order-remark">
      <text class="remark-label">备注：</text>
      <text class="remark-content">{{ order.remark }}</text>
    </view>
    
    <!-- 订单时间 -->
    <view class="order-time">
      <text>{{ formattedTime }}</text>
    </view>
    
    <!-- 操作按钮 -->
    <view class="order-actions">
      <template v-for="action in actions" :key="action.key">
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
import PriceTag from '@/components/PriceTag/index.vue'
import { formatDateTime } from '@/utils/format'
import type { ActionButton, Order, OrderActionKey, OrderPerson, OrderStatus } from './types'

interface Props {
  order: Order
  displayPerson?: OrderPerson
  amountLabel: string
  amountValue: number
  amountUnit: 'cents' | 'yuan'
  actions?: ActionButton[]
  showRemark?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  actions: () => [],
  showRemark: false,
})

const emit = defineEmits<{
  click: []
  'person-click': []
  action: [key: OrderActionKey, order: Order]
}>()

const statusMap: Record<OrderStatus, { text: string; class: string }> = {
  pending: { text: '待支付', class: 'warning' },
  confirmed: { text: '待服务', class: 'info' },
  in_progress: { text: '进行中', class: 'primary' },
  completed: { text: '已完成', class: 'success' },
  canceled: { text: '已取消', class: 'default' },
  refunding: { text: '退款中', class: 'info' },
  refunded: { text: '已退款', class: 'default' },
  disputed: { text: '争议中', class: 'error' },
}

const statusInfo = computed(() => statusMap[props.order.status] || { text: props.order.status, class: 'default' })
const formattedTime = computed(() => formatDateTime(props.order.createdAt))

const handleAction = (key: string) => {
  emit('action', key, props.order)
}
</script>

<style lang="scss" scoped>
.order-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1rpx solid var(--color-border);
  padding: var(--spacing-md);
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease, background 0.2s ease;
  position: relative;
  overflow: hidden;
  cursor: pointer;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-md);
    border-color: var(--color-primary);
  }
  
  &:active {
    transform: translateY(0);
    background: var(--color-bg-secondary);
    border-color: var(--color-border);
  }
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: var(--spacing-xs);
  border-bottom: 1rpx solid var(--color-border);
  margin-bottom: var(--spacing-sm);
}

.order-no {
  flex: 1;
  min-width: 0;
  font-size: var(--font-xs);
  color: var(--color-text-placeholder);
  letter-spacing: 0.2rpx;
  @include text-ellipsis;
}

.order-status {
  font-size: var(--font-sm);
  font-weight: 600;
  padding: 2rpx var(--spacing-xs);
  border-radius: var(--radius-full);
  border: 1rpx solid var(--color-border);
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  
  &::before {
    content: '';
    width: 12rpx;
    height: 12rpx;
    border-radius: 50%;
    background: currentColor;
    opacity: 0.8;
  }
  
  &--warning { color: var(--color-warning); background: var(--color-warning-tint); border-color: var(--color-warning-tint-border); }
  &--info { color: var(--color-info); background: var(--color-info-tint); border-color: var(--color-info-tint-border); }
  &--primary { color: var(--color-primary); background: var(--color-primary-tint); border-color: var(--color-primary-tint-border); }
  &--success { color: var(--color-success); background: var(--color-success-tint); border-color: var(--color-success-tint-border); }
  &--default { color: var(--color-text-placeholder); background: var(--color-bg-secondary); border-color: var(--color-border); }
  &--error { color: var(--color-error); background: var(--color-error-tint); border-color: var(--color-error-tint-border); }
}

.order-person {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-sm);
  padding: 0 var(--spacing-xs);
  cursor: pointer;
  @include press-effect;
}

.person-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
  min-width: 0;
}

.person-name {
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
  @include text-ellipsis;
}

.service-name {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  background: var(--color-bg-secondary);
  border: 1rpx solid var(--color-border);
  padding: 2rpx var(--spacing-xs);
  border-radius: var(--radius-full);
  align-self: flex-start;
  max-width: 100%;
  @include text-ellipsis;
}

.order-info {
  display: flex;
  justify-content: space-between;
  padding: var(--spacing-sm);
  background: var(--color-bg-secondary);
  border-radius: var(--radius-md);
  margin-bottom: var(--spacing-sm);
  border: 1rpx solid var(--color-border);
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
  flex: 1;
  min-width: 0;
  
  &:not(:last-child) {
    border-right: 1rpx solid var(--color-border);
    margin-right: var(--spacing-sm);
    padding-right: var(--spacing-sm);
  }
}

.info-item--right {
  align-items: flex-end;
  text-align: right;
}

.info-label {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}

.info-value {
  display: block;
  font-size: var(--font-sm);
  font-weight: 600;
  color: var(--color-text);
  @include text-ellipsis;
  
  &--price {
    font-weight: 700;
    color: var(--color-primary);
    font-size: var(--font-md);
  }
}

.info-value--price :deep(.amount) {
  font-weight: 700;
}

.order-remark {
  padding: var(--spacing-sm);
  background: var(--color-bg-secondary);
  border-radius: var(--radius-md);
  border: 1rpx solid var(--color-border);
  margin-bottom: var(--spacing-sm);
  
  .remark-label {
    font-size: var(--font-xs);
    color: var(--color-text-secondary);
  }
  
  .remark-content {
    font-size: var(--font-sm);
    color: var(--color-text);
  }
}

.order-time {
  margin-bottom: var(--spacing-xs);
  padding-left: var(--spacing-xs);
  
  text {
    font-size: var(--font-xs);
    color: var(--color-text-placeholder);
    opacity: 0.8;
  }
}

.order-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-xs);
  padding-top: var(--spacing-xs);
  border-top: 1rpx solid var(--color-border);
  flex-wrap: wrap;
  row-gap: var(--spacing-xs);
}
</style>
