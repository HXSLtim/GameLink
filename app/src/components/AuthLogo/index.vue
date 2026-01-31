<template>
  <view class="logo-section">
    <image class="logo" :src="logo" mode="aspectFit" />
    <text class="app-name">{{ appName }}</text>
    <text class="app-slogan">{{ slogan }}</text>
  </view>
</template>

<script setup lang="ts">
interface Props {
  logo?: string
  appName?: string
  slogan?: string
}

withDefaults(defineProps<Props>(), {
  logo: '/static/logo.png',
  appName: 'GameLink',
  slogan: '专业游戏陪玩平台',
})
</script>

<style lang="scss" scoped>
.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 80rpx 0 60rpx;
}

.logo {
  width: 160rpx;
  height: 160rpx;
  margin-bottom: 24rpx;
}

.app-name {
  font-size: 48rpx;
  font-weight: 800;
  color: var(--color-text);
  margin-bottom: 12rpx;
}

.app-slogan {
  font-size: 28rpx;
  color: var(--color-text-secondary);
}
</style>
