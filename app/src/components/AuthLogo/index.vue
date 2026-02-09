<template>
  <view class="logo-section">
    <view v-if="variant === 'icon'" class="logo-icon">
      <text class="logo-icon-text">{{ brandInitial }}</text>
    </view>
    <image v-else class="logo" :src="logo" mode="aspectFit" />
    <text class="app-name">{{ appName }}</text>
    <text class="app-slogan">{{ slogan }}</text>
  </view>
</template>

<script setup lang="ts">
interface Props {
  logo?: string
  appName?: string
  slogan?: string
  variant?: 'image' | 'icon'
}

const props = withDefaults(defineProps<Props>(), {
  logo: '/static/logo.png',
  appName: 'GameLink',
  slogan: '专业游戏陪玩平台',
  variant: 'image',
})

const brandInitial = (props.appName || 'G').trim().slice(0, 1).toUpperCase() || 'G'
</script>

<style lang="scss" scoped>
.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: var(--spacing-lg) 0 var(--spacing-md);
}

.logo {
  width: 120rpx;
  height: 120rpx;
  margin-bottom: var(--spacing-md);
}

.logo-icon {
  width: 120rpx;
  height: 120rpx;
  margin-bottom: var(--spacing-md);
  border-radius: var(--radius-md);
  background: var(--color-icon-accent);
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1rpx solid var(--color-border);
}

.logo-icon-text {
  font-size: var(--font-xl);
  font-weight: 700;
  color: #fff;
  letter-spacing: 1rpx;
}

.app-name {
  font-size: var(--font-xl);
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: var(--spacing-xs);
  letter-spacing: 0.5rpx;
}

.app-slogan {
  font-size: var(--font-md);
  color: var(--color-text-secondary);
}
</style>
