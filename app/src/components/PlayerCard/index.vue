<template>
  <view class="player-card" @click="$emit('click', player)">
    <!-- 头像 -->
    <view class="player-avatar">
      <image 
        v-if="player.avatar" 
        :src="player.avatar" 
        mode="aspectFill" 
      />
      <text v-else class="avatar-placeholder">{{ player.nickname?.[0] || '陪' }}</text>
      
      <!-- 在线状态 -->
      <view 
        class="online-dot" 
        :class="{ 
          online: player.isOnline, 
          busy: player.status === 'busy' 
        }"
      ></view>
    </view>

    <!-- 信息 -->
    <view class="player-info">
      <view class="player-header">
        <text class="player-name">{{ player.nickname }}</text>
        <view v-if="player.isVerified" class="verified-badge">
          <text>认证</text>
        </view>
      </view>
      
      <!-- 游戏标签 -->
      <view class="player-games">
        <view 
          v-for="game in displayGames" 
          :key="game.id" 
          class="game-tag"
        >
          <text>{{ game.name }}</text>
        </view>
        <text v-if="moreGamesCount > 0" class="more-games">
          +{{ moreGamesCount }}
        </text>
      </view>
      
      <!-- 评分和订单 -->
      <view class="player-stats">
        <view class="stat-item">
          <text class="stat-icon">⭐</text>
          <text class="stat-value">{{ player.rating?.toFixed(1) || '5.0' }}</text>
        </view>
        <view class="stat-item">
          <text class="stat-label">接单</text>
          <text class="stat-value">{{ player.orderCount || 0 }}</text>
        </view>
      </view>
    </view>

    <!-- 价格 -->
    <view class="player-price">
      <view class="price-row">
        <text class="price-value">¥{{ player.minPrice || 20 }}</text>
        <text class="price-unit">/局</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

export interface PlayerGame {
  id: number
  name: string
}

export interface Player {
  id: number
  nickname: string
  avatar?: string
  isOnline?: boolean
  isVerified?: boolean
  status?: 'online' | 'busy' | 'offline'
  rating?: number
  orderCount?: number
  minPrice?: number
  games?: PlayerGame[]
}

const props = defineProps<{
  player: Player
  maxGames?: number
}>()

defineEmits<{
  click: [player: Player]
}>()

const maxGames = props.maxGames || 2

const displayGames = computed(() => {
  return props.player.games?.slice(0, maxGames) || []
})

const moreGamesCount = computed(() => {
  const total = props.player.games?.length || 0
  return Math.max(0, total - maxGames)
})
</script>

<style lang="scss" scoped>
// 玩家卡片
.player-card {
  display: flex;
  align-items: center;
  padding: 24rpx 20rpx;
  background: var(--color-bg-card);
  border-radius: 20rpx;
  border: 2rpx solid var(--color-border);
  transition: all 0.2s ease;
  position: relative;
  box-sizing: border-box;
  width: 100%;
  
  &:active {
    transform: scale(0.98);
    border-color: var(--color-primary);
    box-shadow: 0 4rpx 16rpx rgba(0, 210, 106, 0.15);
  }
}

// 头像 - 增大尺寸
.player-avatar {
  position: relative;
  width: 100rpx;
  height: 100rpx;
  border-radius: 50%;
  margin-right: 20rpx;
  flex-shrink: 0;
  
  image {
    width: 100%;
    height: 100%;
    border-radius: 50%;
  }
  
  .avatar-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 40rpx;
    font-weight: 700;
    background: linear-gradient(135deg, var(--color-primary) 0%, #00E676 100%);
    color: #FFFFFF;
    border-radius: 50%;
  }
  
  .online-dot {
    position: absolute;
    right: 4rpx;
    bottom: 4rpx;
    width: 20rpx;
    height: 20rpx;
    border-radius: 50%;
    background: #9CA3AF;
    border: 3rpx solid var(--color-bg-card);
    
    &.online {
      background: #22C55E;
    }
    
    &.busy {
      background: #F59E0B;
    }
  }
}

// 信息区
.player-info {
  flex: 1;
  min-width: 0;
  overflow: hidden;
}

.player-header {
  display: flex;
  align-items: center;
  gap: 8rpx;
  margin-bottom: 8rpx;
  
  .player-name {
    font-size: 32rpx;
    font-weight: 600;
    color: var(--color-text);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  
  .verified-badge {
    flex-shrink: 0;
    padding: 4rpx 10rpx;
    background: rgba(59, 130, 246, 0.1);
    border-radius: 8rpx;
    
    text {
      font-size: 18rpx;
      color: #3B82F6;
    }
  }
}

// 游戏标签
.player-games {
  display: flex;
  align-items: center;
  gap: 8rpx;
  margin-bottom: 10rpx;
  overflow: hidden;
  
  .game-tag {
    flex-shrink: 0;
    padding: 4rpx 12rpx;
    background: rgba(0, 210, 106, 0.1);
    border-radius: 8rpx;
    
    text {
      font-size: 22rpx;
      color: var(--color-primary);
    }
  }
  
  .more-games {
    flex-shrink: 0;
    font-size: 22rpx;
    color: var(--color-text-secondary);
  }
}

// 统计区
.player-stats {
  display: flex;
  align-items: center;
  gap: 16rpx;
  
  .stat-item {
    display: flex;
    align-items: center;
    gap: 6rpx;
    
    .stat-icon {
      font-size: 24rpx;
    }
    
    .stat-label {
      font-size: 22rpx;
      color: var(--color-text-secondary);
    }
    
    .stat-value {
      font-size: 24rpx;
      font-weight: 600;
      color: var(--color-text);
    }
  }
}

// 价格区
.player-price {
  flex-shrink: 0;
  margin-left: 16rpx;
  text-align: right;
  
  .price-row {
    display: flex;
    align-items: baseline;
    justify-content: flex-end;
  }
  
  .price-value {
    font-size: 40rpx;
    font-weight: 800;
    color: #F97316;
    line-height: 1;
  }
  
  .price-unit {
    font-size: 22rpx;
    color: var(--color-text-secondary);
    margin-left: 4rpx;
  }
}
</style>
