<template>
  <SectionCard title="我的订单" margin="var(--spacing-sm) var(--spacing-md)">
    <template #extra>
      <view class="section-more" @tap="$emit('view-all')">
        <text>全部订单</text>
        <uv-icon name="arrow-right" size="14" color="var(--color-text-secondary)"></uv-icon>
      </view>
    </template>
    
    <view class="order-tabs">
      <view
        v-for="tab in tabs"
        :key="tab.key"
        class="order-tab"
        @tap="$emit('click', tab.key)"
      >
        <view class="tab-icon-wrap">
          <uv-icon :name="tab.icon" size="24" :color="tab.iconColor"></uv-icon>
          <view v-if="tab.badge" class="tab-badge">{{ tab.badge }}</view>
        </view>
        <text class="tab-text">{{ tab.label }}</text>
      </view>
    </view>
  </SectionCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import SectionCard from '@/components/SectionCard/index.vue'
import type { OrderCountSummary, OrderQuickEntryStatus } from '@/types/order'

interface Props {
  counts?: OrderCountSummary
}

const props = withDefaults(defineProps<Props>(), {
  counts: () => ({
    pending: 0,
    inProgress: 0,
    toReview: 0,
    refunding: 0,
  }),
})

defineEmits<{
  click: [status: OrderQuickEntryStatus]
  'view-all': []
}>()

// 订单状态图标统一使用中性色，状态由文字标签传达
const iconColor = 'var(--color-text-secondary)'

const tabs = computed<Array<{
  key: OrderQuickEntryStatus
  label: string
  icon: string
  iconColor: string
  badge?: number | string
}>>(() => [
  { 
    key: 'pending', 
    label: '待支付', 
    icon: 'bag', 
    iconColor,
    badge: props.counts?.pending || undefined,
  },
  { 
    key: 'in_progress', 
    label: '进行中', 
    icon: 'clock', 
    iconColor,
    badge: props.counts?.inProgress || undefined,
  },
  { 
    key: 'completed', 
    label: '待评价', 
    icon: 'checkmark-circle', 
    iconColor,
    badge: props.counts?.toReview || undefined,
  },
  {
    key: 'refunding',
    label: '退款中',
    icon: 'reload',
    iconColor,
    badge: props.counts?.refunding || undefined,
  },
])
</script>

<style lang="scss" scoped>
.section-more {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
  cursor: pointer;
  @include press-effect;
}

.order-tabs {
  display: flex;
  padding-top: var(--spacing-xs);
}

.order-tab {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: var(--spacing-xs) var(--spacing-xs);
  border-radius: var(--radius-sm);
  transition: background 0.2s ease, transform 0.15s ease;
  cursor: pointer;
  @include press-effect;

  &:hover {
    background: var(--color-bg-secondary);
  }
  
  &:active {
    background: var(--color-bg-secondary);
    transform: scale(0.96);
  }
}

.tab-icon-wrap {
  position: relative;
  width: 40rpx;
  height: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-secondary);
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
  margin-bottom: var(--spacing-xs);
}

.tab-text {
  font-size: var(--font-xs);
  font-weight: 500;
  color: var(--color-text-secondary);
}

.tab-badge {
  position: absolute;
  top: -8rpx;
  right: -8rpx;
  min-width: 28rpx;
  height: 28rpx;
  padding: 0 var(--spacing-xs);
  background: var(--color-error);
  border-radius: var(--radius-full);
  font-size: var(--font-xs);
  font-weight: 600;
  color: #FFFFFF;
  text-align: center;
  line-height: 28rpx;
}
</style>
