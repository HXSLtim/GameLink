<template>
  <BasePageLayout
    class="player-detail-page"
    padding="0"
    title="陪玩师详情"
    :show-back="true"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
  >
    <template #nav>
      <!-- 顶部导航 -->
      <NavBar title="陪玩师详情" :show-share="true" @back="goBack" @share="handleShare" />
    </template>

    <!-- 页面状态 -->
    <PageState
      :state="pageState"
      :error-message="errorMessage"
      empty-title="陪玩师不存在"
      empty-desc="该陪玩师可能已下架"
      @retry="handleRetry"
    >
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
    </PageState>

    <template #footer>
      <!-- 底部操作栏 -->
      <PlayerActionBar
        :is-favorite="isFavorite"
        :is-online="player.isOnline"
        @favorite="toggleFavorite"
        @chat="goToChat"
        @order="goToOrder"
      />
    </template>
  </BasePageLayout>
</template>

<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
import PageState from '@/components/PageState/index.vue'
// Business 组件
import PlayerDetailHeader from '@/components/PlayerDetailHeader/index.vue'
import PlayerGamesSection from '@/components/PlayerGamesSection/index.vue'
import PlayerServicesSection from '@/components/PlayerServicesSection/index.vue'
import PlayerReviewsSection from '@/components/PlayerReviewsSection/index.vue'
import PlayerActionBar from '@/components/PlayerActionBar/index.vue'
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
.bottom-placeholder {
  height: 160rpx;
}
</style>
