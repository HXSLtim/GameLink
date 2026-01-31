<template>
  <view class="message-item" @tap="$emit('click')">
    <!-- 头像 -->
    <view class="message-avatar">
      <GlAvatar
        :src="message.avatar"
        :text="message.name"
        size="large"
        :badge="unreadBadge"
      />
    </view>
    
    <!-- 内容 -->
    <view class="message-content">
      <view class="message-header">
        <text class="message-name">{{ message.name }}</text>
        <text class="message-time">{{ formattedTime }}</text>
      </view>
      <text class="message-text">{{ message.lastMessage }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlAvatar from '@/components/gl/Avatar/index.vue'

export interface MessageData {
  id: number
  conversationId: number
  avatar: string
  name: string
  lastMessage: string
  lastTime: number
  unread: number
  type: 'chat' | 'system' | 'order'
}

interface Props {
  message: MessageData
}

const props = defineProps<Props>()

defineEmits<{
  click: []
}>()

// 未读徽章
const unreadBadge = computed(() => {
  if (props.message.unread <= 0) return undefined
  if (props.message.unread > 99) return '99+'
  return props.message.unread
})

// 格式化时间
const formattedTime = computed(() => {
  const timestamp = props.message.lastTime
  const now = Date.now()
  const diff = now - timestamp
  
  if (diff < 60 * 1000) return '刚刚'
  if (diff < 60 * 60 * 1000) return `${Math.floor(diff / 60 / 1000)}分钟前`
  if (diff < 24 * 60 * 60 * 1000) return `${Math.floor(diff / 60 / 60 / 1000)}小时前`
  if (diff < 7 * 24 * 60 * 60 * 1000) return `${Math.floor(diff / 24 / 60 / 60 / 1000)}天前`
  
  const date = new Date(timestamp)
  return `${date.getMonth() + 1}/${date.getDate()}`
})
</script>

<style lang="scss" scoped>
.message-item {
  display: flex;
  gap: 28rpx;
  padding: 32rpx;
  background: var(--color-bg-card);
  border-bottom: 1rpx solid var(--color-border);
  transition: all 0.2s;
  
  &:active {
    background: var(--color-bg-secondary);
    transform: scale(0.99);
  }
}

.message-avatar {
  flex-shrink: 0;
}

.message-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12rpx;
}

.message-name {
  font-size: 32rpx;
  font-weight: 600;
  color: var(--color-text);
  letter-spacing: 1rpx;
}

.message-time {
  font-size: 24rpx;
  color: var(--color-text-placeholder);
  font-weight: 500;
}

.message-text {
  font-size: 28rpx;
  color: var(--color-text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.5;
  margin-top: 8rpx;
}
</style>
