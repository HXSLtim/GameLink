<template>
  <BasePageLayout
    class="home-page"
    :scroll="false"
    padding="0"
    title="GameLink"
    :show-back="false"
    :show-tab-bar="true"
    :tab-bar-current="0"
  >
    <template #nav>
      <NavBar title="GameLink" :show-back="false" />
    </template>

    <template #search>
      <!-- 搜索栏 -->
      <SearchBar
        v-model="homeSearchKeyword"
        placeholder="搜索陪玩师、游戏..."
        :show-filter="true"
        @search="handleHomeSearch"
        @filter="handleHomeFilter"
      />
    </template>

    <!-- 移动端：Banner 置顶在 header 区域 -->
    <template v-if="!isPC" #banner>
      <view class="home-banner-stack">
        <view v-if="banners.length > 0" class="home-banner">
          <swiper
            class="banner-swiper"
            :autoplay="banners.length > 1"
            :indicator-dots="false"
            circular
            interval="3500"
            duration="450"
            @change="onBannerChange"
          >
            <swiper-item v-for="banner in banners" :key="banner.id">
              <view class="banner-card" @tap="handleBannerClick(banner)">
                <image class="banner-image" :src="banner.image" mode="aspectFill" />
              </view>
            </swiper-item>
          </swiper>
          <!-- 自定义指示器 -->
          <view v-if="banners.length > 1" class="banner-indicators">
            <view
              v-for="(_, idx) in banners"
              :key="idx"
              class="banner-dot"
              :class="{ active: idx === currentBannerIndex }"
            />
          </view>
        </view>
        <OfflineBanner
          :visible="isOfflineMode"
          :message="hasAnyData ? '网络不可用，显示缓存数据' : '网络不可用，请检查网络后重试'"
          :type="hasAnyData ? 'warning' : 'error'"
          @action="refreshAll"
        />
      </view>
    </template>

    <!-- 主内容区：小程序用 scroll-view 原生滚动，PC 用 overflow 滚动 -->
    <component :is="isPC ? 'view' : 'scroll-view'" class="main-scroll" :class="{ 'pc-scroll': isPC }" :scroll-y="!isPC">
      <view class="content-wrapper fade-in">

        <!-- PC 端：Hero Banner（在滚动内容区内，不置顶） -->
        <view v-if="isPC && banners.length > 0" class="hero-banner">
          <swiper
            class="hero-swiper"
            :autoplay="banners.length > 1"
            :indicator-dots="false"
            circular
            interval="4000"
            duration="500"
            @change="onBannerChange"
          >
            <swiper-item v-for="banner in banners" :key="banner.id">
              <view class="hero-card" @tap="handleBannerClick(banner)">
                <view class="hero-text">
                  <text class="hero-title">{{ banner.title || 'GameLink' }}</text>
                  <text class="hero-desc">{{ banner.description || '' }}</text>
                  <GlButton
                    type="primary"
                    size="small"
                    round
                    @click.stop="handleBannerClick(banner)"
                  >
                    {{ banner.actionText || '了解更多' }}
                  </GlButton>
                </view>
                <view class="hero-image-wrap">
                  <image class="hero-image" :src="banner.image" mode="aspectFill" />
                </view>
              </view>
            </swiper-item>
          </swiper>
          <!-- PC 自定义指示器 -->
          <view v-if="banners.length > 1" class="hero-indicators">
            <view
              v-for="(_, idx) in banners"
              :key="idx"
              class="hero-dot"
              :class="{ active: idx === currentBannerIndex }"
            />
          </view>
        </view>

        <!-- PC 端离线提示 -->
        <OfflineBanner
          v-if="isPC"
          :visible="isOfflineMode"
          :message="hasAnyData ? '网络不可用，显示缓存数据' : '网络不可用，请检查网络后重试'"
          :type="hasAnyData ? 'warning' : 'error'"
          @action="refreshAll"
        />
        <!-- 热门游戏 -->
        <HotGamesScroll
          :games="hotGames"
          :loading="gamesLoading"
          @more="goToGameList"
          @select="(game) => goToGamePlayers(game.id)"
        />

        <!-- 推荐陪玩师 -->
        <view>
          <RecommendPlayersSection
            :players="recommendPlayers"
            :loading="playersLoading"
            @more="goToPlayerList"
            @select="(player) => goToPlayerDetail(player.id)"
            @refresh="refreshAll"
          />
        </view>
      </view>
      
      <!-- 底部垫高，防止被 TabBar 遮挡 -->
      <view class="safe-area-bottom" style="height: 120rpx;"></view>
    </component>
  </BasePageLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useDevice } from '@/composables/useDevice'
import { previewImage } from '@/utils'
// Pattern 组件
import SearchBar from '@/components/SearchBar/index.vue'
import OfflineBanner from '@/components/OfflineBanner/index.vue'
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
// Business 组件
import NavBar from '@/components/NavBar/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
import HotGamesScroll from '@/components/HotGamesScroll/index.vue'
import RecommendPlayersSection from '@/components/RecommendPlayersSection/index.vue'
// Composables
import { useHome } from '@/composables/useHome'
import type { HomeBannerItem } from '@/types/home'

