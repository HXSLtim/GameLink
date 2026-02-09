<template>
  <view
    class="record-list-item"
    :class="{
      'record-list-item--pressable': pressable,
      'record-list-item--divider': divider,
    }"
    :style="itemStyle"
    @tap="$emit('click')"
  >
    <view class="record-list-item__icon">
      <slot name="icon"></slot>
    </view>
    <view class="record-list-item__content">
      <slot name="content"></slot>
    </view>
    <view class="record-list-item__amount">
      <slot name="amount"></slot>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  gap?: string
  padding?: string
  divider?: boolean
  pressable?: boolean
}>(), {
  gap: 'var(--spacing-md)',
  padding: 'var(--spacing-md) 0',
  divider: true,
  pressable: false,
})

defineEmits<{
  click: []
}>()

const itemStyle = computed(() => ({
  gap: props.gap,
  padding: props.padding,
}))
</script>

<style lang="scss" scoped>
.record-list-item {
  display: flex;
  align-items: center;
}

.record-list-item--divider {
  border-bottom: 1rpx solid var(--color-border);

  &:last-child {
    border-bottom: none;
  }
}

.record-list-item--pressable {
  transition: background 0.2s ease;
  cursor: pointer;
  @include press-effect;

  &:active {
    background: var(--color-bg-secondary);
  }
}

.record-list-item__icon {
  flex-shrink: 0;
}

.record-list-item__content {
  flex: 1;
  min-width: 0;
}

.record-list-item__amount {
  flex-shrink: 0;
}
</style>
