<template>
  <RecordListItem pressable>
    <template #icon>
      <view class="record-icon" :class="record.type">
        <uv-icon :name="iconName" size="20" :color="iconColor"></uv-icon>
      </view>
    </template>
    <template #content>
      <view class="record-info">
        <text class="record-title">{{ record.title }}</text>
        <text class="record-desc">{{ record.description }}</text>
        <text class="record-time">{{ formattedTime }}</text>
      </view>
    </template>
    <template #amount>
      <view class="record-amount" :class="{ positive: record.amount > 0 }">
        <text v-if="record.amount > 0" class="amount-sign">+</text>
        <text v-else-if="record.amount < 0" class="amount-sign">-</text>
        <PriceTag :amount="Math.abs(record.amount)" amount-unit="cents" size="medium" :show-currency="false" />
      </view>
    </template>
  </RecordListItem>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { formatMonthDayTime } from '@/utils/format'
import type { TransactionData, TransactionType } from '@/types/wallet'
import PriceTag from '@/components/PriceTag/index.vue'
import RecordListItem from '@/components/RecordListItem/index.vue'

interface Props {
  record: TransactionData
}

const props = defineProps<Props>()

const iconConfigs: Record<TransactionType, { name: string; color: string }> = {
  recharge: { name: 'plus-circle', color: 'var(--color-success)' },
  withdraw: { name: 'red-packet', color: 'var(--color-warning)' },
  withdrawal: { name: 'red-packet', color: 'var(--color-warning)' },
  payment: { name: 'shopping-cart', color: 'var(--color-info)' },
  refund: { name: 'reload', color: 'var(--color-primary)' },
  earning: { name: 'money', color: 'var(--color-success)' },
  bonus: { name: 'gift', color: 'var(--color-error)' },
  commission: { name: 'money', color: 'var(--color-success)' },
}

const iconName = computed(() => iconConfigs[props.record.type]?.name || 'list')
const iconColor = computed(() => iconConfigs[props.record.type]?.color || 'var(--color-text-placeholder)')

const formattedTime = computed(() => formatMonthDayTime(props.record.createdAt))
</script>

<style lang="scss" scoped>
.record-icon {
  width: 64rpx;
  height: 64rpx;
  border-radius: var(--radius-full);
  background: var(--color-bg-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.record-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.record-title {
  font-size: var(--font-base);
  font-weight: 600;
  color: var(--color-text);
}

.record-desc {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}

.record-time {
  font-size: var(--font-xs);
  color: var(--color-text-placeholder);
}

.record-amount {
  display: flex;
  align-items: baseline;
  gap: 2rpx;
  flex-shrink: 0;
  font-size: var(--font-lg);
  font-weight: 600;
  color: var(--color-text);
  
  &.positive {
    color: var(--color-primary);
  }
}

.amount-sign {
  font-size: var(--font-lg);
}

.record-amount :deep(.price-tag) {
  color: inherit;
}
</style>
