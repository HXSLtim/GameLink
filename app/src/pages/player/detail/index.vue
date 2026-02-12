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
      <!-- PC 端布局 -->
      <view v-if="isPC" class="pc-detail-layout">
        <!-- 左侧主要内容 -->
        <view class="pc-detail-main">
          <PlayerDetailHeader :player="player" />
          <PlayerGamesSection :games="player.games" />
          <PlayerReviewsSection
            :rating="player.rating"
            :reviews="displayReviews"
            @more="goToReviews"
          />
        </view>
        
        <!-- 右侧侧边栏 (服务选择 + 操作) -->
        <view class="pc-detail-sidebar">
          <PlayerServicesSection
            :services="player.services"
            :selected-id="selectedService?.id"
            @select="selectService"
          />
          
          <view class="pc-action-bar-wrapper">
            <PlayerActionBar
              :is-favorite="isFavorite"
              :is-online="player.isOnline"
              @favorite="toggleFavorite"
              @chat="goToChat"
              @order="goToOrder"
            />
          </view>
        </view>
      </view>

      <!-- 移动端布局 (垂直堆叠) -->
      <view v-else class="mobile-detail-layout">
        <PlayerDetailHeader :player="player" />
        <PlayerGamesSection :games="player.games" />
        <PlayerServicesSection
          :services="player.services"
          :selected-id="selectedService?.id"
          @select="selectService"
        />
        <PlayerReviewsSection
          :rating="player.rating"
          :reviews="displayReviews"
          @more="goToReviews"
        />
        
        <!-- 底部占位 -->
        <view class="bottom-placeholder"></view>
      </view>
    </PageState>

    <template #footer>
      <!-- 移动端底部操作栏 -->
      <PlayerActionBar
        v-if="!isPC"
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
import { useDevice } from '@/composables/useDevice'

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

const { isPC } = useDevice()

onLoad((options) => {
  const id = Number(options?.id)
  if (id) {
    loadPlayerDetail(id)
  }
})
</script>

<style lang="scss" scoped>
// 底部操作栏（PC端侧边栏用）
.pc-action-bar-wrapper {
  margin-top: var(--spacing-md);
  padding: var(--spacing-md);
  background: var(--color-bg-card);
  border: 1rpx solid var(--color-border);
  border-radius: var(--radius-lg);
}

.pc-detail-layout {
  display: flex;
  gap: var(--spacing-lg);
  align-items: flex-start;
  padding: var(--spacing-md);
  max-width: 1200px;
  margin: 0 auto;
}

.pc-detail-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.pc-detail-sidebar {
  width: 360px;
  flex-shrink: 0;
  position: sticky;
  top: 20px;
}

.mobile-detail-layout {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm); // Mobile gap
}

.bottom-placeholder {
  height: 160rpx;
}
</style>
