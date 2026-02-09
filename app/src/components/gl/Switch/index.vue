<template>
  <view
    class="gl-switch"
    :class="[`gl-switch--${size}`, { 'gl-switch--disabled': disabled }]"
    :style="customStyle"
  >
    <switch
      :checked="modelValue"
      :disabled="disabled"
      :color="color"
      @change="handleChange"
    />
  </view>
</template>

<script setup lang="ts">
type SwitchSize = 'small' | 'medium' | 'large'

interface Props {
  modelValue?: boolean
  disabled?: boolean
  color?: string
  size?: SwitchSize
  customStyle?: string | Record<string, string>
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: false,
  disabled: false,
  color: 'var(--color-primary)',
  size: 'medium',
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  change: [value: boolean]
}>()

const handleChange = (e: any) => {
  const value = Boolean(e?.detail?.value)
  emit('update:modelValue', value)
  emit('change', value)
}
</script>

<style lang="scss" scoped>
.gl-switch {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  
  switch {
    transform-origin: center;
  }
  
  &--small {
    switch {
      transform: scale(0.85);
    }
  }
  
  &--medium {
    switch {
      transform: scale(1);
    }
  }
  
  &--large {
    switch {
      transform: scale(1.15);
    }
  }
  
  &--disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
}
</style>
