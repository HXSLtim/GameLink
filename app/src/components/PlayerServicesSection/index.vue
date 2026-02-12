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
        <view class="service-icon">
          <uv-icon name="server-fill" size="24" color="var(--color-primary)" />
        </view>
        <view class="service-main">
          <text class="service-name">{{ service.name }}</text>
          <text v-if="service.description" class="service-desc">{{ service.description }}</text>
        </view>
        <view class="service-price">
          <PriceTag :amount="service.price" amount-unit="yuan" :unit="service.unit || '局'" size="medium" :color="selectedId === service.id ? 'var(--color-primary)' : ''" />
        </view>
        <view class="check-icon" :class="{ checked: selectedId === service.id }">
          <uv-icon v-if="selectedId === service.id" name="checkbox-mark" size="14" color="#fff"></uv-icon>
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
  gap: var(--spacing-sm);
}

.service-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  background: rgba(255, 255, 255, 0.05); // Glass effect base
  border-radius: var(--radius-lg);
  border: 1rpx solid var(--color-border);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;

  // PC Hover Effect
  @media (hover: hover) {
    &:hover {
      transform: translateY(-2px);
      background: var(--color-bg-secondary);
      border-color: var(--color-primary-light);
      box-shadow: var(--shadow-sm);
    }
  }

  // Selected State (Glow)
  &--selected {
    border-color: var(--color-primary);
    background: rgba(122, 204, 53, 0.1); // Slightly stronger tint
    box-shadow: var(--shadow-primary); // Use the variable

    .service-name {
      color: var(--color-primary);
      text-shadow: 0 0 12rpx rgba(122, 204, 53, 0.3); // Add text glow
    }
  }
}

.service-icon {
  width: 80rpx;
  height: 80rpx;
  border-radius: var(--radius-md);
  background: var(--color-bg-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.service-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.service-name {
  font-size: var(--font-base);
  font-weight: 600;
  color: var(--color-text);
  transition: color 0.2s;
}

.service-desc {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
  @include text-ellipsis;
}

.service-price {
  flex-shrink: 0;
  margin-right: var(--spacing-xs);
}

.check-icon {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  border: 2rpx solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  background: var(--color-bg-card);

  &.checked {
    background: var(--color-primary);
    border-color: var(--color-primary);
    box-shadow: 0 0 10rpx rgba(122, 204, 53, 0.4);
  }
}
</style>
