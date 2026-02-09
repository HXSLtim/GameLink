<template>
  <SectionCard title="收益明细" margin="0 var(--spacing-md)">
    <template #extra>
      <view class="filter-btn" @tap="$emit('filter')">
        <uv-icon name="list" size="16" color="var(--color-text-secondary)"></uv-icon>
        <text>筛选</text>
      </view>
    </template>

    <InfiniteList
      :state="listState"
      :loading="loadingMore"
      :no-more="noMore"
      :show-load-more="items.length > 0"
      empty-title="暂无收益记录"
      padding="0"
      @load-more="$emit('load-more')"
    >
      <EarningsItem
        v-for="item in items"
        :key="item.id"
        :type="item.type"
        :title="item.title"
        :description="item.description"
        :amount="item.amount"
        :created-at="item.createdAt"
        @click="$emit('item-click', item)"
      />
    </InfiniteList>
  </SectionCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import SectionCard from '@/components/SectionCard/index.vue'
import InfiniteList from '@/components/InfiniteList/index.vue'
import EarningsItem from '@/components/EarningsItem/index.vue'
import type { EarningsItem as EarningsListItem } from '@/types/earnings'

interface Props {
  loading: boolean
  loadingMore: boolean
  noMore: boolean
  items: EarningsListItem[]
}

const props = defineProps<Props>()

defineEmits<{
  filter: []
  'load-more': []
  'item-click': [item: EarningsListItem]
}>()

const listState = computed(() => {
  if (props.loading) return 'loading'
  return props.items.length === 0 ? 'empty' : 'content'
})
</script>

<style lang="scss" scoped>
.filter-btn {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  padding: 2rpx var(--spacing-sm);
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
  background: var(--color-bg-secondary);
  border: 1rpx solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  @include press-effect;
}
</style>
