<template>
  <OrderCardBase
    :order="order"
    :display-person="displayPerson"
    amount-label="金额"
    :amount-value="amountValue"
    amount-unit="yuan"
    :actions="actions"
    @click="emit('click')"
    @person-click="emit('person-click')"
    @action="(...args) => emit('action', ...args)"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import OrderCardBase from './OrderCardBase.vue'
import type { ActionButton, Order, OrderActionKey, OrderStatus } from './types'

interface Props {
  order: Order
}

const props = defineProps<Props>()

const emit = defineEmits<{
  click: []
  'person-click': []
  action: [key: OrderActionKey, order: Order]
}>()

const displayPerson = computed(() => props.order.player)
const amountValue = computed(() => props.order.totalAmount ?? 0)

const getUserActions = (status: OrderStatus, reviewed?: boolean): ActionButton[] => {
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
}

const actions = computed(() => getUserActions(props.order.status, props.order.reviewed))
</script>
