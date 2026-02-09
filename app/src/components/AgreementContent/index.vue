<template>
  <view class="agreement-content">
    <!-- 标题区域 -->
    <view class="header-section">
      <text class="main-title">{{ content.title }}</text>
      <text class="update-time">更新时间：{{ content.updateTime }}</text>
    </view>

    <!-- 内容区域 -->
    <view
      v-for="(section, index) in content.sections"
      :key="index"
      class="section"
    >
      <text class="section-title">{{ section.title }}</text>
      <text class="section-text">{{ section.text }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import type { AgreementData } from '@/types/agreement'

interface Props {
  content: AgreementData
}

defineProps<Props>()
</script>

<style lang="scss" scoped>
.agreement-content {
  padding: var(--spacing-md);
}

.header-section {
  margin-bottom: var(--spacing-lg);
  text-align: center;
}

.main-title {
  display: block;
  font-size: var(--font-lg);
  font-weight: bold;
  color: var(--color-text);
  margin-bottom: var(--spacing-xs);
}

.update-time {
  display: block;
  font-size: var(--font-sm);
  color: var(--color-text-placeholder);
}

.section {
  margin-bottom: var(--spacing-lg);
}

.section-title {
  display: block;
  font-size: var(--font-base);
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: var(--spacing-sm);
}

.section-text {
  display: block;
  font-size: var(--font-md);
  color: var(--color-text-secondary);
  line-height: 1.8;
  white-space: pre-wrap;
}
</style>
