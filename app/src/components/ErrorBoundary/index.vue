<template>
  <view class="error-boundary">
    <!-- 网络错误 -->
    <template v-if="type === 'network'">
      <view class="error-icon"><uv-icon name="wifi-off" size="48" color="var(--color-text-placeholder)" /></view>
      <text class="error-title">网络连接失败</text>
      <text class="error-desc">请检查网络设置后重试</text>
    </template>

    <!-- 服务器错误 -->
    <template v-else-if="type === 'server'">
      <view class="error-icon"><uv-icon name="setting" size="48" color="var(--color-text-placeholder)" /></view>
      <text class="error-title">服务暂时不可用</text>
      <text class="error-desc">工程师正在紧急修复中</text>
    </template>

    <!-- 数据加载失败 -->
    <template v-else-if="type === 'data'">
      <view class="error-icon"><uv-icon name="list" size="48" color="var(--color-text-placeholder)" /></view>
      <text class="error-title">数据加载失败</text>
      <text class="error-desc">{{ message || '请稍后重试' }}</text>
    </template>

    <!-- 权限不足 -->
    <template v-else-if="type === 'permission'">
      <view class="error-icon"><uv-icon name="lock" size="48" color="var(--color-text-placeholder)" /></view>
      <text class="error-title">暂无权限访问</text>
      <text class="error-desc">请先登录或联系管理员</text>
    </template>

    <!-- 页面不存在 -->
    <template v-else-if="type === '404'">
      <view class="error-icon"><uv-icon name="search" size="48" color="var(--color-text-placeholder)" /></view>
      <text class="error-title">页面不存在</text>
      <text class="error-desc">请检查链接是否正确</text>
    </template>

    <!-- 通用错误 -->
    <template v-else>
      <view class="error-icon"><uv-icon name="error-circle" size="48" color="var(--color-text-placeholder)" /></view>
      <text class="error-title">{{ title || '出错了' }}</text>
      <text class="error-desc">{{ message || '请稍后重试' }}</text>
    </template>

    <!-- 操作按钮 -->
    <view class="error-actions">
      <view v-if="showRetry" class="action-btn primary" @tap="handleRetry">
        <text>重试</text>
      </view>
      <view v-if="showBack" class="action-btn secondary" @tap="handleBack">
        <text>返回</text>
      </view>
      <view v-if="showHome" class="action-btn secondary" @tap="handleHome">
        <text>回到首页</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
interface Props {
  type?: 'network' | 'server' | 'data' | 'permission' | '404' | 'unknown'
  title?: string
  message?: string
  showRetry?: boolean
  showBack?: boolean
  showHome?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  type: 'unknown',
  showRetry: true,
  showBack: true,
  showHome: false,
})

const emit = defineEmits<{
  retry: []
}>()

const handleRetry = () => {
  emit('retry')
}

const handleBack = () => {
  uni.navigateBack({
    fail: () => {
      uni.switchTab({ url: '/pages/index/index' })
    }
  })
}

const handleHome = () => {
  uni.switchTab({ url: '/pages/index/index' })
}
</script>

<style lang="scss" scoped>
.error-boundary {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-xl) var(--spacing-lg);
  min-height: 60vh;

  .error-icon {
    margin-bottom: var(--spacing-lg);
  }

  .error-title {
    font-size: var(--font-lg);
    font-weight: 600;
    color: var(--color-text);
    margin-bottom: var(--spacing-sm);
    text-align: center;
  }

  .error-desc {
    font-size: var(--font-md);
    color: var(--color-text-secondary);
    text-align: center;
    margin-bottom: var(--spacing-xl);
  }

  .error-actions {
    display: flex;
    gap: var(--spacing-md);

    .action-btn {
      padding: var(--spacing-sm) var(--spacing-xl);
      border-radius: var(--radius-sm);
      font-size: var(--font-md);
      cursor: pointer;
      @include press-effect;

      &.primary {
        background: var(--color-primary);
        color: #FFFFFF;
      }

      &.secondary {
        background: var(--color-bg-card);
        color: var(--color-text);
        border: 1rpx solid var(--color-border);
      }
    }
  }
}
</style>
