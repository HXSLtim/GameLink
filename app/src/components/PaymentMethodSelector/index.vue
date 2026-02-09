<template>
  <GlCard title="支付方式" :shadow="false" bordered>
    <view class="payment-methods">
      <view 
        v-for="method in methods" 
        :key="method.value"
        class="payment-method"
        :class="{ selected: modelValue === method.value, disabled: !method.enabled }"
        @tap="method.enabled && $emit('update:modelValue', method.value)"
      >
        <view class="method-icon"><uv-icon :name="method.icon" size="22" color="var(--color-text-secondary)" /></view>
        <view class="method-info">
          <text class="method-name">{{ method.name }}</text>
          <text v-if="!method.enabled" class="method-tip">{{ method.tip }}</text>
        </view>
        <view class="radio-box" :class="{ checked: modelValue === method.value }">
          <uv-icon v-if="modelValue === method.value" name="checkbox-mark" size="12" color="#fff"></uv-icon>
        </view>
      </view>
    </view>
  </GlCard>
</template>

<script setup lang="ts">
import GlCard from '@/components/gl/Card/index.vue'
import type { PaymentMethod } from '@/types/wallet'

interface Props {
  methods: PaymentMethod[]
  modelValue?: PaymentMethod['value']
}

defineProps<Props>()

defineEmits<{
  'update:modelValue': [value: PaymentMethod['value']]
}>()
</script>

<style lang="scss" scoped>
.payment-methods {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.payment-method {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-sm);
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
  
  &.disabled {
    opacity: 0.5;
  }
}

.method-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.method-info {
  flex: 1;
}

.method-name {
  display: block;
  font-size: var(--font-sm);
  font-weight: 600;
  color: var(--color-text);
}

.method-tip {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}

.radio-box {
  width: 32rpx;
  height: 32rpx;
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  
  &.checked {
    background: var(--color-primary);
    border-color: var(--color-primary);
  }
}
</style>
