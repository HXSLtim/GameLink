<template>
  <SectionCard class="summary-card" margin="var(--spacing-md)" padding="var(--spacing-lg)">
    <view class="summary-header">
      <text class="summary-label">可提现金额（元）</text>
      <GlButton type="primary" size="mini" round @click="$emit('withdraw')">提现</GlButton>
    </view>
    <view class="summary-balance">
      <PriceTag
        class="balance-value"
        :amount="withdrawable"
        amount-unit="cents"
        size="large"
        :show-currency="false"
      />
    </view>
    <view class="summary-row">
      <HeaderStatsRow :items="stats" :show-divider="false" size="md" item-padding="0">
        <template #value="{ item }">
          <PriceTag
            :class="item.key === 'total' ? 'summary-highlight' : 'summary-normal'"
            :amount="item.value"
            amount-unit="cents"
            size="small"
          />
        </template>
      </HeaderStatsRow>
    </view>
  </SectionCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import SectionCard from '@/components/SectionCard/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
import PriceTag from '@/components/PriceTag/index.vue'
import HeaderStatsRow from '@/components/HeaderStatsRow/index.vue'
import type { HeaderStatItem } from '@/types/ui'

interface Props {
  withdrawable: number // 分
  pending: number // 分
  withdrawn: number // 分
  total: number // 分
}

const props = defineProps<Props>()

defineEmits<{
  withdraw: []
}>()

const stats = computed<HeaderStatItem[]>(() => [
  { key: 'pending', label: '待结算', value: props.pending },
  { key: 'withdrawn', label: '已提现', value: props.withdrawn },
  { key: 'total', label: '累计收益', value: props.total },
])

</script>

<style lang="scss" scoped>
.summary-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--spacing-sm);
}

.summary-label {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}

.summary-balance {
  margin-bottom: var(--spacing-lg);
}

.balance-value {
  font-size: var(--font-xl);
  font-weight: 700;
}

.summary-row {
  padding-top: var(--spacing-md);
  border-top: 1rpx solid var(--color-border);
}

.summary-normal :deep(.price-tag) {
  color: var(--color-text);
}

.summary-highlight :deep(.price-tag) {
  color: var(--color-primary);
}

.balance-value :deep(.price-tag) {
  color: var(--color-text);
}

.balance-value :deep(.price-tag .amount) {
  font-weight: 700;
  font-size: var(--font-xl);
}
</style>
