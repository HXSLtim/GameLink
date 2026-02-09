<template>
  <SectionCard :title="title">
    <view class="tips-list">
      <text v-for="(item, index) in items" :key="index" class="tips-item">
        {{ item }}
      </text>
    </view>
  </SectionCard>
</template>

<script setup lang="ts">
import SectionCard from '@/components/SectionCard/index.vue'

interface Props {
  title: string
  items: string[]
}

defineProps<Props>()
</script>

<style lang="scss" scoped>
.tips-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.tips-item {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  line-height: 1.5;
}
</style>
