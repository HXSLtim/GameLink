<template>
  <view class="home-header">
    <view class="header-left">
      <text class="app-name">GameLink</text>
    </view>
    <view class="header-right">
      <!-- 主题切换 -->
      <view class="theme-toggle" @click="toggleTheme">
        <uv-icon :name="isDark ? 'eye-off' : 'eye'" size="22" :color="isDark ? '#FFD700' : '#FF8C00'"></uv-icon>
      </view>
      
      <!-- 登录按钮 / 用户头像 -->
      <GlButton v-if="!isLoggedIn" type="primary" size="small" round @click="$emit('login')">
        登录
      </GlButton>
      <view v-else class="user-mini" @click="$emit('profile')">
        <GlAvatar :src="avatar" :text="nickname" size="small" />
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { useTheme } from '@/composables/useTheme'
import GlButton from '@/components/gl/Button/index.vue'
import GlAvatar from '@/components/gl/Avatar/index.vue'

interface Props {
  isLoggedIn: boolean
  avatar?: string
  nickname?: string
}

defineProps<Props>()

defineEmits<{
  login: []
  profile: []
}>()

const { isDark, toggleTheme } = useTheme()
</script>

<style lang="scss" scoped>
.home-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx 32rpx;
  background: var(--color-bg-card);
  border-bottom: 1rpx solid var(--color-border);
}

.header-left {
  .app-name {
    font-size: 40rpx;
    font-weight: 800;
    background: linear-gradient(135deg, var(--color-primary) 0%, #00E676 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    letter-spacing: 2rpx;
  }
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.theme-toggle {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-secondary);
  border-radius: 50%;
  transition: all 0.3s;
  
  &:active {
    transform: scale(0.9);
  }
}

.user-mini {
  cursor: pointer;
}
</style>
