<template>
  <view class="navbar" :class="{ 'navbar-transparent': transparent }">
    <!-- 状态栏占位 -->
    <view class="navbar-status-bar" :style="{ height: statusBarHeight + 'px' }"></view>
    
    <!-- 导航栏内容 -->
    <view class="navbar-content">
      <!-- 左侧区域 -->
      <view class="navbar-left" @tap="handleBack">
        <slot name="left">
          <view v-if="showBack" class="navbar-back">
            <text class="navbar-back-icon">‹</text>
            <text v-if="backText" class="navbar-back-text">{{ backText }}</text>
          </view>
        </slot>
      </view>
      
      <!-- 标题区域 -->
      <view class="navbar-title">
        <slot name="title">
          <text class="navbar-title-text">{{ title }}</text>
        </slot>
      </view>
      
      <!-- 右侧区域 -->
      <view class="navbar-right">
        <slot name="right"></slot>
      </view>
    </view>
  </view>
  
  <!-- 占位元素，防止内容被导航栏遮挡 -->
  <view v-if="placeholder" class="navbar-placeholder" :style="{ height: navbarHeight + 'px' }"></view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

const props = withDefaults(defineProps<{
  title?: string
  showBack?: boolean
  backText?: string
  transparent?: boolean
  placeholder?: boolean
  fixed?: boolean
}>(), {
  title: '',
  showBack: true,
  backText: '',
  transparent: false,
  placeholder: true,
  fixed: true,
})

const emit = defineEmits<{
  (e: 'back'): void
}>()

// 状态栏高度
const statusBarHeight = ref(20)
// 导航栏内容高度
const navContentHeight = 44
// 总高度
const navbarHeight = computed(() => statusBarHeight.value + navContentHeight)

onMounted(() => {
  // 获取系统信息
  const systemInfo = uni.getSystemInfoSync()
  statusBarHeight.value = systemInfo.statusBarHeight || 20
})

const handleBack = () => {
  if (!props.showBack) return
  
  emit('back')
  
  // 默认返回上一页
  const pages = getCurrentPages()
  if (pages.length > 1) {
    uni.navigateBack()
  } else {
    uni.switchTab({ url: '/pages/index/index' })
  }
}
</script>

<style lang="scss" scoped>
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 999;
  background: var(--bg-card, #FFFFFF);
  
  &-transparent {
    background: transparent;
  }
  
  &-status-bar {
    width: 100%;
  }
  
  &-content {
    display: flex;
    align-items: center;
    height: 88rpx;
    padding: 0 24rpx;
  }
  
  &-left {
    flex-shrink: 0;
    min-width: 120rpx;
  }
  
  &-back {
    display: flex;
    align-items: center;
    
    &-icon {
      font-size: 48rpx;
      font-weight: 300;
      color: var(--text-primary, #1A1A1A);
      line-height: 1;
    }
    
    &-text {
      font-size: 28rpx;
      color: var(--text-primary, #1A1A1A);
      margin-left: 8rpx;
    }
  }
  
  &-title {
    flex: 1;
    text-align: center;
    overflow: hidden;
    
    &-text {
      font-size: 34rpx;
      font-weight: 600;
      color: var(--text-primary, #1A1A1A);
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
  
  &-right {
    flex-shrink: 0;
    min-width: 120rpx;
    display: flex;
    justify-content: flex-end;
  }
}

.navbar-placeholder {
  width: 100%;
}
</style>
