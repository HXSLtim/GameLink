<template>
  <view class="player-header" :class="[`player-header--${mode}`]">
    <!-- Mobile Hero Mode -->
    <template v-if="mode === 'hero'">
      <view class="hero-cover">
        <image v-if="player.coverImage" :src="player.coverImage" mode="aspectFill" class="hero-image" />
        <view v-else class="hero-placeholder"></view>
        <view class="hero-overlay"></view>
      </view>

      <view class="hero-content">
        <view class="hero-user">
          <GlAvatar
            :src="player.avatar"
            :text="player.nickname"
            :size="140"
            :status="player.isOnline ? 'online' : undefined"
            bordered
            class="hero-avatar"
          />
          <view class="hero-info">
            <view class="hero-name-row">
              <text class="hero-nickname">{{ player.nickname }}</text>
              <GlTag v-if="player.isVerified" type="success" size="mini" class="hero-tag">已认证</GlTag>
            </view>
            <view class="hero-meta-row">
              <view v-if="player.gender" class="gender-badge" :class="player.gender">
                <text>{{ player.gender === 'male' ? '♂' : '♀' }}</text>
              </view>
              <text class="hero-id">ID: {{ player.id }}</text>
            </view>
          </view>
        </view>

        <text v-if="player.signature" class="hero-signature">{{ player.signature }}</text>

        <HeaderStatsRow
          class="hero-stats"
          :items="stats"
          size="md"
          item-padding="0"
          theme="dark"
        />
      </view>
    </template>

    <!-- PC Card Mode -->
    <template v-else>
      <view class="header-card">
        <!-- 封面 -->
        <view class="player-cover">
          <image v-if="player.coverImage" :src="player.coverImage" mode="aspectFill" class="cover-image" />
          <view v-else class="cover-placeholder"></view>
        </view>

        <!-- 基本信息 -->
        <view class="player-basic">
          <view class="avatar-wrap">
            <GlAvatar
              :src="player.avatar"
              :text="player.nickname"
              :size="100"
              :status="player.isOnline ? 'online' : undefined"
              bordered
            />
          </view>

          <view class="basic-info">
            <view class="name-row">
              <text class="nickname">{{ player.nickname }}</text>
              <GlTag v-if="player.isVerified" type="success" size="mini">已认证</GlTag>
              <view v-if="player.gender" class="gender-badge" :class="player.gender">
                <text>{{ player.gender === 'male' ? '♂' : '♀' }}</text>
              </view>
            </view>

            <view class="status-row">
              <GlTag :type="player.isOnline ? 'success' : 'default'" size="mini" plain>
                {{ player.isOnline ? '在线' : '离线' }}
              </GlTag>
              <text class="player-id">ID: {{ player.id }}</text>
            </view>
          </view>
        </view>

        <view class="card-signature" v-if="player.signature">
          <text>{{ player.signature }}</text>
        </view>

        <!-- 统计数据 -->
        <HeaderStatsRow
          class="stats-row"
          :items="stats"
          size="md"
          item-padding="0"
        />
      </view>
    </template>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlAvatar from '@/components/gl/Avatar/index.vue'
import GlTag from '@/components/gl/Tag/index.vue'
import HeaderStatsRow from '@/components/HeaderStatsRow/index.vue'
import type { PlayerHeaderData } from '@/types/player'
import type { HeaderStatItem } from '@/types/ui'

interface Props {
  player: PlayerHeaderData
  mode?: 'hero' | 'card'
}

const props = withDefaults(defineProps<Props>(), {
  mode: 'card'
})

const formatJoinDate = computed(() => {
  if (!props.player.createdAt) return '-'
  const date = new Date(props.player.createdAt)
  const now = new Date()
  const months = (now.getFullYear() - date.getFullYear()) * 12 + (now.getMonth() - date.getMonth())
  if (months < 1) return '刚入驻'
  if (months < 12) return `${months}个月`
  return `${Math.floor(months / 12)}年`
})

