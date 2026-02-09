<template>
  <view class="footer">
    <text class="footer-text">
      {{ prefixText }}
      <text class="link-text" @tap="$emit('agreement', 'user')">《用户协议》</text>
      和
      <text class="link-text" @tap="$emit('agreement', 'privacy')">《隐私政策》</text>
    </text>
  </view>
</template>

<script setup lang="ts">
import type { AgreementType } from '@/types/agreement'

type AgreementLinkType = Extract<AgreementType, 'user' | 'privacy'>
interface Props {
  prefixText?: string
}

withDefaults(defineProps<Props>(), {
  prefixText: '登录即代表同意',
})

defineEmits<{
  agreement: [type: AgreementLinkType]
}>()
</script>

<style lang="scss" scoped>
.footer {
  padding: var(--spacing-md);
  text-align: center;
}

.footer-text {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}

.link-text {
  color: var(--color-primary);
}
</style>
