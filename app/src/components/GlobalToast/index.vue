<template>
  <view v-if="toast.state.value.visible" class="toast-container" :class="`toast--${toast.state.value.type}`">
    <view class="toast-content">
      <!-- 图标 -->
      <view class="toast-icon">
        <view v-if="toast.state.value.type === 'loading'" class="loading-spinner"></view>
        <uv-icon 
          v-else-if="toast.state.value.type === 'success'" 
          name="checkmark-circle-fill" 
          size="24" 
          color="var(--color-success)"
        ></uv-icon>
        <uv-icon 
          v-else-if="toast.state.value.type === 'error'" 
          name="close-circle-fill" 
          size="24" 
          color="var(--color-error)"
        ></uv-icon>
        <uv-icon 
          v-else-if="toast.state.value.type === 'warning'" 
          name="error-circle-fill" 
          size="24" 
          color="var(--color-warning)"
        ></uv-icon>
        <uv-icon 
          v-else 
          name="info-circle-fill" 
          size="24" 
          color="var(--color-info)"
        ></uv-icon>
      </view>
      
      <!-- 消息 -->
      <text class="toast-message">{{ toast.state.value.message }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { getGlobalToast } from '@/composables/useToast'

const toast = getGlobalToast()
</script>

<style lang="scss" scoped>
.toast-container {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 10000;
  animation: toastIn 0.25s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

@keyframes toastIn {
  from {
    opacity: 0;
    transform: translate(-50%, -50%) scale(0.8);
  }
  to {
    opacity: 1;
    transform: translate(-50%, -50%) scale(1);
  }
}

.toast-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-lg) var(--spacing-xl);
  background: rgba(0, 0, 0, 0.85);
  border-radius: var(--radius-md);
  border: 1rpx solid rgba(255, 255, 255, 0.06);
}

.toast-icon {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-spinner {
  width: 48rpx;
  height: 48rpx;
  border: 4rpx solid rgba(255, 255, 255, 0.2);
  border-top-color: var(--color-primary);
  border-radius: var(--radius-full);
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.toast-message {
  font-size: var(--font-md);
  color: #fff;
  text-align: center;
  max-width: 400rpx;
  line-height: 1.4;
}

// 不同类型的样式
.toast--success {
  .toast-content {
    border: 1rpx solid rgba(16, 185, 129, 0.3);
  }
}

.toast--error {
  .toast-content {
    border: 1rpx solid rgba(239, 68, 68, 0.3);
  }
}

.toast--warning {
  .toast-content {
    border: 1rpx solid rgba(245, 158, 11, 0.3);
  }
}

.toast--info {
  .toast-content {
    border: 1rpx solid rgba(59, 130, 246, 0.3);
  }
}
</style>
