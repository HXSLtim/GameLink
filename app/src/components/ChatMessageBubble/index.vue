<template>
  <view class="message-wrap">
    <!-- 时间分割线 -->
    <view v-if="showTime" class="time-divider">
      <view class="time-divider-line" />
      <text class="time-divider-text">{{ formattedTime }}</text>
      <view class="time-divider-line" />
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
          <text class="bubble-text">{{ message.content }}</text>
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
          <uv-icon name="volume" size="18" color="currentColor"></uv-icon>
          <view class="voice-wave">
            <view v-for="i in 4" :key="i" class="wave-bar"></view>
          </view>
          <text class="voice-duration">{{ message.duration }}"</text>
        </view>
        
        <!-- 订单消息 -->
        <view v-else-if="message.type === 'order'" class="bubble order-bubble" @tap="$emit('order-click', message.orderId)">
          <view class="order-header">
            <view class="order-icon-wrap">
              <uv-icon name="list" size="14" color="var(--color-primary)" />
            </view>
            <text class="order-title">订单消息</text>
          </view>
          <text class="order-content">{{ message.content }}</text>
          <view class="order-footer">
            <text class="order-action">查看详情</text>
            <uv-icon name="arrow-right" size="12" color="var(--color-primary)" />
          </view>
        </view>
        
        <!-- 消息状态 -->
        <view v-if="isSelf" class="message-status">
          <text v-if="message.status === 'sending'" class="status-sending">发送中...</text>
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
import { formatMonthDayTime, formatTime } from '@/utils/format'
import type { ChatMessageData } from '@/types/message'

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
  return isToday ? formatTime(date) : formatMonthDayTime(date)
})
</script>

<style lang="scss" scoped>
.message-wrap {
  padding: var(--spacing-sm) var(--spacing-md);
}

// ── 时间分割线 ──
.time-divider {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-lg) var(--spacing-xl);
}

.time-divider-line {
  flex: 1;
  height: 1rpx;
  background: var(--color-border);
}

.time-divider-text {
  font-size: var(--font-xs);
  color: var(--color-text-placeholder);
  white-space: nowrap;
}

// ── 系统消息 ──
.system-message {
  text-align: center;
  padding: var(--spacing-md) 0;

  text {
    font-size: var(--font-xs);
    color: var(--color-text-placeholder);
    background: var(--color-bg-secondary);
    padding: 8rpx var(--spacing-md);
    border-radius: 100rpx;
    line-height: 1.6;
  }
}

// ── 消息布局 ──
.message-item {
  display: flex;
  align-items: flex-start;
  gap: var(--spacing-md);

  &.is-self {
    flex-direction: row-reverse;

    .bubble {
      background: var(--color-primary);
      color: #fff;
      border-color: transparent;
      // 自己消息：右上角不圆
      border-radius: var(--radius-lg) 8rpx var(--radius-lg) var(--radius-lg);
    }

    .text-bubble {
      box-shadow: 0 2rpx 8rpx rgba(var(--color-primary-rgb), 0.2);
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
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
  margin-bottom: var(--spacing-xs);
}

// ── 气泡通用 ──
.bubble {
  padding: var(--spacing-md) var(--spacing-lg);
  background: var(--color-bg-card);
  // 他人消息：左上角不圆
  border-radius: 8rpx var(--radius-lg) var(--radius-lg) var(--radius-lg);
  border: 1rpx solid var(--color-border);
  word-break: break-all;
  transition: transform 0.15s ease;

  &:active {
    transform: scale(0.98);
  }
}

// ── 文字气泡 ──
.bubble-text {
  font-size: var(--font-sm);
  line-height: 1.6;
}

// ── 图片气泡 ──
.image-bubble {
  padding: 4rpx;
  background: transparent;
  border: none;
  overflow: hidden;
  border-radius: var(--radius-md);
}

.message-image {
  max-width: 400rpx;
  border-radius: var(--radius-sm);
  display: block;
}

// ── 语音气泡 ──
.voice-bubble {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  min-width: 160rpx;
  cursor: pointer;
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
  opacity: 0.7;
  
  &:nth-child(2) { height: 24rpx; }
  &:nth-child(3) { height: 32rpx; }
  &:nth-child(4) { height: 20rpx; }
}

.voice-bubble.playing .wave-bar {
  animation: wave 0.5s ease-in-out infinite alternate;
  opacity: 1;
  
  &:nth-child(2) { animation-delay: 0.1s; }
  &:nth-child(3) { animation-delay: 0.2s; }
  &:nth-child(4) { animation-delay: 0.3s; }
}

@keyframes wave {
  from { transform: scaleY(0.5); }
  to { transform: scaleY(1); }
}

.voice-duration {
  font-size: var(--font-xs);
  opacity: 0.8;
}

// ── 订单气泡 ──
.order-bubble {
  min-width: 300rpx;
  padding: var(--spacing-md);
  cursor: pointer;
}

.order-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-sm);
}

.order-icon-wrap {
  width: 40rpx;
  height: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(var(--color-primary-rgb), 0.1);
  border-radius: var(--radius-sm);
}

// 自己发的订单卡片内图标
.is-self .order-icon-wrap {
  background: rgba(255, 255, 255, 0.2);
}

.is-self .order-icon-wrap :deep(.uv-icon__icon) {
  color: #fff !important;
}

.is-self .order-footer :deep(.uv-icon__icon) {
  color: rgba(255, 255, 255, 0.8) !important;
}

.order-title {
  font-size: var(--font-sm);
  font-weight: 600;
}

.order-content {
  font-size: var(--font-sm);
  line-height: 1.5;
  opacity: 0.85;
  margin-bottom: var(--spacing-sm);
}

.order-footer {
  display: flex;
  align-items: center;
  gap: 4rpx;
  padding-top: var(--spacing-sm);
  border-top: 1rpx solid var(--color-border);
}

.is-self .order-footer {
  border-color: rgba(255, 255, 255, 0.15);
}

.order-action {
  font-size: var(--font-xs);
  color: var(--color-primary);
  font-weight: 500;
}

.is-self .order-action {
  color: rgba(255, 255, 255, 0.8);
}

// ── 消息状态 ──
.message-status {
  margin-top: 4rpx;
  padding: 0 var(--spacing-xs);
}

.status-sending {
  font-size: var(--font-xs);
  color: var(--color-text-placeholder);
}

.status-read {
  font-size: var(--font-xs);
  color: var(--color-text-placeholder);
}

.status-failed {
  font-size: var(--font-xs);
  color: var(--color-error);
  cursor: pointer;
  
  &:active {
    opacity: 0.7;
  }
}
</style>
