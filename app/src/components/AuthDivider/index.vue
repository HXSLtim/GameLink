<template>
  <view class="divider">
    <view class="divider-line"></view>
    <text class="divider-text">{{ text }}</text>
    <view class="divider-line"></view>
  </view>
</template>

<script setup lang="ts">
interface Props {
  text?: string
}

withDefaults(defineProps<Props>(), {
  text: '或',
})
</script>

<style lang="scss" scoped>
.divider {
  display: flex;
  align-items: center;
  gap: 24rpx;
  margin: 48rpx 0;
}

.divider-line {
  flex: 1;
  height: 1rpx;
  background: var(--color-border);
}

.divider-text {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}
</style>
