<template>
  <SectionCard title="我的评价" margin="var(--spacing-sm) var(--spacing-md)">
    <RatingStars :rating="rating" size="small" />
    <text class="review-text">{{ content }}</text>
  </SectionCard>
</template>

<script setup lang="ts">
import SectionCard from '@/components/SectionCard/index.vue'
import RatingStars from '@/components/RatingStars/index.vue'

interface Props {
  rating: number
  content: string
}

defineProps<Props>()
</script>

<style lang="scss" scoped>
.review-text {
  display: block;
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  margin-top: var(--spacing-xs);
  line-height: 1.5;
}
</style>
