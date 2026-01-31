<template>
  <GlCard :title="title" :shadow="false" bordered class="section-card">
    <view class="info-list">
      <view v-for="item in items" :key="item.label" class="info-row">
        <text class="info-label">{{ item.label }}</text>
        <view class="info-value-wrap">
          <text class="info-value">{{ item.value }}</text>
          <text 
            v-if="item.copyable" 
            class="copy-btn" 
            @tap="handleCopy(item.value)"
          >
            复制
          </text>
        </view>
      </view>
    </view>
  </GlCard>
</template>

<script setup lang="ts">
import GlCard from '@/components/gl/Card/index.vue'

export interface InfoItem {
  label: string
  value: string
  copyable?: boolean
}

interface Props {
  title: string
  items: InfoItem[]
}

defineProps<Props>()

const handleCopy = (text: string) => {
  uni.setClipboardData({
    data: text,
    success: () => {
      uni.showToast({ title: '复制成功', icon: 'success' })
    }
  })
}
</script>

<style lang="scss" scoped>
.section-card {
  margin: 20rpx 24rpx;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.info-label {
  font-size: 28rpx;
  color: var(--color-text-secondary);
  flex-shrink: 0;
}

.info-value-wrap {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.info-value {
  font-size: 28rpx;
  color: var(--color-text);
  text-align: right;
}

.copy-btn {
  font-size: 24rpx;
  color: var(--color-primary);
  padding: 4rpx 16rpx;
  border: 1rpx solid var(--color-primary);
  border-radius: 16rpx;
  
  &:active {
    opacity: 0.7;
  }
}
</style>
