<template>
  <view class="status-card" :class="statusClass" :style="cardStyle">
    <slot />
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  statusClass?: string
  padding?: string
  margin?: string
  radius?: string
  borderWidth?: string
  borderColor?: string
  direction?: string
  justify?: string
  align?: string
  gap?: string
}

const props = withDefaults(defineProps<Props>(), {
  padding: 'var(--spacing-md)',
  margin: 'var(--spacing-md)',
  radius: 'var(--radius-md)',
  direction: 'row',
  justify: 'flex-start',
  align: 'center',
  gap: 'var(--spacing-md)',
})

const cardStyle = computed(() => {
  const styles: Record<string, string> = {
    padding: props.padding,
    margin: props.margin,
    borderRadius: props.radius,
    flexDirection: props.direction,
    justifyContent: props.justify,
    alignItems: props.align,
    gap: props.gap,
  }

  if (props.borderWidth) {
    styles['--status-border-width'] = props.borderWidth
  }
  if (props.borderColor) {
    styles['--status-border-color'] = props.borderColor
  }

  return styles
})
</script>

<style lang="scss" scoped>
.status-card {
  display: flex;
  background: var(--color-bg-card);
  border: 1rpx solid var(--color-border);
  border-left-width: var(--status-border-width, 1rpx);
  border-left-color: var(--status-border-color, var(--color-border));
  color: var(--color-text);
  box-sizing: border-box;
}
</style>
