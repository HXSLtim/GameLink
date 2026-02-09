<template>
  <GlCard :shadow="false" bordered>
    <template #header>
      <view class="chart-header">
        <text class="chart-title">{{ title }}</text>
        <view class="chart-total">
          <text>共</text>
          <PriceTag :amount="total" amount-unit="cents" size="small" />
        </view>
      </view>
    </template>
    
    <view class="chart-container">
      <view class="chart-bars">
        <view 
          v-for="(item, index) in data" 
          :key="index"
          class="chart-bar-wrap"
        >
          <view 
            class="chart-bar" 
            :style="{ height: getBarHeight(item.value) + '%' }"
          ></view>
          <text class="chart-label">{{ item.label }}</text>
        </view>
      </view>
    </view>
  </GlCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlCard from '@/components/gl/Card/index.vue'
import PriceTag from '@/components/PriceTag/index.vue'
import type { ChartDataItem } from '@/types/chart'

interface Props {
  title: string
  total: number // 分
  data: ChartDataItem[]
}

const props = defineProps<Props>()

const maxValue = computed(() => Math.max(...props.data.map(d => d.value), 1))

const getBarHeight = (value: number) => {
  return Math.round((value / maxValue.value) * 100)
}
</script>

<style lang="scss" scoped>
.chart-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.chart-title {
  font-size: var(--font-base);
  font-weight: 600;
  color: var(--color-text);
}

.chart-total {
  display: inline-flex;
  align-items: baseline;
  gap: 4rpx;
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-primary);
}

.chart-total :deep(.price-tag) {
  color: var(--color-primary);
}

.chart-container {
  padding: var(--spacing-md) 0;
}

.chart-bars {
  display: flex;
  align-items: flex-end;
  height: 200rpx;
  gap: var(--spacing-sm);
}

.chart-bar-wrap {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
}

.chart-bar {
  width: 100%;
  max-width: 60rpx;
  background: var(--color-primary);
  border-radius: var(--radius-sm) var(--radius-sm) 0 0;
  min-height: 8rpx;
  margin-top: auto;
  transition: height 0.3s;
}

.chart-label {
  margin-top: var(--spacing-xs);
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}
</style>
