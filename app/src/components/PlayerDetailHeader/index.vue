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
          <view v-if="player.gender" class="gender-badge" :class="player.gender">
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
    <view class="stats-row">
      <view class="stat-item">
        <text class="stat-value">{{ player.rating?.toFixed(1) || '5.0' }}</text>
        <text class="stat-label">评分</text>
      </view>
      <view class="stat-divider"></view>
      <view class="stat-item">
        <text class="stat-value">{{ player.orderCount || 0 }}</text>
        <text class="stat-label">接单数</text>
      </view>
      <view class="stat-divider"></view>
      <view class="stat-item">
        <text class="stat-value">{{ player.favoriteCount || 0 }}</text>
        <text class="stat-label">收藏数</text>
      </view>
      <view class="stat-divider"></view>
      <view class="stat-item">
        <text class="stat-value">{{ formatJoinDate }}</text>
        <text class="stat-label">入驻时间</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlAvatar from '@/components/gl/Avatar/index.vue'
import GlTag from '@/components/gl/Tag/index.vue'

export interface PlayerHeaderData {
  nickname: string
  avatar?: string
  coverImage?: string
  signature?: string
  gender?: 'male' | 'female'
  isOnline: boolean
  isVerified: boolean
  rating: number
  orderCount: number
  favoriteCount: number
  createdAt: string
}

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
</script>

<style lang="scss" scoped>
.header-card {
  position: relative;
  background: var(--color-bg-card);
  border-radius: 0 0 32rpx 32rpx;
  overflow: hidden;
  margin-bottom: 20rpx;
}

.player-cover {
  height: 280rpx;
  background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-primary-light, #4ADE80) 100%);
  
  .cover-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.player-basic {
  display: flex;
  gap: 24rpx;
  padding: 0 32rpx;
  margin-top: -60rpx;
  position: relative;
  z-index: 1;
}

.avatar-wrap {
  flex-shrink: 0;
}

.basic-info {
  flex: 1;
  padding-top: 70rpx;
  min-width: 0;
}

.name-row {
  display: flex;
  align-items: center;
  gap: 12rpx;
  flex-wrap: wrap;
  margin-bottom: 12rpx;
}

.nickname {
  font-size: 36rpx;
  font-weight: 700;
  color: var(--color-text);
}

.gender-badge {
  width: 36rpx;
  height: 36rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20rpx;
  
  &.male {
    background: rgba(59, 130, 246, 0.1);
    color: #3B82F6;
  }
  
  &.female {
    background: rgba(236, 72, 153, 0.1);
    color: #EC4899;
  }
}

.status-row {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.signature {
  font-size: 26rpx;
  color: var(--color-text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.stats-row {
  display: flex;
  padding: 32rpx;
  margin-top: 24rpx;
  border-top: 1rpx solid var(--color-border);
}

.stat-item {
  flex: 1;
  text-align: center;
}

.stat-value {
  display: block;
  font-size: 36rpx;
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: 8rpx;
}

.stat-label {
  font-size: 24rpx;
  color: var(--color-text-secondary);
}

.stat-divider {
  width: 1rpx;
  background: var(--color-border);
  margin: 8rpx 0;
}
</style>
