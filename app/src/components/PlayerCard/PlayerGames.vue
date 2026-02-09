<template>
  <view class="player-games">
    <GlTag v-for="game in games" :key="game.key" size="mini" plain>
      {{ game.name }}
    </GlTag>
    <text v-if="moreCount > 0" class="more-games">+{{ moreCount }}</text>
  </view>
</template>

<script setup lang="ts">
import GlTag from '@/components/gl/Tag/index.vue'
import type { PlayerGameTag } from '@/types/player'

interface Props {
  games: PlayerGameTag[]
  moreCount: number
}

defineProps<Props>()
</script>

<style lang="scss" scoped>
.player-games {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  margin-bottom: var(--spacing-xs);
  flex-wrap: wrap;
}

.more-games {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}
</style>
