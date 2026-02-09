<template>
  <SectionCard title="服务项目">
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
          <PriceTag :amount="service.price" amount-unit="yuan" :unit="service.unit || '局'" size="small" />
        </view>
        <view v-if="selectedId === service.id" class="check-icon">
          <uv-icon name="checkbox-mark" size="16" color="var(--color-primary)"></uv-icon>
        </view>
      </view>
      
      <GlEmpty v-if="!services?.length" title="暂无服务项目" compact />
    </view>
  </SectionCard>
</template>

<script setup lang="ts">
import SectionCard from '@/components/SectionCard/index.vue'
import PriceTag from '@/components/PriceTag/index.vue'
import GlEmpty from '@/components/gl/Empty/index.vue'
import type { PlayerServiceData } from '@/types/player'

interface Props {
  services: PlayerServiceData[]
  selectedId?: number
}

defineProps<Props>()

defineEmits<{
  select: [service: PlayerServiceData]
}>()
</script>

<style lang="scss" scoped>
.services-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.service-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  background: var(--color-bg-card);
  border-radius: var(--radius-lg);
  border: 1rpx solid var(--color-border);
  cursor: pointer;
  transition: border-color 0.2s ease, background 0.2s ease, box-shadow 0.2s ease;
  @include press-effect;

  &:hover {
    border-color: var(--color-primary);
    background: var(--color-bg-secondary);
  }

  &--selected {
    border-color: var(--color-primary);
    background: var(--color-primary-tint);
    box-shadow: var(--shadow-sm);
  }
}

.service-main {
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

.check-icon {
  flex-shrink: 0;
  width: 32rpx;
  height: 32rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-secondary);
  border: 1rpx solid var(--color-border);
  border-radius: var(--radius-md);
}
</style>
