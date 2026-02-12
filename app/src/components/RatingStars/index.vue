<template>
  <view class="rating-stars" :class="[`size-${size}`]">
    <view 
      v-for="i in 5" 
      :key="i"
      class="star"
      :class="{ active: i <= Math.round(modelValue), half: isHalfStar(i) }"
      @click="handleClick(i)"
    >
      <uv-icon
        :name="i <= Math.round(modelValue) ? 'star-fill' : 'star'"
        :size="starSize"
        :color="i <= Math.round(modelValue) ? 'var(--color-gold)' : 'var(--color-text-disabled)'"
      />
    </view>
    <text v-if="showValue" class="rating-value">{{ modelValue.toFixed(1) }}</text>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  modelValue: number
  readonly?: boolean
  size?: 'small' | 'medium' | 'large' | 'mini'
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

const starSize = computed(() => {
  const map = { mini: 10, small: 12, medium: 16, large: 24 }
  return map[props.size] ?? 16
})

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
    transition: opacity 0.2s;
    cursor: pointer;
    @include press-effect;
  }
  
  .rating-value {
    margin-left: 8rpx;
    color: var(--color-text-secondary);
  }
}

// 尺寸
.size-small {
  .rating-value { font-size: var(--font-xs); }
}

.size-medium {
  .rating-value { font-size: var(--font-sm); }
}

.size-large {
  .rating-value { font-size: var(--font-md); }
}
</style>
