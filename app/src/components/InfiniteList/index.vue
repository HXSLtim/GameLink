<template>
  <PageState
    :state="state"
    :error-message="errorMessage"
    :empty-title="emptyTitle"
    :empty-desc="emptyDesc"
    :empty-icon="emptyIcon"
    @retry="$emit('retry')"
  >
    <scroll-view
      class="infinite-list"
      :class="{ 'infinite-list--horizontal': horizontal }"
      :scroll-y="!horizontal"
      :scroll-x="horizontal"
      :refresher-enabled="refresherEnabled"
      :refresher-triggered="refreshing"
      @scrolltolower="handleScrollToLower"
      @refresherrefresh="handleRefresh"
    >
      <view class="list-content" :style="contentStyle">
        <slot></slot>
      </view>
      
      <!-- 加载状态 -->
      <LoadMore
        v-if="showLoadMore"
        :loading="loading"
        :no-more="noMore"
        :loading-text="loadingText"
        :no-more-text="noMoreText"
      />
    </scroll-view>
  </PageState>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import PageState from '@/components/PageState/index.vue'
import LoadMore from '@/components/LoadMore/index.vue'
import type { PageStateType } from '@/types/page'

interface Props {
  // 状态相关
  state?: PageStateType
  loading?: boolean
  noMore?: boolean
  refreshing?: boolean
  errorMessage?: string
  
  // 空状态配置
  emptyTitle?: string
  emptyDesc?: string
  emptyIcon?: string
  
  // 功能开关
  refresherEnabled?: boolean
  horizontal?: boolean
  showLoadMore?: boolean
  
  // 文案
  loadingText?: string
  noMoreText?: string
  
  // 样式
  gap?: string
  padding?: string
}

const props = withDefaults(defineProps<Props>(), {
  state: 'content',
  loading: false,
  noMore: false,
  refreshing: false,
  emptyTitle: '暂无数据',
  emptyDesc: '',
  refresherEnabled: false,
  horizontal: false,
  showLoadMore: true,
  loadingText: '加载中...',
  noMoreText: '没有更多了',
  gap: 'var(--spacing-sm)',
  padding: 'var(--spacing-md)',
})

const emit = defineEmits<{
  loadMore: []
  refresh: []
  retry: []
}>()

const contentStyle = computed(() => ({
  padding: props.padding,
  gap: props.gap,
}))

const handleScrollToLower = () => {
  if (!props.loading && !props.noMore) {
    emit('loadMore')
  }
}

const handleRefresh = () => {
  emit('refresh')
}
</script>

<style lang="scss" scoped>
.infinite-list {
  flex: 1;
  height: 0;
  min-height: 0;
  
  &--horizontal {
    white-space: nowrap;
    
    .list-content {
      display: inline-flex;
      flex-direction: row;
    }
  }
}

.list-content {
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
}
</style>