const {
  isOfflineMode,
  banners,
  hotGames,
  gamesLoading,
  recommendPlayers,
  playersLoading,
  refreshAll,
  goToPlayerList,
  goToGameList,
  goToGamePlayers,
  goToPlayerDetail,
  init,
} = useHome()
const { isPC } = useDevice()
const homeSearchKeyword = ref('')
const currentBannerIndex = ref(0)

// 是否有任何数据（缓存或实时）可供展示
const hasAnyData = computed(() => hotGames.value.length > 0 || recommendPlayers.value.length > 0)

const onBannerChange = (e: any) => {
  currentBannerIndex.value = e.detail?.current ?? 0
}

const handleHomeSearch = (keyword: string) => {
  const trimmed = keyword.trim()
  const url = trimmed
    ? `/pages/player/list/index?keyword=${encodeURIComponent(trimmed)}`
    : '/pages/player/list/index'
  uni.navigateTo({ url })
}

const handleHomeFilter = () => {
  const trimmed = homeSearchKeyword.value.trim()
  const url = trimmed
    ? `/pages/player/list/index?keyword=${encodeURIComponent(trimmed)}`
    : '/pages/player/list/index'
  uni.navigateTo({ url })
}

const handleBannerClick = (banner: HomeBannerItem) => {
  if (banner.type === 'preview') {
    const urls = banner.previewImages?.length ? banner.previewImages : [banner.image]
    previewImage(urls, banner.image)
    return
  }
  if (banner.type === 'link' && banner.link) {
    uni.navigateTo({ url: banner.link })
  }
}

onMounted(init)
</script>

<style lang="scss" scoped>
.home-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  box-sizing: border-box;
}

// ============================================
// 移动端 Swiper Banner
// ============================================
.home-banner-stack {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.home-banner {
  padding: var(--spacing-md) var(--spacing-lg) var(--spacing-sm); // 优化: 左右间距加大，上下微调
  position: relative;
}

.banner-swiper {
  height: 360rpx; // 优化: 增加 40rpx，视觉更突出
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.banner-card {
  height: 100%;
  border-radius: var(--radius-lg);
  overflow: hidden;
  background: var(--color-bg-secondary);
  position: relative;
  cursor: pointer;
  @include press-effect;
}

.banner-image {
  width: 100%;
  height: 100%;
  display: block;
}

// 移动端自定义指示器
.banner-indicators {
  display: flex;
  justify-content: center;
  gap: 8rpx;
  margin-top: var(--spacing-xs);
}

.banner-dot {
  width: 12rpx;
  height: 12rpx;
  border-radius: var(--radius-full);
  background: var(--color-text-placeholder);
  opacity: 0.3;
  transition: all 0.3s ease;

  &.active {
    width: 32rpx;
    background: var(--color-primary);
    opacity: 1;
  }
}

// ============================================
// PC 端 Hero Banner（左文字 + 右图片）
// ============================================
.hero-banner {
  padding: var(--spacing-xl) var(--spacing-lg) var(--spacing-md); // 优化: 增加上下间距，更宽松
  position: relative;
}

.hero-swiper {
  height: 220px; // 优化: 增加 20rpx，视觉更突出
  border-radius: var(--radius-lg);
  overflow: hidden;

  @include desktop-lg {
    height: 280px; // 优化: 增加 40rpx
  }
}

.hero-card {
  display: flex;
  align-items: stretch;
  height: 100%;
  border-radius: var(--radius-lg);
  overflow: hidden;
  background: linear-gradient(135deg, var(--color-bg-card) 0%, var(--color-bg-secondary) 100%);
  cursor: pointer;
  transition: box-shadow 0.35s ease, transform 0.35s ease;

  &:hover {
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
    transform: translateY(-2px);

    .hero-image {
      transform: scale(1.04);
    }
  }
}

.hero-text {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: flex-start;
  gap: 14px;
  padding: 32px 40px;
  min-width: 0;
}

.hero-title {
  font-size: 22px;
  font-weight: 700;
  color: var(--color-text);
  line-height: 1.35;
  letter-spacing: 0.3px;
}

.hero-desc {
  font-size: 14px;
  color: var(--color-text-secondary);
  line-height: 1.6;
  max-width: 320px;
}

.hero-image-wrap {
  width: 40%;
  max-width: 400px;
  flex-shrink: 0;
  overflow: hidden;
  position: relative;
}

.hero-image {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
  transition: transform 0.5s ease;
}

// PC 自定义指示器
.hero-indicators {
  display: flex;
  justify-content: center;
  gap: 6px;
  margin-top: 14px;
}

.hero-dot {
  width: 6px;
  height: 6px;
  border-radius: var(--radius-full);
  background: var(--color-text-placeholder);
  opacity: 0.25;
  cursor: pointer;
  transition: all 0.3s ease;

  &:hover {
    opacity: 0.5;
  }

  &.active {
    width: 24px;
    background: var(--color-primary);
    opacity: 1;
  }
}

.main-scroll {
  flex: 1;
  min-height: 0;
  height: 0; // 配合 flex: 1 确保 scroll-view 获得正确高度
}

.pc-scroll {
  overflow-y: auto; // PC 端使用原生滚动
}

.content-wrapper {
  padding-bottom: env(safe-area-inset-bottom);

  // PC 端 section 间距加大 - 优化
  @include desktop {
    padding: 0 var(--spacing-md);
    > * + * {
      margin-top: var(--spacing-xl); // 优化: section 间距增加到 52rpx
    }
  }
}
</style>
