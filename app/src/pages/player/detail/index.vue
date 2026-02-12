<template>
  <BasePageLayout
    class="player-detail-page"
    padding="0"
    :title="isPC ? '陪玩师详情' : ''"
    :show-back="isPC"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
  >
    <!-- Mobile NavBar: 透明渐变 -->
    <template #nav v-if="!isPC">
      <NavBar
        :title="showNavTitle ? player.nickname : ''"
        :show-share="true"
        :transparent="!showNavBg"
        :fixed="true"
        @back="goBack"
        @share="handleShare"
        :style="{ background: showNavBg ? 'rgba(var(--bg-rgb), 0.95)' : 'transparent' }"
      />
    </template>

    <!-- 页面状态 -->
    <PageState
      :state="pageState"
      :error-message="errorMessage"
      empty-title="陪玩师不存在"
      empty-desc="该陪玩师可能已下架"
      @retry="handleRetry"
    >
      <!-- PC Layout (Dual Column) -->
      <view v-if="isPC" class="pc-detail-layout">
        <!-- Sidebar: Sticky Profile Card -->
        <view class="pc-detail-sidebar">
          <PlayerDetailHeader :player="player" mode="card" />

          <view class="pc-action-card">
            <PlayerActionBar
              :is-favorite="isFavorite"
              :is-online="player.isOnline"
              :price="player.price"
              @favorite="toggleFavorite"
              @chat="goToChat"
              @order="goToOrder"
            />
          </view>
        </view>

        <!-- Main Content: Scrollable -->
        <view class="pc-detail-main">
          <view class="content-card">
            <PlayerServicesSection
              :services="player.services"
              :selected-id="selectedService?.id"
              @select="selectService"
            />
          </view>

          <view class="content-card">
            <PlayerGamesSection :games="player.games" />
          </view>

          <view class="content-card">
            <PlayerReviewsSection
              :rating="player.rating"
              :reviews="displayReviews"
              @more="goToReviews"
            />
          </view>
        </view>
      </view>

      <!-- Mobile Layout (Immersive Stream) -->
      <view v-else class="mobile-detail-layout">
        <PlayerDetailHeader :player="player" mode="hero" />

        <view class="mobile-content-stack">
          <view class="mobile-card">
            <PlayerServicesSection
              :services="player.services"
              :selected-id="selectedService?.id"
              @select="selectService"
            />
          </view>

          <view class="mobile-card">
            <PlayerGamesSection :games="player.games" />
          </view>

          <view class="mobile-card">
            <PlayerReviewsSection
              :rating="player.rating"
              :reviews="displayReviews"
              @more="goToReviews"
            />
          </view>
        </view>

        <!-- 底部占位，防止被 Action Bar 遮挡 -->
        <view class="bottom-placeholder"></view>
      </view>
    </PageState>

    <!-- Mobile Floating Action Bar -->
    <template #footer>
      <view v-if="!isPC" class="mobile-action-bar-container">
        <PlayerActionBar
          :is-favorite="isFavorite"
          :is-online="player.isOnline"
          :price="player.price"
          @favorite="toggleFavorite"
          @chat="goToChat"
          @order="goToOrder"
        />
      </view>
    </template>
  </BasePageLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad, onPageScroll } from '@dcloudio/uni-app'
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

// Scroll handling for Mobile NavBar
const showNavBg = ref(false)
const showNavTitle = ref(false)

onPageScroll((e) => {
  if (isPC.value) return
  const scrollTop = e.scrollTop
  showNavBg.value = scrollTop > 50
  showNavTitle.value = scrollTop > 200
})

onLoad((options) => {
  const id = Number(options?.id)
  if (id) {
    loadPlayerDetail(id)
  }
})
</script>

<style lang="scss" scoped>
// ============================================
// PC Layout
// ============================================
.pc-detail-layout {
  display: flex;
  gap: var(--spacing-lg);
  align-items: flex-start;
  padding: var(--spacing-lg) 0;
  max-width: 1200px;
  margin: 0 auto;
}

.pc-detail-sidebar {
  width: 360px;
  flex-shrink: 0;
  position: sticky;
  top: 80px; // Offset for NavBar
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.pc-action-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  border: 1rpx solid var(--color-border);
}

.pc-detail-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.content-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-lg);
  padding: var(--spacing-lg);
  box-shadow: var(--shadow-sm);
  border: 1rpx solid var(--color-border);
}

// ============================================
// Mobile Layout
// ============================================
.mobile-detail-layout {
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
}

.mobile-content-stack {
  position: relative;
  z-index: 2;
  margin-top: -40rpx; // Slight overlap with hero
  padding: 0 var(--spacing-md);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.mobile-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-lg);
  padding: var(--spacing-md);
  box-shadow: var(--shadow-md);
}

.mobile-action-bar-container {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 100;
  padding: var(--spacing-sm);
  pointer-events: none; // Allow clicks through container

  :deep(.action-bar) {
    pointer-events: auto;
    box-shadow: var(--shadow-lg);
  }
}

.bottom-placeholder {
  height: 180rpx;
}
</style>
