<template>
  <GlCard title="选择充值金额" :shadow="false" bordered>
    <view class="amount-grid">
      <view 
        v-for="option in options" 
        :key="option.value"
        class="amount-option"
        :class="{ selected: modelValue === option.value }"
        @tap="$emit('update:modelValue', option.value)"
      >
        <text class="amount-value">¥{{ option.value }}</text>
        <text v-if="option.bonus" class="amount-bonus">送 ¥{{ option.bonus }}</text>
      </view>
    </view>
    
    <!-- 自定义金额 -->
    <view class="custom-amount">
      <text class="custom-label">自定义金额</text>
      <view class="custom-input-wrap">
        <text class="currency">¥</text>
        <input 
          type="digit"
          :value="customValue"
          class="custom-input"
          :placeholder="placeholder"
          @focus="$emit('update:modelValue', 0)"
          @input="(e: any) => $emit('update:customValue', e.detail.value)"
        />
      </view>
    </view>
    
    <text class="amount-tip">{{ tip }}</text>
  </GlCard>
</template>

<script setup lang="ts">
import GlCard from '@/components/gl/Card/index.vue'

export interface AmountOption {
  value: number
  bonus?: number
}

interface Props {
  options: AmountOption[]
  modelValue?: number
  customValue?: string
  placeholder?: string
  tip?: string
}

withDefaults(defineProps<Props>(), {
  placeholder: '10-10000',
  tip: '充值金额需在 ¥10 - ¥10000 之间',
})

defineEmits<{
  'update:modelValue': [value: number]
  'update:customValue': [value: string]
}>()
</script>

<style lang="scss" scoped>
.amount-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16rpx;
  margin-bottom: 24rpx;
}

.amount-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24rpx 16rpx;
  background: var(--color-bg-secondary);
  border-radius: 16rpx;
  border: 2rpx solid var(--color-border);
  transition: all 0.2s;
  
  &.selected {
    border-color: var(--color-primary);
    background: rgba(0, 210, 106, 0.1);
  }
}

.amount-value {
  font-size: 36rpx;
  font-weight: 700;
  color: var(--color-text);
}

.amount-bonus {
  font-size: 22rpx;
  color: var(--color-primary);
  margin-top: 8rpx;
  font-weight: 600;
}

.custom-amount {
  display: flex;
  align-items: center;
  gap: 20rpx;
  margin-bottom: 16rpx;
}

.custom-label {
  font-size: 28rpx;
  color: var(--color-text);
  flex-shrink: 0;
}

.custom-input-wrap {
  flex: 1;
  display: flex;
  align-items: center;
  padding: 16rpx 20rpx;
  background: var(--color-bg-secondary);
  border-radius: 12rpx;
  border: 2rpx solid var(--color-border);
}

.currency {
  font-size: 32rpx;
  font-weight: 600;
  color: var(--color-text);
  margin-right: 8rpx;
}

.custom-input {
  flex: 1;
  font-size: 32rpx;
  color: var(--color-text);
}

.amount-tip {
  font-size: 24rpx;
  color: var(--color-text-secondary);
}
</style>
