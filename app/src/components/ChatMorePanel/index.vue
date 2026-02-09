<template>
  <view v-if="show" class="more-panel">
    <view class="more-grid">
      <view class="more-item" @tap="$emit('image')">
        <view class="more-icon">
          <uv-icon name="photo" size="24" color="var(--color-primary)" />
        </view>
        <text>相册</text>
      </view>
      <view class="more-item" @tap="$emit('camera')">
        <view class="more-icon">
          <uv-icon name="camera" size="24" color="var(--color-primary)" />
        </view>
        <text>拍照</text>
      </view>
      <view v-if="showOrder" class="more-item" @tap="$emit('order')">
        <view class="more-icon">
          <uv-icon name="order" size="24" color="var(--color-primary)" />
        </view>
        <text>查看订单</text>
      </view>
      <view class="more-item" @tap="$emit('report')">
        <view class="more-icon more-icon--danger">
          <uv-icon name="warning" size="24" color="var(--color-error)" />
        </view>
        <text class="more-text--danger">举报</text>
      </view>
    </view>
    <view class="more-close" @tap="$emit('close')">
      <uv-icon name="arrow-down" size="16" color="var(--color-text-secondary)" />
      <text>收起</text>
    </view>
  </view>
</template>

<script setup lang="ts">
interface Props {
  show: boolean
  showOrder?: boolean
}

withDefaults(defineProps<Props>(), {
  showOrder: false,
})

defineEmits<{
  close: []
  image: []
  camera: []
  order: []
  report: []
}>()
</script>

<style lang="scss" scoped>
.more-panel {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: var(--color-bg-card);
  border-radius: var(--radius-lg) var(--radius-lg) 0 0;
  padding: var(--spacing-lg) var(--spacing-md) var(--spacing-md);
  padding-bottom: calc(var(--spacing-md) + env(safe-area-inset-bottom));
  border-top: 1rpx solid var(--color-border);
  z-index: 200;
  box-shadow: 0 -4rpx 24rpx rgba(0, 0, 0, 0.08);
}

.more-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--spacing-lg) var(--spacing-md);
  margin-bottom: var(--spacing-lg);
}

.more-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-xs);
  cursor: pointer;
  transition: transform 0.2s;
  
  &:active {
    transform: scale(0.92);
  }
}

.more-icon {
  width: 96rpx;
  height: 96rpx;
  background: rgba(var(--color-primary-rgb), 0.08);
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;

  &:hover {
    background: rgba(var(--color-primary-rgb), 0.15);
    transform: translateY(-2rpx);
  }

  &--danger {
    background: rgba(var(--color-error-rgb, 239, 68, 68), 0.08);

    &:hover {
      background: rgba(var(--color-error-rgb, 239, 68, 68), 0.15);
    }
  }
}

.more-item text {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}

.more-text--danger {
  color: var(--color-error) !important;
}

.more-close {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-xs);
  padding: var(--spacing-sm);
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  cursor: pointer;
  border-radius: var(--radius-sm);
  transition: background 0.2s;

  &:hover {
    background: var(--color-bg-secondary);
  }

  &:active {
    transform: scale(0.97);
  }
}
</style>
