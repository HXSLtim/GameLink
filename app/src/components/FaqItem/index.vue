<template>
  <view class="faq-item" @tap="$emit('toggle')">
    <view class="faq-header">
      <text class="faq-question">{{ question }}</text>
      <uv-icon 
        :name="expanded ? 'arrow-up' : 'arrow-down'" 
        size="14" 
        color="var(--color-text-secondary)"
      ></uv-icon>
    </view>
    <view v-if="expanded" class="faq-answer">
      <text>{{ answer }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
interface Props {
  question: string
  answer: string
  expanded?: boolean
}

withDefaults(defineProps<Props>(), {
  expanded: false,
})

defineEmits<{
  toggle: []
}>()
</script>

<style lang="scss" scoped>
.faq-item {
  background: var(--color-bg-card);
  border-radius: 16rpx;
  margin-bottom: 16rpx;
  overflow: hidden;
  border: 2rpx solid var(--color-border);
}

.faq-header {
  display: flex;
  align-items: center;
  padding: 28rpx 24rpx;
}

.faq-question {
  flex: 1;
  font-size: 30rpx;
  font-weight: 500;
  color: var(--color-text);
  line-height: 1.4;
}

.faq-answer {
  padding: 0 24rpx 28rpx;
  border-top: 1rpx solid var(--color-border);
  padding-top: 20rpx;
  
  text {
    font-size: 28rpx;
    color: var(--color-text-secondary);
    line-height: 1.6;
    white-space: pre-wrap;
  }
}
</style>
