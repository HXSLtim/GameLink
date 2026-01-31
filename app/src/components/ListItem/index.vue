<template>
  <view 
    class="list-item" 
    :class="{ 
      'list-item--animated': animated,
      'list-item--clickable': clickable 
    }"
    :style="animationStyle"
    @tap="handleClick"
  >
    <slot></slot>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  index?: number
  animated?: boolean
  clickable?: boolean
  animationDelay?: number
}

const props = withDefaults(defineProps<Props>(), {
  index: 0,
  animated: true,
  clickable: true,
  animationDelay: 0.03,
})

const emit = defineEmits<{
  click: []
}>()

const animationStyle = computed(() => {
  if (!props.animated) return {}
  return {
    animationDelay: `${props.index * props.animationDelay}s`,
  }
})

const handleClick = () => {
  if (props.clickable) {
    emit('click')
  }
}
</script>

<style lang="scss" scoped>
.list-item {
  margin-bottom: 8rpx;
  
  &--animated {
    animation: fadeSlideUp 0.3s ease-out;
    animation-fill-mode: both;
  }
  
  &--clickable {
    &:active {
      transform: scale(0.98);
      opacity: 0.9;
    }
  }
}

@keyframes fadeSlideUp {
  from { 
    opacity: 0; 
    transform: translateY(10rpx); 
  }
  to { 
    opacity: 1; 
    transform: translateY(0); 
  }
}
</style>
