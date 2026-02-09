<template>
  <view class="search-bar">
    <view class="search-input-wrap" :class="{ focused: isFocused }">
      <GlInput
        class="search-input"
        :model-value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :readonly="readonly"
        size="small"
        variant="plain"
        @update:modelValue="handleInput"
        @confirm="handleConfirm"
        @focus="handleFocus"
        @blur="handleBlur"
      />
      <view v-if="modelValue && clearable" class="clear-btn" @tap.stop="handleClear">
        <uv-icon name="close-circle-fill" size="16" color="var(--color-text-placeholder)"></uv-icon>
      </view>
      <view class="search-icon">
        <uv-icon name="search" size="18" color="var(--color-text-secondary)"></uv-icon>
      </view>
    </view>
    
    <!-- 右侧插槽（筛选按钮等） -->
    <slot name="right">
      <view v-if="showFilter" class="filter-btn" @tap="$emit('filter')">
        <uv-icon name="setting" size="16" color="var(--color-text)"></uv-icon>
        <text>{{ filterText }}</text>
      </view>
    </slot>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import GlInput from '@/components/gl/Input/index.vue'

interface Props {
  modelValue: string
  placeholder?: string
  showFilter?: boolean
  filterText?: string
  clearable?: boolean
  disabled?: boolean
  readonly?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '搜索',
  showFilter: false,
  filterText: '筛选',
  clearable: true,
  disabled: false,
  readonly: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  search: [value: string]
  filter: []
  clear: []
  focus: []
  blur: []
}>()

const isFocused = ref(false)

const handleInput = (value: string) => {
  emit('update:modelValue', value)
}

const handleConfirm = () => {
  emit('search', props.modelValue)
}

const handleClear = () => {
  emit('update:modelValue', '')
  emit('clear')
}

const handleFocus = () => {
  isFocused.value = true
  emit('focus')
}

const handleBlur = () => {
  isFocused.value = false
  emit('blur')
}
</script>

<style lang="scss" scoped>
.search-bar {
  display: flex;
  gap: var(--spacing-sm);
  padding: var(--spacing-sm) var(--spacing-md) var(--spacing-xs);
  background: transparent;
}

.search-input-wrap {
  flex: 1;
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: 0 var(--spacing-md);
  height: 64rpx;
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1rpx solid var(--color-border);
  transition: all 0.2s;
  
  &.focused {
    border-color: var(--color-primary);
    background: var(--color-bg-card);
    box-shadow: 0 0 0 3rpx rgba(var(--color-primary-rgb, 122, 204, 53), 0.15);
  }
}

.search-input {
  flex: 1;
  min-width: 0;
  height: 100%;
  
  :deep(.gl-input__field) {
    font-size: var(--font-sm);
    color: var(--color-text);
  }
  
  :deep(.gl-input__field::placeholder) {
    color: var(--color-text-placeholder);
  }
}

.clear-btn {
  padding: var(--spacing-xs);
  @include press-effect;
  
  &:active {
    opacity: 0.7;
  }
}

.search-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.filter-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6rpx;
  padding: 0 var(--spacing-md);
  height: 64rpx;
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  font-size: var(--font-sm);
  color: var(--color-text);
  border: 1rpx solid var(--color-border);
  transition: all 0.2s;
  @include press-effect;
  
  &:active {
    border-color: var(--color-primary);
    background: var(--color-bg-secondary);
    color: var(--color-text);
  }
}
</style>
