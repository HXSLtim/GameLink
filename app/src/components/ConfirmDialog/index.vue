<template>
  <uv-popup :show="show" mode="center" round="16" @close="$emit('close')">
    <view class="confirm-dialog">
      <text class="confirm-title">{{ title }}</text>
      <text class="confirm-content">{{ content }}</text>
      <view class="confirm-actions">
        <GlButton class="action-btn" type="default" plain size="medium" @click="$emit('cancel')">
          {{ cancelText }}
        </GlButton>
        <GlButton class="action-btn" type="primary" size="medium" @click="$emit('confirm')">
          {{ confirmText }}
        </GlButton>
      </view>
    </view>
  </uv-popup>
</template>

<script setup lang="ts">
import GlButton from '@/components/gl/Button/index.vue'
import type { ConfirmOptions } from '@/types/ui'

type Props = ConfirmOptions & {
  show: boolean
}

withDefaults(defineProps<Props>(), {
  title: '提示',
  content: '',
  confirmText: '确定',
  cancelText: '取消',
})

defineEmits<{
  confirm: []
  cancel: []
  close: []
}>()
</script>

<style lang="scss" scoped>
.confirm-dialog {
  width: 560rpx;
  padding: var(--spacing-lg);
  background: var(--color-bg-card);
  border: 1rpx solid var(--color-border);
  border-radius: var(--radius-md);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  box-sizing: border-box;
}

.confirm-title {
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
}

.confirm-content {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  line-height: 1.6;
}

.confirm-actions {
  display: flex;
  gap: var(--spacing-sm);
}

.action-btn {
  flex: 1;
}
</style>
