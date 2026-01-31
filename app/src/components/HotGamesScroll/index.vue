<template>
  <view class="section">
    <SectionHeader title="热门游戏" @more="$emit('more')" />
    
    <scroll-view class="game-scroll" scroll-x>
      <!-- 骨架屏 -->
      <view v-if="loading" class="game-list">
        <view v-for="i in 5" :key="i" class="game-item game-item--skeleton">
          <view class="game-icon-skeleton"></view>
          <view class="game-name-skeleton"></view>
        </view>
      </view>
      
      <!-- 游戏列表 -->
      <view v-else-if="games.length > 0" class="game-list">
        <view 
          v-for="game in games" 
          :key="game.id" 
          class="game-item"
          @click="$emit('select', game)"
        >
          <image v-if="game.icon" :src="game.icon" class="game-icon" mode="aspectFill" />
          <view v-else class="game-icon game-icon--placeholder">
            <text>{{ game.name?.[0] || '游' }}</text>
          </view>
          <text class="game-name">{{ game.name }}</text>
          <text class="game-count">{{ formatCount(game.playerCount) }}人在玩</text>
        </view>
      </view>
      
      <!-- 空状态 -->
      <view v-else class="empty-games">
        <text>暂无热门游戏</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import SectionHeader from '@/components/SectionHeader/index.vue'

export interface HotGameData {
  id: number
  name: string
  icon?: string
  playerCount?: number
}

interface Props {
  games: HotGameData[]
  loading?: boolean
}

withDefaults(defineProps<Props>(), {
  loading: false,
})

defineEmits<{
  more: []
  select: [game: HotGameData]
}>()

const formatCount = (count?: number) => {
  if (!count) return 0
  if (count >= 10000) return `${(count / 10000).toFixed(1)}万`
  if (count >= 1000) return `${(count / 1000).toFixed(1)}k`
  return count
}
</script>

<style lang="scss" scoped>
.section {
  padding: 24rpx;
}

.game-scroll {
  margin: 0 -24rpx;
  padding: 0 24rpx;
  
  &::-webkit-scrollbar {
    display: none;
  }
}

.game-list {
  display: flex;
  gap: 24rpx;
  padding-right: 24rpx;
}

.game-item {
  flex-shrink: 0;
  width: 160rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12rpx;
  padding: 20rpx;
  background: var(--color-bg-card);
  border-radius: 20rpx;
  border: 2rpx solid var(--color-border);
  transition: all 0.2s;
  
  &:active {
    transform: scale(0.96);
    border-color: var(--color-primary);
  }
  
  &--skeleton {
    .game-icon-skeleton {
      width: 80rpx;
      height: 80rpx;
      border-radius: 20rpx;
      background: var(--color-bg-secondary);
      animation: pulse 1.5s infinite;
    }
    
    .game-name-skeleton {
      width: 80rpx;
      height: 24rpx;
      border-radius: 12rpx;
      background: var(--color-bg-secondary);
      animation: pulse 1.5s infinite;
    }
  }
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.game-icon {
  width: 80rpx;
  height: 80rpx;
  border-radius: 20rpx;
  background: var(--color-bg-secondary);
  
  &--placeholder {
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-primary-light, #4ADE80) 100%);
    
    text {
      font-size: 32rpx;
      font-weight: 700;
      color: #FFFFFF;
    }
  }
}

.game-name {
  font-size: 26rpx;
  font-weight: 600;
  color: var(--color-text);
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 140rpx;
}

.game-count {
  font-size: 20rpx;
  color: var(--color-text-secondary);
}

.empty-games {
  padding: 40rpx;
  text-align: center;
  color: var(--color-text-secondary);
  font-size: 26rpx;
}
</style>
