<template>
  <GlCard title="选择服务" required :shadow="false" bordered>
    <view class="services-list">
      <view 
        v-for="service in services" 
        :key="service.id"
        class="service-option"
        :class="{ selected: modelValue === service.id }"
        @tap="$emit('update:modelValue', service.id)"
      >
        <view class="service-info">
          <text class="service-name">{{ service.name }}</text>
          <text class="service-desc">{{ service.description }}</text>
        </view>
        <view class="service-price">
          <text class="price-value">¥{{ service.price }}</text>
          <text class="price-unit">/{{ service.unit || '局' }}</text>
        </view>
        <view class="radio-box" :class="{ checked: modelValue === service.id }">
          <uv-icon v-if="modelValue === service.id" name="checkbox-mark" size="12" color="#fff"></uv-icon>
        </view>
      </view>
    </view>
  </GlCard>
</template>

<script setup lang="ts">
import GlCard from '@/components/gl/Card/index.vue'

export interface ServiceOption {
  id: number
  name: string
  description?: string
  price: number
  unit?: string
}

interface Props {
  services: ServiceOption[]
  modelValue?: number
}

defineProps<Props>()

defineEmits<{
  'update:modelValue': [id: number]
}>()
</script>

<style lang="scss" scoped>
.services-list {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.service-option {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 24rpx;
  background: var(--color-bg-secondary);
  border-radius: 16rpx;
  border: 2rpx solid var(--color-border);
  transition: all 0.2s;
  
  &.selected {
    border-color: var(--color-primary);
    background: rgba(0, 210, 106, 0.08);
  }
}

.service-info {
  flex: 1;
  min-width: 0;
}

.service-name {
  display: block;
  font-size: 30rpx;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 6rpx;
}

.service-desc {
  font-size: 24rpx;
  color: var(--color-text-secondary);
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

.radio-box {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  border: 2rpx solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  
  &.checked {
    background: var(--color-primary);
    border-color: var(--color-primary);
  }
}
</style>
