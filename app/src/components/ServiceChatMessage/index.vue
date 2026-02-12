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
import type { ServiceMessage } from '@/types/CustomerService'

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
  gap: var(--spacing-xs);
  margin-bottom: var(--spacing-md);

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
  background: var(--color-bg-secondary);
  border: 1rpx solid var(--color-border);
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-base);
  flex-shrink: 0;
}

.msg-content {
  max-width: 70%;
}

.msg-bubble {
  padding: var(--spacing-sm) var(--spacing-md);
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  border-top-left-radius: var(--radius-sm);

  &.mine {
    background: var(--color-primary);
    border-radius: var(--radius-md);
    border-top-right-radius: var(--radius-sm);

    text {
      color: #FFFFFF;
    }
  }

  text {
    font-size: var(--font-md);
    color: var(--color-text);
    line-height: 1.6;
    white-space: pre-wrap;
  }
}

.msg-time {
  display: block;
  font-size: var(--font-xs);
  color: var(--color-text-placeholder);
  margin-top: var(--spacing-xs);
}
</style>
