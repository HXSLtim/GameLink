<template>
  <view class="lazy-image" :style="containerStyle">
    <!-- 加载中占位 -->
    <view v-if="loading" class="placeholder" :class="{ 'skeleton-animation': showSkeleton }">
      <view v-if="!showSkeleton" class="placeholder-icon"><uv-icon name="camera" size="32" color="var(--color-text-placeholder)" /></view>
    </view>
    
    <!-- 加载失败 -->
    <view v-else-if="error" class="error-state" @tap="retry">
      <view class="error-icon"><uv-icon name="error-circle" size="32" color="var(--color-text-placeholder)" /></view>
      <text v-if="showRetry" class="retry-text">点击重试</text>
    </view>
    
    <!-- 图片 -->
    <image
      v-show="!loading && !error"
      class="image"
      :src="currentSrc"
      :mode="mode"
      :lazy-load="lazyLoad"
      :webp="webp"
      :show-menu-by-longpress="showMenuByLongpress"
      @load="handleLoad"
      @error="handleError"
    />
  </view>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'

interface Props {
  src: string
  mode?: 'scaleToFill' | 'aspectFit' | 'aspectFill' | 'widthFix' | 'heightFix' | 'top' | 'bottom' | 'center' | 'left' | 'right'
  width?: string | number
  height?: string | number
  radius?: string | number
  lazyLoad?: boolean
  webp?: boolean
  showMenuByLongpress?: boolean
  placeholder?: string
  errorImage?: string
  showSkeleton?: boolean
  showRetry?: boolean
  retryCount?: number
}

const props = withDefaults(defineProps<Props>(), {
  mode: 'aspectFill',
  lazyLoad: true,
  webp: true,
  showMenuByLongpress: false,
  showSkeleton: true,
  showRetry: true,
  retryCount: 3,
})

const emit = defineEmits<{
  load: [e: Event]
  error: [e: Event]
}>()

// 状态
const loading = ref(true)
const error = ref(false)
const retried = ref(0)
const currentSrc = ref('')

// 容器样式
const containerStyle = computed(() => {
  const style: Record<string, string> = {}
  
  if (props.width) {
    style.width = typeof props.width === 'number' ? `${props.width}rpx` : props.width
  }
  if (props.height) {
    style.height = typeof props.height === 'number' ? `${props.height}rpx` : props.height
  }
  if (props.radius) {
    style.borderRadius = typeof props.radius === 'number' ? `${props.radius}rpx` : props.radius
  }
  
  return style
})

// 加载成功
const handleLoad = (e: Event) => {
  loading.value = false
  error.value = false
  emit('load', e)
}

// 加载失败
const handleError = (e: Event) => {
  loading.value = false
  
  // 尝试重试
  if (retried.value < props.retryCount) {
    retried.value++
    loading.value = true
    // 添加时间戳强制重新加载
    currentSrc.value = `${props.src}${props.src.includes('?') ? '&' : '?'}_t=${Date.now()}`
    return
  }
  
  // 显示错误图片或错误状态
  if (props.errorImage) {
    currentSrc.value = props.errorImage
    loading.value = true
  } else {
    error.value = true
  }
  
  emit('error', e)
}

// 手动重试
const retry = () => {
  if (!props.showRetry) return
  retried.value = 0
  error.value = false
  loading.value = true
  currentSrc.value = `${props.src}${props.src.includes('?') ? '&' : '?'}_t=${Date.now()}`
}

// 监听 src 变化
watch(() => props.src, (newSrc) => {
  if (newSrc) {
    loading.value = true
    error.value = false
    retried.value = 0
    currentSrc.value = newSrc
  }
}, { immediate: true })
</script>

<style lang="scss" scoped>
.lazy-image {
  position: relative;
  overflow: hidden;
  background: var(--color-bg-secondary);
  
  .placeholder,
  .error-state,
  .image {
    width: 100%;
    height: 100%;
  }
  
  .placeholder {
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--color-bg-secondary);
    
    &.skeleton-animation {
      background: var(--color-bg-secondary);
      animation: shimmer 1.5s ease-in-out infinite;
    }
    
    .placeholder-icon {
      font-size: var(--font-xl);
      opacity: 0.5;
    }
  }
  
  .error-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: var(--spacing-xs);
    background: var(--color-bg-secondary);
    cursor: pointer;
    @include press-effect;
    
    .error-icon {
      font-size: var(--font-xl);
    }
    
    .retry-text {
      font-size: var(--font-sm);
      color: var(--color-primary);
    }
  }
  
  .image {
    display: block;
  }
}

@keyframes shimmer {
  0% {
    background-color: var(--color-bg-secondary);
  }
  50% {
    background-color: var(--color-bg-card);
  }
  100% {
    background-color: var(--color-bg-secondary);
  }
}
</style>
