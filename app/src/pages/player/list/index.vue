<template>
  <BasePageLayout
    class="player-list-page"
    :scroll="false"
    padding="0"
    title="陪玩师"
    :show-back="false"
    :show-tab-bar="true"
    :tab-bar-current="1"
  >
    <template #search>
      <!-- 搜索栏 -->
      <SearchBar
        v-model="searchKeyword"
        placeholder="搜索陪玩师"
        :show-filter="true"
        @search="handleSearch"
        @filter="showFilter = true"
      />
    </template>

    <template #banner>
      <!-- 离线提示 -->
      <OfflineBanner
        :visible="isOffline"
        message="网络不可用，显示推荐陪玩师"
        @action="refresh"
      />
    </template>
    
    <!-- 陪玩师列表 -->
    <InfiniteList
      :state="pageState"
      :loading="loadingMore"
      :no-more="noMore"
      :error-message="errorMessage"
      empty-title="暂无陪玩师"
      empty-desc="换个条件试试吧"
      padding="var(--spacing-md)"
      @load-more="loadMore"
      @retry="refresh"
    >
      <view class="player-grid">
        <PlayerCard
          v-for="(player, index) in players"
          :key="player.id"
          class="player-grid-item"
          :player="player"
          variant="grid"
          :clickable="true"
          :style="{ animationDelay: `${(index % 8) * 0.04}s` }"
          @click="goToDetail(player.id)"
        />
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
// Pattern 组件
import SearchBar from '@/components/SearchBar/index.vue'
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
import FilterPanel from '@/components/FilterPanel/index.vue'
import OfflineBanner from '@/components/OfflineBanner/index.vue'
import InfiniteList from '@/components/InfiniteList/index.vue'
// Business 组件
import PlayerCard from '@/components/PlayerCard/index.vue'
// Composables
import { usePlayerList } from '@/composables/usePlayerList'

// 使用陪玩师列表 Hook - 所有业务逻辑都封装在这里
const {
  // 数据
  players,
  pageState,
  errorMessage,
  loadingMore,
  noMore,
  isOffline,
  // 筛选状态
  searchKeyword,
  showFilter,
  filterValues,
  // 配置
  filterSections,
  // 方法
  loadMore,
  refresh,
  handleSearch,
  handleFilterApply,
  handleFilterReset,
  goToDetail,
  init,
} = usePlayerList()

onMounted(init)
</script>

<style lang="scss" scoped>
.player-list-page {
  padding-bottom: calc(110rpx + env(safe-area-inset-bottom));

  @include desktop {
    padding-bottom: 0;
  }
}

.player-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: var(--spacing-md);

  @include desktop {
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: var(--spacing-lg);
  }

  @include desktop-lg {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }
}

.player-grid-item {
  animation: fadeSlideUp 0.3s ease-out both;
}

@keyframes fadeSlideUp {
  from {
    opacity: 0;
    transform: translateY(16rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
