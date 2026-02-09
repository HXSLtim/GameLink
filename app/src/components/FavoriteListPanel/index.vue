<template>
  <InfiniteList
    :state="pageState"
    :loading="loadingMore"
    :no-more="noMore"
    :empty-title="emptyTitle"
    :empty-desc="emptyDesc"
    :padding="padding"
    @load-more="$emit('load-more')"
    @retry="$emit('retry')"
  >
    <template #empty-action>
      <GlButton type="primary" size="small" @click="$emit('empty-action')">去看看</GlButton>
    </template>

    <view class="favorite-grid">
      <ListItem
        v-for="(item, index) in favorites"
        :key="item.id"
        :index="index"
        @click="$emit('item-click', item)"
      >
        <PlayerCard
          :player="item"
          variant="favorite"
          price-unit="起"
          :show-select="isEditMode"
          :selected="selectedIds.includes(item.id)"
          :clickable="false"
          @toggle-select="$emit('toggle-select', item.id)"
        />
      </ListItem>
    </view>
  </InfiniteList>
</template>

<script setup lang="ts">
import GlButton from '@/components/gl/Button/index.vue'
import InfiniteList from '@/components/InfiniteList/index.vue'
import ListItem from '@/components/ListItem/index.vue'
import PlayerCard from '@/components/PlayerCard/index.vue'
import type { FavoritePlayerData } from '@/types/player'
import type { PageStateType } from '@/types/page'

interface Props {
  favorites: FavoritePlayerData[]
  pageState: PageStateType
  loadingMore: boolean
  noMore: boolean
  isEditMode: boolean
  selectedIds: number[]
  emptyTitle?: string
  emptyDesc?: string
  padding?: string
}

withDefaults(defineProps<Props>(), {
  emptyTitle: '暂无收藏',
  emptyDesc: '去发现喜欢的陪玩师吧',
  padding: '24rpx',
})

defineEmits<{
  'load-more': []
  retry: []
  'empty-action': []
  'item-click': [item: FavoritePlayerData]
  'toggle-select': [id: number]
}>()
</script>

<style lang="scss" scoped>
.favorite-grid {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);

  @include desktop {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    column-gap: var(--spacing-sm);
    row-gap: var(--spacing-sm);
  }

  @include desktop-lg {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  :deep(.list-item) {
    margin-bottom: 0;
  }
}
</style>
