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
          <PriceTag :amount="service.price" amount-unit="yuan" :unit="service.unit || '局'" size="small" />
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
import PriceTag from '@/components/PriceTag/index.vue'
import type { PlayerServiceData } from '@/types/player'

interface Props {
  services: PlayerServiceData[]
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
  gap: var(--spacing-xs);
}

.service-option {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-sm);
  background: var(--color-bg-card);
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
  transition: all 0.2s;
  cursor: pointer;
  @include press-effect;
  
  &.selected {
    border-color: var(--color-border);
    background: var(--color-bg-secondary);
  }
}

.service-info {
  flex: 1;
  min-width: 0;
}

.service-name {
  display: block;
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: var(--spacing-xs);
}

.service-desc {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
  @include text-ellipsis;
}

.service-price {
  flex-shrink: 0;
  text-align: right;
}

.radio-box {
  width: 32rpx;
  height: 32rpx;
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
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
