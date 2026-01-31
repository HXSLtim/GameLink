<template>
  <view class="wallet-page page-container">
    <!-- 顶部导航 -->
    <NavBar title="我的钱包" @back="goBack" />

    <!-- 余额卡片 -->
    <WalletBalanceCard
      :balance="wallet.balance"
      :total-recharge="wallet.totalRecharge"
      :total-spent="wallet.totalSpent"
      :show-balance="showBalance"
      @toggle-visibility="showBalance = !showBalance"
      @recharge="goToRecharge"
      @withdraw="goToWithdraw"
    />

    <!-- VIP 优惠提示 -->
    <view v-if="wallet.vipLevel > 0" class="vip-tip" @tap="goToVip">
      <GlTag type="warning" size="small">VIP{{ wallet.vipLevel }}</GlTag>
      <text class="vip-text">专属 {{ vipDiscountText }} 折扣</text>
      <text class="vip-more">查看特权 ›</text>
    </view>

    <!-- 快捷入口 -->
    <QuickActions :items="quickActions" @click="handleQuickAction" />

    <!-- 交易记录 -->
    <GlCard title="交易记录" :shadow="false" bordered class="records-card">
      <template #extra>
        <view class="filter-tabs">
          <text 
            v-for="tab in filterTabs" 
            :key="tab.key"
            class="filter-tab"
            :class="{ active: currentFilter === tab.key }"
            @tap="currentFilter = tab.key"
          >
            {{ tab.label }}
          </text>
        </view>
      </template>
      
      <InfiniteList
        :state="loading ? 'loading' : (filteredRecords.length === 0 ? 'empty' : 'content')"
        :loading="loadingMore"
        :no-more="noMore"
        :show-load-more="filteredRecords.length > 0"
        empty-title="暂无交易记录"
        padding="0"
        @load-more="loadMore"
      >
        <TransactionItem
          v-for="record in filteredRecords"
          :key="record.id"
          :record="record"
        />
      </InfiniteList>
    </GlCard>
    
    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import GlCard from '@/components/gl/Card/index.vue'
import GlTag from '@/components/gl/Tag/index.vue'
import InfiniteList from '@/components/InfiniteList/index.vue'
import QuickActions from '@/components/QuickActions/index.vue'
// Business 组件
import WalletBalanceCard from '@/components/WalletBalanceCard/index.vue'
import TransactionItem from '@/components/TransactionItem/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// Composables
import { useWallet } from '@/composables/useWallet'

const {
  loading,
  loadingMore,
  noMore,
  showBalance,
  currentFilter,
  wallet,
  filteredRecords,
  filterTabs,
  quickActions,
  vipDiscountText,
  loadMore,
  handleQuickAction,
  goBack,
  goToRecharge,
  goToWithdraw,
  goToVip,
  init,
} = useWallet()

onMounted(init)

onShow(() => {
  init()
})
</script>

<style lang="scss" scoped>
.wallet-page {
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

.vip-tip {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin: 0 24rpx 20rpx;
  padding: 20rpx 24rpx;
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.1) 0%, rgba(245, 158, 11, 0.05) 100%);
  border: 1rpx solid rgba(245, 158, 11, 0.3);
  border-radius: 16rpx;
}

.vip-text {
  flex: 1;
  font-size: 26rpx;
  color: var(--color-text);
}

.vip-more {
  font-size: 24rpx;
  color: #F59E0B;
}

.records-card {
  margin: 0 24rpx;
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.filter-tabs {
  display: flex;
  gap: 24rpx;
}

.filter-tab {
  font-size: 26rpx;
  color: var(--color-text-secondary);
  padding: 8rpx 0;
  position: relative;
  
  &.active {
    color: var(--color-primary);
    font-weight: 600;
    
    &::after {
      content: '';
      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;
      height: 4rpx;
      background: var(--color-primary);
      border-radius: 2rpx;
    }
  }
}
</style>
