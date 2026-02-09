<template>
  <view class="services-list">
    <template v-if="loading">
      <Skeleton type="list" :rows="3" />
    </template>

    <template v-else-if="services.length > 0">
      <PlayerServiceCard
        v-for="service in services"
        :key="service.id"
        :service="service"
        @toggle-status="$emit('toggle-status', service)"
        @edit="$emit('edit', service)"
        @delete="$emit('delete', service)"
      />
    </template>

    <GlEmpty
      v-else
      title="暂无服务"
      description="添加服务开始接单吧"
      action-text="添加服务"
      @action="$emit('add')"
    />
  </view>
</template>

<script setup lang="ts">
import GlEmpty from '@/components/gl/Empty/index.vue'
import Skeleton from '@/components/Skeleton/index.vue'
import PlayerServiceCard from '@/components/PlayerServiceCard/index.vue'
import type { PlayerServiceCardData } from '@/types/player'

interface Props {
  loading: boolean
  services: PlayerServiceCardData[]
}

defineProps<Props>()

defineEmits<{
  'toggle-status': [service: PlayerServiceCardData]
  edit: [service: PlayerServiceCardData]
  delete: [service: PlayerServiceCardData]
  add: []
}>()
</script>

<style lang="scss" scoped>
.services-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
  padding: 0 var(--spacing-md);

  @include desktop {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    column-gap: var(--spacing-sm);
    row-gap: var(--spacing-sm);
  }

  @include desktop-lg {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}
</style>
