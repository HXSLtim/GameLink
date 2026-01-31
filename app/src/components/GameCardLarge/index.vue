<template>
  <view class="game-card" @tap="$emit('click')">
    <view class="game-cover">
      <image v-if="game.coverImage" :src="game.coverImage" mode="aspectFill" />
      <view v-else class="cover-placeholder">
        <text>{{ game.name?.[0] }}</text>
      </view>
      <GlTag v-if="game.isHot" type="error" size="mini" class="hot-badge">热门</GlTag>
    </view>
    
    <view class="game-info">
      <text class="game-name">{{ game.name }}</text>
      <view class="game-stats">
        <text class="player-count">{{ formatCount(game.playerCount) }}位陪玩</text>
      </view>
      <view class="price-range">
        <text>¥{{ game.minPrice }}-{{ game.maxPrice }}/局</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import GlTag from '@/components/gl/Tag/index.vue'

export interface GameData {
  id: number
  name: string
  coverImage?: string
  categoryId?: string
  isHot?: boolean
  playerCount: number
  minPrice: number
  maxPrice: number
}

interface Props {
  game: GameData
}

defineProps<Props>()

defineEmits<{
  click: []
}>()

const formatCount = (count: number) => {
  if (count >= 10000) {
    return `${(count / 10000).toFixed(1)}万`
  }
  if (count >= 1000) {
    return `${(count / 1000).toFixed(1)}k`
  }
  return count.toString()
}
</script>

<style lang="scss" scoped>
.game-card {
  background: var(--color-bg-card);
  border-radius: 20rpx;
  overflow: hidden;
  border: 2rpx solid var(--color-border);
  transition: all 0.2s;
  
  &:active {
    transform: scale(0.98);
    border-color: var(--color-primary);
  }
}

.game-cover {
  position: relative;
  width: 100%;
  aspect-ratio: 4 / 3;
  background: var(--color-bg-secondary);
  
  image {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.cover-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-primary-light, #4ADE80) 100%);
  
  text {
    font-size: 64rpx;
    font-weight: 800;
    color: #FFFFFF;
  }
}

.hot-badge {
  position: absolute;
  top: 12rpx;
  right: 12rpx;
}

.game-info {
  padding: 20rpx;
}

.game-name {
  font-size: 30rpx;
  font-weight: 600;
  color: var(--color-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: block;
  margin-bottom: 8rpx;
}

.game-stats {
  margin-bottom: 8rpx;
}

.player-count {
  font-size: 24rpx;
  color: var(--color-text-secondary);
}

.price-range {
  text {
    font-size: 26rpx;
    font-weight: 600;
    color: var(--color-primary);
  }
}
</style>
