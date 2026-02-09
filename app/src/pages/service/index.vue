<template>
  <PageShell class="service-page" :scroll="false" padding="0" content-class="service-content">
    <template #header>
      <!-- 顶部导航 -->
      <NavBar title="在线客服" @back="goBack" />
    </template>

    <!-- 客服状态 -->
    <SupportAgentStatus :is-online="isOnline" />

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
        <GlInput
          v-model="inputContent"
          class="chat-input"
          placeholder="请输入您的问题..."
          :disabled="sending"
          size="small"
          variant="plain"
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

    <template #footer>
      <!-- PC 端侧边栏 -->
      <CustomTabBar :show-mobile-tab-bar="false" />
    </template>
  </PageShell>
</template>

<script setup lang="ts">
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import PageShell from '@/components/layout/PageShell/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
import GlInput from '@/components/gl/Input/index.vue'
// Business 组件
import SupportAgentStatus from '@/components/SupportAgentStatus/index.vue'
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
.service-content {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
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
  min-width: 0;
  height: 100%;
  
  :deep(.gl-input__field) {
    font-size: var(--font-md);
    color: var(--color-text);
  }
}
</style>
