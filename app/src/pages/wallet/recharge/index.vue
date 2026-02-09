<template>
  <BasePageLayout
    class="recharge-page"
    padding="0"
    title="充值"
    :show-back="true"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
  >

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

    <view class="bottom-placeholder"></view>

    <template #footer>
      <!-- 底部操作栏 -->
      <RechargeActionBar
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

const tipsList = [
  '1. 充值金额将实时到账，如遇延迟请联系客服',
  '2. 充值赠送金额仅限首次充值该档位',
  '3. 充值金额不可提现，仅用于平台消费',
  '4. 如有疑问请联系在线客服',
]

onMounted(init)
</script>

<style lang="scss" scoped>
.bottom-placeholder {
  height: 180rpx;
}
</style>
