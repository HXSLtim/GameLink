<template>
  <GlCard :shadow="false" bordered class="section-card">
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
      <view v-for="review in reviews" :key="review.id" class="review-item">
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
            @tap="previewImage(review.images!, idx)"
          />
          <view v-if="review.images.length > 3" class="more-images">
            <text>+{{ review.images.length - 3 }}</text>
          </view>
        </view>
      </view>
      
      <GlEmpty v-if="!reviews?.length" title="暂无评价" description="下单体验后来评价吧" compact />
    </view>
  </GlCard>
</template>

<script setup lang="ts">
import GlCard from '@/components/gl/Card/index.vue'
import GlAvatar from '@/components/gl/Avatar/index.vue'
import GlEmpty from '@/components/gl/Empty/index.vue'
import SectionHeader from '@/components/SectionHeader/index.vue'
import RatingStars from '@/components/RatingStars/index.vue'

export interface ReviewData {
  id: number
  userId: number
  userName: string
  userAvatar?: string
  rating: number
  content: string
  images?: string[]
  createdAt: string
}

interface Props {
  rating?: number
  reviews: ReviewData[]
}

defineProps<Props>()

defineEmits<{
  more: []
}>()

const formatTime = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  if (diff < 60 * 1000) return '刚刚'
  if (diff < 60 * 60 * 1000) return `${Math.floor(diff / 60 / 1000)}分钟前`
  if (diff < 24 * 60 * 60 * 1000) return `${Math.floor(diff / 60 / 60 / 1000)}小时前`
  if (diff < 7 * 24 * 60 * 60 * 1000) return `${Math.floor(diff / 24 / 60 / 60 / 1000)}天前`
  
  return `${date.getMonth() + 1}/${date.getDate()}`
}

const previewImage = (images: string[], index: number) => {
  uni.previewImage({
    urls: images,
    current: images[index],
  })
}
</script>

<style lang="scss" scoped>
.section-card {
  margin: 0 24rpx 20rpx;
}

.reviews-summary {
  padding-bottom: 24rpx;
  border-bottom: 1rpx solid var(--color-border);
  margin-bottom: 24rpx;
}

.rating-big {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.rating-value {
  font-size: 56rpx;
  font-weight: 800;
  color: var(--color-primary);
}

.reviews-list {
  display: flex;
  flex-direction: column;
  gap: 24rpx;
}

.review-item {
  padding-bottom: 24rpx;
  border-bottom: 1rpx solid var(--color-border);
  
  &:last-child {
    border-bottom: none;
    padding-bottom: 0;
  }
}

.review-header {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-bottom: 16rpx;
}

.review-user-info {
  flex: 1;
  min-width: 0;
}

.review-user-name {
  display: block;
  font-size: 28rpx;
  font-weight: 500;
  color: var(--color-text);
  margin-bottom: 4rpx;
}

.review-time {
  font-size: 24rpx;
  color: var(--color-text-placeholder);
  flex-shrink: 0;
}

.review-content {
  font-size: 28rpx;
  color: var(--color-text);
  line-height: 1.6;
  margin-bottom: 16rpx;
}

.review-images {
  display: flex;
  gap: 12rpx;
}

.review-image {
  width: 160rpx;
  height: 160rpx;
  border-radius: 12rpx;
  object-fit: cover;
}

.more-images {
  width: 160rpx;
  height: 160rpx;
  border-radius: 12rpx;
  background: var(--color-bg-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  
  text {
    font-size: 28rpx;
    color: var(--color-text-secondary);
  }
}
</style>
