<template>
  <view class="player-list-page page-container">
    <!-- 搜索栏 -->
    <SearchBar
      v-model="searchKeyword"
      placeholder="搜索陪玩师"
      :show-filter="true"
      @search="handleSearch"
      @filter="showFilter = true"
    />
    
    <!-- 游戏分类 -->
    <GameGrid
      v-model="currentGameId"
      :games="games"
      @select="handleGameSelect"
    />
    
    <!-- 离线提示 -->
    <OfflineBanner
      :visible="isOffline"
      message="网络不可用，显示推荐陪玩师"
      @action="refresh"
    />
    
    <!-- 排序栏 -->
    <SortBar
      v-model="currentSort"
      :options="sortOptions"
      @change="handleSortChange"
    />
    
    <!-- 陪玩师列表 -->
    <InfiniteList
      :state="pageState"
      :loading="loadingMore"
      :no-more="noMore"
      :error-message="errorMessage"
      empty-title="暂无陪玩师"
      empty-desc="换个条件试试吧"
      padding="8rpx 12rpx"
      @load-more="loadMore"
      @retry="refresh"
    >
      <ListItem
        v-for="(player, index) in players"
        :key="player.id"
        :index="index"
        @click="goToDetail(player.id)"
      >
        <PlayerCard :player="player" />
      </ListItem>
    </InfiniteList>
    
    <!-- 筛选弹窗 -->
    <FilterPanel
      v-model:visible="showFilter"
      v-model="filterValues"
      :sections="filterSections"
      @apply="handleFilterApply"
      @reset="handleFilterReset"
    />
    
    <!-- 底部导航 -->
    <CustomTabBar :current="1" />
  </view>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
// Pattern 组件
import SearchBar from '@/components/SearchBar/index.vue'
import GameGrid from '@/components/GameGrid/index.vue'
import SortBar from '@/components/SortBar/index.vue'
import FilterPanel from '@/components/FilterPanel/index.vue'
import OfflineBanner from '@/components/OfflineBanner/index.vue'
import InfiniteList from '@/components/InfiniteList/index.vue'
import ListItem from '@/components/ListItem/index.vue'
// Business 组件
import PlayerCard from '@/components/PlayerCard/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
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
  currentGameId,
  currentSort,
  showFilter,
  filterValues,
  games,
  // 配置
  sortOptions,
  filterSections,
  // 方法
  loadMore,
  refresh,
  handleSearch,
  handleGameSelect,
  handleSortChange,
  handleFilterApply,
  handleFilterReset,
  goToDetail,
  init,
} = usePlayerList()

onMounted(init)
</script>

<style lang="scss" scoped>
.player-list-page {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  box-sizing: border-box;
  padding-bottom: calc(110rpx + env(safe-area-inset-bottom));
  
  @include desktop {
    height: 100vh;
    min-height: auto;
    padding-bottom: 0;
    overflow: hidden;
  }
}
</style>
