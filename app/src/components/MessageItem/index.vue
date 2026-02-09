<template>
  <view class="message-item" :class="{ 'has-unread': hasUnread }" @tap="$emit('click')">
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
      <view class="message-preview">
        <!-- 消息类型前缀图标 -->
        <text v-if="typePrefix" class="message-type-prefix">{{ typePrefix }}</text>
        <text class="message-text">{{ message.lastMessage }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlAvatar from '@/components/gl/Avatar/index.vue'
import { formatRelativeTimeShort } from '@/utils/format'
import type { MessageData } from '@/types/message'

interface Props {
  message: MessageData
}

const props = defineProps<Props>()

defineEmits<{
  click: []
}>()

// 是否有未读
const hasUnread = computed(() => (props.message.unread ?? 0) > 0)

// 未读徽章
const unreadBadge = computed(() => {
  if (props.message.unread <= 0) return undefined
  if (props.message.unread > 99) return '99+'
  return props.message.unread
})

// 消息类型前缀
const typePrefix = computed(() => {
  const type = props.message.lastMessageType
  if (!type || type === 'text') return ''
  const map: Record<string, string> = {
    image: '[图片]',
    voice: '[语音]',
    order: '[订单]',
    system: '[系统]',
  }
  return map[type] || ''
})

// 格式化时间
const formattedTime = computed(() => formatRelativeTimeShort(props.message.lastTime))
</script>

<style lang="scss" scoped>
.message-item {
  display: flex;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  background: var(--color-bg-card);
  border: 1rpx solid var(--color-border);
  border-radius: var(--radius-md);
  transition: all 0.2s ease;
  cursor: pointer;

  &:hover {
    background: var(--color-bg-secondary);
    border-color: var(--color-primary);
    box-shadow: var(--shadow-sm);
  }

  &:active {
    transform: scale(0.99);
  }

  // 有未读消息时，左侧加强调线
  &.has-unread {
    border-left: 4rpx solid var(--color-primary);

    .message-name {
      font-weight: 700;
    }

    .message-text {
      color: var(--color-text);
      font-weight: 500;
    }
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
  gap: 6rpx;
}

.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.message-name {
  font-size: var(--font-base);
  font-weight: 600;
  color: var(--color-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.message-time {
  font-size: var(--font-xs);
  color: var(--color-text-placeholder);
  flex-shrink: 0;
  margin-left: var(--spacing-sm);
}

.message-preview {
  display: flex;
  align-items: center;
  gap: 4rpx;
  overflow: hidden;
}

.message-type-prefix {
  font-size: var(--font-sm);
  color: var(--color-primary);
  flex-shrink: 0;
  font-weight: 500;
}

.message-text {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.5;
}
</style>
