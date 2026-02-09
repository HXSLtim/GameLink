<template>
  <GlCard title="游戏认证" :shadow="false" bordered>
    <template #extra>
      <GlButton type="primary" size="mini" plain @click="$emit('add')">+ 添加</GlButton>
    </template>

    <GameCertItem
      v-for="(game, index) in games"
      :key="index"
      :game="game"
      @remove="$emit('remove', index)"
      @select-game="$emit('select-game', index)"
      @select-rank="$emit('select-rank', index)"
      @update:screenshot="(url) => $emit('update:screenshot', index, url)"
    />

    <GlEmpty v-if="games.length === 0" title="请添加至少一个游戏认证" compact />
  </GlCard>
</template>

<script setup lang="ts">
import GlCard from '@/components/gl/Card/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
import GlEmpty from '@/components/gl/Empty/index.vue'
import GameCertItem from '@/components/GameCertItem/index.vue'
import type { GameCertData } from '@/types/certification'

interface Props {
  games: GameCertData[]
}

defineProps<Props>()

defineEmits<{
  add: []
  remove: [index: number]
  'select-game': [index: number]
  'select-rank': [index: number]
  'update:screenshot': [index: number, url: string]
}>()
</script>
