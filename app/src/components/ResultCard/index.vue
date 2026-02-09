<template>
  <StatusCard
    :status-class="type"
    border-width="2rpx"
    :border-color="borderColor"
    direction="column"
    align="center"
    gap="0"
    padding="var(--spacing-lg)"
    margin="var(--spacing-md)"
  >
    <StatusInfo
      :icon="icon"
      :title="title"
      :description="description"
      size="lg"
      direction="column"
      align="center"
      gap="var(--spacing-sm)"
    />
    
    <view v-if="amount" class="amount-info">
      <text class="amount-label">{{ amountLabel }}</text>
      <PriceTag class="amount-value" :amount="amount" amount-unit="cents" size="small" />
    </view>
  </StatusCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import StatusCard from '@/components/StatusCard/index.vue'
import StatusInfo from '@/components/StatusInfo/index.vue'
import PriceTag from '@/components/PriceTag/index.vue'
import { getResultStatusPreset, type ResultStatusType } from '@/components/StatusCard/presets'

interface Props {
  type: ResultStatusType
  title: string
  description?: string
  amount?: number // 分
  amountLabel?: string
}

const props = withDefaults(defineProps<Props>(), {
  amountLabel: '支付金额',
})

const icon = computed(() => getResultStatusPreset(props.type).icon)
const borderColor = computed(() => getResultStatusPreset(props.type).borderColor)
</script>

<style lang="scss" scoped>
.amount-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: var(--spacing-md);
  padding-top: var(--spacing-md);
  border-top: 1rpx solid var(--color-border);
  width: 100%;
}

.amount-label {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  margin-bottom: var(--spacing-xs);
}

.amount-value {
  font-size: var(--font-base);
  font-weight: 600;
  color: var(--color-text);
}

.amount-value :deep(.price-tag) {
  color: var(--color-text);
}
</style>
