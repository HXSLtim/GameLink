<template>
  <RecordListItem gap="20rpx" padding="24rpx 0" @click="$emit('click')">
    <template #icon>
      <view class="detail-icon">
        <uv-icon :name="icon" size="28" color="var(--color-text-secondary)" />
      </view>
    </template>
    <template #content>
      <view class="detail-info">
        <text class="detail-title">{{ title }}</text>
        <text class="detail-desc">{{ description }}</text>
        <text class="detail-time">{{ formattedTime }}</text>
      </view>
    </template>
    <template #amount>
      <view class="detail-amount" :class="{ negative: amount < 0 }">
        <text v-if="amount > 0" class="amount-sign">+</text>
        <PriceTag :amount="Math.abs(amount)" amount-unit="cents" size="small" />
      </view>
    </template>
  </RecordListItem>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { formatMonthDayTime } from '@/utils/format'
import PriceTag from '@/components/PriceTag/index.vue'
import RecordListItem from '@/components/RecordListItem/index.vue'

interface Props {
  type: string
  title: string
  description?: string
  amount: number // 分
  createdAt: string
}

const props = defineProps<Props>()

defineEmits<{
  click: []
}>()

const icons: Record<string, string> = {
  order: 'red-packet',
  withdraw: 'wallet',
  bonus: 'gift',
  refund: 'reload',
}

const icon = computed(() => icons[props.type] || 'red-packet')

const formattedTime = computed(() => formatMonthDayTime(props.createdAt))
</script>

<style lang="scss" scoped>
.detail-icon {
  width: 72rpx;
  height: 72rpx;
  background: var(--color-bg-secondary);
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.detail-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.detail-title {
  font-size: var(--font-base);
  font-weight: 500;
  color: var(--color-text);
}

.detail-desc {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}

.detail-time {
  font-size: var(--font-xs);
  color: var(--color-text-placeholder);
}

.detail-amount {
  display: flex;
  align-items: baseline;
  gap: var(--spacing-xs);
  flex-shrink: 0;
  font-size: var(--font-base);
  font-weight: 600;
  color: var(--color-primary);
  
  &.negative {
    color: var(--color-text);
  }
}

.amount-sign {
  font-size: var(--font-base);
}

.detail-amount :deep(.price-tag) {
  color: inherit;
}
</style>
