<template>
  <view class="favorite-item">
    <!-- 编辑模式选择框 -->
    <view 
      v-if="editMode" 
      class="checkbox"
      :class="{ checked: selected }"
      @tap.stop="$emit('toggle-select')"
    >
      <uv-icon v-if="selected" name="checkbox-mark" size="16" color="#fff"></uv-icon>
    </view>
    
    <!-- 陪玩师卡片 -->
    <view class="player-card" @tap="$emit('click')">
      <GlAvatar :src="player.avatar" :text="player.nickname" size="large" :status="player.isOnline ? 'online' : undefined" />
      
      <view class="player-info">
        <view class="player-header">
          <text class="player-name">{{ player.nickname }}</text>
          <GlTag v-if="player.isOnline" type="success" size="mini">在线</GlTag>
        </view>
        
        <view v-if="player.games?.length" class="player-games">
          <GlTag v-for="game in player.games.slice(0, 2)" :key="game" size="mini" plain>
            {{ game }}
          </GlTag>
        </view>
        
        <view class="player-stats">
          <text class="stat">⭐ {{ player.rating?.toFixed(1) || '5.0' }}</text>
          <text class="stat">接单 {{ player.orderCount || 0 }}</text>
        </view>
      </view>
      
      <view class="player-price">
        <text class="price-value">¥{{ player.minPrice || 20 }}</text>
        <text class="price-unit">起</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import GlAvatar from '@/components/gl/Avatar/index.vue'
import GlTag from '@/components/gl/Tag/index.vue'

export interface FavoritePlayerData {
  id: number
  playerId: number
  nickname: string
  avatar?: string
  isOnline: boolean
  rating: number
  orderCount: number
  minPrice: number
  games: string[]
}

interface Props {
  player: FavoritePlayerData
  editMode?: boolean
  selected?: boolean
}

withDefaults(defineProps<Props>(), {
  editMode: false,
  selected: false,
})

defineEmits<{
  click: []
  'toggle-select': []
}>()
</script>

<style lang="scss" scoped>
.favorite-item {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 24rpx;
  background: var(--color-bg-card);
  margin-bottom: 16rpx;
  border-radius: 20rpx;
  border: 2rpx solid var(--color-border);
}

.checkbox {
  width: 44rpx;
  height: 44rpx;
  border-radius: 50%;
  border: 2rpx solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.2s;
  
  &.checked {
    background: var(--color-primary);
    border-color: var(--color-primary);
  }
}

.player-card {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.player-info {
  flex: 1;
  min-width: 0;
}

.player-header {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 8rpx;
}

.player-name {
  font-size: 32rpx;
  font-weight: 600;
  color: var(--color-text);
}

.player-games {
  display: flex;
  gap: 8rpx;
  margin-bottom: 8rpx;
}

.player-stats {
  display: flex;
  gap: 20rpx;
}

.stat {
  font-size: 24rpx;
  color: var(--color-text-secondary);
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
</style>
