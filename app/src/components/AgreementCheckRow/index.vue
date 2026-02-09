<template>
  <view class="agreement">
    <view class="checkbox" :class="{ checked: modelValue }" @tap="toggle">
      <uv-icon v-if="modelValue" name="checkbox-mark" size="14" color="#fff"></uv-icon>
    </view>
    <text class="agreement-text">
      {{ prefixText }}
      <text class="link" @tap="$emit('link')">《{{ linkText }}》</text>
    </text>
  </view>
</template>

<script setup lang="ts">
interface Props {
  modelValue: boolean
  prefixText?: string
  linkText: string
}

const props = withDefaults(defineProps<Props>(), {
  prefixText: '我已阅读并同意',
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  link: []
}>()

const toggle = () => {
  emit('update:modelValue', !props.modelValue)
}
</script>

<style lang="scss" scoped>
.agreement {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  padding: var(--spacing-md);
  cursor: pointer;
  @include press-effect;
}

.checkbox {
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

.agreement-text {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}

.link {
  color: var(--color-primary);
}
</style>
