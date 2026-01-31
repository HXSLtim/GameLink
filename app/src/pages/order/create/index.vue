<template>
  <view class="order-create-page page-container">
    <!-- 顶部导航 -->
    <NavBar title="下单" @back="goBack" />

    <!-- 加载状态 -->
    <template v-if="loading">
      <view class="loading-wrap">
        <uv-skeleton rows="3" title loading></uv-skeleton>
      </view>
    </template>

    <!-- 内容区域 -->
    <scroll-view v-else class="content-scroll" scroll-y>
      <!-- 陪玩师信息 -->
      <OrderPlayerCard :player="player" />

      <!-- 选择游戏 -->
      <GameSelector v-model="selectedGameId" :games="player.games" />

      <!-- 选择服务 -->
      <ServiceSelector v-model="selectedServiceId" :services="player.services" />

      <!-- 预约时间 -->
      <SchedulePicker v-model:date="scheduledDate" v-model:time="scheduledTime" />

      <!-- 数量选择 -->
      <QuantitySelector v-model="quantity" :title="quantityTitle" :min="1" :max="10" :tip="quantityTip" />

      <!-- 游戏账号 -->
      <GlCard title="游戏账号" :shadow="false" bordered>
        <input v-model="gameAccount" class="text-input" placeholder="请输入您的游戏账号（选填）" />
      </GlCard>

      <!-- 备注信息 -->
      <GlCard title="备注信息" :shadow="false" bordered>
        <textarea v-model="remark" class="textarea-input" placeholder="请输入备注信息（选填）" :maxlength="200" />
        <text class="char-count">{{ remark.length }}/200</text>
      </GlCard>

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
    </scroll-view>

    <!-- 底部操作栏 -->
    <OrderSubmitBar :total="totalFee" :disabled="!canSubmit" :loading="submitting" @submit="submitOrder" />

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import GlCard from '@/components/gl/Card/index.vue'
// Business 组件
import OrderPlayerCard from '@/components/OrderPlayerCard/index.vue'
import GameSelector from '@/components/GameSelector/index.vue'
import ServiceSelector from '@/components/ServiceSelector/index.vue'
import SchedulePicker from '@/components/SchedulePicker/index.vue'
import QuantitySelector from '@/components/QuantitySelector/index.vue'
import CouponSelector from '@/components/CouponSelector/index.vue'
import OrderFeeSection from '@/components/OrderFeeSection/index.vue'
import OrderSubmitBar from '@/components/OrderSubmitBar/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
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
.order-create-page {
  height: 100vh;
  height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  box-sizing: border-box;
  overflow: hidden;
}

.loading-wrap {
  flex: 1;
  padding: 24rpx;
}

.content-scroll {
  flex: 1;
  padding: 24rpx;
  overflow-y: auto;
}

.text-input {
  width: 100%;
  font-size: 28rpx;
  color: var(--color-text);
}

.textarea-input {
  width: 100%;
  min-height: 120rpx;
  font-size: 28rpx;
  color: var(--color-text);
  box-sizing: border-box;
}

.char-count {
  display: block;
  text-align: right;
  font-size: 24rpx;
  color: var(--color-text-placeholder);
  margin-top: 8rpx;
}

.bottom-placeholder {
  height: 160rpx;
}
</style>
