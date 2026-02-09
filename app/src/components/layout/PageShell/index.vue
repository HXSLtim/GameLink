<template>
  <view class="page-shell page-container" :class="{ 'page-shell--with-header': hasHeader }">
    <slot name="header"></slot>

    <scroll-view
      v-if="scroll"
      class="page-shell__scroll"
      scroll-y
      :lower-threshold="lowerThreshold"
      :refresher-enabled="refresherEnabled"
      :refresher-triggered="refresherTriggered"
      :refresher-threshold="refresherThreshold"
      :refresher-background="refresherBackground"
      @scrolltolower="$emit('scrolltolower')"
      @refresherrefresh="$emit('refresherrefresh')"
    >
      <view class="page-shell__content" :class="contentClass" :style="contentStyle">
        <slot></slot>
      </view>
    </scroll-view>

    <view v-else class="page-shell__content page-shell__content--static" :class="contentClass" :style="contentStyle">
      <slot></slot>
    </view>

    <slot name="footer"></slot>

    <ConfirmDialog
      :show="dialogState.show"
      :title="dialogState.title"
      :content="dialogState.content"
      :confirm-text="dialogState.confirmText"
      :cancel-text="dialogState.cancelText"
      @confirm="confirm"
      @cancel="cancel"
      @close="close"
    />
  </view>
</template>

<script setup lang="ts">
import { computed, useSlots } from 'vue'
import { useDevice } from '@/composables/useDevice'
import ConfirmDialog from '@/components/ConfirmDialog/index.vue'
import { useConfirmDialog } from '@/composables/useConfirmDialog'

interface Props {
  scroll?: boolean
  padding?: string
  lowerThreshold?: number
  contentClass?: string | string[] | Record<string, boolean>
  refresherEnabled?: boolean
  refresherTriggered?: boolean
  refresherThreshold?: number
  refresherBackground?: string
}

const props = withDefaults(defineProps<Props>(), {
  scroll: true,
  padding: '24rpx',
  lowerThreshold: 120,
  contentClass: '',
  refresherEnabled: false,
  refresherTriggered: false,
  refresherThreshold: 45,
  refresherBackground: 'transparent',
})

defineEmits<{
  scrolltolower: []
  refresherrefresh: []
}>()

const { isPC } = useDevice()

const contentStyle = computed(() => {
  const padding = props.padding
  const normalized = padding.trim().split(/\s+/)
  const isZeroPadding = normalized.every(part => part === '0' || part === '0px' || part === '0rpx')
  const style: Record<string, string> = { padding }

  if (hasHeader.value && isPC.value && isZeroPadding) {
    style.paddingTop = `calc(${padding} + var(--spacing-sm))`
  }

  return style
})

const slots = useSlots()
const hasHeader = computed(() => Boolean(slots.header))

const { dialogState, confirm, cancel, close } = useConfirmDialog()
</script>

<style lang="scss" scoped>
.page-shell {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  box-sizing: border-box;
  overflow: hidden;

  @include desktop {
    height: 100vh;
    min-height: auto;
  }
}

.page-shell__scroll {
  flex: 1;
  min-height: 0;
  overflow-y: auto;

  // PC 端：更精致的滚动条
  @include desktop {
    @include custom-scrollbar(6px);
    scroll-behavior: smooth;
  }
}

.page-shell__content {
  box-sizing: border-box;

  // 大屏 PC：限制内容最大宽度并居中，提升可读性
  @include desktop-lg {
    max-width: 1400px;
    margin-left: auto;
    margin-right: auto;
  }
}

.page-shell__content--static {
  flex: 1;
  min-height: 0;
  height: 0; // 配合 flex: 1 确保高度正确计算，让内部 scroll-view 可以滚动
  display: flex;
  flex-direction: column;

  // 大屏 PC：限制内容最大宽度并居中
  @include desktop-lg {
    max-width: 1400px;
    margin-left: auto;
    margin-right: auto;
    width: 100%;
  }
}
</style>
