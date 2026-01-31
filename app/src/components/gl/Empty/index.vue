<template>
  <view class="gl-empty" :style="customStyle">
    <view class="gl-empty__icon">
      <slot name="icon">
        <uv-icon :name="icon" :size="iconSize" color="var(--color-text-placeholder)" />
      </slot>
    </view>
    
    <text v-if="title" class="gl-empty__title">{{ title }}</text>
    <text v-if="description" class="gl-empty__desc">{{ description }}</text>
    
    <view v-if="$slots.default" class="gl-empty__action">
      <slot></slot>
    </view>
  </view>
</template>

<script setup lang="ts">
interface Props {
  icon?: string
  iconSize?: number
  title?: string
  description?: string
  customStyle?: string | Record<string, string>
}

withDefaults(defineProps<Props>(), {
  icon: 'empty-data',
  iconSize: 120,
  title: '暂无数据',
  description: '',
})
</script>

<style lang="scss" scoped>
.gl-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80rpx 40rpx;
}

.gl-empty__icon {
  margin-bottom: 24rpx;
  opacity: 0.6;
}

.gl-empty__title {
  font-size: 32rpx;
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-bottom: 12rpx;
}

.gl-empty__desc {
  font-size: 26rpx;
  color: var(--color-text-placeholder);
  text-align: center;
  max-width: 400rpx;
}

.gl-empty__action {
  margin-top: 32rpx;
}
</style>
