<template>
  <BasePageLayout
    class="wallet-page"
    :scroll="false"
    padding="0"
    title="我的钱包"
    :show-back="true"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
  >

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
    <WalletVipTip :vip-level="wallet.vipLevel" :discount-text="vipDiscountText" @click="goToVip" />

    <!-- 快捷入口 -->
    <QuickActions :items="quickActions" @click="handleQuickAction" />

    <!-- 交易记录 -->
    <WalletRecordsSection
      v-model="currentFilter"
      :tabs="filterTabs"
      :loading="loading"
      :loading-more="loadingMore"
      :no-more="noMore"
      :records="filteredRecords"
      @load-more="loadMore"
    />
    
  </BasePageLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
import QuickActions from '@/components/QuickActions/index.vue'
// Business 组件
import WalletBalanceCard from '@/components/WalletBalanceCard/index.vue'
import WalletVipTip from '@/components/WalletVipTip/index.vue'
import WalletRecordsSection from '@/components/WalletRecordsSection/index.vue'
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
</style>
