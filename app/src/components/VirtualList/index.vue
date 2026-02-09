<template>
  <scroll-view
    class="virtual-list"
    scroll-y
    :scroll-top="scrollTop"
    :scroll-into-view="scrollIntoView"
    :scroll-with-animation="scrollWithAnimation"
    :enable-back-to-top="enableBackToTop"
    :refresher-enabled="refresherEnabled"
    :refresher-triggered="refresherTriggered"
    @scroll="handleScroll"
    @scrolltolower="handleScrollToLower"
    @refresherrefresh="handleRefresh"
  >
    <!-- 顶部占位 -->
    <view class="top-placeholder" :style="{ height: `${topHeight}px` }"></view>
    
    <!-- 可见列表项 -->
    <view
      v-for="item in visibleItems"
      :key="getItemKey(item)"
      class="list-item"
      :style="{ height: itemHeight ? `${itemHeight}px` : 'auto' }"
    >
      <slot name="item" :item="item.data" :index="item.index"></slot>
    </view>
    
    <!-- 底部占位 -->
    <view class="bottom-placeholder" :style="{ height: `${bottomHeight}px` }"></view>
    
    <!-- 加载更多 -->
    <view v-if="showLoadMore" class="load-more">
      <slot name="loading" v-if="loading">
        <view class="loading-default">
          <text>加载中...</text>
        </view>
      </slot>
      <slot name="noMore" v-else-if="noMore">
        <view class="no-more-default">
          <text>没有更多了</text>
        </view>
      </slot>
    </view>
  </scroll-view>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'

interface Props {
  // 数据源
  list: any[]
  // 每项高度（固定高度时必填）
  itemHeight?: number
  // 预估高度（动态高度时使用）
  estimatedHeight?: number
  // 缓冲区大小（渲染可视区域外的额外项数）
  buffer?: number
  // 唯一键字段名
  keyField?: string
  // 容器高度
  height?: string | number
  // 是否启用下拉刷新
  refresherEnabled?: boolean
  // 是否正在刷新
  refresherTriggered?: boolean
  // 是否正在加载
  loading?: boolean
  // 是否没有更多
  noMore?: boolean
  // 是否显示加载更多
  showLoadMore?: boolean
  // 滚动动画
  scrollWithAnimation?: boolean
  // 点击顶部状态栏回到顶部
  enableBackToTop?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  itemHeight: 0,
  estimatedHeight: 100,
  buffer: 5,
  keyField: 'id',
  height: '100vh',
  refresherEnabled: false,
  refresherTriggered: false,
  loading: false,
  noMore: false,
  showLoadMore: true,
  scrollWithAnimation: false,
  enableBackToTop: true,
})

const emit = defineEmits<{
  scroll: [e: any]
  scrollToLower: []
  refresh: []
}>()

// 状态
const scrollTop = ref(0)
const scrollIntoView = ref('')
const containerHeight = ref(0)
const currentScrollTop = ref(0)

// 实际使用的项高度
const realItemHeight = computed(() => props.itemHeight || props.estimatedHeight)

// 可见项数量
const visibleCount = computed(() => {
  const height = typeof props.height === 'number' ? props.height : parseInt(props.height) || 600
  return Math.ceil(height / realItemHeight.value) + props.buffer * 2
})

// 起始索引
const startIndex = computed(() => {
  const index = Math.floor(currentScrollTop.value / realItemHeight.value) - props.buffer
  return Math.max(0, index)
})

// 结束索引
const endIndex = computed(() => {
  return Math.min(props.list.length, startIndex.value + visibleCount.value)
})

// 可见项
const visibleItems = computed(() => {
  return props.list.slice(startIndex.value, endIndex.value).map((item, idx) => ({
    data: item,
    index: startIndex.value + idx,
  }))
})

// 顶部占位高度
const topHeight = computed(() => startIndex.value * realItemHeight.value)

// 底部占位高度
const bottomHeight = computed(() => {
  const remainingItems = props.list.length - endIndex.value
  return Math.max(0, remainingItems * realItemHeight.value)
})

// 获取项的唯一键
const getItemKey = (item: { data: any; index: number }) => {
  if (props.keyField && item.data[props.keyField] !== undefined) {
    return item.data[props.keyField]
  }
  return item.index
}

// 滚动事件
const handleScroll = (e: any) => {
  currentScrollTop.value = e.detail.scrollTop
  emit('scroll', e)
}

// 滚动到底部
const handleScrollToLower = () => {
  emit('scrollToLower')
}

// 下拉刷新
const handleRefresh = () => {
  emit('refresh')
}

// 滚动到指定位置
const scrollTo = (top: number) => {
  scrollTop.value = top
}

// 滚动到指定索引
const scrollToIndex = (index: number) => {
  scrollTop.value = index * realItemHeight.value
}

// 暴露方法
defineExpose({
  scrollTo,
  scrollToIndex,
})
</script>

<style lang="scss" scoped>
.virtual-list {
  width: 100%;
  height: v-bind('typeof props.height === "number" ? props.height + "px" : props.height');
  overflow: hidden;
  
  .list-item {
    width: 100%;
  }
  
  .load-more {
    padding: var(--spacing-md) 0;
    
    .loading-default,
    .no-more-default {
      display: flex;
      justify-content: center;
      
      text {
        font-size: var(--font-sm);
        color: var(--color-text-placeholder);
      }
    }
  }
}
</style>
