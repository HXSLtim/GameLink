<template>
  <StatusCard
    :status-class="statusClass"
    border-width="2rpx"
    :border-color="borderColor"
    margin="0"
    radius="0 0 var(--radius-md) var(--radius-md)"
  >
    <StatusInfo
      :icon="statusIcon"
      :title="statusText"
      :description="statusDesc"
      size="sm"
    />
    <view v-if="status === 'pending' && countdown > 0" class="countdown">
      <text>剩余支付时间：{{ formattedCountdown }}</text>
    </view>
  </StatusCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import StatusCard from '@/components/StatusCard/index.vue'
import StatusInfo from '@/components/StatusInfo/index.vue'
import { getOrderStatusPreset } from '@/components/StatusCard/presets'
import type { OrderStatus } from '@/types/order'

interface Props {
  status: OrderStatus | 'cancelled'
  countdown?: number
}

const props = withDefaults(defineProps<Props>(), {
  countdown: 0,
})

const config = computed(() => getOrderStatusPreset(props.status))
const statusIcon = computed(() => config.value.icon)
const statusText = computed(() => config.value.text)
const statusDesc = computed(() => config.value.desc)
const statusClass = computed(() => config.value.className)
const borderColor = computed(() => config.value.borderColor)

const formattedCountdown = computed(() => {
  const minutes = Math.floor(props.countdown / 60)
  const seconds = props.countdown % 60
  return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
})
</script>

<style lang="scss" scoped>
.countdown {
  padding: 2rpx var(--spacing-sm);
  background: var(--color-bg-secondary);
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}
</style>
