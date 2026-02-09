<template>
  <GlCard
    class="section-card"
    :title="title"
    :extra="extra"
    :icon="icon"
    :shadow="shadow"
    :bordered="bordered"
    :padding="padding"
    :custom-style="resolvedStyle"
  >
    <template v-if="$slots.header" #header>
      <slot name="header" />
    </template>
    <template v-if="$slots.extra" #extra>
      <slot name="extra" />
    </template>
    <template v-if="$slots.footer" #footer>
      <slot name="footer" />
    </template>
    <slot />
  </GlCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlCard from '@/components/gl/Card/index.vue'

interface Props {
  title?: string
  extra?: string
  icon?: string
  shadow?: boolean
  bordered?: boolean
  padding?: string
  margin?: string
  customStyle?: string | Record<string, string>
}

const props = withDefaults(defineProps<Props>(), {
  shadow: false,
  bordered: true,
  padding: '24rpx',
  margin: '0 var(--spacing-md) var(--spacing-sm)',
})

const resolvedStyle = computed(() => {
  const baseStyle = { margin: props.margin }
  if (!props.customStyle) return baseStyle
  if (typeof props.customStyle === 'string') {
    return `margin: ${props.margin}; ${props.customStyle}`
  }
  return { ...baseStyle, ...props.customStyle }
})
</script>
