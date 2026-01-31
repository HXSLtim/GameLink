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
export interface AgreementSection {
  title: string
  text: string
}

export interface AgreementData {
  title: string
  updateTime: string
  sections: AgreementSection[]
}

interface Props {
  content: AgreementData
}

defineProps<Props>()
</script>

<style lang="scss" scoped>
.agreement-content {
  padding: 24rpx;
}

.header-section {
  margin-bottom: 32rpx;
  text-align: center;
}

.main-title {
  display: block;
  font-size: 36rpx;
  font-weight: bold;
  color: var(--color-text);
  margin-bottom: 12rpx;
}

.update-time {
  display: block;
  font-size: 24rpx;
  color: var(--color-text-placeholder);
}

.section {
  margin-bottom: 32rpx;
}

.section-title {
  display: block;
  font-size: 30rpx;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 16rpx;
}

.section-text {
  display: block;
  font-size: 28rpx;
  color: var(--color-text-secondary);
  line-height: 1.8;
  white-space: pre-wrap;
}
</style>
