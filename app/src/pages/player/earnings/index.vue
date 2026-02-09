<template>
  <BasePageLayout
    class="earnings-page"
    padding="0"
    title="我的收益"
    :show-back="true"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
  >

    <!-- 收益总览 -->
    <EarningsSummaryCard
      :withdrawable="summary.withdrawable"
      :pending="summary.pending"
      :withdrawn="summary.withdrawn"
      :total="summary.total"
      @withdraw="goToWithdraw"
    />

    <!-- 统计周期切换 -->
    <TabsBar
      v-model="currentPeriod"
      :tabs="periodTabs"
      @change="switchPeriod"
    />

    <!-- 收益统计图表 -->
    <view class="chart-wrap">
      <EarningsChart
        :title="getPeriodTitle()"
        :total="periodTotal"
        :data="chartData"
      />
    </view>

    <!-- 收益明细 -->
    <EarningsDetailSection
      :loading="loading"
      :loading-more="loadingMore"
      :no-more="noMore"
      :items="earningsList"
      @filter="showFilter = true"
      @load-more="loadMore"
      @item-click="(item) => goToOrder(item.orderId)"
    />

    <view class="bottom-placeholder"></view>

    <template #footer>
      <!-- 筛选弹窗 -->
      <FilterPanel
        v-model:visible="showFilter"
        v-model="filterValues"
        title="筛选"
        :sections="filterSections"
        @reset="resetFilter"
        @apply="applyFilter"
      />
    </template>
  </BasePageLayout>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
import TabsBar from '@/components/TabsBar/index.vue'
import FilterPanel from '@/components/FilterPanel/index.vue'
// Business 组件
import EarningsSummaryCard from '@/components/EarningsSummaryCard/index.vue'
import EarningsChart from '@/components/EarningsChart/index.vue'
import EarningsDetailSection from '@/components/EarningsDetailSection/index.vue'
// Composables
import { usePlayerEarnings } from '@/composables/usePlayerEarnings'

const {
  loading,
  loadingMore,
  noMore,
  showFilter,
  summary,
  currentPeriod,
  periodTabs,
  chartData,
  periodTotal,
  earningsList,
  filterType,
  typeOptions,
  getPeriodTitle,
  switchPeriod,
  loadMore,
  applyFilter,
  resetFilter,
  goBack,
  goToWithdraw,
  goToOrder,
  init,
} = usePlayerEarnings()

const filterSections = computed(() => [
  { key: 'type', label: '收益类型', options: typeOptions },
])

const filterValues = computed({
  get: () => ({ type: filterType.value }),
  set: (values) => {
    filterType.value = String(values.type || 'all')
  },
})

onMounted(init)
onShow(init)
</script>

<style lang="scss" scoped>
.chart-wrap {
  padding: 0 24rpx;
  margin-bottom: 20rpx;
}

.bottom-placeholder {
  height: 100rpx;
}
</style>
