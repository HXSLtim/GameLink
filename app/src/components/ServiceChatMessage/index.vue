<template>
  <view class="chat-item" :class="{ mine: message.isMe }">
    <view v-if="!message.isMe" class="msg-avatar">
      <text>👩‍💼</text>
    </view>
    <view class="msg-content">
      <view class="msg-bubble" :class="{ mine: message.isMe }">
        <text>{{ message.content }}</text>
      </view>
      <text class="msg-time">{{ formattedTime }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

export interface ServiceMessage {
  id: number
  content: string
  isMe: boolean
  createdAt: string
}

interface Props {
  message: ServiceMessage
}

const props = defineProps<Props>()

const formattedTime = computed(() => {
  const date = new Date(props.message.createdAt)
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${hours}:${minutes}`
})
</script>

<style lang="scss" scoped>
.chat-item {
  display: flex;
  gap: 12rpx;
  margin-bottom: 24rpx;

  &.mine {
    flex-direction: row-reverse;
    
    .msg-content .msg-time {
      text-align: right;
    }
  }
}

.msg-avatar {
  width: 64rpx;
  height: 64rpx;
  background: linear-gradient(135deg, var(--color-primary), #00B85C);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32rpx;
  flex-shrink: 0;
}

.msg-content {
  max-width: 70%;
}

.msg-bubble {
  padding: 16rpx 20rpx;
  background: var(--color-bg-card);
  border-radius: 16rpx;
  border-top-left-radius: 4rpx;

  &.mine {
    background: var(--color-primary);
    border-radius: 16rpx;
    border-top-right-radius: 4rpx;

    text {
      color: #FFFFFF;
    }
  }

  text {
    font-size: 28rpx;
    color: var(--color-text);
    line-height: 1.6;
    white-space: pre-wrap;
  }
}

.msg-time {
  display: block;
  font-size: 22rpx;
  color: var(--color-text-placeholder);
  margin-top: 8rpx;
}
</style>
