<template>
  <view class="settings-page page-container">
    <!-- 顶部导航 -->
    <NavBar title="设置" @back="goBack" />

    <scroll-view class="content-scroll" scroll-y>
      <!-- 账号安全 -->
      <SettingsSection
        title="账号安全"
        :items="securityItems"
        @click="handleClick"
      />

      <!-- 通知设置 -->
      <SettingsSection
        title="通知设置"
        :items="notificationItems"
        @click="handleClick"
        @switch="handleSwitch"
      />

      <!-- 隐私设置 -->
      <SettingsSection
        title="隐私设置"
        :items="privacyItems"
        @click="handleClick"
        @switch="handleSwitch"
      />

      <!-- 通用设置 -->
      <SettingsSection
        title="通用设置"
        :items="generalItems"
        @click="handleClick"
        @switch="handleSwitch"
      />

      <!-- 关于 -->
      <SettingsSection
        title="关于"
        :items="aboutItems"
        @click="handleClick"
      />

      <!-- 退出登录 -->
      <view class="logout-section">
        <GlButton type="error" plain block round size="large" @click="handleLogout">
          退出登录
        </GlButton>
      </view>

      <view class="bottom-placeholder"></view>
    </scroll-view>

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
// Business 组件
import SettingsSection from '@/components/SettingsSection/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// Composables
import { useSettings } from '@/composables/useSettings'

const {
  securityItems,
  notificationItems,
  privacyItems,
  generalItems,
  aboutItems,
  handleSwitch,
  handleClick,
  handleLogout,
  goBack,
} = useSettings()
</script>

<style lang="scss" scoped>
.settings-page {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  box-sizing: border-box;
  
  @include desktop {
    height: 100vh;
    min-height: auto;
    overflow: hidden;
  }
}

.content-scroll {
  flex: 1;
  overflow-y: auto;
  padding-top: 20rpx;
}

.logout-section {
  padding: 48rpx 24rpx;
}

.bottom-placeholder {
  height: 100rpx;
}
</style>
