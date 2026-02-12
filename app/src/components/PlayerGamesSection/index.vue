<template>
  <SectionCard title="擅长游戏">
    <view class="games-grid">
      <view v-for="game in games" :key="game.id" class="game-item">
        <view class="game-cover-wrap">
          <image v-if="game.icon" :src="game.icon" mode="aspectFill" class="game-cover" />
          <view v-else class="game-cover-placeholder">
            <uv-icon name="gamepad" size="28" color="var(--color-text-secondary)" />
          </view>

          <!-- Rank Badge (Glass effect) -->
          <view v-if="game.rankName" class="game-rank-badge">
            <text>{{ game.rankName }}</text>
          </view>
        </view>

        <view class="game-info">
          <text class="game-name">{{ game.name }}</text>
          <view class="game-meta">
            <PriceTag :amount="game.price" amount-unit="yuan" unit="局" size="small" :show-symbol="false" />
          </view>
        </view>
      </view>

      <GlEmpty v-if="!games?.length" title="暂无游戏信息" compact />
    </view>
  </SectionCard>
</template>

<script setup lang="ts">
import SectionCard from '@/components/SectionCard/index.vue'
import PriceTag from '@/components/PriceTag/index.vue'
import GlEmpty from '@/components/gl/Empty/index.vue'
import type { PlayerGameData } from '@/types/player'

interface Props {
  games: PlayerGameData[]
}

defineProps<Props>()
</script>

<style lang="scss" scoped>
.games-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr); // 移动端双列
  gap: var(--spacing-sm);

  // PC 端三列
  @include desktop {
    grid-template-columns: repeat(3, 1fr);
    gap: var(--spacing-md);
  }
}

.game-item {
  position: relative;
  background: var(--color-bg-card);
  border-radius: var(--radius-lg);
  overflow: hidden;
  border: 1rpx solid var(--color-border);
  transition: all 0.3s ease;
  cursor: pointer;

  &:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-md);
    border-color: var(--color-primary-light);
  }
}

.game-cover-wrap {
  height: 140rpx;
  position: relative;
  background: var(--color-bg-secondary);
}

.game-cover {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.game-cover-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-secondary);
}

.game-rank-badge {
  position: absolute;
  bottom: 8rpx;
  right: 8rpx;
  padding: 4rpx 12rpx;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  border-radius: var(--radius-sm);

  text {
    font-size: 20rpx;
    color: #fff;
    font-weight: 600;
  }
}

.game-info {
  padding: var(--spacing-sm);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.game-name {
  font-size: var(--font-sm);
  font-weight: 600;
  color: var(--color-text);
  @include text-ellipsis;
  flex: 1;
}

.game-meta {
  flex-shrink: 0;
}
</style>
