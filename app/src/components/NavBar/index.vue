<template>
  <view class="navbar" :class="{ 'navbar-transparent': transparent, 'navbar--fixed': fixed }">
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
  <view class="navbar-title" :class="`navbar-title--${resolvedTitleAlign}`">
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
import { ref, computed, onMounted, watch } from 'vue'
import { useDevice } from '@/composables/useDevice'

const props = withDefaults(defineProps<{
  title?: string
  showBack?: boolean
  backText?: string
  transparent?: boolean
  placeholder?: boolean
  fixed?: boolean
  titleAlign?: 'left' | 'center'
}>(), {
  title: '',
  showBack: true,
  backText: '',
  transparent: false,
  placeholder: true,
  fixed: true,
  titleAlign: 'left',
})

const emit = defineEmits<{
  (e: 'back'): void
}>()

const { isPC } = useDevice()

// 状态栏高度
const statusBarHeight = ref(20)
const defaultStatusBarHeight = ref(20)
// 导航栏内容高度
const navContentHeight = 48
// 总高度
const navbarHeight = computed(() => statusBarHeight.value + navContentHeight)
const resolvedTitleAlign = computed(() => (isPC.value ? 'center' : props.titleAlign))

onMounted(() => {
  // 获取系统信息
  const systemInfo = uni.getSystemInfoSync()
  defaultStatusBarHeight.value = systemInfo.statusBarHeight || 20
  statusBarHeight.value = defaultStatusBarHeight.value
})

watch(isPC, (value) => {
  statusBarHeight.value = value ? 0 : defaultStatusBarHeight.value
}, { immediate: true })

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
  position: relative;
  width: 100%;
  background: var(--color-bg);
  border-bottom: 1rpx solid var(--color-border);
  box-shadow: none;

  // PC 端：用柔和阴影替代硬边框
  @include desktop {
    width: 100%;
    border-bottom: none;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  }
  
  &-transparent {
    background: transparent;
    border-bottom: none;
    box-shadow: none;
  }

  &--fixed {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    z-index: 999;

    @include desktop {
      left: var(--sidebar-width);
      right: 0;
      width: calc(100% - var(--sidebar-width));
    }
  }
  
  &-status-bar {
    width: 100%;

    @include desktop {
      height: 0 !important;
      display: none;
    }
  }
  
  &-content {
    display: flex;
    align-items: center;
    height: 96rpx;
    padding: 0 var(--spacing-md);

    @include desktop {
      height: 52px;
      padding: 0 24px;
    }
  }
  
  &-left {
    flex-shrink: 0;
    min-width: 88rpx;

    @include desktop {
      min-width: 64px;
    }
  }
  
  &-back {
    display: flex;
    align-items: center;
    gap: var(--spacing-xs);
    @include press-effect;

    // PC 端：hover 反馈
    @include hover-supported {
      &:hover {
        opacity: 0.7;
      }
    }
    
    &-icon {
      font-size: var(--font-xl);
      font-weight: 400;
      color: var(--color-text);
      line-height: 1;
    }
    
    &-text {
      font-size: var(--font-sm);
      color: var(--color-text);
    }
  }
  
  &-title {
    flex: 1;
    overflow: hidden;
    min-width: 0;
    
    &--left {
      text-align: left;
      padding-left: var(--spacing-sm);

      @include desktop {
        padding-left: 0;
      }
    }
    
    &--center {
      text-align: center;
    }
    
    &-text {
      font-size: var(--font-xl);
      font-weight: 600;
      color: var(--color-text);
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;

      @include desktop {
        font-size: var(--font-md);
        font-weight: 700;
        letter-spacing: -0.2px;
      }
    }
  }
  
  &-right {
    flex-shrink: 0;
    min-width: 88rpx;
    display: flex;
    justify-content: flex-end;
    gap: 8px;

    @include desktop {
      min-width: 64px;
    }
  }
}

.navbar-placeholder {
  width: 100%;
}
</style>
