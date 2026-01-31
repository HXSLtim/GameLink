<template>
  <GlCard title="服务项目" :shadow="false" bordered class="section-card">
    <view class="services-list">
      <view 
        v-for="service in services" 
        :key="service.id"
        class="service-item"
        :class="{ 'service-item--selected': selectedId === service.id }"
        @tap="$emit('select', service)"
      >
        <view class="service-main">
          <text class="service-name">{{ service.name }}</text>
          <text v-if="service.description" class="service-desc">{{ service.description }}</text>
        </view>
        <view class="service-price">
          <text class="price-value">¥{{ service.price }}</text>
          <text class="price-unit">/{{ service.unit || '局' }}</text>
        </view>
        <view v-if="selectedId === service.id" class="check-icon">
          <uv-icon name="checkbox-mark" size="16" color="var(--color-primary)"></uv-icon>
        </view>
      </view>
      
      <GlEmpty v-if="!services?.length" title="暂无服务项目" compact />
    </view>
  </GlCard>
</template>

<script setup lang="ts">
import GlCard from '@/components/gl/Card/index.vue'
import GlEmpty from '@/components/gl/Empty/index.vue'

export interface ServiceData {
  id: number
  name: string
  description?: string
  price: number
  unit?: string
}

interface Props {
  services: ServiceData[]
  selectedId?: number
}

defineProps<Props>()

defineEmits<{
  select: [service: ServiceData]
}>()
</script>

<style lang="scss" scoped>
.section-card {
  margin: 0 24rpx 20rpx;
}

.services-list {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.service-item {
  display: flex;
  align-items: center;
  gap: 16rpx;
  padding: 24rpx;
  background: var(--color-bg-secondary);
  border-radius: 16rpx;
  border: 2rpx solid transparent;
  transition: all 0.2s;
  
  &:active {
    transform: scale(0.99);
  }
  
  &--selected {
    border-color: var(--color-primary);
    background: rgba(0, 210, 106, 0.05);
  }
}

.service-main {
  flex: 1;
  min-width: 0;
}

.service-name {
  display: block;
  font-size: 30rpx;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 8rpx;
}

.service-desc {
  font-size: 24rpx;
  color: var(--color-text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.service-price {
  flex-shrink: 0;
  text-align: right;
}

.price-value {
  font-size: 32rpx;
  font-weight: 700;
  color: var(--color-primary);
}

.price-unit {
  font-size: 22rpx;
  color: var(--color-text-secondary);
}

.check-icon {
  flex-shrink: 0;
  width: 40rpx;
  height: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 210, 106, 0.1);
  border-radius: 50%;
}
</style>
