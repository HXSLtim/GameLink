<template>
  <view
    class="player-card"
    :class="[`player-card--${variant}`, { 'player-card--clickable': clickable }]"
    @tap="handleClick"
  >
    <view
      v-if="showSelect"
      class="player-select"
      :class="{ checked: selected }"
      @tap.stop="emit('toggle-select')"
    >
      <uv-icon v-if="selected" name="checkbox-mark" size="16" color="#fff"></uv-icon>
    </view>

    <template v-if="isGrid">
      <!-- Grid：头像 + 在线状态点 -->
      <view class="grid-top">
        <GlAvatar
          :src="player.avatar"
          :text="player.nickname"
          :size="avatarSize"
          :status="avatarStatus"
        />
      </view>
      
      <view class="grid-body">
        <!-- 名字 + 认证 -->
        <text class="player-name">{{ player.nickname }}</text>

        <!-- 评分（突出显示） + 接单数 -->
        <view v-if="showMeta" class="grid-meta">
          <view v-if="showRating" class="grid-rating">
            <uv-icon name="star-fill" size="12" color="var(--color-gold)" />
            <text class="grid-rating-value">{{ displayRating }}</text>
          </view>
          <text v-if="showOrderCount && (player.orderCount || 0) > 0" class="grid-stat">{{ player.orderCount }}单</text>
          <GlTag v-if="showVerified && player.isVerified" size="mini" plain>认证</GlTag>
        </view>

        <!-- 简介（单行截断） -->
        <text v-if="player.bio" class="player-bio">{{ player.bio }}</text>
      </view>

      <!-- Grid 底部：游戏 + 价格 -->
      <view class="grid-footer">
        <view class="grid-footer-left">
          <PlayerGames v-if="showGames && displayGames.length" :games="displayGames" :more-count="moreGamesCount" />
          <GlTag v-if="showRank && player.rank" size="mini" type="warning">{{ player.rank }}</GlTag>
        </view>
        <view v-if="showPrice && hasPrice" class="player-price player-price--inline">
          <PriceTag
            :amount="priceValue ?? 0"
            amount-unit="yuan"
            :unit="resolvedPriceUnit"
            size="medium"
            :show-decimal="false"
          />
        </view>
      </view>
    </template>

    <template v-else>
      <GlAvatar
        :src="player.avatar"
        :text="player.nickname"
        :size="avatarSize"
        :status="avatarStatus"
      />

      <view class="player-info">
        <view class="player-header">
          <text class="player-name">{{ player.nickname }}</text>
          <PlayerTags
            v-if="showRank || showVerified || showStatusTag"
            :show-rank="showRank"
            :show-verified="showVerified"
            :show-status-tag="showStatusTag"
            :rank="player.rank"
            :is-verified="player.isVerified"
            :is-online="player.isOnline"
          />
        </view>

        <PlayerGames v-if="showGames && displayGames.length" :games="displayGames" :more-count="moreGamesCount" />

        <PlayerMeta
          v-if="showMeta"
          :show-rating="showRating"
          :rating="displayRating"
          :show-order-count="showOrderCount"
          :order-count="player.orderCount || 0"
          :show-main-game="showMainGame"
          :main-game="player.mainGame"
        />
      </view>

      <view v-if="showPrice && hasPrice" class="player-price">
        <PriceTag
          :amount="priceValue ?? 0"
          amount-unit="yuan"
          :unit="resolvedPriceUnit"
          size="medium"
          :show-decimal="false"
        />
      </view>

      <uv-icon v-if="showArrow" name="arrow-right" size="20" color="var(--color-text-secondary)"></uv-icon>
    </template>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlAvatar from '@/components/gl/Avatar/index.vue'
import GlTag from '@/components/gl/Tag/index.vue'
import PlayerTags from './PlayerTags.vue'
import PlayerGames from './PlayerGames.vue'
import PlayerMeta from './PlayerMeta.vue'
import PriceTag from '@/components/PriceTag/index.vue'
import type { PlayerCardData, PlayerGameTag } from '@/types/player'

type PlayerCardVariant = 'default' | 'compact' | 'recommend' | 'favorite' | 'grid'

type PlayerCardVariantConfig = {
  showPrice: boolean
  showGames: boolean
  showRank: boolean
  showMainGame: boolean
  showVerified: boolean
  showOnlineTag: boolean
  showOfflineTag: boolean
  showOrderCount: boolean
  showRating: boolean
}

