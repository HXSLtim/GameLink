<template>
  <view v-if="show" class="more-panel">
    <view class="more-grid">
      <view class="more-item" @tap="$emit('image')">
        <view class="more-icon">🖼️</view>
        <text>图片</text>
      </view>
      <view class="more-item" @tap="$emit('camera')">
        <view class="more-icon">📷</view>
        <text>拍照</text>
      </view>
      <view v-if="showOrder" class="more-item" @tap="$emit('order')">
        <view class="more-icon">📋</view>
        <text>查看订单</text>
      </view>
      <view class="more-item" @tap="$emit('report')">
        <view class="more-icon">⚠️</view>
        <text>举报</text>
      </view>
    </view>
    <view class="more-close" @tap="$emit('close')">
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
  border-radius: 32rpx 32rpx 0 0;
  padding: 32rpx;
  padding-bottom: calc(32rpx + env(safe-area-inset-bottom));
  z-index: 200;
  box-shadow: 0 -8rpx 32rpx rgba(0, 0, 0, 0.1);
}

.more-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 32rpx;
  margin-bottom: 32rpx;
}

.more-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12rpx;
}

.more-icon {
  width: 100rpx;
  height: 100rpx;
  background: var(--color-bg-secondary);
  border-radius: 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 40rpx;
}

.more-item text {
  font-size: 24rpx;
  color: var(--color-text-secondary);
}

.more-close {
  text-align: center;
  padding: 20rpx;
  font-size: 28rpx;
  color: var(--color-text-secondary);
}
</style>
