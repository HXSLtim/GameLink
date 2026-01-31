<template>
  <view class="section">
    <SectionHeader title="推荐陪玩师" @more="$emit('more')" />
    
    <!-- 骨架屏 -->
    <view v-if="loading && players.length === 0" class="player-list">
      <view v-for="i in 4" :key="i" class="player-skeleton">
        <view class="skeleton-avatar"></view>
        <view class="skeleton-info">
          <view class="skeleton-name"></view>
          <view class="skeleton-desc"></view>
        </view>
      </view>
    </view>
    
    <!-- 陪玩师列表 -->
    <view v-else-if="players.length > 0" class="player-list">
      <view 
        v-for="(player, index) in players" 
        :key="player.id" 
        class="player-card"
        :style="{ animationDelay: `${index * 0.05}s` }"
        @click="$emit('select', player)"
      >
        <view class="player-avatar">
          <GlAvatar :src="player.avatar" :text="player.nickname" size="large" :status="player.isOnline ? 'online' : undefined" />
        </view>
        <view class="player-info">
          <view class="player-name-row">
            <text class="player-name">{{ player.nickname }}</text>
            <GlTag v-if="player.rank" size="mini" type="warning">{{ player.rank }}</GlTag>
          </view>
          <view class="player-meta">
            <view class="rating">
              <uv-icon name="star-fill" size="12" color="#F59E0B"></uv-icon>
              <text>{{ player.rating?.toFixed(1) || '5.0' }}</text>
            </view>
            <text class="orders">{{ player.orderCount || 0 }}单</text>
            <text v-if="player.mainGame" class="game">{{ player.mainGame }}</text>
          </view>
        </view>
        <view class="player-price">
          <text class="price-value">¥{{ ((player.hourlyRate || 0) / 100).toFixed(0) }}</text>
          <text class="price-unit">/小时</text>
        </view>
      </view>
    </view>
    
    <!-- 空状态 -->
    <GlEmpty
      v-else
      title="暂无推荐陪玩师"
      description="刷新试试"
      :show-action="true"
      action-text="刷新"
      @action="$emit('refresh')"
    />
    
    <!-- 加载更多 -->
    <view v-if="players.length > 0" class="load-more">
      <text v-if="loading">加载中...</text>
      <text v-else-if="noMore">没有更多了</text>
      <text v-else class="load-more-btn" @click="$emit('load-more')">加载更多</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import SectionHeader from '@/components/SectionHeader/index.vue'
import GlAvatar from '@/components/gl/Avatar/index.vue'
import GlTag from '@/components/gl/Tag/index.vue'
import GlEmpty from '@/components/gl/Empty/index.vue'

export interface RecommendPlayerData {
  id: number
  nickname: string
  avatar?: string
  rank?: string
  rating: number
  hourlyRate: number
  isOnline: boolean
  orderCount: number
  mainGame?: string
}

interface Props {
  players: RecommendPlayerData[]
  loading?: boolean
  noMore?: boolean
}

withDefaults(defineProps<Props>(), {
  loading: false,
  noMore: false,
})

defineEmits<{
  more: []
  select: [player: RecommendPlayerData]
  refresh: []
  'load-more': []
}>()
</script>

<style lang="scss" scoped>
.section {
  padding: 24rpx;
}

.player-list {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.player-skeleton {
  display: flex;
  align-items: center;
  gap: 24rpx;
  padding: 24rpx;
  background: var(--color-bg-card);
  border-radius: 20rpx;
  
  .skeleton-avatar {
    width: 100rpx;
    height: 100rpx;
    border-radius: 50%;
    background: var(--color-bg-secondary);
    animation: pulse 1.5s infinite;
  }
  
  .skeleton-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 16rpx;
  }
  
  .skeleton-name {
    width: 120rpx;
    height: 32rpx;
    border-radius: 16rpx;
    background: var(--color-bg-secondary);
    animation: pulse 1.5s infinite;
  }
  
  .skeleton-desc {
    width: 200rpx;
    height: 24rpx;
    border-radius: 12rpx;
    background: var(--color-bg-secondary);
    animation: pulse 1.5s infinite;
  }
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.player-card {
  display: flex;
  align-items: center;
  gap: 24rpx;
  padding: 24rpx;
  background: var(--color-bg-card);
  border-radius: 20rpx;
  border: 2rpx solid var(--color-border);
  transition: all 0.2s;
  animation: fadeSlideUp 0.3s ease-out both;
  
  &:active {
    transform: scale(0.99);
    border-color: var(--color-primary);
  }
}

@keyframes fadeSlideUp {
  from {
    opacity: 0;
    transform: translateY(20rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.player-avatar {
  flex-shrink: 0;
}

.player-info {
  flex: 1;
  min-width: 0;
}

.player-name-row {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 12rpx;
}

.player-name {
  font-size: 32rpx;
  font-weight: 600;
  color: var(--color-text);
}

.player-meta {
  display: flex;
  align-items: center;
  gap: 20rpx;
  font-size: 24rpx;
  color: var(--color-text-secondary);
}

.rating {
  display: flex;
  align-items: center;
  gap: 4rpx;
}

.player-price {
  flex-shrink: 0;
  text-align: right;
}

.price-value {
  font-size: 36rpx;
  font-weight: 700;
  color: var(--color-primary);
}

.price-unit {
  font-size: 22rpx;
  color: var(--color-text-secondary);
}

.load-more {
  padding: 32rpx;
  text-align: center;
  font-size: 26rpx;
  color: var(--color-text-secondary);
}

.load-more-btn {
  color: var(--color-primary);
  
  &:active {
    opacity: 0.7;
  }
}
</style>
