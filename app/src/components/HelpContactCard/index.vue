<template>
  <view class="contact-section">
    <SectionHeader title="没有找到答案？" :show-more="false" />
    <GlCard :shadow="false" bordered clickable @click="$emit('click')">
      <view class="contact-card">
        <view class="contact-icon"><uv-icon name="chat" size="24" color="var(--color-text-secondary)" /></view>
        <view class="contact-info">
          <text class="contact-title">在线客服</text>
          <text class="contact-desc">工作时间 9:00-22:00</text>
        </view>
        <uv-icon name="arrow-right" size="16" color="var(--color-text-secondary)"></uv-icon>
      </view>
    </GlCard>
  </view>
</template>

<script setup lang="ts">
import SectionHeader from '@/components/SectionHeader/index.vue'
import GlCard from '@/components/gl/Card/index.vue'

defineEmits<{
  click: []
}>()
</script>

<style lang="scss" scoped>
.contact-section {
  padding: var(--spacing-md);
}

.contact-card {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.contact-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.contact-info {
  flex: 1;
}

.contact-title {
  display: block;
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: var(--spacing-xs);
}

.contact-desc {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}
</style>
