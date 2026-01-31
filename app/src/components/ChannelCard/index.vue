<template>
  <view class="channel-card" @tap="$emit('click')">
    <!-- 头像 -->
    <view class="channel-avatar">
      <GlAvatar
        :src="channel.avatar"
        :text="channel.name"
        size="large"
        shape="square"
      />
    </view>
    
    <!-- 信息 -->
    <view class="channel-info">
      <view class="channel-header">
        <text class="channel-name">{{ channel.name }}</text>
        <GlTag v-if="channel.isActive" type="success" size="mini">活跃</GlTag>
      </view>
      <text class="channel-desc">{{ channel.description || '暂无描述' }}</text>
      <view class="channel-meta">
        <view class="meta-item">
          <uv-icon name="account-fill" size="12" color="var(--color-text-secondary)"></uv-icon>
          <text>{{ channel.memberCount }} 人</text>
        </view>
        <view v-if="channel.gameName" class="meta-item">
          <uv-icon name="grid-fill" size="12" color="var(--color-text-secondary)"></uv-icon>
          <text>{{ channel.gameName }}</text>
        </view>
      </view>
    </view>
    
    <!-- 操作按钮 -->
    <view class="channel-action">
      <GlButton
        v-if="channel.isJoined"
        size="small"
        type="default"
        plain
        round
        @click.stop="$emit('leave')"
      >
        已加入
      </GlButton>
      <GlButton
        v-else
        size="small"
        type="primary"
        round
        @click.stop="$emit('join')"
      >
        加入
      </GlButton>
    </view>
  </view>
</template>

<script setup lang="ts">
import GlAvatar from '@/components/gl/Avatar/index.vue'
import GlTag from '@/components/gl/Tag/index.vue'
import GlButton from '@/components/gl/Button/index.vue'

export interface ChannelData {
  id: number
  name: string
  description?: string
  avatar?: string
  memberCount: number
  maxMembers?: number
  isActive: boolean
  isJoined: boolean
  gameId?: number
  gameName?: string
}

interface Props {
  channel: ChannelData
}

defineProps<Props>()

defineEmits<{
  click: []
  join: []
  leave: []
}>()
</script>

<style lang="scss" scoped>
.channel-card {
  display: flex;
  gap: 24rpx;
  padding: 28rpx;
  background: var(--color-bg-card);
  border-radius: 20rpx;
  margin-bottom: 16rpx;
  border: 2rpx solid var(--color-border);
  transition: all 0.2s;
  
  &:active {
    transform: scale(0.99);
    border-color: var(--color-primary);
  }
}

.channel-avatar {
  flex-shrink: 0;
}

.channel-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.channel-header {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.channel-name {
  font-size: 32rpx;
  font-weight: 600;
  color: var(--color-text);
}

.channel-desc {
  font-size: 26rpx;
  color: var(--color-text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.channel-meta {
  display: flex;
  gap: 24rpx;
  margin-top: 8rpx;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6rpx;
  font-size: 22rpx;
  color: var(--color-text-secondary);
}

.channel-action {
  flex-shrink: 0;
  display: flex;
  align-items: center;
}
</style>
