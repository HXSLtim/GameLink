<template>
  <view class="action-bar">
    <view class="select-all" @tap="$emit('toggle-all')">
      <view class="checkbox" :class="{ checked: allSelected }">
        <uv-icon v-if="allSelected" name="checkbox-mark" size="16" color="#fff"></uv-icon>
      </view>
      <text>全选</text>
    </view>
    <view class="action-info">
      <text>已选 {{ selectedCount }} 项</text>
    </view>
    <GlButton
      type="error"
      size="small"
      :disabled="selectedCount === 0"
      @click="$emit('delete')"
    >
      取消收藏
    </GlButton>
  </view>
</template>

<script setup lang="ts">
import GlButton from '@/components/gl/Button/index.vue'

interface Props {
  allSelected: boolean
  selectedCount: number
}

defineProps<Props>()

defineEmits<{
  'toggle-all': []
  delete: []
}>()
</script>

<style lang="scss" scoped>
.action-bar {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-sm) var(--spacing-md);
  padding-bottom: calc(var(--spacing-sm) + env(safe-area-inset-bottom));
  background: var(--color-bg);
  border-top: 1rpx solid var(--color-border);
}

.select-all {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  cursor: pointer;
  @include press-effect;

  text {
    font-size: var(--font-sm);
    color: var(--color-text);
  }
}

.checkbox {
  width: 32rpx;
  height: 32rpx;
  border-radius: var(--radius-full);
  border: 1rpx solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;

  &.checked {
    background: var(--color-primary);
    border-color: var(--color-primary);
  }
}

.action-info {
  flex: 1;
  text-align: center;
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}
</style>
