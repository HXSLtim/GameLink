<template>
  <GlCard title="选择游戏" required :shadow="false" bordered>
    <view class="games-grid">
      <view 
        v-for="game in games" 
        :key="game.id"
        class="game-option"
        :class="{ selected: modelValue === game.id }"
        @tap="$emit('update:modelValue', game.id)"
      >
        <image v-if="game.icon" :src="game.icon" mode="aspectFit" class="game-icon" />
        <text class="game-name">{{ game.name }}</text>
        <view v-if="modelValue === game.id" class="check-mark">
          <uv-icon name="checkbox-mark" size="14" color="#fff"></uv-icon>
        </view>
      </view>
    </view>
  </GlCard>
</template>

<script setup lang="ts">
import GlCard from '@/components/gl/Card/index.vue'
import type { GameOption } from '@/types/game'

interface Props {
  games: GameOption[]
  modelValue?: number
}

defineProps<Props>()

defineEmits<{
  'update:modelValue': [id: number]
}>()
</script>

<style lang="scss" scoped>
.games-grid {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm);
}

.game-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-xs);
  padding: var(--spacing-sm) var(--spacing-md);
  background: var(--color-bg-card);
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
  position: relative;
  min-width: 140rpx;
  transition: all 0.2s;
  cursor: pointer;
  @include press-effect;
  
  &.selected {
    border-color: var(--color-border);
    background: var(--color-bg-secondary);
  }
}

.game-icon {
  width: 56rpx;
  height: 56rpx;
  border-radius: var(--radius-sm);
}

.game-name {
  font-size: var(--font-xs);
  color: var(--color-text);
  font-weight: 500;
  @include text-ellipsis;
}

.check-mark {
  position: absolute;
  top: 8rpx;
  right: 8rpx;
  width: 28rpx;
  height: 28rpx;
  background: var(--color-primary);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