const VARIANT_DEFAULTS: Record<PlayerCardVariant, PlayerCardVariantConfig> = {
  grid: {
    showPrice: true,
    showGames: true,
    showRank: true,
    showMainGame: true,
    showVerified: true,
    showOnlineTag: true,
    showOfflineTag: true,
    showOrderCount: true,
    showRating: true,
  },
  compact: {
    showPrice: false,
    showGames: false,
    showRank: false,
    showMainGame: false,
    showVerified: false,
    showOnlineTag: true,
    showOfflineTag: true,
    showOrderCount: false,
    showRating: true,
  },
  recommend: {
    showPrice: true,
    showGames: false,
    showRank: true,
    showMainGame: true,
    showVerified: false,
    showOnlineTag: false,
    showOfflineTag: false,
    showOrderCount: true,
    showRating: true,
  },
  favorite: {
    showPrice: true,
    showGames: true,
    showRank: false,
    showMainGame: false,
    showVerified: false,
    showOnlineTag: true,
    showOfflineTag: false,
    showOrderCount: true,
    showRating: true,
  },
  default: {
    showPrice: true,
    showGames: true,
    showRank: false,
    showMainGame: false,
    showVerified: true,
    showOnlineTag: false,
    showOfflineTag: false,
    showOrderCount: true,
    showRating: true,
  },
}

interface Props {
  player: PlayerCardData
  variant?: PlayerCardVariant
  maxGames?: number
  clickable?: boolean
  showArrow?: boolean
  showSelect?: boolean
  selected?: boolean
  showPrice?: boolean
  showGames?: boolean
  showRank?: boolean
  showMainGame?: boolean
  showVerified?: boolean
  showOnlineTag?: boolean
  showOfflineTag?: boolean
  showOrderCount?: boolean
  showRating?: boolean
  priceUnit?: string
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  maxGames: 2,
  clickable: false,
  showArrow: false,
  // 显式设为 undefined，阻止 Vue 3 boolean casting (absent → false)
  // 这样 resolveFlag 的 ?? 才能正确回退到 variant 默认值
  showPrice: undefined,
  showGames: undefined,
  showRank: undefined,
  showMainGame: undefined,
  showVerified: undefined,
  showOnlineTag: undefined,
  showOfflineTag: undefined,
  showOrderCount: undefined,
  showRating: undefined,
})

const emit = defineEmits<{
  click: [player: PlayerCardData]
  'toggle-select': []
}>()

const variantDefaults = computed(() => VARIANT_DEFAULTS[props.variant])

const resolveFlag = (value: boolean | undefined, fallback: boolean) => {
  return value ?? fallback
}

const showPrice = computed(() => resolveFlag(props.showPrice, variantDefaults.value.showPrice))
const showGames = computed(() => resolveFlag(props.showGames, variantDefaults.value.showGames))
const showRank = computed(() => resolveFlag(props.showRank, variantDefaults.value.showRank))
const showMainGame = computed(() => resolveFlag(props.showMainGame, variantDefaults.value.showMainGame))
const showVerified = computed(() => resolveFlag(props.showVerified, variantDefaults.value.showVerified))
const showOnlineTag = computed(() => resolveFlag(props.showOnlineTag, variantDefaults.value.showOnlineTag))
const showOfflineTag = computed(() => resolveFlag(props.showOfflineTag, variantDefaults.value.showOfflineTag))
const showOrderCount = computed(() => resolveFlag(props.showOrderCount, variantDefaults.value.showOrderCount))
const showRating = computed(() => resolveFlag(props.showRating, variantDefaults.value.showRating))
const isGrid = computed(() => props.variant === 'grid')

const showStatusTag = computed(() => {
  // 当 isOnline 有明确值时才显示状态标签
  if (props.player.isOnline === undefined) return false
  if (props.player.isOnline) return showOnlineTag.value
  return showOfflineTag.value
})

const normalizedGames = computed<PlayerGameTag[]>(() => {
  return (props.player.games || []).map((game, index) => {
    if (typeof game === 'string') {
      return { key: `${game}-${index}`, name: game }
    }
    return { key: String(game.id ?? `${game.name}-${index}`), name: game.name }
  })
})

const displayGames = computed(() => normalizedGames.value.slice(0, props.maxGames))
const moreGamesCount = computed(() => Math.max(0, normalizedGames.value.length - props.maxGames))

const avatarSize = computed(() => {
  if (props.variant === 'grid') return 'xlarge'
  if (props.variant === 'compact') return 'medium'
  return 'large'
})

const avatarStatus = computed(() => {
  if (props.player.status) return props.player.status
  if (props.player.isOnline === undefined) return undefined
  return props.player.isOnline ? 'online' : 'offline'
})

