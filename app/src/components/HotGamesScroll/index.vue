<template>
  <view class="section">
    <SectionHeader title="热门游戏" @more="$emit('more')" />
    
    <component
      :is="isPC ? 'view' : 'scroll-view'"
      class="game-scroll"
      :class="{ 'game-scroll--pc': isPC }"
      v-bind="isPC ? {} : { 'scroll-x': true }"
    >
      <!-- 骨架屏 -->
      <view v-if="loading" class="game-list">
        <Skeleton type="game" :rows="5" />
      </view>
      
      <!-- 游戏列表 -->
      <view v-else-if="games.length > 0" class="game-list">
        <view 
          v-for="(game, index) in games" 
          :key="game.id" 
          class="game-item"
          :style="{ animationDelay: `${index * 0.04}s` }"
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
        <GlEmpty title="暂无热门游戏" description="稍后再试" compact />
      </view>
    </component>
  </view>
</template>

<script setup lang="ts">
import SectionHeader from '@/components/SectionHeader/index.vue'
import GlEmpty from '@/components/gl/Empty/index.vue'
import Skeleton from '@/components/Skeleton/index.vue'
import { useDevice } from '@/composables/useDevice'
import { formatCount } from '@/utils/format'
import type { HotGameData } from '@/types/game'

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

const { isPC } = useDevice()

</script>

<style lang="scss" scoped>
.section {
  padding: var(--spacing-md);

  @include desktop {
    padding: 0 var(--spacing-lg);
  }
}

.game-scroll {
  margin: 0 calc(-1 * var(--spacing-md));
  padding: 0 var(--spacing-md);
  @include hide-scrollbar;
}

.game-scroll--pc {
  margin: 0;
  padding: 0;

  .game-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
    gap: var(--spacing-md);
    padding-right: 0;
  }

  .game-item {
    width: 100%;
    padding: var(--spacing-md) var(--spacing-sm);
  }

  .game-icon {
    width: 52px;
    height: 52px;
    border-radius: var(--radius-lg);
  }

  .game-icon--placeholder {
    text {
      font-size: 18px;
    }
  }

  .game-name {
    font-size: 13px;
    max-width: 100%;
  }

  .game-count {
    font-size: 11px;
  }

  :deep(.skeleton-game) {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
    gap: var(--spacing-md);
    padding-right: 0;
    width: 100%;
  }

  :deep(.skeleton-game-item) {
    width: 100%;
  }
}

.game-list {
  display: flex;
  gap: var(--spacing-md);
  padding-right: var(--spacing-md);
  padding-bottom: var(--spacing-sm);
}

.game-item {
  flex-shrink: 0;
  width: 140rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-xs);
  padding: var(--spacing-sm) var(--spacing-xs);
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1rpx solid var(--color-border);
  cursor: pointer;
  transition: transform 0.25s ease, box-shadow 0.25s ease, border-color 0.25s ease;
  animation: fadeSlideUp 0.35s ease-out both;
  @include press-effect;

  &:hover {
    transform: translateY(-3px);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
    border-color: rgba(var(--color-primary-rgb), 0.4);
  }
}

.game-icon {
  width: 88rpx;
  height: 88rpx;
  border-radius: var(--radius-md);
  background: var(--color-bg-secondary);
  border: 1rpx solid var(--color-border);
  
  &--placeholder {
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, rgba(var(--color-primary-rgb), 0.08), rgba(var(--color-primary-rgb), 0.18));
    border-color: rgba(var(--color-primary-rgb), 0.15);
    
    text {
      font-size: var(--font-md);
      font-weight: 700;
      color: var(--color-primary);
    }
  }
}

.game-name {
  font-size: var(--font-xs);
  font-weight: 600;
  color: var(--color-text);
  text-align: center;
  @include text-ellipsis;
  max-width: 130rpx;
}

.game-count {
  font-size: var(--font-xs);
  color: var(--color-text-placeholder);
}

.empty-games {
  padding: var(--spacing-lg);
  text-align: center;
}

@keyframes fadeSlideUp {
  from {
    opacity: 0;
    transform: translateY(16rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
