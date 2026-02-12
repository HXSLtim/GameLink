<template>
  <SectionCard>
    <template #header>
      <SectionHeader title="用户评价" @more="$emit('more')" />
    </template>

    <!-- 评分概览 -->
    <view class="reviews-summary">
      <view class="rating-big">
        <text class="rating-value">{{ rating?.toFixed(1) || '5.0' }}</text>
        <RatingStars :rating="rating || 5" size="small" />
      </view>
    </view>

    <!-- 评价列表 -->
    <view class="reviews-list">
      <view v-for="review in reviews" :key="review.id" class="review-card glass-card">
        <view class="review-header">
          <GlAvatar :src="review.userAvatar" :text="review.userName" size="small" />
          <view class="review-user-info">
            <text class="review-user-name">{{ review.userName }}</text>
            <RatingStars :rating="review.rating" size="mini" />
          </view>
          <text class="review-time">{{ formatTime(review.createdAt) }}</text>
        </view>
        <text class="review-content">{{ review.content }}</text>
        <view v-if="review.images?.length" class="review-images">
          <image
            v-for="(img, idx) in review.images.slice(0, 3)"
            :key="idx"
            :src="img"
            mode="aspectFill"
            class="review-image"
            @tap="previewImages(review.images!, idx)"
          />
          <view v-if="review.images.length > 3" class="more-images">
            <text>+{{ review.images.length - 3 }}</text>
          </view>
        </view>
      </view>

      <GlEmpty v-if="!reviews?.length" title="暂无评价" description="下单体验后来评价吧" compact />
    </view>
  </SectionCard>
</template>

<script setup lang="ts">
import SectionCard from '@/components/SectionCard/index.vue'
import GlAvatar from '@/components/gl/Avatar/index.vue'
import GlEmpty from '@/components/gl/Empty/index.vue'
import SectionHeader from '@/components/SectionHeader/index.vue'
import RatingStars from '@/components/RatingStars/index.vue'
import { useImageTools } from '@/composables/useImageTools'
import { formatRelativeTimeShort } from '@/utils/format'
import type { PlayerReviewData } from '@/types/review'

interface Props {
  rating?: number
  reviews: PlayerReviewData[]
}

defineProps<Props>()

defineEmits<{
  more: []
}>()

const { previewImages } = useImageTools()

const formatTime = (dateStr: string) => {
  if (!dateStr) return ''
  return formatRelativeTimeShort(dateStr)
}

</script>

<style lang="scss" scoped>
.reviews-summary {
  padding-bottom: var(--spacing-sm);
  margin-bottom: var(--spacing-md);
  border-bottom: 1rpx solid var(--color-border);
}

.rating-big {
  display: flex;
  align-items: baseline;
  gap: var(--spacing-sm);
}

.rating-value {
  font-size: 56rpx;
  font-weight: 800;
  color: var(--color-primary); // Gaming primary color
  letter-spacing: -1px;
  line-height: 1;
  text-shadow: 0 0 24rpx rgba(122, 204, 53, 0.4); // Neon Glow
}

.reviews-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.review-card {
  padding: var(--spacing-md);
  border-radius: var(--radius-lg);
  border: 1rpx solid var(--color-border);
  background: var(--color-bg-card);
  transition: all 0.2s;

  &.glass-card {
    background: rgba(255, 255, 255, 0.03);
    backdrop-filter: blur(8px);

    @media (prefers-color-scheme: dark) {
      background: rgba(0, 0, 0, 0.1);
      border-color: rgba(255, 255, 255, 0.08);
    }
  }

  &:hover {
    background: var(--color-bg-secondary);
    border-color: var(--color-primary-light);
  }
}

.review-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-md);
}

.review-user-info {
  flex: 1;
  min-width: 0;
}

.review-user-name {
  display: block;
  font-size: var(--font-sm);
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 4rpx;
}

.review-time {
  font-size: var(--font-xs);
  color: var(--color-text-placeholder);
  flex-shrink: 0;
}

.review-content {
  display: block;
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  line-height: 1.6;
  margin-bottom: var(--spacing-md);
}

.review-images {
  display: flex;
  gap: var(--spacing-xs);
  margin-top: var(--spacing-sm);
}

.review-image {
  width: 140rpx;
  height: 140rpx;
  border-radius: var(--radius-md);
  object-fit: cover;
  border: 1rpx solid var(--color-border);
  cursor: pointer;
  transition: transform 0.2s ease;

  &:hover {
    transform: scale(1.05);
    box-shadow: 0 4rpx 12rpx rgba(0,0,0,0.1);
  }
}

.more-images {
  width: 140rpx;
  height: 140rpx;
  border-radius: var(--radius-md);
  background: var(--color-bg-secondary);
  border: 1rpx solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;

  text {
    font-size: var(--font-sm);
    color: var(--color-text-secondary);
    font-weight: 600;
  }
}
</style>
