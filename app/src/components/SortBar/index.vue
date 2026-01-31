<template>
  <view class="sort-bar" :class="{ 'sort-bar--bordered': bordered }">
    <view
      v-for="option in options"
      :key="option.value"
      class="sort-item"
      :class="{ 'sort-item--active': modelValue === option.value }"
      @tap="handleChange(option.value)"
    >
      <text>{{ option.label }}</text>
      <!-- 排序方向指示器 -->
      <view v-if="showDirection && modelValue === option.value" class="sort-direction">
        <uv-icon 
          :name="sortDirection === 'asc' ? 'arrow-up' : 'arrow-down'" 
          size="12" 
          color="var(--color-primary)"
        ></uv-icon>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
export interface SortOption {
  label: string
  value: string
}

interface Props {
  options: SortOption[]
  modelValue: string
  bordered?: boolean
  showDirection?: boolean
  sortDirection?: 'asc' | 'desc'
}

const props = withDefaults(defineProps<Props>(), {
  bordered: true,
  showDirection: false,
  sortDirection: 'desc',
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  change: [value: string]
}>()

const handleChange = (value: string) => {
  emit('update:modelValue', value)
  emit('change', value)
}
</script>

<style lang="scss" scoped>
.sort-bar {
  display: flex;
  padding: 16rpx 24rpx;
  background: var(--color-bg-card);
  
  &--bordered {
    margin-bottom: 12rpx;
    border-bottom: 1rpx solid var(--color-border);
  }
}

.sort-item {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4rpx;
  text-align: center;
  font-size: 28rpx;
  color: var(--color-text-secondary);
  padding: 16rpx 0;
  position: relative;
  transition: all 0.2s;
  font-weight: 500;
  
  &:active {
    color: var(--color-primary);
  }
  
  &--active {
    color: var(--color-primary);
    font-weight: 700;
    
    &::after {
      content: '';
      position: absolute;
      bottom: 0;
      left: 50%;
      transform: translateX(-50%);
      width: 48rpx;
      height: 6rpx;
      background: linear-gradient(90deg, var(--color-primary) 0%, var(--color-primary-light, #4ADE80) 100%);
      border-radius: 3rpx;
    }
  }
}

.sort-direction {
  margin-left: 4rpx;
}
</style>