const stats = computed<HeaderStatItem[]>(() => [
  { label: '评分', value: props.player.rating?.toFixed(1) || '5.0' },
  { label: '接单数', value: props.player.orderCount || 0 },
  { label: '收藏数', value: props.player.favoriteCount || 0 },
  { label: '入驻时间', value: formatJoinDate.value },
])
</script>

<style lang="scss" scoped>
.player-header {
  width: 100%;
}

// ============================================
// Mobile Hero Mode
// ============================================
.player-header--hero {
  position: relative;
  min-height: 560rpx;
  background: var(--color-bg);
  overflow: hidden;
}

.hero-cover {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 0;

  .hero-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .hero-placeholder {
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, #2c3e50, #000000);
  }

  .hero-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(
      to bottom,
      rgba(0,0,0,0.1) 0%,
      rgba(0,0,0,0.3) 50%,
      rgba(15, 23, 42, 0.95) 100%
    );
    backdrop-filter: blur(2px);
  }
}

.hero-content {
  position: relative;
  z-index: 1;
  padding: 0 var(--spacing-lg) var(--spacing-lg);
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  height: 560rpx; // Match min-height
  box-sizing: border-box;
  padding-top: 180rpx; // Space for navbar
}

.hero-user {
  display: flex;
  align-items: flex-end;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.hero-avatar {
  border: 4rpx solid #fff;
  box-shadow: 0 0 32rpx rgba(var(--color-primary-rgb), 0.4);
}

.hero-info {
  flex: 1;
  margin-bottom: var(--spacing-xs);
}

.hero-name-row {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-xs);
  flex-wrap: wrap;
}

.hero-nickname {
  font-size: var(--font-2xl);
  font-weight: 800;
  color: #fff;
  text-shadow: 0 2rpx 4rpx rgba(0,0,0,0.5);
}

.hero-meta-row {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.hero-id {
  font-size: var(--font-xs);
  color: rgba(255, 255, 255, 0.7);
  font-family: monospace;
}

.hero-signature {
  font-size: var(--font-sm);
  color: rgba(255, 255, 255, 0.8);
  margin-bottom: var(--spacing-md);
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.hero-stats {
  :deep(.stat-label) {
    color: rgba(255, 255, 255, 0.6) !important;
  }
  :deep(.stat-value) {
    color: #fff !important;
  }
}

// ============================================
// PC Card Mode
// ============================================
.header-card {
  position: relative;
  background: var(--color-bg-card);
  border-radius: var(--radius-lg);
  overflow: hidden;
  border: 1rpx solid var(--color-border);
}

.player-cover {
  height: 160rpx;
  background: var(--color-bg-secondary);
  position: relative;

  .cover-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.player-basic {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0 var(--spacing-md);
  margin-top: -50rpx;
  position: relative;
  z-index: 1;
  text-align: center;
}

.avatar-wrap {
  margin-bottom: var(--spacing-sm);
}

.basic-info {
  width: 100%;
  margin-bottom: var(--spacing-md);
}

.name-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-xs);
  flex-wrap: wrap;
  margin-bottom: var(--spacing-xs);
}

.nickname {
  font-size: var(--font-lg);
  font-weight: 700;
  color: var(--color-text);
}

.status-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-md);
}

.player-id {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}

.card-signature {
  padding: 0 var(--spacing-lg) var(--spacing-md);
  text-align: center;
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  line-height: 1.5;
}

.stats-row {
  padding: var(--spacing-md);
  border-top: 1rpx solid var(--color-border);
}

.gender-badge {
  width: 32rpx;
  height: 32rpx;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20rpx;

  &.male {
    background: rgba(59, 130, 246, 0.1);
    color: #3B82F6;
  }

  &.female {
    background: rgba(239, 68, 68, 0.1);
    color: #EF4444;
  }
}
</style>
