<template>
  <view class="form-item" :class="{ clickable }" @tap="handleClick">
    <view class="form-label-wrap">
      <text class="form-label">{{ label }}</text>
      <text v-if="required" class="form-required">*</text>
    </view>
    
    <!-- 输入框 -->
    <template v-if="type === 'input'">
      <GlInput
        class="form-input"
        :model-value="modelValue ?? ''"
        :placeholder="placeholder"
        :maxlength="maxlength"
        :disabled="disabled"
        size="small"
        variant="plain"
        align="right"
        @update:modelValue="(value) => $emit('update:modelValue', value)"
      />
    </template>
    
    <!-- 选择器/只读 -->
    <template v-else>
      <view class="form-value-wrap">
        <text class="form-value" :class="{ placeholder: !displayValue }">
          {{ displayValue || placeholder }}
        </text>
        <slot name="extra"></slot>
        <uv-icon v-if="clickable && !$slots.extra" name="arrow-right" size="14" color="var(--color-text-secondary)"></uv-icon>
      </view>
    </template>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlInput from '@/components/gl/Input/index.vue'

interface Props {
  label: string
  type?: 'input' | 'picker' | 'readonly'
  modelValue?: string | number
  displayValue?: string
  placeholder?: string
  maxlength?: number
  clickable?: boolean
  required?: boolean
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  type: 'picker',
  placeholder: '请选择',
  clickable: true,
  required: false,
  disabled: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  click: []
}>()

const handleClick = () => {
  if (props.type !== 'input' && props.clickable) {
    emit('click')
  }
}
</script>

<style lang="scss" scoped>
.form-item {
  display: flex;
  align-items: center;
  padding: var(--spacing-md) 0;
  border-bottom: 1rpx solid var(--color-border);
  transition: background 0.2s ease;
  
  &:last-child {
    border-bottom: none;
  }
  
  &.clickable {
    cursor: pointer;
    
    &:active {
      background: var(--color-bg-secondary);
    }
  }
}

.form-label-wrap {
  display: flex;
  align-items: center;
  width: 160rpx;
  flex-shrink: 0;
}

.form-label {
  font-size: var(--font-base);
  color: var(--color-text);
}

.form-required {
  color: var(--color-error);
  font-size: var(--font-base);
  margin-left: 4rpx;
}

.form-input {
  flex: 1;
  min-width: 0;
  
  :deep(.gl-input__field) {
    font-size: var(--font-base);
    color: var(--color-text);
  }
}

.form-value-wrap {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8rpx;
}

.form-value {
  font-size: var(--font-base);
  color: var(--color-text);
  
  &.placeholder {
    color: var(--color-text-placeholder);
  }
}
</style>
