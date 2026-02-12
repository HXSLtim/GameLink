<template>
  <PageShell
    class="base-page-layout"
    :class="{ 'with-sidebar': showTabBar }"
    :scroll="scroll"
    :padding="padding"
    :lower-threshold="lowerThreshold"
    :refresher-enabled="refresherEnabled"
    :refresher-triggered="refresherTriggered"
    :refresher-threshold="refresherThreshold"
    :refresher-background="refresherBackground"
    @scrolltolower="$emit('scrolltolower')"
    @refresherrefresh="$emit('refresherrefresh')"
  >
    <template #header>
      <slot name="nav">
        <NavBar :title="title" :show-back="showBack" :title-align="titleAlign" />
      </slot>

      <view v-if="showHeader && hasHeaderContent" class="base-page-layout__header">
        <ListHeaderStack class="base-page-layout__header-stack">
          <template #search>
            <slot name="search"></slot>
          </template>
          <template #banner>
            <slot name="banner"></slot>
          </template>
          <template #tabs>
            <slot name="tabs"></slot>
          </template>
          <template #sort>
            <slot name="sort"></slot>
          </template>
          <slot name="header-extra"></slot>
        </ListHeaderStack>
      </view>
    </template>

    <slot></slot>

    <template #footer>
      <slot name="footer"></slot>
      <CustomTabBar
        v-if="showTabBar"
        :current="tabBarCurrent"
        :show-mobile-tab-bar="showMobileTabBar"
      />
    </template>
  </PageShell>
</template>

<script setup lang="ts">
import { computed, useSlots } from 'vue'
import PageShell from '@/components/layout/PageShell/index.vue'
import ListHeaderStack from '@/components/layout/ListHeaderStack/index.vue'
import NavBar from '@/components/NavBar/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'

interface Props {
  title?: string
  showBack?: boolean
  titleAlign?: 'left' | 'center'
  scroll?: boolean
  padding?: string
  lowerThreshold?: number
  refresherEnabled?: boolean
  refresherTriggered?: boolean
  refresherThreshold?: number
  refresherBackground?: string
  showHeader?: boolean
  showTabBar?: boolean
  showMobileTabBar?: boolean
  tabBarCurrent?: number
}

const props = withDefaults(defineProps<Props>(), {
  title: '',
  showBack: true,
  titleAlign: 'left',
  scroll: true,
  padding: '24rpx',
  lowerThreshold: 120,
  refresherEnabled: false,
  refresherTriggered: false,
  refresherThreshold: 45,
  refresherBackground: 'transparent',
  showHeader: true,
  showTabBar: false,
  showMobileTabBar: true,
  tabBarCurrent: 0,
})

defineEmits<{
  scrolltolower: []
  refresherrefresh: []
}>()

const slots = useSlots()
const hasHeaderContent = computed(() => {
  return Boolean(
    slots.search ||
      slots.banner ||
      slots.tabs ||
      slots.sort ||
      slots['header-extra']
  )
})
</script>

<style lang="scss" scoped>
.base-page-layout__header {
  padding: 0;
  background: var(--color-bg);
}

.base-page-layout__header-stack {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}
</style>
