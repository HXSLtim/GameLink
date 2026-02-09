<template>
  <view
    class="status-info"
    :class="[`size-${size}`, `direction-${direction}`]"
    :style="infoStyle"
  >
    <view class="status-icon">
      <uv-icon v-if="iconType === 'uv-icon'" :name="icon" :size="iconSize" color="var(--color-text-secondary)" />
      <text v-else>{{ icon }}</text>
    </view>
    <view class="status-text">
      <text class="status-title">{{ title }}</text>
      <text v-if="description" class="status-desc">{{ description }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  icon: string
  title: string
  description?: string
  iconType?: 'emoji' | 'uv-icon'
  size?: 'sm' | 'md' | 'lg'
  direction?: 'row' | 'column'
  align?: 'start' | 'center' | 'end'
  gap?: string
}>(), {
  description: '',
  iconType: 'uv-icon',
  size: 'sm',
  direction: 'row',
  align: 'center',
  gap: 'var(--spacing-sm)',
})

const iconSize = computed(() => {
  const map = { sm: 20, md: 24, lg: 28 }
  return map[props.size] ?? 20
})

const infoStyle = computed(() => {
  const alignMap: Record<'start' | 'center' | 'end', string> = {
    start: 'flex-start',
    center: 'center',
    end: 'flex-end',
  }
  return {
    flexDirection: props.direction,
    alignItems: alignMap[props.align],
    gap: props.gap,
  }
})
</script>

<style lang="scss" scoped>
.status-info {
  display: flex;
  flex: 1;
  min-width: 0;
}

.status-icon {
  font-size: var(--font-lg);
  flex-shrink: 0;
}

.status-text {
  flex: 1;
  min-width: 0;
}

.status-title {
  display: block;
  font-weight: 600;
  margin-bottom: var(--spacing-xs);
}

.status-desc {
  color: var(--color-text-secondary);
}

.direction-column {
  .status-text {
    text-align: center;
  }
}

.size-sm {
  .status-icon { font-size: var(--font-lg); }
  .status-icon :deep(.uv-icon) { font-size: var(--font-lg); }
  .status-title { font-size: var(--font-sm); }
  .status-desc { font-size: var(--font-xs); }
}

.size-md {
  .status-icon { font-size: var(--font-lg); }
  .status-icon :deep(.uv-icon) { font-size: var(--font-lg); }
  .status-title { font-size: var(--font-md); }
  .status-desc { font-size: var(--font-sm); }
}

.size-lg {
  .status-icon { font-size: var(--font-xl); }
  .status-icon :deep(.uv-icon) { font-size: var(--font-xl); }
  .status-title { font-size: var(--font-md); }
  .status-desc { font-size: var(--font-sm); }
}
</style>
