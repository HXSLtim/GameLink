<template>
  <GlCard title="我的订单" :shadow="false" bordered class="order-section">
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
  </GlCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlCard from '@/components/gl/Card/index.vue'

interface Props {
  pendingCount?: number
  inProgressCount?: number
  toReviewCount?: number
}

const props = withDefaults(defineProps<Props>(), {
  pendingCount: 0,
  inProgressCount: 0,
  toReviewCount: 0,
})

defineEmits<{
  click: [status: string]
  'view-all': []
}>()

const tabs = computed(() => [
  { 
    key: 'pending', 
    label: '待支付', 
    icon: 'bag', 
    iconColor: 'var(--color-primary)',
    badge: props.pendingCount || undefined,
  },
  { 
    key: 'in_progress', 
    label: '进行中', 
    icon: 'clock', 
    iconColor: '#3B82F6',
    badge: props.inProgressCount || undefined,
  },
  { 
    key: 'completed', 
    label: '待评价', 
    icon: 'checkmark-circle', 
    iconColor: '#10B981',
    badge: props.toReviewCount || undefined,
  },
  { 
    key: 'refund', 
    label: '退款', 
    icon: 'reload', 
    iconColor: '#F59E0B',
    badge: undefined,
  },
])
</script>

<style lang="scss" scoped>
.order-section {
  margin: 20rpx 28rpx;
}

.section-more {
  display: flex;
  align-items: center;
  gap: 4rpx;
  font-size: 26rpx;
  color: var(--color-text-secondary);
}

.order-tabs {
  display: flex;
  padding-top: 8rpx;
}

.order-tab {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16rpx 8rpx;
  border-radius: 16rpx;
  transition: all 0.2s;
  
  &:active {
    background: var(--color-bg-secondary);
    transform: scale(0.95);
  }
}

.tab-icon-wrap {
  position: relative;
  width: 56rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-secondary);
  border-radius: 16rpx;
  margin-bottom: 8rpx;
}

.tab-text {
  font-size: 24rpx;
  font-weight: 500;
  color: var(--color-text-secondary);
}

.tab-badge {
  position: absolute;
  top: -8rpx;
  right: -8rpx;
  min-width: 32rpx;
  height: 32rpx;
  padding: 0 8rpx;
  background: linear-gradient(135deg, #EF4444 0%, #DC2626 100%);
  border-radius: 18rpx;
  font-size: 20rpx;
  font-weight: 600;
  color: #FFFFFF;
  text-align: center;
  line-height: 32rpx;
  box-shadow: 0 4rpx 12rpx rgba(239, 68, 68, 0.4);
}
</style>
