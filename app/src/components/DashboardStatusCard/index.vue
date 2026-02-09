<template>
  <StatusCard
    :status-class="status"
    :border-width="borderWidth"
    :border-color="borderColor"
    justify="space-between"
    gap="0"
  >
    <view class="status-left">
      <view class="avatar-wrap">
        <GlAvatar :src="avatar" :text="nickname" :size="80" :status="status === 'online' ? 'online' : status === 'busy' ? 'busy' : undefined" bordered />
      </view>
      <view class="status-info">
        <text class="player-name">{{ nickname }}</text>
        <text class="status-text">{{ statusText }}</text>
      </view>
    </view>
    <GlButton :type="status === 'online' ? 'default' : 'primary'" size="small" round @click="$emit('toggle')">
      {{ status === 'online' ? '下线' : '上线接单' }}
    </GlButton>
  </StatusCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlAvatar from '@/components/gl/Avatar/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
import StatusCard from '@/components/StatusCard/index.vue'
import { getDashboardStatusPreset, type DashboardStatus } from '@/components/StatusCard/presets'

interface Props {
  status: DashboardStatus
  nickname: string
  avatar?: string
}

const props = defineProps<Props>()

defineEmits<{
  toggle: []
}>()

const config = computed(() => getDashboardStatusPreset(props.status))
const statusText = computed(() => config.value.text)
const borderColor = computed(() => config.value.borderColor)
const borderWidth = computed(() => config.value.borderWidth)
</script>

<style lang="scss" scoped>
.status-left {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.avatar-wrap {
  flex-shrink: 0;
}

.status-info {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.player-name {
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
}

.status-text {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}
</style>
