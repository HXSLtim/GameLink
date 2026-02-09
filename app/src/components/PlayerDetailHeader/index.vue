<template>
  <view class="header-card">
    <!-- 封面 -->
    <view class="player-cover">
      <image v-if="player.coverImage" :src="player.coverImage" mode="aspectFill" class="cover-image" />
      <view v-else class="cover-placeholder"></view>
    </view>
    
    <!-- 基本信息 -->
    <view class="player-basic">
      <view class="avatar-wrap">
        <GlAvatar 
          :src="player.avatar" 
          :text="player.nickname" 
          :size="100" 
          :status="player.isOnline ? 'online' : undefined"
          bordered
        />
      </view>
      
      <view class="basic-info">
        <view class="name-row">
          <text class="nickname">{{ player.nickname }}</text>
          <GlTag v-if="player.isVerified" type="success" size="mini">已认证</GlTag>
        <view v-if="player.gender === 'male' || player.gender === 'female'" class="gender-badge" :class="player.gender">
          <text>{{ player.gender === 'male' ? '♂' : '♀' }}</text>
        </view>
        </view>
        
        <view class="status-row">
          <GlTag :type="player.isOnline ? 'success' : 'default'" size="mini" plain>
            {{ player.isOnline ? '在线' : '离线' }}
          </GlTag>
          <text v-if="player.signature" class="signature">{{ player.signature }}</text>
        </view>
      </view>
    </view>
    
    <!-- 统计数据 -->
    <HeaderStatsRow
      class="stats-row"
      :items="stats"
      size="md"
      item-padding="0"
    />
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlAvatar from '@/components/gl/Avatar/index.vue'
import GlTag from '@/components/gl/Tag/index.vue'
import HeaderStatsRow from '@/components/HeaderStatsRow/index.vue'
import type { PlayerHeaderData } from '@/types/player'
import type { HeaderStatItem } from '@/types/ui'

interface Props {
  player: PlayerHeaderData
}

const props = defineProps<Props>()

const formatJoinDate = computed(() => {
  if (!props.player.createdAt) return '-'
  const date = new Date(props.player.createdAt)
  const now = new Date()
  const months = (now.getFullYear() - date.getFullYear()) * 12 + (now.getMonth() - date.getMonth())
  if (months < 1) return '刚入驻'
  if (months < 12) return `${months}个月`
  return `${Math.floor(months / 12)}年`
})

const stats = computed<HeaderStatItem[]>(() => [
  { label: '评分', value: props.player.rating?.toFixed(1) || '5.0' },
  { label: '接单数', value: props.player.orderCount || 0 },
  { label: '收藏数', value: props.player.favoriteCount || 0 },
  { label: '入驻时间', value: formatJoinDate.value },
])
</script>

<style lang="scss" scoped>
.header-card {
  position: relative;
  background: var(--color-bg-card);
  border-radius: 0 0 var(--radius-lg) var(--radius-lg);
  overflow: hidden;
  margin-bottom: var(--spacing-md);
  border-bottom: 1rpx solid var(--color-border);
}

.player-cover {
  height: 280rpx;
  background: var(--color-bg-secondary);
  position: relative;

  .cover-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  // 底部渐变遮罩，让头像区文字更易读
  &::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 100rpx;
    background: linear-gradient(to top, var(--color-bg-card), transparent);
    pointer-events: none;
  }
}

.player-basic {
  display: flex;
  gap: var(--spacing-md);
  padding: 0 var(--spacing-lg);
  margin-top: -56rpx;
  position: relative;
  z-index: 1;
}

.avatar-wrap {
  flex-shrink: 0;
}

.basic-info {
  flex: 1;
  padding-top: 60rpx;
  min-width: 0;
}

.name-row {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  flex-wrap: wrap;
  margin-bottom: var(--spacing-sm);
}

.nickname {
  font-size: var(--font-lg);
  font-weight: 700;
  color: var(--color-text);
}

.gender-badge {
  width: 36rpx;
  height: 36rpx;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-xs);
  
  &.male {
    background: var(--color-info-tint);
    color: var(--color-info);
  }
  
  &.female {
    background: var(--color-error-tint);
    color: var(--color-error);
  }
}

.status-row {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.signature {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  @include text-ellipsis;
}

.stats-row {
  padding: var(--spacing-md) var(--spacing-lg);
  margin-top: var(--spacing-md);
  border-top: 1rpx solid var(--color-border);
}
</style>
