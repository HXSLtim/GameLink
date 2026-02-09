<template>
  <SectionCard title="擅长游戏">
    <view class="games-list">
      <view v-for="game in games" :key="game.id" class="game-item">
        <image v-if="game.icon" :src="game.icon" mode="aspectFit" class="game-icon" />
        <view v-else class="game-icon game-icon--placeholder">
          <text>{{ game.name?.[0] || '游' }}</text>
        </view>
        <view class="game-info">
          <text class="game-name">{{ game.name }}</text>
          <GlTag v-if="game.rankName" size="mini" type="warning">{{ game.rankName }}</GlTag>
        </view>
        <view class="game-price">
          <PriceTag :amount="game.price" amount-unit="yuan" unit="局" size="small" />
        </view>
      </view>
      
      <GlEmpty v-if="!games?.length" title="暂无游戏信息" compact />
    </view>
  </SectionCard>
</template>

<script setup lang="ts">
import SectionCard from '@/components/SectionCard/index.vue'
import PriceTag from '@/components/PriceTag/index.vue'
import GlTag from '@/components/gl/Tag/index.vue'
import GlEmpty from '@/components/gl/Empty/index.vue'
import type { PlayerGameData } from '@/types/player'

interface Props {
  games: PlayerGameData[]
}

defineProps<Props>()
</script>

<style lang="scss" scoped>
.games-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.game-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-sm);
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1rpx solid var(--color-border);
}

.game-icon {
  width: 72rpx;
  height: 72rpx;
  border-radius: var(--radius-md);
  border: 1rpx solid var(--color-border);
  flex-shrink: 0;
  
  &--placeholder {
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--color-bg-secondary);
    
    text {
      font-size: var(--font-base);
      font-weight: 600;
      color: var(--color-text);
    }
  }
}

.game-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.game-name {
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
  @include text-ellipsis;
}

.game-price {
  flex-shrink: 0;
  text-align: right;
}
</style>
