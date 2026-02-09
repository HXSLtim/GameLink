<template>
  <view class="section">
    <SectionHeader title="推荐陪玩师" @more="$emit('more')" />
    
    <!-- 骨架屏 -->
    <view v-if="loading && players.length === 0" class="player-list">
      <Skeleton v-for="i in 4" :key="i" type="card" />
    </view>
    
    <!-- 陪玩师列表 -->
    <view v-else-if="players.length > 0" class="player-list">
      <PlayerCard
        v-for="(player, index) in players"
        :key="player.id"
        class="recommend-card"
        :player="player"
        variant="grid"
        :clickable="true"
        :style="{ animationDelay: `${index * 0.05}s` }"
        @click="$emit('select', $event)"
      />
    </view>
    
    <!-- 空状态 -->
    <GlEmpty
      v-else
      title="暂无推荐陪玩师"
      description="刷新试试"
      :show-action="true"
      action-text="刷新"
      @action="$emit('refresh')"
    />
  </view>
</template>

<script setup lang="ts">
import SectionHeader from '@/components/SectionHeader/index.vue'
import GlEmpty from '@/components/gl/Empty/index.vue'
import Skeleton from '@/components/Skeleton/index.vue'
import PlayerCard from '@/components/PlayerCard/index.vue'
import type { RecommendPlayerData } from '@/types/player'

interface Props {
  players: RecommendPlayerData[]
  loading?: boolean
}

withDefaults(defineProps<Props>(), {
  loading: false,
})

defineEmits<{
  more: []
  select: [player: RecommendPlayerData]
  refresh: []
}>()
</script>

<style lang="scss" scoped>
.section {
  padding: var(--spacing-md);

  @include desktop {
    padding: 0 var(--spacing-lg);
  }
}

.player-list {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: var(--spacing-sm);

  @include desktop {
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: var(--spacing-md);
  }

  @include desktop-lg {
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: 20px;
  }
}

.recommend-card {
  animation: fadeSlideUp 0.3s ease-out both;
}

@keyframes fadeSlideUp {
  from {
    opacity: 0;
    transform: translateY(20rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
