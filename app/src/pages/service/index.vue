<template>
  <view class="service-page page-container">
    <!-- 顶部导航 -->
    <NavBar title="在线客服" @back="goBack" />

    <!-- 客服状态 -->
    <ServiceStatusCard :is-online="isOnline" />

    <!-- 快捷问题 -->
    <QuickQuestionBar :questions="quickQuestions" @select="selectQuestion" />

    <!-- 聊天区域 -->
    <scroll-view 
      class="chat-scroll" 
      scroll-y 
      :scroll-into-view="scrollToView"
      scroll-with-animation
    >
      <view class="chat-list">
        <ServiceChatMessage
          v-for="msg in messages"
          :key="msg.id"
          :message="msg"
        />
      </view>
      <view id="chat-bottom"></view>
    </scroll-view>

    <!-- 输入区域 -->
    <view class="input-bar">
      <view class="input-wrap">
        <input
          v-model="inputContent"
          class="chat-input"
          placeholder="请输入您的问题..."
          :disabled="sending"
          @confirm="sendMessage"
        />
      </view>
      <GlButton 
        type="primary" 
        size="small" 
        round 
        :disabled="!inputContent.trim() || sending"
        @click="sendMessage"
      >
        发送
      </GlButton>
    </view>

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
// Business 组件
import ServiceStatusCard from '@/components/ServiceStatusCard/index.vue'
import QuickQuestionBar from '@/components/QuickQuestionBar/index.vue'
import ServiceChatMessage from '@/components/ServiceChatMessage/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// Composables
import { useCustomerService } from '@/composables/useCustomerService'

const {
  isOnline,
  sending,
  inputContent,
  scrollToView,
  messages,
  quickQuestions,
  selectQuestion,
  sendMessage,
  goBack,
} = useCustomerService()
</script>

<style lang="scss" scoped>
.service-page {
  height: 100vh;
  height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  box-sizing: border-box;
  overflow: hidden;
}

.chat-scroll {
  flex: 1;
  overflow-y: auto;
}

.chat-list {
  padding: 24rpx;
}

.input-bar {
  display: flex;
  align-items: center;
  gap: 16rpx;
  padding: 16rpx 24rpx;
  padding-bottom: calc(16rpx + env(safe-area-inset-bottom));
  background: var(--color-bg-card);
  border-top: 1rpx solid var(--color-border);
}

.input-wrap {
  flex: 1;
  height: 72rpx;
  padding: 0 24rpx;
  background: var(--color-bg-secondary);
  border-radius: 36rpx;
  display: flex;
  align-items: center;
}

.chat-input {
  flex: 1;
  font-size: 28rpx;
  color: var(--color-text);
}
</style>
