<template>
  <view class="price-tag" :class="[`size-${size}`]">
    <text class="currency">¥</text>
    <text class="amount">{{ formattedAmount }}</text>
    <text v-if="unit" class="unit">/{{ unit }}</text>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  amount: number  // 金额（分）
  unit?: string   // 单位（局、小时等）
  size?: 'small' | 'medium' | 'large'
  showDecimal?: boolean
}>(), {
  unit: '',
  size: 'medium',
  showDecimal: true,
})

const formattedAmount = computed(() => {
  const yuan = props.amount / 100
  if (props.showDecimal) {
    return yuan.toFixed(2)
  }
  return Math.floor(yuan).toString()
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
  .currency { font-size: 22rpx; }
  .amount { font-size: 28rpx; }
  .unit { font-size: 20rpx; }
}

.size-medium {
  .currency { font-size: 26rpx; }
  .amount { font-size: 36rpx; }
  .unit { font-size: 24rpx; }
}

.size-large {
  .currency { font-size: 32rpx; }
  .amount { font-size: 48rpx; }
  .unit { font-size: 28rpx; }
}
</style>
