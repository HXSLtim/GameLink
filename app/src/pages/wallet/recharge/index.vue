<template>
  <BasePageLayout
    class="recharge-page"
    padding="0"
    title="充值"
    :show-back="true"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
  >

    <view class="recharge-content">
      <!-- 当前余额 -->
      <RechargeBalanceInfo :amount="currentBalance" amount-unit="cents" />

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
      <AgreementCheckRow v-model="agreeTerms" link-text="充值服务协议" @link="viewAgreement" />

      <!-- 充值说明 -->
      <TipsListCard title="充值说明" :items="tipsList" />

      <!-- 底部占位 (移动端) -->
      <view class="bottom-placeholder"></view>
      
      <!-- PC 端将 Action Bar 移入流中 -->
      <RechargeActionBar
        v-if="isPC"
        class="pc-action-bar"
        :total="finalAmount"
        :bonus="bonusAmount"
        :disabled="!canSubmit"
        :loading="submitting"
        @submit="submitRecharge"
      />
    </view>



    <template #footer>
      <!-- 底部操作栏 (移动端) -->
      <RechargeActionBar
        v-if="!isPC"
        :total="finalAmount"
        :bonus="bonusAmount"
        :disabled="!canSubmit"
        :loading="submitting"
        @submit="submitRecharge"
      />
    </template>
  </BasePageLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
// Pattern 组件
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
// Business 组件
import AmountSelector from '@/components/AmountSelector/index.vue'
import PaymentMethodSelector from '@/components/PaymentMethodSelector/index.vue'
import RechargeBalanceInfo from '@/components/RechargeBalanceInfo/index.vue'
import AgreementCheckRow from '@/components/AgreementCheckRow/index.vue'
import TipsListCard from '@/components/TipsListCard/index.vue'
import RechargeActionBar from '@/components/RechargeActionBar/index.vue'
// Composables
import { useRecharge } from '@/composables/useRecharge'
import { useDevice } from '@/composables/useDevice'

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
  submitRecharge,
  viewAgreement,
  goBack,
  init,
} = useRecharge()

const { isPC } = useDevice()

const tipsList = [
  '1. 充值金额将实时到账，如遇延迟请联系客服',
  '2. 充值赠送金额仅限首次充值该档位',
  '3. 充值金额不可提现，仅用于平台消费',
  '4. 如有疑问请联系在线客服',
]

onMounted(init)
</script>

<style lang="scss" scoped>
.recharge-content {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  padding: var(--spacing-md);

  @include desktop {
    max-width: 600px;
    margin: 0 auto;
    padding: var(--spacing-xl) 0;
  }
}

.bottom-placeholder {
  height: 180rpx;

  @include desktop {
    display: none;
  }
}

// PC 端 Action Bar 样式覆盖
@include desktop {
  :deep(.action-bar) {
    position: static !important;
    width: 100%;
    margin-top: var(--spacing-lg);
    border: 1rpx solid var(--color-border);
    border-radius: var(--radius-lg);
    padding: var(--spacing-md);
    box-shadow: none;
  }
}
</style>
