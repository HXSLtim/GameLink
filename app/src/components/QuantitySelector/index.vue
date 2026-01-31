<template>
  <GlCard :title="title" required :shadow="false" bordered>
    <view class="quantity-selector">
      <view class="quantity-btn" :class="{ disabled: modelValue <= min }" @tap="decrease">
        <uv-icon name="minus" size="18" :color="modelValue <= min ? 'var(--color-text-placeholder)' : 'var(--color-text)'"></uv-icon>
      </view>
      <input 
        type="number" 
        :value="modelValue"
        class="quantity-input"
        @blur="handleInput"
      />
      <view class="quantity-btn" :class="{ disabled: modelValue >= max }" @tap="increase">
        <uv-icon name="plus" size="18" :color="modelValue >= max ? 'var(--color-text-placeholder)' : 'var(--color-text)'"></uv-icon>
      </view>
    </view>
    <text v-if="tip" class="quantity-tip">{{ tip }}</text>
  </GlCard>
</template>

<script setup lang="ts">
import GlCard from '@/components/gl/Card/index.vue'

interface Props {
  title?: string
  modelValue: number
  min?: number
  max?: number
  tip?: string
}

const props = withDefaults(defineProps<Props>(), {
  title: '数量',
  min: 1,
  max: 99,
})

const emit = defineEmits<{
  'update:modelValue': [value: number]
}>()

const decrease = () => {
  if (props.modelValue > props.min) {
    emit('update:modelValue', props.modelValue - 1)
  }
}

const increase = () => {
  if (props.modelValue < props.max) {
    emit('update:modelValue', props.modelValue + 1)
  }
}

const handleInput = (e: any) => {
  let val = parseInt(e.detail.value) || props.min
  val = Math.max(props.min, Math.min(props.max, val))
  emit('update:modelValue', val)
}
</script>

<style lang="scss" scoped>
.quantity-selector {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 24rpx;
}

.quantity-btn {
  width: 64rpx;
  height: 64rpx;
  border-radius: 16rpx;
  background: var(--color-bg-secondary);
  border: 2rpx solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  
  &:active:not(.disabled) {
    background: var(--color-bg-hover);
  }
  
  &.disabled {
    opacity: 0.5;
  }
}

.quantity-input {
  width: 120rpx;
  height: 64rpx;
  text-align: center;
  font-size: 36rpx;
  font-weight: 700;
  color: var(--color-text);
  background: var(--color-bg-secondary);
  border: 2rpx solid var(--color-border);
  border-radius: 16rpx;
}

.quantity-tip {
  display: block;
  text-align: center;
  margin-top: 20rpx;
  font-size: 24rpx;
  color: var(--color-text-secondary);
}
</style>