const displayRating = computed(() => Number(props.player.rating ?? 5).toFixed(1))
const showMeta = computed(() => showRating.value || showOrderCount.value || (showMainGame.value && Boolean(props.player.mainGame)))

const priceValue = computed(() => {
  if (props.player.hourlyRate !== undefined && props.player.hourlyRate !== null) {
    return Math.round(props.player.hourlyRate / 100)
  }
  if (props.player.minPrice !== undefined && props.player.minPrice !== null) {
    return props.player.minPrice
  }
  return undefined
})

const hasPrice = computed(() => priceValue.value !== undefined && priceValue.value !== null)
const resolvedPriceUnit = computed(() => {
  if (props.priceUnit !== undefined) return props.priceUnit.replace(/^\//, '')
  if (props.player.hourlyRate !== undefined && props.player.hourlyRate !== null) return '小时'
  if (props.variant === 'recommend') return '小时'
  return '局'
})

const handleClick = () => {
  if (!props.clickable) return
  emit('click', props.player)
}
</script>

<style lang="scss" scoped>
.player-card {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-sm) var(--spacing-md);
  background: var(--color-bg-card);
  border-radius: var(--radius-lg); // 优化：加大圆角
  border: 1rpx solid var(--color-border);
  width: 100%;
  position: relative;
  overflow: hidden; // 确保子元素不溢出圆角
}

.player-card--compact {
  padding: var(--spacing-sm);
}

.player-card--grid {
  flex-direction: column;
  align-items: center;
  gap: 0;
  padding: var(--spacing-lg) var(--spacing-md);
  text-align: center;
  transition: transform 0.25s ease, box-shadow 0.25s ease, border-color 0.25s ease;

  @include desktop {
    padding: 24px 20px;
  }

  &:hover {
    transform: translateY(-4rpx);
    box-shadow: var(--shadow-lg);
    border-color: var(--color-primary);
  }
}

.grid-top {
  margin-bottom: var(--spacing-md);
}

.grid-body {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-xs);
  margin-bottom: var(--spacing-md);
}

// Grid 内联评分 + 接单数 + 认证
.grid-meta {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
  flex-wrap: wrap;
}

.grid-rating {
  display: inline-flex;
  align-items: center;
  gap: 4rpx;
  line-height: 1;
}

.grid-rating-value {
  font-size: var(--font-xs);
  font-weight: 600;
  color: var(--color-text);
  line-height: 1;
}

.grid-stat {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}

.player-bio {
  font-size: var(--font-xs);
  color: var(--color-text-placeholder);
  @include text-ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-align: center;
  max-width: 100%;
  margin-top: 2rpx;
}

.grid-footer {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: var(--spacing-md);
  border-top: 1rpx solid var(--color-border);
  gap: var(--spacing-sm);
}

.grid-footer-left {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  min-width: 0;
  flex: 1;
  
  :deep(.player-games) {
    margin-bottom: 0;
  }
}

.player-card--clickable {
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); // 优化动画曲线
  @include press-effect;

  &:hover {
    box-shadow: var(--shadow-md);
    border-color: var(--color-primary-light);
    transform: translateY(-2rpx);
  }

  &:active {
    background: var(--color-bg-secondary);
    transform: scale(0.98);
  }
}

.player-select {
  width: 36rpx;
  height: 36rpx;
  border-radius: var(--radius-full);
  border: 1rpx solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: var(--color-bg-card);

  &.checked {
    background: var(--color-primary);
    border-color: var(--color-primary);
  }
}

.player-info {
  flex: 1;
  min-width: 0;
}

.player-info--grid {
  width: 100%;
}

.player-title-row {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  flex-wrap: wrap;
}

.player-title-row .player-name {
  flex-shrink: 0;
}

.player-title-row :deep(.player-tags) {
  margin-bottom: 0;
}

.player-header {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: var(--spacing-xs);
  margin-bottom: var(--spacing-xs);
}

.player-name {
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
  @include text-ellipsis;

  @include desktop {
    font-size: 15px;
  }
}

.player-price {
  flex-shrink: 0;
  text-align: right;
}

.player-price--inline {
  display: flex;
  align-items: baseline;
  gap: 4rpx;
  text-align: right;
}
.player-price :deep(.price-tag) {
  color: var(--color-text);
}

.player-price :deep(.price-tag .amount) {
  font-weight: 700;
  font-size: var(--font-lg);
}

.player-price :deep(.price-tag .unit) {
  font-size: var(--font-xs);
}

.player-card--recommend {
  :deep(.price-tag .amount) {
    font-size: var(--font-md);
  }
}
</style>
