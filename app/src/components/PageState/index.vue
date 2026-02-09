<template>
  <!-- 加载中 -->
  <view v-if="state === 'loading'" class="page-state loading">
    <view class="loading-container">
      <view class="loading-ring"></view>
      <view class="loading-ring loading-ring--delay"></view>
      <uv-icon name="grid-fill" size="32" color="var(--color-primary, #7ACC35)" class="loading-icon"></uv-icon>
    </view>
    <text class="state-text">{{ loadingText }}</text>
  </view>

  <!-- 错误态 -->
  <view v-else-if="state === 'error'" class="page-state error animate-in">
    <view class="state-icon state-icon--error">
      <view class="icon-bg icon-bg--error"></view>
      <uv-icon name="error-circle" size="48" color="var(--color-text-placeholder)"></uv-icon>
    </view>
    <text class="state-title">{{ errorTitle }}</text>
    <text class="state-desc">{{ errorMessage }}</text>
    <view class="state-actions">
      <GlButton type="error" size="small" round icon="reload" @click="handleRetry">
        重新加载
      </GlButton>
    </view>
  </view>

  <!-- 空态 -->
  <GlEmpty
    v-else-if="state === 'empty'"
    class="page-state empty animate-in"
    :icon="emptyIcon"
    :title="emptyTitle"
    :description="emptyDesc"
    :show-action="emptyAction"
    :action-text="emptyActionText"
    @action="handleEmptyAction"
  />

  <!-- 网络断开 -->
  <view v-else-if="state === 'offline'" class="page-state offline animate-in">
    <view class="state-icon state-icon--offline">
      <view class="icon-bg icon-bg--warning"></view>
      <uv-icon name="wifi-off" size="48" color="var(--color-text-placeholder)"></uv-icon>
    </view>
    <text class="state-title">网络连接已断开</text>
    <text class="state-desc">请检查网络设置后重试</text>
    <view class="state-actions">
      <GlButton type="warning" size="small" round icon="reload" @click="handleRetry">
        重新连接
      </GlButton>
    </view>
  </view>

  <!-- 未登录 -->
  <view v-else-if="state === 'login'" class="page-state login animate-in">
    <view class="state-icon state-icon--login">
      <view class="icon-bg icon-bg--primary"></view>
      <uv-icon name="account" size="48" color="var(--color-text-placeholder)"></uv-icon>
    </view>
    <text class="state-title">请先登录</text>
    <text class="state-desc">登录后即可查看更多内容</text>
    <view class="state-actions">
      <GlButton type="primary" size="small" round icon="arrow-right" @click="handleLogin">
        立即登录
      </GlButton>
    </view>
  </view>

  <!-- 正常内容 -->
  <slot v-else></slot>
</template>

<script setup lang="ts">
import GlButton from '@/components/gl/Button/index.vue'
import GlEmpty from '@/components/gl/Empty/index.vue'
import type { PageStateType } from '@/types/page'

interface Props {
  state: PageStateType
  loadingText?: string
  errorTitle?: string
  errorMessage?: string
  emptyTitle?: string
  emptyDesc?: string
  emptyIcon?: string
  emptyAction?: boolean
  emptyActionText?: string
}

const props = withDefaults(defineProps<Props>(), {
  state: 'content',
  loadingText: '加载中...',
  errorTitle: '加载失败',
  errorMessage: '请检查网络连接后重试',
  emptyTitle: '暂无数据',
  emptyDesc: '',
  emptyIcon: 'empty-data',
  emptyAction: false,
  emptyActionText: '去看看',
})

const emit = defineEmits<{
  retry: []
  emptyAction: []
}>()

const handleRetry = () => {
  emit('retry')
}

const handleEmptyAction = () => {
  emit('emptyAction')
}

const handleLogin = () => {
  uni.navigateTo({ url: '/pages/auth/login/index' })
}
</script>

<style lang="scss" scoped>
.page-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-xl) var(--spacing-lg);
  min-height: 400rpx;
}

// 动画入场
.animate-in {
  animation: fadeSlideUp 0.4s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

@keyframes fadeSlideUp {
  from {
    opacity: 0;
    transform: translateY(30rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

// 加载动画
.loading-container {
  position: relative;
  width: 120rpx;
  height: 120rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-ring {
  position: absolute;
  width: 100%;
  height: 100%;
  border: 4rpx solid transparent;
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  
  &--delay {
    width: 80%;
    height: 80%;
    border-top-color: var(--color-primary-light);
    animation-direction: reverse;
    animation-duration: 0.8s;
  }
}

.loading-icon {
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@keyframes pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.6; transform: scale(0.9); }
}

.state-text {
  margin-top: var(--spacing-lg);
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}

// 图标容器
.state-icon {
  position: relative;
  margin-bottom: var(--spacing-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  width: 120rpx;
  height: 120rpx;
  border-radius: var(--radius-full);
  
  &--error .icon-bg {
    background: var(--color-error-tint);
  }
  
  &--empty .icon-bg {
    background: var(--color-bg-secondary);
  }
  
  &--offline .icon-bg {
    background: var(--color-warning-tint);
  }
  
  &--login .icon-bg {
    background: var(--color-primary-tint);
  }
}

.icon-bg {
  position: absolute;
  width: 100%;
  height: 100%;
  border-radius: var(--radius-full);
}

.state-title {
  font-size: var(--font-lg);
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: var(--spacing-xs);
}

.state-desc {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  text-align: center;
  line-height: 1.6;
  margin-bottom: var(--spacing-md);
}

.state-actions {
  display: flex;
  justify-content: center;
}
</style>
