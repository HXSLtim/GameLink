<template>
  <view
    class="gl-input"
    :class="[
      `gl-input--${size}`,
      `gl-input--${variant}`,
      {
        'gl-input--disabled': disabled,
        'gl-input--focused': focused,
        'gl-input--textarea': isTextarea,
      }
    ]"
    :style="customStyle"
  >
    <uv-icon v-if="prefixIcon" :name="prefixIcon" :size="iconSize" color="var(--color-text-secondary)" />
    
    <input
      v-if="!isTextarea"
      :value="modelValue"
      class="gl-input__field"
      :type="nativeType"
      :placeholder="placeholder"
      :maxlength="maxlength"
      :disabled="disabled"
      :readonly="readonly"
      :style="{ textAlign: align }"
      @input="handleInput"
      @focus="handleFocus"
      @blur="handleBlur"
      @confirm="handleConfirm"
    />
    <textarea
      v-else
      :value="modelValue"
      class="gl-input__textarea"
      :placeholder="placeholder"
      :maxlength="maxlength"
      :disabled="disabled"
      :readonly="readonly"
      :auto-height="autoHeight"
      :show-confirm-bar="showConfirmBar"
      :adjust-position="adjustPosition"
      :style="{ textAlign: align }"
      @input="handleInput"
      @focus="handleFocus"
      @blur="handleBlur"
      @confirm="handleConfirm"
    />
    
    <view v-if="clearable && modelValue" class="gl-input__clear" @tap="handleClear">
      <uv-icon name="close-circle-fill" :size="iconSize" color="var(--color-text-placeholder)" />
    </view>
    <view
      v-else-if="suffixIcon"
      class="gl-input__suffix"
      :class="{ 'gl-input__suffix--clickable': suffixClickable }"
      @tap="handleSuffixClick"
    >
      <uv-icon :name="suffixIcon" :size="iconSize" color="var(--color-text-secondary)" />
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

type InputSize = 'small' | 'medium' | 'large'
type InputAlign = 'left' | 'center' | 'right'

interface Props {
  modelValue?: string | number
  type?: 'text' | 'password' | 'number' | 'digit' | 'textarea'
  variant?: 'filled' | 'plain'
  placeholder?: string
  maxlength?: number
  autoHeight?: boolean
  showConfirmBar?: boolean
  adjustPosition?: boolean
  disabled?: boolean
  readonly?: boolean
  clearable?: boolean
  size?: InputSize
  align?: InputAlign
  prefixIcon?: string
  suffixIcon?: string
  suffixClickable?: boolean
  customStyle?: string | Record<string, string>
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  type: 'text',
  variant: 'filled',
  placeholder: '请输入',
  autoHeight: true,
  showConfirmBar: true,
  adjustPosition: true,
  disabled: false,
  readonly: false,
  clearable: false,
  suffixClickable: false,
  size: 'large',
  align: 'left',
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  input: [value: string]
  focus: [e: Event]
  blur: [e: Event]
  confirm: [value: string]
  clear: []
  'suffix-click': []
}>()

const focused = ref(false)

const isTextarea = computed(() => props.type === 'textarea')
const nativeType = computed(() => (props.type === 'textarea' ? 'text' : props.type))

const iconSize = computed(() => {
  const sizes = { small: 14, medium: 16, large: 18 }
  return sizes[props.size]
})

const handleInput = (e: any) => {
  const value = e?.detail?.value ?? ''
  emit('update:modelValue', value)
  emit('input', value)
}

const handleFocus = (e: Event) => {
  focused.value = true
  emit('focus', e)
}

const handleBlur = (e: Event) => {
  focused.value = false
  emit('blur', e)
}

const handleConfirm = (e: any) => {
  const value = e?.detail?.value ?? String(props.modelValue ?? '')
  emit('confirm', value)
}

const handleClear = () => {
  emit('update:modelValue', '')
  emit('clear')
}

const handleSuffixClick = () => {
  if (!props.suffixClickable) return
  emit('suffix-click')
}
</script>

<style lang="scss" scoped>
.gl-input {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  width: 100%;
  padding: 0 var(--spacing-md);
  background: var(--color-bg-secondary);
  border: 1rpx solid var(--color-border);
  border-radius: var(--radius-sm);
  transition: all 0.2s ease;
  cursor: text;
  
  &--small {
    height: 64rpx;
    
    .gl-input__field,
    .gl-input__textarea {
      font-size: var(--font-sm);
    }
  }
  
  &--medium {
    height: 80rpx;
    
    .gl-input__field,
    .gl-input__textarea {
      font-size: var(--font-md);
    }
  }
  
  &--large {
    height: 88rpx;
    
    .gl-input__field,
    .gl-input__textarea {
      font-size: var(--font-base);
    }
  }
  
  &--focused {
    border-color: var(--color-primary);
  }
  
  &--plain {
    background: transparent;
    border: none;
    padding: 0;
    height: auto;
    
    &.gl-input--focused {
      border: none;
    }
  }
  
  &--disabled {
    opacity: 0.6;
    cursor: not-allowed;
    pointer-events: none;
  }
  
  &--textarea {
    align-items: flex-start;
    padding: var(--spacing-sm) var(--spacing-md);
    height: auto;
  }
}

.gl-input__field,
.gl-input__textarea {
  flex: 1;
  width: 100%;
  height: 100%;
  background: transparent;
  color: var(--color-text);
  border: none;
  outline: none;
  
  &::placeholder {
    color: var(--color-text-placeholder);
  }
}

.gl-input__textarea {
  min-height: 120rpx;
  line-height: 1.4;
}

.gl-input__clear {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4rpx;
  border-radius: var(--radius-full);
  cursor: pointer;
  @include press-effect;
}

.gl-input__suffix {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4rpx;
  border-radius: var(--radius-full);

  &--clickable {
    cursor: pointer;
    @include press-effect;
  }
}
</style>
