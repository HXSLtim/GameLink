<template>
  <view class="search-bar">
    <view class="search-input-wrap" :class="{ focused: isFocused }">
      <uv-icon name="search" size="18" color="var(--color-text-secondary)"></uv-icon>
      <input
        :value="modelValue"
        class="search-input"
        :placeholder="placeholder"
        :disabled="disabled"
        @input="handleInput"
        @confirm="handleConfirm"
        @focus="handleFocus"
        @blur="handleBlur"
      />
      <view v-if="modelValue && clearable" class="clear-btn" @tap.stop="handleClear">
        <uv-icon name="close-circle-fill" size="16" color="var(--color-text-placeholder)"></uv-icon>
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

interface Props {
  modelValue: string
  placeholder?: string
  showFilter?: boolean
  filterText?: string
  clearable?: boolean
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '搜索',
  showFilter: false,
  filterText: '筛选',
  clearable: true,
  disabled: false,
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

const handleInput = (e: any) => {
  emit('update:modelValue', e.detail.value)
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
  gap: 12rpx;
  padding: 12rpx 16rpx;
  background: var(--color-bg-card);
  border-bottom: 1rpx solid var(--color-border);
}

.search-input-wrap {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 12rpx;
  padding: 0 20rpx;
  height: 68rpx;
  background: var(--color-bg-secondary);
  border-radius: 34rpx;
  border: 1rpx solid transparent;
  transition: all 0.2s;
  
  &.focused {
    border-color: var(--color-primary);
    background: var(--color-bg-card);
  }
}

.search-input {
  flex: 1;
  font-size: 26rpx;
  color: var(--color-text);
  
  &::placeholder {
    color: var(--color-text-placeholder);
  }
}

.clear-btn {
  padding: 8rpx;
  
  &:active {
    opacity: 0.7;
  }
}

.filter-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6rpx;
  padding: 0 20rpx;
  height: 68rpx;
  background: var(--color-bg-secondary);
  border-radius: 34rpx;
  font-size: 26rpx;
  color: var(--color-text);
  border: 1rpx solid transparent;
  transition: all 0.2s;
  
  &:active {
    border-color: var(--color-primary);
    background: rgba(0, 210, 106, 0.1);
    color: var(--color-primary);
  }
}
</style>
