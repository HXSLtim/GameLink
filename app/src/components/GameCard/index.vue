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
      <view class="game-title-row">
        <text class="game-name">{{ game.name }}</text>
        <view class="game-price">
          <PriceTag
            :amount="game.minPrice"
            amount-unit="yuan"
            size="small"
            :show-decimal="false"
          />
          <text class="price-sep">-</text>
          <PriceTag
            :amount="game.maxPrice"
            amount-unit="yuan"
            size="small"
            :show-decimal="false"
            :show-currency="false"
          />
          <text class="price-unit">/局</text>
        </view>
      </view>
      <text class="game-meta">{{ formatCount(game.playerCount) }}位陪玩</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import GlTag from '@/components/gl/Tag/index.vue'
import PriceTag from '@/components/PriceTag/index.vue'
import { formatCount } from '@/utils/format'
import type { GameCardData } from '@/types/game'

interface Props {
  game: GameCardData
}

defineProps<Props>()

defineEmits<{
  click: []
}>()

</script>

<style lang="scss" scoped>
.game-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  overflow: hidden;
  border: 1rpx solid var(--color-border);
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-md);
    border-color: var(--color-primary);
  }
  
  &:active {
    transform: translateY(0);
    background: var(--color-bg-secondary);
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
  background: var(--color-bg-secondary);
  
  text {
    font-size: var(--font-xl);
    font-weight: 700;
    color: var(--color-text-secondary);
  }
}

.hot-badge {
  position: absolute;
  top: 12rpx;
  right: 12rpx;
}

.game-info {
  padding: var(--spacing-sm);
}

.game-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-xs);
  margin-bottom: var(--spacing-xs);
}

.game-name {
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: block;
  flex: 1;
  min-width: 0;
}

.game-price {
  display: inline-flex;
  align-items: baseline;
  gap: 4rpx;
  flex-shrink: 0;
  font-size: var(--font-sm);
}

.game-price :deep(.price-tag) {
  color: var(--color-text);
}

.price-sep {
  color: var(--color-text);
  font-weight: 600;
}

.price-unit {
  color: var(--color-text-secondary);
  font-weight: 500;
  font-size: var(--font-xs);
}

.game-meta {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}
</style>
