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
          @tap.stop="previewImages(review.images || [], idx)"
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
import { useImageTools } from '@/composables/useImageTools'
import { formatDate } from '@/utils/format'
import type { ReviewCardData } from '@/types/review'

interface Props {
  review: ReviewCardData
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

const { previewImages } = useImageTools()

const formattedTime = computed(() => {
  if (!props.review.createdAt) return ''
  return formatDate(props.review.createdAt)
})

</script>

<style lang="scss" scoped>
.review-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  padding: var(--spacing-md);
  border: 1rpx solid var(--color-border);
}

.order-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding-bottom: var(--spacing-xs);
  border-bottom: 1rpx solid var(--color-border);
  margin-bottom: var(--spacing-sm);
  cursor: pointer;
  @include press-effect;
}

.order-detail {
  flex: 1;
  min-width: 0;
}

.player-name {
  display: block;
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 4rpx;
  @include text-ellipsis;
}

.order-desc {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  @include text-ellipsis;
}

.review-content {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.rating-row {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.rating-text {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  font-weight: 500;
}

.tags-row {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-xs);
}

.content-text {
  font-size: var(--font-sm);
  color: var(--color-text);
  line-height: 1.6;
}

.images-row {
  display: flex;
  gap: var(--spacing-xs);
}

.review-image {
  width: 160rpx;
  height: 160rpx;
  border-radius: var(--radius-md);
  object-fit: cover;
  border: 1rpx solid var(--color-border);
}

.review-time {
  font-size: var(--font-xs);
  color: var(--color-text-placeholder);
  align-self: flex-end;
}

.reply-section {
  margin-top: var(--spacing-sm);
  padding: var(--spacing-sm);
  background: var(--color-bg-secondary);
  border-radius: var(--radius-md);
  border: 1rpx solid var(--color-border);
}

.reply-label {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  margin-bottom: var(--spacing-xs);
  display: block;
}

.reply-content {
  font-size: var(--font-sm);
  color: var(--color-text);
  line-height: 1.5;
}

.pending-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-sm);
  padding-top: var(--spacing-xs);
  border-top: 1rpx solid var(--color-border);
  flex-wrap: wrap;
  row-gap: var(--spacing-xs);
}
</style>
