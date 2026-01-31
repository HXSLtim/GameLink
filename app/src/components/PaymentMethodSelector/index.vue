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
        <text class="method-icon">{{ method.icon }}</text>
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

export interface PaymentMethod {
  value: string
  name: string
  icon: string
  enabled: boolean
  tip?: string
}

interface Props {
  methods: PaymentMethod[]
  modelValue?: string
}

defineProps<Props>()

defineEmits<{
  'update:modelValue': [value: string]
}>()
</script>

<style lang="scss" scoped>
.payment-methods {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.payment-method {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 24rpx;
  background: var(--color-bg-secondary);
  border-radius: 16rpx;
  border: 2rpx solid var(--color-border);
  transition: all 0.2s;
  
  &.selected {
    border-color: var(--color-primary);
    background: rgba(0, 210, 106, 0.08);
  }
  
  &.disabled {
    opacity: 0.5;
  }
}

.method-icon {
  font-size: 40rpx;
}

.method-info {
  flex: 1;
}

.method-name {
  display: block;
  font-size: 30rpx;
  font-weight: 500;
  color: var(--color-text);
}

.method-tip {
  font-size: 24rpx;
  color: var(--color-text-secondary);
}

.radio-box {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  border: 2rpx solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  
  &.checked {
    background: var(--color-primary);
    border-color: var(--color-primary);
  }
}
</style>
