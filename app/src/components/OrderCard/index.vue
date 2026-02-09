<template>
  <OrderCardUser
    v-if="viewMode === 'user'"
    :order="order"
    @click="handleClick"
    @person-click="handlePersonClick"
    @action="handleAction"
  />
  <OrderCardPlayer
    v-else
    :order="order"
    @click="handleClick"
    @person-click="handlePersonClick"
    @action="handleAction"
  />
</template>

<script setup lang="ts">
import OrderCardUser from './OrderCardUser.vue'
import OrderCardPlayer from './OrderCardPlayer.vue'
import type { Order, OrderActionKey, OrderViewMode } from './types'

export type { OrderStatus, OrderViewMode, OrderPerson, Order } from './types'

interface Props {
  order: Order
  viewMode?: OrderViewMode
}

const { order, viewMode } = withDefaults(defineProps<Props>(), {
  viewMode: 'user',
})

const emit = defineEmits<{
  click: []
  'person-click': []
  action: [key: OrderActionKey, order: Order]
}>()

const handleClick = () => emit('click')
const handlePersonClick = () => emit('person-click')
const handleAction = (key: OrderActionKey, order: Order) => emit('action', key, order)
</script>
