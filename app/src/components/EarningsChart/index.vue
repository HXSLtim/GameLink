<template>
  <GlCard :shadow="false" bordered>
    <template #title>
      <view class="chart-header">
        <text class="chart-title">{{ title }}</text>
        <text class="chart-total">共 ¥{{ formatMoney(total) }}</text>
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

export interface ChartDataItem {
  label: string
  value: number
}

interface Props {
  title: string
  total: number // 分
  data: ChartDataItem[]
}

const props = defineProps<Props>()

const formatMoney = (cents: number) => (cents / 100).toFixed(2)

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
  font-size: 30rpx;
  font-weight: 600;
  color: var(--color-text);
}

.chart-total {
  font-size: 28rpx;
  font-weight: 600;
  color: var(--color-primary);
}

.chart-container {
  padding: 24rpx 0;
}

.chart-bars {
  display: flex;
  align-items: flex-end;
  height: 200rpx;
  gap: 16rpx;
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
  background: linear-gradient(180deg, var(--color-primary) 0%, var(--color-primary-light, #4ADE80) 100%);
  border-radius: 8rpx 8rpx 0 0;
  min-height: 8rpx;
  margin-top: auto;
  transition: height 0.3s;
}

.chart-label {
  margin-top: 12rpx;
  font-size: 22rpx;
  color: var(--color-text-secondary);
}
</style>
