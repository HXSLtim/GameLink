<template>
  <uv-popup :show="show" mode="bottom" round="24" @close="closePanel">
    <view class="menu-panel">
      <view v-if="targetId" class="menu-item" @tap="$emit('view-profile', targetId)">
        <uv-icon name="account" size="18" color="var(--color-text-secondary)" />
        <text>查看资料</text>
      </view>
      <view class="menu-item" @tap="$emit('clear')">
        <uv-icon name="trash" size="18" color="var(--color-text-secondary)" />
        <text>清空聊天记录</text>
      </view>
      <view class="menu-item danger" @tap="$emit('block')">
        <uv-icon name="close-circle" size="18" color="var(--color-error)" />
        <text>拉黑用户</text>
      </view>
      <view class="menu-divider" />
      <view class="menu-item cancel" @tap="closePanel">
        <text>取消</text>
      </view>
    </view>
  </uv-popup>
</template>

<script setup lang="ts">
import type { ChatInfo } from '@/types/message'

interface Props {
  show: boolean
  targetId?: ChatInfo['targetId']
}

defineProps<Props>()

const emit = defineEmits<{
  'update:show': [value: boolean]
  'view-profile': [targetId: number]
  clear: []
  block: []
}>()

const closePanel = () => {
  emit('update:show', false)
}
</script>

<style lang="scss" scoped>
.menu-panel {
  padding: var(--spacing-sm) 0;
  padding-bottom: calc(var(--spacing-sm) + env(safe-area-inset-bottom));
}

.menu-item {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
  padding: 28rpx 48rpx;
  font-size: var(--font-md);
  color: var(--color-text);
  transition: background 0.15s ease;
  cursor: pointer;

  &:hover {
    background: var(--color-bg-secondary);
  }

  &:active {
    background: var(--color-bg-secondary);
  }

  &.danger {
    color: var(--color-error);
  }

  &.cancel {
    color: var(--color-text-secondary);
    font-size: var(--font-sm);
  }
}

.menu-divider {
  height: 1rpx;
  background: var(--color-border);
  margin: var(--spacing-xs) var(--spacing-lg);
}
</style>
