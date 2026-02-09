<template>
  <OrderCardBase
    :order="order"
    :display-person="displayPerson"
    amount-label="收益"
    :amount-value="amountValue"
    amount-unit="cents"
    :actions="actions"
    show-remark
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

const displayPerson = computed(() => props.order.user)
const amountValue = computed(() => props.order.earnings ?? 0)

const getPlayerActions = (status: OrderStatus): ActionButton[] => {
  switch (status) {
    case 'pending':
      return [
        { key: 'reject', label: '拒绝', type: 'default', plain: true },
        { key: 'accept', label: '接单', type: 'primary', plain: false },
      ]
    case 'confirmed':
      return [
        { key: 'contact', label: '联系用户', type: 'default', plain: true },
        { key: 'start', label: '开始服务', type: 'primary', plain: false },
      ]
    case 'in_progress':
      return [
        { key: 'contact', label: '联系用户', type: 'default', plain: true },
        { key: 'complete', label: '完成服务', type: 'primary', plain: false },
      ]
    default:
      return [{ key: 'detail', label: '查看详情', type: 'default', plain: true }]
  }
}

const actions = computed(() => getPlayerActions(props.order.status))
</script>
