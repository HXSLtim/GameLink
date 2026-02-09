<template>
  <BasePageLayout
    class="settings-page"
    padding="20rpx 0 0"
    title="设置"
    :show-back="true"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
  >

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

  </BasePageLayout>
</template>

<script setup lang="ts">
// Pattern 组件
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
// Business 组件
import SettingsSection from '@/components/SettingsSection/index.vue'
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
.logout-section {
  padding: 48rpx 24rpx;
}

.bottom-placeholder {
  height: 100rpx;
}
</style>
