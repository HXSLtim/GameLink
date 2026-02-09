<template>
  <view class="header-stats" :class="[`size-${size}`, { 'header-stats--clickable': clickable }]">
    <template v-for="(item, index) in items" :key="item.key || item.label || index">
      <view
        class="header-stat"
        :class="{ 'header-stat--clickable': clickable || item.clickable }"
        :style="itemStyle"
        @tap="handleClick(item, index)"
      >
        <view class="header-stat__value">
          <slot name="value" :item="item">
            <text>{{ item.value }}</text>
          </slot>
        </view>
        <text class="header-stat__label">{{ item.label }}</text>
      </view>
      <view v-if="showDivider && index < items.length - 1" class="header-stat-divider"></view>
    </template>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { HeaderStatItem } from '@/types/ui'

const props = withDefaults(defineProps<{
  items: HeaderStatItem[]
  showDivider?: boolean
  clickable?: boolean
  size?: 'sm' | 'md' | 'lg'
  itemPadding?: string
}>(), {
  showDivider: true,
  clickable: false,
  size: 'md',
  itemPadding: 'var(--spacing-xs) 0',
})

const emit = defineEmits<{
  'item-click': [item: HeaderStatItem, index: number]
}>()

const itemStyle = computed(() => ({
  padding: props.itemPadding,
}))

const handleClick = (item: HeaderStatItem, index: number) => {
  if (!props.clickable && !item.clickable) return
  emit('item-click', item, index)
}
</script>

<style lang="scss" scoped>
.header-stats {
  display: flex;
  align-items: stretch;
}

.header-stat {
  flex: 1;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.header-stat--clickable {
  transition: all 0.2s;
  cursor: pointer;
  @include press-effect;

  &:active {
    transform: scale(0.95);
  }
}

.header-stat__value {
  display: block;
  font-weight: 700;
  color: var(--color-text);
}

.header-stat__value :deep(.price-tag) {
  color: var(--color-text);
}

.header-stat__label {
  font-weight: 500;
  color: var(--color-text-secondary);
  margin-top: var(--spacing-xs);
}

.header-stat-divider {
  width: 1rpx;
  background: var(--color-border);
  margin: var(--spacing-xs) 0;
}

.size-sm {
  .header-stat__value { font-size: var(--font-sm); }
  .header-stat__label { font-size: var(--font-xs); }
}

.size-md {
  .header-stat__value { font-size: var(--font-md); }
  .header-stat__label { font-size: var(--font-xs); }
}

.size-lg {
  .header-stat__value { font-size: var(--font-lg); }
  .header-stat__label { font-size: var(--font-xs); }
}
</style>
