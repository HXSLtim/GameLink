<template>
  <BasePageLayout
    class="order-create-page"
    :scroll="!loading"
    title="下单"
    :show-back="true"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
  >
    <template #nav>
      <!-- 顶部导航 -->
      <NavBar title="下单" @back="goBack" />
    </template>

    <!-- 加载状态 -->
    <template v-if="loading">
      <view class="loading-wrap">
        <Skeleton :rows="3" />
      </view>
    </template>

    <!-- 内容区域 -->
    <template v-else>
      <!-- 陪玩师信息 -->
      <PlayerCard
        class="order-player-card"
        :player="player"
        variant="compact"
        :show-online-tag="true"
        :show-offline-tag="true"
        :clickable="false"
      />

      <!-- 选择游戏 -->
      <GameSelector v-model="selectedGameId" :games="player.games" />

      <!-- 选择服务 -->
      <ServiceSelector v-model="selectedServiceId" :services="player.services" />

      <!-- 预约时间 -->
      <SchedulePicker v-model:date="scheduledDate" v-model:time="scheduledTime" />

      <!-- 数量选择 -->
      <QuantitySelector v-model="quantity" :title="quantityTitle" :min="1" :max="10" :tip="quantityTip" />

      <!-- 游戏账号 -->
      <SectionCard title="游戏账号" margin="var(--spacing-sm) var(--spacing-md)">
        <GlInput
          v-model="gameAccount"
          class="text-input"
          placeholder="请输入您的游戏账号（选填）"
          size="medium"
        />
      </SectionCard>

      <!-- 备注信息 -->
      <SectionCard title="备注信息" margin="var(--spacing-sm) var(--spacing-md)">
        <GlInput
          v-model="remark"
          class="textarea-input"
          type="textarea"
          size="medium"
          placeholder="请输入备注信息（选填）"
          :maxlength="200"
        />
        <text class="char-count">{{ remark.length }}/200</text>
      </SectionCard>

      <!-- 优惠券 -->
      <CouponSelector
        :selected-coupon="selectedCoupon"
        :available-count="availableCoupons.length"
        @click="showCouponPicker = true"
      />

      <!-- 费用明细 -->
      <OrderFeeSection
        :service-fee="serviceFee"
        :coupon-discount="selectedCoupon?.discount"
        :vip-discount="vipDiscount"
        :total="totalFee"
      />

      <!-- 底部占位 -->
      <view class="bottom-placeholder"></view>
    </template>

    <template #footer>
      <!-- 底部操作栏 -->
      <OrderSubmitBar :total="totalFee" :disabled="!canSubmit" :loading="submitting" @submit="submitOrder" />
    </template>
  </BasePageLayout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
import SectionCard from '@/components/SectionCard/index.vue'
import GlInput from '@/components/gl/Input/index.vue'
import Skeleton from '@/components/Skeleton/index.vue'
// Business 组件
import PlayerCard from '@/components/PlayerCard/index.vue'
import GameSelector from '@/components/GameSelector/index.vue'
import ServiceSelector from '@/components/ServiceSelector/index.vue'
import SchedulePicker from '@/components/SchedulePicker/index.vue'
import QuantitySelector from '@/components/QuantitySelector/index.vue'
import CouponSelector from '@/components/CouponSelector/index.vue'
import OrderFeeSection from '@/components/OrderFeeSection/index.vue'
import OrderSubmitBar from '@/components/OrderSubmitBar/index.vue'
// Composables
import { useOrderCreate } from '@/composables/useOrderCreate'

const {
  loading,
  submitting,
  showCouponPicker,
  player,
  selectedGameId,
  selectedServiceId,
  scheduledDate,
  scheduledTime,
  quantity,
  gameAccount,
  remark,
  selectedCoupon,
  availableCoupons,
  serviceFee,
  vipDiscount,
  totalFee,
  canSubmit,
  loadPlayerInfo,
  submitOrder,
  goBack,
} = useOrderCreate()

const quantityTitle = computed(() => {
  const service = player.value.services?.find((s: any) => s.id === selectedServiceId.value)
  return service?.unit === 'hour' ? '服务时长' : '局数'
})

const quantityTip = computed(() => {
  const service = player.value.services?.find((s: any) => s.id === selectedServiceId.value)
  return service?.unit === 'hour' ? '小时' : '局'
})

onLoad((options) => {
  const playerId = Number(options?.playerId)
  if (playerId) {
    loadPlayerInfo(playerId)
  }
})
</script>

<style lang="scss" scoped>
.loading-wrap {
  flex: 1;
  padding: var(--spacing-md);
}

.order-player-card {
  margin-bottom: var(--spacing-md);
}

.text-input {
  width: 100%;
  
  :deep(.gl-input__field) {
    font-size: var(--font-md);
    color: var(--color-text);
  }
}

.textarea-input {
  width: 100%;
  
  :deep(.gl-input__textarea) {
    min-height: 120rpx;
    font-size: var(--font-md);
    color: var(--color-text);
  }
}

.char-count {
  display: block;
  text-align: right;
  font-size: var(--font-xs);
  color: var(--color-text-placeholder);
  margin-top: var(--spacing-sm);
}

.bottom-placeholder {
  height: 160rpx;
}
</style>
