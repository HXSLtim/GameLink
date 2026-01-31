<template>
  <view class="earnings-page page-container">
    <!-- 顶部导航 -->
    <NavBar title="我的收益" @back="goBack" />

    <scroll-view class="content-scroll" scroll-y>
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
      <GlCard title="收益明细" :shadow="false" bordered class="detail-card">
        <template #extra>
          <view class="filter-btn" @tap="showFilter = true">
            <uv-icon name="list" size="16" color="var(--color-text-secondary)"></uv-icon>
            <text>筛选</text>
          </view>
        </template>
        
        <InfiniteList
          :state="loading ? 'loading' : (earningsList.length === 0 ? 'empty' : 'content')"
          :loading="loadingMore"
          :no-more="noMore"
          :show-load-more="earningsList.length > 0"
          empty-title="暂无收益记录"
          padding="0"
          @load-more="loadMore"
        >
          <EarningsItem
            v-for="item in earningsList"
            :key="item.id"
            :type="item.type"
            :title="item.title"
            :description="item.description"
            :amount="item.amount"
            :created-at="item.createdAt"
            @click="goToOrder(item.orderId)"
          />
        </InfiniteList>
      </GlCard>

      <view class="bottom-placeholder"></view>
    </scroll-view>

    <!-- 筛选弹窗 -->
    <uv-popup :show="showFilter" mode="bottom" round="20" @close="showFilter = false">
      <view class="filter-panel">
        <view class="filter-header">
          <text class="filter-title">筛选</text>
          <text class="filter-reset" @tap="resetFilter">重置</text>
        </view>
        
        <view class="filter-section">
          <text class="section-label">收益类型</text>
          <view class="filter-options">
            <view 
              v-for="type in typeOptions" 
              :key="type.value"
              class="filter-option"
              :class="{ active: filterType === type.value }"
              @tap="filterType = type.value"
            >
              <text>{{ type.label }}</text>
            </view>
          </view>
        </view>
        
        <view class="filter-footer">
          <GlButton type="default" plain @click="showFilter = false">取消</GlButton>
          <GlButton type="primary" @click="applyFilter">确定</GlButton>
        </view>
      </view>
    </uv-popup>

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import TabsBar from '@/components/TabsBar/index.vue'
import GlCard from '@/components/gl/Card/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
import InfiniteList from '@/components/InfiniteList/index.vue'
// Business 组件
import EarningsSummaryCard from '@/components/EarningsSummaryCard/index.vue'
import EarningsChart from '@/components/EarningsChart/index.vue'
import EarningsItem from '@/components/EarningsItem/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
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

onMounted(init)
onShow(init)
</script>

<style lang="scss" scoped>
.earnings-page {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  box-sizing: border-box;
  
  @include desktop {
    height: 100vh;
    min-height: auto;
    overflow: hidden;
  }
}

.content-scroll {
  flex: 1;
  overflow-y: auto;
}

.chart-wrap {
  padding: 0 24rpx;
  margin-bottom: 20rpx;
}

.detail-card {
  margin: 0 24rpx;
}

.filter-btn {
  display: flex;
  align-items: center;
  gap: 8rpx;
  font-size: 26rpx;
  color: var(--color-text-secondary);
}

.filter-panel {
  padding: 24rpx;
  padding-bottom: calc(24rpx + env(safe-area-inset-bottom));
}

.filter-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 32rpx;
}

.filter-title {
  font-size: 34rpx;
  font-weight: 700;
  color: var(--color-text);
}

.filter-reset {
  font-size: 28rpx;
  color: var(--color-primary);
}

.filter-section {
  margin-bottom: 32rpx;
}

.section-label {
  display: block;
  font-size: 28rpx;
  color: var(--color-text);
  margin-bottom: 16rpx;
}

.filter-options {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.filter-option {
  padding: 16rpx 32rpx;
  background: var(--color-bg-secondary);
  border-radius: 32rpx;
  border: 2rpx solid var(--color-border);
  font-size: 26rpx;
  color: var(--color-text);
  
  &.active {
    background: rgba(0, 210, 106, 0.1);
    border-color: var(--color-primary);
    color: var(--color-primary);
  }
}

.filter-footer {
  display: flex;
  gap: 24rpx;
  padding-top: 24rpx;
}

.bottom-placeholder {
  height: 100rpx;
}
</style>
