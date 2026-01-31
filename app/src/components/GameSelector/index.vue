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

export interface GameOption {
  id: number
  name: string
  icon?: string
}

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
  gap: 16rpx;
}

.game-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12rpx;
  padding: 20rpx 24rpx;
  background: var(--color-bg-secondary);
  border-radius: 16rpx;
  border: 2rpx solid var(--color-border);
  position: relative;
  min-width: 140rpx;
  transition: all 0.2s;
  
  &.selected {
    border-color: var(--color-primary);
    background: rgba(0, 210, 106, 0.08);
  }
}

.game-icon {
  width: 64rpx;
  height: 64rpx;
  border-radius: 12rpx;
}

.game-name {
  font-size: 26rpx;
  color: var(--color-text);
  font-weight: 500;
}

.check-mark {
  position: absolute;
  top: 8rpx;
  right: 8rpx;
  width: 32rpx;
  height: 32rpx;
  background: var(--color-primary);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
