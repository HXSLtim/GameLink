<template>
  <view class="price-tag" :class="[`size-${size}`]">
    <text v-if="showCurrency" class="currency">¥</text>
    <text class="amount">{{ formattedAmount }}</text>
    <text v-if="unit" class="unit">/{{ unit }}</text>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { formatMoney, formatYuan } from '@/utils/format'

const props = withDefaults(defineProps<{
  amount: number  // 金额
  unit?: string   // 单位（局、小时等）
  size?: 'small' | 'medium' | 'large'
  showDecimal?: boolean
  amountUnit?: 'cents' | 'yuan'
  showCurrency?: boolean
}>(), {
  unit: '',
  size: 'medium',
  showDecimal: true,
  amountUnit: 'cents',
  showCurrency: true,
})

const formattedAmount = computed(() => {
  if (props.amountUnit === 'yuan') {
    return props.showDecimal ? formatYuan(props.amount) : Math.floor(props.amount).toString()
  }
  return props.showDecimal ? formatMoney(props.amount) : Math.floor(props.amount / 100).toString()
})
</script>

<style lang="scss" scoped>
.price-tag {
  display: inline-flex;
  align-items: baseline;
  color: var(--color-primary);
  
  .currency {
    font-weight: 500;
  }
  
  .amount {
    font-weight: 600;
  }
  
  .unit {
    color: var(--color-text-secondary);
    font-weight: 400;
  }
}

// 尺寸
.size-small {
  .currency { font-size: var(--font-xs); }
  .amount { font-size: var(--font-md); }
  .unit { font-size: var(--font-xs); }
}

.size-medium {
  .currency { font-size: var(--font-sm); }
  .amount { font-size: var(--font-lg); }
  .unit { font-size: var(--font-sm); }
}

.size-large {
  .currency { font-size: var(--font-md); }
  .amount { font-size: var(--font-xl); }
  .unit { font-size: var(--font-sm); }
}
</style>
