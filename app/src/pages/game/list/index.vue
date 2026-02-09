<template>
  <BasePageLayout
    class="game-list-page"
    :scroll="false"
    padding="0"
    title="游戏列表"
    :show-back="true"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
  >
    <template #search>
      <!-- 搜索栏 -->
      <SearchBar
        v-model="searchKeyword"
        placeholder="搜索游戏"
        :clearable="true"
        :show-filter="true"
        @search="handleSearch"
        @clear="clearSearch"
        @filter="showFilter = true"
      />
    </template>

    <!-- 游戏列表 -->
    <InfiniteList
      :state="pageState"
      :loading="loadingMore"
      :no-more="noMore"
      :error-message="errorMessage"
      empty-title="暂无游戏"
      :empty-desc="searchKeyword ? '换个关键词试试' : '暂时没有游戏'"
      padding="24rpx"
      @load-more="loadMore"
      @retry="refresh"
    >
      <view class="games-grid">
        <GameCard
          v-for="(game, index) in filteredGames"
          :key="game.id"
          class="game-grid-item"
          :game="game"
          :style="{ animationDelay: `${(index % 8) * 0.04}s` }"
          @click="goToGame(game)"
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
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
import SearchBar from '@/components/SearchBar/index.vue'
import FilterPanel from '@/components/FilterPanel/index.vue'
import InfiniteList from '@/components/InfiniteList/index.vue'
// Business 组件
import GameCard from '@/components/GameCard/index.vue'
// Composables
import { useGameList } from '@/composables/useGameList'

const {
  filteredGames,
  pageState,
  errorMessage,
  loadingMore,
  noMore,
  searchKeyword,
  showFilter,
  filterValues,
  filterSections,
  loadMore,
  refresh,
  handleSearch,
  clearSearch,
  handleFilterApply,
  handleFilterReset,
  goToGame,
  goBack,
  init,
} = useGameList()

onMounted(init)

onShow(() => {
  refresh()
})
</script>

<style lang="scss" scoped>
.games-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-sm);
  
  @include desktop {
    grid-template-columns: repeat(3, 1fr);
    gap: var(--spacing-md);
  }

  @include desktop-lg {
    grid-template-columns: repeat(4, 1fr);
  }
}

.game-grid-item {
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
