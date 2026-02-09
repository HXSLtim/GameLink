<template>
  <view v-if="loading" class="loading-wrap">
    <Skeleton :rows="5" />
  </view>

  <scroll-view
    v-else
    class="message-scroll"
    scroll-y
    :scroll-into-view="scrollToId"
    :scroll-with-animation="true"
    @scrolltoupper="$emit('load-more')"
  >
    <!-- 加载更多历史 -->
    <view v-if="hasMoreHistory" class="load-more-history" @tap="$emit('load-more')">
      <view class="load-more-inner">
        <uv-icon v-if="loadingHistory" name="reload" size="12" color="var(--color-text-secondary)" />
        <text>{{ loadingHistory ? '加载中...' : '查看更早的消息' }}</text>
      </view>
    </view>

    <ChatMessageBubble
      v-for="(message, index) in messages"
      :key="message.id"
      :id="`msg-${message.id}`"
      :message="message"
      :current-user-id="currentUserId"
      :self-name="currentUserName"
      :self-avatar="currentUserAvatar"
      :show-time="showTime(message, index)"
      :show-sender-name="showSenderName"
      :is-playing="playingVoiceId === message.id"
      @avatar-click="$emit('avatar-click', $event)"
      @image-preview="$emit('image-preview', $event)"
      @voice-play="$emit('voice-play', $event)"
      @order-click="$emit('order-click', $event)"
      @resend="$emit('resend', $event)"
    />

    <view id="msg-bottom" class="scroll-bottom"></view>
  </scroll-view>
</template>

<script setup lang="ts">
import Skeleton from '@/components/Skeleton/index.vue'
import ChatMessageBubble from '@/components/ChatMessageBubble/index.vue'
import type { ChatMessageData, ChatShowTimeFn } from '@/types/message'

interface Props {
  loading: boolean
  loadingHistory: boolean
  hasMoreHistory: boolean
  messages: ChatMessageData[]
  scrollToId: string
  currentUserId: number
  currentUserName?: string
  currentUserAvatar?: string
  showTime: ChatShowTimeFn
  showSenderName: boolean
  playingVoiceId?: string | null
}

defineProps<Props>()

defineEmits<{
  'load-more': []
  'avatar-click': [userId: number]
  'image-preview': [url: string]
  'voice-play': [message: ChatMessageData]
  'order-click': [orderId: number | undefined]
  resend: [message: ChatMessageData]
}>()
</script>

<style lang="scss" scoped>
.loading-wrap {
  flex: 1;
  padding: var(--spacing-md);
}

.message-scroll {
  flex: 1;
  overflow-y: auto;
  background: var(--color-bg);
}

.load-more-history {
  display: flex;
  justify-content: center;
  padding: var(--spacing-md);
}

.load-more-inner {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-xs);
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
  background: var(--color-bg-secondary);
  padding: 8rpx var(--spacing-md);
  border-radius: 100rpx;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: var(--color-bg-card);
    color: var(--color-primary);
  }
}

.scroll-bottom {
  height: var(--spacing-md);
}
</style>
