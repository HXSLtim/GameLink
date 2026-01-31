<template>
  <view class="section-header">
    <text class="section-title">{{ title }}</text>
    <view v-if="showMore" class="section-more" @click="$emit('more')">
      <text>{{ moreText }}</text>
      <uv-icon name="arrow-right" size="14" color="var(--color-text-secondary)"></uv-icon>
    </view>
  </view>
</template>

<script setup lang="ts">
interface Props {
  title: string
  showMore?: boolean
  moreText?: string
}

withDefaults(defineProps<Props>(), {
  showMore: true,
  moreText: '更多',
})

defineEmits<{
  more: []
}>()
</script>

<style lang="scss" scoped>
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24rpx;
}

.section-title {
  font-size: 34rpx;
  font-weight: 700;
  color: var(--color-text);
  position: relative;
  padding-left: 20rpx;
  
  &::before {
    content: '';
    position: absolute;
    left: 0;
    top: 50%;
    transform: translateY(-50%);
    width: 6rpx;
    height: 28rpx;
    background: linear-gradient(180deg, var(--color-primary) 0%, var(--color-primary-light, #4ADE80) 100%);
    border-radius: 3rpx;
  }
}

.section-more {
  display: flex;
  align-items: center;
  gap: 4rpx;
  padding: 8rpx 16rpx;
  border-radius: 16rpx;
  transition: all 0.2s;
  
  &:active {
    background: var(--color-bg-secondary);
  }
  
  text {
    font-size: 26rpx;
    color: var(--color-text-secondary);
  }
}
</style>
