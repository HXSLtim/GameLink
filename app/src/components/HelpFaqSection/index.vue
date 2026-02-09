<template>
  <view class="faq-section">
    <SectionHeader title="常见问题" :show-more="false" />

    <view class="faq-list">
      <FaqItem
        v-for="faq in faqs"
        :key="faq.id"
        :question="faq.question"
        :answer="faq.answer"
        :expanded="expandedId === faq.id"
        @toggle="$emit('toggle', faq.id)"
      />

      <GlEmpty v-if="faqs.length === 0" title="未找到相关问题" compact />
    </view>
  </view>
</template>

<script setup lang="ts">
import SectionHeader from '@/components/SectionHeader/index.vue'
import GlEmpty from '@/components/gl/Empty/index.vue'
import FaqItem from '@/components/FaqItem/index.vue'
import type { FaqItemData } from '@/types/faq'

interface Props {
  faqs: FaqItemData[]
  expandedId?: number | null
}

defineProps<Props>()

defineEmits<{
  toggle: [id: number]
}>()
</script>

<style lang="scss" scoped>
.faq-section {
  padding: 0 24rpx;
}

.faq-list {
  margin-top: 16rpx;
}
</style>
