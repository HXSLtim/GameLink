<template>
  <view v-if="visible" class="filter-mask" :class="{ 'filter-mask--pc': isPC }" @tap="handleClose">
    <view
      class="filter-panel"
      :class="{ 'filter-panel--show': visible, 'filter-panel--pc': isPC }"
      @tap.stop
    >
      <!-- 头部 -->
      <view class="filter-header">
        <text class="filter-action" @tap="handleReset">{{ resetText }}</text>
        <text class="filter-title">{{ title }}</text>
        <text class="filter-action filter-action--primary" @tap="handleConfirm">{{ confirmText }}</text>
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
      
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useDevice } from '@/composables/useDevice'
import type { FilterOption, FilterSection, FilterValues } from '@/types/filter'

interface Props {
  visible: boolean
  sections: FilterSection[]
  modelValue: FilterValues
  title?: string
  resetText?: string
  confirmText?: string
}

const props = withDefaults(defineProps<Props>(), {
  title: '筛选条件',
  resetText: '清空',
  confirmText: '完成',
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
const { isPC } = useDevice()

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
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: flex-end;
  animation: fadeIn 0.2s ease-out;
}

.filter-mask--pc {
  align-items: stretch;
  justify-content: flex-end;
}

.filter-mask--pc {
  align-items: stretch;
  justify-content: flex-end;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.filter-panel {
  width: 100%;
  max-height: 75vh;
  background: var(--color-bg-card);
  border-radius: var(--radius-md) var(--radius-md) 0 0;
  border-top: 1rpx solid var(--color-border);
  padding-bottom: calc(env(safe-area-inset-bottom) + var(--spacing-sm));
  animation: slideUp 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
}

.filter-panel--pc {
  width: 360px;
  max-height: 100%;
  height: 100%;
  border-radius: 0;
  border-top: none;
  border-left: 1rpx solid var(--color-border);
  padding-bottom: var(--spacing-lg);
  animation: slideInRight 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes slideUp {
  from { transform: translateY(100%); }
  to { transform: translateY(0); }
}

@keyframes slideInRight {
  from { transform: translateX(100%); }
  to { transform: translateX(0); }
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-lg) var(--spacing-md);
  border-bottom: 1rpx solid var(--color-border);
  flex-shrink: 0;
}

.filter-title {
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
  text-align: center;
  flex: 1;
}

.filter-action {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  font-weight: 500;
  padding: 6rpx var(--spacing-sm);
  border-radius: var(--radius-md);
  border: 1rpx solid var(--color-border);
  background: var(--color-bg-secondary);
  cursor: pointer;
  @include press-effect;
  
  &:active {
    background: var(--color-bg-secondary);
  }
}

.filter-action--primary {
  color: var(--color-primary);
  border-color: var(--color-primary);
  background: var(--color-bg-card);
}

.filter-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

.filter-section {
  padding: var(--spacing-md) var(--spacing-md) var(--spacing-lg);
}

.section-title {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  margin-bottom: var(--spacing-md);
  font-weight: 500;
}

.filter-options {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm);

  @include desktop {
    gap: var(--spacing-md);
  }
}

.filter-option {
  padding: 8rpx var(--spacing-md);
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  font-size: var(--font-sm);
  color: var(--color-text);
  border: 1rpx solid var(--color-border);
  transition: all 0.2s;
  cursor: pointer;
  @include press-effect;

  &:active {
    background: var(--color-bg-secondary);
  }

  &--active {
    background: var(--color-primary);
    color: #fff;
    border-color: var(--color-primary);
    font-weight: 600;
  }
}
</style>
