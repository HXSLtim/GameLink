<template>
  <SectionCard :title="title" margin="var(--spacing-sm) var(--spacing-md)">
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
  </SectionCard>
</template>

<script setup lang="ts">
import SectionCard from '@/components/SectionCard/index.vue'
import type { InfoItem } from '@/types/order'
import { copyToClipboard } from '@/utils'

interface Props {
  title: string
  items: InfoItem[]
}

defineProps<Props>()

const handleCopy = (text: string) => {
  copyToClipboard(text).catch(() => undefined)
}
</script>

<style lang="scss" scoped>
.info-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--spacing-sm);
}

.info-label {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
  flex-shrink: 0;
}

.info-value-wrap {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.info-value {
  font-size: var(--font-sm);
  color: var(--color-text);
  text-align: right;
}

.copy-btn {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
  padding: 2rpx var(--spacing-sm);
  border: 1rpx solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-bg-secondary);
  @include press-effect;
  
  &:active {
    opacity: 0.7;
  }
}
</style>
