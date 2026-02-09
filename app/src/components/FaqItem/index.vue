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
  border-radius: var(--radius-md);
  margin-bottom: var(--spacing-sm);
  overflow: hidden;
  border: 1rpx solid var(--color-border);
  cursor: pointer;
  @include press-effect;
}

.faq-header {
  display: flex;
  align-items: center;
  padding: var(--spacing-md) var(--spacing-md);
}

.faq-question {
  flex: 1;
  font-size: var(--font-md);
  font-weight: 500;
  color: var(--color-text);
  line-height: 1.4;
}

.faq-answer {
  padding: 0 var(--spacing-md) var(--spacing-md);
  border-top: 1rpx solid var(--color-border);
  padding-top: var(--spacing-sm);
  
  text {
    font-size: var(--font-sm);
    color: var(--color-text-secondary);
    line-height: 1.6;
    white-space: pre-wrap;
  }
}
</style>
