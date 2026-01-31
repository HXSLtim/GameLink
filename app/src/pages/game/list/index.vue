<template>
  <view class="game-list-page page-container">
    <!-- 顶部导航 -->
    <NavBar title="游戏列表" @back="goBack" />

    <!-- 搜索栏 -->
    <SearchBar
      v-model="searchKeyword"
      placeholder="搜索游戏"
      :clearable="true"
      @search="handleSearch"
      @clear="clearSearch"
    />

    <!-- 分类标签 -->
    <TabsBar
      v-model="currentCategory"
      :tabs="categoryTabs"
      scrollable
      @change="selectCategory"
    />

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
        <GameCardLarge
          v-for="game in filteredGames"
          :key="game.id"
          :game="game"
          @click="goToGame(game)"
        />
      </view>
    </InfiniteList>

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import SearchBar from '@/components/SearchBar/index.vue'
import TabsBar from '@/components/TabsBar/index.vue'
import InfiniteList from '@/components/InfiniteList/index.vue'
// Business 组件
import GameCardLarge from '@/components/GameCardLarge/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// Composables
import { useGameList } from '@/composables/useGameList'

const {
  filteredGames,
  pageState,
  errorMessage,
  loadingMore,
  noMore,
  searchKeyword,
  currentCategory,
  categories,
  loadMore,
  refresh,
  handleSearch,
  clearSearch,
  selectCategory,
  goToGame,
  goBack,
  init,
} = useGameList()

// 转换为 TabsBar 需要的格式
const categoryTabs = computed(() => 
  categories.value.map(c => ({ key: c.id, label: c.name }))
)

onMounted(init)

onShow(() => {
  refresh()
})
</script>

<style lang="scss" scoped>
.game-list-page {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  box-sizing: border-box;
  
  @include desktop {
    height: 100vh;
    min-height: auto;
    overflow: hidden;
  }
}

.games-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20rpx;
  
  @include desktop {
    grid-template-columns: repeat(3, 1fr);
  }
}
</style>
