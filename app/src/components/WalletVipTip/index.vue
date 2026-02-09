<template>
  <view v-if="vipLevel > 0" class="vip-tip" @tap="$emit('click')">
    <GlTag type="warning" size="small">VIP{{ vipLevel }}</GlTag>
    <text class="vip-text">专属 {{ discountText }} 折扣</text>
    <text class="vip-more">查看特权 ›</text>
  </view>
</template>

<script setup lang="ts">
import GlTag from '@/components/gl/Tag/index.vue'

interface Props {
  vipLevel: number
  discountText: string | number
}

defineProps<Props>()

defineEmits<{
  click: []
}>()
</script>

<style lang="scss" scoped>
.vip-tip {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  margin: 0 var(--spacing-md) var(--spacing-sm);
  padding: var(--spacing-xs) var(--spacing-md);
  background: var(--color-bg-card);
  border: 1rpx solid var(--color-border);
  border-radius: var(--radius-md);
  cursor: pointer;
  @include press-effect;
}

.vip-text {
  flex: 1;
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}

.vip-more {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}
</style>
