<template>
  <view class="status-card" :class="status">
    <view class="status-icon">{{ icon }}</view>
    <view class="status-info">
      <text class="status-title">{{ title }}</text>
      <text class="status-desc">{{ description }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type CertStatus = 'none' | 'pending' | 'approved' | 'rejected'

interface Props {
  status: CertStatus
}

const props = defineProps<Props>()

const config = computed(() => {
  const configs: Record<CertStatus, { icon: string; title: string; desc: string }> = {
    none: { icon: '📝', title: '未认证', desc: '完成认证后即可开始接单' },
    pending: { icon: '⏳', title: '审核中', desc: '预计1-3个工作日内完成审核' },
    approved: { icon: '✅', title: '已认证', desc: '恭喜您已通过认证' },
    rejected: { icon: '❌', title: '认证失败', desc: '请根据反馈修改后重新提交' },
  }
  return configs[props.status]
})

const icon = computed(() => config.value.icon)
const title = computed(() => config.value.title)
const description = computed(() => config.value.desc)
</script>

<style lang="scss" scoped>
.status-card {
  display: flex;
  align-items: center;
  gap: 24rpx;
  padding: 32rpx;
  margin: 24rpx;
  border-radius: 20rpx;
  background: var(--color-bg-card);
  border: 2rpx solid var(--color-border);
  
  &.none {
    background: linear-gradient(135deg, rgba(156, 163, 175, 0.1) 0%, rgba(156, 163, 175, 0.05) 100%);
    border-color: #9CA3AF;
  }
  
  &.pending {
    background: linear-gradient(135deg, rgba(245, 158, 11, 0.1) 0%, rgba(245, 158, 11, 0.05) 100%);
    border-color: #F59E0B;
  }
  
  &.approved {
    background: linear-gradient(135deg, rgba(0, 210, 106, 0.1) 0%, rgba(0, 210, 106, 0.05) 100%);
    border-color: var(--color-primary);
  }
  
  &.rejected {
    background: linear-gradient(135deg, rgba(239, 68, 68, 0.1) 0%, rgba(239, 68, 68, 0.05) 100%);
    border-color: #EF4444;
  }
}

.status-icon {
  font-size: 56rpx;
}

.status-info {
  flex: 1;
}

.status-title {
  display: block;
  font-size: 32rpx;
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: 8rpx;
}

.status-desc {
  font-size: 26rpx;
  color: var(--color-text-secondary);
}
</style>
