<template>
  <GlCard :title="title" required :shadow="false" bordered>
    <view class="quantity-selector">
      <view class="quantity-btn" :class="{ disabled: modelValue <= min }" @tap="decrease">
        <uv-icon name="minus" size="18" :color="modelValue <= min ? 'var(--color-text-placeholder)' : 'var(--color-text)'"></uv-icon>
      </view>
      <GlInput
        class="quantity-input"
        type="number"
        size="small"
        align="center"
        :model-value="inputValue"
        @update:modelValue="handleInput"
        @blur="handleBlur"
      />
      <view class="quantity-btn" :class="{ disabled: modelValue >= max }" @tap="increase">
        <uv-icon name="plus" size="18" :color="modelValue >= max ? 'var(--color-text-placeholder)' : 'var(--color-text)'"></uv-icon>
      </view>
    </view>
    <text v-if="tip" class="quantity-tip">{{ tip }}</text>
  </GlCard>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import GlCard from '@/components/gl/Card/index.vue'
import GlInput from '@/components/gl/Input/index.vue'

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

const inputValue = ref(String(props.modelValue))

watch(
  () => props.modelValue,
  (value) => {
    inputValue.value = String(value ?? '')
  }
)

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

const handleInput = (value: string) => {
  inputValue.value = value
}

const handleBlur = () => {
  let val = parseInt(inputValue.value) || props.min
  val = Math.max(props.min, Math.min(props.max, val))
  emit('update:modelValue', val)
  inputValue.value = String(val)
}
</script>

<style lang="scss" scoped>
.quantity-selector {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
}

.quantity-btn {
  width: 56rpx;
  height: 56rpx;
  border-radius: var(--radius-sm);
  background: var(--color-bg-card);
  border: 1rpx solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  cursor: pointer;
  @include press-effect;
  
  &:active:not(.disabled) {
    background: var(--color-bg-secondary);
  }
  
  &.disabled {
    opacity: 0.5;
  }
}

.quantity-input {
  width: 112rpx;
  min-width: 0;
  
  :deep(.gl-input__field) {
    font-size: var(--font-md);
    font-weight: 600;
    color: var(--color-text);
  }
}

.quantity-tip {
  display: block;
  text-align: center;
  margin-top: var(--spacing-sm);
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}
</style>
