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
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px); /* 磨砂玻璃效果 */
  border-bottom: 1rpx solid rgba(0, 0, 0, 0.05);
  box-shadow: $gl-shadow-sm;
  transition: all 0.3s ease;

  // PC 端适配
  @include desktop {
    width: 100%;
    border-bottom: none;
    box-shadow: $gl-shadow-sm;
    background: #fff;
  }

  &-transparent {
    background: transparent;
    border-bottom: none;
    box-shadow: none;
    backdrop-filter: none;

    .navbar-title-text {
      opacity: 0; /* 透明模式下默认隐藏标题，滚动后显示 */
    }
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
    height: 88rpx; /* 标准高度 */
    padding: 0 $gl-spacing-md;

    @include desktop {
      height: 64px;
      padding: 0 32px;
    }
  }

  &-left {
    flex-shrink: 0;
    min-width: 80rpx;
    display: flex;
    align-items: center;

    @include desktop {
      min-width: 64px;
    }
  }

  &-back {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 64rpx;
    height: 64rpx;
    border-radius: $gl-radius-circle;
    @include press-effect;

    /* 增加触摸区域背景，提升点击体验 */
    &:active {
      background: rgba(0, 0, 0, 0.05);
    }

    // PC 端：hover 反馈
    @include hover-supported {
      &:hover {
        background: rgba(0, 0, 0, 0.05);
      }
    }

    &-icon {
      font-size: 40rpx;
      font-weight: 400;
      color: $gl-text-main;
      line-height: 1;
      margin-top: -4rpx; // 视觉修正
    }

    &-text {
      font-size: 28rpx;
      color: $gl-text-main;
      margin-left: 4rpx;
    }
  }

  &-title {
    flex: 1;
    overflow: hidden;
    min-width: 0;
    display: flex;
    align-items: center;

    &--left {
      justify-content: flex-start;
      padding-left: $gl-spacing-xs;

      @include desktop {
        padding-left: 0;
      }
    }

    &--center {
      justify-content: center;
    }

    &-text {
      font-size: 34rpx;
      font-weight: 600;
      color: $gl-text-main;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      transition: opacity 0.3s;

      @include desktop {
        font-size: 18px;
        font-weight: 700;
        letter-spacing: -0.2px;
      }
    }
  }

  &-right {
    flex-shrink: 0;
    min-width: 80rpx;
    display: flex;
    justify-content: flex-end;
    align-items: center;
    gap: $gl-spacing-sm;

    @include desktop {
      min-width: 64px;
    }
  }
}

.navbar-placeholder {
  width: 100%;
}
</style>