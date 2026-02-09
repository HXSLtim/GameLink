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
  margin-bottom: var(--spacing-sm);

  @include desktop {
    margin-bottom: var(--spacing-md);
  }
}

.section-title {
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
  position: relative;
  padding-left: var(--spacing-sm);
  
  &::before {
    content: '';
    position: absolute;
    left: 0;
    top: 50%;
    transform: translateY(-50%);
    width: 4rpx;
    height: 20rpx;
    background: var(--color-primary);
    border-radius: var(--radius-full);
  }

  @include desktop {
    font-size: 17px;
    font-weight: 700;
    padding-left: 14px;

    &::before {
      width: 3px;
      height: 16px;
    }
  }
}

.section-more {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--radius-md);
  border: 1rpx solid var(--color-border);
  transition: all 0.2s;
  cursor: pointer;
  @include press-effect;
  
  &:hover {
    background: var(--color-bg-secondary);
    border-color: var(--color-primary);
    
    text {
      color: var(--color-primary);
    }
  }
  
  &:active {
    background: var(--color-bg-secondary);
  }
  
  text {
    font-size: var(--font-xs);
    color: var(--color-text-secondary);
    transition: color 0.2s;
  }

  @include desktop {
    padding: 6px 14px;

    text {
      font-size: 13px;
    }
  }
}
</style>
