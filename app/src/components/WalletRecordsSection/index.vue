<template>
  <SectionCard title="交易记录" class="records-card" margin="0 var(--spacing-md)">
    <template #extra>
      <TabsBar
        v-model="currentTab"
        :tabs="tabs"
        @change="(key) => $emit('change', key)"
      />
    </template>

    <InfiniteList
      :state="listState"
      :loading="loadingMore"
      :no-more="noMore"
      :show-load-more="records.length > 0"
      empty-title="暂无交易记录"
      padding="0"
      @load-more="$emit('load-more')"
    >
      <TransactionItem
        v-for="record in records"
        :key="record.id"
        :record="record"
      />
    </InfiniteList>
  </SectionCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import SectionCard from '@/components/SectionCard/index.vue'
import TabsBar from '@/components/TabsBar/index.vue'
import InfiniteList from '@/components/InfiniteList/index.vue'
import TransactionItem from '@/components/TransactionItem/index.vue'
import type { TransactionData } from '@/types/wallet'
import type { TabItem } from '@/types/ui'

interface Props {
  loading: boolean
  loadingMore: boolean
  noMore: boolean
  records: TransactionData[]
  tabs: TabItem[]
  modelValue: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'load-more': []
  change: [value: string]
}>()

const currentTab = computed({
  get: () => props.modelValue,
  set: (value: string) => emit('update:modelValue', value),
})

const listState = computed(() => {
  if (props.loading) return 'loading'
  return props.records.length === 0 ? 'empty' : 'content'
})
</script>

<style lang="scss" scoped>
.records-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
</style>
