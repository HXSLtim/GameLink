<template>
  <uv-popup :show="show" mode="bottom" round="16" @close="$emit('close')">
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
            <view
              v-for="i in 5"
              :key="i"
              class="star"
              :class="{ active: rating >= i }"
              @tap="$emit('update:rating', i)"
            >
              <uv-icon
                :name="rating >= i ? 'star-fill' : 'star'"
                size="24"
                :color="rating >= i ? 'var(--color-gold)' : 'var(--color-text-disabled)'"
              />
            </view>
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
          <GlInput
            class="review-input"
            :model-value="content"
            type="textarea"
            size="small"
            placeholder="分享您的服务体验（选填）"
            :maxlength="500"
            @update:modelValue="(value) => $emit('update:content', value)"
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
import GlInput from '@/components/gl/Input/index.vue'

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
  padding: var(--spacing-md);
  padding-bottom: calc(var(--spacing-md) + env(safe-area-inset-bottom));
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--spacing-md);
}

.modal-title {
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
}

.review-form {
  margin-bottom: var(--spacing-md);
}

.rating-section {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.rating-label,
.tags-label {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}

.rating-stars {
  display: flex;
  gap: var(--spacing-xs);
}

.star {
  transition: opacity 0.2s;
  cursor: pointer;
  @include press-effect;
}

.tags-section {
  margin-bottom: var(--spacing-lg);
}

.tags-label {
  display: block;
  margin-bottom: var(--spacing-sm);
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm);
}

.content-section {
  margin-bottom: var(--spacing-sm);
}

.review-input {
  :deep(.gl-input__textarea) {
    min-height: 200rpx;
  }
}
</style>
