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
export interface MenuItem {
  key: string
  label: string
  icon: string
  iconColor?: string
  iconBg?: string
  badge?: number | string
  value?: string
  disabled?: boolean
}

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
  margin: 20rpx 28rpx;
  border-radius: 20rpx;
  overflow: hidden;
  border: 2rpx solid var(--color-border);
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 28rpx 24rpx;
  border-bottom: 1rpx solid var(--color-border);
  transition: all 0.2s;
  
  &:last-child {
    border-bottom: none;
  }
  
  &:active {
    background: var(--color-bg-secondary);
  }
}

.menu-icon {
  width: 52rpx;
  height: 52rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 24rpx;
  border-radius: 16rpx;
}

.menu-text {
  flex: 1;
  font-size: 30rpx;
  font-weight: 500;
  color: var(--color-text);
}

.menu-badge {
  min-width: 36rpx;
  height: 36rpx;
  padding: 0 10rpx;
  background: var(--color-error);
  border-radius: 18rpx;
  font-size: 22rpx;
  font-weight: 600;
  color: #FFFFFF;
  text-align: center;
  line-height: 36rpx;
}

.menu-value {
  font-size: 26rpx;
  color: var(--color-text-secondary);
  margin-right: 8rpx;
}
</style>
