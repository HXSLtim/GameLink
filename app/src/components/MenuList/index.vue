<template>
  <view class="menu-list">
    <view
      v-for="item in items"
      :key="item.key"
      class="menu-item"
      @tap="handleClick(item)"
    >
      <view class="menu-icon" :style="{ background: item.iconBg || 'var(--color-bg-secondary)' }">
        <uv-icon :name="item.icon" size="22" :color="item.iconColor || 'var(--color-text-secondary)'"></uv-icon>
      </view>
      <text class="menu-text">{{ item.label }}</text>
      
      <!-- 右侧内容插槽 -->
      <slot :name="item.key">
        <view v-if="item.badge" class="menu-badge">{{ item.badge }}</view>
        <view v-else-if="item.value" class="menu-value">{{ item.value }}</view>
        <uv-icon v-else name="arrow-right" size="16" color="var(--color-text-secondary)"></uv-icon>
      </slot>
    </view>
  </view>
</template>

<script setup lang="ts">
import type { MenuItem } from '@/types/ui'

interface Props {
  items: MenuItem[]
}

defineProps<Props>()

const emit = defineEmits<{
  click: [item: MenuItem]
}>()

const handleClick = (item: MenuItem) => {
  if (item.disabled) return
  emit('click', item)
}
</script>

<style lang="scss" scoped>
.menu-list {
  background: var(--color-bg-card);
  margin: var(--spacing-sm) var(--spacing-md);
  border-radius: var(--radius-md);
  overflow: hidden;
  border: 1rpx solid var(--color-border);
}

.menu-item {
  display: flex;
  align-items: center;
  padding: var(--spacing-sm) var(--spacing-md);
  border-bottom: 1rpx solid var(--color-border);
  transition: background 0.2s ease, padding-left 0.2s ease;
  cursor: pointer;
  @include press-effect;
  
  &:last-child {
    border-bottom: none;
  }
  
  &:hover {
    background: var(--color-bg-secondary);
    padding-left: calc(var(--spacing-md) + 4rpx);
  }
  
  &:active {
    background: var(--color-bg-secondary);
  }
}

.menu-icon {
  width: 44rpx;
  height: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: var(--spacing-md);
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
}

.menu-text {
  flex: 1;
  font-size: var(--font-md);
  font-weight: 500;
  color: var(--color-text);
}

.menu-badge {
  min-width: 28rpx;
  height: 28rpx;
  padding: 0 var(--spacing-xs);
  background: var(--color-error);
  border-radius: var(--radius-full);
  font-size: var(--font-xs);
  font-weight: 600;
  color: #FFFFFF;
  text-align: center;
  line-height: 28rpx;
}

.menu-value {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  margin-right: var(--spacing-xs);
}
</style>
