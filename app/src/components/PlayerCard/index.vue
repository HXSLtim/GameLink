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
      <text class="price-value">¥{{ player.minPrice || 20 }}</text>
      <text class="price-unit">/局</text>
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
.player-card {
  display: flex;
  align-items: center;
  padding: 24rpx;
  background: var(--color-bg-card);
  border-radius: 24rpx;
  margin-bottom: 16rpx;
}

// 头像
.player-avatar {
  position: relative;
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  overflow: hidden;
  background: var(--color-bg-secondary);
  margin-right: 24rpx;
  flex-shrink: 0;
  
  image {
    width: 100%;
    height: 100%;
  }
  
  .avatar-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 48rpx;
    color: var(--color-text-placeholder);
    background: var(--color-primary);
    color: #FFFFFF;
  }
  
  .online-dot {
    position: absolute;
    right: 4rpx;
    bottom: 4rpx;
    width: 24rpx;
    height: 24rpx;
    border-radius: 50%;
    background: #9CA3AF;
    border: 4rpx solid var(--color-bg-card);
    
    &.online {
      background: #22C55E;
    }
    
    &.busy {
      background: #F59E0B;
    }
  }
}

// 信息
.player-info {
  flex: 1;
  min-width: 0;
}

.player-header {
  display: flex;
  align-items: center;
  margin-bottom: 12rpx;
  
  .player-name {
    font-size: 32rpx;
    font-weight: 600;
    color: var(--color-text);
    margin-right: 12rpx;
  }
  
  .verified-badge {
    padding: 2rpx 12rpx;
    background: rgba(59, 130, 246, 0.1);
    border-radius: 8rpx;
    
    text {
      font-size: 20rpx;
      color: #3B82F6;
    }
  }
}

// 游戏标签
.player-games {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8rpx;
  margin-bottom: 12rpx;
  
  .game-tag {
    padding: 4rpx 16rpx;
    background: var(--color-bg-secondary);
    border-radius: 8rpx;
    
    text {
      font-size: 22rpx;
      color: var(--color-text-secondary);
    }
  }
  
  .more-games {
    font-size: 22rpx;
    color: var(--color-text-placeholder);
  }
}

// 统计
.player-stats {
  display: flex;
  align-items: center;
  gap: 24rpx;
  
  .stat-item {
    display: flex;
    align-items: center;
    gap: 4rpx;
    
    .stat-icon {
      font-size: 24rpx;
    }
    
    .stat-label {
      font-size: 22rpx;
      color: var(--color-text-placeholder);
    }
    
    .stat-value {
      font-size: 24rpx;
      color: var(--color-text-secondary);
    }
  }
}

// 价格
.player-price {
  display: flex;
  align-items: baseline;
  flex-shrink: 0;
  margin-left: 16rpx;
  
  .price-value {
    font-size: 36rpx;
    font-weight: 600;
    color: var(--color-primary);
  }
  
  .price-unit {
    font-size: 22rpx;
    color: var(--color-text-secondary);
  }
}
</style>
