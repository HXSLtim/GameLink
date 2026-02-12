<template>
  <view class="action-bar" :class="{ 'action-bar--mobile': !isPC, 'action-bar--pc': isPC }">
    <view class="action-left">
      <view class="action-icon-btn" @tap="$emit('favorite')">
        <uv-icon
          :name="isFavorite ? 'heart-fill' : 'heart'"
          size="22"
          :color="isFavorite ? 'var(--color-error)' : 'var(--color-text-secondary)'"
          class="action-icon"
        />
        <text class="action-text" :class="{ 'text-active': isFavorite }">收藏</text>
      </view>
      <view class="action-icon-btn" @tap="$emit('chat')">
        <uv-icon name="chat" size="22" color="var(--color-text-secondary)" class="action-icon" />
        <text class="action-text">聊天</text>
      </view>
    </view>

    <view class="action-right">
      <view class="price-info" v-if="price">
        <text class="price-symbol">¥</text>
        <text class="price-amount">{{ price }}</text>
        <text class="price-unit">/局</text>
      </view>

      <GlButton
        :type="isOnline ? 'primary' : 'default'"
        :disabled="!isOnline"
        size="large"
        round
        block
        class="action-btn"
        @click="$emit('order')"
      >
        {{ isOnline ? '立即下单' : '陪玩师离线' }}
      </GlButton>
    </view>
  </view>
</template>

<script setup lang="ts">
import GlButton from '@/components/gl/Button/index.vue'
import { useDevice } from '@/composables/useDevice'

interface Props {
  isFavorite: boolean
  isOnline: boolean
  price?: number | string
}

defineProps<Props>()

defineEmits<{
  favorite: []
  chat: []
  order: []
}>()

const { isPC } = useDevice()
</script>

<style lang="scss" scoped>
.action-bar {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-sm) var(--spacing-md);
  background: var(--color-bg-card);
  transition: all 0.3s ease;
}

// ============================================
// Mobile: Floating Glassmorphism
// ============================================
.action-bar--mobile {
  margin: 0 var(--spacing-sm) var(--spacing-sm);
  border-radius: var(--radius-full);
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border: 1rpx solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.12);
  padding-bottom: var(--spacing-sm); // Override safe-area if handled by container

  // Dark mode support
  @media (prefers-color-scheme: dark) {
    background: rgba(30, 30, 35, 0.85);
    border-color: rgba(255, 255, 255, 0.1);
  }
}

// ============================================
// PC: Standard Block
// ============================================
.action-bar--pc {
  background: var(--color-bg-card);
  padding: var(--spacing-md);
  border-top: none;
}

.action-left {
  display: flex;
  gap: var(--spacing-sm);
}

.action-icon-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-width: 80rpx;
  gap: 4rpx;
  background: transparent;
  border-radius: var(--radius-md);
  transition: all 0.2s;
  cursor: pointer;
  @include press-effect;

  &:active {
    transform: scale(0.95);
  }
}

.action-icon {
  flex-shrink: 0;
}

.action-text {
  font-size: 20rpx;
  color: var(--color-text-secondary);
  line-height: 1;

  &.text-active {
    color: var(--color-error);
  }
}

.action-right {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: var(--spacing-md);
}

.price-info {
  display: flex;
  align-items: baseline;
  color: var(--color-error);
  margin-right: var(--spacing-xs);

  .price-symbol {
    font-size: var(--font-xs);
  }

  .price-amount {
    font-size: var(--font-xl);
    font-weight: 700;
    line-height: 1;
  }

  .price-unit {
    font-size: var(--font-xs);
    color: var(--color-text-secondary);
    margin-left: 2rpx;
  }
}

.action-btn {
  flex: 1;
  // Make button slightly smaller on mobile to fit nicely
  min-width: 200rpx;
}
</style>