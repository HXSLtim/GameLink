<template>
  <view class="player-detail-page page-container">
    <!-- 顶部导航 -->
    <NavBar title="陪玩师详情" :show-share="true" @back="goBack" @share="handleShare" />

    <!-- 页面状态 -->
    <PageState
      :state="pageState"
      :error-message="errorMessage"
      empty-title="陪玩师不存在"
      empty-desc="该陪玩师可能已下架"
      @retry="handleRetry"
    >
      <!-- 内容区域 -->
      <scroll-view class="content-scroll" scroll-y>
        <!-- 头部信息卡片 -->
        <PlayerDetailHeader :player="player" />

        <!-- 擅长游戏 -->
        <PlayerGamesSection :games="player.games" />

        <!-- 服务项目 -->
        <PlayerServicesSection
          :services="player.services"
          :selected-id="selectedService?.id"
          @select="selectService"
        />

        <!-- 用户评价 -->
        <PlayerReviewsSection
          :rating="player.rating"
          :reviews="displayReviews"
          @more="goToReviews"
        />

        <!-- 底部占位 -->
        <view class="bottom-placeholder"></view>
      </scroll-view>
    </PageState>

    <!-- 底部操作栏 -->
    <PlayerActionBar
      :is-favorite="isFavorite"
      :is-online="player.isOnline"
      @favorite="toggleFavorite"
      @chat="goToChat"
      @order="goToOrder"
    />

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import PageState from '@/components/PageState.vue'
// Business 组件
import PlayerDetailHeader from '@/components/PlayerDetailHeader/index.vue'
import PlayerGamesSection from '@/components/PlayerGamesSection/index.vue'
import PlayerServicesSection from '@/components/PlayerServicesSection/index.vue'
import PlayerReviewsSection from '@/components/PlayerReviewsSection/index.vue'
import PlayerActionBar from '@/components/PlayerActionBar/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// Composables
import { usePlayerDetail } from '@/composables/usePlayerDetail'

const {
  pageState,
  errorMessage,
  player,
  isFavorite,
  selectedService,
  displayReviews,
  loadPlayerDetail,
  handleRetry,
  toggleFavorite,
  selectService,
  goBack,
  handleShare,
  goToReviews,
  goToChat,
  goToOrder,
} = usePlayerDetail()

onLoad((options) => {
  const id = Number(options?.id)
  if (id) {
    loadPlayerDetail(id)
  }
})
</script>

<style lang="scss" scoped>
.player-detail-page {
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

.content-scroll {
  flex: 1;
  overflow-y: auto;
}

.bottom-placeholder {
  height: 180rpx;
}
</style>
