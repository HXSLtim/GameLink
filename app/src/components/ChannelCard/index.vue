<template>
  <view class="channel-card" @tap="$emit('click')">
    <!-- 封面 -->
    <view class="channel-cover">
      <GlAvatar
        :src="channel.avatar"
        :text="channel.name"
        size="xlarge"
        shape="square"
      />
      <GlTag v-if="channel.isActive" type="success" size="mini" class="channel-badge">活跃</GlTag>
    </view>
    
    <!-- 信息 -->
    <view class="channel-body">
      <view class="channel-title-row">
        <text class="channel-name">{{ channel.name }}</text>
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
  </view>
</template>

<script setup lang="ts">
import GlAvatar from '@/components/gl/Avatar/index.vue'
import GlTag from '@/components/gl/Tag/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
import type { ChannelData } from '@/types/community'

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
  flex-direction: column;
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1rpx solid var(--color-border);
  overflow: hidden;
  transition: all 0.2s;
  cursor: pointer;
  
  &:active {
    background: var(--color-bg-secondary);
    border-color: var(--color-border);
  }
}

.channel-cover {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-sm);
  background: var(--color-bg-secondary);
  border-bottom: 1rpx solid var(--color-border);
}

.channel-badge {
  position: absolute;
  top: var(--spacing-sm);
  right: var(--spacing-sm);
}

.channel-body {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
  padding: var(--spacing-sm);
}

.channel-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-sm);
}

.channel-name {
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
  flex: 1;
  min-width: 0;
  @include text-ellipsis;
}

.channel-desc {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  @include text-ellipsis;
}

.channel-meta {
  display: flex;
  gap: var(--spacing-md);
  margin-top: var(--spacing-xs);
}

.meta-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}

.channel-action {
  flex-shrink: 0;
  display: flex;
  align-items: center;
}
</style>
