<template>
  <view class="home-page page-container">
    <!-- 顶部状态栏占位 -->
    <view class="status-bar safe-area-top"></view>
    
    <!-- 头部 -->
    <HomeHeader
      :is-logged-in="isLoggedIn"
      :avatar="userInfo?.avatar"
      :nickname="userInfo?.nickname"
      @login="goToLogin"
      @profile="goToProfile"
    />

    <!-- 搜索栏 -->
    <SearchBar
      placeholder="搜索陪玩师、游戏..."
      :readonly="true"
      @click="goToPlayerList"
    />

    <!-- 离线提示 -->
    <OfflineBanner
      :visible="isOfflineMode"
      message="网络不可用，显示缓存数据"
      @action="refreshAll"
    />

    <!-- 主内容区 -->
    <scroll-view class="main-scroll" scroll-y @scrolltolower="loadMorePlayers">
      <!-- 热门游戏 -->
      <HotGamesScroll
        :games="hotGames"
        :loading="gamesLoading"
        @more="goToGameList"
        @select="(game) => goToGamePlayers(game.id)"
      />

      <!-- 推荐陪玩师 -->
      <RecommendPlayersSection
        :players="recommendPlayers"
        :loading="playersLoading"
        :no-more="noMorePlayers"
        @more="goToPlayerList"
        @select="(player) => goToPlayerDetail(player.id)"
        @refresh="refreshAll"
        @load-more="loadMorePlayers"
      />
    </scroll-view>
    
    <!-- 自定义 TabBar -->
    <CustomTabBar :current="0" />
  </view>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
// Pattern 组件
import SearchBar from '@/components/SearchBar/index.vue'
import OfflineBanner from '@/components/OfflineBanner/index.vue'
// Business 组件
import HomeHeader from '@/components/HomeHeader/index.vue'
import HotGamesScroll from '@/components/HotGamesScroll/index.vue'
import RecommendPlayersSection from '@/components/RecommendPlayersSection/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// Composables
import { useHome } from '@/composables/useHome'

const {
  isOfflineMode,
  isLoggedIn,
  userInfo,
  hotGames,
  gamesLoading,
  recommendPlayers,
  playersLoading,
  noMorePlayers,
  loadMorePlayers,
  refreshAll,
  goToLogin,
  goToProfile,
  goToPlayerList,
  goToGameList,
  goToGamePlayers,
  goToPlayerDetail,
  init,
} = useHome()

onMounted(init)
</script>

<style lang="scss" scoped>
.home-page {
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

.status-bar {
  height: var(--status-bar-height, 44px);
}

.main-scroll {
  flex: 1;
  overflow-y: auto;
}
</style>
