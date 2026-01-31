<template>
  <uv-popup :show="show" mode="bottom" round="20" @close="$emit('close')">
    <view class="review-modal">
      <view class="modal-header">
        <text class="modal-title">评价订单</text>
        <uv-icon name="close" size="24" @click="$emit('close')"></uv-icon>
      </view>
      
      <view class="review-form">
        <!-- 评分 -->
        <view class="rating-section">
          <text class="rating-label">服务评分</text>
          <view class="rating-stars">
            <text 
              v-for="i in 5" 
              :key="i"
              class="star"
              :class="{ active: rating >= i }"
              @tap="$emit('update:rating', i)"
            >★</text>
          </view>
        </view>
        
        <!-- 标签 -->
        <view v-if="tags.length > 0" class="tags-section">
          <text class="tags-label">选择标签</text>
          <view class="tags-list">
            <GlTag 
              v-for="tag in availableTags" 
              :key="tag"
              :type="tags.includes(tag) ? 'primary' : 'default'"
              size="small"
              @click="toggleTag(tag)"
            >
              {{ tag }}
            </GlTag>
          </view>
        </view>
        
        <!-- 内容 -->
        <view class="content-section">
          <textarea 
            :value="content"
            class="review-textarea"
            placeholder="分享您的服务体验（选填）"
            :maxlength="500"
            @input="(e: any) => $emit('update:content', e.detail.value)"
          />
        </view>
      </view>
      
      <GlButton type="primary" block size="large" :loading="loading" @click="$emit('submit')">
        提交评价
      </GlButton>
    </view>
  </uv-popup>
</template>

<script setup lang="ts">
import GlTag from '@/components/gl/Tag/index.vue'
import GlButton from '@/components/gl/Button/index.vue'

interface Props {
  show: boolean
  rating: number
  content?: string
  tags?: string[]
  availableTags?: string[]
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  rating: 5,
  content: '',
  tags: () => [],
  availableTags: () => ['服务态度好', '技术很棒', '沟通顺畅', '准时守约', '非常耐心'],
  loading: false,
})

const emit = defineEmits<{
  close: []
  submit: []
  'update:rating': [value: number]
  'update:content': [value: string]
  'update:tags': [value: string[]]
}>()

const toggleTag = (tag: string) => {
  const newTags = props.tags.includes(tag)
    ? props.tags.filter(t => t !== tag)
    : [...props.tags, tag]
  emit('update:tags', newTags)
}
</script>

<style lang="scss" scoped>
.review-modal {
  padding: 32rpx;
  padding-bottom: calc(32rpx + env(safe-area-inset-bottom));
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 32rpx;
}

.modal-title {
  font-size: 34rpx;
  font-weight: 600;
  color: var(--color-text);
}

.review-form {
  margin-bottom: 32rpx;
}

.rating-section {
  display: flex;
  align-items: center;
  gap: 24rpx;
  margin-bottom: 32rpx;
}

.rating-label,
.tags-label {
  font-size: 28rpx;
  color: var(--color-text-secondary);
}

.rating-stars {
  display: flex;
  gap: 8rpx;
}

.star {
  font-size: 48rpx;
  color: var(--color-border);
  transition: color 0.2s;

  &.active {
    color: #FFB800;
  }
}

.tags-section {
  margin-bottom: 32rpx;
}

.tags-label {
  display: block;
  margin-bottom: 16rpx;
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.content-section {
  margin-bottom: 16rpx;
}

.review-textarea {
  width: 100%;
  height: 200rpx;
  padding: 20rpx;
  background: var(--color-bg-secondary);
  border-radius: 16rpx;
  font-size: 28rpx;
  color: var(--color-text);
  box-sizing: border-box;
}
</style>
