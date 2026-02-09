<template>
  <BasePageLayout
    class="channel-list-page"
    :scroll="false"
    padding="0"
    title="公共频道"
    :show-back="true"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
  >
    <template #search>
      <!-- 搜索栏 -->
      <SearchBar
        v-model="searchKeyword"
        placeholder="搜索频道"
        :show-filter="true"
        @search="handleSearch"
        @filter="showFilter = true"
      />
    </template>

    <template #banner>
      <!-- 离线提示 -->
      <OfflineBanner
        :visible="isOffline"
        message="网络不可用，显示推荐频道"
        @action="refresh"
      />
    </template>

    <!-- 频道列表 -->
    <InfiniteList
      :state="pageState"
      :loading="loadingMore"
      :no-more="noMore"
      :error-message="errorMessage"
      empty-title="暂无频道"
      empty-desc="换个分类试试吧"
      padding="24rpx"
      @load-more="loadMore"
      @retry="refresh"
    >
      <view class="channel-grid">
        <ListItem
          v-for="(channel, index) in channels"
          :key="channel.id"
          :index="index"
          @click="enterChannel(channel)"
        >
          <ChannelCard
            :channel="channel"
            @join="joinChannel(channel)"
            @leave="leaveChannel(channel)"
          />
        </ListItem>
      </view>
    </InfiniteList>
    
    <!-- 筛选弹窗 -->
    <FilterPanel
      v-model:visible="showFilter"
      v-model="filterValues"
      :sections="filterSections"
      @apply="handleFilterApply"
      @reset="handleFilterReset"
    />
  </BasePageLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
import SearchBar from '@/components/SearchBar/index.vue'
import OfflineBanner from '@/components/OfflineBanner/index.vue'
import FilterPanel from '@/components/FilterPanel/index.vue'
import InfiniteList from '@/components/InfiniteList/index.vue'
import ListItem from '@/components/ListItem/index.vue'
// Business 组件
import ChannelCard from '@/components/ChannelCard/index.vue'
// Composables
import { useChannelList } from '@/composables/useChannelList'

const {
  channels,
  pageState,
  errorMessage,
  loadingMore,
  noMore,
  isOffline,
  searchKeyword,
  showFilter,
  filterValues,
  filterSections,
  loadMore,
  refresh,
  handleSearch,
  handleFilterApply,
  handleFilterReset,
  joinChannel,
  leaveChannel,
  enterChannel,
  goBack,
  init,
} = useChannelList()

onMounted(init)

onShow(() => {
  refresh()
})
</script>

<style lang="scss" scoped>
.channel-grid {
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
