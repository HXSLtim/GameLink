<template>
  <GlCard :title="title" :shadow="false" bordered>
    <template v-if="$slots.extra" #extra>
      <slot name="extra" />
    </template>
    <view class="form-list">
      <slot></slot>
    </view>
  </GlCard>
</template>

<script setup lang="ts">
import GlCard from '@/components/gl/Card/index.vue'

interface Props {
  title: string
}

defineProps<Props>()
</script>

<style lang="scss" scoped>
.form-list {
  display: flex;
  flex-direction: column;
}
</style>
