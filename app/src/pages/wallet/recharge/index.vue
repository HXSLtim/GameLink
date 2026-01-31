<template>
  <view class="recharge-page page-container">
    <!-- 顶部导航 -->
    <NavBar title="充值" @back="goBack" />

    <scroll-view class="content-scroll" scroll-y>
      <!-- 当前余额 -->
      <view class="balance-info">
        <text class="balance-label">当前余额</text>
        <text class="balance-value">¥{{ formatBalance(currentBalance) }}</text>
      </view>

      <!-- 充值金额选择 -->
      <AmountSelector
        v-model="selectedAmount"
        v-model:custom-value="customAmount"
        :options="amountOptions"
      />

      <!-- 支付方式 -->
      <PaymentMethodSelector
        v-model="selectedMethod"
        :methods="paymentMethods"
      />

      <!-- 充值协议 -->
      <view class="agreement">
        <view class="checkbox" :class="{ checked: agreeTerms }" @tap="agreeTerms = !agreeTerms">
          <uv-icon v-if="agreeTerms" name="checkbox-mark" size="14" color="#fff"></uv-icon>
        </view>
        <text class="agreement-text">
          我已阅读并同意
          <text class="link" @tap="viewAgreement">《充值服务协议》</text>
        </text>
      </view>

      <!-- 充值说明 -->
      <GlCard title="充值说明" :shadow="false" bordered class="tips-card">
        <view class="tips-list">
          <text class="tips-item">1. 充值金额将实时到账，如遇延迟请联系客服</text>
          <text class="tips-item">2. 充值赠送金额仅限首次充值该档位</text>
          <text class="tips-item">3. 充值金额不可提现，仅用于平台消费</text>
          <text class="tips-item">4. 如有疑问请联系在线客服</text>
        </view>
      </GlCard>

      <view class="bottom-placeholder"></view>
    </scroll-view>

    <!-- 底部操作栏 -->
    <view class="action-bar">
      <view class="total-info">
        <text class="total-label">实付</text>
        <text class="total-value">¥{{ finalAmount.toFixed(2) }}</text>
        <text v-if="bonusAmount > 0" class="bonus-text">（含赠送 ¥{{ bonusAmount }}）</text>
      </view>
      <GlButton 
        type="primary" 
        size="large" 
        round 
        :disabled="!canSubmit"
        :loading="submitting"
        @click="submitRecharge"
      >
        立即充值
      </GlButton>
    </view>

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import GlCard from '@/components/gl/Card/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
// Business 组件
import AmountSelector from '@/components/AmountSelector/index.vue'
import PaymentMethodSelector from '@/components/PaymentMethodSelector/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// Composables
import { useRecharge } from '@/composables/useRecharge'

const {
  submitting,
  agreeTerms,
  currentBalance,
  selectedAmount,
  customAmount,
  selectedMethod,
  amountOptions,
  paymentMethods,
  finalAmount,
  bonusAmount,
  canSubmit,
  formatBalance,
  submitRecharge,
  viewAgreement,
  goBack,
  init,
} = useRecharge()

onMounted(init)
</script>

<style lang="scss" scoped>
.recharge-page {
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
  padding-bottom: 160rpx;
}

.balance-info {
  display: flex;
  align-items: baseline;
  gap: 16rpx;
  padding: 32rpx 24rpx;
}

.balance-label {
  font-size: 28rpx;
  color: var(--color-text-secondary);
}

.balance-value {
  font-size: 48rpx;
  font-weight: 800;
  color: var(--color-text);
}

:deep(.gl-card) {
  margin: 0 24rpx 20rpx;
}

.agreement {
  display: flex;
  align-items: center;
  gap: 12rpx;
  padding: 24rpx;
}

.checkbox {
  width: 36rpx;
  height: 36rpx;
  border-radius: 8rpx;
  border: 2rpx solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  
  &.checked {
    background: var(--color-primary);
    border-color: var(--color-primary);
  }
}

.agreement-text {
  font-size: 26rpx;
  color: var(--color-text-secondary);
}

.link {
  color: var(--color-primary);
}

.tips-list {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.tips-item {
  font-size: 24rpx;
  color: var(--color-text-secondary);
  line-height: 1.5;
}

.action-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx 24rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background: var(--color-bg-card);
  border-top: 1rpx solid var(--color-border);
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 100;
}

.total-info {
  display: flex;
  align-items: baseline;
  gap: 8rpx;
}

.total-label {
  font-size: 26rpx;
  color: var(--color-text-secondary);
}

.total-value {
  font-size: 40rpx;
  font-weight: 800;
  color: var(--color-primary);
}

.bonus-text {
  font-size: 22rpx;
  color: var(--color-text-secondary);
}

.bottom-placeholder {
  height: 180rpx;
}
</style>
