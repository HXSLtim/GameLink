<template>
  <view v-if="visible" class="filter-mask" @tap="handleClose">
    <view class="filter-panel" :class="{ 'filter-panel--show': visible }" @tap.stop>
      <!-- 头部 -->
      <view class="filter-header">
        <text class="filter-title">{{ title }}</text>
        <text class="filter-reset" @tap="handleReset">{{ resetText }}</text>
      </view>
      
      <!-- 筛选内容 -->
      <scroll-view class="filter-content" scroll-y>
        <view v-for="section in sections" :key="section.key" class="filter-section">
          <text class="section-title">{{ section.label }}</text>
          <view class="filter-options">
            <view
              v-for="option in section.options"
              :key="String(option.value)"
              class="filter-option"
              :class="{ 'filter-option--active': isSelected(section.key, option.value) }"
              @tap="handleSelect(section.key, option.value, section.multiple)"
            >
              {{ option.label }}
            </view>
          </view>
        </view>
        
        <!-- 自定义内容插槽 -->
        <slot></slot>
      </scroll-view>
      
      <!-- 底部按钮 -->
      <view class="filter-footer">
        <view class="filter-btn filter-btn--cancel" @tap="handleClose">
          {{ cancelText }}
        </view>
        <view class="filter-btn filter-btn--confirm" @tap="handleConfirm">
          {{ confirmText }}
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

export interface FilterOption {
  label: string
  value: string | number | boolean
}

export interface FilterSection {
  key: string
  label: string
  options: FilterOption[]
  multiple?: boolean
}

export type FilterValues = Record<string, string | number | boolean | Array<string | number | boolean>>

interface Props {
  visible: boolean
  sections: FilterSection[]
  modelValue: FilterValues
  title?: string
  resetText?: string
  cancelText?: string
  confirmText?: string
}

const props = withDefaults(defineProps<Props>(), {
  title: '筛选条件',
  resetText: '重置',
  cancelText: '取消',
  confirmText: '确定',
})

const emit = defineEmits<{
  'update:visible': [visible: boolean]
  'update:modelValue': [values: FilterValues]
  apply: [values: FilterValues]
  reset: []
  close: []
}>()

// 内部状态（用于编辑）
const localValues = ref<FilterValues>({ ...props.modelValue })

// 监听外部值变化
watch(() => props.modelValue, (newVal) => {
  localValues.value = { ...newVal }
}, { deep: true })

// 监听可见性变化，重置本地值
watch(() => props.visible, (visible) => {
  if (visible) {
    localValues.value = { ...props.modelValue }
  }
})

const isSelected = (key: string, value: string | number | boolean): boolean => {
  const current = localValues.value[key]
  if (Array.isArray(current)) {
    return current.includes(value)
  }
  return current === value
}

const handleSelect = (key: string, value: string | number | boolean, multiple?: boolean) => {
  if (multiple) {
    const current = localValues.value[key]
    const arr = Array.isArray(current) ? [...current] : []
    const index = arr.indexOf(value)
    if (index > -1) {
      arr.splice(index, 1)
    } else {
      arr.push(value)
    }
    localValues.value[key] = arr
  } else {
    // 单选：点击已选中的则取消
    localValues.value[key] = localValues.value[key] === value ? '' : value
  }
}

const handleReset = () => {
  const resetValues: FilterValues = {}
  props.sections.forEach(section => {
    resetValues[section.key] = section.multiple ? [] : ''
  })
  localValues.value = resetValues
  emit('reset')
}

const handleClose = () => {
  emit('update:visible', false)
  emit('close')
}

const handleConfirm = () => {
  emit('update:modelValue', { ...localValues.value })
  emit('apply', { ...localValues.value })
  emit('update:visible', false)
}
</script>

<style lang="scss" scoped>
.filter-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1000;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(8rpx);
  display: flex;
  align-items: flex-end;
  animation: fadeIn 0.2s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.filter-panel {
  width: 100%;
  max-height: 75vh;
  background: var(--color-bg-card);
  border-radius: 40rpx 40rpx 0 0;
  padding-bottom: calc(env(safe-area-inset-bottom) + 20rpx);
  animation: slideUp 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
}

@keyframes slideUp {
  from { transform: translateY(100%); }
  to { transform: translateY(0); }
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 36rpx 32rpx;
  border-bottom: 1rpx solid var(--color-border);
  flex-shrink: 0;
}

.filter-title {
  font-size: 36rpx;
  font-weight: 700;
  color: var(--color-text);
}

.filter-reset {
  font-size: 28rpx;
  color: var(--color-primary);
  font-weight: 500;
  padding: 8rpx 20rpx;
  border-radius: 16rpx;
  
  &:active {
    background: rgba(0, 210, 106, 0.1);
  }
}

.filter-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

.filter-section {
  padding: 28rpx 32rpx;
}

.section-title {
  font-size: 28rpx;
  color: var(--color-text-secondary);
  margin-bottom: 24rpx;
  font-weight: 500;
}

.filter-options {
  display: flex;
  flex-wrap: wrap;
  gap: 20rpx;
}

.filter-option {
  padding: 18rpx 36rpx;
  background: var(--color-bg-secondary);
  border-radius: 36rpx;
  font-size: 28rpx;
  color: var(--color-text);
  border: 2rpx solid transparent;
  transition: all 0.2s;
  
  &:active {
    transform: scale(0.95);
  }
  
  &--active {
    background: rgba(0, 210, 106, 0.1);
    color: var(--color-primary);
    border-color: var(--color-primary);
    font-weight: 600;
  }
}

.filter-footer {
  display: flex;
  gap: 24rpx;
  padding: 28rpx 32rpx;
  border-top: 1rpx solid var(--color-border);
  flex-shrink: 0;
}

.filter-btn {
  flex: 1;
  height: 96rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 48rpx;
  font-size: 32rpx;
  font-weight: 600;
  transition: all 0.2s;
  
  &:active {
    transform: scale(0.98);
  }
  
  &--cancel {
    background: var(--color-bg-secondary);
    color: var(--color-text);
    border: 2rpx solid var(--color-border);
  }
  
  &--confirm {
    background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-primary-light, #4ADE80) 100%);
    color: #FFFFFF;
    box-shadow: 0 6rpx 20rpx rgba(0, 210, 106, 0.35);
  }
}
</style>
