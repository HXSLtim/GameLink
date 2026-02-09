<template>
  <view class="player-meta">
    <view v-if="showRating" class="stat stat-rating">
      <uv-icon name="star-fill" size="12" color="var(--color-gold)" />
      <text>{{ rating }}</text>
    </view>
    <text v-if="showOrderCount" class="stat">接单 {{ orderCount }}</text>
    <text v-if="showMainGame && mainGame" class="stat">{{ mainGame }}</text>
  </view>
</template>

<script setup lang="ts">
interface Props {
  showRating: boolean
  rating: string
  showOrderCount: boolean
  orderCount: number
  showMainGame: boolean
  mainGame?: string
}

defineProps<Props>()
</script>

<style lang="scss" scoped>
.player-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  column-gap: var(--spacing-md);
  row-gap: var(--spacing-xs);
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}

.stat {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}

.stat-rating {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-xs);
  line-height: 1;
}
</style>
