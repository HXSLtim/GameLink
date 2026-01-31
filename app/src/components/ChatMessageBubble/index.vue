<template>
  <view class="message-wrap">
    <!-- 时间分割线 -->
    <view v-if="showTime" class="time-divider">
      <text>{{ formattedTime }}</text>
    </view>
    
    <!-- 系统消息 -->
    <view v-if="message.type === 'system'" class="system-message">
      <text>{{ message.content }}</text>
    </view>
    
    <!-- 普通消息 -->
    <view v-else class="message-item" :class="{ 'is-self': isSelf }">
      <!-- 头像（非自己） -->
      <GlAvatar 
        v-if="!isSelf" 
        :src="message.senderAvatar" 
        :text="message.senderName" 
        :size="72"
        @click="$emit('avatar-click', message.senderId)"
      />
      
      <!-- 消息内容 -->
      <view class="message-content">
        <!-- 发送者名称（群聊显示） -->
        <text v-if="showSenderName && !isSelf" class="sender-name">
          {{ message.senderName }}
        </text>
        
        <!-- 文本消息 -->
        <view v-if="message.type === 'text'" class="bubble text-bubble">
          <text>{{ message.content }}</text>
        </view>
        
        <!-- 图片消息 -->
        <view v-else-if="message.type === 'image'" class="bubble image-bubble">
          <image 
            :src="message.content"
            mode="widthFix"
            class="message-image"
            @tap="$emit('image-preview', message.content)"
          />
        </view>
        
        <!-- 语音消息 -->
        <view 
          v-else-if="message.type === 'voice'" 
          class="bubble voice-bubble"
          :class="{ playing: isPlaying }"
          @tap="$emit('voice-play', message)"
        >
          <uv-icon name="volume" size="20" color="currentColor"></uv-icon>
          <view class="voice-wave">
            <view v-for="i in 4" :key="i" class="wave-bar"></view>
          </view>
          <text class="voice-duration">{{ message.duration }}"</text>
        </view>
        
        <!-- 订单消息 -->
        <view v-else-if="message.type === 'order'" class="bubble order-bubble" @tap="$emit('order-click', message.orderId)">
          <view class="order-header">
            <text class="order-icon">📋</text>
            <text class="order-title">订单消息</text>
          </view>
          <text class="order-content">{{ message.content }}</text>
          <text class="order-action">点击查看详情 ›</text>
        </view>
        
        <!-- 消息状态 -->
        <view v-if="isSelf" class="message-status">
          <text v-if="message.status === 'sending'" class="status-sending">发送中</text>
          <text v-else-if="message.status === 'failed'" class="status-failed" @tap="$emit('resend', message)">
            发送失败，点击重试
          </text>
          <text v-else-if="message.status === 'read'" class="status-read">已读</text>
        </view>
      </view>
      
      <!-- 头像（自己） -->
      <GlAvatar v-if="isSelf" :src="selfAvatar" :text="selfName" :size="72" />
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlAvatar from '@/components/gl/Avatar/index.vue'

export interface ChatMessageData {
  id: string
  type: 'text' | 'image' | 'voice' | 'order' | 'system'
  content: string
  senderId: number
  senderName?: string
  senderAvatar?: string
  createdAt: string
  status?: 'sending' | 'sent' | 'read' | 'failed'
  duration?: number
  orderId?: number
}

interface Props {
  message: ChatMessageData
  currentUserId: number
  selfName?: string
  selfAvatar?: string
  showTime?: boolean
  showSenderName?: boolean
  isPlaying?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showTime: false,
  showSenderName: false,
  isPlaying: false,
})

defineEmits<{
  'avatar-click': [userId: number]
  'image-preview': [url: string]
  'voice-play': [message: ChatMessageData]
  'order-click': [orderId: number | undefined]
  'resend': [message: ChatMessageData]
}>()

const isSelf = computed(() => props.message.senderId === props.currentUserId)

const formattedTime = computed(() => {
  const date = new Date(props.message.createdAt)
  const now = new Date()
  const isToday = date.toDateString() === now.toDateString()
  
  const time = `${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
  
  if (isToday) return time
  return `${date.getMonth() + 1}/${date.getDate()} ${time}`
})
</script>

<style lang="scss" scoped>
.message-wrap {
  padding: 8rpx 24rpx;
}

.time-divider {
  text-align: center;
  padding: 16rpx 0;
  
  text {
    font-size: 22rpx;
    color: var(--color-text-placeholder);
    background: var(--color-bg);
    padding: 8rpx 24rpx;
    border-radius: 20rpx;
  }
}

.system-message {
  text-align: center;
  padding: 16rpx 0;
  
  text {
    font-size: 24rpx;
    color: var(--color-text-secondary);
  }
}

.message-item {
  display: flex;
  align-items: flex-start;
  gap: 16rpx;
  
  &.is-self {
    flex-direction: row-reverse;
    
    .bubble {
      background: var(--color-primary);
      color: #fff;
      border-radius: 24rpx 4rpx 24rpx 24rpx;
    }
    
    .message-status {
      text-align: right;
    }
  }
}

.message-content {
  max-width: 70%;
  display: flex;
  flex-direction: column;
}

.sender-name {
  font-size: 22rpx;
  color: var(--color-text-secondary);
  margin-bottom: 8rpx;
}

.bubble {
  padding: 20rpx 24rpx;
  background: var(--color-bg-card);
  border-radius: 4rpx 24rpx 24rpx 24rpx;
  border: 1rpx solid var(--color-border);
  word-break: break-all;
}

.text-bubble text {
  font-size: 30rpx;
  line-height: 1.5;
}

.image-bubble {
  padding: 8rpx;
  background: transparent;
  border: none;
}

.message-image {
  max-width: 400rpx;
  border-radius: 16rpx;
}

.voice-bubble {
  display: flex;
  align-items: center;
  gap: 12rpx;
  min-width: 150rpx;
}

.voice-wave {
  display: flex;
  align-items: center;
  gap: 4rpx;
  height: 32rpx;
}

.wave-bar {
  width: 6rpx;
  height: 16rpx;
  background: currentColor;
  border-radius: 3rpx;
  
  &:nth-child(2) { height: 24rpx; }
  &:nth-child(3) { height: 32rpx; }
  &:nth-child(4) { height: 20rpx; }
}

.voice-bubble.playing .wave-bar {
  animation: wave 0.5s ease-in-out infinite alternate;
  
  &:nth-child(2) { animation-delay: 0.1s; }
  &:nth-child(3) { animation-delay: 0.2s; }
  &:nth-child(4) { animation-delay: 0.3s; }
}

@keyframes wave {
  from { transform: scaleY(0.5); }
  to { transform: scaleY(1); }
}

.voice-duration {
  font-size: 24rpx;
}

.order-bubble {
  min-width: 350rpx;
}

.order-header {
  display: flex;
  align-items: center;
  gap: 8rpx;
  margin-bottom: 12rpx;
}

.order-icon {
  font-size: 28rpx;
}

.order-title {
  font-size: 26rpx;
  font-weight: 600;
}

.order-content {
  font-size: 26rpx;
  color: inherit;
  opacity: 0.9;
  margin-bottom: 12rpx;
}

.order-action {
  font-size: 24rpx;
  opacity: 0.7;
}

.message-status {
  margin-top: 8rpx;
}

.status-sending,
.status-read {
  font-size: 20rpx;
  color: var(--color-text-placeholder);
}

.status-failed {
  font-size: 20rpx;
  color: var(--color-error);
}
</style>
