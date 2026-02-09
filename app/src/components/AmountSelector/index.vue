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
        <PriceTag
          class="amount-value"
          :amount="option.value"
          amount-unit="yuan"
          size="small"
          :show-decimal="false"
        />
        <view v-if="option.bonus" class="amount-bonus">
          <text>送</text>
          <PriceTag
            :amount="option.bonus"
            amount-unit="yuan"
            size="small"
            :show-decimal="false"
          />
        </view>
      </view>
    </view>
    
    <!-- 自定义金额 -->
    <view class="custom-amount">
      <text class="custom-label">自定义金额</text>
      <view class="custom-input-wrap">
        <text class="currency">¥</text>
        <GlInput
          class="custom-input"
          type="digit"
          size="small"
          variant="plain"
          :model-value="customValue ?? ''"
          :placeholder="placeholder"
          @focus="$emit('update:modelValue', 0)"
          @update:modelValue="(value) => $emit('update:customValue', value)"
        />
      </view>
    </view>
    
    <text class="amount-tip">{{ tip }}</text>
  </GlCard>
</template>

<script setup lang="ts">
import GlCard from '@/components/gl/Card/index.vue'
import GlInput from '@/components/gl/Input/index.vue'
import PriceTag from '@/components/PriceTag/index.vue'
import type { AmountOption } from '@/types/wallet'

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
  gap: var(--spacing-xs);
  margin-bottom: var(--spacing-md);
}

.amount-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-sm) var(--spacing-sm);
  background: var(--color-bg-card);
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
  transition: all 0.2s;
  cursor: pointer;
  @include press-effect;
  
  &.selected {
    border-color: var(--color-border);
    background: var(--color-bg-secondary);
  }
}

.amount-value {
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
}

.amount-value :deep(.price-tag) {
  color: var(--color-text);
}

.amount-value :deep(.price-tag .amount) {
  font-weight: 600;
}

.amount-bonus {
  display: inline-flex;
  align-items: baseline;
  gap: 4rpx;
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
  margin-top: var(--spacing-xs);
  font-weight: 600;
}

.amount-bonus :deep(.price-tag) {
  color: var(--color-text-secondary);
}

.custom-amount {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-sm);
}

.custom-label {
  font-size: var(--font-sm);
  color: var(--color-text);
  flex-shrink: 0;
}

.custom-input-wrap {
  flex: 1;
  display: flex;
  align-items: center;
  padding: var(--spacing-xs) var(--spacing-md);
  background: var(--color-bg-card);
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
}

.currency {
  font-size: var(--font-sm);
  font-weight: 600;
  color: var(--color-text);
  margin-right: var(--spacing-xs);
}

.custom-input {
  flex: 1;
  min-width: 0;
  
  :deep(.gl-input__field) {
    font-size: var(--font-base);
    color: var(--color-text);
  }
}

.amount-tip {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}
</style>
