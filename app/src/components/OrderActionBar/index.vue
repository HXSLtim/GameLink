<template>
  <view v-if="actions.length > 0" class="action-bar">
    <GlButton
      v-for="action in actions"
      :key="action.key"
      :type="action.primary ? 'primary' : 'default'"
      :plain="!action.primary"
      size="large"
      round
      @click="$emit('action', action.key)"
    >
      {{ action.label }}
    </GlButton>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlButton from '@/components/gl/Button/index.vue'
import type { OrderActionKey, OrderStatus } from '@/types/order'

interface Props {
  status: OrderStatus | 'cancelled'
  hasReview?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  hasReview: false,
})

defineEmits<{
  action: [key: OrderActionKey]
}>()

interface ActionItem {
  key: OrderActionKey
  label: string
  primary?: boolean
}

const actions = computed((): ActionItem[] => {
  switch (props.status) {
    case 'pending':
      return [
        { key: 'cancel', label: '取消订单' },
        { key: 'pay', label: '立即支付', primary: true },
      ]
    case 'confirmed':
      return [
        { key: 'contact', label: '联系陪玩' },
        { key: 'refund', label: '申请退款' },
      ]
    case 'in_progress':
      return [
        { key: 'contact', label: '联系陪玩' },
        { key: 'refund', label: '申请退款' },
        { key: 'complete', label: '确认完成', primary: true },
      ]
    case 'completed':
      if (!props.hasReview) {
        return [
          { key: 'reorder', label: '再来一单' },
          { key: 'review', label: '去评价', primary: true },
        ]
      }
      return [{ key: 'reorder', label: '再来一单', primary: true }]
    case 'canceled':
    case 'cancelled':
    case 'refunded':
      return [{ key: 'reorder', label: '重新下单', primary: true }]
    default:
      return []
  }
})
</script>

<style lang="scss" scoped>
.action-bar {
  display: flex;
  gap: var(--spacing-md);
  padding: var(--spacing-sm) var(--spacing-md);
  padding-bottom: calc(var(--spacing-sm) + env(safe-area-inset-bottom));
  background: var(--color-bg);
  border-top: 1rpx solid var(--color-border);
  
  :deep(.uv-button) {
    flex: 1;
  }
}
</style>
