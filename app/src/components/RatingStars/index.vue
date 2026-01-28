<template>
  <view class="rating-stars" :class="[`size-${size}`]">
    <view 
      v-for="i in 5" 
      :key="i"
      class="star"
      :class="{ active: i <= Math.round(modelValue), half: isHalfStar(i) }"
      @click="handleClick(i)"
    >
      <text>{{ i <= Math.round(modelValue) ? '★' : '☆' }}</text>
    </view>
    <text v-if="showValue" class="rating-value">{{ modelValue.toFixed(1) }}</text>
  </view>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  modelValue: number
  readonly?: boolean
  size?: 'small' | 'medium' | 'large'
  showValue?: boolean
}>(), {
  modelValue: 0,
  readonly: false,
  size: 'medium',
  showValue: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: number]
}>()

function isHalfStar(index: number): boolean {
  const value = props.modelValue
  return index === Math.ceil(value) && value % 1 !== 0
}

function handleClick(index: number) {
  if (props.readonly) return
  emit('update:modelValue', index)
}
</script>

<style lang="scss" scoped>
.rating-stars {
  display: inline-flex;
  align-items: center;
  gap: 4rpx;
  
  .star {
    color: #D1D5DB;
    transition: color 0.2s;
    
    &.active {
      color: #FBBF24;
    }
  }
  
  .rating-value {
    margin-left: 8rpx;
    color: var(--color-text-secondary);
  }
}

// 尺寸
.size-small {
  .star text { font-size: 24rpx; }
  .rating-value { font-size: 22rpx; }
}

.size-medium {
  .star text { font-size: 32rpx; }
  .rating-value { font-size: 26rpx; }
}

.size-large {
  .star text { font-size: 48rpx; }
  .rating-value { font-size: 32rpx; }
}
</style>
