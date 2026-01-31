<template>
  <view class="review-card">
    <!-- 订单信息 -->
    <view class="order-info" @tap="$emit('order-click')">
      <GlAvatar :src="review.player.avatar" :text="review.player.nickname" size="medium" />
      <view class="order-detail">
        <text class="player-name">{{ review.player.nickname }}</text>
        <text class="order-desc">{{ review.gameName }} · {{ review.serviceName }}</text>
      </view>
      <uv-icon name="arrow-right" size="16" color="var(--color-text-secondary)"></uv-icon>
    </view>
    
    <!-- 评价内容 -->
    <view class="review-content">
      <view class="rating-row">
        <RatingStars :rating="review.rating" size="small" />
        <text class="rating-text">{{ ratingText }}</text>
      </view>
      
      <view v-if="review.tags?.length" class="tags-row">
        <GlTag v-for="tag in review.tags" :key="tag" size="mini" plain>{{ tag }}</GlTag>
      </view>
      
      <text v-if="review.content" class="content-text">{{ review.content }}</text>
      
      <view v-if="review.images?.length" class="images-row">
        <image 
          v-for="(img, idx) in review.images" 
          :key="idx"
          :src="img"
          mode="aspectFill"
          class="review-image"
          @tap.stop="previewImage(idx)"
        />
      </view>
      
      <text class="review-time">{{ formattedTime }}</text>
    </view>
    
    <!-- 商家回复 -->
    <view v-if="review.reply" class="reply-section">
      <text class="reply-label">陪玩师回复：</text>
      <text class="reply-content">{{ review.reply }}</text>
    </view>
    
    <!-- 待评价操作 -->
    <view v-if="showActions" class="pending-actions">
      <GlButton type="default" size="small" plain @click="$emit('skip')">跳过</GlButton>
      <GlButton type="primary" size="small" @click="$emit('write')">去评价</GlButton>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlAvatar from '@/components/gl/Avatar/index.vue'
import GlTag from '@/components/gl/Tag/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
import RatingStars from '@/components/RatingStars/index.vue'

export interface ReviewPlayerInfo {
  id: number
  nickname: string
  avatar?: string
}

export interface ReviewData {
  id: number
  orderId: number
  player: ReviewPlayerInfo
  gameName: string
  serviceName: string
  rating: number
  tags?: string[]
  content?: string
  images?: string[]
  reply?: string
  createdAt: string
}

interface Props {
  review: ReviewData
  showActions?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showActions: false,
})

defineEmits<{
  'order-click': []
  skip: []
  write: []
}>()

const ratingTexts = ['', '非常差', '较差', '一般', '满意', '非常满意']
const ratingText = computed(() => ratingTexts[props.review.rating] || '')

const formattedTime = computed(() => {
  if (!props.review.createdAt) return ''
  const date = new Date(props.review.createdAt)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
})

const previewImage = (index: number) => {
  if (props.review.images) {
    uni.previewImage({
      urls: props.review.images,
      current: props.review.images[index],
    })
  }
}
</script>

<style lang="scss" scoped>
.review-card {
  background: var(--color-bg-card);
  border-radius: 20rpx;
  padding: 24rpx;
  margin-bottom: 20rpx;
  border: 2rpx solid var(--color-border);
}

.order-info {
  display: flex;
  align-items: center;
  gap: 16rpx;
  padding-bottom: 20rpx;
  border-bottom: 1rpx solid var(--color-border);
  margin-bottom: 20rpx;
}

.order-detail {
  flex: 1;
  min-width: 0;
}

.player-name {
  display: block;
  font-size: 30rpx;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 4rpx;
}

.order-desc {
  font-size: 24rpx;
  color: var(--color-text-secondary);
}

.review-content {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.rating-row {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.rating-text {
  font-size: 26rpx;
  color: var(--color-primary);
  font-weight: 500;
}

.tags-row {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
}

.content-text {
  font-size: 28rpx;
  color: var(--color-text);
  line-height: 1.6;
}

.images-row {
  display: flex;
  gap: 12rpx;
}

.review-image {
  width: 160rpx;
  height: 160rpx;
  border-radius: 12rpx;
  object-fit: cover;
}

.review-time {
  font-size: 24rpx;
  color: var(--color-text-placeholder);
}

.reply-section {
  margin-top: 20rpx;
  padding: 20rpx;
  background: var(--color-bg-secondary);
  border-radius: 12rpx;
}

.reply-label {
  font-size: 24rpx;
  color: var(--color-text-secondary);
  margin-bottom: 8rpx;
  display: block;
}

.reply-content {
  font-size: 26rpx;
  color: var(--color-text);
  line-height: 1.5;
}

.pending-actions {
  display: flex;
  justify-content: flex-end;
  gap: 20rpx;
  margin-top: 20rpx;
  padding-top: 20rpx;
  border-top: 1rpx solid var(--color-border);
}
</style>
